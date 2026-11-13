package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCommentLength judges the comment each top-level COMMENT ON statement
// sets, on an object of any kind, against maxLength, counted in bytes. A
// maxLength of 0 or less requires nothing. A comment in a CREATE statement's
// options, or a COMMENT ON under EXPLAIN or in a routine body, is not judged.
//
// Every finding anchors the whole statement. Ported from Bytebase's V1 rule
// SYSTEM_COMMENT_LENGTH.
func SystemCommentLength(stmts []ast.Node, maxLength int) []v1review.Finding {
	if maxLength <= 0 {
		return nil
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		n, ok := stmt.(*ast.CommentStmt)
		if !ok || len(n.Comment) <= maxLength {
			continue
		}
		start, end := span(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.CommentOnTooLong, Start: start, End: end})
	}
	return findings
}
