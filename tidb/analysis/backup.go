package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

const (
	maxTableNameLength = 64
	maxMixedDMLCount   = 5
)

// BackupOptions configures TransformDMLToSelect.
type BackupOptions struct {
	// DefaultDatabase is the database of unqualified tables. Prior backup
	// refuses a statement that changes a table in another database.
	DefaultDatabase string
	// BackupDatabase is where the backup tables are created.
	BackupDatabase string
	// TablePrefix starts each backup table's name.
	TablePrefix string
	// CaseSensitive makes identifier comparisons case-sensitive.
	CaseSensitive bool
	// Catalog is the instance's metadata.
	Catalog Catalog
	// Columns returns a table's generated columns and its other columns. Nil
	// means ClassifyColumns over Catalog.
	Columns func(database, table string) (generated, normal []string, err error)
}

// BackupStatement backs up the rows one or more statements remove or change.
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

type tableReference struct {
	Database string
	Table    string
	Alias    string
	// ExplicitSchema is true when the reference carried an explicit schema
	// qualifier (e.g. db.test). A CTE reference is always unqualified, so a
	// schema-qualified table is never a CTE even if its name matches one.
	ExplicitSchema bool
}

type statementInfo struct {
	offset    int
	statement string
	table     *tableReference
	node      ast.Node
	fullSQL   string
	byteStart int
	byteEnd   int
}

