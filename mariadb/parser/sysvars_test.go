package parser

import "testing"

// TestUnknownSystemVariable pins which SET targets of a stored program read as system
// variables. Each refused case answers 1193 on mariadb:11.8.8 at CREATE time; the others
// create, or are not checked while the program is parsed.
func TestUnknownSystemVariable(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want string // the unknown name, or "" for none
	}{
		{"CREATE PROCEDURE p() SET zzz = 1", "zzz"},
		{"CREATE PROCEDURE p() BEGIN DECLARE a INT; SET a = 1, zzz = 2; END", "zzz"},
		{"CREATE FUNCTION f() RETURNS INT BEGIN SET zzz = 1; RETURN 1; END", "zzz"},
		{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET zzz = 1", "zzz"},
		{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SET zzz = 1", "zzz"},
		{"CREATE PROCEDURE p() BEGIN BEGIN DECLARE zzz INT; END; SET zzz = 1; END", "zzz"},
		{"CREATE PROCEDURE p() SET SQL_MODE = ''", ""},
		{"CREATE PROCEDURE p() SET server_audit_logging = ON", ""},
		{"CREATE PROCEDURE p(zzz INT) SET zzz = 1", ""},
		{"CREATE PROCEDURE p() BEGIN DECLARE zzz INT; BEGIN SET zzz = 1; END; END", ""},
		{"CREATE PROCEDURE p() FOR zzz IN 1..3 DO SET zzz = 1; END FOR", ""},
		{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.zzz = 1", ""},
		{"CREATE PROCEDURE p() SET @zzz = 1", ""},
		// MariaDB refuses these too; only a bare name is checked.
		{"CREATE PROCEDURE p() SET SESSION zzz = 1", ""},
		{"CREATE PROCEDURE p() SET @@zzz = 1", ""},
		{"SET zzz = 1", ""},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		got := ""
		if a := UnknownSystemVariable(list); a != nil {
			got = a.Column.Column
		}
		if got != tc.want {
			t.Errorf("%s: unknown system variable %q, want %q", tc.sql, got, tc.want)
		}
	}
}
