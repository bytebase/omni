package review

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/review"
)

// TestTypeNameMatchesFormatType checks that a DROP FUNCTION message names
// each argument type as the server's format_type does, for every way the
// grammar spells a built-in type: keywords with and without modifiers,
// catalog names, quoted names, qualified names, and arrays.
func TestTypeNameMatchesFormatType(t *testing.T) {
	db := startTypeNameOracle(t)

	spellings := []string{
		"int", "integer", "int4", "pg_catalog.int4", `"int4"`,
		"smallint", "int2", "bigint", "int8",
		"boolean", "bool",
		"real", "float4", "float", "float(10)", "float(30)", "double precision", "float8",
		"numeric", "numeric(10,2)", "decimal", "dec(5)", `"numeric"`,
		"varchar", "varchar(10)", "character varying(3)", "char varying", `"varchar"`, "pg_catalog.varchar",
		"character", "char", "char(3)", "bpchar", "nchar(2)", "national character varying(4)",
		"bit", "bit(3)", "bit varying", "varbit",
		"time", "time(3)", "time with time zone", "timetz", "time(2) without time zone",
		"timestamp", "timestamp(6) with time zone", "timestamptz", "timestamp without time zone",
		"interval", "interval year to month", "interval(3)",
		"json", "jsonb", "text", "pg_catalog.text", `"char"`, "name", "uuid", "bytea", "xml",
		"int[]", "integer[][]", "int ARRAY", "int ARRAY[3]", "text[]", "double precision[]", "timestamptz[]",
	}
	for _, spelling := range spellings {
		var want string
		if err := db.QueryRowContext(context.Background(), "SELECT format_type($1::text::regtype, NULL)", spelling).Scan(&want); err != nil {
			t.Fatalf("format_type(%q): %v", spelling, err)
		}
		sql := "DROP FUNCTION f(" + spelling + ")"
		res, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.DisallowDropObject}}, []review.Target{{}})
		if err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		if len(res.Findings) != 1 {
			t.Fatalf("%s: got %d findings, want 1", sql, len(res.Findings))
		}
		if got, wantMsg := res.Findings[0].Message, "drops function f("+want+")"; got != wantMsg {
			t.Errorf("%s: got %q, want %q", sql, got, wantMsg)
		}
	}
}

func startTypeNameOracle(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	var container *tcpg.PostgresContainer
	var setupErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				setupErr = fmt.Errorf("docker provider panic: %v", r)
			}
		}()
		var err error
		container, err = tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("omni_review"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			setupErr = fmt.Errorf("container start: %w", err)
		}
	}()
	if setupErr != nil {
		t.Fatalf("type name oracle unavailable: %v", setupErr)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("conn string: %v", err)
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return db
}
