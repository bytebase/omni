package analysis

import (
	"testing"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/stretchr/testify/require"
)

// fakeChangeMetadata has schemas public and app; public holds table t with
// index t_idx, and app holds table t and table only_app.
type fakeChangeMetadata struct{}

func (fakeChangeMetadata) HasSchema(name string) bool { return name == "public" || name == "app" }

func (fakeChangeMetadata) RelationSchema(searchPath []string, name string) string {
	for _, schema := range searchPath {
		if (schema == "public" && name == "t") || (schema == "app" && (name == "t" || name == "only_app")) {
			return schema
		}
	}
	return ""
}

func (fakeChangeMetadata) IndexTable(searchPath []string, name string) (string, string, bool) {
	for _, schema := range searchPath {
		if schema == "public" && name == "t_idx" {
			return "public", "t", true
		}
	}
	return "", "", false
}

func changesOf(t *testing.T, script string, searchPath []string) Changes {
	t.Helper()
	var changeStmts []ChangeStatement
	for _, segment := range omnipg.Split(script) {
		if segment.Empty() {
			continue
		}
		stmts, err := omnipg.Parse(segment.Text)
		require.NoError(t, err)
		for _, stmt := range stmts {
			changeStmts = append(changeStmts, ChangeStatement{Node: stmt.AST, Text: segment.Text})
		}
	}
	return ExtractChangedResources(changeStmts, "db", searchPath, fakeChangeMetadata{})
}

func TestExtractChangedResourcesTargets(t *testing.T) {
	tests := []struct {
		script string
		want   []ChangedResource
	}{
		{"CREATE TABLE n (a int)", []ChangedResource{{Database: "db", Schema: "public", Table: "n"}}},
		{"DROP TABLE t, app.only_app", []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}, {Database: "db", Schema: "app", Table: "only_app", AffectsData: true}}},
		{"DROP INDEX t_idx", []ChangedResource{{Database: "db", Schema: "public", Table: "t"}}},
		{"DROP INDEX app.t_idx", nil},
		{"ALTER TABLE only_app ADD COLUMN b int", []ChangedResource{{Database: "db", Schema: "public", Table: "only_app", AffectsData: true}}},
		{"ALTER VIEW v RENAME TO w", nil},
		{"ALTER TABLE t RENAME TO u", []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}, {Database: "db", Schema: "public", Table: "u"}}},
		{"ALTER TABLE t RENAME COLUMN a TO b", []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}}},
		{"CREATE INDEX ON t (a)", []ChangedResource{{Database: "db", Schema: "public", Table: "t"}}},
		{"TRUNCATE t", []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}}},
		{"SELECT * INTO n FROM t", []ChangedResource{{Database: "db", Schema: "public", Table: "n"}}},
		{"WITH d AS (DELETE FROM t RETURNING *) SELECT * FROM d", []ChangedResource{{Database: "db", Schema: "public", Table: "t"}}},
		{"EXPLAIN ANALYZE DELETE FROM other.t", []ChangedResource{{Database: "db", Schema: "other", Table: "t"}}},
		{"EXPLAIN DELETE FROM t", nil},
	}
	for _, tc := range tests {
		require.Equal(t, tc.want, changesOf(t, tc.script, nil).Resources, tc.script)
	}
}

func TestExtractChangedResourcesNewTableSchema(t *testing.T) {
	// With a search path of more than one schema, a new table whose first
	// schema does not exist has no schema.
	require.Equal(t, []ChangedResource{{Database: "db", Schema: "", Table: "n"}}, changesOf(t, "CREATE TABLE n (a int)", []string{"missing", "public"}).Resources)
	require.Equal(t, []ChangedResource{{Database: "db", Schema: "app", Table: "n"}}, changesOf(t, "CREATE TABLE n (a int)", []string{"app", "public"}).Resources)
}

func TestExtractChangedResourcesDML(t *testing.T) {
	changes := changesOf(t, "INSERT INTO t VALUES (1), (2); INSERT INTO t DEFAULT VALUES; INSERT INTO t SELECT * FROM t; UPDATE t SET a = 1; DELETE FROM t", nil)
	require.Equal(t, 3, changes.InsertCount)
	require.Equal(t, 3, changes.DMLCount)
	require.Equal(t, []DMLStatement{{Text: "INSERT INTO t SELECT * FROM t;"}, {Text: "UPDATE t SET a = 1;"}, {Text: "DELETE FROM t;"}}, changes.DMLStatements)

	// A data-modifying CTE makes the INSERT a sample.
	changes = changesOf(t, "WITH d AS (DELETE FROM t RETURNING *) INSERT INTO t VALUES (1)", nil)
	require.Equal(t, 0, changes.InsertCount)
	require.Equal(t, 1, changes.DMLCount)
}

