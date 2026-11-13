package catalog

import (
	"testing"
)

// TestScenario_C24 covers Section C24 "SHOW CREATE TABLE skip_gipk / Invisible
// PK" from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md, as far as it
// reaches OceanBase.
//
// MySQL 8.0.30+ has a Generated Invisible Primary Key (GIPK): with
// `sql_generate_invisible_primary_key=ON`, a CREATE TABLE that declares no
// PRIMARY KEY gets a hidden `my_row_id` column and a PRIMARY KEY over it.
// OceanBase has no such feature. It knows neither of MySQL's two session
// variables, a table without a primary key gets no column added, and
// my_row_id is an ordinary column name.
//
// omni gap: the catalog implements MySQL's feature. It takes both variables
// and, with the first one on, adds my_row_id and reserves the name.
func TestScenario_C24(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	const (
		generate = "SET SESSION sql_generate_invisible_primary_key = 1"
		show     = "SET SESSION show_gipk_in_create_table_and_information_schema = 1"
	)

	// -----------------------------------------------------------------
	// 24.1 No GIPK: the variables are unknown and no column is added
	// -----------------------------------------------------------------
	t.Run("24_1_gipk_hidden_by_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Oracle: ER_UNKNOWN_SYSTEM_VARIABLE (1193) for both variables.
		oracleRefuses(t, mc, generate, 1193)
		oracleRefuses(t, mc, show, 1193)
		// omni gap: the catalog takes both.
		assertIntEq(t, "omni "+generate+", error code", omniCode(c, generate), 0)
		assertIntEq(t, "omni "+show+", error code", omniCode(c, show), 0)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT)")

		// Oracle: the table has the two columns it was given.
		var columns int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'`,
			&columns)
		assertIntEq(t, "oracle 24.1 columns of t", columns, 2)

		// omni gap: the catalog generated my_row_id, since the variable it
		// took is on.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("omni 24.1: table t missing")
		}
		assertIntEq(t, "omni 24.1 columns of t", len(tbl.Columns), 3)
		if col := tbl.GetColumn("my_row_id"); col == nil || !col.GeneratedInvisiblePrimaryKey {
			t.Errorf("omni 24.1: want a generated my_row_id column, got %+v", col)
		}
	})

	// -----------------------------------------------------------------
	// 24.4 my_row_id is an ordinary column name
	// -----------------------------------------------------------------
	t.Run("24_4_gipk_name_collision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		const ddl = "CREATE TABLE t (my_row_id INT, a INT)"
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err != nil {
			t.Errorf("oracle 24.4: %s: %v", ddl, err)
		}
		// Oracle: my_row_id exists as the user's column.
		var name string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&name)
		assertStringEq(t, "oracle 24.4 my_row_id user column name", name, "my_row_id")

		// omni: the catalog takes the table too, while the variable is off.
		assertIntEq(t, "omni 24.4 "+ddl+", error code", omniCode(c, ddl), 0)
		if tbl := c.GetDatabase("testdb").GetTable("t"); tbl == nil || tbl.GetColumn("my_row_id") == nil {
			t.Errorf("omni 24.4: expected user column my_row_id")
		}

		// omni gap: with the variable on, the catalog reserves the name, as
		// MySQL does (ER_GIPK_FAILED_AUTOINC_COLUMN_NAME_RESERVED, 4108).
		reserving := scenarioNewCatalog(t)
		assertIntEq(t, "omni "+generate+", error code", omniCode(reserving, generate), 0)
		assertIntEq(t, "omni 24.4 "+ddl+" with the variable on, error code", omniCode(reserving, ddl), 4108)
	})
}
