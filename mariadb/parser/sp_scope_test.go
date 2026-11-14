package parser

import (
	"strings"
	"testing"
)

// TestStoredProgramLabels holds LEAVE, ITERATE and a block's or loop's label to the labels in
// reach, as MariaDB's grammar actions do (sp_pcontext::find_label). Each refusal names the
// error mariadb:11.8.8 answers to PREPARE of the statement; wantAt is the text the error points
// at.
func TestStoredProgramLabels(t *testing.T) {
	cases := []struct {
		sql    string
		code   int // 0: the server takes it
		wantAt string
	}{
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN LEAVE a; END; END LOOP"},
		{sql: "CREATE PROCEDURE p() a: BEGIN b: BEGIN LEAVE a; END b; END a"},
		{sql: "CREATE PROCEDURE p() a: WHILE 1 DO ITERATE a; END WHILE"},
		{sql: "CREATE PROCEDURE p() a: REPEAT ITERATE A; UNTIL 1 END REPEAT"},
		{sql: "CREATE PROCEDURE p() BEGIN a: LOOP LEAVE a; END LOOP; a: LOOP LEAVE a; END LOOP; END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN DECLARE EXIT HANDLER FOR 1 a: LOOP LEAVE a; END LOOP; END; END LOOP"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 b: BEGIN LEAVE b; END; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT a: BEGIN LEAVE a; RETURN 1; END"},
		{sql: "CREATE TRIGGER t BEFORE INSERT ON t FOR EACH ROW a: BEGIN LEAVE a; END"},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO a: LOOP LEAVE a; END LOOP"},
		{sql: "ALTER EVENT e DO a: LOOP LEAVE a; END LOOP"},

		{sql: "CREATE PROCEDURE p() LEAVE a", code: 1308, wantAt: "a"},
		{sql: "CREATE PROCEDURE p() é: LOOP LEAVE e; END LOOP", code: 1308, wantAt: "e; END"},
		{sql: "ALTER EVENT e DO ITERATE a", code: 1308, wantAt: "a"},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN LEAVE b; END", code: 1308, wantAt: "b; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN ITERATE a; END", code: 1308, wantAt: "a; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN b: BEGIN END b; LEAVE b; END", code: 1308, wantAt: "b; END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN DECLARE EXIT HANDLER FOR 1 LEAVE a; END; END LOOP", code: 1308, wantAt: "a; END; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN a: BEGIN END; END", code: 1309, wantAt: "a: BEGIN END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN A: LOOP LEAVE a; END LOOP; END; END LOOP", code: 1309, wantAt: "A: LOOP"},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 {
				if err != nil {
					t.Fatalf("Parse refused what MariaDB takes: %v", err)
				}
				return
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("Parse took what MariaDB refuses with %d (err %v)", c.code, err)
			}
			if !strings.HasPrefix(c.sql[pe.Position:], c.wantAt) {
				t.Errorf("error at %q, want at %q: %v", c.sql[pe.Position:], c.wantAt, pe)
			}
		})
	}
}

// TestStoredProgramDeclarations holds parameters, variables and cursors to a stored program's
// scope (sp_pcontext::find_variable, find_cursor), and DECLARE to the head of a BEGIN ... END
// block. code is what mariadb:11.8.8 answers to PREPARE of the statement, 0 when the compat
// judge reads that as parsed (compat.PrepareVerdict).
func TestStoredProgramDeclarations(t *testing.T) {
	cases := []struct {
		sql  string
		code int
	}{
		{sql: "CREATE PROCEDURE p(a INT, b INT) BEGIN DECLARE a INT; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; BEGIN DECLARE x INT; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE x CURSOR FOR SELECT 1; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; BEGIN DECLARE c CURSOR FOR SELECT 2; OPEN c; END; CLOSE c; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE EXIT HANDLER FOR 1 CLOSE c; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CURSOR FOR SELECT 1; FETCH NEXT FROM c INTO x; END; END"},
		{sql: "CREATE PROCEDURE p(n INT) BEGIN DECLARE c CURSOR FOR SELECT 1 LIMIT n; SELECT 1 INTO n; GET DIAGNOSTICS n = NUMBER; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE next CURSOR FOR SELECT 1; FETCH NEXT FROM next INTO x; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE next CURSOR FOR SELECT 1; FETCH next INTO x; END"},
		{sql: "CREATE FUNCTION f(a INT) RETURNS INT BEGIN DECLARE b INT; SELECT 1 INTO a; SELECT 2 INTO b LIMIT a; RETURN b; END"},
		{sql: "GET DIAGNOSTICS @n = NUMBER, @r = ROW_COUNT"},
		{sql: "SELECT 1 INTO a", code: 1327},
		{sql: "SELECT 1 LIMIT a", code: 1327},
		{sql: "CREATE PROCEDURE p(a INT, A INT) BEGIN END", code: 1330},
		{sql: "CREATE PROCEDURE p(IN a INT, INOUT a INT) BEGIN END", code: 1330},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x, y, X INT; END", code: 1331},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE C CURSOR FOR SELECT 2; END", code: 1333},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1 INTO @a; END", code: 1323},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1 UNION SELECT 2 INTO @a; END", code: 1323},
		{sql: "CREATE PROCEDURE p() OPEN c", code: 1324},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CURSOR FOR SELECT 1; END; OPEN c; END", code: 1324},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; FETCH c INTO y; END", code: 1327},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH c INTO @x; END", code: 1064},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH NEXT c INTO x; END", code: 1064},
		{sql: "CREATE PROCEDURE p() GET DIAGNOSTICS a = NUMBER", code: 1327},
		{sql: "GET DIAGNOSTICS a = NUMBER", code: 1327},
		{sql: "GET DIAGNOSTICS CONDITION 1 a = MESSAGE_TEXT", code: 1327},
		{sql: "CREATE PROCEDURE p() SELECT 1 INTO a", code: 1327},
		{sql: "CREATE PROCEDURE p() BEGIN BEGIN DECLARE a INT; END; SELECT 1 INTO a; END", code: 1327},
		{sql: "CREATE TRIGGER t BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE a INT; SELECT 1 INTO a, b; END", code: 1327},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1 LIMIT n", code: 1327},
		{sql: "CREATE PROCEDURE p() DECLARE x INT", code: 1064},
		{sql: "CREATE PROCEDURE p() a: LOOP DECLARE x INT; END LOOP", code: 1064},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 DECLARE y INT; END", code: 1064},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 && err != nil {
				t.Fatalf("Parse refused what MariaDB takes: %v", err)
			}
			if c.code != 0 && err == nil {
				t.Fatalf("Parse took what MariaDB refuses with %d", c.code)
			}
		})
	}
}
