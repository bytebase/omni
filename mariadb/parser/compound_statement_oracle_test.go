package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// TestCompoundStatementFormsOracle holds omni to MariaDB 11.8.8 on compound_statement
// (BEGIN NOT ATOMIC, LOOP and FOR standing as statements of their own), on the
// declarations a block holds (ROW, TYPE OF, ROW TYPE OF, cursor parameters), and on what
// the grammar's actions refuse of a declaration while they read it: a name declared twice,
// TYPE OF an undeclared variable, a handler for an undeclared condition or for a condition
// another handler of the block catches.
func TestCompoundStatementFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"BEGIN NOT ATOMIC END", true},
		{"BEGIN NOT ATOMIC SELECT 1; END", true},
		{"BEGIN NOT ATOMIC DECLARE a INT DEFAULT 1; SET a = 2; END", true},
		{"BEGIN NOT ATOMIC BEGIN NOT ATOMIC END; END", true},
		{"BEGIN ATOMIC END", false},
		{"BEGIN NOT ATOMIC END a", false},
		{"BEGIN NOT ATOMIC COMMIT; DECLARE a INT; END", false},
		{"BEGIN NOT ATOMIC RETURN 1; END", false},
		{"EXPLAIN BEGIN NOT ATOMIC END", false},
		{"SET STATEMENT max_statement_time=1 FOR BEGIN NOT ATOMIC DECLARE a INT; END", true},
		{"CREATE PROCEDURE p() BEGIN NOT ATOMIC END", true},
		{"CREATE PROCEDURE p() x: BEGIN NOT ATOMIC END", false},
		{"CREATE FUNCTION f() RETURNS INT BEGIN NOT ATOMIC RETURN 1; END", true},
		{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN NOT ATOMIC DECLARE a INT; END", true},

		{"LOOP SELECT 1; END LOOP", true},
		{"LOOP END LOOP", false},
		{"a: LOOP COMMIT; END LOOP", false},
		{"BEGIN NOT ATOMIC WHILE 1 DO END WHILE; END", false},
		{"BEGIN NOT ATOMIC IF 1 THEN COMMIT; ELSE END IF; END", false},
		{"SET STATEMENT max_statement_time=1 FOR LOOP COMMIT; END LOOP", true},

		{"FOR a IN (SELECT 1) DO COMMIT; END FOR", true},
		{"FOR a IN REVERSE (SELECT 1) DO COMMIT; END FOR", true},
		{"FOR a IN (WITH x AS (SELECT 1) SELECT * FROM x) DO COMMIT; END FOR", true},
		{"FOR a IN (SELECT 1 UNION SELECT 2) DO COMMIT; END FOR", true},
		{"FOR a IN (SELECT 1) DO END FOR", false},
		{"FOR a IN (SELECT 1) DO COMMIT; END FOR a", false},
		{"FOR a IN (SELECT 1 INTO @x) DO COMMIT; END FOR", false},
		{"BEGIN NOT ATOMIC x: FOR a IN (SELECT 1) DO LEAVE x; END FOR x; END", true},
		{"FOR r IN (SELECT 1 AS x) DO BEGIN DECLARE y TYPE OF r; END; END FOR", true},
		{"CREATE PROCEDURE p() FOR a IN (SELECT 1) DO COMMIT; END FOR", true},
		{"CREATE FUNCTION f() RETURNS INT FOR r IN (SELECT 1) DO RETURN 1; END FOR", true},
		{"CREATE PROCEDURE p() DETERMINISTIC FOR r IN (SELECT 1) DO SELECT 1; END FOR", true},
		{"BEGIN NOT /* c */ ATOMIC IF 1 THEN SET @a=1; END IF; END", true},
		{"FOR $r IN (SELECT 1) DO COMMIT; END FOR", true},
		{"FOR résumé IN (SELECT 1) DO COMMIT; END FOR", true},

		// SET STATEMENT ... FOR a compound statement inside a routine keeps the routine's rules.
		{"CREATE FUNCTION f() RETURNS INT BEGIN SET STATEMENT max_statement_time=1 FOR LOOP RETURN 1; END LOOP; END", true},
		{"CREATE FUNCTION f() RETURNS INT BEGIN SET STATEMENT max_statement_time=1 FOR LOOP SELECT 1; END LOOP; RETURN 1; END", false},
		{"CREATE PROCEDURE p() BEGIN DECLARE a INT; SET STATEMENT max_statement_time=1 FOR BEGIN NOT ATOMIC DECLARE b TYPE OF a; END; END", true},

		{"BEGIN NOT ATOMIC DECLARE a, b INT; END", true},
		{"BEGIN NOT ATOMIC DECLARE a, b, A INT; END", false},
		{"BEGIN NOT ATOMIC DECLARE a INT; DECLARE a INT; END", false},
		{"BEGIN NOT ATOMIC DECLARE a INT; BEGIN DECLARE a INT; END; END", true},
		{"BEGIN NOT ATOMIC DECLARE a INT; DECLARE b TYPE OF A; END", true},
		{"BEGIN NOT ATOMIC DECLARE a TYPE OF a; END", false},
		{"BEGIN NOT ATOMIC BEGIN DECLARE a INT; END; BEGIN DECLARE b TYPE OF a; END; END", false},
		{"BEGIN NOT ATOMIC DECLARE b TYPE OF t.c DEFAULT 1; END", true},
		{"BEGIN NOT ATOMIC DECLARE b TYPE OF d.t.c; END", true},
		{"BEGIN NOT ATOMIC DECLARE a ROW TYPE OF t; END", true},
		{"BEGIN NOT ATOMIC DECLARE a ROW TYPE OF d.t.c; END", true},
		{"BEGIN NOT ATOMIC DECLARE a ROW(b INT, c VARCHAR(10) CHARACTER SET latin1) DEFAULT ROW(1, 'x'); END", true},
		{"BEGIN NOT ATOMIC DECLARE a ROW(b INT, B INT); END", false},
		{"BEGIN NOT ATOMIC DECLARE a ROW(); END", false},
		{"BEGIN NOT ATOMIC DECLARE a ROW(b ROW(c INT)); END", false},
		{"BEGIN NOT ATOMIC DECLARE a ROW(b INT NOT NULL); END", false},
		{"BEGIN NOT ATOMIC DECLARE `row` INT; DECLARE type INT; END", true},

		{"BEGIN NOT ATOMIC DECLARE a CONDITION FOR 1; DECLARE a CONDITION FOR 2; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CONDITION FOR 1; BEGIN DECLARE a CONDITION FOR 2; END; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CONDITION FOR SQLSTATE VALUE '01000'; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CONDITION FOR NOT FOUND; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CONDITION FOR b; END", false},
		{"BEGIN NOT ATOMIC DECLARE a INT; DECLARE a CONDITION FOR 1; END", true},

		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR a COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; DECLARE EXIT HANDLER FOR C COMMIT; END", true},
		{"BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; BEGIN DECLARE EXIT HANDLER FOR c COMMIT; END; END", true},
		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CONDITION FOR 1; END; DECLARE EXIT HANDLER FOR c COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1, 1 COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1062 COMMIT; DECLARE CONTINUE HANDLER FOR 01062 COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; DECLARE EXIT HANDLER FOR c, 1 COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE c CONDITION FOR SQLSTATE '01000'; DECLARE EXIT HANDLER FOR c, SQLSTATE '01000' COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR SQLSTATE '01000', SQLWARNING COMMIT; END", true},
		{"BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR NOT FOUND, NOT FOUND COMMIT; END", false},
		{"BEGIN NOT ATOMIC DECLARE c CONDITION FOR SQLSTATE '42000'; DECLARE EXIT HANDLER FOR c COMMIT; BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '42000' COMMIT; END; END", true},

		{"BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1; DECLARE a CURSOR FOR SELECT 1; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b INT, c INT) FOR SELECT b; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (IN b INT) FOR SELECT b; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (OUT b INT) FOR SELECT b; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b INT DEFAULT 1) FOR SELECT b; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR () FOR SELECT 1; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (x INT, X INT) FOR SELECT 1; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b ROW(x INT)) FOR SELECT 1; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b TYPE OF t.c) FOR SELECT 1; END", true},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b TYPE OF a) FOR SELECT 1; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (b ROW TYPE OF a.b.c) FOR SELECT 1; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR (x INT) FOR SELECT 1; DECLARE b TYPE OF x; END", false},
		{"BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1 INTO @x; END", false},

		{"CREATE PROCEDURE p(IN a INT, OUT b ROW(x INT), INOUT c TYPE OF t.c) BEGIN END", true},
		{"CREATE PROCEDURE p(z INT) BEGIN DECLARE a TYPE OF z; END", true},
		{"CREATE PROCEDURE p(z INT) BEGIN DECLARE z INT; END", true},
		{"CREATE PROCEDURE p(z INT, Z INT) BEGIN END", false},
		{"CREATE PROCEDURE p(a TYPE OF a) BEGIN END", false},
		{"CREATE PROCEDURE p(a ROW TYPE OF d.t.c) BEGIN END", false},
		{"CREATE FUNCTION f(a ROW TYPE OF t) RETURNS INT RETURN 1", true},
		{"CREATE FUNCTION f(z INT, z INT) RETURNS INT RETURN 1", false},
	} {
		if got := o.prepareParses(t, c.sql); got != c.accept {
			t.Fatalf("MariaDB prepareParses(%q)=%v, want %v", c.sql, got, c.accept)
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != c.accept {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
				omniAccepts, c.accept, c.sql, perr)
		}
	}
}

