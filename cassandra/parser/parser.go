package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// Parser is the recursive-descent parser for Cassandra CQL.
type Parser struct {
	lexer   *Lexer
	source  string
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

// Parse parses a CQL input containing one or more statements and returns a List of RawStmt.
func Parse(sql string) (*ast.List, error) {
	p := &Parser{
		lexer:  NewLexer(sql),
		source: sql,
	}
	p.advance()

	var items []ast.Node
	for p.cur.Type != tokEOF {
		if p.cur.Type == tokSEMI {
			p.advance()
			continue
		}

		stmtStart := p.cur.Loc
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			break
		}

		stmtEnd := p.prev.End
		raw := &ast.RawStmt{
			Stmt:         stmt,
			StmtLocation: stmtStart,
			StmtLen:      stmtEnd - stmtStart,
		}
		items = append(items, raw)

		// Consume optional trailing semicolons between statements.
		for p.cur.Type == tokSEMI {
			p.advance()
		}
	}

	if p.lexer.Err != nil {
		return nil, p.lexer.Err
	}

	return &ast.List{Items: items}, nil
}

// parseStatement dispatches to the appropriate statement parser.
func (p *Parser) parseStatement() (ast.StmtNode, error) {
	switch p.cur.Type {
	case tokSELECT:
		return p.parseSelect()
	case tokINSERT:
		return p.parseInsert()
	case tokUPDATE:
		return p.parseUpdate()
	case tokDELETE:
		return p.parseDelete()
	case tokBEGIN:
		return p.parseBatch()
	case tokTRUNCATE:
		return p.parseTruncate()
	case tokUSE:
		return p.parseUse()
	case tokCREATE:
		return p.parseCreate()
	case tokALTER:
		return p.parseAlter()
	case tokDROP:
		return p.parseDrop()
	case tokGRANT:
		return p.parseGrant()
	case tokREVOKE:
		return p.parseRevoke()
	case tokLIST:
		return p.parseList()
	case tokAPPLY:
		// APPLY BATCH is handled within parseBatch; standalone APPLY is an error
		return nil, p.errorf("unexpected APPLY without matching BEGIN BATCH")
	default:
		return nil, p.errorf("expected statement, got %s", p.tokenDesc())
	}
}

// ---------------------------------------------------------------------------
// Token manipulation helpers
// ---------------------------------------------------------------------------

func (p *Parser) advance() {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.Next()
	}
}

func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.Next()
		p.hasNext = true
	}
	return p.nextBuf
}

