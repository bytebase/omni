package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	gomysql "github.com/go-sql-driver/mysql"
)

// fulltextParserCase is a statement, whether OceanBase 4.3.5 parses it (false means
// the engine answers 1064), and for a refused one the text the error is at: its last
// occurrence in the statement.
type fulltextParserCase struct {
	sql    string
	parses bool
	at     string
}

// ftTable opens a CREATE TABLE with a fulltext-indexable column c; the case appends
// the index and its options.
const ftTable = "CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk "

// fulltextParserCases: a fulltext index takes WITH PARSER and PARSER_PROPERTIES beside
// the index options, and no other index takes either (parseIndexOption). Read against
// 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02. obft_t has columns id, a, c.
var fulltextParserCases = []fulltextParserCase{
	// The form 4.3.5.6 prints for every fulltext index in SHOW CREATE TABLE.
	{ftTable + "(c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384)", true, ""},
	{ftTable + "(c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384 COMMENT 'x')", true, ""},
	{ftTable + "(c) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=2,max_ngram_size=3) BLOCK_SIZE 16384)", true, ""},
	{ftTable + "(c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=\"smart\") BLOCK_SIZE 16384)", true, ""},
	{ftTable + "(c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) BLOCK_SIZE 16384 /*!80000 INVISIBLE */)", true, ""},
	// The list as users write it: either order, with or without '=' and commas, the keys
	// in any case, repeated, and the option before WITH PARSER or twice.
	{ftTable + "(c) PARSER_PROPERTIES=(max_token_size=50, min_token_size=2))", true, ""},
	{ftTable + "(c, d) WITH PARSER space PARSER_PROPERTIES (min_token_size=2 max_token_size=80))", true, ""},
	{ftTable + "(c) PARSER_PROPERTIES = (min_token_size = 3))", true, ""},
	{ftTable + "(c) PARSER_PROPERTIES=(MIN_TOKEN_SIZE=3, Max_Token_Size=40))", true, ""},
	{ftTable + "(c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3, min_token_size=6))", true, ""},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40))", true, ""},
	{ftTable + "(c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram)", true, ""},
	{ftTable + "(c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='max' '_word'))", true, ""},
	{ftTable + "(c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=\"max_word\", stopword_table='none'))", true, ""},
	{ftTable + "(c) WITH PARSER ik PARSER_PROPERTIES=(stopword_table='none', dict_table='none', quantifier_table='none'))", true, ""},
	{ftTable + "(c) WITH PARSER `ngram`)", true, ""},
	{ftTable + "(c) WITH PARSER NGRAM)", true, ""},
	{ftTable + "(c) WITH PARSER ngram WITH PARSER space)", true, ""},
	{ftTable + "(c) COMMENT 'x' WITH PARSER ngram LOCAL)", true, ""},
	// These parse and are then refused: a key the parser does not take is a 1235, and
	// an unknown parser a 1128.
	{ftTable + "(c) WITH PARSER space PARSER_PROPERTIES=(ngram_token_size=3))", true, ""},
	{ftTable + "(c) WITH PARSER nosuchparser)", true, ""},
	// The fulltext table element takes USING before its columns too.
	{ftTable + "USING BTREE (c))", true, ""},
	{ftTable + "(c) USING BTREE)", true, ""},
	{"CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY USING BTREE fk (c))", false, "fk (c)"},
	// The list is key '=' value with a bare key of nine, an unsigned integer for a size
	// and a string literal for the rest.
	{ftTable + "(c) PARSER_PROPERTIES=())", false, "))"},
	{ftTable + "(c) PARSER_PROPERTIES=(foo=1))", false, "foo"},
	{ftTable + "(c) PARSER_PROPERTIES=(`min_token_size`=3))", false, "`min_token_size`"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size 3))", false, "3"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size='3'))", false, "'3'"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=-3))", false, "-3"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=+3))", false, "+3"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=3.5))", false, "3.5"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=0x3))", false, "0x3"},
	{ftTable + "(c) PARSER_PROPERTIES=(stopword_table=3))", false, "3))"},
	{ftTable + "(c) PARSER_PROPERTIES=(ik_mode=smart))", false, "smart"},
	{ftTable + "(c) PARSER_PROPERTIES=(ik_mode=N'smart'))", false, "N'smart'"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=3,))", false, "))"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=3 , , max_token_size=40))", false, ", max_token_size"},
	{ftTable + "(c) PARSER_PROPERTIES min_token_size=3)", false, "min_token_size"},
	{ftTable + "(c) PARSER_PROPERTIES=)", false, ")"},
	{ftTable + "(c) PARSER_PROPERTIES)", false, ")"},
	{ftTable + "(c) PARSER_PROPERTIES=(min_token_size=3) PARSER_PROPERTIES=())", false, "))"},
	// WITH PARSER takes a name, which may be quoted with backticks but not as a string.
	{ftTable + "(c) WITH PARSER)", false, ")"},
	{ftTable + "(c) WITH PARSER 'ngram')", false, "'ngram'"},
	// CREATE FULLTEXT INDEX and ALTER TABLE ... ADD FULLTEXT take both.
	{"CREATE FULLTEXT INDEX fk ON obft_t (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=4)", true, ""},
	{"CREATE FULLTEXT INDEX fk ON obft_t (c) PARSER_PROPERTIES=(min_token_size=3) LOCAL", true, ""},
	{"CREATE FULLTEXT INDEX fk ON obft_t (c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram LOCAL", true, ""},
	{"ALTER TABLE obft_t ADD FULLTEXT KEY fk (c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=2)", true, ""},
	{"ALTER TABLE obft_t ADD FULLTEXT INDEX fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) COMMENT 'x'", true, ""},
	{"ALTER TABLE obft_t ADD FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3)", true, ""},
	{"ALTER TABLE obft_t ADD FULLTEXT fk (c) WITH PARSER ngram, ADD COLUMN zz INT", true, ""},
	// No other index takes either: WITH PARSER fails at PARSER, where OceanBase's
	// index_option could still go on to WITH ROWID.
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) WITH PARSER ngram)", false, "PARSER ngram"},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a) WITH PARSER ngram)", false, "PARSER ngram"},
	{"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g) WITH PARSER ngram)", false, "PARSER ngram"},
	{"CREATE INDEX k ON obft_t (a) WITH PARSER ngram", false, "PARSER ngram"},
	{"ALTER TABLE obft_t ADD KEY k (a) WITH PARSER ngram", false, "PARSER ngram"},
	{"ALTER TABLE obft_t ADD PRIMARY KEY (id) WITH PARSER ngram", false, "PARSER ngram"},
	{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) PARSER_PROPERTIES=(min_token_size=3))", false, "PARSER_PROPERTIES"},
	{"CREATE INDEX k ON obft_t (a) PARSER_PROPERTIES=(min_token_size=3)", false, "PARSER_PROPERTIES"},
	{"ALTER TABLE obft_t ADD UNIQUE KEY k (a) PARSER_PROPERTIES=(min_token_size=3)", false, "PARSER_PROPERTIES"},
}

