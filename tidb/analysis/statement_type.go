package analysis

import (
	"github.com/bytebase/omni/tidb/ast"
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
