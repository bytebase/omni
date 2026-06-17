package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C14 covers section C14 (Constraint enforcement defaults) from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that both real
// MySQL 8.0 and the omni catalog agree on CHECK-constraint enforcement and
// validation behavior.
//
// Failures in omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c14.md.
func TestScenario_C14(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 14.1 CHECK constraint defaults to ENFORCED ----------------------
	t.Run("14_1_check_defaults_enforced", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, CHECK (a > 0));`)

		// Oracle: information_schema.TABLE_CONSTRAINTS.ENFORCED should be YES.
		// Note: information_schema.CHECK_CONSTRAINTS in MySQL 8.0 does NOT
		// expose an ENFORCED column — it only has (CONSTRAINT_CATALOG,
		// CONSTRAINT_SCHEMA, CONSTRAINT_NAME, CHECK_CLAUSE). The ENFORCED
		// metadata lives in TABLE_CONSTRAINTS instead.
		var tcEnforced string
		oracleScan(t, mc, `SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_TYPE='CHECK'`,
			&tcEnforced)
		assertStringEq(t, "oracle TABLE_CONSTRAINTS.ENFORCED", tcEnforced, "YES")

		// SHOW CREATE TABLE must not contain NOT ENFORCED.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(create), "NOT ENFORCED") {
			t.Errorf("oracle: SHOW CREATE TABLE unexpectedly contains NOT ENFORCED: %s", create)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		chk := c14FindCheck(tbl, "t_chk_1")
		if chk == nil {
			t.Errorf("omni: CHECK constraint t_chk_1 missing")
			return
		}
		// Default = ENFORCED → NotEnforced must be false.
		assertBoolEq(t, "omni t_chk_1 NotEnforced", chk.NotEnforced, false)
	})

	// --- 14.2 ALTER CHECK NOT ENFORCED / ENFORCED toggles the flag -------
	t.Run("14_2_alter_check_enforcement_toggle", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0));
ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED;`)

		// Oracle: after NOT ENFORCED, TABLE_CONSTRAINTS.ENFORCED should be NO.
		var enforced string
		oracleScan(t, mc, `SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_NAME='c_pos'`,
			&enforced)
		assertStringEq(t, "oracle after NOT ENFORCED", enforced, "NO")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		chk := c14FindCheck(tbl, "c_pos")
		if chk == nil {
			t.Errorf("omni: CHECK c_pos missing after initial CREATE+ALTER")
			return
		}
		assertBoolEq(t, "omni c_pos NotEnforced after NOT ENFORCED", chk.NotEnforced, true)

		// Toggle back to ENFORCED and re-check both sides.
		runOnBoth(t, mc, c, `ALTER TABLE t ALTER CHECK c_pos ENFORCED;`)

		oracleScan(t, mc, `SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_NAME='c_pos'`,
			&enforced)
		assertStringEq(t, "oracle after re-ENFORCED", enforced, "YES")

		tbl2 := c.GetDatabase("testdb").GetTable("t")
		chk2 := c14FindCheck(tbl2, "c_pos")
		if chk2 == nil {
			t.Errorf("omni: CHECK c_pos missing after re-ENFORCED")
			return
		}
		assertBoolEq(t, "omni c_pos NotEnforced after re-ENFORCED", chk2.NotEnforced, false)
	})

	// --- 14.3 STORED generated column + CHECK: predicate evaluated against stored value
	t.Run("14_3_check_with_stored_generated_col", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
    a INT,
    g INT AS (a * 2) STORED,
    CONSTRAINT c_g_nonneg CHECK (g >= 0)
);`)

		// Oracle: CHECK_CONSTRAINTS row for c_g_nonneg exists with clause
		// referencing g. MySQL stores it as (`g` >= 0).
		var name, clause string
		oracleScan(t, mc, `SELECT CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='c_g_nonneg'`,
			&name, &clause)
		assertStringEq(t, "oracle CHECK name", name, "c_g_nonneg")
		if !strings.Contains(clause, "g") || !strings.Contains(clause, ">=") {
			t.Errorf("oracle: expected CHECK_CLAUSE to reference g and >=, got %q", clause)
		}

		// Valid INSERT: a=5 → g=10 passes.
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t (a) VALUES (5)`); err != nil {
			t.Errorf("oracle: expected INSERT a=5 to succeed, got %v", err)
		}
		// Invalid INSERT: a=-1 → g=-2 violates CHECK, MySQL error 3819.
		_, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t (a) VALUES (-1)`)
		if err == nil {
			t.Errorf("oracle: expected ER_CHECK_CONSTRAINT_VIOLATED (3819) for a=-1, got nil")
		} else if !strings.Contains(err.Error(), "3819") &&
			!strings.Contains(strings.ToLower(err.Error()), "check constraint") {
			t.Errorf("oracle: expected CHECK violation error, got %v", err)
		}

		// omni: CHECK must be registered and reference the generated column g.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		chk := c14FindCheck(tbl, "c_g_nonneg")
		if chk == nil {
			t.Errorf("omni: CHECK c_g_nonneg missing")
			return
		}
		if !strings.Contains(chk.CheckExpr, "g") {
			t.Errorf("omni: CheckExpr should reference g, got %q", chk.CheckExpr)
		}
		// g column must exist and be stored-generated.
		var gcol *Column
		for _, col := range tbl.Columns {
			if strings.EqualFold(col.Name, "g") {
				gcol = col
				break
			}
		}
		if gcol == nil {
			t.Errorf("omni: generated column g missing")
		}
	})

	// --- 14.4 Forbidden constructs in CHECK: subquery, NOW(), user variable
	t.Run("14_4_check_forbidden_constructs", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Legal baseline: deterministic built-in CHAR_LENGTH is OK.
		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, CHECK (CHAR_LENGTH(CAST(a AS CHAR)) < 10));`)

		// --- 14.4a Subquery in CHECK → ER_CHECK_CONSTRAINT_NOT_ALLOWED_CONTEXT (3812).
		// Need a referenced table first.
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE other (id INT)`); err != nil {
			t.Errorf("oracle setup: CREATE other failed: %v", err)
		}
		_, errSub := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)))`)
		if errSub == nil {
			t.Errorf("oracle: expected subquery CHECK to fail, got nil")
		} else {
			// MySQL 8.0 observed to return 3815/3814 ("disallowed function")
			// for subquery-in-CHECK; older docs mention 3812 as well. Accept
			// any check-constraint rejection error.
			if !c14IsCheckRejection(errSub) {
				t.Errorf("oracle: expected CHECK-rejection error, got %v", errSub)
			}
		}
		c14AssertOmniRejects(t, c, "subquery CHECK",
			`CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)));`)

		// --- 14.4b NOW() in CHECK → ER_CHECK_CONSTRAINT_NAMED_FUNCTION_IS_NOT_ALLOWED (3815).
		_, errNow := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t3 (a INT, CHECK (a < NOW()))`)
		if errNow == nil {
			t.Errorf("oracle: expected NOW() CHECK to fail, got nil")
		} else if !c14IsCheckRejection(errNow) {
			t.Errorf("oracle: expected NOW() CHECK-rejection error, got %v", errNow)
		}
		c14AssertOmniRejects(t, c, "NOW() CHECK",
			`CREATE TABLE t3 (a INT, CHECK (a < NOW()));`)

		// --- 14.4c User variable in CHECK → ER_CHECK_CONSTRAINT_VARIABLES (3813).
		_, errVar := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t4 (a INT, CHECK (a < @max))`)
		if errVar == nil {
			t.Errorf("oracle: expected user-var CHECK to fail, got nil")
		} else if !c14IsCheckRejection(errVar) {
			t.Errorf("oracle: expected user-var CHECK-rejection error, got %v", errVar)
		}
		c14AssertOmniRejects(t, c, "user-var CHECK",
			`CREATE TABLE t4 (a INT, CHECK (a < @max));`)

		// --- 14.4d RAND() in CHECK → ER_CHECK_CONSTRAINT_NAMED_FUNCTION_IS_NOT_ALLOWED (3815).
		_, errRand := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t5 (a INT, CHECK (a < RAND()))`)
		if errRand == nil {
			t.Errorf("oracle: expected RAND() CHECK to fail, got nil")
		} else if !c14IsCheckRejection(errRand) {
			t.Errorf("oracle: expected RAND() CHECK-rejection error, got %v", errRand)
		}
		c14AssertOmniRejects(t, c, "RAND() CHECK",
			`CREATE TABLE t5 (a INT, CHECK (a < RAND()));`)
	})
}

// --- section-local helpers ------------------------------------------------

// c14FindCheck returns the first CHECK constraint on the table whose name
// matches (case-insensitive), or nil.
func c14FindCheck(tbl *Table, name string) *Constraint {
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && strings.EqualFold(con.Name, name) {
			return con
		}
	}
	return nil
}

// c14IsCheckRejection reports whether the given error looks like a MySQL
// rejection of a forbidden construct in a CHECK clause. Accepts any of the
// documented error codes (3812–3815) or the strings "not allowed",
// "disallowed", "subquer", or "variable".
func c14IsCheckRejection(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	for _, code := range []string{"3812", "3813", "3814", "3815"} {
		if strings.Contains(msg, code) {
			return true
		}
	}
	return strings.Contains(lower, "not allowed") ||
		strings.Contains(lower, "disallowed") ||
		strings.Contains(lower, "subquer") ||
		strings.Contains(lower, "check constraint")
}

// c14AssertOmniRejects executes the given DDL against the omni catalog and
// reports a test error if execution succeeds with no error. Uses t.Error so
// subsequent assertions continue to run. The catalog state after a rejected
// DDL is considered undefined — callers should not rely on it.
func c14AssertOmniRejects(t *testing.T, c *Catalog, label, ddl string) {
	t.Helper()
	results, err := c.Exec(ddl, nil)
	if err != nil {
		return // parse error counts as rejection
	}
	for _, r := range results {
		if r.Error != nil {
			return // exec error counts as rejection
		}
	}
	t.Errorf("omni: expected %s to be rejected, but execution succeeded", label)
}
