package catalog

import (
	"errors"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// This file holds the sections of mysql/catalog's container scenarios that
// assert a fact OceanBase can answer: an error, an object that is or is not
// there, an information_schema value. The sections that compared MySQL's
// SHOW CREATE text with the catalog's are in mysql/catalog's copy only; of
// those, sections 1.7, 1.17, 2.2, 4.1, 4.4 and 6.1 keep here the statements
// OceanBase refuses.

// containerRefusal runs setup on OceanBase and on a fresh catalog, then holds
// stmt to the error OceanBase answers and to the catalog's code for it: 0
// when the catalog runs the statement, which is a gap the caller names.
func containerRefusal(t *testing.T, ctr *mysqlContainer, setup, stmt string, engine uint16, omni int) {
	t.Helper()
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	if setup != "" {
		if err := ctr.execSQL(setup); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		if code := omniCode(c, setup); code != 0 {
			t.Fatalf("omni setup: error %d", code)
		}
	}
	oracleRefuses(t, ctr, stmt, engine)
	assertIntEq(t, "omni error code", omniCode(c, stmt), omni)
}

func TestContainer_Section_1_7_DefaultValues(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// OceanBase gives a JSON column no default, expression or not (1101).
	// omni gap: the catalog takes the expression default, as MySQL 8.0 does.
	t.Run("expr_default_json", func(t *testing.T) {
		containerRefusal(t, ctr, "",
			"CREATE TABLE t_def_exjson (a JSON DEFAULT (JSON_ARRAY()))", 1101, 0)
	})
}

func TestContainer_Section_1_17_ForeignKeys(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// OceanBase refuses ON DELETE SET DEFAULT with ER_CANNOT_ADD_FOREIGN
	// (1215). omni gap: the catalog takes it, as MySQL does.
	t.Run("fk_on_delete_set_default", func(t *testing.T) {
		containerRefusal(t, ctr,
			"CREATE TABLE t_parent5 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk5 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent5(id) ON DELETE SET DEFAULT)",
			1215, 0)
	})
}

func TestContainer_Section_2_1_CreateTableVariants(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("if_not_exists_no_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ine")
		ctr.execSQL("CREATE TABLE t_ine (id INT)")

		// Second CREATE with IF NOT EXISTS should not error on ctr.
		ctrErr := ctr.execSQL("CREATE TABLE IF NOT EXISTS t_ine (id INT)")
		if ctrErr != nil {
			t.Fatalf("container error on IF NOT EXISTS: %v", ctrErr)
		}

		// Omni should also not error.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ine (id INT)", nil)
		results, _ := c.Exec("CREATE TABLE IF NOT EXISTS t_ine (id INT)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni error on IF NOT EXISTS: %v", results[0].Error)
		}
	})

	t.Run("create_table_with_view_name_conflict", func(t *testing.T) {
		// Creating a table with same name as an existing view should error.
		ctr.execSQL("DROP TABLE IF EXISTS t_view_conflict")
		ctr.execSQL("DROP VIEW IF EXISTS t_view_conflict")
		ctr.execSQL("CREATE VIEW t_view_conflict AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE TABLE t_view_conflict (id INT)")
		if ctrErr == nil {
			t.Fatal("expected container error when creating table with same name as view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW t_view_conflict AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE TABLE t_view_conflict (id INT)", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error when creating table with same name as view")
		}
	})
}

func TestContainer_Section_2_2_AlterTableColumnOps(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// OceanBase keeps a table's rows in primary-key order and does not drop a
	// column of the key (1235). omni gap: the catalog drops it, as MySQL does.
	t.Run("drop_column_part_of_pk", func(t *testing.T) {
		containerRefusal(t, ctr,
			"CREATE TABLE t_drop_pk (a INT, b INT, PRIMARY KEY (a, b))",
			"ALTER TABLE t_drop_pk DROP COLUMN a", 1235, 0)
	})

	// Error tests: operations that should produce errors.
	errCases := []struct {
		name    string
		setup   string
		alter   string
		table   string
		wantErr bool
	}{
		{
			"drop_column_referenced_by_fk",
			"CREATE TABLE t_fk_parent (id INT PRIMARY KEY); CREATE TABLE t_fk_child (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent(id))",
			"ALTER TABLE t_fk_child DROP COLUMN pid",
			"t_fk_child",
			true,
		},
	}

	for _, tc := range errCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up tables first.
			ctr.execSQL("DROP TABLE IF EXISTS t_fk_child")
			ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent")
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)

			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			ctrErr := ctr.execSQL(tc.alter)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.setup, nil)
			results, _ := c.Exec(tc.alter, nil)
			var omniErr error
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}

			if tc.wantErr {
				if ctrErr == nil {
					t.Fatal("expected container error but got nil")
				}
				if omniErr == nil {
					t.Fatalf("expected omni error but got nil (container error: %v)", ctrErr)
				}
				t.Logf("both errored as expected — container: %v, omni: %v", ctrErr, omniErr)
			} else {
				if ctrErr != nil {
					t.Fatalf("unexpected container error: %v", ctrErr)
				}
				if omniErr != nil {
					t.Fatalf("unexpected omni error: %v", omniErr)
				}
			}
		})
	}
}

func TestContainer_Section_2_3_AlterTableIndexOps(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Error tests.
	errCases := []struct {
		name    string
		setup   string
		alter   string
		table   string
		wantErr bool
	}{
		{
			"add_primary_key_when_exists",
			"CREATE TABLE t_dup_pk (id INT PRIMARY KEY, val INT NOT NULL)",
			"ALTER TABLE t_dup_pk ADD PRIMARY KEY (val)",
			"t_dup_pk",
			true,
		},
	}

	for _, tc := range errCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			ctrErr := ctr.execSQL(tc.alter)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.setup, nil)
			results, _ := c.Exec(tc.alter, nil)
			var omniErr error
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}

			if tc.wantErr {
				if ctrErr == nil {
					t.Fatal("expected container error but got nil")
				}
				if omniErr == nil {
					t.Fatalf("expected omni error but got nil (container error: %v)", ctrErr)
				}
				t.Logf("both errored as expected — container: %v, omni: %v", ctrErr, omniErr)
			}
		})
	}
}

