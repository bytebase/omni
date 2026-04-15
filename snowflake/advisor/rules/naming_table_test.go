package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestNamingTable(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		pattern   string // empty means default
		wantCount int
	}{
		{
			name:      "valid snake_case — no finding",
			sql:       "CREATE TABLE my_table (id INT)",
			wantCount: 0,
		},
		{
			name:      "uppercase name — one finding (default pattern)",
			sql:       "CREATE TABLE MyTable (id INT)",
			wantCount: 1,
		},
		{
			name:      "underscore start (not matching default pattern) — one finding",
			sql:       "CREATE TABLE _bad (id INT)",
			wantCount: 1,
		},
		{
			name:      "custom pattern uppercase — valid",
			sql:       "CREATE TABLE MY_TABLE (id INT)",
			pattern:   `^[A-Z][A-Z0-9_]+$`,
			wantCount: 0,
		},
		{
			name:      "custom pattern uppercase — invalid",
			sql:       "CREATE TABLE my_table (id INT)",
			pattern:   `^[A-Z][A-Z0-9_]+$`,
			wantCount: 1,
		},
		{
			name:      "valid with numbers",
			sql:       "CREATE TABLE order2 (id INT)",
			wantCount: 0,
		},
		{
			name:      "SELECT statement — no finding",
			sql:       "SELECT a FROM t",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := rules.NamingTableConfig{Pattern: tt.pattern}
			rule := rules.NewNamingTableRule(cfg)
			findings := runRule(t, tt.sql, rule)
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestNamingTable_Metadata(t *testing.T) {
	r := rules.NewNamingTableRule(rules.NamingTableConfig{})
	if r.ID() != "snowflake.naming.table" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityWarning {
		t.Errorf("Severity() = %v, want WARNING", r.Severity())
	}
}
