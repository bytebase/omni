package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// RewriteResultLimit rewrites one statement so its outermost query returns at
// most limit rows, splicing into the original text by node offsets. The output
// is the statement with trailing whitespace and ';' trimmed, plus "\n;".
//
// Only the right-most SELECT of the outermost query is touched (the last
// operand of a set operation, the query of a result pipe):
//
//   - LIMIT n, FETCH n and TOP n become min(n, limit); n = 0 becomes limit.
//   - With no limiting clause, " LIMIT <limit>" is appended at the end of that
//     SELECT, inside any wrapping parentheses.
//
// A statement that is not a query is returned trimmed and unchanged. An error
// (a parse error, a non-integer count, OFFSET without LIMIT or FETCH) means
// the caller should wrap the statement instead.
func RewriteResultLimit(statement string, limit int) (string, error) {
	trimmed := strings.TrimRightFunc(statement, isSpaceOrSemicolon)

	file, err := parser.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("statement: %s: %w", statement, err)
	}
	if len(file.Stmts) != 1 {
		return "", fmt.Errorf("expected exactly 1 statement, got %d", len(file.Stmts))
	}

	sel := findRightMostSelect(file.Stmts[0])
	if sel == nil {
		return trimmed + "\n;", nil
	}

	if sel.Limit != nil {
		spliced, err := spliceLimitNumber(trimmed, ast.NodeLoc(sel.Limit), limit)
		if err != nil {
			return "", fmt.Errorf("statement: %s: %w", statement, err)
		}
		return spliced + "\n;", nil
	}

	if sel.Fetch != nil {
		spliced, err := spliceLimitNumber(trimmed, ast.NodeLoc(sel.Fetch.Count), limit)
		if err != nil {
			return "", fmt.Errorf("statement: %s: %w", statement, err)
		}
		return spliced + "\n;", nil
	}

	// LIMIT must precede OFFSET, so a bare OFFSET cannot take an appended LIMIT.
	if sel.Offset != nil {
		return "", fmt.Errorf("statement has OFFSET without LIMIT/FETCH: %s", statement)
	}

	if sel.Top != nil {
		spliced, err := spliceLimitNumber(trimmed, ast.NodeLoc(sel.Top), limit)
		if err != nil {
			return "", fmt.Errorf("statement: %s: %w", statement, err)
		}
		return spliced + "\n;", nil
	}

	end := sel.Loc.End
	if !sel.Loc.IsValid() || end > len(trimmed) {
		return "", fmt.Errorf("invalid select location %v for statement: %s", sel.Loc, statement)
	}
	return trimmed[:end] + fmt.Sprintf(" LIMIT %d", limit) + trimmed[end:] + "\n;", nil
}

// spliceLimitNumber replaces the integer at loc by min(n, limit), where 0
// means no limit and becomes limit. A non-integer at loc is an error.
func spliceLimitNumber(text string, loc ast.Loc, limitCount int) (string, error) {
	if !loc.IsValid() || loc.End > len(text) || loc.Start > loc.End {
		return "", fmt.Errorf("invalid limit number location %v", loc)
	}
	userLimit, err := strconv.Atoi(text[loc.Start:loc.End])
	if err != nil {
		return "", fmt.Errorf("non-integer limit count %q: %w", text[loc.Start:loc.End], err)
	}
	limit := userLimit
	if limit == 0 || limitCount < limit {
		limit = limitCount
	}
	return text[:loc.Start] + strconv.Itoa(limit) + text[loc.End:], nil
}

// findRightMostSelect returns the right-most SELECT of the outermost query, or
// nil for a statement that is not a query.
func findRightMostSelect(node ast.Node) *ast.SelectStmt {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n
	case *ast.SetOperationStmt:
		return findRightMostSelect(n.Right)
	case *ast.ResultScanStmt:
		return findRightMostSelect(n.Query)
	case *ast.ShowStmt:
		if n.Pipe != nil {
			return findRightMostSelect(n.Pipe)
		}
		return nil
	default:
		return nil
	}
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
