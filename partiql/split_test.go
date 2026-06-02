package partiql

import (
	"testing"

	"github.com/bytebase/omni/partiql/parser"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantN     int
		wantText  []string
		wantEmpty []bool
	}{
		{
			name:      "single_statement",
			input:     "SELECT * FROM t",
			wantN:     1,
			wantText:  []string{"SELECT * FROM t"},
			wantEmpty: []bool{false},
		},
		{
			name:      "two_statements",
			input:     "SELECT 1; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT 1;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "trailing_semicolon",
			input:     "SELECT 1;",
			wantN:     1,
			wantText:  []string{"SELECT 1;"},
			wantEmpty: []bool{false},
		},
		{
			name:      "multiple_trailing_semicolons",
			input:     "SELECT 1;;",
			wantN:     2,
			wantText:  []string{"SELECT 1;", ";"},
			wantEmpty: []bool{false, true},
		},
		{
			name:      "semicolon_in_string",
			input:     "SELECT 'a;b'",
			wantN:     1,
			wantText:  []string{"SELECT 'a;b'"},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_quoted_ident",
			input:     `SELECT "a;b"`,
			wantN:     1,
			wantText:  []string{`SELECT "a;b"`},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_ion_literal",
			input:     "SELECT `{a;b}`",
			wantN:     1,
			wantText:  []string{"SELECT `{a;b}`"},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_line_comment",
			input:     "SELECT 1 -- comment; not a split\n; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT 1 -- comment; not a split\n;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "semicolon_in_block_comment",
			input:     "SELECT /* ; */ 1; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT /* ; */ 1;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:  "empty_input",
			input: "",
			wantN: 0,
		},
		{
			name:      "whitespace_only",
			input:     "   \n\t  ",
			wantN:     1,
			wantText:  []string{"   \n\t  "},
			wantEmpty: []bool{true},
		},
		{
			name:      "comment_only",
			input:     "-- just a comment",
			wantN:     1,
			wantText:  []string{"-- just a comment"},
			wantEmpty: []bool{true},
		},
		{
			name:      "doubled_quote_escape",
			input:     "SELECT 'it''s'; SELECT 1",
			wantN:     2,
			wantText:  []string{"SELECT 'it''s';", " SELECT 1"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "nested_block_comment",
			input:     "SELECT /* outer /* inner */ still comment */ 1",
			wantN:     1,
			wantText:  []string{"SELECT /* outer /* inner */ still comment */ 1"},
			wantEmpty: []bool{false},
		},
		{
			name:      "ion_with_backtick_in_short_string",
			input:     "SELECT `\"a`b\"`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `\"a`b\"`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_backtick_in_comment",
			input:     "SELECT `/* not end ` still ion */`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `/* not end ` still ion */`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_backtick_in_symbol",
			input:     "SELECT `'a`b'`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `'a`b'`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_triple_quoted_long_string",
			input:     "SELECT `'''contains ` backtick'''`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `'''contains ` backtick'''`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_simple_no_inner_strings",
			input:     "SELECT `{a: 1}`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `{a: 1}`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		// --- Finding A (round 5): a ';' inside an Ion // LINE comment (inside a
		//     backtick Ion literal) must NOT split. The Ion // line comment runs
		//     to the newline; everything on that line — including the ';' — is
		//     comment content, protected by the Ion literal. The literal closes
		//     at the standalone backtick AFTER the newline, and only the ';' that
		//     follows the literal splits. The old splitter (no // handling) split
		//     at the in-comment ';', disagreeing with the lexer. ---
		{
			name:      "ion_semicolon_in_line_comment",
			input:     "SELECT `// a ; b\n{x:1}`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `// a ; b\n{x:1}`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		// --- Finding A: a backtick inside an Ion // line comment must NOT close
		//     the Ion literal (the // comment protects it through end-of-line),
		//     so the ';' inside the comment must NOT split either. ---
		{
			name:      "ion_backtick_in_line_comment",
			input:     "SELECT `// not ` end ; still\n{x:1}`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `// not ` end ; still\n{x:1}`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		// --- Finding A: a // line comment is the LAST thing in the literal
		//     before the closing backtick on the next line. ---
		{
			name:      "ion_line_comment_then_close",
			input:     "SELECT `1 // c\n`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `1 // c\n`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		// --- Finding A: a base64 blob {{...}} carrying a ';' is one Ion token;
		//     the ';' inside it must NOT split (the old splitter had no blob
		//     handling). Here the blob body is valid base64; the ';' sits in the
		//     SQL after the literal and is the real split point. ---
		{
			name:      "ion_blob_no_split_inside",
			input:     "SELECT `{{ aGVsbG8= }}`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `{{ aGVsbG8= }}`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			segs := Split(tc.input)
			if len(segs) != tc.wantN {
				t.Fatalf("Split(%q) returned %d segments, want %d", tc.input, len(segs), tc.wantN)
			}
			for i, seg := range segs {
				if i < len(tc.wantText) && seg.Text != tc.wantText[i] {
					t.Errorf("seg[%d].Text = %q, want %q", i, seg.Text, tc.wantText[i])
				}
				if i < len(tc.wantEmpty) && seg.Empty() != tc.wantEmpty[i] {
					t.Errorf("seg[%d].Empty() = %v, want %v", i, seg.Empty(), tc.wantEmpty[i])
				}
			}
		})
	}
}

// lexerIonEnd drives the lexer over an input whose FIRST token is a backtick
// Ion literal and reports the lexer's authoritative boundary: when the literal
// closes, closed=true and end is the literal token's Loc.End; when it is
// unterminated (the lexer sets Err and emits no literal token), closed=false.
func lexerIonEnd(t *testing.T, input string) (end int, closed bool) {
	t.Helper()
	lx := parser.NewLexer(input)
	tok := lx.Next()
	if lx.Err != nil {
		// Unterminated Ion literal: the lexer reports an error rather than a
		// literal token. The splitter mirrors this with closed=false.
		return len(input), false
	}
	if tok.Loc.Start != 0 {
		t.Fatalf("lexer first token did not start at the backtick: input=%q start=%d", input, tok.Loc.Start)
	}
	return tok.Loc.End, true
}

// TestSplit_IonBoundaryMatchesLexer is the Finding A (round 5) regression: the
// statement splitter MUST agree with the lexer, byte-for-byte, on where a
// backtick Ion literal ends. They previously diverged because the splitter's
// private Ion scanner did not handle Ion // line comments (nor blobs, nor the
// maximal-munch ION_ANY fallback), so a ';' or backtick inside such a comment
// split at the wrong byte. The fix routes both through parser.ScanIonLiteral;
// this test pins that contract from the consumer side and would fail if the
// splitter ever re-grew a divergent copy.
//
// For each input (which BEGINS with a backtick), we compare:
//   - splitter boundary  = parser.ScanIonLiteral(input, 1)  (the exact call
//     Split makes for a leading backtick), and
//   - lexer boundary      = the first token's Loc.End.
func TestSplit_IonBoundaryMatchesLexer(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		// The historically-divergent construct: Ion // line comments.
		{"line_comment_semicolon", "`// a ; b\n{x:1}` rest"},
		{"line_comment_backtick", "`// not ` end ; still\n{x:1}` rest"},
		{"line_comment_then_close", "`1 // c\n` rest"},
		{"line_comment_eats_closer_unterminated", "`{x:1} // no newline so the backtick is eaten "},
		{"line_comment_cr_terminated", "`// c\r{x:1}` rest"},
		// Block comments (already handled, re-confirmed).
		{"block_comment_backtick", "`/* not ` end */{x:1}` rest"},
		{"block_comment_unterminated_fallback", "`/* a ` b` rest"},
		// Short strings / symbols / long strings.
		{"short_string_backtick", "`\"a`b\"` rest"},
		{"symbol_backtick", "`'a`b'` rest"},
		{"long_string_backtick", "`'''a`b'''` rest"},
		{"unclosed_short_string_fallback", "`\"a`b rest"},
		{"unclosed_long_string_fallback", "`'''a`b rest"},
		// Escapes (valid vs invalid -> ION_ANY fallback).
		{"valid_hex_escape", "`\"a\\x41b`\"` rest"},
		{"invalid_escape_fallback", "`\"a\\qb`c rest"},
		// Base64 blobs / clobs.
		{"blob_valid_base64", "`{{ aGVsbG8= }}` rest"},
		{"blob_slash_slash_content", "`{{ //// }}` rest"},
		{"clob_quoted_string_brace", "`{{ \"}}\" }}` rest"},
		{"invalid_blob_fallback", "`{{ not*base64 }}` rest"},
		// Round-5 disputed VT/FF long-string shape: the literal spans the inner
		// backtick (VT/FF are content, per the oracle) — the splitter must agree.
		{"vtab_long_string_spans", "`'''a\vb`'''` rest"},
		{"formfeed_long_string_spans", "`'''a\fb`'''` rest"},
		{"nul_long_string_closes_first_backtick", "`'''a\x00b`+`'''x'''` rest"},
		// Plain / nested structures and timestamps.
		{"plain_struct", "`{a: 1; b: 2}` rest"},
		{"timestamp", "`2017-09-14T` rest"},
		{"empty_literal", "`` rest"},
		{"unterminated_bare", "`abc"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			splitEnd, splitClosed := parser.ScanIonLiteral(tc.input, 1)
			lexEnd, lexClosed := lexerIonEnd(t, tc.input)
			if splitClosed != lexClosed {
				t.Fatalf("Ion termination disagreement for %q: splitter closed=%v, lexer closed=%v", tc.input, splitClosed, lexClosed)
			}
			if splitEnd != lexEnd {
				t.Fatalf("Ion boundary mismatch for %q:\n  splitter end = %d (literal=%q)\n  lexer    end = %d (literal=%q)",
					tc.input, splitEnd, tc.input[:min(splitEnd, len(tc.input))], lexEnd, tc.input[:min(lexEnd, len(tc.input))])
			}
		})
	}
}
