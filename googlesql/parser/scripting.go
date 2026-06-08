package parser

// This file is the `parser-scripting` DAG node. It implements GoogleSQL's
// procedural language (GoogleSQLParser.g4 §2.11, a hand-port of Google's
// open-source ZetaSQL reference and the grammar bytebase consumes):
//
//	begin_end_block:        BEGIN statement_list? opt_exception_handler? END
//	opt_exception_handler:  EXCEPTION WHEN ERROR THEN statement_list
//	statement_list:         unterminated_statement (';' unterminated_statement)* ';'
//	unterminated_statement: unterminated_sql_statement | unterminated_script_statement
//	if_statement:           IF expr THEN statement_list? elseif_clauses? opt_else? END IF
//	case_statement:         CASE expr? when_then_clauses opt_else? END CASE
//	while_statement:        WHILE expr DO statement_list? END WHILE
//	loop_statement:         LOOP statement_list? END LOOP
//	repeat_statement:       REPEAT statement_list? UNTIL expr END REPEAT
//	for_in_statement:       FOR id IN (query) DO statement_list? END FOR
//	variable_declaration:   DECLARE id_list type opt_default? | DECLARE id_list DEFAULT expr
//	set_statement:          SET TRANSACTION modes | SET id|@p|@@v = expr
//	                        | SET (id_list) = expr | SET id, id = (error alt)
//	execute_immediate:      EXECUTE IMMEDIATE expr (INTO id_list)? (USING arg_list)?
//	break_statement:        BREAK id? | LEAVE id?
//	continue_statement:     CONTINUE id? | ITERATE id?
//	return_statement:       RETURN
//	raise_statement:        RAISE | RAISE USING MESSAGE = expr
//	(labeled)               label ':' (begin|while|loop|repeat|for_in) identifier?
//
// CORRECTNESS (correctness-protocol.md). Scripting is BigQuery-ONLY at the
// GoogleSQL union level, with two oracle-authoritative exceptions confirmed on
// the live Cloud Spanner emulator (scripting_oracle_test.go):
//   - the BEGIN…END envelope (verdict accept, "Statement not supported:
//     BeginStmt"), SET (SingleAssignment / AssignmentFromStruct /
//     ParameterAssignment / SystemVariableAssignment / SetTransactionStatement),
//     and EXECUTE IMMEDIATE (ExecuteImmediateStatement) — the LEADING forms are
//     ACCEPTED on the oracle's authority.
//   - the control-flow statements (IF / CASE / WHILE / LOOP / REPEAT / FOR-IN /
//     DECLARE / BREAK / CONTINUE / RETURN / RAISE / labels) syntax-reject on
//     Spanner; the union parser accepts them on the authority of the legacy .g4
//     + the BigQuery truth1 corpus (SCRIPT-001…019), triangulated.
//
// The emulator's BEGIN…END recognizer is SHALLOW — it accepts ANY token run
// between BEGIN and END (`BEGIN SELECT 1 SELECT 2 END`, `BEGIN; END`), so the
// INTERIOR statement_list grammar is NON-authoritative on Spanner. omni follows
// the .g4: statement_list requires `;`-separated statements and a trailing `;`
// (divergence ledger: Spanner over-accepts the block interior, exactly mirroring
// the transaction-tail divergence in transaction.go).
//
// SPLITTING. The block-aware Split (split.go) keeps a procedural block (and a
// top-level BEGIN…END) whole, so a whole block — internal `;` and all — arrives
// at parseSingle as ONE segment. parseStmt dispatches the leading keyword to one
// of the functions here; the statement_list parser consumes the internal `;`
// itself. The trailing `;` after the block at top level is consumed by Split's
// segmentation, not seen here.

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// statement_list — the procedural body shared by every block form.
// ---------------------------------------------------------------------------

// scriptListEnd is the set of keywords that terminate a statement_list. A
// statement_list is `unterminated_statement (';' unterminated_statement)* ';'`,
// so parsing continues while the parser is at a statement start; it stops at the
// block-closing / clause-introducing keyword that follows the body (END, ELSE,
// ELSEIF, UNTIL, WHEN, EXCEPTION) or at EOF.
//
// The grammar requires a trailing ';' after the last statement, so the body
// statements are separated by ';' AND every statement (including the last) is
// followed by ';' before the terminator. parseStatementList enforces that.

