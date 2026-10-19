package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestDisallowDropObjectDropPeriodIfExists: DROP PERIOD IF EXISTS is still a
// DROP PERIOD; IF EXISTS does not exempt it.
func TestDisallowDropObjectDropPeriodIfExists(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "IF EXISTS does not exempt DROP PERIOD",
			sql:  "ALTER TABLE t DROP PERIOD IF EXISTS FOR app;\nALTER TABLE t DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`, DROP COLUMN n;",
			want: []finding{
				{0, "DROP PERIOD IF EXISTS FOR app", "drops period app of t", []int{0}},
				{1, "DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", "drops period SYSTEM_TIME of t", []int{0}},
				{1, "DROP COLUMN n", "drops column n of t", []int{0}},
			},
		},
	})
}
