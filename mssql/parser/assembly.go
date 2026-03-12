// Package parser - assembly.go implements T-SQL CREATE/ALTER/DROP ASSEMBLY parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateAssemblyStmt parses a CREATE ASSEMBLY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-assembly-transact-sql
//
//	CREATE ASSEMBLY assembly_name
//	    [ AUTHORIZATION owner_name ]
//	    FROM { <client_assembly_specifier> | <assembly_bits> [ ,...n ] }
//	    [ WITH PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE } ]
func (p *Parser) parseCreateAssemblyStmt() *nodes.CreateAssemblyStmt {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.CreateAssemblyStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			stmt.Authorization = p.cur.Str
			p.advance()
		}
	}

	// FROM file_path [, ...]
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FROM") {
		p.advance()
		var files []nodes.Node
		for {
			if p.cur.Type == tokSCONST {
				files = append(files, &nodes.String{Str: p.cur.Str})
				p.advance()
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(files) > 0 {
			stmt.FromFiles = &nodes.List{Items: files}
		}
	}

	// WITH PERMISSION_SET = { SAFE | EXTERNAL_ACCESS | UNSAFE }
	if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PERMISSION_SET") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.isIdentLike() {
				stmt.PermissionSet = strings.ToUpper(p.cur.Str)
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterAssemblyStmt parses an ALTER ASSEMBLY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-assembly-transact-sql
//
//	ALTER ASSEMBLY assembly_name
//	    [ FROM <client_assembly_specifier> | <assembly_bits> ]
//	    [ WITH <assembly_option> [ ,...n ] ]
//	    [ DROP FILE { file_name [ ,...n ] | ALL } ]
//	    [ ADD FILE FROM
//	        <client_file_specifier> [ AS <file_name> ][ ,...n ] ]
func (p *Parser) parseAlterAssemblyStmt() *nodes.AlterAssemblyStmt {
	loc := p.pos()
	// ASSEMBLY keyword already consumed by caller

	stmt := &nodes.AlterAssemblyStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var actions []nodes.Node

	// FROM ...
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FROM") {
		p.advance()
		if p.cur.Type == tokSCONST {
			actions = append(actions, &nodes.String{Str: "FROM=" + p.cur.Str})
			p.advance()
		}
	}

	// WITH PERMISSION_SET | VISIBILITY
	if p.cur.Type == kwWITH {
		p.advance()
		for p.isIdentLike() {
			opt := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.isIdentLike() {
					opt += "=" + strings.ToUpper(p.cur.Str)
					p.advance()
				}
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
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILE") {
			p.advance()
			actions = append(actions, &nodes.String{Str: "DROP FILE"})
			// consume file list
			for p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == kwALL {
				p.advance()
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	// ADD FILE FROM
	if _, ok := p.match(kwADD); ok {
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILE") {
			p.advance()
			actions = append(actions, &nodes.String{Str: "ADD FILE"})
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FROM") {
				p.advance()
				for p.cur.Type == tokSCONST {
					p.advance()
					if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
						p.advance()
						if p.cur.Type == tokSCONST || p.isIdentLike() {
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

	stmt.Loc.End = p.pos()
	return stmt
}
