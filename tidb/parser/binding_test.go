package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's SQL plan bindings, as v8.5.5 reads them.
//
// The corpus in testdata pins the forms pingcap's own parser test writes.
// It writes almost no rejection, and a binding carries two whole
// statements, so the cases that matter most are the ones it leaves out: a
// statement kind TiDB will not bind, a missing FOR or USING, a malformed
// inner statement, a scope word where none belongs, and the digest lists,
// whose shape differs between the three statements that take one. This
// file carries them.
//
// Every case was measured against TiDB v8.5.5 and the expectation is what
// that server answered; TestBindingFormsMatchTiDB holds this file to it,
// so an image upgrade that changes a verdict fails there instead of
// passing silently.
//
// Ref: TiDB's SQL Plan Management,
// https://docs.pingcap.com/tidb/v8.5/sql-plan-management/.

// bindingCase is a statement and whether the engine's parser accepts it.
type bindingCase struct {
	sql    string
	accept bool
}

var bindingCases = []bindingCase{
	// --- CREATE BINDING FOR ... USING ..., and the scope words.
	{"CREATE BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)", true},
	{"CREATE GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)", true},
	{"CREATE SESSION BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)", true},
	{"create global binding for select * from t using select * from t use index(a)", true},
	// GLOBAL and SESSION come before BINDING and nowhere else, one at most.
	{"CREATE BINDING GLOBAL FOR SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE GLOBAL SESSION BINDING FOR SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE SESSION GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE LOCAL BINDING FOR SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE BINDING FOR GLOBAL SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE BINDING FOR SELECT * FROM t USING GLOBAL SELECT * FROM t", false},
	{"DROP BINDING GLOBAL FOR SELECT * FROM t", false},
	{"SET GLOBAL BINDING ENABLED FOR SELECT * FROM t", false},
	{"SET SESSION BINDING ENABLED FOR SELECT * FROM t", false},

	// --- FOR and USING are both required, and neither repeats.
	{"CREATE BINDING FOR SELECT * FROM t", false},
	{"CREATE BINDING FOR SELECT * FROM t USING", false},
	{"CREATE BINDING SELECT * FROM t USING SELECT * FROM t", false},
	{"CREATE BINDING USING", false},
	{"CREATE BINDING FOR SELECT * FROM t USING SELECT * FROM t USING SELECT * FROM t", false},
	{"DROP BINDING FOR SELECT * FROM t USING SELECT * FROM t USING SELECT * FROM t", false},
	{"DROP BINDING FOR", false},
	{"SET BINDING ENABLED", false},
	{"SET BINDING ENABLED SELECT * FROM t", false},
	{"SET BINDING ENABLED FOR USING SELECT * FROM t", false},
	{"SET BINDING FOR SELECT * FROM t", false},
	// ENABLED or DISABLED, spelled exactly.
	{"SET BINDING ENABLED FOR SELECT * FROM t", true},
	{"SET BINDING DISABLED FOR SELECT * FROM t", true},
	{"SET BINDING ENABLE FOR SELECT * FROM t", false},
	{"SET BINDING DISABLE FOR SELECT * FROM t", false},

	// --- CREATE BINDING USING stmt, which reads the shape off the plan.
	{"CREATE BINDING USING SELECT * FROM t", true},
	{"CREATE GLOBAL BINDING USING SELECT * FROM t", true},
	{"CREATE SESSION BINDING USING SELECT * FROM t", true},
	{"CREATE BINDING USING SELECT /*+ USE_INDEX(t, a) */ * FROM t", true},
	// DROP and SET take no such form: USING there needs a FOR before it.
	{"DROP BINDING USING SELECT * FROM t", false},
	{"SET BINDING ENABLED USING SELECT * FROM t", false},

	// --- The statement kinds a binding takes.
	{"CREATE BINDING FOR SELECT 1 USING SELECT 1", true},
	{"CREATE BINDING FOR (SELECT * FROM t) USING (SELECT * FROM t)", true},
	{"CREATE BINDING FOR SELECT 1 UNION SELECT 2 INTERSECT SELECT 3 USING SELECT 1 UNION SELECT 2 INTERSECT SELECT 3", true},
	{"CREATE BINDING FOR SELECT * FROM t FOR UPDATE USING SELECT * FROM t FOR UPDATE", true},
	{"CREATE BINDING FOR TABLE t USING TABLE t", true},
	{"CREATE BINDING FOR VALUES ROW(1, 2) USING VALUES ROW(1, 2)", true},
	{"CREATE BINDING FOR UPDATE t SET a = 1 USING UPDATE t SET a = 1", true},
	{"CREATE BINDING FOR INSERT INTO t1 SELECT * FROM t2 USING INSERT INTO t1 SELECT * FROM t2", true},
	{"CREATE BINDING FOR REPLACE INTO t1 SELECT * FROM t2 USING REPLACE INTO t1 SELECT * FROM t2", true},
	{"CREATE BINDING FOR DELETE FROM t WHERE a = 1 USING DELETE FROM t WHERE a = 1", true},
	// A bindable DELETE with no WHERE clause: the USING that follows it is
	// the binding's, and reading the DELETE must leave that token alone
	// rather than consume it and drop it.
	{"CREATE BINDING FOR DELETE FROM t USING DELETE FROM t", true},
	{"DROP BINDING FOR DELETE FROM t USING DELETE FROM t", true},
	{"SET BINDING ENABLED FOR DELETE FROM t USING DELETE FROM t", true},
	{"CREATE BINDING FOR DELETE t1, t2 FROM t1, t2 USING DELETE t1, t2 FROM t1, t2", true},
	{"CREATE BINDING FOR DELETE FROM t WHERE a = 1 ORDER BY a LIMIT 1 USING DELETE FROM t WHERE a = 1 ORDER BY a LIMIT 1", true},
	{"CREATE BINDING FOR DELETE FROM t1 AS x WHERE x.a = 1 USING DELETE FROM t1 AS x WHERE x.a = 1", true},
	{"CREATE BINDING FOR DELETE FROM t1 PARTITION (p0) WHERE a = 1 USING DELETE FROM t1 PARTITION (p0) WHERE a = 1", true},
	// Nothing else is a bindable statement.
	{"CREATE BINDING FOR SHOW TABLES USING SHOW TABLES", false},
	{"CREATE BINDING FOR EXPLAIN SELECT * FROM t USING EXPLAIN SELECT * FROM t", false},
	{"CREATE BINDING FOR CREATE TABLE zz (a INT) USING CREATE TABLE zz (a INT)", false},
	{"CREATE BINDING FOR SET @a = 1 USING SET @a = 1", false},
	{"CREATE BINDING FOR ANALYZE TABLE t USING ANALYZE TABLE t", false},
	{"CREATE BINDING FOR LOAD DATA INFILE 'a' INTO TABLE t USING LOAD DATA INFILE 'a' INTO TABLE t", false},
	{"CREATE BINDING FOR DO 1 USING DO 1", false},
	{"CREATE BINDING FOR USE test USING USE test", false},
	{"DROP BINDING FOR SHOW TABLES", false},
	{"SET BINDING ENABLED FOR SHOW TABLES", false},
	// A malformed inner statement fails where it is malformed, in either
	// slot: the parser reads both as statements rather than as text.
	{"CREATE BINDING FOR SELECT FROM t USING SELECT * FROM t", false},
	{"CREATE BINDING FOR SELECT * FROM t USING SELECT FROM t", false},
	{"CREATE BINDING FOR SELECT * FROM USING SELECT * FROM t", false},
	{"CREATE BINDING USING UPDATE t SET USING", false},
	{"DROP BINDING FOR UPDATE t SET", false},
	{"SET BINDING ENABLED FOR DELETE FROM WHERE a = 1", false},

	// --- The DELETE forms with a USING clause are not bindable: in
	// `DELETE FROM t1 USING t1, t2` the USING is read as the binding's, so
	// what follows must be a statement, and `t1, t2` is not one. The
	// FROM-first form then names one table, as it does on its own.
	{"CREATE BINDING FOR DELETE FROM t1 USING t1, t2 USING DELETE FROM t1 USING t1, t2", false},
	{"DROP BINDING FOR DELETE FROM t1 USING t1, t2", false},
	{"CREATE BINDING FOR SELECT * FROM t USING DELETE FROM t1 USING t1, t2", false},
	{"SET BINDING ENABLED FOR DELETE FROM t1 USING t1, t2", false},
	{"CREATE BINDING FOR DELETE FROM t1, t2 USING DELETE FROM t1, t2", false},
	{"DROP BINDING FOR DELETE FROM t1, t2", false},

	// --- A WITH clause opens a bindable query and a bindable UPDATE, and
	// nothing else. TiDB accepts `WITH c AS (SELECT 1) DELETE FROM t` on
	// its own but not inside a binding, and has no WITH before INSERT or
	// REPLACE at all.
	{"CREATE BINDING FOR WITH c AS (SELECT 1) SELECT * FROM c USING WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"CREATE BINDING FOR WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM c WHERE n<3) SELECT * FROM c USING WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM c WHERE n<3) SELECT * FROM c", true},
	{"CREATE BINDING FOR WITH c AS (SELECT 1) UPDATE t SET a = 1 USING WITH c AS (SELECT 1) UPDATE t SET a = 1", true},
	{"DROP BINDING FOR WITH c AS (SELECT 1) UPDATE t SET a = 1", true},
	{"CREATE BINDING USING WITH c AS (SELECT 1) UPDATE t SET a = 1", true},
	{"CREATE BINDING FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1 USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1", false},
	{"CREATE BINDING USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1", false},
	{"DROP BINDING FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1", false},
	{"SET BINDING ENABLED FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1", false},
	{"CREATE BINDING FOR WITH c AS (SELECT 1) INSERT INTO t1 SELECT * FROM c USING WITH c AS (SELECT 1) INSERT INTO t1 SELECT * FROM c", false},
	{"CREATE BINDING USING WITH c AS (SELECT 1) REPLACE INTO t1 SELECT * FROM c", false},

	// --- DROP BINDING's optional USING.
	{"DROP BINDING FOR SELECT * FROM t", true},
	{"DROP GLOBAL BINDING FOR SELECT * FROM t", true},
	{"DROP SESSION BINDING FOR SELECT * FROM t", true},
	{"DROP BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)", true},
	{"SET BINDING ENABLED FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)", true},

	// --- CREATE BINDING FROM HISTORY USING PLAN DIGEST, whose list takes a
	// string literal or a user variable, one at least, with no trailing
	// comma.
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 'sss'", true},
	{"CREATE GLOBAL BINDING FROM HISTORY USING PLAN DIGEST 'sss'", true},
	{"CREATE SESSION BINDING FROM HISTORY USING PLAN DIGEST 'sss'", true},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a', 'b'", true},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST @a", true},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST @a, @b, 'test1,test2', @c, 'test333'", true},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a',", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 1", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 1+1", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST x", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST @@global.x", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST CONCAT('a', 'b')", false},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a' USING SELECT * FROM t", false},
	{"CREATE BINDING FROM HISTORY USING SQL DIGEST 'a'", false},
	{"CREATE BINDING FROM HISTORY USING DIGEST 'a'", false},
	{"CREATE BINDING FROM HISTORY PLAN DIGEST 'a'", false},
	{"CREATE BINDING FOR HISTORY USING PLAN DIGEST 'a'", false},
	{"CREATE BINDING FROM USING PLAN DIGEST 'a'", false},

	// --- DROP BINDING FOR SQL DIGEST, the same list shape.
	{"DROP BINDING FOR SQL DIGEST 'a'", true},
	{"DROP GLOBAL BINDING FOR SQL DIGEST 's'", true},
	{"DROP SESSION BINDING FOR SQL DIGEST 'a'", true},
	{"DROP BINDING FOR SQL DIGEST 'a', 'b'", true},
	{"DROP BINDING FOR SQL DIGEST @a, @b, 'test1,test2', @c, 'test333'", true},
	{"DROP BINDING FOR SQL DIGEST @`a`", true},
	{"CREATE BINDING FROM HISTORY USING PLAN DIGEST @`a`", true},
	{"DROP BINDING FOR SQL DIGEST", false},
	{"DROP BINDING FOR SQL DIGEST 1", false},
	{"DROP BINDING FOR SQL DIGEST a", false},
	{"DROP BINDING FOR SQL DIGEST @@global.x", false},
	{"DROP BINDING FOR SQL DIGEST CONCAT('a', 'b')", false},
	{"DROP BINDING FOR SQL DIGEST 'a' USING SELECT * FROM t", false},
	{"DROP BINDING FOR SQL 'a'", false},
	{"DROP BINDING FOR DIGEST 'a'", false},
	{"DROP BINDING FOR PLAN DIGEST 'a'", false},

	// --- SET BINDING ... FOR SQL DIGEST takes one digest, and a string
	// literal only, where DROP BINDING takes a list and a user variable.
	{"SET BINDING ENABLED FOR SQL DIGEST '1'", true},
	{"SET BINDING DISABLED FOR SQL DIGEST '1'", true},
	{"SET BINDING ENABLED FOR SQL DIGEST 'a', 'b'", false},
	{"SET BINDING ENABLED FOR SQL DIGEST @a", false},
	{"SET BINDING DISABLED FOR SQL DIGEST @a", false},
	{"SET BINDING ENABLED FOR SQL DIGEST @a, @b", false},
	{"SET BINDING ENABLED FOR SQL DIGEST", false},
	{"SET BINDING ENABLED FOR PLAN DIGEST 'a'", false},

	// --- The cross-database wildcard, which reaches a table name through
	// a binding but is spelled in TiDB's table-name rule, so it parses
	// wherever a table name does. Only one name may follow the dot.
	{"CREATE BINDING USING SELECT * FROM *.t1", true},
	{"CREATE BINDING USING SELECT * FROM *.t1 WHERE t1.a > (SELECT MAX(a) FROM t2)", true},
	{"CREATE BINDING USING UPDATE *.t1 SET a = 1", true},
	{"SELECT * FROM *.t1", true},
	{"SELECT * FROM * . t1", true},
	{"SELECT * FROM `*`.t1", true},
	{"SELECT a FROM *.t1 JOIN *.t2 ON t1.a = t2.a", true},
	{"SELECT * FROM *.t1 AS x", true},
	{"UPDATE *.t1 SET a = 1", true},
	{"DELETE FROM *.t1", true},
	{"INSERT INTO *.t1 VALUES (1, 2)", true},
	{"REPLACE INTO *.t1 VALUES (1, 2)", true},
	{"CREATE TABLE *.zz (a INT)", true},
	{"DROP TABLE *.zz", true},
	{"ALTER TABLE *.t ADD COLUMN zz INT", true},
	{"TRUNCATE TABLE *.t1", true},
	{"RENAME TABLE *.t1 TO *.t2", true},
	{"ANALYZE TABLE *.t1", true},
	{"CREATE INDEX zz ON *.t1 (a)", true},
	{"SHOW CREATE TABLE *.t1", true},
	{"SELECT * FROM *.*", false},
	{"SELECT * FROM *.t1.a", false},
	{"SELECT *.t1.a FROM *.t1", false},
	{"SELECT * FROM test.*", false},
	// GRANT reads its priv_level with a rule of its own, which this change
	// does not touch; `GRANT SELECT ON *.t1 TO u1` is a 1064 on v8.5.5 and
	// this parser still accepts it, in every MySQL-family engine here.
	// That over-acceptance belongs to the GRANT family.
	// A `*` that no dot follows is still multiplication.
	{"SELECT .78*.21", true},
	{"SELECT 1 * .5", true},

	// --- The five words this grammar adds to the lexer are unreserved, so
	// each still names a table, a column, and an alias.
	{"CREATE TABLE kw (binding INT, digest INT, plan INT, enabled INT, disabled INT)", true},
	{"SELECT binding, digest, plan, enabled, disabled FROM kw", true},
	{"SELECT 1 AS binding, 2 AS digest, 3 AS plan, 4 AS enabled, 5 AS disabled", true},
	{"SELECT * FROM kw binding", true},
	{"SELECT * FROM kw AS plan", true},
	{"DELETE FROM kw binding WHERE binding.a = 1", true},
	{"DELETE FROM kw plan WHERE plan.a = 1", true},
	{"UPDATE kw binding SET binding.binding = 1", true},
	{"SET binding = 1", true},
	{"SET plan = 1", true},
	{"SET enabled = 1", true},
	{"SET @binding = 1", true},
	{"SAVEPOINT binding", true},
	{"CREATE INDEX binding ON kw (binding)", true},
	{"CREATE DATABASE IF NOT EXISTS binding", true},
}

