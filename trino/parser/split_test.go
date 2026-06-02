package parser

import "testing"

// TestSplit_Single verifies a single statement (no trailing semicolon) yields
// one segment spanning the whole input.
func TestSplit_Single(t *testing.T) {
	segs := Split("SELECT 1")
	if len(segs) != 1 {
		t.Fatalf("got %d segments, want 1: %+v", len(segs), segs)
	}
	if segs[0].Text != "SELECT 1" {
		t.Errorf("Text = %q, want %q", segs[0].Text, "SELECT 1")
	}
	if segs[0].ByteStart != 0 || segs[0].ByteEnd != len("SELECT 1") {
		t.Errorf("range = [%d,%d), want [0,%d)", segs[0].ByteStart, segs[0].ByteEnd, len("SELECT 1"))
	}
}

// TestSplit_TwoStatements verifies splitting on a top-level semicolon, with the
// delimiter excluded from each segment's text/range.
func TestSplit_TwoStatements(t *testing.T) {
	in := "SELECT 1; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != "SELECT 1" {
		t.Errorf("seg[0].Text = %q, want %q", segs[0].Text, "SELECT 1")
	}
	if segs[1].Text != " SELECT 2" {
		t.Errorf("seg[1].Text = %q, want %q", segs[1].Text, " SELECT 2")
	}
	// The ';' is at index 8: seg[0] stops before it, seg[1] starts after it.
	if segs[0].ByteEnd != 8 {
		t.Errorf("seg[0].ByteEnd = %d, want 8", segs[0].ByteEnd)
	}
	if segs[1].ByteStart != 9 {
		t.Errorf("seg[1].ByteStart = %d, want 9", segs[1].ByteStart)
	}
}

// TestSplit_EmptySegmentsFiltered verifies lone semicolons and comment-only
// chunks produce no segments.
func TestSplit_EmptySegmentsFiltered(t *testing.T) {
	for _, in := range []string{";", ";;;", "  ;  ", "-- c\n;", "/* c */", ""} {
		segs := Split(in)
		if len(segs) != 0 {
			t.Errorf("Split(%q): got %d segments, want 0: %+v", in, len(segs), segs)
		}
	}
}

// TestSplit_SemicolonInString verifies a ';' inside a single-quoted string
// literal does NOT split the statement.
func TestSplit_SemicolonInString(t *testing.T) {
	in := "SELECT 'a;b' FROM t; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != "SELECT 'a;b' FROM t" {
		t.Errorf("seg[0].Text = %q, want %q", segs[0].Text, "SELECT 'a;b' FROM t")
	}
}

// TestSplit_SemicolonInQuotedIdent verifies a ';' inside a double-quoted
// identifier does NOT split.
func TestSplit_SemicolonInQuotedIdent(t *testing.T) {
	in := `SELECT "a;b" FROM t; SELECT 2`
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != `SELECT "a;b" FROM t` {
		t.Errorf("seg[0].Text = %q, want %q", segs[0].Text, `SELECT "a;b" FROM t`)
	}
}

// TestSplit_SemicolonInComment verifies a ';' inside a line or block comment
// does NOT split.
func TestSplit_SemicolonInComment(t *testing.T) {
	in := "SELECT 1 -- a;b\n; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("line comment: got %d segments, want 2: %+v", len(segs), segs)
	}
	in2 := "SELECT 1 /* a;b */; SELECT 2"
	segs2 := Split(in2)
	if len(segs2) != 2 {
		t.Fatalf("block comment: got %d segments, want 2: %+v", len(segs2), segs2)
	}
}

// TestSplit_TrailingSemicolon verifies a single statement with a trailing ';'
// yields exactly one segment (the empty trailing chunk is filtered).
func TestSplit_TrailingSemicolon(t *testing.T) {
	segs := Split("SELECT 1;")
	if len(segs) != 1 {
		t.Fatalf("got %d segments, want 1: %+v", len(segs), segs)
	}
	if segs[0].Text != "SELECT 1" {
		t.Errorf("Text = %q, want %q", segs[0].Text, "SELECT 1")
	}
}

// TestSplit_RoutineBodyNotSplit verifies that a ';' inside a CREATE FUNCTION
// routine body (BEGIN ... END) does NOT split the statement. The whole
// CREATE FUNCTION must remain a single segment.
func TestSplit_RoutineBodyNotSplit(t *testing.T) {
	in := "CREATE FUNCTION f() RETURNS int BEGIN DECLARE x int; SET x = 1; RETURN x; END"
	segs := Split(in)
	if len(segs) != 1 {
		t.Fatalf("got %d segments, want 1 (routine body must not split): %+v", len(segs), segs)
	}
}

