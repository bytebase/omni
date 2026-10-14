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
		{
			name: "a feature comment for a feature TiDB does not support",
			sql:  "/*T![unsupported_feature_xyz] WEIRD STUFF */;\nSELECT 1;",
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
		{
			name: "EXPLAIN ANALYZE runs the statement, a BATCH dry run does not",
			sql:  "EXPLAIN ANALYZE DELETE FROM t;\nEXPLAIN DELETE FROM t;\nBATCH ON id LIMIT 2 DELETE FROM t;\nBATCH ON id LIMIT 2 DRY RUN DELETE FROM t;",
			want: []finding{
				{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}},
				{2, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}},
			},
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
		{
			name: "a searched CASE statement",
			sql:  "CREATE PROCEDURE pr() BEGIN CASE WHEN @n = NULL THEN SELECT 1; END CASE; END;",
			want: []finding{{0, "@n = NULL", `"@n = NULL" is never true; use IS NULL`, []int{0}}},
		},
		{
			name: "quoted text kept as written",
			sql:  "SELECT * FROM t WHERE `a  b` = NULL;",
			want: []finding{{0, "`a  b` = NULL", "\"`a  b` = NULL\" is never true; use IS NULL", []int{0}}},
		},
		{
			name: "a backslash-escaped quote in a string",
			sql:  "SELECT * FROM t WHERE 'a\\'  b' = NULL;",
			want: []finding{{0, "'a\\'  b' = NULL", "\"'a\\'  b' = NULL\" is never true; use IS NULL", []int{0}}},
		},
		{
			name: "HANDLER ... READ WHERE",
			sql:  "HANDLER t READ NEXT WHERE id = NULL;",
			want: []finding{{0, "id = NULL", "\"id = NULL\" is never true; use IS NULL", []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t, p;\nDROP VIEW v;\nDROP DATABASE d;\nDROP INDEX idx_n ON t;\nDROP PROCEDURE pr;\nDROP PLACEMENT POLICY pp;\nDROP RESOURCE GROUP rg;\nDROP ROLE r1;\nDROP USER u;\nALTER TABLE t DROP COLUMN n, DROP INDEX idx_n, DROP PARTITION p0;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t, p", "drops table t, p", []int{0}},
				{1, "DROP VIEW v", "drops view v", []int{0}},
				{2, "DROP DATABASE d", "drops database d", []int{0}},
				{3, "DROP INDEX idx_n ON t", "drops index idx_n of t", []int{0}},
				{4, "DROP PROCEDURE pr", "drops procedure pr", []int{0}},
				{5, "DROP PLACEMENT POLICY pp", "drops placement policy pp", []int{0}},
				{6, "DROP RESOURCE GROUP rg", "drops resource group rg", []int{0}},
				{7, "DROP ROLE r1", "drops role r1", []int{0}},
				{8, "DROP USER u", "drops user u", []int{0}},
				{9, "DROP COLUMN n", "drops column n of t", []int{0}},
				{9, "DROP INDEX idx_n", "drops index idx_n of t", []int{0}},
				{9, "DROP PARTITION p0", "drops partition p0 of t", []int{0}},
			},
		},
		{
			name: "constraints belong to DisallowDropConstraint",
			sql:  "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n, DROP CONSTRAINT uk_code;\nDROP INDEX `PRIMARY` ON t;",
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
			name: "a statement in a TiDB executable comment",
			sql:  "/*T! DROP TABLE t2 */;",
			want: []finding{{0, "/*T! DROP TABLE t2 */", "drops table t2", []int{0}}},
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
			sql:  "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n, DROP CONSTRAINT uk_code;\nDROP INDEX `PRIMARY` ON p;",
			want: []finding{
				{0, "DROP PRIMARY KEY", "drops the primary key of t", []int{0}},
				{0, "DROP FOREIGN KEY fk_p", "drops foreign key fk_p of t", []int{0}},
				{0, "DROP CHECK ck_n", "drops check constraint ck_n of t", []int{0}},
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
			name: "a table of another database renamed to a synced table's name",
			sql:  "ALTER TABLE other.src RENAME TO t;\nDROP INDEX uk_code ON shop.t;",
			want: []finding{{1, "DROP INDEX uk_code ON shop.t", "drops unique key uk_code of shop.t", []int{0}}},
		},
		{
			name: "an index of a table dropped and created again",
			sql:  "CREATE UNIQUE INDEX ux ON nokey (a);\nDROP TABLE nokey;\nCREATE TABLE nokey (a INT, KEY ux (a));\nDROP INDEX ux ON nokey;",
		},
		{
			name: "a temporary table dropped, which leaves the table of the name as it was",
			sql:  "DROP TEMPORARY TABLE IF EXISTS t;\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "drops unique key uk_code of t", []int{0}}},
		},
		{
			name: "a unique index the change renames",
			sql:  "CREATE UNIQUE INDEX ux ON nokey (a);\nALTER TABLE nokey RENAME INDEX ux TO uy;\nDROP INDEX uy ON nokey;",
			want: []finding{{2, "DROP INDEX uy ON nokey", "drops unique key uy of nokey", []int{0}}},
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
			sql:  "CREATE UNIQUE INDEX ux ON T (c);\nDROP INDEX ux ON t;\nCREATE PROCEDURE pr() CREATE UNIQUE INDEX uy ON nokey (a);\nDROP INDEX uy ON nokey;",
		},
		{
			name: "CREATE UNIQUE INDEX IF NOT EXISTS, which may leave the index as it was",
			sql:  "CREATE UNIQUE INDEX IF NOT EXISTS idx_n ON t (n);\nDROP INDEX idx_n ON t;",
		},
		{
			name: "a key CREATE TABLE IF NOT EXISTS declares, or one touched before a USE",
			sql:  "CREATE TABLE IF NOT EXISTS t (n INT, UNIQUE KEY idx_n (n));\nDROP INDEX idx_n ON t;\nCREATE UNIQUE INDEX ux ON nokey (a);\nUSE other;\nDROP INDEX ux ON nokey;",
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
			want: []finding{{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}}},
		},
		{
			name: "a constraint named PRIMARY",
			sql:  "ALTER TABLE t DROP FOREIGN KEY `PRIMARY`;",
		},
		{
			name: "not settled by IF NOT EXISTS or a temporary table of its name",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE IF NOT EXISTS a (id INT PRIMARY KEY);\nCREATE TABLE b (id INT);\nCREATE TEMPORARY TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT);\nDROP TEMPORARY TABLE IF EXISTS c;",
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}},
				{2, "CREATE TABLE b (id INT)", "creates table b without a primary key", []int{0}},
				{4, "CREATE TABLE c (id INT)", "creates table c without a primary key", []int{0}},
			},
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
	runCases(t, review.OnlineMigration, []ruleCase{
		{name: "not asked for", sql: "ALTER TABLE t ADD COLUMN m INT;"},
		{
			name:   "asked for",
			sql:    "-- gh-ost = {}\nALTER TABLE plain ADD COLUMN m INT;",
			change: review.Change{OnlineMigration: true},
			want:   []finding{{-1, "", "the change requests online migration, which runs gh-ost, and TiDB has no binary log for gh-ost to read", []int{0}}},
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
			name:   "UPDATE and DELETE on one table, which TiDB's backup takes",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			change: backup,
		},
		{
			name:   "statements the backup refuses",
			sql:    "BATCH ON id LIMIT 2 DELETE FROM t WHERE n > 0;\nREPLACE INTO p VALUES (3, 'c');\nINSERT INTO p VALUES (1, 'a') ON DUPLICATE KEY UPDATE note = 'b';\nEXPLAIN ANALYZE UPDATE t SET n = 1 WHERE id = 1;\nEXPLAIN UPDATE t SET n = 1 WHERE id = 1;\nINSERT INTO p VALUES (4, 'd');",
			change: backup,
			want: []finding{
				{0, "BATCH ON id LIMIT 2 DELETE FROM t WHERE n > 0", "prior backup does not support TiDB BATCH (non-transactional DML) statements", []int{0}},
				{1, "REPLACE INTO p VALUES (3, 'c')", "prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements", []int{0}},
				{2, "INSERT INTO p VALUES (1, 'a') ON DUPLICATE KEY UPDATE note = 'b'", "prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements", []int{0}},
				{3, "EXPLAIN ANALYZE UPDATE t SET n = 1 WHERE id = 1", "prior backup does not support EXPLAIN ANALYZE of an UPDATE, DELETE, or BATCH statement", []int{0}},
			},
		},
		{
			name:   "a table of another database",
			sql:    "UPDATE other.t SET n = 1 WHERE id = 1;\nDELETE FROM SHOP.p WHERE id = 2;",
			change: backup,
			want:   []finding{{0, "UPDATE other.t SET n = 1 WHERE id = 1", "prior backup does not support an UPDATE or DELETE of other.t, a table of another database", []int{0}}},
		},
		{
			name:   "a table of another database deleted from twice",
			sql:    "DELETE a, b FROM other.t AS a JOIN other.t AS b ON a.id = b.id WHERE a.id = 1;",
			change: backup,
			want:   []finding{{0, "DELETE a, b FROM other.t AS a JOIN other.t AS b ON a.id = b.id WHERE a.id = 1", "prior backup does not support an UPDATE or DELETE of other.t, a table of another database", []int{0}}},
		},
		{
			name:   "more than five statements over two tables",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nDELETE FROM shop.t WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE p SET note = '' WHERE id = 1;",
			change: backup,
			want:   []finding{{5, "UPDATE p SET note = '' WHERE id = 1", "prior backup cannot back up more than 5 statements over more than one table", []int{0}}},
		},
		{
			name:   "more than five statements on one table, one with WITH",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nWITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x);\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE t SET n = 6 WHERE id = 6;",
			change: backup,
			want:   []finding{{2, "WITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x)", "prior backup does not support a WITH clause when more than 5 statements change one table", []int{0}}},
		},
		{
			name:   "past the limit with a WITH clause on a statement not certain",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nUPDATE t SET n = 3 WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nWITH x AS (SELECT 1 AS id) UPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id IN (SELECT id FROM x);",
			change: backup,
			want:   []finding{{-1, "", "prior backup does not support a WITH clause when it backs up more than 5 statements", []int{0}}},
		},
		{
			name:   "five statements over five tables",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nUPDATE plain SET v = 1 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;\nUPDATE trig SET id = 2 WHERE id = 1;",
			change: backup,
		},
		{
			name:   "six certain statements over two tables, with one not certain",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nUPDATE t SET n = 3 WHERE id = 3;\nUPDATE p SET note = '' WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 2;\nUPDATE p SET note = '' WHERE id = 3;\nUPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id = 4;",
			change: backup,
			want:   []finding{{-1, "", "prior backup cannot back up more than 5 statements over more than one table", []int{0}}},
		},
		{
			name:   "five statements over five tables and a multi-table UPDATE, past the limit",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nUPDATE plain SET v = 1 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;\nUPDATE trig SET id = 2 WHERE id = 1;\nUPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id = 2;",
			change: backup,
			want:   []finding{{-1, "", "prior backup cannot back up more than 5 statements over more than one table", []int{0}}},
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
			name:   "a column USING names",
			sql:    "ALTER TABLE t ADD COLUMN note INT;\nUPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1;\nUPDATE t JOIN p USING (id) SET t.n = 1 WHERE t.id = 2;\nUPDATE t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) SET t.n = 1 WHERE t.id = 3;",
			change: backup,
			want:   []finding{{1, "UPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1", "prior backup runs before the change, when column note of t does not exist yet", []int{0}}},
		},
		{
			name:   "a view as the target",
			sql:    "UPDATE pv SET note = 'z' WHERE id = 1;\nDELETE FROM other.pv WHERE id = 2;",
			change: backup,
			want: []finding{
				{0, "UPDATE pv SET note = 'z' WHERE id = 1", "prior backup copies tables with CREATE TABLE ... LIKE, which refuses view pv", []int{0}},
				{1, "DELETE FROM other.pv WHERE id = 2", "prior backup does not support an UPDATE or DELETE of other.pv, a table of another database", []int{0}},
			},
		},
		{
			name:   "a DELETE target the FROM list does not name",
			sql:    "DELETE T FROM t WHERE id = 1;\nDELETE y FROM t AS x WHERE x.id = 2;",
			change: backup,
			want:   []finding{{1, "DELETE y FROM t AS x WHERE x.id = 2", "prior backup does not find DELETE target y in the FROM list", []int{0}}},
		},
		{
			name:   "a DELETE after a WITH clause, whose targets may be its queries",
			sql:    "WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2;\n/*!80000 WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2 */;",
			change: backup,
		},
		{
			name:   "a view OR REPLACE makes, or a table IF NOT EXISTS makes",
			sql:    "CREATE OR REPLACE VIEW nv AS SELECT 1 AS id;\nUPDATE nv SET id = 2 WHERE id = 1;\nCREATE TABLE IF NOT EXISTS t (gone INT);\nUPDATE t SET n = 1 WHERE gone = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE nv SET id = 2 WHERE id = 1", "prior backup runs before the change, when nv does not exist yet", []int{0}}},
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
