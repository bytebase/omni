package parser

import (
	"strings"
	"testing"
)

// TestLabelMatching holds the end-label rule OceanBase's parser enforces: an end
// label needs a begin label. Two labels that differ parse, because OceanBase
// compares them in its resolver (1308, see checkLabelMatch), so that mismatch
// is oceanbase/validate's label_mismatch.
func TestLabelMatching(t *testing.T) {
	type tc struct {
		name    string
		sql     string
		wantErr string // non-empty ⇒ expect error containing this substring
	}
	cases := []tc{
		// Positive: matching labels.
		{
			name: "BEGIN exact-case match",
			sql: `CREATE PROCEDURE p()
myblock: BEGIN
    SET @a = 1;
END myblock`,
		},
		{
			name: "BEGIN case-insensitive match",
			sql: `CREATE PROCEDURE p()
myBlock: BEGIN
    SET @a = 1;
END MYBLOCK`,
		},
		{
			name: "BEGIN end-label omitted",
			sql: `CREATE PROCEDURE p()
myblock: BEGIN
    SET @a = 1;
END`,
		},
		{
			name: "BEGIN both labels absent",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET @a = 1;
END`,
		},
		{
			name: "labeled WHILE with matching end label",
			sql: `CREATE PROCEDURE p()
BEGIN
    lbl: WHILE @i < 3 DO
        SET @i = @i + 1;
    END WHILE lbl;
END`,
		},
		{
			name: "labeled LOOP with matching end label",
			sql: `CREATE PROCEDURE p()
BEGIN
    myloop: LOOP
        LEAVE myloop;
    END LOOP myloop;
END`,
		},
		{
			name: "labeled REPEAT with matching end label",
			sql: `CREATE PROCEDURE p()
BEGIN
    r1: REPEAT
        SET @i = @i + 1;
    UNTIL @i >= 3
    END REPEAT r1;
END`,
		},

		// Negative: an end label with no begin label.
		{
			name: "BEGIN end-label without begin-label",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET @a = 1;
END lbl`,
			wantErr: `end label "lbl" without matching begin label`,
		},
		{
			// Mismatched labels parse; oceanbase/validate reports them.
			name: "BEGIN mismatched labels",
			sql: `CREATE PROCEDURE p()
foo: BEGIN
    SET @a = 1;
END bar`,
		},
		{
			name: "WHILE mismatched labels",
			sql: `CREATE PROCEDURE p()
BEGIN
    lbl1: WHILE @i < 3 DO
        SET @i = @i + 1;
    END WHILE lbl2;
END`,
		},
		{
			name: "LOOP mismatched labels",
			sql: `CREATE PROCEDURE p()
BEGIN
    foo: LOOP
        LEAVE foo;
    END LOOP myloop;
END`,
		},
		{
			name: "REPEAT mismatched labels",
			sql: `CREATE PROCEDURE p()
BEGIN
    r1: REPEAT
        SET @i = @i + 1;
    UNTIL @i >= 3
    END REPEAT r2;
END`,
		},
		{
			name: "WHILE end-label without begin-label",
			sql: `CREATE PROCEDURE p()
BEGIN
    WHILE @i < 3 DO
        SET @i = @i + 1;
    END WHILE lbl;
END`,
			wantErr: `end label "lbl" without matching begin label`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("Parse failed unexpectedly: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", c.wantErr)
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), c.wantErr)
			}
		})
	}
}
