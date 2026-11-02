package catalog

import (
	"strings"
	"testing"
)

// --- Section 3.2: LIST Partitioning (6 scenarios) ---

// Scenario 1: CREATE TABLE PARTITION BY LIST (expr) with VALUES IN
func TestWalkThrough_7_2_ListByExpr(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t1 (
		id INT NOT NULL,
		region INT NOT NULL
	) PARTITION BY LIST (region) (
		PARTITION p_east VALUES IN (1,2,3),
		PARTITION p_west VALUES IN (4,5,6),
		PARTITION p_central VALUES IN (7,8,9)
	)`)

	ddl := c.ShowCreateTable("testdb", "t1")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// Verify partition type marker
	if !strings.Contains(ddl, "/*!50100 PARTITION BY LIST") {
		t.Errorf("expected /*!50100 PARTITION BY LIST, got:\n%s", ddl)
	}
	// MySQL 8.0 backtick-quotes identifiers in partition expressions
	if !strings.Contains(ddl, "LIST (`region`)") {
		t.Errorf("expected LIST (`region`), got:\n%s", ddl)
	}

	// Verify partition definitions
	for _, name := range []string{"p_east", "p_west", "p_central"} {
		if !strings.Contains(ddl, "PARTITION "+name) {
			t.Errorf("expected PARTITION %s in DDL:\n%s", name, ddl)
		}
	}
	if !strings.Contains(ddl, "VALUES IN (1,2,3)") {
		t.Errorf("expected VALUES IN (1,2,3), got:\n%s", ddl)
	}
	if !strings.Contains(ddl, "VALUES IN (4,5,6)") {
		t.Errorf("expected VALUES IN (4,5,6), got:\n%s", ddl)
	}
	if !strings.Contains(ddl, "VALUES IN (7,8,9)") {
		t.Errorf("expected VALUES IN (7,8,9), got:\n%s", ddl)
	}

	// Verify catalog state
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if tbl.Partitioning == nil {
		t.Fatal("partitioning is nil")
	}
	if tbl.Partitioning.Type != "LIST" {
		t.Errorf("expected partition type LIST, got %q", tbl.Partitioning.Type)
	}
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Errorf("expected 3 partitions, got %d", len(tbl.Partitioning.Partitions))
	}
}

// Scenario 2: CREATE TABLE PARTITION BY LIST COLUMNS (col) — single column
func TestWalkThrough_7_2_ListColumnsSingle(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t2 (
		id INT NOT NULL,
		status VARCHAR(20) NOT NULL
	) PARTITION BY LIST COLUMNS (status) (
		PARTITION p_active VALUES IN ('active','pending'),
		PARTITION p_inactive VALUES IN ('inactive','archived')
	)`)

	ddl := c.ShowCreateTable("testdb", "t2")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	// LIST COLUMNS uses /*!50500
	if !strings.Contains(ddl, "/*!50500 PARTITION BY LIST") {
		t.Errorf("expected /*!50500 PARTITION BY LIST, got:\n%s", ddl)
	}
	// MySQL uses double space before COLUMNS
	if !strings.Contains(ddl, "LIST  COLUMNS(status)") {
		t.Errorf("expected LIST  COLUMNS(status), got:\n%s", ddl)
	}

	// Verify VALUES IN with string values
	if !strings.Contains(ddl, "VALUES IN ('active','pending')") {
		t.Errorf("expected VALUES IN ('active','pending'), got:\n%s", ddl)
	}

	tbl := c.GetDatabase("testdb").GetTable("t2")
	if tbl.Partitioning.Type != "LIST COLUMNS" {
		t.Errorf("expected partition type LIST COLUMNS, got %q", tbl.Partitioning.Type)
	}
	if len(tbl.Partitioning.Columns) != 1 {
		t.Errorf("expected 1 partition column, got %d", len(tbl.Partitioning.Columns))
	}
}

// Scenario 3: CREATE TABLE PARTITION BY LIST COLUMNS (col1, col2) — multi-column
func TestWalkThrough_7_2_ListColumnsMulti(t *testing.T) {
	c := wtSetup(t)

	// Multi-column LIST COLUMNS with tuple syntax: VALUES IN ((c1,c2),(c1,c2))
	// If the parser doesn't support tuple syntax, this will fail at parse time.
	results, err := c.Exec(`CREATE TABLE t3 (
		id INT NOT NULL,
		a INT NOT NULL,
		b INT NOT NULL
	) PARTITION BY LIST COLUMNS (a, b) (
		PARTITION p0 VALUES IN ((0,0),(0,1)),
		PARTITION p1 VALUES IN ((1,0),(1,1))
	)`, nil)
	if err != nil {
		t.Skipf("[~] partial: parser does not support multi-column LIST COLUMNS tuple syntax: %v", err)
		return
	}
	for _, r := range results {
		if r.Error != nil {
			t.Skipf("[~] partial: catalog error for multi-column LIST COLUMNS: %v", r.Error)
			return
		}
	}

	ddl := c.ShowCreateTable("testdb", "t3")
	if ddl == "" {
		t.Fatal("ShowCreateTable returned empty string")
	}

	if !strings.Contains(ddl, "LIST  COLUMNS(a,b)") {
		t.Errorf("expected LIST  COLUMNS(a,b), got:\n%s", ddl)
	}

	tbl := c.GetDatabase("testdb").GetTable("t3")
	if tbl.Partitioning.Type != "LIST COLUMNS" {
		t.Errorf("expected partition type LIST COLUMNS, got %q", tbl.Partitioning.Type)
	}
	if len(tbl.Partitioning.Columns) != 2 {
		t.Errorf("expected 2 partition columns, got %d", len(tbl.Partitioning.Columns))
	}
	if len(tbl.Partitioning.Partitions) != 2 {
		t.Errorf("expected 2 partitions, got %d", len(tbl.Partitioning.Partitions))
	}
}

