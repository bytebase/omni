package parser

import "testing"

// TestNullBackslashN holds omni to live TiDB v8.5.5 on `\N`, MySQL's NULL
// literal (a §3.4 integrationtest corpus gap). It is case-sensitive and exactly
// two bytes: `\n` and `\M` are 1064, and `\Nx` is NULL aliased x. It deparses to
// NULL, which round-trips.
func TestNullBackslashN(t *testing.T) {
	tc := startTiDB(t)
	for _, sql := range []string{
		`select \N`,
		`select \N from test`,
		`select (\N) from test`,
		`select \N + 1`,
		`select \N, 1 from t`,
		`select \Nx`,
		`insert into t values (\N)`,
		`select \n`, // lowercase: not NULL
		`select \M`, // other letter: not NULL
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
