package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseCreateTriggerStmt parses a CREATE TRIGGER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
//
//	CREATE [DEFINER = user] TRIGGER [IF NOT EXISTS] trigger_name
//	    trigger_time trigger_event
//	    ON tbl_name FOR EACH ROW
//	    [trigger_order]
//	    trigger_body
func (p *Parser) parseCreateTriggerStmt() (*nodes.CreateTriggerStmt, error) {
	start := p.pos()
	p.advance() // consume TRIGGER

	// Completion: after CREATE TRIGGER, identifier context.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new trigger name.
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CreateTriggerStmt{Loc: nodes.Loc{Start: start}}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// Trigger name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Completion: after trigger name, offer BEFORE/AFTER.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwBEFORE)
		p.addTokenCandidate(kwAFTER)
		return nil, &ParseError{Message: "collecting"}
	}

	// trigger_time: BEFORE | AFTER
	if p.cur.Type == kwBEFORE {
		stmt.Timing = "BEFORE"
		p.advance()
	} else if p.cur.Type == kwAFTER {
		stmt.Timing = "AFTER"
		p.advance()
	}

	// Completion: after BEFORE/AFTER, offer INSERT/UPDATE/DELETE.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		return nil, &ParseError{Message: "collecting"}
	}

	// trigger_event: INSERT | UPDATE | DELETE
	switch p.cur.Type {
	case kwINSERT:
		stmt.Event = "INSERT"
		p.advance()
	case kwUPDATE:
		stmt.Event = "UPDATE"
		p.advance()
	case kwDELETE:
		stmt.Event = "DELETE"
		p.advance()
	}

	// ON tbl_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	// Completion: after ON, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// FOR EACH ROW
	p.match(kwFOR)
	if p.cur.Type == kwEACH {
		p.advance()
	}
	if p.cur.Type == kwROW {
		p.advance()
	}

	// Optional trigger_order: { FOLLOWS | PRECEDES } other_trigger_name
	if p.cur.Type == kwFOLLOWS || p.cur.Type == kwPRECEDES {
		follows := p.cur.Type == kwFOLLOWS
		p.advance()
		trigName, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Order = &nodes.TriggerOrder{
			Loc:         nodes.Loc{Start: start, End: p.pos()},
			Follows:     follows,
			TriggerName: trigName,
		}
	}

	// Trigger body — consume everything until EOF or ;
	bodyStart := p.pos()
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' && depth == 0 {
			break
		}
		if p.cur.Type == kwBEGIN {
			depth++
		}
		if p.cur.Type == kwEND {
			if depth > 0 {
				depth--
				if depth == 0 {
					p.advance()
					break
				}
			} else {
				break
			}
		}
		p.advance()
	}
	bodyEnd := p.pos()
	if bodyEnd > bodyStart {
		stmt.Body = p.inputText(bodyStart, bodyEnd)
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCreateEventStmt parses a CREATE EVENT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-event.html
//
//	CREATE [DEFINER = user] EVENT [IF NOT EXISTS] event_name
//	    ON SCHEDULE schedule
//	    [ON COMPLETION [NOT] PRESERVE]
//	    [ENABLE | DISABLE | DISABLE ON SLAVE]
//	    [COMMENT 'string']
//	    DO event_body
func (p *Parser) parseCreateEventStmt() (*nodes.CreateEventStmt, error) {
	start := p.pos()
	p.advance() // consume EVENT

	// Completion: after CREATE EVENT, identifier context.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new event name.
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CreateEventStmt{Loc: nodes.Loc{Start: start}}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// Event name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON SCHEDULE
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if p.cur.Type == kwSCHEDULE {
		p.advance()
	}

	// Completion: after ON SCHEDULE, offer AT/EVERY keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwAT)
		return nil, &ParseError{Message: "collecting"}
	}

	sched, err := p.parseEventSchedule()
	if err != nil {
		return nil, err
	}
	stmt.Schedule = sched

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

	// DO event_body
	if p.cur.Type == kwDO {
		p.advance()
	}

	bodyStart := p.pos()
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' && depth == 0 {
			break
		}
		if p.cur.Type == kwBEGIN {
			depth++
		}
		if p.cur.Type == kwEND {
			if depth > 0 {
				depth--
				if depth == 0 {
					p.advance()
					break
				}
			} else {
				break
			}
		}
		p.advance()
	}
	bodyEnd := p.pos()
	if bodyEnd > bodyStart {
		stmt.Body = p.inputText(bodyStart, bodyEnd)
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseEventSchedule parses an event schedule.
//
//	AT timestamp [+ INTERVAL interval] ...
//	EVERY interval [STARTS timestamp] [ENDS timestamp]
func (p *Parser) parseEventSchedule() (*nodes.EventSchedule, error) {
	start := p.pos()
	sched := &nodes.EventSchedule{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwAT {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		sched.At = expr
	} else if p.cur.Type == kwEVERY {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		sched.Every = expr
		// Consume optional interval unit (HOUR, DAY, MINUTE, SECOND, WEEK, MONTH, YEAR, QUARTER)
		// Consume optional interval unit (HOUR, DAY, MINUTE, SECOND, WEEK, MONTH, YEAR, QUARTER)
		if isIntervalUnitToken(p.cur) {
			p.advance()
		}

		// Optional STARTS
		if p.cur.Type == kwSTARTS {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			sched.Starts = expr
		}

		// Optional ENDS
		if p.cur.Type == kwENDS {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			sched.Ends = expr
		}
	}

	sched.Loc.End = p.pos()
	sched.RawText = strings.TrimSpace(p.inputText(start, sched.Loc.End))
	return sched, nil
}

// isIntervalUnitToken returns true if the token represents a MySQL interval unit keyword.
func isIntervalUnitToken(tok Token) bool {
	switch tok.Type {
	case kwYEAR, kwHOUR, kwDAY, kwMINUTE, kwSECOND, kwWEEK, kwMONTH, kwQUARTER,
		kwHOUR_MINUTE, kwHOUR_SECOND, kwHOUR_MICROSECOND,
		kwDAY_HOUR, kwDAY_MINUTE, kwDAY_SECOND, kwDAY_MICROSECOND,
		kwMINUTE_SECOND, kwMINUTE_MICROSECOND, kwSECOND_MICROSECOND:
		return true
	}
	if tok.Type == tokIDENT {
		switch strings.ToUpper(tok.Str) {
		case "HOUR", "DAY", "MINUTE", "SECOND", "WEEK", "MONTH", "QUARTER",
			"HOUR_MINUTE", "HOUR_SECOND", "DAY_HOUR", "DAY_MINUTE", "DAY_SECOND",
			"YEAR_MONTH", "MINUTE_SECOND":
			return true
		}
	}
	return false
}
