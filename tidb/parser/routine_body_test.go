package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Compound routine bodies whose blocks close with a two-word terminator —
// END IF / END CASE / END WHILE / END REPEAT — must not truncate the body at
// the inner END. The body is captured as raw text by balancing every compound
// opener (BEGIN/IF/CASE/WHILE/LOOP/REPEAT) against END, mirroring Split.
//
// CREATE PROCEDURE is the TiDB-grounded case: TiDB v8.5.0 parses these bodies
// (it reports ERROR 8108 = parse-OK, plan-unsupported), while a garbage body is
// rejected (1064), so TiDB validates the body grammar. Verified on TiDB v8.5.0.
// (TiDB rejects CREATE FUNCTION/TRIGGER/EVENT and LOOP-in-procedure outright;
// omni's acceptance of those is a pre-existing over-acceptance, tracked
// separately — not exercised as TiDB-parity here.)

func TestRoutineCompoundBodyAccepted(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN CASE 1 WHEN 1 THEN SELECT 1; END CASE; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE x INT DEFAULT 3; WHILE x>0 DO SET x=x-1; END WHILE; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE x INT DEFAULT 0; REPEAT SET x=x+1; UNTIL x>5 END REPEAT; END`,
		`CREATE PROCEDURE p() BEGIN IF 1>0 THEN BEGIN DECLARE y INT DEFAULT 3; WHILE y>0 DO SET y=y-1; END WHILE; END; END IF; END`,
		`CREATE PROCEDURE p() lbl: BEGIN SELECT 1; END lbl`,
		`CREATE PROCEDURE p() BEGIN END`,
		`CREATE PROCEDURE p() BEGIN SELECT 1; END`,
		// EXISTS (subquery) is a search condition, not a DDL IF EXISTS modifier;
		// the IF must be counted so its END IF balances.
		`CREATE PROCEDURE p() BEGIN IF EXISTS (SELECT 1 FROM t) THEN DELETE FROM t; END IF; END`,
		`CREATE PROCEDURE p() BEGIN IF NOT EXISTS (SELECT 1 FROM t) THEN INSERT INTO t VALUES (1); END IF; END`,
		// A comment between EXISTS and '(' must not flip the IF to a DDL modifier.
		`CREATE PROCEDURE p() BEGIN IF EXISTS /*c*/ (SELECT 1 FROM t) THEN DELETE FROM t; END IF; END`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

// The captured body must span exactly the routine body — the inner END IF must
// not truncate it, and the outer END must not over-consume into a following
// statement. Parsing a routine followed by another statement pins both ends.
func TestRoutineCompoundBodyFullCapture(t *testing.T) {
	const sql = `CREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END; SELECT 999`
	const wantBody = `BEGIN IF 1>0 THEN SELECT 1; END IF; END`

	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Parse(%q) produced %d statements, want 2 (routine + trailing SELECT)", sql, len(list.Items))
	}
	proc, ok := list.Items[0].(*nodes.CreateFunctionStmt)
	if !ok {
		t.Fatalf("first statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
	}
	if !proc.IsProcedure {
		t.Errorf("IsProcedure = false, want true")
	}
	if proc.Body != wantBody {
		t.Errorf("captured body = %q, want %q", proc.Body, wantBody)
	}
}

// findCompoundBodyEnd is the body-capture balancer. These cases exercise it
// directly (decoupled from which statements TiDB accepts): every compound
// opener nets against its END, openers inside string literals or comments are
// ignored, and a top-level ';' terminates the body. The LOOP case verifies the
// balancer counts LOOP/END LOOP even though TiDB rejects LOOP-in-procedure
// (statement-level rejection is a separate concern); here we only assert the
// body end is found correctly.
func TestFindCompoundBodyEnd(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		wantBody string // sql[:findCompoundBodyEnd(sql,0)]
	}{
		{"end_if", `BEGIN IF 1>0 THEN SELECT 1; END IF; END`, `BEGIN IF 1>0 THEN SELECT 1; END IF; END`},
		{"end_case", `BEGIN CASE 1 WHEN 1 THEN SELECT 1; END CASE; END`, `BEGIN CASE 1 WHEN 1 THEN SELECT 1; END CASE; END`},
		{"end_while", `BEGIN WHILE x>0 DO SET x=x-1; END WHILE; END`, `BEGIN WHILE x>0 DO SET x=x-1; END WHILE; END`},
		{"end_repeat", `BEGIN REPEAT SET x=x+1; UNTIL x>5 END REPEAT; END`, `BEGIN REPEAT SET x=x+1; UNTIL x>5 END REPEAT; END`},
		{"end_loop_balancing_only", `BEGIN l: LOOP LEAVE l; END LOOP; END`, `BEGIN l: LOOP LEAVE l; END LOOP; END`},
		{"nested", `BEGIN IF 1>0 THEN BEGIN WHILE y>0 DO SET y=y-1; END WHILE; END; END IF; END`, `BEGIN IF 1>0 THEN BEGIN WHILE y>0 DO SET y=y-1; END WHILE; END; END IF; END`},
		{"string_literal_with_end_if", `BEGIN SELECT 'END IF'; END`, `BEGIN SELECT 'END IF'; END`},
		{"comment_with_end_while", `BEGIN SELECT 1; /* END WHILE */ END`, `BEGIN SELECT 1; /* END WHILE */ END`},
		{"empty", `BEGIN END`, `BEGIN END`},
		{"top_level_semicolon_terminates", `BEGIN SELECT 1; END; SELECT 2`, `BEGIN SELECT 1; END`},
		// IF EXISTS (subquery) is a compound IF — counted — even with a comment
		// between EXISTS and '(', so its END IF balances and the body spans whole.
		{"if_exists_comment_subquery", `BEGIN IF EXISTS /*c*/ (SELECT 1) THEN SELECT 1; END IF; END`, `BEGIN IF EXISTS /*c*/ (SELECT 1) THEN SELECT 1; END IF; END`},
		// A body-internal DDL IF EXISTS <ident> is NOT a compound opener, so the
		// scan stops at the body's top-level ';' instead of over-counting and
		// swallowing the following statement.
		{"if_exists_ddl_stops_at_semicolon", `BEGIN DROP TABLE IF EXISTS t; END; SELECT 2`, `BEGIN DROP TABLE IF EXISTS t; END`},
		// IF(...) / REPEAT(...) functions with a comment before '(' are not
		// compound openers, so the scan stops at the body's top-level ';'.
		{"if_function_comment_before_paren", `BEGIN SELECT IF /*c*/ (1,2,3); END; SELECT 4`, `BEGIN SELECT IF /*c*/ (1,2,3); END`},
		{"repeat_function_comment_before_paren", `BEGIN SELECT REPEAT /*c*/ ('a',2); END; SELECT 4`, `BEGIN SELECT REPEAT /*c*/ ('a',2); END`},
		// Conditional comment before the function paren — a no-op version gate,
		// still a function call, so the scan stops at the body's top-level ';'.
		{"if_function_conditional_comment", `BEGIN SELECT IF /*!50000*/ (1,2,3); END; SELECT 4`, `BEGIN SELECT IF /*!50000*/ (1,2,3); END`},
		{"repeat_function_conditional_comment", `BEGIN SELECT REPEAT /*!50000*/ ('a',2); END; SELECT 4`, `BEGIN SELECT REPEAT /*!50000*/ ('a',2); END`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := findCompoundBodyEnd(c.sql, 0)
			if c.sql[:got] != c.wantBody {
				t.Errorf("findCompoundBodyEnd(%q) ended at %d (%q), want %q", c.sql, got, c.sql[:got], c.wantBody)
			}
		})
	}
}

// A routine body containing a MySQL conditional comment (/*! ... */ or
// /*T! ... */) must parse without panicking. The lexer rewrites its input
// buffer in place when it expands a conditional comment, so the body text must
// be captured against the post-tokenization buffer (clamped), not a stale
// pre-tokenization offset. TiDB v8.5.0 parses these bodies (ERROR 8108).
func TestRoutineCompoundBodyConditionalComment(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p() BEGIN /*!50003 SET @x = 1 */; END`,
		`CREATE PROCEDURE p() BEGIN /*T!90000 SET @x = 1 */; END`,
		`CREATE PROCEDURE p() BEGIN IF 1>0 THEN /*!50003 SET @x = 1 */; END IF; END`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) produced %d statements, want 1", sql, len(list.Items))
			}
			proc, ok := list.Items[0].(*nodes.CreateFunctionStmt)
			if !ok {
				t.Fatalf("statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
			}
			if proc.Body == "" {
				t.Errorf("captured body is empty, want the compound body text")
			}
		})
	}
}

