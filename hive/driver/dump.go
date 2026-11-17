package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var _ driver.Dumper = (*Driver)(nil)

// dumpNullPartition is the directory name Hive gives a partition whose value is
// null. The server refuses both statements that would state it — the literal
// value (error 10111) and dt=null (error 10248) — so a dump leaves such a
// partition out, and TestDumpGaps holds the loss.
const dumpNullPartition = "__HIVE_DEFAULT_PARTITION__"

// dumpKind is what the server says an object is, which supplies its label and
// whether creating it resolves another object.
type dumpKind int

const (
	dumpTable dumpKind = iota
	dumpView
	dumpMaterializedView
)

// label names the object in the section heading the dump writes.
func (k dumpKind) label() string {
	switch k {
	case dumpView:
		return "View"
	case dumpMaterializedView:
		return "Materialized view"
	default:
		return "Table"
	}
}

// dumpObject is one object a dump writes: the name the server lists, the
// statements that recreate it, and the objects of the same database that
// creating it resolves, which orders them.
type dumpObject struct {
	kind  dumpKind
	name  string
	stmts []string
	reads []string
}

// Dump writes the database's DDL as the server states it. The metadata is not
// rendered and does not narrow the dump: DESCRIBE FORMATTED reports a table's
// type, columns, row count and partitions but not the row format, the input and
// output formats, the serde properties, the storage handler, the skew or bucket
// definition, or the table properties, so a synced schema cannot be turned back
// into HiveQL. SHOW CREATE TABLE states all of those, and DESCRIBE FORMATTED
// states the three things it leaves out — a view's or materialized view's query
// and comment, and whether a table is managed — so the dump reads both. The
// caller creates the database; the dump carries the objects in it, in an order
// that replays.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	objects, err := d.dumpObjects(ctx, d.databaseNameOrDefault())
	if err != nil {
		return err
	}
	for _, o := range orderDumpObjects(objects) {
		if _, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n", o.kind.label(), quoteBacktick(o.name)); err != nil {
			return err
		}
		for _, stmt := range o.stmts {
			if _, err := fmt.Fprintf(out, "%s;\n\n", stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// dumpObjects reads the database's objects and the statements that recreate
// each.
func (d *Driver) dumpObjects(ctx context.Context, database string) ([]*dumpObject, error) {
	names, err := d.dumpNames(ctx, database)
	if err != nil {
		return nil, err
	}
	objects := make([]*dumpObject, 0, len(names))
	for _, name := range names {
		object, err := d.dumpObject(ctx, database, name)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}

	// Creating a view resolves its query, so each view is read for the views it
	// resolves. Nothing else needs reading: a table resolves nothing, and a
	// materialized view comes after every table and view, which is all it can
	// read.
	all := make(map[string]bool, len(objects))
	for _, o := range objects {
		all[o.name] = true
	}
	for _, o := range objects {
		if o.kind != dumpView {
			continue
		}
		reads, err := d.dumpReads(ctx, database, o.name, all)
		if err != nil {
			// Hive answers a view whose source is gone with an error here: the
			// statement it states is still the view's, so the dump carries it,
			// ordered by name like any definition the server cannot place. A
			// call that was cancelled or timed out is not that view.
			if ctx.Err() != nil {
				return nil, err
			}
			continue
		}
		o.reads = reads
	}
	return objects, nil
}

// dumpObject returns what recreates one object: the statement SHOW CREATE TABLE
// states, or the one DESCRIBE FORMATTED completes it to, and for a managed
// table the statements that make the table managed and add its partitions.
func (d *Driver) dumpObject(ctx context.Context, database, name string) (*dumpObject, error) {
	info, err := d.getTableInfo(ctx, name, database)
	if err != nil {
		return nil, fmt.Errorf("failed to describe table %s's type: %w", name, err)
	}
	// Every statement the dump writes qualifies the object with the database,
	// and Hive stores a database name in lower case whatever case it was
	// created with: the statements that complete a CREATE name the database the
	// way DESCRIBE FORMATTED does, so that the dump spells it one way.
	if info.database != "" {
		database = info.database
	}
	object := &dumpObject{name: name}
	switch info.tableType {
	case "VIRTUAL_VIEW":
		object.kind = dumpView
		stmt, err := d.dumpStatement(ctx, database, name)
		if err != nil {
			return nil, err
		}
		object.stmts = []string{stmt}
		// SHOW CREATE TABLE states a view's query and no comment, and CREATE
		// VIEW takes one only before the query: the statement that states a
		// view's comment is the ALTER that sets it as a table parameter.
		if info.comment != "" {
			object.stmts = append(object.stmts, "ALTER VIEW "+quoteBacktick(database)+"."+quoteBacktick(name)+
				" SET TBLPROPERTIES('comment'="+dumpString(info.comment)+")")
		}
	case "MATERIALIZED_VIEW":
		// SHOW CREATE TABLE states a materialized view as a plain CREATE TABLE:
		// its comment and storage, and not the query it materializes. Its
		// clauses are the ones CREATE MATERIALIZED VIEW takes, in that order,
		// but for the column list it refuses; DESCRIBE FORMATTED states the
		// query with every identifier quoted, which is the text a replay takes.
		//
		// Hive refuses CREATE MATERIALIZED VIEW without DISABLE REWRITE when a
		// table the view reads is not transactional, and the metadata carries no
		// field for the flag DESCRIBE FORMATTED states, so the dump writes
		// DISABLE REWRITE for every materialized view.
		stmt, err := d.dumpStatement(ctx, database, name)
		if err != nil {
			return nil, err
		}
		clauses, err := dumpClauseStart(stmt)
		if err != nil {
			return nil, fmt.Errorf("materialized view %s: %w", name, err)
		}
		query := dumpMaterializedViewQuery(info.viewDef, info.expandedViewDef, dumpStatesPartitions(stmt))
		if query == "" {
			return nil, fmt.Errorf("materialized view %s: DESCRIBE FORMATTED states no query", name)
		}
		materialized, err := dumpMaterializedViewClauses(stmt[clauses:])
		if err != nil {
			return nil, fmt.Errorf("materialized view %s: %w", name, err)
		}
		object.kind = dumpMaterializedView
		object.stmts = []string{
			"CREATE MATERIALIZED VIEW " + quoteBacktick(database) + "." + quoteBacktick(name) +
				" DISABLE REWRITE" + materialized + "\nAS " + query,
		}
	case "MANAGED_TABLE":
		object.kind = dumpTable
		stmt, err := d.dumpStatement(ctx, database, name)
		if err != nil {
			return nil, err
		}
		object.stmts = append(object.stmts, stmt)
		// SHOW CREATE TABLE states a managed table as a CREATE TABLE with no
		// clause that makes one: with strict managed tables off, Hive 4 stores a
		// CREATE TABLE as an external table, so the stated CREATE alone restores
		// the table external, and only the ALTER stores it managed.
		object.stmts = append(object.stmts,
			"ALTER TABLE "+quoteBacktick(database)+"."+quoteBacktick(name)+" SET TBLPROPERTIES('EXTERNAL'='FALSE')")
		// The statement declares the partitions, which DESCRIBE FORMATTED
		// states only under a setting a session can turn off.
		if dumpStatesPartitions(stmt) {
			partitions, err := d.dumpPartitions(ctx, database, name, info.storage)
			if err != nil {
				return nil, err
			}
			object.stmts = append(object.stmts, partitions...)
		}
	default:
		object.kind = dumpTable
		stmt, err := d.dumpStatement(ctx, database, name)
		if err != nil {
			return nil, err
		}
		object.stmts = append(object.stmts, stmt)
	}
	return object, nil
}

// dumpClauseStart returns the index in a stated CREATE TABLE of the first
// clause after the column list. A name, a type and a comment may each hold a
// parenthesis, so the list opens at the first parenthesis outside a quoted
// identifier or literal and ends at the one that closes it.
func dumpClauseStart(stmt string) (int, error) {
	depth := 0
	for i := 0; i < len(stmt); i++ {
		switch stmt[i] {
		case '`', '\'':
			end, err := dumpQuotedEnd(stmt, i)
			if err != nil {
				return 0, fmt.Errorf("the stated CREATE TABLE has an %w", err)
			}
			i = end
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i + 1, nil
			}
			if depth < 0 {
				return 0, fmt.Errorf("the stated CREATE TABLE has a parenthesis before its column list")
			}
		}
	}
	return 0, fmt.Errorf("the stated CREATE TABLE has no end to its column list")
}

// dumpQuotedEnd returns the index of the character that ends the quoted
// identifier or literal opening at start.
func dumpQuotedEnd(stmt string, start int) (int, error) {
	switch stmt[start] {
	case '`':
		// Hive spells an identifier's own backquote as two.
		for i := start + 1; i < len(stmt); i++ {
			if stmt[i] != '`' {
				continue
			}
			if i+1 < len(stmt) && stmt[i+1] == '`' {
				i++
				continue
			}
			return i, nil
		}
		return 0, fmt.Errorf("unterminated identifier")
	case '\'':
		// A backslash escapes the character after it, a quote included.
		for i := start + 1; i < len(stmt); i++ {
			if stmt[i] == '\\' {
				i++
				continue
			}
			if stmt[i] == '\'' {
				return i, nil
			}
		}
		return 0, fmt.Errorf("unterminated literal")
	}
	return start, nil
}

// dumpStatesPartitions reports whether a stated CREATE TABLE declares
// partitions. The clause is in the statement, where a session cannot hide it:
// DESCRIBE FORMATTED states a table's partition columns only while
// hive.display.partition.cols.separately is on.
func dumpStatesPartitions(stmt string) bool {
	_, ok := dumpKeyword(stmt, "PARTITIONED BY")
	return ok
}

// dumpKeyword returns where stmt states keyword outside a quoted identifier or
// literal, which a comment or a name may spell.
func dumpKeyword(stmt, keyword string) (int, bool) {
	for i := 0; i < len(stmt); i++ {
		switch stmt[i] {
		case '`', '\'':
			end, err := dumpQuotedEnd(stmt, i)
			if err != nil {
				return 0, false
			}
			i = end
		default:
			if strings.HasPrefix(stmt[i:], keyword) {
				return i, true
			}
		}
	}
	return 0, false
}

// dumpMaterializedViewQuery returns the text a CREATE MATERIALIZED VIEW replays
// as its query. Hive states it two ways: the text the view was created with,
// and an expansion that writes every table and column with its qualifier and
// backquotes and, for a partitioned view, wraps the query in the projection
// Hive reads the partitions through. Replaying the expansion of a partitioned
// view would nest that projection, and replaying the text the server states
// loses a name whose backquotes its normalisation dropped, so the expansion
// stands without its projection where it states such a name.
func dumpMaterializedViewQuery(original, expanded string, partitioned bool) string {
	if !partitioned {
		return expanded
	}
	inner := dumpQueryInProjection(expanded)
	if inner == "" {
		return original
	}
	if original == "" || dumpStatesQuotedName(inner) {
		return inner
	}
	return original
}

// dumpQueryInProjection returns the query an expanded definition wraps in
// Hive's partition projection: the text inside the parenthesis its FROM opens,
// and "" when the definition is not a projection over the query.
func dumpQueryInProjection(expanded string) string {
	text := strings.TrimSpace(expanded)
	if !strings.HasPrefix(strings.ToUpper(text), "SELECT ") {
		return ""
	}
	var open, close int
	depth := 0
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '`', '\'':
			end, err := dumpQuotedEnd(text, i)
			if err != nil {
				return ""
			}
			i = end
		case '(':
			if depth == 0 {
				open = i
			}
			depth++
		case ')':
			depth--
			if depth < 0 {
				return ""
			}
			if depth == 0 {
				close = i
			}
		}
	}
	if close == 0 || !strings.HasSuffix(strings.TrimSpace(text[:open]), "FROM") {
		return ""
	}
	if !dumpStatesAlias(strings.TrimSpace(text[close+1:])) {
		return ""
	}
	return strings.TrimSpace(text[open+1 : close])
}

