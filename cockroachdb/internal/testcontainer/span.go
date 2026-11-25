package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// SpanDatabase is the database a span fixture is built in, and the first
// part of every name a span answer gives.
const SpanDatabase = "db"

// Spans asks the pinned CockroachDB what a query reads, for
// cockroachdb/analysis's span.tsv. Its tape holds one raw exchange per
// question, and each answer is what the server stated, reduced to JSON with
// every descriptor ID resolved to a name, so no ID reaches the tape:
//
//   - Fixture builds a starting schema in an empty database and answers the
//     schema it stored, read as cockroachdb/driver's sync reads it, and the
//     table column each view column comes from.
//   - Span prepares a query under a search_path and answers each output
//     column's name and, where the RowDescription names one, the table
//     column it comes from (its table OID and attribute number, resolved
//     through pg_class and pg_attribute); and the relations the query
//     reads: the backward dependencies of a view over the query.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans on the default pin: a replay of its tape, or
// the pinned server in a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	var live func(string) (string, error)
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode != enginetest.ModeReplay {
		pin, err := enginetest.LoadPin(Pins, Default)
		if err != nil {
			t.Fatalf("%v", err)
		}
		l := &spanLive{in: started(t, pin)}
		t.Cleanup(l.close)
		live = func(req string) (string, error) { return l.do(context.Background(), req) }
	}
	return &Spans{raw: tape.OpenRaw(t, server, Default, live)}
}

// SpanFixture is the schema CockroachDB stored for a fixture.
type SpanFixture struct {
	Database string       `json:"database"`
	Schemas  []SpanSchema `json:"schemas"`
}

// SpanSchema is one schema of the fixture's database.
type SpanSchema struct {
	Name              string         `json:"name"`
	Tables            []SpanRelation `json:"tables,omitempty"`
	Views             []SpanRelation `json:"views,omitempty"`
	MaterializedViews []SpanRelation `json:"materializedViews,omitempty"`
}

// SpanRelation is a table or a view and its columns in order, as
// information_schema.columns lists them. A view's definition is
// pg_views.definition (pg_matviews.definition for a materialized view), as
// the sync reads it: CockroachDB stores it with every name qualified by
// database and schema.
type SpanRelation struct {
	Name       string       `json:"name"`
	Columns    []SpanColumn `json:"columns"`
	Definition string       `json:"definition,omitempty"`
}

// SpanColumn is a column, its type, and whether it is hidden: a column such
// as the rowid CockroachDB adds to a table with no primary key, which
// information_schema.columns lists with is_hidden YES and * leaves out. A
// view's column also has the table column its RowDescription names for it,
// as database.schema.table.column, where it names one.
type SpanColumn struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Hidden bool   `json:"hidden,omitempty"`
	Origin string `json:"origin,omitempty"`
}

// SpanAnswer is what CockroachDB states about one query.
type SpanAnswer struct {
	// Error is the SQLSTATE of a query CockroachDB refused; the rest is
	// empty.
	Error string `json:"error,omitempty"`
	// Columns is the output columns, in order.
	Columns []SpanOutput `json:"columns,omitempty"`
	// Tables is every relation the query reads, as database.schema.name.
	Tables []string `json:"tables,omitempty"`
	// TablesUnstated marks a statement no view can hold, such as one with
	// AS OF SYSTEM TIME, an UPSERT, a SHOW or an EXPLAIN, whose relations
	// CockroachDB therefore does not state.
	TablesUnstated bool `json:"tablesUnstated,omitempty"`
}

// SpanOutput is one output column of a query.
type SpanOutput struct {
	Name string `json:"name"`
	// Source is the table column the RowDescription names for it, as
	// database.schema.table.column, and empty where it names none.
	Source string `json:"source,omitempty"`
}

