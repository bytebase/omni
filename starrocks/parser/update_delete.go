package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
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

	// expressionOrDefault takes expression or DEFAULT (StarRocks.g4).
	var val ast.Node
	if p.cur.Kind == kwDEFAULT {
		tok := p.advance()
		val = &ast.Literal{Kind: ast.LitKeyword, Value: "DEFAULT", Loc: tok.Loc}
	} else {
		val, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
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

	// deleteStatement (StarRocks.g4:2197) names no alias for its target:
	// `DELETE FROM a a` and `DELETE FROM a AS b` are syntax errors on the
	// pinned server, while `DELETE a FROM ...` is another statement Doris has
	// and StarRocks does not.

	// partitionNames? (StarRocks.g4:1205)
	if p.atPartitionNames() {
		partitions, temporary, err := p.parsePartitionNames()
		if err != nil {
			return nil, err
		}
		stmt.Partition = partitions
		stmt.TemporaryPartition = temporary
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
		kwINNER, kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwJOIN:
		return true
	}
	return isSelectClauseKeyword(t)
}

// atPartitionNames reports whether cur starts partitionNames
// (StarRocks.g4:2268): TEMPORARY, PARTITION or PARTITIONS.
func (p *Parser) atPartitionNames() bool {
	switch p.cur.Kind {
	case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
		return true
	}
	return false
}

// parsePartitionNames reads partitionNames (StarRocks.g4:2268): a list of
// names or strings, one name or string, or keyPartitions, each name as
// written for a key partition, and whether TEMPORARY came first.
func (p *Parser) parsePartitionNames() ([]string, bool, error) {
	var names ast.AlterTableAction
	if err := p.parseAlterPartitionNames(&names); err != nil {
		return nil, false, err
	}
	if names.PartitionName != "" {
		return []string{names.PartitionName}, names.Temporary, nil
	}
	return names.PartitionList, names.Temporary, nil
}
