package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// optionValueCase is a statement, the error MariaDB raises while parsing it (0
// when it parses), and for a rejection the text the parser's error must point
// at; an empty near means the end of the input.
type optionValueCase struct {
	sql  string
	code uint16
	near string
}

// optionValueCases are table and ALTER TABLE options written without the
// value they take. MariaDB 11.8.8 answers 1064 for every one: the grammar has no
// option that may stand alone, and COMMENT, CONNECTION and PASSWORD take a
// quoted string and nothing else. The same options with a value parse. ta is
// a table the oracle seeds for the ALTER cases.
var optionValueCases = []optionValueCase{
	{"CREATE TABLE t (a INT) KEY_BLOCK_SIZE", 1064, ""},
	{"CREATE TABLE t (a INT) KEY_BLOCK_SIZE =", 1064, ""},
	{"CREATE TABLE t (a INT) ENGINE", 1064, ""},
	{"CREATE TABLE t (a INT) ENGINE =", 1064, ""},
	{"CREATE TABLE t (a INT) ROW_FORMAT", 1064, ""},
	{"CREATE TABLE t (a INT) COMPRESSION", 1064, ""},
	{"CREATE TABLE t (a INT) STATS_PERSISTENT =", 1064, ""},
	{"CREATE TABLE t (a INT) TABLESPACE", 1064, ""},
	{"CREATE TABLE t (a INT) DEFAULT CHARSET =", 1064, ""},
	{"CREATE TABLE t (a INT) COMMENT", 1064, ""},
	{"CREATE TABLE t (a INT) COMMENT =", 1064, ""},
	{"CREATE TABLE t (a INT) COMMENT ENGINE = InnoDB", 1064, "ENGINE"},
	{"CREATE TABLE t (a INT) CONNECTION", 1064, ""},
	{"CREATE TABLE t (a INT) CONNECTION =", 1064, ""},
	{"CREATE TABLE t (a INT) PASSWORD", 1064, ""},
	{"CREATE TABLE t (a INT) PASSWORD =", 1064, ""},
	{"ALTER TABLE ta ENGINE =", 1064, ""},
	{"ALTER TABLE ta KEY_BLOCK_SIZE", 1064, ""},
	{"ALTER TABLE ta COMMENT", 1064, ""},
	{"ALTER TABLE ta COMMENT =", 1064, ""},
	{"ALTER TABLE ta ALGORITHM =", 1064, ""},
	{"ALTER TABLE ta LOCK", 1064, ""},
	{"ALTER TABLE ta CONVERT TO CHARACTER SET", 1064, ""},
	{"ALTER TABLE ta CONVERT TO CHARACTER SET utf8mb4 COLLATE", 1064, ""},
	// CHARACTER is an option only before SET: alone, it used to be dropped.
	{"CREATE TABLE t (a INT) CHARACTER", 1064, ""},
	{"CREATE TABLE t (a INT) CHARACTER ENGINE = InnoDB", 1064, "ENGINE"},
	{"ALTER TABLE ta COMMENT 'x' CHARACTER", 1064, ""},
	// Another table option may follow one without a comma; any other action
	// may not.
	{"ALTER TABLE ta COMMENT 'x' DISABLE KEYS", 1064, "DISABLE"},
	{"ALTER TABLE ta COMMENT 'x' ADD COLUMN b3 INT", 1064, "ADD"},

	{"CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 8", 0, ""},
	{"CREATE TABLE t (a INT) ENGINE = InnoDB ROW_FORMAT = DYNAMIC", 0, ""},
	{"CREATE TABLE t (a INT) COMMENT 'x'", 0, ""},
	{"CREATE TABLE t (a INT) COMMENT = \"x\" ENGINE = InnoDB", 0, ""},
	{"CREATE TABLE t (a INT) CONNECTION 'x'", 0, ""},
	{"CREATE TABLE t (a INT) PASSWORD = 'x' COMMENT = 'y'", 0, ""},
	{"ALTER TABLE ta ENGINE = InnoDB", 0, ""},
	{"ALTER TABLE ta COMMENT 'x'", 0, ""},
	{"ALTER TABLE ta ALGORITHM = INPLACE, LOCK = NONE", 0, ""},
	{"ALTER TABLE ta CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin", 0, ""},
	{"CREATE TABLE t (a INT) CHARACTER SET utf8mb4", 0, ""},
	{"ALTER TABLE ta ENGINE = InnoDB COMMENT 'x'", 0, ""},
	{"ALTER TABLE ta ENGINE = InnoDB ROW_FORMAT = DYNAMIC", 0, ""},
	{"ALTER TABLE ta COMMENT 'x' ENGINE = InnoDB, ALGORITHM = INPLACE", 0, ""},
	{"ALTER TABLE ta ADD COLUMN b2 INT COMMENT 'c', ENGINE = InnoDB COMMENT 'z'", 0, ""},
	{"ALTER TABLE ta DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin", 0, ""},
	{"ALTER TABLE ta AUTO_INCREMENT = 5 COMMENT 'x' ENGINE = InnoDB", 0, ""},
	{"CREATE TABLE t (a INT) DEFAULT CHARACTER SET = utf8mb4 COLLATE utf8mb4_bin", 0, ""},
}

// TestTableOptionValue holds the parser to optionValueCases. A rejection is a
// *ParseError at the token where the value should have been.
func TestTableOptionValue(t *testing.T) {
	for _, tc := range optionValueCases {
		_, err := Parse(tc.sql)
		if (err == nil) != (tc.code == 0) {
			t.Errorf("Parse(%q) error = %v, want rejected %v", tc.sql, err, tc.code != 0)
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

// TestTableOptionValueOracle asserts omni agrees with a live MariaDB 11.8.8 on
// optionValueCases (container test). MariaDB's verdict comes from PREPARE,
// which parses without running the statement.
func TestTableOptionValueOracle(t *testing.T) {
	o := startMariaDB(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE ta (a INT)"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range optionValueCases {
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
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
