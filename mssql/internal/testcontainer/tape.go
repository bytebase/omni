package testcontainer

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// Open returns t's connection to SQL Server: a replay of its tape, or the
// pinned server in a recording. Each call is a stream of its own, one session
// that keeps its SET options and database across requests, starting in the
// database of t's top-level test (Database). In a recording each top-level
// test starts on a server that holds no database but the system's: before
// its first connection the others are dropped, and its own is created empty.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, server, Default)
}

// Database is the name of the database Open connects t to: one per
// top-level test, named for it, so that a test reads the same DB_NAME() in
// every recording. SQL Server keeps 128 characters of a name, so a long test
// name ends in a hash of the whole.
func Database(t testing.TB) string {
	top := enginetest.TopName(t.Name())
	name := "omni_" + notNameChars.ReplaceAllString(top, "_")
	if len(name) > 100 {
		sum := sha256.Sum256([]byte(top))
		name = name[:91] + "_" + hex.EncodeToString(sum[:4])
	}
	return name
}

var notNameChars = regexp.MustCompile(`[^A-Za-z0-9_]`)

// live connects to the database of t's top-level test on the pinned server.
// On the top-level test's first connection it drops every database but the
// system's and creates the test's own.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	in := started(t, pin)
	top := enginetest.TopName(t.Name())
	name := Database(t)
	in.mu.Lock()
	defer in.mu.Unlock()
	if !in.emptied[top] {
		ctx := context.Background()
		db := sql.OpenDB(in.connector)
		defer db.Close()
		var stmts []string
		for _, d := range userDatabases(t, db) {
			stmts = append(stmts, "ALTER DATABASE "+quoteName(d)+" SET SINGLE_USER WITH ROLLBACK IMMEDIATE", "DROP DATABASE "+quoteName(d))
		}
		stmts = append(stmts, "CREATE DATABASE "+quoteName(name))
		for _, stmt := range stmts {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
		in.emptied[top] = true
	}
	u, err := url.Parse(in.dsn)
	if err != nil {
		t.Fatalf("SQL Server's DSN: %v", err)
	}
	q := u.Query()
	q.Set("database", name)
	u.RawQuery = q.Encode()
	c, err := mssql.NewConnector(u.String())
	if err != nil {
		t.Fatalf("connect to SQL Server: %v", err)
	}
	return connector{c}
}

// userDatabases lists the databases that are not the system's.
func userDatabases(t testing.TB, db *sql.DB) []string {
	t.Helper()
	rows, err := db.QueryContext(context.Background(), "SELECT name FROM sys.databases WHERE database_id > 4 ORDER BY name")
	if err != nil {
		t.Fatalf("list SQL Server's databases: %v", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatalf("list SQL Server's databases: %v", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("list SQL Server's databases: %v", err)
	}
	return out
}

// connector is go-mssqldb's, with what the tape needs of a connection.
type connector struct{ c *mssql.Connector }

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	dc, err := c.c.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return conn{dc}, nil
}

func (c connector) Driver() driver.Driver { return c.c.Driver() }

// conn runs a request as go-mssqldb's own statement does: go-mssqldb has no
// ExecContext or QueryContext of its own, and its Prepare prepares nothing
// on the server, so a statement with no arguments is one batch and one with
// arguments a call of sp_executesql.
type conn struct{ driver.Conn }

func (c conn) prepare(ctx context.Context, query string) (driver.Stmt, error) {
	return c.Conn.(driver.ConnPrepareContext).PrepareContext(ctx, query)
}

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	st, err := c.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	return st.(driver.StmtExecContext).ExecContext(ctx, args)
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	st, err := c.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	r, err := st.(driver.StmtQueryContext).QueryContext(ctx, args)
	if err != nil {
		st.Close()
		return nil, err
	}
	return &rows{Rows: r, st: st}, nil
}

func (c conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

func (c conn) Ping(ctx context.Context) error { return c.Conn.(driver.Pinger).Ping(ctx) }

func (c conn) CheckNamedValue(nv *driver.NamedValue) error {
	return c.Conn.(driver.NamedValueChecker).CheckNamedValue(nv)
}

// ResetSession keeps the session: go-mssqldb would send the next batch with
// RESETCONNECTION, which clears SET options and the database context, and a
// tape's stream is one session.
func (c conn) ResetSession(context.Context) error { return nil }

// rows closes its statement with it, and passes on what the tape reads.
type rows struct {
	driver.Rows
	st driver.Stmt
}

func (r *rows) Close() error {
	err := r.Rows.Close()
	r.st.Close()
	return err
}

func (r *rows) ColumnTypeDatabaseTypeName(i int) string {
	if n, ok := r.Rows.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return n.ColumnTypeDatabaseTypeName(i)
	}
	return ""
}

