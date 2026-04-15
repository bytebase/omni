// Package deparse converts Snowflake AST nodes back to SQL strings.
//
// Design choices:
//   - Keywords are always emitted in UPPER CASE.
//   - Identifier casing follows the original source (Ident.Quoted determines
//     whether double-quotes are added; unquoted identifiers keep their
//     original case).
//   - Output is single-line per statement (no cosmetic newlines). This
//     simplifies round-trip testing and downstream rewriting.
//   - If an unsupported node type is encountered an error is returned;
//     the deparser never panics.
package deparse

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// Deparse produces a Snowflake SQL string for the given AST node.
// Returns an error if the node has a type the deparser doesn't handle.
func Deparse(node ast.Node) (string, error) {
	w := &writer{}
	if err := w.writeNode(node); err != nil {
		return "", err
	}
	return w.String(), nil
}

// DeparseFile deparses an *ast.File (all statements, joined by ";\n").
func DeparseFile(file *ast.File) (string, error) {
	parts := make([]string, 0, len(file.Stmts))
	for _, stmt := range file.Stmts {
		out, err := Deparse(stmt)
		if err != nil {
			return "", err
		}
		parts = append(parts, out)
	}
	return strings.Join(parts, ";\n"), nil
}

// writeNode dispatches to the concrete writer for the given node.
func (w *writer) writeNode(node ast.Node) error {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	// --- Query statements ---
	case *ast.SelectStmt:
		return w.writeSelectStmt(n)
	case *ast.SetOperationStmt:
		return w.writeSetOperationStmt(n)

	// --- DML ---
	case *ast.InsertStmt:
		return w.writeInsertStmt(n)
	case *ast.InsertMultiStmt:
		return w.writeInsertMultiStmt(n)
	case *ast.UpdateStmt:
		return w.writeUpdateStmt(n)
	case *ast.DeleteStmt:
		return w.writeDeleteStmt(n)
	case *ast.MergeStmt:
		return w.writeMergeStmt(n)

	// --- DDL: table ---
	case *ast.CreateTableStmt:
		return w.writeCreateTableStmt(n)
	case *ast.AlterTableStmt:
		return w.writeAlterTableStmt(n)

	// --- DDL: database ---
	case *ast.CreateDatabaseStmt:
		return w.writeCreateDatabaseStmt(n)
	case *ast.AlterDatabaseStmt:
		return w.writeAlterDatabaseStmt(n)
	case *ast.DropDatabaseStmt:
		return w.writeDropDatabaseStmt(n)
	case *ast.UndropDatabaseStmt:
		return w.writeUndropDatabaseStmt(n)

	// --- DDL: schema ---
	case *ast.CreateSchemaStmt:
		return w.writeCreateSchemaStmt(n)
	case *ast.AlterSchemaStmt:
		return w.writeAlterSchemaStmt(n)
	case *ast.DropSchemaStmt:
		return w.writeDropSchemaStmt(n)
	case *ast.UndropSchemaStmt:
		return w.writeUndropSchemaStmt(n)

	// --- DDL: drop / undrop ---
	case *ast.DropStmt:
		return w.writeDropStmt(n)
	case *ast.UndropStmt:
		return w.writeUndropStmt(n)

	// --- DDL: view ---
	case *ast.CreateViewStmt:
		return w.writeCreateViewStmt(n)
	case *ast.AlterViewStmt:
		return w.writeAlterViewStmt(n)
	case *ast.CreateMaterializedViewStmt:
		return w.writeCreateMaterializedViewStmt(n)
	case *ast.AlterMaterializedViewStmt:
		return w.writeAlterMaterializedViewStmt(n)

	// --- Expressions (can also be top-level in some contexts) ---
	case *ast.Literal,
		*ast.ColumnRef,
		*ast.StarExpr,
		*ast.BinaryExpr,
		*ast.UnaryExpr,
		*ast.ParenExpr,
		*ast.CastExpr,
		*ast.CaseExpr,
		*ast.FuncCallExpr,
		*ast.IffExpr,
		*ast.CollateExpr,
		*ast.IsExpr,
		*ast.BetweenExpr,
		*ast.InExpr,
		*ast.LikeExpr,
		*ast.AccessExpr,
		*ast.ArrayLiteralExpr,
		*ast.JsonLiteralExpr,
		*ast.LambdaExpr,
		*ast.SubqueryExpr,
		*ast.ExistsExpr:
		return w.writeExpr(node)

	default:
		return fmt.Errorf("deparse: unsupported node type %T", node)
	}
}
