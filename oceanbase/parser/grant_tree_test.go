package parser

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// What a user name leaves in the tree.
//
// The lexer folds `@word` and "@`word`" into a single token, so a reader that
// takes that token whole leaves CREATE USER u1@localhost with the name
// "u1@localhost" and no host at all -- which parses, and which no
// accept-or-reject record can see. tidb/parser/grant_tree_test.go pins the
// same split against the same spellings.

func parseOneUserStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

func parseOneCreateUser(t *testing.T, sql string) *ast.CreateUserStmt {
	t.Helper()
	stmt, ok := parseOneUserStmt(t, sql).(*ast.CreateUserStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	return stmt
}

// TestUserNameSplitsOnTheAt pins that every spelling of user@host comes out as
// the same two halves, with the quoting stripped from each.
func TestUserNameSplitsOnTheAt(t *testing.T) {
	for _, sql := range []string{
		"CREATE USER u1@localhost",
		"CREATE USER `u1`@localhost",
		"CREATE USER 'u1'@localhost",
		"CREATE USER u1@`localhost`",
		"CREATE USER `u1`@`localhost`",
		"CREATE USER 'u1'@`localhost`",
		"CREATE USER u1@'localhost'",
		"CREATE USER `u1`@'localhost'",
		"CREATE USER 'u1'@'localhost'",
		"CREATE USER u1 @localhost",
		"CREATE USER u1 @`localhost`",
		"CREATE USER u1 @'localhost'",
	} {
		t.Run(sql, func(t *testing.T) {
			stmt := parseOneCreateUser(t, sql)
			if len(stmt.Users) != 1 {
				t.Fatalf("Users = %d, want 1", len(stmt.Users))
			}
			if got := stmt.Users[0].Name; got != "u1" {
				t.Errorf("Name = %q, want %q", got, "u1")
			}
			if got := stmt.Users[0].Host; got != "localhost" {
				t.Errorf("Host = %q, want %q", got, "localhost")
			}
		})
	}
}

// TestUserNameHostShapes pins the halves that are not the plain case: a host
// that holds dots, a name that holds the '@' itself, and an '@' with nothing
// after it.
func TestUserNameHostShapes(t *testing.T) {
	for _, c := range []struct {
		sql        string
		name, host string
	}{
		{"CREATE USER u1@10.0.0.1", "u1", "10.0.0.1"},
		{"CREATE USER `u1`@`10.0.0.1`", "u1", "10.0.0.1"},
		{"CREATE USER `user@example.com`@'localhost'", "user@example.com", "localhost"},
		// The host is never read as a keyword, since it arrives on the
		// '@' token rather than on its own.
		{"CREATE USER u1@select", "u1", "select"},
		{"CREATE USER u1@comment", "u1", "comment"},
		// An '@' with nothing after it is an empty host, and not a name
		// that ends in '@'.
		{"CREATE USER u1@", "u1", ""},
		// Case is kept as written, on both halves.
		{"CREATE USER `UEsr1`@`LOCALhost`", "UEsr1", "LOCALhost"},
		{"CREATE USER uesr1@LOCALhost", "uesr1", "LOCALhost"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseOneCreateUser(t, c.sql)
			if len(stmt.Users) != 1 {
				t.Fatalf("Users = %d, want 1", len(stmt.Users))
			}
			if got := stmt.Users[0].Name; got != c.name {
				t.Errorf("Name = %q, want %q", got, c.name)
			}
			if got := stmt.Users[0].Host; got != c.host {
				t.Errorf("Host = %q, want %q", got, c.host)
			}
		})
	}
}

// TestRenameUserSplitsBothSides pins the same split on the statement that
// holds two of them.
func TestRenameUserSplitsBothSides(t *testing.T) {
	stmt, ok := parseOneUserStmt(t, "RENAME USER u1@localhost TO `u2`@`other.host`").(*ast.RenameUserStmt)
	if !ok {
		t.Fatal("Parse gave a statement of another kind")
	}
	if len(stmt.Pairs) != 1 {
		t.Fatalf("Pairs = %d, want 1", len(stmt.Pairs))
	}
	got := [...]string{
		stmt.Pairs[0].OldUser, stmt.Pairs[0].OldHost,
		stmt.Pairs[0].NewUser, stmt.Pairs[0].NewHost,
	}
	want := [...]string{"u1", "localhost", "u2", "other.host"}
	if got != want {
		t.Errorf("pair = %q, want %q", got, want)
	}
}

// TestPrivilegeListElements pins the privilege list element by element: the
// two-word names are one element each and nothing merges with its neighbour.
func TestPrivilegeListElements(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"GRANT SELECT ON *.* TO u1", []string{"SELECT"}},
		{"GRANT GRANT OPTION ON *.* TO u1", []string{"GRANT OPTION"}},
		{"GRANT GRANT OPTION, SELECT ON *.* TO u1", []string{"GRANT OPTION", "SELECT"}},
		{"GRANT SELECT, GRANT OPTION ON *.* TO u1", []string{"SELECT", "GRANT OPTION"}},
		{"grant PROCESS,usage, REPLICATION SLAVE, REPLICATION CLIENT on *.* to u1",
			[]string{"PROCESS", "USAGE", "REPLICATION SLAVE", "REPLICATION CLIENT"}},
		{"GRANT SELECT (a, b), INSERT (c) ON t TO u1", []string{"SELECT(a, b)", "INSERT(c)"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt, ok := parseOneUserStmt(t, c.sql).(*ast.GrantStmt)
			if !ok {
				t.Fatal("Parse gave a statement of another kind")
			}
			assertUserStrings(t, "Privileges", stmt.Privileges, c.want)
		})
	}
}

