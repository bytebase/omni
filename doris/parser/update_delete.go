package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// UPDATE statement parser (T4.2)
// ---------------------------------------------------------------------------

// parseUpdateStmt parses an UPDATE statement:
//
//	UPDATE table [AS alias]
//	    SET col1 = expr1 [, col2 = expr2 ...]
//	    [FROM table_refs]
//	    [WHERE condition]
//
// The UPDATE keyword has already been consumed by the caller.
func (p *Parser) parseUpdateStmt(updateStart int) (*ast.UpdateStmt, error) {
	stmt := &ast.UpdateStmt{
		Loc: ast.Loc{Start: updateStart},
	}

	// table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// optional AS alias (or bare alias) — but not SET/FROM/WHERE
	stmt.TargetAlias = p.parseDMLTargetAlias()

	// SET clause (required)
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	assignments, err := p.parseAssignmentList()
	if err != nil {
		return nil, err
	}
	stmt.Assignments = assignments

	// optional FROM clause
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// optional WHERE clause
	if p.cur.Kind == kwWHERE {
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

// parseAssignmentList parses a comma-separated list of col = expr assignments.
func (p *Parser) parseAssignmentList() ([]*ast.Assignment, error) {
	var assignments []*ast.Assignment

	a, err := p.parseAssignment()
	if err != nil {
		return nil, err
	}
	assignments = append(assignments, a)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		a, err = p.parseAssignment()
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}

	return assignments, nil
}

// parseAssignment parses a single col = expr assignment.
// Supports plain column names and qualified column names (t.col).
func (p *Parser) parseAssignment() (*ast.Assignment, error) {
	startLoc := p.cur.Loc

	col, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{
		Column: col,
		Value:  val,
		Loc:    ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// DELETE statement parser (T4.2)
// ---------------------------------------------------------------------------

// parseDeleteStmt parses a DELETE statement:
//
//	DELETE FROM table [AS alias]
//	    [PARTITION(p1 [, p2 ...])]
//	    [USING table_refs]
//	    [WHERE condition]
//
// The DELETE keyword has already been consumed by the caller.
func (p *Parser) parseDeleteStmt(deleteStart int) (*ast.DeleteStmt, error) {
	stmt := &ast.DeleteStmt{
		Loc: ast.Loc{Start: deleteStart},
	}

	// FROM keyword (required)
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// optional AS alias (or bare alias) — but not PARTITION/PARTITIONS/USING/WHERE
	stmt.TargetAlias = p.parseDMLTargetAlias()

	// optional PARTITION(p1, p2, ...) or PARTITIONS (p1, p2, ...)
	if p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS {
		p.advance() // consume PARTITION / PARTITIONS
		partitions, err := p.parsePartitionNameList()
		if err != nil {
			return nil, err
		}
		stmt.Partition = partitions
	}

	// optional USING clause
	if p.cur.Kind == kwUSING {
		p.advance() // consume USING
		using, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.Using = using
	}

	// optional WHERE clause
	if p.cur.Kind == kwWHERE {
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

// parseDMLTargetAlias parses an optional alias for a DELETE/UPDATE target table.
// It uses the same logic as parseOptionalAlias but additionally excludes
// PARTITION, PARTITIONS, USING, and WHERE which can appear as non-reserved
// keywords but are structural keywords in DML clauses.
func (p *Parser) parseDMLTargetAlias() string {
	// Explicit: AS alias
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return ""
		}
		return name
	}

	// Implicit alias: skip DML structural keywords
	if isDMLClauseKeyword(p.cur.Kind) {
		return ""
	}

	if p.isAliasIdentToken() {
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return ""
		}
		return name
	}

	return ""
}

// isDMLClauseKeyword returns true for keywords that start DML clauses and
// should not be consumed as implicit aliases after a table name.
func isDMLClauseKeyword(t int) bool {
	switch t {
	case kwPARTITION, kwPARTITIONS, kwUSING, kwWHERE, kwSET, kwFROM,
		kwINNER, kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwJOIN, kwNATURAL:
		return true
	}
	return isSelectClauseKeyword(t)
}

// parsePartitionNameList parses (p1 [, p2 ...]) — the partition name list.
// The opening '(' is optional in Doris for single-partition DELETE syntax
// (PARTITION p1) but PARTITIONS always uses '(p1, p2, ...)'. We handle both.
func (p *Parser) parsePartitionNameList() ([]string, error) {
	var names []string

	if p.cur.Kind == int('(') {
		// Parenthesized list: (p1, p2, ...)
		p.advance() // consume '('

		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)

		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			name, _, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			names = append(names, name)
		}

		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	} else {
		// Bare single name: PARTITION p1
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}
