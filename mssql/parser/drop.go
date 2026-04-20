// Package parser - drop.go implements T-SQL DROP statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseDropStmt parses a DROP statement.
//
// BNF (DROP TABLE): mssql/parser/bnf/drop-table-transact-sql.bnf
//
//	DROP TABLE [ IF EXISTS ] { database_name.schema_name.table_name
//	    | schema_name.table_name | table_name } [ ,...n ]
//
// BNF (DROP VIEW): mssql/parser/bnf/drop-view-transact-sql.bnf
//
//	DROP VIEW [ IF EXISTS ] [ schema_name . ] view_name [ ,...n ] [ ; ]
//
// BNF (DROP TRIGGER): mssql/parser/bnf/drop-trigger-transact-sql.bnf
//
//	-- DML Trigger:
//	DROP TRIGGER [ IF EXISTS ] [schema_name.]trigger_name [ ,...n ] [ ; ]
//
//	-- DDL Trigger:
//	DROP TRIGGER [ IF EXISTS ] trigger_name [ ,...n ]
//	ON { DATABASE | ALL SERVER } [ ; ]
//
//	-- Logon Trigger:
//	DROP TRIGGER [ IF EXISTS ] trigger_name [ ,...n ]
//	ON ALL SERVER
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
//
//	<drop_clustered_index_option> ::=
//	{
//	    MAXDOP = max_degree_of_parallelism
//	  | ONLINE = { ON | OFF }
//	  | MOVE TO { partition_scheme_name ( column_name )
//	            | filegroup_name
//	            | "default"
//	            }
//	  [ FILESTREAM_ON { partition_scheme_name
//	            | filestream_filegroup_name
//	            | "default" } ]
//	}
func (p *Parser) parseDropStmt() (*nodes.DropStmt, error) {
	loc := p.pos()
	p.advance() // consume DROP

	stmt := &nodes.DropStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
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
		if p.isAnyKeywordIdent() {
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
			if p.isAnyKeywordIdent() {
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
		if p.cur.Type == kwFULLTEXT {
			p.advance()
			if p.cur.Type == kwINDEX {
				stmt.ObjectType = nodes.DropFulltextIndex
				p.advance()
			} else if p.cur.Type == kwCATALOG {
				stmt.ObjectType = nodes.DropFulltextCatalog
				p.advance()
			}
		}
		// Check for XML SCHEMA COLLECTION
		if p.cur.Type == kwXML {
			p.advance()
			if p.cur.Type == kwSCHEMA {
				p.advance()
				if p.cur.Type == kwCOLLECTION {
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

	// Completion: after DROP <type> [IF EXISTS] → appropriate rule
	if p.collectMode() {
		switch stmt.ObjectType {
		case nodes.DropTable:
			p.addRuleCandidate("table_ref")
		case nodes.DropView, nodes.DropMaterializedView:
			p.addRuleCandidate("view_name")
		case nodes.DropIndex:
			p.addRuleCandidate("index_name")
		case nodes.DropProcedure:
			p.addRuleCandidate("proc_name")
		case nodes.DropFunction:
			p.addRuleCandidate("func_name")
		case nodes.DropDatabase:
			p.addRuleCandidate("database_ref")
		case nodes.DropTrigger:
			p.addRuleCandidate("trigger_name")
		default:
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	// Names (comma-separated)
	var names []nodes.Node
	// For DROP INDEX: OnTables parallels Names; Items[i] is the table after ON
	// for `idx ON tbl` form, or nil for the backward-compatible `tbl.idx` form.
	var onTables []nodes.Node

	// DROP FULLTEXT INDEX ON table does not have a "name" per se - the table IS the target
	if stmt.ObjectType == nodes.DropFulltextIndex {
		// ON table_name
		if _, ok := p.match(kwON); ok {
			if ref, _ := p.parseTableRef(); ref != nil {
				names = append(names, ref)
			}
		}
	} else {
		for {
			name, _ := p.parseTableRef()
			if name == nil {
				break
			}
			// For DROP INDEX: index_name ON table_name [ WITH ( options ) ]
			// Backward-compatible form `tbl.idx` has no ON clause; keep nil placeholder
			// so OnTables.Items[i] aligns with Names.Items[i].
			if stmt.ObjectType == nodes.DropIndex {
				var onTable nodes.Node
				if _, ok := p.match(kwON); ok {
					ref, _ := p.parseTableRef()
					if ref != nil {
						onTable = ref
					}
				}
				onTables = append(onTables, onTable)
				// WITH ( <drop_clustered_index_option> [ ,...n ] )
				if p.cur.Type == kwWITH {
					p.advance()
					if p.cur.Type == '(' {
						stmt.Options, _ = p.parseAlterIndexOptions()
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
	if stmt.ObjectType == nodes.DropIndex && len(onTables) > 0 {
		stmt.OnTables = &nodes.List{Items: onTables}
	}

	// DROP TRIGGER ... ON { DATABASE | ALL SERVER } (for DDL/logon triggers)
	if stmt.ObjectType == nodes.DropTrigger {
		if _, ok := p.match(kwON); ok {
			if p.cur.Type == kwDATABASE {
				stmt.OnDatabase = true
				p.advance()
			} else if p.cur.Type == kwALL {
				p.advance()
				if p.matchIdentCI("SERVER") {
					stmt.OnAllServer = true
				}
			}
		}
	}

	// DROP ASSEMBLY ... WITH NO DEPENDENTS
	if stmt.ObjectType == nodes.DropAssembly && p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwNO {
			p.advance() // consume NO
			if p.cur.Type == kwDEPENDENTS {
				p.advance() // consume DEPENDENTS
				stmt.NoDependents = true
			}
		}
	}

	// Some DROP types also support CASCADE / RESTRICT
	if p.cur.Type == kwCASCADE {
		p.advance()
	} else if p.cur.Type == kwRESTRICT {
		p.advance()
	} else if p.cur.Type == kwCASCADE {
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
