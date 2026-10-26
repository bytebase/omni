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
// APPLY is a single table source; a joined table may be parenthesized; and a
// join hint sits between the join type and JOIN. shape is the outermost
// JoinClause of an accepted statement, printed by joinShape.
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
	{"paren_then_pivot", "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) PIVOT (MAX(col) FOR a IN ([1])) p", true, "(t1 JOIN t2)"},
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
// per JoinClause.
func joinShape(n ast.Node) string {
	switch v := n.(type) {
	case *ast.JoinClause:
		op := joinTypeNames[v.Type]
		if v.Hint != "" {
			op += " " + v.Hint
		}
		return "(" + joinShape(v.Left) + " " + op + " " + joinShape(v.Right) + ")"
	case *ast.TableRef:
		if v.Alias != "" {
			return v.Object + " AS " + v.Alias
		}
		return v.Object
	case *ast.AliasedTableRef:
		return "derived AS " + v.Alias
	default:
		return fmt.Sprintf("%T", n)
	}
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
			join := nthNodeByType(list, "JoinClause", 0)
			if join == nil {
				t.Fatalf("no JoinClause in %s", ast.NodeToString(list))
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
