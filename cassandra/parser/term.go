package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// A term is a value written in a statement: a constant, a collection, UDT or
// tuple literal, a bind marker, a function call, a type hint or a CAST,
// combined with + - * / % (Parser.g:1640-1666). A term never names a column:
// INSERT ... VALUES (a), SET a = b and WHERE a IN (b) are not CQL.

// parserState is what mark saves so that reset can return the parser to an
// earlier token: the parser's tokens and the lexer's position.
type parserState struct {
	cur, prev, nextBuf Token
	hasNext            bool
	pos                int
	err                error
	depth              int
}

// mark saves the parser's position for reset, where the grammar needs more
// lookahead than peekNext gives: a syntactic predicate of Parser.g.
func (p *Parser) mark() parserState {
	return parserState{p.cur, p.prev, p.nextBuf, p.hasNext, p.lexer.pos, p.lexer.Err, p.depth}
}

// reset returns the parser to a position mark saved.
func (p *Parser) reset(s parserState) {
	p.cur, p.prev, p.nextBuf, p.hasNext = s.cur, s.prev, s.nextBuf, s.hasNext
	p.lexer.pos, p.lexer.Err, p.depth = s.pos, s.err, s.depth
}

// isLiteralKeyword reports whether typ is a keyword that is a constant.
func isLiteralKeyword(typ int) bool {
	switch typ {
	case tokTRUE, tokFALSE, tokNULL, tokNAN, tokINFINITY:
		return true
	}
	return false
}

// isNameToken reports whether tok can be a name (Parser.g ident, cident,
// fident, noncol_ident): not a constant keyword such as NULL.
func isNameToken(tok Token) bool {
	return isIdentLike(tok.Type) && !isLiteralKeyword(tok.Type)
}

// isConstantToken reports whether tok starts a constant (Parser.g:1531):
// a string, number, boolean, duration, UUID, hex number, NaN or Infinity.
func isConstantToken(tok Token) bool {
	switch tok.Type {
	case tokSTRING, tokINTEGER, tokFLOAT, tokUUID, tokHEX, tokCODEBLOCK, tokDURATIONLIT,
		tokTRUE, tokFALSE, tokNAN, tokINFINITY:
		return true
	}
	return false
}

// isFunctionNameToken reports whether tok can name a function
// (Parser.g:1611 allowedFunctionName): a name, TOKEN or COUNT, but not one
// of the keywords only a column may be named (TTL, WRITETIME, MAXWRITETIME,
// KEY, CAST, JSON, DISTINCT: Parser.g:1952 unreserved_keyword).
func isFunctionNameToken(tok Token) bool {
	switch {
	case tok.Type == tokKEY, tok.Type == tokCAST, tok.Type == tokJSON, tok.Type == tokDISTINCT,
		isWord(tok, "ttl"), isWord(tok, "writetime"), isWord(tok, "maxwritetime"):
		return false
	case isWord(tok, "token"), isWord(tok, "count"):
		return true
	}
	return isNameToken(tok)
}

// startsTerm reports whether tok, followed by next, can start a term.
func startsTerm(tok, next Token) bool {
	return tok.Type == tokMINUS || startsSimpleTerm(tok, next)
}

// startsSimpleTerm reports whether tok, followed by next, can start a
// simpleTerm (Parser.g:1661).
func startsSimpleTerm(tok, next Token) bool {
	switch {
	case isConstantToken(tok), isSignedConstant(tok, next):
		return true
	case tok.Type == tokNULL, tok.Type == tokQMARK, tok.Type == tokCOLON,
		tok.Type == tokLPAREN, tok.Type == tokLBRACK, tok.Type == tokLBRACE:
		return true
	case tok.Type == tokCAST:
		return next.Type == tokLPAREN
	case isFunctionNameToken(tok):
		return next.Type == tokLPAREN || next.Type == tokDOT
	}
	return false
}

