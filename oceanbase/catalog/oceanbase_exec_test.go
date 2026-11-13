package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// What OceanBase answers to DDL the catalog simulates, where its answer is not MySQL's.
// Each statement below is held to the pinned server and to the catalog, with the error
// number both must answer, 0 for a statement they take. The catalog's rules cite these
// tests.

// engineDatabase opens t's connection to a database of its own on OceanBase, empty, and
// a catalog with a database of the same name selected. The name is the server's to give
// (testcontainer.DatabaseName), so a test reads its schema as DATABASE().
func engineDatabase(t *testing.T, label string) (*sql.DB, *Catalog, string) {
	t.Helper()
	db, name := testcontainer.OpenDatabase(t, label), testcontainer.DatabaseName(label)
	c := New()
	if code := omniCode(c, "CREATE DATABASE "+name+"; USE "+name); code != 0 {
		t.Fatalf("omni: CREATE DATABASE %s: error %d", name, code)
	}
	return db, c, name
}

// runExecSteps runs each case's statements in a database of its own, on OceanBase and on
// a catalog, and compares the error number of every one on both.
func runExecSteps(t *testing.T, label string, cases []struct {
	name  string
	steps []engineStep
}) {
	ctx := context.Background()
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, c, _ := engineDatabase(t, fmt.Sprintf("%s_%02d", label, i))
			for _, step := range tc.steps {
				_, err := db.ExecContext(ctx, step.sql)
				var got uint16
				var myErr *gomysql.MySQLError
				switch {
				case errors.As(err, &myErr):
					got = myErr.Number
				case err != nil:
					t.Fatalf("%s: %v", step.sql, err)
				}
				if got != step.errno {
					t.Errorf("%s\n  OceanBase answered %d (%v), recorded %d", step.sql, got, err, step.errno)
				}
				if code := omniCode(c, step.sql); code != int(step.errno) {
					t.Errorf("%s\n  the catalog answered %d, OceanBase %d", step.sql, code, step.errno)
				}
			}
		})
	}
}

