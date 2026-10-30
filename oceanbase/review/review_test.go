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
				{
					Name:    "mixk",
					Columns: []*metadata.ColumnMetadata{{Name: "email", Nullable: true}, {Name: "code"}},
					Indexes: []*metadata.IndexMetadata{
						{Name: "uk_email", Unique: true, Expressions: []string{"email"}},
						{Name: "uk_f", Unique: true, Expressions: []string{"lower(`code`)"}},
					},
				},
				{
					Name:    "en",
					Columns: []*metadata.ColumnMetadata{{Name: "e", Type: "enum('float','json')"}},
					Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"e"}}},
				},
				{
					Name:    "fl",
					Columns: []*metadata.ColumnMetadata{{Name: "f", Type: "float"}, {Name: "e", Type: "enum('float','other')"}},
					Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"f"}}},
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
			want: []finding{{1, "DROP INDEX uk_code ON t", "drops unique key uk_code of t", []int{0}}},
		},
		{
			name: "an index named with its database, across a USE",
			sql:  "DROP INDEX uk_code ON shop.t;\nCREATE INDEX uk_code ON shop.t (code);\nUSE other;\nDROP INDEX uk_code ON shop.t;",
			want: []finding{{0, "DROP INDEX uk_code ON shop.t", "drops unique key uk_code of shop.t", []int{0}}},
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
			sql:  "CREATE UNIQUE INDEX ux ON T (c);\nDROP INDEX ux ON t;\nCREATE PROCEDURE pr() CREATE UNIQUE INDEX uy ON nokey (a);\nDROP INDEX uy ON nokey;\nALTER EVENT ev DO CREATE UNIQUE INDEX uz ON nokey (a);\nDROP INDEX uz ON nokey;",
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
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}},
				{5, "CREATE TABLE f AS SELECT 1 AS id", "creates table f without a primary key", []int{0}},
			},
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

