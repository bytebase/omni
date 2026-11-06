package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// commentValueCase is a statement, the code MariaDB's PREPARE answers for it
// (1064 when it does not parse, anything else when it does), and for a rejection
// the text the parser's error must point at; an empty near means the end of the
// input.
type commentValueCase struct {
	sql  string
	code uint16
	near string
}

// commentValueCases are the COMMENT clauses outside the table options, and the
// index options of a table's key, written without the value they take. Every
// COMMENT takes a quoted string and nothing else, on a column, an index, a
// routine and an event. MariaDB 11.8.8 answers 1064 for each one without it, and
// parses each one with it. ta is a table the oracle seeds.
//
// MariaDB has no user COMMENT or ATTRIBUTE, no CREATE TABLESPACE and no CREATE
// LOGFILE GROUP, and answers 1064 for them with or without the string. This
// parser reads all four with the string, which is a separate gap, so only the
// forms without it are here: both sides refuse those.
var commentValueCases = []commentValueCase{
	{"CREATE TABLE t (a INT COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a INT COMMENT, b INT)", 1064, ", b INT)"},
	{"CREATE TABLE t (a INT COMMENT NOT NULL)", 1064, "NOT NULL"},
	{"CREATE TABLE t (a INT COMMENT DEFAULT 1)", 1064, "DEFAULT"},
	{"CREATE TABLE t (a INT COMMENT = 'x')", 1064, "= 'x'"},
	{"CREATE TABLE t (a INT AS (1) COMMENT)", 1064, ")"},
	{"ALTER TABLE ta ADD COLUMN b INT COMMENT", 1064, ""},
	{"ALTER TABLE ta MODIFY a INT COMMENT", 1064, ""},
	{"ALTER TABLE ta CHANGE a b INT COMMENT", 1064, ""},

	// A key's index options: CREATE INDEX always refused a bare COMMENT, but a
	// key in CREATE TABLE or ALTER TABLE ... ADD dropped the error and went on.
	{"CREATE INDEX i ON ta (a) COMMENT", 1064, ""},
	{"CREATE TABLE t (a INT, KEY i (a) COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a INT, KEY i (a) COMMENT, b INT)", 1064, ", b INT)"},
	{"CREATE TABLE t (a INT, UNIQUE KEY i (a) COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) COMMENT)", 1064, ")"},
	{"CREATE TABLE t (a INT, KEY i (a) KEY_BLOCK_SIZE =)", 1064, ")"},
	{"CREATE TABLE t (a INT, KEY i (a) USING)", 1064, ")"},
	{"CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER)", 1064, ")"},
	{"CREATE TABLE t (a INT, KEY i (a) ENGINE_ATTRIBUTE)", 1064, ")"},
	{"CREATE TABLE t (a INT, KEY i (a) SECONDARY_ENGINE_ATTRIBUTE =)", 1064, ")"},
	{"ALTER TABLE ta ADD INDEX i (a) COMMENT", 1064, ""},
	{"ALTER TABLE ta ADD UNIQUE (a) COMMENT", 1064, ""},
	{"ALTER TABLE ta ADD PRIMARY KEY (a) COMMENT", 1064, ""},
	{"ALTER TABLE ta ADD INDEX i (a) COMMENT, ADD COLUMN b INT", 1064, ", ADD COLUMN"},

	{"CREATE PROCEDURE p() COMMENT SELECT 1", 1064, "SELECT 1"},
	{"CREATE PROCEDURE p() COMMENT BEGIN END", 1064, "BEGIN END"},
	{"CREATE FUNCTION f() RETURNS INT COMMENT RETURN 1", 1064, "RETURN 1"},
	{"CREATE FUNCTION f() RETURNS INT NO SQL COMMENT RETURN 1", 1064, "RETURN 1"},
	{"ALTER PROCEDURE p COMMENT", 1064, ""},
	{"ALTER FUNCTION f COMMENT", 1064, ""},
	{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY COMMENT DO SELECT 1", 1064, "DO SELECT 1"},
	{"ALTER EVENT e COMMENT", 1064, ""},
	{"ALTER EVENT e COMMENT DO SELECT 1", 1064, "DO SELECT 1"},
	{"ALTER EVENT e ENABLE COMMENT", 1064, ""},
	{"CREATE USER u COMMENT", 1064, ""},
	{"ALTER USER u COMMENT", 1064, ""},
	{"CREATE USER u ATTRIBUTE", 1064, ""},
	{"ALTER USER u ATTRIBUTE", 1064, ""},

	{"CREATE TABLE t (a INT COMMENT 'x', b INT)", 0, ""},
	{"CREATE TABLE t (a INT NOT NULL COMMENT \"x\" DEFAULT 1)", 0, ""},
	{"CREATE TABLE t (a INT AS (1) COMMENT 'x')", 0, ""},
	{"ALTER TABLE ta MODIFY a INT COMMENT 'x'", 0, ""},
	{"ALTER TABLE ta ADD COLUMN b INT COMMENT 'x' FIRST", 0, ""},
	{"CREATE INDEX i ON ta (a) COMMENT 'x'", 0, ""},
	{"CREATE TABLE t (a INT, KEY i (a) COMMENT 'x' USING BTREE KEY_BLOCK_SIZE = 8)", 0, ""},
	{"CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT 'x', UNIQUE (a) COMMENT 'y')", 0, ""},
	// 1128: the ngram parser plugin is not installed, which MariaDB learns after
	// parsing the statement.
	{"CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER ngram)", 1128, ""},
	{"ALTER TABLE ta ADD INDEX i (a) COMMENT 'x', ADD COLUMN b INT", 0, ""},
	{"CREATE PROCEDURE p() COMMENT 'x' SELECT 1", 0, ""},
	{"CREATE FUNCTION f() RETURNS INT COMMENT 'x' RETURN 1", 0, ""},
	{"ALTER PROCEDURE p COMMENT 'x'", 0, ""},
	{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY COMMENT 'x' DO SELECT 1", 0, ""},
	{"ALTER EVENT e COMMENT 'x' DO SELECT 1", 0, ""},
}

// TestCommentValue holds the parser to commentValueCases. A rejection is a
// *ParseError at the token where the value should have been.
func TestCommentValue(t *testing.T) {
	for _, tc := range commentValueCases {
		_, err := Parse(tc.sql)
		if (err == nil) != (tc.code != 1064) {
			t.Errorf("Parse(%q) error = %v, want rejected %v", tc.sql, err, tc.code == 1064)
			continue
		}
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", tc.sql, err)
			continue
		}
		want := len(tc.sql)
		if tc.near != "" {
			want = strings.LastIndex(tc.sql, tc.near)
		}
		if pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", tc.sql, pe.Position, want, err)
		}
	}
}

// TestCommentValueOracle asserts omni agrees with a live MariaDB 11.8.8 on
// commentValueCases (container test). MariaDB's verdict comes from PREPARE,
// which parses without running the statement.
func TestCommentValueOracle(t *testing.T) {
	o := startMariaDB(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE ta (a INT)"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range commentValueCases {
		var code uint16
		_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(tc.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("PREPARE %q: the oracle gave no verdict: %v", tc.sql, err)
			}
			code = myErr.Number
		}
		if code != tc.code {
			t.Errorf("MariaDB code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != (code != 1064) {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
