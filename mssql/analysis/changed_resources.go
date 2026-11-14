package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/mssql/ast"
)

// ChangedResource is a resource a statement writes.
type ChangedResource struct {
	// Database and Schema qualify the resource. An unqualified name is in the
	// current database and schema.
	Database string
	Schema   string
	// Table is the changed table, empty for a database-only change.
	Table string
	// DatabaseOnly marks a non-table object (view, routine, trigger, sequence,
	// synonym, type) whose own name carries a database qualifier. It has no
	// table to scope, only its target database.
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
	// InsertCount counts the rows of INSERT ... VALUES, and one for INSERT ...
	// DEFAULT VALUES.
	InsertCount int
}

// ExtractChangedResources returns the resources one statement writes and its
// DML counts. database and schema are the current ones, and text is the
// statement's text, which a DML statement is sampled from.
func ExtractChangedResources(node ast.Node, database, schema, text string) StatementChanges {
	var changes StatementChanges

	addDML := func() {
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, trimStatement(text))
	}

	addTable := func(ref *ast.TableRef, affectData bool) {
		d, s, t := tableRefTriple(ref, database, schema)
		if t == "" {
			return
		}
		changes.Resources = append(changes.Resources, ChangedResource{Database: d, Schema: s, Table: t, AffectsData: affectData})
	}

	// addObjectDatabase records a database-only write target for a non-table
	// object DDL ONLY when the object's own name carries an explicit database
	// qualifier. Unqualified names record nothing (request-database fallback).
	// The object's own qualifier is used — never an ON-table or body reference.
	addObjectDatabase := func(ref *ast.TableRef) {
		if ref != nil && ref.Database != "" {
			changes.Resources = append(changes.Resources, ChangedResource{Database: ref.Database, DatabaseOnly: true})
		}
	}

	switch n := node.(type) {
	case *ast.CreateTableStmt:
		addTable(n.Name, false)

	case *ast.AlterTableStmt:
		addTable(n.Name, true)

	case *ast.DropStmt:
		switch n.ObjectType {
		case ast.DropTable:
			if n.Names == nil {
				break
			}
			for _, item := range n.Names.Items {
				if ref, ok := item.(*ast.TableRef); ok {
					addTable(ref, true)
				}
			}
		case ast.DropIndex:
			if n.OnTables == nil {
				break
			}
			for _, item := range n.OnTables.Items {
				if ref, ok := item.(*ast.TableRef); ok {
					addTable(ref, false)
				}
			}
		case ast.DropView, ast.DropProcedure, ast.DropFunction, ast.DropTrigger,
			ast.DropSequence, ast.DropSynonym, ast.DropType:
			if n.Names == nil {
				break
			}
			for _, item := range n.Names.Items {
				if ref, ok := item.(*ast.TableRef); ok {
					addObjectDatabase(ref)
				}
			}
		default:
		}

	case *ast.CreateIndexStmt:
		addTable(n.Table, false)

	case *ast.TruncateStmt:
		addTable(n.Table, true)

	case *ast.InsertStmt:
		if ref, ok := n.Relation.(*ast.TableRef); ok {
			addTable(ref, false)
		}
		if n.DefaultValues {
			changes.InsertCount++
			break
		}
		if vc, ok := n.Source.(*ast.ValuesClause); ok {
			changes.InsertCount += vc.Rows.Len()
			break
		}
		// SHOWPLAN estimates one row for INSERT ... EXEC whatever the procedure returns, so the
		// statement is counted without an estimate.
		if _, ok := n.Source.(*ast.ExecStmt); ok {
			changes.DMLCount++
			break
		}
		addDML()

	case *ast.UpdateStmt:
		if ref, ok := n.Relation.(*ast.TableRef); ok {
			addTable(resolveDMLWriteTarget(ref, n.FromClause, n.WithClause), false)
		}
		addDML()

	case *ast.DeleteStmt:
		if ref, ok := n.Relation.(*ast.TableRef); ok {
			addTable(resolveDMLWriteTarget(ref, n.FromClause, n.WithClause), false)
		}
		addDML()

	case *ast.MergeStmt:
		if ref, ok := n.Target.(*ast.TableRef); ok {
			addTable(ref, false)
		}
		addDML()

	case *ast.CreateTableAsSelectStmt:
		addTable(n.Name, false)

	case *ast.SelectStmt:
		// SELECT ... INTO new_table creates a table; the INTO target is the write
		// target (possibly on the first arm of a set operation). INTO #temp is a
		// session-scoped tempdb table, not a database change.
		if target := selectIntoTarget(n); target != nil && !strings.HasPrefix(target.Object, "#") {
			addTable(target, false)
		}

	case *ast.CreateViewStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateProcedureStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateFunctionStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateTriggerStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateSequenceStmt:
		addObjectDatabase(n.Name)
	case *ast.AlterSequenceStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateSynonymStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateTypeStmt:
		addObjectDatabase(n.Name)

	default:
	}
	return changes
}

