package parser

import "testing"

// TestIntervalSignedValue holds omni to live TiDB v8.5.5 on a signed INTERVAL
// value (the larger half of the §3.4 integrationtest corpus's INTERVAL gap, 37
// rows). TiDB's `INTERVAL Expression TimeUnit` takes a leading sign, which omni
// read only as an unsigned primary before. A signed literal deparses as `-2.2`
// (no parentheses), so it stays the interval literal rather than becoming the
// INTERVAL function.
//
// The binary-expression value (`INTERVAL 1 << 37 SECOND`) stays a gap: its
// deparse would parenthesise the value, which TiDB reads as the function.
func TestIntervalSignedValue(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, sql := range []string{
		"SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2.2 MINUTE_MICROSECOND)",
		"SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2.2 SECOND)",
		"SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2 YEAR_MONTH)",
		"SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2 DAY_SECOND)",
		"SELECT '1900-01-01 00:00:00' + INTERVAL -1 DAY",
		"SELECT '1900-01-01 00:00:00' + INTERVAL +1 HOUR",
		// still-valid unsigned forms
		"SELECT DATE_ADD('x', INTERVAL 2.2 MINUTE_MICROSECOND)",
		"SELECT '1900-01-01' + INTERVAL 1 DAY",
	} {
		tidbAccepts, _, err := tc.syntaxVerdict(t, sql, "<interval signed test>")
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
