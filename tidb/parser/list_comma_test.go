package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// listCommaCases are the parenthesized lists of CREATE TABLE and ALTER
// TABLE: table elements, partition and subpartition definitions, and the
// partitions of ADD PARTITION and REORGANIZE PARTITION ... INTO. TiDB
// rejects an empty list and a comma after its last item (1064). pr is
// partitioned by RANGE (a) into p0 and p1.
var listCommaCases = []alterDropCodeCase{
	{"CREATE TABLE t (id INT)", 0},
	{"CREATE TABLE t (id INT,)", 1064},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id),)", 1064},
	{"CREATE TABLE t ()", 1064},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", 0},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10),)", 1064},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) ()", 1064},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))", 0},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0,))", 1064},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) ())", 1064},
	{"ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))", 0},
	{"ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30),)", 1064},
	{"ALTER TABLE pr ADD PARTITION ()", 1064},
	{"ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20))", 0},
	{"ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20),)", 1064},
	{"ALTER TABLE pr REORGANIZE PARTITION p1 INTO ()", 1064},
}

// TestListCommas holds the parser to listCommaCases. A rejected list
// fails at the ')' that closes it early.
func TestListCommas(t *testing.T) {
	for _, tc := range listCommaCases {
		_, err := Parse(tc.sql)
		if (err == nil) != (tc.code == 0) {
			t.Errorf("Parse(%q) error = %v, want rejected %v", tc.sql, err, tc.code != 0)
			continue
		}
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", tc.sql, err)
			continue
		}
		at := strings.Index(tc.sql, ",)")
		if at < 0 {
			at = strings.Index(tc.sql, "()")
		}
		if want := at + 1; pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", tc.sql, pe.Position, want, err)
		}
	}
}

// TestListCommasTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// listCommaCases (container test). TiDB's verdict comes from PREPARE,
// which parses without running the statement.
func TestListCommasTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	// USE and PREPARE are session state: a dedicated connection keeps every
	// statement on the database it selects.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()
	for _, ddl := range []string{
		"DROP DATABASE IF EXISTS omni_list_comma",
		"CREATE DATABASE omni_list_comma",
		"USE omni_list_comma",
		"CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))",
	} {
		if _, err := conn.ExecContext(tc.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_list_comma")

	for _, c := range listCommaCases {
		var code uint16
		_, err := conn.ExecContext(tc.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", c.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != c.code {
			t.Errorf("TiDB code %d, want %d for %q: %v", code, c.code, c.sql, err)
			continue
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, TiDB code %d for %q: %v", omniAccepts, code, c.sql, perr)
		}
	}
}
