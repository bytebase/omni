package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// CREATE ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseCreateRoutineLoad parses createRoutineLoadStatement (StarRocks.g4:1209):
//
//	CREATE ROUTINE LOAD (db=qualifiedName '.')? name=identifier ON table=qualifiedName
//	    (loadProperties (',' loadProperties)*)?
//	    jobProperties?
//	    FROM source=identifier
//	    dataSourceProperties?
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
	var err error
	if stmt.Name, err = p.parseMultipartIdentifier(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if stmt.OnTable, err = p.parseMultipartIdentifier(); err != nil {
		return nil, err
	}
	if stmt.LoadProps, err = p.parseRoutineLoadProperties(); err != nil {
		return nil, err
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
	if p.cur.Kind == int('(') {
		if stmt.DataSourceProperties, err = p.parseParenthesisedProperties(); err != nil {
			return nil, err
		}
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseRoutineLoadProperties reads an optional
// loadProperties (',' loadProperties)* (StarRocks.g4:1228) and returns its
// source text:
//
//	COLUMNS TERMINATED BY string | ROWS TERMINATED BY string
//	| COLUMNS '(' (qualifiedName | assignment) (',' (qualifiedName | assignment))* ')'
//	| WHERE expression | partitionNames
//	assignment : identifier EQ expressionOrDefault
func (p *Parser) parseRoutineLoadProperties() (string, error) {
	switch p.cur.Kind {
	case kwCOLUMNS, kwROWS, kwWHERE, kwTEMPORARY, kwPARTITION, kwPARTITIONS:
	default:
		return "", nil
	}
	start := p.cur.Loc
	for {
		if err := p.parseRoutineLoadProperty(); err != nil {
			return "", err
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return p.input[start.Start-p.baseOffset : p.prev.Loc.End-p.baseOffset], nil
}

// parseRoutineLoadProperty reads one loadProperties.
func (p *Parser) parseRoutineLoadProperty() error {
	switch p.cur.Kind {
	case kwCOLUMNS, kwROWS:
		kw := p.cur.Kind
		p.advance() // consume COLUMNS / ROWS
		if p.cur.Kind == kwTERMINATED {
			p.advance() // consume TERMINATED
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			if p.cur.Kind != tokString {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			return nil
		}
		if kw == kwROWS {
			return p.syntaxErrorAtCur()
		}
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		for {
			name, err := p.parseMultipartIdentifier()
			if err != nil {
				return err
			}
			if p.cur.Kind == int('=') {
				// assignment: identifier EQ expressionOrDefault
				if len(name.Parts) != 1 {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume '='
				if p.cur.Kind == kwDEFAULT {
					p.advance()
				} else if _, err := p.parseExpr(); err != nil {
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
	case kwWHERE:
		p.advance() // consume WHERE
		_, err := p.parseExpr()
		return err
	case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
		var names ast.AlterTableAction
		return p.parseAlterPartitionNames(&names)
	}
	return p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// ALTER ROUTINE LOAD (T6.2)
// ---------------------------------------------------------------------------

// parseAlterRoutineLoad parses alterRoutineLoadStatement (StarRocks.g4:1217):
//
//	ALTER ROUTINE LOAD FOR (db=qualifiedName '.')? name=identifier
//	    (loadProperties (',' loadProperties)*)?
//	    jobProperties?
//	    (FROM source=identifier dataSourceProperties)?
//
// On entry, ALTER has been consumed and cur is ROUTINE.
func (p *Parser) parseAlterRoutineLoad(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	stmt := &ast.AlterRoutineLoadStmt{}
	var err error
	if stmt.Name, err = p.parseMultipartIdentifier(); err != nil {
		return nil, err
	}
	if stmt.LoadProps, err = p.parseRoutineLoadProperties(); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		srcType, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.DataSourceType = strings.ToUpper(srcType)
		if stmt.DataSourceProperties, err = p.parseParenthesisedProperties(); err != nil {
			return nil, err
		}
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
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

// parseShowRoutineLoad parses showRoutineLoadStatement, the task form, and
// showCreateRoutineLoadStatement (StarRocks.g4:1274-1288):
//
//	SHOW ALL? ROUTINE LOAD (FOR name)? (FROM db)? (WHERE expression)?
//	    (ORDER BY sortItem (',' sortItem)*)? limitElement?
//	SHOW ROUTINE LOAD TASK (FROM db)? WHERE expression
//	SHOW CREATE ROUTINE LOAD name
//
// SHOW (and CREATE, when create) has been consumed; cur is ALL or ROUTINE.
func (p *Parser) parseShowRoutineLoad(startLoc ast.Loc, create bool) (ast.Node, error) {
	all := false
	if !create && p.cur.Kind == kwALL {
		all = true
		p.advance() // consume ALL
	}
	if _, err := p.expect(kwROUTINE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	stmt := &ast.ShowRoutineLoadStmt{All: all, Create: create}
	var err error
	if create {
		if stmt.Name, err = p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
		stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
		return stmt, nil
	}
	if !all && p.cur.Kind == kwTASK {
		p.advance() // consume TASK
		return p.parseShowRoutineLoadTask(startLoc)
	}

	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		if stmt.Name, err = p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		db, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = strings.Join(db.Parts, ".")
	}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		if stmt.Where, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if stmt.OrderBy, err = p.parseOrderByList(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwLIMIT {
		if stmt.Limit, stmt.Offset, err = p.parseLimitClause(); err != nil {
			return nil, err
		}
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
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

// parseParenthesisedProperties parses the parenthesised data-source property
// list of a routine load:
//
//	dataSourceProperties : propertyList
//	propertyList         : '(' property (',' property)* ')'
//	property             : key=string '=' value=string
//
// cur must be '(' on entry. The item used to take a bare key and any single
// token as the value, so ("a" = ,), ("a" = )) and ("a" = ( ) parsed, and so
// did (kafka_broker_list = "h:9092") and ("k" = 1), which the engine rejects
// (container-verified).
func (p *Parser) parseParenthesisedProperties() ([]*ast.Property, error) {
	return p.parsePropertyList()
}
