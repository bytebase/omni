package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	tidbcatalog "github.com/bytebase/omni/tidb/catalog"
)

// recallSchema is a synced schema with the facts the partitioning,
// character set, AUTO_RANDOM, and sequence-bounds checks read.
func recallSchema() *metadata.DatabaseSchemaMetadata {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	pk := []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}}
	part := func(typ metadata.TablePartitionMetadata_Type, names ...string) []*metadata.TablePartitionMetadata {
		var out []*metadata.TablePartitionMetadata
		for _, n := range names {
			out = append(out, &metadata.TablePartitionMetadata{Name: n, Type: typ, Expression: "`a`"})
		}
		return out
	}
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{Name: "t", Columns: cols("id", "a", "d"), Indexes: pk, Charset: "utf8mb4", Collation: "utf8mb4_bin"},
				{Name: "latin", Columns: cols("id"), Indexes: pk, Collation: "latin1_swedish_ci"},
				{Name: "pt", Columns: cols("id", "a"), Indexes: pk, Partitions: part(metadata.TablePartitionMetadata_RANGE, "p0", "p1")},
				{Name: "lt", Columns: cols("id", "a"), Indexes: pk, Partitions: part(metadata.TablePartitionMetadata_LIST, "p0")},
				{Name: "ht", Columns: cols("id", "a"), Indexes: pk, Partitions: part(metadata.TablePartitionMetadata_HASH, "p0", "p1")},
				{Name: "ar", Columns: cols("id", "a"), Indexes: pk, ShardingInfo: "PK_AUTO_RANDOM_BITS=5"},
			},
			Views:     []*metadata.ViewMetadata{{Name: "v"}},
			Sequences: []*metadata.SequenceMetadata{{Name: "sq", Start: "1", MinValue: "1", MaxValue: "100", Increment: "1"}},
		}},
	}
}

