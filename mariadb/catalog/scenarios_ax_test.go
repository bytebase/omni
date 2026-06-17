package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// TestScenario_AX covers section AX (ALTER TABLE sub-command implicit
// behaviors) from SCENARIOS-mysql-implicit-behavior.md. Each subtest
// asserts that real MySQL 8.0 and the omni catalog agree on the
// post-ALTER state (column order, index list, constraint presence,
// error behavior) for a given ALTER TABLE sequence.
//
// Failures in omni assertions are NOT proof failures — they are
// recorded in mysql/catalog/scenarios_bug_queue/ax.md.
func TestScenario_AX(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- AX.1 ADD COLUMN append / FIRST / AFTER --------------------------
	t.Run("AX_1_AddColumn_append_first_after", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, c INT);
ALTER TABLE t ADD COLUMN d INT;
ALTER TABLE t ADD COLUMN e INT FIRST;
ALTER TABLE t ADD COLUMN f INT AFTER a;`)

		// Oracle: columns ordered.
		want := []string{"e", "a", "f", "b", "c", "d"}
		oracle := axOracleColumnOrder(t, mc, "t")
		assertStringEq(t, "oracle column order", strings.Join(oracle, ","), strings.Join(want, ","))

		// omni
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omni := axOmniColumnOrder(tbl)
		assertStringEq(t, "omni column order", strings.Join(omni, ","), strings.Join(want, ","))
	})

	// --- AX.2 DROP COLUMN cascades removal of indexes containing column --
	t.Run("AX_2_DropColumn_cascades_indexes", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, c INT, INDEX idx_a(a), INDEX idx_ab(a, b), INDEX idx_bc(b, c));
ALTER TABLE t DROP COLUMN a;`)

		// Oracle: scenario doc claims idx_ab should be removed entirely,
		// but MySQL 8.0 actually strips only the dropped column and keeps
		// idx_ab as an index on the surviving column(s). Verify dual-agreement
		// against observed MySQL behavior rather than the doc's stale claim.
		oracleIdx := axOracleIndexNames(t, mc, "t")
		want := []string{"idx_ab", "idx_bc"}
		assertStringEq(t, "oracle surviving indexes", strings.Join(oracleIdx, ","), strings.Join(want, ","))

		// omni
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omni := axOmniIndexNames(tbl)
		assertStringEq(t, "omni surviving indexes", strings.Join(omni, ","), strings.Join(want, ","))
	})

	// --- AX.3 DROP COLUMN rejects last-column removal --------------------
	t.Run("AX_3_DropColumn_rejects_last_column", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);`)

		// Oracle rejects with ER_CANT_REMOVE_ALL_FIELDS (1090).
		_, oracleErr := mc.db.ExecContext(mc.ctx, `ALTER TABLE t DROP COLUMN a`)
		if oracleErr == nil {
			t.Errorf("oracle: expected error dropping last column, got nil")
		}

		// omni should also reject.
		results, perr := c.Exec(`ALTER TABLE t DROP COLUMN a`, nil)
		omniRejected := perr != nil
		if !omniRejected {
			for _, r := range results {
				if r.Error != nil {
					omniRejected = true
					break
				}
			}
		}
		if !omniRejected {
			t.Errorf("omni: expected error dropping last column, got nil")
		}
	})

	// --- AX.4 DROP COLUMN rejects when referenced by CHECK or GENERATED --
	t.Run("AX_4_DropColumn_rejects_check_or_generated_ref", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, b INT, CHECK (a > 0));
CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a + 1));`)

		// Oracle: check behavior of both DROPs against this container's
		// MySQL build. We then assert omni matches whatever oracle does.
		_, oracleErrCheck := mc.db.ExecContext(mc.ctx, `ALTER TABLE t1 DROP COLUMN a`)
		_, oracleErrGen := mc.db.ExecContext(mc.ctx, `ALTER TABLE t2 DROP COLUMN a`)

		// omni
		r1, _ := c.Exec(`ALTER TABLE t1 DROP COLUMN a`, nil)
		omniBlockedCheck := axExecHasError(r1)
		if (oracleErrCheck != nil) != omniBlockedCheck {
			t.Errorf("omni vs oracle CHECK-ref DROP divergence: oracleErr=%v omniBlocked=%v", oracleErrCheck, omniBlockedCheck)
		}
		r2, _ := c.Exec(`ALTER TABLE t2 DROP COLUMN a`, nil)
		omniBlockedGen := axExecHasError(r2)
		if (oracleErrGen != nil) != omniBlockedGen {
			t.Errorf("omni vs oracle GENERATED-ref DROP divergence: oracleErr=%v omniBlocked=%v", oracleErrGen, omniBlockedGen)
		}
	})

	// --- AX.5 MODIFY COLUMN rewrites spec; attrs NOT inherited ----------
	t.Run("AX_5_ModifyColumn_rewrites_spec", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT PRIMARY KEY, b INT NOT NULL DEFAULT 5 COMMENT 'x');