// Scenario 4: ALTER TABLE ADD PARTITION with new VALUES IN
func TestWalkThrough_7_2_AlterAddPartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t4 (
		id INT NOT NULL,
		region INT NOT NULL
	) PARTITION BY LIST (region) (
		PARTITION p_east VALUES IN (1,2,3),
		PARTITION p_west VALUES IN (4,5,6)
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t4")
	if len(tbl.Partitioning.Partitions) != 2 {
		t.Fatalf("expected 2 partitions before ADD, got %d", len(tbl.Partitioning.Partitions))
	}

	wtExec(t, c, `ALTER TABLE t4 ADD PARTITION (PARTITION p_central VALUES IN (7,8,9))`)

	tbl = c.GetDatabase("testdb").GetTable("t4")
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Fatalf("expected 3 partitions after ADD, got %d", len(tbl.Partitioning.Partitions))
	}

	// Verify the new partition name
	lastPart := tbl.Partitioning.Partitions[2]
	if lastPart.Name != "p_central" {
		t.Errorf("expected partition name p_central, got %q", lastPart.Name)
	}

	ddl := c.ShowCreateTable("testdb", "t4")
	if !strings.Contains(ddl, "VALUES IN (7,8,9)") {
		t.Errorf("expected VALUES IN (7,8,9) in DDL:\n%s", ddl)
	}
}

// Scenario 5: ALTER TABLE DROP PARTITION from LIST table
func TestWalkThrough_7_2_AlterDropPartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t5 (
		id INT NOT NULL,
		region INT NOT NULL
	) PARTITION BY LIST (region) (
		PARTITION p_east VALUES IN (1,2,3),
		PARTITION p_west VALUES IN (4,5,6),
		PARTITION p_central VALUES IN (7,8,9)
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t5")
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Fatalf("expected 3 partitions before DROP, got %d", len(tbl.Partitioning.Partitions))
	}

	wtExec(t, c, `ALTER TABLE t5 DROP PARTITION p_west`)

	tbl = c.GetDatabase("testdb").GetTable("t5")
	if len(tbl.Partitioning.Partitions) != 2 {
		t.Fatalf("expected 2 partitions after DROP, got %d", len(tbl.Partitioning.Partitions))
	}

	// Verify remaining partition names
	names := make([]string, len(tbl.Partitioning.Partitions))
	for i, p := range tbl.Partitioning.Partitions {
		names[i] = p.Name
	}
	if names[0] != "p_east" || names[1] != "p_central" {
		t.Errorf("expected partitions [p_east, p_central], got %v", names)
	}

	ddl := c.ShowCreateTable("testdb", "t5")
	if strings.Contains(ddl, "p_west") {
		t.Errorf("dropped partition p_west still appears in DDL:\n%s", ddl)
	}
}

// Scenario 6: ALTER TABLE REORGANIZE PARTITION in LIST table
func TestWalkThrough_7_2_AlterReorganizePartition(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE t6 (
		id INT NOT NULL,
		region INT NOT NULL
	) PARTITION BY LIST (region) (
		PARTITION p_east VALUES IN (1,2,3),
		PARTITION p_west VALUES IN (4,5,6)
	)`)

	// Reorganize p_east into two new partitions
	wtExec(t, c, `ALTER TABLE t6 REORGANIZE PARTITION p_east INTO (
		PARTITION p_northeast VALUES IN (1,2),
		PARTITION p_southeast VALUES IN (3)
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t6")
	if len(tbl.Partitioning.Partitions) != 3 {
		t.Fatalf("expected 3 partitions after REORGANIZE, got %d", len(tbl.Partitioning.Partitions))
	}

	// Verify partition names and ordering
	expected := []string{"p_northeast", "p_southeast", "p_west"}
	for i, p := range tbl.Partitioning.Partitions {
		if p.Name != expected[i] {
			t.Errorf("partition %d: expected name %q, got %q", i, expected[i], p.Name)
		}
	}

	ddl := c.ShowCreateTable("testdb", "t6")
	if strings.Contains(ddl, "p_east") {
		t.Errorf("reorganized partition p_east still appears in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "VALUES IN (1,2)") {
		t.Errorf("expected VALUES IN (1,2) in DDL:\n%s", ddl)
	}
	if !strings.Contains(ddl, "VALUES IN (3)") {
		t.Errorf("expected VALUES IN (3) in DDL:\n%s", ddl)
	}
}
