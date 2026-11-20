package catalog

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func (c *Catalog) createTrigger(stmt *nodes.CreateTriggerStmt) error {
	// A trigger lives with its table. Its database is the one its name is qualified with, or the
	// current one, and an unqualified table is looked up in it: with `a` current,
	// CREATE TRIGGER b.tr ... ON t reads b.t, and both b.tr ... ON a.t and tr ... ON b.t answer 1435.
	// The table comes first, so a missing one answers 1146 even in the wrong database, and
	// mysql:8.0.46 answers 1049 for a database that is missing.
	triggerSchema := stmt.Schema
	if triggerSchema == "" {
		current, err := c.resolveDatabase("")
		if err != nil {
			return err
		}
		triggerSchema = current.Name
	}
	tableSchema, tableName := triggerSchema, ""
	if stmt.Table != nil {
		tableName = stmt.Table.Name
		if stmt.Table.Schema != "" {
			tableSchema = stmt.Table.Schema
		}
	}
	db := c.GetDatabase(tableSchema)
	if db == nil {
		_, err := c.resolveDatabase(tableSchema)
		return err
	}
	if tableName != "" && db.GetTable(tableName) == nil {
		return errNoSuchTable(db.Name, tableName)
	}
	if toLower(tableSchema) != toLower(triggerSchema) {
		return errTriggerInWrongSchema()
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
	if err := validateTriggerBody(stmt.Timing, stmt.Event, stmt.BodyText); err != nil {
		return err
	}

	trigger := &Trigger{
		Name:                name,
		Database:            db,
		Table:               tableName,
		Timing:              stmt.Timing,
		Event:               stmt.Event,
		Definer:             definer,
		Body:                strings.TrimSpace(stmt.BodyText),
		CharacterSetClient:  c.charsetClient,
		CollationConnection: c.collationConn,
		DatabaseCollation:   db.Collation,
	}

	if stmt.Order != nil {
		trigger.Order = &TriggerOrderInfo{
			Follows:     stmt.Order.Follows,
			TriggerName: stmt.Order.TriggerName,
		}
	}

	assignTriggerActionOrder(db, trigger)
	db.Triggers[key] = trigger
	return nil
}

func validateTriggerBody(timing, event, body string) error {
	lowerBody := strings.ToLower(strings.TrimSpace(body))
	lowerEvent := strings.ToLower(event)
	lowerTiming := strings.ToLower(timing)
	switch lowerEvent {
	case "insert":
		if strings.Contains(lowerBody, "old.") {
			return &Error{Code: 1363, SQLState: "HY000", Message: "There is no OLD row in on INSERT trigger"}
		}
	case "delete":
		if strings.Contains(lowerBody, "new.") {
			return &Error{Code: 1363, SQLState: "HY000", Message: "There is no NEW row in on DELETE trigger"}
		}
	}
	if lowerTiming == "after" && strings.HasPrefix(lowerBody, "set new.") {
		return &Error{Code: 1362, SQLState: "HY000", Message: "Updating of NEW row is not allowed in after trigger"}
	}
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

	tr, exists := db.Triggers[key]
	if !exists {
		if stmt.IfExists {
			return nil
		}
		return errNoSuchTrigger(db.Name, name)
	}

	delete(db.Triggers, key)
	if tr != nil {
		renumberTriggerGroup(db, tr.Table, tr.Timing, tr.Event)
	}
	return nil
}

func triggerGroupPeers(db *Database, tableName, timing, event string) []*Trigger {
	var peers []*Trigger
	for _, tr := range db.Triggers {
		if tr == nil {
			continue
		}
		if strings.EqualFold(tr.Table, tableName) && strings.EqualFold(tr.Timing, timing) && strings.EqualFold(tr.Event, event) {
			peers = append(peers, tr)
		}
	}
	slices.SortFunc(peers, func(a, b *Trigger) int {
		return cmp.Or(cmp.Compare(a.ActionOrder, b.ActionOrder), cmp.Compare(a.Name, b.Name))
	})
	return peers
}

func assignTriggerActionOrder(db *Database, trigger *Trigger) {
	peers := triggerGroupPeers(db, trigger.Table, trigger.Timing, trigger.Event)
	insertIdx := len(peers)
	if trigger.Order != nil {
		for i, peer := range peers {
			if strings.EqualFold(peer.Name, trigger.Order.TriggerName) {
				if trigger.Order.Follows {
					insertIdx = i + 1
				} else {
					insertIdx = i
				}
				break
			}
		}
	}
	peers = slices.Insert(peers, insertIdx, trigger)
	for i, peer := range peers {
		peer.ActionOrder = i + 1
	}
}

func renumberTriggerGroup(db *Database, tableName, timing, event string) {
	for i, peer := range triggerGroupPeers(db, tableName, timing, event) {
		peer.ActionOrder = i + 1
	}
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
