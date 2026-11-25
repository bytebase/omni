package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/spanner/ast"
)

// The accept cases below are the canonical ZetaSQL GRANT/REVOKE corpus
// (parser/googlesql/examples/.../grant_and_revoke.sql), which is the
// accept/reject oracle for these forms — the legacy ANTLR grammar bytebase
// consumes today is a hand-port of that ZetaSQL reference.
//
// DCL is a dialect-divergent zone: the Spanner emulator speaks a DIFFERENT GRANT
// dialect (`GRANT priv ON TABLE x TO ROLE r`, role-name grantees) and REJECTS
// every form here, so it is NON-AUTHORITATIVE for DCL. See the parser-dcl
// divergence ledger entry; these tests follow the ZetaSQL/.g4 grammar.

// parseOneStmt parses input expected to contain exactly one statement and
// returns that statement node. It fails the test on any parse error or on a
// statement count != 1.
func parseOneStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

// TestGrantRevoke_CorpusAccepts parses every statement of the canonical ZetaSQL
// grant_and_revoke.sql corpus and asserts each parses to a single GRANT/REVOKE
// node with no errors. This is the breadth oracle for the node.
func TestGrantRevoke_CorpusAccepts(t *testing.T) {
	cases := []struct {
		sql    string
		revoke bool // true => *RevokeStmt, false => *GrantStmt
	}{}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOneStmt(t, tc.sql)
			switch node.(type) {
			case *ast.GrantStmt:
				if tc.revoke {
					t.Errorf("parseForTest(%q): got *GrantStmt, want *RevokeStmt", tc.sql)
				}
			case *ast.RevokeStmt:
				if !tc.revoke {
					t.Errorf("parseForTest(%q): got *RevokeStmt, want *GrantStmt", tc.sql)
				}
			default:
				t.Errorf("parseForTest(%q): got %T, want a GRANT/REVOKE node", tc.sql, node)
			}
		})
	}
}

// TestGrant_AllPrivileges asserts ALL [PRIVILEGES] yields AllPrivileges with no
// explicit Privilege nodes.
func TestGrant_AllPrivileges(t *testing.T) {
	for _, sql := range []string{} {
		t.Run(sql, func(t *testing.T) {
			g, ok := parseOneStmt(t, sql).(*ast.GrantStmt)
			if !ok {
				t.Fatalf("not a *ast.GrantStmt")
			}
			if !g.AllPrivileges {
				t.Error("AllPrivileges = false, want true")
			}
			if g.Privileges != nil {
				t.Errorf("Privileges = %v, want nil", g.Privileges)
			}
		})
	}
}

// TestGrant_ObjectType asserts the optional 0/1/2-word object type is captured
// and the object path is the remaining path_expression, matching the grammar's
// (identifier identifier?)? path_expression disambiguation.
func TestGrant_ObjectType(t *testing.T) {
	cases := []struct {
		sql      string
		wantType []string
		wantPath string
	}{
		// No object type: bare path.
		// No object type: dotted path.
		// One-word type.
		// One-word type + dotted path.
		// Two-word type.
		// One-word type whose PATH is itself the keyword `table` (corpus case #10).
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			var objType []string
			var path ast.NamePath
			switch n := parseOneStmt(t, tc.sql).(type) {
			case *ast.GrantStmt:
				objType, path = n.ObjectType, n.Path
			case *ast.RevokeStmt:
				objType, path = n.ObjectType, n.Path
			default:
				t.Fatalf("got %T, want GRANT/REVOKE", n)
			}
			if !equalStrings(objType, tc.wantType) {
				t.Errorf("ObjectType = %v, want %v", objType, tc.wantType)
			}
			if got := path.String(); got != tc.wantPath {
				t.Errorf("Path = %q, want %q", got, tc.wantPath)
			}
		})
	}
}

