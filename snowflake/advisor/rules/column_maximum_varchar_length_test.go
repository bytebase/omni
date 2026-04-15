package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestColumnMaximumVarcharLength(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		maxLen    int // 0 means use default
		wantCount int
	}{
		{
			name:      "varchar within limit — no findings",
			sql:       "CREATE TABLE t (name VARCHAR(100) NOT NULL)",
			maxLen:    200,
			wantCount: 0,
		},
		{
			name:      "varchar exceeds limit — one finding",
			sql:       "CREATE TABLE t (name VARCHAR(500) NOT NULL)",
			maxLen:    200,
			wantCount: 1,
		},
		{
			name:      "char exceeds limit — one finding",
			sql:       "CREATE TABLE t (code CHAR(10) NOT NULL)",
			maxLen:    5,
			wantCount: 1,
		},
		{
			name:      "varchar no length param — no findings",
			sql:       "CREATE TABLE t (name VARCHAR NOT NULL)",
			maxLen:    100,
			wantCount: 0,
		},
		{
			name:      "int column — no findings",
			sql:       "CREATE TABLE t (id INT NOT NULL)",
			maxLen:    100,
			wantCount: 0,
		},
		{
			name:      "default max (1024) — exactly at limit — no findings",
			sql:       "CREATE TABLE t (name VARCHAR(1024) NOT NULL)",
			maxLen:    0,
			wantCount: 0,
		},
		{
			name:      "default max (1024) — over limit — one finding",
			sql:       "CREATE TABLE t (name VARCHAR(1025) NOT NULL)",
			maxLen:    0,
			wantCount: 1,
		},
		{
			name:      "multiple columns mixed — one finding",
			sql:       "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50) NOT NULL, bio VARCHAR(5000) NOT NULL)",
			maxLen:    200,
			wantCount: 1,
		},
		{
			name:      "SELECT statement — no findings",
			sql:       "SELECT a FROM t",
			maxLen:    100,
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := rules.ColumnMaximumVarcharLengthConfig{MaxLength: tt.maxLen}
			rule := rules.NewColumnMaximumVarcharLengthRule(cfg)
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

func TestColumnMaximumVarcharLength_Metadata(t *testing.T) {
	r := rules.NewColumnMaximumVarcharLengthRule(rules.ColumnMaximumVarcharLengthConfig{MaxLength: 100})
	if r.ID() != "snowflake.column.maximum-varchar-length" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
