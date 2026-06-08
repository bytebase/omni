package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// CALL / EXECUTE IMMEDIATE / EXECUTE TASK / EXPLAIN (T5.4)
//
// Triangulated oracle: legacy SnowflakeParser.g4 (truth2) + the official
// Snowflake docs (truth1, which win on conflict) + the official example corpus
// under testdata/official/{call,execute-immediate,explain}.
//
// Grammar:
//
//	CALL <proc_name> ( [ <arg> [ , ... ] ] )
//	  -- legacy: CALL object_name '(' expr_list? ')'.  Args are expressions;
//	  -- named arguments (name => value) are documented and accepted here even
//	  -- though the legacy expr_list rule omits them (docs win).
//
//	EXECUTE IMMEDIATE { '<string>' | $$<body>$$ | <variable> | $<session_var> }
//	                  [ USING ( <bind_var> [ , ... ] ) ]
//	  -- The string / dollar body is OPAQUE and raw-scanned VERBATIM (delimiters
//	  -- included).  The legacy grammar's USING list and the docs' bind-variable
//	  -- list are bare identifiers (matched by the official corpus).
//
//	EXECUTE TASK <name> [ USING CONFIG = <config_string> ]
//	EXECUTE TASK <name> RETRY LAST
//	  -- legacy: EXECUTE TASK object_name retry_last? .  USING CONFIG is docs-only.
//
//	EXPLAIN [ USING { TABULAR | JSON | TEXT } ] <statement>
//	  -- The inner statement is parsed structurally via parseStmt.
//
// EXECUTE is NOT an omni keyword (it is absent from the F2 keyword table), so it
// arrives as a plain identifier and is dispatched from parseStmt's default
// branch by its uppercased text (mirroring LS / RM).  CALL and EXPLAIN have
// dedicated dispatch cases.
//
// $-prefixed CALL arguments (CALL p($v)) depend on session-variable expression
// support owned by a separate node (expr-dollar-refs) and are NOT handled here;
// the EXECUTE IMMEDIATE $<session_var> form below reads the $-variable token
// directly (not through the expression parser) and is therefore self-contained.
// ---------------------------------------------------------------------------

