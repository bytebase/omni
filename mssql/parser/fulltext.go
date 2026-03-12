// Package parser - fulltext.go implements T-SQL FULLTEXT INDEX and CATALOG parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateFulltextIndexStmt parses a CREATE FULLTEXT INDEX statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-fulltext-index-transact-sql
//
//	CREATE FULLTEXT INDEX ON table_name
//	    [ ( { column_name
//	          [ TYPE COLUMN type_column_name ]
//	          [ LANGUAGE language_term ]
//	          [ STATISTICAL_SEMANTICS ]
//	        } [ ,...n ]
//	      ) ]
//	    KEY INDEX index_name
//	    [ ON <catalog_filegroup_option> ]
//	    [ WITH ( <with_option> [ ,...n ] ) ]
func (p *Parser) parseCreateFulltextIndexStmt() *nodes.CreateFulltextIndexStmt {
	loc := p.pos()
	// INDEX keyword already consumed by caller

	stmt := &nodes.CreateFulltextIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// ON table_name
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// Column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() {
				col := p.cur.Str
				p.advance()
				// TYPE COLUMN type_col_name
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TYPE") {
					p.advance()
					if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COLUMN") {
						p.advance()
					}
					if p.isIdentLike() {
						p.advance()
					}
				}
				// LANGUAGE term
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LANGUAGE") {
					p.advance()
					if p.isIdentLike() || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
						p.advance()
					}
				}
				// STATISTICAL_SEMANTICS
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STATISTICAL_SEMANTICS") {
					p.advance()
				}
				cols = append(cols, &nodes.String{Str: col})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		stmt.Columns = &nodes.List{Items: cols}
	}

	// KEY INDEX index_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "KEY") {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.KeyIndex = p.cur.Str
			p.advance()
		}
	}

	// ON catalog_name
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() {
			stmt.CatalogName = p.cur.Str
			p.advance()
		}
	}

	// WITH options (can be WITH (opt) or WITH opt without parens)
	if p.cur.Type == kwWITH {
		p.advance()
		useParens := p.cur.Type == '('
		if useParens {
			p.advance()
		}
		var opts []nodes.Node
		for (useParens && p.cur.Type != ')') || (!useParens && p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO) {
			if p.isIdentLike() {
				opt := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						opt += "=" + strings.ToUpper(p.cur.Str)
						p.advance()
					}
				}
				opts = append(opts, &nodes.String{Str: opt})
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if useParens {
			p.match(')')
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterFulltextIndexStmt parses an ALTER FULLTEXT INDEX ON table action.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-fulltext-index-transact-sql
//
//	ALTER FULLTEXT INDEX ON table_name
//	    { ENABLE | DISABLE
//	    | SET CHANGE_TRACKING { MANUAL | AUTO | OFF }
//	    | ADD ( column_name [...] ) [ WITH NO POPULATION ]
//	    | ALTER COLUMN column_name { ADD | DROP } STATISTICAL_SEMANTICS [ WITH NO POPULATION ]
//	    | DROP ( column_name [,...n] ) [ WITH NO POPULATION ]
//	    | START { FULL | INCREMENTAL | UPDATE } POPULATION
//	    | STOP POPULATION
//	    | PAUSE POPULATION
//	    | RESUME POPULATION
//	    }
func (p *Parser) parseAlterFulltextIndexStmt() *nodes.AlterFulltextIndexStmt {
	loc := p.pos()
	// INDEX keyword already consumed by caller

	stmt := &nodes.AlterFulltextIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// ON table_name
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// Action
	var action strings.Builder
	if p.isIdentLike() || p.cur.Type == kwSET || p.cur.Type == kwADD ||
		p.cur.Type == kwALTER || p.cur.Type == kwDROP {
		action.WriteString(strings.ToUpper(p.cur.Str))
		p.advance()
		// consume rest of action
		for p.isIdentLike() || p.cur.Type == kwCHECK || p.cur.Type == kwOFF {
			action.WriteString(" ")
			action.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
			if action.Len() > 50 {
				break
			}
		}
	}
	stmt.Action = action.String()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateFulltextCatalogStmt parses a CREATE FULLTEXT CATALOG statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-fulltext-catalog-transact-sql
//
//	CREATE FULLTEXT CATALOG catalog_name
//	    [ ON FILEGROUP filegroup ]
//	    [ IN PATH 'rootpath' ]
//	    [ WITH ACCENT_SENSITIVITY = { ON | OFF } ]
//	    [ AS DEFAULT ]
//	    [ AUTHORIZATION owner_name ]
func (p *Parser) parseCreateFulltextCatalogStmt() *nodes.CreateFulltextCatalogStmt {
	loc := p.pos()
	// CATALOG keyword already consumed by caller

	stmt := &nodes.CreateFulltextCatalogStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// catalog name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var opts []nodes.Node

	// ON FILEGROUP
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
			p.advance()
			if p.isIdentLike() {
				opts = append(opts, &nodes.String{Str: "FILEGROUP=" + p.cur.Str})
				p.advance()
			}
		}
	}

	// IN PATH
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "IN") {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PATH") {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "PATH=" + p.cur.Str})
			p.advance()
		}
	}

	// WITH ACCENT_SENSITIVITY
	if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentLike() {
			opt := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.cur.Type == kwON || p.cur.Type == kwOFF {
					opt += "=" + strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			opts = append(opts, &nodes.String{Str: opt})
		}
	}

	// AS DEFAULT
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DEFAULT") {
			p.advance()
			opts = append(opts, &nodes.String{Str: "AS DEFAULT"})
		}
	}

	// AUTHORIZATION
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "AUTHORIZATION=" + p.cur.Str})
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterFulltextCatalogStmt parses an ALTER FULLTEXT CATALOG statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-fulltext-catalog-transact-sql
//
//	ALTER FULLTEXT CATALOG catalog_name
//	    { REBUILD [ WITH ACCENT_SENSITIVITY = { ON | OFF } ]
//	    | REORGANIZE
//	    | AS DEFAULT
//	    }
func (p *Parser) parseAlterFulltextCatalogStmt() *nodes.AlterFulltextCatalogStmt {
	loc := p.pos()
	// CATALOG keyword already consumed by caller

	stmt := &nodes.AlterFulltextCatalogStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// catalog name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Action: REBUILD, REORGANIZE, AS DEFAULT
	if p.isIdentLike() {
		action := strings.ToUpper(p.cur.Str)
		p.advance()
		if action == "AS" && p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DEFAULT") {
			action = "AS DEFAULT"
			p.advance()
		} else if action == "REBUILD" {
			// WITH ACCENT_SENSITIVITY = ON|OFF
			if p.cur.Type == kwWITH {
				p.advance()
				for p.isIdentLike() || p.cur.Type == '=' || p.cur.Type == kwON || p.cur.Type == kwOFF {
					p.advance()
				}
			}
		}
		stmt.Action = action
	}

	stmt.Loc.End = p.pos()
	return stmt
}
