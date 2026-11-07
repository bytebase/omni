package catalog

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestScenario_C3 checks TiDB's nullability and default rules against the
// native server, including refusal to enable MySQL's legacy timestamp mode.
func TestScenario_C3(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 3.1 NOT NULL timestamps do not implicitly promote -------------
	t.Run("3_1_first_timestamp_promotion", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// TiDB refuses the legacy mode setting; neither timestamp promotes.
		assertReadonlyTimestampMode(t, mc, c)

		ddl := `CREATE TABLE t (
    ts1 TIMESTAMP NOT NULL,
    ts2 TIMESTAMP NOT NULL
)`
		// Run on container directly — runOnBoth would split sql_mode.
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err != nil {
			t.Errorf("oracle CREATE TABLE failed: %v", err)
		}
		results, err := c.Exec(ddl+";", nil)
		if err != nil {
			t.Errorf("omni parse error: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni exec error: %v", r.Error)
			}
		}

		// Oracle: SHOW CREATE TABLE on container.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if !strings.Contains(lo, "ts1") || strings.Contains(lo, "current_timestamp") {
			t.Errorf("oracle: TIMESTAMP unexpectedly promoted\n%s", create)
		}
		// Neither timestamp acquires an implicit ON UPDATE expression.
		if strings.Count(lo, "on update current_timestamp") != 0 {
			t.Errorf("oracle: unexpected implicit ON UPDATE\n%s", create)
		}

		// The catalog keeps both declarations without implicit defaults.
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
		// Expected omni (default-mode behavior): no promotion.
		if ts1.Nullable || ts2.Nullable || ts1.Default != nil || ts1.OnUpdate != "" {
			t.Errorf("catalog: explicit NOT NULL timestamps changed: %+v %+v", ts1, ts2)
		}
		if ts2.Default != nil && strings.Contains(strings.ToLower(*ts2.Default), "current_timestamp") {
			t.Errorf("omni: ts2 should NEVER auto-promote to CURRENT_TIMESTAMP, got %q", *ts2.Default)
		}
		if ts2.OnUpdate != "" && strings.Contains(strings.ToLower(ts2.OnUpdate), "current_timestamp") {
			t.Errorf("omni: ts2 should NEVER auto-promote ON UPDATE, got %q", ts2.OnUpdate)
		}
		assertScenarioReadback(t, mc, c, "t")
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

	// --- 3.8 Second TIMESTAMP under explicit_defaults_for_timestamp=OFF -----
	t.Run("3_8_second_timestamp_legacy_mode", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		assertReadonlyTimestampMode(t, mc, c)

		ddl := "CREATE TABLE t (ts1 TIMESTAMP, ts2 TIMESTAMP)"
		// Run on container.
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err != nil {
			t.Errorf("oracle CREATE TABLE failed: %v", err)
		}
		// Run the same declaration on the catalog after the rejected setting.
		results, err := c.Exec(ddl+";", nil)
		if err != nil {
			t.Errorf("omni parse error: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni exec error: %v", r.Error)
			}
		}

		// The native server leaves both columns nullable with NULL defaults.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if strings.Contains(lo, "default current_timestamp") {
			t.Errorf("oracle: ts1 unexpectedly promoted\n%s", create)
		}
		if strings.Contains(lo, "0000-00-00 00:00:00") {
			t.Errorf("oracle: ts2 unexpectedly acquired a zero default\n%s", create)
		}

		// Verify both modeled columns and the complete native readback.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing after legacy-mode CREATE TABLE")
			return
		}
		ts1 := tbl.GetColumn("ts1")
		ts2 := tbl.GetColumn("ts2")
		if ts1 == nil || ts2 == nil {
			t.Errorf("omni: ts1/ts2 columns missing")
			return
		}
		if !ts1.Nullable || !ts2.Nullable || ts1.Default != nil || ts2.Default != nil || ts1.OnUpdate != "" || ts2.OnUpdate != "" {
			t.Errorf("catalog: nullable timestamps changed: %+v %+v", ts1, ts2)
		}
		assertScenarioReadback(t, mc, c, "t")
	})
}

// TiDB exposes this variable as read-only. Both values must reject without
// enabling MySQL's old first-timestamp promotion behavior.
func assertReadonlyTimestampMode(t *testing.T, mc *mysqlContainer, c *Catalog) {
	t.Helper()
	stmts := []string{
		"SET SESSION explicit_defaults_for_timestamp=0",
		"SET SESSION explicit_defaults_for_timestamp=1",
	}
	for _, s := range stmts {
		_, err := mc.db.ExecContext(mc.ctx, s)
		var native *gomysql.MySQLError
		if !errors.As(err, &native) || native.Number != 1238 {
			t.Fatalf("native read-only timestamp setting %q: %v", s, err)
		}
		results, err := c.Exec(s, nil)
		if err != nil || len(results) != 1 {
			t.Fatalf("catalog timestamp setting: %v %v", err, results)
		}
		if modeled, ok := results[0].Error.(*Error); !ok || modeled.Code != 1238 {
			t.Fatalf("catalog read-only timestamp setting: %v", results[0].Error)
		}
	}
}
