package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Tests for the `parser-dml` node: INSERT / UPDATE / DELETE / MERGE / TRUNCATE.
//
// CORRECTNESS BASIS (correctness-protocol.md). The accept/reject verdicts are
// proven two ways, both authoritative for the forms they cover:
//   - the canonical ZetaSQL corpus (parser/googlesql/examples/.../dml_*.sql) —
//     the legacy ANTLR grammar bytebase consumes is a hand-port of that ZetaSQL
//     reference, so the corpus is the breadth oracle (TestDML_LegacyCorpusAccepts);
//   - the live Cloud Spanner emulator — the differential is in dml_oracle_test.go
//     (build tag `googlesql_oracle`), covering the SHARED + Spanner-only forms.
//
// BigQuery-only forms (MERGE, TRUNCATE, dashed table paths, the bare `INSERT …
// SELECT … ON CONFLICT` reject) are NON-authoritative on the Spanner emulator
// (MERGE: rejected at statement dispatch before body parse; TRUNCATE: Spanner
// DDL has none; dashed paths: Spanner syntax-rejects '-'). They are covered here
// against the ZetaSQL corpus + the BigQuery docs and recorded in the divergence
// ledger — NOT diffed against the emulator (a Spanner reject there would be a
// false divergence).

// ---------------------------------------------------------------------------
// INSERT — AST structure
// ---------------------------------------------------------------------------

func TestDML_Insert(t *testing.T) {
	t.Run("values with columns", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO Singers (SingerId, FirstName) VALUES (1, 'Marc')")
		if !ins.Into {
			t.Error("Into = false, want true")
		}
		if got := pathString(ins.Target); got != "Singers" {
			t.Errorf("Target = %q, want Singers", got)
		}
		if got := strings.Join(ins.Columns, ","); got != "SingerId,FirstName" {
			t.Errorf("Columns = %q", got)
		}
		if len(ins.Rows) != 1 || len(ins.Rows[0].Values) != 2 {
			t.Fatalf("Rows = %#v", ins.Rows)
		}
		if ins.OrAction != "" {
			t.Errorf("OrAction = %q, want empty", ins.OrAction)
		}
	})

	t.Run("no INTO, no columns", func(t *testing.T) {
		ins := parseInsert(t, "INSERT Singers VALUES (1, 'Marc')")
		if ins.Into {
			t.Error("Into = true, want false")
		}
		if ins.Columns != nil {
			t.Errorf("Columns = %v, want nil", ins.Columns)
		}
	})

	t.Run("multi-row values", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1), (2), (3)")
		if len(ins.Rows) != 3 {
			t.Fatalf("Rows = %d, want 3", len(ins.Rows))
		}
	})

	t.Run("DEFAULT value", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a, b) VALUES (1, DEFAULT)")
		row := ins.Rows[0]
		if _, ok := row.Values[1].(*ast.DefaultExpr); !ok {
			t.Errorf("Values[1] = %T, want *ast.DefaultExpr", row.Values[1])
		}
	})

	t.Run("insert select", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) SELECT x FROM s")
		if ins.Query == nil {
			t.Fatal("Query = nil")
		}
		if _, ok := ins.Query.(*ast.QueryStmt); !ok {
			t.Errorf("Query = %T, want *ast.QueryStmt", ins.Query)
		}
		if ins.Rows != nil {
			t.Error("Rows should be nil for a query source")
		}
		if ins.QueryParens {
			t.Error("QueryParens = true for a bare SELECT source")
		}
	})

	t.Run("insert with leading WITH cte", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO T (c1) WITH q1 AS (SELECT 1) SELECT * FROM q1")
		q := ins.Query.(*ast.QueryStmt)
		if q.With == nil || len(q.With.CTEs) != 1 {
			t.Errorf("expected one CTE, got %#v", q.With)
		}
	})

	t.Run("parenthesized query source", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) (SELECT 1)")
		if !ins.QueryParens {
			t.Error("QueryParens = false, want true for ( query )")
		}
		if ins.OnConflict != nil {
			t.Error("OnConflict should be nil")
		}
	})

	t.Run("set-op bare query source", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a, b) (SELECT 5, 7) UNION ALL SELECT 8, 9")
		if ins.QueryParens {
			t.Error("QueryParens = true; a set-op bare query is NOT the ( query ) form")
		}
		q := ins.Query.(*ast.QueryStmt)
		if _, ok := q.Body.(*ast.SetOperation); !ok {
			t.Errorf("Query.Body = %T, want *ast.SetOperation", q.Body)
		}
	})

	t.Run("parenthesized query with query-level tail is a bare query", func(t *testing.T) {
		// `( query ) LIMIT n` is alt-1 (a bare query), NOT the ON-CONFLICT-eligible
		// `( query )` form — so QueryParens stays false and the trailing LIMIT must
		// be consumed (oracle: accepts; with a trailing ON CONFLICT it rejects).
		for _, sql := range []string{
			"INSERT INTO t (a) (SELECT 1) LIMIT 5",
			"INSERT INTO t (a) (SELECT 1) ORDER BY x",
			"INSERT INTO t (a) (SELECT 1) ORDER BY x LIMIT 5",
		} {
			ins := parseInsert(t, sql)
			if ins.QueryParens {
				t.Errorf("%q: QueryParens = true; a query with a trailing ORDER BY/LIMIT is a bare query", sql)
			}
		}
	})
}

