package analysis

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// StatementType is the kind of change a statement makes. The values mirror
// Bytebase's storepb.StatementType names.
type StatementType int

const (
	StatementTypeUnspecified StatementType = iota
	CreateDatabase
	CreateTable
	CreateIndex
	CreateView
	CreateSequence
	CreateSchema
	CreateFunction
	CreateTrigger
	CreateProcedure
	CreateType
	DropDatabase
	DropTable
	DropView
	DropIndex
	DropSequence
	DropFunction
	DropTrigger
	DropProcedure
	DropType
	AlterDatabase
	AlterTable
	AlterIndex
	AlterView
	AlterSequence
	AlterType
	Truncate
	Rename
	Comment
	Insert
	Update
	Delete
	Execute
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
