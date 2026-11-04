package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parsePlanReplayerStmt parses PLAN REPLAYER, TiDB's plan-reproduction dump.
//
// Ref: TiDB v8.5.5 parser.y:16016-16177 (PlanReplayerStmt,
// PlanReplayerDumpOpt);
// https://docs.pingcap.com/tidb/stable/sql-plan-replayer/
//
//	PLAN REPLAYER DUMP [WITH STATS AS OF TIMESTAMP expr] EXPLAIN [ANALYZE] explainable_stmt
//	PLAN REPLAYER DUMP [WITH STATS AS OF TIMESTAMP expr] EXPLAIN [ANALYZE] SLOW QUERY [WHERE expr] [ORDER BY ...] [LIMIT ...]
//	PLAN REPLAYER DUMP [WITH STATS AS OF TIMESTAMP expr] EXPLAIN [ANALYZE] 'file'
//	PLAN REPLAYER LOAD 'file'
//	PLAN REPLAYER CAPTURE [REMOVE] 'sql_digest' 'plan_digest'
//
// REPLAYER, DUMP, CAPTURE and STATS are read with isBareWord rather than
// registered: each stays usable as a table name, a column name and an
// implicit alias, and a backquoted one is no part of this grammar on either
// side (TestPlanReplayerWordsStayNames). PLAN is already a keyword here,
// registered for the SQL plan bindings, and no other statement starts with
// it, so REPLAYER after it is what makes a statement -- the way WATCH after
// QUERY does.
//
// ANALYZE is EXPLAIN's modifier and carries EXPLAIN's effect: it runs the
// statement, and the statement may be DDL. On pingcap/tidb:v8.5.5
// `PLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN c` adds the
// column, while `PLAN REPLAYER DUMP EXPLAIN ALTER TABLE t ADD COLUMN c`
// writes a dump and leaves the table alone. That is the same split EXPLAIN
// has, measured the same way, and it is why tidb/catalog/exec.go applies the
// ANALYZE form and review's walk follows it.
//
// CAPTURE takes exactly two string literals, and LOAD exactly one: a missing
// one, a third, and an integer in place of either are each a 1064. Every
// verdict in this file was measured against v8.5.5 and is pinned by
// TestPlanReplayerFormsMatchTiDB.
//
// p.cur is kwPLAN.
func (p *Parser) parsePlanReplayerStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume PLAN

	if !p.isBareWord("REPLAYER") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume REPLAYER

	stmt := &nodes.PlanReplayerStmt{Loc: nodes.Loc{Start: start}}

	// This family offers no completion candidates, as the BRIE family does
	// not: two of the three openers are bare words rather than tokens, so a
	// candidate list here could only name LOAD and would read as though
	// DUMP and CAPTURE were not there.

	switch {
	case p.cur.Type == kwLOAD:
		p.advance() // consume LOAD
		tok, err := p.expect(tokSCONST)
		if err != nil {
			return nil, err
		}
		stmt.Load = true
		stmt.File = tok.Str
	case p.isBareWord("CAPTURE"):
		p.advance() // consume CAPTURE
		stmt.Capture = true
		// REMOVE is a keyword here already, registered for QUERY WATCH
		// REMOVE.
		if _, ok := p.match(kwREMOVE); ok {
			stmt.Capture = false
			stmt.Remove = true
		}
		sqlDigest, err := p.expect(tokSCONST)
		if err != nil {
			return nil, err
		}
		planDigest, err := p.expect(tokSCONST)
		if err != nil {
			return nil, err
		}
		stmt.SQLDigest = sqlDigest.Str
		stmt.PlanDigest = planDigest.Str
	case p.isBareWord("DUMP"):
		p.advance() // consume DUMP
		if err := p.parsePlanReplayerDump(stmt); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePlanReplayerDump parses what follows PLAN REPLAYER DUMP: the optional
// historical-statistics clause, EXPLAIN with its optional ANALYZE, and then
// one of the three things EXPLAIN can name here.
//
// Ref: TiDB v8.5.5 parser.y:16017-16122 (the DUMP arms), 16172-16177
// (PlanReplayerDumpOpt).
//
// WITH STATS is not optional in itself: the clause is WITH STATS AS OF
// TIMESTAMP expr whole, so `PLAN REPLAYER DUMP WITH STATS EXPLAIN SELECT 1`
// is a 1064 while `... WITH STATS AS OF TIMESTAMP NOW() EXPLAIN SELECT 1`
// parses.
//
// SLOW QUERY and a quoted file name are arms of EXPLAIN's slot, not
// statements: a file name has to be quoted, so
// `PLAN REPLAYER DUMP EXPLAIN `sql.txt“ is a 1064, and SLOW has to be
// followed by QUERY.
func (p *Parser) parsePlanReplayerDump(stmt *nodes.PlanReplayerStmt) error {
	if _, ok := p.match(kwWITH); ok {
		if !p.isBareWord("STATS") {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume STATS
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		if _, err := p.expect(kwTIMESTAMP); err != nil {
			return err
		}
		ts, err := p.parseExpr()
		if err != nil {
			return err
		}
		if ts == nil {
			return p.syntaxErrorAtCur()
		}
		stmt.HistoricalStats = ts
	}

	if _, err := p.expect(kwEXPLAIN); err != nil {
		return err
	}
	if _, ok := p.match(kwANALYZE); ok {
		stmt.Analyze = true
	}

	// EXPLAIN 'file', which reads its statements from a file rather than
	// carrying one.
	if p.cur.Type == tokSCONST {
		stmt.File = p.advance().Str
		return nil
	}

	// EXPLAIN [ANALYZE] SLOW QUERY, which dumps the slow-query log. SLOW is
	// a keyword here already, registered for ADMIN SHOW SLOW.
	if p.cur.Type == kwSLOW {
		p.advance() // consume SLOW
		if _, err := p.expect(kwQUERY); err != nil {
			return err
		}
		stmt.SlowQuery = true
		return p.parsePlanReplayerSlowQueryTail(stmt)
	}

	inner, err := p.parseExplainableStmt()
	if err != nil {
		return err
	}
	stmt.Stmt = inner
	return nil
}

// parsePlanReplayerSlowQueryTail parses the WHERE, ORDER BY and LIMIT that
// narrow the slow-query dump. All three are optional and in this order:
// upstream spells them WhereClauseOptional OrderByOptional
// SelectStmtLimitOpt, so the LIMIT takes both of SELECT's spellings --
// `LIMIT 1, 10` and `LIMIT 10 OFFSET 5` each parse on v8.5.5 -- and
// `PLAN REPLAYER DUMP EXPLAIN SLOW QUERY` with none of them parses too.
//
// Ref: TiDB v8.5.5 parser.y:16055-16100.
func (p *Parser) parsePlanReplayerSlowQueryTail(stmt *nodes.PlanReplayerStmt) error {
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return err
		}
		if !p.collectMode() && where == nil {
			return p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return err
		}
		stmt.OrderBy = orderBy
	}
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return err
		}
		stmt.Limit = limit
	}
	return nil
}

