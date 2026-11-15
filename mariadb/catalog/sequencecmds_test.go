package catalog

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// TestOracle_Sequences runs each script on the pinned server and on the
// catalog, and compares every statement's error code and, after the script,
// what SHOW CREATE SEQUENCE shows for each sequence named in show.
func TestOracle_Sequences(t *testing.T) {
	cases := []struct {
		name   string
		script []string
		show   []string
	}{
		{"defaults", []string{"CREATE SEQUENCE s"}, []string{"s"}},
		{"options", []string{"CREATE SEQUENCE s START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT 'c''x'"}, []string{"s"}},
		{"descending", []string{"CREATE SEQUENCE s INCREMENT BY -1"}, []string{"s"}},
		{"no limits", []string{"CREATE SEQUENCE s NOCACHE NOMINVALUE NOMAXVALUE ENGINE=MyISAM"}, []string{"s"}},
		{"engine and comment", []string{"CREATE SEQUENCE s START WITH 3 COMMENT='x' ENGINE=Aria"}, []string{"s"}},
		{"types", []string{"CREATE SEQUENCE s AS INT UNSIGNED", "CREATE SEQUENCE u AS TINYINT START WITH -5 MINVALUE -10 MAXVALUE 10"}, []string{"s", "u"}},
		{"out of range", []string{"CREATE SEQUENCE s START WITH 5 MINVALUE 10", "CREATE SEQUENCE u MAXVALUE 1 MINVALUE 1"}, nil},
		{"limits held to the type", []string{"CREATE SEQUENCE v AS TINYINT MAXVALUE 200", "CREATE SEQUENCE w AS TINYINT MINVALUE -200", "CREATE SEQUENCE x AS TINYINT UNSIGNED MAXVALUE 300", "CREATE SEQUENCE y MAXVALUE 9223372036854775807"}, []string{"v", "w", "x", "y"}},
		{"unsigned minimum", []string{"CREATE SEQUENCE z AS TINYINT UNSIGNED MINVALUE 0 START WITH 1", "CREATE SEQUENCE z0 AS TINYINT UNSIGNED MINVALUE 0 START WITH 0", "CREATE SEQUENCE z1 AS INT UNSIGNED START WITH 0"}, []string{"z"}},
		{"alter", []string{"CREATE SEQUENCE s START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000", "ALTER SEQUENCE s START WITH 50", "ALTER SEQUENCE s INCREMENT BY 2 RESTART WITH 3", "ALTER SEQUENCE s MAXVALUE 5", "ALTER SEQUENCE s NOCACHE CYCLE"}, []string{"s"}},
		{"namespace", []string{"CREATE SEQUENCE s", "CREATE SEQUENCE s", "CREATE TABLE s (a INT)", "CREATE VIEW s AS SELECT 1", "CREATE TABLE t (a INT)", "CREATE SEQUENCE t", "DROP SEQUENCE t", "ALTER SEQUENCE t CACHE 5", "ALTER SEQUENCE missing CACHE 5", "DROP SEQUENCE IF EXISTS missing"}, []string{"s"}},
		{"replace", []string{"CREATE SEQUENCE s", "CREATE OR REPLACE SEQUENCE s START WITH 7", "CREATE SEQUENCE IF NOT EXISTS s START WITH 9"}, []string{"s"}},
		{"drop and rename", []string{"CREATE SEQUENCE s", "CREATE SEQUENCE u", "DROP TABLE s", "RENAME TABLE u TO w", "CREATE SEQUENCE s START WITH 2"}, []string{"s", "w"}},
	}
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetDatabase(t, conn, database)
			c := New()
			if _, err := c.Exec("CREATE DATABASE "+quoteIdent(database)+"; USE "+quoteIdent(database)+";", nil); err != nil {
				t.Fatal(err)
			}
			for _, stmt := range tc.script {
				_, serverErr := conn.ExecContext(context.Background(), stmt)
				results, err := c.Exec(stmt, nil)
				if err != nil {
					// The parser refuses a range it can check on its own.
					if serverErr == nil {
						t.Errorf("%s: catalog refuses it (%v), the server takes it", stmt, err)
					}
					continue
				}
				if got, want := errorCode(results[0].Error), errorCode(serverErr); got != want {
					t.Errorf("%s: catalog answers %d, server %d (%v)", stmt, got, want, serverErr)
				}
			}
			for _, name := range tc.show {
				var table, want string
				if err := conn.QueryRowContext(context.Background(), "SHOW CREATE SEQUENCE "+quoteIdent(name)).Scan(&table, &want); err != nil && !errors.Is(err, sql.ErrNoRows) {
					t.Fatal(err)
				}
				if got := c.ShowCreateSequence(database, name); got != want {
					t.Errorf("sequence %s\n--- catalog\n%s\n--- server\n%s", name, got, want)
				}
			}
		})
	}
}

// errorCode is the MariaDB error number err carries, or 0.
func errorCode(err error) int {
	var server *mysql.MySQLError
	if errors.As(err, &server) {
		return int(server.Number)
	}
	var cat *Error
	if errors.As(err, &cat) {
		return cat.Code
	}
	if err != nil {
		return -1
	}
	return 0
}

func TestSequenceUnsignedLimits(t *testing.T) {
	c := New()
	_, err := c.Exec("CREATE DATABASE d; USE d; CREATE SEQUENCE su AS BIGINT UNSIGNED MAXVALUE 18446744073709551614 START WITH 18446744073709551600; CREATE OR REPLACE TABLE su (a INT); CREATE SEQUENCE s; CREATE TABLE t (a INT DEFAULT (NEXTVAL(s) + 1));", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("d").Sequences["su"] != nil || c.GetDatabase("d").Tables["su"] == nil {
		t.Error("CREATE OR REPLACE TABLE must replace the sequence of its name")
	}
	if got := c.ShowCreateTable("d", "t"); !strings.Contains(got, "DEFAULT (nextval(`d`.`s`) + 1)") {
		t.Errorf("sequence call inside an expression default not qualified:\n%s", got)
	}
	c2 := New()
	if _, err := c2.Exec("CREATE DATABASE d; USE d; CREATE SEQUENCE su AS BIGINT UNSIGNED MAXVALUE 18446744073709551614 START WITH 18446744073709551600;", nil); err != nil {
		t.Fatal(err)
	}
	want := "CREATE SEQUENCE `su` as bigint unsigned start with 18446744073709551600 minvalue 1 maxvalue 18446744073709551614 increment by 1 cache 1000 nocycle ENGINE=InnoDB"
	if got := c2.ShowCreateSequence("d", "su"); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
