package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// primaryKeyCase is a statement, whether OceanBase 4.3.5 parses it (false means the
// engine answers 1064), and for a refused one the text the error is at: its last
// occurrence in the statement.
type primaryKeyCase struct {
	sql    string
	parses bool
	at     string
}

// primaryKeyGrammarCases: OceanBase reads a table-level PRIMARY KEY with a grammar
// narrower than MySQL's, and a different one in CREATE TABLE and in ALTER TABLE ...
// ADD (parsePrimaryKey). Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02.
// primaryKeyScopeCases holds the LOCAL, GLOBAL and BLOCK_SIZE side.
var primaryKeyGrammarCases = []primaryKeyCase{
	// CREATE TABLE takes one USING before the columns, one after, then one COMMENT.
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) KEY_BLOCK_SIZE 8)", false, "KEY_BLOCK_SIZE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) KEY_BLOCK_SIZE = 8)", false, "KEY_BLOCK_SIZE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE)", false, "VISIBLE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) INVISIBLE)", false, "INVISIBLE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) WITH PARSER ngram)", false, "WITH"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT 'x' USING BTREE)", false, "USING"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) USING BTREE USING HASH)", false, "USING"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT 'x' COMMENT 'y')", false, "COMMENT"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT = 'x')", false, "= 'x'"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) ENGINE_ATTRIBUTE='{}')", false, "ENGINE_ATTRIBUTE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) SECONDARY_ENGINE_ATTRIBUTE='{}')", false, "SECONDARY_ENGINE_ATTRIBUTE"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE COMMENT 'x')", false, "VISIBLE"},
	{"CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) INVISIBLE)", false, "INVISIBLE"},
	{"CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) KEY_BLOCK_SIZE 8)", false, "KEY_BLOCK_SIZE"},
	{"CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) COMMENT 'x' USING HASH)", false, "USING"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id) USING BTREE COMMENT 'x')", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY USING BTREE (id) USING HASH)", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY USING BTREE (id) USING HASH COMMENT 'x')", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY USING HASH (id) COMMENT 'x')", true, ""},
	{"CREATE TABLE t (id INT, CONSTRAINT `c` PRIMARY KEY (id) USING BTREE COMMENT 'x')", true, ""},
	// An index name, which MySQL and OceanBase ignore, only without CONSTRAINT, and
	// KEY is not optional.
	{"CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY pk (id))", false, "pk (id)"},
	{"CREATE TABLE t (id INT, PRIMARY (id))", false, "(id)"},
	{"CREATE TABLE t (id INT, PRIMARY KEY pk (id))", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY `pk` (id))", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY pk USING BTREE (id))", true, ""},
	{"CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY USING BTREE (id))", true, ""},
	{"CREATE TABLE t (id INT, CONSTRAINT PRIMARY KEY (id))", true, ""},
	// The columns are plain names: no prefix, no ASC or DESC, no expression.
	{"CREATE TABLE t (id INT, PRIMARY KEY (id(3)))", false, "(3)"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id DESC))", false, "DESC"},
	{"CREATE TABLE t (id INT, PRIMARY KEY (id ASC))", false, "ASC"},
	{"CREATE TABLE t (id INT, PRIMARY KEY ((id + 1)))", false, "(id + 1)"},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id, a))", true, ""},
	{"CREATE TABLE t (id INT, PRIMARY KEY (`id`))", true, ""},
	// A column's own PRIMARY KEY is a different production.
	{"CREATE TABLE t (id INT PRIMARY KEY COMMENT 'x')", true, ""},

	// ALTER TABLE ... ADD takes no USING before the columns, and any index option
	// after them.
	{"ALTER TABLE obpk_t ADD PRIMARY KEY USING BTREE (id)", false, "USING"},
	{"ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY USING BTREE (id)", false, "USING"},
	{"ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY pk (id)", false, "pk (id)"},
	{"ALTER TABLE obpk_t ADD PRIMARY (id)", false, "(id)"},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id(3))", false, "(3)"},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id DESC)", false, "DESC"},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id ASC)", false, "ASC"},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY ((id + 1))", false, "(id + 1)"},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY pk (id)", true, ""},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id) KEY_BLOCK_SIZE 8", true, ""},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id) VISIBLE", true, ""},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id) INVISIBLE", true, ""},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id) COMMENT 'x' USING BTREE", true, ""},
	{"ALTER TABLE obpk_t ADD PRIMARY KEY (id) USING BTREE USING HASH", true, ""},
	{"ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY (id) INVISIBLE", true, ""},
}

// TestOceanBasePrimaryKeyGrammar holds the parser to primaryKeyGrammarCases. A refused
// statement fails where OceanBase's "near" says. TestOceanBasePrimaryKeyGrammarOracle
// asks the engine the same statements.
func TestOceanBasePrimaryKeyGrammar(t *testing.T) {
	for _, c := range primaryKeyGrammarCases {
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
		if want := strings.LastIndex(c.sql, c.at); pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d at %q (%v)", c.sql, pe.Position, want, c.at, err)
		}
	}
}

// TestOceanBasePrimaryKeyGrammarOracle asks OceanBase about primaryKeyGrammarCases
// through PREPARE, which parses without creating anything: 1064 is a syntax error, any
// other answer means the statement parsed.
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBasePrimaryKeyGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obpk",
		"USE obpk",
		"DROP TABLE IF EXISTS obpk_t",
		"CREATE TABLE obpk_t (id INT NOT NULL, a INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obpk")
	})
	for _, c := range primaryKeyGrammarCases {
		_, err := db.ExecContext(ctx, "PREPARE obpk_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
