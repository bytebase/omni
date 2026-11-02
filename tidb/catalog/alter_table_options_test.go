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
	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"ROW_FORMAT=DYNAMIC", "COMMENT='x'"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}
