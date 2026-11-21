package analysis

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
)

const (
	// The default schema is 'dbo' for MSSQL.
	defaultSchema      = "dbo"
	maxTableNameLength = 128
)

type dmlType int

const (
	dmlTypeUnknown dmlType = iota
	dmlTypeUpdate
	dmlTypeInsert
	dmlTypeDelete
)

type tableReference struct {
	Database      string
	Schema        string
	Table         string
	Alias         string
	StatementType dmlType
}

type statementInfo struct {
	statement     string
	node          ast.Node
	table         *tableReference
	startPosition Position
	endPosition   Position
}

// ParsedStatement is one statement of a script, as the caller split and
// parsed it.
type ParsedStatement struct {
	Text string
	// Node is the statement's AST, or nil.
	Node ast.Node
	// StartLine is the 1-based line the statement starts on in the script.
	StartLine int
}

// Position is a 1-based line and a column in a script. Column counts the
// characters before the position on its line of the statement's own text,
// so on a statement's first line it does not count what precedes the
// statement.
type Position struct {
	Line   int
	Column int
}

// BackupStatement creates the backup table of one source table, copying the
// rows the script's UPDATE and DELETE statements on it change.
type BackupStatement struct {
	Statement       string
	SourceSchema    string
	SourceTableName string
	TargetTableName string
	// StartPosition is where the first statement backed up starts, and
	// EndPosition where the last one ends.
	StartPosition Position
	EndPosition   Position
}