// parseTerm parses a term (Parser.g:1640 termAddition): termMultiplication
// joined by + and -, left to right.
func (p *Parser) parseTerm() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	left, err := p.parseTermMultiplication()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		// A minus written against a number is that number's sign to
		// Cassandra's lexer (Lexer.g:337), never an operator: 1-2 is two
		// constants, which no term takes.
		if isSignedConstant(p.cur, p.peekNext()) {
			break
		}
		// X = <value> + X (Parser.g:1675) puts a column after a term's +;
		// that + is not the term's, since no term follows it.
		if !p.termFollowsOperator() {
			break
		}
		op := p.cur.Str
		p.advance()
		right, err := p.parseTermMultiplication()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(left.GetLoc().Start)}
	}
	return left, nil
}

// termFollowsOperator reports whether a term follows the operator at p.cur.
func (p *Parser) termFollowsOperator() bool {
	s := p.mark()
	p.advance()
	ok := startsTerm(p.cur, p.peekNext())
	p.reset(s)
	return ok
}

// parseTermMultiplication parses termGroup joined by *, / and %
// (Parser.g:1647).
func (p *Parser) parseTermMultiplication() (ast.ExprNode, error) {
	left, err := p.parseTermGroup()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokSTAR || p.cur.Type == tokSLASH || p.cur.Type == tokPERCENT {
		op := p.cur.Str
		p.advance()
		right, err := p.parseTermGroup()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(left.GetLoc().Start)}
	}
	return left, nil
}

// parseTermGroup parses a simpleTerm, negated by a leading minus
// (Parser.g:1656).
func (p *Parser) parseTermGroup() (ast.ExprNode, error) {
	if p.cur.Type != tokMINUS {
		return p.parseSimpleTerm()
	}
	start := p.curLoc()
	p.advance()
	operand, err := p.parseSimpleTerm()
	if err != nil {
		return nil, err
	}
	return negate(operand, start, p.prevEnd()), nil
}

// negate returns - operand, spanning start to end. A number is negated in
// its literal, as parseConstant does.
func negate(operand ast.ExprNode, start, end int) ast.ExprNode {
	loc := ast.Loc{Start: start, End: end}
	switch lit := operand.(type) {
	case *ast.IntegerLit:
		if !strings.HasPrefix(lit.Val, "-") {
			return &ast.IntegerLit{Val: "-" + lit.Val, Loc: loc}
		}
	case *ast.FloatLit:
		if !strings.HasPrefix(lit.Val, "-") {
			return &ast.FloatLit{Val: "-" + lit.Val, Loc: loc}
		}
	case *ast.DurationLit:
		if !strings.HasPrefix(lit.Val, "-") {
			return &ast.DurationLit{Val: "-" + lit.Val, Loc: loc}
		}
	}
	return &ast.UnaryExpr{Op: "-", Operand: operand, Loc: loc}
}

// parseSimpleTerm parses a simpleTerm (Parser.g:1661): a value, a function
// call, (type) simpleTerm, or CAST(simpleTerm AS native_type).
func (p *Parser) parseSimpleTerm() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	switch {
	case p.cur.Type == tokLPAREN:
		if hint, ok, err := p.tryTypeHint(func(tok, next Token) bool { return startsSimpleTerm(tok, next) }, p.parseSimpleTerm); ok || err != nil {
			return hint, err
		}
		return p.parseTupleLiteral()
	case p.cur.Type == tokCAST && p.peekNext().Type == tokLPAREN:
		return p.parseCastWith(p.parseSimpleTerm)
	case p.cur.Type == tokLBRACK:
		return p.parseListTerm()
	case p.cur.Type == tokLBRACE:
		return p.parseBraceTerm()
	case p.cur.Type == tokQMARK || p.cur.Type == tokCOLON:
		return p.parseBindMarker()
	case p.cur.Type == tokNULL:
		tok := p.cur
		p.advance()
		return &ast.NullLit{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case isConstantToken(p.cur):
		return p.parseConstant()
	case isSignedConstant(p.cur, p.peekNext()):
		start := p.curLoc()
		p.advance() // -
		operand, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		return negate(operand, start, p.prevEnd()), nil
	case isFunctionNameToken(p.cur) && (p.peekNext().Type == tokLPAREN || p.peekNext().Type == tokDOT):
		return p.parseTermFunction()
	}
	return nil, p.errorf("expected term, got %s", p.tokenDesc())
}

