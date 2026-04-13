package parser

import "github.com/bytebase/omni/doris/ast"

// Special tokens.
const (
	tokEOF     = 0
	tokInvalid = 1 // error-recovery token; always accompanied by a LexError
)

// TokenKind identifies a token type. Single-character operators and
// punctuation use their ASCII byte value directly as their TokenKind
// (e.g., ';' = 59, '(' = 40). Named constants below cover multi-char
// operators, literals, and keywords.
type TokenKind = int

// Multi-character operators (500-599).
const (
	tokLessEq     = 500 + iota // <=
	tokGreaterEq               // >=
	tokNotEq                   // <> or !=
	tokNullSafeEq              // <=>
	tokLogicalAnd              // &&
	tokDoublePipes             // ||
	tokShiftLeft               // <<
	tokShiftRight              // >>
	tokArrow                   // ->
	tokDoubleAt                // @@
	tokAssign                  // :=
	tokDotDotDot               // ...
	tokHintStart               // /*+
	tokHintEnd                 // */ — emitted for ANY */ in the token stream; the parser is responsible for validating hint context
)

// Literal tokens (600-699).
const (
	tokInt         = 600 + iota // integer: 42, 0xFF, 0b101
	tokFloat                    // decimal/exponent: 3.14, 1e10, 1.5e-10
	tokString                   // single or double quoted string
	tokIdent                    // unquoted identifier
	tokQuotedIdent              // backtick-quoted identifier
	tokHexLiteral               // X'...' or x'...'
	tokBitLiteral               // B'...' or b'...'
	tokPlaceholder              // ?
)

// Token represents a single lexical token.
type Token struct {
	Kind TokenKind // tok*/kw* constant or ASCII byte value
	Str  string    // content for identifiers, strings, hints
	Ival int64     // parsed integer value for tokInt
	Loc  ast.Loc   // byte offset span in source
}

// LexError records a lexing error with its source position.
type LexError struct {
	Msg string
	Loc ast.Loc
}

// Error messages.
const (
	errUnterminatedString  = "unterminated string literal"
	errUnterminatedQuoted  = "unterminated backtick-quoted identifier"
	errUnterminatedComment = "unterminated block comment"
	errUnknownChar         = "unknown character"
)
