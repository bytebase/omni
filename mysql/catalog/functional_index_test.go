package catalog

import (
	"strings"
	"testing"
)

func TestFunctionalIndexHiddenColumnsVisibleAndHiddenViews(t *testing.T) {
	c := scenarioNewCatalog(t)
	results, err := c.Exec("CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));", nil)
	if err != nil {
		t.Fatalf("exec create table: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec create table result error: %v", r.Error)
		}
	}

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t missing")
	}
	tbl.Columns = append(tbl.Columns, &Column{
		Position:   3,
		Name:       "!hidden!idx_lower!0!0",
		DataType:   "varchar",
		ColumnType: "varchar(64)",
		Hidden:     ColumnHiddenSystem,
		Generated:  &GeneratedColumnInfo{Expr: "lower(`name`)", Stored: false},
	})
	rebuildColIndex(tbl)

	if got := len(tbl.VisibleColumns()); got != 2 {
		t.Fatalf("visible columns = %d, want 2", got)
	}
	if got := len(tbl.HiddenColumns()); got != 1 {
		t.Fatalf("hidden columns = %d, want 1", got)
	}

	show := c.ShowCreateTable("testdb", "t")
	if strings.Contains(show, "!hidden!idx_lower!0!0") {
		t.Fatalf("ShowCreateTable leaked system-hidden column:\n%s", show)
	}
}
