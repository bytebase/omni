// Package deparse — rewrite.go implements AST rewrites applied before deparsing.
// These rewrites match MariaDB's resolver behavior, which is what
// information_schema.VIEWS holds for a view the server accepted.
package deparse

import (
	"strings"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// RewriteExpr applies MariaDB resolver rewrites to an expression AST.
// Currently implements NOT folding: NOT(comparison) → inverted comparison.
// The rewrite is recursive — children are processed first (bottom-up).
func RewriteExpr(node ast.ExprNode) ast.ExprNode {
	if node == nil {
		return nil
	}
	return rewriteExpr(node)
}

func rewriteExpr(node ast.ExprNode) ast.ExprNode {
	switch n := node.(type) {
	case *ast.UnaryExpr:
		// First rewrite the operand recursively
		n.Operand = rewriteExpr(n.Operand)
		if n.Op == ast.UnaryNot {
			return rewriteNot(n.Operand)
		}
		if n.Op == ast.UnaryBinary {
			// The server keeps no BINARY operator: it stores the cast the
			// operator stands for, which is the text CAST(x AS BINARY) gives.
			// Measured on mariadb:11.8.8, `2 * BINARY (1 + 2)` comes back
			// 2 * cast(1 + 2 as char charset binary) and `BINARY 'a' = 'A'`
			// cast('a' as char charset binary) = 'A'.
			return &ast.CastExpr{
				Loc:      n.Loc,
				Expr:     n.Operand,
				TypeName: &ast.DataType{Name: "binary"},
			}
		}
		return n

	case *ast.BinaryExpr:
		n.Left = rewriteExpr(n.Left)
		n.Right = rewriteExpr(n.Right)
		// SOUNDS LIKE is stored as soundex(left) = soundex(right).
		if n.Op == ast.BinOpSoundsLike {
			return &ast.ParenExpr{
				Expr: &ast.BinaryExpr{
					Op: ast.BinOpEq,
					Left: &ast.FuncCallExpr{
						Name:      "soundex",
						HasParens: true,
						Args:      []ast.ExprNode{n.Left},
					},
					Right: &ast.FuncCallExpr{
						Name:      "soundex",
						HasParens: true,
						Args:      []ast.ExprNode{n.Right},
					},
				},
			}
		}
		// Boolean context wrapping: an AND or OR operand that is a bare column
		// is stored as `col <> 0`. Measured on mariadb:11.8.8: `a AND b` over
		// integer columns is `a <> 0 and b <> 0`, while `a AND (b + 1)` keeps
		// `b + 1`, `ifnull(a,0) AND b` keeps the call, `1 AND a` keeps the
		// literal, and XOR wraps neither side. MySQL 8.0 instead wraps every
		// operand that is not itself a boolean expression.
		if n.Op == ast.BinOpAnd || n.Op == ast.BinOpOr {
			if needsBooleanContext(n.Left) {
				n.Left = wrapBooleanContext(n.Left)
			}
			if needsBooleanContext(n.Right) {
				n.Right = wrapBooleanContext(n.Right)
			}
		}
		return n

	case *ast.ParenExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.InExpr:
		n.Expr = rewriteExpr(n.Expr)
		for i, item := range n.List {
			n.List[i] = rewriteExpr(item)
		}
		RewriteSelectStmt(n.Select)
		return n

	case *ast.SubqueryExpr:
		// A subquery's body is rewritten like any other: measured on
		// mariadb:11.8.8, `(SELECT (a = 1) SOUNDS LIKE '0' FROM t LIMIT 1)`
		// is stored as soundex(`t`.`a` = 1) = soundex('0'), the same as in a
		// target of its own.
		RewriteSelectStmt(n.Select)
		return n

	case *ast.ExistsExpr:
		RewriteSelectStmt(n.Select)
		return n

	case *ast.BetweenExpr:
		n.Expr = rewriteExpr(n.Expr)
		n.Low = rewriteExpr(n.Low)
		n.High = rewriteExpr(n.High)
		return n

	case *ast.LikeExpr:
		n.Expr = rewriteExpr(n.Expr)
		n.Pattern = rewriteExpr(n.Pattern)
		if n.Escape != nil {
			n.Escape = rewriteExpr(n.Escape)
		}
		return n

	case *ast.IsExpr:
		n.Expr = rewriteExpr(n.Expr)
		// MariaDB leaves an IS TRUE / IS FALSE operand as written: `a IS TRUE`
		// over an integer column is stored as `a is true`, where MySQL 8.0
		// wraps the operand. Measured on mariadb:11.8.8.
		return n

	case *ast.CaseExpr:
		if n.Operand != nil {
			n.Operand = rewriteExpr(n.Operand)
		}
		for _, w := range n.Whens {
			w.Cond = rewriteExpr(w.Cond)
			w.Result = rewriteExpr(w.Result)
		}
		if n.Default != nil {
			n.Default = rewriteExpr(n.Default)
		}
		return n

	case *ast.FuncCallExpr:
		for i, arg := range n.Args {
			n.Args[i] = rewriteExpr(arg)
		}
		// An aggregate's own ORDER BY is rewritten too: measured on
		// mariadb:11.8.8, GROUP_CONCAT(a ORDER BY NOT (a > 0)) is stored as
		// order by a <= 0. The separator is left alone because it can only be
		// a string literal, which no rewrite touches.
		rewriteOrderBy(n.OrderBy)
		rewriteWindowDef(n.Over)
		return n

	case *ast.RowExpr:
		// A row's items are rewritten each as its own expression: measured,
		// ROW((2=1) SOUNDS LIKE '0',1) is stored as
		// (soundex(2 = 1) = soundex('0'),1). Leaving them kept SOUNDS LIKE at
		// comparison precedence, where the row's commas are the only
		// separators left and the item reads back as the equality alone.
		for i, item := range n.Items {
			n.Items[i] = rewriteExpr(item)
		}
		return n

	case *ast.IntervalExpr:
		// Measured: INTERVAL (NOT (a > 0)) DAY is stored as
		// interval a <= 0 day.
		n.Value = rewriteExpr(n.Value)
		return n

	case *ast.ExtractExpr:
		// Measured: EXTRACT(DAY FROM ((2=1) SOUNDS LIKE '0')) is stored as
		// extract(day from soundex(2 = 1) = soundex('0')).
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.MatchExpr:
		// Measured: MATCH(c) AGAINST (((2=1) SOUNDS LIKE '0')) is stored as
		// against (soundex(2 = 1) = soundex('0')). No test pins this one:
		// the deparser writes no MATCH yet, and that stored body is one the
		// server itself answers 1064 on, since AGAINST takes only a literal.
		n.Against = rewriteExpr(n.Against)
		return n

	case *ast.PositionExpr:
		n.Substr = rewriteExpr(n.Substr)
		n.Str = rewriteExpr(n.Str)
		return n

	case *ast.CastExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.ConvertExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.CollateExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	default:
		// Leaf nodes (literals, column refs, etc.) — no rewriting needed
		return node
	}
}

// rewriteWindowDef applies the rewrites to a window specification's own
// expressions. The server rewrites what a window partitions and orders by
// exactly as it rewrites a target: measured on mariadb:11.8.8,
// OVER (PARTITION BY (a = 1) SOUNDS LIKE '0') is stored as
// partition by soundex(a = 1) = soundex('0'),
// OVER (PARTITION BY NOT (a > 0)) as partition by a <= 0 and
// OVER (PARTITION BY a AND b) as a <> 0 and b <> 0, with the columns
// qualified, and ORDER BY the same. Leaving the specification alone kept
// SOUNDS LIKE at comparison precedence, which drops the parentheses its left
// operand needs and merges partitions the original holds apart.
//
// A frame bound is left alone: it takes no expression at all, so there is
// nothing in one to rewrite. Measured — ROWS BETWEEN (1 + 1) PRECEDING is
// 1064, and so is a bound spelling a predicate.
func rewriteWindowDef(def *ast.WindowDef) {
	if def == nil {
		return
	}
	for i, expr := range def.PartitionBy {
		def.PartitionBy[i] = RewriteExpr(expr)
	}
	rewriteOrderBy(def.OrderBy)
}

// rewriteOrderBy applies the rewrites to the expressions of an ORDER BY list,
// wherever one hangs: a window's, or an aggregate's own.
func rewriteOrderBy(items []*ast.OrderByItem) {
	for _, item := range items {
		if item != nil {
			item.Expr = RewriteExpr(item.Expr)
		}
	}
}

// rewriteTableExpr applies the rewrites inside a FROM item: a derived table's
// own query and a join's ON condition, through both sides of a join.
//
// Measured on mariadb:11.8.8: `JOIN t t2 ON (t1.a = 1) SOUNDS LIKE '0'` is
// stored as `on(soundex(`t1`.`a` = 1) = soundex('0'))` and a derived table's
// body is rewritten as its own query.
func rewriteTableExpr(src ast.TableExpr) {
	switch n := src.(type) {
	case *ast.SubqueryExpr:
		RewriteSelectStmt(n.Select)
	case *ast.JoinClause:
		rewriteTableExpr(n.Left)
		rewriteTableExpr(n.Right)
		if on, ok := n.Condition.(*ast.OnCondition); ok {
			on.Expr = RewriteExpr(on.Expr)
		}
	}
}

// invertOp maps comparison operators to their NOT-inverted counterparts.
var invertOp = map[ast.BinaryOp]ast.BinaryOp{
	ast.BinOpGt: ast.BinOpLe,
	ast.BinOpLt: ast.BinOpGe,
	ast.BinOpGe: ast.BinOpLt,
	ast.BinOpLe: ast.BinOpGt,
	ast.BinOpEq: ast.BinOpNe,
	ast.BinOpNe: ast.BinOpEq,
}

// isComparisonOp returns true if op is a comparison operator that can be inverted.
func isComparisonOp(op ast.BinaryOp) bool {
	_, ok := invertOp[op]
	return ok
}

// unwrapParen strips ParenExpr wrappers to get the inner expression.
func unwrapParen(node ast.ExprNode) ast.ExprNode {
	for {
		p, ok := node.(*ast.ParenExpr)
		if !ok {
			return node
		}
		node = p.Expr
	}
}

// rewriteNot applies NOT folding to the operand of a NOT expression.
// MySQL 8.0's resolver:
//   - NOT (comparison) → inverted comparison (e.g., NOT(a > 0) → (a <= 0))
//   - NOT (LIKE) → not((expr like pattern)) — wraps in not(), doesn't fold
//   - NOT (non-boolean) → (0 = expr) — e.g., NOT(a+1) → (0 = (a+1)), NOT(col) → (0 = col)
func rewriteNot(operand ast.ExprNode) ast.ExprNode {
	// NOT of a quantified comparison inverts the operator and swaps the quantifier:
	// `NOT a > ALL (q)` is `a <= ANY (q)`. Flipping the operator alone would turn `NOT
	// a = ANY (q)` into `a <> ANY (q)`, which holds whenever any row of q differs from
	// a, where the original holds only when every row does. `= ANY` is the exception
	// the server makes, since it stores that one as IN: `(not(a in (q)))`.
	if binExpr, ok := unwrapParen(operand).(*ast.BinaryExpr); ok {
		if sub, ok := binExpr.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" {
			inverted, canInvert := invertOp[binExpr.Op]
			if !canInvert || (binExpr.Op == ast.BinOpEq && sub.Quantifier != "ALL") {
				return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}
			}
			dual := *sub
			dual.Quantifier = "ALL"
			if sub.Quantifier == "ALL" {
				dual.Quantifier = "ANY"
			}
			return &ast.BinaryExpr{Loc: binExpr.Loc, Op: inverted, Left: binExpr.Left, Right: &dual}
		}
	}

	inner := unwrapParen(operand)

	// Case 0: De Morgan. MariaDB pushes a NOT through AND and OR rather than
	// writing one over the whole conjunction: `NOT (a > 0 AND b > 0)` is stored
	// as `a <= 0 or b <= 0`, and the push recurses. MySQL 8.0 writes
	// `(not((...)))`. Measured on mariadb:11.8.8.
	if binExpr, ok := inner.(*ast.BinaryExpr); ok {
		switch binExpr.Op {
		case ast.BinOpAnd, ast.BinOpOr:
			dual := ast.BinOpOr
			if binExpr.Op == ast.BinOpOr {
				dual = ast.BinOpAnd
			}
			left, right := binExpr.Left, binExpr.Right
			if needsBooleanContext(left) {
				left = wrapBooleanContext(left)
			}
			if needsBooleanContext(right) {
				right = wrapBooleanContext(right)
			}
			return &ast.BinaryExpr{
				Loc:   binExpr.Loc,
				Op:    dual,
				Left:  rewriteNot(left),
				Right: rewriteNot(right),
			}
		}
	}

	// Case 1: NOT(comparison) → invert the comparison operator
	if binExpr, ok := inner.(*ast.BinaryExpr); ok {
		if inverted, canInvert := invertOp[binExpr.Op]; canInvert {
			return &ast.BinaryExpr{
				Loc:   binExpr.Loc,
				Op:    inverted,
				Left:  binExpr.Left,
				Right: binExpr.Right,
			}
		}
	}

	// Case 1b: a predicate that carries a NOT of its own takes it, rather than
	// a `!` over the whole: `NOT a IS TRUE` → `a is not true`,
	// `NOT a BETWEEN 1 AND 2` → `a not between 1 and 2`, and
	// `NOT a IN (1,2)` → `a not in (1,2)`. Measured on mariadb:11.8.8. An IN
	// over a subquery is the exception below, which the server stores as
	// `!(a in (q))`.
	switch pred := inner.(type) {
	case *ast.IsExpr:
		flipped := *pred
		flipped.Not = !pred.Not
		return &flipped
	case *ast.BetweenExpr:
		flipped := *pred
		flipped.Not = !pred.Not
		return &flipped
	case *ast.InExpr:
		if pred.Select == nil {
			flipped := *pred
			flipped.Not = !pred.Not
			return &flipped
		}
	}

	// Case 2a: NOT(LIKE) folds into the LIKE, which carries a NOT of its own:
	// `a NOT LIKE 'x'`. Measured on mariadb:11.8.8; MySQL 8.0 wraps it in
	// `not(...)` instead.
	if like, ok := inner.(*ast.LikeExpr); ok {
		flipped := *like
		flipped.Not = !like.Not
		return &flipped
	}

	// Case 2b: NOT(REGEXP) — keep as not() wrapping.
	// MySQL 8.0 rewrites NOT REGEXP to (not(regexp_like(...))).
	if binExpr, ok := inner.(*ast.BinaryExpr); ok && binExpr.Op == ast.BinOpRegexp {
		return &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: operand,
		}
	}

	// Case 2c: NOT(a IN (subquery)) — keep as not() wrapping, which is what the server
	// stores, and what `a <> ALL (q)` prints as. `0 = a in (q)` would mean the same but
	// would not read back as what was written.
	if in, ok := inner.(*ast.InExpr); ok && in.Select != nil && !in.Not {
		return &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: operand,
		}
	}

	// Case 3: NOT over a bare column is stored as `col = 0`, the operand first,
	// where MySQL 8.0 writes the literal first. Anything else keeps the NOT,
	// which prints as `!`: `NOT (a + 1)` → `!(a + 1)` and
	// `NOT ifnull(a,0)` → `!ifnull(a,0)`. Measured on mariadb:11.8.8.
	if needsBooleanContext(operand) {
		return &ast.BinaryExpr{
			Op:    ast.BinOpEq,
			Left:  operand,
			Right: &ast.IntLit{Value: 0},
		}
	}
	return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}
}

