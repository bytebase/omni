package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestWhereRequireUpdateDelete(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "UPDATE with WHERE — no finding",
			sql:       "UPDATE t SET name = 'x' WHERE id = 1",
			wantCount: 0,
		},
		{
			name:      "UPDATE without WHERE — one finding",
			sql:       "UPDATE t SET name = 'x'",
			wantCount: 1,
		},
		{
			name:      "DELETE with WHERE — no finding",
			sql:       "DELETE FROM t WHERE id = 1",
			wantCount: 0,
		},
		{
			name:      "DELETE without WHERE — one finding",
			sql:       "DELETE FROM t",
			wantCount: 1,
		},
		{
			name:      "SELECT — no finding",
			sql:       "SELECT a FROM t",
			wantCount: 0,
		},
		{
			name:      "CREATE TABLE — no finding",
			sql:       "CREATE TABLE t (id INT)",
			wantCount: 0,
		},
		{
			name:      "UPDATE with complex WHERE — no finding",
			sql:       "UPDATE t SET status = 'active' WHERE created_at > '2024-01-01' AND status = 'pending'",
			wantCount: 0,
		},
		{
			name:      "DELETE with subquery WHERE — no finding",
			sql:       "DELETE FROM t WHERE id IN (SELECT id FROM other WHERE flag = true)",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.WhereRequireUpdateDeleteRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestWhereRequireUpdateDelete_Metadata(t *testing.T) {
	r := rules.WhereRequireUpdateDeleteRule{}
	if r.ID() != "snowflake.query.where-require-update-delete" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}

func TestWhereRequireUpdateDelete_FindingMessages(t *testing.T) {
	tests := []struct {
		sql     string
		wantMsg string
	}{
		{
			sql:     "UPDATE t SET a = 1",
			wantMsg: "UPDATE without a WHERE clause will modify all rows",
		},
		{
			sql:     "DELETE FROM t",
			wantMsg: "DELETE without a WHERE clause will remove all rows",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.WhereRequireUpdateDeleteRule{})
			if len(findings) != 1 {
				t.Fatalf("expected 1 finding, got %d", len(findings))
			}
			if findings[0].Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", findings[0].Message, tt.wantMsg)
			}
		})
	}
}
