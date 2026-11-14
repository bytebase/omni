package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateMTMV parses #createMTMV (DorisParser.g4:92):
//
//	CREATE MATERIALIZED VIEW (IF NOT EXISTS)? multipartIdentifier
//	    (LEFT_PAREN simpleColumnDefs RIGHT_PAREN)? buildMode?
//	    (REFRESH refreshMethod? refreshTrigger?)?
//	    ((DUPLICATE)? KEY identifierList)?
//	    (COMMENT STRING_LITERAL)?
//	    (PARTITION BY LEFT_PAREN mvPartition RIGHT_PAREN)?
//	    (DISTRIBUTED BY (HASH identifierList | RANDOM) (BUCKETS (INTEGER_VALUE | AUTO))?)?
//	    propertyClause?
//	    AS? query
//
// On entry, CREATE has been consumed and cur is MATERIALIZED.
func (p *Parser) parseCreateMTMV(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	stmt := &ast.CreateMTMVStmt{}
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
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Kind == int('(') {
		if stmt.Columns, err = p.parseViewColumns(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwBUILD {
		// buildMode (:1009): BUILD (IMMEDIATE | DEFERRED)
		p.advance() // consume BUILD
		if p.cur.Kind != kwIMMEDIATE && p.cur.Kind != kwDEFERRED {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.BuildMode = strings.ToUpper(p.advance().Str)
	}
	if p.cur.Kind == kwREFRESH {
		p.advance() // consume REFRESH
		// refreshMethod (:1023): COMPLETE | AUTO
		if p.cur.Kind == kwCOMPLETE || p.cur.Kind == kwAUTO {
			stmt.RefreshMethod = strings.ToUpper(p.advance().Str)
		}
		if p.cur.Kind == kwON {
			if stmt.RefreshTrigger, err = p.parseMTMVRefreshTrigger(); err != nil {
				return nil, err
			}
		}
	}
	if p.cur.Kind == kwDUPLICATE || p.cur.Kind == kwKEY {
		stmt.DuplicateKey = p.cur.Kind == kwDUPLICATE
		if stmt.DuplicateKey {
			p.advance() // consume DUPLICATE
		}
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		if stmt.Keys, err = p.parseLoadIdentifierList(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.advance().Str
	}
	if p.cur.Kind == kwPARTITION {
		if err := p.parseMTMVPartition(stmt); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwDISTRIBUTED {
		if stmt.DistributedBy, err = p.parseDistributedBy(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}

	// AS? query. Doris takes the query without AS only after a column list,
	// BUILD, a refresh method or trigger, KEY or DISTRIBUTED BY; with none
	// of them before it, or only a bare REFRESH, COMMENT, PARTITION BY or
	// PROPERTIES, it refuses CREATE MATERIALIZED VIEW a SELECT ...
	// (testdata/compat.tsv).
	if p.cur.Kind == kwAS {
		p.advance()
	} else if stmt.Columns == nil && stmt.BuildMode == "" && stmt.RefreshMethod == "" &&
		stmt.RefreshTrigger == nil && stmt.Keys == nil && stmt.DistributedBy == nil {
		return nil, p.syntaxErrorAtCur()
	}
	var queryNode ast.Node
	if p.cur.Kind == kwWITH {
		wq, err := p.parseWithSelect()
		if err != nil {
			return nil, err
		}
		if queryNode, err = p.parseQueryTail(wq); err != nil {
			return nil, err
		}
	} else {
		query, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if queryNode, err = p.parseQueryTail(query); err != nil {
			return nil, err
		}
	}
	stmt.Query = queryNode

	stmt.Loc = startLoc.Merge(ast.NodeLoc(stmt.Query))
	return stmt, nil
}

// parseMTMVPartition reads PARTITION BY LEFT_PAREN mvPartition RIGHT_PAREN
// (DorisParser.g4:97), where mvPartition is an identifier or a function
// call (:1027).
func (p *Parser) parseMTMVPartition(stmt *ast.CreateMTMVStmt) error {
	start := p.advance().Loc // consume PARTITION
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	exprStart := p.cur
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	switch e := expr.(type) {
	case *ast.FuncCallExpr:
		stmt.PartitionExpr = e
		stmt.PartitionBy = &ast.PartitionDesc{FuncExprs: []string{p.sourceText(e.Loc)}}
	case *ast.ColumnRef:
		if len(e.Name.Parts) != 1 {
			return syntaxErrorAt(exprStart)
		}
		stmt.PartitionBy = &ast.PartitionDesc{Columns: []string{e.Name.Parts[0]}}
	default:
		return syntaxErrorAt(exprStart)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return err
	}
	if stmt.PartitionBy != nil {
		stmt.PartitionBy.Loc = ast.Loc{Start: start.Start, End: closeTok.Loc.End}
	}
	return nil
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
		// refreshSchedule (DorisParser.g4:1019): EVERY INTEGER_VALUE identifier
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		count := p.advance().Str
		unit, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		trigger.Interval = count + " " + strings.ToUpper(unit)
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
		// REFRESH (refreshMethod | refreshTrigger | refreshMethod
		// refreshTrigger) (DorisParser.g4:104), refreshMethod COMPLETE |
		// AUTO (DorisParser.g4:1023).
		p.advance() // consume REFRESH
		if p.cur.Kind == kwCOMPLETE || p.cur.Kind == kwAUTO {
			stmt.RefreshMethod = strings.ToUpper(p.advance().Str)
		}
		if p.cur.Kind == kwON {
			trigger, err := p.parseAlterMTMVRefreshTrigger()
			if err != nil {
				return nil, err
			}
			stmt.RefreshTrigger = trigger
		} else if stmt.RefreshMethod == "" {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Loc = startLoc.Merge(p.prev.Loc)

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
//	REFRESH MATERIALIZED VIEW mv_name {COMPLETE | AUTO | PARTITIONS(p1, p2, ...)}
//
// The refresh mode is required (DorisParser.g4:102 #refreshMTMV).
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

	// (partitionSpec | COMPLETE | AUTO), required (DorisParser.g4:102 #refreshMTMV).
	switch p.cur.Kind {
	case kwCOMPLETE:
		stmt.Mode = "COMPLETE"
		stmt.Loc.End = p.advance().Loc.End
	case kwAUTO:
		stmt.Mode = "AUTO"
		stmt.Loc.End = p.advance().Loc.End
	case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
		if stmt.Partitions, _, err = p.parsePartitionSpec(); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
	default:
		return nil, p.syntaxErrorAtCur()
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

// parseMTMVPropertyList parses ALTER MATERIALIZED VIEW ... SET's list,
// LEFT_PAREN propertyItemList RIGHT_PAREN (DorisParser.g4:106), whose items
// take an identifier or a constant on either side: SET (a = a) parses.
// cur must be '(' on entry.
func (p *Parser) parseMTMVPropertyList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	return p.parsePropertyItems(p.parsePropertyItem)
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

// parseAlterMTMVRefreshTrigger reads refreshTrigger (DorisParser.g4:1013):
//
//	ON MANUAL | ON COMMIT | ON SCHEDULE EVERY INTEGER_VALUE identifier [STARTS STRING_LITERAL]
//
// with refreshSchedule (DorisParser.g4:1019). cur is ON.
func (p *Parser) parseAlterMTMVRefreshTrigger() (*ast.MTMVRefreshTrigger, error) {
	trigger := &ast.MTMVRefreshTrigger{Loc: p.cur.Loc}
	p.advance() // consume ON
	switch p.cur.Kind {
	case kwMANUAL:
		trigger.OnManual = true
		p.advance()
	case kwCOMMIT:
		trigger.OnCommit = true
		p.advance()
	case kwSCHEDULE:
		trigger.OnSchedule = true
		p.advance()
		if _, err := p.expect(kwEVERY); err != nil {
			return nil, err
		}
		n, err := p.alterExpectInt()
		if err != nil {
			return nil, err
		}
		unit, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		trigger.Interval = n.Str + " " + strings.ToUpper(unit)
		if p.cur.Kind == kwSTARTS {
			p.advance()
			tok, err := p.alterExpectString()
			if err != nil {
				return nil, err
			}
			trigger.StartsAt = tok.Str
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	trigger.Loc.End = p.prev.Loc.End
	return trigger, nil
}
