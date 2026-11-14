package parser

import (
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// What MySQL 8.0 refuses of the names inside a stored program while it parses. The grammar
// actions keep sp_pcontext, a tree of scopes: the program's own, one for each BEGIN ... END
// block, and one for the statement of each DECLARE ... HANDLER. A name is looked up from the
// innermost scope out, and a label is not looked up past a handler's scope: a handler cannot
// leave or iterate a block or loop around it.
//
// Ref: sql/sp_pcontext.cc, and the sp_labeled_block, sp_unlabeled_block, sp_labeled_control,
// sp_unlabeled_control, sp_decl (HANDLER), sp_proc_stmt_leave and sp_proc_stmt_iterate
// rules of sql/sql_yacc.yy, MySQL 8.0.

// spScope is one sp_pcontext. Variables, cursors and labels are each a namespace of their
// own: a variable and a cursor may share a name.
type spScope struct {
	parent  *spScope
	handler bool // a handler's statement, sp_pcontext::HANDLER_SCOPE
	labels  []spLabel
	vars    []string // the parameters in the program's scope, DECLAREd variables in a block's
	cursors []string
	conds   []spCondition
	// handlerConds holds the condition values of the block's handlers, a named condition as
	// the value it stands for, which sp_pcontext::check_duplicate_handler compares.
	handlerConds []nodes.HandlerCondValue
}

// spCondition is a DECLARE ... CONDITION: its name and the error code or SQLSTATE it names.
type spCondition struct {
	name  string
	value nodes.HandlerCondValue
}

// spLabel is a block's or a loop's label. A block or a loop without one pushes the empty
// name, as MySQL does.
type spLabel struct {
	name string
	loop bool // LOOP, WHILE or REPEAT, which ITERATE may name; a BEGIN block's is not
}

// spNameMatch reports whether a name written in a stored program is one declared there.
// MySQL compares labels and names under utf8mb3_general_ci (my_strcasecmp and my_strnncoll
// on system_charset_info), which also folds accents. omni folds only case, so a lookup of a
// name with a byte past ASCII is taken as found: a refusal must be the server's.
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

// findCondition is sp_pcontext::find_condition, which takes the innermost declaration the
// name folds to. known is false when the name was taken as found only because omni does not
// fold it as the server does (spNameMayMatch): a declaration that may be the one the server
// finds leaves which one it is open, though one past it matches exactly.
func (s *spScope) findCondition(name string) (value nodes.HandlerCondValue, found, known bool) {
	for ; s != nil; s = s.parent {
		for i := len(s.conds) - 1; i >= 0; i-- {
			if spNameMatch(name, s.conds[i].name) {
				return s.conds[i].value, true, true
			}
			if spNameMayMatch(name, s.conds[i].name) {
				return nodes.HandlerCondValue{}, true, false
			}
		}
	}
	return nodes.HandlerCondValue{}, spLenientName(name), false
}

// spSameCondition is sp_condition_value::equals: the same kind, and for an error code or a
// SQLSTATE the same one.
func spSameCondition(a, b nodes.HandlerCondValue) bool {
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
// Ref: sp_labeled_block and sp_labeled_control in MySQL 8.0 sql/sql_yacc.yy.
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
// Ref: sp_proc_stmt_leave and sp_proc_stmt_iterate in MySQL 8.0 sql/sql_yacc.yy.
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
// Ref: sp_pdparam and sp_fdparam in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) checkSPParam(before []string, name string, pos int) error {
	if !p.completing && spHasName(before, name) {
		return &ParseError{Message: "Duplicate parameter: " + name, Position: pos}
	}
	return nil
}

// declareSPVar declares a block's variable: one of the same name in the block is 1331,
// ER_SP_DUP_VAR. One of the same name around the block is hidden by it.
//
// Ref: sp_decl_idents in MySQL 8.0 sql/sql_yacc.yy.
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
// Ref: the DECLARE ... CURSOR arm of sp_decl in MySQL 8.0 sql/sql_yacc.yy.
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
// Ref: sp_proc_stmt_open, sp_proc_stmt_fetch and sp_proc_stmt_close in MySQL 8.0
// sql/sql_yacc.yy.
func (p *Parser) checkSPCursorRef(name string, pos int) error {
	if p.spChecks() && !p.routine.sp.findCursor(name) {
		return &ParseError{Message: "Undefined CURSOR: " + name, Position: pos}
	}
	return nil
}

// checkSPVarRef holds a name that must be a stored program's variable or parameter to the
// ones in reach: 1327, ER_SP_UNDECLARED_VAR. Outside a stored program, outside reports
// whether the server refuses the name too, as GET DIAGNOSTICS's target does in its grammar
// action; where it refuses one only once the statement is read (SELECT ... INTO, LIMIT), the
// check is left to it.
//
// Ref: sp_fetch_list and simple_target_specification in MySQL 8.0 sql/sql_yacc.yy;
// PT_select_sp_var::contextualize in sql/parse_tree_nodes.cc and PTI_limit_option_ident::itemize
// in sql/parse_tree_items.cc,
// which a stored program's statement runs while it parses (MAKE_CMD in sp_proc_stmt_statement).
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

// declareSPCondition declares a block's condition: one of the same name in the block is
// 1332, ER_SP_DUP_COND.
//
// Ref: the DECLARE ... CONDITION arm of sp_decl in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) declareSPCondition(name string, pos int, value nodes.HandlerCondValue) error {
	s := p.routine.sp
	if s == nil {
		return nil
	}
	if p.spChecks() {
		for _, c := range s.conds {
			if spNameMatch(name, c.name) {
				return &ParseError{Message: "Duplicate condition: " + name, Position: pos}
			}
		}
	}
	s.conds = append(s.conds, spCondition{name: name, value: value})
	return nil
}

// declareSPHandlerCond adds a condition value of a handler declared in the current block. A
// named condition must be in reach (1319, ER_SP_COND_MISMATCH), and a value one of the
// block's handlers already has, this one's included, is 1413, ER_SP_DUP_HANDLER.
//
// Ref: sp_hcond and sp_hcond_element in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) declareSPHandlerCond(c nodes.HandlerCondValue, pos int) error {
	s := p.routine.sp
	if !p.spChecks() {
		return nil
	}
	if c.Kind == nodes.HandlerCondName {
		value, found, known := s.findCondition(c.Value)
		if !found {
			return &ParseError{Message: "Undefined CONDITION: " + c.Value, Position: pos}
		}
		if !known {
			return nil
		}
		c = value
	}
	for _, h := range s.handlerConds {
		if spSameCondition(c, h) {
			return &ParseError{Message: "Duplicate handler declared in the same block", Position: pos}
		}
	}
	s.handlerConds = append(s.handlerConds, c)
	return nil
}

// checkSPSignalCond holds the condition SIGNAL or RESIGNAL names to the ones in reach, none
// outside a stored program (1319, ER_SP_COND_MISMATCH), and to one declared for a SQLSTATE
// (1646, ER_SIGNAL_BAD_CONDITION_TYPE).
//
// Ref: signal_value in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) checkSPSignalCond(name string, pos int) error {
	if p.completing {
		return nil
	}
	var value nodes.HandlerCondValue
	found, known := false, false
	if s := p.routine.sp; s != nil {
		value, found, known = s.findCondition(name)
	}
	if !found {
		return &ParseError{Message: "Undefined CONDITION: " + name, Position: pos}
	}
	if known && value.Kind != nodes.HandlerCondSQLState {
		return &ParseError{Message: "SIGNAL/RESIGNAL can only use a CONDITION defined with SQLSTATE", Position: pos}
	}
	return nil
}
