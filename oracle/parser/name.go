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
func (p *Parser) parseIdentifier() string {
	switch p.cur.Type {
	case tokIDENT:
		tok := p.advance()
		return tok.Str
	case tokQIDENT:
		tok := p.advance()
		return tok.Str
	default:
		// Many Oracle keywords can be used as identifiers in non-reserved contexts.
		// If the current token is a keyword, consume it and return as identifier.
		if p.cur.Type >= 2000 {
			tok := p.advance()
			return tok.Str
		}
		return ""
	}
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
func (p *Parser) parseObjectName() *nodes.ObjectName {
	start := p.pos()
	obj := &nodes.ObjectName{
		Loc: nodes.Loc{Start: start},
	}

	name := p.parseIdentifier()
	if name == "" {
		return obj
	}

	// Check for schema.object
	if p.cur.Type == '.' {
		p.advance() // consume '.'
		name2 := p.parseIdentifier()
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
		p.advance() // consume '@'
		obj.DBLink = p.parseIdentifier()
	}

	obj.Loc.End = p.pos()
	return obj
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
func (p *Parser) parseColumnRef() *nodes.ColumnRef {
	start := p.pos()
	col := &nodes.ColumnRef{
		Loc: nodes.Loc{Start: start},
	}

	name1 := p.parseIdentifier()
	if name1 == "" {
		return col
	}

	if p.cur.Type != '.' {
		// Simple column reference
		col.Column = name1
		col.Loc.End = p.pos()
		return col
	}

	// name1.something
	p.advance() // consume '.'

	// Check for name1.*
	if p.cur.Type == '*' {
		p.advance()
		col.Table = name1
		col.Column = "*"
		col.Loc.End = p.pos()
		return col
	}

	name2 := p.parseIdentifier()
	if name2 == "" {
		col.Table = name1
		col.Column = ""
		col.Loc.End = p.pos()
		return col
	}

	if p.cur.Type != '.' {
		// table.column
		col.Table = name1
		col.Column = name2
		col.Loc.End = p.pos()
		return col
	}

	// schema.table.column or schema.table.*
	p.advance() // consume '.'

	if p.cur.Type == '*' {
		p.advance()
		col.Schema = name1
		col.Table = name2
		col.Column = "*"
		col.Loc.End = p.pos()
		return col
	}

	name3 := p.parseIdentifier()
	col.Schema = name1
	col.Table = name2
	col.Column = name3
	col.Loc.End = p.pos()
	return col
}

// parseBindVariable parses a bind variable (:name or :1).
// The current token must be tokBIND.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-language-fundamentals.html
//
//	bind_variable ::= : identifier | : integer
func (p *Parser) parseBindVariable() *nodes.BindVariable {
	if p.cur.Type != tokBIND {
		return nil
	}
	start := p.pos()
	tok := p.advance()
	return &nodes.BindVariable{
		Name: tok.Str,
		Loc:  nodes.Loc{Start: start, End: p.pos()},
	}
}

// parsePseudoColumn parses an Oracle pseudo-column (ROWID, ROWNUM, LEVEL, SYSDATE, SYSTIMESTAMP, USER).
// Returns nil if the current token is not a pseudo-column keyword.
func (p *Parser) parsePseudoColumn() *nodes.PseudoColumn {
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
		return nil
	}

	p.advance()
	return &nodes.PseudoColumn{
		Type: ptype,
		Loc:  nodes.Loc{Start: start, End: p.pos()},
	}
}

// isPseudoColumn returns true if the current token is a pseudo-column keyword.
func (p *Parser) isPseudoColumn() bool {
	switch p.cur.Type {
	case kwROWID, kwROWNUM, kwLEVEL, kwSYSDATE, kwSYSTIMESTAMP, kwUSER:
		return true
	}
	return false
}
