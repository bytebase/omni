package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports every SELECT with a FROM clause and no
// WHERE clause anywhere in the statements: CTE bodies, derived tables,
// subqueries, MERGE sources, and statements nested in scripting blocks
// included. A SELECT without FROM, such as SELECT 1, is not judged. text is
// the statement text the nodes were parsed from.
//
// A finding anchors the SELECT. For a WITH ... SELECT, whose range starts at
// its WITH, it anchors the body's SELECT keyword instead: the first token in
// text after the last CTE, past white space and --, // and /* */ comments,
// as the V1 advisor found it; when no token follows, the WITH. Findings are
// in anchor order, not walk order. WhereMissing is the only code; Name is not
// set. Ported from Bytebase's V1 rule STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node, text string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.Where != nil || sel.From == nil {
				return true
			}
			start, end := stmtAnchor(sel.Loc)
			if k := len(sel.With); k > 0 {
				if last := sel.With[k-1]; last != nil && last.Loc.End >= 0 && last.Loc.End <= len(text) {
					if offset, ok := selectBodyToken(text, last.Loc.End); ok {
						start = offset
					}
				}
			}
			findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Start: start, End: max(end, start)})
			return true
		})
	}
	slices.SortStableFunc(findings, func(a, b v1review.Finding) int { return a.Start - b.Start })
	return findings
}

// selectBodyToken returns the offset of the first token in text at or after
// offset, past white space and --, // and /* */ comments; false when only
// those remain.
func selectBodyToken(text string, offset int) (int, bool) {
	for offset < len(text) {
		switch c := text[offset]; {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f':
			offset++
		case strings.HasPrefix(text[offset:], "--") || strings.HasPrefix(text[offset:], "//"):
			next := strings.IndexByte(text[offset:], '\n')
			if next < 0 {
				return 0, false
			}
			offset += next + 1
		case strings.HasPrefix(text[offset:], "/*"):
			end := strings.Index(text[offset+2:], "*/")
			if end < 0 {
				return 0, false
			}
			offset += 2 + end + 2
		default:
			return offset, true
		}
	}
	return 0, false
}
