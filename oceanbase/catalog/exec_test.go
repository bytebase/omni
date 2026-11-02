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
