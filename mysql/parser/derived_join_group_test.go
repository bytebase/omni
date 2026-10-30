package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// derivedJoinGroupAccepted are join groups and table lists in parentheses
// whose first operand is a derived table, which MySQL 8.0 accepts
// (TestDerivedJoinGroupOracle). MySQL writes the first one for a view over a
// derived table joined to a table.
var derivedJoinGroupAccepted = []string{
	"select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))",
	"select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)",
	// aliases
	"select * from ((select k1 from j1) as d join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d (x) join j2 on d.x = j2.k1)",
	"select * from ((select k1 from j1) as d (x) join j2 on d.x = j2.k1)",
	// no alias: MySQL parses it and then reports 1248, a missing alias
	"select * from ((select k1 from j1) join j2 on 1)",
	// join kinds
	"select * from ((select k1 from j1) d left join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d left outer join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d right join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d inner join j2 on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d cross join j2)",
	"select * from ((select k1 from j1) d join j2)",
	"select * from ((select k1 from j1) d join j2 using (k1))",
	"select * from ((select k1 from j1) d natural join j2)",
	"select * from ((select k1 from j1) d natural left join j2)",
	"select * from ((select k1 from j1) d straight_join j2)",
	"select * from ((select k1 from j1) d straight_join j2 on d.k1 = j2.k1)",
	// table lists
	"select * from ((select k1 from j1) d, j2)",
	"select * from ((select k1 from j1) d, (select k1 from j2) e)",
	"select * from (((select k1 from j1) d, j3) join j2 on 1)",
	// nested groups
	"select * from (((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1)",
	"select * from ((((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1))",
	"select * from ((select k1 from j1) d join (j2 join j3 on j2.k1 = j3.k1) on d.k1 = j2.k1)",
	"select * from ((select k1 from j1) d join ((select k1 from j2) e join j3 on e.k1 = j3.k1) on d.k1 = e.k1)",
	"select * from (j1 join ((select 1) d join j2 on 1) on 1)",
	"select * from j1 join ((select 1) d join j2 on 1) on 1",
	"select * from ((select 1) d join j2 on 1) join j3 on 1",
	"select * from ((select 1) d join j2 on 1), j3",
	// the derived table's query in more parentheses
	"select * from (((select k1 from j1)) d join j2 on d.k1 = j2.k1)",
	"select * from ((((select k1 from j1))) d join j2 on d.k1 = j2.k1)",
	// set operations, WITH, TABLE, VALUES
	"select * from ((select k1 from j1 union select k1 from j2) d join j3 on d.k1 = j3.k1)",
	"select * from (((select k1 from j1) union (select k1 from j2)) d join j3 on d.k1 = j3.k1)",
	"select * from (((select k1 from j1) union (select k1 from j2) order by 1 limit 1) d join j3 on d.k1 = j3.k1)",
	"select * from (((select k1 from j1) except (select k1 from j2)) d join j3 on 1)",
	"select * from (((select k1 from j1) intersect (select k1 from j2)) d join j3 on 1)",
	"select * from ((with c as (select 1 k1) select k1 from c) d join j2 on 1)",
	"select * from ((table j1) d join j2 on 1)",
	"select * from ((values row(1)) d join j2 on 1)",
	// LATERAL
	"select * from (lateral (select 1) d join j2 on 1)",
	"select * from (j1 join lateral (select j1.k1) d on 1)",
	"select * from (j1, lateral (select j1.k1) d)",
	// other statements
	"update ((select k1 from j1) d join j2 on d.k1 = j2.k1) set j2.b = 1",
	"delete j2 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)",
	"insert into j3 select d.k1, 1 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)",
	"with c as (select 1 k1) select * from ((select k1 from c) d join j2 on 1)",
	"select (select 1 from ((select 1) d join j2 on 1) limit 1)",
	"select 1 from j1 where exists (select 1 from ((select 1) d join j2 on 1))",
	// a parenthesized query stays a derived table
	"select * from ((select k1 from j1)) d",
	"select * from (((select k1 from j1))) d",
	"select * from ((select k1 from j1) union (select k1 from j2)) d",
	"select * from ((select k1 from j1) union select k1 from j2) d",
	"select * from ((select k1 from j1) order by 1 limit 1) d",
	"select * from ((select 1) for update) d",
	"select * from ((select 1) lock in share mode) d",
}

