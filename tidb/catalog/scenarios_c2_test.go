package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C2 runs the "Type normalization" section of
// SCENARIOS-mysql-implicit-behavior.md (section C2). Each subtest executes a
// DDL against both a real MySQL 8.0 container and the omni catalog, then
// asserts that both agree on the expected normalized type rendering.
//
// Per the worker protocol, failed omni assertions are NOT test infrastructure
// failures — they are tracked as discovered bugs in scenarios_bug_queue/c2.md.
// The test uses t.Error (not t.Fatal) so every scenario reports all its
// diffs in a single run.
func TestScenario_C2(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// Helper: read oracle COLUMN_TYPE for testdb.t.<col>.
	oracleColumnType := func(t *testing.T, col string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			"SELECT COLUMN_TYPE FROM information_schema.COLUMNS "+
				"WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='"+col+"'",
			&s)
		return strings.ToLower(s)
	}
	oracleDataType := func(t *testing.T, col string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			"SELECT DATA_TYPE FROM information_schema.COLUMNS "+
				"WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='"+col+"'",
			&s)
		return strings.ToLower(s)
	}
	oracleCharset := func(t *testing.T, col string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			"SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS "+
				"WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='"+col+"'",
			&s)
		return strings.ToLower(s)
	}
	omniCol := func(t *testing.T, c *Catalog, col string) *Column {
		t.Helper()
		db := c.GetDatabase("testdb")
		if db == nil {
			t.Error("omni: database testdb not found")
			return nil
		}
		tbl := db.GetTable("t")
		if tbl == nil {
			t.Error("omni: table t not found")
			return nil
		}
		cc := tbl.GetColumn(col)
		if cc == nil {
			t.Errorf("omni: column %q not found", col)
			return nil
		}
		return cc
	}

	// 2.1 REAL → DOUBLE
	t.Run("2_1_REAL_to_DOUBLE", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c REAL)`)

		assertStringEq(t, "oracle DATA_TYPE", oracleDataType(t, "c"), "double")
		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "double")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "double")
		}
	})

	// 2.2 BOOL → TINYINT(1)
	t.Run("2_2_BOOL_to_TINYINT1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c BOOL)`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "tinyint(1)")
		assertStringEq(t, "oracle DATA_TYPE", oracleDataType(t, "c"), "tinyint")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "tinyint")
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "tinyint(1)")
		}
	})

	// 2.3 INTEGER → INT
	t.Run("2_3_INTEGER_to_INT", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c INTEGER)`)

		assertStringEq(t, "oracle DATA_TYPE", oracleDataType(t, "c"), "int")
		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "int")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "int")
		}
	})

	// 2.4 BOOLEAN → TINYINT(1)
	t.Run("2_4_BOOLEAN_to_TINYINT1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c BOOLEAN)`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "tinyint(1)")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "tinyint")
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "tinyint(1)")
		}
	})

	// 2.5 INT1/INT2/INT3/INT4/INT8 → TINYINT/SMALLINT/MEDIUMINT/INT/BIGINT
	t.Run("2_5_INTN_aliases", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT1, b INT2, cc INT3, d INT4, e INT8)`)

		cases := []struct {
			name, want string
		}{
			{"a", "tinyint"},
			{"b", "smallint"},
			{"cc", "mediumint"},
			{"d", "int"},
			{"e", "bigint"},
		}
		for _, cc := range cases {
			assertStringEq(t, "oracle DATA_TYPE "+cc.name, oracleDataType(t, cc.name), cc.want)
			if col := omniCol(t, c, cc.name); col != nil {
				assertStringEq(t, "omni DataType "+cc.name, strings.ToLower(col.DataType), cc.want)
			}
		}
	})

	// 2.6 MIDDLEINT → MEDIUMINT
	t.Run("2_6_MIDDLEINT_to_MEDIUMINT", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c MIDDLEINT)`)

		assertStringEq(t, "oracle DATA_TYPE", oracleDataType(t, "c"), "mediumint")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "mediumint")
		}
	})

	// 2.7 INT(11) display width deprecated → stripped from output
	t.Run("2_7_INT11_width_stripped", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c INT(11))`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "int")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "int")
		}
	})

	// 2.8 INT(N) ZEROFILL → preserves display width + implies UNSIGNED
	t.Run("2_8_INT5_ZEROFILL", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c INT(5) ZEROFILL)`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "int(5) unsigned zerofill")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "int(5) unsigned zerofill")
		}
	})

	// 2.9 SERIAL → BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE
	t.Run("2_9_SERIAL", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c SERIAL)`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "bigint unsigned")
		var nullable, extra string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE, EXTRA FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'`,
			&nullable, &extra)
		assertStringEq(t, "oracle IS_NULLABLE", nullable, "NO")
		assertStringEq(t, "oracle EXTRA", strings.ToLower(extra), "auto_increment")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "bigint unsigned")
			assertBoolEq(t, "omni Nullable", col.Nullable, false)
			assertBoolEq(t, "omni AutoIncrement", col.AutoIncrement, true)
		}
		// implicit UNIQUE
		db := c.GetDatabase("testdb")
		if db != nil {
			tbl := db.GetTable("t")
			if tbl != nil {
				foundUnique := false
				for _, idx := range tbl.Indexes {
					if idx.Unique && len(idx.Columns) == 1 && strings.EqualFold(idx.Columns[0].Name, "c") {
						foundUnique = true
						break
					}
				}
				assertBoolEq(t, "omni SERIAL implicit UNIQUE index", foundUnique, true)
			}
		}
	})

	// 2.10 NUMERIC → DECIMAL
	t.Run("2_10_NUMERIC_to_DECIMAL", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c NUMERIC(10,2))`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "decimal(10,2)")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "decimal(10,2)")
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "decimal")
		}
	})

	// 2.11 DEC and FIXED → DECIMAL
	t.Run("2_11_DEC_FIXED_to_DECIMAL", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a DEC(6,2), b FIXED(6,2))`)

		assertStringEq(t, "oracle COLUMN_TYPE a", oracleColumnType(t, "a"), "decimal(6,2)")
		assertStringEq(t, "oracle COLUMN_TYPE b", oracleColumnType(t, "b"), "decimal(6,2)")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni ColumnType a", strings.ToLower(col.ColumnType), "decimal(6,2)")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni ColumnType b", strings.ToLower(col.ColumnType), "decimal(6,2)")
		}
	})

	// 2.12 DOUBLE PRECISION → DOUBLE
	t.Run("2_12_DOUBLE_PRECISION", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c DOUBLE PRECISION)`)

		assertStringEq(t, "oracle DATA_TYPE", oracleDataType(t, "c"), "double")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni DataType", strings.ToLower(col.DataType), "double")
		}
	})

	// 2.13 FLOAT4 → FLOAT, FLOAT8 → DOUBLE
	t.Run("2_13_FLOAT4_FLOAT8", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a FLOAT4, b FLOAT8)`)

		assertStringEq(t, "oracle a", oracleDataType(t, "a"), "float")
		assertStringEq(t, "oracle b", oracleDataType(t, "b"), "double")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a DataType", strings.ToLower(col.DataType), "float")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b DataType", strings.ToLower(col.DataType), "double")
		}
	})

	// 2.14 FLOAT(p) precision split
	t.Run("2_14_FLOAT_precision_split", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a FLOAT(10), b FLOAT(25))`)

		assertStringEq(t, "oracle a", oracleDataType(t, "a"), "float")
		assertStringEq(t, "oracle b", oracleDataType(t, "b"), "double")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a DataType", strings.ToLower(col.DataType), "float")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b DataType", strings.ToLower(col.DataType), "double")
		}
	})

	// 2.15 FLOAT(M,D) deprecated but preserved
	t.Run("2_15_FLOAT_M_D", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c FLOAT(7,4))`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "float(7,4)")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "float(7,4)")
		}
	})

	// 2.16 CHARACTER → CHAR, CHARACTER VARYING → VARCHAR
	t.Run("2_16_CHARACTER_VARYING", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a CHARACTER(10), b CHARACTER VARYING(20))`)

		assertStringEq(t, "oracle a", oracleColumnType(t, "a"), "char(10)")
		assertStringEq(t, "oracle b", oracleColumnType(t, "b"), "varchar(20)")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a ColumnType", strings.ToLower(col.ColumnType), "char(10)")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b ColumnType", strings.ToLower(col.ColumnType), "varchar(20)")
		}
	})

	// 2.17 NATIONAL CHAR / NCHAR → CHAR utf8mb3
	t.Run("2_17_NATIONAL_CHAR", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a NATIONAL CHAR(10), b NCHAR(10))`)

		assertStringEq(t, "oracle a COLUMN_TYPE", oracleColumnType(t, "a"), "char(10)")
		assertStringEq(t, "oracle b COLUMN_TYPE", oracleColumnType(t, "b"), "char(10)")
		assertStringEq(t, "oracle a CHARSET", oracleCharset(t, "a"), "utf8mb3")
		assertStringEq(t, "oracle b CHARSET", oracleCharset(t, "b"), "utf8mb3")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a ColumnType", strings.ToLower(col.ColumnType), "char(10)")
			assertStringEq(t, "omni a Charset", strings.ToLower(col.Charset), "utf8mb3")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b ColumnType", strings.ToLower(col.ColumnType), "char(10)")
			assertStringEq(t, "omni b Charset", strings.ToLower(col.Charset), "utf8mb3")
		}
	})

	// 2.18 NVARCHAR family → VARCHAR utf8mb3
	t.Run("2_18_NVARCHAR_family", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (
  a NVARCHAR(10),
  b NATIONAL VARCHAR(10),
  cc NCHAR VARCHAR(10),
  d NATIONAL CHAR VARYING(10),
  e NCHAR VARYING(10)
)`)

		for _, name := range []string{"a", "b", "cc", "d", "e"} {
			assertStringEq(t, "oracle "+name+" COLUMN_TYPE", oracleColumnType(t, name), "varchar(10)")
			assertStringEq(t, "oracle "+name+" CHARSET", oracleCharset(t, name), "utf8mb3")

			if col := omniCol(t, c, name); col != nil {
				assertStringEq(t, "omni "+name+" ColumnType", strings.ToLower(col.ColumnType), "varchar(10)")
				assertStringEq(t, "omni "+name+" Charset", strings.ToLower(col.Charset), "utf8mb3")
			}
		}
	})

	// 2.19 LONG / LONG VARCHAR → MEDIUMTEXT; LONG VARBINARY → MEDIUMBLOB
	t.Run("2_19_LONG_aliases", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a LONG, b LONG VARCHAR, cc LONG VARBINARY)`)

		assertStringEq(t, "oracle a", oracleDataType(t, "a"), "mediumtext")
		assertStringEq(t, "oracle b", oracleDataType(t, "b"), "mediumtext")
		assertStringEq(t, "oracle cc", oracleDataType(t, "cc"), "mediumblob")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a DataType", strings.ToLower(col.DataType), "mediumtext")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b DataType", strings.ToLower(col.DataType), "mediumtext")
		}
		if col := omniCol(t, c, "cc"); col != nil {
			assertStringEq(t, "omni cc DataType", strings.ToLower(col.DataType), "mediumblob")
		}
	})

	// 2.20 CHAR and BINARY default to length 1
	t.Run("2_20_CHAR_BINARY_default_length", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a CHAR, b BINARY)`)

		assertStringEq(t, "oracle a COLUMN_TYPE", oracleColumnType(t, "a"), "char(1)")
		assertStringEq(t, "oracle b COLUMN_TYPE", oracleColumnType(t, "b"), "binary(1)")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a ColumnType", strings.ToLower(col.ColumnType), "char(1)")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b ColumnType", strings.ToLower(col.ColumnType), "binary(1)")
		}
	})

	// 2.21 VARCHAR without length is a syntax error
	t.Run("2_21_VARCHAR_no_length_is_error", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Oracle: expect failure.
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE t (c VARCHAR)`); err == nil {
			t.Error("oracle: expected syntax error for bare VARCHAR, got nil")
		}

		// omni: expect failure.
		results, err := c.Exec(`CREATE TABLE t (c VARCHAR);`, nil)
		sawErr := err != nil
		if !sawErr {
			for _, r := range results {
				if r.Error != nil {
					sawErr = true
					break
				}
			}
		}
		if !sawErr {
			t.Error("omni: expected parse/exec error for bare VARCHAR, got success")
		}
	})

	// 2.22 TIMESTAMP/DATETIME/TIME default fsp=0, explicit fsp preserved
	t.Run("2_22_temporal_fsp", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (
  a TIMESTAMP NULL,
  b DATETIME,
  cc TIME,
  d TIMESTAMP(6) NULL,
  e DATETIME(6),
  f TIME(3)
)`)

		assertStringEq(t, "oracle a", oracleColumnType(t, "a"), "timestamp")
		assertStringEq(t, "oracle b", oracleColumnType(t, "b"), "datetime")
		assertStringEq(t, "oracle cc", oracleColumnType(t, "cc"), "time")
		assertStringEq(t, "oracle d", oracleColumnType(t, "d"), "timestamp(6)")
		assertStringEq(t, "oracle e", oracleColumnType(t, "e"), "datetime(6)")
		assertStringEq(t, "oracle f", oracleColumnType(t, "f"), "time(3)")

		if col := omniCol(t, c, "a"); col != nil {
			assertStringEq(t, "omni a ColumnType", strings.ToLower(col.ColumnType), "timestamp")
		}
		if col := omniCol(t, c, "b"); col != nil {
			assertStringEq(t, "omni b ColumnType", strings.ToLower(col.ColumnType), "datetime")
		}
		if col := omniCol(t, c, "cc"); col != nil {
			assertStringEq(t, "omni cc ColumnType", strings.ToLower(col.ColumnType), "time")
		}
		if col := omniCol(t, c, "d"); col != nil {
			assertStringEq(t, "omni d ColumnType", strings.ToLower(col.ColumnType), "timestamp(6)")
		}
		if col := omniCol(t, c, "e"); col != nil {
			assertStringEq(t, "omni e ColumnType", strings.ToLower(col.ColumnType), "datetime(6)")
		}
		if col := omniCol(t, c, "f"); col != nil {
			assertStringEq(t, "omni f ColumnType", strings.ToLower(col.ColumnType), "time(3)")
		}
	})

	// 2.23 YEAR(4) deprecated → stored as YEAR
	t.Run("2_23_YEAR_4_bare_YEAR", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c YEAR(4))`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "year")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "year")
		}
	})

	// 2.24 BIT without length defaults to BIT(1)
	t.Run("2_24_BIT_default_1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (c BIT)`)

		assertStringEq(t, "oracle COLUMN_TYPE", oracleColumnType(t, "c"), "bit(1)")

		if col := omniCol(t, c, "c"); col != nil {
			assertStringEq(t, "omni ColumnType", strings.ToLower(col.ColumnType), "bit(1)")
		}
	})

	// 2.25 VARCHAR(65536) in non-strict → TEXT family
	t.Run("2_25_VARCHAR_overflow_to_text", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Set sql_mode='' on this session (oracle only — omni has no session state).
		if _, err := mc.db.ExecContext(mc.ctx, "SET SESSION sql_mode=''"); err != nil {
			t.Errorf("oracle SET sql_mode: %v", err)
		}
		// Restore strict mode after to avoid leaking.
		defer func() {
			_, _ = mc.db.ExecContext(mc.ctx,
				"SET SESSION sql_mode='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
		}()

		// Oracle side.
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE t (c VARCHAR(65536))`); err != nil {
			t.Errorf("oracle CREATE (non-strict) failed: %v", err)
		}
		var gotType string
		oracleScan(t, mc,
			`SELECT DATA_TYPE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'`,
			&gotType)
		gotType = strings.ToLower(gotType)
		if gotType != "mediumtext" && gotType != "text" {
			t.Errorf("oracle: expected mediumtext or text, got %q", gotType)
		}

		// omni side — we document the current behavior rather than asserting
		// a specific outcome, since omni's byte-length → text promotion is
		// a known gap (scenario 2.25 is pending-verify).
		results, err := c.Exec(`CREATE TABLE t (c VARCHAR(65536));`, nil)
		omniErr := err != nil
		if !omniErr {
			for _, r := range results {
				if r.Error != nil {
					omniErr = true
					break
				}
			}
		}
		if !omniErr {
			if col := omniCol(t, c, "c"); col != nil {
				dt := strings.ToLower(col.DataType)
				if dt != "mediumtext" && dt != "text" {
					t.Errorf("omni: expected mediumtext/text, got DataType=%q ColumnType=%q",
						dt, col.ColumnType)
				}
			}
		} else {
			// omni raised an error — acceptable only if we were in strict mode,
			// but since omni has no session state, this is a diff vs oracle.
			t.Errorf("omni: raised error for VARCHAR(65536) but oracle (non-strict) accepted it as %s", gotType)
		}
	})

	// 2.26 TEXT(N) / BLOB(N) → TINY/TEXT/MEDIUM/LONG by byte count
	t.Run("2_26_TEXT_N_promotion", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (
  a TEXT(100),
  b TEXT(1000),
  cc TEXT(70000),
  d TEXT(20000000)
)`)

		// NOTE: SCENARIOS says a TEXT(100) → tinytext, but with default
		// utf8mb4 charset 100 chars = 400 bytes, exceeding tinytext's
		// 255-byte cap, so MySQL promotes to text. Trust the oracle —
		// see scenarios_bug_queue/c2.md for the scenario-doc mismatch.
		cases := []struct {
			name, want string
		}{
			{"a", "text"},
			{"b", "text"},
			{"cc", "mediumtext"},
			{"d", "longtext"},
		}
		for _, cc := range cases {
			assertStringEq(t, "oracle "+cc.name, oracleDataType(t, cc.name), cc.want)
			if col := omniCol(t, c, cc.name); col != nil {
				assertStringEq(t, "omni "+cc.name+" DataType",
					strings.ToLower(col.DataType), cc.want)
			}
		}
	})
}
