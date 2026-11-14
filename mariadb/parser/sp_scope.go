package parser

import (
	"strings"
	"unicode/utf8"
)

// What MariaDB 11.8 refuses of the names inside a stored program while it parses. The grammar
// actions keep sp_pcontext, a tree of scopes: the program's own, one for each BEGIN ... END
// block, and one for the statement of each DECLARE ... HANDLER. A name is looked up from the
// innermost scope out, and a label is not looked up past a handler's scope: a handler cannot
// leave or iterate a block or loop around it.
//
// Ref: sql/sp_pcontext.cc and sql/sp_pcontext.h; LEX::sp_block_init, sp_push_loop_label,
// sp_leave_statement, sp_iterate_statement and sp_handler_declaration_init in sql/sql_lex.cc;
// sp_block_label in sql/sql_yacc.yy, MariaDB 11.8.

// spScope is one sp_pcontext. Variables, cursors and labels are each a namespace of their
// own: a variable and a cursor may share a name.
type spScope struct {
	parent  *spScope
	handler bool // a handler's statement, sp_pcontext::HANDLER_SCOPE
	labels  []spLabel
	vars    []string // the parameters in the program's scope, DECLAREd variables in a block's
	cursors []string
}

// spLabel is a block's or a loop's label. A block or a loop without one pushes the empty
// name, as MariaDB does.
type spLabel struct {
	name string
	loop bool // LOOP, WHILE or REPEAT, which ITERATE may name; a BEGIN block's is not
}

// spNameMatch reports whether a name written in a stored program is one declared there.
// MariaDB compares labels and names without case, under a collation that knows more of
// Unicode than omni folds. A lookup of a name with a byte past ASCII is taken as found: a
// refusal must be the server's.
func spNameMatch(written, declared string) bool {
	return strings.EqualFold(written, declared)
}

// spNameMayMatch is spNameMatch for a lookup, which must not refuse a name the server finds:
// a name with a byte past ASCII on either side may be one the server's collation folds to
// the other.
func spNameMayMatch(written, declared string) bool {
	return spNameMatch(written, declared) || spLenientName(written) || spLenientName(declared)
}

func spLenientName(name string) bool {
	for i := 0; i < len(name); i++ {
		if name[i] >= utf8.RuneSelf {
			return true
		}
	}
	return false
}

// findLabel is sp_pcontext::find_label: the scope's own labels, innermost first, then the
// parent's unless the scope is a handler's.
func (s *spScope) findLabel(name string) *spLabel {
	for ; s != nil; s = s.parent {
		for i := len(s.labels) - 1; i >= 0; i-- {
			if spNameMatch(name, s.labels[i].name) {
				return &s.labels[i]
			}
		}
		if s.handler {
			return nil
		}
	}
	return nil
}

// findVar is sp_pcontext::find_variable: a parameter or variable in this scope or one
// around it. A handler's scope does not hide them.
func (s *spScope) findVar(name string) bool {
	for ; s != nil; s = s.parent {
		for _, v := range s.vars {
			if spNameMayMatch(name, v) {
				return true
			}
		}
	}
	return false
}

// findCursor is sp_pcontext::find_cursor.
func (s *spScope) findCursor(name string) bool {
	for ; s != nil; s = s.parent {
		for _, c := range s.cursors {
			if spNameMayMatch(name, c) {
				return true
			}
		}
	}
	return false
}

func spHasName(names []string, name string) bool {
	for _, n := range names {
		if spNameMatch(name, n) {
			return true
		}
	}
	return false
}

// spChecks reports whether the parser holds a stored program's names to its scope: inside
// one, and not while collecting completion candidates, where the text is unfinished.
func (p *Parser) spChecks() bool {
	return p.routine.sp != nil && !p.completing
}

// enterStoredProgram starts the scope of a routine, trigger or event body, which holds the
// routine's parameters. The returned function restores the scope around it.
func (p *Parser) enterStoredProgram(r routineScope, params ...string) func() {
	outer := p.routine
	r.sp = &spScope{vars: params}
	p.routine = r
	return func() { p.routine = outer }
}

// pushSPScope opens a BEGIN ... END block's scope, or with handler a handler statement's.
// The returned function closes it.
func (p *Parser) pushSPScope(handler bool) func() {
	outer := p.routine.sp
	if outer == nil {
		return func() {}
	}
	p.routine.sp = &spScope{parent: outer, handler: handler}
	return func() { p.routine.sp = outer }
}

