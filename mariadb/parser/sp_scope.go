package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// spScope is one block of a stored program, or the parameter list around its body, with
// what its declarations named. MariaDB's grammar actions check a declaration against it as
// they read it (sp_pcontext), so a block that repeats a name, or a handler for a condition
// nobody declared, is refused while it parses:
//
//	1327 ER_SP_UNDECLARED_VAR       TYPE OF a variable no block around declares
//	1330 ER_SP_DUP_PARAM            a routine's or a cursor's parameter named twice
//	1331 ER_SP_DUP_VAR              a variable a block declares twice
//	1332 ER_SP_DUP_COND             a condition a block declares twice
//	1333 ER_SP_DUP_CURS             a cursor a block declares twice
//	1319 ER_SP_COND_MISMATCH        a handler for a condition no block around declares
//	1413 ER_SP_DUP_HANDLER          two handlers of a block for the same condition
//
// Names compare case-insensitively. Measured on mariadb:11.8.8.
//
// Ref: sp_decl_idents, sp_decl_body, sp_hcond, sp_param_name and
// opt_parenthesized_cursor_formal_parameters in MariaDB 11.8 sql/sql_yacc.yy.
type spScope struct {
	parent   *spScope
	vars     map[string]bool
	conds    map[string]nodes.HandlerCondValue
	cursors  map[string]bool
	handlers []nodes.HandlerCondValue // what this block's handlers catch, names resolved
}

// pushSpScope opens a block; the returned function closes it.
func (p *Parser) pushSpScope() func() {
	p.spScope = &spScope{parent: p.spScope}
	return func() { p.spScope = p.spScope.parent }
}

// enterStoredProgram starts a stored program's outermost scope, so a block of an earlier
// statement is never searched; the returned function restores the one around it.
func (p *Parser) enterStoredProgram() func() {
	outer := p.spScope
	p.spScope = nil
	return func() { p.spScope = outer }
}

// spChecks reports whether the declaration checks run: inside a stored program, and not
// while completing, where the text is unfinished.
func (p *Parser) spChecks() bool {
	return p.spScope != nil && !p.completing
}

// declareSpVar adds a variable to the current block, refusing a name it already holds.
func (p *Parser) declareSpVar(name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	key := strings.ToLower(name)
	if p.spScope.vars[key] {
		return &ParseError{Message: "Duplicate variable: " + name, Position: pos}
	}
	if p.spScope.vars == nil {
		p.spScope.vars = map[string]bool{}
	}
	p.spScope.vars[key] = true
	return nil
}

// declareSpParam adds a routine's or a cursor's parameter, refusing a repeated name.
func (p *Parser) declareSpParam(name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	key := strings.ToLower(name)
	if p.spScope.vars[key] {
		return &ParseError{Message: "Duplicate parameter: " + name, Position: pos}
	}
	if p.spScope.vars == nil {
		p.spScope.vars = map[string]bool{}
	}
	p.spScope.vars[key] = true
	return nil
}

// checkSpVarDeclared refuses a variable no block around the current one declares.
func (p *Parser) checkSpVarDeclared(name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	key := strings.ToLower(name)
	for s := p.spScope; s != nil; s = s.parent {
		if s.vars[key] {
			return nil
		}
	}
	return &ParseError{Message: "Undeclared variable: " + name, Position: pos}
}

// declareSpCondition adds a condition to the current block.
func (p *Parser) declareSpCondition(name string, val nodes.HandlerCondValue, pos int) error {
	if !p.spChecks() {
		return nil
	}
	key := strings.ToLower(name)
	if _, ok := p.spScope.conds[key]; ok {
		return &ParseError{Message: "Duplicate condition: " + name, Position: pos}
	}
	if p.spScope.conds == nil {
		p.spScope.conds = map[string]nodes.HandlerCondValue{}
	}
	p.spScope.conds[key] = val
	return nil
}

// declareSpCursor adds a cursor to the current block.
func (p *Parser) declareSpCursor(name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	key := strings.ToLower(name)
	if p.spScope.cursors[key] {
		return &ParseError{Message: "Duplicate cursor: " + name, Position: pos}
	}
	if p.spScope.cursors == nil {
		p.spScope.cursors = map[string]bool{}
	}
	p.spScope.cursors[key] = true
	return nil
}

// declareSpHandlerCondition adds what a handler of the current block catches: a named
// condition resolves to the value it was declared for, and two handlers of one block may
// not catch the same value, however each names it.
func (p *Parser) declareSpHandlerCondition(c nodes.HandlerCondValue, pos int) error {
	if !p.spChecks() {
		return nil
	}
	if c.Kind == nodes.HandlerCondName {
		key := strings.ToLower(c.Value)
		found := false
		for s := p.spScope; s != nil && !found; s = s.parent {
			if v, ok := s.conds[key]; ok {
				c, found = v, true
			}
		}
		if !found {
			return &ParseError{Message: "Undefined CONDITION: " + c.Value, Position: pos}
		}
	}
	for _, h := range p.spScope.handlers {
		if sameHandlerCondition(h, c) {
			return &ParseError{Message: "Duplicate handler declared in the same block", Position: pos}
		}
	}
	p.spScope.handlers = append(p.spScope.handlers, c)
	return nil
}

// sameHandlerCondition compares two resolved condition values as sp_condition_value::equals
// does: an error code by its number, an SQLSTATE by its text, a class by its kind.
func sameHandlerCondition(a, b nodes.HandlerCondValue) bool {
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case nodes.HandlerCondErrorCode:
		return strings.TrimLeft(a.Value, "0") == strings.TrimLeft(b.Value, "0")
	case nodes.HandlerCondSQLState:
		return a.Value == b.Value
	}
	return true
}
