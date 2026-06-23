package deparse

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// TestDeparseSystemTime guards that the FOR SYSTEM_TIME temporal clause survives
// a round trip through DeparseSelect (dropping it changes query semantics).
func TestDeparseSystemTime(t *testing.T) {
	cases := []struct {
		sql  string
		want string
	}{
		{"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'", "for system_time as of '2020-01-01'"},
		{"SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' AND '2020-06-01'", "for system_time between '2020-01-01' and '2020-06-01'"},
		{"SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' TO '2020-06-01'", "for system_time from '2020-01-01' to '2020-06-01'"},
		{"SELECT * FROM t FOR SYSTEM_TIME ALL", "for system_time all"},
		{"SELECT * FROM t FOR SYSTEM_TIME AS OF TRANSACTION 12345", "for system_time as of transaction 12345"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := parser.Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse %q: %v", tc.sql, err)
			}
			stmt, ok := list.Items[0].(*ast.SelectStmt)
			if !ok {
				t.Fatalf("expected SelectStmt, got %T", list.Items[0])
			}
			got := DeparseSelect(stmt)
			if !strings.Contains(got, tc.want) {
				t.Errorf("DeparseSelect(%q) = %q, missing %q", tc.sql, got, tc.want)
			}
		})
	}
}
