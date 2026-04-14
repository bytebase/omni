// Package advisor implements the lint-rule framework for Snowflake SQL.
//
// Usage:
//
//	a := advisor.New(myRule1, myRule2)
//	ctx := &advisor.Context{SQL: sql}
//	findings := a.Check(ctx, parsedFile)
//
// The Advisor walks the AST depth-first and fans out each node to every
// registered Rule. Findings from all rules are returned in traversal order.
//
// See rule.go for the Rule interface, and example_rule.go for a minimal
// worked example (NoSelectStarRule).
package advisor

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// Severity classifies how serious a finding is.
// The numeric ordering is intentional: higher value = more severe.
type Severity int

const (
	// SeverityInfo is an informational notice; does not block deployment.
	SeverityInfo Severity = iota
	// SeverityWarning flags a practice that should be reviewed.
	SeverityWarning
	// SeverityError flags a violation that must be corrected.
	SeverityError
)

// String returns the human-readable severity label used in messages and logs.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Finding is a single lint result emitted by a Rule.
type Finding struct {
	// RuleID is the canonical identifier of the rule that produced this finding.
	RuleID string
	// Severity is the severity of this specific finding.
	Severity Severity
	// Loc is the source byte range of the AST node that triggered the finding.
	// Callers can map Loc.Start/End to line/column using a line table.
	Loc ast.Loc
	// Message is a human-readable description of the issue.
	Message string
}

// Context carries per-run information that rules can inspect.
// It is intentionally minimal for T2.6; T2.7 rules may extend it
// when catalog or session-level data is needed.
type Context struct {
	// SQL is the original source text passed to the parser.
	// Rules may extract source snippets using Loc.Start / Loc.End as
	// indices into this string.
	SQL string
}

// Advisor runs a fixed set of Rules against an AST.
// Each Advisor instance is safe to reuse for multiple Check calls.
type Advisor struct {
	rules []Rule
}

// New constructs an Advisor with the given rules. Rules are invoked in
// the order they are provided; order only affects the order of findings
// in the returned slice, not correctness.
func New(rules ...Rule) *Advisor {
	r := make([]Rule, len(rules))
	copy(r, rules)
	return &Advisor{rules: r}
}

// Check walks root depth-first and fans each node out to every Rule.
// The returned slice contains all findings from all rules, in traversal order
// (parent before children, rules in registration order per node).
//
// Check is safe to call concurrently on different Context/root pairs as long
// as no rule itself has shared mutable state.
func (a *Advisor) Check(ctx *Context, root ast.Node) []*Finding {
	run := &advisorRun{advisor: a, ctx: ctx}
	ast.Walk(run, root)
	return run.findings
}

// advisorRun holds per-Check state so that Advisor itself is stateless.
type advisorRun struct {
	advisor  *Advisor
	ctx      *Context
	findings []*Finding
}

// Visit implements ast.Visitor. It is called for each node during the walk.
// A nil node signals the post-order (end-of-children) event; we return nil
// to terminate that branch — no action needed.
func (r *advisorRun) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	for _, rule := range r.advisor.rules {
		r.findings = append(r.findings, rule.Check(r.ctx, node)...)
	}
	return r // continue descending into children
}
