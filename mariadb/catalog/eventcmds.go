package catalog

import (
	"fmt"
	"strings"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// createEvent installs the event stmt defines. authored is true for a
// statement a user runs, false for an event a snapshot already holds, whose
// times MariaDB judged when it was created.
func (c *Catalog) createEvent(stmt *nodes.CreateEventStmt, authored bool) error {
	db, err := c.resolveDatabase(stmt.Schema)
	if err != nil {
		return err
	}

	// MariaDB evaluates ON SCHEDULE expressions before checking whether the
	// event already exists: a bare column reference answers 1054 (measured on
	// mariadb:11.8.8).
	if err := validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}
	if authored {
		if err := validateEventTimes(stmt.Schedule, time.Now()); err != nil {
			return err
		}
	}

	name := stmt.Name
	key := toLower(name)

	if _, exists := db.Events[key]; exists {
		// CREATE OR REPLACE EVENT replaces it (mariadb:11.8.8).
		switch {
		case stmt.OrReplace:
			delete(db.Events, key)
		case stmt.IfNotExists:
			return nil
		default:
			return errDupEvent(name)
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	schedule, at, timed := eventScheduleText(stmt.Schedule)
	enable := stmt.Enable
	// A one-time event whose time has passed is dropped as soon as it is
	// created, or kept DISABLE (DISABLE ON SLAVE stays) ON COMPLETION PRESERVE:
	// measured on mariadb:11.8.8, CREATE EVENT e ON SCHEDULE AT '2000-01-01
	// 00:00:00' DO COMMIT leaves no event.
	if authored && timed && at.Before(time.Now().Add(-eventZoneMargin)) {
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

func (c *Catalog) alterEvent(stmt *nodes.AlterEventStmt) error {
	db, err := c.resolveDatabase(stmt.Schema)
	if err != nil {
		return err
	}

	// MariaDB validates ON SCHEDULE expressions before looking up the event
	// (1054), and rejects renaming an event to its own name (1551) even when
	// the event does not exist (measured on mariadb:11.8.8).
	if err := validateEventSchedule(stmt.Schedule); err != nil {
		return err
	}
	if err := validateEventTimes(stmt.Schedule, time.Now()); err != nil {
		return err
	}

	name := stmt.Name
	key := toLower(name)

	if stmt.RenameTo != "" {
		renameSchema := stmt.RenameSchema
		if renameSchema == "" {
			renameSchema = c.currentDB
		}
		if toLower(renameSchema) == toLower(db.Name) && toLower(stmt.RenameTo) == key {
			return errEventSameName()
		}
	}

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

	// A schedule time that has passed is refused while the event is not kept
	// after it runs (mariadb:11.8.8: ALTER EVENT e ON SCHEDULE AT
	// '2000-01-01 00:00:00' is 1589), before anything changes.
	schedule, at, timed := eventScheduleText(stmt.Schedule)
	if stmt.Schedule != nil {
		onCompletion := event.OnCompletion
		if stmt.OnCompletion != "" {
			onCompletion = stmt.OnCompletion
		}
		if timed && at.Before(time.Now().Add(-eventZoneMargin)) && !strings.EqualFold(onCompletion, "PRESERVE") {
			return &Error{Code: 1589, SQLState: "HY000", Message: "Event execution time is in the past and ON COMPLETION NOT PRESERVE is set. The event was not changed. Specify a time in the future"}
		}
	}

	// Update definer if specified.
	if stmt.Definer != "" {
		event.Definer = stmt.Definer
	}

	// Update schedule if specified.
	if stmt.Schedule != nil {
		event.Schedule = schedule
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

// validateEventSchedule rejects bare column references in ON SCHEDULE
// expressions (AT, EVERY, STARTS, ENDS), which have no tables in scope and
// answer 1054 in MariaDB before checking event existence.
func validateEventSchedule(sched *nodes.EventSchedule) error {
	if sched == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{sched.At, sched.Every, sched.Starts, sched.Ends} {
		if expr == nil {
			continue
		}
		var err error
		nodes.Inspect(expr, func(n nodes.Node) bool {
			if err != nil {
				return false
			}
			switch x := n.(type) {
			case *nodes.SelectStmt:
				return false
			case *nodes.ColumnRef:
				err = errNoSuchColumn(x.Column, "SELECT")
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func showCreateEvent(e *Event) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// DEFINER
	if e.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", formatDefiner(e.Definer)))
	}

	b.WriteString(fmt.Sprintf(" EVENT `%s`", e.Name))

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
