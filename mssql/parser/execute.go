// Package parser - execute.go implements T-SQL EXEC/EXECUTE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseExecStmt parses an EXEC/EXECUTE statement.
//
// BNF: mssql/parser/bnf/execute-transact-sql.bnf
//
//	-- Execute a stored procedure or function:
//	[ { EXEC | EXECUTE } ]
//	    {
//	      [ @return_status = ]
//	      { module_name [ ;number ] | @module_name_var }
//	        [ [ @parameter = ] { value
//	                           | @variable [ OUTPUT ]
//	                           | [ DEFAULT ]
//	                           }
//	        ]
//	      [ ,...n ]
//	      [ WITH <execute_option> [ ,...n ] ]
//	    }
//	[ ; ]
//
//	-- Execute a character string:
//	{ EXEC | EXECUTE }
//	    ( { @string_variable | [ N ]'tsql_string' } [ + ...n ] )
//	    [ AS { LOGIN | USER } = ' name ' ]
//	[ ; ]
//
//	-- Execute a pass-through command against a linked server:
//	{ EXEC | EXECUTE }
//	    ( { @string_variable | [ N ] 'command_string [ ? ]' } [ + ...n ]
//	        [ { , { value | @variable [ OUTPUT ] } } [ ...n ] ]
//	    )
//	    [ AS { LOGIN | USER } = ' name ' ]
//	    [ AT linked_server_name ]
//	    [ AT DATA_SOURCE data_source_name ]
//	[ ; ]
//
//	<execute_option> ::=
//	{
//	        RECOMPILE
//	    | { RESULT SETS UNDEFINED }
//	    | { RESULT SETS NONE }
//	    | { RESULT SETS ( <result_sets_definition> [,...n ] ) }
//	}
//
//	<result_sets_definition> ::=
//	{
//	    (
//	         { column_name
//	           data_type
//	         [ COLLATE collation_name ]
//	         [ NULL | NOT NULL ] }
//	         [,...n ]
//	    )
//	    | AS OBJECT
//	        [ db_name . [ schema_name ] . | schema_name . ]
//	        {table_name | view_name | table_valued_function_name }
//	    | AS TYPE [ schema_name.]table_type_name
//	    | AS FOR XML
//	}
func (p *Parser) parseExecStmt() (*nodes.ExecStmt, error) {
	loc := p.pos()
	p.advance() // consume EXEC or EXECUTE

	stmt := &nodes.ExecStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for EXEC ('string') syntax -- execute character string
	if p.cur.Type == '(' {
		stmt.ExecString = p.parseExecString()

		// AS { LOGIN | USER } = 'name'
		p.parseExecAsContext(stmt)

		// AT linked_server_name | AT DATA_SOURCE data_source_name
		p.parseExecAtClause(stmt)

		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Check for @return_var = proc_name
	if p.cur.Type == tokVARIABLE {
		next := p.peekNext()
		if next.Type == '=' {
			stmt.ReturnVar = p.cur.Str
			p.advance() // consume @var
			p.advance() // consume =
		}
	}

	// Procedure name
	stmt.Name , _ = p.parseTableRef()

	// Optional procedure number: ;number (EXEC sp_test;1)
	if p.cur.Type == ';' {
		next := p.peekNext()
		if next.Type == tokICONST {
			p.advance() // consume ;
			p.advance() // consume number (stored in Name context)
		}
	}

	// Arguments
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != ')' &&
		p.cur.Type != kwGO && p.cur.Type != kwEND &&
		!p.isStatementStart() && !p.isExecWithOption() {
		var args []nodes.Node
		for {
			arg := p.parseExecArg()
			if arg == nil {
				break
			}
			args = append(args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}
	}

	// WITH <execute_option> [,...n]
	if p.isExecWithOption() {
		p.advance() // consume WITH
		stmt.WithOptions = p.parseExecOptionList()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseExecString parses the EXEC ('string' + ...) form.
// Returns the string expression (possibly concatenated).
func (p *Parser) parseExecString() nodes.ExprNode {
	p.advance() // consume '('
	expr, _ := p.parseExpr()
	// Parse additional comma-separated parameters for linked server pass-through
	// ( 'string', param1, param2 )
	for p.cur.Type == ',' {
		p.advance() // consume ','
		p.parseExpr()
	}
	_, _ = p.expect(')')
	return expr
}

// parseExecAsContext parses AS { LOGIN | USER } = 'name'.
func (p *Parser) parseExecAsContext(stmt *nodes.ExecStmt) {
	if p.cur.Type != kwAS {
		return
	}
	next := p.peekNext()
	if next.Type != kwLOGIN && next.Type != kwUSER {
		return
	}
	p.advance() // consume AS
	if p.cur.Type == kwLOGIN {
		p.advance() // consume LOGIN
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.AsLogin = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwUSER {
		p.advance() // consume USER
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.AsUser = p.cur.Str
			p.advance()
		}
	}
}

// parseExecAtClause parses AT linked_server_name | AT DATA_SOURCE name.
func (p *Parser) parseExecAtClause(stmt *nodes.ExecStmt) {
	if p.cur.Type != kwAT {
		return
	}
	p.advance() // consume AT
	// DATA_SOURCE keyword
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "DATA_SOURCE") {
		p.advance() // consume DATA_SOURCE
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			stmt.AtDataSource = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLike() || p.cur.Type == tokIDENT {
		stmt.AtServer = p.cur.Str
		p.advance()
	}
}

// isExecWithOption returns true if current token is WITH followed by an exec option.
func (p *Parser) isExecWithOption() bool {
	if p.cur.Type != kwWITH {
		return false
	}
	next := p.peekNext()
	if next.Str != "" {
		s := strings.ToUpper(next.Str)
		if s == "RECOMPILE" || s == "RESULT" {
			return true
		}
	}
	return false
}

// parseExecOptionList parses WITH execute_option [, execute_option ...].
func (p *Parser) parseExecOptionList() *nodes.List {
	var items []nodes.Node
	for {
		opt := p.parseExecOption()
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

// parseExecOption parses a single execute option.
//
//	RECOMPILE | RESULT SETS UNDEFINED | RESULT SETS NONE
//	| RESULT SETS ( <result_sets_definition> [,...n] )
func (p *Parser) parseExecOption() *nodes.String {
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "RECOMPILE") {
		p.advance()
		return &nodes.String{Str: "RECOMPILE"}
	}
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "RESULT") {
		p.advance() // consume RESULT
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "SETS") {
			p.advance() // consume SETS
		}
		// UNDEFINED | NONE | ( ... )
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "UNDEFINED") {
			p.advance()
			return &nodes.String{Str: "RESULT SETS UNDEFINED"}
		}
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "NONE") {
			p.advance()
			return &nodes.String{Str: "RESULT SETS NONE"}
		}
		if p.cur.Type == '(' {
			// Skip over result set definitions for now (complex nested structure)
			depth := 0
			for p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
					if depth == 0 {
						p.advance()
						break
					}
				}
				p.advance()
			}
			return &nodes.String{Str: "RESULT SETS DEFINED"}
		}
		return &nodes.String{Str: "RESULT SETS"}
	}
	return nil
}

