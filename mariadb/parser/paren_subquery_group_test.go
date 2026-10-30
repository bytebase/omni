package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parenSubqueryAccepted are parentheses whose content begins with a
// parenthesized query but is not itself a query: a join group or table list
// whose first operand is a derived table, and an expression whose first
// operand is a subquery. MariaDB 11.8 accepts them all
// (TestParenSubqueryGroupOracle) and writes several in view definitions.
var parenSubqueryAccepted = []string{
	// join groups, as MariaDB writes a view over a derived table in a join
	"select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on(`d`.`k1` = `j2`.`k1`))",
	"select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) as d join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d (x) join j2 on d.x = j2.k1)",
	"select * from ((select k1 from j1) d left join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d right join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d inner join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d cross join j2)",
	"select * from ((select k1 from j1) d join j2 using (k1))",
	"select * from ((select k1 from j1) d natural join j2)",
	"select * from ((select k1 from j1) d straight_join j2)",
	"select * from ((select k1 from j1) d, j2)",
	"select * from ((select k1 from j1) d, (select k1 from j2) e)",
	"select * from (((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1)",
	"select * from ((select k1 from j1) d join ((select k1 from j2) e join j3 on e.k1 = j3.k1) on d.k1 = e.k1)",
	"select * from (((select k1 from j1)) d join j2 on d.k1 = j2.k1)",
	"select * from (((select k1 from j1) union (select k1 from j2)) d join j3 on d.k1 = j3.k1)",
	"select * from ((select k1 from j1 union select k1 from j2) d join j3 on d.k1 = j3.k1)",
	"select * from ((with c as (select 1 k1) select k1 from c) d join j2 on 1)",
	"select * from j1 join ((select 1) d join j2 on 1) on 1",
	"update ((select k1 from j1) d join j2 on d.k1 = j2.k1) set j2.b = 1",
	"delete j2 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)",
	// MariaDB also parenthesizes a lone derived table
	"select * from ((select k1 from j1) d)",
	"select * from (((select k1 from j1) d))",
	"select * from (((select k1 from j1) d) join j2 on 1)",
	"select * from (((select k1 from j1) d), j2)",
	"select * from (j2 join ((select 1) d) on 1)",
	// a parenthesized query stays a derived table
	"select * from ((select k1 from j1)) d",
	"select * from ((select k1 from j1) union (select k1 from j2)) d",
	"select * from ((select k1 from j1) limit 1) d",
	// expressions, as MariaDB writes a view over a subquery operand
	"select ((select `n1`.`a` from `n1` limit 1) + 0) AS `x`",
	"select ((select 1) + 0)",
	"select ((select a from n1 limit 1) + 0) from n1",
	"select ((select 1) = 1)",
	"select ((select 1) in (1))",
	"select ((select 1) not in (1, 2))",
	"select ((select 1) is null)",
	"select ((select 1) is not true)",
	"select ((select 1) between 0 and 2)",
	"select ((select 1) like 'a')",
	"select ((select 1) and 1)",
	"select ((select 1) <=> 1)",
	"select ((select 1) * 2 + 1)",
	"select ((select 1) - (select 2))",
	"select ((select 1) collate utf8mb4_bin)",
	"select ((select 1) = any (select 1))",
	"select ((select 1), 2) = (1, 2)",
	"select (((select 1) + 0))",
	"select (((select 1)) + 0)",
	"select ((((select 1) + 0) + 0) + 0)",
	"select (((select 1) limit 1) + 0)",
	"select ((with c as (select 1 x) select x from c) = 1)",
	"select abs(((select 1) + 0))",
	"select if(((select count(0) from n1) = 0), 'NO', 'YES')",
	"select coalesce(((select a from n1 limit 1) + 0), 1)",
	"select 1 in ((select 1) + 0, 2)",
	"select 1 in (((select 1) + 0), 2)",
	"select 1 in (((select 1) = 1))",
	"select * from n1 where a in ((select 1) + 0)",
	"select * from n1 where ((select a from n1 limit 1) = 1) and ((select 1) + 0 = 1)",
	"select * from n1 where (a, a) in (((select 1), 1))",
	// a parenthesized query stays a subquery
	"select ((select 1))",
	"select ((select 1)) + 0",
	"select ((select 1) union (select 1))",
	"select 1 in ((select 1) union (select 2))",
	"select ((select 1) union (select 2) limit 1)",
}

// parenSubqueryRejected are rejected by MariaDB 11.8 with 1064.
var parenSubqueryRejected = []string{
	"select * from ((select k1 from j1) d join j2 on 1",
	"select * from ((select k1 from j1) d join)",
	"select * from ((select k1 from j1) d d2 join j2 on 1)",
	"select * from ((select k1 from j1) d join j2 on 1) d2",
	"select * from ((select 1) union (select 2) d join j2 on 1)",
	"select * from ((select 1) d union (select 2))",
	"select ((select 1) + 0",
	"select ((select 1) +)",
	"select ((select 1) + 0 from n1)",
	"select ((select 1) d)",
	"select ((select 1) union (select 2) = 1)",
	"select (((select 1) + 0) union (select 1))",
	"select exists ((select 1) = 1)",
}

