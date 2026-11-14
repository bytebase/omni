package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB 11.8 refuses of a stored program's declarations while it parses, in the
// scopes sp_scope.go keeps. A block names a variable, a condition and a cursor once each, a
// routine or a cursor a parameter once, and the actions refuse while they read:
//
//	1327 ER_SP_UNDECLARED_VAR       TYPE OF a variable no scope in reach declares
//	1330 ER_SP_DUP_PARAM            a routine's or a cursor's parameter named twice
//	1331 ER_SP_DUP_VAR              a variable a block declares twice
//	1332 ER_SP_DUP_COND             a condition a block declares twice
//	1333 ER_SP_DUP_CURS             a cursor a block declares twice
//	1319 ER_SP_COND_MISMATCH        a handler for a condition no scope in reach declares
//	1413 ER_SP_DUP_HANDLER          two handlers of a block for the same condition
//
// A variable and a condition are found across a handler's scope, which hides labels alone.
// As for labels, a name with a byte past ASCII is never refused: MariaDB folds more of
// Unicode than omni does. Measured on mariadb:11.8.8.
//
// Ref: sp_decl_idents, sp_decl_body, sp_hcond, sp_param_name and
// opt_parenthesized_cursor_formal_parameters in MariaDB 11.8 sql/sql_yacc.yy;
// sp_pcontext::find_variable, find_condition and check_duplicate_handler in sql/sp_pcontext.cc.

// declareSPVar adds a variable to the current scope; param says it is a parameter.
func (p *Parser) declareSPVar(name string, pos int, param bool) error {
	if !p.spChecks() {
		return nil
	}
	s := p.routine.sp
	key := strings.ToLower(name)
	if s.vars[key] && !spLenientName(name) {
		msg := "Duplicate variable: "
		if param {
			msg = "Duplicate parameter: "
		}
		return &ParseError{Message: msg + name, Position: pos}
	}
	if s.vars == nil {
		s.vars = map[string]bool{}
	}
	s.vars[key] = true
	return nil
}

// checkSPVarDeclared refuses a variable no scope in reach declares.
func (p *Parser) checkSPVarDeclared(name string, pos int) error {
	if !p.spChecks() || spLenientName(name) {
		return nil
	}
	key := strings.ToLower(name)
	for s := p.routine.sp; s != nil; s = s.parent {
		if s.vars[key] {
			return nil
		}
	}
	return &ParseError{Message: "Undeclared variable: " + name, Position: pos}
}

// declareSPCondition adds a condition to the current scope.
func (p *Parser) declareSPCondition(name string, val nodes.HandlerCondValue, pos int) error {
	if !p.spChecks() {
		return nil
	}
	s := p.routine.sp
	key := strings.ToLower(name)
	if _, ok := s.conds[key]; ok && !spLenientName(name) {
		return &ParseError{Message: "Duplicate condition: " + name, Position: pos}
	}
	if s.conds == nil {
		s.conds = map[string]nodes.HandlerCondValue{}
	}
	s.conds[key] = val
	return nil
}

// declareSPCursor adds a cursor to the current scope.
func (p *Parser) declareSPCursor(name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	s := p.routine.sp
	key := strings.ToLower(name)
	if s.cursors[key] && !spLenientName(name) {
		return &ParseError{Message: "Duplicate cursor: " + name, Position: pos}
	}
	if s.cursors == nil {
		s.cursors = map[string]bool{}
	}
	s.cursors[key] = true
	return nil
}

// declareSPHandlerCondition adds what a handler of the current scope catches: a named
// condition resolves to the value it was declared for, and two handlers of one block may not
// catch the same value, however each names it.
func (p *Parser) declareSPHandlerCondition(c nodes.HandlerCondValue, pos int) error {
	if !p.spChecks() {
		return nil
	}
	if c.Kind == nodes.HandlerCondName {
		if spLenientName(c.Value) {
			return nil
		}
		key := strings.ToLower(c.Value)
		found := false
		for s := p.routine.sp; s != nil && !found; s = s.parent {
			if v, ok := s.conds[key]; ok {
				c, found = v, true
			}
		}
		if !found {
			return &ParseError{Message: "Undefined CONDITION: " + c.Value, Position: pos}
		}
	}
	s := p.routine.sp
	for _, h := range s.handlers {
		if sameHandlerCondition(h, c) {
			return &ParseError{Message: "Duplicate handler declared in the same block", Position: pos}
		}
	}
	s.handlers = append(s.handlers, c)
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
