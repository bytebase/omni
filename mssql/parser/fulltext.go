// Package parser - fulltext.go implements T-SQL FULLTEXT INDEX and CATALOG parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// fulltextIndexWithOptions defines valid option names for CREATE FULLTEXT INDEX ... WITH.
// Matches SqlScriptDOM fulltext index options: CHANGE_TRACKING, STOPLIST, SEARCH (PROPERTY LIST).
var fulltextIndexWithOptions = newOptionSet(kwCHANGE_TRACKING, kwSTOPLIST, kwSEARCH)

// fulltextCatalogWithOptions defines valid option names for CREATE FULLTEXT CATALOG ... WITH.
// Only ACCENT_SENSITIVITY is valid here.
var fulltextCatalogWithOptions = newOptionSet(kwACCENT_SENSITIVITY)

// fulltextCatalogActions defines valid action keywords for ALTER FULLTEXT CATALOG.
var fulltextCatalogActions = newOptionSet(kwREBUILD, kwREORGANIZE, kwAS)

// searchPropertyListWithOptions defines valid option names for ALTER SEARCH PROPERTY LIST ... WITH.
var searchPropertyListWithOptions = newOptionSet().withIdents("PROPERTY_SET_GUID", "PROPERTY_INT_ID", "PROPERTY_DESCRIPTION")

