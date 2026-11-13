package review

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestMain accounts for the package's tapes and removes the PostgreSQL
// servers a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// reviewDatabase is the database the oracle tests ask, which cases name in
// their text.
const reviewDatabase = "omni_review"

// startPostgres returns t's connection to omni_review on PostgreSQL 17: a
// replay of its tape, or the pinned server in a recording, which starts the
// test on a server holding no other test's databases or roles
// (pg/internal/testcontainer).
func startPostgres(t *testing.T) *sql.DB {
	t.Helper()
	return testcontainer.OpenDatabase(t, testcontainer.Default, reviewDatabase)
}

// startPostgresRole is startPostgres on the server of another pin row, for
// the cases a version decides.
func startPostgresRole(t *testing.T, role string) *sql.DB {
	t.Helper()
	return testcontainer.OpenDatabase(t, role, reviewDatabase)
}

// prepared prepares query on conn as name, a statement taking parameters of
// the types params lists, and returns the EXECUTE that runs it, with a $n for
// each parameter, and a func that deallocates it. A query a test sends
// hundreds of times reaches its tape once that way, and a short EXECUTE each
// time after.
func prepared(t testing.TB, conn *sql.Conn, name string, params []string, query string) (string, func()) {
	t.Helper()
	ctx := context.Background()
	execute, head := "EXECUTE "+name, "PREPARE "+name
	if len(params) > 0 {
		head += "(" + strings.Join(params, ", ") + ")"
		var args []string
		for i := range params {
			args = append(args, fmt.Sprintf("$%d", i+1))
		}
		execute += "(" + strings.Join(args, ", ") + ")"
	}
	if _, err := conn.ExecContext(ctx, head+" AS "+query); err != nil {
		t.Fatalf("%s: %v", head, err)
	}
	return execute, func() {
		if _, err := conn.ExecContext(ctx, "DEALLOCATE "+name); err != nil {
			t.Errorf("DEALLOCATE %s: %v", name, err)
		}
	}
}