// keyed is shop with tables for the primary key a DROP COLUMN drops: k1 to
// k3 with a key on id and a column v; pair with a key of two columns, pchk
// with one whose check reads both, and pairu, pairq, pairp, pairn, pairw,
// and pairk with one beside a unique key of b: of whole columns, of a prefix
// of b in the unique key, of a prefix of b in the key, of lengths not given,
// beside one of b and v, and beside one of b that is not unique; gen with a
// generated column, gen2 with one that names the key only in a literal, chk
// with a check of the key and v, solo with a check of the key alone, fx with
// a functional index, and part with partitioning, each reading the key;
// kpart partitioned by KEY (), which reads the key, and kpart0, without a
// key, by KEY (), which reads its unique key and not its other index, and
// kpart2 by KEY () with a unique key of u and one of u and v; inv whose
// other column is invisible, and inv2 with a visible one too; defx whose
// other column's default reads the key, and deflit whose literal default
// names it; uq with a unique key of the key and v, uqv with one of the key
// and v beside one of v, uqvp the same with a prefix of v, uq1 with one of
// the key alone, and kv with a key of v and the key that is not unique; and
// self with a foreign key to itself.
func keyed() []review.Target {
	s := shop()
	pk := func(columns ...string) []*metadata.IndexMetadata {
		return []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: columns}}
	}
	// pairKey is a key of a and b, of the lengths given, beside another
	// index.
	pairKey := func(lengths []int64, other *metadata.IndexMetadata) []*metadata.IndexMetadata {
		return []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a", "b"}, KeyLength: lengths}, other}
	}
	tables := &s.Schemas[0].Tables
	for _, name := range []string{"k1", "k2", "k3"} {
		*tables = append(*tables, &metadata.TableMetadata{Name: name, Columns: columns("id", "v"), Indexes: pk("id")})
	}
	*tables = append(*tables,
		&metadata.TableMetadata{Name: "pair", Columns: columns("a", "b", "v"), Indexes: pk("a", "b")},
		&metadata.TableMetadata{Name: "pairu", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{-1, -1}, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}, KeyLength: []int64{-1}})},
		&metadata.TableMetadata{Name: "pairq", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{-1, -1}, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}, KeyLength: []int64{5}})},
		&metadata.TableMetadata{Name: "pairp", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{-1, 10}, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}, KeyLength: []int64{-1}})},
		&metadata.TableMetadata{Name: "pairn", Columns: columns("a", "b", "v"), Indexes: pairKey(nil, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}})},
		&metadata.TableMetadata{Name: "pairw", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{-1, -1}, &metadata.IndexMetadata{Name: "ubv", Unique: true, Expressions: []string{"b", "v"}, KeyLength: []int64{-1, -1}})},
		&metadata.TableMetadata{Name: "pairk", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{-1, -1}, &metadata.IndexMetadata{Name: "kb", Expressions: []string{"b"}, KeyLength: []int64{-1}})},
		&metadata.TableMetadata{Name: "kpref", Columns: columns("a", "b", "v"), Indexes: pairKey([]int64{10, -1}, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}, KeyLength: []int64{-1}}), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_KEY}}},
		&metadata.TableMetadata{Name: "kprefn", Columns: columns("a", "b", "v"), Indexes: pairKey(nil, &metadata.IndexMetadata{Name: "ub", Unique: true, Expressions: []string{"b"}, KeyLength: []int64{-1}}), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_KEY}}},
		&metadata.TableMetadata{Name: "pchk", Columns: columns("a", "b", "v"), Indexes: pk("a", "b"), CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck", Expression: "(`a` < `b`)"}}},
		&metadata.TableMetadata{
			Name:    "gen",
			Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "g", Generation: &metadata.GenerationMetadata{Expression: "(`id` + 1)"}}},
			Indexes: pk("id"),
		},
		&metadata.TableMetadata{
			Name:    "gen2",
			Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "g", Generation: &metadata.GenerationMetadata{Expression: "concat(_latin1'id',`v`)"}}},
			Indexes: pk("id"),
		},
		&metadata.TableMetadata{Name: "chk", Columns: columns("id", "v"), Indexes: pk("id"), CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck", Expression: "(`ID` < `v`)"}}},
		&metadata.TableMetadata{Name: "solo", Columns: columns("id", "v"), Indexes: pk("id"), CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck", Expression: "(`id` > 0)"}}},
		&metadata.TableMetadata{Name: "fx", Columns: columns("id", "v"), Indexes: append(pk("id"), &metadata.IndexMetadata{Name: "fxi", Expressions: []string{"((`id` + `v`))"}})},
		&metadata.TableMetadata{Name: "part", Columns: columns("id", "v"), Indexes: pk("id"), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Expression: "`id`"}}},
		&metadata.TableMetadata{Name: "kpart", Columns: columns("id", "v"), Indexes: pk("id"), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_KEY}}},
		&metadata.TableMetadata{Name: "kpart0", Columns: columns("u", "v", "w"), Indexes: []*metadata.IndexMetadata{{Name: "uu", Unique: true, Expressions: []string{"u"}}, {Name: "iv", Expressions: []string{"v"}}}, Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_KEY}}},
		&metadata.TableMetadata{Name: "inv", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v", IsInvisible: true}}, Indexes: pk("id")},
		&metadata.TableMetadata{Name: "inv2", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v", IsInvisible: true}, {Name: "w"}}, Indexes: pk("id")},
		&metadata.TableMetadata{Name: "defx", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v", Default: "((`id` + 1))"}}, Indexes: pk("id")},
		&metadata.TableMetadata{Name: "deflit", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v", Default: "'id'"}}, Indexes: pk("id")},
		&metadata.TableMetadata{Name: "uq", Columns: columns("id", "v", "w", "x"), Indexes: append(pk("id"), &metadata.IndexMetadata{Name: "u", Unique: true, Expressions: []string{"id", "v"}})},
		&metadata.TableMetadata{Name: "kv", Columns: columns("id", "v"), Indexes: append(pk("id"), &metadata.IndexMetadata{Name: "k", Expressions: []string{"v", "id"}})},
		&metadata.TableMetadata{Name: "uqv", Columns: columns("id", "v", "w"), Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}, KeyLength: []int64{-1}},
			{Name: "u", Unique: true, Expressions: []string{"v"}, KeyLength: []int64{-1}},
			{Name: "uv", Unique: true, Expressions: []string{"id", "v"}, KeyLength: []int64{-1, -1}},
		}},
		&metadata.TableMetadata{Name: "uqvp", Columns: columns("id", "v", "w"), Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}, KeyLength: []int64{-1}},
			{Name: "u", Unique: true, Expressions: []string{"v"}, KeyLength: []int64{-1}},
			{Name: "uv", Unique: true, Expressions: []string{"id", "v"}, KeyLength: []int64{-1, 5}},
		}},
		&metadata.TableMetadata{Name: "kpart2", Columns: columns("u", "v", "w"), Indexes: []*metadata.IndexMetadata{
			{Name: "uu", Unique: true, Expressions: []string{"u"}, KeyLength: []int64{-1}},
			{Name: "uv", Unique: true, Expressions: []string{"u", "v"}, KeyLength: []int64{-1, -1}},
		}, Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_KEY}}},
		&metadata.TableMetadata{Name: "uq1", Columns: columns("id", "v"), Indexes: append(pk("id"), &metadata.IndexMetadata{Name: "u", Unique: true, Expressions: []string{"id"}})},
		&metadata.TableMetadata{
			Name:        "self",
			Columns:     columns("id", "up"),
			Indexes:     pk("id"),
			ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk_up", Columns: []string{"up"}, ReferencedTable: "self", ReferencedColumns: []string{"id"}}},
		},
	)
	return []review.Target{{Schema: s, BackupDatabaseExists: true}}
}

