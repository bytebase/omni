package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Snowflake Scripting (T7.1) — STRUCTURAL parser
//
// A Snowflake Scripting block is a structured procedural body:
//
//	[ DECLARE <declaration>; ... ] BEGIN <statement>; ... [ EXCEPTION <handler> ... ] END [ <label> ]
//
// The body is parsed STRUCTURALLY (control flow, declarations, cursors,
// exceptions) but NOT semantically — expressions and nested SQL reuse the
// engine's existing parseExpr / parseStmt. The grammar follows the Snowflake
// Scripting reference (docs are authoritative); the legacy SnowflakeParser.g4
// modeled only a thin subset (BEGIN..END, a DECLARE list of `id data_type`,
// `id := expr`, RETURN expr), so this node deliberately exceeds it (recorded as
// a divergence: omni now parses the full scripting surface, legacy ANTLR did
// not).
//
// ---------------------------------------------------------------------------
// Non-termination safety (this is the highest loop-risk node)
//
// Scripting bodies nest arbitrarily (blocks within IF within FOR within ...).
// Two structural guards make an infinite loop impossible:
//
//  1. Depth bound. scriptDepth is incremented on entry to every nested
//     construct and checked against maxScriptDepth; exceeding it returns a
//     ParseError rather than recursing further. This caps stack growth.
//
//  2. Per-iteration progress. Every statement-list / declaration-list loop
//     records p.cur.Loc.Start before parsing an item and, on a successful
//     parse that did NOT advance the token (a structurally-impossible-to-make
//     case, guarded defensively), returns a ParseError. Each loop also treats
//     tokEOF as an unconditional terminator: a list that reaches EOF before its
//     closing keyword (END / END IF / END FOR / ...) yields a "unterminated"
//     ParseError. Together these guarantee every loop both advances and is
//     bounded, so a malformed (e.g. never-closed) block fails fast instead of
//     hanging.
// ---------------------------------------------------------------------------

// maxScriptDepth bounds Snowflake Scripting nesting. Real scripts nest only a
// handful of levels deep; this generous ceiling exists solely to fail fast on
// adversarial / corrupt input rather than recurse until the Go stack overflows.
const maxScriptDepth = 200

// errScriptTooDeep is returned when scriptDepth exceeds maxScriptDepth.
func (p *Parser) errScriptTooDeep() *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: "snowflake scripting block nested too deeply"}
}

// enterScript increments the nesting depth and reports whether the new depth is
// within bounds. Callers must pair a true result with a deferred leaveScript.
func (p *Parser) enterScript() bool {
	p.scriptDepth++
	return p.scriptDepth <= maxScriptDepth
}

func (p *Parser) leaveScript() { p.scriptDepth-- }

