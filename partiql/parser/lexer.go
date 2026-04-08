package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// Lexer is a hand-written tokenizer for PartiQL source code.
//
// Single-pass scanner. The caller drives it via Next(); each call
// returns one token. At end of input or after a lex error, Next()
// returns Token{Type: tokEOF, ...}. The first error encountered is
// stored in Err and all subsequent Next() calls return tokEOF.
type Lexer struct {
	input string // source text
	pos   int    // current read position (next byte to consume)
	start int    // byte offset of token currently being scanned
	Err   error  // first error encountered, nil if none
}

// NewLexer creates a Lexer for the given source string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token from the input.
// At end of input or after a lex error, returns Token{Type: tokEOF, ...}.
// After Err is set, all subsequent calls return tokEOF.
func (l *Lexer) Next() Token {
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.skipWhitespaceAndComments()
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString()
	case ch == '"':
		return l.scanQuotedIdent()
	case ch == '`':
		return l.scanIonLiteral()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]):
		return l.scanNumber() // leading-dot decimal: .5
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	default:
		return l.scanOperator()
	}
}

// skipWhitespaceAndComments advances l.pos past whitespace, line comments,
// and block comments. All three are on the HIDDEN channel per the grammar
// and never appear in the token stream.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace.
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		// Line comment: -- to end of line.
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		// Block comment: /* ... */ (non-nested, greedy-shortest).
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			start := l.pos
			l.pos += 2
			closed := false
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					closed = true
					break
				}
				l.pos++
			}
			if !closed {
				l.Err = fmt.Errorf("unterminated block comment at position %d", start)
				return
			}
			continue
		}

		break
	}
}

// ============================================================================
// Character class helpers.
// ============================================================================

// isIdentStart reports whether ch can begin a PartiQL identifier.
// PartiQL identifiers start with [a-zA-Z_$].
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' || ch == '$'
}

// isIdentContinue reports whether ch can appear in a PartiQL identifier
// after the first character. Adds digits to isIdentStart.
func isIdentContinue(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}

// isDigit reports whether ch is an ASCII decimal digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// ============================================================================
// STUBS — replaced by Tasks 6–10.
//
// These return tokEOF and set l.Err so the package builds at the end of
// Task 5. Each subsequent task removes one stub and adds the real
// implementation alongside its tests.
// ============================================================================

func (l *Lexer) scanString() Token {
	l.Err = fmt.Errorf("scanString not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanQuotedIdent() Token {
	l.Err = fmt.Errorf("scanQuotedIdent not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanIdentOrKeyword() Token {
	l.Err = fmt.Errorf("scanIdentOrKeyword not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanNumber() Token {
	l.Err = fmt.Errorf("scanNumber not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanOperator() Token {
	l.Err = fmt.Errorf("scanOperator not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanIonLiteral() Token {
	l.Err = fmt.Errorf("scanIonLiteral not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// strings is imported but only used by the stubs above (and by future
// scan helpers). Avoid the unused-import error during early tasks by
// keeping a no-op reference. Remove this line in Task 7 when
// scanIdentOrKeyword adds the real strings.ToLower call.
var _ = strings.ToLower
