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
func (p *Parser) parseMergeStmt() *nodes.MergeStmt {
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
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// INTO
	if p.cur.Type == kwINTO {
		p.advance()
	}

	// Target table
	stmt.Target = p.parseObjectName()

	// Optional target alias
	if p.cur.Type == kwAS {
		p.advance()
		stmt.TargetAlias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		stmt.TargetAlias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// USING
	if p.cur.Type == kwUSING {
		p.advance()
	}

	// Source table or subquery
	stmt.Source = p.parseTableRef()

	// Optional source alias — parseTableRef already handles alias on TableRef,
	// but for MergeStmt we want it in SourceAlias. Extract it.
	if tr, ok := stmt.Source.(*nodes.TableRef); ok && tr.Alias != nil {
		stmt.SourceAlias = tr.Alias
		tr.Alias = nil
	}

	// ON ( condition )
	if p.cur.Type == kwON {
		p.advance()
		// The condition may or may not be wrapped in parens
		if p.cur.Type == '(' {
			p.advance()
			stmt.On = p.parseExpr()
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			stmt.On = p.parseExpr()
		}
	}

	// WHEN clauses
	stmt.Clauses = &nodes.List{}
	for p.cur.Type == kwWHEN {
		mc := p.parseMergeClause()
		if mc != nil {
			stmt.Clauses.Items = append(stmt.Clauses.Items, mc)
		}
	}

	// error_logging_clause
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMergeClause parses a single WHEN MATCHED / WHEN NOT MATCHED clause.
func (p *Parser) parseMergeClause() *nodes.MergeClause {
	start := p.pos()
	p.advance() // consume WHEN

	mc := &nodes.MergeClause{
		Loc: nodes.Loc{Start: start},
	}

	// NOT MATCHED or MATCHED
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		mc.Matched = false
	} else {
		mc.Matched = true
	}

	// MATCHED
	if p.cur.Type == kwMATCHED {
		p.advance()
	}

	// Optional AND condition
	if p.cur.Type == kwAND {
		p.advance()
		mc.Condition = p.parseExpr()
	}

	// THEN
	if p.cur.Type == kwTHEN {
		p.advance()
	}

	// UPDATE SET ... | INSERT ... | DELETE
	switch p.cur.Type {
	case kwUPDATE:
		p.advance() // consume UPDATE
		if p.cur.Type == kwSET {
			p.advance() // consume SET
		}
		mc.UpdateSet = p.parseMergeSetList()

		// [ WHERE condition ]
		if p.cur.Type == kwWHERE {
			p.advance()
			mc.UpdateWhere = p.parseExpr()
		}

		// [ DELETE WHERE condition ]
		if p.cur.Type == kwDELETE {
			p.advance()
			if p.cur.Type == kwWHERE {
				p.advance()
			}
			mc.DeleteWhere = p.parseExpr()
		}

	case kwINSERT:
		p.advance() // consume INSERT
		// Optional column list
		if p.cur.Type == '(' {
			p.advance()
			mc.InsertCols = p.parseColumnList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// VALUES ( ... )
		if p.cur.Type == kwVALUES {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				mc.InsertVals = p.parseExprList()
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}

		// [ WHERE condition ]
		if p.cur.Type == kwWHERE {
			p.advance()
			mc.InsertWhere = p.parseExpr()
		}

	case kwDELETE:
		p.advance() // consume DELETE
		mc.IsDelete = true
	}

	mc.Loc.End = p.pos()
	return mc
}

// parseMergeSetList parses a comma-separated list of SET clauses (col = expr).
func (p *Parser) parseMergeSetList() *nodes.List {
	list := &nodes.List{}
	for {
		sc := p.parseMergeSetClause()
		if sc == nil {
			break
		}
		list.Items = append(list.Items, sc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseMergeSetClause parses a single col = { expr | DEFAULT } assignment.
func (p *Parser) parseMergeSetClause() *nodes.SetClause {
	if !p.isIdentLike() {
		return nil
	}
	start := p.pos()
	col := p.parseColumnRef()

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
			Loc:    nodes.Loc{Start: p.pos(), End: p.pos()},
		}
		p.advance()
	} else {
		sc.Value = p.parseExpr()
	}
	sc.Loc.End = p.pos()
	return sc
}

// parseColumnList parses a comma-separated list of column references.
func (p *Parser) parseColumnList() *nodes.List {
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		col := p.parseColumnRef()
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}
