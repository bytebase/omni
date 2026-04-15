package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// IdentifierCase enumerates the supported case conventions for identifiers.
type IdentifierCase int

const (
	// IdentifierCaseLower requires all-lowercase identifiers (snake_case friendly).
	IdentifierCaseLower IdentifierCase = iota
	// IdentifierCaseUpper requires all-uppercase identifiers.
	IdentifierCaseUpper
	// IdentifierCaseCamel requires camelCase identifiers (first char lowercase, at
	// least one uppercase letter, no underscores).
	IdentifierCaseCamel
)

// NamingIdentifierCaseConfig configures the required identifier case convention.
type NamingIdentifierCaseConfig struct {
	// Case is the required case convention. Default is IdentifierCaseLower.
	Case IdentifierCase
}

// namingIdentifierCaseRule fires when an unquoted identifier does not match
// the configured case convention.
//
// Rule ID:  "snowflake.naming.identifier-case"
// Severity: WARNING
//
// Quoted identifiers are exempt — their case is intentional.
//
// Checked locations:
//   - Column names in CREATE TABLE (ColumnDef.Name)
//   - Table name in CREATE TABLE (CreateTableStmt.Name.Name)
//   - Table alias in SelectStmt FROM clause (TableRef.Alias)
//   - SELECT target alias (SelectTarget.Alias)
//   - Table alias in UPDATE (UpdateStmt.Alias)
//   - Table alias in DELETE (DeleteStmt.Alias)
//   - CTE names (SelectStmt.With[i].Name)
type namingIdentifierCaseRule struct {
	caseConvention IdentifierCase
}

const namingIdentifierCaseID = "snowflake.naming.identifier-case"

// NewNamingIdentifierCaseRule constructs the rule.
func NewNamingIdentifierCaseRule(cfg NamingIdentifierCaseConfig) advisor.Rule {
	return &namingIdentifierCaseRule{caseConvention: cfg.Case}
}

// ID implements Rule.
func (*namingIdentifierCaseRule) ID() string { return namingIdentifierCaseID }

// Severity implements Rule.
func (*namingIdentifierCaseRule) Severity() advisor.Severity { return advisor.SeverityWarning }

// Check implements Rule.
func (r *namingIdentifierCaseRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	switch n := node.(type) {
	case *ast.CreateTableStmt:
		return r.checkCreateTable(n)
	case *ast.SelectStmt:
		return r.checkSelect(n)
	case *ast.UpdateStmt:
		return r.checkIdent(n.Alias, n.Loc)
	case *ast.DeleteStmt:
		return r.checkIdent(n.Alias, n.Loc)
	case *ast.TableRef:
		return r.checkIdent(n.Alias, n.Loc)
	}
	return nil
}

func (r *namingIdentifierCaseRule) checkCreateTable(stmt *ast.CreateTableStmt) []*advisor.Finding {
	var findings []*advisor.Finding
	// Table name.
	if f := r.checkIdent(stmt.Name.Name, stmt.Name.Loc); f != nil {
		findings = append(findings, f...)
	}
	// Column names.
	for _, col := range stmt.Columns {
		if f := r.checkIdent(col.Name, col.Loc); f != nil {
			findings = append(findings, f...)
		}
	}
	return findings
}

func (r *namingIdentifierCaseRule) checkSelect(stmt *ast.SelectStmt) []*advisor.Finding {
	var findings []*advisor.Finding
	// CTE names.
	for _, cte := range stmt.With {
		if f := r.checkIdent(cte.Name, stmt.Loc); f != nil {
			findings = append(findings, f...)
		}
	}
	// SELECT target aliases.
	for _, target := range stmt.Targets {
		if f := r.checkIdent(target.Alias, target.Loc); f != nil {
			findings = append(findings, f...)
		}
	}
	return findings
}

// checkIdent checks a single identifier. Returns a finding if it violates the case convention.
func (r *namingIdentifierCaseRule) checkIdent(ident ast.Ident, loc ast.Loc) []*advisor.Finding {
	if ident.IsEmpty() {
		return nil
	}
	if ident.Quoted {
		return nil
	}
	name := ident.Name
	if r.matchesCase(name) {
		return nil
	}
	return []*advisor.Finding{{
		RuleID:   namingIdentifierCaseID,
		Severity: advisor.SeverityWarning,
		Loc:      loc,
		Message: fmt.Sprintf(
			"identifier %q does not match the required %s case convention",
			name, r.caseLabel(),
		),
	}}
}

func (r *namingIdentifierCaseRule) matchesCase(name string) bool {
	switch r.caseConvention {
	case IdentifierCaseLower:
		return name == strings.ToLower(name)
	case IdentifierCaseUpper:
		return name == strings.ToUpper(name)
	case IdentifierCaseCamel:
		return isCamelCase(name)
	default:
		return true
	}
}

func (r *namingIdentifierCaseRule) caseLabel() string {
	switch r.caseConvention {
	case IdentifierCaseLower:
		return "LOWER"
	case IdentifierCaseUpper:
		return "UPPER"
	case IdentifierCaseCamel:
		return "CAMEL"
	default:
		return "LOWER"
	}
}

// isCamelCase reports whether name looks like camelCase:
// - starts with a lowercase letter
// - contains at least one uppercase letter
// - contains no underscores
func isCamelCase(name string) bool {
	if len(name) == 0 {
		return false
	}
	if !unicode.IsLower(rune(name[0])) {
		return false
	}
	if strings.ContainsRune(name, '_') {
		return false
	}
	for _, r := range name {
		if unicode.IsUpper(r) {
			return true
		}
	}
	// All lowercase, no underscores — acceptable as a single-word camel name.
	return true
}
