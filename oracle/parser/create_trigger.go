package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateTriggerStmt parses a CREATE TRIGGER statement after the TRIGGER keyword.
// The caller has already consumed CREATE [OR REPLACE].
//
// BNF: oracle/parser/bnf/CREATE-TRIGGER.bnf
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ]
//	    [ EDITIONABLE | NONEDITIONABLE ]
//	    TRIGGER plsql_trigger_source ;
func (p *Parser) parseCreateTriggerStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) (*nodes.CreateTriggerStmt, error) {
	p.advance() // consume TRIGGER

	stmt := &nodes.CreateTriggerStmt{
		OrReplace:      orReplace,
		IfNotExists:    ifNotExists,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Enable:         true,
		Events:         &nodes.List{},
		Loc:            nodes.Loc{Start: start},
	}
	var parseErr564 error

	// Trigger name
	stmt.Name, parseErr564 = p.parseObjectName()
	if parseErr564 !=

		// Timing: BEFORE | AFTER | INSTEAD OF | FOR (compound trigger)
		nil {
		return nil, parseErr564
	}

	switch p.cur.Type {
	case kwBEFORE:
		stmt.Timing = nodes.TRIGGER_BEFORE
		p.advance()
	case kwAFTER:
		stmt.Timing = nodes.TRIGGER_AFTER
		p.advance()
	case kwINSTEAD:
		stmt.Timing = nodes.TRIGGER_INSTEAD_OF
		p.advance() // consume INSTEAD
		if p.cur.Type == kwOF {
			p.advance() // consume OF
		}
	case kwFOR:
		// Compound trigger syntax: FOR dml_event ON table COMPOUND TRIGGER
		p.advance() // consume FOR
	}
	parseErr565 :=

		// Events: INSERT | UPDATE [OF col, ...] | DELETE | DDL events (CREATE, ALTER, DROP, etc.)
		// separated by OR
		p.parseTriggerEvent(stmt)
	if parseErr565 != nil {
		return nil, parseErr565
	}
	for p.cur.Type == kwOR {
		p.advance()
		parseErr566 := // consume OR
			p.parseTriggerEvent(stmt)
		if parseErr566 !=

			// ON [schema.]table_name | ON DATABASE | ON SCHEMA
			nil {
			return nil, parseErr566
		}
	}

	if p.cur.Type == kwON {
		p.advance() // consume ON
		if p.cur.Type == kwDATABASE {
			stmt.Table = &nodes.ObjectName{Name: "DATABASE"}
			p.advance()
		} else if p.isIdentLikeStr("SCHEMA") {
			stmt.Table = &nodes.ObjectName{Name: "SCHEMA"}
			p.advance()
		} else {
			var parseErr567 error
			stmt.Table, parseErr567 = p.parseObjectName()
			if parseErr567 !=

				// Optional: FOR EACH ROW
				nil {
				return nil, parseErr567
			}
		}
	}

	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if p.cur.Type == kwEACH {
			p.advance() // consume EACH
			if p.cur.Type == kwROW {
				p.advance() // consume ROW
				stmt.ForEachRow = true
			}
		}
	}

	// Optional: WHEN (condition)
	if p.cur.Type == kwWHEN {
		p.advance() // consume WHEN
		if p.cur.Type == '(' {
			p.advance()
			var // consume '('
			parseErr568 error
			stmt.When, parseErr568 = p.parseExpr()
			if parseErr568 != nil {
				return nil, parseErr568
			}
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
		}
	}

	// Optional: ENABLE / DISABLE
	if p.cur.Type == kwENABLE {
		stmt.Enable = true
		p.advance()
	} else if p.cur.Type == kwDISABLE {
		stmt.Enable = false
		p.advance()
	}

	// Compound trigger: FOR ... ON tbl COMPOUND TRIGGER ... END trigger_name;
	if p.cur.Type == kwCOMPOUND {
		p.advance() // consume COMPOUND
		if p.cur.Type == kwTRIGGER {
			p.advance() // consume TRIGGER
		}
		stmt.Compound = true
		// Skip compound trigger body to matching END trigger_name;
		// Track BEGIN/END depth to handle nested blocks.
		depth := 0
		for p.cur.Type != tokEOF {
			if p.cur.Type == kwBEGIN {
				depth++
			} else if p.cur.Type == kwEND {
				if depth > 0 {
					depth--
				} else {
					// Outer END — consume END, optional name, done.
					p.advance() // consume END
					if p.isIdentLike() {
						p.advance() // consume trigger name
					}
					break
				}
			}
			p.advance()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Body: CALL routine_name | plsql_block
	if p.isIdentLikeStr("CALL") {
		p.advance() // consume CALL
		// Parse the routine name as an object name, wrap as a simple expression
		callName, parseErr569 := p.parseObjectName()
		if parseErr569 != nil {
			return nil, parseErr569
		}
		stmt.Body = &nodes.PLSQLBlock{
			Label: "CALL:" + callName.Name,
			Loc:   callName.Loc,
		}
	} else if p.cur.Type == kwDECLARE || p.cur.Type == kwBEGIN || p.cur.Type == tokLABELOPEN {
		var parseErr570 error
		stmt.Body, parseErr570 = p.parsePLSQLBlock()
		if parseErr570 != nil {
			return nil, parseErr570
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTriggerEvent parses a single trigger event:
//
//	DML: INSERT | UPDATE [OF col, ...] | DELETE
//	DDL: CREATE | ALTER | DROP | TRUNCATE | GRANT | REVOKE | LOGON | LOGOFF | ...
func (p *Parser) parseTriggerEvent(stmt *nodes.CreateTriggerStmt) error {
	switch p.cur.Type {
	case kwINSERT:
		p.advance()
		stmt.Events.Items = append(stmt.Events.Items, &nodes.Integer{Ival: int64(nodes.TRIGGER_INSERT)})
	case kwUPDATE:
		p.advance()
		stmt.Events.Items = append(stmt.Events.Items, &nodes.Integer{Ival: int64(nodes.TRIGGER_UPDATE)})
		// Optional: OF col [, col...]
		if p.cur.Type == kwOF {
			p.advance() // consume OF
			// Skip column list — not stored in the AST
			for {
				parseDiscard572, parseErr571 := p.parseIdentifier()
				_ = parseDiscard572
				if parseErr571 != nil {
					return parseErr571
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
	case kwDELETE:
		p.advance()
		stmt.Events.Items = append(stmt.Events.Items, &nodes.Integer{Ival: int64(nodes.TRIGGER_DELETE)})
	default:
		// DDL/database events: CREATE, ALTER, DROP, TRUNCATE, GRANT, REVOKE, etc.
		// Store as a string in a ColumnRef node for now.
		if p.isIdentLike() || p.cur.Type == kwCREATE || p.cur.Type == kwALTER ||
			p.cur.Type == kwDROP || p.cur.Type == kwTRUNCATE ||
			p.cur.Type == kwGRANT || p.cur.Type == kwREVOKE {
			evtName := p.cur.Str
			p.advance()
			stmt.Events.Items = append(stmt.Events.Items, &nodes.ColumnRef{Column: evtName})
		}
	}
	return nil
}
