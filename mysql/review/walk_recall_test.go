package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// recallSchema is a synced schema with the facts the recall checks read:
// column types, foreign keys, CHECK constraints, generated columns,
// partitions, and a secondary engine.
func recallSchema() *metadata.DatabaseSchemaMetadata {
	col := func(name, typ string) *metadata.ColumnMetadata {
		return &metadata.ColumnMetadata{Name: name, Type: typ}
	}
	gen := func(name, expr string, stored bool) *metadata.ColumnMetadata {
		g := metadata.GenerationMetadata_TYPE_VIRTUAL
		if stored {
			g = metadata.GenerationMetadata_TYPE_STORED
		}
		return &metadata.ColumnMetadata{Name: name, Type: "int", Generation: &metadata.GenerationMetadata{Type: g, Expression: expr}}
	}
	pk := func(cols ...string) *metadata.IndexMetadata {
		return &metadata.IndexMetadata{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: cols}
	}
	part := func(typ metadata.TablePartitionMetadata_Type, expr string, parts ...[2]string) []*metadata.TablePartitionMetadata {
		var out []*metadata.TablePartitionMetadata
		for _, p := range parts {
			out = append(out, &metadata.TablePartitionMetadata{Name: p[0], Type: typ, Expression: expr, Value: p[1]})
		}
		return out
	}
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name: "t", Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{col("id", "int"), col("a", "int"), col("s", "varchar(20)"), col("ts", "timestamp"), col("g", "geometry"), col("dt", "date")},
					Indexes: []*metadata.IndexMetadata{pk("id"), {Name: "idx_a", Type: "BTREE", Expressions: []string{"a"}}},
				},
				{Name: "parent", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int"), col("email", "varchar(100)")}, Indexes: []*metadata.IndexMetadata{pk("id")}},
				{
					Name: "child", Engine: "InnoDB",
					Columns:     []*metadata.ColumnMetadata{col("id", "int"), col("pid", "int")},
					Indexes:     []*metadata.IndexMetadata{pk("id"), {Name: "fk_pid", Type: "BTREE", Expressions: []string{"pid"}}},
					ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk_pid", Columns: []string{"pid"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}}},
				},
				{
					Name: "chk", Engine: "InnoDB",
					Columns:          []*metadata.ColumnMetadata{col("id", "int"), col("a", "int"), col("b", "int"), col("c", "int")},
					Indexes:          []*metadata.IndexMetadata{pk("id")},
					CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "k1", Expression: "(`a` > 0)"}, {Name: "k2", Expression: "(`b` > `c`)"}},
				},
				{
					Name: "gen", Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{col("id", "int"), col("a", "int"), gen("v", "(`a` + 1)", false), gen("st", "(`a` + 2)", true)},
					Indexes: []*metadata.IndexMetadata{pk("id")},
				},
				{
					Name: "rng", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int")}, Indexes: []*metadata.IndexMetadata{pk("id")},
					Partitions: part(metadata.TablePartitionMetadata_RANGE, "`id`", [2]string{"p0", "10"}, [2]string{"p1", "20"}),
				},
				{
					Name: "lst", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int")}, Indexes: []*metadata.IndexMetadata{pk("id")},
					Partitions: part(metadata.TablePartitionMetadata_LIST, "`id`", [2]string{"p0", "1,2"}, [2]string{"p1", "3,4"}),
				},
				{
					Name: "hsh", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int"), col("u", "int")}, Indexes: []*metadata.IndexMetadata{pk("id")},
					Partitions: part(metadata.TablePartitionMetadata_HASH, "`id`", [2]string{"p0", ""}, [2]string{"p1", ""}, [2]string{"p2", ""}),
				},
				{
					Name: "sub", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int")}, Indexes: []*metadata.IndexMetadata{pk("id")},
					Partitions: []*metadata.TablePartitionMetadata{{
						Name: "p0", Type: metadata.TablePartitionMetadata_RANGE, Expression: "`id`", Value: "10",
						Subpartitions: []*metadata.TablePartitionMetadata{{Name: "p0sp0", Type: metadata.TablePartitionMetadata_HASH, Expression: "`id`"}},
					}},
				},
				{Name: "sec", Engine: "InnoDB", Columns: []*metadata.ColumnMetadata{col("id", "int")}, Indexes: []*metadata.IndexMetadata{pk("id")}, CreateOptions: `SECONDARY_ENGINE="RAPID"`},
			},
			Views:     []*metadata.ViewMetadata{{Name: "v"}},
			Functions: []*metadata.FunctionMetadata{{Name: "f"}},
		}},
	}
}

