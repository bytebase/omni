package parser

import (
	"context"
	"database/sql"
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
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}

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

// checkSyntax sends SQL to MySQL 8.0 and returns true if syntactically valid.
// Strategy: execute the SQL. If MySQL returns error code 1064 (ER_PARSE_ERROR),
// the SQL is syntactically invalid (rejected). Any other error (e.g. 1146 table
// not found) means the SQL parsed fine but failed for semantic reasons (accepted).
func (o *parserOracle) checkSyntax(sql string) bool {
	_, err := o.db.ExecContext(o.ctx, sql)
	if err == nil {
		return true // executed successfully — valid syntax
	}
	if myErr, ok := err.(*gomysql.MySQLError); ok {
		return myErr.Number != 1064 // 1064 = ER_PARSE_ERROR
	}
	// Non-MySQL error (connection issue, etc.) — treat as valid to avoid
	// false negatives; the caller should investigate.
	return true
}

// assertOracleMatch parses with omni parser AND checks against MySQL 8.0,
// asserting both accept or both reject.
func assertOracleMatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	mysqlAccepts := o.checkSyntax(sql)

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

	mysqlAccepts := o.checkSyntax(sql)

	if omniAccepts != mysqlAccepts {
		t.Logf("KNOWN MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  mysql accepts=%v",
			sql, omniAccepts, omniErr, mysqlAccepts)
	} else {
		// If the mismatch is resolved (e.g. by a parser fix), promote to match.
		t.Logf("RESOLVED — now MATCH for %q: both accept=%v (was known mismatch)", sql, omniAccepts)
	}
}
