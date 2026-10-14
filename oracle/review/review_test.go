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

// shop is a small database as Bytebase syncs one, its objects in the
// schema named "": P with a primary key; T with a primary key, a unique
// constraint, a plain and a unique index, a foreign key to P, a check, and
// a named NOT NULL constraint, which Oracle keeps as a check; NOKEY
// without a key; LONGT with a LONG column; V, a view; MV, a materialized
// view; and SEQ, a sequence.
func shop() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "SHOP",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name:    "P",
					Columns: columns("ID", "NOTE"),
					Indexes: []*metadata.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true}},
				},
				{
					Name:    "T",
					Columns: columns("ID", "CODE", "P_ID", "N"),
					Indexes: []*metadata.IndexMetadata{
						{Name: "PK_T", Primary: true, Unique: true, IsConstraint: true},
						{Name: "UK_CODE", Unique: true, IsConstraint: true},
						{Name: "IX_N"},
						{Name: "UX_N", Unique: true},
					},
					ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "FK_P", ReferencedTable: "P"}},
					CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "CK_N"}, {Name: "NN_CODE", Expression: `"CODE" IS NOT NULL`}},
				},
				{Name: "NOKEY", Columns: columns("A")},
				{
					Name:    "LONGT",
					Columns: []*metadata.ColumnMetadata{{Name: "ID", Type: "NUMBER"}, {Name: "BODY", Type: "LONG"}},
					Indexes: []*metadata.IndexMetadata{{Name: "PK_LONGT", Primary: true, Unique: true, IsConstraint: true}},
				},
			},
			Views:             []*metadata.ViewMetadata{{Name: "V"}, {Name: "LONGV", Columns: []*metadata.ColumnMetadata{{Name: "BODY", Type: "LONG"}}}},
			MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "MV"}},
			Sequences:         []*metadata.SequenceMetadata{{Name: "SEQ"}},
		}},
	}
}

func columns(names ...string) []*metadata.ColumnMetadata {
	out := make([]*metadata.ColumnMetadata, len(names))
	for i, n := range names {
		out[i] = &metadata.ColumnMetadata{Name: n, Type: "NUMBER"}
	}
	return out
}

func TestStatements(t *testing.T) {
	sql := "SET DEFINE OFF\nSELECT 1 FROM dual;\n-- a comment\nBEGIN NULL; END;\n/\nUPDATE t SET n = 1 WHERE id = 1 /* trailing */ ;\n"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range result.Statements {
		got = append(got, sql[r.Start:r.End])
	}
	want := []string{"SELECT 1 FROM dual", "BEGIN NULL; END;", "UPDATE t SET n = 1 WHERE id = 1"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("statements mismatch (-want +got):\n%s", diff)
	}
}

func TestSyntax(t *testing.T) {
	runCases(t, review.Syntax, []ruleCase{
		{name: "valid", sql: "SELECT 1 FROM dual; UPDATE t SET n = 1 WHERE id = 1"},
		{
			name: "one finding for the first error",
			sql:  "SELECT 1 FROM dual;\nSELEC 2;\nSELEC 3;",
			want: []finding{{1, "", `syntax error at or near "SELEC"`, []int{0}}},
		},
	})
	sql := "SELECT 1 FROM dual;\nSELEC 2;\nSELECT 3 FROM dual;"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Statements) != 3 {
		t.Errorf("Statements = %v, want all three listed", result.Statements)
	}
}

func TestRequireWhere(t *testing.T) {
	runCases(t, review.RequireWhere, []ruleCase{
		{
			name: "UPDATE and DELETE without WHERE",
			sql:  "UPDATE t SET n = 1;\nDELETE FROM shop.t;\nUPDATE t SET n = 1 WHERE id = 1;\nDELETE t x WHERE x.id = 1;\nUPDATE (SELECT * FROM t) SET n = 1;",
			want: []finding{
				{0, "UPDATE t SET n = 1", "UPDATE T has no WHERE clause", []int{0}},
				{1, "DELETE FROM shop.t", "DELETE FROM SHOP.T has no WHERE clause", []int{0}},
				{4, "UPDATE (SELECT * FROM t) SET n = 1", "UPDATE an inline view has no WHERE clause", []int{0}},
			},
		},
		{
			name: "inside PL/SQL",
			sql:  "BEGIN\n  IF 1 = 1 THEN DELETE FROM t; END IF;\n  FORALL i IN 1..2 UPDATE p SET note = 'x';\nEND;\n/\nCREATE OR REPLACE TRIGGER trg AFTER INSERT ON p BEGIN DELETE FROM \"Mixed\"; END;\n/\nCREATE PROCEDURE pr AS BEGIN UPDATE t SET n = 1 WHERE id = 2; END;",
			want: []finding{
				{0, "DELETE FROM t", "DELETE FROM T has no WHERE clause", []int{0}},
				{0, "UPDATE p SET note = 'x'", "UPDATE P has no WHERE clause", []int{0}},
				{1, `DELETE FROM "Mixed"`, `DELETE FROM "Mixed" has no WHERE clause`, []int{0}},
			},
		},
		{name: "under EXPLAIN PLAN", sql: "EXPLAIN PLAN FOR DELETE FROM t;"},
	})
}