// TransformDMLToSelect returns the statements that back up the rows a
// script's UPDATE and DELETE statements change into targetDatabase, one per
// table, ordered by where their first statement starts. stmts is the script,
// parsed. sourceDatabase is the database unqualified tables are in, and
// caseSensitive the instance collation. A DML on a temp table is skipped. A
// DML targeting a CTE, statements of different types on one table, or
// statements with different WITH clauses on one table are errors.
func TransformDMLToSelect(stmts []ParsedStatement, sourceDatabase, targetDatabase, tablePrefix string, caseSensitive bool) ([]BackupStatement, error) {
	statementInfoList, err := prepareTransformation(sourceDatabase, stmts, caseSensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(statementInfoList, targetDatabase, tablePrefix)
}

func generateSQL(statementInfoList []statementInfo, targetDatabase string, tablePrefix string) ([]BackupStatement, error) {
	// Tables are visited in the order the script first changes them, so an
	// error names the same table on every run.
	groupByTable := make(map[string][]statementInfo)
	var tableOrder []string
	for _, item := range statementInfoList {
		key := fmt.Sprintf("%s.%s.%s", item.table.Database, item.table.Schema, item.table.Table)
		if _, ok := groupByTable[key]; !ok {
			tableOrder = append(tableOrder, key)
		}
		groupByTable[key] = append(groupByTable[key], item)
	}

	// Check if the statement type is the same for all statements on the same table.
	for _, key := range tableOrder {
		list := groupByTable[key]
		statementType := dmlTypeUnknown
		for _, item := range list {
			if statementType == dmlTypeUnknown {
				statementType = item.table.StatementType
			}
			if statementType != item.table.StatementType {
				return nil, fmt.Errorf("prior backup cannot handle mixed DMLs on the same table %s", key)
			}
		}
	}

	var result []BackupStatement
	for _, key := range tableOrder {
		list := groupByTable[key]
		backupStatement, err := generateSQLForTable(list, targetDatabase, tablePrefix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL for table %s: %w", key, err)
		}
		result = append(result, *backupStatement)
	}

	slices.SortFunc(result, func(a, b BackupStatement) int {
		if a.StartPosition.Line != b.StartPosition.Line {
			if a.StartPosition.Line < b.StartPosition.Line {
				return -1
			}
			return 1
		}
		if a.StartPosition.Column != b.StartPosition.Column {
			if a.StartPosition.Column < b.StartPosition.Column {
				return -1
			}
			return 1
		}
		if a.SourceTableName < b.SourceTableName {
			return -1
		}
		if a.SourceTableName > b.SourceTableName {
			return 1
		}
		return 0
	})

	return result, nil
}

func generateSQLForTable(statementInfoList []statementInfo, targetDatabase string, tablePrefix string) (*BackupStatement, error) {
	table := statementInfoList[0].table

	targetTable := fmt.Sprintf("%s_%s_%s", tablePrefix, table.Table, table.Database)
	targetTable, _ = truncateString(targetTable, maxTableNameLength)
	var buf strings.Builder
	cteClause, err := backupCTEClause(statementInfoList)
	if err != nil {
		return nil, err
	}
	if cteClause != "" {
		if _, err := fmt.Fprintf(&buf, "%s\n", cteClause); err != nil {
			return nil, fmt.Errorf("failed to write buffer: %w", err)
		}
	}
	if _, err := fmt.Fprintf(&buf, "SELECT * INTO [%s].[%s].[%s] FROM (\n", targetDatabase, defaultSchema, targetTable); err != nil {
		return nil, fmt.Errorf("failed to write buffer: %w", err)
	}
	var optionClauses []string
	for i, item := range statementInfoList {
		if i > 0 {
			if _, err := buf.WriteString("\n  UNION\n"); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
		}
		topClause, fromClause, optionClause, err := extractSuffixSelectStatement(item.node, item.statement)
		if err != nil {
			return nil, fmt.Errorf("failed to extract suffix select statement: %w", err)
		}
		if optionClause != "" && !slices.Contains(optionClauses, optionClause) {
			optionClauses = append(optionClauses, optionClause)
		}
		if _, err := buf.WriteString("  SELECT "); err != nil {
			return nil, fmt.Errorf("failed to write buffer: %w", err)
		}
		if len(topClause) > 0 {
			if _, err := buf.WriteString(topClause); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
			if _, err := buf.WriteString(" "); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
		}
		if len(item.table.Alias) == 0 {
			if _, err := fmt.Fprintf(&buf, "[%s].[%s].[%s].* ", item.table.Database, item.table.Schema, item.table.Table); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(&buf, "[%s].* ", item.table.Alias); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
		}
		if len(fromClause) > 0 {
			if _, err := buf.WriteString(fromClause); err != nil {
				return nil, fmt.Errorf("failed to write buffer: %w", err)
			}
		}
	}
	if _, err := buf.WriteString(") AS backup_table"); err != nil {
		return nil, fmt.Errorf("failed to write buffer: %w", err)
	}
	if outerOption := mergeOptionClauses(optionClauses); outerOption != "" {
		if _, err := buf.WriteString(" "); err != nil {
			return nil, fmt.Errorf("failed to write buffer: %w", err)
		}
		if _, err := buf.WriteString(outerOption); err != nil {
			return nil, fmt.Errorf("failed to write buffer: %w", err)
		}
	}
	if _, err := buf.WriteString(";"); err != nil {
		return nil, fmt.Errorf("failed to write buffer: %w", err)
	}
	return &BackupStatement{
		Statement:       buf.String(),
		SourceSchema:    table.Schema,
		SourceTableName: table.Table,
		TargetTableName: targetTable,
		StartPosition:   statementInfoList[0].startPosition,
		EndPosition:     statementInfoList[len(statementInfoList)-1].endPosition,
	}, nil
}

// backupCTEClause returns the WITH clause to place ahead of the backup SELECT
// INTO. T-SQL only allows a CTE list at the top of a statement, so a CTE
// referenced inside a UNION branch must be hoisted to the shared clause. The
// hoisted clause becomes the name scope of every branch, so it is only safe
// when each statement in the group carries the same clause text; a branch
// written without that clause could otherwise resolve a base table name to a
// CTE.
func backupCTEClause(statementInfoList []statementInfo) (string, error) {
	var clause string
	for i, item := range statementInfoList {
		withClause := dmlWithClause(item.node)
		var text string
		if withClause != nil {
			if withClause.XmlNamespaces != nil {
				return "", errors.New("prior backup does not support WITH XMLNAMESPACES")
			}
			text = sourceFromLoc(item.statement, withClause.Loc)
		}
		if i == 0 {
			clause = text
			continue
		}
		if text != clause {
			return "", errors.New("prior backup cannot handle statements with different WITH clauses on the same table")
		}
	}
	return clause, nil
}

func dmlWithClause(node ast.Node) *ast.WithClause {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return n.WithClause
	case *ast.DeleteStmt:
		return n.WithClause
	default:
		return nil
	}
}

