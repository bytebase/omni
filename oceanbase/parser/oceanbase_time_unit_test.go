package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// timeUnitCase is a statement and whether OceanBase 4.3.5 parses it: false means the
// engine answers 1064.
type timeUnitCase struct {
	sql    string
	parses bool
}

// sqlTSIUnitCases: OceanBase takes none of MySQL's SQL_TSI_ unit synonyms, in any
// position a unit goes, while the plain units and SQL_TSI_ words used as names still
// parse. Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02.
var sqlTSIUnitCases = []timeUnitCase{
	{"SELECT TIMESTAMPDIFF(SQL_TSI_SECOND, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_MINUTE, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_HOUR, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_DAY, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_WEEK, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_MONTH, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_QUARTER, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(SQL_TSI_YEAR, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF(sql_tsi_second, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPADD(SQL_TSI_DAY, 1, '2020-01-01')", false},
	{"SELECT TIMESTAMPADD(SQL_TSI_SECOND, 1, '2020-01-01')", false},
	{"SELECT DATE_ADD('2020-01-01', INTERVAL 1 SQL_TSI_DAY)", false},
	{"SELECT DATE_SUB('2020-01-01', INTERVAL 1 SQL_TSI_DAY)", false},
	{"SELECT ADDDATE('2020-01-01', INTERVAL 1 SQL_TSI_DAY)", false},
	{"SELECT SUBDATE('2020-01-01', INTERVAL 1 SQL_TSI_DAY)", false},
	{"SELECT '2020-01-01' + INTERVAL 1 SQL_TSI_DAY", false},
	{"SELECT INTERVAL 1 SQL_TSI_DAY + '2020-01-01'", false},
	{"SELECT '2020-01-01' - INTERVAL 1 SQL_TSI_DAY", false},
	{"SELECT EXTRACT(SQL_TSI_DAY FROM '2020-01-01')", false},
	{"SELECT EXTRACT(SQL_TSI_YEAR FROM '2020-01-01')", false},
	{"SELECT SUM(a) OVER (ORDER BY d RANGE INTERVAL 1 SQL_TSI_DAY PRECEDING) FROM obunit_t", false},
	{"SELECT SUM(a) OVER (ORDER BY d RANGE BETWEEN INTERVAL 1 SQL_TSI_DAY PRECEDING AND CURRENT ROW) FROM obunit_t", false},
	{"CREATE TABLE t (a DATETIME, b INT AS (TIMESTAMPDIFF(SQL_TSI_DAY, a, '2020-01-01')))", false},
	{"CREATE VIEW v AS SELECT TIMESTAMPDIFF(SQL_TSI_DAY, d, d) FROM obunit_t", false},

	// The plain units in the same positions, and a SQL_TSI_ word as a name.
	{"SELECT TIMESTAMPDIFF(SECOND, '2020-01-01', '2020-01-02')", true},
	{"SELECT TIMESTAMPDIFF(MICROSECOND, '2020-01-01', '2020-01-02')", true},
	{"SELECT TIMESTAMPADD(WEEK, 1, '2020-01-01')", true},
	{"SELECT DATE_ADD('2020-01-01', INTERVAL 1 DAY)", true},
	{"SELECT DATE_ADD('2020-01-01', INTERVAL '1 1' DAY_HOUR)", true},
	{"SELECT DATE_ADD('2020-01-01', INTERVAL 1 MICROSECOND)", true},
	{"SELECT EXTRACT(DAY_HOUR FROM '2020-01-01')", true},
	{"SELECT SUM(a) OVER (ORDER BY d RANGE INTERVAL 1 DAY PRECEDING) FROM obunit_t", true},
	{"SELECT sql_tsi_second FROM (SELECT 1 AS sql_tsi_second) t", true},
	// An event schedule refuses a SQL_TSI_ unit too, which the parser already did.
	{"CREATE EVENT e ON SCHEDULE EVERY 1 SQL_TSI_DAY DO SELECT 1", false},
	{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1", true},

	// MySQL refuses a compound unit in TIMESTAMPDIFF and TIMESTAMPADD with 1064.
	// OceanBase reads the same date_unit there as everywhere else, so it parses them.
	{"SELECT TIMESTAMPDIFF(DAY_HOUR, '2020-01-01', '2020-01-02')", true},
	{"SELECT TIMESTAMPDIFF(YEAR_MONTH, '2020-01-01', '2020-01-02')", true},
	{"SELECT TIMESTAMPADD(DAY_HOUR, 1, '2020-01-01')", true},
	{"SELECT TIMESTAMPADD(MINUTE_SECOND, 1, '2020-01-01')", true},
	// A quoted unit or a string is still no unit.
	{"SELECT TIMESTAMPDIFF(`SECOND`, '2020-01-01', '2020-01-02')", false},
	{"SELECT TIMESTAMPDIFF('SECOND', '2020-01-01', '2020-01-02')", false},
}

// TestOceanBaseTimeUnitGrammar holds the parser to sqlTSIUnitCases. A refused unit
// fails at the unit, as OceanBase's "near 'SQL_TSI_DAY ...'" does.
// TestOceanBaseTimeUnitGrammarOracle asks the engine the same statements.
func TestOceanBaseTimeUnitGrammar(t *testing.T) {
	for _, c := range sqlTSIUnitCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.parses {
			t.Errorf("Parse(%q) error = %v, want parses %v", c.sql, err, c.parses)
			continue
		}
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", c.sql, err)
			continue
		}
		want := strings.Index(strings.ToUpper(c.sql), "SQL_TSI_")
		if want < 0 {
			want = strings.IndexAny(c.sql, "`'")
		}
		if pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", c.sql, pe.Position, want, err)
		}
	}
}

// TestOceanBaseTimeUnitGrammarOracle asks OceanBase about sqlTSIUnitCases through
// PREPARE, which parses without running anything: 1064 is a syntax error, any other
// answer means the statement parsed.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseTimeUnitGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obunit",
		"USE obunit",
		"DROP TABLE IF EXISTS obunit_t",
		"CREATE TABLE obunit_t (a INT, d DATETIME)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obunit")
	})
	for _, c := range sqlTSIUnitCases {
		_, err := db.ExecContext(ctx, "PREPARE obunit_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
