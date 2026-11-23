package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/tidb/ast"
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

// StatementTypeOf returns the type of a statement.
func StatementTypeOf(stmt ast.Node) StatementType {
	switch stmt.(type) {
	// DDL statements
	case *ast.CreateDatabaseStmt:
		return CreateDatabase
	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.CreateIndexStmt:
		return CreateIndex

	case *ast.DropDatabaseStmt:
		return DropDatabase
	case *ast.DropTableStmt:
		return DropTable
	case *ast.DropViewStmt:
		return DropView
	case *ast.DropIndexStmt:
		return DropIndex

	case *ast.AlterDatabaseStmt:
		return AlterDatabase
	case *ast.AlterTableStmt:
		return AlterTable

	case *ast.TruncateStmt:
		return Truncate
	case *ast.RenameTableStmt:
		return Rename

	// DML statements
	case *ast.DeleteStmt:
		return Delete
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update

	// A procedure call or a prepared statement runs code the statement text does
	// not show. PREPARE only defines the statement, so it stays unclassified.
	case *ast.CallStmt, *ast.ExecuteStmt:
		return Execute
	}

	return StatementTypeUnspecified
}
