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

// mysqlError is the answer the client library renders to *mysql.MySQLError. The
// numbers and messages are the ones the pinned 3.1.4 server answers with, and
// the ones its kernel without the PL/SQL grammar states, except where a test
// says otherwise.
func mysqlError(number uint16, sqlState, message string) *gomysql.MySQLError {
	err := &gomysql.MySQLError{Number: number, Message: message}
	copy(err.SQLState[:], sqlState)
	return err
}

// TestRoutineProbeRefusals holds the predicate that reads a routine probe's
// failure as an absent catalogue: only a refusal the server states about the
// probe itself counts, so a failure to answer the query still fails the dump
// rather than dropping every function and procedure from it.
func TestRoutineProbeRefusals(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name     string
		ctx      context.Context
		relation string
		err      error
		want     bool
	}{
		{
			name:     "the account has no privilege on the database",
			ctx:      context.Background(),
			relation: "routines",
			err:      mysqlError(1044, "42000", "errCode = 2, detailMessage = Access denied for user 'sync' to database 'omni_dump'"),
			want:     true,
		},
		{
			name:     "the account has no privilege on the object",
			ctx:      context.Background(),
			relation: "routines",
			err:      mysqlError(1142, "42000", "errCode = 2, detailMessage = SHOW CREATE TABLE command denied to user 'sync'@'10.0.0.1' for table 'omni_dump: proc'"),
			want:     true,
		},
		{
			name:     "the account has no privilege on the operation",
			ctx:      context.Background(),
			relation: "routines",
			err:      mysqlError(1227, "42000", "errCode = 2, detailMessage = Access denied; you need (at least one of) the ADMIN privilege(s) for this operation"),
			want:     true,
		},
		{
			name: "the parser refuses the statement",
			ctx:  context.Background(),
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = \nno viable alternative at input 'SHOW FULL FUNCTIONS'(line 1, pos 24)\n"),
			want: true,
		},
		{
			name:     "the catalogue has no such relation",
			ctx:      context.Background(),
			relation: "routines",
			err:      mysqlError(1105, "HY000", "errCode = 2, detailMessage = Table [routines] does not exist in database [information_schema]."),
			want:     true,
		},
		{
			name: "an internal error is not a refusal",
			ctx:  context.Background(),
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = No backend available as scan node, please check the status of your backends."),
			want: false,
		},
		{
			name: "an overload is not a refusal",
			ctx:  context.Background(),
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = Failed to wait for the available backend, timeout"),
			want: false,
		},
		{
			name: "a query timeout is not a refusal",
			ctx:  context.Background(),
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = Query timeout exceeded, 30000ms"),
			want: false,
		},
		{
			name:     "another missing relation says nothing about the catalogue",
			ctx:      context.Background(),
			relation: "routines",
			err:      mysqlError(1105, "HY000", "errCode = 2, detailMessage = Table [alias_abs] does not exist in database [omni_dump]."),
			want:     false,
		},
		{
			name:     "a probe that reads no relation cannot blame one",
			ctx:      context.Background(),
			relation: "",
			err:      mysqlError(1105, "HY000", "errCode = 2, detailMessage = Table [routines] does not exist in database [information_schema]."),
			want:     false,
		},
		{
			name: "an unknown database is not a refusal",
			ctx:  context.Background(),
			err:  mysqlError(1049, "42000", "errCode = 2, detailMessage = Unknown database 'omni_dump'"),
			want: false,
		},
		{
			name: "a dead connection is not a refusal",
			ctx:  context.Background(),
			err:  errors.New("driver: bad connection"),
			want: false,
		},
		{
			name:     "an ended context is not a refusal",
			ctx:      canceled,
			relation: "routines",
			err:      mysqlError(1044, "42000", "errCode = 2, detailMessage = Access denied for user 'sync' to database 'omni_dump'"),
			want:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := routineProbeRefused(test.ctx, test.relation, test.err); got != test.want {
				t.Errorf("routineProbeRefused(%v) = %v, want %v", test.err, got, test.want)
			}
		})
	}
}