// parseExecArg parses a single EXEC argument.
//
//	exec_arg = [@param =] { expr | DEFAULT } [OUTPUT|OUT]
func (p *Parser) parseExecArg() *nodes.ExecArg {
	loc := p.pos()

	arg := &nodes.ExecArg{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for DEFAULT keyword
	if p.cur.Type == kwDEFAULT {
		arg.IsDefault = true
		p.advance()
		arg.Loc.End = p.pos()
		return arg
	}

	// Check for named argument: @param = value
	if p.cur.Type == tokVARIABLE {
		next := p.peekNext()
		if next.Type == '=' {
			arg.Name = p.cur.Str
			p.advance() // consume @param
			p.advance() // consume =
			// Check for DEFAULT after @param =
			if p.cur.Type == kwDEFAULT {
				arg.IsDefault = true
				p.advance()
				arg.Loc.End = p.pos()
				return arg
			}
		}
	}

	// Parse value expression
	arg.Value, _ = p.parseExpr()
	if arg.Value == nil {
		return nil
	}

	// Check for OUTPUT/OUT
	if p.cur.Type == kwOUTPUT || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "out")) {
		arg.Output = true
		p.advance()
	}

	arg.Loc.End = p.pos()
	return arg
}

// isStatementStart returns true if the current token starts a new statement.
func (p *Parser) isStatementStart() bool {
	switch p.cur.Type {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwMERGE,
		kwCREATE, kwALTER, kwDROP, kwTRUNCATE,
		kwDECLARE, kwSET, kwIF, kwWHILE, kwBEGIN,
		kwRETURN, kwBREAK, kwCONTINUE, kwGOTO,
		kwEXEC, kwEXECUTE, kwPRINT, kwRAISERROR, kwTHROW,
		kwGRANT, kwREVOKE, kwDENY, kwUSE, kwWAITFOR,
		kwWITH, kwGO:
		return true
	}
	return false
}
