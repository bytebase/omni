package catalog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// intervalCases run on the catalog and on the server: what the statement answers, and the
// partitions the table has after it, as name=bound. Every one was read from pingcap/tidb:v8.5.5 on
// 2026-11-05.
var intervalCases = []struct {
	name string
	ddl  string
	code int
	want []string
}{
	{name: "int", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", want: []string{"P_LT_10=10", "P_LT_20=20", "P_LT_30=30"}},
	// A sign on a bound is read through; TiDB v8.5.5 takes LESS THAN (+10).
	{name: "signed bounds", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (+10) LAST PARTITION LESS THAN (+30)", want: []string{"P_LT_10=10", "P_LT_20=20", "P_LT_30=30"}},
	{name: "a step written as a sum", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (5+5) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", want: []string{"P_LT_10=10", "P_LT_20=20", "P_LT_30=30"}},
	{name: "negative", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (-10) LAST PARTITION LESS THAN (10)", want: []string{"P_LT_-10=-10", "P_LT_0=0", "P_LT_10=10"}},
	{name: "arithmetic", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1+1) FIRST PARTITION LESS THAN (1+1) LAST PARTITION LESS THAN (2*3)", want: []string{"P_LT_2=2", "P_LT_4=4", "P_LT_6=6"}},
	{name: "one partition", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10)", want: []string{"P_LT_10=10"}},
	{name: "null and maxvalue", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20) NULL PARTITION MAXVALUE PARTITION", want: []string{"P_NULL=-9223372036854775808", "P_LT_10=10", "P_LT_20=20", "P_MAXVALUE=MAXVALUE"}},
	{name: "null, unsigned", ddl: "CREATE TABLE t (a INT UNSIGNED) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=0", "P_LT_10=10"}},
	{name: "null, expression", ddl: "CREATE TABLE t (a TIMESTAMP) PARTITION BY RANGE (UNIX_TIMESTAMP(a)) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=-9223372036854775808", "P_LT_10=10"}},
	{name: "null, columns int", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=-2147483648", "P_LT_10=10"}},
	{name: "null, columns tinyint", ddl: "CREATE TABLE t (a TINYINT) PARTITION BY RANGE COLUMNS (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=-128", "P_LT_10=10"}},
	{name: "null, columns unsigned", ddl: "CREATE TABLE t (a SMALLINT UNSIGNED) PARTITION BY RANGE COLUMNS (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=0", "P_LT_10=10"}},
	{name: "null, columns bigint", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE COLUMNS (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (10) NULL PARTITION", want: []string{"P_NULL=-9223372036854775808", "P_LT_10=10"}},
	{name: "date by day", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-03')", want: []string{"P_LT_2020-01-01='2020-01-01'", "P_LT_2020-01-02='2020-01-02'", "P_LT_2020-01-03='2020-01-03'"}},
	{name: "date by month keeps the end of the month", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-31') LAST PARTITION LESS THAN ('2020-03-31')", want: []string{"P_LT_2020-01-31='2020-01-31'", "P_LT_2020-02-29='2020-02-29'", "P_LT_2020-03-31='2020-03-31'"}},
	{name: "date null", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-02') NULL PARTITION", want: []string{"P_NULL='0000-01-01'", "P_LT_2020-01-01='2020-01-01'", "P_LT_2020-01-02='2020-01-02'"}},
	{name: "datetime by hour", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 HOUR) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-01 02:00:00')", want: []string{"P_LT_2020-01-01 00:00:00='2020-01-01 00:00:00'", "P_LT_2020-01-01 01:00:00='2020-01-01 01:00:00'", "P_LT_2020-01-01 02:00:00='2020-01-01 02:00:00'"}},
	{name: "datetime by year, null and maxvalue", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 YEAR) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01') NULL PARTITION MAXVALUE PARTITION", want: []string{"P_NULL='0000-01-01 00:00:00'", "P_LT_2020-01-01 00:00:00='2020-01-01 00:00:00'", "P_LT_2021-01-01 00:00:00='2021-01-01 00:00:00'", "P_MAXVALUE=MAXVALUE"}},
	{name: "datetime by week", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 WEEK) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-08')", want: []string{"P_LT_2020-01-01 00:00:00='2020-01-01 00:00:00'", "P_LT_2020-01-08 00:00:00='2020-01-08 00:00:00'"}},
	{name: "datetime by quarter", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 QUARTER) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-04-01')", want: []string{"P_LT_2020-01-01 00:00:00='2020-01-01 00:00:00'", "P_LT_2020-04-01 00:00:00='2020-04-01 00:00:00'"}},
	{name: "datetime keeps its time", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 DAY) FIRST PARTITION LESS THAN ('2020-01-01 12:00:00') LAST PARTITION LESS THAN ('2020-01-02 12:00:00')", want: []string{"P_LT_2020-01-01 12:00:00='2020-01-01 12:00:00'", "P_LT_2020-01-02 12:00:00='2020-01-02 12:00:00'"}},
	{name: "written partitions keep their names", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30) (PARTITION a VALUES LESS THAN (10), PARTITION b VALUES LESS THAN (20), PARTITION c VALUES LESS THAN (30))", want: []string{"a=10", "b=20", "c=30"}},
	{name: "PARTITIONS is dropped", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20) PARTITIONS 9", want: []string{"P_LT_10=10", "P_LT_20=20"}},
	{name: "SUBPARTITION BY is dropped", ddl: "CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20) SUBPARTITION BY HASH (b)", want: []string{"P_LT_10=10", "P_LT_20=20"}},

	// Each bound is FIRST + i * INTERVAL in BIGINT, worked out up to LAST, and a date bound is
	// whole seconds or months from FIRST, however many years that spans. Read on 2026-11-06.
	{name: "int at the top of bigint", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (9223372036854775807) LAST PARTITION LESS THAN (9223372036854775807)", want: []string{"P_LT_9223372036854775807=9223372036854775807"}},
	{name: "int up to the top of bigint", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (9223372036854775806) LAST PARTITION LESS THAN (9223372036854775807)", want: []string{"P_LT_9223372036854775806=9223372036854775806", "P_LT_9223372036854775807=9223372036854775807"}},
	{name: "the largest int step", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (9223372036854775807) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (9223372036854775807)", want: []string{"P_LT_0=0", "P_LT_9223372036854775807=9223372036854775807"}},
	{name: "date by 200000 days", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (200000 DAY) FIRST PARTITION LESS THAN ('1000-01-01') LAST PARTITION LESS THAN ('1547-08-02')", want: []string{"P_LT_1000-01-01='1000-01-01'", "P_LT_1547-08-02='1547-08-02'"}},
	{name: "date by 100000 weeks", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (100000 WEEK) FIRST PARTITION LESS THAN ('1000-01-01') LAST PARTITION LESS THAN ('2916-07-15')", want: []string{"P_LT_1000-01-01='1000-01-01'", "P_LT_2916-07-15='2916-07-15'"}},
	{name: "datetime by 10^10 seconds", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (10000000000 SECOND) FIRST PARTITION LESS THAN ('1000-01-01') LAST PARTITION LESS THAN ('1316-11-20 17:46:40')", want: []string{"P_LT_1000-01-01 00:00:00='1000-01-01 00:00:00'", "P_LT_1316-11-20 17:46:40='1316-11-20 17:46:40'"}},
	{name: "date by 5000 years", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (5000 YEAR) FIRST PARTITION LESS THAN ('1000-01-01') LAST PARTITION LESS THAN ('6000-01-01')", want: []string{"P_LT_1000-01-01='1000-01-01'", "P_LT_6000-01-01='6000-01-01'"}},
	{name: "a date step past the range, one partition", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1000000000000 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-01')", want: []string{"P_LT_2020-01-01='2020-01-01'"}},
	{name: "a month step past the range, one partition", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (9223372036854775807 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-01')", want: []string{"P_LT_2020-01-01='2020-01-01'"}},

	{name: "no FIRST and LAST", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10)", code: 8200},
	{name: "no FIRST and LAST, with partitions", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) (PARTITION p VALUES LESS THAN (1))", code: 8200},
	{name: "LAST off the steps", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (35)", code: 8200},
	{name: "LAST below FIRST", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (30) LAST PARTITION LESS THAN (10)", code: 8200},
	{name: "a zero step", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (0) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", code: 8200},
	{name: "a string step", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL ('x') FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)", code: 8200},
	{name: "a string bound", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN ('x') LAST PARTITION LESS THAN (20)", code: 1697},
	{name: "a time unit on an int", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10 DAY) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", code: 1292},
	{name: "no time unit on a datetime", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (10) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01')", code: 1292},
	{name: "a compound unit", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 YEAR_MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01')", code: 8200},
	// TiDB checks the step's text, which must start with a digit 1 to 9
	// (v8.5.5): a sign or a parenthesis before it answers 8200.
	{name: "a signed step", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (+10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", code: 8200},
	{name: "a parenthesized step", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL ((10)) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30)", code: 8200},
	{name: "a signed time step", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (+1 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-03')", code: 8200},
	{name: "a negative time step", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (-1 DAY) FIRST PARTITION LESS THAN ('2020-01-03') LAST PARTITION LESS THAN ('2020-01-01')", code: 8200},
	{name: "a time LAST off the steps", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (2 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-02')", code: 8200},
	{name: "hours on a date repeat a name", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 HOUR) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-02')", code: 1517},
	{name: "a varchar column", ddl: "CREATE TABLE t (a VARCHAR(9)) PARTITION BY RANGE COLUMNS (a) INTERVAL (1) FIRST PARTITION LESS THAN ('a') LAST PARTITION LESS THAN ('c')", code: 8200},
	{name: "two columns", ddl: "CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (3)", code: 8200},
	{name: "an int step past bigint", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (2) FIRST PARTITION LESS THAN (9223372036854775806) LAST PARTITION LESS THAN (9223372036854775807)", code: 1690},
	{name: "int steps across bigint", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (9223372036854775807) FIRST PARTITION LESS THAN (-9223372036854775807) LAST PARTITION LESS THAN (9223372036854775807)", code: 1690},
	{name: "a day step past the range", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1000000000000 DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-02')", code: 1441},
	{name: "a month step past the range", ddl: "CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (9223372036854775807 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-02-01')", code: 1441},
	{name: "a second step past the range", ddl: "CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (9223372036854775807 SECOND) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-01 00:00:01')", code: 1441},
	{name: "an overflowing step", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (9223372036854775807 * 2 + 3) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (2)", code: 1690},
	{name: "an overflowing FIRST", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (9223372036854775807 + 1) LAST PARTITION LESS THAN (2)", code: 1690},
	{name: "an overflowing LAST", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (9223372036854775807 * 2)", code: 1690},
	{name: "the lowest bigint negated", ddl: "CREATE TABLE t (a BIGINT) PARTITION BY RANGE (a) INTERVAL (-(-9223372036854775807 - 1)) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (2)", code: 8200},
	// The limit is 8192 partitions, NULL PARTITION and MAXVALUE PARTITION among them.
	{name: "8191 bounds and MAXVALUE", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (8191) MAXVALUE PARTITION", code: 1499},
	{name: "8191 bounds and NULL", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (8191) NULL PARTITION", code: 1499},
	{name: "8191 bounds, NULL and MAXVALUE", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (8190) NULL PARTITION MAXVALUE PARTITION", code: 1499},
	{name: "too many partitions", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (100000)", code: 1499},
	{name: "partitions that differ", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30) (PARTITION a VALUES LESS THAN (1), PARTITION b VALUES LESS THAN (2), PARTITION c VALUES LESS THAN (3))", code: 8200},
	{name: "fewer partitions", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30) (PARTITION p VALUES LESS THAN (1))", code: 8200},
	{name: "fewer partitions that match as far as they go", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (30) (PARTITION p VALUES LESS THAN (10))", code: 8200},
	{name: "a missing column", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (nosuch) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)", code: 1488},
	{name: "a unique key that misses it", ddl: "CREATE TABLE t (a INT, b INT, UNIQUE KEY u (b)) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)", code: 8264},

	// A reserved word is a partition name when quoted, and stays quoted in SHOW CREATE TABLE.
	{name: "a reserved partition name", ddl: "CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION `ilike` VALUES LESS THAN (1), PARTITION `select` VALUES LESS THAN (2))", want: []string{"ilike=1", "select=2"}},

	// SYSTEM_TIME runs and leaves the table unpartitioned, whatever its keys cover.
	{name: "SYSTEM_TIME", ddl: "CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)"},
	{name: "SYSTEM_TIME LIMIT", ddl: "CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)"},
	{name: "SYSTEM_TIME and a unique key", ddl: "CREATE TABLE t (a INT, b INT, UNIQUE KEY u (b)) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)"},
}

func catalogPartitions(c *Catalog) []string {
	pi := c.GetDatabase("d").Tables["t"].Partitioning
	if pi == nil {
		return nil
	}
	var got []string
	for _, pd := range pi.Partitions {
		got = append(got, pd.Name+"="+pd.ValueExpr)
	}
	return got
}

// TestIntervalPartitions holds the catalog to intervalCases.
// The limit's other side: 8192 partitions, the last of them MAXVALUE, are taken.
func init() {
	want := make([]string, 0, 8192)
	for i := 0; i <= 8190; i++ {
		want = append(want, fmt.Sprintf("P_LT_%d=%d", i, i))
	}
	intervalCases = append(intervalCases, struct {
		name string
		ddl  string
		code int
		want []string
	}{
		name: "8191 bounds and MAXVALUE make 8192",
		ddl:  "CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (0) LAST PARTITION LESS THAN (8190) MAXVALUE PARTITION",
		want: append(want, "P_MAXVALUE=MAXVALUE"),
	})
}

func TestIntervalPartitions(t *testing.T) {
	for _, tc := range intervalCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			execAll(t, c, "CREATE DATABASE d; USE d;")
			if tc.code != 0 {
				var cerr *Error
				if err := execExpectError(t, c, tc.ddl+";"); !errors.As(err, &cerr) || cerr.Code != tc.code {
					t.Fatalf("error %v, want %d", err, tc.code)
				}
				return
			}
			execAll(t, c, tc.ddl+";")
			if got := catalogPartitions(c); strings.Join(got, "|") != strings.Join(tc.want, "|") {
				t.Errorf("partitions %v, want %v", got, tc.want)
			}
			if pi := c.GetDatabase("d").Tables["t"].Partitioning; pi != nil && (pi.NumParts != 0 || pi.SubType != "") {
				t.Errorf("the catalog kept what the server drops: %+v", pi)
			}
			// What SHOW CREATE TABLE writes is DDL the catalog reads back. The lowest bigint is left
			// out: the parser's integer literal stops one short of it, on any bound.
			if len(tc.want) > 0 && !strings.Contains(strings.Join(tc.want, "|"), "-9223372036854775808") {
				show := c.ShowCreateTable("d", "t")
				again := New()
				execAll(t, again, "CREATE DATABASE d; USE d;")
				execAll(t, again, show+";")
				if got := catalogPartitions(again); strings.Join(got, "|") != strings.Join(tc.want, "|") {
					t.Errorf("read back from %s:\npartitions %v, want %v", show, got, tc.want)
				}
			}
		})
	}
}

// TestIntervalPartitionsOracle holds intervalCases to the server.
func TestIntervalPartitionsOracle(t *testing.T) {
	tc := startTiDBForCatalog(t)
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range []string{"CREATE DATABASE IF NOT EXISTS interval_oracle", "USE interval_oracle"} {
		if _, err := conn.ExecContext(tc.ctx, stmt); err != nil {
			t.Fatal(err)
		}
	}
	line := regexp.MustCompile("PARTITION `([^`]+)` VALUES LESS THAN \\((.*?)\\)[,)]?\n?")
	for _, c := range intervalCases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := conn.ExecContext(tc.ctx, "DROP TABLE IF EXISTS t"); err != nil {
				t.Fatal(err)
			}
			var code int
			if _, err := conn.ExecContext(tc.ctx, c.ddl); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatal(err)
				}
				code = int(myErr.Number)
			}
			if code != c.code {
				t.Fatalf("the server answered %d, the record says %d", code, c.code)
			}
			if code != 0 {
				return
			}
			var name, ddl string
			if err := conn.QueryRowContext(tc.ctx, "SHOW CREATE TABLE t").Scan(&name, &ddl); err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, m := range line.FindAllStringSubmatch(ddl, -1) {
				got = append(got, m[1]+"="+m[2])
			}
			if strings.Join(got, "|") != strings.Join(c.want, "|") {
				t.Errorf("the server made %v, the record says %v:\n%s", got, c.want, ddl)
			}
			if len(c.want) == 0 && strings.Contains(ddl, "PARTITION") {
				t.Errorf("the server kept partitions:\n%s", ddl)
			}
		})
	}
	_, _ = conn.ExecContext(tc.ctx, "DROP TABLE IF EXISTS t")
}

// TestIntervalWrittenClause: a bound this package cannot work out keeps the clause as written,
// which the server takes and the catalog reads back.
//
// PARTITIONS n beside it names no partitions of the server's, which generates them and drops the
// count, so the empty ones it would make are not kept.
func TestIntervalWrittenClause(t *testing.T) {
	for _, suffix := range []string{" NULL PARTITION", " NULL PARTITION PARTITIONS 3"} {
		c := New()
		execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (ABS(10)) LAST PARTITION LESS THAN (30)"+suffix+";")
		show := c.ShowCreateTable("d", "t")
		for _, want := range []string{"INTERVAL (10) FIRST PARTITION LESS THAN (", "LAST PARTITION LESS THAN (30) NULL PARTITION"} {
			if !strings.Contains(show, want) {
				t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, show)
			}
		}
		if strings.Contains(show, "PARTITION p0") || strings.Contains(show, "PARTITIONS") {
			t.Errorf("SHOW CREATE TABLE keeps what PARTITIONS made:\n%s", show)
		}
		again := New()
		execAll(t, again, "CREATE DATABASE d; USE d;")
		execAll(t, again, show+";")
		if again.ShowCreateTable("d", "t") != show {
			t.Errorf("the clause did not read back:\n%s\nwas:\n%s", again.ShowCreateTable("d", "t"), show)
		}
	}
}

// TestHistoryPartitionsRefused: HISTORY and CURRENT on a table the catalog partitions at all, as
// v8.5.5 answers them.
func TestHistoryPartitionsRefused(t *testing.T) {
	for _, tc := range []struct {
		setup string
		stmt  string
		code  int
	}{
		{"CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2", "ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)", 1480},
		{"CREATE TABLE t (a INT) PARTITION BY KEY (a) PARTITIONS 2", "ALTER TABLE t ADD PARTITION (PARTITION x CURRENT)", 1480},
		{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1))", "ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)", 4113},
		{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1))", "ALTER TABLE t REORGANIZE PARTITION p INTO (PARTITION x CURRENT)", 4113},
	} {
		c := New()
		execAll(t, c, "CREATE DATABASE d; USE d; "+tc.setup+";")
		before := c.ShowCreateTable("d", "t")
		var cerr *Error
		if err := execExpectError(t, c, tc.stmt+";"); !errors.As(err, &cerr) || cerr.Code != tc.code {
			t.Errorf("%s: error %v, want %d", tc.stmt, err, tc.code)
		}
		if c.ShowCreateTable("d", "t") != before {
			t.Errorf("%s changed the table", tc.stmt)
		}
	}
}

// TestILikeViewText: a stored view keeps ILIKE, and the text reads back.
func TestILikeViewText(t *testing.T) {
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE u (a VARCHAR(9)); CREATE VIEW vv AS SELECT a ILIKE 'b' AS c, a NOT ILIKE 'b' ESCAPE 'x' AS e, a LIKE 'b' AS f FROM u;")
	def := c.GetDatabase("d").Views["vv"].Definition
	// TiDB v8.5.5 stores this body as written below (SHOW CREATE VIEW).
	if want := "SELECT `a` ILIKE _UTF8MB4'b' AS `c`,`a` NOT ILIKE _UTF8MB4'b' ESCAPE 'x' AS `e`,`a` LIKE _UTF8MB4'b' AS `f` FROM `d`.`u`"; def != want {
		t.Errorf("Definition = %s\nwant %s", def, want)
	}
	again := New()
	execAll(t, again, "CREATE DATABASE d; USE d; CREATE TABLE u (a VARCHAR(9)); CREATE VIEW vv AS "+def+";")
	if got := again.GetDatabase("d").Views["vv"].Definition; got != def {
		t.Errorf("read back %s\nwas %s", got, def)
	}
}
