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
	case *ast.InExpr:
		return deparseInExpr(n)
	case *ast.BetweenExpr:
		return deparseBetweenExpr(n)
	case *ast.LikeExpr:
		return deparseLikeExpr(n)
	case *ast.IsExpr:
		return deparseIsExpr(n)
	case *ast.RowExpr:
		return deparseRowExpr(n)
	case *ast.CaseExpr:
		return deparseCaseExpr(n)
	case *ast.CastExpr:
		return deparseCastExpr(n)
	case *ast.ConvertExpr:
		return deparseConvertExpr(n)
	case *ast.IntervalExpr:
		return deparseIntervalExpr(n)
	case *ast.CollateExpr:
		return deparseCollateExpr(n)
	case *ast.FuncCallExpr:
		return deparseFuncCallExpr(n)
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
	left := n.Left
	right := n.Right
	// MySQL normalizes INTERVAL + expr to expr + INTERVAL (interval on the right)
	if _, ok := left.(*ast.IntervalExpr); ok {
		if _, ok2 := right.(*ast.IntervalExpr); !ok2 {
			left, right = right, left
		}
	}
	leftStr := deparseExpr(left)
	rightStr := deparseExpr(right)
	op := binaryOpToString(n.Op)
	return "(" + leftStr + " " + op + " " + rightStr + ")"
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

func deparseInExpr(n *ast.InExpr) string {
	expr := deparseExpr(n.Expr)
	keyword := "in"
	if n.Not {
		keyword = "not in"
	}
	// Build the value list with no spaces after commas
	items := make([]string, len(n.List))
	for i, item := range n.List {
		items[i] = deparseExpr(item)
	}
	return "(" + expr + " " + keyword + " (" + strings.Join(items, ",") + "))"
}

func deparseBetweenExpr(n *ast.BetweenExpr) string {
	expr := deparseExpr(n.Expr)
	low := deparseExpr(n.Low)
	high := deparseExpr(n.High)
	keyword := "between"
	if n.Not {
		keyword = "not between"
	}
	return "(" + expr + " " + keyword + " " + low + " and " + high + ")"
}

func deparseLikeExpr(n *ast.LikeExpr) string {
	expr := deparseExpr(n.Expr)
	pattern := deparseExpr(n.Pattern)
	keyword := "like"
	if n.Not {
		keyword = "not like"
	}
	result := "(" + expr + " " + keyword + " " + pattern
	if n.Escape != nil {
		result += " escape " + deparseExpr(n.Escape)
	}
	return result + ")"
}

func deparseIsExpr(n *ast.IsExpr) string {
	expr := deparseExpr(n.Expr)
	var test string
	switch n.Test {
	case ast.IsNull:
		if n.Not {
			test = "is not null"
		} else {
			test = "is null"
		}
	case ast.IsTrue:
		if n.Not {
			test = "is not true"
		} else {
			test = "is true"
		}
	case ast.IsFalse:
		if n.Not {
			test = "is not false"
		} else {
			test = "is false"
		}
	case ast.IsUnknown:
		if n.Not {
			test = "is not unknown"
		} else {
			test = "is unknown"
		}
	default:
		test = "is ?"
	}
	return "(" + expr + " " + test + ")"
}

func deparseRowExpr(n *ast.RowExpr) string {
	items := make([]string, len(n.Items))
	for i, item := range n.Items {
		items[i] = deparseExpr(item)
	}
	return "row(" + strings.Join(items, ",") + ")"
}

func deparseCaseExpr(n *ast.CaseExpr) string {
	var b strings.Builder
	b.WriteString("(case")
	if n.Operand != nil {
		b.WriteString(" ")
		b.WriteString(deparseExpr(n.Operand))
	}
	for _, w := range n.Whens {
		b.WriteString(" when ")
		b.WriteString(deparseExpr(w.Cond))
		b.WriteString(" then ")
		b.WriteString(deparseExpr(w.Result))
	}
	if n.Default != nil {
		b.WriteString(" else ")
		b.WriteString(deparseExpr(n.Default))
	}
	b.WriteString(" end)")
	return b.String()
}

func deparseCastExpr(n *ast.CastExpr) string {
	expr := deparseExpr(n.Expr)
	typeName := deparseDataType(n.TypeName)
	return "cast(" + expr + " as " + typeName + ")"
}

func deparseConvertExpr(n *ast.ConvertExpr) string {
	expr := deparseExpr(n.Expr)
	// CONVERT(expr USING charset) form
	if n.Charset != "" {
		return "convert(" + expr + " using " + strings.ToLower(n.Charset) + ")"
	}
	// CONVERT(expr, type) form — MySQL rewrites to CAST
	typeName := deparseDataType(n.TypeName)
	return "cast(" + expr + " as " + typeName + ")"
}

