package parser

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/testcontainers/testcontainers-go"
	tcmssql "github.com/testcontainers/testcontainers-go/modules/mssql"
)

// parserOracle wraps a SQL Server container for syntax-level oracle testing.
type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startParserOracle starts SQL Server 2022 via testcontainers and returns a
// parserOracle. The container is cleaned up automatically when the test ends.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}

	ctx := context.Background()

	container, err := tcmssql.Run(ctx, "mcr.microsoft.com/mssql/server:2022-latest",
		tcmssql.WithAcceptEULA(),
		tcmssql.WithPassword("Str0ngPa$$w0rd!"),
	)
	if err != nil {
		t.Fatalf("failed to start SQL Server container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping SQL Server: %v", err)
	}

	// Create a test database and table for queries that reference objects
	setupSQL := []string{
		"CREATE DATABASE testdb",
		"USE testdb",
		"CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)",
	}
	for _, s := range setupSQL {
		if _, err := db.ExecContext(ctx, s); err != nil {
			t.Fatalf("setup SQL failed (%s): %v", s, err)
		}
	}

	return &parserOracle{db: db, ctx: ctx}
}

// canParse tests whether SQL Server accepts the given SQL without execution errors.
// It uses SET PARSEONLY ON to check syntax without executing.
func (o *parserOracle) canParse(sql string) (bool, error) {
	_, err := o.db.ExecContext(o.ctx, "SET PARSEONLY ON")
	if err != nil {
		return false, fmt.Errorf("SET PARSEONLY ON: %w", err)
	}
	defer o.db.ExecContext(o.ctx, "SET PARSEONLY OFF") //nolint:errcheck

	_, err = o.db.ExecContext(o.ctx, sql)
	if err != nil {
		return false, nil // Parse error — SQL Server rejects this syntax
	}
	return true, nil // SQL Server accepts this syntax
}

// TestKeywordOracleOptionPositions verifies whether omni's option parsing
// is too permissive compared to SQL Server.
func TestKeywordOracleOptionPositions(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		name string
		sql  string
	}{
		// Index options — valid
		{"index FILLFACTOR valid", "CREATE INDEX ix ON dbo.t(a) WITH (FILLFACTOR = 80)"},
		{"index PAD_INDEX valid", "CREATE INDEX ix ON dbo.t(a) WITH (PAD_INDEX = ON, FILLFACTOR = 80)"},
		{"index IGNORE_DUP_KEY valid", "CREATE INDEX ix ON dbo.t(a) WITH (IGNORE_DUP_KEY = OFF)"},
		{"index STATISTICS_NORECOMPUTE valid", "CREATE INDEX ix ON dbo.t(a) WITH (STATISTICS_NORECOMPUTE = ON)"},
		// Index options — invalid (core keywords as option names)
		{"index SELECT invalid", "CREATE INDEX ix ON dbo.t(a) WITH (SELECT = 1)"},
		{"index FROM invalid", "CREATE INDEX ix ON dbo.t(a) WITH (FROM = 1)"},
		{"index WHERE invalid", "CREATE INDEX ix ON dbo.t(a) WITH (WHERE = 1)"},
		{"index DROP invalid", "CREATE INDEX ix ON dbo.t(a) WITH (DROP = ON)"},
		// SET options — valid
		{"SET ANSI_NULLS valid", "SET ANSI_NULLS ON"},
		{"SET QUOTED_IDENTIFIER valid", "SET QUOTED_IDENTIFIER OFF"},
		{"SET ARITHABORT valid", "SET ARITHABORT ON"},
		// SET options — invalid
		{"SET SELECT invalid", "SET SELECT ON"},
		{"SET FROM invalid", "SET FROM OFF"},
		{"SET WHERE invalid", "SET WHERE ON"},
		{"SET CREATE invalid", "SET CREATE ON"},
		{"SET INSERT invalid", "SET INSERT OFF"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), tc.sql)
			}
		})
	}
}

// TestKeywordOracleContextDisambiguation verifies that context keywords
// work correctly as identifiers in SQL Server.
func TestKeywordOracleContextDisambiguation(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		name string
		sql  string
	}{
		// Context keywords as table aliases
		{"window as table alias", "SELECT * FROM dbo.t window"},
		{"window as explicit alias", "SELECT * FROM dbo.t AS window"},
		{"encryption as table alias", "SELECT * FROM dbo.t encryption"},
		// Context keywords as column names
		{"partition as column name", "SELECT partition FROM dbo.t"},
		{"partition as qualified column", "SELECT t.partition FROM dbo.t"},
		{"encryption as qualified column", "SELECT t.encryption FROM dbo.t"},
		{"window as column name", "SELECT window FROM dbo.t"},
		// Context keywords as bare aliases
		{"encryption as bare alias", "SELECT 1 encryption FROM dbo.t"},
		{"partition as bare alias", "SELECT 1 partition FROM dbo.t"},
		{"window as bare alias", "SELECT 1 window FROM dbo.t"},
		// Context keywords as explicit aliases
		{"encryption as AS alias", "SELECT 1 AS encryption FROM dbo.t"},
		{"partition as AS alias", "SELECT 1 AS partition FROM dbo.t"},
		// Context keywords as table/column names in DDL
		{"CREATE TABLE encryption", "CREATE TABLE dbo.encryption_test (a INT)"},
		{"column named window", "CREATE TABLE dbo.t_window (window INT)"},
		{"column named partition", "CREATE TABLE dbo.t_partition (partition INT)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), tc.sql)
			}
		})
	}
}

// TestKeywordOracleCoreAsIdentifier verifies that core keywords are correctly
// rejected as unquoted identifiers and accepted when bracket-quoted.
func TestKeywordOracleCoreAsIdentifier(t *testing.T) {
	oracle := startParserOracle(t)

	coreKeywords := []string{
		"select", "from", "where", "insert", "update", "delete",
		"create", "drop", "alter", "order", "group", "having",
		"join", "on", "set", "into", "values", "table",
	}

	for _, kw := range coreKeywords {
		// Unquoted column name — should fail
		t.Run(fmt.Sprintf("column_%s_unquoted", kw), func(t *testing.T) {
			sql := fmt.Sprintf("CREATE TABLE dbo.t_core_%s (%s INT)", kw, kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})

		// Bracket-quoted column name — should succeed
		t.Run(fmt.Sprintf("column_%s_bracketed", kw), func(t *testing.T) {
			sql := fmt.Sprintf("CREATE TABLE dbo.t_core_%s_q ([%s] INT)", kw, kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})

		// Bare alias — should fail
		t.Run(fmt.Sprintf("alias_%s_bare", kw), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1 %s FROM dbo.t", kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})
	}
}

func boolToAcceptReject(b bool) string {
	if b {
		return "ACCEPTS"
	}
	return "REJECTS"
}
