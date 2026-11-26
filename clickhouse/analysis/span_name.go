package analysis

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/clickhouse/parser"
)

// projectionName is the name ClickHouse gives a select-list item: its alias,
// a column's own name, or the expression written as the functions it calls
// (a + 1 is plus(a, 1)), which is how the analyzer names a projection.
func (r *spanResolver) projectionName(e parser.Expr) string {
	if a, ok := e.(*parser.AliasExpr); ok {
		return a.Alias
	}
	if parts, ok := identParts(e); ok {
		return parts[len(parts)-1]
	}
	return (&namer{r: r}).name(e)
}

type namer struct{ r *spanResolver }

// binaryFunctions are the functions ClickHouse's operators call.
var binaryFunctions = map[string]string{
	"+": "plus", "-": "minus", "*": "multiply", "/": "divide", "%": "modulo",
	"||": "concat", "=": "equals", "==": "equals", "!=": "notEquals", "<>": "notEquals",
	"<": "less", ">": "greater", "<=": "lessOrEquals", ">=": "greaterOrEquals",
	"AND": "and", "OR": "or",
	"IS DISTINCT FROM": "isDistinctFrom", "IS NOT DISTINCT FROM": "isNotDistinctFrom",
}

// extractFunctions are the functions EXTRACT(part FROM x) calls.
var extractFunctions = map[string]string{
	"YEAR": "toYear", "QUARTER": "toQuarter", "MONTH": "toMonth", "WEEK": "toISOWeek",
	"DAY": "toDayOfMonth", "HOUR": "toHour", "MINUTE": "toMinute", "SECOND": "toSecond",
}

func (n *namer) call(fn string, args ...parser.Expr) string {
	return fn + "(" + n.list(args) + ")"
}

func (n *namer) list(args []parser.Expr) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = n.name(a)
	}
	return strings.Join(parts, ", ")
}

func (n *namer) name(e parser.Expr) string {
	switch x := e.(type) {
	case nil:
		return "NULL"
	case *parser.AliasExpr:
		return x.Alias
	case *parser.IdentExpr:
		return strings.Join(x.Parts, ".")
	case *parser.LiteralExpr:
		return literalName(x)
	case *parser.UnaryExpr:
		if x.Op == "NOT" {
			return n.call("not", x.Expr)
		}
		if lit, ok := x.Expr.(*parser.LiteralExpr); ok && (lit.Kind == "INT" || lit.Kind == "FLOAT") {
			return "-" + literalName(lit)
		}
		return n.call("negate", x.Expr)
	case *parser.BinaryExpr:
		fn, ok := binaryFunctions[x.Op]
		if !ok {
			break
		}
		if fn == "and" || fn == "or" {
			// ClickHouse folds a chain of AND (or OR) into one call.
			return fn + "(" + n.list(chain(x, x.Op)) + ")"
		}
		return n.call(fn, x.Left, x.Right)
	case *parser.TernaryExpr:
		return n.call("if", x.Cond, x.Then, x.Else)
	case *parser.ParenExpr:
		return n.name(x.Expr)
	case *parser.SubqueryExpr:
		// ClickHouse numbers the subqueries it has to name.
		n.r.subqueries++
		return "_subquery_" + strconv.Itoa(n.r.subqueries)
	case *parser.FuncCallExpr:
		return n.funcName(x)
	case *parser.CastExpr:
		if x.Type == nil {
			break
		}
		return "CAST(" + n.name(x.Expr) + ", " + quote(n.r.text[x.Type.Pos:x.Type.End]) + ")"
	case *parser.CaseExpr:
		var args []parser.Expr
		fn := "multiIf"
		if x.Operand != nil {
			fn = "caseWithExpression"
			args = append(args, x.Operand)
		}
		for _, w := range x.Whens {
			args = append(args, w.Cond, w.Then)
		}
		return n.call(fn, append(args, x.Else)...)
	case *parser.BetweenExpr:
		if x.Not {
			return "or(" + n.call("less", x.Expr, x.Low) + ", " + n.call("greater", x.Expr, x.High) + ")"
		}
		return "and(" + n.call("greaterOrEquals", x.Expr, x.Low) + ", " + n.call("lessOrEquals", x.Expr, x.High) + ")"
	case *parser.InExpr:
		fn := "in"
		if x.Not {
			fn = "notIn"
		}
		if x.Global {
			fn = "global" + strings.ToUpper(fn[:1]) + fn[1:]
		}
		return n.call(fn, x.Left, x.Right)
	case *parser.LikeExpr:
		fn := "like"
		if x.Ilike {
			fn = "ilike"
		}
		if x.Not {
			fn = "not" + strings.ToUpper(fn[:1]) + fn[1:]
		}
		return n.call(fn, x.Left, x.Right)
	case *parser.IsNullExpr:
		if x.Not {
			return n.call("isNotNull", x.Expr)
		}
		return n.call("isNull", x.Expr)
	case *parser.ArrayExpr:
		return "[" + n.list(x.Items) + "]"
	case *parser.TupleExpr:
		if len(x.Items) == 1 {
			return n.call("tuple", x.Items...)
		}
		return "(" + n.list(x.Items) + ")"
	case *parser.IndexExpr:
		return n.call("arrayElement", x.Expr, x.Index)
	case *parser.TupleAccessExpr:
		if parts, ok := identParts(x); ok {
			return strings.Join(parts, ".")
		}
		member := x.Member
		if _, err := strconv.Atoi(member); err != nil {
			member = quote(member)
		}
		return "tupleElement(" + n.name(x.Expr) + ", " + member + ")"
	case *parser.LambdaExpr:
		params := make([]string, len(x.Params))
		copy(params, x.Params)
		return "lambda(tuple(" + strings.Join(params, ", ") + "), " + n.name(x.Body) + ")"
	case *parser.IntervalExpr:
		unit := strings.ToLower(x.Unit)
		if unit == "" {
			break
		}
		return "toInterval" + strings.ToUpper(unit[:1]) + unit[1:] + "(" + n.name(x.Value) + ")"
	case *parser.ExtractExpr:
		if fn, ok := extractFunctions[strings.ToUpper(x.Part)]; ok {
			return n.call(fn, x.From)
		}
	case *parser.SubstringExpr:
		args := []parser.Expr{x.Expr, x.From}
		if x.For != nil {
			args = append(args, x.For)
		}
		return n.call("substring", args...)
	case *parser.TrimExpr:
		fn := map[string]string{"BOTH": "trimBoth", "LEADING": "trimLeft", "TRAILING": "trimRight"}[strings.ToUpper(x.Mode)]
		if fn != "" && x.Chars != nil {
			return n.call(fn, x.From, x.Chars)
		}
	}
	// A form named by no rule above keeps its text.
	return n.r.text[e.StartPos():e.EndPos()]
}

