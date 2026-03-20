// Package deparse converts MySQL AST nodes back to SQL text,
// matching MySQL 8.0's SHOW CREATE VIEW formatting.
package deparse

import (
	"fmt"
	"math/big"
	"strings"

	ast "github.com/bytebase/omni/mysql/ast"
)

// Deparse converts an expression AST node to its SQL text representation,
// matching MySQL 8.0's canonical formatting (as seen in SHOW CREATE VIEW).
func Deparse(node ast.ExprNode) string {
	if node == nil {
		return ""
	}
	return deparseExpr(node)
}

func deparseExpr(node ast.ExprNode) string {
	switch n := node.(type) {
	case *ast.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *ast.FloatLit:
		return n.Value
	case *ast.BoolLit:
		if n.Value {
			return "true"
		}
		return "false"
	case *ast.StringLit:
		return deparseStringLit(n)
	case *ast.NullLit:
		return "NULL"
	case *ast.HexLit:
		return deparseHexLit(n)
	case *ast.BitLit:
		return deparseBitLit(n)
	case *ast.BinaryExpr:
		return deparseBinaryExpr(n)
	case *ast.ColumnRef:
		return deparseColumnRef(n)
	case *ast.UnaryExpr:
		return deparseUnaryExpr(n)
	case *ast.ParenExpr:
		return deparseExpr(n.Expr)
	default:
		return fmt.Sprintf("/* unsupported: %T */", node)
	}
}

func deparseStringLit(n *ast.StringLit) string {
	// MySQL 8.0 uses backslash escaping for single quotes: '' → \'
	// and preserves backslashes as-is.
	escaped := strings.ReplaceAll(n.Value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
	if n.Charset != "" {
		return n.Charset + "'" + escaped + "'"
	}
	return "'" + escaped + "'"
}

func deparseHexLit(n *ast.HexLit) string {
	// MySQL 8.0 normalizes all hex literals to 0x lowercase form.
	// HexLit.Value is either "0xFF" (0x prefix form) or "FF" (X'' form).
	val := n.Value
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		// Already has 0x prefix — just lowercase
		return "0x" + strings.ToLower(val[2:])
	}
	// X'FF' form — value is just the hex digits
	return "0x" + strings.ToLower(val)
}

func deparseBitLit(n *ast.BitLit) string {
	// MySQL 8.0 converts all bit literals to hex form.
	// BitLit.Value is either "0b1010" (0b prefix form) or "1010" (b'' form).
	val := n.Value
	if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
		val = val[2:]
	}
	// Parse binary string to integer, then format as hex
	i := new(big.Int)
	i.SetString(val, 2)
	return "0x" + fmt.Sprintf("%02x", i)
}

func deparseBinaryExpr(n *ast.BinaryExpr) string {
	left := deparseExpr(n.Left)
	right := deparseExpr(n.Right)
	op := binaryOpToString(n.Op)
	return "(" + left + " " + op + " " + right + ")"
}

func deparseColumnRef(n *ast.ColumnRef) string {
	if n.Schema != "" {
		return "`" + n.Schema + "`.`" + n.Table + "`.`" + n.Column + "`"
	}
	if n.Table != "" {
		return "`" + n.Table + "`.`" + n.Column + "`"
	}
	return "`" + n.Column + "`"
}

func binaryOpToString(op ast.BinaryOp) string {
	switch op {
	case ast.BinOpAdd:
		return "+"
	case ast.BinOpSub:
		return "-"
	case ast.BinOpMul:
		return "*"
	case ast.BinOpDiv:
		return "/"
	case ast.BinOpMod:
		return "%"
	case ast.BinOpDivInt:
		return "DIV"
	case ast.BinOpEq:
		return "="
	case ast.BinOpNe:
		return "<>"
	case ast.BinOpLt:
		return "<"
	case ast.BinOpGt:
		return ">"
	case ast.BinOpLe:
		return "<="
	case ast.BinOpGe:
		return ">="
	case ast.BinOpNullSafeEq:
		return "<=>"
	case ast.BinOpAnd:
		return "and"
	case ast.BinOpOr:
		return "or"
	case ast.BinOpBitAnd:
		return "&"
	case ast.BinOpBitOr:
		return "|"
	case ast.BinOpBitXor:
		return "^"
	case ast.BinOpShiftLeft:
		return "<<"
	case ast.BinOpShiftRight:
		return ">>"
	case ast.BinOpSoundsLike:
		return "sounds like"
	default:
		return "?"
	}
}

func deparseUnaryExpr(n *ast.UnaryExpr) string {
	operand := deparseExpr(n.Operand)
	switch n.Op {
	case ast.UnaryMinus:
		return "-" + operand
	case ast.UnaryPlus:
		// MySQL drops unary plus entirely
		return operand
	case ast.UnaryBitNot:
		return "~" + operand
	default:
		return operand
	}
}