// TestOceanBaseFulltextParserGrammar holds the parser to fulltextParserCases. A refused
// statement fails where OceanBase's "near" says. TestOceanBaseFulltextParserGrammarOracle
// asks the engine the same statements.
func TestOceanBaseFulltextParserGrammar(t *testing.T) {
	for _, c := range fulltextParserCases {
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

// TestOceanBaseParserPropertiesAST pins the tree PARSER_PROPERTIES gives: one
// IndexOption holding each property in order, keys lower-cased, sizes as integers and
// adjacent string literals joined.
func TestOceanBaseParserPropertiesAST(t *testing.T) {
	const sql = "CREATE TABLE t (c TEXT, FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(IK_MODE='max' '_word', min_token_size=4))"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	opts := list.Items[0].(*nodes.CreateTableStmt).Constraints[0].IndexOptions
	if len(opts) != 2 || opts[0].Name != "PARSER" || opts[1].Name != "PARSER_PROPERTIES" {
		t.Fatalf("options = %+v, want PARSER then PARSER_PROPERTIES", opts)
	}
	if got := sql[opts[1].Loc.Start:opts[1].Loc.End]; got != "PARSER_PROPERTIES=(IK_MODE='max' '_word', min_token_size=4)" {
		t.Errorf("the option spans %q", got)
	}
	props := opts[1].Properties
	if len(props) != 2 {
		t.Fatalf("%d properties, want 2", len(props))
	}
	if s, ok := props[0].Value.(*nodes.StringLit); props[0].Name != "ik_mode" || !ok || s.Value != "max_word" {
		t.Errorf("property 0 = %s %#v, want ik_mode 'max_word'", props[0].Name, props[0].Value)
	}
	if n, ok := props[1].Value.(*nodes.IntLit); props[1].Name != "min_token_size" || !ok || n.Value != 4 {
		t.Errorf("property 1 = %s %#v, want min_token_size 4", props[1].Name, props[1].Value)
	}
	if got := sql[props[0].Loc.Start:props[0].Loc.End]; got != "IK_MODE='max' '_word'" {
		t.Errorf("property 0 spans %q", got)
	}
}

// TestOceanBaseFulltextParserGrammarOracle asks OceanBase about fulltextParserCases
// through PREPARE, which parses without creating anything: 1064 is a syntax error, any
// other answer means the statement parsed.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseFulltextParserGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obft",
		"USE obft",
		"DROP TABLE IF EXISTS obft_t",
		"CREATE TABLE obft_t (id INT NOT NULL, a INT, c VARCHAR(100))",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obft")
	})
	for _, c := range fulltextParserCases {
		_, err := db.ExecContext(ctx, "PREPARE obft_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
