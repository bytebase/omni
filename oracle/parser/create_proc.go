package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateProcedureStmt parses a CREATE [OR REPLACE] PROCEDURE statement.
//
// BNF: oracle/parser/bnf/CREATE-PROCEDURE.bnf
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ]
//	    [ EDITIONABLE | NONEDITIONABLE ]
//	    PROCEDURE [ schema. ] procedure_name
//	    [ SHARING = { METADATA | NONE } ]
//	    plsql_procedure_source ;
func (p *Parser) parseCreateProcedureStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) (*nodes.CreateProcedureStmt, error) {
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		OrReplace:      orReplace,
		IfNotExists:    ifNotExists,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Loc:            nodes.Loc{Start: start},
	}
	var parseErr454 error

	// Procedure name
	stmt.Name, parseErr454 = p.parseObjectName()
	if parseErr454 !=

		// Optional SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr454
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	// Optional parameter list
	if p.cur.Type == '(' {
		var parseErr455 error
		stmt.Parameters, parseErr455 = p.parseParameterList()
		if parseErr455 !=

			// IS | AS
			nil {
			return nil, parseErr455
		}
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	var parseErr456 error

	// PL/SQL block body (BEGIN ... END)
	stmt.Body, parseErr456 = p.parsePLSQLBlock()
	if parseErr456 != nil {
		return nil, parseErr456
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateFunctionStmt parses a CREATE [OR REPLACE] FUNCTION statement.
//
// BNF: oracle/parser/bnf/CREATE-FUNCTION.bnf
//
//	CREATE [ OR REPLACE ] [ EDITIONABLE | NONEDITIONABLE ]
//	    FUNCTION [ IF NOT EXISTS ] [ schema. ] function_name
//	    [ ( parameter_declaration [, parameter_declaration ]... ) ]
//	    RETURN datatype
//	    [ SHARING = { METADATA | NONE } ]
//	    [ { invoker_rights_clause
//	      | accessible_by_clause
//	      | default_collation_clause
//	      | deterministic_clause
//	      | parallel_enable_clause
//	      | result_cache_clause
//	      | aggregate_clause
//	      | pipelined_clause
//	      | sql_macro_clause }... ]
//	    { IS | AS }
//	    { plsql_function_source | call_spec } ;
func (p *Parser) parseCreateFunctionStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) (*nodes.CreateFunctionStmt, error) {
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		OrReplace:      orReplace,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Loc:            nodes.Loc{Start: start},
	}

	// IF NOT EXISTS (for FUNCTION, it comes after the FUNCTION keyword per BNF)
	if !ifNotExists && p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
				ifNotExists = true
			}
		}
	}
	stmt.IfNotExists = ifNotExists
	var parseErr457 error

	// Function name
	stmt.Name, parseErr457 = p.parseObjectName()
	if parseErr457 !=

		// Optional parameter list
		nil {
		return nil, parseErr457
	}

	if p.cur.Type == '(' {
		var parseErr458 error
		stmt.Parameters, parseErr458 = p.parseParameterList()
		if parseErr458 !=

			// RETURN type
			nil {
			return nil, parseErr458
		}
	}

	if p.cur.Type == kwRETURN {
		p.advance()
		var // consume RETURN
		parseErr459 error
		stmt.ReturnType, parseErr459 = p.parseTypeName()
		if parseErr459 !=

			// Optional SHARING = { METADATA | NONE }
			nil {
			return nil, parseErr459
		}
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}
	parseErr460 :=

		// Optional function properties (can appear in any order before IS/AS)
		p.parseFunctionProperties(stmt)
	if parseErr460 !=

		// IS | AS
		nil {
		return nil, parseErr460
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	var parseErr461 error

	// PL/SQL block body (BEGIN ... END)
	stmt.Body, parseErr461 = p.parsePLSQLBlock()
	if parseErr461 != nil {
		return nil, parseErr461
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFunctionProperties parses optional DETERMINISTIC, PIPELINED, PARALLEL_ENABLE, RESULT_CACHE,
// AGGREGATE USING, SQL_MACRO, AUTHID, and other function property clauses.
func (p *Parser) parseFunctionProperties(stmt *nodes.CreateFunctionStmt) error {
	for {
		switch {
		case p.cur.Type == kwDETERMINISTIC:
			stmt.Deterministic = true
			p.advance()
		case p.cur.Type == kwPIPELINED:
			stmt.Pipelined = true
			p.advance()
		case p.cur.Type == kwPARALLEL_ENABLE:
			stmt.Parallel = true
			p.advance()
		case p.cur.Type == kwRESULT_CACHE:
			stmt.ResultCache = true
			p.advance()
		case p.isIdentLikeStr("AGGREGATE"):
			stmt.Aggregate = true
			p.advance() // consume AGGREGATE
			if p.cur.Type == kwUSING {
				p.advance()
				parseDiscard463, // consume USING
					parseErr462 := p.parseObjectName()
				_ = // consume implementation type
					parseDiscard463
				if parseErr462 != nil {
					return parseErr462
				}
			}
		case p.isIdentLikeStr("SQL_MACRO"):
			stmt.SqlMacro = true
			p.advance() // consume SQL_MACRO
			// Optional ( SCALAR | TABLE )
			if p.cur.Type == '(' {
				p.advance()
				if p.isIdentLike() {
					p.advance() // consume SCALAR or TABLE
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		case p.isIdentLikeStr("AUTHID"):
			// invoker_rights_clause: AUTHID { CURRENT_USER | DEFINER }
			p.advance() // consume AUTHID
			if p.isIdentLike() {
				p.advance() // consume CURRENT_USER or DEFINER
			}
		case p.isIdentLikeStr("ACCESSIBLE"):
			// accessible_by_clause: ACCESSIBLE BY ( accessor [, ...] )
			p.advance() // consume ACCESSIBLE
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					parseDiscard465, parseErr464 := p.parseObjectName()
					_ = parseDiscard465
					if parseErr464 != nil {
						return parseErr464
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		case p.cur.Type == kwDEFAULT && p.isIdentLikeStrAt(p.peekNext(), "COLLATION"):
			// default_collation_clause: DEFAULT COLLATION collation_name
			p.advance() // consume DEFAULT
			p.advance() // consume COLLATION
			if p.isIdentLike() {
				p.advance() // consume collation name
			}
		default:
			return nil
		}
	}
	return nil
}

// parseCreatePackageStmt parses a CREATE [OR REPLACE] PACKAGE [BODY] statement.
//
// BNF: oracle/parser/bnf/CREATE-PACKAGE.bnf
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ]
//	    [ EDITIONABLE | NONEDITIONABLE ]
//	    PACKAGE [ schema. ] package_name
//	    [ SHARING = { METADATA | NONE } ]
//	    [ ACCESSIBLE BY ( accessor [, accessor ]... ) ]
//	    [ invoker_rights_clause ]
//	    { IS | AS }
//	    plsql_package_source ;
//
// BNF: oracle/parser/bnf/CREATE-PACKAGE-BODY.bnf
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ]
//	    [ EDITIONABLE | NONEDITIONABLE ]
//	    PACKAGE BODY [ schema. ] package_name
//	    { IS | AS }
//	    plsql_package_body_source ;
func (p *Parser) parseCreatePackageStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) (*nodes.CreatePackageStmt, error) {
	p.advance() // consume PACKAGE

	stmt := &nodes.CreatePackageStmt{
		OrReplace:      orReplace,
		IfNotExists:    ifNotExists,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Loc:            nodes.Loc{Start: start},
	}

	// Check for BODY keyword
	if p.cur.Type == kwBODY {
		stmt.IsBody = true
		p.advance() // consume BODY
	}
	var parseErr466 error

	// Package name
	stmt.Name, parseErr466 = p.parseObjectName()
	if parseErr466 !=

		// Optional SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr466
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	// Optional ACCESSIBLE BY ( ... ) — skip for package body
	if !stmt.IsBody && p.isIdentLikeStr("ACCESSIBLE") {
		p.advance() // consume ACCESSIBLE
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseDiscard468,
					// Each accessor: [ unit_kind ] [ schema. ] unit_name
					parseErr467 := p.parseObjectName()
				_ = parseDiscard468
				if parseErr467 != nil {
					return nil, parseErr467
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// Optional invoker_rights_clause: AUTHID { CURRENT_USER | DEFINER }
	if !stmt.IsBody && p.isIdentLikeStr("AUTHID") {
		p.advance() // consume AUTHID
		if p.isIdentLike() {
			p.advance() // consume CURRENT_USER or DEFINER
		}
	}

	// IS | AS
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	var parseErr469 error

	// Package declarations/body - collect everything until END
	stmt.Body, parseErr469 = p.parsePackageBody()
	if parseErr469 !=

		// END [name] ;
		nil {
		return nil, parseErr469
	}

	if p.cur.Type == kwEND {
		p.advance() // consume END
	}
	// Optional package name after END
	if p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance() // consume name
	}
	if p.cur.Type == ';' {
		p.advance() // consume ;
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePackageBody parses the declarations inside a package specification or body.
// Stops when END is encountered.
func (p *Parser) parsePackageBody() (*nodes.List, error) {
	decls := &nodes.List{}

	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		// Skip standalone semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}

		// PROCEDURE declaration/definition in package
		if p.cur.Type == kwPROCEDURE {
			decl, parseErr470 := p.parsePackageProcDecl()
			if parseErr470 != nil {
				return nil, parseErr470
			}
			if decl != nil {
				decls.Items = append(decls.Items, decl)
			}
			continue
		}

		// FUNCTION declaration/definition in package
		if p.cur.Type == kwFUNCTION {
			decl, parseErr471 := p.parsePackageFuncDecl()
			if parseErr471 != nil {
				return nil, parseErr471
			}
			if decl != nil {
				decls.Items = append(decls.Items, decl)
			}
			continue
		}

		// BEGIN section (package body initialization)
		if p.cur.Type == kwBEGIN {
			break
		}

		// Variable/type/cursor declarations
		decl, parseErr472 := p.parsePLSQLDeclaration()
		if parseErr472 != nil {
			return nil,

				// If we can't parse anything, skip a token to avoid infinite loop
				parseErr472
		}
		if decl == nil {

			p.advance()
			continue
		}
		decls.Items = append(decls.Items, decl)
	}

	return decls, nil
}

// parsePackageProcDecl parses a PROCEDURE declaration or definition inside a package.
//
//	PROCEDURE name [(params)] ;                    -- specification
//	PROCEDURE name [(params)] IS|AS body ;         -- body definition
func (p *Parser) parsePackageProcDecl() (*nodes.CreateProcedureStmt, error) {
	start := p.pos()
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr473 error

	stmt.Name, parseErr473 = p.parseObjectName()
	if parseErr473 !=

		// Optional parameter list
		nil {
		return nil, parseErr473
	}

	if p.cur.Type == '(' {
		var parseErr474 error
		stmt.Parameters, parseErr474 = p.parseParameterList()
		if parseErr474 !=

			// Check for IS|AS (definition) or ; (declaration)
			nil {
			return nil, parseErr474
		}
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		var parseErr475 error
		stmt.Body, parseErr475 = p.parsePLSQLBlock()
		if parseErr475 != nil {
			return nil, parseErr475
		}
	} else if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePackageFuncDecl parses a FUNCTION declaration or definition inside a package.
//
//	FUNCTION name [(params)] RETURN type ;                    -- specification
//	FUNCTION name [(params)] RETURN type IS|AS body ;         -- body definition
func (p *Parser) parsePackageFuncDecl() (*nodes.CreateFunctionStmt, error) {
	start := p.pos()
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr476 error

	stmt.Name, parseErr476 = p.parseObjectName()
	if parseErr476 !=

		// Optional parameter list
		nil {
		return nil, parseErr476
	}

	if p.cur.Type == '(' {
		var parseErr477 error
		stmt.Parameters, parseErr477 = p.parseParameterList()
		if parseErr477 !=

			// RETURN type
			nil {
			return nil, parseErr477
		}
	}

	if p.cur.Type == kwRETURN {
		p.advance()
		var parseErr478 error
		stmt.ReturnType, parseErr478 = p.parseTypeName()
		if parseErr478 !=

			// Optional function properties
			nil {
			return nil, parseErr478
		}
	}
	parseErr479 := p.parseFunctionProperties(stmt)
	if parseErr479 !=

		// Check for IS|AS (definition) or ; (declaration)
		nil {
		return nil, parseErr479
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		var parseErr480 error
		stmt.Body, parseErr480 = p.parsePLSQLBlock()
		if parseErr480 != nil {
			return nil, parseErr480
		}
	} else if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseParameterList parses a parenthesized parameter list: ( param1, param2, ... )
func (p *Parser) parseParameterList() (*nodes.List, error) {
	params := &nodes.List{}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		param, parseErr481 := p.parseParameter()
		if parseErr481 != nil {
			return nil, parseErr481
		}
		if param != nil {
			params.Items = append(params.Items, param)
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	if p.cur.Type == ')' {
		p.advance() // consume ')'
	}

	return params, nil
}

// parseParameter parses a single parameter declaration.
//
//	name [IN | OUT | IN OUT] [NOCOPY] type [{:= | DEFAULT} expr]
func (p *Parser) parseParameter() (*nodes.Parameter, error) {
	start := p.pos()
	param := &nodes.Parameter{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr482 error

	// Parameter name
	param.Name, parseErr482 = p.parseIdentifier()
	if parseErr482 != nil {
		return nil, parseErr482
	}
	if param.Name == "" {
		return nil, nil
	}

	// Optional mode: IN, OUT, IN OUT
	mode, parseErr483 := p.parseParameterMode()
	if parseErr483 != nil {
		return nil,

			// Type name
			parseErr483
	}
	param.Mode = mode
	var parseErr484 error

	param.TypeName, parseErr484 = p.parseTypeName()
	if parseErr484 !=

		// Optional default value: := expr or DEFAULT expr
		nil {
		return nil, parseErr484
	}

	if p.cur.Type == tokASSIGN {
		p.advance()
		var // consume :=
		parseErr485 error
		param.Default, parseErr485 = p.parseExpr()
		if parseErr485 != nil {
			return nil, parseErr485
		}
	} else if p.cur.Type == kwDEFAULT {
		p.advance()
		var // consume DEFAULT
		parseErr486 error
		param.Default, parseErr486 = p.parseExpr()
		if parseErr486 != nil {
			return nil, parseErr486
		}
	}

	param.Loc.End = p.prev.End
	return param, nil
}

// parseParameterMode parses the optional IN/OUT/IN OUT/NOCOPY mode keywords.
func (p *Parser) parseParameterMode() (string, error) {
	if p.cur.Type == kwIN {
		next := p.peekNext()
		if next.Type == kwOUT {
			p.advance() // consume IN
			p.advance() // consume OUT
			// Optional NOCOPY after IN OUT
			if p.cur.Type == kwNOCOPY {
				p.advance()
				return "IN OUT NOCOPY", nil
			}
			return "IN OUT", nil
		}
		p.advance() // consume IN
		return "IN", nil
	}
	if p.cur.Type == kwOUT {
		p.advance() // consume OUT
		// Optional NOCOPY after OUT
		if p.cur.Type == kwNOCOPY {
			p.advance()
			return "OUT NOCOPY", nil
		}
		return "OUT", nil
	}
	return "", nil
}