// needsBooleanContext reports whether node is an operand MariaDB stores with an
// explicit `<> 0`: a bare column reference, and nothing else. Measured on
// mariadb:11.8.8 — see the AND/OR rewrite above.
func needsBooleanContext(node ast.ExprNode) bool {
	_, ok := unwrapParen(node).(*ast.ColumnRef)
	return ok
}

// wrapBooleanContext wraps a non-boolean expression in `expr <> 0` for boolean
// context (AND/OR/XOR operands, IS TRUE/IS FALSE operands). Measured on
// mariadb:11.8.8: `a AND b` over integer columns is stored as
// `a <> 0 and b <> 0`, the operand first, where MySQL 8.0 writes the literal
// first.
func wrapBooleanContext(node ast.ExprNode) ast.ExprNode {
	return &ast.BinaryExpr{
		Op:    ast.BinOpNe,
		Left:  node,
		Right: &ast.IntLit{Value: 0},
	}
}

// inlineNamedWindows replaces each OVER reference to a window the statement's
// WINDOW clause defines with that definition and drops the clause. MariaDB
// stores no WINDOW clause: `rank() OVER w ... WINDOW w AS (ORDER BY a)` comes
// back as `rank() over ( order by a)`. Measured on mariadb:11.8.8.
func inlineNamedWindows(stmt *ast.SelectStmt) {
	// A query with no clause of its own still walks, because a query nested in
	// it may have one.
	if stmt == nil {
		return
	}
	// A clause entry may itself reference an earlier one, so each is resolved
	// as the map is built and the map holds the resolved window. A reference
	// reaches only backwards, which is what makes one pass enough.
	//
	// Measured on mariadb:11.8.8:
	// `WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)` referenced as w2
	// comes back `over ( partition by `t`.`a` order by `t`.`b`)`, and a chain
	// through a third entry resolves the same way. Taking the entry alone
	// would drop the partitioning it inherits.
	byName := make(map[string]*ast.WindowDef, len(stmt.WindowClause))
	for _, wd := range stmt.WindowClause {
		resolved := wd
		if wd.RefName != "" {
			if def, ok := byName[strings.ToLower(wd.RefName)]; ok {
				resolved = extendWindowDef(def, wd)
				resolved.Name = wd.Name
			}
		}
		byName[strings.ToLower(wd.Name)] = resolved
	}
	// The call can sit anywhere in an expression — `1 + ROW_NUMBER() OVER w`
	// — and in ORDER BY as well as the target list, so every window reference
	// is found by walking rather than by looking only at the top of each
	// target. Dropping the clause while leaving a reference behind would make
	// the stored body name a window it does not define.
	//
	// A window name belongs to the query that declares it, so the walk stops
	// where a nested query begins and that query is inlined against its own
	// clause. Measured on mariadb:11.8.8: under an outer `WINDOW w AS (ORDER
	// BY a)`, the body of
	// `(SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1)` comes
	// back `over ( order by `t`.`b`)`, its own definition, and neither WINDOW
	// clause is left in the stored body.
	inline := func(n ast.Node) bool {
		if nested, ok := n.(*ast.SelectStmt); ok && nested != stmt {
			inlineNamedWindows(nested)
			return false
		}
		call, ok := n.(*ast.FuncCallExpr)
		if !ok || call.Over == nil || call.Over.RefName == "" {
			return true
		}
		if def, ok := byName[strings.ToLower(call.Over.RefName)]; ok {
			call.Over = extendWindowDef(def, call.Over)
		}
		return true
	}
	for _, cte := range stmt.CTEs {
		inlineNamedWindows(cte.Select)
	}
	for _, target := range stmt.TargetList {
		ast.Inspect(target, inline)
	}
	for _, src := range stmt.From {
		ast.Inspect(src, inline)
	}
	if stmt.Where != nil {
		ast.Inspect(stmt.Where, inline)
	}
	for _, item := range stmt.GroupBy {
		ast.Inspect(item.Expr, inline)
	}
	if stmt.Having != nil {
		ast.Inspect(stmt.Having, inline)
	}
	for _, item := range stmt.OrderBy {
		ast.Inspect(item.Expr, inline)
	}
	stmt.WindowClause = nil
}

