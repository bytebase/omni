package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
)

const maxTableNameLength = 64

// DMLType is the change a backed-up statement makes to its table.
type DMLType int

const (
	DMLTypeUnknown DMLType = iota
	DMLTypeUpdate
	DMLTypeInsert
	DMLTypeDelete
)

// TableReference is a table a statement names. Database is the statement's
// default database when the reference is unqualified.
type TableReference struct {
	Database      string
	Table         string
	Alias         string
	StatementType DMLType
	// qualified is true when the reference names its database. Only an
	// unqualified reference can name a CTE.
	qualified bool
}

// isCTE reports whether ref names one of the statement's CTEs rather than a
// table.
func (ref *TableReference) isCTE(cteNames map[string]bool) bool {
	return !ref.qualified && cteNames[ref.Table]
}

// DMLTarget is a table an UPDATE or DELETE changes, with the statement that
// changes it.
type DMLTarget struct {
	// Offset is the index of the statement's split segment.
	Offset    int
	Statement string
	Node      ast.Node
	Table     *TableReference
	// FullSQL is the text Node's locations are relative to.
	FullSQL string
	// ByteStart and ByteEnd are the statement's split segment in the script.
	ByteStart int
	ByteEnd   int
}

// BackupOptions configures TransformDMLToSelect.
type BackupOptions struct {
	// DefaultDatabase is the database of unqualified tables.
	DefaultDatabase string
	// BackupDatabase is where the backup tables are created.
	BackupDatabase string
	// TablePrefix starts each backup table's name.
	TablePrefix string
	// CaseSensitive makes database and table lookups case-sensitive.
	CaseSensitive bool
	// Catalog is the instance's metadata.
	Catalog Catalog
	// Columns returns a table's generated columns and its other columns. Nil
	// means ClassifyColumns over Catalog.
	Columns func(database, table string) (generated, normal []string, err error)
}

// BackupStatement backs up the rows one table loses or changes.
type BackupStatement struct {
	// Statement creates the backup table and copies the rows into it.
	Statement   string
	SourceTable string
	BackupTable string
	// ByteStart is where the first split segment the backup covers starts, and
	// ByteEnd where the last one ends.
	ByteStart int
	ByteEnd   int
}

// TransformDMLToSelect returns, for each table the UPDATE and DELETE
// statements of a script change, a statement that copies the rows they touch
// into a backup table, ordered by position and then table name.
func TransformDMLToSelect(statement string, opts BackupOptions) ([]BackupStatement, error) {
	var database Database
	if opts.Catalog.GetDatabase != nil {
		if db, err := opts.Catalog.GetDatabase(opts.DefaultDatabase); err == nil {
			database = db
		}
	}

	targets, err := prepareTransformation(opts.DefaultDatabase, statement, database)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(targets, opts)
}

func prepareTransformation(databaseName, statement string, database Database) ([]DMLTarget, error) {
	var result []DMLTarget
	for i, segment := range parser.Split(statement) {
		text := segmentText(statement, segment)
		if len(text) == 0 || segment.Empty() {
			continue
		}

		parsed, err := parser.Parse(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sql: %w", err)
		}
		if parsed == nil || len(parsed.Items) == 0 {
			continue
		}

		for _, node := range parsed.Items {
			tables, err := ExtractTables(databaseName, node, text, i, database)
			if err != nil {
				return nil, fmt.Errorf("failed to extract tables: %w", err)
			}
			for _, table := range tables {
				table.ByteStart = segment.ByteStart
				table.ByteEnd = segment.ByteEnd
				result = append(result, table)
			}
		}
	}

	return result, nil
}

// segmentText is the segment's text with its ';', if any.
func segmentText(sql string, segment parser.Segment) string {
	end := segment.ByteEnd
	if end < len(sql) && sql[end] == ';' {
		end++
	}
	return sql[segment.ByteStart:end]
}

// ExtractTables returns the tables a DELETE or UPDATE changes. database, which
// may be nil, is the default database's metadata, used to find the table an
// unqualified SET column belongs to.
func ExtractTables(databaseName string, node ast.Node, fullSQL string, offset int, database Database) ([]DMLTarget, error) {
	switch n := node.(type) {
	case *ast.DeleteStmt:
		return extractTablesFromDelete(databaseName, n, fullSQL, offset)
	case *ast.UpdateStmt:
		return extractTablesFromUpdate(databaseName, n, fullSQL, offset, database)
	default:
		return nil, nil
	}
}

