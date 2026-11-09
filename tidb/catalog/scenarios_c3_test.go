package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C3 covers section C3 (Nullability & default promotion) from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that both real
// MySQL 8.0 and the omni catalog agree on the implicit nullability/default
// rules MySQL applies during CREATE TABLE.
//
// Failures in omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c3.md.
func TestScenario_C3(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 3.1 First TIMESTAMP NOT NULL is NOT auto-promoted under TiDB -------
	t.Run("3_1_first_timestamp_no_promotion", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// TiDB fixes explicit_defaults_for_timestamp=ON (it is a read-only
		// variable in v8.5.5), so the legacy first-TIMESTAMP promotion path
		// never fires: neither column gets DEFAULT CURRENT_TIMESTAMP or
		// ON UPDATE CURRENT_TIMESTAMP.
		ddl := `CREATE TABLE t (
    ts1 TIMESTAMP NOT NULL,
    ts2 TIMESTAMP NOT NULL
)`
		runOnBoth(t, mc, c, ddl)

		// Oracle: SHOW CREATE TABLE on container — no promotion on either column.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if strings.Contains(lo, "default current_timestamp") {
			t.Errorf("oracle: TiDB should NOT promote any TIMESTAMP to DEFAULT CURRENT_TIMESTAMP\n%s", create)
		}
		if strings.Contains(lo, "on update current_timestamp") {
			t.Errorf("oracle: TiDB should NOT add ON UPDATE CURRENT_TIMESTAMP\n%s", create)
		}

		// omni mirrors TiDB's default mode: no promotion on either column.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		ts1 := tbl.GetColumn("ts1")
		ts2 := tbl.GetColumn("ts2")
		if ts1 == nil || ts2 == nil {
			t.Errorf("omni: ts1/ts2 columns missing")
			return
		}
		for name, col := range map[string]*Column{"ts1": ts1, "ts2": ts2} {
			if col.Default != nil && strings.Contains(strings.ToLower(*col.Default), "current_timestamp") {
				t.Errorf("omni: %s should NOT auto-promote to CURRENT_TIMESTAMP, got %q", name, *col.Default)
			}
			if col.OnUpdate != "" && strings.Contains(strings.ToLower(col.OnUpdate), "current_timestamp") {
				t.Errorf("omni: %s should NOT auto-promote ON UPDATE, got %q", name, col.OnUpdate)
			}
		}
	})

	// --- 3.2 PRIMARY KEY implies NOT NULL -----------------------------------
	t.Run("3_2_primary_key_implies_not_null", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT PRIMARY KEY)")

		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`,
			&isNullable)
		assertStringEq(t, "oracle IS_NULLABLE", isNullable, "NO")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Errorf("omni: column a missing")
			return
		}
		assertBoolEq(t, "omni column a Nullable", col.Nullable, false)
	})

	// --- 3.3 AUTO_INCREMENT implies NOT NULL --------------------------------
	t.Run("3_3_auto_increment_implies_not_null", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT AUTO_INCREMENT, KEY (a))")

		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`,
			&isNullable)
		assertStringEq(t, "oracle IS_NULLABLE", isNullable, "NO")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Errorf("omni: column a missing")
			return
		}
		assertBoolEq(t, "omni column a Nullable", col.Nullable, false)
	})

	// --- 3.4 Explicit NULL on PRIMARY KEY column is a hard error ------------
	t.Run("3_4_explicit_null_pk_errors", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Oracle: MySQL must reject.
		_, mysqlErr := mc.db.ExecContext(mc.ctx,
			"CREATE TABLE t (a INT NULL PRIMARY KEY)")
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_PRIMARY_CANT_HAVE_NULL (1171), got nil")
		} else if !strings.Contains(mysqlErr.Error(), "1171") &&
			!strings.Contains(strings.ToLower(mysqlErr.Error()), "primary") {
			t.Errorf("oracle: expected 1171 ER_PRIMARY_CANT_HAVE_NULL, got %v", mysqlErr)
		}

		// omni: should also reject.
		results, err := c.Exec("CREATE TABLE t (a INT NULL PRIMARY KEY);", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects NULL + PK", omniErrored, true)
	})

	// --- 3.5 UNIQUE does NOT imply NOT NULL ---------------------------------
	t.Run("3_5_unique_does_not_imply_not_null", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT UNIQUE)")

		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`,
			&isNullable)
		assertStringEq(t, "oracle IS_NULLABLE (UNIQUE)", isNullable, "YES")

		// Verify a UNIQUE index actually exists on the container side.
		var idxName string
		oracleScan(t, mc,
			`SELECT INDEX_NAME FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0`,
			&idxName)
		if idxName == "" {
			t.Errorf("oracle: expected one UNIQUE index, got none")
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Errorf("omni: column a missing")
			return
		}
		assertBoolEq(t, "omni column a Nullable (UNIQUE)", col.Nullable, true)

		omniIdx := omniUniqueIndexNames(tbl)
		if len(omniIdx) == 0 {
			t.Errorf("omni: expected a UNIQUE index, got none")
		}
	})

	// --- 3.6 Generated column nullability derived from expression -----------
	t.Run("3_6_generated_column_nullability", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
    a INT NULL,
    b INT GENERATED ALWAYS AS (a+1) VIRTUAL
)`
		runOnBoth(t, mc, c, ddl)

		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`,
			&isNullable)
		assertStringEq(t, "oracle IS_NULLABLE (gcol)", isNullable, "YES")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Errorf("omni: column b missing")
			return
		}
		assertBoolEq(t, "omni gcol Nullable", col.Nullable, true)
	})

	// --- 3.7 Explicit NULL + AUTO_INCREMENT --------------------------------
	// SCENARIOS expectation: MySQL errors. Empirically MySQL 8.0.45 accepts
	// the statement and silently promotes id to NOT NULL (AUTO_INCREMENT
	// wins). We test the observed silent-coercion behavior on both sides
	// and document the discrepancy with the SCENARIOS file.
	t.Run("3_7_null_plus_auto_increment_silent_promote", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := "CREATE TABLE t (id INT NULL AUTO_INCREMENT, KEY(id))"
		runOnBoth(t, mc, c, ddl)

		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='id'`,
			&isNullable)
		assertStringEq(t, "oracle id IS_NULLABLE", isNullable, "NO")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("id")
		if col == nil {
			t.Errorf("omni: column id missing")
			return
		}
		assertBoolEq(t, "omni id Nullable (NULL + AUTO_INCREMENT silently promoted)",
			col.Nullable, false)
	})

	// --- 3.8 Bare TIMESTAMP columns under TiDB (explicit_defaults=ON) -------
	t.Run("3_8_bare_timestamp_nullable_default_null", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// explicit_defaults_for_timestamp is fixed ON in TiDB, so a bare
		// TIMESTAMP column (no NULL/NOT NULL, no default) becomes
		// `timestamp NULL DEFAULT NULL` on both columns — no legacy promotion
		// and no implicit '0000-00-00' zero default.
		ddl := "CREATE TABLE t (ts1 TIMESTAMP, ts2 TIMESTAMP)"
		runOnBoth(t, mc, c, ddl)

		// Oracle: both columns nullable with DEFAULT NULL; no promotion.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if strings.Contains(lo, "default current_timestamp") {
			t.Errorf("oracle: TiDB should NOT promote a bare TIMESTAMP to DEFAULT CURRENT_TIMESTAMP\n%s", create)
		}
		if strings.Contains(lo, "0000-00-00 00:00:00") {
			t.Errorf("oracle: TiDB should NOT apply a '0000-00-00 00:00:00' default\n%s", create)
		}
		if strings.Count(lo, "timestamp null default null") != 2 {
			t.Errorf("oracle: expected both columns rendered `timestamp NULL DEFAULT NULL`\n%s", create)
		}

		// omni mirrors TiDB: both columns nullable, no CURRENT_TIMESTAMP default.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		ts1 := tbl.GetColumn("ts1")
		ts2 := tbl.GetColumn("ts2")
		if ts1 == nil || ts2 == nil {
			t.Errorf("omni: ts1/ts2 columns missing")
			return
		}
		for name, col := range map[string]*Column{"ts1": ts1, "ts2": ts2} {
			assertBoolEq(t, "omni "+name+" Nullable", col.Nullable, true)
			if col.Default != nil && strings.Contains(strings.ToLower(*col.Default), "current_timestamp") {
				t.Errorf("omni: %s should NOT auto-promote to CURRENT_TIMESTAMP, got %q", name, *col.Default)
			}
		}
	})
}
