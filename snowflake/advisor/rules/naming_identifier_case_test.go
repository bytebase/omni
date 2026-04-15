package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestNamingIdentifierCase_Lower(t *testing.T) {
	cfg := rules.NamingIdentifierCaseConfig{Case: rules.IdentifierCaseLower}
	rule := rules.NewNamingIdentifierCaseRule(cfg)

	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "all lower — no findings",
			sql:       "CREATE TABLE orders (id INT NOT NULL, user_name VARCHAR(100))",
			wantCount: 0,
		},
		{
			name:      "table name uppercase — one finding",
			sql:       "CREATE TABLE Orders (id INT NOT NULL)",
			wantCount: 1,
		},
		{
			name:      "column name uppercase — one finding",
			sql:       "CREATE TABLE orders (ID INT NOT NULL)",
			wantCount: 1,
		},
		{
			name:      "both uppercase — two findings",
			sql:       "CREATE TABLE Orders (ID INT NOT NULL)",
			wantCount: 2,
		},
		{
			name:      "quoted uppercase — no finding",
			sql:       `CREATE TABLE "Orders" ("ID" INT NOT NULL)`,
			wantCount: 0,
		},
		{
			name:      "select alias lower — no findings",
			sql:       "SELECT a AS alias FROM t",
			wantCount: 0,
		},
		{
			name:      "select alias upper — one finding",
			sql:       "SELECT a AS MyAlias FROM t",
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestNamingIdentifierCase_Upper(t *testing.T) {
	cfg := rules.NamingIdentifierCaseConfig{Case: rules.IdentifierCaseUpper}
	rule := rules.NewNamingIdentifierCaseRule(cfg)

	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "all upper — no findings",
			sql:       "CREATE TABLE ORDERS (ID INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "lowercase table — one finding",
			sql:       "CREATE TABLE orders (ID INT NOT NULL)",
			wantCount: 1,
		},
		{
			name:      "uppercase with underscores — no findings",
			sql:       "CREATE TABLE ORDER_ITEMS (ITEM_ID INT NOT NULL)",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rule)
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
			}
		})
	}
}

func TestNamingIdentifierCase_Metadata(t *testing.T) {
	r := rules.NewNamingIdentifierCaseRule(rules.NamingIdentifierCaseConfig{})
	if r.ID() != "snowflake.naming.identifier-case" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityWarning {
		t.Errorf("Severity() = %v, want WARNING", r.Severity())
	}
}
