package analysis

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/clickhouse/parser"
)

// projectionName is the name ClickHouse gives a select-list item: its alias,
// a column's own name, or the expression written as the functions it calls
// (a + 1 is plus(a, 1)), which is how the analyzer names a projection.
func (r *spanResolver) projectionName(e parser.Expr, sc *spanScope, ctes *cteBinding) string {
	if a, ok := e.(*parser.AliasExpr); ok {
		return a.Alias
	}
	if parts, ok := identParts(e); ok {
		return parts[len(parts)-1]
	}
	return (&namer{r: r, sc: sc, ctes: ctes}).name(e)
}

type namer struct {
	r    *spanResolver
	sc   *spanScope
	ctes *cteBinding
	// subst names a lambda's parameter by the column APPLY applies it to.
	subst map[string]string
}

// apply is the name of a column a star's APPLY transformer applies to: the
// function's call on it, or the lambda's body with the column for its
// parameter.
func (n *namer) apply(t *parser.StarTransformer, column string) string {
	arg := &parser.IdentExpr{Parts: []string{column}}
	switch f := t.ApplyExpr.(type) {
	case *parser.LambdaExpr:
		if len(f.Params) == 1 {
			sub := &namer{r: n.r, sc: n.sc, ctes: n.ctes, subst: map[string]string{f.Params[0]: column}}
			return sub.name(f.Body)
		}
	case *parser.FuncCallExpr:
		// APPLY(f) calls f(column); APPLY(f(p)) calls the parametric f(p)(column).
		call := *f
		call.Params, call.Args = f.Args, []parser.Expr{arg}
		if len(f.Args) == 0 {
			call.Params = nil
		}
		return n.funcName(&call)
	}
	if t.ApplyFunc != "" {
		return t.ApplyFunc + "(" + column + ")"
	}
	return column
}

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
	return strings.Join(n.args(args), ", ")
}

// args names a call's arguments, a star or COLUMNS('re') as the columns it
// expands to.
func (n *namer) args(args []parser.Expr) []string {
	var parts []string
	for _, a := range args {
		if re, ok := columnsMatcher(a); ok && n.sc != nil {
			// COLUMNS('re') is an argument per column it matches.
			for _, c := range n.r.matching(re, n.sc) {
				parts = append(parts, c.name)
			}
			continue
		}
		if s, ok := a.(*parser.StarExpr); ok && n.sc != nil {
			// So is a star.
			for _, c := range n.r.star(s, n.sc, n.ctes) {
				parts = append(parts, c.name)
			}
			continue
		}
		parts = append(parts, n.name(a))
	}
	return parts
}

// maxNamed bounds the text of an expression named by the rules below: a
// longer one keeps its text, so a deep expression is not written out once
// per level.
const maxNamed = 1024

func (n *namer) name(e parser.Expr) string {
	if e != nil && e.EndPos()-e.StartPos() > maxNamed && e.EndPos() <= len(n.r.text) {
		return n.r.text[e.StartPos():e.EndPos()]
	}
	switch x := e.(type) {
	case nil:
		return "NULL"
	case *parser.AliasExpr:
		return x.Alias
	case *parser.IdentExpr:
		if len(x.Parts) == 1 {
			if v, ok := n.subst[x.Parts[0]]; ok {
				return v
			}
		}
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
		if strings.EqualFold(x.Name, "tuple") && len(x.Params) == 0 && x.Over == nil {
			// ClickHouse writes a tuple of more than one element as (a, b).
			if args := n.args(x.Args); len(args) > 1 {
				return "(" + strings.Join(args, ", ") + ")"
			}
		}
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
	if countStar(f) {
		args = nil
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
		// Decimal, or hexadecimal and binary by their prefix: ClickHouse
		// has no octal literals, so 010 is ten.
		digits, base := strings.ReplaceAll(l.Val, "_", ""), 10
		if len(digits) > 2 && digits[0] == '0' {
			switch digits[1] {
			case 'x', 'X':
				digits, base = digits[2:], 16
			case 'b', 'B':
				digits, base = digits[2:], 2
			}
		}
		if v, err := strconv.ParseUint(digits, base, 64); err == nil {
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
