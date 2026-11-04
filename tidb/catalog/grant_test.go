package catalog

import "testing"

// The user and privilege statements reach Exec because a change script holds
// them next to its DDL: a GRANT on the table a migration just made, a CREATE
// USER for the service that will read it, a REVOKE as the last step of
// retiring one.
//
// The catalog models no user, no role and no grant, so what these statements
// must do is run without error and leave everything as it was, and the
// statement after them must still see the table it had. Each is dispatched
// explicitly rather than left to processUtility's default arm, where an
// unmodelled statement and a deliberate no-op read the same.
//
// These tests run unconditionally: they need no container and so are outside
// the TIDB_CATALOG_PARITY quarantine in container_test.go.

func TestExecUserAndPrivilegeStatementsChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`CREATE USER u1@localhost`,
		`CREATE USER 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
		`CREATE USER 'u1'@'localhost' IDENTIFIED BY PASSWORD 'notarealhash'`,
		`CREATE USER 'u1'@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132`,
		`CREATE USER 'u1'@'%' REQUIRE SAN 'DNS:example-user'`,
		`CREATE USER 'u1'@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'`,
		`CREATE USER 'u1'@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'`,
		`ALTER USER 'u1'@'localhost' IDENTIFIED BY PASSWORD 'notarealhash'`,
		`DROP USER u1@localhost`,
		`RENAME USER u1@localhost TO ` + "`u2`@`localhost`",
		`CREATE ROLE 'app_developer'`,
		`DROP ROLE 'app_developer'`,
		`GRANT SELECT ON t TO u1@localhost`,
		`GRANT ALL PRIVILEGES ON testdb.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
		`GRANT GRANT OPTION ON *.* TO u1`,
		`GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO 'u1'@'localhost'`,
		`GRANT 'role1', 'role2' TO 'u1'@'localhost'`,
		`REVOKE SELECT ON t FROM u1@localhost`,
		`REVOKE ALL PRIVILEGES ON testdb.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
		`REVOKE 'role1', 'role2' FROM 'u1'@'localhost'`,
		`SET ROLE ALL`,
		`SET DEFAULT ROLE administrator, developer TO ` + "`joe`@`10.0.0.1`",
		`SET PASSWORD = 'notarealpw'`,
		`SET PASSWORD FOR u1@localhost = 'notarealpw'`,
		"SET PASSWORD FOR `u1`@`localhost` = 'notarealpw'",
	} {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			before := snapshotTable(t, c)
			mustExec(t, c, sql)
			if got := snapshotTable(t, c); got != before {
				t.Errorf("the catalog changed:\nbefore %s\nafter  %s", before, got)
			}

			// The statement is not skipped the way a DML statement is: it
			// is applied, as a no-op.
			results, err := c.Exec(sql, nil)
			if err != nil {
				t.Fatalf("Exec(%q): %v", sql, err)
			}
			if len(results) != 1 {
				t.Fatalf("results = %d, want 1", len(results))
			}
			if results[0].Skipped {
				t.Error("the statement was skipped, want applied as a no-op")
			}
			if results[0].Error != nil {
				t.Errorf("Error = %v, want none", results[0].Error)
			}

			// The DDL that follows still sees the table it had, which is
			// what a catalog that lost the schema here would break.
			mustExec(t, c, `ALTER TABLE t ADD COLUMN c int`)
			if col := c.GetDatabase("testdb").GetTable("t").GetColumn("c"); col == nil {
				t.Error("the column the next statement adds is missing")
			}
		})
	}
}

// TestExecUserStatementsKeepTheirNeighbours pins that these statements in the
// middle of a script change nothing about the statements around them: the
// table one above them made is there below them, and the ALTER below them
// still lands. This is the shape a migration actually has.
func TestExecUserStatementsKeepTheirNeighbours(t *testing.T) {
	c := setupWithDB(t)
	results, err := c.Exec(`CREATE TABLE u (id int);
CREATE USER 'svc'@'10.0.0.1' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132 REQUIRE SAN 'DNS:example-svc';
CREATE ROLE 'reader';
GRANT SELECT, GRANT OPTION ON testdb.* TO 'svc'@'10.0.0.1' IDENTIFIED BY 'notarealpw';
GRANT 'reader' TO 'svc'@'10.0.0.1';
SET DEFAULT ROLE 'reader' TO `+"`svc`@`10.0.0.1`"+`;
REVOKE SELECT ON testdb.* FROM 'svc'@'10.0.0.1';
RENAME USER 'svc'@'10.0.0.1' TO `+"`svc2`@`10.0.0.1`"+`;
DROP USER 'svc2'@'10.0.0.1';
ALTER TABLE u ADD COLUMN v int;`, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != 10 {
		t.Fatalf("results = %d, want 10", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d: %v", i+1, r.Error)
		}
		if r.Skipped {
			t.Errorf("statement %d was skipped, want applied", i+1)
		}
	}
	tbl := c.GetDatabase("testdb").GetTable("u")
	if tbl == nil {
		t.Fatal("table u is missing")
	}
	if tbl.GetColumn("v") == nil {
		t.Error("column v is missing: the ALTER TABLE below the user statements did not land")
	}
}

// TestUserAndPrivilegeExecLines pins the one thing these statements carry
// besides their no-op: their line. Exec reports one ExecResult per statement
// and a caller shows an error against it, and stmtLocStart reads Loc.Start
// through a type switch that answers 0 for a statement it does not list -- so
// a statement missing from that switch reports line 1 wherever it stands, and
// no corpus row can see it.
func TestUserAndPrivilegeExecLines(t *testing.T) {
	stmts := []string{
		"CREATE USER u1@localhost IDENTIFIED BY 'notarealpw'",
		"ALTER USER u1@localhost IDENTIFIED BY PASSWORD 'notarealhash'",
		"CREATE ROLE 'reader'",
		"GRANT SELECT, GRANT OPTION ON *.* TO u1@localhost IDENTIFIED BY 'notarealpw'",
		"GRANT 'reader' TO u1@localhost",
		"SET DEFAULT ROLE 'reader' TO u1@localhost",
		"SET ROLE ALL",
		"SET PASSWORD FOR u1@localhost = 'notarealpw'",
		"REVOKE SELECT ON *.* FROM u1@localhost",
		"REVOKE 'reader' FROM u1@localhost",
		"DROP ROLE 'reader'",
		"RENAME USER u1@localhost TO `u2`@`localhost`",
		"DROP USER u2@localhost",
	}

	// One statement per line, below a statement of another kind, so no
	// statement of this family sits on line 1: a type missing from the
	// switch answers 0, which reads as line 1, and a family statement
	// written there would report the right line by accident.
	sql := "CREATE DATABASE d;\n"
	for _, s := range stmts {
		sql += s + ";\n"
	}

	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != len(stmts)+1 {
		t.Fatalf("got %d results, want %d", len(results), len(stmts)+1)
	}
	for i, r := range results[1:] {
		if r.Error != nil {
			t.Errorf("%s: error %v, want none", stmts[i], r.Error)
		}
		if r.Skipped {
			t.Errorf("%s: skipped, want executed as a utility statement", stmts[i])
		}
		if r.Line != i+2 {
			t.Errorf("%s: line %d, want %d", stmts[i], r.Line, i+2)
		}
	}

	// None of them is a catalog object, so the only database is the one
	// the CREATE DATABASE above them made.
	if got := len(c.Databases()); got != 1 {
		t.Errorf("the catalog holds %d databases, want 1", got)
	}
}
