package completion

import (
	"testing"

	"github.com/bytebase/omni/mssql/parser"
)

func TestResolveTokensToKeywords(t *testing.T) {
	// Build a CandidateSet with some token candidates.
	cs := &parser.CandidateSet{
		Tokens: []int{},
	}
	// Use Collect on "SELECT " to get real token candidates.
	cs = parser.Collect("SELECT ", 7)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}

	result := resolve(cs, nil, "SELECT ", 7)
	if len(result) == 0 {
		t.Fatal("expected non-empty result from token resolution")
	}

	// All results with nil catalog should be keywords (or types from type_name rule).
	for _, c := range result {
		if c.Type != CandidateKeyword && c.Type != CandidateType_ {
			t.Errorf("expected keyword or type candidate with nil catalog, got type %d for %q", c.Type, c.Text)
		}
	}
}

func TestResolveNilCandidateSet(t *testing.T) {
	result := resolve(nil, nil, "", 0)
	if result != nil {
		t.Errorf("expected nil for nil CandidateSet, got %v", result)
	}
}

func TestResolveTypeNameRule(t *testing.T) {
	// Directly test resolveRule for type_name.
	result := resolveRule("type_name", nil, "", 0)
	if len(result) == 0 {
		t.Fatal("expected type candidates for type_name rule")
	}

	// All should be CandidateType_.
	for _, c := range result {
		if c.Type != CandidateType_ {
			t.Errorf("expected CandidateType_ for type_name rule, got %d for %q", c.Type, c.Text)
		}
	}

	// Check some expected types are present.
	typeSet := make(map[string]bool)
	for _, c := range result {
		typeSet[c.Text] = true
	}

	expected := []string{
		"INT", "BIGINT", "SMALLINT", "TINYINT",
		"VARCHAR", "NVARCHAR", "CHAR", "NCHAR",
		"TEXT", "NTEXT",
		"DATETIME", "DATETIME2", "DATE", "TIME",
		"DECIMAL", "NUMERIC", "FLOAT", "REAL",
		"BIT",
		"MONEY", "SMALLMONEY",
		"UNIQUEIDENTIFIER",
		"XML",
		"VARBINARY", "IMAGE",
		"SQL_VARIANT",
	}
	for _, e := range expected {
		if !typeSet[e] {
			t.Errorf("expected type %q in type_name candidates", e)
		}
	}

	if len(result) != len(expected) {
		t.Errorf("expected %d type candidates, got %d", len(expected), len(result))
	}
}

func TestResolveRuleNilCatalog(t *testing.T) {
	// Rules that require catalog should return nil when catalog is nil.
	rules := []string{
		"table_ref", "columnref", "schema_ref", "func_name",
		"proc_ref", "index_ref", "trigger_ref", "database_ref",
		"sequence_ref", "login_ref", "user_ref", "role_ref",
	}
	for _, rule := range rules {
		result := resolveRule(rule, nil, "", 0)
		if result != nil {
			t.Errorf("expected nil for rule %q with nil catalog, got %v", rule, result)
		}
	}
}

func TestResolveUnknownRule(t *testing.T) {
	result := resolveRule("unknown_rule", nil, "", 0)
	if result != nil {
		t.Errorf("expected nil for unknown rule, got %v", result)
	}
}

func TestResolveDedupKeywords(t *testing.T) {
	// Create a CandidateSet with duplicate tokens.
	cs := &parser.CandidateSet{
		Tokens: []int{},
	}
	// Collect from empty input — should get top-level keywords, all deduped.
	cs = parser.Collect("", 0)
	result := resolve(cs, nil, "", 0)

	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool)
	for _, c := range result {
		k := key{c.Text, c.Type}
		if seen[k] {
			t.Errorf("duplicate candidate after resolve: %q (type %d)", c.Text, c.Type)
		}
		seen[k] = true
	}
}

func TestResolveTypeNameWithCandidateSet(t *testing.T) {
	// Simulate a CandidateSet that includes a type_name rule.
	cs := &parser.CandidateSet{
		Tokens: []int{},
		Rules:  []parser.RuleCandidate{{Rule: "type_name"}},
	}
	result := resolve(cs, nil, "", 0)

	// Should contain type candidates.
	hasType := false
	for _, c := range result {
		if c.Type == CandidateType_ {
			hasType = true
			break
		}
	}
	if !hasType {
		t.Error("expected type candidates when type_name rule is in CandidateSet")
	}
}
