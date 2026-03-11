package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/TRUNCATE-TABLE.html
//
//	TRUNCATE TABLE [schema.]table
//	    [{ PRESERVE | PURGE } MATERIALIZED VIEW LOG]
//	    [{ DROP | REUSE } STORAGE]
//	    [CASCADE]
func (p *Parser) parseTruncateStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume TRUNCATE

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: start},
	}

	// TRUNCATE TABLE or TRUNCATE CLUSTER
	if p.cur.Type == kwTABLE {
		p.advance()
	} else if p.cur.Type == kwCLUSTER {
		stmt.Cluster = true
		p.advance()
	}

	// Parse table/cluster name
	stmt.Table = p.parseObjectName()

	// Parse optional clauses
	for {
		if p.cur.Type == kwPURGE {
			// PURGE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				stmt.PurgeMVLog = true
				p.advance() // consume MATERIALIZED
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		} else if p.cur.Type == kwDROP || p.isIdentLike() && p.cur.Str == "REUSE" {
			// DROP STORAGE or REUSE STORAGE
			p.advance()
			if p.cur.Type == kwSTORAGE {
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "PRESERVE" {
			// PRESERVE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				p.advance()
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAnalyzeStmt parses an ANALYZE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ANALYZE.html
//
//	ANALYZE { TABLE | INDEX } [schema.]name
//	    { COMPUTE STATISTICS | ESTIMATE STATISTICS | DELETE STATISTICS | VALIDATE STRUCTURE }
func (p *Parser) parseAnalyzeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ANALYZE

	stmt := &nodes.AnalyzeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Object type: TABLE or INDEX
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
	default:
		stmt.ObjectType = nodes.OBJECT_TABLE
	}

	// Object name
	stmt.Table = p.parseObjectName()

	// Action: COMPUTE STATISTICS, ESTIMATE STATISTICS, DELETE STATISTICS, VALIDATE STRUCTURE
	if p.isIdentLike() {
		action := p.cur.Str
		p.advance()
		// Second word of the action
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	} else if p.cur.Type == kwDELETE {
		p.advance() // consume DELETE
		action := "DELETE"
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	} else if p.cur.Type == kwVALIDATE {
		p.advance() // consume VALIDATE
		action := "VALIDATE"
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExplainPlanStmt parses an EXPLAIN PLAN statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/EXPLAIN-PLAN.html
//
//	EXPLAIN PLAN
//	    [SET STATEMENT_ID = 'id']
//	    [INTO [schema.]table[@dblink]]
//	    FOR statement
func (p *Parser) parseExplainPlanStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume EXPLAIN

	// Expect PLAN
	if p.cur.Type == kwPLAN {
		p.advance()
	}

	stmt := &nodes.ExplainPlanStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional SET STATEMENT_ID = 'id'
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		// STATEMENT_ID is an identifier
		if p.isIdentLike() && p.cur.Str == "STATEMENT_ID" {
			p.advance() // consume STATEMENT_ID
			if p.cur.Type == '=' {
				p.advance() // consume =
			}
			if p.cur.Type == tokSCONST {
				stmt.StatementID = p.cur.Str
				p.advance()
			}
		}
	}

	// Optional INTO [schema.]table
	if p.cur.Type == kwINTO {
		p.advance()
		stmt.Into = p.parseObjectName()
	}

	// FOR statement
	if p.cur.Type == kwFOR {
		p.advance()
		stmt.Statement = p.parseStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFlashbackTableStmt parses a FLASHBACK TABLE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/FLASHBACK-TABLE.html
//
//	FLASHBACK TABLE [schema.]table TO
//	    { SCN expr | TIMESTAMP expr | BEFORE DROP [RENAME TO name] }
func (p *Parser) parseFlashbackTableStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume FLASHBACK

	// Expect TABLE
	if p.cur.Type == kwTABLE {
		p.advance()
	}

	stmt := &nodes.FlashbackTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// SCN expr | TIMESTAMP expr | BEFORE DROP
	switch p.cur.Type {
	case kwSCN:
		p.advance()
		stmt.ToSCN = p.parseExpr()
	case kwTIMESTAMP:
		p.advance()
		stmt.ToTimestamp = p.parseExpr()
	case kwBEFORE:
		p.advance() // consume BEFORE
		if p.cur.Type == kwDROP {
			p.advance() // consume DROP
			stmt.ToBeforeDrop = true
		}
		// Optional RENAME TO name
		if p.cur.Type == kwRENAME {
			p.advance() // consume RENAME
			if p.cur.Type == kwTO {
				p.advance()
			}
			stmt.Rename = p.parseIdentifier()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePurgeStmt parses a PURGE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/PURGE.html
//
//	PURGE { TABLE name | INDEX name | RECYCLEBIN | DBA_RECYCLEBIN | TABLESPACE name }
func (p *Parser) parsePurgeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume PURGE

	stmt := &nodes.PurgeStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwTABLESPACE:
		stmt.ObjectType = nodes.OBJECT_TABLESPACE
		p.advance()
		stmt.Name = p.parseObjectName()
	default:
		// RECYCLEBIN or DBA_RECYCLEBIN (parsed as identifiers)
		if p.isIdentLike() {
			ident := p.cur.Str
			p.advance()
			stmt.Name = &nodes.ObjectName{
				Name: ident,
				Loc:  nodes.Loc{Start: start, End: p.pos()},
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