// TestRequirePrimaryKeyColumnDrop checks DROP COLUMN against the synced
// primary key. An ALTER TABLE of DROP COLUMN subcommands alone is judged:
// dropping every column of the key drops it, and a drop the server refuses,
// or one of only some columns of the key, ends the judging. Any other
// statement on a table leaves the table out, and one the rule cannot bound
// ends the judging.
func TestRequirePrimaryKeyColumnDrop(t *testing.T) {
	lost := func(statement int, text, table string) finding {
		return finding{statement, text, "removes the primary key of " + table + ", and the change adds none back", []int{0}}
	}
	cases := []ruleCase{
		{
			name: "every column of the key",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE pair DROP b, DROP COLUMN a;\nALTER TABLE shop.t DROP COLUMN id, DROP COLUMN code;",
			want: []finding{lost(0, "DROP COLUMN id", "plain"), lost(1, "DROP b", "pair"), lost(2, "DROP COLUMN id", "shop.t")},
		},
		{
			name: "no column of the key",
			sql:  "ALTER TABLE plain DROP COLUMN v;\nALTER TABLE k1 DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "k1")},
		},
		{
			name: "the key, after columns statements before dropped, and what read them",
			sql:  "ALTER TABLE t DROP COLUMN code;\nALTER TABLE t DROP COLUMN id;\nALTER TABLE gen DROP COLUMN g;\nALTER TABLE gen DROP COLUMN id;\nALTER TABLE solo DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "t"), lost(3, "DROP COLUMN id", "gen"), lost(4, "DROP COLUMN id", "solo")},
		},
		{
			name: "the key, before statements the rule passes or cannot follow",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nINSERT INTO plain VALUES (1);\nCREATE VIEW v1 AS SELECT v FROM plain;\nCREATE PROCEDURE pr() ALTER TABLE plain ADD PRIMARY KEY (v);\nSET @x = 1;\nCREATE DATABASE d;\nDROP DATABASE other;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		// What the server refuses, or what follows a drop of only some columns
		// of the key, is not known: no drop after it is judged, and a finding
		// before it stands.
		{
			name: "both columns of a key of two over two statements, a unique key of the other keeping its values apart",
			sql:  "ALTER TABLE pairu DROP COLUMN a;\nALTER TABLE pairu DROP COLUMN b;\nALTER TABLE pairq DROP COLUMN a;\nALTER TABLE pairq DROP COLUMN b;",
			want: []finding{lost(1, "DROP COLUMN b", "pairu"), lost(3, "DROP COLUMN b", "pairq")},
		},
		{name: "one column of a key of two whose other is a prefix, beside a unique key of it", sql: "ALTER TABLE pairp DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "one column of a key of two whose lengths the synced schema lacks, beside a unique key of the other", sql: "ALTER TABLE pairn DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "one column of a key of two, beside a unique key of the other and a column out of the key", sql: "ALTER TABLE pairw DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "one column of a key of two, beside a key of the other that is not unique", sql: "ALTER TABLE pairk DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "some columns of the key", sql: "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE pair DROP COLUMN a;\nALTER TABLE k1 DROP COLUMN id;", want: []finding{lost(0, "DROP COLUMN id", "plain")}},
		{name: "a column the table lacks", sql: "ALTER TABLE plain DROP COLUMN id, DROP COLUMN missing;\nALTER TABLE k1 DROP COLUMN id;"},
		{name: "one column twice", sql: "ALTER TABLE plain DROP COLUMN id, DROP COLUMN ID;"},
		{name: "every column of the table", sql: "ALTER TABLE plain DROP COLUMN id, DROP COLUMN v;"},
		{name: "a column a foreign key of another table references", sql: "ALTER TABLE p DROP COLUMN id;"},
		{name: "a column a foreign key of the table holds", sql: "ALTER TABLE t DROP COLUMN id, DROP COLUMN p_id;"},
		{name: "a column a foreign key of the table to itself references", sql: "ALTER TABLE self DROP COLUMN id;"},
		{name: "a column a generated column reads", sql: "ALTER TABLE gen DROP COLUMN id;"},
		{name: "a column a check reads with another", sql: "ALTER TABLE chk DROP COLUMN id;"},
		{name: "every column of the key, which a check reads", sql: "ALTER TABLE pchk DROP a, DROP b;"},
		{name: "a column a functional index reads", sql: "ALTER TABLE fx DROP COLUMN id;"},
		{name: "a column the table is partitioned by", sql: "ALTER TABLE part DROP COLUMN id;"},
		{name: "a column of a table partitioned by KEY ()", sql: "ALTER TABLE kpart DROP COLUMN id;"},
		{
			name: "the key, after a column of only some unique keys dropped from a table without a primary key partitioned by KEY ()",
			sql:  "ALTER TABLE kpart2 DROP COLUMN v;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, after a column of the key held by a prefix dropped from a table partitioned by KEY (), which leaves the prefix out",
			sql:  "ALTER TABLE kpref DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{name: "a column of the key of a table partitioned by KEY (), of lengths the synced schema lacks", sql: "ALTER TABLE kprefn DROP COLUMN a;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the unique key's column of a table without a primary key partitioned by KEY ()", sql: "ALTER TABLE kpart0 DROP COLUMN u;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, after a column out of any key dropped from a table without a primary key partitioned by KEY ()",
			sql:  "ALTER TABLE kpart0 DROP COLUMN v;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, after a column out of the key dropped from a table partitioned by KEY ()",
			sql:  "ALTER TABLE kpart DROP COLUMN v;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, after a read-only transaction started, which the drop commits",
			sql:  "START TRANSACTION READ ONLY;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{name: "every visible column of the table", sql: "ALTER TABLE inv DROP COLUMN id;"},
		{name: "a column an expression default reads", sql: "ALTER TABLE defx DROP COLUMN id;"},
		{name: "a column of a table the synced schema lacks", sql: "ALTER TABLE missing DROP COLUMN id;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a column of a table of another database", sql: "ALTER TABLE other.plain DROP COLUMN id;\nALTER TABLE plain DROP COLUMN id;"},
		// Any other statement on a table leaves the table out: no later drop
		// of it is judged, and a finding on it is withdrawn.
		{
			name: "another subcommand of the statement",
			sql:  "ALTER TABLE plain DROP COLUMN id, ADD COLUMN w INT;\nALTER TABLE k1 DROP COLUMN id, ALTER CHECK missing ENFORCED;\nALTER TABLE k2 DROP PRIMARY KEY, DROP COLUMN id;",
			want: []finding{{2, "DROP PRIMARY KEY", "removes the primary key of k2, and the change adds none back", []int{0}}},
		},
		{
			name: "a statement after on the table, which may give it a key back, or on another",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE k1 DROP COLUMN id;\nALTER TABLE k2 DROP COLUMN id;\nALTER TABLE plain ADD PRIMARY KEY (v);\nRENAME TABLE k1 TO k9;\nCREATE INDEX i ON nokey (a);",
			want: []finding{lost(2, "DROP COLUMN id", "k2")},
		},
		{
			name: "every column of the key, with an invisible one left besides a visible one",
			sql:  "ALTER TABLE inv2 DROP COLUMN id;",
			want: []finding{lost(0, "DROP COLUMN id", "inv2")},
		},
		{
			name: "a table of another database renamed to its name, before or after, which lands in the current database",
			sql:  "ALTER TABLE other.t RENAME TO k1;\nALTER TABLE k1 DROP COLUMN id;\nALTER TABLE k2 DROP COLUMN id;\nALTER TABLE other.t2 RENAME TO k2;",
		},
		{
			name: "a finding before a statement the rule cannot follow, withdrawn by a later statement on its table",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE k1 DROP COLUMN id;\nCREATE VIEW v1 AS SELECT 1;\nALTER TABLE plain ADD PRIMARY KEY (v);",
			want: []finding{lost(1, "DROP COLUMN id", "k1")},
		},
		{
			name: "findings before a USE of another database, after which a name without one is not the target's, but one with it is",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE k1 DROP COLUMN id;\nUSE other;\nALTER TABLE plain ADD PRIMARY KEY (v);\nALTER TABLE shop.k1 ADD PRIMARY KEY (v);",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, after a USE of the target's database, and not after one naming it in another case, where names compare with case",
			sql:  "USE shop;\nALTER TABLE plain DROP COLUMN id;\nUSE SHOP;\nALTER TABLE k1 DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{name: "a USE of another database, then of the target's, before", sql: "USE other;\nUSE shop;\nALTER TABLE plain DROP COLUMN id;"},
		// Any statement but one of DROP COLUMN alone or one that changes no
		// schema object may fail: no drop after it is judged.
		{name: "another table altered before", sql: "ALTER TABLE nokey ADD COLUMN w INT;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a table made before", sql: "CREATE TABLE c (x INT PRIMARY KEY);\nALTER TABLE plain DROP COLUMN id;"},
		{name: "an index made on another table before", sql: "CREATE INDEX i ON nokey (a);\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a table the database lacks dropped before", sql: "DROP TABLE missing;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "another database dropped before", sql: "DROP DATABASE other;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a view made before", sql: "CREATE VIEW v1 AS SELECT 1;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a procedure made before", sql: "CREATE PROCEDURE pr() SELECT 1;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user made before", sql: "CREATE USER u;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a temporary table of its name dropped before", sql: "DROP TEMPORARY TABLE IF EXISTS plain;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "statements before that cannot fail",
			sql:  "BEGIN;\nSAVEPOINT s;\nSET @x = 1, @sql_require_primary_key = 'a', foreign_key_checks = OFF, @@session.unique_checks = 0, @h = X'41', @h2 = 0x41, @b = b'101', @b2 = 0b1, @n = -1, @f = - +2.5e3, @h3 = X'aB', @h4 = 0xaB, @b3 = B'01', unique_checks = DEFAULT;\nANALYZE TABLE missing;\nANALYZE TABLE nokey UPDATE HISTOGRAM ON missing;\nCOMMIT;\nUNLOCK TABLES;\nROLLBACK;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(8, "DROP COLUMN id", "plain")},
		},
		{name: "a row inserted before", sql: "INSERT INTO nokey VALUES (1);\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a query before", sql: "SELECT a FROM nokey;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the database's tables shown before", sql: "SHOW TABLES FROM shop;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "tables flushed before", sql: "FLUSH TABLES WITH READ LOCK;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, before tables flushed",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nFLUSH TABLES;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "a table truncated before", sql: "TRUNCATE TABLE nokey;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "tables locked before", sql: "LOCK TABLES plain WRITE;\nALTER TABLE plain DROP COLUMN id;\nUNLOCK TABLES;"},
		{name: "a rollback to a savepoint before", sql: "BEGIN;\nROLLBACK TO SAVEPOINT s;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a savepoint released before", sql: "BEGIN;\nSAVEPOINT s;\nRELEASE SAVEPOINT s;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a commit that ends the session before", sql: "COMMIT RELEASE;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a rollback that ends the session before", sql: "ROLLBACK RELEASE;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the transaction's characteristics set before", sql: "SET TRANSACTION ISOLATION LEVEL READ COMMITTED;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "foreign_key_checks set with another variable before", sql: "SET foreign_key_checks = 0, version = 'x';\nALTER TABLE plain DROP COLUMN id;"},
		{name: "foreign_key_checks set globally before", sql: "SET GLOBAL foreign_key_checks = 0;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "foreign_key_checks set globally by its name before", sql: "SET @@global.foreign_key_checks = 0;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "foreign_key_checks set to a value it lacks before", sql: "SET foreign_key_checks = 2;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a query before", sql: "SET @x = (SELECT a FROM nokey);\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set under a scope before", sql: "SET SESSION @x = 1;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, after a scope for a system variable before a user variable",
			sql:  "SET SESSION foreign_key_checks = 1, @x = 1;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{name: "a user variable set to a malformed hexadecimal value before", sql: "SET @x = X'4';\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, after foreign_key_checks set with :=",
			sql:  "SET foreign_key_checks := OFF;\nALTER TABLE plain DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "plain")},
		},
		{name: "a user variable set to a hexadecimal value of an uppercase prefix before", sql: "SET @x = 0X41;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a bit value of an uppercase prefix before", sql: "SET @x = 0B01;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to DEFAULT before", sql: "SET @x = DEFAULT;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a malformed bit value before", sql: "SET @x = b'102';\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a double out of range before", sql: "SET @x = -1e400;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a sign of a string before", sql: "SET @x = -'a';\nALTER TABLE plain DROP COLUMN id;"},
		{name: "a user variable set to a string with a character set before", sql: "SET @x = _utf8mb4'a';\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the character set set before", sql: "SET NAMES utf8mb4;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, before a truncate, tables locked, DML, a SHOW, or a SET that may fail",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nTRUNCATE TABLE plain;\nLOCK TABLES plain READ;\nUNLOCK TABLES;\nINSERT INTO missing VALUES (1);\nSELECT * FROM missing;\nSHOW TABLES FROM shop;\nSET TRANSACTION READ ONLY;\nSET sql_require_primary_key = ON;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before statements on its table that give it no key",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nCREATE INDEX iv ON plain (v);\nALTER TABLE plain ADD COLUMN w INT;\nDROP INDEX iv ON plain;\nCREATE TEMPORARY TABLE plain (a INT);\nDROP TEMPORARY TABLE plain;\nCREATE TABLE IF NOT EXISTS plain (a INT);",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before its table dropped or renamed, or before a table of its name made or renamed onto it, which fail",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE k1 DROP COLUMN id;\nALTER TABLE k2 DROP COLUMN id;\nALTER TABLE k3 DROP COLUMN id;\nALTER TABLE gen2 DROP COLUMN id;\nDROP TABLE plain;\nALTER TABLE k2 RENAME TO k9;\nCREATE TABLE k1 (a INT PRIMARY KEY);\nALTER TABLE nokey RENAME TO k3;\nRENAME TABLE pair TO gen2;",
			want: []finding{lost(1, "DROP COLUMN id", "k1"), lost(3, "DROP COLUMN id", "k3"), lost(4, "DROP COLUMN id", "gen2")},
		},
		// A statement the rule shows to fail withdraws nothing, and the change
		// stops at it.
		{
			name: "the key, before an ALTER TABLE that drops a column the table lacks and adds a key back",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE plain DROP COLUMN missing, ADD PRIMARY KEY (v);",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before an ALTER TABLE that drops a column dropped before and adds a key back",
			sql:  "ALTER TABLE k1 DROP COLUMN id;\nALTER TABLE k1 DROP COLUMN ID, ADD PRIMARY KEY (v);",
			want: []finding{lost(0, "DROP COLUMN id", "k1")},
		},
		{
			name: "the key, before a column dropped by a statement the rule does not judge, then dropped again with a key added back",
			sql:  "ALTER TABLE t DROP COLUMN id;\nCREATE VIEW v1 AS SELECT 1;\nALTER TABLE t DROP COLUMN code;\nALTER TABLE t DROP COLUMN code, ADD PRIMARY KEY (n);",
			want: []finding{lost(0, "DROP COLUMN id", "t")},
		},
		{
			name: "the key, before an ALTER TABLE that fails, and not after it",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE plain DROP COLUMN missing;\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "the key, before a column added and then dropped with a key added back", sql: "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE plain ADD COLUMN missing INT;\nALTER TABLE plain DROP COLUMN missing, ADD PRIMARY KEY (v);"},
		{
			name: "the key, before a rename onto a table the target holds, or onto itself",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO k1;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before an ALTER TABLE that adds a key back and renames onto a table the target holds",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE plain ADD PRIMARY KEY (v), RENAME TO k1;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before a rename of it onto itself",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before a rename that fails, and not a drop after it",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO k1;\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before a rename with a pair onto a name an earlier pair moved a table back onto",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE k1 TO spare, plain TO k1, k2 TO k1;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before a rename with a pair onto a name an earlier pair made",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO spare, k1 TO spare;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "the key, before a rename with pairs onto one name of another database", sql: "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO other.spare, k1 TO other.spare;"},
		{name: "the key, before a USE of another database and a rename with pairs onto one name there", sql: "ALTER TABLE plain DROP COLUMN id;\nUSE other;\nRENAME TABLE shop.plain TO spare, k1 TO spare;"},
		{name: "the key, before a rename with a pair onto a name an earlier pair made and a later one moved away", sql: "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO spare, spare TO s2, k1 TO spare;"},
		{name: "the key, before a rename onto a table renamed away before", sql: "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE k1 TO k9;\nRENAME TABLE plain TO k1;"},
		{name: "the key, before a rename onto a table an earlier pair moves away", sql: "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE k1 TO k9, plain TO k1;"},
		{name: "the key, before an ALTER TABLE that renames it onto itself", sql: "ALTER TABLE plain DROP COLUMN id;\nALTER TABLE plain RENAME TO plain;"},
		{name: "the key, before a rename onto a table dropped before", sql: "ALTER TABLE plain DROP COLUMN id;\nDROP TABLE k1;\nRENAME TABLE plain TO k1;"},
		{name: "the key, before a rename onto a table of another database", sql: "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO other.k1;"},
		// A temporary table of the name takes an ALTER TABLE and a DROP TABLE.
		{
			name: "the key, before a temporary table of its name made, given a key, and dropped",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nALTER TABLE plain ADD PRIMARY KEY (x);\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "the key, before a temporary table of its name dropped, and then the table", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nDROP TABLE plain;\nDROP TABLE plain;"},
		{name: "the key, before a temporary table of its name dropped as one, and then the table", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nDROP TEMPORARY TABLE plain;\nDROP TABLE plain;"},
		{name: "the key, before a temporary table of its name renamed away, and then the table dropped", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nALTER TABLE plain RENAME TO p2;\nDROP TABLE plain;"},
		{
			name: "the key, before a RENAME TABLE of its name onto a table the target holds, which MySQL runs on it, though a temporary table of its name exists",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nRENAME TABLE plain TO k1;\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "the key, before a RENAME TABLE of its name, though a temporary table of it exists", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE plain (x INT);\nRENAME TABLE plain TO p2;"},
		{name: "the key, before a temporary table of its name made in another database, and the table dropped", sql: "ALTER TABLE plain DROP COLUMN id;\nUSE other;\nCREATE TEMPORARY TABLE plain (x INT);\nDROP TABLE shop.plain;"},
		{
			name: "the key, before a temporary table renamed onto its name, given a key, and dropped",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (x INT);\nALTER TABLE tmp RENAME TO plain;\nALTER TABLE plain ADD PRIMARY KEY (x);\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "the key, before a temporary table renamed onto its name and dropped, and then the table", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (x INT);\nALTER TABLE tmp RENAME TO plain;\nDROP TABLE plain;\nDROP TABLE plain;"},
		{name: "the key, before a temporary table renamed onto its name loses a column the table lacks, and then both are dropped", sql: "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (missing INT, y INT);\nALTER TABLE tmp RENAME TO plain;\nALTER TABLE plain DROP COLUMN missing;\nDROP TABLE plain;\nDROP TABLE plain;"},
		{
			name: "the key, before a RENAME TABLE of a temporary table onto its name, which MySQL does not rename, and fails",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (x INT);\nRENAME TABLE tmp TO plain;\nDROP TABLE plain;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "the key, before a statement prepared and deallocated, which runs nothing",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nPREPARE s FROM 'ALTER TABLE plain ADD PRIMARY KEY (v)';\nDEALLOCATE PREPARE s;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{name: "a statement prepared before, which may fail", sql: "PREPARE s FROM 'SELECT 1';\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the key, before a prepared statement executed", sql: "ALTER TABLE plain DROP COLUMN id;\nEXECUTE s;"},
		// A DROP TABLE a foreign key keeps from running is not shown to fail: the
		// synced schema does not give the database the foreign key references.
		{name: "the key, before a DROP TABLE of it, which a foreign key of another table may name", sql: "ALTER TABLE p DROP COLUMN v;\nALTER TABLE plain DROP COLUMN id;\nDROP TABLE plain;"},
		{name: "a column a unique key of it and another column holds", sql: "ALTER TABLE uq DROP COLUMN id;"},
		{
			name: "every column of the key, with a unique key of it and another column whose rest another unique key keeps apart",
			sql:  "ALTER TABLE uqv DROP COLUMN id;",
			want: []finding{lost(0, "DROP COLUMN id", "uqv")},
		},
		{name: "a column a unique key of it and a prefix of another column holds, beside a unique key of the other", sql: "ALTER TABLE uqvp DROP COLUMN id;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "every column of the key, with a unique key of it alone, or of it and another column dropped too, or another key of it and another column",
			sql:  "ALTER TABLE uq DROP COLUMN w;\nALTER TABLE uq1 DROP COLUMN id;\nALTER TABLE uq DROP COLUMN id, DROP COLUMN v;\nALTER TABLE kv DROP COLUMN id;",
			want: []finding{lost(1, "DROP COLUMN id", "uq1"), lost(2, "DROP COLUMN id", "uq"), lost(3, "DROP COLUMN id", "kv")},
		},
		{
			name: "the key, with ALGORITHM or LOCK left at DEFAULT, last or alone, and not before another value",
			sql:  "ALTER TABLE plain DROP COLUMN id, ALGORITHM=DEFAULT;\nALTER TABLE k1 LOCK DEFAULT, DROP COLUMN id, algorithm = default, LOCK=DEFAULT;\nALTER TABLE k2 DROP COLUMN id, ALGORITHM=INSTANT, LOCK=NONE, ALGORITHM=DEFAULT, LOCK=DEFAULT;\nALTER TABLE k3 DROP COLUMN id, ALGORITHM=DEFAULT, ALGORITHM=INSTANT;",
			want: []finding{lost(0, "DROP COLUMN id", "plain"), lost(1, "DROP COLUMN id", "k1"), lost(2, "DROP COLUMN id", "k2")},
		},
		{name: "the key, with LOCK set to another value", sql: "ALTER TABLE plain DROP COLUMN id, LOCK=SHARED;"},
		{name: "the key, with ALGORITHM set to another value before LOCK left at DEFAULT", sql: "ALTER TABLE plain DROP COLUMN id, ALGORITHM=INSTANT, LOCK=DEFAULT;"},
		{name: "the key, with a quoted DEFAULT", sql: "ALTER TABLE plain DROP COLUMN id, ALGORITHM=`DEFAULT`;"},
		{
			name: "a generated column whose expression names the column only in a literal",
			sql:  "ALTER TABLE gen2 DROP COLUMN id;",
			want: []finding{lost(0, "DROP COLUMN id", "gen2")},
		},
		{name: "sql_generate_invisible_primary_key set on before", sql: "SET SESSION sql_generate_invisible_primary_key = ON;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "sql_generate_invisible_primary_key set off before, which needs a privilege", sql: "SET @@session.sql_generate_invisible_primary_key = OFF;\nALTER TABLE plain DROP COLUMN id;"},
		// A statement the rule cannot bound ends the judging.
		{name: "a USE", sql: "USE other;\nALTER TABLE shop.plain DROP COLUMN id;"},
		{name: "a CALL after the drop", sql: "ALTER TABLE plain DROP COLUMN id;\nCALL pr();"},
		{name: "sql_require_primary_key set on before", sql: "SET SESSION sql_require_primary_key = ON;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "sql_require_primary_key set to its default before", sql: "SET sql_require_primary_key = DEFAULT;\nALTER TABLE plain DROP COLUMN id;"},
		{name: "sql_require_primary_key set off before, which needs a privilege", sql: "SET SESSION sql_require_primary_key = OFF;\nALTER TABLE plain DROP COLUMN id;"},
		{
			name: "the key, before sql_require_primary_key set on",
			sql:  "ALTER TABLE plain DROP COLUMN id;\nSET sql_require_primary_key = ON;\nALTER TABLE k2 DROP COLUMN id;",
			want: []finding{lost(0, "DROP COLUMN id", "plain")},
		},
		{
			name: "a literal default that names the column",
			sql:  "ALTER TABLE deflit DROP COLUMN id;",
			want: []finding{lost(0, "DROP COLUMN id", "deflit")},
		},
		{name: "an event", sql: "CREATE EVENT ev ON SCHEDULE EVERY 1 DAY DO ALTER TABLE plain ADD PRIMARY KEY (v);\nALTER TABLE plain DROP COLUMN id;"},
		{name: "the target's database dropped after", sql: "ALTER TABLE plain DROP COLUMN id;\nDROP DATABASE shop;"},
	}
	for i := range cases {
		cases[i].targets = keyed()
	}
	runCases(t, review.RequirePrimaryKey, cases)
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name:    "a table named in another case, where names compare with case and where not",
			sql:     "ALTER TABLE PLAIN DROP COLUMN ID;",
			targets: []review.Target{{Schema: shop()}, {Schema: shop(), LowerCaseTableNames: 1}},
			want:    []finding{{0, "DROP COLUMN ID", "removes the primary key of PLAIN, and the change adds none back", []int{1}}},
		},
		{
			name:    "a rename with pairs onto one name in two cases, which collide only where names compare without case",
			sql:     "ALTER TABLE plain DROP COLUMN id;\nRENAME TABLE plain TO Spare, k1 TO spare;",
			targets: []review.Target{{Schema: keyed()[0].Schema}, {Schema: keyed()[0].Schema, LowerCaseTableNames: 1}},
			want:    []finding{{0, "DROP COLUMN id", "removes the primary key of plain, and the change adds none back", []int{1}}},
		},
		{
			name: "targets whose keys differ",
			sql:  "ALTER TABLE pair DROP COLUMN a;",
			targets: []review.Target{keyed()[0], {Schema: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{
				{Name: "pair", Columns: columns("a", "b"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"a"}}}},
			}}}}}, {}},
			want: []finding{{0, "DROP COLUMN a", "removes the primary key of pair, and the change adds none back", []int{1}}},
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
			name:   "a table of another database, which gh-ost would take from the task's",
			sql:    "ALTER TABLE other.plain ADD COLUMN m INT;\n",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE other.plain ADD COLUMN m INT", "gh-ost migrates plain of the task's database shop, not other.plain as the ALTER names it", []int{0}}},
		},
		{
			name:    "a table named with the task's database, in another case where names fold",
			sql:     "ALTER TABLE SHOP.plain ADD COLUMN m INT;",
			change:  ghost,
			targets: []review.Target{{Schema: shop(), BackupDatabaseExists: true, LowerCaseTableNames: 1}},
		},
		{
			name:   "a database named in another case where names are case-sensitive",
			sql:    "ALTER TABLE SHOP.plain ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE SHOP.plain ADD COLUMN m INT", "gh-ost migrates plain of the task's database shop, not SHOP.plain as the ALTER names it", []int{0}}},
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
			name:   "a column renamed by CHANGE, a name with a space in it",
			sql:    "ALTER TABLE plain CHANGE COLUMN v `v w` INT;",
			change: ghost,
			want:   []finding{{0, "CHANGE COLUMN v `v w` INT", "gh-ost renames column v to `v w` only with approve-renamed-columns, which Bytebase does not pass", []int{0}}},
		},
		{
			name:   "a column kept by CHANGE, a name with a space in it",
			sql:    "ALTER TABLE plain CHANGE COLUMN `v w` `v w` BIGINT;",
			change: ghost,
			want:   []finding{{0, "CHANGE COLUMN `v w` `v w` BIGINT", "gh-ost splits column name `v w` at whitespace and reads CHANGE as a rename, which it takes only with approve-renamed-columns", []int{0}}},
		},
		{
			name:   "a column renamed by RENAME COLUMN",
			sql:    "ALTER TABLE plain RENAME COLUMN v TO w;",
			change: ghost,
			want:   []finding{{0, "RENAME COLUMN v TO w", "gh-ost does not read RENAME COLUMN as a rename, and the migrated table loses the data of column v", []int{0}}},
		},
		{
			name:   "flag values gh-ost cannot read",
			sql:    "-- gh-ost = {\"max-load\":\"Threads_running\"}\nALTER TABLE plain ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{-1, "", `online migration cannot read gh-ost flag "max-load" as load map: "Threads_running"`, []int{0}}},
		},
		{
			name:   "flag values gh-ost's setters refuse",
			sql:    "-- gh-ost = {\"cut-over-lock-timeout-seconds\":\"0\",\"exponential-backoff-max-interval\":\"1\"}\nALTER TABLE plain ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{-1, "", `online migration cannot take gh-ost flag "cut-over-lock-timeout-seconds" = 0, outside 1 to 10`, []int{0}}},
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
			want:   []finding{{0, "ALTER TABLE nul ADD COLUMN m INT", "online migration needs a primary key or a unique key of NOT NULL columns, none of them FLOAT or JSON, on nul, which has none", []int{0}}},
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
			name:   "a nullable unique key and a functional one, neither usable",
			sql:    "ALTER TABLE mixk ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE mixk ADD COLUMN m INT", "online migration needs a primary key or a unique key of NOT NULL columns, none of them FLOAT or JSON, on mixk, which has none", []int{0}}},
		},
		{
			name:   "a unique key of an ENUM naming float",
			sql:    "ALTER TABLE en ADD COLUMN m INT;",
			change: ghost,
		},
		{
			name:   "a primary key of a FLOAT column",
			sql:    "ALTER TABLE fl ADD COLUMN m INT;",
			change: ghost,
			want:   []finding{{0, "ALTER TABLE fl ADD COLUMN m INT", "online migration needs a primary key or a unique key of NOT NULL columns, none of them FLOAT or JSON, on fl, which has none", []int{0}}},
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
			name:    "a target without a synced schema",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			change:  backup,
			targets: []review.Target{{}},
			want: []finding{
				{-1, "", "prior backup needs database bbdataarchive, which does not exist", []int{0}},
				{1, "DELETE FROM t WHERE id = 2", "prior backup cannot back up both UPDATE and DELETE on t", []int{0}},
			},
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
			name:   "a CTE in a subquery, whose name hides a table the change makes",
			sql:    "CREATE TABLE x (id INT);\nUPDATE t SET n = 1 WHERE EXISTS (WITH x AS (SELECT 1 AS id) SELECT 1 FROM x);",
			change: backup,
		},
		{
			name:   "a DELETE after a WITH clause, whose targets may be its queries",
			sql:    "WITH x AS (SELECT 1 AS id) DELETE T FROM t JOIN x ON t.id = x.id;\n/*!80000 WITH x AS (SELECT 1 AS id) DELETE T FROM t JOIN x ON t.id = x.id */;\n/*!80000 -- note\n WITH w AS (SELECT 1 AS id) DELETE y FROM t AS x WHERE x.id = 2 */;",
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
