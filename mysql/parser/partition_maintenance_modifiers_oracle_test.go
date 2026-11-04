package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// A partition-maintenance action in ALTER TABLE takes opt_no_write_to_binlog
// right after the PARTITION keyword, spelled either NO_WRITE_TO_BINLOG or
// LOCAL, and the statement may end with one partitioning clause -- PARTITION
// BY ... or REMOVE PARTITIONING -- that follows the comma-separated action list
// without a comma of its own.
//
// Ref: opt_no_write_to_binlog and alter_commands in MySQL 8.0.46
// sql/sql_yacc.yy.
// https://dev.mysql.com/doc/refman/8.0/en/alter-table-partition-operations.html
//
// The rejections are the point of this file: the modifier in the wrong place,
// two of them at once, a name list left empty by a modifier that swallowed the
// only identifier, and a comma in front of the partitioning clause. The
// expected code pins the engine, so an image upgrade that changes a verdict
// fails here instead of passing silently.
var partitionModifierOracleCases = []alterDropCodeCase{
	// The modifier, in both spellings, on each action that takes one.
	{"ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0", 0},
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL p0, p1", 0},
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL ALL", 0},
	{"ALTER TABLE pt REPAIR PARTITION NO_WRITE_TO_BINLOG p0", 0},
	{"ALTER TABLE pt REPAIR PARTITION LOCAL ALL", 0},
	{"ALTER TABLE pt REBUILD PARTITION LOCAL p0", 0},
	{"ALTER TABLE pt REBUILD PARTITION NO_WRITE_TO_BINLOG ALL", 0},
	{"ALTER TABLE pt ANALYZE PARTITION LOCAL p0", 0},
	{"ALTER TABLE pt ANALYZE PARTITION NO_WRITE_TO_BINLOG p0", 0},
	{"ALTER TABLE pt COALESCE PARTITION LOCAL 3", 0},
	{"ALTER TABLE pt COALESCE PARTITION NO_WRITE_TO_BINLOG 3", 0},
	{"ALTER TABLE pt ADD PARTITION LOCAL", 0},
	{"ALTER TABLE pt ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10", 0},
	{"ALTER TABLE pt ADD PARTITION LOCAL (PARTITION p9 VALUES LESS THAN (100))", 0},
	{"ALTER TABLE pt REORGANIZE PARTITION LOCAL", 0},
	{"ALTER TABLE pt REORGANIZE PARTITION NO_WRITE_TO_BINLOG p0 INTO (PARTITION s0 VALUES LESS THAN (5))", 0},
	// ADD PARTITION and REORGANIZE PARTITION take nothing after the modifier.
	{"ALTER TABLE pt ADD PARTITION", 0},
	{"ALTER TABLE pt REORGANIZE PARTITION", 0},
	// LOCAL is also an ordinary partition name, so it reads as one once the
	// modifier before it has been taken, bare or quoted.
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL local", 0},
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL `local`", 0},
	{"ALTER TABLE pt REBUILD PARTITION LOCAL local, local", 0},
	// ...and so is REMOVE, which the partitioning clause does not claim here.
	{"ALTER TABLE pt REORGANIZE PARTITION LOCAL remove INTO (PARTITION s0 VALUES LESS THAN (5))", 0},
	// One partitioning clause closes the statement, with no comma before it.
	{"ALTER TABLE t REMOVE PARTITIONING", 0},
	{"ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING", 0},
	{"ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n) REMOVE PARTITIONING", 0},
	{"ALTER TABLE t LOCK = DEFAULT REMOVE PARTITIONING", 0},
	{"ALTER TABLE t MODIFY COLUMN n INT AFTER id REMOVE PARTITIONING", 0},
	{"ALTER TABLE t CHANGE COLUMN n m INT AFTER id REMOVE PARTITIONING", 0},
	{"ALTER TABLE t PARTITION BY HASH(id)", 0},
	{"ALTER TABLE t ADD COLUMN a INT PARTITION BY HASH(id)", 0},
	{"ALTER TABLE t COMMENT 'c' PARTITION BY HASH(id)", 0},

	// The modifier belongs after PARTITION, not before it and not after the
	// names. MySQL reads opt_no_write_to_binlog on one side only, so the
	// trailing modifier MariaDB 11.8.8 takes on OPTIMIZE is a 1064 here.
	{"ALTER TABLE pt OPTIMIZE LOCAL PARTITION p0", 1064},
	{"ALTER TABLE pt REBUILD LOCAL PARTITION p0", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION p0 LOCAL", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION ALL NO_WRITE_TO_BINLOG", 1064},
	{"ALTER TABLE pt REPAIR PARTITION p0 LOCAL", 1064},
	{"ALTER TABLE pt REBUILD PARTITION ALL NO_WRITE_TO_BINLOG", 1064},
	{"ALTER TABLE pt ANALYZE PARTITION p0 LOCAL", 1064},
	{"ALTER TABLE pt COALESCE PARTITION 3 LOCAL", 1064},
	{"ALTER TABLE pt ADD PARTITION PARTITIONS 10 LOCAL", 1064},
	{"ALTER TABLE pt ADD PARTITION (PARTITION p9 VALUES LESS THAN (100)) LOCAL", 1064},
	{"ALTER TABLE pt REORGANIZE PARTITION p0 LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))", 1064},
	// One modifier, not two, whichever way they are spelled.
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL NO_WRITE_TO_BINLOG p0", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG LOCAL p0", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL LOCAL p0", 1064},
	{"ALTER TABLE pt ADD PARTITION LOCAL NO_WRITE_TO_BINLOG PARTITIONS 10", 1064},
	{"ALTER TABLE pt COALESCE PARTITION LOCAL 3 LOCAL", 1064},
	// The modifier is not a partition name standing in for the list, so a
	// single LOCAL or NO_WRITE_TO_BINLOG leaves the action without one.
	{"ALTER TABLE pt OPTIMIZE PARTITION LOCAL", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION local", 1064},
	{"ALTER TABLE pt OPTIMIZE PARTITION no_write_to_binlog", 1064},
	{"ALTER TABLE pt REPAIR PARTITION LOCAL", 1064},
	{"ALTER TABLE pt REBUILD PARTITION LOCAL", 1064},
	{"ALTER TABLE pt REORGANIZE PARTITION LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))", 1064},
	// COALESCE PARTITION reads real_ulong_num, which is not optional.
	{"ALTER TABLE pt COALESCE PARTITION", 1064},
	{"ALTER TABLE pt COALESCE PARTITION LOCAL", 1064},
	{"ALTER TABLE pt ADD PARTITION LOCAL PARTITIONS", 1064},
	// ADD PARTITION names no partition and takes no ALL.
	{"ALTER TABLE pt ADD PARTITION LOCAL ALL", 1064},
	// The actions that read no modifier at all.
	{"ALTER TABLE pt CHECK PARTITION LOCAL p0", 1064},
	{"ALTER TABLE pt CHECK PARTITION NO_WRITE_TO_BINLOG p0", 1064},
	{"ALTER TABLE pt TRUNCATE PARTITION LOCAL p0", 1064},
	{"ALTER TABLE pt TRUNCATE PARTITION NO_WRITE_TO_BINLOG p0", 1064},
	{"ALTER TABLE pt DROP PARTITION LOCAL p0", 1064},
	// The partitioning clause takes no comma before it, comes last, and comes
	// once. REMOVE on its own is not one.
	{"ALTER TABLE t ADD COLUMN a INT, REMOVE PARTITIONING", 1064},
	{"ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n), REMOVE PARTITIONING", 1064},
	{"ALTER TABLE t REMOVE PARTITIONING, ADD COLUMN a INT", 1064},
	{"ALTER TABLE t REMOVE PARTITIONING ADD COLUMN a INT", 1064},
	{"ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING REMOVE PARTITIONING", 1064},
	{"ALTER TABLE t ADD COLUMN a INT, PARTITION BY HASH(id)", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id) ADD COLUMN a INT", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id), ADD COLUMN a INT", 1064},
	{"ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING PARTITION BY HASH(id)", 1064},
	{"ALTER TABLE t REMOVE PARTITIONING PARTITION BY HASH(id)", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, ALGORITHM=INPLACE", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, LOCK=NONE", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id) ALGORITHM=INPLACE", 1064},
	{"ALTER TABLE t PARTITION BY HASH(id), ENABLE KEYS", 1064},
	{"ALTER TABLE t REMOVE PARTITIONING, ALGORITHM=INPLACE", 1064},
	{"ALTER TABLE t REMOVE", 1064},
	{"ALTER TABLE t REMOVE PARTITION", 1064},
}

// TestPartitionMaintenanceModifiersOracle asserts omni agrees with a live
// mysql:8.0 on partitionModifierOracleCases (container test). MySQL's verdict
// comes from PREPARE, which parses without running the statement.
func TestPartitionMaintenanceModifiersOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, ddl := range []string{
		"CREATE TABLE t (id INT PRIMARY KEY, n INT)",
		"CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	for _, tc := range partitionModifierOracleCases {
		var code uint16
		_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(tc.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", tc.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != tc.code {
			t.Errorf("MySQL code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, MySQL code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
