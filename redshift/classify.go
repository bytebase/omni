package redshift

import "github.com/bytebase/omni/redshift/ast"

// StatementType is a conservative top-level statement classification.
type StatementType string

const (
	StatementTypeUnknown        StatementType = "UNKNOWN"
	StatementTypeSelect         StatementType = "SELECT"
	StatementTypeDDL            StatementType = "DDL"
	StatementTypeDML            StatementType = "DML"
	StatementTypeShow           StatementType = "SHOW"
	StatementTypeCopy           StatementType = "COPY"
	StatementTypeUnload         StatementType = "UNLOAD"
	StatementTypeExplain        StatementType = "EXPLAIN"
	StatementTypeExplainAnalyze StatementType = "EXPLAIN_ANALYZE"
	StatementTypeSet            StatementType = "SET"
)

// GetStatementTypes parses SQL and returns one type per non-empty statement.
func GetStatementTypes(sql string) ([]StatementType, error) {
	stmts, err := Parse(sql)
	if err != nil {
		return nil, err
	}
	types := make([]StatementType, 0, len(stmts))
	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		types = append(types, ClassifyStatement(stmt.AST))
	}
	return types, nil
}

// ClassifyStatement classifies a parsed Redshift AST node.
func ClassifyStatement(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.SelectStmt:
		if n.IntoClause != nil {
			return StatementTypeDDL
		}
		return StatementTypeSelect
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		return StatementTypeDML
	case *ast.CopyStmt:
		return StatementTypeCopy
	case *ast.UnloadStmt:
		return StatementTypeUnload
	case *ast.RedshiftShowStmt, *ast.VariableShowStmt:
		return StatementTypeShow
	case *ast.ExplainStmt:
		if explainHasAnalyze(n) {
			return StatementTypeExplainAnalyze
		}
		return StatementTypeExplain
	case *ast.VariableSetStmt:
		return StatementTypeSet
	case *ast.RedshiftObjectStmt:
		switch n.Command {
		case "create", "alter", "drop", "attach", "detach":
			return StatementTypeDDL
		default:
			return StatementTypeUnknown
		}
	case *ast.CreateStmt, *ast.CreateTableAsStmt, *ast.ViewStmt, *ast.IndexStmt,
		*ast.DropStmt, *ast.AlterTableStmt, *ast.CreateSchemaStmt, *ast.TruncateStmt,
		*ast.CreateSeqStmt, *ast.AlterSeqStmt, *ast.CreateFunctionStmt, *ast.RenameStmt,
		*ast.AlterObjectSchemaStmt, *ast.AlterOwnerStmt, *ast.CreatedbStmt,
		*ast.AlterDatabaseStmt, *ast.DropdbStmt, *ast.CreateRoleStmt, *ast.AlterRoleStmt,
		*ast.DropRoleStmt:
		return StatementTypeDDL
	default:
		return StatementTypeUnknown
	}
}
