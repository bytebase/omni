package analysis

import (
	"errors"
	"fmt"
	"strings"

	omniast "github.com/bytebase/omni/tidb/ast"
	omnideparse "github.com/bytebase/omni/tidb/deparse"
	"github.com/bytebase/omni/tidb/parser"
)

// StatementWithResultLimit returns the text to send so the server returns at
// most limit rows. It is RewriteResultLimit's rewrite when that succeeds, and
// otherwise statement wrapped as
//
//	WITH result AS (<statement>) SELECT * FROM result LIMIT <limit>;
//
// It never reports an error, so a caller has no fallback to choose.
//
// A statement RewriteResultLimit leaves alone — DML, EXPLAIN, and other
// non-SELECT text — is returned unchanged: wrapping it would change what it
// does. A statement that already carries an integer LIMIT at or below limit is
// returned unchanged; a larger LIMIT becomes limit. TiDB has no version without
// LIMIT, so the wrapper covers only text the rewrite cannot handle, such as a
// nested locking clause or a parse error; a LIMIT 0 statement is kept, and an
// empty statement is wrapped.
func StatementWithResultLimit(statement string, limitCount int) string {
	if stmt, err := RewriteResultLimit(statement, limitCount); err == nil {
		return stmt
	}
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", trimStatement(statement), limitCount)
}

// RewriteResultLimit rewrites one SELECT, including a set operation, so it
// returns at most limit rows, and restores it from the AST. An integer LIMIT
// larger than limit becomes limit, a non-constant one is kept, and a SELECT
// with no LIMIT gets one. A statement that is not a SELECT is returned
// unchanged.
//
// The restore writes no locking clause, so the outermost SELECT's own one is
// written back after it. A statement with a locking clause anywhere else
// returns errNestedLock.
//
// An error (a parse error, other than one statement, a nested locking clause)
// means the caller should wrap the statement instead.
func RewriteResultLimit(statement string, limit int) (string, error) {
	list, err := parser.Parse(statement)
	if err != nil {
		return "", fmt.Errorf("failed to parse tidb statement: %s: %w", statement, err)
	}
	if list == nil || len(list.Items) != 1 {
		return "", fmt.Errorf("expect one single statement in the query, %s", statement)
	}
	// omni folds UNION/INTERSECT/EXCEPT into SelectStmt, so a single case covers
	// both a plain SELECT and a set operation; the trailing LIMIT lives on the
	// outermost SelectStmt either way.
	stmt, ok := list.Items[0].(*omniast.SelectStmt)
	if !ok {
		return statement, nil
	}
	if hasNestedLock(stmt) {
		return "", errNestedLock
	}
	if stmt.Limit != nil && stmt.Limit.Count != nil {
		if v, ok := stmt.Limit.Count.(*omniast.IntLit); ok {
			userLimit := int(v.Value)
			if limit < userLimit {
				userLimit = limit
			}
			stmt.Limit.Count = &omniast.IntLit{Value: int64(userLimit)}
		}
	} else {
		stmt.Limit = &omniast.Limit{
			Count: &omniast.IntLit{Value: int64(limit)},
		}
	}
	// String literals are restored without a charset introducer.
	restored := omnideparse.RestoreSelectWith(stmt, omnideparse.RestoreStringWithoutCharset)
	if stmt.ForUpdate != nil {
		// The restore ends in the space before the clause it leaves out.
		restored += lockingClause(stmt.ForUpdate)
	}
	return restored, nil
}

// errNestedLock is RewriteResultLimit's error for a locking clause it cannot
// write back. TiDB takes the lock wherever the clause is, in a subquery, a CTE,
// a set operation's operand or parentheses, and keeps it when the caller wraps
// the statement.
var errNestedLock = errors.New("cannot restore a locking clause below the outermost SELECT")

// hasNestedLock reports whether a query block other than stmt itself has a
// locking clause. A set operation or a parenthesized query has none of its
// own: the clause belongs to an operand.
func hasNestedLock(stmt *omniast.SelectStmt) bool {
	nested := false
	omniast.Inspect(stmt, func(n omniast.Node) bool {
		if lock, ok := n.(*omniast.ForUpdate); ok && lock != stmt.ForUpdate {
			nested = true
		}
		return !nested
	})
	return nested || (stmt.ForUpdate != nil && (stmt.SetOp != omniast.SetOpNone || stmt.ParenSource != nil))
}

// lockingClause writes a locking clause as TiDB's SelectLockInfo.Restore does.
func lockingClause(lock *omniast.ForUpdate) string {
	if lock.LockInShareMode {
		return "LOCK IN SHARE MODE"
	}
	var b strings.Builder
	if lock.Share {
		b.WriteString("FOR SHARE")
	} else {
		b.WriteString("FOR UPDATE")
	}
	for i, table := range lock.Tables {
		if i == 0 {
			b.WriteString(" OF ")
		} else {
			b.WriteString(", ")
		}
		if table.Schema != "" {
			b.WriteString(quoteIdentifier(table.Schema))
			b.WriteByte('.')
		}
		b.WriteString(quoteIdentifier(table.Name))
	}
	switch {
	case lock.NoWait:
		b.WriteString(" NOWAIT")
	case lock.HasWait:
		b.WriteString(" WAIT " + lock.WaitSeconds)
	case lock.SkipLocked:
		b.WriteString(" SKIP LOCKED")
	}
	return b.String()
}
