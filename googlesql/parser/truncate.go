package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-dml` DAG node. It implements GoogleSQL's
// TRUNCATE TABLE statement (GoogleSQLParser.g4 §2.3 truncate_statement), a
// hand-port of Google's open-source ZetaSQL reference.
//
// Grammar (truncate_statement):
//
//	TRUNCATE TABLE maybe_dashed_path_expression opt_where_expression?
//
// TRUNCATE is a BigQuery DML statement (truth1 BigQuery DML-003). Spanner has no
// TRUNCATE — its DDL path syntax-rejects it (oracle: "Error parsing Spanner DDL
// statement: … Encountered 'TRUNCATE'"), which is NON-authoritative for the
// BigQuery + Spanner union. TRUNCATE forms are therefore oracled against the
// canonical ZetaSQL corpus (truncate.sql) + the BigQuery docs (divergence
// ledger: TRUNCATE non-authoritative on Spanner). The grammar accepts an
// optional trailing WHERE (`TRUNCATE TABLE foo WHERE bar > 3`).

// parseTruncateStmt parses a TRUNCATE TABLE statement. TRUNCATE is the current
// token.
func (p *Parser) parseTruncateStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	p.advance() // TRUNCATE

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &ast.TruncateStmt{Loc: ast.Loc{Start: start}}

	// maybe_dashed_path_expression — a (possibly dashed) table path.
	target, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Target = target
	stmt.Loc.End = target.Loc.End

	// opt_where_expression.
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc.End = nodeLoc(where).End
	}

	return stmt, nil
}
