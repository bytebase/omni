package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	mssqldrv "github.com/microsoft/go-mssqldb"

	nodes "github.com/bytebase/omni/mssql/ast"
)

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

// nestingCases lists every construct the parser descends into recursively.
// shallow marks the ones that must still parse nested 1,000 deep; the rest
// are unclosed, or nest a clause SQL Server does not let nest, and only
// need to fail cleanly when deep.
var nestingCases = []struct {
	name    string
	sql     func(n int) string
	shallow bool
}{
	{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }, true},
	{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }, true},
	{"bitwise NOT", func(n int) string { return "SELECT " + strings.Repeat("~", n) + "1" }, true},
	{"NOT", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("NOT ", n) + "1 = 1" }, true},
	{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }, true},
	{"function call", func(n int) string { return "SELECT " + nest(n, "ABS(", "1", ")") }, true},
	{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS int)") }, true},
	{"CONVERT", func(n int) string { return "SELECT " + nest(n, "CONVERT(int, ", "1", ")") }, true},
	{"IIF", func(n int) string { return "SELECT " + nest(n, "IIF(1 = 1, ", "1", ", 0)") }, true},
	{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN 1 = 1 THEN ", "1", " END") }, true},
	{"BETWEEN operand", func(n int) string { return "SELECT 1 WHERE 1 BETWEEN " + nest(n, "(", "1", ")") + " AND 2" }, true},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }, true},
	{"EXISTS", func(n int) string { return nest(n, "SELECT 1 WHERE EXISTS (", "SELECT 1", ")") }, true},
	{"IN subquery", func(n int) string { return nest(n, "SELECT 1 WHERE 1 IN (", "SELECT 1", ")") }, true},
	{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1 AS a", ") AS x") }, true},
	{"parenthesized join", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON 1 = 1)") }, true},
	{"joins with their ON clauses nested", func(n int) string {
		return "SELECT * FROM t" + strings.Repeat(" JOIN t", n) + strings.Repeat(" ON 1 = 1", n)
	}, true},
	{"left-nested join groups", func(n int) string {
		return "SELECT * FROM " + strings.Repeat("(", n) + "t" + strings.Repeat(" JOIN t ON 1 = 1)", n)
	}, true},
	{"set operation operand", func(n int) string {
		return nest(n, "SELECT 1 UNION SELECT * FROM (", "SELECT 1 AS a", ") AS x")
	}, true},
	{"CTE body", func(n int) string {
		return "WITH c AS (" + nest(n, "SELECT (", "SELECT 1", ")") + " AS a) SELECT * FROM c"
	}, true},
	{"BEGIN...END", func(n int) string { return nest(n, "BEGIN ", "SELECT 1;", " END") }, true},
	{"IF", func(n int) string { return strings.Repeat("IF 1 = 1 ", n) + "SELECT 1" }, true},
	{"IF block", func(n int) string { return nest(n, "IF 1 = 1 BEGIN ", "SELECT 1;", " END") }, true},
	{"ELSE IF", func(n int) string { return "IF 1 = 1 SELECT 1" + strings.Repeat(" ELSE IF 1 = 1 SELECT 1", n) }, true},
	{"WHILE", func(n int) string { return strings.Repeat("WHILE 1 = 1 ", n) + "BREAK" }, true},
	{"TRY block", func(n int) string {
		return nest(n, "BEGIN TRY ", "SELECT 1;", " END TRY BEGIN CATCH SELECT 2; END CATCH")
	}, true},
	{"CATCH block", func(n int) string {
		return nest(n, "BEGIN TRY SELECT 1; END TRY BEGIN CATCH ", "SELECT 2;", " END CATCH")
	}, true},
	{"procedure body", func(n int) string { return "CREATE PROCEDURE p AS " + nest(n, "BEGIN ", "SELECT 1;", " END") }, true},
	{"event session predicate", func(n int) string {
		return "CREATE EVENT SESSION s ON SERVER ADD EVENT sqlserver.sql_batch_starting (WHERE " +
			nest(n, "(", "sqlserver.session_id > 50", ")") + ")"
	}, true},
	{"CREATE SCHEMA elements", func(n int) string { return strings.Repeat("CREATE SCHEMA s ", n) }, false},
	{"ALTER TABLE SET options", func(n int) string { return "ALTER TABLE t SET (" + nest(n, "a = (", "b = 1", ")") + ")" }, false},
	{"availability group options", func(n int) string {
		return "ALTER AVAILABILITY GROUP ag MODIFY REPLICA ON 's' WITH (" + nest(n, "SECONDARY_ROLE (", "ALLOW_CONNECTIONS = ALL", ")") + ")"
	}, false},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }, false},
	{"unclosed function calls", func(n int) string { return "SELECT " + strings.Repeat("ABS(", n) }, false},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }, false},
	{"unclosed derived tables", func(n int) string { return strings.Repeat("SELECT * FROM (", n) }, false},
	{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }, false},
	{"joins waiting for ON", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t", n) }, false},
	{"unclosed blocks", func(n int) string { return strings.Repeat("BEGIN ", n) }, false},
	{"unclosed WAITFOR", func(n int) string { return strings.Repeat("WAITFOR (", n) }, false},
}

