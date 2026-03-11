// Package parser - create_proc.go implements T-SQL CREATE PROCEDURE/FUNCTION statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseCreateProcedureStmt parses a CREATE [OR ALTER] PROCEDURE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-procedure-transact-sql
//
//	CREATE [OR ALTER] PROCEDURE name @param type [= default] [OUTPUT], ... AS BEGIN ... END
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
//	CREATE [OR ALTER] FUNCTION name (@param type, ...) RETURNS type AS BEGIN ... END
//	CREATE [OR ALTER] FUNCTION name (@param type, ...) RETURNS TABLE AS RETURN SELECT ...
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
		} else {
			stmt.ReturnType = p.parseDataType()
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
