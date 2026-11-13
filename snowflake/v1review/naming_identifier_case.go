package v1review

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierCase judges the letter case of the column names
// top-level statements introduce: the same names, in the same order and
// at the same anchors, as NamingIdentifierNoKeyword.
//
// The V1 rule never reads its payload and always requires upper case, so
// the cfg parameter is ignored: CaseConfig{Upper: false} requires upper
// case too. A name is judged as Snowflake resolves it, an unquoted name
// upper-cased and a quoted name as its content, and is reported when
// upper-casing would change it. Only a quoted name can therefore be
// reported, and Finding.Name is its content without the quotes.
// NameCaseMismatch is the only code. Ported from Bytebase's V1 rule
// NAMING_IDENTIFIER_CASE.
func NamingIdentifierCase(stmts []ast.Node, _ v1review.CaseConfig) []v1review.Finding {
	var findings []v1review.Finding
	declaredColumns(stmts, func(ident ast.Ident, anchor ast.Loc) {
		name := ident.Normalize()
		if strings.ToUpper(name) == name {
			return
		}
		start, end := span(anchor)
		findings = append(findings, v1review.Finding{Code: v1review.NameCaseMismatch, Name: name, Start: start, End: end})
	})
	return findings
}
