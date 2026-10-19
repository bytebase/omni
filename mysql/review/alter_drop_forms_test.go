package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestAlterDropForms checks the review of the ALTER TABLE ... DROP forms
// whose parse follows MySQL's alter_list_item: a column drop with RESTRICT or
// CASCADE is a column drop over its whole subcommand, and DROP PRIMARY or
// DROP FOREIGN without KEY, or IF EXISTS on a column, is a syntax error
// rather than a drop.
func TestAlterDropForms(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{{
		name: "a column drop with RESTRICT or CASCADE",
		sql:  "ALTER TABLE t DROP COLUMN n RESTRICT, DROP code CASCADE;",
		want: []finding{
			{0, "DROP COLUMN n RESTRICT", "drops column n of t", []int{0}},
			{0, "DROP code CASCADE", "drops column code of t", []int{0}},
		},
	}})
	runCases(t, review.Syntax, []ruleCase{
		{
			name: "DROP PRIMARY without KEY",
			sql:  "ALTER TABLE t DROP PRIMARY, DROP FOREIGN KEY fk_p;",
			want: []finding{{0, "", "unexpected token", []int{0}}},
		},
		{
			name: "DROP FOREIGN without KEY",
			sql:  "ALTER TABLE t DROP FOREIGN fk_p;",
			want: []finding{{0, "", "unexpected token", []int{0}}},
		},
		{
			name: "IF EXISTS on a column",
			sql:  "ALTER TABLE t DROP COLUMN IF EXISTS n;",
			want: []finding{{0, "", "expected identifier", []int{0}}},
		},
	})
}