ALTER TABLE t MODIFY COLUMN b BIGINT;`)

		var isNullable, colComment string
		var colDefault *string
		oracleScan(t, mc, `SELECT IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`,
			&isNullable, &colDefault, &colComment)
		if isNullable != "YES" {
			t.Errorf("oracle b nullable: got %q, want YES", isNullable)
		}
		if colDefault != nil {
			t.Errorf("oracle b default: got %v, want nil", *colDefault)
		}
		if colComment != "" {
			t.Errorf("oracle b comment: got %q, want empty", colComment)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Errorf("omni: column b missing")
			return
		}
		if !col.Nullable {
			t.Errorf("omni b nullable: got false, want true (MODIFY should not inherit NOT NULL)")
		}
		if col.Default != nil {
			t.Errorf("omni b default: got %q, want nil (MODIFY should not inherit DEFAULT)", *col.Default)
		}
		if col.Comment != "" {
			t.Errorf("omni b comment: got %q, want empty (MODIFY should not inherit COMMENT)", col.Comment)
		}
		if !strings.Contains(strings.ToLower(col.ColumnType), "bigint") && col.DataType != "bigint" {
			t.Errorf("omni b type: got %q/%q, want bigint", col.DataType, col.ColumnType)
		}
	})

	// --- AX.6 CHANGE COLUMN atomic rename+retype ------------------------
	t.Run("AX_6_ChangeColumn_rename_retype", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);
ALTER TABLE t CHANGE COLUMN a b BIGINT NOT NULL;`)

		var colName, dataType, isNullable string
		oracleScan(t, mc, `SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'`,
			&colName, &dataType, &isNullable)
		if colName != "b" || dataType != "bigint" || isNullable != "NO" {
			t.Errorf("oracle: got (%q,%q,%q), want (b,bigint,NO)", colName, dataType, isNullable)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		if tbl.GetColumn("a") != nil {
			t.Errorf("omni: column a still present after CHANGE")
		}
		col := tbl.GetColumn("b")
		if col == nil {
			t.Errorf("omni: column b missing after CHANGE")
			return
		}
		if col.Nullable {
			t.Errorf("omni b nullable: got true, want false")
		}
		if !strings.Contains(strings.ToLower(col.ColumnType+" "+col.DataType), "bigint") {
			t.Errorf("omni b type: got %q/%q, want bigint", col.DataType, col.ColumnType)
		}
	})

	// --- AX.7 ADD INDEX auto-name in ALTER context -----------------------
	t.Run("AX_7_AddIndex_auto_name_alter", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, INDEX (a));
