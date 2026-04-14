package advisor_test

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// mustParse parses sql and fails the test on any parse error.
func mustParse(t *testing.T, sql string) *ast.File {
	t.Helper()
	f, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse(%q): %v", sql, err)
	}
	return f
}

// newCtx returns a Context for the given SQL string.
func newCtx(sql string) *advisor.Context {
	return &advisor.Context{SQL: sql}
}

// ---------------------------------------------------------------------------
// countingRule — records every node it sees (for framework smoke tests).
// ---------------------------------------------------------------------------

type countingRule struct {
	count int
}

func (r *countingRule) ID() string                 { return "test.counting" }
func (r *countingRule) Severity() advisor.Severity { return advisor.SeverityInfo }
func (r *countingRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	r.count++
	return nil
}

// alwaysRule emits one finding for every node it sees.
type alwaysRule struct{}

func (alwaysRule) ID() string                 { return "test.always" }
func (alwaysRule) Severity() advisor.Severity { return advisor.SeverityError }
func (alwaysRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	return []*advisor.Finding{{
		RuleID:   "test.always",
		Severity: advisor.SeverityError,
		Loc:      ast.NodeLoc(node),
		Message:  "always fires",
	}}
}

// ---------------------------------------------------------------------------
// Framework tests
// ---------------------------------------------------------------------------

// TestAdvisor_EmptyRules verifies that an Advisor with no rules produces no
// findings even on a real AST.
func TestAdvisor_EmptyRules(t *testing.T) {
	file := mustParse(t, "SELECT 1")
	a := advisor.New()
	findings := a.Check(newCtx("SELECT 1"), file)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

// TestAdvisor_EmptyAST verifies that no findings are produced for a nil root.
func TestAdvisor_EmptyAST(t *testing.T) {
	a := advisor.New(alwaysRule{})
	findings := a.Check(newCtx(""), nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil root, got %d", len(findings))
	}
}

// TestAdvisor_RuleInvoked verifies that a rule's Check method is called at
// least once for each node in the AST.
func TestAdvisor_RuleInvoked(t *testing.T) {
	sql := "SELECT 1"
	file := mustParse(t, sql)
	rule := &countingRule{}
	a := advisor.New(rule)
	a.Check(newCtx(sql), file)
	if rule.count == 0 {
		t.Error("expected rule.Check to be called at least once")
	}
}

// TestAdvisor_MultipleRulesAllInvoked verifies that every registered rule
// receives every node.
func TestAdvisor_MultipleRulesAllInvoked(t *testing.T) {
	sql := "SELECT 1"
	file := mustParse(t, sql)
	r1, r2, r3 := &countingRule{}, &countingRule{}, &countingRule{}
	a := advisor.New(r1, r2, r3)
	a.Check(newCtx(sql), file)
	if r1.count == 0 || r2.count == 0 || r3.count == 0 {
		t.Errorf("all rules should be invoked; got counts %d/%d/%d", r1.count, r2.count, r3.count)
	}
	if r1.count != r2.count || r2.count != r3.count {
		t.Errorf("all rules should see the same node count; got %d/%d/%d", r1.count, r2.count, r3.count)
	}
}

// TestAdvisor_FindingsCollected verifies that findings from all rules are
// returned and that their RuleID / Severity fields are set correctly.
func TestAdvisor_FindingsCollected(t *testing.T) {
	sql := "SELECT 1"
	file := mustParse(t, sql)
	a := advisor.New(alwaysRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	for _, f := range findings {
		if f.RuleID != "test.always" {
			t.Errorf("RuleID = %q, want %q", f.RuleID, "test.always")
		}
		if f.Severity != advisor.SeverityError {
			t.Errorf("Severity = %v, want ERROR", f.Severity)
		}
	}
}

// TestAdvisor_Reuse verifies that the same Advisor can be used for multiple
// Check calls without cross-contamination.
func TestAdvisor_Reuse(t *testing.T) {
	sql := "SELECT 1"
	file := mustParse(t, sql)
	a := advisor.New(alwaysRule{})
	ctx := newCtx(sql)
	f1 := a.Check(ctx, file)
	f2 := a.Check(ctx, file)
	if len(f1) != len(f2) {
		t.Errorf("reuse: run1=%d findings, run2=%d findings; want equal", len(f1), len(f2))
	}
}

// TestAdvisor_FindingOrderIsRegistrationOrder verifies that findings from
// multiple rules are returned in rule-registration order per node.
func TestAdvisor_FindingOrderIsRegistrationOrder(t *testing.T) {
	var calls []string
	sql := "SELECT 1"
	file := mustParse(t, sql)
	a := advisor.New(
		trackingRule{tag: "A", calls: &calls},
		trackingRule{tag: "B", calls: &calls},
		trackingRule{tag: "C", calls: &calls},
	)
	a.Check(newCtx(sql), file)

	// For each node, A should be called before B before C.
	// calls is: [A, B, C, A, B, C, ...] — one group per node.
	if len(calls) == 0 {
		t.Fatal("no calls recorded")
	}
	for i := 0; i+2 < len(calls); i += 3 {
		if calls[i] != "A" || calls[i+1] != "B" || calls[i+2] != "C" {
			t.Errorf("calls[%d:%d] = %v, want [A B C]", i, i+3, calls[i:i+3])
		}
	}
}

type trackingRule struct {
	tag   string
	calls *[]string
}

func (r trackingRule) ID() string                 { return "test.tracking." + r.tag }
func (r trackingRule) Severity() advisor.Severity { return advisor.SeverityInfo }
func (r trackingRule) Check(_ *advisor.Context, _ ast.Node) []*advisor.Finding {
	*r.calls = append(*r.calls, r.tag)
	return nil
}

// ---------------------------------------------------------------------------
// Severity constant tests
// ---------------------------------------------------------------------------

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		sev  advisor.Severity
		want string
	}{
		{advisor.SeverityInfo, "INFO"},
		{advisor.SeverityWarning, "WARNING"},
		{advisor.SeverityError, "ERROR"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// NoSelectStarRule tests
// ---------------------------------------------------------------------------

// TestNoSelectStar_BareStarFires verifies that SELECT * triggers a WARNING.
func TestNoSelectStar_BareStarFires(t *testing.T) {
	sql := "SELECT * FROM t"
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	f := findings[0]
	if f.RuleID != "snowflake.select.no-select-star" {
		t.Errorf("RuleID = %q, want %q", f.RuleID, "snowflake.select.no-select-star")
	}
	if f.Severity != advisor.SeverityWarning {
		t.Errorf("Severity = %v, want WARNING", f.Severity)
	}
	if !strings.Contains(f.Message, "SELECT *") {
		t.Errorf("Message = %q, should mention SELECT *", f.Message)
	}
}

// TestNoSelectStar_ColumnListNoFire verifies that SELECT a, b produces no
// findings.
func TestNoSelectStar_ColumnListNoFire(t *testing.T) {
	sql := "SELECT a, b, c FROM t"
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for column list, got %d", len(findings))
	}
}

// TestNoSelectStar_QualifiedStarNoFire verifies that SELECT t.* does NOT
// trigger the rule (qualified star is a separate concern).
func TestNoSelectStar_QualifiedStarNoFire(t *testing.T) {
	sql := "SELECT t.* FROM t"
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for qualified star, got %d", len(findings))
	}
}

// TestNoSelectStar_MultipleStarFires verifies that two bare stars in a UNION
// produce two findings.
func TestNoSelectStar_MultipleStarFires(t *testing.T) {
	sql := "SELECT * FROM a UNION ALL SELECT * FROM b"
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings for two bare stars, got %d", len(findings))
	}
}

