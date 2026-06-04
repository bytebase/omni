package parser

import (
	"testing"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the parser-select node's slice of the correctness gate
// (correctness-protocol.md) for the querySpecification / queryNoWith / GROUP BY
// surface. The differential gate (TestSelect_OracleDifferential) is the
// authoritative accept/reject check: omni's Parse accept/reject must equal Trino
// 481's, error-classified (a SYNTAX_ERROR is a reject; any other error —
// TABLE_NOT_FOUND, COLUMN_NOT_FOUND, TYPE_MISMATCH, NOT_SUPPORTED,
// MISSING_ORDER_BY — means Trino's parser ACCEPTED the statement). The oracle
// helpers (connectOracle / oracleAccepts / truncateName) live in
// oracle_foundation_test.go; the relation / CTE / set-op / window corpora live in
// their sibling *_test.go files. The structural tests pin precedence and clause
// shapes the accept/reject gate cannot see.
//
// Corpus sources (correctness-protocol §completeness): every query form of the
// legacy ANTLR grammar (truth2) + every documented Trino 481 form (truth1,
// docs/migration/trino/truth1/query-dml.md) + the legacy examples
// (/Users/h3n4l/OpenSource/parser/trino/examples/*.sql) + oracle-probed edge
// cases. Negative (reject) inputs are required and present.

// selectAcceptCorpus is the querySpecification / queryNoWith / GROUP BY surface
// Trino 481 accepts (its parser does not raise SYNTAX_ERROR). The relation /
// join, CTE, set-op, and window forms have their own accept corpora in the
// sibling files.
var selectAcceptCorpus = []string{
	// --- minimal SELECT ---
	"SELECT 1",
	"SELECT 1, 2, 3",
	"SELECT 1 + 2 * 3",
	"SELECT 'a', TRUE, NULL",
	"SELECT a, b, c FROM t",
	"SELECT * FROM t",

	// --- setQuantifier ---
	"SELECT DISTINCT a FROM t",
	"SELECT ALL a FROM t",
	"SELECT DISTINCT a, b FROM t",

	// --- select items: alias forms ---
	"SELECT a AS x FROM t",
	"SELECT a x FROM t", // bare alias
	"SELECT a + b AS s FROM t",
	"SELECT a x, b AS y, c FROM t",
	"SELECT count(*) AS n FROM t",
	`SELECT a AS "Quoted Alias" FROM t`,

	// --- selectAll variants (the three selectItem alternatives) ---
	"SELECT *",
	"SELECT * FROM t",
	"SELECT t.* FROM t",
	"SELECT a.b.* FROM cat.sch.tbl a",
	"SELECT r.* FROM t",
	"SELECT ROW (1, 'a', true).*",
	"SELECT ROW (1, 'a', true).* AS (f1, f2, f3)",
	"SELECT (CAST(ROW(1, true) AS ROW(field1 bigint, field2 boolean))).* AS (alias1, alias2)",
	"SELECT a.id AS hello FROM t",
	"SELECT col1.f1, col2, col3.f1.f2.f3 FROM table1",
	"SELECT col1.f1[0], col2, col3[2].f2.f3, col4[4] FROM table1",

	// --- WHERE / HAVING ---
	"SELECT a FROM t WHERE a > 1",
	"SELECT a FROM t WHERE a > 1 AND b < 2 OR c = 3",
	"SELECT a FROM t WHERE a IN (1, 2, 3)",
	"SELECT a FROM t WHERE a BETWEEN 1 AND 10",
	"SELECT a, count(*) FROM t GROUP BY a HAVING count(*) > 5",

	// --- GROUP BY: plain, sets, rollup, cube, grouping sets ---
	"SELECT a FROM t GROUP BY a",
	"SELECT a, b FROM t GROUP BY a, b",
	"SELECT a FROM t GROUP BY ()",
	"SELECT a FROM t GROUP BY (a, b)",
	"SELECT a FROM t GROUP BY a, b, (c, d)",
	"SELECT a FROM t GROUP BY ROLLUP (a, b)",
	"SELECT a FROM t GROUP BY ROLLUP ()",
	"SELECT a FROM t GROUP BY CUBE (a, b)",
	"SELECT a FROM t GROUP BY GROUPING SETS (a)",
	"SELECT a, b, GROUPING(a, b) FROM t GROUP BY GROUPING SETS ((a), (b))",
	"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b), ())",
	"SELECT a FROM t GROUP BY ALL GROUPING SETS ((a, b), (a), ()), CUBE (c), ROLLUP (d)",
	"SELECT a FROM t GROUP BY DISTINCT GROUPING SETS ((a, b), (a), ()), CUBE (c), ROLLUP (d)",
	// GROUP BY setQuantifier disambiguation (S3, oracle-probed)
	"SELECT a FROM t GROUP BY ALL x",         // ALL is the quantifier
	"SELECT count(*) FROM t GROUP BY ALL",    // ALL is a grouping expr (column "all")
	"SELECT count(*) FROM t GROUP BY ALL, x", // ALL is a grouping expr (comma after)
	"SELECT a FROM t GROUP BY DISTINCT x",
	"SELECT count(*) FROM t GROUP BY AUTO", // AUTO keyword form (parses)
	// parenthesized-expr vs grouping-set ambiguity (oracle-probed)
	"SELECT a FROM t GROUP BY (a + b) * 2",
	"SELECT a FROM t GROUP BY (a) + 1",
	"SELECT a FROM t GROUP BY (a) IS NULL",

	// --- ORDER BY / OFFSET / LIMIT / FETCH ---
	"SELECT a FROM t ORDER BY a",
	"SELECT a FROM t ORDER BY a DESC",
	"SELECT a FROM t ORDER BY a ASC NULLS FIRST, b DESC NULLS LAST",
	"SELECT a FROM t ORDER BY 1, 2",
	"SELECT a FROM t LIMIT 10",
	"SELECT a FROM t LIMIT ALL",
	"SELECT a FROM t ORDER BY a LIMIT 5",
	"SELECT a FROM t OFFSET 5",
	"SELECT a FROM t OFFSET 5 ROWS",
	"SELECT a FROM t OFFSET 5 ROW",
	"SELECT a FROM t OFFSET ? ROWS LIMIT ?",
	"SELECT a FROM t FETCH FIRST 3 ROWS ONLY",
	"SELECT a FROM t FETCH NEXT ROW ONLY",
	"SELECT a FROM t FETCH FIRST ROW WITH TIES", // ORDER BY (required semantically) goes BEFORE FETCH; see reject corpus
	"SELECT a FROM t ORDER BY a OFFSET 2 ROWS FETCH FIRST 3 ROWS WITH TIES",
	"SELECT a FROM t ORDER BY a OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY",

	// --- queryPrimary: TABLE / VALUES ---
	"TABLE foo",
	"TABLE cat.sch.tbl",
	"VALUES 1, 2, 3",
	"VALUES (1, 'a'), (2, 'b')",
	"VALUES 1",
	"SELECT * FROM (VALUES 42, 13)",

	// --- whole-query decorations on TABLE/VALUES ---
	"TABLE foo ORDER BY 1 LIMIT 5",
	"VALUES (1), (2), (3) ORDER BY 1 DESC",

	// --- nested subqueries in expression position (B1 placeholders) ---
	"SELECT (SELECT 1)",
	"SELECT a FROM t WHERE a IN (SELECT b FROM u)",
	"SELECT a FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.x = t.x)",
	"SELECT a FROM t WHERE a > ALL (SELECT b FROM u)",
}

// selectRejectCorpus is malformed querySpecification / queryNoWith input Trino
// 481 rejects with a SYNTAX_ERROR (required negative coverage).
var selectRejectCorpus = []string{
	"SELECT",                             // no select list
	"SELECT FROM t",                      // empty select list before FROM
	"SELECT 1 FROM",                      // FROM with no relation
	"SELECT 1,",                          // trailing comma in select list
	"SELECT * FROM",                      // FROM with no relation
	"SELECT 1 WHERE",                     // WHERE with no predicate
	"SELECT 1 WHERE FROM t",              // WHERE with no predicate before FROM
	"SELECT 1 GROUP BY",                  // GROUP BY with no element
	"SELECT 1 GROUP",                     // GROUP without BY
	"SELECT 1 ORDER BY",                  // ORDER BY with no key
	"SELECT 1 ORDER",                     // ORDER without BY
	"SELECT 1 HAVING",                    // HAVING with no predicate
	"SELECT a FROM t LIMIT x",            // LIMIT with a non-count token
	"SELECT a FROM t FETCH FIRST 3 ROWS", // FETCH without ONLY/WITH TIES
	"SELECT a FROM t FETCH FIRST ROW WITH TIES ORDER BY a", // ORDER BY must precede FETCH
	"SELECT DISTINCT FROM t",                               // DISTINCT (reserved) with no select item
	"SELECT a FROM t GROUP BY ROLLUP",                      // ROLLUP without parens
	"SELECT a FROM t GROUP BY GROUPING SETS",               // GROUPING SETS without parens
	"VALUES",                                               // VALUES with no rows
	"TABLE",                                                // TABLE with no name
	"SELECT a a a FROM t",                                  // two aliases (trailing-token error)
	"SELECT 1 garbage extra",                               // trailing junk after the statement
}

// TestSelect_AcceptCorpusParses verifies omni accepts every form in
// selectAcceptCorpus (oracle-free completeness smoke test). The oracle
// differential proves the verdicts actually match Trino.
func TestSelect_AcceptCorpusParses(t *testing.T) {
	for _, sql := range selectAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) != 0 {
				t.Errorf("Parse(%q) should accept, got errors: %v", sql, errs)
			}
		})
	}
}

