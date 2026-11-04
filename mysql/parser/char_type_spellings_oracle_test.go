package parser

import "testing"

// MySQL spells CHAR and VARCHAR several ways: CHARACTER is CHAR, VARYING after
// either makes it VARCHAR, VARCHARACTER is VARCHAR in one word, LONG takes the
// VARCHAR spellings as MEDIUMTEXT, and CHAR SET is CHARACTER SET in two words.
// https://dev.mysql.com/doc/refman/8.0/en/char.html
//
// The rule that separates the accepted spellings from the rejected ones is that
// every spelling of VARCHAR carries a length and every spelling of CHAR may
// leave it out. Rather than assert that here, this asks MySQL 8.0 and requires
// the parser to agree with it, so the rule is the engine's and not a memory of
// it.

func TestCharTypeSpellingsMatchMySQL(t *testing.T) {
	cases := []string{
		`CREATE TABLE t (a CHARACTER)`,
		`CREATE TABLE t (a CHARACTER(10))`,
		`CREATE TABLE t (a CHAR VARYING(50))`,
		`CREATE TABLE t (a CHARACTER VARYING(50))`,
		`CREATE TABLE t (a VARCHARACTER(50))`,
		`CREATE TABLE t (a CHARACTER VARYING(10) CHARACTER SET utf8)`,
		`CREATE TABLE t (a VARCHARACTER(50) CHARSET utf8)`,
		`CREATE TABLE t (a NATIONAL CHARACTER)`,
		`CREATE TABLE t (a NATIONAL CHARACTER VARYING(50))`,
		`CREATE TABLE t (a NATIONAL CHAR VARYING(50))`,
		`CREATE TABLE t (a NATIONAL VARCHARACTER(50))`,
		`CREATE TABLE t (a NCHAR VARCHAR(50))`,
		`CREATE TABLE t (a NCHAR VARYING(50))`,
		`CREATE TABLE t (a NCHAR VARCHARACTER(50))`,
		`CREATE TABLE t (a LONG VARCHARACTER)`,
		`CREATE TABLE t (a LONG CHAR VARYING)`,
		`CREATE TABLE t (a LONG CHARACTER VARYING)`,
		`CREATE TABLE t (a LONG CHARACTER VARYING CHARACTER SET utf8)`,
		`CREATE TABLE t (a LONG CHAR SET utf8mb4)`,
		`CREATE TABLE t (a CHAR CHAR SET utf8mb4)`,
		`ALTER TABLE t MODIFY a CHARACTER VARYING(10)`,
		// A VARCHAR spelling with no length.
		`CREATE TABLE t (a VARCHAR)`,
		`CREATE TABLE t (a VARCHARACTER)`,
		`CREATE TABLE t (a CHAR VARYING)`,
		`CREATE TABLE t (a CHARACTER VARYING)`,
		`CREATE TABLE t (a NVARCHAR)`,
		`CREATE TABLE t (a NATIONAL VARCHAR)`,
		`CREATE TABLE t (a NCHAR VARYING)`,
		`CREATE TABLE t (a VARBINARY)`,
		// Parentheses with no integer in them are not a length either, and a
		// type that must carry one says so rather than reading `()` as a
		// length it may leave out.
		`CREATE TABLE t (a VARCHAR())`,
		`CREATE TABLE t (a VARCHARACTER())`,
		`CREATE TABLE t (a CHARACTER VARYING())`,
		`CREATE TABLE t (a CHAR VARYING())`,
		`CREATE TABLE t (a NCHAR VARYING())`,
		`CREATE TABLE t (a NVARCHAR())`,
		`CREATE TABLE t (a VARBINARY())`,
		// LONG carries no length with any spelling.
		`CREATE TABLE t (a LONG CHAR VARYING(10))`,
		`CREATE TABLE t (a LONG VARCHARACTER(10))`,
		// CHARACTER and CHAR only start a charset clause when SET follows.
		`CREATE TABLE t (a LONG CHARACTER)`,
		`CREATE TABLE t (a TEXT CHARACTER)`,
		`CREATE TABLE t (a CHARACTER SET utf8)`,
		`CREATE TABLE t (a CHAR SET utf8)`,
	}
	o := startParserOracle(t)
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertOracleMatch(t, o, sql)
		})
	}
}
