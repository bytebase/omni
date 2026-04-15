package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
	"github.com/bytebase/omni/snowflake/ast"
)

// allTestRules returns all 14 rules configured for integration tests.
func allTestRules() []advisor.Rule {
	return []advisor.Rule{
		rules.ColumnNoNullRule{},
		rules.NewColumnRequireRule(rules.ColumnRequireConfig{Required: []string{"created_at"}}),
		rules.NewColumnMaximumVarcharLengthRule(rules.ColumnMaximumVarcharLengthConfig{MaxLength: 100}),
		rules.TableRequirePKRule{},
		rules.TableNoForeignKeyRule{},
		rules.NewNamingTableRule(rules.NamingTableConfig{}),
		rules.NamingTableNoKeywordRule{},
		rules.NewNamingIdentifierCaseRule(rules.NamingIdentifierCaseConfig{Case: rules.IdentifierCaseLower}),
		rules.NamingIdentifierNoKeywordRule{},
		rules.SelectNoSelectAllRule{},
		rules.WhereRequireSelectRule{},
		rules.WhereRequireUpdateDeleteRule{},
		rules.MigrationCompatibilityRule{},
		rules.NewTableDropNamingConventionRule(rules.TableDropNamingConventionConfig{}),
	}
}

// TestIntegration_EachRuleFiresOnItsOwnTrigger runs every rule individually
// against a SQL that is known to trigger it (or a direct AST node for rules
// that cannot be triggered via parsed SQL) and asserts at least one finding
// with the expected rule ID.
func TestIntegration_EachRuleFiresOnItsOwnTrigger(t *testing.T) {
	rule := rules.NamingIdentifierNoKeywordRule{}
	ctx := &advisor.Context{}
	// Hand-craft an AST node for naming_identifier_no_keyword since the parser
	// rejects reserved keywords in column-name positions.
	reservedKwStmt := &ast.CreateTableStmt{
		Name: &ast.ObjectName{Name: ast.Ident{Name: "orders"}},
		Columns: []*ast.ColumnDef{
			{Name: ast.Ident{Name: "null", Quoted: false}, NotNull: true},
		},
	}
	kwFindings := rule.Check(ctx, reservedKwStmt)

	triggers := []struct {
		ruleID   string
		rule     advisor.Rule
		sql      string
		findings []*advisor.Finding // if non-nil, skip SQL parsing and use these
	}{
		{
			ruleID: "snowflake.column.no-null",
			rule:   rules.ColumnNoNullRule{},
			sql:    "CREATE TABLE t (id INT)",
		},
		{
			ruleID: "snowflake.column.require",
			rule:   rules.NewColumnRequireRule(rules.ColumnRequireConfig{Required: []string{"created_at"}}),
			sql:    "CREATE TABLE t (id INT NOT NULL)",
		},
		{
			ruleID: "snowflake.column.maximum-varchar-length",
			rule:   rules.NewColumnMaximumVarcharLengthRule(rules.ColumnMaximumVarcharLengthConfig{MaxLength: 100}),
			sql:    "CREATE TABLE t (id INT NOT NULL, bio VARCHAR(500) NOT NULL)",
		},
		{
			ruleID: "snowflake.table.require-pk",
			rule:   rules.TableRequirePKRule{},
			sql:    "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50) NOT NULL)",
		},
		{
			ruleID: "snowflake.table.no-foreign-key",
			rule:   rules.TableNoForeignKeyRule{},
			sql:    "CREATE TABLE t (id INT NOT NULL, user_id INT FOREIGN KEY REFERENCES users(id))",
		},
		{
			ruleID: "snowflake.naming.table",
			rule:   rules.NewNamingTableRule(rules.NamingTableConfig{}),
			sql:    "CREATE TABLE BadTable (id INT NOT NULL)",
		},
		{
			ruleID: "snowflake.naming.table-no-keyword",
			rule:   rules.NamingTableNoKeywordRule{},
			sql:    "CREATE TABLE select (id INT NOT NULL)",
		},
		{
			ruleID: "snowflake.naming.identifier-case",
			rule:   rules.NewNamingIdentifierCaseRule(rules.NamingIdentifierCaseConfig{Case: rules.IdentifierCaseLower}),
			sql:    "CREATE TABLE good_table (BadCol INT NOT NULL)",
		},
		{
			// naming_identifier_no_keyword: parser rejects reserved keywords in most
			// identifier positions, so we provide the pre-computed findings directly.
			ruleID:   "snowflake.naming.identifier-no-keyword",
			rule:     rules.NamingIdentifierNoKeywordRule{},
			findings: kwFindings,
		},
		{
			ruleID: "snowflake.select.no-select-all",
			rule:   rules.SelectNoSelectAllRule{},
			sql:    "SELECT * FROM t",
		},
		{
			ruleID: "snowflake.query.where-require-select",
			rule:   rules.WhereRequireSelectRule{},
			sql:    "SELECT a FROM t",
		},
		{
			ruleID: "snowflake.query.where-require-update-delete",
			rule:   rules.WhereRequireUpdateDeleteRule{},
			sql:    "UPDATE t SET a = 1",
		},
		{
			ruleID: "snowflake.migration.compatibility",
			rule:   rules.MigrationCompatibilityRule{},
			sql:    "DROP TABLE bad_table",
		},
		{
			ruleID: "snowflake.table.drop-naming-convention",
			rule:   rules.NewTableDropNamingConventionRule(rules.TableDropNamingConventionConfig{}),
			sql:    "DROP TABLE bad_table",
		},
	}

	for _, tt := range triggers {
		t.Run(tt.ruleID, func(t *testing.T) {
			var findings []*advisor.Finding
			if tt.findings != nil {
				findings = tt.findings
			} else {
				findings = runRule(t, tt.sql, tt.rule)
			}
			fired := false
			for _, f := range findings {
				if f.RuleID == tt.ruleID {
					fired = true
					break
				}
			}
			if !fired {
				t.Errorf("rule %q did not fire (got %d findings)", tt.ruleID, len(findings))
				for _, f := range findings {
					t.Logf("  got finding: ruleID=%s msg=%s", f.RuleID, f.Message)
				}
			}
		})
	}
}

