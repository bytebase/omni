package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bytebase/omni/doris/ast"
	dorisparser "github.com/bytebase/omni/doris/parser"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// The object kinds a dump carries, in the order it writes them.
const (
	dumpTable                = "table"
	dumpSyncMaterializedView = "sync materialized view"
	dumpFunction             = "function"
	dumpView                 = "view"
	dumpMaterializedView     = "materialized view"
	dumpProcedure            = "procedure"
)

// A table's key model (AGGREGATE, UNIQUE or DUPLICATE), DISTRIBUTED BY with its
// buckets, PROPERTIES, partitions, rollups and inline indexes are not in
// information_schema, so the listing queries below find only the objects' names;
// the definitions come from each object's SHOW CREATE form.

const (
	dumpDatabaseQuery = `SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?`

	// A view or materialized view states every reference with three parts, and
	// only the ones under the session's own catalog are this database's objects.
	dumpCatalogQuery = `SELECT current_catalog()`

	dumpTablesQuery = `SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME`

	// information_schema.tables lists an asynchronous materialized view as a
	// BASE TABLE and reports no engine for it, so the materialized views come
	// from Doris's own catalogue function and are removed from the tables. A
	// server that does not answer mv_infos() fails the dump rather than
	// dropping its materialized views from the script.
	dumpMaterializedViewsQuery = `SELECT Name FROM mv_infos("database"="%s") ORDER BY Name`

	// A procedure's body is not in information_schema.routines
	// (ROUTINE_DEFINITION is empty); its name is, in upper case, as the
	// catalogue writes it. The statement from SHOW CREATE PROCEDURE carries the
	// spelling the procedure was created with.
	dumpProceduresQuery = `SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME`

	// A database-scoped function is in no information_schema table either, and
	// SHOW CREATE FUNCTION needs the signature SHOW FULL FUNCTIONS states:
	// SHOW FUNCTIONS lists a name only.
	dumpFunctionsQuery = `SHOW FULL FUNCTIONS FROM %s`

	// A synchronous materialized view is an index on its base table, not a table
	// of its own, so neither information_schema nor mv_infos() lists it: the
	// table's own description names the indexes it carries.
	dumpTableIndexesQuery = `DESC %s.%s ALL`
)

var _ driver.Dumper = (*Driver)(nil)

// dumpObject is one object a dump writes: the statement the server states, and
// the objects in the same database that statement reads, so the dump can order
// definitions whose creation depends on another definition. A synchronous
// materialized view is scoped to its base table, which its kind states alone.
type dumpObject struct {
	kind  string
	name  string
	table string
	stmt  string
	reads []string
}

// Dump writes the database's DDL as the server states it: tables, their
// synchronous materialized views, functions, views, materialized views and
// procedures, in an order that replays. The
// metadata is not rendered and does not narrow the dump, because a Doris table's
// definition is more than its columns and information_schema states none of the
// rest, so the server's own SHOW CREATE statements are the only complete
// definition. The caller creates the database; the dump carries the objects in
// it.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	if d.databaseName == "" {
		return errors.New("dump needs a database")
	}
	if err := d.dumpDatabaseExists(ctx); err != nil {
		return err
	}
	catalog, err := d.dumpCatalog(ctx)
	if err != nil {
		return err
	}
	objects, err := d.dumpObjects(ctx, catalog)
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

// dumpCatalog returns the session's catalog, which a definition's three-part
// references name first.
func (d *Driver) dumpCatalog(ctx context.Context) (string, error) {
	var catalog string
	if err := d.db.QueryRowContext(ctx, dumpCatalogQuery).Scan(&catalog); err != nil {
		return "", fmt.Errorf("%s: %w", dumpCatalogQuery, err)
	}
	return catalog, nil
}