// dumpStatesAlias reports whether text is the one name a projection states
// after the query it wraps: the view's own name, quoted or not.
func dumpStatesAlias(text string) bool {
	if text == "" {
		return false
	}
	if text[0] == '`' {
		end, err := dumpQuotedEnd(text, 0)
		return err == nil && end == len(text)-1
	}
	return !strings.ContainsAny(text, "`'() \t\n")
}

// dumpStatesQuotedName reports whether text states a backquoted name Hive
// cannot read plain, which is the name a statement's own text may have lost the
// backquotes of.
func dumpStatesQuotedName(text string) bool {
	for i := 0; i < len(text); i++ {
		if text[i] != '`' {
			continue
		}
		end, err := dumpQuotedEnd(text, i)
		if err != nil {
			return true
		}
		if !dumpPlainName(text[i+1 : end]) {
			return true
		}
		i = end
	}
	return false
}

// dumpPlainName reports whether name is one Hive reads without backquotes.
func dumpPlainName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c == '_' || c == '$' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
			continue
		}
		return false
	}
	return true
}

// dumpMaterializedViewClauses rewrites the clauses a stated CREATE TABLE holds
// after its column list into the ones CREATE MATERIALIZED VIEW takes: the
// partition clause names its columns without their types, as PARTITIONED ON.
func dumpMaterializedViewClauses(clauses string) (string, error) {
	at, ok := dumpKeyword(clauses, "PARTITIONED BY")
	if !ok {
		return clauses, nil
	}
	open := strings.IndexByte(clauses[at:], '(')
	if open < 0 {
		return "", fmt.Errorf("the stated PARTITIONED BY has no column list")
	}
	open += at
	end, err := dumpClauseStart(clauses[open:])
	if err != nil {
		return "", err
	}
	names, err := dumpPartitionColumnNames(clauses[open+1 : open+end-1])
	if err != nil {
		return "", err
	}
	return clauses[:at] + "PARTITIONED ON (" + strings.Join(names, ", ") + ")" + clauses[open+end:], nil
}

