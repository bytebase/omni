package parser

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// A parenthesized scalar subquery may be the LEFT operand of a binary operator
// while the whole comparison itself sits inside parentheses — function args,
// WHERE parens, view bodies: `if(((SELECT count(0) FROM t) = 0),'NO','YES')`.
// The paren-expression parser used to classify the leading '(' run as a
// parenthesized query expression and hard-expect ')' right after the subquery,
// so the operator was "unexpected token". The stock MySQL 8.0 sys.metrics view
// stores exactly this shape (`if((((select count(0) from ...))) = 0),'NO',...)`),
// which made a canonical dump of sys unloadable.
//
// Oracle evidence (live MySQL 8.0.32 and 5.7.25, see PR body): both versions
// accept every operator after a parenthesized subquery operand inside parens,
// row constructors and IN-lists whose first element starts with a parenthesized
// subquery, and depth-2..4 paren nesting. 8.0 additionally accepts set-op /
// LIMIT / WITH / VALUES / TABLE headed query expressions in that position
// (5.7 has no such syntax); omni parses the 8.0 grammar.

// TestParenSubqueryOperandParses is the acceptance matrix: every form below is
// accepted by live MySQL (8.0.32 and 5.7.25 unless noted 8.0-only).
func TestParenSubqueryOperandParses(t *testing.T) {
	for _, sql := range []string{
		// Comparison operators, subquery LEFT, inside parens.
		"SELECT ((SELECT count(0) FROM t) = 0)",
		"SELECT ((SELECT count(0) FROM t) <> 0)",
		"SELECT ((SELECT count(0) FROM t) != 0)",
		"SELECT ((SELECT count(0) FROM t) < 0)",
		"SELECT ((SELECT count(0) FROM t) > 0)",
		"SELECT ((SELECT count(0) FROM t) <= 0)",
		"SELECT ((SELECT count(0) FROM t) >= 0)",
		"SELECT ((SELECT count(0) FROM t) <=> 0)",
		// Arithmetic / bit operators.
		"SELECT ((SELECT count(0) FROM t) + 1)",
		"SELECT ((SELECT count(0) FROM t) - 1)",
		"SELECT ((SELECT count(0) FROM t) * 2)",
		"SELECT ((SELECT count(0) FROM t) / 2)",
		"SELECT ((SELECT count(0) FROM t) DIV 2)",
		"SELECT ((SELECT count(0) FROM t) MOD 2)",
		"SELECT ((SELECT count(0) FROM t) % 2)",
		"SELECT ((SELECT max(a) FROM t) | 1)",
		"SELECT ((SELECT max(a) FROM t) & 1)",
		"SELECT ((SELECT max(a) FROM t) << 1)",
		"SELECT ((SELECT max(a) FROM t) ^ 1)",
		// Keyword predicates.
		"SELECT ((SELECT max(a) FROM t) LIKE '%1%')",
		"SELECT ((SELECT max(a) FROM t) NOT LIKE '%1%')",
		"SELECT ((SELECT max(a) FROM t) IN (1,2,3))",
		"SELECT ((SELECT max(a) FROM t) NOT IN (1,2))",
		"SELECT ((SELECT max(a) FROM t) IS NULL)",
		"SELECT ((SELECT max(a) FROM t) IS NOT NULL)",
		"SELECT ((SELECT 1) IS TRUE)",
		"SELECT ((SELECT NULL) IS UNKNOWN)",
		"SELECT ((SELECT max(a) FROM t) BETWEEN 0 AND 10)",
		"SELECT ((SELECT max(a) FROM t) AND 1)",
		"SELECT ((SELECT max(a) FROM t) OR 0)",
		"SELECT ((SELECT max(a) FROM t) XOR 0)",
		"SELECT ((SELECT max(a) FROM t) REGEXP '1')",
		"SELECT ((SELECT 'a') COLLATE utf8mb4_bin = 'a')",
		// Nesting depth 2-4 (sys.metrics stores depth 3).
		"SELECT (((SELECT count(0) FROM t)) = 0)",
		"SELECT ((((SELECT count(0) FROM t))) = 0)",
		"SELECT (((((SELECT count(0) FROM t)))) + 1)",
		"SELECT ((((SELECT count(0) FROM t)) = 0))",
		"SELECT (((((SELECT count(0) FROM t)) = 0)))",
		"SELECT ((((SELECT 1) = 0)))",
		"SELECT (((SELECT 1) = 0) AND ((SELECT 2) = 2))",
		"SELECT ((((SELECT 1) = 0)) OR (((SELECT 2)) = 2))",
		"SELECT (((SELECT count(0) FROM t)) DIV 1)",
		// RIGHT-operand order (regression: already worked).
		"SELECT (0 = (SELECT count(0) FROM t))",
		"SELECT (0 = ((SELECT count(0) FROM t)))",
		"SELECT if((0 = ((SELECT count(0) FROM t))),'NO','YES')",
		// Both operands subqueries; chained operators.
		"SELECT ((SELECT count(0) FROM t) = (SELECT count(0) FROM t))",
		"SELECT (((SELECT count(0) FROM t)) = ((SELECT count(0) FROM t)))",
		"SELECT ((SELECT max(a) FROM t) - (SELECT min(a) FROM t) = 1)",
		"SELECT ((SELECT max(a) FROM t) + 1 * 2 - 3)",
		"SELECT ((SELECT 1) AND (SELECT 1))",
		"SELECT ((SELECT min(a) FROM t) BETWEEN (SELECT min(a) FROM t) AND (SELECT max(a) FROM t))",
		// Top level (regression: already worked).
		"SELECT (SELECT count(0) FROM t) = 0",
		"SELECT ((SELECT count(0) FROM t)) = 0",
		"SELECT (((SELECT count(0) FROM t))) = 0",
		// Contexts: function args, WHERE, HAVING, JOIN ON, CASE, ORDER BY.
		"SELECT if(((SELECT count(0) FROM t) = 0),'NO','YES')",
		"SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES')",
		"SELECT if(((((SELECT count(0) FROM t))) = 0),'NO','YES')",
		"SELECT ifnull(((SELECT max(a) FROM t) + 1), 0)",
		"SELECT 1 FROM t WHERE ((SELECT count(0) FROM t) = 2)",
		"SELECT 1 FROM t WHERE (((SELECT count(0) FROM t)) = 2)",
		"SELECT a FROM t GROUP BY a HAVING ((SELECT count(0) FROM t) = 2)",
		"SELECT 1 FROM t t1 JOIN t t2 ON ((SELECT count(0) FROM t) = 2)",
		"SELECT 1 FROM t t1 JOIN t t2 ON (((SELECT count(0) FROM t)) = 2)",
		"SELECT CASE WHEN ((SELECT count(0) FROM t) = 2) THEN 1 ELSE 0 END",
		"SELECT a FROM t ORDER BY ((SELECT count(0) FROM t) = 2)",
		"SELECT 1 FROM t WHERE (((SELECT max(a) FROM t)) IS NOT NULL)",
		"SELECT if(((SELECT count(0) FROM t) BETWEEN 0 AND 9),'a','b')",
		// EXISTS as a parenthesized operand (EXISTS itself is the primary).
		"SELECT (EXISTS(SELECT 1 FROM t)) = 1",
		"SELECT ((EXISTS(SELECT 1 FROM t)) = 1)",
		"SELECT ((EXISTS(SELECT 1 FROM t)) + 1)",
		"SELECT (NOT EXISTS(SELECT 1 FROM t)) = 0",
		// NOT / ! prefix around the comparison.
		"SELECT (NOT (SELECT 1) = 2)",
		"SELECT (!(SELECT 1) = 2)",
		// Row constructors whose first item is a parenthesized subquery.
		"SELECT ((SELECT 1), 2) = ROW(1,2)",
		"SELECT (((SELECT 1)), 2) = ROW(1,2)",
		"SELECT ((SELECT 1), (SELECT 2)) = ROW(1,2)",
		"SELECT ((SELECT 1) = 0, 2) = ROW(0,2)",
		// IN-lists whose first element starts with a parenthesized subquery.
		"SELECT 1 IN ((SELECT 1) = 1)",
		"SELECT 1 IN ((SELECT 1), 2)",
		"SELECT 1 IN ((SELECT 1), (SELECT 2))",
		"SELECT 1 IN ((SELECT max(a) FROM t) + 1, 2)",
		"SELECT 1 IN (((SELECT max(a) FROM t)) + 1, 2)",
		"SELECT 1 IN (((SELECT 1) = 1))",
		"SELECT 1 IN ((((SELECT 1)) = 1))",
		"SELECT 1 IN (((SELECT 1) = 1), 2)",
		// 8.0-only heads in the operand position (5.7 lacks the syntax).
		"SELECT (((SELECT 1) UNION (SELECT 1)) = 1)",
		"SELECT ((SELECT 1 UNION SELECT 1) = 1)",
		"SELECT (((SELECT 1) LIMIT 1) = 1)",
		"SELECT (((SELECT a FROM t ORDER BY a LIMIT 1)) = 1)",
		"SELECT ((WITH cte AS (SELECT 1 AS c) SELECT c FROM cte) = 1)",
		"SELECT ((VALUES ROW(1)) = 1)",
		"SELECT ((TABLE t2) = 1)",
		// Quantified comparison keeps working around these shapes.
		"SELECT (1 = ANY ((SELECT a FROM t)))",
		"SELECT ((SELECT max(a) FROM t) = ANY (SELECT a FROM t))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenSubqueryOperandRejected pins forms live MySQL rejects (1064 on both
// 8.0.32 and 5.7.25): EXISTS takes only a table subquery, and unterminated
// parens stay errors.
func TestParenSubqueryOperandRejected(t *testing.T) {
	for _, sql := range []string{
		"SELECT EXISTS((SELECT 1 FROM t) = 1)",
		"SELECT ((SELECT 1) = ",
		"SELECT ((SELECT 1) = 1",
		"SELECT ((SELECT 1) =)",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

// TestParenSubqueryOperandAST pins the tree shapes the fix produces, and that
// the pure forms keep their pre-fix shapes (no ParenExpr wrapper creep).
func TestParenSubqueryOperandAST(t *testing.T) {
	// ((SELECT ...) = 0) → ParenExpr{BinaryExpr{Left: SubqueryExpr}}.
	item := firstSelectItem(t, "SELECT ((SELECT count(0) FROM t) = 0)")
	paren, ok := item.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("expected *ast.ParenExpr, got %T", item)
	}
	bin, ok := paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr inside parens, got %T", paren.Expr)
	}
	if bin.Op != ast.BinOpEq {
		t.Errorf("Op = %v, want BinOpEq", bin.Op)
	}
	if _, ok := bin.Left.(*ast.SubqueryExpr); !ok {
		t.Errorf("Left = %T, want *ast.SubqueryExpr", bin.Left)
	}
	if _, ok := bin.Right.(*ast.IntLit); !ok {
		t.Errorf("Right = %T, want *ast.IntLit", bin.Right)
	}

	// (((SELECT ...)) = 0) → the depth-2 subquery is still one SubqueryExpr.
	item = firstSelectItem(t, "SELECT (((SELECT count(0) FROM t)) = 0)")
	paren, ok = item.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("depth-2: expected *ast.ParenExpr, got %T", item)
	}
	bin, ok = paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("depth-2: expected *ast.BinaryExpr, got %T", paren.Expr)
	}
	if _, ok := bin.Left.(*ast.SubqueryExpr); !ok {
		t.Errorf("depth-2: Left = %T, want *ast.SubqueryExpr", bin.Left)
	}

	// ((SELECT 1)) with no trailing operator stays a bare SubqueryExpr.
	item = firstSelectItem(t, "SELECT ((SELECT 1))")
	if _, ok := item.(*ast.SubqueryExpr); !ok {
		t.Errorf("pure ((SELECT 1)): got %T, want *ast.SubqueryExpr", item)
	}

	// ((SELECT 1) UNION (SELECT 2)) stays a query expression, not a binary op.
	item = firstSelectItem(t, "SELECT ((SELECT 1) UNION (SELECT 2))")
	if _, ok := item.(*ast.SubqueryExpr); !ok {
		t.Errorf("set-op: got %T, want *ast.SubqueryExpr", item)
	}

	// ((SELECT 1), 2) → RowExpr{[SubqueryExpr, IntLit]}.
	item = firstSelectItem(t, "SELECT ((SELECT 1), 2) = ROW(1,2)")
	bin, ok = item.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("row: expected *ast.BinaryExpr, got %T", item)
	}
	row, ok := bin.Left.(*ast.RowExpr)
	if !ok {
		t.Fatalf("row: Left = %T, want *ast.RowExpr", bin.Left)
	}
	if len(row.Items) != 2 {
		t.Fatalf("row: len(Items) = %d, want 2", len(row.Items))
	}
	if _, ok := row.Items[0].(*ast.SubqueryExpr); !ok {
		t.Errorf("row: Items[0] = %T, want *ast.SubqueryExpr", row.Items[0])
	}

	// IN (subquery) still resolves to InExpr.Select (not a 1-element list).
	item = firstSelectItem(t, "SELECT 1 IN ((SELECT 1))")
	in, ok := item.(*ast.InExpr)
	if !ok {
		t.Fatalf("in-subquery: expected *ast.InExpr, got %T", item)
	}
	if in.Select == nil || len(in.List) != 0 {
		t.Errorf("in-subquery: Select=%v List=%d, want subquery form", in.Select != nil, len(in.List))
	}

	// IN ((SELECT 1) = 1) is a 1-element VALUE list, not an IN-subquery.
	item = firstSelectItem(t, "SELECT 1 IN ((SELECT 1) = 1)")
	in, ok = item.(*ast.InExpr)
	if !ok {
		t.Fatalf("in-compare: expected *ast.InExpr, got %T", item)
	}
	if in.Select != nil || len(in.List) != 1 {
		t.Fatalf("in-compare: Select=%v List=%d, want 1-element list", in.Select != nil, len(in.List))
	}
	bin, ok = in.List[0].(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("in-compare: List[0] = %T, want *ast.BinaryExpr", in.List[0])
	}
	if _, ok := bin.Left.(*ast.SubqueryExpr); !ok {
		t.Errorf("in-compare: List[0].Left = %T, want *ast.SubqueryExpr", bin.Left)
	}

	// IN ((SELECT 1), 2) is a 2-element value list with a subquery head.
	item = firstSelectItem(t, "SELECT 1 IN ((SELECT 1), 2)")
	in, ok = item.(*ast.InExpr)
	if !ok {
		t.Fatalf("in-list: expected *ast.InExpr, got %T", item)
	}
	if in.Select != nil || len(in.List) != 2 {
		t.Fatalf("in-list: Select=%v List=%d, want 2-element list", in.Select != nil, len(in.List))
	}
	if _, ok := in.List[0].(*ast.SubqueryExpr); !ok {
		t.Errorf("in-list: List[0] = %T, want *ast.SubqueryExpr", in.List[0])
	}

	// The north-star shape: if((((select ...))) = 0),'NO','YES') — first arg is
	// the paren-wrapped comparison with a subquery left operand.
	item = firstSelectItem(t, "SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES')")
	fc, ok := item.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("if(): expected *ast.FuncCallExpr, got %T", item)
	}
	if len(fc.Args) != 3 {
		t.Fatalf("if(): len(Args) = %d, want 3", len(fc.Args))
	}
	paren, ok = fc.Args[0].(*ast.ParenExpr)
	if !ok {
		t.Fatalf("if(): Args[0] = %T, want *ast.ParenExpr", fc.Args[0])
	}
	bin, ok = paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("if(): Args[0].Expr = %T, want *ast.BinaryExpr", paren.Expr)
	}
	if _, ok := bin.Left.(*ast.SubqueryExpr); !ok {
		t.Errorf("if(): comparison Left = %T, want *ast.SubqueryExpr", bin.Left)
	}
}
