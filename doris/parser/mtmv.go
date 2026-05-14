package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateMTMV parses a CREATE MATERIALIZED VIEW statement.
// On entry, CREATE has been consumed and cur is MATERIALIZED.
//
// Syntax:
//
//	CREATE MATERIALIZED VIEW [IF NOT EXISTS] mv_name
//	    [BUILD IMMEDIATE | DEFERRED]
//	    [REFRESH COMPLETE | AUTO | NEVER | INCREMENTAL]
//	    [(col_name [COMMENT 'text'], ...)]
//	    [COMMENT 'text']
//	    [PARTITION BY ...]
//	    [DISTRIBUTED BY ...]
//	    [ON SCHEDULE EVERY interval [STARTS ts] | ON COMMIT | ON MANUAL]
//	    [PROPERTIES (...)]
//	    AS query
func (p *Parser) parseCreateMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.CreateMTMVStmt{}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parse optional clauses in a tolerant loop until we see AS or EOF.
	// Clauses can appear in any order per Doris docs.
	for p.cur.Kind != kwAS && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case kwBUILD:
			// BUILD IMMEDIATE | DEFERRED
			p.advance() // consume BUILD
			switch p.cur.Kind {
			case kwIMMEDIATE:
				stmt.BuildMode = "IMMEDIATE"
				p.advance()
			case kwDEFERRED:
				stmt.BuildMode = "DEFERRED"
				p.advance()
			default:
				// tolerate unknown token after BUILD
				stmt.BuildMode = strings.ToUpper(p.cur.Str)
				p.advance()
			}

		case kwREFRESH:
			// REFRESH COMPLETE | AUTO | NEVER | INCREMENTAL
			p.advance() // consume REFRESH
			switch p.cur.Kind {
			case kwCOMPLETE:
				stmt.RefreshMethod = "COMPLETE"
				p.advance()
			case kwAUTO:
				stmt.RefreshMethod = "AUTO"
				p.advance()
			case kwNEVER:
				stmt.RefreshMethod = "NEVER"
				p.advance()
			case kwINCREMENTAL:
				stmt.RefreshMethod = "INCREMENTAL"
				p.advance()
			default:
				// tolerate unknown
				stmt.RefreshMethod = strings.ToUpper(p.cur.Str)
				p.advance()
			}

		case kwON:
			// ON SCHEDULE EVERY interval [STARTS ts]
			// ON COMMIT
			// ON MANUAL
			trigger, err := p.parseMTMVRefreshTrigger()
			if err != nil {
				return nil, err
			}
			stmt.RefreshTrigger = trigger

		case int('('):
			// Optional column list: (col_name [COMMENT 'text'], ...)
			cols, err := p.parseViewColumns()
			if err != nil {
				return nil, err
			}
			stmt.Columns = cols

		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Comment = p.cur.Str
			p.advance()

		case kwPARTITION:
			partDesc, err := p.parsePartitionBy()
			if err != nil {
				return nil, err
			}
			stmt.PartitionBy = partDesc

		case kwAUTO:
			// AUTO PARTITION BY ...
			partDesc, err := p.parsePartitionBy()
			if err != nil {
				return nil, err
			}
			stmt.PartitionBy = partDesc

		case kwDISTRIBUTED:
			distDesc, err := p.parseDistributedBy()
			if err != nil {
				return nil, err
			}
			stmt.DistributedBy = distDesc

		case kwPROPERTIES:
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = props

		default:
			// Tolerate unknown tokens — advance to avoid infinite loop.
			p.advance()
		}
	}

	// AS keyword is required
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Query — SELECT or WITH ... SELECT or set operation
	var queryNode ast.Node
	if p.cur.Kind == kwWITH {
		wq, err := p.parseWithSelect()
		if err != nil {
			return nil, err
		}
		queryNode = wq
	} else {
		query, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		queryNode, err = p.parseSetOpTail(query)
		if err != nil {
			return nil, err
		}
	}
	stmt.Query = queryNode

	stmt.Loc = startLoc.Merge(ast.NodeLoc(stmt.Query))
	return stmt, nil
}

// parseMTMVRefreshTrigger parses the ON clause of CREATE MATERIALIZED VIEW:
//
//	ON SCHEDULE EVERY interval [STARTS timestamp]
//	ON COMMIT
//	ON MANUAL
//
// cur is ON on entry; ON is consumed here.
func (p *Parser) parseMTMVRefreshTrigger() (*ast.MTMVRefreshTrigger, error) {
	startLoc := p.cur.Loc
	p.advance() // consume ON

	trigger := &ast.MTMVRefreshTrigger{Loc: startLoc}

	switch p.cur.Kind {
	case kwSCHEDULE:
		trigger.OnSchedule = true
		p.advance() // consume SCHEDULE

		// EVERY interval_expr
		if _, err := p.expect(kwEVERY); err != nil {
			return nil, err
		}

		// Collect interval: integer + unit keyword (e.g., "1 DAY")
		// The interval is stored as raw text.
		intervalStart := p.cur.Loc.Start
		var parts []string
		// First token should be an integer
		if p.cur.Kind == tokInt {
			parts = append(parts, p.cur.Str)
			p.advance()
		}
		// Unit keyword (DAY, HOUR, MINUTE, SECOND, WEEK, MONTH, YEAR, etc.)
		if p.cur.Kind >= 700 || p.cur.Kind == tokIdent {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		}
		if len(parts) > 0 {
			trigger.Interval = strings.Join(parts, " ")
		}
		_ = intervalStart

		// Optional STARTS timestamp_string
		if p.cur.Kind == kwSTARTS {
			p.advance() // consume STARTS
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			trigger.StartsAt = p.cur.Str
			trigger.Loc.End = p.cur.Loc.End
			p.advance()
		}

	case kwCOMMIT:
		trigger.OnCommit = true
		trigger.Loc.End = p.cur.Loc.End
		p.advance()

	case kwMANUAL:
		trigger.OnManual = true
		trigger.Loc.End = p.cur.Loc.End
		p.advance()

	default:
		// Tolerate unknown ON clause variant.
		trigger.Loc.End = p.cur.Loc.End
		p.advance()
	}

	return trigger, nil
}

