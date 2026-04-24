package parser

import (
	"reflect"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestCollectCompletionJoinUsingIncompleteScope(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN t2 USING ("
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil {
		t.Fatal("expected completion context")
	}
	if ctx.Candidates == nil || !ctx.Candidates.HasRule("columnref") {
		t.Fatal("expected columnref candidate inside JOIN USING")
	}

	got := refNames(ctx.Scope.References)
	want := []string{"t1", "t2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %v, want %v", got, want)
	}
	for _, ref := range ctx.Scope.References {
		if ref.Kind != RangeReferenceRelation {
			t.Fatalf("reference %q kind = %v, want relation", ref.Name, ref.Kind)
		}
	}
}

func TestCollectCompletionJoinOnQualifiedColumnScope(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN t2 ON t1."
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil {
		t.Fatal("expected completion context")
	}
	if ctx.Candidates == nil || !ctx.Candidates.HasRule("columnref") {
		t.Fatal("expected columnref candidate after qualified column prefix")
	}
	got := refNames(ctx.Scope.References)
	want := []string{"t1", "t2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %v, want %v", got, want)
	}
}

func TestCollectCompletionLeftJoinUsingIncompleteScope(t *testing.T) {
	sql := "SELECT * FROM t1 LEFT JOIN t2 USING ("
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	got := refNames(ctx.Scope.References)
	want := []string{"t1", "t2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %v, want %v", got, want)
	}
}

func TestCollectCompletionIncompleteJoinFamilyScope(t *testing.T) {
	cases := []string{
		"SELECT * FROM t1 INNER JOIN t2 ON t2.",
		"SELECT * FROM t1 RIGHT JOIN t2 USING (",
		"SELECT * FROM t1 FULL OUTER JOIN t2 USING (",
		"SELECT * FROM t1 NATURAL JOIN t2 WHERE t1.",
		"SELECT * FROM t1 CROSS JOIN t2 WHERE t2.",
		"SELECT * FROM t1 JOIN t2",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			ctx := CollectCompletion(sql, len(sql))
			if ctx == nil || ctx.Scope == nil {
				t.Fatal("expected completion scope")
			}
			got := refNames(ctx.Scope.References)
			want := []string{"t1", "t2"}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("references = %v, want %v", got, want)
			}
		})
	}
}

func TestCollectCompletionAliasColumnsScope(t *testing.T) {
	sql := "SELECT * FROM public.t AS a(c1, c2) WHERE a."
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.Scope.References) != 1 {
		t.Fatalf("reference count = %d, want 1", len(ctx.Scope.References))
	}
	ref := ctx.Scope.References[0]
	if ref.Kind != RangeReferenceRelation || ref.Schema != "public" || ref.Name != "t" || ref.Alias != "a" {
		t.Fatalf("unexpected reference: %+v", ref)
	}
	if !reflect.DeepEqual(ref.AliasColumns, []string{"c1", "c2"}) {
		t.Fatalf("alias columns = %v, want [c1 c2]", ref.AliasColumns)
	}
}

func TestCollectCompletionIncompleteJoinKeepsLeftScope(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN"
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	got := refNames(ctx.Scope.References)
	want := []string{"t1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %v, want %v", got, want)
	}
}

func TestCollectCompletionSubqueryScope(t *testing.T) {
	sql := "SELECT * FROM (SELECT a FROM t) s WHERE s."
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.Scope.References) != 1 {
		t.Fatalf("reference count = %d, want 1", len(ctx.Scope.References))
	}
	ref := ctx.Scope.References[0]
	if ref.Kind != RangeReferenceSubquery || ref.Alias != "s" {
		t.Fatalf("unexpected subquery reference: %+v", ref)
	}
	if ref.BodyLoc.Start < 0 || ref.BodyLoc.End <= ref.BodyLoc.Start || ref.BodyLoc.End > len(sql) {
		t.Fatalf("invalid body loc: %+v", ref.BodyLoc)
	}
	if got, want := sql[ref.BodyLoc.Start:ref.BodyLoc.End], "SELECT a FROM t"; got != want {
		t.Fatalf("body text = %q, want %q", got, want)
	}
}

