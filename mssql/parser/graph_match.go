// Package parser - graph_match.go implements the graph MATCH predicate.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// The words ScriptDOM reads as identifiers and matches by their text
// (CodeGenerationSupporter): MATCH, SHORTEST_PATH and LAST_NODE are not
// registered keywords, and PATH is one only after FOR.
var (
	graphMatchWord        = newOptionSet().withIdents("MATCH")
	graphShortestPathWord = newOptionSet().withIdents("SHORTEST_PATH")
	graphLastNodeWord     = newOptionSet().withIdents("LAST_NODE")
	graphForPathWord      = newOptionSet().withIdents("PATH")
)

// parseGraphMatchPredicate parses MATCH ( <pattern> ), the current token being
// the parenthesis after MATCH, which the caller consumed at loc.
// SQL Server takes it only in a search condition: WHERE, HAVING, ON, a WHEN
// of CASE, and the condition of IF and WHILE. Elsewhere MATCH is not a
// function, so SELECT MATCH(a -(e)-> b) is an error (Msg 195).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/match-sql-graph
//
//	booleanExpressionPrimary: MATCH ( graphMatchPattern )
func (p *Parser) parseGraphMatchPredicate(loc int) (nodes.ExprNode, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	pattern, err := p.parseGraphMatchPattern()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.GraphMatchPredicate{
		Expression: pattern,
		Loc:        nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseGraphMatchPattern parses one or more graph match expressions, joined
// by AND.
//
//	graphMatchPattern: graphMatchExpressionPattern ( AND graphMatchExpressionPattern )*
func (p *Parser) parseGraphMatchPattern() (nodes.ExprNode, error) {
	left, err := p.parseGraphMatchExpressionPattern()
	if err != nil {
		return nil, err
	}
	for {
		andLoc := p.pos()
		if _, ok := p.match(kwAND); !ok {
			break
		}
		right, err := p.parseGraphMatchExpressionPattern()
		if err != nil {
			return nil, err
		}
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, andLoc), End: p.prevEnd()},
		}
	}
	return left, nil
}

// parseGraphMatchExpressionPattern parses one graph match expression.
//
//	graphMatchExpressionPattern:
//	      SHORTEST_PATH ( ... )          graphMatchRecursivePattern
//	    | LAST_NODE ...                  graphMatchLastNodePattern
//	    | <node> <step> ...              graphCompositeMatchPattern
//
// ScriptDOM tells them apart by the next word (CodeGenerationSupporter
// .ShortestPath and .LastNode), so a name that spells one of them commits to
// that pattern.
func (p *Parser) parseGraphMatchExpressionPattern() (nodes.ExprNode, error) {
	if p.graphWordAhead(graphShortestPathWord) {
		return p.parseGraphMatchRecursivePattern()
	}
	if p.graphWordAhead(graphLastNodeWord) {
		return p.parseGraphMatchLastNodePattern()
	}
	return p.parseGraphCompositeMatchPattern()
}

// graphWordAhead reports whether the current token is the word word,
// undelimited.
func (p *Parser) graphWordAhead(word optionSet) bool {
	return !p.isDelimited(p.cur) && p.isValidOption(word)
}

