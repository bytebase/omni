package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// optionValueCase is a statement, the error TiDB raises while parsing it (0
// when it parses), and for a rejection the text the parser's error must point
// at; an empty near means the end of the input.
type optionValueCase struct {
	sql  string
	code uint16
	near string
}

// optionValueCases are table, database and ALTER TABLE options written without
// the value they take. TiDB answers 1064 for every one: no option may stand
// alone, and COMMENT, CONNECTION and PASSWORD take a quoted string and nothing
// else. The TiDB options at the end of each half go through the same helper as
// the MySQL ones. The same options with a value parse. ta is a table the oracle
// seeds for the ALTER cases.
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
	{"CREATE TABLE t (a INT) SHARD_ROW_ID_BITS", 1064, ""},
	{"CREATE TABLE t (a INT) SHARD_ROW_ID_BITS =", 1064, ""},
	{"CREATE TABLE t (a INT) PRE_SPLIT_REGIONS =", 1064, ""},
	{"CREATE TABLE t (a INT) AUTO_ID_CACHE", 1064, ""},
	{"CREATE TABLE t (a INT) AUTO_RANDOM_BASE =", 1064, ""},
	{"CREATE TABLE t (a INT) TTL_ENABLE =", 1064, ""},
	{"CREATE TABLE t (a INT) PLACEMENT POLICY", 1064, ""},
	{"CREATE TABLE t (a INT) PLACEMENT POLICY =", 1064, ""},
	{"CREATE DATABASE d PLACEMENT POLICY", 1064, ""},
	{"CREATE DATABASE d PLACEMENT POLICY =", 1064, ""},

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
	{"CREATE TABLE t (a INT) SHARD_ROW_ID_BITS = 4", 0, ""},
	{"CREATE TABLE t (a INT) PLACEMENT POLICY = DEFAULT", 0, ""},
	{"CREATE DATABASE d PLACEMENT POLICY DEFAULT", 0, ""},
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

// TestTableOptionValueTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// optionValueCases (container test). TiDB's verdict comes from PREPARE, which
// parses without running the statement.
func TestTableOptionValueTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	// USE and PREPARE are session state: a dedicated connection keeps every
	// statement on the database it selects.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()
	for _, ddl := range []string{
		"DROP DATABASE IF EXISTS omni_option_value",
		"CREATE DATABASE omni_option_value",
		"USE omni_option_value",
		"CREATE TABLE ta (a INT)",
	} {
		if _, err := conn.ExecContext(tc.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_option_value")

	for _, c := range optionValueCases {
		var code uint16
		_, err := conn.ExecContext(tc.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("PREPARE %q: the oracle gave no verdict: %v", c.sql, err)
			}
			code = myErr.Number
		}
		if code != c.code {
			t.Errorf("TiDB code %d, want %d for %q: %v", code, c.code, c.sql, err)
			continue
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, TiDB code %d for %q: %v", omniAccepts, code, c.sql, perr)
		}
	}
}
