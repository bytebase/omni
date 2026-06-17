package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C7 covers section C7 (Index defaults) from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that
// both real MySQL 8.0 and the omni catalog agree on default index
// behaviour for a given DDL input.
//
// Failures in omni assertions are NOT proof failures — they are
// recorded in mysql/catalog/scenarios_bug_queue/c7.md.
func TestScenario_C7(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 7.1 Index algorithm defaults to BTREE --------------------------
	t.Run("7_1_btree_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, KEY (a));`)

		var got string
		oracleScan(t, mc,
			`SELECT INDEX_TYPE FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='a'`,
			&got)
		assertStringEq(t, "oracle index type", got, "BTREE")

		idx := c7findIndex(c, "t", "a")
		if idx == nil {
			t.Errorf("omni: index a missing")
			return
		}
		// Omni either stores "BTREE" explicitly or leaves blank; both round-trip
		// to BTREE. We accept either.
		got2 := strings.ToUpper(idx.IndexType)
		if got2 != "" && got2 != "BTREE" {
			t.Errorf("omni: expected empty or BTREE, got %q", idx.IndexType)
		}
	})

	// --- 7.2 FK creates implicit backing index --------------------------
	t.Run("7_2_fk_implicit_index", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id));`)

		rows := oracleRows(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c'
            ORDER BY INDEX_NAME`)
		var oracleIdxNames []string
		for _, r := range rows {
			oracleIdxNames = append(oracleIdxNames, asString(r[0]))
		}
		want := "c_ibfk_1"
		found := false
		for _, n := range oracleIdxNames {
			if n == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("oracle: expected index %q in %v", want, oracleIdxNames)
		}

		tbl := c.GetDatabase("testdb").GetTable("c")
		if tbl == nil {
			t.Errorf("omni: table c missing")
			return
		}
		// omni: should have an index on column a backing the FK
		hasBackingIdx := false
		for _, idx := range tbl.Indexes {
			if len(idx.Columns) >= 1 && strings.EqualFold(idx.Columns[0].Name, "a") {
				hasBackingIdx = true
				if idx.Name != "c_ibfk_1" {
					t.Errorf("omni: backing index name = %q, want c_ibfk_1", idx.Name)
				}
				break
			}
		}
		assertBoolEq(t, "omni FK backing index exists", hasBackingIdx, true)
	})

	// --- 7.3 USING HASH coerced to BTREE on InnoDB ----------------------
	t.Run("7_3_hash_coerced_to_btree", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, KEY (a) USING HASH) ENGINE=InnoDB;`)

		var got string
		oracleScan(t, mc,
			`SELECT INDEX_TYPE FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='a'`,
			&got)
		assertStringEq(t, "oracle index type after HASH coercion", got, "BTREE")

		idx := c7findIndex(c, "t", "a")
		if idx == nil {
			t.Errorf("omni: index a missing")
			return
		}
		got2 := strings.ToUpper(idx.IndexType)
		if got2 != "" && got2 != "BTREE" {
			t.Errorf("omni: expected empty or BTREE (HASH coercion), got %q", idx.IndexType)
		}
	})

	// --- 7.4 USING BTREE explicit vs implicit rendering -----------------
	t.Run("7_4_using_explicit_vs_implicit", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, KEY (a));
CREATE TABLE t2 (a INT, KEY (a) USING BTREE);`)

		ct1 := oracleShow(t, mc, "SHOW CREATE TABLE t1")
		ct2 := oracleShow(t, mc, "SHOW CREATE TABLE t2")
		// t1 should NOT contain USING BTREE in the KEY line
		if strings.Contains(ct1, "USING BTREE") {
			t.Errorf("oracle: t1 unexpectedly contains USING BTREE: %s", ct1)
		}
		// t2 SHOULD contain USING BTREE
		if !strings.Contains(ct2, "USING BTREE") {
			t.Errorf("oracle: t2 missing USING BTREE: %s", ct2)
		}

		// omni: both tables exist; presence of explicit-flag distinction
		// is an open omni gap. We assert both indexes parse and exist.
		tbl1 := c.GetDatabase("testdb").GetTable("t1")
		tbl2 := c.GetDatabase("testdb").GetTable("t2")
		if tbl1 == nil || tbl2 == nil {
			t.Errorf("omni: t1=%v t2=%v missing", tbl1, tbl2)
			return
		}
		idx1 := c7findIndex(c, "t1", "a")
		idx2 := c7findIndex(c, "t2", "a")
		if idx1 == nil || idx2 == nil {
			t.Errorf("omni: idx1=%v idx2=%v missing", idx1, idx2)
			return
		}
		// omni gap: the catalog has no IndexTypeExplicit field, so it
		// cannot distinguish the two. We assert the distinction here so
		// that the test fails until omni grows the bit.
		t1Explicit := strings.EqualFold(idx1.IndexType, "BTREE")
		t2Explicit := strings.EqualFold(idx2.IndexType, "BTREE")
		if t1Explicit == t2Explicit {
			t.Errorf("omni: cannot distinguish explicit BTREE from default: t1=%q t2=%q",
				idx1.IndexType, idx2.IndexType)
		}
	})

	// --- 7.5 UNIQUE allows multiple NULLs -------------------------------
	t.Run("7_5_unique_multiple_nulls", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, UNIQUE KEY (a));`)

		// Insert three NULLs (oracle only — omni catalog does not execute DML).
		for i := 0; i < 3; i++ {
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO testdb.t VALUES (NULL)"); err != nil {
				t.Errorf("oracle: insert NULL %d failed: %v", i, err)
			}
		}
		var count int
		oracleScan(t, mc, "SELECT COUNT(*) FROM testdb.t", &count)
		assertIntEq(t, "oracle row count after 3 NULL inserts", count, 3)

		// Insert two (1)s — second should fail with duplicate key.
		if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO testdb.t VALUES (1)"); err != nil {
			t.Errorf("oracle: first INSERT 1 failed: %v", err)
		}
		_, err := mc.db.ExecContext(mc.ctx, "INSERT INTO testdb.t VALUES (1)")
		if err == nil {
			t.Errorf("oracle: expected duplicate-key error on second INSERT 1")
		} else if !strings.Contains(err.Error(), "1062") &&
			!strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			t.Errorf("oracle: expected ER_DUP_ENTRY 1062, got %v", err)
		}

		// omni: column a must remain nullable after UNIQUE.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Errorf("omni: column a missing")
			return
		}
		assertBoolEq(t, "omni: UNIQUE keeps column nullable", col.Nullable, true)
	})

	// --- 7.6 VISIBLE default + PK INVISIBLE rejection -------------------
	t.Run("7_6_visible_default_pk_invisible_blocked", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, KEY ix (a));`)

		var visYes string
		oracleScan(t, mc,
			`SELECT IS_VISIBLE FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='ix'`,
			&visYes)
		assertStringEq(t, "oracle IS_VISIBLE default", visYes, "YES")

		idx := c7findIndex(c, "t", "ix")
		if idx == nil {
			t.Errorf("omni: index ix missing")
		} else {
			assertBoolEq(t, "omni: index visible default", idx.Visible, true)
		}

		// PK INVISIBLE must be rejected (oracle 3522). Use fresh table name.
		_, mysqlErr := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t_pkinv (a INT, PRIMARY KEY (a) INVISIBLE)`)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_PK_INDEX_CANT_BE_INVISIBLE, got nil")
		} else if !strings.Contains(mysqlErr.Error(), "3522") &&
			!strings.Contains(strings.ToLower(mysqlErr.Error()), "primary key cannot be invisible") {
			t.Errorf("oracle: expected 3522, got %v", mysqlErr)
		}

		results, err := c.Exec(`CREATE TABLE t_pkinv (a INT, PRIMARY KEY (a) INVISIBLE);`, nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni: rejects PK INVISIBLE", omniErrored, true)
	})

	// --- 7.7 BLOB/TEXT prefix length required ---------------------------
	t.Run("7_7_blob_text_prefix_required", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Form 1: TEXT KEY without length should error (1170).
		_, mysqlErr := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t_noprefix (a TEXT, KEY (a))`)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_BLOB_KEY_WITHOUT_LENGTH, got nil")
		} else if !strings.Contains(mysqlErr.Error(), "1170") &&
			!strings.Contains(strings.ToLower(mysqlErr.Error()), "blob/text column") {
			t.Errorf("oracle: expected 1170, got %v", mysqlErr)
		}
		results, err := c.Exec(`CREATE TABLE t_noprefix (a TEXT, KEY (a));`, nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni: rejects TEXT KEY without prefix", omniErrored, true)

		// Form 2: TEXT KEY with prefix length succeeds; SUB_PART=100.
		runOnBoth(t, mc, c, `CREATE TABLE t_prefix (a TEXT, KEY (a(100)));`)
		var subPart int
		oracleScan(t, mc,
			`SELECT SUB_PART FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_prefix' AND INDEX_NAME='a'`,
			&subPart)
		assertIntEq(t, "oracle SUB_PART", subPart, 100)

		idx := c7findIndex(c, "t_prefix", "a")
		if idx == nil {
			t.Errorf("omni: index a on t_prefix missing")
		} else if len(idx.Columns) != 1 {
			t.Errorf("omni: expected 1 column in idx, got %d", len(idx.Columns))
		} else {
			assertIntEq(t, "omni: prefix length", idx.Columns[0].Length, 100)
		}

		// Form 3: FULLTEXT exempt from prefix requirement.
		runOnBoth(t, mc, c, `CREATE TABLE t_ft (a TEXT, FULLTEXT KEY (a));`)
		var idxType string
		oracleScan(t, mc,
			`SELECT INDEX_TYPE FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_ft' AND INDEX_NAME='a'`,
			&idxType)
		assertStringEq(t, "oracle FULLTEXT index type", idxType, "FULLTEXT")
	})

	// --- 7.8 FULLTEXT WITH PARSER optionality ---------------------------
	t.Run("7_8_fulltext_parser_optional", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a TEXT, FULLTEXT KEY (a));
