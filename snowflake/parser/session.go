package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Session control — USE / SET / UNSET (T6.3)
// ---------------------------------------------------------------------------
//
// These statements switch the session's active context (USE) or manage
// session variables (SET / UNSET). The grammar is small and stable, mirrored
// directly from the legacy ANTLR rules use_command / set / unset and the
// official alter-session / USE docs.

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

// parseUseStmt parses a USE statement:
//
//	USE [DATABASE] <name>
//	USE ROLE <name>
//	USE [SCHEMA] <name>
//	USE WAREHOUSE <name>
//	USE SECONDARY ROLES { ALL | NONE | <role> [, ...] }
//
// The USE keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseUseStmt() (ast.Node, error) {
	useTok := p.advance() // consume USE
	start := useTok.Loc

	stmt := &ast.UseStmt{Loc: ast.Loc{Start: start.Start}}

	switch p.cur.Type {
	case kwROLE:
		p.advance() // consume ROLE
		stmt.Kind = ast.UseRole
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name

	case kwWAREHOUSE:
		p.advance() // consume WAREHOUSE
		stmt.Kind = ast.UseWarehouse
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name

	case kwDATABASE:
		p.advance() // consume DATABASE
		stmt.Kind = ast.UseDatabase
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name

	case kwSCHEMA:
		p.advance() // consume SCHEMA
		stmt.Kind = ast.UseSchema
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name

	case kwSECONDARY:
		p.advance() // consume SECONDARY
		if _, err := p.expect(kwROLES); err != nil {
			return nil, err
		}
		stmt.Kind = ast.UseSecondaryRoles
		roles, err := p.parseSecondaryRoles()
		if err != nil {
			return nil, err
		}
		stmt.SecondaryRoles = roles

	default:
		// Bare USE <name>: defaults to database/schema context in Snowflake; the
		// parser records the bare form without assuming a kind.
		stmt.Kind = ast.UseDefault
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSecondaryRoles parses the tail of USE SECONDARY ROLES:
//
//	ALL
//	NONE
//	<role> [, <role> ...]
//
// ALL/NONE are returned uppercased verbatim. A role-name list is comma-joined.
// On entry the ROLES keyword has already been consumed.
func (p *Parser) parseSecondaryRoles() (string, error) {
	switch p.cur.Type {
	case kwALL:
		p.advance()
		return "ALL", nil
	case kwNONE:
		p.advance()
		return "NONE", nil
	}

	// Role-name list (Snowflake also accepts an explicit list of roles).
	if !p.isObjectTypeWord(p.cur.Type) {
		return "", p.syntaxErrorAtCur()
	}
	var roles []string
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return "", err
		}
		roles = append(roles, name.String())
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return strings.Join(roles, ", "), nil
}

// ---------------------------------------------------------------------------
// SET
// ---------------------------------------------------------------------------

// parseSetStmt parses a SET session-variable assignment:
//
//	SET <var> = <expr>
//	SET ( <var> [, ...] ) = ( <expr> [, ...] )
//
// The SET keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseSetStmt() (ast.Node, error) {
	setTok := p.advance() // consume SET
	start := setTok.Loc

	stmt := &ast.SetStmt{Loc: ast.Loc{Start: start.Start}}

	if p.cur.Type == '(' {
		// SET ( v1, v2, ... ) = ( e1, e2, ... )
		stmt.Paren = true
		names, nameLocs, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		values, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		if len(values) != len(names) {
			return nil, &ParseError{
				Loc: start,
				Msg: "SET variable/value count mismatch: " +
					strconv.Itoa(len(names)) + " variables but " + strconv.Itoa(len(values)) + " values",
			}
		}
		for i := range names {
			stmt.Vars = append(stmt.Vars, &ast.SetVar{
				Name:  names[i],
				Value: values[i],
				Loc:   ast.Loc{Start: nameLocs[i].Start, End: ast.NodeLoc(values[i]).End},
			})
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// SET <var> = <expr>
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Vars = append(stmt.Vars, &ast.SetVar{
		Name:  name,
		Value: value,
		Loc:   ast.Loc{Start: name.Loc.Start, End: ast.NodeLoc(value).End},
	})

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseParenIdentList parses ( id [, id ...] ) and returns the identifiers and
// their locations. On entry cur is '('. Requires at least one identifier.
func (p *Parser) parseParenIdentList() ([]ast.Ident, []ast.Loc, error) {
	if _, err := p.expect('('); err != nil {
		return nil, nil, err
	}
	var names []ast.Ident
	var locs []ast.Loc
	for {
		name, err := p.parseIdent()
		if err != nil {
			return nil, nil, err
		}
		names = append(names, name)
		locs = append(locs, name.Loc)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	if _, err := p.expect(')'); err != nil {
		return nil, nil, err
	}
	return names, locs, nil
}

// (The parenthesized value list reuses parseParenExprList from select.go, which
// has identical semantics: '(' then a comma-separated expression list then ')'.)

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

// parseUnsetStmt parses UNSET of one or more session variables:
//
//	UNSET <var>
//	UNSET ( <var> [, ...] )
//
// The UNSET keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseUnsetStmt() (ast.Node, error) {
	unsetTok := p.advance() // consume UNSET
	start := unsetTok.Loc

	stmt := &ast.UnsetStmt{Loc: ast.Loc{Start: start.Start}}

	if p.cur.Type == '(' {
		stmt.Paren = true
		names, _, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Names = names
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Names = []ast.Ident{name}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
