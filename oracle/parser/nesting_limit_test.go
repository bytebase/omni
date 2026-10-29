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

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

func block(body string) string { return "BEGIN " + body + " END;" }

// nestingShapes are the constructs the parser descends into recursively,
// n levels deep.
var nestingShapes = []struct {
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
	// The operators that associate to the right recurse.
	{"MULTISET operators", func(n int) string {
		return "SELECT a" + strings.Repeat(" MULTISET UNION a", n) + " FROM t"
	}},
	{"exponentiation", func(n int) string {
		return "DECLARE x NUMBER; BEGIN x := 1" + strings.Repeat(" ** 1", n) + "; END;"
	}},
	{"AT TIME ZONE", func(n int) string {
		return "SELECT d" + strings.Repeat(" AT TIME ZONE 'UTC'", n) + " FROM t"
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
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 700,000 nested parentheses were enough.
// The same constructs nested 1,000 deep still parse. Each deep parse runs
// against a deadline.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
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
// Rarg, every query of the chain ending where the last one does, so the
// chain's tree is bounded by maxTreeHeight (TestTreeHeightLimit).
func TestSetOperationChainIsNotNested(t *testing.T) {
	for _, n := range []int{maxNestingDepth + 1, maxTreeHeight / 2} {
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

// walkDepth is a Visitor that records the deepest level ast.Walk reaches.
type walkDepth struct{ depth, max int }

func (w *walkDepth) Visit(n nodes.Node) nodes.Visitor {
	if n == nil {
		w.depth--
		return nil
	}
	w.depth++
	w.max = max(w.max, w.depth)
	return w
}

// chainShapes build left-deep chains of n terms: the parser builds each with
// a loop, so only the tree, not the parse, gets deep.
var chainShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"sum", func(n int) string { return "SELECT 1" + strings.Repeat(" + 1", n) + " FROM dual" }},
	{"difference", func(n int) string { return "SELECT 1" + strings.Repeat(" - 1", n) + " FROM dual" }},
	{"product", func(n int) string { return "SELECT 1" + strings.Repeat(" * 1", n) + " FROM dual" }},
	{"concatenation", func(n int) string { return "SELECT 'a'" + strings.Repeat(" || 'a'", n) + " FROM dual" }},
	{"IS NULL", func(n int) string { return "SELECT 1 FROM t WHERE a" + strings.Repeat(" IS NULL", n) }},
	{"LIKE", func(n int) string { return "SELECT 1 FROM t WHERE a" + strings.Repeat(" LIKE 'x'", n) }},
	{"BETWEEN", func(n int) string { return "SELECT 1 FROM t WHERE a" + strings.Repeat(" BETWEEN 1 AND 2", n) }},
	{"IN", func(n int) string { return "SELECT 1 FROM t WHERE a" + strings.Repeat(" IN (1)", n) }},
	{"cursor attributes", func(n int) string {
		return block("IF c%ROWCOUNT" + strings.Repeat("%ROWCOUNT", n) + " > 0 THEN NULL; END IF;")
	}},
	{"AT LOCAL", func(n int) string { return "SELECT d" + strings.Repeat(" AT LOCAL", n) + " FROM t" }},
	{"method calls", func(n int) string { return "SELECT x.f()" + strings.Repeat(".f()", n) + " FROM t" }},
	{"field access", func(n int) string { return "SELECT (x).a" + strings.Repeat(".a", n) + " FROM t" }},
	{"UNION ALL", func(n int) string {
		return "SELECT 1 FROM dual" + strings.Repeat(" UNION ALL SELECT 1 FROM dual", n)
	}},
	{"INTERSECT", func(n int) string {
		return "SELECT 1 FROM dual" + strings.Repeat(" INTERSECT SELECT 1 FROM dual", n)
	}},
	{"MINUS", func(n int) string { return "SELECT 1 FROM dual" + strings.Repeat(" MINUS SELECT 1 FROM dual", n) }},
	{"EXCEPT ALL", func(n int) string {
		return "SELECT 1 FROM dual" + strings.Repeat(" EXCEPT ALL SELECT 1 FROM dual", n)
	}},
	{"UNION ALL of parenthesized operands", func(n int) string {
		return "(SELECT 1 FROM dual)" + strings.Repeat(" UNION ALL (SELECT 1 FROM dual)", n)
	}},
	{"UNION after a parenthesized operand", func(n int) string {
		return "(SELECT 1 FROM dual)" + strings.Repeat(" UNION SELECT 1 FROM dual", n)
	}},
	{"INSERT of a UNION ALL", func(n int) string {
		return "INSERT INTO t SELECT 1 FROM dual" + strings.Repeat(" UNION ALL SELECT 1 FROM dual", n)
	}},
	{"view of a UNION ALL", func(n int) string {
		return "CREATE VIEW v AS SELECT 1 FROM dual" + strings.Repeat(" UNION ALL SELECT 1 FROM dual", n)
	}},
	{"joins", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(" JOIN t ON 1 = 1", n) }},
	{"cross joins", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(" CROSS JOIN t", n) }},
	{"left joins", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(" LEFT JOIN t ON 1 = 1", n) }},
	{"natural joins", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(" NATURAL JOIN t", n) }},
	{"OUTER APPLY", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(" OUTER APPLY t", n) }},
	{"PIVOT", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" PIVOT (SUM(a) FOR b IN (1))", n) }},
	{"UNPIVOT", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" UNPIVOT (a FOR b IN (c))", n) }},
	{"PL/SQL sum", func(n int) string {
		return "DECLARE x NUMBER; BEGIN x := 1" + strings.Repeat(" + 1", n) + "; END;"
	}},
	{"PL/SQL concatenation", func(n int) string {
		return "DECLARE x VARCHAR2(9); BEGIN x := 'a'" + strings.Repeat(" || 'a'", n) + "; END;"
	}},
	{"procedure sum", func(n int) string {
		return "CREATE PROCEDURE p AS x NUMBER; BEGIN x := 1" + strings.Repeat(" + 1", n) + "; END;"
	}},
	{"embedded SQL sum", func(n int) string {
		return "DECLARE x NUMBER; BEGIN SELECT 1" + strings.Repeat(" + 1", n) + " INTO x FROM dual; END;"
	}},
}

