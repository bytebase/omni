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
// They must also carry their line: Exec reports one ExecResult per statement
// and a caller shows an error against it, and stmtLocStart reads Loc.Start
// through a type switch that answers 0 for a statement it does not list -- so
// a statement missing from that switch reports line 1 wherever it stands.

var userAndPrivilegeStatements = []string{
	`CREATE USER u1@localhost`,
	`CREATE USER 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
	`CREATE USER 'u1'@'%' IDENTIFIED BY PASSWORD 'notarealhash'`,
	`ALTER USER u1@localhost IDENTIFIED BY 'notarealpw'`,
	`DROP USER u1@localhost`,
	"RENAME USER u1@localhost TO `u2`@`localhost`",
	`CREATE ROLE 'reader'`,
	`DROP ROLE 'reader'`,
	`GRANT SELECT ON t TO u1@localhost`,
	`GRANT GRANT OPTION ON *.* TO u1`,
	`GRANT 'reader' TO 'u1'@'localhost'`,
	`REVOKE SELECT ON t FROM u1@localhost`,
	`REVOKE 'reader' FROM 'u1'@'localhost'`,
	`SET ROLE ALL`,
	`SET PASSWORD = 'notarealpw'`,
	`SET PASSWORD FOR u1@localhost = 'notarealpw'`,
}

func TestExecUserAndPrivilegeStatementsChangeNoSchema(t *testing.T) {
	for _, sql := range userAndPrivilegeStatements {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			// The statement is not skipped the way a DML statement is:
			// it is applied, as a no-op.
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

			// The table is as it was, and the DDL that follows still
			// sees it -- which is what a catalog that lost the schema
			// here would break.
			tbl := c.GetDatabase("testdb").GetTable("t")
			if tbl == nil {
				t.Fatal("table t is missing")
			}
			if tbl.GetColumn("a") == nil || tbl.GetColumn("b") == nil {
				t.Error("a column of t is missing")
			}
			if tbl.GetColumn("c") != nil {
				t.Error("t gained a column")
			}
			mustExec(t, c, `ALTER TABLE t ADD COLUMN c int`)
			if col := c.GetDatabase("testdb").GetTable("t").GetColumn("c"); col == nil {
				t.Error("the column the next statement adds is missing")
			}
		})
	}
}

// TestUserAndPrivilegeExecLines pins the line each of them carries. The block
// sits below a statement of another kind so that no statement of this family
// is on line 1: a type missing from the switch answers 0, which reads as line
// 1, and a family statement written there would report the right line by
// accident.
func TestUserAndPrivilegeExecLines(t *testing.T) {
	sql := "CREATE DATABASE d;\n"
	for _, s := range userAndPrivilegeStatements {
		sql += s + ";\n"
	}

	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != len(userAndPrivilegeStatements)+1 {
		t.Fatalf("got %d results, want %d", len(results), len(userAndPrivilegeStatements)+1)
	}
	for i, r := range results[1:] {
		sql := userAndPrivilegeStatements[i]
		if r.Error != nil {
			t.Errorf("%s: error %v, want none", sql, r.Error)
		}
		if r.Skipped {
			t.Errorf("%s: skipped, want executed as a utility statement", sql)
		}
		if r.Line != i+2 {
			t.Errorf("%s: line %d, want %d", sql, r.Line, i+2)
		}
	}

	// None of them is a catalog object, so the only database is the one the
	// CREATE DATABASE above them made.
	if got := len(c.Databases()); got != 1 {
		t.Errorf("the catalog holds %d databases, want 1", got)
	}
}