// TestSelect_RejectCorpusRejected verifies omni rejects every form in
// selectRejectCorpus (required negative coverage).
func TestSelect_RejectCorpusRejected(t *testing.T) {
	for _, sql := range selectRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("Parse(%q) should reject, but accepted", sql)
			}
		})
	}
}

// TestSelect_OracleDifferential is the authoritative accept/reject gate: for
// every form in both corpora, omni's Parse accept/reject must equal Trino 481's
// accept/reject. Skipped when no oracle is reachable.
func TestSelect_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	check := func(t *testing.T, sql string) {
		_, errs := Parse(sql)
		omniAccepts := len(errs) == 0
		trinoAccepts, ok := oracleAccepts(t, o, sql)
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			t.Errorf("MISMATCH sql=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				sql, omniAccepts, errs, trinoAccepts)
		}
	}
	for _, sql := range selectAcceptCorpus {
		t.Run("accept/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
	for _, sql := range selectRejectCorpus {
		t.Run("reject/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
}

// ---------------------------------------------------------------------------
// structural tests (precedence / clause shapes the accept/reject gate can't see)
// ---------------------------------------------------------------------------

// parseOneQuery parses a single-statement query input and returns the *Query,
// failing the test on any parse error or a non-query statement.
func parseOneQuery(t *testing.T, sql string) *Query {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): want 1 statement, got %d", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*QueryStmt)
	if !ok {
		t.Fatalf("Parse(%q): want *QueryStmt, got %T", sql, file.Stmts[0])
	}
	return stmt.Query
}

// querySpec extracts the *QuerySpec body of a query (failing if the body is not
// a plain query specification).
func querySpec(t *testing.T, q *Query) *QuerySpec {
	t.Helper()
	spec, ok := q.Body.(*QuerySpec)
	if !ok {
		t.Fatalf("query body is %T, want *QuerySpec", q.Body)
	}
	return spec
}

func TestSelect_StructureSelectList(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t, "SELECT a, b AS y, * FROM t"))
	if len(spec.Items) != 3 {
		t.Fatalf("want 3 select items, got %d", len(spec.Items))
	}
	if spec.Items[0].Kind != SelectSingle || spec.Items[0].Alias != nil {
		t.Errorf("item 0: want SelectSingle no alias, got kind=%v alias=%v", spec.Items[0].Kind, spec.Items[0].Alias)
	}
	if spec.Items[1].Kind != SelectSingle || spec.Items[1].Alias == nil || spec.Items[1].Alias.Value != "y" {
		t.Errorf("item 1: want SelectSingle alias=y, got kind=%v alias=%v", spec.Items[1].Kind, spec.Items[1].Alias)
	}
	if spec.Items[2].Kind != SelectAll {
		t.Errorf("item 2: want SelectAll, got %v", spec.Items[2].Kind)
	}
}