// parseWithin parses sql on its own goroutine and reports false if it takes
// longer than limit, so a construct whose retries turn exponential at the
// limit fails the test instead of hanging it.
func parseWithin(sql string, limit time.Duration) (error, bool) {
	done := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		done <- err
	}()
	select {
	case err := <-done:
		return err, true
	case <-time.After(limit):
		return nil, false
	}
}

func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == "statement is nested too deeply"
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: 353KB of '(' were enough. The same constructs
// nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingCases {
		if tc.shallow {
			if _, err := Parse(tc.sql(1000)); err != nil {
				t.Errorf("%s: 1,000 levels: %v", tc.name, err)
			}
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			err, ok := parseWithin(tc.sql(n), 10*time.Second)
			if !ok {
				t.Errorf("%s: %d levels: no result within 10s", tc.name, n)
				continue
			}
			if !isNestingError(err) {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// TestNestingLimitError: the error reports its own message, not the
// generic "syntax error at or near" text, at the token where the limit was
// reached.
func TestNestingLimitError(t *testing.T) {
	sql := "SELECT " + strings.Repeat("(", maxNestingDepth)
	_, err := Parse(sql)
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %v, want a *ParseError", err)
	}
	if pe.Error() != "statement is nested too deeply" {
		t.Errorf("Error() = %q", pe.Error())
	}
	if pe.Position <= len("SELECT ") || pe.Position >= len(sql) || sql[pe.Position] != '(' {
		t.Errorf("Position = %d, want the offset of a '(' inside the run", pe.Position)
	}
}

// TestNestingLimitCompletion: completion runs the same parser, and deep
// input must come back without overflowing the stack.
func TestNestingLimitCompletion(t *testing.T) {
	for _, sql := range []string{
		"SELECT " + strings.Repeat("(", 200000),
		strings.Repeat("SELECT (", 200000),
		strings.Repeat("BEGIN ", 200000),
	} {
		done := make(chan struct{})
		go func() {
			Collect(sql, len(sql))
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Errorf("Collect on %.20q...: no result within 10s", sql)
		}
	}
}

// TestFlatListsStayShallow: statements and other flat lists are parsed by
// a loop and add no depth, so neither limit applies to their length.
func TestFlatListsStayShallow(t *testing.T) {
	for name, sql := range map[string]string{
		"statements":          strings.Repeat("SELECT 1; ", 200000),
		"statements in block": "BEGIN " + strings.Repeat("SELECT 1; ", 200000) + "END",
		"select list":         "SELECT 1" + strings.Repeat(", 1", 200000),
		"IN list":             "SELECT 1 WHERE 1 IN (1" + strings.Repeat(", 1", 200000) + ")",
	} {
		err, ok := parseWithin(sql, 10*time.Second)
		if !ok {
			t.Errorf("%s: no result within 10s", name)
		} else if err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// TestLexerLoopsOverComments: the lexer moved from a comment to the next
// token by calling itself, one call per comment, so 3 million comments
// overflowed the stack.
func TestLexerLoopsOverComments(t *testing.T) {
	for name, sql := range map[string]string{
		"block comments": "SELECT 1 " + strings.Repeat("/* c */ ", 3000000),
		"line comments":  "SELECT 1 " + strings.Repeat("-- c\n", 3000000),
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// treeHeight returns how many nodes deep the tree under root is.
func treeHeight(root nodes.Node) int {
	lo, hi := 1, 1
	for treeHeightExceeds(root, hi) {
		lo, hi = hi+1, hi*2
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if treeHeightExceeds(root, mid) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// chainShapes build left-deep chains of n terms: the parser builds each with
// a loop, so only the tree, not the parse, gets deep.
var chainShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1 + ", n) + "1" }},
	{"subtraction", func(n int) string { return "SELECT " + strings.Repeat("1 - ", n) + "1" }},
	{"concatenation", func(n int) string { return "SELECT " + strings.Repeat("'a' + ", n) + "'a'" }},
	{"bitwise operators", func(n int) string { return "SELECT " + strings.Repeat("1 & 1 | ", n) + "1 ^ 1" }},
	{"product", func(n int) string { return "SELECT " + strings.Repeat("1 * ", n) + "1" }},
	{"modulo", func(n int) string { return "SELECT " + strings.Repeat("1 % ", n) + "1" }},
	{"AND", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a = 1 AND ", n) + "a = 1" }},
	{"OR", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a = 1 OR ", n) + "a = 1" }},
	{"COLLATE", func(n int) string { return "SELECT 'a'" + strings.Repeat(" COLLATE Latin1_General_CI_AS", n) }},
	{"AT TIME ZONE", func(n int) string { return "SELECT d" + strings.Repeat(" AT TIME ZONE 'UTC'", n) + " FROM t" }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT 1 UNION ALL ", n) + "SELECT 1" }},
	{"EXCEPT and INTERSECT", func(n int) string {
		return strings.Repeat("SELECT 1 EXCEPT SELECT 1 INTERSECT ", n) + "SELECT 1"
	}},
	{"joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON 1 = 1", n) }},
	{"outer joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" LEFT OUTER JOIN t ON 1 = 1", n) }},
	{"cross joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" CROSS JOIN t", n) }},
	{"APPLY", func(n int) string {
		return "SELECT * FROM t" + strings.Repeat(" CROSS APPLY f() AS x OUTER APPLY f() AS y", n)
	}},
}

// TestTreeHeightLimit: a left-deep chain whose tree is deeper than
// maxTreeHeight is rejected at the statement, because every consumer that
// walks the tree recursively would overflow the stack on it (NodeToString
// did at a million-term sum, review at three million ANDs). A chain within
// the limit still parses, though it is longer than maxNestingDepth: the
// parser itself does not recurse per term.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if err, ok := parseWithin(tc.sql(maxTreeHeight/4), 10*time.Second); !ok || err != nil {
			t.Errorf("%s: %d terms: err = %v, finished = %v", tc.name, maxTreeHeight/4, err, ok)
		}
		err, ok := parseWithin("SELECT 1; "+tc.sql(maxTreeHeight+1), 10*time.Second)
		var pe *ParseError
		if !ok || !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" || pe.Position != len("SELECT 1; ") {
			t.Errorf("%s: %d terms: err = %v, finished = %v, want the nesting error at the statement",
				tc.name, maxTreeHeight+1, err, ok)
		}
	}
}

// TestTreeHeightPerByte: treeTooHigh skips statements shorter than
// maxTreeHeight/16 bytes, on the ground that no construct adds more than 16
// levels of tree per byte of text. Check that across the recursive
// constructs and the chains, and that the measured height of a chain grows
// with its length, as the recursive consumers' depth does.
func TestTreeHeightPerByte(t *testing.T) {
	height := func(name, sql string) int {
		list, err := Parse(sql)
		if err != nil || list == nil || len(list.Items) == 0 {
			return 0
		}
		h := treeHeight(list.Items[0])
		if h > 16*len(sql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(sql))
		}
		return h
	}
	for _, tc := range nestingCases {
		height(tc.name, tc.sql(1000))
	}
	for _, tc := range chainShapes {
		if h := height(tc.name, tc.sql(1000)); h < 1000 {
			t.Errorf("%s: 1,000 terms measured %d deep, want at least 1,000", tc.name, h)
		}
	}
}

// TestNestingLimitBelowServer: SQL Server rejects every recursive construct
// with Msg 191 ("Some part of your SQL statement is nested too deeply")
// well before omni's limit, so the limit never rejects a statement SQL
// Server accepts. 1,100 levels is past SQL Server 2022's limit for each
// construct below (1,016 parentheses, 505 BEGIN...END blocks, 169 scalar
// subqueries, ...) and well inside omni's.
func TestNestingLimitBelowServer(t *testing.T) {
	oracle := startParserOracle(t)
	const n = 1100
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{"parenthesized expression", "SELECT " + nest(n, "(", "1", ")")},
		{"unary minus", "SELECT " + strings.Repeat("- ", n) + "1"},
		{"NOT", "SELECT 1 WHERE " + strings.Repeat("NOT ", n) + "1 = 1"},
		{"function call", "SELECT " + nest(n, "ABS(", "1", ")")},
		{"scalar subquery", nest(n, "SELECT (", "SELECT 1", ")")},
		{"derived table", nest(n, "SELECT * FROM (", "SELECT 1 AS a", ") AS x")},
		{"BEGIN...END", nest(n, "BEGIN ", "SELECT 1;", " END")},
		{"IF", strings.Repeat("IF 1 = 1 ", n) + "SELECT 1"},
		{"ELSE IF", "IF 1 = 1 SELECT 1" + strings.Repeat(" ELSE IF 1 = 1 SELECT 1", n)},
		{"WHILE", strings.Repeat("WHILE 1 = 1 ", n) + "BREAK"},
		{"TRY block", nest(n, "BEGIN TRY ", "SELECT 1;", " END TRY BEGIN CATCH SELECT 2; END CATCH")},
	} {
		serverErr, err := oracle.parseError(tc.sql)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		var me mssqldrv.Error
		if !errors.As(serverErr, &me) || me.Number != 191 {
			t.Errorf("%s: SQL Server: %v, want Msg 191", tc.name, serverErr)
		}
		if _, err := Parse(tc.sql); err != nil {
			t.Errorf("%s: omni: %v", tc.name, err)
		}
	}
	// Flat chains and lists are not nesting: both accept 50,000 terms.
	for name, sql := range map[string]string{
		"UNION":               "SELECT 1" + strings.Repeat(" UNION SELECT 1", 50000),
		"AND":                 "SELECT 1 WHERE " + strings.Repeat("1 = 1 AND ", 50000) + "1 = 1",
		"statements in block": "BEGIN " + strings.Repeat("SELECT 1; ", 50000) + "END",
	} {
		serverErr, err := oracle.parseError(sql)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if serverErr != nil {
			t.Errorf("%s: SQL Server: %v, want it accepted", name, serverErr)
		}
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: omni: %v", name, err)
		}
	}
	// A long sum already exhausts SQL Server's own stack (maxTreeHeight).
	serverErr, err := oracle.parseError("SELECT " + strings.Repeat("1 + ", 50000) + "1")
	if err != nil {
		t.Fatal(err)
	}
	var me mssqldrv.Error
	if !errors.As(serverErr, &me) || me.Number != 8631 {
		t.Errorf("50,000-term sum: SQL Server: %v, want Msg 8631", serverErr)
	}
}
