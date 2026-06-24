package parser

import "testing"

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