// TransformDMLToSelect returns statements that copy the rows the UPDATE and
// DELETE statements of a script touch into backup tables. Up to
// maxMixedDMLCount statements get a backup each; more must all change one
// table and share a backup.
func TransformDMLToSelect(statement string, opts BackupOptions) ([]BackupStatement, error) {
	var database Database
	if opts.Catalog.GetDatabase != nil {
		if db, err := opts.Catalog.GetDatabase(opts.DefaultDatabase); err == nil {
			database = db
		}
	}

	statementInfoList, err := prepareTransformation(opts.DefaultDatabase, statement, database, opts.CaseSensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(statementInfoList, opts)
}

func prepareTransformation(databaseName, statement string, database Database, isCaseSensitive bool) ([]statementInfo, error) {
	var result []statementInfo
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
			tables, err := extractTables(databaseName, node, text, database, isCaseSensitive)
			if err != nil {
				return nil, fmt.Errorf("failed to extract tables: %w", err)
			}
			for _, table := range tables {
				// A BackupStatement carries no source database, and the executor
				// records the task database for the backed-up table. Reject
				// cross-database mutations (DELETE or UPDATE) so a rollback can't
				// be written to the wrong database.
				if !equalIdentifier(table.table.Database, databaseName, isCaseSensitive) {
					return nil, fmt.Errorf("prior backup does not support cross-database mutations: %s.%s (task database %q)", table.table.Database, table.table.Table, databaseName)
				}
				table.offset = i
				table.byteStart = segment.ByteStart
				table.byteEnd = segment.ByteEnd
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

// extractTables extracts the affected table references from a single DML node.
func extractTables(databaseName string, node ast.Node, fullSQL string, database Database, isCaseSensitive bool) ([]statementInfo, error) {
	switch n := node.(type) {
	case *ast.DeleteStmt:
		infos, err := extractTablesFromDelete(databaseName, n, fullSQL, isCaseSensitive)
		return requireTargets(infos, err, "DELETE")
	case *ast.UpdateStmt:
		infos, err := extractTablesFromUpdate(databaseName, n, fullSQL, database, isCaseSensitive)
		return requireTargets(infos, err, "UPDATE")
	case *ast.BatchStmt:
		// TiDB BATCH (non-transactional DML) is not supported by prior backup.
		// Reject it explicitly rather than returning an empty list, which the
		// task executor would treat as a successful no-op and then run the
		// mutation with no backup.
		return nil, errors.New("prior backup does not support TiDB BATCH (non-transactional DML) statements")
	case *ast.ExplainStmt:
		// EXPLAIN ANALYZE executes the wrapped statement (TiDB modifies data for
		// a DML), but prior backup can't unwrap/restore it — reject the ANALYZE
		// form of an UPDATE/DELETE/BATCH. Plain EXPLAIN does not execute, so it
		// needs no backup and is skipped.
		if n.Analyze {
			switch s := n.Stmt.(type) {
			case *ast.UpdateStmt, *ast.DeleteStmt, *ast.BatchStmt:
				return nil, errors.New("prior backup does not support EXPLAIN ANALYZE of an UPDATE/DELETE/BATCH statement")
			case *ast.InsertStmt:
				if isOverwritingInsert(s) {
					return nil, errors.New("prior backup does not support EXPLAIN ANALYZE of a REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statement")
				}
			default:
			}
		}
		return nil, nil
	case *ast.InsertStmt:
		// REPLACE and INSERT ... ON DUPLICATE KEY UPDATE can overwrite or delete
		// existing rows, but prior backup cannot back them up. Reject them rather
		// than returning an empty list, which the task executor would treat as a
		// successful no-op and then run the mutation unprotected. A plain INSERT
		// only adds rows (nothing to restore), so it needs no backup.
		if isOverwritingInsert(n) {
			return nil, errors.New("prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements")
		}
		return nil, nil
	default:
		return nil, nil
	}
}

// isOverwritingInsert reports whether an INSERT can overwrite or delete
// existing rows (REPLACE, or INSERT ... ON DUPLICATE KEY UPDATE), which prior
// backup cannot back up. A plain INSERT only adds rows.
func isOverwritingInsert(n *ast.InsertStmt) bool {
	return n.IsReplace || len(n.OnDuplicateKey) > 0
}

// identifierKey normalizes an identifier for use as a map key/lookup, honoring
// the instance's case sensitivity. TiDB resolves table-name and alias
// qualifiers case-insensitively by default, so a case-mismatched reference
// (UPDATE test AS T SET t.c = ...) must still resolve to the same table.
func identifierKey(name string, isCaseSensitive bool) string {
	if isCaseSensitive {
		return name
	}
	return strings.ToLower(name)
}

// requireTargets fails when a recognized DML statement yields no backup target.
// Returning an empty list with no error would let the task executor treat the
// backup as a successful no-op and run the mutation unprotected (e.g. when a
// table alias collides with a CTE name, so the only target gets filtered out).
func requireTargets(infos []statementInfo, err error, kind string) ([]statementInfo, error) {
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("prior backup could not determine a backup target for the %s statement", kind)
	}
	return infos, nil
}

func extractTablesFromDelete(databaseName string, n *ast.DeleteStmt, fullSQL string, isCaseSensitive bool) ([]statementInfo, error) {
	cteNames := collectCTENames(fullSQL, n.Loc)
	stmtText := extractStatementText(fullSQL, n.Loc)

	singleTables := collectSingleTables(databaseName, cteNames, n.Tables, isCaseSensitive)

	if len(n.Using) == 0 {
		// Single-table DELETE: DELETE FROM t WHERE ...
		if len(singleTables) == 0 {
			return nil, nil
		}
		first := singleTables[firstKey(singleTables)]
		return []statementInfo{{
			statement: stmtText,
			table:     first,
			node:      n,
			fullSQL:   fullSQL,
		}}, nil
	}

	// Multi-table DELETE: DELETE t1, t2 FROM t1 JOIN t2 ... WHERE ...
	// Tables = targets to delete from; Using = the referenced table list.
	refTables := collectSingleTables(databaseName, cteNames, n.Using, isCaseSensitive)

	var result []statementInfo
	for _, target := range singleTables {
		name := target.Table
		if target.Alias != "" {
			name = target.Alias
		}
		ref, ok := refTables[identifierKey(name, isCaseSensitive)]
		if !ok {
			return nil, fmt.Errorf("cannot extract reference table: no matched table %q in referenced table list", name)
		}
		// Skip only if the target resolves to a CTE reference (you can't delete a
		// CTE). A real table aliased with a CTE's name is still a delete target.
		if isCTERef(cteNames, ref) {
			continue
		}
		result = append(result, statementInfo{
			statement: stmtText,
			table:     ref,
			node:      n,
			fullSQL:   fullSQL,
		})
	}
	return result, nil
}

func extractTablesFromUpdate(databaseName string, n *ast.UpdateStmt, fullSQL string, database Database, isCaseSensitive bool) ([]statementInfo, error) {
	cteNames := collectCTENames(fullSQL, n.Loc)
	stmtText := extractStatementText(fullSQL, n.Loc)

	singleTables := collectSingleTables(databaseName, cteNames, n.Tables, isCaseSensitive)

	// Determine which tables the SET clause writes, via column table prefixes.
	updatedTables := make(map[string]bool)
	var unqualifiedColumns []string
	for _, assign := range n.SetList {
		if assign == nil || assign.Column == nil {
			continue
		}
		table := assign.Column.Table
		if table == "" {
			updatedTables[""] = true
			unqualifiedColumns = append(unqualifiedColumns, assign.Column.Column)
			continue
		}
		// Look up the qualifier case-insensitively (TiDB default): the alias/name
		// resolves regardless of the case used in the SET clause.
		key := identifierKey(table, isCaseSensitive)
		// Skip only if the qualifier resolves to a CTE reference in the FROM — a
		// CTE can't be a mutation target. A real table aliased with a CTE's name
		// (its resolved table is real) is still an update target.
		if ref, ok := singleTables[key]; ok && isCTERef(cteNames, ref) {
			continue
		}
		updatedTables[key] = true
	}

	// Resolve unqualified SET columns to their owning table(s) via metadata.
	// Exclude CTE refs from the candidates: a CTE can't be a mutation target,
	// and resolveUnqualifiedColumns' metadata-unavailable/unresolved fallback
	// would otherwise mark every candidate (CTEs included).
	if updatedTables[""] {
		delete(updatedTables, "")
		candidates := make(map[string]*tableReference, len(singleTables))
		for key, ref := range singleTables {
			if isCTERef(cteNames, ref) {
				continue
			}
			candidates[key] = ref
		}
		for t := range resolveUnqualifiedColumns(unqualifiedColumns, candidates, database) {
			updatedTables[t] = true
		}
	}

	// Single-table UPDATE: no explicit qualification needed.
	if len(updatedTables) == 0 && len(singleTables) == 1 {
		for _, ref := range singleTables {
			return []statementInfo{{
				statement: stmtText,
				table:     ref,
				node:      n,
				fullSQL:   fullSQL,
			}}, nil
		}
	}

	var result []statementInfo
	for table := range updatedTables {
		ref, ok := singleTables[table]
		if !ok {
			return nil, fmt.Errorf("cannot extract reference table: no matched updated table %q in referenced table list", table)
		}
		result = append(result, statementInfo{
			statement: stmtText,
			table:     ref,
			node:      n,
			fullSQL:   fullSQL,
		})
	}
	return result, nil
}

// collectSingleTables walks TableExpr slices and collects TableReference entries
// keyed by alias or table name. A table's database is its explicit schema, or
// the statement's database when unqualified.
func collectSingleTables(databaseName string, cteNames map[string]bool, exprs []ast.TableExpr, isCaseSensitive bool) map[string]*tableReference {
	result := make(map[string]*tableReference)
	for _, expr := range exprs {
		collectSingleTablesFromExpr(databaseName, cteNames, expr, isCaseSensitive, result)
	}
	return result
}

func collectSingleTablesFromExpr(databaseName string, cteNames map[string]bool, expr ast.TableExpr, isCaseSensitive bool, out map[string]*tableReference) {
	switch e := expr.(type) {
	case *ast.TableRef:
		db := e.Schema
		if db == "" {
			db = databaseName
		}
		ref := &tableReference{Database: db, Table: e.Name, Alias: e.Alias, ExplicitSchema: e.Schema != ""}
		// Key by alias-or-name, normalized for case-insensitive instances (TiDB
		// default) so a case-mismatched qualifier still resolves. The original
		// Table/Alias on ref are preserved for the emitted SQL.
		key := e.Name
		if e.Alias != "" {
			key = e.Alias
		}
		key = identifierKey(key, isCaseSensitive)
		// A schema-qualified physical table and an unqualified same-named CTE can
		// coexist in one FROM (TiDB keeps them distinct) but collapse to the same
		// map key. Keep the real table: a CTE is never a backup target, and
		// letting it shadow a physical target would drop that target's backup.
		if existing, ok := out[key]; ok && !isCTERef(cteNames, existing) && isCTERef(cteNames, ref) {
			return
		}
		out[key] = ref
	case *ast.JoinClause:
		if e.Left != nil {
			collectSingleTablesFromExpr(databaseName, cteNames, e.Left, isCaseSensitive, out)
		}
		if e.Right != nil {
			collectSingleTablesFromExpr(databaseName, cteNames, e.Right, isCaseSensitive, out)
		}
	default:
	}
}

// resolveUnqualifiedColumns resolves unqualified SET column names to their
// owning table(s) using metadata. Falls back to all referenced tables when
// metadata is unavailable or a column cannot be located.
func resolveUnqualifiedColumns(columns []string, singleTables map[string]*tableReference, database Database) map[string]bool {
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

// hasColumn reports whether table has the column. TiDB column names are
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

func firstKey(m map[string]*tableReference) string {
	for k := range m {
		return k
	}
	return ""
}

// collectCTENames returns the names of the CTEs the statement opens with,
// lower-cased for case-insensitive matching.
func collectCTENames(fullSQL string, stmtLoc ast.Loc) map[string]bool {
	result := make(map[string]bool)
	for _, cte := range parseCTEs(withClause(fullSQL, stmtLoc)) {
		result[strings.ToLower(cte.Name)] = true
	}
	return result
}

// isCTERef reports whether ref resolves to a CTE. A CTE reference is always
// unqualified, so a schema-qualified physical table (db.test) is never a CTE
// even if its name matches one. CTE names match case-insensitively (TiDB
// resolves a case-mismatched reference to the CTE, which shadows a same-named
// real table), consistent with the file's other identifier comparisons.
func isCTERef(cteNames map[string]bool, ref *tableReference) bool {
	return !ref.ExplicitSchema && cteNames[strings.ToLower(ref.Table)]
}

// writeCTEPrefix writes the statement's WITH (CTE) clause before the backup
// SELECT, so a SELECT whose FROM references a CTE stays valid.
func writeCTEPrefix(buf *strings.Builder, node ast.Node, fullSQL string) {
	cte := extractCTE(fullSQL, nodeStmtLoc(node))
	if cte == "" {
		return
	}
	fmt.Fprintf(buf, "%s ", cte)
}

// extractCTE returns the WITH clause the statement opens with, or "" when it
// has none or the clause does not parse.
func extractCTE(fullSQL string, stmtLoc ast.Loc) string {
	clause := withClause(fullSQL, stmtLoc)
	if len(parseCTEs(clause)) == 0 {
		return ""
	}
	return clause
}

// withClause returns the WITH clause that opens the UPDATE or DELETE at loc, or
// "". omni starts the node at its WITH clause and keeps no CTE on it, so the
// clause is the text before the UPDATE or DELETE keyword outside parentheses.
func withClause(sql string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End {
		return ""
	}
	text := sql[loc.Start:loc.End]
	tokens := parser.Tokenize(text)
	if len(tokens) == 0 || parser.TokenName(tokens[0].Type) != "WITH" {
		return ""
	}
	depth := 0
	for _, token := range tokens[1:] {
		switch token.Type {
		case '(':
			depth++
		case ')':
			depth--
		default:
			if depth != 0 {
				continue
			}
			if name := parser.TokenName(token.Type); name == "UPDATE" || name == "DELETE" {
				return strings.TrimSpace(text[:token.Loc])
			}
		}
	}
	return ""
}

// nodeStmtLoc returns the source Loc of a DML statement node.
func nodeStmtLoc(node ast.Node) ast.Loc {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return n.Loc
	case *ast.DeleteStmt:
		return n.Loc
	}
	return ast.Loc{}
}

// parseCTEs returns the CTEs of a WITH clause, or nil when it does not parse.
func parseCTEs(clause string) []*ast.CommonTableExpr {
	if clause == "" {
		return nil
	}
	// Wrap the clause with a dummy SELECT so it parses as a complete statement.
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
	start := max(0, loc.Start)
	end := loc.End
	if end <= start || end > len(fullSQL) {
		end = len(fullSQL)
	}
	return strings.TrimSpace(fullSQL[start:end])
}

func generateSQL(statementInfoList []statementInfo, opts BackupOptions) ([]BackupStatement, error) {
	if len(statementInfoList) <= maxMixedDMLCount {
		return generateSQLForMixedDML(statementInfoList, opts)
	}
	return generateSQLForSingleTable(statementInfoList, opts)
}

func generateSQLForSingleTable(statementInfoList []statementInfo, opts BackupOptions) ([]BackupStatement, error) {
	table := statementInfoList[0].table

	for _, item := range statementInfoList {
		if !equalTable(table, item.table, opts.CaseSensitive) {
			return nil, fmt.Errorf("prior backup cannot handle statements on different tables more than %d", maxMixedDMLCount)
		}
		// This path UNION ALLs the statements into one backup SELECT. A WITH
		// (CTE) prefix can't be emitted per union arm (TiDB rejects WITH after
		// UNION ALL), so reject CTE statements here. CTEs are still supported on
		// the per-statement mixed-DML path (<= maxMixedDMLCount).
		if extractCTE(item.fullSQL, nodeStmtLoc(item.node)) != "" {
			return nil, fmt.Errorf("prior backup does not support WITH (CTE) statements when more than %d DML statements target the same table", maxMixedDMLCount)
		}
	}
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
	for i, item := range statementInfoList {
		if i != 0 {
			buf.WriteString("\n  UNION ALL\n")
		}
		tableNameOrAlias := item.table.Table
		if len(item.table.Alias) > 0 {
			tableNameOrAlias = item.table.Alias
		}
		// Each arm keeps its own ORDER BY and LIMIT only inside parentheses.
		if len(generatedColumns) == 0 {
			fmt.Fprintf(&buf, "  (SELECT `%s`.* FROM ", tableNameOrAlias)
		} else {
			buf.WriteString("  (SELECT ")
			for j, column := range normalColumns {
				if j > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(&buf, "`%s`.`%s`", tableNameOrAlias, column)
			}
			buf.WriteString(" FROM ")
		}
		writeSuffixSelectClause(&buf, item.node, item.fullSQL)
		buf.WriteByte(')')
	}

	buf.WriteByte(';')

	return []BackupStatement{
		{
			Statement:   buf.String(),
			SourceTable: table.Table,
			BackupTable: targetTable,
			ByteStart:   statementInfoList[0].byteStart,
			ByteEnd:     statementInfoList[len(statementInfoList)-1].byteEnd,
		},
	}, nil
}

// equalIdentifier compares two SQL identifiers honoring the instance's case
// sensitivity. Bytebase treats TiDB instances as case-insensitive today, but
// honoring the flag keeps this consistent with ClassifyColumns.
func equalIdentifier(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

func equalTable(a, b *tableReference, caseSensitive bool) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Database != "" && b.Database != "" && !equalIdentifier(a.Database, b.Database, caseSensitive) {
		return false
	}
	return equalIdentifier(a.Table, b.Table, caseSensitive)
}

func generateSQLForMixedDML(statementInfoList []statementInfo, opts BackupOptions) ([]BackupStatement, error) {
	var result []BackupStatement
	offsetLength := 1
	if len(statementInfoList) > 1 {
		offsetLength = digitCount(statementInfoList[len(statementInfoList)-1].offset)
	}

	for _, statementInfo := range statementInfoList {
		table := statementInfo.table
		targetTable := fmt.Sprintf("%s_%0*d_%s", opts.TablePrefix, offsetLength, statementInfo.offset, table.Table)
		targetTable, _ = truncateString(targetTable, maxTableNameLength)
		// If enforce_gtid_consistency = true, we cannot run CREATE TABLE .. AS SELECT.
		// So we create the table first and then run INSERT INTO .. SELECT.
		var buf strings.Builder
		fmt.Fprintf(&buf, "CREATE TABLE `%s`.`%s` LIKE `%s`.`%s`;\n", opts.BackupDatabase, targetTable, table.Database, table.Table)
		generatedColumns, normalColumns, err := opts.classifyColumns(table.Database, table.Table)
		if err != nil {
			return nil, fmt.Errorf("failed to classify columns: %w", err)
		}
		tableNameOrAlias := table.Table
		if len(table.Alias) > 0 {
			tableNameOrAlias = table.Alias
		}
		if len(generatedColumns) == 0 {
			fmt.Fprintf(&buf, "INSERT INTO `%s`.`%s` ", opts.BackupDatabase, targetTable)
			writeCTEPrefix(&buf, statementInfo.node, statementInfo.fullSQL)
			fmt.Fprintf(&buf, "SELECT `%s`.* FROM ", tableNameOrAlias)
		} else {
			fmt.Fprintf(&buf, "INSERT INTO `%s`.`%s` (", opts.BackupDatabase, targetTable)
			for i, column := range normalColumns {
				if i > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(&buf, "`%s`", column)
			}
			buf.WriteString(") ")
			writeCTEPrefix(&buf, statementInfo.node, statementInfo.fullSQL)
			buf.WriteString("SELECT ")
			for i, column := range normalColumns {
				if i > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(&buf, "`%s`.`%s`", tableNameOrAlias, column)
			}
			buf.WriteString(" FROM ")
		}
		writeSuffixSelectClause(&buf, statementInfo.node, statementInfo.fullSQL)
		buf.WriteByte(';')
		result = append(result, BackupStatement{
			Statement:   buf.String(),
			SourceTable: table.Table,
			BackupTable: targetTable,
			ByteStart:   statementInfo.byteStart,
			ByteEnd:     statementInfo.byteEnd,
		})
	}
	return result, nil
}

// digitCount returns the number of decimal digits of n.
func digitCount(n int) int {
	length := 1
	for {
		if n < 10 {
			return length
		}
		n /= 10
		length++
	}
}

// writeSuffixSelectClause writes the "<table refs> WHERE ... ORDER BY ... LIMIT ..."
// suffix of a DML statement, sliced from the original SQL via the node's Loc.
func writeSuffixSelectClause(buf *strings.Builder, node ast.Node, fullSQL string) {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		writeUpdateSuffix(buf, n, fullSQL)
	case *ast.DeleteStmt:
		writeDeleteSuffix(buf, n, fullSQL)
	}
}

func writeUpdateSuffix(buf *strings.Builder, n *ast.UpdateStmt, sql string) {
	// "table_refs": from the first table ref to the last table ref's end.
	if len(n.Tables) > 0 {
		start := nodeLocStart(n.Tables[0])
		end := nodeLocEnd(n.Tables[len(n.Tables)-1])
		if start >= 0 && end > start && end <= len(sql) {
			// omni excludes parentheses wrapping the table-ref list, so a
			// parenthesized join operand can leave the slice unbalanced.
			buf.WriteString(balancedTableRefs(sql[start:end]))
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
		if afterSet >= 0 && afterSet < stmtEnd {
			suffix := strings.TrimSpace(sql[afterSet:stmtEnd])
			if suffix != "" {
				fmt.Fprintf(buf, " %s", suffix)
			}
		}
	}
}

func writeDeleteSuffix(buf *strings.Builder, n *ast.DeleteStmt, sql string) {
	// Single-table DELETE: from the table ref to statement end.
	// Multi-table DELETE: from the USING/reference table list to statement end.
	var start int
	switch {
	case len(n.Using) > 0:
		start = nodeLocStart(n.Using[0])
	case len(n.Tables) > 0:
		start = nodeLocStart(n.Tables[0])
	default:
		return
	}
	end := n.Loc.End
	if end <= 0 || end > len(sql) {
		end = len(sql)
	}
	if start >= 0 && end > start {
		// omni excludes parentheses wrapping the table-ref list, so the slice can
		// carry an unmatched ")"; rebalance it before emitting.
		buf.WriteString(balancedTableRefs(sql[start:end]))
	}
}

// balancedTableRefs returns s (trimmed) with parentheses balanced. omni's
// JoinClause/TableRef Loc excludes a wrapping "(" or ")", so a sliced table-ref
// list can carry an unmatched ")" (a parenthesized left join operand) and/or an
// unmatched "(" (a right operand). The missing parentheses only group, so their
// source position is irrelevant — prepend "(" and append ")" to balance.
func balancedTableRefs(s string) string {
	leading, trailing := parenDeficit(s)
	s = strings.TrimSpace(s)
	if leading > 0 {
		s = strings.Repeat("(", leading) + s
	}
	if trailing > 0 {
		s += strings.Repeat(")", trailing)
	}
	return s
}

// parenDeficit returns how many "(" must be prepended and ")" appended to
// balance s, ignoring parentheses inside string/identifier literals ('...',
// "...", `...`) and comments (/* */, -- , #).
func parenDeficit(s string) (leading int, trailing int) {
	open := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '\'' || c == '"' || c == '`':
			i = skipQuoted(s, i)
		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			i = skipBlockComment(s, i)
		case c == '-' && i+1 < len(s) && s[i+1] == '-' && (i+2 >= len(s) || isASCIISpace(s[i+2])):
			i = skipLineComment(s, i)
		case c == '#':
			i = skipLineComment(s, i)
		case c == '(':
			open++
		case c == ')':
			if open > 0 {
				open--
			} else {
				leading++
			}
		default:
		}
	}
	return leading, open
}

