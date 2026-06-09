package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// TestParse_EmptyInput verifies that parsing empty or whitespace/comment-only
// input yields a non-nil File with no statements and no errors.
func TestParse_EmptyInput(t *testing.T) {
	for _, in := range []string{"", "   ", "\n\t ", "-- just a comment\n", "/* block */", "# pound\n", ";", ";;"} {
		file, errs := Parse(in)
		if file == nil {
			t.Fatalf("Parse(%q): File is nil", in)
		}
		if len(file.Stmts) != 0 {
			t.Errorf("Parse(%q): got %d stmts, want 0", in, len(file.Stmts))
		}
		if len(errs) != 0 {
			t.Errorf("Parse(%q): got errors %v, want none", in, errs)
		}
	}
}

// TestParse_FileLocCoversInput verifies the File node spans the whole input.
func TestParse_FileLocCoversInput(t *testing.T) {
	in := "SELECT 1; SELECT 2"
	file, _ := Parse(in)
	if file.Loc.Start != 0 || file.Loc.End != len(in) {
		t.Errorf("File.Loc = %+v, want {0, %d}", file.Loc, len(in))
	}
}

// TestParse_UnknownStatement verifies a statement starting with a token the
// dispatch switch does not recognize produces a parse error (not a panic, not
// silent acceptance).
func TestParse_UnknownStatement(t *testing.T) {
	_, errs := Parse("FOOBAR 1 2 3")
	if len(errs) == 0 {
		t.Fatal("Parse(\"FOOBAR ...\"): want a parse error, got none")
	}
	if !strings.Contains(errs[0].Msg, "unknown or unsupported statement") {
		t.Errorf("error = %q, want it to mention 'unknown or unsupported statement'", errs[0].Msg)
	}
}

// TestParse_KnownStatementUnsupported verifies a statement whose leading
// keyword IS in the dispatch switch but whose body is still stubbed yields a
// "not yet supported" diagnostic rather than an "unknown statement" one. IMPORT
// is used because it is still stubbed (parser-scripting, pending); the query
// family (SELECT / WITH, parser-select), the DML family (parser-dml), the
// transaction / utility family (BEGIN/COMMIT/ROLLBACK/ASSERT/ANALYZE/DESCRIBE/
// RENAME/CALL, parser-utility) and the data-movement family (EXPORT DATA / MODEL,
// LOAD DATA, CLONE DATA, parser-dml-ext) are now implemented and no longer
// stubbed.
func TestParse_KnownStatementUnsupported(t *testing.T) {
	_, errs := Parse("IMPORT MODULE foo.bar")
	if len(errs) != 1 {
		t.Fatalf("Parse(\"IMPORT ...\"): got %d errors, want 1: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "not yet supported") {
		t.Errorf("error = %q, want 'not yet supported'", errs[0].Msg)
	}
	if !strings.HasPrefix(errs[0].Msg, "IMPORT ") {
		t.Errorf("error = %q, want it to name the IMPORT statement", errs[0].Msg)
	}
}

// TestParse_MultiStatementErrorsCollected verifies ParseBestEffort collects
// errors from every segment, not just the first. The leading `SELECT 1` now
// parses cleanly (parser-select) and the INSERT now parses cleanly
// (parser-dml), so only the two still-stubbed IMPORT segments contribute errors.
func TestParse_MultiStatementErrorsCollected(t *testing.T) {
	res := ParseBestEffort("SELECT 1; INSERT INTO t VALUES (1); IMPORT MODULE a.b; IMPORT MODULE c.d")
	if got := len(res.Errors); got != 2 {
		t.Fatalf("ParseBestEffort: got %d errors, want 2: %v", got, res.Errors)
	}
}

// TestParse_LexErrorPromoted verifies a lexer-level failure (unterminated
// string) surfaces as a parse error with a position, so Diagnose can report it.
func TestParse_LexErrorPromoted(t *testing.T) {
	_, errs := Parse("SELECT 'unterminated")
	found := false
	for _, e := range errs {
		if strings.Contains(e.Msg, errUnterminatedString) {
			found = true
		}
	}
	if !found {
		t.Errorf("errors = %v, want one mentioning %q", errs, errUnterminatedString)
	}
}

// TestParse_LexErrorPositionAbsolute verifies a lex error in the SECOND
// statement carries an absolute byte offset (NewLexerWithOffset shift), not a
// segment-local one.
func TestParse_LexErrorPositionAbsolute(t *testing.T) {
	// The unterminated string begins at byte offset 17 in the full input.
	in := "SELECT 1;\nSELECT 'oops"
	want := strings.Index(in, "'oops")
	_, errs := Parse(in)
	var lexErr *ParseError
	for i := range errs {
		if strings.Contains(errs[i].Msg, errUnterminatedString) {
			lexErr = &errs[i]
		}
	}
	if lexErr == nil {
		t.Fatalf("no unterminated-string error in %v", errs)
	}
	if lexErr.Loc.Start != want {
		t.Errorf("lex error Loc.Start = %d, want %d (absolute offset)", lexErr.Loc.Start, want)
	}
}

// TestParse_StrictReturnsAllErrors verifies the strict Parse entry returns the
// full error slice (matching snowflake/trino Parse, which return []ParseError).
func TestParse_StrictReturnsAllErrors(t *testing.T) {
	_, errs := Parse("FOOBAR; BAZBAR")
	if len(errs) != 2 {
		t.Fatalf("Parse: got %d errors, want 2: %v", len(errs), errs)
	}
}

// TestParse_UnterminatedCommentReported verifies a whole-input unterminated
// block comment surfaces a diagnostic. The comment lexes to EOF, so Split drops
// its (empty) segment — but the lex error must still be reported, because
// bytebase's Diagnose promises to flag every lexer error. Lex errors are
// therefore collected from a single full-input pass, not per surviving segment.
func TestParse_UnterminatedCommentReported(t *testing.T) {
	cases := []string{
		"/* unterminated",            // whole input is the bad comment
		"SELECT 1; /* unterminated",  // bad comment after a (dropped) trailing chunk
		"/* bad */ SELECT 1 /* also", // a closed comment then an unterminated one
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			_, errs := Parse(in)
			found := false
			for _, e := range errs {
				if strings.Contains(e.Msg, errUnterminatedComment) {
					found = true
				}
			}
			if !found {
				t.Errorf("Parse(%q): no unterminated-comment diagnostic in %v", in, errs)
			}
		})
	}
}