// recallCase is a change on recallSchema and its one finding, at
// statement, or none when want is empty.
type recallCase struct {
	sql       string
	statement int
	want      string
}

func runRecallCases(t *testing.T, cases []recallCase) {
	t.Helper()
	target := review.Target{Schema: recallSchema(), ForeignKeyChecks: review.SwitchOn}
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

// TestWalkThroughMySQLViewQueries pins the walk's findings on a view's query,
// as MySQL resolves it (walk_lower_query.go). Each refusal was verified on
// mysql:8.0.46; an empty want is a change the walk lets through, which may
// still fail for a reason it doesn't follow.
func TestWalkThroughMySQLViewQueries(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "CREATE VIEW v2 AS SELECT nope FROM t;", want: "Column `nope` doesn't exist in `t`."},
		{sql: "CREATE VIEW v2 AS SELECT x.a FROM t AS y;", want: "Column `x.a` doesn't exist in the query."},
		{sql: "CREATE VIEW v2 AS SELECT a FROM nosuch;", want: "Table `nosuch` doesn't exist."},
		{sql: "CREATE VIEW v2 AS SELECT ROW_NUMBER() OVER w FROM t;", want: "Window `w` isn't defined in the query's `WINDOW` clause."},
		{sql: "CREATE VIEW v2 AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN 1 FOLLOWING AND CURRENT ROW) FROM t;", want: "A window frame can't end before it starts."},
		{sql: "CREATE VIEW v2 AS SELECT ROW_NUMBER() OVER (GROUPS UNBOUNDED PRECEDING) FROM t;", want: "MySQL doesn't support a `GROUPS` window frame."},
		{sql: "CREATE VIEW v2 AS SELECT a, a FROM t;", want: "This statement names column `a` twice."},
		{sql: "CREATE VIEW v2 AS SELECT @x;", want: "A view's query can't read a variable."},
		{sql: "CREATE VIEW v2 AS SELECT DATE 'x';", want: "`DATE 'x'` isn't a valid date literal."},
		{sql: "CREATE VIEW v2 AS SELECT LINESTRING(a) FROM t;", want: "`LINESTRING` takes geometry values, and its argument is `int`."},
		{sql: "CREATE VIEW v2 AS SELECT JSON_VALUE(a, '$') FROM t;", want: "`JSON_VALUE` reads a JSON document or a string, and its argument is `int`."},
		{sql: "CREATE VIEW v2 AS SELECT CAST(dt AT TIME ZONE 'UTC' AS DATETIME) FROM t;", want: "`CAST(... AT TIME ZONE ...)` takes a `TIMESTAMP`, and its argument is `date`."},
		{sql: "CREATE VIEW v2 AS SELECT CAST(ts AT TIME ZONE 'utc' AS DATETIME) FROM t;", want: "`CAST(... AT TIME ZONE ...)` takes only `'UTC'` or `INTERVAL '+00:00'`, not `'utc'`."},
		{sql: "CREATE VIEW v2 AS SELECT DEFAULT(a) FROM t;", want: "A view can't take `DEFAULT(a)` as a column."},
		{sql: "CREATE VIEW v2 AS VALUES ROW();", want: "Each row of a `VALUES` query outside an `INSERT` needs at least one column."},
		{sql: "CREATE VIEW v2 AS SELECT a FROM t LIMIT n;", want: "`LIMIT` takes a number, not `n`."},
		{sql: "CREATE VIEW v2 AS SELECT a INTO @x FROM t;", want: "A view's query can't hold an `INTO` clause."},
		{sql: "CREATE VIEW v2 AS SELECT s LIKE 'a' ESCAPE s FROM t;", want: "`ESCAPE` takes a constant, not a column."},
		{sql: "CREATE VIEW v2 AS SELECT GROUPING(a) FROM t GROUP BY a;", want: "`GROUPING()` needs a `GROUP BY ... WITH ROLLUP` in its query block."},
		{sql: "ALTER VIEW v AS SELECT nope FROM t;", want: "Column `nope` doesn't exist in `t`."},
		{sql: "CREATE VIEW v2 AS SELECT a FROM t WHERE a > 0;"},
		{sql: "CREATE VIEW v2 AS SELECT ROW_NUMBER() OVER w FROM t WINDOW W AS ();"},
		{sql: "CREATE VIEW v2 AS SELECT t.a, a AS b, shop.t.id FROM t;"},
		{sql: "CREATE VIEW v2 (x, y) AS SELECT a, a FROM t;"},
		{sql: "CREATE VIEW v2 AS WITH c AS (SELECT 1 AS z) SELECT z FROM c;"},
		{sql: "CREATE VIEW v2 AS SELECT CAST(ts AT TIME ZONE 'UTC' AS DATETIME) FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT 1 FROM DUAL;"},
		{sql: "CREATE VIEW v2 AS SELECT a, GROUPING(a) FROM t GROUP BY a WITH ROLLUP;"},
		{sql: "CREATE VIEW v2 AS SELECT LINESTRING(g), JSON_VALUE(s, '$') FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT a FROM t JOIN parent ON parent.id = t.a WHERE nope = 1;"},
		{sql: "CREATE VIEW v2 AS SELECT (SELECT nope FROM parent) FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT * FROM t AS x FOR UPDATE OF x;"},
		{sql: "CREATE VIEW v2 AS WITH t AS (SELECT 1 AS z) SELECT z FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT GET_FORMAT(DATE, 'USA'), CHAR(65 USING utf8mb4), TRIM(BOTH 'x' FROM s), EXTRACT(YEAR FROM dt), DATE_ADD(dt, INTERVAL 1 DAY), CONVERT(s USING latin1), POSITION('a' IN s), TIMESTAMPDIFF(DAY, dt, dt), CURRENT_USER, LOCALTIME FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT NTILE(2) OVER (), LEAD(a, 1) OVER (ORDER BY a) FROM t;"},
		{sql: "CREATE VIEW v2 AS SELECT a FROM t WHERE a = DEFAULT(a);"},
		{sql: "CREATE VIEW v2 AS SELECT * FROM JSON_TABLE('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt;"},
	})
}