// TestSplit_StartTransactionNotABlock verifies that a top-level BEGIN that is a
// transaction start (e.g. "BEGIN; SELECT 1") still splits normally — BEGIN as
// a TCL transaction-start is NOT a compound block opener.
func TestSplit_StartTransactionNotABlock(t *testing.T) {
	in := "START TRANSACTION; SELECT 1; COMMIT"
	segs := Split(in)
	if len(segs) != 3 {
		t.Fatalf("got %d segments, want 3: %+v", len(segs), segs)
	}
}

// TestSplit_EmptyInput verifies Split returns nil for empty input.
func TestSplit_EmptyInput(t *testing.T) {
	if segs := Split(""); segs != nil {
		t.Errorf("Split(\"\") = %+v, want nil", segs)
	}
}

// TestSplit_InlineRoutineBareReturnThenSplit verifies that a WITH FUNCTION
// whose body is a bare RETURN (no BEGIN block) does not swallow the following
// statement's ';'. The inline-routine query is one statement; a second
// statement after ';' must split off.
func TestSplit_InlineRoutineBareReturnThenSplit(t *testing.T) {
	in := "WITH FUNCTION f(x integer) RETURNS integer RETURN x * 2 SELECT f(1); SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != "WITH FUNCTION f(x integer) RETURNS integer RETURN x * 2 SELECT f(1)" {
		t.Errorf("seg[0].Text = %q", segs[0].Text)
	}
}

// TestSplit_CaseExpressionNotABlock verifies a CASE expression in a normal
// query (not a routine body) does NOT confuse the splitter: a ';' after the
// CASE still splits.
func TestSplit_CaseExpressionNotABlock(t *testing.T) {
	in := "SELECT CASE WHEN a THEN 1 ELSE 2 END FROM t; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
}

// TestSplit_DropFunctionNotARoutineBody verifies DROP FUNCTION (FUNCTION at the
// same token index as a routine definition, but a different leading keyword) is
// NOT treated as having a control-flow body — a following ';' splits normally.
func TestSplit_DropFunctionNotARoutineBody(t *testing.T) {
	in := "DROP FUNCTION f; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
}

// TestSplit_FunctionAsColumnNameNotARoutine verifies the non-reserved keyword
// FUNCTION used as a column name inside a CTE does not trip routine-body
// detection: the ';' after a CASE expression still splits.
func TestSplit_FunctionAsColumnNameNotARoutine(t *testing.T) {
	in := "WITH t AS (SELECT 1 AS function) SELECT CASE WHEN x THEN 1 END FROM t; SELECT 2"
	segs := Split(in)
	if len(segs) != 2 {
		t.Fatalf("got %d segments, want 2: %+v", len(segs), segs)
	}
}

// TestSplit_NestedFunctionBlocksOneSegment verifies a CREATE FUNCTION with a
// nested IF inside its BEGIN block (each with its own internal ';') remains a
// single segment.
func TestSplit_NestedFunctionBlocksOneSegment(t *testing.T) {
	in := "CREATE FUNCTION g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END"
	segs := Split(in)
	if len(segs) != 1 {
		t.Fatalf("got %d segments, want 1 (nested blocks must not split): %+v", len(segs), segs)
	}
}

// TestSplit_TypedEndClosersThenSplit verifies typed block closers
// (END IF / END CASE / END WHILE / END LOOP / END REPEAT) inside a routine body
// keep depth tracking balanced, so a ';' that terminates the CREATE FUNCTION
// still splits the following statement off. Regression for an END-suffix
// miscount that swallowed the post-routine ';'.
func TestSplit_TypedEndClosersThenSplit(t *testing.T) {
	cases := []struct {
		desc string
		in   string
	}{
		{"END IF", "CREATE FUNCTION g(x int) RETURNS int BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END; SELECT 2"},
		{"END CASE", "CREATE FUNCTION g(x int) RETURNS int BEGIN CASE WHEN x > 0 THEN RETURN 1; ELSE RETURN 0; END CASE; END; SELECT 2"},
		{"END WHILE", "CREATE FUNCTION h(x int) RETURNS int BEGIN DECLARE i int DEFAULT 0; WHILE i < x DO SET i = i + 1; END WHILE; RETURN i; END; SELECT 9"},
		{"END LOOP + nested END IF", "CREATE FUNCTION lp() RETURNS int BEGIN DECLARE i int DEFAULT 0; abc: LOOP SET i = i + 1; IF i > 10 THEN LEAVE abc; END IF; END LOOP; RETURN i; END; SELECT 1"},
	}
	for _, c := range cases {
		segs := Split(c.in)
		if len(segs) != 2 {
			t.Errorf("[%s] got %d segments, want 2 (routine + trailing SELECT): %+v", c.desc, len(segs), segs)
		}
	}
}

// TestSegment_Empty verifies the Empty predicate.
func TestSegment_Empty(t *testing.T) {
	if !(Segment{Text: "  -- x\n "}).Empty() {
		t.Error("comment/whitespace-only segment should be Empty")
	}
	if (Segment{Text: "SELECT 1"}).Empty() {
		t.Error("statement segment should not be Empty")
	}
}
