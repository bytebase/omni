package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

const maxCommentLength = 1000

// RestoreOptions configures RestoreSQL.
type RestoreOptions struct {
	// Statement is the SQL that changed the table, quoted in a comment that
	// opens the restore SQL.
	Statement string
	// Database and Table are the changed table.
	Database string
	Table    string
	// BackupDatabase and BackupTable hold the backed-up rows.
	BackupDatabase string
	BackupTable    string
	// Columns are the table's columns other than generated ones, which the
	// restore copies back.
	Columns []string
	// Catalog is the instance's metadata. An UPDATE's restore reads the
	// table's unique keys from it.
	Catalog Catalog
}

// FirstDML returns the first UPDATE or DELETE of statement, or nil. It returns
// the parser's error as is.
func FirstDML(statement string) (ast.Node, error) {
	stmtList, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	if stmtList == nil {
		return nil, nil
	}
	for _, item := range stmtList.Items {
		switch item.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			return item, nil
		default:
		}
	}
	return nil, nil
}

// ReferencesTable reports whether an UPDATE or DELETE of statement names
// database.table. An unqualified table is in database. It returns the
// parser's error as is.
func ReferencesTable(statement, database, table string) (bool, error) {
	stmtList, err := parser.Parse(statement)
	if err != nil {
		return false, err
	}
	for _, node := range stmtList.Items {
		if containsTable(node, database, table) {
			return true, nil
		}
	}
	return false, nil
}

// RestoreSQL returns the SQL that copies the backed-up rows of node, an UPDATE
// or DELETE, back into the table. A DELETE's rows are inserted again; an
// UPDATE's overwrite the changed columns, which needs a unique key the UPDATE
// leaves alone.
func RestoreSQL(node ast.Node, opts RestoreOptions) (string, error) {
	var result string
	switch n := node.(type) {
	case *ast.DeleteStmt:
		result = generateDeleteRestore(opts.Database, opts.Table, opts.BackupDatabase, opts.BackupTable, opts.Columns)
	case *ast.UpdateStmt:
		r, err := generateUpdateRestore(opts.Catalog, n, opts.Database, opts.Table, opts.BackupDatabase, opts.BackupTable, opts.Columns)
		if err != nil {
			return "", err
		}
		result = r
	default:
		return "", fmt.Errorf("unexpected statement type: %T", node)
	}

	sqlForComment, truncated := truncateString(opts.Statement, maxCommentLength)
	if truncated {
		sqlForComment += "..."
	}
	// A "*/" in the statement would end the comment early.
	sqlForComment = strings.ReplaceAll(sqlForComment, "*/", "* /")
	return fmt.Sprintf("/*\nOriginal SQL:\n%s\n*/\n%s", sqlForComment, result), nil
}

func generateDeleteRestore(originalDatabase, originalTable, backupDatabase, backupTable string, normalColumns []string) string {
	quotedColumnList := quoteColumns(normalColumns)
	return fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) SELECT %s FROM `%s`.`%s`;", originalDatabase, originalTable, quotedColumnList, quotedColumnList, backupDatabase, backupTable)
}

func generateUpdateRestore(catalog Catalog, stmt *ast.UpdateStmt, originalDatabase, originalTable, backupDatabase, backupTable string, normalColumns []string) (string, error) {
	// A self-join names the table more than once; a SET column may be
	// qualified by any of those references.
	var targets []*TableReference
	for _, table := range extractSingleTablesFromTableExprs(originalDatabase, stmt.Tables) {
		if strings.EqualFold(table.Table, originalTable) && strings.EqualFold(table.Database, originalDatabase) {
			targets = append(targets, table)
		}
	}

	// Extract update column names that belong to the original table.
	updateColumns := extractUpdateColumns(stmt.SetList, originalDatabase, targets, normalColumns)

	has, err := hasDisjointUniqueKey(catalog, originalDatabase, originalTable, updateColumns)
	if err != nil {
		return "", err
	}
	if !has {
		return "", fmt.Errorf("no disjoint unique key found for %s.%s", originalDatabase, originalTable)
	}

	var buf strings.Builder
	quotedColumnList := quoteColumns(normalColumns)
	fmt.Fprintf(&buf, "INSERT INTO `%s`.`%s` (%s) SELECT %s FROM `%s`.`%s` ON DUPLICATE KEY UPDATE ", originalDatabase, originalTable, quotedColumnList, quotedColumnList, backupDatabase, backupTable)

	for i, field := range updateColumns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "`%s` = VALUES(`%s`)", field, field)
	}
	buf.WriteString(";")
	return buf.String(), nil
}

func quoteColumns(columns []string) string {
	var quotedColumns []string
	for _, column := range columns {
		quotedColumns = append(quotedColumns, fmt.Sprintf("`%s`", column))
	}
	return strings.Join(quotedColumns, ", ")
}

