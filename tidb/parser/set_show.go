package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseSetStmt parses a SET statement.
//
// The rules are SetStmt, parser.y:10745, and VariableAssignment, parser.y:10942. GLOBAL,
// SESSION and LOCAL are the only scope words VariableAssignment names, so `SET PERSIST a
// = 1` is a 1064 on TiDB v8.5.5, which reads PERSIST as an ordinary variable name and then
// finds a second word where it needs `=`.
//
//	SET assignment [, assignment] ...
//	assignment: { [GLOBAL | SESSION | LOCAL] name {= | :=} set_expr
//	    | @@[scope.]name {= | :=} set_expr
//	    | @user_var {= | :=} expr
//	    | NAMES {charset_name [COLLATE {collation_name | DEFAULT}] | DEFAULT}
//	    | {CHARACTER SET | CHAR SET | CHARSET} {charset_name | DEFAULT} }
//
// NAMES and the charset words are assignments upstream, so they may sit anywhere in the
// list. Every other form below is a whole statement. The unreserved words that start a
// form are also variable names; see nextEndsVariableName.
func (p *Parser) parseSetStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SET

	if err := p.collectSetAssignment(); err != nil {
		return nil, err
	}

	stmt := &nodes.SetStmt{Loc: nodes.Loc{Start: start}}

	// Whether the first word names a variable, whichever special form it could start.
	named := p.nextEndsVariableName()

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
	if p.cur.Type == kwROLE && !named {
		p.advance() // consume ROLE
		return p.parseSetRoleStmt(start)
	}

	// Check for SET PASSWORD. `=` belongs to this form, so `SET PASSWORD = 1` is a 1064
	// where `SET PASSWORD := 1` and `SET password.a = 1` assign a variable.
	if p.cur.Type == kwPASSWORD && (!named || p.peekNext().Type == '=') {
		return p.parseSetPasswordStmt(start)
	}

	// Check for SET RESOURCE GROUP
	if p.cur.Type == kwRESOURCE && !named {
		return p.parseSetResourceGroupStmt(start)
	}

	// SET BINDING {ENABLED | DISABLED} FOR ..., TiDB's SQL plan binding
	// switch. BINDING is unreserved and `SET binding = 1` is an ordinary
	// variable assignment on v8.5.5, so the status word after it is what
	// picks this statement out.
	if p.cur.Type == kwBINDING && (p.peekNext().Type == kwENABLED || p.peekNext().Type == kwDISABLED) {
		return p.parseSetBindingStmt(start)
	}

	scope, scopeStart, err := p.parseSetScope()
	if err != nil {
		return nil, err
	}

	// SET [GLOBAL|SESSION] TRANSACTION ... Those are the only scopes SetStmt names
	// it with (parser.y:10758-10770), so `SET LOCAL TRANSACTION READ ONLY` is a 1064.
	// TRANSACTION is also an unreserved word, so `SET TRANSACTION = 1` assigns a variable.
	if p.cur.Type == kwTRANSACTION && (scope == "" || scope == "GLOBAL" || scope == "SESSION") && !p.nextEndsVariableName() {
		return p.parseSetTransactionStmt(start, scope)
	}

	stmt.Scope = scope

	// Parse assignment list
	for {
		if err := p.parseSetItem(stmt, scope, scopeStart); err != nil {
			return nil, err
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','

		if err := p.collectSetAssignment(); err != nil {
			return nil, err
		}
		scope, scopeStart, err = p.parseSetScope()
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// collectSetAssignment offers what can start an assignment, under completion.
func (p *Parser) collectSetAssignment() error {
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwGLOBAL)
		p.addTokenCandidate(kwSESSION)
		p.addTokenCandidate(kwCHARACTER)
		p.addRuleCandidate("variable")
		return &ParseError{Message: "collecting"}
	}
	return nil
}

// parseSetScope reads the scope word of one assignment, if it has one, and returns the
// scope with where the assignment starts. Upstream gives every assignment its own scope,
// and one does not carry to the next: in `SET GLOBAL a = 1, b = 2` only a is global.
func (p *Parser) parseSetScope() (string, int, error) {
	start := p.pos()
	scope := ""
	if !p.nextEndsVariableName() {
		switch p.cur.Type {
		case kwGLOBAL:
			scope = "GLOBAL"
		case kwSESSION:
			scope = "SESSION"
		case kwLOCAL:
			scope = "LOCAL"
		default:
		}
	}
	if scope == "" {
		return "", start, nil
	}
	p.advance()

	// Completion: after SET GLOBAL/SESSION, offer variable candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("variable")
		return "", 0, &ParseError{Message: "collecting"}
	}
	return scope, start, nil
}

