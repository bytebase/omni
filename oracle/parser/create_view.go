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

		if p.cur.Type != kwVIEW {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return p.finishCreatePlainViewStmt(stmt)
	}

	return p.finishCreateViewStmt(stmt)
}

// isViewConstraintStart reports whether the current position starts an
// out-of-line view constraint. CONSTRAINT is unambiguous in a view alias
// list; PRIMARY, FOREIGN, and UNIQUE are legal bare aliases (Oracle 23ai
// accepts CREATE VIEW v (primary) AS ...), so they count only when followed
// by their mandatory next token. CHECK constraints do not exist on views.
func (p *Parser) isViewConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT:
		return true
	case kwPRIMARY, kwFOREIGN:
		return p.peekNext().Type == kwKEY
	case kwUNIQUE:
		return p.peekNext().Type == '('
	}
	return false
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
	stmt.Name, err = p.parseReservedCheckedObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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
				return nil, p.syntaxErrorAtCur()
			}
			// Out-of-line view constraint: CONSTRAINT name { PRIMARY KEY |
			// UNIQUE | FOREIGN KEY } (cols) [RELY|NORELY] DISABLE [NOVALIDATE].
			// View constraints are declarative only; the table-only state
			// clauses are rejected — see parseViewConstraintState.
			// Unlike a CREATE TABLE column list, a view alias list has no
			// datatype ambiguity: CONSTRAINT always starts a constraint here
			// (a constraint name may be a nonreserved datatype word like BLOB,
			// while a bare alias named CONSTRAINT is rejected by Oracle).
			// PRIMARY/FOREIGN/UNIQUE, however, are legal bare aliases
			// (verified on 23ai), so they start a constraint only when the
			// required following token is present.
			if p.isViewConstraintStart() {
				viewConstraint, err := p.parseTableConstraintBody()
				if err != nil {
					return nil, err
				}
				if err := p.parseViewConstraintState(); err != nil {
					return nil, err
				}
				viewConstraint.Loc.End = p.prev.End
				if stmt.Constraints == nil {
					stmt.Constraints = &nodes.List{}
				}
				stmt.Constraints.Items = append(stmt.Constraints.Items, viewConstraint)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
				continue
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if name == "" {
				return nil, p.syntaxErrorAtCur()
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
		if stmt.Columns.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
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
		if !selectStmtHasTargets(stmt.Query) {
			return nil, p.syntaxErrorAtCur()
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
	seenRefresh := false
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
			seenRefresh = true
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

		case p.cur.Type == kwUSING && p.peekNext().Type == kwINDEX:
			// USING INDEX index_properties (default index storage for the
			// mview). Properties only: Oracle rejects an index name or a
			// nested CREATE INDEX here (ORA-02000).
			p.advance() // consume USING
			p.advance() // consume INDEX
			var cs constraintState
			if err := p.parseUsingIndexProperties(&cs); err != nil {
				return err
			}

		case p.cur.Type == kwUSING && p.isIdentLikeStrAt(p.peekNext(), "NO"):
			// USING NO INDEX
			p.advance() // consume USING
			p.advance() // consume NO
			if p.cur.Type != kwINDEX {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume INDEX

		case seenRefresh && p.cur.Type == kwUSING &&
			(p.isIdentLikeStrAt(p.peekNext(), "TRUSTED") || p.isIdentLikeStrAt(p.peekNext(), "ENFORCED")):
			// USING { TRUSTED | ENFORCED } CONSTRAINTS — part of create_mv_refresh;
			// Oracle rejects it without a preceding REFRESH clause.
			p.advance() // consume USING
			p.advance() // consume TRUSTED/ENFORCED
			if p.cur.Type != kwCONSTRAINTS {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume CONSTRAINTS

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