func TestRequireIsNull(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "comparisons a condition tests",
			sql:  "SELECT * FROM t WHERE n = NULL OR NULL <> code;\nUPDATE t SET n = 1 WHERE NOT (code != NULL) AND (n ^= (NULL));\nSELECT * FROM t JOIN p ON t.p_id = NULL GROUP BY n HAVING n = NULL;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "NULL <> code", `"NULL <> code" is never true; use IS NOT NULL`, []int{0}},
				{1, "code != NULL", `"code != NULL" is never true; use IS NOT NULL`, []int{0}},
				{1, "n ^= (NULL)", `"n ^= (NULL)" is never true; use IS NOT NULL`, []int{0}},
				{2, "t.p_id = NULL", `"t.p_id = NULL" is never true; use IS NULL`, []int{0}},
				{2, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "the empty string, which Oracle treats as NULL",
			sql:  "SELECT * FROM t WHERE code = '' OR code <> N'' OR code = q'[]';\nSELECT * FROM t WHERE code = ' ';",
			want: []finding{
				{0, "code = ''", `"code = ''" is never true; use IS NULL`, []int{0}},
				{0, "code <> N''", `"code <> N''" is never true; use IS NOT NULL`, []int{0}},
				{0, "code = q'[]'", `"code = q'[]'" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "values, IS NULL, a simple CASE, and IN",
			sql:  "SELECT CASE n WHEN NULL THEN 1 END, DECODE(n, NULL, 1) FROM t WHERE n IS NULL AND n IN (NULL) AND LNNVL(n = NULL);",
		},
		{
			name: "CASE, MERGE, CHECK, START WITH, and multitable INSERT",
			sql:  "SELECT CASE WHEN n = NULL THEN 1 END FROM t START WITH n = '' CONNECT BY PRIOR id = p_id;\nMERGE INTO t USING p ON (t.id = NULL) WHEN MATCHED THEN UPDATE SET t.n = 1 WHERE p.note <> NULL DELETE WHERE t.n = NULL;\nCREATE TABLE c (a NUMBER CHECK (a <> NULL), b NUMBER, CONSTRAINT ck_b CHECK (b = ''));\nINSERT ALL WHEN a = NULL THEN INTO t (id) VALUES (a) SELECT 1 a FROM dual;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "n = ''", `"n = ''" is never true; use IS NULL`, []int{0}},
				{1, "t.id = NULL", `"t.id = NULL" is never true; use IS NULL`, []int{0}},
				{1, "p.note <> NULL", `"p.note <> NULL" is never true; use IS NOT NULL`, []int{0}},
				{1, "t.n = NULL", `"t.n = NULL" is never true; use IS NULL`, []int{0}},
				{2, "a <> NULL", `"a <> NULL" is never true; use IS NOT NULL`, []int{0}},
				{2, "b = ''", `"b = ''" is never true; use IS NULL`, []int{0}},
				{3, "a = NULL", `"a = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "PL/SQL conditions",
			sql:  "DECLARE v NUMBER;\nBEGIN\n  IF v = NULL THEN NULL; ELSIF v <> '' THEN NULL; END IF;\n  WHILE v ~= NULL LOOP EXIT WHEN v = NULL; CONTINUE WHEN v = ''; END LOOP;\n  CASE WHEN v = NULL THEN NULL; ELSE NULL; END CASE;\n  v := CASE WHEN v = NULL THEN 1 END;\nEND;\n/\nCREATE TRIGGER trg BEFORE UPDATE ON t FOR EACH ROW WHEN (new.n = NULL) BEGIN NULL; END;",
			want: []finding{
				{0, "v = NULL", `"v = NULL" is never true; use IS NULL`, []int{0}},
				{0, "v <> ''", `"v <> ''" is never true; use IS NOT NULL`, []int{0}},
				{0, "v ~= NULL", `"v ~= NULL" is never true; use IS NOT NULL`, []int{0}},
				{0, "v = NULL", `"v = NULL" is never true; use IS NULL`, []int{0}},
				{0, "v = ''", `"v = ''" is never true; use IS NULL`, []int{0}},
				{0, "v = NULL", `"v = NULL" is never true; use IS NULL`, []int{0}},
				{0, "v = NULL", `"v = NULL" is never true; use IS NULL`, []int{0}},
				{1, "new.n = NULL", `"new.n = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "MATCH_RECOGNIZE DEFINE",
			sql:  "SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id PATTERN (a) DEFINE a AS a.n = NULL);",
			want: []finding{{0, "a.n = NULL", `"a.n = NULL" is never true; use IS NULL`, []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "ALTER TYPE DROP ATTRIBUTE",
			sql:  "ALTER TYPE person_t DROP ATTRIBUTE (phone) CASCADE INCLUDING TABLE DATA;\nALTER TYPE hr.addr_t DROP ATTRIBUTE (zip, city) CASCADE;\nALTER TYPE person_t ADD ATTRIBUTE (email VARCHAR2(80)) CASCADE;",
			want: []finding{
				{0, "phone", "drops attribute PHONE of type PERSON_T", []int{0}},
				{1, "zip", "drops attribute ZIP of type HR.ADDR_T", []int{0}},
				{1, "city", "drops attribute CITY of type HR.ADDR_T", []int{0}},
			},
		},
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t CASCADE CONSTRAINTS PURGE;\nDROP VIEW v;\nDROP INDEX shop.ix_n;\nDROP PROCEDURE pr;\nDROP PACKAGE BODY pk;\nDROP PUBLIC SYNONYM ps;\nDROP MATERIALIZED VIEW LOG ON t;\nDROP USER u CASCADE;\nDROP TABLESPACE ts INCLUDING CONTENTS;\nDROP AUDIT POLICY ap;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t CASCADE CONSTRAINTS PURGE", "drops table T", []int{0}},
				{1, "DROP VIEW v", "drops view V", []int{0}},
				{2, "DROP INDEX shop.ix_n", "drops index SHOP.IX_N", []int{0}},
				{3, "DROP PROCEDURE pr", "drops procedure PR", []int{0}},
				{4, "DROP PACKAGE BODY pk", "drops package body PK", []int{0}},
				{5, "DROP PUBLIC SYNONYM ps", "drops synonym PS", []int{0}},
				{6, "DROP MATERIALIZED VIEW LOG ON t", "drops materialized view log on T", []int{0}},
				{7, "DROP USER u CASCADE", "drops user U", []int{0}},
				{8, "DROP TABLESPACE ts INCLUDING CONTENTS", "drops tablespace TS", []int{0}},
				{9, "DROP AUDIT POLICY ap", "drops audit policy AP", []int{0}},
			},
		},
		{
			name: "ALTER TABLE clauses that drop part of a table",
			sql:  "ALTER TABLE t DROP COLUMN n;\nALTER TABLE t DROP (n, \"Code\");\nALTER TABLE t DROP UNUSED COLUMNS;\nALTER TABLE t DROP PARTITION p1, p2 UPDATE GLOBAL INDEXES;\nALTER TABLE t DROP PARTITION FOR (10);\nALTER TABLE t DROP SUBPARTITION sp1;\nALTER TABLE t DROP (PERIOD FOR valid);",
			want: []finding{
				{0, "DROP COLUMN n", "drops column N of T", []int{0}},
				{1, `DROP (n, "Code")`, `drops columns N, "Code" of T`, []int{0}},
				{2, "DROP UNUSED COLUMNS", "drops the unused columns of T", []int{0}},
				{3, "DROP PARTITION p1, p2 UPDATE GLOBAL INDEXES", "drops partitions P1, P2 of T", []int{0}},
				{4, "DROP PARTITION FOR (10)", "drops partition FOR (10) of T", []int{0}},
				{5, "DROP SUBPARTITION sp1", "drops subpartition SP1 of T", []int{0}},
				{6, "DROP (PERIOD FOR valid)", "drops period VALID of T", []int{0}},
			},
		},
		{
			name: "constraints belong to DisallowDropConstraint, and SET UNUSED drops nothing",
			sql:  "ALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t SET UNUSED COLUMN n;",
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{
			name: "TRUNCATE TABLE, CLUSTER, and partitions",
			sql:  "TRUNCATE TABLE t;\nTRUNCATE CLUSTER shop.cl;\nALTER TABLE t TRUNCATE PARTITION p1, p2;\nALTER TABLE t TRUNCATE SUBPARTITION FOR (1, 2);",
			want: []finding{
				{0, "TRUNCATE TABLE t", "truncates table T", []int{0}},
				{1, "TRUNCATE CLUSTER shop.cl", "truncates cluster SHOP.CL", []int{0}},
				{2, "TRUNCATE PARTITION p1, p2", "truncates partitions P1, P2 of T", []int{0}},
				{3, "TRUNCATE SUBPARTITION FOR (1, 2)", "truncates subpartition FOR (1, 2) of T", []int{0}},
			},
		},
	})
}

func TestDisallowRename(t *testing.T) {
	runCases(t, review.DisallowRename, []ruleCase{
		{
			name: "ALTER TABLE renames",
			sql:  "ALTER TABLE t RENAME TO t2;\nALTER TABLE p RENAME COLUMN note TO memo;\nALTER TABLE p RENAME CONSTRAINT pk_p TO pk_p2;",
			want: []finding{
				{0, "RENAME TO t2", "renames table T to T2", []int{0}},
				{1, "RENAME COLUMN note TO memo", "renames column NOTE of P to MEMO", []int{0}},
			},
		},
		{
			name: "RENAME of a synced table, view, and sequence",
			sql:  "RENAME t TO t3;\nRENAME v TO v2;\nRENAME seq TO seq2;\nRENAME unknown TO u2;",
			want: []finding{
				{0, "RENAME t TO t3", "renames table T to T3", []int{0}},
				{1, "RENAME v TO v2", "renames view V to V2", []int{0}},
			},
		},
		{
			name: "RENAME of what the change made",
			sql:  "CREATE TABLE n1 (id NUMBER);\nRENAME n1 TO n2;\nRENAME n2 TO n3;\nCREATE VIEW w AS SELECT 1 a FROM dual;\nRENAME w TO w2;\nCREATE SEQUENCE s1;\nRENAME s1 TO s2;\nDROP TABLE nokey;\nRENAME nokey TO nk;",
			want: []finding{
				{1, "RENAME n1 TO n2", "renames table N1 to N2", []int{0}},
				{2, "RENAME n2 TO n3", "renames table N2 to N3", []int{0}},
				{4, "RENAME w TO w2", "renames view W to W2", []int{0}},
			},
		},
		{
			name:    "without a synced schema",
			sql:     "RENAME t TO t3;\nCREATE TABLE n1 (id NUMBER);\nRENAME n1 TO n2;",
			targets: []review.Target{{}},
			want:    []finding{{2, "RENAME n1 TO n2", "renames table N1 to N2", []int{0}}},
		},
		{
			name: "after the session moves to another schema",
			sql:  "ALTER SESSION SET CURRENT_SCHEMA = other;\nRENAME t TO t3;",
		},
		{
			name: "a DROP of another namespace does not hide the table",
			sql:  "DROP INDEX t;\nRENAME t TO t2;",
			want: []finding{{1, "RENAME t TO t2", "renames table T to T2", []int{0}}},
		},
	})
}

func TestDisallowDropConstraint(t *testing.T) {
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{
			name: "constraints of the synced schema",
			sql:  "ALTER TABLE t DROP CONSTRAINT fk_p DROP CONSTRAINT uk_code;\nALTER TABLE shop.t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT nn_code;\nALTER TABLE t DROP CONSTRAINT ux_n;\nALTER TABLE t DROP CONSTRAINT unknown_c;\nALTER TABLE p DROP CONSTRAINT pk_p;",
			want: []finding{
				{0, "DROP CONSTRAINT fk_p", "drops foreign key FK_P of T", []int{0}},
				{0, "DROP CONSTRAINT uk_code", "drops unique constraint UK_CODE of T", []int{0}},
				{1, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of SHOP.T", []int{0}},
				{2, "DROP CONSTRAINT nn_code", "drops check constraint NN_CODE of T", []int{0}},
				{5, "DROP CONSTRAINT pk_p", "drops primary key PK_P of P", []int{0}},
			},
		},
		{
			name: "DROP PRIMARY KEY and DROP UNIQUE, from the SQL alone",
			sql:  "ALTER TABLE nokey DROP PRIMARY KEY;\nALTER TABLE x DROP UNIQUE (a, \"b\") CASCADE;",
			want: []finding{
				{0, "DROP PRIMARY KEY", "drops the primary key of NOKEY", []int{0}},
				{1, `DROP UNIQUE (a, "b") CASCADE`, `drops the unique constraint on (A, "b") of X`, []int{0}},
			},
			targets: []review.Target{{}},
		},
		{
			name: "constraints the change defined",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nCREATE TABLE n1 (id NUMBER CONSTRAINT pk_n1 PRIMARY KEY, r NUMBER CONSTRAINT fk_n1 REFERENCES p);\nALTER TABLE n1 DROP CONSTRAINT fk_n1;\nALTER TABLE n1 DROP CONSTRAINT pk_n1;\nALTER TABLE nokey DROP CONSTRAINT ck_a;",
			want: []finding{
				{1, "DROP CONSTRAINT ck_a", "drops check constraint CK_A of NOKEY", []int{0}},
				{3, "DROP CONSTRAINT fk_n1", "drops foreign key FK_N1 of N1", []int{0}},
				{4, "DROP CONSTRAINT pk_n1", "drops primary key PK_N1 of N1", []int{0}},
			},
		},
		{
			name:    "constraints the change defined, without a synced schema",
			sql:     "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nALTER TABLE t DROP CONSTRAINT ck_n;",
			targets: []review.Target{{}},
			want:    []finding{{1, "DROP CONSTRAINT ck_a", "drops check constraint CK_A of NOKEY", []int{0}}},
		},
		{
			name: "a table the change recreated, renamed, or dropped a column of",
			sql:  "DROP TABLE p CASCADE CONSTRAINTS;\nCREATE TABLE p (id NUMBER);\nALTER TABLE p DROP CONSTRAINT pk_p;\nRENAME t TO t2;\nALTER TABLE t2 DROP CONSTRAINT ck_n;\nALTER TABLE longt DROP COLUMN body;\nALTER TABLE longt DROP CONSTRAINT pk_longt;",
		},
		{
			name: "a name already taken fails its statement, and the scan stops judging",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE t DROP CONSTRAINT uk_code;",
		},
		{
			name: "a name the change freed first",
			sql:  "ALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_n;",
			want: []finding{
				{0, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of T", []int{0}},
				{2, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of NOKEY", []int{0}},
			},
		},
		{
			name: "a primary key dropped out of sight, and a foreign key after CASCADE",
			sql:  "ALTER TABLE t DROP PRIMARY KEY CASCADE;\nALTER TABLE t DROP CONSTRAINT pk_t;\nALTER TABLE t DROP CONSTRAINT fk_p;",
			want: []finding{{0, "DROP PRIMARY KEY CASCADE", "drops the primary key of T", []int{0}}},
		},
		{
			name: "two targets with different schemas",
			sql:  "ALTER TABLE t DROP CONSTRAINT ck_n;",
			targets: []review.Target{
				{Schema: shop()},
				{Schema: &metadata.DatabaseSchemaMetadata{Name: "SHOP"}},
				{Schema: shop()},
			},
			want: []finding{{0, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of T", []int{0, 2}}},
		},
	})
}

func TestRequirePrimaryKey(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "CREATE TABLE without a primary key",
			sql:  "CREATE TABLE a (id NUMBER);\nCREATE TABLE b (id NUMBER PRIMARY KEY);\nCREATE TABLE c (id NUMBER, CONSTRAINT pk_c PRIMARY KEY (id));\nCREATE GLOBAL TEMPORARY TABLE g (id NUMBER);\nCREATE PRIVATE TEMPORARY TABLE ora$ptt_x (id NUMBER);\nCREATE TABLE IF NOT EXISTS d (id NUMBER);\nCREATE TABLE e AS SELECT * FROM p;\nCREATE TABLE f (id PRIMARY KEY) AS SELECT id FROM p;",
			want: []finding{
				{0, "CREATE TABLE a (id NUMBER)", "creates table A without a primary key", []int{0}},
				{6, "CREATE TABLE e AS SELECT * FROM p", "creates table E without a primary key", []int{0}},
			},
		},
		{
			name: "not settled by IF NOT EXISTS",
			sql:  "CREATE TABLE a (id NUMBER);\nCREATE TABLE IF NOT EXISTS a (id NUMBER PRIMARY KEY);",
			want: []finding{{0, "CREATE TABLE a (id NUMBER)", "creates table A without a primary key", []int{0}}},
		},
		{
			name: "settled by a later key, drop, rename, or create",
			sql:  "CREATE TABLE a (id NUMBER);\nALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id);\nCREATE TABLE b (id NUMBER);\nDROP TABLE b;\nCREATE TABLE c (id NUMBER);\nRENAME c TO c2;\nCREATE TABLE d (id NUMBER);\nALTER TABLE d MODIFY (id PRIMARY KEY);\nCREATE TABLE e (id NUMBER);\nALTER TABLE e ADD (k NUMBER PRIMARY KEY);",
		},
		{
			name: "in a CREATE SCHEMA",
			sql:  "CREATE SCHEMA AUTHORIZATION shop CREATE TABLE k (id NUMBER) CREATE VIEW kv AS SELECT * FROM k;",
			want: []finding{{0, "CREATE TABLE k (id NUMBER)", "creates table K without a primary key", []int{0}}},
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE p DROP CONSTRAINT pk_p;\nALTER TABLE p ADD CONSTRAINT pk_p2 PRIMARY KEY (id, note);",
			want: []finding{{0, "DROP PRIMARY KEY", "removes the primary key of T, and the change adds none back", []int{0}}},
		},
		{
			name: "a table that had none before the change",
			sql:  "ALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey DROP PRIMARY KEY;",
		},
		{
			name: "a key the change made, then dropped",
			sql:  "CREATE TABLE a (id NUMBER CONSTRAINT pk_a PRIMARY KEY);\nALTER TABLE a DROP CONSTRAINT pk_a;\nALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t ADD PRIMARY KEY (id);\nALTER TABLE t DROP PRIMARY KEY;",
			want: []finding{
				{1, "DROP CONSTRAINT pk_a", "removes the primary key of A, and the change adds none back", []int{0}},
				{4, "DROP PRIMARY KEY", "removes the primary key of T, and the change adds none back", []int{0}},
			},
		},
		{
			name: "an external table",
			sql:  "CREATE TABLE ext (id NUMBER) ORGANIZATION EXTERNAL (TYPE ORACLE_LOADER DEFAULT DIRECTORY d LOCATION ('x.csv'));",
		},
		{
			name: "a renamed table keeps what the change did to its key",
			sql:  "ALTER TABLE nokey ADD PRIMARY KEY (a);\nRENAME nokey TO nk;\nALTER TABLE nk DROP PRIMARY KEY;",
		},
	})
}

func TestOnlineMigration(t *testing.T) {
	runCases(t, review.OnlineMigration, []ruleCase{
		{name: "not asked for", sql: "ALTER TABLE t ADD (m NUMBER);"},
		{
			name:   "asked for",
			sql:    "ALTER TABLE t ADD (m NUMBER);",
			change: review.Change{OnlineMigration: true},
			want:   []finding{{-1, "", "the change requests online migration, which Oracle does not support", []int{0}}},
		},
	})
}

func TestPriorBackup(t *testing.T) {
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	runCases(t, review.PriorBackup, []ruleCase{
		{name: "a plain backup", sql: "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM p WHERE id = 2;\nALTER TABLE nokey ADD (b NUMBER);", change: backup},
		{
			name:   "a database name with a single quote",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;",
			change: backup,
			targets: []review.Target{{Schema: func() *metadata.DatabaseSchemaMetadata {
				s := shop()
				s.Name = "S'HOP"
				return s
			}(), BackupDatabaseExists: true}},
			want: []finding{{-1, "", "prior backup writes the database name S'HOP unescaped into a quoted comment, which its single quote breaks", []int{0}}},
		},
		{
			name:   "an updatable view with a LONG column",
			sql:    "UPDATE longv SET body = NULL WHERE ROWNUM = 1;",
			change: backup,
			want:   []finding{{0, "UPDATE longv SET body = NULL WHERE ROWNUM = 1", "prior backup copies LONGV with CREATE TABLE AS SELECT, which rejects its LONG or LONG RAW column BODY (ORA-00997)", []int{0}}},
		},
		{
			name:    "no BBDATAARCHIVE",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;",
			change:  backup,
			targets: []review.Target{{Schema: shop()}},
			want:    []finding{{-1, "", "prior backup needs database BBDATAARCHIVE, which does not exist", []int{0}}},
		},
		{
			name:   "over the size limit",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;",
			change: review.Change{PriorBackup: true, MaxBackupSize: 10},
			want:   []finding{{-1, "", "the change is 32 bytes, over the 10-byte limit of prior backup", []int{0}}},
		},
		{
			name:   "nothing to back up",
			sql:    "ALTER TABLE t ADD (m NUMBER);\nINSERT INTO t (id) VALUES (3);\nMERGE INTO t USING p ON (t.id = p.id) WHEN MATCHED THEN UPDATE SET t.n = 1;\nBEGIN DELETE FROM t WHERE id = 1; END;\n/\nUPDATE (SELECT * FROM t) SET n = 1 WHERE id = 1;",
			change: backup,
			want:   []finding{{-1, "", "prior backup copies the rows of top-level UPDATE and DELETE statements on a table, and the change has none, which fails the backup", []int{0}}},
		},
		{
			name:   "a database link, and a quote in a name",
			sql:    "UPDATE t@lnk SET n = 1 WHERE id = 1;\nDELETE FROM \"a'b\" WHERE id = 1;\nDELETE FROM \"a'b\" WHERE id = 2;",
			change: backup,
			want: []finding{
				{0, "UPDATE t@lnk SET n = 1 WHERE id = 1", "prior backup does not support table T@LNK, referenced through a database link", []int{0}},
				{1, `DELETE FROM "a'b" WHERE id = 1`, `prior backup writes the name of table "a'b" unescaped into a quoted comment, which its single quote breaks`, []int{0}},
			},
		},
		{
			name:   "UPDATE and DELETE on one table",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM shop.t WHERE id = 2;\nDELETE FROM other.p WHERE id = 1;\nUPDATE p SET note = 'x' WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "DELETE FROM shop.t WHERE id = 2", "prior backup cannot back up both UPDATE and DELETE on SHOP.T", []int{0}}},
		},
		{
			name:    "UPDATE and DELETE on one table, without a synced schema",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM shop.t WHERE id = 2;\nDELETE FROM t WHERE id = 3;",
			change:  backup,
			targets: []review.Target{{BackupDatabaseExists: true}},
			want:    []finding{{2, "DELETE FROM t WHERE id = 3", "prior backup cannot back up both UPDATE and DELETE on T", []int{0}}},
		},
		{
			name:   "a LONG column",
			sql:    "DELETE FROM longt WHERE id = 1;\nDELETE FROM longt WHERE id = 2;",
			change: backup,
			want:   []finding{{0, "DELETE FROM longt WHERE id = 1", "prior backup copies LONGT with CREATE TABLE AS SELECT, which rejects its LONG or LONG RAW column BODY (ORA-00997)", []int{0}}},
		},
		{
			name:   "a table the change makes, named with its schema",
			sql:    "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE shop.n2 SET id = 2 WHERE id = 1;\nDELETE FROM t WHERE id IN (SELECT id FROM shop.n2);\nUPDATE n2 SET id = 3 WHERE id = 2;",
			change: backup,
			want: []finding{
				{1, "UPDATE shop.n2 SET id = 2 WHERE id = 1", "prior backup runs before the change, when SHOP.N2 does not exist yet", []int{0}},
				{2, "DELETE FROM t WHERE id IN (SELECT id FROM shop.n2)", "prior backup runs before the change, when SHOP.N2 does not exist yet", []int{0}},
				{3, "UPDATE n2 SET id = 3 WHERE id = 2", "prior backup runs before the change, when N2 does not exist yet", []int{0}},
			},
		},
		{
			name:   "a table the change makes in place of a synonym it drops",
			sql:    "DROP SYNONYM n2;\nCREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE shop.n2 SET id = 2 WHERE id = 1;",
			change: backup,
		},
		{
			name:   "a table the change makes, unqualified in a lone statement, which a public synonym may resolve",
			sql:    "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;",
			change: backup,
		},
		{
			name:   "a column the change makes",
			sql:    "ALTER TABLE nokey ADD (gone NUMBER);\nDELETE FROM nokey WHERE gone = 1;\nALTER TABLE p RENAME COLUMN note TO memo;\nUPDATE p x SET x.id = 1 WHERE x.memo = 'a';\nUPDATE t SET n = 1 WHERE n IN (SELECT gone FROM nokey);\nALTER TABLE t ADD (gone2 NUMBER);\nUPDATE t SET gone2 = 1 WHERE id = 1;",
			change: backup,
			want: []finding{
				{1, "DELETE FROM nokey WHERE gone = 1", "prior backup runs before the change, when column GONE of NOKEY does not exist yet", []int{0}},
				{3, "UPDATE p x SET x.id = 1 WHERE x.memo = 'a'", "prior backup runs before the change, when column MEMO of P does not exist yet", []int{0}},
			},
		},
		{
			name:   "a column of a joined table in an UPDATE's FROM list",
			sql:    "ALTER TABLE p ADD (memo NUMBER);\nUPDATE t SET n = 1 FROM p WHERE t.p_id = p.id AND p.memo = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE t SET n = 1 FROM p WHERE t.p_id = p.id AND p.memo = 1", "prior backup runs before the change, when column MEMO of P does not exist yet", []int{0}}},
		},
		{
			name:   "statements after the session moves to another schema",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nALTER SESSION SET CURRENT_SCHEMA = other;\nALTER TABLE nokey ADD (gone NUMBER);\nDELETE FROM nokey WHERE gone = 1;",
			change: backup,
		},
	})
}

func TestNoTargets(t *testing.T) {
	sql := "DROP TABLE t;\nALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nCREATE TABLE a (id NUMBER);\nUPDATE t SET n = 1 WHERE id = 1;"
	result, err := Review(context.Background(), sql, review.Options{Change: review.Change{PriorBackup: true}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range result.Findings {
		if len(f.Targets) != 0 {
			t.Errorf("%s: Targets = %v, want none", f.Rule, f.Targets)
		}
		got = append(got, string(f.Rule)+": "+f.Message)
	}
	want := []string{
		"DISALLOW_DROP_OBJECT: drops table T",
		"DISALLOW_DROP_CONSTRAINT: drops check constraint CK_A of NOKEY",
		"REQUIRE_PRIMARY_KEY: creates table A without a primary key",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

// TestConstraintScanLosesTrack pins the statements after which the scan
// stops judging DROP CONSTRAINT, each because what the constraints are is
// no longer certain.
func TestConstraintScanLosesTrack(t *testing.T) {
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{
			name: "MODIFY ... NULL drops a NOT NULL constraint",
			sql:  "ALTER TABLE t MODIFY (code NULL);\nALTER TABLE t DROP CONSTRAINT nn_code;",
		},
		{
			name: "RENAME CONSTRAINT",
			sql:  "ALTER TABLE t RENAME CONSTRAINT ck_n TO ck_m;\nALTER TABLE t DROP CONSTRAINT ck_m;\nALTER TABLE t DROP CONSTRAINT ck_n;",
		},
		{
			name: "a renamed table keeps the names of its constraints",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT c1 CHECK (a > 0);\nRENAME nokey TO nk;\nCREATE TABLE nokey (a NUMBER CONSTRAINT c1 CHECK (a > 0));\nALTER TABLE nokey DROP CONSTRAINT c1;",
		},
		{
			name: "IF NOT EXISTS of a table the change made keeps what it defined",
			sql:  "CREATE TABLE a (id NUMBER, CONSTRAINT ck_a CHECK (id > 0));\nCREATE TABLE IF NOT EXISTS a (id NUMBER);\nALTER TABLE a DROP CONSTRAINT ck_a;",
			want: []finding{{2, "DROP CONSTRAINT ck_a", "drops check constraint CK_A of A", []int{0}}},
		},
		{
			name: "FLASHBACK TABLE to a time, and DROP USER without CASCADE, keep the scan",
			sql:  "FLASHBACK TABLE t TO TIMESTAMP SYSTIMESTAMP;\nDROP USER missing;\nALTER TABLE t DROP CONSTRAINT ck_n DROP CONSTRAINT fk_p;",
			want: []finding{
				{2, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of T", []int{0}},
				{2, "DROP CONSTRAINT fk_p", "drops foreign key FK_P of T", []int{0}},
			},
		},
		{
			name: "FLASHBACK TABLE restores a table",
			sql:  "FLASHBACK TABLE t TO BEFORE DROP;\nALTER TABLE t DROP CONSTRAINT ck_n;",
		},
		{
			name: "DROP CLUSTER may drop tables it does not name",
			sql:  "DROP CLUSTER cl INCLUDING TABLES;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE p DROP PRIMARY KEY;",
			want: []finding{{2, "DROP PRIMARY KEY", "drops the primary key of P", []int{0}}},
		},
		{
			name: "CREATE TABLE IF NOT EXISTS may leave the table there was",
			sql:  "CREATE TABLE IF NOT EXISTS n1 (id NUMBER CONSTRAINT ck_x CHECK (id > 0));\nALTER TABLE n1 DROP CONSTRAINT ck_x;\nALTER TABLE nokey ADD CONSTRAINT ck_x CHECK (a > 0);\nALTER TABLE t DROP CONSTRAINT ck_n;",
		},
		{
			name: "the session moves to another schema",
			sql:  "ALTER SESSION SET CURRENT_SCHEMA = other;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE shop.t DROP CONSTRAINT uk_code;",
			want: []finding{{2, "DROP CONSTRAINT uk_code", "drops unique constraint UK_CODE of SHOP.T", []int{0}}},
		},
		{
			name:    "a name the change qualifies, the database not known",
			sql:     "ALTER TABLE shop.nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nALTER TABLE x ADD CONSTRAINT ck_b CHECK (a > 0);\nDROP TABLE shop.x;\nALTER TABLE x DROP CONSTRAINT ck_b;",
			targets: []review.Target{{}},
		},
	})
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name:    "a qualified drop settles an unqualified table, the database not known",
			sql:     "CREATE TABLE a (id NUMBER);\nDROP TABLE shop.a;",
			targets: []review.Target{{}},
		},
		{
			name: "DROP USER drops the tables of its schema",
			sql:  "CREATE TABLE other.a (id NUMBER);\nDROP USER other CASCADE;",
		},
	})
}

func TestTargetsMerge(t *testing.T) {
	sql := "DELETE FROM longt WHERE id = 1;\nALTER TABLE t DROP CONSTRAINT ck_n;"
	other := &metadata.DatabaseSchemaMetadata{Name: "SHOP", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "LONGT", Columns: columns("ID")}}}}}
	result, err := Review(context.Background(), sql, review.Options{Change: review.Change{PriorBackup: true}}, []review.Target{
		{Schema: shop(), BackupDatabaseExists: true},
		{Schema: other},
		{Schema: shop(), BackupDatabaseExists: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	var got []finding
	for _, f := range result.Findings {
		got = append(got, finding{f.Statement, sql[f.Range.Start:f.Range.End], f.Message, f.Targets})
	}
	want := []finding{
		{-1, "", "prior backup needs database BBDATAARCHIVE, which does not exist", []int{1}},
		{0, "DELETE FROM longt WHERE id = 1", "prior backup copies LONGT with CREATE TABLE AS SELECT, which rejects its LONG or LONG RAW column BODY (ORA-00997)", []int{0, 2}},
		{1, "DROP CONSTRAINT ck_n", "drops check constraint CK_N of T", []int{0, 2}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

// TestNoPanic feeds every prefix of a change through Review: none may
// panic, and each either parses or ends in one Syntax finding.
func TestNoPanic(t *testing.T) {
	sql := "CREATE TABLE a (id NUMBER CONSTRAINT pk_a PRIMARY KEY, n NUMBER CHECK (n <> ''));\nALTER TABLE a DROP (n) DROP PRIMARY KEY;\nUPDATE t x SET x.n = 1 FROM p WHERE x.p_id = p.id AND p.note = NULL;\nBEGIN IF 1 = NULL THEN DELETE FROM t; END IF; END;\n/\nRENAME t TO u;\nALTER TABLE p TRUNCATE PARTITION FOR (1), p2;"
	opts := review.Options{Change: review.Change{PriorBackup: true, OnlineMigration: true}}
	targets := []review.Target{{Schema: shop()}, {}}
	for i := 0; i <= len(sql); i++ {
		result, err := Review(context.Background(), sql[:i], opts, targets)
		if err != nil {
			t.Fatalf("prefix %d: %v", i, err)
		}
		syntax := 0
		for _, f := range result.Findings {
			if f.Rule == review.Syntax {
				syntax++
			}
		}
		if syntax > 0 && len(result.Findings) != 1 {
			t.Fatalf("prefix %d: a syntax error with other findings: %v", i, result.Findings)
		}
	}
}