func (r *rows) HasNextResultSet() bool {
	if n, ok := r.Rows.(driver.RowsNextResultSet); ok {
		return n.HasNextResultSet()
	}
	return false
}

func (r *rows) NextResultSet() error {
	if n, ok := r.Rows.(driver.RowsNextResultSet); ok {
		return n.NextResultSet()
	}
	return errors.New("no next result set")
}

// Codec renders go-mssqldb's mssql.Error as every error SQL Server reported
// for the request, first to last, one to a line: its number, state, class
// and line, then its procedure where it names one, and its message.
//
//	102 1 15 1 Incorrect syntax near 'FROM'.
//	2627 1 14 3 proc=p Violation of PRIMARY KEY constraint ...
//
// A line break in a message is written \n, and a backslash \\. The server's
// name is left out: it is the container's, and differs on every run.
type Codec struct{}

var (
	escaper   = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\r", `\r`)
	unescaper = strings.NewReplacer(`\\`, `\`, `\n`, "\n", `\r`, "\r")
)

func (Codec) EncodeError(err error) (string, bool) {
	var e mssql.Error
	if !errors.As(err, &e) {
		return "", false
	}
	all := e.All
	if len(all) == 0 {
		all = []mssql.Error{e}
	}
	lines := make([]string, len(all))
	for i, x := range all {
		l := fmt.Sprintf("%d %d %d %d ", x.Number, x.State, x.Class, x.LineNo)
		if x.ProcName != "" {
			l += "proc=" + escaper.Replace(strings.ReplaceAll(x.ProcName, " ", `\s`)) + " "
		}
		lines[i] = l + escaper.Replace(x.Message)
	}
	return strings.Join(lines, "\n"), true
}

func (Codec) DecodeError(text string) error {
	var all []mssql.Error
	for _, line := range strings.Split(text, "\n") {
		f := strings.SplitN(line, " ", 5)
		if len(f) < 5 {
			return fmt.Errorf("tape: a SQL Server error %q", line)
		}
		var nums [4]int64
		for i := range nums {
			n, err := strconv.ParseInt(f[i], 10, 32)
			if err != nil {
				return fmt.Errorf("tape: a SQL Server error %q", line)
			}
			nums[i] = n
		}
		x := mssql.Error{Number: int32(nums[0]), State: uint8(nums[1]), Class: uint8(nums[2]), LineNo: int32(nums[3])}
		msg := f[4]
		if p, rest, ok := strings.Cut(msg, " "); ok && strings.HasPrefix(p, "proc=") {
			x.ProcName = strings.ReplaceAll(unescaper.Replace(strings.TrimPrefix(p, "proc=")), `\s`, " ")
			msg = rest
		}
		x.Message = unescaper.Replace(msg)
		all = append(all, x)
	}
	e := all[len(all)-1]
	e.All = all
	return e
}

// SQL Server makes up two kinds of name, each different on every run, so
// Codec is a tape.Renamer and gives each a stand-in:
//
//   - An object's id, negative for a temporary table, which a test reads
//     as text, 'oid:' + CAST(object_id AS varchar(20)), and sends back as
//     an argument it converts, CAST(SUBSTRING(@p1, 5, 20) AS int). Its
//     stand-in is oid:0 and a number, which no id's text is.
//   - The name of a constraint declared without one, PK__t__3213E83F2F1A1B2C
//     or DF__t__a__5AEE82B9: its kind, its table and column cut short (a
//     space and all), and hex digits of its object's id. Its stand-in keeps all but the digits,
//     which become standin and a number.
var madeUp = regexp.MustCompile(`\boid:-?[0-9]+\b|\b(?:PK|UQ|DF|CK|FK)__[^'"\[\]\n]{1,40}?__(?:[0-9A-F]{8}(?:[0-9A-F]{8})?|standin[0-9]+)\b`)

// MadeUp matches an object id read as text, and a constraint name SQL
// Server made up.
func (Codec) MadeUp() *regexp.Regexp { return madeUp }

// StandIn is the nth made-up name's stand-in.
func (Codec) StandIn(name string, n int) string {
	if strings.HasPrefix(name, "oid:") {
		return fmt.Sprintf("oid:0%d", n)
	}
	i := strings.LastIndex(name, "__")
	return fmt.Sprintf("%s__standin%d", name[:i], n)
}