// TestWalkThroughMySQLTableOptions pins the walk's findings on character
// sets, collations, engines, and the secondary engine
// (walk_lower_options.go). Each refusal was verified on mysql:8.0.46; an
// empty want is a change the walk lets through, which may still fail for a
// reason it doesn't follow.
func TestWalkThroughMySQLTableOptions(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "ALTER TABLE t ENGINE = nope;", want: "MySQL has no storage engine `nope`."},
		{sql: "SET sql_mode = '';\nALTER TABLE t ENGINE = nope;"},
		{sql: "ALTER TABLE t CONVERT TO CHARACTER SET nope;", want: "Character set `nope` doesn't exist."},
		{sql: "ALTER TABLE t CONVERT TO CHARACTER SET latin1 COLLATE utf8mb4_bin;", want: "Collation `utf8mb4_bin` isn't one of character set `latin1`."},
		{sql: "CREATE TABLE x (a VARCHAR(10) CHARACTER SET latin1 COLLATE utf8mb4_bin);", want: "Collation `utf8mb4_bin` isn't one of character set `latin1`."},
		{sql: "CREATE TABLE x (a VARCHAR(10) CHARACTER SET utf8 COLLATE utf8mb3_bin) CHARSET utf8mb4 COLLATE utf8mb4_bin ENGINE heap;"},
		{sql: "ALTER TABLE nosuch ENGINE = nope;", want: "MySQL has no storage engine `nope`."},
		{sql: "ALTER TABLE t SECONDARY_LOAD;", want: "Table `t` has no secondary engine."},
		{sql: "ALTER TABLE sec SECONDARY_LOAD;"},
		{sql: "ALTER TABLE t SECONDARY_ENGINE = RAPID;\nALTER TABLE t SECONDARY_LOAD;"},
		{sql: "CREATE TABLE x (a int);\nALTER TABLE x SECONDARY_UNLOAD;", statement: 1, want: "Table `x` has no secondary engine."},
		{sql: "ALTER TABLE t ADD INDEX i2 (s), ALGORITHM = INSTANT;", want: "InnoDB can't add a key with `ALGORITHM=INSTANT`."},
	})
}

