package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseSetStmt parses a SET statement.
//
// Ref: var_and_val, sys_var_and_val, scope_or_scope_alias, to_or_eq and var_name in
// OceanBase CE v4.3.5_CE_BP6 src/sql/parser/sql_parser_mysql_mode.y (lines 18651, 18698,
// 18707, 18714, 24308).
//
//	SET assignment [, assignment] ...
//	assignment: { [GLOBAL | SESSION] name {= | := | TO} expr
//	    | @@[{GLOBAL | SESSION | LOCAL}.]name {= | := | TO} expr
//	    | @user_var {= | := | TO} expr }
//	SET NAMES {charset_name [COLLATE collation_name] | DEFAULT}
//	SET {CHARACTER SET | CHARSET} {charset_name | DEFAULT}
//
// Every assignment may write a scope, and a scope names the system variable after it, so
// `SET SESSION @a = 1` is a 1064. A name takes no dot, so `SET a.b = 1` is a 1064, but for
// NEW and OLD in a trigger, which the server tells apart with a flag this parser lacks.
// LOCAL, PERSIST and PERSIST_ONLY are no scopes of an assignment: they read as ordinary
// names, so `SET LOCAL a = 1` is a 1064, while set_transaction_stmt takes LOCAL. The words
// that start a form below are also variable names; see nextEndsVariableName.
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

	// Whether the first word names a variable, whichever form it could start.
	named := p.nextEndsVariableName()

	// Check for SET DEFAULT ROLE
	if p.cur.Type == kwDEFAULT && !named {
		p.advance() // consume DEFAULT
		if p.cur.Type == kwROLE {
			p.advance() // consume ROLE
			return p.parseSetDefaultRoleStmt(start)
		}
		return nil, &ParseError{Message: "expected ROLE after SET DEFAULT", Position: p.cur.Loc}
	}

	// Check for SET ROLE
	if p.cur.Type == kwROLE && !named {
		p.advance() // consume ROLE
		return p.parseSetRoleStmt(start)
	}

	// Check for SET PASSWORD
	if p.cur.Type == kwPASSWORD && !named {
		return p.parseSetPasswordStmt(start)
	}

	// The scope the first assignment writes, or SET TRANSACTION's.
	assignStart := p.pos()
	scope := ""
	if !named {
		scope = p.parseSetScope()
	}
	// set_transaction_stmt also takes LOCAL (sql_parser_mysql_mode.y:23994), which no
	// assignment does.
	if scope == "" && p.cur.Type == kwLOCAL && p.peekNext().Type == kwTRANSACTION {
		scope = "LOCAL"
		p.advance()
	}

	// Completion: after SET GLOBAL/SESSION, offer variable candidates.
	if scope != "" {
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("variable")
			return nil, &ParseError{Message: "collecting"}
		}
	}

	// SET [GLOBAL|SESSION|LOCAL] TRANSACTION ... Those are the only scopes
	// set_transaction_stmt names (sql_parser_mysql_mode.y:23975), so `SET PERSIST
	// TRANSACTION READ ONLY` is a 1064, PERSIST reading as a name. TRANSACTION is a name
	// when to_or_eq follows, but only with no scope: the grammar commits to
	// set_transaction_stmt after one, so `SET SESSION TRANSACTION = 1` is a 1064.
	if p.cur.Type == kwTRANSACTION && (scope != "" || !p.nextEndsVariableName()) {
		return p.parseSetTransactionStmt(start, scope)
	}

	stmt.Scope = scope

	// Parse assignment list. Each assignment writes its own scope.
	for {
		if err := p.parseSetItem(stmt, assignStart, scope); err != nil {
			return nil, err
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','

		// Completion: where an assignment starts, offer its scope words and variables.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwGLOBAL)
			p.addTokenCandidate(kwSESSION)
			p.addRuleCandidate("variable")
			return nil, &ParseError{Message: "collecting"}
		}
		assignStart = p.pos()
		scope = ""
		if !p.nextEndsVariableName() {
			scope = p.parseSetScope()
		}
		if scope != "" {
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("variable")
				return nil, &ParseError{Message: "collecting"}
			}
		}
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
		user, err := p.parseUserSpec(false)
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
		// consume '='
		p.match('=')

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

