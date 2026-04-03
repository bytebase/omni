// Package parser - merge.go implements T-SQL MERGE statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseMergeStmt parses a MERGE statement.
//
// BNF: mssql/parser/bnf/merge-transact-sql.bnf
//
//	[ WITH <common_table_expression> [,...n] ]
//	MERGE
//	    [ TOP ( expression ) [ PERCENT ] ]
//	    [ INTO ] <target_table> [ WITH ( <merge_hint> ) ] [ [ AS ] table_alias ]
//	    USING <table_source> [ [ AS ] table_alias ]
//	    ON <merge_search_condition>
//	    [ WHEN MATCHED [ AND <clause_search_condition> ]
//	        THEN <merge_matched> ] [ ...n ]
//	    [ WHEN NOT MATCHED [ BY TARGET ] [ AND <clause_search_condition> ]
//	        THEN <merge_not_matched> ]
//	    [ WHEN NOT MATCHED BY SOURCE [ AND <clause_search_condition> ]
//	        THEN <merge_matched> ] [ ...n ]
//	    [ <output_clause> ]
//	    [ OPTION ( <query_hint> [ ,...n ] ) ]
//	;
//
//	<merge_matched>::= { UPDATE SET <set_clause> | DELETE }
//	<merge_not_matched>::= { INSERT [ ( column_list ) ] { VALUES ( values_list ) | DEFAULT VALUES } }
func (p *Parser) parseMergeStmt() (*nodes.MergeStmt, error) {
	loc := p.pos()
	p.advance() // consume MERGE

	stmt := &nodes.MergeStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		top, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
	}

	// Optional INTO
	p.match(kwINTO)

	// Completion: after MERGE [INTO] → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Target table
	target, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, p.newParseError(p.cur.Loc, "expected target table after MERGE")
	}
	stmt.Target = target

	// Optional WITH ( <merge_hint> ) on target
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		hints, err := p.parseTableHints()
		if err != nil {
			return nil, err
		}
		stmt.Target.Hints = hints
	}

	// Optional alias
	alias := p.parseOptionalAlias()
	if alias != "" {
		stmt.Target.Alias = alias
	}

	// USING source (USING is not a reserved keyword)
	if p.cur.Type == kwUSING {
		p.advance()
		// Completion: after USING → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
	}

	// Parse source table
	source, err := p.parseTableSource()
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, p.newParseError(p.cur.Loc, "expected source table after USING")
	}
	stmt.Source = source

	// Source alias
	sourceAlias := p.parseOptionalAlias()
	if sourceAlias != "" {
		stmt.SourceAlias = sourceAlias
	}

	// ON condition
	if _, ok := p.match(kwON); ok {
		// Completion: after ON in MERGE → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if cond == nil {
			return nil, p.newParseError(p.cur.Loc, "expected condition after ON")
		}
		stmt.OnCondition = cond
	}

	// WHEN clauses
	var whenClauses []nodes.Node
	for p.cur.Type == kwWHEN {
		wc, err := p.parseMergeWhenClause()
		if err != nil {
			return nil, err
		}
		if wc != nil {
			whenClauses = append(whenClauses, wc)
		}
	}
	stmt.WhenClauses = &nodes.List{Items: whenClauses}

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		oc, err := p.parseOutputClause()
		if err != nil {
			return nil, err
		}
		stmt.OutputClause = oc
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		oc, err := p.parseOptionClause()
		if err != nil {
			return nil, err
		}
		stmt.OptionClause = oc
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseMergeWhenClause parses a single WHEN clause in MERGE.
//
//	WHEN MATCHED [ AND <clause_search_condition> ] THEN <merge_matched>
//	WHEN NOT MATCHED [ BY TARGET ] [ AND <clause_search_condition> ] THEN <merge_not_matched>
//	WHEN NOT MATCHED BY SOURCE [ AND <clause_search_condition> ] THEN <merge_matched>
func (p *Parser) parseMergeWhenClause() (*nodes.MergeWhenClause, error) {
	loc := p.pos()
	p.advance() // consume WHEN

	// Completion: after WHEN → MATCHED, NOT
	if p.collectMode() {
		p.addRuleCandidate("matched_keyword")
		p.addTokenCandidate(kwNOT)
		return nil, errCollecting
	}

	wc := &nodes.MergeWhenClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		// MATCHED [BY TARGET] or MATCHED BY SOURCE
		if p.cur.Type == kwMATCHED {
			p.advance()
		}
		wc.Matched = false

		// BY TARGET or BY SOURCE
		if _, ok := p.match(kwBY); ok {
			if p.cur.Type == kwTARGET {
				p.advance()
				wc.ByTarget = true
			} else if p.cur.Type == kwSOURCE {
				p.advance()
				wc.ByTarget = false
			}
		} else {
			wc.ByTarget = true // default NOT MATCHED means BY TARGET
		}
	} else if p.cur.Type == kwMATCHED {
		p.advance()
		wc.Matched = true
	}

	// Optional AND condition
	if _, ok := p.match(kwAND); ok {
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		wc.Condition = cond
	}

	// THEN
	p.match(kwTHEN)

	// Completion: after THEN → action keywords based on match type
	if p.collectMode() {
		if wc.Matched {
			// WHEN MATCHED THEN → UPDATE, DELETE
			p.addTokenCandidate(kwUPDATE)
			p.addTokenCandidate(kwDELETE)
		} else {
			// WHEN NOT MATCHED THEN → INSERT
			p.addTokenCandidate(kwINSERT)
		}
		return nil, errCollecting
	}

	// Action: UPDATE SET, DELETE, or INSERT
	switch {
	case p.cur.Type == kwUPDATE:
		p.advance() // consume UPDATE
		p.match(kwSET)
		setList, err := p.parseSetClauseList()
		if err != nil {
			return nil, err
		}
		wc.Action = &nodes.MergeUpdateAction{
			SetClause: setList,
			Loc:       nodes.Loc{Start: loc, End: -1},
		}
	case p.cur.Type == kwDELETE:
		delLoc := p.pos()
		p.advance() // consume DELETE
		wc.Action = &nodes.MergeDeleteAction{
			Loc: nodes.Loc{Start: delLoc, End: -1},
		}
	case p.cur.Type == kwINSERT:
		action, err := p.parseMergeInsertAction()
		if err != nil {
			return nil, err
		}
		wc.Action = action
	}

	wc.Loc.End = p.prevEnd()
	return wc, nil
}

// parseMergeInsertAction parses INSERT [(cols)] { VALUES (...) | DEFAULT VALUES } in a MERGE WHEN clause.
//
//	INSERT [ ( column_list ) ] { VALUES ( values_list ) | DEFAULT VALUES }
func (p *Parser) parseMergeInsertAction() (*nodes.MergeInsertAction, error) {
	loc := p.pos()
	p.advance() // consume INSERT

	action := &nodes.MergeInsertAction{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			colName, ok := p.parseIdentifier()
			if !ok {
				break
			}
			cols = append(cols, &nodes.String{Str: colName})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		action.Cols = &nodes.List{Items: cols}
	}

	// DEFAULT VALUES or VALUES (...)
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if p.cur.Type == kwVALUES {
			p.advance() // consume VALUES
		}
		action.DefaultValues = true
	} else if _, ok := p.match(kwVALUES); ok {
		if _, err := p.expect('('); err == nil {
			vals, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			action.Values = vals
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}

	action.Loc.End = p.prevEnd()
	return action, nil
}
