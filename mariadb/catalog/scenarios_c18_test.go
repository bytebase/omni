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
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

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
	// 18.3 ENGINE always rendered
	// -----------------------------------------------------------------
	t.Run("18_3_engine_always_rendered", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT)")

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		if !strings.Contains(mysqlCreate, "ENGINE=InnoDB") {
			t.Errorf("oracle: expected ENGINE=InnoDB; got %q", mysqlCreate)
		}
		if !strings.Contains(omniCreate, "ENGINE=InnoDB") {
			t.Errorf("omni: expected ENGINE=InnoDB; got %q", omniCreate)
		}
	})

	// -----------------------------------------------------------------
	// 18.4 AUTO_INCREMENT elided when counter == 1
	// -----------------------------------------------------------------
	t.Run("18_4_auto_increment_elided_when_one", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY)")

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		if strings.Contains(mysqlCreate, "AUTO_INCREMENT=") {
			t.Errorf("oracle: AUTO_INCREMENT= should be elided; got %q", mysqlCreate)
		}
		if strings.Contains(omniCreate, "AUTO_INCREMENT=") {
			t.Errorf("omni: AUTO_INCREMENT= should be elided; got %q", omniCreate)
		}
	})

	// -----------------------------------------------------------------
	// 18.5 DEFAULT CHARSET always rendered
	// -----------------------------------------------------------------
	t.Run("18_5_default_charset_always_rendered", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE tnocs (x INT)")

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE tnocs")
		omniCreate := c.ShowCreateTable("testdb", "tnocs")

		if !strings.Contains(mysqlCreate, "DEFAULT CHARSET=utf8mb4") {
			t.Errorf("oracle: DEFAULT CHARSET=utf8mb4 missing; got %q", mysqlCreate)
		}
		if !strings.Contains(mysqlCreate, "COLLATE=utf8mb4_0900_ai_ci") {
			t.Errorf("oracle: COLLATE=utf8mb4_0900_ai_ci missing; got %q", mysqlCreate)
		}
		if !strings.Contains(omniCreate, "DEFAULT CHARSET=utf8mb4") {
			t.Errorf("omni: DEFAULT CHARSET=utf8mb4 missing; got %q", omniCreate)
		}
		if !strings.Contains(omniCreate, "COLLATE=utf8mb4_0900_ai_ci") {
			t.Errorf("omni: COLLATE=utf8mb4_0900_ai_ci missing; got %q", omniCreate)
		}
	})

	// -----------------------------------------------------------------
	// 18.6 ROW_FORMAT elided when not explicitly specified
	// -----------------------------------------------------------------
	t.Run("18_6_row_format_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t2 (a INT) ROW_FORMAT=DYNAMIC")

		mysqlT := oracleShow(t, mc, "SHOW CREATE TABLE t")
		mysqlT2 := oracleShow(t, mc, "SHOW CREATE TABLE t2")
		omniT := c.ShowCreateTable("testdb", "t")
		omniT2 := c.ShowCreateTable("testdb", "t2")

		if strings.Contains(mysqlT, "ROW_FORMAT=") {
			t.Errorf("oracle: implicit ROW_FORMAT should be elided; got %q", mysqlT)
		}
		if !strings.Contains(mysqlT2, "ROW_FORMAT=DYNAMIC") {
			t.Errorf("oracle: explicit ROW_FORMAT=DYNAMIC missing; got %q", mysqlT2)
		}
		if strings.Contains(omniT, "ROW_FORMAT=") {
			t.Errorf("omni: implicit ROW_FORMAT should be elided; got %q", omniT)
		}
		if !strings.Contains(omniT2, "ROW_FORMAT=DYNAMIC") {
			t.Errorf("omni: explicit ROW_FORMAT=DYNAMIC missing; got %q", omniT2)
		}
	})

	// -----------------------------------------------------------------
	// 18.7 Table-level COLLATE rendering rules
	// -----------------------------------------------------------------
	t.Run("18_7_table_collate_rendering", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_prim (x INT) CHARACTER SET latin1")
		runOnBoth(t, mc, c, "CREATE TABLE t_nonprim (x INT) CHARACTER SET latin1 COLLATE latin1_bin")
		runOnBoth(t, mc, c, "CREATE TABLE t_0900 (x INT) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci")

		mysqlPrim := oracleShow(t, mc, "SHOW CREATE TABLE t_prim")
		mysqlNonPrim := oracleShow(t, mc, "SHOW CREATE TABLE t_nonprim")
		mysql0900 := oracleShow(t, mc, "SHOW CREATE TABLE t_0900")
		omniPrim := c.ShowCreateTable("testdb", "t_prim")
		omniNonPrim := c.ShowCreateTable("testdb", "t_nonprim")
		omni0900 := c.ShowCreateTable("testdb", "t_0900")

		// t_prim: DEFAULT CHARSET=latin1, NO COLLATE=.
		if !strings.Contains(mysqlPrim, "DEFAULT CHARSET=latin1") {
			t.Errorf("oracle t_prim: missing DEFAULT CHARSET=latin1; got %q", mysqlPrim)
		}
		if strings.Contains(mysqlPrim, "COLLATE=") {
			t.Errorf("oracle t_prim: COLLATE= should be elided; got %q", mysqlPrim)
		}
		if strings.Contains(omniPrim, "COLLATE=") {
			t.Errorf("omni t_prim: COLLATE= should be elided; got %q", omniPrim)
		}
		// t_nonprim: COLLATE=latin1_bin.
		if !strings.Contains(mysqlNonPrim, "COLLATE=latin1_bin") {
			t.Errorf("oracle t_nonprim: missing COLLATE=latin1_bin; got %q", mysqlNonPrim)
		}
		if !strings.Contains(omniNonPrim, "COLLATE=latin1_bin") {
			t.Errorf("omni t_nonprim: missing COLLATE=latin1_bin; got %q", omniNonPrim)
		}
		// t_0900: COLLATE=utf8mb4_0900_ai_ci (special case — always rendered).
		if !strings.Contains(mysql0900, "COLLATE=utf8mb4_0900_ai_ci") {
			t.Errorf("oracle t_0900: missing COLLATE=utf8mb4_0900_ai_ci; got %q", mysql0900)
		}
		if !strings.Contains(omni0900, "COLLATE=utf8mb4_0900_ai_ci") {
			t.Errorf("omni t_0900: missing COLLATE=utf8mb4_0900_ai_ci; got %q", omni0900)
		}
	})

	// -----------------------------------------------------------------
	// 18.8 KEY_BLOCK_SIZE elision
	// -----------------------------------------------------------------
	t.Run("18_8_key_block_size_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_nokbs (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t_kbs (a INT) KEY_BLOCK_SIZE=4")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_nokbs")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_kbs")
		omniNo := c.ShowCreateTable("testdb", "t_nokbs")
		omniYes := c.ShowCreateTable("testdb", "t_kbs")

		if strings.Contains(mysqlNo, "KEY_BLOCK_SIZE=") {
			t.Errorf("oracle t_nokbs: KEY_BLOCK_SIZE should be elided; got %q", mysqlNo)
		}
		if !strings.Contains(mysqlYes, "KEY_BLOCK_SIZE=4") {
			t.Errorf("oracle t_kbs: missing KEY_BLOCK_SIZE=4; got %q", mysqlYes)
		}
		if strings.Contains(omniNo, "KEY_BLOCK_SIZE=") {
			t.Errorf("omni t_nokbs: KEY_BLOCK_SIZE should be elided; got %q", omniNo)
		}
		if !strings.Contains(omniYes, "KEY_BLOCK_SIZE=4") {
			t.Errorf("omni t_kbs: missing KEY_BLOCK_SIZE=4; got %q", omniYes)
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
	// 18.10 STATS_PERSISTENT / STATS_AUTO_RECALC / STATS_SAMPLE_PAGES
	// -----------------------------------------------------------------
	t.Run("18_10_stats_clauses_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_nostats (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t_stats (a INT) STATS_PERSISTENT=1 STATS_AUTO_RECALC=0 STATS_SAMPLE_PAGES=32")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_nostats")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_stats")
		omniNo := c.ShowCreateTable("testdb", "t_nostats")
		omniYes := c.ShowCreateTable("testdb", "t_stats")

		for _, clause := range []string{"STATS_PERSISTENT=", "STATS_AUTO_RECALC=", "STATS_SAMPLE_PAGES="} {
			if strings.Contains(mysqlNo, clause) {
				t.Errorf("oracle t_nostats: %s should be elided; got %q", clause, mysqlNo)
			}
			if strings.Contains(omniNo, clause) {
				t.Errorf("omni t_nostats: %s should be elided; got %q", clause, omniNo)
			}
		}
		for _, want := range []string{"STATS_PERSISTENT=1", "STATS_AUTO_RECALC=0", "STATS_SAMPLE_PAGES=32"} {
			if !strings.Contains(mysqlYes, want) {
				t.Errorf("oracle t_stats: missing %s; got %q", want, mysqlYes)
			}
			if !strings.Contains(omniYes, want) {
				t.Errorf("omni t_stats: missing %s; got %q", want, omniYes)
			}
		}
	})

	// -----------------------------------------------------------------
	// 18.11 MIN_ROWS / MAX_ROWS / AVG_ROW_LENGTH elision
	// -----------------------------------------------------------------
	t.Run("18_11_min_max_avg_rows_elision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t_nominmax (a INT)")
		runOnBoth(t, mc, c, "CREATE TABLE t_minmax (a INT) MIN_ROWS=10 MAX_ROWS=1000 AVG_ROW_LENGTH=256")

		mysqlNo := oracleShow(t, mc, "SHOW CREATE TABLE t_nominmax")
		mysqlYes := oracleShow(t, mc, "SHOW CREATE TABLE t_minmax")
		omniNo := c.ShowCreateTable("testdb", "t_nominmax")
		omniYes := c.ShowCreateTable("testdb", "t_minmax")

		for _, clause := range []string{"MIN_ROWS=", "MAX_ROWS=", "AVG_ROW_LENGTH="} {
			if strings.Contains(mysqlNo, clause) {
				t.Errorf("oracle t_nominmax: %s should be elided; got %q", clause, mysqlNo)
			}
			if strings.Contains(omniNo, clause) {
				t.Errorf("omni t_nominmax: %s should be elided; got %q", clause, omniNo)
			}
		}
		for _, want := range []string{"MIN_ROWS=10", "MAX_ROWS=1000", "AVG_ROW_LENGTH=256"} {
			if !strings.Contains(mysqlYes, want) {
				t.Errorf("oracle t_minmax: missing %s; got %q", want, mysqlYes)
			}
			if !strings.Contains(omniYes, want) {
				t.Errorf("omni t_minmax: missing %s; got %q", want, omniYes)
			}
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

	// -----------------------------------------------------------------
	// 18.14 Per-index COMMENT and KEY_BLOCK_SIZE inside index clauses
	// -----------------------------------------------------------------
	t.Run("18_14_per_index_comment_kbs", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			id INT PRIMARY KEY,
			a  INT,
			b  INT,
			KEY ix_plain (a),
			KEY ix_cmt (b) COMMENT 'hello'
		) KEY_BLOCK_SIZE=4`
		runOnBoth(t, mc, c, ddl)

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		// ix_plain: no COMMENT, no per-index KEY_BLOCK_SIZE.
		plainMy := aLine(mysqlCreate, "ix_plain")
		plainOmni := aLine(omniCreate, "ix_plain")
		if strings.Contains(plainMy, "COMMENT") {
			t.Errorf("oracle ix_plain: should not have COMMENT; got %q", plainMy)
		}
		if strings.Contains(plainMy, "KEY_BLOCK_SIZE") {
			t.Errorf("oracle ix_plain: should not have KEY_BLOCK_SIZE; got %q", plainMy)
		}
		if strings.Contains(plainOmni, "COMMENT") {
			t.Errorf("omni ix_plain: should not have COMMENT; got %q", plainOmni)
		}
		if strings.Contains(plainOmni, "KEY_BLOCK_SIZE") {
			t.Errorf("omni ix_plain: should not have KEY_BLOCK_SIZE; got %q", plainOmni)
		}
		// ix_cmt: COMMENT 'hello' present; no KEY_BLOCK_SIZE.
		cmtMy := aLine(mysqlCreate, "ix_cmt")
		cmtOmni := aLine(omniCreate, "ix_cmt")
		if !strings.Contains(cmtMy, "COMMENT 'hello'") {
			t.Errorf("oracle ix_cmt: missing COMMENT 'hello'; got %q", cmtMy)
		}
		if !strings.Contains(cmtOmni, "COMMENT 'hello'") {
			t.Errorf("omni ix_cmt: missing COMMENT 'hello'; got %q", cmtOmni)
		}
		// Table-level KEY_BLOCK_SIZE=4 present.
		if !strings.Contains(mysqlCreate, "KEY_BLOCK_SIZE=4") {
			t.Errorf("oracle: missing table-level KEY_BLOCK_SIZE=4; got %q", mysqlCreate)
		}
		if !strings.Contains(omniCreate, "KEY_BLOCK_SIZE=4") {
			t.Errorf("omni: missing table-level KEY_BLOCK_SIZE=4; got %q", omniCreate)
		}
	})

	// -----------------------------------------------------------------
	// 18.15 USING BTREE/HASH only when algorithm explicit
	// -----------------------------------------------------------------
	t.Run("18_15_using_algorithm_explicit", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			id INT,
			a  INT,
			b  INT,
			KEY ix_default (a),
			KEY ix_btree (a) USING BTREE,
			KEY ix_hash (b) USING HASH
		) ENGINE=InnoDB`
		runOnBoth(t, mc, c, ddl)

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		// ix_default: no USING clause.
		defaultMy := aLine(mysqlCreate, "ix_default")
		defaultOmni := aLine(omniCreate, "ix_default")
		if strings.Contains(defaultMy, "USING") {
			t.Errorf("oracle ix_default: should not have USING clause; got %q", defaultMy)
		}
		if strings.Contains(defaultOmni, "USING") {
			t.Errorf("omni ix_default: should not have USING clause; got %q", defaultOmni)
		}
		// ix_btree: has USING clause.
		btreeMy := aLine(mysqlCreate, "ix_btree")
		btreeOmni := aLine(omniCreate, "ix_btree")
		if !strings.Contains(btreeMy, "USING") {
			t.Errorf("oracle ix_btree: missing USING clause; got %q", btreeMy)
		}
		if !strings.Contains(btreeOmni, "USING") {
			t.Errorf("omni ix_btree: missing USING clause; got %q", btreeOmni)
		}
		// ix_hash: oracle may or may not render USING (InnoDB rewrites HASH→BTREE).
		// Contract: if oracle renders USING, omni must too.
		hashMy := aLine(mysqlCreate, "ix_hash")
		hashOmni := aLine(omniCreate, "ix_hash")
		if strings.Contains(hashMy, "USING") && !strings.Contains(hashOmni, "USING") {
			t.Errorf("omni ix_hash: oracle has USING but omni does not; oracle=%q omni=%q", hashMy, hashOmni)
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