// TestWalkThroughMySQLAlgorithmAndLock pins the walk's findings on ALGORITHM
// and LOCK (walk_lower_algorithm.go). Each refusal was verified on
// mysql:8.0.46; an empty want is a change the walk lets through, which may
// still fail for a reason it doesn't follow.
func TestWalkThroughMySQLAlgorithmAndLock(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "ALTER TABLE t ADD COLUMN z INT, ALGORITHM = INSTANT, LOCK = NONE;", want: "`ALGORITHM=INSTANT` can't take `LOCK=NONE`."},
		{sql: "ALTER TABLE t MODIFY a BIGINT, LOCK = NONE;", want: "InnoDB can't change a column's type with `LOCK=NONE`."},
		{sql: "ALTER TABLE t DROP PRIMARY KEY, ALGORITHM = INPLACE;", want: "InnoDB can't drop the primary key without adding one with `ALGORITHM=INPLACE`."},
		{sql: "ALTER TABLE t ADD CONSTRAINT c5 CHECK (a > 0), ALGORITHM = INPLACE;", want: "InnoDB can't add a `CHECK` constraint with `ALGORITHM=INPLACE`."},
		{sql: "CREATE INDEX i2 ON t (s) ALGORITHM = INSTANT;", want: "InnoDB can't add an index with `ALGORITHM=INSTANT`."},
		{sql: "ALTER TABLE t ADD COLUMN z INT, ALGORITHM = INSTANT;"},
		{sql: "ALTER TABLE t ADD CONSTRAINT c5 CHECK (a > 0) NOT ENFORCED, ALGORITHM = INPLACE;"},
		{sql: "ALTER TABLE t ADD CONSTRAINT c5 CHECK (a > 0) NOT ENFORCED, ALGORITHM = INSTANT;"},
		{sql: "ALTER TABLE t ADD CONSTRAINT c5 CHECK (a > 0) NOT ENFORCED, LOCK = NONE;"},
		{sql: "ALTER TABLE t MODIFY s VARCHAR(40), LOCK = NONE;"},
		{sql: "ALTER TABLE t ADD INDEX i2 (s), LOCK = NONE;"},
		{sql: "ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (id, a), ALGORITHM = INPLACE;"},
		{sql: "DROP TABLE parent;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
	})
}