func extractSuffixSelectStatement(node ast.Node, source string) (string, string, string, error) {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		if _, ok := n.WhereClause.(*ast.CurrentOfExpr); ok {
			return "", "", "", errors.New("UPDATE statement with CURSOR clause is not supported")
		}
		stmtLoc := dmlNodeLoc(n)
		fromSource, searchStart := dmlFromSource(source, n.Relation, n.FromClause, stmtLoc, n.WhereClause, n.OptionClause)
		fromClause, optionClause := buildDMLFromAndOption(source, fromSource, searchStart, stmtLoc, n.WhereClause, n.OptionClause)
		return sourceFromLoc(source, dmlNodeLoc(n.Top)), fromClause, optionClause, nil
	case *ast.DeleteStmt:
		if _, ok := n.WhereClause.(*ast.CurrentOfExpr); ok {
			return "", "", "", errors.New("DELETE statement with CURSOR clause is not supported")
		}
		stmtLoc := dmlNodeLoc(n)
		fromSource, searchStart := dmlFromSource(source, n.Relation, n.FromClause, stmtLoc, n.WhereClause, n.OptionClause)
		fromClause, optionClause := buildDMLFromAndOption(source, fromSource, searchStart, stmtLoc, n.WhereClause, n.OptionClause)
		return sourceFromLoc(source, dmlNodeLoc(n.Top)), fromClause, optionClause, nil
	default:
		return "", "", "", nil
	}
}

func mergeOptionClauses(clauses []string) string {
	switch len(clauses) {
	case 0:
		return ""
	case 1:
		return clauses[0]
	default:
		var hints []string
		for _, clause := range clauses {
			open := strings.IndexByte(clause, '(')
			close := strings.LastIndexByte(clause, ')')
			if open < 0 || close <= open {
				continue
			}
			inner := strings.TrimSpace(clause[open+1 : close])
			if inner != "" && !slices.Contains(hints, inner) {
				hints = append(hints, inner)
			}
		}
		if len(hints) == 0 {
			return clauses[0]
		}
		return "OPTION (" + strings.Join(hints, ", ") + ")"
	}
}

