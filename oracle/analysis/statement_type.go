package analysis

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/oracle/ast"
)

// StatementType is the kind of change a statement makes.
type StatementType = analysis.StatementType

const (
	StatementTypeUnspecified = analysis.StatementTypeUnspecified
	CreateDatabase           = analysis.CreateDatabase
	CreateTable              = analysis.CreateTable
	CreateIndex              = analysis.CreateIndex
	CreateView               = analysis.CreateView
	CreateSequence           = analysis.CreateSequence
	CreateSchema             = analysis.CreateSchema
	CreateFunction           = analysis.CreateFunction
	CreateTrigger            = analysis.CreateTrigger
	CreateProcedure          = analysis.CreateProcedure
	CreateType               = analysis.CreateType
	DropDatabase             = analysis.DropDatabase
	DropTable                = analysis.DropTable
	DropView                 = analysis.DropView
	DropIndex                = analysis.DropIndex
	DropSequence             = analysis.DropSequence
	DropFunction             = analysis.DropFunction
	DropTrigger              = analysis.DropTrigger
	DropProcedure            = analysis.DropProcedure
	DropType                 = analysis.DropType
	AlterDatabase            = analysis.AlterDatabase
	AlterTable               = analysis.AlterTable
	AlterIndex               = analysis.AlterIndex
	AlterView                = analysis.AlterView
	AlterSequence            = analysis.AlterSequence
	AlterType                = analysis.AlterType
	Truncate                 = analysis.Truncate
	Rename                   = analysis.Rename
	Comment                  = analysis.Comment
	Insert                   = analysis.Insert
	Update                   = analysis.Update
	Delete                   = analysis.Delete
	Execute                  = analysis.Execute
)

// StatementTypeOf returns the type of a statement. A statement the switch does
// not name is StatementTypeUnspecified.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.CreateSequenceStmt:
		return CreateSequence
	case *ast.CreateSchemaStmt:
		return CreateSchema
	case *ast.CreateFunctionStmt:
		return CreateFunction
	case *ast.CreateTriggerStmt:
		return CreateTrigger
	case *ast.CreateProcedureStmt:
		return CreateProcedure
	case *ast.CreateTypeStmt:
		return CreateType
	case *ast.DropStmt:
		return dropStatementType(n.ObjectType)
	case *ast.AdminDDLStmt:
		return adminDDLStatementType(n.Action, n.ObjectType)
	case *ast.AlterTableStmt:
		return AlterTable
	case *ast.AlterIndexStmt:
		return AlterIndex
	case *ast.AlterViewStmt:
		return AlterView
	case *ast.AlterSequenceStmt:
		return AlterSequence
	case *ast.AlterTypeStmt:
		return AlterType
	case *ast.TruncateStmt:
		return Truncate
	case *ast.RenameStmt:
		return Rename
	case *ast.CommentStmt:
		return Comment
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	// An anonymous block is how Oracle runs a procedure (SQL*Plus EXEC expands to
	// one), and its body is not classified statement by statement.
	case *ast.CallStmt, *ast.PLSQLBlock:
		return Execute
	default:
		return StatementTypeUnspecified
	}
}

func adminDDLStatementType(action string, objectType ast.ObjectType) StatementType {
	if objectType != ast.OBJECT_DATABASE {
		return StatementTypeUnspecified
	}
	switch strings.ToUpper(action) {
	case "CREATE":
		return CreateDatabase
	case "ALTER":
		return AlterDatabase
	case "DROP":
		return DropDatabase
	default:
		return StatementTypeUnspecified
	}
}

func dropStatementType(objectType ast.ObjectType) StatementType {
	switch objectType {
	case ast.OBJECT_DATABASE:
		return DropDatabase
	case ast.OBJECT_TABLE:
		return DropTable
	case ast.OBJECT_VIEW, ast.OBJECT_MATERIALIZED_VIEW:
		return DropView
	case ast.OBJECT_INDEX:
		return DropIndex
	case ast.OBJECT_SEQUENCE:
		return DropSequence
	case ast.OBJECT_FUNCTION:
		return DropFunction
	case ast.OBJECT_TRIGGER:
		return DropTrigger
	case ast.OBJECT_PROCEDURE:
		return DropProcedure
	case ast.OBJECT_TYPE, ast.OBJECT_TYPE_BODY:
		return DropType
	default:
		return StatementTypeUnspecified
	}
}
