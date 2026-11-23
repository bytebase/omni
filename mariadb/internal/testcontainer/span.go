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

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned MariaDB what a query reads, for mariadb/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on an emptied server and answers the
//     schema MariaDB stored: each database's tables and views with their
//     columns in order, and each view's definition as information_schema
//     holds it.
//   - Span runs a query in a session's database and answers the result set's
//     column names, the query as MariaDB resolved it (the fully qualified
//     text of the note EXPLAIN EXTENDED leaves, after views and derived
//     tables are merged), the plan's rows, and what the optimizer trace says
//     it merged, materialized or eliminated.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned MariaDB in
// a recording.
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
		db := sql.OpenDB(Options{MultiStatements: true}.live(t, pin))
		db.SetMaxOpenConns(1)
		t.Cleanup(func() { _ = db.Close() })
		l := &spanLive{db: db}
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, "default", live)}
}

// SpanFixture is the schema MariaDB stored for a fixture.
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

// SpanAnswer is what MariaDB states about one query.
type SpanAnswer struct {
	// Error is the error number of a query MariaDB refused; the rest is
	// empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column names, in order.
	Columns []string `json:"columns,omitempty"`
	// Resolved is the query as EXPLAIN EXTENDED's note 1003 says MariaDB
	// resolved and optimized it.
	Resolved string `json:"resolved,omitempty"`
	// Plan is each row of EXPLAIN EXTENDED: its table, access type and
	// Extra, separated by tabs.
	Plan []string `json:"plan,omitempty"`
	// Views is each view the optimizer trace merged or materialized.
	Views []string `json:"views,omitempty"`
	// Eliminated is each table the optimizer trace eliminated from a join.
	Eliminated []string `json:"eliminated,omitempty"`
}

// Fixture builds the starting schema sql on a server emptied of every
// database it was not started with, and returns what MariaDB stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span runs statement with database selected and returns what MariaDB states
// it reads.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + database + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

type spanLive struct {
	db *sql.DB
}

func (l *spanLive) do(req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	ctx := context.WithValue(context.Background(), roleKey{}, "default")
	var (
		v   any
		err error
	)
	switch {
	case head == "fixture":
		v, err = l.fixture(ctx, body)
	case strings.HasPrefix(head, "span "):
		v, err = l.span(ctx, strings.TrimPrefix(head, "span "), body)
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
	if err := empty(ctx, l.db); err != nil {
		return fx, err
	}
	if _, err := l.db.ExecContext(ctx, sqlText); err != nil {
		return fx, fmt.Errorf("build the fixture: %w", err)
	}
	// information_schema compares names without case, so the joins compare
	// them as bytes: MixedCase and mixedcase are two tables.
	rows, err := l.db.QueryContext(ctx, `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), c.COLUMN_NAME
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON BINARY c.TABLE_SCHEMA = t.TABLE_SCHEMA AND BINARY c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON BINARY v.TABLE_SCHEMA = t.TABLE_SCHEMA AND BINARY v.TABLE_NAME = t.TABLE_NAME
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

func (l *spanLive) span(ctx context.Context, database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	conn, err := l.db.Conn(ctx)
	if err != nil {
		return a, err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "USE "+quoteName(database)); err != nil {
		return a, fmt.Errorf("select %s: %w", database, err)
	}
	// The result set's names: the fixture's tables are small, so the query
	// runs whole.
	rows, err := conn.QueryContext(ctx, statement)
	var refused *mysql.MySQLError
	if errors.As(err, &refused) {
		a.Error = fmt.Sprint(refused.Number)
		return a, nil
	}
	if err != nil {
		return a, err
	}
	a.Columns, err = rows.Columns()
	if err == nil {
		for rows.Next() {
		}
		err = rows.Err()
	}
	rows.Close()
	if err != nil {
		return a, err
	}

	if _, err := conn.ExecContext(ctx, "SET optimizer_trace = 'enabled=on'"); err != nil {
		return a, err
	}
	plan, err := textRows(ctx, conn, "EXPLAIN EXTENDED "+statement)
	if err != nil {
		return a, fmt.Errorf("EXPLAIN: %w", err)
	}
	for _, r := range plan {
		// id, select_type, table, type, possible_keys, key, key_len, ref,
		// rows, filtered, Extra
		if len(r) != 11 {
			return a, fmt.Errorf("an EXPLAIN row of %d columns", len(r))
		}
		a.Plan = append(a.Plan, r[2]+"\t"+r[3]+"\t"+r[10])
	}
	warnings, err := textRows(ctx, conn, "SHOW WARNINGS")
	if err != nil {
		return a, err
	}
	for _, w := range warnings {
		if len(w) == 3 && w[1] == "1003" {
			a.Resolved = w[2]
		}
	}
	if a.Resolved == "" {
		return a, errors.New("EXPLAIN EXTENDED left no resolved query")
	}
	trace, err := textRows(ctx, conn, "SELECT TRACE, MISSING_BYTES_BEYOND_MAX_MEM_SIZE FROM information_schema.OPTIMIZER_TRACE")
	if err != nil {
		return a, fmt.Errorf("read the optimizer trace: %w", err)
	}
	if _, err := conn.ExecContext(ctx, "SET optimizer_trace = 'enabled=off'"); err != nil {
		return a, err
	}
	if len(trace) != 1 {
		return a, fmt.Errorf("the optimizer trace has %d rows", len(trace))
	}
	if trace[0][1] != "0" {
		return a, fmt.Errorf("the optimizer trace lost %s bytes", trace[0][1])
	}
	var tr any
	if err := json.Unmarshal([]byte(trace[0][0]), &tr); err != nil {
		return a, fmt.Errorf("read the optimizer trace: %w", err)
	}
	walkJSON(tr, func(m map[string]any) {
		if v, ok := m["view"].(map[string]any); ok {
			name, _ := v["table"].(string)
			a.Views = append(a.Views, name)
		}
		if list, ok := m["eliminated_tables"].([]any); ok {
			for _, e := range list {
				name, _ := e.(string)
				a.Eliminated = append(a.Eliminated, name)
			}
		}
	})
	// The trace is a JSON object, which Go reads in no order.
	slices.Sort(a.Views)
	slices.Sort(a.Eliminated)
	return a, nil
}

// textRows reads every row of query as text, NULL as "NULL".
func textRows(ctx context.Context, conn *sql.Conn, query string) ([][]string, error) {
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out [][]string
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			row[i] = v.String
			if !v.Valid {
				row[i] = "NULL"
			}
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// walkJSON calls fn on every object of v.
func walkJSON(v any, fn func(map[string]any)) {
	switch v := v.(type) {
	case map[string]any:
		fn(v)
		for _, c := range v {
			walkJSON(c, fn)
		}
	case []any:
		for _, c := range v {
			walkJSON(c, fn)
		}
	}
}
