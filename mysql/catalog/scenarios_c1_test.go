package catalog

import (
	"sort"
	"strings"
	"testing"
)

// TestScenario_C1 covers section C1 (Name auto-generation) from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that
// both real MySQL 8.0 and the omni catalog agree on the auto-generated
// name for a given DDL input.
//
// Failures in omni assertions are NOT proof failures — they are
// recorded in mysql/catalog/scenarios_bug_queue/c1.md.
func TestScenario_C1(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 1.1 Foreign Key name — CREATE path (fresh counter) --------------
	t.Run("1_1_FK_name_create_path", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE child (
    a INT, CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES p(id),
    b INT, FOREIGN KEY (b) REFERENCES p(id)
);`
		runOnBoth(t, mc, c, ddl)

		// Oracle
		got := oracleFKNames(t, mc, "child")
		want := []string{"child_ibfk_1", "child_ibfk_5"}
		assertStringEq(t, "oracle FK names", strings.Join(got, ","), strings.Join(want, ","))

		// omni
		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Errorf("omni: table child missing")
			return
		}
		omniNames := omniFKNames(tbl)
		assertStringEq(t, "omni FK names", strings.Join(omniNames, ","), strings.Join(want, ","))
	})

	// --- 1.2 Foreign Key name — ALTER path (max+1 counter) ---------------
	t.Run("1_2_FK_name_alter_path", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE child (
    a INT, b INT,
    CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES p(id)
);
ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES p(id);`
		runOnBoth(t, mc, c, ddl)

		got := oracleFKNames(t, mc, "child")
		want := []string{"child_ibfk_20", "child_ibfk_21"}
		assertStringEq(t, "oracle FK names", strings.Join(got, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Errorf("omni: table child missing")
			return
		}
		omniNames := omniFKNames(tbl)
		assertStringEq(t, "omni FK names", strings.Join(omniNames, ","), strings.Join(want, ","))
	})

	// --- 1.3 Partition default naming p0..p{n-1} -------------------------
	t.Run("1_3_partition_default_names", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4;`)

		got := oraclePartitionNames(t, mc, "t")
		want := []string{"p0", "p1", "p2", "p3"}
		assertStringEq(t, "oracle partition names", strings.Join(got, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		var omniNames []string
		if tbl.Partitioning != nil {
			for _, p := range tbl.Partitioning.Partitions {
				omniNames = append(omniNames, p.Name)
			}
		}
		assertStringEq(t, "omni partition names", strings.Join(omniNames, ","), strings.Join(want, ","))
	})

	// --- 1.4 CHECK constraint auto-name (t_chk_1) ------------------------
	t.Run("1_4_check_auto_name", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, CHECK (a > 0));`)

		var got string
		oracleScan(t, mc, `SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'`, &got)
		assertStringEq(t, "oracle CHECK name", got, "t_chk_1")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniChecks := omniCheckNames(tbl)
		assertStringEq(t, "omni CHECK names", strings.Join(omniChecks, ","), "t_chk_1")
	})

	// --- 1.5 UNIQUE KEY auto-name uses field name ------------------------
	t.Run("1_5_unique_auto_name_field", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, UNIQUE KEY (a));`)

		var got string
		oracleScan(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0`, &got)
		assertStringEq(t, "oracle UNIQUE name", got, "a")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniIdx := omniUniqueIndexNames(tbl)
		assertStringEq(t, "omni unique index names", strings.Join(omniIdx, ","), "a")
	})

	// --- 1.6 UNIQUE KEY name collision appends _2 ------------------------
	t.Run("1_6_unique_collision_suffix", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, UNIQUE KEY a (a), UNIQUE KEY (a));`)

		rows := oracleRows(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0
            ORDER BY INDEX_NAME`)
		var got []string
		for _, r := range rows {
			got = append(got, asString(r[0]))
		}
		want := []string{"a", "a_2"}
		assertStringEq(t, "oracle unique index names", strings.Join(got, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniIdx := omniUniqueIndexNames(tbl)
		sort.Strings(omniIdx)
		assertStringEq(t, "omni unique index names", strings.Join(omniIdx, ","), strings.Join(want, ","))
	})

	// --- 1.7 PRIMARY KEY always named PRIMARY ----------------------------
	t.Run("1_7_primary_key_always_named_primary", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
    a INT,
    CONSTRAINT my_pk PRIMARY KEY (a)
);`)

		var got string
		oracleScan(t, mc, `SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_TYPE='PRIMARY KEY'`, &got)
		assertStringEq(t, "oracle PK constraint name", got, "PRIMARY")

		var idxName string
		oracleScan(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='PRIMARY'`, &idxName)
		assertStringEq(t, "oracle PK index name", idxName, "PRIMARY")

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		// No constraint named my_pk should exist.
		for _, con := range tbl.Constraints {
			if strings.EqualFold(con.Name, "my_pk") {
				t.Errorf("omni: unexpected constraint named %q (PK should be renamed to PRIMARY)", con.Name)
			}
		}
		// PK index should be named PRIMARY.
		var pkIdxName string
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				pkIdxName = idx.Name
				break
			}
		}
		assertStringEq(t, "omni PK index name", pkIdxName, "PRIMARY")
	})

	// --- 1.8 Non-PK index cannot be named PRIMARY ------------------------
	t.Run("1_8_primary_name_reserved", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Oracle: expect error from MySQL.
		_, mysqlErr := mc.db.ExecContext(mc.ctx, `CREATE TABLE t (a INT, UNIQUE KEY `+"`PRIMARY`"+` (a))`)
		if mysqlErr == nil {
			t.Errorf("oracle: expected error for UNIQUE KEY `PRIMARY`, got nil")
		} else if !strings.Contains(mysqlErr.Error(), "1280") && !strings.Contains(strings.ToLower(mysqlErr.Error()), "incorrect index name") {
			t.Errorf("oracle: expected ER_WRONG_NAME_FOR_INDEX (1280), got %v", mysqlErr)
		}

		// omni: should also reject. Use fresh catalog per attempt.
		results, err := c.Exec("CREATE TABLE t (a INT, UNIQUE KEY `PRIMARY` (a));", nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects index named PRIMARY", omniErrored, true)
	})

	// --- 1.9 Implicit index name from first key column -------------------
	t.Run("1_9_implicit_index_name_first_col", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
    a INT,
    b INT,
    c INT,
    KEY (b, c)
);`)

		rows := oracleRows(t, mc, `SELECT INDEX_NAME, COLUMN_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY INDEX_NAME, SEQ_IN_INDEX`)
		if len(rows) < 2 {
			t.Errorf("oracle: expected 2 STATISTICS rows, got %d", len(rows))
		}
		if len(rows) >= 1 {
			assertStringEq(t, "oracle index name", asString(rows[0][0]), "b")
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		var firstNonPKIdx *Index
		for _, idx := range tbl.Indexes {
			if !idx.Primary {
				firstNonPKIdx = idx
				break
			}
		}
		if firstNonPKIdx == nil {
			t.Errorf("omni: expected one non-PK index")
		} else {
			assertStringEq(t, "omni index name", firstNonPKIdx.Name, "b")
			if len(firstNonPKIdx.Columns) != 2 ||
				firstNonPKIdx.Columns[0].Name != "b" ||
				firstNonPKIdx.Columns[1].Name != "c" {
				t.Errorf("omni: expected columns [b,c], got %+v", firstNonPKIdx.Columns)
			}
		}
	})

	// --- 1.10 UNIQUE name fallback when first column is "PRIMARY" --------
	t.Run("1_10_unique_primary_column_fallback", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := "CREATE TABLE t (`PRIMARY` INT);\n" +
			"ALTER TABLE t ADD UNIQUE KEY (`PRIMARY`);"
		runOnBoth(t, mc, c, ddl)

		rows := oracleRows(t, mc, `SELECT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0`)
		if len(rows) != 1 {
			t.Errorf("oracle: expected 1 unique index row, got %d", len(rows))
		}
		if len(rows) == 1 {
			assertStringEq(t, "oracle unique index name", asString(rows[0][0]), "PRIMARY_2")
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniIdx := omniUniqueIndexNames(tbl)
		assertStringEq(t, "omni unique index names", strings.Join(omniIdx, ","), "PRIMARY_2")
	})

	// --- 1.11 Functional index auto-name functional_index[_N] ------------
	t.Run("1_11_functional_index_auto_name", func(t *testing.T) {
		t.Skip("structural: depends on functional-index hidden generated column synthesis and MySQL functional_index auto-naming")
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
    a INT,
    INDEX ((a + 1)),
    INDEX ((a * 2))
);`
		runOnBoth(t, mc, c, ddl)

		rows := oracleRows(t, mc, `SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY INDEX_NAME`)
		var got []string
		for _, r := range rows {
			got = append(got, asString(r[0]))
		}
		want := []string{"functional_index", "functional_index_2"}
		assertStringEq(t, "oracle functional index names", strings.Join(got, ","), strings.Join(want, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			// omni almost certainly fails to parse/load functional indexes.
			t.Errorf("omni: table t missing (functional index support gap)")
			return
		}
		var omniIdxNames []string
		for _, idx := range tbl.Indexes {
			if !idx.Primary {
				omniIdxNames = append(omniIdxNames, idx.Name)
			}
		}
		sort.Strings(omniIdxNames)
		assertStringEq(t, "omni functional index names", strings.Join(omniIdxNames, ","), strings.Join(want, ","))
	})

	// --- 1.12 Functional index hidden generated column name --------------
	t.Run("1_12_functional_index_hidden_col", func(t *testing.T) {
		t.Skip("structural: depends on functional-index hidden generated column synthesis and !hidden! name generation")
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT, INDEX fx ((a + 1), (a * 2)));`
		runOnBoth(t, mc, c, ddl)

		// Oracle: hidden columns are not in information_schema.COLUMNS by
		// default. Use a SELECT from the performance_schema / dd dumps via
		// SHOW CREATE TABLE as a loose check. Full name verification needs
		// a data-dictionary dump which we can't easily do from Go; document
		// and fall back to SHOW CREATE TABLE containing the expression.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if !strings.Contains(create, "`fx`") {
			t.Errorf("oracle: SHOW CREATE TABLE missing index fx: %s", create)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing (functional index support gap)")
			return
		}
		// omni will almost certainly not synthesize hidden generated columns
		// with MySQL's !hidden! naming scheme. Check and report.
		hasHidden := false
		for _, col := range tbl.Columns {
			if strings.HasPrefix(col.Name, "!hidden!") {
				hasHidden = true
				break
			}
		}
		assertBoolEq(t, "omni has hidden functional col", hasHidden, true)
	})

	// --- 1.13 CHECK constraint name is schema-scoped ---------------------
	t.Run("1_13_check_name_schema_scoped", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// First table with named check: both should accept.
		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, CONSTRAINT mychk CHECK (a > 0));`)

		// Second table with duplicate named check: MySQL errors with 3822.
		_, mysqlErr := mc.db.ExecContext(mc.ctx,
			`CREATE TABLE t2 (a INT, CONSTRAINT mychk CHECK (a < 100))`)
		if mysqlErr == nil {
			t.Errorf("oracle: expected ER_CHECK_CONSTRAINT_DUP_NAME, got nil")
		} else if !strings.Contains(mysqlErr.Error(), "3822") &&
			!strings.Contains(strings.ToLower(mysqlErr.Error()), "duplicate check constraint") {
			t.Errorf("oracle: expected 3822 Duplicate check constraint name, got %v", mysqlErr)
		}

		// omni: should also error.
		results, err := c.Exec(
			`CREATE TABLE t2 (a INT, CONSTRAINT mychk CHECK (a < 100));`, nil)
		omniErrored := err != nil
		if !omniErrored {
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
					break
				}
			}
		}
		assertBoolEq(t, "omni rejects cross-table duplicate CHECK name", omniErrored, true)
	})
}