func TestSelect_StructureSelectAllFrom(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t, "SELECT t.* FROM t"))
	if len(spec.Items) != 1 || spec.Items[0].Kind != SelectAllFrom {
		t.Fatalf("want one SelectAllFrom item, got %+v", spec.Items)
	}
	// The row expression of `t.*` is a ColumnRef to t.
	if _, ok := spec.Items[0].Expr.(*ColumnRef); !ok {
		t.Errorf("SelectAllFrom expr is %T, want *ColumnRef", spec.Items[0].Expr)
	}

	spec = querySpec(t, parseOneQuery(t, "SELECT a.b.* AS (x, y) FROM cat a"))
	if spec.Items[0].Kind != SelectAllFrom {
		t.Fatalf("want SelectAllFrom, got %v", spec.Items[0].Kind)
	}
	if len(spec.Items[0].Aliases) != 2 {
		t.Errorf("want 2 column aliases, got %d", len(spec.Items[0].Aliases))
	}
	// The row expression of `a.b.*` is a Dereference(ColumnRef a, b).
	if _, ok := spec.Items[0].Expr.(*Dereference); !ok {
		t.Errorf("SelectAllFrom expr is %T, want *Dereference", spec.Items[0].Expr)
	}
}

func TestSelect_StructureClauses(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t,
		"SELECT DISTINCT a FROM t WHERE a > 1 GROUP BY a HAVING count(*) > 2 WINDOW w AS (ORDER BY a)"))
	if spec.SetQuantifier != "DISTINCT" {
		t.Errorf("SetQuantifier=%q, want DISTINCT", spec.SetQuantifier)
	}
	if len(spec.From) != 1 {
		t.Errorf("want 1 FROM relation, got %d", len(spec.From))
	}
	if spec.Where == nil {
		t.Error("Where is nil, want a predicate")
	}
	if spec.GroupBy == nil || len(spec.GroupBy.Elements) != 1 {
		t.Errorf("GroupBy = %+v, want one element", spec.GroupBy)
	}
	if spec.Having == nil {
		t.Error("Having is nil, want a predicate")
	}
	if len(spec.Windows) != 1 || spec.Windows[0].Name.Value != "w" {
		t.Errorf("Windows = %+v, want one named 'w'", spec.Windows)
	}
}

