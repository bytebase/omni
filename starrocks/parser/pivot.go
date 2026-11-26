package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// atPivotClause reports whether cur starts a pivotClause (StarRocks.g4:2202),
// PIVOT '(' ... rather than PIVOT as a relation's alias. PIVOT is nonReserved,
// so `FROM t pivot` names an alias, and after a subquery, a VALUES list or a
// table function `pivot (a, b)` is an alias with column aliases. A pivot's
// first aggregation is a function call, never a name followed by ',' or ')'.
func (p *Parser) atPivotClause() bool {
	if p.cur.Kind != kwPIVOT || p.peekNext().Kind != int('(') {
		return false
	}
	mark := p.save()
	defer p.restore(mark)
	p.advance() // PIVOT
	p.advance() // '('
	if !isIdentifierToken(p.cur.Kind) {
		return true
	}
	next := p.peekNext().Kind
	return next != int(',') && next != int(')')
}

// parsePivotClause parses
//
//	pivotClause
//	    : PIVOT '(' pivotAggregationExpression (',' pivotAggregationExpression)*
//	        FOR (identifier | identifierList) IN '(' pivotValue (',' pivotValue)* ')' ')'
//
// (StarRocks.g4:2202). cur is PIVOT. AstBuilder.visitPivotClause refuses a
// value whose literal count differs from the FOR columns'.
func (p *Parser) parsePivotClause() (*ast.PivotClause, error) {
	pivotTok := p.advance() // PIVOT
	clause := &ast.PivotClause{Loc: ast.Loc{Start: pivotTok.Loc.Start}}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	for {
		agg, err := p.parsePivotAggregation()
		if err != nil {
			return nil, err
		}
		clause.Aggregations = append(clause.Aggregations, agg)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if p.cur.Kind == int('(') {
		cols, err := p.parseColumnAliasList()
		if err != nil {
			return nil, err
		}
		clause.Columns = cols
	} else {
		col, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		clause.Columns = []string{col}
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	for {
		val, err := p.parsePivotValue()
		if err != nil {
			return nil, err
		}
		if len(val.Exprs) != len(clause.Columns) {
			return nil, &ParseError{
				Position: val.Loc.Start, End: val.Loc.End,
				Message: fmt.Sprintf("PIVOT value has %d literals for %d columns", len(val.Exprs), len(clause.Columns)),
			}
		}
		clause.Values = append(clause.Values, val)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End
	return clause, nil
}

// parsePivotAggregation parses pivotAggregationExpression, functionCall
// (AS? (identifier | string))? (StarRocks.g4:2207).
// AstBuilder.visitPivotAggregationExpression casts the call to a
// FunctionCallExpr, so the call forms the visitor builds as another
// expression are refused: see isPivotFunctionCall.
func (p *Parser) parsePivotAggregation() (*ast.PivotAggregation, error) {
	start := p.cur
	fn, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	if !p.isPivotFunctionCall(fn) {
		return nil, &ParseError{
			Position: start.Loc.Start, End: ast.NodeLoc(fn).End,
			Message: "PIVOT takes an aggregate function call",
		}
	}
	agg := &ast.PivotAggregation{Func: fn, Loc: ast.Loc{Start: start.Loc.Start}}
	if agg.Alias, err = p.parsePivotAlias(); err != nil {
		return nil, err
	}
	agg.Loc.End = p.prev.Loc.End
	return agg, nil
}

// isPivotFunctionCall reports whether fn, a functionCall, is one the
// StarRocks 3.4.10 AstBuilder builds as a FunctionCallExpr. It builds no
// FunctionCallExpr for a call with OVER (an AnalyticExpr), an
// informationFunctionExpression (CURRENT_USER, DATABASE() and the rest),
// PASSWORD('x') (a string), TIMESTAMPADD and TIMESTAMPDIFF, or the plain
// calls visitSimpleFunctionCall rewrites (pivotRewrittenCalls).
func (p *Parser) isPivotFunctionCall(fn ast.Node) bool {
	switch n := fn.(type) {
	case *ast.ExtractExpr:
		return true
	case *ast.FuncCallExpr:
		if n.Over != nil || n.Name == nil || len(n.Name.Parts) == 0 {
			return false
		}
		parts := n.Name.Parts
		if len(parts) == 1 && !p.isQuotedName(n.Name) {
			switch strings.ToUpper(parts[0]) {
			case "MOD":
				// specialFunctionExpression's MOD, a FunctionCallExpr.
				return true
			case "CURRENT_USER", "CURRENT_ROLE", "TIMESTAMPADD", "TIMESTAMPDIFF":
				return false
			case "CATALOG", "DATABASE", "SCHEMA", "USER":
				// informationFunctionExpression takes no argument; with one,
				// CATALOG and USER are nonReserved names of a plain call.
				return len(n.Args) > 0 || n.Star
			case "PASSWORD":
				// PASSWORD '(' string ')' builds a string; with any other
				// argument the nonReserved word names a plain call.
				if len(n.Args) == 1 {
					if lit, ok := n.Args[0].(*ast.Literal); ok && lit.Kind == ast.LitString {
						return false
					}
				}
			}
		}
		// FunctionName.createFnName: a db and a function for two parts, the
		// whole name otherwise; lowercased.
		name := strings.Join(parts, ".")
		if len(parts) == 2 {
			name = parts[1]
		}
		return !pivotRewrittenCalls[strings.ToLower(name)]
	}
	return false
}

// pivotRewrittenCalls are the function names AstBuilder.visitSimpleFunctionCall
// (StarRocks 3.4.10) builds as something other than a FunctionCallExpr:
// DATE_FUNCTIONS (a TimestampArithmeticExpr), element_at, isnull and
// isnotnull, ArithmeticExpr's SUPPORT_FUNCTIONS (whose bit-shift keys are
// camel case and never match the lowercased name), connection_id,
// session_user and session_id, map and dict_mapping.
var pivotRewrittenCalls = map[string]bool{
	"date_add": true, "adddate": true, "date_sub": true, "subdate": true, "days_sub": true,
	"element_at": true, "isnull": true, "isnotnull": true,
	"multiply": true, "divide": true, "mod": true, "int_divide": true, "add": true,
	"subtract": true, "bitand": true, "bitor": true, "bitxor": true,
	"connection_id": true, "session_user": true, "session_id": true,
	"map": true, "dict_mapping": true,
}

// parsePivotValue parses pivotValue, (literalExpression |
// literalExpressionList) (AS? (identifier | string))? (StarRocks.g4:2212),
// where literalExpressionList is '(' literalExpression (','
// literalExpression)* ')' (StarRocks.g4:2637).
func (p *Parser) parsePivotValue() (*ast.PivotValue, error) {
	val := &ast.PivotValue{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	if p.cur.Kind == int('(') {
		p.advance()
		for {
			lit, err := p.parsePivotLiteral()
			if err != nil {
				return nil, err
			}
			val.Exprs = append(val.Exprs, lit)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	} else {
		lit, err := p.parsePivotLiteral()
		if err != nil {
			return nil, err
		}
		val.Exprs = []ast.Node{lit}
	}
	val.Loc.End = p.prev.Loc.End
	alias, err := p.parsePivotAlias()
	if err != nil {
		return nil, err
	}
	val.Alias = alias
	val.Loc.End = p.prev.Loc.End
	return val, nil
}

// parsePivotLiteral parses a literalExpression (StarRocks.g4:2428) that
// AstBuilder.visitPivotValue can cast to a LiteralExpr: NULL, a boolean, an
// unsigned number, a date literal, a string, an interval or a binary
// string. A unitBoundary (FLOOR, CEIL) and a parameter (?) build no
// LiteralExpr, and the pinned server refuses both here.
func (p *Parser) parsePivotLiteral() (ast.Node, error) {
	tok := p.cur
	switch tok.Kind {
	case kwNULL, kwTRUE, kwFALSE, tokInt, tokFloat, tokString, tokHexLiteral:
		kind := map[TokenKind]ast.LitKind{
			kwNULL: ast.LitNull, kwTRUE: ast.LitBool, kwFALSE: ast.LitBool,
			tokInt: ast.LitInt, tokFloat: ast.LitFloat, tokString: ast.LitString, tokHexLiteral: ast.LitHex,
		}[tok.Kind]
		p.advance()
		return &ast.Literal{Kind: kind, Value: tok.Str, Loc: tok.Loc}, nil
	case kwDATE, kwDATETIME:
		if p.peekNext().Kind == tokString {
			return p.parseDateLiteral()
		}
	case kwINTERVAL:
		return p.parseIntervalExpr()
	}
	return nil, p.syntaxErrorAtCur()
}

// parsePivotAlias parses the optional alias of a pivot aggregation or value,
// AS? (identifier | string). It returns "" when there is none.
func (p *Parser) parsePivotAlias() (string, error) {
	if p.cur.Kind == kwAS {
		p.advance()
		if p.cur.Kind == tokString {
			return p.advance().Str, nil
		}
		name, _, err := p.parseIdentifier()
		return name, err
	}
	if p.cur.Kind == tokString {
		return p.advance().Str, nil
	}
	if isIdentifierToken(p.cur.Kind) {
		name, _, err := p.parseIdentifier()
		return name, err
	}
	return "", nil
}
