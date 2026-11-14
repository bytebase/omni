package analysis

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
)

// ChangedResource is a resource a statement writes.
type ChangedResource struct {
	// Database is the resource's database. An unqualified table is in the
	// current database.
	Database string
	// Table is the changed table, empty for a database-only change.
	Table string
	// DatabaseOnly marks a non-table object (view, routine, trigger) whose own
	// name carries a database qualifier. It has no table to scope, only its
	// target database.
	DatabaseOnly bool
	// AffectsData marks a table change that touches the table's existing rows,
	// such as DROP, TRUNCATE or ALTER.
	AffectsData bool
}

// StatementChanges is what one statement changes.
type StatementChanges struct {
	// Resources holds the resources the statement writes, in the order it
	// names them, duplicates included.
	Resources []ChangedResource
	// DMLStatements holds the text of a DML statement whose rows can be
	// estimated, ending with ';'.
	DMLStatements []string
	// DMLCount counts the DML statements, including those DMLStatements omits.
	DMLCount int
	// InsertCount counts the rows of INSERT ... VALUES, and one for INSERT ... SET.
	InsertCount int
}

// ExtractChangedResources returns the resources one statement writes and its
// DML counts. database is the current database, and text is the statement's
// text, which a DML statement is sampled from.
func ExtractChangedResources(node ast.Node, database, text string) StatementChanges {
	var changes StatementChanges
	addTable := func(db, table string, affectsData bool) {
		changes.Resources = append(changes.Resources, ChangedResource{Database: db, Table: table, AffectsData: affectsData})
	}
	// addObjectDatabase records a database-only write target for a non-table object DDL
	// (view/function/routine/trigger) ONLY when the object's own name carries an explicit
	// database qualifier (TableRef.Schema). Unqualified names record nothing (request-database
	// fallback). The object's own qualifier is used — never an ON-table reference.
	addObjectDatabase := func(ref *ast.TableRef) {
		if ref != nil && ref.Schema != "" {
			changes.Resources = append(changes.Resources, ChangedResource{Database: ref.Schema, DatabaseOnly: true})
		}
	}

	switch n := node.(type) {
	case *ast.CreateTableStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			// CREATE OR REPLACE drops the existing table first, so its rows go
			// with it and have to be counted. A plain CREATE affects none.
			addTable(db, table, n.OrReplace)
		}

	case *ast.DropTableStmt:
		for _, ref := range n.Tables {
			db, table := tableRef(ref, database)
			addTable(db, table, true)
		}

	case *ast.TruncateStmt:
		for _, ref := range n.Tables {
			db, table := tableRef(ref, database)
			addTable(db, table, true)
		}

	case *ast.LoadDataStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			addTable(db, table, false)
		}

	case *ast.AlterTableStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			addTable(db, table, true)
		}

	case *ast.RenameTableStmt:
		for _, pair := range n.Pairs {
			if pair.Old != nil {
				db, table := tableRef(pair.Old, database)
				addTable(db, table, false)
			}
			if pair.New != nil {
				db, table := tableRef(pair.New, database)
				addTable(db, table, false)
			}
		}

	case *ast.CreateIndexStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			addTable(db, table, false)
		}

	case *ast.DropIndexStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			addTable(db, table, false)
		}

	case *ast.InsertStmt:
		if n.Table != nil {
			db, table := tableRef(n.Table, database)
			addTable(db, table, false)
		}
		if len(n.Values) > 0 {
			changes.InsertCount += len(n.Values)
			break
		}
		// INSERT ... SET inserts exactly one row, and EXPLAIN gives no estimate for it.
		if len(n.SetList) > 0 {
			changes.InsertCount++
			break
		}
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	case *ast.UpdateStmt:
		for _, resource := range extractTableExprs(n.Tables, database) {
			addTable(resource.Database, resource.Table, false)
		}
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	case *ast.DeleteStmt:
		for _, resource := range extractTableExprs(n.Tables, database) {
			addTable(resource.Database, resource.Table, false)
		}
		for _, resource := range extractTableExprs(n.Using, database) {
			addTable(resource.Database, resource.Table, false)
		}
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	// Non-table object DDL: gate by the object's own explicit database qualifier only.
	// CREATE TRIGGER is intentionally skipped — its AST name is unqualified (a bare string);
	// inferring the database from the ON table is out of scope.
	case *ast.CreateViewStmt:
		addObjectDatabase(n.Name)
	case *ast.AlterViewStmt:
		addObjectDatabase(n.Name)
	case *ast.DropViewStmt:
		for _, ref := range n.Views {
			addObjectDatabase(ref)
		}
	case *ast.CreateFunctionStmt:
		addObjectDatabase(n.Name)
	case *ast.AlterRoutineStmt:
		addObjectDatabase(n.Name)
	case *ast.DropRoutineStmt:
		addObjectDatabase(n.Name)
	case *ast.DropTriggerStmt:
		addObjectDatabase(n.Name)

	default:
	}

	return changes
}

// tableRef extracts database and table name from an omni TableRef.
func tableRef(ref *ast.TableRef, defaultDB string) (string, string) {
	if ref == nil {
		return defaultDB, ""
	}
	db := ref.Schema
	if db == "" {
		db = defaultDB
	}
	return db, ref.Name
}

// extractTableExprs collects the tables from a slice of TableExpr.
func extractTableExprs(exprs []ast.TableExpr, defaultDB string) []ChangedResource {
	var result []ChangedResource
	for _, expr := range exprs {
		result = append(result, extractTableExpr(expr, defaultDB)...)
	}
	return result
}

// extractTableExpr recursively extracts the tables from a single TableExpr.
func extractTableExpr(expr ast.TableExpr, defaultDB string) []ChangedResource {
	if expr == nil {
		return nil
	}
	switch n := expr.(type) {
	case *ast.TableRef:
		db, table := tableRef(n, defaultDB)
		return []ChangedResource{{Database: db, Table: table}}
	case *ast.JoinClause:
		var res []ChangedResource
		res = append(res, extractTableExpr(n.Left, defaultDB)...)
		res = append(res, extractTableExpr(n.Right, defaultDB)...)
		return res
	default:
		return nil
	}
}

// statementText returns the text of a statement, ensuring it ends with a semicolon.
func statementText(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasSuffix(text, ";") {
		text += ";"
	}
	return text
}
