package analysis

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

const (
	maxTableNameLengthAfter12_2  = 128
	maxTableNameLengthBefore12_2 = 30
)

// ErrNoBackupableDML reports a script with no UPDATE or DELETE to back up.
var ErrNoBackupableDML = errors.New("no parse results")

// Version is an Oracle release, such as 12.2.
type Version struct {
	First  int
	Second int
}

// GTE returns true if the version is greater than or equal to the base version.
func (v *Version) GTE(base *Version) bool {
	if v.First > base.First {
		return true
	}
	if v.First == base.First {
		return v.Second >= base.Second
	}
	return false
}

// BackupOptions is what TransformDMLToSelect needs besides the script.
type BackupOptions struct {
	// SourceDatabase is the schema unqualified tables are in.
	SourceDatabase string
	// TargetDatabase is the schema the backup tables are created in.
	TargetDatabase string
	// TablePrefix starts each backup table's name.
	TablePrefix string
	// Version is the server's release, which caps the backup table name at 128
	// characters from 12.2 on and at 30 before. A nil Version is 11.0.
	Version *Version
}

// BackupStatement is the CTAS that backs up the rows one table's statements
// change.
type BackupStatement struct {
	Statement       string
	SourceSchema    string
	SourceTableName string
	TargetTableName string
	// Start is the byte offset in the script where the first statement on the
	// table starts, and End the offset of the last byte of the last one.
	Start int
	End   int
}

type statementType int

const (
	statementTypeUnknown statementType = iota
	statementTypeUpdate
	statementTypeDelete
)

type tableReference struct {
	Database      string
	HasSchema     bool
	Schema        string
	Table         string
	Alias         string
	DBLink        string
	StatementType statementType
}

type statementInfo struct {
	offset    int
	statement string
	node      ast.StmtNode
	table     *tableReference
	start     int
	end       int
	fullSQL   string
}

// TransformDMLToSelect turns a script's UPDATE and DELETE statements into one
// CTAS per table that copies the rows they change, ordered by where each
// table's statements start. Statements on one table must be all UPDATE or all
// DELETE, and a table reached through a database link is refused.
func TransformDMLToSelect(statement string, opts BackupOptions) ([]BackupStatement, error) {
	statementInfoList, err := prepareTransformation(opts.SourceDatabase, statement)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(opts, statementInfoList, opts.TargetDatabase, opts.TablePrefix)
}

// ReferencesTable reports whether a script backs up table in database: it has
// an UPDATE or DELETE of it. The error is a parse error.
func ReferencesTable(statement, database, table string) (bool, error) {
	infos, err := prepareTransformation(database, statement)
	if err != nil {
		if err == ErrNoBackupableDML {
			return false, nil
		}
		return false, err
	}
	for _, info := range infos {
		if info.table.Schema == database && info.table.Table == table {
			return true, nil
		}
	}
	return false, nil
}

func generateSQL(opts BackupOptions, statementInfoList []statementInfo, targetDatabase string, tablePrefix string) ([]BackupStatement, error) {
	// Tables are visited in the order the script first changes them, so an
	// error names the same table on every run.
	groupByTable := make(map[string][]statementInfo)
	var tableOrder []string
	for _, item := range statementInfoList {
		key := fmt.Sprintf("%s.%s@%s", item.table.Schema, item.table.Table, item.table.DBLink)
		if _, ok := groupByTable[key]; !ok {
			tableOrder = append(tableOrder, key)
		}
		groupByTable[key] = append(groupByTable[key], item)
	}

	// Check if the statement type is the same for all statements in the group.
	for _, key := range tableOrder {
		list := groupByTable[key]
		statementType := statementTypeUnknown
		for _, item := range list {
			if statementType == statementTypeUnknown {
				statementType = item.table.StatementType
			}
			if statementType != item.table.StatementType {
				return nil, fmt.Errorf("prior backup cannot handle statements with different types on the same table: %s", key)
			}
		}
	}

	var result []BackupStatement
	for _, key := range tableOrder {
		list := groupByTable[key]
		backupStatement, err := generateSQLForTable(opts, list, targetDatabase, tablePrefix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL for table: %s: %w", key, err)
		}
		result = append(result, *backupStatement)
	}

	slices.SortFunc(result, func(i, j BackupStatement) int {
		if i.Start != j.Start {
			if i.Start < j.Start {
				return -1
			}
			return 1
		}
		if i.SourceTableName < j.SourceTableName {
			return -1
		}
		if i.SourceTableName > j.SourceTableName {
			return 1
		}
		return 0
	})

	return result, nil
}

