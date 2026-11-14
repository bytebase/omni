package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// DATA METRIC FUNCTION
// ---------------------------------------------------------------------------
//
// A data metric function (DMF) is a function over one or two tables. Its
// CREATE is read into an ast.CreateRoutineStmt of Kind
// RoutineDataMetricFunction, and its ALTER, which the docs write as ALTER
// FUNCTION with a signature of TABLE ( ... ) argument types, into an
// ast.AlterRoutineStmt with TableArgs. Both follow their pages strictly.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-data-metric-function
// https://docs.snowflake.com/en/sql-reference/sql/alter-function-dmf

// maxDataMetricTables is how many TABLE arguments a DMF takes: the docs
// write one, and an optional second for a referential check.
const maxDataMetricTables = 2

// startsCreateDataMetricFunction reports whether [ SECURE ] DATA METRIC FUNCTION
// starts at cur. It consumes nothing.
func (p *Parser) startsCreateDataMetricFunction() bool {
	return p.atWords("DATA", "METRIC", "FUNCTION") || p.atWords("SECURE", "DATA", "METRIC", "FUNCTION")
}

// parseCreateDataMetricFunction reads
//
//	CREATE [ OR REPLACE ] [ SECURE ] DATA METRIC FUNCTION [ IF NOT EXISTS ] <name>
//	  ( <table_arg> TABLE ( <col_arg> <data_type> [ , ... ] )
//	    [ , <table_arg> TABLE ( <col_arg> <data_type> [ , ... ] ) ] )
//	  RETURNS NUMBER [ [ NOT ] NULL ]
//	  [ LANGUAGE SQL ]
//	  [ COMMENT = '<string>' ]
//	  AS '<expression>'
//
// cur is SECURE or DATA; start is the CREATE token's Loc. IF NOT EXISTS
// excludes OR REPLACE, as the page writes the two forms apart.
func (p *Parser) parseCreateDataMetricFunction(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	if orAlter {
		return p.unsupported("CREATE")
	}
	stmt := &ast.CreateRoutineStmt{
		Kind:      ast.RoutineDataMetricFunction,
		OrReplace: orReplace,
		Secure:    p.acceptWords("SECURE"),
		Loc:       ast.Loc{Start: start.Start},
	}
	if err := p.expectWords("DATA", "METRIC", "FUNCTION"); err != nil {
		return nil, err
	}
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

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		arg, err := p.parseDataMetricTableArg()
		if err != nil {
			return nil, err
		}
		stmt.Args = append(stmt.Args, arg)
		if p.cur.Type != ',' || len(stmt.Args) == maxDataMetricTables {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	if !isWordToken(p.cur, "NUMBER") || p.peekNext().Type == '(' {
		return nil, p.syntaxErrorAtCur()
	}
	ret, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	stmt.ReturnType = ret
	switch {
	case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
		p.advance() // consume NOT
		p.advance() // consume NULL
		stmt.ReturnNotNull = true
	case p.cur.Type == kwNULL:
		p.advance() // consume NULL
		stmt.ReturnNull = true
	}

	if lang := p.cur; p.acceptWords("LANGUAGE") {
		if err := p.expectWords("SQL"); err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, &ast.CopyOption{
			Name: "LANGUAGE", Words: "SQL",
			Loc: ast.Loc{Start: lang.Loc.Start, End: p.prev.Loc.End},
		})
	}
	if p.cur.Type == kwCOMMENT {
		if p.peekNext().Type != '=' {
			return nil, syntaxErrorAtTok(p.peekNext())
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if !propString.matches(opt) {
			return nil, &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error at or near COMMENT"}
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
	stmt.Body = body
	stmt.BodyDollar = dollar
	stmt.Loc.End = end
	return stmt, nil
}

// parseDataMetricTableArg reads <arg> TABLE ( <col> <data_type> [ , ... ] ),
// which names at least one column.
func (p *Parser) parseDataMetricTableArg() (*ast.RoutineArg, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	if p.cur.Type == '(' && p.peekNext().Type == ')' {
		return nil, syntaxErrorAtTok(p.peekNext())
	}
	cols, err := p.parseRoutineTableColumns()
	if err != nil {
		return nil, err
	}
	return &ast.RoutineArg{Name: name, Table: cols, Loc: ast.Loc{Start: name.Loc.Start, End: p.prev.Loc.End}}, nil
}

// startsDataMetricSignature reports whether cur begins a DMF's ALTER
// signature, ( TABLE ( ... ): a plain function's argument type is never
// TABLE. It consumes nothing.
func (p *Parser) startsDataMetricSignature() bool {
	return p.cur.Type == '(' && p.peekNext().Type == kwTABLE
}

// parseAlterDataMetricFunction reads the rest of
//
//	ALTER FUNCTION [ IF EXISTS ] <name>
//	  ( TABLE ( <arg_type> [ , ... ] ) [ , TABLE ( <arg_type> [ , ... ] ) ] )
//	  { RENAME TO <new_name>
//	  | SET SECURE | UNSET SECURE
//	  | SET COMMENT = '<string>' | UNSET COMMENT
//	  | SET TAG <tag> = '<value>' [ , ... ]
//	  | UNSET TAG <tag> [ , ... ] }
//
// once stmt holds IF EXISTS and the name. cur is the signature's '('.
func (p *Parser) parseAlterDataMetricFunction(stmt *ast.AlterRoutineStmt) (ast.Node, error) {
	p.advance() // consume '('
	for {
		argStart := p.cur.Loc.Start
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		arg := &ast.RoutineTableArgType{}
		for {
			typ, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			arg.Types = append(arg.Types, typ)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		arg.Loc = ast.Loc{Start: argStart, End: p.prev.Loc.End}
		stmt.TableArgs = append(stmt.TableArgs, arg)
		if p.cur.Type != ',' || len(stmt.TableArgs) == maxDataMetricTables {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	switch {
	case p.cur.Type == kwRENAME:
		newName, err := p.parseRenameTo()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineRename
		stmt.NewName = newName
	case p.cur.Type == kwSET && p.peekNext().Type == kwTAG:
		p.advance() // consume SET
		p.advance() // consume TAG
		tags, err := p.parseObjectTagList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineSetTag
		stmt.Tags = tags
	case p.cur.Type == kwUNSET && p.peekNext().Type == kwTAG:
		p.advance() // consume UNSET
		p.advance() // consume TAG
		names, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineUnsetTag
		stmt.UnsetTags = names
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		switch {
		case isWordToken(p.cur, "SECURE"):
			tok := p.advance()
			stmt.Options = []*ast.CopyOption{{Name: "SECURE", Bare: true, Loc: tok.Loc}}
		case p.cur.Type == kwCOMMENT && p.peekNext().Type == '=':
			opt, err := p.parseObjectProperty()
			if err != nil {
				return nil, err
			}
			if !propString.matches(opt) {
				return nil, &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error at or near COMMENT"}
			}
			stmt.Options = []*ast.CopyOption{opt}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterRoutineSet
	case p.cur.Type == kwUNSET:
		p.advance() // consume UNSET
		if !isWordToken(p.cur, "SECURE") && !isWordToken(p.cur, "COMMENT") {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterRoutineUnset
		stmt.UnsetProps = []string{strings.ToUpper(p.advance().Str)}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
