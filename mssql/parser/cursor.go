// Package parser - cursor.go implements T-SQL cursor statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// cursorOptions is the set of valid T-SQL cursor option keywords (13 options).
// Matches SqlScriptDOM CursorOptionsHelper.
var cursorOptions = newOptionSet(
	kwLOCAL, kwGLOBAL,
	kwFORWARD_ONLY, kwSCROLL,
	kwSTATIC, kwKEYSET, kwDYNAMIC, kwFAST_FORWARD,
	kwREAD_ONLY, kwREADONLY, kwSCROLL_LOCKS, kwOPTIMISTIC,
	kwTYPE_WARNING,
)

// parseDeclareCursorStmt parses a DECLARE cursor_name CURSOR statement.
//
// BNF: mssql/parser/bnf/declare-cursor-transact-sql.bnf
//
// ISO syntax:
//
//	DECLARE cursor_name [ INSENSITIVE ] [ SCROLL ] CURSOR
//	    FOR select_statement
//	    [ FOR { READ_ONLY | UPDATE [ OF column_name [ , ...n ] ] } ]
//
// Transact-SQL extended syntax:
//
//	DECLARE cursor_name CURSOR [ LOCAL | GLOBAL ]
//	    [ FORWARD_ONLY | SCROLL ]
//	    [ STATIC | KEYSET | DYNAMIC | FAST_FORWARD ]
//	    [ READ_ONLY | SCROLL_LOCKS | OPTIMISTIC ]
//	    [ TYPE_WARNING ]
//	    FOR select_statement
//	    [ FOR UPDATE [ OF column_name [ , ...n ] ] ]
func (p *Parser) parseDeclareCursorStmt() (*nodes.DeclareCursorStmt, error) {
	loc := p.pos()
	p.advance() // consume DECLARE

	// Completion: after DECLARE → @ variable (also covers cursor name context)
	if p.collectMode() {
		p.addRuleCandidate("@variable")
		return nil, errCollecting
	}

	stmt := &nodes.DeclareCursorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Parse cursor name (plain identifier, not @variable). The name, CURSOR
	// and FOR select_statement are required (declareCursorStatement,
	// TSql160.g:16770; cursorDefinitionOptions, TSql160.g:16790).
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// Determine if ISO syntax or T-SQL extended syntax.
	// ISO: cursor_name [ INSENSITIVE ] [ SCROLL ] CURSOR
	// T-SQL: cursor_name CURSOR [ options... ]
	//
	// If we see INSENSITIVE or SCROLL before CURSOR, it's ISO syntax.
	// If we see CURSOR immediately, it's T-SQL extended syntax.

	// Check for ISO keywords before CURSOR
	if p.matchIdentCI("INSENSITIVE") {
		stmt.Insensitive = true
	}

	// SCROLL can appear before CURSOR (ISO) or after CURSOR (T-SQL extended)
	if p.matchIdentCI("SCROLL") {
		stmt.Scroll = true
	}

	// Expect CURSOR keyword
	if _, err := p.expect(kwCURSOR); err != nil {
		return nil, err
	}

	if err := p.parseCursorDefinitionTail(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCursorDefinitionTail parses what follows CURSOR in DECLARE cursor
// and SET @var = CURSOR (cursorDefinitionOptions, TSql160.g:16790): the
// options, FOR select_statement and FOR READ ONLY or FOR UPDATE.
func (p *Parser) parseCursorDefinitionTail(stmt *nodes.DeclareCursorStmt) error {
	// If we already saw INSENSITIVE or SCROLL (ISO syntax), skip to FOR.
	// Otherwise, parse T-SQL extended options between CURSOR and FOR.
	if !stmt.Insensitive && !stmt.Scroll {
		// T-SQL extended syntax: parse options after CURSOR
		if err := p.parseCursorOptions(stmt); err != nil {
			return err
		}
	}

	// FOR select_statement
	if _, err := p.expect(kwFOR); err != nil {
		if p.collectMode() {
			p.addTokenCandidate(kwFOR)
			return errCollecting
		}
		return err
	}
	// Completion: after DECLARE c CURSOR FOR → SELECT
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return errCollecting
	}
	query, err := p.parseSelectStmtWithCTE()
	if err != nil {
		return err
	}
	if isNilNode(query) {
		return p.syntaxErrorAtCur()
	}
	stmt.Query = query

	// Optional: FOR { READ ONLY | UPDATE [ OF column_name [,...n] ] }; SQL
	// Server spells the first READ ONLY here and refuses READ_ONLY.
	if _, ok := p.match(kwFOR); ok {
		switch {
		case p.cur.Type == kwREAD:
			p.advance() // READ
			if err := p.procExpectWord("ONLY"); err != nil {
				return err
			}
			stmt.Concurrency = "READ_ONLY"
		case p.cur.Type == kwUPDATE:
			p.advance() // UPDATE
			stmt.ForUpdate = true
			if _, ok := p.match(kwOF); ok {
				var cols []nodes.Node
				for {
					col, ok := p.parseIdentifier()
					if !ok {
						return p.syntaxErrorAtCur()
					}
					cols = append(cols, &nodes.String{Str: col})
					if _, ok := p.match(','); !ok {
						break
					}
				}
				stmt.UpdateCols = &nodes.List{Items: cols}
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}

	return nil
}

// parseCursorOptions parses T-SQL extended cursor options after CURSOR keyword.
//
//	[ LOCAL | GLOBAL ]
//	[ FORWARD_ONLY | SCROLL ]
//	[ STATIC | KEYSET | DYNAMIC | FAST_FORWARD ]
//	[ READ_ONLY | SCROLL_LOCKS | OPTIMISTIC ]
//	[ TYPE_WARNING ]
//
// Returns an error if an invalid keyword appears in the options position.
func (p *Parser) parseCursorOptions(stmt *nodes.DeclareCursorStmt) error {
	// LOCAL | GLOBAL
	if p.matchIdentCI("LOCAL") {
		stmt.Scope = "LOCAL"
	} else if p.matchIdentCI("GLOBAL") {
		stmt.Scope = "GLOBAL"
	}

	// FORWARD_ONLY | SCROLL
	if p.matchIdentCI("FORWARD_ONLY") {
		stmt.ForwardOnly = true
	} else if p.matchIdentCI("SCROLL") {
		stmt.Scroll = true
	}

	// STATIC | KEYSET | DYNAMIC | FAST_FORWARD
	if p.matchIdentCI("STATIC") {
		stmt.CursorType = "STATIC"
	} else if p.matchIdentCI("KEYSET") {
		stmt.CursorType = "KEYSET"
	} else if p.matchIdentCI("DYNAMIC") {
		stmt.CursorType = "DYNAMIC"
	} else if p.matchIdentCI("FAST_FORWARD") {
		stmt.CursorType = "FAST_FORWARD"
	}

	// READ_ONLY | SCROLL_LOCKS | OPTIMISTIC
	if p.matchIdentCI("READ_ONLY") {
		stmt.Concurrency = "READ_ONLY"
	} else if p.cur.Type == kwREADONLY {
		p.advance()
		stmt.Concurrency = "READ_ONLY"
	} else if p.matchIdentCI("SCROLL_LOCKS") {
		stmt.Concurrency = "SCROLL_LOCKS"
	} else if p.matchIdentCI("OPTIMISTIC") {
		stmt.Concurrency = "OPTIMISTIC"
	}

	// TYPE_WARNING
	if p.matchIdentCI("TYPE_WARNING") {
		stmt.TypeWarning = true
	}

	// Reject invalid keywords in the cursor options position.
	// After consuming all valid options, the next token should be FOR or EOF/semicolon.
	// If it's a keyword or identifier that is NOT a valid cursor option and NOT FOR,
	// it's an error (e.g., DECLARE c CURSOR SELECT FOR SELECT 1).
	if p.cur.Type == kwFOR || p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil
	}
	if p.isIdentLike() && !p.isValidOption(cursorOptions) {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseOpenCursorStmt parses an OPEN cursor statement.
//
// BNF: mssql/parser/bnf/open-transact-sql.bnf
//
//	OPEN { { [ GLOBAL ] cursor_name } | cursor_variable_name }
//
// openCursorStatement (TSql160.g:14487): the cursor is required.
func (p *Parser) parseOpenCursorStmt() (*nodes.OpenCursorStmt, error) {
	loc := p.pos()
	p.advance() // consume OPEN

	// Completion: after OPEN → cursor name
	if p.collectMode() {
		p.addRuleCandidate("cursor_name")
		return nil, errCollecting
	}

	stmt := &nodes.OpenCursorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.Name, stmt.Global, err = p.parseProcCursorRef(); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseFetchCursorStmt parses a FETCH cursor statement.
//
// BNF: mssql/parser/bnf/fetch-transact-sql.bnf
//
//	FETCH
//	    [ [ NEXT | PRIOR | FIRST | LAST
//	            | ABSOLUTE { n | @nvar }
//	            | RELATIVE { n | @nvar }
//	       ]
//	       FROM
//	    ]
//	{ { [ GLOBAL ] cursor_name } | @cursor_variable_name }
//	[ INTO @variable_name [ ,...n ] ]
//
// fetchCursorStatement and rowSelector (TSql160.g:14577, 14597): an
// orientation needs FROM, ABSOLUTE and RELATIVE a signed integer or a
// variable (Msg 153 without), and INTO at least one variable.
func (p *Parser) parseFetchCursorStmt() (*nodes.FetchCursorStmt, error) {
	loc := p.pos()
	p.advance() // consume FETCH

	stmt := &nodes.FetchCursorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Orientation keywords: NEXT, PRIOR, FIRST, LAST, ABSOLUTE, RELATIVE
	for _, w := range []string{"NEXT", "PRIOR", "FIRST", "LAST", "ABSOLUTE", "RELATIVE"} {
		if p.procIsWord(w) {
			stmt.Orientation = w
			p.advance()
			break
		}
	}
	if stmt.Orientation == "ABSOLUTE" || stmt.Orientation == "RELATIVE" {
		var err error
		stmt.FetchOffset, err = p.parseProcValue(procValueKinds{signedInt: true, variable: true})
		if err != nil {
			return nil, err
		}
	}
	if stmt.Orientation != "" {
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
	} else {
		p.match(kwFROM) // optional FROM
	}

	// Completion: after FETCH [NEXT] FROM → cursor name
	if p.collectMode() {
		p.addRuleCandidate("cursor_name")
		return nil, errCollecting
	}

	var err error
	if stmt.Name, stmt.Global, err = p.parseProcCursorRef(); err != nil {
		return nil, err
	}

	// [ INTO @variable_name [,...n] ]
	if _, ok := p.match(kwINTO); ok {
		var vars []nodes.Node
		for {
			if p.cur.Type != tokVARIABLE {
				return nil, p.syntaxErrorAtCur()
			}
			vars = append(vars, &nodes.String{Str: p.cur.Str})
			p.advance()
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.IntoVars = &nodes.List{Items: vars}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCloseCursorStmt parses a CLOSE cursor statement.
//
// BNF: mssql/parser/bnf/close-transact-sql.bnf
//
//	CLOSE { { [ GLOBAL ] cursor_name } | cursor_variable_name }
//
// closeCursorStatement (TSql160.g:14555): the cursor is required.
func (p *Parser) parseCloseCursorStmt() (*nodes.CloseCursorStmt, error) {
	loc := p.pos()
	p.advance() // consume CLOSE

	// Completion: after CLOSE → cursor name
	if p.collectMode() {
		p.addRuleCandidate("cursor_name")
		return nil, errCollecting
	}

	stmt := &nodes.CloseCursorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.Name, stmt.Global, err = p.parseProcCursorRef(); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDeallocateCursorStmt parses a DEALLOCATE cursor statement.
//
// BNF: mssql/parser/bnf/deallocate-transact-sql.bnf
//
//	DEALLOCATE { { [ GLOBAL ] cursor_name } | @cursor_variable_name }
//
// deallocateCursorStatement (TSql160.g:14566): the cursor is required.
func (p *Parser) parseDeallocateCursorStmt() (*nodes.DeallocateCursorStmt, error) {
	loc := p.pos()
	p.advance() // consume DEALLOCATE

	// Completion: after DEALLOCATE → cursor name
	if p.collectMode() {
		p.addRuleCandidate("cursor_name")
		return nil, errCollecting
	}

	stmt := &nodes.DeallocateCursorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.Name, stmt.Global, err = p.parseProcCursorRef(); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// matchIdentCI checks if the current token is an identifier-like token matching
// the given string (case-insensitive). If it matches, consume and return true.
func (p *Parser) matchIdentCI(s string) bool {
	expected := lookupKeyword(s)
	if expected != tokIDENT {
		if p.cur.Type == expected {
			p.advance()
			return true
		}
		return false
	}
	return false
}
