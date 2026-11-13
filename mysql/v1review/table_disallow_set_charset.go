package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowSetCharset judges every top-level CREATE TABLE and ALTER
// TABLE; statements nested in a routine body are not judged. Each
// character set option of a CREATE TABLE ([DEFAULT] CHARSET or CHARACTER
// SET, compared case-insensitively) is one TableCharsetSet, so a statement
// naming two has two findings; a COLLATE option is not judged. In an ALTER
// TABLE, each CONVERT TO CHARACTER SET action is one TableCharsetSet; a
// character set table option in an ALTER TABLE is not judged. Each finding
// anchors the whole statement and names nothing. Ported from Bytebase's V1
// rule TABLE_DISALLOW_SET_CHARSET.
//
// A nil CREATE TABLE option panics, as the V1 advisor did; the parser
// builds none.
func TableDisallowSetCharset(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			for _, opt := range n.Options {
				if opt != nil && strings.EqualFold(opt.Name, "CHARSET") || strings.EqualFold(opt.Name, "CHARACTER SET") || strings.EqualFold(opt.Name, "DEFAULT CHARSET") || strings.EqualFold(opt.Name, "DEFAULT CHARACTER SET") {
					start, end := span(n.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TableCharsetSet, Start: start, End: end})
				}
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATConvertCharset {
					start, end := span(n.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TableCharsetSet, Start: start, End: end})
				}
			}
		default:
		}
	}
	return findings
}
