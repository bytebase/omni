package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseCreateLoadableFunction parses a CREATE [AGGREGATE] FUNCTION for loadable UDFs.
// The caller has already consumed CREATE (and AGGREGATE if present).
// The current token is FUNCTION.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-function-loadable.html
//
//	CREATE [AGGREGATE] FUNCTION [IF NOT EXISTS] function_name
//	    RETURNS {STRING|INTEGER|REAL|DECIMAL}
//	    SONAME shared_library_name
func (p *Parser) parseCreateLoadableFunction(isAggregate bool) (*nodes.CreateFunctionStmt, error) {
	start := p.pos()
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc:         nodes.Loc{Start: start},
		IsAggregate: isAggregate,
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// Function name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	return p.finishLoadableFunction(stmt)
}

// finishLoadableFunction completes parsing a loadable UDF after the function name
// has been consumed. Parses: RETURNS {STRING|INTEGER|REAL|DECIMAL} SONAME 'lib'
func (p *Parser) finishLoadableFunction(stmt *nodes.CreateFunctionStmt) (*nodes.CreateFunctionStmt, error) {
	// RETURNS {STRING|INTEGER|REAL|DECIMAL}
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	retType, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}
	stmt.Returns = &nodes.DataType{
		Loc:  nodes.Loc{Start: p.pos()},
		Name: retType,
	}

	// SONAME 'shared_library_name'
	if _, err := p.expect(kwSONAME); err != nil {
		return nil, err
	}
	if p.cur.Type != tokSCONST {
		return nil, &ParseError{Message: "expected string literal after SONAME", Position: p.cur.Loc}
	}
	stmt.Soname = p.cur.Str
	p.advance()

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCreateFunctionStmt parses a CREATE FUNCTION or CREATE PROCEDURE statement.
// The caller has already consumed CREATE and set isProcedure.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-procedure.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-function.html
//
//	CREATE [DEFINER = user] PROCEDURE sp_name ([proc_parameter[,...]])
//	    [characteristic ...] routine_body
//	CREATE [DEFINER = user] FUNCTION sp_name ([func_parameter[,...]])
//	    RETURNS type_name [characteristic ...] routine_body
func (p *Parser) parseCreateFunctionStmt(isProcedure bool) (*nodes.CreateFunctionStmt, error) {
	start := p.pos()
	p.advance() // consume FUNCTION or PROCEDURE

	// Completion: after CREATE FUNCTION / CREATE PROCEDURE, identifier context.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new name.
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CreateFunctionStmt{
		Loc:         nodes.Loc{Start: start},
		IsProcedure: isProcedure,
	}

	// Optional IF NOT EXISTS (MySQL 8.0.29+, applies to both FUNCTION and PROCEDURE)
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// Function/procedure name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// If this is a function (not procedure) and there's no '(' next, it's a loadable UDF.
	// Loadable UDF: CREATE FUNCTION name RETURNS {STRING|INTEGER|REAL|DECIMAL} SONAME 'lib'
	if !isProcedure && p.cur.Type != '(' {
		return p.finishLoadableFunction(stmt)
	}

	// Parameter list
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside parameter list, offer param direction keywords + type context.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwIN)
		p.addTokenCandidate(kwOUT)
		p.addTokenCandidate(kwINOUT)
		p.addRuleCandidate("type_name")
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type != ')' {
		for {
			param, err := p.parseFuncParam(isProcedure)
			if err != nil {
				return nil, err
			}
			stmt.Params = append(stmt.Params, param)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// RETURNS type (functions only)
	if !isProcedure && p.cur.Type == kwRETURNS {
		p.advance()

		// Completion: after RETURNS, offer type candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("type_name")
			return nil, &ParseError{Message: "collecting"}
		}

		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.Returns = dt
	}

	// Completion: after parameter list (and optional RETURNS), offer characteristics keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwDETERMINISTIC)
		p.addTokenCandidate(kwNO)
		p.addTokenCandidate(kwSQL)
		p.addTokenCandidate(kwCOMMENT)
		p.addTokenCandidate(kwLANGUAGE)
		p.addTokenCandidate(kwRETURNS)
		return nil, &ParseError{Message: "collecting"}
	}

	// Characteristics
	for {
		ch, ok, err := p.parseRoutineCharacteristic()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Characteristics = append(stmt.Characteristics, ch)
	}

	// Routine body — parse via the grammar so the body is disambiguated the
	// same way MySQL's yacc does (statement vs expression context, nested
	// compound statements, labels, DECLAREs, etc.). See
	// docs/plans/2026-04-20-mysql-routine-body-grammar.md.
	//
	// Open the outermost static-validation scope and seed parameters as
	// in-scope variables so SET param_name = ... resolves and same-name
	// DECLAREs in the top BEGIN block conflict (matches MySQL semantics).
	scope := p.pushScope(scopeBlock)
	scope.isFunction = !isProcedure
	for _, fp := range stmt.Params {
		// Synthesize a DeclareVarStmt placeholder to occupy the var slot.
		// Loc points at the parameter; TypeName carries through.
		_ = p.declareVar(fp.Name, &nodes.DeclareVarStmt{
			Loc:      fp.Loc,
			Names:    []string{fp.Name},
			TypeName: fp.TypeName,
		}, fp.Loc.Start)
	}
	defer p.popScope()

	bodyStart := p.pos()
	body, err := p.parseCompoundStmtOrStmt()
	if err != nil {
		return nil, err
	}
	bodyEnd := p.pos()
	stmt.Body = body
	stmt.BodyText = p.inputText(bodyStart, bodyEnd)

	// RETURN coverage: MySQL's CREATE-time check requires at least one
	// RETURN statement to exist somewhere in the function body (ERR 1320
	// "No RETURN found in FUNCTION"). Path analysis is deferred to runtime;
	// SIGNAL/RESIGNAL do not substitute. Loadable UDFs (Soname != "") have
	// no Body to check.
	if !isProcedure && stmt.Soname == "" && body != nil {
		if !containsReturn(body) {
			return nil, &ParseError{
				Message:  "no RETURN found in function body",
				Position: bodyStart,
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseFuncParam parses a function/procedure parameter.
//
//	[IN | OUT | INOUT] param_name type
func (p *Parser) parseFuncParam(isProcedure bool) (*nodes.FuncParam, error) {
	start := p.pos()
	param := &nodes.FuncParam{Loc: nodes.Loc{Start: start}}

	// Optional direction for procedures
	if isProcedure {
		if p.cur.Type == kwIN {
			param.Direction = "IN"
			p.advance()
		} else if p.cur.Type == kwOUT {
			param.Direction = "OUT"
			p.advance()
		} else if p.cur.Type == kwINOUT {
			param.Direction = "INOUT"
			p.advance()
		}
	}

	// Parameter name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	param.Name = name

	// Data type
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	param.TypeName = dt

	param.Loc.End = p.pos()
	return param, nil
}

// parseRoutineCharacteristic parses a routine characteristic.
func (p *Parser) parseRoutineCharacteristic() (*nodes.RoutineCharacteristic, bool, error) {
	start := p.pos()

	switch {
	case p.cur.Type == kwCOMMENT:
		p.advance()
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COMMENT", Value: val,
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwLANGUAGE:
		p.advance()
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, false, err
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "LANGUAGE", Value: name,
		}, true, nil

	case p.cur.Type == kwDETERMINISTIC:
		p.advance()
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DETERMINISTIC", Value: "YES",
		}, true, nil

	case p.cur.Type == kwNOT:
		// NOT DETERMINISTIC
		next := p.peekNext()
		if next.Type == kwDETERMINISTIC {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DETERMINISTIC", Value: "NO",
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwSQL:
		p.advance()
		if p.cur.Type == kwSECURITY {
			p.advance()
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return nil, false, err
			}
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "SQL SECURITY", Value: name,
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwCONTAINS:
		p.advance()
		p.match(kwSQL)
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "CONTAINS SQL",
		}, true, nil

	case p.cur.Type == kwNO:
		if next := p.peekNext(); next.Type == kwSQL {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "NO SQL",
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwREADS:
		p.advance()
		p.match(kwSQL)
		if p.cur.Type == kwDATA {
			p.advance()
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "READS SQL DATA",
		}, true, nil

	case p.cur.Type == kwMODIFIES:
		p.advance()
		p.match(kwSQL)
		if p.cur.Type == kwDATA {
			p.advance()
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "MODIFIES SQL DATA",
		}, true, nil
	}

	return nil, false, nil
}
