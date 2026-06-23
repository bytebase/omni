package parser

import (
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

// parseAssignments parses: assignmentElement (',' assignmentElement)*
func (p *Parser) parseAssignments() ([]*ast.AssignmentElement, error) {
	var assignments []*ast.AssignmentElement
	first, err := p.parseAssignmentElement()
	if err != nil {
		return nil, err
	}
	assignments = append(assignments, first)
	for p.match(tokCOMMA) {
		elem, err := p.parseAssignmentElement()
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, elem)
	}
	return assignments, nil
}

// parseAssignmentElement parses a single assignment in the SET clause.
//
// Forms:
//
//	IDENT '=' expression
//	IDENT '=' IDENT ('+' | '-') expression
//	IDENT '=' expression ('+' | '-') IDENT
//	IDENT '[' expression ']' '=' expression
func (p *Parser) parseAssignmentElement() (*ast.AssignmentElement, error) {
	start := p.curLoc()

	target, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for field access (col.field) or index access (col[idx])
	var assignTarget ast.ExprNode = target
	if p.cur.Type == tokDOT {
		p.advance()
		field, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		assignTarget = &ast.DotAccess{
			Object: target,
			Field:  field,
			Loc:    p.makeLoc(start),
		}
	} else if p.cur.Type == tokLBRACK {
		p.advance() // [
		idx, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		assignTarget = &ast.IndexAccess{
			Collection: target,
			Index:      idx,
			Loc:        p.makeLoc(start),
		}
	}

	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}

	// Parse the right-hand side. We need to detect patterns like:
	//   IDENT (+|-) expression   => counter/collection increment
	//   expression (+|-) IDENT   => collection prepend
	//   plain expression         => simple assignment
	rhs, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Check for arithmetic operator after the first RHS expression.
	if p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		op := p.cur.Str
		p.advance()
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// Determine if this is col = col + val or col = val + col.
		// If lhs is an identifier matching the target, it's col = col + val (operator is += or -=).
		// If rhs is an identifier matching the target, it's col = val + col (also += but collection prepend).
		// In either case, we represent it using the Operator field and the non-target side as Value.
		if ident, ok := rhs.(*ast.Identifier); ok && matchesTarget(ident, target) {
			// col = col + val  =>  Operator: "+"
			return &ast.AssignmentElement{
				Target:   assignTarget,
				Value:    right,
				Operator: op,
				Loc:      p.makeLoc(start),
			}, nil
		}
		if ident, ok := right.(*ast.Identifier); ok && matchesTarget(ident, target) {
			// col = val + col  =>  Operator: "+"
			return &ast.AssignmentElement{
				Target:   assignTarget,
				Value:    rhs,
				Operator: op,
				Loc:      p.makeLoc(start),
			}, nil
		}
		// General arithmetic expression; wrap as binary expression value with simple "=".
		return &ast.AssignmentElement{
			Target:   assignTarget,
			Value:    &ast.BinaryExpr{Left: rhs, Op: op, Right: right, Loc: p.makeLoc(rhs.GetLoc().Start)},
			Operator: "=",
			Loc:      p.makeLoc(start),
		}, nil
	}

	return &ast.AssignmentElement{
		Target:   assignTarget,
		Value:    rhs,
		Operator: "=",
		Loc:      p.makeLoc(start),
	}, nil
}

// matchesTarget checks if an identifier matches the assignment target name.
func matchesTarget(ident *ast.Identifier, target *ast.Identifier) bool {
	return ident.Name == target.Name && ident.Quoted == target.Quoted
}

// parseIfConditions parses IF col op val [AND col op val ...].
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

// parseIfCondition parses a single LWT condition:
//
//	col op value | col IN (values) | col CONTAINS [KEY] value
func (p *Parser) parseIfCondition() (*ast.IfCondition, error) {
	start := p.curLoc()

	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// IN condition
	if p.cur.Type == tokIN {
		p.advance()
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		var values []ast.ExprNode
		if p.cur.Type != tokRPAREN {
			values, err = p.parseExpressionList()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.IfCondition{
			Column:   col,
			Op:       "IN",
			InValues: values,
			Loc:      p.makeLoc(start),
		}, nil
	}

	// CONTAINS [KEY] condition
	if p.cur.Type == tokCONTAINS {
		p.advance()
		op := "CONTAINS"
		if p.cur.Type == tokKEY {
			op = "CONTAINS KEY"
			p.advance()
		}
		val, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &ast.IfCondition{
			Column: col,
			Op:     op,
			Value:  val,
			Loc:    p.makeLoc(start),
		}, nil
	}

	// Comparison operator
	op := p.cur.Str
	if !p.match(tokEQ, tokLT, tokGT, tokLTE, tokGTE, tokNE) {
		return nil, p.errorf("expected comparison operator, IN, or CONTAINS in IF condition, got %s", p.tokenDesc())
	}

	val, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.IfCondition{
		Column: col,
		Op:     op,
		Value:  val,
		Loc:    p.makeLoc(start),
	}, nil
}