// derivedJoinGroupRejected are rejected by MySQL 8.0 with 1064
// (TestDerivedJoinGroupOracle). MySQL does not parenthesize a lone derived
// table or table function.
var derivedJoinGroupRejected = []string{
	"select * from ((select k1 from j1) d)",
	"select * from (((select k1 from j1) d))",
	"select * from ((select k1 from j1) as d)",
	"select * from ((select k1 from j1) d) join j2 on 1",
	"select * from (((select 1) d) join j2 on 1)",
	"select * from (((select 1) d), j2)",
	"select * from j1, (lateral (select j1.k1) d)",
	"select * from (json_table('[1]', '$[*]' columns (x int path '$')) jt)",
	"select * from ((select k1 from j1) d join j2 on 1",
	"select * from ((select k1 from j1) d join)",
	"select * from ((select k1 from j1) d d2 join j2 on 1)",
	"select * from ((select k1 from j1) d join j2 on 1) d2",
	"select * from ((select k1 from j1) d (x join j2 on 1)",
	"select * from ((select k1 from j1) d use index (i) join j2 on 1)",
	"select * from ((select 1) union (select 2) d join j2 on 1)",
	"select * from ((select 1) order by 1 d join j2 on 1)",
}

func TestDerivedJoinGroup(t *testing.T) {
	for _, sql := range derivedJoinGroupAccepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
	for _, sql := range derivedJoinGroupRejected {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestDerivedJoinGroupTree pins the tree of the reported statement: the join
// group's first operand is the derived table, with its alias and its own
// span.
func TestDerivedJoinGroupTree(t *testing.T) {
	sql := "select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)"
	list := ParseAndCheck(t, sql)
	sel := list.Items[0].(*nodes.SelectStmt)
	if len(sel.From) != 1 {
		t.Fatalf("FROM has %d items, want 1", len(sel.From))
	}
	join, ok := sel.From[0].(*nodes.JoinClause)
	if !ok {
		t.Fatalf("FROM is %T, want *ast.JoinClause", sel.From[0])
	}
	sub, ok := join.Left.(*nodes.SubqueryExpr)
	if !ok {
		t.Fatalf("join left operand is %T, want *ast.SubqueryExpr", join.Left)
	}
	if got, want := sql[sub.Loc.Start:sub.Loc.End], "(select j1.k1 from j1) d"; got != want {
		t.Errorf("derived table spans %q, want %q", got, want)
	}
	if sub.Alias != "d" || sub.Select == nil || len(sub.Select.From) != 1 {
		t.Errorf("derived table = alias %q, query %+v", sub.Alias, sub.Select)
	}
	if right, ok := join.Right.(*nodes.TableRef); !ok || right.Name != "j2" {
		t.Errorf("join right operand is %#v, want table j2", join.Right)
	}
	if _, ok := join.Condition.(*nodes.OnCondition); !ok {
		t.Errorf("join condition is %T, want *ast.OnCondition", join.Condition)
	}
}

// TestLoneParenthesizedDerivedTableError: MySQL rejects a lone derived table
// or table function in parentheses at the ')' (8.0.46: "near ')'").
func TestLoneParenthesizedDerivedTableError(t *testing.T) {
	for _, sql := range []string{
		"select * from ((select 1) d)",
		"select * from (lateral (select 1) d)",
		"select * from (json_table('[1]', '$[*]' columns (x int path '$')) jt)",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != len(sql)-1 {
			t.Errorf("Parse(%q) = %v, want a syntax error at the last ')' (offset %d)", sql, err, len(sql)-1)
		}
	}
}

// derivedJoinGroupNesting nest the shapes that classify a '(' by the block
// that opens it.
var derivedJoinGroupNesting = []struct {
	name string
	sql  func(n int) string
}{
	{"derived-first join groups", func(n int) string {
		return "SELECT * FROM " + nest(n, "(", "(SELECT 1) d JOIN t ON 1", ") JOIN t ON 1")
	}},
	{"right-deep derived-first join groups", func(n int) string {
		return "SELECT * FROM " + nest(n, "((SELECT 1) d JOIN ", "t", " ON 1)")
	}},
	{"derived tables over derived-first join groups", func(n int) string {
		return nest(n, "SELECT * FROM ((", "SELECT 1", ") d JOIN t ON 1)")
	}},
	{"derived table in parentheses in a join group", func(n int) string {
		return "SELECT * FROM (" + nest(n, "(", "SELECT 1", ")") + " d JOIN t ON 1)"
	}},
	{"unclosed derived-first join groups", func(n int) string {
		return "SELECT * FROM " + strings.Repeat("((SELECT 1) d JOIN ", n)
	}},
}

// TestDerivedJoinGroupStaysLinear: deciding what a '(' in a FROM clause opens
// scans the block after it, and the verdicts and the scan are recorded, so
// asking again at every level as the parser descends costs nothing.
func TestDerivedJoinGroupStaysLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	for _, tc := range derivedJoinGroupNesting {
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

// TestDerivedJoinGroupNestingLimit: the same shapes stop at maxNestingDepth
// with the nesting error, and parse 1,000 deep.
func TestDerivedJoinGroupNestingLimit(t *testing.T) {
	for _, tc := range derivedJoinGroupNesting {
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