// skipQuoted returns the index of the closing quote of the string/identifier
// literal opened at i (s[i] is the opening quote), or len(s)-1 if unterminated.
func skipQuoted(s string, i int) int {
	quote := s[i]
	for j := i + 1; j < len(s); {
		switch {
		case s[j] == '\\' && quote != '`':
			j += 2 // backslash escape ('...'/"..." only)
		case s[j] == quote && j+1 < len(s) && s[j+1] == quote:
			j += 2 // doubled-quote escape ('' "" ``)
		case s[j] == quote:
			return j
		default:
			j++
		}
	}
	return len(s) - 1
}

// skipBlockComment returns the index of the closing "/" of the /* */ comment
// opened at i, or len(s)-1 if unterminated.
func skipBlockComment(s string, i int) int {
	for j := i + 2; j+1 < len(s); j++ {
		if s[j] == '*' && s[j+1] == '/' {
			return j + 1
		}
	}
	return len(s) - 1
}

// skipLineComment returns the index of the newline ending the -- or # comment
// opened at i, or len(s)-1 if it runs to the end.
func skipLineComment(s string, i int) int {
	for j := i; j < len(s); j++ {
		if s[j] == '\n' {
			return j
		}
	}
	return len(s) - 1
}

func isASCIISpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func nodeLocStart(expr ast.TableExpr) int {
	switch e := expr.(type) {
	case *ast.TableRef:
		return e.Loc.Start
	case *ast.JoinClause:
		return e.Loc.Start
	case *ast.SubqueryExpr:
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
	case *ast.SubqueryExpr:
		return e.Loc.End
	}
	return -1
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

	tableSchema := db.GetTable(table)
	if tableSchema == nil {
		return nil, nil, fmt.Errorf("failed to get table metadata for table %q", table)
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
