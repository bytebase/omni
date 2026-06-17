package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C25 covers Section C25 "DECIMAL defaults" from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md.
//
// MySQL canonicalizes DECIMAL columns:
//   - DECIMAL              → DECIMAL(10,0)
//   - DECIMAL(M)           → DECIMAL(M,0)
//   - DECIMAL(M,D)         → DECIMAL(M,D), with 1 <= M <= 65 and 0 <= D <= 30, D <= M
//   - NUMERIC              → synonym, stored as decimal in information_schema
//   - UNSIGNED / ZEROFILL  → flags after the (M,D) spec
//
// information_schema.COLUMNS exposes NUMERIC_PRECISION/NUMERIC_SCALE and
// COLUMN_TYPE; SHOW CREATE TABLE always renders the fully-qualified
// `decimal(M,D)` form, never the bare DECIMAL or DECIMAL(M) shorthand.
//
// Failures in omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c25.md.
func TestScenario_C25(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// -----------------------------------------------------------------
	// 25.1 DECIMAL with no precision/scale → DECIMAL(10,0)
	// -----------------------------------------------------------------
	t.Run("25_1_decimal_default_10_0", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (d DECIMAL)")

		// Oracle: COLUMN_TYPE / NUMERIC_PRECISION / NUMERIC_SCALE.
		var colType string
		var prec, scale int
		oracleScan(t, mc,
			`SELECT COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
			 FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'`,
			&colType, &prec, &scale)
		assertStringEq(t, "oracle column_type", strings.ToLower(colType), "decimal(10,0)")
		assertIntEq(t, "oracle numeric_precision", prec, 10)
		assertIntEq(t, "oracle numeric_scale", scale, 0)

		// omni: column type renders as decimal(10,0).
		col := c25Col(t, c, "t", "d")
		if col == nil {
			return
		}
		assertStringEq(t, "omni column_type", strings.ToLower(col.ColumnType), "decimal(10,0)")
	})

	// -----------------------------------------------------------------
	// 25.2 DECIMAL precision-only → scale defaults to 0
	// -----------------------------------------------------------------
	t.Run("25_2_precision_only_scale_zero", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			d5  DECIMAL(5),
			d15 DECIMAL(15),
			d65 DECIMAL(65)
		)`)

		// Oracle: COLUMN_TYPE rows.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
			 FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION`)
		got := map[string]string{}
		for _, r := range rows {
			got[asString(r[0])] = strings.ToLower(asString(r[1]))
		}
		assertStringEq(t, "oracle d5",  got["d5"],  "decimal(5,0)")
		assertStringEq(t, "oracle d15", got["d15"], "decimal(15,0)")
		assertStringEq(t, "oracle d65", got["d65"], "decimal(65,0)")

		// Oracle: SHOW CREATE TABLE must also use the explicit zero scale.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lower := strings.ToLower(create)
		for _, want := range []string{"decimal(5,0)", "decimal(15,0)", "decimal(65,0)"} {
			if !strings.Contains(lower, want) {
				t.Errorf("oracle SHOW CREATE TABLE: missing %q in:\n%s", want, create)
			}
		}

		// omni
		for _, want := range []struct{ name, typ string }{
			{"d5", "decimal(5,0)"},
			{"d15", "decimal(15,0)"},
			{"d65", "decimal(65,0)"},
		} {
			col := c25Col(t, c, "t", want.name)
			if col == nil {
				continue
			}
			assertStringEq(t, "omni "+want.name+" column_type",
				strings.ToLower(col.ColumnType), want.typ)
		}
	})

	// -----------------------------------------------------------------
	// 25.3 DECIMAL bounds: max P=65, S=30, scale > precision rejection
	// -----------------------------------------------------------------
	t.Run("25_3_bounds_max_p_s_and_scale_gt_precision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// 25.3a: DECIMAL(65,30) accepted on both sides.
		runOnBoth(t, mc, c, "CREATE TABLE ok_max_p (d DECIMAL(65,30))")

		var colType string
		oracleScan(t, mc,
			`SELECT COLUMN_TYPE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='ok_max_p' AND COLUMN_NAME='d'`,
			&colType)
		assertStringEq(t, "oracle ok_max_p column_type",
			strings.ToLower(colType), "decimal(65,30)")
		if col := c25Col(t, c, "ok_max_p", "d"); col != nil {
			assertStringEq(t, "omni ok_max_p column_type",
				strings.ToLower(col.ColumnType), "decimal(65,30)")
		}

		// 25.3b: DECIMAL(66,0) → ER_TOO_BIG_PRECISION (1426).
		c25AssertBothError(t, mc, c,
			"CREATE TABLE err_p_gt_65 (d DECIMAL(66, 0))",
			"precision 66", "25.3b precision > 65")

		// 25.3c: DECIMAL(40,31) → ER_TOO_BIG_SCALE (1425).
		c25AssertBothError(t, mc, c,
			"CREATE TABLE err_s_gt_30 (d DECIMAL(40, 31))",
			"scale 31", "25.3c scale > 30")

		// 25.3d: DECIMAL(5,6) → ER_M_BIGGER_THAN_D (1427).
		c25AssertBothError(t, mc, c,
			"CREATE TABLE err_s_gt_p (d DECIMAL(5, 6))",
			"M must be >= D", "25.3d scale > precision")

		// 25.3e: DECIMAL(-1,0) → parse error on both sides.
		c25AssertBothError(t, mc, c,
			"CREATE TABLE err_neg (d DECIMAL(-1, 0))",
			"", "25.3e negative precision")
	})

	// -----------------------------------------------------------------
	// 25.4 UNSIGNED / ZEROFILL / NUMERIC synonym
	// -----------------------------------------------------------------
	t.Run("25_4_unsigned_zerofill_numeric_synonym", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			d1 DECIMAL(10,2) UNSIGNED,
			d2 DECIMAL(10,2) UNSIGNED ZEROFILL,
			d3 NUMERIC(10,2) UNSIGNED
		)`)

		// Oracle: NUMERIC reported as decimal; precision=10, scale=2 for all.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME, COLUMN_TYPE, DATA_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
			 FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION`)
		want := map[string]string{
			"d1": "decimal(10,2) unsigned",
			"d2": "decimal(10,2) unsigned zerofill",
			"d3": "decimal(10,2) unsigned",
		}
		for _, r := range rows {
			name := asString(r[0])
			ct := strings.ToLower(asString(r[1]))
			dt := strings.ToLower(asString(r[2]))
			assertStringEq(t, "oracle "+name+" column_type", ct, want[name])
			assertStringEq(t, "oracle "+name+" data_type", dt, "decimal")
		}

		// omni: catalog should round-trip these as well.
		for _, name := range []string{"d1", "d2", "d3"} {
			col := c25Col(t, c, "t", name)
			if col == nil {
				continue
			}
			ct := strings.ToLower(col.ColumnType)
			if !strings.Contains(ct, "decimal(10,2)") {
				t.Errorf("omni %s column_type: got %q, want substring decimal(10,2)", name, ct)
			}
			if !strings.Contains(ct, "unsigned") {
				t.Errorf("omni %s column_type: got %q, want substring unsigned", name, ct)
			}
			if name == "d2" && !strings.Contains(ct, "zerofill") {
				t.Errorf("omni d2 column_type: got %q, want substring zerofill", ct)
			}
		}
	})

	// -----------------------------------------------------------------
	// 25.5 Zero-scale rendering: DECIMAL / DECIMAL(5) / DECIMAL(5,0) all collapse
	// -----------------------------------------------------------------
	t.Run("25_5_zero_scale_rendering", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a DECIMAL,
			b DECIMAL(5),
			c DECIMAL(5,0),
			d DECIMAL(10,0)
		)`)

		// Oracle: information_schema rows.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION`)
		expected := map[string]string{
			"a": "decimal(10,0)",
			"b": "decimal(5,0)",
			"c": "decimal(5,0)",
			"d": "decimal(10,0)",
		}
		for _, r := range rows {
			n := asString(r[0])
			ct := strings.ToLower(asString(r[1]))
			assertStringEq(t, "oracle "+n+" column_type", ct, expected[n])
		}

		// Oracle: SHOW CREATE TABLE renders explicit (M,0).
		create := strings.ToLower(oracleShow(t, mc, "SHOW CREATE TABLE t"))
		for _, want := range []string{"decimal(10,0)", "decimal(5,0)"} {
			if !strings.Contains(create, want) {
				t.Errorf("oracle SHOW CREATE TABLE: missing %q in:\n%s", want, create)
			}
		}
		// And must NOT render the bare or single-arg form.
		for _, bad := range []string{"decimal(5)", "decimal(10)", "decimal,"} {
			if strings.Contains(create, bad) {
				t.Errorf("oracle SHOW CREATE TABLE: should not contain %q in:\n%s", bad, create)
			}
		}

		// omni
		for name, want := range expected {
			col := c25Col(t, c, "t", name)
			if col == nil {
				continue
			}
			assertStringEq(t, "omni "+name+" column_type",
				strings.ToLower(col.ColumnType), want)
		}
	})
}