func extractTablesFromDelete(databaseName string, n *ast.DeleteStmt, fullSQL string, offset int) ([]DMLTarget, error) {
	cteNames := collectCTENames(fullSQL, n.Loc)
	stmtText := extractStatementText(fullSQL, n.Loc)

	singleTables := collectSingleTables(databaseName, n.Tables)

	if len(n.Using) == 0 {
		// Single-table DELETE: DELETE FROM t WHERE ...
		if len(singleTables) == 0 {
			return nil, nil
		}
		for _, ref := range singleTables {
			ref.StatementType = DMLTypeDelete
		}
		first := singleTables[firstKey(singleTables)]
		return []DMLTarget{{
			Offset:    offset,
			Statement: stmtText,
			Node:      n,
			Table:     first,
			FullSQL:   fullSQL,
		}}, nil
	}

	// Multi-table DELETE: DELETE t1, t2 FROM t1 JOIN t2 ... WHERE ...
	// Tables = targets to delete from, Using = referenced table list.
	refTables := collectSingleTables(databaseName, n.Using)

	var result []DMLTarget
	for _, target := range singleTables {
		name := target.Table
		if target.Alias != "" {
			name = target.Alias
		}
		ref, ok := refTables[name]
		if !ok {
			if cteNames[name] {
				continue
			}
			return nil, fmt.Errorf("cannot extract reference table: no matched table %q in referenced table list", name)
		}
		// A table whose alias is a CTE's name is still a target; only a
		// reference to the CTE itself is not.
		if ref.isCTE(cteNames) {
			continue
		}
		ref.StatementType = DMLTypeDelete
		result = append(result, DMLTarget{
			Offset:    offset,
			Statement: stmtText,
			Node:      n,
			Table:     ref,
			FullSQL:   fullSQL,
		})
	}
	return result, nil
}

func extractTablesFromUpdate(databaseName string, n *ast.UpdateStmt, fullSQL string, offset int, database Database) ([]DMLTarget, error) {
	cteNames := collectCTENames(fullSQL, n.Loc)
	stmtText := extractStatementText(fullSQL, n.Loc)
	singleTables := collectSingleTables(databaseName, n.Tables)

	// Determine which tables are updated via SET clause column prefixes.
	updatedTables := make(map[string]bool)
	var unqualifiedColumns []string
	for _, assign := range n.SetList {
		if assign.Column != nil {
			table := assign.Column.Table
			ref, ok := singleTables[table]
			if (ok && !ref.isCTE(cteNames)) || (!ok && !cteNames[table]) {
				updatedTables[table] = true
				if table == "" {
					unqualifiedColumns = append(unqualifiedColumns, assign.Column.Column)
				}
			}
		}
	}

	// Resolve unqualified columns using metadata.
	if updatedTables[""] {
		delete(updatedTables, "")
		resolved := resolveUnqualifiedColumns(unqualifiedColumns, singleTables, database)
		for t := range resolved {
			updatedTables[t] = true
		}
	}

	// Single-table UPDATE (only one table, no explicit qualification needed).
	if len(updatedTables) == 0 && len(singleTables) == 1 {
		for _, ref := range singleTables {
			ref.StatementType = DMLTypeUpdate
			return []DMLTarget{{
				Offset:    offset,
				Statement: stmtText,
				Node:      n,
				Table:     ref,
				FullSQL:   fullSQL,
			}}, nil
		}
	}

	var result []DMLTarget
	for table := range updatedTables {
		ref, ok := singleTables[table]
		if !ok {
			return nil, fmt.Errorf("cannot extract reference table: no matched updated table %q in referenced table list", table)
		}
		ref.StatementType = DMLTypeUpdate
		result = append(result, DMLTarget{
			Offset:    offset,
			Statement: stmtText,
			Node:      n,
			Table:     ref,
			FullSQL:   fullSQL,
		})
	}
	return result, nil
}