func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.cur.Type == t {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expect(typ int) (Token, error) {
	if p.cur.Type == typ {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, p.errorf("expected %s, got %s", tokenName(typ), p.tokenDesc())
}

func (p *Parser) expectKeyword(typ int) error {
	_, err := p.expect(typ)
	return err
}

func (p *Parser) curLoc() int {
	return p.cur.Loc
}

func (p *Parser) prevEnd() int {
	return p.prev.End
}

func (p *Parser) makeLoc(start int) ast.Loc {
	return ast.Loc{Start: start, End: p.prevEnd()}
}

func (p *Parser) errorf(format string, args ...any) *ParseError {
	if pe, ok := p.lexer.Err.(*ParseError); ok {
		return pe
	}
	line, col := offsetToLineCol(p.lexer.lineIdx, p.cur.Loc)
	near := p.cur.Str
	if near == "" && p.cur.Type != tokEOF {
		near = p.extractNear(p.cur.Loc)
	}
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Loc:     ast.Loc{Start: p.cur.Loc, End: p.cur.End},
		Line:    line,
		Column:  col,
		Near:    near,
	}
}

func (p *Parser) extractNear(offset int) string {
	if offset >= len(p.source) {
		return ""
	}
	end := offset
	for end < len(p.source) && end-offset < 30 && p.source[end] != ' ' && p.source[end] != '\n' && p.source[end] != '\t' {
		end++
	}
	return p.source[offset:end]
}

func (p *Parser) tokenDesc() string {
	switch p.cur.Type {
	case tokEOF:
		return "end of input"
	case tokILLEGAL:
		return fmt.Sprintf("illegal character %q", p.cur.Str)
	default:
		if p.cur.Str != "" {
			return fmt.Sprintf("%q", p.cur.Str)
		}
		return tokenName(p.cur.Type)
	}
}

// ---------------------------------------------------------------------------
// Common parsing helpers
// ---------------------------------------------------------------------------

// parseIdentifier parses an identifier (unquoted, quoted, or keyword-as-ident).
func (p *Parser) parseIdentifier() (*ast.Identifier, error) {
	tok := p.cur
	switch {
	case tok.Type == tokIDENT:
		p.advance()
		return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokQUOTED:
		p.advance()
		return &ast.Identifier{Name: tok.Str, Quoted: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case isKeyword(tok.Type):
		p.advance()
		return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	default:
		return nil, p.errorf("expected identifier, got %s", p.tokenDesc())
	}
}

// parseQualifiedName parses name or keyspace.name.
func (p *Parser) parseQualifiedName() (*ast.QualifiedName, error) {
	start := p.curLoc()
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	parts := []*ast.Identifier{first}
	if p.cur.Type == tokDOT {
		p.advance()
		second, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		parts = append(parts, second)
	}
	return &ast.QualifiedName{Parts: parts, Loc: p.makeLoc(start)}, nil
}

// parseIfNotExists parses optional IF NOT EXISTS, returning (found, error).
func (p *Parser) parseIfNotExists() (bool, error) {
	if p.cur.Type == tokIF {
		next := p.peekNext()
		if next.Type == tokNOT {
			p.advance() // IF
			p.advance() // NOT
			if err := p.expectKeyword(tokEXISTS); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

// parseIfExists parses optional IF EXISTS.
func (p *Parser) parseIfExists() bool {
	if p.cur.Type == tokIF && p.peekNext().Type == tokEXISTS {
		p.advance() // IF
		p.advance() // EXISTS
		return true
	}
	return false
}

// parseUsingClause parses optional USING TTL n [AND TIMESTAMP m] / USING TIMESTAMP m [AND TTL n].
func (p *Parser) parseUsingClause() (*ast.UsingClause, error) {
	if p.cur.Type != tokUSING {
		return nil, nil
	}
	start := p.curLoc()
	p.advance() // USING

	clause := &ast.UsingClause{}
	for {
		switch p.cur.Type {
		case tokTTL:
			p.advance()
			val, err := p.parseConstant()
			if err != nil {
				return nil, err
			}
			clause.TTL = val
		case tokTIMESTAMP:
			p.advance()
			val, err := p.parseConstant()
			if err != nil {
				return nil, err
			}
			clause.Timestamp = val
		default:
			return nil, p.errorf("expected TTL or TIMESTAMP after USING")
		}
		if !p.match(tokAND) {
			break
		}
	}
	clause.Loc = p.makeLoc(start)
	return clause, nil
}

// parseWhereClause parses WHERE relationElement (AND relationElement)*.
func (p *Parser) parseWhereClause() ([]ast.ExprNode, error) {
	if err := p.expectKeyword(tokWHERE); err != nil {
		return nil, err
	}
	return p.parseRelationElements()
}

// parseRelationElements parses relation (AND relation)*.
func (p *Parser) parseRelationElements() ([]ast.ExprNode, error) {
	var relations []ast.ExprNode
	first, err := p.parseRelationElement()
	if err != nil {
		return nil, err
	}
	relations = append(relations, first)
	for p.match(tokAND) {
		rel, err := p.parseRelationElement()
		if err != nil {
			return nil, err
		}
		relations = append(relations, rel)
	}
	return relations, nil
}

// parseRelationElement parses a single WHERE condition.
func (p *Parser) parseRelationElement() (ast.ExprNode, error) {
	start := p.curLoc()

	// Handle tuple comparison: (col1, col2, ...) op/IN (...)
	if p.cur.Type == tokLPAREN {
		return p.parseTupleRelation(start)
	}

	// Could be: identifier op constant, identifier IN (...), identifier CONTAINS [KEY] constant,
	// or function op constant/function.
	left, err := p.parseRelationLeft()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case tokEQ, tokLT, tokGT, tokLTE, tokGTE, tokNE:
		op := p.cur.Str
		p.advance()
		right, err := p.parseRelationRight()
		if err != nil {
			return nil, err
		}
		return &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(start)}, nil

	case tokIN:
		p.advance()
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		var values []ast.ExprNode
		if p.cur.Type != tokRPAREN {
			values, err = p.parseExpressionList()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.InExpr{Column: left, Values: values, Loc: p.makeLoc(start)}, nil

	case tokCONTAINS:
		p.advance()
		isKey := false
		if p.cur.Type == tokKEY {
			isKey = true
			p.advance()
		}
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		return &ast.ContainsExpr{Column: left, Value: val, IsKey: isKey, Loc: p.makeLoc(start)}, nil

	default:
		return nil, p.errorf("expected operator after expression in WHERE clause, got %s", p.tokenDesc())
	}
}

func (p *Parser) parseRelationLeft() (ast.ExprNode, error) {
	if isIdentLike(p.cur.Type) {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Check for function call: name(...)
		if p.cur.Type == tokLPAREN {
			return p.parseFunctionCallWithName(name)
		}
		// Check for dotted name: name.field
		if p.cur.Type == tokDOT {
			start := name.Loc.Start
			p.advance()
			field, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &ast.DotAccess{Object: name, Field: field, Loc: p.makeLoc(start)}, nil
		}
		return name, nil
	}
	return nil, p.errorf("expected identifier or function call in relation, got %s", p.tokenDesc())
}

func (p *Parser) parseRelationRight() (ast.ExprNode, error) {
	// Could be a constant or a function call.
	if isIdentLike(p.cur.Type) && p.peekNext().Type == tokLPAREN {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return p.parseFunctionCallWithName(name)
	}
	return p.parseConstant()
}

func (p *Parser) parseTupleRelation(start int) (ast.ExprNode, error) {
	p.advance() // (
	var cols []ast.ExprNode
	for {
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if !p.match(tokCOMMA) {
			break
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	if p.cur.Type == tokIN {
		p.advance()
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		var tuples []*ast.TupleLit
		for {
			t, err := p.parseTupleLit()
			if err != nil {
				return nil, err
			}
			tuples = append(tuples, t)
			if !p.match(tokCOMMA) {
				break
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.TupleInExpr{Columns: cols, Tuples: tuples, Loc: p.makeLoc(start)}, nil
	}

	// Tuple comparison: (col1, col2) op (val1, val2)
	op := p.cur.Str
	if !p.match(tokEQ, tokLT, tokGT, tokLTE, tokGTE, tokNE) {
		return nil, p.errorf("expected operator or IN after tuple columns")
	}
	var values []ast.ExprNode
	// Could be a single tuple or multiple tuples
	t, err := p.parseTupleLit()
	if err != nil {
		return nil, err
	}
	values = t.Elements
	return &ast.TupleCompareExpr{Columns: cols, Op: op, Values: values, Loc: p.makeLoc(start)}, nil
}

func (p *Parser) parseTupleLit() (*ast.TupleLit, error) {
	start := p.curLoc()
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	elems, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// tokenName returns a human-readable name for a token type.
func tokenName(typ int) string {
	switch {
	case typ == tokEOF:
		return "EOF"
	case typ == tokILLEGAL:
		return "ILLEGAL"
	case typ == tokIDENT:
		return "identifier"
	case typ == tokQUOTED:
		return "quoted identifier"
	case typ == tokSTRING:
		return "string"
	case typ == tokINTEGER:
		return "integer"
	case typ == tokFLOAT:
		return "float"
	case typ == tokUUID:
		return "UUID"
	case typ == tokHEX:
		return "hex literal"
	case typ == tokCODEBLOCK:
		return "code block"
	case typ == tokSEMI:
		return "';'"
	case typ == tokLPAREN:
		return "'('"
	case typ == tokRPAREN:
		return "')'"
	case typ == tokLBRACE:
		return "'{'"
	case typ == tokRBRACE:
		return "'}'"
	case typ == tokLBRACK:
		return "'['"
	case typ == tokRBRACK:
		return "']'"
	case typ == tokCOMMA:
		return "','"
	case typ == tokDOT:
		return "'.'"
	case typ == tokSTAR:
		return "'*'"
	case typ == tokEQ:
		return "'='"
	case typ == tokLT:
		return "'<'"
	case typ == tokGT:
		return "'>'"
	case typ == tokNE:
		return "'!='"
	default:
		// For keywords, reverse lookup
		for name, t := range keywords {
			if t == typ {
				return strings.ToUpper(name)
			}
		}
		return fmt.Sprintf("token(%d)", typ)
	}
}
