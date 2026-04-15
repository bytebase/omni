package rules

import (
	"fmt"
	"regexp"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

const (
	// DefaultTableNamePattern is the default regex for table names.
	DefaultTableNamePattern = `^[a-z][a-z0-9_]{0,63}$`
)

// NamingTableConfig holds the regex pattern for table names.
type NamingTableConfig struct {
	// Pattern is a regular expression that table names must match.
	// An empty string means use DefaultTableNamePattern.
	Pattern string
}

// namingTableRule fires when a CREATE TABLE name does not match the configured pattern.
//
// Rule ID:  "snowflake.naming.table"
// Severity: WARNING
//
// The check is applied to the last (local) part of the table name.
// Quoted identifiers are checked as-is (their case is intentional).
// Unquoted identifiers are checked with their source case (pre-resolution).
type namingTableRule struct {
	re *regexp.Regexp
}

const namingTableID = "snowflake.naming.table"

// NewNamingTableRule constructs the rule.
func NewNamingTableRule(cfg NamingTableConfig) advisor.Rule {
	pattern := cfg.Pattern
	if pattern == "" {
		pattern = DefaultTableNamePattern
	}
	return &namingTableRule{re: regexp.MustCompile(pattern)}
}

// ID implements Rule.
func (*namingTableRule) ID() string { return namingTableID }

// Severity implements Rule.
func (*namingTableRule) Severity() advisor.Severity { return advisor.SeverityWarning }

// Check implements Rule.
func (r *namingTableRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	tableName := stmt.Name.Name.Name
	if !r.re.MatchString(tableName) {
		return []*advisor.Finding{{
			RuleID:   namingTableID,
			Severity: advisor.SeverityWarning,
			Loc:      stmt.Name.Loc,
			Message: fmt.Sprintf(
				"table name %q does not match naming convention %q",
				tableName, r.re.String(),
			),
		}}
	}
	return nil
}
