package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	gomysql "github.com/go-sql-driver/mysql"
)

// alterTableSetCase is a statement, whether OceanBase 4.3.5 parses it (false means the
// engine answers 1064), and for a refused one the text the error is at: its last
// occurrence in the statement, or "" for the end of the input.
type alterTableSetCase struct {
	sql    string
	parses bool
	at     string
}

// alterTableSetCases: an ALTER TABLE action may be `[SET] table_option [table_option
// ...]`, the options without commas between them (parseAlterTableSet). Read against
// 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02. obset_t is a plain table and obset_p one
// partitioned by RANGE COLUMNS over a DATETIME.
var alterTableSetCases = []alterTableSetCase{
	{"ALTER TABLE obset_t SET COMMENT 'x'", true, ""},
	{"ALTER TABLE obset_t SET COMMENT = 'y'", true, ""},
	{"ALTER TABLE obset_t SET PCTFREE = 5", true, ""},
	{"ALTER TABLE obset_t SET PCTFREE = 6 COMMENT 'z'", true, ""},
	{"ALTER TABLE obset_t SET PCTFREE = 7, COMMENT 'w'", true, ""},
	{"ALTER TABLE obset_t SET PCTFREE = 8, SET COMMENT 'v'", true, ""},
	{"ALTER TABLE obset_t SET AUTO_INCREMENT = 100", true, ""},
	{"ALTER TABLE obset_t SET ENGINE = InnoDB", true, ""},
	{"ALTER TABLE obset_t SET CHARSET = utf8mb4", true, ""},
	{"ALTER TABLE obset_t SET DEFAULT CHARSET utf8mb4", true, ""},
	{"ALTER TABLE obset_t SET CHARACTER SET utf8mb4", true, ""},
	{"ALTER TABLE obset_t SET DEFAULT CHARACTER SET = gbk COLLATE gbk_bin", true, ""},
	{"ALTER TABLE obset_t SET ROW_FORMAT = COMPRESSED", true, ""},
	{"ALTER TABLE obset_t SET COLLATE utf8mb4_bin", true, ""},
	{"ALTER TABLE obset_t SET BLOCK_SIZE = 32768 REPLICA_NUM = 1", true, ""},
	{"ALTER TABLE obset_t SET USE_BLOOM_FILTER = TRUE", true, ""},
	{"ALTER TABLE obset_t SET KEY_BLOCK_SIZE = 8", true, ""},
	{"ALTER TABLE obset_t SET AUTO_INCREMENT_MODE = 'ORDER'", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x' ENGINE = InnoDB ROW_FORMAT = DYNAMIC", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x', ADD COLUMN b INT", true, ""},
	{"ALTER TABLE obset_t ADD COLUMN b INT, SET COMMENT 'x'", true, ""},
	// OceanBase parses this and refuses it with 1235.
	{"ALTER TABLE obset_t SET ORGANIZATION HEAP", true, ""},
	{"ALTER TABLE obset_p SET DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)", true, ""},
	{"ALTER TABLE obset_p SET DYNAMIC_PARTITION_POLICY = (ENABLE = FALSE)", true, ""},
	{"ALTER TABLE obset_p DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x' COMMENT 'y'", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x', ENGINE = InnoDB, SET ROW_FORMAT = DYNAMIC", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x', DROP COLUMN a", true, ""},
	{"ALTER TABLE obset_t SET COMMENT 'x', ALGORITHM = INPLACE", true, ""},
	{"ALTER TABLE obset_t ALGORITHM = INPLACE, SET COMMENT 'x'", true, ""},
	{"ALTER TABLE obset_t SET DEFAULT COLLATE = utf8mb4_bin", true, ""},
	{"ALTER TABLE obset_t SET STATS_PERSISTENT = 0", true, ""},
	{"ALTER TABLE obset_t SET MAX_ROWS = 10", true, ""},
	{"ALTER TABLE obset_t SET CHECKSUM = 1", true, ""},
	{"ALTER TABLE obset_t SET COMPRESSION = 'zstd_1.3.8'", true, ""},
	{"ALTER TABLE obset_t SET TABLET_SIZE = 134217728", true, ""},
	{"ALTER TABLE obset_t SET AUTO_INCREMENT 5", true, ""},
	{"ALTER TABLE obset_t set comment 'x'", true, ""},
	// OceanBase parses this and refuses it with 1235.
	{"ALTER TABLE obset_t SET ORGANIZATION = INDEX", true, ""},
	// Without SET the options follow one another the same way.
	{"ALTER TABLE obset_t PCTFREE = 9 COMMENT 'u'", true, ""},
	{"ALTER TABLE obset_t ENGINE = InnoDB COMMENT 'x'", true, ""},
	{"ALTER TABLE obset_t COMMENT 'x' PCTFREE 5", true, ""},
	// SET takes table options and nothing else.
	{"ALTER TABLE obset_t SET", false, ""},
	{"ALTER TABLE obset_t SET ALGORITHM = INPLACE", false, "ALGORITHM"},
	{"ALTER TABLE obset_t SET LOCK = NONE", false, "LOCK"},
	{"ALTER TABLE obset_t SET ADD COLUMN b INT", false, "ADD"},
	{"ALTER TABLE obset_t SET SET COMMENT 'x'", false, "SET"},
	{"ALTER TABLE obset_t SET COMMENT 'x' ADD COLUMN b INT", false, "ADD"},
	{"ALTER TABLE obset_t SET, COMMENT 'x'", false, ", COMMENT"},
	{"ALTER TABLE obset_t SET CONVERT TO CHARACTER SET utf8mb4", false, "CONVERT"},
	{"ALTER TABLE obset_t SET PARTITION BY HASH(a) PARTITIONS 2", false, "PARTITION BY"},
	{"ALTER TABLE obset_t `SET` COMMENT 'x'", false, "`SET`"},
}

// TestOceanBaseAlterTableSetGrammar holds the parser to alterTableSetCases. A refused
// statement fails where OceanBase's "near" says. TestOceanBaseAlterTableSetGrammarOracle
// asks the engine the same statements.
func TestOceanBaseAlterTableSetGrammar(t *testing.T) {
	for _, c := range alterTableSetCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.parses {
			t.Errorf("Parse(%q) error = %v, want parses %v", c.sql, err, c.parses)
			continue
		}
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", c.sql, err)
			continue
		}
		want := len(c.sql)
		if c.at != "" {
			want = strings.LastIndex(c.sql, c.at)
		}
		if pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", c.sql, pe.Position, want, err)
		}
	}
}

