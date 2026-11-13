package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of PostgreSQL's gram.y and its near-misses, and the suites
// testdata/README.md lists, with omni's verdict and PostgreSQL's. A
// rejection carries the byte it is reported at. The test reads the file and
// starts no server. OMNI_RECORD=1 judges the rows that have no verdict yet
// on the pinned PostgreSQL.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join("..", "testdata", "gram.y"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict is Parse's answer as a compat.tsv cell: accept, or reject@ and
// the byte offset of the error. PostgreSQL reports no position for a few
// errors; ParseError has none to give then either (-1).
func omniVerdict(sql string) string {
	_, err := Parse(sql)
	if err == nil {
		return compat.Accept
	}
	var pe *ParseError
	if errors.As(err, &pe) && pe.Position >= 0 {
		return compat.RejectAt(pe.Position)
	}
	return compat.Reject
}

// compatJudge asks the pinned PostgreSQL whether its raw parser accepts each
// statement, with the probe compat.RawParseVerdict reads: the statement after
// SELECT 1/0 in one simple query, which the server parses whole before it
// runs any of it. Nothing of the statement runs, so the judge needs no
// transaction and leaves nothing behind. It runs only in a recording, on a
// connection outside any tape.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	db := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		// The protocol sends a query as a C string: the server would read a
		// statement only up to its first zero byte.
		if strings.ContainsRune(sql, 0) {
			return "", fmt.Errorf("the statement holds a zero byte, which the protocol cannot carry")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err := db.ExecContext(ctx, compat.RawParseHead+sql)
		return compat.RawParseVerdict(sql, err, pgState)
	}
}

// pgState reads a server error for compat.RawParseVerdict.
func pgState(err error) (sqlstate string, position int, ok bool) {
	var pe *pgconn.PgError
	if !errors.As(err, &pe) {
		return "", 0, false
	}
	return pe.Code, int(pe.Position), true
}
