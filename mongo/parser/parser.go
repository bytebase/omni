package parser

import (
	"fmt"

	"github.com/bytebase/omni/mongo/ast"
)

// Parser is a recursive descent parser for mongosh commands.
type Parser struct {
	lexer   *Lexer
	input   string
	cur     Token // current token
	prev    Token // previous token
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// ParseResult holds the outcome of a best-effort parse.
type ParseResult struct {
	Nodes  []ast.Node    // successfully parsed statement nodes
	Errors []*ParseError // errors encountered during parsing
}

// Parse parses a mongosh input string into a list of AST nodes.
// Each semicolon-separated or newline-separated command becomes one node.
//
// On syntax error the parser recovers by skipping to the next statement
// boundary and continues parsing. Returns a non-nil error only when no
// statements could be parsed at all; otherwise partial results are returned
// and errors are available via ParseBestEffort.
func Parse(input string) ([]ast.Node, error) {
	result := ParseBestEffort(input)
	if len(result.Nodes) == 0 && len(result.Errors) > 0 {
		return nil, result.Errors[0]
	}
	return result.Nodes, nil
}

// ParseBestEffort parses a mongosh input string, recovering from errors.
// It always returns as many successfully parsed statements as possible,
// along with any errors encountered.
func ParseBestEffort(input string) *ParseResult {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	result := &ParseResult{}
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		node, err := p.parseStatement()
		if err != nil {
			var parseErr *ParseError
			if e, ok := err.(*ParseError); ok {
				parseErr = e
			} else {
				parseErr = &ParseError{Message: err.Error()}
			}
			result.Errors = append(result.Errors, parseErr)
			p.skipToNextStatement()
			continue
		}
		if node != nil {
			result.Nodes = append(result.Nodes, node)
		}
	}
	return result
}

// skipToNextStatement advances past the current erroneous statement to the
// start of the next one. It always consumes at least one token to guarantee
// progress. It skips tokens until it finds a semicolon or a newline at
// nesting depth 0, or reaches EOF.
func (p *Parser) skipToNextStatement() {
	// Always consume at least one token to avoid infinite loops.
	p.advance()

	depth := 0
	for p.cur.Type != tokEOF {
		switch p.cur.Type {
		case '{', '[', '(':
			depth++
		case '}', ']', ')':
			if depth > 0 {
				depth--
			}
		case ';':
			p.advance()
			return
		}
		// Check if there's a newline between the previous and current token
		// at depth 0 — the current token starts a new statement.
		if depth == 0 && p.prev.End > 0 && p.cur.Loc > p.prev.End {
			for i := p.prev.End; i < p.cur.Loc && i < len(p.input); i++ {
				if p.input[i] == '\n' {
					return
				}
			}
		}
		p.advance()
	}
}

// parseStatement dispatches to the appropriate family parser based on the first token.
func (p *Parser) parseStatement() (ast.Node, error) {
	switch p.cur.Type {
	case kwShow:
		return p.parseShowCommand()
	case kwDb:
		return p.parseDbStatement()
	case kwRs:
		return p.parseRsStatement()
	case kwSh:
		return p.parseShStatement()
	case kwSp:
		return p.parseSpStatement()
	case kwMongo:
		return p.parseConnectionStatement()
	case kwConnect:
		return p.parseConnectionStatement()
	default:
		// Check for native functions: identifier-like token followed by '('
		if p.isIdentLike(p.cur.Type) && p.peekNext().Type == '(' {
			return p.parseNativeFunctionCall()
		}
		return nil, p.syntaxErrorAtCur()
	}
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

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.unexpectedTokenError(tokenType)
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

// isIdentLike returns true if the given token type is an identifier or a keyword
// that can be used as an identifier in certain contexts (e.g., method names).
func (p *Parser) isIdentLike(tokenType int) bool {
	return tokenType == tokIdent || tokenType >= 700
}

// expectIdent consumes the current token if it's an identifier or keyword-as-identifier.
// Returns the token string and an error if not identifier-like.
func (p *Parser) expectIdent() (string, Token, error) {
	if p.isIdentLike(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok, nil
	}
	return "", Token{}, p.unexpectedTokenError(tokIdent)
}

// parseArguments parses a parenthesized argument list: ( expr, expr, ... )
// A single trailing comma before the closing paren is allowed, matching the
// behavior of mongosh / JavaScript (ES2017+).
// Assumes the opening '(' has NOT been consumed yet.
func (p *Parser) parseArguments() ([]ast.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var args []ast.Node
	if p.cur.Type == ')' {
		p.advance()
		return args, nil
	}

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if p.cur.Type == ')' {
			p.advance()
			return args, nil
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		// Allow a single trailing comma before the closing paren:
		//   db.coll.find({…},).sort(…)  → equivalent to db.coll.find({…}).sort(…)
		if p.cur.Type == ')' {
			p.advance()
			return args, nil
		}
	}
}

// parseMethodCall parses .methodName(args...) assuming the dot has been consumed.
func (p *Parser) parseMethodCall() (*ast.MethodCall, error) {
	name, tok, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	start := tok.Loc

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &ast.MethodCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: p.prev.End},
	}, nil
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
	Line     int // 1-based line number
	Column   int // 1-based column number
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s (line %d, column %d)", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// syntaxErrorAtCur returns a ParseError at the current token position.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	line, col := p.lineCol(p.cur.Loc)
	var msg string
	if p.cur.Type == tokEOF {
		msg = "syntax error at end of input"
	} else {
		text := p.cur.Str
		if text == "" {
			text = string(rune(p.cur.Type))
		}
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Line:     line,
		Column:   col,
	}
}

// unexpectedTokenError returns a ParseError for an unexpected token.
func (p *Parser) unexpectedTokenError(expected int) *ParseError {
	line, col := p.lineCol(p.cur.Loc)
	var msg string
	if p.cur.Type == tokEOF {
		msg = fmt.Sprintf("expected %s but reached end of input", tokenName(expected))
	} else {
		text := p.cur.Str
		if text == "" {
			text = string(rune(p.cur.Type))
		}
		msg = fmt.Sprintf("expected %s, got %q", tokenName(expected), text)
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Line:     line,
		Column:   col,
	}
}

// lineCol computes the 1-based line and column for a byte offset.
func (p *Parser) lineCol(offset int) (int, int) {
	if offset > len(p.input) {
		offset = len(p.input)
	}
	line := 1
	lastNewline := -1
	for i := 0; i < offset; i++ {
		if p.input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col := offset - lastNewline
	return line, col
}
