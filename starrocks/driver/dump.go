package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/starrocks/ast"
	srparser "github.com/bytebase/omni/starrocks/parser"
)

// The object kinds a dump carries, in the order it writes them.
const (
	dumpTable            = "table"
	dumpView             = "view"
	dumpMaterializedView = "materialized view"
)

// A table's key model, DISTRIBUTED BY, PROPERTIES, partitions and indexes are not
// in information_schema, so the listing queries below find only the objects'
// names; the definitions come from each object's SHOW CREATE form.

const (
	dumpDatabaseQuery = `SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?`

	dumpTablesQuery = `SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME`

	// information_schema.tables lists an asynchronous materialized view as a
	// VIEW, indistinguishable from a view, so the materialized views come from
	// information_schema.materialized_views and are removed from the views. A
	// synchronous materialized view (a rollup index) is in that catalogue too,
	// with REFRESH_TYPE 'ROLLUP', but is not a table of any kind: see the driver
	// package's note on what the dump leaves out.
	dumpMaterializedViewsQuery = `SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME`
)

var _ driver.Dumper = (*Driver)(nil)

// dumpObject is one object a dump writes: the statement the server states, and
// the objects in the same database that statement reads, so the dump can order
// definitions whose creation depends on another definition.
type dumpObject struct {
	kind  string
	name  string
	stmt  string
	reads []string
}

// Dump writes the database's DDL as the server states it: tables, then views and
// asynchronous materialized views in an order that replays. The metadata is not
// rendered and does not narrow the dump, because a StarRocks table's definition
// is more than its columns and information_schema states none of the rest, so
// the server's own SHOW CREATE statements are the only complete definition. The
// caller creates the database; the dump carries the objects in it.
//
// Every statement the dump writes takes effect when the server answers it: a
// table, a view and an asynchronous materialized view are all created
// synchronously, so a script applies back-to-back. A rollup is created by an
// asynchronous schema change of its base table, and Dump carries none, so the
// pacing a Doris dump documents for two synchronous views on one table has no
// counterpart here. TestDumpGaps records that omission, and the one for a
// user-defined function, which no statement states; StarRocks has no stored
// procedure.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	if d.databaseName == "" {
		return errors.New("dump needs a database")
	}
	if err := d.dumpDatabaseExists(ctx); err != nil {
		return err
	}
	objects, err := d.dumpObjects(ctx)
	if err != nil {
		return err
	}
	for _, o := range orderDumpObjects(objects) {
		if err := writeDumpObject(out, o); err != nil {
			return err
		}
	}
	return nil
}

// dumpDatabaseExists reports a database the connection names but the server does
// not have as driver.ErrNotFound, so a typo does not produce an empty script.
func (d *Driver) dumpDatabaseExists(ctx context.Context) error {
	var name string
	err := d.db.QueryRowContext(ctx, dumpDatabaseQuery, d.databaseName).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: database %q", driver.ErrNotFound, d.databaseName)
	}
	return err
}

