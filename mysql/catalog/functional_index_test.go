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

func TestFunctionalIndexCreateTableAutoNamesAndHiddenColumns(t *testing.T) {
	c := scenarioNewCatalog(t)
	results, err := c.Exec(`
		CREATE TABLE t (
			a INT,
			INDEX ((a + 1)),
			INDEX ((a * 2))
		);
	`, nil)
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
	if got := len(tbl.HiddenColumns()); got != 2 {
		t.Fatalf("hidden columns = %d, want 2", got)
	}

	wantIndexes := []struct {
		name      string
		hiddenCol string
		exprSub   string
	}{
		{"functional_index", "!hidden!functional_index!0!0", "`a` + 1"},
		{"functional_index_2", "!hidden!functional_index_2!0!0", "`a` * 2"},
	}
	for _, want := range wantIndexes {
		idx := findIndexForTest(tbl, want.name)
		if idx == nil {
			t.Fatalf("index %q missing", want.name)
		}
		if len(idx.Columns) != 1 {
			t.Fatalf("index %q columns = %d, want 1", want.name, len(idx.Columns))
		}
		if idx.Columns[0].Name != want.hiddenCol {
			t.Fatalf("index %q column name = %q, want %q", want.name, idx.Columns[0].Name, want.hiddenCol)
		}
		if !strings.Contains(idx.Columns[0].Expr, want.exprSub) {
			t.Fatalf("index %q expr = %q, want containing %q", want.name, idx.Columns[0].Expr, want.exprSub)
		}
		hidden := tbl.GetColumn(want.hiddenCol)
		if hidden == nil {
			t.Fatalf("hidden column %q missing", want.hiddenCol)
		}
		if hidden.Hidden != ColumnHiddenSystem {
			t.Fatalf("hidden column %q Hidden = %v, want ColumnHiddenSystem", want.hiddenCol, hidden.Hidden)
		}
		if hidden.Generated == nil || hidden.Generated.Stored {
			t.Fatalf("hidden column %q should be virtual generated column", want.hiddenCol)
		}
	}
}

func TestFunctionalIndexCreateTableNamedMultiPartHiddenColumns(t *testing.T) {
	c := scenarioNewCatalog(t)
	results, err := c.Exec("CREATE TABLE t (a INT, INDEX fx ((a + 1), (a * 2)));", nil)
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
	idx := findIndexForTest(tbl, "fx")
	if idx == nil {
		t.Fatal("index fx missing")
	}
	if len(idx.Columns) != 2 {
		t.Fatalf("index fx columns = %d, want 2", len(idx.Columns))
	}
	for i, wantName := range []string{"!hidden!fx!0!0", "!hidden!fx!1!0"} {
		if idx.Columns[i].Name != wantName {
			t.Fatalf("key part %d name = %q, want %q", i, idx.Columns[i].Name, wantName)
		}
		if col := tbl.GetColumn(wantName); col == nil || col.Hidden != ColumnHiddenSystem {
			t.Fatalf("hidden column %q missing or not system-hidden", wantName)
		}
	}
}

func findIndexForTest(tbl *Table, name string) *Index {
	for _, idx := range tbl.Indexes {
		if strings.EqualFold(idx.Name, name) {
			return idx
		}
	}
	return nil
}
