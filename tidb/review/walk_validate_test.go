package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

func TestWalkThroughTiDBAlterSequence(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Valid ALTER SEQUENCE on an existing sequence is accepted and keeps walking.
	noFinding(t, "ALTER SEQUENCE sq INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE;", shop)
	noFinding(t, "ALTER SEQUENCE sq RESTART WITH 5;\nALTER SEQUENCE sq NOMINVALUE NOMAXVALUE NOCACHE NOCYCLE;", shop)
	oneFinding(t, "ALTER SEQUENCE sq INCREMENT 2;\nCREATE SEQUENCE sq;", 1, "Sequence `sq` already exists.", shop)
	// Missing sequence (even with IF EXISTS on TiDB v8.5.5) and wrong object kind.
	oneFinding(t, "ALTER SEQUENCE nosuch INCREMENT 1;", 0, "Sequence `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER SEQUENCE IF EXISTS nosuch INCREMENT 1;", 0, "Sequence `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER SEQUENCE t INCREMENT 1;", 0, "`t` is a table, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE v NOMINVALUE;", 0, "`v` is a view, not a sequence.", shop)
	// Temporary table does not shadow a missing sequence.
	oneFinding(t, "CREATE TEMPORARY TABLE tmp (a int);\nALTER SEQUENCE tmp INCREMENT 1;", 1, "Sequence `tmp` doesn't exist.", shop)
	// Self-contradictory sequence options (4136).
	oneFinding(t, "ALTER SEQUENCE sq CACHE 0;", 0, "Sequence `sq` has conflicting options.", shop)
	oneFinding(t, "ALTER SEQUENCE sq INCREMENT 0;", 0, "Sequence `sq` has conflicting options.", shop)
	oneFinding(t, "CREATE SEQUENCE s2 CACHE 0;", 0, "Sequence `s2` has conflicting options.", shop)
	oneFinding(t, "CREATE SEQUENCE s2 INCREMENT 0;", 0, "Sequence `s2` has conflicting options.", shop)
	oneFinding(t, "CREATE SEQUENCE s2 INCREMENT 1 INCREMENT 0;", 0, "Sequence `s2` has conflicting options.", shop)
	oneFinding(t, "CREATE SEQUENCE s2 MINVALUE 10 MAXVALUE 5;", 0, "Sequence `s2` has conflicting options.", shop)
	noFinding(t, "CREATE SEQUENCE s2 CACHE 10 INCREMENT 2;", shop)
	// Repeated INCREMENT/MINVALUE/MAXVALUE options apply in order; a later valid
	// option overrides an earlier invalid one (while CACHE <= 0 is refused per option).
	noFinding(t, "CREATE SEQUENCE s2 INCREMENT 0 INCREMENT 1;", shop)
	noFinding(t, "ALTER SEQUENCE sq INCREMENT 0 INCREMENT 1;", shop)
	oneFinding(t, "CREATE SEQUENCE s2 CACHE 0 CACHE 10;", 0, "Sequence `s2` has conflicting options.", shop)
	// When creation is uncertain (IF NOT EXISTS under OtherSpecs), skip option validation.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE SEQUENCE IF NOT EXISTS sq CACHE 0;", other)
}

