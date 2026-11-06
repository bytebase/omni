package parser

import (
	"sort"
	"testing"
)

// TestNameTokenFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// nameTokenFormCases and flushFormCases (container test).
func TestNameTokenFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, cases := range [][]parserFormCase{nameTokenFormCases, flushFormCases} {
		for _, tc := range cases {
			t.Run(tc.sql, func(t *testing.T) {
				if got := o.prepareParses(t, tc.sql); got != tc.accept {
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
}

// TestIntroducerWordsMatchMariaDB sweeps every character set and collation the server
// lists, and every name charsetIntroducerNames holds, through a name position as an
// underscored word, and holds omni's verdict to the server's. A word the server reads as
// an introducer is a 1064 there and nowhere else, so this catches a set that is missing a
// character set or holds one the server lacks. After a separating dot no word introduces,
// which each introducer is checked for too. The qualifier names no database: the verdict
// is the parser's, and an unknown database is refused only after parsing.
func TestIntroducerWordsMatchMariaDB(t *testing.T) {
	o := startMariaDB(t)
	words := map[string]bool{"utf8": true, "filename": true, "gb18030": true, "default": true, "foo": true}
	for name := range charsetIntroducerNames {
		words[name] = true
	}
	rows, err := o.db.QueryContext(o.ctx,
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
			sqls = append(sqls, "SELECT * FROM d._"+w)
		}
		for _, sql := range sqls {
			server := o.prepareParses(t, sql)
			_, err := Parse(sql)
			if omni := err == nil; omni != server {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v", omni, server, sql, err)
			}
			if !server {
				introducers++
			}
		}
	}
	// The 40 character sets of SHOW CHARACTER SET, utf8 and filename.
	if introducers != 42 {
		t.Errorf("MariaDB refused %d words as introducers, want 42", introducers)
	}
}