func generateSQLForTable(opts BackupOptions, statementInfoList []statementInfo, targetDatabase string, tablePrefix string) (*BackupStatement, error) {
	table := statementInfoList[0].table

	// Prior backup cannot round-trip a table referenced through a database
	// link: the backup CTAS could read the remote rows, but the restore path
	// only records schema+table and would write the rows back into a LOCAL
	// table of the same name — silent cross-database corruption on rollback.
	// That was already true before the ROWID rework (the link was dropped
	// from the reconstructed source), so fail loudly instead of producing a
	// backup that cannot be restored correctly.
	if table.DBLink != "" {
		return nil, fmt.Errorf("prior backup does not support tables referenced through a database link: %s.%s@%s; disable prior backup for this issue", table.Schema, table.Table, table.DBLink)
	}

	version := opts.Version
	if version == nil {
		version = &Version{
			First:  11,
			Second: 0,
		}
	}

	targetTable := fmt.Sprintf("%s_%s_%s", tablePrefix, table.Table, table.Schema)
	if version.GTE(&Version{First: 12, Second: 2}) {
		targetTable, _ = truncateString(targetTable, maxTableNameLengthAfter12_2)
	} else {
		targetTable, _ = truncateString(targetTable, maxTableNameLengthBefore12_2)
	}

	var buf strings.Builder
	if _, err := fmt.Fprintf(&buf, `CREATE TABLE "%s"."%s" AS`+"\n", targetDatabase, targetTable); err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}
	if len(statementInfoList) == 1 && !joinsTables(statementInfoList[0].node) {
		// Single DML: a plain CTAS needs no deduplication.
		if err := writeBranchSelect(&buf, statementInfoList[0], "*", "  "); err != nil {
			return nil, err
		}
	} else {
		// Multiple DMLs on the same table: deduplicate by ROWID instead of a
		// value-based UNION. UNION's implicit DISTINCT compares every column,
		// which fails on LOB columns (ORA-22848) and silently collapses
		// value-identical physical rows into one backup row. ROWID is always
		// comparable, keeps LOB columns out of the comparison key, and
		// deduplicates by physical row — the semantics a backup wants.
		if _, err := fmt.Fprintf(&buf, `  SELECT "%s"."%s".* FROM "%s"."%s" WHERE ROWID IN (`+"\n",
			table.Schema, table.Table, table.Schema, table.Table); err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
		for i, info := range statementInfoList {
			if i != 0 {
				if _, err := buf.WriteString("\n    UNION\n"); err != nil {
					return nil, fmt.Errorf("failed to write to buffer: %w", err)
				}
			}
			if err := writeBranchSelect(&buf, info, "ROWID", "    "); err != nil {
				return nil, err
			}
		}
		if _, err := buf.WriteString("\n  )"); err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}

	if _, err := buf.WriteString(";"); err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	return &BackupStatement{
		Statement:       buf.String(),
		SourceSchema:    table.Schema,
		SourceTableName: table.Table,
		TargetTableName: targetTable,
		Start:           statementInfoList[0].start,
		End:             statementInfoList[len(statementInfoList)-1].end,
	}, nil
}

