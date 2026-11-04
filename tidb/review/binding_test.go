package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughTiDBBinding pins what the walk does with a SQL plan
// binding: it follows it and records nothing. A statement the walk does not
// lower ends the walk before itself, so a binding in the middle of a change
// would leave every statement after it unjudged -- which is what the second
// case here would miss.
func TestWalkThroughTiDBBinding(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}

	// A binding itself is never a finding, whatever its inner statements
	// name: nothing in one runs, and the names are the application's.
	noFinding(t, "CREATE GLOBAL BINDING FOR SELECT * FROM nosuch USING SELECT * FROM nosuch;", shop)
	noFinding(t, "CREATE BINDING USING UPDATE nosuch SET a = 1;", shop)
	noFinding(t, "CREATE BINDING FOR DELETE FROM nosuch USING DELETE FROM nosuch;", shop)
	noFinding(t, "DROP BINDING FOR SQL DIGEST 'abc';", shop)
	noFinding(t, "CREATE BINDING FROM HISTORY USING PLAN DIGEST 'abc';", shop)
	noFinding(t, "SET BINDING DISABLED FOR SELECT * FROM nosuch;", shop)

	// The walk carries on past one, so the statement after it is judged.
	oneFinding(t, "CREATE BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(i);\nALTER TABLE t ADD COLUMN a int;", 1,
		"Column `a` already exists in `t`.", shop)
	oneFinding(t, "DROP BINDING FOR SELECT * FROM t;\nSET BINDING ENABLED FOR SELECT * FROM t;\nCREATE TABLE t (a int);", 2,
		"Table `t` already exists.", shop)

	// A binding carrying a CREATE-shaped name makes nothing, so the real
	// statement after it is the first to make it.
	noFinding(t, "CREATE BINDING USING INSERT INTO x SELECT 1;\nCREATE TABLE x (a int);", shop)
}

// TestRequireWhereBinding pins that REQUIRE_WHERE looks at a binding and
// not through it: the UPDATE or DELETE a binding carries updates and
// deletes nothing, where the same statement on its own is a finding.
func TestRequireWhereBinding(t *testing.T) {
	runCases(t, review.RequireWhere, []ruleCase{
		{
			name: "no finding inside a binding",
			sql: "CREATE GLOBAL BINDING FOR DELETE FROM t USING DELETE FROM t;\n" +
				"CREATE BINDING USING UPDATE t SET n = 1;\n" +
				"DROP BINDING FOR DELETE FROM t;\n" +
				"SET BINDING ENABLED FOR UPDATE t SET n = 1;",
		},
		{
			name: "the same statements outside one are findings",
			sql:  "CREATE BINDING FOR DELETE FROM t USING DELETE FROM t;\nDELETE FROM t;",
			want: []finding{{1, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}}},
		},
	})
}

// TestRequireIsNullBinding pins the other side of that choice: a predicate
// written inside a binding is still a predicate of a real query, so
// REQUIRE_IS_NULL reports it, as it reports one under a plain EXPLAIN.
// REQUIRE_WHERE is about what a statement destroys and a binding destroys
// nothing; this rule is about a comparison that can never be true, and one
// bound into a plan shape never will be either.
func TestRequireIsNullBinding(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "a comparison in a binding's own statements",
			sql:  "CREATE BINDING FOR SELECT * FROM t WHERE n = NULL USING SELECT * FROM t USE INDEX(idx_n) WHERE n = NULL;",
			want: []finding{
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			name: "and in DROP and SET",
			sql:  "DROP BINDING FOR SELECT * FROM t WHERE n = NULL;\nSET BINDING ENABLED FOR SELECT * FROM t WHERE n = NULL;",
			want: []finding{
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{1, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
	})
}