// chain flattens a left-deep chain of op into its operands.
func chain(e parser.Expr, op string) []parser.Expr {
	if b, ok := e.(*parser.BinaryExpr); ok && b.Op == op {
		return append(chain(b.Left, op), chain(b.Right, op)...)
	}
	return []parser.Expr{e}
}

func (n *namer) funcName(f *parser.FuncCallExpr) string {
	name := f.Name
	args := f.Args
	if len(args) == 1 {
		if s, ok := args[0].(*parser.StarExpr); ok && len(s.Qualifier) == 0 && strings.EqualFold(name, "count") {
			args = nil
		}
	}
	if f.Distinct {
		name += "Distinct"
	}
	out := name
	if len(f.Params) > 0 {
		out += "(" + n.list(f.Params) + ")"
	}
	out += "(" + n.list(args) + ")"
	switch {
	case f.Over != nil:
		out += " OVER (" + n.window(f.Over) + ")"
	case f.OverName != "":
		out += " OVER " + f.OverName
	}
	return out
}

func (n *namer) window(w *parser.WindowSpec) string {
	var parts []string
	if w.BaseWindow != "" {
		parts = append(parts, w.BaseWindow)
	}
	if len(w.PartitionBy) > 0 {
		parts = append(parts, "PARTITION BY "+n.list(w.PartitionBy))
	}
	if len(w.OrderBy) > 0 {
		items := make([]string, len(w.OrderBy))
		for i, o := range w.OrderBy {
			dir := " ASC"
			if o.Desc {
				dir = " DESC"
			}
			items[i] = n.name(o.Expr) + dir
		}
		parts = append(parts, "ORDER BY "+strings.Join(items, ", "))
	}
	if w.FrameUnits != "" {
		frame := strings.ToUpper(w.FrameUnits) + " "
		if w.FrameEnd != "" {
			frame += "BETWEEN " + w.FrameStart + " AND " + w.FrameEnd
		} else {
			frame += w.FrameStart
		}
		parts = append(parts, frame)
	}
	return strings.Join(parts, " ")
}

func literalName(l *parser.LiteralExpr) string {
	switch l.Kind {
	case "STRING":
		return quote(l.Val)
	case "INT":
		if v, err := strconv.ParseUint(strings.ReplaceAll(l.Val, "_", ""), 0, 64); err == nil {
			return strconv.FormatUint(v, 10)
		}
	case "FLOAT":
		if v, err := strconv.ParseFloat(l.Val, 64); err == nil {
			s := strconv.FormatFloat(v, 'f', -1, 64)
			if !strings.ContainsAny(s, ".e") {
				s += "."
			}
			return s
		}
	case "NULL":
		return "NULL"
	case "BOOL":
		return strings.ToLower(l.Val)
	case "INF":
		return "inf"
	case "NAN":
		return "nan"
	}
	return l.Raw
}

// quote writes s as a ClickHouse string literal.
func quote(s string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s) + "'"
}
