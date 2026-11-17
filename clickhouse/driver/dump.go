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
	if _, err := conn.ExecContext(ctx, "SET format_display_secrets_in_show_and_select = 1"); err != nil {
		return wrapDatabaseNotFound(formatErrorWithQuery(err, "SET format_display_secrets_in_show_and_select = 1"))
	}
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
	// A materialized view with its own engine has a table the server named
	// itself, ".inner_id.<uuid>". It is the view's storage, not an object of the
	// database, and a dump that wrote it would create a second one beside the
	// view's. The shape is matched as the server writes it rather than by a
	// prefix a user's own name could take.
	const query = `
		SELECT
			name
		FROM system.tables
		WHERE database = ? AND NOT match(name, '^[.]inner_id[.][0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$')
		ORDER BY name`
	rows, err := conn.QueryContext(ctx, query, database)
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
			o.reads = dumpDictionaryTableReads(o, byName)
			continue
		}
		o.reads = dumpReads(o, database, definitions)
	}
	for _, o := range objects {
		if o.kind == dumpMaterializedView || o.kind == dumpDictionaryTable {
			o.reads = dumpExpandedReads(o, definitions)
		}
	}
	return objects, nil
}

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

// dumpReads returns the other definitions in definitions that o reads.
// ClickHouse's rendering names every table a statement reads with its database,
// so one is found by the qualified spelling the dump would use for it, and o's
// own name in its CREATE statement is not a read.
func dumpReads(o *dumpObject, database string, definitions map[string]*dumpObject) []string {
	var reads []string
	for name := range definitions {
		if name != o.name && dumpReferences(o.stmt, database, name) {
			reads = append(reads, name)
		}
	}
	sort.Strings(reads)
	return reads
}

// dumpDictionaryEngine matches the dictionary a table on the Dictionary engine
// reads, which the server states as ENGINE = Dictionary('<database>.<name>').
var dumpDictionaryEngine = regexp.MustCompile(`ENGINE = Dictionary\('([^']+)'\)`)

// dumpDictionaryTableReads returns the objects a table on the Dictionary engine
// reads: creating the table loads the dictionary, and the dictionary reads its
// source, which may be a table, a view or a materialized view of this database.
// Both are stated as string literals — the table's ENGINE clause and the
// dictionary's SOURCE — so the names come from those literals.
func dumpDictionaryTableReads(o *dumpObject, byName map[string]*dumpObject) []string {
	engine := dumpDictionaryEngine.FindStringSubmatch(o.stmt)
	if engine == nil {
		return nil
	}
	dictionary := byName[dumpQualifiedName(engine[1])]
	if dictionary == nil {
		return nil
	}
	seen := make(map[string]bool, len(byName))
	for _, literal := range dumpLiterals(dictionary.stmt) {
		name := dumpQualifiedName(literal)
		if name != o.name && byName[name] != nil {
			seen[name] = true
		}
	}
	reads := make([]string, 0, len(seen))
	for name := range seen {
		reads = append(reads, name)
	}
	sort.Strings(reads)
	return reads
}

// dumpQualifiedName returns the object part of a name the server states, which
// may be qualified with its database.
func dumpQualifiedName(qualified string) string {
	return qualified[strings.LastIndex(qualified, ".")+1:]
}

// dumpLiterals returns the contents of the single-quoted strings in stmt, with a
// backslash escape resolved, which is how a dictionary states its source and a
// table states the dictionary it reads.
func dumpLiterals(stmt string) []string {
	var literals []string
	var text strings.Builder
	inString := false
	for i := 0; i < len(stmt); i++ {
		switch c := stmt[i]; {
		case !inString && c == '\'':
			inString, text = true, strings.Builder{}
		case inString && c == '\\' && i+1 < len(stmt):
			i++
			text.WriteByte(stmt[i])
		case inString && c == '\'':
			inString = false
			literals = append(literals, text.String())
		case inString:
			text.WriteByte(c)
		}
	}
	return literals
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

// dumpReferences reports whether stmt names database.name outside a string
// literal, and not as the tail of a longer identifier. A reference inside a
// literal is not a read: an edge it added could close a cycle and make the order
// fall back to names.
func dumpReferences(stmt, database, name string) bool {
	// ClickHouse writes each part of a reference bare where it can and in
	// backquotes where it must — and in backquotes for the keywords it quotes
	// anyway — so every spelling is looked for.
	for _, db := range []string{database, dumpIdentifier(database)} {
		for _, object := range []string{name, dumpIdentifier(name)} {
			if dumpNamesUnquoted(stmt, db+"."+object) {
				return true
			}
		}
	}
	return false
}

// dumpNamesUnquoted reports whether stmt holds reference outside a quote, and not
// as the tail of a longer identifier.
func dumpNamesUnquoted(stmt, reference string) bool {
	for offset := 0; offset < len(stmt); {
		start := strings.Index(stmt[offset:], reference)
		if start < 0 {
			return false
		}
		start += offset
		end := start + len(reference)
		if !dumpReferenceRune(stmt, start-1) && !dumpReferenceRune(stmt, end) && dumpQuoted(stmt, start) == 0 {
			return true
		}
		offset = end
	}
	return false
}

// dumpReferenceRune reports whether the byte at i continues an identifier or a
// qualified name, so a match beside it is part of a longer one.
func dumpReferenceRune(s string, i int) bool {
	if i < 0 || i >= len(s) {
		return false
	}
	c := s[i]
	return c == '_' || c == '$' || c == '.' || c == '`' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// dumpQuoted returns the quote the byte at want is inside, and 0 when it is not
// quoted: ClickHouse's rendering escapes a quote inside a string or a name with a
// backslash, writes a string in single quotes, and writes a name in backquotes or
// double quotes. A reference inside either is not one: inside a string it is
// text, and inside a name it is part of that name.
func dumpQuoted(stmt string, want int) byte {
	quote := byte(0)
	for i := 0; i < want && i < len(stmt); i++ {
		switch c := stmt[i]; {
		case quote == 0 && (c == '\'' || c == '`' || c == '"'):
			quote = c
		case quote != 0 && c == '\\':
			i++
		case quote != 0 && c == quote:
			quote = 0
		}
	}
	return quote
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
