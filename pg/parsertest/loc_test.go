package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

func TestLocSelectStmt(t *testing.T) {
	sql := "SELECT a, b FROM t WHERE x > 0"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if sel.Loc.Start == -1 || sel.Loc.End == -1 {
		t.Fatalf("SelectStmt Loc not set: %+v", sel.Loc)
	}
	got := sql[sel.Loc.Start:sel.Loc.End]
	if got != sql {
		t.Errorf("SelectStmt text = %q, want %q", got, sql)
	}
}

func TestLocSelectWithParens(t *testing.T) {
	// (SELECT 1) is not a valid top-level statement in this parser,
	// so test via UNION where the right side is a parenthesized select.
	sql := "SELECT 1 UNION (SELECT 2)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	// The top-level is UNION; check that Rarg (parenthesized select) has Loc for inner SELECT.
	rarg := sel.Rarg
	if rarg.Loc.Start == -1 || rarg.Loc.End == -1 {
		t.Fatalf("parenthesized SelectStmt Loc not set: %+v", rarg.Loc)
	}
	got := sql[rarg.Loc.Start:rarg.Loc.End]
	want := "SELECT 2"
	if got != want {
		t.Errorf("parenthesized SelectStmt text = %q, want %q", got, want)
	}
}

func TestLocSelectUnion(t *testing.T) {
	sql := "SELECT 1 UNION SELECT 2"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if sel.Loc.Start == -1 || sel.Loc.End == -1 {
		t.Fatalf("union SelectStmt Loc not set: %+v", sel.Loc)
	}
	got := sql[sel.Loc.Start:sel.Loc.End]
	if got != sql {
		t.Errorf("union text = %q, want %q", got, sql)
	}
}

func TestLocSelectMultiStmt(t *testing.T) {
	sql := "SELECT 1; SELECT 2"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range list.Items {
		raw := item.(*nodes.RawStmt)
		sel := raw.Stmt.(*nodes.SelectStmt)
		if sel.Loc.Start == -1 || sel.Loc.End == -1 {
			t.Errorf("stmt %d: SelectStmt Loc not set: %+v", i, sel.Loc)
		}
	}
}

func TestLocInsertStmt(t *testing.T) {
	sql := "INSERT INTO t VALUES (1, 2)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	ins := raw.Stmt.(*nodes.InsertStmt)
	if ins.Loc.Start == -1 || ins.Loc.End == -1 {
		t.Fatalf("InsertStmt Loc not set: %+v", ins.Loc)
	}
	got := sql[ins.Loc.Start:ins.Loc.End]
	if got != sql {
		t.Errorf("InsertStmt text = %q, want %q", got, sql)
	}
}

func TestLocInsertWithCTE(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) INSERT INTO t SELECT * FROM cte"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	ins := raw.Stmt.(*nodes.InsertStmt)
	if ins.Loc.Start == -1 || ins.Loc.End == -1 {
		t.Fatalf("InsertStmt Loc not set: %+v", ins.Loc)
	}
	got := sql[ins.Loc.Start:ins.Loc.End]
	if got != sql {
		t.Errorf("InsertStmt text = %q, want %q", got, sql)
	}
}

func TestLocUpdateStmt(t *testing.T) {
	sql := "UPDATE t SET a = 1 WHERE id = 5"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	upd := raw.Stmt.(*nodes.UpdateStmt)
	if upd.Loc.Start == -1 || upd.Loc.End == -1 {
		t.Fatalf("UpdateStmt Loc not set: %+v", upd.Loc)
	}
	got := sql[upd.Loc.Start:upd.Loc.End]
	if got != sql {
		t.Errorf("UpdateStmt text = %q, want %q", got, sql)
	}
}

func TestLocDeleteStmt(t *testing.T) {
	sql := "DELETE FROM t WHERE id = 5"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	del := raw.Stmt.(*nodes.DeleteStmt)
	if del.Loc.Start == -1 || del.Loc.End == -1 {
		t.Fatalf("DeleteStmt Loc not set: %+v", del.Loc)
	}
	got := sql[del.Loc.Start:del.Loc.End]
	if got != sql {
		t.Errorf("DeleteStmt text = %q, want %q", got, sql)
	}
}