// TestNoSelectStar_LocationIsValid verifies that the Loc on a finding covers
// the source bytes of the * token.
func TestNoSelectStar_LocationIsValid(t *testing.T) {
	sql := "SELECT * FROM t"
	// The * is at byte offset 7 (0-indexed: "SELECT " is 7 bytes, * is at 7).
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	loc := findings[0].Loc
	if !loc.IsValid() {
		t.Fatalf("Loc is not valid: %+v", loc)
	}
	// The SelectTarget covering * should start at/before the * character and
	// end after it.  "SELECT * FROM t" — * is at index 7.
	if loc.Start > 7 || loc.End <= 7 {
		t.Errorf("Loc = {%d, %d}, want Start<=7 and End>7 to cover '*' at offset 7",
			loc.Start, loc.End)
	}
	// The source substring covered by Loc should contain *.
	sub := sql[loc.Start:loc.End]
	if !strings.Contains(sub, "*") {
		t.Errorf("source[%d:%d] = %q, should contain '*'", loc.Start, loc.End, sub)
	}
}

// TestNoSelectStar_IDAndSeverity verifies the rule's static metadata.
func TestNoSelectStar_IDAndSeverity(t *testing.T) {
	rule := advisor.NoSelectStarRule{}
	if rule.ID() != "snowflake.select.no-select-star" {
		t.Errorf("ID() = %q, want %q", rule.ID(), "snowflake.select.no-select-star")
	}
	if rule.Severity() != advisor.SeverityWarning {
		t.Errorf("Severity() = %v, want WARNING", rule.Severity())
	}
}

// TestNoSelectStar_NoSelectNoFire verifies a CREATE TABLE produces no finding.
func TestNoSelectStar_NoSelectNoFire(t *testing.T) {
	sql := "CREATE TABLE t (id INT)"
	file := mustParse(t, sql)
	a := advisor.New(advisor.NoSelectStarRule{})
	findings := a.Check(newCtx(sql), file)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for DDL, got %d", len(findings))
	}
}
