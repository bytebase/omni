package parser

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// TestApplicationTimePeriodCreate covers CREATE TABLE with a named (application-
// time) PERIOD FOR <name> (start, end) — container-verified accepted by
// mariadb:11.8.8, including alongside SYSTEM_TIME (bitemporal).
func TestApplicationTimePeriodCreate(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e))",
		"CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time (s, e))",
		// The period name follows the identifier rule — non-reserved keywords
		// (history, current, ...) are valid period names.
		"CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR history(s, e))",
		"CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR current(s, e))",
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, " +
			"re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, s DATE, e DATE, " +
			"PERIOD FOR SYSTEM_TIME(rs, re), PERIOD FOR app_time(s, e)) WITH SYSTEM VERSIONING",
	} {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestWithoutOverlapsParse covers WITHOUT OVERLAPS on a UNIQUE / PRIMARY KEY
// application-time period key part — container-verified accepted by 11.8.8.
func TestWithoutOverlapsParse(t *testing.T) {
	const p = "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), "
	for _, sql := range []string{
		p + "UNIQUE (id, app_time WITHOUT OVERLAPS))",
		p + "PRIMARY KEY (id, app_time WITHOUT OVERLAPS))",
		p + "UNIQUE KEY uk (id, app_time WITHOUT OVERLAPS))",
	} {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestWithoutOverlapsReject: a WITHOUT OVERLAPS key part needs an ordinary key
// column (only-period key => 1064 vs 11.8.8).
func TestWithoutOverlapsReject(t *testing.T) {
	ParseExpectError(t, "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (app_time WITHOUT OVERLAPS))")
}

// TestWithoutOverlapsGrammarReject: WITHOUT OVERLAPS is only valid on a
// UNIQUE/PRIMARY KEY and needs an ordinary key column (else 1064 vs 11.8.8),
// across table constraints, CREATE INDEX, and ALTER.
func TestWithoutOverlapsGrammarReject(t *testing.T) {
	const sv = "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), "
	for _, sql := range []string{
		sv + "KEY k (id, app_time WITHOUT OVERLAPS))",             // non-unique table KEY
		"CREATE INDEX k ON t (id, app_time WITHOUT OVERLAPS)",     // non-unique CREATE INDEX
		"ALTER TABLE t ADD KEY k (id, app_time WITHOUT OVERLAPS)", // non-unique ALTER ADD KEY
		"CREATE UNIQUE INDEX ux ON t (app_time WITHOUT OVERLAPS)", // unique, only-period
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestWithoutOverlapsBareOnly: WITHOUT OVERLAPS is valid only on a bare column
// key part — no functional expression, prefix length, or ASC/DESC (1064 vs 11.8.8).
func TestWithoutOverlapsBareOnly(t *testing.T) {
	const sv = "CREATE TABLE t (id INT, p VARCHAR(20), s DATE, e DATE, PERIOD FOR app_time(s, e), "
	for _, sql := range []string{
		sv + "UNIQUE (id, (id + 1) WITHOUT OVERLAPS))",      // functional
		sv + "UNIQUE (id, p(1) WITHOUT OVERLAPS))",          // prefix length
		sv + "UNIQUE (id, app_time DESC WITHOUT OVERLAPS))", // DESC
		sv + "UNIQUE (id, app_time ASC WITHOUT OVERLAPS))",  // ASC
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestWithoutOverlapsPositionReject: WITHOUT OVERLAPS must appear exactly once
// as the last key part (else 1064 vs 11.8.8).
func TestWithoutOverlapsPositionReject(t *testing.T) {
	const sv = "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), "
	for _, sql := range []string{
		sv + "UNIQUE (id, app_time WITHOUT OVERLAPS, e WITHOUT OVERLAPS))", // two overlaps parts
		sv + "UNIQUE (app_time WITHOUT OVERLAPS, id))",                     // not the last part
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestEmptyPrefixParensReject: an empty prefix "()" is invalid (1064 vs 11.8.8),
// generally and under WITHOUT OVERLAPS.
func TestEmptyPrefixParensReject(t *testing.T) {
	const base = "CREATE TABLE t (id INT, p VARCHAR(20), s DATE, e DATE, PERIOD FOR app_time(s, e), "
	for _, sql := range []string{
		base + "UNIQUE (id, p()))",                         // general empty prefix
		base + "UNIQUE (id, app_time() WITHOUT OVERLAPS))", // empty prefix bypassing bare check
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestWithoutOverlapsQuotedKeywordReject: a backtick-quoted `OVERLAPS` is an
// identifier, not the keyword, so WITHOUT `OVERLAPS` is 1064 (vs 11.8.8).
func TestWithoutOverlapsQuotedKeywordReject(t *testing.T) {
	bt := "`"
	sql := "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), " +
		"UNIQUE (id, app_time WITHOUT " + bt + "OVERLAPS" + bt + "))"
	ParseExpectError(t, sql)
}

// TestForPortionOf covers the application-time FOR PORTION OF <period> FROM x TO y
// clause on UPDATE / DELETE — container-verified accepted by mariadb:11.8.8.
func TestForPortionOf(t *testing.T) {
	for _, sql := range []string{
		"UPDATE t FOR PORTION OF app_time FROM '2020-01-01' TO '2021-01-01' SET id = 1",
		"DELETE FROM t FOR PORTION OF app_time FROM '2020-01-01' TO '2021-01-01'",
		// Non-reserved keyword period names are valid.
		"UPDATE t FOR PORTION OF history FROM '2020-01-01' TO '2021-01-01' SET id = 1",
	} {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestForPortionOfAST: the period name and FROM/TO bounds are captured.
func TestForPortionOfAST(t *testing.T) {
	u := parseSeqStmt[*ast.UpdateStmt](t, "UPDATE t FOR PORTION OF app_time FROM '2020-01-01' TO '2021-01-01' SET id = 1")
	if u.ForPortionOf == nil || u.ForPortionOf.PeriodName != "app_time" || u.ForPortionOf.From == nil || u.ForPortionOf.To == nil {
		t.Fatalf("UPDATE ForPortionOf = %+v", u.ForPortionOf)
	}
	if got := ast.NodeToString(u); !strings.Contains(got, ":for_portion_of {FOR_PORTION_OF :period app_time") {
		t.Errorf("NodeToString missing for_portion_of:\n%s", got)
	}
	d := parseSeqStmt[*ast.DeleteStmt](t, "DELETE FROM t FOR PORTION OF app_time FROM '2020-01-01' TO '2021-01-01'")
	if d.ForPortionOf == nil || d.ForPortionOf.PeriodName != "app_time" {
		t.Fatalf("DELETE ForPortionOf = %+v", d.ForPortionOf)
	}
}

// TestForPortionOfReject: SYSTEM_TIME is not a valid FOR PORTION OF period, and
// PORTION must be the keyword (a quoted `PORTION` does not match) — 1064 vs 11.8.8.
func TestForPortionOfReject(t *testing.T) {
	for _, sql := range []string{
		"UPDATE t FOR PORTION OF SYSTEM_TIME FROM '2020-01-01' TO '2021-01-01' SET id = 1",
		"DELETE FROM t FOR PORTION OF SYSTEM_TIME FROM '2020-01-01' TO '2021-01-01'",
		"UPDATE t FOR `PORTION` OF app_time FROM '2020-01-01' TO '2021-01-01' SET id = 1",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestForPortionOfAliasPosition: the table alias goes AFTER FOR PORTION OF for
// UPDATE (the target before the clause must be a bare base table — no alias, no
// join) but BEFORE it for DELETE. Grounded vs mariadb:11.8.8.
func TestForPortionOfAliasPosition(t *testing.T) {
	const p = "FROM '2020-01-01' TO '2021-01-01'"
	for _, sql := range []string{
		"UPDATE t FOR PORTION OF app_time " + p + " AS x SET x.id = 1", // UPDATE alias after (AS)
		"UPDATE t FOR PORTION OF app_time " + p + " x SET x.id = 1",    // UPDATE alias after (implicit)
		"DELETE FROM t AS x FOR PORTION OF app_time " + p,              // DELETE alias before
	} {
		t.Run("ok/"+sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
	for _, sql := range []string{
		"UPDATE t AS x FOR PORTION OF app_time " + p + " SET x.id = 1",                  // UPDATE alias before (AS)
		"UPDATE t x FOR PORTION OF app_time " + p + " SET x.id = 1",                     // UPDATE alias before (implicit)
		"UPDATE t JOIN u ON t.id = u.id FOR PORTION OF app_time " + p + " SET t.id = 1", // joined target
		"DELETE FROM t FOR PORTION OF app_time " + p + " AS x",                          // DELETE alias after
	} {
		t.Run("reject/"+sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
	// The alias after the clause is captured on the target table.
	u := parseSeqStmt[*ast.UpdateStmt](t, "UPDATE t FOR PORTION OF app_time "+p+" AS x SET x.id = 1")
	if tref, ok := u.Tables[0].(*ast.TableRef); !ok || tref.Alias != "x" || u.ForPortionOf == nil {
		t.Fatalf("alias not captured after clause: %+v", u.Tables[0])
	}
}
