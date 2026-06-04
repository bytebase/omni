package parser

import (
	"testing"

	"github.com/bytebase/omni/trino/ast"
)

// parseIdentFrom is a tiny harness: it lexes input, primes a Parser, and runs
// the given parse method, returning the result and the parser so callers can
// inspect the trailing token.
func newParserForTest(input string) *Parser {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p
}

// TestParseIdentifier_Bare verifies an unquoted identifier parses with Quoted
// false and case preserved.
func TestParseIdentifier_Bare(t *testing.T) {
	p := newParserForTest("MyCol")
	id, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("parseIdentifier: %v", err)
	}
	if id.Value != "MyCol" || id.Quoted {
		t.Errorf("got {Value:%q Quoted:%v}, want {MyCol false}", id.Value, id.Quoted)
	}
	if id.Loc.Start != 0 || id.Loc.End != 5 {
		t.Errorf("Loc = %+v, want [0,5)", id.Loc)
	}
}

// TestParseIdentifier_DoubleQuoted verifies a "double-quoted" identifier parses
// with Quoted true, QuoteRune '"', and the surrounding quotes stripped.
func TestParseIdentifier_DoubleQuoted(t *testing.T) {
	p := newParserForTest(`"My Col"`)
	id, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("parseIdentifier: %v", err)
	}
	if id.Value != "My Col" || !id.Quoted || id.QuoteRune != '"' {
		t.Errorf("got {Value:%q Quoted:%v Quote:%q}, want {My Col true \"}", id.Value, id.Quoted, string(id.QuoteRune))
	}
}

// TestParseIdentifier_BackquotedRejected verifies a `backtick`-quoted token is
// REJECTED as an identifier. The legacy grammar lists BACKQUOTED_IDENTIFIER_,
// and the lexer still tokenizes it, but Trino 481 rejects backtick quoting
// everywhere (oracle-confirmed: `SELECT * FROM `+"`bt table`"+` SYNTAX_ERROR).
// Only "double-quotes" are valid identifier quoting. See divergence ledger.
func TestParseIdentifier_BackquotedRejected(t *testing.T) {
	p := newParserForTest("`bt col`")
	if _, err := p.parseIdentifier(); err == nil {
		t.Fatal("parseIdentifier(`bt col`): want error (backtick rejected by Trino 481), got nil")
	}
}

// TestParseIdentifier_NonReservedKeyword verifies a non-reserved keyword (e.g.
// ZONE) is accepted as an identifier in the first position. The legacy
// non_reserved.sql example "SELECT zone FROM t" exercises exactly this.
func TestParseIdentifier_NonReservedKeyword(t *testing.T) {
	p := newParserForTest("zone")
	id, err := p.parseIdentifier()
	if err != nil {
		t.Fatalf("parseIdentifier(zone): %v", err)
	}
	if id.Value != "zone" || id.Quoted {
		t.Errorf("got {Value:%q Quoted:%v}, want {zone false}", id.Value, id.Quoted)
	}
}

// TestParseIdentifier_ReservedKeywordRejected verifies a reserved keyword (e.g.
// FROM) is NOT accepted as a bare identifier in the first position.
func TestParseIdentifier_ReservedKeywordRejected(t *testing.T) {
	p := newParserForTest("from")
	if _, err := p.parseIdentifier(); err == nil {
		t.Fatal("parseIdentifier(from): want error (FROM is reserved), got nil")
	}
}

// TestParseIdentifier_DigitLeadingRejected verifies a digit-leading identifier
// (DIGIT_IDENTIFIER_, e.g. 1abc) is REJECTED in identifier position. Although
// the legacy TrinoParser.g4 `identifier` rule lists DIGIT_IDENTIFIER_ as an
// alternative, Trino 481 rejects unquoted digit-leading identifiers everywhere
// (oracle-confirmed: `SELECT * FROM 1abc`, `CREATE TABLE 1abc (x int)`,
// `SELECT t.1abc FROM t` all SYNTAX_ERROR). See divergence ledger.
func TestParseIdentifier_DigitLeadingRejected(t *testing.T) {
	p := newParserForTest("1abc")
	if _, err := p.parseIdentifier(); err == nil {
		t.Fatal("parseIdentifier(1abc): want error (digit-leading rejected by Trino 481), got nil")
	}
}

// TestParseQualifiedName_OnePart verifies a single-part name.
func TestParseQualifiedName_OnePart(t *testing.T) {
	p := newParserForTest("orders")
	qn, err := p.parseQualifiedName()
	if err != nil {
		t.Fatalf("parseQualifiedName: %v", err)
	}
	if got := qn.Normalize(); got != "orders" {
		t.Errorf("Normalize = %q, want %q", got, "orders")
	}
	if len(qn.Parts) != 1 {
		t.Errorf("len(Parts) = %d, want 1", len(qn.Parts))
	}
}