func TestDML_InsertOrAction(t *testing.T) {
	cases := map[string]string{
		"INSERT OR IGNORE INTO t (a) VALUES (1)":  "OR IGNORE",
		"INSERT OR UPDATE INTO t (a) VALUES (1)":  "OR UPDATE",
		"INSERT OR REPLACE INTO t (a) VALUES (1)": "OR REPLACE",
		"INSERT IGNORE INTO t (a) VALUES (1)":     "IGNORE",
		"INSERT REPLACE INTO t (a) VALUES (1)":    "REPLACE",
		"INSERT UPDATE INTO t (a) VALUES (1)":     "UPDATE",
	}
	for sql, want := range cases {
		t.Run(want, func(t *testing.T) {
			ins := parseInsert(t, sql)
			if ins.OrAction != want {
				t.Errorf("OrAction = %q, want %q", ins.OrAction, want)
			}
		})
	}
}

func TestDML_InsertOnConflict(t *testing.T) {
	t.Run("do nothing", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) ON CONFLICT DO NOTHING")
		if ins.OnConflict == nil || !ins.OnConflict.DoNothing {
			t.Fatalf("OnConflict = %#v", ins.OnConflict)
		}
	})
	t.Run("target columns + do update set where", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) ON CONFLICT (a) DO UPDATE SET a = a + excluded.a WHERE a > 1")
		oc := ins.OnConflict
		if oc == nil || oc.DoNothing {
			t.Fatalf("OnConflict = %#v", oc)
		}
		if strings.Join(oc.Columns, ",") != "a" {
			t.Errorf("Columns = %v", oc.Columns)
		}
		if len(oc.SetItems) != 1 {
			t.Errorf("SetItems = %#v", oc.SetItems)
		}
		if oc.Where == nil {
			t.Error("Where = nil")
		}
	})
	t.Run("on unique constraint", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) ON CONFLICT ON UNIQUE CONSTRAINT uc DO NOTHING")
		if ins.OnConflict.ConstraintName != "uc" {
			t.Errorf("ConstraintName = %q, want uc", ins.OnConflict.ConstraintName)
		}
	})
	t.Run("parenthesized query + on conflict", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) (SELECT 1) ON CONFLICT DO NOTHING")
		if !ins.QueryParens || ins.OnConflict == nil {
			t.Fatalf("QueryParens=%v OnConflict=%#v", ins.QueryParens, ins.OnConflict)
		}
	})
}

