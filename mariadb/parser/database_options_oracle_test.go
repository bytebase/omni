package parser

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/compat"
)

// The option list CREATE DATABASE and ALTER DATABASE read, and the names its
// CHARACTER SET, CHARSET and COLLATE take. The name is MySQL's ident_or_text
// plus the reserved word BINARY -- an identifier, a non-reserved keyword, a
// backquoted identifier or a string literal in either quote -- so an integer
// and every other reserved keyword are a 1064.
//
// The rejections are the point of this file. The corpus in tidb/parser pins the
// forms pingcap's own parser test writes, and it writes none of these: a bare
// COLLATE, a conflicting charset, a comma between options, an option list left
// empty, a DEFAULT with nothing after it, and the options each engine does not
// have. Every case here was measured against MariaDB 11.8.8, and the expectation
// is what that server answered, so an image upgrade that changes a verdict
// fails here instead of passing silently.

// databaseOptionCase is a statement and whether the engine's parser accepts it.
type databaseOptionCase struct {
	sql    string
	accept bool
}

var databaseOptionCases = []databaseOptionCase{
	// The charset and collation name in every shape the engines read.
	{"ALTER DATABASE d CHARACTER SET = 'utf8'", true},
	{"ALTER DATABASE d CHARACTER SET 'utf8'", true},
	{`ALTER DATABASE d CHARACTER SET = "utf8"`, true},
	{"ALTER DATABASE d CHARACTER SET = `utf8`", true},
	{"ALTER DATABASE d CHARACTER SET utf8", true},
	{"ALTER DATABASE d DEFAULT CHARACTER SET = 'utf8'", true},
	{"ALTER SCHEMA d DEFAULT CHARSET = 'UTF8'", true},
	{"ALTER DATABASE d CHARSET = binary", true},
	{"ALTER DATABASE d CHARSET binary", true},
	{"ALTER DATABASE d CHARACTER SET binary", true},
	{"ALTER DATABASE d COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE d COLLATE 'utf8_bin'", true},
	{"ALTER DATABASE d COLLATE = `utf8_bin`", true},
	{"ALTER DATABASE d DEFAULT COLLATE = 'utf8_bin'", true},
	// A name the server does not know is a 1115 when the statement runs, not a
	// 1064 when it parses.
	{"ALTER DATABASE d CHARSET = 'nosuchcharset'", true},
	// Several options at once, with and without DEFAULT, and conflicting
	// charsets, which the server settles after parsing rather than refusing.
	{"ALTER DATABASE d CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE d CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'", false},
	{"ALTER DATABASE d DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'", false},
	// The nameless arm, which the reserved words COLLATE and DEFAULT reach.
	{"ALTER DATABASE COLLATE = 'utf8_bin'", true},
	{"ALTER SCHEMA COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", true},
	// An empty quoted name is a name: the server parses it and answers 1102.
	{"ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'", true},
	// CREATE DATABASE reads the same names and needs no option at all.
	{"CREATE DATABASE d", true},
	{"CREATE DATABASE d CHARACTER SET = 'utf8'", true},
	{"CREATE DATABASE d COLLATE = 'utf8_bin'", true},
	{"CREATE DATABASE d CHARSET = binary", true},
	{"CREATE SCHEMA d DEFAULT COLLATE = 'UTF8_BiN'", true},
	{"CREATE DATABASE d DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_bin'", true},

	// BINARY is a collation name as well as a charset name.
	{"ALTER DATABASE d COLLATE = binary", true},
	{"ALTER DATABASE d COLLATE binary", true},
	{"ALTER DATABASE d CHARSET=binary COLLATE = binary", true},

	// ENCRYPTION is not an option on this engine, on either statement.
	{"ALTER DATABASE d ENCRYPTION = 'N'", false},
	{"ALTER DATABASE d ENCRYPTION 'N'", false},
	{"ALTER DATABASE d DEFAULT ENCRYPTION = 'N'", false},
	{"ALTER DATABASE d ENCRYPTION", false},
	{"ALTER DATABASE ENCRYPTION = 'N'", false},
	{"CREATE DATABASE d ENCRYPTION = 'N'", false},
	{"CREATE DATABASE d DEFAULT ENCRYPTION = 'Y'", false},

	// READ ONLY is not an option on this engine, in any spelling.
	{"ALTER DATABASE d READ ONLY = 1", false},
	{"ALTER DATABASE d READ ONLY = 0", false},
	{"ALTER DATABASE d READ ONLY = DEFAULT", false},
	{"ALTER DATABASE d READ ONLY 1", false},
	{"ALTER DATABASE d READ ONLY", false},
	{"ALTER DATABASE READ ONLY = 1", false},
	{"ALTER DATABASE d READ", false},
	{"ALTER DATABASE d READ WRITE", false},
	{"ALTER DATABASE d DEFAULT READ ONLY = 1", false},
	{"ALTER DATABASE d CHARSET = utf8 READ ONLY", false},

	// The name is required wherever an option takes one.
	{"ALTER DATABASE d COLLATE", false},
	{"ALTER DATABASE d COLLATE =", false},
	{"ALTER DATABASE d CHARSET =", false},
	{"ALTER DATABASE d CHARACTER SET", false},
	{"ALTER DATABASE d CHARACTER SET =", false},
	{"ALTER DATABASE d CHARACTER SET = 'utf8' COLLATE", false},
	// An integer is no name, and neither is a reserved keyword other than
	// BINARY.
	{"ALTER DATABASE d CHARSET = 1", false},
	{"ALTER DATABASE d COLLATE = 1", false},
	{"ALTER DATABASE d CHARSET = select", false},
	{"ALTER DATABASE d COLLATE = select", false},
	{"CREATE DATABASE d CHARSET = 1", false},
	{"CREATE DATABASE d COLLATE = select", false},
	// The option list carries no commas, and no comma stands after the name.
	{"ALTER DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'", false},
	{"ALTER DATABASE d, CHARSET = 'utf8'", false},
	{"CREATE DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'", false},
	// ALTER DATABASE needs one option; CREATE DATABASE is the one that does not.
	{"ALTER DATABASE", false},
	{"ALTER DATABASE d", false},
	{"ALTER SCHEMA", false},
	{"ALTER SCHEMA d", false},
	// DEFAULT is a prefix, not an option, in front of one and after one.
	{"ALTER DATABASE d DEFAULT", false},
	{"ALTER DATABASE d DEFAULT = 'utf8'", false},
	{"ALTER DATABASE d CHARSET = utf8 DEFAULT", false},
	// CHARACTER starts an option only when SET follows it.
	{"ALTER DATABASE d CHARSET = utf8 CHARACTER", false},
	{"CREATE DATABASE d CHARSET = utf8 CHARACTER", false},
	// A reserved word is no database name.
	{"ALTER DATABASE binary CHARSET = 'utf8'", false},
	// READ ONLY is no CREATE DATABASE option on any of these engines.
	{"CREATE DATABASE d READ ONLY = 1", false},
	{"CREATE DATABASE d READ ONLY 0", false},
	{"CREATE DATABASE d READ ONLY = DEFAULT", false},

	// CHARSET is not reserved, so it reads as the database name and the
	// statement then runs out of options.
	{"ALTER DATABASE CHARSET = 'utf8'", false},
	{"ALTER DATABASE CHARSET utf8", false},

	// TiDB's placement and TiFlash options are not this engine's.
	{"ALTER DATABASE d PLACEMENT POLICY SET DEFAULT", false},
	{"ALTER DATABASE d PLACEMENT POLICY = 'p'", false},
	{"ALTER DATABASE PLACEMENT POLICY = 'p'", false},
	{"ALTER DATABASE d PLACEMENT", false},
	{"ALTER DATABASE d SET TIFLASH REPLICA 1", false},
	{"ALTER DATABASE SET TIFLASH REPLICA 1", false},
}