// parseCallStmt parses `CALL <proc_name> ( [ <arg> [ , ... ] ] )`. The CALL
// keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseCallStmt() (ast.Node, error) {
	callTok := p.advance() // consume CALL
	stmt := &ast.CallStmt{
		Loc: ast.Loc{Start: callTok.Loc.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Empty argument list: CALL p().
	if p.cur.Type == ')' {
		closeTok := p.advance() // consume ')'
		stmt.Loc.End = closeTok.Loc.End
		return stmt, nil
	}

	for {
		arg, err := p.parseCallArg()
		if err != nil {
			return nil, err
		}
		stmt.Args = append(stmt.Args, arg)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = closeTok.Loc.End
	return stmt, nil
}

// parseCallArg parses one CALL argument: either a positional expression or a
// named argument `<name> => <expr>`. A named argument is detected by a leading
// identifier immediately followed by the `=>` (tokAssoc) operator.
func (p *Parser) parseCallArg() (*ast.CallArg, error) {
	start := p.cur.Loc.Start

	// Named argument: <name> => <expr>.
	if p.isIdentToken() && p.peekNext().Type == tokAssoc {
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		p.advance() // consume '=>'
		value, err := p.parseCallArgValue()
		if err != nil {
			return nil, err
		}
		return &ast.CallArg{
			Name:  name,
			Value: value,
			Loc:   ast.Loc{Start: start, End: p.prev.Loc.End},
		}, nil
	}

	value, err := p.parseCallArgValue()
	if err != nil {
		return nil, err
	}
	return &ast.CallArg{
		Value: value,
		Loc:   ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// parseCallArgValue parses a single CALL argument value. Arguments are
// expressions; the official docs additionally permit a bare scalar subquery
// (CALL p(SELECT COUNT(*) FROM t)) without surrounding parentheses, so a leading
// SELECT / WITH is parsed as a query expression. (The legacy expr_list rule does
// not allow this; docs win.)
func (p *Parser) parseCallArgValue() (ast.Node, error) {
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		if p.cur.Type == kwWITH {
			return p.parseWithQueryExpr()
		}
		return p.parseQueryExpr()
	}
	return p.parseExpr()
}

// ---------------------------------------------------------------------------
// EXECUTE dispatch (EXECUTE IMMEDIATE | EXECUTE TASK)
// ---------------------------------------------------------------------------

// parseExecuteStmt dispatches an EXECUTE statement on the keyword that follows
// EXECUTE (IMMEDIATE or TASK). The EXECUTE word (a plain identifier, not a
// keyword) has NOT yet been consumed when this function is called.
func (p *Parser) parseExecuteStmt() (ast.Node, error) {
	execTok := p.advance() // consume EXECUTE
	switch p.cur.Type {
	case kwIMMEDIATE:
		return p.parseExecuteImmediateStmt(execTok.Loc)
	case kwTASK:
		return p.parseExecuteTaskStmt(execTok.Loc)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseExecuteImmediateStmt parses
//
//	EXECUTE IMMEDIATE { '<string>' | $$<body>$$ | <variable> | $<session_var> }
//	                  [ USING ( <bind_var> [ , ... ] ) ]
//
// On entry cur is the IMMEDIATE keyword. start is the Loc of the EXECUTE word.
func (p *Parser) parseExecuteImmediateStmt(start ast.Loc) (ast.Node, error) {
	// Snapshot the lexer error count BEFORE consuming IMMEDIATE — consuming it
	// primes cur with the body's first token, and for a multi-line single-quoted
	// body that priming makes scanString append a spurious "unterminated string
	// literal" error that scanExecImmBody must later drop.
	errLenBefore := len(p.lexer.errors)

	p.advance() // consume IMMEDIATE
	stmt := &ast.ExecuteImmediateStmt{
		Loc: ast.Loc{Start: start.Start},
	}

	switch {
	case p.cur.Type == tokVariable:
		// $<session_variable>. The lexer strips the leading $; Str is the name.
		varTok := p.advance()
		stmt.Source = ast.ExecImmSessionVar
		stmt.Var = ast.Ident{Name: varTok.Str, Loc: varTok.Loc}
		stmt.Loc.End = varTok.Loc.End

	case p.curIsStringOrDollarBody():
		// '<string>' or $$<body>$$ — opaque body captured verbatim by raw-scan.
		body, dollar, end, err := p.scanExecImmBody(errLenBefore)
		if err != nil {
			return nil, err
		}
		if dollar {
			stmt.Source = ast.ExecImmDollar
		} else {
			stmt.Source = ast.ExecImmString
		}
		stmt.Body = body
		stmt.Loc.End = end

	case p.isIdentToken():
		// <variable> — a bare Snowflake Scripting local variable.
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Source = ast.ExecImmVariable
		stmt.Var = name
		stmt.Loc.End = name.Loc.End

	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Optional USING ( <bind_variable> [ , ... ] ).
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		using, end, err := p.parseExecImmUsing()
		if err != nil {
			return nil, err
		}
		stmt.Using = using
		stmt.Loc.End = end
	}

	return stmt, nil
}

// curIsStringOrDollarBody reports whether the current token begins an EXECUTE
// IMMEDIATE string ('...') or dollar ($$...$$) body. The lexer emits tokString
// for both a single-line single-quoted string and a complete $$...$$ run; a
// multi-line single-quoted body fails to lex (emitting tokInvalid) but still
// begins with a single quote in the source, so the source byte is consulted as
// well.
func (p *Parser) curIsStringOrDollarBody() bool {
	if p.cur.Type == tokString {
		return true
	}
	i := p.cur.Loc.Start - p.base
	if i < 0 || i >= len(p.input) {
		return false
	}
	return p.input[i] == '\'' || p.input[i] == '$'
}

// scanExecImmBody raw-scans the EXECUTE IMMEDIATE string / dollar body directly
// from the source text and returns it VERBATIM (delimiters included), whether
// the dollar form was used, and the absolute end offset. On entry cur is the
// body's first token.
//
// The body is raw-scanned rather than read from the token stream for the same
// reason routine bodies are (see function_procedure.go): omni's lexer rejects a
// single-quoted string that spans newlines, and the dollar form is scanned the
// same way for uniformity so the body never depends on the token stream.
//
// LOOP-GUARD: each scan loop is bounded by len(p.input) and advances its cursor
// by at least one byte every iteration; the bound and the unconditional j++
// guarantee termination. An unterminated or malformed body returns a ParseError
// rather than looping. A negative test (TestExecuteImmediate_Reject) feeds an
// unterminated body and asserts a fast error return.
//
// errLenBefore is the lexer's error count captured before cur was primed with
// the body token (see parseExecuteImmediateStmt); any spurious lex error the
// discarded body token produced is truncated back to it after a successful scan.
func (p *Parser) scanExecImmBody(errLenBefore int) (body string, dollar bool, end int, err error) {
	startLocal := p.cur.Loc.Start - p.base
	if startLocal < 0 || startLocal >= len(p.input) {
		return "", false, 0, p.syntaxErrorAtCur()
	}

	var endLocal int
	switch p.input[startLocal] {
	case '$':
		// Dollar-quoted body: $$ ... $$ (no nesting, no escapes).
		if startLocal+1 >= len(p.input) || p.input[startLocal+1] != '$' {
			return "", false, 0, p.syntaxErrorAtCur()
		}
		j := startLocal + 2
		closed := false
		for j+1 < len(p.input) {
			if p.input[j] == '$' && p.input[j+1] == '$' {
				j += 2 // include the closing $$
				closed = true
				break
			}
			j++ // LOOP-GUARD: always advances; bounded by len(p.input).
		}
		if !closed {
			return "", false, 0, &ParseError{
				Loc: p.cur.Loc,
				Msg: "unterminated dollar-quoted EXECUTE IMMEDIATE body",
			}
		}
		endLocal = j
		dollar = true

	case '\'':
		// Single-quoted body: ' ... ' where '' is an escaped quote and the body
		// may span newlines. Scan to the matching unescaped closing quote.
		j := startLocal + 1
		closed := false
		for j < len(p.input) {
			if p.input[j] == '\'' {
				if j+1 < len(p.input) && p.input[j+1] == '\'' {
					j += 2 // '' escaped quote; LOOP-GUARD: advances by two.
					continue
				}
				j++ // consume closing quote
				closed = true
				break
			}
			j++ // LOOP-GUARD: always advances; bounded by len(p.input).
		}
		if !closed {
			return "", false, 0, &ParseError{
				Loc: p.cur.Loc,
				Msg: "unterminated EXECUTE IMMEDIATE body string literal",
			}
		}
		endLocal = j
		dollar = false

	default:
		return "", false, 0, p.syntaxErrorAtCur()
	}

	body = p.input[startLocal:endLocal]
	end = endLocal + p.base

	// Drop any spurious lex error a discarded multi-line single-quoted body
	// token produced.
	if len(p.lexer.errors) > errLenBefore {
		p.lexer.errors = p.lexer.errors[:errLenBefore]
	}

	// Re-sync the lexer past the body and re-prime cur, discarding the buffered
	// lookahead the bypassed token stream left behind.
	p.lexer.pos = endLocal
	p.hasNext = false
	p.advance()

	return body, dollar, end, nil
}

// parseExecImmUsing parses the EXECUTE IMMEDIATE `USING ( <bind_variable>
// [ , ... ] )` list. The USING keyword has already been consumed. Returns the
// bind-variable identifiers and the absolute end offset (the closing ')'). The
// list must be non-empty.
func (p *Parser) parseExecImmUsing() ([]ast.Ident, int, error) {
	if _, err := p.expect('('); err != nil {
		return nil, 0, err
	}

	var binds []ast.Ident
	for {
		name, err := p.parseIdent()
		if err != nil {
			return nil, 0, err
		}
		binds = append(binds, name)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, 0, err
	}
	return binds, closeTok.Loc.End, nil
}

// parseExecuteTaskStmt parses
//
//	EXECUTE TASK <name> [ USING CONFIG = <config_string> ]
//	EXECUTE TASK <name> RETRY LAST
//
// On entry cur is the TASK keyword. start is the Loc of the EXECUTE word.
func (p *Parser) parseExecuteTaskStmt(start ast.Loc) (ast.Node, error) {
	p.advance() // consume TASK
	stmt := &ast.ExecuteTaskStmt{
		Loc: ast.Loc{Start: start.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.cur.Type == kwRETRY:
		p.advance() // consume RETRY
		if _, err := p.expect(kwLAST); err != nil {
			return nil, err
		}
		stmt.RetryLast = true

	case p.cur.Type == kwUSING:
		p.advance() // consume USING
		if !p.curIsWord("CONFIG") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume CONFIG
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		cfgTok := p.advance() // consume the config string literal
		stmt.UsingConfig = p.srcSlice(cfgTok.Loc.Start, cfgTok.Loc.End)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

// parseExplainStmt parses `EXPLAIN [ USING { TABULAR | JSON | TEXT } ]
// <statement>`. The EXPLAIN keyword has NOT yet been consumed when this function
// is called. The inner statement is parsed structurally via parseStmt.
func (p *Parser) parseExplainStmt() (ast.Node, error) {
	explainTok := p.advance() // consume EXPLAIN
	stmt := &ast.ExplainStmt{
		Format: ast.ExplainDefault,
		Loc:    ast.Loc{Start: explainTok.Loc.Start},
	}

	// Optional USING { TABULAR | JSON | TEXT }.
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		switch p.cur.Type {
		case kwTABULAR:
			stmt.Format = ast.ExplainTabular
		case kwJSON:
			stmt.Format = ast.ExplainJSON
		case kwTEXT:
			stmt.Format = ast.ExplainText
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume the format keyword
	}

	// The inner statement, parsed via the top-level statement dispatcher.
	inner, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	stmt.Stmt = inner
	// Use the last consumed token's end rather than NodeLoc(inner): the inner
	// statement may be a node type not enumerated in ast.NodeLoc (which is a
	// partial switch), and p.prev always carries a valid end offset.
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