// TestParse_LexErrorAfterRecoveryStopReported verifies a lex error positioned
// AFTER the point where statement error-recovery stops is still reported. The
// per-statement parser stops at the first ';' boundary on an unsupported/invalid
// statement, so a pull-based per-segment lexer would never reach a later
// unterminated string; the full-input lex pass catches it.
func TestParse_LexErrorAfterRecoveryStopReported(t *testing.T) {
	cases := []string{
		"CREATE PROCEDURE p() BEGIN SELECT 1; SELECT 'oops", // unterminated string deep in a block segment
		"FOOBAR a b c 'oops", // unknown stmt, then a late unterminated string
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			_, errs := Parse(in)
			found := false
			for _, e := range errs {
				if strings.Contains(e.Msg, errUnterminatedString) {
					found = true
				}
			}
			if !found {
				t.Errorf("Parse(%q): no unterminated-string diagnostic in %v", in, errs)
			}
		})
	}
}

// TestParse_UnterminatedHintReported verifies a malformed statement-level hint
// whose body never closes (`@{ … <EOF>` or `@[ … @] { … <EOF>`) yields a
// diagnostic rather than silently swallowing the statement. Without it, the hint
// skip would consume to EOF and dispatch would never run, so an invalid input
// would draw no diagnostic at all (oracle: Spanner rejects malformed hints).
func TestParse_UnterminatedHintReported(t *testing.T) {
	bad := []string{"@{k=1 SELECT 1", "@{", "@{a={b=1}", "@[5@]{k=1", "@[5@] SELECT 1"}
	for _, in := range bad {
		t.Run("bad/"+in, func(t *testing.T) {
			_, errs := Parse(in)
			if len(errs) == 0 || !strings.Contains(errs[0].Msg, "unterminated statement hint") {
				t.Errorf("Parse(%q): want an unterminated-hint diagnostic, got %v", in, errs)
			}
		})
	}
	// A well-formed hint must NOT be flagged; it should dispatch to the (stubbed)
	// statement keyword after the hint.
	good := []string{"@{USE_ADDITIONAL_PARALLELISM=TRUE} SELECT 1", "@5 SELECT 1", "@[5@]{key=1} SELECT 1"}
	for _, in := range good {
		t.Run("good/"+in, func(t *testing.T) {
			_, errs := Parse(in)
			for _, e := range errs {
				if strings.Contains(e.Msg, "unterminated statement hint") {
					t.Errorf("Parse(%q): well-formed hint wrongly flagged: %v", in, errs)
				}
			}
		})
	}
}

