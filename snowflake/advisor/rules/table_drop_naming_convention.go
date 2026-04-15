package rules

import (
	"fmt"
	"regexp"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

const (
	// DefaultDropNamingPattern is the default regex for dropped table names.
	// Matches tables whose name starts with "deleted_" or whose name is "_deleted".
	DefaultDropNamingPattern = `^_deleted$|^deleted_`
)

// TableDropNamingConventionConfig holds the regex pattern for dropped table names.
type TableDropNamingConventionConfig struct {
	// Pattern is a regular expression that the dropped table name must match.
	// An empty string means use DefaultDropNamingPattern.
	Pattern string
}

// tableDropNamingConventionRule fires when a DROP TABLE target does not match
// the configured archival naming convention.
//
// Rule ID:  "snowflake.table.drop-naming-convention"
// Severity: WARNING
//
// The intent is to require that tables be renamed to an archival name (e.g.
// deleted_orders or orders_deleted) before dropping, so that other teams
// have a chance to notice the deletion.
type tableDropNamingConventionRule struct {
	re *regexp.Regexp
}

const tableDropNamingConventionID = "snowflake.table.drop-naming-convention"

// NewTableDropNamingConventionRule constructs the rule.
func NewTableDropNamingConventionRule(cfg TableDropNamingConventionConfig) advisor.Rule {
	pattern := cfg.Pattern
	if pattern == "" {
		pattern = DefaultDropNamingPattern
	}
	return &tableDropNamingConventionRule{re: regexp.MustCompile(pattern)}
}

// ID implements Rule.
func (*tableDropNamingConventionRule) ID() string { return tableDropNamingConventionID }

// Severity implements Rule.
func (*tableDropNamingConventionRule) Severity() advisor.Severity { return advisor.SeverityWarning }

// Check implements Rule.
func (r *tableDropNamingConventionRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	drop, ok := node.(*ast.DropStmt)
	if !ok {
		return nil
	}
	if drop.Kind != ast.DropTable {
		return nil
	}
	tableName := drop.Name.Name.Name
	if r.re.MatchString(tableName) {
		return nil
	}
	return []*advisor.Finding{{
		RuleID:   tableDropNamingConventionID,
		Severity: advisor.SeverityWarning,
		Loc:      drop.Loc,
		Message: fmt.Sprintf(
			"table %q being dropped does not match the archival naming convention %q; "+
				"consider renaming to an archival name before dropping",
			tableName, r.re.String(),
		),
	}}
}