// TestWalkThroughMySQLForeignKeyDependents pins the walk's findings on
// foreign keys of other tables and their column types (walk_lower_fk.go).
// Each refusal was verified on mysql:8.0.46; an empty want is a change the
// walk lets through, which may still fail for a reason it doesn't follow.
func TestWalkThroughMySQLForeignKeyDependents(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "DROP TABLE parent, child;"},
		{sql: "SELECT f();\nDROP TEMPORARY TABLE IF EXISTS parent;\nDROP TABLE parent, child;"},
		{sql: "DROP TABLE IF EXISTS parent, nosuch;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "TRUNCATE TABLE parent;", want: "Table `parent` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "TRUNCATE TABLE child;"},
		{sql: "ALTER TABLE parent ENGINE = MyISAM;", want: "Table `parent` is in foreign key `fk_pid` of `child`, and engine `MyISAM` keeps no foreign keys."},
		{sql: "ALTER TABLE child ENGINE = MyISAM;", want: "Table `child` has foreign key `fk_pid`, and engine `MyISAM` keeps no foreign keys."},
		{sql: "SET foreign_key_checks = 0;\nDROP TABLE parent;"},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid;\nDROP TABLE parent;"},
		{sql: "ALTER TABLE parent MODIFY id BIGINT;", want: "Foreign key `fk_pid` pairs `pid` (`int`) with `parent`.`id` (`bigint`), types that don't match."},
		{sql: "ALTER TABLE child MODIFY pid VARCHAR(10);", want: "Foreign key `fk_pid` pairs `pid` (`varchar`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "ALTER TABLE child ADD CONSTRAINT f2 FOREIGN KEY (id) REFERENCES parent (nope);", want: "Referenced column `nope` doesn't exist in `parent`."},
		{sql: "ALTER TABLE child ADD CONSTRAINT f2 FOREIGN KEY (id) REFERENCES parent (email);", want: "Foreign key `f2` pairs `id` (`int`) with `parent`.`email` (`varchar`), types that don't match."},
		{sql: "CREATE TABLE x (id BIGINT, FOREIGN KEY (id) REFERENCES parent (id));", want: "A foreign key pairs `id` (`bigint`) with `parent`.`id` (`int`), types that don't match."},
		{sql: "ALTER TABLE child DROP FOREIGN KEY fk_pid, MODIFY pid BIGINT;"},
		{sql: "ALTER TABLE parent DROP COLUMN id;", want: "Column `id` is referenced by foreign key `fk_pid` of `child`."},
		{sql: "ALTER TABLE parent MODIFY id INT UNSIGNED;"},
		{sql: "ALTER TABLE rng ADD PARTITION (PARTITION p2 VALUES LESS THAN (15));", want: "Partition `p2`'s bound must be greater than the partition's before it."},
		{sql: "ALTER TABLE rng ADD PARTITION (PARTITION p2 VALUES IN (30));", want: "Partition `p2` can't take `VALUES IN`: only `LIST` partitioning takes it."},
	})
}

