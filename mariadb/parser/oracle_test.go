package parser

import (
	"context"
	"database/sql"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// parserOracle wraps a SQL container for syntax-level oracle testing.
type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startParserOracle starts MySQL 8.0 via testcontainers. Used by the routine
// alignment oracle (routine_alignment_test.go). The MariaDB 11.8 divergence
// inventory uses startMariaDB (oracle_corpus_test.go) instead — repointing the
// routine oracle to MariaDB is follow-up (P1+) work, out of PR1 scope.
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
