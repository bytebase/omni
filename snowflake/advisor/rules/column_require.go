package rules

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// ColumnRequireConfig holds the list of column names that every CREATE TABLE
// must define.
type ColumnRequireConfig struct {
	// Required is a list of column names (case-insensitive). An empty list
	// means the rule is a no-op.
	Required []string
}

// columnRequireRule fires when a CREATE TABLE is missing one or more required
// columns.
//
// Rule ID:  "snowflake.column.require"
// Severity: ERROR
type columnRequireRule struct {
	required []string // normalized (uppercased) required column names
}

const columnRequireID = "snowflake.column.require"

// NewColumnRequireRule constructs the rule with the given configuration.
func NewColumnRequireRule(cfg ColumnRequireConfig) advisor.Rule {
	normalized := make([]string, len(cfg.Required))
	for i, n := range cfg.Required {
		normalized[i] = strings.ToUpper(n)
	}
	return &columnRequireRule{required: normalized}
}

// ID implements Rule.
func (*columnRequireRule) ID() string { return columnRequireID }

// Severity implements Rule.
func (*columnRequireRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (r *columnRequireRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	if len(r.required) == 0 {
		return nil
	}
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	// Build a set of present column names.
	present := make(map[string]bool, len(stmt.Columns))
	for _, col := range stmt.Columns {
		present[strings.ToUpper(col.Name.Name)] = true
	}
	var findings []*advisor.Finding
	for _, req := range r.required {
		if !present[req] {
			findings = append(findings, &advisor.Finding{
				RuleID:   columnRequireID,
				Severity: advisor.SeverityError,
				Loc:      stmt.Loc,
				Message:  fmt.Sprintf("table %q must have a column named %q", stmt.Name.String(), req),
			})
		}
	}
	return findings
}
