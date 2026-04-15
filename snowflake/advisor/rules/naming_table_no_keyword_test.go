package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestNamingTableNoKeyword(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "normal name — no finding",
			sql:       "CREATE TABLE orders (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "reserved keyword SELECT — one finding",
			sql:       `CREATE TABLE select (id INT NOT NULL)`,
			wantCount: 1,
		},
		{
			name:      "reserved keyword TABLE — one finding",
			sql:       `CREATE TABLE table (id INT NOT NULL)`,
			wantCount: 1,
		},
		{
			name:      "quoted reserved keyword — no finding (safe)",
			sql:       `CREATE TABLE "select" (id INT NOT NULL)`,
			wantCount: 0,
		},
		{
			name:      "non-reserved keyword (ABORT is not reserved) — no finding",
			sql:       "CREATE TABLE abort (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "reserved keyword WHERE — one finding",
			sql:       `CREATE TABLE where (id INT NOT NULL)`,
			wantCount: 1,
		},
		{
			name:      "SELECT statement — no finding",
			sql:       "SELECT a FROM t",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.NamingTableNoKeywordRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestNamingTableNoKeyword_Metadata(t *testing.T) {
	r := rules.NamingTableNoKeywordRule{}
	if r.ID() != "snowflake.naming.table-no-keyword" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
