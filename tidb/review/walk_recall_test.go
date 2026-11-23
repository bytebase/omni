package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
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
		{"ALTER TABLE t PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10));\nALTER TABLE t COALESCE PARTITION 1;", "`COALESCE PARTITION` works only on `HASH` or `KEY` partitioning, and `t` has `RANGE`."},
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
		"ALTER TABLE latin CHARSET utf8;",
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
		"CREATE VIEW w AS SELECT f(a) FROM t;",
	} {
		t.Run("ok/"+sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}
