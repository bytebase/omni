package catalog

import (
	"strings"
	"testing"
)

// TestAlterTableSpaceSeparatedOptions: table options that follow one another
// without commas in ALTER TABLE each apply, as they do separated by commas.
func TestAlterTableSpaceSeparatedOptions(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (a INT); ALTER TABLE t ROW_FORMAT = DYNAMIC COMMENT 'x' AUTO_INCREMENT = 5;", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec: %v", r.Error)
		}
	}
	// TiDB v8.5.5 keeps ROW_FORMAT without writing it in SHOW CREATE TABLE
	// (read on 2026-11-06), so the catalog's field is where it shows.
	tbl := c.GetDatabase("d").GetTable("t")
	if tbl.RowFormat != "DYNAMIC" || tbl.Comment != "x" || tbl.AutoIncrement != 5 {
		t.Errorf("options: ROW_FORMAT=%q COMMENT=%q AUTO_INCREMENT=%d", tbl.RowFormat, tbl.Comment, tbl.AutoIncrement)
	}
	ddl := c.ShowCreateTable("d", "t")
	if !strings.Contains(ddl, "COMMENT='x'") || strings.Contains(ddl, "ROW_FORMAT") {
		t.Errorf("SHOW CREATE TABLE:\n%s", ddl)
	}
}
