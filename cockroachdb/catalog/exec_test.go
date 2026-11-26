package catalog_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
)

// execLast runs a script and returns the refusal of its last statement,
// failing on any refusal before it.
func execLast(t *testing.T, script string) (*catalog.Catalog, error) {
	t.Helper()
	c := catalog.New()
	results, err := c.Exec(script)
	require.NoError(t, err)
	for i, r := range results {
		if r.Error != nil && i < len(results)-1 {
			t.Fatalf("line %d: %s: %v", r.Line, r.SQL, r.Error)
		}
	}
	return c, results[len(results)-1].Error
}

func code(err error) string {
	var refusal *catalog.Error
	if errors.As(err, &refusal) {
		return refusal.Code
	}
	return ""
}

func TestExecRefusals(t *testing.T) {
	for _, tc := range []struct{ script, code string }{
		{"ALTER TABLE missing.t ADD COLUMN x INT8", "42P01"},
		{"COMMENT ON TABLE missing.t IS 'x'", "42P01"},
		{"CREATE TABLE t (a INT8); CREATE INDEX i ON t (a) STORING (missing)", "42703"},
		{"CREATE TABLE t (a INT8); CREATE INDEX i ON t (a) WHERE missing > 0", "42703"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (id STRING REFERENCES p (id))", "42804"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (s STRING); ALTER TABLE c ADD FOREIGN KEY (s) REFERENCES p", "42804"},
	} {
		_, err := execLast(t, tc.script)
		require.Equal(t, tc.code, code(err), "%s: %v", tc.script, err)
	}
	_, err := execLast(t, "CREATE SCHEMA s; SET search_path = s")
	require.Error(t, err)
}

func TestExecRenameColumnRewritesComputedColumns(t *testing.T) {
	c, err := execLast(t, "CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t RENAME COLUMN a TO x")
	require.NoError(t, err)
	require.Equal(t, "x * 2", c.Metadata().Schemas[0].Tables[0].Columns[1].Generation.Expression)
}

func TestExecDropIndexCascadeDropsForeignKeys(t *testing.T) {
	c, err := execLast(t, "CREATE TABLE p (id INT8 PRIMARY KEY, u INT8 UNIQUE); CREATE TABLE c (u INT8 REFERENCES p (u)); DROP INDEX p@p_u_key CASCADE")
	require.NoError(t, err)
	for _, table := range c.Metadata().Schemas[0].Tables {
		require.Empty(t, table.ForeignKeys, table.Name)
	}
}

// TestSDLComputedAnnotations keeps an annotation that picks another overload
// and drops one CockroachDB writes on a literal's own type.
func TestSDLComputedAnnotations(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 AS (length('é':::BYTES)) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 AS (length('é':::STRING)) STORED);",
		catalog.SDLMigrationOptions{})
	require.Error(t, err, "%v", plan)
	plan, err = catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100:::DECIMAL) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100) STORED);",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)
}
