package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C24 covers Section C24 "SHOW CREATE TABLE skip_gipk / Invisible
// PK" from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md.
//
// MySQL 8.0.30+ supports a Generated Invisible Primary Key (GIPK): when
// `sql_generate_invisible_primary_key=ON` and a CREATE TABLE has no PRIMARY
// KEY declared, MySQL silently inserts a hidden column
// `my_row_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT INVISIBLE` at position 0
// and adds a PRIMARY KEY over it. The presence of this column is hidden from
// SHOW CREATE TABLE / information_schema unless
// `show_gipk_in_create_table_and_information_schema=ON` is also set.
//
// All session settings are issued on the pinned single-conn pool from
// scenarioContainer so they persist for the duration of each subtest. The
// session vars are reset to OFF at the end of each subtest so other workers
// (or later subtests) start from a known state. The same settings are also
// issued on the catalog instance under test, because GIPK is session-state
// driven on both sides.
func TestScenario_C24(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// Helper: set a session variable on the pinned container connection.
	setSession := func(t *testing.T, stmt string) {
		t.Helper()
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Errorf("session set %q: %v", stmt, err)
		}
	}

	setCatalogSession := func(t *testing.T, c *Catalog, stmt string) {
		t.Helper()
		catalogStmt := strings.ReplaceAll(stmt, " = ON", " = 1")
		catalogStmt = strings.ReplaceAll(catalogStmt, " = OFF", " = 0")
		results, err := c.Exec(catalogStmt, nil)
		if err != nil {
			t.Fatalf("catalog session set %q parse/exec: %v", stmt, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("catalog session set %q result error: %v", stmt, r.Error)
			}
		}
	}

	// Helper: restore both GIPK session vars to OFF at end of subtest.
	resetSession := func(t *testing.T) {
		t.Helper()
		setSession(t, "SET SESSION sql_generate_invisible_primary_key = OFF")
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = OFF")
	}

	// -----------------------------------------------------------------
	// 24.1 GIPK omitted from SHOW CREATE TABLE by default
	// -----------------------------------------------------------------
	t.Run("24_1_gipk_hidden_by_default", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)
		c := scenarioNewCatalog(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")
		setCatalogSession(t, c, "SET SESSION sql_generate_invisible_primary_key = ON")
		// Explicitly set the visibility flag OFF for this subtest. Note: in
		// MySQL 8.0.32+ the default for this var flipped to ON in some
		// distributions, so we cannot rely on the session inheriting OFF.
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = OFF")
		setCatalogSession(t, c, "SET SESSION show_gipk_in_create_table_and_information_schema = OFF")

		ddl := "CREATE TABLE t (a INT, b INT)"
		runOnBoth(t, mc, c, ddl)

		// Oracle: SHOW CREATE TABLE under default visibility must NOT include
		// my_row_id.
		hidden := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(hidden, "my_row_id") {
			t.Errorf("oracle 24.1: SHOW CREATE TABLE under default visibility should hide my_row_id, got:\n%s", hidden)
		}

		// information_schema.COLUMNS under default visibility also hides it.
		var colsHidden int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&colsHidden)
		assertIntEq(t, "oracle 24.1 information_schema hidden", colsHidden, 0)

		// Toggle visibility ON: SHOW CREATE TABLE now reveals my_row_id.
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = ON")
		shown := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if !strings.Contains(shown, "my_row_id") {
			t.Errorf("oracle 24.1: SHOW CREATE TABLE under visibility=ON should reveal my_row_id, got:\n%s", shown)
		}
		var colsShown int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&colsShown)
		assertIntEq(t, "oracle 24.1 information_schema visible", colsShown, 1)

		// omni: catalog should have generated my_row_id at position 0 with a
		// PRIMARY KEY index.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.1: table t missing")
			return
		}
		if tbl.GetColumn("my_row_id") == nil {
			t.Errorf("omni 24.1: expected catalog to generate my_row_id GIPK column (sql_generate_invisible_primary_key not honored)")
		}

		omniHidden := c.ShowCreateTable("testdb", "t")
		if strings.Contains(omniHidden, "my_row_id") {
			t.Errorf("omni 24.1: ShowCreateTable should hide GIPK while visibility=OFF, got:\n%s", omniHidden)
		}
		setCatalogSession(t, c, "SET SESSION show_gipk_in_create_table_and_information_schema = ON")
		omniShown := c.ShowCreateTable("testdb", "t")
		if !strings.Contains(omniShown, "my_row_id") {
			t.Errorf("omni 24.1: ShowCreateTable should reveal GIPK while visibility=ON, got:\n%s", omniShown)
		}
	})

	// -----------------------------------------------------------------
	// 24.4 my_row_id name collision with user-defined column
	// -----------------------------------------------------------------
	t.Run("24_4_gipk_name_collision", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")

		// Oracle: CREATE TABLE with user-declared my_row_id under GIPK=ON
		// must fail with ER_GIPK_FAILED_AUTOINC_COLUMN_NAME_RESERVED (4108).
		_, oracleErrOn := mc.db.ExecContext(mc.ctx,
			"CREATE TABLE t (my_row_id INT, a INT)")
		if oracleErrOn == nil {
			t.Errorf("oracle 24.4: expected error creating table with my_row_id while GIPK=ON, got nil")
		} else if !strings.Contains(oracleErrOn.Error(), "my_row_id") {
			t.Errorf("oracle 24.4: error message should mention my_row_id, got: %v", oracleErrOn)
		}

		// omni: catalog should reject the same CREATE TABLE while
		// sql_generate_invisible_primary_key is ON.
		c := scenarioNewCatalog(t)
		setCatalogSession(t, c, "SET SESSION sql_generate_invisible_primary_key = ON")
		results, err := c.Exec("CREATE TABLE t (my_row_id INT, a INT);", nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni 24.4: expected error rejecting user-declared my_row_id under GIPK=ON")
		}

		// Now turn GIPK off and verify both sides accept the same CREATE
		// TABLE. We use a fresh table name to avoid clashing with whatever
		// omni or oracle may have left behind from the failing path.
		setSession(t, "SET SESSION sql_generate_invisible_primary_key = OFF")
		scenarioReset(t, mc)
		c2 := scenarioNewCatalog(t)
		setCatalogSession(t, c2, "SET SESSION sql_generate_invisible_primary_key = OFF")
		runOnBoth(t, mc, c2, "CREATE TABLE t (my_row_id INT, a INT)")

		// Oracle: my_row_id should exist as an ordinary user column.
		var name string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&name)
		assertStringEq(t, "oracle 24.4 my_row_id user column name", name, "my_row_id")

		// omni: catalog should also have it as an ordinary column.
		tbl := c2.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.4: table t missing under GIPK=OFF")
			return
		}
		if tbl.GetColumn("my_row_id") == nil {
			t.Errorf("omni 24.4: expected user column my_row_id under GIPK=OFF")
		}
	})
}
