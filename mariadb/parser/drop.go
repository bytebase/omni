package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseDropTableStmt parses a DROP TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-table.html
//
//	DROP [TEMPORARY] TABLE [IF EXISTS] tbl_name [, tbl_name] ...
//	    [RESTRICT | CASCADE]
func (p *Parser) parseDropTableStmt(start int, temporary bool) (*nodes.DropTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.DropTableStmt{
		Loc:       nodes.Loc{Start: start},
		Temporary: temporary,
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Completion: after DROP TABLE [IF EXISTS], offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Table list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// RESTRICT | CASCADE
	if _, ok := p.match(kwRESTRICT); ok {
		stmt.Restrict = true
	} else if _, ok := p.match(kwCASCADE); ok {
		stmt.Cascade = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropIndexStmt parses a DROP INDEX statement.
//
// Ref: https://mariadb.com/kb/en/drop-index/
//
//	DROP INDEX [IF EXISTS] index_name ON tbl_name [WAIT n | NOWAIT]
//
// MariaDB reads IF EXISTS here, which MySQL does not, and reads neither of the
// clauses MySQL closes the statement with: mariadb:11.8.8 answers 0 for DROP
// INDEX IF EXISTS i ON t and 1064 for ALGORITHM and LOCK in every spelling,
// DEFAULT and INPLACE alike (read on 2026-11-04). It closes the statement with
// [WAIT n | NOWAIT] instead.
func (p *Parser) parseDropIndexStmt(start int) (*nodes.DropIndexStmt, error) {
	p.advance() // consume INDEX

	// Completion: after DROP INDEX, offer index_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("index_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.DropIndexStmt{Loc: nodes.Loc{Start: start}}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Index name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON table
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	// Completion: after DROP INDEX idx ON, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	if stmt.NoWait, stmt.HasWait, stmt.WaitSeconds, err = p.parseIndexLockWait(); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropViewStmt parses a DROP VIEW statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-view.html
//
//	DROP VIEW [IF EXISTS] view_name [, view_name] ... [RESTRICT | CASCADE]
func (p *Parser) parseDropViewStmt(start int) (*nodes.DropViewStmt, error) {
	p.advance() // consume VIEW

	stmt := &nodes.DropViewStmt{Loc: nodes.Loc{Start: start}}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Completion: after DROP VIEW [IF EXISTS], offer view_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("view_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// View list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Views = append(stmt.Views, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// RESTRICT | CASCADE
	if _, ok := p.match(kwRESTRICT); ok {
		stmt.Restrict = true
	} else if _, ok := p.match(kwCASCADE); ok {
		stmt.Cascade = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/truncate-table.html
//
//	TRUNCATE [TABLE] tbl_name
func (p *Parser) parseTruncateStmt() (*nodes.TruncateStmt, error) {
	start := p.pos()
	p.advance() // consume TRUNCATE

	// Optional TABLE keyword
	p.match(kwTABLE)

	// Completion: after TRUNCATE [TABLE], offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.TruncateStmt{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Tables: []*nodes.TableRef{ref},
	}
	return stmt, nil
}

// parseRenameTableStmt parses a RENAME TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/rename-table.html
//
//	RENAME TABLE tbl_name TO new_tbl_name [, tbl_name2 TO new_tbl_name2] ...
func (p *Parser) parseRenameTableStmt(start int) (*nodes.RenameTableStmt, error) {
	// Already consumed RENAME
	p.match(kwTABLE)

	// Completion: after RENAME TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.RenameTableStmt{Loc: nodes.Loc{Start: start}}

	for {
		pairStart := p.pos()
		old, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		// Completion: after RENAME TABLE t TO, identifier context (new name).
		p.checkCursor()
		if p.collectMode() {
			// No specific candidates — user defines a new name.
			return nil, &ParseError{Message: "collecting"}
		}
		newRef, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Pairs = append(stmt.Pairs, &nodes.RenameTablePair{
			Loc: nodes.Loc{Start: pairStart, End: p.prev.End},
			Old: old,
			New: newRef,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