// TestCompoundStatementTree pins what a compound statement and MariaDB's declarations
// read into: a statement's Loc runs from its first token to its last, and the anchored and
// ROW types land in their own fields.
func TestCompoundStatementTree(t *testing.T) {
	sql := "BEGIN NOT ATOMIC DECLARE a ROW(b INT); DECLARE c TYPE OF a; " +
		"DECLARE d CURSOR (e INT) FOR SELECT e; END"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	block, ok := list.Items[0].(*nodes.BeginEndBlock)
	if !ok || !block.NotAtomic {
		t.Fatalf("got %T %+v, want a NOT ATOMIC BeginEndBlock", list.Items[0], list.Items[0])
	}
	if block.Loc.Start != 0 || block.Loc.End != len(sql) {
		t.Errorf("block Loc = %+v, want {0 %d}", block.Loc, len(sql))
	}
	row := block.Stmts[0].(*nodes.DeclareVarStmt)
	if row.TypeName != nil || len(row.RowFields) != 1 || row.RowFields[0].Name != "b" {
		t.Errorf("ROW declaration = %+v", row)
	}
	anchored := block.Stmts[1].(*nodes.DeclareVarStmt)
	if anchored.Anchor == nil || anchored.Anchor.Row || len(anchored.Anchor.Names) != 1 {
		t.Errorf("TYPE OF declaration = %+v", anchored)
	}
	if got := sql[anchored.Anchor.Loc.Start:anchored.Anchor.Loc.End]; got != "TYPE OF a" {
		t.Errorf("anchor text = %q", got)
	}
	cursor := block.Stmts[2].(*nodes.DeclareCursorStmt)
	if len(cursor.Params) != 1 || cursor.Params[0].Name != "e" {
		t.Errorf("cursor params = %+v", cursor.Params)
	}

	sql = "FOR r IN (SELECT 1) DO COMMIT; END FOR"
	list, err = Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	loop, ok := list.Items[0].(*nodes.ForStmt)
	if !ok || loop.Var != "r" || loop.Select == nil || len(loop.Stmts) != 1 {
		t.Fatalf("got %T %+v, want a ForStmt", list.Items[0], list.Items[0])
	}
	if loop.Loc.Start != 0 || loop.Loc.End != len(sql) {
		t.Errorf("FOR Loc = %+v, want {0 %d}", loop.Loc, len(sql))
	}
}