func TestLocUpdateWithCTE(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) UPDATE t SET a = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	upd := raw.Stmt.(*nodes.UpdateStmt)
	got := sql[upd.Loc.Start:upd.Loc.End]
	if got != sql {
		t.Errorf("UpdateStmt text = %q, want %q", got, sql)
	}
}

func TestLocDeleteWithUsing(t *testing.T) {
	sql := "DELETE FROM t USING t2 WHERE t.id = t2.id"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	del := raw.Stmt.(*nodes.DeleteStmt)
	got := sql[del.Loc.Start:del.Loc.End]
	if got != sql {
		t.Errorf("DeleteStmt text = %q, want %q", got, sql)
	}
}

func TestLocMergeStmt(t *testing.T) {
	sql := "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a WHEN NOT MATCHED THEN INSERT (a) VALUES (s.a)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	merge := raw.Stmt.(*nodes.MergeStmt)
	if merge.Loc.Start == -1 || merge.Loc.End == -1 {
		t.Fatalf("MergeStmt Loc not set: %+v", merge.Loc)
	}
	got := sql[merge.Loc.Start:merge.Loc.End]
	if got != sql {
		t.Errorf("MergeStmt text = %q, want %q", got, sql)
	}
}

// Integration tests: extract sub-clause text via Loc slicing

func TestLocWhereClauseExtraction(t *testing.T) {
	sql := "UPDATE t SET a = 1 WHERE id > 5"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	upd := raw.Stmt.(*nodes.UpdateStmt)

	loc := nodes.NodeLoc(upd.WhereClause)
	if loc.Start == -1 {
		t.Fatal("WhereClause has no Loc")
	}
	got := sql[loc.Start:loc.End]
	if got != "id > 5" {
		t.Errorf("WHERE expression = %q, want %q", got, "id > 5")
	}
}

func TestLocWithClauseExtraction(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) SELECT * FROM cte"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	if sel.WithClause == nil {
		t.Fatal("WithClause is nil")
	}
	got := sql[sel.WithClause.Loc.Start:sel.WithClause.Loc.End]
	// Loc.End may include trailing whitespace (points to next token start).
	// Verify the extracted text starts with the expected clause.
	want := "WITH cte AS (SELECT 1)"
	if len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("WITH clause = %q, want prefix %q", got, want)
	}
}

func TestLocFromClauseExtraction(t *testing.T) {
	sql := "SELECT * FROM t1, t2"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	span := nodes.ListSpan(sel.FromClause)
	if span.Start == -1 {
		t.Fatal("FromClause has no span")
	}
	got := sql[span.Start:span.End]
	if got != "t1, t2" {
		t.Errorf("FROM clause = %q, want %q", got, "t1, t2")
	}
}

func TestLocLimitCountExtraction(t *testing.T) {
	sql := "SELECT * FROM t LIMIT 100"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	loc := nodes.NodeLoc(sel.LimitCount)
	if loc.Start == -1 {
		t.Fatal("LimitCount has no Loc")
	}
	got := sql[loc.Start:loc.End]
	if got != "100" {
		t.Errorf("LIMIT value = %q, want %q", got, "100")
	}
}

func TestLocSortClauseExtraction(t *testing.T) {
	sql := "SELECT * FROM t ORDER BY a, b DESC"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	span := nodes.ListSpan(sel.SortClause)
	if span.Start == -1 {
		t.Fatal("SortClause has no span")
	}
	got := sql[span.Start:span.End]
	if got != "a, b DESC" {
		t.Errorf("ORDER BY items = %q, want %q", got, "a, b DESC")
	}
}

