package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseSetStmt parses a SET statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-variable.html
//
//	SET variable = expr [, variable = expr] ...
//	variable: { user_var_name | param_name | local_var_name
//	    | {GLOBAL | @@GLOBAL.} system_var_name
//	    | {PERSIST | @@PERSIST.} system_var_name
//	    | {PERSIST_ONLY | @@PERSIST_ONLY.} system_var_name
//	    | [SESSION | @@SESSION. | @@] system_var_name }
//	SET NAMES {charset_name [COLLATE collation_name] | DEFAULT}
//	SET {CHARACTER SET | CHARSET} {charset_name | DEFAULT}
func (p *Parser) parseSetStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SET

	// Completion: after SET, offer variable/scope candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwGLOBAL)
		p.addTokenCandidate(kwSESSION)
		p.addTokenCandidate(kwCHARACTER)
		p.addRuleCandidate("variable")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.SetStmt{Loc: nodes.Loc{Start: start}}

	// Check for NAMES special form
	// SET NAMES {charset_name [COLLATE collation_name] | DEFAULT}
	if p.cur.Type == kwNAMES {
		p.advance() // consume NAMES

		// Completion: after SET NAMES, offer charset candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("charset")
			return nil, &ParseError{Message: "collecting"}
		}
		// Parse charset name or DEFAULT
		charsetLoc := p.pos()
		var charset string
		if p.cur.Type == kwDEFAULT {
			charset = "DEFAULT"
			p.advance()
		} else {
			var err error
			charset, charsetLoc, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
		}
		// Build assignment: NAMES = charset
		stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
			Loc:    nodes.Loc{Start: charsetLoc, End: p.prev.End},
			Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc, End: p.prev.End}, Column: "NAMES"},
			Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc, End: p.prev.End}, Value: charset},
		})
		// Optional COLLATE (not valid with DEFAULT, but we parse it anyway)
		if _, ok := p.match(kwCOLLATE); ok {
			collation, collLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: collLoc, End: p.prev.End},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: collLoc, End: p.prev.End}, Column: "COLLATE"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: collLoc, End: p.prev.End}, Value: collation},
			})
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Check for CHARACTER SET / CHARSET special form
	// SET {CHARACTER SET | CHARSET} {charset_name | DEFAULT}
	if p.cur.Type == kwCHARACTER || p.cur.Type == kwCHARSET {
		isCharset := p.cur.Type == kwCHARSET
		p.advance() // consume CHARACTER or CHARSET
		if isCharset || p.cur.Type == kwSET {
			if !isCharset {
				p.advance() // consume SET
			}

			// Completion: after SET CHARACTER SET, offer charset candidates.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("charset")
				return nil, &ParseError{Message: "collecting"}
			}

			charsetLoc := p.pos()
			var charset string
			if p.cur.Type == kwDEFAULT {
				charset = "DEFAULT"
				p.advance()
			} else {
				var err error
				charset, charsetLoc, err = p.parseIdentifier()
				if err != nil {
					return nil, err
				}
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: charsetLoc, End: p.prev.End},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc, End: p.prev.End}, Column: "CHARACTER SET"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc, End: p.prev.End}, Value: charset},
			})
			stmt.Loc.End = p.prev.End
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
	case kwPERSIST_ONLY:
		scope = "PERSIST_ONLY"
		p.advance()
	default:
	}

	// Completion: after SET GLOBAL/SESSION, offer variable candidates.
	if scope != "" {
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("variable")
			return nil, &ParseError{Message: "collecting"}
		}
	}

	// SET [GLOBAL|SESSION|LOCAL] TRANSACTION ... option_type has no PERSIST or PERSIST_ONLY
	// (sql_yacc.yy:17137), so `SET PERSIST TRANSACTION READ ONLY` is a 1064.
	if p.cur.Type == kwTRANSACTION {
		if scope == "PERSIST" || scope == "PERSIST_ONLY" {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseSetTransactionStmt(start, scope)
	}

	stmt.Scope = scope

	// Parse assignment list.
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetPasswordStmt parses a SET PASSWORD statement.
// SET has already been consumed; p.cur is PASSWORD.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-password.html
//
//	SET PASSWORD [FOR user] auth_option
//	    [REPLACE 'current_auth_string']
//	    [RETAIN CURRENT PASSWORD]
//
//	auth_option: {
//	    = 'auth_string'
//	  | TO RANDOM
//	}
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

	// auth_option: = 'auth_string' | TO RANDOM
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if p.cur.Type == kwRANDOM {
			p.advance() // consume RANDOM
			stmt.ToRandom = true
		}
	} else {
		// The = or := is required here: PASSWORD equal, and a bare string is a 1064
		// (measured).
		if err := p.expectEqual(); err != nil {
			return nil, err
		}

		// Password value: either a string literal or PASSWORD('string')
		if p.cur.Type == kwPASSWORD {
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
	}

	// [REPLACE 'current_auth_string']
	if p.cur.Type == kwREPLACE {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.Replace = p.cur.Str
			p.advance()
		}
	}

	// [RETAIN CURRENT PASSWORD]
	if p.cur.Type == kwRETAIN {
		p.advance() // consume RETAIN
		p.advance() // consume CURRENT
		p.advance() // consume PASSWORD
		stmt.RetainCurrentPassword = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetAssignment parses a single SET assignment: var = expr.
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
	} else if p.cur.Type == kwNEW || p.cur.Type == kwOLD {
		// Trigger-only assignment form: SET NEW.col = val / SET OLD.col = val.
		// NEW/OLD are reserved keywords that parseLvalueIdent rejects, so we
		// special-case them here as qualifier tokens.
		colStart := p.pos()
		qualTok := p.advance()
		qual := strings.ToUpper(qualTok.Str)
		if _, err := p.expect('.'); err != nil {
			return nil, err
		}
		colName, _, err := p.parseLvalueIdent()
		if err != nil {
			return nil, err
		}
		col = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: colStart, End: p.prev.End},
			Table:  qual,
			Column: colName,
		}
	} else {
		// SET variable names use lvalue_ident context: excludes ambiguous_4
		// keywords (GLOBAL, SESSION, LOCAL) which are scope modifiers, not variables.
		colStart := p.pos()
		name, _, err := p.parseLvalueIdent()
		if err != nil {
			return nil, err
		}
		col = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: colStart},
			Column: name,
		}
		// Support dot-qualified variable names: schema.var
		if p.cur.Type == '.' {
			p.advance()
			name2, _, err := p.parseLvalueIdent()
			if err != nil {
				return nil, err
			}
			col.Table = name
			col.Column = name2
		}
		col.Loc.End = p.prev.End
	}

	// Expect '=' or ':='
	if p.cur.Type == tokAssign {
		p.advance()
	} else if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// Parse value expression
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
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

	// Completion: after SHOW, offer keyword candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwTABLES, kwCOLUMNS, kwINDEX, kwDATABASES,
			kwCREATE, kwSTATUS, kwVARIABLES, kwPROCESSLIST,
			kwWARNINGS, kwERRORS, kwENGINE, kwENGINES,
			kwGLOBAL, kwSESSION, kwFULL, kwBINARY,
		} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.ShowStmt{Loc: nodes.Loc{Start: start}}

	// SHOW COUNT(*) {ERRORS | WARNINGS}. The parens are part of the
	// production (sql_yacc.yy: SHOW COUNT_SYM '(' '*' ')' WARNINGS_SYM), so
	// neither they nor the target word is optional.
	if p.cur.Type == kwCOUNT {
		p.advance() // consume COUNT
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect('*'); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case kwERRORS:
			stmt.Type = "COUNT ERRORS"
		case kwWARNINGS:
			stmt.Type = "COUNT WARNINGS"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

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
		// Optional FROM|IN db
		if p.matchFromOrIn() {
			// Completion: after SHOW TABLES FROM, offer database_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("database_ref")
				return nil, &ParseError{Message: "collecting"}
			}
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
			// Optional FROM|IN db
			if p.matchFromOrIn() {
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.From = ref
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwPROCESSLIST {
			stmt.Type = "FULL PROCESSLIST"
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	case kwCOLUMNS, kwFIELDS:
		stmt.Type = "COLUMNS"
		p.advance()
		// FROM|IN tbl
		if _, err := p.expectFromOrIn(); err != nil {
			return nil, err
		}
		// Completion: after SHOW COLUMNS FROM, offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
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
			// Completion: after SHOW CREATE TABLE, offer table_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("table_ref")
				return nil, &ParseError{Message: "collecting"}
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwDATABASE, kwSCHEMA:
			stmt.Type = "CREATE DATABASE"
			p.advance()
			// Optional IF NOT EXISTS, whose three words are one token
			// in the grammar: a bare IF, or an IF NOT, is a 1064.
			if p.cur.Type == kwIF {
				p.advance()
				if _, err := p.expect(kwNOT); err != nil {
					return nil, err
				}
				if _, err := p.expect(kwEXISTS_KW); err != nil {
					return nil, err
				}
				stmt.IfNotExists = true
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwVIEW:
			stmt.Type = "CREATE VIEW"
			p.advance()
			// Completion: after SHOW CREATE VIEW, offer view_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("view_ref")
				return nil, &ParseError{Message: "collecting"}
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwPROCEDURE:
			stmt.Type = "CREATE PROCEDURE"
			p.advance()
			// Completion: after SHOW CREATE PROCEDURE, offer procedure_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("procedure_ref")
				return nil, &ParseError{Message: "collecting"}
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwFUNCTION:
			stmt.Type = "CREATE FUNCTION"
			p.advance()
			// Completion: after SHOW CREATE FUNCTION, offer function_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("function_ref")
				return nil, &ParseError{Message: "collecting"}
			}
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
		case kwSEQUENCE:
			// SHOW CREATE SEQUENCE seq_name. MariaDB has SEQUENCE objects
			// (10.3+); MySQL 8.0 has none and answers 1064 here.
			//
			// Ref: https://mariadb.com/kb/en/show-create-sequence/
			stmt.Type = "CREATE SEQUENCE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwUSER:
			stmt.Type = "CREATE USER"
			p.advance()
			// The user is optional on MariaDB, which then reports the
			// current one; MySQL 8.0 requires it and answers 1064 on the
			// bare form.
			//
			// Ref: https://mariadb.com/kb/en/show-create-user/
			if p.cur.Type == tokEOF || p.cur.Type == ';' {
				break
			}
			// Handle CURRENT_USER / CURRENT_USER()
			if p.cur.Type == kwCURRENT_USER {
				start := p.pos()
				p.advance()
				if p.cur.Type == '(' {
					p.advance()
					p.match(')')
				}
				stmt.ForUser = &nodes.UserSpec{
					Loc:  nodes.Loc{Start: start, End: p.prev.End},
					Name: "CURRENT_USER",
				}
			} else {
				user, err := p.parseUserSpec()
				if err != nil {
					return nil, err
				}
				stmt.ForUser = user
			}
		default:
			// SHOW CREATE with nothing after it, or with an object this
			// engine has no SHOW CREATE for, is a 1064.
			return nil, p.syntaxErrorAtCur()
		}

	// MariaDB has no SHOW EXTENDED: the EXTENDED prefix MySQL 8.0 reads on
	// SHOW {COLUMNS | FIELDS | INDEX | KEYS | TABLES} is a 1064 on MariaDB
	// 11.8.8 in every one of those spellings, so kwEXTENDED is not an arm of
	// this switch.

	case kwINDEX:
		// INDEX, INDEXES and KEYS are the three spellings (ShowIndexKwd);
		// KEY on its own is not one of them and is a 1064.
		stmt.Type = "INDEX"
		p.advance()
		// FROM|IN tbl
		if _, err := p.expectFromOrIn(); err != nil {
			return nil, err
		}
		// Completion: after SHOW INDEX FROM, offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
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
		} else {
			return nil, p.syntaxErrorAtCur()
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
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			first, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if first == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
				count, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if count == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.LimitOffset = first
				stmt.LimitCount = count
			} else {
				stmt.LimitCount = first
			}
		}

	case kwERRORS:
		stmt.Type = "ERRORS"
		p.advance()
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			first, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if first == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
				count, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if count == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.LimitOffset = first
				stmt.LimitCount = count
			} else {
				stmt.LimitCount = first
			}
		}

	case kwSTORAGE:
		// SHOW STORAGE ENGINES
		p.advance() // consume STORAGE
		if _, err := p.expect(kwENGINES); err != nil {
			return nil, err
		}
		stmt.Type = "ENGINES"

	case kwENGINES:
		stmt.Type = "ENGINES"
		p.advance()

	case kwENGINE:
		// SHOW ENGINE engine_name {STATUS | MUTEX | LOGS}. LOGS is in
		// sql_yacc.yy (SHOW ENGINE_SYM engine_or_all LOGS_SYM) though the
		// manual documents only the other two, and MySQL 8.0 parses it.
		p.advance() // consume ENGINE
		engineName, nameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = &nodes.TableRef{
			Loc:  nodes.Loc{Start: nameLoc, End: p.prev.End},
			Name: engineName,
		}
		switch p.cur.Type {
		case kwSTATUS:
			stmt.Type = "ENGINE STATUS"
		case kwMUTEX:
			stmt.Type = "ENGINE MUTEX"
		case kwLOGS:
			stmt.Type = "ENGINE LOGS"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

	case kwPLUGINS:
		stmt.Type = "PLUGINS"
		p.advance()

	case kwMASTER:
		p.advance() // consume MASTER
		if p.cur.Type == kwSTATUS {
			stmt.Type = "MASTER STATUS"
			p.advance()
		} else if p.cur.Type == kwLOGS {
			// SHOW MASTER LOGS (synonym for SHOW BINARY LOGS)
			stmt.Type = "BINARY LOGS"
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	case kwSLAVE:
		p.advance() // consume SLAVE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "SLAVE STATUS"
			p.advance()
			if err := p.parseShowChannel(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwHOSTS {
			stmt.Type = "REPLICAS"
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	case kwREPLICA:
		p.advance() // consume REPLICA
		if p.cur.Type == kwHOSTS {
			// SHOW REPLICA HOSTS, MariaDB 10.5.1's alias for SHOW SLAVE
			// HOSTS. MySQL spells it SHOW REPLICAS instead, which is a 1064
			// here, so kwREPLICAS is not an arm of this switch.
			//
			// Ref: https://mariadb.com/kb/en/show-replica-hosts/
			p.advance()
			stmt.Type = "REPLICAS"
			break
		}
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "REPLICA STATUS"
		if err := p.parseShowChannel(stmt); err != nil {
			return nil, err
		}

	case kwBINARY:
		p.advance() // consume BINARY
		if _, err := p.expect(kwLOGS); err != nil {
			return nil, err
		}
		stmt.Type = "BINARY LOGS"

	case kwBINLOG:
		p.advance() // consume BINLOG
		if p.cur.Type == kwSTATUS {
			// SHOW BINLOG STATUS, MariaDB 10.5.2's name for SHOW MASTER
			// STATUS. MySQL spells the replacement SHOW BINARY LOG STATUS
			// and answers 1064 here.
			//
			// Ref: https://mariadb.com/kb/en/show-binlog-status/
			p.advance()
			stmt.Type = "BINLOG STATUS"
			break
		}
		if _, err := p.expect(kwEVENTS); err != nil {
			return nil, err
		}
		stmt.Type = "BINLOG EVENTS"
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.FromPos = expr
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			first, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if first == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
				count, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if count == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.LimitOffset = first
				stmt.LimitCount = count
			} else {
				stmt.LimitCount = first
			}
		}

	case kwTABLE:
		// STATUS and nothing else: SHOW TABLE tbl {REGIONS | DISTRIBUTIONS |
		// NEXT_ROW_ID} is TiDB's, and a 1064 here.
		p.advance() // consume TABLE
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "TABLE STATUS"
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwTRIGGERS:
		// Plural only: SHOW TRIGGER is a 1064.
		stmt.Type = "TRIGGERS"
		p.advance()
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwEVENTS:
		// Plural only: SHOW EVENT is a 1064.
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
		} else if p.cur.Type == kwCODE {
			stmt.Type = "PROCEDURE CODE"
			p.advance() // consume CODE
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	case kwFUNCTION:
		p.advance() // consume FUNCTION
		if p.cur.Type == kwSTATUS {
			stmt.Type = "FUNCTION STATUS"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwCODE {
			stmt.Type = "FUNCTION CODE"
			p.advance() // consume CODE
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	case kwOPEN:
		p.advance() // consume OPEN
		if _, err := p.expect(kwTABLES); err != nil {
			return nil, err
		}
		stmt.Type = "OPEN TABLES"
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwPRIVILEGES:
		stmt.Type = "PRIVILEGES"
		p.advance()

	case kwPROFILES:
		stmt.Type = "PROFILES"
		p.advance()

	case kwRELAYLOG:
		p.advance() // consume RELAYLOG
		if _, err := p.expect(kwEVENTS); err != nil {
			return nil, err
		}
		stmt.Type = "RELAYLOG EVENTS"
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.FromPos = expr
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			first, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if first == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
				count, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if count == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.LimitOffset = first
				stmt.LimitCount = count
			} else {
				stmt.LimitCount = first
			}
		}
		if err := p.parseShowChannel(stmt); err != nil {
			return nil, err
		}

	case kwCHARACTER:
		p.advance() // consume CHARACTER
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		stmt.Type = "CHARACTER SET"
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwCHARSET:
		// SHOW CHARSET (synonym for SHOW CHARACTER SET)
		stmt.Type = "CHARACTER SET"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwCOLLATION:
		stmt.Type = "COLLATION"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	default:
		// SHOW PROFILE [type [, type] ...] [FOR QUERY n] [LIMIT row_count [OFFSET offset]]
		if p.cur.Type == kwPROFILE {
			stmt.Type = "PROFILE"
			p.advance() // consume PROFILE
			if err := p.parseShowProfileOptions(stmt); err != nil {
				return nil, err
			}
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}

		// Handle GRANTS and PROCESSLIST as identifier-based keywords
		if p.cur.Type == kwGRANT || p.cur.Type == kwGRANTS {
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
						Loc:  nodes.Loc{Start: start, End: p.prev.End},
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
		} else if p.cur.Type == kwPROCESSLIST {
			stmt.Type = "PROCESSLIST"
			p.advance()
		} else if p.cur.Type == kwKEYS || p.cur.Type == kwINDEXES {
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
		} else if p.cur.Type == kwSCHEMAS {
			// SHOW SCHEMAS (synonym for SHOW DATABASES)
			stmt.Type = "DATABASES"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type != tokEOF && p.cur.Type != ';' {
		// A SHOW statement ends where its target ends: Split has already cut
		// on the delimiter, so anything left in this segment is a second
		// statement, and the statement loop takes one without needing a
		// separator. That turned `SHOW DATABASES TABLE t` into a SHOW plus a
		// `TABLE t` -- a 1064 on MariaDB 11.8.8, two statements here -- and did the
		// same for every other target with a statement after it. The check
		// belongs to SHOW rather than to the loop: `CREATE TABLE ta TABLE tb`
		// is one statement the loop reads as two, and closing it there needs
		// CREATE TABLE ... [AS] {TABLE t | VALUES ...} first.
		return nil, p.syntaxErrorAtCur()
	}

	if stmt.Type == "" {
		// Unreachable: every arm above either sets Type or returns an
		// error. The guard is here because the arms used to fall through
		// silently, which left an empty SHOW in front of the rest of the
		// text -- `SHOW FULL` parsed as a statement, and `SHOW ANALYZE
		// STATUS`, which is TiDB's, parsed as an empty SHOW plus an
		// ANALYZE TABLE.
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseShowChannel reads the optional FOR CHANNEL of SHOW {SLAVE | REPLICA}
// STATUS and SHOW RELAYLOG EVENTS. The channel is MySQL's
// TEXT_STRING_sys_nonewline (sql_yacc.yy, opt_channel): a string literal, so
// `FOR CHANNEL c` names no channel and is a 1064, and a FOR with no CHANNEL
// after it is one too.
func (p *Parser) parseShowChannel(stmt *nodes.ShowStmt) error {
	if p.cur.Type != kwFOR {
		return nil
	}
	p.advance() // consume FOR
	if _, err := p.expect(kwCHANNEL); err != nil {
		return err
	}
	if p.cur.Type != tokSCONST {
		return p.syntaxErrorAtCur()
	}
	stmt.Channel = p.advance().Str
	return nil
}

// parseShowLikeOrWhere parses the optional LIKE or WHERE of a SHOW statement.
// The two read different things: LIKE takes MySQL's TEXT_STRING_sys
// (sql_yacc.yy, opt_wild_or_where: LIKE_SYM TEXT_STRING_sys), one string
// literal and nothing else, while WHERE takes a whole expression. Reading an
// expression after LIKE took far more than the engine does -- `LIKE @v`,
// `LIKE concat('a','%')`, `LIKE 1`, `LIKE ('t')`, `LIKE -1`, `LIKE NULL`, and
// a second filter in `LIKE 't' LIKE 'u'` and `LIKE 't' AND 1`, every one of
// them a 1064.
func (p *Parser) parseShowLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.cur.Type == kwLIKE {
		p.advance()
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		tok := p.advance()
		stmt.Like = &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}
	} else if p.cur.Type == kwWHERE {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		if expr == nil {
			return p.syntaxErrorAtCur()
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
			if expr == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.ForQuery = expr
		}
	}

	// Optional LIMIT {row_count | offset ',' row_count | row_count OFFSET
	// offset}. SHOW PROFILE reads a whole opt_limit_clause, which is the
	// comma form as well as the OFFSET one, and the two cannot be combined.
	if _, ok := p.match(kwLIMIT); ok {
		first, err := p.parseExpr()
		if err != nil {
			return err
		}
		if first == nil {
			return p.syntaxErrorAtCur()
		}
		switch p.cur.Type {
		case ',':
			p.advance() // consume ','
			count, err := p.parseExpr()
			if err != nil {
				return err
			}
			if count == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.LimitOffset = first
			stmt.LimitCount = count
		case kwOFFSET:
			p.advance() // consume OFFSET
			off, err := p.parseExpr()
			if err != nil {
				return err
			}
			if off == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.LimitCount = first
			stmt.LimitOffset = off
		default:
			stmt.LimitCount = first
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
	// Handle profile types (now registered as keyword tokens)
	if p.cur.Type == kwCPU {
		p.advance()
		return "CPU"
	}
	if p.cur.Type == kwIPC {
		p.advance()
		return "IPC"
	}
	if p.cur.Type == kwMEMORY {
		p.advance()
		return "MEMORY"
	}
	if p.cur.Type == kwSOURCE {
		p.advance()
		return "SOURCE"
	}
	if p.cur.Type == kwSWAPS {
		p.advance()
		return "SWAPS"
	}
	if p.cur.Type == kwBLOCK {
		p.advance() // consume BLOCK
		if p.cur.Type == kwIO {
			p.advance()
		}
		return "BLOCK_IO"
	}
	if p.cur.Type == kwCONTEXT {
		p.advance() // consume CONTEXT
		if p.cur.Type == kwSWITCHES {
			p.advance()
		}
		return "CONTEXT_SWITCHES"
	}
	if p.cur.Type == kwPAGE {
		p.advance() // consume PAGE
		if p.cur.Type == kwFAULTS {
			p.advance()
		}
		return "PAGE_FAULTS"
	}
	return ""
}

// isIdentLike checks if the current token is an identifier (or unreserved keyword)
// matching the given name (case-insensitive).
func (p *Parser) isIdentLike(name string) bool {
	if p.cur.Type == tokIDENT {
		return eqFold(p.cur.Str, name)
	}
	// Also match non-reserved keyword tokens by their string value
	if p.cur.Type >= 700 && !isReserved(p.cur.Type) {
		return eqFold(p.cur.Str, name)
	}
	return false
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

	// Completion: after USE, offer database_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.UseStmt{
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
		Database: name,
	}, nil
}

// parseExplainStmt parses an EXPLAIN, DESCRIBE or DESC statement.
//
// Ref: https://mariadb.com/kb/en/explain/; describe_command and
// opt_describe_column in MariaDB 11.8 sql/sql_yacc.yy. Measured against
// mariadb:11.8.8.
//
//	{EXPLAIN | DESCRIBE | DESC} tbl_name [col_name | wild]
//	{EXPLAIN | DESCRIBE | DESC} [EXTENDED | PARTITIONS | FORMAT = format] explainable_stmt
//	{EXPLAIN | DESCRIBE | DESC} [FORMAT = format] FOR CONNECTION connection_id
//	format: ident_or_text
//	explainable_stmt: {SELECT | VALUES | DELETE | INSERT | REPLACE | UPDATE}
//
// The modifier list is MariaDB's, and it is not MySQL 8.0's. EXTENDED and
// PARTITIONS are still here, where MySQL 8.0.46 and TiDB each answer 1064 for
// both, and they are one choice rather than a list: EXPLAIN EXTENDED
// PARTITIONS and EXPLAIN EXTENDED FORMAT = JSON are each a 1064, and neither
// reaches FOR CONNECTION.
//
// There is no EXPLAIN ANALYZE: MariaDB spells that ANALYZE [FORMAT = JSON]
// stmt, a statement of its own, so EXPLAIN ANALYZE SELECT 1 is a 1064 here
// where MySQL and TiDB both take it. That spelling is additive MariaDB surface
// this parser does not read yet, not a form to port from MySQL.
//
// No DDL is explainable, and TABLE is no statement here: EXPLAIN TABLE t and
// EXPLAIN ALTER TABLE t ADD INDEX (a) are both syntax errors. A modifier rules
// the table form out, so FORMAT = JSON t is a 1064.
//
// DESC and DESCRIBE read only the table form here. MariaDB takes them as full
// EXPLAIN synonyms -- DESC FORMAT = JSON SELECT 1 parses on 11.8.8 -- which is
// a gap this file records rather than closes.
func (p *Parser) parseExplainStmt() (*nodes.ExplainStmt, error) {
	start := p.pos()
	isDescribe := p.cur.Type == kwDESCRIBE || p.cur.Type == kwDESC
	p.advance() // consume EXPLAIN, DESCRIBE, or DESC

	stmt := &nodes.ExplainStmt{Loc: nodes.Loc{Start: start}}

	if isDescribe {
		// Completion: after DESCRIBE/DESC, offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		// DESCRIBE tbl_name [col_name | wild]
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
		if err := p.parseExplainColumn(showStmt); err != nil {
			return nil, err
		}
		showStmt.Loc.End = p.prev.End
		stmt.Stmt = showStmt
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// One modifier at most: EXTENDED, PARTITIONS or FORMAT = format. The
	// format name is ident_or_text -- an identifier, a non-reserved keyword, a
	// backquoted identifier, or a string literal in either quote -- and
	// MariaDB checks the name itself after parsing (1791), so FORMAT = tree
	// and FORMAT = 'tree' both parse while FORMAT = 1 and the reserved
	// FORMAT = binary are each a 1064.
	//
	// EXTENDED and PARTITIONS are not reserved, so each is also a table name:
	// EXPLAIN EXTENDED t reads as the table form, describing a table called
	// `extended`, and parses on 11.8.8. One token of lookahead settles which
	// arm the word starts, as the engine's own grammar does.
	hasFormat := false
	switch p.cur.Type {
	case kwEXTENDED:
		if !p.explainModifierAhead() {
			break
		}
		stmt.Extended = true
		p.advance()
	case kwPARTITIONS:
		if !p.explainModifierAhead() {
			break
		}
		stmt.Partitions = true
		p.advance()
	case kwFORMAT:
		hasFormat = true
		p.advance() // consume FORMAT
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		name, _, err := p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
		stmt.Format = name
	}

	// EXPLAIN [FORMAT = format] FOR CONNECTION connection_id, which EXTENDED
	// and PARTITIONS do not reach. CONNECTION and the id are both required:
	// FOR 1 and FOR CONNECTION are each a 1064.
	//
	// The id is read as an unsigned integer, which is the documented form.
	// MariaDB 11.8.8 also parses a signed number, a string literal and a user
	// variable there -- FOR CONNECTION -1, FOR CONNECTION 'x' and
	// FOR CONNECTION @a all reach the server, where MySQL and TiDB refuse all
	// three -- and none of them is a connection id any layer above this parser
	// could use, so they stay rejected here and the divergence is pinned in
	// explain_forms_test.go instead.
	if p.cur.Type == kwFOR && !stmt.Extended && !stmt.Partitions {
		p.advance() // consume FOR
		if _, err := p.expect(kwCONNECTION); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokICONST)
		if err != nil {
			return nil, err
		}
		stmt.ForConnection = tok.Ival
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Completion: after EXPLAIN [options], offer explainable statement keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwEXTENDED, kwPARTITIONS, kwFORMAT, kwFOR,
			kwSELECT, kwWITH, kwINSERT, kwUPDATE, kwDELETE, kwREPLACE, kwVALUES,
		} {
			p.addTokenCandidate(t)
		}
		if !hasFormat {
			// The table form is still reachable: EXPLAIN t is as valid as
			// EXPLAIN SELECT, and a cursor right after EXTENDED or
			// PARTITIONS is the one place that word could still be the
			// table's name. Only a FORMAT rules the form out.
			p.addRuleCandidate("table_ref")
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// Parse the explainable statement
	switch p.cur.Type {
	case kwSELECT, kwVALUES, kwWITH, '(':
		// A query expression, as at statement level (MariaDB's
		// explainable_command is select, insert, replace, update, or delete).
		sel, err := p.parseQueryStatement()
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
		// EXPLAIN tbl_name [col_name | wild], which no modifier stands in
		// front of.
		// hasFormat rather than Format != "": the name may be an empty
		// string, and FORMAT = '' is a modifier all the same.
		if stmt.Extended || stmt.Partitions || hasFormat {
			return nil, p.syntaxErrorAtCur()
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		showStmt := &nodes.ShowStmt{
			Loc:  nodes.Loc{Start: ref.Loc.Start, End: p.prev.End},
			Type: "COLUMNS",
			From: ref,
		}
		if err := p.parseExplainColumn(showStmt); err != nil {
			return nil, err
		}
		showStmt.Loc.End = p.prev.End
		stmt.Stmt = showStmt
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// startsExplainableStmt reports whether a token can begin the statement
// EXPLAIN explains. It tells EXPLAIN's one-token-lookahead arms apart: the
// word after EXTENDED or PARTITIONS decides whether that word was a modifier
// or the table name of the table form.
func startsExplainableStmt(tokenType int) bool {
	switch tokenType {
	case kwSELECT, kwVALUES, kwWITH, '(', kwINSERT, kwUPDATE, kwDELETE, kwREPLACE:
		return true
	}
	return false
}

// explainModifierAhead reports whether EXTENDED or PARTITIONS at the current
// token is the modifier rather than the table name of the table form, which
// what follows it decides.
//
// Completion asks from the other side. With the cursor right after the word,
// as in "EXPLAIN EXTENDED |", nothing follows it yet: strict parsing reads the
// table form, which is what the server does with EXPLAIN EXTENDED on its own,
// while the caller is about to type the statement the modifier takes. So a
// cursor there picks the modifier, and the candidates that arm collects name
// the table form too.
func (p *Parser) explainModifierAhead() bool {
	next := p.peekNext()
	if startsExplainableStmt(next.Type) {
		return true
	}
	return p.completing && next.Loc >= p.cursorOff
}

// parseExplainColumn parses the column name or wild pattern EXPLAIN and
// DESCRIBE take after a table name, into the ShowStmt the table form lowers
// to. MariaDB's slot is opt_describe_column, an identifier or a string, and
// nothing else: EXPLAIN t a, EXPLAIN t `a` and EXPLAIN t 'a%' all parse on
// 11.8.8, while a qualified name, an integer and a reserved word are each a
// 1064 -- EXPLAIN t a.b and EXPLAIN t 123 are syntax errors, where TiDB takes
// the qualified name and refuses the string.
//
// It is held in ShowStmt.Like, which is where the EXPLAIN table form has always
// put it; ShowStmt has no column field of its own.
func (p *Parser) parseExplainColumn(showStmt *nodes.ShowStmt) error {
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil
	}
	if p.cur.Type == tokSCONST {
		tok := p.advance()
		showStmt.Like = &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: tok.Str}
		return nil
	}
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	showStmt.Like = &nodes.ColumnRef{Loc: nodes.Loc{Start: start, End: p.prev.End}, Column: name}
	return nil
}
