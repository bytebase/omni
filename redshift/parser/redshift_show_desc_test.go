package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftShowParse(t *testing.T) {
	tests := []string{
		"SHOW DATABASES;",
		"SHOW DATABASES LIKE 'dev%' LIMIT 10;",
		"SHOW SCHEMAS FROM DATABASE dev LIMIT 10;",
		"SHOW TABLES FROM SCHEMA dev.public LIKE 't%' LIMIT 10;",
		"SHOW COLUMNS FROM TABLE public.t;",
		"SHOW TABLE users;",
		"SHOW EXTERNAL TABLE spectrum.users;",
		"SHOW GRANTS ON TABLE public.t;",
		"SHOW DATASHARES;",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			stmt := firstRedshiftShowStmt(t, sql)
			if stmt.Command == "" {
				t.Fatalf("expected non-empty command")
			}
		})
	}
}

func TestRedshiftShowVariableParse(t *testing.T) {
	tests := []struct {
		sql  string
		name string
	}{
		{sql: "SHOW current_user;", name: "current_user"},
		{sql: "SHOW session_user;", name: "session_user"},
		{sql: "SHOW current_database;", name: "current_database"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstVariableShowStmt(t, tt.sql)
			if stmt.Name != tt.name {
				t.Fatalf("Name = %q, want %q", stmt.Name, tt.name)
			}
		})
	}
}

func TestRedshiftDescParse(t *testing.T) {
	tests := []struct {
		sql        string
		objectType string
	}{
		{sql: "DESC DATASHARE share_name;", objectType: "datashare"},
		{sql: "DESC IDENTITY PROVIDER provider_name;", objectType: "identity provider"},
		{sql: "DESCRIBE DATASHARE share_name;", objectType: "datashare"},
		{sql: "DESCRIBE IDENTITY PROVIDER provider_name;", objectType: "identity provider"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftDescStmt(t, tt.sql)
			if stmt.ObjectType != tt.objectType {
				t.Fatalf("expected object type %q, got %q", tt.objectType, stmt.ObjectType)
			}
			if stmt.Name == nil {
				t.Fatalf("expected object name")
			}
		})
	}
}

func firstRedshiftShowStmt(t *testing.T, sql string) *nodes.RedshiftShowStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.RedshiftShowStmt)
	if !ok {
		t.Fatalf("expected RedshiftShowStmt, got %T", raw.Stmt)
	}
	return stmt
}

func firstVariableShowStmt(t *testing.T, sql string) *nodes.VariableShowStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected VariableShowStmt, got %T", raw.Stmt)
	}
	return stmt
}

func firstRedshiftDescStmt(t *testing.T, sql string) *nodes.RedshiftDescStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.RedshiftDescStmt)
	if !ok {
		t.Fatalf("expected RedshiftDescStmt, got %T", raw.Stmt)
	}
	return stmt
}
