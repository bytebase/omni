package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mysqlparser "github.com/bytebase/omni/mariadb/parser"
)

// createTrigger installs a trigger. source is the script the statement was
// parsed from, which SHOW CREATE TRIGGER echoes; it is empty where the caller
// has no script, as DefineTrigger does for a synced trigger. A synced trigger,
// and one a declarative schema states, is installed as the server holds it:
// its table may have lost a column it reads since it was made, and the trigger
// it follows may be installed after it, so neither is checked.
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

	// CREATE OR REPLACE replaces a trigger of the name (mariadb:11.8.8).
	if _, exists := db.Triggers[key]; exists && !stmt.OrReplace {
		if !stmt.IfNotExists {
			return errDupTrigger(name)
		}
		return nil
	}
	authored := source != "" && !c.loadingSDL
	if authored {
		if err := checkTriggerRowColumns(stmt.Body, db.GetTable(tableName)); err != nil {
			return err
		}
	}
	// FOLLOWS and PRECEDES name a trigger of the same table, time and event
	// (4031, ER_REFERENCED_TRG_DOES_NOT_EXIST, mariadb:11.8.8).
	if o := stmt.Order; o != nil && authored {
		ref := db.Triggers[toLower(o.TriggerName)]
		if ref == nil || toLower(ref.Table) != toLower(tableName) ||
			!strings.EqualFold(ref.Timing, stmt.Timing) || !strings.EqualFold(ref.Event, stmt.Event) {
			return &Error{Code: 4031, SQLState: "HY000",
				Message: fmt.Sprintf("Referenced trigger '%s' for the given action time and event type does not exist", o.TriggerName)}
		}
	}
	delete(db.Triggers, key)

	definer := programDefiner(stmt.Definer)

	c.triggersCreated++
	sourceText, tableSpan := triggerSourceDDL(stmt, source)
	trigger := &Trigger{
		created:             c.triggersCreated,
		Name:                name,
		Database:            db,
		Table:               tableName,
		Timing:              stmt.Timing,
		Event:               stmt.Event,
		Definer:             definer,
		Body:                strings.TrimSpace(stmt.BodyText),
		SourceText:          sourceText,
		tableSpan:           tableSpan,
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

// dropTableTriggers drops the triggers of a table db drops, as DROP TABLE
// does (mariadb:11.8.8).
func dropTableTriggers(db *Database, table string) {
	for name, trg := range db.Triggers {
		if toLower(trg.Table) == toLower(table) {
			delete(db.Triggers, name)
		}
	}
}

// renameTableTriggers moves the triggers of a table renamed within its
// database to the new name. MariaDB writes the name into the statement SHOW
// CREATE TRIGGER echoes, quoted: measured on mariadb:11.8.8, a trigger written
// `... ON t FOR EACH ROW ...` reads `... ON `+"`t2`"+` FOR EACH ROW ...` after
// RENAME TABLE t TO t2.
func renameTableTriggers(db *Database, oldName, newName string) {
	for _, trg := range db.Triggers {
		if toLower(trg.Table) != toLower(oldName) {
			continue
		}
		trg.Table = newName
		if s, e := trg.tableSpan[0], trg.tableSpan[1]; e > s && e <= len(trg.SourceText) {
			quoted := backtickIDQ(newName)
			trg.SourceText = trg.SourceText[:s] + quoted + trg.SourceText[e:]
			trg.tableSpan[1] = s + len(quoted)
		}
	}
}

// checkTriggerRowColumns holds NEW.col and OLD.col in a trigger's body to the
// columns of its table: 1054, measured on mariadb:11.8.8 for SET NEW.zzz = 1.
// A query in the body is left alone, since a table it reads may be named NEW.
func checkTriggerRowColumns(body nodes.Node, tbl *Table) error {
	if body == nil || tbl == nil {
		return nil
	}
	var err error
	nodes.Inspect(body, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.SelectStmt:
			return false
		case *nodes.ColumnRef:
			row := strings.ToUpper(x.Table)
			if x.Schema == "" && !x.Star && (row == "NEW" || row == "OLD") && tbl.GetColumn(x.Column) == nil {
				err = errNoSuchColumn(x.Column, row)
			}
		}
		return true
	})
	return err
}

// triggerSourceDDL returns the part of the CREATE TRIGGER statement that
// SHOW CREATE TRIGGER echoes: the text from the TRIGGER keyword to the end of
// the statement, as written, without the FOLLOWS or PRECEDES clause, and where
// the table's name stands in it, which a rename of the table rewrites.
//
// Measured on mariadb:11.8.8. The shown statement is the original one with
// only the definer replaced, so a statement written
// "create   trigger   tr1   before  insert on   tg" keeps its spacing and its
// lower case, one written with backticks around the trigger and the table
// keeps both, and one written without them gains neither. What sits before
// the TRIGGER keyword is dropped, OR REPLACE and a comment with it, while
// IF NOT EXISTS, which sits after it, is kept. The ordering clause is dropped
// wherever it stood.
func triggerSourceDDL(stmt *nodes.CreateTriggerStmt, source string) (string, [2]int) {
	start, end := stmt.Loc.Start, stmt.Loc.End
	if source == "" || start < 0 || start >= end || end > len(source) {
		return "", [2]int{}
	}
	text := source[start:end]
	kw := triggerKeywordOffset(text)
	if kw < 0 {
		return "", [2]int{}
	}
	// The table's name comes before the ordering clause, so its place in the
	// text from the keyword on holds whether or not the clause is dropped.
	var span [2]int
	if t := stmt.Table; t != nil && t.Loc.Start-start > kw && t.Loc.End > t.Loc.Start && t.Loc.End-start <= len(text) {
		span = [2]int{t.Loc.Start - start - kw, t.Loc.End - start - kw}
	}
	return triggerSourceText(stmt, text, start, kw), span
}

func triggerSourceText(stmt *nodes.CreateTriggerStmt, text string, start, kw int) string {
	if stmt.Order == nil {
		return strings.TrimSpace(text[kw:])
	}
	// The ordering clause is a span of the same statement, so its offsets are
	// relative to the script and come back to this slice.
	from, to := stmt.Order.Loc.Start-start, stmt.Order.Loc.End-start
	if from < kw || from >= to || to > len(text) {
		return strings.TrimSpace(text[kw:])
	}
	// The text before the clause is kept as it was written, whitespace and
	// all: trimming it would take the newline that ends a line comment with
	// it and comment out the body. What goes is the clause and the whitespace
	// after it, which is what the server drops: measured on mariadb:11.8.8,
	// `FOR EACH ROW\nFOLLOWS tr1\nSET NEW.a = 2` comes back as
	// `FOR EACH ROW\nSET NEW.a = 2`.
	before := text[kw:from]
	after := strings.TrimLeft(text[to:], " \t\r\n")
	if after == "" {
		return strings.TrimSpace(before)
	}
	if n := len(before); n > 0 && !isSpaceByte(before[n-1]) {
		before += " "
	}
	return strings.TrimSpace(before + after)
}

// isSpaceByte reports whether c is one of the bytes that separate SQL tokens.
func isSpaceByte(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n', '\v', '\f':
		return true
	}
	return false
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
		// A trigger is found by its file in the database's directory, and a
		// missing directory holds none: DROP TRIGGER nosuch.tr is 1360
		// (mariadb:11.8.8).
		if schema != "" {
			return errNoSuchTrigger(schema, stmt.Name.Name)
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
		b.WriteString(fmt.Sprintf(" DEFINER=%s", formatDefiner(tr.Definer)))
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
