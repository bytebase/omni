package driver

import (
	"strings"
	"testing"
)

func TestDMLStatement(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		wantErr   string
	}{
		{name: "insert", statement: "INSERT INTO t VALUES (1)"},
		{name: "replace", statement: "REPLACE INTO t VALUES (1)"},
		{name: "update", statement: "UPDATE t SET c = 1 WHERE id = 1"},
		{name: "delete", statement: "DELETE FROM t WHERE id = 1"},
		{name: "select", statement: "SELECT * FROM t", wantErr: "expected an INSERT, UPDATE or DELETE"},
		{name: "script", statement: "DELETE FROM t WHERE id = 1; SELECT 1", wantErr: "expected exactly one statement"},
		{name: "empty", statement: "", wantErr: "expected exactly one statement"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := dmlStatement(tc.statement)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("dmlStatement(%q) = (%T, %v), want error containing %q", tc.statement, stmt, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("dmlStatement(%q) = %v", tc.statement, err)
			}
			if stmt == nil {
				t.Fatalf("dmlStatement(%q) returned no statement", tc.statement)
			}
		})
	}
}
