package analysis

import (
	"github.com/bytebase/omni/mssql/ast"
)

// StatementType is the kind of change a statement makes. The values mirror
// Bytebase's storepb.StatementType names.
type StatementType int

const (
	StatementTypeUnspecified StatementType = iota
	AlterDatabase
	AlterIndex
	AlterTable
	CreateDatabase
	CreateIndex
	CreateSchema
	CreateTable
	CreateView
	DropDatabase
	DropIndex
	DropSchema
	DropTable
	DropView
	Truncate
	Delete
	Insert
	Update
	Execute
)

// StatementTypeOf returns the type of a statement. A statement the switch does
// not name is StatementTypeUnspecified.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.AlterDatabaseStmt:
		return AlterDatabase
	case *ast.AlterIndexStmt:
		return AlterIndex
	case *ast.AlterTableStmt:
		return AlterTable
	case *ast.CreateDatabaseStmt:
		return CreateDatabase
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.CreateSchemaStmt:
		return CreateSchema
	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.DropStmt:
		switch n.ObjectType {
		case ast.DropDatabase:
			return DropDatabase
		case ast.DropIndex:
			return DropIndex
		case ast.DropSchema:
			return DropSchema
		case ast.DropTable:
			return DropTable
		case ast.DropView:
			return DropView
		default:
		}
	case *ast.TruncateStmt:
		return Truncate
	case *ast.DeleteStmt:
		return Delete
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.ExecStmt:
		return Execute
	default:
	}
	return StatementTypeUnspecified
}
