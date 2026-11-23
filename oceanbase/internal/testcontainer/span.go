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
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned OceanBase what a query reads, for
// oceanbase/analysis's span.tsv. Its tape holds one raw exchange per
// question, and each answer is what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on an emptied server and answers the
//     schema OceanBase stored: each database's tables and views with their
//     columns in order, and each view's definition as information_schema
//     holds it.
//   - Span prepares a query in a session's database and answers the result
//     set's column definitions, and the tables and views OceanBase checks
//     the SELECT privilege on: an account that holds none runs the query
//     again and again, and is granted each table a refusal (1142) names,
//     until the query runs. OceanBase checks privileges when a statement
//     runs, not when it is prepared.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned OceanBase
// in a recording.
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
		dsn, err := OceanBase.DSN(context.Background())
		if err != nil {
			t.Fatalf("OceanBase: %v", err)
		}
		cfg, err := gomysql.ParseDSN(dsn)
		if err != nil {
			t.Fatalf("%v", err)
		}
		l := &spanLive{db: db, addr: cfg.Addr, user: cfg.User}
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, "default", live)}
}

// SpanFixture is the schema OceanBase stored for a fixture.
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

// SpanAnswer is what OceanBase states about one query.
type SpanAnswer struct {
	// Error is the error number of a query OceanBase refused; the rest is
	// empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column definitions, in order.
	Columns []SpanColumn `json:"columns,omitempty"`
	// Tables is every table and view OceanBase checks the SELECT privilege
	// on for the query, as db.table.
	Tables []string `json:"tables,omitempty"`
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
// database and account it was not started with, and returns what OceanBase
// stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span prepares statement with database selected and returns what OceanBase
// states it reads.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + database + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

// spanUser is the account the privilege walk prepares queries as. It holds
// SHOW VIEW on each fixture database, which lets it select the database and
// read nothing.
const spanUser = "omni_span"

type spanLive struct {
	db         *sql.DB
	addr, user string
	root       *wire
}

func (l *spanLive) close() {
	if l.root != nil {
		l.root.Close()
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
	if err := empty(ctx, l.db); err != nil {
		return fx, err
	}
	if _, err := l.db.ExecContext(ctx, sqlText); err != nil {
		return fx, fmt.Errorf("build the fixture: %w", err)
	}
	rows, err := l.db.QueryContext(ctx, `
		SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE, COALESCE(v.VIEW_DEFINITION, ''), c.COLUMN_NAME
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		LEFT JOIN information_schema.VIEWS v ON v.TABLE_SCHEMA = t.TABLE_SCHEMA AND v.TABLE_NAME = t.TABLE_NAME
		WHERE LOWER(t.TABLE_SCHEMA) NOT IN ('mysql', 'information_schema', 'oceanbase', 'sys', 'test')
		  AND t.TABLE_SCHEMA NOT LIKE '\_\_%'
		ORDER BY BINARY t.TABLE_SCHEMA, BINARY t.TABLE_NAME, c.ORDINAL_POSITION`)
	if err != nil {
		return fx, err
	}
	for rows.Next() {
		var db, table, kind, def, column string
		if err := rows.Scan(&db, &table, &kind, &def, &column); err != nil {
			rows.Close()
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
	rows.Close()
	if err := rows.Err(); err != nil {
		return fx, err
	}
	stmts := []string{"CREATE USER " + spanUser}
	for _, d := range fx.Databases {
		stmts = append(stmts, "GRANT SHOW VIEW ON `"+d.Name+"`.* TO "+spanUser)
	}
	for _, stmt := range stmts {
		if _, err := l.db.ExecContext(ctx, stmt); err != nil {
			return fx, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return fx, nil
}

// deniedTable reads the table of OceanBase's 1142: "SELECT command denied to
// user 'u'@'%' for table 't'".
var deniedTable = regexp.MustCompile(`for table '([^']+)'$`)

func (l *spanLive) span(ctx context.Context, database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	if l.root == nil {
		w, err := dialWire(l.addr, l.user, "")
		if err != nil {
			return a, err
		}
		l.root = w
	}
	if err := l.root.Use(database); err != nil {
		return a, fmt.Errorf("select %s: %w", database, err)
	}
	cols, err := l.root.Prepare(statement)
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

	// The privilege walk.
	tenant := ""
	if i := strings.IndexByte(l.user, '@'); i >= 0 {
		tenant = l.user[i:]
	}
	granted := map[string]bool{}
	defer func() {
		for t := range granted {
			db, table, _ := strings.Cut(t, ".")
			_, _ = l.db.ExecContext(ctx, "REVOKE SELECT ON `"+db+"`.`"+table+"` FROM "+spanUser)
		}
	}()
	for attempt := 0; ; attempt++ {
		if attempt > 32 {
			return a, errors.New("the privilege walk does not end")
		}
		denied, err := l.denied(database, statement, spanUser+tenant)
		if err != nil {
			return a, err
		}
		if denied == "" {
			break
		}
		var dbs []string
		rows, err := l.db.QueryContext(ctx, `
			SELECT TABLE_SCHEMA FROM information_schema.TABLES
			WHERE LOWER(TABLE_NAME) = LOWER(?) AND LOWER(TABLE_SCHEMA) NOT IN ('mysql', 'information_schema', 'oceanbase', 'sys', 'test')`, denied)
		if err != nil {
			return a, err
		}
		for rows.Next() {
			var db string
			if err := rows.Scan(&db); err != nil {
				rows.Close()
				return a, err
			}
			dbs = append(dbs, db)
		}
		rows.Close()
		if len(dbs) != 1 {
			return a, fmt.Errorf("the refusal names %s, which %d fixture databases hold", denied, len(dbs))
		}
		key := dbs[0] + "." + denied
		if granted[key] {
			// A grant reaches a new session at once; give it a moment
			// before taking the refusal for a loop.
			time.Sleep(time.Second)
			continue
		}
		granted[key] = true
		a.Tables = append(a.Tables, key)
		if _, err := l.db.ExecContext(ctx, "GRANT SELECT ON `"+dbs[0]+"`.`"+denied+"` TO "+spanUser); err != nil {
			return a, err
		}
	}
	slices.Sort(a.Tables)
	return a, nil
}

// denied runs statement as user on a session of its own and returns the
// table OceanBase refuses it for, or "" when the query runs. The fixture's
// tables are small, so it runs whole.
func (l *spanLive) denied(database, statement, user string) (string, error) {
	w, err := dialWire(l.addr, user, "")
	if err != nil {
		return "", err
	}
	defer w.Close()
	if err := w.Use(database); err != nil {
		return "", fmt.Errorf("select %s as %s: %w", database, user, err)
	}
	_, err = w.Query(statement)
	var refused *wireError
	switch {
	case err == nil:
		return "", nil
	case errors.As(err, &refused) && refused.Number == 1142:
		m := deniedTable.FindStringSubmatch(refused.Message)
		if m == nil {
			return "", fmt.Errorf("a refusal without a table: %v", refused)
		}
		return m[1], nil
	default:
		return "", fmt.Errorf("run as %s: %w", user, err)
	}
}
