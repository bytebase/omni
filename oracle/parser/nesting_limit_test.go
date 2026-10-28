package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	setNestingLimit(t, 1<<20)
}

// setNestingLimit sets maxNestingDepth for the rest of the test.
func setNestingLimit(t *testing.T, limit int) {
	saved := maxNestingDepth
	maxNestingDepth = limit
	t.Cleanup(func() { maxNestingDepth = saved })
}

// parseWithin parses sql and fails when that takes over ten seconds, so a
// speculative parse that retried at every level past the limit fails the
// test instead of hanging it.
func parseWithin(sql string) error {
	parsed := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		parsed <- err
	}()
	select {
	case err := <-parsed:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("Parse took over 10s")
	}
}

// isNestingError reports whether err is the nesting error, which a lexical
// error at a directive's condition reports with the token it is at.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && strings.HasPrefix(pe.Message, nestingErrorMessage)
}

// briefErr shortens an error for a test message: a syntax error quotes the
// text it is at, which can be most of a deep statement.
func briefErr(err error) string {
	if err == nil {
		return "<nil>"
	}
	if msg := err.Error(); len(msg) > 200 {
		return msg[:200] + "..."
	}
	return err.Error()
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 700,000 nested parentheses were enough.
// The same constructs nested 1,000 deep still parse. Each deep parse runs
// against a deadline.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	block := func(body string) string { return "BEGIN " + body + " END;" }
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") + " FROM dual" }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1 FROM dual" }},
		{"unary plus", func(n int) string { return "SELECT " + strings.Repeat("+ ", n) + "1 FROM dual" }},
		{"NOT", func(n int) string { return "SELECT 1 FROM dual WHERE " + strings.Repeat("NOT ", n) + "1 = 1" }},
		{"PRIOR", func(n int) string {
			return "SELECT 1 FROM t CONNECT BY " + strings.Repeat("PRIOR ", n) + "a = b"
		}},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") + " FROM dual" }},
		{"MULTISET operators", func(n int) string {
			return "SELECT a" + strings.Repeat(" MULTISET UNION a", n) + " FROM t"
		}},
		{"function call", func(n int) string { return "SELECT " + nest(n, "abs(", "1", ")") + " FROM dual" }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS NUMBER)") + " FROM dual" }},
		{"CASE", func(n int) string {
			return "SELECT " + nest(n, "CASE WHEN 1 = 1 THEN ", "1", " END") + " FROM dual"
		}},
		{"DECODE", func(n int) string { return "SELECT " + nest(n, "DECODE(", "1", ", 1, 1)") + " FROM dual" }},
		{"IN list", func(n int) string { return "SELECT 1 FROM dual WHERE " + nest(n, "1 IN (", "1", ")") + " = 1" }},
		{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1 FROM dual", ") FROM dual") }},
		{"EXISTS", func(n int) string {
			return nest(n, "SELECT 1 FROM dual WHERE EXISTS (", "SELECT 1 FROM dual", ")")
		}},
		{"IN subquery", func(n int) string {
			return nest(n, "SELECT 1 FROM dual WHERE 1 IN (", "SELECT 1 FROM dual", ")")
		}},
		{"CURSOR expression", func(n int) string { return nest(n, "SELECT CURSOR(", "SELECT 1 FROM dual", ") FROM dual") }},
		{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1 FROM dual", ") x") }},
		{"LATERAL derived table", func(n int) string {
			return nest(n, "SELECT * FROM t, LATERAL (", "SELECT 1 FROM dual", ") x")
		}},
		{"join group", func(n int) string { return "SELECT 1 FROM " + nest(n, "(t JOIN ", "t", " ON 1 = 1)") }},
		{"parenthesized table", func(n int) string { return "SELECT 1 FROM " + nest(n, "(", "t", ")") }},
		{"parenthesized query", func(n int) string { return nest(n, "(", "SELECT 1 FROM dual", ")") }},
		{"parenthesized right operand", func(n int) string {
			return nest(n, "(SELECT 1 FROM dual UNION ", "SELECT 1 FROM dual", ")")
		}},
		{"parenthesized left operand", func(n int) string {
			return nest(n, "(", "SELECT 1 FROM dual", ") EXCEPT SELECT 1 FROM dual")
		}},
		{"set operation in a derived table", func(n int) string {
			return nest(n, "SELECT 1 FROM dual UNION SELECT * FROM (", "SELECT 1 FROM dual", ")")
		}},
		{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1 FROM dual", ") SELECT * FROM c") }},
		{"GROUP BY grouping sets", func(n int) string {
			return "SELECT 1 FROM t GROUP BY " + nest(n, "GROUPING SETS (", "a", ")")
		}},
		{"JSON_TABLE nested columns", func(n int) string {
			return "SELECT * FROM JSON_TABLE('{}', '$' COLUMNS (" +
				nest(n, "NESTED PATH '$' COLUMNS (", "a NUMBER PATH '$.a'", ")") + ")) jt"
		}},
		{"EXPLAIN PLAN", func(n int) string { return strings.Repeat("EXPLAIN PLAN FOR ", n) + "SELECT 1 FROM dual" }},
		{"hierarchy levels", func(n int) string {
			return "CREATE HIERARCHY h USING d (" + strings.Repeat("l CHILD OF ", n) + "l)"
		}},
		{"subpartitions", func(n int) string {
			return "CREATE TABLE t (a NUMBER) PARTITION BY RANGE (a) " +
				strings.Repeat("SUBPARTITION BY HASH (a) ", n) + "(PARTITION p VALUES LESS THAN (1))"
		}},
		{"PL/SQL block", func(n int) string { return nest(n, "BEGIN ", "NULL;", " END;") }},
		{"PL/SQL IF", func(n int) string { return block(nest(n, "IF 1 = 1 THEN ", "NULL;", " END IF;")) }},
		{"PL/SQL LOOP", func(n int) string { return block(nest(n, "LOOP ", "EXIT;", " END LOOP;")) }},
		{"PL/SQL WHILE", func(n int) string { return block(nest(n, "WHILE 1 = 1 LOOP ", "EXIT;", " END LOOP;")) }},
		{"PL/SQL FOR", func(n int) string { return block(nest(n, "FOR i IN 1..2 LOOP ", "NULL;", " END LOOP;")) }},
		{"PL/SQL CASE", func(n int) string {
			return block(nest(n, "CASE WHEN 1 = 1 THEN ", "NULL;", " END CASE;"))
		}},
		{"exception handler", func(n int) string {
			return nest(n, "BEGIN NULL; EXCEPTION WHEN OTHERS THEN ", "NULL;", " END;")
		}},
		{"nested subprograms", func(n int) string {
			// The last body closes the DECLARE.
			return "DECLARE " + nest(n, "PROCEDURE p IS ", "BEGIN NULL; END;", " BEGIN NULL; END;")
		}},
		{"PL/SQL expression", func(n int) string {
			return "DECLARE x NUMBER; BEGIN x := " + nest(n, "(", "1", ")") + "; END;"
		}},
		{"directive condition", func(n int) string {
			return block("$IF " + nest(n, "(", "TRUE", ")") + " $THEN NULL; $END")
		}},
		{"unclosed blocks", func(n int) string { return strings.Repeat("BEGIN ", n) + "NULL;" }},
		{"unclosed IFs", func(n int) string { return "BEGIN " + strings.Repeat("IF 1 = 1 THEN ", n) }},
		{"unclosed subprograms", func(n int) string { return "DECLARE " + strings.Repeat("PROCEDURE p IS ", n) }},
	} {
		if err := parseWithin(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, briefErr(err))
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			if err := parseWithin(tc.sql(n)); !isNestingError(err) {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, briefErr(err))
			}
		}
	}
}