func TestSelect_StructureGroupByQuantifier(t *testing.T) {
	// `GROUP BY ALL x` — ALL is the setQuantifier, x the grouping element.
	gb := querySpec(t, parseOneQuery(t, "SELECT a FROM t GROUP BY ALL x")).GroupBy
	if gb.SetQuantifier != "ALL" {
		t.Errorf("`GROUP BY ALL x`: SetQuantifier=%q, want ALL", gb.SetQuantifier)
	}
	if len(gb.Elements) != 1 {
		t.Errorf("`GROUP BY ALL x`: want 1 element, got %d", len(gb.Elements))
	}

	// `GROUP BY ALL` — ALL is a grouping expression (no quantifier).
	gb = querySpec(t, parseOneQuery(t, "SELECT count(*) FROM t GROUP BY ALL")).GroupBy
	if gb.SetQuantifier != "" {
		t.Errorf("`GROUP BY ALL`: SetQuantifier=%q, want empty", gb.SetQuantifier)
	}
	if len(gb.Elements) != 1 || gb.Elements[0].Kind != GroupingExprKind {
		t.Errorf("`GROUP BY ALL`: want 1 expr element, got %+v", gb.Elements)
	}

	// `GROUP BY ALL, x` — ALL is a grouping expression; two elements.
	gb = querySpec(t, parseOneQuery(t, "SELECT count(*) FROM t GROUP BY ALL, x")).GroupBy
	if gb.SetQuantifier != "" {
		t.Errorf("`GROUP BY ALL, x`: SetQuantifier=%q, want empty", gb.SetQuantifier)
	}
	if len(gb.Elements) != 2 {
		t.Errorf("`GROUP BY ALL, x`: want 2 elements, got %d", len(gb.Elements))
	}
}

