package advisor

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// Rule is the interface implemented by every lint rule.
//
// Rules are stateless with respect to a single Check call: all per-run
// context arrives via the ctx argument, and all output is returned as
// a slice (never stored on the Rule). This makes rules safe to reuse
// across multiple Advisor.Check calls without re-registration.
type Rule interface {
	// ID returns the canonical identifier for this rule.
	// Convention: "snowflake.<category>.<short-name>", e.g.
	// "snowflake.select.no-select-star".
	ID() string

	// Severity returns the default severity for findings produced by this rule.
	// Implementations may vary severity per finding if needed; the value
	// returned here is used for documentation, filtering, and tests.
	Severity() Severity

	// Check inspects node and returns zero or more Findings.
	// It is called once per AST node during a depth-first walk.
	// Check must be non-blocking and must not mutate ctx or node.
	Check(ctx *Context, node ast.Node) []*Finding
}
