package parser

import (
	"strings"
	"testing"
)

// TestSemanticValidation locks in the inline static-validation behavior
// added to mirror MySQL's sp_pcontext / sp_head::check_return semantics.
//
// Each case asserts: parser accepts a valid construct, or rejects an
// invalid one with a specific error substring. Categories cover the
// behavior the routine-body audit (TestRoutineBodyAudit) calls out as
// V-class plus the scope-barrier / kind-distinction edges that bare
// audit probes don't exercise.
func TestSemanticValidation(t *testing.T) {
	type tc struct {
		name    string
		sql     string
		wantErr string // empty = expect accept
	}
	cases := []tc{
		// ---------- Symbol resolution: labels (LEAVE / ITERATE) ----------
		{
			name: "LEAVE outer label from nested loop",
			sql: `CREATE PROCEDURE p()
BEGIN
    outer_l: LOOP
        inner_l: WHILE @i < 3 DO
            LEAVE outer_l;
        END WHILE inner_l;
    END LOOP outer_l;
END`,
		},
		{
			name: "LEAVE BEGIN label (BEGIN-kind labels are leave-able)",
			sql: `CREATE PROCEDURE p()
my_block: BEGIN
    LEAVE my_block;
END my_block`,
		},
		{
			name: "ITERATE BEGIN label is rejected (loop-only)",
			sql: `CREATE PROCEDURE p()
my_block: BEGIN
    ITERATE my_block;
END my_block`,
			wantErr: "ITERATE references undeclared loop label",
		},
		{
			name: "ITERATE loop label",
			sql: `CREATE PROCEDURE p()
BEGIN
    lp: LOOP
        ITERATE lp;
    END LOOP lp;
END`,
		},
		{
			name: "LEAVE undeclared label",
			sql: `CREATE PROCEDURE p()
BEGIN
    LEAVE nowhere;
END`,
			wantErr: `LEAVE references undeclared label: nowhere`,
		},
		{
			name: "label barrier — LEAVE outer label from inside HANDLER body",
			sql: `CREATE PROCEDURE p()
my_block: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        LEAVE my_block;
    END;
    SELECT 1;
END my_block`,
			wantErr: `LEAVE references undeclared label: my_block`,
		},

		// ---------- Symbol resolution: variables ----------
		{
			name: "SET local var (declared)",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT;
    SET x = 1;
END`,
		},
		{
			name: "SET parameter resolves",
			sql: `CREATE PROCEDURE p(IN n INT)
BEGIN
    SET n = n + 1;
END`,
		},
		{
			name: "SET undeclared bare ident is rejected",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET nope = 1;
END`,
			wantErr: `undeclared variable: nope`,
		},
		{
			name: "SET @user_var no scope check",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET @anything = 1;
END`,
		},
		{
			name: "SET SESSION sql_mode no scope check",
			sql: `CREATE PROCEDURE p()
BEGIN
    SET SESSION sql_mode = '';
END`,
		},

		// ---------- Symbol resolution: cursors ----------
		{
			name: "OPEN/FETCH/CLOSE declared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE c CURSOR FOR SELECT 1;
    OPEN c;
    FETCH c INTO @x;
    CLOSE c;
END`,
		},
		{
			name: "OPEN undeclared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    OPEN nope;
END`,
			wantErr: `undeclared cursor: nope`,
		},
		{
			name: "FETCH undeclared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    FETCH nope INTO @x;
END`,
			wantErr: `undeclared cursor: nope`,
		},
		{
			name: "CLOSE undeclared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    CLOSE nope;
END`,
			wantErr: `undeclared cursor: nope`,
		},

		// ---------- Symbol resolution: HANDLER condition name ----------
		{
			name: "HANDLER references declared CONDITION",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key SET @e = 1;
    SELECT 1;
END`,
		},
		{
			name: "HANDLER references undeclared condition name",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE EXIT HANDLER FOR nope SET @e = 1;
    SELECT 1;
END`,
			wantErr: `undeclared condition: nope`,
		},
		{
			name: "HANDLER for SQLSTATE literal does not need declaration",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE EXIT HANDLER FOR SQLSTATE '23000' SET @e = 1;
    SELECT 1;
END`,
		},

		// ---------- Duplicate detection ----------
		{
			name: "duplicate DECLARE VAR same scope",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT;
    DECLARE x INT;
    SELECT 1;
END`,
			wantErr: `duplicate variable declaration: x`,
		},
		{
			name: "DECLARE VAR shadowed in nested block is OK",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE x INT;
    BEGIN
        DECLARE x INT;
        SELECT 1;
    END;
END`,
		},
		{
			// Shadowing a parameter in the body's BEGIN block is permitted —
			// the params live in the outer routine scope, the BEGIN body
			// opens a child scope, so the inner DECLARE doesn't conflict.
			name: "DECLARE VAR shadows parameter in nested BEGIN OK",
			sql: `CREATE PROCEDURE p(IN n INT)
BEGIN
    DECLARE n INT;
    SET n = 1;
END`,
		},
		{
			name: "duplicate DECLARE CURSOR same scope",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE c CURSOR FOR SELECT 1;
    DECLARE c CURSOR FOR SELECT 2;
    SELECT 1;
END`,
			wantErr: `duplicate cursor declaration: c`,
		},
		{
			name: "duplicate DECLARE CONDITION same scope",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE dk CONDITION FOR SQLSTATE '23000';
    DECLARE dk CONDITION FOR SQLSTATE '23001';
    SELECT 1;
END`,
			wantErr: `duplicate condition declaration: dk`,
		},
		{
			name: "duplicate label nested",
			sql: `CREATE PROCEDURE p()
lbl: BEGIN
    lbl: BEGIN
        SELECT 1;
    END lbl;
END lbl`,
			wantErr: `duplicate label: lbl`,
		},
		{
			name: "label name reused across handler barrier OK",
			sql: `CREATE PROCEDURE p()
my_block: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    my_block: BEGIN
        SELECT 1;
    END my_block;
    SELECT 2;
END my_block`,
		},
		{
			name: "duplicate HANDLER condition value within one DECLARE",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e = 1;
    SELECT 1;
END`,
			wantErr: `duplicate condition value`,
		},
		{
			name: "HANDLER different conditions same kind OK",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23001' SET @e = 1;
    SELECT 1;
END`,
		},
		{
			name: "HANDLER duplicate built-in category",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE CONTINUE HANDLER FOR NOT FOUND, NOT FOUND SET @d = 1;
    SELECT 1;
END`,
			wantErr: `duplicate condition value`,
		},

		// ---------- Function RETURN coverage ----------
		{
			name: "function with single RETURN body",
			sql:  `CREATE FUNCTION f() RETURNS INT RETURN 1`,
		},
		{
			name: "function with BEGIN ending in RETURN",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    SET x = x + 1;
    RETURN x;
END`,
		},
		{
			name: "function with IF ELSE both RETURN",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    IF x > 0 THEN RETURN 1;
    ELSE RETURN 0;
    END IF;
END`,
		},
		{
			name: "function missing RETURN on else path",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    IF x > 0 THEN RETURN 1;
    END IF;
END`,
			wantErr: `function body does not RETURN on all paths`,
		},
		{
			name: "function CASE with all branches RETURN + ELSE",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    CASE x
        WHEN 1 THEN RETURN 1;
        WHEN 2 THEN RETURN 2;
        ELSE RETURN 0;
    END CASE;
END`,
		},
		{
			name: "function CASE without ELSE rejected",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    CASE x
        WHEN 1 THEN RETURN 1;
    END CASE;
END`,
			wantErr: `function body does not RETURN on all paths`,
		},
		{
			name: "SIGNAL in function body counts as terminal",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    SIGNAL SQLSTATE '45000';
END`,
		},
		{
			name: "function whose body is just BEGIN ... SELECT (no RETURN)",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    SELECT 1;
END`,
			wantErr: `function body does not RETURN on all paths`,
		},
		{
			name: "function with HANDLER + RETURN at end",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET @done = 1;
    RETURN 0;
END`,
		},

		// ---------- RETURN outside function ----------
		{
			name: "RETURN inside procedure rejected",
			sql: `CREATE PROCEDURE p()
BEGIN
    RETURN 1;
END`,
			wantErr: `RETURN is only allowed inside a function body`,
		},
		{
			name: "RETURN inside trigger rejected",
			sql: `CREATE TRIGGER trg BEFORE INSERT ON tbl FOR EACH ROW
BEGIN
    RETURN 1;
END`,
			wantErr: `RETURN is only allowed inside a function body`,
		},
		{
			name: "RETURN inside event rejected",
			sql: `CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO
BEGIN
    RETURN 1;
END`,
			wantErr: `RETURN is only allowed inside a function body`,
		},
		{
			name: "RETURN nested inside function (in HANDLER body)",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION RETURN -1;
    RETURN 0;
END`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("Parse failed: %v", err)
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
