package parser

import "testing"

// TestSpatialOptions pins what mariadb:11.8.8 parses of a spatial column's
// REF_SYSTEM_ID and of a SPATIAL key's parts: REF_SYSTEM_ID = n right after a
// spatial type, once, and one key part.
func TestSpatialOptions(t *testing.T) {
	for _, tt := range []struct {
		sql string
		ok  bool
	}{
		{"CREATE TABLE t (a POINT REF_SYSTEM_ID = 4326 NOT NULL)", true},
		{"CREATE TABLE t (a GEOMETRY REF_SYSTEM_ID=0)", true},
		{"ALTER TABLE t ADD COLUMN g POLYGON REF_SYSTEM_ID = 3857", true},
		{"CREATE TABLE t (a INT REF_SYSTEM_ID = 1)", false},
		{"CREATE TABLE t (a POINT NOT NULL REF_SYSTEM_ID = 1)", false},
		{"CREATE TABLE t (a POINT REF_SYSTEM_ID = 1 REF_SYSTEM_ID = 2)", false},
		{"CREATE TABLE t (a POINT REF_SYSTEM_ID = -1)", false},
		{"CREATE TABLE t (a INT) REF_SYSTEM_ID = 1", false},
		{"CREATE TABLE t (a INT `REF_SYSTEM_ID` = 1)", true},
		{"CREATE TABLE t (g POINT NOT NULL, SPATIAL KEY (g))", true},
		{"CREATE TABLE t (g POINT NOT NULL, a INT, SPATIAL KEY (g, a))", false},
		{"ALTER TABLE t ADD SPATIAL INDEX s (g, a)", false},
		{"CREATE SPATIAL INDEX s ON t (g)", true},
		{"CREATE SPATIAL INDEX s ON t (g, g)", false},
	} {
		_, err := Parse(tt.sql)
		if (err == nil) != tt.ok {
			t.Errorf("%s: got %v, want ok=%v", tt.sql, err, tt.ok)
		}
	}
}
