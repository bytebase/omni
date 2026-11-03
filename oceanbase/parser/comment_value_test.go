package parser

import "testing"

// TestCommentRequiresValue pins that a COMMENT outside the table options takes its
// quoted string, as the table option does since #583, on a column, an index, a
// routine, an event, a user, a tablespace and a logfile group, and that a key's index
// option inside CREATE TABLE or ALTER TABLE ... ADD fails when its value is missing.
// A missing string used to become an empty comment, and a key's failed option was
// dropped rather than reported, so review passed DDL the server answers 1064 for.
// TestOceanBaseGrammarInventory's truncatedGrammar holds the same forms to OceanBase.
func TestCommentRequiresValue(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT COMMENT)",
		"CREATE TABLE t (a INT COMMENT, b INT)",
		"CREATE TABLE t (a INT COMMENT NOT NULL)",
		"ALTER TABLE t ADD COLUMN b INT COMMENT",
		"ALTER TABLE t MODIFY a INT COMMENT",
		"CREATE TABLE t (a INT, KEY i (a) COMMENT)",
		"CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT, b INT)",
		"CREATE TABLE t (a INT, KEY i (a) USING)",
		"CREATE TABLE t (a INT, KEY i (a) BLOCK_SIZE)",
		"ALTER TABLE t ADD INDEX i (a) COMMENT",
		"ALTER TABLE t ADD INDEX i (a) COMMENT, ADD COLUMN b INT",
		"CREATE PROCEDURE p() COMMENT SELECT 1",
		"ALTER PROCEDURE p COMMENT",
		"ALTER EVENT e COMMENT",
		"CREATE USER u COMMENT",
		"ALTER USER u ATTRIBUTE",
		"CREATE TABLESPACE ts ADD DATAFILE 'ts.ibd' COMMENT",
		"CREATE LOGFILE GROUP lg ADD UNDOFILE 'u.dat' COMMENT =",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("parsed %q, but the option has no value", sql)
		}
	}
	// The same clauses with values must still parse, or the fix is too strict.
	for _, sql := range []string{
		"CREATE TABLE t (a INT COMMENT 'x', b INT NOT NULL COMMENT \"y\" DEFAULT 1)",
		"ALTER TABLE t MODIFY a INT COMMENT 'x'",
		"CREATE TABLE t (a INT, KEY i (a) COMMENT 'x' USING BTREE BLOCK_SIZE 16384, PRIMARY KEY (a) COMMENT 'y')",
		"ALTER TABLE t ADD INDEX i (a) COMMENT 'x', ADD COLUMN b INT",
		"CREATE PROCEDURE p() COMMENT 'x' SELECT 1",
		"ALTER PROCEDURE p COMMENT 'x'",
		"ALTER EVENT e COMMENT 'x'",
		"CREATE USER u COMMENT 'x'",
		"CREATE TABLESPACE ts ADD DATAFILE 'ts.ibd' COMMENT = 'x'",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%q: %v", sql, err)
		}
	}
}