// parseSetItem reads one assignment of the list into stmt, after its scope. NAMES and the
// charset forms are assignments, so they may sit anywhere in the list, and a written scope
// takes a name and nothing else, so `SET GLOBAL NAMES utf8` is a 1064.
func (p *Parser) parseSetItem(stmt *nodes.SetStmt, start int, scope string) error {
	if scope == "" {
		switch {
		// NAMES is also a name here: `SET NAMES = 1` parses.
		case p.cur.Type == kwNAMES && !p.nextIsAssignOp():
			return p.parseSetNames(stmt)
		case p.cur.Type == kwCHARSET && !p.nextIsAssignOp(),
			(p.cur.Type == kwCHARACTER || p.cur.Type == kwCHAR) && p.peekNext().Type == kwSET:
			return p.parseSetCharset(stmt)
		}
	}
	asgn, err := p.parseSetVariableAssignment(start, scope)
	if err != nil {
		return err
	}
	stmt.Assignments = append(stmt.Assignments, asgn)
	return nil
}

// nextIsAssignOp reports whether to_or_eq follows p.cur: `=`, `:=` or TO.
func (p *Parser) nextIsAssignOp() bool {
	next := p.peekNext().Type
	return next == '=' || next == tokAssign || next == kwTO
}

// parseSetNames reads NAMES {charset_name | DEFAULT} [COLLATE collation_name] as one assignment
// for the charset and one for the collation. A collation is a string or a name, so `SET NAMES a
// COLLATE BINARY` and `SET NAMES a COLLATE DEFAULT` are 1064s. p.cur is NAMES.
func (p *Parser) parseSetNames(stmt *nodes.SetStmt) error {
	p.advance() // consume NAMES

	// Completion: after SET NAMES, offer charset candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("charset")
		return &ParseError{Message: "collecting"}
	}

	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		appendSetWord(stmt, "NAMES", "DEFAULT", tok.Loc, p.prev.End)
	} else {
		charset, loc, err := p.parseSetCharsetName()
		if err != nil {
			return err
		}
		appendSetWord(stmt, "NAMES", charset, loc, p.prev.End)
	}

	if _, ok := p.match(kwCOLLATE); !ok {
		return nil
	}
	if p.isVariableRef() {
		return p.syntaxErrorAtCur()
	}
	collation, loc, err := p.parseIdentOrText()
	if err != nil {
		return err
	}
	appendSetWord(stmt, "COLLATE", collation, loc, p.prev.End)
	return nil
}

// parseSetCharset reads {CHARACTER SET | CHAR SET | CHARSET} {charset_name | DEFAULT}. p.cur is
// its first word, and the caller has seen that CHARACTER and CHAR have their SET.
func (p *Parser) parseSetCharset(stmt *nodes.SetStmt) error {
	if p.cur.Type != kwCHARSET {
		p.advance() // consume CHARACTER or CHAR
	}
	p.advance() // consume SET or CHARSET

	// Completion: after SET CHARACTER SET, offer charset candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("charset")
		return &ParseError{Message: "collecting"}
	}

	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		appendSetWord(stmt, "CHARACTER SET", "DEFAULT", tok.Loc, p.prev.End)
		return nil
	}
	charset, loc, err := p.parseSetCharsetName()
	if err != nil {
		return err
	}
	appendSetWord(stmt, "CHARACTER SET", charset, loc, p.prev.End)
	return nil
}

// parseSetCharsetName reads a charset or collation name: a string, a name, or BINARY. The
// lexer hands a user variable over as an identifier, and `SET NAMES @a` is a 1064.
func (p *Parser) parseSetCharsetName() (string, int, error) {
	if p.cur.Type == kwBINARY {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.isVariableRef() {
		return "", 0, p.syntaxErrorAtCur()
	}
	return p.parseIdentOrText()
}

// appendSetWord adds the assignment a NAMES or charset word stands for.
func appendSetWord(stmt *nodes.SetStmt, column, value string, start, end int) {
	loc := nodes.Loc{Start: start, End: end}
	stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
		Loc:    loc,
		Column: &nodes.ColumnRef{Loc: loc, Column: column},
		Value:  &nodes.StringLit{Loc: loc, Value: value},
	})
}

