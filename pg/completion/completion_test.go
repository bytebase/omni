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
