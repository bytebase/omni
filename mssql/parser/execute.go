// Package parser - execute.go implements T-SQL EXEC/EXECUTE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseExecStmt parses an EXEC/EXECUTE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/execute-transact-sql
//
//	EXEC [@return_var =] proc_name [arg, @param = value [OUTPUT], ...]
func (p *Parser) parseExecStmt() *nodes.ExecStmt {
	loc := p.pos()
	p.advance() // consume EXEC or EXECUTE

	stmt := &nodes.ExecStmt{
		Loc: nodes.Loc{Start: loc},
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
	stmt.Name = p.parseTableRef()

	// Arguments
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != ')' &&
		p.cur.Type != kwGO && p.cur.Type != kwEND &&
		!p.isStatementStart() {
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExecArg parses a single EXEC argument.
//
//	exec_arg = [@param =] expr [OUTPUT|OUT]
func (p *Parser) parseExecArg() *nodes.ExecArg {
	loc := p.pos()

	arg := &nodes.ExecArg{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for named argument: @param = value
	if p.cur.Type == tokVARIABLE {
		next := p.peekNext()
		if next.Type == '=' {
			arg.Name = p.cur.Str
			p.advance() // consume @param
			p.advance() // consume =
		}
	}

	// Parse value expression
	arg.Value = p.parseExpr()
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
