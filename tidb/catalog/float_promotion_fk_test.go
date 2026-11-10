package catalog

import "testing"

// TestFloatPromotionFKCompatibility checks that a FLOAT(M>24) column added
// through ALTER gets the same promoted DOUBLE base type as one created through
// CREATE TABLE, so a foreign key between them is accepted (both DOUBLE) instead
// of rejected with 3780 for mismatched base types. Regression test for the
// CREATE/ALTER promotion inconsistency.
func TestFloatPromotionFKCompatibility(t *testing.T) {
	c := New()
	exec := func(sql string) {
		t.Helper()
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("Exec(%q) parse error: %v", sql, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("Exec(%q) error: %v", sql, r.Error)
			}
		}
	}
	exec("CREATE DATABASE d")
	c.SetCurrentDatabase("d")
	exec("CREATE TABLE p (a FLOAT(30), KEY (a))")
	exec("CREATE TABLE ch (id INT)")
	// FLOAT(30) added via ALTER must promote to DOUBLE, like the CREATE path.
	exec("ALTER TABLE ch ADD COLUMN a FLOAT(30)")
	// The FK references p(a); both columns are DOUBLE, so it must be accepted.
	exec("ALTER TABLE ch ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)")

	for _, tn := range []string{"p", "ch"} {
		tbl := c.GetDatabase("d").GetTable(tn)
		if tbl == nil {
			t.Fatalf("table %s missing", tn)
		}
		col := tbl.GetColumn("a")
		if col == nil {
			t.Fatalf("%s.a missing", tn)
		}
		if col.DataType != "double" {
			t.Errorf("%s.a DataType = %q, want double", tn, col.DataType)
		}
	}
}
