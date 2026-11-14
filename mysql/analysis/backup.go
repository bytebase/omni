package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
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
	// Columns returns what the backup needs of a table's columns and keys.
	// Nil means ClassifyColumns over Catalog.
	Columns func(database, table string) (TableColumns, error)
}

// TableColumns is what a backup and its restore need of a table.
type TableColumns struct {
	// Generated are the generated columns, which a backup leaves out.
	Generated []string
	// Normal are the other columns, invisible ones included, in table order.
	Normal []string
	// RowKey reports whether a primary key, or a unique key whose parts are
	// all non-nullable columns, tells the table's rows apart.
	RowKey bool
	// RowKeyColumns are the parts of that key.
	RowKeyColumns []string
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
	targets, err := prepareTransformation(opts.DefaultDatabase, statement, opts.Catalog, opts.CaseSensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(targets, opts)
}

func prepareTransformation(databaseName, statement string, catalog Catalog, caseSensitive bool) ([]DMLTarget, error) {
	var result []DMLTarget
	// written holds the tables the statements before the current one change.
	var written []ColumnResource
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
			tables, err := ExtractTables(databaseName, node, text, i, catalog)
			if err != nil {
				return nil, fmt.Errorf("failed to extract tables: %w", err)
			}
			// The backups read their rows before any statement runs, so a
			// statement must not read a table an earlier one changes.
			if len(tables) > 0 {
				for read := range statementReads(databaseName, node, text, tables) {
					for _, changed := range written {
						if equalIdentifier(read.Database, changed.Database, caseSensitive) && equalIdentifier(read.Table, changed.Table, caseSensitive) {
							return nil, fmt.Errorf("prior backup cannot handle a statement that reads %s.%s after an earlier statement changes it", read.Database, read.Table)
						}
					}
				}
			}
			written = append(written, statementWrites(databaseName, node, tables)...)
			for _, table := range tables {
				table.ByteStart = segment.ByteStart
				table.ByteEnd = segment.ByteEnd
				result = append(result, table)
			}
		}
	}

	return result, nil
}

// statementReads returns the tables an UPDATE or DELETE reads other than
// through its targets: subqueries, CTEs, and the other tables it joins.
func statementReads(databaseName string, node ast.Node, fullSQL string, targets []DMLTarget) ColumnSet {
	result := make(ColumnSet)
	target := make(map[string]bool)
	for _, t := range targets {
		key := t.Table.Table
		if t.Table.Alias != "" {
			key = t.Table.Alias
		}
		target[key] = true
	}
	var loc ast.Loc
	var tables []ast.TableExpr
	switch n := node.(type) {
	case *ast.UpdateStmt:
		loc, tables = n.Loc, n.Tables
		for _, assign := range n.SetList {
			collectAccessTablesFromExpr(result, assign.Value, databaseName)
		}
		collectAccessTablesFromExpr(result, n.Where, databaseName)
		collectAccessTablesFromOrderBy(result, n.OrderBy, databaseName)
	case *ast.DeleteStmt:
		loc, tables = n.Loc, n.Tables
		if len(n.Using) > 0 {
			tables = n.Using
		}
		collectAccessTablesFromExpr(result, n.Where, databaseName)
		collectAccessTablesFromOrderBy(result, n.OrderBy, databaseName)
	default:
		return result
	}
	cteNames := make(map[string]bool)
	for _, cte := range parseCTEs(withClause(fullSQL, loc)) {
		cteNames[cte.Name] = true
		collectAccessTablesFromNode(result, cte.Select, databaseName)
	}
	var walk func(expr ast.TableExpr)
	walk = func(expr ast.TableExpr) {
		switch e := expr.(type) {
		case *ast.TableRef:
			key := e.Name
			if e.Alias != "" {
				key = e.Alias
			}
			if target[key] || (e.Schema == "" && cteNames[e.Name]) {
				return
			}
			collectAccessTablesFromTableRef(result, e, databaseName)
		case *ast.JoinClause:
			walk(e.Left)
			walk(e.Right)
			if on, ok := e.Condition.(*ast.OnCondition); ok {
				collectAccessTablesFromExpr(result, on.Expr, databaseName)
			}
		default:
			collectAccessTablesFromTableExpr(result, expr, databaseName)
		}
	}
	for _, table := range tables {
		walk(table)
	}
	return result
}