ALTER TABLE t ADD INDEX (b);
ALTER TABLE t ADD INDEX (a, b);`)

		want := []string{"a", "a_2", "b"}
		oracle := axOracleIndexNames(t, mc, "t")
		assertStringEq(t, "oracle index names", strings.Join(oracle, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omni := axOmniIndexNames(tbl)
		assertStringEq(t, "omni index names", strings.Join(omni, ","), strings.Join(want, ","))
	})

	// --- AX.8 ADD UNIQUE / KEY / FULLTEXT implicit index types ----------
	t.Run("AX_8_AddUnique_Key_Fulltext_types", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b VARCHAR(100)) ENGINE=InnoDB;
ALTER TABLE t ADD UNIQUE (a);
ALTER TABLE t ADD KEY (b);
ALTER TABLE t ADD FULLTEXT (b);`)

		// Oracle check: examine STATISTICS for counts.
		rows := oracleRows(t, mc, `SELECT INDEX_NAME, NON_UNIQUE, INDEX_TYPE FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY INDEX_NAME`)
		if len(rows) < 3 {
			t.Errorf("oracle: expected >=3 index rows, got %d", len(rows))
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		var haveUnique, haveFulltext, havePlain bool
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				continue
			}
			if idx.Unique {
				haveUnique = true
			}
			if idx.Fulltext {
				haveFulltext = true
			}
			if !idx.Unique && !idx.Fulltext && !idx.Spatial && !idx.Primary {
				havePlain = true
			}
		}
		if !haveUnique {
			t.Errorf("omni: expected a UNIQUE index after ADD UNIQUE")
		}
		if !havePlain {
			t.Errorf("omni: expected a plain KEY index after ADD KEY")
		}
		if !haveFulltext {
			t.Errorf("omni: expected a FULLTEXT index after ADD FULLTEXT")
		}
	})

	// --- AX.9 FK column-level shorthand silent-ignored in CREATE --------
	t.Run("AX_9_FK_column_shorthand_silent_ignore", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE parent (id INT PRIMARY KEY);
CREATE TABLE t1 (a INT);
ALTER TABLE t1 ADD FOREIGN KEY (a) REFERENCES parent(id);
CREATE TABLE t2 (a INT REFERENCES parent(id));`)

		// Oracle: t1 has FK, t2 has none.
		oracleT1 := oracleFKNames(t, mc, "t1")
		if len(oracleT1) != 1 {
			t.Errorf("oracle t1 FK count: got %d, want 1", len(oracleT1))
		}
		oracleT2 := oracleFKNames(t, mc, "t2")
		if len(oracleT2) != 0 {
			t.Errorf("oracle t2 FK count: got %d, want 0 (column-level REFERENCES should be ignored)", len(oracleT2))
		}

		tbl1 := c.GetDatabase("testdb").GetTable("t1")
		tbl2 := c.GetDatabase("testdb").GetTable("t2")
		if tbl1 == nil || tbl2 == nil {
			t.Errorf("omni: t1 or t2 missing")
			return
		}
		if n := len(omniFKNames(tbl1)); n != 1 {
			t.Errorf("omni t1 FK count: got %d, want 1", n)
		}
		if n := len(omniFKNames(tbl2)); n != 0 {
			t.Errorf("omni t2 FK count: got %d, want 0 (column-level REFERENCES should be parsed-but-ignored)", n)
		}
	})

	// --- AX.10 RENAME COLUMN preserves attributes ------------------------
	t.Run("AX_10_RenameColumn_preserves_attrs", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT NOT NULL DEFAULT 5 COMMENT 'hi', INDEX (a));
ALTER TABLE t RENAME COLUMN a TO aa;`)

		var colName, isNullable, colComment string
		var colDefault *string
		oracleScan(t, mc, `SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT
            FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'`,
			&colName, &isNullable, &colDefault, &colComment)
		if colName != "aa" || isNullable != "NO" || colComment != "hi" {
			t.Errorf("oracle: got (%q,%q,%q), want (aa,NO,hi)", colName, isNullable, colComment)
		}
		if colDefault == nil || *colDefault != "5" {
			t.Errorf("oracle default: got %v, want 5", colDefault)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		if tbl.GetColumn("a") != nil {
			t.Errorf("omni: column a still present after RENAME COLUMN")
		}
		col := tbl.GetColumn("aa")
		if col == nil {
			t.Errorf("omni: column aa missing")
			return
		}
		if col.Nullable {
			t.Errorf("omni aa nullable: got true, want false (attrs preserved)")
		}
		if col.Default == nil || *col.Default != "5" {
			t.Errorf("omni aa default: got %v, want 5", col.Default)
		}
		if col.Comment != "hi" {
			t.Errorf("omni aa comment: got %q, want hi", col.Comment)
		}
	})

	// --- AX.11 RENAME INDEX ---------------------------------------------
	t.Run("AX_11_RenameIndex", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, INDEX old_name (a, b));
