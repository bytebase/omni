package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestGraphMatch pins the shapes SQL Server 2022 takes in MATCH, and the ones
// it refuses, as the pinned server answered them under SET PARSEONLY ON:
// statements of the family beginning "IF MATCH(" and "SELECT * FROM a WHERE
// MATCH(" in testdata/compat.tsv.
func TestGraphMatch(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM a WHERE MATCH(a -(e)-> b)",
		"SELECT * FROM a WHERE MATCH(a <-(e)- b)",
		"SELECT * FROM a WHERE MATCH(a -(e)-> b -(e)-> c)",
		"SELECT * FROM a WHERE MATCH(a -(e)-> b AND b -(e)-> c)",
		"SELECT * FROM a WHERE MATCH(LAST_NODE -(e)-> a)",
		"SELECT * FROM a WHERE MATCH(a -(e)-> LAST_NODE(b))",
		"SELECT * FROM a WHERE MATCH(LAST_NODE(a) = LAST_NODE(b))",
		"SELECT * FROM a WHERE MATCH(LAST_NODE(a) = LAST_NODE(b) AND LAST_NODE(b) = LAST_NODE(c))",
		"SELECT * FROM a WHERE MATCH(SHORTEST_PATH(a(-(e)->a) +))",
		"SELECT * FROM a WHERE MATCH(SHORTEST_PATH(a(-(e)->a){1,2}))",
		"SELECT * FROM a WHERE MATCH(SHORTEST_PATH((a -(e)->) + b))",
		"SELECT * FROM a WHERE MATCH(a -(e)-> b) AND a = 1",
		"SELECT * FROM a WHERE NOT MATCH(a -(e)-> b)",
		"IF MATCH(a -(e)-> b) SELECT 1",
		"IF MATCH(LAST_NODE(a) = LAST_NODE(b)) SELECT 1",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// MATCH is a predicate, not a function.
		"SELECT MATCH(a -(e)-> b)",
		// Both sides of the equality must use LAST_NODE, and only "=" compares.
		"IF MATCH(LAST_NODE = a) SELECT 1",
		"IF MATCH(LAST_NODE = LAST_NODE(a)) SELECT 1",
		"IF MATCH(a = b) SELECT 1",
		"IF MATCH(LAST_NODE(a) <> LAST_NODE(b)) SELECT 1",
		"IF MATCH((a) = a) SELECT 1",
		// An edge keeps its identifier and its arrows.
		"IF MATCH(LAST_NODE -(a) > a) SELECT 1",
		"IF MATCH(LAST_NODE -(a) - a) SELECT 1",
		"IF MATCH(a -(a) > a) SELECT 1",
		"IF MATCH(a <(a) - a) SELECT 1",
		"IF MATCH(a < -() - a) SELECT 1",
		// SHORTEST_PATH takes at least one step, and LAST_NODE is refused in it.
		"IF MATCH(SHORTEST_PATH(a() +)) SELECT 1",
		"IF MATCH(SHORTEST_PATH(() + a)) SELECT 1",
		"IF MATCH(SHORTEST_PATH(a(LAST_NODE(b) -(e)-> b) +)) SELECT 1",
		// The quantifier is + or { n , m }.
		"IF MATCH(SHORTEST_PATH(a(-(e)->a){1})) SELECT 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want SQL Server's rejection", sql)
		}
	}
}

// TestGraphMatchAST pins the tree a MATCH builds.
func TestGraphMatchAST(t *testing.T) {
	list, err := Parse("SELECT * FROM a WHERE MATCH(a -(e)-> b)")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	match, ok := sel.WhereClause.(*ast.GraphMatchPredicate)
	if !ok {
		t.Fatalf("WHERE = %T, want *GraphMatchPredicate", sel.WhereClause)
	}
	step, ok := match.Expression.(*ast.GraphMatchCompositeExpr)
	if !ok {
		t.Fatalf("MATCH pattern = %T, want *GraphMatchCompositeExpr", match.Expression)
	}
	if step.Left.Node != "a" || step.Edge != "e" || step.Right.Node != "b" {
		t.Errorf("step = %s-(%s)-> %s, want a-(e)-> b", step.Left.Node, step.Edge, step.Right.Node)
	}
	if !step.ArrowOnRight {
		t.Error("ArrowOnRight = false, want true for ->")
	}
	if got := ast.NodeToString(step); got == "" {
		t.Error("NodeToString(step) = empty")
	}

	list, err = Parse("SELECT * FROM a WHERE MATCH(LAST_NODE(a) = LAST_NODE(b))")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	match = list.Items[0].(*ast.SelectStmt).WhereClause.(*ast.GraphMatchPredicate)
	last, ok := match.Expression.(*ast.GraphMatchLastNodePredicate)
	if !ok {
		t.Fatalf("MATCH pattern = %T, want *GraphMatchLastNodePredicate", match.Expression)
	}
	if !last.Left.UsesLastNode || !last.Right.UsesLastNode {
		t.Errorf("UsesLastNode = %t, %t, want true, true", last.Left.UsesLastNode, last.Right.UsesLastNode)
	}

	list, err = Parse("SELECT * FROM a WHERE MATCH(SHORTEST_PATH(a(-(e)->a){1,2}))")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	match = list.Items[0].(*ast.SelectStmt).WhereClause.(*ast.GraphMatchPredicate)
	rec, ok := match.Expression.(*ast.GraphMatchRecursivePredicate)
	if !ok {
		t.Fatalf("MATCH pattern = %T, want *GraphMatchRecursivePredicate", match.Expression)
	}
	if !rec.AnchorOnLeft || rec.Outer.Node != "a" {
		t.Errorf("recursive predicate = anchorOnLeft %t outer %v, want true and a", rec.AnchorOnLeft, rec.Outer)
	}
	if rec.IsPlusSign || rec.Lower == nil || rec.Upper == nil || rec.Lower.Ival != 1 || rec.Upper.Ival != 2 {
		t.Errorf("quantifier = plus %t, %v, %v, want {1,2}", rec.IsPlusSign, rec.Lower, rec.Upper)
	}
}

