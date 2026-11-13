package catalog

import (
	"errors"
	"testing"
)

// TestOracle_OceanBaseKeyColumns holds validateIndexColumns to the engine: a key's
// columns exist, in a functional key part too, unless the same ALTER TABLE adds them
// later, and a FULLTEXT key's columns are text.
func TestOracle_OceanBaseKeyColumns(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), t TEXT, vb VARBINARY(10), j JSON)"
	runExecSteps(t, "oceanbase_catalog_key_columns", []struct {
		name  string
		steps []engineStep
	}{
		{"a missing column", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD KEY (zz)", 1072},
			{"ALTER TABLE a ADD KEY ((c + zz))", 1072},
			{"ALTER TABLE a ADD FULLTEXT (zz, c)", 1072},
			{"CREATE TABLE f (b INT PRIMARY KEY, FULLTEXT (zz))", 1072},
			{"ALTER TABLE a ADD KEY k (zz), ADD zz INT", taken},
		}},
		{"FULLTEXT", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD FULLTEXT (c, zz)", 1283},
			{"ALTER TABLE a ADD FULLTEXT (vb)", 1283},
			{"ALTER TABLE a ADD FULLTEXT (j)", 1283},
			{"CREATE FULLTEXT INDEX i ON a (c)", 1283},
			{"CREATE TABLE f (b INT PRIMARY KEY, c INT, FULLTEXT (c))", 1283},
			{"ALTER TABLE a ADD FULLTEXT (v, t)", taken},
		}},
	})
}

// TestFulltextColumnSQLState holds the SQLSTATE of each 1283 to what
// TestOracle_OceanBaseKeyColumns's tape recorded, which runExecSteps does not compare.
func TestFulltextColumnSQLState(t *testing.T) {
	c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE a (b INT PRIMARY KEY, c INT, vb VARBINARY(10), j JSON)")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ sql, state string }{
		{"ALTER TABLE a ADD FULLTEXT (c)", "42000"},
		{"ALTER TABLE a ADD FULLTEXT (j)", "42000"},
		{"ALTER TABLE a ADD FULLTEXT (vb)", "HY000"},
	} {
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		var refusal *Error
		if !errors.As(results[0].Error, &refusal) || refusal.Code != 1283 || refusal.SQLState != tc.state {
			t.Errorf("%s: got %v, want 1283 (%s)", tc.sql, results[0].Error, tc.state)
		}
	}
}
