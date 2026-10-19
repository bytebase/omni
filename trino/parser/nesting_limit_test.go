package parser

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

// raiseNestingLimit lifts maxNestingDepth, and maxTreeHeight with it, for a
// test that nests deeper on purpose, and restores them when the test ends.
func raiseNestingLimit(t *testing.T) {
	setNestingLimit(t, 1<<20)
	saved := maxTreeHeight
	maxTreeHeight = 1 << 30
	t.Cleanup(func() { maxTreeHeight = saved })
}

// setNestingLimit sets maxNestingDepth for the rest of the test.
func setNestingLimit(t *testing.T, n int) {
	saved := maxNestingDepth
	maxNestingDepth = n
	t.Cleanup(func() { maxNestingDepth = saved })
}

// isNestingError reports whether err carries the nesting error.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == nestingTooDeep
}

// parseWithin runs Parse and ParseBestEffort on sql and returns Parse's
// error, failing the test when they take longer than the deadline: a
// speculative parse that retried its readings at the limit would never
// return instead of failing.
func parseWithin(t *testing.T, name string, sql string) error {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		ParseBestEffort(sql)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("%s: Parse and ParseBestEffort of %d bytes took over 10s", name, len(sql))
		return nil
	}
}

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

func routine(body string) string {
	return "CREATE FUNCTION f() RETURNS int BEGIN " + body + " END"
}