func TestParenSubqueryGroup(t *testing.T) {
	for _, sql := range parenSubqueryAccepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
	for _, sql := range parenSubqueryRejected {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestParenSubqueryGroupTree pins the trees of the two reported shapes.
func TestParenSubqueryGroupTree(t *testing.T) {
	sql := "select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)"
	sel := ParseAndCheck(t, sql).Items[0].(*nodes.SelectStmt)
	join, ok := sel.From[0].(*nodes.JoinClause)
	if !ok {
		t.Fatalf("FROM is %T, want *ast.JoinClause", sel.From[0])
	}
	sub, ok := join.Left.(*nodes.SubqueryExpr)
	if !ok || sub.Alias != "d" || sql[sub.Loc.Start:sub.Loc.End] != "(select j1.k1 from j1) d" {
		t.Errorf("join left operand = %#v, want the derived table d", join.Left)
	}

	sql = "select ((select a from n1 limit 1) + 0) from n1"
	sel = ParseAndCheck(t, sql).Items[0].(*nodes.SelectStmt)
	paren, ok := sel.TargetList[0].(*nodes.ParenExpr)
	if !ok {
		t.Fatalf("select item is %T, want *ast.ParenExpr", sel.TargetList[0])
	}
	bin, ok := paren.Expr.(*nodes.BinaryExpr)
	if !ok {
		t.Fatalf("parenthesized expression is %T, want *ast.BinaryExpr", paren.Expr)
	}
	if sub, ok := bin.Left.(*nodes.SubqueryExpr); !ok || sql[sub.Loc.Start:sub.Loc.End] != "(select a from n1 limit 1)" {
		t.Errorf("left operand = %#v, want the subquery", bin.Left)
	}
}

// parenSubqueryNesting nest the shapes that classify a '(' by the block that
// opens it.
var parenSubqueryNesting = []struct {
	name string
	sql  func(n int) string
}{
	{"subquery operands", func(n int) string { return "SELECT " + nest(n, "((SELECT 1) + ", "1", ")") }},
	{"left-deep subquery operands", func(n int) string { return "SELECT " + nest(n, "(", "(SELECT 1) = 0", ") = 0") }},
	{"subquery operand in parentheses", func(n int) string { return "SELECT " + nest(n, "(", "(SELECT 1) = 0", ")") }},
	{"subquery operand in an IN list", func(n int) string { return "SELECT " + nest(n, "1 IN ((SELECT 1) + ", "1", ")") }},
	{"derived-first join groups", func(n int) string {
		return "SELECT * FROM " + nest(n, "(", "(SELECT 1) d JOIN t ON 1", ") JOIN t ON 1")
	}},
	{"right-deep derived-first join groups", func(n int) string {
		return "SELECT * FROM " + nest(n, "((SELECT 1) d JOIN ", "t", " ON 1)")
	}},
	{"derived tables over derived-first join groups", func(n int) string {
		return nest(n, "SELECT * FROM ((", "SELECT 1", ") d JOIN t ON 1)")
	}},
	{"unclosed subquery operands", func(n int) string { return "SELECT " + strings.Repeat("((SELECT 1) + ", n) }},
	{"unclosed derived-first join groups", func(n int) string {
		return "SELECT * FROM " + strings.Repeat("((SELECT 1) d JOIN ", n)
	}},
}

// TestParenSubqueryGroupStaysLinear: deciding what a '(' opens scans the
// block after it, and the verdicts and the scan are recorded, so asking again
// at every level as the parser descends costs nothing.
func TestParenSubqueryGroupStaysLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	for _, tc := range parenSubqueryNesting {
		start := time.Now()
		_, err := Parse(tc.sql(depth))
		if unclosed := strings.HasPrefix(tc.name, "unclosed"); (err != nil) != unclosed {
			t.Fatalf("%s: Parse error = %v, want error %v", tc.name, briefErr(err), unclosed)
		}
		if elapsed := time.Since(start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse through %d levels took %v", tc.name, depth, elapsed)
		}
	}
}

// TestParenSubqueryGroupNestingLimit: the same shapes stop at maxNestingDepth
// with the nesting error, and parse 1,000 deep.
func TestParenSubqueryGroupNestingLimit(t *testing.T) {
	for _, tc := range parenSubqueryNesting {
		if err := parseWithin(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, briefErr(err))
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			err := parseWithin(tc.sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, briefErr(err))
			}
		}
	}
}
