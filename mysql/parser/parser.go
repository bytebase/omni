package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// Parser is a recursive descent parser for MySQL SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid

	// Completion mode fields.
	completing bool          // true when collecting completion candidates
	cursorOff  int           // byte offset of the cursor in source
	candidates *CandidateSet // collected candidates
	collecting bool          // true once cursor position is reached
	maxCollect int           // max exploration depth
}

// Parse parses a SQL string into an AST list.
// Returns a list of statement nodes.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance()

	var stmts []nodes.Node
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, p.enrichError(err)
		}
		if stmt == nil {
			// parseStmt returned nil without error — unconsumed tokens remain.
			if p.cur.Type != tokEOF {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		stmts = append(stmts, stmt)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	if p.cur.Type == tokEOF {
		return Token{}, p.syntaxErrorAtCur()
	}
	return Token{}, &ParseError{
		Message:  "unexpected token",
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// inputText returns a substring of the original input between start and end byte positions.
func (p *Parser) inputText(start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(p.lexer.input) {
		end = len(p.lexer.input)
	}
	if start >= end {
		return ""
	}
	return p.lexer.input[start:end]
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message     string
	Position    int
	Line        int    // 1-based line number
	Column      int    // 1-based column number
	RelatedText string // source text of the line where the error occurred
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		if e.RelatedText != "" {
			return fmt.Sprintf("%s (line %d, column %d)\nrelated text: %s", e.Message, e.Line, e.Column, e.RelatedText)
		}
		return fmt.Sprintf("%s (line %d, column %d)", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// syntaxErrorAtCur returns a ParseError describing a syntax error at the current token.
// If the current token is EOF, the message says "at end of input".
// Otherwise it says "at or near" with the token text.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	return p.syntaxErrorAtTok(p.cur)
}

// syntaxErrorAtTok returns a ParseError describing a syntax error at the given token.
func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	line, col := p.lineCol(tok.Loc)
	var msg string
	if tok.Type == tokEOF {
		msg = "syntax error at end of input"
	} else {
		text := tok.Str
		if text == "" {
			text = string(rune(tok.Type))
		}
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
		Line:     line,
		Column:   col,
	}
}

// lineCol computes the 1-based line and column for a byte offset in the input.
func (p *Parser) lineCol(offset int) (int, int) {
	input := p.lexer.input
	if offset > len(input) {
		offset = len(input)
	}
	line := 1
	lastNewline := -1
	for i := 0; i < offset; i++ {
		if input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col := offset - lastNewline
	return line, col
}

// lineColStatic computes line and column from raw input (for use outside a Parser).
func lineColStatic(input string, offset int) (int, int) {
	if offset > len(input) {
		offset = len(input)
	}
	line := 1
	lastNewline := -1
	for i := 0; i < offset; i++ {
		if input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col := offset - lastNewline
	return line, col
}

// enrichError fills in Line/Column and RelatedText on a ParseError using the parser's input.
func (p *Parser) enrichError(err error) error {
	if pe, ok := err.(*ParseError); ok {
		if pe.Line == 0 {
			pe.Line, pe.Column = p.lineCol(pe.Position)
		}
		if pe.RelatedText == "" {
			pe.RelatedText = p.extractLine(pe.Position)
		}
	}
	return err
}

// extractLine returns the source line containing the given byte offset.
func (p *Parser) extractLine(offset int) string {
	input := p.lexer.input
	if offset < 0 || offset > len(input) {
		return ""
	}
	// Find start of line
	start := offset
	for start > 0 && input[start-1] != '\n' {
		start--
	}
	// Find end of line
	end := offset
	for end < len(input) && input[end] != '\n' {
		end++
	}
	line := input[start:end]
	// Trim to reasonable length
	if len(line) > 200 {
		line = line[:200] + "..."
	}
	return line
}

