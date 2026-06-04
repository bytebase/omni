package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-dml DAG node (with insert.go and merge.go):
// it implements Trino's DELETE and UPDATE data-manipulation statements as
// hand-written recursive-descent parsers over the token stream.
//
// The statement structs live here in package parser and satisfy ast.Node with
// the placeholder tag ast.T_Invalid; see the insert.go file header for why DML
// statement nodes use T_Invalid rather than a dedicated trino/ast/nodetags.go
// tag.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | DELETE_ FROM_ qualifiedName (WHERE_ booleanExpression)?                   # delete
//	    | UPDATE_ qualifiedName SET_ updateAssignment (COMMA_ updateAssignment)*
//	          (WHERE_ where = booleanExpression)?                                   # update
//	    ;
//	updateAssignment : identifier EQ_ expression ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Oracle-confirmed facts (Trino 481) baked in:
//
//	D-DEL1 (DELETE requires FROM). `DELETE orders` is a SYNTAX_ERROR; the FROM
//	   keyword is mandatory. The WHERE clause is optional (`DELETE FROM orders`
//	   deletes all rows). A dangling `DELETE FROM orders WHERE` is rejected.
//	D-UPD1 (UPDATE assignment target is a single identifier, value is an
//	   expression). `UPDATE t SET (a, b) = (1, 2)` — a row assignment — is a
//	   SYNTAX_ERROR in Trino 481: each assignment is exactly `identifier = expr`.
//	D-UPD2 (UPDATE requires at least one assignment after SET). `UPDATE t SET`
//	   and `UPDATE t a = 1` (missing SET) are SYNTAX_ERRORs.
//	D-DML-NAME (target is at most catalog.schema.table). A 4-part target name
//	   (`a.b.c.d`) is a SYNTAX_ERROR for both DELETE and UPDATE.
//
// The WHERE / value / assignment expressions use the shared expression entry
// point parseExpr (== the grammar's `expression` / `booleanExpression`), so a
// subquery predicate (`DELETE … WHERE k IN (SELECT …)`) and a scalar-subquery
// assignment value (`UPDATE … SET m = (SELECT …)`) parse via the expression
// node's SubqueryExpr placeholder, exactly as in SELECT predicates.
//
// FLAGGED divergence (ledger #5): `@ branch_name` after the target table —
// `DELETE FROM orders @ audit`, `UPDATE purchases @ audit SET …` — is accepted
// by Trino 481 but rejected by omni because the lexer emits no '@' token. See
// the insert.go file header.

// ---------------------------------------------------------------------------
// DELETE
// ---------------------------------------------------------------------------

// DeleteStmt is a `DELETE FROM qualifiedName [WHERE booleanExpression]`
// statement (the delete alternative). Target is the table to delete from (1-3
// parts). Where is the optional row filter (nil when absent — meaning delete
// all rows).
type DeleteStmt struct {
	Target *ast.QualifiedName
	Where  Expr // nil when there is no WHERE clause
	Loc    ast.Loc
}

// Tag implements ast.Node.
func (n *DeleteStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Span returns the source byte range.
func (n *DeleteStmt) Span() ast.Loc { return n.Loc }

// ---------------------------------------------------------------------------
// UPDATE
// ---------------------------------------------------------------------------

// UpdateAssignment is one `identifier = expression` SET assignment (the
// updateAssignment rule). Column is the single target column name; Value is the
// new-value expression (which may be a scalar subquery placeholder).
type UpdateAssignment struct {
	Column *ast.Identifier
	Value  Expr
	Loc    ast.Loc
}

// UpdateStmt is an `UPDATE qualifiedName SET assignment (, assignment)*
// [WHERE booleanExpression]` statement (the update alternative). Target is the
// table to update (1-3 parts). Assignments is the non-empty SET list. Where is
// the optional row filter (nil when absent).
type UpdateStmt struct {
	Target      *ast.QualifiedName
	Assignments []UpdateAssignment
	Where       Expr // nil when there is no WHERE clause
	Loc         ast.Loc
}

// Tag implements ast.Node.
func (n *UpdateStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Span returns the source byte range.
func (n *UpdateStmt) Span() ast.Loc { return n.Loc }

// Compile-time assertions that the statement types satisfy ast.Node.
var (
	_ ast.Node = (*DeleteStmt)(nil)
	_ ast.Node = (*UpdateStmt)(nil)
)

// ---------------------------------------------------------------------------
// Parsers
// ---------------------------------------------------------------------------

// parseDeleteStmt parses `DELETE FROM qualifiedName [WHERE booleanExpression]`.
// On entry cur is the DELETE keyword. D-DEL1: FROM is mandatory; the WHERE
// clause is optional.
func (p *Parser) parseDeleteStmt() (ast.Node, error) {
	deleteTok := p.advance() // consume DELETE
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// D-DML-NAME: target is at most catalog.schema.table.
	target, err := p.parseBoundedQualifiedName(3, "delete target")
	if err != nil {
		return nil, err
	}

	stmt := &DeleteStmt{Target: target, Loc: ast.Loc{Start: deleteTok.Loc.Start, End: target.Loc.End}}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc.End = where.Span().End
	}

	return stmt, nil
}

// parseUpdateStmt parses
//
//	UPDATE qualifiedName SET updateAssignment (, updateAssignment)*
//	    [WHERE booleanExpression]
//
// On entry cur is the UPDATE keyword. D-UPD2: SET and at least one assignment
// are mandatory; the WHERE clause is optional.
func (p *Parser) parseUpdateStmt() (ast.Node, error) {
	updateTok := p.advance() // consume UPDATE

	// D-DML-NAME: target is at most catalog.schema.table.
	target, err := p.parseBoundedQualifiedName(3, "update target")
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	stmt := &UpdateStmt{Target: target, Loc: ast.Loc{Start: updateTok.Loc.Start, End: target.Loc.End}}

	// One or more comma-separated assignments (the list is non-empty: D-UPD2).
	for {
		asgn, err := p.parseUpdateAssignment()
		if err != nil {
			return nil, err
		}
		stmt.Assignments = append(stmt.Assignments, asgn)
		stmt.Loc.End = asgn.Loc.End
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc.End = where.Span().End
	}

	return stmt, nil
}

// parseUpdateAssignment parses one `identifier = expression` (the
// updateAssignment rule). D-UPD1: the target is a single bare identifier — a
// parenthesized row target like `(a, b) = (1, 2)` is rejected because '(' is
// not an identifier start, surfacing as a syntax error at the '('.
func (p *Parser) parseUpdateAssignment() (UpdateAssignment, error) {
	col, err := p.parseIdentifier()
	if err != nil {
		return UpdateAssignment{}, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return UpdateAssignment{}, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return UpdateAssignment{}, err
	}
	return UpdateAssignment{
		Column: col,
		Value:  value,
		Loc:    ast.Loc{Start: col.Loc.Start, End: value.Span().End},
	}, nil
}
