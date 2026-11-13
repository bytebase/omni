package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableComment judges the COMMENT option of every top-level CREATE TABLE,
// including CREATE TABLE ... LIKE and ... AS SELECT; statements nested in a
// routine body are not judged. With cfg.Required, a statement with no
// COMMENT option is CommentMissing; an empty comment satisfies it. A
// comment longer than cfg.MaxLength bytes is CommentTooLong, and a negative
// MaxLength means no limit. The comment judged for length is the first
// COMMENT option's value. A statement can have both findings, missing first.
// Name is the table name as parsed, without its database; each finding
// anchors the whole statement. Ported from Bytebase's V1 rule TABLE_COMMENT.
//
// A nil cfg panics at the first judged statement, as the V1 advisor did
// when the rule had no payload.
func TableComment(stmts []ast.Node, cfg *v1review.TableCommentConfig) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.CreateTableStmt)
		if !ok || n.Table == nil {
			continue
		}
		comment, exists := "", false
		for _, opt := range n.Options {
			if strings.EqualFold(opt.Name, "COMMENT") {
				comment = opt.Value
				break
			}
		}
		for _, opt := range n.Options {
			if opt != nil && opt.Name == "COMMENT" {
				exists = true
				break
			}
		}
		exists = exists || comment != ""
		start, end := span(n.Loc)
		if cfg.Required && !exists {
			findings = append(findings, v1review.Finding{Code: v1review.TableCommentMissing, Name: n.Table.Name, Start: start, End: end})
		}
		if cfg.MaxLength >= 0 && len(comment) > cfg.MaxLength {
			findings = append(findings, v1review.Finding{Code: v1review.TableCommentTooLong, Name: n.Table.Name, Start: start, End: end})
		}
	}
	return findings
}
