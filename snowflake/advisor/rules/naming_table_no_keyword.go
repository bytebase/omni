package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// NamingTableNoKeywordRule fires when a CREATE TABLE uses a reserved keyword
// as its table name without quoting.
//
// Rule ID:  "snowflake.naming.table-no-keyword"
// Severity: ERROR
//
// Quoted identifiers are always permitted (Snowflake allows "SELECT" as an
// identifier when quoted). Only unquoted identifiers that collide with
// reserved keywords are flagged.
type NamingTableNoKeywordRule struct{}

const namingTableNoKeywordID = "snowflake.naming.table-no-keyword"

// ID implements Rule.
func (NamingTableNoKeywordRule) ID() string { return namingTableNoKeywordID }

// Severity implements Rule.
func (NamingTableNoKeywordRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (NamingTableNoKeywordRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	nameIdent := stmt.Name.Name
	// Quoted identifiers are always safe.
	if nameIdent.Quoted {
		return nil
	}
	if parser.IsReservedKeyword(nameIdent.Name) {
		return []*advisor.Finding{{
			RuleID:   namingTableNoKeywordID,
			Severity: advisor.SeverityError,
			Loc:      stmt.Name.Loc,
			Message: fmt.Sprintf(
				"table name %q is a reserved keyword; use a different name or quote it",
				nameIdent.Name,
			),
		}}
	}
	return nil
}
