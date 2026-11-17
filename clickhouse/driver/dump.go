package driver

import (
	"context"
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
	objects, err := d.dumpObjects(ctx, database)
	if err != nil {
		return wrapDatabaseNotFound(err)
	}
	for _, o := range orderDumpObjects(objects) {
		if _, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", dumpKindLabel(o.kind), dumpLabel(o.name), o.stmt); err != nil {
			return err
		}
	}
	return nil
}

// dumpObjects reads every object of database and the statement the server
// states for it.
func (d *Driver) dumpObjects(ctx context.Context, database string) ([]*dumpObject, error) {
	// A materialized view with its own engine has a table the server named
	// itself, ".inner_id.<uuid>" (".inner.<name>" before 21.4). It is the
	// view's storage, not an object of the database, and a dump that wrote it
	// would create a second one beside the view's.
	const query = `
		SELECT
			name
		FROM system.tables
		WHERE database = ? AND name NOT LIKE '.inner_id.%' AND name NOT LIKE '.inner.%'
		ORDER BY name`
	rows, err := d.db.QueryContext(ctx, query, database)
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
		if o.stmt, err = d.dumpStatement(ctx, database, o.name); err != nil {
			return nil, err
		}
		o.kind = dumpKind(o.stmt)
	}
	materializedViews := make(map[string]bool, len(objects))
	for _, o := range objects {
		if o.kind == dumpMaterializedView {
			materializedViews[o.name] = true
		}
	}
	for _, o := range objects {
		if o.kind == dumpMaterializedView {
			o.reads = dumpReads(o, database, materializedViews)
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
func (d *Driver) dumpStatement(ctx context.Context, database, name string) (string, error) {
	query := "SHOW CREATE TABLE " + dumpIdentifier(database) + "." + dumpIdentifier(name)
	rows, err := d.db.QueryContext(ctx, query)
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
//   - tables on the Dictionary engine, which load their dictionary when they
//     are created, and with it the dictionary's source;
//   - views, whose creation does not resolve the query they hold, so their own
//     order does not matter;
//   - materialized views, which do resolve their source and their TO target
//     when they are created, each after the materialized views it reads.
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
	ordered = append(ordered, sortDumpObjects(dictionaryTables)...)
	ordered = append(ordered, sortDumpObjects(views)...)
	return append(ordered, orderDumpMaterializedViews(materializedViews)...)
}

// orderDumpMaterializedViews orders materialized views so that one follows every
// materialized view it reads, names breaking ties so the dump is deterministic.
// A reference cycle cannot exist on the server, which refuses the creation that
// would close it; were one found anyway, the remaining views are written in name
// order so the dump still carries them.
func orderDumpMaterializedViews(materializedViews []*dumpObject) []*dumpObject {
	pending := sortDumpObjects(materializedViews)
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

// dumpReads returns the other materialized views in names that o reads.
// ClickHouse's rendering names every table a materialized view reads with its
// database, so one is found by the qualified spelling the dump would use for it,
// and o's own name in its CREATE statement is not a read.
func dumpReads(o *dumpObject, database string, names map[string]bool) []string {
	var reads []string
	for name := range names {
		if name != o.name && dumpReferences(o.stmt, database, name) {
			reads = append(reads, name)
		}
	}
	sort.Strings(reads)
	return reads
}

// dumpReferences reports whether stmt names database.name outside a string
// literal, and not as the tail of a longer identifier. A reference inside a
// literal is not a read: an edge it added could close a cycle and make the order
// fall back to names.
func dumpReferences(stmt, database, name string) bool {
	reference := dumpIdentifier(database) + "." + dumpIdentifier(name)
	for offset := 0; offset < len(stmt); {
		start := strings.Index(stmt[offset:], reference)
		if start < 0 {
			return false
		}
		start += offset
		end := start + len(reference)
		if !dumpReferenceRune(stmt, start-1) && !dumpReferenceRune(stmt, end) && !dumpInString(stmt, start) {
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

// dumpInString reports whether the byte at want is inside a string literal.
// ClickHouse's rendering escapes a quote inside a string with a backslash, and
// writes a name in backquotes or double quotes, where a quote needs none.
func dumpInString(stmt string, want int) bool {
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
	return quote == '\''
}

var plainIdentifier = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

// dumpIdentifier writes name the way ClickHouse's own rendering writes it: bare
// when it is a plain identifier, and in backquotes with backslashes and
// backquotes escaped otherwise.
func dumpIdentifier(name string) string {
	if plainIdentifier.MatchString(name) {
		return name
	}
	return "`" + strings.NewReplacer(`\`, `\\`, "`", "\\`").Replace(name) + "`"
}

// dumpLabel writes an object's name for the header of its statement. An
// identifier may hold a line break, which a -- comment would not survive, so it
// is escaped here.
func dumpLabel(name string) string {
	return strings.NewReplacer("\n", `\n`, "\r", `\r`, "\t", `\t`).Replace(dumpIdentifier(name))
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
