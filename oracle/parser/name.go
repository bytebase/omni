package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseIdentifier parses an identifier (unquoted or "double-quoted").
// Returns the identifier string. Unquoted identifiers are already uppercased by the lexer.
// Double-quoted identifiers preserve their original case.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Database-Object-Names-and-Qualifiers.html
//
//	identifier ::= unquoted_identifier | quoted_identifier
func (p *Parser) parseIdentifier() (string, error) {
	switch p.cur.Type {
	case tokIDENT:
		tok := p.advance()
		return tok.Str, nil
	case tokQIDENT:
		tok := p.advance()
		return tok.Str, nil
	default:
		// Many Oracle keywords can be used as identifiers in non-reserved contexts.
		// If the current token is a keyword, consume it and return as identifier.
		if p.cur.Type >= 2000 {
			tok := p.advance()
			return tok.Str, nil
		}
		return "", nil
	}
}

func (p *Parser) parseObjectNameIdentifier() (string, error) {
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return "", err
	}
	return p.parseIdentifier()
}

// parseAlias parses an alias name and records its exact source span.
func (p *Parser) parseAlias() (*nodes.Alias, error) {
	start := p.pos()
	name, parseErr820 := p.parseIdentifier()
	if parseErr820 != nil {
		return nil, parseErr820
	}
	if name == "" {
		return nil, nil
	}
	return &nodes.Alias{
		Name: name,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// isIdentLike returns true if the current token can be treated as an identifier.
// This includes actual identifiers, quoted identifiers, and non-reserved keywords.
func (p *Parser) isIdentLike() bool {
	return p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.cur.Type >= 2000
}

// parseObjectName parses a possibly schema-qualified object name with optional @dblink.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Database-Object-Names-and-Qualifiers.html
//
//	object_name ::= [ schema . ] name [ @dblink ]
func (p *Parser) parseObjectName() (*nodes.ObjectName, error) {
	return p.parseObjectNameWith(p.parseIdentifier)
}

func (p *Parser) parseReservedCheckedObjectName() (*nodes.ObjectName, error) {
	return p.parseObjectNameWith(p.parseObjectNameIdentifier)
}

func (p *Parser) parseObjectNameWith(parseComponent func() (string, error)) (*nodes.ObjectName, error) {
	start := p.pos()
	obj := &nodes.ObjectName{
		Loc: nodes.Loc{Start: start},
	}

	name, parseErr821 := parseComponent()
	if parseErr821 != nil {
		return nil, parseErr821

		// Check for schema.object
	}
	if name == "" {
		return obj, nil
	}

	if p.cur.Type == '.' {
		p.advance() // consume '.'
		name2, parseErr822 := parseComponent()
		if parseErr822 != nil {
			return nil, parseErr822
		}
		if name2 != "" {
			obj.Schema = name
			obj.Name = name2
		} else {
			// Just "name." with no continuation — treat name as the object name
			obj.Name = name
		}
	} else {
		obj.Name = name
	}

	// Check for @dblink
	if p.cur.Type == '@' {
		p.advance()
		var // consume '@'
		parseErr823 error
		obj.DBLink, parseErr823 = parseComponent()
		if parseErr823 != nil {
			return nil, parseErr823
		}
	}

	obj.Loc.End = p.prev.End
	return obj, nil
}

// parseColumnRef parses a column reference which can be:
//
//	column
//	table.column
//	schema.table.column
//	table.*
//	schema.table.*
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Column-Expressions.html
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error) {
	start := p.pos()
	col := &nodes.ColumnRef{
		Loc: nodes.Loc{Start: start},
	}

	name1, parseErr824 := p.parseIdentifier()
	if parseErr824 != nil {
		return nil, parseErr824
	}
	if name1 == "" {
		return col, nil
	}

	if p.cur.Type != '.' {
		// Simple column reference
		col.Column = name1
		col.Loc.End = p.prev.End
		return col, nil
	}

	// name1.something
	p.advance() // consume '.'

	// Check for name1.*
	if p.cur.Type == '*' {
		p.advance()
		col.Table = name1
		col.Column = "*"
		col.Loc.End = p.prev.End
		return col, nil
	}

	name2, parseErr825 := p.parseIdentifier()
	if parseErr825 != nil {
		return nil, parseErr825
	}
	if name2 == "" {
		col.Table = name1
		col.Column = ""
		col.Loc.End = p.prev.End
		return col, nil
	}

	if p.cur.Type != '.' {
		// table.column
		col.Table = name1
		col.Column = name2
		col.Loc.End = p.prev.End
		return col, nil
	}

	// schema.table.column or schema.table.*
	p.advance() // consume '.'

	if p.cur.Type == '*' {
		p.advance()
		col.Schema = name1
		col.Table = name2
		col.Column = "*"
		col.Loc.End = p.prev.End
		return col, nil
	}

	name3, parseErr826 := p.parseIdentifier()
	if parseErr826 != nil {
		return nil, parseErr826
	}
	col.Schema = name1
	col.Table = name2
	col.Column = name3
	col.Loc.End = p.prev.End
	return col, nil
}

// parseBindVariable parses a bind variable (:name or :1).
// The current token must be tokBIND.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-language-fundamentals.html
//
//	bind_variable ::= : identifier | : integer
func (p *Parser) parseBindVariable() (*nodes.BindVariable, error) {
	if p.cur.Type != tokBIND {
		return nil, nil
	}
	start := p.pos()
	tok := p.advance()
	bv := &nodes.BindVariable{
		Name: tok.Str,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}
	// Handle :name.member (e.g., :NEW.created_date in trigger bodies)
	if p.cur.Type == '.' && p.peekNext().Type != '*' {
		p.advance() // consume '.'
		if p.isIdentLike() || p.cur.Type == tokQIDENT {
			var parseErr827 error
			bv.Member, parseErr827 = p.parseIdentifier()
			if parseErr827 != nil {
				return nil, parseErr827
			}
			bv.Loc.End = p.prev.End
		}
	}
	return bv, nil
}

// parsePseudoColumn parses an Oracle pseudo-column (ROWID, ROWNUM, LEVEL, SYSDATE, SYSTIMESTAMP, USER).
// Returns nil if the current token is not a pseudo-column keyword.
func (p *Parser) parsePseudoColumn() (*nodes.PseudoColumn, error) {
	start := p.pos()
	var ptype nodes.PseudoColumnType

	switch p.cur.Type {
	case kwROWID:
		ptype = nodes.PSEUDO_ROWID
	case kwROWNUM:
		ptype = nodes.PSEUDO_ROWNUM
	case kwLEVEL:
		ptype = nodes.PSEUDO_LEVEL
	case kwSYSDATE:
		ptype = nodes.PSEUDO_SYSDATE
	case kwSYSTIMESTAMP:
		ptype = nodes.PSEUDO_SYSTIMESTAMP
	case kwUSER:
		ptype = nodes.PSEUDO_USER
	default:
		return nil, nil
	}

	p.advance()
	return &nodes.PseudoColumn{
		Type: ptype,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// isPseudoColumn returns true if the current token is a pseudo-column keyword.
func (p *Parser) isPseudoColumn() bool {
	switch p.cur.Type {
	case kwROWID, kwROWNUM, kwLEVEL, kwSYSDATE, kwSYSTIMESTAMP, kwUSER:
		return true
	}
	return false
}