// --- section-local helpers ------------------------------------------------

// oracleFKNames returns the FK constraint names on the given table, ordered
// alphabetically by CONSTRAINT_NAME.
func oracleFKNames(t *testing.T, mc *mysqlContainer, tableName string) []string {
	t.Helper()
	rows := oracleRows(t, mc, `SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
        WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+tableName+`'
              AND CONSTRAINT_TYPE='FOREIGN KEY'
        ORDER BY CONSTRAINT_NAME`)
	var out []string
	for _, r := range rows {
		out = append(out, asString(r[0]))
	}
	return out
}

// oraclePartitionNames returns partition names ordered by ordinal position.
func oraclePartitionNames(t *testing.T, mc *mysqlContainer, tableName string) []string {
	t.Helper()
	rows := oracleRows(t, mc, `SELECT PARTITION_NAME FROM information_schema.PARTITIONS
        WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+tableName+`'
        ORDER BY PARTITION_ORDINAL_POSITION`)
	var out []string
	for _, r := range rows {
		out = append(out, asString(r[0]))
	}
	return out
}

// omniFKNames returns FK constraint names for the table sorted alphabetically.
func omniFKNames(tbl *Table) []string {
	var out []string
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			out = append(out, con.Name)
		}
	}
	sort.Strings(out)
	return out
}

// omniCheckNames returns check constraint names sorted alphabetically.
func omniCheckNames(tbl *Table) []string {
	var out []string
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck {
			out = append(out, con.Name)
		}
	}
	sort.Strings(out)
	return out
}

// omniUniqueIndexNames returns names of unique indexes (excluding PK) sorted.
func omniUniqueIndexNames(tbl *Table) []string {
	var out []string
	for _, idx := range tbl.Indexes {
		if idx.Unique && !idx.Primary {
			out = append(out, idx.Name)
		}
	}
	sort.Strings(out)
	return out
}
