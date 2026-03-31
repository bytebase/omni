package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseFromClause parses: FROM from_specification
// Returns (TableExpr, []*JoinExpr, error).
func (p *Parser) parseFromClause() (nodes.TableExpr, []*nodes.JoinExpr, error) {
	p.advance() // consume FROM

	source, err := p.parseFromSource()
	if err != nil {
		return nil, nil, err
	}

	// Parse JOIN clauses.
	var joins []*nodes.JoinExpr
	for p.cur.Type == tokJOIN {
		j, err := p.parseJoinClause()
		if err != nil {
			return nil, nil, err
		}
		joins = append(joins, j)
	}

	return source, joins, nil
}

// parseJoinClause parses: JOIN from_source.
func (p *Parser) parseJoinClause() (*nodes.JoinExpr, error) {
	startLoc := p.pos()
	p.advance() // consume JOIN

	source, err := p.parseFromSource()
	if err != nil {
		return nil, err
	}

	endLoc := p.prev.End
	return &nodes.JoinExpr{
		Source: source,
		Loc:    nodes.Loc{Start: startLoc, End: endLoc},
	}, nil
}

// parseFromSource parses either:
//   - container_expr [AS? alias]
//   - identifier IN container_expr  (array iteration)
func (p *Parser) parseFromSource() (nodes.TableExpr, error) {
	startLoc := p.pos()

	// Check for iteration form: identifier IN container_expr
	if isIdentLike(p.cur.Type) && p.peekNext().Type == tokIN {
		alias := p.cur.Str
		p.advance() // consume identifier
		p.advance() // consume IN
		container, err := p.parseContainerExpr()
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.End
		return &nodes.ArrayIterationExpr{
			Alias:  alias,
			Source: container,
			Loc:    nodes.Loc{Start: startLoc, End: endLoc},
		}, nil
	}

	// Otherwise: container_expr [AS? alias]
	container, err := p.parseContainerExpr()
	if err != nil {
		return nil, err
	}

	// Optional alias.
	if p.cur.Type == tokAS {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.End
		return &nodes.AliasedTableExpr{
			Source: container,
			Alias:  name,
			Loc:    nodes.Loc{Start: startLoc, End: endLoc},
		}, nil
	}
	if isIdentLike(p.cur.Type) && !isClauseStart(p.cur.Type) {
		name := p.cur.Str
		p.advance()
		endLoc := p.prev.End
		return &nodes.AliasedTableExpr{
			Source: container,
			Alias:  name,
			Loc:    nodes.Loc{Start: startLoc, End: endLoc},
		}, nil
	}

	return container, nil
}

// parseContainerExpr parses: ROOT | container_name | expr.property | expr[index] | (SELECT ...)
func (p *Parser) parseContainerExpr() (nodes.TableExpr, error) {
	startLoc := p.pos()
	var base nodes.TableExpr

	switch p.cur.Type {
	case tokROOT:
		p.advance()
		base = &nodes.ContainerRef{
			Root: true,
			Loc:  nodes.Loc{Start: startLoc, End: p.prev.End},
		}
	case tokLPAREN:
		// Subquery: (SELECT ...)
		p.advance() // consume (
		sub, err := p.parseSelect()
		if err != nil {
			return nil, err
		}
		rparen, err := p.expect(tokRPAREN)
		if err != nil {
			return nil, err
		}
		base = &nodes.SubqueryExpr{
			Select: sub,
			Loc:    nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
		}
	default:
		if !isIdentLike(p.cur.Type) {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected container name, got %q", p.cur.Str),
				Pos:     p.cur.Loc,
			}
		}
		name := p.cur.Str
		p.advance()
		base = &nodes.ContainerRef{
			Name: name,
			Loc:  nodes.Loc{Start: startLoc, End: p.prev.End},
		}
	}

	// Postfix: . and []
	for {
		if p.cur.Type == tokDOT {
			p.advance()
			prop, _, err := p.parsePropertyName()
			if err != nil {
				return nil, err
			}
			endLoc := p.prev.End
			base = &nodes.DotAccessExpr{
				Expr:     exprFromTableExpr(base),
				Property: prop,
				Loc:      nodes.Loc{Start: tableExprStart(base), End: endLoc},
			}
		} else if p.cur.Type == tokLBRACK {
			p.advance()
			index, err := p.parseBracketIndex()
			if err != nil {
				return nil, err
			}
			rbracket, err := p.expect(tokRBRACK)
			if err != nil {
				return nil, err
			}
			base = &nodes.BracketAccessExpr{
				Expr:  exprFromTableExpr(base),
				Index: index,
				Loc:   nodes.Loc{Start: tableExprStart(base), End: rbracket.Loc + 1},
			}
		} else {
			break
		}
	}

	return base, nil
}

// parseBracketIndex parses the index inside [...] in FROM context.
// The grammar's container_expression alt4 restricts this to string literals,
// array indexes (integers), and parameter names. The scalar_expression version
// (in expr.go) is more permissive and allows any expression.
func (p *Parser) parseBracketIndex() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokSCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.StringLit{Val: val, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case tokDCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.StringLit{Val: val, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case tokICONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.NumberLit{Val: val, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case tokPARAM:
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.ParamRef{Name: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected string, integer, or parameter in bracket index, got %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
}

// exprFromTableExpr converts a TableExpr to ExprNode.
func exprFromTableExpr(te nodes.TableExpr) nodes.ExprNode {
	switch t := te.(type) {
	case *nodes.ContainerRef:
		return &nodes.ColumnRef{Name: t.Name, Loc: t.Loc}
	case *nodes.DotAccessExpr:
		return t
	case *nodes.BracketAccessExpr:
		return t
	case *nodes.AliasedTableExpr:
		return exprFromTableExpr(t.Source)
	default:
		// Should not happen in valid FROM expressions.
		return &nodes.ColumnRef{Name: "?", Loc: nodes.Loc{Start: -1, End: -1}}
	}
}

// tableExprStart returns the start location of a TableExpr.
// NOTE: When adding new TableExpr types, add a case here.
func tableExprStart(te nodes.TableExpr) int {
	switch t := te.(type) {
	case *nodes.ContainerRef:
		return t.Loc.Start
	case *nodes.DotAccessExpr:
		return t.Loc.Start
	case *nodes.BracketAccessExpr:
		return t.Loc.Start
	case *nodes.AliasedTableExpr:
		return t.Loc.Start
	case *nodes.ArrayIterationExpr:
		return t.Loc.Start
	case *nodes.SubqueryExpr:
		return t.Loc.Start
	default:
		return -1
	}
}

// isClauseStart returns true if the token type starts a clause.
func isClauseStart(tokType int) bool {
	switch tokType {
	case tokFROM, tokWHERE, tokGROUP, tokHAVING, tokORDER, tokOFFSET, tokLIMIT, tokJOIN:
		return true
	}
	return false
}

// isIdentLike returns true if the token can serve as an identifier.
func isIdentLike(tokType int) bool {
	switch tokType {
	case tokIDENT,
		tokIN, tokBETWEEN, tokTOP, tokVALUE, tokORDER, tokBY,
		tokGROUP, tokOFFSET, tokLIMIT, tokASC, tokDESC, tokEXISTS,
		tokLIKE, tokHAVING, tokJOIN, tokESCAPE, tokARRAY, tokROOT, tokRANK:
		return true
	}
	return false
}