// c25Col fetches a column from testdb.<table>, reporting via t.Error when
// missing rather than aborting the test.
func c25Col(t *testing.T, c *Catalog, table, col string) *Column {
	t.Helper()
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Errorf("omni: database testdb missing")
		return nil
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		t.Errorf("omni: table %s missing", table)
		return nil
	}
	column := tbl.GetColumn(col)
	if column == nil {
		t.Errorf("omni: column %s.%s missing", table, col)
		return nil
	}
	return column
}

// c25AssertBothError runs the same DDL on the MySQL container and the omni
// catalog, asserting that BOTH return an error. If wantSubstr is non-empty it
// must appear in the MySQL container error message (case-insensitive). label
// is used only in error messages for easier triage.
func c25AssertBothError(t *testing.T, mc *mysqlContainer, c *Catalog, ddl, wantSubstr, label string) {
	t.Helper()

	_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
	if oracleErr == nil {
		t.Errorf("oracle %s: expected error for %q, got nil", label, ddl)
	} else if wantSubstr != "" &&
		!strings.Contains(strings.ToLower(oracleErr.Error()), strings.ToLower(wantSubstr)) {
		t.Errorf("oracle %s: error %q missing substring %q", label, oracleErr.Error(), wantSubstr)
	}

	results, err := c.Exec(ddl+";", nil)
	if err != nil {
		// Parse-level rejection is fine — counts as an error from omni.
		return
	}
	sawErr := false
	for _, r := range results {
		if r.Error != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Errorf("omni %s: expected error for %q, got nil", label, ddl)
	}
}
