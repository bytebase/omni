package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseIdent is a helper that constructs a Parser from input and
// calls parseIdent. Since parseIdent is unexported, tests in the same
// package can call it directly through this helper.
func testParseIdent(input string) (ast.Ident, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseIdent()
}

// testParseIdentStrict wraps parseIdentStrict the same way.
func testParseIdentStrict(input string) (ast.Ident, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseIdentStrict()
}

// testParseObjectName wraps parseObjectName the same way.
func testParseObjectName(input string) (*ast.ObjectName, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseObjectName()
}

// ---------------------------------------------------------------------------
// parseIdent tests
// ---------------------------------------------------------------------------

func TestParseIdent_Bare(t *testing.T) {
	id, err := testParseIdent("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || id.Quoted {
		t.Errorf("got %+v, want Ident{foo, unquoted}", id)
	}
}

func TestParseIdent_Quoted(t *testing.T) {
	id, err := testParseIdent(`"foo"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || !id.Quoted {
		t.Errorf("got %+v, want Ident{foo, quoted}", id)
	}
}

func TestParseIdent_QuotedWithSpaces(t *testing.T) {
	id, err := testParseIdent(`"my table"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "my table" || !id.Quoted {
		t.Errorf("got %+v, want Ident{my table, quoted}", id)
	}
}

func TestParseIdent_QuotedEscape(t *testing.T) {
	// F2's lexer un-escapes "" → " in Token.Str.
	id, err := testParseIdent(`"a""b"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != `a"b` || !id.Quoted {
		t.Errorf("got %+v, want Ident{a\"b, quoted}", id)
	}
}

func TestParseIdent_NonReservedKeyword(t *testing.T) {
	// ALERT is a non-reserved keyword — should be accepted as identifier.
	id, err := testParseIdent("ALERT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "ALERT" || id.Quoted {
		t.Errorf("got %+v, want Ident{ALERT, unquoted}", id)
	}
}

func TestParseIdent_ReservedKeyword(t *testing.T) {
	// SELECT is a reserved keyword — should be rejected.
	_, err := testParseIdent("SELECT")
	if err == nil {
		t.Fatal("expected error for reserved keyword SELECT, got nil")
	}
	if pe, ok := err.(*ParseError); ok {
		if !strings.Contains(pe.Msg, "expected identifier") {
			t.Errorf("error Msg = %q, want 'expected identifier'", pe.Msg)
		}
	} else {
		t.Errorf("expected *ParseError, got %T", err)
	}
}

func TestParseIdent_EOF(t *testing.T) {
	_, err := testParseIdent("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestParseIdent_LocCorrect(t *testing.T) {
	// "   foo" — the identifier starts at byte 3.
	id, err := testParseIdent("   foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Loc.Start != 3 || id.Loc.End != 6 {
		t.Errorf("Loc = %+v, want {3, 6}", id.Loc)
	}
}

// ---------------------------------------------------------------------------
// parseIdentStrict tests
// ---------------------------------------------------------------------------

func TestParseIdentStrict_AcceptsBare(t *testing.T) {
	id, err := testParseIdentStrict("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || id.Quoted {
		t.Errorf("got %+v, want Ident{foo, unquoted}", id)
	}
}

func TestParseIdentStrict_AcceptsQuoted(t *testing.T) {
	id, err := testParseIdentStrict(`"foo"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || !id.Quoted {
		t.Errorf("got %+v, want Ident{foo, quoted}", id)
	}
}

func TestParseIdentStrict_RejectsNonReservedKeyword(t *testing.T) {
	// parseIdentStrict rejects ALL keywords, even non-reserved ones.
	_, err := testParseIdentStrict("ALERT")
	if err == nil {
		t.Fatal("expected error for keyword ALERT in strict mode, got nil")
	}
}

// ---------------------------------------------------------------------------
// parseObjectName tests
// ---------------------------------------------------------------------------

func TestParseObjectName_OnePart(t *testing.T) {
	obj, err := testParseObjectName("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Name.Name != "foo" || !obj.Schema.IsEmpty() || !obj.Database.IsEmpty() {
		t.Errorf("got %+v, want 1-part name foo", obj)
	}
}

func TestParseObjectName_TwoParts(t *testing.T) {
	obj, err := testParseObjectName("sch.tbl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Schema.Name != "sch" || obj.Name.Name != "tbl" || !obj.Database.IsEmpty() {
		t.Errorf("got Schema=%q Name=%q DB=%q, want sch.tbl",
			obj.Schema.Name, obj.Name.Name, obj.Database.Name)
	}
	// Loc should span both parts.
	if obj.Loc.Start != 0 || obj.Loc.End != 7 {
		t.Errorf("Loc = %+v, want {0, 7}", obj.Loc)
	}
}

func TestParseObjectName_ThreeParts(t *testing.T) {
	obj, err := testParseObjectName("db.sch.tbl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Database.Name != "db" || obj.Schema.Name != "sch" || obj.Name.Name != "tbl" {
		t.Errorf("got DB=%q Schema=%q Name=%q, want db.sch.tbl",
			obj.Database.Name, obj.Schema.Name, obj.Name.Name)
	}
	if obj.Loc.Start != 0 || obj.Loc.End != 10 {
		t.Errorf("Loc = %+v, want {0, 10}", obj.Loc)
	}
}

func TestParseObjectName_Mixed(t *testing.T) {
	// "My DB".sch."Quoted Table"
	obj, err := testParseObjectName(`"My DB".sch."Quoted Table"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !obj.Database.Quoted || obj.Database.Name != "My DB" {
		t.Errorf("Database = %+v, want quoted 'My DB'", obj.Database)
	}
	if obj.Schema.Quoted || obj.Schema.Name != "sch" {
		t.Errorf("Schema = %+v, want unquoted 'sch'", obj.Schema)
	}
	if !obj.Name.Quoted || obj.Name.Name != "Quoted Table" {
		t.Errorf("Name = %+v, want quoted 'Quoted Table'", obj.Name)
	}
}

func TestParseObjectName_TrailingDot(t *testing.T) {
	// "foo." — error after the dot (no identifier follows).
	_, err := testParseObjectName("foo.")
	if err == nil {
		t.Fatal("expected error for trailing dot, got nil")
	}
}

func TestParseObjectName_LeadingDot(t *testing.T) {
	// ".foo" — the dot is not an identifier.
	_, err := testParseObjectName(".foo")
	if err == nil {
		t.Fatal("expected error for leading dot, got nil")
	}
}

func TestParseObjectName_DoubleDot(t *testing.T) {
	// "foo..bar" — error at the second dot.
	_, err := testParseObjectName("foo..bar")
	if err == nil {
		t.Fatal("expected error for double dot, got nil")
	}
}

func TestParseObjectName_Freestanding(t *testing.T) {
	// ParseObjectName is the public freestanding helper.
	obj, errs := ParseObjectName("db.sch.tbl")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if obj.Database.Name != "db" || obj.Schema.Name != "sch" || obj.Name.Name != "tbl" {
		t.Errorf("ParseObjectName mismatch: %+v", obj)
	}
}

func TestParseObjectName_FreestandingTrailingTokens(t *testing.T) {
	// "foo bar" — parses "foo" then complains about trailing "bar".
	obj, errs := ParseObjectName("foo bar")
	if obj == nil {
		t.Fatal("expected non-nil ObjectName even with trailing tokens")
	}
	if obj.Name.Name != "foo" {
		t.Errorf("Name = %q, want foo", obj.Name.Name)
	}
	if len(errs) != 1 || !strings.Contains(errs[0].Msg, "unexpected token") {
		t.Errorf("expected 1 'unexpected token' error, got %+v", errs)
	}
}

func TestParseObjectName_FreestandingEmpty(t *testing.T) {
	_, errs := ParseObjectName("")
	if len(errs) == 0 {
		t.Error("expected error for empty input, got none")
	}
}

func TestParseObjectName_KeywordAsPart(t *testing.T) {
	// Non-reserved keyword ALERT should be usable as any part.
	obj, err := testParseObjectName("ALERT.ALERT.ALERT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Database.Name != "ALERT" || obj.Schema.Name != "ALERT" || obj.Name.Name != "ALERT" {
		t.Errorf("got %+v, want ALERT.ALERT.ALERT", obj)
	}
}