func prepareTransformation(databaseName string, parsedStatements []ParsedStatement, caseSensitive bool) ([]statementInfo, error) {
	var dmls []statementInfo
	// otherDatabase is the database a USE switched to, when it is not
	// databaseName: the backup reads databaseName's tables, so a later DML,
	// which runs in the other one, cannot be backed up.
	otherDatabase := ""
	for _, parsedStatement := range parsedStatements {
		node := parsedStatement.Node
		if node == nil {
			continue
		}
		if use, ok := node.(*ast.UseStmt); ok {
			otherDatabase = use.Database
			if use.Database == databaseName || !caseSensitive && strings.EqualFold(use.Database, databaseName) {
				otherDatabase = ""
			}
			continue
		}
		var (
			table         *tableReference
			targetRef     *ast.TableRef
			statementType dmlType
			err           error
		)
		switch n := node.(type) {
		case *ast.UpdateStmt:
			table, targetRef, err = resolveDMLTargetTable(n.Relation, n.FromClause, databaseName)
			statementType = dmlTypeUpdate
		case *ast.DeleteStmt:
			table, targetRef, err = resolveDMLTargetTable(n.Relation, n.FromClause, databaseName)
			statementType = dmlTypeDelete
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
		if table == nil {
			return nil, fmt.Errorf("failed to resolve DML target table")
		}
		if strings.HasPrefix(table.Table, "#") {
			continue
		}
		if otherDatabase != "" {
			return nil, fmt.Errorf("prior backup does not support a DML after USE %s, as the backup reads the tables of %s; remove the USE, or disable or skip prior backup", otherDatabase, databaseName)
		}
		if isCTETarget(targetRef, dmlWithClause(node), caseSensitive) {
			return nil, fmt.Errorf("prior backup does not support DML targeting CTE %q, target the base table directly", table.Table)
		}
		table.StatementType = statementType
		loc := dmlNodeLoc(node)
		dmls = append(dmls, statementInfo{
			statement:     parsedStatement.Text,
			node:          node,
			table:         table,
			startPosition: positionFromByteOffset(parsedStatement.StartLine, parsedStatement.Text, loc.Start),
			endPosition:   positionFromByteOffset(parsedStatement.StartLine, parsedStatement.Text, dmlEndOffset(parsedStatement.Text, loc)),
		})
	}

	return dmls, nil
}

// resolveDMLTargetTable returns the DML target and the TableRef it was
// resolved from. When the target is an alias declared in the FROM clause, both
// describe the aliased physical table.
func resolveDMLTargetTable(relation ast.TableExpr, fromClause *ast.List, databaseName string) (*tableReference, *ast.TableRef, error) {
	ref, ok := relation.(*ast.TableRef)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported DML target table source %T", relation)
	}
	table := tableReferenceFromTableRef(ref, databaseName, defaultSchema)
	if fromClause != nil && table.Database == databaseName && table.Schema == defaultSchema {
		if physical, physicalRef := findPhysicalTableForAlias(fromClause, table); physical != nil {
			return physical, physicalRef, nil
		}
	}
	return table, ref, nil
}

func tableReferenceFromTableRef(ref *ast.TableRef, defaultDatabase, defaultSchema string) *tableReference {
	schemaName := defaultSchema
	if ref.Schema != "" {
		schemaName = ref.Schema
	}
	databaseName := defaultDatabase
	if ref.Database != "" {
		databaseName = ref.Database
	}
	return &tableReference{
		Database: databaseName,
		Schema:   schemaName,
		Table:    ref.Object,
		Alias:    ref.Alias,
	}
}

func findPhysicalTableForAlias(list *ast.List, table *tableReference) (*tableReference, *ast.TableRef) {
	if list == nil || table == nil {
		return nil, nil
	}
	for _, item := range list.Items {
		if result, ref := findPhysicalTableForAliasInNode(item, table); result != nil {
			return result, ref
		}
	}
	return nil, nil
}

func findPhysicalTableForAliasInNode(node ast.Node, table *tableReference) (*tableReference, *ast.TableRef) {
	switch n := node.(type) {
	case *ast.TableRef:
		if n.Alias != "" && n.Alias == table.Table {
			result := tableReferenceFromTableRef(n, table.Database, table.Schema)
			result.Alias = n.Alias
			return result, n
		}
	case *ast.AliasedTableRef:
		if ref, ok := n.Table.(*ast.TableRef); ok && n.Alias == table.Table {
			result := tableReferenceFromTableRef(ref, table.Database, table.Schema)
			result.Alias = n.Alias
			return result, ref
		}
	case *ast.JoinClause:
		if result, ref := findPhysicalTableForAliasInNode(n.Left, table); result != nil {
			return result, ref
		}
		return findPhysicalTableForAliasInNode(n.Right, table)
	default:
	}
	return nil, nil
}

func dmlFromSource(source string, relation ast.TableExpr, fromClause *ast.List, stmtLoc ast.Loc, where ast.ExprNode, option *ast.List) (string, int) {
	if fromClause != nil && len(fromClause.Items) > 0 {
		loc := listLoc(fromClause)
		end := dmlTailStart(source, loc.End, stmtLoc, where, option)
		if end < 0 {
			end = trimDMLStatementEnd(source, stmtLoc.End)
		}
		return strings.TrimSpace(source[loc.Start:end]), end
	}
	loc := dmlNodeLoc(relation)
	return sourceFromLoc(source, loc), loc.End
}

