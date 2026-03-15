// Package parser - drop.go implements T-SQL DROP statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseDropStmt parses a DROP statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-table-transact-sql
//
//	DROP TABLE|VIEW|INDEX|PROCEDURE|FUNCTION|DATABASE [IF EXISTS] name [, ...]
//
// BNF (DROP INDEX): mssql/parser/bnf/drop-index-transact-sql.bnf
//
//	DROP INDEX [ IF EXISTS ]
//	{ <drop_relational_or_xml_or_spatial_index> [ ,...n ]
//	| <drop_backward_compatible_index> [ ,...n ]
//	}
//
//	<drop_relational_or_xml_or_spatial_index> ::=
//	    index_name ON <object>
//	    [ WITH ( <drop_clustered_index_option> [ ,...n ] ) ]
//
//	<drop_backward_compatible_index> ::=
//	    [ owner_name. ] table_or_view_name.index_name
func (p *Parser) parseDropStmt() *nodes.DropStmt {
	loc := p.pos()
	p.advance() // consume DROP

	stmt := &nodes.DropStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Object type
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.DropTable
		p.advance()
	case kwMATERIALIZED:
		// DROP MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			stmt.ObjectType = nodes.DropMaterializedView
			p.advance()
		}
	case kwVIEW:
		stmt.ObjectType = nodes.DropView
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.DropIndex
		p.advance()
	case kwPROCEDURE, kwPROC:
		stmt.ObjectType = nodes.DropProcedure
		p.advance()
	case kwFUNCTION:
		stmt.ObjectType = nodes.DropFunction
		p.advance()
	case kwDATABASE:
		stmt.ObjectType = nodes.DropDatabase
		p.advance()
	case kwSCHEMA:
		stmt.ObjectType = nodes.DropSchema
		p.advance()
	case kwTRIGGER:
		stmt.ObjectType = nodes.DropTrigger
		p.advance()
	case kwTYPE:
		stmt.ObjectType = nodes.DropType
		p.advance()
	case kwSTATISTICS:
		// DROP STATISTICS is handled specially below
		stmt.ObjectType = nodes.DropStatistics
		p.advance()
	case kwRULE:
		stmt.ObjectType = nodes.DropRule
		p.advance()
	default:
		if p.isIdentLike() {
			switch strings.ToUpper(p.cur.Str) {
			case "SEQUENCE":
				stmt.ObjectType = nodes.DropSequence
				p.advance()
			case "SYNONYM":
				stmt.ObjectType = nodes.DropSynonym
				p.advance()
			case "ASSEMBLY":
				stmt.ObjectType = nodes.DropAssembly
				p.advance()
			case "DEFAULT":
				stmt.ObjectType = nodes.DropDefault
				p.advance()
			}
		}
		// Check for PARTITION FUNCTION/SCHEME
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if p.isIdentLike() {
				switch strings.ToUpper(p.cur.Str) {
				case "FUNCTION":
					stmt.ObjectType = nodes.DropPartitionFunction
					p.advance()
				case "SCHEME":
					stmt.ObjectType = nodes.DropPartitionScheme
					p.advance()
				}
			}
		}
		// Check for FULLTEXT INDEX/CATALOG
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "FULLTEXT") {
			p.advance()
			if p.cur.Type == kwINDEX {
				stmt.ObjectType = nodes.DropFulltextIndex
				p.advance()
			} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CATALOG") {
				stmt.ObjectType = nodes.DropFulltextCatalog
				p.advance()
			}
		}
		// Check for XML SCHEMA COLLECTION
		if p.cur.Type == kwXML {
			p.advance()
			if p.isIdentLike() && strings.EqualFold(p.cur.Str, "SCHEMA") {
				p.advance()
				if p.isIdentLike() && strings.EqualFold(p.cur.Str, "COLLECTION") {
					stmt.ObjectType = nodes.DropXmlSchemaCollection
					p.advance()
				}
			}
		}
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		}
	}

	// Names (comma-separated)
	var names []nodes.Node

	// DROP FULLTEXT INDEX ON table does not have a "name" per se - the table IS the target
	if stmt.ObjectType == nodes.DropFulltextIndex {
		// ON table_name
		if _, ok := p.match(kwON); ok {
			if ref := p.parseTableRef(); ref != nil {
				names = append(names, ref)
			}
		}
	} else {
		for {
			name := p.parseTableRef()
			if name == nil {
				break
			}
			// For DROP INDEX: index_name ON table_name [ WITH ( options ) ]
			if stmt.ObjectType == nodes.DropIndex {
				if _, ok := p.match(kwON); ok {
					p.parseTableRef() // consume table name
				}
				// WITH ( <drop_clustered_index_option> [ ,...n ] )
				if p.cur.Type == kwWITH {
					p.advance()
					if p.cur.Type == '(' {
						stmt.Options = p.parseAlterIndexOptions()
					}
				}
			}
			names = append(names, name)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	stmt.Names = &nodes.List{Items: names}

	// Some DROP types also support CASCADE / RESTRICT
	if p.cur.Type == kwCASCADE {
		p.advance()
	} else if p.cur.Type == kwRESTRICT {
		p.advance()
	} else if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "cascade") {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