// isSignedConstant reports whether tok is a minus written against the
// number, duration, NaN or Infinity next: Cassandra's lexer reads the sign as
// part of that one token (Lexer.g:337 INTEGER, FLOAT, DURATION,
// K_NEGATIVE_NAN, K_NEGATIVE_INFINITY), so it is a constant wherever a
// simpleTerm stands, such as CAST(-1 AS int). omni's lexer reads the minus
// apart.
func isSignedConstant(tok, next Token) bool {
	if tok.Type != tokMINUS || next.Loc != tok.End {
		return false
	}
	switch next.Type {
	case tokINTEGER, tokFLOAT, tokDURATIONLIT, tokNAN, tokINFINITY:
		return true
	}
	return false
}

// tryTypeHint parses '(' comparatorType ')' and what follows it, read by
// parse, when p.cur is '(' and a type in parentheses is followed by a token
// starts accepts. Otherwise it reads nothing and returns ok false: the '('
// opens a tuple or a parenthesised selector.
func (p *Parser) tryTypeHint(starts func(tok, next Token) bool, parse func() (ast.ExprNode, error)) (ast.ExprNode, bool, error) {
	s := p.mark()
	start := p.curLoc()
	p.advance() // (
	dt, err := p.parseDataType()
	if err != nil || p.cur.Type != tokRPAREN {
		p.reset(s)
		return nil, false, nil
	}
	p.advance() // )
	if !starts(p.cur, p.peekNext()) {
		p.reset(s)
		return nil, false, nil
	}
	expr, err := parse()
	if err != nil {
		return nil, false, err
	}
	return &ast.TypeHintExpr{Type: dt, Expr: expr, Loc: p.makeLoc(start)}, true, nil
}

// parseCastWith parses CAST '(' operand AS native_type ')', the operand read
// by parse.
func (p *Parser) parseCastWith(parse func() (ast.ExprNode, error)) (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // CAST
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	expr, err := parse()
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokAS); err != nil {
		return nil, err
	}
	dt, err := p.parseNativeType()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.CastExpr{Expr: expr, Type: dt, Loc: p.makeLoc(start)}, nil
}

// parseNativeType parses a native_type (Parser.g:1872), the only types CAST
// takes.
func (p *Parser) parseNativeType() (*ast.DataType, error) {
	switch p.cur.Type {
	case tokASCII, tokBIGINT, tokBLOB, tokBOOLEAN, tokCOUNTER, tokDECIMAL, tokDOUBLE,
		tokDURATION, tokFLOATKW, tokINET, tokINT, tokSMALLINT, tokTEXT, tokTIMESTAMP,
		tokTINYINT, tokUUID_KW, tokVARCHAR, tokVARINT, tokTIMEUUID, tokDATE, tokTIME:
		tok := p.cur
		p.advance()
		loc := ast.Loc{Start: tok.Loc, End: tok.End}
		return &ast.DataType{Name: &ast.Identifier{Name: tok.Str, Loc: loc}, Loc: loc}, nil
	}
	return nil, p.errorf("expected native type, got %s", p.tokenDesc())
}