// extendWindowDef resolves a reference to a named window. A reference may
// carry a specification of its own, which extends the named window instead of
// replacing it: the name supplies what it declares and the reference adds what
// the name leaves out.
//
// Measured on mariadb:11.8.8. `OVER (w ORDER BY b)` over
// `WINDOW w AS (PARTITION BY a)` comes back `over ( partition by `t`.`a` order
// by `t`.`b`)`, and `OVER (w ROWS UNBOUNDED PRECEDING)` over
// `WINDOW w AS (ORDER BY a)` keeps the ordering and the frame together. The
// server refuses the cases where the two would collide rather than choosing
// between them: a reference adding PARTITION BY is 4011, one adding ORDER BY
// to a name that has one is 4012, and one adding a frame to a name that has
// one is 4013. So each field here can only come from one side.
func extendWindowDef(def, ref *ast.WindowDef) *ast.WindowDef {
	extended := *def
	extended.Name = ""
	extended.RefName = ""
	if len(ref.PartitionBy) > 0 {
		extended.PartitionBy = ref.PartitionBy
	}
	if len(ref.OrderBy) > 0 {
		extended.OrderBy = ref.OrderBy
	}
	if ref.Frame != nil {
		extended.Frame = ref.Frame
	}
	return &extended
}

// applyCTEColumnNames renames the output aliases of each CTE's query to the
// names the CTE declares, which is where MariaDB keeps them:
// `WITH cte(x, y) AS (SELECT a, b FROM t)` comes back as
// `with cte(`x`,`y`) as (select `t`.`a` AS `x`,`t`.`b` AS `y` from `t`)`.
// Measured on mariadb:11.8.8.
func applyCTEColumnNames(ctes []*ast.CommonTableExpr) {
	for _, cte := range ctes {
		if cte == nil || len(cte.Columns) == 0 || cte.Select == nil {
			continue
		}
		leaf := ast.LeftmostQueryLeaf(cte.Select)
		if leaf == nil {
			continue
		}
		renamed := make(map[string]string, len(cte.Columns))
		for i, target := range leaf.TargetList {
			if i >= len(cte.Columns) {
				break
			}
			if rt, ok := target.(*ast.ResTarget); ok {
				if rt.Name != "" && !strings.EqualFold(rt.Name, cte.Columns[i]) {
					renamed[strings.ToLower(rt.Name)] = cte.Columns[i]
				}
				rt.Name = cte.Columns[i]
			}
		}
		// A HAVING reference to an output column is kept by name, so the name
		// the CTE declares replaces the one the target carried: measured on
		// mariadb:11.8.8, `WITH c(x) AS (SELECT a AS z FROM t HAVING z > 0)`
		// comes back having `x` > 0. An ORDER BY or GROUP BY mention of the
		// alias is resolved to the target's expression instead, so it carries
		// no name to rename.
		if len(renamed) > 0 && leaf.Having != nil {
			RenameOutputRefs(leaf.Having, renamed)
		}
	}
}

