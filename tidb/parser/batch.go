package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseBatchStmt parses a TiDB non-transactional DML statement:
//
//	BATCH [ON <col>] LIMIT <n> [DRY RUN [QUERY]] {DELETE | UPDATE | INSERT | REPLACE}
//
// Grammar ref: pingcap parser.y NonTransactionalDMLStmt —
// "BATCH" OptionalShardColumn "LIMIT" NUM DryRunOptions ShardableStmt.
// LIMIT takes a bare integer (NUM), not an expression.
func (p *Parser) parseBatchStmt() (*nodes.BatchStmt, error) {
	start := p.pos()
	p.advance() // consume BATCH

	stmt := &nodes.BatchStmt{Loc: nodes.Loc{Start: start}}

	// Completion: after BATCH, offer ON / LIMIT.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwON)
		p.addTokenCandidate(kwLIMIT)
		return nil, &ParseError{Message: "collecting"}
	}

	// OptionalShardColumn: [ON ColumnName]
	if _, ok := p.match(kwON); ok {
		col, err := p.parseColumnRef()
		if err != nil {
			return nil, err
		}
		stmt.ShardColumn = col
	}

	// LIMIT NUM (mandatory). Grammar is NUM, so reject expressions, placeholders,
	// and signed/parenthesized forms by requiring an integer-literal token.
	if _, err := p.expect(kwLIMIT); err != nil {
		return nil, err
	}
	if p.cur.Type != tokICONST {
		return nil, p.syntaxErrorAtCur()
	}
	numTok := p.advance()
	limit, err := strconv.ParseUint(numTok.Str, 10, 64)
	if err != nil {
		return nil, &ParseError{Message: "invalid BATCH LIMIT value", Position: numTok.Loc}
	}
	stmt.Limit = limit

	// DryRunOptions: [] | DRY RUN | DRY RUN QUERY
	if _, ok := p.match(kwDRY); ok {
		if _, err := p.expect(kwRUN); err != nil {
			return nil, err
		}
		if _, ok := p.match(kwQUERY); ok {
			stmt.DryRun = nodes.BatchDryRunQuery
		} else {
			stmt.DryRun = nodes.BatchDryRunSplitDML
		}
	}

	// Completion: before the DML, offer DRY plus the shardable statement starters.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwDRY, kwDELETE, kwUPDATE, kwINSERT, kwREPLACE} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// ShardableStmt: DELETE | UPDATE | INSERT | REPLACE.
	var dml nodes.StmtNode
	switch p.cur.Type {
	case kwDELETE:
		dml, err = p.parseDeleteStmt()
	case kwUPDATE:
		dml, err = p.parseUpdateStmt()
	case kwINSERT:
		dml, err = p.parseInsertStmt()
	case kwREPLACE:
		dml, err = p.parseReplaceStmt()
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	stmt.DML = dml

	stmt.Loc.End = p.pos()
	return stmt, nil
}