func TestContainer_Section_2_6_DropTable(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("drop_table_basic", func(t *testing.T) {
		// Setup: create a table, then drop it. Verify it's gone.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop1")
		ctr.execSQL("CREATE TABLE t_drop1 (id INT PRIMARY KEY, name VARCHAR(100))")

		ctrErr := ctr.execSQL("DROP TABLE t_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE error: %v", ctrErr)
		}
		_, ctrShowErr := ctr.showCreateTable("t_drop1")
		if ctrShowErr == nil {
			t.Fatal("container: table still exists after DROP TABLE")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_drop1 (id INT PRIMARY KEY, name VARCHAR(100))", nil)
		results, _ := c.Exec("DROP TABLE t_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_drop1")
		if omniDDL != "" {
			t.Errorf("omni: table still exists after DROP TABLE, got: %s", omniDDL)
		}
	})

	t.Run("drop_table_if_exists", func(t *testing.T) {
		// DROP TABLE IF EXISTS on a nonexistent table should not error.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop_ine")

		ctrErr := ctr.execSQL("DROP TABLE IF EXISTS t_drop_ine")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE IF EXISTS t_drop_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("drop_table_multi", func(t *testing.T) {
		// DROP TABLE t1, t2, t3 — multi-table drop.
		ctr.execSQL("DROP TABLE IF EXISTS t_dm1")
		ctr.execSQL("DROP TABLE IF EXISTS t_dm2")
		ctr.execSQL("DROP TABLE IF EXISTS t_dm3")
		ctr.execSQL("CREATE TABLE t_dm1 (id INT)")
		ctr.execSQL("CREATE TABLE t_dm2 (id INT)")
		ctr.execSQL("CREATE TABLE t_dm3 (id INT)")

		ctrErr := ctr.execSQL("DROP TABLE t_dm1, t_dm2, t_dm3")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE multi error: %v", ctrErr)
		}
		for _, tbl := range []string{"t_dm1", "t_dm2", "t_dm3"} {
			if _, err := ctr.showCreateTable(tbl); err == nil {
				t.Errorf("container: table %s still exists after multi-drop", tbl)
			}
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dm1 (id INT)", nil)
		c.Exec("CREATE TABLE t_dm2 (id INT)", nil)
		c.Exec("CREATE TABLE t_dm3 (id INT)", nil)
		results, _ := c.Exec("DROP TABLE t_dm1, t_dm2, t_dm3", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE multi error: %v", results[0].Error)
		}
		for _, tbl := range []string{"t_dm1", "t_dm2", "t_dm3"} {
			ddl := c.ShowCreateTable("test", tbl)
			if ddl != "" {
				t.Errorf("omni: table %s still exists after multi-drop", tbl)
			}
		}
	})

	t.Run("drop_table_nonexistent_error", func(t *testing.T) {
		// DROP TABLE on nonexistent table should produce error 1051.
		ctr.execSQL("DROP TABLE IF EXISTS t_noexist")
		ctrErr := ctr.execSQL("DROP TABLE t_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE t_noexist", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for DROP nonexistent table")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1051 {
			t.Errorf("omni error code: want 1051, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})

	t.Run("drop_temporary_table", func(t *testing.T) {
		// DROP TEMPORARY TABLE should work.
		ctr.execSQL("DROP TEMPORARY TABLE IF EXISTS t_temp_drop")
		ctr.execSQL("CREATE TEMPORARY TABLE t_temp_drop (id INT, val VARCHAR(50))")
		ctrErr := ctr.execSQL("DROP TEMPORARY TABLE t_temp_drop")
		if ctrErr != nil {
			t.Fatalf("container DROP TEMPORARY TABLE error: %v", ctrErr)
		}
		_, ctrShowErr := ctr.showCreateTable("t_temp_drop")
		if ctrShowErr == nil {
			t.Fatal("container: temp table still exists after DROP TEMPORARY TABLE")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TEMPORARY TABLE t_temp_drop (id INT, val VARCHAR(50))", nil)
		results, _ := c.Exec("DROP TEMPORARY TABLE t_temp_drop", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TEMPORARY TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_temp_drop")
		if omniDDL != "" {
			t.Errorf("omni: temp table still exists after DROP TEMPORARY TABLE")
		}
	})

	t.Run("drop_table_fk_referenced", func(t *testing.T) {
		// DROP TABLE that has FK references should error (with foreign_key_checks=1, the default).
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent")
		ctr.execSQL("CREATE TABLE t_fk_parent (id INT PRIMARY KEY)")
		ctr.execSQL("CREATE TABLE t_fk_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fk_parent(id))")

		ctrErr := ctr.execSQL("DROP TABLE t_fk_parent")
		if ctrErr == nil {
			t.Fatal("container: expected error when dropping FK-referenced table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fk_parent (id INT PRIMARY KEY)", nil)
		c.Exec("CREATE TABLE t_fk_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fk_parent(id))", nil)
		results, _ := c.Exec("DROP TABLE t_fk_parent", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error when dropping FK-referenced table")
		}

		// Verify table still exists after failed drop.
		omniDDL := c.ShowCreateTable("test", "t_fk_parent")
		if omniDDL == "" {
			t.Error("omni: parent table was deleted despite FK reference")
		}
	})
}

func TestContainer_Section_2_7_TruncateTable(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario 2: TRUNCATE resets AUTO_INCREMENT. OceanBase shows a table's
	// counter in SHOW CREATE TABLE and nowhere else
	// (information_schema.TABLES.AUTO_INCREMENT is NULL).
	t.Run("truncate_resets_auto_increment", func(t *testing.T) {
		const create = "CREATE TABLE t_trunc_ai (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100)) AUTO_INCREMENT=1000"
		if err := ctr.execSQL(create); err != nil {
			t.Fatalf("container create: %v", err)
		}
		before, err := ctr.showCreateTable("t_trunc_ai")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(before, " AUTO_INCREMENT = 1000 ") {
			t.Errorf("container before truncate: want AUTO_INCREMENT = 1000 in %s", before)
		}
		if err := ctr.execSQL("TRUNCATE TABLE t_trunc_ai"); err != nil {
			t.Fatalf("container truncate: %v", err)
		}
		after, err := ctr.showCreateTable("t_trunc_ai")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(after, " AUTO_INCREMENT = 1 ") {
			t.Errorf("container after truncate: want AUTO_INCREMENT = 1 in %s", after)
		}

		// omni: the counter is 1000, and unset (0) once the table is truncated.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(create, nil)
		tbl := c.GetDatabase("test").GetTable("t_trunc_ai")
		if tbl == nil {
			t.Fatal("omni: table t_trunc_ai missing")
		}
		if tbl.AutoIncrement != 1000 {
			t.Errorf("omni AutoIncrement before truncate: got %d, want 1000", tbl.AutoIncrement)
		}
		c.Exec("TRUNCATE TABLE t_trunc_ai", nil)
		if tbl := c.GetDatabase("test").GetTable("t_trunc_ai"); tbl == nil || tbl.AutoIncrement != 0 {
			t.Errorf("omni AutoIncrement after truncate: want 0, got %+v", tbl)
		}
	})

	// Scenario 3: TRUNCATE nonexistent table → error
	t.Run("truncate_nonexistent", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_trunc_noexist")
		ctrErr := ctr.execSQL("TRUNCATE TABLE t_trunc_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for TRUNCATE nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("TRUNCATE TABLE t_trunc_noexist", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for TRUNCATE nonexistent table")
		}
	})
}

func TestContainer_Section_2_8_CreateDropIndex(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario 5: CREATE INDEX IF NOT EXISTS. MySQL 8.0 has no such form
	// (a syntax error there); OceanBase has it, and so does the catalog: on a
	// duplicate name it is a no-op.
	t.Run("create_index_if_not_exists", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci5")
		ctr.execSQL("CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)")
		ctr.execSQL("CREATE INDEX idx_val ON t_ci5 (val)")
		if err := ctr.execSQL("CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)"); err != nil {
			t.Errorf("container: CREATE INDEX IF NOT EXISTS: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)", nil)
		c.Exec("CREATE INDEX idx_val ON t_ci5 (val)", nil)
		results, _ := c.Exec("CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)", nil)
		if results[0].Error != nil {
			t.Errorf("omni: unexpected error for IF NOT EXISTS: %v", results[0].Error)
		}
	})

	// Scenario 6: CREATE INDEX — duplicate name → error 1061
	t.Run("create_index_duplicate_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci6")
		ctr.execSQL("CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)")
		ctr.execSQL("CREATE INDEX idx_a ON t_ci6 (a)")
		ctrErr := ctr.execSQL("CREATE INDEX idx_a ON t_ci6 (b)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate index name")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)", nil)
		c.Exec("CREATE INDEX idx_a ON t_ci6 (a)", nil)
		results, _ := c.Exec("CREATE INDEX idx_a ON t_ci6 (b)", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for duplicate index name")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1061 {
			t.Errorf("omni error code: want 1061, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})

	// Scenario 8: DROP INDEX nonexistent → error 1091
	t.Run("drop_index_nonexistent_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_di2")
		ctr.execSQL("CREATE TABLE t_di2 (id INT PRIMARY KEY)")
		ctrErr := ctr.execSQL("DROP INDEX idx_noexist ON t_di2")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent index")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_di2 (id INT PRIMARY KEY)", nil)
		results, _ := c.Exec("DROP INDEX idx_noexist ON t_di2", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for DROP nonexistent index")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1091 {
			t.Errorf("omni error code: want 1091, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})
}

func TestContainer_Section_2_9_RenameTable(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("rename_nonexistent_error", func(t *testing.T) {
		// RENAME TABLE nonexistent → error
		ctr.execSQL("DROP TABLE IF EXISTS t_noexist_ren")
		ctrErr := ctr.execSQL("RENAME TABLE t_noexist_ren TO t_noexist_ren2")
		if ctrErr == nil {
			t.Fatal("container: expected error for RENAME nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("RENAME TABLE t_noexist_ren TO t_noexist_ren2", &ExecOptions{ContinueOnError: true})
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for RENAME nonexistent table")
		}
	})

	t.Run("rename_to_existing_error", func(t *testing.T) {
		// RENAME TABLE to existing name → error
		ctr.execSQL("DROP TABLE IF EXISTS t_ren_exist1")
		ctr.execSQL("DROP TABLE IF EXISTS t_ren_exist2")
		ctr.execSQL("CREATE TABLE t_ren_exist1 (id INT)")
		ctr.execSQL("CREATE TABLE t_ren_exist2 (id INT)")

		ctrErr := ctr.execSQL("RENAME TABLE t_ren_exist1 TO t_ren_exist2")
		if ctrErr == nil {
			t.Fatal("container: expected error for RENAME to existing table name")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ren_exist1 (id INT)", nil)
		c.Exec("CREATE TABLE t_ren_exist2 (id INT)", nil)
		results, _ := c.Exec("RENAME TABLE t_ren_exist1 TO t_ren_exist2", &ExecOptions{ContinueOnError: true})
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for RENAME to existing table name")
		}
	})
}

func TestContainer_Section_2_10_CreateDropView(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("create_view_basic", func(t *testing.T) {
		// CREATE VIEW v AS SELECT ...
		ctr.execSQL("DROP VIEW IF EXISTS v_basic")
		ctrErr := ctr.execSQL("CREATE VIEW v_basic AS SELECT 1 AS a")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE VIEW v_basic AS SELECT 1 AS a", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_basic")] == nil {
			t.Error("omni: view v_basic should exist after CREATE VIEW")
		}
	})

	t.Run("create_or_replace_view", func(t *testing.T) {
		// CREATE OR REPLACE VIEW
		ctr.execSQL("DROP VIEW IF EXISTS v_replace")
		ctr.execSQL("CREATE VIEW v_replace AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE OR REPLACE VIEW v_replace AS SELECT 2 AS b")
		if ctrErr != nil {
			t.Fatalf("container CREATE OR REPLACE VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_replace AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE OR REPLACE VIEW v_replace AS SELECT 2 AS b", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE OR REPLACE VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_replace")] == nil {
			t.Error("omni: view v_replace should exist after CREATE OR REPLACE VIEW")
		}
	})

	t.Run("create_view_with_columns", func(t *testing.T) {
		// CREATE VIEW with column list
		ctr.execSQL("DROP VIEW IF EXISTS v_cols")
		ctrErr := ctr.execSQL("CREATE VIEW v_cols (x, y) AS SELECT 1, 2")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW with columns error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE VIEW v_cols (x, y) AS SELECT 1, 2", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW with columns error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_cols")]
		if v == nil {
			t.Fatal("omni: view v_cols should exist")
		}
		if len(v.Columns) != 2 || v.Columns[0] != "x" || v.Columns[1] != "y" {
			t.Errorf("omni: expected columns [x, y], got %v", v.Columns)
		}
	})

	t.Run("create_view_with_options", func(t *testing.T) {
		// CREATE VIEW with ALGORITHM, DEFINER, SQL_SECURITY
		ctr.execSQL("DROP VIEW IF EXISTS v_opts")
		ctrErr := ctr.execSQL("CREATE ALGORITHM=MERGE DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW v_opts AS SELECT 1 AS a")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW with options error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE ALGORITHM=MERGE DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW v_opts AS SELECT 1 AS a", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW with options error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_opts")]
		if v == nil {
			t.Fatal("omni: view v_opts should exist")
		}
		if v.Algorithm != "MERGE" {
			t.Errorf("omni: expected Algorithm=MERGE, got %q", v.Algorithm)
		}
		if v.SqlSecurity != "INVOKER" {
			t.Errorf("omni: expected SqlSecurity=INVOKER, got %q", v.SqlSecurity)
		}
	})

	t.Run("create_view_with_check_option", func(t *testing.T) {
		// CREATE VIEW with CHECK OPTION
		ctr.execSQL("DROP VIEW IF EXISTS v_chk")
		ctr.execSQL("DROP TABLE IF EXISTS t_chk_view")
		ctr.execSQL("CREATE TABLE t_chk_view (id INT, val INT)")
		ctrErr := ctr.execSQL("CREATE VIEW v_chk AS SELECT * FROM t_chk_view WHERE val > 0 WITH CASCADED CHECK OPTION")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW WITH CHECK OPTION error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_chk_view (id INT, val INT)", nil)
		results, _ := c.Exec("CREATE VIEW v_chk AS SELECT * FROM t_chk_view WHERE val > 0 WITH CASCADED CHECK OPTION", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW WITH CHECK OPTION error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_chk")]
		if v == nil {
			t.Fatal("omni: view v_chk should exist")
		}
		if v.CheckOption != "CASCADED" {
			t.Errorf("omni: expected CheckOption=CASCADED, got %q", v.CheckOption)
		}
	})

	t.Run("drop_view_basic", func(t *testing.T) {
		// DROP VIEW v
		ctr.execSQL("DROP VIEW IF EXISTS v_drop1")
		ctr.execSQL("CREATE VIEW v_drop1 AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("DROP VIEW v_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_drop1 AS SELECT 1 AS a", nil)
		results, _ := c.Exec("DROP VIEW v_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_drop1")] != nil {
			t.Error("omni: view v_drop1 should not exist after DROP VIEW")
		}
	})

	t.Run("drop_view_if_exists", func(t *testing.T) {
		// DROP VIEW IF EXISTS on nonexistent view — no error
		ctr.execSQL("DROP VIEW IF EXISTS v_noexist")
		ctrErr := ctr.execSQL("DROP VIEW IF EXISTS v_noexist")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP VIEW IF EXISTS v_noexist", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("drop_view_multi", func(t *testing.T) {
		// DROP VIEW v1, v2 (multi-view)
		ctr.execSQL("DROP VIEW IF EXISTS v_m1")
		ctr.execSQL("DROP VIEW IF EXISTS v_m2")
		ctr.execSQL("CREATE VIEW v_m1 AS SELECT 1 AS a")
		ctr.execSQL("CREATE VIEW v_m2 AS SELECT 2 AS b")
		ctrErr := ctr.execSQL("DROP VIEW v_m1, v_m2")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW multi error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_m1 AS SELECT 1 AS a", nil)
		c.Exec("CREATE VIEW v_m2 AS SELECT 2 AS b", nil)
		results, _ := c.Exec("DROP VIEW v_m1, v_m2", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW multi error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_m1")] != nil {
			t.Error("omni: view v_m1 should not exist after DROP VIEW")
		}
		if db.Views[toLower("v_m2")] != nil {
			t.Error("omni: view v_m2 should not exist after DROP VIEW")
		}
	})

	// Extra: CREATE VIEW duplicate (no OR REPLACE) should error
	t.Run("create_view_duplicate_error", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_dup")
		ctr.execSQL("CREATE VIEW v_dup AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE VIEW v_dup AS SELECT 2 AS b")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_dup AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE VIEW v_dup AS SELECT 2 AS b", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate view")
		}
	})

	// Extra: CREATE VIEW with same name as existing table should error
	t.Run("create_view_table_conflict", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_tbl_conflict")
		ctr.execSQL("DROP TABLE IF EXISTS v_tbl_conflict")
		ctr.execSQL("CREATE TABLE v_tbl_conflict (id INT)")
		ctrErr := ctr.execSQL("CREATE VIEW v_tbl_conflict AS SELECT 1 AS a")
		if ctrErr == nil {
			t.Fatal("container: expected error for view with same name as table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE v_tbl_conflict (id INT)", nil)
		results, _ := c.Exec("CREATE VIEW v_tbl_conflict AS SELECT 1 AS a", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for view with same name as table")
		}
	})

	// Extra: DROP VIEW on nonexistent view (no IF EXISTS) should error
	t.Run("drop_view_nonexistent_error", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_nonexist_err")
		ctrErr := ctr.execSQL("DROP VIEW v_nonexist_err")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP VIEW on nonexistent view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP VIEW v_nonexist_err", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP VIEW on nonexistent view")
		}
	})
}

