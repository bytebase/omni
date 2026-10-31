package parser

import "testing"

// TestTableOptionRequiresValue pins that an option routed through consumeOptionValue must
// have a value. The helper used to return an empty string with no error, so a truncated
// option parsed and review passed DDL the server answers 1064 for. The MySQL options are
// here with the OceanBase ones because the helper is shared and the hole was theirs first.
func TestTableOptionRequiresValue(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT) KEY_BLOCK_SIZE",
		"CREATE TABLE t (a INT) KEY_BLOCK_SIZE =",
		"CREATE TABLE t (a INT) STATS_SAMPLE_PAGES =;",
		"CREATE TABLE t (a INT) REPLICA_NUM =;",
		"CREATE TABLE t (a INT) BLOCK_SIZE",
		"CREATE TABLE t (a INT) USE_BLOOM_FILTER =",
		"CREATE TABLE t (a INT) TABLET_SIZE =",
		"CREATE TABLE t (a INT) AUTO_INCREMENT_MODE",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("parsed %q, but the option has no value", sql)
		}
	}
	// The same options with values must still parse, or the fix is too strict.
	for _, sql := range []string{
		"CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 8",
		"CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = DEFAULT",
		"CREATE TABLE t (a INT) REPLICA_NUM = 1 BLOCK_SIZE = 16384",
		"CREATE TABLE t (a INT) USE_BLOOM_FILTER = FALSE AUTO_INCREMENT_MODE = 'ORDER'",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%q: %v", sql, err)
		}
	}
}
