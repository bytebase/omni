package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestMigrationCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "CREATE TABLE — no finding",
			sql:       "CREATE TABLE t (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "ALTER TABLE ADD COLUMN — no finding",
			sql:       "ALTER TABLE t ADD COLUMN new_col INT",
			wantCount: 0,
		},
		{
			name:      "DROP TABLE — one finding",
			sql:       "DROP TABLE orders",
			wantCount: 1,
		},
		{
			name:      "DROP TABLE IF EXISTS — one finding",
			sql:       "DROP TABLE IF EXISTS orders",
			wantCount: 1,
		},
		{
			name:      "DROP VIEW — no finding (only DROP TABLE is destructive by this rule)",
			sql:       "DROP VIEW v",
			wantCount: 0,
		},
		{
			name:      "ALTER TABLE DROP COLUMN — one finding",
			sql:       "ALTER TABLE t DROP COLUMN old_col",
			wantCount: 1,
		},
		{
			name:      "ALTER TABLE DROP COLUMN multiple — multiple findings",
			sql:       "ALTER TABLE t DROP COLUMN a, DROP COLUMN b",
			wantCount: 2,
		},
		{
			name:      "ALTER TABLE RENAME — no finding",
			sql:       "ALTER TABLE t RENAME TO t2",
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
			findings := runRule(t, tt.sql, rules.MigrationCompatibilityRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestMigrationCompatibility_Metadata(t *testing.T) {
	r := rules.MigrationCompatibilityRule{}
	if r.ID() != "snowflake.migration.compatibility" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
