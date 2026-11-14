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
		{"CREATE TABLE n (a int)", []ChangedResource{{"db", "public", "n", false}}},
		{"DROP TABLE t, app.only_app", []ChangedResource{{"db", "public", "t", true}, {"db", "app", "only_app", true}}},
		{"DROP INDEX t_idx", []ChangedResource{{"db", "public", "t", false}}},
		{"DROP INDEX app.t_idx", nil},
		{"ALTER TABLE only_app ADD COLUMN b int", []ChangedResource{{"db", "public", "only_app", true}}},
		{"ALTER VIEW v RENAME TO w", nil},
		{"ALTER TABLE t RENAME TO u", []ChangedResource{{"db", "public", "t", true}, {"db", "public", "u", false}}},
		{"ALTER TABLE t RENAME COLUMN a TO b", []ChangedResource{{"db", "public", "t", true}}},
		{"CREATE INDEX ON t (a)", []ChangedResource{{"db", "public", "t", false}}},
		{"TRUNCATE t", []ChangedResource{{"db", "public", "t", true}}},
		{"SELECT * INTO n FROM t", []ChangedResource{{"db", "public", "n", false}}},
		{"WITH d AS (DELETE FROM t RETURNING *) SELECT * FROM d", []ChangedResource{{"db", "public", "t", false}}},
		{"EXPLAIN ANALYZE DELETE FROM other.t", []ChangedResource{{"db", "other", "t", false}}},
		{"EXPLAIN DELETE FROM t", nil},
	}
	for _, tc := range tests {
		require.Equal(t, tc.want, changesOf(t, tc.script, nil).Resources, tc.script)
	}
}

func TestExtractChangedResourcesNewTableSchema(t *testing.T) {
	// With a search path of more than one schema, a new table whose first
	// schema does not exist has no schema.
	require.Equal(t, []ChangedResource{{"db", "", "n", false}}, changesOf(t, "CREATE TABLE n (a int)", []string{"missing", "public"}).Resources)
	require.Equal(t, []ChangedResource{{"db", "app", "n", false}}, changesOf(t, "CREATE TABLE n (a int)", []string{"app", "public"}).Resources)
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
		{"db", "app", "t", false},
		{"db", "public", "t", false},
		{"db", "app", "t", false},
		{"db", "public", "t", false},
		{"db", "public", "t", false},
		{"db", "public", "t", false},
	}, changes.Resources)
}

func TestExtractChangedResourcesUnwrapsExplainAnalyze(t *testing.T) {
	changes := changesOf(t, "SELECT 1; EXPLAIN ANALYZE DELETE FROM t;", nil)
	require.Equal(t, []DMLStatement{{Text: "DELETE FROM t;"}}, changes.DMLStatements)
}
