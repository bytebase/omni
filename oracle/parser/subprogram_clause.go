package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// Clauses of subprogram and package headings.
//
// Each clause stands at most once in a heading: Oracle 23ai rejects a second
// AUTHID, ACCESSIBLE BY, DETERMINISTIC, PIPELINED, PARALLEL_ENABLE,
// RESULT_CACHE, SQL_MACRO, ORDER BY, or CLUSTER BY with PLS-00371.

// clauseSeen records the clauses a heading has used, for the PLS-00371 check.
type clauseSeen map[string]bool

// subprogramLevel says where a subprogram heading stands, which decides the
// clauses it takes. On Oracle 23ai AUTHID stands only on a schema-level unit
// (PLS-00157 on a packaged, nested, or object type method), and ACCESSIBLE
// BY not on an object type method (PLS-00262). DEFAULT COLLATION is
// documented for schema-level units only.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/DEFAULT-COLLATION-clause.html
type subprogramLevel int

const (
	subprogramSchema   subprogramLevel = iota // CREATE PROCEDURE or FUNCTION
	subprogramPackaged                        // in a package, or nested in a block
	subprogramMethod                          // an object type method
)

// firstClause reports a syntax error at the current token, which starts clause,
// when the heading already has that clause.
func (p *Parser) firstClause(seen clauseSeen, clause string) error {
	if seen[clause] {
		return p.syntaxErrorAtCur()
	}
	seen[clause] = true
	return nil
}

