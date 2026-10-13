package review

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// finding is a finding as the tests compare it: the statement, the text
// its range anchors, the message, and the targets.
type finding struct {
	Statement int
	Text      string
	Message   string
	Targets   []int
}

func run(t *testing.T, rule review.Rule, sql string, change review.Change, targets []review.Target) []finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{rule}, Change: change}, targets)
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if len(result.Failures) != 0 {
		t.Fatalf("Failures = %v", result.Failures)
	}
	var got []finding
	for _, f := range result.Findings {
		if f.Rule == review.Syntax && rule != review.Syntax {
			t.Fatalf("syntax finding: %s", f.Message)
		}
		if f.Rule != rule {
			continue
		}
		got = append(got, finding{f.Statement, sql[f.Range.Start:f.Range.End], f.Message, f.Targets})
	}
	return got
}

type ruleCase struct {
	name    string
	sql     string
	change  review.Change
	targets []review.Target
	want    []finding
}

func runCases(t *testing.T, rule review.Rule, cases []ruleCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			targets := tt.targets
			if targets == nil {
				targets = []review.Target{{Schema: shop(), BackupDatabaseExists: true}}
			}
			got := run(t, rule, tt.sql, tt.change, targets)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// shop is a small database: t with a primary key, a unique key, a check,
// and a foreign key to p; nokey without a key; plain with only a primary
// key; and trig with a trigger.
func shop() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name:    "p",
					Columns: columns("id", "note"),
					Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}},
				},
				{
					Name:    "t",
					Columns: columns("id", "code", "p_id", "n"),
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
						{Name: "uk_code", Unique: true, Expressions: []string{"code"}},
						{Name: "idx_n", Expressions: []string{"n"}},
					},
					ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "fk_p", Columns: []string{"p_id"}, ReferencedTable: "p", ReferencedColumns: []string{"id"}}},
					CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_n", Expression: "n > 0"}},
				},
				{
					Name:    "nokey",
					Columns: columns("a"),
				},
				{
					Name:    "plain",
					Columns: columns("id", "v"),
					Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}},
				},
				{
					Name:     "trig",
					Columns:  columns("id"),
					Indexes:  []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}},
					Triggers: []*metadata.TriggerMetadata{{Name: "trig_bi"}},
				},
				{
					Name:    "nul",
					Columns: []*metadata.ColumnMetadata{{Name: "email", Nullable: true}, {Name: "code"}},
					Indexes: []*metadata.IndexMetadata{{Name: "uk_email", Unique: true, Expressions: []string{"email"}}},
				},
				{
					Name:    "fkey",
					Columns: columns("id", "code"),
					Indexes: []*metadata.IndexMetadata{{Name: "uk_f", Unique: true, Expressions: []string{"lower(`code`)"}}},
				},
			},
			Views: []*metadata.ViewMetadata{{Name: "pv"}},
		}},
	}
}

func columns(names ...string) []*metadata.ColumnMetadata {
	out := make([]*metadata.ColumnMetadata, len(names))
	for i, n := range names {
		out[i] = &metadata.ColumnMetadata{Name: n}
	}
	return out
}

func TestSyntax(t *testing.T) {
	runCases(t, review.Syntax, []ruleCase{
		{name: "valid", sql: "SELECT 1; UPDATE t SET n = 1 WHERE id = 1;"},
		{
			name: "one finding for the first error",
			sql:  "SELECT 1;\nSELEC 2;\nSELEC 3;",
			want: []finding{{1, "", `syntax error at or near "SELEC"`, []int{0}}},
		},
	})
}

func TestRequireWhere(t *testing.T) {
	runCases(t, review.RequireWhere, []ruleCase{
		{
			name: "UPDATE and DELETE without WHERE",
			sql:  "UPDATE t SET n = 1;\nDELETE FROM t;\nUPDATE t SET n = 1 WHERE id = 1;\nDELETE t, p FROM t JOIN p ON t.p_id = p.id;",
			want: []finding{
				{0, "UPDATE t SET n = 1", "UPDATE t has no WHERE clause", []int{0}},
				{1, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}},
				{3, "DELETE t, p FROM t JOIN p ON t.p_id = p.id", "DELETE FROM t, p has no WHERE clause", []int{0}},
			},
		},
		{
			name: "inside a procedure, not under EXPLAIN",
			sql:  "CREATE PROCEDURE pr() BEGIN DELETE FROM t; END;\nEXPLAIN DELETE FROM t;",
			want: []finding{{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}}},
		},
	})
}

