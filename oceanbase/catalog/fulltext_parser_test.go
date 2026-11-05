package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// fulltextTable wraps a fulltext index definition in the table every case uses.
func fulltextTable(index string) string {
	return "CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), " + index + ")"
}

// fulltextEchoCases pair a fulltext index as DDL writes it with the line 4.3.5.6 prints
// for it in SHOW CREATE TABLE (read on 2026-11-02). TestOracle_OceanBaseFulltextParserEcho
// asks the engine for each readback again.
var fulltextEchoCases = []struct{ user, readback string }{
	{"FULLTEXT KEY fk (c)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram2",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=2,max_ngram_size=3) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER beng",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ik",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=\"smart\") BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=5)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=5,max_token_size=84) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) PARSER_PROPERTIES=(max_token_size=50, min_token_size=2)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=2,max_token_size=50) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) COMMENT 'x' WITH PARSER ngram LOCAL",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384 COMMENT 'x'"},
	{"FULLTEXT KEY fk (c) WITH PARSER NGRAM",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='MAX_WORD')",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=\"max_word\") BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='max_word', stopword_table='none')",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=\"max_word\") BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=1, max_ngram_size=4)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=1,max_ngram_size=4) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) INVISIBLE",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) BLOCK_SIZE 16384 /*!80000 INVISIBLE */"},
	{"FULLTEXT KEY fk (c, d) WITH PARSER space PARSER_PROPERTIES (min_token_size=2 max_token_size=80)",
		"FULLTEXT KEY `fk` (`c`, `d`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=2,max_token_size=80) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=40) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3, min_token_size=6)",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=6,max_token_size=84) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram WITH PARSER space",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384"},
	{"FULLTEXT KEY fk (c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram",
		"FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384"},
}

// fulltextChangeCases pair two fulltext indexes OceanBase stores differently, which the
// diff must report as a change.
var fulltextChangeCases = []struct{ from, to string }{
	{"FULLTEXT KEY fk (c)", "FULLTEXT KEY fk (c) WITH PARSER ngram"},
	{"FULLTEXT KEY fk (c) WITH PARSER ngram", "FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3)"},
	{"FULLTEXT KEY fk (c) WITH PARSER ik", "FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='max_word')"},
	{"FULLTEXT KEY fk (c)", "FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=4)"},
	{"FULLTEXT KEY fk (c) WITH PARSER space", "FULLTEXT KEY fk (c) WITH PARSER beng"},
}

// indexDiffs returns the index entries of a diff from one table definition to another.
func indexDiffs(t *testing.T, fromDDL, toDDL string) []IndexDiffEntry {
	t.Helper()
	sc := serverCharsetFor(MySQL80)
	from, _ := loadOneTable(t, sc, fromDDL, "t")
	to, _ := loadOneTable(t, sc, toDDL, "t")
	var out []IndexDiffEntry
	for _, te := range DiffWithNormalizer(from, to, NormalizerFor(MySQL80)).Tables {
		out = append(out, te.Indexes...)
	}
	return out
}

// TestFulltextParserReadbackCanonicalizes pins that each fulltext index and the line
// OceanBase prints for it compare equal, in both directions, and that the readback
// loads and compares equal to itself.
func TestFulltextParserReadbackCanonicalizes(t *testing.T) {
	for _, tc := range fulltextEchoCases {
		t.Run(tc.user, func(t *testing.T) {
			user, stored := fulltextTable(tc.user), fulltextTable(tc.readback)
			for _, pair := range [][2]string{{user, stored}, {stored, user}, {stored, stored}} {
				if d := indexDiffs(t, pair[0], pair[1]); len(d) > 0 {
					t.Errorf("%s %s:\n  from: %s\n  to:   %s", d[0].Name, d[0].Action, pair[0], pair[1])
				}
			}
		})
	}
}

// TestFulltextParserChangeDiffs pins that a change of parser or of a property the engine
// shows is a change, and that the migration carries the new parser.
func TestFulltextParserChangeDiffs(t *testing.T) {
	for _, tc := range fulltextChangeCases {
		t.Run(tc.from+" -> "+tc.to, func(t *testing.T) {
			if d := indexDiffs(t, fulltextTable(tc.from), fulltextTable(tc.to)); len(d) != 1 || d[0].Action != DiffModify {
				t.Errorf("diff = %+v, want one index MODIFY", d)
			}
			sc := serverCharsetFor(MySQL80)
			from, _ := loadOneTable(t, sc, fulltextTable(tc.from), "t")
			to, toTbl := loadOneTable(t, sc, fulltextTable(tc.to), "t")
			n := NormalizerFor(MySQL80)
			plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n).SQL()
			var ft *Index
			for _, idx := range toTbl.Indexes {
				if idx.Fulltext {
					ft = idx
				}
			}
			if want := formatFulltextParser(ft); !strings.Contains(plan, want) {
				t.Errorf("the migration does not carry %q:\n%s", want, plan)
			}
		})
	}
}