// TestBindingForms holds this parser to the verdicts the table records.
func TestBindingForms(t *testing.T) {
	for _, c := range bindingCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, this file records %s (err=%v)",
					verdictName(accepts), c.sql, verdictName(c.accept), err)
			}
		})
	}
}

// parseOneBinding parses sql and returns its single statement.
func parseOneBinding(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if list.Len() != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, list.Len())
	}
	return list.Items[0]
}

// stmtText is the source text a statement node covers, which is how a tree
// test says which of two inner statements it is looking at.
func stmtText(t *testing.T, sql string, n nodes.Node) string {
	t.Helper()
	if n == nil {
		return ""
	}
	loc := bindingLoc(t, n)
	if loc.Start < 0 || loc.End > len(sql) || loc.End <= loc.Start {
		t.Fatalf("node %T has Loc %+v, which is no range of %q", n, loc, sql)
	}
	return sql[loc.Start:loc.End]
}

func bindingLoc(t *testing.T, n nodes.Node) nodes.Loc {
	t.Helper()
	switch v := n.(type) {
	case *nodes.SelectStmt:
		return v.Loc
	case *nodes.InsertStmt:
		return v.Loc
	case *nodes.UpdateStmt:
		return v.Loc
	case *nodes.DeleteStmt:
		return v.Loc
	case *nodes.TableStmt:
		return v.Loc
	case *nodes.ValuesStmt:
		return v.Loc
	}
	t.Fatalf("node %T is not a statement a binding carries", n)
	return nodes.NoLoc()
}

