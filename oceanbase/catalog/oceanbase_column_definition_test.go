package catalog

import "testing"

// TestOracle_OceanBaseColumnDefinition holds checkColumnDefinition, checkAddedColumn
// and the literal defaults literalDefaultFits reads to the engine, and that each is
// refused before ADD COLUMN's 1060 for a column the table has.
func TestOracle_OceanBaseColumnDefinition(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT)"
	runExecSteps(t, "oceanbase_catalog_column_definition", []struct {
		name  string
		steps []engineStep
	}{
		{"before the name", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD c TINYINT DEFAULT 's'", 1067},
			{"ALTER TABLE a ADD c TINYINT DEFAULT CURRENT_TIMESTAMP", 1067},
			{"ALTER TABLE a ADD c TINYINT DEFAULT (b)", 1235},
			{"ALTER TABLE a ADD c TINYINT AS (zz)", 1054},
			{"ALTER TABLE a ADD c TINYINT AS (c) CHECK (c)", 1235},
			{"ALTER TABLE a ADD c TINYINT SRID 1", 1221},
			{"ALTER TABLE a ADD c TINYINT PRIMARY KEY", 1068},
			{"ALTER TABLE a ADD c TINYINT CHECK (c)", 3812},
			{"ALTER TABLE a ADD c TINYINT CHECK (c > 0)", 1060},
			{"ALTER TABLE a ADD e TINYINT AS (e)", 1054},
			{"ALTER TABLE a ADD e TINYINT REFERENCES a (b) CHECK (e)", taken},
		}},
		{"literal defaults", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD e1 TINYINT DEFAULT '1.5'", taken},
			{"ALTER TABLE a ADD e2 TINYINT DEFAULT 128", 1067},
			{"ALTER TABLE a ADD e2 TINYINT UNSIGNED DEFAULT 255", taken},
			{"ALTER TABLE a ADD e3 INT DEFAULT ''", 1067},
			{"ALTER TABLE a ADD e3 TINYINT DEFAULT 0x10", taken},
			{"ALTER TABLE a ADD e4 DECIMAL(3,1) DEFAULT 100", 1067},
			{"ALTER TABLE a ADD e4 BIT(1) DEFAULT 2", 1067},
			{"ALTER TABLE a ADD e5 CHAR(2) DEFAULT 'abc'", 1067},
			{"ALTER TABLE a ADD e5 VARCHAR(2) DEFAULT 'éé'", taken},
			{"ALTER TABLE a ADD e6 ENUM('a') DEFAULT 'x'", 1067},
			{"ALTER TABLE a ADD e6 SET('a','b') DEFAULT 'a,b'", taken},
			{"ALTER TABLE a ADD e7 DATETIME DEFAULT '2020-13-01'", 1067},
			{"ALTER TABLE a ADD e7 DATE DEFAULT '20200101'", taken},
			{"ALTER TABLE a ADD e8 TIME DEFAULT 's'", 1067},
			{"ALTER TABLE a ADD e8 YEAR DEFAULT 's'", taken},
			{"ALTER TABLE a ADD e9 DATE DEFAULT CURRENT_TIMESTAMP", 1067},
			{"ALTER TABLE a ALTER c SET DEFAULT 's'", 1067},
			{"ALTER TABLE a ALTER zz SET DEFAULT (c)", 1235},
			{"ALTER TABLE a MODIFY c TINYINT DEFAULT 's'", 1067},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, c TINYINT DEFAULT 's')", 1067},
		}},
	})
}