// dumpPartitionColumnNames reads the names a stated PARTITIONED BY clause lists,
// which CREATE MATERIALIZED VIEW takes without their types.
func dumpPartitionColumnNames(columns string) ([]string, error) {
	var names []string
	var element strings.Builder
	depth := 0
	add := func() error {
		text := strings.TrimSpace(element.String())
		element.Reset()
		if text == "" {
			return nil
		}
		name, err := dumpColumnName(text)
		if err != nil {
			return err
		}
		names = append(names, name)
		return nil
	}
	for i := 0; i < len(columns); i++ {
		switch columns[i] {
		case '`', '\'':
			end, err := dumpQuotedEnd(columns, i)
			if err != nil {
				return nil, err
			}
			element.WriteString(columns[i : end+1])
			i = end
		case '(':
			depth++
			element.WriteByte(columns[i])
		case ')':
			depth--
			element.WriteByte(columns[i])
		case ',':
			if depth != 0 {
				element.WriteByte(columns[i])
				continue
			}
			if err := add(); err != nil {
				return nil, err
			}
		default:
			element.WriteByte(columns[i])
		}
	}
	if err := add(); err != nil {
		return nil, err
	}
	return names, nil
}

// dumpColumnName reads the name a stated partition column begins with.
func dumpColumnName(column string) (string, error) {
	if strings.HasPrefix(column, "`") {
		end, err := dumpQuotedEnd(column, 0)
		if err != nil {
			return "", err
		}
		return column[:end+1], nil
	}
	if end := strings.IndexAny(column, " \t\n"); end >= 0 {
		return column[:end], nil
	}
	return column, nil
}