// TestOracle_OceanBaseKeyErrors holds to the engine the errors of a second primary key and
// of a foreign key OceanBase will not add, and the foreign keys it takes that MySQL refuses.
func TestOracle_OceanBaseKeyErrors(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_keys", []struct {
		name  string
		steps []engineStep
	}{
		{"a second primary key", []engineStep{
			{"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))", 1062},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT PRIMARY KEY)", 1068},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT, PRIMARY KEY (b))", 1062},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL)", taken},
			{"ALTER TABLE t ADD PRIMARY KEY (b)", 1068},
			{"ALTER TABLE t ADD COLUMN c INT PRIMARY KEY", 1068},
		}},
		{"a foreign key to nothing", []engineStep{
			{"CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES nosuch (id))", 1146},
			{"CREATE TABLE p (id INT PRIMARY KEY, v INT)", taken},
			{"CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES p (nosuch))", 5031},
			{"CREATE TABLE c (pid INT, FOREIGN KEY (nosuch) REFERENCES p (id))", 5031},
			{"CREATE TABLE c (pid INT)", taken},
			{"ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES nosuch (id)", 1146},
			{"ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES p (nosuch)", 5031},
			{"ALTER TABLE c ADD FOREIGN KEY (nosuch) REFERENCES p (id)", 5031},
		}},
		// The referenced columns lead an index, unique or not, in any order.
		{"a foreign key to columns without a key", []engineStep{
			{"CREATE TABLE p (id INT, v INT, w INT, x INT, y INT, KEY kv (v), UNIQUE KEY uw (w), UNIQUE KEY uxy (x, y))", taken},
			{"CREATE TABLE c1 (pid INT, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (pid INT, FOREIGN KEY (pid) REFERENCES p (v))", taken},
			{"CREATE TABLE c3 (pid INT, FOREIGN KEY (pid) REFERENCES p (w))", taken},
			{"CREATE TABLE c4 (pid INT, FOREIGN KEY (pid) REFERENCES p (x))", taken},
			{"CREATE TABLE c5 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (x, y))", taken},
			{"CREATE TABLE c6 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (y, x))", taken},
			{"CREATE TABLE c7 (pid INT)", taken},
			{"ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (v)", taken},
			{"ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (id)", 1215},
			{"CREATE TABLE c8 (pid INT, FOREIGN KEY (pid) REFERENCES p (y))", 1215},
			{"CREATE TABLE c9 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (v, w))", 1215},
			{"CREATE TABLE q (a INT, b INT, c INT, PRIMARY KEY (a, b), KEY kc (c, a))", taken},
			{"CREATE TABLE d1 (x INT, FOREIGN KEY (x) REFERENCES q (b))", 1215},
			{"CREATE TABLE d2 (x INT, FOREIGN KEY (x) REFERENCES q (a))", taken},
			{"CREATE TABLE d3 (x INT, y INT, FOREIGN KEY (x, y) REFERENCES q (a, c))", taken},
			{"CREATE TABLE d4 (x INT, y INT, z INT, FOREIGN KEY (x, y, z) REFERENCES q (a, b, c))", 1215},
			{"CREATE TABLE s (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES s (id))", taken},
			{"CREATE TABLE s2 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES s2 (id))", 1215},
		}},
		{"a foreign key between types", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (pid BIGINT, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c3 (pid INT UNSIGNED, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c4 (pid INT)", taken},
			{"ALTER TABLE c4 ADD FOREIGN KEY (pid) REFERENCES p (id)", taken},
			{"CREATE TABLE c5 (pid BIGINT)", taken},
			{"ALTER TABLE c5 ADD FOREIGN KEY (pid) REFERENCES p (id)", 1215},
		}},
		{"SET DEFAULT is no action of a foreign key", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT)", 1215},
			{"CREATE TABLE c2 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET DEFAULT)", 1215},
			{"CREATE TABLE c3 (pid INT DEFAULT 0)", taken},
			{"ALTER TABLE c3 ADD FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT", 1215},
		}},
		{"SET NULL needs a nullable column", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET NULL)", 1830},
			{"CREATE TABLE c2 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET NULL)", 1830},
		}},
		{"a foreign key over a generated column", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY, v INT GENERATED ALWAYS AS (id + 1) VIRTUAL, s INT GENERATED ALWAYS AS (id + 1) STORED, UNIQUE KEY uv (v), UNIQUE KEY us (s))", taken},
			{"CREATE TABLE c1 (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL, FOREIGN KEY (b) REFERENCES p (id))", taken},
			{"CREATE TABLE c2 (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED, FOREIGN KEY (b) REFERENCES p (id))", taken},
			{"CREATE TABLE c3 (a INT, FOREIGN KEY (a) REFERENCES p (v))", taken},
			{"CREATE TABLE c4 (a INT, FOREIGN KEY (a) REFERENCES p (s))", taken},
		}},
		{"a foreign key's name twice", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id), CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))", taken},
			{"CREATE TABLE c3 (a INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))", 1022},
			{"ALTER TABLE c2 ADD CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id)", 1022},
		}},
	})
}

