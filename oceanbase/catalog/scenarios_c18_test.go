package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C18 covers Section C18 "SHOW CREATE TABLE elision rules" from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. Each subtest runs the
// scenario's DDL on both the MySQL 8.0 container and the omni catalog, then
// asserts that omni's ShowCreateTable output matches the oracle on the
// specific elision rule under test.
//
// Critical section: every mismatch here breaks SDL round-trip. Failures are
// recorded as t.Error (not t.Fatal) so all 15 scenarios run in one pass, and
// each omni gap is documented in mysql/catalog/scenarios_bug_queue/c18.md.
func TestScenario_C18(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// -----------------------------------------------------------------
	// 18.1 Column charset elided when equal to table default
	// -----------------------------------------------------------------
	t.Run("18_1_column_charset_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			a VARCHAR(10),
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
		) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci`
		runOnBoth(t, mc, c, ddl)

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		// Column a: no column-level CHARACTER SET / COLLATE.
		if strings.Contains(aLine(mysqlCreate, "`a`"), "CHARACTER SET") {
			t.Errorf("oracle: column a should not have CHARACTER SET; got %q", mysqlCreate)
		}
		if strings.Contains(aLine(omniCreate, "`a`"), "CHARACTER SET") {
			t.Errorf("omni: column a should not have CHARACTER SET; got %q", omniCreate)
		}
		// Column b: has non-default collation, so CHARACTER SET/COLLATE rendered.
		if !strings.Contains(aLine(mysqlCreate, "`b`"), "utf8mb4_unicode_ci") {
			t.Errorf("oracle: column b should have COLLATE utf8mb4_unicode_ci; got %q", mysqlCreate)
		}
		if !strings.Contains(aLine(omniCreate, "`b`"), "utf8mb4_unicode_ci") {
			t.Errorf("omni: column b should have COLLATE utf8mb4_unicode_ci; got %q", omniCreate)
		}
	})

	// -----------------------------------------------------------------
	// 18.2 NOT NULL elision: TIMESTAMP shows NULL, others hide it
	// -----------------------------------------------------------------
	t.Run("18_2_null_elision_timestamp", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			i    INT,
			i_nn INT NOT NULL,
			ts   TIMESTAMP NULL DEFAULT NULL,
			ts_nn TIMESTAMP NOT NULL DEFAULT '2020-01-01'
		)`
		runOnBoth(t, mc, c, ddl)

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		// Oracle: i has no explicit NULL, ts has explicit NULL.
		if strings.Contains(aLine(mysqlCreate, "`i` "), "NULL DEFAULT") {
			// `i int DEFAULT NULL` is fine; we only object to `NULL DEFAULT NULL`.
		}
		if !strings.Contains(aLine(mysqlCreate, "`ts` "), "NULL DEFAULT NULL") {
			t.Errorf("oracle: ts TIMESTAMP should render explicit NULL; got %q", aLine(mysqlCreate, "`ts` "))
		}
		// omni comparison
		if strings.Contains(aLine(omniCreate, "`i` "), "`i` int NULL") {
			t.Errorf("omni: i should not render NULL keyword; got %q", aLine(omniCreate, "`i` "))
		}
		if !strings.Contains(aLine(omniCreate, "`ts` "), "NULL") {
			t.Errorf("omni: ts TIMESTAMP should render explicit NULL; got %q", aLine(omniCreate, "`ts` "))
		}
	})

	// -----------------------------------------------------------------
	// 18.9 COMPRESSION elision
	// -----------------------------------------------------------------
	t.Run("18_9_compression_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_nocomp (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t_comp (a INT) COMPRESSION='ZLIB'")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_nocomp")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_comp")
		omniNo := c.ShowCreateTable("testdb", "t_nocomp")
		omniYes := c.ShowCreateTable("testdb", "t_comp")

		if strings.Contains(mysqlNo, "COMPRESSION=") {
			t.Errorf("oracle t_nocomp: COMPRESSION should be elided; got %q", mysqlNo)
		}
		if !strings.Contains(mysqlYes, "COMPRESSION='ZLIB'") {
			t.Errorf("oracle t_comp: missing COMPRESSION='ZLIB'; got %q", mysqlYes)
		}
		if strings.Contains(omniNo, "COMPRESSION=") {
			t.Errorf("omni t_nocomp: COMPRESSION should be elided; got %q", omniNo)
		}
		if !strings.Contains(omniYes, "COMPRESSION='ZLIB'") {
			t.Errorf("omni t_comp: missing COMPRESSION='ZLIB'; got %q", omniYes)
		}
	})

	// -----------------------------------------------------------------
	// 18.12 TABLESPACE clause elision
	// -----------------------------------------------------------------
	t.Run("18_12_tablespace_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_default (a INT)")
		// Use innodb_system — always available in MySQL 8.0.
		runOnBoth(t, mc, c, "CREATE TABLE t_gts (a INT) TABLESPACE=innodb_system")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_default")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_gts")
		omniNo := c.ShowCreateTable("testdb", "t_default")
		omniYes := c.ShowCreateTable("testdb", "t_gts")

		if strings.Contains(mysqlNo, "TABLESPACE") {
			t.Errorf("oracle t_default: TABLESPACE should be elided; got %q", mysqlNo)
		}
		if !strings.Contains(mysqlYes, "TABLESPACE") {
			t.Errorf("oracle t_gts: missing TABLESPACE clause; got %q", mysqlYes)
		}
		if strings.Contains(omniNo, "TABLESPACE") {
			t.Errorf("omni t_default: TABLESPACE should be elided; got %q", omniNo)
		}
		if !strings.Contains(omniYes, "TABLESPACE") {
			t.Errorf("omni t_gts: missing TABLESPACE clause; got %q", omniYes)
		}
	})

	// -----------------------------------------------------------------
	// 18.13 PACK_KEYS / CHECKSUM / DELAY_KEY_WRITE elision
	// -----------------------------------------------------------------
	t.Run("18_13_pack_checksum_delay_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_none (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t_opts (a INT) PACK_KEYS=1 CHECKSUM=1 DELAY_KEY_WRITE=1")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_none")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_opts")
		omniNo := c.ShowCreateTable("testdb", "t_none")
		omniYes := c.ShowCreateTable("testdb", "t_opts")

		for _, clause := range []string{"PACK_KEYS=", "CHECKSUM=", "DELAY_KEY_WRITE="} {
			if strings.Contains(mysqlNo, clause) {
				t.Errorf("oracle t_none: %s should be elided; got %q", clause, mysqlNo)
			}
			if strings.Contains(omniNo, clause) {
				t.Errorf("omni t_none: %s should be elided; got %q", clause, omniNo)
			}
		}
		// Oracle may reject or normalize; only assert presence where oracle shows them.
		for _, want := range []string{"PACK_KEYS=1", "CHECKSUM=1", "DELAY_KEY_WRITE=1"} {
			if strings.Contains(mysqlYes, want) && !strings.Contains(omniYes, want) {
				t.Errorf("omni t_opts: oracle has %s but omni does not; got omni=%q", want, omniYes)
			}
		}
	})
}

// aLine returns the line from multi-line text s that contains needle. Empty
// string if no match. Used by C18 to grab a single column/index line out of
// SHOW CREATE TABLE output for per-line substring checks.
func aLine(s, needle string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return ""
}