// dumpString renders s as a Hive string literal. Hive concatenates adjacent
// literals rather than failing on them, so doubling a quote would lose it; the
// backslash escapes the characters that would otherwise end the literal or the
// line it is on, which the splitter cuts a statement at.
func dumpString(s string) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
	escaped = strings.ReplaceAll(escaped, "\r", `\r`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	return "'" + escaped + "'"
}

// dumpNames returns the database's object names, in name order. Every name SHOW
// TABLES answers is an object to dump: Hive 4 parses no CREATE INDEX, so no name
// is one of Hive's own index tables.
func (d *Driver) dumpNames(ctx context.Context, database string) ([]string, error) {
	query := "SHOW TABLES IN " + quoteBacktick(database)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	sort.Strings(names)
	return names, nil
}

// dumpStatement returns the statement SHOW CREATE TABLE states for one object,
// without its terminating semicolon. Hive states it in one row; a server that
// states it a line per row is joined back into one statement.
func (d *Driver) dumpStatement(ctx context.Context, database, name string) (string, error) {
	query := "SHOW CREATE TABLE " + quoteBacktick(database) + "." + quoteBacktick(name)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var line sql.NullString
		if err := rows.Scan(&line); err != nil {
			return "", fmt.Errorf("%s: %w", query, err)
		}
		if line.String != "" {
			lines = append(lines, line.String)
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("%s: no statement", query)
	}
	return strings.TrimSuffix(strings.TrimSpace(strings.Join(lines, "\n")), ";"), nil
}

