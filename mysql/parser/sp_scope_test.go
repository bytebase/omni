package parser

import (
	"strings"
	"testing"
)

// TestStoredProgramLabels holds LEAVE, ITERATE and a block's or loop's label to the labels in
// reach, as MySQL's grammar actions do (sp_pcontext::find_label). Each refusal names the
// error mysql:8.0.46 answers to PREPARE of the statement; wantAt is the text the error points
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
					t.Fatalf("Parse refused what MySQL takes: %v", err)
				}
				return
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("Parse took what MySQL refuses with %d (err %v)", c.code, err)
			}
			if !strings.HasPrefix(c.sql[pe.Position:], c.wantAt) {
				t.Errorf("error at %q, want at %q: %v", c.sql[pe.Position:], c.wantAt, pe)
			}
		})
	}
}
