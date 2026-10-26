package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// joinNestingCases covers how joins nest in a FROM clause. The right operand
// of INNER, LEFT, RIGHT, and FULL JOIN is a table source with joins of its
// own and is followed by a required ON; the right operand of CROSS JOIN and
// APPLY is a single table source; a joined table may be parenthesized; a
// join hint sits between the join type and JOIN; and PIVOT and UNPIVOT wrap
// everything to their left. shape is the outermost JoinClause, PivotExpr, or
// UnpivotExpr of an accepted statement, printed by joinShape.
var joinNestingCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// Nested joins, parenthesized and not.
	{"paren_right", "SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) ON t.a = t1.a", true, "(t JOIN (t1 JOIN t2))"},
	{"bare_right", "SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t JOIN (t1 JOIN t2))"},
	{"paren_whole", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a)", true, "(t1 JOIN t2)"},
	{"paren_double", "SELECT * FROM ((t1 JOIN t2 ON t1.a = t2.a))", true, "(t1 JOIN t2)"},
	{"paren_left", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) JOIN t ON t.a = t1.a", true, "((t1 JOIN t2) JOIN t)"},
	{"paren_both", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) JOIN (t3 JOIN t4 ON t3.a = t4.a) ON t1.a = t3.a", true, "((t1 JOIN t2) JOIN (t3 JOIN t4))"},
	{"bare_right_mixed_types", "SELECT * FROM t LEFT JOIN t1 INNER JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t LEFT (t1 JOIN t2))"},
	{"bare_right_outer_types", "SELECT * FROM t LEFT OUTER JOIN t1 RIGHT JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t LEFT (t1 RIGHT t2))"},
	{"bare_right_full", "SELECT * FROM t FULL JOIN t1 FULL OUTER JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t FULL (t1 FULL t2))"},
	{"bare_three_levels", "SELECT * FROM t JOIN t1 JOIN t2 JOIN t3 ON t2.a = t3.a ON t1.a = t2.a ON t.a = t1.a", true, "(t JOIN (t1 JOIN (t2 JOIN t3)))"},
	{"sequential_left_deep", "SELECT * FROM t JOIN t1 ON t.a = t1.a JOIN t2 ON t1.a = t2.a", true, "((t JOIN t1) JOIN t2)"},
	{"nested_then_sequential", "SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a JOIN t3 ON t.a = t3.a", true, "((t JOIN (t1 JOIN t2)) JOIN t3)"},
	{"sequential_inside_right", "SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t2.a = t3.a ON t.a = t1.a", true, "(t JOIN ((t1 JOIN t2) JOIN t3))"},
	{"cross_join_inside_right", "SELECT * FROM t JOIN t1 CROSS JOIN t2 ON t.a = t1.a", true, "(t JOIN (t1 CROSS t2))"},
	{"cross_apply_inside_right", "SELECT * FROM t JOIN t1 CROSS APPLY (SELECT 1 AS x) ca ON t.a = t1.a", true, "(t JOIN (t1 CROSS APPLY derived AS ca))"},
	{"outer_apply_inside_right", "SELECT * FROM t LEFT JOIN t1 OUTER APPLY (SELECT 1 AS x) oa ON t.a = t1.a", true, "(t LEFT (t1 OUTER APPLY derived AS oa))"},
	{"cross_join_then_join", "SELECT * FROM t CROSS JOIN t1 JOIN t2 ON t1.a = t2.a", true, "((t CROSS t1) JOIN t2)"},
	{"cross_join_paren_right", "SELECT * FROM t CROSS JOIN (t1 JOIN t2 ON t1.a = t2.a)", true, "(t CROSS (t1 JOIN t2))"},
	{"paren_cross_join", "SELECT * FROM t JOIN (t1 CROSS JOIN t2) ON t.a = t1.a", true, "(t JOIN (t1 CROSS t2))"},
	{"paren_derived_table_first", "SELECT * FROM t JOIN ((SELECT 1 AS a) s JOIN t2 ON s.a = t2.a) ON t.a = s.a", true, "(t JOIN (derived AS s JOIN t2))"},
	{"paren_values_first", "SELECT * FROM ((VALUES (1)) v(a) JOIN t2 ON v.a = t2.a)", true, "(derived AS v JOIN t2)"},
	{"paren_after_comma", "SELECT * FROM t, (t1 JOIN t2 ON t1.a = t2.a)", true, "(t1 JOIN t2)"},
	{"paren_then_apply", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) CROSS APPLY (SELECT 1 AS x) ca", true, "((t1 JOIN t2) CROSS APPLY derived AS ca)"},
	{"paren_then_pivot", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t1 JOIN t2) PIVOT p)"},
	{"table_hint_inside_right", "SELECT * FROM t JOIN t1 WITH (NOLOCK) JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t JOIN (t1 JOIN t2))"},
	{"aliases_inside_right", "SELECT * FROM t AS x JOIN t1 y JOIN t2 AS z ON y.a = z.a ON x.a = y.a", true, "(t AS x JOIN (t1 AS y JOIN t2 AS z))"},
	{"update_from_paren", "UPDATE t SET a = 1 FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) ON t.a = t1.a", true, "(t JOIN (t1 JOIN t2))"},
	{"delete_from_bare", "DELETE t FROM t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t JOIN (t1 JOIN t2))"},
	{"merge_using_paren", "MERGE INTO dbo.t USING (t1 JOIN t2 ON t1.a = t2.a) ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "(t1 JOIN t2)"},
	{"merge_using_bare", "MERGE INTO dbo.t USING t1 JOIN t2 ON t1.a = t2.a ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "(t1 JOIN t2)"},
	{"derived_table_then_join", "SELECT * FROM (SELECT a FROM dbo.t) d JOIN t1 ON d.a = t1.a", true, "(derived AS d JOIN t1)"},
	{"values_then_join", "SELECT * FROM (VALUES (1)) v(a) JOIN t1 ON v.a = t1.a", true, "(derived AS v JOIN t1)"},

	// A parenthesized joined table holds a join and takes no alias.
	{"paren_alias_as", "SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) AS x ON t.a = t1.a", false, ""},
	{"paren_alias_bare", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) x", false, ""},
	{"paren_single_table", "SELECT * FROM (t1)", false, ""},
	{"paren_single_table_right", "SELECT * FROM t JOIN (t1) ON t.a = t1.a", false, ""},
	{"paren_aliased_derived_table", "SELECT * FROM ((SELECT 1 AS a) s)", false, ""},
	{"paren_empty", "SELECT * FROM ()", false, ""},
	{"paren_unclosed", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a", false, ""},

	// ON is required after a qualified join and forbidden after CROSS JOIN.
	{"missing_on", "SELECT * FROM t JOIN t1", false, ""},
	{"missing_on_left", "SELECT * FROM t LEFT JOIN t1", false, ""},
	{"missing_on_before_where", "SELECT * FROM t JOIN t1 WHERE 1 = 1", false, ""},
	{"missing_outer_on", "SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a", false, ""},
	{"missing_on_after_paren", "SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a)", false, ""},
	{"extra_on", "SELECT * FROM t JOIN t1 ON t.a = t1.a ON t.a = t1.a", false, ""},
	{"cross_join_right_is_single", "SELECT * FROM t CROSS JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a", false, ""},
	{"cross_join_with_on", "SELECT * FROM t CROSS JOIN t1 ON t.a = t1.a", false, ""},
	{"missing_on_update_from", "UPDATE t SET a = 1 FROM t JOIN t1", false, ""},

	// The join type needs JOIN after it.
	{"inner_without_join", "SELECT * FROM t INNER t1 ON t.a = t1.a", false, ""},
	{"left_outer_without_join", "SELECT * FROM t LEFT OUTER t1 ON t.a = t1.a", false, ""},
	{"cross_without_join", "SELECT * FROM t CROSS t1", false, ""},
	{"hint_without_join", "SELECT * FROM t INNER HASH t1 ON t.a = t1.a", false, ""},

	// Join hints.
	{"hint_inner_hash", "SELECT * FROM t INNER HASH JOIN t1 ON t.a = t1.a", true, "(t JOIN HASH t1)"},
	{"hint_inner_loop", "SELECT * FROM t INNER LOOP JOIN t1 ON t.a = t1.a", true, "(t JOIN LOOP t1)"},
	{"hint_inner_merge", "SELECT * FROM t INNER MERGE JOIN t1 ON t.a = t1.a", true, "(t JOIN MERGE t1)"},
	{"hint_inner_remote", "SELECT * FROM t INNER REMOTE JOIN t1 ON t.a = t1.a", true, "(t JOIN REMOTE t1)"},
	{"hint_left_loop", "SELECT * FROM t LEFT LOOP JOIN t1 ON t.a = t1.a", true, "(t LEFT LOOP t1)"},
	{"hint_left_outer_merge", "SELECT * FROM t LEFT OUTER MERGE JOIN t1 ON t.a = t1.a", true, "(t LEFT MERGE t1)"},
	{"hint_right_hash", "SELECT * FROM t RIGHT HASH JOIN t1 ON t.a = t1.a", true, "(t RIGHT HASH t1)"},
	{"hint_full_outer_hash", "SELECT * FROM t FULL OUTER HASH JOIN t1 ON t.a = t1.a", true, "(t FULL HASH t1)"},
	{"hint_inside_right", "SELECT * FROM t JOIN t1 INNER HASH JOIN t2 ON t1.a = t2.a ON t.a = t1.a", true, "(t JOIN (t1 JOIN HASH t2))"},
	{"hint_lower_case", "SELECT * FROM t inner hash join t1 ON t.a = t1.a", true, "(t JOIN HASH t1)"},
	{"hint_reduce", "SELECT * FROM t INNER REDUCE JOIN t1 ON t.a = t1.a", true, "(t JOIN REDUCE t1)"},
	{"hint_replicate", "SELECT * FROM t LEFT OUTER replicate JOIN t1 ON t.a = t1.a", true, "(t LEFT REPLICATE t1)"},
	{"hint_redistribute", "SELECT * FROM t RIGHT REDISTRIBUTE JOIN t1 ON t.a = t1.a", true, "(t RIGHT REDISTRIBUTE t1)"},
	{"hint_word_as_alias", "SELECT * FROM t HASH JOIN t1 ON t.a = t1.a", true, "(t AS HASH JOIN t1)"},
	{"hint_remote_left", "SELECT * FROM t LEFT REMOTE JOIN t1 ON t.a = t1.a", false, ""},
	{"hint_remote_full_outer", "SELECT * FROM t FULL OUTER REMOTE JOIN t1 ON t.a = t1.a", false, ""},
	{"hint_on_cross_join", "SELECT * FROM t CROSS HASH JOIN t1", false, ""},
	{"hint_unknown", "SELECT * FROM t INNER FOO JOIN t1 ON t.a = t1.a", false, ""},
	{"hint_bracketed", "SELECT * FROM t INNER [HASH] JOIN t1 ON t.a = t1.a", false, ""},
	{"hint_quoted", `SELECT * FROM t INNER "REDUCE" JOIN t1 ON t.a = t1.a`, false, ""},
	{"hint_with_argument", "SELECT * FROM t INNER REDISTRIBUTE(2) JOIN t1 ON t.a = t1.a", false, ""},
	{"hint_twice", "SELECT * FROM t INNER HASH LOOP JOIN t1 ON t.a = t1.a", false, ""},

	// PIVOT and UNPIVOT follow a table source like a join and wrap everything
	// to their left, CROSS JOIN and APPLY included: SQL Server binds no column
	// of t outside p in t CROSS JOIN t1 PIVOT (...) p. Inside a qualified
	// join's right operand they wrap only that operand.
	{"pivot_after_join", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t JOIN t1) PIVOT p)"},
	{"pivot_after_join_as", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) AS p", true, "((t JOIN t1) PIVOT p)"},
	{"unpivot_after_join", "SELECT * FROM t JOIN t1 ON t.a = t1.a UNPIVOT (v FOR c IN (a, col)) u", true, "((t JOIN t1) UNPIVOT u)"},
	{"unpivot_after_left_join", "SELECT * FROM t LEFT JOIN t1 ON t.a = t1.a UNPIVOT (v FOR c IN (a, col)) u", true, "((t LEFT t1) UNPIVOT u)"},
	{"pivot_after_hinted_join", "SELECT * FROM t INNER HASH JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t JOIN HASH t1) PIVOT p)"},
	{"pivot_after_join_chain", "SELECT * FROM t JOIN t1 ON t.a = t1.a JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "(((t JOIN t1) JOIN t2) PIVOT p)"},
	{"pivot_after_cross_join", "SELECT * FROM t CROSS JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t CROSS t1) PIVOT p)"},
	{"pivot_after_cross_apply", "SELECT * FROM t CROSS APPLY (SELECT 1 AS x) ca PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t CROSS APPLY derived AS ca) PIVOT p)"},
	{"unpivot_after_outer_apply", "SELECT * FROM t OUTER APPLY (SELECT 1 AS x) oa UNPIVOT (v FOR c IN (a, col)) u", true, "((t OUTER APPLY derived AS oa) UNPIVOT u)"},
	{"pivot_inside_right", "SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p ON t.a = p.a", true, "(t JOIN (t1 PIVOT p))"},
	{"pivot_after_join_inside_right", "SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p ON t.a = p.a", true, "(t JOIN ((t1 JOIN t2) PIVOT p))"},
	{"pivot_then_join_inside_right", "SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a ON t.a = p.a", true, "(t JOIN ((t1 PIVOT p) JOIN t2))"},
	{"unpivot_after_pivot_inside_right", "SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p UNPIVOT (v FOR c IN (a, col)) u ON t.a = u.a", true, "(t JOIN ((t1 PIVOT p) UNPIVOT u))"},
	{"pivot_after_cross_join_inside_right", "SELECT * FROM t JOIN t1 CROSS JOIN t2 PIVOT (MAX(col) FOR a IN ([1])) p ON t.a = p.a", true, "(t JOIN ((t1 CROSS t2) PIVOT p))"},
	{"join_after_pivot", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a", true, "(((t JOIN t1) PIVOT p) JOIN t2)"},
	{"cross_join_after_pivot", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p CROSS JOIN t2", true, "(((t JOIN t1) PIVOT p) CROSS t2)"},
	{"join_after_cross_join_pivot", "SELECT * FROM t CROSS JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a", true, "(((t CROSS t1) PIVOT p) JOIN t2)"},
	{"pivot_twice_after_join", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p PIVOT (MAX(col) FOR a IN ([2])) q", true, "(((t JOIN t1) PIVOT p) PIVOT q)"},
	{"unpivot_after_pivot_after_join", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p UNPIVOT (v FOR c IN (a, col)) u", true, "(((t JOIN t1) PIVOT p) UNPIVOT u)"},
	{"unpivot_after_join_after_pivot", "SELECT * FROM t1 PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a UNPIVOT (v FOR c IN (a, col)) u", true, "(((t1 PIVOT p) JOIN t2) UNPIVOT u)"},
	{"pivot_after_join_then_where", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p WHERE p.a = 1", true, "((t JOIN t1) PIVOT p)"},
	{"pivot_after_join_then_comma", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p, t2", true, "((t JOIN t1) PIVOT p)"},
	{"pivot_after_join_after_comma", "SELECT * FROM t, t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t1 JOIN t2) PIVOT p)"},
	{"paren_then_pivot_then_join", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) PIVOT (MAX(col) FOR a IN ([1])) p JOIN t3 ON p.a = t3.a", true, "(((t1 JOIN t2) PIVOT p) JOIN t3)"},
	{"pivot_after_openjson", "SELECT * FROM OPENJSON('[]') j PIVOT (MAX(col) FOR a IN ([1])) p", true, "(derived AS j PIVOT p)"},
	{"pivot_after_apply_openjson", "SELECT * FROM t JOIN t1 ON t.a = t1.a CROSS APPLY OPENJSON('[]') j PIVOT (MAX(col) FOR a IN ([1])) p", true, "(((t JOIN t1) CROSS APPLY derived AS j) PIVOT p)"},
	{"pivot_after_table_variable", "DECLARE @tv TABLE (a INT, col INT); SELECT * FROM @tv v PIVOT (MAX(col) FOR a IN ([1])) p", true, "(@tv AS v PIVOT p)"},
	{"pivot_after_cross_join_table_variable", "DECLARE @tv TABLE (a INT, col INT); SELECT * FROM t CROSS JOIN @tv v PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t CROSS @tv AS v) PIVOT p)"},
	{"update_from_pivot_after_join", "UPDATE t SET a = 1 FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t JOIN t1) PIVOT p)"},
	{"delete_from_pivot_after_join", "DELETE t FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p", true, "((t JOIN t1) PIVOT p)"},
	{"merge_using_pivot_after_join", "MERGE INTO dbo.t USING t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p ON dbo.t.a = p.a WHEN MATCHED THEN DELETE;", true, "((t1 JOIN t2) PIVOT p)"},

	// A parenthesized joined table may hold a PIVOT or UNPIVOT but must end
	// with a join.
	{"paren_pivot_then_join", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p JOIN t3 ON p.a = t3.a)", true, "(((t1 JOIN t2) PIVOT p) JOIN t3)"},
	{"paren_pivot_first", "SELECT * FROM (t1 PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a)", true, "((t1 PIVOT p) JOIN t2)"},
	{"paren_unpivot_first", "SELECT * FROM (t1 UNPIVOT (v FOR c IN (a, col)) u JOIN t2 ON u.a = t2.a)", true, "((t1 UNPIVOT u) JOIN t2)"},
	{"paren_pivot_inside_right", "SELECT * FROM (t1 JOIN t2 PIVOT (MAX(col) FOR a IN ([1])) p ON t1.a = p.a)", true, "(t1 JOIN (t2 PIVOT p))"},
	{"paren_cross_join_pivot_cross_join", "SELECT * FROM (t1 CROSS JOIN t2 PIVOT (MAX(col) FOR a IN ([1])) p CROSS JOIN t3)", true, "(((t1 CROSS t2) PIVOT p) CROSS t3)"},
	{"paren_paren_pivot_then_join", "SELECT * FROM ((t1 JOIN t2 ON t1.a = t2.a) PIVOT (MAX(col) FOR a IN ([1])) p JOIN t3 ON p.a = t3.a)", true, "(((t1 JOIN t2) PIVOT p) JOIN t3)"},
	{"paren_right_pivot_first", "SELECT * FROM t JOIN (t1 PIVOT (MAX(col) FOR a IN ([1])) p JOIN t2 ON p.a = t2.a) ON t.a = p.a", true, "(t JOIN ((t1 PIVOT p) JOIN t2))"},
	{"paren_ends_with_pivot", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p)", false, ""},
	{"paren_ends_with_unpivot", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a UNPIVOT (v FOR c IN (a, col)) u)", false, ""},
	{"paren_right_ends_with_pivot", "SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p) ON t.a = p.a", false, ""},
	{"paren_cross_join_ends_with_pivot", "SELECT * FROM (t1 CROSS JOIN t2 PIVOT (MAX(col) FOR a IN ([1])) p)", false, ""},
	{"paren_single_pivot", "SELECT * FROM (t1 PIVOT (MAX(col) FOR a IN ([1])) p)", false, ""},
	{"paren_paren_join_then_pivot", "SELECT * FROM ((t1 JOIN t2 ON t1.a = t2.a) PIVOT (MAX(col) FOR a IN ([1])) p)", false, ""},

	// A PIVOT after a join takes nothing a join or a table name would.
	{"pivot_after_join_then_on", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p ON t.a = t1.a", false, ""},
	{"pivot_inside_right_missing_on", "SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p", false, ""},
	{"pivot_after_join_table_hint", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p WITH (NOLOCK)", false, ""},
	{"pivot_after_join_column_list", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p (x)", false, ""},
	{"pivot_after_join_unfinished", "SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT", false, ""},
}

var joinTypeNames = map[ast.JoinType]string{
	ast.JoinInner:      "JOIN",
	ast.JoinLeft:       "LEFT",
	ast.JoinRight:      "RIGHT",
	ast.JoinFull:       "FULL",
	ast.JoinCross:      "CROSS",
	ast.JoinCrossApply: "CROSS APPLY",
	ast.JoinOuterApply: "OUTER APPLY",
}

// joinShape prints a table source as its join tree, one parenthesized group
// per JoinClause, PivotExpr, and UnpivotExpr.
func joinShape(n ast.Node) string {
	switch v := n.(type) {
	case *ast.JoinClause:
		op := joinTypeNames[v.Type]
		if v.Hint != "" {
			op += " " + v.Hint
		}
		return "(" + joinShape(v.Left) + " " + op + " " + joinShape(v.Right) + ")"
	case *ast.PivotExpr:
		return "(" + joinShape(v.Source) + " PIVOT " + v.Alias + ")"
	case *ast.UnpivotExpr:
		return "(" + joinShape(v.Source) + " UNPIVOT " + v.Alias + ")"
	case *ast.TableRef:
		if v.Alias != "" {
			return v.Object + " AS " + v.Alias
		}
		return v.Object
	case *ast.TableVarRef:
		return v.Name + " AS " + v.Alias
	case *ast.AliasedTableRef:
		return "derived AS " + v.Alias
	default:
		return fmt.Sprintf("%T", n)
	}
}

// outermostJoin returns the first JoinClause, PivotExpr, or UnpivotExpr in
// pre-order, which is the outermost one of the first table source that has
// one.
func outermostJoin(root ast.Node) ast.Node {
	var found ast.Node
	ast.Inspect(root, func(n ast.Node) bool {
		if found != nil {
			return false
		}
		switch n.(type) {
		case *ast.JoinClause, *ast.PivotExpr, *ast.UnpivotExpr:
			found = n
			return false
		}
		return true
	})
	return found
}

// TestJoinNesting checks that omni accepts the join forms SQL Server accepts,
// builds the join tree they denote, and keeps every Loc inside its parent.
func TestJoinNesting(t *testing.T) {
	for _, tc := range joinNestingCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error:\n%s", tc.sql, ast.NodeToString(list))
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			join := outermostJoin(list)
			if join == nil {
				t.Fatalf("no JoinClause, PivotExpr, or UnpivotExpr in %s", ast.NodeToString(list))
			}
			if got := joinShape(join); got != tc.shape {
				t.Errorf("join tree = %s, want %s", got, tc.shape)
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestJoinNestingOracle checks accept/reject parity with SQL Server for every
// case in joinNestingCases.
func TestJoinNestingOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range joinNestingCases {
		t.Run(tc.name, func(t *testing.T) {
			ssErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, tc.sql)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

// TestJoinNestingLoc checks the span of each JoinClause: a parenthesized
// joined table spans its parentheses, and a join that starts with one starts
// at its '('.
func TestJoinNestingLoc(t *testing.T) {
	tests := []struct {
		sql   string
		index int // JoinClause in pre-order
		want  string
	}{
		{"SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) ON t.a = t1.a", 0, "t JOIN (t1 JOIN t2 ON t1.a = t2.a) ON t.a = t1.a"},
		{"SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) ON t.a = t1.a", 1, "(t1 JOIN t2 ON t1.a = t2.a)"},
		{"SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a", 0, "t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a"},
		{"SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a ON t.a = t1.a", 1, "t1 JOIN t2 ON t1.a = t2.a"},
		{"SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) JOIN t ON t.a = t1.a", 0, "(t1 JOIN t2 ON t1.a = t2.a) JOIN t ON t.a = t1.a"},
		{"SELECT * FROM ((t1 JOIN t2 ON t1.a = t2.a))", 0, "((t1 JOIN t2 ON t1.a = t2.a))"},
		{"SELECT * FROM t LEFT OUTER HASH JOIN t1 ON t.a = t1.a", 0, "t LEFT OUTER HASH JOIN t1 ON t.a = t1.a"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			join := nthNodeByType(list, "JoinClause", tc.index)
			if join == nil {
				t.Fatalf("no JoinClause #%d in %s", tc.index, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(join)); got != tc.want {
				t.Errorf("JoinClause #%d text = %q, want %q", tc.index, got, tc.want)
			}
		})
	}
}

// TestJoinPivotLoc checks the span of a PIVOT or UNPIVOT among joins: it
// starts at its source, so it covers the joins it wraps, and a join that
// takes it as the left operand starts where it does.
func TestJoinPivotLoc(t *testing.T) {
	const pivot = "PIVOT (MAX(col) FOR a IN ([1])) p"
	tests := []struct {
		sql   string
		node  string
		index int // node of that type in pre-order
		want  string
	}{
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot, "PivotExpr", 0, "t JOIN t1 ON t.a = t1.a " + pivot},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot, "JoinClause", 0, "t JOIN t1 ON t.a = t1.a"},
		{"SELECT * FROM t LEFT JOIN t1 ON t.a = t1.a UNPIVOT (v FOR c IN (a, col)) u WHERE u.a = 1", "UnpivotExpr", 0, "t LEFT JOIN t1 ON t.a = t1.a UNPIVOT (v FOR c IN (a, col)) u"},
		{"SELECT * FROM t CROSS JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) AS p", "PivotExpr", 0, "t CROSS JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) AS p"},
		{"SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a " + pivot + " ON t.a = p.a", "JoinClause", 0, "t JOIN t1 JOIN t2 ON t1.a = t2.a " + pivot + " ON t.a = p.a"},
		{"SELECT * FROM t JOIN t1 JOIN t2 ON t1.a = t2.a " + pivot + " ON t.a = p.a", "PivotExpr", 0, "t1 JOIN t2 ON t1.a = t2.a " + pivot},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot + " JOIN t2 ON p.a = t2.a", "JoinClause", 0, "t JOIN t1 ON t.a = t1.a " + pivot + " JOIN t2 ON p.a = t2.a"},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot + " JOIN t2 ON p.a = t2.a", "PivotExpr", 0, "t JOIN t1 ON t.a = t1.a " + pivot},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot + " PIVOT (MAX(col) FOR a IN ([2])) q", "PivotExpr", 0, "t JOIN t1 ON t.a = t1.a " + pivot + " PIVOT (MAX(col) FOR a IN ([2])) q"},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a " + pivot + " PIVOT (MAX(col) FOR a IN ([2])) q", "PivotExpr", 1, "t JOIN t1 ON t.a = t1.a " + pivot},
		{"SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a " + pivot + " JOIN t3 ON p.a = t3.a)", "JoinClause", 0, "(t1 JOIN t2 ON t1.a = t2.a " + pivot + " JOIN t3 ON p.a = t3.a)"},
		{"SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a " + pivot + " JOIN t3 ON p.a = t3.a)", "PivotExpr", 0, "t1 JOIN t2 ON t1.a = t2.a " + pivot},
		{"SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) " + pivot, "PivotExpr", 0, "(t1 JOIN t2 ON t1.a = t2.a) " + pivot},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s#%d/%s", tc.node, tc.index, tc.want), func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			n := nthNodeByType(list, tc.node, tc.index)
			if n == nil {
				t.Fatalf("no %s #%d in %s", tc.node, tc.index, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s #%d text = %q, want %q", tc.node, tc.index, got, tc.want)
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestJoinNestingErrorPosition checks that a missing ON or JOIN is reported
// at the token that stands in its place.
func TestJoinNestingErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"SELECT * FROM t JOIN t1 WHERE 1 = 1", `syntax error at or near "WHERE"`},
		{"SELECT * FROM t JOIN t1", "syntax error at end of input"},
		{"SELECT * FROM t JOIN (t1 JOIN t2 ON t1.a = t2.a) x ON t.a = t1.a", `syntax error at or near "x"`},
		{"SELECT * FROM (t1)", `syntax error at or near ")"`},
		{"SELECT * FROM t INNER t1 ON t.a = t1.a", `syntax error at or near "t1"`},
		{"SELECT * FROM t CROSS t1", `syntax error at or near "t1"`},
		{"SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a PIVOT (MAX(col) FOR a IN ([1])) p)", `syntax error at or near ")"`},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p ON t.a = t1.a", `syntax error at or near "ON"`},
		{"SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p", "syntax error at end of input"},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want %s", tc.sql, tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Parse(%q) error = %q, want %s", tc.sql, err.Error(), tc.want)
			}
		})
	}
}

// TestJoinNestingCompletion checks the candidates offered inside a join
// operator and after the right operand of a qualified join.
func TestJoinNestingCompletion(t *testing.T) {
	tests := []struct {
		sql       string
		wantToks  []int
		wantRules []string
	}{
		{"SELECT * FROM t JOIN t1 ", []int{kwON, kwJOIN, kwCROSS}, nil},
		{"SELECT * FROM t JOIN t1 JOIN t2 ", []int{kwON, kwJOIN}, nil},
		{"SELECT * FROM t INNER ", []int{kwJOIN, kwLOOP, kwHASH, kwMERGE, kwREMOTE}, nil},
		{"SELECT * FROM t LEFT ", []int{kwOUTER, kwJOIN, kwHASH}, nil},
		{"SELECT * FROM t LEFT OUTER ", []int{kwJOIN, kwLOOP}, nil},
		{"SELECT * FROM t FULL HASH ", []int{kwJOIN}, nil},
		{"SELECT * FROM t CROSS ", []int{kwJOIN, kwAPPLY}, nil},
		{"SELECT * FROM (", []int{kwSELECT, kwVALUES}, []string{"table_ref"}},
		{"SELECT * FROM t JOIN (", []int{kwSELECT, kwVALUES}, []string{"table_ref"}},
		{"SELECT * FROM t JOIN (t1 JOIN t2 ON ", nil, []string{"columnref"}},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (", nil, []string{"func_name"}},
		{"SELECT * FROM t JOIN t1 ON t.a = t1.a PIVOT (MAX(col) FOR a IN ([1])) p ", []int{kwJOIN, kwCROSS, kwWHERE}, nil},
		{"SELECT * FROM t JOIN t1 PIVOT (MAX(col) FOR a IN ([1])) p ", []int{kwON, kwJOIN}, nil},
	}
	for _, tc := range tests {
		t.Run(tc.sql+"|", func(t *testing.T) {
			cs := Collect(tc.sql, len(tc.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, tok := range tc.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s", TokenName(tok))
				}
			}
			for _, r := range tc.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
		})
	}
	// ON is a candidate only where a qualified join still owes one.
	if cs := Collect("SELECT * FROM t ", len("SELECT * FROM t ")); cs.HasToken(kwON) {
		t.Errorf("SELECT * FROM t | offers ON")
	}
}
