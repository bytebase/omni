package catalog

import "testing"

// TestCurrentTimestampDefaultSpelling pins that the catalog records a
// NOW()/CURRENT_TIMESTAMP-family DEFAULT and ON UPDATE as TiDB's
// information_schema spells them (CURRENT_TIMESTAMP, with precision), so a
// column built from parsed DDL compares equal to one a live TiDB sync reads
// back. Bytebase's schema/tidb metadata path depends on this.
func TestCurrentTimestampDefaultSpelling(t *testing.T) {
	cases := []struct {
		ddl          string
		col          string
		wantDefault  string
		wantOnUpdate string
	}{
		{"CREATE TABLE t (ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP)", "ts", "CURRENT_TIMESTAMP", ""},
		{"CREATE TABLE t (ts TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3))", "ts", "CURRENT_TIMESTAMP(3)", ""},
		{"CREATE TABLE t (dt DATETIME DEFAULT CURRENT_TIMESTAMP)", "dt", "CURRENT_TIMESTAMP", ""},
		{"CREATE TABLE t (dt DATETIME DEFAULT NOW())", "dt", "CURRENT_TIMESTAMP", ""},
		{
			"CREATE TABLE t (ts TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6))",
			"ts", "CURRENT_TIMESTAMP(6)", "CURRENT_TIMESTAMP(6)",
		},
		// Constants and ordinary values are untouched.
		{"CREATE TABLE t (a INT DEFAULT 5)", "a", "5", ""},
		{"CREATE TABLE t (s VARCHAR(8) DEFAULT 'hi')", "s", "'hi'", ""},
	}
	for _, tc := range cases {
		c := New()
		mustOK(t, c, "CREATE DATABASE d")
		c.SetCurrentDatabase("d")
		mustOK(t, c, tc.ddl)
		col := c.GetDatabase("d").GetTable("t").GetColumn(tc.col)
		if col == nil {
			t.Fatalf("%s: column %s missing", tc.ddl, tc.col)
		}
		got := "<nil>"
		if col.Default != nil {
			got = *col.Default
		}
		if got != tc.wantDefault {
			t.Errorf("%s: Default=%q, want %q", tc.ddl, got, tc.wantDefault)
		}
		if col.OnUpdate != tc.wantOnUpdate {
			t.Errorf("%s: OnUpdate=%q, want %q", tc.ddl, col.OnUpdate, tc.wantOnUpdate)
		}
	}

	// The ALTER column builder records the same spelling as CREATE.
	c := New()
	mustOK(t, c, "CREATE DATABASE d")
	c.SetCurrentDatabase("d")
	mustOK(t, c, "CREATE TABLE t (id INT)")
	mustOK(t, c, "ALTER TABLE t ADD COLUMN ts TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)")
	col := c.GetDatabase("d").GetTable("t").GetColumn("ts")
	if col == nil || col.Default == nil || *col.Default != "CURRENT_TIMESTAMP(3)" {
		t.Errorf("ALTER ADD: Default = %v, want CURRENT_TIMESTAMP(3)", col)
	}
	if col != nil && col.OnUpdate != "CURRENT_TIMESTAMP(3)" {
		t.Errorf("ALTER ADD: OnUpdate = %q, want CURRENT_TIMESTAMP(3)", col.OnUpdate)
	}
}

func mustOK(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec(%q): %v", sql, r.Error)
		}
	}
}
