package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// obOptionCase is a statement and whether OceanBase 4.3.5 parses it: false
// means the engine answers 1064.
type obOptionCase struct {
	sql    string
	parses bool
}

// tableOptionValueCases holds OceanBase's storage options to the value each
// takes. Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02. A value the
// engine parses and then refuses (BLOCK_SIZE = 1 is 5258, REPLICA_NUM = 0 is
// 1235, AUTO_INCREMENT_MODE = 'x' is 1149) still parses here.
func tableOptionValueCases() []obOptionCase {
	ints := []string{"= 'x'", "= '16'", "= TRUE", "= abc", "= SELECT", "= DEFAULT", "= NULL", "= ON", "= -1", "= 1.5"}
	bools := []string{"= 1", "= 0", "= 'x'", "= abc", "= SELECT", "= DEFAULT", "= NULL", "= ON", "= -1"}
	strs := []string{"= ORDER", "= NOORDER", "= TRUE", "= 1", "= abc", "= SELECT", "= DEFAULT", "= NULL"}
	var cases []obOptionCase
	for _, o := range []struct {
		name string
		good []string
		bad  []string
	}{
		{"BLOCK_SIZE", []string{"= 16384", "16384", "= 1"}, ints},
		{"REPLICA_NUM", []string{"= 1", "1", "= 0"}, ints},
		{"TABLET_SIZE", []string{"= 134217728", "134217728", "= 0"}, ints},
		{"PCTFREE", []string{"= 10", "10", "= 0"}, ints},
		{"USE_BLOOM_FILTER", []string{"= FALSE", "FALSE", "= TRUE", "= true"}, bools},
		{"ENABLE_MACRO_BLOCK_BLOOM_FILTER", []string{"= FALSE", "FALSE", "= TRUE", "= true"}, bools},
		{"AUTO_INCREMENT_MODE", []string{"= 'ORDER'", "'ORDER'", "= 'NOORDER'", "= 'x'"}, strs},
	} {
		for _, v := range o.good {
			cases = append(cases, obOptionCase{"CREATE TABLE t (a INT PRIMARY KEY) " + o.name + " " + v, true})
		}
		for _, v := range o.bad {
			cases = append(cases, obOptionCase{"CREATE TABLE t (a INT PRIMARY KEY) " + o.name + " " + v, false})
		}
	}
	return append(cases,
		obOptionCase{"ALTER TABLE obopt_t PCTFREE = 20", true},
		obOptionCase{"ALTER TABLE obopt_t BLOCK_SIZE = 32768", true},
		obOptionCase{"ALTER TABLE obopt_t USE_BLOOM_FILTER = TRUE", true},
		obOptionCase{"ALTER TABLE obopt_t PCTFREE = abc", false},
		obOptionCase{"ALTER TABLE obopt_t TABLET_SIZE = abc", false},
	)
}

// organizationCases: ORGANIZATION takes INDEX or HEAP, bare, with or without
// the '=', and nothing else.
var organizationCases = []obOptionCase{
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION INDEX", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION = INDEX", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION HEAP", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION = HEAP", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION index", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION heap", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION INDEX ORGANIZATION HEAP", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION HEAP COMMENT 'c'", true},
	{"CREATE TABLE t (a INT) ORGANIZATION HEAP", true},
	{"ALTER TABLE obopt_t ORGANIZATION HEAP", true},
	{"ALTER TABLE obopt_t ORGANIZATION = INDEX", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION = 'HEAP'", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION 'INDEX'", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION `HEAP`", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION foo", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION = foo", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION =", false},
}

