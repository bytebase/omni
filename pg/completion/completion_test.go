package completion

import (
	"testing"

	"github.com/bytebase/omni/pg/catalog"
)

func TestCompleteCandidateTypes(t *testing.T) {
	// Verify all candidate types are distinct.
	types := []CandidateType{
		CandidateKeyword, CandidateSchema, CandidateTable, CandidateView,
		CandidateMaterializedView, CandidateColumn, CandidateFunction,
		CandidateSequence, CandidateIndex, CandidateType_, CandidateTrigger,
		CandidatePolicy, CandidateExtension,
	}
	seen := make(map[CandidateType]bool)
	for _, ct := range types {
		if seen[ct] {
			t.Errorf("duplicate CandidateType value: %d", ct)
		}
		seen[ct] = true
	}
	if len(seen) != 13 {
		t.Errorf("expected 13 distinct types, got %d", len(seen))
	}
}

func TestCompleteEmpty(t *testing.T) {
	// Complete with empty input and no catalog should not panic.
	result := Complete("", 0, nil)
	// Result may be nil or empty for now; just verify no panic.
	_ = result
}

func TestCompleteTableFromCatalog(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	candidates := Complete("SELECT * FROM ", 14, cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateTable && c.Text == "users" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'users' table candidate")
	}
}

func TestCompleteKeywords(t *testing.T) {
	candidates := Complete("", 0, nil)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT keyword candidate")
	}
}

func TestCompletePrefixFilter(t *testing.T) {
	candidates := Complete("SEL", 3, nil)
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "INSERT" {
			t.Error("INSERT should be filtered out by prefix 'SEL'")
		}
	}
	found := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT to match prefix 'SEL'")
	}
}

func TestCompleteCTE(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	sql := "WITH active AS (SELECT * FROM users) SELECT * FROM "
	offset := len(sql)
	candidates := Complete(sql, offset, cat)

	// Should include CTE name "active" as a table ref
	found := false
	for _, c := range candidates {
		if c.Text == "active" {
			found = true
		}
	}
	if !found {
		t.Error("expected CTE 'active' in table candidates")
	}
}

func TestTrickyCompletePartialSQL(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE orders (id int, amount numeric);", nil)

	// "SELECT * FROM " — ends abruptly, standard may or may not work
	// but tricky should handle it
	candidates := Complete("SELECT * FROM ", 14, cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateTable && c.Text == "orders" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'orders' table in tricky completion")
	}
}
