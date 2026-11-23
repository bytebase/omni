package testcontainer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned StarRocks what a query reads, for
// starrocks/analysis's span.tsv. Its tape holds one raw exchange per
// question, and each answer is what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on a server emptied of every database
//     the image did not start with, and answers the schema StarRocks
//     stored: each database's tables, views and materialized views with
//     their columns in order, and each view's and materialized view's
//     definition as information_schema holds it.
//   - Span runs a query in a session's database and answers the result
//     set's column definitions, the tables and views the analyzer bound
//     (the query dump's table_meta and view_meta), and the query's logical
//     plan (EXPLAIN LOGICAL), or the error StarRocks refused it with.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned
// StarRocks in a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	var live func(string) (string, error)
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode != enginetest.ModeReplay {
		pin, err := enginetest.LoadPin(Pins, "default")
		if err != nil {
			t.Fatalf("%v", err)
		}
		cfg := config(t, pin)
		// The fixture drops every database the image did not start with,
		// the judge's own among them, so no connection selects one.
		cfg.DBName = ""
		cfg.MultiStatements = false
		cfg.ReadTimeout, cfg.WriteTimeout = 2*time.Minute, 2*time.Minute
		l := &spanLive{cfg: cfg, http: shared.httpAddr, root: sql.OpenDB(must(mysql.NewConnector(cfg)))}
		l.root.SetMaxOpenConns(1)
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, "default", live)}
}

// SpanFixture is the schema StarRocks stored for a fixture.
type SpanFixture struct {
	Databases []SpanDatabase `json:"databases"`
}

// SpanDatabase is one database a fixture made.
type SpanDatabase struct {
	Name              string         `json:"name"`
	Tables            []SpanRelation `json:"tables"`
	Views             []SpanRelation `json:"views"`
	MaterializedViews []SpanRelation `json:"materializedViews,omitempty"`
}

// SpanRelation is a table, view or materialized view and its columns in
// order.
type SpanRelation struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	// Definition is a view's VIEW_DEFINITION, or a materialized view's
	// MATERIALIZED_VIEW_DEFINITION.
	Definition string `json:"definition,omitempty"`
}

// SpanAnswer is what StarRocks states about one query.
type SpanAnswer struct {
	// Error is the error StarRocks refused the query with, as "<number>
	// <message>"; the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column names, in order.
	Columns []string `json:"columns,omitempty"`
	// Tables and Views are what the query dump says the analyzer bound, as
	// database.name, in order.
	Tables []string `json:"tables,omitempty"`
	Views  []string `json:"views,omitempty"`
	// Plan is the logical plan, without its estimates.
	Plan []string `json:"plan,omitempty"`
}

// Fixture builds the starting schema sql, one statement per line, on a
// server emptied of every database it was not started with, and returns
// what StarRocks stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span runs statement with database selected and returns what StarRocks
// states it reads.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + database + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

type spanLive struct {
	cfg  *mysql.Config
	http string
	root *sql.DB
	wire *wire
}

func (l *spanLive) close() {
	if l.wire != nil {
		l.wire.Close()
	}
	l.root.Close()
}

