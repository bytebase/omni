package parser

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startMariaDB returns the MariaDB oracle with an empty `test` database
// selected, as a server of its own would be: a replay of the test's tape, or
// the pinned server in a recording (see testcontainer.OpenEmpty). It connects
// as root, so a case it runs must not change the server under the tests that
// follow in a way OpenEmpty does not undo: no SET GLOBAL, no KILL.
//
// The server aborts a slow or locked statement itself (max_statement_time,
// lock_wait_timeout), so the connection stays usable after one.
func startMariaDB(t *testing.T) *parserOracle {
	t.Helper()
	db := testcontainer.OpenEmptyWith(t, testcontainer.Options{
		Database:        "test",
		MultiStatements: true,
		ParseTime:       true,
		Session:         map[string]string{"max_statement_time": "5", "lock_wait_timeout": "2"},
	})
	return &parserOracle{db: db, ctx: context.Background()}
}

// TestMain runs the tests through the tape harness, and removes the server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
