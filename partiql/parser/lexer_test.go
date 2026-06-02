package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// tokenStreamCase is one row in the TestLexer_Tokens table.
type tokenStreamCase struct {
	name   string
	input  string
	tokens []Token // expected, excluding the trailing EOF
}

// runTokenStreamCase drains l.Next() until tokEOF and asserts the captured
// stream matches tc.tokens token-by-token. Asserts l.Err == nil.
func runTokenStreamCase(t *testing.T, tc tokenStreamCase) {
	t.Helper()
	l := NewLexer(tc.input)
	var got []Token
	for {
		tok := l.Next()
		if tok.Type == tokEOF {
			break
		}
		got = append(got, tok)
	}
	if l.Err != nil {
		t.Fatalf("unexpected error: %v", l.Err)
	}
	if len(got) != len(tc.tokens) {
		t.Errorf("token count mismatch: got %d tokens, want %d\n got: %+v\nwant: %+v",
			len(got), len(tc.tokens), got, tc.tokens)
		return
	}
	for i := range got {
		if got[i] != tc.tokens[i] {
			t.Errorf("token[%d] mismatch:\n got: %+v\nwant: %+v", i, got[i], tc.tokens[i])
		}
	}
}

// TestLexer_Tokens is the master golden-test table for the lexer.
// Sections cover whitespace/comments, string literals, quoted identifiers,
// unquoted identifiers, keywords, numeric literals, operators, punctuation,
// multi-token integration sequences, and Ion literals.
//
// All cases assert token-by-token equality and l.Err == nil.
func TestLexer_Tokens(t *testing.T) {
	cases := []tokenStreamCase{
		// =============================================================
		// Empty input + whitespace + comments
		// =============================================================
		{"empty", "", nil},
		{"whitespace_spaces", "   ", nil},
		{"whitespace_tabs", "\t\t", nil},
		{"whitespace_newlines", "\n\n", nil},
		{"whitespace_mixed", " \t\n\r ", nil},
		{"line_comment_only", "-- a comment\n", nil},
		{"line_comment_at_eof", "-- a comment without newline", nil},
		{"block_comment_only", "/* a comment */", nil},
		{"block_comment_multiline", "/*\nmulti\nline\n*/", nil},

		// =============================================================
		// String literals (Task 6)
		// =============================================================
		{
			"string_simple",
			"'hello'",
			[]Token{{tokSCONST, "hello", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"string_empty",
			"''",
			[]Token{{tokSCONST, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"string_doubled_quote",
			"'it''s'",
			[]Token{{tokSCONST, "it's", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"string_with_whitespace",
			"'  '",
			[]Token{{tokSCONST, "  ", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"string_with_special_chars",
			"'a!@#%^&*()'",
			[]Token{{tokSCONST, "a!@#%^&*()", ast.Loc{Start: 0, End: 12}}},
		},
		{
			"string_only_doubled_quotes",
			"''''''",
			[]Token{{tokSCONST, "''", ast.Loc{Start: 0, End: 6}}},
		},

		// =============================================================
		// Quoted identifiers (Task 6)
		// =============================================================
		{
			"quoted_ident_simple",
			`"Foo"`,
			[]Token{{tokIDENT_QUOTED, "Foo", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"quoted_ident_empty",
			`""`,
			[]Token{{tokIDENT_QUOTED, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"quoted_ident_doubled_quote",
			`"a""b"`,
			[]Token{{tokIDENT_QUOTED, `a"b`, ast.Loc{Start: 0, End: 6}}},
		},
		{
			"quoted_ident_with_space",
			`"Foo Bar"`,
			[]Token{{tokIDENT_QUOTED, "Foo Bar", ast.Loc{Start: 0, End: 9}}},
		},
		{
			"quoted_ident_case_preserved",
			`"FoO"`,
			[]Token{{tokIDENT_QUOTED, "FoO", ast.Loc{Start: 0, End: 5}}},
		},

		// =============================================================
		// Unquoted identifiers (Task 7)
		// =============================================================
		{
			"unquoted_ident_simple",
			"foo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"unquoted_ident_with_underscore",
			"_foo",
			[]Token{{tokIDENT, "_foo", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"unquoted_ident_with_dollar",
			"x$y",
			[]Token{{tokIDENT, "x$y", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"unquoted_ident_with_digit_in_middle",
			"a1b2",
			[]Token{{tokIDENT, "a1b2", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"unquoted_ident_uppercase_preserved",
			"FOO",
			[]Token{{tokIDENT, "FOO", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"unquoted_ident_leading_dollar",
			"$rowid",
			[]Token{{tokIDENT, "$rowid", ast.Loc{Start: 0, End: 6}}},
		},

		// =============================================================
		// Keywords (case-insensitive lookup, raw text preserved) (Task 7)
		// =============================================================
		{
			"keyword_select_lower",
			"select",
			[]Token{{tokSELECT, "select", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_select_upper",
			"SELECT",
			[]Token{{tokSELECT, "SELECT", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_select_mixed",
			"Select",
			[]Token{{tokSELECT, "Select", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_from",
			"FROM",
			[]Token{{tokFROM, "FROM", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"keyword_where",
			"WHERE",
			[]Token{{tokWHERE, "WHERE", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"keyword_pivot_partiql_unique",
			"PIVOT",
			[]Token{{tokPIVOT, "PIVOT", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"keyword_missing_partiql_unique",
			"MISSING",
			[]Token{{tokMISSING, "MISSING", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"keyword_bag_data_type",
			"BAG",
			[]Token{{tokBAG, "BAG", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"keyword_can_lossless_cast_underscored",
			"CAN_LOSSLESS_CAST",
			[]Token{{tokCAN_LOSSLESS_CAST, "CAN_LOSSLESS_CAST", ast.Loc{Start: 0, End: 17}}},
		},
		{
			"keyword_can_lossless_cast_mixed_case",
			"Can_Lossless_Cast",
			[]Token{{tokCAN_LOSSLESS_CAST, "Can_Lossless_Cast", ast.Loc{Start: 0, End: 17}}},
		},

		// =============================================================
		// Identifier vs keyword cases after whitespace/comments (Task 7)
		// =============================================================
		{
			"ident_after_line_comment",
			"-- skipped\nfoo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 11, End: 14}}},
		},
		{
			"ident_after_block_comment",
			"/* x */ foo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 8, End: 11}}},
		},

		// =============================================================
		// Numeric literals (Task 8)
		// =============================================================
		{
			"integer",
			"42",
			[]Token{{tokICONST, "42", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"integer_zero",
			"0",
			[]Token{{tokICONST, "0", ast.Loc{Start: 0, End: 1}}},
		},
		{
			"integer_large",
			"1234567890",
			[]Token{{tokICONST, "1234567890", ast.Loc{Start: 0, End: 10}}},
		},
		{
			"decimal_dot",
			"3.14",
			[]Token{{tokFCONST, "3.14", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_leading_dot",
			".5",
			[]Token{{tokFCONST, ".5", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"decimal_trailing_dot",
			"42.",
			[]Token{{tokFCONST, "42.", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"decimal_scientific_lower_e",
			"1e10",
			[]Token{{tokFCONST, "1e10", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_scientific_upper_e",
			"1E10",
			[]Token{{tokFCONST, "1E10", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_scientific_negative_exp",
			"1.5e-3",
			[]Token{{tokFCONST, "1.5e-3", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"decimal_scientific_positive_exp",
			"2.5e+4",
			[]Token{{tokFCONST, "2.5e+4", ast.Loc{Start: 0, End: 6}}},
		},

		// =============================================================
		// Single-character operators (Task 9)
		// =============================================================
		{"op_plus", "+", []Token{{tokPLUS, "+", ast.Loc{Start: 0, End: 1}}}},
		{"op_minus", "-", []Token{{tokMINUS, "-", ast.Loc{Start: 0, End: 1}}}},
		{"op_asterisk", "*", []Token{{tokASTERISK, "*", ast.Loc{Start: 0, End: 1}}}},
		{"op_slash_forward", "/", []Token{{tokSLASH_FORWARD, "/", ast.Loc{Start: 0, End: 1}}}},
		{"op_percent", "%", []Token{{tokPERCENT, "%", ast.Loc{Start: 0, End: 1}}}},
		{"op_caret", "^", []Token{{tokCARET, "^", ast.Loc{Start: 0, End: 1}}}},
		{"op_tilde", "~", []Token{{tokTILDE, "~", ast.Loc{Start: 0, End: 1}}}},
		{"op_at_sign", "@", []Token{{tokAT_SIGN, "@", ast.Loc{Start: 0, End: 1}}}},
		{"op_eq", "=", []Token{{tokEQ, "=", ast.Loc{Start: 0, End: 1}}}},
		{"op_lt", "<", []Token{{tokLT, "<", ast.Loc{Start: 0, End: 1}}}},
		{"op_gt", ">", []Token{{tokGT, ">", ast.Loc{Start: 0, End: 1}}}},

		// =============================================================
		// Punctuation (Task 9)
		// =============================================================
		{"punct_paren_left", "(", []Token{{tokPAREN_LEFT, "(", ast.Loc{Start: 0, End: 1}}}},
		{"punct_paren_right", ")", []Token{{tokPAREN_RIGHT, ")", ast.Loc{Start: 0, End: 1}}}},
		{"punct_bracket_left", "[", []Token{{tokBRACKET_LEFT, "[", ast.Loc{Start: 0, End: 1}}}},
		{"punct_bracket_right", "]", []Token{{tokBRACKET_RIGHT, "]", ast.Loc{Start: 0, End: 1}}}},
		{"punct_brace_left", "{", []Token{{tokBRACE_LEFT, "{", ast.Loc{Start: 0, End: 1}}}},
		{"punct_brace_right", "}", []Token{{tokBRACE_RIGHT, "}", ast.Loc{Start: 0, End: 1}}}},
		{"punct_colon", ":", []Token{{tokCOLON, ":", ast.Loc{Start: 0, End: 1}}}},
		{"punct_colon_semi", ";", []Token{{tokCOLON_SEMI, ";", ast.Loc{Start: 0, End: 1}}}},
		{"punct_comma", ",", []Token{{tokCOMMA, ",", ast.Loc{Start: 0, End: 1}}}},
		{"punct_period", ".", []Token{{tokPERIOD, ".", ast.Loc{Start: 0, End: 1}}}},
		{"punct_question_mark", "?", []Token{{tokQUESTION_MARK, "?", ast.Loc{Start: 0, End: 1}}}},

		// =============================================================
		// Two-character operators (Task 9)
		// =============================================================
		{"op_lt_eq", "<=", []Token{{tokLT_EQ, "<=", ast.Loc{Start: 0, End: 2}}}},
		{"op_gt_eq", ">=", []Token{{tokGT_EQ, ">=", ast.Loc{Start: 0, End: 2}}}},
		{"op_neq_angle", "<>", []Token{{tokNEQ, "<>", ast.Loc{Start: 0, End: 2}}}},
		{"op_neq_bang", "!=", []Token{{tokNEQ, "!=", ast.Loc{Start: 0, End: 2}}}},
		{"op_concat", "||", []Token{{tokCONCAT, "||", ast.Loc{Start: 0, End: 2}}}},
		{"op_angle_double_left", "<<", []Token{{tokANGLE_DOUBLE_LEFT, "<<", ast.Loc{Start: 0, End: 2}}}},
		{"op_angle_double_right", ">>", []Token{{tokANGLE_DOUBLE_RIGHT, ">>", ast.Loc{Start: 0, End: 2}}}},

		// =============================================================
		// Multi-token integration sequences (Task 9)
		// =============================================================
		{
			"select_star_from_table",
			"SELECT * FROM Music",
			[]Token{
				{tokSELECT, "SELECT", ast.Loc{Start: 0, End: 6}},
				{tokASTERISK, "*", ast.Loc{Start: 7, End: 8}},
				{tokFROM, "FROM", ast.Loc{Start: 9, End: 13}},
				{tokIDENT, "Music", ast.Loc{Start: 14, End: 19}},
			},
		},
		{
			"where_with_string_literal",
			"WHERE Artist='Pink Floyd'",
			[]Token{
				{tokWHERE, "WHERE", ast.Loc{Start: 0, End: 5}},
				{tokIDENT, "Artist", ast.Loc{Start: 6, End: 12}},
				{tokEQ, "=", ast.Loc{Start: 12, End: 13}},
				{tokSCONST, "Pink Floyd", ast.Loc{Start: 13, End: 25}},
			},
		},
		{
			"bag_literal",
			"<<1, 2, 3>>",
			[]Token{
				{tokANGLE_DOUBLE_LEFT, "<<", ast.Loc{Start: 0, End: 2}},
				{tokICONST, "1", ast.Loc{Start: 2, End: 3}},
				{tokCOMMA, ",", ast.Loc{Start: 3, End: 4}},
				{tokICONST, "2", ast.Loc{Start: 5, End: 6}},
				{tokCOMMA, ",", ast.Loc{Start: 6, End: 7}},
				{tokICONST, "3", ast.Loc{Start: 8, End: 9}},
				{tokANGLE_DOUBLE_RIGHT, ">>", ast.Loc{Start: 9, End: 11}},
			},
		},
		{
			"two_char_op_in_expr",
			"a<=5",
			[]Token{
				{tokIDENT, "a", ast.Loc{Start: 0, End: 1}},
				{tokLT_EQ, "<=", ast.Loc{Start: 1, End: 3}},
				{tokICONST, "5", ast.Loc{Start: 3, End: 4}},
			},
		},
		{
			"path_expression_dot",
			"t.foo",
			[]Token{
				{tokIDENT, "t", ast.Loc{Start: 0, End: 1}},
				{tokPERIOD, ".", ast.Loc{Start: 1, End: 2}},
				{tokIDENT, "foo", ast.Loc{Start: 2, End: 5}},
			},
		},

		// =============================================================
		// Ion literals — base forms (no inner backtick)
		// =============================================================
		{
			"ion_simple",
			"`{a: 1}`",
			[]Token{{tokION_LITERAL, "{a: 1}", ast.Loc{Start: 0, End: 8}}},
		},
		{
			"ion_empty",
			"``",
			[]Token{{tokION_LITERAL, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"ion_with_whitespace",
			"`  abc  `",
			[]Token{{tokION_LITERAL, "  abc  ", ast.Loc{Start: 0, End: 9}}},
		},
		// A single-quoted Ion symbol with no inner backtick. The closing
		// backtick after the symbol terminates the literal.
		{
			"ion_quoted_symbol_plain",
			"`'hello'`",
			[]Token{{tokION_LITERAL, "'hello'", ast.Loc{Start: 0, End: 9}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runTokenStreamCase(t, tc)
		})
	}
}

// TestLexer_IonMode pins the Ion-mode-aware backtick scanner
// (node parser-ion-literals). On a backtick the lexer enters Ion
// sub-mode and consumes an Ion value verbatim up to a *standalone*
// closing backtick, emitting one tokION_LITERAL. Crucially, a
// backtick that appears INSIDE an Ion sub-token — single-quoted
// symbol, double-quoted short string, triple-quoted long string,
// {{ }} lob, // line comment, or /* */ block comment — does NOT
// terminate the literal.
//
// Grammar: PartiQLLexer.g4 ION mode (lines 406-430). Token.Str is
// the verbatim inner content (no decoding); Token.Loc covers the
// full `...` range including both backticks.
func TestLexer_IonMode(t *testing.T) {
	cases := []tokenStreamCase{
		// --- backtick inside a single-quoted symbol ---
		{
			"ion_backtick_in_symbol",
			"`'a`b'`",
			[]Token{{tokION_LITERAL, "'a`b'", ast.Loc{Start: 0, End: 7}}},
		},
		// --- backtick inside a double-quoted short string ---
		{
			"ion_backtick_in_short_string",
			"`\"x`y\"`",
			[]Token{{tokION_LITERAL, "\"x`y\"", ast.Loc{Start: 0, End: 7}}},
		},
		// --- backtick inside a triple-quoted long string ---
		{
			"ion_backtick_in_long_string",
			"`'''a`b'''`",
			[]Token{{tokION_LITERAL, "'''a`b'''", ast.Loc{Start: 0, End: 11}}},
		},
		// --- long string spanning a newline (no inner backtick) ---
		{
			"ion_long_string_multiline",
			"`'''a\nb'''`",
			[]Token{{tokION_LITERAL, "'''a\nb'''", ast.Loc{Start: 0, End: 11}}},
		},
		// --- four single-quotes: ANTLR maximal-munch can't match a
		//     triple-quoted long string (no closing '''), so the quotes
		//     are two empty symbols ('' then ''), then the backtick
		//     closes. This is the lexer-grammar fallback, not a long
		//     string. ---
		{
			"ion_four_quotes_are_two_empty_symbols",
			"`''''`",
			[]Token{{tokION_LITERAL, "''''", ast.Loc{Start: 0, End: 6}}},
		},
		// --- six single-quotes IS an empty long string (open + close) ---
		{
			"ion_six_quotes_empty_long_string",
			"`''''''`",
			[]Token{{tokION_LITERAL, "''''''", ast.Loc{Start: 0, End: 8}}},
		},
		// --- a long string whose body contains a backtick, closes with
		//     the triple-quote, then trailing content before the literal
		//     closes. The inner backtick is captured, not a terminator. ---
		{
			"ion_long_string_with_backtick_then_more",
			"`'''a`b''' x`",
			[]Token{{tokION_LITERAL, "'''a`b''' x", ast.Loc{Start: 0, End: 13}}},
		},
		// --- backtick inside a // line comment ---
		{
			"ion_backtick_in_line_comment",
			"`// c`omment\n42`",
			[]Token{{tokION_LITERAL, "// c`omment\n42", ast.Loc{Start: 0, End: 16}}},
		},
		// --- backtick inside a /* */ block comment ---
		{
			"ion_backtick_in_block_comment",
			"`/* c`c */1`",
			[]Token{{tokION_LITERAL, "/* c`c */1", ast.Loc{Start: 0, End: 12}}},
		},
		// --- a lob blob: {{ base64 }}. A valid ION_BLOB (g4 414-415:
		//     LOB_START (BASE_64_QUARTET | WS)* BASE_64_PAD? WS* LOB_END) is
		//     matched whole, so the closing backtick after the }} closes the
		//     literal. ---
		{
			"ion_lob_blob",
			"`{{ aGVsbG8= }}`",
			[]Token{{tokION_LITERAL, "{{ aGVsbG8= }}", ast.Loc{Start: 0, End: 16}}},
		},
		// --- a BLOB whose base64 content contains '//'. Because the whole
		//     {{ //// }} matches ION_BLOB as a single maximal-munch token
		//     (base64 INCLUDES '/', and '////' is a valid BASE_64_QUARTET),
		//     the '//' is base64 CONTENT, never an ION_INLINE_COMMENT. The
		//     backtick after the }} closes the literal. Codex P2 regression
		//     case; verified against the generated ANTLR lexer oracle:
		//       `{{ //// }}`  ->  one ION_CLOSURE spanning the whole input.
		{
			"ion_blob_double_slash_is_content",
			"`{{ //// }}`",
			[]Token{{tokION_LITERAL, "{{ //// }}", ast.Loc{Start: 0, End: 12}}},
		},
		// --- a BLOB with '=' padding (BASE_64_PAD2: '//==' is CHAR CHAR '=' '=')
		//     and interior whitespace between quartets. Still one ION_BLOB. ---
		{
			"ion_blob_slash_with_padding",
			"`{{ //// //// //// //== }}`",
			[]Token{{tokION_LITERAL, "{{ //// //// //// //== }}", ast.Loc{Start: 0, End: 27}}},
		},
		// --- a BLOB with no surrounding whitespace: {{////}} is a valid
		//     ION_BLOB ('////' = one quartet, no WS, immediate LOB_END). ---
		{
			"ion_blob_slash_tight",
			"`{{////}}`",
			[]Token{{tokION_LITERAL, "{{////}}", ast.Loc{Start: 0, End: 10}}},
		},
		// --- two adjacent BLOBs are each matched as ION_BLOB, then the
		//     backtick closes. ---
		{
			"ion_two_blobs_adjacent",
			"`{{////}}{{////}}`",
			[]Token{{tokION_LITERAL, "{{////}}{{////}}", ast.Loc{Start: 0, End: 18}}},
		},
		// --- a CLOB whose quoted string contains '//': the '{{ ... }}' is NOT
		//     a valid ION_BLOB (the '"' is not base64), so the braces are
		//     ION_ANY and the '"a//b"' is a SHORT_QUOTED_STRING in which '//'
		//     is plain content (not a comment). The trailing backtick closes.
		//     Confirms the blob and string scanners compose. Oracle: single
		//     literal. ---
		{
			"ion_clob_double_slash_in_string",
			"`{{ \"a//b\" }}`",
			[]Token{{tokION_LITERAL, "{{ \"a//b\" }}", ast.Loc{Start: 0, End: 14}}},
		},
		// --- a backtick in a {{ }} region is NOT protected: braces are
		//     ION_ANY (the grammar's ION_BLOB matches base64 only, which
		//     contains no backtick), so the backtick after "x" is a
		//     standalone ION_CLOSURE and the literal ends there. The rest
		//     (`y`) lexes as an ordinary identifier. Verified against the
		//     generated ANTLR lexer oracle:
		//       `{{ x`  -> ION_CLOSURE; y -> IDENTIFIER.
		{
			"ion_backtick_in_lob_region_closes",
			"`{{ x`y",
			[]Token{
				{tokION_LITERAL, "{{ x", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "y", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- CLOB with }} inside a double-quoted short string: the }} is
		//     protected by SHORT_QUOTED_STRING, so the lob does not end
		//     early; only the trailing backtick (after the real }}) closes.
		//     This is the Codex P2 regression case. Oracle: single literal. ---
		{
			"ion_clob_double_brace_in_short_string",
			"`{{ \"}}\" }}`",
			[]Token{{tokION_LITERAL, "{{ \"}}\" }}", ast.Loc{Start: 0, End: 12}}},
		},
		// --- CLOB with }} mid short string then trailing content ---
		{
			"ion_clob_double_brace_midstring",
			"`{{ \"a}}b\" }}`",
			[]Token{{tokION_LITERAL, "{{ \"a}}b\" }}", ast.Loc{Start: 0, End: 14}}},
		},
		// --- CLOB with }} inside a triple-quoted long string ---
		{
			"ion_clob_double_brace_in_long_string",
			"`{{ '''}}''' }}`",
			[]Token{{tokION_LITERAL, "{{ '''}}''' }}", ast.Loc{Start: 0, End: 16}}},
		},
		// --- }} inside a single-quoted symbol within the lob region ---
		{
			"ion_clob_double_brace_in_symbol",
			"`{{ '}}' }}`",
			[]Token{{tokION_LITERAL, "{{ '}}' }}", ast.Loc{Start: 0, End: 12}}},
		},
		// --- escaped quote then }} inside a short string in the lob region;
		//     the backslash escape keeps the string open across the inner
		//     quote, so the }} is still string content. ---
		{
			"ion_clob_escaped_quote_then_double_brace",
			"`{{ \"a\\\"}}\" }}`",
			[]Token{{tokION_LITERAL, "{{ \"a\\\"}}\" }}", ast.Loc{Start: 0, End: 15}}},
		},
		// --- bare {{ with no base64 and no closer is just two ION_ANY
		//     braces; the backtick closes. ---
		{
			"ion_bare_double_open_brace",
			"`{{`",
			[]Token{{tokION_LITERAL, "{{", ast.Loc{Start: 0, End: 4}}},
		},
		// --- bare }} is two ION_ANY braces, then the backtick closes. ---
		{
			"ion_bare_double_close_brace",
			"`}}`",
			[]Token{{tokION_LITERAL, "}}", ast.Loc{Start: 0, End: 4}}},
		},
		// --- an unclosed " inside the lob region degrades to ION_ANY (no
		//     closing "), so the }} after it is ordinary content and the
		//     backtick closes. Oracle: single literal `{{ "unclosed }}`. ---
		{
			"ion_lob_unclosed_string_degrades",
			"`{{ \"unclosed }}`",
			[]Token{{tokION_LITERAL, "{{ \"unclosed }}", ast.Loc{Start: 0, End: 17}}},
		},
		// --- escaped quote inside a symbol must not end the symbol ---
		{
			"ion_escaped_quote_in_symbol",
			"`'a\\'`b'`",
			[]Token{{tokION_LITERAL, "'a\\'`b'", ast.Loc{Start: 0, End: 9}}},
		},
		// --- escaped quote inside a short string ---
		{
			"ion_escaped_quote_in_short_string",
			"`\"a\\\"`b\"`",
			[]Token{{tokION_LITERAL, "\"a\\\"`b\"", ast.Loc{Start: 0, End: 9}}},
		},
		// --- INVALID escape does not protect the following byte. Per
		//     PartiQLLexer.g4 TEXT_ESCAPE (g4 465-466) only COMMON_ESCAPE /
		//     HEX_ESCAPE / UNICODE_ESCAPE are escapes; '\q' is none of these.
		//     So inside a short string, '"\q' cannot be SHORT_QUOTED_STRING
		//     content, the '"' degrades to ION_ANY, and the backtick after
		//     '\q' is a standalone ION_CLOSURE that ends the literal. The
		//     trailing 'x' is an ordinary identifier. Codex P2 regression
		//     case; oracle: `"\q` -> ION_CLOSURE (content "\q); x -> IDENT.
		{
			"ion_invalid_escape_short_string_closes",
			"`\"\\q`x",
			[]Token{
				{tokION_LITERAL, "\"\\q", ast.Loc{Start: 0, End: 5}},
				{tokIDENT, "x", ast.Loc{Start: 5, End: 6}},
			},
		},
		// --- same for a single-quoted symbol: '\q' is not a valid escape,
		//     the "'" degrades to ION_ANY, and the backtick closes. ---
		{
			"ion_invalid_escape_symbol_closes",
			"`'\\q`x",
			[]Token{
				{tokION_LITERAL, "'\\q", ast.Loc{Start: 0, End: 5}},
				{tokIDENT, "x", ast.Loc{Start: 5, End: 6}},
			},
		},
		// --- INVALID escape with a LATER quote. This is the exact Codex P2
		//     shape: a short string with '\q', then a backtick, then a later
		//     '"'. A buggy scanner that treats '\q' as a protected escape
		//     would keep the string open, close it at the LATER '"', and thus
		//     swallow the first backtick — emitting the wrong content (or
		//     erroring). The grammar forbids this: '\q' is not a TEXT_ESCAPE,
		//     so the '"' is ION_ANY and the first backtick closes the literal.
		//     The trailing '"z"' is then an ordinary quoted identifier. Oracle:
		//       `"\q` -> ION_CLOSURE (content "\q); "z" -> IDENTIFIER_QUOTED.
		{
			"ion_invalid_escape_short_string_later_quote",
			"`\"\\q` \"z\"",
			[]Token{
				{tokION_LITERAL, "\"\\q", ast.Loc{Start: 0, End: 5}},
				{tokIDENT_QUOTED, "z", ast.Loc{Start: 6, End: 9}},
			},
		},
		// --- same shape for a single-quoted symbol with a later "'". A buggy
		//     escape would close the symbol at the later quote; the grammar
		//     closes the literal at the first backtick. The trailing "'z'" is
		//     an ordinary single-quoted string. Oracle: `'\q` -> ION_CLOSURE;
		//     'z' -> LITERAL_STRING. ---
		{
			"ion_invalid_escape_symbol_later_quote",
			"`'\\q` 'z'",
			[]Token{
				{tokION_LITERAL, "'\\q", ast.Loc{Start: 0, End: 5}},
				{tokSCONST, "z", ast.Loc{Start: 6, End: 9}},
			},
		},
		// --- a backslash immediately before a backtick is NOT an escape: the
		//     backtick is not a COMMON_ESCAPE_CODE (g4 504-519). So '\`' does
		//     not protect the backtick; the '"' and '\' are ION_ANY and the
		//     backtick closes. Oracle: `"\` -> ION_CLOSURE (content "\); x ->
		//     IDENT. ---
		{
			"ion_backslash_then_backtick_closes",
			"`\"\\`x",
			[]Token{
				{tokION_LITERAL, "\"\\", ast.Loc{Start: 0, End: 4}},
				{tokIDENT, "x", ast.Loc{Start: 4, End: 5}},
			},
		},
		// --- a VALID escape DOES protect the following backtick (regression
		//     guard for the fix above). '\n' is a COMMON_ESCAPE, so the short
		//     string "\n`" stays open across the inner backtick and the literal
		//     closes only at the final backtick. ---
		{
			"ion_valid_escape_protects_backtick",
			"`\"\\n`\"`",
			[]Token{{tokION_LITERAL, "\"\\n`\"", ast.Loc{Start: 0, End: 7}}},
		},
		// --- a VALID escape spelled uppercase still protects: caseInsensitive
		//     (g4 line 4) makes '\N' == '\n' a COMMON_ESCAPE. ---
		{
			"ion_valid_escape_uppercase_protects",
			"`\"\\N`\"`",
			[]Token{{tokION_LITERAL, "\"\\N`\"", ast.Loc{Start: 0, End: 7}}},
		},
		// --- a valid HEX_ESCAPE (\xNN, g4 521-522) protects the backtick. ---
		{
			"ion_valid_hex_escape_protects",
			"`\"\\x41`\"`",
			[]Token{{tokION_LITERAL, "\"\\x41`\"", ast.Loc{Start: 0, End: 9}}},
		},
		// --- RAW NEWLINE inside a SHORT string is NOT valid content. Per
		//     PartiQLLexer.g4 STRING_SHORT_TEXT_ALLOWED (g4 451-456), the
		//     content set starts at U+0020 plus only WS_NOT_NL (tab / vtab /
		//     form-feed, g4 536-541); a raw U+000A (LF) must be escaped. So a
		//     SHORT string carrying a raw newline cannot match: the '"' degrades
		//     to ION_ANY, the raw bytes are ION_ANY, and the FIRST backtick is a
		//     standalone ION_CLOSURE that ends the literal. A buggy scanner that
		//     accepted the newline as content would keep the string open across
		//     the backtick and close at the LATER '"', swallowing the trailing
		//     SQL. Codex round-3 P2 case; verified against the generated ANTLR
		//     lexer oracle: `"a\nb` -> ION_CLOSURE (content "a\nb); "z" ->
		//     IDENTIFIER_QUOTED.
		{
			"ion_raw_newline_short_string_later_quote",
			"`\"a\nb` \"z\"",
			[]Token{
				{tokION_LITERAL, "\"a\nb", ast.Loc{Start: 0, End: 6}},
				{tokIDENT_QUOTED, "z", ast.Loc{Start: 7, End: 10}},
			},
		},
		// --- same shape for a single-quoted SYMBOL: SYMBOL_TEXT_ALLOWED
		//     (g4 494-499) excludes raw newline too, so the "'" degrades to
		//     ION_ANY and the first backtick closes; the trailing "'z'" is an
		//     ordinary single-quoted string. Oracle: `'a\nb` -> ION_CLOSURE;
		//     'z' -> LITERAL_STRING. ---
		{
			"ion_raw_newline_symbol_later_quote",
			"`'a\nb` 'z'",
			[]Token{
				{tokION_LITERAL, "'a\nb", ast.Loc{Start: 0, End: 6}},
				{tokSCONST, "z", ast.Loc{Start: 7, End: 10}},
			},
		},
		// --- raw newline in a SHORT string, simple close (no later quote): the
		//     '"' is ION_ANY and the backtick closes; trailing 'x' is an
		//     identifier. Oracle: `"a\nb` -> ION_CLOSURE; x -> IDENTIFIER. ---
		{
			"ion_raw_newline_short_string_closes",
			"`\"a\nb`x",
			[]Token{
				{tokION_LITERAL, "\"a\nb", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "x", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- raw newline in a SYMBOL, simple close. ---
		{
			"ion_raw_newline_symbol_closes",
			"`'a\nb`x",
			[]Token{
				{tokION_LITERAL, "'a\nb", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "x", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- raw CARRIAGE RETURN (U+000D) is excluded from the SHORT set
		//     exactly like LF (both are ION_NEWLINE, g4 432-436, not in
		//     WS_NOT_NL). The '"' degrades and the backtick closes. ---
		{
			"ion_raw_cr_short_string_closes",
			"`\"a\rb`x",
			[]Token{
				{tokION_LITERAL, "\"a\rb", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "x", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- a raw NUL (U+0000), a C0 control byte excluded from the SHORT set,
		//     also fails the match; the '"' degrades and the backtick closes. ---
		{
			"ion_raw_nul_short_string_closes",
			"`\"a\x00b`x",
			[]Token{
				{tokION_LITERAL, "\"a\x00b", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "x", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- a raw U+001F (unit separator), the highest excluded C0 control
		//     byte, fails the match too. Boundary guard just below U+0020. ---
		{
			"ion_raw_us_short_string_closes",
			"`\"a\x1fb`x",
			[]Token{
				{tokION_LITERAL, "\"a\x1fb", ast.Loc{Start: 0, End: 6}},
				{tokIDENT, "x", ast.Loc{Start: 6, End: 7}},
			},
		},
		// --- a raw TAB (U+0009) IS valid SHORT content: it is in WS_NOT_NL
		//     (g4 536-541). So the string stays open ACROSS the inner backtick
		//     and closes only at the final '"'+backtick. Regression guard that
		//     the control-byte rejection does NOT over-reject the allowed
		//     whitespace. Oracle: `"a\tb`"` -> one ION_CLOSURE (content
		//     "a\tb`"). ---
		{
			"ion_raw_tab_short_string_is_content",
			"`\"a\tb`\"`",
			[]Token{{tokION_LITERAL, "\"a\tb`\"", ast.Loc{Start: 0, End: 8}}},
		},
		// --- a raw VERTICAL TAB (U+000B) is WS_NOT_NL too -> content. ---
		{
			"ion_raw_vtab_short_string_is_content",
			"`\"a\vb`\"`",
			[]Token{{tokION_LITERAL, "\"a\vb`\"", ast.Loc{Start: 0, End: 8}}},
		},
		// --- a raw FORM FEED (U+000C) is WS_NOT_NL too -> content. ---
		{
			"ion_raw_formfeed_short_string_is_content",
			"`\"a\fb`\"`",
			[]Token{{tokION_LITERAL, "\"a\fb`\"", ast.Loc{Start: 0, End: 8}}},
		},
		// --- a raw TAB inside a SYMBOL is likewise WS_NOT_NL content; the
		//     symbol spans the backtick and closes at the final "'"+backtick. ---
		{
			"ion_raw_tab_symbol_is_content",
			"`'a\tb`'`",
			[]Token{{tokION_LITERAL, "'a\tb`'", ast.Loc{Start: 0, End: 8}}},
		},
		// --- a LONG string (''' ''') legitimately SPANS raw newlines: its
		//     STRING_LONG_TEXT_ALLOWED (g4 459-463) admits WS = [ \r\n\t]
		//     (g4 390-391). So a raw newline inside a long string is content,
		//     NOT a match failure — long strings span lines by design. The
		//     long-string control-byte rejection (ion_raw_*_long_string_closes
		//     below) deliberately exempts LF/CR. Regression guard. (Mirrors
		//     ion_long_string_multiline above; pinned here for contrast.) ---
		{
			"ion_raw_newline_long_string_is_content",
			"`'''a\nb'''`",
			[]Token{{tokION_LITERAL, "'''a\nb'''", ast.Loc{Start: 0, End: 11}}},
		},
		// --- a raw CR inside a LONG string is also content (WS includes \r). ---
		{
			"ion_raw_cr_long_string_is_content",
			"`'''a\rb'''`",
			[]Token{{tokION_LITERAL, "'''a\rb'''", ast.Loc{Start: 0, End: 11}}},
		},
		// --- LONG-STRING raw C0 control rejection (Codex round-4 P2). A
		//     triple-quoted long string carrying a raw C0 control byte OUTSIDE
		//     its allowed set must NOT match: the rule fails, ANTLR maximal-munch
		//     degrades the opening quotes through QUOTED_SYMBOL to ION_ANY, and
		//     the FIRST standalone backtick is the ION_CLOSURE. A buggy scanner
		//     that swallowed the control byte as content would keep the long
		//     string open ACROSS that backtick and close only at a LATER '''',
		//     hiding the standalone backtick and swallowing the following SQL.
		//     The long-string allowed set DIFFERS from the SHORT set: long DOES
		//     admit raw LF/CR (newlines, above) and — per the generated ANTLR
		//     lexer oracle — VT/FF too (below), but rejects NUL and the other C0
		//     controls. All boundaries verified byte-for-byte against the
		//     generated ANTLR lexer (github.com/bytebase/parser/partiql).
		// raw NUL (U+0000) inside the long string, simple close: the long-string
		// rule fails on the NUL, the opening quotes degrade through QUOTED_SYMBOL
		// to ION_ANY, and the backtick closes the literal; trailing 'x' is an
		// identifier. Oracle: `'''a\x00b` -> ION_CLOSURE; x -> IDENTIFIER. (The
		// later-'''' variant is the next case.)
		{
			"ion_raw_nul_long_string_closes",
			"`'''a\x00b`x",
			[]Token{
				{tokION_LITERAL, "'''a\x00b", ast.Loc{Start: 0, End: 8}},
				{tokIDENT, "x", ast.Loc{Start: 8, End: 9}},
			},
		},
		// raw NUL with a LATER '''' present: a buggy span would close at that
		// later '''' (inside the second literal below), hiding the standalone
		// backtick; the correct close is the FIRST backtick. This is the exact
		// finding shape — the standalone backtick must NOT be swallowed. After
		// the close: '+' is an operator and `'''x'''` is a SEPARATE, complete
		// Ion literal. Oracle: ION_CLOSURE[0:8], PLUS, ION_CLOSURE[9:18].
		{
			"ion_raw_nul_long_string_later_triple_quote",
			"`'''a\x00b`+`'''x'''`",
			[]Token{
				{tokION_LITERAL, "'''a\x00b", ast.Loc{Start: 0, End: 8}},
				{tokPLUS, "+", ast.Loc{Start: 8, End: 9}},
				{tokION_LITERAL, "'''x'''", ast.Loc{Start: 9, End: 18}},
			},
		},
		// raw FORM FEED (U+000C) variant: FF IS allowed long-string content
		// (oracle), so the long string SPANS the inner backtick and closes at
		// the LATER ''''+backtick — the opposite outcome from NUL. This is the
		// form-feed contrast case requested in the finding; it guards that the
		// rejection does NOT over-reject FF.
		{
			"ion_raw_formfeed_long_string_is_content",
			"`'''a\x0cb`'''`x",
			[]Token{
				{tokION_LITERAL, "'''a\x0cb`'''", ast.Loc{Start: 0, End: 12}},
				{tokIDENT, "x", ast.Loc{Start: 12, End: 13}},
			},
		},
		// raw VERTICAL TAB (U+000B) is likewise allowed (oracle) -> spans. ---
		{
			"ion_raw_vtab_long_string_is_content",
			"`'''a\vb`'''`x",
			[]Token{
				{tokION_LITERAL, "'''a\vb`'''", ast.Loc{Start: 0, End: 12}},
				{tokIDENT, "x", ast.Loc{Start: 12, End: 13}},
			},
		},
		// raw U+001F (unit separator), the highest disallowed C0 control byte —
		// boundary guard just below U+0020. Closes at the first backtick. ---
		{
			"ion_raw_us_long_string_closes",
			"`'''a\x1fb`x",
			[]Token{
				{tokION_LITERAL, "'''a\x1fb", ast.Loc{Start: 0, End: 8}},
				{tokIDENT, "x", ast.Loc{Start: 8, End: 9}},
			},
		},
		// raw U+0001 (start of heading), the lowest disallowed C0 control byte
		// above NUL — other boundary guard. Closes at the first backtick. ---
		{
			"ion_raw_soh_long_string_closes",
			"`'''a\x01b`x",
			[]Token{
				{tokION_LITERAL, "'''a\x01b", ast.Loc{Start: 0, End: 8}},
				{tokIDENT, "x", ast.Loc{Start: 8, End: 9}},
			},
		},
		// a normal multi-line long string carrying \r \n \t together still
		// closes at its OWN trailing '''' (one literal), unaffected by the
		// rejection — the core regression guard for spanning newlines. ---
		{
			"ion_multiline_long_string_crlf_tab",
			"`'''l1\r\n\tl2'''`x",
			[]Token{
				{tokION_LITERAL, "'''l1\r\n\tl2'''", ast.Loc{Start: 0, End: 15}},
				{tokIDENT, "x", ast.Loc{Start: 15, End: 16}},
			},
		},
		// --- ROUND-5 DISPUTED-SHAPE settlement (VT/FF). Codex argued VT (0x0B)
		//     and FF (0x0C) should be REJECTED from the long-string content set
		//     (grammar text WS=[space \r \n \t]). The EXECUTABLE ORACLE (generated
		//     ANTLR PartiQLLexer, authoritative over grammar text) was queried
		//     differentially with the EXACT disputed shape below — a long string
		//     carrying a raw control byte before a standalone backtick, then a
		//     LATER ''' and the closing backtick, with a trailing identifier:
		//       in: `'''<ctrl>X`Y''' ` TAIL
		//     Oracle result per control byte (the single ION_CLOSURE token text):
		//       CR  0x0D (allowed)    -> `'''\rX`Y''' `  (HELD; spans inner `)
		//       NUL 0x00 (disallowed) -> `'''\x00X`      (FAILED; first ` closes)
		//       VT  0x0B (disputed)   -> `'''\vX`Y''' `  (HELD; spans inner `)
		//       FF  0x0C (disputed)   -> `'''\fX`Y''' `  (HELD; spans inner `)
		//     VT/FF behave EXACTLY like the indisputably-allowed CR and OPPOSITE
		//     the disallowed NUL: the long string KEEPS them as content and spans
		//     the inner backtick, closing at the LATER '''+backtick. So omni's
		//     round-4 code is CORRECT and Codex's reading is REBUTTED. These cases
		//     pin that oracle verdict in code. (Note: the inner "Y''' " yields,
		//     after the closing backtick, an IDENTIFIER "Y" then an EMPTY symbol
		//     literal `''' ` — three quotes = one empty long string with a leading
		//     WS — then the standalone backtick; we assert only through the held
		//     ION_LITERAL, which is the disputed boundary.)
		// VT held: literal spans to the LATER '''+backtick (content is
		// input[1:12] = "'''<VT>X`Y''' "; closing backtick at index 12).
		{
			"ion_round5_vtab_disputed_shape_held",
			"`'''\vX`Y''' `Z",
			[]Token{
				{tokION_LITERAL, "'''\vX`Y''' ", ast.Loc{Start: 0, End: 13}},
				{tokIDENT, "Z", ast.Loc{Start: 13, End: 14}},
			},
		},
		// FF held: literal spans to the LATER '''+backtick.
		{
			"ion_round5_formfeed_disputed_shape_held",
			"`'''\fX`Y''' `Z",
			[]Token{
				{tokION_LITERAL, "'''\fX`Y''' ", ast.Loc{Start: 0, End: 13}},
				{tokIDENT, "Z", ast.Loc{Start: 13, End: 14}},
			},
		},
		// CR (indisputably allowed) bracket: same shape, also HELD — confirms the
		// disputed-shape harness logic, not just VT/FF.
		{
			"ion_round5_cr_bracket_held",
			"`'''\rX`Y''' `Z",
			[]Token{
				{tokION_LITERAL, "'''\rX`Y''' ", ast.Loc{Start: 0, End: 13}},
				{tokIDENT, "Z", ast.Loc{Start: 13, End: 14}},
			},
		},
		// NUL (indisputably disallowed) bracket: same shape, but FAILED — the
		// FIRST standalone backtick (index 6) closes the literal (content
		// input[1:6] = "'''<NUL>X"); the contrast that proves VT/FF/CR genuinely
		// HELD above rather than the scanner being permissive. After the close
		// the trailing 'Y' is a top-level IDENTIFIER.
		{
			"ion_round5_nul_bracket_closes_first_backtick",
			"`'''\x00X`Y",
			[]Token{
				{tokION_LITERAL, "'''\x00X", ast.Loc{Start: 0, End: 7}},
				{tokIDENT, "Y", ast.Loc{Start: 7, End: 8}},
			},
		},
		// --- a real Ion timestamp value (PartiQL spec example) ---
		{
			"ion_timestamp",
			"`2017-09-14T`",
			[]Token{{tokION_LITERAL, "2017-09-14T", ast.Loc{Start: 0, End: 13}}},
		},
		// --- adjacent backtick after a closed symbol opens a new literal ---
		{
			"ion_two_literals_back_to_back",
			"`'a'``'b'`",
			[]Token{
				{tokION_LITERAL, "'a'", ast.Loc{Start: 0, End: 5}},
				{tokION_LITERAL, "'b'", ast.Loc{Start: 5, End: 10}},
			},
		},
		// --- a struct with a string field containing a backtick ---
		{
			"ion_struct_with_backtick_string",
			"`{a: \"x`y\"}`",
			[]Token{{tokION_LITERAL, "{a: \"x`y\"}", ast.Loc{Start: 0, End: 12}}},
		},
		// --- empty quoted symbol then close ---
		{
			"ion_empty_symbol",
			"`''`",
			[]Token{{tokION_LITERAL, "''", ast.Loc{Start: 0, End: 4}}},
		},
		// --- empty short string then close ---
		{
			"ion_empty_short_string",
			"`\"\"`",
			[]Token{{tokION_LITERAL, "\"\"", ast.Loc{Start: 0, End: 4}}},
		},
		// --- single { is not a lob; struct braces pass through as ION_ANY ---
		{
			"ion_single_brace_struct",
			"`{ }`",
			[]Token{{tokION_LITERAL, "{ }", ast.Loc{Start: 0, End: 5}}},
		},
		// --- a lone '/' (not // or /*) is just ION_ANY content ---
		{
			"ion_lone_slash",
			"`a/b`",
			[]Token{{tokION_LITERAL, "a/b", ast.Loc{Start: 0, End: 5}}},
		},
		// --- a // line comment terminated by a newline; the backtick on
		//     the following line closes the literal. Grammar: the inline
		//     comment ends at ION_NEWLINE (g4 408-409). ---
		{
			"ion_line_comment_then_newline_then_close",
			"`a // c\n`",
			[]Token{{tokION_LITERAL, "a // c\n", ast.Loc{Start: 0, End: 9}}},
		},
		// --- ANTLR maximal-munch FALLBACK to ION_ANY. A quote/comment
		//     opener that never finds its closer does NOT match its
		//     multi-byte rule; the opener degrades to a single-byte
		//     ION_ANY and the next standalone backtick still closes the
		//     literal. These were previously (incorrectly) treated as
		//     unterminated errors; the generated ANTLR lexer oracle accepts
		//     each as a single literal. ---
		// open " with no closing " before the backtick: '"' is ION_ANY.
		{
			"ion_unclosed_short_string_falls_back",
			"`\"abc`",
			[]Token{{tokION_LITERAL, "\"abc", ast.Loc{Start: 0, End: 6}}},
		},
		// open ' with no closing ' before the backtick: "'" is ION_ANY.
		{
			"ion_unclosed_symbol_falls_back",
			"`'abc`",
			[]Token{{tokION_LITERAL, "'abc", ast.Loc{Start: 0, End: 6}}},
		},
		// open ''' with no closing ''': degrades through symbol to ION_ANY;
		// the backtick closes. Oracle content: '''abc.
		{
			"ion_unclosed_long_string_falls_back",
			"`'''abc`",
			[]Token{{tokION_LITERAL, "'''abc", ast.Loc{Start: 0, End: 8}}},
		},
		// open /* with no closing */: '/' is ION_ANY; the backtick closes.
		{
			"ion_unclosed_block_comment_falls_back",
			"`/* abc`",
			[]Token{{tokION_LITERAL, "/* abc", ast.Loc{Start: 0, End: 8}}},
		},
		// open {{ with no closing }} (and no quote): braces are ION_ANY; the
		// backtick closes. (ION_BLOB needs base64 + }}, which is absent.)
		{
			"ion_unclosed_lob_falls_back",
			"`{{ abc`",
			[]Token{{tokION_LITERAL, "{{ abc", ast.Loc{Start: 0, End: 8}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runTokenStreamCase(t, tc)
		})
	}
}

// TestLexer_IonMode_Errors covers the reject forms: an Ion literal that
// runs to EOF with no standalone closing backtick. This happens either
// because there simply is no further backtick, or because an inner
// construct that DID match (a closing-quote-less string is NOT such a
// construct — see TestLexer_IonMode fallbacks) swallowed the backtick:
// a // line comment running to EOF, or a trailing backslash escape inside
// a string. The generated ANTLR lexer oracle produces zero tokens for
// each of these (the ION mode is never popped). NOTE: an unclosed string
// / symbol / block-comment / lob whose opener is followed by a later
// standalone backtick is NOT an error — ANTLR maximal-munch degrades the
// opener to ION_ANY and the backtick closes; those accept cases live in
// TestLexer_IonMode.
func TestLexer_IonMode_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "unterminated_bare",
			input:     "`abc",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// A trailing backslash inside an open symbol consumes the EOF
			// boundary; no closing ' and no later backtick, so the literal
			// never closes. Oracle: zero tokens.
			name:      "unterminated_trailing_escape",
			input:     "`'abc\\",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// '''  with nothing after it: no closing ''', no later
			// backtick anywhere, so the literal runs to EOF unterminated.
			name:      "unterminated_open_triple_quote",
			input:     "`'''",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// '  with nothing after it: no closing ' and no later backtick,
			// so the literal runs to EOF unterminated.
			name:      "unterminated_open_single_quote",
			input:     "`'",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// A // line comment with no trailing newline runs to EOF,
			// swallowing the backtick that would otherwise have closed the
			// literal (g4 408-409: '//' .*? (NL|EOF)). Oracle: zero tokens.
			name:      "unterminated_line_comment_eats_closer",
			input:     "`a // c`",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// // comment to EOF with no backtick at all: also unterminated.
			name:      "unterminated_line_comment_no_closer",
			input:     "`a//b",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// '{{ //// ' is NOT a valid ION_BLOB (no closing '}}' LOB_END), so
			// the braces are ION_ANY and the '//' starts an ION_INLINE_COMMENT
			// that runs to EOF, swallowing the closing backtick. Distinguishes
			// "valid blob shields '//'" from "invalid blob leaves '//' a
			// comment". Oracle: zero tokens (mode never popped).
			name:      "unterminated_invalid_blob_slash_comment",
			input:     "`{{ //// ",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// '{{//}}' is NOT a valid ION_BLOB ('//' is only two base64 chars,
			// not a BASE_64_QUARTET and not a BASE_64_PAD), so the '//' begins
			// an ION_INLINE_COMMENT that eats '}}' + the backtick to EOF.
			// Oracle: zero tokens.
			name:      "unterminated_short_blob_slash_comment",
			input:     "`{{//}}`",
			wantErrIn: "unterminated Ion literal",
		},
		{
			// A valid ION_BLOB is consumed whole, but a '//' that appears
			// AFTER the lob is an ordinary ION_INLINE_COMMENT: it runs to EOF
			// and swallows the closing backtick. Confirms the blob matcher
			// only shields '//' that is base64 CONTENT, not '//' outside the
			// lob. Oracle: zero tokens.
			name:      "unterminated_comment_after_blob",
			input:     "`{{ //// }} // c",
			wantErrIn: "unterminated Ion literal",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			for {
				tok := l.Next()
				if tok.Type == tokEOF {
					break
				}
			}
			if l.Err == nil {
				t.Errorf("expected error, got nil")
				return
			}
			if !strings.Contains(l.Err.Error(), tc.wantErrIn) {
				t.Errorf("error mismatch\n got: %v\nwant substring: %q", l.Err, tc.wantErrIn)
			}
		})
	}
}

// TestLexer_IonMode_RoundTrip is the structural gate for the Ion
// literal scanner. No reference parser exists for partiql, so it relies
// on source-reconstruction stability: for any accepted Ion literal,
// wrapping Token.Str back in backticks must reproduce the exact original
// source, and re-lexing that reconstruction must yield an identical
// token. This catches off-by-one slicing and any silent content
// mutation in the sub-mode scanners.
func TestLexer_IonMode_RoundTrip(t *testing.T) {
	inputs := []string{
		"`{a: 1}`",
		"``",
		"`  abc  `",
		"`'hello'`",
		"`'a`b'`",
		"`\"x`y\"`",
		"`'''a`b'''`",
		"`'''a\nb'''`",
		"`// c`omment\n42`",
		"`/* c`c */1`",
		"`{{ aGVsbG8= }}`",
		"`{{ //// }}`",
		"`{{////}}`",
		"`{{ \"a//b\" }}`",
		"`{{ \"}}\" }}`",
		"`{{ '''}}''' }}`",
		"`'a\\'`b'`",
		"`2017-09-14T`",
		"`{a: \"x`y\"}`",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			l := NewLexer(in)
			tok := l.Next()
			if l.Err != nil {
				t.Fatalf("unexpected error: %v", l.Err)
			}
			if tok.Type != tokION_LITERAL {
				t.Fatalf("token type = %s, want ION_LITERAL", tokenName(tok.Type))
			}
			// The token must span the whole input and reconstruct it.
			if got := in[tok.Loc.Start:tok.Loc.End]; got != in {
				t.Errorf("loc slice = %q, want full input %q", got, in)
			}
			recon := "`" + tok.Str + "`"
			if recon != in {
				t.Errorf("reconstruction = %q, want %q", recon, in)
			}
			// Re-lex the reconstruction: it must yield an identical token.
			l2 := NewLexer(recon)
			tok2 := l2.Next()
			if l2.Err != nil {
				t.Fatalf("re-lex error: %v", l2.Err)
			}
			if tok2 != tok {
				t.Errorf("re-lex token mismatch\n got: %+v\nwant: %+v", tok2, tok)
			}
			if next := l2.Next(); next.Type != tokEOF {
				t.Errorf("expected EOF after re-lexed literal, got %s", tokenName(next.Type))
			}
		})
	}
}

// TestLexer_AWSCorpus loads every .partiql file from
// partiql/parser/testdata/aws-corpus/, filters out the 2 known-bad
// syntax-skeleton files, and asserts each one lexes to a non-error
// token stream ending with tokEOF. Catches "does the lexer tokenize
// at all" regressions on real AWS DynamoDB PartiQL examples.
//
// Skipped files:
//   - select-001.partiql: SELECT syntax skeleton with bracket placeholders
//   - insert-002.partiql: INSERT syntax skeleton with backtick placeholder
//
// Both are flagged in testdata/aws-corpus/index.json as not-real-PartiQL.
// The skip list is hard-coded here for clarity.
func TestLexer_AWSCorpus(t *testing.T) {
	skip := map[string]bool{
		"select-001.partiql": true,
		"insert-002.partiql": true,
	}

	files, err := filepath.Glob("testdata/aws-corpus/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no corpus files found — testdata/aws-corpus/ missing or empty?")
	}

	var lexed, skipped int
	for _, f := range files {
		name := filepath.Base(f)
		if skip[name] {
			skipped++
			continue
		}
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			l := NewLexer(string(data))
			tokens := 0
			for {
				tok := l.Next()
				if tok.Type == tokEOF {
					break
				}
				tokens++
				if tokens > 100000 {
					t.Fatalf("token stream did not terminate after %d tokens", tokens)
				}
			}
			if l.Err != nil {
				t.Errorf("lexer error: %v", l.Err)
			}
			if tokens == 0 {
				t.Errorf("lexed to zero tokens")
			}
		})
		lexed++
	}
	t.Logf("AWS corpus: %d files lexed, %d skipped", lexed, skipped)
}

// TestLexer_Errors covers the 5 error triggers in the lexer. Each case
// drains Next() until tokEOF and asserts l.Err is set with the expected
// error message substring.
func TestLexer_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string // substring of the expected error message
	}{
		{
			name:      "unterminated_string",
			input:     "'hello",
			wantErrIn: "unterminated string literal",
		},
		{
			name:      "unterminated_quoted_ident",
			input:     `"foo`,
			wantErrIn: "unterminated quoted identifier",
		},
		{
			name:      "unterminated_ion_literal",
			input:     "`abc",
			wantErrIn: "unterminated Ion literal",
		},
		{
			name:      "unterminated_block_comment",
			input:     "/* nope",
			wantErrIn: "unterminated block comment",
		},
		{
			name:      "unexpected_character_lone_bang",
			input:     "! 1",
			wantErrIn: "unexpected character",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			for {
				tok := l.Next()
				if tok.Type == tokEOF {
					break
				}
			}
			if l.Err == nil {
				t.Errorf("expected error, got nil")
				return
			}
			if !strings.Contains(l.Err.Error(), tc.wantErrIn) {
				t.Errorf("error mismatch\n got: %v\nwant substring: %q", l.Err, tc.wantErrIn)
			}
		})
	}
}

// TestTokenName_AllCovered walks every tok* constant declared in token.go
// and asserts tokenName returns a non-empty string. If a future
// contributor adds a new tok* constant without wiring it into tokenName,
// this test fails.
//
// LIMITATION: this list is a manual mirror of token.go. The length check
// catches slice drift, and the loop catches missing tokenName arms — but
// neither catches a constant that was added to BOTH token.go and tokenName
// without being added here. New tok* constants must be added to all three
// places.
func TestTokenName_AllCovered(t *testing.T) {
	// Total: 2 specials + 6 literals + 28 operators/punctuation + 266 keywords = 302.
	all := []int{
		// Specials.
		tokEOF, tokInvalid,
		// Literals.
		tokSCONST, tokICONST, tokFCONST, tokIDENT, tokIDENT_QUOTED, tokION_LITERAL,
		// Operators / punctuation.
		tokPLUS, tokMINUS, tokASTERISK, tokSLASH_FORWARD, tokPERCENT,
		tokCARET, tokTILDE, tokAT_SIGN, tokEQ, tokNEQ, tokLT, tokGT,
		tokLT_EQ, tokGT_EQ, tokCONCAT, tokANGLE_DOUBLE_LEFT, tokANGLE_DOUBLE_RIGHT,
		tokPAREN_LEFT, tokPAREN_RIGHT, tokBRACKET_LEFT, tokBRACKET_RIGHT,
		tokBRACE_LEFT, tokBRACE_RIGHT, tokCOLON, tokCOLON_SEMI, tokCOMMA,
		tokPERIOD, tokQUESTION_MARK,
		// Keywords (alphabetical).
		tokABSOLUTE, tokACTION, tokADD, tokALL, tokALLOCATE, tokALTER, tokAND,
		tokANY, tokARE, tokAS, tokASC, tokASSERTION, tokAT, tokAUTHORIZATION,
		tokAVG, tokBAG, tokBEGIN, tokBETWEEN, tokBIGINT, tokBIT, tokBIT_LENGTH,
		tokBLOB, tokBOOL, tokBOOLEAN, tokBY, tokCAN_CAST, tokCAN_LOSSLESS_CAST,
		tokCASCADE, tokCASCADED, tokCASE, tokCAST, tokCATALOG, tokCHAR,
		tokCHAR_LENGTH, tokCHARACTER, tokCHARACTER_LENGTH, tokCHECK, tokCLOB,
		tokCLOSE, tokCOALESCE, tokCOLLATE, tokCOLLATION, tokCOLUMN, tokCOMMIT,
		tokCONFLICT, tokCONNECT, tokCONNECTION, tokCONSTRAINT, tokCONSTRAINTS,
		tokCONTINUE, tokCONVERT, tokCORRESPONDING, tokCOUNT, tokCREATE, tokCROSS,
		tokCURRENT, tokCURRENT_DATE, tokCURRENT_TIME, tokCURRENT_TIMESTAMP,
		tokCURRENT_USER, tokCURSOR, tokDATE, tokDATE_ADD, tokDATE_DIFF,
		tokDEALLOCATE, tokDEC, tokDECIMAL, tokDECLARE, tokDEFAULT, tokDEFERRABLE,
		tokDEFERRED, tokDELETE, tokDESC, tokDESCRIBE, tokDESCRIPTOR,
		tokDIAGNOSTICS, tokDISCONNECT, tokDISTINCT, tokDO, tokDOMAIN, tokDOUBLE,
		tokDROP, tokELSE, tokEND, tokEND_EXEC, tokESCAPE, tokEXCEPT, tokEXCEPTION,
		tokEXCLUDED, tokEXEC, tokEXECUTE, tokEXISTS, tokEXPLAIN, tokEXTERNAL,
		tokEXTRACT, tokFALSE, tokFETCH, tokFIRST, tokFLOAT, tokFOR, tokFOREIGN,
		tokFOUND, tokFROM, tokFULL, tokGET, tokGLOBAL, tokGO, tokGOTO, tokGRANT,
		tokGROUP, tokHAVING, tokIDENTITY, tokIMMEDIATE, tokIN, tokINDEX,
		tokINDICATOR, tokINITIALLY, tokINNER, tokINPUT, tokINSENSITIVE, tokINSERT,
		tokINT, tokINT2, tokINT4, tokINT8, tokINTEGER, tokINTEGER2, tokINTEGER4,
		tokINTEGER8, tokINTERSECT, tokINTERVAL, tokINTO, tokIS, tokISOLATION,
		tokJOIN, tokKEY, tokLAG, tokLANGUAGE, tokLAST, tokLATERAL, tokLEAD,
		tokLEFT, tokLET, tokLEVEL, tokLIKE, tokLIMIT, tokLIST, tokLOCAL, tokLOWER,
		tokMATCH, tokMAX, tokMIN, tokMISSING, tokMODIFIED, tokMODULE, tokNAMES,
		tokNATIONAL, tokNATURAL, tokNCHAR, tokNEW, tokNEXT, tokNO, tokNOT,
		tokNOTHING, tokNULL, tokNULLIF, tokNULLS, tokNUMERIC, tokOCTET_LENGTH,
		tokOF, tokOFFSET, tokOLD, tokON, tokONLY, tokOPEN, tokOPTION, tokOR,
		tokORDER, tokOUTER, tokOUTPUT, tokOVER, tokOVERLAPS, tokOVERLAY, tokPAD,
		tokPARTIAL, tokPARTITION, tokPIVOT, tokPLACING, tokPOSITION, tokPRECISION,
		tokPREPARE, tokPRESERVE, tokPRIMARY, tokPRIOR, tokPRIVILEGES, tokPROCEDURE,
		tokPUBLIC, tokREAD, tokREAL, tokREFERENCES, tokRELATIVE, tokREMOVE,
		tokREPLACE, tokRESTRICT, tokRETURNING, tokREVOKE, tokRIGHT, tokROLLBACK,
		tokROWS, tokSCHEMA, tokSCROLL, tokSECTION, tokSELECT, tokSESSION,
		tokSESSION_USER, tokSET, tokSEXP, tokSHORTEST, tokSIZE, tokSMALLINT,
		tokSOME, tokSPACE, tokSQL, tokSQLCODE, tokSQLERROR, tokSQLSTATE, tokSTRING,
		tokSTRUCT, tokSUBSTRING, tokSUM, tokSYMBOL, tokSYSTEM_USER, tokTABLE,
		tokTEMPORARY, tokTHEN, tokTIME, tokTIMESTAMP, tokTO, tokTRANSACTION,
		tokTRANSLATE, tokTRANSLATION, tokTRIM, tokTRUE, tokTUPLE, tokUNION,
		tokUNIQUE, tokUNKNOWN, tokUNPIVOT, tokUPDATE, tokUPPER, tokUPSERT,
		tokUSAGE, tokUSER, tokUSING, tokVALUE, tokVALUES, tokVARCHAR, tokVARYING,
		tokVIEW, tokWHEN, tokWHENEVER, tokWHERE, tokWITH, tokWORK, tokWRITE,
		tokZONE,
	}
	if got := len(all); got != 302 {
		t.Errorf("test list has %d entries, want 302 — did a tok* constant get added or removed without updating this test?", got)
	}
	for _, tt := range all {
		name := tokenName(tt)
		if name == "" {
			t.Errorf("tokenName(%d) returned empty string — missing switch arm in token.go?", tt)
		}
	}
}

// TestKeywords_LenMatchesConstants asserts that the keywords map in
// keywords.go has exactly 266 entries — the same number as the keyword
// constants in token.go — and that every map value resolves to a
// non-empty tokenName. The length check catches add/remove drift; the
// per-value check catches the case where a keyword maps to a deleted
// or renamed constant (which leaves the lengths equal but the mapping
// stale).
func TestKeywords_LenMatchesConstants(t *testing.T) {
	const expectedKeywordCount = 266
	if got := len(keywords); got != expectedKeywordCount {
		t.Errorf("len(keywords) = %d, want %d — did a tok* keyword constant get added or removed without updating the keywords map?", got, expectedKeywordCount)
	}
	for word, tt := range keywords {
		if tokenName(tt) == "" {
			t.Errorf("keywords[%q] = %d has no tokenName entry", word, tt)
		}
	}
}
