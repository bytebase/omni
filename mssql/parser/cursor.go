// Package parser - cursor.go implements T-SQL cursor statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
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
func (p *Parser) parseDeclareCursorStmt() *nodes.DeclareCursorStmt {
	loc := p.pos()
	p.advance() // consume DECLARE

	stmt := &nodes.DeclareCursorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Parse cursor name (plain identifier, not @variable)
	name, ok := p.parseIdentifier()
	if !ok {
		return stmt
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
	if _, ok := p.match(kwCURSOR); !ok {
		stmt.Loc.End = p.pos()
		return stmt
	}

	// If we already saw INSENSITIVE or SCROLL (ISO syntax), skip to FOR.
	// Otherwise, parse T-SQL extended options between CURSOR and FOR.
	if !stmt.Insensitive && !stmt.Scroll {
		// T-SQL extended syntax: parse options after CURSOR
		p.parseCursorOptions(stmt)
	}

	// FOR select_statement
	if _, ok := p.match(kwFOR); ok {
		stmt.Query = p.parseSelectStmt()
	}

	// Optional: FOR { READ_ONLY | UPDATE [ OF column_name [,...n] ] }
	if _, ok := p.match(kwFOR); ok {
		if p.matchIdentCI("READ_ONLY") || p.cur.Type == kwREADONLY {
			if p.cur.Type == kwREADONLY {
				p.advance()
			}
			stmt.Concurrency = "READ_ONLY"
		} else if _, ok := p.match(kwUPDATE); ok {
			stmt.ForUpdate = true
			if _, ok := p.match(kwOF); ok {
				var cols []nodes.Node
				for {
					col, ok := p.parseIdentifier()
					if !ok {
						break
					}
					cols = append(cols, &nodes.String{Str: col})
					if _, ok := p.match(','); !ok {
						break
					}
				}
				if len(cols) > 0 {
					stmt.UpdateCols = &nodes.List{Items: cols}
				}
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCursorOptions parses T-SQL extended cursor options after CURSOR keyword.
//
//	[ LOCAL | GLOBAL ]
//	[ FORWARD_ONLY | SCROLL ]
//	[ STATIC | KEYSET | DYNAMIC | FAST_FORWARD ]
//	[ READ_ONLY | SCROLL_LOCKS | OPTIMISTIC ]
//	[ TYPE_WARNING ]
func (p *Parser) parseCursorOptions(stmt *nodes.DeclareCursorStmt) {
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
}

// parseOpenCursorStmt parses an OPEN cursor statement.
//
// BNF: mssql/parser/bnf/open-transact-sql.bnf
//
//	OPEN { { [ GLOBAL ] cursor_name } | cursor_variable_name }
func (p *Parser) parseOpenCursorStmt() *nodes.OpenCursorStmt {
	loc := p.pos()
	p.advance() // consume OPEN

	stmt := &nodes.OpenCursorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for @cursor_variable
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
		stmt.Loc.End = p.pos()
		return stmt
	}

	// [ GLOBAL ] cursor_name
	if p.matchIdentCI("GLOBAL") {
		stmt.Global = true
	}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseFetchCursorStmt() *nodes.FetchCursorStmt {
	loc := p.pos()
	p.advance() // consume FETCH

	stmt := &nodes.FetchCursorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Try to parse orientation [ FROM ]
	// Orientation keywords: NEXT, PRIOR, FIRST, LAST, ABSOLUTE, RELATIVE
	orientation := ""
	switch {
	case p.matchIdentCI("NEXT"):
		orientation = "NEXT"
	case p.matchIdentCI("PRIOR"):
		orientation = "PRIOR"
	case p.matchIdentCI("FIRST"):
		orientation = "FIRST"
	case p.matchIdentCI("LAST"):
		orientation = "LAST"
	case p.matchIdentCI("ABSOLUTE"):
		orientation = "ABSOLUTE"
		stmt.FetchOffset, _ = p.parseExpr()
	case p.matchIdentCI("RELATIVE"):
		orientation = "RELATIVE"
		stmt.FetchOffset, _ = p.parseExpr()
	}

	if orientation != "" {
		stmt.Orientation = orientation
		// Expect FROM after orientation
		p.match(kwFROM)
	}

	// If no orientation was specified, check for FROM keyword (FETCH FROM cursor_name)
	// or just the cursor name directly.
	if orientation == "" {
		// Could be: FETCH FROM cursor_name or FETCH cursor_name or FETCH @var
		p.match(kwFROM) // optional FROM
	}

	// Parse cursor reference: [ GLOBAL ] cursor_name | @cursor_variable
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	} else {
		if p.matchIdentCI("GLOBAL") {
			stmt.Global = true
		}
		name, _ := p.parseIdentifier()
		stmt.Name = name
	}

	// [ INTO @variable_name [,...n] ]
	if _, ok := p.match(kwINTO); ok {
		var vars []nodes.Node
		for {
			if p.cur.Type != tokVARIABLE {
				break
			}
			vars = append(vars, &nodes.String{Str: p.cur.Str})
			p.advance()
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(vars) > 0 {
			stmt.IntoVars = &nodes.List{Items: vars}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCloseCursorStmt parses a CLOSE cursor statement.
//
// BNF: mssql/parser/bnf/close-transact-sql.bnf
//
//	CLOSE { { [ GLOBAL ] cursor_name } | cursor_variable_name }
func (p *Parser) parseCloseCursorStmt() *nodes.CloseCursorStmt {
	loc := p.pos()
	p.advance() // consume CLOSE

	stmt := &nodes.CloseCursorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for @cursor_variable
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
		stmt.Loc.End = p.pos()
		return stmt
	}

	// [ GLOBAL ] cursor_name
	if p.matchIdentCI("GLOBAL") {
		stmt.Global = true
	}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDeallocateCursorStmt parses a DEALLOCATE cursor statement.
//
// BNF: mssql/parser/bnf/deallocate-transact-sql.bnf
//
//	DEALLOCATE { { [ GLOBAL ] cursor_name } | @cursor_variable_name }
func (p *Parser) parseDeallocateCursorStmt() *nodes.DeallocateCursorStmt {
	loc := p.pos()
	p.advance() // consume DEALLOCATE

	stmt := &nodes.DeallocateCursorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Check for @cursor_variable
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
		stmt.Loc.End = p.pos()
		return stmt
	}

	// [ GLOBAL ] cursor_name
	if p.matchIdentCI("GLOBAL") {
		stmt.Global = true
	}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt
}

// matchIdentCI checks if the current token is an identifier-like token matching
// the given string (case-insensitive). If it matches, consume and return true.
func (p *Parser) matchIdentCI(s string) bool {
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, s) {
		p.advance()
		return true
	}
	return false
}
