package completion

import "testing"

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
