package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseAlterEventStmt parses an ALTER EVENT statement.
// The caller has already consumed ALTER. The current token is EVENT.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-event.html
//
//	ALTER [DEFINER = user] EVENT event_name
//	    [ON SCHEDULE schedule_clause]
//	    [ON COMPLETION [NOT] PRESERVE]
//	    [RENAME TO new_event_name]
//	    [ENABLE | DISABLE | DISABLE ON SLAVE]
//	    [COMMENT 'string']
//	    [DO event_body]
func (p *Parser) parseAlterEventStmt() (*nodes.AlterEventStmt, error) {
	start := p.pos()
	p.advance() // consume EVENT

	// Completion: after ALTER EVENT, offer event_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("event_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.AlterEventStmt{Loc: nodes.Loc{Start: start}}

	// Event name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional: ON SCHEDULE
	if p.cur.Type == kwON {
		next := p.peekNext()
		if next.Type == kwSCHEDULE {
			p.advance() // ON
			p.advance() // SCHEDULE
			sched, err := p.parseEventSchedule()
			if err != nil {
				return nil, err
			}
			stmt.Schedule = sched
		}
	}

	// Optional: ON COMPLETION [NOT] PRESERVE
	if p.cur.Type == kwON {
		next := p.peekNext()
		if next.Type == kwCOMPLETION {
			p.advance() // ON
			p.advance() // COMPLETION
			if _, ok := p.match(kwNOT); ok {
				stmt.OnCompletion = "NOT PRESERVE"
			} else {
				stmt.OnCompletion = "PRESERVE"
			}
			if p.cur.Type == kwPRESERVE {
				p.advance()
			}
		}
	}

	// Optional: RENAME TO new_event_name
	if p.cur.Type == kwRENAME {
		p.advance()
		p.match(kwTO)
		newName, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.RenameTo = newName
	}

	// Optional: ENABLE | DISABLE [ON SLAVE]
	if p.cur.Type == kwENABLE {
		stmt.Enable = "ENABLE"
		p.advance()
	} else if p.cur.Type == kwDISABLE {
		stmt.Enable = "DISABLE"
		p.advance()
		if p.cur.Type == kwON {
			next := p.peekNext()
			if next.Type == kwSLAVE {
				p.advance()
				p.advance()
				stmt.Enable = "DISABLE ON SLAVE"
			}
		}
	}

	// Optional: COMMENT 'string'
	if p.cur.Type == kwCOMMENT {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.Comment = p.cur.Str
			p.advance()
		}
	}

	// Optional: DO event_body
	if p.cur.Type == kwDO {
		p.advance()
		bodyStart := p.pos()
		body, err := p.parseCompoundStmtOrStmt()
		if err != nil {
			return nil, err
		}
		bodyEnd := p.pos()
		stmt.Body = body
		stmt.BodyText = p.inputText(bodyStart, bodyEnd)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterRoutineStmt parses an ALTER FUNCTION or ALTER PROCEDURE statement.
// The caller has already consumed ALTER. The current token is FUNCTION or PROCEDURE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-function.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-procedure.html
//
//	ALTER {FUNCTION | PROCEDURE} sp_name [characteristic ...]
func (p *Parser) parseAlterRoutineStmt(isProcedure bool) (*nodes.AlterRoutineStmt, error) {
	start := p.pos()
	p.advance() // consume FUNCTION or PROCEDURE

	// Completion: after ALTER FUNCTION/PROCEDURE, offer function_ref/procedure_ref.
	p.checkCursor()
	if p.collectMode() {
		if isProcedure {
			p.addRuleCandidate("procedure_ref")
		} else {
			p.addRuleCandidate("function_ref")
		}
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.AlterRoutineStmt{
		Loc:         nodes.Loc{Start: start},
		IsProcedure: isProcedure,
	}

	// Routine name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// Characteristics
	for {
		ch, ok, err := p.parseRoutineCharacteristic()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Characteristics = append(stmt.Characteristics, ch)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropRoutineStmt parses a DROP FUNCTION or DROP PROCEDURE statement.
// The caller has already consumed DROP. The current token is FUNCTION or PROCEDURE.
//
//	DROP {FUNCTION | PROCEDURE} [IF EXISTS] sp_name
func (p *Parser) parseDropRoutineStmt(start int, isProcedure bool) (*nodes.DropRoutineStmt, error) {
	p.advance() // consume FUNCTION or PROCEDURE

	stmt := &nodes.DropRoutineStmt{
		Loc:         nodes.Loc{Start: start},
		IsProcedure: isProcedure,
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// Completion: after DROP FUNCTION/PROCEDURE [IF EXISTS], offer function_ref/procedure_ref.
	p.checkCursor()
	if p.collectMode() {
		if isProcedure {
			p.addRuleCandidate("procedure_ref")
		} else {
			p.addRuleCandidate("function_ref")
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// Routine name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropTriggerStmt parses a DROP TRIGGER statement.
// The caller has already consumed DROP. The current token is TRIGGER.
//
//	DROP TRIGGER [IF EXISTS] [schema_name.]trigger_name
func (p *Parser) parseDropTriggerStmt() (*nodes.DropTriggerStmt, error) {
	start := p.pos()
	p.advance() // consume TRIGGER

	stmt := &nodes.DropTriggerStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// Completion: after DROP TRIGGER [IF EXISTS], offer trigger_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("trigger_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Trigger name (can be schema.trigger_name)
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropEventStmt parses a DROP EVENT statement.
// The caller has already consumed DROP. The current token is EVENT.
//
//	DROP EVENT [IF EXISTS] event_name
func (p *Parser) parseDropEventStmt() (*nodes.DropEventStmt, error) {
	start := p.pos()
	p.advance() // consume EVENT

	stmt := &nodes.DropEventStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// Completion: after DROP EVENT [IF EXISTS], offer event_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("event_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Event name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