func buildDMLFromAndOption(source, fromSource string, searchStart int, stmtLoc ast.Loc, where ast.ExprNode, option *ast.List) (string, string) {
	if fromSource == "" {
		return "", ""
	}
	end := trimDMLStatementEnd(source, stmtLoc.End)
	optionStart := dmlOptionStart(source, searchStart, end, option)
	var optionClause string
	if optionStart >= 0 && optionStart < end {
		optionClause = strings.TrimSpace(source[optionStart:end])
		end = optionStart
	}
	var whereClause string
	if where != nil {
		if whereStart := findKeywordBefore(source, "WHERE", dmlNodeLoc(where).Start, searchStart); whereStart >= 0 && whereStart < end {
			whereClause = strings.TrimSpace(source[whereStart:end])
		}
	}
	if whereClause == "" {
		return "FROM " + fromSource, optionClause
	}
	return "FROM " + fromSource + " " + whereClause, optionClause
}

func dmlOptionStart(source string, searchStart, end int, option *ast.List) int {
	if option == nil {
		return -1
	}
	if loc := listLoc(option); loc.Start >= 0 {
		if start := findKeywordBefore(source, "OPTION", loc.Start, searchStart); start >= 0 {
			return start
		}
	}
	return findKeywordAfter(source, "OPTION", searchStart, end)
}

func dmlTailStart(source string, searchStart int, stmtLoc ast.Loc, where ast.ExprNode, option *ast.List) int {
	if where != nil {
		if start := findKeywordBefore(source, "WHERE", dmlNodeLoc(where).Start, searchStart); start >= 0 {
			return start
		}
	}
	return dmlOptionStart(source, searchStart, trimDMLStatementEnd(source, stmtLoc.End), option)
}

func trimDMLStatementEnd(source string, end int) int {
	if end > len(source) {
		end = len(source)
	}
	for end > 0 {
		switch source[end-1] {
		case ';', ' ', '\t', '\r', '\n':
			end--
		default:
			return end
		}
	}
	return end
}

func findKeywordBefore(source, keyword string, before int, after int) int {
	if before > len(source) {
		before = len(source)
	}
	if after < 0 {
		after = 0
	}
	return findKeywordOutsideSQL(source, keyword, after, before, true)
}

func findKeywordAfter(source, keyword string, start int, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(source) {
		end = len(source)
	}
	return findKeywordOutsideSQL(source, keyword, start, end, false)
}

func findKeywordOutsideSQL(source, keyword string, start, end int, last bool) int {
	if start >= end {
		return -1
	}
	result := -1
	for i := start; i < end; {
		switch source[i] {
		case '\'':
			i = skipQuotedSQL(source, i, end, '\'')
		case '"':
			i = skipQuotedSQL(source, i, end, '"')
		case '[':
			i = skipQuotedSQL(source, i, end, ']')
		case '-':
			if i+1 < end && source[i+1] == '-' {
				i = skipLineComment(source, i+2, end)
				continue
			}
			if keywordMatchAt(source, keyword, i, start, end) {
				if !last {
					return i
				}
				result = i
				i += len(keyword)
				continue
			}
			i++
		case '/':
			if i+1 < end && source[i+1] == '*' {
				i = skipBlockComment(source, i+2, end)
				continue
			}
			if keywordMatchAt(source, keyword, i, start, end) {
				if !last {
					return i
				}
				result = i
				i += len(keyword)
				continue
			}
			i++
		default:
			if keywordMatchAt(source, keyword, i, start, end) {
				if !last {
					return i
				}
				result = i
				i += len(keyword)
				continue
			}
			i++
		}
	}
	return result
}