// TestSplitCompoundStatements pins where Split cuts around BEGIN NOT ATOMIC and FOR: a
// comment inside BEGIN NOT ATOMIC, and a FOR loop that is a routine's whole body after its
// RETURNS or its characteristics, keep their statements in one segment.
func TestSplitCompoundStatements(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"BEGIN NOT /* c */ ATOMIC IF 1 THEN SET @a=1; END IF; END; SELECT 2",
			[]string{"BEGIN NOT /* c */ ATOMIC IF 1 THEN SET @a=1; END IF; END", " SELECT 2"}},
		{"CREATE FUNCTION f() RETURNS INT FOR r IN (SELECT 1) DO RETURN 1; END FOR; SELECT 2",
			[]string{"CREATE FUNCTION f() RETURNS INT FOR r IN (SELECT 1) DO RETURN 1; END FOR", " SELECT 2"}},
		{"CREATE PROCEDURE p() DETERMINISTIC FOR `r` IN (SELECT 1) DO SELECT 1; END FOR; SELECT 2",
			[]string{"CREATE PROCEDURE p() DETERMINISTIC FOR `r` IN (SELECT 1) DO SELECT 1; END FOR", " SELECT 2"}},
		{"FOR $r IN (SELECT 1) DO COMMIT; END FOR; SELECT 2",
			[]string{"FOR $r IN (SELECT 1) DO COMMIT; END FOR", " SELECT 2"}},
		{"FOR résumé IN (SELECT 1) DO COMMIT; END FOR; SELECT 2",
			[]string{"FOR résumé IN (SELECT 1) DO COMMIT; END FOR", " SELECT 2"}},
		{"SELECT for$r IN (SELECT 1); SELECT 2",
			[]string{"SELECT for$r IN (SELECT 1)", " SELECT 2"}},
		{"SELECT forë IN (SELECT 1); SELECT 2",
			[]string{"SELECT forë IN (SELECT 1)", " SELECT 2"}},
		{"SELECT * FROM t FOR UPDATE; SET STATEMENT a=1 FOR SELECT 1; SELECT 2",
			[]string{"SELECT * FROM t FOR UPDATE", " SET STATEMENT a=1 FOR SELECT 1", " SELECT 2"}},
	} {
		var got []string
		for _, s := range Split(c.sql) {
			got = append(got, s.Text)
		}
		if strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("Split(%q) = %q, want %q", c.sql, got, c.want)
		}
	}
}
