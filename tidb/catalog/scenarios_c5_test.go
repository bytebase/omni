package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C5 covers section C5 (Constraint defaults) from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. It checks FK and CHECK
// constraint defaults: ON DELETE/ON UPDATE/MATCH defaults, FK
// SET DEFAULT InnoDB rejection, FK column type compatibility,
// FK on virtual gcol rejection, CHECK ENFORCED default,
// column-level vs table-level CHECK equivalence, column-level CHECK
// cross-column reference rejection.
//
// Failures in omni assertions are NOT proof failures — they are
// recorded in mysql/catalog/scenarios_bug_queue/c5.md.
func TestScenario_C5(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 5.1 FK ON DELETE default — RESTRICT internally / NO ACTION reported ---
	t.Run("5_1_fk_on_delete_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id));`
		runOnBoth(t, mc, c, ddl)

		// Oracle: REFERENTIAL_CONSTRAINTS rules.
		var delRule, updRule, matchOpt string
		oracleScan(t, mc, `SELECT DELETE_RULE, UPDATE_RULE, MATCH_OPTION
            FROM information_schema.REFERENTIAL_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'`,
			&delRule, &updRule, &matchOpt)
		assertStringEq(t, "oracle DELETE_RULE", delRule, "NO ACTION")
		assertStringEq(t, "oracle UPDATE_RULE", updRule, "NO ACTION")
		assertStringEq(t, "oracle MATCH_OPTION", matchOpt, "NONE")

		// SHOW CREATE TABLE should not contain ON DELETE / ON UPDATE / MATCH.
		create := oracleShow(t, mc, "SHOW CREATE TABLE c")
		if strings.Contains(strings.ToUpper(create), "ON DELETE") {
			t.Errorf("oracle: SHOW CREATE TABLE should omit ON DELETE: %s", create)
		}
		if strings.Contains(strings.ToUpper(create), "ON UPDATE") {
			t.Errorf("oracle: SHOW CREATE TABLE should omit ON UPDATE: %s", create)
		}
		if strings.Contains(strings.ToUpper(create), "MATCH") {
			t.Errorf("oracle: SHOW CREATE TABLE should omit MATCH: %s", create)
		}

		// omni: FK should have default OnDelete/OnUpdate (empty or NO ACTION /
		// RESTRICT) and the deparsed table should not render those clauses.
		tbl := c.GetDatabase("testdb").GetTable("c")
		if tbl == nil {
			t.Errorf("omni: table c missing")
			return
		}
		fk := c5FirstFK(tbl)
		if fk == nil {
			t.Errorf("omni: no FK constraint on c")
			return
		}
		if !c5IsFKDefault(fk.OnDelete) {
			t.Errorf("omni: FK OnDelete should be default, got %q", fk.OnDelete)
		}
		if !c5IsFKDefault(fk.OnUpdate) {
			t.Errorf("omni: FK OnUpdate should be default, got %q", fk.OnUpdate)
		}
		omniCreate := c5OmniShowCreate(t, c, "c")
		if strings.Contains(strings.ToUpper(omniCreate), "ON DELETE") {
			t.Errorf("omni: deparse should omit ON DELETE: %s", omniCreate)
		}
		if strings.Contains(strings.ToUpper(omniCreate), "ON UPDATE") {
			t.Errorf("omni: deparse should omit ON UPDATE: %s", omniCreate)
		}
	})

	// --- 5.2 FK ON DELETE SET NULL on a NOT NULL column errors --------------
	t.Run("5_2_fk_set_null_requires_nullable", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		setup := `CREATE TABLE p (id INT PRIMARY KEY);`
		if _, err := mc.db.ExecContext(mc.ctx, setup); err != nil {
			t.Errorf("oracle setup: %v", err)
		}
		_, _ = c.Exec(setup, nil)

		bad := `CREATE TABLE c (
            a INT NOT NULL,
            FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL
        )`
		_, mysqlErr := mc.db.ExecContext(mc.ctx, bad)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_FK_COLUMN_NOT_NULL, got nil")
		}

		results, err := c.Exec(bad+";", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects SET NULL on NOT NULL column", omniErrored, true)
	})

	// --- 5.3 FK MATCH default rendered as NONE in information_schema --------
	t.Run("5_3_fk_match_default_none", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id));`)

		var matchOpt string
		oracleScan(t, mc, `SELECT MATCH_OPTION FROM information_schema.REFERENTIAL_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'`, &matchOpt)
		assertStringEq(t, "oracle MATCH_OPTION", matchOpt, "NONE")

		tbl := c.GetDatabase("testdb").GetTable("c")
		if tbl == nil {
			t.Errorf("omni: table c missing")
			return
		}
		fk := c5FirstFK(tbl)
		if fk == nil {
			t.Errorf("omni: no FK on c")
			return
		}
		// omni MatchType should be empty / NONE / SIMPLE — anything matching default.
		if fk.MatchType != "" && !strings.EqualFold(fk.MatchType, "NONE") &&
			!strings.EqualFold(fk.MatchType, "SIMPLE") {
			t.Errorf("omni: FK MatchType should be default (empty/NONE/SIMPLE), got %q", fk.MatchType)
		}
	})

	// --- 5.4 FK ON UPDATE default independent of ON DELETE -----------------
	t.Run("5_4_fk_on_update_independent", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id) ON DELETE CASCADE);`
		runOnBoth(t, mc, c, ddl)

		var delRule, updRule string
		oracleScan(t, mc, `SELECT DELETE_RULE, UPDATE_RULE
            FROM information_schema.REFERENTIAL_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'`,
			&delRule, &updRule)
		assertStringEq(t, "oracle DELETE_RULE", delRule, "CASCADE")
		assertStringEq(t, "oracle UPDATE_RULE", updRule, "NO ACTION")

		// SHOW CREATE renders ON DELETE CASCADE only.
		create := oracleShow(t, mc, "SHOW CREATE TABLE c")
		if !strings.Contains(strings.ToUpper(create), "ON DELETE CASCADE") {
			t.Errorf("oracle: expected ON DELETE CASCADE, got %s", create)
		}
		if strings.Contains(strings.ToUpper(create), "ON UPDATE") {
			t.Errorf("oracle: expected no ON UPDATE clause, got %s", create)
		}

		tbl := c.GetDatabase("testdb").GetTable("c")
		if tbl == nil {
			t.Errorf("omni: table c missing")
			return
		}
		fk := c5FirstFK(tbl)
		if fk == nil {
			t.Errorf("omni: no FK on c")
			return
		}
		if !strings.EqualFold(fk.OnDelete, "CASCADE") {
			t.Errorf("omni: FK OnDelete should be CASCADE, got %q", fk.OnDelete)
		}
		if !c5IsFKDefault(fk.OnUpdate) {
			t.Errorf("omni: FK OnUpdate should be default (empty/NO ACTION/RESTRICT), got %q", fk.OnUpdate)
		}
		omniCreate := c5OmniShowCreate(t, c, "c")
		if !strings.Contains(strings.ToUpper(omniCreate), "ON DELETE CASCADE") {
			t.Errorf("omni: deparse missing ON DELETE CASCADE: %s", omniCreate)
		}
		if strings.Contains(strings.ToUpper(omniCreate), "ON UPDATE") {
			t.Errorf("omni: deparse should omit ON UPDATE: %s", omniCreate)
		}
	})

	// --- 5.5 FK SET DEFAULT rejected by InnoDB ------------------------------
	t.Run("5_5_fk_set_default_innodb_limitation", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		setup := `CREATE TABLE p (id INT PRIMARY KEY);`
		if _, err := mc.db.ExecContext(mc.ctx, setup); err != nil {
			t.Errorf("oracle setup: %v", err)
		}
		_, _ = c.Exec(setup, nil)

		bad := `CREATE TABLE c (
            a INT DEFAULT 0,
            FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT
        )`

		// MySQL/InnoDB returns an error or warning for SET DEFAULT. Capture
		// whichever is observed rather than asserting the precise behavior so
		// the test reflects the real oracle.
		_, mysqlErr := mc.db.ExecContext(mc.ctx, bad)
		oracleRejected := mysqlErr != nil

		results, err := c.Exec(bad+";", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}

		// Both should agree.
		assertBoolEq(t, "omni FK SET DEFAULT matches MySQL rejection",
			omniErrored, oracleRejected)
	})

	// --- 5.6 FK column type/size/sign must match parent --------------------
	t.Run("5_6_fk_column_type_compat", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		setup := `CREATE TABLE p (id BIGINT PRIMARY KEY);`
		if _, err := mc.db.ExecContext(mc.ctx, setup); err != nil {
			t.Errorf("oracle setup: %v", err)
		}
		_, _ = c.Exec(setup, nil)

		bad := `CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))`
		_, mysqlErr := mc.db.ExecContext(mc.ctx, bad)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_FK_INCOMPATIBLE_COLUMNS (3780), got nil")
		}

		results, err := c.Exec(bad+";", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects FK with mismatched column types", omniErrored, true)
	})

	// --- 5.7 FK on a VIRTUAL generated column rejected ---------------------
	t.Run("5_7_fk_on_virtual_gcol_rejected", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		setup := `CREATE TABLE p (id INT PRIMARY KEY);`
		if _, err := mc.db.ExecContext(mc.ctx, setup); err != nil {
			t.Errorf("oracle setup: %v", err)
		}
		_, _ = c.Exec(setup, nil)

		bad := `CREATE TABLE c (
            a INT,
            b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
            FOREIGN KEY (b) REFERENCES p(id)
        )`
		_, mysqlErr := mc.db.ExecContext(mc.ctx, bad)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_FK_CANNOT_USE_VIRTUAL_COLUMN (3104), got nil")
		}

		results, err := c.Exec(bad+";", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects FK on VIRTUAL gcol", omniErrored, true)
	})

	// --- 5.8 CHECK discarded under TiDB default ----------------------------
	t.Run("5_8_check_discarded_under_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
            a INT,
            CONSTRAINT chk_pos CHECK (a > 0)
        );`)

		// tidb_enable_check_constraint is OFF by default in v8.5.5: the CHECK
		// is discarded, so there is no CHECK_CONSTRAINTS row and SHOW CREATE
		// omits it on both engines.
		var n int
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='chk_pos'`, &n)
		assertIntEq(t, "oracle CHECK_CONSTRAINTS rows", n, 0)

		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(create), "CHECK") {
			t.Errorf("oracle: SHOW CREATE should not retain CHECK under default: %s", create)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		assertIntEq(t, "omni CHECK count", len(omniCheckNames(tbl)), 0)

		omniCreate := c5OmniShowCreate(t, c, "t")
		if strings.Contains(strings.ToUpper(omniCreate), "CHECK") {
			t.Errorf("omni: deparse should not retain CHECK under default: %s", omniCreate)
		}
	})

	// --- 5.9 column-level and table-level CHECK both discarded -------------
	t.Run("5_9_check_column_and_table_level_discarded", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT CHECK (a > 0));
