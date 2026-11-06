package catalog

import "testing"

// TestLegacyTypeWords: the legacy type words store what TiDB v8.5.5 stores for
// them (SHOW CREATE TABLE, 2026-11-06). ASCII is latin1 and BYTE binary, so TEXT
// BYTE is a blob; YEAR keeps none of its modifiers; any UNSIGNED or ZEROFILL in a
// list makes a number unsigned, BOOL included; STORAGE is ignored; and UNICODE
// names ucs2, which TiDB parses and then refuses with 1115.
func TestLegacyTypeWords(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t ("+
		"a TEXT ASCII, b TEXT BYTE, c LONG BYTE, d LONG ASCII, e ENUM('x') ASCII, f SET('x') BYTE, "+
		"g MEDIUMTEXT ASCII COLLATE latin1_bin, h YEAR ZEROFILL, i YEAR(4) SIGNED UNSIGNED ZEROFILL, "+
		"j INT SIGNED UNSIGNED, k INT ZEROFILL UNSIGNED, l BOOL UNSIGNED, m DECIMAL(5,2) ZEROFILL, "+
		"n INT STORAGE DEFAULT, o LONG VARCHAR BYTE, p TINYTEXT BYTE, q BOOL ZEROFILL, r BOOL SIGNED); "+
		"CREATE TABLE u (a TEXT UNICODE);", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results[:3] {
		if r.Error != nil {
			t.Fatalf("statement %d: %v", r.Index, r.Error)
		}
	}
	if e, ok := results[3].Error.(*Error); !ok || e.Code != 1115 {
		t.Errorf("CREATE TABLE u (a TEXT UNICODE): err = %v, want 1115", results[3].Error)
	}
	want := "CREATE TABLE `t` (\n" +
		"  `a` text CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,\n" +
		"  `b` blob DEFAULT NULL,\n" +
		"  `c` mediumblob DEFAULT NULL,\n" +
		"  `d` mediumtext CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,\n" +
		"  `e` enum('x') CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,\n" +
		"  `f` set('x') CHARACTER SET binary COLLATE binary DEFAULT NULL,\n" +
		"  `g` mediumtext CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,\n" +
		"  `h` year(4) DEFAULT NULL,\n" +
		"  `i` year(4) DEFAULT NULL,\n" +
		"  `j` int unsigned DEFAULT NULL,\n" +
		"  `k` int(10) unsigned zerofill DEFAULT NULL,\n" +
		"  `l` tinyint(1) unsigned DEFAULT NULL,\n" +
		"  `m` decimal(5,2) unsigned zerofill DEFAULT NULL,\n" +
		"  `n` int DEFAULT NULL,\n" +
		"  `o` mediumblob DEFAULT NULL,\n" +
		"  `p` tinyblob DEFAULT NULL,\n" +
		"  `q` tinyint(1) unsigned zerofill DEFAULT NULL,\n" +
		"  `r` tinyint(1) DEFAULT NULL\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
	if got := c.ShowCreateTable("d", "t"); got != want {
		t.Errorf("SHOW CREATE TABLE =\n%s\nwant\n%s", got, want)
	}
}
