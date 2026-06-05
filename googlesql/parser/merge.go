package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-dml` DAG node. It implements GoogleSQL's
// MERGE statement (GoogleSQLParser.g4 §2.7 merge_statement), a hand-port of
// Google's open-source ZetaSQL reference. Shared DML helpers (column list,
// update_item_list, insert_values_row) live in insert.go.
//
// Grammar (merge_statement):
//
//	MERGE INTO? <target> as_alias? USING merge_source ON expression (merge_when_clause)+
//	merge_source: table_path_expression | table_subquery
//	merge_when_clause:
//	    WHEN MATCHED [AND expr] THEN merge_action
//	  | WHEN NOT MATCHED [BY TARGET] [AND expr] THEN merge_action
//	  | WHEN NOT MATCHED BY SOURCE [AND expr] THEN merge_action
//	merge_action:
//	    INSERT column_list? (VALUES insert_values_row | ROW)
//	  | UPDATE SET update_item_list
//	  | DELETE
//
// MERGE is BigQuery-only among the bytebase consumers: Spanner feature-rejects
// it AFTER parse (the emulator answers "Statement not supported: MergeStmt" at
// statement dispatch, BEFORE validating the body — so the emulator gives NO
// authoritative syntax verdict on MERGE bodies). MERGE forms are therefore
// oracled against the canonical ZetaSQL corpus (dml_merge.sql) + the BigQuery
// docs, NOT the emulator (divergence ledger: MERGE non-authoritative on Spanner;
// the grammar requires (merge_when_clause)+ but the emulator accepts a MERGE
// with zero WHEN clauses — we follow the .g4 and require >= 1).

// parseMergeStmt parses a MERGE statement. MERGE is the current token.
func (p *Parser) parseMergeStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	p.advance() // MERGE

	stmt := &ast.MergeStmt{Loc: ast.Loc{Start: start}}

	// opt_into.
	p.match(kwINTO)

	// Target (maybe_dashed_path_expression — a table name).
	target, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Optional `[AS] alias`. USING (a reserved keyword) terminates the target, so
	// a bare identifier before USING is the implicit alias.
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		stmt.Alias = alias
	}

	// USING merge_source.
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}
	source, err := p.parseMergeSource()
	if err != nil {
		return nil, err
	}
	stmt.Source = source

	// ON expression.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	on, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.On = on

	// (merge_when_clause)+ — at least one is required.
	for p.cur.Type == kwWHEN {
		when, err := p.parseMergeWhen()
		if err != nil {
			return nil, err
		}
		stmt.Whens = append(stmt.Whens, when)
	}
	if len(stmt.Whens) == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseMergeSource parses merge_source: `table_path_expression | table_subquery`.
// A leading '(' is a parenthesized subquery; otherwise it is a table path
// (optionally aliased). The result is a *TableExpr (Path or Subquery set).
func (p *Parser) parseMergeSource() (ast.Node, error) {
	if p.cur.Type == int('(') {
		openTok := p.advance() // '('
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		te := &ast.TableExpr{Subquery: query, Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}
		// table_subquery takes an optional `[AS] alias`.
		if err := p.parseTableAliasOnly(te); err != nil {
			return nil, err
		}
		return te, nil
	}

	// table_path_expression: a table path with an optional `[AS] alias`.
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	te := &ast.TableExpr{Path: path, Loc: path.Loc}
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		te.Alias = alias
		te.Loc.End = p.prev.Loc.End
	}
	return te, nil
}

// parseMergeWhen parses one merge_when_clause. WHEN is the current token.
func (p *Parser) parseMergeWhen() (*ast.MergeWhen, error) {
	whenTok := p.advance() // WHEN
	when := &ast.MergeWhen{Loc: whenTok.Loc}

	if _, ok := p.match(kwNOT); ok {
		// WHEN NOT MATCHED [BY TARGET | BY SOURCE] [AND expr] THEN …
		if _, err := p.expect(kwMATCHED); err != nil {
			return nil, err
		}
		when.NotMatched = true
		if p.cur.Type == kwBY {
			p.advance() // BY
			switch p.cur.Type {
			case kwTARGET:
				p.advance()
				when.ByTarget = true
			case kwSOURCE:
				p.advance()
				when.BySource = true
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}
	} else {
		// WHEN MATCHED [AND expr] THEN …
		if _, err := p.expect(kwMATCHED); err != nil {
			return nil, err
		}
		when.Matched = true
	}

	// opt_and_expression: AND expr.
	if _, ok := p.match(kwAND); ok {
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		when.And = cond
	}

	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}
	action, err := p.parseMergeAction()
	if err != nil {
		return nil, err
	}
	when.Action = action
	when.Loc.End = p.prev.Loc.End
	return when, nil
}

// parseMergeAction parses merge_action:
//
//	INSERT column_list? (VALUES insert_values_row | ROW)
//	UPDATE SET update_item_list
//	DELETE
//
// The current token is the action keyword (INSERT / UPDATE / DELETE).
func (p *Parser) parseMergeAction() (*ast.MergeAction, error) {
	start := p.cur.Loc
	switch p.cur.Type {
	case kwINSERT:
		p.advance() // INSERT
		act := &ast.MergeAction{Kind: ast.MergeInsert, Loc: start}
		// Optional column_list.
		if p.cur.Type == int('(') {
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			act.Columns = cols
		}
		// merge_insert_value_list_or_source_row: VALUES insert_values_row | ROW.
		switch p.cur.Type {
		case kwVALUES:
			p.advance() // VALUES
			row, err := p.parseInsertValuesRow()
			if err != nil {
				return nil, err
			}
			act.InsertRow = row
		case kwROW:
			p.advance() // ROW
			act.SourceRow = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		act.Loc.End = p.prev.Loc.End
		return act, nil

	case kwUPDATE:
		p.advance() // UPDATE
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		items, err := p.parseUpdateItemList()
		if err != nil {
			return nil, err
		}
		return &ast.MergeAction{Kind: ast.MergeUpdate, SetItems: items, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil

	case kwDELETE:
		delTok := p.advance() // DELETE
		return &ast.MergeAction{Kind: ast.MergeDelete, Loc: delTok.Loc}, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}