// nestingShapes are the constructs the parser descends into recursively, n
// levels deep.
var nestingShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
	{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
	{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true" }},
	{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
	{"function call", func(n int) string { return "SELECT " + nest(n, "f(", "1", ")") }},
	{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS int)") }},
	{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN true THEN ", "1", " END") }},
	{"ARRAY", func(n int) string { return "SELECT " + nest(n, "ARRAY[", "1", "]") }},
	{"ROW", func(n int) string { return "SELECT " + nest(n, "ROW(", "1", ")") }},
	{"subscript", func(n int) string { return "SELECT " + nest(n, "a[", "1", "]") }},
	{"AT TIME ZONE", func(n int) string { return "SELECT " + nest(n, "(", "a", " AT TIME ZONE 'UTC')") }},
	{"lambda", func(n int) string { return "SELECT " + nest(n, "transform(a, x -> ", "x", ")") }},
	{"window", func(n int) string { return "SELECT " + nest(n, "count(*) OVER (ORDER BY ", "1", ")") }},
	{"JSON_ARRAY", func(n int) string { return "SELECT " + nest(n, "JSON_ARRAY(", "1", ")") }},
	{"JSON_OBJECT", func(n int) string { return "SELECT " + nest(n, "JSON_OBJECT('k' VALUE ", "1", ")") }},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
	{"EXISTS", func(n int) string { return nest(n, "SELECT 1 WHERE EXISTS (", "SELECT 1", ")") }},
	{"IN subquery", func(n int) string { return nest(n, "SELECT 1 WHERE 1 IN (", "SELECT 1", ")") }},
	{"quantified comparison", func(n int) string { return nest(n, "SELECT 1 WHERE 1 = ANY (", "SELECT 1", ")") }},
	{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
	{"LATERAL", func(n int) string { return nest(n, "SELECT * FROM LATERAL (", "SELECT 1", ")") }},
	{"table function argument", func(n int) string {
		return nest(n, "SELECT * FROM TABLE(f(input => TABLE(", "SELECT 1", ")))")
	}},
	{"FROM parentheses", func(n int) string { return "SELECT * FROM " + nest(n, "(", "t", ")") }},
	{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON true)") }},
	{"right-deep join", func(n int) string {
		return "SELECT * FROM " + strings.Repeat("t JOIN ", n) + "t" + strings.Repeat(" ON true", n)
	}},
	{"parenthesized query", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
	{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
	{"INTERSECT operand", func(n int) string { return nest(n, "SELECT 1 INTERSECT (", "SELECT 1", ")") }},
	{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
	{"array type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "array(", "int", ")") + ")" }},
	{"legacy array type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "array<", "int", ">") + ")" }},
	{"map type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "map(int, ", "int", ")") + ")" }},
	{"ROW type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "ROW(a ", "int", ")") + ")" }},
	{"row pattern", func(n int) string {
		return "SELECT * FROM t MATCH_RECOGNIZE (PATTERN (" + nest(n, "(", "A", ")") + ") DEFINE A AS true)"
	}},
	{"PERMUTE", func(n int) string {
		return "SELECT * FROM t MATCH_RECOGNIZE (PATTERN (" + nest(n, "PERMUTE(", "A", ")") + ") DEFINE A AS true)"
	}},
	{"compound statement", func(n int) string { return routine(nest(n, "BEGIN ", "RETURN 1;", " END;")) }},
	{"IF statement", func(n int) string { return routine(nest(n, "IF true THEN ", "RETURN 1;", " END IF;")) }},
	{"LOOP statement", func(n int) string { return routine(nest(n, "LOOP ", "RETURN 1;", " END LOOP;")) }},
	{"EXPLAIN", func(n int) string { return strings.Repeat("EXPLAIN ", n) + "SELECT 1" }},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
	{"unclosed function calls", func(n int) string { return "SELECT " + strings.Repeat("f(", n) }},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
	{"unclosed EXISTS", func(n int) string { return "SELECT " + strings.Repeat("EXISTS (SELECT (", n) }},
	{"unclosed derived tables", func(n int) string { return strings.Repeat("SELECT * FROM (", n) }},
	{"unclosed FROM parentheses", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
	{"unclosed types", func(n int) string { return "SELECT CAST(NULL AS " + strings.Repeat("array(", n) }},
	{"unclosed routine blocks", func(n int) string { return routine(strings.Repeat("BEGIN ", n)) }},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: 384KB of unclosed '(' were enough. The same
// constructs nested 1,000 deep still parse; Trino 482 itself rejects most of
// them at fewer levels than that.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		if !strings.HasPrefix(tc.name, "unclosed") {
			if err := parseWithin(t, tc.name, tc.sql(1000)); err != nil {
				t.Errorf("%s: 1,000 levels: %v", tc.name, err)
			}
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			if err := parseWithin(t, tc.name, tc.sql(n)); !isNestingError(err) {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// TestNestingLimitStandaloneEntryPoints: ParseExpression and ParseDataType
// stop at the limit too.
func TestNestingLimitStandaloneEntryPoints(t *testing.T) {
	const n = 200000
	if _, errs := ParseExpression(strings.Repeat("(", n) + "1" + strings.Repeat(")", n)); len(errs) != 1 || errs[0].Message != nestingTooDeep {
		t.Errorf("ParseExpression: errs = %v, want the nesting error", errs)
	}
	if _, errs := ParseDataType(strings.Repeat("array(", n) + "int" + strings.Repeat(")", n)); len(errs) != 1 || errs[0].Message != nestingTooDeep {
		t.Errorf("ParseDataType: errs = %v, want the nesting error", errs)
	}
}

// TestNestingLimitCountsReparsedBodies: a PREPARE's statement, and a subquery
// body the parse captured as raw text, are parsed again by a parser of their
// own, and the bodies nested inside are parsed again by theirs. Each such
// parser starts at the depth of the body it parses, so the limit counts the
// levels around the body and the reparses stop at it; started at depth 0,
// they recursed into each other without bound. Every level reparses the text
// of all the levels inside it, so the limit is lowered here to keep the
// reparses short.
func TestNestingLimitCountsReparsedBodies(t *testing.T) {
	setNestingLimit(t, 100)
	for _, tc := range []struct {
		name string
		sql  func(n int) string
		ok   bool // whether the statement parses when shallow
	}{
		{"PREPARE", func(n int) string { return strings.Repeat("PREPARE p FROM ", n) + "SELECT 1" }, true},
		// The in-place subquery parse rejects WITH FUNCTION, so every body is
		// captured raw and parsed again as a statement.
		{"WITH FUNCTION subquery", func(n int) string {
			return "SELECT " + strings.Repeat("(WITH FUNCTION f() RETURNS int RETURN 1 SELECT ", n) + "1" + strings.Repeat(")", n)
		}, false},
	} {
		err := parseWithin(t, tc.name, tc.sql(10))
		if tc.ok && err != nil || !tc.ok && (err == nil || isNestingError(err)) {
			t.Errorf("%s: 10 levels: err = %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 300} {
			if err := parseWithin(t, tc.name, tc.sql(n)); !isNestingError(err) {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// statementHeight parses sql, one statement, and returns how many nodes deep
// its tree is, counting the bodies of its subqueries. A trailing comment
// lengthens the text, without adding to the tree, past the length below
// which the parser does not record the bodies' heights.
func statementHeight(t *testing.T, sql string) int {
	t.Helper()
	sql += " /*" + strings.Repeat(" ", maxTreeHeight/16) + "*/"
	p := &Parser{lexer: NewLexer(sql), input: sql, strictTrailing: true}
	p.advance()
	node, err := p.parseStmt()
	if err != nil {
		t.Fatalf("%.60s...: %v", sql, err)
	}
	return treeHeight(node, math.MaxInt-1)
}

// chainShapes build left-deep chains of n terms: the parser builds each with
// a loop, so only the tree, not the parse, gets deep.
var chainShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1+", n) + "1" }},
	{"difference", func(n int) string { return "SELECT " + strings.Repeat("1-", n) + "1" }},
	{"product", func(n int) string { return "SELECT " + strings.Repeat("1*", n) + "1" }},
	{"concatenation", func(n int) string { return "SELECT " + strings.Repeat("'a'||", n) + "'a'" }},
	{"OR", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a OR ", n) + "a" }},
	{"AND", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a AND ", n) + "a" }},
	{"AT TIME ZONE", func(n int) string { return "SELECT a" + strings.Repeat(" AT TIME ZONE 'UTC'", n) }},
	{"subscripts", func(n int) string { return "SELECT a" + strings.Repeat("[1]", n) }},
	{"field references", func(n int) string { return "SELECT f(a)" + strings.Repeat(".b", n) }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT 1 UNION ALL ", n) + "SELECT 1" }},
	{"EXCEPT", func(n int) string { return strings.Repeat("SELECT 1 EXCEPT ", n) + "SELECT 1" }},
	{"INTERSECT", func(n int) string { return strings.Repeat("SELECT 1 INTERSECT ", n) + "SELECT 1" }},
	{"joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON true", n) }},
	{"cross joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" CROSS JOIN t", n) }},
	{"natural joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" NATURAL JOIN t", n) }},
	{"array type suffixes", func(n int) string { return "SELECT CAST(NULL AS int" + strings.Repeat(" ARRAY", n) + ")" }},
	{"row pattern alternation", func(n int) string {
		return "SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A" + strings.Repeat(" | A", n) + ") DEFINE A AS true)"
	}},
	{"row pattern concatenation", func(n int) string {
		return "SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A" + strings.Repeat(" A", n) + ") DEFINE A AS true)"
	}},
	// A chain in each of four nested subqueries, the subquery at the bottom
	// of the chain around it: trino/analysis parses each body again and
	// walks its tree from there, so the heights add up.
	{"chains in nested subqueries", func(n int) string {
		const levels = 4
		chain := strings.Repeat("+1", n/levels+1)
		return strings.Repeat("SELECT (", levels) + "SELECT 1" + strings.Repeat(")"+chain, levels)
	}},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it (trino/analysis did at three million terms). A
// chain within the limit still parses.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if err := parseWithin(t, tc.name, tc.sql(maxTreeHeight/4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		err := parseWithin(t, tc.name, tc.sql(maxTreeHeight+1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingTooDeep || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, err)
		}
		if tc.name == "chains in nested subqueries" {
			// Best effort does not parse subquery bodies: its tree holds
			// them as text, and nothing deep.
			continue
		}
		if res := ParseBestEffort(tc.sql(maxTreeHeight + 1)); len(res.File.Stmts) != 0 || len(res.Errors) == 0 || res.Errors[0].Message != nestingTooDeep {
			t.Errorf("%s: %d terms: best effort kept %d statements, errors %v", tc.name, maxTreeHeight+1, len(res.File.Stmts), res.Errors)
		}
	}
	n := maxTreeHeight + 1
	if _, errs := ParseExpression(strings.Repeat("1+", n) + "1"); len(errs) != 1 || errs[0].Message != nestingTooDeep {
		t.Errorf("ParseExpression: errs = %v, want the nesting error", errs)
	}
	if _, errs := ParseDataType("int" + strings.Repeat(" ARRAY", n)); len(errs) != 1 || errs[0].Message != nestingTooDeep {
		t.Errorf("ParseDataType: errs = %v, want the nesting error", errs)
	}
}

// TestTreeHeightGrowsWithChains: the height counts the nodes the
// consumers walk, the parser's own expression, query, relation, type, and
// row pattern nodes among them: a chain of n terms is at least n deep.
func TestTreeHeightGrowsWithChains(t *testing.T) {
	for _, tc := range chainShapes {
		if h := statementHeight(t, tc.sql(1000)); h < 1000 {
			t.Errorf("%s: 1,000 terms measure %d deep", tc.name, h)
		}
	}
}

// TestTreeHeightPerByte: treeTooHigh skips text shorter than maxTreeHeight/16
// bytes, on the ground that no construct adds more than 16 levels of tree
// per byte of text. Check that across the recursive constructs and the
// chains.
func TestTreeHeightPerByte(t *testing.T) {
	check := func(name, sql string) {
		if _, err := Parse(sql); err != nil {
			return
		}
		if h := statementHeight(t, sql); h > 16*len(sql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(sql))
		}
	}
	for _, tc := range nestingShapes {
		check(tc.name, tc.sql(1000))
	}
	for _, tc := range chainShapes {
		check(tc.name, tc.sql(1000))
	}
}
