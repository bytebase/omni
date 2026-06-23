package parser

import "testing"

// TestSystemTimePartitioning covers PARTITION BY SYSTEM_TIME on a system-versioned
// table — all container-verified accepted by mariadb:11.8.8: bare, INTERVAL,
// LIMIT, PARTITIONS n, STARTS, AUTO, and explicit HISTORY/CURRENT partitions.
func TestSystemTimePartitioning(t *testing.T) {
	const sv = "CREATE TABLE t (x INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME"
	for _, sql := range []string{
		sv + " (PARTITION p0 HISTORY, PARTITION pn CURRENT)",
		sv + " INTERVAL 1 WEEK (PARTITION p0 HISTORY, PARTITION pn CURRENT)",
		sv + " LIMIT 1000 (PARTITION p0 HISTORY, PARTITION pn CURRENT)",
		sv + " INTERVAL 1 MONTH PARTITIONS 12",
		sv + " LIMIT 1000 PARTITIONS 5",
		sv + " INTERVAL 1 MONTH STARTS '2020-01-01 00:00:00' PARTITIONS 4",
		sv + " INTERVAL 1 MONTH AUTO",
		sv + " LIMIT 1000 AUTO",
		sv + " INTERVAL 1 MONTH STARTS '2020-01-01 00:00:00' AUTO",
	} {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSystemTimePartitioningReject pins the SYSTEM_TIME option grammar:
// INTERVAL value unit [STARTS expr] [AUTO] | LIMIT n [AUTO] | (no modifier).
// All container-verified 1064 vs mariadb:11.8.8.
func TestSystemTimePartitioningReject(t *testing.T) {
	const sv = "CREATE TABLE t (x INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME"
	for _, sql := range []string{
		sv + " INTERVAL 1 (PARTITION p0 HISTORY, PARTITION pn CURRENT)",                   // missing unit
		sv + " INTERVAL 1 FOOBAR (PARTITION p0 HISTORY, PARTITION pn CURRENT)",            // invalid unit
		sv + " STARTS '2020-01-01 00:00:00' (PARTITION p0 HISTORY, PARTITION pn CURRENT)", // bare STARTS
		sv + " AUTO", // bare AUTO
		sv + " LIMIT 1000 STARTS '2020-01-01 00:00:00' (PARTITION p0 HISTORY, PARTITION pn CURRENT)", // STARTS after LIMIT
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}