func prepareTransformation(databaseName, statement string) ([]statementInfo, error) {
	statements := SplitStatements(statement)

	var result []statementInfo
	for i, stmt := range statements {
		if stmt.Empty {
			continue
		}
		list, err := Parse(stmt.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PLSQL: %w", err)
		}
		for _, item := range list.Items {
			raw, ok := item.(*ast.RawStmt)
			if !ok || raw.Stmt == nil {
				continue
			}
			info := extractDML(databaseName, raw.Stmt, stmt.Text)
			if info == nil {
				continue
			}
			info.offset = i
			info.start = stmt.Start + raw.Loc.Start
			info.end = stmt.Start + max(raw.Loc.End-1, raw.Loc.Start)
			info.fullSQL = stmt.Text
			result = append(result, *info)
		}
	}

	if len(result) == 0 {
		return nil, ErrNoBackupableDML
	}
	return result, nil
}

func extractDML(databaseName string, node ast.StmtNode, fullSQL string) *statementInfo {
	switch n := node.(type) {
	case *ast.DeleteStmt:
		table := dMLTableReference(databaseName, n.Table, n.Alias, statementTypeDelete)
		if table == nil {
			return nil
		}
		return &statementInfo{
			statement: extractStatementText(fullSQL, n.Loc),
			node:      n,
			table:     table,
		}
	case *ast.UpdateStmt:
		table := dMLTableReference(databaseName, n.Table, n.Alias, statementTypeUpdate)
		if table == nil {
			return nil
		}
		return &statementInfo{
			statement: extractStatementText(fullSQL, n.Loc),
			node:      n,
			table:     table,
		}
	default:
		return nil
	}
}

func dMLTableReference(databaseName string, name *ast.ObjectName, alias *ast.Alias, statementType statementType) *tableReference {
	if name == nil || name.Name == "" {
		return nil
	}
	schema := name.Schema
	hasSchema := schema != ""
	if schema == "" {
		schema = databaseName
	}
	table := &tableReference{
		Database:      schema,
		HasSchema:     hasSchema,
		Schema:        schema,
		Table:         name.Name,
		DBLink:        name.DBLink,
		StatementType: statementType,
	}
	if alias != nil {
		table.Alias = alias.Name
	}
	return table
}

// writeBranchSelect writes one `SELECT <qualifier>.<projection> FROM <dml
// target and predicate>` branch. projection is "*" for a plain CTAS body and
// "ROWID" for the deduplication subquery branches.
// joinsTables reports whether a DML statement joins its target to the
// tables of a FROM list, so that a target row may match more than once and
// its backup needs the ROWID deduplication.
func joinsTables(node ast.StmtNode) bool {
	switch n := node.(type) {
	case *ast.DeleteStmt:
		return n.FromClause != nil && n.FromClause.Len() > 0
	case *ast.UpdateStmt:
		return n.FromClause != nil && n.FromClause.Len() > 0 && n.WhereClause != nil
	}
	return false
}

func writeBranchSelect(buf *strings.Builder, info statementInfo, projection string, indent string) error {
	t := info.table
	switch {
	case t.Alias != "":
		if _, err := fmt.Fprintf(buf, `%sSELECT "%s".%s FROM `, indent, t.Alias, projection); err != nil {
			return fmt.Errorf("failed to write to buffer: %w", err)
		}
	case t.HasSchema:
		if _, err := fmt.Fprintf(buf, `%sSELECT "%s"."%s".%s FROM `, indent, t.Schema, t.Table, projection); err != nil {
			return fmt.Errorf("failed to write to buffer: %w", err)
		}
	default:
		if _, err := fmt.Fprintf(buf, `%sSELECT "%s".%s FROM `, indent, t.Table, projection); err != nil {
			return fmt.Errorf("failed to write to buffer: %w", err)
		}
	}
	if err := writeSuffixSelectClause(buf, info.node, info.fullSQL); err != nil {
		return fmt.Errorf("failed to write suffix select clause: %w", err)
	}
	return nil
}

