package testcontainer

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPGKindOf(t *testing.T) {
	cases := []struct {
		sql  string
		want string
	}{
		{"SELECT 1", "query"},
		{"  -- leading comment\nSELECT 1", "query"},
		{"/* outer /* nested */ outer */ SELECT 1", "query"},
		{"/* unclosed comment", "query"},
		{"WITH cte AS (SELECT 1) SELECT * FROM cte", "query"},
		{"BEGIN", "query"},
		{"", "query"},
		{";", "query"},
		{"INSERT INTO t VALUES (1)", "dml"},
		{"-- c\nupdate t set a = 1", "dml"},
		{"DELETE FROM t", "dml"},
		{"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", "dml"},
		{"CREATE TABLE t (id bigint PRIMARY KEY)", "ddl"},
		{"/* ddl */ ALTER TABLE t ADD COLUMN c text", "ddl"},
		{"DROP TABLE t", "ddl"},
		{"GRANT SELECT ON TABLE t TO r", "ddl"},
		{"REVOKE SELECT ON TABLE t FROM r", "ddl"},
		{"TRUNCATE TABLE t", "ddl"},
		{"ANALYZE", "ddl"},
		{"RENAME TABLE a TO b", "ddl"},
	}
	for _, tc := range cases {
		if got := pgKindOf(tc.sql); got != tc.want {
			t.Errorf("pgKindOf(%q) = %q, want %q", tc.sql, got, tc.want)
		}
	}
}

func TestIsSingleUseDMLError(t *testing.T) {
	if isSingleUseDMLError(nil) {
		t.Fatal("nil should not match")
	}
	dmlErr := status.Error(codes.InvalidArgument, "DML statements may not be performed in single-use transactions, to avoid replay.")
	if !isSingleUseDMLError(dmlErr) {
		t.Fatal("expected single-use DML error to match")
	}
	otherErr := status.Error(codes.InvalidArgument, `[ERROR] syntax error at or near "zzz"`)
	if isSingleUseDMLError(otherErr) {
		t.Fatal("syntax error should not match single-use DML error")
	}
}