// collectSingleTables walks TableExpr slices and collects TableRef entries keyed by alias or name.
func collectSingleTables(databaseName string, exprs []ast.TableExpr) map[string]*TableReference {
	result := make(map[string]*TableReference)
	for _, expr := range exprs {
		collectSingleTablesFromExpr(databaseName, expr, result)
	}
	return result
}

func collectSingleTablesFromExpr(databaseName string, expr ast.TableExpr, out map[string]*TableReference) {
	switch e := expr.(type) {
	case *ast.TableRef:
		db := e.Schema
		if db == "" {
			db = databaseName
		}
		ref := &TableReference{
			Database:  db,
			Table:     e.Name,
			Alias:     e.Alias,
			qualified: e.Schema != "",
		}
		key := e.Name
		if e.Alias != "" {
			key = e.Alias
		}
		out[key] = ref
	case *ast.JoinClause:
		if e.Left != nil {
			collectSingleTablesFromExpr(databaseName, e.Left, out)
		}
		if e.Right != nil {
			collectSingleTablesFromExpr(databaseName, e.Right, out)
		}
	default:
	}
}

// collectCTENames returns the names of the CTEs the statement at stmtLoc opens with.
func collectCTENames(fullSQL string, stmtLoc ast.Loc) map[string]bool {
	result := make(map[string]bool)
	for _, cte := range parseCTEs(withClause(fullSQL, stmtLoc)) {
		result[cte.Name] = true
	}
	return result
}

// extractCTE returns the WITH clause the statement at stmtLoc opens with, or ""
// when it has none or the clause does not parse.
func extractCTE(fullSQL string, stmtLoc ast.Loc) string {
	clause := withClause(fullSQL, stmtLoc)
	if len(parseCTEs(clause)) == 0 {
		return ""
	}
	return clause
}

// withClause returns the WITH clause that opens the UPDATE or DELETE at loc, or "".
func withClause(sql string, loc ast.Loc) string {
	start := dmlKeywordStart(sql, loc)
	if start <= loc.Start {
		return ""
	}
	return strings.TrimSpace(sql[loc.Start:start])
}

// parseCTEs returns the CTEs of a WITH clause, or nil when it does not parse.
func parseCTEs(clause string) []*ast.CommonTableExpr {
	if clause == "" {
		return nil
	}
	// Wrap the clause with a dummy SELECT to make it parseable.
	parsed, err := parser.Parse(clause + " SELECT 1")
	if err != nil || parsed == nil {
		return nil
	}
	for _, item := range parsed.Items {
		if sel, ok := item.(*ast.SelectStmt); ok && len(sel.CTEs) > 0 {
			return sel.CTEs
		}
	}
	return nil
}

func extractStatementText(fullSQL string, loc ast.Loc) string {
	end := loc.End
	if end <= 0 || end > len(fullSQL) {
		end = len(fullSQL)
	}
	start := 0
	if loc.Start > 0 {
		start = loc.Start
	}
	return strings.TrimSpace(fullSQL[start:end])
}

// writeSuffixSelectClause writes the FROM ... WHERE ... suffix of a DML statement.
func writeSuffixSelectClause(buf *strings.Builder, node ast.Node, fullSQL string) error {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		writeUpdateSuffix(buf, n, fullSQL)
	case *ast.DeleteStmt:
		return writeDeleteSuffix(buf, n, fullSQL)
	}
	return nil
}

func writeUpdateSuffix(buf *strings.Builder, n *ast.UpdateStmt, sql string) {
	// For UPDATE: write "table_refs WHERE ... ORDER BY ... LIMIT ..."
	// Use the last table's Loc.End as the boundary — everything from first table
	// to last table's end is the table reference list.
	if len(n.Tables) > 0 {
		start := nodeLocStart(n.Tables[0])
		end := nodeLocEnd(n.Tables[len(n.Tables)-1])
		if start >= 0 && end > start && end <= len(sql) {
			buf.WriteString(strings.TrimSpace(sql[start:end]))
		}
	}

	// Everything after the last SET assignment is WHERE + ORDER BY + LIMIT.
	if len(n.SetList) > 0 {
		lastAssign := n.SetList[len(n.SetList)-1]
		afterSet := lastAssign.Loc.End
		stmtEnd := n.Loc.End
		if stmtEnd <= 0 || stmtEnd > len(sql) {
			stmtEnd = len(sql)
		}
		suffix := strings.TrimSpace(sql[afterSet:stmtEnd])
		if suffix != "" {
			fmt.Fprintf(buf, " %s", suffix)
		}
	}
}