// Fixture builds the starting schema sql, one statement per line, in an
// empty database and returns what CockroachDB stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span prepares statement under searchPath, a search_path value, and returns
// what CockroachDB states it reads.
func (s *Spans) Span(searchPath, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + searchPath + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

// spanLive answers the requests of a recording on one connection to the
// fixture's database.
type spanLive struct {
	in   *instance
	conn *pgx.Conn
}

func (l *spanLive) close() {
	if l.conn != nil {
		_ = l.conn.Close(context.Background())
		l.conn = nil
	}
}

func (l *spanLive) do(ctx context.Context, req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	var (
		v   any
		err error
	)
	switch kind, arg, _ := strings.Cut(head, " "); kind {
	case "fixture":
		v, err = l.fixture(ctx, body)
	case "span":
		v, err = l.span(ctx, arg, body)
	default:
		return "", fmt.Errorf("testcontainer: unknown span request %q", head)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

// spanSchemas leaves out the schemas the sync leaves out.
const spanSchemas = `n.nspname NOT IN ('crdb_internal', 'information_schema', 'pg_catalog', 'pg_extension')`

func (l *spanLive) fixture(ctx context.Context, sql string) (SpanFixture, error) {
	fx := SpanFixture{Database: SpanDatabase}
	l.close()
	admin, err := pgx.ConnectConfig(ctx, l.in.config("defaultdb"))
	if err != nil {
		return fx, err
	}
	defer admin.Close(ctx)
	for _, stmt := range []string{"DROP DATABASE IF EXISTS " + quote(SpanDatabase) + " CASCADE", "CREATE DATABASE " + quote(SpanDatabase)} {
		if _, err := admin.Exec(ctx, stmt); err != nil {
			return fx, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	cfg := l.in.config(SpanDatabase)
	// Prepare describes a query through the extended protocol.
	cfg.DefaultQueryExecMode = pgx.QueryExecModeExec
	if l.conn, err = pgx.ConnectConfig(ctx, cfg); err != nil {
		return fx, err
	}
	// One statement at a time: CockroachDB runs a script of several in one
	// implicit transaction, whose schema changes a later statement may not
	// see.
	for _, stmt := range strings.Split(sql, "\n") {
		if stmt = strings.TrimSpace(stmt); stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if _, err := l.conn.Exec(ctx, stmt); err != nil {
			return fx, fmt.Errorf("fixture: %s: %w", stmt, err)
		}
	}
	rows, err := l.conn.Query(ctx, `SELECT n.nspname FROM pg_catalog.pg_namespace n WHERE `+spanSchemas+` ORDER BY 1`)
	if err != nil {
		return fx, err
	}
	schemas, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return fx, err
	}
	bySchema := map[string]*SpanSchema{}
	for _, name := range schemas {
		fx.Schemas = append(fx.Schemas, SpanSchema{Name: name})
	}
	for i := range fx.Schemas {
		bySchema[fx.Schemas[i].Name] = &fx.Schemas[i]
	}
	type relation struct {
		Schema, Name, Kind string
		Definition         *string
	}
	rows, err = l.conn.Query(ctx, `SELECT n.nspname, c.relname, c.relkind::STRING,
	COALESCE((SELECT v.definition FROM pg_catalog.pg_views v WHERE v.schemaname = n.nspname AND v.viewname = c.relname),
		(SELECT m.definition FROM pg_catalog.pg_matviews m WHERE m.schemaname = n.nspname AND m.matviewname = c.relname))
FROM pg_catalog.pg_class c
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('r', 'v', 'm') AND `+spanSchemas+`
ORDER BY 1, 2`)
	if err != nil {
		return fx, err
	}
	rels, err := pgx.CollectRows(rows, pgx.RowToStructByPos[relation])
	if err != nil {
		return fx, err
	}
	for _, r := range rels {
		s := bySchema[r.Schema]
		if s == nil {
			continue
		}
		rel := SpanRelation{Name: r.Name}
		if r.Definition != nil {
			rel.Definition = *r.Definition
		}
		// The columns as the sync reads them.
		type column struct {
			Name, Type, Hidden string
		}
		rows, err := l.conn.Query(ctx, `SELECT column_name, data_type, is_hidden FROM information_schema.columns
WHERE table_schema = $1 AND table_name = $2 ORDER BY ordinal_position`, r.Schema, r.Name)
		if err != nil {
			return fx, err
		}
		cols, err := pgx.CollectRows(rows, pgx.RowToStructByPos[column])
		if err != nil {
			return fx, err
		}
		for _, c := range cols {
			rel.Columns = append(rel.Columns, SpanColumn{Name: c.Name, Type: c.Type, Hidden: c.Hidden == "YES"})
		}
		switch r.Kind {
		case "v":
			// The table column the RowDescription names for each of the
			// view's columns.
			sd, err := l.conn.Prepare(ctx, "", "SELECT * FROM "+quote(r.Schema)+"."+quote(r.Name))
			if err != nil {
				return fx, err
			}
			if len(sd.Fields) != len(rel.Columns) {
				return fx, fmt.Errorf("view %s.%s: %d columns, and * reads %d", r.Schema, r.Name, len(rel.Columns), len(sd.Fields))
			}
			for i, f := range sd.Fields {
				if rel.Columns[i].Origin, err = l.origin(ctx, f); err != nil {
					return fx, err
				}
			}
			s.Views = append(s.Views, rel)
		case "m":
			s.MaterializedViews = append(s.MaterializedViews, rel)
		default:
			s.Tables = append(s.Tables, rel)
		}
	}
	return fx, nil
}

// origin is the table column a RowDescription field names, as
// database.schema.table.column, or empty where it names none.
func (l *spanLive) origin(ctx context.Context, f pgconn.FieldDescription) (string, error) {
	if f.TableOID == 0 {
		return "", nil
	}
	rows, err := l.conn.Query(ctx, `SELECT pg_catalog.current_database() || '.' || n.nspname || '.' || c.relname || '.' || a.attname
FROM pg_catalog.pg_class c
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
	JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid
WHERE c.oid = $1 AND a.attnum = $2`, f.TableOID, int16(f.TableAttributeNumber))
	if err != nil {
		return "", err
	}
	names, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return "", err
	}
	if len(names) != 1 {
		return "", fmt.Errorf("column %s's origin %d.%d names %d columns", f.Name, f.TableOID, int16(f.TableAttributeNumber), len(names))
	}
	return names[0], nil
}

// probeView is the view whose dependencies name the relations a query reads.
const probeView = "omni_span_probe"

func (l *spanLive) span(ctx context.Context, searchPath, sql string) (SpanAnswer, error) {
	var a SpanAnswer
	if l.conn == nil {
		return a, errors.New("testcontainer: a span request before any fixture")
	}
	if _, err := l.conn.Exec(ctx, "SELECT pg_catalog.set_config('search_path', $1, false)", searchPath); err != nil {
		return a, err
	}
	sd, err := l.conn.Prepare(ctx, "", sql)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr):
		a.Error = pgErr.Code
		return a, nil
	case err != nil:
		return a, err
	}
	for _, f := range sd.Fields {
		out := SpanOutput{Name: f.Name}
		if out.Source, err = l.origin(ctx, f); err != nil {
			return a, err
		}
		a.Columns = append(a.Columns, out)
	}
	// A view records as its dependencies every table and view its query
	// names, in any clause, subquery or common table expression, and a view
	// itself rather than what the view reads. Its own column names keep two
	// output columns of one name apart. It is made in public, outside the
	// search_path's reach of the query it holds, which it has already read.
	view := "public." + probeView
	list := ""
	if len(sd.Fields) > 0 {
		list = " (" + strings.Join(namesByPosition(len(sd.Fields)), ", ") + ")"
	}
	if _, err := l.conn.Exec(ctx, "CREATE VIEW "+view+list+" AS "+sql); err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "42601" {
			a.TablesUnstated = true
			return a, nil
		}
		return a, fmt.Errorf("a view over the query, which reads its relations: %w", err)
	}
	defer l.conn.Exec(ctx, "DROP VIEW IF EXISTS "+view)
	rows, err := l.conn.Query(ctx, `SELECT DISTINCT pg_catalog.current_database() || '.' || n.nspname || '.' || c.relname
FROM crdb_internal.backward_dependencies d
	JOIN pg_catalog.pg_class c ON c.oid = d.dependson_id
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE d.descriptor_id = $1::REGCLASS::OID AND c.relkind IN ('r', 'v', 'm')
ORDER BY 1`, view)
	if err != nil {
		return a, err
	}
	if a.Tables, err = pgx.CollectRows(rows, pgx.RowTo[string]); err != nil {
		return a, err
	}
	return a, nil
}

// namesByPosition names n columns c1, c2, and so on.
func namesByPosition(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("c%d", i+1)
	}
	return out
}
