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
				targets = []review.Target{{Schema: shop()}}
			}
			got := run(t, rule, tt.sql, tt.change, targets)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// shop is a small synced database, its names stored as Snowflake stores
// them: PUBLIC.T with a primary key, a unique constraint, a plain index, a
// foreign key to PUBLIC.P, and a check; PUBLIC.P with a primary key;
// PUBLIC."lower", created with a quoted name; PUBLIC.V, a view; and DUP, a
// table in both PUBLIC and S.
func shop() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "SHOP",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "PUBLIC",
				Tables: []*metadata.TableMetadata{
					{
						Name:    "P",
						Indexes: []*metadata.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true}},
					},
					{
						Name: "T",
						Indexes: []*metadata.IndexMetadata{
							{Name: "PK_T", Primary: true, Unique: true},
							{Name: "UK_CODE", Unique: true},
							{Name: "IX_N"},
						},
						ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "FK_P", ReferencedTable: "P"}},
						CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "CK_N"}},
					},
					{
						Name:    "lower",
						Indexes: []*metadata.IndexMetadata{{Name: "pk_lower", Primary: true, Unique: true}},
					},
					{
						Name:    "DUP",
						Indexes: []*metadata.IndexMetadata{{Name: "PK_DUP", Primary: true, Unique: true}},
					},
				},
				Views: []*metadata.ViewMetadata{{Name: "V"}},
			},
			{
				Name: "S",
				Tables: []*metadata.TableMetadata{
					{
						Name:    "DUP",
						Indexes: []*metadata.IndexMetadata{{Name: "PK_DUP", Primary: true, Unique: true}},
					},
				},
			},
		},
	}
}

func TestStatements(t *testing.T) {
	sql := "  -- lead\nSELECT 1 /* trail */ ;\n;\nBEGIN\n  DELETE FROM t WHERE id = 1;\nEND;\n// note\nUPDATE t SET n = 1 WHERE id = 1"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var got []string
	for _, rng := range result.Statements {
		got = append(got, sql[rng.Start:rng.End])
	}
	want := []string{"SELECT 1", "BEGIN\n  DELETE FROM t WHERE id = 1;\nEND", "UPDATE t SET n = 1 WHERE id = 1"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("statements mismatch (-want +got):\n%s", diff)
	}
	if len(result.Findings) != 0 {
		t.Errorf("Findings = %v", result.Findings)
	}
}

