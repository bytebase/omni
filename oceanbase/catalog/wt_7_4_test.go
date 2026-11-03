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
	if strings.Contains(ddl, "SUBPARTITIONS ") {
		t.Errorf("explicit subpartition definitions should not render SUBPARTITIONS count:\n%s", ddl)
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

// Scenario 5: Partition with COMMENT option. MySQL keeps a COMMENT per partition and
// renders it, which this test held when the package was forked from mysql. OceanBase
// takes only ID and ENGINE = InnoDB on a partition and answers 1064 for COMMENT, with
// or without '=' (parser partitionOptionCases), so the catalog must refuse the table.
func TestWalkThrough_7_4_PartitionCommentOption(t *testing.T) {
	c := wtSetup(t)

	if _, err := c.Exec(`CREATE TABLE t5 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY RANGE (id) (
	    PARTITION p0 VALUES LESS THAN (100) COMMENT='first partition',
	    PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT='last partition'
	  )`, nil); err == nil {
		t.Error("a partition COMMENT parsed; OceanBase answers 1064")
	}
	if c.GetDatabase("testdb").GetTable("t5") != nil {
		t.Error("the refused table was created")
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
	if strings.Contains(ddl, "SUBPARTITION p0sp") || strings.Contains(ddl, "SUBPARTITION p1sp") {
		t.Errorf("default subpartition names should not be rendered when SUBPARTITIONS count is used:\n%s", ddl)
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

func TestWalkThrough_7_4_DefaultSubpartitionsDoNotRenderImplicitCount(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t_default_subparts (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  (
	    PARTITION p0 VALUES LESS THAN (2000),
	    PARTITION p1 VALUES LESS THAN MAXVALUE
	  )`)

	ddl := c.ShowCreateTable("testdb", "t_default_subparts")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}
	if strings.Contains(ddl, "SUBPARTITIONS ") {
		t.Errorf("implicit default subpartition count should not be rendered:\n%s", ddl)
	}
	if strings.Contains(ddl, "SUBPARTITION p0sp") || strings.Contains(ddl, "SUBPARTITION p1sp") {
		t.Errorf("implicit default subpartition names should not be rendered:\n%s", ddl)
	}
}

// TestWalkThrough_7_4_ExplicitSubpartitionsInheritParentOptions: a partition that lists
// its subpartitions renders its options on them, not before the list. MySQL's form of
// this test carried COMMENT, which OceanBase refuses on a partition and a subpartition
// (parser partitionOptionCases), so it uses ENGINE = InnoDB, the one option both take.
func TestWalkThrough_7_4_ExplicitSubpartitionsInheritParentOptions(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t_subpart_engine (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  (
	    PARTITION p0 VALUES LESS THAN (2000) ENGINE = InnoDB (
	      SUBPARTITION s0 ENGINE = InnoDB,
	      SUBPARTITION s1
	    ),
	    PARTITION p1 VALUES LESS THAN MAXVALUE (
	      SUBPARTITION s2,
	      SUBPARTITION s3
	    )
	  )`)

	ddl := c.ShowCreateTable("testdb", "t_subpart_engine")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}
	if strings.Contains(ddl, "PARTITION p0 VALUES LESS THAN (2000) ENGINE = InnoDB") {
		t.Errorf("parent partition options should not be rendered directly before explicit subpartition list:\n%s", ddl)
	}
	for _, sub := range []string{"s0", "s1", "s2", "s3"} {
		if !strings.Contains(ddl, "SUBPARTITION "+sub+" ENGINE = InnoDB") {
			t.Errorf("expected subpartition %s to carry ENGINE = InnoDB:\n%s", sub, ddl)
		}
	}

	if _, err := c.Exec(`CREATE TABLE t_subpart_comment (id INT NOT NULL, purchased DATE NOT NULL)
	  PARTITION BY RANGE (YEAR(purchased)) SUBPARTITION BY HASH (TO_DAYS(purchased))
	  (PARTITION p0 VALUES LESS THAN (2000) (SUBPARTITION s0 COMMENT='s0c', SUBPARTITION s1))`, nil); err == nil {
		t.Error("a subpartition COMMENT parsed; OceanBase answers 1064")
	}
}

// TestWalkThrough_7_4_ExplicitSubpartitionsMayDifferFromDeclaredCount: MySQL refuses
// explicit subpartition lists that disagree with SUBPARTITIONS N (1064), and this test
// held that rule when the package was forked from mysql. OceanBase takes this statement
// as written: 5.7.25-OceanBase_CE-v4.3.5.6 creates p0 with s0, s1 and p1 with s2, s3.
func TestWalkThrough_7_4_ExplicitSubpartitionsMayDifferFromDeclaredCount(t *testing.T) {
	c := wtSetup(t)

	results, err := c.Exec(`CREATE TABLE t_subparts (
		id INT NOT NULL,
		purchased DATE NOT NULL
	) PARTITION BY RANGE (YEAR(purchased))
	  SUBPARTITION BY HASH (TO_DAYS(purchased))
	  SUBPARTITIONS 3
	  (
	    PARTITION p0 VALUES LESS THAN (2000) (
	      SUBPARTITION s0,
	      SUBPARTITION s1
	    ),
	    PARTITION p1 VALUES LESS THAN MAXVALUE (
	      SUBPARTITION s2,
	      SUBPARTITION s3
	    )
	  )`, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}
	if results[0].Error != nil {
		t.Fatalf("OceanBase accepts the statement, the catalog refused it: %v", results[0].Error)
	}
	var got []string
	for _, pd := range c.GetDatabase("testdb").GetTable("t_subparts").Partitioning.Partitions {
		for _, sp := range pd.SubPartitions {
			got = append(got, pd.Name+"."+sp.Name)
		}
	}
	if want := "p0.s0 p0.s1 p1.s2 p1.s3"; strings.Join(got, " ") != want {
		t.Errorf("subpartitions = %v, want %s", got, want)
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
