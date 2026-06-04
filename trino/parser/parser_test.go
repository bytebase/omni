package parser

import (
	"strings"
	"testing"
)

// TestParse_EmptyInput verifies that parsing empty or whitespace-only input
// yields a non-nil File with no statements and no errors.
func TestParse_EmptyInput(t *testing.T) {
	for _, in := range []string{"", "   ", "\n\t ", "-- just a comment\n", "/* block */"} {
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

// TestParse_UnknownStatement verifies a statement starting with a token the
// dispatch switch does not recognize produces a parse error (not a panic, not
// silent acceptance).
func TestParse_UnknownStatement(t *testing.T) {
	file, errs := Parse("FOOBAR 1 2 3")
	if file == nil {
		t.Fatal("Parse: File is nil")
	}
	if len(errs) == 0 {
		t.Fatal("Parse(\"FOOBAR ...\"): want a parse error, got none")
	}
	if !strings.Contains(errs[0].Msg, "unknown or unsupported statement") {
		t.Errorf("error = %q, want it to mention 'unknown or unsupported statement'", errs[0].Msg)
	}
}

// TestParse_KnownStatementUnsupported verifies that a statement whose leading
// keyword IS in the dispatch switch but whose body is STILL stubbed (no DAG node
// has implemented it yet) yields a "not yet supported" diagnostic rather than an
// "unknown statement" one. INSERT (parser-dml's job) is still stubbed; SELECT and
// the rest of the query layer are now implemented by parser-select, so this test
// uses a statement that remains a foundation stub.
func TestParse_KnownStatementUnsupported(t *testing.T) {
	file, errs := Parse("INSERT INTO t VALUES (1)")
	if file == nil {
		t.Fatal("Parse: File is nil")
	}
	if len(errs) != 1 {
		t.Fatalf("Parse(\"INSERT ...\"): got %d errors, want 1: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "not yet supported") {
		t.Errorf("error = %q, want it to mention 'not yet supported'", errs[0].Msg)
	}
}

// TestParse_MultiStatementErrorsCollected verifies that ParseBestEffort
// collects errors from every segment, not just the first. Two still-stubbed
// statements (INSERT and DELETE, parser-dml's job) separated by ';' must yield
// two errors.
func TestParse_MultiStatementErrorsCollected(t *testing.T) {
	res := ParseBestEffort("INSERT INTO t VALUES (1); DELETE FROM u")
	if got := len(res.Errors); got != 2 {
		t.Fatalf("ParseBestEffort: got %d errors, want 2: %v", got, res.Errors)
	}
}

// TestParse_LexErrorPromoted verifies that a lexer-level failure (unterminated
// string) surfaces as a parse error with a position, so Diagnose can report it.
func TestParse_LexErrorPromoted(t *testing.T) {
	_, errs := Parse("SELECT 'unterminated")
	if len(errs) == 0 {
		t.Fatal("Parse with unterminated string: want a (lex-promoted) error, got none")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Msg, "unterminated string") {
			found = true
		}
	}
	if !found {
		t.Errorf("errors = %v, want one mentioning 'unterminated string'", errs)
	}
}

// TestParse_StrictReturnsAllErrors verifies the strict Parse entry returns the
// full error slice (matching doris's Parse signature, which returns
// []ParseError rather than a single error).
func TestParse_StrictReturnsAllErrors(t *testing.T) {
	_, errs := Parse("FOOBAR; BAZBAR")
	if len(errs) != 2 {
		t.Fatalf("Parse: got %d errors, want 2: %v", len(errs), errs)
	}
}

// TestParse_DispatchKeywordsRecognized verifies that every documented top-level
// statement keyword is routed by the dispatch switch — i.e. it produces a
// "not yet supported" (stubbed) error, NOT an "unknown statement" error. This
// is the foundation's completeness contract: all 81 statement forms must be
// reachable for Diagnose to avoid false "unknown statement" diagnostics.
func TestParse_DispatchKeywordsRecognized(t *testing.T) {
	// One representative prefix per first-keyword in the grammar's `statement`
	// and `rootQuery`/`queryPrimary` rules. The body after the keyword does not
	// matter — the foundation stubs the body — so we only need the leading
	// keyword(s) to land in a real dispatch case rather than the default.
	prefixes := []string{
		"SELECT", "WITH", "TABLE", "VALUES", "(",
		"USE", "CREATE", "DROP", "ALTER", "INSERT", "DELETE", "UPDATE",
		"MERGE", "TRUNCATE", "COMMENT", "ANALYZE", "REFRESH", "CALL",
		"GRANT", "REVOKE", "DENY", "SET", "RESET", "SHOW", "EXPLAIN",
		"DESCRIBE", "DESC", "START", "COMMIT", "ROLLBACK", "PREPARE",
		"DEALLOCATE", "EXECUTE",
	}
	for _, p := range prefixes {
		// Use a minimal-but-plausible tail so the leading token is the keyword.
		// EXPLAIN is special: it is now a real parser (parser-utility node) that
		// recurses into the inner statement, so `EXPLAIN x` would route the inner
		// bare `x` to the default branch (correctly — Trino also rejects
		// `EXPLAIN x`). Probe it with an inner keyword that is itself dispatched
		// so the assertion checks EXPLAIN's own dispatch, not the inner token.
		sql := p + " x"
		if p == "EXPLAIN" {
			sql = "EXPLAIN SELECT 1"
		}
		_, errs := Parse(sql)
		for _, e := range errs {
			if strings.Contains(e.Msg, "unknown or unsupported statement") {
				t.Errorf("keyword %q routed to default (unknown) dispatch: %q -> %q", p, sql, e.Msg)
			}
		}
	}
}
