package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestWhereRequireSelect(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "SELECT with WHERE — no finding",
			sql:       "SELECT a FROM t WHERE id = 1",
			wantCount: 0,
		},
		{
			name:      "SELECT without WHERE — one finding",
			sql:       "SELECT a FROM t",
			wantCount: 1,
		},
		{
			name:      "SELECT without FROM (table-less) — no finding",
			sql:       "SELECT 1",
			wantCount: 0,
		},
		{
			name:      "SELECT CURRENT_TIMESTAMP() — no finding",
			sql:       "SELECT CURRENT_TIMESTAMP()",
			wantCount: 0,
		},
		{
			name:      "UNION — two findings (each branch lacks WHERE)",
			sql:       "SELECT a FROM t1 UNION ALL SELECT b FROM t2",
			wantCount: 2,
		},
		{
			name:      "UNION — one branch has WHERE — one finding",
			sql:       "SELECT a FROM t1 WHERE id = 1 UNION ALL SELECT b FROM t2",
			wantCount: 1,
		},
		{
			name:      "CREATE TABLE — no finding",
			sql:       "CREATE TABLE t (id INT)",
			wantCount: 0,
		},
		{
			name:      "SELECT with GROUP BY but no WHERE — one finding",
			sql:       "SELECT dept, COUNT(*) FROM employees GROUP BY dept",
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.WhereRequireSelectRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestWhereRequireSelect_Metadata(t *testing.T) {
	r := rules.WhereRequireSelectRule{}
	if r.ID() != "snowflake.query.where-require-select" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityWarning {
		t.Errorf("Severity() = %v, want WARNING", r.Severity())
	}
}