func TestDML_InsertTrailers(t *testing.T) {
	t.Run("assert rows modified", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED 5")
		if ins.AssertRows == nil {
			t.Error("AssertRows = nil")
		}
	})
	t.Run("then return", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) THEN RETURN a, b")
		if ins.Returning == nil || len(ins.Returning.Items) != 2 {
			t.Fatalf("Returning = %#v", ins.Returning)
		}
		if ins.Returning.WithAction {
			t.Error("WithAction = true, want false")
		}
	})
	t.Run("then return with action", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) THEN RETURN WITH ACTION AS act *")
		r := ins.Returning
		if r == nil || !r.WithAction || r.ActionAlias != "act" {
			t.Fatalf("Returning = %#v", r)
		}
		if len(r.Items) != 1 || !r.Items[0].Star {
			t.Errorf("Items = %#v, want a single star", r.Items)
		}
	})
	t.Run("then return with inline WITH expr item", func(t *testing.T) {
		// `THEN RETURN WITH(x AS 1, x)` is a select-list item that is an inline
		// WITH(...) expression — NOT the WITH ACTION modifier. Regression for the
		// over-eager WITH-ACTION branch (oracle: accepts).
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) THEN RETURN WITH(x AS 1, x)")
		if ins.Returning == nil || ins.Returning.WithAction {
			t.Fatalf("Returning = %#v; WITH(...) must NOT be read as WITH ACTION", ins.Returning)
		}
		if len(ins.Returning.Items) != 1 {
			t.Errorf("Items = %d, want 1 (the WITH expression)", len(ins.Returning.Items))
		}
	})
	t.Run("assert rows restricted operand", func(t *testing.T) {
		// ASSERT_ROWS_MODIFIED takes only an int literal / parameter / @@var / CAST
		// of one — NOT an arbitrary expression. Regression: these must reject
		// (oracle-confirmed syntax rejects).
		for _, sql := range []string{
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED 1 + 1",
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED 'x'",
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED (SELECT 1)",
		} {
			if _, errs := Parse(sql); len(errs) == 0 {
				t.Errorf("Parse(%q): expected a syntax error (operand is not int/param/CAST)", sql)
			}
		}
		// The accepted operand forms.
		for _, sql := range []string{
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED 5",
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED @p",
			"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED CAST(@p AS INT64)",
		} {
			parseInsert(t, sql)
		}
	})
	t.Run("keyword column name", func(t *testing.T) {
		// `graph` is a non-reserved keyword and a valid column name; the column-list
		// vs query-source discriminator must not treat it as a query head.
		ins := parseInsert(t, "INSERT INTO t (graph) VALUES (1)")
		if len(ins.Columns) != 1 {
			t.Errorf("Columns = %v, want one (graph)", ins.Columns)
		}
	})
	t.Run("malformed embedded subquery rejected", func(t *testing.T) {
		// A DML expression-embedded subquery is re-parsed (fillSubqueries), so a
		// malformed one surfaces a diagnostic. Regression (oracle: rejects on `b`).
		if _, errs := Parse("UPDATE t SET x = (SELECT 1 FROM s a b) WHERE id = 1"); len(errs) == 0 {
			t.Error("expected a syntax error for the malformed embedded subquery")
		}
	})
	t.Run("assert + returning together", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO t (a) VALUES (1) ON CONFLICT DO NOTHING ASSERT_ROWS_MODIFIED 1 THEN RETURN WITH ACTION AS ACTION a, b")
		if ins.OnConflict == nil || ins.AssertRows == nil || ins.Returning == nil {
			t.Fatalf("ins = %#v", ins)
		}
		if ins.Returning.ActionAlias != "ACTION" {
			t.Errorf("ActionAlias = %q (ACTION is a valid non-reserved alias)", ins.Returning.ActionAlias)
		}
	})
}

// ---------------------------------------------------------------------------
// UPDATE — AST structure
// ---------------------------------------------------------------------------

