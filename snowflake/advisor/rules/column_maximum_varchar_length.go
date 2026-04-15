package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

const (
	// DefaultMaxVarcharLength is the default maximum allowed length for VARCHAR/CHAR columns.
	DefaultMaxVarcharLength = 1024
)

// ColumnMaximumVarcharLengthConfig holds the maximum allowed VARCHAR/CHAR length.
type ColumnMaximumVarcharLengthConfig struct {
	// MaxLength is the maximum allowed length. Zero means use DefaultMaxVarcharLength.
	MaxLength int
}

// columnMaximumVarcharLengthRule fires when a VARCHAR or CHAR column declares a
// length greater than the configured maximum.
//
// Rule ID:  "snowflake.column.maximum-varchar-length"
// Severity: ERROR
//
// Columns with no explicit length parameter (unbounded VARCHAR) are not flagged
// because they have an implicit length determined by Snowflake (16MB for VARCHAR).
// Callers who want to forbid unbounded types should combine this rule with a
// separate type restriction rule.
type columnMaximumVarcharLengthRule struct {
	maxLength int
}

const columnMaximumVarcharLengthID = "snowflake.column.maximum-varchar-length"

// NewColumnMaximumVarcharLengthRule constructs the rule.
func NewColumnMaximumVarcharLengthRule(cfg ColumnMaximumVarcharLengthConfig) advisor.Rule {
	max := cfg.MaxLength
	if max <= 0 {
		max = DefaultMaxVarcharLength
	}
	return &columnMaximumVarcharLengthRule{maxLength: max}
}

// ID implements Rule.
func (*columnMaximumVarcharLengthRule) ID() string { return columnMaximumVarcharLengthID }

// Severity implements Rule.
func (*columnMaximumVarcharLengthRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (r *columnMaximumVarcharLengthRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	var findings []*advisor.Finding
	for _, col := range stmt.Columns {
		dt := col.DataType
		if dt == nil {
			continue
		}
		if dt.Kind != ast.TypeVarchar && dt.Kind != ast.TypeChar {
			continue
		}
		if len(dt.Params) == 0 {
			// No explicit length — skip.
			continue
		}
		length := dt.Params[0]
		if length > r.maxLength {
			findings = append(findings, &advisor.Finding{
				RuleID:   columnMaximumVarcharLengthID,
				Severity: advisor.SeverityError,
				Loc:      col.Loc,
				Message: fmt.Sprintf(
					"column %q declares %s(%d) which exceeds the maximum allowed length %d",
					col.Name.Name, dt.Name, length, r.maxLength,
				),
			})
		}
	}
	return findings
}
