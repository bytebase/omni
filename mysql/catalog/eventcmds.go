package catalog

import (
	"fmt"
	"strings"
	"time"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func (c *Catalog) createEvent(stmt *nodes.CreateEventStmt) error {
	db, err := c.resolveDatabase(stmt.Schema)
	if err != nil {
		return err
	}

	// The schedule's expressions are read before the event's name is looked up:
	// over an existing event ev, `CREATE EVENT ev ON SCHEDULE AT nosuchcol DO
	// COMMIT` is 1054 and not the 1537 its name answers (mysql:8.0.46).
	if err := c.validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}
	if err := validateEventTimes(stmt.Schedule, time.Now().UTC()); err != nil {
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

	schedule, at, timed := eventScheduleText(stmt.Schedule)
	enable := stmt.Enable
	// A one-time event whose time has passed is dropped as soon as it is
	// created, or kept DISABLE ON COMPLETION PRESERVE (mysql:8.0.46: CREATE
	// EVENT e ON SCHEDULE AT '2000-01-01 00:00:00' DO COMMIT leaves no event).
	if timed && at.Before(time.Now().UTC()) {
		if !strings.EqualFold(stmt.OnCompletion, "PRESERVE") {
			return nil
		}
		if enable == "" || strings.EqualFold(enable, "ENABLE") {
			enable = "DISABLE"
		}
	}

	event := &Event{
		Name:         name,
		Database:     db,
		Definer:      definer,
		Schedule:     schedule,
		OnCompletion: stmt.OnCompletion,
		Enable:       enable,
		Comment:      stmt.Comment,
		Body:         strings.TrimSpace(stmt.BodyText),
	}

	db.Events[key] = event
	return nil
}

// validateEventSchedule refuses a schedule expression that reads a column. An
// event's AT, EVERY, STARTS and ENDS are evaluated on their own, with no table
// in scope, so `AT utf8mb4` names no column and MySQL answers 1054, `Unknown
// column 'utf8mb4' in 'field list'` (mysql:8.0.46). Each expression is read as
// a view body's no-FROM target list is, which is where the code and its message
// come from. CREATE EVENT and ALTER EVENT both answer it, and before the event
// the statement names: `ALTER EVENT nosuchevent ON SCHEDULE AT nosuchcol` is
// 1054 and not the 1537 or 1146 its name would answer.
func (c *Catalog) validateEventSchedule(schedule *nodes.EventSchedule) error {
	if schedule == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{schedule.At, schedule.Every, schedule.Starts, schedule.Ends} {
		if expr == nil || !exprNamesColumn(expr) {
			continue
		}
		sel := &nodes.SelectStmt{TargetList: []nodes.ExprNode{&nodes.ResTarget{Val: expr}}}
		if _, err := c.analyzeViewBody(sel, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) alterEvent(stmt *nodes.AlterEventStmt) error {
	db, err := c.resolveDatabase(stmt.Schema)
	if err != nil {
		return err
	}

	// The schedule is read before the event is looked up, as CREATE EVENT's is.
	if err := c.validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}
	if err := validateEventTimes(stmt.Schedule, time.Now().UTC()); err != nil {
		return err
	}

	// MySQL refuses renaming an event to its own (schema, name) with 1551
	// (ER_EVENT_SAME_NAME) before looking up whether the event exists.
	if stmt.RenameTo != "" {
		renameSchema := stmt.RenameSchema
		if renameSchema == "" {
			if cur, curErr := c.resolveDatabase(""); curErr == nil {
				renameSchema = cur.Name
			}
		}
		if toLower(db.Name) == toLower(renameSchema) && toLower(stmt.Name) == toLower(stmt.RenameTo) {
			return errEventSameName()
		}
	}

	name := stmt.Name
	key := toLower(name)

	event, exists := db.Events[key]
	if !exists {
		return errNoSuchEvent(db.Name, name)
	}

	// The new name is checked before anything changes, so a statement that fails leaves the
	// event as it was.
	var target *Database
	if stmt.RenameTo != "" {
		// A new name with no database goes to the current one, whichever database the event was
		// in: with `a` current, ALTER EVENT b.e RENAME TO f makes a.f (mysql:8.0.46, mariadb:11.8.8).
		if target, err = c.resolveDatabase(stmt.RenameSchema); err != nil {
			return err
		}
		if _, taken := target.Events[toLower(stmt.RenameTo)]; taken {
			return errDupEvent(stmt.RenameTo)
		}
	}

	// Update definer if specified.
	if stmt.Definer != "" {
		event.Definer = stmt.Definer
	}

	// Update schedule if specified.
	if stmt.Schedule != nil {
		event.Schedule, _, _ = eventScheduleText(stmt.Schedule)
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
		event.Body = strings.TrimSpace(stmt.BodyText)
	}

	// Handle RENAME TO.
	if target != nil {
		newKey := toLower(stmt.RenameTo)
		delete(db.Events, key)
		event.Name = stmt.RenameTo
		event.Database = target
		target.Events[newKey] = event
	}

	return nil
}

func (c *Catalog) dropEvent(stmt *nodes.DropEventStmt) error {
	db, err := c.resolveDatabase(stmt.Schema)
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
