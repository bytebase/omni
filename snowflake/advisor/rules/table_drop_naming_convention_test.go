package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestTableDropNamingConvention(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		pattern   string // empty means default
		wantCount int
	}{
		{
			name:      "normal table name — one finding (default pattern)",
			sql:       "DROP TABLE orders",
			wantCount: 1,
		},
		{
			name:      "deleted_ prefix — no finding",
			sql:       "DROP TABLE deleted_orders",
			wantCount: 0,
		},
		{
			name:      "_deleted name — no finding",
			sql:       "DROP TABLE _deleted",
			wantCount: 0,
		},
		{
			name:      "has _deleted in middle — one finding (does not match prefix/exact)",
			sql:       "DROP TABLE orders_deleted_2024",
			wantCount: 1,
		},
		{
			name:      "custom pattern — _archive suffix",
			sql:       "DROP TABLE orders_archive",
			pattern:   `_archive$`,
			wantCount: 0,
		},
		{
			name:      "custom pattern — does not match",
			sql:       "DROP TABLE orders",
			pattern:   `_archive$`,
			wantCount: 1,
		},
		{
			name:      "DROP VIEW — no finding (only DROP TABLE applies)",
			sql:       "DROP VIEW v",
			wantCount: 0,
		},
		{
			name:      "CREATE TABLE — no finding",
			sql:       "CREATE TABLE t (id INT)",
			wantCount: 0,
		},
		{
			name:      "SELECT — no finding",
			sql:       "SELECT a FROM t",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := rules.TableDropNamingConventionConfig{Pattern: tt.pattern}
			rule := rules.NewTableDropNamingConventionRule(cfg)
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

func TestTableDropNamingConvention_Metadata(t *testing.T) {
	r := rules.NewTableDropNamingConventionRule(rules.TableDropNamingConventionConfig{})
	if r.ID() != "snowflake.table.drop-naming-convention" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityWarning {
		t.Errorf("Severity() = %v, want WARNING", r.Severity())
	}
}