// flatShapes are long lists the parser keeps flat, however many terms.
var flatShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"AND", func(n int) string { return "SELECT 1 FROM t WHERE a = 1" + strings.Repeat(" AND a = 1", n) }},
	{"OR", func(n int) string { return "SELECT 1 FROM t WHERE a = 1" + strings.Repeat(" OR a = 1", n) }},
	{"IN list", func(n int) string { return "SELECT 1 FROM t WHERE a IN (1" + strings.Repeat(", 1", n) + ")" }},
	{"select list", func(n int) string { return "SELECT 1" + strings.Repeat(", 1", n) + " FROM dual" }},
	{"FROM list", func(n int) string { return "SELECT 1 FROM t" + strings.Repeat(", t", n) }},
	{"CASE branches", func(n int) string {
		return "SELECT CASE" + strings.Repeat(" WHEN a = 1 THEN 1", n) + " END FROM t"
	}},
	{"ELSIF branches", func(n int) string {
		return block("IF 1 = 1 THEN NULL;" + strings.Repeat(" ELSIF 1 = 1 THEN NULL;", n) + " END IF;")
	}},
	{"PL/SQL statements", func(n int) string { return block(strings.Repeat("NULL; ", n)) }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected at its statement, because every consumer that walks the tree
// recursively would overflow the stack on it (ast.NodeToString did at a
// million terms). A chain within the limit still parses, and a list the
// parser keeps flat parses at any length.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, briefErr(err))
		}
		_, err := Parse(tc.sql(maxTreeHeight + 1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingErrorMessage || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, briefErr(err))
		}
	}
	for _, tc := range flatShapes {
		if _, err := Parse(tc.sql(maxTreeHeight + 1)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight+1, briefErr(err))
		}
	}

	// The static evaluators of a $IF condition walk its tree recursively
	// while the lexer reads it, before its statement is measured: they
	// overflowed the stack on a sum of two million terms. Oracle 23ai
	// evaluates a sum of 10,001.
	cond := func(n int) string {
		return "BEGIN $IF 1" + strings.Repeat(" + 1", n) + " > 0 $THEN NULL; $END NULL; END;"
	}
	if _, err := Parse(cond(10001)); err != nil {
		t.Errorf("a condition of 10,001 terms: %v", briefErr(err))
	}
	_, err := Parse(cond(maxTreeHeight + 1))
	var pe *ParseError
	if !errors.As(err, &pe) || !strings.HasPrefix(pe.Message, nestingErrorMessage) || pe.Position != len("BEGIN $IF ") {
		t.Errorf("a condition of %d terms: err = %v, want the nesting error at the condition", maxTreeHeight+1, briefErr(err))
	}
}