func TestDML_Update(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE Singers SET FirstName = 'A' WHERE SingerId = 1")
		if pathString(u.Target) != "Singers" {
			t.Errorf("Target = %q", pathString(u.Target))
		}
		if len(u.Items) != 1 {
			t.Fatalf("Items = %d", len(u.Items))
		}
		if u.Items[0].Path != "FirstName" {
			t.Errorf("Items[0].Path = %q", u.Items[0].Path)
		}
		if u.Where == nil {
			t.Error("Where = nil")
		}
	})
	t.Run("alias forms", func(t *testing.T) {
		for _, sql := range []string{
			"UPDATE t a SET a.x = 1 WHERE a.id = 2",
			"UPDATE t AS a SET x = 1 WHERE id = 2",
		} {
			u := parseUpdate(t, sql)
			if u.Alias != "a" {
				t.Errorf("%q: Alias = %q, want a", sql, u.Alias)
			}
		}
	})
	t.Run("multiple set items + DEFAULT rhs", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE t SET x = 1, y = DEFAULT WHERE id = 1")
		if len(u.Items) != 2 {
			t.Fatalf("Items = %d", len(u.Items))
		}
		if _, ok := u.Items[1].Value.(*ast.DefaultExpr); !ok {
			t.Errorf("Items[1].Value = %T, want *ast.DefaultExpr", u.Items[1].Value)
		}
	})
	t.Run("from clause join-update", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE t a SET x = b.y FROM s b WHERE a.id = b.id")
		if len(u.From) != 1 {
			t.Fatalf("From = %#v", u.From)
		}
	})
	t.Run("with offset alias", func(t *testing.T) {
		// Non-keyword alias `o` — see the DELETE with-offset note (divergence #84).
		u := parseUpdate(t, "UPDATE T WITH OFFSET AS o SET x = y")
		if !u.WithOffset || u.WithOffsetAlias != "o" {
			t.Errorf("WithOffset=%v alias=%q", u.WithOffset, u.WithOffsetAlias)
		}
	})
	t.Run("nested DML set item", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE T SET (DELETE x), z = DEFAULT")
		if len(u.Items) != 2 {
			t.Fatalf("Items = %d", len(u.Items))
		}
		if _, ok := u.Items[0].Nested.(*ast.DeleteStmt); !ok {
			t.Errorf("Items[0].Nested = %T, want *ast.DeleteStmt", u.Items[0].Nested)
		}
	})
	t.Run("generalized lhs paths", func(t *testing.T) {
		// These are the ZetaSQL update_set_value LHS forms.
		for _, sql := range []string{
			"UPDATE T SET id.(path.x.extension) = 5",
			"UPDATE T SET id[0] = DEFAULT",
			"UPDATE T SET id1.id2.(path.x.extension) = 5",
			"UPDATE T SET id1[0].(a.b.c).id1.(d.e.f)[1].id3 = 5",
			"UPDATE T SET id1[0][1] = 5",
		} {
			parseUpdate(t, sql)
		}
	})
	t.Run("assert rows + then return", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE t SET x = 1 WHERE id = 1 ASSERT_ROWS_MODIFIED 1 THEN RETURN x")
		if u.AssertRows == nil || u.Returning == nil {
			t.Fatalf("u = %#v", u)
		}
	})
}

// ---------------------------------------------------------------------------
// DELETE — AST structure
// ---------------------------------------------------------------------------