// parseSetValue reads the value of an assignment to col. A name and @@name take
// a set_expr: ON, BINARY, SYSTEM, DEFAULT or an expression. ON is a true, as
// the review walk reads it, and the other words are the value as text. @name takes an
// expression alone, so `SET @a = ON` and `SET @a = DEFAULT` are 1064s.
//
// BINARY and SYSTEM are the value only where the assignment ends, because `SET a =
// BINARY 'x'` is an expression. ON and a bare DEFAULT are the value and nothing more, so `SET a
// = ON + 1` and `SET a = DEFAULT + 1` are 1064s, where parseExpr would go on to read the sum.
// DEFAULT(b) is an expression.
func (p *Parser) parseSetValue(col *nodes.ColumnRef) (nodes.ExprNode, error) {
	next := p.peekNext().Type
	bareDefault := p.cur.Type == kwDEFAULT && next != '('
	if strings.HasPrefix(col.Column, "@") && !strings.HasPrefix(col.Column, "@@") {
		if bareDefault {
			return nil, p.syntaxErrorAtCur()
		}
	} else {
		ends := next == ',' || next == ';' || next == tokEOF
		switch {
		case p.cur.Type == kwON:
			tok := p.advance()
			return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: true}, nil
		case (p.cur.Type == kwBINARY || p.cur.Type == kwSYSTEM) && ends:
			tok := p.advance()
			return &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: strings.ToUpper(tok.Str)}, nil
		case bareDefault:
			return p.parseDefaultExpr()
		}
	}
	// An expr holds no bare DEFAULT, so `SET a = (DEFAULT)` and `SET a = 1 + DEFAULT` are 1064s.
	outer := p.noBareDefault
	p.noBareDefault = true
	val, err := p.parseExpr()
	p.noBareDefault = outer
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return val, nil
}

// nextEndsVariableName reports whether the word at p.cur is a variable's name rather than
// the keyword of a SET form. var_name takes the scope words, ROLE, CHARSET, NAMES and
// TRANSACTION as unreserved keywords, so to_or_eq after one ends a name.
//
// A next token at or past the completion cursor is text the user has not reached, so
// `SET TRANSACTION <cursor> = 1` still offers what follows the keyword.
func (p *Parser) nextEndsVariableName() bool {
	next := p.peekNext()
	if p.completing && next.Loc >= p.cursorOff {
		return false
	}
	switch p.cur.Type {
	case kwGLOBAL, kwSESSION, kwLOCAL, kwPERSIST, kwPERSIST_ONLY, kwROLE, kwCHARSET, kwNAMES, kwTRANSACTION:
		return next.Type == '=' || next.Type == tokAssign || next.Type == kwTO
	}
	return false
}

// parseSetScope consumes the scope word at p.cur, if there is one.
func (p *Parser) parseSetScope() string {
	scope := ""
	switch p.cur.Type {
	case kwGLOBAL:
		scope = "GLOBAL"
	case kwSESSION:
		scope = "SESSION"
	default:
		return ""
	}
	p.advance()
	return scope
}

// parseSetVariableAssignment parses one assignment of a SET statement, after the scope it
// wrote, if any. start is where the assignment began.
func (p *Parser) parseSetVariableAssignment(start int, scope string) (*nodes.Assignment, error) {
	var col *nodes.ColumnRef

	if p.isVariableRef() {
		// A scope names the system variable after it, never a variable reference.
		if scope != "" {
			return nil, p.syntaxErrorAtCur()
		}
		// @@ takes a scope word only before a dot.
		atScope := false
		if p.cur.Type == tokAt2 {
			switch p.peekNext().Type {
			case kwGLOBAL, kwSESSION, kwLOCAL:
				atScope = true
			}
		}
		vref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		if atScope && vref.Scope == "" {
			return nil, &ParseError{Message: "unexpected scope", Position: vref.Loc.Start}
		}
		// var_name takes no dot, and SYSTEM_VARIABLE's scopes are GLOBAL_ALIAS and
		// SESSION_ALIAS, which @@local. also lexes as.
		if vref.System && (strings.Contains(vref.Name, ".") || vref.Scope == "PERSIST" || vref.Scope == "PERSIST_ONLY") {
			return nil, &ParseError{Message: "unexpected variable name", Position: vref.Loc.Start}
		}
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
	} else if (p.cur.Type == kwNEW || p.cur.Type == kwOLD) && p.peekNext().Type == '.' && scope == "" {
		// Trigger-only assignment form: SET NEW.col = val / SET OLD.col = val.
		colStart := p.pos()
		qual := strings.ToUpper(p.advance().Str)
		p.advance() // consume '.'
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
		colStart := p.pos()
		var name string
		// var_name: NAME_OB or unreserved_keyword_normal. PASSWORD is a form of its own
		// with no scope, so a later `PASSWORD = 1` is a 1064.
		if scope == "" && p.cur.Type == kwPASSWORD {
			return nil, p.syntaxErrorAtCur()
		}
		var err error
		name, _, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		col = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: colStart},
			Column: name,
		}
		col.Loc.End = p.prev.End
	}

	// to_or_eq: `=`, `:=` or TO.
	if p.cur.Type == tokAssign || p.cur.Type == kwTO {
		p.advance()
	} else if _, err := p.expect('='); err != nil {
		return nil, err
	}

	val, err := p.parseSetValue(col)
	if err != nil {
		return nil, err
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Scope:  scope,
		Column: col,
		Value:  val,
	}, nil
}

