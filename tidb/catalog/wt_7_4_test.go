package catalog

import (
	"strings"
	"testing"
)

// --- Section 3.4: Subpartitions and Partition Options (8 scenarios) ---

// Scenario 1: RANGE partitions with SUBPARTITION BY HASH
// TiDB v8.5.5 reads the SUBPARTITION clause and keeps none of it (read on
// 2026-11-06).
func TestWalkThrough_7_4_RangeSubpartByHash(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t1 (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  SUBPARTITIONS 2
	  (
	    PARTITION p0 VALUES LESS THAN (2000),
	    PARTITION p1 VALUES LESS THAN (2010),
	    PARTITION p2 VALUES LESS THAN MAXVALUE
	  )`)

	assertNoSubpartitions(t, c, "t1")
}

// Scenario 2: RANGE partitions with SUBPARTITION BY KEY
// TiDB v8.5.5 reads the SUBPARTITION clause and keeps none of it (read on
// 2026-11-06).
func TestWalkThrough_7_4_RangeSubpartByKey(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t2 (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY KEY (id)
	  SUBPARTITIONS 2
	  (
	    PARTITION p0 VALUES LESS THAN (2000),
	    PARTITION p1 VALUES LESS THAN MAXVALUE
	  )`)

	assertNoSubpartitions(t, c, "t2")
}

// Scenario 3: Explicit subpartition definitions with names
// TiDB v8.5.5 reads the SUBPARTITION clause and keeps none of it (read on
// 2026-11-06).
func TestWalkThrough_7_4_ExplicitSubpartNames(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t3 (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  (
	    PARTITION p0 VALUES LESS THAN (2000) (
	      SUBPARTITION s0,
	      SUBPARTITION s1
	    ),
	    PARTITION p1 VALUES LESS THAN MAXVALUE (
	      SUBPARTITION s2,
	      SUBPARTITION s3
	    )
	  )`)

	assertNoSubpartitions(t, c, "t3")
}

// Scenario 4: Partition with ENGINE option
func TestWalkThrough_7_4_PartitionEngineOption(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t4 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY RANGE (id) (
	    PARTITION p0 VALUES LESS THAN (100) ENGINE=InnoDB,
	    PARTITION p1 VALUES LESS THAN (200) ENGINE=InnoDB,
	    PARTITION p2 VALUES LESS THAN MAXVALUE ENGINE=InnoDB
	  )`)

	ddl := c.ShowCreateTable("testdb", "t4")
	// TiDB keeps no partition ENGINE: v8.5.5 writes these partitions with
	// none (read on 2026-11-06).
	if strings.Contains(ddl, "ENGINE = ") || !strings.Contains(ddl, "PARTITION `p0` VALUES LESS THAN (100),") {
		t.Errorf("expected partitions with no ENGINE in DDL:\n%s", ddl)
	}
	tbl := c.GetDatabase("testdb").GetTable("t4")
	if tbl.Partitioning == nil || len(tbl.Partitioning.Partitions) != 3 {
		t.Fatalf("partitioning: %+v", tbl.Partitioning)
	}
}

// Scenario 5: Partition with COMMENT option
func TestWalkThrough_7_4_PartitionCommentOption(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t5 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY RANGE (id) (
	    PARTITION p0 VALUES LESS THAN (100) COMMENT='first partition',
	    PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT='last partition'
	  )`)

	ddl := c.ShowCreateTable("testdb", "t5")
	// TiDB writes a partition comment as COMMENT '...', with no equals sign.
	if !strings.Contains(ddl, "PARTITION `p0` VALUES LESS THAN (100) COMMENT 'first partition'") {
		t.Errorf("expected COMMENT 'first partition' in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITION `p1` VALUES LESS THAN (MAXVALUE) COMMENT 'last partition'") {
		t.Errorf("expected COMMENT 'last partition' in DDL:\n%s", ddl)
	}
	tbl := c.GetDatabase("testdb").GetTable("t5")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if p0 := tbl.Partitioning.Partitions[0]; p0.Comment != "first partition" {
		t.Errorf("expected comment 'first partition', got %q", p0.Comment)
	}
}

// Scenario 6: SUBPARTITIONS N count without explicit defs
// TiDB v8.5.5 reads the SUBPARTITION clause and keeps none of it (read on
// 2026-11-06).
func TestWalkThrough_7_4_SubpartitionsNCount(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t6 (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  SUBPARTITIONS 3
	  (
	    PARTITION p0 VALUES LESS THAN (2000),
	    PARTITION p1 VALUES LESS THAN MAXVALUE
	  )`)

	assertNoSubpartitions(t, c, "t6")
}

