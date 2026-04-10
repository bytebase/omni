package catalog

import "testing"

func TestCatalog(t *testing.T) {
	c := New()

	// Empty catalog.
	if len(c.Tables()) != 0 {
		t.Errorf("expected 0 tables, got %d", len(c.Tables()))
	}
	if c.HasTable("Music") {
		t.Error("HasTable returned true for non-existent table")
	}

	// Add tables.
	c.AddTable("Music")
	c.AddTable("Albums")
	c.AddTable("Music") // duplicate

	// Tables returns sorted, deduplicated.
	tables := c.Tables()
	if len(tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables))
	}
	if tables[0] != "Albums" || tables[1] != "Music" {
		t.Errorf("Tables() = %v, want [Albums Music]", tables)
	}

	// HasTable.
	if !c.HasTable("Music") {
		t.Error("HasTable returned false for existing table")
	}
	if c.HasTable("NonExistent") {
		t.Error("HasTable returned true for non-existent table")
	}
}