func TestExtractChangedResourcesReplaysSearchPath(t *testing.T) {
	changes := changesOf(t, `SET search_path TO app; DELETE FROM t; RESET search_path; DELETE FROM t;
BEGIN; SET LOCAL search_path = app; UPDATE t SET a = 1; ROLLBACK; UPDATE t SET a = 2;
BEGIN; SET search_path = app; ROLLBACK; UPDATE t SET a = 3;
SET search_path = app; DISCARD ALL; UPDATE t SET a = 4`, nil)
	require.Equal(t, []DMLStatement{
		{Text: "DELETE FROM t;", SearchPath: []string{"app"}},
		{Text: "DELETE FROM t;"},
		{Text: "UPDATE t SET a = 1;", SearchPath: []string{"app"}},
		{Text: "UPDATE t SET a = 2;"},
		{Text: "UPDATE t SET a = 3;"},
		{Text: "UPDATE t SET a = 4;"},
	}, changes.DMLStatements)
	require.Equal(t, []ChangedResource{
		{Database: "db", Schema: "app", Table: "t"},
		{Database: "db", Schema: "public", Table: "t"},
		{Database: "db", Schema: "app", Table: "t"},
		{Database: "db", Schema: "public", Table: "t"},
		{Database: "db", Schema: "public", Table: "t"},
		{Database: "db", Schema: "public", Table: "t"},
	}, changes.Resources)
}

func TestExtractChangedResourcesUnwrapsExplainAnalyze(t *testing.T) {
	changes := changesOf(t, "SELECT 1; EXPLAIN ANALYZE DELETE FROM t;", nil)
	require.Equal(t, []DMLStatement{{Text: "DELETE FROM t;"}}, changes.DMLStatements)

	const multi = `EXPLAIN (ANALYZE, BUFFERS) WITH d AS (DELETE FROM t RETURNING id) INSERT INTO t2 SELECT id FROM d;
EXPLAIN ANALYZE MERGE INTO t3 USING t ON t3.id = t.id WHEN MATCHED THEN DELETE;
EXPLAIN UPDATE t4 SET c = 1;
EXPLAIN (ANALYZE FALSE) DELETE FROM t5;
EXPLAIN ANALYZE WITH d AS (DELETE FROM t6 RETURNING id) SELECT count(*) FROM d ORDER BY 1;`
	changes = changesOf(t, multi, nil)
	require.Equal(t, 3, changes.DMLCount)
	require.Equal(t, []DMLStatement{
		{Text: "WITH d AS (DELETE FROM t RETURNING id) INSERT INTO t2 SELECT id FROM d;"},
		{Text: "MERGE INTO t3 USING t ON t3.id = t.id WHEN MATCHED THEN DELETE;"},
		{Text: "WITH d AS (DELETE FROM t6 RETURNING id) SELECT count(*) FROM d ORDER BY 1;"},
	}, changes.DMLStatements)
	require.Equal(t, []ChangedResource{
		{Database: "db", Schema: "public", Table: "t"},
		{Database: "db", Schema: "public", Table: "t2"},
		{Database: "db", Schema: "public", Table: "t3"},
		{Database: "db", Schema: "public", Table: "t6"},
	}, changes.Resources)
}

func TestExtractChangedResourcesSearchPathTransactions(t *testing.T) {
	const statement = `UPDATE t SET c = 1;
SET search_path TO app;
UPDATE t SET c = 2;
SET search_path TO public;
UPDATE t SET c = 3;
SET search_path TO app;
RESET search_path;
UPDATE t SET c = 4;
SET search_path TO app;
SET search_path TO DEFAULT;
UPDATE t SET c = 5;
BEGIN;
SET LOCAL search_path TO app;
UPDATE t SET c = 6;
COMMIT;
UPDATE t SET c = 7;
BEGIN;
SET search_path TO app;
ROLLBACK;
UPDATE t SET c = 8;
BEGIN;
SET search_path TO app;
COMMIT AND CHAIN;
ROLLBACK;
UPDATE t SET c = 9;
DISCARD ALL;
UPDATE t SET c = 10;
SET search_path TO app;
BEGIN;
SET search_path TO other;
BEGIN;
SET search_path TO third;
ROLLBACK;
UPDATE t SET c = 11;
SET search_path TO other;
SET LOCAL search_path TO app;
ROLLBACK;
UPDATE t SET c = 12;
BEGIN;
SET LOCAL search_path TO app;
SET search_path FROM CURRENT;
COMMIT;
UPDATE t SET c = 13;`
	changes := changesOf(t, statement, []string{"public"})
	var formatted []string
	for _, stmt := range changes.DMLStatements {
		text := stmt.Text
		if stmt.SearchPath != nil {
			text = WithSearchPath(text, stmt.SearchPath)
		}
		formatted = append(formatted, text)
	}
	require.Equal(t, []string{
		"UPDATE t SET c = 1;",
		"SET LOCAL search_path TO \"app\";\nUPDATE t SET c = 2;",
		"UPDATE t SET c = 3;",
		"UPDATE t SET c = 4;",
		"UPDATE t SET c = 5;",
		"SET LOCAL search_path TO \"app\";\nUPDATE t SET c = 6;",
		"UPDATE t SET c = 7;",
		"UPDATE t SET c = 8;",
		"SET LOCAL search_path TO \"app\";\nUPDATE t SET c = 9;",
		"UPDATE t SET c = 10;",
		"SET LOCAL search_path TO \"app\";\nUPDATE t SET c = 11;",
		"SET LOCAL search_path TO \"other\";\nUPDATE t SET c = 12;",
		"SET LOCAL search_path TO \"app\";\nUPDATE t SET c = 13;",
	}, formatted)
}

