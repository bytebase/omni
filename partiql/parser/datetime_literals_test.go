package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestDateTimeLiteral_Date covers the DATE LITERAL_STRING form
// (PartiQLParser.g4 line 670, literal#LiteralDate).
//
// The grammar accepts any string body after DATE — value validation
// (real calendar dates) is a semantic concern, not a syntax one, and
// the legacy ANTLR parser does not validate it either. So `DATE
// 'garbage'` parses syntactically; we assert that to match antlr.
func TestDateTimeLiteral_Date(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantVal string
	}{
		{"iso", "DATE '2026-01-01'", "2026-01-01"},
		{"min", "DATE '0001-01-01'", "0001-01-01"},
		{"leap", "DATE '2024-02-29'", "2024-02-29"},
		// Grammar-accepted, semantically-odd bodies (following antlr: no
		// syntax-layer value validation).
		{"empty_body", "DATE ''", ""},
		{"non_date_body", "DATE 'not-a-date'", "not-a-date"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("ParseExpr(%q) error: %v", tc.input, err)
			}
			lit, ok := expr.(*ast.DateLit)
			if !ok {
				t.Fatalf("ParseExpr(%q) = %T, want *ast.DateLit", tc.input, expr)
			}
			if lit.Val != tc.wantVal {
				t.Errorf("DateLit.Val = %q, want %q", lit.Val, tc.wantVal)
			}
			// Loc must span from the DATE keyword (offset 0) through the
			// closing quote of the string body.
			if lit.Loc.Start != 0 {
				t.Errorf("DateLit.Loc.Start = %d, want 0", lit.Loc.Start)
			}
			if lit.Loc.End != len(tc.input) {
				t.Errorf("DateLit.Loc.End = %d, want %d", lit.Loc.End, len(tc.input))
			}
		})
	}
}

// TestDateTimeLiteral_Time covers the full TIME literal form
// (PartiQLParser.g4 line 671, literal#LiteralTime):
//
//	TIME ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING
func TestDateTimeLiteral_Time(t *testing.T) {
	p3 := 3
	p6 := 6
	p0 := 0
	cases := []struct {
		name      string
		input     string
		wantVal   string
		wantPrec  *int
		wantWithZ bool
	}{
		{"plain", "TIME '12:00:00'", "12:00:00", nil, false},
		{"frac", "TIME '12:00:00.123'", "12:00:00.123", nil, false},
		{"precision", "TIME (3) '12:00:00.123'", "12:00:00.123", &p3, false},
		{"precision_no_space", "TIME(6) '23:59:59.999999'", "23:59:59.999999", &p6, false},
		{"precision_zero", "TIME(0) '01:02:03'", "01:02:03", &p0, false},
		{"with_tz", "TIME WITH TIME ZONE '12:00:00'", "12:00:00", nil, true},
		{"precision_with_tz", "TIME (3) WITH TIME ZONE '12:00:00-08:00'", "12:00:00-08:00", &p3, true},
		// Grammar-accepted odd body (no syntax-layer validation).
		{"empty_body", "TIME ''", "", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("ParseExpr(%q) error: %v", tc.input, err)
			}
			lit, ok := expr.(*ast.TimeLit)
			if !ok {
				t.Fatalf("ParseExpr(%q) = %T, want *ast.TimeLit", tc.input, expr)
			}
			if lit.Val != tc.wantVal {
				t.Errorf("TimeLit.Val = %q, want %q", lit.Val, tc.wantVal)
			}
			switch {
			case tc.wantPrec == nil && lit.Precision != nil:
				t.Errorf("TimeLit.Precision = %d, want nil", *lit.Precision)
			case tc.wantPrec != nil && lit.Precision == nil:
				t.Errorf("TimeLit.Precision = nil, want %d", *tc.wantPrec)
			case tc.wantPrec != nil && lit.Precision != nil && *lit.Precision != *tc.wantPrec:
				t.Errorf("TimeLit.Precision = %d, want %d", *lit.Precision, *tc.wantPrec)
			}
			if lit.WithTimeZone != tc.wantWithZ {
				t.Errorf("TimeLit.WithTimeZone = %v, want %v", lit.WithTimeZone, tc.wantWithZ)
			}
			if lit.Loc.Start != 0 {
				t.Errorf("TimeLit.Loc.Start = %d, want 0", lit.Loc.Start)
			}
			if lit.Loc.End != len(tc.input) {
				t.Errorf("TimeLit.Loc.End = %d, want %d", lit.Loc.End, len(tc.input))
			}
		})
	}
}

