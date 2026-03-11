// Package parser implements a recursive descent SQL parser for Oracle PL/SQL.
//
// This parser produces AST nodes from the oracle/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Parser is a recursive descent parser for Oracle SQL/PL/SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a SQL string into an AST list.
// Each statement is wrapped in a *RawStmt.
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

		stmtLoc := p.pos()
		stmt := p.parseStmt()
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, &ParseError{
					Message:  fmt.Sprintf("unexpected token %q at position %d", p.cur.Str, p.cur.Loc),
					Position: p.cur.Loc,
				}
			}
			break
		}

		raw := &nodes.RawStmt{
			Stmt:         stmt,
			StmtLocation: stmtLoc,
			StmtLen:      p.pos() - stmtLoc,
		}
		stmts = append(stmts, raw)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// parseStmt dispatches to statement-specific parsers.
// Minimal implementation for batch 0 (infrastructure).
// Full dispatch will be wired in batch 23.
func (p *Parser) parseStmt() nodes.StmtNode {
	switch p.cur.Type {
	case kwSELECT:
		return nil // placeholder: will be implemented in batch 4
	default:
		return nil
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
	return Token{}, &ParseError{
		Message:  fmt.Sprintf("expected token type %d, got %d (%q)", tokenType, p.cur.Type, p.cur.Str),
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// isKeyword checks whether the current token is a specific keyword.
func (p *Parser) isKeyword(kw int) bool {
	return p.cur.Type == kw
}

// matchKeyword consumes the token if it is the given keyword.
func (p *Parser) matchKeyword(kw int) (Token, bool) {
	if p.cur.Type == kw {
		return p.advance(), true
	}
	return Token{}, false
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}