// parseStatementList parses a statement_list: a sequence of
// `unterminated_statement ';'` pairs, stopping when the current token is one of
// the terminator keywords (or EOF). Returns the parsed statements (empty slice
// when the optional statement_list is absent).
//
// Per the .g4, statement_list is `unterminated_non_empty_statement_list ';'`,
// i.e. EVERY statement — the last one included — must be followed by ';'. So
// after each parsed statement we REQUIRE a ';' before the next statement or the
// terminator. An empty body (terminator immediately after the opener) yields an
// empty slice (the enclosing rule's `statement_list?`).
func (p *Parser) parseStatementList(terminators ...int) ([]ast.Node, error) {
	var stmts []ast.Node
	for {
		if p.atStatementListEnd(terminators) {
			return stmts, nil
		}
		stmt, err := p.parseScriptElement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		// statement_list requires a ';' after every statement (the trailing
		// SEMI_SYMBOL of unterminated_non_empty_statement_list ';'). Without it the
		// body is not a well-formed statement_list.
		if _, err := p.expect(int(';')); err != nil {
			return nil, err
		}
	}
}

// atStatementListEnd reports whether the current token ends a statement_list:
// EOF or one of the caller-supplied terminator keywords.
func (p *Parser) atStatementListEnd(terminators []int) bool {
	if p.cur.Type == tokEOF {
		return true
	}
	for _, t := range terminators {
		if p.cur.Type == t {
			return true
		}
	}
	return false
}

// parseScriptElement parses one unterminated_statement inside a statement_list:
// either an unterminated_sql_statement (any top-level SQL statement) or an
// unterminated_script_statement (a nested control-flow / variable statement).
//
// It reuses parseStmt — the same dispatcher used at the top level — so a nested
// SQL statement (SELECT / INSERT / CALL / …) and a nested script statement
// (IF / WHILE / DECLARE / …) are both handled. parseStmt does NOT consume the
// trailing ';'; the enclosing parseStatementList does. A leading statement-level
// hint (`@{…}` / `@n`) on a body statement is skipped first, matching
// unterminated_sql_statement's `statement_level_hint?`.
func (p *Parser) parseScriptElement() (ast.Node, error) {
	if p.cur.Type == int('@') {
		if hintErr := p.skipStatementLevelHint(); hintErr != nil {
			return nil, hintErr
		}
	}
	return p.parseStmt()
}

// ---------------------------------------------------------------------------
// BEGIN … END block
// ---------------------------------------------------------------------------

// parseBeginEndBlock parses a begin_end_block:
//
//	BEGIN statement_list? opt_exception_handler? END
//	opt_exception_handler: EXCEPTION WHEN ERROR THEN statement_list
//
// BEGIN is the current token. This is reached from parseStmt's BEGIN arm when the
// follower is not a TCL transaction follower (isTCLBeginFollower), and from
// parseScriptElement for a nested block.
func (p *Parser) parseBeginEndBlock() (ast.Node, error) {
	begin := p.advance() // BEGIN
	block := &ast.BeginEndBlock{Loc: begin.Loc}

	// statement_list? — the main body, terminated by EXCEPTION or END.
	body, err := p.parseStatementList(kwEXCEPTION, kwEND)
	if err != nil {
		return nil, err
	}
	block.Body = body

	// opt_exception_handler? — EXCEPTION WHEN ERROR THEN statement_list.
	if p.cur.Type == kwEXCEPTION {
		p.advance() // EXCEPTION
		if _, err := p.expect(kwWHEN); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwERROR); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		// The handler statement_list is REQUIRED (the grammar has statement_list,
		// not statement_list?); it is terminated by END.
		handler, err := p.parseStatementList(kwEND)
		if err != nil {
			return nil, err
		}
		block.HasException = true
		block.Exception = handler
	}

	end, err := p.expect(kwEND)
	if err != nil {
		return nil, err
	}
	block.Loc.End = end.Loc.End
	return block, nil
}

// ---------------------------------------------------------------------------
// IF
// ---------------------------------------------------------------------------

