package parser

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/internal/sharedserver"
)

// parserOracle wraps a MySQL 8.0 container for syntax-level oracle testing.
type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startParserOracle returns the MySQL 8.0 oracle with a database of this test's
// selected, dropped when the test ends. A statement it runs can change nothing
// outside that database: see sharedserver.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	ctx := context.Background()
	db, err := sharedserver.MySQL80.Database(ctx, "tidb_parser")
	if err != nil {
		t.Fatalf("MySQL 8.0 oracle: %v", err)
	}
	t.Cleanup(func() { db.Close(ctx) })
	return &parserOracle{db: db.DB, ctx: ctx}
}

// checkSyntax reports whether the oracle parses sql. It runs the statement and
// reads the error: 1064 is a parse error, anything else means the statement parsed
// and failed for a semantic reason, such as a missing table.
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
		return myErr.Number != 1064
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
