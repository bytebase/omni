// Package parser - create_trigger.go implements T-SQL CREATE TRIGGER statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateTriggerStmt parses a CREATE [OR ALTER] TRIGGER statement.
//
// BNF: mssql/parser/bnf/create-trigger-transact-sql.bnf
//
// DML trigger:
//
//	CREATE [ OR ALTER ] TRIGGER [ schema_name . ] trigger_name
//	ON { table | view }
//	[ WITH <dml_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER | INSTEAD OF }
//	{ [ INSERT ] [ , ] [ UPDATE ] [ , ] [ DELETE ] }
//	[ WITH APPEND ]
//	[ NOT FOR REPLICATION ]
//	AS { sql_statement [ ; ] [ , ...n ] | EXTERNAL NAME <method_specifier> [ ; ] }
//
//	<dml_trigger_option> ::=
//	    [ ENCRYPTION ]
//	    [ EXECUTE AS Clause ]
//
//	<method_specifier> ::=
//	    assembly_name.class_name.method_name
//
// DML trigger on memory-optimized tables:
//
//	CREATE [ OR ALTER ] TRIGGER [ schema_name . ] trigger_name
//	ON { table }
//	[ WITH <dml_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER }
//	{ [ INSERT ] [ , ] [ UPDATE ] [ , ] [ DELETE ] }
//	AS { sql_statement [ ; ] [ , ...n ] }
//
//	<dml_trigger_option> ::=
//	    [ NATIVE_COMPILATION ]
//	    [ SCHEMABINDING ]
//	    [ EXECUTE AS Clause ]
//
// DDL trigger:
//
//	CREATE [ OR ALTER ] TRIGGER trigger_name
//	ON { ALL SERVER | DATABASE }
//	[ WITH <ddl_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER } { event_type | event_group } [ , ...n ]
//	AS { sql_statement [ ; ] [ , ...n ] | EXTERNAL NAME <method_specifier> [ ; ] }
//
//	<ddl_trigger_option> ::=
//	    [ ENCRYPTION ]
//	    [ EXECUTE AS Clause ]
//
// Logon trigger:
//
//	CREATE [ OR ALTER ] TRIGGER trigger_name
//	ON ALL SERVER
//	[ WITH <logon_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER } LOGON
//	AS { sql_statement [ ; ] [ , ...n ] | EXTERNAL NAME <method_specifier> [ ; ] }
//
//	<logon_trigger_option> ::=
//	    [ ENCRYPTION ]
//	    [ EXECUTE AS Clause ]
func (p *Parser) parseCreateTriggerStmt(orAlter bool) (*nodes.CreateTriggerStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER TRIGGER → identifier or trigger name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("trigger_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateTriggerStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Trigger name (possibly schema-qualified)
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON clause
	if _, ok := p.match(kwON); !ok {
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// Completion: after CREATE TRIGGER tr ON → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Determine target: table/view, DATABASE, or ALL SERVER
	if p.cur.Type == kwDATABASE {
		stmt.OnDatabase = true
		p.advance()
	} else if p.cur.Type == kwALL {
		p.advance()
		// Expect SERVER
		if p.matchIdentCI("SERVER") {
			stmt.OnAllServer = true
		}
	} else {
		// DML trigger: ON table_or_view
		tbl, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Table = tbl
	}

	// Completion: after CREATE TRIGGER tr ON t → FOR/AFTER/INSTEAD OF
	if p.collectMode() {
		p.addTokenCandidate(kwFOR)
		p.addRuleCandidate("AFTER")
		p.addRuleCandidate("INSTEAD OF")
		p.addTokenCandidate(kwWITH)
		return nil, errCollecting
	}

	// Optional WITH clause (trigger options: ENCRYPTION, EXECUTE AS, NATIVE_COMPILATION, SCHEMABINDING)
	//
	//  <dml_trigger_option> ::=
	//      [ ENCRYPTION ]
	//      [ EXECUTE AS Clause ]
	//      [ NATIVE_COMPILATION ]
	//      [ SCHEMABINDING ]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		// Distinguish WITH ENCRYPTION/EXECUTE AS from WITH APPEND
		if next.Type != kwAPPEND {
			p.advance() // consume WITH
			stmt.TriggerOptions = p.parseTriggerWithOptions()
		}
	}

	// Handle case where matchIdentCI in parseTriggerWithOptions already consumed AFTER/INSTEAD
	// by checking if TriggerType was set during option parsing
	// { FOR | AFTER | INSTEAD OF }
	if p.cur.Type == kwFOR {
		stmt.TriggerType = "FOR"
		p.advance()
	} else if p.matchIdentCI("AFTER") {
		stmt.TriggerType = "AFTER"
	} else if p.matchIdentCI("INSTEAD") {
		// INSTEAD OF
		p.match(kwOF)
		stmt.TriggerType = "INSTEAD OF"
	}

	// Completion: after FOR/AFTER/INSTEAD OF → INSERT/UPDATE/DELETE
	if p.collectMode() {
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		return nil, errCollecting
	}

	// Parse event list
	var events []nodes.Node
	for {
		evtLoc := p.pos()
		if p.cur.Type == kwINSERT {
			events = append(events, &nodes.TriggerEvent{Name: "INSERT", Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
			p.advance()
		} else if p.cur.Type == kwUPDATE {
			events = append(events, &nodes.TriggerEvent{Name: "UPDATE", Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
			p.advance()
		} else if p.cur.Type == kwDELETE {
			events = append(events, &nodes.TriggerEvent{Name: "DELETE", Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
			p.advance()
		} else if p.matchIdentCI("LOGON") {
			events = append(events, &nodes.TriggerEvent{Name: "LOGON", Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
		} else if p.isIdentLike() {
			// DDL event type or event group (e.g., CREATE_TABLE, DDL_TABLE_EVENTS)
			events = append(events, &nodes.TriggerEvent{Name: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
			p.advance()
		} else {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(events) > 0 {
		stmt.Events = &nodes.List{Items: events}
	}

	// [ WITH APPEND ] (for DML triggers)
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwAPPEND {
			p.advance() // WITH
			p.advance() // APPEND
			stmt.WithAppend = true
		}
	}

	// [ NOT FOR REPLICATION ]
	if p.cur.Type == kwNOT {
		next := p.peekNext()
		if next.Type == kwFOR {
			p.advance() // NOT
			p.advance() // FOR
			if p.matchIdentCI("REPLICATION") {
				stmt.NotForReplication = true
			}
		}
	}

	// AS
	p.match(kwAS)

	// EXTERNAL NAME assembly_name.class_name.method_name (CLR trigger)
	if p.cur.Type == kwEXTERNAL {
		p.advance() // consume EXTERNAL
		if p.cur.Type == kwNAME {
			p.advance() // consume NAME
		}
		var parts []string
		for {
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				parts = append(parts, p.cur.Str)
				p.advance()
			} else {
				break
			}
			if p.cur.Type != '.' {
				break
			}
			p.advance() // consume '.'
		}
		stmt.ExternalName = joinDots(parts)
	} else {
		// Body: parse statement(s)
		body, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseTriggerWithOptions parses comma-separated trigger options after WITH.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-trigger-transact-sql
//
//	<dml_trigger_option> ::=
//	    ENCRYPTION
//	  | EXECUTE AS { CALLER | SELF | OWNER | 'user_name' }
//	  | NATIVE_COMPILATION
//	  | SCHEMABINDING
func (p *Parser) parseTriggerWithOptions() *nodes.List {
	var items []nodes.Node
	for {
		if p.cur.Type == kwFOR || p.cur.Type == tokEOF {
			break
		}
		// Check for AFTER / INSTEAD which mark end of options
		if p.cur.Type == kwAFTER || p.cur.Type == kwINSTEAD {
			break
		}

		optLoc := p.pos()
		if p.cur.Type == kwENCRYPTION {
			items = append(items, &nodes.TriggerOption{Name: "ENCRYPTION", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		} else if p.cur.Type == kwEXEC || p.cur.Type == kwEXECUTE {
			p.advance() // consume EXECUTE/EXEC
			if p.cur.Type == kwAS {
				p.advance() // consume AS
			}
			// CALLER | SELF | OWNER | 'user_name'
			var asVal string
			if p.cur.Type == tokSCONST {
				asVal = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				asVal = strings.ToUpper(p.cur.Str)
				p.advance()
			}
			items = append(items, &nodes.TriggerOption{Name: "EXECUTE AS", Value: asVal, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		} else if p.cur.Type == kwNATIVE_COMPILATION {
			items = append(items, &nodes.TriggerOption{Name: "NATIVE_COMPILATION", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		} else if p.cur.Type == kwSCHEMABINDING {
			items = append(items, &nodes.TriggerOption{Name: "SCHEMABINDING", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		} else {
			// Unknown option — break to avoid infinite loop
			break
		}

		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}
