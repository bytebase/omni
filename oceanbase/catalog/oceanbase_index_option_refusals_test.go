package catalog

import "testing"

// TestOracle_OceanBaseIndexOptionRefusals holds checkIndexOptions to the engine: the index
// options a key that is not a vector index may not take, and that they are refused
// before the key's name or columns are read.
func TestOracle_OceanBaseIndexOptionRefusals(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), KEY idx_c (c))"
	runExecSteps(t, "oceanbase_catalog_index_option_refusals", []struct {
		name  string
		steps []engineStep
	}{
		{"ALTER TABLE", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD KEY (c) WITH (a = 1)", 1235},
			{"ALTER TABLE a ADD KEY (zz) DATA_TABLE_ID 1", 1235},
			{"ALTER TABLE a ADD KEY idx_c (c) INDEX_TABLE_ID 1", 1235},
			{"ALTER TABLE a ADD UNIQUE (c) STORAGE_CACHE_POLICY (GLOBAL = 'hot')", 1235},
			{"ALTER TABLE a ADD b INT, ADD KEY (c) WITH (a = 1)", 1235},
			{"ALTER TABLE a ADD FULLTEXT (v) WITH (a = 1)", 1235},
		}},
		{"CREATE TABLE and CREATE INDEX", []engineStep{
			{table, taken},
			{"CREATE TABLE a (b INT PRIMARY KEY, KEY (b) WITH (a = 1))", 1235},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b) DATA_TABLE_ID 1)", 1235},
			{"CREATE INDEX i ON a (c) WITH (a = 1)", 1235},
			{"CREATE INDEX i ON a (c) DATA_TABLE_ID 1", 4179},
			{"CREATE INDEX i ON a (c) STORAGE_CACHE_POLICY (GLOBAL = 'hot')", 1235},
		}},
	})
}
