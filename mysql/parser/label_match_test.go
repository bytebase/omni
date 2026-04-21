package parser

import (
	"strings"
	"testing"
)

// TestLabelMatching enforces MySQL's end-label matching rule (ERR 1310).
// Per docs/plans/2026-04-20-mysql-routine-body-grammar.md category 2, matches
// are case-insensitive and enforced by parseBeginEndBlock, parseWhileStmt,
// parseLoopStmt, parseRepeatStmt.
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

		// Negative: mismatched labels.
		{
			name: "BEGIN end-label without begin-label",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET @a = 1;
END lbl`,
			wantErr: `end label "lbl" without matching begin label`,
		},
		{
			name: "BEGIN mismatched labels",
			sql: `CREATE PROCEDURE p()
foo: BEGIN
    SET @a = 1;
END bar`,
			wantErr: `end label "bar" does not match begin label "foo"`,
		},
		{
			name: "WHILE mismatched labels",
			sql: `CREATE PROCEDURE p()
BEGIN
    lbl1: WHILE @i < 3 DO
        SET @i = @i + 1;
    END WHILE lbl2;
END`,
			wantErr: `end label "lbl2" does not match begin label "lbl1"`,
		},
		{
			name: "LOOP end-label without begin-label",
			sql: `CREATE PROCEDURE p()
BEGIN
    foo: LOOP
        LEAVE foo;
    END LOOP myloop;
END`,
			wantErr: `end label "myloop" does not match begin label "foo"`,
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
			wantErr: `end label "r2" does not match begin label "r1"`,
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
