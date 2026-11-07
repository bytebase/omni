package catalog

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestScenario_C24 verifies that TiDB rejects MySQL's GIPK settings and never
// invents columns or indexes after those rejections. User-declared my_row_id
// columns remain ordinary columns, regardless of attempted setting changes.
func TestScenario_C24(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	setSession := func(t *testing.T, stmt string) {
		t.Helper()
		_, err := mc.db.ExecContext(mc.ctx, stmt)
		var native *gomysql.MySQLError
		if !errors.As(err, &native) || native.Number != 1193 {
			t.Fatalf("TiDB must reject MySQL GIPK setting %q with 1193: %v", stmt, err)
		}
		assertRejectedDDL(t, scenarioNewCatalog(t), stmt, 1193)
	}

	resetSession := func(t *testing.T) {
		t.Helper()
		setSession(t, "SET SESSION sql_generate_invisible_primary_key = OFF")
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = OFF")
	}

	t.Run("24_1_gipk_hidden_by_default", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)
		c := scenarioNewCatalog(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = OFF")

		ddl := "CREATE TABLE t (a INT, b INT)"
		runOnBoth(t, mc, c, ddl)

		hidden := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(hidden, "my_row_id") {
			t.Errorf("oracle 24.1: SHOW CREATE TABLE under default visibility should hide my_row_id, got:\n%s", hidden)
		}

		var colsHidden int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&colsHidden)
		assertIntEq(t, "oracle 24.1 information_schema hidden", colsHidden, 0)

		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = ON")
		shown := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(shown, "my_row_id") || shown != hidden {
			t.Errorf("rejected visibility setting changed SHOW CREATE: before %s, after %s", hidden, shown)
		}
		var colsShown int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&colsShown)
		assertIntEq(t, "native column count after rejected visibility setting", colsShown, 0)

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.1: table t missing")
			return
		}
		if tbl.GetColumn("my_row_id") != nil {
			t.Error("catalog invented a MySQL GIPK column")
		}
		assertScenarioReadback(t, mc, c, "t")
	})

	t.Run("24_2_gipk_column_spec", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)
		c := scenarioNewCatalog(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = ON")

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT)")

		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, EXTRA, COLUMN_KEY
			 FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`)
		if len(rows) != 0 {
			t.Errorf("native engine invented a MySQL GIPK column: %v", rows)
		}

		var firstCol string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION LIMIT 1`,
			&firstCol)
		assertStringEq(t, "native first declared column", firstCol, "a")

		idxRows := oracleRows(t, mc, "SHOW INDEX FROM t")
		foundPK := false
		for _, r := range idxRows {
			if len(r) >= 5 && asString(r[2]) == "PRIMARY" && asString(r[4]) == "my_row_id" {
				foundPK = true
				break
			}
		}
		assertBoolEq(t, "native has no fabricated PRIMARY index", foundPK, false)

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni 24.2: table t missing")
			return
		}
		if tbl.GetColumn("my_row_id") != nil || len(tbl.Columns) != 1 || tbl.Columns[0].Name != "a" {
			t.Errorf("catalog must expose exactly the declared column: %+v", tbl.Columns)
		}
		pkOK := false
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				if len(idx.Columns) == 1 && strings.EqualFold(idx.Columns[0].Name, "my_row_id") {
					pkOK = true
				}
				break
			}
		}
		if pkOK {
			t.Error("catalog invented PRIMARY KEY(my_row_id)")
		}
		assertScenarioReadback(t, mc, c, "t")
	})

	t.Run("24_3_gipk_suppressed_only_by_pk", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)
		c := scenarioNewCatalog(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")
		setSession(t, "SET SESSION show_gipk_in_create_table_and_information_schema = ON")

		runOnBoth(t, mc, c, "CREATE TABLE t1 (id INT PRIMARY KEY, a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t2 (id INT NOT NULL UNIQUE, a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t3 (id INT AUTO_INCREMENT, a INT, PRIMARY KEY (id))")

		var n1, n2, n3 int
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='my_row_id'`,
			&n1)
		assertIntEq(t, "oracle 24.3 t1 has no GIPK", n1, 0)

		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND COLUMN_NAME='my_row_id'`,
			&n2)
		assertIntEq(t, "native UNIQUE NOT NULL table has no GIPK", n2, 0)

		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t3' AND COLUMN_NAME='my_row_id'`,
			&n3)
		assertIntEq(t, "oracle 24.3 t3 has no GIPK", n3, 0)

		t1 := c.GetDatabase("testdb").GetTable("t1")
		if t1 != nil && t1.GetColumn("my_row_id") != nil {
			t.Errorf("omni 24.3: t1 should not have my_row_id (explicit PK present)")
		}
		t2 := c.GetDatabase("testdb").GetTable("t2")
		if t2 == nil {
			t.Errorf("omni 24.3: table t2 missing")
		} else if t2.GetColumn("my_row_id") != nil {
			t.Error("catalog invented a GIPK on a UNIQUE NOT NULL table")
		}
		t3 := c.GetDatabase("testdb").GetTable("t3")
		if t3 != nil && t3.GetColumn("my_row_id") != nil {
			t.Errorf("omni 24.3: t3 should not have my_row_id (table-level PK present)")
		}
		for _, table := range []string{"t1", "t2", "t3"} {
			assertScenarioReadback(t, mc, c, table)
		}
	})

	t.Run("24_4_gipk_name_collision", func(t *testing.T) {
		scenarioReset(t, mc)
		defer resetSession(t)

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = ON")

		_, oracleErrOn := mc.db.ExecContext(mc.ctx,
			"CREATE TABLE t (my_row_id INT, a INT)")
		if oracleErrOn != nil {
			t.Fatalf("native user-declared my_row_id: %v", oracleErrOn)
		}

		c := scenarioNewCatalog(t)
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
		if omniErr != nil {
			t.Fatalf("catalog user-declared my_row_id: %v", omniErr)
		}
		assertScenarioReadback(t, mc, c, "t")

		setSession(t, "SET SESSION sql_generate_invisible_primary_key = OFF")
		scenarioReset(t, mc)
		c2 := scenarioNewCatalog(t)
		runOnBoth(t, mc, c2, "CREATE TABLE t (my_row_id INT, a INT)")

		var name string
		oracleScan(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'`,
			&name)
		assertStringEq(t, "oracle 24.4 my_row_id user column name", name, "my_row_id")

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
