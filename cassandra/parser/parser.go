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
	depth   int // nesting depth of the productions being parsed (see enterNesting)
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

		first := p.cur
		stmtStart := p.cur.Loc
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			break
		}

		stmtEnd := p.prev.End
		if treeTooHigh(stmt, stmtEnd-stmtStart) {
			return nil, p.treeTooDeepError(first)
		}
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
	case tokDESCRIBE, tokDESC:
		return p.parseDescribe()
	case tokADD:
		return p.parseAddIdentity()
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
		Message:  fmt.Sprintf(format, args...),
		Position: p.cur.Loc, End: p.cur.End,
		Line:   line,
		Column: col,
		Near:   near,
	}
}

// errorAt is errorf for an earlier token: a refusal the parser decides only
// after reading past what it refuses, such as a column declared twice.
func (p *Parser) errorAt(tok Token, format string, args ...any) *ParseError {
	if pe, ok := p.lexer.Err.(*ParseError); ok {
		return pe
	}
	line, col := offsetToLineCol(p.lexer.lineIdx, tok.Loc)
	near := tok.Str
	if near == "" && tok.Type != tokEOF {
		near = p.extractNear(tok.Loc)
	}
	return &ParseError{
		Message:  fmt.Sprintf(format, args...),
		Position: tok.Loc, End: tok.End,
		Line:   line,
		Column: col,
		Near:   near,
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
	case reservedKeywords[tokWord(tok)]:
		// A word Lexer.g reads as a keyword that unreserved_keyword
		// (Parser.g:1946) leaves out is never a name; quote it to use it.
		return nil, p.errorf("expected identifier, got reserved word %s", p.tokenDesc())
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

// isWord reports whether tok is the unquoted word w, read as an identifier
// or as a keyword, in any case. CQL's unreserved keywords that omni's lexer
// reads as identifiers (CLUSTER, MASKED, USERS, IDENTITY, ...) are matched by
// their text, so they stay names everywhere else.
func isWord(tok Token, w string) bool {
	return (tok.Type == tokIDENT || isKeyword(tok.Type)) && strings.EqualFold(tok.Str, w)
}

// curIs reports whether the current token is the word w (see isWord).
func (p *Parser) curIs(w string) bool { return isWord(p.cur, w) }

// identKey is the name Cassandra compares an identifier by: an unquoted name
// folds to lower case, a quoted one is kept as written (ColumnIdentifier).
func identKey(id *ast.Identifier) string {
	if id.Quoted {
		return id.Name
	}
	return strings.ToLower(id.Name)
}

// parseRoleName parses roleName (Parser.g:1519): an identifier, a quoted
// name, an unreserved keyword or a string literal. A string comes back as a
// quoted Identifier, since Cassandra keeps its case.
func (p *Parser) parseRoleName() (*ast.Identifier, error) {
	if isStringLiteral(p.cur.Type) {
		tok := p.cur
		p.advance()
		return &ast.Identifier{Name: tok.Str, Quoted: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	}
	return p.parseIdentifier()
}

// parseUsername parses username (Parser.g:1921) and identity (Parser.g:1927):
// an identifier or a string literal. A quoted name is refused by the
// grammar's own action, and no keyword of Lexer.g is a user name; a word
// only omni's lexer makes a keyword (REPLICATION, NOW, SAI) is an IDENT to
// Cassandra.
func (p *Parser) parseUsername(what string) (*ast.Identifier, error) {
	tok := p.cur
	switch {
	case tokWord(tok) != "" && !isLexerKeyword(tokWord(tok)):
		p.advance()
		return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case isStringLiteral(tok.Type):
		p.advance()
		return &ast.Identifier{Name: tok.Str, Quoted: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	default:
		return nil, p.errorf("expected %s (identifier or string), got %s", what, p.tokenDesc())
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
	case typ == tokQMARK:
		return "'?'"
	case typ == tokRANGE:
		return "'..'"
	case typ == tokPLUSEQ:
		return "'+='"
	case typ == tokMINUSEQ:
		return "'-='"
	case typ == tokDURATIONLIT:
		return "duration"
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
