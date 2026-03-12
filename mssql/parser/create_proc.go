// Package parser - create_proc.go implements T-SQL CREATE PROCEDURE/FUNCTION statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateProcedureStmt parses a CREATE [OR ALTER] PROCEDURE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-procedure-transact-sql
//
//	CREATE [ OR ALTER ] { PROC | PROCEDURE } [ schema_name. ] procedure_name
//	    [ { @parameter [ type_schema_name. ] data_type }
//	      [ VARYING ] [ = default ] [ OUT | OUTPUT | READONLY ]
//	    ] [ ,...n ]
//	[ WITH <procedure_option> [ ,...n ] ]
//	[ FOR REPLICATION ]
//	AS { [ BEGIN ] sql_statement [;] [ ...n ] [ END ] }
//
//	<procedure_option> ::=
//	    [ ENCRYPTION ]
//	    [ RECOMPILE ]
//	    [ EXECUTE AS { CALLER | SELF | OWNER | 'user_name' } ]
//	    [ NATIVE_COMPILATION ]
//	    [ SCHEMABINDING ]
func (p *Parser) parseCreateProcedureStmt(orAlter bool) *nodes.CreateProcedureStmt {
	loc := p.pos()

	stmt := &nodes.CreateProcedureStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc},
	}

	// Procedure name
	stmt.Name = p.parseTableRef()

	// Parameters
	if p.cur.Type == tokVARIABLE {
		stmt.Params = p.parseParamDefList()
	}

	// WITH <procedure_option> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isRoutineOption(next) {
			p.advance() // consume WITH
			stmt.Options = p.parseRoutineOptionList()
		}
	}

	// AS
	p.match(kwAS)

	// Body (BEGIN...END block or single statement)
	stmt.Body = p.parseStmt()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateFunctionStmt parses a CREATE [OR ALTER] FUNCTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-function-transact-sql
//
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	    ( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type
//	        [ = default_value ] [ READONLY ] } [ ,...n ] ] )
//	RETURNS return_data_type
//	    [ WITH <function_option> [ ,...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN scalar_expression
//	    END
//
//	RETURNS TABLE
//	    [ WITH <function_option> [ ,...n ] ]
//	    [ AS ]
//	    RETURN ( select_stmt )
//
//	RETURNS @return_variable TABLE <table_type_definition>
//	    [ WITH <function_option> [ ,...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN
//	    END
//
//	<function_option> ::=
//	    { ENCRYPTION | SCHEMABINDING | RETURNS NULL ON NULL INPUT
//	      | CALLED ON NULL INPUT | EXECUTE AS Clause
//	      | INLINE = { ON | OFF } | NATIVE_COMPILATION }
func (p *Parser) parseCreateFunctionStmt(orAlter bool) *nodes.CreateFunctionStmt {
	loc := p.pos()

	stmt := &nodes.CreateFunctionStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc},
	}

	// Function name
	stmt.Name = p.parseTableRef()

	// Parameters in parentheses
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ')' {
			stmt.Params = p.parseParamDefList()
		}
		_, _ = p.expect(')')
	}

	// RETURNS type or RETURNS TABLE
	if _, ok := p.match(kwRETURNS); ok {
		if p.cur.Type == kwTABLE {
			p.advance()
			stmt.ReturnsTable = &nodes.ReturnsTableDef{
				Loc: nodes.Loc{Start: p.pos()},
			}
			// Check for table definition (multi-statement TVF)
			if p.cur.Type == '(' {
				p.advance()
				var cols []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					col := p.parseColumnDef()
					if col != nil {
						cols = append(cols, col)
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				_, _ = p.expect(')')
				stmt.ReturnsTable.Columns = &nodes.List{Items: cols}
			}
			stmt.ReturnsTable.Loc.End = p.pos()
		} else if p.cur.Type == tokVARIABLE {
			// RETURNS @var TABLE (...)
			varName := p.cur.Str
			p.advance() // consume @var
			if p.cur.Type == kwTABLE {
				p.advance() // consume TABLE
				stmt.ReturnsTable = &nodes.ReturnsTableDef{
					Variable: varName,
					Loc:      nodes.Loc{Start: p.pos()},
				}
				if p.cur.Type == '(' {
					p.advance()
					var cols []nodes.Node
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						col := p.parseColumnDef()
						if col != nil {
							cols = append(cols, col)
						}
						if _, ok := p.match(','); !ok {
							break
						}
					}
					_, _ = p.expect(')')
					stmt.ReturnsTable.Columns = &nodes.List{Items: cols}
				}
				stmt.ReturnsTable.Loc.End = p.pos()
			}
		} else {
			stmt.ReturnType = p.parseDataType()
		}
	}

	// WITH <function_option> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isRoutineOption(next) {
			p.advance() // consume WITH
			stmt.Options = p.parseRoutineOptionList()
		}
	}

	// AS
	p.match(kwAS)

	// Body: BEGIN...END or RETURN SELECT...
	if p.cur.Type == kwRETURN {
		// Inline table-valued function: RETURN SELECT ...
		retLoc := p.pos()
		p.advance() // consume RETURN
		selectStmt := p.parseSelectStmt()
		stmt.Body = &nodes.ReturnStmt{
			Value: &nodes.SubqueryExpr{Query: selectStmt, Loc: nodes.Loc{Start: retLoc}},
			Loc:   nodes.Loc{Start: retLoc},
		}
	} else {
		stmt.Body = p.parseStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseParamDefList parses a comma-separated list of parameter definitions.
//
//	param_def_list = param_def { ',' param_def }
//	param_def = @name type [= default] [OUTPUT|READONLY]
func (p *Parser) parseParamDefList() *nodes.List {
	var params []nodes.Node
	for {
		param := p.parseParamDef()
		if param == nil {
			break
		}
		params = append(params, param)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: params}
}

// parseParamDef parses a single parameter definition.
func (p *Parser) parseParamDef() *nodes.ParamDef {
	if p.cur.Type != tokVARIABLE {
		return nil
	}

	loc := p.pos()
	param := &nodes.ParamDef{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume @param

	// Data type
	param.DataType = p.parseDataType()

	// Default value
	if p.cur.Type == '=' {
		p.advance()
		param.Default = p.parseExpr()
	}

	// OUTPUT
	if p.cur.Type == kwOUTPUT || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "out")) {
		param.Output = true
		p.advance()
	}

	// READONLY
	if p.cur.Type == kwREADONLY {
		param.ReadOnly = true
		p.advance()
	}

	param.Loc.End = p.pos()
	return param
}

