package catalog

import "testing"

// TestWalkThrough_5_4 covers section 1.4 — Error Semantics in Multi-Command ALTER.
// Multi-command ALTER TABLE is atomic in MySQL: if any sub-command fails, the
// entire ALTER is rolled back and the table state is unchanged.
func TestWalkThrough_5_4(t *testing.T) {
	t.Run("dup_column_in_same_alter", func(t *testing.T) {
		// ADD COLUMN x INT, ADD COLUMN x INT — duplicate column error, verify rollback
		c := setupTestTable(t) // t1: id INT NOT NULL, name VARCHAR(100), age INT

		results, _ := c.Exec(
			"ALTER TABLE t1 ADD COLUMN x INT, ADD COLUMN x INT",
			&ExecOptions{ContinueOnError: true},
		)
		assertError(t, results[0].Error, ErrDupColumn)

		// Verify rollback: table should be unchanged.
		tbl := c.GetDatabase("test").GetTable("t1")
		if len(tbl.Columns) != 3 {
			t.Errorf("expected 3 columns after rollback, got %d", len(tbl.Columns))
		}
		if tbl.GetColumn("x") != nil {
			t.Error("column 'x' should not exist after rollback")
		}
	})

	t.Run("add_then_drop_nonexistent", func(t *testing.T) {
		// ADD COLUMN x INT, DROP COLUMN nonexistent — first succeeds, second errors;
		// verify x was NOT added (MySQL rolls back entire ALTER)
		c := setupTestTable(t)

		results, _ := c.Exec(
			"ALTER TABLE t1 ADD COLUMN x INT, DROP COLUMN nonexistent",
			&ExecOptions{ContinueOnError: true},
		)
		if results[0].Error == nil {
			t.Fatal("expected error from DROP COLUMN nonexistent")
		}

		tbl := c.GetDatabase("test").GetTable("t1")
		if len(tbl.Columns) != 3 {
			t.Errorf("expected 3 columns after rollback, got %d", len(tbl.Columns))
		}
		if tbl.GetColumn("x") != nil {
			t.Error("column 'x' should not exist after rollback")
		}
	})

	t.Run("modify_nonexistent_then_add", func(t *testing.T) {
		// MODIFY COLUMN nonexistent INT, ADD COLUMN y INT — first errors, second never runs
		c := setupTestTable(t)

		results, _ := c.Exec(
			"ALTER TABLE t1 MODIFY COLUMN nonexistent INT, ADD COLUMN y INT",
			&ExecOptions{ContinueOnError: true},
		)
		assertError(t, results[0].Error, ErrNoSuchColumn)

		tbl := c.GetDatabase("test").GetTable("t1")
		if len(tbl.Columns) != 3 {
			t.Errorf("expected 3 columns after rollback, got %d", len(tbl.Columns))
		}
		if tbl.GetColumn("y") != nil {
			t.Error("column 'y' should not exist — entire ALTER rolled back")
		}
	})

	t.Run("drop_nonexistent_index_then_add_column", func(t *testing.T) {
		// DROP INDEX nonexistent, ADD COLUMN y INT — error on first, verify y not added
		c := setupTestTable(t)

		results, _ := c.Exec(
			"ALTER TABLE t1 DROP INDEX nonexistent, ADD COLUMN y INT",
			&ExecOptions{ContinueOnError: true},
		)
		assertError(t, results[0].Error, ErrCantDropKey)

		tbl := c.GetDatabase("test").GetTable("t1")
		if len(tbl.Columns) != 3 {
			t.Errorf("expected 3 columns after rollback, got %d", len(tbl.Columns))
		}
		if tbl.GetColumn("y") != nil {
			t.Error("column 'y' should not exist — entire ALTER rolled back")
		}
	})

	t.Run("add_index_then_drop_same_column", func(t *testing.T) {
		// ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x — MySQL resolves
		// the DROP COLUMN against the table the statement found, which has no
		// x, so it refuses with 1091 and the ALTER rolls back. (This case used
		// to record that MySQL takes the statement, which the pinned server
		// does not: mysql:8.0.46 answers 1091 "Can't DROP 'x'".)
		c := setupTestTable(t)

		results, _ := c.Exec(
			"ALTER TABLE t1 ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x",
			&ExecOptions{ContinueOnError: true},
		)
		assertError(t, results[0].Error, ErrCantDropKey)

		tbl := c.GetDatabase("test").GetTable("t1")
		if len(tbl.Columns) != 3 {
			t.Errorf("expected 3 columns after rollback, got %d", len(tbl.Columns))
		}
		if tbl.GetColumn("x") != nil {
			t.Error("column 'x' should not exist after rollback")
		}
	})
}
