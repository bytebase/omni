package parser

import "testing"

// expect describes one expected token in the output stream. Loc is omitted
// here — position assertions live in position_test.go.
type expect struct {
	Type int
	Str  string
}

// runLexerCases drives a table-driven test. For each row, it tokenizes
// row.input and asserts the resulting non-EOF tokens match row.want by
// (Type, Str). Errors must be empty.
func runLexerCases(t *testing.T, cases []struct {
	name  string
	input string
	want  []expect
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			// Strip the trailing EOF token before comparison.
			if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
				t.Fatalf("expected stream to end with tokEOF, got %+v", tokens)
			}
			tokens = tokens[:len(tokens)-1]
			if len(tokens) != len(c.want) {
				t.Fatalf("token count mismatch: got %d, want %d\n got=%+v\nwant=%+v",
					len(tokens), len(c.want), tokens, c.want)
			}
			for i, tok := range tokens {
				if tok.Type != c.want[i].Type {
					t.Errorf("token[%d] Type = %d (%s), want %d (%s)",
						i, tok.Type, TokenName(tok.Type), c.want[i].Type, TokenName(c.want[i].Type))
				}
				if tok.Str != c.want[i].Str {
					t.Errorf("token[%d] Str = %q, want %q", i, tok.Str, c.want[i].Str)
				}
			}
		})
	}
}

