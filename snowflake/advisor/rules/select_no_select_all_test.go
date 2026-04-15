package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestSelectNoSelectAll(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "explicit column list — no findings",
			sql:       "SELECT a, b, c FROM t",
			wantCount: 0,
		},
		{
			name:      "bare star — one finding",
			sql:       "SELECT * FROM t",
			wantCount: 1,
		},
		{
			name:      "qualified star — no finding",
			sql:       "SELECT t.* FROM t",
			wantCount: 0,
		},
		{
			name:      "UNION with two bare stars — two findings",
			sql:       "SELECT * FROM a UNION ALL SELECT * FROM b",
			wantCount: 2,
		},
		{
			name:      "mix of star and column — one finding",
			sql:       "SELECT *, a FROM t",
			wantCount: 1,
		},
		{
			name:      "CREATE TABLE — no finding",
			sql:       "CREATE TABLE t (id INT)",
			wantCount: 0,
		},
		{
			name:      "SELECT 1 (no table) — no finding",
			sql:       "SELECT 1",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.SelectNoSelectAllRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestSelectNoSelectAll_Metadata(t *testing.T) {
	r := rules.SelectNoSelectAllRule{}
	if r.ID() != "snowflake.select.no-select-all" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}

func TestSelectNoSelectAll_SeverityIsError(t *testing.T) {
	findings := runRule(t, "SELECT * FROM t", rules.SelectNoSelectAllRule{})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding")
	}
	if findings[0].Severity != advisor.SeverityError {
		t.Errorf("Severity = %v, want ERROR", findings[0].Severity)
	}
}