func keywordMatchAt(source, keyword string, offset, start, end int) bool {
	if offset+len(keyword) > end || !strings.EqualFold(source[offset:offset+len(keyword)], keyword) {
		return false
	}
	if offset > start && isSQLIdentifierChar(source[offset-1]) {
		return false
	}
	if offset+len(keyword) < end && isSQLIdentifierChar(source[offset+len(keyword)]) {
		return false
	}
	return true
}

func isSQLIdentifierChar(ch byte) bool {
	return ch == '_' || ch == '@' || ch == '#' || ch == '$' ||
		(ch >= '0' && ch <= '9') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z')
}

func skipQuotedSQL(source string, start, end int, closeCh byte) int {
	i := start + 1
	for i < end {
		if source[i] == closeCh {
			if i+1 < end && source[i+1] == closeCh {
				i += 2
				continue
			}
			return i + 1
		}
		i++
	}
	return end
}

func skipLineComment(source string, start, end int) int {
	for start < end && source[start] != '\n' {
		start++
	}
	return start
}

func skipBlockComment(source string, start, end int) int {
	depth := 1
	for start+1 < end {
		switch {
		case source[start] == '/' && source[start+1] == '*':
			depth++
			start += 2
		case source[start] == '*' && source[start+1] == '/':
			depth--
			start += 2
			if depth == 0 {
				return start
			}
		default:
			start++
		}
	}
	return end
}

func sourceFromLoc(source string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End < 0 || loc.Start >= len(source) {
		return ""
	}
	end := min(loc.End, len(source))
	return strings.TrimSpace(source[loc.Start:end])
}

func positionFromByteOffset(startLine int, source string, offset int) Position {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(source) && len(source) > 0 {
		offset = len(source) - 1
	}
	line, column := calculateLineAndColumn(source, offset)
	return Position{
		Line:   startLine + line,
		Column: column,
	}
}

// calculateLineAndColumn returns the 0-based line and the column, in runes,
// of a byte offset. \r\n is one line break and a lone \r is one too.
func calculateLineAndColumn(statement string, byteOffset int) (line, column int) {
	if byteOffset > len(statement) {
		byteOffset = len(statement)
	}
	s := statement[:byteOffset]
	for i, r := range s {
		switch r {
		case '\r':
			line++
			column = 0
		case '\n':
			if i > 0 && s[i-1] == '\r' {
				continue
			}
			line++
			column = 0
		default:
			column++
		}
	}
	return line, column
}

// truncateString returns str's first limit characters, and whether it cut
// anything.
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

func dmlEndOffset(source string, loc ast.Loc) int {
	if loc.End >= 0 && loc.End < len(source) && source[loc.End] == ';' {
		return loc.End
	}
	return loc.End - 1
}

func listLoc(list *ast.List) ast.Loc {
	if list == nil || len(list.Items) == 0 {
		return ast.NoLoc()
	}
	return ast.Loc{
		Start: dmlNodeLoc(list.Items[0]).Start,
		End:   dmlNodeLoc(list.Items[len(list.Items)-1]).End,
	}
}

func dmlNodeLoc(node ast.Node) ast.Loc {
	if isNil(node) {
		return ast.NoLoc()
	}
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return n.Loc
	case *ast.DeleteStmt:
		return n.Loc
	case *ast.TopClause:
		return n.Loc
	case *ast.TableRef:
		return n.Loc
	case *ast.TableVarRef:
		return n.Loc
	case *ast.TableVarMethodCallRef:
		return n.Loc
	case *ast.AliasedTableRef:
		return n.Loc
	case *ast.JoinClause:
		loc := n.Loc
		leftLoc := dmlNodeLoc(n.Left)
		if leftLoc.Start >= 0 && (loc.Start < 0 || leftLoc.Start < loc.Start) {
			loc.Start = leftLoc.Start
		}
		return loc
	case *ast.SetExpr:
		return n.Loc
	default:
		loc := nodeLoc(node)
		if loc.Start >= 0 {
			return loc
		}
		return ast.NoLoc()
	}
}

func isNil(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
