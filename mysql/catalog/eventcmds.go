package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func (c *Catalog) createEvent(stmt *nodes.CreateEventStmt) error {
	db, err := c.resolveDatabase("")
	if err != nil {
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
		Body:         strings.TrimSpace(stmt.BodyText),
	}

	db.Events[key] = event
	return nil
}

func (c *Catalog) alterEvent(stmt *nodes.AlterEventStmt) error {
	db, err := c.resolveDatabase("")
	if err != nil {
		return err
	}

	name := stmt.Name
	key := toLower(name)

	event, exists := db.Events[key]
	if !exists {
		return errNoSuchEvent(db.Name, name)
	}

	// Update definer if specified.
	if stmt.Definer != "" {
		event.Definer = stmt.Definer
	}

	// Update schedule if specified.
	if stmt.Schedule != nil {
		event.Schedule = stmt.Schedule.RawText
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
