package catalog

import (
	"strings"
	"testing"
)

// --- Section 3.4: Subpartitions and Partition Options (8 scenarios) ---

// Scenario 1: RANGE partitions with SUBPARTITION BY HASH
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

	ddl := c.ShowCreateTable("testdb", "t1")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify SUBPARTITION BY HASH rendering
	if !strings.Contains(ddl, "SUBPARTITION BY HASH") {
		t.Errorf("expected SUBPARTITION BY HASH in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "SUBPARTITIONS 2") {
		t.Errorf("expected SUBPARTITIONS 2 in DDL:\n%s", ddl)
	}
	// Verify RANGE partitions still rendered
	if !strings.Contains(ddl, "PARTITION p0 VALUES LESS THAN (2000)") {
		t.Errorf("expected PARTITION p0 in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITION p2 VALUES LESS THAN MAXVALUE") {
		t.Errorf("expected PARTITION p2 with MAXVALUE in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.SubType != "HASH" {
		t.Errorf("expected SubType HASH, got %q", tbl.Partitioning.SubType)
	}
	if tbl.Partitioning.NumSubParts != 2 {
		t.Errorf("expected NumSubParts 2, got %d", tbl.Partitioning.NumSubParts)
	}
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Errorf("expected 3 partitions, got %d", len(tbl.Partitioning.Partitions))
	}
}

// Scenario 2: RANGE partitions with SUBPARTITION BY KEY
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

	ddl := c.ShowCreateTable("testdb", "t2")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify SUBPARTITION BY KEY rendering
	if !strings.Contains(ddl, "SUBPARTITION BY KEY") {
		t.Errorf("expected SUBPARTITION BY KEY in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "KEY (`id`)") {
		t.Errorf("expected KEY (`id`) in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "SUBPARTITIONS 2") {
		t.Errorf("expected SUBPARTITIONS 2 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t2")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.SubType != "KEY" {
		t.Errorf("expected SubType KEY, got %q", tbl.Partitioning.SubType)
	}
	if len(tbl.Partitioning.SubColumns) != 1 || tbl.Partitioning.SubColumns[0] != "id" {
		t.Errorf("expected SubColumns [id], got %v", tbl.Partitioning.SubColumns)
	}
	if tbl.Partitioning.NumSubParts != 2 {
		t.Errorf("expected NumSubParts 2, got %d", tbl.Partitioning.NumSubParts)
	}
}

// Scenario 3: Explicit subpartition definitions with names
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

	ddl := c.ShowCreateTable("testdb", "t3")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify explicit subpartition names rendered
	if !strings.Contains(ddl, "SUBPARTITION s0") {
		t.Errorf("expected SUBPARTITION s0 in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "SUBPARTITION s1") {
		t.Errorf("expected SUBPARTITION s1 in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "SUBPARTITION s2") {
		t.Errorf("expected SUBPARTITION s2 in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "SUBPARTITION s3") {
		t.Errorf("expected SUBPARTITION s3 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t3")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if len(tbl.Partitioning.Partitions) != 2 {
		t.Fatalf("expected 2 partitions, got %d", len(tbl.Partitioning.Partitions))
	}
	p0 := tbl.Partitioning.Partitions[0]
	if len(p0.SubPartitions) != 2 {
		t.Fatalf("expected 2 subpartitions for p0, got %d", len(p0.SubPartitions))
	}
	if p0.SubPartitions[0].Name != "s0" {
		t.Errorf("expected subpartition name s0, got %q", p0.SubPartitions[0].Name)
	}
	if p0.SubPartitions[1].Name != "s1" {
		t.Errorf("expected subpartition name s1, got %q", p0.SubPartitions[1].Name)
	}
	p1 := tbl.Partitioning.Partitions[1]
	if len(p1.SubPartitions) != 2 {
		t.Fatalf("expected 2 subpartitions for p1, got %d", len(p1.SubPartitions))
	}
	if p1.SubPartitions[0].Name != "s2" {
		t.Errorf("expected subpartition name s2, got %q", p1.SubPartitions[0].Name)
	}
	if p1.SubPartitions[1].Name != "s3" {
		t.Errorf("expected subpartition name s3, got %q", p1.SubPartitions[1].Name)
	}
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
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify ENGINE = InnoDB rendered per partition
	if !strings.Contains(ddl, "ENGINE = InnoDB") {
		t.Errorf("expected ENGINE = InnoDB per partition in DDL:\n%s", ddl)
	}

	// Count occurrences of ENGINE = InnoDB — should be 3 (one per partition)
	count := strings.Count(ddl, "ENGINE = InnoDB")
	if count < 3 {
		t.Errorf("expected at least 3 occurrences of ENGINE = InnoDB, got %d in DDL:\n%s", count, ddl)
	}

	// Verify catalog state: engine defaults to InnoDB even when not specified
	tbl := c.GetDatabase("testdb").GetTable("t4")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Errorf("expected 3 partitions, got %d", len(tbl.Partitioning.Partitions))
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
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify COMMENT rendered per partition
	if !strings.Contains(ddl, "COMMENT = 'first partition'") {
		t.Errorf("expected COMMENT = 'first partition' in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "COMMENT = 'last partition'") {
		t.Errorf("expected COMMENT = 'last partition' in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t5")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	p0 := tbl.Partitioning.Partitions[0]
	if p0.Comment != "first partition" {
		t.Errorf("expected comment 'first partition', got %q", p0.Comment)
	}
	p1 := tbl.Partitioning.Partitions[1]
	if p1.Comment != "last partition" {
		t.Errorf("expected comment 'last partition', got %q", p1.Comment)
	}
}

// Scenario 6: SUBPARTITIONS N count without explicit defs
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

	ddl := c.ShowCreateTable("testdb", "t6")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify SUBPARTITIONS 3 rendered (no explicit subpartition names)
	if !strings.Contains(ddl, "SUBPARTITIONS 3") {
		t.Errorf("expected SUBPARTITIONS 3 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t6")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.NumSubParts != 3 {
		t.Errorf("expected NumSubParts 3, got %d", tbl.Partitioning.NumSubParts)
	}
	// Auto-generated subpartitions should be present on each partition def.
	for _, p := range tbl.Partitioning.Partitions {
		if len(p.SubPartitions) != 3 {
			t.Errorf("expected 3 auto-generated subpartitions for partition %s, got %d", p.Name, len(p.SubPartitions))
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