func TestWalkThroughTiDBUnsupportedAlterActionsAndOptions(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"ALTER TABLE t FORCE;", "TiDB doesn't support `FORCE` in `ALTER TABLE`."},
		{"ALTER TABLE t CHECK PARTITION ALL;", "TiDB doesn't support `CHECK PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t OPTIMIZE PARTITION ALL;", "TiDB doesn't support `OPTIMIZE PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t REBUILD PARTITION ALL;", "TiDB doesn't support `REBUILD PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t REPAIR PARTITION ALL;", "TiDB doesn't support `REPAIR PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t DISCARD TABLESPACE;", "TiDB doesn't support `DISCARD TABLESPACE` in `ALTER TABLE`."},
		{"ALTER TABLE t IMPORT TABLESPACE;", "TiDB doesn't support `IMPORT TABLESPACE` in `ALTER TABLE`."},
		{"ALTER TABLE t DISCARD PARTITION ALL TABLESPACE;", "TiDB doesn't support `DISCARD PARTITION TABLESPACE` in `ALTER TABLE`."},
		{"ALTER TABLE t IMPORT PARTITION ALL TABLESPACE;", "TiDB doesn't support `IMPORT PARTITION TABLESPACE` in `ALTER TABLE`."},
		{"ALTER TABLE t SECONDARY_LOAD;", "TiDB doesn't support `SECONDARY_LOAD` in `ALTER TABLE`."},
		{"ALTER TABLE t SECONDARY_UNLOAD;", "TiDB doesn't support `SECONDARY_UNLOAD` in `ALTER TABLE`."},
		{"ALTER TABLE t STATS_OPTIONS = DEFAULT;", "TiDB doesn't support `STATS_OPTIONS` in `ALTER TABLE`."},
		{"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (100);", "TiDB doesn't support `SPLIT MAXVALUE PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t MERGE FIRST PARTITION LESS THAN (100);", "TiDB doesn't support `MERGE FIRST PARTITION` in `ALTER TABLE`."},
		{"ALTER TABLE t UNION = ();", "TiDB doesn't support table option `UNION`."},
		{"ALTER TABLE t INSERT_METHOD = LAST;", "TiDB doesn't support table option `INSERT_METHOD`."},
		{"ALTER TABLE t ENGINE = 'bogus';", "TiDB doesn't recognize storage engine `bogus`."},
		{"ALTER TABLE t AFFINITY = 'bogus';", "TiDB doesn't recognize `AFFINITY` value `bogus`."},
		{"ALTER TABLE t AVG_ROW_LENGTH = 10;", "TiDB doesn't support table option `AVG_ROW_LENGTH` in `ALTER TABLE`."},
		{"ALTER TABLE t STATS_SAMPLE_RATE = 1.0;", "TiDB doesn't support table option `STATS_SAMPLE_RATE` in `ALTER TABLE`."},
		{"ALTER TABLE t STORAGE DISK;", "TiDB doesn't support table option `STORAGE` in `ALTER TABLE`."},
		{"ALTER TABLE t CONVERT TO CHARSET DEFAULT;", "TiDB doesn't support `CONVERT TO CHARSET DEFAULT`."},
		{"ALTER TABLE t CHARSET = ucs2;", "TiDB doesn't recognize character set `ucs2`."},
		{"ALTER TABLE t ATTRIBUTES = 'not_kv';", "TiDB requires `ATTRIBUTES` entries in `key=value` format without conflicts."},
		{"ALTER TABLE t ATTRIBUTES = 'k=1,k=2';", "TiDB requires `ATTRIBUTES` entries in `key=value` format without conflicts."},
		{"ALTER TABLE t PARTITION p0 COMMENT = 'x';", "TiDB doesn't support partition option `COMMENT` in `ALTER TABLE`."},
		{"CREATE TABLE x (a int) ENGINE = 'bogus';", "TiDB doesn't recognize storage engine `bogus`."},
		{"CREATE TABLE x (a int) UNION = ();", "TiDB doesn't support table option `UNION`."},
		{"CREATE TABLE x (a int) INSERT_METHOD = FIRST;", "TiDB doesn't support table option `INSERT_METHOD`."},
	} {
		t.Run(tc.sql, func(t *testing.T) { oneFinding(t, tc.sql, 0, tc.want, shop) })
	}
	// Valid controls.
	for _, sql := range []string{
		"ALTER TABLE t ENGINE = InnoDB, ROW_FORMAT = DYNAMIC, COMMENT = 'ok', AUTO_INCREMENT = 10;",
		"ALTER TABLE t FORCE AUTO_INCREMENT = 5;",
		"ALTER TABLE t AFFINITY = 'table';",
		"ALTER TABLE t ATTRIBUTES = 'merge_option=allow,k=v';",
		"ALTER TABLE t CHARSET = utf8mb4;",
		"ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;",
		"ALTER TABLE t TTL = a + INTERVAL 1 DAY;",
		"CREATE TABLE x (a int) ENGINE = InnoDB CHARSET = utf8mb4 AFFINITY = 'none';",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

func TestWalkThroughTiDBColumnValidation(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"ALTER TABLE t ADD COLUMN c int AUTO_RANDOM;", "TiDB doesn't support adding column `c` with `AUTO_RANDOM`."},
		{"ALTER TABLE t ADD COLUMN c int AUTO_INCREMENT;", "TiDB doesn't support adding column `c` with `AUTO_INCREMENT`."},
		{"ALTER TABLE t ADD COLUMN c SERIAL;", "TiDB doesn't support adding column `c` with `AUTO_INCREMENT`."},
		{"ALTER TABLE t ADD COLUMN c int UNIQUE;", "TiDB doesn't support adding column `c` with `UNIQUE`."},
		{"ALTER TABLE nokey ADD COLUMN c int PRIMARY KEY;", "TiDB doesn't support adding column `c` with `PRIMARY KEY`."},
		{"ALTER TABLE t ADD COLUMN c int AS (a + 1) STORED;", "TiDB doesn't support a `STORED` generated column in `ALTER TABLE`."},
		{"ALTER TABLE t ADD COLUMN c tinyint DEFAULT CURRENT_TIMESTAMP;", "`c` has an invalid default value for its type."},
		{"ALTER TABLE t ADD COLUMN c datetime(3) DEFAULT CURRENT_TIMESTAMP(6);", "`c` has an invalid default value for its type."},
		{"ALTER TABLE t ADD COLUMN c tinyint DEFAULT CURDATE();", "`c` has an invalid default value for its type."},
		{"ALTER TABLE t ADD COLUMN c tinyint ON UPDATE CURRENT_TIMESTAMP;", "`c` has an invalid `ON UPDATE` clause for its type."},
		{"ALTER TABLE t ADD COLUMN c datetime(3) ON UPDATE CURRENT_TIMESTAMP(6);", "`c` has an invalid `ON UPDATE` clause for its type."},
		{"ALTER TABLE t ADD COLUMN c ENUM('a', 'a');", "`c` declares duplicate ENUM value `a`."},
		{"ALTER TABLE t ADD COLUMN c SET('x', 'x');", "`c` declares duplicate SET value `x`."},
		{"ALTER TABLE t ADD COLUMN c TINYTEXT UNICODE;", "TiDB doesn't recognize character set `ucs2`."},
		{"ALTER TABLE t ADD COLUMN c int DEFAULT my_func();", "TiDB doesn't allow `MY_FUNC` in a column `DEFAULT` expression."},
		{"ALTER TABLE t ADD COLUMN c int DEFAULT NEXT VALUE FOR nosuch;", "Sequence `nosuch` doesn't exist."},
		{"ALTER TABLE t ADD COLUMN c int DEFAULT NEXT VALUE FOR r;", "`r` is a table, not a sequence."},
		{"ALTER TABLE t ADD COLUMN c int DEFAULT NEXT VALUE FOR sq;", "TiDB doesn't support a sequence default in `ALTER TABLE ADD COLUMN`."},
		{"ALTER TABLE t ADD COLUMN c int AS (nosuch + 1) VIRTUAL;", "Column `nosuch` doesn't exist in `t`."},
		{"CREATE TABLE x (c int DEFAULT CURRENT_TIMESTAMP);", "`c` has an invalid default value for its type."},
	} {
		t.Run(tc.sql, func(t *testing.T) { oneFinding(t, tc.sql, 0, tc.want, shop) })
	}
	// Valid controls.
	for _, sql := range []string{
		"ALTER TABLE t ADD COLUMN c datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;",
		"ALTER TABLE t ADD COLUMN c datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);",
		"ALTER TABLE t ADD COLUMN c date DEFAULT (CURDATE());",
		"ALTER TABLE t ADD COLUMN c ENUM('a', 'b');",
		"ALTER TABLE t ADD COLUMN c int AS (a + 1) VIRTUAL;",
		"CREATE TABLE x (id bigint PRIMARY KEY AUTO_RANDOM, u int UNIQUE, g int AS (u + 1) STORED);",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

func TestWalkThroughTiDBIndexesTablesViewsFKsAndPartitions(t *testing.T) {
	shop := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	for _, tc := range []struct {
		sql  string
		want string
	}{
		// Indexes.
		{"CREATE FULLTEXT INDEX ft ON t (a);", "TiDB doesn't support FULLTEXT indexes."},
		{"CREATE SPATIAL INDEX sp ON t (a);", "TiDB doesn't support SPATIAL indexes."},
		{"CREATE FULLTEXT INDEX ft ON v (a);", "TiDB doesn't support FULLTEXT indexes."},
		{"CREATE INDEX idx_v ON v (id);", "TiDB doesn't support column indexes on view `v`."},
		{"CREATE INDEX idx_missing ON t (nosuch);", "Column `nosuch` doesn't exist in `t`."},
		{"ALTER TABLE t ADD INDEX idx_missing (nosuch);", "Column `nosuch` doesn't exist in `t`."},
		{"ALTER TABLE nokey ADD PRIMARY KEY (nosuch);", "Column `nosuch` doesn't exist in `nokey`."},
		// CREATE TABLE / CREATE VIEW / FOREIGN KEY / RENAME TABLE.
		{"CREATE TABLE x;", "`CREATE TABLE` must define at least one column."},
		{"CREATE TEMPORARY TABLE tmp;", "`CREATE TABLE` must define at least one column."},
		{"CREATE TABLE x AS SELECT 1 AS a;", "TiDB doesn't support `CREATE TABLE ... SELECT`."},
		{"CREATE VIEW w AS SELECT *;", "`CREATE VIEW` expands `*` with no table to read from."},
		{"CREATE OR REPLACE VIEW v AS SELECT *;", "`CREATE VIEW` expands `*` with no table to read from."},
		{"CREATE VIEW w AS SELECT x.* FROM t;", "`CREATE VIEW` expands `*` with no table to read from."},
		{"CREATE VIEW w AS SELECT other.t.* FROM t;", "`CREATE VIEW` expands `*` with no table to read from."},
		{"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES r;", "`FOREIGN KEY` must name the referenced columns."},
		{"CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES r);", "`FOREIGN KEY` must name the referenced columns."},
		{"RENAME TABLE t TO t, t TO t;", "Table `t` already exists."},
		{"RENAME TABLE nosuch TO nosuch, nosuch TO nosuch;", "`nosuch` doesn't exist."},
		// Partitions.
		{"ALTER TABLE t PARTITION BY HASH (nosuch);", "Column `nosuch` doesn't exist in `t`."},
		{"ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (10);", "TiDB requires `FIRST` and `LAST` partitions in `INTERVAL` partitioning."},
		{"ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (a) LAST PARTITION LESS THAN (100);", "Partition bound expression can't reference column `a`."},
		{"ALTER TABLE pt FIRST PARTITION LESS THAN (a);", "Partition bound expression can't reference column `a`."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p1 HISTORY);", "TiDB doesn't support `HISTORY` partitions."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p1 CURRENT);", "TiDB doesn't support `CURRENT` partitions."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p1 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20));", "This statement declares partition `p1` twice."},
		// TiDB reads the missing VALUES LESS THAN first (1479, verified).
		{"ALTER TABLE pt ADD PARTITION (PARTITION p1 (SUBPARTITION s1, SUBPARTITION s1));", "Partition `p1` needs `VALUES LESS THAN` for `RANGE` partitioning."},
		{"ALTER TABLE pt REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10) (SUBPARTITION s1, SUBPARTITION s1));", "This statement declares subpartition `s1` twice."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p1 VALUES LESS THAN (a));", "Partition bound expression can't reference column `a`."},
		{"ALTER TABLE pt ANALYZE PARTITION p0 INDEX nosuch;", "Index `nosuch` doesn't exist in `pt`."},
		{"ALTER TABLE pt EXCHANGE PARTITION p0 WITH TABLE pt;", "TiDB can't exchange a partition of `pt` with itself."},
		{"ALTER TABLE pt EXCHANGE PARTITION p0 WITH TABLE nosuch;", "Table `nosuch` doesn't exist."},
		{"ALTER TABLE pt EXCHANGE PARTITION p0 WITH TABLE v;", "`v` is a view, not a table."},
	} {
		t.Run(tc.sql, func(t *testing.T) { oneFinding(t, tc.sql, 0, tc.want, shop) })
	}
	// Valid controls.
	for _, sql := range []string{
		"CREATE INDEX idx_a ON t (a);",
		"CREATE INDEX idx_expr ON v ((1 + 1));",
		"ALTER TABLE t ADD INDEX idx_a (a);",
		"CREATE VIEW w AS SELECT * FROM t;",
		"CREATE VIEW w AS SELECT t.* FROM t;",
		"CREATE VIEW w AS SELECT shop.t.* FROM t;",
		"CREATE VIEW w AS SELECT t.* FROM shop.t;",
		"CREATE VIEW w AS SELECT 1 AS x;",
		"ALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES r (id);",
		"RENAME TABLE t TO tmp, r TO t, tmp TO r;",
		"ALTER TABLE t PARTITION BY HASH (a) PARTITIONS 4;",
		"ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (100);",
		"ALTER TABLE pt ADD PARTITION (PARTITION p1 VALUES LESS THAN (100), PARTITION p2 VALUES LESS THAN MAXVALUE);",
		"ALTER TABLE pt ANALYZE PARTITION p0 INDEX i;",
		"ALTER TABLE pt EXCHANGE PARTITION p0 WITH TABLE nokey;",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}
