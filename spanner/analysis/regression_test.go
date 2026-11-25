package analysis

import (
	"sort"
	"testing"
)

// This file holds one permanent regression test per review finding (the Codex
// lens, review-gate.md), so a future change that reintroduces the defect fails
// loudly. Each test names the finding it pins.

// Finding #2: Spanner 3-part db.schema.table puts db in Database, not Catalog
// (Spanner has no project/catalog layer).
func TestRegression_Spanner3PartDatabase(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM db.sch.t")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	ta := span.AccessTables[0]
	if ta.Database != "db" || ta.Schema != "sch" || ta.Table != "t" || ta.Catalog != "" {
		t.Errorf("Spanner db.sch.t = {Catalog:%q Database:%q Schema:%q Table:%q}, want {\"\" db sch t}",
			ta.Catalog, ta.Database, ta.Schema, ta.Table)
	}
}

// Finding #3: DROP/ALTER of non-table objects (SCHEMA/DATABASE/INDEX) must NOT
// produce fake AccessTables entries.
func TestRegression_DDLNonTableObjects(t *testing.T) {
	tests := []struct {
		sql  string
		want int // expected AccessTables count
	}{
		{"DROP SCHEMA s", 0},
		{"DROP INDEX idx", 0},
		{"ALTER SCHEMA s SET OPTIONS (default_leader = 'us')", 0},
		{"DROP TABLE realtbl", 1}, // a real table IS recorded
		{"DROP VIEW realview", 1}, // a view IS recorded
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			if len(span.AccessTables) != tc.want {
				t.Errorf("AccessTables = %v (%d), want %d", tableNames(span), len(span.AccessTables), tc.want)
			}
		})
	}
}

// Finding #4: CREATE TABLE source-table references (LIKE / CLONE / COPY /
// INTERLEAVE parent / FOREIGN KEY references) must appear in AccessTables.
func TestRegression_CreateTableSourceRefs(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantTables []string
	}{
		{
			// The merged parser requires the ON DELETE action on INTERLEAVE (it
			// rejects the bare `INTERLEAVE IN PARENT parent`, which the Spanner
			// emulator actually ACCEPTS — a parser-grammar over-reject flagged to the
			// ledger, owned by parser-ddl-spanner, NOT this node). Use the parseable
			// form here; the analysis records the parent either way.
			name:       "interleave parent",
			sql:        "CREATE TABLE child (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT parent ON DELETE CASCADE",
			wantTables: []string{"child", "parent"},
		},
		{
			name:       "foreign key reference",
			sql:        "CREATE TABLE child (id INT64, pid INT64, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES parent (id)) PRIMARY KEY (id)",
			wantTables: []string{"child", "parent"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			got := tableNames(span)
			if !eqStrings(got, tc.wantTables) {
				t.Errorf("tables = %v, want %v", got, tc.wantTables)
			}
		})
	}
}

// Finding #6: TVF arguments contribute predicate columns (the doc-promised
// behavior), not just subquery discovery.
func TestRegression_TVFArgsPredicateColumns(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM my_tvf(t.id) AS x")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	got := predicateColumnNames(span)
	if !eqStrings(got, []string{"id"}) {
		t.Errorf("predicate columns = %v, want [id] (TVF arg t.id)", got)
	}
}

// Finding #7: WITH-expression local variables are excluded from the body's
// column refs, like lambda parameters.
func TestRegression_WithExprVarScoping(t *testing.T) {
	// `x` is a WITH-expr local binding; `y` is a table column. The body `x + y`
	// must contribute only `y`. The binding `x AS y_plus_1`'s value `y + 1` also
	// references `y` (a real column).
	span, err := GetQuerySpan("SELECT WITH(x AS y + 1, x + z) AS v FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("want 1 result, got %d", len(span.Results))
	}
	var cols []string
	for _, sc := range span.Results[0].SourceColumns {
		cols = append(cols, sc.Column)
	}
	sort.Strings(cols)
	// y (from the binding value), z (free in the body) — but NOT x (bound name).
	if !eqStrings(cols, []string{"y", "z"}) {
		t.Errorf("source columns = %v, want [y z] (WITH-expr var x excluded)", cols)
	}
}