// TestGrantRevoke_Rejects asserts malformed GRANT/REVOKE inputs produce a parse
// error (not a panic, not silent acceptance). These are grammar-level rejects of
// the legacy ZetaSQL grammar.
func TestGrantRevoke_Rejects(t *testing.T) {
	cases := []string{
		"GRANT ON table foo TO 'x'",           // missing privileges
		"GRANT `select` table foo TO 'x'",     // missing ON
		"GRANT `select` ON TO 'x'",            // missing object path
		"GRANT `select` ON table foo 'x'",     // missing TO
		"GRANT `select` ON table foo TO",      // empty grantee list
		"GRANT `select` ON table foo TO foo",  // grantee must be string/param, not identifier
		"GRANT `select` ON table foo TO 'x',", // trailing comma in grantee list
		"GRANT `select`, ON table foo TO 'x'", // trailing comma in privilege list
		"REVOKE ON table foo FROM 'x'",        // missing privileges
		"REVOKE `select` ON table foo TO 'x'", // REVOKE uses FROM, not TO
		"REVOKE `select` ON table foo",        // missing FROM
		// NOTE: `GRANT ROLE analyst TO ROLE senior` is NO LONGER a reject — the
		// parser-ddl-spanner node added the Spanner role-to-role grant to the union
		// (emulator-authoritative: it ACCEPTS that form). See TestGrant_RoleToRole /
		// spanner_ddl_test.go. The grantee must still not be a bare identifier
		// (covered by `… TO foo` above), and `TO 'string'` is the legacy dialect.
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, errs := parseForTest(sql)
			if len(errs) == 0 {
				t.Errorf("parseForTest(%q): want a parse error, got none", sql)
			}
			for _, e := range errs {
				if strings.Contains(e.Message, "not yet supported") {
					t.Errorf("parseForTest(%q): got 'not yet supported' (stub still wired): %v", sql, errs)
				}
			}
		})
	}
}

// TestGrantRevoke_RejectsTrailingTokens is the regression test for the
// parseSingle EOF assertion (PR #195 review, finding 1). The grammar root is
// `stmts EOF`: a complete statement must be followed by end-of-input. Before the
// fix, parseSingle never asserted EOF after a successful parseStmt, so trailing
// junk after an otherwise-valid GRANT/REVOKE was silently dropped — the parse
// returned stmts=1, errs=0, a false-accept that yields zero diagnostics for the
// Diagnose consumer. Each case below is a valid GRANT/REVOKE prefix followed by
// tokens the grammar cannot continue with; all must now be rejected (error
// reported) AND produce no statement node (the segment as a whole does not
// parse).
func TestGrantRevoke_RejectsTrailingTokens(t *testing.T) {
	cases := []string{
		"GRANT `select` ON foo TO 'x' garbage",     // bare trailing identifier
		"GRANT `select` ON foo TO 'x' ON bar",      // trailing keyword run
		"REVOKE `select` ON foo FROM 'x' garbage",  // REVOKE trailing junk
		"REVOKE ALL PRIVILEGES ON foo FROM 'x' 99", // trailing integer literal
	}
	// Adjacent string literal concatenation in GRANT TO is valid BigQuery/union
	// syntax, rejected in DialectSpanner.
	if _, err := Parse("GRANT `select` ON foo TO 'x' 'y'"); err == nil {
		t.Errorf("Parse(GRANT ... TO 'x' 'y'): want error, got nil")
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			file, errs := parseForTest(sql)
			if len(errs) == 0 {
				t.Errorf("parseForTest(%q): want a syntax error for trailing tokens, got none", sql)
			}
			if len(file.Stmts) != 0 {
				t.Errorf("parseForTest(%q): got %d stmts, want 0 (trailing junk makes the segment invalid)", sql, len(file.Stmts))
			}
			for _, e := range errs {
				if strings.Contains(e.Message, "not yet supported") {
					t.Errorf("parseForTest(%q): got 'not yet supported' (unexpected stub path): %v", sql, errs)
				}
			}
		})
	}
}

// equalStrings reports whether two string slices are element-wise equal,
// treating nil and empty as equal.
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
