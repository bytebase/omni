package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestAlterDropForms checks the review of the ALTER TABLE ... DROP forms
// whose parse follows TiDB's AlterTableSpec: a column drop with RESTRICT or
// CASCADE, with or without IF EXISTS, is a column drop over its whole
// subcommand, and DROP PRIMARY or DROP FOREIGN without KEY is a syntax error
// rather than a drop.
func TestAlterDropForms(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{{
		name: "a column drop with RESTRICT or CASCADE",
		sql:  "ALTER TABLE t DROP COLUMN n RESTRICT, DROP IF EXISTS code CASCADE;",
		want: []finding{
			{0, "DROP COLUMN n RESTRICT", "drops column n of t", []int{0}},
			{0, "DROP IF EXISTS code CASCADE", "drops column code of t", []int{0}},
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
	})
}
