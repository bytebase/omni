package parser

import (
	"strings"

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
	stmt.Name, parseErr454 = p.parseReservedCheckedObjectName()
	if parseErr454 !=

		// Optional SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr454
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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

	if tn := p.tablePseudoType(stmt.Parameters, nil); tn != nil {
		return nil, p.syntaxErrorAtType(tn)
	}
	if err := p.checkCharsetSources(stmt.Parameters, nil); err != nil {
		return nil, err
	}
	if err := p.parseProcedureProperties(stmt, subprogramSchema); err != nil {
		return nil, err
	}

	if p.isIdentLikeStr("WRAPPED") {
		return p.parseWrappedProcedure(stmt)
	}

	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// PL/SQL block body (BEGIN ... END) or call spec
	var parseErr456 error
	stmt.Body, stmt.CallSpec, parseErr456 = p.parseSubprogramImplementation()
	if parseErr456 != nil {
		return nil, parseErr456
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func (p *Parser) parseWrappedProcedure(stmt *nodes.CreateProcedureStmt) (*nodes.CreateProcedureStmt, error) {
	wrappedTok := p.cur
	wrappedStart := p.cur.Loc
	stmt.Wrapped = true

	// The wrapped body runs to the next ';' or to the end of the parsed range;
	// text past p.lexer.end belongs to a later segment.
	wrappedEnd := p.lexer.end
	if idx := strings.IndexByte(p.source[wrappedTok.End:p.lexer.end], ';'); idx >= 0 {
		wrappedEnd = wrappedTok.End + idx
	}
	wrappedSourceEnd := trimRightSpace(p.source, wrappedEnd)
	if strings.TrimSpace(p.source[wrappedTok.End:wrappedSourceEnd]) == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if wrappedStart >= 0 && wrappedSourceEnd <= len(p.source) && wrappedStart < wrappedSourceEnd {
		stmt.WrappedSource = p.source[wrappedStart:wrappedSourceEnd]
	}
	stmt.Loc.End = wrappedSourceEnd
	p.prev = Token{Type: tokIDENT, Str: "WRAPPED", Loc: wrappedStart, End: wrappedSourceEnd}
	p.hasNext = false
	if wrappedEnd < p.lexer.end && p.source[wrappedEnd] == ';' {
		p.cur = Token{Type: ';', Str: ";", Loc: wrappedEnd, End: wrappedEnd + 1}
		p.lexer.pos = wrappedEnd + 1
	} else {
		p.cur = Token{Type: tokEOF, Loc: wrappedEnd, End: wrappedEnd}
		p.lexer.pos = wrappedEnd
	}
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
	stmt.Name, parseErr457 = p.parseReservedCheckedObjectName()
	if parseErr457 !=

		// Optional parameter list
		nil {
		return nil, parseErr457
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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

	if p.cur.Type != kwRETURN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr459 error
	stmt.ReturnType, parseErr459 = p.parsePLSQLDatatype(typeModsNone, charsetFlexible)
	if parseErr459 != nil {
		return nil, parseErr459
	}
	if stmt.ReturnType == nil || stmt.ReturnType.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
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
		p.parseFunctionProperties(stmt, subprogramSchema)
	if parseErr460 !=

		// IS | AS
		nil {
		return nil, parseErr460
	}
	if stmt.Implementation != nil {
		return p.finishTypeImplementedFunction(stmt)
	}

	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr461 error

	// PL/SQL block body (BEGIN ... END) or call spec
	stmt.Body, stmt.CallSpec, parseErr461 = p.parseSubprogramImplementation()
	if parseErr461 != nil {
		return nil, parseErr461
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFunctionProperties parses the clauses between a function's RETURN
// type and its IS | AS or ';', in any order, each at most once (PLS-00371).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-FUNCTION-statement.html
//
//	[ { invoker_rights_clause | accessible_by_clause | default_collation_clause
//	  | deterministic_clause | parallel_enable_clause | result_cache_clause
//	  | aggregate_clause | pipelined_clause | sql_macro_clause
//	  | streaming_clause }... ]
//	pipelined_clause    ::= PIPELINED [ { ROW | TABLE } POLYMORPHIC ]
//	result_cache_clause ::= RESULT_CACHE [ RELIES_ON ( [ data_source [, data_source ]... ] ) ]
//	sql_macro_clause    ::= SQL_MACRO [ ( [ TYPE => ] { SCALAR | TABLE } ) ]
//	aggregate_clause    ::= AGGREGATE USING [ schema. ] implementation_type
//
// A type can implement the function: AGGREGATE USING type, or USING type
// after PIPELINED. The clause list ends there and the function has no IS | AS
// body; Oracle 23ai rejects a following IS or clause (PLS-00103), USING
// without PIPELINED (PLS-00624), and AGGREGATE without USING after it. level
// limits AUTHID, ACCESSIBLE BY, and DEFAULT COLLATION as subprogramLevel says.
func (p *Parser) parseFunctionProperties(stmt *nodes.CreateFunctionStmt, level subprogramLevel) error {
	seen := clauseSeen{}
	// notPolymorphic is the first clause a polymorphic table function does
	// not take: Oracle 23ai rejects DETERMINISTIC, PARALLEL_ENABLE, ORDER
	// BY, CLUSTER BY, and AUTHID on one with PLS-00760, and RESULT_CACHE
	// with PLS-00999 (its TABLE parameter), before or after PIPELINED ...
	// POLYMORPHIC.
	//
	// notMacro is the first clause a SQL macro does not take: those, and
	// PIPELINED, whose collection result a macro cannot return (PLS-00778,
	// PLS-00630). A macro returns a character type: NUMBER, LONG, or a
	// %ROWTYPE result is PLS-00776, while CHAR, VARCHAR2, CLOB, and the
	// national types compile.
	var notPolymorphic, notMacro, macroTok *Token
	var partTok Token
	var streamToks []Token
	note := func() {
		tok := p.cur
		if notPolymorphic == nil {
			notPolymorphic = &tok
		}
		if notMacro == nil {
			notMacro = &tok
		}
	}
	checkPolymorphic := func(err error) error {
		switch {
		case err != nil:
			return err
		case stmt.Polymorphic != "" && notPolymorphic != nil:
			return p.syntaxErrorAtTok(*notPolymorphic)
		case stmt.SqlMacro && notMacro != nil:
			return p.syntaxErrorAtTok(*notMacro)
		case stmt.SqlMacro && !mayBeCharacterType(stmt.ReturnType):
			return p.syntaxErrorAtTok(*macroTok)
		}
		if stmt.Polymorphic == "" {
			if tn := p.tablePseudoType(stmt.Parameters, stmt.ReturnType); tn != nil {
				return p.syntaxErrorAtType(tn)
			}
		} else if err := p.checkPolymorphicSignature(stmt); err != nil {
			return err
		}
		if err := p.checkCharsetSources(stmt.Parameters, stmt.ReturnType); err != nil {
			return err
		}
		return p.checkParallelArguments(stmt, partTok, streamToks)
	}
	for {
		switch {
		case p.cur.Type == kwDETERMINISTIC:
			if err := p.firstClause(seen, "DETERMINISTIC"); err != nil {
				return err
			}
			note()
			stmt.Deterministic = true
			p.advance()
		case p.cur.Type == kwPIPELINED:
			if err := p.firstClause(seen, "PIPELINED"); err != nil {
				return err
			}
			if notMacro == nil {
				tok := p.cur
				notMacro = &tok
			}
			stmt.Pipelined = true
			p.advance()
			if (p.cur.Type == kwROW || p.cur.Type == kwTABLE) && p.isIdentLikeStrAt(p.peekNext(), "POLYMORPHIC") {
				// The row kind is required (PIPELINED POLYMORPHIC is
				// PLS-00103), and an object type method is not polymorphic
				// (PLS-00765); a nested function may be.
				if level == subprogramMethod {
					return p.syntaxErrorAtCur()
				}
				stmt.Polymorphic = "ROW"
				if p.cur.Type == kwTABLE {
					stmt.Polymorphic = "TABLE"
				}
				p.advance() // consume ROW or TABLE
				p.advance() // consume POLYMORPHIC
			}
		case p.cur.Type == kwPARALLEL_ENABLE:
			if err := p.firstClause(seen, "PARALLEL_ENABLE"); err != nil {
				return err
			}
			note()
			stmt.Parallel = true
			p.advance()
			if p.cur.Type == '(' {
				spec, argTok, err := p.parseParallelEnableSpec()
				if err != nil {
					return err
				}
				stmt.ParallelSpec = spec
				partTok = argTok
			}
		case p.cur.Type == kwORDER, p.cur.Type == kwCLUSTER:
			kind := "ORDER BY"
			if p.cur.Type == kwCLUSTER {
				kind = "CLUSTER BY"
			}
			if err := p.firstClause(seen, kind); err != nil {
				return err
			}
			note()
			clause, argTok, err := p.parseStreamingClause()
			if err != nil {
				return err
			}
			streamToks = append(streamToks, argTok)
			if stmt.Streaming == nil {
				stmt.Streaming = &nodes.List{}
			}
			stmt.Streaming.Items = append(stmt.Streaming.Items, clause)
		case p.cur.Type == kwRESULT_CACHE:
			if err := p.firstClause(seen, "RESULT_CACHE"); err != nil {
				return err
			}
			note()
			stmt.ResultCache = true
			p.advance()
			if p.isKeywordStr("RELIES_ON") {
				list, err := p.parseResultCacheReliesOn()
				if err != nil {
					return err
				}
				stmt.ReliesOn = list
			}
		case p.isKeywordStr("AGGREGATE"):
			aggTok := p.cur
			// PIPELINED and AGGREGATE exclude each other (PLS-00371), unless
			// the function is polymorphic: Oracle 23ai compiles PIPELINED
			// ROW POLYMORPHIC AGGREGATE USING type, and SQL_MACRO with it.
			if stmt.Pipelined && stmt.Polymorphic == "" {
				return p.syntaxErrorAtCur()
			}
			stmt.Aggregate = true
			p.advance() // consume AGGREGATE
			if p.cur.Type != kwUSING {
				return p.syntaxErrorAtCur()
			}
			if err := p.parseImplementationType(stmt); err != nil {
				return err
			}
			// An aggregate takes an argument: Oracle 23ai rejects one
			// without parameters (PLS-00652) and compiles one with two, or
			// with an OUT parameter, against its implementation type.
			if stmt.Parameters == nil || stmt.Parameters.Len() == 0 {
				return p.syntaxErrorAtTok(aggTok)
			}
			return checkPolymorphic(nil)
		case p.cur.Type == kwUSING:
			if !stmt.Pipelined {
				return p.syntaxErrorAtCur()
			}
			return checkPolymorphic(p.parseImplementationType(stmt))
		case p.isKeywordStr("SQL_MACRO"):
			if err := p.firstClause(seen, "SQL_MACRO"); err != nil {
				return err
			}
			// A type method is not a SQL macro (PLS-00781).
			if level == subprogramMethod {
				return p.syntaxErrorAtCur()
			}
			tok := p.cur
			macroTok = &tok
			stmt.SqlMacro = true
			p.advance() // consume SQL_MACRO
			if p.cur.Type == '(' {
				kind, err := p.parseSqlMacroType()
				if err != nil {
					return err
				}
				stmt.SqlMacroType = kind
			}
		case p.isKeywordStr("AUTHID") && level == subprogramSchema:
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			note()
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE") && level.takesAccessibleBy():
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation() && level == subprogramSchema:
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			collation, err := p.parseDefaultCollationClause()
			if err != nil {
				return err
			}
			stmt.DefaultCollation = collation
		default:
			return checkPolymorphic(nil)
		}
	}
}

// mayBeCharacterType reports whether tn may hold character data: it is not a
// predefined non-character type (the national character types count as
// character types) and not a %ROWTYPE record. A SQL macro returns such a
// type (PLS-00776), and item%CHARSET takes its character set from one
// (PLS-00550). A name the text does not settle is left to the engine.
func mayBeCharacterType(tn *nodes.TypeName) bool {
	if tn == nil {
		return true
	}
	if tn.IsPercRowtype {
		return false
	}
	if tn.IsPercType || tn.Names.Len() != 1 {
		return true
	}
	first, ok := tn.Names.Items[0].(*nodes.String)
	if !ok {
		return true
	}
	switch first.Str {
	case "NCHAR", "NVARCHAR2", "NCLOB":
		return true
	}
	return !plsqlNonCharacterTypes[first.Str]
}

// parseImplementationType parses USING [ schema. ] implementation_type
// [ @dblink ], the current token at USING. Every part needs its name: Oracle
// 23ai rejects USING impl. and USING impl@ (PLS-00103), and takes a dotted
// database link name (impl@lnk.dom).
func (p *Parser) parseImplementationType(stmt *nodes.CreateFunctionStmt) error {
	p.advance() // consume USING
	name, err := p.parseUnitName()
	if err != nil {
		return err
	}
	if p.cur.Type == '@' {
		p.advance()
		if !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		link, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		for p.cur.Type == '.' {
			p.advance()
			if !p.isIdentLike() {
				return p.syntaxErrorAtCur()
			}
			part, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			link += "." + part
		}
		name.DBLink = link
		name.Loc.End = p.prev.End
	}
	stmt.Implementation = name
	return nil
}

// finishTypeImplementedFunction ends a function a type implements, whose
// clause list ends at the type name: only its ';' follows.
func (p *Parser) finishTypeImplementedFunction(stmt *nodes.CreateFunctionStmt) (*nodes.CreateFunctionStmt, error) {
	if p.cur.Type != ';' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func (p *Parser) parseOptionalAuthID() (string, error) {
	if !p.isIdentLikeStr("AUTHID") {
		return "", nil
	}
	p.advance()
	if p.isIdentLikeStr("CURRENT_USER") {
		p.advance()
		return "CURRENT_USER", nil
	}
	if p.cur.Type == kwCURRENT {
		p.advance()
		if p.isIdentLikeStr("USER") {
			p.advance()
		}
		return "CURRENT_USER", nil
	}
	if p.isIdentLikeStr("DEFINER") {
		p.advance()
		return "DEFINER", nil
	}
	return "", p.syntaxErrorAtCur()
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
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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

	// ACCESSIBLE BY, AUTHID, and DEFAULT COLLATION head a specification,
	// in any order; a body takes none of them.
	if !stmt.IsBody {
		if err := p.parsePackageProperties(stmt); err != nil {
			return nil, err
		}
	}

	// IS | AS
	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr469 error

	// Package declarations/body - collect everything until END
	stmt.Body, parseErr469 = p.parsePackageBody()
	if parseErr469 !=

		// END [name] ;
		nil {
		return nil, parseErr469
	}

	if p.cur.Type != kwEND {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume END
	// Optional package name after END
	if p.atEndName() {
		p.advance() // consume name
	}
	if p.cur.Type != ';' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ;

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePackageBody parses the declarations inside a package specification or body.
// Stops when END is encountered.
func (p *Parser) parsePackageBody() (*nodes.List, error) {
	decls := &nodes.List{}

	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		// PROCEDURE declaration/definition in package
		if p.cur.Type == kwPROCEDURE {
			decl, parseErr470 := p.parsePackageProcDecl(subprogramPackaged)
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
			decl, parseErr471 := p.parsePackageFuncDecl(subprogramPackaged)
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

		// Variable/type/cursor declarations. Each ends with ';', and a
		// stray ';' between items is PLS-00103 on Oracle 23ai.
		decl, parseErr472 := p.parsePLSQLDeclaration()
		if parseErr472 != nil {
			return nil, parseErr472
		}
		if decl == nil || p.prev.Type != ';' {
			return nil, p.syntaxErrorAtCur()
		}
		decls.Items = append(decls.Items, decl)
	}

	return decls, nil
}

// parsePackageProcDecl parses a PROCEDURE declaration or definition inside a package.
//
//	PROCEDURE name [(params)] ;                    -- specification
//	PROCEDURE name [(params)] IS|AS body ;         -- body definition
func (p *Parser) parsePackageProcDecl(level subprogramLevel) (*nodes.CreateProcedureStmt, error) {
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
	if tn := p.tablePseudoType(stmt.Parameters, nil); tn != nil {
		return nil, p.syntaxErrorAtType(tn)
	}
	if err := p.checkCharsetSources(stmt.Parameters, nil); err != nil {
		return nil, err
	}
	if err := p.parseProcedureProperties(stmt, level); err != nil {
		return nil, err
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		var parseErr475 error
		stmt.Body, stmt.CallSpec, parseErr475 = p.parseSubprogramImplementation()
		if parseErr475 != nil {
			return nil, parseErr475
		}
	} else if p.cur.Type == ';' {
		p.advance()
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePackageFuncDecl parses a FUNCTION declaration or definition inside a package.
//
//	FUNCTION name [(params)] RETURN type ;                    -- specification
//	FUNCTION name [(params)] RETURN type IS|AS body ;         -- body definition
func (p *Parser) parsePackageFuncDecl(level subprogramLevel) (*nodes.CreateFunctionStmt, error) {
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
		stmt.ReturnType, parseErr478 = p.parsePLSQLDatatype(typeModsNone, charsetFlexible)
		if parseErr478 !=

			// Optional function properties
			nil {
			return nil, parseErr478
		}
		if stmt.ReturnType == nil || stmt.ReturnType.Names.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	} else {
		return nil, p.syntaxErrorAtCur()
	}
	parseErr479 := p.parseFunctionProperties(stmt, level)
	if parseErr479 !=

		// Check for IS|AS (definition) or ; (declaration)
		nil {
		return nil, parseErr479
	}
	if stmt.Implementation != nil {
		return p.finishTypeImplementedFunction(stmt)
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
		var parseErr480 error
		stmt.Body, stmt.CallSpec, parseErr480 = p.parseSubprogramImplementation()
		if parseErr480 != nil {
			return nil, parseErr480
		}
	} else if p.cur.Type == ';' {
		p.advance()
	} else {
		return nil, p.syntaxErrorAtCur()
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

	// A parameter type is unconstrained and may name ANY_CS.
	param.TypeName, parseErr484 = p.parsePLSQLDatatype(typeModsNone, charsetAnyCS)
	if parseErr484 !=

		// Optional default value: := expr or DEFAULT expr
		nil {
		return nil, parseErr484
	}
	if param.TypeName == nil || param.TypeName.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
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

// parseSubprogramImplementation parses what follows the IS|AS of a procedure
// or function: a call spec, or a PL/SQL block. Exactly one result is non-nil.
func (p *Parser) parseSubprogramImplementation() (nodes.StmtNode, *nodes.CallSpec, error) {
	if p.isCallSpecWord("LANGUAGE", "EXTERNAL", "MLE") {
		spec, err := p.parseCallSpec(true)
		if err != nil {
			return nil, nil, err
		}
		return nil, spec, nil
	}
	block, err := p.parsePLSQLBlock()
	if err != nil {
		return nil, nil, err
	}
	return block, nil, nil
}

// isCallSpecWord reports whether the current token, the first after IS|AS,
// is one of words unquoted; see isCallSpecWordToken.
func (p *Parser) isCallSpecWord(words ...string) bool {
	return isCallSpecWordToken(p.cur, words...)
}

// isCallSpecWordToken reports whether tok, the first token after IS|AS, is
// one of words unquoted. Oracle 23ai commits to a call spec on that word
// alone: a first declaration of a variable named LANGUAGE, EXTERNAL, or MLE
// is PLS-00103 in a standalone, nested, package, or type body subprogram
// alike, while a quoted "MLE" declares one. The parser and the splitter both
// decide by it, so that both end the unit at the same ';'.
func isCallSpecWordToken(tok Token, words ...string) bool {
	if tok.Type != tokIDENT {
		return false
	}
	for _, word := range words {
		if tok.Str == word {
			return true
		}
	}
	return false
}

// callSpecCClauseOf names the C call spec clause tok starts, or "" when it
// starts none.
func callSpecCClauseOf(tok Token) string {
	switch tok.Type {
	case kwNAME:
		return "NAME"
	case kwLIBRARY:
		return "LIBRARY"
	case kwWITH:
		return "WITH"
	case tokIDENT:
		switch tok.Str {
		case "LANGUAGE", "CALLING", "AGENT", "PARAMETERS":
			return tok.Str
		}
	}
	return ""
}

// parseCallSpec parses a call_spec, which publishes a Java method or a C
// function as the implementation of a procedure or function.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/lnpls/call-specification.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/adfns/external-procedures.html
//
//	call_spec ::= LANGUAGE { java_declaration | c_declaration }
//	java_declaration ::= JAVA NAME string
//	c_declaration ::= C { [ NAME name ] LIBRARY lib_name | LIBRARY lib_name [ NAME name ] }
//	    [ AGENT IN ( argument [, argument ]... ) ]
//	    [ WITH CONTEXT ]
//	    [ PARAMETERS ( external_parameter [, external_parameter ]... ) ]
//
// EXTERNAL, the superseded spelling of LANGUAGE C, is still accepted, with
// its own LANGUAGE C and CALLING STANDARD { C | PASCAL } clauses. Oracle 23ai
// takes the C clauses of either form in any order, rejects a repeated one
// (PLS-00139, -00140, -00142 through -00145, -00171), and requires LIBRARY
// (PLS-00247).
//
// terminated says whether a ';' must end the call spec; see finishCallSpec.
func (p *Parser) parseCallSpec(terminated bool) (*nodes.CallSpec, error) {
	spec := &nodes.CallSpec{Loc: nodes.Loc{Start: p.pos()}}
	if p.isKeywordStr("MLE") {
		return p.parseMLECallSpec(spec, terminated)
	}
	seen := make(map[string]bool)

	if p.isKeywordStr("EXTERNAL") {
		spec.Language = "C"
		spec.External = true
		p.advance() // consume EXTERNAL
	} else {
		p.advance() // consume LANGUAGE
		if p.cur.Type == kwJAVA {
			spec.Language = "JAVA"
			p.advance() // consume JAVA
			if p.cur.Type != kwNAME {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume NAME
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			spec.Name = p.cur.Str
			p.advance()
			return p.finishCallSpec(spec, terminated)
		}
		if p.cur.Type != tokIDENT || p.cur.Str != "C" {
			return nil, p.syntaxErrorAtCur()
		}
		spec.Language = "C"
		seen["LANGUAGE"] = true
		p.advance() // consume C
	}

	var err error
	for clause := callSpecCClauseOf(p.cur); clause != ""; clause = callSpecCClauseOf(p.cur) {
		if seen[clause] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[clause] = true
		p.advance() // consume the clause keyword
		switch clause {
		case "NAME":
			spec.Name, err = p.parseCallSpecCName()
		case "LIBRARY":
			spec.Library, err = p.parseObjectName()
			if err == nil && (spec.Library == nil || spec.Library.Name == "") {
				err = p.syntaxErrorAtCur()
			}
		case "LANGUAGE":
			// Only C follows LANGUAGE here (PLS-00103 for anything else).
			if p.cur.Type != tokIDENT || p.cur.Str != "C" {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume C
		case "CALLING":
			if !p.isKeywordStr("STANDARD") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume STANDARD
			if p.cur.Type != tokIDENT || (p.cur.Str != "C" && p.cur.Str != "PASCAL") {
				return nil, p.syntaxErrorAtCur()
			}
			spec.CallingStandard = p.cur.Str
			p.advance()
		case "AGENT":
			if p.cur.Type != kwIN {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume IN
			spec.AgentIn, err = p.parseCallSpecAgentIn()
		case "WITH":
			if p.cur.Type != kwCONTEXT {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume CONTEXT
			spec.WithContext = true
		case "PARAMETERS":
			spec.Parameters, err = p.parseCallSpecParameters()
		}
		if err != nil {
			return nil, err
		}
	}
	if !seen["LIBRARY"] {
		return nil, p.syntaxErrorAtCur()
	}
	return p.finishCallSpec(spec, terminated)
}

// parseMLECallSpec parses an MLE call spec, which publishes a JavaScript
// function as the implementation of a procedure or function. The current
// token is MLE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/mlejs/call-specifications-functions.html
//
//	MLE MODULE [ schema. ] module [ ENV [ schema. ] env ] SIGNATURE 'signature'
//	MLE LANGUAGE language_name [ PURE ] delimited_code
//
// Verified on Oracle 23ai: SIGNATURE is required and takes a string, PURE
// does not belong to the MODULE form, and both forms end with a required ';'
// (PLS-00103). The inline code is taken verbatim; see mleInlineCodeEnd.
func (p *Parser) parseMLECallSpec(spec *nodes.CallSpec, terminated bool) (*nodes.CallSpec, error) {
	spec.MLE = true
	p.advance() // consume MLE
	if p.isKeywordStr("MODULE") {
		p.advance() // consume MODULE
		var err error
		if spec.Module, err = p.parseObjectName(); err != nil {
			return nil, err
		}
		if spec.Module == nil || spec.Module.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		if p.isKeywordStr("ENV") {
			p.advance() // consume ENV
			if spec.Env, err = p.parseObjectName(); err != nil {
				return nil, err
			}
			if spec.Env == nil || spec.Env.Name == "" {
				return nil, p.syntaxErrorAtCur()
			}
		}
		if !p.isKeywordStr("SIGNATURE") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume SIGNATURE
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		spec.Name = p.cur.Str
		p.advance()
		return p.finishCallSpec(spec, terminated)
	}

	if !p.isKeywordStr("LANGUAGE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume LANGUAGE
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	langTok := p.cur
	spec.Language = langTok.Str
	codeStart, codeEnd, pure, ok := mleInlineCodeEnd(p.source, langTok.End, p.lexer.end)
	if !ok {
		return nil, &ParseError{
			Message:  "syntax error: missing closing delimiter for MLE language code",
			Position: codeStart,
		}
	}
	spec.Pure = pure
	spec.Code = p.source[codeStart:codeEnd]
	// Resume lexing after the closing delimiter.
	p.hasNext = false
	p.lexer.pos = codeEnd
	p.prev = Token{Type: tokIDENT, Loc: codeStart, End: codeEnd}
	p.cur = p.lexer.NextToken()
	return p.finishCallSpec(spec, terminated)
}

// parseCallSpecCName parses the C function name of a C call spec, an
// identifier whose case is kept when it is quoted.
func (p *Parser) parseCallSpecCName() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	return p.parseIdentifier()
}

// parseCallSpecAgentIn parses the parenthesized argument of AGENT IN: the
// name of the one parameter that carries the agent name. Oracle 23ai rejects a
// second argument with PLS-00103 although the reference shows a list.
func (p *Parser) parseCallSpecAgentIn() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume (
	name, err := p.parseCallSpecCName()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume )
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}, nil
}

// parseCallSpecParameters parses the parenthesized external parameters of a C
// call spec. Each entry is kept as its words joined by spaces.
func (p *Parser) parseCallSpecParameters() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume (
	list := &nodes.List{}
	for {
		param, err := p.parseExternalParameter()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, &nodes.String{Str: param})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume )
	return list, nil
}

// parseExternalParameter parses one external_parameter of a C call spec and
// returns its words joined by spaces.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/adfns/external-procedures.html
//
//	external_parameter ::= { CONTEXT | SELF | RETURN | parameter_name }
//	    [ NATIVE | property [ IN | OUT ] | INDICATOR STRUCT ]
//	    [ BY { REFERENCE | VALUE } ] [ external_datatype ]
//	property ::= INDICATOR | LENGTH | MAXLEN | DURATION | CHARSETID | CHARSETFORM | TDO
//
// The shape follows the token sets Oracle 23ai expects at each position
// (PLS-00103): one property at most, INDICATOR STRUCT ends the entry, IN and
// OUT only after a property, one BY clause, and one datatype from a fixed set.
// Whether a combination suits the parameter is a later check (PLS-00235,
// PLS-00250, PLS-00253), not a syntax rule.
func (p *Parser) parseExternalParameter() (string, error) {
	var words []string
	take := func() {
		words = append(words, p.cur.Str)
		p.advance()
	}

	// The head: CONTEXT, SELF, RETURN, or a parameter name.
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	if p.cur.Type == tokQIDENT {
		name, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		words = append(words, name)
	} else {
		take()
	}

	switch {
	case p.isKeywordStr("NATIVE"):
		take()
	case p.isKeywordStr("INDICATOR"):
		take()
		if p.isKeywordStr("STRUCT") {
			take()
			return strings.Join(words, " "), nil
		}
		if p.cur.Type == kwIN || p.cur.Type == kwOUT {
			take()
		}
	case p.isKeywordStr("LENGTH"), p.isKeywordStr("MAXLEN"), p.isKeywordStr("DURATION"),
		p.isKeywordStr("CHARSETID"), p.isKeywordStr("CHARSETFORM"), p.isKeywordStr("TDO"):
		take()
		if p.cur.Type == kwIN || p.cur.Type == kwOUT {
			take()
		}
	}

	if p.cur.Type == kwBY {
		take()
		if !p.isKeywordStr("REFERENCE") && !p.isKeywordStr("VALUE") {
			return "", p.syntaxErrorAtCur()
		}
		take()
	}

	switch {
	case p.isKeywordStr("UNSIGNED"):
		take()
		if !p.isKeywordStr("LONG") && !p.isKeywordStr("CHAR") && !p.isKeywordStr("SHORT") && !p.isKeywordStr("INT") {
			return "", p.syntaxErrorAtCur()
		}
		take()
	case p.cur.Type != tokQIDENT && p.isIdentLike() && externalDatatypes[p.cur.Str]:
		take()
	}
	return strings.Join(words, " "), nil
}

// externalDatatypes holds the external datatypes of a C call spec parameter,
// as Oracle 23ai lists them in its PLS-00103 message; UNSIGNED takes LONG,
// CHAR, SHORT, or INT and is handled apart.
var externalDatatypes = map[string]bool{
	"ARRAY": true, "CHAR": true, "DOUBLE": true, "FLOAT": true, "INT": true,
	"LONG": true, "OCICOLL": true, "OCIDATE": true, "OCIDATETIME": true,
	"OCIDURATION": true, "OCIINTERVAL": true, "OCILOBLOCATOR": true,
	"OCINUMBER": true, "OCIRAW": true, "OCIREF": true, "OCIREFCURSOR": true,
	"OCIROWID": true, "OCISTRING": true, "OCITYPE": true, "ORLANY": true,
	"ORLVARY": true, "RAW": true, "SB1": true, "SB2": true, "SB4": true,
	"SHORT": true, "SIZE_T": true, "STRING": true, "STRUCT": true, "UB1": true,
	"UB2": true, "UB4": true, "VALIST": true, "VOID": true,
}

// finishCallSpec closes the call spec's Loc at its last token. A call spec
// that implements a stored subprogram, a package member, or a type body method
// is PL/SQL text ended by a required ';' (Oracle compiles the unit with
// PLS-00103 without it), which is consumed. In an object type specification
// (terminated false) the ',' or ')' after it ends it instead, and a ';' there
// is a syntax error.
func (p *Parser) finishCallSpec(spec *nodes.CallSpec, terminated bool) (*nodes.CallSpec, error) {
	spec.Loc.End = p.prev.End
	if !terminated {
		return spec, nil
	}
	if p.cur.Type != ';' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ;
	return spec, nil
}
