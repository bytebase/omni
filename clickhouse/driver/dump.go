package driver

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var _ driver.Dumper = (*Driver)(nil)

// The object kinds a dump names in its headers.
const (
	dumpDictionary       = "dictionary"
	dumpDictionaryTable  = "dictionary table"
	dumpTable            = "table"
	dumpView             = "view"
	dumpMaterializedView = "materialized view"
)

// dumpObject is one object a dump writes: the kind its statement makes it, the
// name the server lists, the statement the server states, and the materialized
// views in the same database that statement reads, which orders them.
type dumpObject struct {
	kind  string
	name  string
	stmt  string
	reads []string
}

// Dump writes the database's DDL as the server states it. The metadata is not
// rendered and does not narrow the dump: information_schema and system.tables
// do not carry a table's settings, data-skipping indexes, TTL or column codecs,
// so a synced schema cannot be turned back into ClickHouse DDL, and the
// server's own statements are the only complete definition. The caller creates
// the database; the dump carries the objects in it, in an order that replays.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	database, err := d.syncDatabaseName(ctx)
	if err != nil {
		return wrapDatabaseNotFound(err)
	}
	// The dump runs on one session so it can ask the server to state a secret in
	// full. Whether one is shown is the server's decision — it needs the
	// display_secrets_in_show_and_select setting and the
	// displaySecretsInShowAndSelect privilege — and the server states the
	// placeholder when it does not show it.
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return wrapDatabaseNotFound(err)
	}
	defer func() { _ = conn.Close() }()
	// A session that may not change a setting — readonly = 1 — answers with an
	// error and still states every object, with the placeholder for a secret.
	_, _ = conn.ExecContext(ctx, "SET format_display_secrets_in_show_and_select = 1")
	objects, err := d.dumpObjects(ctx, conn, database)
	if err != nil {
		return wrapDatabaseNotFound(err)
	}
	for _, o := range orderDumpObjects(objects) {
		if _, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", dumpKindLabel(o.kind), dumpIdentifier(o.name), o.stmt); err != nil {
			return err
		}
	}
	return nil
}