// TestParse_EmptyHintReported verifies a balanced but EMPTY statement-level hint
// body (`@{}`, or whitespace/comment-only between the braces) draws a diagnostic
// rather than being silently consumed. GoogleSQL requires at least one hint
// entry: oracle-confirmed against the Spanner emulator, `@{} SELECT 1` rejects
// with `Syntax error: Unexpected "}"`. Without this check the empty hint is
// skipped, parseSingle reaches EOF, and the (invalid) input draws no diagnostic
// at all — which bytebase's Diagnose must not allow.
func TestParse_EmptyHintReported(t *testing.T) {
	// Empty hint bodies — must be flagged.
	bad := []string{"@{}", "@{} SELECT 1", "@{   }", "@{ /* c */ } SELECT 1", "@[5@]{}", "@[5@]{} SELECT 1"}
	for _, in := range bad {
		t.Run("empty/"+in, func(t *testing.T) {
			_, errs := Parse(in)
			if len(errs) == 0 || !strings.Contains(errs[0].Msg, "empty statement hint") {
				t.Errorf("Parse(%q): want an empty-hint diagnostic, got %v", in, errs)
			}
		})
	}
	// A hint with content between the braces is NOT empty: the foundation only
	// SKIPS the hint to reach the statement keyword (validating the entry's
	// internal `key=value` shape is the hint-parsing node's job). So `@{k} …`
	// must dispatch to the statement, not draw an empty-hint diagnostic. (The
	// oracle does ultimately reject `@{k}` for a missing `=value`, but that is a
	// deeper hint-body parse the foundation deliberately defers.)
	nonEmpty := []string{"@{k} SELECT 1", "@{k=1} SELECT 1", "@{USE_ADDITIONAL_PARALLELISM=TRUE} SELECT 1"}
	for _, in := range nonEmpty {
		t.Run("nonempty/"+in, func(t *testing.T) {
			_, errs := Parse(in)
			for _, e := range errs {
				if strings.Contains(e.Msg, "empty statement hint") {
					t.Errorf("Parse(%q): non-empty hint wrongly flagged as empty: %v", in, errs)
				}
			}
		})
	}
}

// dispatchPrefix is one representative leading token sequence per documented
// GoogleSQL top-level statement kind (antlr_rules.md §4: sql_statement_body
// alternatives + the procedural-script forms recognized at top level). The body
// after the keyword does not matter for the leading-keyword dispatch check — we
// only need the leading keyword(s) to land in a real dispatch case rather than
// default. A bare " x" tail suffices for most; the three block forms whose body
// is a statement_list that begins immediately after the opener (BEGIN / LOOP /
// REPEAT) need a benign complete tail so the parsed body does not itself produce
// an inner "unknown statement starting with x" (the now-implemented
// parser-scripting bodies parse their statement_list, so a bare `x` body is a
// real — but irrelevant-to-this-test — reject).
var dispatchPrefixes = []struct{ prefix, tail string }{
	// Query / GQL.
	{"SELECT", "x"}, {"WITH", "x"}, {"GRAPH", "x"}, {"FROM", "x"}, {"(", "x"},
	// DDL.
	{"CREATE", "x"}, {"ALTER", "x"}, {"DROP", "x"}, {"RENAME", "x"},
	{"UNDROP", "x"}, {"TRUNCATE", "x"}, {"DEFINE", "x"},
	// DML.
	{"INSERT", "x"}, {"UPDATE", "x"}, {"DELETE", "x"}, {"MERGE", "x"},
	// DCL.
	{"GRANT", "x"}, {"REVOKE", "x"},
	// Transactions / batch / session. A top-level BEGIN is a TCL transaction only
	// when followed by a TCL follower (here TRANSACTION); a non-follower opens a
	// BEGIN…END block (covered by the scripting `{"BEGIN","END"}` entry below).
	{"BEGIN", "TRANSACTION"}, {"START", "x"}, {"COMMIT", "x"}, {"ROLLBACK", "x"},
	{"SET", "x"}, {"RUN", "x"}, {"ABORT", "x"},
	// Utility / metadata.
	{"EXPLAIN", "x"}, {"DESCRIBE", "x"}, {"DESC", "x"}, {"SHOW", "x"},
	{"ANALYZE", "x"}, {"ASSERT", "x"}, {"CALL", "x"}, {"EXECUTE", "x"},
	{"IMPORT", "x"}, {"MODULE", "x"}, {"EXPORT", "x"}, {"LOAD", "x"}, {"CLONE", "x"},
	// Procedural / scripting (now implemented by parser-scripting). BEGIN / LOOP /
	// REPEAT have a statement_list body that starts right after the opener, so
	// their tails are complete benign blocks rather than a bare `x` (which would
	// be parsed as an invalid body statement, not a leading-keyword miss).
	{"IF", "x"}, {"CASE", "x"}, {"WHILE", "x"}, {"LOOP", "END LOOP"},
	{"REPEAT", "UNTIL TRUE END REPEAT"}, {"FOR", "x"}, {"DECLARE", "x"},
	{"BREAK", "x"}, {"LEAVE", "x"}, {"CONTINUE", "x"}, {"ITERATE", "x"},
	{"RETURN", "x"}, {"RAISE", "x"}, {"BEGIN", "END"},
}

