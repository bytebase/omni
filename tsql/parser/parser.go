// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the tsql/ast package.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// Parser is a recursive descent parser for T-SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a T-SQL string into an AST list.
// Currently supports basic infrastructure; statement dispatch will be
// implemented incrementally across batches.
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
		stmt := p.parseStmt()
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, &ParseError{
					Message:  "unexpected token in statement",
					Position: p.cur.Loc,
				}
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

// parseStmt dispatches to statement-specific parsers.
// Minimal implementation for batch 0 - only SELECT is supported initially.
// Full dispatch will be implemented in batch 22.
func (p *Parser) parseStmt() nodes.StmtNode {
	switch p.cur.Type {
	case kwSELECT:
		return p.parseSimpleSelectForExpr()
	default:
		return nil
	}
}

// parseSimpleSelectForExpr is a minimal SELECT parser for expression testing.
// It parses: SELECT expr [, expr]...
// Full SELECT support will come in batch 4.
func (p *Parser) parseSimpleSelectForExpr() *nodes.SelectStmt {
	loc := p.pos()
	p.advance() // consume SELECT

	var targets []nodes.Node
	for {
		targetLoc := p.pos() // capture BEFORE parsing expression
		expr := p.parseExpr()
		if expr == nil {
			break
		}
		target := &nodes.ResTarget{
			Val: expr,
			Loc: nodes.Loc{Start: targetLoc},
		}
		// Check for alias: AS name or just name
		if _, ok := p.match(kwAS); ok {
			if tok, ok := p.match(tokIDENT); ok {
				target.Name = tok.Str
			}
		} else if tok, ok := p.match(tokIDENT); ok {
			target.Name = tok.Str
		}
		targets = append(targets, target)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	return &nodes.SelectStmt{
		TargetList: &nodes.List{Items: targets},
		Loc:        nodes.Loc{Start: loc},
	}
}

// parseExpr parses an expression using precedence climbing.
// This is a minimal implementation for batch 0 infrastructure testing.
// Full expression support will come in batch 3.
func (p *Parser) parseExpr() nodes.ExprNode {
	return p.parseOr()
}

func (p *Parser) parseOr() nodes.ExprNode {
	left := p.parseAnd()
	for p.cur.Type == kwOR {
		loc := p.pos()
		p.advance()
		right := p.parseAnd()
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpOr,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

func (p *Parser) parseAnd() nodes.ExprNode {
	left := p.parseNot()
	for p.cur.Type == kwAND {
		loc := p.pos()
		p.advance()
		right := p.parseNot()
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

func (p *Parser) parseNot() nodes.ExprNode {
	if p.cur.Type == kwNOT {
		loc := p.pos()
		p.advance()
		operand := p.parseNot()
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryNot,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc},
		}
	}
	return p.parseComparison()
}

func (p *Parser) parseComparison() nodes.ExprNode {
	left := p.parseAddition()
	if left == nil {
		return nil
	}

	// IS [NOT] NULL
	if p.cur.Type == kwIS {
		loc := p.pos()
		p.advance()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		if p.cur.Type == kwNULL {
			p.advance()
			testType := nodes.IsNull
			if not {
				testType = nodes.IsNotNull
			}
			return &nodes.IsExpr{
				Expr:     left,
				TestType: testType,
				Loc:      nodes.Loc{Start: loc},
			}
		}
	}

	// BETWEEN
	if p.cur.Type == kwBETWEEN || (p.cur.Type == kwNOT && p.peekNext().Type == kwBETWEEN) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume BETWEEN
		low := p.parseAddition()
		if _, err := p.expect(kwAND); err != nil {
			return left
		}
		high := p.parseAddition()
		return &nodes.BetweenExpr{
			Expr: left,
			Low:  low,
			High: high,
			Not:  not,
			Loc:  nodes.Loc{Start: loc},
		}
	}

	// [NOT] IN
	if p.cur.Type == kwIN || (p.cur.Type == kwNOT && p.peekNext().Type == kwIN) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume IN
		if _, err := p.expect('('); err != nil {
			return left
		}
		var items []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			expr := p.parseExpr()
			items = append(items, expr)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		_, _ = p.expect(')')
		return &nodes.InExpr{
			Expr: left,
			List: &nodes.List{Items: items},
			Not:  not,
			Loc:  nodes.Loc{Start: loc},
		}
	}

	// [NOT] LIKE
	if p.cur.Type == kwLIKE || (p.cur.Type == kwNOT && p.peekNext().Type == kwLIKE) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume LIKE
		pattern := p.parseAddition()
		var escape nodes.ExprNode
		if p.cur.Type == kwESCAPE {
			p.advance()
			escape = p.parseAddition()
		}
		return &nodes.LikeExpr{
			Expr:    left,
			Pattern: pattern,
			Escape:  escape,
			Not:     not,
			Loc:     nodes.Loc{Start: loc},
		}
	}

	// Comparison operators
	op, ok := comparisonOp(p.cur.Type)
	if ok {
		loc := p.pos()
		p.advance()
		right := p.parseAddition()
		return &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}

	return left
}

