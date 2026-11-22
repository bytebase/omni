package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

func (c *Catalog) createEvent(stmt *nodes.CreateEventStmt) error {
	db, err := c.resolveDatabase("")
	if err != nil {
		return err
	}

	// OceanBase 4.3.5.6 validates schedule timestamp expressions (AT, STARTS,
	// ENDS) before checking whether the event name already exists, answering
	// 1525 (Incorrect AT value) for a bare identifier.
	if err := validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}

	name := stmt.Name
	key := toLower(name)

	if _, exists := db.Events[key]; exists {
		if !stmt.IfNotExists {
			return errDupEvent(name)
		}
		return nil
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Extract raw schedule text from the AST.
	schedule := ""
	if stmt.Schedule != nil {
		schedule = stmt.Schedule.RawText
	}

	event := &Event{
		Name:         name,
		Database:     db,
		Definer:      definer,
		Schedule:     schedule,
		OnCompletion: stmt.OnCompletion,
		Enable:       stmt.Enable,
		Comment:      stmt.Comment,
		Body:         normalizeEventBody(stmt.Body, stmt.BodyText),
		started:      startsWhenMade(stmt.Schedule),
	}

	db.Events[key] = event
	return nil
}

// normalizeEventBody trims an event body and unwraps a single-statement
// unlabeled BEGIN <stmt>; END wrapper, matching OceanBase 4.3.5.6's
// information_schema.EVENTS.EVENT_DEFINITION storage.
func normalizeEventBody(body nodes.Node, bodyText string) string {
	trimmed := strings.TrimSpace(bodyText)
	blk, ok := body.(*nodes.BeginEndBlock)
	if !ok || blk.Label != "" || len(blk.Stmts) != 1 {
		return trimmed
	}
	if len(trimmed) < len("BEGIN")+len("END") {
		return trimmed
	}
	if !strings.EqualFold(trimmed[:len("BEGIN")], "BEGIN") || !strings.EqualFold(trimmed[len(trimmed)-len("END"):], "END") {
		return trimmed
	}
	inner := strings.TrimSpace(trimmed[len("BEGIN") : len(trimmed)-len("END")])
	return strings.TrimSpace(strings.TrimSuffix(inner, ";"))
}

// validateEventSchedule checks that schedule timestamp expressions (AT, STARTS,
// ENDS) do not reference bare columns. OceanBase 4.3.5.6 evaluates them before
// event name lookup and answers 1525 (Incorrect AT value: '<col>').
func validateEventSchedule(schedule *nodes.EventSchedule) error {
	if schedule == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{schedule.At, schedule.Starts, schedule.Ends} {
		if expr == nil {
			continue
		}
		if col, ok := firstColumnRef(expr); ok {
			return errEventWrongValue(col)
		}
	}
	return nil
}

func firstColumnRef(expr nodes.ExprNode) (string, bool) {
	if expr == nil {
		return "", false
	}
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return e.Column, true
	case *nodes.ParenExpr:
		return firstColumnRef(e.Expr)
	case *nodes.UnaryExpr:
		return firstColumnRef(e.Operand)
	case *nodes.BinaryExpr:
		if col, ok := firstColumnRef(e.Left); ok {
			return col, true
		}
		return firstColumnRef(e.Right)
	}
	return "", false
}

