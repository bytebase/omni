package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// parseUpdate parses an UPDATE statement:
//
//	UPDATE [keyspace.]table [USING TTL n AND TIMESTAMP m] SET assignments WHERE relationElements [IF EXISTS | IF ifConditionList]
func (p *Parser) parseUpdate() (*ast.UpdateStmt, error) {
	start := p.curLoc()
	if err := p.expectKeyword(tokUPDATE); err != nil {
		return nil, err
	}

	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.UpdateStmt{Table: table}

	// Optional USING clause (before SET)
	using, err := p.parseUsingClause()
	if err != nil {
		return nil, err
	}
	stmt.Using = using

	// SET assignments
	if err := p.expectKeyword(tokSET); err != nil {
		return nil, err
	}

	assignments, err := p.parseAssignments()
	if err != nil {
		return nil, err
	}
	stmt.Assignments = assignments

	// WHERE clause
	where, err := p.parseWhereClause()
	if err != nil {
		return nil, err
	}
	stmt.Where = where

	// Optional IF EXISTS or IF conditions
	if p.cur.Type == tokIF {
		if p.peekNext().Type == tokEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		} else {
			conds, err := p.parseIfConditions()
			if err != nil {
				return nil, err
			}
			stmt.IfConditions = conds
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseAssignments parses columnOperation (',' columnOperation)*
// (Parser.g:552). Setting a column twice in ways that do not combine is
// refused as Parser.g:144 addRawUpdate refuses it: a column set with = may
// not be updated again, nor a field set twice.
func (p *Parser) parseAssignments() ([]*ast.AssignmentElement, error) {
	var assignments []*ast.AssignmentElement
	seen := map[string]*columnUpdates{}
	for {
		startTok := p.cur
		elem, err := p.parseAssignmentElement()
		if err != nil {
			return nil, err
		}
		u := rawUpdateOf(elem)
		prev := seen[u.column]
		if prev == nil {
			prev = &columnUpdates{fields: map[string]bool{}}
			seen[u.column] = prev
		}
		if !prev.admits(u) {
			return nil, p.errorAt(startTok, "multiple incompatible setting of column %s", u.column)
		}
		prev.add(u)
		assignments = append(assignments, elem)
		if !p.match(tokCOMMA) {
			return assignments, nil
		}
	}
}

// rawUpdate is the kind of update an assignment makes to a column, as
// src/java/org/apache/cassandra/cql3/Operation.java tells them apart.
type rawUpdate struct {
	column string
	kind   int    // updateSetValue, updateSetElement, updateSetField, updateOther
	field  string // for updateSetField
}

const (
	updateSetValue = iota
	updateSetElement
	updateSetField
	updateOther // addition, subtraction, prepend
)

func rawUpdateOf(a *ast.AssignmentElement) rawUpdate {
	switch t := a.Target.(type) {
	case *ast.IndexAccess:
		return rawUpdate{column: columnKey(t.Collection.(*ast.Identifier)), kind: updateSetElement}
	case *ast.DotAccess:
		return rawUpdate{column: columnKey(t.Object.(*ast.Identifier)), kind: updateSetField, field: columnKey(t.Field)}
	}
	kind := updateOther
	if a.Operator == "=" {
		kind = updateSetValue
	}
	return rawUpdate{column: columnKey(a.Target.(*ast.Identifier)), kind: kind}
}

// columnUpdates sums up the updates made so far to one column, so that a
// later one is checked against all of them at once.
type columnUpdates struct {
	any      bool            // some update was made
	setValue bool            // one of them set the whole column
	fields   map[string]bool // the fields set
}

// admits reports whether o is compatible with every earlier update of the
// column, as Operation.java's isCompatibleWith decides for each pair: a
// column set whole takes no other update, and a field is set at most once.
func (c *columnUpdates) admits(o rawUpdate) bool {
	switch {
	case c.setValue:
		return false
	case o.kind == updateSetValue:
		return !c.any
	case o.kind == updateSetField:
		return !c.fields[o.field]
	}
	return true
}

func (c *columnUpdates) add(o rawUpdate) {
	c.any = true
	switch o.kind {
	case updateSetValue:
		c.setValue = true
	case updateSetField:
		c.fields[o.field] = true
	}
}

// columnKey is the name Cassandra knows a column by: an unquoted name in
// lower case, a quoted one as written.
func columnKey(id *ast.Identifier) string {
	if id.Quoted {
		return id.Name
	}
	return strings.ToLower(id.Name)
}

// parseAssignmentElement parses a columnOperation (Parser.g:1668):
//
//	cident '=' term                     set the column
//	cident '=' term '+' cident          prepend: X = <value> + X
//	cident '=' cident ('+' | '-') term  X = X + <value>, X = X - <value>
//	cident '=' cident INTEGER           X = X 3, as Cassandra lexes X = X -3
//	cident ('+=' | '-=') term
//	cident '[' term ']' '=' term
//	cident '.' fident '=' term
//
// The column on both sides must be the same one (Parser.g:1675-1700).
func (p *Parser) parseAssignmentElement() (*ast.AssignmentElement, error) {
	start := p.curLoc()

	if !isNameToken(p.cur) {
		return nil, p.errorf("expected column name, got %s", p.tokenDesc())
	}
	target, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	elem := &ast.AssignmentElement{Target: target, Operator: "="}
	switch p.cur.Type {
	case tokPLUSEQ, tokMINUSEQ:
		elem.Operator = p.cur.Str[:1]
		p.advance()
		if elem.Value, err = p.parseTerm(); err != nil {
			return nil, err
		}
		elem.Loc = p.makeLoc(start)
		return elem, nil

	case tokLBRACK:
		p.advance()
		idx, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		elem.Target = &ast.IndexAccess{Collection: target, Index: idx, Loc: p.makeLoc(start)}

	case tokDOT:
		p.advance()
		if !isNameToken(p.cur) {
			return nil, p.errorf("expected field name, got %s", p.tokenDesc())
		}
		field, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		elem.Target = &ast.DotAccess{Object: target, Field: field, Loc: p.makeLoc(start)}
	}

	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}

	if elem.Target != ast.ExprNode(target) {
		// col[key] = term, col.field = term (Parser.g:1714, 1721).
		if elem.Value, err = p.parseTerm(); err != nil {
			return nil, err
		}
		elem.Loc = p.makeLoc(start)
		return elem, nil
	}

	// A name not followed by '(' or '.' is no term: it is the column again.
	if isNameToken(p.cur) {
		if next := p.peekNext(); next.Type != tokLPAREN && next.Type != tokDOT {
			colTok := p.cur
			col, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			switch {
			case isSignedConstant(p.cur, p.peekNext()):
				// X = X -3 is [X, '=', X, INTEGER] to Cassandra's lexer
				// (Parser.g:1694, Lexer.g:337): an addition of the
				// integer -3, which no operator may follow.
				if p.peekNext().Type != tokINTEGER {
					return nil, p.errorf("expected '+' or '-' after column %s, got %s", col.Name, p.tokenDesc())
				}
				elem.Operator = "+"
				minus := p.curLoc()
				p.advance() // -
				lit, err := p.parseConstant()
				if err != nil {
					return nil, err
				}
				elem.Value = negate(lit, minus, p.prevEnd())
			case p.cur.Type == tokPLUS, p.cur.Type == tokMINUS:
				elem.Operator = p.cur.Str
				p.advance()
				if elem.Value, err = p.parseTerm(); err != nil {
					return nil, err
				}
			case p.cur.Type == tokINTEGER:
				elem.Operator = "+"
				if elem.Value, err = p.parseConstant(); err != nil {
					return nil, err
				}
			default:
				return nil, p.errorf("expected '+' or '-' after column %s, got %s", col.Name, p.tokenDesc())
			}
			if columnKey(col) != columnKey(target) {
				return nil, p.errorAt(colTok, "only expressions of the form X = X %s <value> are supported", elem.Operator)
			}
			elem.Loc = p.makeLoc(start)
			return elem, nil
		}
	}

	if elem.Value, err = p.parseTerm(); err != nil {
		return nil, err
	}
	if p.cur.Type == tokPLUS {
		// X = <value> + X: parseTerm stopped before this '+', as no term
		// follows it.
		p.advance()
		colTok := p.cur
		if !isNameToken(p.cur) {
			return nil, p.errorf("expected column name, got %s", p.tokenDesc())
		}
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if columnKey(col) != columnKey(target) {
			return nil, p.errorAt(colTok, "only expressions of the form X = <value> + X are supported")
		}
		elem.Operator = "+"
	}
	elem.Loc = p.makeLoc(start)
	return elem, nil
}

// parseIfConditions parses IF columnCondition (AND columnCondition)*
// (Parser.g:566).
func (p *Parser) parseIfConditions() ([]*ast.IfCondition, error) {
	if err := p.expectKeyword(tokIF); err != nil {
		return nil, err
	}

	var conditions []*ast.IfCondition
	for {
		cond, err := p.parseIfCondition()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, cond)
		if !p.match(tokAND) {
			break
		}
	}
	return conditions, nil
}

