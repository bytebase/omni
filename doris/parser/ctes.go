package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// WITH / CTE parsing (T1.6)
// ---------------------------------------------------------------------------

// parseWithSelect parses a WITH clause followed by the query it prefixes:
//
//	WITH cte_name [(col, ...)] AS (query)
//	  [, cte_name2 [(col, ...)] AS (query2)]
//	{ SELECT ... | VALUES ... }
//
// (DorisParser.g4:1167 query is `cte? queryTerm queryOrganization`, and a
// VALUES query is a queryTerm alternative, so a WITH clause prefixes either.)
// The WITH keyword has NOT yet been consumed when this is called.
func (p *Parser) parseWithSelect() (ast.Node, error) {
	with, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}
	return p.parseWithQueryBody(with)
}

// parseWithQueryBody parses the query a WITH clause prefixes and attaches the
// clause to it, extending the node's range to cover WITH. The caller reads any
// set-operation and trailing-clause tail.
func (p *Parser) parseWithQueryBody(with *ast.WithClause) (ast.Node, error) {
	if p.cur.Kind == kwVALUES {
		stmt, err := p.parseValuesStmt()
		if err != nil {
			return nil, err
		}
		stmt.With, stmt.Loc.Start = with, with.Loc.Start
		return stmt, nil
	}

	stmt, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.With, stmt.Loc.Start = with, with.Loc.Start
	return stmt, nil
}

// parseWithStatement parses a statement that begins with a WITH clause. Doris
// prefixes a query (SELECT, and a set operation after it) or a DML statement
// with it:
//
//	explain? cte? INSERT ...   (DorisParser.g4:138)
//	explain? cte? UPDATE ...   (DorisParser.g4:145)
//	explain? cte? DELETE ...   (DorisParser.g4:149)
//
// The WITH keyword has NOT yet been consumed when this is called.
func (p *Parser) parseWithStatement() (ast.Node, error) {
	with, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	switch p.cur.Kind {
	case kwINSERT:
		stmt, err := p.parseInsert()
		if err != nil {
			return nil, err
		}
		stmt.With = with
		stmt.Loc.Start = with.Loc.Start
		return stmt, nil
	case kwUPDATE:
		updateTok := p.advance() // consume UPDATE
		stmt, err := p.parseUpdateStmt(updateTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		stmt.With = with
		stmt.Loc.Start = with.Loc.Start
		return stmt, nil
	case kwDELETE:
		deleteTok := p.advance() // consume DELETE
		stmt, err := p.parseDeleteStmt(deleteTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		stmt.With = with
		stmt.Loc.Start = with.Loc.Start
		return stmt, nil
	}

	stmt, err := p.parseWithQueryBody(with)
	if err != nil {
		return nil, err
	}

	// WITH ... {SELECT|VALUES} continues with set operators and their
	// trailing clauses as a bare query does (engine-verified).
	return p.parseQueryStmt(stmt)
}

// parseWithClause parses WITH cte (',' cte)* and returns the
// clause, leaving the cursor at the statement keyword that follows it. The
// WITH keyword has NOT yet been consumed when this is called.
func (p *Parser) parseWithClause() (*ast.WithClause, error) {
	withTok, err := p.expect(kwWITH)
	if err != nil {
		return nil, err
	}

	with := &ast.WithClause{
		Loc: ast.Loc{Start: withTok.Loc.Start},
	}

	// No RECURSIVE: DorisParser.g4's cte is WITH aliasQuery (COMMA
	// aliasQuery)* at the pinned 3.1.4, and RECURSIVE is no token of its
	// lexer, so WITH recursive AS (...) names a CTE recursive.

	// Parse comma-separated CTE definitions.
	cte, err := p.parseCTE()
	if err != nil {
		return nil, err
	}
	with.CTEs = append(with.CTEs, cte)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		cte, err = p.parseCTE()
		if err != nil {
			return nil, err
		}
		with.CTEs = append(with.CTEs, cte)
	}

	with.Loc.End = p.prev.Loc.End
	return with, nil
}

// parseCTE parses one CTE definition:
//
//	cte_name [(col1, col2, ...)] AS ( query )
func (p *Parser) parseCTE() (*ast.CTE, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	startLoc := p.cur.Loc

	// CTE name — must be a plain or quoted identifier.
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	cte := &ast.CTE{
		Name: name,
		Loc:  ast.Loc{Start: startLoc.Start},
	}

	// Optional column list: (col1, col2, ...)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('

		colName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cte.Columns = append(cte.Columns, colName)

		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			colName, _, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			cte.Columns = append(cte.Columns, colName)
		}

		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	// AS keyword.
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Opening paren of the inner query.
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Inner SELECT statement. If it starts with WITH, recurse; otherwise parse
	// a plain SELECT followed by optional set-operator tail (UNION/INTERSECT/EXCEPT).
	var innerStmt ast.Node
	switch p.cur.Kind {
	case kwWITH:
		innerStmt, err = p.parseWithSelect()
	case kwVALUES:
		// WITH c AS (VALUES (1)) SELECT 1 — a CTE body is a query, and a
		// VALUES query is one (DorisParser.g4:1167, :1186).
		innerStmt, err = p.parseValuesStmt()
	case int('('):
		// WITH c AS ((SELECT 1)) ... is engine-verified valid.
		innerStmt, err = p.parseParenQueryOperand()
	default:
		innerStmt, err = p.parseSelectStmt()
	}
	if err != nil {
		return nil, err
	}
	// Either operand form continues with set operators and, after one,
	// trailing clauses — WITH c AS (SELECT 1 UNION (SELECT 2) LIMIT 5) ...
	innerStmt, err = p.parseQueryTail(innerStmt)
	if err != nil {
		return nil, err
	}
	cte.Query = innerStmt

	// Closing paren.
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	cte.Loc.End = p.prev.Loc.End
	return cte, nil
}
