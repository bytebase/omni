package testcontainer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// SpanRole is the server mysql/analysis's span.tsv is judged on: 8.4 is the
// first release whose EXPLAIN FORMAT=JSON (explain_json_format_version 2)
// names the schema and table of each access, where 8.0's names its alias.
const SpanRole = MySQL84

// Spans asks the pinned MySQL what a query reads, for mysql/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on an emptied server and answers the
//     schema MySQL stored: each database's tables and views with their
//     columns in order, and each view's definition as information_schema
//     holds it.
//   - Span prepares a query in a session's database and answers the result
//     set's column definitions, the query as MySQL resolved it (the
//     fully qualified text EXPLAIN reports, after views and derived tables
//     are merged), the tables the plan accesses, and whether the query read
//     a view (the optimizer trace's view steps) or the plan left out a table
//     it read while optimizing.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned MySQL 8.4
// in a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	var live func(string) (string, error)
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode != enginetest.ModeReplay {
		pin, err := enginetest.LoadPin(Pins, SpanRole)
		if err != nil {
			t.Fatalf("%v", err)
		}
		server := servers[SpanRole]
		if err := server.start(context.Background(), pin); err != nil {
			t.Fatalf("MySQL (%s): %v", SpanRole, err)
		}
		l := &spanLive{server: server, root: sql.OpenDB(server.connector(t, pin, Options{MultiStatements: true}))}
		l.root.SetMaxOpenConns(1)
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, SpanRole, live)}
}

// SpanFixture is the schema MySQL stored for a fixture.
type SpanFixture struct {
	Databases []SpanDatabase `json:"databases"`
}

// SpanDatabase is one database a fixture made.
type SpanDatabase struct {
	Name   string         `json:"name"`
	Tables []SpanRelation `json:"tables"`
	Views  []SpanRelation `json:"views"`
}

// SpanRelation is a table or a view and its columns in order.
type SpanRelation struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	// Definition is a view's VIEW_DEFINITION.
	Definition string `json:"definition,omitempty"`
}

// SpanAnswer is what MySQL states about one query.
type SpanAnswer struct {
	// Error is the error number of a query MySQL refused; the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column definitions, in order.
	Columns []SpanColumn `json:"columns,omitempty"`
	// Resolved is the query as EXPLAIN reports MySQL resolved it.
	Resolved string `json:"resolved,omitempty"`
	// Tables is every table the plan accesses, as database.table.
	Tables []string `json:"tables,omitempty"`
	// ReadsView is set when the optimizer trace merged or materialized a
	// view: the plan names a view's tables, never the view.
	ReadsView bool `json:"readsView,omitempty"`
	// Hidden is set when the plan leaves out a table it read while
	// optimizing, such as a constant table or one an aggregate was read
	// from.
	Hidden bool `json:"hidden,omitempty"`
}

// SpanColumn is one column definition of a result set.
type SpanColumn struct {
	Name     string `json:"name"`
	Schema   string `json:"schema,omitempty"`
	Table    string `json:"table,omitempty"`
	OrgTable string `json:"orgTable,omitempty"`
	OrgName  string `json:"orgName,omitempty"`
}