func writeSuffixSelectClause(buf *strings.Builder, node ast.StmtNode, fullSQL string) error {
	suffix := ""
	switch n := node.(type) {
	case *ast.UpdateStmt:
		suffix = oracleDMLTargetText(fullSQL, n.Table, n.PartitionExt, n.Alias)
		if n.WhereClause != nil && n.SetClauses != nil && n.SetClauses.Len() > 0 {
			if setClause, ok := n.SetClauses.Items[n.SetClauses.Len()-1].(*ast.SetClause); ok {
				whereStart := setClause.Loc.End
				if fromText, fromEnd, ok := oracleListText(fullSQL, n.FromClause); ok {
					suffix += ", " + fromText
					whereStart = fromEnd
				}
				suffix = joinOracleSuffix(suffix, oracleWhereSuffix(fullSQL, whereStart, n.WhereClause))
			}
		}
	case *ast.DeleteStmt:
		targetEnd := n.Table.Loc.End
		if n.PartitionExt != nil {
			targetEnd = n.PartitionExt.Loc.End
		}
		if n.Alias != nil {
			targetEnd = n.Alias.Loc.End
		}
		suffix = oracleDMLTargetText(fullSQL, n.Table, n.PartitionExt, n.Alias)
		// DELETE t FROM u ...: the FROM's tables join the target's.
		if fromText, fromEnd, ok := oracleListText(fullSQL, n.FromClause); ok {
			suffix += ", " + fromText
			targetEnd = fromEnd
		}
		suffix = joinOracleSuffix(suffix, oracleWhereSuffix(fullSQL, targetEnd, n.WhereClause))
	default:
	}
	_, err := buf.WriteString(suffix)
	return err
}

func oracleDMLTargetText(sql string, name *ast.ObjectName, partitionExt *ast.PartitionExtClause, alias *ast.Alias) string {
	if name == nil {
		return ""
	}
	end := name.Loc.End
	if partitionExt != nil {
		end = partitionExt.Loc.End
	}
	if alias != nil {
		end = alias.Loc.End
	}
	return strings.TrimSpace(oracleLocText(sql, name.Loc.Start, end))
}

func oracleDMLTrailingText(sql string, start, end int) string {
	text := strings.TrimSpace(oracleLocText(sql, start, end))
	return strings.TrimSpace(strings.TrimSuffix(text, ";"))
}

func oracleListText(sql string, list *ast.List) (string, int, bool) {
	if list == nil || list.Len() == 0 {
		return "", 0, false
	}
	start, ok := oracleNodeLoc(list.Items[0])
	if !ok {
		return "", 0, false
	}
	end, ok := oracleNodeLoc(list.Items[list.Len()-1])
	if !ok {
		return "", 0, false
	}
	text := strings.TrimSpace(oracleLocText(sql, start.Start, end.End))
	if text == "" {
		return "", 0, false
	}
	return text, end.End, true
}

func oracleWhereSuffix(sql string, start int, where ast.ExprNode) string {
	loc, ok := oracleNodeLoc(where)
	if !ok {
		return ""
	}
	return oracleDMLTrailingText(sql, start, loc.End)
}

func oracleNodeLoc(node ast.Node) (ast.Loc, bool) {
	if node == nil {
		return ast.Loc{}, false
	}
	value := reflect.ValueOf(node)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return ast.Loc{}, false
	}
	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return ast.Loc{}, false
	}
	field := elem.FieldByName("Loc")
	if !field.IsValid() || field.Type() != reflect.TypeOf(ast.Loc{}) {
		return ast.Loc{}, false
	}
	loc, ok := field.Interface().(ast.Loc)
	if !ok || loc.End <= loc.Start {
		return ast.Loc{}, false
	}
	return loc, true
}

func oracleLocText(sql string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end <= 0 || end > len(sql) {
		end = len(sql)
	}
	if start >= end || start >= len(sql) {
		return ""
	}
	return sql[start:end]
}

func joinOracleSuffix(parts ...string) string {
	var nonEmpty []string
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, " ")
}

func extractStatementText(sql string, loc ast.Loc) string {
	return strings.TrimSpace(strings.TrimSuffix(oracleLocText(sql, loc.Start, loc.End), ";"))
}

// truncateString keeps the first limit characters of str.
func truncateString(str string, limit int) (string, bool) {
	chars := 0
	for i := range str {
		if chars >= limit {
			return str[:i], true
		}
		chars++
	}
	return str, false
}