// pushSPLabel declares the label of a block or loop starting here, the empty name when it
// has none, and returns the function that drops it at the block's or loop's end. A label
// already in reach is 1309, ER_SP_LABEL_REDEFINE.
//
// Ref: sp_pcontext::block_label_declare and LEX::sp_push_loop_label, MariaDB 11.8.
func (p *Parser) pushSPLabel(name string, pos int, loop bool) (func(), error) {
	s := p.routine.sp
	if s == nil {
		return func() {}, nil
	}
	if name != "" && p.spChecks() && s.findLabel(name) != nil {
		return nil, &ParseError{Message: "Redefining label " + name, Position: pos}
	}
	n := len(s.labels)
	s.labels = append(s.labels, spLabel{name: name, loop: loop})
	return func() { s.labels = s.labels[:n] }, nil
}

// checkSPLabelRef holds LEAVE's and ITERATE's label to the blocks and loops around it: 1308,
// ER_SP_LILABEL_MISMATCH, when none in reach has it, and for ITERATE when it names a block.
//
// Ref: LEX::sp_leave_statement and LEX::sp_iterate_statement in MariaDB 11.8 sql/sql_lex.cc.
func (p *Parser) checkSPLabelRef(verb, name string, pos int, iterate bool) error {
	if !p.spChecks() {
		return nil
	}
	lab := p.routine.sp.findLabel(name)
	if lab == nil && spLenientName(name) {
		return nil
	}
	if lab == nil || (iterate && !lab.loop) {
		return &ParseError{Message: verb + " with no matching label: " + name, Position: pos}
	}
	return nil
}

// checkSPParam refuses a parameter named as one before it: 1330, ER_SP_DUP_PARAM.
//
// Ref: LEX::sp_param_init in MariaDB 11.8 sql/sql_lex.cc.
func (p *Parser) checkSPParam(before []string, name string, pos int) error {
	if !p.completing && spHasName(before, name) {
		return &ParseError{Message: "Duplicate parameter: " + name, Position: pos}
	}
	return nil
}

// declareSPVar declares a block's variable: one of the same name in the block is 1331,
// ER_SP_DUP_VAR. One of the same name around the block is hidden by it.
//
// Ref: sp_decl_idents in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) declareSPVar(name string, pos int) error {
	s := p.routine.sp
	if s == nil {
		return nil
	}
	if p.spChecks() && spHasName(s.vars, name) {
		return &ParseError{Message: "Duplicate variable: " + name, Position: pos}
	}
	s.vars = append(s.vars, name)
	return nil
}

// declareSPCursor declares a block's cursor: one of the same name in the block is 1333,
// ER_SP_DUP_CURS.
//
// Ref: LEX::sp_declare_cursor in MariaDB 11.8 sql/sql_lex.cc.
func (p *Parser) declareSPCursor(name string, pos int) error {
	s := p.routine.sp
	if s == nil {
		return nil
	}
	if p.spChecks() && spHasName(s.cursors, name) {
		return &ParseError{Message: "Duplicate cursor: " + name, Position: pos}
	}
	s.cursors = append(s.cursors, name)
	return nil
}

// checkSPCursorRef holds OPEN's, FETCH's and CLOSE's cursor to the cursors in reach: 1324,
// ER_SP_CURSOR_MISMATCH.
//
// Ref: sp_proc_stmt_open and sp_proc_stmt_close in MariaDB 11.8 sql/sql_yacc.yy, and
// LEX::sp_add_instr_cfetch in sql/sql_lex.cc.
func (p *Parser) checkSPCursorRef(name string, pos int) error {
	if p.spChecks() && !p.routine.sp.findCursor(name) {
		return &ParseError{Message: "Undefined CURSOR: " + name, Position: pos}
	}
	return nil
}

// checkSPVarRef holds a name that must be a stored program's variable or parameter to the
// ones in reach: 1327, ER_SP_UNDECLARED_VAR. Outside a stored program, outside reports
// whether the server refuses the name too, as it does a FETCH's, GET DIAGNOSTICS's, SELECT
// ... INTO's or LIMIT's while it reads it.
//
// Ref: LEX::make_fetch_target, LEX::create_item_for_sp_var, LEX::create_outvar and
// LEX::create_item_limit in MariaDB 11.8 sql/sql_lex.cc.
func (p *Parser) checkSPVarRef(name string, pos int, outside bool) error {
	if p.completing {
		return nil
	}
	if s := p.routine.sp; s != nil {
		if s.findVar(name) {
			return nil
		}
	} else if !outside {
		return nil
	}
	return &ParseError{Message: "Undeclared variable: " + name, Position: pos}
}
