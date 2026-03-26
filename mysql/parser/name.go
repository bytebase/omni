package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseIdentifier parses a plain identifier (unquoted or backtick-quoted).
// Returns the identifier string and its start position.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifiers.html
//
//	identifier:
//	    unquoted_identifier
//	    | backtick_quoted_identifier
func (p *Parser) parseIdentifier() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	// Many MySQL keywords can also be used as identifiers in certain contexts.
	// Accept any keyword token as an identifier (non-reserved keyword handling).
	if p.cur.Type >= 700 {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseColumnRef parses a column reference, which may be qualified:
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	column_ref:
//	    identifier
//	    | identifier '.' identifier
//	    | identifier '.' identifier '.' identifier
//	    | identifier '.' '*'
//	    | identifier '.' identifier '.' '*'
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error) {
	start := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	ref := &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start},
		Column: name,
	}

	// Check for dot-qualification
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Check for table.* or schema.table.*
		if p.cur.Type == '*' {
			p.advance()
			ref.Table = name
			ref.Column = ""
			ref.Star = true
			ref.Loc.End = p.pos()
			return ref, nil
		}

		name2, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		// Check for second dot: schema.table.col or schema.table.*
		if p.cur.Type == '.' {
			p.advance() // consume second '.'

			if p.cur.Type == '*' {
				p.advance()
				ref.Schema = name
				ref.Table = name2
				ref.Column = ""
				ref.Star = true
				ref.Loc.End = p.pos()
				return ref, nil
			}

			name3, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			ref.Schema = name
			ref.Table = name2
			ref.Column = name3
			ref.Loc.End = p.pos()
			return ref, nil
		}

		// table.col
		ref.Table = name
		ref.Column = name2
	}

	ref.Loc.End = p.pos()
	return ref, nil
}

// parseTableRef parses a table reference (possibly qualified with schema).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	table_ref:
//	    identifier
//	    | identifier '.' identifier
func (p *Parser) parseTableRef() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for schema.table
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Completion: after "db.", offer table_ref qualified with database.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		name2, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		ref.Schema = name
		ref.Name = name2
	}

	ref.Loc.End = p.pos()
	return ref, nil
}

// parseTableRefWithAlias parses a table reference with an optional alias.
//
//	table_ref_alias:
//	    table_ref [AS identifier | identifier]
func (p *Parser) parseTableRefWithAlias() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// Optional PARTITION (p0, p1, ...)
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		ref.Partitions = parts
		ref.Loc.End = p.pos()
	}

	// Optional AS alias
	if _, ok := p.match(kwAS); ok {
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		ref.Alias = alias
		ref.Loc.End = p.pos()
	} else if p.cur.Type == tokIDENT {
		// Alias without AS keyword
		alias, _, _ := p.parseIdentifier()
		ref.Alias = alias
		ref.Loc.End = p.pos()
	}

	// Optional index hints: USE/FORCE/IGNORE {INDEX|KEY} ...
	if p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		hints, err := p.parseIndexHints()
		if err != nil {
			return nil, err
		}
		ref.IndexHints = hints
		ref.Loc.End = p.pos()
	}

	return ref, nil
}

// parseVariableRef parses a user variable (@var) or system variable (@@var).
// The lexer emits these as tokIDENT with "@" or "@@" prefix in the Str.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/user-variables.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html
//
//	variable_ref:
//	    '@' identifier
//	    | '@@' [GLOBAL | SESSION | LOCAL] '.' identifier
func (p *Parser) parseVariableRef() (*nodes.VariableRef, error) {
	if p.cur.Type != tokIDENT {
		return nil, &ParseError{
			Message:  "expected variable reference",
			Position: p.cur.Loc,
		}
	}

	tok := p.cur
	str := tok.Str

	// Check for @@ prefix (system variable)
	if len(str) > 2 && str[0] == '@' && str[1] == '@' {
		p.advance()
		name := str[2:]
		ref := &nodes.VariableRef{
			Loc:    nodes.Loc{Start: tok.Loc},
			System: true,
		}

		// Check for scope prefix: @@global.var, @@session.var, @@local.var
		if dotIdx := indexOf(name, '.'); dotIdx >= 0 {
			scope := name[:dotIdx]
			varName := name[dotIdx+1:]
			switch {
			case eqFold(scope, "global"):
				ref.Scope = "GLOBAL"
			case eqFold(scope, "session"):
				ref.Scope = "SESSION"
			case eqFold(scope, "local"):
				ref.Scope = "LOCAL"
			default:
				// Not a scope, treat as qualified name
				ref.Name = name
				ref.Loc.End = p.pos()
				return ref, nil
			}
			ref.Name = varName
		} else {
			ref.Name = name
		}

		ref.Loc.End = p.pos()
		return ref, nil
	}

	// Check for @ prefix (user variable)
	if len(str) > 1 && str[0] == '@' {
		p.advance()
		ref := &nodes.VariableRef{
			Loc:  nodes.Loc{Start: tok.Loc},
			Name: str[1:],
		}
		ref.Loc.End = p.pos()
		return ref, nil
	}

	return nil, &ParseError{
		Message:  "expected variable reference",
		Position: p.cur.Loc,
	}
}

// isIdentToken returns true if the current token can be used as an identifier.
func (p *Parser) isIdentToken() bool {
	return p.cur.Type == tokIDENT || p.cur.Type >= 700
}

// isVariableRef returns true if the current token is a variable reference.
func (p *Parser) isVariableRef() bool {
	if p.cur.Type != tokIDENT {
		return false
	}
	return len(p.cur.Str) > 0 && p.cur.Str[0] == '@'
}

// indexOf returns the index of the first occurrence of ch in s, or -1.
func indexOf(s string, ch byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return i
		}
	}
	return -1
}

// eqFold reports whether s and t are equal under Unicode case-folding (ASCII only).
func eqFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}