func writeDeleteSuffix(buf *strings.Builder, n *ast.DeleteStmt, sql string) error {
	// For single-table DELETE: everything from the table ref to statement end.
	// For multi-table DELETE: everything from the USING/reference table list to end.
	end, err := endBeforeReturning(n.Returning, sql, n.Loc.End)
	if err != nil {
		return err
	}
	var start int
	switch {
	case len(n.Using) > 0:
		start = nodeLocStart(n.Using[0])
	case len(n.Tables) > 0:
		start = nodeLocStart(n.Tables[0])
	default:
		return nil
	}
	if start < 0 || end <= start {
		return nil
	}
	buf.WriteString(strings.TrimSpace(sql[start:end]))
	return nil
}

func nodeLocStart(expr ast.TableExpr) int {
	switch e := expr.(type) {
	case *ast.TableRef:
		return e.Loc.Start
	case *ast.JoinClause:
		return e.Loc.Start
	}
	return -1
}

func nodeLocEnd(expr ast.TableExpr) int {
	switch e := expr.(type) {
	case *ast.TableRef:
		return e.Loc.End
	case *ast.JoinClause:
		return e.Loc.End
	}
	return -1
}

// resolveUnqualifiedColumns resolves unqualified column names to owning tables using metadata.
func resolveUnqualifiedColumns(columns []string, singleTables map[string]*TableReference, database Database) map[string]bool {
	result := make(map[string]bool)

	if database == nil || !database.HasSchema() {
		for t := range singleTables {
			result[t] = true
		}
		return result
	}

	for _, col := range columns {
		resolved := false
		for aliasOrName, ref := range singleTables {
			tableMetadata := database.GetTable(ref.Table)
			if tableMetadata == nil {
				continue
			}
			if hasColumn(tableMetadata, col) {
				result[aliasOrName] = true
				resolved = true
			}
		}
		if !resolved {
			for t := range singleTables {
				result[t] = true
			}
			return result
		}
	}
	return result
}

// hasColumn reports whether table has the column. MariaDB column names are
// case-insensitive.
func hasColumn(table *metadata.TableMetadata, name string) bool {
	name = strings.ToLower(name)
	for _, column := range table.GetColumns() {
		if strings.ToLower(column.GetName()) == name {
			return true
		}
	}
	return false
}

func firstKey(m map[string]*TableReference) string {
	for k := range m {
		return k
	}
	return ""
}

func generateSQL(targets []DMLTarget, opts BackupOptions) ([]BackupStatement, error) {
	groupByTable := make(map[string][]DMLTarget)
	for _, item := range targets {
		key := fmt.Sprintf("%s.%s", item.Table.Database, item.Table.Table)
		groupByTable[key] = append(groupByTable[key], item)
	}

	for key, list := range groupByTable {
		stmtType := DMLTypeUnknown
		for _, item := range list {
			if stmtType == DMLTypeUnknown {
				stmtType = item.Table.StatementType
			}
			if stmtType != item.Table.StatementType {
				return nil, fmt.Errorf("prior backup cannot handle mixed DML statements on the same table %s", key)
			}
		}
	}

	var result []BackupStatement
	for key, list := range groupByTable {
		backupStatement, err := generateSQLForTable(list, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL for table %s: %w", key, err)
		}
		result = append(result, *backupStatement)
	}

	slices.SortFunc(result, func(i, j BackupStatement) int {
		if i.ByteStart != j.ByteStart {
			if i.ByteStart < j.ByteStart {
				return -1
			}
			return 1
		}
		if i.SourceTable < j.SourceTable {
			return -1
		}
		if i.SourceTable > j.SourceTable {
			return 1
		}
		return 0
	})

	return result, nil
}