func TestLocRelationExtraction(t *testing.T) {
	sql := "DELETE FROM public.users WHERE id = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	del := raw.Stmt.(*nodes.DeleteStmt)

	got := sql[del.Relation.Loc.Start:del.Relation.Loc.End]
	// Loc.End may include trailing whitespace.
	want := "public.users"
	if len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("relation = %q, want prefix %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Alias.Loc tests — verify that the new Loc field on Alias correctly spans
// the alias text (including AS keyword when present).

func TestLocAliasUpdateWithAS(t *testing.T) {
	sql := "UPDATE test AS t1 SET name = 'new'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.UpdateStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	alias := stmt.Relation.Alias
	got := sql[alias.Loc.Start:alias.Loc.End]
	if got != "AS t1" {
		t.Errorf("alias text = %q, want %q", got, "AS t1")
	}

	// Combined: relation + alias
	full := sql[stmt.Relation.Loc.Start:alias.Loc.End]
	if full != "test AS t1" {
		t.Errorf("relation+alias text = %q, want %q", full, "test AS t1")
	}
}

func TestLocAliasUpdateBare(t *testing.T) {
	sql := "UPDATE test t1 SET name = 'new'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.UpdateStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	alias := stmt.Relation.Alias
	got := sql[alias.Loc.Start:alias.Loc.End]
	if got != "t1" {
		t.Errorf("alias text = %q, want %q", got, "t1")
	}

	// Combined: relation + alias
	full := sql[stmt.Relation.Loc.Start:alias.Loc.End]
	if full != "test t1" {
		t.Errorf("relation+alias text = %q, want %q", full, "test t1")
	}
}

func TestLocAliasDeleteWithAS(t *testing.T) {
	sql := "DELETE FROM test AS t1 WHERE t1.id = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DeleteStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	alias := stmt.Relation.Alias
	got := sql[alias.Loc.Start:alias.Loc.End]
	if got != "AS t1" {
		t.Errorf("alias text = %q, want %q", got, "AS t1")
	}
}

func TestLocAliasDeleteBare(t *testing.T) {
	sql := "DELETE FROM test t1 WHERE t1.id = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DeleteStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	alias := stmt.Relation.Alias
	got := sql[alias.Loc.Start:alias.Loc.End]
	if got != "t1" {
		t.Errorf("alias text = %q, want %q", got, "t1")
	}
}

func TestLocAliasInsertWithAS(t *testing.T) {
	sql := "INSERT INTO test AS t1 VALUES (1)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.InsertStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	alias := stmt.Relation.Alias
	got := sql[alias.Loc.Start:alias.Loc.End]
	if got != "AS t1" {
		t.Errorf("alias text = %q, want %q", got, "AS t1")
	}
}

func TestLocAliasMergeWithAS(t *testing.T) {
	sql := "MERGE INTO target AS t USING source AS s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.MergeStmt)

	if stmt.Relation.Alias == nil {
		t.Fatal("expected target Alias to be set")
	}
	got := sql[stmt.Relation.Alias.Loc.Start:stmt.Relation.Alias.Loc.End]
	if got != "AS t" {
		t.Errorf("target alias text = %q, want %q", got, "AS t")
	}
}

func TestLocAliasSelectFromWithAS(t *testing.T) {
	sql := "SELECT * FROM users AS u WHERE u.id = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	if sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
		t.Fatal("expected FROM clause")
	}
	rv := sel.FromClause.Items[0].(*nodes.RangeVar)
	if rv.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	got := sql[rv.Alias.Loc.Start:rv.Alias.Loc.End]
	if got != "AS u" {
		t.Errorf("alias text = %q, want %q", got, "AS u")
	}

	// Combined: relation + alias
	full := sql[rv.Loc.Start:rv.Alias.Loc.End]
	if full != "users AS u" {
		t.Errorf("relation+alias text = %q, want %q", full, "users AS u")
	}
}

func TestLocAliasSelectFromBare(t *testing.T) {
	sql := "SELECT * FROM users u WHERE u.id = 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	rv := sel.FromClause.Items[0].(*nodes.RangeVar)
	if rv.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	got := sql[rv.Alias.Loc.Start:rv.Alias.Loc.End]
	if got != "u" {
		t.Errorf("alias text = %q, want %q", got, "u")
	}
}

func TestLocAliasSelectFromWithColumnAliases(t *testing.T) {
	sql := "SELECT * FROM users AS u(id, name)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)

	rv := sel.FromClause.Items[0].(*nodes.RangeVar)
	if rv.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	got := sql[rv.Alias.Loc.Start:rv.Alias.Loc.End]
	if got != "AS u(id, name)" {
		t.Errorf("alias text = %q, want %q", got, "AS u(id, name)")
	}
}

func TestLocAliasNoAlias(t *testing.T) {
	sql := "UPDATE test SET name = 'new'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.UpdateStmt)

	if stmt.Relation.Alias != nil {
		t.Errorf("expected no alias, got %+v", stmt.Relation.Alias)
	}
}