// parseIfCondition parses a columnCondition (Parser.g:1728):
//
//	cident relationType term
//	cident CONTAINS KEY? term
//	cident IN (singleColumnInValues | inMarker)
//	cident '[' term ']' (relationType term | IN ...)
//	cident '.' fident (relationType term | IN ...)
func (p *Parser) parseIfCondition() (*ast.IfCondition, error) {
	start := p.curLoc()

	if !isNameToken(p.cur) {
		return nil, p.errorf("expected column name in IF condition, got %s", p.tokenDesc())
	}
	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cond := &ast.IfCondition{Column: col}

	switch p.cur.Type {
	case tokCONTAINS:
		p.advance()
		cond.Op = "CONTAINS"
		if p.match(tokKEY) {
			cond.Op = "CONTAINS KEY"
		}
		if cond.Value, err = p.parseTerm(); err != nil {
			return nil, err
		}
		cond.Loc = p.makeLoc(start)
		return cond, nil
	case tokLBRACK:
		p.advance()
		if cond.Element, err = p.parseTerm(); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
	case tokDOT:
		p.advance()
		if !isNameToken(p.cur) {
			return nil, p.errorf("expected field name, got %s", p.tokenDesc())
		}
		if cond.Field, err = p.parseIdentifier(); err != nil {
			return nil, err
		}
	}

	switch {
	case p.cur.Type == tokIN:
		p.advance()
		cond.Op = "IN"
		if cond.InValues, cond.InMarker, err = p.parseInValues(); err != nil {
			return nil, err
		}
	case isRelationType(p.cur.Type):
		cond.Op = p.cur.Str
		p.advance()
		if cond.Value, err = p.parseTerm(); err != nil {
			return nil, err
		}
	default:
		return nil, p.errorf("expected comparison operator, IN, or CONTAINS in IF condition, got %s", p.tokenDesc())
	}
	cond.Loc = p.makeLoc(start)
	return cond, nil
}
