package parser

import (
	"strconv"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The tree each operator form parses to. operator_forms_test.go says whether a
// statement parses; none of that sees a wrong grouping, and the grouping is what
// a consumer reads. `NOT a = b` read as `(NOT a) = b` parsed fine for as long as it
// was wrong, and every view and generated column over it stored the wrong meaning.
//
// Every expected shape here is held to MariaDB 11.8.8 by TestOperatorTreesMatchEngine,
// which is what makes these more than a recording of what this parser does today.

// exprShape prints e with every operator parenthesized, so that two trees print
// the same only if they group the same. The output is SQL the engine reads.
func exprShape(t testing.TB, e nodes.ExprNode) string {
	t.Helper()
	not := func(b bool) string {
		if b {
			return "NOT "
		}
		return ""
	}
	switch n := e.(type) {
	case *nodes.ColumnRef:
		return n.Column
	case *nodes.IntLit:
		return strconv.FormatInt(n.Value, 10)
	case *nodes.StringLit:
		return "'" + n.Value + "'"
	case *nodes.NullLit:
		return "NULL"
	case *nodes.VariableRef:
		return "@" + n.Name
	case *nodes.ParenExpr:
		return exprShape(t, n.Expr)
	case *nodes.CollateExpr:
		return "(" + exprShape(t, n.Expr) + " COLLATE " + n.Collation + ")"
	case *nodes.BinaryExpr:
		return "(" + exprShape(t, n.Left) + " " + binaryOpText(t, n.Op) + " " + exprShape(t, n.Right) + ")"
	case *nodes.UnaryExpr:
		op, ok := map[nodes.UnaryOp]string{
			nodes.UnaryMinus:  "- ",
			nodes.UnaryPlus:   "+ ",
			nodes.UnaryNot:    "NOT ",
			nodes.UnaryBitNot: "~ ",
		}[n.Op]
		if !ok {
			t.Fatalf("exprShape: unary operator %d", n.Op)
		}
		// NOT REGEXP is a NOT over a REGEXP in this tree, and prints as the
		// operator it was written as so a reader can compare it with the input.
		if b, isBin := n.Operand.(*nodes.BinaryExpr); isBin && n.Op == nodes.UnaryNot && b.Op == nodes.BinOpRegexp {
			return "(" + exprShape(t, b.Left) + " NOT REGEXP " + exprShape(t, b.Right) + ")"
		}
		return "(" + op + exprShape(t, n.Operand) + ")"
	case *nodes.IsExpr:
		test, ok := map[nodes.IsTestType]string{
			nodes.IsNull:    "NULL",
			nodes.IsTrue:    "TRUE",
			nodes.IsFalse:   "FALSE",
			nodes.IsUnknown: "UNKNOWN",
		}[n.Test]
		if !ok {
			t.Fatalf("exprShape: IS test %d", n.Test)
		}
		return "(" + exprShape(t, n.Expr) + " IS " + not(n.Not) + test + ")"
	case *nodes.BetweenExpr:
		return "(" + exprShape(t, n.Expr) + " " + not(n.Not) + "BETWEEN " + exprShape(t, n.Low) + " AND " + exprShape(t, n.High) + ")"
	case *nodes.InExpr:
		items := make([]string, 0, len(n.List))
		for _, item := range n.List {
			items = append(items, exprShape(t, item))
		}
		return "(" + exprShape(t, n.Expr) + " " + not(n.Not) + "IN (" + strings.Join(items, ", ") + "))"
	case *nodes.LikeExpr:
		return "(" + exprShape(t, n.Expr) + " " + not(n.Not) + "LIKE " + exprShape(t, n.Pattern) + ")"
	}
	t.Fatalf("exprShape: node %T", e)
	return ""
}

// binaryOpText is the operator's SQL spelling, whichever way it was written.
func binaryOpText(t testing.TB, op nodes.BinaryOp) string {
	t.Helper()
	text, ok := map[nodes.BinaryOp]string{
		nodes.BinOpAdd:        "+",
		nodes.BinOpSub:        "-",
		nodes.BinOpMul:        "*",
		nodes.BinOpDiv:        "/",
		nodes.BinOpMod:        "%",
		nodes.BinOpDivInt:     "DIV",
		nodes.BinOpEq:         "=",
		nodes.BinOpNe:         "<>",
		nodes.BinOpLt:         "<",
		nodes.BinOpGt:         ">",
		nodes.BinOpLe:         "<=",
		nodes.BinOpGe:         ">=",
		nodes.BinOpNullSafeEq: "<=>",
		nodes.BinOpAnd:        "AND",
		nodes.BinOpOr:         "OR",
		nodes.BinOpXor:        "XOR",
		nodes.BinOpBitAnd:     "&",
		nodes.BinOpBitOr:      "|",
		nodes.BinOpBitXor:     "^",
		nodes.BinOpShiftLeft:  "<<",
		nodes.BinOpShiftRight: ">>",
		nodes.BinOpRegexp:     "REGEXP",
		nodes.BinOpAssign:     ":=",
		nodes.BinOpSoundsLike: "SOUNDS LIKE",
	}[op]
	if !ok {
		t.Fatalf("binaryOpText: operator %d", op)
	}
	return text
}

// selectExpr parses `SELECT expr` and returns the expression.
func selectExpr(t testing.TB, expr string) nodes.ExprNode {
	t.Helper()
	list, err := Parse("SELECT " + expr)
	if err != nil {
		t.Fatalf("Parse(%q): %v", "SELECT "+expr, err)
	}
	stmt, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok || len(stmt.TargetList) != 1 {
		t.Fatalf("Parse(%q): not a one-column SELECT", "SELECT "+expr)
	}
	if target, ok := stmt.TargetList[0].(*nodes.ResTarget); ok {
		return target.Val
	}
	return stmt.TargetList[0]
}

// operatorTreeCases is an expression and the grouping it parses to.
var operatorTreeCases = []struct{ expr, shape string }{
	{`a || b`, `(a OR b)`},
	{`a && b`, `(a AND b)`},
	{`a || b || c`, `((a OR b) OR c)`},
	{`a && b && c`, `((a AND b) AND c)`},
	{`a || b && c`, `(a OR (b AND c))`},
	{`a && b || c`, `((a AND b) OR c)`},
	{`a || b xor c`, `(a OR (b XOR c))`},
	{`a xor b && c`, `(a XOR (b AND c))`},
	{`a or b || c`, `((a OR b) OR c)`},
	{`a || b or c`, `((a OR b) OR c)`},
	{`a and b && c`, `((a AND b) AND c)`},
	{`a && b and c`, `((a AND b) AND c)`},
	{`a = 1 || b = 2 && c = 3`, `((a = 1) OR ((b = 2) AND (c = 3)))`},
	{`not a || b`, `((NOT a) OR b)`},
	{`a || not b`, `(a OR (NOT b))`},
	{`a && not b`, `(a AND (NOT b))`},
	{`! a || b`, `((NOT a) OR b)`},
	{`a | b || c`, `((a | b) OR c)`},
	{`a || b | c`, `(a OR (b | c))`},
	{`not a = b`, `(NOT (a = b))`},
	{`not a <> b`, `(NOT (a <> b))`},
	{`not a < b`, `(NOT (a < b))`},
	{`not a in (1)`, `(NOT (a IN (1)))`},
	{`not a is null`, `(NOT (a IS NULL))`},
	{`not a is true`, `(NOT (a IS TRUE))`},
	{`not a between 1 and 2`, `(NOT (a BETWEEN 1 AND 2))`},
	{`not a like b`, `(NOT (a LIKE b))`},
	{`not a + 1`, `(NOT (a + 1))`},
	{`not a and b`, `((NOT a) AND b)`},
	{`not a or b`, `((NOT a) OR b)`},
	{`not not a`, `(NOT (NOT a))`},
	{`! a = b`, `((NOT a) = b)`},
	{`! a + 1`, `((NOT a) + 1)`},
	{`- a = b`, `((- a) = b)`},
	{`a = b in (1)`, `(a = (b IN (1)))`},
	{`a in (1) = b`, `((a IN (1)) = b)`},
	{`a = b in (1) = c`, `((a = (b IN (1))) = c)`},
	{`a = b like c`, `(a = (b LIKE c))`},
	{`a like b = c`, `((a LIKE b) = c)`},
	{`a <> b not in (1)`, `(a <> (b NOT IN (1)))`},
	{`a < b between 1 and 2`, `(a < (b BETWEEN 1 AND 2))`},
	{`a between 1 and 2 = b`, `((a BETWEEN 1 AND 2) = b)`},
	{`a = b is null`, `((a = b) IS NULL)`},
	{`a is null = b`, `((a IS NULL) = b)`},
	{`a = b is true`, `((a = b) IS TRUE)`},
	{`a + 1 in (1)`, `((a + 1) IN (1))`},
	{`a | 1 in (1)`, `((a | 1) IN (1))`},
	{`a + 1 like b`, `((a + 1) LIKE b)`},
	{`a between 1 and 2 | 3`, `(a BETWEEN 1 AND (2 | 3))`},
	{`a between 1 | 2 and 3`, `(a BETWEEN (1 | 2) AND 3)`},
	{`a between 1 and 2 + 3`, `(a BETWEEN 1 AND (2 + 3))`},
	{`a between 1 and 2 and b`, `((a BETWEEN 1 AND 2) AND b)`},
	{`a between 1 and 2 or b`, `((a BETWEEN 1 AND 2) OR b)`},
	{`a between 1 and b in (1)`, `(a BETWEEN 1 AND (b IN (1)))`},
	{`@a := 1`, `(@a := 1)`},
	{`@a := @b := 1`, `(@a := (@b := 1))`},
	{`@a := 1 + 1`, `(@a := (1 + 1))`},
	{`@a := 1 or 2`, `(@a := (1 OR 2))`},
	{`@a := 1 = 2`, `(@a := (1 = 2))`},
	{`not @a := 1`, `(NOT (@a := 1))`},
	{`1 or @a := 2 or 3`, `(1 OR (@a := (2 OR 3)))`},
	{`1 and @a := 2 and 3`, `(1 AND (@a := (2 AND 3)))`},
	{`1 + @a := 2 + 3`, `(1 + (@a := (2 + 3)))`},
	{`- @a := 2 + 3`, `(- (@a := (2 + 3)))`},
	{`1 = @a := 2 = 3`, `(1 = (@a := (2 = 3)))`},
	{`a in (1) in (2)`, `((a IN (1)) IN (2))`},
	{`a like b like c`, `((a LIKE b) LIKE c)`},
	{`a like b + c`, `(a LIKE (b + c))`},
	{`a like b not like c`, `(a LIKE (b NOT LIKE c))`},
	{`a like b not regexp c`, `(a LIKE (b NOT REGEXP c))`},
	{`a regexp b not in (1)`, `(a REGEXP (b NOT IN (1)))`},
	{`a like b regexp c`, `((a LIKE b) REGEXP c)`},
	{`a like b in (1)`, `((a LIKE b) IN (1))`},
	{`a is true = b`, `((a IS TRUE) = b)`},
	{`a is true is null`, `((a IS TRUE) IS NULL)`},
	{`a between 1 and 2 between 3 and 4`, `(a BETWEEN 1 AND (2 BETWEEN 3 AND 4))`},
	{`a between b in (1) and c`, `(a BETWEEN (b IN (1)) AND c)`},
	{`a sounds like b sounds like c`, `((a SOUNDS LIKE b) SOUNDS LIKE c)`},
}

func TestOperatorTrees(t *testing.T) {
	for _, c := range operatorTreeCases {
		t.Run(c.expr, func(t *testing.T) {
			if got := exprShape(t, selectExpr(t, c.expr)); got != c.shape {
				t.Errorf("SELECT %s\n got %s\nwant %s", c.expr, got, c.shape)
			}
		})
	}
}

// TestOperatorSpelling pins what a shape cannot: `||` and `&&` are OR and AND with
// their spelling kept, which is what the auto-alias of a view column prints and how
// a consumer can tell them apart when the session has PIPES_AS_CONCAT.
func TestOperatorSpelling(t *testing.T) {
	for _, c := range []struct {
		expr     string
		op       nodes.BinaryOp
		original string
	}{
		{"a || b", nodes.BinOpOr, "||"},
		{"a OR b", nodes.BinOpOr, ""},
		{"a && b", nodes.BinOpAnd, "&&"},
		{"a AND b", nodes.BinOpAnd, ""},
	} {
		bin, ok := selectExpr(t, c.expr).(*nodes.BinaryExpr)
		if !ok {
			t.Fatalf("%s: not a BinaryExpr", c.expr)
		}
		if bin.Op != c.op || bin.OriginalOp != c.original {
			t.Errorf("%s: op %d spelled %q, want op %d spelled %q", c.expr, bin.Op, bin.OriginalOp, c.op, c.original)
		}
	}
}

// TestOperatorLocs holds the engine contract's Loc rule on the nodes this grammar
// builds: from the first token of the left operand to the last token of the whole.
func TestOperatorLocs(t *testing.T) {
	loc := func(e nodes.ExprNode) nodes.Loc {
		switch n := e.(type) {
		case *nodes.BinaryExpr:
			return n.Loc
		case *nodes.UnaryExpr:
			return n.Loc
		}
		t.Fatalf("no Loc read for %T", e)
		return nodes.Loc{}
	}
	for _, expr := range []string{"a || b", "a && b", "not a = b", "@a := 1 + 1", "a || b && c", "not a || b"} {
		sql := "SELECT " + expr + " "
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		val := list.Items[0].(*nodes.SelectStmt).TargetList[0]
		if target, ok := val.(*nodes.ResTarget); ok {
			val = target.Val
		}
		got := loc(val)
		if got.Start != len("SELECT ") || got.End != len("SELECT ")+len(expr) {
			t.Errorf("%q: Loc %d..%d covers %q, want the whole expression", expr, got.Start, got.End, sql[got.Start:got.End])
		}
	}
}

// TestQuantifierKept pins that the quantifier reaches the tree: a consumer that
// loses it reads `a > ALL (q)` as a scalar comparison against q.
func TestQuantifierKept(t *testing.T) {
	for _, c := range []struct{ expr, quantifier string }{
		{"a = ANY (SELECT 1)", "ANY"},
		{"a = some (SELECT 1)", "SOME"},
		{"a > ALL (SELECT 1)", "ALL"},
		{"a = (SELECT 1)", ""},
	} {
		bin, ok := selectExpr(t, c.expr).(*nodes.BinaryExpr)
		if !ok {
			t.Fatalf("%s: not a BinaryExpr", c.expr)
		}
		sub, ok := bin.Right.(*nodes.SubqueryExpr)
		if !ok {
			t.Fatalf("%s: right side is %T", c.expr, bin.Right)
		}
		if sub.Quantifier != c.quantifier {
			t.Errorf("%s: quantifier %q, want %q", c.expr, sub.Quantifier, c.quantifier)
		}
	}
}
