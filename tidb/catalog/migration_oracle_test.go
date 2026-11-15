package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// migrationOraclePairs add, to diffOraclePairs, changes whose statements
// depend on one another.
var migrationOraclePairs = []struct{ name, a, b, syncDrops string }{
	{"column added between", "CREATE TABLE t (a int, c int)", "CREATE TABLE t (a int, b int, c int)", ""},
	{"columns added first and chained", "CREATE TABLE t (c int)", "CREATE TABLE t (a int, b int, c int, d int, e int)", ""},
	{"column dropped under a composite key", "CREATE TABLE t (a int, b int, KEY k (a, b))", "CREATE TABLE t (a int)", ""},
	{"column dropped and key narrowed", "CREATE TABLE t (a int, b int, KEY k (a, b))", "CREATE TABLE t (a int, KEY k (a))", ""},
	{"key over a changed column", "CREATE TABLE t (a int, KEY k (a))", "CREATE TABLE t (a bigint, KEY k (a, a2), a2 int)", ""},
	{"column under an expression key", "CREATE TABLE t (a varchar(10), KEY k ((lower(a))))", "CREATE TABLE t (a varchar(20), KEY k ((lower(a))))", ""},
	{"column under a foreign key", "CREATE TABLE p (id int PRIMARY KEY NONCLUSTERED); CREATE TABLE c (id int PRIMARY KEY, pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))",
		"CREATE TABLE p (id bigint PRIMARY KEY NONCLUSTERED); CREATE TABLE c (id int PRIMARY KEY, pid bigint, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", ""},
	{"foreign key to a new table", "CREATE TABLE c (id int PRIMARY KEY, pid int)",
		"CREATE TABLE c (id int PRIMARY KEY, pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id)); CREATE TABLE p (id int PRIMARY KEY)", ""},
	{"tables referencing each other dropped", "CREATE TABLE a (id int PRIMARY KEY, bid int); CREATE TABLE b (id int PRIMARY KEY, aid int, CONSTRAINT fb FOREIGN KEY (aid) REFERENCES a (id)); ALTER TABLE a ADD CONSTRAINT fa FOREIGN KEY (bid) REFERENCES b (id)", "CREATE TABLE keep (a int)", ""},
	{"new tables referencing each other", "CREATE TABLE keep (a int)", "CREATE TABLE keep (a int); CREATE TABLE a (id int PRIMARY KEY, bid int); CREATE TABLE b (id int PRIMARY KEY, aid int, CONSTRAINT fb FOREIGN KEY (aid) REFERENCES a (id)); ALTER TABLE a ADD CONSTRAINT fa FOREIGN KEY (bid) REFERENCES b (id)", ""},
	{"foreign key dropped with its key", "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (pid int)", ""},
	{"nonclustered primary key changed", "CREATE TABLE t (a int, b int, PRIMARY KEY (a) NONCLUSTERED)", "CREATE TABLE t (a int, b int, PRIMARY KEY (a, b) NONCLUSTERED)", ""},
	{"views read new views", "CREATE TABLE t (a int)", "CREATE TABLE t (a int); CREATE VIEW vc AS SELECT a FROM t; CREATE VIEW vb AS SELECT a FROM vc; CREATE VIEW va AS SELECT a FROM vb", ""},
	{"view reads a new column", "CREATE TABLE t (a int); CREATE VIEW v AS SELECT a FROM t", "CREATE TABLE t (a int, b int); CREATE VIEW v AS SELECT a, b FROM t", ""},
	{"table becomes a view", "CREATE TABLE base (a int); CREATE TABLE x (a int)", "CREATE TABLE base (a int); CREATE VIEW x AS SELECT a FROM base", ""},
	{"view becomes a table", "CREATE TABLE base (a int); CREATE VIEW x AS SELECT a FROM base", "CREATE TABLE base (a int); CREATE TABLE x (a int)", ""},
	{"table becomes a sequence", "CREATE TABLE x (a int)", "CREATE SEQUENCE x", ""},
	{"sequence options", "CREATE SEQUENCE s START 10 MINVALUE 5", "CREATE SEQUENCE s START 50 MINVALUE 20 MAXVALUE 1000 INCREMENT 3 CACHE 10 CYCLE", ""},
	{"sequence read by a new table", "CREATE TABLE keep (a int)", "CREATE TABLE keep (a int); CREATE SEQUENCE s; CREATE TABLE t (id int DEFAULT (NEXTVAL(s)))", "the database a default's sequence is in, which TiDB writes and the loader keeps"},
	{"range partitions appended", "CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", "CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION pm VALUES LESS THAN MAXVALUE)", ""},
	{"list partition appended", "CREATE TABLE t (a int) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1, 2))", "CREATE TABLE t (a int) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1, 2), PARTITION p1 VALUES IN (3))", ""},
	{"partitioning added", "CREATE TABLE t (a int)", "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 3", ""},
	{"partitioning removed", "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 3", "CREATE TABLE t (a int)", ""},
	{"partitioning replaced", "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 3", "CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION pm VALUES LESS THAN MAXVALUE)", ""},
	{"ttl added", "CREATE TABLE t (a datetime)", "CREATE TABLE t (a datetime) TTL = a + INTERVAL 1 DAY TTL_ENABLE = 'OFF' TTL_JOB_INTERVAL = '2h'", "TTL"},
	{"ttl removed", "CREATE TABLE t (a datetime) TTL = a + INTERVAL 1 DAY", "CREATE TABLE t (a datetime)", "TTL"},
	{"options together", "CREATE TABLE t (a int, b varchar(5))", "CREATE TABLE t (a int, b varchar(5)) COMMENT 'it''s \\\\ x' SHARD_ROW_ID_BITS=3 AUTO_ID_CACHE=10", "AUTO_ID_CACHE"},
	{"checks replaced", "CREATE TABLE t (a int, b int, CONSTRAINT c1 CHECK (a > 0))", "CREATE TABLE t (a int, b int, CONSTRAINT c1 CHECK (a > 1), CONSTRAINT c2 CHECK (b > 0), CONSTRAINT c3 CHECK (b < 10))", ""},
	{"column defaults and comments", "CREATE TABLE t (a int, b varchar(10), c timestamp)",
		"CREATE TABLE t (a int NOT NULL DEFAULT 7 COMMENT 'it''s', b varchar(10) DEFAULT 'x''y\\\\z', c timestamp(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3))", ""},
	{"every column replaced", "CREATE TABLE t (a int, b int)", "CREATE TABLE t (c int, d int)", ""},
	{"generated column dropped with its input", "CREATE TABLE t (a int, z int AS (a + 1) VIRTUAL, y int AS (z + 1) VIRTUAL, keep int)", "CREATE TABLE t (keep int)", ""},
	{"view aliased like the view that reads it", "CREATE TABLE t (a int)", "CREATE TABLE t (a int); CREATE VIEW vb AS SELECT a AS va FROM t; CREATE VIEW va AS SELECT va FROM vb", ""},
	{"partitioning removed for a unique key", "CREATE TABLE t (a int, b int) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (a int, b int, UNIQUE KEY uk (b))", ""},
	{"partitioning replaced with a unique key", "CREATE TABLE t (a int, b int) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (a int, b int, UNIQUE KEY uk (b)) PARTITION BY HASH (b) PARTITIONS 3", ""},
	{"generated column retargeted off a dropped column", "CREATE TABLE t (a int, z int AS (a + 1) VIRTUAL, keep int)", "CREATE TABLE t (z int AS (keep + 1) VIRTUAL, keep int)", ""},
	{"ttl on a newly added column", "CREATE TABLE t (id int)", "CREATE TABLE t (id int, ts datetime) TTL = ts + INTERVAL 1 DAY", "TTL"},
	{"ordinary to global temporary", "CREATE TABLE t (a int)", "CREATE GLOBAL TEMPORARY TABLE t (a int) ON COMMIT DELETE ROWS", "global temporary (no metadata field)"},
	{"global temporary to ordinary", "CREATE GLOBAL TEMPORARY TABLE t (a int) ON COMMIT DELETE ROWS", "CREATE TABLE t (a int)", "global temporary (no metadata field)"},
	{"ttl column dropped with its ttl", "CREATE TABLE t (id int, ts datetime) TTL = ts + INTERVAL 1 DAY", "CREATE TABLE t (id int)", "TTL"},
	{"ttl removed keeping its column", "CREATE TABLE t (id int, ts datetime) TTL = ts + INTERVAL 1 DAY", "CREATE TABLE t (id int, ts datetime)", "TTL"},
	{"ttl retargeted off a dropped column", "CREATE TABLE t (id int, ts datetime, ts2 datetime) TTL = ts + INTERVAL 1 DAY", "CREATE TABLE t (id int, ts2 datetime) TTL = ts2 + INTERVAL 2 DAY", "TTL"},
	{"fk backing index widened", "CREATE TABLE par (id int PRIMARY KEY); CREATE TABLE ch (a int, b int, KEY k (a), CONSTRAINT fk FOREIGN KEY (a) REFERENCES par (id))", "CREATE TABLE par (id int PRIMARY KEY); CREATE TABLE ch (a int, b int, KEY k (a, b), CONSTRAINT fk FOREIGN KEY (a) REFERENCES par (id))", ""},
	{"fk backing index dropped for another", "CREATE TABLE par (id int PRIMARY KEY); CREATE TABLE ch (a int, b int, KEY k (a), CONSTRAINT fk FOREIGN KEY (a) REFERENCES par (id))", "CREATE TABLE par (id int PRIMARY KEY); CREATE TABLE ch (a int, b int, KEY k2 (a, b), CONSTRAINT fk FOREIGN KEY (a) REFERENCES par (id))", ""},
	{"generated column added", "CREATE TABLE t (a int)", "CREATE TABLE t (a int, b int AS (a * 2) VIRTUAL, KEY kb (b))", ""},

	// Changes TiDB makes only by rebuilding the table.
	{"stored generated column changed", "CREATE TABLE t (a int, b int AS (a + 1) STORED)", "CREATE TABLE t (a int, b int AS (a + 2) STORED)", ""},
	{"stored generated column added", "CREATE TABLE t (a int)", "CREATE TABLE t (a int, b int AS (a + 1) STORED)", ""},
	{"auto increment added", "CREATE TABLE t (id int PRIMARY KEY)", "CREATE TABLE t (id int PRIMARY KEY AUTO_INCREMENT)", ""},
	{"auto increment removed", "CREATE TABLE t (id int PRIMARY KEY AUTO_INCREMENT)", "CREATE TABLE t (id int PRIMARY KEY)", ""},
	{"auto random added", "CREATE TABLE t (id bigint PRIMARY KEY)", "CREATE TABLE t (id bigint PRIMARY KEY AUTO_RANDOM(5))", ""},
	{"charset narrowed", "CREATE TABLE t (a varchar(5))", "CREATE TABLE t (a varchar(5) CHARSET latin1)", ""},
	{"clustered key column widened", "CREATE TABLE t (id int, v varchar(5), PRIMARY KEY (id, v))", "CREATE TABLE t (id bigint, v varchar(10), PRIMARY KEY (id, v))", ""},
	{"clustered key column narrowed", "CREATE TABLE t (id bigint PRIMARY KEY)", "CREATE TABLE t (id int PRIMARY KEY)", ""},
	{"column under a generated column retyped", "CREATE TABLE t (a int, b int AS (a + 1) VIRTUAL)", "CREATE TABLE t (a bigint, b int AS (a + 1) VIRTUAL)", ""},
	{"clustered primary key added", "CREATE TABLE t (id int NOT NULL)", "CREATE TABLE t (id int NOT NULL, PRIMARY KEY (id) CLUSTERED)", ""},
}

