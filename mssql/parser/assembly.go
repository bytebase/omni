// Package parser - assembly.go implements T-SQL CREATE/ALTER/DROP ASSEMBLY parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

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
func (p *Parser) parseCreateAssemblyStmt() (*nodes.CreateAssemblyStmt, error) {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.CreateAssemblyStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			stmt.Authorization = p.cur.Str
			p.advance()
		}
	}

	// FROM { <client_assembly_specifier> | <assembly_bits> [ ,...n ] }
	if p.cur.Type == kwFROM {
		p.advance()
		var files []nodes.Node
		for {
			switch p.cur.Type {
			case tokSCONST, tokNSCONST:
				files = append(files, &nodes.String{Str: p.cur.Str})
				p.advance()
			case tokICONST, tokFCONST:
				// varbinary_literal (e.g., 0xABC123)
				files = append(files, &nodes.String{Str: p.cur.Str})
				p.advance()
			default:
				goto doneFiles
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	doneFiles:
		if len(files) > 0 {
			stmt.FromFiles = &nodes.List{Items: files}
		}
	}

	// WITH PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE }
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwPERMISSION_SET {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.isAnyKeywordIdent() {
				stmt.PermissionSet = strings.ToUpper(p.cur.Str)
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
func (p *Parser) parseAlterAssemblyStmt() (*nodes.AlterAssemblyStmt, error) {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.AlterAssemblyStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var actions []nodes.Node

	// FROM ...
	if p.cur.Type == kwFROM {
		p.advance()
		if p.cur.Type == tokSCONST {
			actions = append(actions, &nodes.String{Str: "FROM=" + p.cur.Str})
			p.advance()
		}
	}

	// WITH PERMISSION_SET | VISIBILITY | UNCHECKED DATA
	if p.cur.Type == kwWITH {
		p.advance()
		for p.isAnyKeywordIdent() {
			opt := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == kwON || p.cur.Type == kwOFF {
					opt += "=" + strings.ToUpper(p.cur.Str)
					p.advance()
				}
			} else if opt == "UNCHECKED" && p.cur.Type == kwDATA {
				// UNCHECKED DATA (two-word option)
				opt = "UNCHECKED DATA"
				p.advance()
			}
			actions = append(actions, &nodes.String{Str: opt})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	// DROP FILE
	if p.cur.Type == kwDROP {
		p.advance()
		if p.cur.Type == kwFILE {
			p.advance()
			actions = append(actions, &nodes.String{Str: "DROP FILE"})
			// consume file list
			for p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == kwALL {
				p.advance()
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	// ADD FILE FROM
	if _, ok := p.match(kwADD); ok {
		if p.cur.Type == kwFILE {
			p.advance()
			actions = append(actions, &nodes.String{Str: "ADD FILE"})
			if p.cur.Type == kwFROM {
				p.advance()
				for p.cur.Type == tokSCONST {
					p.advance()
					if p.cur.Type == kwAS {
						p.advance()
						if p.cur.Type == tokSCONST || p.isAnyKeywordIdent() {
							p.advance()
						}
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
			}
		}
	}

	if len(actions) > 0 {
		stmt.Actions = &nodes.List{Items: actions}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
