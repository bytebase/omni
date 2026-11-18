package parser

import "github.com/bytebase/omni/trino/ast"

// This file implements rootQueryWithSession, the `WITH SESSION name = value …
// <query>` statement prefix (statementDefault). It is only a statement prefix:
// a subquery's WITH is the CTE form.

// SessionProperty is one `qualifiedName = expression` of a WITH SESSION prefix.
type SessionProperty struct {
	Name  *ast.QualifiedName
	Value Expr
	Loc   ast.Loc
}

// parseQueryWithSessionStmt parses `WITH SESSION prop (, prop)* query`. WITH is
// current and SESSION is the next token. The result is an ordinary QueryStmt
// carrying the session properties, so every analysis consumer keeps working.
func (p *Parser) parseQueryWithSessionStmt() (ast.Node, error) {
	withTok := p.advance() // consume WITH
	p.advance()            // consume SESSION
	s := &QueryStmt{Loc: ast.Loc{Start: withTok.Loc.Start}}
	first, err := p.parseSessionProperty()
	if err != nil {
		return nil, err
	}
	s.Session = append(s.Session, first)
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSessionProperty()
		if err != nil {
			return nil, err
		}
		s.Session = append(s.Session, next)
	}
	q, err := p.parseRootQuery()
	if err != nil {
		return nil, err
	}
	s.Query = q
	s.Loc.End = q.Loc.End
	return s, nil
}

// parseSessionProperty parses `qualifiedName = expression`.
func (p *Parser) parseSessionProperty() (SessionProperty, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return SessionProperty{}, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return SessionProperty{}, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return SessionProperty{}, err
	}
	return SessionProperty{
		Name:  name,
		Value: value,
		Loc:   ast.Loc{Start: name.Loc.Start, End: value.Span().End},
	}, nil
}
