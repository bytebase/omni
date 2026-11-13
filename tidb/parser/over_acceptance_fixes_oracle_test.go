package parser

import "testing"

// TestOverAcceptanceFixes holds omni to live TiDB v8.5.5 on the four
// over-acceptance families the §3.4 integrationtest corpus surfaced (omni
// accepted what TiDB rejects). Each statement is measured against the server and
// omni must give the same verdict, so the rejections and the forms that stay
// valid are both pinned.
func TestOverAcceptanceFixes(t *testing.T) {
	tc := startTiDB(t)
	for _, sql := range []string{
		// A. function-name keyword immediately before `(` is a function, not a
		// table name; a space before the paren, or a quoted name, is a table.
		"CREATE TABLE COUNT(a bigint)",
		"CREATE TABLE BIT_AND(a bigint)",
		"CREATE TABLE NOW(a bigint)",
		"CREATE TABLE SUM(a int)",
		"CREATE TABLE GROUP_CONCAT(a int)",
		"CREATE TABLE COUNT (a bigint)",
		"CREATE TABLE `COUNT` (a bigint)",
		"SELECT COUNT(*) FROM t",
		"SELECT * FROM COUNT",
		"CREATE TABLE t (COUNT int)",
		// B. ALTER USER USER() takes only IDENTIFIED BY <string>, not BY PASSWORD.
		"ALTER USER USER() IDENTIFIED BY PASSWORD '*B50FBDB37F1256824274912F2A1CE648082C3F1F'",
		"ALTER USER USER() IDENTIFIED BY 'plain'",
		"CREATE USER u IDENTIFIED BY PASSWORD '*hash'",
		"ALTER USER u IDENTIFIED BY PASSWORD '*hash'",
		// C1. LIMIT/OFFSET is a uint64.
		"select * from t limit 18446744073709551616 offset 3",
		"select * from t limit 18446744073709551615 offset 3",
		"select 18446744073709551616",
		// C2. a float literal that overflows the double is a syntax error.
		"insert into tt1 values(4556414e723532)",
		"select 1e400",
		"select 1e-400",
		"select 1.5e10",
	} {
		tidbAccepts, err := tc.prepareSyntax(sql)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", sql, err)
		}
		_, perr := Parse(sql)
		if omniAccepts := perr == nil; omniAccepts != tidbAccepts {
			t.Errorf("DIVERGENCE omni accepts=%v, TiDB accepts=%v for %q: %v",
				omniAccepts, tidbAccepts, sql, perr)
		}
	}
}
