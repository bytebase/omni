package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateViewStmt parses a CREATE [OR REPLACE] VIEW or CREATE MATERIALIZED VIEW statement.
// The CREATE keyword has already been consumed. orReplace is set if OR REPLACE was parsed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-VIEW.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-MATERIALIZED-VIEW.html
//
//	CREATE [ OR REPLACE ]
//	    [ { FORCE | NOFORCE } ]
//	    [ EDITIONING ]
//	    [ { EDITIONABLE | NONEDITIONABLE } ]
//	    VIEW [ schema. ] view_name
//	    [ IF NOT EXISTS ]
//	    [ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
//	    [ ( column_alias [,...] ) ]
//	    [ DEFAULT COLLATION collation_name ]
//	    [ BEQUEATH { CURRENT_USER | DEFINER } ]
//	    AS select_statement
//	    [ WITH { READ ONLY | CHECK OPTION } [ CONSTRAINT constraint_name ] ]
//	    [ CONTAINER_MAP | CONTAINERS_DEFAULT ]
//
//	CREATE [ OR REPLACE ] MATERIALIZED VIEW [schema.]mview_name
//	    [ IF NOT EXISTS ]
//	    [ { BUILD IMMEDIATE | BUILD DEFERRED } ]
//	    [ REFRESH ... ]
//	    [ ON PREBUILT TABLE ... ]
//	    [ NEVER REFRESH ]
//	    [ ENABLE | DISABLE QUERY REWRITE ]
//	    AS select_statement
func (p *Parser) parseCreateViewStmt(start int, orReplace bool) (*nodes.CreateViewStmt, error) {
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
		// [FORCE | NOFORCE] [EDITIONING] [EDITIONABLE | NONEDITIONABLE] VIEW
		if p.cur.Type == kwFORCE {
			stmt.Force = true
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
			stmt.NoForce = true
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NO" {
			// NO FORCE
			stmt.NoForce = true
			p.advance()
			if p.cur.Type == kwFORCE {
				p.advance()
			}
		}

		// EDITIONING
		if p.isIdentLike() && p.cur.Str == "EDITIONING" {
			stmt.Editioning = true
			p.advance()
		}

		// EDITIONABLE | NONEDITIONABLE
		if p.isIdentLike() && p.cur.Str == "EDITIONABLE" {
			stmt.Editionable = "EDITIONABLE"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NONEDITIONABLE" {
			stmt.Editionable = "NONEDITIONABLE"
			p.advance()
		}

		if p.cur.Type == kwVIEW {
			p.advance()
		}
	}

	return p.finishCreateViewStmt(stmt)
}

