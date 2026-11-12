package parser

import "testing"

// The forms the user and privilege statements take on this engine, and -- the
// point of this file -- the forms they do not.
//
// This is the same survey of the family's surface that tidb/parser records,
// run against oceanbase-ce:4.3.5, so the two files read side by side: where a
// verdict differs, the engines differ. The quoting of a user name is where the
// family is easiest to get wrong, because the lexer folds `@word` and
// "@`word`" into one token and leaves `@'word'` as two, so every combination
// is here: which half may be bare, quoted or backquoted, where whitespace is
// allowed around the '@', and what an '@' with nothing after it means.
//
// Every case below was measured against oceanbase-ce:4.3.5 through PREPARE on
// 2026-11-04, and the recorded verdict is what that server answered.
// TestGrantFormsOracle is what keeps it a recording: an image upgrade that
// changes an answer fails there rather than leaving a stale expectation
// behind.

// grantFormCase is a statement and a verdict on it.
type grantFormCase struct {
	sql    string
	accept bool
}

var grantFormCases = []grantFormCase{
	// --- user@host: which half may be written how ---
	// Username is `StringName ['@' StringName]`, and StringName is a string
	// literal or an identifier, so both halves take all three spellings.
	{"CREATE USER u1@localhost", true},
	{"CREATE USER `u1`@localhost", true},
	{"CREATE USER 'u1'@localhost", true},
	{"CREATE USER u1@`localhost`", true},
	{"CREATE USER `u1`@`localhost`", true},
	{"CREATE USER 'u1'@`localhost`", true},
	{"CREATE USER u1@'localhost'", true},
	{"CREATE USER `u1`@'localhost'", true},
	{"CREATE USER 'u1'@'localhost'", true},
	// A backquoted name holds an '@' of its own, which is how an account
	// named after an email address is written.
	{"CREATE USER `user@example.com`@'localhost'", true},
	// A host is a name, so it holds dots and may be backquoted around a
	// space; neither makes it two names.
	{"CREATE USER u1@10.0.0.1", true},
	{"CREATE USER u1@`local host`", false},
	// The name is a StringName and so not a reserved word, while the host
	// rides in on the '@' token and is never read as one.
	{"CREATE USER select@localhost", false},
	{"CREATE USER comment@localhost", true},
	{"CREATE USER u1@select", true},
	{"CREATE USER u1@comment", true},
	// The '@' has to touch its host: `@word` and "@`word`" are one token,
	// and the two-token spelling is only reached when the host is quoted.
	{"CREATE USER u1 @localhost", true},
	{"CREATE USER u1 @`localhost`", true},
	{"CREATE USER u1 @'localhost'", true},
	{"CREATE USER u1@ localhost", false},
	{"CREATE USER u1@ 'localhost'", false},
	{"CREATE USER u1 @ localhost", false},
	{"CREATE USER u1 @ 'localhost'", false},
	// An '@' with nothing after it is a user with an empty host.
	{"CREATE USER u1@", true},
	{"DROP USER u1@", true},
	{"ALTER USER u1@ IDENTIFIED BY 'notarealpw'", true},
	{"GRANT SELECT ON *.* TO u1@", true},
	{"SET DEFAULT ROLE role1 TO u1@", true},
	// The same spellings reach every statement that names a user.
	{"DROP USER u1@localhost, `u2`@`localhost`", true},
	{"RENAME USER u1@localhost TO `u2`@`localhost`", true},
	{"SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`", true},

	// --- IDENTIFIED: the auth options ---
	// AuthPlugin is a StringName, so the plugin may be quoted as well as
	// bare. A backquoted plugin is an identifier and so a plugin too.
	{"CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password'", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH `mysql_native_password`", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS ''", true},
	// HashString is a string literal or a hex literal; the hex spelling is
	// what a dumped caching_sha2_password account carries.
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132333435", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password AS 0x2A3132333435", false},
	// BY PASSWORD names the stored hash. MySQL 8.0 dropped the form.
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132333435", false},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD ''", true},
	{"ALTER USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'", false},
	// Each arm needs its string, and no arm combines with another.
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY", false},
	{"CREATE USER u1@'%' IDENTIFIED BY", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 123", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY PASSWORD 'notarealhash'", true},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD PASSWORD 'notarealhash'", false},
	// A user list carries one auth option per user, not one for the list.
	{"CREATE USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"ALTER USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'", false},

	// --- REQUIRE ---
	// NONE, SSL and X509 each stand alone, and the value-carrying options
	// make a list that AND joins optionally.
	{"CREATE USER u1@'%' REQUIRE NONE", true},
	{"CREATE USER u1@'%' REQUIRE SSL", true},
	{"CREATE USER u1@'%' REQUIRE X509", true},
	{"CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA'", true},
	{"CREATE USER u1@'%' REQUIRE SUBJECT 's' ISSUER 'i' CIPHER 'c'", true},
	{"CREATE USER u1@'%' REQUIRE SUBJECT 's' AND ISSUER 'i'", true},
	{"CREATE USER u1@'%' REQUIRE CIPHER 'c' AND SUBJECT 's' AND ISSUER 'i'", true},
	// SAN and TOKEN_ISSUER are TiDB's own, and belong to the list.
	{"CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user, URI:spiffe://example.org/myservice'", false},
	{"CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'", false},
	{"CREATE USER u1@'%' REQUIRE token_issuer 'issuer-abc'", false},
	{"CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA' SAN 'DNS:example-user'", false},
	{"CREATE USER u1@'%' REQUIRE ISSUER 'i' AND SAN 'DNS:example-user'", false},
	{"CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'", false},
	{"CREATE USER u1@'%' REQUIRE ISSUER 'i' SUBJECT 's' CIPHER 'c' SAN 'DNS:x' TOKEN_ISSUER 'ti'", false},
	// A standalone option is standalone: it joins no list, with or without
	// the AND.
	{"CREATE USER u1@'%' REQUIRE SSL AND X509", false},
	{"CREATE USER u1@'%' REQUIRE SSL AND SAN 'DNS:example-user'", false},
	{"CREATE USER u1@'%' REQUIRE X509 SAN 'DNS:example-user'", false},
	{"CREATE USER u1@'%' REQUIRE NONE SSL", false},
	{"CREATE USER u1@'%' REQUIRE NONE AND SSL", false},
	// The list is never empty, every element carries a string literal, and
	// an AND is never trailing.
	{"CREATE USER u1@'%' REQUIRE", false},
	{"CREATE USER u1@'%' REQUIRE CIPHER", false},
	{"CREATE USER u1@'%' REQUIRE SAN", false},
	{"CREATE USER u1@'%' REQUIRE TOKEN_ISSUER", false},
	{"CREATE USER u1@'%' REQUIRE CIPHER 'c' AND", false},
	{"CREATE USER u1@'%' REQUIRE SAN 0x2A31", false},
	// SAN and TOKEN_ISSUER are the words the rule names, not names: a
	// backquoted one is an identifier and no option at all.
	{"CREATE USER u1@'%' REQUIRE `SAN` 'DNS:example-user'", false},
	{"CREATE USER u1@'%' REQUIRE `TOKEN_ISSUER` 'issuer-abc'", false},

	// --- the account options, which the auth clauses sit in front of ---
	{"CREATE USER u1@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{\"email\": \"user@example.com\"}'", false},

	// --- GRANT and REVOKE: the grantee list ---
	// The privilege forms take a UserSpecList, so each grantee may carry an
	// auth option of its own.
	{"GRANT ALL PRIVILEGES ON zz.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'", true},
	{"GRANT PROCESS, USAGE ON *.* TO 'u1'@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"REVOKE ALL PRIVILEGES ON zz.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'", false},
	{"REVOKE SELECT ON *.* FROM u1 IDENTIFIED BY 'notarealpw'", false},
	{"GRANT ALL ON *.* TO u1 IDENTIFIED BY 'notarealpw' REQUIRE SSL WITH GRANT OPTION", false},
	// A GRANT carries neither a REQUIRE clause nor a resource option
	// here; MariaDB and TiDB take both.
	{"GRANT ALL ON *.* TO u1 REQUIRE SSL", false},
	{"GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 10", false},
	{"GRANT ALL ON *.* TO u1 REQUIRE SAN 'DNS:example-user'", false},
	// The role forms take a UsernameList instead, and so does GRANT PROXY.
	{"GRANT 'role1' TO u1 IDENTIFIED BY 'notarealpw'", false},
	{"REVOKE 'role1' FROM u1 IDENTIFIED BY 'notarealpw'", false},
	{"GRANT PROXY ON u1 TO u2 IDENTIFIED BY 'notarealpw'", false},

	// --- GRANT: the WITH clauses (measured 2026-11-12) ---
	// grant_options is `WITH GRANT OPTION`, on a privilege grant and on GRANT
	// PROXY, and opt_with_admin_option `WITH ADMIN OPTION`
	// (sql_parser_mysql_mode.y). A WITH takes all of its words.
	{"GRANT ALL ON *.* TO u1 WITH GRANT OPTION", true},
	{"GRANT ALL ON *.* TO u1 WITH GRANT", false},
	{"GRANT ALL ON *.* TO u1 WITH", false},
	{"GRANT ALL ON *.* TO u1 WITH OPTION", false},
	{"GRANT SELECT ON db1.t TO u1 WITH GRANT", false},
	{"GRANT ALL ON *.* TO u1 WITH GRANT OPTION GRANT OPTION", false},
	// GRANT, ADMIN and OPTION are the words the rules name, not names.
	{"GRANT ALL ON *.* TO u1 WITH `GRANT` OPTION", false},
	{"GRANT ALL ON *.* TO u1 WITH GRANT `OPTION`", false},
	{"GRANT role1 TO u1 WITH `ADMIN` OPTION", false},
	{"GRANT PROXY ON u1 TO u2 WITH GRANT OPTION", true},
	{"GRANT PROXY ON u1 TO u2 WITH GRANT", false},
	{"GRANT PROXY ON u1 TO u2 WITH", false},
	{"GRANT role1 TO u1 WITH ADMIN OPTION", true},
	{"GRANT role1 TO u1 WITH ADMIN", false},
	{"GRANT role1 TO u1 WITH", false},
	{"GRANT role1 TO u1 WITH GRANT OPTION", false},
	// AS user and its WITH ROLE are MySQL 8.0's: a privilege grant ends at its
	// WITH clause here.
	{"GRANT ALL ON *.* TO u1 AS u2", false},
	{"GRANT ALL ON *.* TO u1 AS 'u2'@'%'", false},
	{"GRANT SELECT ON db1.t TO u1 AS u2", false},
	{"GRANT ALL ON *.* TO u1 AS u2 WITH ROLE DEFAULT", false},
	{"GRANT ALL ON *.* TO u1 AS u2 WITH ROLE NONE", false},
	{"GRANT ALL ON *.* TO u1 AS u2 WITH ROLE ALL EXCEPT role1", false},
	{"GRANT ALL ON *.* TO u1 AS u2 WITH ROLE role1, role2", false},
	{"GRANT ALL ON *.* TO u1 WITH GRANT OPTION AS u2", false},

	// --- GRANT and REVOKE: role or privilege ---
	// A string literal is a role name and nothing else, with or without an
	// @host of its own.
	{"GRANT 'role1' TO 'u1'@'localhost'", true},
	{"GRANT 'role1', 'role2' TO 'u1'@'localhost', 'u2'@'localhost'", true},
	{"GRANT 'role1'@'%', 'role2'@'%' TO 'u1'@'localhost'", true},
	{"GRANT 'u1' TO 'u1'", true},
	{"REVOKE 'role1' FROM u1", true},
	{"REVOKE 'role1', 'role2' FROM 'u1'@'localhost', 'u2'@'localhost'", true},
	{"SET DEFAULT ROLE a, b TO u1@localhost", true},
	// So the privilege reading is closed to it.
	{"GRANT 'role1' ON *.* TO u1", true},
	{"GRANT 'role1'@'%' ON *.* TO u1", true},
	{"REVOKE 'role1' ON *.* FROM u1", true},
	{"GRANT SELECT, 'role1' ON *.* TO u1", true},
	// A word with an @host after it is a role too, keyword or not, and the
	// privilege reading is closed to it as well.
	{"GRANT role1@localhost TO u1", true},
	{"GRANT role1@`localhost` TO u1", true},
	{"GRANT role1@'localhost' TO u1", true},
	{"GRANT comment@localhost TO u1", true},
	{"REVOKE role1@localhost FROM u1", true},
	{"REVOKE comment@localhost FROM u1", true},
	{"GRANT role1@localhost ON *.* TO u1", true},
	{"GRANT comment@localhost ON *.* TO u1", true},
	// A word is a role name only when it is an identifier. A keyword
	// reduces to a privilege and the role reading is closed to it, which is
	// also why the two-word privileges are no role names.
	{"GRANT role1 TO u1", true},
	{"GRANT `select` TO u1", true},
	{"GRANT COMMENT TO u1", true},
	{"GRANT CREATE VIEW TO u1", false},
	{"GRANT REPLICATION SLAVE TO u1", false},
	{"GRANT GRANT OPTION TO u1", false},

	// --- the privilege names ---
	// GRANT OPTION is the privilege; GRANT on its own is not one.
	{"GRANT GRANT OPTION ON *.* TO u1", true},
	{"GRANT GRANT OPTION ON `db`.* TO u1", true},
	{"GRANT GRANT OPTION, SELECT ON *.* TO u1", true},
	{"GRANT SELECT, GRANT OPTION ON *.* TO u1", true},
	{"GRANT GRANT ON *.* TO u1", false},
	// ExtendedPriv is a run of identifiers, which is how the long names are
	// spelled.
	{"GRANT INVOKE LAMBDA ON *.* TO u1", false},
	{"GRANT INVOKE SAGEMAKER ON *.* TO u1", false},
	{"GRANT INVOKE COMPREHEND ON *.* TO u1", false},
	{"GRANT INVOKE LAMBDA, INVOKE SAGEMAKER ON *.* TO u1", false},
	{"GRANT BACKUP_ADMIN ON *.* TO u1", true},
	{"GRANT SELECT, BACKUP_ADMIN ON *.* TO u1", true},
	// The run takes identifiers, backquoted or not, and nothing else.
	{"GRANT INVOKE `LAMBDA` ON *.* TO u1", false},
	// LOAD FROM S3 and SELECT INTO S3 are spelled out in the grammar
	// because their first two words are reserved, so neither reaches the
	// run of identifiers. S3 is the word the rule names: a backquoted one
	// is an identifier, and the two words do not recombine.
	{"GRANT LOAD FROM S3 ON *.* TO u1", false},
	{"GRANT SELECT INTO S3 ON *.* TO u1", false},
	{"GRANT LOAD FROM s3 ON *.* TO u1", false},
	{"GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1", false},
	{"GRANT LOAD FROM `S3` ON *.* TO u1", false},
	{"GRANT LOAD S3 ON *.* TO u1", false},
	{"GRANT SELECT S3 ON *.* TO u1", false},
	{"GRANT LOAD FROM SELECT ON *.* TO u1", false},
	// An ExtendedPriv is a role name too, since the reading is still open
	// when the TO arrives.
	{"GRANT LOAD FROM S3 TO u1", false},
	{"GRANT INVOKE LAMBDA TO u1", false},
	// --- a name that spells a privilege word, and a clause that is not one ---
	// A word written backquoted is a name, not the keyword, so it is a role
	// name even where the keyword is a privilege or is reserved. The
	// multi-word completions are closed to it for the same reason.
	{"GRANT `grant` TO u1", true},
	{"REVOKE `grant` FROM u1", true},
	{"GRANT `option` TO u1", true},
	{"GRANT `select` TO u1", true},
	// A backquoted name holding a space is one name, and so still a role.
	{"GRANT `grant option` TO u1", true},
	{"GRANT `app reader` TO u1", true},
	{"REVOKE `app reader` FROM u1", true},
	{"GRANT `app reader` ON *.* TO u1", true},
	// So the second word of a two-word privilege is not read after one: the
	// VIEW and the SLAVE below are left over, and the statement ends there.
	{"GRANT `create` VIEW ON *.* TO u1", false},
	{"GRANT `replication` SLAVE ON *.* TO u1", false},
	// An @host is an @host only where the source says so. The lexer strips
	// the backquotes, so "`@localhost`" reaches this parser as an
	// identifier whose text is "@localhost", exactly like a real one.
	{"CREATE USER 'u1' `@localhost`", false},
	{"DROP USER 'u1' `@localhost`", false},
	// A bare IDENTIFIED is not a clause, on a user spec or on a grantee.
	{"CREATE USER u1@'%' IDENTIFIED", false},
	{"GRANT SELECT ON *.* TO u1 IDENTIFIED", false},
	{"REVOKE SELECT ON *.* FROM u1 IDENTIFIED", false},
	{"GRANT SELECT ON *.* TO u1 IDENTIFIED BY", false},
	// CURRENT_USER names the connected user wherever a user is named, but
	// oceanbase-ce:4.3.5 answers 1064 for it as a RENAME USER target and after
	// SET PASSWORD FOR (measured), where mysql:8.0 reads it. The positions that
	// do accept it execute against the connected account, so they are measured by
	// probe rather than run here; these two stay rejected on both sides.
	{"RENAME USER CURRENT_USER TO CURRENT_USER", false},
	{"RENAME USER u1@localhost TO CURRENT_USER", false},
	{"SET PASSWORD FOR CURRENT_USER = 'notarealpw'", false},
	// CURRENT_USER is a Username and not a Rolename, so oceanbase-ce:4.3.5 answers
	// 1064 for the keyword wherever a role is named. DROP ROLE current_user is the
	// one the engine reads, and the accepting positions execute against a real
	// account, so neither is run here.
	{"CREATE ROLE current_user", false},
	{"SET ROLE current_user", false},
	{"SET DEFAULT ROLE current_user TO u1@localhost", false},
}

// TestGrantForms holds this parser to the verdicts above. It needs no
// container.
func TestGrantForms(t *testing.T) {
	for _, c := range grantFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.accept && err != nil {
				t.Errorf("Parse(%q): %v, want accepted", c.sql, err)
			}
			if !c.accept && err == nil {
				t.Errorf("Parse(%q) was accepted, want rejected", c.sql)
			}
		})
	}
}

