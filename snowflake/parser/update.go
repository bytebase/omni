package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// UPDATE statement parser
// ---------------------------------------------------------------------------

// parseUpdateStmt parses an UPDATE statement:
//
//	UPDATE table [alias] SET col = expr [, col = expr ...]
//	  [FROM table_refs]
//	  [WHERE cond]
func (p *Parser) parseUpdateStmt() (*ast.UpdateStmt, error) {
	updateTok, err := p.expect(kwUPDATE)
	if err != nil {
		return nil, err
	}

	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.UpdateStmt{
		Target: target,
		Loc:    ast.Loc{Start: updateTok.Loc.Start},
	}

	// Optional table alias (not AS — Snowflake allows bare alias)
	alias, hasAlias := p.parseOptionalAlias()
	if hasAlias {
		stmt.Alias = alias
	}

	// SET col = expr [, ...]
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	sets, err := p.parseUpdateSetList()
	if err != nil {
		return nil, err
	}
	stmt.Sets = sets

	// Optional FROM clause (Snowflake extension)
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// Optional WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseUpdateSetList parses comma-separated col = expr assignments.
func (p *Parser) parseUpdateSetList() ([]*ast.UpdateSet, error) {
	var sets []*ast.UpdateSet

	set, err := p.parseUpdateSet()
	if err != nil {
		return nil, err
	}
	sets = append(sets, set)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		set, err = p.parseUpdateSet()
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, nil
}

// parseUpdateSet parses one col = expr assignment.
// The column may be qualified (e.g. alias.column).
func (p *Parser) parseUpdateSet() (*ast.UpdateSet, error) {
	col, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	startLoc := col.Loc.Start

	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.UpdateSet{
		Column: col,
		Value:  val,
		Loc:    ast.Loc{Start: startLoc, End: p.prev.Loc.End},
	}, nil
}
