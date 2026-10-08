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
// each argument type as the server's format_type does: a type known to be
// pg_catalog's under its SQL name, for every way the grammar spells one,
// and an unqualified name as written, which is what format_type writes
// when the search path finds a type of that name before pg_catalog.
func TestTypeNameMatchesFormatType(t *testing.T) {
	ctx := context.Background()
	db := startTypeNameOracle(t)
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	defer conn.Close()

	// Keywords with and without modifiers, which the grammar qualifies
	// with pg_catalog, names the SQL qualifies itself, and arrays. The
	// unqualified names here name no other type in this database.
	assertFormatType(t, conn, []string{
		"int", "integer", "pg_catalog.int4", "smallint", "bigint", "boolean", "pg_catalog.bool",
		"real", "float", "float(10)", "float(30)", "double precision", "pg_catalog.float8",
		"numeric", "numeric(10,2)", "decimal", "dec(5)",
		"varchar", "varchar(10)", "character varying(3)", "char varying", "pg_catalog.varchar",
		"character", "char", "char(3)", "nchar(2)", "national character varying(4)", "pg_catalog.bpchar",
		"bit", "bit(3)", "bit varying", "pg_catalog.varbit",
		"time", "time(3)", "time with time zone", "time(2) without time zone", "pg_catalog.timetz",
		"timestamp", "timestamp(6) with time zone", "timestamp without time zone", "pg_catalog.timestamptz",
		"interval", "interval year to month", "interval(3)",
		"json", "jsonb", "text", "pg_catalog.text", `"char"`, "name", "uuid", "bytea", "xml",
		"int[]", "integer[][]", "int ARRAY", "int ARRAY[3]", "text[]", "double precision[]", "pg_catalog.int8[]",
	})

	// A schema ahead of pg_catalog with a type of every catalog name a
	// user can write unqualified: the unqualified name now finds the
	// user's type, while a keyword or pg_catalog-qualified spelling still
	// finds the built-in one.
	shadowed := []string{
		"int2", "int4", "int8", "bool", "float4", "float8", "bpchar", "varbit", "timetz", "timestamptz",
		`"varchar"`, `"numeric"`, `"bit"`, `"time"`, `"timestamp"`, `"interval"`, `"json"`,
	}
	if _, err := conn.ExecContext(ctx, "CREATE SCHEMA s"); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	for _, name := range shadowed {
		if _, err := conn.ExecContext(ctx, "CREATE TYPE s."+name+" AS (a int)"); err != nil {
			t.Fatalf("create type s.%s: %v", name, err)
		}
	}
	if _, err := conn.ExecContext(ctx, "SET search_path = s, pg_catalog"); err != nil {
		t.Fatalf("set search_path: %v", err)
	}
	assertFormatType(t, conn, append(shadowed,
		"int4[]", "timestamptz[]",
		"integer", "character varying(3)", "timestamp with time zone", "bit varying", "pg_catalog.int4",
	))
}

// assertFormatType reviews DROP FUNCTION f(spelling) for each spelling and
// checks that the message names the type as format_type does on conn.
func assertFormatType(t *testing.T, conn *sql.Conn, spellings []string) {
	t.Helper()
	ctx := context.Background()
	for _, spelling := range spellings {
		var want string
		if err := conn.QueryRowContext(ctx, "SELECT format_type($1::text::regtype, NULL)", spelling).Scan(&want); err != nil {
			t.Fatalf("format_type(%q): %v", spelling, err)
		}
		sql := "DROP FUNCTION f(" + spelling + ")"
		res, err := Review(ctx, sql, review.Options{Rules: []review.Rule{review.DisallowDropObject}}, []review.Target{{}})
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