// parseCreateFulltextIndexStmt parses a CREATE FULLTEXT INDEX statement.
//
// BNF: mssql/parser/bnf/create-fulltext-index-transact-sql.bnf
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
//	    [ WITH [ ( ] <with_option> [ ,...n ] [ ) ] ]
//
//	<catalog_filegroup_option> ::=
//	    fulltext_catalog_name
//	  | ( fulltext_catalog_name , FILEGROUP filegroup_name )
//	  | ( FILEGROUP filegroup_name , fulltext_catalog_name )
//	  | ( FILEGROUP filegroup_name )
//
//	<with_option> ::=
//	    CHANGE_TRACKING [ = ] { MANUAL | AUTO | OFF [ , NO POPULATION ] }
//	  | STOPLIST [ = ] { OFF | SYSTEM | stoplist_name }
//	  | SEARCH PROPERTY LIST [ = ] property_list_name
func (p *Parser) parseCreateFulltextIndexStmt() (*nodes.CreateFulltextIndexStmt, error) {
	loc := p.pos()
	// INDEX keyword already consumed by caller

	stmt := &nodes.CreateFulltextIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// ON table_name
	if _, ok := p.match(kwON); ok {
		stmt.Table , _ = p.parseTableRef()
	}

	// Column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isAnyKeywordIdent() {
				col := p.cur.Str
				p.advance()
				// TYPE COLUMN type_col_name
				if p.cur.Type == kwTYPE {
					p.advance()
					if p.cur.Type == kwCOLUMN {
						p.advance()
					}
					if p.isAnyKeywordIdent() {
						p.advance()
					}
				}
				// LANGUAGE term
				if p.cur.Type == kwLANGUAGE {
					p.advance()
					if p.isAnyKeywordIdent() || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
						p.advance()
					}
				}
				// STATISTICAL_SEMANTICS
				if p.cur.Type == kwSTATISTICAL_SEMANTICS {
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
	if p.cur.Type == kwKEY {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
		if p.isAnyKeywordIdent() {
			stmt.KeyIndex = p.cur.Str
			p.advance()
		}
	}

	// ON <catalog_filegroup_option>
	//   catalog_name
	//   | ( catalog_name , FILEGROUP filegroup_name )
	//   | ( FILEGROUP filegroup_name , catalog_name )
	//   | ( FILEGROUP filegroup_name )
	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == '(' {
			// Parenthesized form
			p.advance()
			if p.cur.Type == kwFILEGROUP {
				p.advance()
				if p.isAnyKeywordIdent() {
					p.advance() // filegroup name (not captured in AST)
				}
				if _, ok := p.match(','); ok {
					if p.isIdentLike() {
						stmt.CatalogName = p.cur.Str
						p.advance()
					}
				}
			} else if p.isIdentLike() {
				stmt.CatalogName = p.cur.Str
				p.advance()
				if _, ok := p.match(','); ok {
					if p.cur.Type == kwFILEGROUP {
						p.advance()
						if p.isAnyKeywordIdent() {
							p.advance() // filegroup name
						}
					}
				}
			}
			p.match(')')
		} else if p.isIdentLike() {
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
			if p.isValidOption(fulltextIndexWithOptions) {
				opt := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if p.isAnyKeywordIdent() || p.cur.Type == kwON || p.cur.Type == kwOFF {
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterFulltextIndexStmt parses an ALTER FULLTEXT INDEX ON table action.
//
// BNF: mssql/parser/bnf/alter-fulltext-index-transact-sql.bnf
//
//	ALTER FULLTEXT INDEX ON table_name
//	    { ENABLE
//	    | DISABLE
//	    | SET CHANGE_TRACKING [ = ] { MANUAL | AUTO | OFF }
//	    | SET STOPLIST [ = ] { OFF | SYSTEM | stoplist_name } [ WITH NO POPULATION ]
//	    | SET SEARCH PROPERTY LIST [ = ] { OFF | property_list_name } [ WITH NO POPULATION ]
//	    | ADD ( column_name
//	          [ TYPE COLUMN type_column_name ]
//	          [ LANGUAGE language_term ]
//	          [ STATISTICAL_SEMANTICS ]
//	        [ ,...n ] ) [ WITH NO POPULATION ]
//	    | ALTER COLUMN column_name { ADD | DROP } STATISTICAL_SEMANTICS [ WITH NO POPULATION ]
//	    | DROP ( column_name [,...n] ) [ WITH NO POPULATION ]
//	    | START { FULL | INCREMENTAL | UPDATE } POPULATION
//	    | { STOP | PAUSE | RESUME } POPULATION
//	    }
func (p *Parser) parseAlterFulltextIndexStmt() (*nodes.AlterFulltextIndexStmt, error) {
	loc := p.pos()
	// INDEX keyword already consumed by caller

	stmt := &nodes.AlterFulltextIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// ON table_name
	if _, ok := p.match(kwON); ok {
		stmt.Table , _ = p.parseTableRef()
	}

	// Action dispatch
	switch {
	case p.cur.Type == kwENABLE:
		stmt.Action = "ENABLE"
		p.advance()

	case p.cur.Type == kwDISABLE:
		stmt.Action = "DISABLE"
		p.advance()

	case p.cur.Type == kwSET:
		stmt.Action = "SET"
		p.advance()
		if p.cur.Type == kwCHANGE_TRACKING {
			// SET CHANGE_TRACKING [ = ] { MANUAL | AUTO | OFF }
			p.advance()
			p.match('=') // optional =
			if p.cur.Type == kwMANUAL {
				stmt.ChangeTracking = "MANUAL"
				p.advance()
			} else if p.cur.Type == kwAUTO {
				stmt.ChangeTracking = "AUTO"
				p.advance()
			} else if p.cur.Type == kwOFF {
				stmt.ChangeTracking = "OFF"
				p.advance()
			}
		} else if p.cur.Type == kwSTOPLIST {
			// SET STOPLIST [ = ] { OFF | SYSTEM | stoplist_name } [ WITH NO POPULATION ]
			p.advance()
			p.match('=') // optional =
			var opts []nodes.Node
			if p.cur.Type == kwOFF {
				opts = append(opts, &nodes.String{Str: "STOPLIST=OFF"})
				p.advance()
			} else if p.cur.Type == kwSYSTEM {
				opts = append(opts, &nodes.String{Str: "STOPLIST=SYSTEM"})
				p.advance()
			} else if p.isAnyKeywordIdent() {
				opts = append(opts, &nodes.String{Str: "STOPLIST=" + p.cur.Str})
				p.advance()
			}
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
			stmt.WithNoPopulation = p.parseWithNoPopulation()
		} else if p.cur.Type == kwSEARCH {
			// SET SEARCH PROPERTY LIST [ = ] { OFF | property_list_name } [ WITH NO POPULATION ]
			p.advance() // consume SEARCH
			if p.cur.Type == kwPROPERTY {
				p.advance()
			}
			if p.cur.Type == kwLIST {
				p.advance()
			}
			p.match('=') // optional =
			var opts []nodes.Node
			if p.cur.Type == kwOFF {
				opts = append(opts, &nodes.String{Str: "SEARCH_PROPERTY_LIST=OFF"})
				p.advance()
			} else if p.isAnyKeywordIdent() {
				opts = append(opts, &nodes.String{Str: "SEARCH_PROPERTY_LIST=" + p.cur.Str})
				p.advance()
			}
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
			stmt.WithNoPopulation = p.parseWithNoPopulation()
		}

	case p.cur.Type == kwADD:
		stmt.Action = "ADD"
		p.advance()
		// ( column_name [ TYPE COLUMN ... ] [ LANGUAGE ... ] [ STATISTICAL_SEMANTICS ] [,...n] )
		if p.cur.Type == '(' {
			p.advance()
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isAnyKeywordIdent() {
					col := p.cur.Str
					p.advance()
					// TYPE COLUMN type_col_name
					if p.cur.Type == kwTYPE {
						p.advance()
						if p.cur.Type == kwCOLUMN {
							p.advance()
						}
						if p.isAnyKeywordIdent() {
							p.advance()
						}
					}
					// LANGUAGE term
					if p.cur.Type == kwLANGUAGE {
						p.advance()
						if p.isAnyKeywordIdent() || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
							p.advance()
						}
					}
					// STATISTICAL_SEMANTICS
					if p.cur.Type == kwSTATISTICAL_SEMANTICS {
						p.advance()
					}
					cols = append(cols, &nodes.String{Str: col})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(cols) > 0 {
				stmt.Columns = &nodes.List{Items: cols}
			}
		}
		// WITH NO POPULATION
		stmt.WithNoPopulation = p.parseWithNoPopulation()

	case p.cur.Type == kwALTER:
		stmt.Action = "ALTER"
		p.advance()
		// COLUMN column_name { ADD | DROP } STATISTICAL_SEMANTICS
		if p.cur.Type == kwCOLUMN {
			p.advance()
			if p.isIdentLike() {
				stmt.ColumnName = p.cur.Str
				p.advance()
			}
			if p.cur.Type == kwADD {
				stmt.ColumnAction = "ADD"
				p.advance()
			} else if p.cur.Type == kwDROP {
				stmt.ColumnAction = "DROP"
				p.advance()
			}
			// STATISTICAL_SEMANTICS
			if p.cur.Type == kwSTATISTICAL_SEMANTICS {
				p.advance()
			}
		}
		// WITH NO POPULATION
		stmt.WithNoPopulation = p.parseWithNoPopulation()

	case p.cur.Type == kwDROP:
		stmt.Action = "DROP"
		p.advance()
		// ( column_name [,...n] )
		if p.cur.Type == '(' {
			p.advance()
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isAnyKeywordIdent() {
					cols = append(cols, &nodes.String{Str: p.cur.Str})
					p.advance()
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(cols) > 0 {
				stmt.Columns = &nodes.List{Items: cols}
			}
		}
		// WITH NO POPULATION
		stmt.WithNoPopulation = p.parseWithNoPopulation()

	case p.cur.Type == kwSTART:
		stmt.Action = "START"
		p.advance()
		// { FULL | INCREMENTAL | UPDATE } POPULATION
		if p.isAnyKeywordIdent() {
			stmt.PopulationType = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		// POPULATION
		if p.cur.Type == kwPOPULATION {
			p.advance()
		}

	case p.cur.Type == kwSTOP:
		stmt.Action = "STOP"
		p.advance()
		if p.cur.Type == kwPOPULATION {
			p.advance()
		}

	case p.cur.Type == kwPAUSE:
		stmt.Action = "PAUSE"
		p.advance()
		if p.cur.Type == kwPOPULATION {
			p.advance()
		}

	case p.cur.Type == kwRESUME:
		stmt.Action = "RESUME"
		p.advance()
		if p.cur.Type == kwPOPULATION {
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseWithNoPopulation checks for and consumes WITH NO POPULATION.
func (p *Parser) parseWithNoPopulation() bool {
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwNO {
			p.advance() // consume WITH
			p.advance() // consume NO
			if p.cur.Type == kwPOPULATION {
				p.advance()
			}
			return true
		}
	}
	return false
}

// parseCreateFulltextCatalogStmt parses a CREATE FULLTEXT CATALOG statement.
//
// BNF: mssql/parser/bnf/create-fulltext-catalog-transact-sql.bnf
//
//	CREATE FULLTEXT CATALOG catalog_name
//	    [ ON FILEGROUP filegroup ]
//	    [ IN PATH 'rootpath' ]
//	    [ WITH ACCENT_SENSITIVITY = { ON | OFF } ]
//	    [ AS DEFAULT ]
//	    [ AUTHORIZATION owner_name ]
func (p *Parser) parseCreateFulltextCatalogStmt() (*nodes.CreateFulltextCatalogStmt, error) {
	loc := p.pos()
	// CATALOG keyword already consumed by caller

	stmt := &nodes.CreateFulltextCatalogStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
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
		if p.cur.Type == kwFILEGROUP {
			p.advance()
			if p.isAnyKeywordIdent() {
				opts = append(opts, &nodes.String{Str: "FILEGROUP=" + p.cur.Str})
				p.advance()
			}
		}
	}

	// IN PATH
	if p.cur.Type == kwIN {
		p.advance()
		if p.cur.Type == kwPATH {
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
		if p.isValidOption(fulltextCatalogWithOptions) {
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
	if p.cur.Type == kwAS {
		p.advance()
		if p.cur.Type == kwDEFAULT {
			p.advance()
			opts = append(opts, &nodes.String{Str: "AS DEFAULT"})
		}
	}

	// AUTHORIZATION
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			opts = append(opts, &nodes.String{Str: "AUTHORIZATION=" + p.cur.Str})
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateFulltextStoplistStmt parses a CREATE FULLTEXT STOPLIST statement.
//
// BNF: mssql/parser/bnf/create-fulltext-stoplist-transact-sql.bnf
//
//	CREATE FULLTEXT STOPLIST stoplist_name
//	    [ FROM { [ database_name . ] source_stoplist_name } | SYSTEM STOPLIST ]
//	    [ AUTHORIZATION owner_name ]
func (p *Parser) parseCreateFulltextStoplistStmt() (*nodes.CreateFulltextStoplistStmt, error) {
	loc := p.pos()
	// STOPLIST keyword already consumed by caller

	stmt := &nodes.CreateFulltextStoplistStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// stoplist_name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FROM clause
	if p.cur.Type == kwFROM {
		p.advance()
		// SYSTEM STOPLIST
		if p.cur.Type == kwSYSTEM {
			p.advance()
			if p.cur.Type == kwSTOPLIST {
				p.advance()
			}
			stmt.SystemStoplist = true
		} else if p.isAnyKeywordIdent() {
			// [ database_name . ] source_stoplist_name
			name1 := p.cur.Str
			p.advance()
			if p.cur.Type == '.' {
				p.advance()
				stmt.SourceDB = name1
				if p.isAnyKeywordIdent() {
					stmt.SourceList = p.cur.Str
					p.advance()
				}
			} else {
				stmt.SourceList = name1
			}
		}
	}

	// AUTHORIZATION owner_name
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			stmt.Authorization = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterFulltextStoplistStmt parses an ALTER FULLTEXT STOPLIST statement.
//
// BNF: mssql/parser/bnf/alter-fulltext-stoplist-transact-sql.bnf
//
//	ALTER FULLTEXT STOPLIST stoplist_name
//	{
//	    ADD [N] 'stopword' LANGUAGE language_term
//	  | DROP
//	    {
//	        'stopword' LANGUAGE language_term
//	      | ALL LANGUAGE language_term
//	      | ALL
//	    }
//	}
func (p *Parser) parseAlterFulltextStoplistStmt() (*nodes.AlterFulltextStoplistStmt, error) {
	loc := p.pos()
	// STOPLIST keyword already consumed by caller

	stmt := &nodes.AlterFulltextStoplistStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// stoplist_name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ADD or DROP
	if p.cur.Type == kwADD {
		stmt.Action = "ADD"
		p.advance()

		// 'stopword' or N'stopword'
		if p.cur.Type == tokNSCONST {
			stmt.IsNStr = true
			stmt.Stopword = p.cur.Str
			p.advance()
		} else if p.cur.Type == tokSCONST {
			stmt.Stopword = p.cur.Str
			p.advance()
		}

		// LANGUAGE language_term
		if p.cur.Type == kwLANGUAGE {
			p.advance()
			if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
				stmt.Language = p.cur.Str
				p.advance()
			}
		}
	} else if p.cur.Type == kwDROP {
		stmt.Action = "DROP"
		p.advance()

		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			// DROP [N]'stopword' LANGUAGE language_term
			if p.cur.Type == tokNSCONST {
				stmt.IsNStr = true
			}
			stmt.Stopword = p.cur.Str
			p.advance()
			if p.cur.Type == kwLANGUAGE {
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
					stmt.Language = p.cur.Str
					p.advance()
				}
			}
		} else if p.cur.Type == kwALL {
			p.advance()
			stmt.DropAll = true
			// ALL LANGUAGE language_term  or just ALL
			if p.cur.Type == kwLANGUAGE {
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
					stmt.Language = p.cur.Str
					p.advance()
				}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropFulltextStoplistStmt parses a DROP FULLTEXT STOPLIST statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-fulltext-stoplist-transact-sql
//
//	DROP FULLTEXT STOPLIST stoplist_name
func (p *Parser) parseDropFulltextStoplistStmt() (*nodes.DropFulltextStoplistStmt, error) {
	loc := p.pos()
	// STOPLIST keyword already consumed by caller

	stmt := &nodes.DropFulltextStoplistStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateSearchPropertyListStmt parses a CREATE SEARCH PROPERTY LIST statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-search-property-list-transact-sql
//
//	CREATE SEARCH PROPERTY LIST new_list_name
//	    [ FROM [ database_name . ] source_list_name ]
//	    [ AUTHORIZATION owner_name ]
func (p *Parser) parseCreateSearchPropertyListStmt() (*nodes.CreateSearchPropertyListStmt, error) {
	loc := p.pos()
	// LIST keyword already consumed by caller

	stmt := &nodes.CreateSearchPropertyListStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// list name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FROM clause
	if p.cur.Type == kwFROM {
		p.advance()
		if p.isAnyKeywordIdent() {
			name1 := p.cur.Str
			p.advance()
			if p.cur.Type == '.' {
				p.advance()
				stmt.SourceDB = name1
				if p.isAnyKeywordIdent() {
					stmt.SourceList = p.cur.Str
					p.advance()
				}
			} else {
				stmt.SourceList = name1
			}
		}
	}

	// AUTHORIZATION owner_name
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			stmt.Authorization = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterSearchPropertyListStmt parses an ALTER SEARCH PROPERTY LIST statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-search-property-list-transact-sql
//
//	ALTER SEARCH PROPERTY LIST list_name
//	{
//	    ADD 'property_name'
//	      WITH
//	      (
//	          PROPERTY_SET_GUID = 'property_set_guid'
//	        , PROPERTY_INT_ID = property_int_id
//	        [ , PROPERTY_DESCRIPTION = 'property_description' ]
//	      )
//	  | DROP 'property_name'
//	}
func (p *Parser) parseAlterSearchPropertyListStmt() (*nodes.AlterSearchPropertyListStmt, error) {
	loc := p.pos()
	// LIST keyword already consumed by caller

	stmt := &nodes.AlterSearchPropertyListStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// list name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ADD or DROP
	if p.cur.Type == kwADD {
		stmt.Action = "ADD"
		p.advance()

		// 'property_name'
		if p.cur.Type == tokSCONST {
			stmt.PropertyName = p.cur.Str
			p.advance()
		}

		// WITH ( ... )
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isValidOption(searchPropertyListWithOptions) {
						key := strings.ToUpper(p.cur.Str)
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
						}
						switch key {
						case "PROPERTY_SET_GUID":
							if p.cur.Type == tokSCONST {
								stmt.PropertySetGUID = p.cur.Str
								p.advance()
							}
						case "PROPERTY_INT_ID":
							if p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
								stmt.PropertyIntID = p.cur.Str
								p.advance()
							}
						case "PROPERTY_DESCRIPTION":
							if p.cur.Type == tokSCONST {
								stmt.PropertyDesc = p.cur.Str
								p.advance()
							}
						default:
							// unknown option - break out rather than silently skipping
						}
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				p.match(')')
			}
		}
	} else if p.cur.Type == kwDROP {
		stmt.Action = "DROP"
		p.advance()

		// 'property_name'
		if p.cur.Type == tokSCONST {
			stmt.PropertyName = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropSearchPropertyListStmt parses a DROP SEARCH PROPERTY LIST statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-search-property-list-transact-sql
//
//	DROP SEARCH PROPERTY LIST property_list_name
func (p *Parser) parseDropSearchPropertyListStmt() (*nodes.DropSearchPropertyListStmt, error) {
	loc := p.pos()
	// LIST keyword already consumed by caller

	stmt := &nodes.DropSearchPropertyListStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterFulltextCatalogStmt parses an ALTER FULLTEXT CATALOG statement.
//
// BNF: mssql/parser/bnf/alter-fulltext-catalog-transact-sql.bnf
//
//	ALTER FULLTEXT CATALOG catalog_name
//	    { REBUILD [ WITH ACCENT_SENSITIVITY = { ON | OFF } ]
//	    | REORGANIZE
//	    | AS DEFAULT
//	    }
func (p *Parser) parseAlterFulltextCatalogStmt() (*nodes.AlterFulltextCatalogStmt, error) {
	loc := p.pos()
	// CATALOG keyword already consumed by caller

	stmt := &nodes.AlterFulltextCatalogStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// catalog name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Action: REBUILD, REORGANIZE, AS DEFAULT
	if p.isValidOption(fulltextCatalogActions) {
		action := strings.ToUpper(p.cur.Str)
		p.advance()
		if action == "AS" && p.cur.Type == kwDEFAULT {
			action = "AS DEFAULT"
			p.advance()
		} else if action == "REBUILD" {
			// WITH ACCENT_SENSITIVITY = { ON | OFF }
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwACCENT_SENSITIVITY {
					optName := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						if p.cur.Type == kwON || p.cur.Type == kwOFF {
							optName += "=" + strings.ToUpper(p.cur.Str)
							p.advance()
						}
					}
					stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: optName}}}
				}
			}
		}
		stmt.Action = action
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
