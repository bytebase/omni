package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateSchemaQuotaParse(t *testing.T) {
	tests := []struct {
		sql      string
		name     string
		quota    string
		authrole string
	}{
		{
			sql:      "CREATE SCHEMA us_sales AUTHORIZATION dwuser QUOTA 50 GB;",
			name:     "us_sales",
			quota:    "50 GB",
			authrole: "dwuser",
		},
		{
			sql:      "CREATE SCHEMA IF NOT EXISTS reporting AUTHORIZATION report_user QUOTA 100 GB;",
			name:     "reporting",
			quota:    "100 GB",
			authrole: "report_user",
		},
		{
			sql:   "CREATE SCHEMA temp_schema QUOTA UNLIMITED;",
			name:  "temp_schema",
			quota: "UNLIMITED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstCreateSchemaStmt(t, tt.sql)
			if stmt.Schemaname != tt.name {
				t.Fatalf("expected schema name %q, got %q", tt.name, stmt.Schemaname)
			}
			if tt.authrole != "" {
				if stmt.Authrole == nil || stmt.Authrole.Rolename != tt.authrole {
					t.Fatalf("expected auth role %q, got %#v", tt.authrole, stmt.Authrole)
				}
			}
			assertDefElemString(t, stmt.Options, "quota", tt.quota)
		})
	}
}

func TestRedshiftAlterSchemaQuotaParse(t *testing.T) {
	tests := []struct {
		sql   string
		name  string
		quota string
	}{
		{sql: "ALTER SCHEMA us_sales QUOTA 300 GB;", name: "us_sales", quota: "300 GB"},
		{sql: "ALTER SCHEMA test_schema QUOTA UNLIMITED;", name: "test_schema", quota: "UNLIMITED"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftObjectStmt(t, tt.sql)
			if stmt.Command != "alter" {
				t.Fatalf("expected alter command, got %q", stmt.Command)
			}
			if stmt.ObjectType != "schema" {
				t.Fatalf("expected schema object type, got %q", stmt.ObjectType)
			}
			if stmt.Name == nil || len(stmt.Name.Items) != 1 {
				t.Fatalf("expected one schema name, got %#v", stmt.Name)
			}
			name, ok := stmt.Name.Items[0].(*nodes.String)
			if !ok || name.Str != tt.name {
				t.Fatalf("expected schema name %q, got %#v", tt.name, stmt.Name.Items[0])
			}
			assertDefElemString(t, stmt.Options, "quota", tt.quota)
		})
	}
}

func firstCreateSchemaStmt(t *testing.T, sql string) *nodes.CreateSchemaStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected CreateSchemaStmt, got %T", raw.Stmt)
	}
	return stmt
}