// TestBindingTrees pins what the parser builds, not only that it parsed:
// the node kind, the scope, and which inner statement landed in which
// slot. A binding whose two statements were swapped, or whose inner
// statement was kept as text, would parse just as well.
func TestBindingTrees(t *testing.T) {
	t.Run("create for using", func(t *testing.T) {
		sql := "CREATE GLOBAL BINDING FOR SELECT * FROM t WHERE a = 1 USING SELECT * FROM t USE INDEX(i) WHERE a = 1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.CreateBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built %T, want *ast.CreateBindingStmt", sql, parseOneBinding(t, sql))
		}
		if stmt.Scope != nodes.BindingScopeGlobal {
			t.Errorf("scope = %v, want GLOBAL", stmt.Scope)
		}
		if got, want := stmtText(t, sql, stmt.OriginStmt), "SELECT * FROM t WHERE a = 1"; got != want {
			t.Errorf("OriginStmt covers %q, want %q", got, want)
		}
		if got, want := stmtText(t, sql, stmt.HintedStmt), "SELECT * FROM t USE INDEX(i) WHERE a = 1"; got != want {
			t.Errorf("HintedStmt covers %q, want %q", got, want)
		}
		// Both slots hold a parsed statement, not text: the hinted one's
		// index hint is a node of its own.
		hinted, ok := stmt.HintedStmt.(*nodes.SelectStmt)
		if !ok {
			t.Fatalf("HintedStmt is %T, want *ast.SelectStmt", stmt.HintedStmt)
		}
		if len(hinted.From) != 1 {
			t.Fatalf("HintedStmt has %d FROM entries, want 1", len(hinted.From))
		}
		ref, ok := hinted.From[0].(*nodes.TableRef)
		if !ok {
			t.Fatalf("HintedStmt's FROM entry is %T, want *ast.TableRef", hinted.From[0])
		}
		if len(ref.IndexHints) != 1 || len(ref.IndexHints[0].Indexes) != 1 || ref.IndexHints[0].Indexes[0] != "i" {
			t.Errorf("HintedStmt's index hints = %+v, want USE INDEX(i)", ref.IndexHints)
		}
		if stmt.PlanDigests != nil {
			t.Errorf("PlanDigests = %v, want none", stmt.PlanDigests)
		}
	})

	t.Run("scope words", func(t *testing.T) {
		for _, c := range []struct {
			sql  string
			want nodes.BindingScope
		}{
			{"CREATE BINDING FOR SELECT 1 USING SELECT 1", nodes.BindingScopeUnset},
			{"CREATE GLOBAL BINDING FOR SELECT 1 USING SELECT 1", nodes.BindingScopeGlobal},
			{"CREATE SESSION BINDING FOR SELECT 1 USING SELECT 1", nodes.BindingScopeSession},
			{"create session binding using select 1", nodes.BindingScopeSession},
			{"CREATE SESSION BINDING FROM HISTORY USING PLAN DIGEST 'a'", nodes.BindingScopeSession},
		} {
			stmt, ok := parseOneBinding(t, c.sql).(*nodes.CreateBindingStmt)
			if !ok {
				t.Fatalf("Parse(%q) built no *ast.CreateBindingStmt", c.sql)
			}
			if stmt.Scope != c.want {
				t.Errorf("Parse(%q): scope = %v, want %v", c.sql, stmt.Scope, c.want)
			}
		}
		for _, c := range []struct {
			sql  string
			want nodes.BindingScope
		}{
			{"DROP BINDING FOR SELECT 1", nodes.BindingScopeUnset},
			{"DROP GLOBAL BINDING FOR SELECT 1", nodes.BindingScopeGlobal},
			{"DROP SESSION BINDING FOR SQL DIGEST 'a'", nodes.BindingScopeSession},
		} {
			stmt, ok := parseOneBinding(t, c.sql).(*nodes.DropBindingStmt)
			if !ok {
				t.Fatalf("Parse(%q) built no *ast.DropBindingStmt", c.sql)
			}
			if stmt.Scope != c.want {
				t.Errorf("Parse(%q): scope = %v, want %v", c.sql, stmt.Scope, c.want)
			}
		}
	})

	t.Run("create using only", func(t *testing.T) {
		sql := "CREATE BINDING USING SELECT * FROM t USE INDEX(a)"
		stmt, ok := parseOneBinding(t, sql).(*nodes.CreateBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.CreateBindingStmt", sql)
		}
		// The form has no FOR statement, so the FOR slot stays empty
		// rather than holding a copy of the plan.
		if stmt.OriginStmt != nil {
			t.Errorf("OriginStmt = %T, want none", stmt.OriginStmt)
		}
		if got, want := stmtText(t, sql, stmt.HintedStmt), "SELECT * FROM t USE INDEX(a)"; got != want {
			t.Errorf("HintedStmt covers %q, want %q", got, want)
		}
	})

	t.Run("create from history", func(t *testing.T) {
		sql := "CREATE SESSION BINDING FROM HISTORY USING PLAN DIGEST 'a', @b"
		stmt, ok := parseOneBinding(t, sql).(*nodes.CreateBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.CreateBindingStmt", sql)
		}
		if stmt.OriginStmt != nil || stmt.HintedStmt != nil {
			t.Errorf("the FROM HISTORY form carries statements: origin %T, hinted %T", stmt.OriginStmt, stmt.HintedStmt)
		}
		if len(stmt.PlanDigests) != 2 {
			t.Fatalf("PlanDigests has %d entries, want 2", len(stmt.PlanDigests))
		}
		lit, ok := stmt.PlanDigests[0].(*nodes.StringLit)
		if !ok || lit.Value != "a" {
			t.Errorf("PlanDigests[0] = %#v, want the string literal a", stmt.PlanDigests[0])
		}
		ref, ok := stmt.PlanDigests[1].(*nodes.VariableRef)
		if !ok || ref.Name != "b" || ref.System {
			t.Errorf("PlanDigests[1] = %#v, want the user variable @b", stmt.PlanDigests[1])
		}
	})

	t.Run("drop for using", func(t *testing.T) {
		sql := "DROP SESSION BINDING FOR UPDATE t SET a = 1 USING UPDATE /*+ INL_JOIN(t) */ t SET a = 1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.DropBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.DropBindingStmt", sql)
		}
		if stmt.Scope != nodes.BindingScopeSession {
			t.Errorf("scope = %v, want SESSION", stmt.Scope)
		}
		if got, want := stmtText(t, sql, stmt.OriginStmt), "UPDATE t SET a = 1"; got != want {
			t.Errorf("OriginStmt covers %q, want %q", got, want)
		}
		if got, want := stmtText(t, sql, stmt.HintedStmt), "UPDATE /*+ INL_JOIN(t) */ t SET a = 1"; got != want {
			t.Errorf("HintedStmt covers %q, want %q", got, want)
		}
		if _, ok := stmt.OriginStmt.(*nodes.UpdateStmt); !ok {
			t.Errorf("OriginStmt is %T, want *ast.UpdateStmt", stmt.OriginStmt)
		}
	})

	t.Run("drop without using", func(t *testing.T) {
		sql := "DROP BINDING FOR DELETE FROM t WHERE a = 1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.DropBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.DropBindingStmt", sql)
		}
		if got, want := stmtText(t, sql, stmt.OriginStmt), "DELETE FROM t WHERE a = 1"; got != want {
			t.Errorf("OriginStmt covers %q, want %q", got, want)
		}
		if stmt.HintedStmt != nil {
			t.Errorf("HintedStmt = %T, want none", stmt.HintedStmt)
		}
		if stmt.SQLDigests != nil {
			t.Errorf("SQLDigests = %v, want none", stmt.SQLDigests)
		}
	})

	t.Run("drop for sql digest", func(t *testing.T) {
		sql := "DROP BINDING FOR SQL DIGEST @a, 'b'"
		stmt, ok := parseOneBinding(t, sql).(*nodes.DropBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.DropBindingStmt", sql)
		}
		if stmt.OriginStmt != nil || stmt.HintedStmt != nil {
			t.Errorf("the SQL DIGEST form carries statements: origin %T, hinted %T", stmt.OriginStmt, stmt.HintedStmt)
		}
		if len(stmt.SQLDigests) != 2 {
			t.Fatalf("SQLDigests has %d entries, want 2", len(stmt.SQLDigests))
		}
		if ref, ok := stmt.SQLDigests[0].(*nodes.VariableRef); !ok || ref.Name != "a" {
			t.Errorf("SQLDigests[0] = %#v, want the user variable @a", stmt.SQLDigests[0])
		}
		if lit, ok := stmt.SQLDigests[1].(*nodes.StringLit); !ok || lit.Value != "b" {
			t.Errorf("SQLDigests[1] = %#v, want the string literal b", stmt.SQLDigests[1])
		}
	})

	t.Run("set binding", func(t *testing.T) {
		sql := "SET BINDING ENABLED FOR SELECT * FROM t WHERE a = 1 USING SELECT * FROM t USE INDEX(i) WHERE a = 1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.SetBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.SetBindingStmt", sql)
		}
		if !stmt.Enabled {
			t.Errorf("Enabled = false, want true")
		}
		if got, want := stmtText(t, sql, stmt.OriginStmt), "SELECT * FROM t WHERE a = 1"; got != want {
			t.Errorf("OriginStmt covers %q, want %q", got, want)
		}
		if got, want := stmtText(t, sql, stmt.HintedStmt), "SELECT * FROM t USE INDEX(i) WHERE a = 1"; got != want {
			t.Errorf("HintedStmt covers %q, want %q", got, want)
		}

		sql = "SET BINDING DISABLED FOR SQL DIGEST 'abc'"
		disabled, ok := parseOneBinding(t, sql).(*nodes.SetBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.SetBindingStmt", sql)
		}
		if disabled.Enabled {
			t.Errorf("Enabled = true, want false")
		}
		if lit, ok := disabled.SQLDigest.(*nodes.StringLit); !ok || lit.Value != "abc" {
			t.Errorf("SQLDigest = %#v, want the string literal abc", disabled.SQLDigest)
		}
	})

	t.Run("set binding is not a variable assignment", func(t *testing.T) {
		// BINDING is unreserved, so `SET binding = 1` stays an ordinary
		// assignment and only the status word picks the binding statement.
		sql := "SET binding = 1"
		if stmt, ok := parseOneBinding(t, sql).(*nodes.SetStmt); !ok {
			t.Errorf("Parse(%q) built %T, want *ast.SetStmt", sql, parseOneBinding(t, sql))
		} else if len(stmt.Assignments) != 1 || stmt.Assignments[0].Column.Column != "binding" {
			t.Errorf("Parse(%q) assigned %+v, want binding = 1", sql, stmt.Assignments)
		}
	})

	t.Run("wildcard database", func(t *testing.T) {
		sql := "CREATE BINDING USING SELECT * FROM *.t1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.CreateBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.CreateBindingStmt", sql)
		}
		sel, ok := stmt.HintedStmt.(*nodes.SelectStmt)
		if !ok {
			t.Fatalf("HintedStmt is %T, want *ast.SelectStmt", stmt.HintedStmt)
		}
		ref, ok := sel.From[0].(*nodes.TableRef)
		if !ok {
			t.Fatalf("the FROM entry is %T, want *ast.TableRef", sel.From[0])
		}
		if ref.Schema != "*" || ref.Name != "t1" {
			t.Errorf("table reference = %q.%q, want \"*\".\"t1\"", ref.Schema, ref.Name)
		}
		if got, want := sql[ref.Loc.Start:ref.Loc.End], "*.t1"; got != want {
			t.Errorf("the reference covers %q, want %q", got, want)
		}
	})

	t.Run("bindable delete keeps no using", func(t *testing.T) {
		// The DELETE a binding takes has no USING clause of its own, so
		// the one in the text belongs to the binding.
		sql := "CREATE BINDING FOR DELETE FROM t WHERE a = 1 USING DELETE FROM t WHERE a = 1"
		stmt, ok := parseOneBinding(t, sql).(*nodes.CreateBindingStmt)
		if !ok {
			t.Fatalf("Parse(%q) built no *ast.CreateBindingStmt", sql)
		}
		origin, ok := stmt.OriginStmt.(*nodes.DeleteStmt)
		if !ok {
			t.Fatalf("OriginStmt is %T, want *ast.DeleteStmt", stmt.OriginStmt)
		}
		if origin.Using != nil {
			t.Errorf("OriginStmt has a USING clause: %+v", origin.Using)
		}
		if got, want := stmtText(t, sql, stmt.OriginStmt), "DELETE FROM t WHERE a = 1"; got != want {
			t.Errorf("OriginStmt covers %q, want %q", got, want)
		}
		if got, want := stmtText(t, sql, stmt.HintedStmt), "DELETE FROM t WHERE a = 1"; got != want {
			t.Errorf("HintedStmt covers %q, want %q", got, want)
		}
		// A plain DELETE still reads its own USING clause.
		plain := "DELETE FROM t1 USING t1, t2"
		del, ok := parseOneBinding(t, plain).(*nodes.DeleteStmt)
		if !ok {
			t.Fatalf("Parse(%q) built %T, want *ast.DeleteStmt", plain, parseOneBinding(t, plain))
		}
		if len(del.Using) != 2 {
			t.Errorf("Parse(%q): USING has %d entries, want 2", plain, len(del.Using))
		}
	})
}

