package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterDatabaseOptionsParse(t *testing.T) {
	tests := []struct {
		sql       string
		option    string
		wantValue string
	}{
		{"ALTER DATABASE sampledb COLLATE CASE_SENSITIVE;", "collate", "case_sensitive"},
		{"ALTER DATABASE reports_db COLLATE CS;", "collate", "cs"},
		{"ALTER DATABASE analytics COLLATE CASE_INSENSITIVE;", "collate", "case_insensitive"},
		{"ALTER DATABASE staging_db COLLATE CI;", "collate", "ci"},
		{"ALTER DATABASE production_db CONNECTION LIMIT UNLIMITED;", "connection_limit", "unlimited"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := singleStmt(t, tt.sql)
			alter, ok := stmt.(*nodes.AlterDatabaseStmt)
			if !ok {
				t.Fatalf("expected AlterDatabaseStmt, got %T", stmt)
			}
			assertDatabaseOptionString(t, alter.Options, tt.option, tt.wantValue)
		})
	}
}

func TestRedshiftCreateDatabaseOptionsParse(t *testing.T) {
	tests := []struct {
		sql       string
		option    string
		wantValue string
	}{
		{"CREATE DATABASE sampledb ISOLATION LEVEL SNAPSHOT;", "isolation_level", "snapshot"},
		{"CREATE DATABASE transactionaldb ISOLATION LEVEL SERIALIZABLE;", "isolation_level", "serializable"},
		{"CREATE DATABASE cidb COLLATE CASE_INSENSITIVE;", "collate", "case_insensitive"},
		{"CREATE DATABASE cidb2 COLLATE CI;", "collate", "ci"},
		{"CREATE DATABASE csdb COLLATE CASE_SENSITIVE;", "collate", "case_sensitive"},
		{"CREATE DATABASE csdb2 COLLATE CS;", "collate", "cs"},
		{"CREATE DATABASE enterprisedb OWNER = admin CONNECTION LIMIT 200 COLLATE CI ISOLATION LEVEL SNAPSHOT;", "isolation_level", "snapshot"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := singleStmt(t, tt.sql)
			create, ok := stmt.(*nodes.CreatedbStmt)
			if !ok {
				t.Fatalf("expected CreatedbStmt, got %T", stmt)
			}
			assertDatabaseOptionString(t, create.Options, tt.option, tt.wantValue)
		})
	}
}

func assertDatabaseOptionString(t *testing.T, opts *nodes.List, option string, wantValue string) {
	t.Helper()
	elem := findDefElem(opts, option)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", option, opts)
	}
	value, ok := elem.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected option %q string arg, got %T", option, elem.Arg)
	}
	if value.Str != wantValue {
		t.Fatalf("expected option %q=%q, got %q", option, wantValue, value.Str)
	}
}
