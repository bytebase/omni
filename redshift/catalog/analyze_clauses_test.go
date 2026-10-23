package catalog

import (
	"errors"
	"slices"
	"strings"
	"testing"

	pg "github.com/bytebase/omni/redshift"
	nodes "github.com/bytebase/omni/redshift/ast"
)

const analyzeClausesSchema = `
	CREATE TABLE t (id int, masked text, x text);
	CREATE TABLE u (id int, y text);
	CREATE TABLE w (id int, x text, z int);
`

func newAnalyzeClausesCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	if _, err := c.Exec(analyzeClausesSchema, nil); err != nil {
		t.Fatal(err)
	}
	return c
}

func analyzeWithOptions(t *testing.T, c *Catalog, sql string, opts AnalyzeOptions) (*Query, error) {
	t.Helper()
	stmts, err := pg.Parse(sql)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("want one statement in %q, got %d", sql, len(stmts))
	}
	sel, ok := stmts[0].AST.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("%q is not a SELECT", sql)
	}
	return c.AnalyzeSelectStmtWithOptions(sel, opts)
}

// readColumns names every table column an expression reads, as
// "table.column", resolving Vars of joins through their alias expressions.
func readColumns(q *Query, expr AnalyzedExpr) []string {
	var cols []string
	var visit func(e AnalyzedExpr) bool
	visit = func(e AnalyzedExpr) bool {
		v, ok := e.(*VarExpr)
		if !ok || v.LevelsUp != 0 || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
			return true
		}
		rte := q.RangeTable[v.RangeIdx]
		switch rte.Kind {
		case RTERelation:
			cols = append(cols, rte.RelName+"."+rte.ColNames[v.AttNum-1])
		case RTEJoin:
			WalkAnalyzedExpr(rte.JoinAliasVars[v.AttNum-1], visit)
		}
		return true
	}
	WalkAnalyzedExpr(expr, visit)
	slices.Sort(cols)
	return slices.Compact(cols)
}

// targetExpr returns the expression of output column i (0-based), skipping
// the junk entries ORDER BY, GROUP BY, and window clauses add.
func targetExpr(t *testing.T, q *Query, i int) AnalyzedExpr {
	t.Helper()
	output := nonJunkTargetEntries(q.TargetList)
	if i >= len(output) {
		t.Fatalf("query has %d output columns, want column %d", len(output), i+1)
	}
	return output[i].Expr
}

func TestAnalyzeAggregateKeepsFilterAndOrder(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT
		sum(id) FILTER (WHERE masked = 'a'),
		count(*) FILTER (WHERE masked IS NULL),
		string_agg(x, ',' ORDER BY masked DESC NULLS LAST)
		FROM t`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range [][]string{
		{"t.id", "t.masked"},
		{"t.masked"},
		{"t.masked", "t.x"},
	} {
		if got := readColumns(q, targetExpr(t, q, i)); !slices.Equal(got, want) {
			t.Errorf("target %d reads %v, want %v", i+1, got, want)
		}
	}

	count := targetExpr(t, q, 1).(*AggExpr)
	if count.AggFuncOID != 2803 {
		t.Errorf("count(*) resolved to OID %d, want 2803", count.AggFuncOID)
	}
	agg := targetExpr(t, q, 2).(*AggExpr)
	if len(agg.Args) != 2 || len(agg.AggOrder) != 1 || agg.AggWithinGroup {
		t.Fatalf("string_agg: %d args, %d order items, within group %v", len(agg.Args), len(agg.AggOrder), agg.AggWithinGroup)
	}
	if sb := agg.AggOrder[0]; !sb.Descending || sb.NullsFirst {
		t.Errorf("ORDER BY masked DESC NULLS LAST analyzed as descending=%v nullsFirst=%v", sb.Descending, sb.NullsFirst)
	}
}

func TestAnalyzeWithinGroupSplitsDirectArgs(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) FROM t`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	agg := targetExpr(t, q, 0).(*AggExpr)
	if agg.AggFuncOID != 3974 || !agg.AggWithinGroup || len(agg.Args) != 1 || len(agg.AggOrder) != 1 {
		t.Fatalf("got OID %d, within group %v, %d direct args, %d order items", agg.AggFuncOID, agg.AggWithinGroup, len(agg.Args), len(agg.AggOrder))
	}
	// Both the direct argument and the ordering are coerced to the
	// aggregate's float8 parameters.
	if got := agg.Args[0].exprType(); got != FLOAT8OID {
		t.Errorf("direct argument type %d, want float8", got)
	}
	if got := agg.AggOrder[0].Expr.exprType(); got != FLOAT8OID {
		t.Errorf("ordering type %d, want float8", got)
	}
	if got := readColumns(q, agg); !slices.Equal(got, []string{"t.id"}) {
		t.Errorf("reads %v, want [t.id]", got)
	}
}