func TestWalkThroughTiDBPartitioning(t *testing.T) {
	shop := review.Target{Schema: recallSchema()}
	for _, tc := range []struct{ sql, want string }{
		{"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (10));", "Table `t` is not partitioned."},
		{"ALTER TABLE t ADD PARTITION IF NOT EXISTS (PARTITION p2 VALUES LESS THAN (10));", "Table `t` is not partitioned."},
		{"ALTER TABLE t DROP PARTITION IF EXISTS p0;", "Table `t` is not partitioned."},
		{"ALTER TABLE t TRUNCATE PARTITION ALL;", "Table `t` is not partitioned."},
		{"ALTER TABLE t ANALYZE PARTITION p0;", "Table `t` is not partitioned."},
		{"ALTER TABLE t COALESCE PARTITION 1;", "Table `t` is not partitioned."},
		{"ALTER TABLE t REMOVE PARTITIONING;", "Table `t` is not partitioned."},
		{"ALTER TABLE t REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10));", "Table `t` is not partitioned."},
		{"ALTER TABLE t PARTITION p0 ATTRIBUTES = 'k=v';", "Table `t` is not partitioned."},
		{"ALTER TABLE t COMPACT PARTITION p0;", "Table `t` is not partitioned."},
		{"ALTER TABLE v COMPACT PARTITION p0;", "`v` is a view, not a table."},
		{"ALTER TABLE pt CACHE;", "TiDB can't cache `pt`, which is partitioned."},
		{"ALTER TABLE pt COALESCE PARTITION 1;", "`COALESCE PARTITION` works only on `HASH` or `KEY` partitioning, and `pt` has `RANGE`."},
		{"ALTER TABLE pt ADD PARTITION PARTITIONS 2;", "Each partition of `pt` has to be defined for `RANGE` partitioning."},
		{"ALTER TABLE lt ADD PARTITION;", "Each partition of `lt` has to be defined for `LIST` partitioning."},
		{"ALTER TABLE pt REORGANIZE PARTITION;", "`REORGANIZE PARTITION` of `pt` has to name the partitions it reorganizes."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p2);", "Partition `p2` needs `VALUES LESS THAN` for `RANGE` partitioning."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p2 VALUES IN (1));", "Partition `p2` can't use `VALUES IN` with `RANGE` partitioning."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p2 DEFAULT);", "Partition `p2` can't use `VALUES IN` with `RANGE` partitioning."},
		{"ALTER TABLE lt ADD PARTITION (PARTITION p2 VALUES LESS THAN (10));", "Partition `p2` can't use `VALUES LESS THAN` with `LIST` partitioning."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p2 VALUES LESS THAN (10, 20));", "Partition `p2` can't have more than one `VALUES LESS THAN` value for `RANGE` partitioning."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION P1 VALUES LESS THAN (100));", "Partition `P1` already exists in `pt`."},
		// The walk follows what a statement does to the partitioning.
		{"ALTER TABLE pt REMOVE PARTITIONING;\nALTER TABLE pt DROP PARTITION p0;", "Table `pt` is not partitioned."},
		{"ALTER TABLE t PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10));\nALTER TABLE t COALESCE PARTITION 1;", "`COALESCE PARTITION` works only on `HASH` or `KEY` partitioning, and `t` has `RANGE`."},
		{"CREATE TABLE n (a int);\nALTER TABLE n DROP PARTITION p0;", "Table `n` is not partitioned."},
		{"ALTER TABLE pt ADD PARTITION (PARTITION p2 VALUES LESS THAN (300));\nALTER TABLE pt ADD PARTITION (PARTITION p2 VALUES LESS THAN (400));", "Partition `p2` already exists in `pt`."},
		{"RENAME TABLE t TO t2;\nALTER TABLE t2 TRUNCATE PARTITION ALL;", "Table `t2` is not partitioned."},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			got := walkFindings(t, tc.sql, shop)
			if len(got) != 1 || got[0].Message != tc.want {
				t.Errorf("%s\nfindings = %+v, want %q", tc.sql, got, tc.want)
			}
		})
	}
	for _, sql := range []string{
		"ALTER TABLE t PARTITION p0;",
		"ALTER TABLE t COMPACT;",
		"ALTER TABLE t CACHE;",
		"ALTER TABLE pt ADD PARTITION (PARTITION p2 VALUES LESS THAN (300));",
		"ALTER TABLE pt ADD PARTITION IF NOT EXISTS (PARTITION p1 VALUES LESS THAN (300));",
		"ALTER TABLE pt DROP PARTITION p0;\nALTER TABLE pt ADD PARTITION (PARTITION p0 VALUES LESS THAN (300));",
		"ALTER TABLE pt REORGANIZE PARTITION p1 INTO (PARTITION p1 VALUES LESS THAN (100), PARTITION p2 VALUES LESS THAN (200));",
		"ALTER TABLE pt TRUNCATE PARTITION p0;",
		"ALTER TABLE pt COMPACT PARTITION p0;",
		"ALTER TABLE ht COALESCE PARTITION 1;",
		"ALTER TABLE ht ADD PARTITION PARTITIONS 2;",
		"ALTER TABLE lt ADD PARTITION (PARTITION p1 VALUES IN (5));",
		"ALTER TABLE lt ADD PARTITION (PARTITION p1 DEFAULT);",
		"ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 4;\nALTER TABLE t COALESCE PARTITION 2;",
		"ALTER TABLE pt REMOVE PARTITIONING;\nALTER TABLE pt CACHE;",
		// A table the walk cannot read the partitioning of claims nothing.
		"CREATE TABLE c LIKE pt;\nALTER TABLE c DROP PARTITION p0;",
		"CREATE TABLE s (a int) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (1));\nALTER TABLE s DROP PARTITION p0;",
		"DROP TABLE t;\nCREATE TABLE t (a int) PARTITION BY HASH (a);\nALTER TABLE t TRUNCATE PARTITION ALL;",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

func TestWalkThroughTiDBTableFacts(t *testing.T) {
	shop := review.Target{Schema: recallSchema()}
	for _, tc := range []struct{ sql, want string }{
		{"ALTER TABLE t CHARSET binary;", "TiDB can't change the character set of `t` from `utf8mb4` to `binary`."},
		{"ALTER TABLE t COLLATE binary;", "TiDB can't change the character set of `t` from `utf8mb4` to `binary`."},
		{"ALTER TABLE t CONVERT TO CHARSET latin1;", "TiDB can't change the character set of `t` from `utf8mb4` to `latin1`."},
		{"ALTER TABLE t CHARSET utf8mb3;", "TiDB can't change the character set of `t` from `utf8mb4` to `utf8`."},
		{"ALTER TABLE latin CHARSET binary;", "TiDB can't change the character set of `latin` from `latin1` to `binary`."},
		{"ALTER TABLE latin CHARSET utf8;", "TiDB can't change the character set of `latin` from `latin1` to `utf8`."},
		{"CREATE TABLE n (a int) CHARSET latin1;\nALTER TABLE n CHARSET binary;", "TiDB can't change the character set of `n` from `latin1` to `binary`."},
		{"ALTER TABLE t AUTO_RANDOM_BASE = 100;", "Table `t` has no `AUTO_RANDOM` column to rebase."},
		{"ALTER TABLE t TTL = nosuch + INTERVAL 1 DAY;", "Column `nosuch` doesn't exist in `t`."},
		{"ALTER TABLE t ALTER a SET DEFAULT (d);", "Column `d` is unknown in a `DEFAULT` expression, which reads no columns."},
		{"DROP INDEX i ON v;", "Index `i` doesn't exist in `v`."},
		{"ALTER SEQUENCE sq MAXVALUE 1;", "Sequence `sq` has conflicting options."},
		{"ALTER SEQUENCE sq START WITH 200;", "Sequence `sq` has conflicting options."},
		{"ALTER SEQUENCE sq MINVALUE 50;", "Sequence `sq` has conflicting options."},
		{"ALTER SEQUENCE sq MAXVALUE 50 START WITH 10;\nALTER SEQUENCE sq START WITH 60;", "Sequence `sq` has conflicting options."},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			got := walkFindings(t, tc.sql, shop)
			if len(got) != 1 || got[0].Message != tc.want {
				t.Errorf("%s\nfindings = %+v, want %q", tc.sql, got, tc.want)
			}
		})
	}
	for _, sql := range []string{
		"ALTER TABLE t CHARSET utf8mb4 COLLATE utf8mb4_general_ci;",
		"ALTER TABLE latin CONVERT TO CHARSET utf8mb4;",
		"ALTER TABLE t CHARSET utf8mb4;\nALTER TABLE t CHARSET utf8mb4;",
		// A table created without a character set takes the database's,
		// which the walk does not read.
		"CREATE TABLE n (a int);\nALTER TABLE n CHARSET binary;",
		"ALTER TABLE ar AUTO_RANDOM_BASE = 100;",
		"CREATE TABLE n (id bigint AUTO_RANDOM PRIMARY KEY);\nALTER TABLE n AUTO_RANDOM_BASE = 100;",
		"ALTER TABLE t MODIFY id bigint AUTO_RANDOM;\nALTER TABLE t AUTO_RANDOM_BASE = 100;",
		"ALTER TABLE t ALTER a SET DEFAULT (1 + 1);",
		"ALTER TABLE t ALTER a SET DEFAULT 5;",
		"DROP INDEX IF EXISTS i ON v;",
		"ALTER SEQUENCE sq MAXVALUE 50;",
		"ALTER SEQUENCE sq NOMAXVALUE;\nALTER SEQUENCE sq MAXVALUE 1000;",
		"ALTER SEQUENCE sq RESTART WITH 5;",
		"DROP SEQUENCE sq;\nCREATE SEQUENCE sq START WITH 500 MAXVALUE 1000;\nALTER SEQUENCE sq MINVALUE 100;",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

func TestWalkThroughTiDBViewBodies(t *testing.T) {
	shop := review.Target{Schema: recallSchema()}
	for _, tc := range []struct{ sql, want string }{
		{"CREATE VIEW w AS SELECT @x FROM t;", "A view's query can't read a variable."},
		{"CREATE VIEW w AS SELECT @@sql_mode;", "A view's query can't read a variable."},
		{"CREATE VIEW w AS SELECT a FROM t WHERE a = ?;", "A view's query can't hold a parameter marker `?`."},
		{"CREATE VIEW w AS SELECT YEAR() FROM t;", "`YEAR()` needs at least one argument."},
		{"CREATE VIEW w AS SELECT NTILE(a) OVER () FROM t;", "`NTILE` takes a constant, not a column, as its count."},
		{"CREATE VIEW w AS SELECT RANK() OVER w1 FROM t;", "Window `w1` isn't defined in the query's `WINDOW` clause."},
		{"CREATE VIEW w AS SELECT RANK() OVER w1 FROM t WINDOW w2 AS (w3);", "Window `w3` isn't defined in the query's `WINDOW` clause."},
		{"CREATE VIEW w AS SELECT a FROM t WINDOW w1 AS (), w1 AS ();", "Window `w1` is defined twice."},
		{"CREATE VIEW w AS SELECT (SELECT RANK() OVER w1 FROM t) FROM t WINDOW w1 AS ();", "Window `w1` isn't defined in the query's `WINDOW` clause."},
		{"CREATE VIEW w AS SELECT SUM(a) OVER (GROUPS UNBOUNDED PRECEDING) FROM t;", "TiDB doesn't support a `GROUPS` window frame."},
		{"CREATE VIEW w AS SELECT SUM(a) OVER (ROWS BETWEEN UNBOUNDED FOLLOWING AND CURRENT ROW) FROM t;", "A window frame can't start at `UNBOUNDED FOLLOWING`."},
		{"CREATE VIEW w AS SELECT SUM(a) OVER (ROWS BETWEEN CURRENT ROW AND 1 PRECEDING) FROM t;", "A window frame can't end before it starts."},
		{"CREATE VIEW w AS SELECT SUM(a) OVER (ORDER BY a ROWS INTERVAL 1 DAY PRECEDING) FROM t;", "A `ROWS` window frame can't take an `INTERVAL` bound."},
		{"CREATE VIEW w AS SELECT a, a FROM t;", "This statement names view column `a` twice."},
		{"CREATE VIEW w AS SELECT a, d AS A FROM t;", "This statement names view column `A` twice."},
		{"CREATE VIEW w AS SELECT NEXTVAL(t);", "`t` is a table, not a sequence."},
		{"CREATE VIEW w AS SELECT a;", "Column `a` doesn't exist in a query that reads no table."},
		{"CREATE VIEW w AS SELECT nosuch FROM t;", "Column `nosuch` doesn't exist in `t`."},
		{"CREATE VIEW w AS SELECT x.nosuch FROM t AS x WHERE x.a = 1;", "Column `nosuch` doesn't exist in `t`."},
		{"CREATE VIEW w AS SELECT a FROM t UNION SELECT nosuch FROM t;", "Column `nosuch` doesn't exist in `t`."},
		{"CREATE VIEW w AS SELECT MATCH(a) AGAINST ('x') FROM t;", "TiDB doesn't support `MATCH ... AGAINST`."},
		{"CREATE VIEW w AS SELECT CAST(a AS UNSIGNED ARRAY) FROM t;", "TiDB doesn't support `CAST(... AS ... ARRAY)` outside an index."},
		{"CREATE VIEW w AS SELECT DATE 'none';", "`DATE 'none'` isn't a valid date literal."},
		{"CREATE VIEW w AS SELECT shop.f(a) FROM t;", "Function `shop.f` doesn't exist; TiDB has no stored functions."},
		{"CREATE VIEW w AS SELECT f(a) FROM t;", "Function `f` doesn't exist; TiDB has no stored functions."},
		{"CREATE VIEW w AS SELECT ROW(a, d) FROM t;", "A select item can't be a row of 2 columns."},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			got := walkFindings(t, tc.sql, shop)
			if len(got) != 1 || got[0].Message != tc.want {
				t.Errorf("%s\nfindings = %+v, want %q", tc.sql, got, tc.want)
			}
		})
	}
	for _, sql := range []string{
		"CREATE VIEW w AS SELECT a, d FROM t;",
		"CREATE VIEW w (x, y) AS SELECT a, a FROM t;",
		"CREATE VIEW w AS SELECT a, a + 1 FROM t;",
		"CREATE VIEW w AS SELECT RANK() OVER w1 FROM t WINDOW w1 AS (ORDER BY a);",
		"CREATE VIEW w AS SELECT RANK() OVER (W1 ORDER BY a) FROM t WINDOW w1 AS (PARTITION BY d);",
		"CREATE VIEW w AS SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t;",
		"CREATE VIEW w AS SELECT SUM(a) OVER (ORDER BY a RANGE BETWEEN INTERVAL 1 DAY PRECEDING AND CURRENT ROW) FROM t;",
		"CREATE VIEW w AS SELECT NTH_VALUE(a, 2) OVER () FROM t;",
		"CREATE VIEW w AS SELECT NEXTVAL(sq);",
		"CREATE VIEW w AS SELECT NEXT VALUE FOR sq FROM t;",
		"CREATE VIEW w AS SELECT YEAR(NOW()), CURRENT_TIMESTAMP, CURRENT_DATE, PI(), NOW();",
		"CREATE VIEW w AS SELECT a AS x FROM t ORDER BY x;",
		"CREATE VIEW w AS SELECT a FROM t GROUP BY a HAVING a > 1;",
		"CREATE VIEW w AS SELECT t.a, x.d FROM t, t AS x;",
		"CREATE VIEW w AS SELECT nosuch FROM v;",
		"CREATE VIEW w AS SELECT (SELECT a) FROM t;",
		"CREATE VIEW w AS SELECT a FROM t WHERE EXISTS (SELECT 1 FROM v WHERE v.x = t.a);",
		"CREATE VIEW w AS WITH c AS (SELECT 1 AS z) SELECT z FROM c;",
		"CREATE VIEW w AS SELECT EXTRACT(DAY FROM d), TIMESTAMPADD(MINUTE, 1, d), DATE_ADD(d, INTERVAL 1 DAY) FROM t;",
		"CREATE VIEW w AS SELECT DATE '2020-01-01', TIME '10:00:00';",
		"CREATE VIEW w AS SELECT ROW(1, 2) = ROW(a, d) FROM t;",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

// shapeDDL builds the synced schema the shape, foreign key, and view type
// checks read (walk_table_shape.go, walk_fk.go, walk_columns.go), through
// tidb/catalog as the simulation does.
const shapeDDL = `CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(20), ts TIMESTAMP NULL, j JSON, dt DATE, KEY idx_a (a));
CREATE TABLE parent (id INT PRIMARY KEY, email VARCHAR(100));
CREATE TABLE child (id INT PRIMARY KEY, pid INT, KEY fk_pid (pid), CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent (id));
CREATE TABLE hsh (id INT PRIMARY KEY, u INT) PARTITION BY HASH (id) PARTITIONS 3;
CREATE TABLE nc (id INT, a INT, PRIMARY KEY (id) NONCLUSTERED);
CREATE TABLE comp (id INT, a INT, b INT, PRIMARY KEY (id), KEY ab (a, b));
CREATE TABLE ac (id INT PRIMARY KEY) AUTO_ID_CACHE 1;
CREATE VIEW v AS SELECT id FROM t;`

// shapeCase is a change on shapeDDL and its one finding, at statement, or
// none when want is empty. Each refusal was verified on v8.5.5; an empty
// want is a change the walk lets through, which may still fail for a
// reason it doesn't follow (noted beside it).
type shapeCase struct {
	sql       string
	statement int
	want      string
}

func runShapeCases(t *testing.T, cases []shapeCase) {
	t.Helper()
	c, err := tidbcatalog.LoadSQL("CREATE DATABASE shop;\nUSE shop;\n" + shapeDDL)
	if err != nil {
		t.Fatal(err)
	}
	target := review.Target{Schema: c.DatabaseMetadata("shop"), ForeignKeyChecks: review.SwitchOn}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			if tc.want == "" {
				noFinding(t, tc.sql, target)
				return
			}
			oneFinding(t, tc.sql, tc.statement, tc.want, target)
		})
	}
}