// parseAlterMTMV parses an ALTER MATERIALIZED VIEW statement.
// On entry, ALTER has been consumed and cur is MATERIALIZED.
//
// Syntax:
//
//	ALTER MATERIALIZED VIEW mv_name
//	    { RENAME new_name
//	    | REFRESH new_method
//	    | REPLACE WITH MATERIALIZED VIEW other_mv
//	    | SET PROPERTIES (...) }
func (p *Parser) parseAlterMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.AlterMTMVStmt{}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action
	switch p.cur.Kind {
	case kwRENAME:
		p.advance() // consume RENAME
		newName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = newName
		stmt.Loc = startLoc.Merge(ast.NodeLoc(newName))

	case kwREFRESH:
		p.advance() // consume REFRESH
		switch p.cur.Kind {
		case kwCOMPLETE:
			stmt.RefreshMethod = "COMPLETE"
		case kwAUTO:
			stmt.RefreshMethod = "AUTO"
		case kwNEVER:
			stmt.RefreshMethod = "NEVER"
		case kwINCREMENTAL:
			stmt.RefreshMethod = "INCREMENTAL"
		default:
			stmt.RefreshMethod = strings.ToUpper(p.cur.Str)
		}
		stmt.Loc = startLoc.Merge(p.cur.Loc)
		p.advance()

	case kwREPLACE:
		// REPLACE WITH MATERIALIZED VIEW other_mv
		p.advance() // consume REPLACE
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwMATERIALIZED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Replace = true
		stmt.ReplaceTarget = target
		stmt.Loc = startLoc.Merge(ast.NodeLoc(target))

	case kwSET:
		p.advance() // consume SET
		// SET ("key"="val", ...) — parse the parenthesised key=value list directly.
		props, err := p.parseMTMVPropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		stmt.Loc = startLoc.Merge(p.prev.Loc)

	default:
		stmt.Loc = startLoc.Merge(p.cur.Loc)
	}

	return stmt, nil
}

// parseDropMTMV parses a DROP MATERIALIZED VIEW statement.
// On entry, DROP has been consumed and cur is MATERIALIZED.
//
// Syntax:
//
//	DROP MATERIALIZED VIEW [IF EXISTS] mv_name [ON base_table]
func (p *Parser) parseDropMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.DropMTMVStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))

	// Optional ON base_table (for sync MVs)
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		baseName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.OnBase = baseName
		stmt.Loc = startLoc.Merge(ast.NodeLoc(baseName))
	}

	return stmt, nil
}

// parseRefreshMTMV parses a REFRESH MATERIALIZED VIEW statement.
// On entry, REFRESH has been consumed and cur is MATERIALIZED.
//
// Syntax:
//
//	REFRESH MATERIALIZED VIEW mv_name [COMPLETE | AUTO | PARTITIONS(p1, p2, ...)]
func (p *Parser) parseRefreshMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.RefreshMTMVStmt{}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))

	// Optional mode or PARTITIONS
	switch p.cur.Kind {
	case kwCOMPLETE:
		stmt.Mode = "COMPLETE"
		stmt.Loc.End = p.cur.Loc.End
		p.advance()
	case kwAUTO:
		stmt.Mode = "AUTO"
		stmt.Loc.End = p.cur.Loc.End
		p.advance()
	case kwPARTITIONS:
		p.advance() // consume PARTITIONS
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		// Parse comma-separated partition names
		for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
			partName, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Partitions = append(stmt.Partitions, partName)
			if p.cur.Kind == int(',') {
				p.advance()
			}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = closeTok.Loc.End
	}

	return stmt, nil
}

// parsePauseMTMVJob parses a PAUSE MATERIALIZED VIEW JOB ON mv_name statement.
// On entry, PAUSE has been consumed and cur is MATERIALIZED.
func (p *Parser) parsePauseMTMVJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}
	// Consume ON
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	stmt := &ast.PauseMTMVJobStmt{}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parseResumeMTMVJob parses a RESUME MATERIALIZED VIEW JOB ON mv_name statement.
// On entry, RESUME has been consumed and cur is MATERIALIZED.
func (p *Parser) parseResumeMTMVJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}
	// Consume ON
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	stmt := &ast.ResumeMTMVJobStmt{}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parseMTMVPropertyList parses a parenthesised key=value list:
//
//	( "key" = "val" [, "key" = "val"] ... )
//
// cur must be '(' on entry. Used for ALTER MATERIALIZED VIEW ... SET (...).
func (p *Parser) parseMTMVPropertyList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var props []*ast.Property

	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		startLoc := p.cur.Loc

		// Key — string literal or bare identifier
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

		// Value — string literal
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		endLoc := p.cur.Loc
		p.advance()

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

// parseCancelMTMVTask parses a CANCEL MATERIALIZED VIEW TASK task_id ON mv_name statement.
// On entry, CANCEL has been consumed and cur is MATERIALIZED.
func (p *Parser) parseCancelMTMVTask(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	// Consume TASK
	if _, err := p.expect(kwTASK); err != nil {
		return nil, err
	}

	stmt := &ast.CancelMTMVTaskStmt{}

	// task_id — integer
	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.TaskID = p.cur.Ival
	p.advance()

	// ON mv_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}