// parseSetItem reads one assignment of the list into stmt, after its scope. A scope takes
// a VariableName and nothing else, parser.y:10947-10958, so `SET GLOBAL @a = 1` and `SET
// GLOBAL NAMES utf8` are 1064s.
func (p *Parser) parseSetItem(stmt *nodes.SetStmt, scope string, start int) error {
	if scope != "" {
		asgn, err := p.parseSetVariableAssignment(start)
		if err != nil {
			return err
		}
		asgn.Scope = scope
		stmt.Assignments = append(stmt.Assignments, asgn)
		return nil
	}

	named := p.nextEndsVariableName()
	switch {
	case p.cur.Type == kwNAMES && !named:
		return p.parseSetNames(stmt)
	case p.cur.Type == kwCHARSET && !named,
		(p.cur.Type == kwCHARACTER || p.cur.Type == kwCHAR) && p.peekNext().Type == kwSET:
		return p.parseSetCharset(stmt)
	}

	asgn, err := p.parseSetAssignment()
	if err != nil {
		return err
	}
	stmt.Assignments = append(stmt.Assignments, asgn)
	return nil
}

// parseSetNames reads NAMES {CharsetName [COLLATE {DEFAULT | StringName}] | DEFAULT},
// parser.y:10981-11003, as one assignment for the charset and one for the collation. p.cur
// is NAMES.
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
		return nil
	}
	charset, loc, err := p.parseSetCharsetName()
	if err != nil {
		return err
	}
	appendSetWord(stmt, "NAMES", charset, loc, p.prev.End)

	if _, ok := p.match(kwCOLLATE); !ok {
		return nil
	}
	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		appendSetWord(stmt, "COLLATE", "DEFAULT", tok.Loc, p.prev.End)
		return nil
	}
	collation, loc, err := p.parseSetStringName()
	if err != nil {
		return err
	}
	appendSetWord(stmt, "COLLATE", collation, loc, p.prev.End)
	return nil
}

// parseSetCharset reads CharsetKw CharsetNameOrDefault, parser.y:11008 and 13556. p.cur is
// the first word of CharsetKw, and the caller has seen that CHARACTER and CHAR have their
// SET.
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

// parseSetCharsetName reads CharsetName, parser.y:11023: a StringName or BINARY.
func (p *Parser) parseSetCharsetName() (string, int, error) {
	if p.cur.Type == kwBINARY {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseSetStringName()
}

// parseSetStringName reads StringName, parser.y:13604: a string literal or an Identifier.
// The lexer hands a user variable over as an identifier, and `SET NAMES @a` is a 1064.
func (p *Parser) parseSetStringName() (string, int, error) {
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

// nextEndsVariableName reports whether the word at p.cur is a variable's name rather
// than the keyword of a SET form. VariableName (parser.y:10924) is always followed by
// `=`, `:=` or `.`, and no form continues its keyword with one of those, except the `=`
// of SET PASSWORD, which the caller settles.
//
// A next token at or past the completion cursor is text the user has not reached, so
// `SET TRANSACTION <cursor> = 1` still offers what follows the keyword.
func (p *Parser) nextEndsVariableName() bool {
	next := p.peekNext()
	if p.completing && next.Loc >= p.cursorOff {
		return false
	}
	return next.Type == '=' || next.Type == tokAssign || next.Type == '.'
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

// parseSetAssignment parses a single SET assignment with no scope word: var = expr
func (p *Parser) parseSetAssignment() (*nodes.Assignment, error) {
	start := p.pos()

	// Handle @var and @@var references
	if !p.isVariableRef() {
		return p.parseSetVariableAssignment(start)
	}
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
	col := &nodes.ColumnRef{
		Loc:    vref.Loc,
		Column: prefix + vref.Name,
	}
	if vref.System {
		return p.parseSetAssignmentValue(start, col)
	}

	// A user variable takes an Expression, parser.y:10975, not a SetExpr: `SET @a = ON`,
	// `SET @a = BINARY` and `SET @a = DEFAULT` are 1064s.
	if err := p.parseSetAssignOp(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwDEFAULT && p.peekNext().Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseSetExpression(start, col)
}

// parseSetVariableAssignment parses VariableName EqOrAssignmentEq SetExpr. start is where
// the assignment began, at its scope word when it has one.
func (p *Parser) parseSetVariableAssignment(start int) (*nodes.Assignment, error) {
	// The lexer hands a user variable over as an identifier, and after a scope
	// VariableName is all that may follow.
	if p.isVariableRef() {
		return nil, p.syntaxErrorAtCur()
	}

	// VariableName, parser.y:10924: Identifier, or two of them with a dot. An
	// Identifier is any unreserved word, parser.y:6774, the scope words included:
	// parseSetStmt has already read a scope when one was meant, and MySQL's narrower
	// lvalue_ident is not TiDB's rule here.
	colStart := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	col := &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: colStart},
		Column: name,
	}
	// Support dot-qualified variable names: schema.var
	if p.cur.Type == '.' {
		p.advance()
		if p.isVariableRef() {
			return nil, p.syntaxErrorAtCur()
		}
		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		col.Table = name
		col.Column = name2
	}
	col.Loc.End = p.prev.End

	return p.parseSetAssignmentValue(start, col)
}

// parseSetAssignmentValue parses EqOrAssignmentEq SetExpr, parser.y:10909: ON, BINARY,
// DEFAULT or an Expression. ON and BINARY are the value as text, as upstream has them, so
// a reader of `SET a = ON` and `SET a = 'ON'` sees one thing.
//
// BINARY also starts an Expression, as in `SET a = BINARY 'x'`, so it is the value only
// where the assignment ends. DEFAULT is the value unless a column follows it, so `SET a =
// DEFAULT + 1` is a 1064 where parseExpr would go on to read the sum.
func (p *Parser) parseSetAssignmentValue(start int, col *nodes.ColumnRef) (*nodes.Assignment, error) {
	if err := p.parseSetAssignOp(); err != nil {
		return nil, err
	}

	var val nodes.ExprNode
	switch next := p.peekNext().Type; {
	case p.cur.Type == kwON,
		p.cur.Type == kwBINARY && (next == ',' || next == ';' || next == tokEOF):
		tok := p.advance()
		val = &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: strings.ToUpper(tok.Str)}
	case p.cur.Type == kwDEFAULT && next != '(':
		expr, err := p.parseDefaultExpr()
		if err != nil {
			return nil, err
		}
		val = expr
	default:
		return p.parseSetExpression(start, col)
	}
	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Column: col,
		Value:  val,
	}, nil
}

