package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C14 covers section C14 (CHECK-constraint behavior) re-baselined
// to TiDB v8.5.5. The fixtures were forked from mysql/catalog, where CHECK
// constraints are retained and enforced. TiDB v8.5.5 ships with
// tidb_enable_check_constraint OFF by default, so a CHECK is parsed and then
// silently discarded: it never reaches information_schema.CHECK_CONSTRAINTS,
// TABLE_CONSTRAINTS has no CHECK row (and no ENFORCED column at all, unlike
// MySQL 8.0), SHOW CREATE TABLE omits it, and the predicate is never
// evaluated. omni's catalog discards CHECK constraints identically, so each
// subtest asserts that neither engine retains or enforces them.
func TestScenario_C14(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 14.1 CHECK parsed then discarded under default ------------------
	t.Run("14_1_check_discarded_under_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, CHECK (a > 0));`)

		// TiDB: no CHECK_CONSTRAINTS row, and SHOW CREATE omits the CHECK.
		var n int
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'`, &n)
		assertIntEq(t, "oracle CHECK_CONSTRAINTS rows", n, 0)

		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(create), "CHECK") {
			t.Errorf("oracle: SHOW CREATE should not retain CHECK under default: %s", create)
		}

		// omni discards the CHECK too.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		assertIntEq(t, "omni CHECK count", len(omniCheckNames(tbl)), 0)
	})

	// --- 14.2 ALTER ... ALTER CHECK ENFORCED toggle is a no-op -----------
	t.Run("14_2_alter_check_enforcement_noop", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// The CHECK is discarded at CREATE; ALTER ... ALTER CHECK c_pos
		// NOT ENFORCED / ENFORCED are accepted as no-ops by both engines.
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0));
ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED;`)

		var n int
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'`, &n)
		assertIntEq(t, "oracle CHECK_CONSTRAINTS rows after NOT ENFORCED", n, 0)

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		assertIntEq(t, "omni CHECK count after NOT ENFORCED", len(omniCheckNames(tbl)), 0)

		runOnBoth(t, mc, c, `ALTER TABLE t ALTER CHECK c_pos ENFORCED;`)

		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'`, &n)
		assertIntEq(t, "oracle CHECK_CONSTRAINTS rows after re-ENFORCED", n, 0)

		tbl2 := c.GetDatabase("testdb").GetTable("t")
		assertIntEq(t, "omni CHECK count after re-ENFORCED", len(omniCheckNames(tbl2)), 0)
	})

	// --- 14.3 STORED generated column retained, CHECK discarded ----------
	t.Run("14_3_check_with_stored_generated_col", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
    a INT,
    g INT AS (a * 2) STORED,
    CONSTRAINT c_g_nonneg CHECK (g >= 0)
);`)

		// The generated column is retained; the CHECK is not.
		var n int
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='c_g_nonneg'`, &n)
		assertIntEq(t, "oracle CHECK_CONSTRAINTS rows", n, 0)

		// With the CHECK discarded, a=-1 (g=-2) is accepted — no enforcement.
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t (a) VALUES (5)`); err != nil {
			t.Errorf("oracle: expected INSERT a=5 to succeed, got %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t (a) VALUES (-1)`); err != nil {
			t.Errorf("oracle: expected INSERT a=-1 to succeed (CHECK discarded), got %v", err)
		}

		// omni: no CHECK registered, but the generated column g must exist.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		assertIntEq(t, "omni CHECK count", len(omniCheckNames(tbl)), 0)
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

	// --- 14.4 Constructs MySQL forbids in CHECK are accepted -------------
	t.Run("14_4_check_forbidden_constructs_accepted", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// With the CHECK body never validated, constructs MySQL rejects
		// (subquery, NOW(), user variable, RAND()) are accepted and
		// discarded by both engines.
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE other (id INT)`); err != nil {
			t.Errorf("oracle setup: CREATE other failed: %v", err)
		}
		cases := []struct{ label, ddl string }{
			{"subquery CHECK", `CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)))`},
			{"NOW() CHECK", `CREATE TABLE t3 (a INT, CHECK (a < NOW()))`},
			{"user-var CHECK", `CREATE TABLE t4 (a INT, CHECK (a < @max))`},
			{"RAND() CHECK", `CREATE TABLE t5 (a INT, CHECK (a < RAND()))`},
		}
		for _, tc := range cases {
			if _, err := mc.db.ExecContext(mc.ctx, tc.ddl); err != nil {
				t.Errorf("oracle: expected %s to be accepted (CHECK discarded), got %v", tc.label, err)
			}
			c14AssertOmniAccepts(t, c, tc.label, tc.ddl+";")
		}
	})
}

// --- section-local helpers ------------------------------------------------

// c14AssertOmniAccepts executes the given DDL against the omni catalog and
// reports a test error if parsing or execution returns an error.
func c14AssertOmniAccepts(t *testing.T, c *Catalog, label, ddl string) {
	t.Helper()
	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Errorf("omni: expected %s to be accepted, got parse error: %v", label, err)
		return
	}
	for _, r := range results {
		if r.Error != nil {
			t.Errorf("omni: expected %s to be accepted, got exec error: %v", label, r.Error)
			return
		}
	}
}
