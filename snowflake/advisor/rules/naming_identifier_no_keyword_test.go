package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
	"github.com/bytebase/omni/snowflake/ast"
)

func TestNamingIdentifierNoKeyword(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "normal identifiers — no findings",
			sql:       "CREATE TABLE orders (id INT NOT NULL, status VARCHAR(50))",
			wantCount: 0,
		},
		{
			name:      "quoted reserved keyword column — no finding",
			sql:       `CREATE TABLE t ("null" INT NOT NULL)`,
			wantCount: 0,
		},
		{
			name:      "non-reserved keyword column — no finding",
			sql:       "CREATE TABLE t (access INT NOT NULL)",
			wantCount: 0,
		},
		{
			name:      "normal SELECT — no findings",
			sql:       "SELECT a AS my_alias FROM t AS tbl WHERE 1=1",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.NamingIdentifierNoKeywordRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

// TestNamingIdentifierNoKeyword_DirectAST tests the rule with a hand-crafted AST
// node to verify the core logic fires when a reserved keyword appears as an
// unquoted identifier. This bypasses the parser which rejects reserved keywords
// in most identifier positions.
func TestNamingIdentifierNoKeyword_DirectAST(t *testing.T) {
	rule := rules.NamingIdentifierNoKeywordRule{}
	ctx := &advisor.Context{SQL: ""}

	// Simulate a CreateTableStmt where column name is the reserved keyword "null".
	stmt := &ast.CreateTableStmt{
		Name: &ast.ObjectName{Name: ast.Ident{Name: "orders"}},
		Columns: []*ast.ColumnDef{
			{
				Name:    ast.Ident{Name: "null", Quoted: false}, // unquoted reserved keyword
				NotNull: true,
			},
		},
	}
	findings := rule.Check(ctx, stmt)
	if len(findings) != 1 {
		t.Errorf("expected 1 finding for reserved keyword column, got %d", len(findings))
		return
	}
	if findings[0].RuleID != "snowflake.naming.identifier-no-keyword" {
		t.Errorf("RuleID = %q", findings[0].RuleID)
	}
}

// TestNamingIdentifierNoKeyword_QuotedSafe verifies that quoted reserved keyword
// identifiers do not trigger the rule.
func TestNamingIdentifierNoKeyword_QuotedSafe(t *testing.T) {
	rule := rules.NamingIdentifierNoKeywordRule{}
	ctx := &advisor.Context{SQL: ""}

	stmt := &ast.CreateTableStmt{
		Name: &ast.ObjectName{Name: ast.Ident{Name: "orders"}},
		Columns: []*ast.ColumnDef{
			{
				Name:    ast.Ident{Name: "select", Quoted: true}, // quoted — safe
				NotNull: true,
			},
		},
	}
	findings := rule.Check(ctx, stmt)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for quoted identifier, got %d", len(findings))
	}
}

func TestNamingIdentifierNoKeyword_Metadata(t *testing.T) {
	r := rules.NamingIdentifierNoKeywordRule{}
	if r.ID() != "snowflake.naming.identifier-no-keyword" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
