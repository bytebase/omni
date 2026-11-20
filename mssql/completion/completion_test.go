package completion

import (
	"strings"
	"testing"
)

func TestCompleteCandidateTypes(t *testing.T) {
	// Verify all candidate types are distinct.
	types := []CandidateType{
		CandidateKeyword, CandidateSchema, CandidateTable, CandidateView,
		CandidateColumn, CandidateFunction, CandidateProcedure, CandidateIndex,
		CandidateTrigger, CandidateSequence, CandidateType_, CandidateLogin,
		CandidateUser, CandidateRole, CandidateDatabase,
	}
	seen := make(map[CandidateType]bool)
	for _, ct := range types {
		if seen[ct] {
			t.Errorf("duplicate CandidateType value: %d", ct)
		}
		seen[ct] = true
	}
	if len(seen) != 15 {
		t.Errorf("expected 15 distinct types, got %d", len(seen))
	}
}

func TestCandidateStructFields(t *testing.T) {
	c := Candidate{
		Text:       "SELECT",
		Type:       CandidateKeyword,
		Definition: "SQL keyword",
		Comment:    "Retrieves rows",
	}
	if c.Text != "SELECT" || c.Type != CandidateKeyword || c.Definition != "SQL keyword" || c.Comment != "Retrieves rows" {
		t.Error("Candidate struct fields not set correctly")
	}
}

func TestCompleteSignature(t *testing.T) {
	// Complete(sql, cursorOffset, catalog) returns []Candidate
	result := Complete("", 0, nil)
	// Should return keyword candidates even with nil catalog.
	if result == nil {
		t.Fatal("expected non-nil result for empty input")
	}
	var _ []Candidate = result
}

func TestCompleteNilCatalogReturnsKeywordsOnly(t *testing.T) {
	result := Complete("", 0, nil)
	for _, c := range result {
		if c.Type != CandidateKeyword {
			t.Errorf("expected only keyword candidates with nil catalog, got type %d for %q", c.Type, c.Text)
		}
	}
}

func TestCompleteEmptyReturnsTopLevelKeywords(t *testing.T) {
	result := Complete("", 0, nil)
	if len(result) == 0 {
		t.Fatal("expected top-level keywords for empty input")
	}
	// Should contain SELECT
	found := false
	for _, c := range result {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT keyword candidate for empty input")
	}
}

func TestCompletePrefixFiltering(t *testing.T) {
	// "SEL" at offset 3 should match SELECT but not INSERT.
	result := Complete("SEL", 3, nil)
	foundSelect := false
	for _, c := range result {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			foundSelect = true
		}
		if c.Type == CandidateKeyword && c.Text == "INSERT" {
			t.Error("INSERT should be filtered out by prefix 'SEL'")
		}
	}
	if !foundSelect {
		t.Error("expected SELECT to match prefix 'SEL'")
	}
}

func TestCompletePrefixCaseInsensitive(t *testing.T) {
	// Lowercase prefix should still match uppercase keywords.
	result := Complete("sel", 3, nil)
	found := false
	for _, c := range result {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT to match lowercase prefix 'sel'")
	}
}

func TestCompleteDeduplication(t *testing.T) {
	result := Complete("", 0, nil)
	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool)
	for _, c := range result {
		k := key{strings.ToLower(c.Text), c.Type}
		if seen[k] {
			t.Errorf("duplicate candidate: %q (type %d)", c.Text, c.Type)
		}
		seen[k] = true
	}
}

// --- Tricky completion tests ---

func TestTrickyCompleteTrailingSpace(t *testing.T) {
	// "SELECT * FROM " with trailing space should produce candidates via placeholder.
	sql := "SELECT * FROM "
	result := Complete(sql, len(sql), nil)
	if len(result) == 0 {
		t.Fatal("expected candidates for 'SELECT * FROM ' (trailing space)")
	}
}

func TestTrickyCompleteTruncatedKeyword(t *testing.T) {
	// "SELE" at offset 4 should prefix-filter against keywords and include SELECT.
	result := Complete("SELE", 4, nil)
	found := false
	for _, c := range result {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT to match truncated keyword 'SELE'")
	}
}

func TestTrickyCompleteAfterComma(t *testing.T) {
	// "SELECT a," at end should produce candidates (column or keyword).
	sql := "SELECT a,"
	result := Complete(sql, len(sql), nil)
	if len(result) == 0 {
		t.Fatal("expected candidates after comma in 'SELECT a,'")
	}
}

func TestTrickyCompleteAfterOperator(t *testing.T) {
	// "SELECT * FROM t WHERE a >" at end should produce candidates.
	sql := "SELECT * FROM t WHERE a >"
	result := Complete(sql, len(sql), nil)
	if len(result) == 0 {
		t.Fatal("expected candidates after operator in 'WHERE a >'")
	}
}

func TestTrickyCompleteMultipleStrategies(t *testing.T) {
	// Ensure trickyComplete tries strategies in order and returns first non-empty.
	// We test indirectly: if standardComplete fails but trickyComplete succeeds,
	// Complete should still return results.
	sql := "INSERT INTO "
	result := Complete(sql, len(sql), nil)
	if len(result) == 0 {
		t.Fatal("expected candidates for 'INSERT INTO ' (trailing space)")
	}
}

func TestTrickyCompleteFallbackBestEffort(t *testing.T) {
	// Even with a difficult pattern, we should get best-effort results or nil.
	// This just exercises the fallback path — no crash.
	sql := "EXEC sp_"
	result := Complete(sql, len(sql), nil)
	// May or may not return results, but should not panic.
	_ = result
}

func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		sql    string
		offset int
		want   string
	}{
		{"SELECT", 3, "SEL"},
		{"SELECT ", 7, ""},
		{"SELECT foo", 10, "foo"},
		{"", 0, ""},
		{"SELECT", 100, "SELECT"}, // offset beyond length
	}
	for _, tt := range tests {
		got := extractPrefix(tt.sql, tt.offset)
		if got != tt.want {
			t.Errorf("extractPrefix(%q, %d) = %q, want %q", tt.sql, tt.offset, got, tt.want)
		}
	}
}