// dumpObjects reads every object's statement from the server.
func (d *Driver) dumpObjects(ctx context.Context, catalog string) ([]*dumpObject, error) {
	materializedViews, err := d.dumpMaterializedViewNames(ctx)
	if err != nil {
		return nil, err
	}

	var tables []string
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
		if kind == dumpTable {
			tables = append(tables, name)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	syncMaterializedViews, err := d.dumpSyncMaterializedViews(ctx, tables)
	if err != nil {
		return nil, err
	}
	objects = append(objects, syncMaterializedViews...)

	for name := range materializedViews {
		objects = append(objects, &dumpObject{kind: dumpMaterializedView, name: name})
	}

	functions, err := d.dumpFunctionSignatures(ctx)
	if err != nil {
		return nil, err
	}
	for _, signature := range functions {
		// A variadic function's definition does not replay: the catalogue states
		// the ellipsis against the last type (fv(int...)), Doris's parser refuses
		// that argument list, and the definition it states for the accepted form
		// (CREATE ALIAS FUNCTION fv(int...) ...) is refused when applied.
		// TestDumpGaps records the gap.
		if strings.Contains(signature, "...") {
			continue
		}
		objects = append(objects, &dumpObject{kind: dumpFunction, name: signature})
	}

	procedures, err := d.dumpProcedureNames(ctx)
	if err != nil {
		return nil, err
	}
	for _, name := range procedures {
		objects = append(objects, &dumpObject{kind: dumpProcedure, name: name})
	}

	// The statement queries go in a fixed order, so a recording of them does
	// not depend on the map iteration above.
	objects = sortDumpObjects(objects)
	defined := objects[:0]
	for _, o := range objects {
		stmt, err := d.dumpStatement(ctx, o.kind, o.name, o.table)
		if errors.Is(err, errNoStatement) && o.kind == dumpSyncMaterializedView {
			// A rollup index is listed like a synchronous materialized view, and
			// no statement states its definition; TestDumpGaps records the gap.
			continue
		}
		if err != nil {
			return nil, err
		}
		o.stmt = stmt
		if o.kind == dumpView || o.kind == dumpMaterializedView {
			if o.reads, err = dumpDefinitionReads(stmt, catalog, d.databaseName); err != nil {
				return nil, fmt.Errorf("%s %q: %w", o.kind, o.name, err)
			}
		}
		defined = append(defined, o)
	}
	return defined, nil
}

// errNoStatement reports a form the server answers no row for, which a
// synchronous materialized view's index can be.
var errNoStatement = errors.New("no statement")

// dumpMaterializedViewNames returns the set of the database's materialized
// views.
func (d *Driver) dumpMaterializedViewNames(ctx context.Context) (map[string]bool, error) {
	query := fmt.Sprintf(dumpMaterializedViewsQuery, dumpInfoString(d.databaseName))
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
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

// dumpSyncMaterializedViews returns the synchronous materialized views of the
// given tables, each scoped to its base table. The table is the one that reads
// them back: SHOW CREATE MATERIALIZED VIEW takes the view's name and its table.
func (d *Driver) dumpSyncMaterializedViews(ctx context.Context, tables []string) ([]*dumpObject, error) {
	var objects []*dumpObject
	for _, table := range tables {
		query := fmt.Sprintf(dumpTableIndexesQuery, quoteBacktick(d.databaseName), quoteBacktick(table))
		names, err := dumpIndexNames(ctx, d.db, query)
		if err != nil {
			return nil, err
		}
		for _, name := range names {
			// The table's own index is the base index; every other index is a
			// synchronous materialized view.
			if name == "" || strings.EqualFold(name, table) {
				continue
			}
			objects = append(objects, &dumpObject{kind: dumpSyncMaterializedView, name: name, table: table})
		}
	}
	return objects, nil
}

// dumpIndexNames returns the first column of a description of a table's
// indexes: the index name, repeated for each column it holds.
func dumpIndexNames(ctx context.Context, q interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, query string) ([]string, error) {
	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		name := values[0].String
		if len(names) > 0 && names[len(names)-1] == name {
			continue
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// dumpFunctionSignatures returns the signatures of the database's functions,
// in the order the server lists them. SHOW CREATE FUNCTION states a definition
// only for the session's database, which the connection selects.
func (d *Driver) dumpFunctionSignatures(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf(dumpFunctionsQuery, quoteBacktick(d.databaseName))
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var signatures []string
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		if values[0].Valid && values[0].String != "" {
			signatures = append(signatures, values[0].String)
		}
	}
	return signatures, rows.Err()
}

// dumpProcedureNames returns the database's procedure names, sorted.
func (d *Driver) dumpProcedureNames(ctx context.Context) ([]string, error) {
	rows, err := d.db.QueryContext(ctx, dumpProceduresQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dumpProceduresQuery, err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// dumpStatement returns the statement the server states for one object, without
// its terminating semicolon. Each SHOW CREATE form returns the object's name
// first and its statement second, followed by session columns a view adds.
func (d *Driver) dumpStatement(ctx context.Context, kind, name, table string) (string, error) {
	qualified := quoteBacktick(d.databaseName) + "." + quoteBacktick(name)
	var query string
	statementColumn := 1
	switch kind {
	case dumpTable:
		query = "SHOW CREATE TABLE " + qualified
	case dumpSyncMaterializedView:
		// The view is scoped to its base table, which the statement names.
		query = "SHOW CREATE MATERIALIZED VIEW " + quoteBacktick(name) + " ON " + quoteBacktick(d.databaseName) + "." + quoteBacktick(table)
		// TableName, ViewName, CreateStmt.
		statementColumn = 2
	case dumpFunction:
		// A signature is the server's own rendering; a variadic one is left out
		// of the object list, so it never reaches this query.
		query = "SHOW CREATE FUNCTION " + name
	case dumpView:
		query = "SHOW CREATE VIEW " + qualified
	case dumpMaterializedView:
		query = "SHOW CREATE MATERIALIZED VIEW " + qualified
	case dumpProcedure:
		query = "SHOW CREATE PROCEDURE " + qualified
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
		return "", fmt.Errorf("%s: %w", query, errNoStatement)
	}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(columns) <= statementColumn {
		return "", fmt.Errorf("%s: %d columns, want the statement at %d", query, len(columns), statementColumn)
	}
	values := make([]sql.NullString, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	if err := rows.Scan(dest...); err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	return strings.TrimSuffix(strings.TrimSpace(values[statementColumn].String), ";"), nil
}

// dumpDefinitionReads returns the objects in database that a view's or a
// materialized view's definition reads. The definition is the server's own
// rendering, so every reference is qualified with its catalog and database and
// the parser sees the statement the server stored.
func dumpDefinitionReads(stmt, catalog, database string) ([]string, error) {
	file, err := dorisparser.Parse(stmt)
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
				// A table-valued function such as mv_infos() is not an object
				// the dump creates.
				if n.Func == nil {
					if name := dumpReferencedName(n.Name, catalog, database); name != "" {
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
// database under catalog, and "" when it is not: a reference to another database
// or another catalog is not the dump's to create. SHOW CREATE qualifies every object a view or materialized
// view reads, so an unqualified name is a CTE or a table-valued function, and
// taking one for an object would invent a dependency that can reverse a real
// one: a view named like a CTE in another definition's WITH clause would read
// that definition and be read by it.
func dumpReferencedName(name *ast.ObjectName, catalog, database string) string {
	switch len(name.Parts) {
	case 2:
		if strings.EqualFold(name.Parts[0], database) {
			return name.Parts[1]
		}
	case 3:
		// catalog.database.table, and only under the session's own catalog: a
		// reference into an external catalog whose database has the same name is
		// not this database's object.
		if strings.EqualFold(name.Parts[0], catalog) && strings.EqualFold(name.Parts[1], database) {
			return name.Parts[2]
		}
	}
	return ""
}

// orderDumpObjects puts the objects in the order a restore replays: tables
// first, then functions, then views and materialized views so that each
// definition reads only objects already created, then procedures. A materialized
// view's definition cannot contain a view, but a view's may contain a
// materialized view, and either may call a function or read another of its kind,
// so definitions follow those references.
func orderDumpObjects(objects []*dumpObject) []*dumpObject {
	var tables, syncMaterializedViews, functions, definitions, procedures []*dumpObject
	for _, o := range objects {
		switch o.kind {
		case dumpTable:
			tables = append(tables, o)
		case dumpSyncMaterializedView:
			syncMaterializedViews = append(syncMaterializedViews, o)
		case dumpFunction:
			functions = append(functions, o)
		case dumpProcedure:
			procedures = append(procedures, o)
		default:
			definitions = append(definitions, o)
		}
	}
	ordered := append(sortDumpObjects(tables), sortSyncMaterializedViews(syncMaterializedViews)...)
	ordered = append(ordered, sortDumpObjects(functions)...)
	ordered = append(ordered, orderDumpDefinitions(definitions)...)
	return append(ordered, sortDumpObjects(procedures)...)
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

// sortSyncMaterializedViews orders synchronous materialized views by their base
// table, then their name: the same view name may repeat across tables.
func sortSyncMaterializedViews(objects []*dumpObject) []*dumpObject {
	sorted := append([]*dumpObject(nil), objects...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].table != sorted[j].table {
			return sorted[i].table < sorted[j].table
		}
		return sorted[i].name < sorted[j].name
	})
	return sorted
}

func sortDumpObjects(objects []*dumpObject) []*dumpObject {
	sorted := append([]*dumpObject(nil), objects...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].name < sorted[j].name })
	return sorted
}

// writeDumpObject writes one object's statement under a header naming it. A
// procedure's body holds semicolons, so it is written between DELIMITER
// directives, as a client applies it.
func writeDumpObject(out io.Writer, o *dumpObject) error {
	if o.kind == dumpProcedure {
		_, err := fmt.Fprintf(out, "--\n-- Procedure structure for %s\n--\nDELIMITER ;;\n%s;;\nDELIMITER ;\n\n", quoteBacktick(o.name), o.stmt)
		return err
	}
	switch o.kind {
	case dumpSyncMaterializedView:
		// No header: Doris keeps a comment that precedes CREATE MATERIALIZED
		// VIEW as part of the view's definition, so a header would come back
		// inside the statement of the next dump.
		_, err := fmt.Fprintf(out, "%s;\n\n", o.stmt)
		return err
	case dumpFunction:
		_, err := fmt.Fprintf(out, "--\n-- Function structure for %s\n--\n%s;\n\n", o.name, o.stmt)
		return err
	}
	var label string
	switch o.kind {
	case dumpTable:
		label = "Table"
	case dumpView:
		label = "View"
	case dumpMaterializedView:
		label = "Materialized view"
	}
	_, err := fmt.Fprintf(out, "--\n-- %s structure for %s\n--\n%s;\n\n", label, quoteBacktick(o.name), o.stmt)
	return err
}

// dumpInfoString writes s as the body of the double-quoted string literal that
// mv_infos() takes its database argument as.
func dumpInfoString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}