// TestParse_DispatchKeywordsRecognized verifies every documented top-level
// statement keyword is routed by the dispatch switch — it does NOT hit the
// default "unknown statement" branch. This is the foundation's completeness
// contract: all statement forms must be reachable so bytebase's Diagnose never
// emits a false "unknown statement" diagnostic for valid GoogleSQL.
func TestParse_DispatchKeywordsRecognized(t *testing.T) {
	for _, p := range dispatchPrefixes {
		t.Run(p.prefix, func(t *testing.T) {
			sql := p.prefix + " " + p.tail // leading token is the keyword
			_, errs := Parse(sql)
			for _, e := range errs {
				if strings.Contains(e.Msg, "unknown or unsupported statement") {
					t.Errorf("Parse(%q): leading keyword %q hit the UNKNOWN branch; "+
						"it must be in the dispatch switch (got %q)", sql, p.prefix, e.Msg)
				}
			}
		})
	}
}

// TestParse_StatementLevelHintSkipped verifies a leading statement-level hint
// (@{...} / @N / @[N@]{...}) is skipped so dispatch sees the real statement
// keyword. Without the skip, `@{...} SELECT` would be reported as an unknown
// statement starting with '@'. Now that parser-select implements the query
// grammar, a hinted SELECT parses cleanly (zero diagnostics) and yields one
// statement node — the load-bearing assertion is that the hint was skipped (no
// "unknown statement starting with '@'" error).
func TestParse_StatementLevelHintSkipped(t *testing.T) {
	cases := []string{
		"@{USE_ADDITIONAL_PARALLELISM=TRUE} SELECT 1",
		"@{OPTIMIZER_VERSION=2, OPTIMIZER_STATISTICS_PACKAGE='auto'} SELECT * FROM T",
		"@5 SELECT 1",
		"@[5@]{key=1} SELECT 1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			file, errs := Parse(sql)
			if len(errs) != 0 {
				t.Fatalf("Parse(%q): got %d errors, want 0 (hinted SELECT parses): %v", sql, len(errs), errs)
			}
			if len(file.Stmts) != 1 {
				t.Errorf("Parse(%q): File.Stmts = %d, want 1 (the SELECT after the hint)", sql, len(file.Stmts))
			}
		})
	}
}

// TestParse_QueryStatementsParse documents that the parser-select, parser-dml,
// and parser-utility nodes flip the foundation's "all bodies stubbed" invariant:
// a valid SELECT, INSERT, and CALL all now parse to real AST nodes, while a
// still-stubbed segment (IMPORT MODULE) yields its "not yet supported"
// diagnostic. (Pre-parser-select this test asserted File.Stmts stayed empty for
// SELECT.)
func TestParse_QueryStatementsParse(t *testing.T) {
	file, errs := Parse("SELECT 1; INSERT INTO t VALUES (1); CALL p(); IMPORT MODULE a.b")
	if len(file.Stmts) != 3 {
		t.Errorf("File.Stmts = %d, want 3 (SELECT + INSERT + CALL parse; IMPORT is still stubbed)", len(file.Stmts))
	}
	if len(errs) != 1 {
		t.Errorf("errors = %d, want 1 (the stubbed IMPORT): %v", len(errs), errs)
	}
}

// TestParse_ProceduralBodyParsedAsOneSegment verifies the parse driver feeds a
// procedural BEGIN/END body to a single parseSingle call (block-aware Split): a
// stored procedure whose body has an inner ';' must be ONE segment, not split at
// the inner statement boundary. Now that the parser-ddl-bigquery node parses
// CREATE PROCEDURE, the whole block parses to exactly ONE statement with no
// errors — had Split cut at the inner ';', the trailing fragments
// (`SELECT 2`, `END`) would have failed to parse, so a clean single-statement
// parse is the proof the block stayed whole.
func TestParse_ProceduralBodyParsedAsOneSegment(t *testing.T) {
	res := ParseBestEffort("CREATE PROCEDURE p() BEGIN SELECT 1; SELECT 2; END")
	if len(res.Errors) != 0 {
		t.Fatalf("got %d errors, want 0 (block is one segment, CREATE PROCEDURE parses): %v", len(res.Errors), res.Errors)
	}
	if len(res.File.Stmts) != 1 {
		t.Fatalf("got %d statements, want 1 (the whole procedure)", len(res.File.Stmts))
	}
	if _, ok := res.File.Stmts[0].(*ast.CreateProcedureStmt); !ok {
		t.Errorf("statement is %T, want *ast.CreateProcedureStmt", res.File.Stmts[0])
	}
}
