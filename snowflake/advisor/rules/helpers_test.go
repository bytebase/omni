package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// runRule is a test helper that parses sql, runs a single Rule via the
// Advisor framework, and returns the findings.
func runRule(t *testing.T, sql string, rule advisor.Rule) []*advisor.Finding {
	t.Helper()
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse(%q): %v", sql, err)
	}
	a := advisor.New(rule)
	ctx := &advisor.Context{SQL: sql}
	return a.Check(ctx, file)
}

// runRules is like runRule but accepts multiple rules.
func runRules(t *testing.T, sql string, ruleList ...advisor.Rule) []*advisor.Finding {
	t.Helper()
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse(%q): %v", sql, err)
	}
	a := advisor.New(ruleList...)
	ctx := &advisor.Context{SQL: sql}
	return a.Check(ctx, file)
}

// findingsByRule returns the subset of findings whose RuleID equals id.
func findingsByRule(findings []*advisor.Finding, id string) []*advisor.Finding {
	var out []*advisor.Finding
	for _, f := range findings {
		if f.RuleID == id {
			out = append(out, f)
		}
	}
	return out
}

// mustParse parses sql and fails the test on any parse error.
// Kept here so individual test files can use it if needed.
func mustParse(t *testing.T, sql string) *ast.File {
	t.Helper()
	f, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse(%q): %v", sql, err)
	}
	return f
}
