package catalog

import (
	"strings"
	"testing"
)

// TestEnumDuplicateValueMessage holds 1291's message to the server's: the value
// named is the first spelling, which MariaDB compares case-insensitively.
// Measured on mariadb:11.8.8: `CREATE TABLE t (e ENUM('s','S'))` answers
// "Column 'e' has duplicated value 's' in ENUM", and a SET the same with SET.
func TestEnumDuplicateValueMessage(t *testing.T) {
	for _, tc := range []struct{ ddl, want string }{
		{"CREATE TABLE t (e ENUM('s', 'S'))", "Column 'e' has duplicated value 's' in ENUM"},
		{"CREATE TABLE t (e SET('a', 'b', 'A'))", "Column 'e' has duplicated value 'a' in SET"},
		{"CREATE TABLE t (e ENUM('a ', 'a'))", "Column 'e' has duplicated value 'a' in ENUM"},
		{"CREATE TABLE t (e ENUM('S ', 's'))", "Column 'e' has duplicated value 'S' in ENUM"},
	} {
		err := execErr(t, tc.ddl)
		catErr, ok := err.(*Error)
		if !ok || catErr.Code != 1291 {
			t.Fatalf("%s: expected 1291, got %v", tc.ddl, err)
		}
		if !strings.Contains(catErr.Message, tc.want) {
			t.Errorf("%s: Message = %q, want %q", tc.ddl, catErr.Message, tc.want)
		}
	}

	// Outside strict sql_mode, 1291 is a note and the table is created: measured
	// on mariadb:11.8.8 under `SET sql_mode = ''` and `'ANSI'`.
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; SET sql_mode = ''; CREATE TABLE t (e ENUM('s', 's'));")
	if c.GetDatabase("d").GetTable("t") == nil {
		t.Fatal("non-strict sql_mode rejected ENUM('s', 's')")
	}
}
