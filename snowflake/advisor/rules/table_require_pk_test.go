package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestTableRequirePK(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "inline PK — no finding",
			sql:       "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(100))",
			wantCount: 0,
		},
		{
			name:      "table-level PK — no finding",
			sql:       "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100), PRIMARY KEY (id))",
			wantCount: 0,
		},
		{
			name:      "no PK — one finding",
			sql:       "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100))",
			wantCount: 1,
		},
		{
			name:      "TEMPORARY table — exempt",
			sql:       "CREATE TEMPORARY TABLE t (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "TRANSIENT table — exempt",
			sql:       "CREATE TRANSIENT TABLE t (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "VOLATILE table — exempt",
			sql:       "CREATE VOLATILE TABLE t (id INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "AS SELECT — exempt",
			sql:       "CREATE TABLE t AS SELECT id, name FROM other",
			wantCount: 0,
		},
		{
			name:      "LIKE — exempt",
			sql:       "CREATE TABLE t LIKE other_table",
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
			findings := runRule(t, tt.sql, rules.TableRequirePKRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestTableRequirePK_Metadata(t *testing.T) {
	r := rules.TableRequirePKRule{}
	if r.ID() != "snowflake.table.require-pk" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
