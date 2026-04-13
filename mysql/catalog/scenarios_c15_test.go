package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C15 covers Section C15 "Column positioning defaults" from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL against both
// a real MySQL 8.0 container and the omni catalog, then asserts that both
// agree on the resulting column ordering.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c15.md.
func TestScenario_C15(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c15OracleColumnOrder fetches the column names from MySQL's
	// information_schema.COLUMNS for testdb.<table>, ordered by
	// ORDINAL_POSITION.
	c15OracleColumnOrder := func(t *testing.T, table string) []string {
		t.Helper()
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+table+`'
			 ORDER BY ORDINAL_POSITION`)
		out := make([]string, 0, len(rows))
		for _, r := range rows {
			if len(r) == 0 {
				continue
			}
			out = append(out, asString(r[0]))
		}
		return out
	}

	// c15OmniColumnOrder fetches column names from the omni catalog in
	// declaration order (Table.Columns slice order).
	c15OmniColumnOrder := func(c *Catalog, table string) []string {
		db := c.GetDatabase("testdb")
		if db == nil {
			return nil
		}
		tbl := db.GetTable(table)
		if tbl == nil {
			return nil
		}
		out := make([]string, 0, len(tbl.Columns))
		for _, col := range tbl.Columns {
			out = append(out, col.Name)
		}
		return out
	}

	c15AssertOrder := func(t *testing.T, label string, got, want []string) {
		t.Helper()
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("%s: got %v, want %v", label, got, want)
		}
	}

	// --- 15.1 ADD COLUMN appends to end -----------------------------------
	t.Run("15_1_add_column_appends_end", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);
			ALTER TABLE t ADD COLUMN c INT;`)

		want := []string{"a", "b", "c"}
		c15AssertOrder(t, "oracle order after ADD COLUMN c", c15OracleColumnOrder(t, "t"), want)
		c15AssertOrder(t, "omni order after ADD COLUMN c", c15OmniColumnOrder(c, "t"), want)
	})

	// --- 15.2 ADD COLUMN ... FIRST ----------------------------------------
	t.Run("15_2_add_column_first", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);
			ALTER TABLE t ADD COLUMN c INT FIRST;`)

		want := []string{"c", "a", "b"}
		c15AssertOrder(t, "oracle order after ADD COLUMN c FIRST", c15OracleColumnOrder(t, "t"), want)
		c15AssertOrder(t, "omni order after ADD COLUMN c FIRST", c15OmniColumnOrder(c, "t"), want)
	})

	// --- 15.3 ADD COLUMN ... AFTER col ------------------------------------
	t.Run("15_3_add_column_after", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, c INT);
			ALTER TABLE t ADD COLUMN x INT AFTER a;`)

		want := []string{"a", "x", "b", "c"}
		c15AssertOrder(t, "oracle order after ADD COLUMN x AFTER a", c15OracleColumnOrder(t, "t"), want)
		c15AssertOrder(t, "omni order after ADD COLUMN x AFTER a", c15OmniColumnOrder(c, "t"), want)
	})

	// --- 15.4 MODIFY retains position unless FIRST/AFTER ------------------
	t.Run("15_4_modify_retains_position", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, c INT);
			ALTER TABLE t MODIFY b BIGINT;`)

		want := []string{"a", "b", "c"}
		c15AssertOrder(t, "oracle order after MODIFY b BIGINT", c15OracleColumnOrder(t, "t"), want)
		c15AssertOrder(t, "omni order after MODIFY b BIGINT", c15OmniColumnOrder(c, "t"), want)

		// Also verify the type actually changed on both sides.
		var dataType string
		oracleScan(t, mc,
			`SELECT DATA_TYPE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`,
			&dataType)
		if dataType != "bigint" {
			t.Errorf("oracle DATA_TYPE for b: got %q, want %q", dataType, "bigint")
		}

		if tbl := c.GetDatabase("testdb").GetTable("t"); tbl != nil {
			if col := tbl.GetColumn("b"); col != nil {
				if !strings.EqualFold(col.DataType, "bigint") {
					t.Errorf("omni DataType for b: got %q, want %q", col.DataType, "bigint")
				}
			} else {
				t.Errorf("omni: column b missing after MODIFY")
			}
		}
	})

	// --- 15.5 multiple ADD COLUMN in one ALTER, left-to-right resolution ---
	t.Run("15_5_multi_add_left_to_right", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);
			ALTER TABLE t ADD COLUMN x INT AFTER a, ADD COLUMN y INT AFTER x;`)

		want := []string{"a", "x", "y", "b"}
		c15AssertOrder(t, "oracle order after multi-ADD", c15OracleColumnOrder(t, "t"), want)
		c15AssertOrder(t, "omni order after multi-ADD", c15OmniColumnOrder(c, "t"), want)
	})
}