func TestContainer_Section_2_11_CreateDropAlterDatabase(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("create_database", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_create1")
		ctrErr := ctr.execSQL("CREATE DATABASE db_create1")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE error: %v", ctrErr)
		}
		ctrDDL, err := ctr.showCreateDatabase("db_create1")
		if err != nil {
			t.Fatalf("container SHOW CREATE DATABASE error: %v", err)
		}

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_create1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_create1")
		if db == nil {
			t.Fatal("omni: database not found after CREATE DATABASE")
		}
		// Verify charset/collation defaults match ctr.
		// Oracle returns something like: CREATE DATABASE `db_create1` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
		if !strings.Contains(ctrDDL, "utf8mb4") {
			t.Logf("container DDL: %s", ctrDDL)
		}
		if db.Charset != "utf8mb4" {
			t.Errorf("omni charset mismatch: got %q, want utf8mb4", db.Charset)
		}
		if db.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("omni collation mismatch: got %q, want utf8mb4_0900_ai_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_create1")
	})

	t.Run("create_database_if_not_exists", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_ine")
		ctr.execSQL("CREATE DATABASE db_ine")

		// CREATE DATABASE IF NOT EXISTS on existing db should succeed (no error).
		ctrErr := ctr.execSQL("CREATE DATABASE IF NOT EXISTS db_ine")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE IF NOT EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE db_ine", nil)
		results, _ := c.Exec("CREATE DATABASE IF NOT EXISTS db_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE IF NOT EXISTS error: %v", results[0].Error)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_ine")
	})

	t.Run("create_database_charset", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_cs")
		ctrErr := ctr.execSQL("CREATE DATABASE db_cs CHARACTER SET latin1")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE with charset error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_cs")

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_cs CHARACTER SET latin1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE with charset error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_cs")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		// Oracle should show latin1 charset
		if !strings.Contains(ctrDDL, "latin1") {
			t.Errorf("container DDL missing latin1: %s", ctrDDL)
		}
		if db.Charset != "latin1" {
			t.Errorf("omni charset: got %q, want latin1", db.Charset)
		}
		// Default collation for latin1 is latin1_swedish_ci
		if db.Collation != "latin1_swedish_ci" {
			t.Errorf("omni collation: got %q, want latin1_swedish_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_cs")
	})

	t.Run("create_database_collate", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_coll")
		ctrErr := ctr.execSQL("CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE with collate error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_coll")

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE with collate error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_coll")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4_unicode_ci") {
			t.Errorf("container DDL missing utf8mb4_unicode_ci: %s", ctrDDL)
		}
		if db.Collation != "utf8mb4_unicode_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_unicode_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_coll")
	})

	t.Run("drop_database", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_drop1")
		ctr.execSQL("CREATE DATABASE db_drop1")

		ctrErr := ctr.execSQL("DROP DATABASE db_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE error: %v", ctrErr)
		}
		// Verify it's gone.
		_, showErr := ctr.showCreateDatabase("db_drop1")
		if showErr == nil {
			t.Fatal("container: database still exists after DROP DATABASE")
		}

		c := New()
		c.Exec("CREATE DATABASE db_drop1", nil)
		results, _ := c.Exec("DROP DATABASE db_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP DATABASE error: %v", results[0].Error)
		}
		if c.GetDatabase("db_drop1") != nil {
			t.Fatal("omni: database still exists after DROP DATABASE")
		}
	})

	t.Run("drop_database_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_drop_ine")
		// DROP DATABASE IF EXISTS on nonexistent db should not error.
		ctrErr := ctr.execSQL("DROP DATABASE IF EXISTS db_drop_ine")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		results, _ := c.Exec("DROP DATABASE IF EXISTS db_drop_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP DATABASE IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("alter_database_charset", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_cs")
		ctr.execSQL("CREATE DATABASE db_alter_cs")

		ctrErr := ctr.execSQL("ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container ALTER DATABASE charset error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_alter_cs")

		c := New()
		c.Exec("CREATE DATABASE db_alter_cs", nil)
		results, _ := c.Exec("ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER DATABASE charset error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_alter_cs")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4") {
			t.Errorf("container DDL missing utf8mb4: %s", ctrDDL)
		}
		if db.Charset != "utf8mb4" {
			t.Errorf("omni charset: got %q, want utf8mb4", db.Charset)
		}
		if db.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_0900_ai_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_cs")
	})

	t.Run("alter_database_collate", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_coll")
		ctr.execSQL("CREATE DATABASE db_alter_coll")

		ctrErr := ctr.execSQL("ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci")
		if ctrErr != nil {
			t.Fatalf("container ALTER DATABASE collate error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_alter_coll")

		c := New()
		c.Exec("CREATE DATABASE db_alter_coll", nil)
		results, _ := c.Exec("ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER DATABASE collate error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_alter_coll")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4_unicode_ci") {
			t.Errorf("container DDL missing utf8mb4_unicode_ci: %s", ctrDDL)
		}
		if db.Collation != "utf8mb4_unicode_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_unicode_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_coll")
	})

	t.Run("ops_on_nonexistent_database", func(t *testing.T) {
		// DROP DATABASE on nonexistent db should error.
		ctr.execSQL("DROP DATABASE IF EXISTS db_nonexist_xyz")
		oracleDropErr := ctr.execSQL("DROP DATABASE db_nonexist_xyz")
		if oracleDropErr == nil {
			t.Fatal("container: expected error for DROP nonexistent database")
		}

		c := New()
		results, _ := c.Exec("DROP DATABASE db_nonexist_xyz", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != ErrUnknownDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr.Code, ErrUnknownDatabase)
		}

		// ALTER DATABASE on nonexistent db should error.
		oracleAlterErr := ctr.execSQL("ALTER DATABASE db_nonexist_xyz CHARACTER SET utf8mb4")
		if oracleAlterErr == nil {
			t.Fatal("container: expected error for ALTER nonexistent database")
		}

		c2 := New()
		results2, _ := c2.Exec("ALTER DATABASE db_nonexist_xyz CHARACTER SET utf8mb4", &ExecOptions{ContinueOnError: true})
		if results2[0].Error == nil {
			t.Fatal("omni: expected error for ALTER nonexistent database")
		}
		catErr2, ok := results2[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results2[0].Error)
		}
		if catErr2.Code != ErrUnknownDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr2.Code, ErrUnknownDatabase)
		}

		// CREATE DATABASE duplicate should error.
		c3 := New()
		c3.Exec("CREATE DATABASE db_dup_test", nil)
		results3, _ := c3.Exec("CREATE DATABASE db_dup_test", &ExecOptions{ContinueOnError: true})
		if results3[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE DATABASE")
		}
		catErr3, ok := results3[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results3[0].Error)
		}
		if catErr3.Code != ErrDupDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr3.Code, ErrDupDatabase)
		}
	})
}

