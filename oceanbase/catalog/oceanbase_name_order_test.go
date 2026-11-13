package catalog

import "testing"

// TestOracle_OceanBaseNameOrder holds to the engine when OceanBase looks a name up:
// CREATE TABLE last, after its whole definition; LIKE, DROP TABLE, DROP VIEW,
// RENAME TABLE and the key commands with their own orders and codes.
func TestOracle_OceanBaseNameOrder(t *testing.T) {
	const setup = "CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY idx_c (c))"
	runExecSteps(t, "oceanbase_catalog_name_order", []struct {
		name  string
		steps []engineStep
	}{
		{"CREATE TABLE", []engineStep{
			{setup, taken},
			{"CREATE TABLE a (x INT, x INT)", 1060},
			{"CREATE TABLE IF NOT EXISTS a (x INT, x INT)", 1060},
			{"CREATE TABLE a (x INT, KEY k (x), KEY k (x))", 1061},
			{"CREATE TABLE a (x VARCHAR(300000))", 1074},
			{"CREATE TABLE IF NOT EXISTS a (x INT)", taken},
			{"CREATE TABLE a LIKE a", 1066},
			{"CREATE TABLE IF NOT EXISTS a LIKE a", 1066},
			{"CREATE TABLE zz LIKE zz", 1066},
			{"CREATE TABLE a LIKE qq", 1146},
		}},
		{"DROP and RENAME", []engineStep{
			{setup, taken},
			{"CREATE VIEW v AS SELECT 1 AS x", taken},
			{"DROP TABLE zz, a, zz", 1066},
			{"DROP TABLE IF EXISTS a, a", 1066},
			{"DROP VIEW a", 1347},
			{"DROP VIEW IF EXISTS v, a", 1347},
			{"RENAME TABLE zz TO qq", 1017},
			{"RENAME TABLE zz TO a", 1050},
		}},
		{"keys", []engineStep{
			{setup, taken},
			{"ALTER TABLE a ALTER INDEX zz VISIBLE", 1176},
			{"ALTER TABLE a RENAME KEY zz TO idx_c", 1176},
			{"ALTER TABLE a RENAME KEY idx_c TO idx_c", 1061},
			{"ALTER TABLE a RENAME KEY idx_c TO idx_d", taken},
		}},
	})
}
