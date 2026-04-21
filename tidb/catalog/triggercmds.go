package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func (c *Catalog) createTrigger(stmt *nodes.CreateTriggerStmt) error {
	// Resolve database from the table reference.
	schema := ""
	if stmt.Table != nil {
		schema = stmt.Table.Schema
	}
	db, err := c.resolveDatabase(schema)
	if err != nil {
		return err
	}

	// Verify the table exists.
	tableName := ""
	if stmt.Table != nil {
		tableName = stmt.Table.Name
	}
	if tableName != "" {
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return errNoSuchTable(db.Name, tableName)
		}
	}

	name := stmt.Name
	key := toLower(name)

	if _, exists := db.Triggers[key]; exists {
		if !stmt.IfNotExists {
			return errDupTrigger(name)
		}
		return nil
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	trigger := &Trigger{
		Name:     name,
		Database: db,
		Table:    tableName,
		Timing:   stmt.Timing,
		Event:    stmt.Event,
		Definer:  definer,
		Body:     strings.TrimSpace(stmt.Body),
	}

	if stmt.Order != nil {
		trigger.Order = &TriggerOrderInfo{
			Follows:     stmt.Order.Follows,
			TriggerName: stmt.Order.TriggerName,
		}
	}

	db.Triggers[key] = trigger
	return nil
}

func (c *Catalog) dropTrigger(stmt *nodes.DropTriggerStmt) error {
	schema := ""
	if stmt.Name != nil {
		schema = stmt.Name.Schema
	}
	db, err := c.resolveDatabase(schema)
	if err != nil {
		if stmt.IfExists {
			return nil
		}
		return err
	}

	name := ""
	if stmt.Name != nil {
		name = stmt.Name.Name
	}
	key := toLower(name)

	if _, exists := db.Triggers[key]; !exists {
		if stmt.IfExists {
			return nil
		}
		return errNoSuchTrigger(db.Name, name)
	}

	delete(db.Triggers, key)
	return nil
}

// ShowCreateTrigger produces MySQL 8.0-compatible SHOW CREATE TRIGGER output.
//
// MySQL 8.0 SHOW CREATE TRIGGER format:
//
//	CREATE DEFINER=`root`@`%` TRIGGER `trigger_name` BEFORE INSERT ON `table_name` FOR EACH ROW trigger_body
func (c *Catalog) ShowCreateTrigger(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	trigger := db.Triggers[toLower(name)]
	if trigger == nil {
		return ""
	}
	return showCreateTrigger(trigger)
}

func showCreateTrigger(tr *Trigger) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// DEFINER
	if tr.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", tr.Definer))
	}

	b.WriteString(fmt.Sprintf(" TRIGGER `%s` %s %s ON `%s` FOR EACH ROW",
		tr.Name, tr.Timing, tr.Event, tr.Table))

	// Note: MySQL 8.0 SHOW CREATE TRIGGER does NOT include FOLLOWS/PRECEDES.

	// Body
	if tr.Body != "" {
		b.WriteString(fmt.Sprintf(" %s", tr.Body))
	}

	return b.String()
}
