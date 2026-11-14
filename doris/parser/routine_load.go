package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseCreateRoutineLoad parses CREATE ROUTINE LOAD as both of Doris's
// parsers take it (DorisParser.g4:381; the CUP parser's
// create_routine_load_stmt, DorisLegacyParser.g4:814):
//
//	CREATE ROUTINE LOAD [db.]job_name [ON table_name]
//	    [WITH (APPEND | DELETE | MERGE)]
//	    [load_property (',' load_property)*]
//	    [PROPERTIES ('k' = 'v', ...)]
//	    FROM type ('k' = 'v', ...)
//	    [COMMENT 'text']
//
// On entry, CREATE has been consumed and cur is ROUTINE.
func (p *Parser) parseCreateRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	stmt := &ast.CreateRoutineLoadStmt{}

	// job_label: ident ('.' ident)?
	label, err := p.parseLoadLabel()
	if err != nil {
		return nil, err
	}
	stmt.Name = &ast.ObjectName{Parts: strings.Split(label, "."), Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}

	if p.cur.Kind == kwON {
		p.advance() // consume ON
		name, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.OnTable = &ast.ObjectName{Parts: []string{name}, Loc: loc}
	}
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		switch p.cur.Kind {
		case kwAPPEND, kwDELETE, kwMERGE:
			stmt.MergeType = strings.ToUpper(p.advance().Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind != kwPROPERTIES && p.cur.Kind != kwFROM {
		start := p.cur.Loc
		for {
			if err := p.parseRoutineLoadProperty(); err != nil {
				return nil, err
			}
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		stmt.LoadProps = p.sourceText(ast.Loc{Start: start.Start, End: p.prev.Loc.End})
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.JobProperties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	srcType, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.DataSourceType = strings.ToUpper(srcType)
	if stmt.DataSourceProperties, err = p.parseParenthesisedProperties(); err != nil {
		return nil, err
	}

	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.advance().Str
	}

	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseRoutineLoadProperty reads one loadProperty (DorisParser.g4:403) as the
// CUP parser's load_property (DorisLegacyParser.g4:831) also takes it:
//
//	COLUMNS TERMINATED BY STRING_LITERAL
//	| COLUMNS '(' import_column_desc (',' import_column_desc)* ')'
//	| PRECEDING FILTER expression | WHERE expression | DELETE ON expression
//	| ORDER BY identifier | partitionSpec
//
// where import_column_desc is name ('=' expression)?, and, after the first,
// may sit in parentheses (the CUP parser's import_column_descs).
func (p *Parser) parseRoutineLoadProperty() error {
	switch p.cur.Kind {
	case kwCOLUMNS:
		if p.peekNext().Kind == kwTERMINATED {
			_, err := p.parseTerminatedBy(kwCOLUMNS)
			return err
		}
		p.advance() // consume COLUMNS
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		for first := true; ; first = false {
			paren := !first && p.cur.Kind == int('(')
			if paren {
				p.advance()
			}
			if _, _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if p.cur.Kind == int('=') {
				p.advance()
				if _, err := p.parseExpr(); err != nil {
					return err
				}
			}
			if paren {
				if _, err := p.expect(int(')')); err != nil {
					return err
				}
			}
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		_, err := p.expect(int(')'))
		return err
	case kwPRECEDING:
		p.advance()
		if _, err := p.expect(kwFILTER); err != nil {
			return err
		}
		_, err := p.parseExpr()
		return err
	case kwWHERE:
		p.advance()
		_, err := p.parseExpr()
		return err
	case kwDELETE:
		p.advance()
		if _, err := p.expect(kwON); err != nil {
			return err
		}
		_, err := p.parseExpr()
		return err
	case kwORDER:
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		_, _, err := p.parseIdentifier()
		return err
	case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
		// partitionSpec, as a LOAD data description reads it.
		var spec ast.LoadDataDesc
		return p.parseLoadPartitionSpec(&spec)
	}
	return p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// ALTER ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseAlterRoutineLoad parses an ALTER ROUTINE LOAD statement.
// On entry, ALTER has been consumed and cur is ROUTINE.
//
// Syntax:
//
//	ALTER ROUTINE LOAD FOR [db.]job_name
//	    [PROPERTIES (...)]
//	    [FROM KAFKA (...)]
func (p *Parser) parseAlterRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	// Consume ROUTINE
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	// Consume LOAD
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	// Consume FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	stmt := &ast.AlterRoutineLoadStmt{}

	// Job name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional PROPERTIES and FROM clauses
	for p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case kwPROPERTIES:
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = props
		case kwFROM:
			p.advance() // consume FROM
			srcType, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.DataSourceType = strings.ToUpper(srcType)
			if p.cur.Kind == int('(') {
				props, err := p.parseParenthesisedProperties()
				if err != nil {
					return nil, err
				}
				stmt.DataSourceProperties = props
			}
		default:
			// No more recognised clauses
			goto done
		}
	}
done:
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// PAUSE ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parsePauseRoutineLoad parses a PAUSE ROUTINE LOAD statement.
// On entry, PAUSE has been consumed and cur is ROUTINE.
//
// Syntax:
//
//	PAUSE ROUTINE LOAD FOR [db.]job_name
//	PAUSE ALL ROUTINE LOAD [FOR db]
func (p *Parser) parsePauseRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	return p.parseRoutineLoadControl(startLoc, "PAUSE")
}

// ---------------------------------------------------------------------------
// RESUME ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseResumeRoutineLoad parses a RESUME ROUTINE LOAD statement.
// On entry, RESUME has been consumed and cur is ROUTINE.
//
// Syntax:
//
//	RESUME ROUTINE LOAD FOR [db.]job_name
//	RESUME ALL ROUTINE LOAD [FOR db]
func (p *Parser) parseResumeRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	return p.parseRoutineLoadControl(startLoc, "RESUME")
}

// parseRoutineLoadControl is the shared implementation for PAUSE/RESUME
// ROUTINE LOAD. kind is "PAUSE" or "RESUME"; on entry PAUSE/RESUME has been
// consumed and cur is ROUTINE.
func (p *Parser) parseRoutineLoadControl(startLoc ast.Loc, kind string) (ast.Node, error) {
	// Consume ROUTINE
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	// Consume LOAD
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	// Consume FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// Job name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	endLoc := ast.NodeLoc(name)

	if kind == "PAUSE" {
		stmt := &ast.PauseRoutineLoadStmt{
			Name: name,
			Loc:  startLoc.Merge(endLoc),
		}
		return stmt, nil
	}
	stmt := &ast.ResumeRoutineLoadStmt{
		Name: name,
		Loc:  startLoc.Merge(endLoc),
	}
	return stmt, nil
}

// parsePauseAllRoutineLoad parses PAUSE ALL ROUTINE LOAD [FOR db].
// On entry, PAUSE ALL has been consumed and cur is ROUTINE.
func (p *Parser) parsePauseAllRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.PauseRoutineLoadStmt{All: true}
	if err := p.consumeRoutineLoad(); err != nil {
		return nil, err
	}
	// Optional FOR db
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.For = dbName
	}
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseResumeAllRoutineLoad parses RESUME ALL ROUTINE LOAD [FOR db].
// On entry, RESUME ALL has been consumed and cur is ROUTINE.
func (p *Parser) parseResumeAllRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.ResumeRoutineLoadStmt{All: true}
	if err := p.consumeRoutineLoad(); err != nil {
		return nil, err
	}
	// Optional FOR db
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.For = dbName
	}
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// consumeRoutineLoad consumes ROUTINE LOAD keywords; cur must be ROUTINE on entry.
func (p *Parser) consumeRoutineLoad() error {
	if _, err := p.expect(kwROUTINE); err != nil {
		return err
	}
	if _, err := p.expect(kwLOAD); err != nil {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// STOP ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseStopRoutineLoad parses a STOP ROUTINE LOAD statement.
// On entry, STOP has been consumed and cur is ROUTINE.
//
// Syntax:
//
//	STOP ROUTINE LOAD FOR [db.]job_name
func (p *Parser) parseStopRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	// Consume ROUTINE
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	// Consume LOAD
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	// Consume FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.StopRoutineLoadStmt{
		Name: name,
		Loc:  startLoc.Merge(ast.NodeLoc(name)),
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SHOW ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseShowRoutineLoad parses SHOW [ALL] ROUTINE LOAD ... statements.
// On entry, SHOW has been consumed and cur is ROUTINE (or ALL for SHOW ALL).
//
// Syntax:
//
//	SHOW [ALL] ROUTINE LOAD [FOR [db.]name | LIKE 'pattern' | FROM db]
//	SHOW ROUTINE LOAD TASK FROM db [WHERE condition]
func (p *Parser) parseShowRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	var all bool

	// Check for SHOW ALL ROUTINE LOAD
	if p.cur.Kind == kwALL {
		all = true
		p.advance() // consume ALL
	}

	// Consume ROUTINE
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	// Consume LOAD
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}

	// Check for SHOW ROUTINE LOAD TASK
	if p.cur.Kind == kwTASK {
		p.advance() // consume TASK
		return p.parseShowRoutineLoadTask(startLoc)
	}

	stmt := &ast.ShowRoutineLoadStmt{All: all}

	// Optional qualifiers: FOR [db.]name | LIKE 'pattern' | FROM db
	switch p.cur.Kind {
	case kwFOR:
		p.advance() // consume FOR
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name

	case kwLIKE:
		p.advance() // consume LIKE
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Like = p.cur.Str
		p.advance()

	case kwFROM:
		p.advance() // consume FROM
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = dbName
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseShowRoutineLoadTask parses SHOW ROUTINE LOAD TASK ... statement.
// On entry, SHOW ROUTINE LOAD TASK has been consumed; cur is at FROM or WHERE or EOF.
func (p *Parser) parseShowRoutineLoadTask(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.ShowRoutineLoadTaskStmt{}

	// Optional (FROM | IN) db
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM / IN
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = dbName
	}

	// Optional WHERE condition
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SYNC (T6.2)
// ---------------------------------------------------------------------------

// parseSyncStmt parses a SYNC statement.
// On entry, SYNC has been consumed.
func (p *Parser) parseSyncStmt(startLoc ast.Loc) (ast.Node, error) {
	return &ast.SyncStmt{Loc: startLoc}, nil
}

// ---------------------------------------------------------------------------
// Shared helper
// ---------------------------------------------------------------------------

// parseParenthesisedProperties parses the parenthesised data-source property
// list of a routine load:
//
//	FROM KAFKA ( "key" = "value" [, "key" = "value"] ... )
//
// cur must be '(' on entry.
//
// The grammar gives the list as a propertyItemList, whose items may take an
// identifier or any constant on either side, but the engine holds a routine
// load to a string on both: (kafka_broker_list = "h:9092"), ("k" = 1),
// ("k" = tp), ("k" = NULL) and ("k" = true) are all rejects, reported by its
// legacy parser (container-verified). The item used to take any single token
// as the value, so ("a" = ,), ("a" = )) and ("a" = ( ) parsed.
func (p *Parser) parseParenthesisedProperties() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	return p.parsePropertyItems(p.parseStringProperty)
}