// assertNoSubpartitions holds a table to what TiDB keeps of a SUBPARTITION
// clause: nothing, in the catalog or in SHOW CREATE TABLE.
func assertNoSubpartitions(t *testing.T, c *Catalog, table string) {
	t.Helper()
	ddl := c.ShowCreateTable("testdb", table)
	if strings.Contains(ddl, "SUBPARTITION") {
		t.Errorf("TiDB keeps no subpartitions, the DDL has them:\n%s", ddl)
	}
	pi := c.GetDatabase("testdb").GetTable(table).Partitioning
	if pi == nil {
		t.Fatal("partitioning is nil")
	}
	if pi.SubType != "" || pi.NumSubParts != 0 {
		t.Errorf("SubType=%q NumSubParts=%d", pi.SubType, pi.NumSubParts)
	}
	for _, p := range pi.Partitions {
		if len(p.SubPartitions) != 0 {
			t.Errorf("partition %s kept %d subpartitions", p.Name, len(p.SubPartitions))
		}
	}
}

// Scenario 7: ALTER TABLE TRUNCATE PARTITION — no structural change
func TestWalkThrough_7_4_TruncatePartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t7 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY RANGE (id) (
	    PARTITION p0 VALUES LESS THAN (100),
	    PARTITION p1 VALUES LESS THAN (200),
	    PARTITION p2 VALUES LESS THAN MAXVALUE
	  )`)

	ddlBefore := c.ShowCreateTable("testdb", "t7")
	if ddlBefore == "" {
		t.Fatal("ShowCreateTable returned empty string before TRUNCATE")
	}

	wtExec(t, c, `ALTER TABLE t7 TRUNCATE PARTITION p1`)

	ddlAfter := c.ShowCreateTable("testdb", "t7")

	// TRUNCATE PARTITION is a data-only operation; structure unchanged
	if ddlBefore != ddlAfter {
		t.Errorf("SHOW CREATE TABLE changed after TRUNCATE PARTITION:\n--- before ---\n%s\n--- after ---\n%s",
			ddlBefore, ddlAfter)
	}

	// Verify catalog state unchanged
	tbl := c.GetDatabase("testdb").GetTable("t7")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil after TRUNCATE")
	}
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Errorf("expected 3 partitions after TRUNCATE, got %d", len(tbl.Partitioning.Partitions))
	}
}

// Scenario 8: ALTER TABLE EXCHANGE PARTITION — validation only, structure unchanged
func TestWalkThrough_7_4_ExchangePartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t8 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY RANGE (id) (
	    PARTITION p0 VALUES LESS THAN (100),
	    PARTITION p1 VALUES LESS THAN (200),
	    PARTITION p2 VALUES LESS THAN MAXVALUE
	  )`)

	// Create the exchange target table (non-partitioned)
	wtExec(t, c, `CREATE TABLE t8_swap (
		id INT NOT NULL,
		val INT NOT NULL
	)`)

	ddlBefore := c.ShowCreateTable("testdb", "t8")
	ddlSwapBefore := c.ShowCreateTable("testdb", "t8_swap")

	wtExec(t, c, `ALTER TABLE t8 EXCHANGE PARTITION p1 WITH TABLE t8_swap`)

	ddlAfter := c.ShowCreateTable("testdb", "t8")
	ddlSwapAfter := c.ShowCreateTable("testdb", "t8_swap")

	// EXCHANGE PARTITION is a data operation; structure unchanged for both tables
	if ddlBefore != ddlAfter {
		t.Errorf("partitioned table structure changed after EXCHANGE:\n--- before ---\n%s\n--- after ---\n%s",
			ddlBefore, ddlAfter)
	}
	if ddlSwapBefore != ddlSwapAfter {
		t.Errorf("swap table structure changed after EXCHANGE:\n--- before ---\n%s\n--- after ---\n%s",
			ddlSwapBefore, ddlSwapAfter)
	}

	// Verify both tables still exist
	if c.GetDatabase("testdb").GetTable("t8") == nil {
		t.Error("partitioned table t8 should still exist")
	}
	if c.GetDatabase("testdb").GetTable("t8_swap") == nil {
		t.Error("swap table t8_swap should still exist")
	}

	// Verify exchange with nonexistent table errors
	results, err := c.Exec(`ALTER TABLE t8 EXCHANGE PARTITION p1 WITH TABLE nonexistent`, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Error("expected error when exchanging with nonexistent table")
	}
}
