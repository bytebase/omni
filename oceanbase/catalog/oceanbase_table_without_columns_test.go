package catalog

import "testing"

// TestOracle_OceanBaseTableWithoutColumns holds checkSomeColumnVisible to the engine:
// a table element list with no visible column is 4028 before anything else.
func TestOracle_OceanBaseTableWithoutColumns(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_table_without_columns", []struct {
		name  string
		steps []engineStep
	}{
		{"no visible column", []engineStep{
			{"CREATE TABLE a (b INT PRIMARY KEY)", taken},
			{"CREATE TABLE a (KEY (b))", 4028},
			{"CREATE TABLE IF NOT EXISTS a (KEY (b))", 4028},
			{"CREATE TABLE z (KEY (zz))", 4028},
			{"CREATE TABLE z (FOREIGN KEY (b) REFERENCES a (b))", 4028},
			{"CREATE TABLE z (KEY (b)) SELECT zz", 4028},
			{"CREATE TABLE z (b INT INVISIBLE)", 4028},
			{"CREATE TABLE z (b INT INVISIBLE, c INT)", taken},
		}},
	})
}
