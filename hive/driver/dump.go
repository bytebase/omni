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
		stmt, err := d.dumpViewStatement(ctx, database, name, info)
		if err != nil {
			return nil, err
		}
		object.stmts = []string{stmt}
	case "MATERIALIZED_VIEW":
		// SHOW CREATE TABLE states a materialized view as a plain CREATE TABLE:
		// its columns, storage and comment, not the query it materializes.
		// DESCRIBE FORMATTED states the query, and only CREATE MATERIALIZED VIEW
		// takes it.
		//
		// Whenever Hive refuses to enable rewriting on a table no transaction
		// manager covers — every table of a server without ACID — it refuses
		// CREATE MATERIALIZED VIEW without DISABLE REWRITE, and the metadata
		// carries no field for the flag DESCRIBE FORMATTED states, so the dump
		// writes DISABLE REWRITE for every materialized view.
		object.kind = dumpMaterializedView
		object.stmts = []string{
			"CREATE MATERIALIZED VIEW " + quoteBacktick(database) + "." + quoteBacktick(name) +
				" DISABLE REWRITE" + dumpComment(info.comment) + " AS " + info.viewDef,
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
		partitions, err := d.dumpPartitions(ctx, database, name)
		if err != nil {
			return nil, err
		}
		object.stmts = append(object.stmts, partitions...)
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

// dumpViewStatement is the statement that recreates a view: the one SHOW CREATE
// TABLE states, which is the view's expanded query, unless the view has a
// comment — SHOW CREATE TABLE states none, and CREATE VIEW takes one only
// before AS.
func (d *Driver) dumpViewStatement(ctx context.Context, database, name string, info *TableInfo) (string, error) {
	stmt, err := d.dumpStatement(ctx, database, name)
	if err != nil {
		return "", err
	}
	if info.comment == "" {
		return stmt, nil
	}
	return "CREATE VIEW " + quoteBacktick(database) + "." + quoteBacktick(name) +
		dumpComment(info.comment) + " AS " + info.viewDef, nil
}

// dumpComment states a comment as CREATE VIEW and CREATE MATERIALIZED VIEW take
// one, or nothing for an object without.
func dumpComment(comment string) string {
	if comment == "" {
		return ""
	}
	return " COMMENT " + dumpString(comment)
}

// dumpString renders s as a Hive string literal. Hive concatenates adjacent
// literals rather than failing on them, so doubling a quote would lose it; the
// backslash escapes the two characters that would otherwise end the literal.
func dumpString(s string) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
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
// takes and is left out.
func (d *Driver) dumpPartitions(ctx context.Context, database, name string) ([]string, error) {
	query := "SHOW PARTITIONS " + quoteBacktick(database) + "." + quoteBacktick(name)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		// A table that is not partitioned answers with an error, which is the
		// same as having no partitions.
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		return nil, nil
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
		location, err := d.dumpPartitionLocation(ctx, database, name, spec)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, "ALTER TABLE "+quoteBacktick(database)+"."+quoteBacktick(name)+
			" ADD PARTITION ("+spec+") LOCATION "+dumpString(location))
	}
	return stmts, nil
}

// dumpPartitionLocation reads the location a partition lives at, which
// DESCRIBE FORMATTED states for the partition and not for the table.
func (d *Driver) dumpPartitionLocation(ctx context.Context, database, name, spec string) (string, error) {
	query := "DESCRIBE FORMATTED " + quoteBacktick(database) + "." + quoteBacktick(name) + " PARTITION (" + spec + ")"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	location := ""
	for rows.Next() {
		var colName, dataType sql.NullString
		if err := rows.Scan(&colName, &dataType, new(sql.NullString)); err != nil {
			return "", fmt.Errorf("%s: %w", query, err)
		}
		if strings.TrimSpace(colName.String) == "Location:" {
			location = strings.TrimSpace(dataType.String)
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	if location == "" {
		return "", fmt.Errorf("%s: no location", query)
	}
	return location, nil
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
