package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestColumnNoNull(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "all NOT NULL — no findings",
			sql:       "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100) NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "one nullable column — one finding",
			sql:       "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100))",
			wantCount: 1,
		},
		{
			name:      "two nullable columns — two findings",
			sql:       "CREATE TABLE t (id INT, name VARCHAR(100))",
			wantCount: 2,
		},
		{
			name:      "explicit NULL — one finding",
			sql:       "CREATE TABLE t (id INT NULL)",
			wantCount: 1,
		},
		{
			name:      "virtual column exempt",
			sql:       "CREATE TABLE t (id INT NOT NULL, full_name VARCHAR(200) AS (id::VARCHAR))",
			wantCount: 0,
		},
		{
			name:      "mix of virtual and nullable — only real cols fire",
			sql:       "CREATE TABLE t (id INT, full_name VARCHAR(200) AS (id::VARCHAR))",
			wantCount: 1,
		},
		{
			name:      "SELECT statement — no findings",
			sql:       "SELECT a, b FROM t",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.ColumnNoNullRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestColumnNoNull_Metadata(t *testing.T) {
	r := rules.ColumnNoNullRule{}
	if r.ID() != "snowflake.column.no-null" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}

func TestColumnNoNull_FindingHasCorrectRuleID(t *testing.T) {
	findings := runRule(t, "CREATE TABLE t (id INT)", rules.ColumnNoNullRule{})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding")
	}
	if findings[0].RuleID != "snowflake.column.no-null" {
		t.Errorf("RuleID = %q", findings[0].RuleID)
	}
}