func TestExtractChangedResourcesAdditionalShapes(t *testing.T) {
	tests := []struct {
		name            string
		script          string
		wantResources   []ChangedResource
		wantDMLCount    int
		wantInsertCount int
	}{
		{
			name:          "truncate multi-schema",
			script:        "TRUNCATE TABLE public.t1, myschema.t2;",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t1", AffectsData: true}, {Database: "db", Schema: "myschema", Table: "t2", AffectsData: true}},
		},
		{
			name:          "merge",
			script:        "MERGE INTO myschema.tgt t USING src s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a;",
			wantResources: []ChangedResource{{Database: "db", Schema: "myschema", Table: "tgt"}},
			wantDMLCount:  1,
		},
		{
			name:            "insert values on conflict",
			script:          "INSERT INTO t VALUES (1), (2) ON CONFLICT (id) DO UPDATE SET c = 1;",
			wantResources:   []ChangedResource{{Database: "db", Schema: "public", Table: "t"}},
			wantInsertCount: 2,
		},
		{
			name:          "insert select with read-only CTE",
			script:        "WITH x AS (SELECT id FROM t2) INSERT INTO t SELECT id FROM x;",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t"}},
			wantDMLCount:  1,
		},
		{
			name:          "data-modifying CTE feeding insert values",
			script:        "WITH i AS (INSERT INTO t VALUES (1) RETURNING id) INSERT INTO t2 VALUES (2), (3);",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t"}, {Database: "db", Schema: "public", Table: "t2"}},
			wantDMLCount:  1,
		},
		{
			name:          "data-modifying CTEs feeding delete",
			script:        "WITH u AS (UPDATE t SET c = 1 RETURNING id), m AS (MERGE INTO t3 USING t2 ON t3.id = t2.id WHEN MATCHED THEN DELETE) DELETE FROM t2 WHERE id IN (SELECT id FROM u);",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t"}, {Database: "db", Schema: "public", Table: "t3"}, {Database: "db", Schema: "public", Table: "t2"}},
			wantDMLCount:  1,
		},
		{
			name:          "data-modifying CTE in CTAS",
			script:        "CREATE TABLE t_moved AS WITH d AS (DELETE FROM t RETURNING *) SELECT * FROM d;",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t_moved"}, {Database: "db", Schema: "public", Table: "t"}},
			wantDMLCount:  1,
		},
		{
			name:          "read-only CTE in select",
			script:        "WITH x AS (SELECT id FROM t) SELECT * FROM x;",
			wantResources: nil,
		},
		{
			name:          "alter table rename constraint",
			script:        "ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_pk;",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}},
		},
		{
			name:          "alter trigger rename",
			script:        "ALTER TRIGGER trg ON t RENAME TO trg2;",
			wantResources: []ChangedResource{{Database: "db", Schema: "public", Table: "t", AffectsData: true}},
		},
		{
			name:          "alter index rename",
			script:        "ALTER INDEX t_idx RENAME TO t_idx2;",
			wantResources: nil,
		},
		{
			name:          "create table as",
			script:        "CREATE TABLE myschema.t1 AS SELECT * FROM src;",
			wantResources: []ChangedResource{{Database: "db", Schema: "myschema", Table: "t1"}},
		},
		{
			name:          "select into in union first arm",
			script:        "SELECT a INTO myschema.t3 FROM src UNION SELECT a FROM src2;",
			wantResources: []ChangedResource{{Database: "db", Schema: "myschema", Table: "t3"}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			changes := changesOf(t, tc.script, nil)
			require.Equal(t, tc.wantResources, changes.Resources)
			require.Equal(t, tc.wantDMLCount, changes.DMLCount)
			require.Equal(t, tc.wantInsertCount, changes.InsertCount)
		})
	}
}
