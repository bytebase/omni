package catalog

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmariadb "github.com/testcontainers/testcontainers-go/modules/mariadb"
)

// TestDropPeriodIfExistsOracle runs ALTER TABLE ... DROP PERIOD IF EXISTS on a
// live mariadb:11.8.8 and on the catalog, from the same table, and asserts
// both give the same error and leave the same periods. The expected code pins
// the engine, so an image upgrade that changes a result fails here.
func TestDropPeriodIfExistsOracle(t *testing.T) {
	const (
		app      = "CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))"
		plain    = "CREATE TABLE ta (id INT, s DATE, e DATE)"
		implicit = "CREATE TABLE ta (id INT) WITH SYSTEM VERSIONING"
		explicit = "CREATE TABLE ta (id INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING"
		both     = "CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e)) WITH SYSTEM VERSIONING"
		overlaps = "CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e), UNIQUE KEY k (id, app WITHOUT OVERLAPS))"
	)
	cases := []struct {
		name, create, alter string
		code                int
	}{
		{"drops the period", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app", 0},
		{"name is case-insensitive", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR APP", 0},
		{"quoted name", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR `app`", 0},
		{"other name", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR nope", 0},
		{"other name without IF EXISTS", app, "ALTER TABLE ta DROP PERIOD FOR nope", 1091},
		{"no period", plain, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app", 0},
		{"no period without IF EXISTS", plain, "ALTER TABLE ta DROP PERIOD FOR app", 1091},
		// Quoted, SYSTEM_TIME names an application-time period, never the
		// system-versioning one.
		{"quoted SYSTEM_TIME, application-time period", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", 0},
		{"quoted SYSTEM_TIME, implicit system period", implicit, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", 0},
		{"quoted SYSTEM_TIME, explicit system period", explicit, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", 0},
		{"quoted SYSTEM_TIME, no period", plain, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", 0},
		{"bitemporal keeps versioning", both, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app", 0},
		// IF EXISTS skips only an absent period, not a period a key uses.
		{"WITHOUT OVERLAPS key", overlaps, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app", 4156},
		{"WITHOUT OVERLAPS key dropped too", overlaps, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP KEY k", 0},
		{"other name beside WITHOUT OVERLAPS key", overlaps, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR nope", 0},
		// Commands apply in order.
		{"then without IF EXISTS", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP PERIOD FOR app", 1091},
		{"after without IF EXISTS", app, "ALTER TABLE ta DROP PERIOD FOR app, DROP PERIOD IF EXISTS FOR app", 0},
		{"twice", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP PERIOD IF EXISTS FOR app", 0},
		{"then ADD PERIOD", app, "ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, ADD PERIOD FOR p2(s, e)", 0},
	}

	conn := startDropPeriodMariaDB(t)
	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, stmt := range []string{"DROP DATABASE IF EXISTS test", "CREATE DATABASE test", "USE test", tc.create} {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("MariaDB %q: %v", stmt, err)
				}
			}
			wantCode, wantMsg := 0, ""
			if _, err := conn.ExecContext(ctx, tc.alter); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatalf("MariaDB %q: %v", tc.alter, err)
				}
				wantCode, wantMsg = int(myErr.Number), myErr.Message
			}
			if wantCode != tc.code {
				t.Fatalf("MariaDB code %d, want %d: %s", wantCode, tc.code, wantMsg)
			}
			var table, engineDDL string
			if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE ta").Scan(&table, &engineDDL); err != nil {
				t.Fatalf("MariaDB SHOW CREATE TABLE: %v", err)
			}

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			if r, err := c.Exec(tc.create, nil); err != nil || r[0].Error != nil {
				t.Fatalf("catalog %q: %v %v", tc.create, err, r)
			}
			r, err := c.Exec(tc.alter, &ExecOptions{ContinueOnError: true})
			if err != nil {
				t.Fatalf("catalog %q: %v", tc.alter, err)
			}
			gotCode, gotMsg := 0, ""
			if r[0].Error != nil {
				var catErr *Error
				if !errors.As(r[0].Error, &catErr) {
					t.Fatalf("catalog %q: %v", tc.alter, r[0].Error)
				}
				gotCode, gotMsg = catErr.Code, catErr.Message
			}
			if gotCode != wantCode || gotMsg != wantMsg {
				t.Errorf("catalog error %d %q, MariaDB %d %q", gotCode, gotMsg, wantCode, wantMsg)
			}
			if got, want := periodState(c.ShowCreateTable("test", "ta")), periodState(engineDDL); got != want {
				t.Errorf("catalog periods:\n%s\nMariaDB periods:\n%s", got, want)
			}
		})
	}
}

// periodState is the part of SHOW CREATE TABLE a DROP PERIOD changes: the
// PERIOD FOR clauses, the WITHOUT OVERLAPS keys, and WITH SYSTEM VERSIONING.
// The rest of the text differs between the catalog and MariaDB 11.8 for
// reasons unrelated to periods (int display width, default collation).
func periodState(ddl string) string {
	var lines []string
	for _, line := range strings.Split(ddl, "\n") {
		line = strings.TrimSuffix(strings.TrimSpace(line), ",")
		if strings.Contains(line, "PERIOD FOR") || strings.Contains(line, "WITHOUT OVERLAPS") {
			lines = append(lines, line)
		}
	}
	if strings.Contains(ddl, "WITH SYSTEM VERSIONING") {
		lines = append(lines, "WITH SYSTEM VERSIONING")
	}
	return strings.Join(lines, "\n")
}

// startDropPeriodMariaDB starts a mariadb:11.8.8 container for this test and
// returns one connection to it, so USE holds across statements.
func startDropPeriodMariaDB(t *testing.T) *sql.Conn {
	t.Helper()
	ctx := context.Background()
	container, err := tcmariadb.Run(ctx, "mariadb:11.8.8",
		tcmariadb.WithDatabase("test"),
		tcmariadb.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start MariaDB container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("MariaDB connection string: %v", err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open MariaDB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("connect to MariaDB: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}