// TestDateTimeLiteral_Reject covers the negative/reject cases required
// by the correctness protocol. These are forms the grammar does NOT
// accept; the parser must reject each with a syntax error.
func TestDateTimeLiteral_Reject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// DATE with no string body.
		{"date_no_string", "DATE", "expected string literal"},
		// DATE followed by a non-string token.
		{"date_number_body", "DATE 123", "expected string literal"},
		{"date_ident_body", "DATE foo", "expected string literal"},
		// TIME with no string body.
		{"time_no_string", "TIME", "expected string literal"},
		// TIME precision must be an integer.
		{"time_precision_non_int", "TIME ('x') '12:00:00'", "expected integer"},
		// TIME precision paren must close.
		{"time_precision_unclosed", "TIME (3 '12:00:00'", "expected"},
		// WITH must be followed by TIME ZONE.
		{"time_with_no_zone", "TIME WITH TIME '12:00:00'", "expected ZONE"},
		// TIMESTAMP literal form is NOT part of PartiQL (oracle-confirmed:
		// canonical grammar + AWS DynamoDB). It must be rejected.
		{"timestamp_literal", "TIMESTAMP '2026-01-01 00:00:00'", "TIMESTAMP literal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) succeeded, want error containing %q", tc.input, tc.wantErrIn)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("ParseExpr(%q) error = %q, want substring %q", tc.input, err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestDateTimeLiteral_TrailingLexError asserts that a DATE/TIME literal
// followed by un-lexable trailing content is REJECTED, not silently
// accepted (Codex review finding on PR #181).
//
// Mechanism: the string body is the final token of every DATE/TIME form,
// so the advance() that consumes it is where the lexer first scans
// whatever follows. If that trailing content is itself un-lexable — an
// unrecognized character (`#`) or an unterminated block comment (`/*`) —
// the lexer's first-error-and-stop contract sets Lexer.Err and yields
// tokEOF instead of the offending token. Before the fix, ParseExpr's
// post-parse `cur == tokEOF` check mistook this for a clean end and
// returned the DateLit/TimeLit as success, dropping the lexer error.
// parseDateLiteral/parseTimeLiteral now surface the pending lexer error
// (via expectStringLiteral -> checkLexerErr) so the parse fails.
//
// Oracle (antlr_fallback): PartiQLLexer.g4 has no standalone `#` token —
// `#` matches only the catch-all `UNRECOGNIZED : . ;` rule, and an
// unterminated `/*` cannot match COMMENT_BLOCK (`'/*' .*? '*/'`). Either
// way ANTLR yields trailing input the `literal` rule cannot consume, so
// the canonical grammar likewise REJECTS these inputs. omni surfaces the
// rejection one layer earlier (lexer error vs ANTLR's extra-token parser
// error), but the accept/reject verdict matches: both reject.
func TestDateTimeLiteral_TrailingLexError(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// The two exact inputs from the Codex finding.
		{"date_trailing_hash", "DATE '2026-01-01' #", "unexpected character '#'"},
		{"time_trailing_block_comment", "TIME '12:00:00' /*", "unterminated block comment"},
		// Same hazard on the other end (DATE + unterminated comment,
		// TIME + unrecognized char) and on the precision/with-tz TIME
		// variants, which all funnel through expectStringLiteral.
		{"date_trailing_block_comment", "DATE '2026-01-01' /*", "unterminated block comment"},
		{"time_trailing_hash", "TIME '12:00:00' #", "unexpected character '#'"},
		{"time_prec_trailing_hash", "TIME (3) '12:00:00.123' #", "unexpected character '#'"},
		{"time_tz_trailing_block_comment", "TIME WITH TIME ZONE '12:00:00' /*", "unterminated block comment"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) = %T, want error containing %q (trailing lexer error must not be swallowed)", tc.input, expr, tc.wantErrIn)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("ParseExpr(%q) error = %q, want substring %q", tc.input, err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestDateTimeLiteral_TrailingValidTokenStillRejected guards the
// neighbouring case so the trailing-lex-error fix is not mistaken for the
// EOF check: a DATE/TIME literal followed by a WELL-FORMED extra token
// (which leaves Lexer.Err nil) must still be rejected, by the existing
// "unexpected token after expression" EOF check rather than by the new
// lexer-error surfacing. This pins both rejection paths.
func TestDateTimeLiteral_TrailingValidTokenStillRejected(t *testing.T) {
	cases := []string{
		"DATE '2026-01-01' 99",
		"TIME '12:00:00' foo",
		"DATE '2026-01-01' @",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			p := NewParser(in)
			expr, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) = %T, want error (trailing token after literal)", in, expr)
			}
			if !strings.Contains(err.Error(), "after expression") {
				t.Errorf("ParseExpr(%q) error = %q, want the EOF-check message %q", in, err.Error(), "after expression")
			}
		})
	}
}

// TestDateTimeLiteral_Goldens iterates every .partiql file under
// testdata/datetime/ and compares the parser's pretty-printed output
// (via ast.NodeToString) against the matching .golden file.
//
// Run with `go test -update -run TestDateTimeLiteral_Goldens ./partiql/parser/...`
// to regenerate goldens after intentional AST shape changes.
func TestDateTimeLiteral_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/datetime/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/datetime/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got+"\n" != string(want) {
				t.Errorf("AST mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}
