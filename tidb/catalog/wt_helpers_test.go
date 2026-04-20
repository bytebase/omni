package catalog

import "testing"

// wtSetup creates a Catalog with database "testdb" selected, ready for walk-through tests.
func wtSetup(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE testdb; USE testdb;", nil)
	if err != nil {
		t.Fatalf("wtSetup parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("wtSetup exec error: %v", r.Error)
		}
	}
	return c
}

// wtExec executes SQL on the catalog and fatals on any error.
func wtExec(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("wtExec parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("wtExec exec error on stmt %d: %v", r.Index, r.Error)
		}
	}
}

// assertError asserts that err is a *catalog.Error with the given MySQL error code.
func assertError(t *testing.T, err error, code int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %d, got nil", code)
	}
	catErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *catalog.Error, got %T: %v", err, err)
	}
	if catErr.Code != code {
		t.Errorf("expected error code %d, got %d: %s", code, catErr.Code, catErr.Message)
	}
}

// assertNoError asserts that err is nil.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
