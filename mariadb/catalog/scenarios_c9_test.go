package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C9 covers Section C9 "Generated column defaults" from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL against both
// a real MySQL 8.0 container and the omni catalog, then asserts that both
// agree on the effective default for a given generated-column behavior.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c9.md.
func TestScenario_C9(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c9OracleExtra returns the EXTRA column for a given column in testdb.
	c9OracleExtra := func(t *testing.T, table, col string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			`SELECT IFNULL(EXTRA,'') FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+table+`' AND COLUMN_NAME='`+col+`'`,
			&s)
		return s
	}

	// c9OmniExec runs a multi-statement DDL on omni and returns (parseErr,
	// anyStmtErrored). Used by scenarios that expect omni to reject DDL.
	c9OmniExec := func(c *Catalog, ddl string) (bool, error) {
		results, err := c.Exec(ddl, nil)
		if err != nil {
			return true, err
		}
		for _, r := range results {
			if r.Error != nil {
				return true, r.Error
			}
		}
		return false, nil
	}

	// --- 9.1 Generated column storage defaults to VIRTUAL --------------------
	t.Run("9_1_gcol_storage_defaults_virtual", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a+1))`)

		extra := c9OracleExtra(t, "t", "b")
		if !strings.Contains(strings.ToUpper(extra), "VIRTUAL GENERATED") {
			t.Errorf("oracle EXTRA for b: got %q, want contains %q", extra, "VIRTUAL GENERATED")
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Error("omni: column b not found")
			return
		}
		if col.Generated == nil {
			t.Error("omni: column b should be generated")
			return
		}
		// Default storage is VIRTUAL (Stored == false).
		assertBoolEq(t, "omni b Generated.Stored", col.Generated.Stored, false)
	})

	// --- 9.2 FK on generated column — dual-agreement test -------------------
	//
	// SCENARIOS-mysql-implicit-behavior.md claims MySQL rejects FK where the
	// child column is a STORED generated column (`ER_FK_CANNOT_USE_VIRTUAL_COLUMN`).
	// Empirical oracle check: MySQL 8.0.45 ALLOWS FK on a child STORED gcol;
	// only VIRTUAL gcols are rejected on the child side. The scenario's
	// expected error is outdated. This test asserts oracle and omni agree
	// (both should accept) and documents the scenario discrepancy.
	t.Run("9_2_fk_on_stored_gcol_child", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		parent := `CREATE TABLE p (id INT PRIMARY KEY)`
		if _, err := mc.db.ExecContext(mc.ctx, parent); err != nil {
			t.Errorf("oracle setup parent: %v", err)
		}
		if _, err := c.Exec(parent+";", nil); err != nil {
			t.Errorf("omni setup parent: %v", err)
		}

		ddl := `CREATE TABLE c (
            a INT,
            b INT GENERATED ALWAYS AS (a+1) STORED,
            FOREIGN KEY (b) REFERENCES p(id)
        )`
		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		oracleAccepted := oracleErr == nil
		omniErrored, _ := c9OmniExec(c, ddl+";")
		omniAccepted := !omniErrored

		// Both systems must agree.
		assertBoolEq(t, "oracle vs omni agreement on FK on STORED gcol",
			omniAccepted, oracleAccepted)

		// Also verify the VIRTUAL gcol case is rejected by both (this is the
		// real rule — `ER_FK_CANNOT_USE_VIRTUAL_COLUMN`).
		scenarioReset(t, mc)
		c2 := scenarioNewCatalog(t)
		if _, err := mc.db.ExecContext(mc.ctx, parent); err != nil {
			t.Errorf("oracle setup parent (virtual case): %v", err)
		}
		if _, err := c2.Exec(parent+";", nil); err != nil {
			t.Errorf("omni setup parent (virtual case): %v", err)
		}
		badVirtual := `CREATE TABLE cv (
            a INT,
            b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
            FOREIGN KEY (b) REFERENCES p(id)
        )`
		_, oracleVirtErr := mc.db.ExecContext(mc.ctx, badVirtual)
		if oracleVirtErr == nil {
			t.Errorf("oracle: expected FK-on-VIRTUAL-gcol rejection, got nil error")
		}
		omniVirtErrored, _ := c9OmniExec(c2, badVirtual+";")
		assertBoolEq(t, "omni rejects FK on VIRTUAL gcol (child)", omniVirtErrored, true)
	})

	// --- 9.3 VIRTUAL gcol in PK rejected; secondary KEY allowed --------------
	t.Run("9_3_virtual_gcol_pk_rejected_secondary_allowed", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// t1: secondary KEY on virtual gcol should succeed.
		good := `CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, KEY (b))`
		if _, err := mc.db.ExecContext(mc.ctx, good); err != nil {
			t.Errorf("oracle: expected t1 to succeed, got: %v", err)
		}
		if omniErr, _ := c9OmniExec(c, good+";"); omniErr {
			t.Errorf("omni: expected t1 to succeed (VIRTUAL gcol as secondary KEY)")
		}
		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl != nil {
			found := false
			for _, idx := range tbl.Indexes {
				for _, col := range idx.Columns {
					if strings.EqualFold(col.Name, "b") {
						found = true
					}
				}
			}
			if !found {
				t.Errorf("omni: secondary index on b missing from t1")
			}
		} else {
			t.Errorf("omni: t1 table missing after CREATE")
		}

		// t2: PRIMARY KEY on virtual gcol should error.
		bad := `CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL PRIMARY KEY)`
		_, oracleErr := mc.db.ExecContext(mc.ctx, bad)
		if oracleErr == nil {
			t.Errorf("oracle: expected VIRTUAL gcol PK rejection, got nil error")
		}
		omniErrored, _ := c9OmniExec(c, bad+";")
		assertBoolEq(t, "omni rejects VIRTUAL gcol as PRIMARY KEY", omniErrored, true)
	})

	// --- 9.4 Gcol expression must be deterministic (NOW() rejected) ---------
	t.Run("9_4_gcol_expr_must_be_deterministic", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		bad := `CREATE TABLE t (a INT, b TIMESTAMP GENERATED ALWAYS AS (NOW()) VIRTUAL)`
		_, oracleErr := mc.db.ExecContext(mc.ctx, bad)
		if oracleErr == nil {
			t.Errorf("oracle: expected non-deterministic gcol rejection, got nil error")
		}

		omniErrored, _ := c9OmniExec(c, bad+";")
		assertBoolEq(t, "omni rejects NOW() in gcol expression", omniErrored, true)
	})

	// --- 9.5 UNIQUE on gcol allowed (both STORED and VIRTUAL under InnoDB) ---
	t.Run("9_5_unique_on_gcol_allowed", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl1 := `CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) STORED UNIQUE)`
		ddl2 := `CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL UNIQUE)`

		if _, err := mc.db.ExecContext(mc.ctx, ddl1); err != nil {
			t.Errorf("oracle t1 STORED UNIQUE: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, ddl2); err != nil {
			t.Errorf("oracle t2 VIRTUAL UNIQUE: %v", err)
		}

		if omniErr, err := c9OmniExec(c, ddl1+";"); omniErr {
			t.Errorf("omni t1 STORED UNIQUE: %v", err)
		}
		if omniErr, err := c9OmniExec(c, ddl2+";"); omniErr {
			t.Errorf("omni t2 VIRTUAL UNIQUE: %v", err)
		}

		// Oracle: both tables should have a UNIQUE index covering b.
		for _, table := range []string{"t1", "t2"} {
			var idxCount int64
			oracleScan(t, mc,
				`SELECT COUNT(*) FROM information_schema.STATISTICS
                 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+table+`'
                 AND COLUMN_NAME='b' AND NON_UNIQUE=0`,
				&idxCount)
			if idxCount < 1 {
				t.Errorf("oracle: %s should have a UNIQUE index on b", table)
			}
		}

		// Omni: both tables should have a UNIQUE index on b.
		for _, table := range []string{"t1", "t2"} {
			tbl := c.GetDatabase("testdb").GetTable(table)
			if tbl == nil {
				t.Errorf("omni: %s missing", table)
				continue
			}
			hasUnique := false
			for _, idx := range tbl.Indexes {
				if !idx.Unique {
					continue
				}
				for _, col := range idx.Columns {
					if strings.EqualFold(col.Name, "b") {
						hasUnique = true
					}
				}
			}
			if !hasUnique {
				t.Errorf("omni: %s missing UNIQUE index on b", table)
			}
		}
	})

	// --- 9.6 Gcol NOT NULL declaration accepted at CREATE time --------------
	t.Run("9_6_gcol_not_null_accepted", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
            a INT NULL,
            b INT GENERATED ALWAYS AS (a+1) VIRTUAL NOT NULL
        )`
		runOnBoth(t, mc, c, ddl)

		// Oracle: IS_NULLABLE for b is 'NO'.
		var isNullable string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`,
			&isNullable)
		assertStringEq(t, "oracle IS_NULLABLE for b", isNullable, "NO")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t missing")
			return
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Error("omni: column b missing")
			return
		}
		assertBoolEq(t, "omni b Nullable", col.Nullable, false)
		if col.Generated == nil {
			t.Error("omni: b should still be generated")
		}
	})

	// --- 9.7 FK child referencing STORED gcol parent is allowed -------------
	t.Run("9_7_fk_parent_stored_gcol_allowed", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE p (
            a INT,
            b INT GENERATED ALWAYS AS (a+1) STORED,
            UNIQUE KEY (b)
        );
CREATE TABLE c (x INT, FOREIGN KEY (x) REFERENCES p(b));`
		runOnBoth(t, mc, c, ddl)

		// Oracle: FK exists in REFERENTIAL_CONSTRAINTS.
		var fkCount int64
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
             AND REFERENCED_TABLE_NAME='p'`,
			&fkCount)
		if fkCount != 1 {
			t.Errorf("oracle: expected 1 FK on c referencing p, got %d", fkCount)
		}

		tblC := c.GetDatabase("testdb").GetTable("c")
		if tblC == nil {
			t.Error("omni: table c missing")
			return
		}
		foundFK := false
		for _, con := range tblC.Constraints {
			if con.Type == ConForeignKey && strings.EqualFold(con.RefTable, "p") {
				foundFK = true
			}
		}
		if !foundFK {
			t.Errorf("omni: no FK on c referencing p (parent STORED gcol should be allowed)")
		}
	})

	// --- 9.8 Gcol charset derived from expression inputs --------------------
	t.Run("9_8_gcol_charset_from_expression", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
            a VARCHAR(10) CHARACTER SET latin1,
            b VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a, 'x')) VIRTUAL
        )`
		runOnBoth(t, mc, c, ddl)

		// Oracle: query b's CHARACTER_SET_NAME. SCENARIOS claims this
		// should be latin1 (inherited from expression inputs). Empirical
		// oracle on MySQL 8.0.45 returns utf8mb4 (the table default),
		// because when the gcol column is declared as VARCHAR without an
		// explicit charset it picks up the table charset rather than the
		// expression's result-coercion charset. Dual-assertion: omni and
		// oracle must agree on whatever MySQL actually does.
		var cs string
		oracleScan(t, mc,
			`SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`,
			&cs)
		oracleCharset := strings.ToLower(cs)

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t missing")
			return
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Error("omni: column b missing")
			return
		}
		omniCharset := strings.ToLower(col.Charset)
		// Omni may leave the field empty when the column inherits from the
		// table default — treat that as an acceptable default state.
		if omniCharset == "" {
			omniCharset = "utf8mb4"
		}
		assertStringEq(t, "omni b Charset agrees with oracle", omniCharset, oracleCharset)
	})
}
