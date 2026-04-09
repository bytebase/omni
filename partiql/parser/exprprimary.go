package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parsePrimary produces a primary expression and attaches any trailing
// path steps. It's the entry point for the exprPrimary + pathStep+
// combination in the grammar (line 528):
//
//	exprPrimary pathStep+      # ExprPrimaryPath
//
// Grammar: exprPrimary (lines 514-534).
func (p *Parser) parsePrimary() (ast.ExprNode, error) {
	base, err := p.parsePrimaryBase()
	if err != nil {
		return nil, err
	}
	if !isPathStepStart(p.cur.Type) {
		return base, nil
	}
	return p.parsePathSteps(base)
}

// parsePrimaryBase dispatches on the current token to produce one of
// the 16 primary-expression alternatives. Foundation handles exprTerm
// alternatives directly; every other alternative calls deferredFeature
// to return a stub error pointing at the owning DAG node.
//
// Grammar: exprPrimary base alternatives (lines 514-534) + exprTerm
// alternatives (lines 542-549).
func (p *Parser) parsePrimaryBase() (ast.ExprNode, error) {
	switch p.cur.Type {
	// ------------------------------------------------------------------
	// Real: literal primary forms (delegates to literals.go).
	// ------------------------------------------------------------------
	case tokNULL, tokMISSING, tokTRUE, tokFALSE,
		tokSCONST, tokICONST, tokFCONST, tokION_LITERAL,
		tokDATE, tokTIME:
		return p.parseLiteral()

	// ------------------------------------------------------------------
	// Real: parenthesized expr (valueList and (SELECT ...) are stubbed
	// inside parseParenExpr).
	// ------------------------------------------------------------------
	case tokPAREN_LEFT:
		return p.parseParenExpr()

	// ------------------------------------------------------------------
	// Real: collection literals.
	// ------------------------------------------------------------------
	case tokBRACKET_LEFT:
		return p.parseArrayLit()
	case tokANGLE_DOUBLE_LEFT:
		return p.parseBagLit()
	case tokBRACE_LEFT:
		return p.parseTupleLit()

	// ------------------------------------------------------------------
	// Real: parameter and varRef.
	// ------------------------------------------------------------------
	case tokQUESTION_MARK:
		return p.parseParamRef()
	case tokAT_SIGN, tokIDENT, tokIDENT_QUOTED:
		return p.parseVarRef()

	// ------------------------------------------------------------------
	// Stub: VALUES row list (no AST node yet).
	// ------------------------------------------------------------------
	case tokVALUES:
		return nil, p.deferredFeature("VALUES", "parser-dml (DAG node 6)")

	// ------------------------------------------------------------------
	// Stub: CAST family → parser-builtins
	// ------------------------------------------------------------------
	case tokCAST:
		return nil, p.deferredFeature("CAST", "parser-builtins (DAG node 15)")
	case tokCAN_CAST:
		return nil, p.deferredFeature("CAN_CAST", "parser-builtins (DAG node 15)")
	case tokCAN_LOSSLESS_CAST:
		return nil, p.deferredFeature("CAN_LOSSLESS_CAST", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: CASE expression → parser-builtins
	// ------------------------------------------------------------------
	case tokCASE:
		return nil, p.deferredFeature("CASE", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: keyword-bearing builtin functions → parser-builtins
	// ------------------------------------------------------------------
	case tokCOALESCE:
		return nil, p.deferredFeature("COALESCE", "parser-builtins (DAG node 15)")
	case tokNULLIF:
		return nil, p.deferredFeature("NULLIF", "parser-builtins (DAG node 15)")
	case tokSUBSTRING:
		return nil, p.deferredFeature("SUBSTRING", "parser-builtins (DAG node 15)")
	case tokTRIM:
		return nil, p.deferredFeature("TRIM", "parser-builtins (DAG node 15)")
	case tokEXTRACT:
		return nil, p.deferredFeature("EXTRACT", "parser-builtins (DAG node 15)")
	case tokDATE_ADD:
		return nil, p.deferredFeature("DATE_ADD", "parser-builtins (DAG node 15)")
	case tokDATE_DIFF:
		return nil, p.deferredFeature("DATE_DIFF", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: reserved-name scalar functions → parser-builtins
	// ------------------------------------------------------------------
	case tokCHAR_LENGTH:
		return nil, p.deferredFeature("CHAR_LENGTH", "parser-builtins (DAG node 15)")
	case tokCHARACTER_LENGTH:
		return nil, p.deferredFeature("CHARACTER_LENGTH", "parser-builtins (DAG node 15)")
	case tokOCTET_LENGTH:
		return nil, p.deferredFeature("OCTET_LENGTH", "parser-builtins (DAG node 15)")
	case tokBIT_LENGTH:
		return nil, p.deferredFeature("BIT_LENGTH", "parser-builtins (DAG node 15)")
	case tokUPPER:
		return nil, p.deferredFeature("UPPER", "parser-builtins (DAG node 15)")
	case tokLOWER:
		return nil, p.deferredFeature("LOWER", "parser-builtins (DAG node 15)")
	case tokSIZE:
		return nil, p.deferredFeature("SIZE", "parser-builtins (DAG node 15)")
	case tokEXISTS:
		return nil, p.deferredFeature("EXISTS", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: sequenceConstructor (LIST/SEXP) → parser-builtins
	// ------------------------------------------------------------------
	case tokLIST:
		return nil, p.deferredFeature("LIST() constructor", "parser-builtins (DAG node 15)")
	case tokSEXP:
		return nil, p.deferredFeature("SEXP() constructor", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: aggregates → parser-aggregates
	// ------------------------------------------------------------------
	case tokCOUNT:
		return nil, p.deferredFeature("COUNT() aggregate", "parser-aggregates (DAG node 14)")
	case tokMAX:
		return nil, p.deferredFeature("MAX() aggregate", "parser-aggregates (DAG node 14)")
	case tokMIN:
		return nil, p.deferredFeature("MIN() aggregate", "parser-aggregates (DAG node 14)")
	case tokSUM:
		return nil, p.deferredFeature("SUM() aggregate", "parser-aggregates (DAG node 14)")
	case tokAVG:
		return nil, p.deferredFeature("AVG() aggregate", "parser-aggregates (DAG node 14)")

	// ------------------------------------------------------------------
	// Stub: window functions → parser-window
	// ------------------------------------------------------------------
	case tokLAG:
		return nil, p.deferredFeature("LAG() window", "parser-window (DAG node 13)")
	case tokLEAD:
		return nil, p.deferredFeature("LEAD() window", "parser-window (DAG node 13)")
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("unexpected token %q in expression", p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// parseParenExpr handles tokPAREN_LEFT dispatch for primary expressions:
//
//	(expr)             → plain parenthesized expression, returns inner expr
//	(SELECT ...)       → STUB via parseSelectExpr (Task 9 adds stub)
//	(expr, expr, ...)  → STUB: valueList deferred to parser-dml (no AST node)
//	(expr MATCH ...)   → STUB: graph match deferred to parser-graph (node 16)
//
// At this task, SELECT is not yet stubbed (Task 9 does that). If the
// caller writes `(SELECT ...)`, the parser will try to parse SELECT as
// an expression and fall through to the default "unexpected token" error
// in parsePrimaryBase. Task 9 tightens this.
//
// Grammar: exprTerm#ExprTermWrappedQuery (line 543) + exprGraphMatchMany
// (lines 625-626).
func (p *Parser) parseParenExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume (

	// Empty parens `()` are not valid in any PartiQL primary position.
	if p.cur.Type == tokPAREN_RIGHT {
		return nil, &ParseError{
			Message: "empty parenthesized expression",
			Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
		}
	}

	first, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Graph match: (expr MATCH ...) — deferred to parser-graph.
	if p.cur.Type == tokMATCH {
		return nil, p.deferredFeature("graph MATCH expression", "parser-graph (DAG node 16)")
	}

	// valueList: (expr, expr, ...) — deferred to parser-dml.
	if p.cur.Type == tokCOMMA {
		return nil, p.deferredFeature("valueList", "parser-dml (DAG node 6)")
	}

	// Plain (expr): consume the closing paren and return the inner expr.
	// Note: the returned expression does NOT get a new wrapping node —
	// parentheses are purely syntactic. The inner expression's Loc is
	// preserved as-is (we don't extend it to cover the outer parens).
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, err
	}
	return first, nil
}

// parseArrayLit parses `[expr, expr, ...]`. Empty brackets `[]` are
// valid (empty array literal).
//
// Grammar: array (line 649-650).
func (p *Parser) parseArrayLit() (*ast.ListLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume [
	var items []ast.ExprNode
	if p.cur.Type != tokBRACKET_RIGHT {
		for {
			item, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACKET_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.ListLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseBagLit parses `<<expr, expr, ...>>`. Empty `<<>>` is valid.
//
// Grammar: bag (line 652-653).
func (p *Parser) parseBagLit() (*ast.BagLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume <<
	var items []ast.ExprNode
	if p.cur.Type != tokANGLE_DOUBLE_RIGHT {
		for {
			item, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokANGLE_DOUBLE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.BagLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTupleLit parses `{key: value, key: value, ...}`. Empty `{}` is
// valid.
//
// Grammar: tuple (line 655-656) + pair (line 658-659).
func (p *Parser) parseTupleLit() (*ast.TupleLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume {
	var pairs []*ast.TuplePair
	if p.cur.Type != tokBRACE_RIGHT {
		for {
			pair, err := p.parseTuplePair()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, pair)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.TupleLit{
		Pairs: pairs,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTuplePair parses one `key: value` entry inside a tuple literal.
func (p *Parser) parseTuplePair() (*ast.TuplePair, error) {
	key, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}
	value, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &ast.TuplePair{
		Key:   key,
		Value: value,
		Loc:   ast.Loc{Start: key.GetLoc().Start, End: value.GetLoc().End},
	}, nil
}

// parseParamRef consumes a QUESTION_MARK and returns an ast.ParamRef.
//
// Grammar: parameter (line 632-633).
func (p *Parser) parseParamRef() (*ast.ParamRef, error) {
	loc := p.cur.Loc
	p.advance()
	return &ast.ParamRef{Loc: loc}, nil
}
