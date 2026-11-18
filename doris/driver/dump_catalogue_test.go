package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"fmt"
	"io"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// A Doris kernel built without the PL/SQL grammar (SelectDB 4.0.12+, 4.1.8+ and
// 5.0.5+, apache/doris#58700) answers a routine catalogue probe with an error
// instead of an empty set, and so does an account without the catalogue's
// privilege. No public image of such a kernel exists to record from, so the
// dump's answer to that refusal is held here against a scripted connection:
// the objects it can state, and none of the routines.
func TestDumpToleratesRefusedRoutineCatalogues(t *testing.T) {
	refusal := &gomysql.MySQLError{Number: 1105, SQLState: [5]byte{'H', 'Y', '0', '0', '0'}, Message: "errCode = 2, detailMessage = no viable alternative at input 'SHOW FULL FUNCTIONS'"}
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		switch {
		case strings.HasPrefix(query, "SHOW FULL FUNCTIONS"), strings.Contains(query, "information_schema.routines"):
			return nil, nil, refusal
		default:
			return dumpCatalogueAnswer(query)
		}
	})
	d := &Driver{db: db, databaseName: "omni_stripped"}

	var out strings.Builder
	if err := d.Dump(context.Background(), &out, nil); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	dump := out.String()
	if !strings.Contains(dump, "CREATE TABLE `orders`") {
		t.Errorf("the table is not in the dump:\n%s", dump)
	}
	for _, absent := range []string{"-- Function structure", "-- Procedure structure", "CREATE PROCEDURE"} {
		if strings.Contains(dump, absent) {
			t.Errorf("the dump carries %q although the server states no routines:\n%s", absent, dump)
		}
	}
}

// TestDumpFailsOnARoutineProbeTheServerDidNotAnswer pins the other side of the
// rule: only the server's own refusal is tolerated, so a connection that dies or
// a context that ends still fails the dump rather than writing a schema that
// silently lost its routines.
func TestDumpFailsOnARoutineProbeTheServerDidNotAnswer(t *testing.T) {
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		if strings.HasPrefix(query, "SHOW FULL FUNCTIONS") {
			return nil, nil, errors.New("driver: bad connection")
		}
		return dumpCatalogueAnswer(query)
	})
	d := &Driver{db: db, databaseName: "omni_stripped"}

	var out strings.Builder
	if err := d.Dump(context.Background(), &out, nil); err == nil {
		t.Fatalf("Dump: want an error, got a dump:\n%s", out.String())
	}
}

// dumpCatalogueAnswer answers the dump's other queries for one database with one
// table, so a routine probe is the only thing a test can break.
func dumpCatalogueAnswer(query string) ([]string, [][]sqldriver.Value, error) {
	switch {
	case strings.HasPrefix(query, "SELECT SCHEMA_NAME FROM information_schema.schemata"):
		return []string{"SCHEMA_NAME"}, [][]sqldriver.Value{{"omni_stripped"}}, nil
	case strings.HasPrefix(query, "SELECT current_catalog()"):
		return []string{"current_catalog()"}, [][]sqldriver.Value{{"internal"}}, nil
	case strings.Contains(query, "mv_infos("):
		return []string{"Name"}, nil, nil
	case strings.Contains(query, "information_schema.tables"):
		return []string{"TABLE_NAME", "TABLE_TYPE"}, [][]sqldriver.Value{{"orders", "BASE TABLE"}}, nil
	case strings.HasPrefix(query, "DESC "):
		return []string{"IndexName", "IndexKeysType"}, [][]sqldriver.Value{{"orders", "DUPLICATE_KEYS"}}, nil
	case strings.HasPrefix(query, "SHOW CREATE TABLE"):
		return []string{"Table", "Create Table"},
			[][]sqldriver.Value{{"orders", "CREATE TABLE `orders` (\n  `id` int NULL\n) ENGINE=OLAP\nDUPLICATE KEY(`id`)"}}, nil
	default:
		return nil, nil, fmt.Errorf("unexpected statement: %q", query)
	}
}

// stubDatabase returns a database/sql handle whose every query is answered by
// answer, so a test reaches a probe's failure without a server.
func stubDatabase(t *testing.T, answer func(query string) ([]string, [][]sqldriver.Value, error)) *sql.DB {
	t.Helper()
	db := sql.OpenDB(stubConnector{answer: answer})
	t.Cleanup(func() { _ = db.Close() })
	return db
}

type stubConnector struct {
	answer func(query string) ([]string, [][]sqldriver.Value, error)
}

func (c stubConnector) Connect(context.Context) (sqldriver.Conn, error) {
	return &stubConn{answer: c.answer}, nil
}

func (stubConnector) Driver() sqldriver.Driver { return stubDriver{} }

type stubDriver struct{}

func (stubDriver) Open(string) (sqldriver.Conn, error) {
	return nil, errors.New("stubDriver.Open is not used: the connector builds its own connection")
}

type stubConn struct {
	answer func(query string) ([]string, [][]sqldriver.Value, error)
}

func (*stubConn) Prepare(string) (sqldriver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}
func (*stubConn) Close() error { return nil }
func (*stubConn) Begin() (sqldriver.Tx, error) {
	return nil, errors.New("transactions are not supported")
}

// QueryContext answers both Query and QueryRow without Prepare, and renders the
// arguments into the statement the way the driver's DSN does (interpolateParams).
func (c *stubConn) QueryContext(_ context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Rows, error) {
	for _, arg := range args {
		query = strings.Replace(query, "?", fmt.Sprintf("'%v'", arg.Value), 1)
	}
	columns, rows, err := c.answer(query)
	if err != nil {
		return nil, err
	}
	return &stubRows{columns: columns, rows: rows}, nil
}

type stubRows struct {
	columns []string
	rows    [][]sqldriver.Value
	next    int
}

func (r *stubRows) Columns() []string { return r.columns }
func (*stubRows) Close() error        { return nil }
func (r *stubRows) Next(dest []sqldriver.Value) error {
	if r.next >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.next])
	r.next++
	return nil
}