func TestCollectCompletionCTEScope(t *testing.T) {
	sql := "WITH x(a) AS (SELECT a FROM t) SELECT x. FROM x"
	cursor := len("WITH x(a) AS (SELECT a FROM t) SELECT x.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.Scope.References) != 1 {
		t.Fatalf("reference count = %d, want 1", len(ctx.Scope.References))
	}
	ref := ctx.Scope.References[0]
	if ref.Kind != RangeReferenceCTE || ref.Name != "x" {
		t.Fatalf("unexpected CTE reference: %+v", ref)
	}
	if !reflect.DeepEqual(ref.AliasColumns, []string{"a"}) {
		t.Fatalf("CTE columns = %v, want [a]", ref.AliasColumns)
	}
}

func TestCollectCompletionFunctionScope(t *testing.T) {
	sql := "SELECT * FROM generate_series(1, 3) AS g(n) WHERE g."
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.Scope.References) != 1 {
		t.Fatalf("reference count = %d, want 1", len(ctx.Scope.References))
	}
	ref := ctx.Scope.References[0]
	if ref.Kind != RangeReferenceFunction || ref.Alias != "g" {
		t.Fatalf("unexpected function reference: %+v", ref)
	}
	if !reflect.DeepEqual(ref.AliasColumns, []string{"n"}) {
		t.Fatalf("function alias columns = %v, want [n]", ref.AliasColumns)
	}
}

func TestCollectCompletionCorrelatedSubqueryScopeStack(t *testing.T) {
	sql := "SELECT * FROM t1 WHERE EXISTS (SELECT | FROM t2)"
	cursor := len("SELECT * FROM t1 WHERE EXISTS (SELECT ")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}

	if got, want := refNames(ctx.Scope.LocalReferences), []string{"t2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if len(ctx.Scope.OuterReferences) != 1 {
		t.Fatalf("outer reference levels = %d, want 1", len(ctx.Scope.OuterReferences))
	}
	if got, want := refNames(ctx.Scope.OuterReferences[0]), []string{"t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("outer references = %v, want %v", got, want)
	}
	if got, want := refNames(ctx.Scope.References), []string{"t2", "t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("visible references = %v, want %v", got, want)
	}
}

