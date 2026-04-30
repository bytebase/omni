package mssql

import "testing"

func TestCollectCompletionWrapper(t *testing.T) {
	sql := "SELECT v. FROM (VALUES (1, 'a')) AS v(Id, Label)"
	ctx := CollectCompletion(sql, len("SELECT v."))
	if ctx == nil || ctx.Candidates == nil || ctx.Scope == nil {
		t.Fatal("expected completion context")
	}
	if len(ctx.Scope.LocalReferences) != 1 {
		t.Fatalf("local reference count = %d, want 1", len(ctx.Scope.LocalReferences))
	}
	if ctx.Scope.LocalReferences[0].Kind != RangeReferenceValues {
		t.Fatalf("reference kind = %v, want values", ctx.Scope.LocalReferences[0].Kind)
	}
}