// TestOracle_OceanBaseRefusals holds to the engine the statements OceanBase refuses that
// MySQL 8.0 takes: a default on a JSON column, the drop of a primary key's column, a view
// with two columns of one name, the rename of an event that has started, a compressor that
// is not OceanBase's, and REORGANIZE PARTITION on a table without a primary key.
func TestOracle_OceanBaseRefusals(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_refusals", []struct {
		name  string
		steps []engineStep
	}{
		{"a default on a JSON column", []engineStep{
			{"CREATE TABLE t1 (a JSON DEFAULT (JSON_ARRAY()))", 1101},
			{"CREATE TABLE t2 (a JSON DEFAULT '[]')", 1101},
			{"CREATE TABLE t3 (a JSON DEFAULT NULL)", taken},
			{"CREATE TABLE t4 (a BLOB DEFAULT (SUBSTRING('abcdef', 1, 3)))", 1067},
			{"CREATE TABLE t5 (a TEXT DEFAULT ('x'))", 1067},
			{"CREATE TABLE t6 (a TEXT DEFAULT 'x')", 1067},
			{"CREATE TABLE t6 (a TEXT DEFAULT NULL, b BLOB DEFAULT NULL, c POINT DEFAULT NULL)", taken},
			{"CREATE TABLE t7 (a POINT DEFAULT (POINT(0, 0)))", 1101},
			{"CREATE TABLE t8 (a BINARY(16) DEFAULT (UUID_TO_BIN(UUID())))", taken},
			{"CREATE TABLE t9 (id INT)", taken},
			{"ALTER TABLE t9 ADD COLUMN a JSON DEFAULT (JSON_ARRAY())", 1101},
			{"ALTER TABLE t9 ADD COLUMN b JSON", taken},
			{"ALTER TABLE t9 ALTER COLUMN b SET DEFAULT (JSON_ARRAY())", 1101},
			{"ALTER TABLE t9 MODIFY COLUMN b JSON DEFAULT (JSON_OBJECT())", 1101},
		}},
		{"the drop of a primary key's column", []engineStep{
			{"CREATE TABLE t1 (a INT, b INT, c INT, PRIMARY KEY (a, b), UNIQUE KEY uc (c))", taken},
			{"ALTER TABLE t1 DROP COLUMN a", 1235},
			{"ALTER TABLE t1 DROP COLUMN b", 1235},
			{"ALTER TABLE t1 DROP COLUMN c", taken},
			{"CREATE TABLE t2 (a INT PRIMARY KEY, b INT)", taken},
			{"ALTER TABLE t2 DROP COLUMN a", 1235},
			{"ALTER TABLE t2 DROP COLUMN a, DROP PRIMARY KEY", 1235},
		}},
		{"a view with two columns of one name", []engineStep{
			{"CREATE TABLE t1 (a INT, b INT)", taken},
			{"CREATE TABLE t2 (a INT, b INT)", taken},
			{"CREATE VIEW v1 AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a", 1060},
			{"CREATE VIEW v2 AS SELECT a, a FROM t1", 1060},
			{"CREATE VIEW v3 AS SELECT a AS x, b AS X FROM t1", 1060},
			{"CREATE VIEW v4 (x, x) AS SELECT a, b FROM t1", 1060},
			{"CREATE VIEW v5 (x, y) AS SELECT a, a FROM t1", taken},
			{"CREATE VIEW v6 AS SELECT * FROM t1 JOIN t2 ON t1.a = t2.a", 1060},
			{"CREATE VIEW v7 AS SELECT 1, 1", taken},
			{"CREATE VIEW v8 AS SELECT a + 1, a + 1 FROM t1", taken},
			{"CREATE VIEW v9 AS SELECT a, 1 AS a FROM t1", 1060},
			{"CREATE OR REPLACE VIEW v5 AS SELECT a, a FROM t1", 1060},
			{"CREATE VIEW v10 AS SELECT a FROM t1 UNION SELECT b FROM t1", taken},
		}},
		{"a view that compares two collations it cannot", []engineStep{
			{"CREATE VIEW v1 AS SELECT ('a' COLLATE utf8mb4_bin = 'b' COLLATE utf8mb4_general_ci) AS n", 1267},
			{"CREATE VIEW v2 AS SELECT IF('a' COLLATE utf8mb4_bin = 'b' COLLATE utf8mb4_general_ci, 1, 0) AS n", 1267},
			{"CREATE VIEW v3 AS SELECT CASE WHEN 'a' COLLATE utf8mb4_bin = 'b' COLLATE utf8mb4_general_ci THEN 1 ELSE 0 END AS n", 1267},
			{"CREATE VIEW v4 AS SELECT CASE 'a' COLLATE utf8mb4_bin WHEN 'b' COLLATE utf8mb4_general_ci THEN 1 ELSE 0 END AS n", 1267},
			{"CREATE VIEW v5 AS SELECT CASE 'a' COLLATE utf8mb4_bin WHEN 'b' THEN 1 ELSE 0 END AS n", taken},
		}},
		{"the rename of an event that has started", []engineStep{
			{"CREATE EVENT e1 ON SCHEDULE EVERY 1 HOUR DO SELECT 1", taken},
			{"ALTER EVENT e1 RENAME TO e1b", 1235},
			{"CREATE EVENT e2 ON SCHEDULE EVERY 1 HOUR STARTS CURRENT_TIMESTAMP + INTERVAL 1 DAY DO SELECT 1", taken},
			{"ALTER EVENT e2 RENAME TO e2b", taken},
			{"CREATE EVENT e3 ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO SELECT 1", taken},
			{"ALTER EVENT e3 RENAME TO e3b", taken},
			{"CREATE EVENT e4 ON SCHEDULE EVERY 1 HOUR DISABLE DO SELECT 1", taken},
			{"ALTER EVENT e4 RENAME TO e4b", 1235},
			{"ALTER EVENT e1 COMMENT 'x'", taken},
			{"ALTER EVENT e2b ON SCHEDULE EVERY 1 HOUR", taken},
			{"ALTER EVENT e2b RENAME TO e2c", 1235},
		}},
		{"a compressor that is not OceanBase's", []engineStep{
			{"CREATE TABLE t1 (a INT) COMPRESSION='ZLIB'", 1210},
			{"CREATE TABLE t2 (a INT) COMPRESSION='LZ4'", 1210},
			{"CREATE TABLE t3 (a INT) COMPRESSION='NONE'", taken},
			{"CREATE TABLE t4 (a INT) COMPRESSION='none'", taken},
			{"CREATE TABLE t5 (a INT) COMPRESSION='zstd_1.3.8'", taken},
			{"CREATE TABLE t6 (a INT) COMPRESSION='zstd_1.0'", taken},
			{"CREATE TABLE t7 (a INT) COMPRESSION='lz4_1.0'", taken},
			{"CREATE TABLE t8 (a INT) COMPRESSION='snappy_1.0'", taken},
			{"CREATE TABLE t9 (a INT) COMPRESSION='zlib_1.0'", 1235},
			{"CREATE TABLE t10 (a INT) COMPRESSION='lz4_1.9.1'", taken},
			{"CREATE TABLE t11 (a INT) COMPRESSION='ZSTD_1.3.8'", taken},
			{"CREATE TABLE t12 (a INT) COMPRESSION=''", 1210},
			{"ALTER TABLE t3 COMPRESSION='ZLIB'", 1210},
			{"ALTER TABLE t3 COMPRESSION='lz4_1.0'", taken},
		}},
		{"REORGANIZE PARTITION on a table without a primary key", []engineStep{
			{"CREATE TABLE t1 (id INT NOT NULL, val INT NOT NULL) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t1 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", 1235},
			{"CREATE TABLE t2 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (id, val)) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t2 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", 1235},
			{"CREATE TABLE t3 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val, id)) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t3 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", taken},
			{"CREATE TABLE t4 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val)) PARTITION BY RANGE COLUMNS (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t4 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", taken},
			{"CREATE TABLE t5 (a INT NOT NULL, b INT NOT NULL, c INT NOT NULL, PRIMARY KEY (a, b, c)) PARTITION BY RANGE COLUMNS (b, a) (PARTITION p0 VALUES LESS THAN (10, 10), PARTITION p1 VALUES LESS THAN (30, 30))", taken},
			{"ALTER TABLE t5 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20, 20), PARTITION p1b VALUES LESS THAN (30, 30))", 1235},
			{"CREATE TABLE t6 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val, id)) PARTITION BY RANGE (val + 1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t6 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", 1235},
			{"CREATE TABLE t7 (a INT NOT NULL, b INT NOT NULL, c INT NOT NULL, PRIMARY KEY (a, b, c)) PARTITION BY RANGE COLUMNS (a, b) (PARTITION p0 VALUES LESS THAN (10, 10), PARTITION p1 VALUES LESS THAN (30, 30))", taken},
			{"ALTER TABLE t7 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20, 20), PARTITION p1b VALUES LESS THAN (30, 30))", taken},
			// A key written in backticks is a column like any other, whatever its name holds.
			{"CREATE TABLE t8 (`val-id` INT NOT NULL, `a``b c` INT NOT NULL, PRIMARY KEY (`val-id`, `a``b c`)) PARTITION BY RANGE (`val-id`) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t8 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", taken},
			{"CREATE TABLE t9 (`val-id` INT NOT NULL, `a``b c` INT NOT NULL, PRIMARY KEY (`a``b c`, `val-id`)) PARTITION BY RANGE (`a``b c`) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t9 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", taken},
			{"CREATE TABLE t10 (`val-id` INT NOT NULL, id INT NOT NULL, PRIMARY KEY (id, `val-id`)) PARTITION BY RANGE (`val-id`) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))", taken},
			{"ALTER TABLE t10 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))", 1235},
		}},
		{"MySQL's invisible primary key", []engineStep{
			{"SET SESSION sql_generate_invisible_primary_key = 1", 1193},
			{"SET SESSION show_gipk_in_create_table_and_information_schema = 1", 1193},
			{"SET sql_generate_invisible_primary_key = ON", 1193},
			{"SET @@session.sql_generate_invisible_primary_key = 0", 1193},
			{"CREATE TABLE t (my_row_id INT, a INT)", taken},
		}},
	})
}

