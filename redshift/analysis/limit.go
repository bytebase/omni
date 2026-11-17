package analysis

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
)

// StatementWithResultLimit returns the text to send so the server returns at
// most limit rows. It is RewriteResultLimit's rewrite when that succeeds, and
// otherwise statement wrapped as
//
//	WITH result AS (
//	<statement>
//	) SELECT * FROM result LIMIT <limit>;
//
// It never reports an error, so a caller has no fallback to choose.
func StatementWithResultLimit(statement string, limit int) string {
	stmt, err := RewriteResultLimit(statement, limit)
	if err != nil {
		// To handle cases where there are comments in the query.
		// eg. select * from t1 -- this is comment;
		// Add two new line symbols here.
		return fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", trimStatement(statement), limit)
	}
	return stmt
}

// RewriteResultLimit rewrites one SELECT so it returns at most limit rows,
// splicing into the original text by node offsets. A statement that is not a
// SELECT is returned unchanged.
//
//   - An integer LIMIT at most limit is kept; a larger one becomes limit.
//   - With no LIMIT, "LIMIT <limit>" goes before a trailing FOR UPDATE/SHARE
//     clause, or after everything else.
//
// An error (empty input, a parse error, more than one statement, a LIMIT
// expression that is not an integer) means the caller should wrap the
// statement instead.
func RewriteResultLimit(statement string, limit int) (string, error) {
	if strings.TrimSpace(statement) == "" {
		return "", errors.New("empty statement")
	}

	stmts, err := redshift.Parse(statement)
	if err != nil {
		return "", fmt.Errorf("failed to parse statement: %w", err)
	}
	if len(stmts) != 1 {
		return "", fmt.Errorf("expected exactly one statement, got %d", len(stmts))
	}

	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		// Non-SELECT statement, return as-is.
		return statement, nil
	}

	return rewriteSelectLimit(statement, sel, limit)
}

// rewriteSelectLimit adds or adjusts the LIMIT clause of a SELECT statement
// using byte-offset positions from the omni AST to surgically edit the original SQL text.
func rewriteSelectLimit(sql string, sel *ast.SelectStmt, limitCount int) (string, error) {
	if sel.LimitCount != nil {
		existingLimit, isInteger := extractIntFromNode(sel.LimitCount)
		loc := nodeLocOf(sel.LimitCount)
		if isInteger && existingLimit >= 0 && existingLimit <= limitCount {
			return sql, nil // existing limit is already lower or equal, keep it
		}
		if isInteger && loc.Start >= 0 && loc.End > loc.Start && loc.End <= len(sql) {
			return sql[:loc.Start] + fmt.Sprintf("%d", limitCount) + sql[loc.End:], nil
		}
		// ALL/NULL and non-constant limits lack reliable replacement spans.
		// Cap them with the outer query instead of rewriting their text.
		return "", errors.New("cannot rewrite LIMIT expression in place")
	}

	// No LIMIT clause — find the right insertion point.
	// Redshift's PostgreSQL-derived grammar order: ... ORDER BY ... LIMIT ... FOR UPDATE ...
	// LIMIT goes BEFORE FOR UPDATE but AFTER everything else.
	insertPos, beforeLocking := findLimitInsertPosition(sel)
	if beforeLocking {
		// Inserting at the start of FOR UPDATE/SHARE. The original whitespace
		// before FOR becomes the separator before LIMIT; we add a trailing
		// space to separate the limit value from FOR.
		return sql[:insertPos] + fmt.Sprintf("LIMIT %d ", limitCount) + sql[insertPos:], nil
	}
	return sql[:insertPos] + fmt.Sprintf(" LIMIT %d", limitCount) + sql[insertPos:], nil
}

// findLimitInsertPosition returns the byte offset where " LIMIT N" should be inserted,
// and whether the insertion is before a locking clause (FOR UPDATE/SHARE).
func findLimitInsertPosition(sel *ast.SelectStmt) (int, bool) {
	// LIMIT must appear before FOR UPDATE/SHARE.
	if sel.LockingClause != nil {
		span := ast.ListSpan(sel.LockingClause)
		if span.Start > 0 {
			return span.Start, true
		}
	}

	// Otherwise insert at SelectStmt.Loc.End (after everything, including outer parens).
	// For CTE queries, omni may report SelectStmt.Loc.End before the outer ORDER BY,
	// so account for the sort clause explicitly.
	end := sel.Loc.End
	if span := ast.ListSpan(sel.SortClause); span.End > end {
		end = span.End
	}
	if end <= 0 {
		return 0, false
	}
	return end, false
}

// extractIntFromNode distinguishes integer limits, including zero, from unlimited
// ALL/NULL constants and expressions that require the fallback wrapper.
func extractIntFromNode(node ast.Node) (int, bool) {
	switch n := node.(type) {
	case *ast.Integer:
		return int(n.Ival), true
	case *ast.A_Const:
		if iv, ok := n.Val.(*ast.Integer); ok {
			return int(iv.Ival), true
		}
		return 0, false
	default:
		return 0, false
	}
}

// nodeLocOf returns the Loc of a node, handling common wrapper types.
func nodeLocOf(node ast.Node) ast.Loc {
	switch n := node.(type) {
	case *ast.A_Const:
		return n.Loc
	default:
		return ast.Loc{Start: -1, End: -1}
	}
}

// trimStatement returns statement without leading space and trailing space or
// semicolons.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}