// beginIsTransaction reports whether a BEGIN at the current position opens a
// TCL transaction (BEGIN [WORK|TRANSACTION] [NAME ...]) rather than a Snowflake
// Scripting block. cur must be the BEGIN keyword. The decision is made on the
// token AFTER BEGIN: WORK / TRANSACTION / NAME, or a statement boundary
// (`;` / EOF), means TCL; anything else (a statement opener) means a scripting
// block. This keeps every documented TCL BEGIN form working while routing a
// `BEGIN <statement> ... END` body to the scripting parser.
func (p *Parser) beginIsTransaction() bool {
	switch p.peekNext().Type {
	case kwWORK, kwTRANSACTION, kwNAME, ';', tokEOF:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Block
// ---------------------------------------------------------------------------

// parseScriptBlock parses a Snowflake Scripting block. cur is either DECLARE or
// BEGIN (already known to open a block, not a TCL BEGIN):
//
//	[ DECLARE <declaration>; ... ] BEGIN <statement>; ... [ EXCEPTION <handler> ... ] END [ <label> ]
func (p *Parser) parseScriptBlock() (ast.Node, error) {
	if !p.enterScript() {
		return nil, p.errScriptTooDeep()
	}
	defer p.leaveScript()

	start := p.cur.Loc.Start
	stmt := &ast.ScriptBlockStmt{Loc: ast.Loc{Start: start}}

	// Optional DECLARE section.
	if p.cur.Type == kwDECLARE {
		p.advance() // consume DECLARE
		decls, err := p.parseScriptDeclarationList()
		if err != nil {
			return nil, err
		}
		stmt.Decls = decls
	}

	// BEGIN <statement_list>.
	if _, err := p.expect(kwBEGIN); err != nil {
		return nil, err
	}
	body, err := p.parseScriptStatementList(kwEND, kwSCRIPT_EXCEPTION)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	// Optional EXCEPTION section.
	if p.cur.Type == kwSCRIPT_EXCEPTION {
		p.advance() // consume EXCEPTION
		handlers, err := p.parseScriptExceptionHandlers()
		if err != nil {
			return nil, err
		}
		stmt.Handlers = handlers
	}

	// END [ <label> ].
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// tryParseScriptLabel consumes an optional trailing block/loop label after an
// END (or END FOR / END LOOP / ...) and reports whether one was present.
//
// A label is accepted ONLY as a bare or quoted identifier (tokIdent /
// tokQuotedIdent) — NOT as a non-reserved keyword. This is deliberate: keyword
// tokens are not valid labels (a user label that happens to collide with a
// keyword must be quoted), and accepting them would let a structural keyword be
// swallowed as a label. In particular, in a nested block `BEGIN BEGIN ... END
// END` (no `;` between the inner and outer END) the inner block's
// tryParseScriptLabel must NOT consume the outer END (a non-reserved keyword)
// as a label, which would leave the outer block unterminated.
func (p *Parser) tryParseScriptLabel() (ast.Ident, bool) {
	switch p.cur.Type {
	case tokIdent, tokQuotedIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: tok.Type == tokQuotedIdent, Loc: tok.Loc}, true
	default:
		return ast.Ident{}, false
	}
}

// ---------------------------------------------------------------------------
// DECLARE section
// ---------------------------------------------------------------------------

// parseScriptDeclarationList parses one or more declarations, each terminated
// by a semicolon, stopping at BEGIN:
//
//	<declaration>; [ <declaration>; ... ] BEGIN
//
// Per the docs every declaration (including the last) is `;`-terminated.
func (p *Parser) parseScriptDeclarationList() ([]ast.Node, error) {
	var decls []ast.Node
	for p.cur.Type != kwBEGIN {
		if p.cur.Type == tokEOF {
			return nil, p.unterminated("DECLARE section (expected BEGIN)")
		}
		startPos := p.cur.Loc.Start

		decl, err := p.parseScriptDeclaration()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)

		// Each declaration must be terminated by ';'.
		if _, err := p.expect(';'); err != nil {
			return nil, err
		}

		if p.cur.Loc.Start == startPos {
			// Defensive: a declaration that consumed nothing would loop forever.
			return nil, p.noProgress("DECLARE section")
		}
	}
	return decls, nil
}

// parseScriptDeclaration parses one DECLARE entry. The form is selected by the
// token that follows the declared name:
//
//	<name> CURSOR FOR <query>
//	<name> RESULTSET [ { DEFAULT | := } [ ASYNC ] ( <query> ) ]
//	<name> EXCEPTION [ ( <number>, '<message>' ) ]
//	<name> [ <type> ] [ { DEFAULT | := } <expr> ]      (variable; type optional)
func (p *Parser) parseScriptDeclaration() (*ast.ScriptDeclaration, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	decl := &ast.ScriptDeclaration{Name: name, Loc: ast.Loc{Start: name.Loc.Start, End: name.Loc.End}}

	switch p.cur.Type {
	case kwSCRIPT_CURSOR:
		p.advance() // consume CURSOR
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		q, err := p.parseScriptQuery()
		if err != nil {
			return nil, err
		}
		decl.Kind = ast.ScriptDeclCursor
		decl.Query = q
		decl.Loc.End = ast.NodeLoc(q).End

	case kwSCRIPT_RESULTSET:
		p.advance() // consume RESULTSET
		decl.Kind = ast.ScriptDeclResultset
		// Optional { DEFAULT | := } [ ASYNC ] ( <query> ).
		if p.tryParseAssignOp() {
			async, q, end, err := p.parseResultsetValue()
			if err != nil {
				return nil, err
			}
			decl.Async = async
			decl.Query = q
			decl.Loc.End = end
		} else {
			decl.Loc.End = p.prev.Loc.End
		}

	case kwSCRIPT_EXCEPTION:
		p.advance() // consume EXCEPTION
		decl.Kind = ast.ScriptDeclException
		// Optional ( <number>, '<message>' ).
		if p.cur.Type == '(' {
			args, end, err := p.parseScriptParenExprList()
			if err != nil {
				return nil, err
			}
			decl.ExcArgs = args
			decl.Loc.End = end
		} else {
			decl.Loc.End = p.prev.Loc.End
		}

	default:
		// Variable: <name> [ <type> ] [ { DEFAULT | := } <expr> ]. The type is
		// optional only when a value follows (inferred); a bare `<name>;` with
		// neither type nor value is not a valid declaration, but the assignment
		// operator below is checked first so `<name> := <expr>` works typeless.
		decl.Kind = ast.ScriptDeclVar
		if !p.startsAssignOp() {
			typ, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			decl.Type = typ
			decl.Loc.End = typ.Loc.End
		}
		if p.tryParseAssignOp() {
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			decl.Default = val
			decl.Loc.End = ast.NodeLoc(val).End
		}
	}

	return decl, nil
}

// parseResultsetValue parses the value tail of a RESULTSET, after the
// `{ DEFAULT | := }` operator has been consumed:
//
//	[ ASYNC ] ( <query> )
//
// Returns whether ASYNC was present, the query node, and the end offset.
func (p *Parser) parseResultsetValue() (async bool, query ast.Node, end int, err error) {
	if p.cur.Type == kwSCRIPT_ASYNC {
		p.advance() // consume ASYNC
		async = true
	}
	if _, err := p.expect('('); err != nil {
		return false, nil, 0, err
	}
	q, err := p.parseScriptQuery()
	if err != nil {
		return false, nil, 0, err
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return false, nil, 0, err
	}
	return async, q, closeTok.Loc.End, nil
}

// ---------------------------------------------------------------------------
// EXCEPTION section
// ---------------------------------------------------------------------------

// parseScriptExceptionHandlers parses one or more WHEN handlers, stopping at
// END:
//
//	WHEN <exc_name> [ OR <exc_name> ... ] THEN <stmts>
//	WHEN OTHER THEN <stmts>
func (p *Parser) parseScriptExceptionHandlers() ([]ast.Node, error) {
	var handlers []ast.Node
	for p.cur.Type == kwWHEN {
		startPos := p.cur.Loc.Start
		whenTok := p.advance() // consume WHEN

		h := &ast.ScriptExceptionHandler{Loc: ast.Loc{Start: whenTok.Loc.Start}}

		if p.cur.Type == kwSCRIPT_OTHER {
			p.advance() // consume OTHER
			h.Other = true
		} else {
			// <exc_name> [ OR <exc_name> ... ].
			for {
				n, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				h.Names = append(h.Names, n)
				if p.cur.Type != kwOR {
					break
				}
				p.advance() // consume OR
			}
		}

		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		body, err := p.parseScriptStatementList(kwWHEN, kwEND)
		if err != nil {
			return nil, err
		}
		h.Body = body
		h.Loc.End = p.prev.Loc.End
		handlers = append(handlers, h)

		if p.cur.Loc.Start == startPos {
			return nil, p.noProgress("EXCEPTION section")
		}
	}
	if len(handlers) == 0 {
		// EXCEPTION with no WHEN handler is malformed.
		return nil, p.syntaxErrorAtCur()
	}
	return handlers, nil
}

// ---------------------------------------------------------------------------
// Statement list
// ---------------------------------------------------------------------------

// parseScriptStatementList parses a `;`-separated run of scripting statements,
// stopping (without consuming) at any of the given terminator keywords or EOF.
//
// Snowflake Scripting requires each statement to be `;`-terminated; this parser
// also tolerates a missing final `;` before the terminator (lenient, matching
// the legacy task_scripting_statement_list `SEMI?`). EOF before a terminator is
// an "unterminated" error — that is the structural guarantee against a runaway
// loop on a never-closed block.
func (p *Parser) parseScriptStatementList(terminators ...int) ([]ast.Node, error) {
	var stmts []ast.Node
	for {
		if p.cur.Type == tokEOF {
			return nil, p.unterminated("scripting block")
		}
		if p.isScriptTerminator(p.cur.Type, terminators) {
			return stmts, nil
		}

		startPos := p.cur.Loc.Start
		s, err := p.parseScriptStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)

		if p.cur.Loc.Start == startPos {
			// Defensive: a statement that consumed nothing would loop forever.
			return nil, p.noProgress("scripting block")
		}

		// Consume the statement separator. A statement may be immediately
		// followed by a terminator with no trailing ';' (lenient).
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		if p.isScriptTerminator(p.cur.Type, terminators) {
			return stmts, nil
		}
		if p.cur.Type == tokEOF {
			return nil, p.unterminated("scripting block")
		}
		// Anything else after a complete statement (and not a separator /
		// terminator) is a syntax error rather than a silent stop.
		return nil, p.syntaxErrorAtCur()
	}
}

