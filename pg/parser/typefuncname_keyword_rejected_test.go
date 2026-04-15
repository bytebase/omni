package parser

import "testing"

// TestTypeFuncNameKeywordRejectedAsExpression locks in that
// TypeFuncNameKeyword tokens (INNER, LEFT, JOIN, CROSS, NATURAL, etc.)
// are NOT valid expression atoms in omni's parser. They never were
// accepted as bare expressions — the previous `isColId() ||
// isTypeFunctionName()` predicates at expr.go:1475,
// create_table.go:901, and create_index.go:273 admitted them only to
// have parseColId reject them one call deeper. After removing those
// dead-code ORs, this test confirms the behavior is unchanged.
//
// PG's grammar:  a_expr → c_expr → columnref → ColId
// (TypeFuncNameKeyword is NOT in ColId, so PG also rejects.)
func TestTypeFuncNameKeywordRejectedAsExpression(t *testing.T) {
	sqls := []string{
		`SELECT inner`,
		`SELECT left`,
		`SELECT right`,
		`SELECT join`,
		`SELECT cross`,
		`SELECT outer`,
		`SELECT full`,
		`SELECT natural`,
		`SELECT verbose`,
	}
	for _, sql := range sqls {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			if err == nil {
				t.Errorf("expected parse error for %q (TypeFuncNameKeyword as expression atom), got nil", sql)
			}
		})
	}
}