func comparisonOp(tok int) (nodes.BinaryOp, bool) {
	switch tok {
	case '=':
		return nodes.BinOpEq, true
	case tokNOTEQUAL:
		return nodes.BinOpNeq, true
	case '<':
		return nodes.BinOpLt, true
	case '>':
		return nodes.BinOpGt, true
	case tokLESSEQUAL:
		return nodes.BinOpLte, true
	case tokGREATEQUAL:
		return nodes.BinOpGte, true
	case tokNOTLESS:
		return nodes.BinOpNotLt, true
	case tokNOTGREATER:
		return nodes.BinOpNotGt, true
	default:
		return 0, false
	}
}

func (p *Parser) parseAddition() nodes.ExprNode {
	left := p.parseMultiplication()
	for p.cur.Type == '+' || p.cur.Type == '-' || p.cur.Type == '&' || p.cur.Type == '|' || p.cur.Type == '^' {
		loc := p.pos()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '+':
			op = nodes.BinOpAdd
		case '-':
			op = nodes.BinOpSub
		case '&':
			op = nodes.BinOpBitAnd
		case '|':
			op = nodes.BinOpBitOr
		case '^':
			op = nodes.BinOpBitXor
		}
		p.advance()
		right := p.parseMultiplication()
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

func (p *Parser) parseMultiplication() nodes.ExprNode {
	left := p.parseUnary()
	for p.cur.Type == '*' || p.cur.Type == '/' || p.cur.Type == '%' {
		loc := p.pos()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '*':
			op = nodes.BinOpMul
		case '/':
			op = nodes.BinOpDiv
		case '%':
			op = nodes.BinOpMod
		}
		p.advance()
		right := p.parseUnary()
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

func (p *Parser) parseUnary() nodes.ExprNode {
	if p.cur.Type == '+' || p.cur.Type == '-' || p.cur.Type == '~' {
		loc := p.pos()
		var op nodes.UnaryOp
		switch p.cur.Type {
		case '+':
			op = nodes.UnaryPlus
		case '-':
			op = nodes.UnaryMinus
		case '~':
			op = nodes.UnaryBitNot
		}
		p.advance()
		operand := p.parseUnary()
		return &nodes.UnaryExpr{
			Op:      op,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc},
		}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() nodes.ExprNode {
	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitInteger,
			Ival: tok.Ival,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokFCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitFloat,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitString,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokNSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type:    nodes.LitString,
			Str:     tok.Str,
			IsNChar: true,
			Loc:     nodes.Loc{Start: tok.Loc},
		}
	case kwNULL:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitNull,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case kwDEFAULT:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitDefault,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokSYSVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case '*':
		tok := p.advance()
		return &nodes.StarExpr{
			Loc: nodes.Loc{Start: tok.Loc},
		}
	case '(':
		loc := p.pos()
		p.advance()
		expr := p.parseExpr()
		_, _ = p.expect(')')
		return &nodes.ParenExpr{
			Expr: expr,
			Loc:  nodes.Loc{Start: loc},
		}
	case kwCAST:
		return p.parseCast()
	case kwCONVERT:
		return p.parseConvert()
	case kwCASE:
		return p.parseCaseExpr()
	case kwCOALESCE:
		return p.parseCoalesce()
	case kwNULLIF:
		return p.parseNullif()
	case kwIIF:
		return p.parseIif()
	case kwEXISTS:
		return p.parseExists()
	case kwTRY_CAST:
		return p.parseTryCast()
	case kwTRY_CONVERT:
		return p.parseTryConvert()
	case tokIDENT:
		// Could be column ref, function call, or qualified name
		return p.parseIdentExpr()
	default:
		// Many keywords can also be used as identifiers or function names in T-SQL.
		if p.isIdentLike() {
			next := p.peekNext()
			if next.Type == '(' || next.Type == '.' {
				return p.parseIdentExpr()
			}
		}
		return nil
	}
}

func (p *Parser) parseCast() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil
	}
	expr := p.parseExpr()
	if _, err := p.expect(kwAS); err != nil {
		return nil
	}
	dt := p.parseDataType()
	_, _ = p.expect(')')
	return &nodes.CastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseTryCast() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume TRY_CAST
	if _, err := p.expect('('); err != nil {
		return nil
	}
	expr := p.parseExpr()
	if _, err := p.expect(kwAS); err != nil {
		return nil
	}
	dt := p.parseDataType()
	_, _ = p.expect(')')
	return &nodes.TryCastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseConvert() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CONVERT
	if _, err := p.expect('('); err != nil {
		return nil
	}
	dt := p.parseDataType()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	expr := p.parseExpr()
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		style = p.parseExpr()
	}
	_, _ = p.expect(')')
	return &nodes.ConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseTryConvert() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume TRY_CONVERT
	if _, err := p.expect('('); err != nil {
		return nil
	}
	dt := p.parseDataType()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	expr := p.parseExpr()
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		style = p.parseExpr()
	}
	_, _ = p.expect(')')
	return &nodes.TryConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseCaseExpr() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CASE

	var arg nodes.ExprNode
	// Simple CASE vs searched CASE
	if p.cur.Type != kwWHEN {
		arg = p.parseExpr()
	}

	var whenList []nodes.Node
	for p.cur.Type == kwWHEN {
		wloc := p.pos()
		p.advance() // consume WHEN
		cond := p.parseExpr()
		if _, err := p.expect(kwTHEN); err != nil {
			break
		}
		result := p.parseExpr()
		whenList = append(whenList, &nodes.CaseWhen{
			Condition: cond,
			Result:    result,
			Loc:       nodes.Loc{Start: wloc},
		})
	}

	var elseExpr nodes.ExprNode
	if _, ok := p.match(kwELSE); ok {
		elseExpr = p.parseExpr()
	}

	_, _ = p.expect(kwEND)

	return &nodes.CaseExpr{
		Arg:      arg,
		WhenList: &nodes.List{Items: whenList},
		ElseExpr: elseExpr,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseCoalesce() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume COALESCE
	if _, err := p.expect('('); err != nil {
		return nil
	}
	var args []nodes.Node
	for {
		arg := p.parseExpr()
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.CoalesceExpr{
		Args: &nodes.List{Items: args},
		Loc:  nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseNullif() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume NULLIF
	if _, err := p.expect('('); err != nil {
		return nil
	}
	left := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	right := p.parseExpr()
	_, _ = p.expect(')')
	return &nodes.NullifExpr{
		Left:  left,
		Right: right,
		Loc:   nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseIif() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume IIF
	if _, err := p.expect('('); err != nil {
		return nil
	}
	cond := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	trueVal := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	falseVal := p.parseExpr()
	_, _ = p.expect(')')
	return &nodes.IifExpr{
		Condition: cond,
		TrueVal:   trueVal,
		FalseVal:  falseVal,
		Loc:       nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseExists() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume EXISTS
	if _, err := p.expect('('); err != nil {
		return nil
	}
	// For now, parse inner SELECT as a simple select
	query := p.parseStmt()
	_, _ = p.expect(')')
	return &nodes.ExistsExpr{
		Subquery: query,
		Loc:      nodes.Loc{Start: loc},
	}
}

func (p *Parser) parseFuncCall(name string, loc int) nodes.ExprNode {
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Object: name, Loc: nodes.Loc{Start: loc}},
		Loc:  nodes.Loc{Start: loc},
	}

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		_, _ = p.expect(')')
		return fc
	}

	if p.cur.Type == ')' {
		p.advance()
		return fc
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	var args []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		arg := p.parseExpr()
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	fc.Args = &nodes.List{Items: args}

	_, _ = p.expect(')')

	// Check for OVER clause
	if p.cur.Type == kwOVER {
		fc.Over = p.parseOverClause()
	}

	return fc
}

func (p *Parser) parseOverClause() *nodes.OverClause {
	loc := p.pos()
	p.advance() // consume OVER
	if _, err := p.expect('('); err != nil {
		return nil
	}

	over := &nodes.OverClause{
		Loc: nodes.Loc{Start: loc},
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over
		}
		var parts []nodes.Node
		for {
			expr := p.parseExpr()
			parts = append(parts, expr)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		over.PartitionBy = &nodes.List{Items: parts}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over
		}
		var orders []nodes.Node
		for {
			expr := p.parseExpr()
			dir := nodes.SortDefault
			if _, ok := p.match(kwASC); ok {
				dir = nodes.SortAsc
			} else if _, ok := p.match(kwDESC); ok {
				dir = nodes.SortDesc
			}
			orders = append(orders, &nodes.OrderByItem{
				Expr:    expr,
				SortDir: dir,
				Loc:     nodes.Loc{Start: p.pos()},
			})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		over.OrderBy = &nodes.List{Items: orders}
	}

	_, _ = p.expect(')')
	return over
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
		Message:  "unexpected token",
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}