// TestWalkThroughMySQLPartitionOperations pins the walk's findings on
// partition clauses and operations (walk_lower_partition.go). Each refusal
// was verified on mysql:8.0.46; an empty want is a change the walk lets
// through, which may still fail for a reason it doesn't follow.
func TestWalkThroughMySQLPartitionOperations(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "ALTER TABLE rng DROP PARTITION nope;", want: "Partition `nope` doesn't exist in `rng`."},
		{sql: "ALTER TABLE sub TRUNCATE PARTITION p0sp0;"},
		{sql: "ALTER TABLE rng TRUNCATE PARTITION p0, nope;", want: "Partition `nope` doesn't exist in `rng`."},
		{sql: "ALTER TABLE lst ADD PARTITION (PARTITION p2 VALUES IN (4, 5));", want: "The partitions' `VALUES IN` name `4` twice."},
		{sql: "ALTER TABLE hsh DROP PARTITION p0;", want: "`DROP PARTITION` works only on `RANGE` or `LIST` partitioning."},
		{sql: "ALTER TABLE hsh COALESCE PARTITION 3;", want: "Cannot remove all partitions from `hsh`."},
		{sql: "ALTER TABLE hsh COALESCE PARTITION 0;", want: "`COALESCE PARTITION` must remove at least one partition."},
		{sql: "ALTER TABLE hsh REORGANIZE PARTITION p0, p1 INTO (PARTITION p01);", want: "`REORGANIZE PARTITION` of `HASH` or `KEY` partitioning must keep the number of partitions."},
		{sql: "ALTER TABLE hsh ADD UNIQUE KEY uk (u);", want: "Every unique key of `hsh` must include all columns in its partitioning function."},
		{sql: "ALTER TABLE rng EXCHANGE PARTITION p0 WITH TABLE hsh;", want: "Table `hsh` is partitioned, and a partition can be exchanged only with a table that isn't."},
		{sql: "CREATE TABLE x (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p0 VALUES LESS THAN (20));", want: "This statement names partition `p0` twice."},
		{sql: "CREATE TABLE x (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (5));", want: "Partition `p1`'s bound must be greater than the partition's before it."},
		{sql: "CREATE TABLE x (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES LESS THAN (10));", want: "Partition `p0` can't take `VALUES LESS THAN`: only `RANGE` partitioning takes it."},
		{sql: "CREATE TABLE x (a VARCHAR(10)) PARTITION BY HASH (a);", want: "Partition column `a` is `varchar`, and the partitioning function must be an integer."},
		{sql: "CREATE TABLE x (a TEXT) PARTITION BY KEY (a);", want: "Partition column `a` is `text`, which no partitioning can read."},
		{sql: "ALTER TABLE rng ADD PARTITION (PARTITION p2 VALUES LESS THAN (30));"},
		{sql: "ALTER TABLE lst ADD PARTITION (PARTITION p2 VALUES IN (5));"},
		{sql: "ALTER TABLE hsh COALESCE PARTITION 2;"},
		{sql: "ALTER TABLE hsh ADD UNIQUE KEY uk (id, u);"},
		{sql: "ALTER TABLE rng ADD PARTITION (PARTITION p2 VALUES LESS THAN (30));\nALTER TABLE rng ADD PARTITION (PARTITION p3 VALUES LESS THAN (25));"},
		{sql: "ALTER TABLE chk DROP COLUMN b;", want: "Check constraint `k2` uses column `b`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE chk RENAME COLUMN a TO a2;", want: "Check constraint `k1` uses column `a`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE chk DROP COLUMN a;"},
	})
}