// TestTreeHeightPerByte: treeTooHigh skips statements shorter than
// maxTreeHeight/16 bytes, on the ground that no construct adds more than 16
// levels of tree per byte of text. Check that across the recursive
// constructs and the chains, and that a chain's height grows with its
// terms, as the depth ast.Walk reaches does.
func TestTreeHeightPerByte(t *testing.T) {
	const n = 1000
	check := func(name, sql string) int {
		list, err := Parse(sql)
		if err != nil || len(list.Items) == 0 {
			return 0
		}
		h := treeHeight(list)
		if h > 16*len(sql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(sql))
		}
		w := &walkDepth{}
		nodes.Walk(w, list)
		if w.max != h {
			t.Errorf("%s: treeHeight %d, ast.Walk reaches %d", name, h, w.max)
		}
		return h
	}
	for _, tc := range nestingShapes {
		check(tc.name, tc.sql(n))
	}
	for _, tc := range chainShapes {
		if h := check(tc.name, tc.sql(n)); h < n {
			t.Errorf("%s: %d terms, tree only %d deep", tc.name, n, h)
		}
	}
	for _, tc := range flatShapes {
		if h := check(tc.name, tc.sql(n)); h > 20 {
			t.Errorf("%s: %d terms, tree %d deep", tc.name, n, h)
		}
	}
}

// TestInspectMatchesAstInspect: inspect, which keeps its own stack, visits
// the nodes ast.Inspect does, in the same order, pruned where f says.
func TestInspectMatchesAstInspect(t *testing.T) {
	list, err := Parse(`DECLARE x NUMBER; BEGIN
		IF x > 1 THEN BEGIN x := f(x, 2) + 3; END; ELSE x := CASE WHEN x = 1 THEN 2 END; END IF;
		SELECT a + 1 INTO x FROM t WHERE b IN (SELECT c FROM u);
	EXCEPTION WHEN OTHERS THEN BEGIN NULL; END; END;`)
	if err != nil {
		t.Fatal(err)
	}
	for _, prune := range []bool{false, true} {
		visit := func(seen *[]nodes.Node) func(nodes.Node) bool {
			return func(n nodes.Node) bool {
				*seen = append(*seen, n)
				_, binary := n.(*nodes.BinaryExpr)
				return !(prune && binary)
			}
		}
		var want, got []nodes.Node
		nodes.Inspect(list, visit(&want))
		inspect(list, visit(&got))
		if len(got) != len(want) {
			t.Fatalf("prune %v: inspect visited %d nodes, ast.Inspect %d", prune, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("prune %v: node %d is %T, want %T", prune, i, got[i], want[i])
			}
		}
	}
}
