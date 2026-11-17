package analysis

import "github.com/bytebase/omni/googlesql/ast"

// StatementType is the kind of change a GoogleSQL statement makes. The values
// mirror Bytebase's storepb.StatementType names so the adapter maps them 1:1.
type StatementType int

const (
	StatementTypeUnspecified StatementType = iota
	Insert
	Update
	Delete
	Merge
	Truncate
	CreateTable
	CreateView
	CreateIndex
	CreateSchema
	CreateDatabase
	CreateFunction
	CreateProcedure
	CreateSequence
	AlterTable
	AlterView
	AlterIndex
	AlterDatabase
	AlterSequence
	DropTable
	DropView
	DropIndex
	DropSchema
	DropDatabase
	DropSequence
	DropFunction
	DropProcedure
	Rename
	Execute
)

// StatementTypeOf returns the type of a statement.
//
// CALL and EXECUTE IMMEDIATE report Execute: what they change is decided by
// the procedure body or the run-time SQL string, which the statement does not
// show.
//
// Deferral: every GoogleSQL statement the enum has no value for reports
// StatementTypeUnspecified, so a rule selecting DML skips it. The gap includes
// statements that do write rows — LOAD DATA INTO, CLONE DATA INTO, scripting
// bodies (BEGIN…END, IF, FOR, WHILE, each of which parses as one statement, so a
// CALL or EXECUTE IMMEDIATE nested in one is not seen either). It also includes
// ALTER SCHEMA, row access policies, generic entities, and the Spanner-only DDL
// set (change streams, roles, GRANT / REVOKE, proto bundles), which do not write
// rows.
//
// A new row-writing statement needs its own enum value before it lands here.
// ClassifySQL returns Unknown for these too, so it cannot gate the list.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	case *ast.MergeStmt:
		return Merge
	case *ast.TruncateStmt:
		return Truncate

	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.CreateMaterializedViewStmt:
		return CreateView
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.SearchVectorIndexStmt:
		return CreateIndex
	case *ast.CreateSchemaStmt:
		return CreateSchema
	case *ast.CreateDatabaseStmt:
		return CreateDatabase
	case *ast.CreateFunctionStmt:
		return CreateFunction
	case *ast.CreateProcedureStmt:
		return CreateProcedure
	case *ast.CreateSequenceStmt:
		return CreateSequence
	case *ast.CreateSnapshotStmt:
		return CreateTable

	case *ast.AlterStmt:
		return alterStatementType(n.Object)
	case *ast.DropStmt:
		return dropStatementType(n.Object)
	case *ast.AlterSequenceStmt:
		return AlterSequence
	case *ast.DropSequenceStmt:
		return DropSequence
	case *ast.BQAlterStmt:
		return bqAlterStatementType(n.Object)
	case *ast.BQDropStmt:
		return bqDropStatementType(n.Object)

	case *ast.RenameStmt:
		return Rename

	case *ast.CallStmt, *ast.ExecuteImmediateStmt:
		return Execute

	default:
		return StatementTypeUnspecified
	}
}

func alterStatementType(kind ast.AlterObjectKind) StatementType {
	switch kind {
	case ast.AlterTable:
		return AlterTable
	case ast.AlterView:
		return AlterView
	case ast.AlterIndex, ast.AlterSearchIndex:
		return AlterIndex
	case ast.AlterDatabase:
		return AlterDatabase
	default:
		// ALTER SCHEMA has no enum value; see the deferral on StatementTypeOf.
		return StatementTypeUnspecified
	}
}

func dropStatementType(kind ast.DropObjectKind) StatementType {
	switch kind {
	case ast.DropTable:
		return DropTable
	case ast.DropView:
		return DropView
	case ast.DropIndex:
		return DropIndex
	case ast.DropSchema:
		return DropSchema
	case ast.DropDatabase:
		return DropDatabase
	default:
		return StatementTypeUnspecified
	}
}

func bqAlterStatementType(kind ast.BQAlterObjectKind) StatementType {
	switch kind {
	case ast.BQAlterMaterializedView, ast.BQAlterApproxView:
		return AlterView
	case ast.BQAlterVectorIndex:
		return AlterIndex
	default:
		return StatementTypeUnspecified
	}
}

func bqDropStatementType(kind ast.BQDropObjectKind) StatementType {
	switch kind {
	case ast.BQDropFunction, ast.BQDropTableFunction:
		return DropFunction
	case ast.BQDropProcedure:
		return DropProcedure
	case ast.BQDropMaterializedView, ast.BQDropApproxView:
		return DropView
	case ast.BQDropSnapshotTable:
		return DropTable
	case ast.BQDropSearchIndex, ast.BQDropVectorIndex:
		return DropIndex
	default:
		return StatementTypeUnspecified
	}
}