// parseIfStmt parses an if_statement:
//
//	IF expr THEN statement_list? elseif_clauses? opt_else? END IF
//
// IF is the current token.
func (p *Parser) parseIfStmt() (ast.Node, error) {
	ifTok := p.advance() // IF
	stmt := &ast.IfStmt{Loc: ifTok.Loc}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Cond = cond

	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}
	// statement_list? terminated by ELSEIF / ELSE / END.
	then, err := p.parseStatementList(kwELSEIF, kwELSE, kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Then = then

	// elseif_clauses? — (ELSEIF expr THEN statement_list?)+.
	for p.cur.Type == kwELSEIF {
		elseifTok := p.advance() // ELSEIF
		clause := &ast.ElseIfClause{Loc: elseifTok.Loc}
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		clause.Cond = cond
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		body, err := p.parseStatementList(kwELSEIF, kwELSE, kwEND)
		if err != nil {
			return nil, err
		}
		clause.Then = body
		clause.Loc.End = p.prev.Loc.End
		stmt.ElseIf = append(stmt.ElseIf, clause)
	}

	// opt_else? — ELSE statement_list?.
	if p.cur.Type == kwELSE {
		p.advance() // ELSE
		elseBody, err := p.parseStatementList(kwEND)
		if err != nil {
			return nil, err
		}
		stmt.HasElse = true
		stmt.Else = elseBody
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endIf, err := p.expect(kwIF)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endIf.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CASE (statement form)
// ---------------------------------------------------------------------------

// parseCaseStmt parses a case_statement:
//
//	CASE expr? when_then_clauses opt_else? END CASE
//	when_then_clauses: (WHEN expr THEN statement_list?)+
//
// CASE is the current token. The optional operand distinguishes the simple form
// (`CASE x WHEN …`) from the searched form (`CASE WHEN …`); it is present when
// the token after CASE is not WHEN.
func (p *Parser) parseCaseStmt() (ast.Node, error) {
	caseTok := p.advance() // CASE
	stmt := &ast.CaseStmt{Loc: caseTok.Loc}

	// expr? — the operand of the simple form; absent when WHEN follows directly.
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Operand = operand
	}

	// when_then_clauses — at least one WHEN…THEN.
	if p.cur.Type != kwWHEN {
		return nil, p.syntaxErrorAtCur()
	}
	for p.cur.Type == kwWHEN {
		whenTok := p.advance() // WHEN
		clause := &ast.WhenThenClause{Loc: whenTok.Loc}
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		clause.Cond = cond
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		body, err := p.parseStatementList(kwWHEN, kwELSE, kwEND)
		if err != nil {
			return nil, err
		}
		clause.Then = body
		clause.Loc.End = p.prev.Loc.End
		stmt.Whens = append(stmt.Whens, clause)
	}

	// opt_else? — ELSE statement_list?.
	if p.cur.Type == kwELSE {
		p.advance() // ELSE
		elseBody, err := p.parseStatementList(kwEND)
		if err != nil {
			return nil, err
		}
		stmt.HasElse = true
		stmt.Else = elseBody
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endCase, err := p.expect(kwCASE)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endCase.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// WHILE / LOOP / REPEAT
// ---------------------------------------------------------------------------

// parseWhileStmt parses a while_statement: `WHILE expr DO statement_list? END
// WHILE`. WHILE is the current token.
func (p *Parser) parseWhileStmt() (ast.Node, error) {
	whileTok := p.advance() // WHILE
	stmt := &ast.WhileStmt{Loc: whileTok.Loc}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Cond = cond

	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}
	body, err := p.parseStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endWhile, err := p.expect(kwWHILE)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endWhile.Loc.End
	return stmt, nil
}

// parseLoopStmt parses a loop_statement: `LOOP statement_list? END LOOP`. LOOP is
// the current token.
func (p *Parser) parseLoopStmt() (ast.Node, error) {
	loopTok := p.advance() // LOOP
	stmt := &ast.LoopStmt{Loc: loopTok.Loc}

	body, err := p.parseStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endLoop, err := p.expect(kwLOOP)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endLoop.Loc.End
	return stmt, nil
}

// parseRepeatStmt parses a repeat_statement: `REPEAT statement_list? until_clause
// END REPEAT`, where until_clause is `UNTIL expr`. REPEAT is the current token.
func (p *Parser) parseRepeatStmt() (ast.Node, error) {
	repeatTok := p.advance() // REPEAT
	stmt := &ast.RepeatStmt{Loc: repeatTok.Loc}

	body, err := p.parseStatementList(kwUNTIL)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	// until_clause — UNTIL expr (required).
	if _, err := p.expect(kwUNTIL); err != nil {
		return nil, err
	}
	until, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Until = until

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endRepeat, err := p.expect(kwREPEAT)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endRepeat.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// FOR … IN
// ---------------------------------------------------------------------------

// parseForInStmt parses a for_in_statement:
//
//	FOR identifier IN (query) DO statement_list? END FOR
//
// FOR is the current token. The loop source is a parenthesized_query.
func (p *Parser) parseForInStmt() (ast.Node, error) {
	forTok := p.advance() // FOR
	stmt := &ast.ForInStmt{Loc: forTok.Loc}

	// loop variable.
	varTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(varTok)
	if err != nil {
		return nil, err
	}
	stmt.Var = &ast.Identifier{Name: name, Loc: varTok.Loc}

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	// parenthesized_query: '(' query ')'.
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	p.fillSubqueries(query)
	stmt.Query = query

	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}
	body, err := p.parseStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endFor, err := p.expect(kwFOR)
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = endFor.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DECLARE
// ---------------------------------------------------------------------------

// parseDeclareStmt parses a variable_declaration:
//
//	DECLARE identifier_list type opt_default_expression?
//	DECLARE identifier_list DEFAULT expression
//
// DECLARE is the current token. The two alternatives are disambiguated after the
// identifier_list: a DEFAULT keyword next selects the type-less alternative; any
// other token begins the required type.
func (p *Parser) parseDeclareStmt() (ast.Node, error) {
	declareTok := p.advance() // DECLARE
	stmt := &ast.DeclareStmt{Loc: declareTok.Loc}

	names, err := p.parseScriptIdentifierList()
	if err != nil {
		return nil, err
	}
	stmt.Names = names

	if p.cur.Type == kwDEFAULT {
		// DECLARE id_list DEFAULT expr — no type.
		p.advance() // DEFAULT
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Default = def
	} else {
		// DECLARE id_list type opt_default_expression? — type required.
		dt, err := p.parseType()
		if err != nil {
			return nil, err
		}
		stmt.Type = &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
		// opt_default_expression: DEFAULT expression.
		if p.cur.Type == kwDEFAULT {
			p.advance() // DEFAULT
			def, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Default = def
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	// The DEFAULT expression may embed subqueries (e.g.
	// `DECLARE x DEFAULT (SELECT … LIMIT 1)`); fill them so the query-span walker
	// can reach the source query.
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseScriptIdentifierList parses an identifier_list: `identifier (',' identifier)*`,
// returning the names as *ast.Identifier nodes (≥1). Used by DECLARE and the SET
// tuple LHS.
func (p *Parser) parseScriptIdentifierList() ([]*ast.Identifier, error) {
	var ids []*ast.Identifier
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		ids = append(ids, &ast.Identifier{Name: name, Loc: tok.Loc})
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return ids, nil
}

// ---------------------------------------------------------------------------
// SET
// ---------------------------------------------------------------------------

// parseSetStmt parses a set_statement:
//
//	SET TRANSACTION transaction_mode_list
//	SET identifier = expression
//	SET named_parameter_expression = expression          (@p = expr)
//	SET system_variable_expression = expression          (@@v = expr)
//	SET ( identifier_list ) = expression
//	SET identifier ',' identifier '=' …                  (error alt: requires parens)
//
// SET is the current token. All forms are reachable; the error alt is reported
// with the ZetaSQL message (the .g4 NotifyErrorListeners text, which the Spanner
// emulator echoes verbatim) so the over-permissive bare multi-variable form is
// rejected rather than mis-parsed.
func (p *Parser) parseSetStmt() (ast.Node, error) {
	setTok := p.advance() // SET
	stmt := &ast.SetStmt{Loc: setTok.Loc}

	switch p.cur.Type {
	case kwTRANSACTION:
		// SET TRANSACTION transaction_mode_list — reuse the transaction node's mode
		// parser (transaction.go). The mode list is required (≥1 mode).
		p.advance() // TRANSACTION
		if !p.atTransactionModeStart() {
			return nil, p.syntaxErrorAtCur()
		}
		modes, err := p.parseTransactionModeList()
		if err != nil {
			return nil, err
		}
		stmt.Kind = ast.SetTransaction
		stmt.Modes = modes
		stmt.Loc.End = modes[len(modes)-1].Loc.End
		return stmt, nil

	case int('('):
		// SET ( identifier_list ) = expression.
		p.advance() // '('
		ids, err := p.parseScriptIdentifierList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Kind = ast.SetTuple
		stmt.Targets = make([]ast.Node, len(ids))
		for i, id := range ids {
			stmt.Targets[i] = id
		}
		stmt.Value = value
		stmt.Loc.End = p.prev.Loc.End
		p.fillSubqueries(stmt)
		return stmt, nil

	case int('@'), tokAtAt:
		// SET @p = expr (named parameter, '@') or SET @@v = expr (system variable,
		// '@@'). The lexer emits '@' as int('@') and '@@' as tokAtAt.
		target, err := p.parseSetParameterTarget()
		if err != nil {
			return nil, err
		}
		return p.finishSetVariable(stmt, target)

	default:
		// SET identifier = expression — OR the error alt `SET id, id = …`.
		if !isIdentifierStart(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		idTok := p.advance()
		name, err := p.identifierText(idTok)
		if err != nil {
			return nil, err
		}
		// A ',' here is the error alt: multiple bare variables without parentheses.
		// The .g4 emits a fixed message (echoed verbatim by the Spanner emulator).
		if p.cur.Type == int(',') {
			return nil, &ParseError{
				Loc: idTok.Loc,
				Msg: "Using SET with multiple variables requires parentheses around the variable list",
			}
		}
		target := &ast.Identifier{Name: name, Loc: idTok.Loc}
		return p.finishSetVariable(stmt, target)
	}
}

// parseSetParameterTarget parses the LHS of a SET when it begins with '@':
// `@identifier` (named_parameter_expression) → *ast.Parameter, or `@@path`
// (system_variable_expression) → *ast.SystemVariable. The current token is '@'
// or '@@'.
func (p *Parser) parseSetParameterTarget() (ast.Node, error) {
	if p.cur.Type == tokAtAt {
		return p.parseSystemVariable()
	}
	return p.parseNamedParameter()
}

// finishSetVariable consumes `= expression` for the SetVariable forms (plain
// identifier / @p / @@v) and returns the completed SetStmt. target is the parsed
// LHS.
func (p *Parser) finishSetVariable(stmt *ast.SetStmt, target ast.Node) (ast.Node, error) {
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Kind = ast.SetVariable
	stmt.Targets = []ast.Node{target}
	stmt.Value = value
	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// EXECUTE IMMEDIATE
// ---------------------------------------------------------------------------

// parseExecuteImmediateStmt parses an execute_immediate:
//
//	EXECUTE IMMEDIATE expression (INTO identifier_list)? (USING arg_list)?
//	arg: expression (AS alias)?
//
// EXECUTE is the current token (parseStmt's EXECUTE arm routes here after
// confirming nothing — the IMMEDIATE keyword is required and validated here).
func (p *Parser) parseExecuteImmediateStmt() (ast.Node, error) {
	execTok := p.advance() // EXECUTE
	if _, err := p.expect(kwIMMEDIATE); err != nil {
		return nil, err
	}
	stmt := &ast.ExecuteImmediateStmt{Loc: execTok.Loc}

	sql, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.SQL = sql

	// (INTO identifier_list)?
	if p.cur.Type == kwINTO {
		p.advance() // INTO
		ids, err := p.parseScriptIdentifierList()
		if err != nil {
			return nil, err
		}
		stmt.Into = ids
	}

	// (USING arg_list)? — each arg is `expression (AS alias)?`.
	if p.cur.Type == kwUSING {
		p.advance() // USING
		for {
			arg, err := p.parseUsingArg()
			if err != nil {
				return nil, err
			}
			stmt.Using = append(stmt.Using, arg)
			if _, ok := p.match(int(',')); ok {
				continue
			}
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	// The dynamic SQL expression and USING argument expressions may embed
	// subqueries; fill them for the query-span walker.
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseUsingArg parses one EXECUTE IMMEDIATE USING argument:
// `expression (AS alias)?`. The alias binds a named @parameter in the dynamic
// SQL; it is recorded as the spelled name.
func (p *Parser) parseUsingArg() (*ast.UsingArg, error) {
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	arg := &ast.UsingArg{Value: value}
	arg.Loc = ast.NodeLoc(value)
	if p.cur.Type == kwAS {
		p.advance() // AS
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		arg.Alias = name
		arg.Loc.End = aliasTok.Loc.End
	}
	return arg, nil
}

// ---------------------------------------------------------------------------
// BREAK / LEAVE / CONTINUE / ITERATE / RETURN / RAISE
// ---------------------------------------------------------------------------

// parseBreakStmt parses a break_statement: `BREAK identifier?` or
// `LEAVE identifier?`. The current token is BREAK or LEAVE.
func (p *Parser) parseBreakStmt() (ast.Node, error) {
	tok := p.advance() // BREAK | LEAVE
	stmt := &ast.BreakStmt{IsLeave: tok.Type == kwLEAVE, Loc: tok.Loc}
	if isIdentifierStart(p.cur.Type) {
		labelTok := p.advance()
		name, err := p.identifierText(labelTok)
		if err != nil {
			return nil, err
		}
		stmt.Label = name
		stmt.Loc.End = labelTok.Loc.End
	}
	return stmt, nil
}

// parseContinueStmt parses a continue_statement: `CONTINUE identifier?` or
// `ITERATE identifier?`. The current token is CONTINUE or ITERATE.
func (p *Parser) parseContinueStmt() (ast.Node, error) {
	tok := p.advance() // CONTINUE | ITERATE
	stmt := &ast.ContinueStmt{IsIterate: tok.Type == kwITERATE, Loc: tok.Loc}
	if isIdentifierStart(p.cur.Type) {
		labelTok := p.advance()
		name, err := p.identifierText(labelTok)
		if err != nil {
			return nil, err
		}
		stmt.Label = name
		stmt.Loc.End = labelTok.Loc.End
	}
	return stmt, nil
}

// parseReturnStmt parses a return_statement: a bare `RETURN`. RETURN is the
// current token; GoogleSQL's RETURN carries no value.
func (p *Parser) parseReturnStmt() (ast.Node, error) {
	tok := p.advance() // RETURN
	return &ast.ReturnStmt{Loc: tok.Loc}, nil
}

// parseRaiseStmt parses a raise_statement:
//
//	RAISE
//	RAISE USING MESSAGE = expression
//
// RAISE is the current token.
func (p *Parser) parseRaiseStmt() (ast.Node, error) {
	raiseTok := p.advance() // RAISE
	stmt := &ast.RaiseStmt{Loc: raiseTok.Loc}
	if p.cur.Type == kwUSING {
		p.advance() // USING
		if _, err := p.expect(kwMESSAGE); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		msg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Message = msg
		stmt.Loc.End = p.prev.Loc.End
		p.fillSubqueries(stmt)
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Labeled statements
// ---------------------------------------------------------------------------

// parseLabeledStmt parses the labeled-statement alternative of
// unterminated_script_statement:
//
//	label ':' unterminated_unlabeled_script_statement identifier?
//	unterminated_unlabeled_script_statement: begin_end_block | while | loop | repeat | for_in
//
// The current token is the label identifier; the next token is ':'. Only the
// block / loop forms may be labeled (NOT IF / CASE / DECLARE / SET / …). A
// trailing identifier after the block's END is the optional end-label.
func (p *Parser) parseLabeledStmt() (ast.Node, error) {
	labelTok := p.advance() // label identifier
	name, err := p.identifierText(labelTok)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(':')); err != nil {
		return nil, err
	}

	stmt := &ast.LabeledStmt{Label: name, Loc: labelTok.Loc}

	var inner ast.Node
	switch p.cur.Type {
	case kwBEGIN:
		inner, err = p.parseBeginEndBlock()
	case kwWHILE:
		inner, err = p.parseWhileStmt()
	case kwLOOP:
		inner, err = p.parseLoopStmt()
	case kwREPEAT:
		inner, err = p.parseRepeatStmt()
	case kwFOR:
		inner, err = p.parseForInStmt()
	default:
		// Only the unlabeled block forms can carry a label.
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	stmt.Stmt = inner
	stmt.Loc.End = p.prev.Loc.End

	// Optional trailing end-label (identifier?). It must not be consumed if it is
	// actually the next statement — but at this point the block has fully closed,
	// and inside a statement_list the next token after a statement is ';' or a
	// terminator, never a bare identifier. So a bare identifier here is the
	// end-label.
	if isIdentifierStart(p.cur.Type) {
		endTok := p.advance()
		endName, err := p.identifierText(endTok)
		if err != nil {
			return nil, err
		}
		stmt.EndLabel = endName
		stmt.Loc.End = endTok.Loc.End
	}
	return stmt, nil
}

// isScriptLabelStart reports whether the current position begins a labeled
// script statement: an identifier immediately followed by ':'. Used by parseStmt
// to route a `label: <block>` before the ordinary keyword dispatch.
func (p *Parser) isScriptLabelStart() bool {
	return isIdentifierStart(p.cur.Type) && p.peekNext().Type == int(':')
}
