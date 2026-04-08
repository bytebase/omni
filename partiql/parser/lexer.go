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
// Scan helpers — strings, quoted identifiers, and keywords (Tasks 6-7).
// ============================================================================

// scanString scans a single-quoted PartiQL string literal and returns a
// tokSCONST token. The only escape mechanism in PartiQL strings is the
// doubled-apostrophe form: two consecutive apostrophes in a row stand for
// a single apostrophe in the decoded value. There are no backslash escapes.
// Token.Str holds the fully decoded string content (without the surrounding
// quotes). Next() has already set l.start to the offset of the opening quote
// before calling this function, so scanString must not modify l.start.
func (l *Lexer) scanString() Token {
	l.pos++ // skip the opening apostrophe
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			// Two apostrophes in a row encode a single apostrophe.
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			// Single closing apostrophe — end of the literal.
			l.pos++
			return Token{Type: tokSCONST, Str: buf.String(), Loc: ast.Loc{Start: l.start, End: l.pos}}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	// Reached end of input without finding the closing apostrophe.
	l.Err = fmt.Errorf("unterminated string literal at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanQuotedIdent scans a double-quoted PartiQL identifier and returns a
// tokIDENT_QUOTED token. Quoted identifiers preserve case and are NOT looked
// up in the keyword map — they are always identifiers, never keywords. The
// only escape mechanism is the doubled double-quote form: "" within the
// identifier stands for a single " in the decoded name. Token.Str holds the
// decoded identifier text (without the surrounding quotes). Next() has
// already set l.start to the offset of the opening quote before calling this
// function, so scanQuotedIdent must not modify l.start.
func (l *Lexer) scanQuotedIdent() Token {
	l.pos++ // skip the opening double-quote
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			// Two double-quotes in a row encode a single double-quote.
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			// Single closing double-quote — end of the identifier.
			l.pos++
			return Token{Type: tokIDENT_QUOTED, Str: buf.String(), Loc: ast.Loc{Start: l.start, End: l.pos}}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	// Reached end of input without finding the closing double-quote.
	l.Err = fmt.Errorf("unterminated quoted identifier at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanIdentOrKeyword consumes an unquoted identifier and looks it up in
// the keywords map. If the lowercased text matches a keyword, returns
// that keyword token; otherwise returns tokIDENT.
//
// Token.Str preserves the original case (so the parser/AST can render
// identifiers as written). Keyword matching is case-insensitive per
// PartiQLLexer.g4 caseInsensitive=true.
//
// Grammar: IDENTIFIER : [A-Z$_][A-Z0-9$_]*;
//
//	(with caseInsensitive=true expanding to [a-zA-Z$_][a-zA-Z0-9$_]*)
func (l *Lexer) scanIdentOrKeyword() Token {
	for l.pos < len(l.input) && isIdentContinue(l.input[l.pos]) {
		l.pos++
	}
	raw := l.input[l.start:l.pos]
	lower := strings.ToLower(raw)
	if tt, ok := keywords[lower]; ok {
		return Token{
			Type: tt,
			Str:  raw,
			Loc:  ast.Loc{Start: l.start, End: l.pos},
		}
	}
	return Token{
		Type: tokIDENT,
		Str:  raw,
		Loc:  ast.Loc{Start: l.start, End: l.pos},
	}
}

// ============================================================================
// Remaining stubs — replaced by Tasks 8-10.
//
// These return tokEOF and set l.Err so the package builds. Each subsequent
// task removes one stub and adds the real implementation alongside its
// tests. Tasks 6-7 already replaced scanString, scanQuotedIdent, and
// scanIdentOrKeyword.
// ============================================================================

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