func (c *Catalog) alterEvent(stmt *nodes.AlterEventStmt) error {
	db, err := c.resolveDatabase("")
	if err != nil {
		return err
	}

	// Validate schedule expressions before looking up the event, matching
	// OceanBase 4.3.5.6.
	if err := validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}

	// OceanBase 4.3.5.6 refuses renaming an event to its own name with 1551
	// (Same old and new event name) before checking whether the event exists.
	if stmt.RenameTo != "" && toLower(stmt.Name) == toLower(stmt.RenameTo) {
		return errEventSameName()
	}

	name := stmt.Name
	key := toLower(name)

	event, exists := db.Events[key]
	if !exists {
		return errNoSuchEvent(db.Name, name)
	}
	// OceanBase renames no event that has started, enabled or not, where MySQL renames
	// any (TestOracle_OceanBaseRefusals).
	if stmt.RenameTo != "" && event.started {
		return &Error{Code: 1235, SQLState: "0A000", Message: "modify EVENT NAME that has already started not supported"}
	}
	if stmt.RenameTo != "" {
		if _, taken := db.Events[toLower(stmt.RenameTo)]; taken {
			return errDupEvent(stmt.RenameTo)
		}
	}

	// Update definer if specified.
	if stmt.Definer != "" {
		event.Definer = stmt.Definer
	}

	// Update schedule if specified.
	if stmt.Schedule != nil {
		event.Schedule = stmt.Schedule.RawText
		event.started = startsWhenMade(stmt.Schedule)
	}

	// Update ON COMPLETION if specified.
	if stmt.OnCompletion != "" {
		event.OnCompletion = stmt.OnCompletion
	}

	// Update enable/disable if specified.
	if stmt.Enable != "" {
		event.Enable = stmt.Enable
	}

	// Update comment if specified.
	if stmt.Comment != "" {
		event.Comment = stmt.Comment
	}

	// Update body if specified.
	if stmt.BodyText != "" {
		event.Body = normalizeEventBody(stmt.Body, stmt.BodyText)
	}

	// Handle RENAME TO.
	if stmt.RenameTo != "" {
		newKey := toLower(stmt.RenameTo)
		delete(db.Events, key)
		event.Name = stmt.RenameTo
		db.Events[newKey] = event
	}

	return nil
}

func (c *Catalog) dropEvent(stmt *nodes.DropEventStmt) error {
	db, err := c.resolveDatabase("")
	if err != nil {
		if stmt.IfExists {
			return nil
		}
		return err
	}

	name := stmt.Name
	key := toLower(name)

	if _, exists := db.Events[key]; !exists {
		if stmt.IfExists {
			return nil
		}
		return errNoSuchEvent(db.Name, name)
	}

	delete(db.Events, key)
	return nil
}

// ShowCreateEvent produces MySQL 8.0-compatible SHOW CREATE EVENT output.
//
// MySQL 8.0 SHOW CREATE EVENT format:
//
//	CREATE DEFINER=`root`@`%` EVENT `event_name` ON SCHEDULE schedule ON COMPLETION [NOT] PRESERVE [ENABLE|DISABLE|DISABLE ON SLAVE] [COMMENT 'string'] DO event_body
func (c *Catalog) ShowCreateEvent(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	event := db.Events[toLower(name)]
	if event == nil {
		return ""
	}
	return showCreateEvent(event)
}

func showCreateEvent(e *Event) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// DEFINER
	if e.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", e.Definer))
	}

	b.WriteString(fmt.Sprintf(" EVENT %s", quoteIdent(e.Name)))

	// ON SCHEDULE
	if e.Schedule != "" {
		b.WriteString(fmt.Sprintf(" ON SCHEDULE %s", e.Schedule))
	}

	// ON COMPLETION
	if e.OnCompletion == "NOT PRESERVE" {
		b.WriteString(" ON COMPLETION NOT PRESERVE")
	} else if e.OnCompletion == "PRESERVE" {
		b.WriteString(" ON COMPLETION PRESERVE")
	} else {
		// MySQL default: NOT PRESERVE, shown explicitly in SHOW CREATE EVENT
		b.WriteString(" ON COMPLETION NOT PRESERVE")
	}

	// ENABLE / DISABLE
	if e.Enable == "DISABLE" {
		b.WriteString(" DISABLE")
	} else if e.Enable == "DISABLE ON SLAVE" {
		b.WriteString(" DISABLE ON SLAVE")
	} else {
		// MySQL default: ENABLE, shown explicitly in SHOW CREATE EVENT
		b.WriteString(" ENABLE")
	}

	// COMMENT
	if e.Comment != "" {
		b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(e.Comment)))
	}

	// DO event_body
	if e.Body != "" {
		b.WriteString(fmt.Sprintf(" DO %s", e.Body))
	}

	return b.String()
}

// startsWhenMade reports whether a schedule starts its event at once: EVERY with no STARTS.
func startsWhenMade(schedule *nodes.EventSchedule) bool {
	return schedule != nil && schedule.Every != nil && schedule.Starts == nil
}
