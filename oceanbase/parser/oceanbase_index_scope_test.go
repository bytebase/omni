package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// indexScopeCase is a statement and whether OceanBase 4.3.5 parses it: false means
// the engine answers 1064.
type indexScopeCase struct {
	sql    string
	parses bool
}

// primaryKeyScopeCases: in CREATE TABLE a PRIMARY KEY takes USING and COMMENT and
// none of OceanBase's index scope or block size, while ALTER TABLE ... ADD PRIMARY KEY
// takes them. Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02.
var primaryKeyScopeCases = []indexScopeCase{
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) GLOBAL)", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) LOCAL)", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (a) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (a) LOCAL) PARTITION BY HASH(id) PARTITIONS 2", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id, a) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2", false},
	{"CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) GLOBAL)", false},
	{"CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) LOCAL)", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) COMMENT 'x' GLOBAL)", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) BLOCK_SIZE 16384)", false},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) USING BTREE)", true},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY USING BTREE (id))", true},
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) COMMENT 'x')", true},
	{"CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) COMMENT 'x')", true},
	// A primary key that leaves out the partitioning column parses, and is then
	// refused with 1503.
	{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (a)) PARTITION BY HASH(id) PARTITIONS 2", true},
	{"ALTER TABLE obscope_t ADD PRIMARY KEY (id) GLOBAL", true},
	{"ALTER TABLE obscope_t ADD PRIMARY KEY (id) LOCAL", true},
	{"ALTER TABLE obscope_t ADD PRIMARY KEY (id) BLOCK_SIZE 16384", true},
	{"ALTER TABLE obscope_t ADD CONSTRAINT pkc PRIMARY KEY (id) GLOBAL", true},
}

// indexScopeCases are the scope on the other index kinds, which OceanBase's grammar
// takes wherever it takes an index option: FULLTEXT and SPATIAL included, and repeated,
// in which case the last one stands. A GLOBAL fulltext index, and a GLOBAL spatial
// index on a partitioned table, parse and are then refused with 1235.
var indexScopeCases = []indexScopeCase{
	{"CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) LOCAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) GLOBAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ngram LOCAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) LOCAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) GLOBAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2", true},
	{"CREATE FULLTEXT INDEX fk ON obscope_t (c) LOCAL", true},
	{"CREATE FULLTEXT INDEX fk ON obscope_t (c) GLOBAL", true},
	{"CREATE SPATIAL INDEX sk ON obscope_t (g) GLOBAL", true},
	{"ALTER TABLE obscope_t ADD FULLTEXT KEY fk (c) GLOBAL", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL GLOBAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) GLOBAL LOCAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) GLOBAL GLOBAL)", true},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL COMMENT 'x' GLOBAL)", true},
	{"CREATE INDEX k ON obscope_t (a) LOCAL GLOBAL", true},
	{"ALTER TABLE obscope_t ADD UNIQUE KEY uk (a) GLOBAL LOCAL", true},
}

// TestOceanBaseIndexScopeGrammar holds the parser to the OceanBase verdicts above.
// TestOceanBaseIndexScopeGrammarOracle asks the engine the same statements.
func TestOceanBaseIndexScopeGrammar(t *testing.T) {
	for _, c := range append(append([]indexScopeCase{}, primaryKeyScopeCases...), indexScopeCases...) {
		_, err := Parse(c.sql)
		if (err == nil) != c.parses {
			t.Errorf("Parse(%q) error = %v, want parses %v", c.sql, err, c.parses)
			continue
		}
		if err == nil {
			continue
		}
		// A refused scope fails at the scope, as OceanBase's "near 'GLOBAL)'" does.
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", c.sql, err)
			continue
		}
		at := -1
		for _, word := range []string{" GLOBAL", " LOCAL", " BLOCK_SIZE"} {
			if i := strings.Index(c.sql, word); i >= 0 && (at < 0 || i < at) {
				at = i
			}
		}
		if want := at + 1; pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", c.sql, pe.Position, want, err)
		}
	}
}

// TestOceanBaseIndexScopeGrammarOracle asks OceanBase about every case of
// TestOceanBaseIndexScopeGrammar. The verdict comes from PREPARE, which parses
// without creating anything: 1064 is a syntax error, any other answer means the
// statement parsed.
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseIndexScopeGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obscope",
		"USE obscope",
		"DROP TABLE IF EXISTS obscope_t",
		"CREATE TABLE obscope_t (id INT NOT NULL, a INT, c VARCHAR(100), g GEOMETRY NOT NULL SRID 0)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obscope")
	})
	for _, c := range append(append([]indexScopeCase{}, primaryKeyScopeCases...), indexScopeCases...) {
		_, err := db.ExecContext(ctx, "PREPARE obscope_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
