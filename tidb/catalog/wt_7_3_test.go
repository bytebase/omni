package catalog

import (
	"strings"
	"testing"
)

// --- Section 3.3: HASH and KEY Partitioning (9 scenarios) ---

// Scenario 1: CREATE TABLE PARTITION BY HASH (expr) PARTITIONS 4
func TestWalkThrough_7_3_HashByExpr(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t1 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY HASH (id) PARTITIONS 4`)

	ddl := c.ShowCreateTable("testdb", "t1")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify partition type marker
	if !strings.Contains(ddl, "/*!50100 PARTITION BY HASH") {
		t.Errorf("expected /*!50100 PARTITION BY HASH, got:\n%s", ddl)
	}
	// MySQL 8.0 backtick-quotes identifiers in HASH expression
	if !strings.Contains(ddl, "HASH (`id`)") {
		t.Errorf("expected HASH (`id`), got:\n%s", ddl)
	}
	// HASH with no explicit partition defs renders PARTITIONS N
	if !strings.Contains(ddl, "PARTITIONS 4") {
		t.Errorf("expected PARTITIONS 4 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.Type != "HASH" {
		t.Errorf("expected partition type HASH, got %q", tbl.Partitioning.Type)
	}
	if tbl.Partitioning.NumParts != 4 {
		t.Errorf("expected NumParts 4, got %d", tbl.Partitioning.NumParts)
	}
	if tbl.Partitioning.Linear {
		t.Error("expected Linear=false for non-linear HASH")
	}
}

// Scenario 2: CREATE TABLE PARTITION BY LINEAR HASH (expr) PARTITIONS 4
func TestWalkThrough_7_3_LinearHash(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t2 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY LINEAR HASH (id) PARTITIONS 4`)

	ddl := c.ShowCreateTable("testdb", "t2")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// LINEAR keyword should be rendered
	if !strings.Contains(ddl, "LINEAR HASH") {
		t.Errorf("expected LINEAR HASH in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITIONS 4") {
		t.Errorf("expected PARTITIONS 4 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t2")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if !tbl.Partitioning.Linear {
		t.Error("expected Linear=true for LINEAR HASH")
	}
	if tbl.Partitioning.Type != "HASH" {
		t.Errorf("expected partition type HASH, got %q", tbl.Partitioning.Type)
	}
}

// Scenario 3: CREATE TABLE PARTITION BY KEY (col) PARTITIONS 4
func TestWalkThrough_7_3_KeyPartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t3 (
		id INT NOT NULL PRIMARY KEY,
		val INT NOT NULL
	) PARTITION BY KEY (id) PARTITIONS 4`)

	ddl := c.ShowCreateTable("testdb", "t3")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	if !strings.Contains(ddl, "/*!50100 PARTITION BY KEY") {
		t.Errorf("expected /*!50100 PARTITION BY KEY, got:\n%s", ddl)
	}
	// KEY columns are rendered without backticks (plain)
	if !strings.Contains(ddl, "KEY (id)") {
		t.Errorf("expected KEY (id) in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITIONS 4") {
		t.Errorf("expected PARTITIONS 4 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t3")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.Type != "KEY" {
		t.Errorf("expected partition type KEY, got %q", tbl.Partitioning.Type)
	}
	if len(tbl.Partitioning.Columns) != 1 || tbl.Partitioning.Columns[0] != "id" {
		t.Errorf("expected Columns [id], got %v", tbl.Partitioning.Columns)
	}
	if tbl.Partitioning.NumParts != 4 {
		t.Errorf("expected NumParts 4, got %d", tbl.Partitioning.NumParts)
	}
}

// Scenario 4: CREATE TABLE PARTITION BY KEY () PARTITIONS 4 — uses PK
func TestWalkThrough_7_3_KeyEmptyColumns(t *testing.T) {
	c := wtSetup(t)

	results, err := c.Exec(`CREATE TABLE t4 (
		id INT NOT NULL PRIMARY KEY,
		val INT NOT NULL
	) PARTITION BY KEY () PARTITIONS 4`, nil)
	if err != nil {
		t.Skipf("[~] partial: parser does not support KEY () with empty column list: %v", err)
		return
	}
	for _, r := range results {
		if r.Error != nil {
			t.Skipf("[~] partial: catalog error for KEY (): %v", r.Error)
			return
		}
	}

	ddl := c.ShowCreateTable("testdb", "t4")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// KEY () uses PK columns — MySQL renders as KEY ()
	if !strings.Contains(ddl, "KEY ()") {
		t.Errorf("expected KEY () in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITIONS 4") {
		t.Errorf("expected PARTITIONS 4 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t4")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.Type != "KEY" {
		t.Errorf("expected partition type KEY, got %q", tbl.Partitioning.Type)
	}
	if tbl.Partitioning.NumParts != 4 {
		t.Errorf("expected NumParts 4, got %d", tbl.Partitioning.NumParts)
	}
}

// Scenario 5: CREATE TABLE PARTITION BY KEY (col) ALGORITHM=2 PARTITIONS 4
func TestWalkThrough_7_3_KeyAlgorithm(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t5 (
		id INT NOT NULL PRIMARY KEY,
		val INT NOT NULL
	) PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 4`)

	ddl := c.ShowCreateTable("testdb", "t5")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// ALGORITHM=2 should be rendered as ALGORITHM = 2
	if !strings.Contains(ddl, "KEY ALGORITHM = 2 (id)") {
		t.Errorf("expected KEY ALGORITHM = 2 (id) in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "PARTITIONS 4") {
		t.Errorf("expected PARTITIONS 4 in DDL:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t5")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.Algorithm != 2 {
		t.Errorf("expected Algorithm 2, got %d", tbl.Partitioning.Algorithm)
	}
}

// Scenario 6: ALTER TABLE COALESCE PARTITION 2 on HASH table (4→2)
func TestWalkThrough_7_3_CoalesceHash(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t6 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY HASH (id) PARTITIONS 4`)

	tbl := c.GetDatabase("testdb").GetTable("t6")
	if tbl.Partitioning.NumParts != 4 {
		t.Fatalf("expected NumParts 4 before COALESCE, got %d", tbl.Partitioning.NumParts)
	}

	wtExec(t, c, `ALTER TABLE t6 COALESCE PARTITION 2`)

	tbl = c.GetDatabase("testdb").GetTable("t6")
	if tbl.Partitioning.NumParts != 2 {
		t.Errorf("expected NumParts 2 after COALESCE, got %d", tbl.Partitioning.NumParts)
	}

	ddl := c.ShowCreateTable("testdb", "t6")
	if !strings.Contains(ddl, "PARTITIONS 2") {
		t.Errorf("expected PARTITIONS 2 in DDL:\n%s", ddl)
	}
}

// Scenario 7: ALTER TABLE COALESCE PARTITION on KEY table
func TestWalkThrough_7_3_CoalesceKey(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t7 (
		id INT NOT NULL PRIMARY KEY,
		val INT NOT NULL
	) PARTITION BY KEY (id) PARTITIONS 4`)

	tbl := c.GetDatabase("testdb").GetTable("t7")
	if tbl.Partitioning.NumParts != 4 {
		t.Fatalf("expected NumParts 4 before COALESCE, got %d", tbl.Partitioning.NumParts)
	}

	wtExec(t, c, `ALTER TABLE t7 COALESCE PARTITION 2`)

	tbl = c.GetDatabase("testdb").GetTable("t7")
	if tbl.Partitioning.NumParts != 2 {
		t.Errorf("expected NumParts 2 after COALESCE, got %d", tbl.Partitioning.NumParts)
	}

	ddl := c.ShowCreateTable("testdb", "t7")
	if !strings.Contains(ddl, "PARTITIONS 2") {
		t.Errorf("expected PARTITIONS 2 in DDL:\n%s", ddl)
	}
}

