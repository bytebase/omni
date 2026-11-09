package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C7 covers section C7 (Index defaults). Each subtest asserts that
// live TiDB v8.5.5 and the omni catalog agree on default index behaviour for a
// given DDL input. Forked from mysql/catalog; subtests that assumed MySQL 8.0
// semantics (USING BTREE rendering, FULLTEXT retention, GEOMETRY/SPATIAL) were
// re-baselined to TiDB v8.5.5 behaviour measured against the container.
func TestScenario_C7(t *testing.T) {
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

	// --- 7.4 USING BTREE is not rendered (TiDB) -------------------------
	t.Run("7_4_using_btree_not_rendered", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, KEY (a));
CREATE TABLE t2 (a INT, KEY (a) USING BTREE);`)

		// TiDB has no hash index and writes no USING clause for a BTREE index,
		// so an explicit USING BTREE renders the same as the default — neither
		// table carries USING BTREE in SHOW CREATE (read against
		// pingcap/tidb:v8.5.5 on 2026-11-09).
		ct1 := oracleShow(t, mc, "SHOW CREATE TABLE t1")
		ct2 := oracleShow(t, mc, "SHOW CREATE TABLE t2")
		if strings.Contains(ct1, "USING BTREE") {
			t.Errorf("oracle: t1 unexpectedly contains USING BTREE: %s", ct1)
		}
		if strings.Contains(ct2, "USING BTREE") {
			t.Errorf("oracle: t2 unexpectedly contains USING BTREE: %s", ct2)
		}

		// omni matches: its SHOW CREATE writes no USING clause for either table.
		idx1 := c7findIndex(c, "t1", "a")
		idx2 := c7findIndex(c, "t2", "a")
		if idx1 == nil || idx2 == nil {
			t.Errorf("omni: idx1=%v idx2=%v missing", idx1, idx2)
			return
		}
		omni1 := c.ShowCreateTable("testdb", "t1")
		omni2 := c.ShowCreateTable("testdb", "t2")
		if strings.Contains(omni1, "USING BTREE") {
			t.Errorf("omni: t1 unexpectedly contains USING BTREE: %s", omni1)
		}
		if strings.Contains(omni2, "USING BTREE") {
			t.Errorf("omni: t2 unexpectedly contains USING BTREE: %s", omni2)
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

		// Form 3: a FULLTEXT key is exempt from the prefix rule — TiDB accepts
		// the statement — but TiDB reads it in CREATE TABLE and keeps nothing:
		// SHOW CREATE lists no index and STATISTICS has no row (read against
		// pingcap/tidb:v8.5.5 on 2026-11-09). omni keeps no index either.
		runOnBoth(t, mc, c, `CREATE TABLE t_ft (a TEXT, FULLTEXT KEY (a));`)
		rows := oracleRows(t, mc,
			`SELECT INDEX_NAME FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_ft'`)
		if len(rows) != 0 {
			t.Errorf("oracle: expected no index on t_ft, got %v", rows)
		}
		if idx := c7findIndex(c, "t_ft", "a"); idx != nil {
			t.Errorf("omni: expected no index on t_ft, got %+v", idx)
		}
	})

	// --- 7.8 FULLTEXT (with or without parser) dropped in CREATE TABLE ---
	t.Run("7_8_fulltext_dropped_in_create_table", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a TEXT, FULLTEXT KEY (a));
CREATE TABLE t2 (a TEXT, FULLTEXT KEY (a) WITH PARSER ngram);`)

		// TiDB accepts a FULLTEXT key in CREATE TABLE and keeps nothing, parser
		// clause and all: neither SHOW CREATE carries FULLTEXT or WITH PARSER
		// (read against pingcap/tidb:v8.5.5 on 2026-11-09).
		ct1 := oracleShow(t, mc, "SHOW CREATE TABLE t1")
		ct2 := oracleShow(t, mc, "SHOW CREATE TABLE t2")
		if strings.Contains(ct1, "FULLTEXT") {
			t.Errorf("oracle: t1 unexpectedly carries a FULLTEXT index: %s", ct1)
		}
		if strings.Contains(ct2, "FULLTEXT") || strings.Contains(ct2, "WITH PARSER") {
			t.Errorf("oracle: t2 unexpectedly carries FULLTEXT/WITH PARSER: %s", ct2)
		}

		// omni keeps no FULLTEXT index from CREATE TABLE either.
		if idx := c7findIndex(c, "t1", "a"); idx != nil {
			t.Errorf("omni: t1 expected no index, got %+v", idx)
		}
		if idx := c7findIndex(c, "t2", "a"); idx != nil {
			t.Errorf("omni: t2 expected no index, got %+v", idx)
		}
	})

	// --- 7.9 GEOMETRY type and SPATIAL index unsupported (TiDB) ----------
	t.Run("7_9_geometry_spatial_unsupported", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// TiDB v8.5.5 has no GEOMETRY type and no SPATIAL index: the type is a
		// parse error (1064), not the MySQL spatial-column semantics (1252/3500).
		// omni refuses the type at parse time as well (read against
		// pingcap/tidb:v8.5.5 on 2026-11-09).
		for _, ddl := range []string{
			`CREATE TABLE t_ok (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g))`,
			`CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g))`,
			`CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)`,
		} {
			_, mysqlErr := mc.db.ExecContext(mc.ctx, ddl)
			if mysqlErr == nil {
				t.Errorf("oracle: expected rejection of GEOMETRY/SPATIAL, got nil for %q", ddl)
			} else if !strings.Contains(mysqlErr.Error(), "1064") {
				t.Errorf("oracle: expected 1064 for %q, got %v", ddl, mysqlErr)
			}

			results, err := c.Exec(ddl+";", nil)
			omniErrored := err != nil
			if !omniErrored {
				for _, r := range results {
					if r.Error != nil {
						omniErrored = true
						break
					}
				}
			}
			assertBoolEq(t, "omni: rejects GEOMETRY/SPATIAL for "+ddl, omniErrored, true)
		}
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