// TestOracle_OceanBaseAlterTableOptions holds checkAlterTableOptions to the engine:
// which values of OceanBase's table options ALTER TABLE takes, and that the options
// are checked before the statement's other commands, in the order written.
func TestOracle_OceanBaseAlterTableOptions(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT)"
	runExecSteps(t, "oceanbase_catalog_alter_options", []struct {
		name  string
		steps []engineStep
	}{
		{"values", []engineStep{
			{table, taken},
			{"ALTER TABLE a TABLE_MODE 'queuing'", taken},
			{"ALTER TABLE a TABLE_MODE ''", 1235},
			{"ALTER TABLE a DUPLICATE_SCOPE 'CLUSTER'", taken},
			{"ALTER TABLE a DUPLICATE_SCOPE ''", 1210},
			{"ALTER TABLE a DUPLICATE_READ_CONSISTENCY 'weak'", taken},
			{"ALTER TABLE a SEMISTRUCT_ENCODING_TYPE ''", taken},
			{"ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'none'", 1210},
			{"ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'encoding'", 1235},
			{"ALTER TABLE a AUTO_INCREMENT_MODE 'ORDER'", taken},
			{"ALTER TABLE a AUTO_INCREMENT_MODE ''", 1149},
			{"ALTER TABLE a KV_ATTRIBUTES '{\"Hbase\": {}}'", taken},
			{"ALTER TABLE a BLOCK_SIZE 1024", taken},
			{"ALTER TABLE a BLOCK_SIZE 1048576", taken},
			{"ALTER TABLE a BLOCK_SIZE 1048577", 5258},
			{"ALTER TABLE a STORAGE_FORMAT_VERSION 4", taken},
			{"ALTER TABLE a STORAGE_FORMAT_VERSION 3", 4179},
			{"ALTER TABLE a STORAGE_FORMAT_VERSION 5", 1210},
			{"ALTER TABLE a TABLET_SIZE 0", taken},
			{"ALTER TABLE a TABLET_SIZE 4194304", taken},
			{"ALTER TABLE a TABLET_SIZE 3145728", 4147},
			{"ALTER TABLE a ORGANIZATION HEAP", 1235},
			{"ALTER TABLE a FORMAT(TYPE = 'csv')", 4016},
		}},
		{"a ROW_FORMAT = COMPRESSED table takes semistructured encoding", []engineStep{
			{table + " ROW_FORMAT = COMPRESSED", taken},
			{"ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'encoding'", taken},
		}},
		{"a dynamic partition policy", []engineStep{
			{table, taken},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'minute')", 1210},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = '1 hours')", 1210},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = '+13:01')", 1298},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's', TIME_UNIT = 's')", 1298},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's', TIME_ZONE = 's')", 1210},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour' TIME_UNIT = 'hour')", 1210},
			{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = '-12:59', PRECREATE_TIME = '+1day')", 1235},
			{"CREATE TABLE h (b INT PRIMARY KEY) PARTITION BY HASH(b) PARTITIONS 2", taken},
			{"ALTER TABLE h DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour')", 1235},
			{"CREATE TABLE r (b DATETIME PRIMARY KEY) PARTITION BY RANGE COLUMNS(b) (PARTITION p VALUES LESS THAN ('2026-01-01'))", taken},
			{"ALTER TABLE r DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour')", taken},
		}},
		{"options first, in order", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD b INT, BLOCK_SIZE 1", 5258},
			{"ALTER TABLE a DROP COLUMN zz, DUPLICATE_SCOPE 's'", 1210},
			{"ALTER TABLE a ADD KEY(zz), TABLE_MODE 's'", 1235},
			{"ALTER TABLE a TABLE_MODE 's', LOCATION 's'", 1235},
			{"ALTER TABLE a LOCATION 's', TABLE_MODE 's'", 4016},
		}},
	})
}