// TestBindingKeywordsStayIdentifiers pins the hazard of registering a
// keyword: a word the lexer now returns as a keyword token must still
// reach every position that reads a name, or a table called plan and a
// column called digest stop parsing. Each of these is accepted by v8.5.5
// (TestBindingFormsMatchTiDB measures them there too).
func TestBindingKeywordsStayIdentifiers(t *testing.T) {
	for _, word := range bindingKeywords {
		for _, form := range bindingKeywordForms {
			sql := fillKeyword(form, word)
			t.Run(sql, func(t *testing.T) {
				if _, err := Parse(sql); err != nil {
					t.Errorf("Parse(%q): %v", sql, err)
				}
			})
		}
	}
}

// bindingKeywords are the words this grammar adds to the lexer.
var bindingKeywords = []string{"binding", "digest", "plan", "enabled", "disabled"}

// bindingKeywordForms are name positions, one %s for each place the word
// goes.
//
// A role name is not among them: v8.5.5 answers 1064 for `CREATE ROLE
// binding` and for `CREATE ROLE status` alike, so its role names take no
// unreserved keyword at all and say nothing about these five words. This
// parser accepts both, which belongs to the role family.
var bindingKeywordForms = []string{
	"CREATE TABLE %s (%s INT)",
	"SELECT %s FROM t",
	"SELECT t.%s FROM t",
	"SELECT 1 AS %s",
	"SELECT * FROM t %s",
	"SELECT * FROM t AS %s",
	"SELECT * FROM %s.t",
	"DELETE FROM t %s WHERE %s.a = 1",
	"UPDATE t %s SET %s.a = 1",
	"INSERT INTO %s (%s) VALUES (1)",
	"CREATE INDEX %s ON t (%s)",
	"DROP INDEX %s ON t",
	"ALTER TABLE t ADD COLUMN %s INT",
	"ALTER TABLE t RENAME COLUMN a TO %s",
	"CREATE DATABASE %s",
	"SAVEPOINT %s",
	"RELEASE SAVEPOINT %s",
	"SET %s = 1",
	"SET @%s = 1",
	"CREATE VIEW %s AS SELECT 1",
	"SELECT * FROM t ORDER BY %s",
	"SELECT * FROM t GROUP BY %s HAVING %s > 1",
	"ANALYZE TABLE %s",
	"TRUNCATE TABLE %s",
}

// fillKeyword substitutes word for every %s in form.
func fillKeyword(form, word string) string {
	out := ""
	for i := 0; i < len(form); i++ {
		if form[i] == '%' && i+1 < len(form) && form[i+1] == 's' {
			out += word
			i++
			continue
		}
		out += string(form[i])
	}
	return out
}