// TestWalkThroughMySQLColumnsAndChecks pins the walk's findings on CHECK
// constraints, generated columns, keys, and column types
// (walk_lower_columns.go), and event schedules. Each refusal was verified on
// mysql:8.0.46; an empty want is a change the walk lets through, which may
// still fail for a reason it doesn't follow.
func TestWalkThroughMySQLColumnsAndChecks(t *testing.T) {
	runRecallCases(t, []recallCase{
		{sql: "ALTER TABLE chk DROP CHECK k2, DROP COLUMN b;"},
		{sql: "ALTER TABLE chk DROP CHECK k1;\nALTER TABLE chk ADD CONSTRAINT k1 CHECK (b > 0);\nALTER TABLE chk RENAME COLUMN a TO a2;"},
		{sql: "ALTER TABLE chk DROP CHECK k1, ADD CONSTRAINT k1 CHECK (b > 0);"},
		{sql: "CREATE TABLE x (a VARCHAR(20000) COLLATE latin1_bin) DEFAULT CHARSET = utf8mb4;"},
		{sql: "CREATE TABLE x (a VARCHAR(20000) COLLATE utf8mb4_bin) DEFAULT CHARSET = latin1;", want: "Column `a` is longer than 16383 characters, MySQL's limit for its type; use `TEXT` or `BLOB`."},
		{sql: "ALTER TABLE chk DROP COLUMN b, DROP COLUMN c;", want: "Check constraint `k2` uses column `b`, which can't be dropped or renamed."},
		{sql: "ALTER TABLE chk MODIFY b BIGINT;"},
		{sql: "ALTER TABLE t ADD CONSTRAINT K1 CHECK (a > 0);", want: "Check constraint `K1` already exists in `chk`."},
		{sql: "CREATE TABLE x (a INT, CONSTRAINT k2 CHECK (a > 0));", want: "Check constraint `k2` already exists in `chk`."},
		{sql: "ALTER TABLE chk DROP CHECK k1;\nALTER TABLE t ADD CONSTRAINT k1 CHECK (a > 0);"},
		{sql: "ALTER TABLE gen MODIFY v INT;", want: "Column `v` can't change between a virtual generated column and a stored or plain one."},
		{sql: "ALTER TABLE gen ALTER COLUMN v SET DEFAULT 1;", want: "Generated column `v` can't take a default."},
		{sql: "ALTER TABLE gen MODIFY st INT;"},
		{sql: "ALTER TABLE gen MODIFY a INT AS (id + 1) VIRTUAL;", want: "Column `a` can't change between a virtual generated column and a stored or plain one."},
		{sql: "ALTER TABLE gen MODIFY st INT;\nALTER TABLE gen ALTER COLUMN st SET DEFAULT 1;"},
		{sql: "CREATE TABLE x (a INT, b DOUBLE AS (RAND()));", want: "Generated column `b` can't call `RAND()`, whose value isn't deterministic."},
		{sql: "ALTER TABLE t ADD CONSTRAINT c9 CHECK (a < UUID_SHORT());", want: "A `CHECK` constraint can't call `UUID_SHORT()`, whose value isn't deterministic."},
		{sql: "ALTER TABLE t ADD INDEX fi ((NOW()));", want: "A functional key part can't call `NOW()`, whose value isn't deterministic."},
		{sql: "CREATE TABLE x (a INT INVISIBLE);", want: "Table `x` must have at least one visible column."},
		{sql: "CREATE TABLE x (a ENUM('s','s'));", want: "Column `a` lists `'s'` twice in its `ENUM`."},
		{sql: "CREATE TABLE x (a VARCHAR(20000)) CHARSET utf8mb4;", want: "Column `a` is longer than 16383 characters, MySQL's limit for its type; use `TEXT` or `BLOB`."},
		{sql: "CREATE TABLE x (a VARCHAR(16383)) CHARSET utf8mb4;"},
		{sql: "SET sql_mode = '';\nCREATE TABLE x (a ENUM('s','s'));"},
		{sql: "ALTER TABLE t ADD COLUMN g2 INT AS (nope + 1);", want: "Column `nope` doesn't exist in `t`."},
		{sql: "ALTER TABLE t ADD INDEX fi ((nope + 1));", want: "Column `nope` doesn't exist in `t`."},
		{sql: "ALTER TABLE t ALTER INDEX `PRIMARY` INVISIBLE;", want: "A primary key can't be invisible."},
		{sql: "CREATE TABLE x (a INT, PRIMARY KEY ((a + 1)));", want: "A primary key can't have a functional key part."},
		{sql: "CREATE TABLE x (a INT AUTO_INCREMENT KEY, b INT AS (a + 1));", want: "Generated column `b` can't read `AUTO_INCREMENT` column `a`."},
		{sql: "CREATE TABLE x (a INT, b INT AS (c + 1), c INT AS (a + 1));", want: "Generated column `b` can't read generated column `c`, defined after it."},
		{sql: "CREATE TABLE x (a INT, b INT AS (c + 1), c INT);"},
		{sql: "CREATE TABLE x (a INT AUTO_INCREMENT KEY, CHECK (a > 0));", want: "A `CHECK` constraint can't read `AUTO_INCREMENT` column `a`."},
		{sql: "ALTER TABLE t ADD COLUMN a ENUM('s','s');", want: "Column `a` already exists in `t`."},
		{sql: "CREATE EVENT e2 ON SCHEDULE EVERY 0 DAY DO COMMIT;", want: "An event's `EVERY` interval must be positive."},
		{sql: "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY STARTS '2031-01-01' ENDS '2030-01-01' DO COMMIT;", want: "An event's `ENDS` comes before its `STARTS`."},
		{sql: "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01' ENDS '2031-01-01 00:00:00' DO COMMIT;"},
	})
}
