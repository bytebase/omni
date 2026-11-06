package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
	"strings"
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
func (p *Parser) parseCreateLoadableFunction(start int, isAggregate bool) (*nodes.CreateFunctionStmt, error) {
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
// The caller has already consumed CREATE and set isProcedure.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-procedure.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-function.html
//
//	CREATE [DEFINER = user] PROCEDURE sp_name ([proc_parameter[,...]])
//	    [characteristic ...] routine_body
//	CREATE [DEFINER = user] FUNCTION sp_name ([func_parameter[,...]])
//	    RETURNS type_name [characteristic ...] routine_body
func (p *Parser) parseCreateFunctionStmt(start int, isProcedure bool) (*nodes.CreateFunctionStmt, error) {
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
	ref, err := p.parseNamedRoutineRef()
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

		dt, err := p.parseRoutineDataType()
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
	outerRoutine := p.routine
	p.routine = routineScope{noResultSet: !isProcedure, function: !isProcedure}
	defer func() { p.routine = outerRoutine }()
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
	dt, err := p.parseRoutineDataType()
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
		// COMMENT takes a quoted string. Reporting "no characteristic" here, with
		// COMMENT already consumed, let the body or the end of the statement follow.
		return nil, false, p.syntaxErrorAtCur()

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

// parseRoutineDataType reads the type of a routine's parameter, return value or local
// variable. A collation there needs a character set beside it: the action answers 1235,
// ER_NOT_SUPPORTED_YET, for `a CHAR(5) COLLATE ascii_bin` and runs it with CHARSET ascii.
//
// Ref: sp_param_type_with_opt_collate and type_with_opt_collate in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseRoutineDataType() (*nodes.DataType, error) {
	start := p.pos()
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt.Collate != "" && dt.Charset == "" {
		return nil, &ParseError{Message: "COLLATE with no CHARACTER SET in SP parameters, RETURNS, DECLARE", Position: start}
	}
	// A column's limits wait for the table to be prepared, and a routine's types have none
	// to wait for: the action checks them as they are read.
	hasScale := strings.Contains(p.inputText(start, p.prev.End), ",")
	if err := refuseRoutineTypeLimits(dt, hasScale, start); err != nil {
		return nil, err
	}
	return dt, nil
}

// refuseRoutineTypeLimits holds a routine's type to what a column would be held to later: 1439
// for a display width or length, 1074 for a CHAR or BINARY past 255, 1426 for a precision, 1425
// for a scale and 1427 for a scale past its precision.
//
// Ref: Create_field::init in MySQL 8.0 sql/create_field.cc.
func refuseRoutineTypeLimits(dt *nodes.DataType, hasScale bool, pos int) error {
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: pos}
	}
	switch dt.Name {
	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT":
		if dt.Length > 255 {
			return refuse("Display width out of range")
		}
	case "BIT":
		if dt.Length > 64 {
			return refuse("Display width out of range")
		}
	case "CHAR", "BINARY":
		if dt.Length > 255 {
			return refuse("Column length too big")
		}
	case "VARCHAR", "VARBINARY", "TEXT", "BLOB":
		if dt.Length > 4294967295 {
			return refuse("Display width out of range")
		}
	case "DECIMAL", "NUMERIC":
		switch {
		case dt.Length > 65:
			return refuse("Too-big precision")
		case dt.Scale > 30:
			return refuse("Too big scale")
		case hasScale && dt.Scale > dt.Length:
			return refuse("For float(M,D), double(M,D) or decimal(M,D), M must be >= D")
		}
	case "FLOAT", "DOUBLE":
		switch {
		case !hasScale && dt.Length > 53:
			return refuse("Incorrect column specifier")
		case hasScale && dt.Length > 255:
			return refuse("Display width out of range")
		case hasScale && (dt.Scale > 30 || dt.Scale > dt.Length):
			return refuse("For float(M,D), double(M,D) or decimal(M,D), M must be >= D")
		}
	case "TIME", "DATETIME", "TIMESTAMP":
		if dt.Length > 6 {
			return refuse("Too-big precision")
		}
	}
	return nil
}