func TestSelect_StructureGroupingElements(t *testing.T) {
	gb := querySpec(t, parseOneQuery(t,
		"SELECT a FROM t GROUP BY a, ROLLUP (b, c), CUBE (d), GROUPING SETS ((e), ())")).GroupBy
	if len(gb.Elements) != 4 {
		t.Fatalf("want 4 grouping elements, got %d", len(gb.Elements))
	}
	wantKinds := []GroupingElementKind{GroupingExprKind, GroupingRollup, GroupingCube, GroupingSets}
	for i, want := range wantKinds {
		if gb.Elements[i].Kind != want {
			t.Errorf("element %d kind=%v, want %v", i, gb.Elements[i].Kind, want)
		}
	}
	// ROLLUP (b, c) → two expressions.
	if len(gb.Elements[1].Exprs) != 2 {
		t.Errorf("ROLLUP element: want 2 exprs, got %d", len(gb.Elements[1].Exprs))
	}
	// GROUPING SETS ((e), ()) → two sets, second empty.
	if len(gb.Elements[3].Sets) != 2 {
		t.Errorf("GROUPING SETS: want 2 sets, got %d", len(gb.Elements[3].Sets))
	}
	if len(gb.Elements[3].Sets[1]) != 0 {
		t.Errorf("GROUPING SETS second set: want empty, got %d", len(gb.Elements[3].Sets[1]))
	}
}

func TestSelect_StructureLimitFetchOffset(t *testing.T) {
	q := parseOneQuery(t, "SELECT a FROM t ORDER BY a OFFSET 2 ROWS FETCH FIRST 3 ROWS WITH TIES")
	if len(q.OrderBy) != 1 {
		t.Errorf("want 1 ORDER BY item, got %d", len(q.OrderBy))
	}
	if q.Offset == nil || q.Offset.RowKeyword != "ROWS" {
		t.Errorf("Offset = %+v, want ROWS", q.Offset)
	}
	if q.Fetch == nil || q.Fetch.Keyword != "FIRST" || !q.Fetch.WithTies {
		t.Errorf("Fetch = %+v, want FIRST WITH TIES", q.Fetch)
	}

	q = parseOneQuery(t, "SELECT a FROM t LIMIT ALL")
	if q.Limit == nil || !q.Limit.All {
		t.Errorf("Limit = %+v, want LIMIT ALL", q.Limit)
	}

	q = parseOneQuery(t, "SELECT a FROM t LIMIT 10")
	if q.Limit == nil || q.Limit.All || q.Limit.Count == nil {
		t.Errorf("Limit = %+v, want LIMIT 10", q.Limit)
	}
}

func TestSelect_StructureTableAndValues(t *testing.T) {
	q := parseOneQuery(t, "TABLE cat.sch.tbl")
	tq, ok := q.Body.(*TableQuery)
	if !ok {
		t.Fatalf("body is %T, want *TableQuery", q.Body)
	}
	if got := tq.Name.Normalize(); got != "cat.sch.tbl" {
		t.Errorf("TABLE name = %q, want cat.sch.tbl", got)
	}

	q = parseOneQuery(t, "VALUES (1, 'a'), (2, 'b')")
	vq, ok := q.Body.(*ValuesQuery)
	if !ok {
		t.Fatalf("body is %T, want *ValuesQuery", q.Body)
	}
	if len(vq.Rows) != 2 {
		t.Errorf("VALUES rows = %d, want 2", len(vq.Rows))
	}
}

// TestSelect_LocSpans checks the parsed query's Loc covers the whole statement.
func TestSelect_LocSpans(t *testing.T) {
	sql := "SELECT a, b FROM t WHERE a > 1"
	q := parseOneQuery(t, sql)
	if q.Loc.Start != 0 {
		t.Errorf("query Loc.Start = %d, want 0", q.Loc.Start)
	}
	if q.Loc.End != len(sql) {
		t.Errorf("query Loc.End = %d, want %d", q.Loc.End, len(sql))
	}
	// every select item carries a valid span.
	for i, it := range querySpec(t, q).Items {
		if !it.Loc.IsValid() {
			t.Errorf("select item %d has invalid Loc %+v", i, it.Loc)
		}
		_ = ast.Loc(it.Loc)
	}
}
