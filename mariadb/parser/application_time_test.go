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
