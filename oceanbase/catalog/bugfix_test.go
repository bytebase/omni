package catalog

import (
	"testing"
)

func TestBugFix_InExprNodeToSQL(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (category VARCHAR(50), CONSTRAINT chk CHECK (category IN ('a','b','c')))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	var found bool
	for _, con := range tbl.Constraints {
		if con.Name == "chk" {
			found = true
			// The expression should contain IN, not "(?)"
			if con.CheckExpr == "(?)" || con.CheckExpr == "" {
				t.Errorf("CHECK expression was not properly deparsed: got %q", con.CheckExpr)
			}
			t.Logf("CHECK expression: %s", con.CheckExpr)
		}
	}
	if !found {
		t.Error("constraint chk not found")
	}
}

func TestBugFix_FKIndexOrder(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, `CREATE TABLE child (
		id INT AUTO_INCREMENT PRIMARY KEY,
		parent_id INT NOT NULL,
		CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id),
		INDEX idx_parent (parent_id)
	)`)
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	// Should have 2 indexes: PRIMARY + idx_parent (FK should reuse idx_parent, not create a 3rd)
	if len(tbl.Indexes) != 2 {
		names := make([]string, len(tbl.Indexes))
		for i, idx := range tbl.Indexes {
			names[i] = idx.Name
		}
		t.Errorf("expected 2 indexes (PRIMARY + idx_parent), got %d: %v", len(tbl.Indexes), names)
	}
}

func TestBugFix_PartitionAutoGen(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	if tbl.Partitioning == nil {
		t.Fatal("partitioning info is nil")
	}
	if len(tbl.Partitioning.Partitions) != 4 {
		t.Errorf("expected 4 partitions, got %d", len(tbl.Partitioning.Partitions))
	}
	for i, p := range tbl.Partitioning.Partitions {
		expected := "p" + string(rune('0'+i))
		if p.Name != expected {
			t.Errorf("partition %d: expected name %q, got %q", i, expected, p.Name)
		}
	}
}
