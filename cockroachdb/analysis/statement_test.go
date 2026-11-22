package analysis

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassifyExecution(t *testing.T) {
	for _, tc := range []struct {
		name string
		stmt string
		want ExecutionFlags
	}{
		{
			name: "drop database is non-transactional",
			stmt: "DROP DATABASE mydb",
			want: ExecutionFlags{NonTransactional: true},
		},
		{
			name: "create index concurrently is non-transactional",
			stmt: "CREATE INDEX CONCURRENTLY idx ON t (a)",
			want: ExecutionFlags{NonTransactional: true},
		},
		{
			name: "create index without concurrently is transactional",
			stmt: "CREATE INDEX idx ON t (a)",
			want: ExecutionFlags{},
		},
		{
			name: "drop index concurrently is non-transactional",
			stmt: "DROP INDEX CONCURRENTLY t@idx",
			want: ExecutionFlags{NonTransactional: true},
		},
		{
			name: "drop index without concurrently is transactional",
			stmt: "DROP INDEX t@idx",
			want: ExecutionFlags{},
		},
		{
			name: "set role is session role",
			stmt: "SET ROLE app_user",
			want: ExecutionFlags{SessionRole: true},
		},
		{
			name: "set local role is not session role",
			stmt: "SET LOCAL ROLE app_user",
			want: ExecutionFlags{},
		},
		{
			name: "set other var is not session role",
			stmt: "SET search_path = app",
			want: ExecutionFlags{},
		},
		{
			name: "grant needs superuser",
			stmt: "GRANT SELECT ON TABLE t TO app_user",
			want: ExecutionFlags{Superuser: true},
		},
		{
			name: "grant role needs superuser",
			stmt: "GRANT admin TO app_user",
			want: ExecutionFlags{Superuser: true},
		},
		{
			name: "create extension needs superuser",
			stmt: "CREATE EXTENSION IF NOT EXISTS pgcrypto",
			want: ExecutionFlags{Superuser: true},
		},
		{
			name: "keyword in string literal does not trigger flag",
			stmt: "INSERT INTO t VALUES ('DROP DATABASE')",
			want: ExecutionFlags{},
		},
		{
			name: "multiple statements return zero flags",
			stmt: "DROP DATABASE a; DROP DATABASE b",
			want: ExecutionFlags{},
		},
		{
			name: "invalid sql returns zero flags",
			stmt: "NOT VALID SQL",
			want: ExecutionFlags{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifyExecution(tc.stmt)
			if got != tc.want {
				t.Errorf("ClassifyExecution(%q) = %+v, want %+v", tc.stmt, got, tc.want)
			}
		})
	}
}

func TestCreateDatabaseName(t *testing.T) {
	for _, tc := range []struct {
		stmt    string
		want    string
		wantErr bool
	}{
		{`CREATE DATABASE "hello" ENCODING "UTF8";`, "hello", false},
		{`CREATE DATABASE "hello";`, "hello", false},
		{`CREATE DATABASE hello;`, "hello", false},
		{`create database "team data";`, "team data", false},
		{`CREATE DATABASE;`, "", true},
		{`SELECT 1;`, "", true},
	} {
		t.Run(tc.stmt, func(t *testing.T) {
			got, err := CreateDatabaseName(tc.stmt)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("CreateDatabaseName(%q) = %q, want error", tc.stmt, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateDatabaseName(%q): %v", tc.stmt, err)
			}
			if got != tc.want {
				t.Errorf("CreateDatabaseName(%q) = %q, want %q", tc.stmt, got, tc.want)
			}
		})
	}
}

func TestIsDoBlock(t *testing.T) {
	for _, tc := range []struct {
		stmt string
		want bool
	}{
		{`DO $$ BEGIN DELETE FROM t; END $$`, true},
		{`SELECT 1`, false},
		{`DO $$ BEGIN DELETE FROM t; END $$; SELECT 1`, false},
		{`SELECT FROM`, false},
	} {
		if got := IsDoBlock(tc.stmt); got != tc.want {
			t.Errorf("IsDoBlock(%q) = %v, want %v", tc.stmt, got, tc.want)
		}
	}
}

func TestReturnsRows(t *testing.T) {
	for _, tc := range []struct {
		stmt    string
		want    bool
		wantErr bool
	}{
		{"SELECT 1", true, false},
		{"SELECT 1:::INT8", true, false},
		{"SHOW CLUSTER SETTING version", true, false},
		{"SHOW CREATE TABLE app.people", true, false},
		{"EXPLAIN ANALYZE UPDATE t SET v = v + 1", true, false},
		{"INSERT INTO t VALUES (1) RETURNING id", true, false},
		{"WITH d AS (DELETE FROM t RETURNING id) SELECT count(*) FROM d", true, false},
		{"UPDATE t SET v = 1", false, false},
		{"CREATE TABLE t (id INT8 PRIMARY KEY)", false, false},
		{"SELECT FROM", false, true},
	} {
		t.Run(tc.stmt, func(t *testing.T) {
			got, err := ReturnsRows(tc.stmt)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ReturnsRows(%q) succeeded, want error", tc.stmt)
				}
				return
			}
			if err != nil {
				t.Fatalf("ReturnsRows(%q): %v", tc.stmt, err)
			}
			if got != tc.want {
				t.Errorf("ReturnsRows(%q) = %v, want %v", tc.stmt, got, tc.want)
			}
		})
	}
}

func TestIsSelectStatement(t *testing.T) {
	for _, tc := range []struct {
		stmt string
		want bool
	}{
		{"SELECT 1:::INT8", true},
		{"(SELECT 1)", true},
		{"SHOW TABLES", false},
		{"INSERT INTO t VALUES (1) RETURNING id", false},
		{"SELECT 1; SELECT 2", false},
		{"SELECT FROM", false},
	} {
		if got := IsSelectStatement(tc.stmt); got != tc.want {
			t.Errorf("IsSelectStatement(%q) = %v, want %v", tc.stmt, got, tc.want)
		}
	}
}

func TestParseIndexDefinition(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stmt    string
		want    *IndexDefinition
		wantErr bool
	}{
		{
			name: "unique multi-column with direction",
			stmt: "CREATE UNIQUE INDEX idx_user_email ON public.users (email ASC, created_at DESC)",
			want: &IndexDefinition{
				Expressions: []string{"email", "created_at"},
				Descending:  []bool{false, true},
				Unique:      true,
			},
		},
		{
			name: "expression index",
			stmt: "CREATE INDEX idx_lower_name ON public.users (lower(name) DESC)",
			want: &IndexDefinition{
				Expressions: []string{"lower(name)"},
				Descending:  []bool{true},
				Unique:      false,
			},
		},
		{
			name:    "not a create index",
			stmt:    "SELECT 1",
			wantErr: true,
		},
		{
			name:    "multiple statements",
			stmt:    "CREATE INDEX a ON t (x); CREATE INDEX b ON t (y)",
			wantErr: true,
		},
		{
			name:    "syntax error",
			stmt:    "CREATE INDEX",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseIndexDefinition(tc.stmt)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseIndexDefinition(%q) = %+v, want error", tc.stmt, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseIndexDefinition(%q): %v", tc.stmt, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ParseIndexDefinition(%q) mismatch (-want +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
