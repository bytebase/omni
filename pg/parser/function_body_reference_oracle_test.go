package parser_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytebase/omni/pg/parser"
)

func TestFunctionBodyReferenceOracle(t *testing.T) {
	db, cleanup := startFunctionBodyOracle(t)
	defer cleanup()

	fixtures := []string{
		`SET check_function_bodies = on`,
		`CREATE TABLE t(id int PRIMARY KEY, v int)`,
		`CREATE TABLE archive(id int, v int)`,
		`CREATE SCHEMA api`,
		`CREATE MATERIALIZED VIEW mv AS SELECT 1 AS id`,
	}
	for _, sql := range fixtures {
		if _, err := db.ExecContext(context.Background(), sql); err != nil {
			t.Fatalf("fixture %q: %v", sql, err)
		}
	}

	tests := []struct {
		name   string
		sql    string
		accept bool
	}{
		{
			name:   "language sql select body",
			sql:    `CREATE FUNCTION ref_sql_select(x int) RETURNS int LANGUAGE sql AS $$ SELECT x + 1 $$`,
			accept: true,
		},
		{
			name:   "language sql dml returning body",
			sql:    `CREATE FUNCTION ref_sql_insert() RETURNS int LANGUAGE sql AS $$ INSERT INTO t(id, v) VALUES (1, 1) RETURNING id $$`,
			accept: true,
		},
		{
			name:   "language sql rejects plpgsql if",
			sql:    `CREATE FUNCTION ref_sql_bad_if() RETURNS int LANGUAGE sql AS $$ IF true THEN SELECT 1; END IF $$`,
			accept: false,
		},
		{
			name:   "language plpgsql returning into",
			sql:    `CREATE FUNCTION ref_plpgsql_returning() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x int; BEGIN INSERT INTO t(id, v) VALUES (2, 2) RETURNING id INTO x; RETURN x; END $$`,
			accept: true,
		},
		{
			name:   "language plpgsql utility sql",
			sql:    `CREATE FUNCTION ref_plpgsql_utility() RETURNS void LANGUAGE plpgsql AS $$ BEGIN REFRESH MATERIALIZED VIEW mv; GRANT USAGE ON SCHEMA api TO public; END $$`,
			accept: true,
		},
		{
			name:   "language plpgsql rejects missing into target",
			sql:    `CREATE FUNCTION ref_plpgsql_bad_into() RETURNS int LANGUAGE plpgsql AS $$ BEGIN SELECT 1 INTO; END $$`,
			accept: false,
		},
		{
			name:   "sql standard return body",
			sql:    `CREATE FUNCTION ref_sql_standard_return() RETURNS int LANGUAGE sql RETURN 1`,
			accept: true,
		},
		{
			name:   "sql standard rejects malformed atomic body",
			sql:    `CREATE FUNCTION ref_sql_standard_bad() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT FROM; END`,
			accept: false,
		},
		{
			name:   "sql standard rejects missing atomic keyword",
			sql:    `CREATE FUNCTION ref_sql_standard_missing_atomic() RETURNS int LANGUAGE sql BEGIN RETURN 1; END`,
			accept: false,
		},
		{
			name:   "sql standard rejects missing end keyword",
			sql:    `CREATE FUNCTION ref_sql_standard_missing_end() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1;`,
			accept: false,
		},
		{
			name:   "rejects AS without string body",
			sql:    `CREATE FUNCTION ref_bad_as() RETURNS int AS LANGUAGE sql`,
			accept: false,
		},
		{
			name:   "rejects LANGUAGE without name",
			sql:    `CREATE FUNCTION ref_bad_language() RETURNS int LANGUAGE`,
			accept: false,
		},
		{
			name:   "rejects malformed TRANSFORM clause",
			sql:    `CREATE FUNCTION ref_bad_transform() RETURNS int TRANSFORM int LANGUAGE sql AS 'SELECT 1'`,
			accept: false,
		},
		{
			name:   "rejects malformed SECURITY clause",
			sql:    `CREATE FUNCTION ref_bad_security() RETURNS int SECURITY LANGUAGE sql AS 'SELECT 1'`,
			accept: false,
		},
		{
			name:   "rejects malformed SET TIME clause",
			sql:    `CREATE FUNCTION ref_bad_set_time() RETURNS int SET TIME 'UTC' LANGUAGE sql AS 'SELECT 1'`,
			accept: false,
		},
		{
			name:   "rejects malformed SET SESSION clause",
			sql:    `CREATE FUNCTION ref_bad_set_session() RETURNS int SET SESSION ROLE LANGUAGE sql AS 'SELECT 1'`,
			accept: false,
		},
		{
			name:   "rejects malformed SET FROM clause",
			sql:    `CREATE FUNCTION ref_bad_set_from() RETURNS int SET search_path FROM 'public' LANGUAGE sql AS 'SELECT 1'`,
			accept: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, omniErr := pgParseWithBodyCheck(tt.sql)
			omniAccept := omniErr == nil
			if omniAccept != tt.accept {
				t.Fatalf("omni accept = %v, want %v, err=%v", omniAccept, tt.accept, omniErr)
			}

			_, pgErr := db.ExecContext(context.Background(), tt.sql)
			pgAccept := pgErr == nil
			if pgAccept != tt.accept {
				t.Fatalf("postgres accept = %v, want %v, err=%v", pgAccept, tt.accept, pgErr)
			}
		})
	}
}

func startFunctionBodyOracle(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	db, err := parser.OracleDatabase(context.Background(), "omni_body")
	if err != nil {
		t.Fatalf("function body oracle unavailable in CI: %v", err)
	}
	return db, func() { db.Close() }
}
