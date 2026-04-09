// Package parser provides a hand-written recursive-descent parser for
// PartiQL. This file holds the Parser struct, token buffer helpers,
// error construction, and the shared utility parsers
// (parseSymbolPrimitive, parseVarRef, parseType).
//
// Expression parsing is split across expr.go (precedence ladder) and
// exprprimary.go (primary-expression dispatch). Literal parsing lives
// in literals.go; path-step chains in path.go.
//
// Upstream: this parser consumes tokens produced by Lexer (lexer.go,
// token.go, keywords.go). The lexer uses first-error-and-stop via
// Lexer.Err; the parser surfaces lexer errors as ParseError.
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// ParseError represents a syntax error during parsing.
//
// Loc is the full byte range of the current (offending) token when the
// error was raised. Using the full token span rather than a zero-width
// point lets future error-rendering code highlight the problematic
// region rather than a bare position marker. The parser uses fail-fast
// semantics — the first ParseError aborts the parse and is returned to
// the caller.
type ParseError struct {
	Message string
	Loc     ast.Loc
}

// Error renders a human-readable message including the byte position.
func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Loc.Start, e.Message)
}

// Parser is the recursive-descent parser for PartiQL.
//
// Fields:
//   - lexer:   the token source
//   - cur:     current lookahead token (the "peek" slot)
//   - prev:    most recently consumed token, used to compute end Locs
//   - nextBuf: one-token lookahead buffer for peekNext()
//   - hasNext: true when nextBuf holds a valid token
type Parser struct {
	lexer   *Lexer
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

// NewParser creates a Parser that reads from the given input string.
// The first token is primed into cur before the constructor returns.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime cur with the first token
	return p
}

// advance moves cur to prev and reads the next token from the lexer.
// When the lexer's first-error-and-stop contract fires, subsequent
// calls produce tokEOF with a nil Str — callers must invoke
// checkLexerErr() at strategic points (function entry, after expect)
// to surface the lexer error. This function does NOT embed the error
// message in cur.Str.
func (p *Parser) advance() {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.Next()
	}
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// peekNext returns the token AFTER cur without consuming either cur or
// the returned token. Uses the nextBuf slot; subsequent calls are
// idempotent until advance() is called.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.Next()
		p.hasNext = true
	}
	return p.nextBuf
}

// match consumes cur if its Type matches any of the given types.
// Returns true on match, false otherwise. Idiomatic for optional
// keywords like DISTINCT/ALL or NOT.
func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.cur.Type == t {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes cur if its Type matches tokenType, otherwise returns
// a *ParseError. Used for required tokens (PAREN_RIGHT, COMMA, etc.).
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, &ParseError{
		Message: fmt.Sprintf("expected %s, got %q", tokenName(tokenType), p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// checkLexerErr returns a *ParseError wrapping the lexer's error if
// any. Parser methods call this at strategic points (function entry,
// after expect) to propagate lexer errors into the parse result.
func (p *Parser) checkLexerErr() error {
	if p.lexer.Err != nil {
		return &ParseError{
			Message: p.lexer.Err.Error(),
			Loc:     p.cur.Loc,
		}
	}
	return nil
}

// deferredFeature constructs a stub error pointing at the current
// token. Every non-foundation case in parsePrimaryBase (and a few
// other dispatch points) calls this helper so the error format is
// uniform across the whole parser.
//
// The error message format is the grep contract for future DAG node
// implementers: running
//
//	grep -rn "deferred to parser-builtins" partiql/parser/
//
// returns the full list of stub call sites node 15 needs to replace.
func (p *Parser) deferredFeature(feature, ownerNode string) error {
	return &ParseError{
		Message: fmt.Sprintf("%s is deferred to %s", feature, ownerNode),
		Loc:     p.cur.Loc,
	}
}

// parseSymbolPrimitive consumes an unquoted or double-quoted identifier
// and returns the name, case-sensitivity flag, and location. Rejects
// bare keywords. Matches grammar rule symbolPrimitive (line 742).
func (p *Parser) parseSymbolPrimitive() (name string, caseSensitive bool, loc ast.Loc, err error) {
	switch p.cur.Type {
	case tokIDENT:
		name = p.cur.Str
		caseSensitive = false
		loc = p.cur.Loc
		p.advance()
		return
	case tokIDENT_QUOTED:
		name = p.cur.Str
		caseSensitive = true
		loc = p.cur.Loc
		p.advance()
		return
	default:
		err = &ParseError{
			Message: fmt.Sprintf("expected identifier, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
		return
	}
}

// parseVarRef handles optional @-prefix plus symbolPrimitive. Matches
// grammar rule varRefExpr (lines 635-636).
//
// NOTE: Task 10 upgrades this function to detect IDENT followed by
// PAREN_LEFT (function call form) and return a deferred-feature stub
// error. The Task 1 version handles only bare varref.
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	atPrefixed := false
	if p.cur.Type == tokAT_SIGN {
		atPrefixed = true
		p.advance()
	}
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	return &ast.VarRef{
		Name:          name,
		AtPrefixed:    atPrefixed,
		CaseSensitive: caseSensitive,
		Loc:           ast.Loc{Start: start, End: nameLoc.End},
	}, nil
}