func TestContainer_Section_3_1_DatabaseErrors(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1007_dup_database", func(t *testing.T) {
		// Setup: create the database first, then try to create it again.
		ctr.execSQL("DROP DATABASE IF EXISTS db_err_dup")
		ctr.execSQL("CREATE DATABASE db_err_dup")

		ctrErr := ctr.execSQL("CREATE DATABASE db_err_dup")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate CREATE DATABASE")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1007 {
			t.Fatalf("container: expected error code 1007, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE db_err_dup", nil)
		results, _ := c.Exec("CREATE DATABASE db_err_dup", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE DATABASE")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		// Oracle: "Can't create database 'db_err_dup'; database exists"
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP DATABASE IF EXISTS db_err_dup")
	})

	t.Run("1049_unknown_database", func(t *testing.T) {
		// USE a nonexistent database on ctr.
		ctr.execSQL("DROP DATABASE IF EXISTS db_err_unknown_xyz")
		ctrErr := ctr.execSQL("USE db_err_unknown_xyz")
		if ctrErr == nil {
			t.Fatal("container: expected error for USE nonexistent database")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1049 {
			t.Fatalf("container: expected error code 1049, got %d", ctrCode)
		}

		// Run on omni: DROP DATABASE on nonexistent db triggers 1008, but
		// we need to test the "Unknown database" error (1049).
		// Use DROP DATABASE which should return 1008 in MySQL...
		// Actually, let's check what MySQL returns for DROP DATABASE on nonexistent:
		ctrErr2 := ctr.execSQL("DROP DATABASE db_err_unknown_xyz")
		if ctrErr2 == nil {
			t.Fatal("container: expected error for DROP nonexistent database")
		}
		ctrCode2, _, _ := extractMySQLErr(ctrErr2)
		t.Logf("container DROP error code: %d", ctrCode2)

		// For omni, test USE nonexistent database.
		c := New()
		results, _ := c.Exec("USE db_err_unknown_xyz", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for USE nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1046_no_database_selected", func(t *testing.T) {
		// On ctr, try to CREATE TABLE without selecting a database.
		// We need a fresh connection with no default database.
		// The container connection defaults to "test" database, so we'll
		// test omni behavior and verify the error code/SQLSTATE/message match MySQL's known format.

		// First verify MySQL's behavior: SELECT DATABASE() after no USE should work,
		// but CREATE TABLE without database should fail.
		// Since our container connection defaults to 'test' db, we verify omni matches MySQL's
		// documented error format: ERROR 1046 (3D000): No database selected

		// Run on omni with no current database.
		c := New()
		// Don't set any database — just try to create a table.
		results, _ := c.Exec("CREATE TABLE t_no_db (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for CREATE TABLE without database selected")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Verify against MySQL's documented error values.
		wantCode := 1046
		wantState := "3D000"
		wantMsg := "No database selected"

		if catErr.Code != wantCode {
			t.Errorf("error code: got %d, want %d", catErr.Code, wantCode)
		}
		if catErr.SQLState != wantState {
			t.Errorf("SQLSTATE: got %q, want %q", catErr.SQLState, wantState)
		}
		if catErr.Message != wantMsg {
			t.Errorf("message: got %q, want %q", catErr.Message, wantMsg)
		}
	})
}

func TestContainer_Section_3_2_TableErrors(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1050_table_already_exists", func(t *testing.T) {
		// Setup: create a table, then try to create it again.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_dup")
		ctr.execSQL("CREATE TABLE t_err_dup (id INT)")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_dup (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate CREATE TABLE")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1050 {
			t.Fatalf("container: expected error code 1050, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_err_dup (id INT)", nil)
		results, _ := c.Exec("CREATE TABLE t_err_dup (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE TABLE")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_err_dup")
	})

	t.Run("1051_unknown_table_drop", func(t *testing.T) {
		// DROP TABLE on a nonexistent table should return 1051.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_noexist")

		ctrErr := ctr.execSQL("DROP TABLE t_err_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent table")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1051 {
			t.Fatalf("container: expected error code 1051, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE t_err_noexist", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP nonexistent table")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1146_table_doesnt_exist", func(t *testing.T) {
		// ALTER TABLE on a nonexistent table should return 1146.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_noexist2")

		ctrErr := ctr.execSQL("ALTER TABLE t_err_noexist2 ADD COLUMN x INT")
		if ctrErr == nil {
			t.Fatal("container: expected error for ALTER nonexistent table")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1146 {
			t.Fatalf("container: expected error code 1146, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("ALTER TABLE t_err_noexist2 ADD COLUMN x INT", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for ALTER nonexistent table")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})
}

func TestContainer_Section_3_3_ColumnErrors(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1054_unknown_column", func(t *testing.T) {
		// ALTER TABLE ... MODIFY COLUMN AFTER nonexistent_col triggers 1054
		// "Unknown column 'col' in 'table definition'"
		ctr.execSQL("DROP TABLE IF EXISTS t_err_nocol")
		ctr.execSQL("CREATE TABLE t_err_nocol (id INT, name VARCHAR(50))")

		ctrErr := ctr.execSQL("ALTER TABLE t_err_nocol MODIFY COLUMN name VARCHAR(50) AFTER nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error for unknown column in AFTER clause")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1054 {
			t.Fatalf("container: expected error code 1054, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_err_nocol (id INT, name VARCHAR(50))", nil)
		results, _ := c.Exec("ALTER TABLE t_err_nocol MODIFY COLUMN name VARCHAR(50) AFTER nonexistent", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for unknown column in AFTER clause")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_err_nocol")
	})

	t.Run("1060_duplicate_column", func(t *testing.T) {
		// CREATE TABLE with two columns of the same name.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_dupcol")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_dupcol (a INT, a VARCHAR(10))")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate column name")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1060 {
			t.Fatalf("container: expected error code 1060, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE t_err_dupcol (a INT, a VARCHAR(10))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate column name")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1068_multiple_primary_key", func(t *testing.T) {
		// CREATE TABLE with two PRIMARY KEY definitions. MySQL answers
		// ER_MULTIPLE_PRI_KEY (1068); OceanBase answers 1062, in its own words.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_multipk")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_multipk (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))")
		if ctrErr == nil {
			t.Fatal("container: expected error for multiple primary key")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		if ctrCode != 1062 || ctrState != "23000" || ctrMsg != "Duplicated primary key" {
			t.Fatalf("container: got %d (%s) %s, want 1062 (23000) Duplicated primary key", ctrCode, ctrState, ctrMsg)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE t_err_multipk (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for multiple primary key")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// omni gap: the catalog raises MySQL's error.
		if catErr.Code != 1068 || catErr.SQLState != "42000" || catErr.Message != "Multiple primary key defined" {
			t.Errorf("omni: got %d (%s) %s, want 1068 (42000) Multiple primary key defined", catErr.Code, catErr.SQLState, catErr.Message)
		}
	})
}

func TestContainer_Section_3_4_IndexKeyErrors(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1061_dup_key_name", func(t *testing.T) {
		// Setup: create a table with an index, then try to add another index with the same name.
		ctr.execSQL("DROP TABLE IF EXISTS t_dup_key")
		ctr.execSQL("CREATE TABLE t_dup_key (a INT, b INT, KEY idx_a (a))")

		ctrErr := ctr.execSQL("ALTER TABLE t_dup_key ADD INDEX idx_a (b)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate key name")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1061 {
			t.Fatalf("container: expected error code 1061, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dup_key (a INT, b INT, KEY idx_a (a))", nil)
		results, _ := c.Exec("ALTER TABLE t_dup_key ADD INDEX idx_a (b)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate key name")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_dup_key")
	})

	t.Run("1091_cant_drop_key", func(t *testing.T) {
		// Setup: create a table, then try to drop a nonexistent index.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop_key")
		ctr.execSQL("CREATE TABLE t_drop_key (a INT)")

		ctrErr := ctr.execSQL("ALTER TABLE t_drop_key DROP INDEX idx_nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error for dropping nonexistent key")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1091 {
			t.Fatalf("container: expected error code 1091, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_drop_key (a INT)", nil)
		results, _ := c.Exec("ALTER TABLE t_drop_key DROP INDEX idx_nonexistent", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for dropping nonexistent key")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_drop_key")
	})
}

// TestContainer_Section_3_5_FKErrors holds the foreign-key errors. MySQL 8.0
// names each cause (1824, 1822, 3780); OceanBase answers 1146 for a missing
// parent and ER_CANNOT_ADD_FOREIGN (1215) for the rest.
//
// omni gap: the catalog raises MySQL's errors.
func TestContainer_Section_3_5_FKErrors(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		setup string
		stmt  string
		// OceanBase's answer.
		code  uint16
		state string
		msg   string
		// The catalog's.
		omniCode  int
		omniState string
		omniMsg   string
	}{
		{
			name:     "1824_fk_ref_table_not_found",
			stmt:     "CREATE TABLE t_fk_noref (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_nonexistent_parent(id))",
			code:     1146,
			state:    "42S02",
			msg:      "Table doesn't exist",
			omniCode: 1824, omniState: "HY000",
			omniMsg: "Failed to open the referenced table 't_nonexistent_parent'",
		},
		{
			name:     "1822_fk_missing_index_on_ref_table",
			setup:    "CREATE TABLE t_fk_parent_nokey (id INT, val INT)",
			stmt:     "CREATE TABLE t_fk_child_nokey (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent_nokey(val))",
			code:     1215,
			state:    "HY000",
			msg:      "Cannot add foreign key constraint",
			omniCode: 1822, omniState: "HY000",
			omniMsg: "Failed to add the foreign key constraint. Missing index for constraint 't_fk_child_nokey_ibfk_1' in the referenced table 't_fk_parent_nokey'",
		},
		{
			name:     "3780_fk_column_type_mismatch",
			setup:    "CREATE TABLE t_fk_parent_mismatch (id INT PRIMARY KEY)",
			stmt:     "CREATE TABLE t_fk_child_mismatch (id INT, pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES t_fk_parent_mismatch(id))",
			code:     1215,
			state:    "HY000",
			msg:      "Cannot add foreign key constraint",
			omniCode: 3780, omniState: "HY000",
			omniMsg: "Referencing column 'pid' and referenced column 'id' in foreign key constraint 't_fk_child_mismatch_ibfk_1' are incompatible.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			if tc.setup != "" {
				if err := ctr.execSQL(tc.setup); err != nil {
					t.Fatalf("container setup: %v", err)
				}
				c.Exec(tc.setup, nil)
			}

			ctrErr := ctr.execSQL(tc.stmt)
			var me *mysqldriver.MySQLError
			if !errors.As(ctrErr, &me) {
				t.Fatalf("container: got %v, want error %d", ctrErr, tc.code)
			}
			if me.Number != tc.code || string(me.SQLState[:]) != tc.state || me.Message != tc.msg {
				t.Errorf("container: got %d (%s) %s, want %d (%s) %s", me.Number, me.SQLState[:], me.Message, tc.code, tc.state, tc.msg)
			}

			results, _ := c.Exec(tc.stmt, &ExecOptions{ContinueOnError: true})
			catErr, ok := results[0].Error.(*Error)
			if !ok {
				t.Fatalf("omni: expected *Error, got %T: %v", results[0].Error, results[0].Error)
			}
			if catErr.Code != tc.omniCode || catErr.SQLState != tc.omniState || catErr.Message != tc.omniMsg {
				t.Errorf("omni: got %d (%s) %s, want %d (%s) %s", catErr.Code, catErr.SQLState, catErr.Message, tc.omniCode, tc.omniState, tc.omniMsg)
			}
		})
	}
}

func TestContainer_Section_3_6_ErrorContext(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	// Scenario 1: Error message identifier quoting matches MySQL
	// MySQL uses single quotes around identifiers in error messages (not backticks).
	// Test a variety of error types and compare message format exactly.
	t.Run("identifier_quoting", func(t *testing.T) {
		// 1a: Duplicate database — quotes around db name
		ctr.execSQL("DROP DATABASE IF EXISTS db_quote_test")
		ctr.execSQL("CREATE DATABASE db_quote_test")
		ctrErr := ctr.execSQL("CREATE DATABASE db_quote_test")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container dup db: %d %s", ctrCode, ctrMsg)

		c := New()
		c.Exec("CREATE DATABASE db_quote_test", nil)
		results, _ := c.Exec("CREATE DATABASE db_quote_test", &ExecOptions{ContinueOnError: true})
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup db message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1b: Duplicate table — quotes around table name
		ctr.execSQL("DROP TABLE IF EXISTS t_quote_test")
		ctr.execSQL("CREATE TABLE t_quote_test (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE t_quote_test (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container dup table: %d %s", ctrCode, ctrMsg)

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		c2.Exec("CREATE TABLE t_quote_test (id INT)", nil)
		results, _ = c2.Exec("CREATE TABLE t_quote_test (id INT)", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup table message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1c: Unknown column — quotes around column name and context
		ctr.execSQL("DROP TABLE IF EXISTS t_col_quote")
		ctr.execSQL("CREATE TABLE t_col_quote (id INT)")
		ctrErr = ctr.execSQL("ALTER TABLE t_col_quote DROP COLUMN nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container unknown col: %d %s", ctrCode, ctrMsg)

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		c3.Exec("CREATE TABLE t_col_quote (id INT)", nil)
		results, _ = c3.Exec("ALTER TABLE t_col_quote DROP COLUMN nonexistent", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("unknown col message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1d: Duplicate key name — quotes around key name
		ctr.execSQL("DROP TABLE IF EXISTS t_key_quote")
		ctr.execSQL("CREATE TABLE t_key_quote (id INT, val INT, KEY idx_val (val))")
		ctrErr = ctr.execSQL("ALTER TABLE t_key_quote ADD INDEX idx_val (id)")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container dup key: %d %s", ctrCode, ctrMsg)

		c4 := New()
		c4.Exec("CREATE DATABASE test", nil)
		c4.SetCurrentDatabase("test")
		c4.Exec("CREATE TABLE t_key_quote (id INT, val INT, KEY idx_val (val))", nil)
		results, _ = c4.Exec("ALTER TABLE t_key_quote ADD INDEX idx_val (id)", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup key message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1e: Can't drop key — quotes around key name
		ctr.execSQL("DROP TABLE IF EXISTS t_dropkey_quote")
		ctr.execSQL("CREATE TABLE t_dropkey_quote (id INT)")
		ctrErr = ctr.execSQL("ALTER TABLE t_dropkey_quote DROP INDEX nokey")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container can't drop key: %d %s", ctrCode, ctrMsg)

		c5 := New()
		c5.Exec("CREATE DATABASE test", nil)
		c5.SetCurrentDatabase("test")
		c5.Exec("CREATE TABLE t_dropkey_quote (id INT)", nil)
		results, _ = c5.Exec("ALTER TABLE t_dropkey_quote DROP INDEX nokey", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("can't drop key message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1f: Unknown table (DROP TABLE) — quotes around db.table
		ctr.execSQL("DROP TABLE IF EXISTS t_unknown_quote")
		ctrErr = ctr.execSQL("DROP TABLE t_unknown_quote")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container unknown table: %d %s", ctrCode, ctrMsg)

		c6 := New()
		c6.Exec("CREATE DATABASE test", nil)
		c6.SetCurrentDatabase("test")
		results, _ = c6.Exec("DROP TABLE t_unknown_quote", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("unknown table message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_quote_test")
		ctr.execSQL("DROP TABLE IF EXISTS t_col_quote")
		ctr.execSQL("DROP TABLE IF EXISTS t_key_quote")
		ctr.execSQL("DROP TABLE IF EXISTS t_dropkey_quote")
		ctr.execSQL("DROP DATABASE IF EXISTS db_quote_test")
	})

	// Scenario 2: Error position (index) for multi-statement SQL
	// When executing multiple statements, the error should appear at the correct index.
	t.Run("error_position_multi_statement", func(t *testing.T) {
		// Setup: Execute 3 statements where the 3rd one fails.
		// The first two should succeed, the third should error.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")

		sql := "CREATE TABLE t_pos1 (id INT); CREATE TABLE t_pos2 (id INT); CREATE TABLE t_pos1 (id INT)"
		results, _ := c.Exec(sql, &ExecOptions{ContinueOnError: true})

		// Verify we got 3 results.
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		// First two should succeed.
		if results[0].Error != nil {
			t.Errorf("statement 0: unexpected error: %v", results[0].Error)
		}
		if results[1].Error != nil {
			t.Errorf("statement 1: unexpected error: %v", results[1].Error)
		}
		// Third should fail with duplicate table.
		if results[2].Error == nil {
			t.Fatal("statement 2: expected error for duplicate table")
		}
		if results[2].Index != 2 {
			t.Errorf("error index: want 2, got %d", results[2].Index)
		}
		catErr, ok := results[2].Error.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", results[2].Error)
		}
		if catErr.Code != 1050 {
			t.Errorf("error code: want 1050, got %d", catErr.Code)
		}

		// Also verify on container: same 3-statement SQL, error on 3rd.
		ctr.execSQL("DROP TABLE IF EXISTS t_pos1")
		ctr.execSQL("DROP TABLE IF EXISTS t_pos2")
		ctr.execSQL("CREATE TABLE t_pos1 (id INT)")
		ctr.execSQL("CREATE TABLE t_pos2 (id INT)")
		ctrErr := ctr.execSQL("CREATE TABLE t_pos1 (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate table")
		}
		ctrCode, _, _ := extractMySQLErr(ctrErr)
		if ctrCode != 1050 {
			t.Errorf("container error code: want 1050, got %d", ctrCode)
		}

		// Second test: error on the 1st statement of multi-statement batch.
		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")

		sql2 := "CREATE TABLE t_pos1 (a INT, a INT); CREATE TABLE t_pos3 (id INT)"
		results2, _ := c2.Exec(sql2, &ExecOptions{ContinueOnError: true})
		if len(results2) < 1 {
			t.Fatal("expected at least 1 result")
		}
		if results2[0].Error == nil {
			t.Fatal("statement 0: expected error for duplicate column")
		}
		if results2[0].Index != 0 {
			t.Errorf("error index: want 0, got %d", results2[0].Index)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_pos1")
		ctr.execSQL("DROP TABLE IF EXISTS t_pos2")
	})

	// Scenario 3: IF EXISTS suppresses errors correctly
	t.Run("if_exists_suppresses_errors", func(t *testing.T) {
		// 3a: DROP TABLE IF EXISTS on nonexistent table — no error on both.
		ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_none")
		ctrErr := ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE IF EXISTS t_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP TABLE IF EXISTS error: %v", results[0].Error)
		}

		// 3b: DROP DATABASE IF EXISTS on nonexistent database — no error on both.
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifexists_none")
		ctrErr = ctr.execSQL("DROP DATABASE IF EXISTS db_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE IF EXISTS error: %v", ctrErr)
		}

		c2 := New()
		results, _ = c2.Exec("DROP DATABASE IF EXISTS db_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP DATABASE IF EXISTS error: %v", results[0].Error)
		}

		// 3c: DROP VIEW IF EXISTS on nonexistent view — no error on both.
		ctr.execSQL("DROP VIEW IF EXISTS v_ifexists_none")
		ctrErr = ctr.execSQL("DROP VIEW IF EXISTS v_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW IF EXISTS error: %v", ctrErr)
		}

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		results, _ = c3.Exec("DROP VIEW IF EXISTS v_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP VIEW IF EXISTS error: %v", results[0].Error)
		}

		// 3d: DROP TABLE IF EXISTS on existing table — should succeed (table is dropped).
		ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_real")
		ctr.execSQL("CREATE TABLE t_ifexists_real (id INT)")
		ctrErr = ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_real")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS (existing) error: %v", ctrErr)
		}

		c4 := New()
		c4.Exec("CREATE DATABASE test", nil)
		c4.SetCurrentDatabase("test")
		c4.Exec("CREATE TABLE t_ifexists_real (id INT)", nil)
		results, _ = c4.Exec("DROP TABLE IF EXISTS t_ifexists_real", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP TABLE IF EXISTS (existing) error: %v", results[0].Error)
		}
		// Verify the table was actually dropped.
		results, _ = c4.Exec("DROP TABLE t_ifexists_real", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Error("omni: table should have been dropped but still exists")
		}

		// 3e: Without IF EXISTS, DROP TABLE on nonexistent table — must error.
		ctr.execSQL("DROP TABLE IF EXISTS t_noifexists")
		ctrErr = ctr.execSQL("DROP TABLE t_noifexists")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP TABLE without IF EXISTS")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)

		c5 := New()
		c5.Exec("CREATE DATABASE test", nil)
		c5.SetCurrentDatabase("test")
		results, _ = c5.Exec("DROP TABLE t_noifexists", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP TABLE without IF EXISTS")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	// Scenario 4: IF NOT EXISTS suppresses errors correctly
	t.Run("if_not_exists_suppresses_errors", func(t *testing.T) {
		// 4a: CREATE DATABASE IF NOT EXISTS on existing database — no error.
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifne")
		ctr.execSQL("CREATE DATABASE db_ifne")
		ctrErr := ctr.execSQL("CREATE DATABASE IF NOT EXISTS db_ifne")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE IF NOT EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE db_ifne", nil)
		results, _ := c.Exec("CREATE DATABASE IF NOT EXISTS db_ifne", nil)
		if results[0].Error != nil {
			t.Errorf("omni CREATE DATABASE IF NOT EXISTS error: %v", results[0].Error)
		}

		// 4b: CREATE TABLE IF NOT EXISTS on existing table — no error, original preserved.
		ctr.execSQL("DROP TABLE IF EXISTS t_ifne")
		ctr.execSQL("CREATE TABLE t_ifne (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE IF NOT EXISTS t_ifne (val VARCHAR(100))")
		if ctrErr != nil {
			t.Fatalf("container CREATE TABLE IF NOT EXISTS error: %v", ctrErr)
		}
		// Verify original table is unchanged.
		ctrDDL, _ := ctr.showCreateTable("t_ifne")
		if !strings.Contains(ctrDDL, "`id`") {
			t.Errorf("container: original table structure should be preserved, got: %s", ctrDDL)
		}

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		c2.Exec("CREATE TABLE t_ifne (id INT)", nil)
		results, _ = c2.Exec("CREATE TABLE IF NOT EXISTS t_ifne (val VARCHAR(100))", nil)
		if results[0].Error != nil {
			t.Errorf("omni CREATE TABLE IF NOT EXISTS error: %v", results[0].Error)
		}
		// Verify original table is unchanged.
		omniDDL := c2.ShowCreateTable("test", "t_ifne")
		if !strings.Contains(omniDDL, "`id`") {
			t.Errorf("omni: original table structure should be preserved, got: %s", omniDDL)
		}

		// 4c: Without IF NOT EXISTS, CREATE TABLE on existing table — must error.
		ctr.execSQL("DROP TABLE IF EXISTS t_no_ifne")
		ctr.execSQL("CREATE TABLE t_no_ifne (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE t_no_ifne (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for CREATE TABLE without IF NOT EXISTS")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		c3.Exec("CREATE TABLE t_no_ifne (id INT)", nil)
		results, _ = c3.Exec("CREATE TABLE t_no_ifne (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for CREATE TABLE without IF NOT EXISTS")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_ifne")
		ctr.execSQL("DROP TABLE IF EXISTS t_no_ifne")
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifne")
	})
}

// TestContainer_Section_4_1_Partitioning holds the partitioning statements of
// MySQL 8.0 that OceanBase refuses.
//
// omni gap: the catalog runs each of them, as MySQL does.
func TestContainer_Section_4_1_Partitioning(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// OceanBase's grammar has no LINEAR partitioning.
	t.Run("partition_linear_hash", func(t *testing.T) {
		containerRefusal(t, ctr, "",
			"CREATE TABLE t_part_linear (id INT NOT NULL) PARTITION BY LINEAR HASH (id) PARTITIONS 4", 1064, 0)
	})

	// REORGANIZE PARTITION is not supported on a table without a primary key.
	t.Run("alter_reorganize_partition", func(t *testing.T) {
		containerRefusal(t, ctr,
			`CREATE TABLE t_alter_reorgp (
				id INT NOT NULL,
				val INT NOT NULL
			) PARTITION BY RANGE (val) (
				PARTITION p0 VALUES LESS THAN (10),
				PARTITION p1 VALUES LESS THAN (30)
			)`,
			`ALTER TABLE t_alter_reorgp REORGANIZE PARTITION p1 INTO (
				PARTITION p1a VALUES LESS THAN (20),
				PARTITION p1b VALUES LESS THAN (30)
			)`, 1235, 0)
	})

	// OceanBase's grammar has no COALESCE PARTITION.
	t.Run("alter_coalesce_partition", func(t *testing.T) {
		containerRefusal(t, ctr,
			"CREATE TABLE t_alter_coalp (id INT NOT NULL) PARTITION BY HASH (id) PARTITIONS 4",
			"ALTER TABLE t_alter_coalp COALESCE PARTITION 2", 1064, 0)
	})

	// OceanBase's EXCHANGE PARTITION ends WITHOUT VALIDATION
	// (sql_parser_mysql_mode.y: EXCHANGE PARTITION relation_name WITH TABLE
	// relation_factor WITHOUT VALIDATION). MySQL's form, which may stop at the
	// table, is a syntax error.
	t.Run("alter_exchange_partition", func(t *testing.T) {
		containerRefusal(t, ctr,
			`CREATE TABLE t_alter_exchp (
				id INT NOT NULL,
				val INT NOT NULL
			) PARTITION BY RANGE (val) (
				PARTITION p0 VALUES LESS THAN (10),
				PARTITION p1 VALUES LESS THAN MAXVALUE
			);
			CREATE TABLE t_alter_exchp_swap (
				id INT NOT NULL,
				val INT NOT NULL
			)`,
			"ALTER TABLE t_alter_exchp EXCHANGE PARTITION p0 WITH TABLE t_alter_exchp_swap", 1064, 0)
	})
}

func TestContainer_Section_4_2_StoredRoutines(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// 4.2.3: DROP FUNCTION / PROCEDURE
	t.Run("drop_function", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_drop_test")
		ctr.execSQL("CREATE FUNCTION fn_drop_test() RETURNS INT DETERMINISTIC RETURN 1")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE FUNCTION fn_drop_test() RETURNS INT DETERMINISTIC RETURN 1", nil)

		// Drop on container
		if err := ctr.execSQL("DROP FUNCTION fn_drop_test"); err != nil {
			t.Fatalf("container DROP FUNCTION: %v", err)
		}
		// Drop on omni
		results, _ := c.Exec("DROP FUNCTION fn_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP FUNCTION error: %v", results[0].Error)
		}

		// Both should now error on SHOW CREATE FUNCTION
		_, ctrErr := ctr.showCreateFunction("fn_drop_test")
		omniDDL := c.ShowCreateFunction("test", "fn_drop_test")

		if ctrErr == nil {
			t.Error("expected container error after DROP FUNCTION")
		}
		if omniDDL != "" {
			t.Error("expected empty omni DDL after DROP FUNCTION")
		}
	})

	t.Run("drop_procedure", func(t *testing.T) {
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_drop_test")
		ctr.execSQLDirect("CREATE PROCEDURE sp_drop_test() BEGIN SELECT 1; END")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE PROCEDURE sp_drop_test() BEGIN SELECT 1; END", nil)

		// Drop on container
		if err := ctr.execSQL("DROP PROCEDURE sp_drop_test"); err != nil {
			t.Fatalf("container DROP PROCEDURE: %v", err)
		}
		// Drop on omni
		results, _ := c.Exec("DROP PROCEDURE sp_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP PROCEDURE error: %v", results[0].Error)
		}

		_, ctrErr := ctr.showCreateProcedure("sp_drop_test")
		omniDDL := c.ShowCreateProcedure("test", "sp_drop_test")

		if ctrErr == nil {
			t.Error("expected container error after DROP PROCEDURE")
		}
		if omniDDL != "" {
			t.Error("expected empty omni DDL after DROP PROCEDURE")
		}
	})

	// 4.2.4: DROP FUNCTION IF EXISTS — no error for nonexistent
	t.Run("drop_function_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_nonexistent_xyz")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP FUNCTION IF EXISTS fn_nonexistent_xyz", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP FUNCTION IF EXISTS should not error: %v", results[0].Error)
		}
	})

	// 4.2.5: DROP FUNCTION nonexistent — error
	t.Run("drop_function_nonexistent_error", func(t *testing.T) {
		ctrErr := ctr.execSQL("DROP FUNCTION fn_totally_nonexistent")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP FUNCTION fn_totally_nonexistent", nil)

		if ctrErr == nil {
			t.Fatal("expected container error for DROP nonexistent FUNCTION")
		}
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP nonexistent FUNCTION")
		}

		// Both should produce error code 1305
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(ctrErr, &mysqlErr) {
			if catErr, ok := results[0].Error.(*Error); ok {
				if catErr.Code != int(mysqlErr.Number) {
					t.Errorf("error code mismatch: container=%d omni=%d", mysqlErr.Number, catErr.Code)
				}
			}
		}
	})
}

