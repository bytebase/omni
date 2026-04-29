package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func newFoundationTestParser(sql string) *Parser {
	p := &Parser{
		lexer:  NewLexer(sql),
		source: sql,
	}
	p.advance()
	return p
}

func TestFoundationSyntaxErrorHelpers(t *testing.T) {
	t.Run("current token", func(t *testing.T) {
		p := newFoundationTestParser("SELECT 1")
		err := p.syntaxErrorAtCur()
		if err == nil {
			t.Fatal("syntaxErrorAtCur returned nil")
		}
		if err.Position != 0 {
			t.Fatalf("Position = %d, want 0", err.Position)
		}
		if !strings.Contains(err.Error(), `syntax error at or near "SELECT"`) {
			t.Fatalf("error = %q, want syntax error at current token", err.Error())
		}
	})

	t.Run("EOF", func(t *testing.T) {
		p := newFoundationTestParser("")
		err := p.syntaxErrorAtCur()
		if err == nil {
			t.Fatal("syntaxErrorAtCur returned nil")
		}
		if err.Position != 0 {
			t.Fatalf("Position = %d, want 0", err.Position)
		}
		if !strings.Contains(err.Error(), "syntax error at end of input") {
			t.Fatalf("error = %q, want EOF syntax error", err.Error())
		}
	})
}

func TestFoundationExpectUsesSyntaxError(t *testing.T) {
	p := newFoundationTestParser("")
	_, err := p.expect(kwSELECT)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "syntax error at end of input") {
		t.Fatalf("error = %q, want EOF syntax error", err.Error())
	}
}

func TestFoundationLexerErrorsPropagate(t *testing.T) {
	cases := []string{
		"SELECT 'unterminated",
		"SELECT \"unterminated",
		"SELECT q'[unterminated",
		"SELECT /* unterminated",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want lexer error", sql)
			}
			if !strings.Contains(err.Error(), "unterminated") {
				t.Fatalf("error = %q, want lexer error context", err.Error())
			}
		})
	}
}

func TestFoundationStatementSeparatorRequired(t *testing.T) {
	_, err := Parse("SELECT 1 SELECT 2")
	if err == nil {
		t.Fatal("expected error for adjacent statements without separator")
	}
	if !strings.Contains(err.Error(), `syntax error at or near "SELECT"`) {
		t.Fatalf("error = %q, want syntax error at second SELECT", err.Error())
	}
}

func TestFoundationLocEndUsesConsumedTokenEnd(t *testing.T) {
	sql := "SELECT /*+ FULL(t) */ * FROM t"
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SelectStmt)
	if stmt.Loc.End != len(sql) {
		t.Fatalf("SelectStmt Loc.End = %d, want %d", stmt.Loc.End, len(sql))
	}
	if stmt.Hints == nil || stmt.Hints.Len() != 1 {
		t.Fatalf("expected one hint, got %#v", stmt.Hints)
	}
	hint := stmt.Hints.Items[0].(*ast.Hint)
	if hint.Loc.End <= hint.Loc.Start {
		t.Fatalf("Hint Loc = %+v, want non-empty span", hint.Loc)
	}
	if got := sql[hint.Loc.Start:hint.Loc.End]; !strings.HasPrefix(got, "/*+") {
		t.Fatalf("hint Loc text = %q, want hint comment span", got)
	}
}

func TestFoundationTokenEndSpans(t *testing.T) {
	cases := []struct {
		sql  string
		typ  int
		text string
	}{
		{"'abc'", tokSCONST, "'abc'"},
		{`"MixedCase"`, tokQIDENT, `"MixedCase"`},
		{"/*+ FULL(t) */ SELECT 1", tokHINT, "/*+ FULL(t) */"},
	}
	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			tok := NewLexer(tc.sql).NextToken()
			if tok.Type != tc.typ {
				t.Fatalf("token type = %d, want %d", tok.Type, tc.typ)
			}
			if tok.End < tok.Loc {
				t.Fatalf("token span = [%d,%d), want End >= Loc", tok.Loc, tok.End)
			}
			if got := tc.sql[tok.Loc:tok.End]; got != tc.text {
				t.Fatalf("token span text = %q, want %q", got, tc.text)
			}
		})
	}
}

func TestFoundationMultiStatementRawSpans(t *testing.T) {
	sql := "SELECT 1; SELECT 2"
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if result.Len() != 2 {
		t.Fatalf("statement count = %d, want 2", result.Len())
	}

	first := result.Items[0].(*ast.RawStmt)
	if first.Loc != (ast.Loc{Start: 0, End: len("SELECT 1")}) {
		t.Fatalf("first RawStmt Loc = %+v, want [0,%d)", first.Loc, len("SELECT 1"))
	}

	second := result.Items[1].(*ast.RawStmt)
	if second.Loc != (ast.Loc{Start: len("SELECT 1; "), End: len(sql)}) {
		t.Fatalf("second RawStmt Loc = %+v, want [%d,%d)", second.Loc, len("SELECT 1; "), len(sql))
	}
}