func (l *spanLive) do(req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	var (
		v   any
		err error
	)
	switch {
	case head == "fixture":
		v, err = l.fixture(context.Background(), body)
	case strings.HasPrefix(head, "span "):
		v, err = l.span(strings.TrimPrefix(head, "span "), body)
	default:
		return "", fmt.Errorf("spans: an unknown request %q", head)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

func (l *spanLive) fixture(ctx context.Context, sqlText string) (SpanFixture, error) {
	var fx SpanFixture
	databases, err := queryRows(ctx, l.root, "SHOW DATABASES")
	if err != nil {
		return fx, err
	}
	for _, name := range databases {
		if shared.databases[strings.ToLower(name)] {
			continue
		}
		if _, err := l.root.ExecContext(ctx, "DROP DATABASE "+quoteName(name)+" FORCE"); err != nil {
			return fx, fmt.Errorf("drop %s: %w", name, err)
		}
	}
	for _, stmt := range strings.Split(sqlText, "\n") {
		stmt = strings.TrimSuffix(strings.TrimSpace(stmt), ";")
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if _, err := l.root.ExecContext(ctx, stmt); err != nil {
			return fx, fmt.Errorf("build the fixture: %s: %w", stmt, err)
		}
	}
	type column struct {
		db, table, kind, def, mvDef, name string
		pos                               int
	}
	// materialized_views asks a back end, which a new server may not have
	// registered yet.
	var rows *sql.Rows
	deadline := time.Now().Add(3 * time.Minute)
	for {
		rows, err = l.root.QueryContext(ctx, schemaQuery)
		if err == nil || time.Now().After(deadline) || !strings.Contains(err.Error(), "listMaterializedViewStatus") {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return fx, err
	}
	defer rows.Close()
	var columns []column
	for rows.Next() {
		var c column
		if err := rows.Scan(&c.db, &c.table, &c.kind, &c.def, &c.mvDef, &c.name, &c.pos); err != nil {
			return fx, err
		}
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		return fx, err
	}
	// StarRocks answers in no fixed order.
	slices.SortFunc(columns, func(a, b column) int {
		if n := strings.Compare(a.db, b.db); n != 0 {
			return n
		}
		if n := strings.Compare(a.table, b.table); n != 0 {
			return n
		}
		return a.pos - b.pos
	})
	for _, c := range columns {
		if n := len(fx.Databases); n == 0 || fx.Databases[n-1].Name != c.db {
			fx.Databases = append(fx.Databases, SpanDatabase{Name: c.db})
		}
		d := &fx.Databases[len(fx.Databases)-1]
		list, def := &d.Tables, ""
		switch {
		case c.mvDef != "":
			list, def = &d.MaterializedViews, c.mvDef
		case c.kind == "VIEW":
			list, def = &d.Views, c.def
		}
		if n := len(*list); n == 0 || (*list)[n-1].Name != c.table {
			*list = append(*list, SpanRelation{Name: c.table, Definition: def})
		}
		r := &(*list)[len(*list)-1]
		r.Columns = append(r.Columns, c.name)
	}
	return fx, nil
}

// schemaQuery reads the schema a fixture stored. A materialized view is a
// VIEW in information_schema.TABLES, with no row in VIEWS;
// materialized_views holds its definition.
const schemaQuery = `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), COALESCE(m.MATERIALIZED_VIEW_DEFINITION, ''), c.COLUMN_NAME, c.ORDINAL_POSITION
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON v.TABLE_SCHEMA = t.TABLE_SCHEMA AND v.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.materialized_views m ON m.TABLE_SCHEMA = t.TABLE_SCHEMA AND m.TABLE_NAME = t.TABLE_NAME
		WHERE t.TABLE_SCHEMA NOT IN ('information_schema', 'sys', '_statistics_')`

func (l *spanLive) span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	if l.wire == nil {
		w, err := dialWire(l.cfg.Addr, l.cfg.User)
		if err != nil {
			return a, err
		}
		// A rewrite to a materialized view would swap the tables a plan
		// scans for one the query does not name.
		// A busy machine plans slowly, and the planner gives up after 3s.
		for _, stmt := range []string{"SET enable_materialized_view_rewrite = false", "SET new_planner_optimize_timeout = 60000"} {
			if err := w.Exec(stmt); err != nil {
				w.Close()
				return a, fmt.Errorf("%s: %w", stmt, err)
			}
		}
		l.wire = w
	}
	w := l.wire
	if err := w.Use(database); err != nil {
		return a, fmt.Errorf("select %s: %w", database, err)
	}
	cols, _, err := w.Query(statement)
	for try := 0; try < 30 && transient(err); try++ {
		time.Sleep(2 * time.Second)
		cols, _, err = w.Query(statement)
	}
	var refused *wireError
	if errors.As(err, &refused) {
		a.Error = fmt.Sprintf("%d %s", refused.Number, refused.Message)
		return a, nil
	}
	if err != nil {
		return a, err
	}
	for _, c := range cols {
		a.Columns = append(a.Columns, c.Name)
	}
	_, plan, err := w.Query("EXPLAIN LOGICAL " + statement)
	for try := 0; try < 30 && transient(err); try++ {
		time.Sleep(2 * time.Second)
		_, plan, err = w.Query("EXPLAIN LOGICAL " + statement)
	}
	if err != nil {
		return a, fmt.Errorf("EXPLAIN LOGICAL: %w", err)
	}
	for _, row := range plan {
		for _, line := range strings.Split(strings.Join(row, " "), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "Estimates:") || strings.HasPrefix(trimmed, "partitionRatio:") {
				continue
			}
			a.Plan = append(a.Plan, line)
		}
	}
	dump, err := l.dump(database, statement)
	if err != nil {
		return a, fmt.Errorf("dump the query: %w", err)
	}
	for name := range dump.Tables {
		a.Tables = append(a.Tables, name)
	}
	for name := range dump.Views {
		a.Views = append(a.Views, name)
	}
	// The dump is a JSON object, which Go reads in no order.
	slices.Sort(a.Tables)
	slices.Sort(a.Views)
	return a, nil
}

// transient reports whether err is the server's own trouble rather than
// its answer to the query: a back end it briefly blacklisted.
func transient(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Backend node not found")
}

// queryDump is the part of a query dump (/api/query_dump) the suite reads:
// the CREATE statement of every table, and the definition of every view,
// the analyzer bound, keyed by database.name.
type queryDump struct {
	Tables    map[string]string `json:"table_meta"`
	Views     map[string]string `json:"view_meta"`
	Exception []string          `json:"exception"`
}

func (l *spanLive) dump(database, statement string) (queryDump, error) {
	var d queryDump
	// mock=false keeps the names; a dump masks them by default.
	u := "http://" + l.http + "/api/query_dump?mock=false&db=" + url.QueryEscape(database)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(statement))
	if err != nil {
		return d, err
	}
	req.SetBasicAuth(l.cfg.User, l.cfg.Passwd)
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}
	if resp.StatusCode != http.StatusOK {
		return d, fmt.Errorf("%s: %s", resp.Status, body)
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return d, fmt.Errorf("%w: %s", err, body)
	}
	if len(d.Exception) > 0 {
		return d, fmt.Errorf("the dump failed: %v", d.Exception)
	}
	return d, nil
}
