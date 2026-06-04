package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
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
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
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
	}{
		{"GRANT `select`, `update` ON table foo TO 'user@google.com'", false},
		{"GRANT ALL PRIVILEGES ON view foo TO @user1, @@user2, 'mdbuser/bar1', 'mdbuser/bar2'", false},
		{"GRANT `select` ON materialized view foo TO 'user@google.com'", false},
		{"REVOKE `select` ON materialized view foo FROM 'user@google.com'", true},
		{"GRANT ALL PRIVILEGES ON datascape.foo TO 'bar'", false},
		{"GRANT `select`, insert(col1, col2, col3), `update`(col2) ON foo TO 'mdbgroup/bar'", false},
		{"GRANT execute ON script datascape.script_foo TO 'group@google.com'", false},
		{"REVOKE ALL PRIVILEGES ON foo FROM 'bar'", true},
		{"REVOKE delete ON table foo FROM 'mdbuser/bar'", true},
		{"REVOKE ALL PRIVILEGES ON table table FROM 'mdbuser/user', @user2, 'user3', @@user4", true},
		{"REVOKE delete, `update`(col2) ON view foo FROM 'mdbgroup/bar'", true},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOneStmt(t, tc.sql)
			switch node.(type) {
			case *ast.GrantStmt:
				if tc.revoke {
					t.Errorf("Parse(%q): got *GrantStmt, want *RevokeStmt", tc.sql)
				}
			case *ast.RevokeStmt:
				if !tc.revoke {
					t.Errorf("Parse(%q): got *RevokeStmt, want *GrantStmt", tc.sql)
				}
			default:
				t.Errorf("Parse(%q): got %T, want a GRANT/REVOKE node", tc.sql, node)
			}
		})
	}
}

// TestGrant_PrivilegeList asserts the explicit privilege list, including
// backtick-quoted names, bare non-reserved keyword names, and per-privilege
// column lists, is parsed structurally.
func TestGrant_PrivilegeList(t *testing.T) {
	node := parseOneStmt(t, "GRANT `select`, insert(col1, col2, col3), `update`(col2) ON foo TO 'mdbgroup/bar'")
	g, ok := node.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.GrantStmt", node)
	}
	if g.AllPrivileges {
		t.Error("AllPrivileges = true, want false (explicit list)")
	}
	if len(g.Privileges) != 3 {
		t.Fatalf("got %d privileges, want 3", len(g.Privileges))
	}
	// `select` — backtick-quoted, normalized to "select", no columns.
	if g.Privileges[0].Name != "select" {
		t.Errorf("priv[0].Name = %q, want %q", g.Privileges[0].Name, "select")
	}
	if g.Privileges[0].Columns != nil {
		t.Errorf("priv[0].Columns = %v, want nil", g.Privileges[0].Columns)
	}
	// insert(col1, col2, col3) — three columns.
	if g.Privileges[1].Name != "insert" {
		t.Errorf("priv[1].Name = %q, want %q", g.Privileges[1].Name, "insert")
	}
	if len(g.Privileges[1].Columns) != 3 {
		t.Fatalf("priv[1] got %d columns, want 3", len(g.Privileges[1].Columns))
	}
	wantCols := []string{"col1", "col2", "col3"}
	for i, want := range wantCols {
		if got := g.Privileges[1].Columns[i].String(); got != want {
			t.Errorf("priv[1].Columns[%d] = %q, want %q", i, got, want)
		}
	}
	// `update`(col2) — one column.
	if g.Privileges[2].Name != "update" {
		t.Errorf("priv[2].Name = %q, want %q", g.Privileges[2].Name, "update")
	}
	if len(g.Privileges[2].Columns) != 1 || g.Privileges[2].Columns[0].String() != "col2" {
		t.Errorf("priv[2].Columns = %v, want [col2]", g.Privileges[2].Columns)
	}
}