// dynamicPartitionPolicyCases: the policy's keys are fixed and bare, each
// assignment needs its '=', ENABLE takes TRUE or FALSE, and every other key a
// string, which may be adjacent literals or double-quoted. A key may repeat, and so
// may the option. Three statements parse and are then refused with 1210: an
// EXPIRE_TIME of '-1day', the option given twice, and assignments with no comma
// between them. OceanBase prints the policy back with every key, as
// `DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY', PRECREATE_TIME =
// '-1', EXPIRE_TIME = '-1', TIME_ZONE = 'DEFAULT', BIGINT_PRECISION = 'NONE')`.
func dynamicPartitionPolicyCases() []obOptionCase {
	var cases []obOptionCase
	for _, c := range []obOptionCase{
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY')", true},
		{"DYNAMIC_PARTITION_POLICY (ENABLE = TRUE, TIME_UNIT = 'DAY')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY', ENABLE = FALSE)", true},
		{"DYNAMIC_PARTITION_POLICY = (enable = true, time_unit = 'day')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'hour', PRECREATE_TIME = '3hour', EXPIRE_TIME = '1day', TIME_ZONE = '+8:00', BIGINT_PRECISION = 'none')", true},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY', ENABLE = FALSE)", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY') DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE)", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY', EXPIRE_TIME = '-1day')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY') COMMENT 'c'", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = \"DAY\")", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DA' 'Y')", true},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE TIME_UNIT = 'DAY')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY' PRECREATE_TIME = '1day')", true},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY', PRECREATE_TIME = '-1', EXPIRE_TIME = '-1', TIME_ZONE = 'DEFAULT', BIGINT_PRECISION = 'NONE')", true},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY', START_TIME = '1day')", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY',, ENABLE = TRUE)", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = _utf8mb4'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT := 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE TRUE, TIME_UNIT 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (x y)", false},
		{"DYNAMIC_PARTITION_POLICY = (x = 1)", false},
		{"DYNAMIC_PARTITION_POLICY = (x = 'a')", false},
		{"DYNAMIC_PARTITION_POLICY = ('TIME_UNIT' = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (`TIME_UNIT` = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = DAY)", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 1.5)", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY', PRECREATE_TIME = 3)", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY', PRECREATE_TIME = -1)", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = 'true', TIME_UNIT = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = 1, TIME_UNIT = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = ON, TIME_UNIT = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = (ENABLE = NULL, TIME_UNIT = 'DAY')", false},
		{"DYNAMIC_PARTITION_POLICY = ()", false},
		{"DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY',)", false},
		{"DYNAMIC_PARTITION_POLICY = ENABLE", false},
		{"DYNAMIC_PARTITION_POLICY = 'x'", false},
		{"DYNAMIC_PARTITION_POLICY", false},
	} {
		cases = append(cases, obOptionCase{
			"CREATE TABLE t (a INT, b DATETIME) " + c.sql + " PARTITION BY RANGE COLUMNS(b) (PARTITION p0 VALUES LESS THAN ('2026-01-01'))",
			c.parses,
		})
	}
	return append(cases,
		obOptionCase{"ALTER TABLE obopt_dp DYNAMIC_PARTITION_POLICY = (ENABLE = FALSE)", true},
		obOptionCase{"ALTER TABLE obopt_dp DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY')", true},
		obOptionCase{"ALTER TABLE obopt_dp DYNAMIC_PARTITION_POLICY = (TIME_UNIT = DAY)", false},
	)
}

// storageAndStatsOptionCases holds the STORAGE and STATS_PERSISTENT table
// options to the value each takes, measured against OceanBase CE 4.3.5.
// Unlike the mysql and mariadb forks, OceanBase has no STORAGE table option at
// all: `STORAGE DISK` and `STORAGE MEMORY` are a 1064 here, as is every other
// value, so every STORAGE row rejects. STATS_PERSISTENT is `[=] {DEFAULT |
// number}`, so a string and a bare identifier are each a 1064. The out-of-[0,1]
// STATS_PERSISTENT numbers the mysql fork rejects are left out: this change
// only closes the string/identifier hole and still takes any integer, so it
// asserts no bound it has not measured.
var storageAndStatsOptionCases = []obOptionCase{
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE DISK", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE MEMORY", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE = DISK", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE 's'", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE foo", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE = 1", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STORAGE DEFAULT", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT = DEFAULT", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT = 0", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT 1", true},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT = 's'", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT DEFINED", false},
	{"CREATE TABLE t (a INT PRIMARY KEY) STATS_PERSISTENT = foo", false},
}

func allTableOptionCases() []obOptionCase {
	cases := tableOptionValueCases()
	cases = append(cases, organizationCases...)
	cases = append(cases, storageAndStatsOptionCases...)
	return append(cases, dynamicPartitionPolicyCases()...)
}

// TestOceanBaseTableOptionGrammar holds the parser to the OceanBase verdicts
// above. TestOceanBaseTableOptionGrammarOracle asks the engine the same
// statements.
func TestOceanBaseTableOptionGrammar(t *testing.T) {
	for _, c := range allTableOptionCases() {
		if _, err := Parse(c.sql); (err == nil) != c.parses {
			t.Errorf("Parse(%q) error = %v, want parses %v", c.sql, err, c.parses)
		}
	}
}

// TestOceanBaseTableOptionGrammarOracle asks OceanBase about every case of
// TestOceanBaseTableOptionGrammar. The verdict comes from PREPARE, which
// parses without creating anything: 1064 is a syntax error, any other answer
// means the statement parsed.
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseTableOptionGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obopt",
		"USE obopt",
		"DROP TABLE IF EXISTS obopt_t, obopt_dp",
		"CREATE TABLE obopt_t (a INT PRIMARY KEY)",
		"CREATE TABLE obopt_dp (a INT, b DATETIME) DYNAMIC_PARTITION_POLICY = (TIME_UNIT = 'DAY') PARTITION BY RANGE COLUMNS(b) (PARTITION p0 VALUES LESS THAN ('2026-01-01'))",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obopt")
	})
	for _, c := range allTableOptionCases() {
		_, err := db.ExecContext(ctx, "PREPARE obopt_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		engineParses := err == nil
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			engineParses = myErr.Number != 1064
		}
		if engineParses != c.parses {
			t.Errorf("OceanBase parses=%v, want %v for %q: %v", engineParses, c.parses, c.sql, err)
		}
		if _, perr := Parse(c.sql); (perr == nil) != engineParses {
			t.Errorf("DIVERGENCE omni parses=%v, OceanBase parses=%v for %q: %v", perr == nil, engineParses, c.sql, perr)
		}
	}
}