// TestUnclosedParenthesesFailFast: the parentheses of a statement must
// balance before it is parsed, so unclosed ones fail with a syntax error at
// any depth without descending into them.
func TestUnclosedParenthesesFailFast(t *testing.T) {
	for name, sql := range map[string]string{
		"parentheses": "SELECT " + strings.Repeat("(", 200000),
		"subqueries":  strings.Repeat("SELECT (", 200000),
		"join groups": "SELECT * FROM " + strings.Repeat("(", 200000),
	} {
		var pe *ParseError
		if err := parseWithin(sql); !errors.As(err, &pe) {
			t.Errorf("unclosed %s: err = %v, want a syntax error", name, briefErr(err))
		}
	}
}

// TestSetOperationChainIsNotNested: a chain of set operations is long, not
// deep, and Oracle 23ai runs one of 10,001 UNION ALL arms, so its operands
// do not count toward the nesting limit. They still chain right-deep through
// Rarg, every query of the chain ending where the last one does.
func TestSetOperationChainIsNotNested(t *testing.T) {
	for _, n := range []int{maxNestingDepth + 1, 200000} {
		sql := "SELECT 1 FROM dual" + strings.Repeat(" UNION ALL SELECT 1 FROM dual", n)
		if _, err := Parse(sql); err != nil {
			t.Errorf("%d UNION ALL arms: %v", n, briefErr(err))
		}
	}
	sql := "SELECT 1 FROM a UNION ALL SELECT 2 FROM b INTERSECT SELECT 3 FROM c MINUS SELECT 4 FROM d"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	sel := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	for i, want := range []nodes.SetOperation{nodes.SETOP_UNION, nodes.SETOP_INTERSECT, nodes.SETOP_MINUS, nodes.SETOP_NONE} {
		if sel == nil {
			t.Fatalf("operand %d missing", i)
		}
		if sel.Op != want || sel.SetAll != (i == 0) || sel.Loc.End != len(sql) {
			t.Errorf("operand %d: Op %v, SetAll %v, End %d; want %v, %v, %d", i, sel.Op, sel.SetAll, sel.Loc.End, want, i == 0, len(sql))
		}
		sel = sel.Rarg
	}
	if sel != nil {
		t.Errorf("a fifth operand: %+v", sel)
	}
}

