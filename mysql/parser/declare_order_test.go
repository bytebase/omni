package parser

import (
	"strings"
	"testing"
)

// TestDeclareOrdering enforces MySQL's DECLARE ordering rule in BEGIN...END
// blocks (ERR 1337 / ERR 1338). Per docs/plans/
// 2026-04-20-mysql-routine-body-grammar.md category 3, the ordering is:
//   DECLARE var/condition  →  DECLARE cursor  →  DECLARE handler  →  stmts
// Each kind may advance the phase but not regress. Per-BEGIN scope.
// Container-verified against MySQL 8.0 on 2026-04-20.
func TestDeclareOrdering(t *testing.T) {
	type tc struct {
		name    string
		sql     string
		wantErr string // non-empty ⇒ expect error containing this substring
	}
	cases := []tc{
		// Positive: canonical orderings.
		{
			name: "var → cursor → handler → stmt",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT DEFAULT 0;
    DECLARE cur CURSOR FOR SELECT 1;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    OPEN cur;
    CLOSE cur;
END`,
		},
		{
			name: "condition → handler (skip cursor)",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key SET @err = 1;
    INSERT INTO t(id) VALUES (1);
END`,
		},
		{
			name: "var → handler (skip cursor)",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    SELECT 1;
END`,
		},
		{
			name: "multiple vars then multiple cursors then multiple handlers",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE a INT;
    DECLARE b INT;
    DECLARE c1 CURSOR FOR SELECT 1;
    DECLARE c2 CURSOR FOR SELECT 2;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION SET @err = 1;
    SELECT 1;
END`,
		},
		{
			name: "nested BEGIN resets phase",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT DEFAULT 0;
    SELECT 1;
    BEGIN
        DECLARE y INT DEFAULT 0;
        SELECT 2;
    END;
END`,
		},
		{
			name: "only DECLAREs, no statements",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT;
    DECLARE c CURSOR FOR SELECT 1;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION SET @err = 1;
END`,
		},

		// Negative: order violations.
		{
			name: "var after cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE cur CURSOR FOR SELECT 1;
    DECLARE x INT;
    OPEN cur;
END`,
			wantErr: "variable or condition declaration after cursor or handler declaration",
		},
		{
			name: "condition after cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE cur CURSOR FOR SELECT 1;
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    OPEN cur;
END`,
			wantErr: "variable or condition declaration after cursor or handler declaration",
		},
		{
			name: "var after handler",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    DECLARE x INT;
    SELECT 1;
END`,
			wantErr: "variable or condition declaration after cursor or handler declaration",
		},
		{
			name: "cursor after handler",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    DECLARE cur CURSOR FOR SELECT 1;
    OPEN cur;
END`,
			wantErr: "cursor declaration after handler declaration",
		},
		{
			name: "handler after regular statement",
			sql: `CREATE PROCEDURE p()
BEGIN
    SELECT 1;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    SELECT 2;
END`,
			wantErr: "handler declaration after regular statement",
		},
		{
			name: "cursor after regular statement",
			sql: `CREATE PROCEDURE p()
BEGIN
    SELECT 1;
    DECLARE cur CURSOR FOR SELECT 2;
END`,
			wantErr: "cursor declaration after handler declaration",
		},
		{
			name: "var after regular statement",
			sql: `CREATE PROCEDURE p()
BEGIN
    SELECT 1;
    DECLARE x INT;
END`,
			wantErr: "variable or condition declaration after cursor or handler declaration",
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
