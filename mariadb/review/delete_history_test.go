package review

import (
	"context"
	"testing"

	"github.com/bytebase/omni/review"
)

// TestRequireWhereSkipsDeleteHistory: DELETE HISTORY takes no WHERE, since
// BEFORE SYSTEM_TIME bounds it, so RequireWhere has nothing to ask of it.
func TestRequireWhereSkipsDeleteHistory(t *testing.T) {
	sql := "DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP '2020-01-01 00:00:00';"
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.RequireWhere}}, []review.Target{{Schema: shop()}})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Errorf("findings = %+v, want none", result.Findings)
	}
}