// resolveDMLWriteTarget returns the table an UPDATE/DELETE target writes. SQL Server binds an
// unqualified target to the FROM table with that exposed name: its alias, or its own name when
// it has none. A qualified target names its own table, even when a FROM entry is aliased to its
// bare name. Bytebase's synced metadata reports every SQL Server database as case-insensitive,
// even under a case-sensitive collation, so a target binds only when both collations agree: the
// FROM match is exact and the CTE check folds case. Other targets stay as written, such as a
// differently cased name or a CTE or derived-table target.
func resolveDMLWriteTarget(target *ast.TableRef, fromClause *ast.List, withClause *ast.WithClause) *ast.TableRef {
	if target.Server != "" || target.Database != "" || target.Schema != "" || fromClause == nil {
		return target
	}
	source := tableByExposedName(fromClause, target.Object)
	if source == nil || isCTETarget(source, withClause, false /* caseSensitive */) {
		return target
	}
	return source
}

func tableByExposedName(node ast.Node, name string) *ast.TableRef {
	switch n := node.(type) {
	case *ast.List:
		for _, item := range n.Items {
			if ref := tableByExposedName(item, name); ref != nil {
				return ref
			}
		}
	case *ast.TableRef:
		if n.Alias == name || (n.Alias == "" && n.Object == name) {
			return n
		}
	case *ast.AliasedTableRef:
		// The parser wraps an unaliased table that has TABLESAMPLE in an AliasedTableRef.
		if n.Alias == "" {
			return tableByExposedName(n.Table, name)
		}
		if ref, ok := n.Table.(*ast.TableRef); ok && n.Alias == name {
			return ref
		}
	case *ast.JoinClause:
		if ref := tableByExposedName(n.Left, name); ref != nil {
			return ref
		}
		return tableByExposedName(n.Right, name)
	default:
	}
	return nil
}

// isCTETarget reports whether an unqualified DML target names a CTE declared
// by the statement's own WITH clause. Name comparison follows the instance
// collation: a case-sensitive collation keeps c and C distinct.
func isCTETarget(ref *ast.TableRef, withClause *ast.WithClause, caseSensitive bool) bool {
	if ref == nil || withClause == nil || withClause.CTEs == nil || ref.Database != "" || ref.Schema != "" {
		return false
	}
	for _, cteNode := range withClause.CTEs.Items {
		cte, ok := cteNode.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		if cte.Name == ref.Object || (!caseSensitive && strings.EqualFold(cte.Name, ref.Object)) {
			return true
		}
	}
	return false
}

func tableRefTriple(ref *ast.TableRef, defaultDatabase, defaultSchema string) (string, string, string) {
	if ref == nil {
		return defaultDatabase, defaultSchema, ""
	}
	db := ref.Database
	if db == "" {
		db = defaultDatabase
	}
	schema := ref.Schema
	if schema == "" {
		schema = defaultSchema
	}
	return db, schema, ref.Object
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace) + ";"
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
