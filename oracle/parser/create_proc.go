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
func (p *Parser) parseCreateProcedureStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) *nodes.CreateProcedureStmt {
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		OrReplace:      orReplace,
		IfNotExists:    ifNotExists,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Loc:            nodes.Loc{Start: start},
	}

	// Procedure name
	stmt.Name = p.parseObjectName()

	// Optional SHARING = { METADATA | NONE }
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
		stmt.Parameters = p.parseParameterList()
	}

	// IS | AS
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}

	// PL/SQL block body (BEGIN ... END)
	stmt.Body = p.parsePLSQLBlock()

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseCreateFunctionStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) *nodes.CreateFunctionStmt {
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

	// Function name
	stmt.Name = p.parseObjectName()

	// Optional parameter list
	if p.cur.Type == '(' {
		stmt.Parameters = p.parseParameterList()
	}

	// RETURN type
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		stmt.ReturnType = p.parseTypeName()
	}

	// Optional SHARING = { METADATA | NONE }
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

	// Optional function properties (can appear in any order before IS/AS)
	p.parseFunctionProperties(stmt)

	// IS | AS
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}

	// PL/SQL block body (BEGIN ... END)
	stmt.Body = p.parsePLSQLBlock()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFunctionProperties parses optional DETERMINISTIC, PIPELINED, PARALLEL_ENABLE, RESULT_CACHE,
// AGGREGATE USING, SQL_MACRO, AUTHID, and other function property clauses.
func (p *Parser) parseFunctionProperties(stmt *nodes.CreateFunctionStmt) {
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
				p.advance() // consume USING
				p.parseObjectName() // consume implementation type
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
					p.parseObjectName()
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
			return
		}
	}
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
func (p *Parser) parseCreatePackageStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) *nodes.CreatePackageStmt {
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

	// Package name
	stmt.Name = p.parseObjectName()

	// Optional SHARING = { METADATA | NONE }
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
				// Each accessor: [ unit_kind ] [ schema. ] unit_name
				p.parseObjectName()
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

	// Package declarations/body - collect everything until END
	stmt.Body = p.parsePackageBody()

	// END [name] ;
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePackageBody parses the declarations inside a package specification or body.
// Stops when END is encountered.
func (p *Parser) parsePackageBody() *nodes.List {
	decls := &nodes.List{}

	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		// Skip standalone semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}

		// PROCEDURE declaration/definition in package
		if p.cur.Type == kwPROCEDURE {
			decl := p.parsePackageProcDecl()
			if decl != nil {
				decls.Items = append(decls.Items, decl)
			}
			continue
		}

		// FUNCTION declaration/definition in package
		if p.cur.Type == kwFUNCTION {
			decl := p.parsePackageFuncDecl()
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
		decl := p.parsePLSQLDeclaration()
		if decl == nil {
			// If we can't parse anything, skip a token to avoid infinite loop
			p.advance()
			continue
		}
		decls.Items = append(decls.Items, decl)
	}

	return decls
}

// parsePackageProcDecl parses a PROCEDURE declaration or definition inside a package.
//
//	PROCEDURE name [(params)] ;                    -- specification
//	PROCEDURE name [(params)] IS|AS body ;         -- body definition
func (p *Parser) parsePackageProcDecl() *nodes.CreateProcedureStmt {
	start := p.pos()
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional parameter list
	if p.cur.Type == '(' {
		stmt.Parameters = p.parseParameterList()
	}

	// Check for IS|AS (definition) or ; (declaration)
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		stmt.Body = p.parsePLSQLBlock()
	} else if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePackageFuncDecl parses a FUNCTION declaration or definition inside a package.
//
//	FUNCTION name [(params)] RETURN type ;                    -- specification
//	FUNCTION name [(params)] RETURN type IS|AS body ;         -- body definition
func (p *Parser) parsePackageFuncDecl() *nodes.CreateFunctionStmt {
	start := p.pos()
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional parameter list
	if p.cur.Type == '(' {
		stmt.Parameters = p.parseParameterList()
	}

	// RETURN type
	if p.cur.Type == kwRETURN {
		p.advance()
		stmt.ReturnType = p.parseTypeName()
	}

	// Optional function properties
	p.parseFunctionProperties(stmt)

	// Check for IS|AS (definition) or ; (declaration)
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		stmt.Body = p.parsePLSQLBlock()
	} else if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseParameterList parses a parenthesized parameter list: ( param1, param2, ... )
func (p *Parser) parseParameterList() *nodes.List {
	params := &nodes.List{}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		param := p.parseParameter()
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

	return params
}

// parseParameter parses a single parameter declaration.
//
//	name [IN | OUT | IN OUT] [NOCOPY] type [{:= | DEFAULT} expr]
func (p *Parser) parseParameter() *nodes.Parameter {
	start := p.pos()
	param := &nodes.Parameter{
		Loc: nodes.Loc{Start: start},
	}

	// Parameter name
	param.Name = p.parseIdentifier()
	if param.Name == "" {
		return nil
	}

	// Optional mode: IN, OUT, IN OUT
	mode := p.parseParameterMode()
	param.Mode = mode

	// Type name
	param.TypeName = p.parseTypeName()

	// Optional default value: := expr or DEFAULT expr
	if p.cur.Type == tokASSIGN {
		p.advance() // consume :=
		param.Default = p.parseExpr()
	} else if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		param.Default = p.parseExpr()
	}

	param.Loc.End = p.pos()
	return param
}

// parseParameterMode parses the optional IN/OUT/IN OUT/NOCOPY mode keywords.
func (p *Parser) parseParameterMode() string {
	if p.cur.Type == kwIN {
		next := p.peekNext()
		if next.Type == kwOUT {
			p.advance() // consume IN
			p.advance() // consume OUT
			// Optional NOCOPY after IN OUT
			if p.cur.Type == kwNOCOPY {
				p.advance()
				return "IN OUT NOCOPY"
			}
			return "IN OUT"
		}
		p.advance() // consume IN
		return "IN"
	}
	if p.cur.Type == kwOUT {
		p.advance() // consume OUT
		// Optional NOCOPY after OUT
		if p.cur.Type == kwNOCOPY {
			p.advance()
			return "OUT NOCOPY"
		}
		return "OUT"
	}
	return ""
}
