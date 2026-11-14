package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// CREATE AI FUNCTION
// ---------------------------------------------------------------------------

// parseCreateAiFunctionStmt parses
//
//	CREATE [ OR REPLACE ] AI FUNCTION [ IF NOT EXISTS ] <name> (
//	    [ <arg_name> <arg_data_type> [ , ... ] ] )
//	  RETURNS <result_data_type>
//	  [ COMMENT = '<string_literal>' ]
//	  AS $$ <sql_expression_calling_an_ai_function> $$
//
//	CREATE [ OR REPLACE ] AI FUNCTION [ IF NOT EXISTS ] <name> (
//	    [ <arg_name> <arg_data_type> [ , ... ] ] )
//	  RETURNS <result_data_type>
//	  FROM EXPERIMENT <experiment_name> RUN <run_name>
//
// IF NOT EXISTS excludes OR REPLACE, as the docs grammar writes the two forms
// apart. The body may be a '...' string as well as a $$ block, as any routine
// body may. cur is AI; start is the CREATE token's Loc.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-ai-function
func (p *Parser) parseCreateAiFunctionStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume AI
	p.advance() // consume FUNCTION
	stmt := &ast.CreateAiFunctionStmt{OrReplace: orReplace, Loc: ast.Loc{Start: start.Start}}
	if !orReplace {
		if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
			return nil, err
		}
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	args, err := p.parseRoutineArgList(false)
	if err != nil {
		return nil, err
	}
	for _, a := range args {
		if a.Default != nil {
			return nil, &ParseError{Position: a.Loc.Start, End: a.Loc.End, Message: "syntax error: an AI function's argument takes no DEFAULT"}
		}
	}
	stmt.Args = args
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	if stmt.ReturnType, err = p.parseDataType(); err != nil {
		return nil, err
	}

	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		if err := p.expectWords("EXPERIMENT"); err != nil {
			return nil, err
		}
		if stmt.Experiment, err = p.parseObjectName(); err != nil {
			return nil, err
		}
		if err := p.expectWords("RUN"); err != nil {
			return nil, err
		}
		if stmt.Run, err = p.parseIdent(); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	if p.curIsWord("COMMENT") && p.startsObjectProperty() {
		opt, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if err := aiCheckOptions([]*ast.CopyOption{opt}, objectProps{"COMMENT": propString}); err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
	}
	if p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	body, dollar, end, err := p.parseRoutineBody()
	if err != nil {
		return nil, err
	}
	stmt.Body, stmt.BodyDollar = body, dollar
	stmt.Loc.End = end
	return stmt, nil
}
