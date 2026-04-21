package validate

import (
	"testing"

	parser "github.com/bytebase/omni/mysql/parser"
)

// TestSemanticValidationSQL mirrors the MySQL sp_pcontext / sp_head::check_return
// semantics via SQL-fed cases. Each case asserts that Validate() either emits
// the expected diagnostic code or reports none. Parsing must always succeed:
// after the Phase 5 parser/validator split, grammar errors stay in the parser
// and all semantic errors move here.
//
// This file was migrated from mysql/parser/semantic_test.go as part of the
// Phase 5 structural PR (docs/plans/2026-04-21-mysql-parser-validator-split.md).
func TestSemanticValidationSQL(t *testing.T) {
	type tc struct {
		name     string
		sql      string
		wantCode string // empty = expect no matching diagnostic
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
			wantCode: "undeclared_loop_label",
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
			wantCode: "undeclared_label",
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
			wantCode: "undeclared_label",
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
			wantCode: "undeclared_variable",
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
			wantCode: "undeclared_cursor",
		},
		{
			name: "FETCH undeclared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    FETCH nope INTO @x;
END`,
			wantCode: "undeclared_cursor",
		},
		{
			name: "CLOSE undeclared cursor",
			sql: `CREATE PROCEDURE p()
BEGIN
    CLOSE nope;
END`,
			wantCode: "undeclared_cursor",
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
			wantCode: "undeclared_condition",
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
			wantCode: "duplicate_variable",
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
			wantCode: "duplicate_cursor",
		},
		{
			name: "duplicate DECLARE CONDITION same scope",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE dk CONDITION FOR SQLSTATE '23000';
    DECLARE dk CONDITION FOR SQLSTATE '23001';
    SELECT 1;
END`,
			wantCode: "duplicate_condition",
		},
		{
			name: "duplicate label nested",
			sql: `CREATE PROCEDURE p()
lbl: BEGIN
    lbl: BEGIN
        SELECT 1;
    END lbl;
END lbl`,
			wantCode: "duplicate_label",
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
			wantCode: "duplicate_handler_condition",
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
			wantCode: "duplicate_handler_condition",
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
			// MySQL 8.0 CREATE-time check is presence-only ("No RETURN
			// found in FUNCTION"). Path analysis is deferred to runtime,
			// so a function that returns only on the THEN path is accepted
			// at CREATE.
			name: "function missing RETURN on else path is accepted (matches MySQL)",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    IF x > 0 THEN RETURN 1;
    END IF;
END`,
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
			// Same as above: CASE without ELSE is accepted at CREATE
			// because a RETURN exists somewhere in the body.
			name: "function CASE without ELSE accepted (matches MySQL)",
			sql: `CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    CASE x
        WHEN 1 THEN RETURN 1;
    END CASE;
END`,
		},
		{
			// SIGNAL/RESIGNAL alone does NOT satisfy MySQL's check
			// (sp_head's HAS_RETURN flag is set only by RETURN).
			name: "function with only SIGNAL is rejected (MySQL ERR 1320)",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    SIGNAL SQLSTATE '45000';
END`,
			wantCode: "function_missing_return",
		},
		{
			name: "function whose body is just BEGIN ... SELECT (no RETURN)",
			sql: `CREATE FUNCTION f() RETURNS INT
BEGIN
    SELECT 1;
END`,
			wantCode: "function_missing_return",
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
			wantCode: "return_outside_function",
		},
		{
			name: "RETURN inside trigger rejected",
			sql: `CREATE TRIGGER trg BEFORE INSERT ON tbl FOR EACH ROW
BEGIN
    RETURN 1;
END`,
			wantCode: "return_outside_function",
		},
		{
			name: "RETURN inside event rejected",
			sql: `CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO
BEGIN
    RETURN 1;
END`,
			wantCode: "return_outside_function",
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
			list, err := parser.Parse(c.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			diags := Validate(list, Options{})
			if c.wantCode == "" {
				for _, d := range diags {
					t.Errorf("unexpected diagnostic: %+v", d)
				}
				return
			}
			for _, d := range diags {
				if d.Code == c.wantCode {
					return
				}
			}
			t.Fatalf("expected diagnostic code %q, got %v", c.wantCode, diags)
		})
	}
}