// dumpPartitions states a managed table's partitions: the names SHOW PARTITIONS
// answers carry the values, and DESCRIBE FORMATTED carries the location each
// one lives at, which a partition added with its own LOCATION does not share
// with the table. A partition whose value is null has no statement the server
// takes and is left out, and one that reads its files otherwise than the table
// does has no statement ADD PARTITION carries either, so it fails the dump
// rather than coming back with the table's storage.
func (d *Driver) dumpPartitions(ctx context.Context, database, name string, table dumpStorage) ([]string, error) {
	query := "SHOW PARTITIONS " + quoteBacktick(database) + "." + quoteBacktick(name)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var partition string
		if err := rows.Scan(&partition); err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		names = append(names, partition)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}

	var stmts []string
	for _, partition := range names {
		spec, ok := dumpPartitionSpec(partition)
		if !ok {
			continue
		}
		location, storage, err := d.dumpPartition(ctx, database, name, spec)
		if err != nil {
			return nil, err
		}
		stmt, err := dumpPartitionStatement(database, name, partition, spec, location, table, storage)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// dumpPartitionStatement states one partition of a managed table where it lives.
// A partition that reads its files otherwise than the table does states nothing
// here: ADD PARTITION carries the table's storage, so the statement would
// register it under a format its files are not in.
func dumpPartitionStatement(database, name, partition, spec, location string, table, storage dumpStorage) (string, error) {
	if diff := dumpStorageDiff(table, storage); diff != "" {
		return "", fmt.Errorf("partition %s of %s.%s states its own %s, which ADD PARTITION does not carry",
			partition, database, name, diff)
	}
	return "ALTER TABLE " + quoteBacktick(database) + "." + quoteBacktick(name) +
		" ADD PARTITION (" + spec + ") LOCATION " + dumpString(location), nil
}

// dumpPartition reads the location a partition lives at, which DESCRIBE
// FORMATTED states for the partition and not for the table, and the storage it
// reads its files with.
func (d *Driver) dumpPartition(ctx context.Context, database, name, spec string) (string, dumpStorage, error) {
	query := "DESCRIBE FORMATTED " + quoteBacktick(database) + "." + quoteBacktick(name) + " PARTITION (" + spec + ")"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", dumpStorage{}, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var storage dumpStorage
	location := ""
	for rows.Next() {
		var colName, dataType, comment sql.NullString
		if err := rows.Scan(&colName, &dataType, &comment); err != nil {
			return "", dumpStorage{}, fmt.Errorf("%s: %w", query, err)
		}
		if strings.TrimSpace(colName.String) == "Location:" {
			location = strings.TrimSpace(dataType.String)
			continue
		}
		dumpStorageField(&storage, strings.TrimSpace(colName.String), strings.TrimSpace(dataType.String))
	}
	if err := rows.Err(); err != nil {
		return "", dumpStorage{}, fmt.Errorf("%s: %w", query, err)
	}
	if location == "" {
		return "", dumpStorage{}, fmt.Errorf("%s: no location", query)
	}
	return location, storage, nil
}

// dumpStorage is what DESCRIBE FORMATTED states about an object's storage: the
// classes Hive reads and writes its files with.
type dumpStorage struct {
	serde        string
	inputFormat  string
	outputFormat string
}

// dumpStorageField reads one row of a DESCRIBE FORMATTED storage section into
// storage, and reports whether the row stated one of the three.
func dumpStorageField(storage *dumpStorage, label, value string) bool {
	switch label {
	case "SerDe Library:":
		storage.serde = value
	case "InputFormat:":
		storage.inputFormat = value
	case "OutputFormat:":
		storage.outputFormat = value
	default:
		return false
	}
	return true
}

// dumpStorageDiff names what a partition's storage states that the table's does
// not, and is "" when the partition reads its files as the table does.
func dumpStorageDiff(table, partition dumpStorage) string {
	var differ []string
	if partition.serde != table.serde {
		differ = append(differ, "SerDe Library "+partition.serde)
	}
	if partition.inputFormat != table.inputFormat {
		differ = append(differ, "InputFormat "+partition.inputFormat)
	}
	if partition.outputFormat != table.outputFormat {
		differ = append(differ, "OutputFormat "+partition.outputFormat)
	}
	return strings.Join(differ, ", ")
}

// dumpPartitionSpec renders a partition name as the specification ADD PARTITION
// takes, and false for a partition whose value is null. Hive names a partition
// by its directory: each level states a column and its value, which Hive
// escapes so that neither the separator nor the level's end can appear in a
// value.
func dumpPartitionSpec(partition string) (string, bool) {
	var pairs []string
	for _, level := range strings.Split(partition, "/") {
		column, value, ok := strings.Cut(level, "=")
		if !ok {
			return "", false
		}
		value = dumpUnescape(value)
		if value == dumpNullPartition {
			return "", false
		}
		pairs = append(pairs, quoteBacktick(column)+"="+dumpString(value))
	}
	return strings.Join(pairs, ", "), true
}

// dumpUnescape reads a value Hive escaped for a directory name. A percent sign
// is escaped too, so every %XX is the character it spells.
func dumpUnescape(value string) string {
	if !strings.Contains(value, "%") {
		return value
	}
	var b strings.Builder
	for i := 0; i < len(value); i++ {
		if value[i] == '%' && i+2 < len(value) {
			if code, err := strconv.ParseUint(value[i+1:i+3], 16, 8); err == nil {
				b.WriteByte(byte(code))
				i += 2
				continue
			}
		}
		b.WriteByte(value[i])
	}
	return b.String()
}

// orderDumpObjects puts the objects in the order a restore replays them. Tables
// resolve nothing, so they come first in name order; a view follows the views it
// resolves; and a materialized view comes last, in name order. Hive refuses a
// view or a materialized view whose definition reads a materialized view (error
// 40000, "View definition references materialized view"), and EXPLAIN
// DEPENDENCY answers a materialized view's own name and not the tables behind
// it — so a materialized view resolves tables and views alone, and every table
// and view is already there.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var tables, views, materializedViews []*dumpObject
	for _, o := range objects {
		switch o.kind {
		case dumpView:
			views = append(views, o)
		case dumpMaterializedView:
			materializedViews = append(materializedViews, o)
		default:
			tables = append(tables, o)
		}
	}
	return append(append(sortDumpObjects(tables), orderDumpViews(views)...), sortDumpObjects(materializedViews)...)
}

