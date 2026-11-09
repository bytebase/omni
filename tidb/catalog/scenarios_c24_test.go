package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C24 covers Section C24 "Generated Invisible Primary Key (GIPK)"
// from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md.
//
// GIPK is a MySQL 8.0.30+ feature: with `sql_generate_invisible_primary_key=ON`
// a CREATE TABLE lacking a PRIMARY KEY silently gains a hidden
// `my_row_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT INVISIBLE` column plus a
// PRIMARY KEY over it.
//
// TiDB v8.5.5 does NOT implement GIPK. Measured against the container: both
// `sql_generate_invisible_primary_key` and
// `show_gipk_in_create_table_and_information_schema` are unknown system
// variables (SET fails with 1193), so no invisible primary key is ever
// generated and `my_row_id` is an ordinary identifier. These scenarios were
// forked from MySQL and are re-baselined here to assert TiDB's actual behavior.
// omni, which likewise has no GIPK generation, matches TiDB on every point —
// there is no omni gap to close and GIPK must NOT be implemented in omni.
func TestScenario_C24(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// unknownSysVar returns true when err reports an unknown system variable
	// (1193), which is how TiDB v8.5.5 rejects the MySQL-only GIPK vars.
	unknownSysVar := func(err error) bool {
		return err != nil && (strings.Contains(err.Error(), "1193") ||
			strings.Contains(strings.ToLower(err.Error()), "unknown system variable"))
	}

	// -----------------------------------------------------------------
	// 24.1 GIPK session variables do not exist in TiDB
	// -----------------------------------------------------------------
	t.Run("24_1_gipk_vars_unknown", func(t *testing.T) {
		scenarioReset(t, mc)

		_, errGen := mc.db.ExecContext(mc.ctx,
			"SET SESSION sql_generate_invisible_primary_key = ON")
		if !unknownSysVar(errGen) {
			t.Errorf("oracle 24.1: expected 1193 unknown system variable for sql_generate_invisible_primary_key, got: %v", errGen)
		}
		_, errShow := mc.db.ExecContext(mc.ctx,
			"SET SESSION show_gipk_in_create_table_and_information_schema = ON")
		if !unknownSysVar(errShow) {
			t.Errorf("oracle 24.1: expected 1193 unknown system variable for show_gipk_in_create_table_and_information_schema, got: %v", errShow)
		}
	})

	// -----------------------------------------------------------------
	// 24.2 CREATE TABLE without a PK generates no invisible my_row_id
	// -----------------------------------------------------------------
	t.Run("24_2_no_invisible_pk_generated", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT)")

		// Oracle: no my_row_id column, and SHOW CREATE does not mention it.
		var cols int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&cols)
		assertIntEq(t, "oracle 24.2 no my_row_id column", cols, 0)
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(create, "my_row_id") {
			t.Errorf("oracle 24.2: SHOW CREATE must not generate my_row_id, got:\n%s", create)
		}

		// Oracle: the first column is the user's own `a`, not a generated PK.
		var firstCol string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION LIMIT 1`,
			&firstCol)
		assertStringEq(t, "oracle 24.2 first column", firstCol, "a")

		// omni: catalog matches — no my_row_id, no synthesized PRIMARY KEY.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.2: table t missing")
			return
		}
		if tbl.GetColumn("my_row_id") != nil {
			t.Errorf("omni 24.2: catalog must not generate my_row_id (GIPK is not a TiDB feature)")
		}
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				t.Errorf("omni 24.2: unexpected PRIMARY KEY synthesized on a PK-less table")
			}
		}
	})

	// -----------------------------------------------------------------
	// 24.3 No invisible PK regardless of UNIQUE NOT NULL / explicit PK
	// -----------------------------------------------------------------
	t.Run("24_3_no_gipk_for_any_shape", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t1 (id INT PRIMARY KEY, a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t2 (id INT NOT NULL UNIQUE, a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t3 (id INT AUTO_INCREMENT, a INT, PRIMARY KEY (id))")

		// Oracle: none of the three tables gains a my_row_id — TiDB never
		// generates one, including for t2 where MySQL's GIPK would have.
		for _, name := range []string{"t1", "t2", "t3"} {
			var n int
			oracleScan(t, mc,
				`SELECT COUNT(*) FROM information_schema.COLUMNS
				 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+name+`' AND COLUMN_NAME='my_row_id'`,
				&n)
			assertIntEq(t, "oracle 24.3 "+name+" has no GIPK", n, 0)
		}

		// omni: catalog matches — no my_row_id on any of them.
		for _, name := range []string{"t1", "t2", "t3"} {
			tbl := c.GetDatabase("testdb").GetTable(name)
			if tbl == nil {
				t.Errorf("omni 24.3: table %s missing", name)
				continue
			}
			if tbl.GetColumn("my_row_id") != nil {
				t.Errorf("omni 24.3: %s must not have my_row_id (GIPK is not a TiDB feature)", name)
			}
		}
	})

	// -----------------------------------------------------------------
	// 24.4 my_row_id is an ordinary identifier (not reserved by GIPK)
	// -----------------------------------------------------------------
	t.Run("24_4_my_row_id_is_ordinary_identifier", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// With no GIPK, a user may freely declare a column named my_row_id.
		runOnBoth(t, mc, c, "CREATE TABLE t (my_row_id INT, a INT)")

		// Oracle: my_row_id exists as an ordinary user column.
		var name string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&name)
		assertStringEq(t, "oracle 24.4 my_row_id user column name", name, "my_row_id")

		// omni: catalog also keeps it as an ordinary column.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.4: table t missing")
			return
		}
		if tbl.GetColumn("my_row_id") == nil {
			t.Errorf("omni 24.4: expected user column my_row_id")
		}
	})
}
