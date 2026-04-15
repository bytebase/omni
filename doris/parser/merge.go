package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// MERGE INTO statement parser (T4.3)
// ---------------------------------------------------------------------------

// parseMergeStmt parses a MERGE INTO statement.
//
// Syntax:
//
//	MERGE INTO target [AS alias]
//	  USING source [AS alias]
//	  ON condition
//	  WHEN [NOT] MATCHED [AND condition] THEN action
//	  [WHEN ...]
//
// The MERGE keyword has already been consumed by the caller.
func (p *Parser) parseMergeStmt(mergeLoc ast.Loc) (*ast.MergeStmt, error) {
	// Consume INTO
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	// Parse target table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.MergeStmt{
		Target: target,
		Loc:    ast.Loc{Start: mergeLoc.Start},
	}

	// Optional target alias: [AS] alias
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		alias, _, err := p.parseAliasIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.TargetAlias = alias
	} else if p.isAliasIdentToken() {
		alias, _, err := p.parseAliasIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.TargetAlias = alias
	}

	// USING source
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}

	// Source: either (subquery) or table_name
	source, sourceAlias, err := p.parseMergeSource()
	if err != nil {
		return nil, err
	}
	stmt.Source = source
	stmt.SourceAlias = sourceAlias

	// ON condition
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	onExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.On = onExpr

	// One or more WHEN clauses
	for p.cur.Kind == kwWHEN {
		clause, err := p.parseMergeClause()
		if err != nil {
			return nil, err
		}
		stmt.Clauses = append(stmt.Clauses, clause)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseMergeSource parses the USING source (table reference or subquery)
// and optional alias. Returns (source node, alias string, error).
func (p *Parser) parseMergeSource() (ast.Node, string, error) {
	// Parenthesized subquery: (SELECT ...)
	if p.cur.Kind == int('(') {
		openTok := p.advance() // consume '('
		subq, err := p.parseSubqueryPlaceholder(openTok.Loc.Start)
		if err != nil {
			return nil, "", err
		}
		alias := p.parseOptionalAlias()
		return subq, alias, nil
	}

	// Table reference
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, "", err
	}
	ref := &ast.TableRef{
		Name: name,
		Loc:  name.Loc,
	}
	alias := p.parseOptionalAlias()
	if alias != "" {
		ref.Alias = alias
		ref.Loc.End = p.prev.Loc.End
	}
	return ref, alias, nil
}

// parseMergeClause parses one WHEN [NOT] MATCHED [AND condition] THEN action clause.
// Consumes the leading WHEN token.
func (p *Parser) parseMergeClause() (*ast.MergeClause, error) {
	startTok, err := p.expect(kwWHEN)
	if err != nil {
		return nil, err
	}

	clause := &ast.MergeClause{
		Loc: ast.Loc{Start: startTok.Loc.Start},
	}

	// Optional NOT
	if p.cur.Kind == kwNOT {
		p.advance() // consume NOT
		clause.NotMatched = true
	}

	// MATCHED
	if _, err := p.expect(kwMATCHED); err != nil {
		return nil, err
	}

	// Optional AND condition
	if p.cur.Kind == kwAND {
		p.advance() // consume AND
		andExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		clause.And = andExpr
	}

	// THEN
	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}

	// Action: DELETE | UPDATE SET ... | INSERT ... | DO NOTHING
	if err := p.parseMergeAction(clause); err != nil {
		return nil, err
	}

	clause.Loc.End = p.prev.Loc.End
	return clause, nil
}

