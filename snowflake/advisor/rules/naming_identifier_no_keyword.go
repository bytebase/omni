package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// NamingIdentifierNoKeywordRule fires when an unquoted identifier is a
// reserved keyword. This complements NamingTableNoKeywordRule by covering
// additional identifier locations: column names, aliases, and CTE names.
//
// Rule ID:  "snowflake.naming.identifier-no-keyword"
// Severity: ERROR
//
// Quoted identifiers are exempt. Checked locations:
//   - Column names in CREATE TABLE
//   - Table alias in FROM / JOIN (TableRef.Alias)
//   - SELECT target alias (SelectTarget.Alias)
//   - CTE name (SelectStmt.With[i].Name)
type NamingIdentifierNoKeywordRule struct{}

const namingIdentifierNoKeywordID = "snowflake.naming.identifier-no-keyword"

// ID implements Rule.
func (NamingIdentifierNoKeywordRule) ID() string { return namingIdentifierNoKeywordID }

// Severity implements Rule.
func (NamingIdentifierNoKeywordRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (NamingIdentifierNoKeywordRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	switch n := node.(type) {
	case *ast.CreateTableStmt:
		return checkKeywordIdents(n)
	case *ast.SelectStmt:
		return checkKeywordSelectIdents(n)
	case *ast.TableRef:
		return checkKeywordIdent(n.Alias, n.Loc)
	}
	return nil
}

func checkKeywordIdents(stmt *ast.CreateTableStmt) []*advisor.Finding {
	var findings []*advisor.Finding
	for _, col := range stmt.Columns {
		findings = append(findings, checkKeywordIdent(col.Name, col.Loc)...)
	}
	return findings
}

func checkKeywordSelectIdents(stmt *ast.SelectStmt) []*advisor.Finding {
	var findings []*advisor.Finding
	for _, cte := range stmt.With {
		findings = append(findings, checkKeywordIdent(cte.Name, stmt.Loc)...)
	}
	for _, target := range stmt.Targets {
		findings = append(findings, checkKeywordIdent(target.Alias, target.Loc)...)
	}
	return findings
}

func checkKeywordIdent(ident ast.Ident, loc ast.Loc) []*advisor.Finding {
	if ident.IsEmpty() {
		return nil
	}
	if ident.Quoted {
		return nil
	}
	if parser.IsReservedKeyword(ident.Name) {
		return []*advisor.Finding{{
			RuleID:   namingIdentifierNoKeywordID,
			Severity: advisor.SeverityError,
			Loc:      loc,
			Message: fmt.Sprintf(
				"identifier %q is a reserved keyword; use a different name or quote it",
				ident.Name,
			),
		}}
	}
	return nil
}
