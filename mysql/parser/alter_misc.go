package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseAlterViewStmt parses an ALTER VIEW statement.
// The caller has already consumed ALTER. The current token may be ALGORITHM, DEFINER, SQL, or VIEW.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-view.html
//
//	ALTER
//	    [ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE}]
//	    [DEFINER = user]
//	    [SQL SECURITY { DEFINER | INVOKER }]
//	    VIEW view_name [(column_list)]
//	    AS select_statement
//	    [WITH [CASCADED | LOCAL] CHECK OPTION]
func (p *Parser) parseAlterViewStmt() (*nodes.AlterViewStmt, error) {
	start := p.pos()
	stmt := &nodes.AlterViewStmt{
		Loc: nodes.Loc{Start: start},
	}

	// ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE}
	if p.cur.Type == kwALGORITHM {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Algorithm = name
	}

	// DEFINER = user
	if p.cur.Type == kwDEFINER {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		definer := name
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			host, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			definer = definer + "@" + host
		}
		stmt.Definer = definer
	}

	// SQL SECURITY { DEFINER | INVOKER }
	if p.cur.Type == kwSQL {
		p.advance()
		if p.cur.Type == kwSECURITY || (p.isIdentToken() && eqFold(p.cur.Str, "security")) {
			p.advance()
			name, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.SqlSecurity = name
		}
	}

	// VIEW keyword
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	// View name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// Optional column list
	if p.cur.Type == '(' {
		cols, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// SELECT statement — capture raw text
	selectStart := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Select = sel
	stmt.SelectText = strings.TrimSpace(p.inputText(selectStart, p.pos()))

	// WITH [CASCADED | LOCAL] CHECK OPTION
	if p.cur.Type == kwWITH {
		p.advance()
		checkOption := "CASCADED" // default
		if p.cur.Type == kwCASCADED {
			p.advance()
		} else if p.cur.Type == kwLOCAL {
			checkOption = "LOCAL"
			p.advance()
		}
		if p.cur.Type == kwCHECK {
			p.advance()
			if eqFold(p.cur.Str, "option") {
				p.advance()
			}
		}
		stmt.CheckOption = checkOption
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

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

	stmt := &nodes.AlterEventStmt{Loc: nodes.Loc{Start: start}}

	// Event name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional: ON SCHEDULE
	if p.cur.Type == kwON {
		next := p.peekNext()
		if next.Type == tokIDENT && eqFold(next.Str, "schedule") {
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
		if next.Type == tokIDENT && eqFold(next.Str, "completion") {
			p.advance() // ON
			p.advance() // COMPLETION
			if _, ok := p.match(kwNOT); ok {
				stmt.OnCompletion = "NOT PRESERVE"
			} else {
				stmt.OnCompletion = "PRESERVE"
			}
			if p.isIdentToken() && eqFold(p.cur.Str, "preserve") {
				p.advance()
			}
		}
	}

	// Optional: RENAME TO new_event_name
	if p.cur.Type == kwRENAME {
		p.advance()
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
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
			if next.Type == kwSLAVE || (next.Type == tokIDENT && eqFold(next.Str, "slave")) {
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
			stmt.Body = p.lexer.input[bodyStart:bodyEnd]
		}
	}

	stmt.Loc.End = p.pos()
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

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropRoutineStmt parses a DROP FUNCTION or DROP PROCEDURE statement.
// The caller has already consumed DROP. The current token is FUNCTION or PROCEDURE.
//
//	DROP {FUNCTION | PROCEDURE} [IF EXISTS] sp_name
func (p *Parser) parseDropRoutineStmt(isProcedure bool) (*nodes.DropRoutineStmt, error) {
	start := p.pos()
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

	// Routine name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	stmt.Loc.End = p.pos()
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

	// Trigger name (can be schema.trigger_name)
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	stmt.Loc.End = p.pos()
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

	// Event name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt, nil
}