// isScriptTerminator reports whether tokType is one of terminators.
func (p *Parser) isScriptTerminator(tokType int, terminators []int) bool {
	for _, t := range terminators {
		if tokType == t {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Statement dispatch
// ---------------------------------------------------------------------------

// parseScriptStatement parses ONE scripting statement and guarantees it
// consumes at least one token. It dispatches on the leading token:
//
//	nested block ...................... DECLARE / BEGIN
//	IF / CASE / FOR / WHILE / REPEAT / LOOP ... control flow
//	BREAK / CONTINUE / RETURN / RAISE ... jump statements
//	LET ............................... declaration-statement
//	OPEN / FETCH / CLOSE .............. cursor statements
//	<var> := <expr> ................... assignment
//	<sql> ............................. any SQL statement (reuses parseStmt)
func (p *Parser) parseScriptStatement() (ast.Node, error) {
	if !p.enterScript() {
		return nil, p.errScriptTooDeep()
	}
	defer p.leaveScript()

	switch p.cur.Type {
	case kwDECLARE:
		return p.parseScriptBlock()
	case kwBEGIN:
		// A bare BEGIN inside a block always opens a nested scripting block
		// (transaction control is not a scripting statement). Route directly.
		return p.parseScriptBlock()
	case kwIF:
		return p.parseScriptIf()
	case kwCASE:
		return p.parseScriptCase()
	case kwFOR:
		return p.parseScriptFor()
	case kwSCRIPT_WHILE:
		return p.parseScriptWhile()
	case kwSCRIPT_REPEAT:
		return p.parseScriptRepeat()
	case kwSCRIPT_LOOP:
		return p.parseScriptLoop()
	case kwSCRIPT_BREAK, kwSCRIPT_EXIT:
		return p.parseScriptBreak()
	case kwCONTINUE, kwSCRIPT_ITERATE:
		return p.parseScriptContinue()
	case kwRETURN:
		return p.parseScriptReturn()
	case kwSCRIPT_RAISE:
		return p.parseScriptRaise()
	case kwSCRIPT_LET:
		return p.parseScriptLet()
	case kwSCRIPT_OPEN:
		return p.parseScriptOpen()
	case kwFETCH:
		return p.parseScriptFetch()
	case kwSCRIPT_CLOSE:
		return p.parseScriptClose()
	default:
		// Either an assignment (<var> := <expr>) or a nested SQL statement.
		// An assignment is detected by an identifier-like leading token whose
		// successor is the `:=` operator; everything else is delegated to the
		// top-level statement parser (SELECT / INSERT / CALL / EXECUTE / ...).
		if p.isIdentToken() && p.assignFollows() {
			return p.parseScriptAssign()
		}
		return p.parseStmt()
	}
}

// assignFollows reports whether the token AFTER the current identifier-like
// token begins the `:=` assignment operator (a ':' followed by '='). cur is the
// candidate target identifier.
func (p *Parser) assignFollows() bool {
	return p.peekNext().Type == ':'
}

// parseScriptAssign parses `<var> := <expr>`. cur is the target identifier.
func (p *Parser) parseScriptAssign() (ast.Node, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if !p.tryParseAssignOp() {
		return nil, p.syntaxErrorAtCur()
	}
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.ScriptAssignStmt{
		Target: name,
		Value:  val,
		Loc:    ast.Loc{Start: name.Loc.Start, End: ast.NodeLoc(val).End},
	}, nil
}

// ---------------------------------------------------------------------------
// LET
// ---------------------------------------------------------------------------

// parseScriptLet parses a LET declaration-statement. cur is LET:
//
//	LET <var> [ <type> ] { DEFAULT | := } <expr>
//	LET <cursor> CURSOR FOR <query>
//	LET <resultset> RESULTSET { DEFAULT | := } [ ASYNC ] ( <query> )
func (p *Parser) parseScriptLet() (ast.Node, error) {
	letTok := p.advance() // consume LET
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt := &ast.ScriptLetStmt{Name: name, Loc: ast.Loc{Start: letTok.Loc.Start}}

	switch p.cur.Type {
	case kwSCRIPT_CURSOR:
		p.advance() // consume CURSOR
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		q, err := p.parseScriptQuery()
		if err != nil {
			return nil, err
		}
		stmt.Kind = ast.ScriptDeclCursor
		stmt.Query = q
		stmt.Loc.End = ast.NodeLoc(q).End

	case kwSCRIPT_RESULTSET:
		p.advance() // consume RESULTSET
		stmt.Kind = ast.ScriptDeclResultset
		if !p.tryParseAssignOp() {
			return nil, p.syntaxErrorAtCur()
		}
		async, q, end, err := p.parseResultsetValue()
		if err != nil {
			return nil, err
		}
		stmt.Async = async
		stmt.Query = q
		stmt.Loc.End = end

	default:
		// Variable: LET <var> [ <type> ] { DEFAULT | := } <expr>.
		stmt.Kind = ast.ScriptDeclVar
		if !p.startsAssignOp() {
			typ, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.Type = typ
		}
		if !p.tryParseAssignOp() {
			return nil, p.syntaxErrorAtCur()
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Default = val
		stmt.Loc.End = ast.NodeLoc(val).End
	}

	return stmt, nil
}

// ---------------------------------------------------------------------------
// IF
// ---------------------------------------------------------------------------

// parseScriptIf parses an IF statement. cur is IF:
//
//	IF ( <cond> ) THEN <stmts> [ ELSEIF ( <cond> ) THEN <stmts> ]* [ ELSE <stmts> ] END IF
func (p *Parser) parseScriptIf() (ast.Node, error) {
	ifTok := p.advance() // consume IF
	stmt := &ast.ScriptIfStmt{Loc: ast.Loc{Start: ifTok.Loc.Start}}

	// Leading IF arm.
	branch, err := p.parseScriptIfBranch()
	if err != nil {
		return nil, err
	}
	stmt.Branches = append(stmt.Branches, branch)

	// Zero or more ELSEIF arms.
	for p.cur.Type == kwSCRIPT_ELSEIF {
		p.advance() // consume ELSEIF
		b, err := p.parseScriptIfBranch()
		if err != nil {
			return nil, err
		}
		stmt.Branches = append(stmt.Branches, b)
	}

	// Optional ELSE.
	if p.cur.Type == kwELSE {
		p.advance() // consume ELSE
		body, err := p.parseScriptStatementList(kwEND)
		if err != nil {
			return nil, err
		}
		stmt.Else = body
	}

	// END IF.
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIF); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseScriptIfBranch parses `( <cond> ) THEN <stmts>` for one IF / ELSEIF arm.
// The body terminates at ELSEIF / ELSE / END (whichever comes first).
func (p *Parser) parseScriptIfBranch() (*ast.ScriptIfBranch, error) {
	startPos := p.cur.Loc.Start
	cond, err := p.parseScriptParenCondition()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}
	body, err := p.parseScriptStatementList(kwSCRIPT_ELSEIF, kwELSE, kwEND)
	if err != nil {
		return nil, err
	}
	return &ast.ScriptIfBranch{
		Cond: cond,
		Body: body,
		Loc:  ast.Loc{Start: startPos, End: p.prev.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// CASE
// ---------------------------------------------------------------------------

// parseScriptCase parses a CASE statement (simple or searched). cur is CASE:
//
//	CASE [ ( <operand> ) ] WHEN <expr> THEN <stmts> [ WHEN ... ]* [ ELSE <stmts> ] END [ CASE ]
//
// The simple form has an operand expression after CASE (commonly parenthesized,
// e.g. `CASE (x)`); the searched form goes straight to WHEN. The operand is
// detected as "anything that is not the WHEN keyword".
func (p *Parser) parseScriptCase() (ast.Node, error) {
	caseTok := p.advance() // consume CASE
	stmt := &ast.ScriptCaseStmt{Loc: ast.Loc{Start: caseTok.Loc.Start}}

	// Optional operand (simple CASE). Searched CASE begins directly with WHEN.
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Operand = operand
	}

	// One or more WHEN <expr> THEN <stmts> arms.
	for p.cur.Type == kwWHEN {
		startPos := p.cur.Loc.Start
		p.advance() // consume WHEN
		match, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		body, err := p.parseScriptStatementList(kwWHEN, kwELSE, kwEND)
		if err != nil {
			return nil, err
		}
		stmt.Whens = append(stmt.Whens, &ast.ScriptCaseWhen{
			Match: match,
			Body:  body,
			Loc:   ast.Loc{Start: startPos, End: p.prev.Loc.End},
		})
	}
	if len(stmt.Whens) == 0 {
		// CASE with no WHEN arm is malformed.
		return nil, p.syntaxErrorAtCur()
	}

	// Optional ELSE.
	if p.cur.Type == kwELSE {
		p.advance() // consume ELSE
		body, err := p.parseScriptStatementList(kwEND)
		if err != nil {
			return nil, err
		}
		stmt.Else = body
	}

	// END [ CASE ].
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if p.cur.Type == kwCASE {
		p.advance() // consume optional CASE
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// FOR
// ---------------------------------------------------------------------------

// parseScriptFor parses a FOR loop in either form. cur is FOR:
//
//	FOR <counter> IN [ REVERSE ] <start> TO <end> { DO | LOOP } <stmts> END { FOR | LOOP } [ <label> ]
//	FOR <row> IN { <cursor> | <resultset> } DO <stmts> END FOR [ <label> ]
//
// The two forms are disambiguated AFTER `IN`: a REVERSE modifier or a `<start>
// TO <end>` range marks the counter form; a bare cursor/resultset name (i.e.
// the loop body opener `DO`/`LOOP` follows the name) marks the cursor form.
func (p *Parser) parseScriptFor() (ast.Node, error) {
	forTok := p.advance() // consume FOR
	stmt := &ast.ScriptForStmt{Loc: ast.Loc{Start: forTok.Loc.Start}}

	v, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Var = v

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}

	// REVERSE forces the counter form.
	if p.cur.Type == kwSCRIPT_REVERSE {
		p.advance() // consume REVERSE
		stmt.Kind = ast.ScriptForCounter
		stmt.Reverse = true
		if err := p.parseScriptForCounterRange(stmt); err != nil {
			return nil, err
		}
	} else {
		// Disambiguate counter vs cursor: a cursor/resultset source is a single
		// identifier immediately followed by the loop opener DO / LOOP. Anything
		// else is a counter range expression `<start> TO <end>`.
		if p.isIdentToken() && p.scriptForSourceFollows() {
			src, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.Kind = ast.ScriptForCursor
			stmt.Source = src
		} else {
			stmt.Kind = ast.ScriptForCounter
			if err := p.parseScriptForCounterRange(stmt); err != nil {
				return nil, err
			}
		}
	}

	// Loop opener: DO or LOOP.
	if _, err := p.expectLoopOpener(); err != nil {
		return nil, err
	}

	// Body terminates at END.
	body, err := p.parseScriptStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	// END { FOR | LOOP } [ <label> ].
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expectForLoopCloser(); err != nil {
		return nil, err
	}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// scriptForSourceFollows reports whether the current identifier is a bare
// FOR-loop cursor/resultset source — i.e. the token AFTER it is a loop opener
// (DO / LOOP). A counter range instead has `TO` (or a larger expression) after
// the start value, so the loop opener does not immediately follow.
func (p *Parser) scriptForSourceFollows() bool {
	switch p.peekNext().Type {
	case kwDO, kwSCRIPT_LOOP:
		return true
	default:
		return false
	}
}

// parseScriptForCounterRange parses `<start> TO <end>` into the counter-form
// stmt fields.
func (p *Parser) parseScriptForCounterRange(stmt *ast.ScriptForStmt) error {
	start, err := p.parseExpr()
	if err != nil {
		return err
	}
	stmt.Start = start
	if _, err := p.expect(kwTO); err != nil {
		return err
	}
	end, err := p.parseExpr()
	if err != nil {
		return err
	}
	stmt.End = end
	return nil
}

// ---------------------------------------------------------------------------
// WHILE / REPEAT / LOOP
// ---------------------------------------------------------------------------

// parseScriptWhile parses `WHILE ( <cond> ) { DO | LOOP } <stmts> END { WHILE | LOOP } [ <label> ]`.
func (p *Parser) parseScriptWhile() (ast.Node, error) {
	whileTok := p.advance() // consume WHILE
	stmt := &ast.ScriptWhileStmt{Loc: ast.Loc{Start: whileTok.Loc.Start}}

	cond, err := p.parseScriptParenCondition()
	if err != nil {
		return nil, err
	}
	stmt.Cond = cond

	if _, err := p.expectLoopOpener(); err != nil {
		return nil, err
	}
	body, err := p.parseScriptStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expectWhileLoopCloser(); err != nil {
		return nil, err
	}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseScriptRepeat parses `REPEAT <stmts> UNTIL ( <cond> ) END REPEAT [ <label> ]`.
func (p *Parser) parseScriptRepeat() (ast.Node, error) {
	repeatTok := p.advance() // consume REPEAT
	stmt := &ast.ScriptRepeatStmt{Loc: ast.Loc{Start: repeatTok.Loc.Start}}

	body, err := p.parseScriptStatementList(kwSCRIPT_UNTIL)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwSCRIPT_UNTIL); err != nil {
		return nil, err
	}
	cond, err := p.parseScriptParenCondition()
	if err != nil {
		return nil, err
	}
	stmt.Cond = cond

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSCRIPT_REPEAT); err != nil {
		return nil, err
	}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseScriptLoop parses `LOOP <stmts> END LOOP [ <label> ]`.
func (p *Parser) parseScriptLoop() (ast.Node, error) {
	loopTok := p.advance() // consume LOOP
	stmt := &ast.ScriptLoopStmt{Loc: ast.Loc{Start: loopTok.Loc.Start}}

	body, err := p.parseScriptStatementList(kwEND)
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSCRIPT_LOOP); err != nil {
		return nil, err
	}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Jump statements: BREAK / CONTINUE / RETURN / RAISE
// ---------------------------------------------------------------------------

// parseScriptBreak parses `BREAK [ <label> ]` (alias EXIT). cur is BREAK/EXIT.
func (p *Parser) parseScriptBreak() (ast.Node, error) {
	tok := p.advance() // consume BREAK / EXIT
	stmt := &ast.ScriptBreakStmt{Loc: tok.Loc}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
		stmt.Loc.End = label.Loc.End
	}
	return stmt, nil
}

// parseScriptContinue parses `CONTINUE [ <label> ]` (alias ITERATE). cur is
// CONTINUE/ITERATE.
func (p *Parser) parseScriptContinue() (ast.Node, error) {
	tok := p.advance() // consume CONTINUE / ITERATE
	stmt := &ast.ScriptContinueStmt{Loc: tok.Loc}
	if label, ok := p.tryParseScriptLabel(); ok {
		stmt.Label = label
		stmt.Loc.End = label.Loc.End
	}
	return stmt, nil
}

// parseScriptReturn parses `RETURN [ <expr> ]`. cur is RETURN. A bare RETURN
// (immediately followed by a statement separator / terminator) has no value.
func (p *Parser) parseScriptReturn() (ast.Node, error) {
	tok := p.advance() // consume RETURN
	stmt := &ast.ScriptReturnStmt{Loc: tok.Loc}
	if p.scriptExprEnds() {
		return stmt, nil
	}
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Value = val
	stmt.Loc.End = ast.NodeLoc(val).End
	return stmt, nil
}

// parseScriptRaise parses `RAISE [ <exc_name> ]`. cur is RAISE. A bare RAISE
// re-raises the current exception (valid only inside a handler, but parsed
// structurally regardless).
func (p *Parser) parseScriptRaise() (ast.Node, error) {
	tok := p.advance() // consume RAISE
	stmt := &ast.ScriptRaiseStmt{Loc: tok.Loc}
	if p.scriptExprEnds() {
		return stmt, nil
	}
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = name.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Cursor statements: OPEN / FETCH / CLOSE
// ---------------------------------------------------------------------------

// parseScriptOpen parses `OPEN <cursor> [ USING ( <bind> [ , ... ] ) ]`. cur is OPEN.
func (p *Parser) parseScriptOpen() (ast.Node, error) {
	openTok := p.advance() // consume OPEN
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt := &ast.ScriptOpenStmt{Cursor: name, Loc: ast.Loc{Start: openTok.Loc.Start, End: name.Loc.End}}

	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		args, end, err := p.parseScriptParenExprList()
		if err != nil {
			return nil, err
		}
		stmt.Using = args
		stmt.Loc.End = end
	}
	return stmt, nil
}

// parseScriptFetch parses `FETCH <cursor> INTO <var> [ , ... ]`. cur is FETCH.
func (p *Parser) parseScriptFetch() (ast.Node, error) {
	fetchTok := p.advance() // consume FETCH
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt := &ast.ScriptFetchStmt{Cursor: name, Loc: ast.Loc{Start: fetchTok.Loc.Start, End: name.Loc.End}}

	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	for {
		v, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Into = append(stmt.Into, v)
		stmt.Loc.End = v.Loc.End
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return stmt, nil
}

// parseScriptClose parses `CLOSE <cursor>`. cur is CLOSE.
func (p *Parser) parseScriptClose() (ast.Node, error) {
	closeTok := p.advance() // consume CLOSE
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &ast.ScriptCloseStmt{
		Cursor: name,
		Loc:    ast.Loc{Start: closeTok.Loc.Start, End: name.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// parseScriptParenCondition parses `( <cond> )` — the parenthesized boolean
// condition shared by IF / ELSEIF / WHILE / UNTIL. The parentheses are required
// by the docs. The returned node is the inner condition expression (the
// surrounding ParenExpr is preserved so the source round-trips).
func (p *Parser) parseScriptParenCondition() (ast.Node, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	// parseExpr handles a leading '(' as a ParenExpr, keeping the parens in the
	// AST. This is sufficient for structural fidelity.
	return p.parseExpr()
}

// parseScriptQuery parses a nested query (the body of a CURSOR FOR / RESULTSET).
// It accepts SELECT, WITH ... SELECT, and a parenthesized query, reusing the
// engine's query-expression parser.
func (p *Parser) parseScriptQuery() (ast.Node, error) {
	switch p.cur.Type {
	case kwWITH:
		return p.parseWithQueryExpr()
	default:
		return p.parseQueryExpr()
	}
}

// parseScriptParenExprList parses `( <expr> [ , <expr> ... ] )` and returns the
// expressions plus the end offset of the closing ')'. An empty list `()` is
// permitted (yields nil). cur must be '('.
func (p *Parser) parseScriptParenExprList() ([]ast.Node, int, error) {
	if _, err := p.expect('('); err != nil {
		return nil, 0, err
	}
	if p.cur.Type == ')' {
		closeTok := p.advance()
		return nil, closeTok.Loc.End, nil
	}
	var args []ast.Node
	for {
		e, err := p.parseExpr()
		if err != nil {
			return nil, 0, err
		}
		args = append(args, e)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return nil, 0, err
	}
	return args, closeTok.Loc.End, nil
}

// startsAssignOp reports whether cur begins an assignment operator: either the
// `:=` operator (a ':' followed by '=') or the DEFAULT keyword. Used to decide
// whether a declaration / LET has reached its value (so the optional type slot
// is skipped).
func (p *Parser) startsAssignOp() bool {
	if p.cur.Type == kwDEFAULT {
		return true
	}
	return p.cur.Type == ':' && p.peekNext().Type == '='
}

// tryParseAssignOp consumes an assignment operator — `:=` or the DEFAULT keyword
// — if present, and reports whether one was consumed.
func (p *Parser) tryParseAssignOp() bool {
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		return true
	}
	if p.cur.Type == ':' && p.peekNext().Type == '=' {
		p.advance() // consume ':'
		p.advance() // consume '='
		return true
	}
	return false
}

// scriptExprEnds reports whether the current token ends an optional-expression
// statement (RETURN / RAISE) without a value — i.e. a statement separator,
// terminator keyword, or EOF follows immediately.
func (p *Parser) scriptExprEnds() bool {
	switch p.cur.Type {
	case ';', tokEOF,
		kwEND, kwELSE, kwSCRIPT_ELSEIF, kwWHEN, kwSCRIPT_UNTIL:
		return true
	default:
		return false
	}
}

// expectLoopOpener consumes the loop-body opener DO or LOOP (shared by FOR /
// WHILE). Returns a syntax error if neither is present.
func (p *Parser) expectLoopOpener() (Token, error) {
	switch p.cur.Type {
	case kwDO, kwSCRIPT_LOOP:
		return p.advance(), nil
	default:
		return Token{}, p.syntaxErrorAtCur()
	}
}

// expectForLoopCloser consumes the FOR-loop closer FOR or LOOP after END.
func (p *Parser) expectForLoopCloser() (Token, error) {
	switch p.cur.Type {
	case kwFOR, kwSCRIPT_LOOP:
		return p.advance(), nil
	default:
		return Token{}, p.syntaxErrorAtCur()
	}
}

// expectWhileLoopCloser consumes the WHILE-loop closer WHILE or LOOP after END.
func (p *Parser) expectWhileLoopCloser() (Token, error) {
	switch p.cur.Type {
	case kwSCRIPT_WHILE, kwSCRIPT_LOOP:
		return p.advance(), nil
	default:
		return Token{}, p.syntaxErrorAtCur()
	}
}

// unterminated returns a ParseError for a scripting construct that hit EOF
// before its closing keyword. This is the fast-fail path that prevents a
// never-closed block from looping.
func (p *Parser) unterminated(what string) *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: "unterminated " + what + " at end of input"}
}

// noProgress returns a ParseError for the defensive case where a list loop
// parsed an item without consuming any token. It can only fire on a parser bug;
// returning an error (rather than spinning) keeps non-termination structurally
// impossible.
func (p *Parser) noProgress(what string) *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: "internal: no progress parsing " + what}
}