// TestFulltextParserRendering pins how the catalog renders a fulltext parser: as the DDL
// wrote it, keys in OceanBase's order, ik_mode lower-cased, a later list replacing an
// earlier one, and nothing for an index that names none. The rendering parses back to
// the same index.
func TestFulltextParserRendering(t *testing.T) {
	for _, tc := range []struct{ index, want string }{
		{"FULLTEXT KEY fk (c)", "FULLTEXT KEY `fk` (`c`)"},
		{"FULLTEXT KEY fk (c) WITH PARSER NGRAM", "FULLTEXT KEY `fk` (`c`) WITH PARSER ngram"},
		{"FULLTEXT KEY fk (c) PARSER_PROPERTIES=(max_token_size=50, MIN_TOKEN_SIZE=2) COMMENT 'x'",
			"FULLTEXT KEY `fk` (`c`) PARSER_PROPERTIES=(min_token_size=2,max_token_size=50) COMMENT 'x'"},
		{"FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(stopword_table='none', ik_mode='MAX_WORD')",
			"FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='max_word',stopword_table='none')"},
		{"FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40)",
			"FULLTEXT KEY `fk` (`c`) PARSER_PROPERTIES=(max_token_size=40)"},
	} {
		t.Run(tc.index, func(t *testing.T) {
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+fulltextTable(tc.index)+";", nil); err != nil {
				t.Fatal(err)
			}
			ddl := c.ShowCreateTable("d", "t")
			if !strings.Contains(ddl, tc.want+",") && !strings.Contains(ddl, tc.want+"\n") {
				t.Errorf("SHOW CREATE TABLE renders\n%s\nwant the line %s", ddl, tc.want)
			}
			if d := indexDiffs(t, fulltextTable(tc.index), ddl); len(d) > 0 {
				t.Errorf("the rendering does not read back as the index: %s %s\n%s", d[0].Name, d[0].Action, ddl)
			}
		})
	}
}

// fulltextLikeSource is a table whose fulltext index names a parser and its properties,
// which CREATE TABLE ... LIKE copies.
const fulltextLikeSource = "CREATE TABLE src (id INT NOT NULL, c VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3))"

// TestFulltextParserCreateTableLike pins that CREATE TABLE ... LIKE copies a fulltext
// index's parser and properties, and that the copy does not share them with the source.
func TestFulltextParserCreateTableLike(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d; "+fulltextLikeSource+"; CREATE TABLE t LIKE src;", nil); err != nil {
		t.Fatal(err)
	}
	var src, dst *Index
	for _, idx := range c.GetDatabase("d").GetTable("src").Indexes {
		if idx.Fulltext {
			src = idx
		}
	}
	for _, idx := range c.GetDatabase("d").GetTable("t").Indexes {
		if idx.Fulltext {
			dst = idx
		}
	}
	if dst == nil || dst.Parser != "ngram" || dst.ParserProperties["ngram_token_size"] != "3" {
		t.Fatalf("the copy's fulltext index is %+v, want ngram with ngram_token_size 3", dst)
	}
	dst.ParserProperties["ngram_token_size"] = "4"
	if src.ParserProperties["ngram_token_size"] != "3" {
		t.Error("the copy shares its parser properties with the source")
	}
}

// TestOracle_OceanBaseFulltextParserEcho creates each of fulltextEchoCases on OceanBase,
// requires the engine to print the recorded line, and requires the DDL and the readback
// to compare equal on their indexes. Each of fulltextChangeCases, applied as from and
// read back, must still differ from to.
//
//	go test ./oceanbase/catalog/ -run OceanBase -count=1
func TestOracle_OceanBaseFulltextParserEcho(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for i, tc := range fulltextEchoCases {
		t.Run(tc.user, func(t *testing.T) {
			dbName := fmt.Sprintf("obft_%d", i)
			readback, err := storeAndRead(ctx, o, dbName, fulltextTable(tc.user), "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the index: %v", err)
			}
			t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })
			if !strings.Contains(readback, tc.readback) {
				t.Errorf("OceanBase prints\n%s\nthe case records %s", readback, tc.readback)
			}
			user := fulltextTable(tc.user)
			for _, pair := range [][2]string{{user, readback}, {readback, user}, {readback, readback}} {
				if d := indexDiffs(t, pair[0], pair[1]); len(d) > 0 {
					t.Errorf("%s %s:\n  from: %s\n  to:   %s", d[0].Name, d[0].Action, pair[0], pair[1])
				}
			}
		})
	}
	t.Run("create table like", func(t *testing.T) {
		if _, err := storeAndRead(ctx, o, "obft_like", fulltextLikeSource, "src"); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS obft_like") })
		if _, err := o.db.ExecContext(ctx, "CREATE TABLE t LIKE src"); err != nil {
			t.Fatalf("OceanBase refused CREATE TABLE LIKE: %v", err)
		}
		var readback string
		if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE obft_like.t").Scan(new(string), &readback); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(readback, "WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3)") {
			t.Errorf("OceanBase does not copy the parser:\n%s", readback)
		}
		if d := indexDiffs(t, fulltextLikeSource+";\nCREATE TABLE t LIKE src;", readback); len(d) > 0 {
			t.Errorf("%s %s against the copy's readback\n%s", d[0].Name, d[0].Action, readback)
		}
	})
	for i, tc := range fulltextChangeCases {
		t.Run(tc.from+" -> "+tc.to, func(t *testing.T) {
			dbName := fmt.Sprintf("obftchg_%d", i)
			readback, err := storeAndRead(ctx, o, dbName, fulltextTable(tc.from), "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the index: %v", err)
			}
			t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })
			if d := indexDiffs(t, readback, fulltextTable(tc.to)); len(d) != 1 || d[0].Action != DiffModify {
				t.Errorf("diff from the readback = %+v, want one index MODIFY\n%s", d, readback)
			}
		})
	}
}