// TestLexerLoopsOverComments: the lexer moved from a comment, or a hint
// that stands where it is a comment, to the next token by calling itself,
// one call per comment, so a long script of comments overflowed the stack.
func TestLexerLoopsOverComments(t *testing.T) {
	for name, sql := range map[string]string{
		"block comments":     "SELECT 1 " + strings.Repeat("/* c */ ", 2000000) + "FROM dual",
		"line comments":      "SELECT 1 " + strings.Repeat("-- c\n", 2000000) + "FROM dual",
		"misplaced hints":    "BEGIN " + strings.Repeat("/*+ c */ ", 2000000) + "NULL; END;",
		"comments in a unit": "BEGIN " + strings.Repeat("/* c */ ", 2000000) + "NULL; END;",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: %v", name, briefErr(err))
		}
	}
}

// TestCompletionNestingBounds: completion reads the branches of nested $IF
// directives, and the columns of nested derived tables, recursing per
// level; past the nesting limit it leaves the deeper ones out instead of
// recursing until the stack overflows (4 million nested $IF directives
// did). A lowered limit shows what is left out.
func TestCompletionNestingBounds(t *testing.T) {
	setNestingLimit(t, 3)
	texts := func(tokens []Token) string {
		var words []string
		for _, tok := range tokens {
			words = append(words, tok.Str)
		}
		return strings.Join(words, " ")
	}
	sql := "SELECT a $IF TRUE $THEN , b $IF TRUE $THEN , c $IF TRUE $THEN , d $IF TRUE $THEN , e $ERROR 'x' $END $END , f $END , g $END , h $END , i FROM t"
	if got, want := texts(completionTokens(sql, len(sql))), "SELECT A , B , C , D , F , G , H , I FROM T"; got != want {
		t.Errorf("completion tokens = %q, want %q", got, want)
	}

	query := func(depth int) []Token {
		return Tokenize(strings.Repeat("SELECT * FROM (", depth) + "SELECT 1 a FROM dual" + strings.Repeat(") x", depth))
	}
	if got := extractOracleSelectListColumns(query(2), nil); len(got) != 1 || got[0] != "A" {
		t.Errorf("columns 2 derived tables deep = %q, want [A]", got)
	}
	if got := extractOracleSelectListColumns(query(3), nil); got != nil {
		t.Errorf("columns 3 derived tables deep = %q, want none", got)
	}
	// The wrappers of a query come off in one pass, so a run of them
	// counts once.
	wrapped := Tokenize("SELECT * FROM (" + strings.Repeat("(", 5) + "SELECT 1 a FROM dual" + strings.Repeat(")", 5) + ") x")
	if got := extractOracleSelectListColumns(wrapped, nil); len(got) != 1 || got[0] != "A" {
		t.Errorf("columns through 5 wrappers = %q, want [A]", got)
	}
}
