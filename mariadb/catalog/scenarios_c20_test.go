package catalog

import (
	"database/sql"
	"strings"
	"testing"
)

// TestScenario_C20 covers Section C20 "Field-type-specific implicit defaults"
// from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. Each subtest runs
// the scenario's DDL on both a MySQL 8.0 container and omni's catalog, then
// asserts:
//
//   1. The catalog state (Column.Default pointer, Nullable) — no implicit
//      default synthesized into the AST.
//   2. information_schema.COLUMNS.COLUMN_DEFAULT and SHOW CREATE TABLE
//      rendering match between oracle and omni.
//
// For sections 20.6 / 20.8 (error scenarios), both the oracle and omni
// must reject the DDL — we compare only the fact that an error is raised,
// not the exact message.
//
// All failures use t.Error rather than t.Fatal so the whole section runs
// and each omni gap is captured in scenarios_bug_queue/c20.md.
func TestScenario_C20(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 20.1 INT NOT NULL, no DEFAULT → implicit 0 -----------------------
	t.Run("20_1_int_notnull_no_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE c20_1 (id INT NOT NULL)")

		// Oracle: COLUMN_DEFAULT is NULL (no explicit default stored).
		colDef := c20oracleColumnDefault(t, mc, "c20_1", "id")
		if colDef.Valid {
			t.Errorf("oracle: c20_1.id COLUMN_DEFAULT expected NULL, got %q", colDef.String)
		}
		// Oracle SHOW CREATE must NOT render DEFAULT 0.
		my := oracleShow(t, mc, "SHOW CREATE TABLE c20_1")
		if strings.Contains(aLine(my, "`id`"), "DEFAULT") {
			t.Errorf("oracle: c20_1 SHOW CREATE should not render DEFAULT on id; got %q", aLine(my, "`id`"))
		}

		// omni catalog state
		col := c20getColumn(t, c, "c20_1", "id")
		if col == nil {
			return
		}
		if col.Default != nil {
			t.Errorf("omni: c20_1.id Default expected nil, got %q", *col.Default)
		}
		if col.Nullable {
			t.Errorf("omni: c20_1.id Nullable expected false, got true")
		}
		// omni SHOW CREATE must NOT render DEFAULT.
		omniCreate := c.ShowCreateTable("testdb", "c20_1")
		if strings.Contains(aLine(omniCreate, "`id`"), "DEFAULT") {
			t.Errorf("omni: c20_1 SHOW CREATE should not render DEFAULT on id; got %q", aLine(omniCreate, "`id`"))
		}
	})

	// --- 20.2 INT nullable, no DEFAULT → synthesized DEFAULT NULL ---------
	t.Run("20_2_int_nullable_default_null_synthesis", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE c20_2 (id INT)")

		// Oracle: SHOW CREATE renders "`id` int DEFAULT NULL".
		my := oracleShow(t, mc, "SHOW CREATE TABLE c20_2")
		idLine := aLine(my, "`id`")
		if !strings.Contains(idLine, "DEFAULT NULL") {
			t.Errorf("oracle: c20_2 SHOW CREATE expected `DEFAULT NULL` on id; got %q", idLine)
		}
		// information_schema.COLUMN_DEFAULT is NULL (no string default stored).
		colDef := c20oracleColumnDefault(t, mc, "c20_2", "id")
		if colDef.Valid {
			t.Errorf("oracle: c20_2.id COLUMN_DEFAULT expected NULL, got %q", colDef.String)
		}

		// omni catalog: Default nil, Nullable true.
		col := c20getColumn(t, c, "c20_2", "id")
		if col == nil {
			return
		}
		if col.Default != nil {
			t.Errorf("omni: c20_2.id Default expected nil, got %q", *col.Default)
		}
		if !col.Nullable {
			t.Errorf("omni: c20_2.id Nullable expected true, got false")
		}
		// omni deparse must synthesize DEFAULT NULL.
		omniCreate := c.ShowCreateTable("testdb", "c20_2")
		if !strings.Contains(aLine(omniCreate, "`id`"), "DEFAULT NULL") {
			t.Errorf("omni: c20_2 SHOW CREATE expected DEFAULT NULL on id; got %q", aLine(omniCreate, "`id`"))
		}
	})

	// --- 20.3 VARCHAR/CHAR NOT NULL, no DEFAULT → implicit '' -------------
	t.Run("20_3_string_notnull_no_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE c20_3 (name VARCHAR(64) NOT NULL, code CHAR(4) NOT NULL)")

		for _, cname := range []string{"name", "code"} {
			colDef := c20oracleColumnDefault(t, mc, "c20_3", cname)
			if colDef.Valid {
				t.Errorf("oracle: c20_3.%s COLUMN_DEFAULT expected NULL, got %q", cname, colDef.String)
			}
			col := c20getColumn(t, c, "c20_3", cname)
			if col == nil {
				continue
			}
			if col.Default != nil {
				t.Errorf("omni: c20_3.%s Default expected nil, got %q", cname, *col.Default)
			}
			if col.Nullable {
				t.Errorf("omni: c20_3.%s Nullable expected false", cname)
			}
		}

		my := oracleShow(t, mc, "SHOW CREATE TABLE c20_3")
		omniCreate := c.ShowCreateTable("testdb", "c20_3")
		for _, cname := range []string{"`name`", "`code`"} {
			if strings.Contains(aLine(my, cname), "DEFAULT") {
				t.Errorf("oracle: c20_3 %s line should not render DEFAULT; got %q", cname, aLine(my, cname))
			}
			if strings.Contains(aLine(omniCreate, cname), "DEFAULT") {
				t.Errorf("omni: c20_3 %s line should not render DEFAULT; got %q", cname, aLine(omniCreate, cname))
			}
		}
	})

	// --- 20.4 ENUM NOT NULL (default=first) & nullable (default NULL) -----
	t.Run("20_4_enum_notnull_first_value", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE c20_4 (
			status ENUM('active','archived','deleted') NOT NULL,
			kind   ENUM('a','b','c')
		)`)

		// Oracle: status COLUMN_DEFAULT is NULL (catalog property — even
		// though runtime fills in 'active'). kind COLUMN_DEFAULT is NULL too.
		for _, cname := range []string{"status", "kind"} {
			colDef := c20oracleColumnDefault(t, mc, "c20_4", cname)
			if colDef.Valid {
				t.Errorf("oracle: c20_4.%s COLUMN_DEFAULT expected NULL, got %q", cname, colDef.String)
			}
		}

		// omni catalog state
		if col := c20getColumn(t, c, "c20_4", "status"); col != nil {
			if col.Default != nil {
				t.Errorf("omni: c20_4.status Default expected nil, got %q", *col.Default)
			}
			if col.Nullable {
				t.Errorf("omni: c20_4.status Nullable expected false")
			}
		}
		if col := c20getColumn(t, c, "c20_4", "kind"); col != nil {
			if col.Default != nil {
				t.Errorf("omni: c20_4.kind Default expected nil, got %q", *col.Default)
			}
			if !col.Nullable {
				t.Errorf("omni: c20_4.kind Nullable expected true")
			}
		}

		// Oracle SHOW CREATE: status has no DEFAULT, kind has DEFAULT NULL.
		my := oracleShow(t, mc, "SHOW CREATE TABLE c20_4")
		if strings.Contains(aLine(my, "`status`"), "DEFAULT") {
			t.Errorf("oracle: c20_4 status should not render DEFAULT; got %q", aLine(my, "`status`"))
		}
		if !strings.Contains(aLine(my, "`kind`"), "DEFAULT NULL") {
			t.Errorf("oracle: c20_4 kind expected DEFAULT NULL; got %q", aLine(my, "`kind`"))
		}

		omniCreate := c.ShowCreateTable("testdb", "c20_4")
		if strings.Contains(aLine(omniCreate, "`status`"), "DEFAULT") {
			t.Errorf("omni: c20_4 status should not render DEFAULT; got %q", aLine(omniCreate, "`status`"))
		}
		if !strings.Contains(aLine(omniCreate, "`kind`"), "DEFAULT NULL") {
			t.Errorf("omni: c20_4 kind expected DEFAULT NULL; got %q", aLine(omniCreate, "`kind`"))
		}
	})

	// --- 20.5 DATETIME/DATE NOT NULL, no DEFAULT --------------------------
	t.Run("20_5_datetime_notnull_no_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE c20_5 (
			created_at DATETIME NOT NULL,
			birthday   DATE     NOT NULL
		)`)

		// Oracle: COLUMN_DEFAULT is NULL; catalog does not pre-apply zero-date.
		for _, cname := range []string{"created_at", "birthday"} {
			colDef := c20oracleColumnDefault(t, mc, "c20_5", cname)
			if colDef.Valid {
				t.Errorf("oracle: c20_5.%s COLUMN_DEFAULT expected NULL, got %q", cname, colDef.String)
			}
		}

		// Oracle SHOW CREATE: no DEFAULT rendered.
		my := oracleShow(t, mc, "SHOW CREATE TABLE c20_5")
		if strings.Contains(aLine(my, "`created_at`"), "DEFAULT") {
			t.Errorf("oracle: c20_5 created_at should not render DEFAULT; got %q", aLine(my, "`created_at`"))
		}
		if strings.Contains(aLine(my, "`birthday`"), "DEFAULT") {
			t.Errorf("oracle: c20_5 birthday should not render DEFAULT; got %q", aLine(my, "`birthday`"))
		}

		for _, cname := range []string{"created_at", "birthday"} {
			col := c20getColumn(t, c, "c20_5", cname)
			if col == nil {
				continue
			}
			if col.Default != nil {
				t.Errorf("omni: c20_5.%s Default expected nil, got %q", cname, *col.Default)
			}
			if col.Nullable {
				t.Errorf("omni: c20_5.%s Nullable expected false", cname)
			}
		}

		omniCreate := c.ShowCreateTable("testdb", "c20_5")
		if strings.Contains(aLine(omniCreate, "`created_at`"), "DEFAULT") {
			t.Errorf("omni: c20_5 created_at should not render DEFAULT; got %q", aLine(omniCreate, "`created_at`"))
		}
		if strings.Contains(aLine(omniCreate, "`birthday`"), "DEFAULT") {
			t.Errorf("omni: c20_5 birthday should not render DEFAULT; got %q", aLine(omniCreate, "`birthday`"))
		}
	})

	// --- 20.6 BLOB/TEXT/JSON/GEOMETRY literal DEFAULT → ER 1101 -----------
	t.Run("20_6_blob_text_literal_default_rejected", func(t *testing.T) {
		type caze struct {
			name string
			ddl  string
		}
		cases := []caze{
			{"c20_6a", "CREATE TABLE c20_6a (b BLOB DEFAULT 'abc')"},
			{"c20_6b", "CREATE TABLE c20_6b (t TEXT DEFAULT 'hello')"},
			{"c20_6c", "CREATE TABLE c20_6c (g GEOMETRY DEFAULT 'x')"},
			{"c20_6d", "CREATE TABLE c20_6d (j JSON DEFAULT '[]')"},
		}
		for _, tc := range cases {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			_, mysqlErr := mc.db.ExecContext(mc.ctx, tc.ddl)
			if mysqlErr == nil {
				t.Errorf("oracle: expected ER_BLOB_CANT_HAVE_DEFAULT (1101) for %s, got nil", tc.name)
			} else if !strings.Contains(mysqlErr.Error(), "1101") &&
				!strings.Contains(strings.ToLower(mysqlErr.Error()), "can't have a default value") {
				t.Errorf("oracle: expected 1101 for %s, got %v", tc.name, mysqlErr)
			}

			omniErrored := c20execExpectError(t, c, tc.ddl+";")
			if !omniErrored {
				t.Errorf("omni: expected rejection of %s DDL; got success", tc.name)
			}
			// No table row should exist.
			if c20getTable(c, tc.name) != nil {
				t.Errorf("omni: %s should not be in catalog after rejected DDL", tc.name)
			}
		}
	})

	// --- 20.7 JSON/BLOB expression DEFAULT (8.0.13+) accepted -------------
	t.Run("20_7_expression_default_accepted", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE c20_7 (
			id INT PRIMARY KEY,
			tags  JSON      DEFAULT (JSON_ARRAY()),
			meta  JSON      DEFAULT (JSON_OBJECT('v', 1)),
			blob1 BLOB      DEFAULT (SUBSTRING('abcdef', 1, 3)),
			pt    POINT     DEFAULT (POINT(0, 0)),
			uuid  BINARY(16) DEFAULT (UUID_TO_BIN(UUID()))
		)`
		runOnBoth(t, mc, c, ddl)

		// Oracle: each column must have a non-NULL COLUMN_DEFAULT (the
		// expression text) and EXTRA containing "DEFAULT_GENERATED".
		expressionCols := []string{"tags", "meta", "blob1", "pt", "uuid"}
		for _, cname := range expressionCols {
			var colDef sql.NullString
			var extra string
			oracleScan(t, mc, `SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='`+cname+`'`,
				&colDef, &extra)
			if !colDef.Valid || colDef.String == "" {
				t.Errorf("oracle: c20_7.%s COLUMN_DEFAULT expected non-empty expression text, got NULL", cname)
			}
			if !strings.Contains(extra, "DEFAULT_GENERATED") {
				t.Errorf("oracle: c20_7.%s EXTRA expected DEFAULT_GENERATED, got %q", cname, extra)
			}
		}

		// Oracle: tables exist in omni too.
		if c20getTable(c, "c20_7") == nil {
			t.Errorf("omni: c20_7 expected to exist in catalog after DDL accepted")
		}
		for _, cname := range expressionCols {
			col := c20getColumn(t, c, "c20_7", cname)
			if col == nil {
				continue
			}
			if col.Default == nil || *col.Default == "" {
				t.Errorf("omni: c20_7.%s Default expected non-nil expression, got nil/empty", cname)
			}
		}
	})

	// --- 20.8 Generated column with DEFAULT clause → error ---------------
	t.Run("20_8_generated_with_default_rejected", func(t *testing.T) {
		cases := []struct {
			name string
			ddl  string
		}{
			{"c20_8a", "CREATE TABLE c20_8a (a INT, b INT AS (a + 1) DEFAULT 0)"},
			{"c20_8b", "CREATE TABLE c20_8b (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL DEFAULT 0)"},
			{"c20_8c", "CREATE TABLE c20_8c (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED DEFAULT (a * 2))"},
		}
		for _, tc := range cases {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			_, mysqlErr := mc.db.ExecContext(mc.ctx, tc.ddl)
			if mysqlErr == nil {
				t.Errorf("oracle: expected parse/usage error for %s, got nil", tc.name)
			} else if !strings.Contains(mysqlErr.Error(), "1064") &&
				!strings.Contains(mysqlErr.Error(), "1221") &&
				!strings.Contains(strings.ToLower(mysqlErr.Error()), "default") {
				t.Errorf("oracle: expected 1064/1221 for %s, got %v", tc.name, mysqlErr)
			}

			omniErrored := c20execExpectError(t, c, tc.ddl+";")
			if !omniErrored {
				t.Errorf("omni: expected rejection of %s DDL; got success", tc.name)
			}
			if c20getTable(c, tc.name) != nil {
				t.Errorf("omni: %s should not be in catalog after rejected DDL", tc.name)
			}
		}
	})
}

// -------------------- C20 local helpers --------------------

// c20oracleColumnDefault reads information_schema.COLUMNS.COLUMN_DEFAULT for a
// column on the testdb database. Returns a sql.NullString so the caller can
// distinguish "no default stored" (Valid=false) from "default is empty string".
func c20oracleColumnDefault(t *testing.T, mc *mysqlContainer, table, column string) sql.NullString {
	t.Helper()
	var colDef sql.NullString
	row := mc.db.QueryRowContext(mc.ctx,
		`SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?`,
		table, column)
	if err := row.Scan(&colDef); err != nil {
		t.Errorf("c20oracleColumnDefault %s.%s: %v", table, column, err)
	}
	return colDef
}

// c20getColumn returns the omni catalog column or nil (reporting t.Error).
func c20getColumn(t *testing.T, c *Catalog, table, column string) *Column {
	t.Helper()
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Errorf("omni: database testdb missing")
		return nil
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		t.Errorf("omni: table %s missing", table)
		return nil
	}
	col := tbl.GetColumn(column)
	if col == nil {
		t.Errorf("omni: column %s.%s missing", table, column)
		return nil
	}
	return col
}

// c20getTable returns the omni catalog table or nil without reporting (used
// to confirm absence after rejected DDL).
func c20getTable(c *Catalog, table string) *Table {
	db := c.GetDatabase("testdb")
	if db == nil {
		return nil
	}
	return db.GetTable(table)
}

// c20execExpectError runs a DDL against omni and returns true if any statement
// in the batch raised an error (parse or exec).
func c20execExpectError(t *testing.T, c *Catalog, ddl string) bool {
	t.Helper()
	results, err := c.Exec(ddl, nil)
	if err != nil {
		return true
	}
	for _, r := range results {
		if r.Error != nil {
			return true
		}
	}
	return false
}
