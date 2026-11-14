package analysis

import (
	"errors"
	"fmt"
	"strings"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// RewriteResultLimit rewrites one SELECT so it returns at most limitCount
// rows, splicing into the original text by node offsets. A statement that is
// not a SELECT is returned unchanged.
//
//   - A non-negative integer LIMIT at most limitCount is kept; any other integer becomes
//     limitCount.
//   - With no LIMIT, "LIMIT <limitCount>" goes before a FOR UPDATE/SHARE
//     clause, or after the statement's last clause.
//
// An error (empty input, a parse error, more than one statement, a LIMIT ALL,
// NULL or expression) means the caller should wrap the statement instead.
func RewriteResultLimit(statement string, limitCount int) (string, error) {
	if strings.TrimSpace(statement) == "" {
		return "", errors.New("empty statement")
	}

	stmts, err := omnipg.Parse(statement)
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

	return rewriteSelectLimit(statement, sel, limitCount)
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
		if isInteger {
			if end, ok := integerLiteralEnd(sql, loc); ok {
				return sql[:loc.Start] + fmt.Sprintf("%d", limitCount) + sql[end:], nil
			}
		}
		// ALL/NULL and non-constant limits lack reliable replacement spans.
		// Cap them with the outer query instead of rewriting their text.
		return "", errors.New("cannot rewrite LIMIT expression in place")
	}

	// No LIMIT clause — find the right insertion point.
	// PostgreSQL grammar order: ... ORDER BY ... LIMIT ... FOR UPDATE ...
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

// integerLiteralEnd returns the exclusive end byte offset of the integer
// literal. The AST span may include whitespace before the next token, which
// must remain when replacing the literal.
func integerLiteralEnd(sql string, loc ast.Loc) (int, bool) {
	if loc.Start < 0 || loc.End <= loc.Start || loc.End > len(sql) {
		return 0, false
	}
	literal := strings.TrimRight(sql[loc.Start:loc.End], " \t\r\n\f\v")
	return loc.Start + len(literal), true
}
