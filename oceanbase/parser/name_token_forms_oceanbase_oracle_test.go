package parser

import (
	"context"
	"database/sql"
	"sort"
	"testing"
)

// nameTokenOceanBase opens the OceanBase oracle with a database of its own selected, which
// the dot cases need: they name a table in it.
func nameTokenOceanBase(t *testing.T) (*sql.DB, context.Context) {
	t.Helper()
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS omni_name_token_forms",
		"USE omni_name_token_forms",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return db, ctx
}

// TestNameTokenFormsOceanBaseOracle asserts omni agrees with a live OceanBase CE 4.3.5.6
// on nameTokenFormCases and flushFormCases (DSN-gated, see oceanBaseOracle). The verdict
// comes from PREPARE, which parses the statement without running it.
func TestNameTokenFormsOceanBaseOracle(t *testing.T) {
	db, ctx := nameTokenOceanBase(t)
	for _, cases := range [][]parserFormCase{nameTokenFormCases, flushFormCases} {
		for _, tc := range cases {
			t.Run(tc.sql, func(t *testing.T) {
				if got := prepareAcceptsOceanBase(t, ctx, db, tc.sql); got != tc.accept {
					t.Fatalf("OceanBase accepts=%v, want %v for %q", got, tc.accept, tc.sql)
				}
				_, err := Parse(tc.sql)
				if omniAccepts := err == nil; omniAccepts != tc.accept {
					t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
						omniAccepts, tc.accept, tc.sql, err)
				}
			})
		}
	}
}

// mysqlCharacterSets are the 41 character sets of SHOW CHARACTER SET on MySQL 8.0, most
// of which OceanBase lacks, swept to show they are names there.
var mysqlCharacterSets = []string{
	"armscii8", "ascii", "big5", "binary", "cp1250", "cp1251", "cp1256", "cp1257", "cp850",
	"cp852", "cp866", "cp932", "dec8", "eucjpms", "euckr", "gb18030", "gb2312", "gbk",
	"geostd8", "greek", "hebrew", "hp8", "keybcs2", "koi8r", "koi8u", "latin1", "latin2",
	"latin5", "latin7", "macce", "macroman", "sjis", "swe7", "tis620", "ucs2", "ujis",
	"utf16", "utf16le", "utf32", "utf8mb3", "utf8mb4",
}

// TestIntroducerWordsMatchOceanBase sweeps every character set and collation the server
// lists, MySQL's character sets, and every name charsetIntroducerNames holds, through a
// name position as an underscored word, and holds omni's verdict to the server's. A word
// the server reads as an introducer is a 1064 there and nowhere else, so this catches a
// set that is missing one or holds one the server lacks. Joined to another name by a dot
// no word introduces, which each introducer is checked for too.
func TestIntroducerWordsMatchOceanBase(t *testing.T) {
	db, ctx := nameTokenOceanBase(t)
	words := map[string]bool{"utf8": true, "filename": true, "default": true, "foo": true}
	for _, name := range mysqlCharacterSets {
		words[name] = true
	}
	for name := range charsetIntroducerNames {
		words[name] = true
	}
	rows, err := db.QueryContext(ctx,
		"SELECT CHARACTER_SET_NAME FROM information_schema.CHARACTER_SETS "+
			"UNION SELECT COLLATION_NAME FROM information_schema.COLLATIONS")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		words[name] = true
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	sorted := make([]string, 0, len(words))
	for w := range words {
		sorted = append(sorted, w)
	}
	sort.Strings(sorted)

	introducers := 0
	for _, w := range sorted {
		sqls := []string{"CREATE TABLE _" + w + " (b int)"}
		if charsetIntroducerNames[w] {
			sqls = append(sqls, "SELECT * FROM t._"+w, "SELECT * FROM _"+w+".t")
		}
		for _, sql := range sqls {
			server := prepareAcceptsOceanBase(t, ctx, db, sql)
			_, err := Parse(sql)
			if omni := err == nil; omni != server {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v", omni, server, sql, err)
			}
			if !server {
				introducers++
			}
		}
	}
	// The 26 introducer tokens of sql_parser_mysql_mode.l.
	if introducers != 26 {
		t.Errorf("OceanBase refused %d words as introducers, want 26", introducers)
	}
}