func TestRequireIsNull(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "comparisons a predicate tests",
			sql:  "SELECT * FROM t WHERE n = NULL OR NULL <> code;\nUPDATE t SET n = 1 WHERE NOT (code != NULL);\nSELECT * FROM t JOIN p ON t.p_id = NULL;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "NULL <> code", `"NULL <> code" is never true; use IS NOT NULL`, []int{0}},
				{1, "code != NULL", `"code != NULL" is never true; use IS NOT NULL`, []int{0}},
				{2, "t.p_id = NULL", `"t.p_id = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "a value, the null-safe operator, and IS NULL",
			sql:  "SELECT n = NULL FROM t;\nSELECT * FROM t WHERE n <=> NULL;\nSELECT * FROM t WHERE n IS NULL;\nSELECT * FROM t WHERE IFNULL(n = NULL, 1);",
		},
		{
			name: "IF's condition",
			sql:  "SELECT IF(n = NULL, 1, 0) FROM t;",
			want: []finding{{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}}},
		},
		{
			name: "SHOW ... WHERE",
			sql:  "SHOW INDEX FROM t WHERE Key_name = NULL;",
			want: []finding{{0, "Key_name = NULL", `"Key_name = NULL" is never true; use IS NULL`, []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t, p;\nDROP VIEW v;\nDROP DATABASE d;\nDROP INDEX idx_n ON t;\nDROP PROCEDURE pr;\nDROP FUNCTION f;\nDROP TRIGGER tr;\nDROP EVENT e;\nDROP USER u;\nALTER TABLE t DROP COLUMN n, DROP INDEX idx_n, DROP PARTITION p0;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t, p", "drops table t, p", []int{0}},
				{1, "DROP VIEW v", "drops view v", []int{0}},
				{2, "DROP DATABASE d", "drops database d", []int{0}},
				{3, "DROP INDEX idx_n ON t", "drops index idx_n of t", []int{0}},
				{4, "DROP PROCEDURE pr", "drops procedure pr", []int{0}},
				{5, "DROP FUNCTION f", "drops function f", []int{0}},
				{6, "DROP TRIGGER tr", "drops trigger tr", []int{0}},
				{7, "DROP EVENT e", "drops event e", []int{0}},
				{8, "DROP USER u", "drops user u", []int{0}},
				{9, "DROP COLUMN n", "drops column n of t", []int{0}},
				{9, "DROP INDEX idx_n", "drops index idx_n of t", []int{0}},
				{9, "DROP PARTITION p0", "drops partition p0 of t", []int{0}},
			},
		},
		{
			name: "constraints belong to DisallowDropConstraint",
			sql:  "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CONSTRAINT ck_n, DROP CONSTRAINT uk_code;\nDROP INDEX `PRIMARY` ON t;",
		},
		{
			name: "a name with a line break",
			sql:  "DROP TABLE `a\nb`;",
			want: []finding{{0, "DROP TABLE `a\nb`", "drops table `a\\nb`", []int{0}}},
		},
		{
			name: "an account with a line break",
			sql:  "DROP USER 'a\nb'@'%';",
			want: []finding{{0, "DROP USER 'a\nb'@'%'", "drops user a\\nb@%", []int{0}}},
		},
		{
			name: "MariaDB objects and subcommands",
			sql:  "DROP SEQUENCE s;\nALTER TABLE t DROP SYSTEM VERSIONING;\nALTER TABLE t DROP PERIOD FOR SYSTEM_TIME;",
			want: []finding{
				{0, "DROP SEQUENCE s", "drops sequence s", []int{0}},
				{1, "DROP SYSTEM VERSIONING", "drops system versioning of t", []int{0}},
				{2, "DROP PERIOD FOR SYSTEM_TIME", "drops period SYSTEM_TIME of t", []int{0}},
			},
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{
			name: "TRUNCATE and TRUNCATE PARTITION",
			sql:  "TRUNCATE TABLE t;\nALTER TABLE t TRUNCATE PARTITION p0, p1;\nALTER TABLE t TRUNCATE PARTITION ALL;",
			want: []finding{
				{0, "TRUNCATE TABLE t", "truncates t", []int{0}},
				{1, "TRUNCATE PARTITION p0, p1", "truncates partition p0, p1 of t", []int{0}},
				{2, "TRUNCATE PARTITION ALL", "truncates all partitions of t", []int{0}},
			},
		},
	})
}

func TestDisallowRename(t *testing.T) {
	runCases(t, review.DisallowRename, []ruleCase{
		{
			name: "tables and columns",
			sql:  "RENAME TABLE t TO t2, p TO p2;\nALTER TABLE t RENAME TO t3;\nALTER TABLE t RENAME COLUMN n TO m;\nALTER TABLE t CHANGE COLUMN n m INT;\nALTER TABLE t CHANGE n N BIGINT;\nALTER TABLE t RENAME INDEX idx_n TO idx_m;",
			want: []finding{
				{0, "t TO t2", "renames table t to t2", []int{0}},
				{0, "p TO p2", "renames table p to p2", []int{0}},
				{1, "RENAME TO t3", "renames table t to t3", []int{0}},
				{2, "RENAME COLUMN n TO m", "renames column n of t to m", []int{0}},
				{3, "CHANGE COLUMN n m INT", "renames column n of t to m", []int{0}},
			},
		},
	})
}

func TestDisallowDropConstraint(t *testing.T) {
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{
			name: "the kinds the SQL names",
			sql:  "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CONSTRAINT ck_n, DROP CONSTRAINT uk_code;\nDROP INDEX `PRIMARY` ON p;",
			want: []finding{
				{0, "DROP PRIMARY KEY", "drops the primary key of t", []int{0}},
				{0, "DROP FOREIGN KEY fk_p", "drops foreign key fk_p of t", []int{0}},
				{0, "DROP CONSTRAINT ck_n", "drops constraint ck_n of t", []int{0}},
				{0, "DROP CONSTRAINT uk_code", "drops constraint uk_code of t", []int{0}},
				{1, "DROP INDEX `PRIMARY` ON p", "drops the primary key of p", []int{0}},
			},
		},
		{
			name: "a unique index of the synced schema, or one the change made",
			sql:  "ALTER TABLE t DROP INDEX uk_code;\nDROP INDEX idx_n ON t;\nALTER TABLE nokey ADD UNIQUE KEY uk_a (a);\nALTER TABLE nokey DROP INDEX uk_a;",
			want: []finding{
				{0, "DROP INDEX uk_code", "drops unique key uk_code of t", []int{0}},
				{3, "DROP INDEX uk_a", "drops unique key uk_a of nokey", []int{0}},
			},
		},
		{
			name: "a unique index named another way, in the database named",
			sql:  "DROP INDEX UK_CODE ON shop.t;\nDROP INDEX uk_code ON other.t;",
			want: []finding{{0, "DROP INDEX UK_CODE ON shop.t", "drops unique key UK_CODE of shop.t", []int{0}}},
		},
		{
			name: "an index the change touched before, or after a USE",
			sql:  "ALTER TABLE t DROP INDEX uk_code, ADD INDEX uk_code (code);\nALTER TABLE t DROP INDEX uk_code;\nUSE other;\nDROP INDEX uk_code ON t;",
			want: []finding{{0, "DROP INDEX uk_code", "drops unique key uk_code of t", []int{0}}},
		},
		{
			name: "a check constraint of an index's name",
			sql:  "ALTER TABLE t ADD CONSTRAINT uk_code CHECK (n > 0);\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "drops unique key uk_code of t", []int{0}}},
		},
		{
			name: "an index of another database's table",
			sql:  "CREATE UNIQUE INDEX idx_n ON other.t (n);\nDROP INDEX idx_n ON t;",
		},
		{
			name: "a table named in another case, or a routine body's index",
			sql:  "CREATE UNIQUE INDEX ux ON T (c);\nDROP INDEX ux ON t;\nCREATE PROCEDURE pr() CREATE UNIQUE INDEX uy ON nokey (a);\nDROP INDEX uy ON nokey;\nALTER EVENT ev DO CREATE UNIQUE INDEX uz ON nokey (a);\nDROP INDEX uz ON nokey;",
		},
		{
			name: "CREATE UNIQUE INDEX IF NOT EXISTS, which may leave the index as it was",
			sql:  "CREATE UNIQUE INDEX IF NOT EXISTS idx_n ON t (n);\nDROP INDEX idx_n ON t;",
		},
		{
			name: "an index of a table CREATE TABLE IF NOT EXISTS leaves",
			sql:  "CREATE TABLE IF NOT EXISTS t (id INT);\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "drops unique key uk_code of t", []int{0}}},
		},
		{
			name: "an index of a table the change creates or renames to",
			sql:  "RENAME TABLE t TO t_old, plain TO t;\nDROP INDEX uk_code ON t;\nDROP TABLE t;\nCREATE TABLE t (code INT, KEY uk_code (code));\nALTER TABLE t DROP INDEX uk_code;",
		},
	})
}

func TestRequirePrimaryKey(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "CREATE TABLE without a primary key",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT, PRIMARY KEY (id));\nCREATE TEMPORARY TABLE d (id INT);\nCREATE TABLE e LIKE t;\nCREATE TABLE f AS SELECT 1 AS id;\nCREATE TABLE IF NOT EXISTS g (id INT);\nCREATE TABLE h (PRIMARY KEY (id)) AS SELECT 1 AS id;",
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}},
				{5, "CREATE TABLE f AS SELECT 1 AS id", "creates table f without a primary key", []int{0}},
			},
		},
		{
			name: "CREATE OR REPLACE TABLE",
			sql:  "CREATE OR REPLACE TABLE plain (id INT);",
			want: []finding{{0, "CREATE OR REPLACE TABLE plain (id INT)", "creates table plain without a primary key", []int{0}}},
		},
		{
			name: "settled by a later key, drop, or rename",
			sql:  "CREATE TABLE a (id INT);\nALTER TABLE a ADD PRIMARY KEY (id);\nCREATE TABLE b (id INT);\nDROP TABLE b;\nCREATE TABLE c (id INT);\nRENAME TABLE c TO c2;\nCREATE DATABASE d;\nCREATE TABLE d.e (id INT);\nDROP DATABASE d;",
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE p DROP PRIMARY KEY, ADD PRIMARY KEY (id, note);\nALTER TABLE trig DROP PRIMARY KEY;\nALTER TABLE trig ADD PRIMARY KEY (id);",
			want: []finding{{0, "DROP PRIMARY KEY", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
	})
}

func TestOnlineMigration(t *testing.T) {
	ghost := review.Change{OnlineMigration: true}
	runCases(t, review.OnlineMigration, []ruleCase{
		{name: "one ALTER TABLE of a table gh-ost can copy", sql: "-- gh-ost = {\"chunk-size\":\"500\"}\nALTER TABLE plain ADD COLUMN m INT;", change: ghost},
		{name: "not asked for", sql: "UPDATE t SET n = 1; UPDATE t SET n = 2;"},
		{
			name:   "more than one statement",
			sql:    "ALTER TABLE t ADD COLUMN m INT;\nALTER TABLE p ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{-1, "", "online migration runs the change as one ALTER TABLE statement, and the change has 2 statements", []int{0}}},
		},
		{
			name:   "not an ALTER TABLE",
			sql:    "CREATE INDEX i ON t (n);",
			change: ghost,
			want:   []finding{{0, "CREATE INDEX i ON t (n)", "online migration runs only an ALTER TABLE statement", []int{0}}},
		},
		{
			name:   "flags Bytebase cannot read",
			sql:    "-- gh-ost = {\"chunk-size\":\"many\",\"no-such-flag\":\"1\"}\nALTER TABLE plain ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{-1, "", `online migration cannot read gh-ost flag "chunk-size" as int: "many"`, []int{0}}},
		},
		{
			name:   "an unknown flag",
			sql:    "-- gh-ost = {\"include-triggers\":\"true\"}\nALTER TABLE trig ADD COLUMN m INT;",
			change: ghost,
			want: []finding{
				{-1, "", `online migration does not support gh-ost flag "include-triggers"`, []int{0}},
				{0, "ALTER TABLE trig ADD COLUMN m INT", "online migration does not support trig, which has a trigger", []int{0}},
			},
		},
		{
			name:   "a rename of the table",
			sql:    "ALTER TABLE plain ADD COLUMN m INT, RENAME TO plain2;",
			change: ghost,
			want:   []finding{{0, "RENAME TO plain2", "gh-ost does not rename the table it migrates", []int{0}}},
		},
		{
			name:   "a column renamed by CHANGE",
			sql:    "ALTER TABLE plain CHANGE COLUMN v w INT, CHANGE id id BIGINT;",
			change: ghost,
			want:   []finding{{0, "CHANGE COLUMN v w INT", "gh-ost renames column v to w only with approve-renamed-columns, which Bytebase does not pass", []int{0}}},
		},
		{
			name:   "flag values gh-ost cannot read",
			sql:    "-- gh-ost = {\"max-load\":\"Threads_running\"}\nALTER TABLE plain ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{-1, "", `online migration cannot read gh-ost flag "max-load" as load map: "Threads_running"`, []int{0}}},
		},
		{
			name:   "flag values gh-ost reads",
			sql:    "-- gh-ost = {\"max-load\":\"Threads_running=25,Threads_connected=100\",\"throttle-control-replicas\":\"r1:3306,r2,[::1]:3307\"}\nALTER TABLE plain ADD COLUMN m INT;",
			change: ghost,
		},
		{
			name:   "only a nullable unique key",
			sql:    "ALTER TABLE nul ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE nul ADD COLUMN m INT", "online migration needs a primary key or a unique key of NOT NULL columns on nul, which has none", []int{0}}},
		},
		{
			name:   "the table of the task's database, whatever the ALTER names",
			sql:    "ALTER TABLE other.nokey ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE other.nokey ADD COLUMN m INT", "online migration needs a primary or unique key on other.nokey, which has none", []int{0}}},
		},
		{
			name:   "a view",
			sql:    "ALTER TABLE pv ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE pv ADD COLUMN m INT", "online migration needs a table, and pv is a view", []int{0}}},
		},
		{
			name:   "a table whose only unique key is of expressions",
			sql:    "ALTER TABLE fkey ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE fkey ADD COLUMN m INT", "online migration needs a primary or unique key on fkey, which has none", []int{0}}},
		},
		{
			name:   "tables gh-ost refuses",
			sql:    "ALTER TABLE nokey ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE nokey ADD COLUMN m INT", "online migration needs a primary or unique key on nokey, which has none", []int{0}}},
		},
		{
			name:   "a foreign key of its own",
			sql:    "ALTER TABLE t ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE t ADD COLUMN m INT", "online migration does not support t, which has a foreign key", []int{0}}},
		},
		{
			name:   "a foreign key to it",
			sql:    "ALTER TABLE p ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE p ADD COLUMN m INT", "online migration does not support p, which a foreign key references", []int{0}}},
		},
		{
			name:    "no bbdataarchive",
			sql:     "ALTER TABLE trig ENGINE = InnoDB;",
			change:  ghost,
			targets: []review.Target{{Schema: shop(), BackupDatabaseExists: false}},
			want: []finding{
				{-1, "", "online migration needs database bbdataarchive, which does not exist", []int{0}},
				{0, "ALTER TABLE trig ENGINE = InnoDB", "online migration does not support trig, which has a trigger", []int{0}},
			},
		},
	})
}

func TestPriorBackup(t *testing.T) {
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	runCases(t, review.PriorBackup, []ruleCase{
		{name: "a plain backup", sql: "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM p WHERE id = 2;", change: backup},
		{name: "not asked for", sql: "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;"},
		{
			name:    "no bbdataarchive, with or without DML",
			sql:     "ALTER TABLE t ADD COLUMN m INT;",
			change:  backup,
			targets: []review.Target{{Schema: shop()}},
			want:    []finding{{-1, "", "prior backup needs database bbdataarchive, which does not exist", []int{0}}},
		},
		{
			name:   "over the size limit",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;",
			change: review.Change{PriorBackup: true, MaxBackupSize: 10},
			want:   []finding{{-1, "", "the change is 32 bytes, over the 10-byte limit of prior backup", []int{0}}},
		},
		{
			name:   "UPDATE and DELETE on one table, a name without a database the target's",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;\nDELETE FROM shop.p WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nDELETE FROM other.plain WHERE id = 1;\nUPDATE plain SET v = 1 WHERE id = 1;",
			change: backup,
			want: []finding{
				{1, "DELETE FROM t WHERE id = 2", "prior backup cannot back up both UPDATE and DELETE on t", []int{0}},
				{3, "UPDATE p SET note = '' WHERE id = 1", "prior backup cannot back up both UPDATE and DELETE on p", []int{0}},
			},
		},
		{
			name:   "a table the change creates",
			sql:    "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nDELETE t FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1;\nDELETE FROM nope WHERE a = 1;",
			change: backup,
			want: []finding{
				{1, "UPDATE n2 SET id = 2 WHERE id = 1", "prior backup runs before the change, when n2 does not exist yet", []int{0}},
				{2, "DELETE t FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1", "prior backup runs before the change, when n2 does not exist yet", []int{0}},
			},
		},
		{
			name:   "a column the change adds or renames to",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;\nUPDATE t SET n = 2 WHERE t.gone = 2 ORDER BY id;\nALTER TABLE p RENAME COLUMN note TO memo;\nUPDATE t JOIN p ON p.memo = t.code SET t.n = 1 WHERE t.id = 3;",
			change: backup,
			want: []finding{
				{1, "UPDATE t SET n = 1 WHERE gone = 1", "prior backup runs before the change, when column gone of t does not exist yet", []int{0}},
				{2, "UPDATE t SET n = 2 WHERE t.gone = 2 ORDER BY id", "prior backup runs before the change, when column gone of t does not exist yet", []int{0}},
				{4, "UPDATE t JOIN p ON p.memo = t.code SET t.n = 1 WHERE t.id = 3", "prior backup runs before the change, when column memo of p does not exist yet", []int{0}},
			},
		},
		{
			name:   "UPDATE and DELETE on one table, a DELETE of several tables",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE t FROM t JOIN p ON p.id = t.p_id WHERE t.id = 2;\nUPDATE plain SET v = 1 WHERE id = 1;\nDELETE x FROM plain AS x JOIN p ON p.id = x.id WHERE x.id = 2;",
			change: backup,
			want: []finding{
				{1, "DELETE t FROM t JOIN p ON p.id = t.p_id WHERE t.id = 2", "prior backup cannot back up both UPDATE and DELETE on t", []int{0}},
				{3, "DELETE x FROM plain AS x JOIN p ON p.id = x.id WHERE x.id = 2", "prior backup cannot back up both UPDATE and DELETE on plain", []int{0}},
			},
		},
		{
			name:   "a column USING names",
			sql:    "ALTER TABLE t ADD COLUMN note INT;\nUPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1;\nUPDATE t JOIN p USING (id) SET t.n = 1 WHERE t.id = 2;\nUPDATE t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) SET t.n = 1 WHERE t.id = 3;",
			change: backup,
			want:   []finding{{1, "UPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1", "prior backup runs before the change, when column note of t does not exist yet", []int{0}}},
		},
		{
			name:   "a view as the target",
			sql:    "UPDATE pv SET note = 'z' WHERE id = 1;\nDELETE FROM other.pv WHERE id = 2;",
			change: backup,
			want:   []finding{{0, "UPDATE pv SET note = 'z' WHERE id = 1", "prior backup copies tables with CREATE TABLE ... LIKE, which refuses view pv", []int{0}}},
		},
		{
			name:   "a DELETE target the FROM list does not name the same way",
			sql:    "DELETE T FROM t WHERE id = 1;\nDELETE x FROM t AS x WHERE x.id = 2;",
			change: backup,
			want:   []finding{{0, "DELETE T FROM t WHERE id = 1", "prior backup does not find DELETE target T in the FROM list, written the same way", []int{0}}},
		},
		{
			name:   "statements before a USE",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;\nUPDATE n3 SET id = 1 WHERE id = 2;\nUSE other;\nCREATE TABLE n3 (id INT PRIMARY KEY);\nUPDATE t SET n = 2 WHERE gone = 2;",
			change: backup,
			want:   []finding{{1, "UPDATE t SET n = 1 WHERE gone = 1", "prior backup runs before the change, when column gone of t does not exist yet", []int{0}}},
		},
		{
			name:   "what the rule cannot tell",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE id IN (SELECT gone FROM t);\nUPDATE t SET gone = 1 WHERE id = 1;\nWITH x AS (SELECT 1 AS gone) UPDATE t SET n = 3 WHERE id IN (SELECT gone FROM x);",
			change: backup,
		},
	})
}
