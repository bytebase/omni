package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// What MySQL 8.0 refuses of a trigger's or stored function's body while it parses. The
// action for a column named through NEW or OLD holds it to the trigger's event and time:
// 1363, ER_TRG_NO_SUCH_ROW_IN_TRG, for a row the event has no such thing as, and 1362,
// ER_TRG_CANT_CHANGE_ROW, for assigning to a row that can't change. The action for a
// statement that returns a result set answers 1415, ER_SP_NO_RETSET, in a trigger or a
// function, where a procedure and an event may.
//
// Ref: the NEW_SYM and OLD_SYM handling of simple_ident_q and the sp_proc_stmt_statement arm
// that raises ER_SP_NO_RETSET, in MySQL 8.0 sql/sql_yacc.yy.

// routineScope is what the parser needs to know of the stored program it is inside.
type routineScope struct {
	triggerTiming string // BEFORE or AFTER, "" outside a trigger
	triggerEvent  string // INSERT, UPDATE or DELETE
	noResultSet   bool   // inside a function or a trigger
	function      bool   // inside a function, the one place RETURN may stand (1313 elsewhere)
}

// refuseTriggerRowRead holds NEW.col and OLD.col, in any case and however quoted, to the
// trigger's event.
func (p *Parser) refuseTriggerRowRead(qualifier string, pos int) error {
	if p.routine.triggerEvent == "" {
		return nil
	}
	switch {
	case strings.EqualFold(qualifier, "NEW") && p.routine.triggerEvent == "DELETE":
		return &ParseError{Message: "There is no NEW row in on DELETE trigger", Position: pos}
	case strings.EqualFold(qualifier, "OLD") && p.routine.triggerEvent == "INSERT":
		return &ParseError{Message: "There is no OLD row in on INSERT trigger", Position: pos}
	}
	return nil
}

// refuseTriggerRowWrite holds SET NEW.col and SET OLD.col to it: only a BEFORE trigger's NEW
// row changes.
func (p *Parser) refuseTriggerRowWrite(qualifier string, pos int) error {
	if p.routine.triggerEvent == "" {
		return nil
	}
	if strings.EqualFold(qualifier, "NEW") && p.routine.triggerEvent == "DELETE" {
		return &ParseError{Message: "There is no NEW row in on DELETE trigger", Position: pos}
	}
	if strings.EqualFold(qualifier, "OLD") || p.routine.triggerTiming == "AFTER" {
		return &ParseError{Message: "Updating of " + strings.ToUpper(qualifier) + " row is not allowed in trigger", Position: pos}
	}
	return nil
}

// refuseResultSet is 1415 for a statement of a function's or trigger's body.
func (p *Parser) refuseResultSet(stmt nodes.Node, pos int) error {
	if !p.routine.noResultSet {
		return nil
	}
	returns := false
	switch s := stmt.(type) {
	case *nodes.ShowStmt, *nodes.ExplainStmt:
		returns = true
	case *nodes.SelectStmt:
		returns = !selectHasInto(s)
	}
	if returns {
		return &ParseError{Message: "Not allowed to return a result set from a function or trigger", Position: pos}
	}
	return nil
}

// selectHasInto reports an INTO on any branch of a query.
func selectHasInto(sel *nodes.SelectStmt) bool {
	stack := []*nodes.SelectStmt{sel}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if s == nil {
			continue
		}
		if s.Into != nil {
			return true
		}
		stack = append(stack, s.Left, s.Right, s.ParenSource)
	}
	return false
}
