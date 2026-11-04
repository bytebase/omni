package parser

import "testing"

// A partition-maintenance action in ALTER TABLE takes MySQL's
// opt_no_write_to_binlog right after the PARTITION keyword, spelled either
// NO_WRITE_TO_BINLOG or LOCAL, and the statement may end with one partitioning
// clause -- PARTITION BY ... or REMOVE PARTITIONING -- that follows the
// comma-separated action list without a comma of its own.
// https://dev.mysql.com/doc/refman/8.0/en/alter-table-partition-operations.html
//
// The corpus in testdata carries the accepted spellings TiDB's own parser test
// writes. What it cannot carry is the shape of the rejections, because pingcap
// never writes them: the modifier in the wrong place, two of them at once, a
// name list left empty by a modifier that swallowed the only identifier, and a
// comma in front of the partitioning clause. Those are below, and
// TestPartitionMaintenanceModifiersMatchTiDB asks the engine about every one,
// so the verdicts recorded here are TiDB v8.5.5's and not a memory of them.

var partitionMaintenanceModifiersAccepted = []string{
	// The modifier, in both spellings, on each action that takes one.
	`ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL p0, p1`,
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL ALL`,
	`ALTER TABLE pt REPAIR PARTITION NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt REPAIR PARTITION LOCAL ALL`,
	`ALTER TABLE pt REBUILD PARTITION LOCAL p0`,
	`ALTER TABLE pt REBUILD PARTITION NO_WRITE_TO_BINLOG ALL`,
	`ALTER TABLE pt COALESCE PARTITION LOCAL 3`,
	`ALTER TABLE pt COALESCE PARTITION NO_WRITE_TO_BINLOG 3`,
	`ALTER TABLE pt ADD PARTITION LOCAL`,
	`ALTER TABLE pt ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10`,
	`ALTER TABLE pt ADD PARTITION LOCAL (PARTITION p9 VALUES LESS THAN (100))`,
	`ALTER TABLE pt REORGANIZE PARTITION LOCAL`,
	`ALTER TABLE pt REORGANIZE PARTITION NO_WRITE_TO_BINLOG p0 INTO (PARTITION s0 VALUES LESS THAN (5))`,
	// ADD PARTITION and REORGANIZE PARTITION take nothing after the modifier.
	`ALTER TABLE pt ADD PARTITION`,
	`ALTER TABLE pt REORGANIZE PARTITION`,
	// LOCAL is also an ordinary partition name, so it reads as one once the
	// modifier before it has been taken, bare or quoted.
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL local`,
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL ` + "`local`",
	`ALTER TABLE pt REBUILD PARTITION LOCAL local, local`,
	// ...and so is REMOVE, which the partitioning clause does not claim here.
	`ALTER TABLE pt REORGANIZE PARTITION LOCAL remove INTO (PARTITION s0 VALUES LESS THAN (5))`,
	// One partitioning clause closes the statement, with no comma before it.
	`ALTER TABLE t REMOVE PARTITIONING`,
	`ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING`,
	`ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n) REMOVE PARTITIONING`,
	`ALTER TABLE t LOCK = DEFAULT REMOVE PARTITIONING`,
	`ALTER TABLE t MODIFY COLUMN n INT AFTER id REMOVE PARTITIONING`,
	`ALTER TABLE t CHANGE COLUMN n m INT AFTER id REMOVE PARTITIONING`,
	`ALTER TABLE t PARTITION BY HASH(id)`,
	`ALTER TABLE t ADD COLUMN a INT PARTITION BY HASH(id)`,
	`ALTER TABLE t COMMENT 'c' PARTITION BY HASH(id)`,
}

var partitionMaintenanceModifiersRejected = []string{
	// The modifier belongs after PARTITION, not before it and not after the
	// names. TiDB v8.5.5 reads opt_no_write_to_binlog on one side only, so the
	// trailing modifier MariaDB 11.8.8 takes on OPTIMIZE is a 1064 here.
	`ALTER TABLE pt OPTIMIZE LOCAL PARTITION p0`,
	`ALTER TABLE pt REBUILD LOCAL PARTITION p0`,
	`ALTER TABLE pt OPTIMIZE PARTITION p0 LOCAL`,
	`ALTER TABLE pt OPTIMIZE PARTITION ALL NO_WRITE_TO_BINLOG`,
	`ALTER TABLE pt REPAIR PARTITION p0 LOCAL`,
	`ALTER TABLE pt REBUILD PARTITION ALL NO_WRITE_TO_BINLOG`,
	`ALTER TABLE pt COALESCE PARTITION 3 LOCAL`,
	`ALTER TABLE pt ADD PARTITION PARTITIONS 10 LOCAL`,
	`ALTER TABLE pt ADD PARTITION (PARTITION p9 VALUES LESS THAN (100)) LOCAL`,
	`ALTER TABLE pt REORGANIZE PARTITION p0 LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))`,
	// One modifier, not two, whichever way they are spelled.
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG LOCAL p0`,
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL LOCAL p0`,
	`ALTER TABLE pt ADD PARTITION LOCAL NO_WRITE_TO_BINLOG PARTITIONS 10`,
	`ALTER TABLE pt COALESCE PARTITION LOCAL 3 LOCAL`,
	// The modifier is not a partition name standing in for the list, so a
	// single LOCAL or NO_WRITE_TO_BINLOG leaves the action without one.
	`ALTER TABLE pt OPTIMIZE PARTITION LOCAL`,
	`ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG`,
	`ALTER TABLE pt OPTIMIZE PARTITION local`,
	`ALTER TABLE pt OPTIMIZE PARTITION no_write_to_binlog`,
	`ALTER TABLE pt REPAIR PARTITION LOCAL`,
	`ALTER TABLE pt REBUILD PARTITION LOCAL`,
	`ALTER TABLE pt REORGANIZE PARTITION LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))`,
	// COALESCE PARTITION reads real_ulong_num, which is not optional.
	`ALTER TABLE pt COALESCE PARTITION`,
	`ALTER TABLE pt COALESCE PARTITION LOCAL`,
	`ALTER TABLE pt ADD PARTITION LOCAL PARTITIONS`,
	// ADD PARTITION names no partition and takes no ALL.
	`ALTER TABLE pt ADD PARTITION LOCAL ALL`,
	// The actions that read no modifier at all.
	`ALTER TABLE pt CHECK PARTITION LOCAL p0`,
	`ALTER TABLE pt CHECK PARTITION NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt TRUNCATE PARTITION LOCAL p0`,
	`ALTER TABLE pt TRUNCATE PARTITION NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt DROP PARTITION LOCAL p0`,
	`ALTER TABLE pt ANALYZE PARTITION LOCAL p0`,
	`ALTER TABLE pt ANALYZE PARTITION NO_WRITE_TO_BINLOG p0`,
	`ALTER TABLE pt ANALYZE PARTITION p0 LOCAL`,
	// The partitioning clause takes no comma before it, comes last, and comes
	// once. REMOVE on its own is not one.
	`ALTER TABLE t ADD COLUMN a INT, REMOVE PARTITIONING`,
	`ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n), REMOVE PARTITIONING`,
	`ALTER TABLE t REMOVE PARTITIONING, ADD COLUMN a INT`,
	`ALTER TABLE t REMOVE PARTITIONING ADD COLUMN a INT`,
	`ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING REMOVE PARTITIONING`,
	`ALTER TABLE t ADD COLUMN a INT, PARTITION BY HASH(id)`,
	`ALTER TABLE t PARTITION BY HASH(id) ADD COLUMN a INT`,
	`ALTER TABLE t PARTITION BY HASH(id), ADD COLUMN a INT`,
	`ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING PARTITION BY HASH(id)`,
	`ALTER TABLE t REMOVE PARTITIONING PARTITION BY HASH(id)`,
	`ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, ALGORITHM=INPLACE`,
	`ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, LOCK=NONE`,
	`ALTER TABLE t PARTITION BY HASH(id) ALGORITHM=INPLACE`,
	`ALTER TABLE t PARTITION BY HASH(id), ENABLE KEYS`,
	`ALTER TABLE t REMOVE PARTITIONING, ALGORITHM=INPLACE`,
	`ALTER TABLE t REMOVE`,
	`ALTER TABLE t REMOVE PARTITION`,
}

func TestPartitionMaintenanceModifiersAccepted(t *testing.T) {
	for _, sql := range partitionMaintenanceModifiersAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestPartitionMaintenanceModifiersRejected(t *testing.T) {
	for _, sql := range partitionMaintenanceModifiersRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.5 rejects it", sql)
			}
		})
	}
}

// TestPartitionMaintenanceModifiersMatchTiDB holds the two tables above to the
// engine itself (container test).
func TestPartitionMaintenanceModifiersMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for want, cases := range map[bool][]string{
		true:  partitionMaintenanceModifiersAccepted,
		false: partitionMaintenanceModifiersRejected,
	} {
		for _, sql := range cases {
			t.Run(sql, func(t *testing.T) {
				accepts, err := tc.prepareSyntax(sql)
				if err != nil {
					t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
				}
				if accepts != want {
					t.Errorf("TiDB %ss %q, this file records %s",
						verdictName(accepts), sql, verdictName(want))
				}
			})
		}
	}
}
