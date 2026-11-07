package deparse

import "testing"

// TestQuantifiedComparison pins that a comparison keeps its quantifier through the
// rewrite and the writer, in MariaDB 11.8's stored form. Dropping ALL reads `a > ALL
// (q)` as a scalar comparison against q, and inverting a NOT without swapping the quantifier reads
// `NOT a = ANY (q)` as `a <> ANY (q)`; both parse and print without complaint, and a
// view stored through either means something else.
func TestQuantifiedComparison(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"a = ANY (SELECT 1)", "`a` in (select 1)"},
		{"a = SOME (SELECT 1)", "`a` in (select 1)"},
		{"a = ALL (SELECT 1)", "`a` = all (select 1)"},
		{"a <> ANY (SELECT 1)", "`a` <> any (select 1)"},
		{"a <> SOME (SELECT 1)", "`a` <> any (select 1)"},
		{"a <> ALL (SELECT 1)", "!(`a` in (select 1))"},
		{"a > ANY (SELECT 1)", "`a` > any (select 1)"},
		{"a > SOME (SELECT 1)", "`a` > any (select 1)"},
		{"a > ALL (SELECT 1)", "`a` > all (select 1)"},
		{"a <= ALL (SELECT 1)", "`a` <= all (select 1)"},
		{"NOT a = ANY (SELECT 1)", "!(`a` in (select 1))"},
		{"NOT a > ALL (SELECT 1)", "`a` <= any (select 1)"},
		{"NOT a > ANY (SELECT 1)", "`a` <= all (select 1)"},
		{"NOT a = ALL (SELECT 1)", "`a` <> any (select 1)"},
		{"NOT a <> ANY (SELECT 1)", "`a` = all (select 1)"},
		{"NOT a <> ALL (SELECT 1)", "`a` in (select 1)"},
		{"NOT a < ALL (SELECT 1)", "`a` >= any (select 1)"},
		{"NOT a <= ANY (SELECT 1)", "`a` > all (select 1)"},
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

// TestNestedQuantifiedComparison pins that a nested quantified comparison still means
// what was written: its text parses, and printing that parse again gives the same text.
// These writers leave `a > ALL (q)` and `a IN (q)` bare where the grammar allows it, so
// a missing pair of parentheses is either a syntax error or, as in `1 = (a > ALL (q))`
// read back as `(1 = a) > ALL (q)`, a view that means something else.
func TestNestedQuantifiedComparison(t *testing.T) {
	for _, input := range nestedQuantifiedShapes {
		t.Run(input, func(t *testing.T) {
			once := Deparse(RewriteExpr(parseExpr(t, input)))
			if twice := Deparse(RewriteExpr(parseExpr(t, once))); twice != once {
				t.Errorf("Deparse(%q) = %q, which reads back as %q", input, once, twice)
			}
		})
	}
}

// TestNestedQuantifiedComparisonParentheses pins where the parentheses go: only where
// the grammar needs them, so `(a > ALL (q)) = 1` stays bare as the server stores it.
func TestNestedQuantifiedComparisonParentheses(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"(a > ALL (SELECT 1)) + 1", "(`a` > all (select 1)) + 1"},
		{"1 + (a > ALL (SELECT 1))", "1 + (`a` > all (select 1))"},
		{"1 = (a > ALL (SELECT 1))", "1 = (`a` > all (select 1))"},
		{"(a > ALL (SELECT 1)) = 1", "`a` > all (select 1) = 1"},
		{"(a = ANY (SELECT 1)) + 1", "(`a` in (select 1)) + 1"},
		{"1 = (a = ANY (SELECT 1))", "1 = `a` in (select 1)"},
		{"(a = ANY (SELECT 1)) IN (SELECT 1)", "(`a` in (select 1)) in (select 1)"},
	} {
		t.Run(c.input, func(t *testing.T) {
			if got := Deparse(RewriteExpr(parseExpr(t, c.input))); got != c.want {
				t.Errorf("Deparse(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

var nestedQuantifiedShapes = []string{
	"(a > ALL (SELECT 1)) + 1",
	"1 + (a > ALL (SELECT 1))",
	"1 = (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) = 1",
	"(a = ANY (SELECT 1)) + 1",
	"1 = (a = ANY (SELECT 1))",
	"(a <> ALL (SELECT 1)) + 1",
	"1 = (a <> ALL (SELECT 1))",
	"(a > ANY (SELECT 1)) AND b",
	"b OR (a > ALL (SELECT 1))",
	"(a = ANY (SELECT 1)) AND b",
	"(a > ALL (SELECT 1)) IS NULL",
	"(a = ANY (SELECT 1)) IS NULL",
	"(a > ALL (SELECT 1)) IN (1, 2)",
	"(a = ANY (SELECT 1)) IN (1, 2)",
	"(a > ALL (SELECT 1)) IN (SELECT 1)",
	"(a = ANY (SELECT 1)) IN (SELECT 1)",
	"(a > ALL (SELECT 1)) BETWEEN 1 AND 2",
	"1 BETWEEN (a > ALL (SELECT 1)) AND 2",
	"1 BETWEEN 0 AND (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) LIKE 'x'",
	"'x' LIKE (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) COLLATE utf8mb4_bin",
	"(a > ALL (SELECT 1)) * (b < ANY (SELECT 1))",
	"-(a > ALL (SELECT 1))",
	"+(a > ALL (SELECT 1)) + 1",
	"(a > ALL (SELECT 1)) > ALL (SELECT 1)",
	"(a = ANY (SELECT 1)) = ANY (SELECT 1)",
	"(a <> ALL (SELECT 1)) > ALL (SELECT 1)",
	"(a IN (SELECT 1)) + 1",
	"(a NOT IN (SELECT 1)) + 1",
	"1 = (a IN (SELECT 1))",
	"(a IN (SELECT 1)) IN (SELECT 1)",
	"(a IN (SELECT 1)) BETWEEN 1 AND 2",
	"(a IN (SELECT 1)) LIKE 'x'",
	"(a IN (SELECT 1)) COLLATE utf8mb4_bin",
	"NOT a IN (SELECT 1)",
	"abs(a > ALL (SELECT 1))",
	"CASE WHEN a > ALL (SELECT 1) THEN 1 ELSE (a > ANY (SELECT 1)) + 1 END",
}
