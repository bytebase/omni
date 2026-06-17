package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C8 covers Section C8 "Table option defaults" from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL against both
// a real MySQL 8.0 container and the omni catalog, then asserts that both
// agree on the effective default for a given table-level option.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c8.md.
func TestScenario_C8(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c8OracleTableScalar runs a single-scalar SELECT against
	// information_schema.TABLES for testdb.<table>, selecting a single
	// string/int column. Uses IFNULL so NULL columns come back as empty
	// string (the cases we care about are TABLE_COLLATION and CREATE_OPTIONS).
	c8OracleStr := func(t *testing.T, col, table string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			"SELECT IFNULL("+col+",'') FROM information_schema.TABLES "+
				"WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='"+table+"'",
			&s)
		return s
	}

	// c8ResetDBWithCharset drops testdb, recreates with a specific
	// charset, USEs it, and returns a fresh omni catalog with the same
	// initial state.
	c8ResetDBWithCharset := func(t *testing.T, charset string) *Catalog {
		t.Helper()
		if _, err := mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS testdb"); err != nil {
			t.Fatalf("oracle DROP DATABASE: %v", err)
		}
		createStmt := "CREATE DATABASE testdb CHARACTER SET " + charset
		if _, err := mc.db.ExecContext(mc.ctx, createStmt); err != nil {
			t.Fatalf("oracle CREATE DATABASE: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, "USE testdb"); err != nil {
			t.Fatalf("oracle USE testdb: %v", err)
		}
		c := New()
		results, err := c.Exec(createStmt+"; USE testdb;", nil)
		if err != nil {
			t.Errorf("omni parse error for %q: %v", createStmt, err)
			return c
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni exec error on stmt %d: %v", r.Index, r.Error)
			}
		}
		return c
	}

	// --- 8.1 Storage engine defaults to InnoDB ---------------------------
	t.Run("8_1_engine_defaults_innodb", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		got := c8OracleStr(t, "ENGINE", "t")
		assertStringEq(t, "oracle ENGINE", strings.ToLower(got), "innodb")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		// Omni may store empty string (meaning default) or "InnoDB".
		omniEngine := strings.ToLower(tbl.Engine)
		if omniEngine != "" && omniEngine != "innodb" {
			t.Errorf("omni Engine: got %q, want \"innodb\" or empty default", tbl.Engine)
		}
	})

	// --- 8.2 ROW_FORMAT defaults to DYNAMIC ------------------------------
	t.Run("8_2_row_format_defaults_dynamic", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		got := c8OracleStr(t, "ROW_FORMAT", "t")
		assertStringEq(t, "oracle ROW_FORMAT", strings.ToLower(got), "dynamic")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		// Omni may store empty string (default) or "DYNAMIC".
		omniRF := strings.ToLower(tbl.RowFormat)
		if omniRF != "" && omniRF != "dynamic" {
			t.Errorf("omni RowFormat: got %q, want \"dynamic\" or empty default", tbl.RowFormat)
		}
	})

	// --- 8.3 AUTO_INCREMENT starts at 1, elided from SHOW CREATE ---------
	t.Run("8_3_auto_increment_starts_at_one", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY)`)

		// Oracle: information_schema.TABLES.AUTO_INCREMENT is the
		// next counter value — reported as NULL (no rows yet) or 1
		// depending on engine state. IFNULL normalises to 0. Either
		// 0 (unset/NULL) or 1 confirms "starts at 1".
		var ai int64
		oracleScan(t, mc,
			`SELECT IFNULL(AUTO_INCREMENT,0) FROM information_schema.TABLES
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'`,
			&ai)
		if ai != 0 && ai != 1 {
			t.Errorf("oracle AUTO_INCREMENT: got %d, want 0 (NULL) or 1", ai)
		}

		// Oracle: SHOW CREATE TABLE elides AUTO_INCREMENT= clause (C18.4).
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "AUTO_INCREMENT=") {
			t.Errorf("oracle SHOW CREATE TABLE should elide AUTO_INCREMENT= clause; got:\n%s", show)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		// Omni's AutoIncrement may be 0 (unset sentinel) or 1.
		if tbl.AutoIncrement != 0 && tbl.AutoIncrement != 1 {
			t.Errorf("omni AutoIncrement: got %d, want 0 or 1", tbl.AutoIncrement)
		}
	})

	// --- 8.4 CHARSET inherits from database default ----------------------
	t.Run("8_4_charset_inherits_from_db", func(t *testing.T) {
		c := c8ResetDBWithCharset(t, "latin1")

		runOnBoth(t, mc, c, `CREATE TABLE t (a VARCHAR(10))`)

		// Oracle: information_schema.TABLES.TABLE_COLLATION should be
		// latin1_swedish_ci (the default collation of latin1).
		got := c8OracleStr(t, "TABLE_COLLATION", "t")
		assertStringEq(t, "oracle TABLE_COLLATION", strings.ToLower(got), "latin1_swedish_ci")

		// Oracle: column a charset is latin1.
		var colCS string
		oracleScan(t, mc,
			`SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`,
			&colCS)
		assertStringEq(t, "oracle column CHARACTER_SET_NAME",
			strings.ToLower(colCS), "latin1")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		assertStringEq(t, "omni table Charset",
			strings.ToLower(tbl.Charset), "latin1")
	})

	// --- 8.5 COLLATE alone derives CHARSET -------------------------------
	t.Run("8_5_collate_alone_derives_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT) COLLATE=latin1_german2_ci`)

		// Oracle: TABLE_COLLATION should be latin1_german2_ci.
		got := c8OracleStr(t, "TABLE_COLLATION", "t")
		assertStringEq(t, "oracle TABLE_COLLATION", strings.ToLower(got), "latin1_german2_ci")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		assertStringEq(t, "omni table Charset",
			strings.ToLower(tbl.Charset), "latin1")
		assertStringEq(t, "omni table Collation",
			strings.ToLower(tbl.Collation), "latin1_german2_ci")
	})

	// --- 8.6 KEY_BLOCK_SIZE defaults to 0 and elided in SHOW -------------
	t.Run("8_6_key_block_size_default_zero", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		// Oracle: CREATE_OPTIONS does not mention key_block_size.
		opts := c8OracleStr(t, "CREATE_OPTIONS", "t")
		if strings.Contains(strings.ToLower(opts), "key_block_size") {
			t.Errorf("oracle CREATE_OPTIONS unexpectedly contains key_block_size: %q", opts)
		}
		// Oracle: SHOW CREATE TABLE omits the clause.
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "KEY_BLOCK_SIZE") {
			t.Errorf("oracle SHOW CREATE TABLE should omit KEY_BLOCK_SIZE; got:\n%s", show)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return
		}
		assertIntEq(t, "omni KeyBlockSize", tbl.KeyBlockSize, 0)
	})

	// --- 8.7 COMPRESSION default (None) ----------------------------------
	t.Run("8_7_compression_default_none", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Setup: plain CREATE without COMPRESSION.
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		// Oracle: SHOW CREATE TABLE omits COMPRESSION=.
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "COMPRESSION=") {
			t.Errorf("oracle SHOW CREATE TABLE should omit COMPRESSION=; got:\n%s", show)
		}
		opts := c8OracleStr(t, "CREATE_OPTIONS", "t")
		if strings.Contains(strings.ToLower(opts), "compress") {
			t.Errorf("oracle CREATE_OPTIONS unexpectedly contains compression: %q", opts)
		}

		// Omni: no Compression field exists. We still verify that a
		// CREATE with an explicit COMPRESSION option parses without
		// error. This exercises the omni parser path even though the
		// value is dropped.
		results, err := c.Exec(`CREATE TABLE t_cmp (a INT) COMPRESSION='NONE';`, nil)
		if err != nil {
			t.Errorf("omni: parse error for COMPRESSION='NONE': %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error for COMPRESSION='NONE': %v", r.Error)
			}
		}
		// Intentional: document that Compression is not modeled.
		// This is a MED-severity omni gap. See scenarios_bug_queue/c8.md.
		_ = c.GetDatabase("testdb").GetTable("t_cmp")
	})

	// --- 8.8 ENCRYPTION depends on server default_table_encryption ------
	t.Run("8_8_encryption_default_off", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		// Oracle: with default_table_encryption=OFF (testcontainer
		// default), SHOW CREATE TABLE omits ENCRYPTION and
		// CREATE_OPTIONS has no encryption entry.
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "ENCRYPTION=") {
			t.Errorf("oracle SHOW CREATE TABLE should omit ENCRYPTION=; got:\n%s", show)
		}
		opts := c8OracleStr(t, "CREATE_OPTIONS", "t")
		if strings.Contains(strings.ToLower(opts), "encrypt") {
			t.Errorf("oracle CREATE_OPTIONS unexpectedly contains encryption: %q", opts)
		}

		// Omni gap: no Encryption field. Verify parser accepts the
		// option without crashing.
		results, err := c.Exec(`CREATE TABLE t_enc (a INT) ENCRYPTION='N';`, nil)
		if err != nil {
			t.Errorf("omni: parse error for ENCRYPTION='N': %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error for ENCRYPTION='N': %v", r.Error)
			}
		}
	})

	// --- 8.9 STATS_PERSISTENT defaults to DEFAULT -----------------------
	t.Run("8_9_stats_persistent_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		// Oracle: CREATE_OPTIONS does not mention stats_persistent;
		// SHOW CREATE TABLE omits the clause.
		opts := c8OracleStr(t, "CREATE_OPTIONS", "t")
		if strings.Contains(strings.ToLower(opts), "stats_persistent") {
			t.Errorf("oracle CREATE_OPTIONS unexpectedly contains stats_persistent: %q", opts)
		}
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "STATS_PERSISTENT") {
			t.Errorf("oracle SHOW CREATE TABLE should omit STATS_PERSISTENT; got:\n%s", show)
		}

		// Omni gap: no StatsPersistent field. Verify parser accepts.
		results, err := c.Exec(`CREATE TABLE t_sp (a INT) STATS_PERSISTENT=DEFAULT;`, nil)
		if err != nil {
			t.Errorf("omni: parse error for STATS_PERSISTENT=DEFAULT: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error for STATS_PERSISTENT=DEFAULT: %v", r.Error)
			}
		}
	})

	// --- 8.10 TABLESPACE defaults to innodb_file_per_table --------------
	t.Run("8_10_tablespace_default_file_per_table", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT)`)

		// Oracle: SHOW CREATE TABLE omits TABLESPACE= clause by default.
		show := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(strings.ToUpper(show), "TABLESPACE=") {
			t.Errorf("oracle SHOW CREATE TABLE should omit TABLESPACE=; got:\n%s", show)
		}
		// Oracle: information_schema.INNODB_TABLES has a row for the
		// table, and its SPACE column is non-zero (each
		// file_per_table tablespace has its own id).
		var space int64
		oracleScan(t, mc,
			`SELECT IFNULL(SPACE,0) FROM information_schema.INNODB_TABLES
            WHERE NAME='testdb/t'`,
			&space)
		if space == 0 {
			t.Errorf("oracle INNODB_TABLES.SPACE for testdb/t: got 0, want non-zero (file_per_table)")
		}

		// Omni gap: no Tablespace field. Verify parser accepts an
		// explicit TABLESPACE clause without crashing.
		results, err := c.Exec(`CREATE TABLE t_ts (a INT) TABLESPACE=innodb_file_per_table;`, nil)
		if err != nil {
			t.Errorf("omni: parse error for TABLESPACE=innodb_file_per_table: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: exec error for TABLESPACE=innodb_file_per_table: %v", r.Error)
			}
		}
	})
}