// TestParseQualifiedName_ThreeParts verifies catalog.schema.table with mixed
// quoting normalizes correctly (unquoted lowered, quoted preserved).
func TestParseQualifiedName_ThreeParts(t *testing.T) {
	p := newParserForTest(`Cat."My Schema".TBL`)
	qn, err := p.parseQualifiedName()
	if err != nil {
		t.Fatalf("parseQualifiedName: %v", err)
	}
	if len(qn.Parts) != 3 {
		t.Fatalf("len(Parts) = %d, want 3", len(qn.Parts))
	}
	if got := qn.Normalize(); got != "cat.My Schema.tbl" {
		t.Errorf("Normalize = %q, want %q", got, "cat.My Schema.tbl")
	}
}

// TestParseQualifiedName_ReservedAfterDotRejected verifies a reserved keyword
// is REJECTED as a qualifiedName part even after a dot. Trino's grammar is
// `qualifiedName : identifier (DOT_ identifier)*` and `identifier` has no
// reserved-word alternative, so reserved words are rejected in every position
// (oracle-confirmed: `SELECT s.from FROM t`, `SELECT id FROM public.order`
// both SYNTAX_ERROR). This is unlike doris/snowflake, which accept reserved
// words post-dot. See divergence ledger.
func TestParseQualifiedName_ReservedAfterDotRejected(t *testing.T) {
	p := newParserForTest("s.from")
	if _, err := p.parseQualifiedName(); err == nil {
		t.Fatal("parseQualifiedName(s.from): want error (FROM reserved, rejected as name part), got nil")
	}
}

// TestParseQualifiedName_NonReservedAfterDot verifies a non-reserved keyword
// (ZONE) IS accepted as a dotted name part (oracle: `SELECT t.zone FROM t`
// ACCEPT).
func TestParseQualifiedName_NonReservedAfterDot(t *testing.T) {
	p := newParserForTest("t.zone")
	qn, err := p.parseQualifiedName()
	if err != nil {
		t.Fatalf("parseQualifiedName(t.zone): %v", err)
	}
	if got := qn.Normalize(); got != "t.zone" {
		t.Errorf("Normalize = %q, want %q", got, "t.zone")
	}
}

// TestNormalizeTrinoIdentifier verifies the exported normalizer is a faithful
// port of the legacy bytebase helper: unquoted (and backtick) spellings fold to
// lower case, only surrounding double-quotes are stripped (case preserved), and
// a doubled-quote escape is NOT collapsed (legacy slices the inter-quote bytes
// verbatim).
func TestNormalizeTrinoIdentifier(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"MyCol", "mycol"},
		{`"MyCol"`, "MyCol"},
		{"UPPER", "upper"},
		{`"keep THIS"`, "keep THIS"},
		// Legacy parity: backticks are NOT stripped (not valid Trino quoting);
		// the whole spelling is lower-cased.
		{"`MyCol`", "`mycol`"},
		// Legacy parity: the "" escape is left as-is, not collapsed to ".
		{`"a""b"`, `a""b`},
	}
	for _, c := range cases {
		if got := NormalizeTrinoIdentifier(c.in); got != c.want {
			t.Errorf("NormalizeTrinoIdentifier(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestExtractQualifiedNameParts verifies the exported helper splits a dotted
// name into normalized parts, handling quoting per Trino rules. Mirrors the
// legacy bytebase ExtractQualifiedNameParts helper.
func TestExtractQualifiedNameParts(t *testing.T) {
	parts, err := ExtractQualifiedNameParts(`Cat."My Schema".TBL`)
	if err != nil {
		t.Fatalf("ExtractQualifiedNameParts: %v", err)
	}
	want := []string{"cat", "My Schema", "tbl"}
	if len(parts) != len(want) {
		t.Fatalf("got %v, want %v", parts, want)
	}
	for i := range want {
		if parts[i] != want[i] {
			t.Errorf("part %d = %q, want %q", i, parts[i], want[i])
		}
	}
}

// TestParseQualifiedNameString verifies the standalone ParseQualifiedName entry
// returns a node whose String() is source-faithful (re-quoting preserved).
func TestParseQualifiedNameString(t *testing.T) {
	qn, errs := ParseQualifiedName(`a."B c".d`)
	if len(errs) != 0 {
		t.Fatalf("ParseQualifiedName: %v", errs)
	}
	if got := qn.String(); got != `a."B c".d` {
		t.Errorf("String = %q, want %q", got, `a."B c".d`)
	}
}

// TestParseQualifiedName_TrailingTokens verifies ParseQualifiedName reports an
// error when input has tokens after the name.
func TestParseQualifiedName_TrailingTokens(t *testing.T) {
	_, errs := ParseQualifiedName("a b")
	if len(errs) == 0 {
		t.Fatal("ParseQualifiedName(\"a b\"): want trailing-token error, got none")
	}
}

// TestParseQualifiedName_LocSpansAllParts checks the qualified-name Loc covers
// from the first part start to the last part end.
func TestParseQualifiedName_LocSpansAllParts(t *testing.T) {
	in := "cat.sch.tbl"
	p := newParserForTest(in)
	qn, err := p.parseQualifiedName()
	if err != nil {
		t.Fatalf("parseQualifiedName: %v", err)
	}
	if qn.Loc.Start != 0 || qn.Loc.End != len(in) {
		t.Errorf("Loc = %+v, want [0,%d)", qn.Loc, len(in))
	}
	if !qn.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
	_ = ast.NoLoc()
}