func TestDML_Delete(t *testing.T) {
	t.Run("with FROM", func(t *testing.T) {
		d := parseDelete(t, "DELETE FROM Singers WHERE SingerId = 1")
		if !d.From {
			t.Error("From = false, want true")
		}
		if d.Where == nil {
			t.Error("Where = nil")
		}
	})
	t.Run("without FROM", func(t *testing.T) {
		d := parseDelete(t, "DELETE Singers WHERE SingerId = 1")
		if d.From {
			t.Error("From = true, want false")
		}
	})
	t.Run("bare delete, no where", func(t *testing.T) {
		d := parseDelete(t, "DELETE T")
		if d.From || d.Where != nil {
			t.Errorf("d = %#v", d)
		}
	})
	t.Run("alias + then return", func(t *testing.T) {
		d := parseDelete(t, "DELETE Singers AS s WHERE s.SingerId = 1 THEN RETURN s.FirstName")
		if d.Alias != "s" || d.Returning == nil {
			t.Fatalf("d = %#v", d)
		}
	})
	t.Run("with offset", func(t *testing.T) {
		// Use a non-keyword alias `o`: the keyword-spelled `offset` alias is
		// affected by the known identifierText case-folding gap (divergence #84,
		// out of this node's scope), which would mask the WITH OFFSET assertion.
		d := parseDelete(t, "DELETE T WITH OFFSET AS o WHERE true")
		if !d.WithOffset || d.WithOffsetAlias != "o" {
			t.Errorf("WithOffset=%v alias=%q", d.WithOffset, d.WithOffsetAlias)
		}
	})
	t.Run("assert rows cast forms", func(t *testing.T) {
		for _, sql := range []string{
			"DELETE T ASSERT_ROWS_MODIFIED 10",
			"DELETE T ASSERT_ROWS_MODIFIED CAST(@param1 AS int64)",
			"DELETE T ASSERT_ROWS_MODIFIED CAST(@@sysvar AS int64)",
			"DELETE T ASSERT_ROWS_MODIFIED @row_count",
		} {
			d := parseDelete(t, sql)
			if d.AssertRows == nil {
				t.Errorf("%q: AssertRows = nil", sql)
			}
		}
	})
	t.Run("generalized target", func(t *testing.T) {
		parseDelete(t, "DELETE T.(a.b).c WHERE true")
		parseDelete(t, "DELETE T.a[0].b WHERE true")
	})
}

// ---------------------------------------------------------------------------
// MERGE — AST structure (BigQuery; ZetaSQL-corpus oracle)
// ---------------------------------------------------------------------------

func TestDML_Merge(t *testing.T) {
	t.Run("matched update", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 WHEN MATCHED AND T.T1 = 5 THEN UPDATE SET T1 = T1 + 10, T2 = T.T1")
		if pathString(m.Target) != "T" {
			t.Errorf("Target = %q", pathString(m.Target))
		}
		if len(m.Whens) != 1 {
			t.Fatalf("Whens = %d", len(m.Whens))
		}
		w := m.Whens[0]
		if !w.Matched || w.And == nil {
			t.Errorf("when = %#v", w)
		}
		if w.Action.Kind != ast.MergeUpdate || len(w.Action.SetItems) != 2 {
			t.Errorf("action = %#v", w.Action)
		}
	})
	t.Run("not matched by target insert", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 WHEN NOT MATCHED BY TARGET THEN INSERT (t1, t2) VALUES (10, S.C3)")
		w := m.Whens[0]
		if !w.NotMatched || !w.ByTarget {
			t.Errorf("when = %#v", w)
		}
		if w.Action.Kind != ast.MergeInsert || len(w.Action.Columns) != 2 {
			t.Errorf("action = %#v", w.Action)
		}
	})
	t.Run("not matched by source delete", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 WHEN NOT MATCHED BY SOURCE THEN DELETE")
		w := m.Whens[0]
		if !w.NotMatched || !w.BySource || w.Action.Kind != ast.MergeDelete {
			t.Errorf("when = %#v", w)
		}
	})
	t.Run("not matched default (no BY)", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 WHEN NOT MATCHED THEN INSERT (t1) VALUES (1)")
		w := m.Whens[0]
		if !w.NotMatched || w.ByTarget || w.BySource {
			t.Errorf("when = %#v (default NOT MATCHED is neither BY TARGET nor BY SOURCE)", w)
		}
	})
	t.Run("insert ROW source", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 WHEN NOT MATCHED BY SOURCE THEN INSERT ROW")
		a := m.Whens[0].Action
		if a.Kind != ast.MergeInsert || !a.SourceRow || a.InsertRow != nil {
			t.Errorf("action = %#v", a)
		}
	})
	t.Run("subquery source + alias", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T AS X USING (SELECT * FROM Y JOIN Z ON Y.C1 = Z.C1) AS S ON X.t1 = S.s1 WHEN MATCHED THEN DELETE")
		if m.Alias != "X" {
			t.Errorf("Alias = %q, want X", m.Alias)
		}
		src, ok := m.Source.(*ast.TableExpr)
		if !ok || src.Subquery == nil || src.Alias != "S" {
			t.Errorf("Source = %#v", m.Source)
		}
	})
	t.Run("multiple when clauses", func(t *testing.T) {
		m := parseMerge(t, "MERGE INTO T USING S ON t1 = s1 "+
			"WHEN MATCHED AND T.T1 = 5 THEN UPDATE SET T1 = T1 + 10 "+
			"WHEN NOT MATCHED BY TARGET THEN INSERT (t1) VALUES (10) "+
			"WHEN NOT MATCHED BY SOURCE THEN DELETE")
		if len(m.Whens) != 3 {
			t.Fatalf("Whens = %d, want 3", len(m.Whens))
		}
	})
}

