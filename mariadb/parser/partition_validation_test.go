package parser

import "testing"

// TestPartitionCountValidation covers general partition-clause leniencies that
// omni accepted but MariaDB 11.8.8 rejects with 1064 (verified across types):
// PARTITIONS without a count, an empty partition list, and a PARTITIONS count
// that mismatches the explicit partition list length.
func TestPartitionCountValidation(t *testing.T) {
	t.Run("reject", func(t *testing.T) {
		for _, sql := range []string{
			"CREATE TABLE t (x INT) PARTITION BY HASH(x) PARTITIONS",                 // no count
			"CREATE TABLE t (x INT) PARTITION BY HASH(x) ()",                         // empty list
			"CREATE TABLE t (x INT) PARTITION BY HASH(x) PARTITIONS 2 (PARTITION a)", // count > list
			"CREATE TABLE t (x INT) PARTITION BY RANGE(x) ()",                        // empty list (range)
		} {
			t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
		}
	})
	t.Run("accept", func(t *testing.T) {
		for _, sql := range []string{
			"CREATE TABLE t (x INT) PARTITION BY HASH(x) PARTITIONS 2",
			"CREATE TABLE t (x INT) PARTITION BY HASH(x) PARTITIONS 2 (PARTITION a, PARTITION b)",
			"CREATE TABLE t (x INT) PARTITION BY RANGE(x) (PARTITION p VALUES LESS THAN (10))",
		} {
			t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
		}
	})
}
