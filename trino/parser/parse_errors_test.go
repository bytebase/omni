package parser

import (
	"errors"
	"testing"
)

// TestParse_ReturnsParseErrors pins the error contract every omni engine
// shares: Parse returns nil on success, otherwise an error that errors.As
// resolves to the first *ParseError with a byte Position, while AllErrors
// still yields every error a multi-statement script produced.
func TestParse_ReturnsParseErrors(t *testing.T) {
	if _, err := Parse("SELECT 1"); err != nil {
		t.Fatalf("Parse(valid) = %v, want nil", err)
	}

	file, err := Parse("SELECT 1 ))); SELECT 2; SELECT 3 )))")
	if err == nil {
		t.Fatal("Parse(two bad segments) returned nil error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("errors.As(*ParseError) failed on %T", err)
	}
	if pe.Position != len("SELECT 1 ") {
		t.Errorf("first error Position = %d, want %d (the first ')')", pe.Position, len("SELECT 1 "))
	}
	if pe.Message == "" {
		t.Error("first error has an empty Message")
	}
	if got := AllErrors(err); len(got) < 2 {
		t.Errorf("AllErrors = %d errors, want at least 2 (one per bad segment): %v", len(got), got)
	}
	if file == nil || len(file.Stmts) != 1 {
		t.Errorf("File keeps the statements that parsed: got %v", file)
	}
	if AllErrors(nil) != nil || AllErrors(errors.New("other")) != nil {
		t.Error("AllErrors must be nil for nil and for non-parse errors")
	}
}

// TestParseBestEffort_ObjectKeywordDrainsSegment pins what an invalid CREATE or
// DROP object keyword leaves behind in the best-effort parse: the syntax error
// at that keyword, and the lexical errors of the rest of the segment, which the
// lazy lexer reaches only because the rejection drains it. Strict Parse lexes
// the whole input itself, so it reports the same two.
func TestParseBestEffort_ObjectKeywordDrainsSegment(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"CREATE INDEX 'unterminated", []string{"syntax error at or near INDEX", "unterminated string literal"}},
		{"DROP FOO 'unterminated", []string{"syntax error at or near FOO", "unterminated string literal"}},
	} {
		res := ParseBestEffort(tc.sql)
		if len(res.Errors) != len(tc.want) {
			t.Fatalf("ParseBestEffort(%q) = %v, want %v", tc.sql, res.Errors, tc.want)
		}
		for i, want := range tc.want {
			if res.Errors[i].Message != want {
				t.Errorf("ParseBestEffort(%q) error %d = %q, want %q", tc.sql, i, res.Errors[i].Message, want)
			}
		}
		if _, err := Parse(tc.sql); len(AllErrors(err)) != len(tc.want) {
			t.Errorf("Parse(%q) = %v, want %v", tc.sql, AllErrors(err), tc.want)
		}
	}
}

// TestParseBestEffort_KeepsExplainedNodeWithDeferredCheck pins the node a
// deferred check (see Parser.deferred) leaves behind: the explained statement is
// invalid, and best effort still returns the EXPLAIN node around it, as it does
// for an inner parse error.
func TestParseBestEffort_KeepsExplainedNodeWithDeferredCheck(t *testing.T) {
	for _, sql := range []string{
		"EXPLAIN SELECT TRIM(FROM 1.5)",
		"EXPLAIN GRANT SELECT TO a",
	} {
		res := ParseBestEffort(sql)
		if len(res.Errors) != 1 {
			t.Fatalf("ParseBestEffort(%q) errors = %v, want 1", sql, res.Errors)
		}
		if len(res.File.Stmts) != 1 {
			t.Fatalf("ParseBestEffort(%q) statements = %v, want the EXPLAIN node", sql, res.File.Stmts)
		}
		if _, ok := res.File.Stmts[0].(*ExplainStmt); !ok {
			t.Errorf("ParseBestEffort(%q) statement is %T, want *ExplainStmt", sql, res.File.Stmts[0])
		}
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want the deferred check's", sql)
		}
	}
}