// TestRoleListElements pins the role list the same way.
func TestRoleListElements(t *testing.T) {
	for _, c := range []struct {
		sql       string
		roles, to []string
	}{
		{"GRANT 'role1' TO 'u1'@'localhost'",
			[]string{"role1"}, []string{"u1@localhost"}},
		{"GRANT 'role1', 'role2' TO 'u1'@'LOCalhost', 'u2'@'LOcalhost'",
			[]string{"role1", "role2"}, []string{"u1@LOCalhost", "u2@LOcalhost"}},
		{"GRANT 'app_read'@'%','app_write'@'%' TO 'rw_user1'@'localhost'",
			[]string{"app_read@%", "app_write@%"}, []string{"rw_user1@localhost"}},
		{"GRANT role1@localhost, 'role2' TO u1",
			[]string{"role1@localhost", "role2"}, []string{"u1"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt, ok := parseOneUserStmt(t, c.sql).(*ast.GrantRoleStmt)
			if !ok {
				t.Fatal("Parse gave a statement of another kind")
			}
			assertUserStrings(t, "Roles", stmt.Roles, c.roles)
			assertUserStrings(t, "To", stmt.To, c.to)
		})
	}
}

// TestAuthOptionFields pins which field each arm of the auth option lands in.
// Every arm writes a credential, so reading one as another is the defect that
// matters here.
func TestAuthOptionFields(t *testing.T) {
	for _, c := range []struct {
		sql                    string
		plugin, password, hash string
	}{
		{"CREATE USER u1@'%' IDENTIFIED BY 'notarealpw'", "", "notarealpw", ""},
		{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'", "", "", "notarealhash"},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'",
			"caching_sha2_password", "notarealpw", ""},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY PASSWORD 'notarealhash'",
			"caching_sha2_password", "", "notarealhash"},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'",
			"caching_sha2_password", "", "notarealhash"},
		// An empty hash is a hash, not a missing one.
		{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS ''",
			"mysql_native_password", "", ""},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseOneCreateUser(t, c.sql)
			if len(stmt.Users) != 1 {
				t.Fatalf("Users = %d, want 1", len(stmt.Users))
			}
			u := stmt.Users[0]
			if u.AuthPlugin != c.plugin {
				t.Errorf("AuthPlugin = %q, want %q", u.AuthPlugin, c.plugin)
			}
			if u.Password != c.password {
				t.Errorf("Password = %q, want %q", u.Password, c.password)
			}
			if u.AuthHash != c.hash {
				t.Errorf("AuthHash = %q, want %q", u.AuthHash, c.hash)
			}
		})
	}
}

func assertUserStrings(t *testing.T, field string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %d entries %q, want %d %q", field, len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d] = %q, want %q (whole list %q)", field, i, got[i], want[i], got)
		}
	}
}

// TestGrantCatalogObjectType holds CATALOG to its two readings: the object type
// before a privilege level, and a name before TO, FROM or a dot.
func TestGrantCatalogObjectType(t *testing.T) {
	for _, tc := range []struct{ sql, typ, schema, name string }{
		{"GRANT SELECT ON CATALOG * TO u", "CATALOG", "", "*"},
		{"GRANT SELECT ON catalog TO u", "", "", "catalog"},
		{"GRANT SELECT ON catalog.* TO u", "", "catalog", "*"},
		{"GRANT SELECT ON catalog.t TO u", "", "catalog", "t"},
		{"REVOKE SELECT ON catalog.t FROM u", "", "catalog", "t"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		var on *ast.GrantTarget
		switch s := list.Items[0].(type) {
		case *ast.GrantStmt:
			on = s.On
		case *ast.RevokeStmt:
			on = s.On
		}
		if on == nil || on.Type != tc.typ || on.Name.Schema != tc.schema || on.Name.Name != tc.name {
			t.Errorf("Parse(%q).On = %+v, want type %q on %q.%q", tc.sql, on, tc.typ, tc.schema, tc.name)
		}
	}
}

// TestRoleGrantNames holds a role grant's names to unreserved_keyword_for_role_name:
// every unreserved keyword but unreserved_keyword_ambiguous_roles, so RESTART names
// a role and SUPER, a privilege word, and DUAL, a reserved one, do not.
func TestRoleGrantNames(t *testing.T) {
	for _, sql := range []string{"GRANT restart TO u", "REVOKE restart FROM u", "GRANT r1, `super` TO u"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{"GRANT super TO u", "GRANT DUAL TO u", "GRANT OPTION TO u", "GRANT select TO u"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
