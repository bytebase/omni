package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `parser-select` DAG node (with select.go, setops.go,
// relation.go and ctes.go): it implements Trino's named-window grammar — the
// `WINDOW` clause of a querySpecification and its `windowDefinition` rule.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	querySpecification: SELECT … (WINDOW windowDefinition (, windowDefinition)*)? ;
//	windowDefinition  : name=identifier AS ( windowSpecification ) ;
//	windowSpecification: (existingWindowName)? (PARTITION BY …)?
//	                       (ORDER BY sortItem (, …)*)? windowFrame? ;
//
// The windowSpecification itself — PARTITION BY / ORDER BY / the window frame
// (ROWS|RANGE|GROUPS …) — was implemented by the expressions DAG node for the
// inline OVER clause (function.go: WindowSpec / WindowFrame / WindowBound and
// parseWindowSpecification). This file reuses that parser for the named-window
// definitions, so a name defined in the WINDOW clause and a name referenced by
// `OVER w` share one WindowSpec model. The pattern-recognition frame additions
// remain deferred to parser-match-recognize (B2 in expr.go).

// WindowDefinition is one `name AS ( windowSpecification )` of a querySpecification
// WINDOW clause (the windowDefinition rule). Name is the window's name; Spec is
// the (reused) inline-OVER window specification.
type WindowDefinition struct {
	Name *ast.Identifier
	Spec *WindowSpec
	Loc  ast.Loc
}

// parseWindowClause parses `WINDOW windowDefinition (, windowDefinition)*`
// (WINDOW is current) and returns the non-empty list of named window definitions.
func (p *Parser) parseWindowClause() ([]WindowDefinition, error) {
	p.advance() // consume WINDOW
	first, err := p.parseWindowDefinition()
	if err != nil {
		return nil, err
	}
	defs := []WindowDefinition{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseWindowDefinition()
		if err != nil {
			return nil, err
		}
		defs = append(defs, next)
	}
	return defs, nil
}

// parseWindowDefinition parses one `name=identifier AS ( windowSpecification )`
// (the windowDefinition rule). The window specification body is parsed by the
// shared parseWindowSpecification (function.go), which consumes the opening '('
// and the closing ')'.
func (p *Parser) parseWindowDefinition() (WindowDefinition, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return WindowDefinition{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return WindowDefinition{}, err
	}
	if p.cur.Kind != int('(') {
		return WindowDefinition{}, p.exprErrorAt("expected ( after AS in window definition")
	}
	// parseWindowSpecification consumes the '(' … ')'. Pass the name's start as
	// the span origin so the definition's WindowSpec.Loc covers `name AS ( … )`.
	spec, err := p.parseWindowSpecification(name.Loc.Start)
	if err != nil {
		return WindowDefinition{}, err
	}
	return WindowDefinition{
		Name: name,
		Spec: spec,
		Loc:  ast.Loc{Start: name.Loc.Start, End: spec.Loc.End},
	}, nil
}
