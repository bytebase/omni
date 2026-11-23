package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mssql/ast"
)

// StatementType is the kind of change a statement makes.
type StatementType = analysis.StatementType

const (
	StatementTypeUnspecified = analysis.StatementTypeUnspecified
	AlterDatabase            = analysis.AlterDatabase
	AlterIndex               = analysis.AlterIndex
	AlterTable               = analysis.AlterTable
	CreateDatabase           = analysis.CreateDatabase
	CreateIndex              = analysis.CreateIndex
	CreateSchema             = analysis.CreateSchema
	CreateTable              = analysis.CreateTable
	CreateView               = analysis.CreateView
	DropDatabase             = analysis.DropDatabase
	DropIndex                = analysis.DropIndex
	DropSchema               = analysis.DropSchema
	DropTable                = analysis.DropTable
	DropView                 = analysis.DropView
	Truncate                 = analysis.Truncate
	Delete                   = analysis.Delete
	Insert                   = analysis.Insert
	Update                   = analysis.Update
	Execute                  = analysis.Execute
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
