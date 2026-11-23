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

// TestGrantPermissionsApplicable holds GRANT and REVOKE to what cassandra:5.0.9
// answered for each permission on each resource: filterPermissions
// (Parser.g:149) refuses a statement none of whose permissions the resource
// takes. LIST does not filter.
func TestGrantPermissionsApplicable(t *testing.T) {
	all := []string{"CREATE", "ALTER", "DROP", "SELECT", "MODIFY", "AUTHORIZE", "DESCRIBE", "EXECUTE", "UNMASK", "SELECT_MASKED"}
	accepted := map[string]string{
		"ALL KEYSPACES":               "CREATE ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"KEYSPACE k":                  "CREATE ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"ALL TABLES IN KEYSPACE k":    "CREATE ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"a":                           "ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"TABLE k.a":                   "ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"ROLE":                        "ALTER DROP SELECT MODIFY AUTHORIZE UNMASK SELECT_MASKED",
		"ALL ROLES":                   "CREATE ALTER DROP AUTHORIZE DESCRIBE",
		"ROLE r":                      "ALTER DROP AUTHORIZE DESCRIBE",
		"ALL FUNCTIONS":               "CREATE ALTER DROP AUTHORIZE EXECUTE",
		"ALL FUNCTIONS IN KEYSPACE k": "CREATE ALTER DROP AUTHORIZE EXECUTE",
		"ALL MBEANS":                  "SELECT MODIFY AUTHORIZE DESCRIBE EXECUTE",
		"MBEAN 'org.apache:type=a'":   "SELECT MODIFY AUTHORIZE DESCRIBE EXECUTE",
		"MBEANS 'org.apache:type=*'":  "SELECT MODIFY AUTHORIZE DESCRIBE EXECUTE",
	}
	for res, perms := range accepted {
		ok := map[string]bool{}
		for _, perm := range strings.Fields(perms) {
			ok[perm] = true
		}
		for _, perm := range all {
			for _, sql := range []string{"GRANT " + perm + " ON " + res + " TO r", "REVOKE " + perm + " ON " + res + " FROM r"} {
				if _, err := Parse(sql); (err == nil) != ok[perm] {
					t.Errorf("Parse(%q) error = %v, want accepted %v", sql, err, ok[perm])
				}
			}
			if sql := "LIST " + perm + " ON " + res; !parses(sql) {
				t.Errorf("Parse(%q) failed: LIST does not filter", sql)
			}
		}
		for _, sql := range []string{"GRANT ALL ON " + res + " TO r", "GRANT ALL PERMISSIONS ON " + res + " TO r"} {
			if !parses(sql) {
				t.Errorf("Parse(%q) failed", sql)
			}
		}
	}
	// One permission the resource takes is enough.
	for sql, want := range map[string]bool{
		"GRANT CREATE, SELECT ON a TO r":                   true,
		"GRANT CREATE PERMISSION, DESCRIBE ON k.a TO r":    false,
		"GRANT CREATE, SELECT ON ROLE r TO r":              false,
		"GRANT CREATE PERMISSION, DESCRIBE ON ROLE r TO r": true,
	} {
		if parses(sql) != want {
			t.Errorf("Parse(%q) accepted = %v, want %v", sql, !want, want)
		}
	}
}

func parses(sql string) bool {
	_, err := Parse(sql)
	return err == nil
}

// TestIntValueAndSignedConstants holds intValue (Parser.g:1601), DELETE's
// USING (Parser.g:612) and a property's constant to what cassandra:5.0.9
// answered: a '-' belongs to a number only written against it (Lexer.g:337).
func TestIntValueAndSignedConstants(t *testing.T) {
	for sql, want := range map[string]bool{
		"INSERT INTO a (a) VALUES (1) USING TTL ?":                           true,
		"INSERT INTO a (a) VALUES (1) USING TTL : a":                         true,
		"INSERT INTO a (a) VALUES (1) USING TTL -1":                          true,
		"INSERT INTO a (a) VALUES (1) USING TTL 1 AND TIMESTAMP 2 AND TTL 3": true,
		"INSERT INTO a (a) VALUES (1) USING TTL 'a'":                         false,
		"INSERT INTO a (a) VALUES (1) USING TTL null":                        false,
		"INSERT INTO a (a) VALUES (1) USING TTL 1.5":                         false,
		"INSERT INTO a (a) VALUES (1) USING TTL 1d":                          false,
		"INSERT INTO a (a) VALUES (1) USING TIMESTAMP - 1":                   false,
		"UPDATE a USING TIMESTAMP 1 AND TTL ? SET a = 1 WHERE a = 1":         true,
		"DELETE FROM a USING TIMESTAMP :a WHERE a = 1":                       true,
		"DELETE FROM a USING TIMESTAMP -1 WHERE a = 1":                       true,
		"DELETE FROM a USING TTL 1 WHERE a = 1":                              false,
		"DELETE FROM a USING TIMESTAMP 1 AND TIMESTAMP 2 WHERE a = 1":        false,
		"BEGIN BATCH USING TIMESTAMP ? AND TTL 1 APPLY BATCH":                true,
		"BEGIN BATCH USING TIMESTAMP 'a' APPLY BATCH":                        false,
		"SELECT a FROM a LIMIT -1":                                           true,
		"SELECT a FROM a PER PARTITION LIMIT -1":                             true,
		"SELECT a FROM a LIMIT - 1":                                          false,
		"CREATE KEYSPACE a WITH a = 1d":                                      true,
		"CREATE KEYSPACE a WITH a = -P1D":                                    true,
		"CREATE KEYSPACE a WITH replication = {'a': 1d}":                     true,
		"CREATE KEYSPACE a WITH a = -1":                                      true,
		"CREATE KEYSPACE a WITH a = - 1":                                     false,
		"CREATE KEYSPACE a WITH a = - 1d":                                    false,
		"GRANT ALL ON ROLE TO a":                                             true,
		"REVOKE ALL ON ROLE FROM a":                                          true,
		"GRANT ALL ON ROLE KEY TO a":                                         true,
		"GRANT ALL ON ROLE FROM TO a":                                        false,
	} {
		if parses(sql) != want {
			t.Errorf("Parse(%q) accepted = %v, want %v", sql, !want, want)
		}
	}
}