// A Doris kernel built without the PL/SQL grammar (SelectDB 4.0.12+, 4.1.8+ and
// 5.0.5+, apache/doris#58700) answers a routine catalogue probe with the
// parser's refusal instead of an empty set. No public image of such a kernel
// exists to record from, so the dump's answer to that refusal is held here
// against a scripted connection: the objects it can state, and none of the
// routines.
func TestDumpToleratesRefusedRoutineCatalogues(t *testing.T) {
	refusal := mysqlError(1105, "HY000", "errCode = 2, detailMessage = \nno viable alternative at input 'SHOW FULL FUNCTIONS'(line 1, pos 42)\n")
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		switch {
		case isRoutineProbe(query):
			return nil, nil, refusal
		default:
			return dumpCatalogueAnswer(query, dumpTableAnswer{name: "orders", kind: "BASE TABLE"})
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
// rule: only the refusals above are tolerated, so a probe the server failed to
// answer still fails the dump rather than writing a schema that silently lost
// its routines.
func TestDumpFailsOnARoutineProbeTheServerDidNotAnswer(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "an internal error",
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = No backend available as scan node, please check the status of your backends."),
		},
		{
			name: "a query timeout",
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = Query timeout exceeded, 30000ms"),
		},
		{
			name: "another missing table",
			err:  mysqlError(1105, "HY000", "errCode = 2, detailMessage = Table [routines] does not exist in database [omni_dump]."),
		},
		{
			name: "the database is gone",
			err:  mysqlError(1049, "42000", "errCode = 2, detailMessage = Unknown database 'omni_stripped'"),
		},
		{
			name: "the connection died",
			err:  errors.New("driver: bad connection"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
				if isRoutineProbe(query) {
					return nil, nil, test.err
				}
				return dumpCatalogueAnswer(query, dumpTableAnswer{name: "orders", kind: "BASE TABLE"})
			})
			d := &Driver{db: db, databaseName: "omni_stripped"}

			var out strings.Builder
			if err := d.Dump(context.Background(), &out, nil); err == nil {
				t.Fatalf("Dump: want an error, got a dump:\n%s", out.String())
			}
		})
	}
}

// TestDumpRefusesAViewThatCallsAnOmittedFunction holds the dump to a replayable
// script: when the server refused the function catalogue, a view whose
// definition calls one of its functions names a function the script never
// creates, so the dump reports that instead of handing the caller a script that
// fails to restore. The corpus's af_function_view calls alias_abs this way.
func TestDumpRefusesAViewThatCallsAnOmittedFunction(t *testing.T) {
	refusal := mysqlError(1105, "HY000", "errCode = 2, detailMessage = \nno viable alternative at input 'SHOW FULL FUNCTIONS'(line 1, pos 42)\n")
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		switch {
		case isRoutineProbe(query):
			return nil, nil, refusal
		default:
			return dumpCatalogueAnswer(query, dumpTableAnswer{name: "orders", kind: "BASE TABLE"},
				dumpTableAnswer{
					name: "af_function_view",
					kind: "VIEW",
					stmt: "CREATE VIEW `af_function_view`\n(k)\n AS SELECT `omni_stripped`.`alias_abs`(`internal`.`omni_stripped`.`orders`.`id`) AS `k` FROM `internal`.`omni_stripped`.`orders`",
				})
		}
	})
	d := &Driver{db: db, databaseName: "omni_stripped"}

	var out strings.Builder
	err := d.Dump(context.Background(), &out, nil)
	if err == nil {
		t.Fatalf("Dump: want an error, got a dump:\n%s", out.String())
	}
	for _, want := range []string{"af_function_view", "alias_abs"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Dump error %q does not name %q", err, want)
		}
	}
}

