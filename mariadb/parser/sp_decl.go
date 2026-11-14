package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB 11.8 refuses of a block's declarations while it parses, beyond the variables,
// parameters and cursors sp_scope.go holds:
//
//	1327 ER_SP_UNDECLARED_VAR       TYPE OF a variable no scope in reach declares
//	1332 ER_SP_DUP_COND             a condition a block declares twice
//	1319 ER_SP_COND_MISMATCH        a handler for a condition no scope in reach declares
//	1413 ER_SP_DUP_HANDLER          two handlers of a block for the same condition
//
// A condition is found across a handler's scope, which hides labels alone. As for the other
// names, a lookup of a name with a byte past ASCII is never refused. Measured on
// mariadb:11.8.8.
//
// Ref: sp_decl_variable_list_anchored, sp_decl_body and sp_hcond in MariaDB 11.8
// sql/sql_yacc.yy; sp_pcontext::find_condition and check_duplicate_handler in
// sql/sp_pcontext.cc.

// checkSPVarFree refuses a variable a DECLARE names that the block, or the DECLARE before it,
// already names: 1331. The names are declared once the DECLARE's type is read, so its TYPE OF
// does not see them (declareSPVars).
func (p *Parser) checkSPVarFree(before []string, name string, pos int) error {
	s := p.routine.sp
	if !p.spChecks() {
		return nil
	}
	if spHasName(s.vars, name) || spHasName(before, name) {
		return &ParseError{Message: "Duplicate variable: " + name, Position: pos}
	}
	return nil
}

// declareSPVars adds a DECLARE's variables to the current block.
func (p *Parser) declareSPVars(names []string) {
	if s := p.routine.sp; s != nil {
		s.vars = append(s.vars, names...)
	}
}

// checkSPVarDeclared refuses TYPE OF a variable no scope in reach declares.
func (p *Parser) checkSPVarDeclared(name string, pos int) error {
	if p.spChecks() && !p.routine.sp.findVar(name) {
		return &ParseError{Message: "Undeclared variable: " + name, Position: pos}
	}
	return nil
}

// declareSPCondition adds a condition to the current scope.
func (p *Parser) declareSPCondition(name string, val nodes.HandlerCondValue, pos int) error {
	s := p.routine.sp
	if s == nil {
		return nil
	}
	key := strings.ToLower(name)
	if _, ok := s.conds[key]; ok && p.spChecks() && !spLenientName(name) {
		return &ParseError{Message: "Duplicate condition: " + name, Position: pos}
	}
	if s.conds == nil {
		s.conds = map[string]nodes.HandlerCondValue{}
	}
	s.conds[key] = val
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
