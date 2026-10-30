package parser

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// parserOracle wraps a MySQL 8.0 container for syntax-level oracle testing.
type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startParserOracle starts MySQL 8.0 via testcontainers and returns a
// parserOracle. The container is cleaned up automatically when the test ends.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()

	ctx := context.Background()

	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	connStr, err := container.ConnectionString(ctx, "parseTime=true", "multiStatements=true")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	return &parserOracle{db: db, ctx: ctx}
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
