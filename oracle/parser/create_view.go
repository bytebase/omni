package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateViewStmt parses a CREATE [OR REPLACE] [[NO] FORCE] VIEW statement
// or CREATE [OR REPLACE] MATERIALIZED VIEW statement.
// The CREATE keyword has already been consumed. The caller has already parsed
// OR REPLACE if present and passes orReplace.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-VIEW.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-MATERIALIZED-VIEW.html
//
//	CREATE [ OR REPLACE ] [ FORCE | NO FORCE ] VIEW [ schema. ] view_name
//	    [ ( column_alias [, ...] ) ]
//	    AS select_statement
//	    [ WITH CHECK OPTION | WITH READ ONLY ]
//
//	CREATE [ OR REPLACE ] MATERIALIZED VIEW [ schema. ] view_name
//	    [ BUILD IMMEDIATE | BUILD DEFERRED ]
//	    [ REFRESH ... ]
//	    [ ENABLE QUERY REWRITE ]
//	    AS select_statement
func (p *Parser) parseCreateViewStmt(start int, orReplace bool) *nodes.CreateViewStmt {
	stmt := &nodes.CreateViewStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// MATERIALIZED VIEW
	if p.cur.Type == kwMATERIALIZED {
		stmt.Materialized = true
		p.advance()
		if p.cur.Type == kwVIEW {
			p.advance()
		}
	} else {
		// [FORCE | NO FORCE] VIEW
		if p.cur.Type == kwFORCE {
			stmt.Force = true
			p.advance()
		} else if p.cur.Type == kwNOT || (p.isIdentLikeStr("NO") && p.peekNext().Type == kwFORCE) {
			// NO FORCE
			stmt.NoForce = true
			p.advance() // consume NO
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
			}
		}

		if p.cur.Type == kwVIEW {
			p.advance()
		}
	}

	// View name
	stmt.Name = p.parseObjectName()

	// Optional column alias list: ( col1, col2, ... )
	if p.cur.Type == '(' {
		p.advance()
		stmt.Columns = &nodes.List{}
		for {
			name := p.parseIdentifier()
			if name == "" {
				break
			}
			stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Materialized view options (before AS)
	if stmt.Materialized {
		p.parseMaterializedViewOptions(stmt)
	}

	// AS select_statement
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Query = p.parseSelectStmt()
	}

	// WITH CHECK OPTION | WITH READ ONLY
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwCHECK {
			p.advance()
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithCheckOpt = true
		} else if p.cur.Type == kwREAD {
			p.advance()
			if p.cur.Type == kwONLY {
				p.advance()
			}
			stmt.WithReadOnly = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMaterializedViewOptions parses BUILD, REFRESH, and ENABLE QUERY REWRITE options.
func (p *Parser) parseMaterializedViewOptions(stmt *nodes.CreateViewStmt) {
	for {
		if p.isIdentLikeStr("BUILD") {
			p.advance()
			if p.cur.Type == kwIMMEDIATE {
				stmt.BuildMode = "IMMEDIATE"
				p.advance()
			} else if p.cur.Type == kwDEFERRED {
				stmt.BuildMode = "DEFERRED"
				p.advance()
			}
		} else if p.cur.Type == kwREFRESH {
			p.advance()
			// FAST | COMPLETE | FORCE | ON DEMAND | ON COMMIT
			if p.isIdentLikeStr("FAST") {
				stmt.RefreshMode = "FAST"
				p.advance()
			} else if p.isIdentLikeStr("COMPLETE") {
				stmt.RefreshMode = "COMPLETE"
				p.advance()
			} else if p.cur.Type == kwFORCE {
				stmt.RefreshMode = "FORCE"
				p.advance()
			} else if p.cur.Type == kwON {
				p.advance()
				if p.isIdentLikeStr("DEMAND") {
					stmt.RefreshMode = "ON DEMAND"
					p.advance()
				} else if p.cur.Type == kwCOMMIT {
					stmt.RefreshMode = "ON COMMIT"
					p.advance()
				}
			}
		} else if p.cur.Type == kwENABLE {
			p.advance()
			if p.isIdentLikeStr("QUERY") {
				p.advance()
				if p.cur.Type == kwREWRITE {
					p.advance()
				}
				stmt.EnableQuery = true
			}
		} else {
			break
		}
	}
}