CREATE TABLE t2 (a INT, CHECK (a > 0));`)

		// Neither the column-level nor the table-level CHECK is retained.
		rows1 := oracleRows(t, mc, `SELECT CONSTRAINT_NAME, CHECK_CLAUSE
            FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't1_chk_%'`)
		rows2 := oracleRows(t, mc, `SELECT CONSTRAINT_NAME, CHECK_CLAUSE
            FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't2_chk_%'`)
		assertIntEq(t, "oracle t1 CHECK rows", len(rows1), 0)
		assertIntEq(t, "oracle t2 CHECK rows", len(rows2), 0)

		t1 := c.GetDatabase("testdb").GetTable("t1")
		t2 := c.GetDatabase("testdb").GetTable("t2")
		if t1 == nil || t2 == nil {
			t.Errorf("omni: t1 or t2 missing")
			return
		}
		assertIntEq(t, "omni t1 check count", len(omniCheckNames(t1)), 0)
		assertIntEq(t, "omni t2 check count", len(omniCheckNames(t2)), 0)
	})

	// --- 5.10 column-level CHECK with cross-column ref discarded -----------
	t.Run("5_10_column_check_cross_ref_discarded", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// MySQL rejects a column-level CHECK referencing another column
		// (ER_CHECK_CONSTRAINT_REFERS, 3823). Under the default OFF setting
		// the CHECK is discarded without validation, so both engines accept.
		bad := `CREATE TABLE t (a INT CHECK (a > b), b INT)`
		if _, oracleErr := mc.db.ExecContext(mc.ctx, bad); oracleErr != nil {
			t.Errorf("oracle: expected success (CHECK discarded), got %v", oracleErr)
		}

		results, err := c.Exec(bad+";", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni accepts column-level CHECK referencing other column (discarded)", omniErrored, false)

		// A table-level CHECK with the same expression is likewise accepted.
		scenarioReset(t, mc)
		c2 := scenarioNewCatalog(t)
		good := `CREATE TABLE t (a INT, b INT, CHECK (a > b));`
		runOnBoth(t, mc, c2, good)
	})
}

// --- C5 section helpers ---------------------------------------------------

// c5FirstFK returns the first FK constraint of the table or nil.
func c5FirstFK(tbl *Table) *Constraint {
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			return con
		}
	}
	return nil
}

// c5IsFKDefault returns true when the FK action string is one of the values
// MySQL treats as the implicit default (empty, NO ACTION, RESTRICT).
func c5IsFKDefault(action string) bool {
	upper := strings.ToUpper(strings.TrimSpace(action))
	return upper == "" || upper == "NO ACTION" || upper == "RESTRICT"
}

// c5OmniShowCreate returns omni's SHOW CREATE TABLE rendering for the named
// table in testdb, or the empty string if the table is missing.
func c5OmniShowCreate(t *testing.T, c *Catalog, table string) string {
	t.Helper()
	return c.ShowCreateTable("testdb", table)
}
