package driver

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var _ driver.Dumper = (*Driver)(nil)

// dumpViewPrefix begins every statement SHOW CREATE TABLE states for a view.
const dumpViewPrefix = "CREATE VIEW "

// dumpObject is one object a dump writes: the name the server lists, the
// statement the server states for it, whether that statement is a view's, and
// the views in the same database the statement reads, which orders the views.
type dumpObject struct {
	name   string
	stmt   string
	isView bool
	reads  []string
}

// Dump writes the database's DDL as the server states it. The metadata is not
// rendered and does not narrow the dump: DESCRIBE FORMATTED reports a table's
// type, columns, row count and partitions but not the row format, the input and
// output formats, the serde properties, the storage handler, the skew or bucket
// definition, or the table properties, so a synced schema cannot be turned back
// into HiveQL. SHOW CREATE TABLE states all of them, and it is what the caller
// gets. The caller creates the database; the dump carries the objects in it, in
// an order that replays.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	objects, err := d.dumpObjects(ctx, d.databaseNameOrDefault())
	if err != nil {
		return err
	}
	for _, o := range orderDumpObjects(objects) {
		label := "Table"
		if o.isView {
			label = "View"
		}
		if _, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", label, quoteBacktick(o.name), o.stmt); err != nil {
			return err
		}
	}
	return nil
}

// dumpObjects reads the database's tables and views and the statement the
// server states for each.
func (d *Driver) dumpObjects(ctx context.Context, database string) ([]*dumpObject, error) {
	names, err := d.dumpNames(ctx, database)
	if err != nil {
		return nil, err
	}
	var objects []*dumpObject
	for _, name := range names {
		stmt, err := d.dumpStatement(ctx, database, name)
		if err != nil {
			return nil, err
		}
		objects = append(objects, &dumpObject{
			name:   name,
			stmt:   stmt,
			isView: strings.HasPrefix(stmt, dumpViewPrefix),
		})
	}
	views := make(map[string]bool, len(objects))
	for _, o := range objects {
		if o.isView {
			views[o.name] = true
		}
	}
	for _, o := range objects {
		if o.isView {
			o.reads = dumpReads(o, database, views)
		}
	}
	return objects, nil
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

// orderDumpObjects puts the objects in the order a restore replays them: every
// object that is not a view first, since a table states no dependency and a
// materialized view is stated as one, then the views, each after the views in
// its database it reads. Hive resolves a view's query when it creates the view
// and states no dependency anywhere, so the order comes from the statements.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var other, views []*dumpObject
	for _, o := range objects {
		if o.isView {
			views = append(views, o)
		} else {
			other = append(other, o)
		}
	}
	return append(sortDumpObjects(other), orderDumpViews(views)...)
}

// orderDumpViews orders views so that a view follows every view in its database
// it reads, names breaking ties so the dump is deterministic. A reference cycle
// cannot exist on the server, which refuses the creation that would close it;
// were one found anyway, the remaining views are written in name order so the
// dump still carries them.
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

// dumpReads returns the other views in views that o reads. Hive's rendering
// writes every reference in a view's statement in backquotes and qualified with
// its database, which is what the dump would write for it, so a reference is a
// substring bounded by backquotes, and o's own name in its CREATE statement is
// not a read.
func dumpReads(o *dumpObject, database string, views map[string]bool) []string {
	var reads []string
	for name := range views {
		if name != o.name && dumpReferences(o.stmt, database, name) {
			reads = append(reads, name)
		}
	}
	sort.Strings(reads)
	return reads
}

// dumpReferences reports whether stmt names database.name outside a string
// literal. A name inside a literal is not a read, and the edge it would add can
// close a cycle and leave the order to names.
func dumpReferences(stmt, database, name string) bool {
	reference := quoteBacktick(database) + "." + quoteBacktick(name)
	for offset := 0; offset < len(stmt); {
		start := strings.Index(stmt[offset:], reference)
		if start < 0 {
			return false
		}
		start += offset
		if !dumpInString(stmt, start) {
			return true
		}
		offset = start + len(reference)
	}
	return false
}

// dumpInString reports whether the byte at want is inside a string literal.
// Hive's rendering quotes a string with ' or ", escapes with a backslash, and
// writes a name in backquotes, where neither quote needs escaping.
func dumpInString(stmt string, want int) bool {
	quote := byte(0)
	for i := 0; i < want && i < len(stmt); i++ {
		switch c := stmt[i]; {
		case quote == 0 && (c == '\'' || c == '"'):
			quote = c
		case quote != 0 && c == '\\':
			i++
		case quote != 0 && c == quote:
			quote = 0
		}
	}
	return quote != 0
}
