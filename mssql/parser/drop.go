// Package parser - drop.go implements T-SQL DROP statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// DROP ident-based subcommand sets — matching SqlScriptDOM's DROP variants
// that are lexed as identifiers rather than dedicated keywords.
var (
	dropIdentDispatch     = newOptionSet().withIdents("SEQUENCE", "SYNONYM", "ASSEMBLY", "DEFAULT")
	dropPartitionDispatch = newOptionSet().withIdents("FUNCTION", "SCHEME")
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
		dropStart := p.cur
		// DROP SEQUENCE/SYNONYM/ASSEMBLY/DEFAULT dispatch — closed set.
		if p.isValidOption(dropIdentDispatch) {
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
		// PARTITION takes FUNCTION or SCHEME, and one of them is required:
		// DROP PARTITION a is Msg 102 (measured).
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if !p.isValidOption(dropPartitionDispatch) {
				return nil, p.syntaxErrorAtCur()
			}
			switch strings.ToUpper(p.cur.Str) {
			case "FUNCTION":
				stmt.ObjectType = nodes.DropPartitionFunction
				p.advance()
			case "SCHEME":
				stmt.ObjectType = nodes.DropPartitionScheme
				p.advance()
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
		// XML SCHEMA COLLECTION, whose three words are all required:
		// DROP XML SCHEMA a is Msg 102 (measured).
		if p.cur.Type == kwXML {
			p.advance()
			if _, err := p.expect(kwSCHEMA); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwCOLLECTION); err != nil {
				return nil, err
			}
			stmt.ObjectType = nodes.DropXmlSchemaCollection
		}
		// A DROP of an object type none of these names is Msg 343, "Unknown
		// object type 'a' used in a CREATE, DROP, or ALTER statement"
		// (measured): every arm above consumes the words it names, so a
		// token still at the front is not an object type.
		if p.cur == dropStart {
			return nil, p.newParseError(p.cur.Loc, "unknown object type after DROP")
		}
	}

	// IF EXISTS: a bare IF is not part of the statement, and SQL Server
	// refuses DROP DATABASE IF a with Msg 102 (measured).
	if p.cur.Type == kwIF {
		p.advance() // IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
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
			p.addRuleCandidate("index_ref")
		case nodes.DropProcedure:
			p.addRuleCandidate("proc_name")
			p.addRuleCandidate("proc_ref")
		case nodes.DropFunction:
			p.addRuleCandidate("func_name")
		case nodes.DropDatabase:
			p.addRuleCandidate("database_ref")
		case nodes.DropTrigger:
			p.addRuleCandidate("trigger_name")
			p.addRuleCandidate("trigger_ref")
		case nodes.DropSequence:
			p.addRuleCandidate("sequence_ref")
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
	} else if stmt.ObjectType == nodes.DropIndex {
		// DROP INDEX clauses are read by parseDropIndexClauses (create_index.go).
		if err := p.parseDropIndexClauses(stmt); err != nil {
			return nil, err
		}
		names = stmt.Names.Items
		onTables = stmt.OnTables.Items
	} else {
		// Every object this arm drops names at least one, and a comma promises
		// another: `DROP TABLE` and `DROP TABLE a,` are both the start of nothing
		// (SQL Server answers "Incorrect syntax near the end of the statement").
		name, err := p.parseObjectRef(dropObjectCompletionRule(stmt.ObjectType))
		if err != nil {
			return nil, err
		}
		if name == nil {
			return nil, p.syntaxErrorAtCur()
		}
		names = append(names, name)
		for {
			if _, ok := p.match(','); !ok {
				break
			}
			next, err := p.parseObjectRef(dropObjectCompletionRule(stmt.ObjectType))
			if err != nil {
				return nil, err
			}
			if next == nil {
				return nil, p.syntaxErrorAtCur()
			}
			names = append(names, next)
		}
	}
	stmt.Names = &nodes.List{Items: names}
	if stmt.ObjectType == nodes.DropIndex && len(onTables) > 0 {
		stmt.OnTables = &nodes.List{Items: onTables}
	}

	// DROP TRIGGER ... ON { DATABASE | ALL SERVER } (for DDL/logon triggers).
	// The ON takes one of the two scopes, and DROP TRIGGER a ON is Msg 102
	// (measured).
	if stmt.ObjectType == nodes.DropTrigger {
		if _, ok := p.match(kwON); ok {
			if p.cur.Type == kwDATABASE {
				stmt.OnDatabase = true
				p.advance()
			} else if p.cur.Type == kwALL {
				p.advance()
				if !p.matchIdentCI("SERVER") {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.OnAllServer = true
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}

	// DROP ASSEMBLY ... WITH NO DEPENDENTS, whose three words are all
	// required: DROP ASSEMBLY a WITH NO is Msg 102 (measured).
	if stmt.ObjectType == nodes.DropAssembly && p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect(kwNO); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwDEPENDENTS); err != nil {
			return nil, err
		}
		stmt.NoDependents = true
	}

	// SQL Server takes no CASCADE or RESTRICT on a DROP: DROP SCHEMA a
	// CASCADE and ... RESTRICT are Msg 156 (measured), and ScriptDOM's
	// dropSchemaStatement carries them for another engine.

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

func dropObjectCompletionRule(objectType nodes.DropObjectType) string {
	switch objectType {
	case nodes.DropTable:
		return "table_ref"
	case nodes.DropView, nodes.DropMaterializedView:
		return "view_name"
	case nodes.DropIndex:
		return "index_ref"
	case nodes.DropProcedure:
		return "proc_ref"
	case nodes.DropFunction:
		return "func_name"
	case nodes.DropDatabase:
		return "database_ref"
	case nodes.DropTrigger:
		return "trigger_ref"
	case nodes.DropSequence:
		return "sequence_ref"
	default:
		return "identifier"
	}
}