// RenameOutputRefs renames an unqualified column reference that names a
// renamed output column, which is how a view's or a CTE's own column list
// reaches the HAVING that mentions an alias. A reference to a table's own
// column is qualified by the resolver, so an unqualified one left in the body
// is the output's.
//
// The walk stops where a nested query begins, since that query names its own
// outputs: measured on mariadb:11.8.8,
// `CREATE VIEW v (x) AS SELECT a AS z FROM t HAVING z > 0 AND
// EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)` comes back with the
// outer reference as `x` and the derived table's own still `z`.
func RenameOutputRefs(node ast.Node, renamed map[string]string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.SelectStmt); ok {
			return false
		}
		cr, ok := n.(*ast.ColumnRef)
		if !ok || cr.Table != "" || cr.Schema != "" {
			return true
		}
		if name, ok := renamed[strings.ToLower(cr.Column)]; ok {
			cr.Column = name
		}
		return true
	})
}

// RewriteSelectStmt applies RewriteExpr to all expression positions in a
// SelectStmt, which includes the queries nested in it: a CTE's body, a derived
// table's, and a subquery's in an expression are each rewritten as a query of
// their own. Measured on mariadb:11.8.8.
// This should be called after resolver but before deparsing.
func RewriteSelectStmt(stmt *ast.SelectStmt) {
	if stmt == nil {
		return
	}
	inlineNamedWindows(stmt)
	applyCTEColumnNames(stmt.CTEs)
	for _, cte := range stmt.CTEs {
		if cte != nil {
			RewriteSelectStmt(cte.Select)
		}
	}
	if stmt.SetOp != ast.SetOpNone {
		RewriteSelectStmt(stmt.Left)
		RewriteSelectStmt(stmt.Right)
		// The ORDER BY over the whole set operation belongs to neither side,
		// and the server rewrites it as it does any other: measured on
		// mariadb:11.8.8, `SELECT 2 AS a UNION SELECT 3 ORDER BY ((a=2)
		// SOUNDS LIKE '0'), a DESC` is stored as
		// order by soundex(`a` = 2) = soundex('0'),`a` desc.
		rewriteOrderBy(stmt.OrderBy)
		return
	}
	if stmt.ParenSource != nil {
		RewriteSelectStmt(stmt.ParenSource)
		rewriteOrderBy(stmt.OrderBy)
		return
	}
	for i, target := range stmt.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok {
			rt.Val = RewriteExpr(rt.Val)
		} else {
			stmt.TargetList[i] = RewriteExpr(target)
		}
	}
	for _, src := range stmt.From {
		rewriteTableExpr(src)
	}
	if stmt.Where != nil {
		stmt.Where = RewriteExpr(stmt.Where)
	}
	for _, item := range stmt.GroupBy {
		item.Expr = RewriteExpr(item.Expr)
	}
	if stmt.Having != nil {
		stmt.Having = RewriteExpr(stmt.Having)
	}
	rewriteOrderBy(stmt.OrderBy)
}
