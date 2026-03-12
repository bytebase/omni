package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

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

	stmt := &nodes.CreateFunctionStmt{
		Loc:         nodes.Loc{Start: start},
		IsProcedure: isProcedure,
	}

	// Function/procedure name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// Parameter list
	if _, err := p.expect('('); err != nil {
		return nil, err
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
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.Returns = dt
	}

	// Characteristics
	for {
		ch, ok := p.parseRoutineCharacteristic()
		if !ok {
			break
		}
		stmt.Characteristics = append(stmt.Characteristics, ch)
	}

	// Routine body — consume everything until we hit a semicolon or EOF
	bodyStart := p.pos()
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' && depth == 0 {
			break
		}
		if p.cur.Type == kwBEGIN {
			depth++
		}
		if p.cur.Type == kwEND {
			if depth > 0 {
				depth--
				if depth == 0 {
					p.advance()
					break
				}
			} else {
				break
			}
		}
		p.advance()
	}
	bodyEnd := p.pos()
	if bodyEnd > bodyStart {
		stmt.Body = p.lexer.input[bodyStart:bodyEnd]
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
func (p *Parser) parseRoutineCharacteristic() (*nodes.RoutineCharacteristic, bool) {
	start := p.pos()

	switch {
	case p.cur.Type == kwCOMMENT:
		p.advance()
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COMMENT", Value: val,
			}, true
		}
		return nil, false

	case p.cur.Type == kwLANGUAGE:
		p.advance()
		name, _, _ := p.parseIdentifier()
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "LANGUAGE", Value: name,
		}, true

	case p.cur.Type == tokIDENT && eqFold(p.cur.Str, "deterministic"):
		p.advance()
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DETERMINISTIC", Value: "YES",
		}, true

	case p.cur.Type == kwNOT:
		// NOT DETERMINISTIC
		if next := p.peekNext(); next.Type == tokIDENT && eqFold(next.Str, "deterministic") {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DETERMINISTIC", Value: "NO",
			}, true
		}
		return nil, false

	case p.cur.Type == kwSQL:
		p.advance()
		if p.cur.Type == kwSECURITY || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "security")) {
			p.advance()
			name, _, _ := p.parseIdentifier()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "SQL SECURITY", Value: name,
			}, true
		}
		return nil, false

	case p.cur.Type == tokIDENT && eqFold(p.cur.Str, "contains"):
		p.advance()
		p.match(kwSQL)
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "CONTAINS SQL",
		}, true

	case p.cur.Type == kwNO:
		if next := p.peekNext(); next.Type == kwSQL {
			p.advance()
			p.advance()
			return &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "NO SQL",
			}, true
		}
		return nil, false

	case p.cur.Type == tokIDENT && eqFold(p.cur.Str, "reads"):
		p.advance()
		p.match(kwSQL)
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "data") {
			p.advance()
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "READS SQL DATA",
		}, true

	case p.cur.Type == tokIDENT && eqFold(p.cur.Str, "modifies"):
		p.advance()
		p.match(kwSQL)
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "data") {
			p.advance()
		}
		return &nodes.RoutineCharacteristic{
			Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA ACCESS", Value: "MODIFIES SQL DATA",
		}, true
	}

	return nil, false
}
