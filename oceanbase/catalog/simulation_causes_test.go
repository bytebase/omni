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
}