// dumpObjects reads every object's statement from the server.
func (d *Driver) dumpObjects(ctx context.Context) ([]*dumpObject, error) {
	materializedViews, err := d.dumpMaterializedViewNames(ctx)
	if err != nil {
		return nil, err
	}

	var objects []*dumpObject
	rows, err := d.db.QueryContext(ctx, dumpTablesQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dumpTablesQuery, err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, tableType string
		if err := rows.Scan(&name, &tableType); err != nil {
			return nil, err
		}
		var kind string
		switch {
		case materializedViews[name]:
			// Emitted below, from the materialized-view list.
			continue
		case tableType == baseTableType:
			kind = dumpTable
		case tableType == viewTableType:
			kind = dumpView
		default:
			// information_schema's own SYSTEM VIEWs and any kind a newer
			// version adds are not this database's schema.
			continue
		}
		objects = append(objects, &dumpObject{kind: kind, name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// The statement queries go in a fixed order, so a recording of them does
	// not depend on the map iteration above.
	for name := range materializedViews {
		objects = append(objects, &dumpObject{kind: dumpMaterializedView, name: name})
	}
	objects = sortDumpObjects(objects)
	for _, o := range objects {
		if o.stmt, err = d.dumpStatement(ctx, o.kind, o.name); err != nil {
			return nil, err
		}
		if o.kind == dumpView || o.kind == dumpMaterializedView {
			if o.reads, err = dumpDefinitionReads(o.stmt, d.databaseName); err != nil {
				return nil, fmt.Errorf("%s %q: %w", o.kind, o.name, err)
			}
		}
	}
	return objects, nil
}

// dumpMaterializedViewNames returns the set of the database's asynchronous
// materialized views, which information_schema.tables states as plain views.
func (d *Driver) dumpMaterializedViewNames(ctx context.Context) (map[string]bool, error) {
	rows, err := d.db.QueryContext(ctx, dumpMaterializedViewsQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dumpMaterializedViewsQuery, err)
	}
	defer rows.Close()
	names := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names[name] = true
	}
	return names, rows.Err()
}

// dumpStatement returns the statement the server states for one object, without
// its terminating semicolon. Each SHOW CREATE form returns the object's name
// first and its statement second, and a view adds session columns after them.
func (d *Driver) dumpStatement(ctx context.Context, kind, name string) (string, error) {
	qualified := quoteBacktick(d.databaseName) + "." + quoteBacktick(name)
	var query string
	switch kind {
	case dumpTable:
		query = "SHOW CREATE TABLE " + qualified
	case dumpView:
		query = "SHOW CREATE VIEW " + qualified
	case dumpMaterializedView:
		query = "SHOW CREATE MATERIALIZED VIEW " + qualified
	default:
		return "", fmt.Errorf("no definition statement for a %s", kind)
	}
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s: no row", query)
	}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(columns) < 2 {
		return "", fmt.Errorf("%s: %d columns, want the statement at 1", query, len(columns))
	}
	values := make([]sql.NullString, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	if err := rows.Scan(dest...); err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	return strings.TrimSuffix(strings.TrimSpace(values[1].String), ";"), nil
}

// dumpDefinitionReads returns the objects in database that a view's or a
// materialized view's definition reads. The definition is the server's own
// rendering, so every object reference is qualified with its database and the
// parser sees the statement the server stored.
func dumpDefinitionReads(stmt, database string) ([]string, error) {
	file, err := srparser.Parse(stmt)
	if err != nil {
		return nil, err
	}
	if len(file.Stmts) != 1 {
		return nil, fmt.Errorf("%d statements, want one", len(file.Stmts))
	}
	var query ast.Node
	switch s := file.Stmts[0].(type) {
	case *ast.CreateViewStmt:
		query = s.Query
	case *ast.CreateMTMVStmt:
		query = s.Query
	default:
		return nil, fmt.Errorf("%T is not a view or materialized view definition", file.Stmts[0])
	}
	var reads []string
	var collect func(node ast.Node)
	collect = func(node ast.Node) {
		ast.Inspect(node, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.TableRef:
				// A FROM-subquery packs its text into Name; Walk visits the
				// Subquery node instead, which the case below reads.
				if n.Subquery == nil {
					if name := dumpReferencedName(n.Name, database); name != "" {
						reads = append(reads, name)
					}
				}
			case *ast.SubqueryExpr:
				// Walk keeps a subquery body as a leaf and holds the parsed body
				// in Query, so a reference inside a scalar, EXISTS or IN
				// subquery is a reference all the same.
				collect(n.Query)
				return false
			}
			return true
		})
	}
	collect(query)
	return reads, nil
}

// dumpReferencedName returns the object a table reference names when it is in
// database, and "" when it is not: a reference to another database is not the
// dump's to create. SHOW CREATE states every object a view or materialized view
// reads with two parts, the definition having been read in the session's
// catalog, so an unqualified name is a CTE and taking one for an object would
// invent a dependency that can reverse a real one: a view named like a CTE in
// another definition's WITH clause would read that definition and be read by it.
func dumpReferencedName(name *ast.ObjectName, database string) string {
	if name == nil || len(name.Parts) != 2 {
		return ""
	}
	if !strings.EqualFold(name.Parts[0], database) {
		return ""
	}
	return name.Parts[1]
}

// orderDumpObjects puts the objects in the order a restore replays: tables
// first, then views and materialized views so that each definition reads only
// objects already created. A view may read a materialized view and a
// materialized view may read a view, so the two kinds share one dependency
// order, and either may read another of its kind.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var tables, definitions []*dumpObject
	for _, o := range objects {
		if o.kind == dumpTable {
			tables = append(tables, o)
			continue
		}
		definitions = append(definitions, o)
	}
	return append(sortDumpObjects(tables), orderDumpDefinitions(definitions)...)
}

// orderDumpDefinitions orders views and materialized views so that a definition
// follows every definition in its database that it reads, names breaking ties so
// the dump is deterministic. A reference cycle cannot exist on the server, which
// refuses the creation that would close it; were one extracted anyway, the
// remaining objects are written in name order so the dump still carries them.
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

// writeDumpObject writes one object's statement under a header naming it.
func writeDumpObject(out io.Writer, o *dumpObject) error {
	var label string
	switch o.kind {
	case dumpTable:
		label = "Table"
	case dumpView:
		label = "View"
	case dumpMaterializedView:
		label = "Materialized view"
	default:
		return fmt.Errorf("no header for a %s", o.kind)
	}
	_, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", label, quoteBacktick(o.name), o.stmt)
	return err
}
