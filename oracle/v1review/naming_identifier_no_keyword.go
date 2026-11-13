package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierNoKeyword reports each name a statement defines that is a
// keyword in the V1 advisor's table (see isKeyword), compared without
// regard to case. The names are those identifiers collects: definitions at
// any depth, never references. A finding's Name is the name as the parser
// stores it and its range is the Loc identifiers pairs the name with, which
// for most names is the declaring node rather than the name itself.
// Findings are in statement order, then walk order. Ported from Bytebase's
// V1 rule NAMING_IDENTIFIER_NO_KEYWORD.
func NamingIdentifierNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		for _, ident := range identifiers(stmt) {
			if isKeyword(ident.name) {
				start, end := span(ident.loc)
				findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: ident.name, Start: start, End: end})
			}
		}
	}
	return findings
}
