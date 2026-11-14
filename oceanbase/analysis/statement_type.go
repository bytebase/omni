package analysis

import (
	"github.com/bytebase/omni/oceanbase/ast"
)

// StatementType is the kind of change a statement makes. The values mirror
// Bytebase's storepb.StatementType names.
type StatementType int

const (
	StatementTypeUnspecified StatementType = iota
	CreateDatabase
	CreateTable
	CreateView
	CreateIndex
	CreateEvent
	CreateTrigger
	CreateFunction
	CreateProcedure
	DropDatabase
	DropTable
	DropView
	DropIndex
	DropEvent
	DropTrigger
	DropFunction
	DropProcedure
	AlterDatabase
	AlterTable
	AlterView
	AlterEvent
	Truncate
	Rename
	Insert
	Update
	Delete
	Execute
)

// StatementTypeOf returns the type of a statement. A statement the switch does
// not name is StatementTypeUnspecified.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	// CREATE
	case *ast.CreateDatabaseStmt:
		return CreateDatabase
	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.CreateEventStmt:
		return CreateEvent
	case *ast.CreateTriggerStmt:
		return CreateTrigger
	case *ast.CreateFunctionStmt:
		if n.IsProcedure {
			return CreateProcedure
		}
		return CreateFunction

	// DROP
	case *ast.DropDatabaseStmt:
		return DropDatabase
	case *ast.DropTableStmt:
		return DropTable
	case *ast.DropIndexStmt:
		return DropIndex
	case *ast.DropViewStmt:
		return DropView
	case *ast.DropEventStmt:
		return DropEvent
	case *ast.DropTriggerStmt:
		return DropTrigger
	case *ast.DropRoutineStmt:
		if n.IsProcedure {
			return DropProcedure
		}
		return DropFunction

	// ALTER
	case *ast.AlterTableStmt:
		return AlterTable
	case *ast.AlterDatabaseStmt:
		return AlterDatabase
	case *ast.AlterViewStmt:
		return AlterView
	case *ast.AlterEventStmt:
		return AlterEvent

	// OTHER DDL
	case *ast.TruncateStmt:
		return Truncate
	case *ast.RenameTableStmt:
		return Rename

	// DML
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete

	// A procedure call or a prepared statement runs code the statement text does
	// not show. PREPARE only defines the statement, so it stays unclassified.
	case *ast.CallStmt, *ast.ExecuteStmt:
		return Execute

	default:
		return StatementTypeUnspecified
	}
}
