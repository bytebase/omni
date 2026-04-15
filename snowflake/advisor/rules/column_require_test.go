package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestColumnRequire(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		required  []string
		wantCount int
	}{
		{
			name:      "empty config — no findings",
			sql:       "CREATE TABLE t (id INT)",
			required:  nil,
			wantCount: 0,
		},
		{
			name:      "required column present — no findings",
			sql:       "CREATE TABLE t (id INT NOT NULL, created_at TIMESTAMP NOT NULL)",
			required:  []string{"id", "created_at"},
			wantCount: 0,
		},
		{
			name:      "required column missing — one finding",
			sql:       "CREATE TABLE t (id INT NOT NULL)",
			required:  []string{"id", "created_at"},
			wantCount: 1,
		},
		{
			name:      "all required missing — two findings",
			sql:       "CREATE TABLE t (name VARCHAR(100))",
			required:  []string{"id", "created_at"},
			wantCount: 2,
		},
		{
			name:      "case-insensitive match — no findings",
			sql:       "CREATE TABLE t (ID INT NOT NULL)",
			required:  []string{"id"},
			wantCount: 0,
		},
		{
			name:      "SELECT statement — no findings",
			sql:       "SELECT a FROM t",
			required:  []string{"id"},
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := rules.ColumnRequireConfig{Required: tt.required}
			rule := rules.NewColumnRequireRule(cfg)
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

func TestColumnRequire_Metadata(t *testing.T) {
	r := rules.NewColumnRequireRule(rules.ColumnRequireConfig{Required: []string{"id"}})
	if r.ID() != "snowflake.column.require" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
