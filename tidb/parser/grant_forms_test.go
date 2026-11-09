package parser

import "testing"

// The forms the user and privilege statements take, and -- the point of this
// file -- the forms they do not.
//
// The corpus in testdata/pingcap_parser_corpus.tsv is pingcap's own parser
// test, and for this family it is almost all accepting rows: 30 of them moved
// when this family was closed and none of them is a rejection. A reader that
// took any word after IDENTIFIED WITH, any option after REQUIRE and any
// spelling of user@host would turn every one of those rows green while
// accepting `REQUIRE SSL AND X509` and `CREATE USER u@ 'h'`, so the
// rejections are written here by hand.
//
// The quoting of a user name is where this family is easiest to get wrong,
// because the lexer folds `@word` and "@`word`" into one token and leaves
// `@'word'` as two. So every combination is here: which half may be bare,
// quoted or backquoted, where whitespace is allowed around the '@', and what
// an '@' with nothing after it means.
//
// Every case below was measured against pingcap/tidb:v8.5.5 through PREPARE
// on 2026-11-04, and the recorded verdict is what that server answered.
// TestGrantFormsMatchTiDB is what keeps it a recording: an image upgrade that
// changes an answer fails there rather than leaving a stale expectation
// behind.
//
// Ref: Username, UserSpec, AuthOption, RequireClause, RoleOrPrivElem and
// PrivType in pingcap/tidb pkg/parser/parser.y.

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
	{"CREATE USER u1@`local host`", true},
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
	{"CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH `mysql_native_password`", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS ''", true},
	// HashString is a string literal or a hex literal; the hex spelling is
	// what a dumped caching_sha2_password account carries.
	{"CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132333435", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password AS 0x2A3132333435", true},
	// BY PASSWORD names the stored hash. MySQL 8.0 dropped the form.
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132333435", true},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD ''", true},
	{"ALTER USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	// Each arm needs its string, and no arm combines with another.
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY", false},
	{"CREATE USER u1@'%' IDENTIFIED BY", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 123", false},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY PASSWORD 'notarealhash'", false},
	{"CREATE USER u1@'%' IDENTIFIED BY PASSWORD PASSWORD 'notarealhash'", false},
	// A user list carries one auth option per user, not one for the list.
	{"CREATE USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"ALTER USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},

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
	{"CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user, URI:spiffe://example.org/myservice'", true},
	{"CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'", true},
	{"CREATE USER u1@'%' REQUIRE token_issuer 'issuer-abc'", true},
	{"CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA' SAN 'DNS:example-user'", true},
	{"CREATE USER u1@'%' REQUIRE ISSUER 'i' AND SAN 'DNS:example-user'", true},
	{"CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'", true},
	{"CREATE USER u1@'%' REQUIRE ISSUER 'i' SUBJECT 's' CIPHER 'c' SAN 'DNS:x' TOKEN_ISSUER 'ti'", true},
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
	{"CREATE USER u1@'%' PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK", true},
	{"CREATE USER u1@'%' ATTRIBUTE '{\"email\": \"user@example.com\"}'", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{\"email\": \"user@example.com\"}'", true},

	// --- GRANT and REVOKE: the grantee list ---
	// The privilege forms take a UserSpecList, so each grantee may carry an
	// auth option of its own.
	{"GRANT ALL PRIVILEGES ON zz.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'", true},
	{"GRANT PROCESS, USAGE ON *.* TO 'u1'@'%' IDENTIFIED BY PASSWORD 'notarealhash'", true},
	{"REVOKE ALL PRIVILEGES ON zz.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'", true},
	{"REVOKE SELECT ON *.* FROM u1 IDENTIFIED BY 'notarealpw'", true},
	{"GRANT ALL ON *.* TO u1 IDENTIFIED BY 'notarealpw' REQUIRE SSL WITH GRANT OPTION", true},
	// A GRANT carries a REQUIRE clause and a resource option here;
	// mysql:8.0 and oceanbase-ce:4.3.5 have neither.
	{"GRANT ALL ON *.* TO u1 REQUIRE SSL", true},
	{"GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 10", true},
	{"GRANT ALL ON *.* TO u1 REQUIRE SAN 'DNS:example-user'", true},
	// The role forms take a UsernameList instead, and so does GRANT PROXY.
	{"GRANT 'role1' TO u1 IDENTIFIED BY 'notarealpw'", false},
	{"REVOKE 'role1' FROM u1 IDENTIFIED BY 'notarealpw'", false},
	{"GRANT PROXY ON u1 TO u2 IDENTIFIED BY 'notarealpw'", false},

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
	{"GRANT 'role1' ON *.* TO u1", false},
	{"GRANT 'role1'@'%' ON *.* TO u1", false},
	{"REVOKE 'role1' ON *.* FROM u1", false},
	{"GRANT SELECT, 'role1' ON *.* TO u1", false},
	// A word with an @host after it is a role too, keyword or not, and the
	// privilege reading is closed to it as well.
	{"GRANT role1@localhost TO u1", true},
	{"GRANT role1@`localhost` TO u1", true},
	{"GRANT role1@'localhost' TO u1", true},
	{"GRANT comment@localhost TO u1", true},
	{"REVOKE role1@localhost FROM u1", true},
	{"REVOKE comment@localhost FROM u1", true},
	{"GRANT role1@localhost ON *.* TO u1", false},
	{"GRANT comment@localhost ON *.* TO u1", false},
	// A word is a role name only when it is an identifier. A keyword
	// reduces to a privilege and the role reading is closed to it, which is
	// also why the two-word privileges are no role names.
	{"GRANT role1 TO u1", true},
	{"GRANT `select` TO u1", true},
	{"GRANT SELECT TO u1", false},
	{"GRANT COMMENT TO u1", false},
	// S3 is a keyword in TiDB and so no ExtendedPriv identifier, and an unreserved
	// keyword is neither a PrivType nor an identifier.
	{"GRANT S3 ON *.* TO u1", false},
	{"GRANT COMMENT ON *.* TO u1", false},
	{"GRANT COMMENT, SELECT ON *.* TO u1", false},
	{"GRANT usage TO u1", false},
	{"GRANT process TO u1", false},
	{"GRANT CREATE VIEW TO u1", false},
	{"GRANT REPLICATION SLAVE TO u1", false},
	{"GRANT GRANT OPTION TO u1", false},
	// SEND_CREDENTIALS_TO_TIKV is a BRIEBooleanOptionName keyword, so no
	// ExtendedPriv identifier and no PrivType; it is neither a grantee nor a
	// privilege. A backquoted spelling is a plain identifier and so a role.
	{"GRANT SEND_CREDENTIALS_TO_TIKV TO CURRENT_USER", false},
	{"GRANT SEND_CREDENTIALS_TO_TIKV TO u1", false},
	{"GRANT SEND_CREDENTIALS_TO_TIKV ON *.* TO u1", false},
	{"GRANT `send_credentials_to_tikv` TO u1", true},

	// --- the privilege names ---
	// GRANT OPTION is the privilege; GRANT on its own is not one.
	{"GRANT GRANT OPTION ON *.* TO u1", true},
	{"GRANT GRANT OPTION ON `db`.* TO u1", true},
	{"GRANT GRANT OPTION, SELECT ON *.* TO u1", true},
	{"GRANT SELECT, GRANT OPTION ON *.* TO u1", true},
	{"GRANT GRANT ON *.* TO u1", false},
	// ExtendedPriv is a run of identifiers, which is how the long names are
	// spelled.
	{"GRANT INVOKE LAMBDA ON *.* TO u1", true},
	{"GRANT INVOKE SAGEMAKER ON *.* TO u1", true},
	{"GRANT INVOKE COMPREHEND ON *.* TO u1", true},
	{"GRANT INVOKE LAMBDA, INVOKE SAGEMAKER ON *.* TO u1", true},
	{"GRANT BACKUP_ADMIN ON *.* TO u1", true},
	{"GRANT SELECT, BACKUP_ADMIN ON *.* TO u1", true},
	// The run takes identifiers, backquoted or not, and nothing else.
	{"GRANT INVOKE `LAMBDA` ON *.* TO u1", true},
	// LOAD FROM S3 and SELECT INTO S3 are spelled out in the grammar
	// because their first two words are reserved, so neither reaches the
	// run of identifiers. S3 is the word the rule names: a backquoted one
	// is an identifier, and the two words do not recombine.
	{"GRANT LOAD FROM S3 ON *.* TO u1", true},
	{"GRANT SELECT INTO S3 ON *.* TO u1", true},
	{"GRANT LOAD FROM s3 ON *.* TO u1", true},
	{"GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1", true},
	{"GRANT LOAD FROM `S3` ON *.* TO u1", false},
	{"GRANT LOAD S3 ON *.* TO u1", false},
	{"GRANT SELECT S3 ON *.* TO u1", false},
	{"GRANT LOAD FROM SELECT ON *.* TO u1", false},
	// An ExtendedPriv is a role name too, since the reading is still open
	// when the TO arrives.
	{"GRANT LOAD FROM S3 TO u1", true},
	{"GRANT INVOKE LAMBDA TO u1", true},
	{"GRANT BACKUP_ADMIN TO u1", true},
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
	// CURRENT_USER names the connected user wherever a user is named, its own
	// alternative of Username (parser.y:11346), with an optional empty "()".
	// TiDB v8.5.5 reads it in every account statement, SET DEFAULT ROLE and SET
	// PASSWORD among them; the parser.y corpus pins the full set.
	{"CREATE USER CURRENT_USER", true},
	{"CREATE USER CURRENT_USER()", true},
	{"ALTER USER CURRENT_USER IDENTIFIED BY 'notarealpw'", true},
	{"DROP USER CURRENT_USER", true},
	{"RENAME USER CURRENT_USER TO CURRENT_USER", true},
	{"GRANT 'role1' TO CURRENT_USER", true},
	{"REVOKE 'role1' FROM CURRENT_USER", true},
	{"GRANT PROXY ON CURRENT_USER TO CURRENT_USER", true},
	{"GRANT SELECT ON *.* TO CURRENT_USER", true},
	{"SET PASSWORD FOR CURRENT_USER = 'notarealpw'", true},
	{"SET DEFAULT ROLE NONE TO CURRENT_USER", true},
	// CURRENT_USER is a Username (parser.y:11346) and not a Rolename, so v8.5.5
	// answers 1064 for the keyword in the DEFAULT ROLE clause. The role-name
	// positions of CREATE/DROP/SET ROLE and SHOW GRANTS ... USING are pinned by
	// the role-name sweep.
	{"CREATE USER u1@localhost DEFAULT ROLE current_user", false},
	{"ALTER USER u1@localhost DEFAULT ROLE current_user", false},

	// ALTER USER USER() is the current-user function spec, `"USER" '(' ')'`
	// (parser.y:13745): both parens are required, so the close paren cannot be
	// dropped. The complete form stays valid.
	{"ALTER USER USER() IDENTIFIED BY 'notarealpw'", true},
	{"ALTER USER USER(IDENTIFIED BY 'notarealpw'", false},

	// --- IF [NOT] EXISTS and RENAME USER's TO are required where they appear ---
	// CREATE USER takes IfNotExists (`/* empty */ | "IF" NotSym "EXISTS"`),
	// DROP/ALTER USER take IfExists (`/* empty */ | "IF" "EXISTS"`), and
	// UserToUserList's UserToUser is `Username "TO" Username`, so a half-written
	// IF run or a pair missing TO is 1064 (parser.y IfNotExists / IfExists /
	// RenameUserStmt / UserToUserList). The complete forms stay valid.
	{"CREATE USER IF NOT EXISTS u", true},
	{"ALTER USER IF EXISTS u IDENTIFIED BY 'notarealpw'", true},
	{"DROP USER IF EXISTS u", true},
	{"RENAME USER a TO b, c TO d", true},
	{"CREATE USER IF EXISTS u", false},
	{"CREATE USER IF NOT u", false},
	{"ALTER USER IF u", false},
	{"DROP USER IF u", false},
	{"RENAME USER a b", false},
	{"RENAME USER a b, c TO d", false},
	{"RENAME USER a TO b, c d", false},
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

// grantWordsAsNames is the other half of the decision not to register SAN and
// TOKEN_ISSUER as keywords. Both are read as bare words, so both are still
// ordinary names everywhere else -- including the implicit-alias position,
// which is where a registered keyword hurts most: that branch takes only a
// plain identifier token, so registering a word there once broke `DELETE FROM
// t config WHERE config.id = 1`.
var grantWordsAsNames = []string{
	// As a table name, and as an implicit alias of one.
	"SELECT * FROM san",
	"SELECT * FROM token_issuer",
	"SELECT * FROM t san WHERE san.id = 1",
	"SELECT * FROM t token_issuer WHERE token_issuer.id = 1",
	"DELETE FROM t san WHERE san.id = 1",
	"UPDATE t token_issuer SET token_issuer.a = 1",
	// As a column name and a select-list alias.
	"SELECT san, token_issuer FROM t",
	"SELECT 1 AS san",
	"SELECT 1 AS token_issuer",
	// As names in DDL.
	"CREATE TABLE san (token_issuer int)",
	"CREATE TABLE t (san varchar(64), token_issuer varchar(64), INDEX san (san))",
	"ALTER TABLE t ADD COLUMN san varchar(64)",
	"DROP TABLE token_issuer",
	// And beside the statements of this family, which is what makes the
	// bare-word reading load-bearing rather than incidental.
	"CREATE USER san@'%' REQUIRE SAN 'DNS:example-user'",
	"CREATE USER token_issuer@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'",
	"GRANT SELECT ON san.* TO token_issuer",
}

// TestGrantWordsStayNames holds the other half of the keyword decision: every
// word this family reads bare is still an ordinary name everywhere else. It
// needs no container.
func TestGrantWordsStayNames(t *testing.T) {
	for _, sql := range grantWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}

// grantKnownDivergences records the statements of this family this parser
// still answers differently from pingcap/tidb:v8.5.5, with what this parser
// answers, so each is a measured and named gap rather than an untested one.
// Every one of them needs knowledge this parser does not hold -- which user
// options that server has -- rather than a clause of the family this change
// closed, so each belongs to its own change.
var grantKnownDivergences = []grantFormCase{
	// RANDOM PASSWORD is MySQL 8.0's. Retiring it moves this package's
	// own tests, so it belongs to its own change.
	{"CREATE USER u1@'%' IDENTIFIED BY RANDOM PASSWORD", true},
	{"CREATE USER u1@'%' IDENTIFIED WITH caching_sha2_password BY RANDOM PASSWORD", true},
	{"ALTER USER u1@'%' IDENTIFIED BY RANDOM PASSWORD RETAIN CURRENT PASSWORD", true},
	// ALL ends the privilege list here, where TiDB reads it as one element
	// of the list -- a coverage gap rather than a false accept.
	{"GRANT ALL, GRANT OPTION ON *.* TO u1", false},
	{"GRANT ALL PRIVILEGES, GRANT OPTION ON *.* TO u1", false},
}

// TestGrantKnownDivergences holds the list above to this parser, so a case
// that has started agreeing with pingcap/tidb:v8.5.5 fails here and is moved
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
