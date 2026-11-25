package catalog

import (
	"errors"
	"testing"
)

// TestExecResults pins what Exec reports beside the schema it stores: one
// result per statement with its line, a stop at the first refusal unless
// ContinueOnError, and Trino's errorName on a refusal.
func TestExecResults(t *testing.T) {
	newMemory := func() *Catalog {
		c := New()
		c.EnsureCatalog("memory").EnsureSchema("s")
		c.SetCurrentCatalog("memory")
		c.SetCurrentSchema("s")
		return c
	}
	sql := "CREATE TABLE t (a integer);\nDROP TABLE nope;\nCREATE TABLE \"U\" (\"B\" bigint)"
	results, err := newMemory().Exec(sql, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].Error != nil || results[1].Line != 2 {
		t.Fatalf("got %+v, want two results, the second on line 2", results)
	}
	var e *Error
	if !errors.As(results[1].Error, &e) || e.Name != "TABLE_NOT_FOUND" {
		t.Fatalf("got %v, want TABLE_NOT_FOUND", results[1].Error)
	}

	c := newMemory()
	results, err = c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil || len(results) != 3 || results[2].Error != nil {
		t.Fatalf("got %+v, %v, want three results", results, err)
	}
	// Trino stores a quoted name in lower case too.
	u := c.GetCatalog("memory").GetSchema("s").GetTable("u")
	if u == nil || u.GetColumn("b") == nil {
		t.Fatalf("want table u with column b")
	}
	if got, want := u.ShowCreate("memory"), "CREATE TABLE memory.s.u (\n   b bigint\n)"; got != want {
		t.Errorf("ShowCreate = %q, want %q", got, want)
	}

	if _, err := newMemory().Exec("CREATE TABLE (", nil); err == nil {
		t.Errorf("want a parse error")
	}
}