// ---------------------------------------------------------------------------
// TRUNCATE — AST structure (BigQuery; ZetaSQL-corpus oracle)
// ---------------------------------------------------------------------------

func TestDML_Truncate(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		tr := parseTruncate(t, "TRUNCATE TABLE foo")
		if pathString(tr.Target) != "foo" || tr.Where != nil {
			t.Errorf("tr = %#v", tr)
		}
	})
	t.Run("with where", func(t *testing.T) {
		tr := parseTruncate(t, "TRUNCATE TABLE foo WHERE bar > 3")
		if tr.Where == nil {
			t.Error("Where = nil")
		}
	})
	t.Run("qualified path", func(t *testing.T) {
		// Use non-keyword path parts: `project`/`dataset` are non-reserved keywords
		// affected by the identifierText case-folding gap (divergence #84).
		tr := parseTruncate(t, "TRUNCATE TABLE myproj.ds.foo")
		if pathString(tr.Target) != "myproj.ds.foo" {
			t.Errorf("Target = %q", pathString(tr.Target))
		}
	})
}

// ---------------------------------------------------------------------------
// Dashed BigQuery paths (BigQuery-only; non-authoritative on Spanner — see ledger)
// ---------------------------------------------------------------------------

func TestDML_DashedPaths(t *testing.T) {
	// Dashed table paths are BigQuery-valid (dashed_path_expression) and omni
	// accepts them. The Spanner emulator syntax-rejects '-' in a table name, so
	// these are covered here, not in the differential (divergence ledger #85).
	t.Run("insert", func(t *testing.T) {
		ins := parseInsert(t, "INSERT INTO my-project.ds.tbl (a) VALUES (1)")
		if pathString(ins.Target) != "my-project.ds.tbl" {
			t.Errorf("Target = %q", pathString(ins.Target))
		}
	})
	t.Run("update", func(t *testing.T) {
		u := parseUpdate(t, "UPDATE my-project.ds.tbl SET x = 1 WHERE id = 1")
		if pathString(u.Target) != "my-project.ds.tbl" {
			t.Errorf("Target = %q", pathString(u.Target))
		}
	})
	t.Run("truncate", func(t *testing.T) {
		parseTruncate(t, "TRUNCATE TABLE my-project.ds.tbl")
	})
}

// ---------------------------------------------------------------------------
// Negative cases — omni must reject (oracle-confirmed syntax rejects)
// ---------------------------------------------------------------------------

