package catalog

import (
	"errors"
	"testing"
)

func TestDDLValidation(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		execAll(t, c, `
			CREATE DATABASE db;
			USE db;
			CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d));
			CREATE TABLE tp (b INT, c INT, PRIMARY KEY (b, c))
				PARTITION BY RANGE (c) (
					PARTITION p0 VALUES LESS THAN (100),
					PARTITION p1 VALUES LESS THAN (200),
					PARTITION p2 VALUES LESS THAN (300)
				);
			CREATE TABLE t_nopk (a INT);
			CREATE VIEW v AS SELECT b FROM t;
		`)
		return c
	}

	assertCode := func(t *testing.T, c *Catalog, sql string, wantCode int) {
		t.Helper()
		res, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("Exec(%q) parse err: %v, want code %d", sql, err, wantCode)
		}
		if len(res) != 1 || res[0].Error == nil {
			t.Fatalf("Exec(%q) succeeded, want error code %d", sql, wantCode)
		}
		var catErr *Error
		if !errors.As(res[0].Error, &catErr) {
			t.Fatalf("Exec(%q) err = %T (%v), want *Error %d", sql, res[0].Error, res[0].Error, wantCode)
		}
		if catErr.Code != wantCode {
			t.Errorf("Exec(%q) code = %d (%s), want %d", sql, catErr.Code, catErr.Message, wantCode)
		}
	}

	t.Run("ctas_unsupported_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE TABLE fresh SELECT *", ErrUnknownError)
		// Fires before duplicate-table 1050.
		assertCode(t, c, "CREATE TABLE t SELECT *", ErrUnknownError)
		assertCode(t, c, "CREATE TABLE t WITH cte AS (SELECT *) SELECT *", ErrUnknownError)
	})

	t.Run("preprocess_table_options_1286_8232_8233", func(t *testing.T) {
		c := setup(t)
		// Unknown engine fails before missing table (1146) or view (1347).
		assertCode(t, c, "ALTER TABLE missing ENGINE 'nope'", ErrUnknownStorageEngine)
		assertCode(t, c, "ALTER TABLE v STORAGE ENGINE 'nope'", ErrUnknownStorageEngine)
		assertCode(t, c, "ALTER TABLE t ENGINE 'TiFlash'", ErrUnknownStorageEngine)
		assertCode(t, c, "CREATE TABLE u (a INT) ENGINE = 'nope'", ErrUnknownStorageEngine)
		// Valid engines succeed.
		execAll(t, c, "ALTER TABLE t ENGINE = InnoDB; ALTER TABLE t STORAGE ENGINE = MyISAM;")

		// UNION and INSERT_METHOD fail before table lookup.
		assertCode(t, c, "ALTER TABLE missing UNION ()", ErrTableOptionUnionUnsupported)
		assertCode(t, c, "ALTER TABLE v UNION ()", ErrTableOptionUnionUnsupported)
		assertCode(t, c, "ALTER TABLE missing INSERT_METHOD 'FIRST'", ErrTableOptionInsertMethodUnsupported)
		assertCode(t, c, "ALTER TABLE v INSERT_METHOD 'FIRST'", ErrTableOptionInsertMethodUnsupported)
	})

	t.Run("add_column_inline_constraints_8200", func(t *testing.T) {
		c := setup(t)
		// Fires ahead of duplicate column (1060) and multiple PK (1068).
		assertCode(t, c, "ALTER TABLE t ADD e TINYINT PRIMARY KEY", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD c TINYINT PRIMARY KEY", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD e TINYINT UNIQUE", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD c TINYINT UNIQUE KEY", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD e TINYINT AUTO_INCREMENT", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD e SERIAL", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD e TINYINT SERIAL DEFAULT VALUE", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t ADD COLUMN IF NOT EXISTS c INT PRIMARY KEY", ErrUnsupportedDDL)
	})

	t.Run("modify_primary_key_column_type_8200", func(t *testing.T) {
		c := setup(t)
		// Narrowing a PK column fails with 8200.
		assertCode(t, c, "ALTER TABLE tp MODIFY c TINYINT", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE tp CHANGE c c TINYINT", ErrUnsupportedDDL)
		// Changing signedness of a PK column fails with 8200.
		assertCode(t, c, "ALTER TABLE t MODIFY b INT UNSIGNED", ErrUnsupportedDDL)
		// Widening a PK integer column or modifying a non-PK column succeeds.
		execAll(t, c, "ALTER TABLE t MODIFY b BIGINT; ALTER TABLE t MODIFY c TINYINT;")
	})

	t.Run("alter_charset_compatibility_8200", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t CHARSET 'binary'", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t COLLATE 'binary'", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t CONVERT TO CHARSET 'binary'", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE t CHARSET 'latin1'", ErrUnsupportedDDL)
		execAll(t, c, "ALTER TABLE t CHARSET 'utf8mb4' COLLATE 'utf8mb4_general_ci';")
	})

	t.Run("unsupported_alter_specs_and_options_8200", func(t *testing.T) {
		c := setup(t)
		for _, sql := range []string{
			"ALTER TABLE t CHECK PARTITION ALL",
			"ALTER TABLE t OPTIMIZE PARTITION ALL",
			"ALTER TABLE t REBUILD PARTITION ALL",
			"ALTER TABLE t REPAIR PARTITION ALL",
			"ALTER TABLE t DISCARD TABLESPACE",
			"ALTER TABLE t IMPORT TABLESPACE",
			"ALTER TABLE t FORCE",
			"ALTER TABLE t SECONDARY_LOAD",
			"ALTER TABLE t SECONDARY_UNLOAD",
			"ALTER TABLE t STATS_OPTIONS 'x'",
			"ALTER TABLE tp SPLIT MAXVALUE PARTITION LESS THAN (400)",
			"ALTER TABLE tp MERGE FIRST PARTITION LESS THAN (50)",
			"ALTER TABLE t AVG_ROW_LENGTH 1",
			"ALTER TABLE t CHECKSUM 1",
			"ALTER TABLE t COMPRESSION 'ZLIB'",
			"ALTER TABLE t KEY_BLOCK_SIZE 4",
			"ALTER TABLE t PRE_SPLIT_REGIONS 2",
			"ALTER TABLE t SECONDARY_ENGINE NULL",
			"ALTER TABLE t STATS_PERSISTENT 1",
			"ALTER TABLE t STORAGE DISK",
			"ALTER TABLE t TABLESPACE ts1",
		} {
			assertCode(t, c, sql, ErrUnsupportedDDL)
		}
	})

	t.Run("shard_row_id_bits_8200", func(t *testing.T) {
		c := setup(t)
		// Clustered PK table rejects SHARD_ROW_ID_BITS > 0.
		assertCode(t, c, "ALTER TABLE t SHARD_ROW_ID_BITS 4", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE tp SHARD_ROW_ID_BITS 4", ErrUnsupportedDDL)
		// Non-PK table (uses _tidb_rowid) accepts SHARD_ROW_ID_BITS.
		execAll(t, c, "ALTER TABLE t_nopk SHARD_ROW_ID_BITS 4; ALTER TABLE t SHARD_ROW_ID_BITS 0;")
	})

	t.Run("reorganize_partition_range_adjacency_8200", func(t *testing.T) {
		c := setup(t)
		// Bare REORGANIZE PARTITION and non-adjacent RANGE partitions fail with 8200.
		assertCode(t, c, "ALTER TABLE tp REORGANIZE PARTITION", ErrUnsupportedDDL)
		assertCode(t, c, "ALTER TABLE tp REORGANIZE PARTITION p0, p2 INTO (PARTITION p02 VALUES LESS THAN (300))", ErrUnsupportedDDL)
		// Adjacent RANGE partitions succeed.
		execAll(t, c, "ALTER TABLE tp REORGANIZE PARTITION p0, p1 INTO (PARTITION p01 VALUES LESS THAN (200));")
	})

	t.Run("auto_id_cache_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t AUTO_ID_CACHE 1", ErrUnknownError)
		execAll(t, c, "ALTER TABLE t AUTO_ID_CACHE 100; ALTER TABLE t AUTO_ID_CACHE 0;")
		execAll(t, c, "CREATE TABLE t_cache1 (a INT) AUTO_ID_CACHE = 1;")
		assertCode(t, c, "ALTER TABLE t_cache1 AUTO_ID_CACHE 0", ErrUnknownError)
		execAll(t, c, "ALTER TABLE t_cache1 AUTO_ID_CACHE 1;")
	})

	t.Run("extended_stats_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t ADD STATS_EXTENDED s1 CARDINALITY(b, c)", ErrUnknownError)
		assertCode(t, c, "ALTER TABLE t DROP STATS_EXTENDED s1", ErrUnknownError)
		execAll(t, c, "SET tidb_enable_extended_stats = ON;")
		execAll(t, c, "ALTER TABLE t ADD STATS_EXTENDED s1 CARDINALITY(b, c); ALTER TABLE t DROP STATS_EXTENDED s1;")
	})

	t.Run("compact_partition_and_partition_options_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t COMPACT PARTITION p0", ErrUnknownError)
		assertCode(t, c, "ALTER TABLE v COMPACT PARTITION p0 TIFLASH REPLICA", ErrUnknownError)
		execAll(t, c, "ALTER TABLE t COMPACT; ALTER TABLE v COMPACT; ALTER TABLE tp COMPACT PARTITION p0;")
		assertCode(t, c, "ALTER TABLE tp PARTITION p0 COMMENT 'x'", ErrUnknownError)
	})

	t.Run("tiflash_server_count_1105", func(t *testing.T) {
		c := setup(t)
		// Default unconstrained (-1) accepts any replica count.
		execAll(t, c, "ALTER TABLE t SET TIFLASH REPLICA 2;")
		// With 0 TiFlash servers, replica > 0 fails with 1105 while 0 and HYPO succeed.
		c.SetTiFlashServerCount(0)
		assertCode(t, c, "ALTER TABLE t SET TIFLASH REPLICA 1", ErrUnknownError)
		execAll(t, c, "ALTER TABLE t SET TIFLASH REPLICA 0; ALTER TABLE t SET HYPO TIFLASH REPLICA 1;")
	})

	t.Run("attributes_8237_and_affinity_8266", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t ATTRIBUTES 'binary'", ErrInvalidAttributes)
		assertCode(t, c, "ALTER TABLE t ATTRIBUTES 'k=a,k=b'", ErrInvalidAttributes)
		assertCode(t, c, "ALTER TABLE tp PARTITION p0 ATTRIBUTES 'binary'", ErrInvalidAttributes)
		execAll(t, c, "ALTER TABLE t ATTRIBUTES 'k=a,k=a'; ALTER TABLE t ATTRIBUTES DEFAULT; ALTER TABLE t ATTRIBUTES '';")

		assertCode(t, c, "ALTER TABLE t AFFINITY 'binary'", ErrInvalidTableAffinity)
		execAll(t, c, "ALTER TABLE t AFFINITY 'table'; ALTER TABLE t AFFINITY 'partition'; ALTER TABLE t AFFINITY 'none'; ALTER TABLE t AFFINITY '';")
	})

	t.Run("curdate_default_small_int_1067", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "ALTER TABLE t ADD e TINYINT DEFAULT CURDATE()", ErrInvalidDefault)
		assertCode(t, c, "ALTER TABLE t ADD e SMALLINT DEFAULT CURRENT_DATE", ErrInvalidDefault)
		assertCode(t, c, "ALTER TABLE t ADD e MEDIUMINT DEFAULT CURRENT_DATE()", ErrInvalidDefault)
		// Duplicate column 1060 still fires before default value evaluation.
		assertCode(t, c, "ALTER TABLE t ADD c TINYINT DEFAULT CURRENT_DATE", ErrDupColumn)
		// INT / BIGINT / DATE columns accept CURDATE() default.
		execAll(t, c, "ALTER TABLE t ADD e INT DEFAULT CURDATE(); ALTER TABLE t ADD f DATE DEFAULT CURRENT_DATE;")
	})
}
