package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C3 covers section C3 (Nullability & default promotion) from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that both real
// MySQL 8.0 and the omni catalog agree on the implicit nullability/default
// rules MySQL applies during CREATE TABLE.
//
// Failures in omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c3.md.
func TestScenario_C3(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 3.1 First TIMESTAMP NOT NULL auto-promotes (legacy mode) -------
	t.Run("3_1_first_timestamp_promotion", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Force legacy mode on the container so the first-TS promotion fires.
		// Also clear sql_mode so MySQL accepts the implicit '0000-00-00' zero
		// default for the second TIMESTAMP column. omni does not track these
		// session vars; we still verify omni's static behavior.
		setLegacyTimestampMode(t, mc)
		defer restoreTimestampMode(t, mc)

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
		if !strings.Contains(lo, "ts1") ||
			!strings.Contains(lo, "default current_timestamp") ||
			!strings.Contains(lo, "on update current_timestamp") {
			t.Errorf("oracle: expected ts1 promoted to CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n%s", create)
		}
		// ts2 should NOT carry an ON UPDATE CURRENT_TIMESTAMP — only one occurrence overall.
		if strings.Count(lo, "on update current_timestamp") != 1 {
			t.Errorf("oracle: expected exactly one ON UPDATE CURRENT_TIMESTAMP (ts1 only)\n%s", create)
		}

		// omni: this is omni's "static" view. omni does not track
		// explicit_defaults_for_timestamp, so the expectation here documents
		// what omni *currently* does. We accept either:
		//   (a) omni promotes the first TIMESTAMP NOT NULL  (matches legacy)
		//   (b) omni leaves both ts1 and ts2 alone        (matches default mode)
		// Anything else is a bug. Document the asymmetry rather than failing.
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
		if ts1.Default != nil && strings.Contains(strings.ToLower(*ts1.Default), "current_timestamp") {
			// Acceptable: omni mirrors legacy mode.
			t.Logf("omni: ts1 has CURRENT_TIMESTAMP default — promotion happens unconditionally")
		}
		if ts2.Default != nil && strings.Contains(strings.ToLower(*ts2.Default), "current_timestamp") {
			t.Errorf("omni: ts2 should NEVER auto-promote to CURRENT_TIMESTAMP, got %q", *ts2.Default)
		}
		if ts2.OnUpdate != "" && strings.Contains(strings.ToLower(ts2.OnUpdate), "current_timestamp") {
			t.Errorf("omni: ts2 should NEVER auto-promote ON UPDATE, got %q", ts2.OnUpdate)
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

	// --- 3.8 Second TIMESTAMP under explicit_defaults_for_timestamp=OFF -----
	t.Run("3_8_second_timestamp_legacy_mode", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Container side: legacy mode + permissive sql_mode.
		setLegacyTimestampMode(t, mc)
		defer restoreTimestampMode(t, mc)

		ddl := "CREATE TABLE t (ts1 TIMESTAMP, ts2 TIMESTAMP)"
		// Run on container.
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err != nil {
			t.Errorf("oracle CREATE TABLE failed: %v", err)
		}
		// Run on omni separately — omni does not track the session var so
		// its observed state is whatever omni's default branch produces.
		results, err := c.Exec(ddl+";", nil)
		if err != nil {
			t.Errorf("omni parse error: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni exec error: %v", r.Error)
			}
		}

		// Oracle: ts1 promoted, ts2 NOT NULL DEFAULT zero literal.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if !strings.Contains(lo, "default current_timestamp") {
			t.Errorf("oracle: legacy-mode ts1 should have DEFAULT CURRENT_TIMESTAMP\n%s", create)
		}
		if !strings.Contains(lo, "0000-00-00 00:00:00") {
			t.Errorf("oracle: legacy-mode ts2 should have DEFAULT '0000-00-00 00:00:00'\n%s", create)
		}

		// omni asymmetry: omni does not track explicit_defaults_for_timestamp,
		// so we cannot expect it to mirror the legacy transform. Document
		// what omni does for posterity. Anything goes here EXCEPT crashing.
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
		t.Logf("omni asymmetry (no session var tracking): ts1.Nullable=%v ts1.Default=%v ts2.Nullable=%v ts2.Default=%v",
			ts1.Nullable, derefStr(ts1.Default), ts2.Nullable, derefStr(ts2.Default))
	})
}

func derefStr(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

// setLegacyTimestampMode flips the container into the deprecated
// pre-5.6 TIMESTAMP semantics (explicit_defaults_for_timestamp=OFF) and
// drops the strict-zero-date sql_mode flags so '0000-00-00 00:00:00'
// defaults are accepted. This is required for C3.1 and C3.8 which exercise
// the legacy TIMESTAMP promotion path inside `promote_first_timestamp_column`.
func setLegacyTimestampMode(t *testing.T, mc *mysqlContainer) {
	t.Helper()
	stmts := []string{
		"SET SESSION explicit_defaults_for_timestamp=0",
		"SET SESSION sql_mode=''",
	}
	for _, s := range stmts {
		if _, err := mc.db.ExecContext(mc.ctx, s); err != nil {
			t.Fatalf("setLegacyTimestampMode %q: %v", s, err)
		}
	}
}

// restoreTimestampMode reverts the session vars touched by
// setLegacyTimestampMode back to MySQL 8.0 defaults.
func restoreTimestampMode(t *testing.T, mc *mysqlContainer) {
	t.Helper()
	stmts := []string{
		"SET SESSION explicit_defaults_for_timestamp=1",
		"SET SESSION sql_mode=DEFAULT",
	}
	for _, s := range stmts {
		_, _ = mc.db.ExecContext(mc.ctx, s)
	}
}
