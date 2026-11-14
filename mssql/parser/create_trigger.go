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
//
// The grammar is triggerStatementBody (TSql160.g:25231), with
// dmlTriggerMidSection (TSql160.g:25292) after a table and
// ddlTriggerMidSection (TSql160.g:25346) after DATABASE or ALL SERVER.
// SQL Server 2022 refuses WITH APPEND ("'APPEND' is not a recognized
// option", Msg 195), which the documentation still lists.
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
	name, err := p.routineParseName("table_ref")
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON clause
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	// Completion: after CREATE TRIGGER tr ON → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Determine target: table/view, DATABASE, or ALL SERVER (triggerObject,
	// TSql160.g:25394).
	switch p.cur.Type {
	case kwDATABASE:
		stmt.OnDatabase = true
		p.advance()
	case kwALL:
		p.advance()
		if err := p.routineExpectWord("SERVER"); err != nil {
			return nil, err
		}
		stmt.OnAllServer = true
	default:
		// DML trigger: ON table_or_view
		tbl, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if tbl == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// At most database.schema.object (Msg 117).
		if tbl.Server != "" {
			return nil, p.newParseError(tbl.Loc.Start, "the object name contains more than the maximum number of prefixes")
		}
		stmt.Table = tbl
	}
	ddl := stmt.Table == nil

	// A database or server trigger has no schema (Msg 1094).
	if ddl && stmt.Name.Schema != "" {
		return nil, p.newParseError(stmt.Name.Loc.Start, "a schema name is not allowed as a prefix to the name of a database or server trigger")
	}

	// Completion: after CREATE TRIGGER tr ON t → FOR/AFTER/INSTEAD OF
	if p.collectMode() {
		p.addTokenCandidate(kwFOR)
		p.addRuleCandidate("AFTER")
		p.addRuleCandidate("INSTEAD OF")
		p.addTokenCandidate(kwWITH)
		return nil, errCollecting
	}

	// Optional WITH clause (trigger options: ENCRYPTION, EXECUTE AS,
	// NATIVE_COMPILATION, SCHEMABINDING).
	var opts routineOptions
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err = p.routineParseOptions(routineTriggerOptions)
		if err != nil {
			return nil, err
		}
		stmt.TriggerOptions = opts.list()
	}

	// { FOR | AFTER | INSTEAD OF }; a DDL trigger takes FOR or AFTER only.
	switch {
	case p.cur.Type == kwFOR:
		stmt.TriggerType = "FOR"
		p.advance()
	case p.cur.Type == kwAFTER:
		stmt.TriggerType = "AFTER"
		p.advance()
	case !ddl && p.cur.Type == kwINSTEAD:
		p.advance()
		if _, err := p.expect(kwOF); err != nil {
			return nil, err
		}
		stmt.TriggerType = "INSTEAD OF"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Completion: after FOR/AFTER/INSTEAD OF → INSERT/UPDATE/DELETE
	if p.collectMode() {
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		return nil, errCollecting
	}

	if ddl {
		events, err := p.triggerParseDDLEvents(stmt.OnAllServer)
		if err != nil {
			return nil, err
		}
		stmt.Events = events
	} else {
		events, err := p.triggerParseDMLActions()
		if err != nil {
			return nil, err
		}
		stmt.Events = events

		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT {
			p.advance() // NOT
			if _, err := p.expect(kwFOR); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwREPLICATION); err != nil {
				return nil, err
			}
			stmt.NotForReplication = true
		}
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// A natively compiled trigger takes neither ENCRYPTION nor NOT FOR
	// REPLICATION (Msg 10794).
	if opts.native && (opts.has["ENCRYPTION"] || stmt.NotForReplication) {
		return nil, p.newParseError(loc, "an option is not supported with natively compiled modules")
	}

	switch {
	case opts.native:
		body, err := p.routineParseAtomicBody()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	case p.cur.Type == kwEXTERNAL:
		// EXTERNAL NAME assembly_name.class_name.method_name (CLR trigger)
		ext, err := p.routineParseMethodSpecifier()
		if err != nil {
			return nil, err
		}
		stmt.ExternalName = ext
	default:
		// Body: ( Semicolon )* statementList, which is not empty.
		body, err := p.routineParseStmtBody(true)
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// triggerParseDMLActions parses the actions of a DML trigger:
// { INSERT | UPDATE | DELETE } [ ,...n ] (dmlTriggerAction, TSql160.g:25429),
// each at most once (Msg 1034).
func (p *Parser) triggerParseDMLActions() (*nodes.List, error) {
	var events []nodes.Node
	seen := map[string]bool{}
	for {
		evtLoc := p.pos()
		var name string
		switch p.cur.Type {
		case kwINSERT:
			name = "INSERT"
		case kwUPDATE:
			name = "UPDATE"
		case kwDELETE:
			name = "DELETE"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if seen[name] {
			return nil, p.newParseError(evtLoc, "duplicate specification of the action "+name+" in the trigger declaration")
		}
		seen[name] = true
		p.advance()
		events = append(events, &nodes.TriggerEvent{Name: name, Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: events}, nil
}

// triggerParseDDLEvents parses the events of a DDL or logon trigger:
// { event_type | event_group | LOGON } [ ,...n ] (ddlTriggerAction,
// TSql160.g:25450). ScriptDOM reads an event as any identifier; SQL Server
// knows the names in sys.trigger_event_types (Msg 1084 "is an invalid event
// type"), takes the server-level ones and LOGON only ON ALL SERVER
// (Msg 1098), refuses a delimited name (Msg 102), and refuses an event named
// twice, by itself or in a group (Msg 1095).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/triggers/ddl-event-groups
func (p *Parser) triggerParseDDLEvents(onAllServer bool) (*nodes.List, error) {
	var events []nodes.Node
	var names []string
	for {
		evtLoc := p.pos()
		if p.cur.Str == "" || !(p.cur.Type == tokIDENT || p.cur.Type >= kwABSENT) || p.isDelimited(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		name := strings.ToUpper(p.cur.Str)
		_, known := triggerEventParents[name]
		switch {
		case known && triggerEventIsDatabaseLevel(name):
		case onAllServer && (known || name == "LOGON"):
		default:
			return nil, p.newParseError(evtLoc, "'"+p.cur.Str+"' is an invalid event type for this trigger")
		}
		for _, prev := range names {
			if triggerEventHolds(prev, name) || triggerEventHolds(name, prev) {
				return nil, p.newParseError(evtLoc, "'"+p.cur.Str+"' has already been specified as an event type")
			}
		}
		names = append(names, name)
		p.advance()
		events = append(events, &nodes.TriggerEvent{Name: name, Loc: nodes.Loc{Start: evtLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: events}, nil
}
