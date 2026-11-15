package analysis

import (
	"fmt"
	"slices"
	"strings"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

const (
	maxTableNameLength = 63
)

// dmlType is the change a backed-up statement makes to its table.
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

func (t *tableReference) String() string {
	if t.Database != "" {
		return quoteIdentifier(t.Database) + "." + quoteIdentifier(t.Schema) + "." + quoteIdentifier(t.Table)
	}

	if t.Schema != "" {
		return quoteIdentifier(t.Schema) + "." + quoteIdentifier(t.Table)
	}

	return quoteIdentifier(t.Table)
}

// quoteIdentifier writes name as a quoted identifier, doubling any '"' in it.
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

type statementInfo struct {
	offset    int
	statement string
	node      ast.Node
	table     *tableReference
	start     int
	end       int
	fullSQL   string
}

// BackupOptions configures TransformDMLToSelect.
type BackupOptions struct {
	// SearchPath resolves unqualified tables; empty means public.
	SearchPath []string
	// SearchRelation returns the schema of the first relation (table, view,
	// materialized view, external table or sequence) named name in
	// searchPath's schemas, or "".
	SearchRelation func(searchPath []string, name string) string
	// TargetSchema is where the backup tables are created.
	TargetSchema string
	// TablePrefix starts each backup table's name.
	TablePrefix string
}

// BackupStatement creates the backup table of one source table, copying the
// rows the script's UPDATE and DELETE statements on it change.
type BackupStatement struct {
	Statement       string
	SourceSchema    string
	SourceTableName string
	TargetTableName string
	// StartOffset is the byte offset where the first statement backed up
	// starts, and EndOffset where the last one ends.
	StartOffset int
	EndOffset   int
}

// TransformDMLToSelect returns the statements that back up the rows a
// script's UPDATE and DELETE statements change, one per table, ordered by
// where their first statement starts. stmts is the script statement parsed,
// in order. A SET search_path in the script changes how later unqualified
// tables resolve. A statement with a data-modifying CTE, an unqualified table
// the search path does not find, or statements of different types on one
// table are errors.
func TransformDMLToSelect(stmts []omnipg.Statement, statement string, opts BackupOptions) ([]BackupStatement, error) {
	statementInfoList, err := prepareTransformation(stmts, statement, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: %w", err)
	}

	return generateSQL(statementInfoList, opts.TargetSchema, opts.TablePrefix)
}

func generateSQL(statementInfoList []statementInfo, targetSchema string, tablePrefix string) ([]BackupStatement, error) {
	groupByTable := make(map[string][]statementInfo)
	for _, item := range statementInfoList {
		key := fmt.Sprintf("%s.%s", item.table.Schema, item.table.Table)
		groupByTable[key] = append(groupByTable[key], item)
	}

	// Check if the statement type is the same for all statements on the same table.
	for key, list := range groupByTable {
		statementType := dmlTypeUnknown
		for _, item := range list {
			if statementType == dmlTypeUnknown {
				statementType = item.table.StatementType
			}
			if statementType != item.table.StatementType {
				return nil, fmt.Errorf("The statement type is not the same for all statements on the same table %q", key)
			}
		}
	}

	var result []BackupStatement
	for key, list := range groupByTable {
		backupStatement, err := generateSQLForTable(list, targetSchema, tablePrefix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL for table %q: %w", key, err)
		}
		result = append(result, *backupStatement)
	}

	slices.SortFunc(result, func(a, b BackupStatement) int {
		if a.StartOffset != b.StartOffset {
			if a.StartOffset < b.StartOffset {
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

func generateSQLForTable(statementInfoList []statementInfo, targetSchema string, tablePrefix string) (*BackupStatement, error) {
	table := statementInfoList[0].table

	targetTable := fmt.Sprintf("%s_%s_%s", tablePrefix, table.Table, table.Schema)
	targetTable, _ = truncateString(targetTable, maxTableNameLength)
	var buf strings.Builder
	if _, err := fmt.Fprintf(&buf, "CREATE TABLE %s.%s AS\n", quoteIdentifier(targetSchema), quoteIdentifier(targetTable)); err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	// In PostgreSQL, WITH clause must appear once before the first SELECT.
	cteClause := extractCTE(statementInfoList[0].node, statementInfoList[0].fullSQL)
	if cteClause != "" {
		if _, err := fmt.Fprintf(&buf, "%s\n", cteClause); err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}

	for i, item := range statementInfoList {
		if i != 0 {
			if _, err := buf.WriteString("\n  UNION\n"); err != nil {
				return nil, fmt.Errorf("failed to write to buffer: %w", err)
			}
		}
		if table.Alias != "" {
			if _, err := fmt.Fprintf(&buf, "  SELECT %s.* ", quoteIdentifier(table.Alias)); err != nil {
				return nil, fmt.Errorf("failed to write to buffer: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(&buf, `  SELECT %s.* `, table.String()); err != nil {
				return nil, fmt.Errorf("failed to write to buffer: %w", err)
			}
		}

		if err := writeSuffixSelectClause(&buf, item.node, item.fullSQL); err != nil {
			return nil, fmt.Errorf("failed to write string with new line: %w", err)
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
		StartOffset:     statementInfoList[0].start,
		EndOffset:       statementInfoList[len(statementInfoList)-1].end,
	}, nil
}

// DataModifyingCTENames returns the names of the CTEs in a statement's WITH
// clause that run INSERT, UPDATE, DELETE, or MERGE. PostgreSQL only allows
// such CTEs in the top-level WITH clause, so nested queries are not searched.
func DataModifyingCTENames(node ast.Node) []string {
	var wc *ast.WithClause
	switch n := node.(type) {
	case *ast.SelectStmt:
		wc = n.WithClause
	case *ast.InsertStmt:
		wc = n.WithClause
	case *ast.UpdateStmt:
		wc = n.WithClause
	case *ast.DeleteStmt:
		wc = n.WithClause
	case *ast.MergeStmt:
		wc = n.WithClause
	default:
	}
	if wc == nil || wc.Ctes == nil {
		return nil
	}
	var names []string
	for _, item := range wc.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		if _, isSelect := cte.Ctequery.(*ast.SelectStmt); !isSelect {
			names = append(names, cte.Ctename)
		}
	}
	return names
}

// rejectDataModifyingCTE fails a backup whose statement has a data-modifying
// CTE: the backup copies the WITH clause, so running it would apply that
// change once before the statement itself runs.
func rejectDataModifyingCTE(node ast.Node) error {
	names := DataModifyingCTENames(node)
	if len(names) == 0 {
		return nil
	}
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}
	return fmt.Errorf("prior backup does not support data-modifying WITH query %s; rewrite the statement without it, or disable or skip prior backup", strings.Join(quoted, ", "))
}

func extractCTE(node ast.Node, fullSQL string) string {
	var wc *ast.WithClause
	switch n := node.(type) {
	case *ast.UpdateStmt:
		wc = n.WithClause
	case *ast.DeleteStmt:
		wc = n.WithClause
	default:
	}
	if wc == nil {
		return ""
	}
	return extractNodeText(wc.Loc, fullSQL)
}

func writeSuffixSelectClause(buf *strings.Builder, node ast.Node, fullSQL string) error {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return writeUpdateSuffix(buf, n, fullSQL)
	case *ast.DeleteStmt:
		return writeDeleteSuffix(buf, n, fullSQL)
	}
	return nil
}

func writeUpdateSuffix(buf *strings.Builder, n *ast.UpdateStmt, sql string) error {
	if n.Relation == nil {
		return nil
	}

	relText := extractNodeText(n.Relation.Loc, sql)
	if _, err := fmt.Fprintf(buf, "FROM %s", relText); err != nil {
		return fmt.Errorf("failed to write to buffer: %w", err)
	}

	if n.FromClause != nil && n.FromClause.Len() > 0 {
		fromText := extractNodeText(ast.ListSpan(n.FromClause), sql)
		if fromText != "" {
			if _, err := fmt.Fprintf(buf, ", %s", fromText); err != nil {
				return fmt.Errorf("failed to write to buffer: %w", err)
			}
		}
	}

	if n.WhereClause != nil {
		whereText := extractNodeText(ast.NodeLoc(n.WhereClause), sql)
		if whereText != "" {
			if _, err := fmt.Fprintf(buf, " WHERE %s", whereText); err != nil {
				return fmt.Errorf("failed to write to buffer: %w", err)
			}
		}
	}
	return nil
}

func writeDeleteSuffix(buf *strings.Builder, n *ast.DeleteStmt, sql string) error {
	if n.Relation == nil {
		return nil
	}

	relText := extractNodeText(n.Relation.Loc, sql)
	if _, err := fmt.Fprintf(buf, "FROM %s", relText); err != nil {
		return fmt.Errorf("failed to write to buffer: %w", err)
	}

	if n.UsingClause != nil && n.UsingClause.Len() > 0 {
		usingText := extractNodeText(ast.ListSpan(n.UsingClause), sql)
		if usingText != "" {
			if _, err := fmt.Fprintf(buf, ", %s", usingText); err != nil {
				return fmt.Errorf("failed to write to buffer: %w", err)
			}
		}
	}

	if n.WhereClause != nil {
		whereText := extractNodeText(ast.NodeLoc(n.WhereClause), sql)
		if whereText != "" {
			if _, err := fmt.Fprintf(buf, " WHERE %s", whereText); err != nil {
				return fmt.Errorf("failed to write to buffer: %w", err)
			}
		}
	}
	return nil
}

// extractNodeText extracts trimmed text from SQL using a Loc range.
func extractNodeText(loc ast.Loc, sql string) string {
	if loc.Start < 0 || loc.End < 0 || loc.Start >= loc.End || loc.End > len(sql) {
		return ""
	}
	return strings.TrimSpace(sql[loc.Start:loc.End])
}

func prepareTransformation(stmts []omnipg.Statement, statement string, opts BackupOptions) ([]statementInfo, error) {
	searchPath := opts.SearchPath
	var dmls []statementInfo
	for i, stmt := range stmts {
		if stmt.Empty() {
			continue
		}

		switch n := stmt.AST.(type) {
		case *ast.VariableSetStmt:
			if strings.EqualFold(n.Name, "search_path") && n.Args != nil {
				var newSearchPath []string
				for _, arg := range n.Args.Items {
					if ac, ok := arg.(*ast.A_Const); ok {
						if s, ok := ac.Val.(*ast.String); ok {
							newSearchPath = append(newSearchPath, s.Str)
						}
					}
				}
				if len(newSearchPath) > 0 {
					searchPath = newSearchPath
				}
			}
		case *ast.UpdateStmt:
			if err := rejectDataModifyingCTE(n); err != nil {
				return nil, err
			}
			table, err := extractTableReferenceFromRangeVar(n.Relation, opts.SearchRelation, searchPath)
			if err != nil {
				return nil, fmt.Errorf("failed to extract table reference from update statement at offset %d: %w", i, err)
			}
			if table == nil {
				continue
			}
			table.StatementType = dmlTypeUpdate
			nodeLoc := n.Loc
			dmls = append(dmls, statementInfo{
				offset:    i,
				statement: stmt.Text,
				node:      n,
				table:     table,
				start:     nodeLoc.Start,
				end:       nodeLoc.End,
				fullSQL:   statement,
			})
		case *ast.DeleteStmt:
			if err := rejectDataModifyingCTE(n); err != nil {
				return nil, err
			}
			table, err := extractTableReferenceFromRangeVar(n.Relation, opts.SearchRelation, searchPath)
			if err != nil {
				return nil, fmt.Errorf("failed to extract table reference from delete statement at offset %d: %w", i, err)
			}
			if table == nil {
				continue
			}
			table.StatementType = dmlTypeDelete
			nodeLoc := n.Loc
			dmls = append(dmls, statementInfo{
				offset:    i,
				statement: stmt.Text,
				node:      n,
				table:     table,
				start:     nodeLoc.Start,
				end:       nodeLoc.End,
				fullSQL:   statement,
			})
		default:
		}
	}

	return dmls, nil
}

func extractTableReferenceFromRangeVar(rv *ast.RangeVar, searchRelation func(searchPath []string, name string) string, searchPath []string) (*tableReference, error) {
	if rv == nil {
		return nil, nil
	}

	table := tableReference{}

	switch {
	case rv.Catalogname != "":
		table.Database = rv.Catalogname
		table.Schema = rv.Schemaname
		table.Table = rv.Relname
	case rv.Schemaname != "":
		table.Schema = rv.Schemaname
		table.Table = rv.Relname
	default:
		// TODO: remove it in the future.
		// Handle the case where the search path is not synchronized with the metadata.
		if len(searchPath) == 0 {
			searchPath = []string{"public"}
		}
		schemaName := searchRelation(searchPath, rv.Relname)
		if schemaName == "" {
			return nil, fmt.Errorf("Table %q not found in metadata with search path %v", rv.Relname, searchPath)
		}
		table.Schema = schemaName
		table.Table = rv.Relname
	}

	if rv.Alias != nil && rv.Alias.Aliasname != "" {
		table.Alias = rv.Alias.Aliasname
	}
	return &table, nil
}

// truncateString truncates str to at most limit runes and reports whether it
// did.
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
