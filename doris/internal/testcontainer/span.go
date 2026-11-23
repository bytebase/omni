package testcontainer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned Doris what a query reads, for doris/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on a server emptied of every database
//     the image did not start with, and answers the schema Doris stored:
//     each database's tables and views with their columns in order, and each
//     view's definition as information_schema holds it.
//   - Span runs a query in a session's database and answers the result set's
//     column definitions and the tables and views its analyzed plan
//     (EXPLAIN ANALYZED PLAN) reads, or the error Doris refused it with.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned Doris in
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
		cfg := config(t, pin)
		// The fixture drops every database the image did not start with,
		// the judge's own among them, so no connection selects one.
		cfg.DBName = ""
		cfg.MultiStatements = false
		l := &spanLive{cfg: cfg, root: sql.OpenDB(must(mysql.NewConnector(cfg)))}
		l.root.SetMaxOpenConns(1)
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, "default", live)}
}

// SpanFixture is the schema Doris stored for a fixture.
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

// SpanAnswer is what Doris states about one query.
type SpanAnswer struct {
	// Error is the message of the error Doris refused the query with; the
	// rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column definitions, in order.
	Columns []SpanColumn `json:"columns,omitempty"`
	// Relations is every table and view the analyzed plan reads, as
	// database.name, in order.
	Relations []string `json:"relations,omitempty"`
}

// SpanColumn is one column definition of a result set.
type SpanColumn struct {
	Name     string `json:"name"`
	Schema   string `json:"schema,omitempty"`
	Table    string `json:"table,omitempty"`
	OrgTable string `json:"orgTable,omitempty"`
	OrgName  string `json:"orgName,omitempty"`
}

// Fixture builds the starting schema sql, one statement per line, on a
// server emptied of every database it was not started with, and returns
// what Doris stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span runs statement with database selected and returns what Doris states
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
	cfg  *mysql.Config
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
		db, table, kind, def, name string
		pos                        int
	}
	rows, err := l.root.QueryContext(ctx, `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), c.COLUMN_NAME, c.ORDINAL_POSITION
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON v.TABLE_SCHEMA = t.TABLE_SCHEMA AND v.TABLE_NAME = t.TABLE_NAME
		WHERE t.TABLE_SCHEMA NOT IN ('mysql', 'information_schema', '__internal_schema')`)
	if err != nil {
		return fx, err
	}
	defer rows.Close()
	var columns []column
	for rows.Next() {
		var c column
		if err := rows.Scan(&c.db, &c.table, &c.kind, &c.def, &c.name, &c.pos); err != nil {
			return fx, err
		}
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		return fx, err
	}
	// Doris answers in no fixed order.
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
		list := &d.Tables
		if c.kind == "VIEW" {
			list = &d.Views
		}
		if n := len(*list); n == 0 || (*list)[n-1].Name != c.table {
			*list = append(*list, SpanRelation{Name: c.table, Definition: c.def})
		}
		r := &(*list)[len(*list)-1]
		r.Columns = append(r.Columns, c.name)
	}
	return fx, nil
}

// The relations an analyzed plan reads: an OLAP table's scan names it, and
// a view, an information_schema table or any other catalog table is the
// three-part qualifier of the alias above its node.
var (
	olapScanRe = regexp.MustCompile(`LogicalOlapScan \( qualified=([^,\s]+)`)
	relationRe = regexp.MustCompile(`LogicalSubQueryAlias \( qualifier=\[([^,\]]+), ([^,\]]+), ([^,\]]+)\] \)`)
)

func (l *spanLive) span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	if l.wire == nil {
		w, err := dialWire(l.cfg.Addr, l.cfg.User)
		if err != nil {
			return a, err
		}
		l.wire = w
	}
	w := l.wire
	if err := w.Use(database); err != nil {
		return a, fmt.Errorf("select %s: %w", database, err)
	}
	cols, _, err := w.Query(statement)
	var refused *wireError
	if errors.As(err, &refused) {
		a.Error = refused.Message
		return a, nil
	}
	if err != nil {
		return a, err
	}
	for _, c := range cols {
		a.Columns = append(a.Columns, SpanColumn{Name: c.Name, Schema: c.Schema, Table: c.Table, OrgTable: c.OrgTable, OrgName: c.OrgName})
	}
	_, plan, err := w.Query("EXPLAIN ANALYZED PLAN " + statement)
	if err != nil {
		return a, fmt.Errorf("EXPLAIN: %w", err)
	}
	seen := map[string]bool{}
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			a.Relations = append(a.Relations, name)
		}
	}
	for _, row := range plan {
		line := strings.Join(row, " ")
		if m := olapScanRe.FindStringSubmatch(line); m != nil {
			parts := strings.Split(m[1], ".")
			if len(parts) != 3 {
				return a, fmt.Errorf("a scan of %q", m[1])
			}
			add(parts[1] + "." + parts[2])
		}
		if m := relationRe.FindStringSubmatch(line); m != nil {
			add(m[2] + "." + m[3])
		}
	}
	return a, nil
}
