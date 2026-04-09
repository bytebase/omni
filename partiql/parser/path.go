package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// isPathStepStart reports whether the current token can begin a path
// step. Path steps appear after a primary expression in the form:
//
//	.field  ."Field"  .*  [expr]  [*]
//
// Grammar: pathStep (PartiQLParser.g4 lines 618-623).
func isPathStepStart(t int) bool {
	return t == tokPERIOD || t == tokBRACKET_LEFT
}

// parsePathSteps consumes a sequence of path steps and wraps the base
// expression in an ast.PathExpr. Called by parsePrimary when the
// current token after the base is a path-step start.
//
// This function assumes at least one path step is present (the caller
// checks isPathStepStart before invoking). If called with no path
// step to consume, returns the base unchanged wrapped in a PathExpr
// with an empty Steps slice — callers should avoid that degenerate
// case.
//
// Grammar: exprPrimary#ExprPrimaryPath (line 528) + pathStep (lines
// 618-623).
func (p *Parser) parsePathSteps(base ast.ExprNode) (*ast.PathExpr, error) {
	var steps []ast.PathStep
	for isPathStepStart(p.cur.Type) {
		step, err := p.parsePathStep()
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	endLoc := base.GetLoc().End
	if len(steps) > 0 {
		endLoc = steps[len(steps)-1].GetLoc().End
	}
	return &ast.PathExpr{
		Root:  base,
		Steps: steps,
		Loc:   ast.Loc{Start: base.GetLoc().Start, End: endLoc},
	}, nil
}

// parsePathStep consumes exactly one path step.
//
// Grammar: pathStep (lines 618-623):
//
//	BRACKET_LEFT key=expr BRACKET_RIGHT        # PathStepIndexExpr
//	BRACKET_LEFT ASTERISK BRACKET_RIGHT        # PathStepIndexAll
//	PERIOD key=symbolPrimitive                 # PathStepDotExpr
//	PERIOD ASTERISK                            # PathStepDotAll
func (p *Parser) parsePathStep() (ast.PathStep, error) {
	switch p.cur.Type {
	case tokPERIOD:
		start := p.cur.Loc.Start
		p.advance() // consume .
		// Dot-star: .*
		if p.cur.Type == tokASTERISK {
			end := p.cur.Loc.End
			p.advance()
			return &ast.AllFieldsStep{
				Loc: ast.Loc{Start: start, End: end},
			}, nil
		}
		// Dot-field: .foo or ."Foo"
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DotStep{
			Field:         name,
			CaseSensitive: caseSensitive,
			Loc:           ast.Loc{Start: start, End: nameLoc.End},
		}, nil

	case tokBRACKET_LEFT:
		start := p.cur.Loc.Start
		p.advance() // consume [
		// Bracket-star: [*]
		if p.cur.Type == tokASTERISK && p.peekNext().Type == tokBRACKET_RIGHT {
			p.advance() // consume *
			end := p.cur.Loc.End
			p.advance() // consume ]
			return &ast.WildcardStep{
				Loc: ast.Loc{Start: start, End: end},
			}, nil
		}
		// Bracket-expr: [expr]
		idx, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(tokBRACKET_RIGHT)
		if err != nil {
			return nil, err
		}
		return &ast.IndexStep{
			Index: idx,
			Loc:   ast.Loc{Start: start, End: rp.Loc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: "expected path step (. or [)",
		Loc:     p.cur.Loc,
	}
}