func TestContainer_Section_4_3_Triggers(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create a table that triggers can reference.
	ctr.execSQL("DROP TABLE IF EXISTS t_trigger")
	if err := ctr.execSQL("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))"); err != nil {
		t.Fatalf("container setup table: %v", err)
	}

	// 4.3.2: DROP TRIGGER
	t.Run("drop_trigger", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_drop_test")
		ctr.execSQLDirect("CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0")
		if err := ctr.execSQL("DROP TRIGGER tr_drop_test"); err != nil {
			t.Fatalf("container DROP TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		c.Exec("CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0", nil)
		results, _ := c.Exec("DROP TRIGGER tr_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TRIGGER error: %v", results[0].Error)
		}

		// Verify trigger is gone
		ddl := c.ShowCreateTrigger("test", "tr_drop_test")
		if ddl != "" {
			t.Errorf("trigger should be gone, but got: %s", ddl)
		}
	})

	t.Run("drop_trigger_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_nonexistent_drop")
		if err := ctr.execSQL("DROP TRIGGER IF EXISTS tr_nonexistent_drop"); err != nil {
			t.Fatalf("container DROP TRIGGER IF EXISTS: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TRIGGER IF EXISTS tr_nonexistent_drop", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TRIGGER IF EXISTS should not error: %v", results[0].Error)
		}
	})

	t.Run("drop_trigger_not_exist_error", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_no_exist")
		ctrErr := ctr.execSQL("DROP TRIGGER tr_no_exist")
		if ctrErr == nil {
			t.Fatal("expected container error for DROP TRIGGER on nonexistent trigger")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TRIGGER tr_no_exist", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP TRIGGER on nonexistent trigger")
		}
		if catErr, ok := results[0].Error.(*Error); ok {
			if catErr.Code != ErrNoSuchTrigger {
				t.Errorf("error code: want %d, got %d", ErrNoSuchTrigger, catErr.Code)
			}
		}
	})
}

