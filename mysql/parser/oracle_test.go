package parser

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startParserOracle returns the oracle with its database empty and selected, as a
// server of its own would be: a replay of the test's tape, or the pinned server in
// a recording (see testcontainer.OpenDatabase). A statement it runs can change
// nothing outside that database: see testcontainer.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	return &parserOracle{db: testcontainer.OpenDatabase(t, "mysql_parser"), ctx: context.Background()}
}

// TestMain runs the tests through the tape harness, and removes the server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// checkSyntax reports whether the oracle parses sql. It runs the statement and
// reads the error: 1064 is a parse error, and anything else goes to prepareParses, which
// tells a code raised while parsing from one raised for a semantic reason after it, such
// as a missing table.
//
// Anything that is not a server error — a dropped connection, a timeout — is no
// verdict at all and fails the test. Reporting it as acceptance would pass every
// case that expects the engine to accept, against a server that never answered.
func (o *parserOracle) checkSyntax(t *testing.T, sql string) bool {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, sql)
	if err == nil {
		return true
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number != 1064 && o.prepareParses(t, sql)
	}
	t.Fatalf("the oracle gave no verdict for %q: %v", sql, err)
	return false
}

// assertOracleMatch parses with omni parser AND checks against MySQL 8.0,
// asserting both accept or both reject.
func assertOracleMatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	mysqlAccepts := o.checkSyntax(t, sql)

	if omniAccepts != mysqlAccepts {
		t.Errorf("MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  mysql accepts=%v",
			sql, omniAccepts, omniErr, mysqlAccepts)
	} else {
		t.Logf("MATCH for %q: both accept=%v", sql, omniAccepts)
	}
}

// assertOracleMismatch logs a known mismatch between omni and MySQL 8.0 without
// failing the test. These are tracked as [~] partial in the SCENARIOS file.
func assertOracleMismatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	mysqlAccepts := o.checkSyntax(t, sql)

	if omniAccepts != mysqlAccepts {
		t.Logf("KNOWN MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  mysql accepts=%v",
			sql, omniAccepts, omniErr, mysqlAccepts)
	} else {
		// If the mismatch is resolved (e.g. by a parser fix), promote to match.
		t.Logf("RESOLVED — now MATCH for %q: both accept=%v (was known mismatch)", sql, omniAccepts)
	}
}