// extractSingleTablesFromTableExprs walks omni TableExpr nodes and returns
// a map of alias-or-name to TableReference for simple table references.
func extractSingleTablesFromTableExprs(databaseName string, exprs []ast.TableExpr) map[string]*TableReference {
	result := make(map[string]*TableReference)
	for _, expr := range exprs {
		collectTableRefs(databaseName, expr, result)
	}
	return result
}

func collectTableRefs(databaseName string, expr ast.TableExpr, result map[string]*TableReference) {
	switch n := expr.(type) {
	case *ast.TableRef:
		database := n.Schema
		if database == "" {
			database = databaseName
		}
		table := &TableReference{
			Database: database,
			Table:    n.Name,
			Alias:    n.Alias,
		}
		if n.Alias != "" {
			result[n.Alias] = table
		} else {
			result[n.Name] = table
		}
	case *ast.JoinClause:
		collectTableRefs(databaseName, n.Left, result)
		collectTableRefs(databaseName, n.Right, result)
	default:
	}
}

// extractUpdateColumns extracts column names from SET assignments that belong
// to the original table.
func extractUpdateColumns(setList []*ast.Assignment, database string, targets []*TableReference, normalColumns []string) []string {
	var result []string
	for _, assignment := range setList {
		col := assignment.Column
		if col == nil {
			continue
		}

		if col.Schema != "" && !strings.EqualFold(col.Schema, database) {
			continue
		}

		if col.Table == "" {
			// Unqualified column: check if it belongs to the normalColumns set.
			for _, c := range normalColumns {
				if strings.EqualFold(c, col.Column) {
					result = append(result, col.Column)
					break
				}
			}
			continue
		}

		for _, target := range targets {
			if col.Table == target.Alias || strings.EqualFold(col.Table, target.Table) {
				result = append(result, col.Column)
				break
			}
		}
	}
	return result
}

func hasDisjointUniqueKey(catalog Catalog, originalDatabase, originalTable string, updateColumns []string) (bool, error) {
	columnMap := make(map[string]bool)
	for _, column := range updateColumns {
		columnMap[strings.ToLower(column)] = true
	}

	if catalog.GetDatabase == nil {
		return false, errors.New("the catalog has no GetDatabase")
	}

	database, err := catalog.GetDatabase(originalDatabase)
	if err != nil {
		return false, fmt.Errorf("failed to get database metadata for %s: %w", originalDatabase, err)
	}
	if database == nil {
		return false, fmt.Errorf("database metadata is nil for %s", originalDatabase)
	}

	if !database.HasSchema() {
		return false, fmt.Errorf("schema is nil for %s", originalDatabase)
	}

	tableMetadata := database.GetTable(originalTable)
	if tableMetadata == nil {
		return false, fmt.Errorf("table metadata is nil for %s.%s", originalDatabase, originalTable)
	}

	regularColumns := make(map[string]bool)
	for _, column := range tableMetadata.Columns {
		if column == nil || column.Generation != nil {
			continue
		}
		regularColumns[strings.ToLower(column.Name)] = true
	}

	for _, index := range tableMetadata.Indexes {
		if !index.Primary && !index.Unique {
			continue
		}
		// A key on a generated column or an expression can change with
		// the columns it reads, which string comparison with the SET
		// columns cannot see, so it never identifies a row to restore.
		// Neither does a key with no parts, which is vacuously disjoint.
		if len(index.Expressions) == 0 || !allRegular(index.Expressions, regularColumns) {
			continue
		}
		if disjoint(index.Expressions, columnMap) {
			return true, nil
		}
	}

	return false, nil
}

// allRegular reports whether every key part is a column that is not
// generated.
func allRegular(parts []string, regularColumns map[string]bool) bool {
	for _, part := range parts {
		if !regularColumns[strings.ToLower(part)] {
			return false
		}
	}
	return true
}

func disjoint(a []string, b map[string]bool) bool {
	for _, item := range a {
		if _, ok := b[strings.ToLower(item)]; ok {
			return false
		}
	}
	return true
}

// containsTable checks whether the given AST node references the specified table.
func containsTable(node ast.Node, database, table string) bool {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		for _, expr := range n.Tables {
			if tableExprReferences(expr, database, table) {
				return true
			}
		}
	case *ast.DeleteStmt:
		for _, expr := range n.Tables {
			if tableExprReferences(expr, database, table) {
				return true
			}
		}
		for _, expr := range n.Using {
			if tableExprReferences(expr, database, table) {
				return true
			}
		}
	default:
	}
	return false
}

func tableExprReferences(expr ast.TableExpr, database, table string) bool {
	switch n := expr.(type) {
	case *ast.TableRef:
		db := n.Schema
		if db == "" {
			db = database
		}
		return db == database && strings.EqualFold(n.Name, table)
	case *ast.JoinClause:
		return tableExprReferences(n.Left, database, table) || tableExprReferences(n.Right, database, table)
	}
	return false
}