// TestGrant_AllPrivileges asserts ALL [PRIVILEGES] yields AllPrivileges with no
// explicit Privilege nodes.
func TestGrant_AllPrivileges(t *testing.T) {
	for _, sql := range []string{
		"GRANT ALL PRIVILEGES ON datascape.foo TO 'bar'",
		"GRANT ALL ON datascape.foo TO 'bar'", // PRIVILEGES is optional
	} {
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
		{"GRANT ALL PRIVILEGES ON foo TO 'bar'", nil, "foo"},
		// No object type: dotted path.
		{"GRANT ALL PRIVILEGES ON datascape.foo TO 'bar'", nil, "datascape.foo"},
		// One-word type.
		{"GRANT `select` ON table foo TO 'x'", []string{"table"}, "foo"},
		// One-word type + dotted path.
		{"GRANT execute ON script datascape.script_foo TO 'x'", []string{"script"}, "datascape.script_foo"},
		// Two-word type.
		{"GRANT `select` ON materialized view foo TO 'x'", []string{"materialized", "view"}, "foo"},
		// One-word type whose PATH is itself the keyword `table` (corpus case #10).
		{"REVOKE ALL PRIVILEGES ON table table FROM 'x'", []string{"table"}, "table"},
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

// TestGrant_GranteeKinds asserts each grantee shape of string_literal_or_parameter
// is parsed with the right kind and payload.
func TestGrant_GranteeKinds(t *testing.T) {
	node := parseOneStmt(t, "GRANT ALL PRIVILEGES ON view foo TO @user1, @@user2, 'mdbuser/bar1', ?")
	g := node.(*ast.GrantStmt)
	if len(g.Grantees) != 4 {
		t.Fatalf("got %d grantees, want 4", len(g.Grantees))
	}
	want := []struct {
		kind  ast.GranteeKind
		value string
	}{
		{ast.GranteeNamedParameter, "user1"},
		{ast.GranteeSystemVariable, "user2"},
		{ast.GranteeString, "mdbuser/bar1"},
		{ast.GranteePositionalParameter, ""},
	}
	for i, w := range want {
		if g.Grantees[i].Kind != w.kind {
			t.Errorf("grantee[%d].Kind = %v, want %v", i, g.Grantees[i].Kind, w.kind)
		}
		if g.Grantees[i].Value != w.value {
			t.Errorf("grantee[%d].Value = %q, want %q", i, g.Grantees[i].Value, w.value)
		}
	}
}

// TestRevoke_Shape asserts REVOKE parses to a *RevokeStmt with the FROM grantee
// list, mirroring GRANT.
func TestRevoke_Shape(t *testing.T) {
	node := parseOneStmt(t, "REVOKE delete, `update`(col2) ON view foo FROM 'mdbgroup/bar', @p")
	r, ok := node.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.RevokeStmt", node)
	}
	if len(r.Privileges) != 2 {
		t.Errorf("got %d privileges, want 2", len(r.Privileges))
	}
	if !equalStrings(r.ObjectType, []string{"view"}) {
		t.Errorf("ObjectType = %v, want [view]", r.ObjectType)
	}
	if len(r.Grantees) != 2 {
		t.Errorf("got %d grantees, want 2", len(r.Grantees))
	}
}

// TestGrantRevoke_SelectIsPrivilegeName asserts the reserved keyword SELECT is
// accepted as a (bare) privilege name — the grammar's privilege_name:
// identifier | SELECT special case (SELECT is reserved, so it would not match
// the generic identifier rule).
func TestGrantRevoke_SelectIsPrivilegeName(t *testing.T) {
	g := parseOneStmt(t, "GRANT select ON table foo TO 'x'").(*ast.GrantStmt)
	if len(g.Privileges) != 1 || g.Privileges[0].Name != "select" {
		t.Fatalf("Privileges = %v, want a single 'select'", g.Privileges)
	}
	// SELECT with a column list, too.
	g2 := parseOneStmt(t, "GRANT select(col1) ON table foo TO 'x'").(*ast.GrantStmt)
	if len(g2.Privileges) != 1 || len(g2.Privileges[0].Columns) != 1 {
		t.Fatalf("Privileges = %v, want select(col1)", g2.Privileges)
	}
}

// TestGrantRevoke_Loc asserts the statement Loc spans from the GRANT/REVOKE
// keyword through the final grantee.
func TestGrantRevoke_Loc(t *testing.T) {
	sql := "GRANT `select` ON table foo TO 'x'"
	g := parseOneStmt(t, sql).(*ast.GrantStmt)
	if g.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", g.Loc.Start)
	}
	if g.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", g.Loc.End, len(sql))
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
		"GRANT ROLE analyst TO ROLE senior",   // Spanner role form: NOT in the ZetaSQL grammar
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("Parse(%q): want a parse error, got none", sql)
			}
			for _, e := range errs {
				if strings.Contains(e.Msg, "not yet supported") {
					t.Errorf("Parse(%q): got 'not yet supported' (stub still wired): %v", sql, errs)
				}
			}
		})
	}
}

// TestGrantRevoke_RejectsStopAtBoundary asserts a malformed GRANT in a
// multi-statement input does not swallow the following statement — error
// recovery stops at the ';' boundary.
func TestGrantRevoke_RejectsStopAtBoundary(t *testing.T) {
	res := ParseBestEffort("GRANT ON foo TO 'x'; GRANT `select` ON table foo TO 'y'")
	// First statement is malformed (error); second is valid (one stmt).
	if len(res.File.Stmts) != 1 {
		t.Errorf("got %d stmts, want 1 (only the valid second)", len(res.File.Stmts))
	}
	if len(res.Errors) == 0 {
		t.Error("want at least one error from the malformed first statement")
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