// finishCreateViewStmt finishes parsing a CREATE VIEW statement after the
// MATERIALIZED/FORCE/VIEW prefix has been consumed.
func (p *Parser) finishCreateViewStmt(stmt *nodes.CreateViewStmt) (*nodes.CreateViewStmt, error) {
	// IF NOT EXISTS (before view name)
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if p.cur.Type == kwEXISTS {
			p.advance() // consume EXISTS
		}
		stmt.IfNotExists = true
	}

	// View name
	var err error
	stmt.Name, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}

	// SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
	if p.isIdentLike() && p.cur.Str == "SHARING" {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLike() && p.cur.Str == "METADATA" {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "DATA" {
			stmt.Sharing = "DATA"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "EXTENDED" {
			p.advance() // consume EXTENDED
			if p.isIdentLike() && p.cur.Str == "DATA" {
				p.advance() // consume DATA
			}
			stmt.Sharing = "EXTENDED DATA"
		} else if p.isIdentLike() && p.cur.Str == "NONE" {
			stmt.Sharing = "NONE"
			p.advance()
		}
	}

	// Optional column alias list: ( col1, col2, ... )
	if p.cur.Type == '(' && !stmt.Materialized {
		p.advance()
		stmt.Columns = &nodes.List{}
		for {
			if p.cur.Type == ')' || p.cur.Type == tokEOF {
				break
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if name == "" {
				break
			}
			stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
			// consume VISIBLE/INVISIBLE and inline constraints
			for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF && p.cur.Type != kwAS {
				if p.cur.Type == kwCONSTRAINT || p.isIdentLikeStr("VISIBLE") || p.isIdentLikeStr("INVISIBLE") {
					p.advance()
					continue
				}
				if p.cur.Type == kwNOT || p.cur.Type == kwNULL || p.cur.Type == kwUNIQUE ||
					p.cur.Type == kwPRIMARY || p.cur.Type == kwREFERENCES || p.cur.Type == kwCHECK {
					p.advance()
					continue
				}
				if p.isIdentLike() {
					p.advance()
					continue
				}
				break
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// DEFAULT COLLATION collation_name
	if p.isIdentLike() && p.cur.Str == "DEFAULT" {
		next := p.peekNext()
		if next.Type == tokIDENT || next.Str == "COLLATION" {
			p.advance() // consume DEFAULT
			if p.isIdentLike() && p.cur.Str == "COLLATION" {
				p.advance() // consume COLLATION
				stmt.DefaultCollation, err = p.parseIdentifier()
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// BEQUEATH { CURRENT_USER | DEFINER }
	if p.isIdentLike() && p.cur.Str == "BEQUEATH" {
		p.advance() // consume BEQUEATH
		if p.isIdentLike() && (p.cur.Str == "CURRENT_USER" || p.cur.Str == "CURRENT") {
			stmt.Bequeath = "CURRENT_USER"
			p.advance()
			if p.cur.Type == '_' || (p.isIdentLike() && p.cur.Str == "USER") {
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "DEFINER" {
			stmt.Bequeath = "DEFINER"
			p.advance()
		}
	}

	// CONTAINER_MAP | CONTAINERS_DEFAULT (can appear before AS)
	if p.isIdentLike() && p.cur.Str == "CONTAINER_MAP" {
		stmt.ContainerMap = true
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "CONTAINERS_DEFAULT" {
		stmt.ContainersDefault = true
		p.advance()
	}

	// Materialized view options (before AS)
	if stmt.Materialized {
		if err := p.parseMaterializedViewOptions(stmt); err != nil {
			return nil, err
		}
	}

	// AS select_statement
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Query, err = p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
	}

	// WITH CHECK OPTION | WITH READ ONLY [CONSTRAINT name]
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
		// CONSTRAINT constraint_name
		if p.cur.Type == kwCONSTRAINT {
			p.advance()
			stmt.ConstraintName, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMaterializedViewOptions parses BUILD, REFRESH, and other options for materialized views.
func (p *Parser) parseMaterializedViewOptions(stmt *nodes.CreateViewStmt) error {
	for {
		switch {
		case p.isIdentLike() && p.cur.Str == "BUILD":
			p.advance() // consume BUILD
			if p.cur.Type == kwIMMEDIATE || (p.isIdentLike() && p.cur.Str == "IMMEDIATE") {
				stmt.BuildMode = "IMMEDIATE"
				p.advance()
			} else if p.cur.Type == kwDEFERRED || (p.isIdentLike() && p.cur.Str == "DEFERRED") {
				stmt.BuildMode = "DEFERRED"
				p.advance()
			}

		case p.cur.Type == kwON && p.isIdentLikeStrAt(p.peekNext(), "PREBUILT"):
			// ON PREBUILT TABLE [ { WITH | WITHOUT } REDUCED PRECISION ]
			p.advance() // consume ON
			p.advance() // consume PREBUILT
			if p.isIdentLike() && p.cur.Str == "TABLE" {
				p.advance()
			}
			stmt.OnPrebuilt = true
			if p.cur.Type == kwWITH {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "REDUCED" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "PRECISION" {
						p.advance()
					}
				}
				stmt.ReducedPrec = "WITH_REDUCED"
			} else if p.isIdentLike() && p.cur.Str == "WITHOUT" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "REDUCED" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "PRECISION" {
						p.advance()
					}
				}
				stmt.ReducedPrec = "WITHOUT_REDUCED"
			}

		case p.isIdentLike() && p.cur.Str == "NEVER":
			// NEVER REFRESH
			p.advance()
			if p.cur.Type == kwREFRESH {
				p.advance()
			}
			stmt.NeverRefresh = true

		case p.cur.Type == kwREFRESH:
			p.advance()
			// FAST | COMPLETE | FORCE
			if p.isIdentLike() && p.cur.Str == "FAST" {
				stmt.RefreshMethod = "FAST"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "COMPLETE" {
				stmt.RefreshMethod = "COMPLETE"
				p.advance()
			} else if p.cur.Type == kwFORCE {
				stmt.RefreshMethod = "FORCE"
				p.advance()
			}
			// ON COMMIT | ON DEMAND | ON STATEMENT
			if p.cur.Type == kwON {
				p.advance()
				if p.cur.Type == kwCOMMIT {
					stmt.RefreshMode = "ON COMMIT"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "DEMAND" {
					stmt.RefreshMode = "ON DEMAND"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "STATEMENT" {
					stmt.RefreshMode = "ON STATEMENT"
					stmt.RefreshOnStmt = true
					p.advance()
				}
			}
			// START WITH expr
			if p.cur.Type == kwSTART {
				p.advance()
				if p.cur.Type == kwWITH {
					p.advance()
				}
				var parseErr592 error
				stmt.StartWith, parseErr592 = p.parseExpr()
				if parseErr592 !=

					// NEXT expr
					nil {
					return parseErr592
				}
			}

			if p.cur.Type == kwNEXT {
				p.advance()
				var parseErr593 error
				stmt.Next, parseErr593 = p.parseExpr()
				if parseErr593 !=

					// WITH PRIMARY KEY | WITH ROWID
					nil {
					return parseErr593
				}
			}

			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwPRIMARY {
					p.advance()
					if p.cur.Type == kwKEY {
						p.advance()
					}
					stmt.WithPK = true
				} else if p.cur.Type == kwROWID {
					stmt.WithRowID = true
					p.advance()
				}
			}
			// USING ... (skip)
			if p.cur.Type == kwUSING {
				p.advance()
				// skip rollback segment or constraints spec
				for p.cur.Type != ';' && p.cur.Type != tokEOF &&
					p.cur.Type != kwENABLE && p.cur.Type != kwDISABLE &&
					p.cur.Type != kwAS && !(p.isIdentLike() && p.cur.Str == "BUILD") &&
					p.cur.Type != kwREFRESH && p.cur.Type != kwCACHE &&
					p.cur.Type != kwNOCACHE && p.cur.Type != kwPARALLEL &&
					p.cur.Type != kwNOPARALLEL {
					p.advance()
				}
			}

		case p.cur.Type == kwENABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "QUERY" {
				p.advance()
				if p.cur.Type == kwREWRITE {
					p.advance()
				}
				stmt.EnableQuery = true
			} else if p.cur.Type == kwON {
				// ENABLE ON QUERY COMPUTATION
				p.advance()
				if p.isIdentLike() && p.cur.Str == "QUERY" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "COMPUTATION" {
						p.advance()
					}
				}
				stmt.EnableOnQueryComputation = true
			} else {
				stmt.EnableQuery = true
			}

		case p.cur.Type == kwDISABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "QUERY" {
				p.advance()
				if p.cur.Type == kwREWRITE {
					p.advance()
				}
				stmt.DisableQuery = true
			} else if p.cur.Type == kwON {
				// DISABLE ON QUERY COMPUTATION
				p.advance()
				if p.isIdentLike() && p.cur.Str == "QUERY" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "COMPUTATION" {
						p.advance()
					}
				}
				stmt.DisableOnQueryComputation = true
			} else {
				stmt.DisableQuery = true
			}

		case p.cur.Type == kwCACHE:
			stmt.CacheMode = "CACHE"
			p.advance()

		case p.cur.Type == kwNOCACHE:
			stmt.CacheMode = "NOCACHE"
			p.advance()

		case p.cur.Type == kwPARALLEL:
			stmt.ParallelMode = "PARALLEL"
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.ParallelDegree = p.cur.Str
				p.advance()
			}

		case p.cur.Type == kwNOPARALLEL:
			stmt.ParallelMode = "NOPARALLEL"
			p.advance()

		case p.cur.Type == kwLOGGING:
			p.advance()

		case p.cur.Type == kwNOLOGGING:
			p.advance()

		case p.isIdentLike() && p.cur.Str == "SEGMENT":
			// SEGMENT CREATION IMMEDIATE | DEFERRED
			p.advance()
			if p.isIdentLike() && p.cur.Str == "CREATION" {
				p.advance()
				if p.cur.Type == kwIMMEDIATE || p.cur.Type == kwDEFERRED {
					p.advance()
				}
			}

		case p.isIdentLike() && (p.cur.Str == "PCTFREE" || p.cur.Str == "PCTUSED" ||
			p.cur.Str == "INITRANS" || p.cur.Str == "MAXTRANS" ||
			p.cur.Str == "TABLESPACE" || p.cur.Str == "STORAGE"):
			// physical attributes - skip keyword and value
			p.advance()
			if p.cur.Type == tokICONST || p.isIdentLike() {
				p.advance()
				// for STORAGE (...) or TABLESPACE name
				if p.cur.Type == '(' {
					p.advance()
					depth := 1
					for depth > 0 && p.cur.Type != tokEOF {
						if p.cur.Type == '(' {
							depth++
						} else if p.cur.Type == ')' {
							depth--
						}
						p.advance()
					}
				}
			}

		default:
			return nil
		}
	}
	return nil
}

// isIdentLikeStrAt checks if the given token is an identifier-like token with the given string.
func (p *Parser) isIdentLikeStrAt(tok Token, s string) bool {
	if tok.Type == tokIDENT || tok.Type >= kwAS {
		return tok.Str == s
	}
	return false
}

// parseCreateMviewLogStmt parses a CREATE MATERIALIZED VIEW LOG statement.
// Called after CREATE MATERIALIZED VIEW LOG has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-MATERIALIZED-VIEW-LOG.html
//
//	CREATE MATERIALIZED VIEW LOG
//	    [ IF NOT EXISTS ]
//	    ON [ schema. ] table
//	    [ SHARING = { METADATA | NONE } ]
//	    [ physical_attributes_clause ]
//	    [ TABLESPACE tablespace ]
//	    [ logging_clause ]
//	    [ { CACHE | NOCACHE } ]
//	    [ parallel_clause ]
//	    [ WITH [ { PRIMARY KEY | ROWID | OBJECT ID } ]
//	            [ SEQUENCE ]
//	            [ ( column [, column ]... ) ]
//	            [ COMMIT SCN ]
//	    ]
//	    [ new_values_clause ]
//	    [ mv_log_purge_clause ]
//	    [ for_refresh_clause ]
func (p *Parser) parseCreateMviewLogStmt(start int) (*nodes.CreateMviewLogStmt, error) {
	stmt := &nodes.CreateMviewLogStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if p.cur.Type == kwEXISTS {
			p.advance()
		}
		stmt.IfNotExists = true
	}

	// ON table
	if p.cur.Type == kwON {
		p.advance()
		var parseErr594 error
		stmt.OnTable, parseErr594 = p.parseObjectName()
		if parseErr594 !=

			// SHARING = { METADATA | NONE }
			nil {
			return nil, parseErr594
		}
	}

	if p.isIdentLike() && p.cur.Str == "SHARING" {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "METADATA" {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NONE" {
			stmt.Sharing = "NONE"
			p.advance()
		}
	}

	// Skip physical attributes, TABLESPACE, logging, CACHE, PARALLEL until WITH or new_values or purge
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwWITH || p.isIdentLike() && (p.cur.Str == "INCLUDING" || p.cur.Str == "EXCLUDING") ||
			p.isIdentLike() && p.cur.Str == "PURGE" || p.isIdentLike() && p.cur.Str == "FOR" {
			break
		}
		if p.isIdentLike() && (p.cur.Str == "PCTFREE" || p.cur.Str == "PCTUSED" ||
			p.cur.Str == "INITRANS" || p.cur.Str == "MAXTRANS" || p.cur.Str == "TABLESPACE" ||
			p.cur.Str == "STORAGE") {
			p.advance()
			if p.cur.Type == tokICONST || p.isIdentLike() {
				p.advance()
			}
			if p.cur.Type == '(' {
				p.advance()
				depth := 1
				for depth > 0 && p.cur.Type != tokEOF {
					if p.cur.Type == '(' {
						depth++
					} else if p.cur.Type == ')' {
						depth--
					}
					p.advance()
				}
			}
			continue
		}
		if p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING ||
			p.cur.Type == kwCACHE || p.cur.Type == kwNOCACHE ||
			p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL {
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			}
			continue
		}
		break
	}

	// WITH [ { PRIMARY KEY | ROWID | OBJECT ID } ] [ SEQUENCE ] [ (cols) ] [ COMMIT SCN ]
	if p.cur.Type == kwWITH {
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			if p.cur.Type == kwPRIMARY {
				p.advance()
				if p.cur.Type == kwKEY {
					p.advance()
				}
				stmt.WithPK = true
			} else if p.cur.Type == kwROWID {
				stmt.WithRowID = true
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "OBJECT" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "ID" {
					p.advance()
				}
				stmt.WithOID = true
			} else if p.cur.Type == kwSEQUENCE {
				stmt.WithSeq = true
				p.advance()
			} else if p.cur.Type == '(' {
				// column list
				p.advance()
				stmt.Columns = &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					name, parseErr595 := p.parseIdentifier()
					if parseErr595 != nil {
						return nil, parseErr595
					}
					if name != "" {
						stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
					}
					if p.cur.Type == ',' {
						p.advance()
					} else {
						break
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			} else if p.cur.Type == kwCOMMIT {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "SCN" {
					p.advance()
				}
				stmt.CommitSCN = true
			} else if p.cur.Type == ',' {
				p.advance()
			} else {
				break
			}
		}
	}

	// new_values_clause: INCLUDING NEW VALUES | EXCLUDING NEW VALUES
	if p.isIdentLike() && p.cur.Str == "INCLUDING" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "NEW" {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "VALUES" {
			p.advance()
		}
		stmt.Including = true
	} else if p.isIdentLike() && p.cur.Str == "EXCLUDING" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "NEW" {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "VALUES" {
			p.advance()
		}
		stmt.Excluding = true
	}

	// mv_log_purge_clause: PURGE { IMMEDIATE [ SYNCHRONOUS | ASYNCHRONOUS ] | START WITH expr [ NEXT expr ] }
	if p.cur.Type == kwPURGE {
		p.advance()
		if p.cur.Type == kwIMMEDIATE {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "SYNCHRONOUS" {
				stmt.PurgeMode = "IMMEDIATE_SYNC"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "ASYNCHRONOUS" {
				stmt.PurgeMode = "IMMEDIATE_ASYNC"
				p.advance()
			} else {
				stmt.PurgeMode = "IMMEDIATE_SYNC"
			}
		} else if p.cur.Type == kwSTART {
			stmt.PurgeMode = "START_WITH"
			p.advance()
			if p.cur.Type == kwWITH {
				p.advance()
			}
			var parseErr596 error
			stmt.PurgeStart, parseErr596 = p.parseExpr()
			if parseErr596 != nil {
				return nil, parseErr596
			}
			if p.cur.Type == kwNEXT {
				p.advance()
				var parseErr597 error
				stmt.PurgeNext, parseErr597 = p.parseExpr()
				if parseErr597 != nil {

					// for_refresh_clause: FOR { FAST | SYNCHRONOUS } REFRESH [ USING staging_log ]
					return nil, parseErr597
				}
			}
		}
	}

	if p.isIdentLike() && p.cur.Str == "FOR" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "FAST" {
			stmt.ForRefresh = "FAST"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "SYNCHRONOUS" {
			stmt.ForRefresh = "SYNCHRONOUS"
			p.advance()
		}
		if p.cur.Type == kwREFRESH {
			p.advance()
		}
		if p.cur.Type == kwUSING {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "STAGING" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "LOG" {
					p.advance()
				}
			}
			if p.isIdentLike() {
				stmt.StagingLog = p.cur.Str
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterMviewLogStmt parses an ALTER MATERIALIZED VIEW LOG statement.
// Called after ALTER MATERIALIZED VIEW LOG has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-MATERIALIZED-VIEW-LOG.html
//
//	ALTER MATERIALIZED VIEW LOG [ IF EXISTS ]
//	    [ FORCE ]
//	    ON [ schema. ] table
//	    { physical_attributes_clause | add_mv_log_column_clause | parallel_clause | logging_clause |
//	      allocate_extent_clause | shrink_clause | move_mv_log_clause |
//	      mv_log_augmentation | mv_log_purge_clause | for_refresh_clause }
func (p *Parser) parseAlterMviewLogStmt(start int) (*nodes.AlterMviewLogStmt, error) {
	stmt := &nodes.AlterMviewLogStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}

	// FORCE
	if p.cur.Type == kwFORCE {
		stmt.Force = true
		p.advance()
	}

	// ON table
	if p.cur.Type == kwON {
		p.advance()
		var parseErr598 error
		stmt.OnTable, parseErr598 = p.parseObjectName()
		if parseErr598 !=

			// Parse action
			nil {
			return nil, parseErr598
		}
	}

	switch {
	case p.cur.Type == kwADD:
		p.advance()
		stmt.Action = "ADD"
		// add_mv_log_column_clause: ADD ( column [,...] )
		if p.cur.Type == '(' {
			p.advance()
			stmt.Columns = &nodes.List{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				name, parseErr599 := p.parseIdentifier()
				if parseErr599 != nil {
					return nil, parseErr599
				}
				if name != "" {
					stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
				}
				if p.cur.Type == ',' {
					p.advance()
				} else {
					break
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		} else if p.cur.Type == kwPRIMARY {
			stmt.Action = "ADD_PRIMARY_KEY"
			p.advance()
			if p.cur.Type == kwKEY {
				p.advance()
			}
		} else if p.cur.Type == kwROWID {
			stmt.Action = "ADD_ROWID"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "OBJECT" {
			stmt.Action = "ADD_OBJECT_ID"
			p.advance()
			if p.isIdentLike() && p.cur.Str == "ID" {
				p.advance()
			}
		} else if p.cur.Type == kwSEQUENCE {
			stmt.Action = "ADD_SEQUENCE"
			p.advance()
		}

	case p.isIdentLike() && p.cur.Str == "SHRINK":
		stmt.Action = "SHRINK"
		p.advance()
		if p.isIdentLike() && p.cur.Str == "SPACE" {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "COMPACT" {
			p.advance()
		} else if p.cur.Type == kwCASCADE {
			p.advance()
		}

	case p.isIdentLike() && p.cur.Str == "MOVE":
		stmt.Action = "MOVE"
		p.advance()
		// skip optional tablespace
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}

	case p.cur.Type == kwPURGE:
		stmt.Action = "PURGE"
		p.advance()
		if p.cur.Type == kwIMMEDIATE {
			p.advance()
		} else if p.cur.Type == kwSTART {
			p.advance()
			if p.cur.Type == kwWITH {
				p.advance()
			}
			parseDiscard601, parseErr600 := p.parseExpr()
			_ = // skip expr
				parseDiscard601
			if parseErr600 != nil {
				return nil, parseErr600
			}
			if p.cur.Type == kwNEXT {
				p.advance()
				parseDiscard603, parseErr602 := p.parseExpr()
				_ = parseDiscard603
				if parseErr602 != nil {
					return nil, parseErr602
				}
			}
		}

	case p.cur.Type == kwPARALLEL:
		stmt.Action = "PARALLEL"
		p.advance()
		if p.cur.Type == tokICONST {
			p.advance()
		}

	case p.cur.Type == kwNOPARALLEL:
		stmt.Action = "NOPARALLEL"
		p.advance()

	case p.cur.Type == kwLOGGING:
		stmt.Action = "LOGGING"
		p.advance()

	case p.cur.Type == kwNOLOGGING:
		stmt.Action = "NOLOGGING"
		p.advance()

	case p.isIdentLike() && p.cur.Str == "ALLOCATE":
		stmt.Action = "ALLOCATE_EXTENT"
		p.advance()
		if p.isIdentLike() && p.cur.Str == "EXTENT" {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
				}
				p.advance()
			}
		}

	case p.isIdentLike() && p.cur.Str == "FOR":
		stmt.Action = "FOR_REFRESH"
		p.advance()
		if p.isIdentLike() && (p.cur.Str == "FAST" || p.cur.Str == "SYNCHRONOUS") {
			p.advance()
		}
		if p.cur.Type == kwREFRESH {
			p.advance()
		}
		if p.cur.Type == kwUSING {
			p.advance()
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
		}

	default:
		// skip unknown clauses
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateAnalyticViewStmt parses a CREATE ANALYTIC VIEW statement.
// Called after CREATE [OR REPLACE] [FORCE|NOFORCE] ANALYTIC VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-ANALYTIC-VIEW.html
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] ANALYTIC VIEW
//	    [ IF NOT EXISTS ] [ schema. ] analytic_view_name
//	    [ SHARING = { METADATA | NONE } ]
//	    using_clause
//	    dim_by_clause
//	    measures_clause
//	    [ default_measure_clause ]
//	    [ default_aggregate_clause ]
//	    [ cache_clause ]
//	    [ fact_columns_clause ]
//	    [ qry_transform_clause ]
func (p *Parser) parseCreateAnalyticViewStmt(start int, orReplace, force, noForce bool) (*nodes.CreateAnalyticViewStmt, error) {
	stmt := &nodes.CreateAnalyticViewStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance()
		p.advance()
		if p.cur.Type == kwEXISTS {
			p.advance()
		}
		stmt.IfNotExists = true
	}
	var parseErr604 error

	// view name
	stmt.Name, parseErr604 = p.parseObjectName()
	if parseErr604 !=

		// SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr604
	}

	if p.isIdentLike() && p.cur.Str == "SHARING" {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "METADATA" {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NONE" {
			stmt.Sharing = "NONE"
			p.advance()
		}
	}

	// USING table_name [ AS alias ]
	if p.cur.Type == kwUSING {
		p.advance()
		var parseErr605 error
		stmt.UsingTable, parseErr605 = p.parseObjectName()
		if parseErr605 != nil {
			return nil, parseErr605
		}
		if p.cur.Type == kwAS {
			p.advance()
			var parseErr606 error
			stmt.UsingAlias, parseErr606 = p.parseIdentifier()
			if parseErr606 !=

				// DIMENSION BY ( ... )
				nil {
				return nil, parseErr606
			}
		}
	}

	if p.isIdentLike() && p.cur.Str == "DIMENSION" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "BY" {
			p.advance()
		}
		stmt.DimBy = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			var buf string
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				if buf != "" {
					buf += " "
				}
				buf += p.cur.Str
				p.advance()
			}
			if buf != "" {
				stmt.DimBy.Items = append(stmt.DimBy.Items, &nodes.String{Str: buf})
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// MEASURES ( ... )
	if p.isIdentLike() && p.cur.Str == "MEASURES" {
		p.advance()
		stmt.Measures = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			var buf string
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				if buf != "" {
					buf += " "
				}
				buf += p.cur.Str
				p.advance()
			}
			if buf != "" {
				stmt.Measures.Items = append(stmt.Measures.Items, &nodes.String{Str: buf})
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// DEFAULT MEASURE name
	if p.isIdentLike() && p.cur.Str == "DEFAULT" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "MEASURE" {
			p.advance()
			var parseErr607 error
			stmt.DefaultMeasure, parseErr607 = p.parseIdentifier()
			if parseErr607 != nil {
				return nil, parseErr607
			}
		} else if p.isIdentLike() && p.cur.Str == "AGGREGATE" {
			// DEFAULT AGGREGATE BY function
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BY" {
				p.advance()
			}
			var parseErr608 error
			stmt.DefaultAggregate, parseErr608 = p.parseIdentifier()
			if parseErr608 !=

				// skip remaining (cache, fact, qry_transform)
				nil {
				return nil, parseErr608
			}
		}
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterAnalyticViewStmt parses an ALTER ANALYTIC VIEW statement.
// Called after ALTER ANALYTIC VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-ANALYTIC-VIEW.html
//
//	ALTER ANALYTIC VIEW [ IF EXISTS ] [ schema . ] analytic_view_name
//	    { RENAME TO new_av_name | COMPILE | alter_add_cache_clause | alter_drop_cache_clause }
func (p *Parser) parseAlterAnalyticViewStmt(start int) (*nodes.AlterAnalyticViewStmt, error) {
	stmt := &nodes.AlterAnalyticViewStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var parseErr609 error

	stmt.Name, parseErr609 = p.parseObjectName()
	if parseErr609 != nil {
		return nil, parseErr609
	}

	switch {
	case p.isIdentLike() && p.cur.Str == "COMPILE":
		stmt.Action = "COMPILE"
		p.advance()

	case p.cur.Type == kwRENAME || (p.isIdentLike() && p.cur.Str == "RENAME"):
		stmt.Action = "RENAME"
		p.advance() // consume RENAME
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		var parseErr610 error
		stmt.NewName, parseErr610 = p.parseObjectName()
		if parseErr610 != nil {
			return nil, parseErr610
		}

	case p.isIdentLike() && p.cur.Str == "ADD":
		stmt.Action = "ADD_CACHE"
		p.advance()
		// skip CACHE clause
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}

	case p.isIdentLike() && p.cur.Str == "DROP":
		stmt.Action = "DROP_CACHE"
		p.advance()
		// skip CACHE clause
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}

	default:
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateJsonDualityViewStmt parses a CREATE JSON RELATIONAL DUALITY VIEW statement.
// Called after CREATE [OR REPLACE | IF NOT EXISTS] JSON RELATIONAL DUALITY VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/create-json-relational-duality-view.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] JSON RELATIONAL DUALITY VIEW
//	    [ schema. ] view_name
//	    [ duality_view_replication_clause ]
//	    AS duality_view_subquery
func (p *Parser) parseCreateJsonDualityViewStmt(start int, orReplace, ifNotExists bool) (*nodes.CreateJsonDualityViewStmt, error) {
	stmt := &nodes.CreateJsonDualityViewStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}
	var parseErr611 error

	stmt.Name, parseErr611 = p.parseObjectName()
	if parseErr611 !=

		// duality_view_replication_clause: ENABLE | DISABLE LOGICAL REPLICATION
		nil {
		return nil, parseErr611
	}

	if p.cur.Type == kwENABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "LOGICAL" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REPLICATION" {
				p.advance()
			}
		}
		stmt.EnableLogicalReplication = true
	} else if p.cur.Type == kwDISABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "LOGICAL" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REPLICATION" {
				p.advance()
			}
		}
		stmt.DisableLogicalReplication = true
	}

	// AS duality_view_subquery
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr612 error
		stmt.Query, parseErr612 = p.parseSelectStmt()
		if parseErr612 != nil {
			return nil, parseErr612
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterJsonDualityViewStmt parses an ALTER JSON RELATIONAL DUALITY VIEW statement.
// Called after ALTER JSON RELATIONAL DUALITY VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/alter-json-relational-duality-view.html
//
//	ALTER JSON RELATIONAL DUALITY VIEW [ schema. ] view_name
//	    { ENABLE | DISABLE } LOGICAL REPLICATION
func (p *Parser) parseAlterJsonDualityViewStmt(start int) (*nodes.AlterJsonDualityViewStmt, error) {
	stmt := &nodes.AlterJsonDualityViewStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr613 error

	stmt.Name, parseErr613 = p.parseObjectName()
	if parseErr613 != nil {
		return nil, parseErr613
	}

	if p.cur.Type == kwENABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "LOGICAL" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REPLICATION" {
				p.advance()
			}
		}
		stmt.Action = "ENABLE_LOGICAL_REPLICATION"
	} else if p.cur.Type == kwDISABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "LOGICAL" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REPLICATION" {
				p.advance()
			}
		}
		stmt.Action = "DISABLE_LOGICAL_REPLICATION"
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
