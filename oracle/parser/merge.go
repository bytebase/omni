package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseMergeStmt parses a MERGE INTO statement.
//
// BNF: oracle/parser/bnf/MERGE.bnf
//
//	MERGE [ hint ] INTO [ schema. ] { table | view } [ t_alias ]
//	    USING { [ schema. ] { table | view }
//	          | subquery
//	          } [ t_alias ]
//	    ON ( condition )
//	    [ merge_update_clause ]
//	    [ merge_insert_clause ]
//	    [ error_logging_clause ] ;
//
//	merge_update_clause::=
//	    WHEN MATCHED THEN
//	    UPDATE SET column = { expr | DEFAULT }
//	               [, column = { expr | DEFAULT } ]...
//	    [ WHERE condition ]
//	    [ DELETE WHERE condition ]
//
//	merge_insert_clause::=
//	    WHEN NOT MATCHED THEN
//	    INSERT [ ( column [, column ]... ) ]
//	    VALUES ( { expr | DEFAULT } [, { expr | DEFAULT } ]... )
//	    [ WHERE condition ]
//
//	error_logging_clause::=
//	    LOG ERRORS [ INTO [ schema. ] table ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
func (p *Parser) parseMergeStmt() (*nodes.MergeStmt, error) {
	start := p.pos()
	p.advance() // consume MERGE

	stmt := &nodes.MergeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
		})
		p.advance()
	}

	// INTO
	if p.cur.Type == kwINTO {
		p.advance()
	}
	var parseErr801 error

	// Target table
	stmt.Target, parseErr801 = p.parseObjectName()
	if parseErr801 !=

		// Optional target alias
		nil {
		return nil, parseErr801
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr802 error
		stmt.TargetAlias, parseErr802 = p.parseAlias()
		if parseErr802 != nil {
			return nil, parseErr802
		}
	} else if p.isTableAliasCandidate() {
		var parseErr803 error
		stmt.TargetAlias, parseErr803 = p.parseAlias()
		if parseErr803 !=

			// USING
			nil {
			return nil, parseErr803
		}
	}

	if p.cur.Type == kwUSING {
		p.advance()
	}
	var parseErr804 error

	// Source table or subquery
	stmt.Source, parseErr804 = p.parseTableRef()
	if parseErr804 !=

		// Optional source alias — parseTableRef already handles alias on TableRef,
		// but for MergeStmt we want it in SourceAlias. Extract it.
		nil {
		return nil, parseErr804
	}

	if tr, ok := stmt.Source.(*nodes.TableRef); ok && tr.Alias != nil {
		stmt.SourceAlias = tr.Alias
		tr.Alias = nil
	}

	// ON ( condition )
	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	// The condition may or may not be wrapped in parens.
	wrappedOn := false
	if p.cur.Type == '(' {
		wrappedOn = true
		p.advance()
	}
	var parseErr805 error
	stmt.On, parseErr805 = p.parseExpr()
	if parseErr805 != nil {
		return nil, parseErr805
	}
	if stmt.On == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if wrappedOn {
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	stmt.Clauses = &nodes.List{}
	for p.cur.Type == kwWHEN {
		mc, parseErr807 := p.parseMergeClause()
		if parseErr807 != nil {
			return nil, parseErr807
		}
		if mc != nil {
			stmt.Clauses.Items = append(stmt.Clauses.Items, mc)
		}
	}
	if stmt.Clauses.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	// error_logging_clause
	if p.cur.Type == kwLOG {
		var parseErr808 error
		stmt.ErrorLog, parseErr808 = p.parseErrorLogClause()
		if parseErr808 != nil {
			return nil, parseErr808
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMergeClause parses a single WHEN MATCHED / WHEN NOT MATCHED clause.
func (p *Parser) parseMergeClause() (*nodes.MergeClause, error) {
	start := p.pos()
	p.advance() // consume WHEN

	mc := &nodes.MergeClause{
		Loc: nodes.Loc{Start: start},
	}

	// NOT MATCHED or MATCHED
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		mc.Matched = false
		if p.cur.Type != kwMATCHED {
			return nil, p.syntaxErrorAtCur()
		}
	} else {
		mc.Matched = true
		if p.cur.Type != kwMATCHED {
			return nil, p.syntaxErrorAtCur()
		}
	}

	// MATCHED
	p.advance()

	// Optional AND condition
	if p.cur.Type == kwAND {
		p.advance()
		var parseErr809 error
		mc.Condition, parseErr809 = p.parseExpr()
		if parseErr809 !=

			// THEN
			nil {
			return nil, parseErr809
		}
	}

	if p.cur.Type != kwTHEN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// UPDATE SET ... | INSERT ... | DELETE
	switch p.cur.Type {
	case kwUPDATE:
		p.advance() // consume UPDATE
		if p.cur.Type != kwSET {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume SET
		var parseErr810 error
		mc.UpdateSet, parseErr810 = p.parseMergeSetList()
		if parseErr810 !=

			// [ WHERE condition ]
			nil {
			return nil, parseErr810
		}
		if mc.UpdateSet.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}

		if p.cur.Type == kwWHERE {
			p.advance()
			var parseErr811 error
			mc.UpdateWhere, parseErr811 = p.parseExpr()
			if parseErr811 !=

				// [ DELETE WHERE condition ]
				nil {
				return nil, parseErr811
			}
		}

		if p.cur.Type == kwDELETE {
			p.advance()
			if p.cur.Type == kwWHERE {
				p.advance()
			}
			var parseErr812 error
			mc.DeleteWhere, parseErr812 = p.parseExpr()
			if parseErr812 != nil {
				return nil, parseErr812

				// consume INSERT
			}
		}

	case kwINSERT:
		p.advance()
		// Optional column list
		if p.cur.Type == '(' {
			p.advance()
			var parseErr813 error
			mc.InsertCols, parseErr813 = p.parseColumnList()
			if parseErr813 != nil {
				return nil, parseErr813
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// VALUES ( ... )
		if p.cur.Type != kwVALUES {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr814 error
		mc.InsertVals, parseErr814 = p.parseExprList()
		if parseErr814 != nil {
			return nil, parseErr814
		}
		if mc.InsertVals.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

		// [ WHERE condition ]
		if p.cur.Type == kwWHERE {
			p.advance()
			var parseErr815 error
			mc.InsertWhere, parseErr815 = p.parseExpr()
			if parseErr815 != nil {
				return nil, parseErr815

				// consume DELETE
			}
		}

	case kwDELETE:
		p.advance()
		mc.IsDelete = true
	default:
		return nil, p.syntaxErrorAtCur()
	}

	mc.Loc.End = p.prev.End
	return mc, nil
}

// parseMergeSetList parses a comma-separated list of SET clauses (col = expr).
func (p *Parser) parseMergeSetList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		sc, parseErr816 := p.parseMergeSetClause()
		if parseErr816 != nil {
			return nil, parseErr816
		}
		if sc == nil {
			break
		}
		list.Items = append(list.Items, sc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseMergeSetClause parses a single col = { expr | DEFAULT } assignment.
func (p *Parser) parseMergeSetClause() (*nodes.SetClause, error) {
	if !p.isIdentLike() {
		return nil, nil
	}
	start := p.pos()
	col, parseErr817 := p.parseColumnRef()
	if parseErr817 != nil {
		return nil, parseErr817
	}

	sc := &nodes.SetClause{
		Column: col,
		Loc:    nodes.Loc{Start: start},
	}

	// =
	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT or expression
	if p.cur.Type == kwDEFAULT {
		sc.Value = &nodes.ColumnRef{
			Column: "DEFAULT",
			Loc:    nodes.Loc{Start: p.pos(), End: p.prev.End},
		}
		p.advance()
	} else {
		var parseErr818 error
		sc.Value, parseErr818 = p.parseExpr()
		if parseErr818 != nil {
			return nil, parseErr818
		}
	}
	sc.Loc.End = p.prev.End
	return sc, nil
}

// parseColumnList parses a comma-separated list of column references.
func (p *Parser) parseColumnList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		col, parseErr819 := p.parseColumnRef()
		if parseErr819 != nil {
			return nil, parseErr819
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}