// parseAccessibleByClause parses ACCESSIBLE BY ( accessor [, accessor ]... ),
// the current token at ACCESSIBLE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/ACCESSIBLE-BY-clause.html
//
//	accessor  ::= [ unit_kind ] [ schema. ] unit_name
//	unit_kind ::= { FUNCTION | PROCEDURE | PACKAGE | TRIGGER | TYPE }
//
// Oracle 23ai reads an unquoted unit kind word as the kind, so ACCESSIBLE BY
// (PACKAGE) is PLS-00103, as are an empty list, a missing parenthesis, and an
// accessor with a database link.
func (p *Parser) parseAccessibleByClause() (*nodes.List, error) {
	p.advance() // consume ACCESSIBLE
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for {
		acc := &nodes.Accessor{Loc: nodes.Loc{Start: p.pos()}}
		switch p.cur.Type {
		case kwFUNCTION:
			acc.UnitKind = "FUNCTION"
		case kwPROCEDURE:
			acc.UnitKind = "PROCEDURE"
		case kwPACKAGE:
			acc.UnitKind = "PACKAGE"
		case kwTRIGGER:
			acc.UnitKind = "TRIGGER"
		case kwTYPE:
			acc.UnitKind = "TYPE"
		}
		if acc.UnitKind != "" {
			p.advance()
		}
		name, err := p.parseUnitName()
		if err != nil {
			return nil, err
		}
		acc.Name = name
		acc.Loc.End = p.prev.End
		list.Items = append(list.Items, acc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseUnitName parses [ schema. ] name, without a database link.
func (p *Parser) parseUnitName() (*nodes.ObjectName, error) {
	start := p.pos()
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	name := &nodes.ObjectName{Name: first, Loc: nodes.Loc{Start: start}}
	if p.cur.Type == '.' {
		p.advance()
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		second, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		name.Schema, name.Name = first, second
	}
	name.Loc.End = p.prev.End
	return name, nil
}

// parseDefaultCollationClause parses DEFAULT COLLATION collation_option, the
// current token at DEFAULT.
func (p *Parser) parseDefaultCollationClause() error {
	p.advance() // consume DEFAULT
	p.advance() // consume COLLATION
	if !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume the collation option (USING_NLS_COMP)
	return nil
}

// atDefaultCollation reports whether the current tokens are DEFAULT COLLATION.
func (p *Parser) atDefaultCollation() bool {
	return p.cur.Type == kwDEFAULT && p.isIdentLikeStrAt(p.peekNext(), "COLLATION")
}

// parseProcedureProperties parses the clauses between a procedure's
// parameters and its IS | AS or ';', in any order; level limits them as
// subprogramLevel says.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-PROCEDURE-statement.html
//
//	[ { invoker_rights_clause | accessible_by_clause | default_collation_clause }... ]
func (p *Parser) parseProcedureProperties(stmt *nodes.CreateProcedureStmt, level subprogramLevel) error {
	seen := clauseSeen{}
	for {
		switch {
		case p.isKeywordStr("AUTHID") && level == subprogramSchema:
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE") && level != subprogramMethod:
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation() && level == subprogramSchema:
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			if err := p.parseDefaultCollationClause(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// parsePackageProperties parses the clauses between a package
// specification's name and its IS | AS, in any order.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-PACKAGE-statement.html
//
//	[ { default_collation_clause | invoker_rights_clause | accessible_by_clause }... ]
func (p *Parser) parsePackageProperties(stmt *nodes.CreatePackageStmt) error {
	seen := clauseSeen{}
	for {
		switch {
		case p.isKeywordStr("AUTHID"):
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE"):
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation():
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			if err := p.parseDefaultCollationClause(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// parseParallelEnableSpec parses the partitioning of PARALLEL_ENABLE, the
// current token at '('.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/PARALLEL_ENABLE-clause.html
//
//	( PARTITION argument BY { ANY | { HASH | RANGE | VALUE } ( expr [, expr ]... ) } )
//
// Oracle 23ai rejects an empty partitioning, PARTITION without BY, other
// methods (LIST), and an empty column list with PLS-00103. That VALUE takes a
// single column is checked when the unit compiles.
func (p *Parser) parseParallelEnableSpec() (*nodes.ParallelEnableClause, error) {
	spec := &nodes.ParallelEnableClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // consume '('
	if p.cur.Type != kwPARTITION {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	arg, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	spec.Argument = arg
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.cur.Type == kwANY:
		spec.PartitionBy = "ANY"
		p.advance()
	case p.cur.Type == kwHASH, p.cur.Type == kwRANGE, p.isKeywordStr("VALUE"):
		spec.PartitionBy = p.cur.Str
		if p.cur.Type == kwHASH {
			spec.PartitionBy = "HASH"
		} else if p.cur.Type == kwRANGE {
			spec.PartitionBy = "RANGE"
		}
		p.advance()
		cols, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		spec.Columns = cols
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	spec.Loc.End = p.prev.End
	return spec, nil
}

// parseStreamingClause parses { ORDER | CLUSTER } argument BY ( expr [, expr ]... ),
// the current token at ORDER or CLUSTER. Oracle 23ai takes it anywhere in a
// function's clause list, with or without PARALLEL_ENABLE, and requires a
// plain parameter name (ORDER c.x BY is PLS-00103).
func (p *Parser) parseStreamingClause() (*nodes.StreamingClause, error) {
	clause := &nodes.StreamingClause{Kind: "ORDER", Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == kwCLUSTER {
		clause.Kind = "CLUSTER"
	}
	p.advance()
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	arg, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	clause.Argument = arg
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	cols, err := p.parseParenExprList()
	if err != nil {
		return nil, err
	}
	clause.Columns = cols
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseParenExprList parses ( expr [, expr ]... ) with at least one
// expression, the current token at '('.
func (p *Parser) parseParenExprList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		list.Items = append(list.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseResultCacheReliesOn parses RELIES_ON ( [ data_source [, data_source ]... ] ),
// the current token at RELIES_ON. The list may be empty on Oracle 23ai.
func (p *Parser) parseResultCacheReliesOn() (*nodes.List, error) {
	p.advance() // consume RELIES_ON
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for p.cur.Type != ')' {
		// A ',' needs a data source after it: RELIES_ON (t1,) is PLS-00103.
		name, err := p.parseUnitName()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == ')' {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseSqlMacroType parses ( [ TYPE => ] { SCALAR | TABLE } ) after SQL_MACRO,
// the current token at '('.
func (p *Parser) parseSqlMacroType() (string, error) {
	p.advance() // consume '('
	if p.cur.Type == kwTYPE && p.peekNext().Type == tokASSOC {
		p.advance() // consume TYPE
		p.advance() // consume =>
	}
	var kind string
	switch {
	case p.cur.Type == kwTABLE:
		kind = "TABLE"
	case p.isKeywordStr("SCALAR"):
		kind = "SCALAR"
	default:
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return kind, nil
}