func TestDML_Rejects(t *testing.T) {
	cases := []string{
		// Truncated / incomplete.
		"INSERT INTO",
		"INSERT t",
		"INSERT INTO t (a) VALUES",
		"UPDATE",
		"UPDATE t SET",
		"DELETE",
		"UPDATE t SET x = 1 FROM",
		"INSERT INTO t (a) VALUES (1) THEN RETURN",
		// Trailing comma in SET list (oracle: Expected "(" but got keyword WHERE).
		"UPDATE t SET x = 1, WHERE id = 1",
		// ON CONFLICT after a BARE query source is rejected (oracle: Expected end
		// of input but got keyword ON) — only VALUES / TABLE / ( query ) allow it.
		"INSERT INTO t (a) SELECT 1 ON CONFLICT DO NOTHING",
		// A parenthesized query WITH a query-level tail (LIMIT) is a bare query, so
		// a trailing ON CONFLICT is also rejected.
		"INSERT INTO t (a) (SELECT 1) LIMIT 5 ON CONFLICT DO NOTHING",
		// MERGE requires >= 1 WHEN clause (the .g4's (merge_when_clause)+).
		"MERGE INTO t USING s ON t.a = s.a",
		// MERGE missing USING / ON.
		"MERGE INTO t",
		"MERGE INTO t USING s ON",
		// Slashed target is not a DML target.
		"UPDATE /a/b SET x = 1 WHERE id = 1",
		// TRUNCATE without TABLE keyword.
		"TRUNCATE foo",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("Parse(%q): expected a syntax error, got none", sql)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Breadth: the entire canonical ZetaSQL DML corpus must parse cleanly
// ---------------------------------------------------------------------------

// TestDML_LegacyCorpusAccepts parses every statement of the canonical ZetaSQL
// DML corpus files and asserts each parses to a single DML node with no errors.
// This is the breadth oracle for the node (correctness-protocol completeness
// gate): the legacy grammar bytebase consumes is a hand-port of this ZetaSQL
// reference, so full corpus parity is the bar. Skips if the legacy checkout is
// absent (CI without it).
func TestDML_LegacyCorpusAccepts(t *testing.T) {
	dir := filepath.Join(legacyCorpusRoot, "zetasql", "parser", "testdata")
	files := []string{
		"dml_insert.sql",
		"dml_update.sql",
		"dml_delete.sql",
		"dml_merge.sql",
		"truncate.sql",
		"dml_insert_on_conflict_clause.sql",
	}
	for _, name := range files {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Skipf("legacy corpus file %s not available: %v", path, err)
		}
		// The truncate.sql corpus mixes in SELECTs that USE `truncate` as an
		// identifier; those are valid query statements, not DML. We only assert a
		// clean parse to a single statement (no error), not the node kind, for
		// that file.
		assertDMLNode := name != "truncate.sql"
		for i, raw := range splitStatements(string(data)) {
			stmt := strings.TrimSpace(raw)
			if stmt == "" {
				continue
			}
			file, errs := Parse(stmt)
			if len(errs) != 0 {
				t.Errorf("%s[%d] parse failed:\n%s\nerrs=%v", name, i, stmt, errs)
				continue
			}
			if len(file.Stmts) != 1 {
				t.Errorf("%s[%d]: got %d stmts:\n%s", name, i, len(file.Stmts), stmt)
				continue
			}
			if assertDMLNode && !isDMLNode(file.Stmts[0]) {
				t.Errorf("%s[%d]: got %T, want a DML node:\n%s", name, i, file.Stmts[0], stmt)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseInsert(t *testing.T, sql string) *ast.InsertStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	ins, ok := n.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.InsertStmt", sql, n)
	}
	return ins
}

func parseUpdate(t *testing.T, sql string) *ast.UpdateStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	u, ok := n.(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.UpdateStmt", sql, n)
	}
	return u
}

func parseDelete(t *testing.T, sql string) *ast.DeleteStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	d, ok := n.(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.DeleteStmt", sql, n)
	}
	return d
}

func parseMerge(t *testing.T, sql string) *ast.MergeStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	m, ok := n.(*ast.MergeStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.MergeStmt", sql, n)
	}
	return m
}

func parseTruncate(t *testing.T, sql string) *ast.TruncateStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	tr, ok := n.(*ast.TruncateStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.TruncateStmt", sql, n)
	}
	return tr
}

func pathString(p *ast.PathExpr) string {
	if p == nil {
		return ""
	}
	return strings.Join(p.Parts, ".")
}

func isDMLNode(n ast.Node) bool {
	switch n.(type) {
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt, *ast.TruncateStmt:
		return true
	}
	return false
}

// splitStatements splits a corpus file on top-level ';' (the corpus is
// formatted so no ';' appears inside a statement except as a terminator).
func splitStatements(src string) []string {
	return strings.Split(src, ";")
}