func TestAnalyzeNaturalJoin(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT *, id FROM t NATURAL JOIN w`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	jn := q.JoinTree.FromList[0].(*JoinExprNode)
	if !slices.Equal(jn.UsingClause, []string{"id", "x"}) {
		t.Errorf("USING columns %v, want [id x]", jn.UsingClause)
	}
	if got := readColumns(q, jn.Quals); !slices.Equal(got, []string{"t.id", "t.x", "w.id", "w.x"}) {
		t.Errorf("join condition reads %v", got)
	}
	var names []string
	for _, te := range q.TargetList {
		names = append(names, te.ResName)
	}
	if want := []string{"id", "x", "masked", "z", "id"}; !slices.Equal(names, want) {
		t.Errorf("output columns %v, want %v", names, want)
	}
	// The merged column of an inner join is the left input's.
	v := targetExpr(t, q, 4).(*VarExpr)
	if v.RangeIdx != 0 || v.AttNum != 1 {
		t.Errorf("unqualified id analyzed as RangeIdx %d AttNum %d, want t.id", v.RangeIdx, v.AttNum)
	}
}

func TestAnalyzeUsingColumnSides(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT id FROM t RIGHT JOIN u USING (id)`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := readColumns(q, targetExpr(t, q, 0)); !slices.Equal(got, []string{"u.id"}) {
		t.Errorf("RIGHT JOIN USING column reads %v, want [u.id]", got)
	}

	q, err = analyzeWithOptions(t, c, `SELECT id FROM t FULL JOIN u USING (id)`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	v := targetExpr(t, q, 0).(*VarExpr)
	if q.RangeTable[v.RangeIdx].Kind != RTEJoin {
		t.Fatalf("FULL JOIN USING column analyzed as a Var of %v, want the join", q.RangeTable[v.RangeIdx].Kind)
	}
	if _, ok := q.RangeTable[v.RangeIdx].JoinAliasVars[v.AttNum-1].(*CoalesceExprQ); !ok {
		t.Errorf("FULL JOIN USING column alias is %T, want COALESCE", q.RangeTable[v.RangeIdx].JoinAliasVars[v.AttNum-1])
	}
	if got := readColumns(q, v); !slices.Equal(got, []string{"t.id", "u.id"}) {
		t.Errorf("FULL JOIN USING column reads %v, want both inputs", got)
	}
}

func TestAnalyzeJoinColumnsResolveToInputs(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT y, masked FROM t JOIN u ON t.id = u.id`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"u.y", "t.masked"} {
		v := targetExpr(t, q, i).(*VarExpr)
		if rte := q.RangeTable[v.RangeIdx]; rte.Kind != RTERelation || rte.RelName+"."+rte.ColNames[v.AttNum-1] != want {
			t.Errorf("target %d is a Var of %s (kind %v), want %s", i+1, rte.RelName, rte.Kind, want)
		}
	}
}

func TestAnalyzeSetOpOutputColumns(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT id, x FROM t UNION SELECT id, y FROM u`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for i, te := range q.TargetList {
		col, ok := te.Expr.(*SetOpColumnExpr)
		if !ok {
			t.Fatalf("output column %d is %T, want *SetOpColumnExpr", i+1, te.Expr)
		}
		if int(col.ResNo) != i+1 {
			t.Errorf("output column %d has ResNo %d", i+1, col.ResNo)
		}
	}
	if coll := q.TargetList[1].Expr.exprCollation(); coll != DEFAULT_COLLATION_OID {
		t.Errorf("text output column collation %d, want default", coll)
	}
}

func TestAnalyzeBuiltinArgDefaults(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT jsonb_path_exists('{}'::jsonb, '$.a'), jsonb_set('{}'::jsonb, '{a}', '1'), random_normal()`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []uint32{4005, 3305, 6212} {
		if f := targetExpr(t, q, i).(*FuncCallExpr); f.FuncOID != want {
			t.Errorf("target %d resolved to OID %d, want %d", i+1, f.FuncOID, want)
		}
	}
	if got := targetExpr(t, q, 0).(*FuncCallExpr).Args[1].exprType(); got != JSONPATHOID {
		t.Errorf("jsonb_path_exists path argument type %d, want jsonpath", got)
	}
}

func TestAnalyzeNamedArgsAndVariadicDefaults(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT
		jsonb_path_exists(path => '$.a', target => '{}'::jsonb),
		make_interval(1, mins => id),
		pg_logical_slot_get_changes('slot', NULL, NULL)
		FROM t`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	f := targetExpr(t, q, 0).(*FuncCallExpr)
	if f.FuncOID != 4005 {
		t.Fatalf("jsonb_path_exists resolved to OID %d", f.FuncOID)
	}
	// Each named argument records its parameter and is coerced to its type.
	path, target := f.Args[0].(*NamedArgExprQ), f.Args[1].(*NamedArgExprQ)
	if path.ArgNum != 1 || target.ArgNum != 0 || path.exprType() != JSONPATHOID {
		t.Errorf("path => ArgNum %d type %d, target => ArgNum %d", path.ArgNum, path.exprType(), target.ArgNum)
	}
	if mi := targetExpr(t, q, 1).(*FuncCallExpr); mi.FuncOID != 3464 || mi.Args[1].(*NamedArgExprQ).ArgNum != 5 {
		t.Errorf("make_interval(1, mins => id) analyzed as %#v", mi)
	}
	if f := targetExpr(t, q, 2).(*FuncCallExpr); f.FuncOID != 3782 {
		t.Errorf("pg_logical_slot_get_changes resolved to OID %d, want 3782", f.FuncOID)
	}
}

func TestAnalyzeOpaqueUnknownFunctions(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	sql := `SELECT levenshtein(x, 'a'), my_agg(x ORDER BY masked) FILTER (WHERE id > 1),
		my_win(x) OVER (PARTITION BY masked), lower(x, 1)
		FROM t WHERE levenshtein(x, masked) < 3`
	var catErr *Error
	if _, err := analyzeWithOptions(t, c, sql, AnalyzeOptions{}); !errors.As(err, &catErr) || catErr.Code != CodeUndefinedFunction {
		t.Fatalf("without the option: %v, want 42883", err)
	}

	q, err := analyzeWithOptions(t, c, sql, AnalyzeOptions{OpaqueUnknownFunctions: true})
	if err != nil {
		t.Fatal(err)
	}
	if f, ok := targetExpr(t, q, 0).(*FuncCallExpr); !ok || f.FuncOID != 0 || f.FuncName != "levenshtein" || f.ResultType != UNKNOWNOID {
		t.Errorf("levenshtein analyzed as %#v", targetExpr(t, q, 0))
	}
	if a, ok := targetExpr(t, q, 1).(*AggExpr); !ok || a.AggFuncOID != 0 || len(a.AggOrder) != 1 || a.AggFilter == nil {
		t.Errorf("my_agg analyzed as %#v", targetExpr(t, q, 1))
	}
	if w, ok := targetExpr(t, q, 2).(*WindowFuncExpr); !ok || w.FuncOID != 0 {
		t.Errorf("my_win analyzed as %#v", targetExpr(t, q, 2))
	}
	// A known name with no matching overload is opaque too.
	if f, ok := targetExpr(t, q, 3).(*FuncCallExpr); !ok || f.FuncOID != 0 {
		t.Errorf("lower(x, 1) analyzed as %#v", targetExpr(t, q, 3))
	}
	for i, want := range [][]string{{"t.x"}, {"t.id", "t.masked", "t.x"}, {"t.x"}, {"t.x"}} {
		if got := readColumns(q, targetExpr(t, q, i)); !slices.Equal(got, want) {
			t.Errorf("target %d reads %v, want %v", i+1, got, want)
		}
	}
	if got := readColumns(q, q.JoinTree.Quals); !slices.Equal(got, []string{"t.masked", "t.x"}) {
		t.Errorf("WHERE reads %v", got)
	}

	// The non-recursive term of a recursive CTE is analyzed with the
	// options too.
	if _, err := analyzeWithOptions(t, c, `WITH RECURSIVE r(n) AS (SELECT ext(id)::int FROM t UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r`, AnalyzeOptions{OpaqueUnknownFunctions: true}); err != nil {
		t.Errorf("recursive CTE seed: %v", err)
	}
}

func TestAnalyzeRejectsQualifiedAggSortOperator(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	var catErr *Error
	_, err := analyzeWithOptions(t, c, `SELECT array_agg(id ORDER BY id USING OPERATOR(my_schema.>)) FROM t`, AnalyzeOptions{})
	if !errors.As(err, &catErr) || catErr.Code != CodeFeatureNotSupported {
		t.Fatalf("got %v, want 0A000: a user operator's ordering cannot be kept as ASC or DESC", err)
	}
	if _, err := analyzeWithOptions(t, c, `SELECT array_agg(id ORDER BY id USING OPERATOR(pg_catalog.>)) FROM t`, AnalyzeOptions{}); err != nil {
		t.Fatal(err)
	}
}

// TestAnalyzeJunkEntriesFollowOutputColumns checks that the junk entries a
// window's PARTITION BY and ORDER BY add come after the output columns, so
// a subquery column's AttNum indexes its target entry.
func TestAnalyzeJunkEntriesFollowOutputColumns(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	q, err := analyzeWithOptions(t, c, `SELECT a.x2 FROM (
		SELECT id, sum(id) OVER (PARTITION BY masked ORDER BY x) AS s, x AS x2 FROM t ORDER BY masked || x
	) a`, AnalyzeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	sub := q.RangeTable[0].Subquery
	var names []string
	for i, te := range sub.TargetList {
		if int(te.ResNo) != i+1 {
			t.Errorf("entry %d has ResNo %d", i, te.ResNo)
		}
		if !te.ResJunk {
			names = append(names, te.ResName)
		} else if len(names) < 3 {
			t.Errorf("junk entry %q at %d precedes an output column", te.ResName, i)
		}
	}
	if !slices.Equal(names, []string{"id", "s", "x2"}) {
		t.Errorf("output columns %v", names)
	}
	v := targetExpr(t, q, 0).(*VarExpr)
	if got := readColumns(sub, sub.TargetList[v.AttNum-1].Expr); !slices.Equal(got, []string{"t.x"}) {
		t.Errorf("a.x2 (AttNum %d) reads %v through TargetList[AttNum-1], want [t.x]", v.AttNum, got)
	}
}

// TestAnalyzeOpaqueFunctionsInFrom checks how OpaqueUnknownFunctions
// analyzes a function in FROM that the catalog does not know.
func TestAnalyzeOpaqueFunctionsInFrom(t *testing.T) {
	c := newAnalyzeClausesCatalog(t)
	opaque := AnalyzeOptions{OpaqueUnknownFunctions: true}
	columns := func(q *Query, name string) []string {
		for _, rte := range q.RangeTable {
			if rte.ERef == name {
				if !rte.UnknownColumns {
					t.Errorf("%s does not have UnknownColumns", name)
				}
				return rte.ColNames
			}
		}
		t.Fatalf("no range table entry %s", name)
		return nil
	}

	// Without the option, the call fails as in PostgreSQL.
	var catErr *Error
	if _, err := analyzeWithOptions(t, c, `SELECT 1 FROM t CROSS JOIN unknown_rows() f`, AnalyzeOptions{}); !errors.As(err, &catErr) || catErr.Code != CodeUndefinedFunction {
		t.Fatalf("without the option: %v, want 42883", err)
	}

	// Each column the query names becomes a column of the function; a
	// name a table has stays the table's.
	q, err := analyzeWithOptions(t, c, `SELECT t.x, f.a, id, b FROM t CROSS JOIN LATERAL unknown_rows(t.id) f WHERE f.c > 1`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	if got := columns(q, "f"); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("f has columns %v, want [a b c]", got)
	}
	if got := readColumns(q, targetExpr(t, q, 2)); !slices.Equal(got, []string{"t.id"}) {
		t.Errorf("id reads %v, want [t.id]", got)
	}
	fn := q.RangeTable[1].FuncExprs[0].(*FuncCallExpr)
	if got := readColumns(q, fn); !slices.Equal(got, []string{"t.id"}) {
		t.Errorf("the function call reads %v, want [t.id]", got)
	}

	// The alias names the first columns; a correlated reference adds one
	// at the level that owns the function.
	q, err = analyzeWithOptions(t, c, `SELECT (SELECT f.extra), a FROM unknown_rows() AS f(a, b)`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	if got := columns(q, "f"); !slices.Equal(got, []string{"a", "b", "extra"}) {
		t.Errorf("f has columns %v, want [a b extra]", got)
	}
	sub := targetExpr(t, q, 0).(*SubLinkExpr).SubQuery
	if v := sub.TargetList[0].Expr.(*VarExpr); v.LevelsUp != 1 || v.AttNum != 3 {
		t.Errorf("f.extra in the subquery is %#v, want LevelsUp 1 AttNum 3", v)
	}

	// A column definition list gives the exact columns, so * works.
	q, err = analyzeWithOptions(t, c, `SELECT * FROM unknown_rows() AS f(a int, b text)`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.TargetList) != 2 || q.RangeTable[0].UnknownColumns || q.TargetList[1].Expr.exprType() != TEXTOID {
		t.Errorf("column definition list gave %d columns, UnknownColumns %v", len(q.TargetList), q.RangeTable[0].UnknownColumns)
	}

	for _, sql := range []string{
		// * cannot list unknown columns.
		`SELECT * FROM t CROSS JOIN unknown_rows() f`,
		`SELECT f.* FROM unknown_rows() f`,
		// An unqualified name two such functions could have is unresolved.
		`SELECT a FROM unknown_rows() f, other_rows() g`,
		// ORDINALITY follows columns whose number is unknown.
		`SELECT f.a FROM unknown_rows() WITH ORDINALITY f`,
	} {
		if _, err := analyzeWithOptions(t, c, sql, opaque); err == nil {
			t.Errorf("%s: analyzed, want an error", sql)
		}
	}
	// Columns the query text says exist, by alias or USING, resolve like a
	// table's: over inferred ones, and merged by USING into one.
	for _, sql := range []string{
		`SELECT f.a, g.a FROM unknown_rows() f, other_rows() g`,
		`SELECT a FROM unknown_rows() f(a), other_rows() g`,
		`SELECT id FROM unknown_rows() f JOIN other_rows() g USING (id)`,
	} {
		if _, err := analyzeWithOptions(t, c, sql, opaque); err != nil {
			t.Errorf("%s: %v", sql, err)
		}
	}

	// A USING name is a column of the function, on either side.
	for _, sql := range []string{
		`SELECT id FROM unknown_rows() f JOIN t USING (id)`,
		`SELECT id FROM t JOIN unknown_rows() f USING (id)`,
		`SELECT x FROM (unknown_rows() f CROSS JOIN u) JOIN t USING (x)`,
	} {
		q, err := analyzeWithOptions(t, c, sql, opaque)
		if err != nil {
			t.Errorf("%s: %v", sql, err)
			continue
		}
		if jn := q.JoinTree.FromList[0].(*JoinExprNode); jn.Quals == nil {
			t.Errorf("%s: no join condition", sql)
		}
	}

	// A real column of an unqualified name wins over a function column the
	// query named, whatever the order of the references and the join.
	for _, sql := range []string{
		`SELECT f.id, id FROM t CROSS JOIN unknown_rows() f`,
		`SELECT id, f.id FROM t CROSS JOIN unknown_rows() f`,
		`SELECT f.id, id FROM t, unknown_rows() f`,
		`SELECT id, f.id FROM t, unknown_rows() f`,
	} {
		q, err := analyzeWithOptions(t, c, sql, opaque)
		if err != nil {
			t.Errorf("%s: %v", sql, err)
			continue
		}
		var id AnalyzedExpr
		for _, te := range nonJunkTargetEntries(q.TargetList) {
			if v, ok := te.Expr.(*VarExpr); ok && !q.RangeTable[v.RangeIdx].UnknownColumns {
				id = v
			}
		}
		if got := readColumns(q, id); !slices.Equal(got, []string{"t.id"}) {
			t.Errorf("%s: id reads %v, want [t.id]", sql, got)
		}
	}
	// Two functions' columns of one name are still ambiguous.
	if _, err := analyzeWithOptions(t, c, `SELECT f.a, g.a, a FROM unknown_rows() f CROSS JOIN other_rows() g`, opaque); !errors.As(err, &catErr) || catErr.Code != CodeAmbiguousColumn {
		t.Errorf("a over two functions' columns: %v, want 42702", err)
	}

	// An alias this level has owns a qualified reference: the column is
	// the inner function's, not the outer table's of the same alias.
	q, err = analyzeWithOptions(t, c, `SELECT (SELECT f.masked FROM unknown_rows() f) FROM t f`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	if v := targetExpr(t, q, 0).(*SubLinkExpr).SubQuery.TargetList[0].Expr.(*VarExpr); v.LevelsUp != 0 {
		t.Errorf("f.masked in the subquery reads the outer t (LevelsUp %d)", v.LevelsUp)
	}

	// An ON condition cannot see an alias-seeded column of another FROM
	// item, and a NATURAL join cannot match unknown columns. A local alias
	// shadows an enclosing function's; a USING join that matched a name
	// keeps a function below it from gaining it; and an unqualified name
	// stays ambiguous while another function could have it, whatever the
	// order of the references.
	for _, sql := range []string{
		`SELECT 1 FROM unknown_rows() f(x), t JOIN u ON f.x = t.id`,
		`SELECT f.id, id FROM unknown_rows() f NATURAL JOIN t`,
		`SELECT (SELECT f.zz FROM u f) FROM unknown_rows() f`,
		`SELECT f.id FROM (t CROSS JOIN unknown_rows() f) JOIN u USING (id)`,
		`SELECT f.a, a FROM unknown_rows() f, other_rows() g`,
		`SELECT a, f.a FROM unknown_rows() f, other_rows() g`,
		`SELECT id FROM t, unknown_rows() f(id)`,
		`SELECT id FROM t, (unknown_rows() f JOIN u USING (id))`,
		`SELECT f.a, (SELECT a FROM unknown_one() g, unknown_two() h) FROM unknown_outer() f`,
		`SELECT (SELECT a FROM unknown_one() g, unknown_two() h), f.a FROM unknown_outer() f`,
	} {
		if _, err := analyzeWithOptions(t, c, sql, opaque); err == nil {
			t.Errorf("%s: analyzed, want an error", sql)
		}
	}

	// An ON condition sees only the join's inputs.
	if _, err := analyzeWithOptions(t, c, `SELECT 1 FROM unknown_rows() f, t JOIN u ON f.x = t.id`, opaque); err == nil {
		t.Error("f.x outside the ON condition's join: analyzed, want an error")
	}
	if _, err := analyzeWithOptions(t, c, `SELECT 1 FROM t JOIN unknown_rows() f ON f.x = t.id`, opaque); err != nil {
		t.Errorf("f.x inside the ON condition's join: %v", err)
	}

	// A name two functions at the inner level could own is unresolved,
	// not taken by an enclosing level's function.
	q, err = analyzeWithOptions(t, c, `SELECT (SELECT a FROM unknown_one() f, unknown_two() g) FROM unknown_outer() h`, opaque)
	if err == nil {
		t.Errorf("ambiguous inner name: analyzed, h has columns %v", q.RangeTable[0].ColNames)
	}

	// A nearer function takes a name an enclosing function only inferred.
	q, err = analyzeWithOptions(t, c, `SELECT f.a, (SELECT a FROM unknown_one() g) FROM unknown_outer() f`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	if v := targetExpr(t, q, 1).(*SubLinkExpr).SubQuery.TargetList[0].Expr.(*VarExpr); v.LevelsUp != 0 {
		t.Errorf("inner a reads the outer f (LevelsUp %d), want the inner g", v.LevelsUp)
	}

	// Deparse gives only the alias's column names, which are positional.
	q, err = analyzeWithOptions(t, c, `SELECT (SELECT f.b FROM ext_rows() f), (SELECT g.c FROM ext_rows() AS g(a, b)),
		(SELECT h.a FROM ext_rows() AS h(a int, b varchar(5))),
		(SELECT k.a FROM ext_rows() AS k(a text COLLATE tenant.mycoll, b text COLLATE pg_catalog."C"))`, opaque)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{
		"FROM ext_rows() f)",
		"FROM ext_rows() g(a, b))",
		"FROM ext_rows() h(a integer, b character varying(5)))",
		`FROM ext_rows() k(a text COLLATE tenant.mycoll, b text COLLATE "C"))`,
	} {
		if got := c.DeparseExpr(targetExpr(t, q, i), q.RangeTable, false); !strings.Contains(got, want) {
			t.Errorf("target %d deparses as %s, want %s", i+1, got, want)
		}
	}
}
