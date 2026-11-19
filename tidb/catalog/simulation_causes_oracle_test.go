package catalog

import (
	"context"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// pingcap/tidb:v8.5.5 for a partition action on a table that has no partitions,
// the one cause those rows name: TiDB answers 1505 for every action that needs
// partitions, including ANALYZE PARTITION and the interval partitions FIRST and
// LAST name, which it takes on a partitioned table. Each case is answered by the
// pinned server and by the catalog, and both must agree with the code the case
// states.
func TestSimulationCausesOracle(t *testing.T) {
	conn, cleanup := startContainer(t)
	defer cleanup()
	ctx := context.Background()
	// A table named a is what every case leaves, and DROP VIEW fails on it, so
	// the reset drops tables only.
	reset := func(tb testing.TB) {
		for _, stmt := range []string{"DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS s"} {
			if _, err := conn.db.ExecContext(ctx, stmt); err != nil {
				tb.Errorf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(func() { reset(t) })

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	partitioned := "CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))"
	// The interval partitions FIRST and LAST name need the INTERVAL clause: on
	// a table that only ranges, TiDB answers 8200 for them.
	interval := "CREATE TABLE a (b INT) PARTITION BY RANGE (b) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)"
	cases := []struct {
		name  string
		setup []string
		stmt  string
		code  int
	}{
		{
			name:  "REMOVE PARTITIONING of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  1505,
		},
		{
			name:  "ADD PARTITION to a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (10))",
			code:  1505,
		},
		{
			name:  "ADD PARTITION PARTITIONS to a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION PARTITIONS 2",
			code:  1505,
		},
		{
			name:  "DROP PARTITION from a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP PARTITION p0",
			code:  1505,
		},
		{
			// IF EXISTS covers a missing partition, not a missing partitioning.
			name:  "DROP PARTITION IF EXISTS from a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP PARTITION IF EXISTS p0",
			code:  1505,
		},
		{
			name:  "TRUNCATE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TRUNCATE PARTITION ALL",
			code:  1505,
		},
		{
			name:  "COALESCE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COALESCE PARTITION 1",
			code:  1505,
		},
		{
			name:  "REORGANIZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)",
			code:  1505,
		},
		{
			name:  "EXCHANGE PARTITION with another table on a table with none",
			setup: []string{oneTable, "CREATE TABLE s (b INT)"},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s",
			code:  1505,
		},
		{
			// TiDB answers 1505 here where MySQL 8.0.46 answers 1066.
			name:  "EXCHANGE PARTITION with the table itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a",
			code:  1505,
		},
		{
			name:  "ANALYZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ANALYZE PARTITION p0",
			code:  1505,
		},
		{
			name:  "PARTITION ATTRIBUTES of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT",
			code:  1505,
		},
		{
			name:  "LAST PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a LAST PARTITION LESS THAN (10)",
			code:  1505,
		},
		{
			name:  "FIRST PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a FIRST PARTITION LESS THAN (20)",
			code:  1505,
		},
		{
			name:  "REMOVE PARTITIONING of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  0,
		},
		{
			name:  "ANALYZE PARTITION of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a ANALYZE PARTITION p0",
			code:  0,
		},
		{
			name:  "PARTITION ATTRIBUTES of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT",
			code:  0,
		},
		{
			name:  "LAST PARTITION of an interval-partitioned table",
			setup: []string{interval},
			stmt:  "ALTER TABLE a LAST PARTITION LESS THAN (1000)",
			code:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset(t)
			for _, stmt := range tc.setup {
				if _, err := conn.db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			_, err := conn.db.ExecContext(ctx, tc.stmt)
			engineCode, engineMessage := simulationEngineRefusal(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("tidb: %s answered %d (%s), want %d", tc.stmt, engineCode, engineMessage, tc.code)
			}

			c := New()
			execAll(t, c, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				execAll(t, c, stmt)
			}
			catalogCode, catalogMessage := simulationCatalogRefusal(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d (%s), tidb answered %d (%s)",
					tc.stmt, catalogCode, catalogMessage, engineCode, engineMessage)
			}
			if tc.code == 0 && tc.stmt == "ALTER TABLE a REMOVE PARTITIONING" {
				db := c.GetDatabase("d")
				if db == nil || db.GetTable("a") == nil || db.GetTable("a").Partitioning != nil {
					t.Fatalf("catalog kept the table partitioned after %s", tc.stmt)
				}
			}
		})
	}
}

// simulationEngineRefusal reads the server's answer as the error number and
// message, 0 for a statement it took.
func simulationEngineRefusal(t *testing.T, stmt string, err error) (int, string) {
	t.Helper()
	if err == nil {
		return 0, ""
	}
	var refusal *gomysql.MySQLError
	if !errors.As(err, &refusal) {
		t.Fatalf("%s: not a TiDB error: %v", stmt, err)
	}
	return int(refusal.Number), refusal.Message
}

// simulationCatalogRefusal reads the catalog's answer as the error number and
// message, 0 for a statement it took.
func simulationCatalogRefusal(t *testing.T, c *Catalog, stmt string) (int, string) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", stmt, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want 1", stmt, len(results))
	}
	if results[0].Error == nil {
		return 0, ""
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) {
		t.Fatalf("%s: not a catalog error: %v", stmt, results[0].Error)
	}
	return refusal.Code, refusal.Message
}