// TestOceanBaseAlterTableSetCommands pins the commands an action with SET and several
// options gives: one ATTableOption per option, in order, the first spanning the SET.
func TestOceanBaseAlterTableSetCommands(t *testing.T) {
	const sql = "ALTER TABLE t SET PCTFREE = 6 COMMENT 'z', SET ROW_FORMAT = DYNAMIC, ENGINE = InnoDB KEY_BLOCK_SIZE = 8"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	stmt := list.Items[0].(*nodes.AlterTableStmt)
	want := []struct{ name, text string }{
		{"PCTFREE", "SET PCTFREE = 6"},
		{"COMMENT", "COMMENT 'z'"},
		{"ROW_FORMAT", "SET ROW_FORMAT = DYNAMIC"},
		{"ENGINE", "ENGINE = InnoDB"},
		{"KEY_BLOCK_SIZE", "KEY_BLOCK_SIZE = 8"},
	}
	if len(stmt.Commands) != len(want) {
		t.Fatalf("%d commands, want %d", len(stmt.Commands), len(want))
	}
	for i, w := range want {
		cmd := stmt.Commands[i]
		if cmd.Type != nodes.ATTableOption || cmd.Option == nil || !strings.EqualFold(cmd.Option.Name, w.name) {
			t.Errorf("command %d = %+v, want the %s option", i, cmd, w.name)
			continue
		}
		if got := sql[cmd.Loc.Start:cmd.Loc.End]; got != w.text {
			t.Errorf("command %d spans %q, want %q", i, got, w.text)
		}
	}
}

// TestOceanBaseAlterTableSetGrammarOracle asks OceanBase about alterTableSetCases
// through PREPARE, which parses without changing anything: 1064 is a syntax error, any
// other answer means the statement parsed.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseAlterTableSetGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obset",
		"USE obset",
		"DROP TABLE IF EXISTS obset_t",
		"DROP TABLE IF EXISTS obset_p",
		"CREATE TABLE obset_t (a INT PRIMARY KEY)",
		"CREATE TABLE obset_p (a DATETIME NOT NULL, PRIMARY KEY (a)) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY') PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('2026-01-01'))",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obset")
	})
	for _, c := range alterTableSetCases {
		_, err := db.ExecContext(ctx, "PREPARE obset_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		engineParses := err == nil
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			engineParses = myErr.Number != 1064
		}
		if engineParses != c.parses {
			t.Errorf("OceanBase parses=%v, want %v for %q: %v", engineParses, c.parses, c.sql, err)
		}
		if _, perr := Parse(c.sql); (perr == nil) != engineParses {
			t.Errorf("DIVERGENCE omni parses=%v, OceanBase parses=%v for %q: %v", perr == nil, engineParses, c.sql, perr)
		}
	}
}
