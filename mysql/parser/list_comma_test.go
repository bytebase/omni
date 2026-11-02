package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// listCommaCases are the parenthesized lists of CREATE TABLE and ALTER
// TABLE: table elements, partition and subpartition definitions, and the
// partitions of ADD PARTITION and REORGANIZE PARTITION ... INTO. MySQL
// rejects an empty list and a comma after its last item (1064), as MariaDB
// and TiDB do. pr is partitioned by RANGE (a) into p0 and p1.
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

// TestListCommasOracle asserts omni agrees with a live mysql:8.0 on
// listCommaCases (container test). MySQL's verdict comes from PREPARE,
// which parses without running the statement.
func TestListCommasOracle(t *testing.T) {
	o := startParserOracle(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range listCommaCases {
		var code uint16
		_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(tc.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", tc.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != tc.code {
			t.Errorf("MySQL code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, MySQL code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