// TestIsDistinctFrom pins the IS [NOT] DISTINCT FROM forms, which SQL Server
// takes in a search condition only.
func TestIsDistinctFrom(t *testing.T) {
	for _, sql := range []string{
		"IF 1 IS DISTINCT FROM 1 SELECT 1",
		"IF 1 IS NOT DISTINCT FROM 1 SELECT 1",
		"IF 1 IS DISTINCT FROM NULL SELECT 1",
		"IF 1 IS NOT DISTINCT FROM NULL SELECT 1",
		"IF 1 IS DISTINCT FROM ALL(SELECT 1) SELECT 1",
		"IF 1 IS NOT DISTINCT FROM ALL(SELECT 1) SELECT 1",
		"IF 1 IS DISTINCT FROM (SELECT 1) SELECT 1",
		"SELECT * FROM a WHERE a IS DISTINCT FROM b",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		// A scalar expression, not a search condition.
		"SELECT 1 IS DISTINCT FROM 2",
		"SELECT 1 IS NOT DISTINCT FROM 2",
		// The quantified form takes a subquery.
		"IF 1 IS DISTINCT FROM ALL() SELECT 1",
		"IF 1 IS DISTINCT FROM SOME() SELECT 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want SQL Server's rejection", sql)
		}
	}

	list, err := Parse("SELECT * FROM a WHERE a IS NOT DISTINCT FROM b")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	is, ok := list.Items[0].(*ast.SelectStmt).WhereClause.(*ast.IsExpr)
	if !ok {
		t.Fatalf("WHERE = %T, want *IsExpr", list.Items[0].(*ast.SelectStmt).WhereClause)
	}
	if is.TestType != ast.IsNotDistinctFrom {
		t.Errorf("TestType = %v, want IsNotDistinctFrom", is.TestType)
	}

	list, err = Parse("IF 1 IS DISTINCT FROM ALL(SELECT 1) SELECT 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sub, ok := list.Items[0].(*ast.IfStmt).Condition.(*ast.SubqueryComparisonExpr)
	if !ok {
		t.Fatalf("IF condition = %T, want *SubqueryComparisonExpr", list.Items[0].(*ast.IfStmt).Condition)
	}
	if sub.Distinct != ast.DistinctFrom || sub.Quantifier != "ALL" {
		t.Errorf("predicate = distinct %v quantifier %q, want DistinctFrom and ALL", sub.Distinct, sub.Quantifier)
	}
}

// TestSplitComparisonOperators pins the comparison operators whose two
// characters are separate tokens. SQL Server reads them in a search condition
// only.
func TestSplitComparisonOperators(t *testing.T) {
	for _, sql := range []string{
		"IF 1 ! = 1 SELECT 1",
		"IF 1 ! < 1 SELECT 1",
		"IF 1 ! > 1 SELECT 1",
		"IF 1 < = 1 SELECT 1",
		"IF 1 < > 1 SELECT 1",
		"IF 1 > = 1 SELECT 1",
		"SELECT * FROM a WHERE a ! = 1",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"SELECT 1 ! = 2",
		"SELECT 1 < = 2",
		"SELECT 1 > = 2",
		"SELECT 1 < > 2",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want SQL Server's rejection", sql)
		}
	}

	list, err := Parse("SELECT * FROM a WHERE a ! = 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	cmp, ok := list.Items[0].(*ast.SelectStmt).WhereClause.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("WHERE = %T, want *BinaryExpr", list.Items[0].(*ast.SelectStmt).WhereClause)
	}
	if cmp.Op != ast.BinOpNeq {
		t.Errorf("Op = %v, want BinOpNeq", cmp.Op)
	}
}