func generateSQLForTable(targets []DMLTarget, opts BackupOptions) (*BackupStatement, error) {
	table := targets[0].Table

	generatedColumns, normalColumns, err := opts.classifyColumns(table.Database, table.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to classify columns: %w", err)
	}

	targetTable := fmt.Sprintf("%s_%s_%s", opts.TablePrefix, table.Table, table.Database)
	targetTable, _ = truncateString(targetTable, maxTableNameLength)
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE TABLE `%s`.`%s` LIKE `%s`.`%s`;\n", opts.BackupDatabase, targetTable, table.Database, table.Table)

	fmt.Fprintf(&buf, "INSERT INTO `%s`.`%s`", opts.BackupDatabase, targetTable)
	if len(generatedColumns) > 0 {
		buf.WriteString(" (")
		for i, column := range normalColumns {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(&buf, "`%s`", column)
		}
		buf.WriteString(")")
	}
	for i, item := range targets {
		if i != 0 {
			buf.WriteString("\n  UNION DISTINCT\n")
		}
		tableNameOrAlias := item.Table.Table
		if len(item.Table.Alias) > 0 {
			tableNameOrAlias = item.Table.Alias
		}
		// Each arm keeps its own ORDER BY and LIMIT only inside parentheses.
		buf.WriteString("  ")
		if len(targets) > 1 {
			buf.WriteByte('(')
		}

		stmtLoc := nodeStmtLoc(item.Node)
		cteString := extractCTE(item.FullSQL, stmtLoc)
		if len(cteString) > 0 {
			fmt.Fprintf(&buf, "%s ", cteString)
		}

		if len(generatedColumns) == 0 {
			fmt.Fprintf(&buf, "SELECT `%s`.* FROM ", tableNameOrAlias)
		} else {
			buf.WriteString("SELECT ")
			for j, column := range normalColumns {
				if j > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(&buf, "`%s`.`%s`", tableNameOrAlias, column)
			}
			buf.WriteString(" FROM ")
		}
		if err := writeSuffixSelectClause(&buf, item.Node, item.FullSQL); err != nil {
			return nil, fmt.Errorf("failed to extract suffix select statement: %w", err)
		}
		if len(targets) > 1 {
			buf.WriteByte(')')
		}
	}

	buf.WriteByte(';')

	return &BackupStatement{
		Statement:   buf.String(),
		SourceTable: table.Table,
		BackupTable: targetTable,
		ByteStart:   targets[0].ByteStart,
		ByteEnd:     targets[len(targets)-1].ByteEnd,
	}, nil
}

func nodeStmtLoc(node ast.Node) ast.Loc {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return n.Loc
	case *ast.DeleteStmt:
		return n.Loc
	}
	return ast.Loc{}
}

func (opts BackupOptions) classifyColumns(database, table string) ([]string, []string, error) {
	if opts.Columns != nil {
		return opts.Columns(database, table)
	}
	return ClassifyColumns(opts.Catalog, opts.CaseSensitive, database, table)
}

// ClassifyColumns returns the generated columns of database.table and its
// other columns, in table order. The database is matched against
// ListDatabaseNames; a database the catalog does not have is a
// *NotFoundError.
func ClassifyColumns(catalog Catalog, caseSensitive bool, database, table string) (generated, normal []string, err error) {
	if catalog.GetDatabase == nil {
		return nil, nil, errors.New("the catalog has no GetDatabase")
	}

	var db Database
	allDatabaseNames, err := catalog.ListDatabaseNames()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list databases names: %w", err)
	}
	for _, name := range allDatabaseNames {
		if equalIdentifier(name, database, caseSensitive) {
			db, err = catalog.GetDatabase(name)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get database metadata for database %q: %w", name, err)
			}
			break
		}
	}
	if db == nil {
		return nil, nil, &NotFoundError{Database: &database}
	}

	if !db.HasSchema() {
		return nil, nil, errors.New("failed to get schema metadata")
	}

	var tableSchema *metadata.TableMetadata
	if !caseSensitive {
		for _, tableName := range db.ListTableNames() {
			if strings.EqualFold(tableName, table) {
				tableSchema = db.GetTable(tableName)
				break
			}
		}
	} else {
		tableSchema = db.GetTable(table)
	}
	if tableSchema == nil {
		return nil, nil, fmt.Errorf("table %s not found in schema", table)
	}

	for _, column := range tableSchema.GetColumns() {
		if column.GetGeneration() != nil {
			generated = append(generated, column.GetName())
		} else {
			normal = append(normal, column.GetName())
		}
	}

	return generated, normal, nil
}

func equalIdentifier(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// truncateString returns str cut to limit runes, and whether it was cut.
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
