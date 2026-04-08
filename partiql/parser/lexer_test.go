package parser

import (
	"os"
	"path/filepath"
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
// Tasks 6-10 each append a section of cases as their scan helper lands.
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
		// Ion literals — base lexer (Task 10)
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
		// Pins the simplified scanner's behavior on Ion content containing
		// a single quote. Once Ion-mode-aware scanning lands (DAG node 17),
		// this case will continue to pass — the inner ' characters are not
		// the literal terminator. Documents the current behavior in
		// executable form for the future refactor.
		{
			"ion_with_single_quote_known_limitation",
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
