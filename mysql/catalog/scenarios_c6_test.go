package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C6 covers Section C6 (Partition defaults) of the
// mysql-implicit-behavior starmap. Each subtest runs DDL against both a real
// MySQL 8.0 container and the omni catalog and asserts the observable state
// agrees.
//
// Uses helpers from scenarios_helpers_test.go and reuses:
//   - oraclePartitionNames (scenarios_c1_test.go)
//
// Section-local helpers use the `c6` prefix to avoid collisions.
func TestScenario_C6(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 6.1 HASH without PARTITIONS defaults to 1 ----------------------
	t.Run("6_1_HASH_partitions_default_1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT) PARTITION BY HASH(id)`)

		names := oraclePartitionNames(t, mc, "t")
		wantNames := []string{"p0"}
		assertStringEq(t, "oracle partition names",
			strings.Join(names, ","), strings.Join(wantNames, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		omniNames := c6OmniPartitionNames(tbl)
		assertStringEq(t, "omni partition names",
			strings.Join(omniNames, ","), strings.Join(wantNames, ","))
	})

	// --- 6.2 SUBPARTITIONS default to 1 if not specified ----------------
	t.Run("6_2_subpartitions_default_1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT, d DATE)
            PARTITION BY RANGE(YEAR(d))
            SUBPARTITION BY HASH(id)
            (PARTITION p0 VALUES LESS THAN (2000),
             PARTITION p1 VALUES LESS THAN MAXVALUE);`)

		// Oracle: expect 2 subpartitions total (1 per parent).
		rows := oracleRows(t, mc, `SELECT PARTITION_NAME, SUBPARTITION_NAME
            FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION`)
		if len(rows) != 2 {
			t.Errorf("oracle subpartition rows: got %d, want 2", len(rows))
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		total := 0
		for _, p := range tbl.Partitioning.Partitions {
			total += len(p.SubPartitions)
		}
		assertIntEq(t, "omni subpartition count", total, 2)
	})

	// --- 6.3 Partition ENGINE defaults to table ENGINE ------------------
	t.Run("6_3_partition_engine_defaults_to_table", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) ENGINE=InnoDB PARTITION BY HASH(id) PARTITIONS 2`)

		// Oracle: SHOW CREATE TABLE renders table-level ENGINE=InnoDB.
		sc := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if !strings.Contains(sc, "ENGINE=InnoDB") {
			t.Errorf("oracle SHOW CREATE: expected ENGINE=InnoDB, got %s", sc)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		for i, p := range tbl.Partitioning.Partitions {
			// Empty per-partition engine is OK as long as table-level engine
			// renders it in SHOW CREATE; but a concrete non-InnoDB would be
			// wrong.
			if p.Engine != "" && !strings.EqualFold(p.Engine, "InnoDB") {
				t.Errorf("omni: partition %d engine %q != InnoDB", i, p.Engine)
			}
		}
	})

	// --- 6.4 KEY ALGORITHM default 2 ------------------------------------
	t.Run("6_4_key_algorithm_default_2", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT) PARTITION BY KEY(id) PARTITIONS 4`)

		// Oracle SHOW CREATE TABLE should NOT mention ALGORITHM=1 (algo 2 is
		// default, so it is elided from the output).
		sc := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(sc, "ALGORITHM=1") {
			t.Errorf("oracle SHOW CREATE should not have ALGORITHM=1: %s", sc)
		}

		// omni catalog: algorithm should default to 2 (or be 0 == unset;
		// either way it must not be 1).
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		if tbl.Partitioning.Algorithm == 1 {
			t.Errorf("omni: Algorithm=1 leaked as default, want 2 or 0")
		}
	})

	// --- 6.5 KEY() empty column list → PK columns ----------------------
	t.Run("6_5_key_empty_columns_defaults_to_PK", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT PRIMARY KEY, v INT) PARTITION BY KEY() PARTITIONS 4`)

		// Oracle: PARTITION_EXPRESSION column should name `id`.
		var expr string
		oracleScan(t, mc, `SELECT COALESCE(PARTITION_EXPRESSION,'')
            FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            LIMIT 1`, &expr)
		if !strings.Contains(expr, "id") {
			t.Errorf("oracle partition expression: got %q, want containing 'id'", expr)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		got := strings.Join(tbl.Partitioning.Columns, ",")
		if got != "id" {
			t.Errorf("omni partition columns: got %q, want %q", got, "id")
		}
	})

	// --- 6.6 LINEAR HASH / LINEAR KEY preserved -------------------------
	t.Run("6_6_linear_preserved", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) PARTITION BY LINEAR HASH(id) PARTITIONS 4`)

		sc := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if !strings.Contains(sc, "LINEAR HASH") {
			t.Errorf("oracle SHOW CREATE: expected LINEAR HASH, got %s", sc)
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		assertBoolEq(t, "omni Linear", tbl.Partitioning.Linear, true)
	})

	// --- 6.7 RANGE/LIST require explicit partition definitions ---------
	t.Run("6_7_range_partitions_n_shortcut_rejected", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (id INT) PARTITION BY RANGE(id) PARTITIONS 4`

		// Oracle: must reject.
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err == nil {
			t.Errorf("oracle: expected error for RANGE without definitions")
		}

		// omni: should also reject.
		results, err := c.Exec(ddl, nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: expected error for RANGE without definitions, got nil")
		}
	})

	// --- 6.8 MAXVALUE must appear in the last RANGE partition only -----
	t.Run("6_8_maxvalue_last_only", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (id INT) PARTITION BY RANGE(id)
            (PARTITION p0 VALUES LESS THAN MAXVALUE,
             PARTITION p1 VALUES LESS THAN (100))`

		if _, err := mc.db.ExecContext(mc.ctx, ddl); err == nil {
			t.Errorf("oracle: expected error for misplaced MAXVALUE")
		}

		results, err := c.Exec(ddl, nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: expected error for misplaced MAXVALUE, got nil")
		}
	})

	// --- 6.9 LIST comparison semantics (docs + round-trip) ------------
	t.Run("6_9_list_equality_roundtrip", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (c INT) PARTITION BY LIST(c)
             (PARTITION p0 VALUES IN (1,2), PARTITION p1 VALUES IN (3,4))`)

		names := oraclePartitionNames(t, mc, "t")
		wantNames := []string{"p0", "p1"}
		assertStringEq(t, "oracle partition names",
			strings.Join(names, ","), strings.Join(wantNames, ","))

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table or partitioning missing")
			return
		}
		assertStringEq(t, "omni partition type", tbl.Partitioning.Type, "LIST")
		assertIntEq(t, "omni partition count",
			len(tbl.Partitioning.Partitions), 2)
	})

	// --- 6.10 LIST DEFAULT partition --------------------------------
	t.Run("6_10_list_default_partition", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (c INT) PARTITION BY LIST(c)
             (PARTITION p0 VALUES IN (1,2), PARTITION pd VALUES IN (DEFAULT))`)

		names := oraclePartitionNames(t, mc, "t")
		wantNames := []string{"p0", "pd"}
		assertStringEq(t, "oracle partition names",
			strings.Join(names, ","), strings.Join(wantNames, ","))

		// Oracle: pd's PARTITION_DESCRIPTION should be DEFAULT.
		var desc string
		oracleScan(t, mc, `SELECT COALESCE(PARTITION_DESCRIPTION,'')
            FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND PARTITION_NAME='pd'`,
			&desc)
		if !strings.EqualFold(desc, "DEFAULT") {
			t.Errorf("oracle: pd description = %q, want DEFAULT", desc)
		}

		// omni: expect DEFAULT token to round-trip on the last partition's
		// ValueExpr.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil || len(tbl.Partitioning.Partitions) < 2 {
			t.Errorf("omni: table/partitioning missing or short")
			return
		}
		got := strings.ToUpper(tbl.Partitioning.Partitions[1].ValueExpr)
		if !strings.Contains(got, "DEFAULT") {
			t.Errorf("omni: pd ValueExpr = %q, want containing DEFAULT", got)
		}
	})

	// --- 6.11 Partition function result must be INTEGER --------------
	t.Run("6_11_partition_expr_non_integer_rejected", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a VARCHAR(10), b VARCHAR(10))
            PARTITION BY RANGE(CONCAT(a,b))
            (PARTITION p0 VALUES LESS THAN ('m'),
             PARTITION p1 VALUES LESS THAN MAXVALUE)`

		if _, err := mc.db.ExecContext(mc.ctx, ddl); err == nil {
			t.Errorf("oracle: expected error for non-integer partition expr")
		}

		results, err := c.Exec(ddl, nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: expected error for non-integer partition expr, got nil")
		}
	})

	// --- 6.12 TIMESTAMP requires UNIX_TIMESTAMP wrapping --------------
	t.Run("6_12_timestamp_requires_unix_timestamp", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		badDDL := `CREATE TABLE t1 (ts TIMESTAMP) PARTITION BY RANGE(ts)
            (PARTITION p0 VALUES LESS THAN (100))`

		if _, err := mc.db.ExecContext(mc.ctx, badDDL); err == nil {
			t.Errorf("oracle: expected error for bare TIMESTAMP partition expr")
		}
		results, err := c.Exec(badDDL, nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: expected error for bare TIMESTAMP partition expr")
		}

		// The wrapped form must succeed on both.
		goodDDL := `CREATE TABLE t2 (ts TIMESTAMP NOT NULL)
            PARTITION BY RANGE(UNIX_TIMESTAMP(ts))
            (PARTITION p0 VALUES LESS THAN (100),
             PARTITION p1 VALUES LESS THAN MAXVALUE)`
		runOnBoth(t, mc, c, goodDDL)

		tbl := c.GetDatabase("testdb").GetTable("t2")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: t2 partitioning missing")
			return
		}
		if !strings.Contains(strings.ToUpper(tbl.Partitioning.Expr), "UNIX_TIMESTAMP") {
			t.Errorf("omni: expected expr containing UNIX_TIMESTAMP, got %q",
				tbl.Partitioning.Expr)
		}
	})

	// --- 6.13 UNIQUE KEY must cover partition expression columns -----
	t.Run("6_13_unique_key_must_cover_partition_cols", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT, b INT, UNIQUE KEY (a))
            PARTITION BY HASH(b) PARTITIONS 4`

		if _, err := mc.db.ExecContext(mc.ctx, ddl); err == nil {
			t.Errorf("oracle: expected error when UNIQUE KEY excludes partition col")
		}

		results, err := c.Exec(ddl, nil)
		omniErr := err
		if omniErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: expected error, got nil (schema silently accepted)")
		}
	})

	// --- 6.14 Per-partition options preserved verbatim --------------
	t.Run("6_14_per_partition_options_preserved", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) PARTITION BY HASH(id)
             (PARTITION p0 COMMENT='first' ENGINE=InnoDB,
              PARTITION p1 COMMENT='second' ENGINE=InnoDB)`)

		// Oracle: PARTITION_COMMENT should be set per partition.
		rows := oracleRows(t, mc, `SELECT PARTITION_NAME, PARTITION_COMMENT
            FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY PARTITION_ORDINAL_POSITION`)
		if len(rows) != 2 {
			t.Errorf("oracle rows: got %d, want 2", len(rows))
		} else {
			if s, _ := rows[0][1].(string); s != "first" {
				t.Errorf("oracle p0 comment: got %q, want %q", s, "first")
			}
			if s, _ := rows[1][1].(string); s != "second" {
				t.Errorf("oracle p1 comment: got %q, want %q", s, "second")
			}
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil || len(tbl.Partitioning.Partitions) < 2 {
			t.Errorf("omni: table/partitioning missing or short")
			return
		}
		assertStringEq(t, "omni p0 comment", tbl.Partitioning.Partitions[0].Comment, "first")
		assertStringEq(t, "omni p1 comment", tbl.Partitioning.Partitions[1].Comment, "second")
	})

	// --- 6.15 Subpartition options inherit handling ----------------
	t.Run("6_15_subpartition_options", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (id INT, d DATE)
            ENGINE=InnoDB
            PARTITION BY RANGE(YEAR(d))
            SUBPARTITION BY HASH(id) SUBPARTITIONS 2
            (PARTITION p0 VALUES LESS THAN (2000)
                (SUBPARTITION s0 COMMENT='sa', SUBPARTITION s1 COMMENT='sb'),
             PARTITION p1 VALUES LESS THAN MAXVALUE
                (SUBPARTITION s2, SUBPARTITION s3))`)

		// Oracle: 4 subpartitions total.
		rows := oracleRows(t, mc, `SELECT SUBPARTITION_NAME
            FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION`)
		assertIntEq(t, "oracle subpartition count", len(rows), 4)

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table/partitioning missing")
			return
		}
		total := 0
		for _, p := range tbl.Partitioning.Partitions {
			total += len(p.SubPartitions)
		}
		assertIntEq(t, "omni subpartition count", total, 4)
	})

	// --- 6.16 ALTER ADD PARTITION auto-naming ---------------------
	t.Run("6_16_add_partition_auto_naming", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3`)

		// Try ALTER on both sides. omni may not support this — we report
		// either outcome as assertion failures, not panics.
		addDDL := `ALTER TABLE t ADD PARTITION PARTITIONS 2`
		if _, err := mc.db.ExecContext(mc.ctx, addDDL); err != nil {
			t.Errorf("oracle: ALTER ADD PARTITION failed: %v", err)
		}
		names := oraclePartitionNames(t, mc, "t")
		want := []string{"p0", "p1", "p2", "p3", "p4"}
		assertStringEq(t, "oracle partition names after ADD",
			strings.Join(names, ","), strings.Join(want, ","))

		// omni side: tolerant — report but do not panic.
		results, err := c.Exec(addDDL, nil)
		if err != nil {
			t.Errorf("omni: ALTER ADD PARTITION parse error: %v", err)
			return
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: ALTER ADD PARTITION exec error: %v", r.Error)
			}
		}
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table/partitioning missing")
			return
		}
		omniNames := c6OmniPartitionNames(tbl)
		assertStringEq(t, "omni partition names after ADD",
			strings.Join(omniNames, ","), strings.Join(want, ","))
	})

	// --- 6.17 COALESCE PARTITION removes last N ------------------
	t.Run("6_17_coalesce_partition_tail", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 6`)

		coalesce := `ALTER TABLE t COALESCE PARTITION 2`
		if _, err := mc.db.ExecContext(mc.ctx, coalesce); err != nil {
			t.Errorf("oracle: ALTER COALESCE PARTITION failed: %v", err)
		}
		names := oraclePartitionNames(t, mc, "t")
		want := []string{"p0", "p1", "p2", "p3"}
		assertStringEq(t, "oracle partition names after COALESCE",
			strings.Join(names, ","), strings.Join(want, ","))

		results, err := c.Exec(coalesce, nil)
		if err != nil {
			t.Errorf("omni: COALESCE parse error: %v", err)
			return
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni: COALESCE exec error: %v", r.Error)
			}
		}
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil || tbl.Partitioning == nil {
			t.Errorf("omni: table/partitioning missing")
			return
		}
		omniNames := c6OmniPartitionNames(tbl)
		assertStringEq(t, "omni partition names after COALESCE",
			strings.Join(omniNames, ","), strings.Join(want, ","))
	})
}

// c6OmniPartitionNames returns the partition names in order from the omni
// catalog table. Returns nil if partitioning is missing.
func c6OmniPartitionNames(tbl *Table) []string {
	if tbl == nil || tbl.Partitioning == nil {
		return nil
	}
	names := make([]string, 0, len(tbl.Partitioning.Partitions))
	for _, p := range tbl.Partitioning.Partitions {
		names = append(names, p.Name)
	}
	return names
}

