package catalog

import "testing"

// --- Section 9.1 (Phase 9): SET Variable Effects (7 scenarios) ---
// File target: wt_13_1_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_13_1"

func TestWalkThrough_13_1_SetVariableEffects(t *testing.T) {
	// Scenario 1: SET foreign_key_checks = 0 then CREATE TABLE with invalid FK — succeeds
	t.Run("fk_checks_0_allows_invalid_fk", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET foreign_key_checks = 0")
		// Create a table with FK referencing a non-existent table — should succeed.
		wtExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES nonexistent(id)
		)`)
		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child should exist")
		}
	})

	// Scenario 2: SET foreign_key_checks = 1 then CREATE TABLE with invalid FK — fails
	t.Run("fk_checks_1_rejects_invalid_fk", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET foreign_key_checks = 1")
		results, _ := c.Exec(`CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES nonexistent(id)
		)`, nil)
		if len(results) == 0 {
			t.Fatal("expected result from CREATE TABLE")
		}
		assertError(t, results[0].Error, ErrFKNoRefTable)
	})

	// Scenario 3: SET foreign_key_checks = OFF — accepts OFF as 0
	t.Run("fk_checks_off_as_0", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET foreign_key_checks = OFF")
		if c.ForeignKeyChecks() {
			t.Error("foreign_key_checks should be false after SET OFF")
		}
		// Verify it actually allows invalid FK.
		wtExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES nonexistent(id)
		)`)
		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child should exist")
		}
	})

	// Scenario 4: SET NAMES utf8mb4 — silently accepted
	t.Run("set_names_silently_accepted", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET NAMES utf8mb4")
		// No state change — catalog still works normally.
		wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, PRIMARY KEY (id))")
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("table t should exist after SET NAMES")
		}
	})

	// Scenario 5: SET CHARACTER SET latin1 — silently accepted
	t.Run("set_character_set_silently_accepted", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET CHARACTER SET latin1")
		// No state change — catalog still works normally.
		wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, PRIMARY KEY (id))")
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("table t should exist after SET CHARACTER SET")
		}
	})

	// Scenario 6: SET sql_mode = 'STRICT_TRANS_TABLES' — silently accepted
	t.Run("set_sql_mode_silently_accepted", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET sql_mode = 'STRICT_TRANS_TABLES'")
		// No state change — catalog still works normally.
		wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, PRIMARY KEY (id))")
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("table t should exist after SET sql_mode")
		}
	})

	// Scenario 7: SET unknown_variable = 'value' — silently accepted
	t.Run("set_unknown_variable_silently_accepted", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET unknown_variable = 'value'")
		// No state change — catalog still works normally.
		wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, PRIMARY KEY (id))")
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("table t should exist after SET unknown_variable")
		}
	})
}
