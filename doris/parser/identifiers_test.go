package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// makeParser creates a Parser primed on the given input, ready to parse.
func makeParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance() // prime cur
	return p
}

// ---------------------------------------------------------------------------
// parseIdentifier
// ---------------------------------------------------------------------------

func TestParseIdentifier_Unquoted(t *testing.T) {
	p := makeParser("myTable")
	name, loc, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "myTable" {
		t.Errorf("got %q, want %q", name, "myTable")
	}
	if loc.Start != 0 || loc.End != 7 {
		t.Errorf("loc = %v, want {0,7}", loc)
	}
	// parser should now be at EOF
	if p.cur.Kind != tokEOF {
		t.Errorf("cur after parse = %d, want EOF", p.cur.Kind)
	}
}

func TestParseIdentifier_BacktickQuoted(t *testing.T) {
	p := makeParser("`my table`")
	name, _, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// lexer strips backticks, so we get the raw content
	if name != "my table" {
		t.Errorf("got %q, want %q", name, "my table")
	}
}

func TestParseIdentifier_BacktickPreservesCase(t *testing.T) {
	p := makeParser("`MyMixedCase`")
	name, _, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "MyMixedCase" {
		t.Errorf("got %q, want %q", name, "MyMixedCase")
	}
}

func TestParseIdentifier_NonReservedKeyword(t *testing.T) {
	// COMMENT is non-reserved in Doris and should parse as identifier.
	p := makeParser("comment")
	name, _, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "comment" {
		t.Errorf("got %q, want %q", name, "comment")
	}
}

func TestParseIdentifier_ReservedKeywordFails(t *testing.T) {
	// SELECT is reserved and must not parse as an identifier.
	p := makeParser("select")
	_, _, err := p.parseIdentifier()
	if err == nil {
		t.Fatal("expected error for reserved keyword SELECT, got nil")
	}
}

func TestParseIdentifier_EOF_Fails(t *testing.T) {
	p := makeParser("")
	_, _, err := p.parseIdentifier()
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

// ---------------------------------------------------------------------------
// parseMultipartIdentifier
// ---------------------------------------------------------------------------

func TestParseMultipartIdentifier_Single(t *testing.T) {
	p := makeParser("myTable")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(name.Parts) != 1 {
		t.Fatalf("got %d parts, want 1", len(name.Parts))
	}
	if name.Parts[0] != "myTable" {
		t.Errorf("Parts[0] = %q, want %q", name.Parts[0], "myTable")
	}
}

func TestParseMultipartIdentifier_TwoPart(t *testing.T) {
	p := makeParser("db.table")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"db", "table"}
	if len(name.Parts) != len(want) {
		t.Fatalf("got %d parts, want %d: %v", len(name.Parts), len(want), name.Parts)
	}
	for i, w := range want {
		if name.Parts[i] != w {
			t.Errorf("Parts[%d] = %q, want %q", i, name.Parts[i], w)
		}
	}
}

func TestParseMultipartIdentifier_ThreePart(t *testing.T) {
	p := makeParser("catalog.db.table")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"catalog", "db", "table"}
	if len(name.Parts) != 3 {
		t.Fatalf("got %d parts, want 3: %v", len(name.Parts), name.Parts)
	}
	for i, w := range want {
		if name.Parts[i] != w {
			t.Errorf("Parts[%d] = %q, want %q", i, name.Parts[i], w)
		}
	}
}

func TestParseMultipartIdentifier_QuotedSingle(t *testing.T) {
	p := makeParser("`my table`")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(name.Parts) != 1 || name.Parts[0] != "my table" {
		t.Errorf("got Parts=%v, want [\"my table\"]", name.Parts)
	}
}

func TestParseMultipartIdentifier_MixedQuoted(t *testing.T) {
	// db.`my table`
	p := makeParser("db.`my table`")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"db", "my table"}
	if len(name.Parts) != 2 {
		t.Fatalf("got %d parts, want 2: %v", len(name.Parts), name.Parts)
	}
	for i, w := range want {
		if name.Parts[i] != w {
			t.Errorf("Parts[%d] = %q, want %q", i, name.Parts[i], w)
		}
	}
}

func TestParseMultipartIdentifier_NonReservedKeyword(t *testing.T) {
	// COMMENT is non-reserved and valid as a table name.
	p := makeParser("db.comment")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"db", "comment"}
	if len(name.Parts) != 2 {
		t.Fatalf("got %d parts, want 2: %v", len(name.Parts), name.Parts)
	}
	for i, w := range want {
		if name.Parts[i] != w {
			t.Errorf("Parts[%d] = %q, want %q", i, name.Parts[i], w)
		}
	}
}