ALTER TABLE t RENAME INDEX old_name TO new_name;`)

		oracle := axOracleIndexNames(t, mc, "t")
		want := []string{"new_name"}
		assertStringEq(t, "oracle index names", strings.Join(oracle, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omni := axOmniIndexNames(tbl)
		assertStringEq(t, "omni index names", strings.Join(omni, ","), strings.Join(want, ","))
	})

	// --- AX.12 RENAME TO (table rename via ALTER) -----------------------
	t.Run("AX_12_RenameTable_via_ALTER", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);
ALTER TABLE t RENAME TO t2;`)

		// Oracle: t gone, t2 present.
		var cnt int64
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.TABLES
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'`, &cnt)
		if cnt != 0 {
			t.Errorf("oracle: expected t to be gone, count=%d", cnt)
		}
		oracleScan(t, mc, `SELECT COUNT(*) FROM information_schema.TABLES
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2'`, &cnt)
		if cnt != 1 {
			t.Errorf("oracle: expected t2 to exist, count=%d", cnt)
		}

		db := c.GetDatabase("testdb")
		if db.GetTable("t") != nil {
			t.Errorf("omni: table t still present after RENAME TO t2")
		}
		if db.GetTable("t2") == nil {
			t.Errorf("omni: table t2 missing after RENAME TO")
		}
	})

	// --- AX.13 ALTER COLUMN SET/DROP DEFAULT, SET INVISIBLE/VISIBLE -----
	t.Run("AX_13_AlterColumn_default_visibility", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT NOT NULL, b INT);
ALTER TABLE t ALTER COLUMN a SET DEFAULT 5;`)

		var colDefault *string
		oracleScan(t, mc, `SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`, &colDefault)
		if colDefault == nil || *colDefault != "5" {
			t.Errorf("oracle a default after SET: got %v, want 5", colDefault)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		colA := tbl.GetColumn("a")
		if colA == nil {
			t.Errorf("omni: column a missing")
			return
		}
		if colA.Default == nil || *colA.Default != "5" {
			t.Errorf("omni a default after SET: got %v, want 5", colA.Default)
		}

		runOnBoth(t, mc, c, `ALTER TABLE t ALTER COLUMN a DROP DEFAULT;`)
		oracleScan(t, mc, `SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`, &colDefault)
		if colDefault != nil {
			t.Errorf("oracle a default after DROP: got %q, want nil", *colDefault)
		}
		colA = c.GetDatabase("testdb").GetTable("t").GetColumn("a")
		if colA.Default != nil {
			t.Errorf("omni a default after DROP: got %q, want nil", *colA.Default)
		}

		runOnBoth(t, mc, c, `ALTER TABLE t ALTER COLUMN b SET INVISIBLE;`)
		var extra string
		oracleScan(t, mc, `SELECT EXTRA FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`, &extra)
		if !strings.Contains(strings.ToUpper(extra), "INVISIBLE") {
			t.Errorf("oracle b EXTRA after SET INVISIBLE: got %q, want INVISIBLE", extra)
		}
		colB := c.GetDatabase("testdb").GetTable("t").GetColumn("b")
		if !colB.Invisible {
			t.Errorf("omni b Invisible after SET INVISIBLE: got false, want true")
		}

		runOnBoth(t, mc, c, `ALTER TABLE t ALTER COLUMN b SET VISIBLE;`)
		oracleScan(t, mc, `SELECT EXTRA FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'`, &extra)
		if strings.Contains(strings.ToUpper(extra), "INVISIBLE") {
			t.Errorf("oracle b EXTRA after SET VISIBLE: got %q, want no INVISIBLE", extra)
		}
		colB = c.GetDatabase("testdb").GetTable("t").GetColumn("b")
		if colB.Invisible {
			t.Errorf("omni b Invisible after SET VISIBLE: got true, want false")
		}
	})

	// --- AX.14 ALTER INDEX VISIBLE / INVISIBLE ---------------------------
	t.Run("AX_14_AlterIndex_visibility", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, INDEX idx_a (a));
