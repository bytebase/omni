package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What the user and privilege statements leave in the tree.
//
// The corpus records only whether a statement parses, which is why a family
// can close every one of its rows and still ship a defect. Here the concrete
// risk is the user name: the lexer folds `@word` and "@`word`" into a single
// token, so a reader that takes that token whole leaves CREATE USER
// u1@localhost with the name "u1@localhost" and no host at all -- which
// parses, and which nothing in the corpus says must not happen. The privilege
// and role lists are the other half: a list whose boundaries are lost reads
// the same as one that kept them.

func parseOneCreateUser(t *testing.T, sql string) *ast.CreateUserStmt {
	t.Helper()
	stmt, ok := parseOneStmt(t, sql).(*ast.CreateUserStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	return stmt
}

// TestUserNameSplitsOnTheAt pins that every spelling of user@host comes out
// as the same two halves, with the quoting stripped from each.
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
// that holds dots or a space, a name that holds the '@' itself, and an '@'
// with nothing after it.
func TestUserNameHostShapes(t *testing.T) {
	for _, c := range []struct {
		sql        string
		name, host string
	}{
		{"CREATE USER u1@10.0.0.1", "u1", "10.0.0.1"},
		{"CREATE USER `u1`@`10.0.0.1`", "u1", "10.0.0.1"},
		{"CREATE USER u1@`local host`", "u1", "local host"},
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
	stmt, ok := parseOneStmt(t, "RENAME USER u1@localhost TO `u2`@`other host`").(*ast.RenameUserStmt)
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
	want := [...]string{"u1", "localhost", "u2", "other host"}
	if got != want {
		t.Errorf("pair = %q, want %q", got, want)
	}
}

// TestUserListKeepsItsOwnClauses pins that a user list carries one auth
// option per user rather than one for the list, which is what a reader that
// hoisted the last clause onto every user would break.
func TestUserListKeepsItsOwnClauses(t *testing.T) {
	stmt := parseOneCreateUser(t,
		"CREATE USER 'root'@'localhost' IDENTIFIED BY 'notarealpw', "+
			"'root'@'127.0.0.1' IDENTIFIED BY PASSWORD 'notarealhash'")
	if len(stmt.Users) != 2 {
		t.Fatalf("Users = %d, want 2", len(stmt.Users))
	}
	first, second := stmt.Users[0], stmt.Users[1]
	if first.Host != "localhost" || second.Host != "127.0.0.1" {
		t.Errorf("hosts = %q, %q, want %q, %q", first.Host, second.Host, "localhost", "127.0.0.1")
	}
	// The first user has a password and no hash; the second has a hash and
	// no password. A reader that read BY PASSWORD as BY would give the
	// second user the password "notarealhash".
	if first.Password != "notarealpw" || first.AuthHash != "" {
		t.Errorf("first user: Password = %q, AuthHash = %q, want %q and empty",
			first.Password, first.AuthHash, "notarealpw")
	}
	if second.AuthHash != "notarealhash" || second.Password != "" {
		t.Errorf("second user: AuthHash = %q, Password = %q, want %q and empty",
			second.AuthHash, second.Password, "notarealhash")
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
		{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132", "", "", "0x2A3132"},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password'", "caching_sha2_password", "", ""},
		{"CREATE USER u1@'%' IDENTIFIED WITH caching_sha2_password", "caching_sha2_password", "", ""},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'",
			"caching_sha2_password", "notarealpw", ""},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'",
			"caching_sha2_password", "", "notarealhash"},
		{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132",
			"caching_sha2_password", "", "0x2A3132"},
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

// TestRequireClauseFields pins which REQUIRE option lands in which field, and
// that a list of them keeps every element rather than the last.
func TestRequireClauseFields(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want ast.RequireClause
	}{
		{"CREATE USER u1@'%' REQUIRE NONE", ast.RequireClause{None: true}},
		{"CREATE USER u1@'%' REQUIRE SSL", ast.RequireClause{SSL: true}},
		{"CREATE USER u1@'%' REQUIRE X509", ast.RequireClause{X509: true}},
		{"CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user'",
			ast.RequireClause{SAN: "DNS:example-user"}},
		{"CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'",
			ast.RequireClause{TokenIssuer: "issuer-abc"}},
		// SAN and TOKEN_ISSUER are distinct options, and neither is the
		// ISSUER one.
		{"CREATE USER u1@'%' REQUIRE ISSUER 'i' SAN 'DNS:x' TOKEN_ISSUER 'ti'",
			ast.RequireClause{Issuer: "i", SAN: "DNS:x", TokenIssuer: "ti"}},
		{"CREATE USER u1@'%' REQUIRE CIPHER 'c' AND SUBJECT 's' AND ISSUER 'i'",
			ast.RequireClause{Cipher: "c", Subject: "s", Issuer: "i"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseOneCreateUser(t, c.sql)
			if stmt.Require == nil {
				t.Fatal("Require is nil")
			}
			got := *stmt.Require
			got.Loc = ast.Loc{}
			if got != c.want {
				t.Errorf("Require = %+v, want %+v", got, c.want)
			}
		})
	}
}

// TestPrivilegeListElements pins the privilege list element by element: the
// two-word names are one element each, the long names keep both of their
// words, and nothing merges with its neighbour.
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
		{"GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA, INVOKE SAGEMAKER, INVOKE COMPREHEND ON *.* TO u1",
			[]string{"LOAD FROM S3", "SELECT INTO S3", "INVOKE LAMBDA", "INVOKE SAGEMAKER", "INVOKE COMPREHEND"}},
		{"GRANT LOAD FROM s3 ON *.* TO u1", []string{"LOAD FROM S3"}},
		{"GRANT INVOKE `LAMBDA` ON *.* TO u1", []string{"INVOKE LAMBDA"}},
		{"GRANT BACKUP_ADMIN, SELECT ON *.* TO u1", []string{"BACKUP_ADMIN", "SELECT"}},
		{"GRANT SELECT (a, b), INSERT (c) ON t TO u1", []string{"SELECT(a, b)", "INSERT(c)"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, c.sql).(*ast.GrantStmt)
			if !ok {
				t.Fatal("Parse gave a statement of another kind")
			}
			assertStrings(t, "Privileges", stmt.Privileges, c.want)
		})
	}
}

