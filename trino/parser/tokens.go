// Package parser implements a hand-written recursive-descent parser for Trino
// SQL (Trino 481 grammar, derived from the official SqlBase.g4). This file
// defines the token model produced by the lexer; keyword tokens live in
// keywords.go and the scanner in lexer.go.
//
// The package mirrors omni's doris/parser and snowflake/parser conventions: a
// stateless lexer emits Token{Kind, Str, Ival, Loc}; single-byte operators and
// punctuation use their ASCII byte value directly as the Kind; multi-byte
// operators, literals, and keywords use the named constants below. Source
// positions are tracked as ast.Loc{Start, End} byte offsets so the statement
// splitter and downstream parser/completion can map tokens back to the input.
package parser

import "github.com/bytebase/omni/trino/ast"

// Special tokens.
const (
	tokEOF     = 0
	tokInvalid = 1 // error-recovery token; always accompanied by a LexError
)

// TokenKind identifies a token type. Single-character operators and
// punctuation use their ASCII byte value directly as their TokenKind
// (e.g., ';' = 59, '(' = 40, '.' = 46). The named constants below cover
// multi-character operators, literals, and keywords (kw* in keywords.go).
type TokenKind = int

// Multi-character operators and punctuation (500-599).
//
// Names and shapes follow Trino's TrinoLexer.g4 token vocabulary. Single-byte
// punctuation (EQ_ '=', LT_ '<', GT_ '>', PLUS_ '+', MINUS_ '-', ASTERISK_ '*',
// SLASH_ '/', PERCENT_ '%', SEMICOLON_ ';', DOT_ '.', COLON_ ':', COMMA_ ',',
// LPAREN_ '(', RPAREN_ ')', LSQUARE_ '[', RSQUARE_ ']', LCURLY_ '{',
// RCURLY_ '}', VBAR_ '|', DOLLAR_ '$', CARET_ '^', QUESTION_MARK_ '?') is
// emitted with Kind == int(byte) rather than a named constant here.
const (
	tokNotEq        = 500 + iota // <> or != (NEQ_)
	tokLessEq                    // <=        (LTE_)
	tokGreaterEq                 // >=        (GTE_)
	tokConcat                    // ||        (CONCAT_)
	tokArrow                     // ->        (RARROW_)
	tokLeftArrow                 // <-        (LARROW_)
	tokDoubleArrow               // =>        (RDOUBLEARROW_)
	tokLCurlyHyphen              // {-        (LCURLYHYPHEN_, MATCH_RECOGNIZE excluded pattern)
	tokRCurlyHyphen              // -}        (RCURLYHYPHEN_, MATCH_RECOGNIZE excluded pattern)
)

// Literal and identifier tokens (600-699).
//
// The four identifier kinds are kept distinct because Trino's grammar (and the
// bytebase completion consumer) treats them differently: unquoted identifiers
// may collide with non-reserved keywords, quoted/back-quoted identifiers never
// do, and a digit-leading identifier is a separate lexer rule.
const (
	tokInteger         = 600 + iota // INTEGER_VALUE_: DIGIT+
	tokDecimal                      // DECIMAL_VALUE_: DIGIT+ '.' DIGIT* | '.' DIGIT+
	tokDouble                       // DOUBLE_VALUE_:  ... EXPONENT
	tokString                       // STRING_:        '...'  ('' escape)
	tokUnicodeString                // UNICODE_STRING_: U&'...' [UESCAPE '...']
	tokBinaryLiteral                // BINARY_LITERAL_: X'...'
	tokIdent                        // IDENTIFIER_:     (LETTER|'_') (LETTER|DIGIT|'_')*
	tokQuotedIdent                  // QUOTED_IDENTIFIER_:     "..."  ("" escape)
	tokBackquotedIdent              // BACKQUOTED_IDENTIFIER_: `...`  (`` escape)
	tokDigitIdent                   // DIGIT_IDENTIFIER_:      DIGIT (LETTER|DIGIT|'_')+
	tokQuestion                     // QUESTION_MARK_: ?  (positional parameter / placeholder)
)

// Token represents a single lexical token.
//
// For string/identifier kinds, Str holds the decoded content (escapes
// resolved, surrounding quotes removed). For keyword and unquoted-identifier
// kinds, Str holds the original source text (preserving case) so consumers can
// echo or re-normalize it. Ival holds the parsed value for tokInteger only.
type Token struct {
	Kind TokenKind // tok*/kw* constant or ASCII byte value
	Str  string    // decoded content for literals/identifiers; source text for keywords/idents
	Ival int64     // parsed integer value for tokInteger (best-effort; 0 on overflow)
	Loc  ast.Loc   // byte-offset span [Start,End) in the source text
}

// LexError records a lexing error with its source position.
type LexError struct {
	Msg string
	Loc ast.Loc
}

// Error messages emitted by the lexer.
const (
	errUnterminatedString  = "unterminated string literal"
	errUnterminatedQuoted  = "unterminated quoted identifier"
	errUnterminatedBinary  = "unterminated binary literal"
	errUnterminatedComment = "unterminated bracketed comment"
	errUnknownChar         = "unrecognized character"
)