// A routine defined with a custom DELIMITER (mysqldump style) followed by
// another statement must parse as two statements — the compound body must not
// swallow the following statement.
func TestRoutineBodyWithDelimiter(t *testing.T) {
	const sql = "DELIMITER //\nCREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END//\nSELECT 2//"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Parse produced %d statements, want 2 (procedure + SELECT)", len(list.Items))
	}
	if _, ok := list.Items[0].(*nodes.CreateFunctionStmt); !ok {
		t.Errorf("first statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
	}
}

// The IF [NOT] EXISTS tightening must not change how Split delimits ordinary
// DDL: IF EXISTS / IF NOT EXISTS followed by an identifier is a DDL modifier
// (not a compound opener), so these statements stay at top level and split on
// their own ';'.
func TestSplitPreservesIfExistsDDL(t *testing.T) {
	const sql = `DROP TABLE IF EXISTS t; CREATE TABLE IF NOT EXISTS u (a INT); SELECT 1`
	segs := Split(sql)
	if len(segs) != 3 {
		var texts []string
		for _, s := range segs {
			texts = append(texts, s.Text)
		}
		t.Fatalf("Split(%q) = %d segments %q, want 3", sql, len(segs), texts)
	}
}

// The body capture must not swallow a statement following the routine — the
// split-first guarantee the over-count analysis relies on. Covers the two body
// shapes whose scan could mis-balance: a conditional comment, and a
// body-internal DDL IF EXISTS, each followed by another statement.
func TestRoutineBodyDoesNotSwallowTrailingStatement(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		// Conditional-comment body + trailing SELECT (default ';' delimiter).
		// TiDB parses both (procedure 8108, SELECT executes).
		{"conditional_comment", `CREATE PROCEDURE p() BEGIN /*!50003 SET @x = 1 */; END; SELECT 2`},
		// DDL IF EXISTS in the body + trailing SELECT, custom delimiter so the
		// routine and SELECT share ONE segment (this is what exercises
		// findCompoundBodyEnd's top-level ';' detection). omni text-captures the
		// body; the point is that the trailing SELECT is not swallowed. TiDB
		// rejects DDL inside a procedure body, so this pins omni's split
		// correctness, not TiDB parity.
		{"ddl_if_exists_custom_delimiter", "DELIMITER //\nCREATE PROCEDURE p() BEGIN DROP TABLE IF EXISTS t; END; SELECT 2//"},
		// IF(...) / REPEAT(...) are built-in functions, not compound openers,
		// even with a comment between the keyword and '('. TiDB v8.5.0 accepts
		// SELECT IF /*c*/ (1,2,3) and SELECT REPEAT /*c*/ ('a',2).
		{"if_function_comment_before_paren", `CREATE PROCEDURE p() BEGIN SELECT IF /*c*/ (1,2,3); END; SELECT 4`},
		{"repeat_function_comment_before_paren", `CREATE PROCEDURE p() BEGIN SELECT REPEAT /*c*/ ('a',2); END; SELECT 4`},
		// A conditional comment (/*!...*/) is a no-op version gate before the
		// function paren — still a function call, not a compound opener.
		{"if_function_conditional_comment", `CREATE PROCEDURE p() BEGIN SELECT IF /*!50000*/ (1,2,3); END; SELECT 4`},
		{"repeat_function_conditional_comment", `CREATE PROCEDURE p() BEGIN SELECT REPEAT /*!50000*/ ('a',2); END; SELECT 4`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.sql, err)
			}
			if len(list.Items) != 2 {
				t.Fatalf("Parse(%q) produced %d statements, want 2 (routine + trailing)", c.sql, len(list.Items))
			}
			if _, ok := list.Items[0].(*nodes.CreateFunctionStmt); !ok {
				t.Errorf("first statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
			}
		})
	}
}