// parseBindMarker parses ? or :name.
func (p *Parser) parseBindMarker() (*ast.BindMarker, error) {
	if p.cur.Type == tokQMARK {
		m := &ast.BindMarker{Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		return m, nil
	}
	start := p.curLoc()
	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}
	if !isUnreservedName(p.cur) {
		return nil, p.errorf("expected bind marker name, got %s", p.tokenDesc())
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.BindMarker{Name: name.Name, Loc: p.makeLoc(start)}, nil
}

// parseTupleLiteral parses '(' term (',' term)* ')' (Parser.g:1587).
func (p *Parser) parseTupleLiteral() (*ast.TupleLit, error) {
	start := p.curLoc()
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	elems, err := p.parseTermList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseTermList parses term (',' term)*.
func (p *Parser) parseTermList() ([]ast.ExprNode, error) {
	var list []ast.ExprNode
	for {
		t, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		list = append(list, t)
		if !p.match(tokCOMMA) {
			return list, nil
		}
	}
}

// parseListTerm parses '[' (term (',' term)*)? ']' (Parser.g:1572).
func (p *Parser) parseListTerm() (*ast.ListLit, error) {
	start := p.curLoc()
	p.advance() // [
	var elems []ast.ExprNode
	if p.cur.Type != tokRBRACK {
		var err error
		if elems, err = p.parseTermList(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}
	return &ast.ListLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseBraceTerm parses a set, map or UDT literal (Parser.g:1564
// collectionLiteral, 1578 usertypeLiteral). A UDT literal's keys are field
// names, kept as Identifier keys of a MapLit; {} is an empty MapLit.
func (p *Parser) parseBraceTerm() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // {
	if p.match(tokRBRACE) {
		return &ast.MapLit{Loc: p.makeLoc(start)}, nil
	}
	if isNameToken(p.cur) && p.peekNext().Type == tokCOLON {
		var keys, values []ast.ExprNode
		for {
			k, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokCOLON); err != nil {
				return nil, err
			}
			v, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			keys, values = append(keys, k), append(values, v)
			if !p.match(tokCOMMA) {
				break
			}
		}
		if _, err := p.expect(tokRBRACE); err != nil {
			return nil, err
		}
		return &ast.MapLit{Keys: keys, Values: values, Loc: p.makeLoc(start)}, nil
	}
	first, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	if p.match(tokCOLON) {
		v, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		keys, values := []ast.ExprNode{first}, []ast.ExprNode{v}
		for p.match(tokCOMMA) {
			k, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokCOLON); err != nil {
				return nil, err
			}
			v, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			keys, values = append(keys, k), append(values, v)
		}
		if _, err := p.expect(tokRBRACE); err != nil {
			return nil, err
		}
		return &ast.MapLit{Keys: keys, Values: values, Loc: p.makeLoc(start)}, nil
	}
	elems := []ast.ExprNode{first}
	for p.match(tokCOMMA) {
		e, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		elems = append(elems, e)
	}
	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	return &ast.SetLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseFunctionName parses (keyspaceName '.')? allowedFunctionName
// (Parser.g:1605). A built-in function keyword is named in lower case.
func (p *Parser) parseFunctionName() (*ast.Identifier, error) {
	if !isFunctionNameToken(p.cur) {
		return nil, p.errorf("expected function name, got %s", p.tokenDesc())
	}
	if p.peekNext().Type == tokDOT {
		ks, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		p.advance() // .
		if !isFunctionNameToken(p.cur) {
			return nil, p.errorf("expected function name, got %s", p.tokenDesc())
		}
		fn, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.Identifier{Name: ks.Name + "." + fn.Name, Loc: ast.Loc{Start: ks.Loc.Start, End: fn.Loc.End}}, nil
	}
	if name, ok := builtinFuncToken[p.cur.Type]; ok {
		tok := p.cur
		p.advance()
		return &ast.Identifier{Name: name, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	}
	return p.parseIdentifier()
}

// parseTermFunction parses functionName '(' (term (',' term)*)? ')'
// (Parser.g:1619).
func (p *Parser) parseTermFunction() (*ast.FunctionCall, error) {
	name, err := p.parseFunctionName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	var args []ast.ExprNode
	if p.cur.Type != tokRPAREN {
		if args, err = p.parseTermList(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.FunctionCall{Name: name, Args: args, Loc: p.makeLoc(name.Loc.Start)}, nil
}

// isUnreservedName reports whether tok can be a name where Parser.g wants a
// field name (fident) or a bind marker name (noncol_ident): a name that is
// not a reserved word.
func isUnreservedName(tok Token) bool {
	if tok.Type == tokQUOTED {
		return true
	}
	return isNameToken(tok) && !reservedKeywords[strings.ToLower(tok.Str)]
}