// TestIntegration_AllRulesNoFalsePositive verifies that a well-written CREATE
// TABLE with no violations produces no findings from any of the 14 rules.
func TestIntegration_AllRulesNoFalsePositive(t *testing.T) {
	// This SQL should satisfy all column/table/naming rules:
	// - columns all have NOT NULL
	// - required column created_at is present
	// - no VARCHAR over 100 chars
	// - has a PRIMARY KEY
	// - no foreign keys
	// - table name matches default pattern
	// - not a reserved keyword
	// - all identifiers are lowercase
	// - column names are not reserved keywords
	sql := `CREATE TABLE orders (
		id INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		status VARCHAR(50) NOT NULL,
		PRIMARY KEY (id)
	)`
	findings := runRules(t, sql, allTestRules()...)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for clean CREATE TABLE, got %d:", len(findings))
		for _, f := range findings {
			t.Logf("  ruleID=%s msg=%s", f.RuleID, f.Message)
		}
	}
}

// TestIntegration_AllRulesCount verifies that exactly 14 unique rule IDs are
// registered in allTestRules().
func TestIntegration_AllRulesCount(t *testing.T) {
	ruleList := allTestRules()
	seen := make(map[string]bool)
	for _, r := range ruleList {
		id := r.ID()
		if seen[id] {
			t.Errorf("duplicate rule ID: %q", id)
		}
		seen[id] = true
	}
	const want = 14
	if len(seen) != want {
		t.Errorf("got %d unique rule IDs, want %d", len(seen), want)
		for id := range seen {
			t.Logf("  %s", id)
		}
	}
}

// TestIntegration_MixedStatements verifies that rules correctly fire across
// multiple statement types in the same advisor run (each statement independently).
func TestIntegration_MixedStatements(t *testing.T) {
	tests := []struct {
		sql    string
		ruleID string
		want   int
	}{
		{"SELECT * FROM t", "snowflake.select.no-select-all", 1},
		{"SELECT a FROM t WHERE id = 1", "snowflake.select.no-select-all", 0},
		{"UPDATE t SET a = 1", "snowflake.query.where-require-update-delete", 1},
		{"UPDATE t SET a = 1 WHERE id = 1", "snowflake.query.where-require-update-delete", 0},
		{"DELETE FROM t", "snowflake.query.where-require-update-delete", 1},
		{"DROP TABLE my_table", "snowflake.migration.compatibility", 1},
		{"DROP TABLE deleted_orders", "snowflake.table.drop-naming-convention", 0},
		{"DROP TABLE my_table", "snowflake.table.drop-naming-convention", 1},
	}
	for _, tt := range tests {
		t.Run(tt.sql+"/"+tt.ruleID, func(t *testing.T) {
			// Find the specific rule.
			var targetRule advisor.Rule
			for _, r := range allTestRules() {
				if r.ID() == tt.ruleID {
					targetRule = r
					break
				}
			}
			if targetRule == nil {
				t.Fatalf("rule %q not found", tt.ruleID)
			}
			findings := findingsByRule(runRule(t, tt.sql, targetRule), tt.ruleID)
			if len(findings) != tt.want {
				t.Errorf("got %d findings, want %d", len(findings), tt.want)
			}
		})
	}
}
