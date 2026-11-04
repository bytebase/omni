package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// The user and privilege statements in the WALK_THROUGH walk.
//
// None creates, drops, or renames a name the walk keeps -- the walk models no
// user, no role and no grant -- so each is followed and nothing is recorded.
// What that has to earn is the statement after it. A node the walk does not
// lower returns nil from lower(), which ends the walk there, so a GRANT in
// the middle of a change script would hide every error below it: a change
// that adds a column twice would go unreported because the walk never reached
// the second statement.
//
// These statements were already listed in the walk, but the forms below are
// new to this parser, so the walk reaches them for the first time here.
//
// These tests need no container and run unconditionally.

// grantWalkStatements are the forms the walk must follow, written against
// walkSchema's table t.
var grantWalkStatements = []string{
	`CREATE USER u1@localhost`,
	`CREATE USER 'u1'@'localhost' IDENTIFIED BY PASSWORD 'notarealhash'`,
	`CREATE USER 'u1'@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132`,
	`CREATE USER 'u1'@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'`,
	`CREATE USER 'u1'@'%' REQUIRE SAN 'DNS:example-user'`,
	`CREATE USER 'u1'@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'`,
	`ALTER USER 'u1'@'localhost' IDENTIFIED BY PASSWORD 'notarealhash'`,
	`DROP USER u1@localhost`,
	"RENAME USER u1@localhost TO `u2`@`localhost`",
	`GRANT SELECT ON t TO u1@localhost`,
	`GRANT ALL PRIVILEGES ON shop.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
	`GRANT GRANT OPTION ON *.* TO u1`,
	`GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO 'u1'@'localhost'`,
	`GRANT 'role1', 'role2' TO 'u1'@'localhost'`,
	`REVOKE SELECT ON t FROM u1@localhost`,
	`REVOKE ALL PRIVILEGES ON shop.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'`,
	`REVOKE 'role1', 'role2' FROM 'u1'@'localhost'`,
	"SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`",
	`SET PASSWORD = 'notarealpw'`,
	`SET PASSWORD FOR u1@localhost = 'notarealpw'`,
	"SET PASSWORD FOR `u1`@`localhost` = 'notarealpw'",
}

func TestWalkThroughFollowsUserAndPrivilegeStatements(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range grantWalkStatements {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing: it changes no name.
			noFinding(t, sql+";", shop)

			// The walk carries on past it, so the error below is still
			// found, and it is reported against statement 2 -- the
			// statement's own index, not the GRANT's.
			oneFinding(t, sql+";\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
				2, "Column `c` already exists in `t`.", shop)

			// And the names a statement above it set are still there
			// below it, so the walk did not lose its state either.
			noFinding(t, "ALTER TABLE t ADD COLUMN c int;\n"+sql+";\nALTER TABLE t DROP COLUMN c;", shop)
		})
	}
}

// TestWalkThroughGrantKeepsTheTableName pins that one of these statements in
// the middle does not make the walk forget a rename either: a GRANT names a
// table, and the walk must read it as a name rather than as a new object.
func TestWalkThroughGrantKeepsTheTableName(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nGRANT SELECT ON u TO u1@localhost;\nALTER TABLE t ADD COLUMN c int;",
		2, "Table `t` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nGRANT SELECT ON u TO u1@localhost;\nCREATE TABLE u (a int);",
		2, "Table `u` already exists.", shop)
}

// TestWalkThroughGrantClaimsNoName pins the other half: these statements
// claim nothing, so a table one of them names that is not there is not a
// finding. v8.5.5 runs `GRANT SELECT ON nosuch TO u1` against a schema with
// no such table, so the walk must not report one.
func TestWalkThroughGrantClaimsNoName(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "GRANT SELECT ON nosuch TO u1@localhost;", shop)
	noFinding(t, "REVOKE SELECT ON nosuch FROM u1@localhost;", shop)
	noFinding(t, "DROP TABLE t;\nGRANT SELECT ON t TO u1@localhost;", shop)
}
