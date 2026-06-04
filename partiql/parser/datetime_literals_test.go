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
