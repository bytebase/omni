package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the tables top-level DROP TABLE statements drop
// against cfg.Pattern: every table the statement lists, in the order
// written. A DROP TABLE nested in a routine body or an IF or BEGIN block is
// not judged, nor are DROP EXTERNAL TABLE, DROP VIEW, and TRUNCATE TABLE.
// The name is judged without its server, database or schema qualifier and
// without its brackets or double quotes, lower-cased with strings.ToLower,
// a temporary table's # prefix included; the finding's Name is that
// lower-cased text.
//
// cfg.MaxLength is not read: the rule has no length limit. Every finding
// anchors the whole statement. NameMismatch is the only code. Ported from
// Bytebase's V1 rule TABLE_DROP_NAMING_CONVENTION.
func TableDropNaming(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil {
		return nil, nil
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		drop, ok := stmt.(*ast.DropStmt)
		if !ok || drop.ObjectType != ast.DropTable || drop.Names == nil {
			continue
		}
		start, end := span(drop.Loc)
		for _, item := range drop.Names.Items {
			ref, ok := item.(*ast.TableRef)
			if !ok || ref == nil {
				continue
			}
			name := strings.ToLower(ref.Object)
			if name == "" || pattern.MatchString(name) {
				continue
			}
			findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
		}
	}
	return findings, nil
}
