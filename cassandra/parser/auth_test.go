package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
)

// TestFunctionResourceRequiresKeyspace holds functionResource (Parser.g:1191)
// to Cassandra's parse-time check: FunctionResource.functionFromCql
// unconditionally rejects an unqualified function name in GRANT, REVOKE and
// LIST statements, while a table named FUNCTION still parses as a dataResource.
func TestFunctionResourceRequiresKeyspace(t *testing.T) {
	accepted := []struct {
		sql      string
		wantType string
		wantArgs int
	}{
		{"GRANT EXECUTE ON FUNCTION ks.a() TO r", "FUNCTION", 0},
		{"GRANT EXECUTE ON FUNCTION ks.a(int) TO r", "FUNCTION", 1},
		{"GRANT EXECUTE ON FUNCTION ks.a(ascii, text) TO r", "FUNCTION", 2},
		{"GRANT ALL ON FUNCTION ks.token(int) TO r", "FUNCTION", 1},
		{"GRANT ALL ON FUNCTION ks.count() TO r", "FUNCTION", 0},
		{`GRANT EXECUTE ON FUNCTION "ks"."fn"(int) TO r`, "FUNCTION", 1},
		{"REVOKE EXECUTE ON FUNCTION ks.a() FROM r", "FUNCTION", 0},
		{"REVOKE ALL ON FUNCTION ks.a(ascii, ascii) FROM r", "FUNCTION", 2},
		{"LIST ALL ON FUNCTION ks.a()", "FUNCTION", 0},
		{"LIST ALL ON FUNCTION ks.a(ascii)", "FUNCTION", 1},
		{"LIST ALL ON FUNCTION ks.a(ascii, ascii) OF r NORECURSIVE", "FUNCTION", 2},
		{"LIST ALL ON ALL FUNCTIONS", "ALL FUNCTIONS", 0},
		{"LIST ALL ON ALL FUNCTIONS IN KEYSPACE ks", "ALL FUNCTIONS IN KEYSPACE", 0},
		// FUNCTION is an unreserved keyword, so it remains a valid bare table
		// name when not followed by a function resource name.
		{"LIST ALL ON FUNCTION", "TABLE", 0},
		{"LIST ALL ON FUNCTION OF r", "TABLE", 0},
		{"LIST ALL ON FUNCTION NORECURSIVE", "TABLE", 0},
		{"LIST ALL ON FUNCTION.t", "TABLE", 0},
		{"GRANT SELECT ON FUNCTION TO r", "TABLE", 0},
		{"REVOKE SELECT ON FUNCTION FROM r", "TABLE", 0},
	}

	for _, tc := range accepted {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tc.sql, err)
			continue
		}
		if list == nil || len(list.Items) != 1 {
			t.Errorf("Parse(%q) returned %v items, want 1", tc.sql, list)
			continue
		}
		var res *ast.Resource
		switch s := list.Items[0].(*ast.RawStmt).Stmt.(type) {
		case *ast.GrantStmt:
			res = s.Resource
		case *ast.RevokeStmt:
			res = s.Resource
		case *ast.ListPermissionsStmt:
			res = s.Resource
		default:
			t.Errorf("Parse(%q) produced unexpected stmt %T", tc.sql, s)
			continue
		}
		if res == nil || res.Type != tc.wantType {
			t.Errorf("Parse(%q) resource = %+v, want Type=%q", tc.sql, res, tc.wantType)
		}
		if res != nil && len(res.ArgTypes) != tc.wantArgs {
			t.Errorf("Parse(%q) ArgTypes len = %d, want %d", tc.sql, len(res.ArgTypes), tc.wantArgs)
		}
	}

	rejected := []string{
		"LIST ALL ON FUNCTION a()",
		"LIST ALL ON FUNCTION a(zzz)",
		"LIST ALL ON FUNCTION a(ASCII)",
		"LIST ALL ON FUNCTION a(ASCII, ASCII)",
		`LIST ALL ON FUNCTION "a"()`,
		"LIST ALL ON FUNCTION token()",
		"LIST ALL ON FUNCTION count(int)",
		"GRANT EXECUTE ON FUNCTION a() TO r",
		"GRANT EXECUTE ON FUNCTION a(int, text) TO r",
		"REVOKE EXECUTE ON FUNCTION a() FROM r",
		"REVOKE EXECUTE ON FUNCTION a(int) FROM r",
		"LIST ALL ON FUNCTION ks.a",
		"GRANT EXECUTE ON FUNCTION ks.a TO r",
	}

	for _, sql := range rejected {
		_, err := Parse(sql)
		if err == nil {
			t.Errorf("Parse(%q) = nil error, want *ParseError", sql)
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error = %T (%v), want *ParseError", sql, err, err)
			continue
		}
		if strings.Contains(sql, "(") && !strings.Contains(sql, ".") {
			if !strings.Contains(pe.Message, "keyspace") {
				t.Errorf("Parse(%q) error = %q, want mention of keyspace", sql, pe.Message)
			}
		}
	}
}
