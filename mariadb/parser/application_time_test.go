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
