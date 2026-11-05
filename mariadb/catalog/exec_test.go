package catalog

import "testing"

func TestExecSkipsDML(t *testing.T) {
	c := New()
	results, err := c.Exec("SELECT 1; INSERT INTO t VALUES (1)", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if !r.Skipped {
			t.Errorf("expected DML to be skipped")
		}
	}
}

func TestExecEmpty(t *testing.T) {
	c := New()
	results, err := c.Exec("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty SQL")
	}
}

// TestExecSkipsEmptyRows: INSERT and REPLACE with an empty column list or row are
// DML like any other.
func TestExecSkipsEmptyRows(t *testing.T) {
	c := New()
	results, err := c.Exec("INSERT INTO t () VALUES (); REPLACE INTO t VALUES (), (); INSERT INTO t () SELECT 1;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	for _, r := range results {
		if !r.Skipped || r.Error != nil {
			t.Errorf("statement %d: Skipped = %v, Error = %v", r.Index, r.Skipped, r.Error)
		}
	}
}