// TestWalkThroughTiDBForeignKeyDependents pins the foreign key checks (walk_fk.go).
func TestWalkThroughTiDBForeignKeyDependents(t *testing.T) {
	runShapeCases(t, []shapeCase{
		{sql: "DROP TABLE parent;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "DROP TABLE IF EXISTS parent, nosuch;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "TRUNCATE TABLE parent;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "DROP TABLE parent, child;"},
		{sql: "TRUNCATE TABLE child;"},
		{sql: "SET foreign_key_checks = 0;\nDROP TABLE parent;"},
		{sql: "SET foreign_key_checks = 0;\nTRUNCATE TABLE parent;"},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid;\nDROP TABLE parent;"},
		// TiDB follows the rename and refuses (3730); the walk loses the
		// foreign key.
		{sql: "ALTER TABLE child RENAME TO kid;\nDROP TABLE parent;"},
		{sql: "ALTER TABLE child DROP COLUMN pid;", want: "Column `pid` is in foreign key `fk_pid` of `child`, and can't be dropped."},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid, DROP COLUMN pid;", want: "Column `pid` is in foreign key `fk_pid` of `child`, and can't be dropped."},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid;\nALTER TABLE child DROP COLUMN pid;"},
		{sql: "ALTER TABLE parent DROP COLUMN id;", want: "Column `id` is referenced by foreign key `fk_pid` of `child`, and can't be dropped."},
		{sql: "ALTER TABLE child DROP INDEX fk_pid;", want: "Index `fk_pid` of `child` serves foreign key `fk_pid`, and no other index does."},
		{sql: "DROP INDEX fk_pid ON child;", want: "Index `fk_pid` of `child` serves foreign key `fk_pid`, and no other index does."},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid, DROP INDEX fk_pid;", want: "Index `fk_pid` of `child` serves foreign key `fk_pid`, and no other index does."},
		{sql: "ALTER TABLE child ADD INDEX i2 (pid, id);\nALTER TABLE child DROP INDEX fk_pid;"},
		// TiDB refuses (1553); the walk no longer knows the indexes.
		{sql: "ALTER TABLE child RENAME INDEX fk_pid TO zz;\nALTER TABLE child DROP INDEX zz;"},
		{sql: "ALTER TABLE child MODIFY pid BIGINT;", want: "Foreign key `fk_pid` would pair `pid` (`bigint`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "ALTER TABLE child CHANGE pid pid VARCHAR(10);", want: "Foreign key `fk_pid` would pair `pid` (`varchar`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "ALTER TABLE parent MODIFY id BIGINT;", want: "Foreign key `fk_pid` would pair `id` (`bigint`) with `child`.`pid` (`int`), types that don't match."},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid, MODIFY pid BIGINT;", want: "Foreign key `fk_pid` would pair `pid` (`bigint`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "ALTER TABLE child MODIFY pid INT UNSIGNED;"},
		{sql: "ALTER TABLE child MODIFY pid INT(3);"},
		// TiDB follows the rename and refuses (3780); the walk claims
		// nothing of a renamed column's foreign keys.
		{sql: "ALTER TABLE child RENAME COLUMN pid TO p2;\nALTER TABLE child MODIFY p2 BIGINT;"},
		{sql: "ALTER TABLE child ADD CONSTRAINT f2 FOREIGN KEY (id) REFERENCES parent (email);", want: "A foreign key would pair `id` (`int`) with `parent`.`email` (`varchar`), types that don't match."},
		{sql: "CREATE TABLE x (id BIGINT, FOREIGN KEY (id) REFERENCES parent (id));", want: "A foreign key would pair `id` (`bigint`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "SET foreign_key_checks = 0;\nCREATE TABLE x (id BIGINT, FOREIGN KEY (id) REFERENCES parent (id));", statement: 1, want: "A foreign key would pair `id` (`bigint`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "CREATE TABLE x (id INT(3), FOREIGN KEY (id) REFERENCES parent (id));"},
		// TiDB refuses a sign apart here (3780), which the walk doesn't
		// claim: MODIFY takes it.
		{sql: "CREATE TABLE x (id INT UNSIGNED, FOREIGN KEY (id) REFERENCES parent (id));"},
	})
}

// TestWalkThroughTiDBColumnsAndKeys pins the checks of a table's keys,
// column types, and options (walk_columns.go).
func TestWalkThroughTiDBColumnsAndKeys(t *testing.T) {
	runShapeCases(t, []shapeCase{
		{sql: "ALTER TABLE comp DROP COLUMN a;", want: "TiDB can't drop column `a`, which index `ab` holds with other columns."},
		{sql: "ALTER TABLE comp DROP INDEX ab, DROP COLUMN a;", want: "TiDB can't drop column `a`, which index `ab` holds with other columns."},
		{sql: "ALTER TABLE comp DROP INDEX ab;\nALTER TABLE comp DROP COLUMN a;"},
		{sql: "ALTER TABLE nc DROP COLUMN id;", want: "TiDB can't drop column `id`, which is in the primary key of `nc`."},
		{sql: "ALTER TABLE t DROP COLUMN id;", want: "TiDB can't drop column `id`, which is in the primary key of `t`."},
		{sql: "ALTER TABLE t DROP COLUMN a;"},
		{sql: "ALTER TABLE t MODIFY id TINYINT;", want: "TiDB can't change primary key column `id` from `int` to `tinyint`, which rewrites its data."},
		{sql: "ALTER TABLE nc MODIFY id INT UNSIGNED;", want: "TiDB can't change primary key column `id` from `int` to `int unsigned`, which rewrites its data."},
		{sql: "ALTER TABLE t MODIFY id BIGINT;"},
		{sql: "ALTER TABLE hsh MODIFY u VARCHAR(10);", want: "TiDB can't change column `u` from `int` to `varchar`, which rewrites its data, in partitioned table `hsh`."},
		{sql: "ALTER TABLE hsh MODIFY u BIGINT;"},
		{sql: "ALTER TABLE comp MODIFY a TINYINT;"},
		{sql: "ALTER TABLE t TTL = a + INTERVAL 1 DAY;", want: "TTL column `a` is `int`; TiDB takes only a `DATETIME`, `DATE`, or `TIMESTAMP` column."},
		{sql: "ALTER TABLE t TTL = ts + INTERVAL 1 DAY;"},
		{sql: "CREATE TABLE x (a INT) TTL = a + INTERVAL 1 DAY;", want: "TTL column `a` is `int`; TiDB takes only a `DATETIME`, `DATE`, or `TIMESTAMP` column."},
		{sql: "ALTER TABLE t AUTO_ID_CACHE 1;", want: "TiDB can't move `AUTO_ID_CACHE` of `t` between 1 and another value."},
		{sql: "ALTER TABLE t AUTO_ID_CACHE 100;"},
		{sql: "ALTER TABLE ac AUTO_ID_CACHE 100;", want: "TiDB can't move `AUTO_ID_CACHE` of `ac` between 1 and another value."},
		{sql: "ALTER TABLE ac AUTO_ID_CACHE 1;"},
		{sql: "CREATE TABLE x (id INT PRIMARY KEY) AUTO_ID_CACHE 1;\nALTER TABLE x AUTO_ID_CACHE 100;", statement: 1, want: "TiDB can't move `AUTO_ID_CACHE` of `x` between 1 and another value."},
		{sql: "ALTER TABLE t DROP PRIMARY KEY;", want: "TiDB can't drop the primary key of `t`, which clusters the table."},
		{sql: "ALTER TABLE comp DROP INDEX `PRIMARY`;", want: "TiDB can't drop the primary key of `comp`, which clusters the table."},
		{sql: "DROP INDEX `PRIMARY` ON t;", want: "TiDB can't drop the primary key of `t`, which clusters the table."},
		{sql: "ALTER TABLE t RENAME TO t2;\nALTER TABLE t2 DROP PRIMARY KEY;", statement: 1, want: "TiDB can't drop the primary key of `t2`, which clusters the table."},
		{sql: "ALTER TABLE nc DROP PRIMARY KEY;"},
		// TiDB refuses this too (8200); the walk leaves IF EXISTS alone.
		{sql: "ALTER TABLE t DROP INDEX IF EXISTS `PRIMARY`;"},
		// Whether TiDB clusters the key is tidb_enable_clustered_index,
		// which the walk does not read from CREATE TABLE.
		{sql: "CREATE TABLE x (id INT PRIMARY KEY);\nALTER TABLE x DROP PRIMARY KEY;"},
		{sql: "DROP TABLE t;\nCREATE TABLE t (id INT, a INT);\nALTER TABLE t SHARD_ROW_ID_BITS = 4;"},
		{sql: "ALTER TABLE t SHARD_ROW_ID_BITS = 4;", want: "TiDB can't set `SHARD_ROW_ID_BITS` on `t`, which is clustered by its primary key and has no row id to shard."},
		{sql: "ALTER TABLE comp SHARD_ROW_ID_BITS 1, COMMENT 'c';", want: "TiDB can't set `SHARD_ROW_ID_BITS` on `comp`, which is clustered by its primary key and has no row id to shard."},
		{sql: "ALTER TABLE t SHARD_ROW_ID_BITS = 0;"},
		{sql: "ALTER TABLE nc SHARD_ROW_ID_BITS = 4;"},
		// TiDB refuses this too (8200); the walk claims only an option that
		// comes first.
		{sql: "ALTER TABLE t COMMENT 'c', SHARD_ROW_ID_BITS = 4;"},
		{sql: "ALTER TABLE t PARTITION BY HASH (a);", want: "Every unique key of `t` must hold all columns of its partitioning, and the primary key doesn't."},
		{sql: "ALTER TABLE nc PARTITION BY HASH (a);", want: "Every unique key of `nc` must hold all columns of its partitioning, and the primary key doesn't."},
		{sql: "ALTER TABLE t PARTITION BY HASH (id);"},
		// TiDB refuses a clustered primary key global (1503); the walk
		// leaves UPDATE INDEXES alone.
		{sql: "ALTER TABLE t PARTITION BY HASH (a) UPDATE INDEXES (`PRIMARY` GLOBAL);"},
		{sql: "CREATE TABLE x (a INT PRIMARY KEY, b INT) PARTITION BY HASH (b);", want: "Every unique key of `x` must hold all columns of its partitioning, and the primary key doesn't."},
		{sql: "CREATE TABLE x (a INT, b INT, UNIQUE KEY uk (a)) PARTITION BY HASH (b);", want: "Every unique key of `x` must hold all columns of its partitioning, and unique key `uk` doesn't."},
		{sql: "CREATE TABLE x (a INT PRIMARY KEY, b DATE) PARTITION BY RANGE (YEAR(b)) (PARTITION p0 VALUES LESS THAN (2000));", want: "Every unique key of `x` must hold all columns of its partitioning, and the primary key doesn't."},
		{sql: "CREATE TABLE x (a INT, b INT, UNIQUE KEY (a) GLOBAL) PARTITION BY HASH (b);"},
		{sql: "CREATE TABLE x (a INT, b INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS (b) (PARTITION p0 VALUES LESS THAN (10));"},
		{sql: "ALTER TABLE t ADD INDEX i3 (a), ALGORITHM = INSTANT;", want: "TiDB can't add an index or a constraint with `ALGORITHM=INSTANT`; it takes `INPLACE`."},
		{sql: "ALTER TABLE t ADD INDEX i3 (a), ALGORITHM = INPLACE;"},
		{sql: "ALTER TABLE t ADD COLUMN z INT, ALGORITHM = INSTANT;"},
		{sql: "CREATE INDEX i3 ON t (a) ALGORITHM = INSTANT;"},
		{sql: "ALTER TABLE t CONVERT TO CHARACTER SET latin1 COLLATE utf8mb4_bin;", want: "Collation `utf8mb4_bin` isn't one of character set `latin1`."},
		{sql: "ALTER TABLE t CHARSET utf8mb4 COLLATE utf8_bin;", want: "Collation `utf8_bin` isn't one of character set `utf8mb4`."},
		{sql: "CREATE TABLE x (a VARCHAR(10) CHARACTER SET latin1 COLLATE utf8mb4_bin);", want: "Collation `utf8mb4_bin` isn't one of character set `latin1`."},
		{sql: "CREATE TABLE x (a INT) DEFAULT CHARSET = latin1 DEFAULT COLLATE = utf8mb4_bin;", want: "Collation `utf8mb4_bin` isn't one of character set `latin1`."},
		{sql: "CREATE TABLE x (a VARCHAR(10) CHARACTER SET utf8 COLLATE utf8mb3_bin) CHARSET utf8mb4 COLLATE utf8mb4_bin;"},
		{sql: "CREATE TABLE x (a INT, b INT AS (RAND()));", want: "Generated column `b` can't call `RAND`, whose value isn't deterministic."},
		{sql: "CREATE TABLE x (a INT, b INT AS (CURRENT_TIMESTAMP));", want: "Generated column `b` can't call `CURRENT_TIMESTAMP`, whose value isn't deterministic."},
		{sql: "ALTER TABLE t ADD COLUMN z INT AS (UUID());", want: "Generated column `z` can't call `UUID`, whose value isn't deterministic."},
		{sql: "ALTER TABLE t ADD INDEX fi ((UUID()));", want: "An expression index can't call `UUID`, whose value isn't deterministic."},
		{sql: "CREATE INDEX fi ON t ((RAND()));", want: "An expression index can't call `RAND`, whose value isn't deterministic."},
		{sql: "CREATE TABLE x (a INT, INDEX ((NOW())));", want: "An expression index can't call `NOW`, whose value isn't deterministic."},
		{sql: "CREATE TABLE x (a INT, b INT AS (a + 1), INDEX ((a * 2)));"},
		{sql: "CREATE TABLE x (`current_timestamp` INT, b INT AS (`current_timestamp`), INDEX ((`current_timestamp` + 1)));"},
	})
}

// TestWalkThroughTiDBViewTypes pins the view checks that read functions
// and column types (walk_view.go).
func TestWalkThroughTiDBViewTypes(t *testing.T) {
	runShapeCases(t, []shapeCase{
		{sql: "CREATE VIEW w AS SELECT nope(a) FROM t;", want: "Function `nope` doesn't exist; TiDB has no stored functions."},
		{sql: "CREATE VIEW w AS SELECT POINT(1, 2);", want: "Function `point` doesn't exist; TiDB has no stored functions."},
		{sql: "CREATE VIEW w AS WITH c AS (SELECT nope()) SELECT 1;", want: "Function `nope` doesn't exist; TiDB has no stored functions."},
		{sql: "CREATE VIEW w AS SELECT APPROX_PERCENTILE(a) FROM t;", want: "`APPROX_PERCENTILE` takes 2 arguments."},
		{sql: "CREATE VIEW w AS SELECT APPROX_PERCENTILE(a, 50), ABS(a), NOW(), JSON_LENGTH(a) FROM t;"},
		{sql: "CREATE VIEW w AS SELECT a -> '$.x' FROM t;", want: "`->` reads a JSON document or a string, and column `a` is `int`."},
		{sql: "CREATE VIEW w AS SELECT x.ts ->> '$' FROM t AS x;", want: "`->>` reads a JSON document or a string, and column `ts` is `timestamp`."},
		{sql: "CREATE VIEW w AS SELECT 1 MEMBER OF (a) FROM t;", want: "`MEMBER OF` reads a JSON document or a string, and column `a` is `int`."},
		{sql: "CREATE VIEW w AS SELECT JSON_EXTRACT(dt, '$') FROM t;", want: "`JSON_EXTRACT` reads a JSON document or a string, and column `dt` is `date`."},
		{sql: "CREATE VIEW w AS WITH c AS (SELECT a -> '$' FROM t) SELECT 1;", want: "`->` reads a JSON document or a string, and column `a` is `int`."},
		{sql: "CREATE VIEW w AS SELECT id FROM t WHERE a -> '$' = 1;", want: "`->` reads a JSON document or a string, and column `a` is `int`."},
		{sql: "CREATE VIEW w AS SELECT s -> '$.x', j ->> '$', a MEMBER OF (j) FROM t;"},
		{sql: "CREATE VIEW w AS WITH t AS (SELECT j AS a FROM t) SELECT a -> '$' FROM t;"},
		{sql: "CREATE VIEW w AS SELECT nope.t.a FROM t;", want: "Column `nope.t.a` doesn't exist in the query."},
		{sql: "CREATE VIEW w AS SELECT shop.t.a FROM t AS y;", want: "Column `shop.t.a` doesn't exist in the query."},
		{sql: "CREATE VIEW w AS SELECT shop.y.a, SHOP.y.s FROM t AS y;"},
		{sql: "CREATE VIEW w AS SELECT shop.t.a FROM shop.t;"},
		{sql: "ALTER TABLE t CHANGE COLUMN IF EXISTS missing j INT;\nCREATE VIEW w AS SELECT JSON_EXTRACT(j, '$') FROM t;"},
		{sql: "ALTER TABLE t MODIFY COLUMN IF EXISTS missing INT;\nALTER TABLE t ADD COLUMN missing JSON;\nCREATE VIEW w AS SELECT JSON_EXTRACT(missing, '$') FROM t;"},
	})
}

// TestWalkThroughTiDBCheckColumns pins the CHECK constraint checks
// (walk_checks.go) on a sync whose CHECKs were made while
// tidb_enable_check_constraint was on. The refusals hold with the switch on
// or off (verified on v8.5.5); DROP CHECK beside other subcommands fails
// with 8200 only while it is on, so the walk says nothing of a statement
// that drops a CHECK before the column.
func TestWalkThroughTiDBCheckColumns(t *testing.T) {
	c, err := tidbcatalog.LoadSQL("SET tidb_enable_check_constraint = 1;\nCREATE DATABASE shop;\nUSE shop;\n" +
		"CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT CHECK (guests > 0), nights INT, note VARCHAR(100), " +
		"CONSTRAINT chk_dates CHECK (starts <= ends), CONSTRAINT chk_note CHECK (note <> 'starts'), CONSTRAINT chk_always CHECK (1 = 1));")
	if err != nil {
		t.Fatal(err)
	}
	target := review.Target{Schema: c.DatabaseMetadata("shop"), ForeignKeyChecks: review.SwitchOn}
	for _, tc := range []shapeCase{
		{sql: "ALTER TABLE bookings DROP COLUMN starts;", want: "Check constraint `chk_dates` uses column `starts`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE bookings DROP COLUMN ends, DROP COLUMN starts;", want: "Check constraint `chk_dates` uses column `ends`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE bookings DROP COLUMN starts, DROP CHECK chk_dates;", want: "Check constraint `chk_dates` uses column `starts`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE bookings RENAME COLUMN guests TO people;", want: "Check constraint `bookings_chk_1` uses column `guests`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE bookings RENAME COLUMN note TO memo;", want: "Check constraint `chk_note` uses column `note`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE bookings DROP COLUMN guests;"},
		{sql: "ALTER TABLE bookings DROP COLUMN note;"},
		{sql: "ALTER TABLE bookings DROP COLUMN nights;"},
		{sql: "ALTER TABLE bookings RENAME COLUMN nights TO stays;"},
		// 8200 with the switch on, 3959 with it off: the walk can't tell.
		{sql: "ALTER TABLE bookings DROP CHECK chk_dates, DROP COLUMN starts;"},
		{sql: "ALTER TABLE bookings DROP CONSTRAINT chk_dates, DROP COLUMN starts;"},
		// TiDB takes CHANGE COLUMN and drops the CHECK; the walk follows the
		// table no further.
		{sql: "ALTER TABLE bookings CHANGE COLUMN starts begins DATE;\nALTER TABLE bookings DROP COLUMN ends;"},
		{sql: "DROP TABLE bookings;\nCREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE);\nALTER TABLE bookings DROP COLUMN starts;"},
		{sql: "RENAME TABLE bookings TO b2;\nALTER TABLE b2 DROP COLUMN starts;", statement: 1, want: "Check constraint `chk_dates` uses column `starts`, which can't be dropped or renamed."},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			if tc.want == "" {
				noFinding(t, tc.sql, target)
				return
			}
			oneFinding(t, tc.sql, tc.statement, tc.want, target)
		})
	}
}

// TestWalkThroughTiDBColumnSet pins the checks a table's synced columns
// settle (walk_column_set.go): character set moves, generated columns, the
// integer row handle, and invisible primary keys. Each refusal was
// verified on v8.5.5 in tidb/catalog's simulate.tsv.
func TestWalkThroughTiDBColumnSet(t *testing.T) {
	c, err := tidbcatalog.LoadSQL("CREATE DATABASE shop;\nUSE shop;\n" +
		"CREATE TABLE notes (id INT PRIMARY KEY, title VARCHAR(100), body TEXT, code CHAR(10) CHARACTER SET ascii, tag VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin) DEFAULT CHARSET=latin1;\n" +
		"CREATE TABLE people (id INT PRIMARY KEY, first_name VARCHAR(50) NOT NULL, last_name VARCHAR(50) NOT NULL, full_name VARCHAR(101) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL, name_len INT AS (CHAR_LENGTH(first_name)) STORED, age INT, KEY idx_full (full_name), KEY idx_lower ((LOWER(last_name))));\n" +
		"CREATE TABLE customers (id INT PRIMARY KEY, email VARCHAR(100) NOT NULL, UNIQUE KEY uk_email (email));")
	if err != nil {
		t.Fatal(err)
	}
	target := review.Target{Schema: c.DatabaseMetadata("shop"), ForeignKeyChecks: review.SwitchOn}
	for _, tc := range []shapeCase{
		{sql: "ALTER TABLE notes CONVERT TO CHARACTER SET utf8mb4;", want: "TiDB can't change the character set of column `code` of `notes` from `ascii` to `utf8mb4`."},
		{sql: "ALTER TABLE notes CONVERT TO CHARACTER SET latin1;", want: "TiDB can't change the character set of column `code` of `notes` from `ascii` to `latin1`."},
		{sql: "ALTER TABLE notes CONVERT TO CHARACTER SET utf8mb3;", want: "TiDB can't change the character set of `notes` from `latin1` to `utf8`."},
		{sql: "ALTER TABLE notes MODIFY tag VARCHAR(20);", want: "TiDB can't change the character set of column `tag` of `notes` from `utf8mb4` to `latin1`."},
		{sql: "ALTER TABLE notes MODIFY title VARCHAR(100) CHARACTER SET utf8mb4;"},
		{sql: "ALTER TABLE notes MODIFY body TEXT CHARACTER SET utf8mb4;"},
		{sql: "ALTER TABLE notes DEFAULT CHARSET = utf8mb4;"},
		// The walk follows a table's columns no further than the sync.
		{sql: "ALTER TABLE notes DROP COLUMN code;\nALTER TABLE notes CONVERT TO CHARACTER SET utf8mb4;"},
		{sql: "ALTER TABLE people MODIFY full_name VARCHAR(101);", want: "TiDB can't change whether generated column `full_name` is stored."},
		{sql: "ALTER TABLE people MODIFY age INT AS (id + 1) VIRTUAL;", want: "TiDB can't change whether generated column `age` is stored."},
		{sql: "ALTER TABLE people MODIFY full_name VARCHAR(200) AS (CONCAT(last_name, ', ', first_name)) VIRTUAL;", want: "TiDB can't change the expression of generated column `full_name`, which an index holds."},
		{sql: "ALTER TABLE people CHANGE first_name given_name VARCHAR(50) NOT NULL;", want: "Column `first_name` is read by generated column `full_name`, and can't be renamed."},
		{sql: "ALTER TABLE people MODIFY name_len INT;"},
		// The walk follows which columns an index holds no further than the sync.
		{sql: "DROP INDEX idx_full ON people;\nALTER TABLE people MODIFY full_name VARCHAR(200) AS (CONCAT(last_name, ', ', first_name)) VIRTUAL;"},
		{sql: "CREATE INDEX idx_low2 ON people ((LOWER(first_name)));\nALTER TABLE people CHANGE age years INT;"},
		{sql: "ALTER TABLE people CHANGE age years INT;"},
		// TiDB answers 3108 here, but the walk can't tell when the column
		// of an expression index reading last_name comes first (3837).
		{sql: "ALTER TABLE people CHANGE last_name surname VARCHAR(50) NOT NULL;"},
		{sql: "ALTER TABLE people MODIFY name_len INT AS (CHAR_LENGTH(first_name)) VIRTUAL;", want: "TiDB can't change whether generated column `name_len` is stored."},
		{sql: "CREATE TABLE inv (a INT NOT NULL, b INT NOT NULL UNIQUE, UNIQUE KEY u (a) INVISIBLE);", want: "Unique key `u`, the primary key of `inv`, can't be invisible."},
		{sql: "ALTER TABLE people ALTER INDEX `PRIMARY` INVISIBLE;", want: "Index `PRIMARY` of `people` is its integer row handle, which TiDB keeps no index for."},
		{sql: "CREATE TABLE refunds (id INT PRIMARY KEY, email VARCHAR(50) CHARACTER SET latin1, FOREIGN KEY (email) REFERENCES customers (email));", want: "A foreign key would pair `email` (`varchar` in `latin1`) with `customers`.`email` (`varchar` in `utf8mb4`), types that don't match."},
		{sql: "CREATE TABLE refunds (id INT PRIMARY KEY, email VARCHAR(50) CHARACTER SET utf8mb4, FOREIGN KEY (email) REFERENCES customers (email));"},
		{sql: "CREATE TABLE inv (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE);", want: "Unique key `u`, the primary key of `inv`, can't be invisible."},
		{sql: "CREATE TABLE inv (a INT, PRIMARY KEY (a) INVISIBLE);", want: "The primary key of `inv` can't be invisible."},
		{sql: "CREATE TABLE inv (a INT, UNIQUE KEY u (a) INVISIBLE);"},
		{sql: "CREATE TABLE inv (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE);"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			if tc.want == "" {
				noFinding(t, tc.sql, target)
				return
			}
			oneFinding(t, tc.sql, tc.statement, tc.want, target)
		})
	}
}