// parseSetAssignOp consumes `=` or `:=`.
func (p *Parser) parseSetAssignOp() error {
	if p.cur.Type == tokAssign {
		p.advance()
		return nil
	}
	_, err := p.expect('=')
	return err
}

// parseSetExpression parses an assignment's value as an Expression.
func (p *Parser) parseSetExpression(start int, col *nodes.ColumnRef) (*nodes.Assignment, error) {
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

// showOneWordTargets are TiDB's filterable SHOW targets that are a single
// unreserved keyword, each mapped to the ShowStmt type it sets. They are a
// table rather than switch arms because nothing distinguishes them but the
// word (ShowTargetFilterable, parser.y:12232-12310).
var showOneWordTargets = map[int]string{
	kwBACKUPS:              "BACKUPS",
	kwRESTORES:             "RESTORES",
	kwCONFIG:               "CONFIG",
	kwSESSION_STATES:       "SESSION_STATES",
	kwSTATS_EXTENDED:       "STATS_EXTENDED",
	kwSTATS_META:           "STATS_META",
	kwSTATS_HISTOGRAMS:     "STATS_HISTOGRAMS",
	kwSTATS_TOPN:           "STATS_TOPN",
	kwSTATS_BUCKETS:        "STATS_BUCKETS",
	kwSTATS_HEALTHY:        "STATS_HEALTHY",
	kwSTATS_LOCKED:         "STATS_LOCKED",
	kwHISTOGRAMS_IN_FLIGHT: "HISTOGRAMS_IN_FLIGHT",
	kwCOLUMN_STATS_USAGE:   "COLUMN_STATS_USAGE",
	kwBINDINGS:             "BINDINGS",
}

// parseShowStmt parses a SHOW statement.
//
// TiDB's SHOW is MySQL's target table less the statements TiDB has no
// subsystem for, plus its own. Upstream splits it in two: ShowStmt
// (parser.y:11755) holds the targets that take no filter, and
// ShowTargetFilterable (parser.y:12067) the ones ShowLikeOrWhereOpt may
// follow. Every arm below is one of those two.
//
// A target MySQL has that is in neither is a 1064 on TiDB v8.5.5, measured
// in showTargetCases: SHOW ENGINE name {STATUS | MUTEX | LOGS}, SHOW STORAGE
// ENGINES, SHOW {MASTER | BINARY} LOGS, SHOW SLAVE HOSTS, SHOW REPLICAS, SHOW
// {BINLOG | RELAYLOG} EVENTS, SHOW EXTENDED {INDEX | KEYS | TABLES}, the LIMIT
// of SHOW WARNINGS and SHOW ERRORS, and the FOR CHANNEL of SHOW {SLAVE |
// REPLICA} STATUS. None is parsed here.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/show.html;
// https://docs.pingcap.com/tidb/stable/sql-statement-show-table-regions/
func (p *Parser) parseShowStmt() (*nodes.ShowStmt, error) {
	start := p.pos()
	p.advance() // consume SHOW

	// Completion: after SHOW, offer keyword candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwTABLES, kwCOLUMNS, kwINDEX, kwDATABASES,
			kwCREATE, kwSTATUS, kwVARIABLES, kwPROCESSLIST,
			kwWARNINGS, kwERRORS, kwENGINES, kwBINARY,
			kwGLOBAL, kwSESSION, kwFULL, kwPLACEMENT,
			kwCONFIG, kwSTATS_META, kwBINDINGS,
		} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.ShowStmt{Loc: nodes.Loc{Start: start}}

	// SHOW COUNT(*) ERRORS | SHOW COUNT(*) WARNINGS. The parens are part of
	// the production (builtinCount '(' '*' ')', parser.y:12155), so dropping
	// them is a 1064 rather than a bare SHOW COUNT.
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
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}
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
			ref, err := p.parseShowDatabaseName()
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
		switch p.cur.Type {
		case kwCOLUMNS, kwFIELDS:
			stmt.Type = "FULL COLUMNS"
			p.advance()
			if err := p.parseShowColumnsTarget(stmt); err != nil {
				return nil, err
			}
		case kwTABLES:
			stmt.Type = "FULL TABLES"
			p.advance()
			// Optional FROM|IN db
			if p.matchFromOrIn() {
				ref, err := p.parseShowDatabaseName()
				if err != nil {
					return nil, err
				}
				stmt.From = ref
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		case kwPROCESSLIST:
			stmt.Type = "FULL PROCESSLIST"
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwCOLUMNS, kwFIELDS:
		stmt.Type = "COLUMNS"
		p.advance()
		if err := p.parseShowColumnsTarget(stmt); err != nil {
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
			// Optional IF NOT EXISTS
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
			ref, err := p.parseShowDatabaseName()
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
		case kwSEQUENCE:
			// SHOW CREATE SEQUENCE TableName (parser.y:11789). TiDB has
			// sequences; MySQL 8.0 has none and answers 1064 here.
			stmt.Type = "CREATE SEQUENCE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwPLACEMENT:
			// SHOW CREATE PLACEMENT POLICY PolicyName (parser.y:11796),
			// where PolicyName is a bare Identifier: no qualifier, no IF
			// EXISTS, one name.
			p.advance() // consume PLACEMENT
			if _, err := p.expect(kwPOLICY); err != nil {
				return nil, err
			}
			name, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Type = "CREATE PLACEMENT POLICY"
			stmt.PolicyName = name
		case kwRESOURCE:
			// SHOW CREATE RESOURCE GROUP ResourceGroupName
			// (parser.y:11803), where the name is an Identifier or the
			// keyword DEFAULT (parser.y:4495).
			p.advance() // consume RESOURCE
			if _, err := p.expect(kwGROUP); err != nil {
				return nil, err
			}
			if p.cur.Type == kwDEFAULT {
				p.advance()
				stmt.Type = "CREATE RESOURCE GROUP"
				stmt.ResourceGroupName = "DEFAULT"
				break
			}
			name, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Type = "CREATE RESOURCE GROUP"
			stmt.ResourceGroupName = name
		case kwUSER:
			stmt.Type = "CREATE USER"
			p.advance()
			// Handle CURRENT_USER / CURRENT_USER()
			if p.cur.Type == kwCURRENT_USER {
				start := p.pos()
				p.advance()
				if p.cur.Type == '(' {
					p.advance()
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
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
			// TiDB v8.5.5 has no functions, triggers or events, so SHOW
			// CREATE {FUNCTION,TRIGGER,EVENT} is a 1064 at the object
			// keyword, and so is a CREATE with nothing after it.
			return nil, p.syntaxErrorAtCur()
		}

	case kwEXTENDED:
		// EXTENDED OptFull FieldsOrColumns ... (parser.y:12141) is the only
		// arm TiDB has: SHOW EXTENDED {INDEX | KEYS | TABLES} is MySQL's and
		// a 1064 here.
		p.advance() // consume EXTENDED
		if p.cur.Type == kwFULL {
			p.advance()
			stmt.Type = "EXTENDED FULL COLUMNS"
		} else {
			stmt.Type = "EXTENDED COLUMNS"
		}
		if p.cur.Type != kwCOLUMNS && p.cur.Type != kwFIELDS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.parseShowColumnsTarget(stmt); err != nil {
			return nil, err
		}

	case kwINDEX:
		// INDEX, INDEXES and KEYS are the three spellings (ShowIndexKwd,
		// parser.y:12048); KEY on its own is not one of them and is a 1064.
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
		if err := p.parseShowTableAndDatabase(stmt); err != nil {
			return nil, err
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
		switch p.cur.Type {
		case kwVARIABLES:
			stmt.Type = scope + " VARIABLES"
		case kwSTATUS:
			stmt.Type = scope + " STATUS"
		case kwBINDINGS:
			// GlobalScope "BINDINGS" (parser.y:12218) — TiDB's SQL
			// bindings, which MySQL does not have.
			stmt.Type = scope + " BINDINGS"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
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
		// WARNINGS and ERRORS are filterable targets and nothing more:
		// TiDB has no LIMIT on either, where MySQL does.
		stmt.Type = "WARNINGS"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwERRORS:
		stmt.Type = "ERRORS"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwENGINES:
		stmt.Type = "ENGINES"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwPLUGINS:
		stmt.Type = "PLUGINS"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwMASTER:
		// "SHOW" "MASTER" "STATUS" (parser.y:11871). SHOW MASTER LOGS is
		// MySQL's and a 1064 here.
		p.advance() // consume MASTER
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "MASTER STATUS"

	case kwBINARY:
		// "SHOW" "BINARY" "LOG" "STATUS" (parser.y:11879), the MySQL 8.2
		// spelling of SHOW MASTER STATUS. SHOW BINARY LOGS is a 1064.
		p.advance() // consume BINARY
		if _, err := p.expect(kwLOG); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "BINARY LOG STATUS"

	case kwSLAVE, kwREPLICA:
		// "SHOW" Replica "STATUS" (parser.y:11885), with no FOR CHANNEL and
		// no SLAVE HOSTS.
		replica := p.cur.Type == kwREPLICA
		p.advance()
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		if replica {
			stmt.Type = "REPLICA STATUS"
		} else {
			stmt.Type = "SLAVE STATUS"
		}

	case kwTABLE:
		p.advance() // consume TABLE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "TABLE STATUS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
			break
		}
		if err := p.parseShowTableTarget(stmt); err != nil {
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
		// TiDB v8.5.5 accepts only SHOW PROCEDURE STATUS — not CODE, not
		// the bare form.
		p.advance() // consume PROCEDURE
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "PROCEDURE STATUS"
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwFUNCTION:
		// As with PROCEDURE: STATUS and nothing else.
		p.advance() // consume FUNCTION
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "FUNCTION STATUS"
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
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

	case kwPROFILE:
		stmt.Type = "PROFILE"
		p.advance() // consume PROFILE
		if err := p.parseShowProfileOptions(stmt); err != nil {
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

	case kwPLACEMENT:
		// SHOW PLACEMENT and SHOW PLACEMENT LABELS are filterable
		// (parser.y:12297-12302); SHOW PLACEMENT FOR ShowPlacementTarget
		// (parser.y:11932) is its own arm and takes no filter.
		p.advance() // consume PLACEMENT
		switch p.cur.Type {
		case kwLABELS:
			p.advance()
			stmt.Type = "PLACEMENT LABELS"
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		case kwFOR:
			p.advance()
			if err := p.parseShowPlacementTarget(stmt); err != nil {
				return nil, err
			}
		default:
			stmt.Type = "PLACEMENT"
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwANALYZE:
		// "ANALYZE" "STATUS" (parser.y:12277), filterable. ANALYZE is
		// reserved, so nothing else reaches this arm.
		p.advance() // consume ANALYZE
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "ANALYZE STATUS"
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwBINDING_CACHE:
		// "BINDING_CACHE" "STATUS" (parser.y:12204).
		p.advance() // consume BINDING_CACHE
		if _, err := p.expect(kwSTATUS); err != nil {
			return nil, err
		}
		stmt.Type = "BINDING_CACHE STATUS"
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwDISTRIBUTION:
		// "DISTRIBUTION" "JOBS" is filterable (parser.y:12306);
		// "SHOW" "DISTRIBUTION" "JOB" Int64Num is its own arm
		// (parser.y:11945) and takes an integer literal and no filter.
		p.advance() // consume DISTRIBUTION
		switch p.cur.Type {
		case kwJOBS:
			p.advance()
			stmt.Type = "DISTRIBUTION JOBS"
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		case kwJOB:
			p.advance()
			id, err := p.parseShowJobID()
			if err != nil {
				return nil, err
			}
			stmt.Type = "DISTRIBUTION JOB"
			stmt.JobID = id
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwBUILTINS:
		// Its own arm (parser.y:11926), so no filter follows it.
		stmt.Type = "BUILTINS"
		p.advance()

	default:
		// The remaining targets are single unreserved keywords, every one
		// of them filterable.
		if name, ok := showOneWordTargets[p.cur.Type]; ok {
			stmt.Type = name
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
			break
		}

		// SHOW REPLICAS is MySQL's SHOW SLAVE HOSTS and a 1064 on TiDB, so
		// kwREPLICAS is not an arm.

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
						if _, err := p.expect(')'); err != nil {
							return nil, err
						}
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
						if err := p.checkRoleName(); err != nil {
							return nil, err
						}
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
			if err := p.parseShowTableAndDatabase(stmt); err != nil {
				return nil, err
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
		// separator. That turned `SHOW PLACEMENT FOR DATABASE db1 TABLE tb1`
		// into a SHOW plus a `TABLE tb1` -- a 1064 on TiDB v8.5.5, two
		// statements here -- and did the same for every other target with a
		// statement after it, `SHOW DATABASES TABLE t` included. The check
		// belongs to SHOW rather than to the loop: `CREATE TABLE ta TABLE tb`
		// is one statement the loop reads as two, and closing it there needs
		// CREATE TABLE ... [AS] {TABLE t | VALUES ...} first.
		return nil, p.syntaxErrorAtCur()
	}

	if stmt.Type == "" {
		// Unreachable: every arm above either sets Type or returns an
		// error. The guard is here because the arms used to fall through
		// silently, which left an empty SHOW in front of the rest of the
		// text -- `show analyze status` parsed as an empty SHOW plus an
		// ANALYZE TABLE, and `SHOW FULL` parsed as a statement.
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseShowColumnsTarget reads the rest of SHOW [EXTENDED] [FULL]
// {COLUMNS | FIELDS}: ShowTableAliasOpt ShowDatabaseNameOpt
// ShowLikeOrWhereOpt (parser.y:12131). The table is required; the database
// is the second FROM|IN and lands on the table's schema.
func (p *Parser) parseShowColumnsTarget(stmt *nodes.ShowStmt) error {
	if _, err := p.expectFromOrIn(); err != nil {
		return err
	}
	// Completion: after SHOW COLUMNS FROM, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return &ParseError{Message: "collecting"}
	}
	if err := p.parseShowTableAndDatabase(stmt); err != nil {
		return err
	}
	return p.parseShowLikeOrWhere(stmt)
}

// parseShowTableAndDatabase reads a table reference and the optional
// FROM|IN database that may follow it, merging the database into the
// reference's schema.
func (p *Parser) parseShowTableAndDatabase(stmt *nodes.ShowStmt) error {
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	stmt.From = ref
	if p.matchFromOrIn() {
		dbRef, err := p.parseShowDatabaseName()
		if err != nil {
			return err
		}
		stmt.From.Schema = dbRef.Name
	}
	return nil
}

// parseShowDatabaseName reads the DBName of ShowDatabaseNameOpt, parser.y:12128 and 4439:
// one Identifier, so `SHOW TABLES FROM a.b` is a 1064. It holds the name the way a table
// reference would, because stmt.From is where a SHOW's database goes.
func (p *Parser) parseShowDatabaseName() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.TableRef{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name}, nil
}

// parseShowTableTarget reads TiDB's four SHOW TABLE forms, which share a
// table name and a partition list:
//
//	SHOW TABLE tbl [PARTITION (p, ...)] REGIONS [WHERE expr]
//	SHOW TABLE tbl [PARTITION (p, ...)] INDEX idx REGIONS [WHERE expr]
//	SHOW TABLE tbl [PARTITION (p, ...)] DISTRIBUTIONS [WHERE expr]
//	SHOW TABLE tbl NEXT_ROW_ID
//
// The partition list rides on the table reference, where upstream keeps it.
// NEXT_ROW_ID takes no partition list and none of the four takes LIKE: the
// filter is WhereClauseOptional, not ShowLikeOrWhereOpt.
//
// Ref: parser.y:11818-11849;
// https://docs.pingcap.com/tidb/stable/sql-statement-show-table-regions/
func (p *Parser) parseShowTableTarget(stmt *nodes.ShowStmt) error {
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	stmt.From = ref

	partitioned := false
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		parts, err := p.parseParenIdentList()
		if err != nil {
			return err
		}
		stmt.From.Partitions = parts
		stmt.From.Loc.End = p.prev.End
		partitioned = true
	}

	switch p.cur.Type {
	case kwREGIONS:
		p.advance()
		stmt.Type = "TABLE REGIONS"
	case kwDISTRIBUTIONS:
		p.advance()
		stmt.Type = "TABLE DISTRIBUTIONS"
	case kwINDEX:
		p.advance() // consume INDEX
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		if _, err := p.expect(kwREGIONS); err != nil {
			return err
		}
		stmt.Type = "TABLE INDEX REGIONS"
		stmt.IndexName = name
	case kwNEXT_ROW_ID:
		if partitioned {
			// The NEXT_ROW_ID arm reads no PartitionNameListOpt.
			return p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.Type = "TABLE NEXT_ROW_ID"
		return nil
	default:
		return p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwWHERE {
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

// parseShowPlacementTarget reads what follows SHOW PLACEMENT FOR:
//
//	{DATABASE | SCHEMA} db
//	TABLE tbl [PARTITION p]
//
// Ref: ShowPlacementTarget, parser.y:11964.
func (p *Parser) parseShowPlacementTarget(stmt *nodes.ShowStmt) error {
	switch p.cur.Type {
	case kwDATABASE, kwSCHEMA:
		p.advance()
		name, nameStart, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		stmt.Type = "PLACEMENT FOR DATABASE"
		stmt.From = &nodes.TableRef{
			Loc:  nodes.Loc{Start: nameStart, End: p.prev.End},
			Name: name,
		}
		return nil
	case kwTABLE:
		p.advance()
		ref, err := p.parseTableRef()
		if err != nil {
			return err
		}
		stmt.Type = "PLACEMENT FOR TABLE"
		stmt.From = ref
		if p.cur.Type == kwPARTITION {
			p.advance()
			name, _, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			stmt.Type = "PLACEMENT FOR PARTITION"
			stmt.Partition = name
		}
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// parseShowJobID reads the Int64Num a job id takes (parser.y:6189): an
// integer literal, so a string, a name or a negative number is a 1064.
func (p *Parser) parseShowJobID() (nodes.ExprNode, error) {
	start := p.cur.Loc
	id, err := p.parseAdminInt64()
	if err != nil {
		return nil, err
	}
	return &nodes.IntLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: id}, nil
}

// parseShowLikeOrWhere parses the optional LIKE or WHERE of a filterable SHOW
// target. The two take different expressions upstream
// (ShowLikeOrWhereOpt, parser.y:12310): LIKE takes a SimpleExpr and WHERE a
// whole Expression, so `SHOW DATABASES LIKE 't' LIKE 'u'` and `SHOW DATABASES
// LIKE 't' AND 1` are a 1064 while the same text after WHERE parses. Reading
// a whole expression after LIKE took both.
func (p *Parser) parseShowLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.cur.Type == kwLIKE {
		p.advance()
		expr, err := p.parseShowLikePattern()
		if err != nil {
			return err
		}
		if expr == nil {
			return p.syntaxErrorAtCur()
		}
		stmt.Like = expr
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

// parseShowLikePattern reads the SimpleExpr a SHOW filter's LIKE takes
// (parser.y:8071). SimpleExpr is prefix and atomic forms plus two infix ones,
// `SimpleExpr "COLLATE" CollationName` (parser.y:8076) and `SimpleExpr pipes
// SimpleExpr` (parser.y:8104), and the arithmetic, bitwise, comparison and
// boolean operators of a whole Expression are none of it.
//
// Of the two infix forms only COLLATE is read here. TiDB v8.5.5 answers 1064
// on `SHOW DATABASES LIKE 'a' || 'b'`: outside PIPES_AS_CONCAT the lexer
// gives `||` as OR, which is no SimpleExpr, so the grammar's pipes arm is
// unreachable in the mode the server runs in. `+`, `*`, `|` and `->` after
// the pattern are a 1064 for the plainer reason that SimpleExpr has no arm
// for them, and the prefix parse already stops short of all four.
//
// The collation is a CollationName, so a quoted name and the keyword BINARY
// are as good as an identifier: `COLLATE 'utf8mb4_bin'` and `COLLATE binary`
// both parse, where reading a bare identifier took neither.
func (p *Parser) parseShowLikePattern() (nodes.ExprNode, error) {
	start := p.pos()
	expr, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	for p.cur.Type == kwCOLLATE {
		p.advance() // consume COLLATE
		collation, err := p.parseCollationName()
		if err != nil {
			return nil, err
		}
		expr = &nodes.CollateExpr{
			Loc:       nodes.Loc{Start: start, End: p.prev.End},
			Expr:      expr,
			Collation: collation,
		}
	}
	return expr, nil
}

// parseShowFromLikeOrWhere parses optional FROM|IN db, LIKE, or WHERE for SHOW statements.
func (p *Parser) parseShowFromLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.matchFromOrIn() {
		ref, err := p.parseShowDatabaseName()
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
			// Int64Num (parser.y:11817), so a name, an expression and 2^63 are 1064.
			expr, err := p.parseShowJobID()
			if err != nil {
				return err
			}
			stmt.ForQuery = expr
		}
	}

	// Optional LIMIT {row_count | offset ',' row_count | row_count OFFSET
	// offset}. SHOW PROFILE reads a whole SelectStmtLimitOpt
	// (parser.y:11911), which is the comma form as well as the OFFSET one,
	// and the two cannot be combined.
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

// explainFormatNames are the bare names TiDB reads after FORMAT =. The list is
// a closed set of keyword tokens in the engine's own grammar
// (ExplainFormatType), so a name outside it, a backquoted name and a reserved
// word are each a 1064: `FORMAT = tree`, `FORMAT = nosuchformat`,
// "FORMAT = `row`" and `FORMAT = binary` are syntax errors on TiDB v8.5.5,
// while every one of these eight parses. A quoted name is a separate arm and
// carries any text at all, checked as a name when the statement runs (1791),
// so `FORMAT = 'tree'` parses and `FORMAT = tree` does not.
//
// Ref: ExplainFormatType in TiDB v8.5.5 parser/parser.y.
var explainFormatNames = map[string]bool{
	"TRADITIONAL":    true,
	"JSON":           true,
	"ROW":            true,
	"DOT":            true,
	"BRIEF":          true,
	"VERBOSE":        true,
	"TRUE_CARD_COST": true,
	"TIDB_JSON":      true,
}

// parseExplainFormatName parses the name after EXPLAIN's FORMAT =. It is a
// string literal in either quote, holding any text, or one of the eight bare
// names in explainFormatNames. TRADITIONAL, DOT, BRIEF, VERBOSE,
// TRUE_CARD_COST and TIDB_JSON are not keywords in this lexer, so the bare arm
// matches on the token's text and turns a backquoted name away by reading the
// source: a backquoted identifier is not a keyword in the engine's grammar
// either.
func (p *Parser) parseExplainFormatName() (string, error) {
	if p.cur.Type == tokSCONST {
		return p.advance().Str, nil
	}
	if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
		if name := strings.ToUpper(p.cur.Str); explainFormatNames[name] && !p.backtickQuotedAt(p.cur.Loc) {
			return p.advance().Str, nil
		}
	}
	if p.cur.Type == tokEOF {
		return "", p.syntaxErrorAtCur()
	}
	return "", &ParseError{
		Message:  "expected an EXPLAIN format name or a quoted one",
		Position: p.cur.Loc,
	}
}

// parseExplainStmt parses an EXPLAIN, DESCRIBE or DESC statement.
//
// Ref: ExplainStmt, ExplainFormatType and ExplainableStmt in TiDB v8.5.5
// parser/parser.y; https://docs.pingcap.com/tidb/stable/sql-statement-explain/
//
//	{EXPLAIN | DESCRIBE | DESC} tbl_name [col_name]
//	{EXPLAIN | DESCRIBE | DESC} [FORMAT = format] {explainable_stmt | FOR CONNECTION connection_id}
//	{EXPLAIN | DESCRIBE | DESC} ANALYZE [FORMAT = format] explainable_stmt
//	format: {TRADITIONAL | JSON | ROW | DOT | BRIEF | VERBOSE | TRUE_CARD_COST | TIDB_JSON | 'any text'}
//	explainable_stmt: {SELECT | TABLE | VALUES | DELETE | INSERT | REPLACE | UPDATE | ALTER TABLE}
//
// The modifier list is the engine's, not MySQL's. TiDB has no EXPLAIN EXTENDED
// and no EXPLAIN PARTITIONS -- MySQL 8.0 dropped both and TiDB never had them,
// so each is a 1064 here, as it is on MySQL 8.0.46. ANALYZE reaches only an
// explainable statement: EXPLAIN ANALYZE FOR CONNECTION 1 and
// EXPLAIN ANALYZE t are syntax errors, which MySQL takes and TiDB does not.
// A modifier also rules the table form out, so FORMAT = ROW t is a 1064.
//
// ALTER TABLE is explainable on TiDB and on no other engine in this family.
// EXPLAIN ALTER TABLE reports the plan for the change without making it;
// EXPLAIN ANALYZE ALTER TABLE makes it, which is why review's walk follows the
// statement ANALYZE carries.
//
// DESC and DESCRIBE read only the table form here. TiDB takes them as full
// EXPLAIN synonyms -- DESC FORMAT = 'row' SELECT 1 parses on v8.5.5 -- which is
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
		if err := p.parseExplainColumn(showStmt); err != nil {
			return nil, err
		}
		showStmt.Loc.End = p.prev.End
		stmt.Stmt = showStmt
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if p.cur.Type == kwANALYZE {
		stmt.Analyze = true
		p.advance()
	}

	hasFormat := false
	if p.cur.Type == kwFORMAT {
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

	// EXPLAIN [FORMAT = format] FOR CONNECTION connection_id. CONNECTION and
	// the id are both required, and the id is an unsigned integer: FOR 1,
	// FOR CONNECTION and FOR CONNECTION 'x' are each a 1064. ANALYZE does not
	// reach this arm, so EXPLAIN ANALYZE FOR CONNECTION 1 falls through to the
	// explainable statement below and fails there.
	if p.cur.Type == kwFOR && !stmt.Analyze {
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
			kwANALYZE, kwFORMAT, kwFOR,
			kwSELECT, kwWITH, kwTABLE, kwINSERT, kwUPDATE, kwDELETE, kwREPLACE, kwVALUES, kwALTER,
		} {
			p.addTokenCandidate(t)
		}
		if !stmt.Analyze && !hasFormat {
			// The table form is still reachable here, and EXPLAIN t is as
			// valid as EXPLAIN SELECT.
			p.addRuleCandidate("table_ref")
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// The explainable statement, which PLAN REPLAYER DUMP EXPLAIN names too:
	// both slots are upstream's one ExplainableStmt production, and
	// parseExplainableStmt is it. ALTER TABLE is the one DDL in the set --
	// ALTER DATABASE, CREATE TABLE, TRUNCATE and RENAME are each a 1064.
	if p.atExplainableStmt() {
		inner, err := p.parseExplainableStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = inner
	} else {
		// EXPLAIN tbl_name [col_name], which no modifier stands in front of.
		// hasFormat rather than Format != "": the name may be an empty
		// string, and FORMAT = '' is a modifier all the same.
		if stmt.Analyze || hasFormat {
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

// parseExplainColumn parses the column name EXPLAIN and DESCRIBE take after a
// table name, into the ShowStmt the table form lowers to. TiDB's arm is
// ExplainSym TableName ColumnName, so the name is one to three dotted
// identifiers and nothing else: EXPLAIN t 'abc' and EXPLAIN t 123 are each a
// 1064, where MySQL 8.0.46 reads the string as a wild pattern. A qualified name
// parses, as ColumnName carries one.
//
// It is held in ShowStmt.Like, which is where the EXPLAIN table form has always
// put it; ShowStmt has no column field of its own.
func (p *Parser) parseExplainColumn(showStmt *nodes.ShowStmt) error {
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil
	}
	col, err := p.parseColumnRef()
	if err != nil {
		return err
	}
	if col == nil || col.Star {
		return p.syntaxErrorAtCur()
	}
	showStmt.Like = col
	return nil
}