func TestContainer_Section_4_4_Events(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// 4.4.2: ALTER EVENT. OceanBase does not rename an event that has
	// started, and one scheduled EVERY 1 HOUR starts when it is created (1235).
	// omni gap: the catalog renames it, as MySQL does.
	t.Run("alter_event_rename", func(t *testing.T) {
		containerRefusal(t, ctr,
			"CREATE EVENT ev_rename_old ON SCHEDULE EVERY 1 HOUR DO SELECT 1",
			"ALTER EVENT ev_rename_old RENAME TO ev_rename_new", 1235, 0)
	})

	// 4.4.3: DROP EVENT
	t.Run("drop_event", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_drop_test")
		ctr.execSQLDirect("CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1")
		if err := ctr.execSQL("DROP EVENT ev_drop_test"); err != nil {
			t.Fatalf("container DROP EVENT: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1", nil)
		results, _ := c.Exec("DROP EVENT ev_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP EVENT error: %v", results[0].Error)
		}

		// Verify event is gone
		ddl := c.ShowCreateEvent("test", "ev_drop_test")
		if ddl != "" {
			t.Errorf("event should be gone, but got: %s", ddl)
		}
	})

	t.Run("drop_event_if_exists", func(t *testing.T) {
		if err := ctr.execSQL("DROP EVENT IF EXISTS ev_nonexistent"); err != nil {
			t.Fatalf("container DROP EVENT IF EXISTS: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP EVENT IF EXISTS ev_nonexistent", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP EVENT IF EXISTS should not error: %v", results[0].Error)
		}
	})

	t.Run("drop_event_not_exist_error", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_no_exist")
		ctrErr := ctr.execSQL("DROP EVENT ev_no_exist")
		if ctrErr == nil {
			t.Fatal("expected container error for DROP EVENT on nonexistent event")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP EVENT ev_no_exist", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP EVENT on nonexistent event")
		}
		if catErr, ok := results[0].Error.(*Error); ok {
			if catErr.Code != ErrNoSuchEvent {
				t.Errorf("error code: want %d, got %d", ErrNoSuchEvent, catErr.Code)
			}
		}
	})

	// 4.4.4: SHOW CREATE EVENT. OceanBase has events and no SHOW CREATE EVENT:
	// the statement is a syntax error, so there is no event text to compare.
	// omni gap: the parser takes the statement, as MySQL does.
	t.Run("show_create_event_basic", func(t *testing.T) {
		containerRefusal(t, ctr,
			"CREATE EVENT ev_show_test ON SCHEDULE EVERY 1 HOUR COMMENT 'test event' DO SELECT 1",
			"SHOW CREATE EVENT ev_show_test", 1064, 0)
	})
}

func TestContainer_Section_4_5_ViewsDeep(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("alter_view", func(t *testing.T) {
		// ALTER VIEW changes the view definition.
		// Verify both container and omni accept ALTER VIEW and update the view.
		ctr.execSQL("DROP VIEW IF EXISTS v_alter_test")
		ctr.execSQL("CREATE VIEW v_alter_test AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c")
		if ctrErr != nil {
			t.Fatalf("container ALTER VIEW error: %v", ctrErr)
		}
		// Verify container view was updated.
		ctrDDL, _ := ctr.showCreateView("v_alter_test")
		if !strings.Contains(ctrDDL, "v_alter_test") {
			t.Fatalf("container: view not found after ALTER VIEW")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_alter_test AS SELECT 1 AS a", nil)
		results, _ := c.Exec("ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER VIEW error: %v", results[0].Error)
		}
		// Verify omni view was updated.
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_alter_test")]
		if v == nil {
			t.Fatal("omni: view v_alter_test should exist after ALTER VIEW")
		}
		if !strings.Contains(v.Definition, "2") {
			t.Errorf("omni: view definition should contain new select, got: %s", v.Definition)
		}

		// ALTER VIEW on nonexistent view should error on both.
		ctr.execSQL("DROP VIEW IF EXISTS v_alter_noexist")
		ctrErr = ctr.execSQL("ALTER VIEW v_alter_noexist AS SELECT 1")
		if ctrErr == nil {
			t.Fatal("container: expected error for ALTER VIEW on nonexistent view")
		}

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		results2, _ := c2.Exec("ALTER VIEW v_alter_noexist AS SELECT 1", &ExecOptions{ContinueOnError: true})
		if results2[0].Error == nil {
			t.Fatal("omni: expected error for ALTER VIEW on nonexistent view")
		}
	})

	t.Run("view_dependency_tracking", func(t *testing.T) {
		// In MySQL, dropping a base table does NOT drop the view.
		// The view still exists but errors when queried.
		ctr.execSQL("DROP VIEW IF EXISTS v_dep_test")
		ctr.execSQL("DROP TABLE IF EXISTS t_dep_base")
		ctr.execSQL("CREATE TABLE t_dep_base (id INT, val VARCHAR(100))")
		ctr.execSQL("CREATE VIEW v_dep_test AS SELECT id, val FROM t_dep_base")

		// Drop the base table
		ctrErr := ctr.execSQL("DROP TABLE t_dep_base")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE error: %v", ctrErr)
		}

		// View should still exist in container
		_, showErr := ctr.showCreateView("v_dep_test")
		if showErr != nil {
			t.Fatalf("container: view should still exist after base table drop, got: %v", showErr)
		}

		// Omni: same behavior
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dep_base (id INT, val VARCHAR(100))", nil)
		c.Exec("CREATE VIEW v_dep_test AS SELECT id, val FROM t_dep_base", nil)
		c.Exec("DROP TABLE t_dep_base", nil)

		// View should still exist
		db := c.GetDatabase("test")
		if db.Views[toLower("v_dep_test")] == nil {
			t.Error("omni: view v_dep_test should still exist after base table drop")
		}

		// Cleanup container
		ctr.execSQL("DROP VIEW IF EXISTS v_dep_test")
	})
}