// dumpObjects reads every object of database and the statement the server
// states for it.
func (d *Driver) dumpObjects(ctx context.Context, conn *sql.Conn, database string) ([]*dumpObject, error) {
	// The server keeps an object's own tables in the ".inner_id." namespace: a
	// materialized view with its own engine has ".inner_id.<uuid>", and a
	// TimeSeries table has one per role, ".inner_id.<role>.<uuid>". They are the
	// object's storage, not objects of the database, and a dump that wrote them
	// would make a second one beside the object's own.
	const query = `
		SELECT
			name
		FROM system.tables
		WHERE database = ? AND NOT match(name, ?)
		ORDER BY name`
	rows, err := conn.QueryContext(ctx, query, database, dumpStorageTablePattern)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	var objects []*dumpObject
	for rows.Next() {
		var o dumpObject
		if err := rows.Scan(&o.name); err != nil {
			return nil, err
		}
		objects = append(objects, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	// The statement queries go in name order, which is the order the rows
	// arrived in, so a recording of them does not depend on a map.
	for _, o := range objects {
		if o.stmt, err = d.dumpStatement(ctx, conn, database, o.name); err != nil {
			return nil, err
		}
		o.kind = dumpKind(o.stmt)
	}
	byName := make(map[string]*dumpObject, len(objects))
	for _, o := range objects {
		byName[o.name] = o
	}
	definitions := make(map[string]*dumpObject, len(objects))
	for _, o := range objects {
		switch o.kind {
		case dumpView, dumpMaterializedView, dumpDictionaryTable:
			definitions[o.name] = o
		}
	}
	for _, o := range objects {
		if _, ok := definitions[o.name]; !ok {
			continue
		}
		if o.kind == dumpDictionaryTable {
			o.reads = dumpDictionaryTableReads(o, database, byName)
			if query, ok := dumpDictionarySourceQuery(o, database, byName); ok {
				reads, err := d.dumpReads(ctx, conn, database, o.name, query, definitions)
				if err != nil {
					return nil, err
				}
				o.reads = append(o.reads, reads...)
			}
			continue
		}
		reads, err := d.dumpReads(ctx, conn, database, o.name, o.stmt, definitions)
		if err != nil {
			return nil, err
		}
		o.reads = reads
	}
	for _, o := range objects {
		if o.kind == dumpMaterializedView || o.kind == dumpDictionaryTable {
			o.reads = dumpExpandedReads(o, definitions)
		}
	}
	return objects, nil
}

// dumpStorageTablePattern matches the tables the server makes for an object. The
// shape is matched as the server writes it rather than by a prefix a user's own
// name could take.
const dumpStorageTablePattern = `^[.]inner_id[.]([a-z_]+[.])?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

// dumpKind names the object a statement creates. SHOW CREATE TABLE states a
// dictionary as CREATE DICTIONARY, a view and a materialized view as their own
// statements, and everything else as CREATE TABLE — a table on the Dictionary
// engine included, which is a table that reads a dictionary and is stated as
// one.
func dumpKind(stmt string) string {
	switch {
	case strings.HasPrefix(stmt, "CREATE DICTIONARY "):
		return dumpDictionary
	case strings.HasPrefix(stmt, "CREATE MATERIALIZED VIEW "):
		return dumpMaterializedView
	case strings.HasPrefix(stmt, "CREATE TABLE ") && strings.Contains(stmt, "ENGINE = Dictionary("):
		return dumpDictionaryTable
	case strings.HasPrefix(stmt, "CREATE TABLE "):
		return dumpTable
	default:
		// A view states CREATE VIEW, and a release that adds another kind of
		// view states it as a CREATE of its own; either way it is a definition
		// rather than a table.
		return dumpView
	}
}

// dumpStatement returns the statement the server states for one object, without
// its terminating semicolon. SHOW CREATE TABLE answers one row in one column
// for every kind of object, a dictionary included.
func (d *Driver) dumpStatement(ctx context.Context, conn *sql.Conn, database, name string) (string, error) {
	query := "SHOW CREATE TABLE " + dumpIdentifier(database) + "." + dumpIdentifier(name)
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return "", formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", formatErrorWithQuery(err, query)
		}
		return "", fmt.Errorf("%s: no row", query)
	}
	var stmt string
	if err := rows.Scan(&stmt); err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSpace(stmt), ";"), nil
}

// orderDumpObjects puts the objects in the order a restore replays them:
//
//   - tables, whose creation resolves nothing;
//   - dictionaries, whose creation resolves nothing either, so they follow the
//     tables their sources usually are;
//   - views, whose creation does not resolve the query they hold, so their own
//     order does not matter;
//   - materialized views and tables on the Dictionary engine, which resolve what
//     they read when they are created: a materialized view its source and its TO
//     target, and a table on the Dictionary engine its dictionary, which loads
//     and reads the dictionary's source. Each follows what it reads.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var tables, dictionaries, dictionaryTables, views, materializedViews []*dumpObject
	for _, o := range objects {
		switch o.kind {
		case dumpTable:
			tables = append(tables, o)
		case dumpDictionary:
			dictionaries = append(dictionaries, o)
		case dumpDictionaryTable:
			dictionaryTables = append(dictionaryTables, o)
		case dumpMaterializedView:
			materializedViews = append(materializedViews, o)
		default:
			views = append(views, o)
		}
	}
	ordered := sortDumpObjects(tables)
	ordered = append(ordered, sortDumpObjects(dictionaries)...)
	ordered = append(ordered, sortDumpObjects(views)...)
	return append(ordered, orderDumpResolved(append(materializedViews, dictionaryTables...))...)
}

// orderDumpResolved orders the objects that resolve what they read when they are
// created, so that one follows every one it reads, directly or through a view,
// names breaking ties so the dump is deterministic. A reference cycle cannot
// exist on the server, which refuses the creation that would close it; were one
// found anyway, the remaining objects are written in name order so the dump
// still carries them.
func orderDumpResolved(objects []*dumpObject) []*dumpObject {
	pending := sortDumpObjects(objects)
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

// dumpReads returns the other definitions in definitions that o reads, as the
// server parses o's statement: EXPLAIN AST renders every table a statement names
// after the marker TableIdentifier, so a name inside a string literal, a column
// or an alias is not a read, and neither is a name that merely starts another.
func (d *Driver) dumpReads(ctx context.Context, conn *sql.Conn, database, self, statement string, definitions map[string]*dumpObject) ([]string, error) {
	query := "EXPLAIN AST " + statement
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	var ast strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, formatErrorWithQuery(err, query)
		}
		ast.WriteString(line)
		ast.WriteString("\n")
	}
	if err := rows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, query)
	}

	return dumpASTReads(ast.String(), database, self, definitions), nil
}

// dumpASTReads returns the definitions in definitions that the AST names as a
// table, and not the statement's own object.
func dumpASTReads(ast, database, self string, definitions map[string]*dumpObject) []string {
	var reads []string
	for name := range definitions {
		if name != self && dumpASTNames(ast, "TableIdentifier "+database+"."+name) {
			reads = append(reads, name)
		}
	}
	sort.Strings(reads)
	return reads
}

// dumpASTNames reports whether the AST names reference as a node of its own: the
// marker starts a line once its indentation is skipped — so a string literal that
// spells it out, which is escaped and mid-line, is not one — and the name ends
// the node, at a line break or before the alias the rendering appends.
func dumpASTNames(ast, reference string) bool {
	indentation := func(start int) bool {
		for i := start - 1; i >= 0 && ast[i] != '\n'; i-- {
			if ast[i] != ' ' && ast[i] != '\t' {
				return false
			}
		}
		return true
	}
	for offset := 0; offset < len(ast); {
		start := strings.Index(ast[offset:], reference)
		if start < 0 {
			return false
		}
		start += offset
		end := start + len(reference)
		if indentation(start) && (end == len(ast) || ast[end] == '\n' || strings.HasPrefix(ast[end:], " (alias ")) {
			return true
		}
		offset = end
	}
	return false
}

// The clauses that state what a table on the Dictionary engine reads: the ENGINE
// argument names the dictionary, and the dictionary's SOURCE names the table and
// database it reads, as string literals.
var (
	dumpDictionaryEngine   = regexp.MustCompile(`ENGINE = Dictionary\('((?:[^'\\]|\\.)*)'\)`)
	dumpDictionaryTableKey = regexp.MustCompile(`(?i)TABLE '((?:[^'\\]|\\.)*)'`)
	dumpDictionaryDBKey    = regexp.MustCompile(`(?i)DB '((?:[^'\\]|\\.)*)'`)
	dumpDictionaryQuery    = regexp.MustCompile(`(?i)QUERY '((?:[^'\\]|\\.)*)'`)
)

