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

// Spans asks a pinned PostgreSQL what a query reads, for pg/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON with every OID resolved to a name,
// so no OID reaches the tape:
//
//   - Fixture builds a starting schema in an empty database and answers the
//     schema it stored, read as pg/driver's sync reads it.
//   - Span prepares a query under a search_path and answers each output
//     column's name and, where the RowDescription names one, the table
//     column it comes from (its table OID and attribute number, resolved
//     through pg_class and pg_attribute); and the relations the query
//     reads: the relations the rewrite rule of a temporary view over the
//     query depends on.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans on the server of pin row role: a replay of its
// tape, or the pinned server in a recording.
func OpenSpans(t testing.TB, role string) *Spans {
	t.Helper()
	var live func(string) (string, error)
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode != enginetest.ModeReplay {
		pin, err := enginetest.LoadPin(Pins, role)
		if err != nil {
			t.Fatalf("%v", err)
		}
		in := started(t, pin)
		in.empty(t)
		l := &spanLive{in: in}
		t.Cleanup(l.close)
		live = func(req string) (string, error) { return l.do(context.Background(), req) }
	}
	return &Spans{raw: tape.OpenRaw(t, server, role, live)}
}

// SpanFixture is the schema PostgreSQL stored for a fixture.
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

// SpanRelation is a table or a view and its columns in order. A view's
// definition is pg_views.definition read with an empty search_path, as the
// sync reads it, so every name in it is qualified. A partitioned table lists
// its partitions, which the sync stores with it rather than as tables.
type SpanRelation struct {
	Name       string       `json:"name"`
	Columns    []SpanColumn `json:"columns"`
	Definition string       `json:"definition,omitempty"`
	Partitions []string     `json:"partitions,omitempty"`
}

// SpanColumn is a column and its type.
type SpanColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SpanAnswer is what PostgreSQL states about one query.
type SpanAnswer struct {
	// Error is the SQLSTATE of a query PostgreSQL refused; the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the output columns, in order.
	Columns []SpanOutput `json:"columns,omitempty"`
	// Tables is every relation the query reads, as database.schema.name.
	Tables []string `json:"tables,omitempty"`
}

// SpanOutput is one output column of a query.
type SpanOutput struct {
	Name string `json:"name"`
	// Sources is the table column the RowDescription names for it, as
	// database.schema.table.column, or every column of the relation it
	// names for a whole-row reference (attribute number 0), and empty where
	// it names none.
	Sources []string `json:"sources,omitempty"`
}

// Fixture builds the starting schema sql in an empty database and returns
// what PostgreSQL stored. sql runs in one simple query.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span prepares statement under searchPath, a search_path value, and returns
// what PostgreSQL states it reads.
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

