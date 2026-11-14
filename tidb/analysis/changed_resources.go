package analysis

import (
	"slices"
	"strings"
	"unicode"

	"github.com/bytebase/omni/tidb/ast"
)

// ChangedResource is a resource a statement writes.
type ChangedResource struct {
	// Database is the resource's database. An unqualified table is in the
	// current database.
	Database string
	// Table is the changed table, empty for a database-only change.
	Table string
	// DatabaseOnly marks a non-table object (view, sequence) whose own name
	// carries a database qualifier. It has no table to scope, only its target
	// database.
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
	// estimated, without its trailing ';'. For a BATCH statement it is the
	// DML inside it.
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
	recordResourceChanges(database, node, &changes)
	countDML(&changes, node, text)
	return changes
}

// recordResourceChanges records the tables and databases a statement changes.
//
// A table-level DDL records the table it changes (DROP/ALTER/TRUNCATE also mark
// its data affected). A DML statement records the tables it writes. A non-table
// object DDL (view/sequence) records only its target database, and only when its
// name carries an explicit qualifier; an unqualified name records nothing and
// falls back to the request database.
func recordResourceChanges(database string, node ast.Node, changes *StatementChanges) {
	addTable := func(ref *ast.TableRef, affectsData bool) {
		if ref == nil {
			return
		}
		d := ref.Schema
		if d == "" {
			d = database
		}
		changes.Resources = append(changes.Resources, ChangedResource{Database: d, Table: ref.Name, AffectsData: affectsData})
	}
	objectDatabase := func(ref *ast.TableRef) {
		if ref != nil && ref.Schema != "" {
			changes.Resources = append(changes.Resources, ChangedResource{Database: ref.Schema, DatabaseOnly: true})
		}
	}

	switch n := node.(type) {
	case *ast.CreateTableStmt:
		addTable(n.Table, false)
	case *ast.DropTableStmt:
		for _, name := range n.Tables {
			addTable(name, true)
		}
	case *ast.DropViewStmt:
		for _, view := range n.Views {
			objectDatabase(view)
		}
	case *ast.AlterTableStmt:
		addTable(n.Table, true)
	case *ast.TruncateStmt:
		for _, name := range n.Tables {
			addTable(name, true)
		}
	case *ast.LoadDataStmt:
		addTable(n.Table, false)
	case *ast.ImportIntoStmt:
		addTable(n.Table, false)
	case *ast.RenameTableStmt:
		for _, pair := range n.Pairs {
			addTable(pair.Old, false)
			addTable(pair.New, false)
		}
	case *ast.CreateIndexStmt:
		addTable(n.Table, false)
	case *ast.DropIndexStmt:
		addTable(n.Table, false)
	case *ast.CreateViewStmt:
		objectDatabase(n.Name)
	case *ast.CreateSequenceStmt:
		objectDatabase(n.Name)
	case *ast.AlterSequenceStmt:
		objectDatabase(n.Name)
	case *ast.DropSequenceStmt:
		for _, seq := range n.Sequences {
			objectDatabase(seq)
		}
	case *ast.InsertStmt:
		if n.Table != nil {
			addDMLTables(addTable, []ast.TableExpr{n.Table})
		}
	case *ast.UpdateStmt:
		addDMLTables(addTable, n.Tables)
	case *ast.DeleteStmt:
		addDMLTables(addTable, slices.Concat(n.Tables, n.Using))
	case *ast.BatchStmt:
		if n.DryRun == ast.BatchDryRunNone {
			recordResourceChanges(database, n.DML, changes)
		}
	default:
	}
}

// addDMLTables records every base table named by a DML statement's table
// expressions, descending through joins. A derived table (subquery) names no
// base table.
func addDMLTables(addTable func(*ast.TableRef, bool), tableExprs []ast.TableExpr) {
	var add func(expr ast.TableExpr)
	add = func(expr ast.TableExpr) {
		switch e := expr.(type) {
		case *ast.TableRef:
			addTable(e, false)
		case *ast.JoinClause:
			add(e.Left)
			add(e.Right)
		default:
		}
	}
	for _, expr := range tableExprs {
		add(expr)
	}
}

func countDML(changes *StatementChanges, node ast.Node, text string) {
	switch n := node.(type) {
	case *ast.InsertStmt:
		switch {
		case len(n.Values) > 0:
			changes.InsertCount += len(n.Values)
		case len(n.SetList) > 0:
			changes.InsertCount++
		default:
			addSampleDML(changes, text)
		}
	case *ast.UpdateStmt, *ast.DeleteStmt:
		addSampleDML(changes, text)
	case *ast.BatchStmt:
		if n.DryRun == ast.BatchDryRunNone {
			addBatchSampleDML(changes, n, text)
		}
	default:
	}
}

// addBatchSampleDML samples the DML inside a BATCH statement, because EXPLAIN
// rejects the BATCH syntax.
func addBatchSampleDML(changes *StatementChanges, batch *ast.BatchStmt, text string) {
	// The parser reads the body of an executable comment, so the DML's location can fall inside one
	// and the slice would leave the comment unclosed.
	if strings.Contains(text, "/*!") || strings.Contains(text, "/*T!") || strings.Contains(text, "/*M!") {
		changes.DMLCount++
		return
	}
	var loc ast.Loc
	switch dml := batch.DML.(type) {
	case *ast.DeleteStmt:
		loc = dml.Loc
	case *ast.UpdateStmt:
		loc = dml.Loc
	case *ast.InsertStmt:
		loc = dml.Loc
	default:
		changes.DMLCount++
		return
	}
	end := loc.End
	if end <= loc.Start || end > len(text) {
		end = len(text)
	}
	if loc.Start <= 0 || loc.Start >= end {
		changes.DMLCount++
		return
	}
	addSampleDML(changes, text[loc.Start:end])
}

func addSampleDML(changes *StatementChanges, text string) {
	changes.DMLCount++
	changes.DMLStatements = append(changes.DMLStatements, trimStatement(text))
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