ALTER TABLE t ALTER INDEX idx_a INVISIBLE;`)

		var isVisible string
		oracleScan(t, mc, `SELECT IS_VISIBLE FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='idx_a' LIMIT 1`, &isVisible)
		if isVisible != "NO" {
			t.Errorf("oracle idx_a IS_VISIBLE: got %q, want NO", isVisible)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		var omniIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "idx_a" {
				omniIdx = idx
				break
			}
		}
		if omniIdx == nil {
			t.Errorf("omni: index idx_a missing")
			return
		}
		if omniIdx.Visible {
			t.Errorf("omni idx_a Visible: got true, want false after INVISIBLE")
		}
	})

	// --- AX.15 Multi-sub-command ALTER — drop/add/rename composed -------
	t.Run("AX_15_MultiSubCommand_compose", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Use a variant the scenario doc's single-pass claim actually holds
		// for: MySQL evaluates RENAME COLUMN after ADD COLUMN's AFTER, so
		// `ADD COLUMN c INT AFTER a` + `RENAME a TO aa` in the same ALTER
		// raises "Unknown column a" (observed in MySQL 8.0). We split the
		// scenario into a simpler composition that both engines accept and
		// verify positional + index resolution after DROP + ADD + RENAME.
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE t (a INT, b INT)`); err != nil {
			t.Fatalf("oracle create t: %v", err)
		}
		if _, err := c.Exec(`CREATE TABLE t (a INT, b INT);`, nil); err != nil {
			t.Fatalf("omni create t: %v", err)
		}

		ddl := `ALTER TABLE t
  ADD COLUMN c INT,
  DROP COLUMN b,
  ADD INDEX (c),
  RENAME COLUMN a TO aa`
		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		if oracleErr != nil {
			t.Errorf("oracle ALTER failed: %v", oracleErr)
		}

		// omni
		results, perr := c.Exec(ddl, nil)
		if perr != nil {
			t.Errorf("omni parse error: %v", perr)
		}
		if axExecHasError(results) {
			t.Errorf("omni ALTER exec error: %v", results)
		}

		want := []string{"aa", "c"}
		oracle := axOracleColumnOrder(t, mc, "t")
		assertStringEq(t, "oracle column order", strings.Join(oracle, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omni := axOmniColumnOrder(tbl)
		assertStringEq(t, "omni column order", strings.Join(omni, ","), strings.Join(want, ","))

		// And the (c) index exists.
		oracleIdx := axOracleIndexNames(t, mc, "t")
		foundOracle := false
		for _, n := range oracleIdx {
			if n == "c" {
				foundOracle = true
				break
			}
		}
		if !foundOracle {
			t.Errorf("oracle: expected index named c, got %v", oracleIdx)
		}
		foundOmni := false
		for _, idx := range tbl.Indexes {
			if idx.Name == "c" {
				foundOmni = true
				break
			}
		}
		if !foundOmni {
			t.Errorf("omni: expected index named c, got %v", axOmniIndexNames(tbl))
		}
	})
}

// axOracleColumnOrder returns the column names of the given table in
// ORDINAL_POSITION order, as seen by the MySQL container.
func axOracleColumnOrder(t *testing.T, mc *mysqlContainer, tableName string) []string {
	t.Helper()
	rows := oracleRows(t, mc, fmt.Sprintf(
		`SELECT COLUMN_NAME FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=%q ORDER BY ORDINAL_POSITION`, tableName))
	var names []string
	for _, row := range rows {
		if len(row) > 0 {
			names = append(names, asString(row[0]))
		}
	}
	return names
}

// axOmniColumnOrder returns the column names of the given omni table
// in positional order.
func axOmniColumnOrder(tbl *Table) []string {
	names := make([]string, 0, len(tbl.Columns))
	for _, col := range tbl.Columns {
		names = append(names, col.Name)
	}
	return names
}

// axOracleIndexNames returns sorted distinct index names for a table
// from information_schema.STATISTICS (PRIMARY excluded).
func axOracleIndexNames(t *testing.T, mc *mysqlContainer, tableName string) []string {
	t.Helper()
	rows := oracleRows(t, mc, fmt.Sprintf(
		`SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=%q AND INDEX_NAME <> 'PRIMARY'
         ORDER BY INDEX_NAME`, tableName))
	var names []string
	for _, row := range rows {
		if len(row) > 0 {
			names = append(names, asString(row[0]))
		}
	}
	return names
}

// axOmniIndexNames returns sorted non-primary index names for an omni table.
func axOmniIndexNames(tbl *Table) []string {
	var names []string
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			continue
		}
		names = append(names, idx.Name)
	}
	// Simple insertion sort to keep dep-free; tests compare joined strings.
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && names[j-1] > names[j]; j-- {
			names[j-1], names[j] = names[j], names[j-1]
		}
	}
	return names
}

// axExecHasError reports whether any ExecResult carries an error.
func axExecHasError(results []ExecResult) bool {
	for _, r := range results {
		if r.Error != nil {
			return true
		}
	}
	return false
}
