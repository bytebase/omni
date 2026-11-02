package catalog

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// rangeColumnsCase is a RANGE COLUMNS table with MAXVALUE in its bounds: code is
// the error the engine answers, 0 when it runs the statement.
type rangeColumnsCase struct {
	name   string
	setup  string
	sql    string
	table  string
	code   int
	bounds []string // the bounds the engine prints, for a statement it runs
}

// rangeColumnsCases are MySQL 8.0's verdicts. Asserting every bound of a value
// list to an expression panicked Exec on MAXVALUE beside other bounds, and an
// all-MAXVALUE tuple was read back as one MAXVALUE.
var rangeColumnsCases = []rangeColumnsCase{
	{"maxvalue beside a bound", "", "CREATE TABLE rc1 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE))", "rc1", 0, []string{"5,MAXVALUE", "MAXVALUE,MAXVALUE"}},
	{"three columns", "", "CREATE TABLE rc2 (a INT, b VARCHAR(10), c DATE) PARTITION BY RANGE COLUMNS(a, b, c) (PARTITION p0 VALUES LESS THAN (1, 'm', MAXVALUE), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE, MAXVALUE))", "rc2", 0, []string{"1,'m',MAXVALUE", "MAXVALUE,MAXVALUE,MAXVALUE"}},
	{"added partition", "CREATE TABLE rc3 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, 10))", "ALTER TABLE rc3 ADD PARTITION (PARTITION p1 VALUES LESS THAN (10, MAXVALUE))", "rc3", 0, []string{"5,10", "10,MAXVALUE"}},
	{"maxvalue before a greater bound", "", "CREATE TABLE rc4 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE), PARTITION p1 VALUES LESS THAN (6, 1))", "rc4", 0, []string{"5,MAXVALUE", "6,1"}},
	{"string bounds", "", "CREATE TABLE rc5 (a VARCHAR(10), b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN ('a', MAXVALUE), PARTITION p1 VALUES LESS THAN ('b', 1))", "rc5", 0, []string{"'a',MAXVALUE", "'b',1"}},
	{"maxvalue then a value", "", "CREATE TABLE rr1 (a INT) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (MAXVALUE), PARTITION p1 VALUES LESS THAN (10))", "rr1", 1493, nil},
	{"maxvalue tuples", "", "CREATE TABLE rr2 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (MAXVALUE, MAXVALUE), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE))", "rr2", 1493, nil},
	{"maxvalue ends the comparison", "", "CREATE TABLE rr3 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (MAXVALUE, 5), PARTITION p1 VALUES LESS THAN (MAXVALUE, 10))", "rr3", 1493, nil},
	{"maxvalue after an equal column", "", "CREATE TABLE rr4 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE), PARTITION p1 VALUES LESS THAN (5, 10))", "rr4", 1493, nil},
	{"string bounds decreasing", "", "CREATE TABLE rr5 (a VARCHAR(10), b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN ('b', MAXVALUE), PARTITION p1 VALUES LESS THAN ('a', 1))", "rr5", 1493, nil},
	{"language collation", "", "CREATE TABLE rc6 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_spanish2_ci) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('cz'), PARTITION p1 VALUES LESS THAN ('ch'))", "rc6", 0, []string{"'cz'", "'ch'"}},
	{"unpadded dates", "", "CREATE TABLE rc8 (a DATE) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('2020-9-1'), PARTITION p1 VALUES LESS THAN ('2020-10-01'))", "rc8", 0, []string{"'2020-9-1'", "'2020-10-01'"}},
	{"padded dates decreasing", "", "CREATE TABLE rr7 (a DATE) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('2021-01-01'), PARTITION p1 VALUES LESS THAN ('2020-12-31'))", "rr7", 1493, nil},
	{"range maxvalue not last", "", "CREATE TABLE rr6 (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN MAXVALUE, PARTITION p1 VALUES LESS THAN (10))", "rr6", 1481, nil},
}

