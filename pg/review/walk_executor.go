package review

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// Bytebase's PostgreSQL executor does not always run a change in the order
// it is written. It runs the statements PostgreSQL refuses inside a
// transaction block after its transaction commits, or, in auto-commit
// mode, before every other statement; a SET ROLE it runs twice, in place
// and with them. What it moves it decides by matching each statement's
// text, as its splitter cuts it, against the patterns below, and a
// directive at the top of the change chooses the mode. The walk follows
// the change in the order Bytebase runs it.
//
// bytebase: backend/plugin/db/pg/pg.go — Execute, executeInTransactionMode,
// executeInAutoCommitMode, IsNonTransactionStatement, isSetRoleStatement,
// isSuperuserStatement; backend/plugin/parser/base/transaction_mode.go —
// ParseTransactionConfig; backend/plugin/db/transaction/config.go —
// DefaultMode (transaction mode).
var (
	bbDropDatabase = regexp.MustCompile(`(?i)DROP DATABASE`)
	bbCreateIndex  = regexp.MustCompile(`(?i)CREATE(\s+(UNIQUE\s+)?)INDEX(\s+)CONCURRENTLY`)
	bbDropIndex    = regexp.MustCompile(`(?i)DROP(\s+)INDEX(\s+)CONCURRENTLY`)
	bbVacuum       = regexp.MustCompile(`(?i)^\s*VACUUM`)
	bbSetRole      = regexp.MustCompile(`(?i)SET\s+((SESSION|LOCAL)\s+)?ROLE`)
	bbTxnMode      = regexp.MustCompile(`(?i)^\s*--\s*txn-mode\s*=\s*(on|off)\s*$`)
	bbTxnIsolation = regexp.MustCompile(`(?i)^\s*--\s*txn-isolation\s*=\s*(READ\s+UNCOMMITTED|READ\s+COMMITTED|REPEATABLE\s+READ|SERIALIZABLE)\s*$`)
)

// bbSuperuserPrefixes begin the statements Bytebase runs as the session
// user in tenant mode, between SET LOCAL ROLE NONE and SET LOCAL ROLE of
// the database owner.
var bbSuperuserPrefixes = []string{"GRANT", "CREATE EXTENSION", "CREATE EVENT TRIGGER", "COMMENT ON EVENT TRIGGER", "COMMENT ON EXTENSION"}

// run is a change in the order Bytebase runs it.
type run struct {
	// autocommit is true when the change asks for auto-commit mode: each
	// statement commits by itself. Otherwise Bytebase runs the statements
	// it does not move in one transaction it opens.
	autocommit bool
	steps      []step
}

// step is one statement as Bytebase runs it.
type step struct {
	// stmt is the statement's index in the change.
	stmt int
	// again marks the second run of a statement Bytebase runs twice.
	again bool
	// committed marks the first statement Bytebase runs after it commits
	// its transaction, and then, in tenant mode, sets the session's role
	// to the database owner.
	committed bool
	// wrapped marks a statement Bytebase runs, in tenant mode, as the
	// session user, and after which it sets the role back to the database
	// owner until its transaction ends.
	wrapped bool
}

// planRun returns the order Bytebase runs the change in, or nil when the
// text it runs does not split into the statements the review parsed.
func planRun(sql string, stmts []statement) *run {
	autocommit, text := bytebaseScript(sql)
	var segs []string
	for _, seg := range pg.Split(text) {
		if !seg.Empty() {
			segs = append(segs, seg.Text)
		}
	}
	if len(segs) != len(statementRanges(sql)) {
		return nil
	}
	setRole := make([]bool, len(stmts))
	moved := make([]bool, len(stmts))
	wrapped := make([]bool, len(stmts))
	for i := range stmts {
		seg := segs[stmts[i].index]
		switch {
		case bbSetRole.MatchString(seg):
			setRole[i] = true
		case bbDropDatabase.MatchString(seg) || bbCreateIndex.MatchString(seg) || bbDropIndex.MatchString(seg) || bbVacuum.MatchString(seg):
			moved[i] = true
		default:
			upper := strings.ToUpper(strings.TrimLeftFunc(seg, unicode.IsSpace))
			for _, prefix := range bbSuperuserPrefixes {
				if strings.HasPrefix(upper, prefix) {
					wrapped[i] = true
				}
			}
		}
	}
	r := &run{autocommit: autocommit}
	inPlace := func(again bool) {
		for i := range stmts {
			if !moved[i] {
				r.steps = append(r.steps, step{stmt: i, again: again && setRole[i], wrapped: wrapped[i]})
			}
		}
	}
	apart := func(again, committed bool) {
		for i := range stmts {
			if setRole[i] || moved[i] {
				r.steps = append(r.steps, step{stmt: i, again: again && setRole[i], committed: committed})
				committed = false
			}
		}
	}
	if autocommit {
		apart(false, false)
		inPlace(true)
	} else {
		inPlace(false)
		apart(true, true)
	}
	return r
}

// bytebaseScript reads the directives at the top of a change as Bytebase
// does: it returns whether the change asks for auto-commit mode, and the
// text Bytebase runs, which leaves the directive lines out.
func bytebaseScript(sql string) (autocommit bool, text string) {
	lines := strings.Split(sql, "\n")
	directive := make(map[int]bool)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "--") {
			break
		}
		if m := bbTxnMode.FindStringSubmatch(line); len(m) == 2 {
			autocommit = strings.EqualFold(m[1], "off")
			directive[i] = true
			continue
		}
		if bbTxnIsolation.MatchString(line) {
			directive[i] = true
		}
	}
	if len(directive) == 0 {
		return autocommit, sql
	}
	var kept []string
	for i, line := range lines {
		if !directive[i] {
			kept = append(kept, line)
		}
	}
	return autocommit, strings.TrimSpace(strings.Join(kept, "\n"))
}

// refusedInBlock reports whether PostgreSQL may refuse the statement inside
// a transaction block (25001), as it does CONCURRENTLY, VACUUM, and some
// forms of REINDEX and CLUSTER. Bytebase runs one its patterns miss, such
// as CREATE /* ... */ INDEX CONCURRENTLY, inside its transaction, where it
// fails; the walk ends before it, as a fix of it may move it.
//
// pg: src/backend/tcop/utility.c — PreventInTransactionBlock callers
func refusedInBlock(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.IndexStmt:
		return v.Concurrent
	case *ast.DropStmt:
		return v.Concurrent
	case *ast.VacuumStmt:
		return v.IsVacuumCmd
	case *ast.ReindexStmt, *ast.ClusterStmt:
		return true
	}
	return false
}
