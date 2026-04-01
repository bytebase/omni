package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/mongo/analysis"
)

func TestOperationIsRead(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want bool
	}{
		{analysis.OpFind, true},
		{analysis.OpFindOne, true},
		{analysis.OpAggregate, true},
		{analysis.OpCount, true},
		{analysis.OpDistinct, true},
		{analysis.OpRead, true},
		{analysis.OpWrite, false},
		{analysis.OpAdmin, false},
		{analysis.OpInfo, false},
		{analysis.OpExplain, false},
		{analysis.OpUnknown, false},
	}
	for _, tc := range tests {
		if got := tc.op.IsRead(); got != tc.want {
			t.Errorf("%v.IsRead() = %v, want %v", tc.op, got, tc.want)
		}
	}
}

func TestOperationString(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want string
	}{
		{analysis.OpFind, "find"},
		{analysis.OpFindOne, "findOne"},
		{analysis.OpAggregate, "aggregate"},
		{analysis.OpCount, "count"},
		{analysis.OpDistinct, "distinct"},
		{analysis.OpRead, "read"},
		{analysis.OpWrite, "write"},
		{analysis.OpAdmin, "admin"},
		{analysis.OpInfo, "info"},
		{analysis.OpExplain, "explain"},
		{analysis.OpUnknown, "unknown"},
	}
	for _, tc := range tests {
		if got := tc.op.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.op, got, tc.want)
		}
	}
}
