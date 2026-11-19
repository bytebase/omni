package catalog

import (
	"errors"
	"testing"
)

// TestAlterRepartitionRefusedBeforeStoring holds the ALTER path to the checks
// CREATE TABLE makes: a repartitioning statement is judged before the catalog
// installs it, so a layout the engine would refuse leaves the table as it was.
func TestAlterRepartitionRefusedBeforeStoring(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE a (b INT PRIMARY KEY, c INT)")
	err := execError(t, c, "ALTER TABLE a PARTITION BY RANGE(b)(PARTITION p0 VALUES LESS THAN MAXVALUE, PARTITION p1 VALUES LESS THAN (10))")
	var refusal *Error
	if !errors.As(err, &refusal) || refusal.Code != 1481 {
		t.Fatalf("ALTER repartition: error %v, want 1481", err)
	}
	if tbl := c.GetDatabase("d").GetTable("a"); tbl.Partitioning != nil {
		t.Fatalf("the catalog installed the refused partitioning: %+v", tbl.Partitioning)
	}
	// Only RANGE takes no partition list as automatic; a LIST without one is
	// refused rather than stored as an empty layout.
	err = execError(t, c, "ALTER TABLE a PARTITION BY LIST(c)")
	if !errors.As(err, &refusal) || refusal.Code != 1492 {
		t.Fatalf("ALTER LIST without a partition list: error %v, want 1492", err)
	}
}

// TestAlterRepartitionJudgesTheTableItLeaves: a sibling command may drop the
// column the partitioning names. The engine refuses the whole statement at parse
// time, so there is no stored layout to compare; the catalog still must not leave
// a partitioning that names a column the statement dropped.
func TestAlterRepartitionJudgesTheTableItLeaves(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))")
	err := execError(t, c, "ALTER TABLE a DROP COLUMN a, PARTITION BY HASH(a)")
	var refusal *Error
	if !errors.As(err, &refusal) || refusal.Code != 1054 {
		t.Fatalf("ALTER with a dropped partition key: error %v, want 1054", err)
	}
	tbl := c.GetDatabase("d").GetTable("a")
	if tbl.GetColumn("a") == nil {
		t.Fatalf("the failed ALTER kept the column dropped: %+v", tbl.Columns)
	}
	if tbl.Partitioning != nil {
		t.Fatalf("the catalog installed a partitioning naming a dropped column: %+v", tbl.Partitioning)
	}
}