// Fixture builds the starting schema sql on a server emptied of every
// database it was not started with, and returns what MySQL stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span prepares statement with database selected and returns what MySQL
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
	server *Server
	root   *sql.DB
	wire   *wire
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
	dbs, err := queryPairs(ctx, l.root, "SELECT SCHEMA_NAME, '' FROM information_schema.SCHEMATA")
	if err != nil {
		return fx, err
	}
	var stmts []string
	for _, d := range dbs {
		if !systemDatabases[strings.ToLower(d[0])] {
			stmts = append(stmts, "DROP DATABASE `"+strings.ReplaceAll(d[0], "`", "``")+"`")
		}
	}
	for _, stmt := range stmts {
		if _, err := l.root.ExecContext(ctx, stmt); err != nil {
			return fx, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	if _, err := l.root.ExecContext(ctx, sqlText); err != nil {
		return fx, fmt.Errorf("build the fixture: %w", err)
	}
	if _, err := l.root.ExecContext(ctx, "USE mysql"); err != nil {
		return fx, err
	}
	rows, err := l.root.QueryContext(ctx, `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), c.COLUMN_NAME
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON v.TABLE_SCHEMA = t.TABLE_SCHEMA AND v.TABLE_NAME = t.TABLE_NAME
		WHERE t.TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
		ORDER BY BINARY t.TABLE_SCHEMA, BINARY t.TABLE_NAME, c.ORDINAL_POSITION`)
	if err != nil {
		return fx, err
	}
	defer rows.Close()
	for rows.Next() {
		var db, table, kind, def, column string
		if err := rows.Scan(&db, &table, &kind, &def, &column); err != nil {
			return fx, err
		}
		if n := len(fx.Databases); n == 0 || fx.Databases[n-1].Name != db {
			fx.Databases = append(fx.Databases, SpanDatabase{Name: db})
		}
		d := &fx.Databases[len(fx.Databases)-1]
		list := &d.Tables
		if kind == "VIEW" {
			list = &d.Views
		}
		if n := len(*list); n == 0 || (*list)[n-1].Name != table {
			*list = append(*list, SpanRelation{Name: table, Definition: def})
		}
		r := &(*list)[len(*list)-1]
		r.Columns = append(r.Columns, column)
	}
	return fx, rows.Err()
}

func (l *spanLive) span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	if l.wire == nil {
		// As root: the optimizer trace of a query that reads a view is
		// withheld from any account but the view's definer.
		w, err := dialWire(l.server.config.Addr, l.server.config.User, l.server.config.Passwd)
		if err != nil {
			return a, err
		}
		for _, stmt := range []string{
			"SET explain_json_format_version = 2",
			"SET optimizer_trace_max_mem_size = 67108864",
		} {
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
	cols, err := w.Prepare(statement)
	var refused *wireError
	if errors.As(err, &refused) {
		a.Error = fmt.Sprint(refused.Number)
		return a, nil
	}
	if err != nil {
		return a, err
	}
	for _, c := range cols {
		a.Columns = append(a.Columns, SpanColumn{Name: c.Name, Schema: c.Schema, Table: c.Table, OrgTable: c.OrgTable, OrgName: c.OrgName})
	}

	if err := w.Exec("SET optimizer_trace = 'enabled=on'"); err != nil {
		return a, err
	}
	plan, err := w.Query("EXPLAIN FORMAT=JSON " + statement)
	if err != nil {
		return a, fmt.Errorf("EXPLAIN: %w", err)
	}
	trace, err := w.Query("SELECT TRACE, MISSING_BYTES_BEYOND_MAX_MEM_SIZE FROM information_schema.OPTIMIZER_TRACE")
	if err != nil {
		return a, fmt.Errorf("read the optimizer trace: %w", err)
	}
	if err := w.Exec("SET optimizer_trace = 'enabled=off'"); err != nil {
		return a, err
	}
	if len(plan) != 1 || len(trace) != 1 {
		return a, fmt.Errorf("EXPLAIN answered %d rows and the trace %d", len(plan), len(trace))
	}
	if trace[0][1] != "0" {
		return a, fmt.Errorf("the optimizer trace lost %s bytes", trace[0][1])
	}
	var tr any
	if err := json.Unmarshal([]byte(trace[0][0]), &tr); err != nil {
		return a, fmt.Errorf("read the optimizer trace: %w", err)
	}
	a.ReadsView = hasKey(tr, "view")

	var p map[string]any
	if err := json.Unmarshal([]byte(plan[0][0]), &p); err != nil {
		return a, fmt.Errorf("read the plan: %w", err)
	}
	a.Resolved, _ = p["query"].(string)
	if a.Resolved == "" {
		return a, errors.New("the plan names no resolved query")
	}
	seen := map[string]bool{}
	walkPlan(p, func(node map[string]any) {
		schema, _ := node["schema_name"].(string)
		table, _ := node["table_name"].(string)
		if schema != "" && table != "" && !seen[schema+"."+table] {
			seen[schema+"."+table] = true
			a.Tables = append(a.Tables, schema+"."+table)
		}
		op, _ := node["operation"].(string)
		for _, hint := range []string{"before execution", "optimized away", "Impossible", "no matching row", "Zero rows"} {
			if strings.Contains(op, hint) {
				a.Hidden = true
			}
		}
	})
	// The plan is a JSON object, which Go reads in no order.
	slices.Sort(a.Tables)
	return a, nil
}

// walkPlan calls fn on every object of an EXPLAIN FORMAT=JSON plan.
func walkPlan(v any, fn func(map[string]any)) {
	switch v := v.(type) {
	case map[string]any:
		fn(v)
		for _, c := range v {
			walkPlan(c, fn)
		}
	case []any:
		for _, c := range v {
			walkPlan(c, fn)
		}
	}
}

// hasKey reports whether any object of v holds key.
func hasKey(v any, key string) bool {
	found := false
	walkPlan(v, func(m map[string]any) {
		if _, ok := m[key]; ok {
			found = true
		}
	})
	return found
}