// statementWrites returns the tables a statement changes: an UPDATE's or
// DELETE's targets, or an INSERT's or REPLACE's table.
func statementWrites(databaseName string, node ast.Node, targets []DMLTarget) []ColumnResource {
	var result []ColumnResource
	for _, t := range targets {
		result = append(result, ColumnResource{Database: t.Table.Database, Table: t.Table.Table})
	}
	if insert, ok := node.(*ast.InsertStmt); ok && insert.Table != nil {
		database := insert.Table.Schema
		if database == "" {
			database = databaseName
		}
		result = append(result, ColumnResource{Database: database, Table: insert.Table.Name})
	}
	return result
}

// segmentText is the segment's text with its ';', if any.
func segmentText(sql string, segment parser.Segment) string {
	end := segment.ByteEnd
	if end < len(sql) && sql[end] == ';' {
		end++
	}
	return sql[segment.ByteStart:end]
}

// ExtractTables returns the tables a DELETE or UPDATE changes. catalog's
// metadata, where it has the database of each table the statement names,
// finds the table an unqualified SET column belongs to.
func ExtractTables(databaseName string, node ast.Node, fullSQL string, offset int, catalog Catalog) ([]DMLTarget, error) {
	switch n := node.(type) {
	case *ast.DeleteStmt:
		return extractTablesFromDelete(databaseName, n, fullSQL, offset)
	case *ast.UpdateStmt:
		return extractTablesFromUpdate(databaseName, n, fullSQL, offset, catalog)
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

func extractTablesFromUpdate(databaseName string, n *ast.UpdateStmt, fullSQL string, offset int, catalog Catalog) ([]DMLTarget, error) {
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
		resolved := resolveUnqualifiedColumns(unqualifiedColumns, singleTables, catalog)
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

// writeSuffixSelectClause writes the FROM ... WHERE ... suffix of a DML
// statement, with filter, if any, added to its WHERE.
func writeSuffixSelectClause(buf *strings.Builder, node ast.Node, fullSQL string, filter string) {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		writeUpdateSuffix(buf, n, fullSQL, filter)
	case *ast.DeleteStmt:
		writeDeleteSuffix(buf, n, fullSQL, filter)
	}
}

// withFilter returns sql[start:end] with filter ANDed to where, which lies in
// that range, or appended as a WHERE when where is nil. A statement whose
// target an outer join can leave NULL joins tables, so it has no ORDER BY or
// LIMIT after its WHERE.
func withFilter(sql string, start, end int, where ast.ExprNode, filter string) string {
	if filter == "" {
		return sql[start:end]
	}
	if where == nil {
		return strings.TrimSpace(sql[start:end]) + " WHERE " + filter
	}
	loc, _ := nodeLoc(where)
	return sql[start:loc.Start] + "(" + sql[loc.Start:loc.End] + ") AND " + filter + sql[loc.End:end]
}

func writeUpdateSuffix(buf *strings.Builder, n *ast.UpdateStmt, sql string, filter string) {
	// For UPDATE: write "table_refs WHERE ... ORDER BY ... LIMIT ..."
	// Use the last table's Loc.End as the boundary — everything from first table
	// to last table's end is the table reference list.
	if len(n.Tables) > 0 {
		start := nodeLocStart(n.Tables[0])
		end := nodeLocEnd(n.Tables[len(n.Tables)-1])
		if start >= 0 && end > start && end <= len(sql) {
			// A table location leaves out the parentheses that group it, so
			// the slice can hold only one of a pair.
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
		suffix := strings.TrimSpace(withFilter(sql, afterSet, stmtEnd, n.Where, filter))
		if suffix != "" {
			fmt.Fprintf(buf, " %s", suffix)
		}
	}
}

func writeDeleteSuffix(buf *strings.Builder, n *ast.DeleteStmt, sql string, filter string) {
	// For single-table DELETE: everything from the table ref to statement end.
	// For multi-table DELETE: everything from the USING/reference table list to end.
	if len(n.Using) > 0 {
		start := nodeLocStart(n.Using[0])
		end := n.Loc.End
		if end <= 0 || end > len(sql) {
			end = len(sql)
		}
		if start >= 0 && end > start {
			buf.WriteString(balancedTableRefs(withFilter(sql, start, end, n.Where, filter)))
		}
	} else if len(n.Tables) > 0 {
		start := nodeLocStart(n.Tables[0])
		end := n.Loc.End
		if end <= 0 || end > len(sql) {
			end = len(sql)
		}
		if start >= 0 && end > start {
			buf.WriteString(balancedTableRefs(withFilter(sql, start, end, n.Where, filter)))
		}
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

// resolveUnqualifiedColumns resolves unqualified column names to owning tables
// using each table's own database's metadata. Without the metadata of every
// table's database, every table may own the columns.
func resolveUnqualifiedColumns(columns []string, singleTables map[string]*TableReference, catalog Catalog) map[string]bool {
	result := make(map[string]bool)

	databases := make(map[string]Database)
	for _, ref := range singleTables {
		if _, ok := databases[ref.Database]; ok {
			continue
		}
		var database Database
		if catalog.GetDatabase != nil {
			if db, err := catalog.GetDatabase(ref.Database); err == nil {
				database = db
			}
		}
		if database == nil || !database.HasSchema() {
			for t := range singleTables {
				result[t] = true
			}
			return result
		}
		databases[ref.Database] = database
	}

	for _, col := range columns {
		resolved := false
		for aliasOrName, ref := range singleTables {
			tableMetadata := databases[ref.Database].GetTable(ref.Table)
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

// hasColumn reports whether table has the column. MySQL column names are
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
		// The backup reads every statement's rows before any runs, so a later
		// statement's LIMIT may pick rows the backup never read.
		for i, item := range list {
			if i > 0 && hasLimit(item.Node) {
				return nil, fmt.Errorf("prior backup cannot handle a statement with LIMIT after another statement on the same table %s", key)
			}
		}
	}

	type keyed struct {
		key       string
		statement BackupStatement
	}
	var backups []keyed
	for key, list := range groupByTable {
		backupStatement, err := generateSQLForTable(list, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL for table %s: %w", key, err)
		}
		backups = append(backups, keyed{key: key, statement: *backupStatement})
	}

	slices.SortFunc(backups, func(a, b keyed) int {
		i, j := a.statement, b.statement
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

	// Truncation can give two tables' backups one name.
	var result []BackupStatement
	for i, backup := range backups {
		for _, earlier := range backups[:i] {
			if equalIdentifier(backup.statement.BackupTable, earlier.statement.BackupTable, opts.CaseSensitive) {
				first, second := min(earlier.key, backup.key), max(earlier.key, backup.key)
				return nil, fmt.Errorf("prior backup cannot name the backups of %s and %s apart: both are %s", first, second, backup.statement.BackupTable)
			}
		}
		result = append(result, backup.statement)
	}

	return result, nil
}

func generateSQLForTable(targets []DMLTarget, opts BackupOptions) (*BackupStatement, error) {
	table := targets[0].Table

	columns, err := opts.classifyColumns(table.Database, table.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to classify columns: %w", err)
	}
	// Without a row key, a row a join repeats or two statements both read
	// cannot be told from rows that are equal.
	if !columns.RowKey {
		if len(targets) > 1 {
			return nil, errors.New("prior backup cannot combine several statements on a table without a primary key or a unique key of non-null columns")
		}
		if joins(targets[0].Node) {
			return nil, errors.New("prior backup cannot back up a table without a primary key or a unique key of non-null columns from a statement that joins other tables")
		}
	}
	if err := checkForeignKeys(opts, table, targets); err != nil {
		return nil, err
	}

	targetTable := fmt.Sprintf("%s_%s_%s", opts.TablePrefix, table.Table, table.Database)
	targetTable, _ = truncateString(targetTable, maxTableNameLength)
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE TABLE %s.%s LIKE %s.%s;\n", quoteIdentifier(opts.BackupDatabase), quoteIdentifier(targetTable), quoteIdentifier(table.Database), quoteIdentifier(table.Table))

	// The columns are listed, as a wildcard leaves invisible columns out.
	fmt.Fprintf(&buf, "INSERT INTO %s.%s", quoteIdentifier(opts.BackupDatabase), quoteIdentifier(targetTable))
	if len(columns.Normal) > 0 {
		fmt.Fprintf(&buf, " (%s)", quoteColumnList(columns.Normal, "", ","))
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

		buf.WriteString("SELECT ")
		// A join repeats a row once per match; the row key makes the copies
		// equal, so DISTINCT keeps one.
		if joins(item.Node) {
			buf.WriteString("DISTINCT ")
		}
		if len(columns.Normal) == 0 {
			fmt.Fprintf(&buf, "%s.*", quoteIdentifier(tableNameOrAlias))
		} else {
			buf.WriteString(quoteColumnList(columns.Normal, quoteIdentifier(tableNameOrAlias)+".", ","))
		}
		buf.WriteString(" FROM ")
		filter, err := outerJoinFilter(item, tableNameOrAlias, columns)
		if err != nil {
			return nil, err
		}
		writeSuffixSelectClause(&buf, item.Node, item.FullSQL, filter)
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

// hasLimit reports whether an UPDATE or DELETE has a LIMIT.
func hasLimit(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return n.Limit != nil
	case *ast.DeleteStmt:
		return n.Limit != nil
	}
	return false
}

// outerJoinFilter returns the condition that drops the rows an outer join
// makes up for item's target, which hold NULL in every column, or "" when the
// target is on no outer join's nullable side.
func outerJoinFilter(item DMLTarget, tableNameOrAlias string, columns TableColumns) (string, error) {
	var tables []ast.TableExpr
	var where ast.ExprNode
	switch n := item.Node.(type) {
	case *ast.UpdateStmt:
		tables, where = n.Tables, n.Where
	case *ast.DeleteStmt:
		tables, where = n.Tables, n.Where
		if len(n.Using) > 0 {
			tables = n.Using
		}
	}
	if !onNullableSide(tables, tableNameOrAlias, false) {
		return "", nil
	}
	if where != nil {
		if _, ok := nodeLoc(where); !ok {
			return "", fmt.Errorf("prior backup cannot place the filter for %s's outer join", tableNameOrAlias)
		}
	}
	if len(columns.RowKeyColumns) == 0 {
		return "", fmt.Errorf("prior backup cannot tell %s's rows from the rows its outer join makes up", tableNameOrAlias)
	}
	return quoteIdentifier(tableNameOrAlias) + "." + quoteIdentifier(columns.RowKeyColumns[0]) + " IS NOT NULL", nil
}

// onNullableSide reports whether the table named name (its alias, if any) is
// on the side of an outer join that the join fills with NULLs.
func onNullableSide(exprs []ast.TableExpr, name string, nullable bool) bool {
	for _, expr := range exprs {
		switch e := expr.(type) {
		case *ast.TableRef:
			key := e.Name
			if e.Alias != "" {
				key = e.Alias
			}
			if key == name {
				return nullable
			}
		case *ast.JoinClause:
			left, right := nullable, nullable
			switch e.Type {
			case ast.JoinLeft, ast.JoinNaturalLeft:
				right = true
			case ast.JoinRight, ast.JoinNaturalRight:
				left = true
			default:
			}
			if onNullableSide([]ast.TableExpr{e.Left}, name, left) || onNullableSide([]ast.TableExpr{e.Right}, name, right) {
				return true
			}
		}
	}
	return false
}

// checkForeignKeys refuses a backup when a foreign key of another table in
// the database cascades the change to rows the backup leaves out: ON DELETE
// for a DELETE, and ON UPDATE for an UPDATE that sets a column the key
// references. A catalog without GetDatabase has no foreign keys to read.
func checkForeignKeys(opts BackupOptions, table *TableReference, targets []DMLTarget) error {
	if opts.Catalog.GetDatabase == nil {
		return nil
	}
	db, err := findDatabase(opts.Catalog, opts.CaseSensitive, table.Database)
	if err != nil {
		return fmt.Errorf("prior backup cannot read the foreign keys of %s.%s: %w", table.Database, table.Table, err)
	}
	if db == nil || !db.HasSchema() {
		return fmt.Errorf("prior backup cannot read the foreign keys of %s.%s: the database has no metadata", table.Database, table.Table)
	}
	setColumns := make(map[string]bool)
	for _, item := range targets {
		if update, ok := item.Node.(*ast.UpdateStmt); ok {
			for _, assign := range update.SetList {
				if assign.Column != nil {
					setColumns[strings.ToLower(assign.Column.Column)] = true
				}
			}
		}
	}
	for _, name := range db.ListTableNames() {
		child := db.GetTable(name)
		for _, fk := range child.GetForeignKeys() {
			if !equalIdentifier(fk.GetReferencedTable(), table.Table, opts.CaseSensitive) {
				continue
			}
			if fk.GetReferencedSchema() != "" && !equalIdentifier(fk.GetReferencedSchema(), table.Database, opts.CaseSensitive) {
				continue
			}
			action, event := fk.GetOnDelete(), "DELETE"
			if table.StatementType == DMLTypeUpdate {
				action, event = fk.GetOnUpdate(), "UPDATE"
				referenced := false
				for _, column := range fk.GetReferencedColumns() {
					referenced = referenced || setColumns[strings.ToLower(column)]
				}
				if !referenced {
					continue
				}
			}
			switch strings.ToUpper(action) {
			case "CASCADE", "SET NULL", "SET DEFAULT":
				return fmt.Errorf("prior backup cannot back up %s.%s: foreign key %s of %s is ON %s %s, and the rows it changes are not backed up", table.Database, table.Table, fk.GetName(), child.GetName(), event, strings.ToUpper(action))
			default:
			}
		}
	}
	return nil
}

// findDatabase returns the database the catalog lists under name, or nil.
func findDatabase(catalog Catalog, caseSensitive bool, name string) (Database, error) {
	if catalog.ListDatabaseNames == nil {
		return nil, errors.New("the catalog has no ListDatabaseNames")
	}
	names, err := catalog.ListDatabaseNames()
	if err != nil {
		return nil, err
	}
	for _, listed := range names {
		if equalIdentifier(listed, name, caseSensitive) {
			return catalog.GetDatabase(listed)
		}
	}
	return nil, nil
}

// joins reports whether an UPDATE or DELETE reads more than one table, so a
// row it changes may come once per match.
func joins(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return !singleTable(n.Tables)
	case *ast.DeleteStmt:
		if len(n.Using) > 0 {
			return !singleTable(n.Using)
		}
		return !singleTable(n.Tables)
	}
	return false
}

func singleTable(exprs []ast.TableExpr) bool {
	if len(exprs) != 1 {
		return false
	}
	_, ok := exprs[0].(*ast.TableRef)
	return ok
}

// quoteIdentifier quotes name in backticks, doubling any backtick in it.
func quoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// quoteColumnList quotes each column, after prefix, and joins them with sep.
func quoteColumnList(columns []string, prefix, sep string) string {
	quoted := make([]string, 0, len(columns))
	for _, column := range columns {
		quoted = append(quoted, prefix+quoteIdentifier(column))
	}
	return strings.Join(quoted, sep)
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

func (opts BackupOptions) classifyColumns(database, table string) (TableColumns, error) {
	if opts.Columns != nil {
		return opts.Columns(database, table)
	}
	return ClassifyColumns(opts.Catalog, opts.CaseSensitive, database, table)
}

// ClassifyColumns returns what a backup needs of database.table: its
// generated and other columns, in table order, and whether a key tells its
// rows apart. The database is matched against ListDatabaseNames; a database
// the catalog does not have is a *NotFoundError.
func ClassifyColumns(catalog Catalog, caseSensitive bool, database, table string) (TableColumns, error) {
	if catalog.GetDatabase == nil {
		return TableColumns{}, errors.New("the catalog has no GetDatabase")
	}

	var db Database
	allDatabaseNames, err := catalog.ListDatabaseNames()
	if err != nil {
		return TableColumns{}, fmt.Errorf("failed to list databases names: %w", err)
	}
	for _, name := range allDatabaseNames {
		if equalIdentifier(name, database, caseSensitive) {
			db, err = catalog.GetDatabase(name)
			if err != nil {
				return TableColumns{}, fmt.Errorf("failed to get database metadata for database %q: %w", name, err)
			}
			break
		}
	}
	if db == nil {
		return TableColumns{}, &NotFoundError{Database: &database}
	}

	if !db.HasSchema() {
		return TableColumns{}, errors.New("failed to get schema metadata")
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
		return TableColumns{}, fmt.Errorf("table %s not found in schema", table)
	}

	var columns TableColumns
	for _, column := range tableSchema.GetColumns() {
		if column.GetGeneration() != nil {
			columns.Generated = append(columns.Generated, column.GetName())
		} else {
			columns.Normal = append(columns.Normal, column.GetName())
		}
	}
	columns.RowKeyColumns = rowKey(tableSchema)
	columns.RowKey = columns.RowKeyColumns != nil
	return columns, nil
}

// rowKey returns the parts of a primary key, or of a unique key whose parts
// are all non-nullable columns, which tells table's rows apart, or nil. NULLs
// never collide, and an expression part may be NULL.
func rowKey(table *metadata.TableMetadata) []string {
	nullable := make(map[string]bool)
	for _, column := range table.GetColumns() {
		nullable[strings.ToLower(column.GetName())] = column.GetNullable()
	}
	for _, index := range table.GetIndexes() {
		if index.GetPrimary() {
			return index.GetExpressions()
		}
		if !index.GetUnique() || len(index.GetExpressions()) == 0 {
			continue
		}
		key := true
		for _, part := range index.GetExpressions() {
			if isNullable, ok := nullable[strings.ToLower(part)]; !ok || isNullable {
				key = false
				break
			}
		}
		if key {
			return index.GetExpressions()
		}
	}
	return nil
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
