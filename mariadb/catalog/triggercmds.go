package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mysqlparser "github.com/bytebase/omni/mariadb/parser"
)

// createTrigger installs a trigger. source is the script the statement was
// parsed from, which SHOW CREATE TRIGGER echoes; it is empty where the caller
// has no script, as DefineTrigger does for a synced trigger.
func (c *Catalog) createTrigger(stmt *nodes.CreateTriggerStmt, source string) error {
	// A trigger lives with its table. Its database is the one its name is qualified with, or the
	// current one, and an unqualified table is looked up in it: with `a` current,
	// CREATE TRIGGER b.tr ... ON t reads b.t, and both b.tr ... ON a.t and tr ... ON b.t answer 1435.
	// The table comes first, so a missing one answers 1146 even in the wrong database, and
	// mariadb:11.8.8 answers 1146 for a database that is missing.
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
		return errNoSuchTable(tableSchema, tableName)
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
		SourceText:          triggerSourceDDL(stmt, source),
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

	db.Triggers[key] = trigger
	return nil
}

// triggerSourceDDL returns the part of the CREATE TRIGGER statement that
// SHOW CREATE TRIGGER echoes: the text from the TRIGGER keyword to the end of
// the statement, as written, without the FOLLOWS or PRECEDES clause.
//
// Measured on mariadb:11.8.8. The shown statement is the original one with
// only the definer replaced, so `create   trigger   tr1   before  insert on
// tg` keeps its spacing and its lower case, `` CREATE TRIGGER `tr2` ... ON
// `tg` `` keeps both backticks and `CREATE TRIGGER tr0 ... ON tg` keeps
// neither. What sits before the TRIGGER keyword is dropped, OR REPLACE and a
// comment with it, while IF NOT EXISTS, which sits after it, is kept. The
// ordering clause is dropped wherever it stood.
func triggerSourceDDL(stmt *nodes.CreateTriggerStmt, source string) string {
	start, end := stmt.Loc.Start, stmt.Loc.End
	if source == "" || start < 0 || start >= end || end > len(source) {
		return ""
	}
	text := source[start:end]
	kw := triggerKeywordOffset(text)
	if kw < 0 {
		return ""
	}
	if stmt.Order == nil {
		return strings.TrimSpace(text[kw:])
	}
	// The ordering clause is a span of the same statement, so its offsets are
	// relative to the script and come back to this slice.
	from, to := stmt.Order.Loc.Start-start, stmt.Order.Loc.End-start
	if from < kw || from >= to || to > len(text) {
		return strings.TrimSpace(text[kw:])
	}
	before := strings.TrimSpace(text[kw:from])
	after := strings.TrimSpace(text[to:])
	if after == "" {
		return before
	}
	return before + " " + after
}

// triggerKeywordOffset returns the offset of the TRIGGER keyword in a CREATE
// TRIGGER statement, or -1. The statement is tokenized rather than searched,
// so a definer or a body holding the word does not match: only a bare
// `trigger` token before the name can be the keyword, since the word is
// reserved and any other use of it has to be quoted.
func triggerKeywordOffset(text string) int {
	for _, tok := range mysqlparser.Tokenize(text) {
		if tok.Loc < 0 || tok.End <= tok.Loc || tok.End > len(text) {
			continue
		}
		word := text[tok.Loc:tok.End]
		switch word[0] {
		case '`', '\'', '"':
			continue
		}
		if strings.EqualFold(word, "trigger") {
			return tok.Loc
		}
	}
	return -1
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

	// MariaDB echoes the statement the trigger was created by, so where this
	// catalog holds it, that is the answer: it carries the writer's own
	// spacing, case and quoting, down to the backticks they did or did not
	// write. See triggerSourceDDL.
	if tr.SourceText != "" {
		b.WriteString(" ")
		b.WriteString(tr.SourceText)
		return b.String()
	}

	// A synced trigger has no statement to echo, so the DDL is rebuilt from
	// what the catalog holds. Identifiers are backticked, which the writer may
	// not have done, because a name that needs them makes the output
	// unparseable without them: a trigger `my trigger` on a table `order`
	// would otherwise come out as `TRIGGER my trigger ... ON order`.
	b.WriteString(fmt.Sprintf(" TRIGGER %s %s %s ON %s FOR EACH ROW",
		backtickIDQ(tr.Name), tr.Timing, tr.Event, backtickIDQ(tr.Table)))

	// Note: SHOW CREATE TRIGGER does NOT include FOLLOWS/PRECEDES.

	// Body
	if tr.Body != "" {
		b.WriteString(fmt.Sprintf(" %s", tr.Body))
	}

	return b.String()
}