// orderDumpViews puts the views in the order a restore replays them: each view
// after every view in its database that creating it resolves. Hive resolves a
// view's query when it creates the view and states no dependency anywhere, so
// the order comes from the statements.
func orderDumpViews(views []*dumpObject) []*dumpObject {
	pending := sortDumpObjects(views)
	names := make(map[string]bool, len(pending))
	for _, o := range pending {
		names[o.name] = true
	}
	emitted := make(map[string]bool, len(pending))
	ordered := make([]*dumpObject, 0, len(pending))
	for len(pending) > 0 {
		progress := false
		for i := 0; i < len(pending); {
			ready := true
			for _, read := range pending[i].reads {
				if names[read] && !emitted[read] {
					ready = false
					break
				}
			}
			if !ready {
				i++
				continue
			}
			ordered = append(ordered, pending[i])
			emitted[pending[i].name] = true
			pending = append(pending[:i], pending[i+1:]...)
			progress = true
		}
		if progress {
			continue
		}
		// A reference cycle cannot exist on the server, which refuses the
		// creation that would close it; were one found anyway, the remaining
		// objects are written in name order so the dump still carries them.
		ordered = append(ordered, pending[0])
		emitted[pending[0].name] = true
		pending = pending[1:]
	}
	return ordered
}

func sortDumpObjects(objects []*dumpObject) []*dumpObject {
	sorted := append([]*dumpObject(nil), objects...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].name < sorted[j].name })
	return sorted
}

// dumpReads returns the objects in names that creating o resolves, as the
// server states them: the objects of this database that the query reaches,
// which is what EXPLAIN DEPENDENCY names. The statement Hive renders cannot be
// read for this: a qualified column of a table aliased with the database's name
// is spelled like a reference to another view, and so is a name inside a string
// literal.
func (d *Driver) dumpReads(ctx context.Context, database, name string, names map[string]bool) ([]string, error) {
	query := "EXPLAIN DEPENDENCY SELECT * FROM " + quoteBacktick(database) + "." + quoteBacktick(name)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var dependencies dumpDependencies
	if rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		if err := json.Unmarshal([]byte(text), &dependencies); err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}

	// Hive states each table as database@table, in whatever case it was created,
	// and the list holds everything the query reaches, through an object in
	// another database included. Only this database's objects order the dump,
	// and every one the query reaches has to be there before it.
	var reads []string
	for _, table := range dependencies.InputTables {
		readDatabase, readName, ok := strings.Cut(table.Name, "@")
		if !ok || !strings.EqualFold(readDatabase, database) || readName == name || !names[readName] {
			continue
		}
		reads = append(reads, readName)
	}
	sort.Strings(reads)
	return reads, nil
}

// dumpDependencies is what EXPLAIN DEPENDENCY answers: every table a query
// reaches, as database@table.
type dumpDependencies struct {
	InputTables []struct {
		Name string `json:"tablename"`
	} `json:"input_tables"`
}
