// Package parser - assembly.go implements T-SQL CREATE/ALTER/DROP ASSEMBLY parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// assemblyPermissionSetValues matches SqlScriptDOM's AssemblyPermissionSet enum
// for CREATE/ALTER ASSEMBLY ... WITH PERMISSION_SET = ...
var assemblyPermissionSetValues = newOptionSet().withIdents("SAFE", "EXTERNAL_ACCESS", "UNSAFE")

// parseCreateAssemblyStmt parses a CREATE ASSEMBLY statement.
//
// BNF: mssql/parser/bnf/create-assembly-transact-sql.bnf
//
//	CREATE ASSEMBLY assembly_name
//	    [ AUTHORIZATION owner_name ]
//	    FROM { <client_assembly_specifier> | <assembly_bits> [ ,...n ] }
//	    [ WITH PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE } ]
//	    [ ; ]
//
//	<client_assembly_specifier> ::=
//	    '[ \\computer_name\ ] share_name\ [ path\ ] manifest_file_name'
//	    | '[ local_path\ ] manifest_file_name'
//
//	<assembly_bits> ::=
//	    { varbinary_literal | varbinary_expression }
//
// The grammar is createAssemblyStatement (TSql160.g:1050): the name and
// FROM expressionList are required.
func (p *Parser) parseCreateAssemblyStmt() (*nodes.CreateAssemblyStmt, error) {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.CreateAssemblyStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		owner, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Authorization = owner
	}

	// FROM { <client_assembly_specifier> | <assembly_bits> [ ,...n ] }
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	var files []nodes.Node
	for {
		f, err := p.assemblyParseSource()
		if err != nil {
			return nil, err
		}
		files = append(files, f)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.FromFiles = &nodes.List{Items: files}

	// WITH PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE }
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type != kwPERMISSION_SET {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		v, err := p.assemblyParsePermissionSet()
		if err != nil {
			return nil, err
		}
		stmt.PermissionSet = v
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// assemblyParseSource parses one FROM item of CREATE or ALTER ASSEMBLY, an
// expression (expressionList, alterAssemblyFromClause, TSql160.g:11193): a
// path, the assembly bits, or an expression that yields them. A literal
// keeps its value as the String; anything else keeps its source text.
func (p *Parser) assemblyParseSource() (*nodes.String, error) {
	start := p.pos()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if lit, ok := expr.(*nodes.Literal); ok && lit.Str != "" {
		return &nodes.String{Str: lit.Str}, nil
	}
	return &nodes.String{Str: p.source[start:p.prevEnd()]}, nil
}

// assemblyParsePermissionSet parses = { SAFE | EXTERNAL_ACCESS | UNSAFE }
// (assemblyPermissionSetOption, TSql160.g:11230) after PERMISSION_SET.
func (p *Parser) assemblyParsePermissionSet() (string, error) {
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	if !p.isValidOption(assemblyPermissionSetValues) || p.isDelimited(p.cur) {
		return "", p.syntaxErrorAtCur()
	}
	v := strings.ToUpper(p.cur.Str)
	p.advance()
	return v, nil
}

// parseAlterAssemblyStmt parses an ALTER ASSEMBLY statement.
//
// BNF: mssql/parser/bnf/alter-assembly-transact-sql.bnf
//
//	ALTER ASSEMBLY assembly_name
//	    [ FROM <client_assembly_specifier> | <assembly_bits> ]
//	    [ WITH <assembly_option> [ ,...n ] ]
//	    [ DROP FILE { file_name [ ,...n ] | ALL } ]
//	    [ ADD FILE FROM
//	        {
//	            client_file_specifier [ AS file_name ]
//	          | file_bits AS file_name
//	        } [ ,...n ]
//	    ] [ ; ]
//
//	<assembly_option> ::=
//	    PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE }
//	  | VISIBILITY = { ON | OFF }
//	  | UNCHECKED DATA
//
// The grammar is alterAssemblyStatement (TSql160.g:11150) and the rules
// under it: FROM takes one expression, DROP FILE ALL or string literals,
// ADD FILE FROM expressions each with an optional AS string literal. SQL
// Server wants at least one clause and each option once (Msg 102).
func (p *Parser) parseAlterAssemblyStmt() (*nodes.AlterAssemblyStmt, error) {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.AlterAssemblyStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	var actions []nodes.Node

	// FROM expression (alterAssemblyFromClause, TSql160.g:11193)
	if p.cur.Type == kwFROM {
		p.advance()
		src, err := p.assemblyParseSource()
		if err != nil {
			return nil, err
		}
		actions = append(actions, &nodes.String{Str: "FROM=" + src.Str})
	}

	// WITH PERMISSION_SET | VISIBILITY | UNCHECKED DATA (alterAssemblyWith,
	// TSql160.g:11197)
	if p.cur.Type == kwWITH {
		p.advance()
		seen := map[string]bool{}
		for {
			optLoc := p.pos()
			var opt string
			switch {
			case p.cur.Type == kwPERMISSION_SET:
				p.advance()
				v, err := p.assemblyParsePermissionSet()
				if err != nil {
					return nil, err
				}
				opt = "PERMISSION_SET=" + v
				if seen["PERMISSION_SET"] {
					return nil, p.newParseError(optLoc, "PERMISSION_SET is specified more than once")
				}
				seen["PERMISSION_SET"] = true
			case p.ddlIsWord("VISIBILITY"):
				p.advance()
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				switch p.cur.Type {
				case kwON:
					opt = "VISIBILITY=ON"
				case kwOFF:
					opt = "VISIBILITY=OFF"
				default:
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
				if seen["VISIBILITY"] {
					return nil, p.newParseError(optLoc, "VISIBILITY is specified more than once")
				}
				seen["VISIBILITY"] = true
			case p.ddlIsWord("UNCHECKED"):
				p.advance()
				if _, err := p.expect(kwDATA); err != nil {
					return nil, err
				}
				opt = "UNCHECKED DATA"
				if seen[opt] {
					return nil, p.newParseError(optLoc, "UNCHECKED DATA is specified more than once")
				}
				seen[opt] = true
			default:
				return nil, p.syntaxErrorAtCur()
			}
			actions = append(actions, &nodes.String{Str: opt})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	// DROP FILE { ALL | 'file_name' [ ,...n ] } (alterAssemblyDropFile,
	// TSql160.g:11266)
	if p.cur.Type == kwDROP {
		p.advance()
		if _, err := p.expect(kwFILE); err != nil {
			return nil, err
		}
		actions = append(actions, &nodes.String{Str: "DROP FILE"})
		if p.cur.Type == kwALL {
			p.advance()
		} else {
			for {
				if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	// ADD FILE FROM expression [ AS 'file_name' ] [ ,...n ]
	// (alterAssemblyAddFile, TSql160.g:11290)
	if p.cur.Type == kwADD {
		p.advance()
		if _, err := p.expect(kwFILE); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		actions = append(actions, &nodes.String{Str: "ADD FILE"})
		for {
			if _, err := p.assemblyParseSource(); err != nil {
				return nil, err
			}
			if _, ok := p.match(kwAS); ok {
				if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(actions) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Actions = &nodes.List{Items: actions}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
