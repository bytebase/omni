package pg

import (
	"testing"

	"github.com/bytebase/omni/pg/parser"
)

func TestCollectCompletionWrapper(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN t2 USING ("
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Candidates == nil || ctx.Scope == nil {
		t.Fatal("expected completion context")
	}
	if !ctx.Candidates.HasRule("columnref") {
		t.Fatal("expected columnref candidate")
	}
	if len(ctx.Scope.References) != 2 {
		t.Fatalf("reference count = %d, want 2", len(ctx.Scope.References))
	}
	if ctx.Scope.References[0].Kind != parser.RangeReferenceRelation {
		t.Fatalf("first reference kind = %v, want relation", ctx.Scope.References[0].Kind)
	}
}