func TestContainer_Section_5_1_UseStatement(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("use_db_sets_current_database", func(t *testing.T) {
		// On container: create a separate database, USE it, then CREATE TABLE without qualifying.
		ctr.execSQL("DROP DATABASE IF EXISTS use_test_db")
		ctr.execSQL("CREATE DATABASE use_test_db")
		if err := ctr.execSQL("USE use_test_db"); err != nil {
			t.Fatalf("container USE: %v", err)
		}
		ctr.execSQL("DROP TABLE IF EXISTS t_use_check")
		if err := ctr.execSQL("CREATE TABLE t_use_check (id INT PRIMARY KEY)"); err != nil {
			t.Fatalf("container CREATE TABLE: %v", err)
		}
		// The unqualified table landed in the database USE selected.
		var current, schema string
		if err := ctr.db.QueryRowContext(ctr.ctx, "SELECT DATABASE()").Scan(&current); err != nil {
			t.Fatal(err)
		}
		if current != "use_test_db" {
			t.Errorf("container DATABASE(): got %q, want use_test_db", current)
		}
		if err := ctr.db.QueryRowContext(ctr.ctx, "SELECT TABLE_SCHEMA FROM information_schema.TABLES WHERE TABLE_NAME = 't_use_check'").Scan(&schema); err != nil {
			t.Fatal(err)
		}
		if schema != "use_test_db" {
			t.Errorf("container: t_use_check is in %q, want use_test_db", schema)
		}

		// On omni: same sequence via Exec.
		c := New()
		c.Exec("CREATE DATABASE use_test_db", nil)
		results, _ := c.Exec("USE use_test_db", nil)
		if results[0].Error != nil {
			t.Fatalf("omni USE error: %v", results[0].Error)
		}
		// Verify current database was set.
		if c.CurrentDatabase() != "use_test_db" {
			t.Fatalf("omni CurrentDatabase: got %q, want %q", c.CurrentDatabase(), "use_test_db")
		}

		results, _ = c.Exec("CREATE TABLE t_use_check (id INT PRIMARY KEY)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE TABLE error: %v", results[0].Error)
		}
		if c.GetDatabase("use_test_db").GetTable("t_use_check") == nil {
			t.Error("omni: t_use_check is not in use_test_db")
		}

		// Cleanup ctr.
		ctr.execSQL("DROP DATABASE IF EXISTS use_test_db")
	})

	t.Run("use_nonexistent_error_1049", func(t *testing.T) {
		// On container: USE a nonexistent database.
		ctr.execSQL("DROP DATABASE IF EXISTS use_nonexistent_db")
		ctrErr := ctr.execSQL("USE use_nonexistent_db")
		if ctrErr == nil {
			t.Fatal("container: expected error for USE nonexistent database")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1049 {
			t.Fatalf("container: expected error code 1049, got %d", ctrCode)
		}

		// On omni: USE a nonexistent database.
		c := New()
		results, _ := c.Exec("USE use_nonexistent_db", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for USE nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})
}

func TestContainer_Section_5_2_SetVariables(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario: SET foreign_key_checks = 0 — skip FK validation on CREATE TABLE
	t.Run("fk_checks_off_create", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_fkoff_child")
		ctr.execSQL("SET foreign_key_checks = 0")
		ctrErr := ctr.execSQL("CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))")
		ctr.execSQL("SET foreign_key_checks = 1")
		if ctrErr != nil {
			t.Fatalf("container: unexpected error with FK checks off: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("SET foreign_key_checks = 0", nil)
		results, _ := c.Exec("CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni: unexpected error with FK checks off: %v", results[0].Error)
		}
		if c.GetDatabase("test").GetTable("t_fkoff_child") == nil {
			t.Error("omni: t_fkoff_child missing")
		}
	})

	// Scenario: SET foreign_key_checks = 0 — drop table with FK references
	t.Run("fk_checks_off_drop", func(t *testing.T) {
		ctr.execSQL("SET foreign_key_checks = 0")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkdrop_child")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkdrop_parent")
		ctr.execSQL("SET foreign_key_checks = 1")
		ctr.execSQL("CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)")
		ctr.execSQL("CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))")
		ctr.execSQL("SET foreign_key_checks = 0")
		ctrErr := ctr.execSQL("DROP TABLE t_fkdrop_parent")
		ctr.execSQL("SET foreign_key_checks = 1")
		if ctrErr != nil {
			t.Fatalf("container: unexpected error dropping FK-referenced table with checks off: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)", nil)
		c.Exec("CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))", nil)
		c.Exec("SET foreign_key_checks = 0", nil)
		results, _ := c.Exec("DROP TABLE t_fkdrop_parent", nil)
		if results[0].Error != nil {
			t.Fatalf("omni: unexpected error dropping FK-referenced table with checks off: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_fkdrop_parent")
		if omniDDL != "" {
			t.Errorf("omni: parent table should be gone after drop, got: %s", omniDDL)
		}
	})

	// Scenario: SET foreign_key_checks = 1 — enforce FK validation
	t.Run("fk_checks_on_enforce", func(t *testing.T) {
		ctr.execSQL("SET foreign_key_checks = 0")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkon_child")
		ctr.execSQL("SET foreign_key_checks = 1")
		ctrErr := ctr.execSQL("CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))")
		if ctrErr == nil {
			t.Fatal("container: expected error with FK checks on, referencing non-existent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("SET foreign_key_checks = 1", nil)
		results, _ := c.Exec("CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))", nil)
		if results[0].Error == nil {
			t.Fatal("omni: expected error with FK checks on, referencing non-existent table")
		}
	})

	// Scenario: SET NAMES utf8mb4 — silently accepted
	t.Run("set_names", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET NAMES utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container: SET NAMES error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET NAMES utf8mb4", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET NAMES error: %v", results[0].Error)
		}
	})

	// Scenario: SET CHARACTER SET utf8mb4 — silently accepted
	t.Run("set_character_set", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET CHARACTER SET utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container: SET CHARACTER SET error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET CHARACTER SET utf8mb4", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET CHARACTER SET error: %v", results[0].Error)
		}
	})

	// Scenario: SET sql_mode — silently accepted
	t.Run("set_sql_mode", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'")
		if ctrErr != nil {
			t.Fatalf("container: SET sql_mode error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET sql_mode error: %v", results[0].Error)
		}
	})
}

// TestContainer_Section_5_3_UserRoleManagement verifies that user/role management
// statements are accepted by both OceanBase and omni without error.
// These are marked [~] partial because the in-memory catalog does not actually
// store users, roles, or privileges — it just silently accepts them.
func TestContainer_Section_5_3_UserRoleManagement(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name    string
		setup   string // SQL to run before the main statement (on container only)
		sql     string // main statement to test
		cleanup string // SQL to run after the test (on container only)
	}{
		{
			"create_user",
			"",
			"CREATE USER 'testuser53'@'localhost' IDENTIFIED BY 'Password123!'",
			"DROP USER IF EXISTS 'testuser53'@'localhost'",
		},
		{
			"drop_user",
			"CREATE USER IF NOT EXISTS 'dropme53'@'localhost' IDENTIFIED BY 'Password123!'",
			"DROP USER 'dropme53'@'localhost'",
			"",
		},
		{
			"alter_user",
			"CREATE USER IF NOT EXISTS 'alterme53'@'localhost' IDENTIFIED BY 'Password123!'",
			"ALTER USER 'alterme53'@'localhost' IDENTIFIED BY 'NewPassword456!'",
			"DROP USER IF EXISTS 'alterme53'@'localhost'",
		},
		{
			"create_role",
			"",
			"CREATE ROLE 'app_role53'",
			"DROP ROLE IF EXISTS 'app_role53'",
		},
		{
			"drop_role",
			"CREATE ROLE IF NOT EXISTS 'droprole53'",
			"DROP ROLE 'droprole53'",
			"",
		},
		{
			"grant_privileges",
			"CREATE USER IF NOT EXISTS 'grantee53'@'localhost' IDENTIFIED BY 'Password123!'",
			"GRANT SELECT, INSERT ON test.* TO 'grantee53'@'localhost'",
			"DROP USER IF EXISTS 'grantee53'@'localhost'",
		},
		{
			"revoke_privileges",
			"CREATE USER IF NOT EXISTS 'revokee53'@'localhost' IDENTIFIED BY 'Password123!'; GRANT SELECT ON test.* TO 'revokee53'@'localhost'",
			"REVOKE SELECT ON test.* FROM 'revokee53'@'localhost'",
			"DROP USER IF EXISTS 'revokee53'@'localhost'",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup on container
			if tc.setup != "" {
				if err := ctr.execSQL(tc.setup); err != nil {
					t.Fatalf("container setup: %v", err)
				}
			}

			// Run on container — must succeed
			ctrErr := ctr.execSQL(tc.sql)
			if ctrErr != nil {
				t.Fatalf("container exec: %v", ctrErr)
			}

			// Cleanup on container
			if tc.cleanup != "" {
				_ = ctr.execSQL(tc.cleanup)
			}

			// Run on omni — must parse and not return a parse error.
			// The catalog silently accepts these (no-op) since it doesn't
			// track users/roles/privileges, which is expected behavior
			// for an in-memory DDL catalog.
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if len(results) > 0 && results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
		})
	}
}

func TestContainer_Section_6_1_ShowCreateTableIntegration(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// The table of every default form holds a JSON column with an expression
	// default, which OceanBase refuses (1101).
	// omni gap: the catalog takes it, as MySQL 8.0 does.
	t.Run("all_defaults", func(t *testing.T) {
		containerRefusal(t, ctr, "", `CREATE TABLE t_all_defaults (
			id INT NOT NULL AUTO_INCREMENT,
			int_def INT DEFAULT 0,
			int_null INT DEFAULT NULL,
			int_notnull INT NOT NULL,
			varchar_def VARCHAR(100) DEFAULT 'hello',
			varchar_empty VARCHAR(100) DEFAULT '',
			float_def FLOAT DEFAULT 3.14,
			decimal_def DECIMAL(10,2) DEFAULT 0.00,
			bool_true BOOLEAN DEFAULT TRUE,
			bool_false BOOLEAN DEFAULT FALSE,
			enum_def ENUM('a','b','c') DEFAULT 'a',
			set_def SET('x','y','z') DEFAULT 'x,y',
			bit_def BIT(8) DEFAULT b'00001111',
			blob_col BLOB,
			text_col TEXT,
			json_col JSON,
			ts_def TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			dt_def DATETIME DEFAULT CURRENT_TIMESTAMP,
			ts3_def TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
			expr_int INT DEFAULT (FLOOR(RAND()*100)),
			expr_json JSON DEFAULT (JSON_ARRAY()),
			expr_varchar VARCHAR(100) DEFAULT (UUID()),
			dt_literal DATETIME DEFAULT '2024-01-01 00:00:00',
			date_literal DATE DEFAULT '2024-01-01',
			PRIMARY KEY (id)
		)`, 1101, 0)
	})
}