// grantKnownDivergences records the statements of this family this parser
// still answers differently from oceanbase-ce:4.3.5, with what this parser
// answers, so each is a measured and named gap rather than an untested one.
// Every one of them needs knowledge this parser does not hold -- which words
// that server counts as privilege names, which user options it has -- rather
// than a clause of the family this change closed, so each belongs to its own
// change.
var grantKnownDivergences = []grantFormCase{
	// The password and lock options and ATTRIBUTE are MySQL 8.0's, and
	// OceanBase has none of them.
	{"CREATE USER u1@'%' PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK", true},
	{"CREATE USER u1@'%' ATTRIBUTE '{\"email\": \"user@example.com\"}'", true},
	// A privilege keyword is no role name, but this parser holds no list
	// of which keywords are privileges; USAGE and PROCESS reach it as
	// identifiers, so neither is read as the privilege it is.
	{"GRANT SELECT TO u1", true},
	{"GRANT usage TO u1", true},
	{"GRANT process TO u1", true},
}

// TestGrantKnownDivergences holds the list above to this parser, so a case
// that has started agreeing with oceanbase-ce:4.3.5 fails here and is moved
// into grantFormCases. It needs no container; the oracle test is what says
// the server has not changed its mind.
func TestGrantKnownDivergences(t *testing.T) {
	for _, c := range grantKnownDivergences {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("this parser now %ss %q; the list records %s, so the case belongs in grantFormCases (err=%v)",
					grantVerdict(accepts), c.sql, grantVerdict(c.accept), err)
			}
		})
	}
}

// grantVerdict names a parse outcome, for the message above.
func grantVerdict(accepts bool) string {
	if accepts {
		return "accept"
	}
	return "reject"
}
