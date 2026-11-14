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

// spScope is one sp_pcontext.
type spScope struct {
	parent  *spScope
	handler bool // a handler's statement, sp_pcontext::HANDLER_SCOPE
	labels  []spLabel
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

// spChecks reports whether the parser holds a stored program's names to its scope: inside
// one, and not while collecting completion candidates, where the text is unfinished.
func (p *Parser) spChecks() bool {
	return p.routine.sp != nil && !p.completing
}

// enterStoredProgram starts the scope of a routine, trigger or event body. The returned
// function restores the scope around it.
func (p *Parser) enterStoredProgram(r routineScope) func() {
	outer := p.routine
	r.sp = &spScope{}
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
