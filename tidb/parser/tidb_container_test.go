package parser

import (
	"context"
	"database/sql"
	"errors"
	"iter"
	"os"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// tidbRelease is the TiDB release the package's tapes and compat.tsv were
// recorded on, as tidb/testdata/pins.tsv pins it. Messages name it.
const tidbRelease = "v8.5.5"

// TestMain runs the tests through the tape harness, and removes the server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// tidbContainer is a test's connection to TiDB: a replay of its tape, or the
// pinned server in a recording (testcontainer.Open). The pool holds one
// connection, at session defaults with the database test selected, so USE and
// PREPARE hold across statements.
type tidbContainer struct {
	db  *sql.DB
	ctx context.Context
}

func startTiDB(t testing.TB) *tidbContainer {
	t.Helper()
	return &tidbContainer{db: testcontainer.Open(t), ctx: context.Background()}
}

// canExecute runs the statement and returns the server's error.
func (tc *tidbContainer) canExecute(sqlStr string) error {
	_, err := tc.db.ExecContext(tc.ctx, sqlStr)
	return err
}

// prepareSyntax reports whether TiDB accepts the statement's syntax, without
// running it: PREPARE parses the text and stops there. A non-nil error means
// the server gave no verdict at all: anything that is not a server error.
//
// ER_PARSE_ERROR (1064) is the parser's rejection, and it is not the only one: a
// grammar action refuses text its productions read, with the code it chooses, while
// the parse is still running. CreateIndexStmt's answers 1149 for a VECTOR INDEX over
// a plain column, whether or not the table exists. Other errors come after the parse:
// 1295 for a statement the prepared-statement protocol does not carry, 1146 for a
// table that does not exist. The code alone cannot tell the two apart (1105 is both),
// so the verdict asks the server which it was: see prepareVerdict.
//
// The statement is the one argument of one request, which the connection
// inlines as an escaped literal, and a PREPARE under the same name replaces the
// one before it, so each verdict is one exchange on the tape (two when it
// needs the probe).
func (tc *tidbContainer) prepareSyntax(stmt string) (bool, error) {
	prepare := func(text string) error {
		_, err := tc.db.ExecContext(tc.ctx, "PREPARE omni_stmt FROM ?", text)
		return err
	}
	return prepareVerdict(prepare(stmt), func() error { return prepare(stmt + compat.ProbeTail) })
}

// prepareVerdict reads PREPARE's answer through compat.PrepareVerdict, which
// says how the probe tells a grammar action's refusal from an error raised after
// the parse.
func prepareVerdict(err error, probe func() error) (bool, error) {
	return compat.PrepareVerdict(err, mysqlErrorCode, probe)
}

// mysqlErrorCode returns the server's error number, and false for an error the
// server did not send.
func mysqlErrorCode(err error) (uint16, bool) {
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		return 0, false
	}
	return myErr.Number, true
}

// byVerdict ranges over m's accepted cases, then its rejected ones, so what a
// test sends does not depend on map order.
func byVerdict[T any](m map[bool]T) iter.Seq2[bool, T] {
	return func(yield func(bool, T) bool) {
		for _, want := range []bool{true, false} {
			if v, ok := m[want]; ok && !yield(want, v) {
				return
			}
		}
	}
}

func verdictName(accepts bool) string {
	if accepts {
		return "accept"
	}
	return "reject"
}

// TestSyntaxProbeControls holds the probe to statements whose answer is known from
// the grammar and the server's own order of work. Each pair differs in nothing but
// the shape a grammar action looks at, on a table that does not exist.
func TestSyntaxProbeControls(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range []struct {
		stmt string
		want bool
	}{
		{"CREATE VECTOR INDEX i ON nosuch (a)", false},  // 1149, CreateIndexStmt's action
		{"CREATE VECTOR INDEX i ON nosuch ((a))", true}, // 1146
		{"SELECT nosuch", true},                         // 1054
		{"SELECT 1 FROM nosuch -- a comment", true},     // 1146, and the tail survives the comment
		{"CREATE TABLE t (a VECTOR(16384))", true},      // 1105, from the DDL, not the parser
		{"SELECT @@nosuch", true},                       // 1193, a name it looks up
		{"SELEKT 1", false},                             // 1064
		{"SELECT DISTINCT ALL * FROM nosuch", false},    // 1221, the action on SelectStmtOpts
		{"ALTER TABLE t ALGORITHM = nosuch", false},     // 1800
	} {
		accepts, err := tc.prepareSyntax(c.stmt)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", c.stmt, err)
		}
		if accepts != c.want {
			t.Errorf("%q: accepts=%v, want %v", c.stmt, accepts, c.want)
		}
	}
}
