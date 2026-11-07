package review

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/parser"
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
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
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
				{0, "UPDATE t SET n = 1", "`UPDATE` of `t` has no `WHERE` clause.", []int{0}},
				{1, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}},
				{3, "DELETE t, p FROM t JOIN p ON t.p_id = p.id", "`DELETE` from `t` and `p` has no `WHERE` clause.", []int{0}},
			},
		},
		{
			name: "inside a procedure, not under EXPLAIN",
			sql:  "CREATE PROCEDURE pr() BEGIN DELETE FROM t; END;\nEXPLAIN DELETE FROM t;",
			want: []finding{{0, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}}},
		},
		{
			name: "EXPLAIN ANALYZE runs the statement, a BATCH dry run does not",
			sql:  "EXPLAIN ANALYZE DELETE FROM t;\nEXPLAIN DELETE FROM t;\nBATCH ON id LIMIT 2 DELETE FROM t;\nBATCH ON id LIMIT 2 DRY RUN DELETE FROM t;",
			want: []finding{
				{0, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}},
				{2, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}},
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
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{0, "NULL <> code", "`NULL <> code` is never true.\n\nUse `IS NOT NULL` instead.", []int{0}},
				{1, "code != NULL", "`code != NULL` is never true.\n\nUse `IS NOT NULL` instead.", []int{0}},
				{2, "t.p_id = NULL", "`t.p_id = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			name: "a value, the null-safe operator, and IS NULL",
			sql:  "SELECT n = NULL FROM t;\nSELECT * FROM t WHERE n <=> NULL;\nSELECT * FROM t WHERE n IS NULL;\nSELECT * FROM t WHERE IFNULL(n = NULL, 1);",
		},
		{
			name: "IF's condition",
			sql:  "SELECT IF(n = NULL, 1, 0) FROM t;",
			want: []finding{{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "SHOW ... WHERE",
			sql:  "SHOW INDEX FROM t WHERE Key_name = NULL;",
			want: []finding{{0, "Key_name = NULL", "`Key_name = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			// TiDB's own SHOW statements carry a WHERE in the same field, so
			// the rule reaches them without knowing they exist. A new SHOW
			// variant whose filter landed anywhere else would be a silent
			// gap here.
			name: "TiDB's own SHOW ... WHERE",
			sql:  "SHOW TABLE t1 REGIONS WHERE a = NULL;\nSHOW CONFIG WHERE type = NULL;\nSHOW DISTRIBUTION JOBS WHERE id = NULL;\nSHOW PLACEMENT WHERE Target = NULL;",
			want: []finding{
				{0, "a = NULL", "`a = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{1, "type = NULL", "`type = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{2, "id = NULL", "`id = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{3, "Target = NULL", "`Target = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			// The resource control statements carry expressions of their
			// own, and the walker reaches them: a subquery inside a QUERY
			// WATCH pattern is a condition position like any other. The
			// timestamp CALIBRATE RESOURCE takes is not a condition, so
			// `START_TIME a = NULL` is no finding -- it is a value, not a
			// test -- and a rule that reported it there would be wrong.
			name: "the expressions the resource control statements carry",
			sql:  "QUERY WATCH ADD SQL DIGEST (SELECT d FROM t WHERE a = NULL);\nCALIBRATE RESOURCE START_TIME a = NULL;",
			want: []finding{
				{0, "a = NULL", "`a = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			name: "a searched CASE statement",
			sql:  "CREATE PROCEDURE pr() BEGIN CASE WHEN @n = NULL THEN SELECT 1; END CASE; END;",
			want: []finding{{0, "@n = NULL", "`@n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "quoted text kept as written",
			sql:  "SELECT * FROM t WHERE `a  b` = NULL;",
			want: []finding{{0, "`a  b` = NULL", "`` `a  b` = NULL `` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "a backslash-escaped quote in a string",
			sql:  "SELECT * FROM t WHERE 'a\\'  b' = NULL;",
			want: []finding{{0, "'a\\'  b' = NULL", "`'a\\'  b' = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "HANDLER ... READ WHERE",
			sql:  "HANDLER t READ NEXT WHERE id = NULL;",
			want: []finding{{0, "id = NULL", "`id = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t, p;\nDROP VIEW v;\nDROP DATABASE d;\nDROP INDEX idx_n ON t;\nDROP PROCEDURE pr;\nDROP PLACEMENT POLICY pp;\nDROP RESOURCE GROUP rg;\nDROP ROLE r1;\nDROP USER u;\nALTER TABLE t DROP COLUMN n, DROP INDEX idx_n, DROP PARTITION p0;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t, p", "Drops tables `t` and `p`.", []int{0}},
				{1, "DROP VIEW v", "Drops view `v`.", []int{0}},
				{2, "DROP DATABASE d", "Drops database `d`.", []int{0}},
				{3, "DROP INDEX idx_n ON t", "Drops index `idx_n` of `t`.", []int{0}},
				{4, "DROP PROCEDURE pr", "Drops procedure `pr`.", []int{0}},
				{5, "DROP PLACEMENT POLICY pp", "Drops placement policy `pp`.", []int{0}},
				{6, "DROP RESOURCE GROUP rg", "Drops resource group `rg`.", []int{0}},
				{7, "DROP ROLE r1", "Drops role `r1`.", []int{0}},
				{8, "DROP USER u", "Drops user `u`.", []int{0}},
				{9, "DROP COLUMN n", "Drops column `n` of `t`.", []int{0}},
				{9, "DROP INDEX idx_n", "Drops index `idx_n` of `t`.", []int{0}},
				{9, "DROP PARTITION p0", "Drops partition `p0` of `t`.", []int{0}},
			},
		},
		{
			// The rule reads DropResourceGroupStmt, whose shape this change
			// replaced: TiDB's DROP takes IF EXISTS and DEFAULT as a name
			// where MySQL's took FORCE.
			name: "a resource group with IF EXISTS",
			sql:  "DROP RESOURCE GROUP IF EXISTS rg;\nDROP RESOURCE GROUP DEFAULT;",
			want: []finding{
				{0, "DROP RESOURCE GROUP IF EXISTS rg", "Drops resource group `rg`.", []int{0}},
				{1, "DROP RESOURCE GROUP DEFAULT", "Drops resource group `` `DEFAULT` ``.", []int{0}},
			},
		},
		{
			name: "sequences",
			sql:  "DROP SEQUENCE s;\nDROP SEQUENCE IF EXISTS s1, shop.s2;",
			want: []finding{
				{0, "DROP SEQUENCE s", "Drops sequence `s`.", []int{0}},
				{1, "DROP SEQUENCE IF EXISTS s1, shop.s2", "Drops sequences `s1` and `shop.s2`.", []int{0}},
			},
		},
		{
			name: "constraints belong to DisallowDropConstraint",
			sql:  "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n, DROP CONSTRAINT uk_code;\nDROP INDEX `PRIMARY` ON t;",
		},
		{
			name: "a name with a line break",
			sql:  "DROP TABLE `a\nb`;",
			want: []finding{{0, "DROP TABLE `a\nb`", "Drops table `` `a\\nb` ``.", []int{0}}},
		},
		{
			name: "an account with a line break",
			sql:  "DROP USER 'a\nb'@'%';",
			want: []finding{{0, "DROP USER 'a\nb'@'%'", "Drops user `a\\nb@%`.", []int{0}}},
		},
		{
			name: "a statement in a TiDB executable comment",
			sql:  "/*T! DROP TABLE t2 */;",
			want: []finding{{0, "DROP TABLE t2", "Drops table `t2`.", []int{0}}},
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{
			name: "TRUNCATE and TRUNCATE PARTITION",
			sql:  "TRUNCATE TABLE t;\nALTER TABLE t TRUNCATE PARTITION p0, p1;\nALTER TABLE t TRUNCATE PARTITION ALL;",
			want: []finding{
				{0, "TRUNCATE TABLE t", "Truncates `t`.", []int{0}},
				{1, "TRUNCATE PARTITION p0, p1", "Truncates partitions `p0` and `p1` of `t`.", []int{0}},
				{2, "TRUNCATE PARTITION ALL", "Truncates all partitions of `t`.", []int{0}},
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
				{0, "t TO t2", "Renames table `t` to `t2`.", []int{0}},
				{0, "p TO p2", "Renames table `p` to `p2`.", []int{0}},
				{1, "RENAME TO t3", "Renames table `t` to `t3`.", []int{0}},
				{2, "RENAME COLUMN n TO m", "Renames column `n` of `t` to `m`.", []int{0}},
				{3, "CHANGE COLUMN n m INT", "Renames column `n` of `t` to `m`.", []int{0}},
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
				{0, "DROP PRIMARY KEY", "Drops the primary key of `t`.", []int{0}},
				{0, "DROP FOREIGN KEY fk_p", "Drops foreign key `fk_p` of `t`.", []int{0}},
				{0, "DROP CHECK ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}},
				{0, "DROP CONSTRAINT uk_code", "Drops constraint `uk_code` of `t`.", []int{0}},
				{1, "DROP INDEX `PRIMARY` ON p", "Drops the primary key of `p`.", []int{0}},
			},
		},
		{
			name: "a unique index of the synced schema, or one the change made",
			sql:  "ALTER TABLE t DROP INDEX uk_code;\nDROP INDEX idx_n ON t;\nALTER TABLE nokey ADD UNIQUE KEY uk_a (a);\nALTER TABLE nokey DROP INDEX uk_a;",
			want: []finding{
				{0, "DROP INDEX uk_code", "Drops unique constraint `uk_code` of `t`.", []int{0}},
				{3, "DROP INDEX uk_a", "Drops unique constraint `uk_a` of `nokey`.", []int{0}},
			},
		},
		{
			name: "a unique index named another way, in the database named",
			sql:  "DROP INDEX UK_CODE ON shop.t;\nDROP INDEX uk_code ON other.t;",
			want: []finding{{0, "DROP INDEX UK_CODE ON shop.t", "Drops unique constraint `UK_CODE` of `shop.t`.", []int{0}}},
		},
		{
			name: "an index the change touched before, or after a USE",
			sql:  "ALTER TABLE t DROP INDEX uk_code, ADD INDEX uk_code (code);\nALTER TABLE t DROP INDEX uk_code;\nUSE other;\nDROP INDEX uk_code ON t;",
			want: []finding{{0, "DROP INDEX uk_code", "Drops unique constraint `uk_code` of `t`.", []int{0}}},
		},
		{
			name: "a table of another database renamed to a synced table's name",
			sql:  "ALTER TABLE other.src RENAME TO t;\nDROP INDEX uk_code ON shop.t;",
			want: []finding{{1, "DROP INDEX uk_code ON shop.t", "Drops unique constraint `uk_code` of `shop.t`.", []int{0}}},
		},
		{
			name: "an index of a table dropped and created again",
			sql:  "CREATE UNIQUE INDEX ux ON nokey (a);\nDROP TABLE nokey;\nCREATE TABLE nokey (a INT, KEY ux (a));\nDROP INDEX ux ON nokey;",
		},
		{
			name: "a routine of another database, whose body resolves a bare table there",
			sql:  "CREATE PROCEDURE other.pr() DROP INDEX uk_code ON t;",
		},
		{
			name: "a unique index made in one branch of a routine body and dropped in another",
			sql:  "DELIMITER //\nCREATE PROCEDURE pr() BEGIN IF 1 THEN CREATE UNIQUE INDEX ux ON nokey (a); ELSE DROP INDEX ux ON nokey; END IF; END//\nDELIMITER ;",
		},
		{
			name: "a routine body's own DDL, which keeps its later drops from being judged and touches nothing after it",
			sql:  "DELIMITER //\nCREATE PROCEDURE pr() BEGIN CREATE TEMPORARY TABLE t (id INT, KEY uk_code (id)); DROP INDEX uk_code ON t; END//\nDELIMITER ;\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "Drops unique constraint `uk_code` of `t`.", []int{0}}},
		},
		{
			name: "an index named with its database, across a USE",
			sql:  "DROP INDEX uk_code ON shop.t;\nCREATE INDEX uk_code ON shop.t (code);\nUSE other;\nDROP INDEX uk_code ON shop.t;",
			want: []finding{{0, "DROP INDEX uk_code ON shop.t", "Drops unique constraint `uk_code` of `shop.t`.", []int{0}}},
		},
		{
			name: "a temporary table dropped, which leaves the table of the name as it was",
			sql:  "DROP TEMPORARY TABLE IF EXISTS t;\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "Drops unique constraint `uk_code` of `t`.", []int{0}}},
		},
		{
			name: "a unique index the change renames",
			sql:  "CREATE UNIQUE INDEX ux ON nokey (a);\nALTER TABLE nokey RENAME INDEX ux TO uy;\nDROP INDEX uy ON nokey;",
			want: []finding{{2, "DROP INDEX uy ON nokey", "Drops unique constraint `uy` of `nokey`.", []int{0}}},
		},
		{
			name: "a check constraint of an index's name",
			sql:  "ALTER TABLE t ADD CONSTRAINT uk_code CHECK (n > 0);\nDROP INDEX uk_code ON t;",
			want: []finding{{1, "DROP INDEX uk_code ON t", "Drops unique constraint `uk_code` of `t`.", []int{0}}},
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
			want: []finding{{1, "DROP INDEX uk_code ON t", "Drops unique constraint `uk_code` of `t`.", []int{0}}},
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
			want: []finding{{0, "CREATE TABLE a (id INT)", "Creates table `a` without a primary key.", []int{0}}},
		},
		{
			name: "a constraint named PRIMARY",
			sql:  "ALTER TABLE t DROP FOREIGN KEY `PRIMARY`;",
		},
		{
			name: "not settled by IF NOT EXISTS or a temporary table of its name",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE IF NOT EXISTS a (id INT PRIMARY KEY);\nCREATE TABLE b (id INT);\nCREATE TEMPORARY TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT);\nDROP TEMPORARY TABLE IF EXISTS c;",
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "Creates table `a` without a primary key.", []int{0}},
				{2, "CREATE TABLE b (id INT)", "Creates table `b` without a primary key.", []int{0}},
				{4, "CREATE TABLE c (id INT)", "Creates table `c` without a primary key.", []int{0}},
			},
		},
		{
			name: "settled by a later key, drop, or rename",
			sql:  "CREATE TABLE a (id INT);\nALTER TABLE a ADD PRIMARY KEY (id);\nCREATE TABLE b (id INT);\nDROP TABLE b;\nCREATE TABLE c (id INT);\nRENAME TABLE c TO c2;\nCREATE DATABASE d;\nCREATE TABLE d.e (id INT);\nDROP DATABASE d;",
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE p DROP PRIMARY KEY, ADD PRIMARY KEY (id, note);\nALTER TABLE trig DROP PRIMARY KEY;\nALTER TABLE trig ADD PRIMARY KEY (id);",
			want: []finding{{0, "DROP PRIMARY KEY", "Removes the primary key of `t`, and the change doesn't add one back.", []int{0}}},
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
			want:   []finding{{-1, "", "TiDB doesn't support online migration.\n\ngh-ost needs a binary log, which TiDB doesn't have. Turn off **Online migration** for this change.", []int{0}}},
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
			want:    []finding{{-1, "", "Prior backup needs the backup database `bbdataarchive`, which doesn't exist on this instance.\n\nCreate it, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "over the size limit",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;",
			change: review.Change{PriorBackup: true, MaxBackupSize: 10},
			want:   []finding{{-1, "", "The change is over prior backup's size limit of 10 bytes.\n\nSplit it into smaller changes, or turn off **Prior backup**.", []int{0}}},
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
				{0, "BATCH ON id LIMIT 2 DELETE FROM t WHERE n > 0", "Prior backup can't back up a TiDB `BATCH` statement.\n\nTurn off **Prior backup** for this change.", []int{0}},
				{1, "REPLACE INTO p VALUES (3, 'c')", "Prior backup can't back up `REPLACE` or `INSERT ... ON DUPLICATE KEY UPDATE`.\n\nTurn off **Prior backup** for this change.", []int{0}},
				{2, "INSERT INTO p VALUES (1, 'a') ON DUPLICATE KEY UPDATE note = 'b'", "Prior backup can't back up `REPLACE` or `INSERT ... ON DUPLICATE KEY UPDATE`.\n\nTurn off **Prior backup** for this change.", []int{0}},
				{3, "EXPLAIN ANALYZE UPDATE t SET n = 1 WHERE id = 1", "Prior backup can't back up `EXPLAIN ANALYZE` of a statement that changes rows.\n\nRemove `EXPLAIN ANALYZE`, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a table of another database",
			sql:    "UPDATE other.t SET n = 1 WHERE id = 1;\nDELETE FROM SHOP.p WHERE id = 2;",
			change: backup,
			want:   []finding{{0, "UPDATE other.t SET n = 1 WHERE id = 1", "Prior backup can't back up `other.t`, a table of another database.\n\nRun that statement in a change on its own database, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a table of another database deleted from twice",
			sql:    "DELETE a, b FROM other.t AS a JOIN other.t AS b ON a.id = b.id WHERE a.id = 1;",
			change: backup,
			want:   []finding{{0, "DELETE a, b FROM other.t AS a JOIN other.t AS b ON a.id = b.id WHERE a.id = 1", "Prior backup can't back up `other.t`, a table of another database.\n\nRun that statement in a change on its own database, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "more than five statements over two tables",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nDELETE FROM shop.t WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE p SET note = '' WHERE id = 1;",
			change: backup,
			want:   []finding{{5, "UPDATE p SET note = '' WHERE id = 1", "Prior backup can't back up more than 5 statements across several tables.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "more than five statements on one table, one with WITH",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nWITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x);\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE t SET n = 6 WHERE id = 6;",
			change: backup,
			want:   []finding{{2, "WITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x)", "Prior backup can't back up a `WITH` clause when more than 5 statements change one table.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "more than five statements on one table, with WITH in feature gates",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\n/*T![clustered_index] WITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x) */;\n/*T![clustered_index] /* note */ WITH x AS (SELECT 4 AS id) UPDATE t SET n = 4 WHERE id IN (SELECT id FROM x) */;\n/*T![unsupported_feature_xyz] x */ /*T![clustered_index] WITH x AS (SELECT 5 AS id) UPDATE t SET n = 5 WHERE id IN (SELECT id FROM x) */;\n/*T![clustered_index] */ /* c */ /*T![clustered_index] WITH x AS (SELECT 6 AS id) UPDATE t SET n = 6 WHERE id IN (SELECT id FROM x) */;\n/*T![unsupported_feature_xyz] x */ WITH y AS (SELECT 7 AS id) UPDATE t SET n = 7 WHERE id IN (SELECT id FROM y);",
			change: backup,
			want:   []finding{{6, "WITH y AS (SELECT 7 AS id) UPDATE t SET n = 7 WHERE id IN (SELECT id FROM y)", "Prior backup can't back up a `WITH` clause when more than 5 statements change one table.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			// TiDB drops the opener of a gate that never closes and runs the
			// WITH after it as ordinary SQL, so the rule reads it.
			name:   "more than five statements on one table, with WITH after a gate that never closes",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nUPDATE t SET n = 3 WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\n/*T![clustered_index] WITH x AS (SELECT 6 AS id) UPDATE t SET n = 6 WHERE id IN (SELECT id FROM x);",
			change: backup,
			want:   []finding{{5, "WITH x AS (SELECT 6 AS id) UPDATE t SET n = 6 WHERE id IN (SELECT id FROM x)", "Prior backup can't back up a `WITH` clause when more than 5 statements change one table.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "past the limit with a WITH clause on a statement not certain",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nUPDATE t SET n = 3 WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nWITH x AS (SELECT 1 AS id) UPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id IN (SELECT id FROM x);",
			change: backup,
			want:   []finding{{-1, "", "Prior backup can't back up a `WITH` clause when the change has more than 5 statements to back up.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
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
			want:   []finding{{-1, "", "Prior backup can't back up more than 5 statements across several tables.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "five statements over five tables and a multi-table UPDATE, past the limit",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nUPDATE plain SET v = 1 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;\nUPDATE trig SET id = 2 WHERE id = 1;\nUPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id = 2;",
			change: backup,
			want:   []finding{{-1, "", "Prior backup can't back up more than 5 statements across several tables.\n\nSplit the change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a table the change creates",
			sql:    "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nDELETE t FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1;\nDELETE FROM nope WHERE a = 1;",
			change: backup,
			want: []finding{
				{1, "UPDATE n2 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `n2` doesn't exist yet.\n\nCreate `n2` in an earlier change, or turn off **Prior backup**.", []int{0}},
				{2, "DELETE t FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1", "Prior backup runs before the change, when `n2` doesn't exist yet.\n\nCreate `n2` in an earlier change, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a column the change adds or renames to",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;\nUPDATE t SET n = 2 WHERE t.gone = 2 ORDER BY id;\nALTER TABLE p RENAME COLUMN note TO memo;\nUPDATE t JOIN p ON p.memo = t.code SET t.n = 1 WHERE t.id = 3;",
			change: backup,
			want: []finding{
				{1, "UPDATE t SET n = 1 WHERE gone = 1", "Prior backup runs before the change, when column `gone` of `t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
				{2, "UPDATE t SET n = 2 WHERE t.gone = 2 ORDER BY id", "Prior backup runs before the change, when column `gone` of `t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
				{4, "UPDATE t JOIN p ON p.memo = t.code SET t.n = 1 WHERE t.id = 3", "Prior backup runs before the change, when column `memo` of `p` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a column USING names",
			sql:    "ALTER TABLE t ADD COLUMN note INT;\nUPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1;\nUPDATE t JOIN p USING (id) SET t.n = 1 WHERE t.id = 2;\nUPDATE t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) SET t.n = 1 WHERE t.id = 3;",
			change: backup,
			want:   []finding{{1, "UPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1", "Prior backup runs before the change, when column `note` of `t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a view as the target",
			sql:    "UPDATE pv SET note = 'z' WHERE id = 1;\nDELETE FROM other.pv WHERE id = 2;",
			change: backup,
			want: []finding{
				{0, "UPDATE pv SET note = 'z' WHERE id = 1", "Prior backup can't back up `pv`: it is a view.\n\nTurn off **Prior backup** for this change.", []int{0}},
				{1, "DELETE FROM other.pv WHERE id = 2", "Prior backup can't back up `other.pv`, a table of another database.\n\nRun that statement in a change on its own database, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a DELETE target the FROM list does not name",
			sql:    "DELETE T FROM t WHERE id = 1;\nDELETE y FROM t AS x WHERE x.id = 2;",
			change: backup,
			want:   []finding{{1, "DELETE y FROM t AS x WHERE x.id = 2", "Prior backup can't find the `DELETE` target `y` in the `FROM` list.\n\nWrite the target as the `FROM` list writes it.", []int{0}}},
		},
		{
			name:   "a CTE in a subquery, whose name hides a table the change makes",
			sql:    "CREATE TABLE x (id INT);\nUPDATE t SET n = 1 WHERE EXISTS (WITH x AS (SELECT 1 AS id) SELECT 1 FROM x);",
			change: backup,
		},
		{
			name:   "a DELETE after a WITH clause, whose targets may be its queries",
			sql:    "WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2;\n/*!80000 WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2 */;\n/*T! WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2 */;\n/*!80000 -- note\n WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2 */;",
			change: backup,
		},
		{
			name:   "a view OR REPLACE makes, or a table IF NOT EXISTS makes",
			sql:    "CREATE OR REPLACE VIEW nv AS SELECT 1 AS id;\nUPDATE nv SET id = 2 WHERE id = 1;\nCREATE TABLE IF NOT EXISTS t (gone INT);\nUPDATE t SET n = 1 WHERE gone = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE nv SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `nv` doesn't exist yet.\n\nCreate `nv` in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "statements before a USE",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;\nUPDATE n3 SET id = 1 WHERE id = 2;\nUSE other;\nCREATE TABLE n3 (id INT PRIMARY KEY);\nUPDATE t SET n = 2 WHERE gone = 2;",
			change: backup,
			want:   []finding{{1, "UPDATE t SET n = 1 WHERE gone = 1", "Prior backup runs before the change, when column `gone` of `t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "what the rule cannot tell",
			sql:    "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE id IN (SELECT gone FROM t);\nUPDATE t SET gone = 1 WHERE id = 1;\nWITH x AS (SELECT 1 AS gone) UPDATE t SET n = 3 WHERE id IN (SELECT gone FROM x);",
			change: backup,
		},
	})
}

// TestInFeatureGateStaysLinear: openers of gates that never close, ahead of a
// WITH, are answered by one scan. Asking parser.FeatureGate again for each one
// took seconds for 8,000 of them.
// The budget below is CPU time, not clock time.
func TestInFeatureGateStaysLinear(t *testing.T) {
	sql := strings.Repeat("/*T![clustered_index] ", 20000) + "WITH x AS (SELECT 1 AS id) SELECT id FROM x"
	start := cpuTime(t)
	if inFeatureGate(sql, review.Range{Start: 0, End: len(sql)}, strings.Index(sql, "WITH")) {
		t.Errorf("a WITH after gates that never close reads as gated")
	}
	if elapsed := cpuSince(t, start); elapsed > 2*time.Second {
		t.Errorf("inFeatureGate over 20,000 unclosed gates took %v of CPU time", elapsed)
	}
}

// TestInFeatureGateAgrees holds inFeatureGate, which stops at the first gate
// that never closes, to the loop that reads on past it, from every offset of
// random gate-heavy input.
func TestInFeatureGateAgrees(t *testing.T) {
	readOn := func(sql string, r review.Range, offset int) bool {
		i := r.Start
		for {
			i = skipTrivia(sql, i, r.End)
			if i >= offset {
				return false
			}
			body, closeStart, ok := parser.FeatureGate(sql[:r.End], i)
			if !ok {
				return false
			}
			if closeStart < 0 {
				i = body
				continue
			}
			if offset < closeStart {
				return true
			}
			i = closeStart + len("*/")
		}
	}
	pieces := []string{
		"/*T![clustered_index] ", "/*T![no_such_feature] ", "/*T! ", "/*! ", "/* c */ ", "/* ", "*/ ",
		"'", "'*/' ", "-- */\n", "# x\n", " ", "WITH ", "x ", "\n",
	}
	rng := rand.New(rand.NewSource(1))
	for iter := 0; iter < 20000; iter++ {
		var sb strings.Builder
		for k := rng.Intn(12); k >= 0; k-- {
			sb.WriteString(pieces[rng.Intn(len(pieces))])
		}
		sql := sb.String()
		r := review.Range{Start: 0, End: len(sql)}
		for offset := 0; offset <= len(sql); offset++ {
			if got, want := inFeatureGate(sql, r, offset), readOn(sql, r, offset); got != want {
				t.Fatalf("inFeatureGate(%q, %d) = %v, reading on gives %v", sql, offset, got, want)
			}
		}
	}
}

// TestStatementBodies pins Result.Bodies: each statement from its first
// token to its last, without the whitespace, comments, and semicolons
// around it.
func TestStatementBodies(t *testing.T) {
	sql := "SELECT 1;\n-- a note\nCREATE TABLE t (a INT PRIMARY KEY /*T![clustered_index] CLUSTERED */) /* trailing */;"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Bodies) != len(result.Statements) {
		t.Fatalf("%d bodies for %d statements", len(result.Bodies), len(result.Statements))
	}
	var got []string
	for i, body := range result.Bodies {
		if statement := result.Statements[i]; body.Start < statement.Start || body.End > statement.End {
			t.Errorf("body %v lies outside statement %v", body, statement)
		}
		got = append(got, sql[body.Start:body.End])
	}
	want := []string{"SELECT 1", "CREATE TABLE t (a INT PRIMARY KEY /*T![clustered_index] CLUSTERED */)"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bodies mismatch (-want +got):\n%s", diff)
	}
}