func TestCollectCompletionLateralSubquerySeesLeftSibling(t *testing.T) {
	sql := "SELECT * FROM t2 x, LATERAL (SELECT x.| FROM t1) sub"
	cursor := len("SELECT * FROM t2 x, LATERAL (SELECT x.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := refNames(ctx.Scope.LocalReferences), []string{"t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if len(ctx.Scope.OuterReferences) != 1 {
		t.Fatalf("outer reference levels = %d, want 1", len(ctx.Scope.OuterReferences))
	}
	if got, want := refNames(ctx.Scope.OuterReferences[0]), []string{"x"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("outer references = %v, want %v", got, want)
	}
	if got, want := refNames(ctx.Scope.References), []string{"t1", "x"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("visible references = %v, want %v", got, want)
	}
}

func TestCollectCompletionPlainFromSubqueryDoesNotSeeSibling(t *testing.T) {
	sql := "SELECT * FROM t2 x, (SELECT x.| FROM t1) sub"
	cursor := len("SELECT * FROM t2 x, (SELECT x.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := refNames(ctx.Scope.LocalReferences), []string{"t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if len(ctx.Scope.OuterReferences) != 0 {
		t.Fatalf("outer reference levels = %d, want 0", len(ctx.Scope.OuterReferences))
	}
	if got, want := refNames(ctx.Scope.References), []string{"t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("visible references = %v, want %v", got, want)
	}
}

func TestCollectCompletionInheritedCTEInNestedSelect(t *testing.T) {
	sql := "WITH x(a) AS (SELECT a FROM t) SELECT | FROM (SELECT * FROM x) sub"
	cursor := len("WITH x(a) AS (SELECT a FROM t) SELECT ")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.CTEs) != 1 || ctx.CTEs[0].Name != "x" {
		t.Fatalf("CTEs = %+v, want x", ctx.CTEs)
	}
	if !reflect.DeepEqual(ctx.CTEs[0].AliasColumns, []string{"a"}) {
		t.Fatalf("CTE columns = %v, want [a]", ctx.CTEs[0].AliasColumns)
	}
	if got, want := refNames(ctx.Scope.References), []string{"sub"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("visible references = %v, want %v", got, want)
	}
}

func TestCollectCompletionRecursiveCTEVisibleInsideRecursiveTerm(t *testing.T) {
	sql := "WITH RECURSIVE x(a) AS (SELECT a FROM t UNION ALL SELECT x.| FROM x JOIN t1 ON x.a = t1.c1) SELECT * FROM x"
	cursor := len("WITH RECURSIVE x(a) AS (SELECT a FROM t UNION ALL SELECT x.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.CTEs) != 1 || ctx.CTEs[0].Name != "x" {
		t.Fatalf("CTEs = %+v, want x", ctx.CTEs)
	}
	if got, want := refNames(ctx.Scope.LocalReferences), []string{"x", "t1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if ctx.Scope.LocalReferences[0].Kind != RangeReferenceCTE {
		t.Fatalf("first local reference kind = %v, want CTE", ctx.Scope.LocalReferences[0].Kind)
	}
}

func TestCollectCompletionDMLSelectInheritsStatementCTE(t *testing.T) {
	sql := "WITH x(a) AS (SELECT a FROM t) INSERT INTO dst SELECT | FROM x"
	cursor := len("WITH x(a) AS (SELECT a FROM t) INSERT INTO dst SELECT ")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.CTEs) != 1 || ctx.CTEs[0].Name != "x" {
		t.Fatalf("CTEs = %+v, want x", ctx.CTEs)
	}
	if got, want := refNames(ctx.Scope.LocalReferences), []string{"x"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if ctx.Scope.LocalReferences[0].Kind != RangeReferenceCTE {
		t.Fatalf("local reference kind = %v, want CTE", ctx.Scope.LocalReferences[0].Kind)
	}
}

func TestCollectCompletionMultiStatementUsesCursorStatement(t *testing.T) {
	sql := "SELECT * FROM t1; SELECT | FROM t2"
	cursor := len("SELECT * FROM t1; SELECT ")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := refNames(ctx.Scope.References), []string{"t2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("visible references = %v, want %v", got, want)
	}
}

func TestCollectCompletionParenthesizedJoinAliasScope(t *testing.T) {
	sql := "SELECT jt.| FROM (t1 JOIN t2 USING (c1)) jt(c1, c2)"
	cursor := len("SELECT jt.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if len(ctx.Scope.References) != 1 {
		t.Fatalf("reference count = %d, want 1", len(ctx.Scope.References))
	}
	ref := ctx.Scope.References[0]
	if ref.Kind != RangeReferenceJoinAlias || ref.Alias != "jt" {
		t.Fatalf("unexpected join alias reference: %+v", ref)
	}
	if !reflect.DeepEqual(ref.AliasColumns, []string{"c1", "c2"}) {
		t.Fatalf("join alias columns = %v, want [c1 c2]", ref.AliasColumns)
	}
}

func TestCollectCompletionUsesAstLocType(t *testing.T) {
	var _ nodes.Loc = RangeReference{}.Loc
	var _ nodes.Loc = RangeReference{}.BodyLoc
}

func TestParseOrdinaryJoinWithoutQualStaysStrict(t *testing.T) {
	if _, err := Parse("SELECT * FROM t1 JOIN t2"); err == nil {
		t.Fatal("expected ordinary Parse to reject JOIN without ON/USING")
	}
}

func refNames(refs []RangeReference) []string {
	var names []string
	for _, ref := range refs {
		if ref.Alias != "" {
			names = append(names, ref.Alias)
			continue
		}
		names = append(names, ref.Name)
	}
	return names
}