// isRoutineOption checks if a token looks like a routine option keyword.
// Used to disambiguate WITH <option> from WITH in other contexts.
func (p *Parser) isRoutineOption(tok Token) bool {
	// Check keyword token types
	switch tok.Type {
	case kwSCHEMABINDING, kwEXEC, kwEXECUTE, kwRETURNS:
		return true
	}
	// Check string values for context-sensitive keywords
	if tok.Str != "" {
		s := strings.ToUpper(tok.Str)
		switch s {
		case "RECOMPILE", "ENCRYPTION", "NATIVE_COMPILATION",
			"VIEW_METADATA", "INLINE", "CALLED",
			"SCHEMABINDING", "EXECUTE", "EXEC", "RETURNS":
			return true
		}
	}
	return false
}

// parseRoutineOptionList parses a comma-separated list of routine options.
//
//	<procedure_option> ::=
//	    ENCRYPTION | RECOMPILE | EXECUTE AS { CALLER | SELF | OWNER | 'user_name' }
//	    | NATIVE_COMPILATION | SCHEMABINDING
//
//	<function_option> ::=
//	    ENCRYPTION | SCHEMABINDING | RETURNS NULL ON NULL INPUT
//	    | CALLED ON NULL INPUT | EXECUTE AS Clause
//	    | INLINE = { ON | OFF } | NATIVE_COMPILATION
//
//	<view_option> ::=
//	    ENCRYPTION | SCHEMABINDING | VIEW_METADATA
func (p *Parser) parseRoutineOptionList() *nodes.List {
	var items []nodes.Node
	for {
		opt := p.parseRoutineOption()
		if opt == nil {
			break
		}
		items = append(items, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}

// parseRoutineOption parses a single routine option.
func (p *Parser) parseRoutineOption() *nodes.String {
	// EXECUTE AS { CALLER | SELF | OWNER | 'user_name' }
	if p.cur.Type == kwEXECUTE || p.cur.Type == kwEXEC {
		p.advance() // consume EXECUTE/EXEC
		if p.cur.Type == kwAS {
			p.advance() // consume AS
		}
		// Principal: CALLER, SELF, OWNER, or 'string'
		principal := ""
		if p.cur.Type == tokSCONST {
			principal = p.cur.Str
			p.advance()
		} else if p.isIdentLike() {
			principal = p.cur.Str
			p.advance()
		}
		return &nodes.String{Str: "EXECUTE AS " + principal}
	}

	// RETURNS NULL ON NULL INPUT
	if p.cur.Type == kwRETURNS {
		p.advance() // consume RETURNS
		// NULL ON NULL INPUT
		if p.cur.Type == kwNULL {
			p.advance() // NULL
			if p.cur.Type == kwON {
				p.advance() // ON
			}
			if p.cur.Type == kwNULL {
				p.advance() // NULL
			}
			if p.isIdentLike() && strings.EqualFold(p.cur.Str, "INPUT") {
				p.advance() // INPUT
			}
			return &nodes.String{Str: "RETURNS NULL ON NULL INPUT"}
		}
		return &nodes.String{Str: "RETURNS"}
	}

	// CALLED ON NULL INPUT
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CALLED") {
		p.advance() // CALLED
		if p.cur.Type == kwON {
			p.advance() // ON
		}
		if p.cur.Type == kwNULL {
			p.advance() // NULL
		}
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "INPUT") {
			p.advance() // INPUT
		}
		return &nodes.String{Str: "CALLED ON NULL INPUT"}
	}

	// SCHEMABINDING
	if p.cur.Type == kwSCHEMABINDING {
		p.advance()
		return &nodes.String{Str: "SCHEMABINDING"}
	}

	// Simple identifier options: RECOMPILE, ENCRYPTION, VIEW_METADATA, NATIVE_COMPILATION, INLINE
	if p.isIdentLike() {
		s := strings.ToUpper(p.cur.Str)
		switch s {
		case "RECOMPILE", "ENCRYPTION", "VIEW_METADATA", "NATIVE_COMPILATION":
			p.advance()
			return &nodes.String{Str: s}
		case "INLINE":
			p.advance()
			// INLINE = { ON | OFF }
			if p.cur.Type == '=' {
				p.advance()
				val := "ON"
				if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.cur.Type == kwON {
					p.advance()
				}
				return &nodes.String{Str: "INLINE = " + val}
			}
			return &nodes.String{Str: "INLINE"}
		}
	}

	return nil
}
