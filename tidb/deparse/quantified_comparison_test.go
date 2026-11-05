package deparse

import "testing"

// TestQuantifiedComparison pins that a comparison keeps its quantifier through the
// rewrite and the writer, in TiDB v8.5.5's meaning, in this package's MySQL-shaped text, stored form. Dropping ALL reads `a > ALL
// (q)` as a scalar comparison against q, and folding NOT into the operator reads `NOT
// a = ANY (q)` as `a <> ANY (q)`; both parse and print without complaint, and a
// view stored through either means something else.
func TestQuantifiedComparison(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"a = ANY (SELECT 1)", "(`a` = any (select 1))"},
		{"a = SOME (SELECT 1)", "(`a` = any (select 1))"},
		{"a = ALL (SELECT 1)", "(`a` = all (select 1))"},
		{"a <> ANY (SELECT 1)", "(`a` <> any (select 1))"},
		{"a <> SOME (SELECT 1)", "(`a` <> any (select 1))"},
		{"a <> ALL (SELECT 1)", "(`a` <> all (select 1))"},
		{"a > ANY (SELECT 1)", "(`a` > any (select 1))"},
		{"a > SOME (SELECT 1)", "(`a` > any (select 1))"},
		{"a > ALL (SELECT 1)", "(`a` > all (select 1))"},
		{"a <= ALL (SELECT 1)", "(`a` <= all (select 1))"},
		{"NOT a = ANY (SELECT 1)", "(not((`a` = any (select 1))))"},
		{"NOT a > ALL (SELECT 1)", "(not((`a` > all (select 1))))"},
	} {
		t.Run(c.input, func(t *testing.T) {
			node := parseExpr(t, c.input)
			rewritten := RewriteExpr(node)
			if got := Deparse(rewritten); got != c.want {
				t.Errorf("Deparse(RewriteExpr(%q)) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