CREATE TABLE t2 (a TEXT, FULLTEXT KEY (a) WITH PARSER ngram);`)

		ct1 := oracleShow(t, mc, "SHOW CREATE TABLE t1")
		ct2 := oracleShow(t, mc, "SHOW CREATE TABLE t2")
		if strings.Contains(ct1, "WITH PARSER") {
			t.Errorf("oracle: t1 (no parser) unexpectedly contains WITH PARSER: %s", ct1)
		}
		if !strings.Contains(ct2, "WITH PARSER") {
			t.Errorf("oracle: t2 missing WITH PARSER ngram: %s", ct2)
		}

		// omni: catalog Index has no ParserName field — gap. We assert
		// only that both tables and FULLTEXT indexes exist.
		idx1 := c7findIndex(c, "t1", "a")
		idx2 := c7findIndex(c, "t2", "a")
		if idx1 == nil || idx2 == nil {
			t.Errorf("omni: idx1=%v idx2=%v missing", idx1, idx2)
			return
		}
		assertBoolEq(t, "omni t1 fulltext", idx1.Fulltext, true)
		assertBoolEq(t, "omni t2 fulltext", idx2.Fulltext, true)
		// Distinction (parser name) is omni gap — flag for bug queue.
		// We can't read ParserName since the field doesn't exist, so log.
		t.Logf("omni: ParserName field not present on Index; cannot distinguish t1 vs t2 (expected gap)")
	})

	// --- 7.9 SPATIAL requires NOT NULL ----------------------------------
	t.Run("7_9_spatial_not_null", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// OK case: NOT NULL geometry with SRID.
		runOnBoth(t, mc, c, `CREATE TABLE t_ok (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g));`)

		// Error case 1: nullable geometry — ER_SPATIAL_CANT_HAVE_NULL 1252.
		_, mysqlErr := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g))`)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_SPATIAL_CANT_HAVE_NULL, got nil")
		} else if !strings.Contains(mysqlErr.Error(), "1252") &&
			!strings.Contains(strings.ToLower(mysqlErr.Error()), "all parts of a spatial index") {
			t.Errorf("oracle: expected 1252, got %v", mysqlErr)
		}
		results, err := c.Exec(`CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g));`, nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni: rejects nullable SPATIAL", omniErrored, true)

		// Error case 2: SPATIAL with USING BTREE — ER_INDEX_TYPE_NOT_SUPPORTED_FOR_SPATIAL_INDEX 3500.
		_, mysqlErr2 := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)`)
		if mysqlErr2 == nil {
			t.Errorf("oracle: expected rejection of SPATIAL USING BTREE, got nil")
		} else {
			// MySQL 8.0 parser rejects `USING BTREE` on a SPATIAL index at parse
			// time with a plain syntax error (1064) rather than the semantic
			// ER_INDEX_TYPE_NOT_SUPPORTED_FOR_SPATIAL_INDEX (3500). Accept either
			// so the scenario is version-robust.
			msg := strings.ToLower(mysqlErr2.Error())
			if !strings.Contains(mysqlErr2.Error(), "3500") &&
				!strings.Contains(mysqlErr2.Error(), "1064") &&
				!strings.Contains(msg, "not supported") &&
				!strings.Contains(msg, "syntax") {
				t.Errorf("oracle: expected 3500 or 1064/syntax, got %v", mysqlErr2)
			}
		}
		results2, err2 := c.Exec(
			`CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE);`, nil)
		omniErrored2 := err2 != nil
		if !omniErrored2 {
			for _, r := range results2 {
				if r.Error != nil {
					omniErrored2 = true
					break
				}
			}
		}
		assertBoolEq(t, "omni: rejects SPATIAL USING BTREE", omniErrored2, true)
	})

	// --- 7.10 PK + UNIQUE coexistence on same column --------------------
	t.Run("7_10_pk_and_unique_coexist", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY uk (a));`)

		rows := oracleRows(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY INDEX_NAME`)
		var got []string
		for _, r := range rows {
			got = append(got, asString(r[0]))
		}
		// Expect both PRIMARY and uk.
		havePrimary := false
		haveUk := false
		for _, n := range got {
			if n == "PRIMARY" {
				havePrimary = true
			}
			if n == "uk" {
				haveUk = true
			}
		}
		if !havePrimary || !haveUk {
			t.Errorf("oracle: expected both PRIMARY and uk indexes, got %v", got)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniHasPK := false
		omniHasUk := false
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				omniHasPK = true
			}
			if idx.Unique && !idx.Primary && idx.Name == "uk" {
				omniHasUk = true
			}
		}
		assertBoolEq(t, "omni: has PRIMARY index", omniHasPK, true)
		assertBoolEq(t, "omni: has uk unique index", omniHasUk, true)
	})
}

// --- section-local helpers ------------------------------------------------

// c7findIndex returns the named index on the given table in testdb, or nil.
func c7findIndex(c *Catalog, tableName, indexName string) *Index {
	db := c.GetDatabase("testdb")
	if db == nil {
		return nil
	}
	tbl := db.GetTable(tableName)
	if tbl == nil {
		return nil
	}
	for _, idx := range tbl.Indexes {
		if strings.EqualFold(idx.Name, indexName) {
			return idx
		}
	}
	return nil
}
