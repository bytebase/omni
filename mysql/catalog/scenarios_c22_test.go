package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C22 covers Section C22 "ALTER TABLE algorithm / lock defaults"
// from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md.
//
// omni's catalog intentionally does NOT track ALGORITHM= / LOCK= clauses —
// they are execution-time concerns, not persisted schema state. So these
// tests verify three things per scenario:
//
//  1. omni parses and accepts the ALTER ... ALGORITHM=... / LOCK=... form.
//  2. omni applies the underlying column/index change to catalog state
//     identically whether or not an ALGORITHM/LOCK clause is present.
//  3. MySQL 8.0 actually accepts or rejects the statement as the scenario
//     claims (oracle sanity check — a regression in MySQL wouldn't be an
//     omni bug but would invalidate the assertion).
//
// Algorithm selection (scenarios 22.1/22.2) is not directly observable from
// the client side without inspecting internal counters, so those tests only
// verify that DEFAULT clauses parse and succeed on both sides.
//
// Failures in omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c22.md.
func TestScenario_C22(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// -----------------------------------------------------------------
	// 22.1 ALGORITHM=DEFAULT picks fastest supported (oracle-only parse check)
	// -----------------------------------------------------------------
	t.Run("22_1_algorithm_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB")
		// Bare ADD COLUMN — DEFAULT algorithm, DEFAULT lock.
		runOnBoth(t, mc, c, "ALTER TABLE t1 ADD COLUMN b INT")
		// Explicit ALGORITHM=DEFAULT, LOCK=DEFAULT.
		runOnBoth(t, mc, c,
			"ALTER TABLE t1 ADD COLUMN c INT, ALGORITHM=DEFAULT, LOCK=DEFAULT")

		// Oracle: both columns present.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1'
			 ORDER BY ORDINAL_POSITION`)
		var oracleCols []string
		for _, r := range rows {
			oracleCols = append(oracleCols, asString(r[0]))
		}
		assertStringEq(t, "oracle columns", strings.Join(oracleCols, ","),
			"id,a,b,c")

		// omni: both columns applied to catalog.
		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		var omniCols []string
		for _, col := range tbl.Columns {
			omniCols = append(omniCols, col.Name)
		}
		assertStringEq(t, "omni columns", strings.Join(omniCols, ","),
			"id,a,b,c")
	})

	// -----------------------------------------------------------------
	// 22.2 LOCK=DEFAULT picks least restrictive (oracle-only parse check)
	// -----------------------------------------------------------------
	t.Run("22_2_lock_default_add_index", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB")
		runOnBoth(t, mc, c, "ALTER TABLE t1 ADD INDEX ix_a (a)")
		runOnBoth(t, mc, c,
			"ALTER TABLE t1 ADD INDEX ix_a2 (a), LOCK=DEFAULT")

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		foundIxA, foundIxA2 := false, false
		for _, idx := range tbl.Indexes {
			switch idx.Name {
			case "ix_a":
				foundIxA = true
			case "ix_a2":
				foundIxA2 = true
			}
		}
		assertBoolEq(t, "omni ix_a present", foundIxA, true)
		assertBoolEq(t, "omni ix_a2 present", foundIxA2, true)
	})

	// -----------------------------------------------------------------
	// 22.3 ADD COLUMN trailing nullable is INSTANT in 8.0.12+
	// -----------------------------------------------------------------
	t.Run("22_3_add_column_instant", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB ROW_FORMAT=DYNAMIC")
		// 22.3a: bare ADD COLUMN.
		runOnBoth(t, mc, c, "ALTER TABLE t1 ADD COLUMN b VARCHAR(32) NULL")
		// 22.3b: explicit ALGORITHM=INSTANT.
		runOnBoth(t, mc, c,
			"ALTER TABLE t1 ADD COLUMN c INT NULL, ALGORITHM=INSTANT")

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		var names []string
		for _, col := range tbl.Columns {
			names = append(names, col.Name)
		}
		assertStringEq(t, "omni columns after INSTANT adds",
			strings.Join(names, ","), "id,a,b,c")

		// Oracle sanity: ALGORITHM=INSTANT statement must have succeeded.
		var n int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='c'`,
			&n)
		assertIntEq(t, "oracle column c exists", n, 1)
	})

	// -----------------------------------------------------------------
	// 22.4 DROP COLUMN INSTANT (8.0.29+) else INPLACE
	// -----------------------------------------------------------------
	t.Run("22_4_drop_column", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT, b INT) ENGINE=InnoDB")
		// 22.4a: bare DROP COLUMN.
		runOnBoth(t, mc, c, "ALTER TABLE t1 DROP COLUMN b")

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		if tbl.GetColumn("b") != nil {
			t.Errorf("omni: column b should be dropped")
		}
		if tbl.GetColumn("a") == nil {
			t.Errorf("omni: column a should still exist")
		}

		// 22.4b: explicit ALGORITHM=INSTANT on a fresh table. On MySQL 8.0.29+
		// this succeeds; on earlier versions it errors. Either way, omni's
		// parser must accept it, and if oracle accepts it, omni must apply
		// the drop.
		runOnBoth(t, mc, c,
			"CREATE TABLE t2 (id INT PRIMARY KEY, a INT, b INT) ENGINE=InnoDB")
		_, oracleErr := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t2 DROP COLUMN b, ALGORITHM=INSTANT")

		// omni: apply against catalog regardless.
		results, err := c.Exec("ALTER TABLE t2 DROP COLUMN b, ALGORITHM=INSTANT;", nil)
		if err != nil {
			t.Errorf("omni: parse error for DROP COLUMN + ALGORITHM=INSTANT: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}

		if oracleErr == nil {
			// 8.0.29+ — drop succeeded, omni must have the column gone too.
			if tbl2 := c.GetDatabase("testdb").GetTable("t2"); tbl2 != nil &&
				tbl2.GetColumn("b") != nil {
				t.Errorf("omni: t2.b should be dropped after ALGORITHM=INSTANT")
			}
		} else {
			// Pre-8.0.29 — oracle rejected. omni catalog is algorithm-oblivious
			// and still drops the column, which is the documented correct
			// behavior (SDL diff classifier must filter these, not the
			// catalog). Only assert the oracle error kind.
			if !strings.Contains(oracleErr.Error(), "ALGORITHM") &&
				!strings.Contains(oracleErr.Error(), "not supported") {
				t.Errorf("oracle: unexpected error form: %v", oracleErr)
			}
		}
	})

	// -----------------------------------------------------------------
	// 22.5 RENAME COLUMN metadata-only (INPLACE / INSTANT on 8.0.29+)
	// -----------------------------------------------------------------
	t.Run("22_5_rename_column", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, old_name INT) ENGINE=InnoDB")
		// 22.5a: RENAME COLUMN form.
		runOnBoth(t, mc, c, "ALTER TABLE t1 RENAME COLUMN old_name TO new_name")

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		if tbl.GetColumn("old_name") != nil {
			t.Errorf("omni: old_name should be gone after rename")
		}
		if tbl.GetColumn("new_name") == nil {
			t.Errorf("omni: new_name should exist after rename")
		}

		// Oracle: INFORMATION_SCHEMA.COLUMNS should have new_name.
		var n int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='new_name'`,
			&n)
		assertIntEq(t, "oracle new_name present", n, 1)

		// 22.5b: CHANGE COLUMN form (type unchanged).
		runOnBoth(t, mc, c,
			"ALTER TABLE t1 CHANGE COLUMN new_name newer_name INT")
		if tbl.GetColumn("new_name") != nil {
			t.Errorf("omni: new_name should be gone after CHANGE rename")
		}
		if tbl.GetColumn("newer_name") == nil {
			t.Errorf("omni: newer_name should exist after CHANGE rename")
		}
	})

	// -----------------------------------------------------------------
	// 22.6 CHANGE COLUMN type forces COPY; explicit INPLACE/LOCK=NONE error
	// -----------------------------------------------------------------
	t.Run("22_6_change_column_type_copy", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB")

		// 22.6a: bare CHANGE COLUMN type. DEFAULT → COPY, succeeds.
		runOnBoth(t, mc, c, "ALTER TABLE t1 CHANGE COLUMN a a BIGINT")

		// Oracle: column a now has type bigint.
		var oracleType string
		oracleScan(t, mc,
			`SELECT DATA_TYPE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='a'`,
			&oracleType)
		if strings.ToLower(oracleType) != "bigint" {
			t.Errorf("oracle: column a type got %q want bigint", oracleType)
		}

		// omni: column a's type should be bigint.
		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni: table t1 missing")
			return
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Errorf("omni: column a missing after CHANGE")
		} else if !strings.Contains(strings.ToLower(col.DataType), "bigint") {
			t.Errorf("omni: column a type got %q want bigint", col.DataType)
		}

		// 22.6b: MODIFY + ALGORITHM=INPLACE → MySQL rejects.
		_, oracleErr := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t1 MODIFY COLUMN a INT, ALGORITHM=INPLACE")
		if oracleErr == nil {
			t.Errorf("oracle: expected error for MODIFY ... ALGORITHM=INPLACE, got nil")
		} else if !strings.Contains(oracleErr.Error(), "ALGORITHM") &&
			!strings.Contains(oracleErr.Error(), "not supported") {
			t.Errorf("oracle: unexpected error form: %v", oracleErr)
		}
		// omni must at least parse and (per catalog contract) apply the change.
		results, err := c.Exec("ALTER TABLE t1 MODIFY COLUMN a INT, ALGORITHM=INPLACE;", nil)
		if err != nil {
			t.Errorf("omni: parse error on MODIFY + ALGORITHM=INPLACE: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}

		// 22.6c: MODIFY with a genuine type change + LOCK=NONE → MySQL
		// rejects. Column `a` is currently bigint (from 22.6a); changing it
		// to INT UNSIGNED is a real type change, which forces COPY, which is
		// incompatible with LOCK=NONE. (A no-op MODIFY to the same type is
		// silently optimized away and would not error.)
		_, oracleErr2 := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t1 MODIFY COLUMN a INT UNSIGNED, LOCK=NONE")
		if oracleErr2 == nil {
			t.Errorf("oracle: expected error for MODIFY <typechange> ... LOCK=NONE, got nil")
		}
		results2, err2 := c.Exec("ALTER TABLE t1 MODIFY COLUMN a INT UNSIGNED, LOCK=NONE;", nil)
		if err2 != nil {
			t.Errorf("omni: parse error on MODIFY + LOCK=NONE: %v", err2)
		}
		for _, r := range results2 {
			if r.Error != nil {
				t.Errorf("omni: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}
	})

	// -----------------------------------------------------------------
	// 22.7 ALGORITHM=INSTANT on unsupported operation → hard error
	// -----------------------------------------------------------------
	t.Run("22_7_instant_unsupported_hard_error", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB")

		// PK rebuild + ALGORITHM=INSTANT → MySQL rejects with
		// ER_ALTER_OPERATION_NOT_SUPPORTED_REASON.
		_, oracleErr := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t1 DROP PRIMARY KEY, ADD PRIMARY KEY (a), ALGORITHM=INSTANT")
		if oracleErr == nil {
			t.Errorf("oracle: expected error for PK rebuild ALGORITHM=INSTANT, got nil")
		} else if !strings.Contains(oracleErr.Error(), "ALGORITHM") &&
			!strings.Contains(oracleErr.Error(), "not supported") {
			t.Errorf("oracle: unexpected error form: %v", oracleErr)
		}

		// omni: parser must accept; catalog applies PK swap regardless.
		results, err := c.Exec(
			"ALTER TABLE t1 DROP PRIMARY KEY, ADD PRIMARY KEY (a), ALGORITHM=INSTANT;", nil)
		if err != nil {
			t.Errorf("omni: parse error on PK rebuild + ALGORITHM=INSTANT: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}
	})

	// -----------------------------------------------------------------
	// 22.8 LOCK=NONE on COPY-only op errors; LOCK=... with INSTANT errors
	// -----------------------------------------------------------------
	t.Run("22_8_lock_none_copy_and_instant", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			"CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB")

		// 22.8a: COPY-only operation + LOCK=NONE → hard error. CHANGE COLUMN
		// with a real type change forces COPY, and COPY cannot honor
		// LOCK=NONE.
		_, oracleErr := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t1 CHANGE COLUMN a a BIGINT UNSIGNED, LOCK=NONE")
		if oracleErr == nil {
			t.Errorf("oracle 22.8a: expected error for COPY op + LOCK=NONE, got nil")
		}
		results, _ := c.Exec(
			"ALTER TABLE t1 CHANGE COLUMN a a BIGINT UNSIGNED, LOCK=NONE;", nil)
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni 22.8a: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}

		// 22.8b: ADD COLUMN + ALGORITHM=INSTANT + LOCK=NONE → hard error
		// ("Only LOCK=DEFAULT is permitted for operations using ALGORITHM=INSTANT").
		_, oracleErr2 := mc.db.ExecContext(mc.ctx,
			"ALTER TABLE t1 ADD COLUMN b INT, ALGORITHM=INSTANT, LOCK=NONE")
		if oracleErr2 == nil {
			t.Errorf("oracle 22.8b: expected error for INSTANT + LOCK=NONE, got nil")
		}
		results2, err2 := c.Exec(
			"ALTER TABLE t1 ADD COLUMN b INT, ALGORITHM=INSTANT, LOCK=NONE;", nil)
		if err2 != nil {
			t.Errorf("omni 22.8b: parse error: %v", err2)
		}
		for _, r := range results2 {
			if r.Error != nil {
				t.Errorf("omni 22.8b: exec error on stmt %d: %v", r.Index, r.Error)
			}
		}

		// 22.8c: ADD COLUMN + ALGORITHM=INSTANT (no LOCK) → succeeds.
		// Note: b may already exist in omni catalog from 22.8b above; use c.
		runOnBoth(t, mc, c,
			"ALTER TABLE t1 ADD COLUMN c INT, ALGORITHM=INSTANT")

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Errorf("omni 22.8c: table t1 missing")
			return
		}
		if tbl.GetColumn("c") == nil {
			t.Errorf("omni 22.8c: column c should be added")
		}
	})
}