func TestParseMultipartIdentifier_DotNotFollowedByIdent(t *testing.T) {
	// "mytable.(" — the '.' is followed by '(' which is not an identifier,
	// so we stop at the first part and leave '.(' for the caller.
	// (Note: "mytable.123" would not work because the lexer folds ".123" into
	// a tokFloat when it immediately follows an identifier.)
	p := makeParser("mytable.(")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(name.Parts) != 1 || name.Parts[0] != "mytable" {
		t.Errorf("got Parts=%v, want [\"mytable\"]", name.Parts)
	}
	// The '.' should still be in the stream.
	if p.cur.Kind != int('.') {
		t.Errorf("cur after partial parse = %d, want '.'", p.cur.Kind)
	}
}

func TestParseMultipartIdentifier_ReservedKeywordFails(t *testing.T) {
	// SELECT is reserved — parsing it as an identifier must fail.
	p := makeParser("select")
	_, err := p.parseMultipartIdentifier()
	if err == nil {
		t.Fatal("expected error for reserved keyword SELECT, got nil")
	}
}

func TestParseMultipartIdentifier_LocSpan(t *testing.T) {
	// "db.table": start=0, end=8
	p := makeParser("db.table")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", name.Loc.Start)
	}
	if name.Loc.End != 8 {
		t.Errorf("Loc.End = %d, want 8", name.Loc.End)
	}
}

func TestParseMultipartIdentifier_String(t *testing.T) {
	p := makeParser("catalog.db.table")
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name.String() != "catalog.db.table" {
		t.Errorf("String() = %q, want %q", name.String(), "catalog.db.table")
	}
}

// ---------------------------------------------------------------------------
// parseIdentifierOrString
// ---------------------------------------------------------------------------

func TestParseIdentifierOrString_Identifier(t *testing.T) {
	p := makeParser("mydb")
	val, _, err := p.parseIdentifierOrString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "mydb" {
		t.Errorf("got %q, want %q", val, "mydb")
	}
}

func TestParseIdentifierOrString_StringLiteral(t *testing.T) {
	p := makeParser("'mydb'")
	val, _, err := p.parseIdentifierOrString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "mydb" {
		t.Errorf("got %q, want %q", val, "mydb")
	}
}

func TestParseIdentifierOrString_DoubleQuotedString(t *testing.T) {
	p := makeParser(`"mydb"`)
	val, _, err := p.parseIdentifierOrString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "mydb" {
		t.Errorf("got %q, want %q", val, "mydb")
	}
}

// ---------------------------------------------------------------------------
// NormalizeIdentifier
// ---------------------------------------------------------------------------

func TestNormalizeIdentifier_UnquotedLowercases(t *testing.T) {
	got := NormalizeIdentifier("MyTable", false)
	if got != "mytable" {
		t.Errorf("NormalizeIdentifier(%q, false) = %q, want %q", "MyTable", got, "mytable")
	}
}

func TestNormalizeIdentifier_QuotedPreservesCase(t *testing.T) {
	got := NormalizeIdentifier("MyTable", true)
	if got != "MyTable" {
		t.Errorf("NormalizeIdentifier(%q, true) = %q, want %q", "MyTable", got, "MyTable")
	}
}

func TestNormalizeIdentifier_AlreadyLower(t *testing.T) {
	got := NormalizeIdentifier("mytable", false)
	if got != "mytable" {
		t.Errorf("NormalizeIdentifier(%q, false) = %q, want %q", "mytable", got, "mytable")
	}
}

func TestNormalizeIdentifier_AllUpper(t *testing.T) {
	got := NormalizeIdentifier("MYTABLE", false)
	if got != "mytable" {
		t.Errorf("NormalizeIdentifier(%q, false) = %q, want %q", "MYTABLE", got, "mytable")
	}
}

// ---------------------------------------------------------------------------
// NormalizeObjectName
// ---------------------------------------------------------------------------

func TestNormalizeObjectName_SinglePart(t *testing.T) {
	name := &ast.ObjectName{Parts: []string{"MyTable"}}
	got := NormalizeObjectName(name)
	if got != "mytable" {
		t.Errorf("NormalizeObjectName = %q, want %q", got, "mytable")
	}
}

func TestNormalizeObjectName_MultiPart(t *testing.T) {
	name := &ast.ObjectName{Parts: []string{"MyCatalog", "MyDB", "MyTable"}}
	got := NormalizeObjectName(name)
	if got != "mycatalog.mydb.mytable" {
		t.Errorf("NormalizeObjectName = %q, want %q", got, "mycatalog.mydb.mytable")
	}
}

func TestNormalizeObjectName_Nil(t *testing.T) {
	got := NormalizeObjectName(nil)
	if got != "" {
		t.Errorf("NormalizeObjectName(nil) = %q, want %q", got, "")
	}
}

func TestNormalizeObjectName_AlreadyLower(t *testing.T) {
	name := &ast.ObjectName{Parts: []string{"catalog", "db", "table"}}
	got := NormalizeObjectName(name)
	if got != "catalog.db.table" {
		t.Errorf("NormalizeObjectName = %q, want %q", got, "catalog.db.table")
	}
}
