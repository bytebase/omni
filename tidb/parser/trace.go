package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseTraceStmt parses TRACE, TiDB's execution and optimizer tracer.
//
// Ref: TiDB v8.5.5 parser.y:5521-5559 (TraceStmt), 12628-12656
// (TraceableStmt); https://docs.pingcap.com/tidb/stable/sql-statement-trace/
//
//	TRACE traceable_stmt
//	TRACE FORMAT = 'format' traceable_stmt
//	TRACE PLAN traceable_stmt
//	TRACE PLAN TARGET = 'target' traceable_stmt
//
// FORMAT and PLAN are alternatives, not a list: they do not combine in either
// order, and TARGET belongs to PLAN alone. v8.5.5
// answers 1064 for `TRACE PLAN FORMAT = 'row' SELECT 1`, for
// `TRACE FORMAT = 'row' PLAN SELECT 1`, and for a bare TRACE with no
// statement after it.
//
// Both values are string literals and nothing else, which is where TRACE
// parts from EXPLAIN: `TRACE FORMAT = row SELECT 1` and
// `TRACE PLAN TARGET = estimation SELECT 1` are each a 1064, where
// `EXPLAIN FORMAT = ROW SELECT 1` is the bare-name arm of EXPLAIN's grammar.
// The text is not checked here either -- `TARGET = 'arandomstring'` parses,
// and v8.5.5 rejects the name when the statement runs.
//
// TRACE is not a keyword in this lexer and needs none: the dispatcher reads
// it with isBareWord, so `trace` stays usable as a table name, a column name
// and an implicit alias. TARGET is read the same way. Registering either would
// put every site that reads a name at risk, the implicit alias first, whose
// branch takes tokIDENT alone (TestTraceWordsStayNames).
//
// p.cur is the bare word TRACE.
func (p *Parser) parseTraceStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume TRACE

	stmt := &nodes.TraceStmt{Loc: nodes.Loc{Start: start}}

	// hasModifier rather than a test on the fields: FORMAT = '' is a
	// modifier all the same, and its name is the empty string.
	hasModifier := false

	switch {
	case p.cur.Type == kwFORMAT:
		hasModifier = true
		p.advance() // consume FORMAT
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokSCONST)
		if err != nil {
			return nil, err
		}
		stmt.Format = tok.Str
	case p.cur.Type == kwPLAN:
		hasModifier = true
		p.advance() // consume PLAN
		stmt.TracePlan = true
		// TARGET is read as a bare word: nothing but this clause reads it,
		// and upstream's keyword token takes no backquoted name either.
		if p.isBareWord("TARGET") {
			p.advance() // consume TARGET
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			tok, err := p.expect(tokSCONST)
			if err != nil {
				return nil, err
			}
			stmt.Target = tok.Str
		}
	}

	// Completion: after TRACE [options], offer the traceable statement
	// keywords, and the two modifiers while neither has been read.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwSELECT, kwWITH, kwTABLE, kwVALUES,
			kwINSERT, kwREPLACE, kwUPDATE, kwDELETE,
			kwLOAD, kwANALYZE, kwSET,
			kwBEGIN, kwSTART, kwCOMMIT, kwROLLBACK, kwSAVEPOINT, kwRELEASE,
		} {
			p.addTokenCandidate(t)
		}
		if !hasModifier {
			p.addTokenCandidate(kwFORMAT)
			p.addTokenCandidate(kwPLAN)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	inner, err := p.parseTraceableStmt()
	if err != nil {
		return nil, err
	}
	stmt.Stmt = inner
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTraceableStmt parses the statement TRACE carries: upstream's
// TraceableStmt.
//
// Ref: TiDB v8.5.5 parser.y:12628-12656.
//
//	traceable_stmt: {
//	    SELECT | TABLE | VALUES | WITH ... SELECT | ( query )
//	  | INSERT | REPLACE | UPDATE | DELETE
//	  | LOAD DATA | ANALYZE TABLE | SET
//	  | BEGIN | START TRANSACTION | COMMIT | ROLLBACK
//	  | SAVEPOINT | RELEASE SAVEPOINT
//	}
//
// The set holds no DDL at all, which is what makes TRACE unlike EXPLAIN:
// ALTER TABLE is explainable on TiDB and not traceable, so
// `TRACE ALTER TABLE t ADD COLUMN c` is a 1064 and the catalog never has a
// traced DDL to apply. CREATE TABLE, DROP, TRUNCATE and RENAME are each a
// 1064 too, and so are EXPLAIN, SHOW, ADMIN, IMPORT INTO, DO, PREPARE,
// EXECUTE, KILL and FLUSH. Every verdict here was measured against
// pingcap/tidb:v8.5.5 and is pinned by TestTraceFormsMatchTiDB.
//
// Three arms are narrower than the keyword that opens them, because upstream
// names one production where omni's dispatcher reads a family:
//
//   - LOAD reaches LoadDataStmt only. `TRACE LOAD XML INFILE 'x' INTO TABLE t`
//     and `TRACE LOAD STATS 'x'` are both a 1064.
//   - START reaches BeginTransactionStmt only. `TRACE START REPLICA` is a 1064.
//   - SET reaches upstream's SetStmt, which spells SET, SET PASSWORD,
//     SET [GLOBAL|SESSION] TRANSACTION, SET CONFIG, SET SESSION_STATES and
//     SET RESOURCE GROUP. SET ROLE and SET DEFAULT ROLE are productions of
//     their own and not in it, so `TRACE SET ROLE r` is a 1064. Of the six in
//     it, SET CONFIG and SET SESSION_STATES are not read by this parser at
//     all -- `SET SESSION_STATES 'x'` is a gap row of its own -- so each is a
//     rejection here for that reason rather than for anything TRACE does.
//
// Each is enforced on the node the family parser returned rather than by
// looking further ahead, so the inner statement is read by the one parser
// that owns it.
//
// A traceable statement this parser cannot read on its own is still a
// rejection here: `TRACE BEGIN OPTIMISTIC` fails because `BEGIN OPTIMISTIC`
// does, which is a gap in the transaction grammar that the corpus records as
// a row of its own. These comments record such gaps rather than close them;
// each belongs to the grammar of the statement TRACE carries.
func (p *Parser) parseTraceableStmt() (nodes.StmtNode, error) {
	// open is the keyword the statement starts with, which is where a
	// narrowed arm reports its rejection: the family parser read a
	// statement upstream's TraceableStmt does not name, and upstream's own
	// parser fails on the token that cannot reduce.
	open := p.cur
	switch p.cur.Type {
	case kwSELECT, kwWITH, '(', kwTABLE, kwVALUES:
		// One query expression, which parseSelectStmt reads whole: a SELECT
		// leaf, a WITH clause in front of one, a parenthesized body, a TABLE
		// or VALUES primary, or a set operation over them (SelectStmt,
		// SelectStmtWithClause, SetOprStmt and SubSelect in TraceableStmt).
		// A TABLE or VALUES target is the query expression it begins rather
		// than a leaf, so `TRACE TABLE t1 UNION TABLE t2` traces the set
		// operation; it is accepted on v8.5.5, as `TRACE VALUES ROW(1) UNION
		// VALUES ROW(2)` is.
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwREPLACE:
		return p.parseReplaceStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwANALYZE:
		return p.parseAnalyzeTableStmt()
	case kwBEGIN:
		return p.parseBeginStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVEPOINT:
		return p.parseSavepointStmt()
	case kwRELEASE:
		return p.parseReleaseSavepointStmt()
	case kwLOAD:
		node, err := p.parseLoadDispatch()
		if err != nil {
			return nil, err
		}
		load, ok := node.(*nodes.LoadDataStmt)
		if !ok || load.IsXML {
			return nil, p.syntaxErrorAtTok(open)
		}
		return load, nil
	case kwSTART:
		node, err := p.parseStartDispatch()
		if err != nil {
			return nil, err
		}
		begin, ok := node.(*nodes.BeginStmt)
		if !ok {
			return nil, p.syntaxErrorAtTok(open)
		}
		return begin, nil
	case kwSET:
		node, err := p.parseSetStmt()
		if err != nil {
			return nil, err
		}
		switch node.(type) {
		case *nodes.SetRoleStmt, *nodes.SetDefaultRoleStmt:
			return nil, p.syntaxErrorAtTok(open)
		}
		inner, ok := node.(nodes.StmtNode)
		if !ok {
			return nil, p.syntaxErrorAtTok(open)
		}
		return inner, nil
	}
	return nil, p.syntaxErrorAtCur()
}