// prepareAccepts asks the oracle whether sql parses, without running it: the name the
// forms tables read prepareParses by.
func (o *parserOracle) prepareAccepts(t *testing.T, sql string) bool {
	t.Helper()
	return o.prepareParses(t, sql)
}

// prepareParses asks the oracle whether sql parses, without running it, and reads an
// error beyond 1064 by where it was raised. The statement is bound to a user variable, so
// its own quoting stays the server's business. MariaDB refuses some text while it parses with the code of the
// check that refused it, as in VECTOR(0)'s 1063, which prepareAccepts reads as parsed.
// The server reads all of the text before it looks at anything, so a statement refused
// after parsing answers 1064 once a syntax error follows it as a second statement, and
// a statement refused while parsing answers its own code again.
//
// The names it looks up while parsing count as read, as omni holds no such list: a
// system variable (1193), a character set (1115), a collation (1273, 1253), a plugin or
// function (1128), a storage engine (1286) and a structured variable (4081). So does a form a build leaves out (1289,
// SHOW PROCEDURE CODE without debugging).
func (o *parserOracle) prepareParses(t *testing.T, sql string) bool {
	t.Helper()
	parses, err := o.prepareParsesErr(sql)
	if err != nil {
		t.Fatalf("no verdict for %q: %v", sql, err)
	}
	return parses
}

// prepareParsesErr is prepareParses for a caller with no test to fail. The
// statement travels as PREPARE's argument, which the connection writes into it
// as a string literal, and a PREPARE of the same name replaces the last one:
// the question is one request, or two with the probe, which keeps a tape small.
func (o *parserOracle) prepareParsesErr(sql string) (bool, error) {
	prepare := func(text string) (uint16, error) {
		_, err := o.db.ExecContext(o.ctx, "PREPARE omni_parses_stmt FROM ?", text)
		if err == nil {
			return 0, nil
		}
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			return 0, err
		}
		return myErr.Number, nil
	}
	code, err := prepare(sql)
	if err != nil {
		return false, err
	}
	switch {
	case code == 0, mariaDBLookupErrors[code]:
		return true, nil
	case code == 1064:
		return false, nil
	}
	again, err := prepare(sql + compat.ProbeTail)
	if err != nil {
		return false, err
	}
	return again != code, nil
}

// TestPrepareParsesControls holds prepareParses to statements whose answer is known.
func TestPrepareParsesControls(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT 1 FROM nosuch", true},              // 1146
		{"SELECT nosuch", true},                     // 1054
		{"SELECT 1 FROM nosuch -- a comment", true}, // 1146, the tail survives the comment
		{"SELEKT 1", false},                         // 1064
		{"CREATE TABLE t (a VECTOR(0))", false},     // 1063, while parsing
		{"CREATE TABLE t (a VECTOR(16384))", false}, // 1074, while parsing
		{"SELECT CAST(1 AS POINT)", false},          // 4162, while parsing
		{"SELECT @@nosuch", true},                   // 1193, a name it looks up
	} {
		if got := o.prepareParses(t, c.sql); got != c.accept {
			t.Errorf("%q: prepareParses=%v, want %v", c.sql, got, c.accept)
		}
	}
}

// TestDatabaseOptionsOracle asserts omni agrees with a live
// MariaDB 11.8.8 on databaseOptionCases (container test).
func TestDatabaseOptionsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range databaseOptionCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