func TestSyntax(t *testing.T) {
	runCases(t, review.Syntax, []ruleCase{
		{name: "valid", sql: "SELECT 1; UPDATE t SET n = 1 WHERE id = 1"},
		{
			name: "one finding for the first error",
			sql:  "SELECT 1;\nSELEC 2;\nSELEC 3;",
			want: []finding{{1, "", "unknown or unsupported statement starting with SELEC", []int{0}}},
		},
		{
			name: "a lex error after the last statement",
			sql:  "SELECT 1; /* open",
			want: []finding{{0, "", "unterminated block comment", []int{0}}},
		},
		{
			name: "a lex error alone",
			sql:  "/* open",
			want: []finding{{-1, "", "unterminated block comment", []int{0}}},
		},
	})

	// A syntax error ends the review: no other rule reports, and every
	// statement is still listed.
	sql := "DROP TABLE t;\nUPDATE t SET n = 1;\nSELEC 3;"
	result, err := Review(context.Background(), sql, review.Options{}, []review.Target{{Schema: shop()}, {}})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if len(result.Statements) != 3 {
		t.Errorf("Statements = %v, want 3", result.Statements)
	}
	want := []review.Finding{{Rule: review.Syntax, Statement: 2, Range: review.Range{Start: 34, End: 34}, Message: "unknown or unsupported statement starting with SELEC", Targets: []int{0, 1}}}
	if diff := cmp.Diff(want, result.Findings); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

func TestOptions(t *testing.T) {
	sql := "DROP TABLE t;\nUPDATE t SET n = 1;"
	rulesOf := func(opts review.Options) []review.Rule {
		result, err := Review(context.Background(), sql, opts, []review.Target{{Schema: shop()}})
		if err != nil {
			t.Fatalf("Review: %v", err)
		}
		var rules []review.Rule
		for _, f := range result.Findings {
			rules = append(rules, f.Rule)
		}
		return rules
	}
	if diff := cmp.Diff([]review.Rule{review.DisallowDropObject, review.RequireWhere}, rulesOf(review.Options{})); diff != "" {
		t.Errorf("every rule (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]review.Rule{review.RequireWhere}, rulesOf(review.Options{Rules: []review.Rule{review.WalkThrough, review.RequireWhere, "NO_SUCH_RULE"}})); diff != "" {
		t.Errorf("one rule (-want +got):\n%s", diff)
	}
	if got := rulesOf(review.Options{Rules: []review.Rule{}}); len(got) != 0 {
		t.Errorf("no rules: %v", got)
	}
}

func TestRequireWhere(t *testing.T) {
	runCases(t, review.RequireWhere, []ruleCase{
		{
			name: "UPDATE and DELETE without WHERE",
			sql:  "UPDATE t SET n = 1;\nDELETE FROM public.t;\nUPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t USING p WHERE t.p_id = p.id;\nUPDATE t SET n = p.n FROM p;",
			want: []finding{
				{0, "UPDATE t SET n = 1", "UPDATE t has no WHERE clause", []int{0}},
				{1, "DELETE FROM public.t", "DELETE FROM public.t has no WHERE clause", []int{0}},
				{4, "UPDATE t SET n = p.n FROM p", "UPDATE t has no WHERE clause", []int{0}},
			},
		},
		{
			name: "inside a block, a task, and an alert",
			sql:  "BEGIN\n  IF (1 = 1) THEN\n    DELETE FROM t;\n  END IF;\nEND;\nCREATE TASK tk WAREHOUSE = w SCHEDULE = '1 minute' AS UPDATE t SET n = 1;\nCREATE ALERT al WAREHOUSE = w SCHEDULE = '1 minute' IF (EXISTS (SELECT 1)) THEN DELETE FROM \"Log\";",
			want: []finding{
				{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}},
				{1, "UPDATE t SET n = 1", "UPDATE t has no WHERE clause", []int{0}},
				{2, `DELETE FROM "Log"`, `DELETE FROM "Log" has no WHERE clause`, []int{0}},
			},
		},
		{
			name: "a routine body and EXECUTE IMMEDIATE are text",
			sql:  "CREATE PROCEDURE pr() RETURNS INT LANGUAGE SQL AS $$ BEGIN DELETE FROM t; RETURN 1; END $$;\nEXECUTE IMMEDIATE 'DELETE FROM t';",
		},
		{
			name: "a statement under EXPLAIN, which does not run it",
			sql:  "EXPLAIN UPDATE t SET n = 1;\nEXPLAIN DELETE FROM t;",
		},
	})
}

func TestRequireIsNull(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "comparisons a predicate tests",
			sql:  "SELECT * FROM t WHERE n = NULL OR NULL <> code;\nUPDATE t SET n = 1 WHERE NOT (code != NULL);\nDELETE FROM t WHERE (id = 1 AND (n = (NULL)));",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "NULL <> code", `"NULL <> code" is never true; use IS NOT NULL`, []int{0}},
				{1, "code != NULL", `"code != NULL" is never true; use IS NOT NULL`, []int{0}},
				{2, "n = (NULL)", `"n = (NULL)" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "HAVING, QUALIFY, join conditions, and hierarchical queries",
			sql:  "SELECT p.id FROM t JOIN p ON t.p_id = NULL GROUP BY p.id HAVING MAX(t.n) = NULL QUALIFY ROW_NUMBER() OVER (ORDER BY p.id) = NULL;\nSELECT id FROM t START WITH p_id = NULL CONNECT BY PRIOR id = p_id;",
			want: []finding{
				{0, "t.p_id = NULL", `"t.p_id = NULL" is never true; use IS NULL`, []int{0}},
				{0, "MAX(t.n) = NULL", `"MAX(t.n) = NULL" is never true; use IS NULL`, []int{0}},
				{0, "ROW_NUMBER() OVER (ORDER BY p.id) = NULL", `"ROW_NUMBER() OVER (ORDER BY p.id) = NULL" is never true; use IS NULL`, []int{0}},
				{1, "p_id = NULL", `"p_id = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "IFF and a searched CASE, not a simple CASE",
			sql:  "SELECT IFF(n = NULL, 1, 0), CASE WHEN code <> NULL THEN 1 END, CASE n WHEN NULL THEN 1 END FROM t;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "code <> NULL", `"code <> NULL" is never true; use IS NOT NULL`, []int{0}},
			},
		},
		{
			name: "MERGE and INSERT ALL conditions",
			sql:  "MERGE INTO t USING p ON t.id = NULL WHEN MATCHED AND p.note <> NULL THEN DELETE;\nINSERT ALL WHEN id = NULL THEN INTO t SELECT id FROM p;",
			want: []finding{
				{0, "t.id = NULL", `"t.id = NULL" is never true; use IS NULL`, []int{0}},
				{0, "p.note <> NULL", `"p.note <> NULL" is never true; use IS NOT NULL`, []int{0}},
				{1, "id = NULL", `"id = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "Snowflake Scripting conditions",
			sql:  "DECLARE c INT DEFAULT 1;\nBEGIN\n  IF (c = NULL) THEN\n    RETURN 0;\n  ELSEIF (c <> NULL) THEN\n    RETURN 1;\n  END IF;\n  WHILE (c != NULL) DO\n    c := c + 1;\n  END WHILE;\n  RETURN c;\nEND",
			want: []finding{
				{0, "c = NULL", `"c = NULL" is never true; use IS NULL`, []int{0}},
				{0, "c <> NULL", `"c <> NULL" is never true; use IS NOT NULL`, []int{0}},
				{0, "c != NULL", `"c != NULL" is never true; use IS NOT NULL`, []int{0}},
			},
		},
		{
			name: "a value, IS NULL, a cast, and a comparison inside another",
			sql:  "SELECT n = NULL FROM t WHERE n IS NULL AND code IS NOT NULL;\nSELECT * FROM t WHERE n = NULL::INT;\nSELECT * FROM t WHERE (n = NULL) = FALSE;",
		},
		{
			name: "a task's WHEN, MATCH_RECOGNIZE DEFINE, and quoted text kept as written",
			sql:  "CREATE TASK tk WAREHOUSE = w WHEN SYSTEM$STREAM_HAS_DATA('s') = NULL AS DELETE FROM t WHERE id = 1;\nSELECT * FROM t MATCH_RECOGNIZE (PATTERN (a) DEFINE a AS a.n = NULL);\nSELECT * FROM t WHERE \"a  b\" = NULL;",
			want: []finding{
				{0, "SYSTEM$STREAM_HAS_DATA('s') = NULL", `"SYSTEM$STREAM_HAS_DATA('s') = NULL" is never true; use IS NULL`, []int{0}},
				{1, "a.n = NULL", `"a.n = NULL" is never true; use IS NULL`, []int{0}},
				{2, "\"a  b\" = NULL", `""a  b" = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "a backslash-escaped quote in a string",
			sql:  "SELECT * FROM t WHERE 'a\\'  b' = NULL;",
			want: []finding{{0, "'a\\'  b' = NULL", "\"'a\\'  b' = NULL\" is never true; use IS NULL", []int{0}}},
		},
		{
			name: "a row access policy's body and a dynamic table's IMMUTABLE WHERE",
			sql:  "CREATE ROW ACCESS POLICY rp AS (n INT) RETURNS BOOLEAN -> n = NULL;\nCREATE DYNAMIC TABLE dt TARGET_LAG = '1 minute' WAREHOUSE = w IMMUTABLE WHERE (n = NULL) AS SELECT n FROM t;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{1, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "RECLUSTER WHERE",
			sql:  "ALTER TABLE t RECLUSTER WHERE n = NULL;",
			want: []finding{{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{name: "a statement under EXPLAIN, which does not run it", sql: "EXPLAIN DROP TABLE t;"},
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t;\nDROP VIEW v;\nDROP MATERIALIZED VIEW mv;\nDROP DYNAMIC TABLE dt;\nDROP EXTERNAL TABLE et;\nDROP STAGE st;\nDROP STREAM sm;\nDROP TASK tk;\nDROP PIPE pp;\nDROP SEQUENCE sq;\nDROP FUNCTION f(NUMBER);\nDROP PROCEDURE pr();\nDROP USER u;\nDROP ROLE r;\nDROP WAREHOUSE w;\nDROP SCHEMA s CASCADE;\nDROP DATABASE d;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t", "drops table t", []int{0}},
				{1, "DROP VIEW v", "drops view v", []int{0}},
				{2, "DROP MATERIALIZED VIEW mv", "drops materialized view mv", []int{0}},
				{3, "DROP DYNAMIC TABLE dt", "drops dynamic table dt", []int{0}},
				{4, "DROP EXTERNAL TABLE et", "drops external table et", []int{0}},
				{5, "DROP STAGE st", "drops stage st", []int{0}},
				{6, "DROP STREAM sm", "drops stream sm", []int{0}},
				{7, "DROP TASK tk", "drops task tk", []int{0}},
				{8, "DROP PIPE pp", "drops pipe pp", []int{0}},
				{9, "DROP SEQUENCE sq", "drops sequence sq", []int{0}},
				{10, "DROP FUNCTION f(NUMBER)", "drops function f(NUMBER)", []int{0}},
				{11, "DROP PROCEDURE pr()", "drops procedure pr()", []int{0}},
				{12, "DROP USER u", "drops user u", []int{0}},
				{13, "DROP ROLE r", "drops role r", []int{0}},
				{14, "DROP WAREHOUSE w", "drops warehouse w", []int{0}},
				{15, "DROP SCHEMA s CASCADE", "drops schema s", []int{0}},
				{16, "DROP DATABASE d", "drops database d", []int{0}},
			},
		},
		{
			name: "ALTER TABLE ... DROP COLUMN anchors each column",
			sql:  "ALTER TABLE shop.public.t DROP COLUMN \"Code\", n;",
			want: []finding{
				{0, `"Code"`, `drops column "Code" of shop.public.t`, []int{0}},
				{0, "n", "drops column n of shop.public.t", []int{0}},
			},
		},
		{
			// The parser starts DROP SCHEMA's location at SCHEMA; at top
			// level the finding anchors the statement instead.
			name: "inside a block",
			sql:  "BEGIN\n  DROP TABLE t;\n  DROP SCHEMA s;\nEND",
			want: []finding{
				{0, "DROP TABLE t", "drops table t", []int{0}},
				{0, "SCHEMA s", "drops schema s", []int{0}},
			},
		},
		{
			name: "a quoted name",
			sql:  `DROP TABLE "My ""T""";`,
			want: []finding{{0, `DROP TABLE "My ""T"""`, `drops table "My ""T"""`, []int{0}}},
		},
		{
			name: "UNDROP, constraints, and properties are not object drops",
			sql:  "UNDROP TABLE t;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t DROP CLUSTERING KEY;\nALTER TABLE t DROP SEARCH OPTIMIZATION;\nALTER TABLE t DROP ROW ACCESS POLICY rap;",
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{name: "a statement under EXPLAIN, which does not run it", sql: "EXPLAIN TRUNCATE TABLE t;"},
		{
			name: "every form",
			sql:  "TRUNCATE TABLE t;\nTRUNCATE IF EXISTS public.p;\nTRUNCATE MATERIALIZED VIEW mv;\nTRUNCATE ERROR_TABLE(t);\nBEGIN\n  TRUNCATE t;\nEND",
			want: []finding{
				{0, "TRUNCATE TABLE t", "truncates t", []int{0}},
				{1, "TRUNCATE IF EXISTS public.p", "truncates public.p", []int{0}},
				{2, "TRUNCATE MATERIALIZED VIEW mv", "truncates materialized view mv", []int{0}},
				{3, "TRUNCATE ERROR_TABLE(t)", "truncates the error table of t", []int{0}},
				{4, "TRUNCATE t", "truncates t", []int{0}},
			},
		},
		{name: "DELETE is not TRUNCATE", sql: "DELETE FROM t WHERE id = 1;"},
	})
}

func TestDisallowRename(t *testing.T) {
	runCases(t, review.DisallowRename, []ruleCase{
		{name: "a statement under EXPLAIN, which does not run it", sql: "EXPLAIN ALTER TABLE t RENAME TO u;"},
		{
			name: "tables, views, and columns",
			sql:  "ALTER TABLE t RENAME TO s.t3;\nALTER TABLE IF EXISTS t SWAP WITH p;\nALTER TABLE t RENAME COLUMN n TO \"N2\";\nALTER VIEW v RENAME TO v2;\nALTER MATERIALIZED VIEW mv RENAME TO mv2;\nALTER DYNAMIC TABLE dt RENAME TO dt2;\nALTER DYNAMIC TABLE dt SWAP WITH dt3;\nALTER SEMANTIC VIEW sv RENAME TO sv2;",
			want: []finding{
				{0, "RENAME TO s.t3", "renames table t to s.t3", []int{0}},
				{1, "SWAP WITH p", "swaps the names of tables t and p", []int{0}},
				{2, `RENAME COLUMN n TO "N2"`, `renames column n of t to "N2"`, []int{0}},
				{3, "ALTER VIEW v RENAME TO v2", "renames view v to v2", []int{0}},
				{4, "ALTER MATERIALIZED VIEW mv RENAME TO mv2", "renames materialized view mv to mv2", []int{0}},
				{5, "ALTER DYNAMIC TABLE dt RENAME TO dt2", "renames dynamic table dt to dt2", []int{0}},
				{6, "ALTER DYNAMIC TABLE dt SWAP WITH dt3", "swaps the names of dynamic tables dt and dt3", []int{0}},
				{7, "ALTER SEMANTIC VIEW sv RENAME TO sv2", "renames semantic view sv to sv2", []int{0}},
			},
		},
		{
			name: "inside a block",
			sql:  "BEGIN\n  ALTER TABLE t RENAME TO t4;\nEND",
			want: []finding{{0, "RENAME TO t4", "renames table t to t4", []int{0}}},
		},
		{
			name: "other objects are not in scope",
			sql:  "ALTER DATABASE d RENAME TO d2;\nALTER SCHEMA s RENAME TO s2;\nALTER SEQUENCE sq RENAME TO sq2;\nALTER STAGE st RENAME TO st2;\nALTER TABLE t RENAME CONSTRAINT ck_n TO ck_m;",
		},
	})
}

func TestDisallowDropConstraint(t *testing.T) {
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{name: "a statement under EXPLAIN, which does not run it", sql: "EXPLAIN ALTER TABLE t DROP PRIMARY KEY;"},
		{
			name: "from the SQL alone",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t DROP UNIQUE (code);\nALTER TABLE t DROP FOREIGN KEY (p_id) CASCADE;",
			want: []finding{
				{0, "DROP PRIMARY KEY", "drops the primary key of t", []int{0}},
				{1, "DROP UNIQUE (code)", "drops a unique constraint of t", []int{0}},
				{2, "DROP FOREIGN KEY (p_id) CASCADE", "drops a foreign key of t", []int{0}},
			},
		},
		{
			name: "constraints of the synced schema",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t, DROP CONSTRAINT fk_p;\nALTER TABLE public.t DROP CONSTRAINT \"UK_CODE\";\nALTER TABLE shop.public.t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT ix_n;\nALTER TABLE t DROP CONSTRAINT nope;",
			want: []finding{
				{0, "DROP CONSTRAINT pk_t", "drops primary key pk_t of t", []int{0}},
				{0, "DROP CONSTRAINT fk_p", "drops foreign key fk_p of t", []int{0}},
				{1, `DROP CONSTRAINT "UK_CODE"`, `drops unique constraint "UK_CODE" of public.t`, []int{0}},
				{2, "DROP CONSTRAINT ck_n", "drops check constraint ck_n of shop.public.t", []int{0}},
			},
		},
		{
			name: "names fold to upper case unless quoted",
			sql:  "ALTER TABLE \"lower\" DROP CONSTRAINT \"pk_lower\";\nALTER TABLE lower DROP CONSTRAINT pk_lower;\nALTER TABLE \"t\" DROP CONSTRAINT pk_t;\nALTER TABLE t DROP CONSTRAINT \"pk_t\";",
			want: []finding{{0, `DROP CONSTRAINT "pk_lower"`, `drops primary key "pk_lower" of "lower"`, []int{0}}},
		},
		{
			name: "a table the synced schema does not settle",
			sql:  "ALTER TABLE dup DROP CONSTRAINT pk_dup;\nALTER TABLE other.public.t DROP CONSTRAINT pk_t;\nALTER TABLE s.dup DROP CONSTRAINT pk_dup;",
			want: []finding{{2, "DROP CONSTRAINT pk_dup", "drops primary key pk_dup of s.dup", []int{0}}},
		},
		{
			name: "constraints the change defined",
			sql:  "CREATE TABLE n1 (a INT CONSTRAINT pk_n1 PRIMARY KEY, b INT, CONSTRAINT uk_b UNIQUE (b));\nALTER TABLE n1 DROP CONSTRAINT uk_b;\nALTER TABLE p ADD CONSTRAINT fk_t FOREIGN KEY (id) REFERENCES t (id);\nALTER TABLE p DROP CONSTRAINT fk_t;\nALTER TABLE n1 RENAME CONSTRAINT pk_n1 TO pk_n2;\nALTER TABLE n1 DROP CONSTRAINT pk_n2;",
			want: []finding{
				{1, "DROP CONSTRAINT uk_b", "drops unique constraint uk_b of n1", []int{0}},
				{3, "DROP CONSTRAINT fk_t", "drops foreign key fk_t of p", []int{0}},
				{5, "DROP CONSTRAINT pk_n2", "drops primary key pk_n2 of n1", []int{0}},
			},
		},
		{
			name: "a table the change replaced, or a name after USE",
			sql:  "CREATE OR REPLACE TABLE t (id INT);\nALTER TABLE t DROP CONSTRAINT uk_code;\nALTER TABLE p RENAME TO p2;\nALTER TABLE p2 DROP CONSTRAINT pk_p;\nUSE SCHEMA s;\nALTER TABLE dup DROP CONSTRAINT pk_dup;",
		},
		{
			name: "a constraint defined after USE is the change's",
			sql:  "USE DATABASE other;\nCREATE TABLE n1 (a INT, CONSTRAINT uk_a UNIQUE (a));\nALTER TABLE n1 DROP CONSTRAINT uk_a;\nUSE ROLE admin;\nALTER TABLE n1 DROP CONSTRAINT uk_a;",
			want: []finding{{2, "DROP CONSTRAINT uk_a", "drops unique constraint uk_a of n1", []int{0}}},
		},
		{
			name: "a constraint defined inside a block is not known",
			sql:  "BEGIN\n  ALTER TABLE p ADD CONSTRAINT uk_note UNIQUE (note);\nEND;\nALTER TABLE p DROP CONSTRAINT uk_note;\nBEGIN\n  ALTER TABLE t DROP CONSTRAINT fk_p;\nEND",
			want: []finding{{2, "DROP CONSTRAINT fk_p", "drops foreign key fk_p of t", []int{0}}},
		},
		{
			name:    "per target",
			sql:     "ALTER TABLE t DROP CONSTRAINT pk_t;\nALTER TABLE t DROP PRIMARY KEY;",
			targets: []review.Target{{Schema: shop()}, {}, {Schema: shop()}},
			want: []finding{
				{0, "DROP CONSTRAINT pk_t", "drops primary key pk_t of t", []int{0, 2}},
				{1, "DROP PRIMARY KEY", "drops the primary key of t", []int{0, 1, 2}},
			},
		},
	})
}