// migrationRefusals lists the pairs whose plan TiDB refuses, and why. The
// refused statement must carry a Warning, and an entry whose plan converges
// fails, so a fix removes it.
var migrationRefusals = map[string]string{
	"clustered primary key":                   "a clustered primary key cannot be dropped (8200)",
	"stored generated column changed":         "a stored generated column cannot be changed (3106)",
	"stored generated column added":           "a stored generated column cannot be added (3106)",
	"auto increment added":                    "AUTO_INCREMENT cannot be set on a column (8200)",
	"auto increment removed":                  "AUTO_INCREMENT is removed only with tidb_allow_remove_auto_inc on (8200)",
	"auto random added":                       "AUTO_RANDOM cannot be added to a column (8216)",
	"charset narrowed":                        "a charset changes only from utf8 to utf8mb4 (8200)",
	"clustered key column narrowed":           "a clustered primary key column's type changes only where no row is rewritten (8200)",
	"clustered primary key added":             "a clustered primary key cannot be added (8200)",
	"column under a generated column retyped": "a column a generated column reads keeps its type (3106)",
}

// TestOracle_MigrationConverges applies, for each pair, the plan from A to B
// to a database holding A and requires TiDB to store what it stores for B
// (storedSchema). Plans come from catalogs built by Exec and from catalogs
// loaded from the snapshots Bytebase's sync takes.
func TestOracle_MigrationConverges(t *testing.T) {
	ctr, done := startContainer(t)
	defer done()
	exec := func(t *testing.T, stmt string) {
		t.Helper()
		if _, err := ctr.db.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("TiDB: %q: %v", stmt, err)
		}
	}
	exec(t, "SET GLOBAL tidb_enable_check_constraint = ON")
	defer exec(t, "SET GLOBAL tidb_enable_check_constraint = OFF")
	defer exec(t, "USE test")
	setUp := func(t *testing.T, db, script string) {
		t.Helper()
		exec(t, "DROP DATABASE IF EXISTS "+db)
		exec(t, "CREATE DATABASE "+db)
		exec(t, "USE "+db)
		exec(t, "SET foreign_key_checks = 0")
		defer exec(t, "SET foreign_key_checks = 1")
		for _, stmt := range splitStatements(strings.ReplaceAll(script, diffOracleDBB+".", db+".")) {
			if stmt = strings.TrimSpace(stmt); stmt != "" {
				exec(t, stmt)
			}
		}
	}
	type pair struct{ name, a, b, syncDrops string }
	var pairs []pair
	for _, p := range diffOraclePairs {
		pairs = append(pairs, pair{p.name, p.a, p.b, p.syncDrops})
	}
	for _, p := range migrationOraclePairs {
		pairs = append(pairs, pair{p.name, p.a, p.b, p.syncDrops})
	}
	const target, applied = "omni_mig_b", "omni_mig_a"
	for _, p := range pairs {
		t.Run(p.name, func(t *testing.T) {
			setUp(t, target, p.b)
			want := storedSchema(t, ctr, target)
			setUp(t, applied, p.a)
			snapA := syncTiDBSnapshotWithSequences(t, ctr, applied)
			snapB := syncTiDBSnapshotWithSequences(t, ctr, target)

			paths := []struct {
				name     string
				from, to *Catalog
			}{{"Exec", execCatalogIn(t, applied, p.a), execCatalogIn(t, applied, p.b)}}
			if p.syncDrops == "" {
				paths = append(paths, struct {
					name     string
					from, to *Catalog
				}{"snapshot", loadSnapshotIn(t, snapA, applied), loadSnapshotIn(t, snapB, applied)})
			}
			reason, refused := migrationRefusals[p.name]
			for i, path := range paths {
				if i > 0 {
					setUp(t, applied, p.a)
				}
				plan := GenerateMigration(path.from, path.to, Diff(path.from, path.to))
				exec(t, "USE "+applied)
				var failed *MigrationOp
				var failure error
				for i, op := range plan.Ops {
					if _, err := ctr.db.ExecContext(ctr.ctx, op.SQL); err != nil {
						failed, failure = &plan.Ops[i], err
						break
					}
				}
				got := storedSchema(t, ctr, applied)
				switch {
				case refused && failed == nil && got == want:
					t.Errorf("%s: listed as refused (%s) but converges: remove the entry", path.name, reason)
				case refused && failed != nil && failed.Warning == "":
					t.Errorf("%s: TiDB refuses %q (%v) and the op carries no warning", path.name, failed.SQL, failure)
				case refused:
				case failed != nil:
					t.Errorf("%s: TiDB refuses %q: %v\nplan:\n%s", path.name, failed.SQL, failure, plan.SQL())
				case got != want:
					t.Errorf("%s: stored schema differs\nplan:\n%s\n--- got\n%s\n--- want\n%s", path.name, plan.SQL(), got, want)
				}
			}
		})
	}
}

func execCatalogIn(t *testing.T, db, script string) *Catalog {
	t.Helper()
	c := New()
	c.SetCheckConstraints(true)
	c.SetForeignKeyChecks(false)
	results, err := c.Exec("CREATE DATABASE "+db+"; USE "+db+";\n"+strings.ReplaceAll(script, diffOracleDBB+".", db+"."), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("catalog: line %d: %v", r.Line, r.Error)
		}
	}
	return c
}

func loadSnapshotIn(t *testing.T, meta *metadata.DatabaseSchemaMetadata, db string) *Catalog {
	t.Helper()
	c := New()
	report, err := c.LoadMetadataInto(context.Background(), meta, db)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(report.Degraded) + len(report.Missing) + len(report.Incomplete); n > 0 {
		t.Fatalf("snapshot %s: degraded %v, missing %v, incomplete %v", meta.GetName(), report.Degraded, report.Missing, report.Incomplete)
	}
	return c
}