// parseExplainableStmt parses the statement PLAN REPLAYER DUMP EXPLAIN
// carries: upstream's ExplainableStmt, the same production EXPLAIN's own
// statement slot names.
//
// Ref: TiDB v8.5.5 parser.y:12658-12680.
//
//	explainable_stmt: {
//	    SELECT | TABLE | VALUES | WITH ... SELECT | ( query )
//	  | INSERT | REPLACE | UPDATE | DELETE
//	  | ALTER TABLE | IMPORT INTO
//	}
//
// ALTER TABLE is the one DDL in the set, and the only reason the catalog has
// to read what a dump carries. CREATE TABLE, TRUNCATE, SET, BEGIN and
// LOAD DATA are each a 1064 -- which is TRACE's set almost inverted, since
// TraceableStmt holds the transaction statements, SET and LOAD DATA and holds
// no DDL.
//
// IMPORT INTO is the other half of the set and is read here too, now that
// this parser has it: `PLAN REPLAYER DUMP EXPLAIN IMPORT INTO t FROM 'f'`
// parses on v8.5.5 and here.
// atExplainableStmt reports whether the tokens at p.cur open an
// ExplainableStmt. EXPLAIN's own statement slot needs the question answered
// before it commits, because EXPLAIN tbl_name is the arm it falls back to.
func (p *Parser) atExplainableStmt() bool {
	switch p.cur.Type {
	case kwSELECT, kwWITH, '(', kwTABLE, kwVALUES, kwINSERT, kwREPLACE, kwUPDATE, kwDELETE, kwALTER, kwIMPORT:
		return true
	}
	return false
}

func (p *Parser) parseExplainableStmt() (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwSELECT, kwWITH, '(', kwTABLE, kwVALUES:
		// One query expression, which parseSelectStmt reads whole. A TABLE
		// or VALUES target is the query expression it begins, not a leaf, so
		// `... EXPLAIN TABLE t1 UNION TABLE t2` explains the set operation
		// (accepted on v8.5.5).
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwREPLACE:
		return p.parseReplaceStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwALTER:
		alterStart := p.pos()
		p.advance() // consume ALTER
		if p.cur.Type != kwTABLE {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseAlterTableStmt(alterStart)
	case kwIMPORT:
		// parseImportDispatch owns the IMPORT INTO spelling, so the one
		// statement IMPORT opens is read by the one parser that has it.
		node, err := p.parseImportDispatch()
		if err != nil {
			return nil, err
		}
		inner, ok := node.(nodes.StmtNode)
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		return inner, nil
	}
	return nil, p.syntaxErrorAtCur()
}
