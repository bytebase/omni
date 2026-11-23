package analysis

import (
	"slices"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/pg/ast"
)

// StatementType is the kind of change a statement makes.
type StatementType = analysis.StatementType

// StatementTypeUnspecified is a statement the classifier does not name.
const StatementTypeUnspecified = analysis.StatementTypeUnspecified

const (
	AlterSequence   = analysis.AlterSequence
	AlterTable      = analysis.AlterTable
	AlterType       = analysis.AlterType
	AlterView       = analysis.AlterView
	Comment         = analysis.Comment
	CreateDatabase  = analysis.CreateDatabase
	CreateExtension = analysis.CreateExtension
	CreateFunction  = analysis.CreateFunction
	CreateIndex     = analysis.CreateIndex
	CreateSchema    = analysis.CreateSchema
	CreateSequence  = analysis.CreateSequence
	CreateTable     = analysis.CreateTable
	CreateTrigger   = analysis.CreateTrigger
	CreateType      = analysis.CreateType
	CreateView      = analysis.CreateView
	Delete          = analysis.Delete
	DropDatabase    = analysis.DropDatabase
	DropExtension   = analysis.DropExtension
	DropFunction    = analysis.DropFunction
	DropIndex       = analysis.DropIndex
	DropSchema      = analysis.DropSchema
	DropSequence    = analysis.DropSequence
	DropTable       = analysis.DropTable
	DropTrigger     = analysis.DropTrigger
	DropType        = analysis.DropType
	DropView        = analysis.DropView
	Execute         = analysis.Execute
	Insert          = analysis.Insert
	Merge           = analysis.Merge
	RenameIndex     = analysis.RenameIndex
	RenameSchema    = analysis.RenameSchema
	RenameSequence  = analysis.RenameSequence
	Truncate        = analysis.Truncate
	Update          = analysis.Update
)

// StatementTypes returns the type of the statement followed by the types of its
// data-modifying CTEs, each type once. EXPLAIN ANALYZE takes the types of the
// statement it runs. A statement the classifier does not name adds none.
func StatementTypes(node ast.Node) []StatementType {
	node, _ = UnwrapExplainAnalyze(node, "")
	var types []StatementType
	add := func(statementType StatementType) {
		if statementType != StatementTypeUnspecified && !slices.Contains(types, statementType) {
			types = append(types, statementType)
		}
	}
	add(classifyStatementType(node))
	if with := getWithClause(node); with != nil && with.Ctes != nil {
		for _, item := range with.Ctes.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				add(classifyStatementType(cte.Ctequery))
			}
		}
	}
	return types
}

// classifyStatementType returns the statement type for an omni AST node.
func classifyStatementType(node ast.Node) StatementType {
	switch n := node.(type) {
	// DDL - CREATE
	case *ast.CreateStmt:
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
	case *ast.CreateTrigStmt:
		return CreateTrigger
	case *ast.CreateExtensionStmt:
		return CreateExtension
	case *ast.CreatedbStmt:
		return CreateDatabase
	case *ast.DefineStmt:
		if n.Kind == ast.OBJECT_TYPE {
			return CreateType
		}
		return StatementTypeUnspecified
	case *ast.CreateEnumStmt:
		return CreateType
	case *ast.CreateTableAsStmt:
		if n.Objtype == ast.OBJECT_MATVIEW {
			return CreateView
		}
		return CreateTable
	case *ast.SelectStmt:
		if hasIntoClause(n) {
			return CreateTable
		}
		return StatementTypeUnspecified

	// DDL - DROP
	case *ast.DropStmt:
		return getDropStatementTypeFrom(n)
	case *ast.DropdbStmt:
		return DropDatabase
	case *ast.TruncateStmt:
		return Truncate

	// DDL - ALTER
	case *ast.AlterTableStmt:
		if ast.ObjectType(n.ObjType) == ast.OBJECT_VIEW {
			return AlterView
		}
		return AlterTable
	case *ast.AlterSeqStmt:
		return AlterSequence
	case *ast.AlterEnumStmt:
		return AlterType

	// DDL - RENAME
	case *ast.RenameStmt:
		return getRenameStatementType(n)

	// DDL - COMMENT
	case *ast.CommentStmt:
		return Comment

	// DDL - DROP FUNCTION (AlterFunctionStmt with drop action, or separate node)
	case *ast.AlterFunctionStmt:
		// This is ALTER FUNCTION, not DROP FUNCTION. Not classified in original code.
		return StatementTypeUnspecified

	// DML
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	case *ast.MergeStmt:
		return Merge

	// A procedure, DO block or prepared statement runs code this classifier does not look into.
	case *ast.CallStmt, *ast.DoStmt, *ast.ExecuteStmt:
		return Execute

	default:
		return StatementTypeUnspecified
	}
}

// getDropStatementTypeFrom determines the specific DROP statement type from omni AST.
func getDropStatementTypeFrom(n *ast.DropStmt) StatementType {
	objType := ast.ObjectType(n.RemoveType)
	switch objType {
	case ast.OBJECT_TABLE:
		return DropTable
	case ast.OBJECT_VIEW:
		return DropView
	case ast.OBJECT_MATVIEW:
		// DROP MATERIALIZED VIEW holds data — treat as DROP_TABLE for risk assessment.
		return DropTable
	case ast.OBJECT_INDEX:
		return DropIndex
	case ast.OBJECT_SEQUENCE:
		return DropSequence
	case ast.OBJECT_SCHEMA:
		return DropSchema
	case ast.OBJECT_EXTENSION:
		return DropExtension
	case ast.OBJECT_TYPE:
		return DropType
	case ast.OBJECT_TRIGGER:
		return DropTrigger
	case ast.OBJECT_FUNCTION:
		return DropFunction
	default:
		return StatementTypeUnspecified
	}
}

// getRenameStatementType determines the statement type for a RenameStmt.
func getRenameStatementType(n *ast.RenameStmt) StatementType {
	switch n.RenameType {
	case ast.OBJECT_INDEX:
		return RenameIndex
	case ast.OBJECT_SCHEMA:
		return RenameSchema
	case ast.OBJECT_SEQUENCE:
		return RenameSequence
	case ast.OBJECT_VIEW:
		return AlterView
	case ast.OBJECT_TABLE, ast.OBJECT_COLUMN, ast.OBJECT_TABCONSTRAINT:
		// RENAME TABLE, RENAME COLUMN, RENAME CONSTRAINT all return ALTER_TABLE
		// Check RelationType for VIEW
		if n.RelationType == ast.OBJECT_VIEW {
			return AlterView
		}
		return AlterTable
	default:
		return StatementTypeUnspecified
	}
}
