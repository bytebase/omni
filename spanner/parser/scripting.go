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
// the live Cloud Spanner emulator (the suite:scripting rows of testdata/compat.tsv):
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
	"github.com/bytebase/omni/spanner/ast"
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

// ---------------------------------------------------------------------------
// BEGIN … END block
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// IF
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// CASE (statement form)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WHILE / LOOP / REPEAT
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// FOR … IN
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// DECLARE
// ---------------------------------------------------------------------------

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
				Position: idTok.Loc.Start, End: idTok.Loc.End,
				Message: "Using SET with multiple variables requires parentheses around the variable list",
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
	// subqueries; the statement dispatch (parseMaybePipeSuffixStmt) fills them
	// for the query-span walker.
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
	return nil, p.syntaxErrorAtCur()
}

// isScriptLabelStart reports whether the current position begins a labeled
// script statement: an identifier immediately followed by ':'. Used by parseStmt
// to route a `label: <block>` before the ordinary keyword dispatch.
func (p *Parser) isScriptLabelStart() bool {
	return isIdentifierStart(p.cur.Type) && p.peekNext().Type == int(':')
}
