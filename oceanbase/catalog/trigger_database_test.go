package catalog

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestOracle_QualifiedTriggerDatabase: a trigger named with a database is taken or refused
// as OceanBase takes it, with the server's code. The database it names has to exist
// (1049), then its table (1146), and then the two have to be one database (1435). The table
// is looked up in its own database or the current one, never in the trigger's.
func TestOracle_QualifiedTriggerDatabase(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	setup := []string{
		"DROP DATABASE IF EXISTS zz_trig_a",
		"DROP DATABASE IF EXISTS zz_trig_b",
		"CREATE DATABASE zz_trig_a",
		"CREATE DATABASE zz_trig_b",
		"CREATE TABLE zz_trig_a.t (id INT PRIMARY KEY, a INT)",
		"CREATE TABLE zz_trig_a.only_a (id INT PRIMARY KEY, a INT)",
		"CREATE TABLE zz_trig_b.t (id INT PRIMARY KEY, a INT)",
		"USE zz_trig_a",
	}
	for _, stmt := range setup {
		if _, err := o.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	c := New()
	if _, err := c.Exec(strings.Join(setup[2:], "; "), nil); err != nil {
		t.Fatal(err)
	}

	const body = " FOR EACH ROW SET NEW.a = 1"
	for _, probe := range []struct{ id, sql string }{
		{"its own database", "CREATE TRIGGER zz_trig_a.q1 BEFORE INSERT ON t" + body},
		{"both qualified", "CREATE TRIGGER zz_trig_a.q2 BEFORE INSERT ON zz_trig_a.t" + body},
		{"both in another database", "CREATE TRIGGER zz_trig_b.q3 BEFORE INSERT ON zz_trig_b.t" + body},
		{"database in capitals", "CREATE TRIGGER ZZ_TRIG_A.q4 BEFORE INSERT ON t" + body},
		{"if not exists", "CREATE TRIGGER IF NOT EXISTS zz_trig_a.q1 BEFORE INSERT ON t" + body},
		{"table in the current database", "CREATE TRIGGER zz_trig_b.q5 BEFORE INSERT ON t" + body},
		{"table qualified with another", "CREATE TRIGGER zz_trig_b.q6 BEFORE INSERT ON zz_trig_a.t" + body},
		{"table only in the current database", "CREATE TRIGGER zz_trig_b.q7 BEFORE INSERT ON only_a" + body},
		{"no such table", "CREATE TRIGGER zz_trig_a.q8 BEFORE INSERT ON no_such_table" + body},
		{"no such table, another database", "CREATE TRIGGER zz_trig_b.q9 BEFORE INSERT ON no_such_table" + body},
		{"no such database", "CREATE TRIGGER zz_no_such_database.q10 BEFORE INSERT ON t" + body},
		{"no such database or table", "CREATE TRIGGER zz_no_such_database.q11 BEFORE INSERT ON no_such_table" + body},
		{"drop", "DROP TRIGGER zz_trig_b.q3"},
	} {
		t.Run(probe.id, gapped(func(t testing.TB) {
			serverCode := 0
			if _, err := o.db.ExecContext(ctx, probe.sql); err != nil {
				var merr *gomysql.MySQLError
				if !errors.As(err, &merr) {
					t.Fatal(err)
				}
				serverCode = int(merr.Number)
			}
			results, err := c.Exec(probe.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			catalogCode := 0
			var cerr *Error
			if errors.As(results[len(results)-1].Error, &cerr) {
				catalogCode = cerr.Code
			}
			if catalogCode != serverCode {
				t.Errorf("catalog answered %d, the server %d, for %s", catalogCode, serverCode, probe.sql)
			}
		}))
	}

	// The triggers that were made are the same ones, each in the database its name said.
	rows, err := o.db.QueryContext(ctx, "SELECT TRIGGER_SCHEMA, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('zz_trig_a', 'zz_trig_b') ORDER BY TRIGGER_SCHEMA, TRIGGER_NAME")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var server []string
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			t.Fatal(err)
		}
		server = append(server, schema+"."+name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	var omni []string
	for _, schema := range []string{"zz_trig_a", "zz_trig_b"} {
		var names []string
		for _, trigger := range c.GetDatabase(schema).Triggers {
			names = append(names, schema+"."+trigger.Name)
		}
		slices.Sort(names)
		omni = append(omni, names...)
	}
	if strings.Join(omni, " ") != strings.Join(server, " ") {
		t.Errorf("catalog holds %v, the server %v", omni, server)
	}
}
