package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateTriggerStmt parses a CREATE TRIGGER statement after the TRIGGER keyword.
// The caller has already consumed CREATE [OR REPLACE].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TRIGGER.html
//
//	TRIGGER [schema.]trigger_name
//	  { BEFORE | AFTER | INSTEAD OF }
//	  { INSERT | UPDATE [OF col [, col...]] | DELETE }
//	  [OR { INSERT | UPDATE [OF col [, col...]] | DELETE } ...]
//	  ON [schema.]table_name
//	  [FOR EACH ROW]
//	  [WHEN (condition)]
//	  { plsql_block | CALL routine_name }
func (p *Parser) parseCreateTriggerStmt(start int, orReplace bool) *nodes.CreateTriggerStmt {
	p.advance() // consume TRIGGER

	stmt := &nodes.CreateTriggerStmt{
		OrReplace: orReplace,
		Enable:    true,
		Events:    &nodes.List{},
		Loc:       nodes.Loc{Start: start},
	}

	// Trigger name
	stmt.Name = p.parseObjectName()

	// Timing: BEFORE | AFTER | INSTEAD OF | FOR (compound trigger)
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

	// Events: INSERT | UPDATE [OF col, ...] | DELETE | DDL events (CREATE, ALTER, DROP, etc.)
	// separated by OR
	p.parseTriggerEvent(stmt)
	for p.cur.Type == kwOR {
		p.advance() // consume OR
		p.parseTriggerEvent(stmt)
	}

	// ON [schema.]table_name | ON DATABASE | ON SCHEMA
	if p.cur.Type == kwON {
		p.advance() // consume ON
		if p.cur.Type == kwDATABASE {
			stmt.Table = &nodes.ObjectName{Name: "DATABASE"}
			p.advance()
		} else if p.isIdentLikeStr("SCHEMA") {
			stmt.Table = &nodes.ObjectName{Name: "SCHEMA"}
			p.advance()
		} else {
			stmt.Table = p.parseObjectName()
		}
	}

	// Optional: FOR EACH ROW
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
			p.advance() // consume '('
			stmt.When = p.parseExpr()
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
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Body: CALL routine_name | plsql_block
	if p.isIdentLikeStr("CALL") {
		p.advance() // consume CALL
		// Parse the routine name as an object name, wrap as a simple expression
		callName := p.parseObjectName()
		stmt.Body = &nodes.PLSQLBlock{
			Label: "CALL:" + callName.Name,
			Loc:   callName.Loc,
		}
	} else if p.cur.Type == kwDECLARE || p.cur.Type == kwBEGIN || p.cur.Type == tokLABELOPEN {
		stmt.Body = p.parsePLSQLBlock()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTriggerEvent parses a single trigger event:
//
//	DML: INSERT | UPDATE [OF col, ...] | DELETE
//	DDL: CREATE | ALTER | DROP | TRUNCATE | GRANT | REVOKE | LOGON | LOGOFF | ...
func (p *Parser) parseTriggerEvent(stmt *nodes.CreateTriggerStmt) {
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
				p.parseIdentifier()
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
}
