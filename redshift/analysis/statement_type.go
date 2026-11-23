package analysis

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/redshift/ast"
)

// StatementType is the kind of change a statement makes.
type StatementType = analysis.StatementType

const (
	StatementTypeUnspecified = analysis.StatementTypeUnspecified
	CreateDatabase           = analysis.CreateDatabase
	CreateTable              = analysis.CreateTable
	CreateView               = analysis.CreateView
	CreateIndex              = analysis.CreateIndex
	CreateSequence           = analysis.CreateSequence
	CreateSchema             = analysis.CreateSchema
	CreateFunction           = analysis.CreateFunction
	DropDatabase             = analysis.DropDatabase
	DropTable                = analysis.DropTable
	DropView                 = analysis.DropView
	DropSchema               = analysis.DropSchema
	DropIndex                = analysis.DropIndex
	DropSequence             = analysis.DropSequence
	DropFunction             = analysis.DropFunction
	AlterDatabase            = analysis.AlterDatabase
	AlterTable               = analysis.AlterTable
	AlterView                = analysis.AlterView
	AlterIndex               = analysis.AlterIndex
	AlterSequence            = analysis.AlterSequence
	RenameIndex              = analysis.RenameIndex
	RenameSchema             = analysis.RenameSchema
	RenameSequence           = analysis.RenameSequence
	Comment                  = analysis.Comment
	Truncate                 = analysis.Truncate
	Insert                   = analysis.Insert
	Update                   = analysis.Update
	Delete                   = analysis.Delete
	Merge                    = analysis.Merge
	Execute                  = analysis.Execute
)

// StatementTypeOf returns the type of a statement. A statement this switch does
// not name is StatementTypeUnspecified.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.CreateStmt:
		return CreateTable
	case *ast.CreateTableAsStmt:
		if ast.ObjectType(n.Objtype) == ast.OBJECT_MATVIEW {
			return CreateView
		}
		return CreateTable
	case *ast.ViewStmt:
		return CreateView
	case *ast.IndexStmt:
		return CreateIndex
	case *ast.CreateSeqStmt:
		return CreateSequence
	case *ast.CreateSchemaStmt:
		return CreateSchema
	case *ast.CreateFunctionStmt:
		return CreateFunction
	case *ast.CreatedbStmt:
		return CreateDatabase
	case *ast.DropdbStmt:
		return DropDatabase
	case *ast.DropStmt:
		return dropStatementType(ast.ObjectType(n.RemoveType))
	case *ast.AlterTableStmt:
		if ast.ObjectType(n.ObjType) == ast.OBJECT_VIEW || ast.ObjectType(n.ObjType) == ast.OBJECT_MATVIEW {
			return AlterView
		}
		return AlterTable
	case *ast.AlterSeqStmt:
		return AlterSequence
	case *ast.RenameStmt:
		return renameStatementType(n.RenameType)
	case *ast.CommentStmt:
		return Comment
	case *ast.TruncateStmt:
		return Truncate
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	case *ast.MergeStmt:
		return Merge
	// A procedure or prepared statement runs code this classifier does not look
	// into. Redshift has no DO.
	case *ast.CallStmt, *ast.ExecuteStmt:
		return Execute
	case *ast.RedshiftObjectStmt:
		return redshiftObjectStatementType(n.Command, n.ObjectType)
	default:
		return StatementTypeUnspecified
	}
}

func dropStatementType(objectType ast.ObjectType) StatementType {
	switch objectType {
	case ast.OBJECT_TABLE:
		return DropTable
	case ast.OBJECT_VIEW:
		return DropView
	case ast.OBJECT_MATVIEW:
		return DropTable
	case ast.OBJECT_SCHEMA:
		return DropSchema
	case ast.OBJECT_INDEX:
		return DropIndex
	case ast.OBJECT_SEQUENCE:
		return DropSequence
	case ast.OBJECT_FUNCTION:
		return DropFunction
	default:
		return StatementTypeUnspecified
	}
}

func renameStatementType(objectType ast.ObjectType) StatementType {
	switch objectType {
	case ast.OBJECT_INDEX:
		return RenameIndex
	case ast.OBJECT_SCHEMA:
		return RenameSchema
	case ast.OBJECT_SEQUENCE:
		return RenameSequence
	case ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
		return AlterView
	case ast.OBJECT_TABLE:
		return AlterTable
	default:
		return StatementTypeUnspecified
	}
}

func redshiftObjectStatementType(command, objectType string) StatementType {
	switch strings.ToLower(command) {
	case "create":
		switch strings.ToLower(objectType) {
		case "database":
			return CreateDatabase
		case "schema", "external schema":
			return CreateSchema
		default:
			return StatementTypeUnspecified
		}
	case "alter":
		switch strings.ToLower(objectType) {
		case "database":
			return AlterDatabase
		case "index":
			return AlterIndex
		case "sequence":
			return AlterSequence
		case "table":
			return AlterTable
		case "view", "external view", "materialized view":
			return AlterView
		default:
			return StatementTypeUnspecified
		}
	case "drop":
		switch strings.ToLower(objectType) {
		case "database":
			return DropDatabase
		case "schema", "external schema":
			return DropSchema
		default:
			return StatementTypeUnspecified
		}
	default:
		return StatementTypeUnspecified
	}
}
