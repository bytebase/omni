package parser

import (
	"fmt"
	"strconv"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// ParseError represents a syntax error during parsing.
type ParseError struct {
	Message string
	Pos     int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Pos, e.Message)
}

// Parser is the recursive-descent parser for Cosmos DB SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

// Parse parses a single Cosmos DB SQL SELECT statement and returns a List of RawStmt.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance()
	if p.lexer.Err != nil {
		return nil, &ParseError{Message: p.lexer.Err.Error(), Pos: 0}
	}

	stmt, err := p.parseSelect()
	if err != nil {
		return nil, err
	}

	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q after statement", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}

	raw := &nodes.RawStmt{
		Stmt:         stmt,
		StmtLocation: stmt.Loc.Start,
		StmtLen:      stmt.Loc.End - stmt.Loc.Start,
	}

	return &nodes.List{Items: []nodes.Node{raw}}, nil
}

// advance moves to the next token.
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

// peekNext returns the token after the current one without consuming either.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.Next()
		p.hasNext = true
	}
	return p.nextBuf
}

// match checks if the current token matches any of the given types and consumes it if so.
func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.cur.Type == t {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes the current token if it matches the given type, or returns an error.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, &ParseError{
		Message: fmt.Sprintf("expected %s, got %q", tokenName(tokenType), p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// pos returns the byte offset of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// parseIdentifier accepts tokIDENT plus a set of keywords that can be used as identifiers.
// Returns (name, location, error).
func (p *Parser) parseIdentifier() (string, int, error) {
	switch p.cur.Type {
	case tokIDENT,
		tokIN, tokBETWEEN, tokTOP, tokVALUE, tokORDER, tokBY,
		tokGROUP, tokOFFSET, tokLIMIT, tokASC, tokDESC, tokEXISTS,
		tokLIKE, tokHAVING, tokJOIN, tokESCAPE, tokARRAY, tokROOT, tokRANK:
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return name, loc, nil
	}
	return "", p.cur.Loc, &ParseError{
		Message: fmt.Sprintf("expected identifier, got %q", p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// parsePropertyName accepts everything parseIdentifier does plus more keywords.
func (p *Parser) parsePropertyName() (string, int, error) {
	switch p.cur.Type {
	case tokIDENT,
		tokIN, tokBETWEEN, tokTOP, tokVALUE, tokORDER, tokBY,
		tokGROUP, tokOFFSET, tokLIMIT, tokASC, tokDESC, tokEXISTS,
		tokLIKE, tokHAVING, tokJOIN, tokESCAPE, tokARRAY, tokROOT, tokRANK,
		// Additional keywords allowed as property names:
		tokSELECT, tokFROM, tokWHERE, tokNOT, tokAND, tokOR,
		tokAS, tokTRUE, tokFALSE, tokNULL, tokUNDEFINED, tokUDF, tokDISTINCT:
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return name, loc, nil
	}
	return "", p.cur.Loc, &ParseError{
		Message: fmt.Sprintf("expected property name, got %q", p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// parseIntLiteral parses an integer literal and returns its value and location.
func (p *Parser) parseIntLiteral() (int, int, error) {
	if p.cur.Type != tokICONST {
		return 0, p.cur.Loc, &ParseError{
			Message: fmt.Sprintf("expected integer, got %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
	val, err := strconv.Atoi(p.cur.Str)
	if err != nil {
		return 0, p.cur.Loc, &ParseError{
			Message: fmt.Sprintf("invalid integer %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
	loc := p.cur.Loc
	p.advance()
	return val, loc, nil
}

// parseSelect parses a complete SELECT statement.
func (p *Parser) parseSelect() (*nodes.SelectStmt, error) {
	startLoc := p.pos()
	if _, err := p.expect(tokSELECT); err != nil {
		return nil, err
	}

	stmt := &nodes.SelectStmt{}

	// TOP n
	if p.cur.Type == tokTOP {
		top, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
	}

	// DISTINCT
	if p.cur.Type == tokDISTINCT {
		stmt.Distinct = true
		p.advance()
	}

	// VALUE
	if p.cur.Type == tokVALUE {
		stmt.Value = true
		p.advance()
	}

	// * or target list
	if p.cur.Type == tokSTAR {
		stmt.Star = true
		p.advance()
	} else {
		targets, err := p.parseTargetList()
		if err != nil {
			return nil, err
		}
		stmt.Targets = targets
	}

	// FROM
	if p.cur.Type == tokFROM {
		from, joins, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
		stmt.Joins = joins
	}

	// WHERE
	if p.cur.Type == tokWHERE {
		p.advance()
		where, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY
	if p.cur.Type == tokGROUP {
		groupBy, err := p.parseGroupBy()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
	}

	// HAVING
	if p.cur.Type == tokHAVING {
		p.advance()
		having, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// ORDER BY
	if p.cur.Type == tokORDER {
		orderBy, err := p.parseOrderBy()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// OFFSET ... LIMIT ...
	if p.cur.Type == tokOFFSET {
		ol, err := p.parseOffsetLimit()
		if err != nil {
			return nil, err
		}
		stmt.OffsetLimit = ol
	}

	endLoc := p.prev.Loc + len(p.prev.Str)
	stmt.Loc = nodes.Loc{Start: startLoc, End: endLoc}
	return stmt, nil
}

// tokenName returns a human-readable name for a token type.
func tokenName(t int) string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokICONST:
		return "integer"
	case tokFCONST:
		return "float"
	case tokHCONST:
		return "hex"
	case tokSCONST:
		return "string"
	case tokDCONST:
		return "double-quoted string"
	case tokIDENT:
		return "identifier"
	case tokPARAM:
		return "parameter"
	case tokLPAREN:
		return "'('"
	case tokRPAREN:
		return "')'"
	case tokLBRACK:
		return "'['"
	case tokRBRACK:
		return "']'"
	case tokLBRACE:
		return "'{'"
	case tokRBRACE:
		return "'}'"
	case tokCOMMA:
		return "','"
	case tokCOLON:
		return "':'"
	case tokDOT:
		return "'.'"
	case tokSTAR:
		return "'*'"
	case tokSELECT:
		return "SELECT"
	case tokFROM:
		return "FROM"
	case tokWHERE:
		return "WHERE"
	case tokAND:
		return "AND"
	case tokOR:
		return "OR"
	case tokNOT:
		return "NOT"
	case tokIN:
		return "IN"
	case tokBETWEEN:
		return "BETWEEN"
	case tokLIKE:
		return "LIKE"
	case tokAS:
		return "AS"
	case tokBY:
		return "BY"
	default:
		return fmt.Sprintf("token(%d)", t)
	}
}