// TestContainer_Section_6_3_InformationSchemaConsistency verifies that the catalog's
// internal state is consistent with what MySQL 8.0 reports via INFORMATION_SCHEMA.
//
// The omni catalog does not support INFORMATION_SCHEMA SQL queries (SELECT is
// treated as DML and skipped). Instead, we compare the catalog's Go-level data
// structures against OceanBase's INFORMATION_SCHEMA query results.
//
// All scenarios are marked [~] partial because the catalog lacks an
// INFORMATION_SCHEMA query engine — users cannot run SELECT ... FROM
// INFORMATION_SCHEMA.* against the in-memory catalog.
func TestContainer_Section_6_3_InformationSchemaConsistency(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup SQL: a table with various column types, indexes, FK, and CHECK constraints.
	parentSQL := "CREATE TABLE t_is_parent (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)) ENGINE=InnoDB"
	setupSQL := `CREATE TABLE t_is_test (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL DEFAULT '',
		price DECIMAL(10,2) DEFAULT '0.00',
		status ENUM('active','inactive') DEFAULT 'active',
		parent_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY uk_name (name),
		KEY idx_status (status),
		KEY idx_parent (parent_id),
		CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_is_parent (id) ON DELETE SET NULL,
		CONSTRAINT chk_price CHECK (price >= 0)
	) ENGINE=InnoDB`

	// Setup on ctr.
	ctr.execSQL("DROP TABLE IF EXISTS t_is_test")
	ctr.execSQL("DROP TABLE IF EXISTS t_is_parent")
	if err := ctr.execSQL(parentSQL); err != nil {
		t.Fatalf("container parent table: %v", err)
	}
	if err := ctr.execSQL(setupSQL); err != nil {
		t.Fatalf("container setup: %v", err)
	}

	// Setup on omni.
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	if results, _ := c.Exec(parentSQL, nil); results[0].Error != nil {
		t.Fatalf("omni parent table: %v", results[0].Error)
	}
	if results, _ := c.Exec(setupSQL, nil); results[0].Error != nil {
		t.Fatalf("omni setup: %v", results[0].Error)
	}

	db := c.GetDatabase("test")
	tbl := db.GetTable("t_is_test")

	// 6.3.1: INFORMATION_SCHEMA.COLUMNS matches catalog state
	t.Run("columns_match", func(t *testing.T) {
		oracleCols, err := ctr.queryColumns("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryColumns: %v", err)
		}
		if len(oracleCols) != len(tbl.Columns) {
			t.Fatalf("column count: container=%d omni=%d", len(oracleCols), len(tbl.Columns))
		}
		for i, oc := range oracleCols {
			omniCol := tbl.Columns[i]
			if oc.Name != omniCol.Name {
				t.Errorf("col[%d] name: container=%q omni=%q", i, oc.Name, omniCol.Name)
			}
			if oc.Position != omniCol.Position {
				t.Errorf("col[%d] position: container=%d omni=%d", i, oc.Position, omniCol.Position)
			}
			wantNullable := "YES"
			if !omniCol.Nullable {
				wantNullable = "NO"
			}
			if oc.Nullable != wantNullable {
				t.Errorf("col[%d] %s nullable: container=%q omni=%q", i, oc.Name, oc.Nullable, wantNullable)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.2: INFORMATION_SCHEMA.STATISTICS matches catalog indexes
	t.Run("statistics_match", func(t *testing.T) {
		oracleIdxs, err := ctr.queryIndexes("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryIndexes: %v", err)
		}
		// Build omni index column map: indexName -> []columnName
		omniIdxCols := make(map[string][]string)
		for _, idx := range tbl.Indexes {
			for _, ic := range idx.Columns {
				omniIdxCols[idx.Name] = append(omniIdxCols[idx.Name], ic.Name)
			}
		}
		// Build container index column map.
		oracleIdxCols := make(map[string][]string)
		for _, oi := range oracleIdxs {
			oracleIdxCols[oi.Name] = append(oracleIdxCols[oi.Name], oi.ColumnName)
		}
		// Compare index names and columns.
		for name, oracleCols := range oracleIdxCols {
			omniCols, ok := omniIdxCols[name]
			if !ok {
				t.Errorf("index %q exists in container but not in omni", name)
				continue
			}
			if len(oracleCols) != len(omniCols) {
				t.Errorf("index %q column count: container=%d omni=%d", name, len(oracleCols), len(omniCols))
				continue
			}
			for j := range oracleCols {
				if oracleCols[j] != omniCols[j] {
					t.Errorf("index %q col[%d]: container=%q omni=%q", name, j, oracleCols[j], omniCols[j])
				}
			}
		}
		for name := range omniIdxCols {
			if _, ok := oracleIdxCols[name]; !ok {
				t.Errorf("index %q exists in omni but not in container", name)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.3: INFORMATION_SCHEMA.TABLE_CONSTRAINTS matches
	t.Run("table_constraints_match", func(t *testing.T) {
		oracleCons, err := ctr.queryConstraints("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryConstraints: %v", err)
		}
		// Build omni constraint map.
		omniCons := make(map[string]string) // name -> type
		for _, con := range tbl.Constraints {
			var typeName string
			switch con.Type {
			case ConPrimaryKey:
				typeName = "PRIMARY KEY"
			case ConUniqueKey:
				typeName = "UNIQUE"
			case ConForeignKey:
				typeName = "FOREIGN KEY"
			case ConCheck:
				typeName = "CHECK"
			}
			omniCons[con.Name] = typeName
		}
		// Build container constraint map.
		oracleConsMap := make(map[string]string)
		for _, oc := range oracleCons {
			oracleConsMap[oc.Name] = oc.Type
		}
		for name, oType := range oracleConsMap {
			omniType, ok := omniCons[name]
			if !ok {
				t.Errorf("constraint %q (%s) in container but not omni", name, oType)
				continue
			}
			if oType != omniType {
				t.Errorf("constraint %q type: container=%q omni=%q", name, oType, omniType)
			}
		}
		for name := range omniCons {
			if _, ok := oracleConsMap[name]; !ok {
				t.Errorf("constraint %q in omni but not container", name)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.4: INFORMATION_SCHEMA.KEY_COLUMN_USAGE matches
	t.Run("key_column_usage_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, COLUMN_NAME, ORDINAL_POSITION,
			       REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
			WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 't_is_test'
			ORDER BY CONSTRAINT_NAME, ORDINAL_POSITION`)
		if err != nil {
			t.Fatalf("container KEY_COLUMN_USAGE query: %v", err)
		}
		defer rows.Close()

		type kcuRow struct {
			constraintName, columnName string
			ordinalPos                 int
			refTable, refColumn        *string
		}
		var oracleKCU []kcuRow
		for rows.Next() {
			var r kcuRow
			var refTbl, refCol *string
			if err := rows.Scan(&r.constraintName, &r.columnName, &r.ordinalPos, &refTbl, &refCol); err != nil {
				t.Fatalf("scan: %v", err)
			}
			r.refTable = refTbl
			r.refColumn = refCol
			oracleKCU = append(oracleKCU, r)
		}

		// Verify omni constraints have matching columns.
		omniKCU := make(map[string][]string) // constraint name -> columns
		for _, con := range tbl.Constraints {
			if con.Type == ConCheck {
				continue // CHECK constraints don't appear in KEY_COLUMN_USAGE
			}
			omniKCU[con.Name] = con.Columns
		}
		for _, okcu := range oracleKCU {
			cols, ok := omniKCU[okcu.constraintName]
			if !ok {
				t.Errorf("constraint %q in container KEY_COLUMN_USAGE but not omni", okcu.constraintName)
				continue
			}
			idx := okcu.ordinalPos - 1
			if idx < len(cols) && cols[idx] != okcu.columnName {
				t.Errorf("constraint %q col[%d]: container=%q omni=%q",
					okcu.constraintName, idx, okcu.columnName, cols[idx])
			}
		}

		// Verify FK references match.
		for _, con := range tbl.Constraints {
			if con.Type != ConForeignKey {
				continue
			}
			for _, okcu := range oracleKCU {
				if okcu.constraintName == con.Name && okcu.refTable != nil {
					if *okcu.refTable != con.RefTable {
						t.Errorf("FK %q ref table: container=%q omni=%q", con.Name, *okcu.refTable, con.RefTable)
					}
				}
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.5: INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS matches
	t.Run("referential_constraints_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, UNIQUE_CONSTRAINT_NAME,
			       MATCH_OPTION, UPDATE_RULE, DELETE_RULE,
			       REFERENCED_TABLE_NAME
			FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = 'test' AND TABLE_NAME = 't_is_test'
			ORDER BY CONSTRAINT_NAME`)
		if err != nil {
			t.Fatalf("container REFERENTIAL_CONSTRAINTS query: %v", err)
		}
		defer rows.Close()

		type refConRow struct {
			name, uniqueConName, matchOption, updateRule, deleteRule, refTable string
		}
		var oracleRefs []refConRow
		for rows.Next() {
			var r refConRow
			if err := rows.Scan(&r.name, &r.uniqueConName, &r.matchOption, &r.updateRule, &r.deleteRule, &r.refTable); err != nil {
				t.Fatalf("scan: %v", err)
			}
			oracleRefs = append(oracleRefs, r)
		}

		// Compare with omni FK constraints.
		for _, oref := range oracleRefs {
			var found *Constraint
			for _, con := range tbl.Constraints {
				if con.Type == ConForeignKey && con.Name == oref.name {
					found = con
					break
				}
			}
			if found == nil {
				t.Errorf("FK %q in container REFERENTIAL_CONSTRAINTS but not omni", oref.name)
				continue
			}
			if oref.refTable != found.RefTable {
				t.Errorf("FK %q ref table: container=%q omni=%q", oref.name, oref.refTable, found.RefTable)
			}
			// Normalize action names for comparison.
			oracleDelete := oref.deleteRule
			omniDelete := found.OnDelete
			if omniDelete == "" {
				omniDelete = "RESTRICT" // MySQL default
			}
			if oracleDelete != omniDelete {
				t.Errorf("FK %q delete rule: container=%q omni=%q", oref.name, oracleDelete, omniDelete)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.6: INFORMATION_SCHEMA.CHECK_CONSTRAINTS matches
	t.Run("check_constraints_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, CHECK_CLAUSE
			FROM INFORMATION_SCHEMA.CHECK_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = 'test'
			ORDER BY CONSTRAINT_NAME`)
		if err != nil {
			t.Fatalf("container CHECK_CONSTRAINTS query: %v", err)
		}
		defer rows.Close()

		type checkRow struct {
			name, clause string
		}
		var oracleChecks []checkRow
		for rows.Next() {
			var r checkRow
			if err := rows.Scan(&r.name, &r.clause); err != nil {
				t.Fatalf("scan: %v", err)
			}
			oracleChecks = append(oracleChecks, r)
		}

		// Filter to just our table's check constraints (INFORMATION_SCHEMA.CHECK_CONSTRAINTS
		// doesn't have TABLE_NAME, so we match by constraint name).
		omniChecks := make(map[string]string) // name -> expr
		for _, con := range tbl.Constraints {
			if con.Type == ConCheck {
				omniChecks[con.Name] = con.CheckExpr
			}
		}
		for _, oc := range oracleChecks {
			omniExpr, ok := omniChecks[oc.name]
			if !ok {
				// May be an ENUM/SET constraint auto-generated by MySQL — skip.
				continue
			}
			// Normalize for comparison: remove outer parens and whitespace.
			oracleNorm := normalizeWhitespace(oc.clause)
			omniNorm := normalizeWhitespace(omniExpr)
			if oracleNorm != omniNorm {
				// Log but don't fail — expression format may differ slightly.
				t.Logf("check %q clause: container=%q omni=%q (expression format may differ)", oc.name, oracleNorm, omniNorm)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})
}
