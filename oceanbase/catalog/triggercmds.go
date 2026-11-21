package catalog

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

func (c *Catalog) createTrigger(stmt *nodes.CreateTriggerStmt) error {
	// A trigger named with a database lives there, and so must its table. OceanBase
	// 4.3.5.6 answers 1049 for a database that is missing, before it looks for the
	// table, and 1435 where the table is in another one. The table's database is its own
	// or the current one, never the trigger's: with b current, CREATE TRIGGER a.tr ... ON
	// t reads b.t and is a 1435, where MySQL reads a.t.
	var named *Database
	if stmt.Schema != "" {
		if named = c.GetDatabase(stmt.Schema); named == nil {
			return errUnknownDatabase(stmt.Schema)
		}
	}

	if err := rejectRemoteRef(stmt.Table); err != nil {
		return err
	}

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
	if named != nil && named != db {
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

	b.WriteString(fmt.Sprintf(" TRIGGER %s %s %s ON %s FOR EACH ROW",
		quoteIdent(tr.Name), tr.Timing, tr.Event, quoteIdent(tr.Table)))

	// Note: MySQL 8.0 SHOW CREATE TRIGGER does NOT include FOLLOWS/PRECEDES.

	// Body
	if tr.Body != "" {
		b.WriteString(fmt.Sprintf(" %s", tr.Body))
	}

	return b.String()
}
