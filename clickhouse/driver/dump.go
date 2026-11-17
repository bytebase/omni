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

// The object kinds a dump carries.
const (
	dumpTable            = "table"
	dumpDictionary       = "dictionary"
	dumpView             = "view"
	dumpMaterializedView = "materialized view"
)

// dumpObject is one object a dump writes: the kind the server reports for it,
// its name, the statement the server states, and the definitions in the same
// database that statement reads, which orders the definitions.
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
		if err := writeDumpObject(out, o); err != nil {
			return err
		}
	}
	return nil
}

// dumpObjects reads every object of database and the statement the server
// states for it.
func (d *Driver) dumpObjects(ctx context.Context, database string) ([]*dumpObject, error) {
	// A materialized view with its own engine has a table the server named
	// itself, ".inner_id.<uuid>"; it is the view's storage, not an object of
	// the database, and a dump that wrote it would create a second one beside
	// the view's.
	const query = `
		SELECT
			name,
			engine
		FROM system.tables
		WHERE database = ? AND name NOT LIKE '.%'
		ORDER BY name`
	rows, err := d.db.QueryContext(ctx, query, database)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	var objects []*dumpObject
	for rows.Next() {
		var name, engine string
		if err := rows.Scan(&name, &engine); err != nil {
			return nil, err
		}
		objects = append(objects, &dumpObject{kind: dumpKind(engine), name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	definitions := make(map[string]bool, len(objects))
	for _, o := range objects {
		if o.kind == dumpView || o.kind == dumpMaterializedView {
			definitions[o.name] = true
		}
	}
	// The statement queries go in name order, which is the order the rows
	// arrived in, so a recording of them does not depend on a map.
	for _, o := range objects {
		if o.stmt, err = d.dumpStatement(ctx, database, o); err != nil {
			return nil, err
		}
		if definitions[o.name] {
			o.reads = dumpReads(o, database, definitions)
		}
	}
	return objects, nil
}

// dumpKind names the kind of object a table's engine makes it.
func dumpKind(engine string) string {
	switch engine {
	case "View":
		return dumpView
	case "MaterializedView":
		return dumpMaterializedView
	case "Dictionary":
		return dumpDictionary
	default:
		return dumpTable
	}
}

// dumpStatement returns the statement the server states for one object, without
// its terminating semicolon. Both SHOW CREATE forms answer one row in one
// column.
func (d *Driver) dumpStatement(ctx context.Context, database string, o *dumpObject) (string, error) {
	qualified := dumpIdentifier(database) + "." + dumpIdentifier(o.name)
	query := "SHOW CREATE TABLE " + qualified
	if o.kind == dumpDictionary {
		query = "SHOW CREATE DICTIONARY " + qualified
	}
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

// orderDumpObjects puts the objects in the order a restore replays: tables,
// then dictionaries, whose creation resolves no source, then views and
// materialized views, each after the definitions in its database it reads.
// Only a materialized view has to be ordered: it resolves its source and its TO
// target when it is created, while a view states its columns, and ClickHouse
// creates a view that states them without resolving what it reads. Neither the
// server nor the metadata states what a materialized view reads, so the order
// comes from the statements.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var tables, dictionaries, definitions []*dumpObject
	for _, o := range objects {
		switch o.kind {
		case dumpTable:
			tables = append(tables, o)
		case dumpDictionary:
			dictionaries = append(dictionaries, o)
		default:
			definitions = append(definitions, o)
		}
	}
	ordered := append(sortDumpObjects(tables), sortDumpObjects(dictionaries)...)
	return append(ordered, orderDumpDefinitions(definitions)...)
}

// orderDumpDefinitions orders views and materialized views so that a definition
// follows every definition in its database that it reads, names breaking ties
// so the dump is deterministic. A reference cycle cannot exist on the server,
// which refuses the creation that would close it; were one found anyway, the
// remaining definitions are written in name order so the dump still carries
// them.
func orderDumpDefinitions(definitions []*dumpObject) []*dumpObject {
	pending := sortDumpObjects(definitions)
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

// dumpReads returns the other definitions in names that o reads. ClickHouse's
// rendering names every table a materialized view reads with its database, so a
// definition is found by the qualified spelling the dump would use for it, and
// o's own name in its CREATE statement is not a read. A reference inside a
// string literal adds an edge that is not there, which only moves a definition
// later.
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

// dumpReferences reports whether stmt names database.name, and not the tail of
// a longer identifier.
func dumpReferences(stmt, database, name string) bool {
	reference := dumpIdentifier(database) + "." + dumpIdentifier(name)
	for offset := 0; offset < len(stmt); {
		start := strings.Index(stmt[offset:], reference)
		if start < 0 {
			return false
		}
		start += offset
		end := start + len(reference)
		if !dumpReferenceRune(stmt, start-1) && !dumpReferenceRune(stmt, end) {
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

// writeDumpObject writes one object's statement under a header naming it.
func writeDumpObject(out io.Writer, o *dumpObject) error {
	label := "Table"
	switch o.kind {
	case dumpDictionary:
		label = "Dictionary"
	case dumpView:
		label = "View"
	case dumpMaterializedView:
		label = "Materialized view"
	}
	_, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", label, dumpIdentifier(o.name), o.stmt)
	return err
}
