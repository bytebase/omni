package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseSetStmt parses a SET statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-variable.html
//
//	SET [GLOBAL | SESSION | LOCAL] var = expr [, var = expr] ...
//	SET NAMES charset [COLLATE collation]
//	SET CHARACTER SET charset
func (p *Parser) parseSetStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SET

	stmt := &nodes.SetStmt{Loc: nodes.Loc{Start: start}}

	// Check for NAMES special form
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "names") {
		p.advance() // consume NAMES
		// Parse charset name
		charset, charsetLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Build assignment: NAMES = charset
		stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
			Loc:    nodes.Loc{Start: charsetLoc},
			Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc}, Column: "NAMES"},
			Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc}, Value: charset},
		})
		// Optional COLLATE
		if _, ok := p.match(kwCOLLATE); ok {
			collation, collLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: collLoc},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: collLoc}, Column: "COLLATE"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: collLoc}, Value: collation},
			})
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Check for CHARACTER SET special form
	if p.cur.Type == kwCHARACTER {
		p.advance() // consume CHARACTER
		if _, ok := p.match(kwSET); ok {
			charset, charsetLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: charsetLoc},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc}, Column: "CHARACTER SET"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc}, Value: charset},
			})
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	// Check for SET DEFAULT ROLE
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if p.cur.Type == kwROLE {
			p.advance() // consume ROLE
			return p.parseSetDefaultRoleStmt(start)
		}
		return nil, &ParseError{Message: "expected ROLE after SET DEFAULT", Position: p.cur.Loc}
	}

	// Check for SET ROLE
	if p.cur.Type == kwROLE {
		p.advance() // consume ROLE
		return p.parseSetRoleStmt(start)
	}

	// Check for SET PASSWORD
	if p.cur.Type == kwPASSWORD {
		return p.parseSetPasswordStmt(start)
	}

	// Check for SET RESOURCE GROUP
	if p.cur.Type == kwRESOURCE {
		return p.parseSetResourceGroupStmt(start)
	}

	// Check for GLOBAL / SESSION / LOCAL / PERSIST / PERSIST_ONLY scope
	scope := ""
	switch p.cur.Type {
	case kwGLOBAL:
		scope = "GLOBAL"
		p.advance()
	case kwSESSION:
		scope = "SESSION"
		p.advance()
	case kwLOCAL:
		scope = "LOCAL"
		p.advance()
	case kwPERSIST:
		scope = "PERSIST"
		p.advance()
	default:
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "persist_only") {
			scope = "PERSIST_ONLY"
			p.advance()
		}
	}

	// SET [GLOBAL|SESSION] TRANSACTION ...
	if p.cur.Type == kwTRANSACTION {
		return p.parseSetTransactionStmt(start, scope)
	}

	stmt.Scope = scope

	// Parse assignment list
	for {
		asgn, err := p.parseSetAssignment()
		if err != nil {
			return nil, err
		}
		stmt.Assignments = append(stmt.Assignments, asgn)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSetPasswordStmt parses a SET PASSWORD statement.
// SET has already been consumed; p.cur is PASSWORD.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-password.html
//
//	SET PASSWORD [FOR user] = 'auth_string'
//	SET PASSWORD [FOR user] = PASSWORD('auth_string')
func (p *Parser) parseSetPasswordStmt(start int) (*nodes.SetPasswordStmt, error) {
	p.advance() // consume PASSWORD

	stmt := &nodes.SetPasswordStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional FOR user[@host]
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		user, err := p.parseUserSpec()
		if err != nil {
			return nil, err
		}
		stmt.User = user
	}

	// consume '='
	p.match('=')

	// Password value: either a string literal or PASSWORD('string')
	if p.isIdentToken() && eqFold(p.cur.Str, "PASSWORD") {
		// PASSWORD('auth_string') form
		p.advance() // consume PASSWORD
		p.match('(')
		if p.cur.Type == tokSCONST {
			stmt.Password = "PASSWORD(" + p.cur.Str + ")"
			p.advance()
		}
		p.match(')')
	} else if p.cur.Type == tokSCONST {
		stmt.Password = p.cur.Str
		p.advance()
	} else {
		return nil, &ParseError{Message: "expected password string", Position: p.cur.Loc}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSetAssignment parses a single SET assignment: var = expr
func (p *Parser) parseSetAssignment() (*nodes.Assignment, error) {
	start := p.pos()

	var col *nodes.ColumnRef

	// Handle @var and @@var references
	if p.isVariableRef() {
		vref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		// Convert VariableRef to ColumnRef for the assignment
		prefix := "@"
		if vref.System {
			prefix = "@@"
			if vref.Scope != "" {
				prefix = "@@" + vref.Scope + "."
			}
		}
		col = &nodes.ColumnRef{
			Loc:    vref.Loc,
			Column: prefix + vref.Name,
		}
	} else {
		var err error
		col, err = p.parseColumnRef()
		if err != nil {
			return nil, err
		}
	}

	// Expect '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// Parse value expression
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Column: col,
		Value:  val,
	}, nil
}

// parseShowStmt parses a SHOW statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/show.html
func (p *Parser) parseShowStmt() (*nodes.ShowStmt, error) {
	start := p.pos()
	p.advance() // consume SHOW

	stmt := &nodes.ShowStmt{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwDATABASES:
		stmt.Type = "DATABASES"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwTABLES:
		stmt.Type = "TABLES"
		p.advance()
		// Optional FROM db
		if _, ok := p.match(kwFROM); ok {
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwFULL:
		p.advance() // consume FULL
		if p.cur.Type == kwCOLUMNS || p.cur.Type == kwFIELDS {
			stmt.Type = "FULL COLUMNS"
			p.advance()
			// FROM|IN tbl
			if _, err := p.expectFromOrIn(); err != nil {
				return nil, err
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			// Optional FROM|IN db
			if p.matchFromOrIn() {
				dbRef, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				// Merge: set schema on From
				stmt.From.Schema = dbRef.Name
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwTABLES {
			stmt.Type = "FULL TABLES"
			p.advance()
			// Optional FROM db
			if _, ok := p.match(kwFROM); ok {
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From = ref
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwPROCESSLIST || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "processlist")) {
			stmt.Type = "FULL PROCESSLIST"
			p.advance()
		}

	case kwCOLUMNS, kwFIELDS:
		stmt.Type = "COLUMNS"
		p.advance()
		// FROM|IN tbl
		if _, err := p.expectFromOrIn(); err != nil {
			return nil, err
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.From = ref
		// Optional FROM|IN db
		if p.matchFromOrIn() {
			dbRef, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From.Schema = dbRef.Name
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwCREATE:
		p.advance() // consume CREATE
		switch p.cur.Type {
		case kwTABLE:
			stmt.Type = "CREATE TABLE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwDATABASE:
			stmt.Type = "CREATE DATABASE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwVIEW:
			stmt.Type = "CREATE VIEW"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwPROCEDURE:
			stmt.Type = "CREATE PROCEDURE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwFUNCTION:
			stmt.Type = "CREATE FUNCTION"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwTRIGGER:
			stmt.Type = "CREATE TRIGGER"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwEVENT:
			stmt.Type = "CREATE EVENT"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwUSER:
			stmt.Type = "CREATE USER"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			// Skip optional @'host' part
			if p.cur.Type == '@' {
				p.advance()
				p.advance() // skip host
			}
		}

	case kwEXTENDED:
		p.advance() // consume EXTENDED
		if p.cur.Type == kwINDEX || p.cur.Type == kwKEY || p.cur.Type == kwKEYS || (p.cur.Type == tokIDENT && (eqFold(p.cur.Str, "indexes") || eqFold(p.cur.Str, "keys"))) {
			// SHOW EXTENDED {INDEX | INDEXES | KEYS} ...
			stmt.Type = "EXTENDED INDEX"
			p.advance()
			if _, err := p.expectFromOrIn(); err != nil {
				return nil, err
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			if p.matchFromOrIn() {
				dbRef, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From.Schema = dbRef.Name
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwFULL {
			p.advance() // consume FULL
			// SHOW EXTENDED FULL COLUMNS ...
			if p.cur.Type == kwCOLUMNS || p.cur.Type == kwFIELDS || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "fields")) {
				stmt.Type = "EXTENDED FULL COLUMNS"
				p.advance()
				if _, err := p.expectFromOrIn(); err != nil {
					return nil, err
				}
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From = ref
				if p.matchFromOrIn() {
					dbRef, err := p.parseTableRef()
					if err != nil {
						return nil, err
					}
					stmt.From.Schema = dbRef.Name
				}
				if err := p.parseShowLikeOrWhere(stmt); err != nil {
					return nil, err
				}
			}
		} else if p.cur.Type == kwCOLUMNS || p.cur.Type == kwFIELDS || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "fields")) {
			// SHOW EXTENDED COLUMNS ...
			stmt.Type = "EXTENDED COLUMNS"
			p.advance()
			if _, err := p.expectFromOrIn(); err != nil {
				return nil, err
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			if p.matchFromOrIn() {
				dbRef, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From.Schema = dbRef.Name
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwINDEX, kwKEY:
		stmt.Type = "INDEX"
		p.advance()
		// FROM|IN tbl
		if _, err := p.expectFromOrIn(); err != nil {
			return nil, err
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.From = ref
		// Optional FROM|IN db
		if p.matchFromOrIn() {
			dbRef, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From.Schema = dbRef.Name
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwGLOBAL, kwSESSION:
		scope := "GLOBAL"
		if p.cur.Type == kwSESSION {
			scope = "SESSION"
		}
		p.advance()
		if p.cur.Type == kwVARIABLES {
			stmt.Type = scope + " VARIABLES"
			p.advance()
		} else if p.cur.Type == kwSTATUS {
			stmt.Type = scope + " STATUS"
			p.advance()
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwVARIABLES:
		stmt.Type = "VARIABLES"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwSTATUS:
		stmt.Type = "STATUS"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwWARNINGS:
		stmt.Type = "WARNINGS"
		p.advance()
		// Optional LIMIT
		if _, ok := p.match(kwLIMIT); ok {
			// Just skip the count for now
			p.advance()
		}

	case kwERRORS:
		stmt.Type = "ERRORS"
		p.advance()
		// Optional LIMIT
		if _, ok := p.match(kwLIMIT); ok {
			p.advance()
		}

	case kwENGINES:
		stmt.Type = "ENGINES"
		p.advance()

	case kwENGINE:
		// SHOW ENGINE engine_name {STATUS | MUTEX}
		p.advance() // consume ENGINE
		engineName, nameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = &nodes.TableRef{
			Loc:  nodes.Loc{Start: nameLoc, End: p.pos()},
			Name: engineName,
		}
		if p.cur.Type == kwSTATUS {
			stmt.Type = "ENGINE STATUS"
			p.advance()
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "mutex") {
			stmt.Type = "ENGINE MUTEX"
			p.advance()
		}

	case kwPLUGINS:
		stmt.Type = "PLUGINS"
		p.advance()

	case kwMASTER:
		p.advance() // consume MASTER
		if p.cur.Type == kwSTATUS {
			stmt.Type = "MASTER STATUS"
			p.advance()
		}

	case kwSLAVE:
		p.advance() // consume SLAVE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "SLAVE STATUS"
			p.advance()
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "hosts") {
			stmt.Type = "REPLICAS"
			p.advance()
		}

	case kwREPLICA:
		p.advance() // consume REPLICA
		if p.cur.Type == kwSTATUS {
			stmt.Type = "REPLICA STATUS"
			p.advance()
		}

	case kwBINARY:
		p.advance() // consume BINARY
		if p.cur.Type == kwLOGS {
			stmt.Type = "BINARY LOGS"
			p.advance()
		}

	case kwBINLOG:
		p.advance() // consume BINLOG
		stmt.Type = "BINLOG EVENTS"
		// Expect EVENTS (as identifier since it may not be a keyword)
		if p.cur.Type == kwEVENT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events")) {
			p.advance()
		}
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			// skip the position value
			p.advance()
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			// skip limit value(s)
			p.advance()
			if p.cur.Type == ',' {
				p.advance() // consume ','
				p.advance() // consume count
			}
		}

	case kwTABLE:
		p.advance() // consume TABLE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "TABLE STATUS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwTRIGGER:
		// SHOW TRIGGERS (kwTRIGGER won't match plural; handled in default)
		stmt.Type = "TRIGGERS"
		p.advance()
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwEVENT:
		// SHOW EVENTS (kwEVENT won't match plural; handled in default)
		stmt.Type = "EVENTS"
		p.advance()
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwPROCEDURE:
		p.advance() // consume PROCEDURE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "PROCEDURE STATUS"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "code") {
			stmt.Type = "PROCEDURE CODE"
			p.advance() // consume CODE
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		}

	case kwFUNCTION:
		p.advance() // consume FUNCTION
		if p.cur.Type == kwSTATUS {
			stmt.Type = "FUNCTION STATUS"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "code") {
			stmt.Type = "FUNCTION CODE"
			p.advance() // consume CODE
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		}

	case kwOPEN:
		p.advance() // consume OPEN
		if p.cur.Type == kwTABLES {
			stmt.Type = "OPEN TABLES"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwPRIVILEGES:
		stmt.Type = "PRIVILEGES"
		p.advance()

	case kwPROFILES:
		stmt.Type = "PROFILES"
		p.advance()

	case kwRELAYLOG:
		p.advance() // consume RELAYLOG
		stmt.Type = "RELAYLOG EVENTS"
		if p.cur.Type == kwEVENT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events")) {
			p.advance()
		}
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			p.advance()
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
				p.advance()
			}
		}

	case kwCHARACTER:
		p.advance() // consume CHARACTER
		if p.cur.Type == kwSET {
			stmt.Type = "CHARACTER SET"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwCOLLATION:
		stmt.Type = "COLLATION"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	default:
		// SHOW PROFILE [type [, type] ...] [FOR QUERY n] [LIMIT row_count [OFFSET offset]]
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "profile") {
			stmt.Type = "PROFILE"
			p.advance() // consume PROFILE
			if err := p.parseShowProfileOptions(stmt); err != nil {
				return nil, err
			}
			stmt.Loc.End = p.pos()
			return stmt, nil
		}

		// SHOW REPLICAS
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "replicas") {
			stmt.Type = "REPLICAS"
			p.advance()
			stmt.Loc.End = p.pos()
			return stmt, nil
		}

		// Handle GRANTS and PROCESSLIST as identifier-based keywords
		if p.cur.Type == kwGRANT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "grants")) {
			stmt.Type = "GRANTS"
			p.advance()
			// Optional FOR user_or_role
			if _, ok := p.match(kwFOR); ok {
				// Handle CURRENT_USER / CURRENT_USER()
				if p.cur.Type == kwCURRENT_USER {
					start := p.pos()
					p.advance()
					// Optional ()
					if p.cur.Type == '(' {
						p.advance()
						p.match(')')
					}
					stmt.ForUser = &nodes.UserSpec{
						Loc:  nodes.Loc{Start: start, End: p.pos()},
						Name: "CURRENT_USER",
					}
				} else {
					user, err := p.parseUserSpec()
					if err != nil {
						return nil, err
					}
					stmt.ForUser = user
				}
				// Optional USING role [, role] ...
				if _, ok := p.match(kwUSING); ok {
					for {
						role, err := p.parseUserSpec()
						if err != nil {
							return nil, err
						}
						stmt.Using = append(stmt.Using, role)
						if p.cur.Type != ',' {
							break
						}
						p.advance() // consume ','
					}
				}
			}
		} else if p.cur.Type == kwPROCESSLIST || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "processlist")) {
			stmt.Type = "PROCESSLIST"
			p.advance()
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "triggers") {
			stmt.Type = "TRIGGERS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events") {
			stmt.Type = "EVENTS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwKEYS || (p.cur.Type == tokIDENT && (eqFold(p.cur.Str, "indexes") || eqFold(p.cur.Str, "keys"))) {
			// SHOW INDEXES|KEYS (synonyms for SHOW INDEX)
			stmt.Type = "INDEX"
			p.advance()
			if _, err := p.expectFromOrIn(); err != nil {
				return nil, err
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			if p.matchFromOrIn() {
				dbRef, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From.Schema = dbRef.Name
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "schemas") {
			// SHOW SCHEMAS (synonym for SHOW DATABASES)
			stmt.Type = "DATABASES"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseShowLikeOrWhere parses optional LIKE or WHERE clause for SHOW statements.
func (p *Parser) parseShowLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.cur.Type == kwLIKE {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Like = expr
	} else if p.cur.Type == kwWHERE {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Where = expr
	}
	return nil
}

// parseShowFromLikeOrWhere parses optional FROM|IN db, LIKE, or WHERE for SHOW statements.
func (p *Parser) parseShowFromLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.matchFromOrIn() {
		ref, err := p.parseTableRef()
		if err != nil {
			return err
		}
		stmt.From = ref
	}
	return p.parseShowLikeOrWhere(stmt)
}

// parseShowProfileOptions parses the optional parts of SHOW PROFILE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/show-profile.html
//
//	SHOW PROFILE [type [, type] ... ]
//	    [FOR QUERY n]
//	    [LIMIT row_count [OFFSET offset]]
//
//	type: {
//	    ALL
//	  | BLOCK IO
//	  | CONTEXT SWITCHES
//	  | CPU
//	  | IPC
//	  | MEMORY
//	  | PAGE FAULTS
//	  | SOURCE
//	  | SWAPS
//	}
func (p *Parser) parseShowProfileOptions(stmt *nodes.ShowStmt) error {
	// Parse optional profile type list
	for {
		pt := p.parseShowProfileType()
		if pt == "" {
			break
		}
		stmt.ProfileTypes = append(stmt.ProfileTypes, pt)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// Optional FOR QUERY n
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		// Expect QUERY
		if p.cur.Type == kwQUERY {
			p.advance() // consume QUERY
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			stmt.ForQuery = expr
		}
	}

	// Optional LIMIT row_count [OFFSET offset]
	if _, ok := p.match(kwLIMIT); ok {
		p.advance() // skip count
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "offset") {
			p.advance() // consume OFFSET
			p.advance() // skip offset value
		}
	}

	return nil
}

// parseShowProfileType tries to parse a single SHOW PROFILE type keyword.
// Returns "" if the current token is not a valid profile type.
func (p *Parser) parseShowProfileType() string {
	if p.cur.Type == kwALL {
		p.advance()
		return "ALL"
	}
	// Handle profile types that may be keywords or identifiers
	if p.isIdentLike("cpu") {
		p.advance()
		return "CPU"
	}
	if p.isIdentLike("ipc") {
		p.advance()
		return "IPC"
	}
	if p.isIdentLike("memory") {
		p.advance()
		return "MEMORY"
	}
	if p.cur.Type == kwSOURCE || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "source")) {
		p.advance()
		return "SOURCE"
	}
	if p.isIdentLike("swaps") {
		p.advance()
		return "SWAPS"
	}
	if p.isIdentLike("block") {
		p.advance() // consume BLOCK
		if p.isIdentLike("io") {
			p.advance()
		}
		return "BLOCK_IO"
	}
	if p.isIdentLike("context") {
		p.advance() // consume CONTEXT
		if p.isIdentLike("switches") {
			p.advance()
		}
		return "CONTEXT_SWITCHES"
	}
	if p.isIdentLike("page") {
		p.advance() // consume PAGE
		if p.isIdentLike("faults") {
			p.advance()
		}
		return "PAGE_FAULTS"
	}
	return ""
}

// isIdentLike checks if the current token is an identifier (or unreserved keyword)
// matching the given name (case-insensitive).
func (p *Parser) isIdentLike(name string) bool {
	return p.cur.Type == tokIDENT && eqFold(p.cur.Str, name)
}

// expectFromOrIn expects FROM or IN keyword (both are equivalent in SHOW statements).
func (p *Parser) expectFromOrIn() (Token, error) {
	if p.cur.Type == kwFROM || p.cur.Type == kwIN {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, &ParseError{Message: "expected FROM or IN", Position: p.cur.Loc}
}

// matchFromOrIn matches FROM or IN keyword (both are equivalent in SHOW statements).
func (p *Parser) matchFromOrIn() bool {
	if p.cur.Type == kwFROM || p.cur.Type == kwIN {
		p.advance()
		return true
	}
	return false
}

// parseUseStmt parses a USE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/use.html
//
//	USE db_name
func (p *Parser) parseUseStmt() (*nodes.UseStmt, error) {
	start := p.pos()
	p.advance() // consume USE

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.UseStmt{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Database: name,
	}, nil
}

// parseExplainStmt parses an EXPLAIN or DESCRIBE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/explain.html
//
//	EXPLAIN [EXTENDED | PARTITIONS] {SELECT|INSERT|UPDATE|DELETE|REPLACE} ...
//	EXPLAIN ANALYZE SELECT ...
//	EXPLAIN FORMAT = {TRADITIONAL|JSON|TREE} {SELECT|INSERT|UPDATE|DELETE|REPLACE} ...
//	DESCRIBE tbl_name [col_name | wild]
func (p *Parser) parseExplainStmt() (*nodes.ExplainStmt, error) {
	start := p.pos()
	isDescribe := p.cur.Type == kwDESCRIBE
	p.advance() // consume EXPLAIN or DESCRIBE

	stmt := &nodes.ExplainStmt{Loc: nodes.Loc{Start: start}}

	if isDescribe {
		// DESCRIBE tbl_name [col_name]
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		// Wrap as a ShowStmt for DESCRIBE (which is equivalent to SHOW COLUMNS FROM tbl)
		showStmt := &nodes.ShowStmt{
			Loc:  nodes.Loc{Start: start},
			Type: "COLUMNS",
			From: ref,
		}
		// Optional column name
		if p.cur.Type != tokEOF && p.cur.Type != ';' {
			colExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			showStmt.Like = colExpr
		}
		showStmt.Loc.End = p.pos()
		stmt.Stmt = showStmt
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// EXPLAIN [ANALYZE] [EXTENDED] [PARTITIONS] [FORMAT = value] stmt

	// Check for ANALYZE
	if p.cur.Type == kwANALYZE {
		stmt.Analyze = true
		p.advance()
	}

	// Check for EXTENDED (deprecated in 8.0 but still parsed)
	if p.cur.Type == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	}

	// Check for PARTITIONS (deprecated in 8.0 but still parsed)
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "partitions") {
		stmt.Partitions = true
		p.advance()
	}

	// Check for FORMAT = value
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "format") {
		p.advance() // consume FORMAT
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		// Parse format value: TRADITIONAL, JSON, TREE
		formatName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Format = formatName
	}

	// EXPLAIN FOR CONNECTION connection_id
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if p.cur.Type == kwCONNECTION || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "connection")) {
			p.advance() // consume CONNECTION
		}
		if p.cur.Type == tokICONST {
			stmt.ForConnection = p.cur.Ival
			p.advance()
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Parse the explainable statement
	switch p.cur.Type {
	case kwSELECT:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = sel
	case kwINSERT:
		ins, err := p.parseInsertStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = ins
	case kwUPDATE:
		upd, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = upd
	case kwDELETE:
		del, err := p.parseDeleteStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = del
	case kwREPLACE:
		rep, err := p.parseReplaceStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = rep
	default:
		// For other tokens, try to parse as a table ref (EXPLAIN table_name)
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		showStmt := &nodes.ShowStmt{
			Loc:  nodes.Loc{Start: ref.Loc.Start, End: p.pos()},
			Type: "COLUMNS",
			From: ref,
		}
		stmt.Stmt = showStmt
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