func TestRequirePrimaryKey(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{name: "a statement under EXPLAIN, which does not run it", sql: "EXPLAIN CREATE TABLE x (id INT);"},
		{
			name: "CREATE TABLE without a primary key",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT, CONSTRAINT pk_c PRIMARY KEY (id));\nCREATE OR REPLACE TRANSIENT TABLE d (id INT UNIQUE);\nCREATE HYBRID TABLE e (id INT, CONSTRAINT pk_e PRIMARY KEY (id));",
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}},
				{3, "CREATE OR REPLACE TRANSIENT TABLE d (id INT UNIQUE)", "creates table d without a primary key", []int{0}},
			},
		},
		{
			name: "CTAS is judged by what it declares",
			sql:  "CREATE TABLE a AS SELECT id FROM t;\nCREATE TABLE b (id INT PRIMARY KEY) AS SELECT id FROM t;",
			want: []finding{{0, "CREATE TABLE a AS SELECT id FROM t", "creates table a without a primary key", []int{0}}},
		},
		{
			name: "tables not judged",
			sql:  "CREATE TEMPORARY TABLE a (id INT);\nCREATE LOCAL TEMPORARY TABLE b (id INT);\nCREATE VOLATILE TABLE c (id INT);\nCREATE TABLE IF NOT EXISTS d (id INT);\nCREATE TABLE e LIKE t;\nCREATE TABLE f CLONE t;\nCREATE OR ALTER TABLE g (id INT);\nBEGIN\n  CREATE TABLE h (id INT);\nEND",
		},
		{
			name: "settled by a later key, drop, rename, swap, or create",
			sql:  "CREATE TABLE a (id INT NOT NULL);\nALTER TABLE a ADD PRIMARY KEY (id);\nCREATE TABLE b (id INT);\nDROP TABLE b;\nCREATE TABLE c (id INT);\nALTER TABLE c RENAME TO c2;\nCREATE TABLE d (id INT);\nALTER TABLE p SWAP WITH d;\nCREATE TABLE e (id INT);\nCREATE OR REPLACE TABLE e (id INT PRIMARY KEY);\nCREATE TABLE s2.f (id INT);\nDROP SCHEMA s2;\nCREATE TABLE g (id INT);\nALTER TABLE g ADD COLUMN k INT PRIMARY KEY;",
		},
		{
			name: "settled by a column added PRIMARY KEY with another inline constraint",
			sql:  "CREATE TABLE a (x INT);\nALTER TABLE a ADD COLUMN id INT PRIMARY KEY UNIQUE;",
		},
		{
			name: "not settled by a keyless CREATE OR ALTER, nor judged under IF EXISTS",
			sql:  "CREATE TABLE a (id INT);\nCREATE OR ALTER TABLE a (id INT);\nALTER TABLE IF EXISTS missing DROP PRIMARY KEY;",
			want: []finding{{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}}},
		},
		{
			name: "not settled by a schema of another database",
			sql:  "CREATE TABLE db1.s.a (id INT);\nDROP SCHEMA db2.s;",
			want: []finding{{0, "CREATE TABLE db1.s.a (id INT)", "creates table db1.s.a without a primary key", []int{0}}},
		},
		{
			name: "a column declared PRIMARY KEY with another inline constraint",
			sql:  "CREATE TABLE a (id INT PRIMARY KEY UNIQUE NOT NULL);",
		},
		{
			name: "not settled by IF NOT EXISTS or a temporary table of its name",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE IF NOT EXISTS a (id INT PRIMARY KEY);\nCREATE TABLE b (id INT);\nCREATE TEMPORARY TABLE b (id INT PRIMARY KEY);",
			want: []finding{
				{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}},
				{2, "CREATE TABLE b (id INT)", "creates table b without a primary key", []int{0}},
			},
		},
		{
			name: "not settled by another table or a column rename",
			sql:  "CREATE TABLE public.a (id INT);\nALTER TABLE b ADD PRIMARY KEY (id);\nALTER TABLE a RENAME COLUMN id TO id2;\nDROP SCHEMA s2;\nDROP TABLE s2.a;",
			want: []finding{{0, "CREATE TABLE public.a (id INT)", "creates table public.a without a primary key", []int{0}}},
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE p DROP PRIMARY KEY, ADD PRIMARY KEY (id, note);\nALTER TABLE dup DROP PRIMARY KEY;\nALTER TABLE dup ADD CONSTRAINT pk_dup2 PRIMARY KEY (id);",
			want: []finding{{0, "DROP PRIMARY KEY", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name: "a primary key dropped by name",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t;\nALTER TABLE p DROP CONSTRAINT pk_p;\nDROP TABLE p;\nCREATE TABLE n1 (a INT CONSTRAINT pk_n1 PRIMARY KEY);\nALTER TABLE n1 DROP CONSTRAINT pk_n1;\nALTER TABLE t DROP CONSTRAINT uk_code;",
			want: []finding{
				{0, "DROP CONSTRAINT pk_t", "removes the primary key of t, and the change adds none back", []int{0}},
				{4, "DROP CONSTRAINT pk_n1", "removes the primary key of n1, and the change adds none back", []int{0}},
			},
		},
		{
			name: "a drop inside a block is not judged",
			sql:  "BEGIN\n  ALTER TABLE t DROP PRIMARY KEY;\n  ALTER TABLE p DROP CONSTRAINT pk_p;\nEND",
		},
	})
}

func TestOnlineMigration(t *testing.T) {
	runCases(t, review.OnlineMigration, []ruleCase{
		{name: "not asked for", sql: "ALTER TABLE t ADD COLUMN m INT;"},
		{
			name:    "asked for",
			sql:     "ALTER TABLE t ADD COLUMN m INT;",
			change:  review.Change{OnlineMigration: true},
			targets: []review.Target{{Schema: shop()}, {}},
			want:    []finding{{-1, "", "the change requests online migration, which Snowflake does not support", []int{0, 1}}},
		},
	})
}

func TestPriorBackup(t *testing.T) {
	runCases(t, review.PriorBackup, []ruleCase{
		{name: "not asked for", sql: "UPDATE t SET n = 1 WHERE id = 1;"},
		{
			name:    "asked for",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;",
			change:  review.Change{PriorBackup: true, MaxBackupSize: 1 << 20},
			targets: []review.Target{{Schema: shop(), BackupDatabaseExists: true}, {}},
			want:    []finding{{-1, "", "the change enables prior backup, which Bytebase does not take on Snowflake: the change runs without a backup", []int{0, 1}}},
		},
	})
}
