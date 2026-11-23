package testcontainer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned TiDB what a query reads, for tidb/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on an emptied server and answers the
//     schema TiDB stored: each database's tables and views with their
//     columns in order, and each view's definition as information_schema
//     holds it.
//   - Span prepares a query in a session's database and answers the result
//     set's column definitions, the rows of its plan (EXPLAIN FORMAT =
//     'brief'), and the tables the statement summary records for it once it
//     ran.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned TiDB in a
// recording.
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
		db := sql.OpenDB(Options{}.live(t, pin))
		db.SetMaxOpenConns(1)
		l := &spanLive{db: db, addr: shared.config.Addr, user: shared.config.User}
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, "default", live)}
}

// SpanFixture is the schema TiDB stored for a fixture.
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

// SpanAnswer is what TiDB states about one query.
type SpanAnswer struct {
	// Error is the error number of a query TiDB refused; the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column definitions, in order.
	Columns []SpanColumn `json:"columns,omitempty"`
	// Plan is each row of EXPLAIN FORMAT = 'brief' but its estimate: the
	// operator (with the tree drawn before it), its task, access object
	// and operator info, separated by tabs.
	Plan []string `json:"plan,omitempty"`
	// Tables is the TABLE_NAMES the statement summary records for the
	// query: the tables and views it names, as db.table.
	Tables string `json:"tables,omitempty"`
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
// database it was not started with, and returns what TiDB stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span prepares statement with database selected and returns what TiDB
// states it reads.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + database + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

// spanSystemDatabases are the databases a TiDB server starts with.
var spanSystemDatabases = map[string]bool{
	"mysql": true, "information_schema": true, "performance_schema": true,
	"metrics_schema": true, "sys": true, "test": true,
}

type spanLive struct {
	db         *sql.DB
	addr, user string
	wire       *wire
}

func (l *spanLive) close() {
	if l.wire != nil {
		l.wire.Close()
	}
	l.db.Close()
}

func (l *spanLive) do(req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	ctx := context.Background()
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
	rows, err := l.db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA")
	if err != nil {
		return fx, err
	}
	var drop []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return fx, err
		}
		if !spanSystemDatabases[strings.ToLower(name)] {
			drop = append(drop, "DROP DATABASE `"+strings.ReplaceAll(name, "`", "``")+"`")
		}
	}
	rows.Close()
	for _, stmt := range drop {
		if _, err := l.db.ExecContext(ctx, stmt); err != nil {
			return fx, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	if _, err := l.db.ExecContext(ctx, sqlText); err != nil {
		return fx, fmt.Errorf("build the fixture: %w", err)
	}
	rows, err = l.db.QueryContext(ctx, `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), c.COLUMN_NAME
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON v.TABLE_SCHEMA = t.TABLE_SCHEMA AND v.TABLE_NAME = t.TABLE_NAME
		WHERE LOWER(t.TABLE_SCHEMA) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys', 'test')
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
	if l.wire == nil {
		w, err := dialWire(l.addr, l.user, "")
		if err != nil {
			return a, err
		}
		l.wire = w
	}
	if err := l.wire.Use(database); err != nil {
		return a, fmt.Errorf("select %s: %w", database, err)
	}
	cols, err := l.wire.Prepare(statement)
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

	conn, err := l.db.Conn(ctx)
	if err != nil {
		return a, err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "USE `"+strings.ReplaceAll(database, "`", "``")+"`"); err != nil {
		return a, err
	}
	plan, err := conn.QueryContext(ctx, "EXPLAIN FORMAT = 'brief' "+statement)
	if err != nil {
		return a, fmt.Errorf("EXPLAIN: %w", err)
	}
	for plan.Next() {
		var id, est, task, object, info string
		if err := plan.Scan(&id, &est, &task, &object, &info); err != nil {
			plan.Close()
			return a, err
		}
		a.Plan = append(a.Plan, id+"\t"+task+"\t"+object+"\t"+info)
	}
	plan.Close()
	if err := plan.Err(); err != nil {
		return a, err
	}
	// The statement summary records the tables a statement names once it
	// has run; the fixture's tables are small, so it runs whole. Turning
	// the summary off and on empties it, so the query's row is its own.
	for _, stmt := range []string{"SET GLOBAL tidb_enable_stmt_summary = OFF", "SET GLOBAL tidb_enable_stmt_summary = ON"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return a, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	rows, err := conn.QueryContext(ctx, statement)
	if err != nil {
		return a, fmt.Errorf("run: %w", err)
	}
	for rows.Next() {
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return a, fmt.Errorf("run: %w", err)
	}
	var tables sql.NullString
	err = conn.QueryRowContext(ctx, `
		SELECT TABLE_NAMES FROM information_schema.STATEMENTS_SUMMARY
		WHERE QUERY_SAMPLE_TEXT = ? AND SCHEMA_NAME = ?`, statement, database).Scan(&tables)
	if err != nil {
		return a, fmt.Errorf("read the statement summary: %w", err)
	}
	a.Tables = tables.String
	return a, nil
}
