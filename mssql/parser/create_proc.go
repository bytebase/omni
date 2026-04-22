// Package parser - create_proc.go implements T-SQL CREATE PROCEDURE/FUNCTION statement parsing.
package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateProcedureStmt parses a CREATE [OR ALTER] PROCEDURE statement.
//
// BNF: mssql/parser/bnf/create-procedure-transact-sql.bnf
//
//	CREATE [ OR ALTER ] { PROC | PROCEDURE }
//	    [ schema_name. ] procedure_name [ ; number ]
//	    [ { @parameter_name [ type_schema_name. ] data_type }
//	      [ VARYING ] [ NULL ] [ = default ] [ OUT | OUTPUT | READONLY ]
//	    ] [ ,...n ]
//	[ WITH <procedure_option> [ ,...n ] ]
//	[ FOR REPLICATION ]
//	AS { [ BEGIN ] sql_statement [;] [ ...n ] [ END ] }
//	[;]
//
//	-- CLR syntax:
//	CREATE [ OR ALTER ] { PROC | PROCEDURE }
//	    [ schema_name. ] procedure_name [ ; number ]
//	    [ { @parameter_name [ type_schema_name. ] data_type }
//	      [ = default ] [ OUT | OUTPUT ] [ READONLY ]
//	    ] [ ,...n ]
//	[ WITH EXECUTE AS Clause ]
//	AS { EXTERNAL NAME assembly_name.class_name.method_name }
//	[;]
//
//	<procedure_option> ::=
//	    [ ENCRYPTION ]
//	    [ RECOMPILE ]
//	    [ EXECUTE AS { CALLER | SELF | OWNER | 'user_name' } ]
//	    [ NATIVE_COMPILATION ]
//	    [ SCHEMABINDING ]
func (p *Parser) parseCreateProcedureStmt(orAlter bool) (*nodes.CreateProcedureStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER PROCEDURE → identifier or proc name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("proc_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateProcedureStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Procedure name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional procedure number: ; number
	if p.cur.Type == ';' {
		next := p.peekNext()
		if next.Type == tokICONST {
			p.advance() // consume ;
			if n, err := strconv.Atoi(p.cur.Str); err == nil {
				stmt.Number = n
			}
			p.advance() // consume number
		}
	}

	// Parameters (may or may not be in parentheses per BNF)
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ')' {
			params, err := p.parseParamDefList()
			if err != nil {
				return nil, err
			}
			stmt.Params = params
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == tokVARIABLE {
		params, err := p.parseParamDefList()
		if err != nil {
			return nil, err
		}
		stmt.Params = params
	}

	// Completion: after CREATE PROCEDURE p [params] → AS/WITH
	if p.collectMode() {
		p.addTokenCandidate(kwAS)
		p.addTokenCandidate(kwWITH)
		return nil, errCollecting
	}

	// WITH <procedure_option> [,...n]
	// In CREATE PROC context, WITH always introduces procedure options (never a CTE).
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err := p.parseRoutineOptionList()
		if err != nil {
			return nil, err
		}
		if opts == nil || len(opts.Items) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Options = opts
	}

	// FOR REPLICATION
	if p.cur.Type == kwFOR {
		next := p.peekNext()
		if next.Type == kwREPLICATION {
			p.advance() // consume FOR
			p.advance() // consume REPLICATION
			stmt.ForReplication = true
		}
	}

	// AS
	p.match(kwAS)

	// Completion: after CREATE PROCEDURE p AS → statement keywords
	if p.collectMode() {
		for _, kw := range topLevelKeywords {
			p.addTokenCandidate(kw)
		}
		return nil, errCollecting
	}

	// EXTERNAL NAME assembly.class.method (CLR)
	if p.cur.Type == kwEXTERNAL {
		p.advance() // consume EXTERNAL
		if p.cur.Type == kwNAME {
			p.advance() // consume NAME
		}
		stmt.ExternalName = p.parseMethodSpecifier()
	} else {
		// Body (BEGIN...END block or single statement)
		body, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseMethodSpecifier parses assembly_name.class_name.method_name.
func (p *Parser) parseMethodSpecifier() string {
	var parts []string
	for {
		if p.isIdentLike() {
			parts = append(parts, p.cur.Str)
			p.advance()
		} else {
			break
		}
		if p.cur.Type != '.' {
			break
		}
		p.advance() // consume dot
	}
	return strings.Join(parts, ".")
}

// parseCreateFunctionStmt parses a CREATE [OR ALTER] FUNCTION statement.
//
// BNF: mssql/parser/bnf/create-function-transact-sql.bnf
//
//	-- Scalar Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS return_data_type
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN scalar_expression
//	    END
//	[ ; ]
//
//	-- Inline Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS TABLE
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    RETURN [ ( ] select_stmt [ ) ]
//	[ ; ]
//
//	-- Multi-Statement Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS @return_variable TABLE <table_type_definition>
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN
//	    END
//	[ ; ]
//
//	-- CLR Scalar Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] }
//	    [ , ...n ]
//	)
//	RETURNS { return_data_type }
//	    [ WITH <clr_function_option> [ , ...n ] ]
//	    [ AS ] EXTERNAL NAME <method_specifier>
//	[ ; ]
//
//	-- CLR Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] }
//	    [ , ...n ]
//	)
//	RETURNS TABLE <clr_table_type_definition>
//	    [ WITH <clr_function_option> [ , ...n ] ]
//	    [ ORDER ( <order_clause> ) ]
//	    [ AS ] EXTERNAL NAME <method_specifier>
//	[ ; ]
//
//	<function_option> ::=
//	{
//	    [ ENCRYPTION ]
//	  | [ SCHEMABINDING ]
//	  | [ RETURNS NULL ON NULL INPUT | CALLED ON NULL INPUT ]
//	  | [ EXECUTE_AS_Clause ]
//	  | [ INLINE = { ON | OFF } ]
//	}
//
//	<method_specifier> ::= assembly_name.class_name.method_name
func (p *Parser) parseCreateFunctionStmt(orAlter bool) (*nodes.CreateFunctionStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER FUNCTION → identifier or func name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("func_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateFunctionStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Function name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parameters in parentheses
	if p.cur.Type == '(' {
		p.advance()
		// Completion: inside function params → @ variable
		if p.collectMode() {
			p.addRuleCandidate("variable")
			return nil, errCollecting
		}
		if p.cur.Type != ')' {
			params, err := p.parseParamDefList()
			if err != nil {
				return nil, err
			}
			stmt.Params = params
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// RETURNS type or RETURNS TABLE
	if _, ok := p.match(kwRETURNS); ok {
		// Completion: after RETURNS → type_name
		if p.collectMode() {
			p.addRuleCandidate("type_name")
			p.addTokenCandidate(kwTABLE)
			return nil, errCollecting
		}
		if p.cur.Type == kwTABLE {
			p.advance()
			stmt.ReturnsTable = &nodes.ReturnsTableDef{
				Loc: nodes.Loc{Start: p.pos(), End: -1},
			}
			// Check for table definition (CLR or multi-statement TVF)
			if p.cur.Type == '(' {
				p.advance()
				var cols []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					col, err := p.parseColumnDef()
					if err != nil {
						return nil, err
					}
					if col != nil {
						cols = append(cols, col)
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				stmt.ReturnsTable.Columns = &nodes.List{Items: cols}
			}
			stmt.ReturnsTable.Loc.End = p.prevEnd()
		} else if p.cur.Type == tokVARIABLE {
			// RETURNS @var TABLE (...)
			varName := p.cur.Str
			p.advance() // consume @var
			if p.cur.Type == kwTABLE {
				p.advance() // consume TABLE
				stmt.ReturnsTable = &nodes.ReturnsTableDef{
					Variable: varName,
					Loc:      nodes.Loc{Start: p.pos(), End: -1},
				}
				if p.cur.Type == '(' {
					p.advance()
					var cols []nodes.Node
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						col, err := p.parseColumnDef()
						if err != nil {
							return nil, err
						}
						if col != nil {
							cols = append(cols, col)
						}
						if _, ok := p.match(','); !ok {
							break
						}
					}
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
					stmt.ReturnsTable.Columns = &nodes.List{Items: cols}
				}
				stmt.ReturnsTable.Loc.End = p.prevEnd()
			}
		} else {
			dt, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.ReturnType = dt
		}
	}

	// WITH <function_option> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isRoutineOption(next) {
			p.advance() // consume WITH
			opts, err := p.parseRoutineOptionList()
			if err != nil {
				return nil, err
			}
			stmt.Options = opts
		}
	}

	// ORDER ( column ASC|DESC [,...n] ) for CLR table-valued functions
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == '(' {
			p.advance()
			orderBy, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			stmt.OrderClause = orderBy
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}

	// AS
	p.match(kwAS)

	// EXTERNAL NAME assembly.class.method (CLR)
	if p.cur.Type == kwEXTERNAL {
		p.advance() // consume EXTERNAL
		if p.cur.Type == kwNAME {
			p.advance() // consume NAME
		}
		stmt.ExternalName = p.parseMethodSpecifier()
	} else if p.cur.Type == kwRETURN {
		// Inline table-valued function: RETURN [ ( ] select_stmt [ ) ]
		retLoc := p.pos()
		p.advance() // consume RETURN
		hasParen := false
		if p.cur.Type == '(' {
			hasParen = true
			p.advance()
		}
		selectStmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if hasParen {
			p.match(')')
		}
		stmt.Body = &nodes.ReturnStmt{
			Value: &nodes.SubqueryExpr{Query: selectStmt, Loc: nodes.Loc{Start: retLoc, End: -1}},
			Loc:   nodes.Loc{Start: retLoc, End: -1},
		}
	} else {
		body, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseParamDefList parses a comma-separated list of parameter definitions.
//
//	param_def_list = param_def { ',' param_def }
//	param_def = @name [AS] type [VARYING] [NULL] [= default] [OUT|OUTPUT|READONLY]
func (p *Parser) parseParamDefList() (*nodes.List, error) {
	var params []nodes.Node
	for {
		param, err := p.parseParamDef()
		if err != nil {
			return nil, err
		}
		if param == nil {
			break
		}
		params = append(params, param)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: params}, nil
}

// parseParamDef parses a single parameter definition.
//
//	@parameter_name [ AS ] [ type_schema_name. ] data_type
//	    [ VARYING ] [ NULL ] [ = default ] [ OUT | OUTPUT | READONLY ]
func (p *Parser) parseParamDef() (*nodes.ParamDef, error) {
	if p.cur.Type != tokVARIABLE {
		return nil, nil
	}

	loc := p.pos()
	param := &nodes.ParamDef{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}
	p.advance() // consume @param

	// Optional AS keyword (for function parameters)
	p.match(kwAS)

	// Data type
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	param.DataType = dt

	// VARYING (for cursor parameters)
	if p.cur.Type == kwVARYING {
		param.Varying = true
		p.advance()
	}

	// NULL
	if p.cur.Type == kwNULL {
		param.Null = true
		p.advance()
	}

	// Default value
	if p.cur.Type == '=' {
		p.advance()
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		param.Default = def
	}

	// OUTPUT / OUT
	if p.cur.Type == kwOUTPUT || p.cur.Type == kwOUT {
		param.Output = true
		p.advance()
	}

	// READONLY
	if p.cur.Type == kwREADONLY {
		param.ReadOnly = true
		p.advance()
	}

	param.Loc.End = p.prevEnd()
	return param, nil
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
func (p *Parser) parseRoutineOptionList() (*nodes.List, error) {
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
	return &nodes.List{Items: items}, nil
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
			if p.cur.Type == kwINPUT {
				p.advance() // INPUT
			}
			return &nodes.String{Str: "RETURNS NULL ON NULL INPUT"}
		}
		return &nodes.String{Str: "RETURNS"}
	}

	// CALLED ON NULL INPUT
	if p.cur.Type == kwCALLED {
		p.advance() // CALLED
		if p.cur.Type == kwON {
			p.advance() // ON
		}
		if p.cur.Type == kwNULL {
			p.advance() // NULL
		}
		if p.cur.Type == kwINPUT {
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
