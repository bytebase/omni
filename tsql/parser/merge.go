// Package parser - merge.go implements T-SQL MERGE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseMergeStmt parses a MERGE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/merge-transact-sql
//
//	MERGE [INTO] target [AS alias]
//	USING source [AS alias] ON condition
//	WHEN MATCHED [AND condition] THEN UPDATE SET ... | DELETE
//	WHEN NOT MATCHED [BY TARGET] [AND condition] THEN INSERT [(cols)] VALUES (...)
//	WHEN NOT MATCHED BY SOURCE [AND condition] THEN UPDATE SET ... | DELETE
//	[OUTPUT ...]
func (p *Parser) parseMergeStmt() *nodes.MergeStmt {
	loc := p.pos()
	p.advance() // consume MERGE

	stmt := &nodes.MergeStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional INTO
	p.match(kwINTO)

	// Target table
	stmt.Target = p.parseTableRef()

	// Optional alias
	alias := p.parseOptionalAlias()
	if alias != "" {
		stmt.Target.Alias = alias
	}

	// USING source (USING is not a reserved keyword)
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "using") {
		p.advance()
	}

	// Parse source table
	source := p.parseTableSource()
	stmt.Source = source

	// Source alias
	sourceAlias := p.parseOptionalAlias()
	if sourceAlias != "" {
		stmt.SourceAlias = sourceAlias
	}

	// ON condition
	if _, ok := p.match(kwON); ok {
		stmt.OnCondition = p.parseExpr()
	}

	// WHEN clauses
	var whenClauses []nodes.Node
	for p.cur.Type == kwWHEN {
		wc := p.parseMergeWhenClause()
		if wc != nil {
			whenClauses = append(whenClauses, wc)
		}
	}
	stmt.WhenClauses = &nodes.List{Items: whenClauses}

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMergeWhenClause parses a single WHEN clause in MERGE.
func (p *Parser) parseMergeWhenClause() *nodes.MergeWhenClause {
	loc := p.pos()
	p.advance() // consume WHEN

	wc := &nodes.MergeWhenClause{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		// MATCHED [BY TARGET] or MATCHED BY SOURCE
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "matched") {
			p.advance()
		}
		wc.Matched = false

		// BY TARGET or BY SOURCE
		if _, ok := p.match(kwBY); ok {
			if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "target") {
				p.advance()
				wc.ByTarget = true
			} else if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "source") {
				p.advance()
				wc.ByTarget = false
			}
		} else {
			wc.ByTarget = true // default NOT MATCHED means BY TARGET
		}
	} else if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "matched") {
		p.advance()
		wc.Matched = true
	}

	// Optional AND condition
	if _, ok := p.match(kwAND); ok {
		wc.Condition = p.parseExpr()
	}

	// THEN
	p.match(kwTHEN)

	// Action: UPDATE SET, DELETE, or INSERT
	switch {
	case p.cur.Type == kwUPDATE:
		p.advance() // consume UPDATE
		p.match(kwSET)
		setList := p.parseSetClauseList()
		wc.Action = &nodes.MergeUpdateAction{
			SetClause: setList,
			Loc:       nodes.Loc{Start: loc},
		}
	case p.cur.Type == kwDELETE:
		delLoc := p.pos()
		p.advance() // consume DELETE
		wc.Action = &nodes.MergeDeleteAction{
			Loc: nodes.Loc{Start: delLoc},
		}
	case p.cur.Type == kwINSERT:
		wc.Action = p.parseMergeInsertAction()
	}

	wc.Loc.End = p.pos()
	return wc
}

// parseMergeInsertAction parses INSERT [(cols)] VALUES (...) in a MERGE WHEN clause.
func (p *Parser) parseMergeInsertAction() *nodes.MergeInsertAction {
	loc := p.pos()
	p.advance() // consume INSERT

	action := &nodes.MergeInsertAction{
		Loc: nodes.Loc{Start: loc},
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
		_, _ = p.expect(')')
		action.Cols = &nodes.List{Items: cols}
	}

	// VALUES (...)
	if _, ok := p.match(kwVALUES); ok {
		if _, err := p.expect('('); err == nil {
			action.Values = p.parseExprList()
			_, _ = p.expect(')')
		}
	}

	action.Loc.End = p.pos()
	return action
}
