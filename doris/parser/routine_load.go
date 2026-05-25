package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseCreateRoutineLoad parses a CREATE ROUTINE LOAD statement.
// On entry, CREATE has been consumed and cur is ROUTINE.
//
// Syntax:
//
//	CREATE ROUTINE LOAD [db.]job_name ON table_name
//	    [load_properties]
//	    [PROPERTIES (...)]
//	    FROM {KAFKA | S3 | HDFS} (data_source_properties)
//	    [COMMENT 'text']
func (p *Parser) parseCreateRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	// Consume ROUTINE
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	// Consume LOAD
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}

	stmt := &ast.CreateRoutineLoadStmt{}

	// Job name: [db.]job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON table_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	tableName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.OnTable = tableName

	// Parse optional clauses until FROM (which introduces the data source).
	// load_properties can include COLUMNS, WHERE, PARTITION, etc. We collect
	// them as raw text (best-effort) since their syntax is complex.
	var loadPropParts []string
	for p.cur.Kind != kwFROM && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case kwPROPERTIES:
			// PROPERTIES (...) — job properties
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.JobProperties = props
		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Kind == tokString {
				stmt.Comment = p.cur.Str
				p.advance()
			}
		default:
			// Capture as raw load properties text (COLUMNS, WHERE, SET, etc.)
			loadPropParts = append(loadPropParts, p.cur.Str)
			p.advance()
		}
	}
	if len(loadPropParts) > 0 {
		stmt.LoadProps = strings.Join(loadPropParts, " ")
	}

	// FROM data_source_type (data_source_properties)
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM

		// Data source type: KAFKA, S3, HDFS (treated as non-reserved identifiers)
		srcType, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.DataSourceType = strings.ToUpper(srcType)

		// (data_source_properties)
		if p.cur.Kind == int('(') {
			props, err := p.parseParenthesisedProperties()
			if err != nil {
				return nil, err
			}
			stmt.DataSourceProperties = props
		}
	}

	// Optional trailing COMMENT
	if p.cur.Kind == kwCOMMENT && stmt.Comment == "" {
		p.advance() // consume COMMENT
		if p.cur.Kind == tokString {
			stmt.Comment = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
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

	// Optional FROM db
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
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

// parseParenthesisedProperties parses a parenthesised key=value list where
// values may be string literals or plain identifiers/numbers. This is used for
// the data source properties in ROUTINE LOAD:
//
//	( "key" = "value" [, "key" = "value"] ... )
//
// cur must be '(' on entry.
func (p *Parser) parseParenthesisedProperties() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var props []*ast.Property
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		startLoc := p.cur.Loc

		// Key — string literal or identifier/keyword
		var key string
		if p.cur.Kind == tokString {
			key = p.cur.Str
			p.advance()
		} else {
			var err error
			key, _, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
		}

		// '='
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}

		// Value — string literal, integer, or identifier
		var val string
		endLoc := p.cur.Loc
		switch p.cur.Kind {
		case tokString:
			val = p.cur.Str
			p.advance()
		case tokInt:
			val = p.cur.Str
			p.advance()
		default:
			// Treat keyword or identifier as raw text
			val = p.cur.Str
			p.advance()
		}

		props = append(props, &ast.Property{
			Key:   key,
			Value: val,
			Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
		})

		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return props, nil
}
