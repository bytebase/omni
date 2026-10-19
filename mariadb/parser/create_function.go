package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseCreateLoadableFunction parses a CREATE [AGGREGATE] FUNCTION for loadable UDFs.
// The caller has already consumed CREATE [OR REPLACE] (and AGGREGATE if present).
// The current token is FUNCTION.
//
// Ref: https://mariadb.com/kb/en/create-function-udf/
//
//	CREATE [OR REPLACE] [AGGREGATE] FUNCTION [IF NOT EXISTS] function_name
//	    RETURNS {STRING|INTEGER|REAL|DECIMAL}
//	    SONAME shared_library_name
func (p *Parser) parseCreateLoadableFunction(start int, orReplace, isAggregate bool) (*nodes.CreateFunctionStmt, error) {
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc:         nodes.Loc{Start: start},
		OrReplace:   orReplace,
		IsAggregate: isAggregate,
	}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

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
	retStart := p.pos()
	retType, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}
	stmt.Returns = &nodes.DataType{
		Loc:  nodes.Loc{Start: retStart, End: p.prev.End},
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateFunctionStmt parses a CREATE FUNCTION or CREATE PROCEDURE statement.
// The caller has already consumed CREATE [OR REPLACE] [DEFINER = user] and set
// isProcedure.
//
// Ref: https://mariadb.com/kb/en/create-procedure/
// Ref: https://mariadb.com/kb/en/create-function/
//
//	CREATE [OR REPLACE] [DEFINER = user] PROCEDURE [IF NOT EXISTS] sp_name
//	    ([proc_parameter[,...]]) [characteristic ...] routine_body
//	CREATE [OR REPLACE] [DEFINER = user] FUNCTION [IF NOT EXISTS] sp_name
//	    ([func_parameter[,...]]) RETURNS type_name [characteristic ...] routine_body
func (p *Parser) parseCreateFunctionStmt(start int, orReplace, isProcedure bool) (*nodes.CreateFunctionStmt, error) {
	p.advance() // consume FUNCTION or PROCEDURE

	// Completion: after CREATE FUNCTION / CREATE PROCEDURE, identifier context.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new name.
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CreateFunctionStmt{
		Loc:         nodes.Loc{Start: start},
		OrReplace:   orReplace,
		IsProcedure: isProcedure,
	}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

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
	bodyStart := p.pos()
	body, err := p.parseCompoundStmtOrStmt()
	if err != nil {
		return nil, err
	}
	bodyEnd := p.pos()
	stmt.Body = body
	stmt.BodyText = p.inputText(bodyStart, bodyEnd)

	stmt.Loc.End = p.prev.End
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

	param.Loc.End = p.prev.End
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
				Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COMMENT", Value: val,
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
			Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "LANGUAGE", Value: name,
		}, true, nil

	case p.cur.Type == kwDETERMINISTIC:
		p.advance()
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DETERMINISTIC", Value: "YES",
		}, true, nil

	case p.cur.Type == kwNOT:
		// NOT DETERMINISTIC
		next := p.peekNext()
		if next.Type == kwDETERMINISTIC {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DETERMINISTIC", Value: "NO",
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
				Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "SQL SECURITY", Value: name,
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwCONTAINS:
		p.advance()
		p.match(kwSQL)
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA ACCESS", Value: "CONTAINS SQL",
		}, true, nil

	case p.cur.Type == kwNO:
		if next := p.peekNext(); next.Type == kwSQL {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA ACCESS", Value: "NO SQL",
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
			Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA ACCESS", Value: "READS SQL DATA",
		}, true, nil

	case p.cur.Type == kwMODIFIES:
		p.advance()
		p.match(kwSQL)
		if p.cur.Type == kwDATA {
			p.advance()
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA ACCESS", Value: "MODIFIES SQL DATA",
		}, true, nil
	}

	return nil, false, nil
}
