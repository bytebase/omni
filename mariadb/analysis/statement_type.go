package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mariadb/ast"
)

// StatementType is the kind of change a statement makes.
type StatementType = analysis.StatementType

const (
	StatementTypeUnspecified = analysis.StatementTypeUnspecified
	CreateDatabase           = analysis.CreateDatabase
	CreateTable              = analysis.CreateTable
	CreateView               = analysis.CreateView
	CreateIndex              = analysis.CreateIndex
	CreateEvent              = analysis.CreateEvent
	CreateTrigger            = analysis.CreateTrigger
	CreateFunction           = analysis.CreateFunction
	CreateProcedure          = analysis.CreateProcedure
	DropDatabase             = analysis.DropDatabase
	DropTable                = analysis.DropTable
	DropView                 = analysis.DropView
	DropIndex                = analysis.DropIndex
	DropEvent                = analysis.DropEvent
	DropTrigger              = analysis.DropTrigger
	DropFunction             = analysis.DropFunction
	DropProcedure            = analysis.DropProcedure
	AlterDatabase            = analysis.AlterDatabase
	AlterTable               = analysis.AlterTable
	AlterView                = analysis.AlterView
	AlterEvent               = analysis.AlterEvent
	Truncate                 = analysis.Truncate
	Rename                   = analysis.Rename
	Insert                   = analysis.Insert
	Update                   = analysis.Update
	Delete                   = analysis.Delete
	Execute                  = analysis.Execute
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