// parseMergeAction parses the action part after THEN and fills in the clause fields.
func (p *Parser) parseMergeAction(clause *ast.MergeClause) error {
	switch p.cur.Kind {
	case kwDELETE:
		p.advance() // consume DELETE
		clause.Action = ast.MergeActionDelete
		return nil

	case kwUPDATE:
		return p.parseMergeUpdateAction(clause)

	case kwINSERT:
		return p.parseMergeInsertAction(clause)

	case kwDO:
		p.advance() // consume DO
		// NOTHING — expect an identifier token with value "nothing"
		// NOTHING is not a keyword in Doris, so it is a plain identifier.
		if p.cur.Kind != tokIdent || p.cur.Str != "nothing" {
			// Tolerate any casing
			if p.cur.Kind != tokIdent {
				return p.syntaxErrorAtCur()
			}
		}
		p.advance() // consume NOTHING
		clause.Action = ast.MergeActionDoNothing
		return nil

	default:
		return p.syntaxErrorAtCur()
	}
}

// parseMergeUpdateAction parses UPDATE SET col = expr [, ...] or UPDATE SET *.
// The UPDATE keyword has NOT yet been consumed.
func (p *Parser) parseMergeUpdateAction(clause *ast.MergeClause) error {
	p.advance() // consume UPDATE

	if _, err := p.expect(kwSET); err != nil {
		return err
	}

	clause.Action = ast.MergeActionUpdate

	// UPDATE SET * — wildcard update
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		clause.UpdateAll = true
		return nil
	}

	// UPDATE SET col = expr [, col = expr, ...]
	assignments, err := p.parseMergeAssignments()
	if err != nil {
		return err
	}
	clause.Assignments = assignments
	return nil
}

// parseMergeAssignments parses one or more col = expr pairs separated by commas.
func (p *Parser) parseMergeAssignments() ([]*ast.Assignment, error) {
	var assignments []*ast.Assignment

	a, err := p.parseMergeAssignment()
	if err != nil {
		return nil, err
	}
	assignments = append(assignments, a)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		a, err = p.parseMergeAssignment()
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}

	return assignments, nil
}

// parseMergeAssignment parses one col = expr assignment.
// The column reference may be qualified (e.g. t.c), but we store just the column name.
func (p *Parser) parseMergeAssignment() (*ast.Assignment, error) {
	startLoc := p.cur.Loc

	// Parse the left-hand side: may be qualified (table.col)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// Expect '='
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}

	// Parse value expression
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{
		Column: name,
		Value:  val,
		Loc:    ast.Loc{Start: startLoc.Start, End: ast.NodeLoc(val).End},
	}, nil
}

// parseMergeInsertAction parses INSERT [(cols)] VALUES (exprs) or INSERT *.
// The INSERT keyword has NOT yet been consumed.
func (p *Parser) parseMergeInsertAction(clause *ast.MergeClause) error {
	p.advance() // consume INSERT

	clause.Action = ast.MergeActionInsert

	// INSERT * — wildcard insert
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		clause.InsertAll = true
		return nil
	}

	// Optional column list: (col1, col2, ...)
	if p.cur.Kind == int('(') {
		// Check if next token could be a column name — if so, it's a column list.
		// (Otherwise it might be VALUES directly without column list.)
		cols, err := p.parseMergeInsertColumns()
		if err != nil {
			return err
		}
		clause.Columns = cols
	}

	// VALUES (expr, ...)
	if _, err := p.expect(kwVALUES); err != nil {
		return err
	}

	if _, err := p.expect(int('(')); err != nil {
		return err
	}

	var vals []ast.Node
	// Allow empty VALUES () ?  In practice Doris requires at least one value.
	if p.cur.Kind != int(')') {
		val, err := p.parseExpr()
		if err != nil {
			return err
		}
		vals = append(vals, val)

		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			val, err = p.parseExpr()
			if err != nil {
				return err
			}
			vals = append(vals, val)
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return err
	}

	clause.Values = vals
	return nil
}

// parseMergeInsertColumns parses (col1, col2, ...) — the optional column list
// in a MERGE INSERT clause.
func (p *Parser) parseMergeInsertColumns() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var cols []string
	col, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cols = append(cols, col)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		col, _, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return cols, nil
}