// dumpDictionaryTableReads returns the objects a table on the Dictionary engine
// reads: creating the table loads the dictionary, and the dictionary reads its
// source, which may be a table, a view or a materialized view of this database.
// Only those clauses are read; a comment or a default that happens to spell an
// object's name does not make one.
func dumpDictionaryTableReads(o *dumpObject, database string, byName map[string]*dumpObject) []string {
	engine := dumpDictionaryEngine.FindStringSubmatch(o.stmt)
	if engine == nil {
		return nil
	}
	dictionary, _ := dumpLiteralObject(dumpLiteralValue(engine[1]), database, byName)
	if dictionary == nil {
		return nil
	}
	table := dumpDictionaryTableKey.FindStringSubmatch(dictionary.stmt)
	if table == nil {
		return nil
	}
	if db := dumpDictionaryDBKey.FindStringSubmatch(dictionary.stmt); db != nil && !strings.EqualFold(dumpLiteralValue(db[1]), database) {
		// A table of another database is not an object of this one, whatever it
		// is called.
		return nil
	}
	if read, name := dumpLiteralObject(dumpLiteralValue(table[1]), database, byName); read != nil && name != o.name {
		return []string{name}
	}
	return nil
}

// dumpDictionarySourceQuery returns the query a dictionary states as its source
// instead of a table, which the server parses like any other statement.
func dumpDictionarySourceQuery(o *dumpObject, database string, byName map[string]*dumpObject) (string, bool) {
	engine := dumpDictionaryEngine.FindStringSubmatch(o.stmt)
	if engine == nil {
		return "", false
	}
	dictionary, _ := dumpLiteralObject(dumpLiteralValue(engine[1]), database, byName)
	if dictionary == nil {
		return "", false
	}
	query := dumpDictionaryQuery.FindStringSubmatch(dictionary.stmt)
	if query == nil {
		return "", false
	}
	return dumpLiteralValue(query[1]), true
}

