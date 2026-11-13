package parser

import "testing"

// TestOceanBaseColumnAttributes holds OceanBase's own column attributes to their
// productions. Whether each form parses is held to the server by the
// column_attribute rows of testdata/compat.tsv; these add the values no row writes.
func TestOceanBaseColumnAttributes(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT ORIG_DEFAULT 1)",
		"CREATE TABLE t (a INT ORIG_DEFAULT -1.5)",
		"CREATE TABLE t (a VARCHAR(3) ORIG_DEFAULT 'x')",
		"CREATE TABLE t (a INT ORIG_DEFAULT NULL)",
		"CREATE TABLE t (a DATETIME ORIG_DEFAULT CURRENT_TIMESTAMP(6))",
		"CREATE TABLE t (a DATETIME ORIG_DEFAULT NOW())",
		"CREATE TABLE t (a INT SKIP_INDEX(MIN_MAX, SUM) CHUNK 8)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"CREATE TABLE t (a INT ORIG_DEFAULT 1 + 2)",         // no expression
		"CREATE TABLE t (a INT ORIG_DEFAULT -CURRENT_USER)", // a sign takes a number
		"CREATE TABLE t (a INT ORIG_DEFAULT (1))",
		"CREATE TABLE t (a DATETIME ORIG_DEFAULT NOW)", // NOW takes its parentheses
		"CREATE TABLE t (a INT ORIG_DEFAULT CURRENT_USER)",
		"CREATE TABLE t (a INT SKIP_INDEX(MAX))",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