func (l *spanLive) fixture(ctx context.Context, sql string) (SpanFixture, error) {
	fx := SpanFixture{Database: SpanDatabase}
	l.close()
	admin, err := pgx.ConnectConfig(ctx, l.in.config("postgres"))
	if err != nil {
		return fx, err
	}
	defer admin.Close(ctx)
	for _, stmt := range append(l.in.dropDatabase(SpanDatabase, "IF EXISTS "), "CREATE DATABASE "+quote(SpanDatabase)) {
		if _, err := admin.Exec(ctx, stmt); err != nil {
			return fx, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	if l.conn, err = pgx.ConnectConfig(ctx, l.in.config(SpanDatabase)); err != nil {
		return fx, err
	}
	if _, err := l.conn.Exec(ctx, sql); err != nil {
		return fx, fmt.Errorf("fixture: %w", err)
	}
	// What follows reads the stored schema as pg/driver's sync does: with an
	// empty search_path, so a view's definition qualifies every name.
	if _, err := l.conn.Exec(ctx, "RESET ALL; SET search_path = ''"); err != nil {
		return fx, err
	}
	defer l.conn.Exec(ctx, "RESET search_path")
	rows, err := l.conn.Query(ctx, `SELECT n.nspname FROM pg_catalog.pg_namespace n
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast') AND n.nspname NOT LIKE 'pg\_temp\_%' AND n.nspname NOT LIKE 'pg\_toast\_temp\_%'
ORDER BY 1`)
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
		Partitions         []string
		Columns            []string
		Types              []string
	}
	rows, err = l.conn.Query(ctx, `SELECT n.nspname, c.relname, c.relkind::text,
	CASE WHEN c.relkind IN ('v', 'm') THEN pg_catalog.pg_get_viewdef(c.oid) END,
	ARRAY(SELECT p.relname::text FROM pg_catalog.pg_inherits i JOIN pg_catalog.pg_class p ON p.oid = i.inhrelid WHERE i.inhparent = c.oid ORDER BY 1),
	array_agg(a.attname::text ORDER BY a.attnum), array_agg(pg_catalog.format_type(a.atttypid, NULL) ORDER BY a.attnum)
FROM pg_catalog.pg_class c
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
	JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid AND a.attnum > 0 AND NOT a.attisdropped
WHERE c.relkind IN ('r', 'p', 'v', 'm', 'f') AND NOT c.relispartition
GROUP BY n.nspname, c.relname, c.relkind, c.oid
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
		rel := SpanRelation{Name: r.Name, Partitions: r.Partitions}
		for i, c := range r.Columns {
			rel.Columns = append(rel.Columns, SpanColumn{Name: c, Type: r.Types[i]})
		}
		if r.Definition != nil {
			rel.Definition = *r.Definition
		}
		switch r.Kind {
		case "v":
			s.Views = append(s.Views, rel)
		case "m":
			s.MaterializedViews = append(s.MaterializedViews, rel)
		default:
			s.Tables = append(s.Tables, rel)
		}
	}
	return fx, nil
}

// probeView is the temporary view whose dependencies name the relations a
// query reads.
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
		if f.TableOID != 0 {
			// Attribute number 0 is a whole-row reference: every column of
			// the relation.
			rows, err := l.conn.Query(ctx, `SELECT pg_catalog.current_database() || '.' || n.nspname || '.' || c.relname || '.' || a.attname
FROM pg_catalog.pg_class c
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
	JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid
WHERE c.oid = $1 AND (a.attnum = $2 OR $2 = 0 AND a.attnum > 0 AND NOT a.attisdropped)
ORDER BY a.attnum`, f.TableOID, int16(f.TableAttributeNumber))
			if err != nil {
				return a, err
			}
			if out.Sources, err = pgx.CollectRows(rows, pgx.RowTo[string]); err != nil {
				return a, err
			}
			if len(out.Sources) == 0 {
				return a, fmt.Errorf("column %s's origin %d.%d names no column", f.Name, f.TableOID, int16(f.TableAttributeNumber))
			}
		}
		a.Columns = append(a.Columns, out)
	}
	// A view's rewrite rule depends on every relation its query names, in
	// any clause, subquery or common table expression, and on a view itself
	// rather than what the view reads. Its own column names keep two output
	// columns of one name apart.
	view := "pg_temp." + probeView
	list := ""
	if len(sd.Fields) > 0 {
		list = " (" + strings.Join(namesByPosition(len(sd.Fields)), ", ") + ")"
	}
	if _, err := l.conn.Exec(ctx, "CREATE TEMPORARY VIEW "+probeView+list+" AS "+sql); err != nil {
		return a, fmt.Errorf("a temporary view over the query, which reads its relations: %w", err)
	}
	defer l.conn.Exec(ctx, "DROP VIEW IF EXISTS "+view)
	rows, err := l.conn.Query(ctx, `SELECT DISTINCT pg_catalog.current_database() || '.' || n.nspname || '.' || c.relname
FROM pg_catalog.pg_rewrite r
	JOIN pg_catalog.pg_depend d ON d.classid = 'pg_catalog.pg_rewrite'::pg_catalog.regclass AND d.objid = r.oid
	JOIN pg_catalog.pg_class c ON d.refclassid = 'pg_catalog.pg_class'::pg_catalog.regclass AND d.refobjid = c.oid
	JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE r.ev_class = $1::pg_catalog.regclass AND c.oid <> r.ev_class AND c.relkind IN ('r', 'p', 'v', 'm', 'f')
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