// parseGraphMatchLastNodePattern parses a pattern that starts with LAST_NODE:
// LAST_NODE (a) = LAST_NODE (b), or LAST_NODE (a) followed by steps.
//
//	graphMatchLastNodePattern:
//	    <node> = <node>                 (both sides must use LAST_NODE)
//	  | <node> <step> ( <step> )*
func (p *Parser) parseGraphMatchLastNodePattern() (nodes.ExprNode, error) {
	loc := p.pos()
	left, err := p.parseGraphMatchNodeExpression(true)
	if err != nil {
		return nil, err
	}
	if _, ok := p.match('='); ok {
		right, err := p.parseGraphMatchNodeExpression(true)
		if err != nil {
			return nil, err
		}
		if !left.UsesLastNode || !right.UsesLastNode {
			return nil, p.newParseError(left.Loc.Start, "both sides of a graph match equality must be LAST_NODE ( <node> )")
		}
		return &nodes.GraphMatchLastNodePredicate{
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	return p.parseGraphMatchSteps(loc, left, true)
}

// parseGraphCompositeMatchPattern parses a pattern that starts with a node
// and continues over one or more edges.
//
//	graphCompositeMatchPattern: <node> <step> ( <step> )*
func (p *Parser) parseGraphCompositeMatchPattern() (nodes.ExprNode, error) {
	loc := p.pos()
	left, err := p.parseGraphMatchNodeExpression(true)
	if err != nil {
		return nil, err
	}
	return p.parseGraphMatchSteps(loc, left, true)
}

// parseGraphMatchSteps parses one or more steps from left, joining them with
// AND, and reports the pattern it read. A step's right node continues the
// next step.
func (p *Parser) parseGraphMatchSteps(loc int, left *nodes.GraphMatchNodeExpr, canUseLastNode bool) (nodes.ExprNode, error) {
	var result nodes.ExprNode
	for {
		step, err := p.parseGraphSingleMatchPattern(left, canUseLastNode)
		if err != nil {
			return nil, err
		}
		if result == nil {
			result = step
		} else {
			result = &nodes.BinaryExpr{
				Op:    nodes.BinOpAnd,
				Left:  result,
				Right: step,
				Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
			}
		}
		if !p.atGraphStepStart() {
			break
		}
		left = step.Right
	}
	return result, nil
}

// atGraphStepStart reports whether the current tokens begin a step of a graph
// pattern: - ( <edge> ) -> or < - ( <edge> ) -.
func (p *Parser) atGraphStepStart() bool {
	if p.cur.Type == '-' {
		return p.peekNext().Type == '('
	}
	if p.cur.Type == '<' {
		return p.peekNext().Type == '-' && p.peekSecond().Type == '('
	}
	return false
}

// parseGraphMatchNodeExpression parses a node of a graph pattern.
//
//	graphMatchNodeExpression: LAST_NODE ( <identifier> ) | <identifier>
//
// LAST_NODE is refused where ScriptDOM's vCanUseLastNode is false, which is
// inside SHORTEST_PATH: the node the recursion ended at is not known there.
func (p *Parser) parseGraphMatchNodeExpression(canUseLastNode bool) (*nodes.GraphMatchNodeExpr, error) {
	loc := p.pos()
	if p.graphWordAhead(graphLastNodeWord) && p.peekNext().Type == '(' {
		if !canUseLastNode {
			return nil, p.newParseError(loc, "LAST_NODE cannot name a node inside SHORTEST_PATH")
		}
		p.advance() // consume LAST_NODE
		p.advance() // consume (
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.GraphMatchNodeExpr{
			Node:         name,
			UsesLastNode: true,
			Loc:          nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	return &nodes.GraphMatchNodeExpr{
		Node: name,
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseGraphSingleMatchPattern parses one step of a graph pattern, from left.
//
//	graphSingleMatchPattern:
//	    - ( <edge> ) -> <node>
//	  | < - ( <edge> ) - <node>
func (p *Parser) parseGraphSingleMatchPattern(left *nodes.GraphMatchNodeExpr, canUseLastNode bool) (*nodes.GraphMatchCompositeExpr, error) {
	loc := p.pos()
	arrowOnRight, edge, err := p.parseGraphEdge()
	if err != nil {
		return nil, err
	}
	right, err := p.parseGraphMatchNodeExpression(canUseLastNode)
	if err != nil {
		return nil, err
	}
	return &nodes.GraphMatchCompositeExpr{
		Left:         left,
		Edge:         edge,
		Right:        right,
		ArrowOnRight: arrowOnRight,
		Loc:          nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
	}, nil
}

// parseGraphSingleLeftMatchPattern parses one step of a left-anchored
// SHORTEST_PATH pattern, which has no right node.
//
//	graphSingleLeftMatchPattern:
//	    <node> - ( <edge> ) ->
//	  | <node> < - ( <edge> ) -
func (p *Parser) parseGraphSingleLeftMatchPattern(canUseLastNode bool) (*nodes.GraphMatchCompositeExpr, error) {
	loc := p.pos()
	left, err := p.parseGraphMatchNodeExpression(canUseLastNode)
	if err != nil {
		return nil, err
	}
	arrowOnRight, edge, err := p.parseGraphEdge()
	if err != nil {
		return nil, err
	}
	return &nodes.GraphMatchCompositeExpr{
		Left:         left,
		Edge:         edge,
		ArrowOnRight: arrowOnRight,
		Loc:          nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
	}, nil
}

// parseGraphEdge parses the edge of a step: - ( <edge> ) -> or
// < - ( <edge> ) -, and reports its alias and whether the arrow points right.
func (p *Parser) parseGraphEdge() (arrowOnRight bool, edge string, err error) {
	switch p.cur.Type {
	case '-':
		p.advance()
		arrowOnRight = true
	case '<':
		p.advance()
		arrowOnRight = false
		if _, err := p.expect('-'); err != nil {
			return false, "", err
		}
	default:
		return false, "", p.unexpectedToken()
	}
	if _, err := p.expect('('); err != nil {
		return false, "", err
	}
	edge, ok := p.parseIdentifier()
	if !ok {
		return false, "", p.unexpectedToken()
	}
	if _, err := p.expect(')'); err != nil {
		return false, "", err
	}
	if _, err := p.expect('-'); err != nil {
		return false, "", err
	}
	if arrowOnRight {
		if _, err := p.expect('>'); err != nil {
			return false, "", err
		}
	}
	return arrowOnRight, edge, nil
}

// parseGraphMatchRecursivePattern parses SHORTEST_PATH ( ... ).
//
//	graphMatchRecursivePattern:
//	    SHORTEST_PATH ( <node> ( <pattern> ) <quantifier> )
//	  | SHORTEST_PATH ( ( <pattern> ) <quantifier> <node> )
func (p *Parser) parseGraphMatchRecursivePattern() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume SHORTEST_PATH
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	pred := &nodes.GraphMatchRecursivePredicate{Loc: nodes.Loc{Start: loc, End: -1}}
	if p.cur.Type == '(' {
		// SHORTEST_PATH ( ( <pattern> ) <quantifier> <node> )
		p.advance() // consume (
		steps, err := p.parseGraphLeftRecursivePatterns()
		if err != nil {
			return nil, err
		}
		if len(steps) == 0 {
			return nil, p.newParseError(p.pos(), "SHORTEST_PATH takes at least one pattern")
		}
		pred.Expression = &nodes.List{Items: steps}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if err := p.parseGraphRecursiveQuantifier(pred); err != nil {
			return nil, err
		}
		outer, err := p.parseGraphMatchNodeExpression(true)
		if err != nil {
			return nil, err
		}
		pred.Outer = outer
	} else {
		// SHORTEST_PATH ( <node> ( <pattern> ) <quantifier> )
		outer, err := p.parseGraphMatchNodeExpression(true)
		if err != nil {
			return nil, err
		}
		pred.Outer = outer
		pred.AnchorOnLeft = true
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		steps, err := p.parseGraphRightRecursivePatterns(outer)
		if err != nil {
			return nil, err
		}
		if len(steps) == 0 {
			return nil, p.newParseError(p.pos(), "SHORTEST_PATH takes at least one pattern")
		}
		pred.Expression = &nodes.List{Items: steps}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if err := p.parseGraphRecursiveQuantifier(pred); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	pred.Loc.End = p.prevEnd()
	return pred, nil
}

// parseGraphRightRecursivePatterns parses the steps of a right-anchored
// SHORTEST_PATH, which run from outer.
//
//	graphMatchSPRightRecursivePattern: ( graphSingleMatchPattern )*
func (p *Parser) parseGraphRightRecursivePatterns(outer *nodes.GraphMatchNodeExpr) ([]nodes.Node, error) {
	var steps []nodes.Node
	left := outer
	for p.atGraphStepStart() {
		step, err := p.parseGraphSingleMatchPattern(left, false)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
		left = step.Right
	}
	return steps, nil
}

// parseGraphLeftRecursivePatterns parses the steps of a left-anchored
// SHORTEST_PATH, which have no right node.
//
//	graphMatchSPLeftRecursivePattern: ( graphSingleLeftMatchPattern )*
func (p *Parser) parseGraphLeftRecursivePatterns() ([]nodes.Node, error) {
	var steps []nodes.Node
	for p.isGraphNodeStart() {
		step, err := p.parseGraphSingleLeftMatchPattern(false)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

// isGraphNodeStart reports whether the current token begins a node of a
// graph pattern: an identifier, or LAST_NODE.
func (p *Parser) isGraphNodeStart() bool {
	if p.graphWordAhead(graphLastNodeWord) && p.peekNext().Type == '(' {
		return true
	}
	return p.isIdentLike() && !p.isDelimited(p.cur)
}

// parseGraphRecursiveQuantifier parses the quantifier of SHORTEST_PATH: + or
// { <n> , <m> }.
func (p *Parser) parseGraphRecursiveQuantifier(pred *nodes.GraphMatchRecursivePredicate) error {
	switch p.cur.Type {
	case '+':
		p.advance()
		pred.IsPlusSign = true
		return nil
	case '{':
		p.advance()
		lower, err := p.parseGraphLimit()
		if err != nil {
			return err
		}
		if _, err := p.expect(','); err != nil {
			return err
		}
		upper, err := p.parseGraphLimit()
		if err != nil {
			return err
		}
		if _, err := p.expect('}'); err != nil {
			return err
		}
		pred.Lower, pred.Upper = lower, upper
		return nil
	}
	return p.newParseError(p.pos(), "SHORTEST_PATH takes + or { n , m } after its pattern")
}

// parseGraphLimit parses one limit of a { n , m } quantifier, which the
// grammar takes as an integer.
func (p *Parser) parseGraphLimit() (*nodes.Literal, error) {
	if procIsBinary(p.cur) || p.cur.Type != tokICONST {
		return nil, p.unexpectedToken()
	}
	tok := p.advance()
	return &nodes.Literal{
		Type: nodes.LitInteger,
		Ival: tok.Ival,
		Str:  tok.Str,
		Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
	}, nil
}