// Re-review finding #1: a later WITH-expr binding may reference an earlier
// binding, which must also be treated as a local variable (not a table column).
func TestRegression_WithExprProgressiveScoping(t *testing.T) {
	// a AS c1; b AS a + c2 (a is the earlier local binding, c2 a table column);
	// body b + c3 (b is local, c3 a table column). Source columns: c1, c2, c3.
	span, err := GetQuerySpan("SELECT WITH(a AS c1, b AS a + c2, b + c3) AS v FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("want 1 result, got %d", len(span.Results))
	}
	var cols []string
	for _, sc := range span.Results[0].SourceColumns {
		cols = append(cols, sc.Column)
	}
	sort.Strings(cols)
	if !eqStrings(cols, []string{"c1", "c2", "c3"}) {
		t.Errorf("source columns = %v, want [c1 c2 c3] (locals a, b excluded)", cols)
	}
}

// Round-3 finding #1: non-recursive CTE visibility is sequential — a forward
// reference to a LATER sibling resolves to a base table, not the CTE.
func TestRegression_NonRecursiveCTEForwardRef(t *testing.T) {
	// CTE `a`'s body references `b`, which is a LATER sibling CTE. In
	// non-recursive GoogleSQL `b` is not yet visible inside `a`, so it is the base
	// table `b`. Tables: b (the real base table referenced from a) + src (from b's
	// body). Neither `a` nor the CTE-`b` reference is a base table.
	sql := "WITH a AS (SELECT * FROM b), b AS (SELECT * FROM src) SELECT * FROM a"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := tableNames(span); !eqStrings(got, []string{"b", "src"}) {
		t.Errorf("tables = %v, want [b src] (forward ref to later sibling b is a base table)", got)
	}
}

// TestRegression_RecursiveCTESelfRef confirms a RECURSIVE CTE's self-reference
// is filtered (the name is visible inside its own body).
func TestRegression_RecursiveCTESelfRef(t *testing.T) {
	sql := "WITH RECURSIVE c AS ((SELECT 1 AS n) UNION ALL (SELECT n + 1 FROM c WHERE n < 3)) SELECT * FROM c"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	// `c` references only itself; no base tables.
	if len(span.AccessTables) != 0 {
		t.Errorf("tables = %v, want [] (recursive self-ref c is not a base table)", tableNames(span))
	}
}

// Round-3 finding #2: a correlated array/field path off an earlier FROM alias
// (`FROM T t, t.arr`) is NOT a base table.
func TestRegression_CorrelatedArrayPath(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM T1 t1, t1.array_column")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	// Only T1 is a base table; t1.array_column is a correlated unnest of its alias.
	if got := tableNames(span); !eqStrings(got, []string{"T1"}) {
		t.Errorf("tables = %v, want [T1] (t1.array_column is a correlated path, not a table)", got)
	}
	// A genuine 2-part table (dataset.table, root not a FROM alias) is still a
	// base table — the alias check must not over-suppress real qualified tables.
	span2, _ := GetQuerySpan("SELECT * FROM ds.users")
	if got := tableNames(span2); !eqStrings(got, []string{"ds.users"}) {
		t.Errorf("tables = %v, want [ds.users] (real qualified table)", got)
	}
}

// Round-4 finding: a correlated array path off a DERIVED (subquery) alias or a
// CHAINED correlated alias is also not a base table.
func TestRegression_CorrelatedPathDerivedAndChained(t *testing.T) {
	// Off a subquery alias: `s.arr` correlates the derived relation `s`.
	span, err := GetQuerySpan("SELECT * FROM (SELECT [1, 2] AS arr) s, s.arr")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.AccessTables) != 0 {
		t.Errorf("derived-alias correlated path: tables = %v, want [] (subquery has no base table; s.arr is correlated)", tableNames(span))
	}

	// Chained: `t.arr a` registers alias `a`, so `a.child` correlates off it.
	span2, err := GetQuerySpan("SELECT * FROM T t, t.arr a, a.child")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := tableNames(span2); !eqStrings(got, []string{"T"}) {
		t.Errorf("chained correlated path: tables = %v, want [T] (t.arr and a.child are correlated)", got)
	}
}
