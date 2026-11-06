package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// The rules still read a statement whose expressions carry the function
// argument forms this engine learned: an aggregate quantifier, EXTRACT's
// compound unit, the INTERVAL function and POSITION(a IN b).
//
// POSITION is a node kind of its own, and the walk reaches an unknown node by
// a default arm that abandons the statement — after which every later rule
// works from an Unknown table. So each case puts a finding on the far side of
// the new form: the rule has to walk past it to report.
func TestFunctionArgumentForms(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "a predicate beside each new form",
			sql: "SELECT * FROM t WHERE POSITION('x' IN code) > 0 AND n = NULL;\n" +
				"SELECT * FROM t WHERE EXTRACT(DAY_HOUR FROM n) > 0 AND n = NULL;\n" +
				"SELECT * FROM t WHERE INTERVAL(n, 1, 2) > 0 AND n = NULL;\n" +
				"SELECT * FROM t WHERE CAST(code AS SIGNED) > 0 AND n = NULL;\n" +
				"SELECT COUNT(DISTINCTROW n) FROM t WHERE n = NULL;\n" +
				"SELECT SUM(ALL n) FROM t WHERE n = NULL;",
			want: []finding{
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{1, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{2, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{3, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{4, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{5, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			name: "inside POSITION's own operands",
			sql:  "SELECT * FROM t WHERE POSITION('x' IN IF(n = NULL, code, code)) > 0;",
			want: []finding{{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
	})

	// A functional key part carries an expression where the walk reads the
	// table it creates.
	runCases(t, review.RequirePrimaryKey, []ruleCase{})
}