// TestDumpKeepsAViewThatCallsNoFunction pins the other side: a view that calls
// only a built-in reads nothing the dump omits, so a refused function catalogue
// still leaves a replayable script.
func TestDumpKeepsAViewThatCallsNoFunction(t *testing.T) {
	refusal := mysqlError(1044, "42000", "errCode = 2, detailMessage = Access denied for user 'sync' to database 'omni_stripped'")
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		switch {
		case isRoutineProbe(query):
			return nil, nil, refusal
		default:
			return dumpCatalogueAnswer(query, dumpTableAnswer{name: "orders", kind: "BASE TABLE"},
				dumpTableAnswer{
					name: "abs_view",
					kind: "VIEW",
					stmt: "CREATE VIEW `abs_view`\n(k)\n AS SELECT abs(`internal`.`omni_stripped`.`orders`.`id`) AS `k` FROM `internal`.`omni_stripped`.`orders`",
				})
		}
	})
	d := &Driver{db: db, databaseName: "omni_stripped"}

	var out strings.Builder
	if err := d.Dump(context.Background(), &out, nil); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	dump := out.String()
	for _, want := range []string{"CREATE TABLE `orders`", "CREATE VIEW `abs_view`"} {
		if !strings.Contains(dump, want) {
			t.Errorf("the dump does not carry %q:\n%s", want, dump)
		}
	}
}

// isRoutineProbe reports whether a statement is one of the two probes of the
// routine catalogues.
func isRoutineProbe(query string) bool {
	return strings.HasPrefix(query, "SHOW FULL FUNCTIONS") || strings.Contains(query, "information_schema.routines")
}

// dumpTableAnswer is one object the scripted server states for a database.
type dumpTableAnswer struct {
	name string
	kind string
	stmt string
}

// dumpCatalogueAnswer answers the dump's other queries for one database holding
// the given objects, so a routine probe is the only thing a test can break.
func dumpCatalogueAnswer(query string, objects ...dumpTableAnswer) ([]string, [][]sqldriver.Value, error) {
	switch {
	case strings.HasPrefix(query, "SELECT SCHEMA_NAME FROM information_schema.schemata"):
		return []string{"SCHEMA_NAME"}, [][]sqldriver.Value{{"omni_stripped"}}, nil
	case strings.HasPrefix(query, "SELECT current_catalog()"):
		return []string{"current_catalog()"}, [][]sqldriver.Value{{"internal"}}, nil
	case strings.Contains(query, "mv_infos("):
		return []string{"Name"}, nil, nil
	case strings.Contains(query, "information_schema.tables"):
		rows := make([][]sqldriver.Value, 0, len(objects))
		for _, o := range objects {
			rows = append(rows, []sqldriver.Value{o.name, o.kind})
		}
		return []string{"TABLE_NAME", "TABLE_TYPE"}, rows, nil
	case strings.HasPrefix(query, "DESC "):
		// No synchronous materialized views: the table's own index only.
		for _, o := range objects {
			if o.kind == "BASE TABLE" && strings.Contains(query, "`"+o.name+"`") {
				return []string{"IndexName", "IndexKeysType"}, [][]sqldriver.Value{{o.name, "DUPLICATE_KEYS"}}, nil
			}
		}
		return nil, nil, fmt.Errorf("no such table: %q", query)
	case strings.HasPrefix(query, "SHOW CREATE TABLE"):
		return []string{"Table", "Create Table"},
			[][]sqldriver.Value{{"orders", "CREATE TABLE `orders` (\n  `id` int NULL\n) ENGINE=OLAP\nDUPLICATE KEY(`id`)"}}, nil
	case strings.HasPrefix(query, "SHOW CREATE VIEW"):
		for _, o := range objects {
			if o.stmt != "" && strings.Contains(query, o.name) {
				return []string{"View", "Create View", "character_set_client", "collation_connection"},
					[][]sqldriver.Value{{o.name, o.stmt, "utf8", "utf8_bin"}}, nil
			}
		}
		return nil, nil, fmt.Errorf("no such view: %q", query)
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
