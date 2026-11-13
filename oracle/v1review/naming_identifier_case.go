package v1review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierCase reports each name a statement defines that is not in
// the letter case cfg requires: a name that strings.ToUpper would change
// when cfg.Upper, one that strings.ToLower would change otherwise. The
// names are those identifiers collects: definitions at any depth, never
// references. A name is judged as the parser stores it, so an unquoted name
// is always upper case, whatever case the SQL wrote it in: it never fails
// the upper-case rule, and fails the lower-case rule unless it has no
// letters. A finding's Name is that stored name and its range is the Loc
// identifiers pairs the name with, which for most names is the declaring
// node rather than the name itself. Findings are in statement order, then
// walk order. Ported from Bytebase's V1 rule NAMING_IDENTIFIER_CASE.
func NamingIdentifierCase(stmts []ast.Node, cfg v1review.CaseConfig) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		for _, ident := range identifiers(stmt) {
			want := strings.ToLower(ident.name)
			if cfg.Upper {
				want = strings.ToUpper(ident.name)
			}
			if ident.name != want {
				start, end := span(ident.loc)
				findings = append(findings, v1review.Finding{Code: v1review.NameCaseMismatch, Name: ident.name, Start: start, End: end})
			}
		}
	}
	return findings
}