func deparseDataType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	name := strings.ToLower(dt.Name)
	switch name {
	case "char":
		result := "char"
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
		// MySQL adds charset for CHAR in CAST
		charset := dt.Charset
		if charset == "" {
			charset = "utf8mb4"
		}
		result += " charset " + strings.ToLower(charset)
		return result
	case "binary":
		// CAST to BINARY becomes cast(x as char charset binary)
		result := "char"
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
		result += " charset binary"
		return result
	case "signed", "signed integer":
		return "signed"
	case "unsigned", "unsigned integer":
		return "unsigned"
	case "decimal":
		if dt.Scale > 0 {
			return fmt.Sprintf("decimal(%d,%d)", dt.Length, dt.Scale)
		}
		if dt.Length > 0 {
			return fmt.Sprintf("decimal(%d)", dt.Length)
		}
		return "decimal"
	case "date":
		return "date"
	case "datetime":
		if dt.Length > 0 {
			return fmt.Sprintf("datetime(%d)", dt.Length)
		}
		return "datetime"
	case "time":
		if dt.Length > 0 {
			return fmt.Sprintf("time(%d)", dt.Length)
		}
		return "time"
	case "json":
		return "json"
	case "float":
		return "float"
	case "double":
		return "double"
	default:
		return name
	}
}

func deparseIntervalExpr(n *ast.IntervalExpr) string {
	val := deparseExpr(n.Value)
	return "interval " + val + " " + strings.ToLower(n.Unit)
}

func deparseCollateExpr(n *ast.CollateExpr) string {
	expr := deparseExpr(n.Expr)
	return "(" + expr + " collate " + n.Collation + ")"
}

// funcNameRewrites maps uppercase function names to their MySQL 8.0 canonical forms.
// These rewrites are applied by SHOW CREATE VIEW in MySQL 8.0.
var funcNameRewrites = map[string]string{
	"SUBSTRING":         "substr",
	"CURRENT_TIMESTAMP":  "now",
	"CURRENT_DATE":       "curdate",
	"CURRENT_TIME":       "curtime",
	"CURRENT_USER":       "current_user",
	"NOW":                "now",
	"LOCALTIME":          "now",
	"LOCALTIMESTAMP":     "now",
}

// deparseTrimDirectional handles TRIM(LEADING|TRAILING|BOTH remstr FROM str).
// MySQL 8.0 SHOW CREATE VIEW format: trim(leading 'x' from `a`)
func deparseTrimDirectional(direction string, args []ast.ExprNode) string {
	if len(args) == 2 {
		remstr := deparseExpr(args[0])
		str := deparseExpr(args[1])
		return "trim(" + direction + " " + remstr + " from " + str + ")"
	}
	// Fallback: single arg (shouldn't happen for directional, but be safe)
	if len(args) == 1 {
		return "trim(" + direction + " " + deparseExpr(args[0]) + ")"
	}
	return "trim()"
}

func deparseFuncCallExpr(n *ast.FuncCallExpr) string {
	// Handle TRIM special forms: TRIM_LEADING, TRIM_TRAILING, TRIM_BOTH
	// Parser encodes these as FuncCallExpr with Name="TRIM_LEADING" etc.
	// Args: [remstr, str] for directional forms
	name := strings.ToUpper(n.Name)
	switch name {
	case "TRIM_LEADING":
		return deparseTrimDirectional("leading", n.Args)
	case "TRIM_TRAILING":
		return deparseTrimDirectional("trailing", n.Args)
	case "TRIM_BOTH":
		return deparseTrimDirectional("both", n.Args)
	}

	// GROUP_CONCAT has special formatting
	if name == "GROUP_CONCAT" {
		return deparseGroupConcat(n)
	}

	// Determine the canonical function name
	canonical, ok := funcNameRewrites[name]
	if !ok {
		canonical = strings.ToLower(n.Name)
	}

	// Schema-qualified name
	if n.Schema != "" {
		canonical = strings.ToLower(n.Schema) + "." + canonical
	}

	// Zero-arg functions (CURRENT_TIMESTAMP, NOW(), etc.) — always emit parens
	if len(n.Args) == 0 && !n.Star {
		return canonical + "()"
	}

	// COUNT(*) — MySQL 8.0 rewrites COUNT(*) to count(0)
	if n.Star {
		return canonical + "(0)"
	}

	// Build argument list with no spaces after commas
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		args[i] = deparseExpr(arg)
	}

	var argStr string
	if n.Distinct {
		argStr = "distinct " + strings.Join(args, ",")
	} else {
		argStr = strings.Join(args, ",")
	}

	return canonical + "(" + argStr + ")"
}

// deparseGroupConcat handles GROUP_CONCAT with its special syntax:
// group_concat([distinct] expr [order by expr ASC|DESC] separator 'str')
// MySQL 8.0 always shows the separator (default ',') and explicit ASC in ORDER BY.
func deparseGroupConcat(n *ast.FuncCallExpr) string {
	var b strings.Builder
	b.WriteString("group_concat(")

	// DISTINCT
	if n.Distinct {
		b.WriteString("distinct ")
	}

	// Arguments
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(deparseExpr(arg))
	}

	// ORDER BY — MySQL 8.0 always shows explicit ASC/DESC
	if len(n.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range n.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(item.Expr))
			if item.Desc {
				b.WriteString(" DESC")
			} else {
				b.WriteString(" ASC")
			}
		}
	}

	// SEPARATOR — always shown; default is ','
	b.WriteString(" separator ")
	if n.Separator != nil {
		b.WriteString(deparseExpr(n.Separator))
	} else {
		b.WriteString("','")
	}

	b.WriteString(")")
	return b.String()
}