// parseSetAssignment parses a single var = expr of INSTALL COMPONENT ... SET.
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
		case kwUSER:
			stmt.Type = "CREATE USER"
			p.advance()
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
				user, err := p.parseUserSpec(false)
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

	case kwEXTENDED:
		p.advance() // consume EXTENDED
		if p.cur.Type == kwINDEX || p.cur.Type == kwKEY || p.cur.Type == kwKEYS || p.cur.Type == kwINDEXES {
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
			if p.cur.Type == kwCOLUMNS || p.cur.Type == kwFIELDS {
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
			} else if p.cur.Type == kwTABLES {
				// SHOW EXTENDED FULL TABLES ...
				stmt.Type = "EXTENDED FULL TABLES"
				p.advance()
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
			}
		} else if p.cur.Type == kwCOLUMNS || p.cur.Type == kwFIELDS {
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
		} else if p.cur.Type == kwTABLES {
			// SHOW EXTENDED TABLES ...
			stmt.Type = "EXTENDED TABLES"
			p.advance()
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
		} else {
			return nil, p.syntaxErrorAtCur()
		}

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
		// SHOW ENGINE engine_name {STATUS | MUTEX}. MySQL 8.0 reads a third
		// arm, LOGS (sql_yacc.yy, SHOW ENGINE_SYM engine_or_all LOGS_SYM),
		// which OceanBase CE 4.3.5.6 answers 1064 on, so it is not here.
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

		// SHOW REPLICAS
		if p.cur.Type == kwREPLICAS {
			stmt.Type = "REPLICAS"
			p.advance()
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
					user, err := p.parseUserSpec(false)
					if err != nil {
						return nil, err
					}
					stmt.ForUser = user
				}
				// Optional USING role [, role] ...
				if _, ok := p.match(kwUSING); ok {
					for {
						role, err := p.parseUserSpec(false)
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
		// `TABLE t` -- a 1064 on OceanBase CE 4.3.5.6, two statements here -- and did the
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

	// Optional LIMIT row_count [OFFSET offset]. MySQL 8.0 reads a whole
	// opt_limit_clause here, the `LIMIT offset, row_count` form included;
	// OceanBase CE 4.3.5.6 reads no SHOW PROFILE argument at all and
	// answers 1064 on every one of them, so the comma form is not ported.
	if _, ok := p.match(kwLIMIT); ok {
		count, err := p.parseExpr()
		if err != nil {
			return err
		}
		if count == nil {
			return p.syntaxErrorAtCur()
		}
		stmt.LimitCount = count
		if p.cur.Type == kwOFFSET {
			p.advance() // consume OFFSET
			off, err := p.parseExpr()
			if err != nil {
				return err
			}
			if off == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.LimitOffset = off
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

// explainFormatNames are the bare names OceanBase reads after FORMAT =. The
// list is closed and the name is never quoted: FORMAT = 'traditional',
// "FORMAT = `TRADITIONAL`", FORMAT = ROW and FORMAT = nosuchformat are each a
// 1064 on 5.7.25-OceanBase_CE-v4.3.5.6, where MySQL 8.0.46 takes any
// ident_or_text and checks the name after parsing (1791). OceanBase's own extra
// output modes are modifiers rather than format names -- EXPLAIN BASIC and
// EXPLAIN OUTLINE parse, FORMAT = BASIC does not.
var explainFormatNames = map[string]bool{
	"TRADITIONAL": true,
	"JSON":        true,
}

// parseExplainFormatName parses the name after EXPLAIN's FORMAT =. Only the
// bare TRADITIONAL and JSON reach it, and a backquoted name does not: a
// backquoted identifier is no keyword in OceanBase's grammar either.
func (p *Parser) parseExplainFormatName() (string, error) {
	if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
		if name := strings.ToUpper(p.cur.Str); explainFormatNames[name] && !p.backtickQuotedAt(p.cur.Loc) {
			return p.advance().Str, nil
		}
	}
	if p.cur.Type == tokEOF {
		return "", p.syntaxErrorAtCur()
	}
	return "", &ParseError{
		Message:  "expected TRADITIONAL or JSON as the EXPLAIN format name",
		Position: p.cur.Loc,
	}
}

// parseExplainStmt parses an EXPLAIN, DESCRIBE or DESC statement.
//
// Ref: https://en.oceanbase.com/docs/common-oceanbase-database-10000000002017657
// (EXPLAIN); measured against 5.7.25-OceanBase_CE-v4.3.5.6.
//
//	{EXPLAIN | DESCRIBE | DESC} tbl_name [col_name | wild]
//	{EXPLAIN | DESCRIBE | DESC} [EXTENDED | PARTITIONS | FORMAT = format] explainable_stmt
//	format: {TRADITIONAL | JSON}
//	explainable_stmt: {SELECT | VALUES | DELETE | INSERT | REPLACE | UPDATE}
//
// The modifier list is OceanBase's, and it is neither MySQL 8.0's nor the one
// this parser inherited. EXTENDED and PARTITIONS are still here, as they are in
// MySQL 5.7, whose protocol version this server reports, while MySQL 8.0.46 and
// TiDB each answer 1064 for both. They are one choice rather than a list:
// EXPLAIN EXTENDED PARTITIONS, EXPLAIN EXTENDED FORMAT = TRADITIONAL and
// EXPLAIN FORMAT = TRADITIONAL EXTENDED are each a 1064.
//
// There is no EXPLAIN ANALYZE -- EXPLAIN ANALYZE SELECT 1 is a 1064, where
// MySQL and TiDB both take it -- and no FOR CONNECTION, so EXPLAIN FOR
// CONNECTION 1 is a 1064 too. Neither is any DDL explainable, and TABLE is no
// statement here: EXPLAIN TABLE t and EXPLAIN ALTER TABLE t ADD INDEX (a) are
// both syntax errors. A modifier rules the table form out, so
// FORMAT = TRADITIONAL t is a 1064.
//
// OceanBase has two further modifiers this parser does not read yet, EXPLAIN
// BASIC and EXPLAIN OUTLINE; both parse on the server and are recorded as
// pending in the OceanBase grammar inventory.
//
// DESC and DESCRIBE read only the table form here. OceanBase takes them as full
// EXPLAIN synonyms -- DESC FORMAT = TRADITIONAL SELECT 1 parses -- which is a
// gap this file records rather than closes.
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

	// One modifier at most: EXTENDED, PARTITIONS or FORMAT = format.
	//
	// EXTENDED and PARTITIONS are not reserved, so each is also a table name:
	// EXPLAIN EXTENDED t reads as the table form, describing a table called
	// `extended`, and parses on the server. One token of lookahead settles
	// which arm the word starts, as the engine's own grammar does.
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
		name, err := p.parseExplainFormatName()
		if err != nil {
			return nil, err
		}
		stmt.Format = name
	}

	// Completion: after EXPLAIN [options], offer explainable statement keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwEXTENDED, kwPARTITIONS, kwFORMAT,
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
		// A query expression, which parseQueryStatement reads whole: a SELECT
		// leaf, a VALUES primary, a WITH clause in front of one, a parenthesized
		// body, or a set operation over them.
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
// to. The slot is MySQL's opt_describe_column, an identifier or a string, and
// nothing else: EXPLAIN t a, EXPLAIN t `a` and EXPLAIN t 'a%' all parse on
// 5.7.25-OceanBase_CE-v4.3.5.6, while a qualified name, an integer and a
// reserved word are each a 1064.
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