// TestRoleListElements pins the role list the same way, including the host
// each role may carry.
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
		{"GRANT INVOKE LAMBDA TO u1@localhost",
			[]string{"INVOKE LAMBDA"}, []string{"u1@localhost"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, c.sql).(*ast.GrantRoleStmt)
			if !ok {
				t.Fatal("Parse gave a statement of another kind")
			}
			assertStrings(t, "Roles", stmt.Roles, c.roles)
			assertStrings(t, "To", stmt.To, c.to)
		})
	}
}

// TestRevokeRoleListElements pins the REVOKE side of the same list.
func TestRevokeRoleListElements(t *testing.T) {
	stmt, ok := parseOneStmt(t,
		"REVOKE 'role1', 'role2' FROM 'user1'@'localhost', 'user2'@'localhost'").(*ast.RevokeRoleStmt)
	if !ok {
		t.Fatal("Parse gave a statement of another kind")
	}
	assertStrings(t, "Roles", stmt.Roles, []string{"role1", "role2"})
	assertStrings(t, "From", stmt.From, []string{"user1@localhost", "user2@localhost"})
}

// TestGranteeAuthOptionIsDropped pins that a grantee's auth option is read
// and dropped rather than folded into the grantee's name, which is how an
// `IDENTIFIED BY` tail would otherwise reach a GRANT's user list.
func TestGranteeAuthOptionIsDropped(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"grant all privileges on zabbix.* to 'zabbix'@'localhost' identified by 'notarealpw'",
			[]string{"zabbix@localhost"}},
		{"GRANT PROCESS ON *.* TO 'u1'@'%' identified by password 'notarealhash'",
			[]string{"u1@%"}},
		{"GRANT SELECT ON t TO u1 IDENTIFIED BY 'notarealpw', u2 IDENTIFIED BY 'notarealpw2'",
			[]string{"u1", "u2"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, c.sql).(*ast.GrantStmt)
			if !ok {
				t.Fatal("Parse gave a statement of another kind")
			}
			assertStrings(t, "To", stmt.To, c.want)
		})
	}
}

// TestSetDefaultRoleListElements pins the role and user lists of SET DEFAULT
// ROLE, where both sides are lists.
func TestSetDefaultRoleListElements(t *testing.T) {
	stmt, ok := parseOneStmt(t,
		"SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`").(*ast.SetDefaultRoleStmt)
	if !ok {
		t.Fatal("Parse gave a statement of another kind")
	}
	assertStrings(t, "Roles", stmt.Roles, []string{"administrator", "developer"})
	assertStrings(t, "To", stmt.To, []string{"joe@10.0.0.1"})
}

func assertStrings(t *testing.T, field string, got, want []string) {
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

// TestPrivilegeNamesAreUppercased pins the one normalisation this family
// does: a privilege name is folded to upper case, so a script that writes
// them in lower case reads the same.
func TestPrivilegeNamesAreUppercased(t *testing.T) {
	stmt, ok := parseOneStmt(t, "grant select, invoke lambda on *.* to u1").(*ast.GrantStmt)
	if !ok {
		t.Fatal("Parse gave a statement of another kind")
	}
	for _, name := range stmt.Privileges {
		if name != strings.ToUpper(name) {
			t.Errorf("privilege %q is not upper case", name)
		}
	}
}