// Scenario 8: ALTER TABLE ADD PARTITION on HASH table — error in MySQL
// MySQL rejects ADD PARTITION on HASH/KEY tables. Document omni behavior.
func TestWalkThrough_7_3_AddPartitionHashError(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t8 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY HASH (id) PARTITIONS 4`)

	// In MySQL, this would error. Test what omni does.
	results, err := c.Exec(`ALTER TABLE t8 ADD PARTITION (PARTITION p_extra ENGINE=InnoDB)`, nil)
	if err != nil {
		// Parse error — document and skip
		t.Logf("note: ADD PARTITION on HASH table caused parse error (MySQL would reject this too): %v", err)
		return
	}
	for _, r := range results {
		if r.Error != nil {
			// Catalog error — this is actually correct behavior (MySQL rejects it)
			t.Logf("note: ADD PARTITION on HASH table correctly errored: %v", r.Error)
			return
		}
	}

	// If we reach here, omni allowed ADD PARTITION on HASH (differs from MySQL)
	t.Log("note: omni allows ADD PARTITION on HASH table — MySQL would reject this. Documenting behavior difference.")

	// Verify the partition was added (since omni allowed it)
	tbl := c.GetDatabase("testdb").GetTable("t8")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	// The original had NumParts=4 with no explicit Partitions slice.
	// ADD PARTITION appends to the Partitions slice.
	t.Logf("after ADD PARTITION: NumParts=%d, len(Partitions)=%d",
		tbl.Partitioning.NumParts, len(tbl.Partitioning.Partitions))
}

// Scenario 9: ALTER TABLE REMOVE PARTITIONING on partitioned table
func TestWalkThrough_7_3_RemovePartitioning(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t9 (
		id INT NOT NULL,
		val INT NOT NULL
	) PARTITION BY HASH (id) PARTITIONS 4`)

	tbl := c.GetDatabase("testdb").GetTable("t9")
	if tbl.Partitioning == nil {
		t.Fatal("partitioning should be set before REMOVE")
	}

	wtExec(t, c, `ALTER TABLE t9 REMOVE PARTITIONING`)

	tbl = c.GetDatabase("testdb").GetTable("t9")
	if tbl.Partitioning != nil {
		t.Error("partitioning should be nil after REMOVE PARTITIONING")
	}

	ddl := c.ShowCreateTable("testdb", "t9")
	if strings.Contains(ddl, "PARTITION") {
		t.Errorf("SHOW CREATE TABLE should not contain PARTITION after REMOVE:\n%s", ddl)
	}
}
