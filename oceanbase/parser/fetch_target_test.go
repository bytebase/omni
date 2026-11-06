package parser

import "testing"

// TestFetchTargetInLaterSegment holds the user-variable check on FETCH targets to
// where the routine sits in a script: a statement after the first is parsed with
// absolute token positions over its own segment's text, so the check must read
// the byte at the position relative to the segment.
func TestFetchTargetInLaterSegment(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"DO 111;CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH c INTO x; SET @a=1; END", true},
		{"SELECT 1;CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; FETCH c INTO @x; END", false},
		{"SELECT 1;CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH c INTO x, @y; END", false},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) accepts=%v, want %v: %v", tc.sql, err == nil, tc.accept, err)
		}
	}
}
