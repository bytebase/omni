package catalog

import "testing"

// --- Section 7.1 (Phase 7): Catalog State Isolation (8 scenarios) ---
// File target: wt_11_1_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_11_1"

func TestWalkThrough_11_1_CatalogStateIsolation(t *testing.T) {
	t.Run("separate_catalogs_independent", func(t *testing.T) {
		// Scenario 1: Execute DDL on catalog A, verify catalog B (separate New()) is unaffected.
		catA := wtSetup(t)
		catB := wtSetup(t)

		wtExec(t, catA, "CREATE TABLE t1 (id INT NOT NULL)")

		tblA := catA.GetDatabase("testdb").GetTable("t1")
		if tblA == nil {
			t.Fatal("catalog A should have table t1")
		}

		tblB := catB.GetDatabase("testdb").GetTable("t1")
		if tblB != nil {
			t.Fatal("catalog B should NOT have table t1")
		}
	})

	t.Run("dropped_reference_not_reusable", func(t *testing.T) {
		// Scenario 2: Create table, get reference, execute DROP — reference should not
		// be reusable on new table with same name.
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL)")

		oldRef := c.GetDatabase("testdb").GetTable("t1")
		if oldRef == nil {
			t.Fatal("table t1 should exist")
		}

		wtExec(t, c, "DROP TABLE t1")
		if c.GetDatabase("testdb").GetTable("t1") != nil {
			t.Fatal("t1 should be dropped")
		}

		// Re-create with different schema.
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(50))")

		newRef := c.GetDatabase("testdb").GetTable("t1")
		if newRef == nil {
			t.Fatal("new t1 should exist after re-creation")
		}

		// The old reference should differ from the new one.
		if oldRef == newRef {
			t.Fatal("old reference and new reference should be different pointers")
		}
		if len(oldRef.Columns) == len(newRef.Columns) {
			t.Fatal("old table had 1 column, new table has 2 — they should differ")
		}
	})

	t.Run("incremental_state_accumulation", func(t *testing.T) {
		// Scenario 3: Two Exec() calls building up state incrementally — state accumulates.
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL)")
		wtExec(t, c, "CREATE TABLE t2 (id INT NOT NULL, val VARCHAR(100))")

		db := c.GetDatabase("testdb")
		if db.GetTable("t1") == nil {
			t.Fatal("t1 should exist after first Exec")
		}
		if db.GetTable("t2") == nil {
			t.Fatal("t2 should exist after second Exec")
		}

		// Further accumulation via ALTER.
		wtExec(t, c, "ALTER TABLE t1 ADD COLUMN name VARCHAR(50)")
		tbl := db.GetTable("t1")
		if len(tbl.Columns) != 2 {
			t.Fatalf("expected 2 columns after ALTER, got %d", len(tbl.Columns))
		}
	})

	t.Run("continue_on_error_partial_apply", func(t *testing.T) {
		// Scenario 4: Exec with ContinueOnError — successful statements applied, failed ones not.
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL)")

		sql := `CREATE TABLE t2 (id INT NOT NULL);
CREATE TABLE t1 (id INT NOT NULL);
CREATE TABLE t3 (id INT NOT NULL);`

		results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		// t2 should succeed (index 0).
		if results[0].Error != nil {
			t.Fatalf("stmt 0 (CREATE t2) should succeed, got: %v", results[0].Error)
		}
		// t1 should fail — duplicate (index 1).
		if results[1].Error == nil {
			t.Fatal("stmt 1 (CREATE t1 duplicate) should fail")
		}
		// t3 should succeed (index 2).
		if results[2].Error != nil {
			t.Fatalf("stmt 2 (CREATE t3) should succeed, got: %v", results[2].Error)
		}

		db := c.GetDatabase("testdb")
		if db.GetTable("t2") == nil {
			t.Fatal("t2 should exist (successful stmt before error)")
		}
		if db.GetTable("t3") == nil {
			t.Fatal("t3 should exist (successful stmt after error with ContinueOnError)")
		}
	})

	t.Run("exec_empty_string", func(t *testing.T) {
		// Scenario 5: Exec empty string — no state change, no error.
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL)")

		results, err := c.Exec("", nil)
		if err != nil {
			t.Fatalf("expected no parse error, got: %v", err)
		}
		if results != nil {
			t.Fatalf("expected nil results for empty string, got %d results", len(results))
		}

		// State unchanged.
		if c.GetDatabase("testdb").GetTable("t1") == nil {
			t.Fatal("t1 should still exist after empty Exec")
		}
	})

	t.Run("exec_dml_only_skipped", func(t *testing.T) {
		// Scenario 6: Exec with only DML — no state change, statements skipped.
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL)")

		results, err := c.Exec("SELECT 1; INSERT INTO t1 VALUES (1);", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		for i, r := range results {
			if r.Error != nil {
				t.Fatalf("stmt %d should not error, got: %v", i, r.Error)
			}
			if !r.Skipped {
				t.Fatalf("stmt %d should be skipped (DML)", i)
			}
		}

		// State unchanged — table still has same structure.
		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Fatal("t1 should still exist")
		}
		if len(tbl.Columns) != 1 {
			t.Fatalf("expected 1 column, got %d", len(tbl.Columns))
		}
	})

	t.Run("drop_database_cascade", func(t *testing.T) {
		// Scenario 7: DROP DATABASE cascade — all tables, views, routines, triggers, events removed.
		c := wtSetup(t)

		// Create various objects.
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, val INT)")
		wtExec(t, c, "CREATE VIEW v1 AS SELECT id FROM t1")
		wtExec(t, c, "CREATE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1")
		wtExec(t, c, "CREATE TRIGGER tr1 BEFORE INSERT ON t1 FOR EACH ROW SET NEW.val = 0")
		wtExec(t, c, "CREATE EVENT ev1 ON SCHEDULE EVERY 1 HOUR DO SELECT 1")

		// Verify objects exist.
		db := c.GetDatabase("testdb")
		if db == nil {
			t.Fatal("testdb should exist")
		}
		if db.GetTable("t1") == nil {
			t.Fatal("t1 should exist before DROP DATABASE")
		}
		if db.Views[toLower("v1")] == nil {
			t.Fatal("v1 should exist before DROP DATABASE")
		}
		if db.Functions[toLower("f1")] == nil {
			t.Fatal("f1 should exist before DROP DATABASE")
		}
		if db.Triggers[toLower("tr1")] == nil {
			t.Fatal("tr1 should exist before DROP DATABASE")
		}
		if db.Events[toLower("ev1")] == nil {
			t.Fatal("ev1 should exist before DROP DATABASE")
		}

		// DROP DATABASE removes everything.
		wtExec(t, c, "DROP DATABASE testdb")

		if c.GetDatabase("testdb") != nil {
			t.Fatal("testdb should not exist after DROP DATABASE")
		}
		// Current database should be cleared.
		if c.CurrentDatabase() != "" {
			t.Fatalf("current database should be empty after DROP, got %q", c.CurrentDatabase())
		}
	})

	t.Run("operations_across_two_databases", func(t *testing.T) {
		// Scenario 8: Operations across two databases — CREATE TABLE in db1 and db2 independently.
		c := New()
		wtExec(t, c, "CREATE DATABASE db1; CREATE DATABASE db2;")

		c.SetCurrentDatabase("db1")
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, a VARCHAR(50))")

		c.SetCurrentDatabase("db2")
		wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, b INT, c INT)")

		// Verify db1.t1 has 2 columns.
		db1 := c.GetDatabase("db1")
		tbl1 := db1.GetTable("t1")
		if tbl1 == nil {
			t.Fatal("db1.t1 should exist")
		}
		if len(tbl1.Columns) != 2 {
			t.Fatalf("db1.t1 expected 2 columns, got %d", len(tbl1.Columns))
		}

		// Verify db2.t1 has 3 columns.
		db2 := c.GetDatabase("db2")
		tbl2 := db2.GetTable("t1")
		if tbl2 == nil {
			t.Fatal("db2.t1 should exist")
		}
		if len(tbl2.Columns) != 3 {
			t.Fatalf("db2.t1 expected 3 columns, got %d", len(tbl2.Columns))
		}

		// Dropping table in db1 should not affect db2.
		c.SetCurrentDatabase("db1")
		wtExec(t, c, "DROP TABLE t1")
		if db1.GetTable("t1") != nil {
			t.Fatal("db1.t1 should be dropped")
		}
		if db2.GetTable("t1") == nil {
			t.Fatal("db2.t1 should still exist after dropping db1.t1")
		}
	})
}
