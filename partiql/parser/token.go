// Package parser implements the PartiQL parser. This file declares the
// token type used by the lexer (lexer.go) and the full set of token
// type constants. Token positions use ast.Loc directly to eliminate the
// need for the parser to convert at every AST node construction site.
package parser

import "github.com/bytebase/omni/partiql/ast"

// Token is a single PartiQL lexer token.
//
// Type is one of the tok* constants below.
//
// Str holds the raw source text for most tokens. For tokSCONST (single-quoted
// string literals) and tokIDENT_QUOTED (double-quoted identifiers), Str holds
// the *decoded* value with the doubled-quote escape collapsed: two consecutive
// quote characters inside the literal represent a single quote in the decoded
// value (e.g., the SCONST source spelled i, t, quote, quote, s with surrounding
// single quotes decodes to the Go string "it's"). For tokION_LITERAL, Str is
// the verbatim inner content between the backticks (no decoding).
//
// Loc is the half-open byte range covering the token in the source string.
type Token struct {
	Type int
	Str  string
	Loc  ast.Loc
}

// ===========================================================================
// Special tokens.
// ===========================================================================

const (
	tokEOF     = 0 // end of input or after lex error
	tokInvalid = 1 // sentinel for unknown token type (never returned by Next)
)

// ===========================================================================
// Literal tokens — group 1000.
// ===========================================================================

const (
	tokSCONST       = iota + 1000 // single-quoted string literal: 'hello'
	tokICONST                     // integer literal: 42
	tokFCONST                     // decimal/float literal: 3.14, 1e10, .5
	tokIDENT                      // unquoted identifier (case-insensitive lookup)
	tokIDENT_QUOTED               // double-quoted identifier (case-sensitive): "Foo"
	tokION_LITERAL                // backtick-delimited Ion blob (body deferred to DAG node 17)
)

// ===========================================================================
// Operator and punctuation tokens — group 2000.
//
// Names follow PartiQLLexer.g4 rule names verbatim for traceability against
// the grammar.
// ===========================================================================

const (
	tokPLUS               = iota + 2000 // +
	tokMINUS                            // -
	tokASTERISK                         // *
	tokSLASH_FORWARD                    // /
	tokPERCENT                          // %
	tokCARET                            // ^
	tokTILDE                            // ~
	tokAT_SIGN                          // @
	tokEQ                               // =
	tokNEQ                              // <> or !=
	tokLT                               // < (ANGLE_LEFT in grammar)
	tokGT                               // > (ANGLE_RIGHT in grammar)
	tokLT_EQ                            // <=
	tokGT_EQ                            // >=
	tokCONCAT                           // ||
	tokANGLE_DOUBLE_LEFT                // <<  (PartiQL bag-literal start)
	tokANGLE_DOUBLE_RIGHT               // >>  (PartiQL bag-literal end)
	tokPAREN_LEFT                       // (
	tokPAREN_RIGHT                      // )
	tokBRACKET_LEFT                     // [
	tokBRACKET_RIGHT                    // ]
	tokBRACE_LEFT                       // {
	tokBRACE_RIGHT                      // }
	tokCOLON                            // :
	tokCOLON_SEMI                       // ;
	tokCOMMA                            // ,
	tokPERIOD                           // .
	tokQUESTION_MARK                    // ?
)

// tokenName returns the canonical printable name for a token type constant.
// Used by error messages, test failure output, and future debugging.
//
// Task 3 expands this switch to cover all 302 constants. For now it
// covers only the 36 non-keyword constants from this file.
func tokenName(t int) string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"
	case tokSCONST:
		return "SCONST"
	case tokICONST:
		return "ICONST"
	case tokFCONST:
		return "FCONST"
	case tokIDENT:
		return "IDENT"
	case tokIDENT_QUOTED:
		return "IDENT_QUOTED"
	case tokION_LITERAL:
		return "ION_LITERAL"
	case tokPLUS:
		return "PLUS"
	case tokMINUS:
		return "MINUS"
	case tokASTERISK:
		return "ASTERISK"
	case tokSLASH_FORWARD:
		return "SLASH_FORWARD"
	case tokPERCENT:
		return "PERCENT"
	case tokCARET:
		return "CARET"
	case tokTILDE:
		return "TILDE"
	case tokAT_SIGN:
		return "AT_SIGN"
	case tokEQ:
		return "EQ"
	case tokNEQ:
		return "NEQ"
	case tokLT:
		return "LT"
	case tokGT:
		return "GT"
	case tokLT_EQ:
		return "LT_EQ"
	case tokGT_EQ:
		return "GT_EQ"
	case tokCONCAT:
		return "CONCAT"
	case tokANGLE_DOUBLE_LEFT:
		return "ANGLE_DOUBLE_LEFT"
	case tokANGLE_DOUBLE_RIGHT:
		return "ANGLE_DOUBLE_RIGHT"
	case tokPAREN_LEFT:
		return "PAREN_LEFT"
	case tokPAREN_RIGHT:
		return "PAREN_RIGHT"
	case tokBRACKET_LEFT:
		return "BRACKET_LEFT"
	case tokBRACKET_RIGHT:
		return "BRACKET_RIGHT"
	case tokBRACE_LEFT:
		return "BRACE_LEFT"
	case tokBRACE_RIGHT:
		return "BRACE_RIGHT"
	case tokCOLON:
		return "COLON"
	case tokCOLON_SEMI:
		return "COLON_SEMI"
	case tokCOMMA:
		return "COMMA"
	case tokPERIOD:
		return "PERIOD"
	case tokQUESTION_MARK:
		return "QUESTION_MARK"
	}
	return ""
}