// dumpLiteralValue returns what the text of a ClickHouse string literal stands
// for: a quote, a backslash and a control character are escaped with a backslash.
func dumpLiteralValue(text string) string {
	if !strings.Contains(text, `\`) {
		return text
	}
	var value strings.Builder
	for i := 0; i < len(text); i++ {
		if text[i] != '\\' || i+1 == len(text) {
			value.WriteByte(text[i])
			continue
		}
		i++
		switch c := text[i]; c {
		case 'n':
			value.WriteByte('\n')
		case 'r':
			value.WriteByte('\r')
		case 't':
			value.WriteByte('\t')
		case 'b':
			value.WriteByte('\b')
		case 'f':
			value.WriteByte('\f')
		case '0':
			value.WriteByte(0)
		default:
			value.WriteByte(c)
		}
	}
	return value.String()
}

// dumpLiteralObject returns the object of byName that literal names. The server
// writes a literal as it was given, so it may be qualified with its database and
// may hold a dot of its own: the literal is tried whole, then without the
// database, and only then by its last part.
func dumpLiteralObject(literal, database string, byName map[string]*dumpObject) (*dumpObject, string) {
	for _, name := range []string{
		literal,
		strings.TrimPrefix(literal, database+"."),
		literal[strings.LastIndex(literal, ".")+1:],
	} {
		if object := byName[name]; object != nil {
			return object, name
		}
	}
	return nil, ""
}

// dumpExpandedReads returns every definition o reads, directly or through a view:
// creating a view does not resolve the query it holds, so a view may be created
// before what it reads, but creating a materialized view resolves its own query
// through the views it names.
func dumpExpandedReads(o *dumpObject, definitions map[string]*dumpObject) []string {
	seen := make(map[string]bool, len(definitions))
	var walk func(reads []string)
	walk = func(reads []string) {
		for _, name := range reads {
			if seen[name] || name == o.name {
				continue
			}
			seen[name] = true
			if read := definitions[name]; read != nil {
				walk(read.reads)
			}
		}
	}
	walk(o.reads)
	expanded := make([]string, 0, len(seen))
	for name := range seen {
		expanded = append(expanded, name)
	}
	sort.Strings(expanded)
	return expanded
}

// dumpIdentifier writes name in backquotes, with the escapes ClickHouse writes
// for a backslash, a backquote and the control characters an identifier may
// hold. Quoting always is what makes the statement the dump asks for, the header
// above it — a -- comment a line break would end — and the reference another
// statement holds safe without a list of the keywords ClickHouse quotes anyway.
func dumpIdentifier(name string) string {
	var b strings.Builder
	b.WriteByte('`')
	for i := 0; i < len(name); i++ {
		switch c := name[i]; c {
		case '\\', '`':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('`')
	return b.String()
}

// dumpKindLabel names a kind of object for the header of its statement.
func dumpKindLabel(kind string) string {
	switch kind {
	case dumpDictionary:
		return "Dictionary"
	case dumpView:
		return "View"
	case dumpMaterializedView:
		return "Materialized view"
	default:
		return "Table"
	}
}
