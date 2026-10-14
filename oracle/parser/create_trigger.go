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
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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
	if stmt.Events.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
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
		if stmt.Table == nil || stmt.Table.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.isIdentLikeStr("REFERENCING") {
		var parseErr568 error
		stmt.Referencing, parseErr568 = p.parseTriggerReferencingClause()
		if parseErr568 != nil {
			return nil, parseErr568
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
			parseErr569 error
			stmt.When, parseErr569 = p.parseExpr()
			if parseErr569 != nil {
				return nil, parseErr569
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
					// Outer END — consume END, optional name, and the ';'
					// that ends the trigger (PLS-00103 without it).
					p.advance() // consume END
					if p.atEndName() {
						p.advance() // consume trigger name
					}
					if p.cur.Type != ';' {
						return nil, p.syntaxErrorAtCur()
					}
					p.advance() // consume ;
					break
				}
			}
			p.advance()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Body: CALL routine [ ( argument [, argument ]... ) ] | plsql_block
	if p.isIdentLikeStr("CALL") {
		call, err := p.parseTriggerCall()
		if err != nil {
			return nil, err
		}
		stmt.Body = call
	} else if p.cur.Type == kwDECLARE || p.cur.Type == kwBEGIN || p.cur.Type == tokLABELOPEN {
		block, parseErr571 := p.parsePLSQLBlock()
		if parseErr571 != nil {
			return nil, parseErr571
		}
		if err := p.checkBlockCharsetScopes(block); err != nil {
			return nil, err
		}
		stmt.Body = block
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func (p *Parser) parseTriggerReferencingClause() (*nodes.TriggerReferencingClause, error) {
	start := p.pos()
	p.advance() // consume REFERENCING
	clause := &nodes.TriggerReferencingClause{Loc: nodes.Loc{Start: start}}
	for {
		if p.isIdentLikeStr("OLD") {
			p.advance()
			if p.cur.Type == kwAS {
				p.advance()
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			clause.OldAlias = name
		} else if p.isIdentLikeStr("NEW") {
			p.advance()
			if p.cur.Type == kwAS {
				p.advance()
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			clause.NewAlias = name
		} else if p.isIdentLikeStr("PARENT") {
			p.advance()
			if p.cur.Type == kwAS {
				p.advance()
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			clause.ParentAlias = name
		} else {
			break
		}
	}
	clause.Loc.End = p.prev.End
	return clause, nil
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
		if p.cur.Type == kwON || p.cur.Type == kwBEGIN || p.cur.Type == kwDECLARE || p.cur.Type == tokEOF {
			return p.syntaxErrorAtCur()
		}
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

// parseTriggerCall parses the CALL body of a trigger, the current token at
// CALL. The arguments are expressions, :OLD and :NEW references among them,
// in positional or named notation. Oracle 23ai rejects an empty argument
// list (ORA-00936), an INTO clause (ORA-00911), and a ';' after the call
// (PLS-00103), takes schema.package.routine, and rejects a database link
// (ORA-00911 for CALL p@lnk).
//
//	CALL [ schema . ] [ package . ] routine [ ( argument [, argument ]... ) ]
func (p *Parser) parseTriggerCall() (*nodes.CallStmt, error) {
	call := &nodes.CallStmt{Args: &nodes.List{}, Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // consume CALL
	name, pkg, err := p.parseRoutineName()
	if err != nil {
		return nil, err
	}
	// The routine has a name (ORA-03050 for CALL (1)).
	if name == nil || name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	if name.DBLink != "" {
		return nil, p.syntaxErrorAtNode(name)
	}
	call.Name, call.Package = name, pkg
	if p.cur.Type == '(' {
		p.advance() // consume (
		named := false
		for {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == tokASSOC {
				if arg, err = p.parseNamedArgExpr(arg); err != nil {
					return nil, err
				}
				named = true
			} else if named {
				// PLS-00312: a positional parameter association may not
				// follow a named association
				return nil, p.syntaxErrorAtNode(arg)
			}
			call.Args.Items = append(call.Args.Items, arg)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ,
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume )
	}
	call.Loc.End = p.prev.End
	return call, nil
}