// TestRangeColumnsMaxvalue holds the catalog to rangeColumnsCases without an
// engine: no panic, the engine's verdict, and every bound read back.
func TestRangeColumnsMaxvalue(t *testing.T) {
	for _, tc := range rangeColumnsCases {
		t.Run(tc.name, func(t *testing.T) {
			code, ddl := rangeColumnsOmni(t, tc)
			if code != tc.code {
				t.Fatalf("catalog answered %d, want %d", code, tc.code)
			}
			if got := rangeBounds(ddl); code == 0 && strings.Join(got, " ") != strings.Join(tc.bounds, " ") {
				t.Errorf("bounds %q, want %q", got, tc.bounds)
			}
		})
	}
}

// TestRangeColumnsMaxvalue_Oracle proves rangeColumnsCases against MySQL 8.0,
// and that the catalog reads back the partition bounds the engine prints.
func TestRangeColumnsMaxvalue_Oracle(t *testing.T) {
	db := func() *sql.DB { ctr, cleanup := startContainer(t); t.Cleanup(cleanup); return ctr.db }()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, q := range []string{"CREATE DATABASE IF NOT EXISTS rc_test", "USE rc_test"} {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	defer conn.ExecContext(ctx, "DROP DATABASE IF EXISTS rc_test")
	for _, tc := range rangeColumnsCases {
		t.Run(tc.name, func(t *testing.T) {
			conn.ExecContext(ctx, "DROP TABLE IF EXISTS "+tc.table)
			if tc.setup != "" {
				if _, err := conn.ExecContext(ctx, tc.setup); err != nil {
					t.Fatalf("setup: %v", err)
				}
			}
			code := 0
			if _, err := conn.ExecContext(ctx, tc.sql); err != nil {
				var me *gomysql.MySQLError
				if !errors.As(err, &me) {
					t.Fatalf("the engine gave no verdict: %v", err)
				}
				code = int(me.Number)
			}
			if code != tc.code {
				t.Fatalf("the engine answered %d, the case says %d", code, tc.code)
			}
			omniCode, omniDDL := rangeColumnsOmni(t, tc)
			if omniCode != code {
				t.Fatalf("catalog answered %d, the engine %d", omniCode, code)
			}
			if code != 0 {
				return
			}
			var name, engineDDL string
			if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE "+tc.table).Scan(&name, &engineDDL); err != nil {
				t.Fatal(err)
			}
			if got, want := rangeBounds(omniDDL), rangeBounds(engineDDL); strings.Join(got, " ") != strings.Join(want, " ") {
				t.Errorf("bounds %q, the engine prints %q", got, want)
			}
			if want := rangeBounds(engineDDL); strings.Join(want, " ") != strings.Join(tc.bounds, " ") {
				t.Errorf("the engine prints %q, the case says %q", want, tc.bounds)
			}
		})
	}
}

var rangeBoundRE = regexp.MustCompile(`VALUES LESS THAN \(([^)]*)\)`)

func rangeBounds(ddl string) []string {
	var out []string
	for _, m := range rangeBoundRE.FindAllStringSubmatch(ddl, -1) {
		out = append(out, m[1])
	}
	return out
}

func rangeColumnsOmni(t *testing.T, tc rangeColumnsCase) (int, string) {
	t.Helper()
	c := New()
	script := "CREATE DATABASE rc_test; USE rc_test; "
	if tc.setup != "" {
		script += tc.setup + "; "
	}
	results, err := c.Exec(script+tc.sql+";", nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	for _, r := range results {
		var ce *Error
		if errors.As(r.Error, &ce) {
			return ce.Code, ""
		}
		if r.Error != nil {
			t.Fatalf("Exec: %v", r.Error)
		}
	}
	return 0, c.ShowCreateTable("rc_test", tc.table)
}

// TestRangeColumnsUnsignedPastInt64: the lexer clamps an integer literal past
// int64, so two different BIGINT UNSIGNED bounds there can read as one value;
// the check leaves them open rather than refuse the table, which the engine
// creates. (The engine then answers 1654 to SHOW CREATE TABLE for it, so the
// oracle does not read it back.)
func TestRangeColumnsUnsignedPastInt64(t *testing.T) {
	tc := rangeColumnsCase{"unsigned past int64", "", "CREATE TABLE ru1 (a BIGINT UNSIGNED) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (9223372036854775807), PARTITION p1 VALUES LESS THAN (9223372036854775808))", "ru1", 0, nil}
	if code, _ := rangeColumnsOmni(t, tc); code != 0 {
		t.Errorf("catalog answered %d for bounds the engine takes", code)
	}
}
