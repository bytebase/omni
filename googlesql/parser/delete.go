package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-dml` DAG node. It implements GoogleSQL's
// DELETE statement (GoogleSQLParser.g4 §2.7 delete_statement), a hand-port of
// Google's open-source ZetaSQL reference. Shared DML helpers (target path,
// ASSERT_ROWS_MODIFIED, THEN RETURN, WITH OFFSET) live in insert.go / update.go.
//
// Grammar (delete_statement):
//
//	DELETE FROM? <target> hint? as_alias? opt_with_offset_and_alias?
//	  opt_where_expression? opt_assert? opt_returning?
//
// The FROM keyword is OPTIONAL (ZetaSQL: `DELETE T WHERE …` parses, as does
// `DELETE FROM T …`). <target> is maybe_dashed_generalized_path_expression.

// parseDeleteStmt parses a DELETE statement. DELETE is the current token.
func (p *Parser) parseDeleteStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	p.advance() // DELETE

	stmt := &ast.DeleteStmt{Loc: ast.Loc{Start: start}}

	// Optional FROM.
	if _, ok := p.match(kwFROM); ok {
		stmt.From = true
	}

	// Target.
	target, err := p.parseDMLTargetPath()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Optional per-target `@{…}` hint.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	// Optional `[AS] alias`.
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		stmt.Alias = alias
	}

	// Optional `WITH OFFSET [[AS] name]`.
	if err := p.parseWithOffsetAndAlias(&stmt.WithOffset, &stmt.WithOffsetAlias); err != nil {
		return nil, err
	}

	// Optional WHERE.
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// opt_assert_rows_modified.
	if p.cur.Type == kwASSERT_ROWS_MODIFIED {
		ar, err := p.parseAssertRowsModified()
		if err != nil {
			return nil, err
		}
		stmt.AssertRows = ar
	}

	// opt_returning_clause.
	if p.cur.Type == kwTHEN {
		ret, err := p.parseReturning()
		if err != nil {
			return nil, err
		}
		stmt.Returning = ret
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
