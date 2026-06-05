package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-dml` DAG node. It implements GoogleSQL's
// UPDATE statement (GoogleSQLParser.g4 §2.7 update_statement), a hand-port of
// Google's open-source ZetaSQL reference. Shared DML helpers (target path,
// update_item_list, ASSERT_ROWS_MODIFIED, THEN RETURN) live in insert.go.
//
// Grammar (update_statement):
//
//	UPDATE <target> hint? as_alias? opt_with_offset_and_alias?
//	  SET update_item_list from_clause? opt_where_expression? opt_assert? opt_returning?
//
// <target> is maybe_dashed_generalized_path_expression; opt_with_offset_and_alias
// is `WITH OFFSET [[AS] name]`; the SET list is non-empty. The FROM clause is the
// join-update source list (`UPDATE T SET x = S.y FROM S WHERE …`).

// parseUpdateStmt parses an UPDATE statement. UPDATE is the current token.
func (p *Parser) parseUpdateStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	p.advance() // UPDATE

	stmt := &ast.UpdateStmt{Loc: ast.Loc{Start: start}}

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

	// Optional `[AS] alias`. SET (a reserved keyword) terminates the target, so a
	// bare identifier before SET is the implicit alias.
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		stmt.Alias = alias
	}

	// Optional `WITH OFFSET [[AS] name]`.
	if err := p.parseWithOffsetAndAlias(&stmt.WithOffset, &stmt.WithOffsetAlias); err != nil {
		return nil, err
	}

	// SET update_item_list.
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	items, err := p.parseUpdateItemList()
	if err != nil {
		return nil, err
	}
	stmt.Items = items

	// Optional FROM clause (join-update source list).
	if p.cur.Type == kwFROM {
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
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

// parseWithOffsetAndAlias parses opt_with_offset_and_alias:
// `WITH OFFSET [[AS] name]`. It sets *have to true and fills *alias when the
// clause is present, and is a no-op otherwise. Used by UPDATE and DELETE.
func (p *Parser) parseWithOffsetAndAlias(have *bool, alias *string) error {
	if !(p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET) {
		return nil
	}
	p.advance() // WITH
	p.advance() // OFFSET
	*have = true
	if a, ok, err := p.tryParseTableAlias(); err != nil {
		return err
	} else if ok {
		*alias = a
	}
	return nil
}