func TestLexer_SingleCharOperators(t *testing.T) {
	// GoogleSQLLexer.g4 single-char *_OPERATOR / *_SYMBOL tokens that have no
	// multi-char extension at this byte.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"plus", "+", []expect{{int('+'), ""}}},
		{"minus", "-", []expect{{int('-'), ""}}},
		{"star", "*", []expect{{int('*'), ""}}},
		{"slash", "/", []expect{{int('/'), ""}}},
		{"percent", "%", []expect{{int('%'), ""}}},
		{"tilde", "~", []expect{{int('~'), ""}}},
		{"lparen", "(", []expect{{int('('), ""}}},
		{"rparen", ")", []expect{{int(')'), ""}}},
		{"lbracket", "[", []expect{{int('['), ""}}},
		{"rbracket", "]", []expect{{int(']'), ""}}},
		{"lbrace", "{", []expect{{int('{'), ""}}},
		{"rbrace", "}", []expect{{int('}'), ""}}},
		{"comma", ",", []expect{{int(','), ""}}},
		{"semi", ";", []expect{{int(';'), ""}}},
		{"dot", ".", []expect{{int('.'), ""}}},
		{"at", "@", []expect{{int('@'), ""}}},
		{"colon", ":", []expect{{int(':'), ""}}},
		{"lt", "<", []expect{{int('<'), ""}}},
		{"gt", ">", []expect{{int('>'), ""}}},
		{"eq", "=", []expect{{int('='), ""}}},
		{"bang", "!", []expect{{int('!'), ""}}},
		{"caret", "^", []expect{{int('^'), ""}}},
		{"ampersand", "&", []expect{{int('&'), ""}}},
		{"question", "?", []expect{{int('?'), ""}}},
		{"stroke", "|", []expect{{int('|'), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_MultiCharOperators(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"not_eq_bang", "!=", []expect{{tokNotEqual, ""}}},
		{"not_eq_angle", "<>", []expect{{tokNotEqual2, ""}}},
		{"less_eq", "<=", []expect{{tokLessEqual, ""}}},
		{"greater_eq", ">=", []expect{{tokGreaterEqual, ""}}},
		{"shift_left", "<<", []expect{{tokShiftLeft, ""}}},
		{"shift_right", ">>", []expect{{tokShiftRight, ""}}},
		{"arrow", "->", []expect{{tokArrow, ""}}},
		{"fat_arrow", "=>", []expect{{tokFatArrow, ""}}},
		{"plus_eq", "+=", []expect{{tokPlusEqual, ""}}},
		{"minus_eq", "-=", []expect{{tokMinusEqual, ""}}},
		{"pipe_op", "|>", []expect{{tokPipe, ""}}},
		{"bool_or", "||", []expect{{tokBoolOr, ""}}},
		{"atat", "@@", []expect{{tokAtAt, ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_OperatorBoundaries(t *testing.T) {
	// Ensure greedy multi-char lexing does not over-consume.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		// "<<=" must lex as "<<" then "=" (KL_OPERATOR is the longest prefix).
		{"shift_left_then_eq", "<<=", []expect{{tokShiftLeft, ""}, {int('='), ""}}},
		// ">>=" -> ">>" then "="
		{"shift_right_then_eq", ">>=", []expect{{tokShiftRight, ""}, {int('='), ""}}},
		// "|>>" -> "|>" then ">"
		{"pipe_then_gt", "|>>", []expect{{tokPipe, ""}, {int('>'), ""}}},
		// "->>" -> "->" then ">" (GoogleSQL has no FLOW token, unlike Snowflake).
		{"arrow_then_gt", "->>", []expect{{tokArrow, ""}, {int('>'), ""}}},
		// "@@@" -> "@@" then "@"
		{"atat_then_at", "@@@", []expect{{tokAtAt, ""}, {int('@'), ""}}},
		// "===" -> "=" "=" "=" (no "==" token)
		{"triple_eq", "===", []expect{{int('='), ""}, {int('='), ""}, {int('='), ""}}},
		// "||>" -> "||" then ">" (BOOL_OR is the longest prefix, not |>).
		{"bool_or_then_gt", "||>", []expect{{tokBoolOr, ""}, {int('>'), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_NestedGenericClose(t *testing.T) {
	// PARSER-HANDOFF CONTRACT: a nested generic type close like
	// ARRAY<ARRAY<INT64>> produces a single tokShiftRight ('>>') token, NOT two
	// '>' tokens — matching the legacy GoogleSQLLexer.g4 (KR_OPERATOR: '>>' is
	// greedy) and ZetaSQL's tokenizer. The parser's template_type_close rule is
	// GT_OPERATOR (a single '>'), so the downstream types/expressions node must
	// SPLIT a '>>' (and '>=') at a type-close position into '>' + remainder.
	//
	// VERIFIED against the Spanner emulator: "SELECT CAST(x AS ARRAY<ARRAY<INT64>>)"
	// parses (semantic "Arrays of arrays are not supported"), confirming the
	// real tokenizer+parser accept the '>>' close. The lexer's job is only to
	// emit '>>' as one token; the split lives in the parser.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"array_array", "ARRAY<ARRAY<INT64>>", []expect{
			{kwARRAY, "ARRAY"}, {int('<'), ""}, {kwARRAY, "ARRAY"}, {int('<'), ""},
			{tokIdentifier, "INT64"}, {tokShiftRight, ""},
		}},
		{"single_close", "ARRAY<INT64>", []expect{
			{kwARRAY, "ARRAY"}, {int('<'), ""}, {tokIdentifier, "INT64"}, {int('>'), ""},
		}},
	}
	runLexerCases(t, cases)
}

func TestLexer_UnquotedIdentifiers(t *testing.T) {
	// GoogleSQLLexer.g4 UNQUOTED_IDENTIFIER: [A-Z_][A-Z0-9_]* (case-insensitive).
	// Unlike Snowflake, @ and $ are NOT identifier continuation characters.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple_lower", "abc", []expect{{tokIdentifier, "abc"}}},
		{"simple_upper", "ABC", []expect{{tokIdentifier, "ABC"}}},
		{"underscore_start", "_x", []expect{{tokIdentifier, "_x"}}},
		{"with_digits", "x123", []expect{{tokIdentifier, "x123"}}},
		{"all_underscore", "___", []expect{{tokIdentifier, "___"}}},
		// '$' is not an identifier char in GoogleSQL: "x$y" -> "x", then bad byte.
		// (covered in error tests) — here we check '@' splits identifiers.
		{"at_splits", "x@y", []expect{{tokIdentifier, "x"}, {int('@'), ""}, {tokIdentifier, "y"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_BacktickIdentifiers(t *testing.T) {
	// BQTEXT: backtick-quoted identifier. Body preserves source bytes verbatim
	// (Str excludes the surrounding backticks). Escapes via ANY_ESCAPE (\x).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple", "`my_col`", []expect{{tokIdentifier, "my_col"}}},
		{"with_space", "`my col`", []expect{{tokIdentifier, "my col"}}},
		{"with_dash", "`my-proj`", []expect{{tokIdentifier, "my-proj"}}},
		{"with_dot", "`a.b.c`", []expect{{tokIdentifier, "a.b.c"}}},
		{"with_keyword", "`select`", []expect{{tokIdentifier, "select"}}},
		{"empty", "``", []expect{{tokIdentifier, ""}}},
		{"escaped_backtick", "`a\\`b`", []expect{{tokIdentifier, "a\\`b"}}},
		{"escaped_backslash", "`a\\\\b`", []expect{{tokIdentifier, "a\\\\b"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_StringLiterals(t *testing.T) {
	// STRING_LITERAL: R? (SQTEXT | DQTEXT | SQ3TEXT | DQ3TEXT). Str preserves
	// the verbatim body (no escape processing at the lexer layer — escapes are
	// resolved by the parser/analyzer). Quotes are stripped.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"single", `'hello'`, []expect{{tokString, "hello"}}},
		{"single_empty", `''`, []expect{{tokString, ""}}},
		{"double", `"hello"`, []expect{{tokString, "hello"}}},
		{"double_empty", `""`, []expect{{tokString, ""}}},
		// Backslash escapes are kept verbatim in Str (lexer does not unescape).
		{"single_escape_quote", `'a\'b'`, []expect{{tokString, `a\'b`}}},
		{"double_escape_quote", `"a\"b"`, []expect{{tokString, `a\"b`}}},
		{"single_escape_n", `'a\nb'`, []expect{{tokString, `a\nb`}}},
		{"single_backslash", `'a\\b'`, []expect{{tokString, `a\\b`}}},
		{"triple_single", `'''abc'''`, []expect{{tokString, "abc"}}},
		{"triple_double", `"""abc"""`, []expect{{tokString, "abc"}}},
		{"triple_with_newline", "'''a\nb'''", []expect{{tokString, "a\nb"}}},
		{"triple_with_quote", `'''a'b'''`, []expect{{tokString, "a'b"}}},
		{"triple_double_with_quote", `"""a"b"""`, []expect{{tokString, `a"b`}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_TripleQuoteGreedy(t *testing.T) {
	// The tokenizer prefers a triple-quote opener over two empty strings.
	// VERIFIED against the Spanner emulator:
	//   "''''"   -> REJECT "Unclosed triple-quoted string literal" (''' opens a triple)
	//   "''''''" -> ACCEPT (''' + ''' = one empty triple-quoted string)
	//   "'' ''"  -> ACCEPT, two strings ("Concatenation ... not supported")
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"six_singles_empty_triple", "''''''", []expect{{tokString, ""}}},
		{"six_doubles_empty_triple", `""""""`, []expect{{tokString, ""}}},
		{"two_empty_strings", "'' ''", []expect{{tokString, ""}, {tokString, ""}}},
		{"triple_embedded_quotes", `'''a''b'''`, []expect{{tokString, "a''b"}}},
	}
	runLexerCases(t, cases)

	// Four single quotes: ''' opens a triple that is never closed -> error.
	_, errs := Tokenize("''''")
	if len(errs) != 1 {
		t.Errorf("\"''''\" expected 1 unterminated-triple error, got %d: %+v", len(errs), errs)
	}
}

func TestLexer_RawStringLiterals(t *testing.T) {
	// R prefix marks a raw string. Backslashes are NOT escapes inside a raw
	// string. The lexer records IsRaw=true and preserves the body verbatim.
	tokens, errs := Tokenize(`r'a\nb'`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 2 || tokens[0].Type != tokString {
		t.Fatalf("expected single tokString + EOF, got %+v", tokens)
	}
	if !tokens[0].IsRaw {
		t.Errorf("expected IsRaw=true for r'...'")
	}
	if tokens[0].Str != `a\nb` {
		t.Errorf("Str = %q, want %q", tokens[0].Str, `a\nb`)
	}
	// Uppercase R prefix.
	tokens, _ = Tokenize(`R"x"`)
	if !tokens[0].IsRaw || tokens[0].Type != tokString || tokens[0].Str != "x" {
		t.Errorf("R-prefix: got %+v", tokens[0])
	}
	// Raw triple-quoted.
	tokens, _ = Tokenize(`r'''a\b'''`)
	if !tokens[0].IsRaw || tokens[0].Str != `a\b` {
		t.Errorf("raw triple: got %+v", tokens[0])
	}
}

func TestLexer_BytesLiterals(t *testing.T) {
	// BYTES_LITERAL: (B | RB | BR) (SQTEXT | DQTEXT | SQ3TEXT | DQ3TEXT).
	cases := []struct {
		name  string
		input string
		raw   bool
		str   string
	}{
		{"bytes_single", `b'abc'`, false, "abc"},
		{"bytes_double", `b"abc"`, false, "abc"},
		{"bytes_upper", `B'abc'`, false, "abc"},
		{"bytes_triple", `b'''abc'''`, false, "abc"},
		{"raw_bytes_rb", `rb'a\x'`, true, `a\x`},
		{"raw_bytes_br", `br'a\x'`, true, `a\x`},
		{"raw_bytes_upper", `RB'a'`, true, "a"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			if tokens[0].Type != tokBytes {
				t.Fatalf("Type = %s, want tokBytes", TokenName(tokens[0].Type))
			}
			if tokens[0].IsRaw != c.raw {
				t.Errorf("IsRaw = %v, want %v", tokens[0].IsRaw, c.raw)
			}
			if tokens[0].Str != c.str {
				t.Errorf("Str = %q, want %q", tokens[0].Str, c.str)
			}
		})
	}
}

func TestLexer_PrefixLetterIsIdentifierWhenNotQuote(t *testing.T) {
	// r/R/b/B and the rb/br pairs only introduce a string/bytes literal when
	// immediately followed by a quote. Otherwise they are ordinary identifiers.
	// VERIFIED against the Spanner emulator: "SELECT rbx'y'" -> "got string
	// literal 'y'" (rbx is an alias/ident, 'y' a separate string).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"bare_r", "r", []expect{{tokIdentifier, "r"}}},
		{"bare_b", "b", []expect{{tokIdentifier, "b"}}},
		{"rate", "rate", []expect{{tokIdentifier, "rate"}}},
		{"b2", "b2", []expect{{tokIdentifier, "b2"}}},
		{"rbx_then_string", "rbx'y'", []expect{{tokIdentifier, "rbx"}, {tokString, "y"}}},
		{"keyword_then_string", "from'a'", []expect{{kwFROM, "from"}, {tokString, "a"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_IntegerLiterals(t *testing.T) {
	// INTEGER_LITERAL: DECIMAL_DIGITS | HEX_DIGITS ('0x' hexdigit+).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"zero", "0", []expect{{tokInteger, "0"}}},
		{"small", "42", []expect{{tokInteger, "42"}}},
		{"large", "1234567890", []expect{{tokInteger, "1234567890"}}},
		{"hex_lower", "0x1f", []expect{{tokInteger, "0x1f"}}},
		{"hex_upper", "0XABCDEF", []expect{{tokInteger, "0XABCDEF"}}},
		{"hex_mixed", "0xAbCd", []expect{{tokInteger, "0xAbCd"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_IntegerValue(t *testing.T) {
	// Ival is populated for decimal integers.
	tokens, _ := Tokenize("42")
	if tokens[0].Ival != 42 {
		t.Errorf("Ival = %d, want 42", tokens[0].Ival)
	}
	// Hex integers populate Ival too.
	tokens, _ = Tokenize("0x1F")
	if tokens[0].Ival != 31 {
		t.Errorf("hex Ival = %d, want 31", tokens[0].Ival)
	}
}

func TestLexer_FloatLiterals(t *testing.T) {
	// FLOATING_POINT_LITERAL: digits.digits?[E...] | digits?.digits[E...] |
	// digits E digits.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple", "1.5", []expect{{tokFloat, "1.5"}}},
		{"trailing_dot", "1.", []expect{{tokFloat, "1."}}},
		{"leading_dot", ".5", []expect{{tokFloat, ".5"}}},
		{"exp_int_base", "1e10", []expect{{tokFloat, "1e10"}}},
		{"exp_upper", "1E10", []expect{{tokFloat, "1E10"}}},
		{"exp_pos", "1e+10", []expect{{tokFloat, "1e+10"}}},
		{"exp_neg", "1e-10", []expect{{tokFloat, "1e-10"}}},
		{"float_with_exp", "1.5e10", []expect{{tokFloat, "1.5e10"}}},
		{"dot_base_exp", ".5e10", []expect{{tokFloat, ".5e10"}}},
		{"trailing_dot_exp", "1.e10", []expect{{tokFloat, "1.e10"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_NumberBoundaries(t *testing.T) {
	// GoogleSQL's FLOATING_POINT_LITERAL is maximal-munch: "1.5" is one float,
	// and "1." (digits '.' digits?) is also a complete float.
	//
	// "1.2.3" tokenizes as the float "1.2" followed by the float ".3"
	// (the digits? '.' digits form). VERIFIED against the Spanner emulator,
	// whose parse error reads: got floating point literal ".3". This matches
	// the legacy ANTLR / ZetaSQL flex tokenizer.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"int_dot_int", "1.5", []expect{{tokFloat, "1.5"}}},
		{"version_like", "1.2.3", []expect{{tokFloat, "1.2"}, {tokFloat, ".3"}}},
		// integer followed by dotted ident: "0x1f" stays hex.
		{"hex_full", "0x1f", []expect{{tokInteger, "0x1f"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_SignNotAbsorbedIntoNumber(t *testing.T) {
	// DIVERGENCE FROM LEGACY ANTLR (oracle-decided): the legacy
	// GoogleSQLLexer.g4 FLOATING_POINT_LITERAL rule carries an optional leading
	// (PLUS|MINUS) sign, which under ANTLR maximal-munch would absorb the '-'
	// in "1.5-2.5" into the second float, breaking subtraction. The
	// authoritative ZetaSQL tokenizer (and the live Spanner emulator) does NOT
	// do this — the sign is a separate operator token handled by the parser's
	// signed_numeric_literal rule.
	//
	// VERIFIED on the Spanner emulator:
	//   "SELECT 1.5 2.5"  -> REJECT (two adjacent floats: "got floating point literal 2.5")
	//   "SELECT 1.5-2.5"  -> ACCEPT (parses as subtraction)
	// => the '-' must lex as a standalone MINUS_OPERATOR.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"minus_float", "1.5-2.5", []expect{{tokFloat, "1.5"}, {int('-'), ""}, {tokFloat, "2.5"}}},
		{"minus_int", "5-5", []expect{{tokInteger, "5"}, {int('-'), ""}, {tokInteger, "5"}}},
		{"leading_minus_float", "-1.5", []expect{{int('-'), ""}, {tokFloat, "1.5"}}},
		{"leading_plus_int", "+5", []expect{{int('+'), ""}, {tokInteger, "5"}}},
		{"plus_float", "1.5+2.5", []expect{{tokFloat, "1.5"}, {int('+'), ""}, {tokFloat, "2.5"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_QueryParameters(t *testing.T) {
	// @ is AT_SYMBOL, @@ is ATAT_SYMBOL, ? is QUESTION_SYMBOL. The lexer emits
	// these as standalone tokens; the parser assembles @name / @@path.
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"named_param", "@p", []expect{{int('@'), ""}, {tokIdentifier, "p"}}},
		{"positional_param", "?", []expect{{int('?'), ""}}},
		{"system_var", "@@x", []expect{{tokAtAt, ""}, {tokIdentifier, "x"}}},
		// Hint open: '@' followed by '{' — two tokens (no combined @{ token).
		{"hint_open", "@{", []expect{{int('@'), ""}, {int('{'), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_Comments(t *testing.T) {
	// All three comment forms produce empty token streams (just EOF).
	// BLOCK_COMMENT does NOT nest in GoogleSQL (unlike Snowflake).
	cases := []string{
		"-- line comment\n",
		"# pound comment\n",
		"/* block comment */",
		"/**/",
		"/* a */ /* b */",
		"-- comment 1\n-- comment 2\n",
		"# c1\n# c2\n",
		"-- no trailing newline",
		"# no trailing newline",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			tokens, errs := Tokenize(input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			if len(tokens) != 1 || tokens[0].Type != tokEOF {
				t.Errorf("expected only EOF token, got %+v", tokens)
			}
		})
	}
}

func TestLexer_BlockCommentDoesNotNest(t *testing.T) {
	// GoogleSQL BLOCK_COMMENT: '/**/' | '/*' ~[!] .*? '*/' — non-greedy, NO
	// nesting. "/* outer /* inner */ tail */" closes at the FIRST */, leaving
	// "tail */" as real tokens.
	tokens, errs := Tokenize("/* outer /* inner */ tail */")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	// After the first */, the remaining "tail */" lexes as: IDENT(tail) * /
	want := []expect{{tokIdentifier, "tail"}, {int('*'), ""}, {int('/'), ""}}
	tokens = tokens[:len(tokens)-1] // strip EOF
	if len(tokens) != len(want) {
		t.Fatalf("token count = %d, want %d: %+v", len(tokens), len(want), tokens)
	}
	for i, tok := range tokens {
		if tok.Type != want[i].Type || tok.Str != want[i].Str {
			t.Errorf("token[%d] = (%s,%q), want (%s,%q)",
				i, TokenName(tok.Type), tok.Str, TokenName(want[i].Type), want[i].Str)
		}
	}
}

func TestLexer_BlockCommentBangNotComment(t *testing.T) {
	// BLOCK_COMMENT excludes '/*!' (the ~[!] guard). "/*!x*/" is therefore NOT
	// a comment: it lexes as / * ! x * / (operators + ident).
	tokens, errs := Tokenize("/*!x*/")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	tokens = tokens[:len(tokens)-1]
	want := []expect{{int('/'), ""}, {int('*'), ""}, {int('!'), ""}, {tokIdentifier, "x"}, {int('*'), ""}, {int('/'), ""}}
	if len(tokens) != len(want) {
		t.Fatalf("token count = %d, want %d: %+v", len(tokens), len(want), tokens)
	}
	for i, tok := range tokens {
		if tok.Type != want[i].Type {
			t.Errorf("token[%d] = %s, want %s", i, TokenName(tok.Type), TokenName(want[i].Type))
		}
	}
}

func TestLexer_CommentsBetweenTokens(t *testing.T) {
	tokens, errs := Tokenize("SELECT /* hidden */ FROM")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens (SELECT FROM EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != kwSELECT || tokens[1].Type != kwFROM {
		t.Errorf("expected [SELECT, FROM, EOF], got %+v", tokens)
	}
}

func TestLexer_Keywords(t *testing.T) {
	// Spot-check keywords from each category; full coverage in
	// keyword_completeness_test.go (checks every keyword from the .g4).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"select_lower", "select", []expect{{kwSELECT, "select"}}},
		{"select_upper", "SELECT", []expect{{kwSELECT, "SELECT"}}},
		{"select_mixed", "SeLeCt", []expect{{kwSELECT, "SeLeCt"}}},
		{"from", "FROM", []expect{{kwFROM, "FROM"}}},
		{"where", "WHERE", []expect{{kwWHERE, "WHERE"}}},
		{"create", "CREATE", []expect{{kwCREATE, "CREATE"}}},
		{"table", "TABLE", []expect{{kwTABLE, "TABLE"}}},
		{"struct", "STRUCT", []expect{{kwSTRUCT, "STRUCT"}}},
		{"array", "ARRAY", []expect{{kwARRAY, "ARRAY"}}},
		{"interleave", "INTERLEAVE", []expect{{kwINTERLEAVE, "INTERLEAVE"}}},
		{"qualify", "QUALIFY", []expect{{kwQUALIFY, "QUALIFY"}}},
		{"unnest", "UNNEST", []expect{{kwUNNEST, "UNNEST"}}},
		{"graph", "GRAPH", []expect{{kwGRAPH, "GRAPH"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_KeywordReservedClassification(t *testing.T) {
	// IsReservedKeyword against the .g4 common_keyword_as_identifier split.
	reserved := []string{"SELECT", "FROM", "WHERE", "CREATE", "JOIN", "WITH", "UNION", "STRUCT", "PARTITION"}
	for _, kw := range reserved {
		if !IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = false, want true", kw)
		}
	}
	// These ARE word-keywords but appear in common_keyword_as_identifier, so
	// they are usable as bare identifiers (non-reserved).
	nonReserved := []string{"TABLE", "VIEW", "INDEX", "OPTIONS", "INTERLEAVE", "SIMPLE", "SYSTEM", "VALUE", "OFFSET", "GRAPH"}
	for _, kw := range nonReserved {
		if IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = true, want false", kw)
		}
	}
	// A plain identifier is never reserved.
	if IsReservedKeyword("my_column") {
		t.Errorf("IsReservedKeyword(my_column) = true, want false")
	}
}

func TestLexer_NonReservedKeywordIsIdentToken(t *testing.T) {
	// Non-reserved keywords still lex AS their keyword token (the
	// reserved/non-reserved distinction is enforced by the parser via
	// keyword_as_identifier, not by the lexer). The lexer always returns the
	// kw* constant for any keyword.
	tokens, _ := Tokenize("TABLE")
	if tokens[0].Type != kwTABLE {
		t.Errorf("non-reserved keyword TABLE lexed as %s, want kwTABLE", TokenName(tokens[0].Type))
	}
}

func TestLexer_FullStatement(t *testing.T) {
	// A realistic GoogleSQL statement covering keywords, idents, string,
	// number, operators, and backtick path.
	tokens, errs := Tokenize("SELECT a, b FROM `proj.ds.t` WHERE c >= 10 AND name = 'x'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	want := []expect{
		{kwSELECT, "SELECT"},
		{tokIdentifier, "a"},
		{int(','), ""},
		{tokIdentifier, "b"},
		{kwFROM, "FROM"},
		{tokIdentifier, "proj.ds.t"},
		{kwWHERE, "WHERE"},
		{tokIdentifier, "c"},
		{tokGreaterEqual, ""},
		{tokInteger, "10"},
		{kwAND, "AND"},
		{tokIdentifier, "name"},
		{int('='), ""},
		{tokString, "x"},
	}
	tokens = tokens[:len(tokens)-1] // strip EOF
	if len(tokens) != len(want) {
		t.Fatalf("token count = %d, want %d:\n%+v", len(tokens), len(want), tokens)
	}
	for i, tok := range tokens {
		if tok.Type != want[i].Type || tok.Str != want[i].Str {
			t.Errorf("token[%d] = (%s,%q), want (%s,%q)",
				i, TokenName(tok.Type), tok.Str, TokenName(want[i].Type), want[i].Str)
		}
	}
}

func TestLexer_EmptyInput(t *testing.T) {
	tokens, errs := Tokenize("")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF, got %+v", tokens)
	}
}

func TestLexer_WhitespaceOnly(t *testing.T) {
	// WHITESPACE: [ \t\f\r\n]. Includes form-feed (\f).
	tokens, errs := Tokenize(" \t\f\r\n ")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF, got %+v", tokens)
	}
}

// ---- Error cases ----

func TestLexer_UnterminatedString(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"single", `'unterminated`},
		{"double", `"unterminated`},
		{"single_newline", "'has\nnewline'"}, // single-quote string cannot span newline
		{"double_newline", "\"has\nnewline\""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) == 0 {
				t.Fatalf("expected a lex error for %q, got none (tokens=%+v)", c.input, tokens)
			}
			// First token (or one of them) must be tokInvalid.
			foundInvalid := false
			for _, tok := range tokens {
				if tok.Type == tokInvalid {
					foundInvalid = true
				}
			}
			if !foundInvalid {
				t.Errorf("expected a tokInvalid token, got %+v", tokens)
			}
		})
	}
}

func TestLexer_UnterminatedTripleString(t *testing.T) {
	_, errs := Tokenize(`'''unterminated triple`)
	if len(errs) == 0 {
		t.Fatalf("expected a lex error for unterminated triple string")
	}
}

func TestLexer_UnterminatedBacktick(t *testing.T) {
	tokens, errs := Tokenize("`unterminated")
	if len(errs) == 0 {
		t.Fatalf("expected a lex error for unterminated backtick identifier, tokens=%+v", tokens)
	}
}

func TestLexer_UnterminatedBlockComment(t *testing.T) {
	_, errs := Tokenize("/* never ends")
	if len(errs) == 0 {
		t.Fatalf("expected a lex error for unterminated block comment")
	}
}

func TestLexer_InvalidByte(t *testing.T) {
	// '$' is not a valid GoogleSQL byte (no $ token, not an identifier char).
	tokens, errs := Tokenize("$")
	if len(errs) == 0 {
		t.Fatalf("expected a lex error for '$', tokens=%+v", tokens)
	}
	if tokens[0].Type != tokInvalid {
		t.Errorf("expected tokInvalid, got %s", TokenName(tokens[0].Type))
	}
}

func TestLexer_DollarInIdentifierIsError(t *testing.T) {
	// "x$y" -> "x" (ident), then '$' is invalid, then "y" (ident).
	tokens, errs := Tokenize("x$y")
	if len(errs) != 1 {
		t.Fatalf("expected exactly 1 lex error, got %d: %+v", len(errs), errs)
	}
	if tokens[0].Type != tokIdentifier || tokens[0].Str != "x" {
		t.Errorf("token[0] = %+v, want ident x", tokens[0])
	}
}

func TestLexer_ErrorRecoveryContinues(t *testing.T) {
	// After a bad byte, lexing should continue and produce subsequent tokens.
	tokens, errs := Tokenize("SELECT $ FROM t")
	if len(errs) != 1 {
		t.Fatalf("expected 1 lex error, got %d: %+v", len(errs), errs)
	}
	// Should still see SELECT, (invalid), FROM, t, EOF.
	if tokens[0].Type != kwSELECT {
		t.Errorf("token[0] = %s, want SELECT", TokenName(tokens[0].Type))
	}
	var sawFrom, sawT bool
	for _, tok := range tokens {
		if tok.Type == kwFROM {
			sawFrom = true
		}
		if tok.Type == tokIdentifier && tok.Str == "t" {
			sawT = true
		}
	}
	if !sawFrom || !sawT {
		t.Errorf("error recovery failed to continue: %+v", tokens)
	}
}
