package analysis

import "testing"

// tableSig is a compact, order-insensitive form of TableAccess used in tests.
type tableSig struct {
	Catalog string
	Schema  string
	Table   string
	Alias   string
}

func toSigs(tables []TableAccess) []tableSig {
	out := make([]tableSig, len(tables))
	for i, tbl := range tables {
		out[i] = tableSig{Catalog: tbl.Catalog, Schema: tbl.Schema, Table: tbl.Table, Alias: tbl.Alias}
	}
	return out
}

func containsSig(haystack []tableSig, needle tableSig) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestGetQuerySpan_Basics(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		wantType    QueryType
		wantTables  []tableSig
		wantResults []string // output column names, in order
		wantCTEs    []string
	}{
		{
			name:        "simple select",
			sql:         "SELECT a FROM t",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "t"}},
			wantResults: []string{"a"},
		},
		{
			name:        "qualified schema.table",
			sql:         "SELECT * FROM s.t",
			wantType:    Select,
			wantTables:  []tableSig{{Schema: "s", Table: "t"}},
			wantResults: []string{"*"},
		},
		{
			name:        "three-part catalog.schema.table",
			sql:         "SELECT * FROM cat.sch.tbl",
			wantType:    Select,
			wantTables:  []tableSig{{Catalog: "cat", Schema: "sch", Table: "tbl"}},
			wantResults: []string{"*"},
		},
		{
			name:        "table alias",
			sql:         "SELECT a FROM t AS x",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "t", Alias: "x"}},
			wantResults: []string{"a"},
		},
		{
			name:        "table alias without AS",
			sql:         "SELECT a FROM t x",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "t", Alias: "x"}},
			wantResults: []string{"a"},
		},
		{
			name:        "select item alias",
			sql:         "SELECT a AS col_a, b FROM t",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "t"}},
			wantResults: []string{"col_a", "b"},
		},
		{
			name:     "join of two tables",
			sql:      "SELECT a, b FROM t1 JOIN t2 ON t1.id = t2.id",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"a", "b"},
		},
		{
			name:     "inner and left join with qualified table",
			sql:      "SELECT a FROM t1 INNER JOIN t2 ON t1.id = t2.id LEFT JOIN s.t3 ON t1.x = t3.x",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
				{Schema: "s", Table: "t3"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "cross join",
			sql:      "SELECT * FROM t1 CROSS JOIN t2",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"*"},
		},
		{
			name:     "comma join (implicit cross)",
			sql:      "SELECT * FROM t1, t2",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"*"},
		},
		{
			name:        "cte exclusion",
			sql:         "WITH c AS (SELECT * FROM real_t) SELECT * FROM c",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "real_t"}},
			wantResults: []string{"*"},
			wantCTEs:    []string{"c"},
		},
		{
			name: "multiple ctes chained",
			sql: `WITH a AS (SELECT id FROM t1),
			      b AS (SELECT id FROM a JOIN t2 ON a.id = t2.id)
			      SELECT id FROM b`,
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"id"},
			wantCTEs:    []string{"a", "b"},
		},
		{
			name:     "subquery in where IN",
			sql:      "SELECT a FROM t WHERE b IN (SELECT b FROM other)",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "other"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "scalar subquery in select",
			sql:      "SELECT a, (SELECT max(x) FROM agg) FROM t",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "agg"},
			},
		},
		{
			name:     "exists subquery",
			sql:      "SELECT a FROM t WHERE EXISTS (SELECT 1 FROM other WHERE other.id = t.id)",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "other"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "from subquery (parsed into real Query)",
			sql:      "SELECT x.a FROM (SELECT a FROM inner_t) AS x",
			wantType: Select,
			// The derived-table alias `x` applies to the derived table, not to
			// its underlying base table — so inner_t is recorded unaliased.
			wantTables: []tableSig{
				{Table: "inner_t"},
			},
		},
		{
			name:     "union two tables (left wins for results)",
			sql:      "SELECT a FROM t1 UNION SELECT a FROM t2",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "union all three tables",
			sql:      "SELECT a FROM t1 UNION ALL SELECT a FROM t2 UNION ALL SELECT a FROM s.t3",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
				{Schema: "s", Table: "t3"},
			},
			wantResults: []string{"a"},
		},
		{
			name:        "table query shorthand",
			sql:         "TABLE orders",
			wantType:    Select,
			wantTables:  []tableSig{{Table: "orders"}},
			wantResults: nil,
		},
		{
			name:        "unnest in from",
			sql:         "SELECT v FROM UNNEST(ARRAY[1,2,3]) AS u(v)",
			wantType:    Select,
			wantTables:  nil, // UNNEST is not a base table
			wantResults: []string{"v"},
		},
		{
			name:     "lateral subquery",
			sql:      "SELECT t.a, l.b FROM t, LATERAL (SELECT b FROM u WHERE u.id = t.id) AS l",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "u"},
			},
		},
		{
			name:     "order by limit do not add tables",
			sql:      "SELECT a FROM t ORDER BY a LIMIT 10",
			wantType: Select,
			wantTables: []tableSig{
				{Table: "t"},
			},
			wantResults: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span, err := GetQuerySpan(tt.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan returned error: %v", err)
			}
			if span.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", span.Type, tt.wantType)
			}

			gotSigs := toSigs(span.AccessTables)
			if len(gotSigs) != len(tt.wantTables) {
				t.Errorf("AccessTables len = %d (%v), want %d (%v)",
					len(gotSigs), gotSigs, len(tt.wantTables), tt.wantTables)
			}
			for _, want := range tt.wantTables {
				if !containsSig(gotSigs, want) {
					t.Errorf("AccessTables missing %+v (got %+v)", want, gotSigs)
				}
			}

			if tt.wantResults != nil {
				if len(span.Results) != len(tt.wantResults) {
					t.Fatalf("Results len = %d, want %d (got %+v)",
						len(span.Results), len(tt.wantResults), span.Results)
				}
				for i, want := range tt.wantResults {
					if span.Results[i].Name != want {
						t.Errorf("Results[%d].Name = %q, want %q", i, span.Results[i].Name, want)
					}
				}
			}

			if tt.wantCTEs != nil {
				if len(span.CTEs) != len(tt.wantCTEs) {
					t.Fatalf("CTEs = %v, want %v", span.CTEs, tt.wantCTEs)
				}
				for i, want := range tt.wantCTEs {
					if span.CTEs[i] != want {
						t.Errorf("CTEs[%d] = %q, want %q", i, span.CTEs[i], want)
					}
				}
			}
		})
	}
}

func TestGetQuerySpan_NonQuery(t *testing.T) {
	tests := []struct {
		sql      string
		wantType QueryType
	}{
		{"CREATE TABLE t (id INT)", DDL},
		{"SHOW TABLES", SelectInfoSchema},
		{"INSERT INTO t VALUES (1)", DML},
		{"EXPLAIN SELECT 1", Explain},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			span, err := GetQuerySpan(tt.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan returned error: %v", err)
			}
			if span.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", span.Type, tt.wantType)
			}
		})
	}
}

func TestGetQuerySpan_Empty(t *testing.T) {
	span, err := GetQuerySpan("")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if span.Type != Unknown {
		t.Errorf("Type = %v, want Unknown", span.Type)
	}
	if len(span.AccessTables) != 0 {
		t.Errorf("AccessTables = %v, want empty", span.AccessTables)
	}
	if len(span.Results) != 0 {
		t.Errorf("Results = %v, want empty", span.Results)
	}
}

func TestGetQuerySpan_DedupTables(t *testing.T) {
	// Same unaliased table referenced twice — should appear once.
	span, err := GetQuerySpan("SELECT a FROM t WHERE b IN (SELECT b FROM t)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.AccessTables) != 1 {
		t.Errorf("AccessTables = %v, want one entry for t", span.AccessTables)
	}
}

func TestGetQuerySpan_SameTableDifferentAliases(t *testing.T) {
	// Self-join — two aliases of the same physical table should each appear.
	span, err := GetQuerySpan("SELECT * FROM t AS a JOIN t AS b ON a.id = b.id")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.AccessTables) != 2 {
		t.Errorf("AccessTables = %+v, want two (a, b)", span.AccessTables)
	}
	sigs := toSigs(span.AccessTables)
	if !containsSig(sigs, tableSig{Table: "t", Alias: "a"}) {
		t.Errorf("missing alias a: %+v", sigs)
	}
	if !containsSig(sigs, tableSig{Table: "t", Alias: "b"}) {
		t.Errorf("missing alias b: %+v", sigs)
	}
}

func TestGetQuerySpan_SourceColumns(t *testing.T) {
	span, err := GetQuerySpan("SELECT a + b AS total FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results len = %d, want 1 (%+v)", len(span.Results), span.Results)
	}
	r := span.Results[0]
	if r.Name != "total" {
		t.Errorf("Results[0].Name = %q, want total", r.Name)
	}
	if len(r.SourceColumns) != 2 {
		t.Fatalf("SourceColumns len = %d, want 2 (%+v)", len(r.SourceColumns), r.SourceColumns)
	}
	want := []ColumnRef{{Column: "a"}, {Column: "b"}}
	for i, w := range want {
		if r.SourceColumns[i] != w {
			t.Errorf("SourceColumns[%d] = %+v, want %+v", i, r.SourceColumns[i], w)
		}
	}
}

func TestGetQuerySpan_DirectColumnBesideSubquery(t *testing.T) {
	// A select item whose expression has a direct column reference alongside a
	// subquery (`a IN (SELECT ...)`) must surface the direct column `a` as a
	// source column, and the subquery's table must still be discovered.
	span, err := GetQuerySpan("SELECT (a IN (SELECT x FROM u)) AS flag FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 || span.Results[0].Name != "flag" {
		t.Fatalf("Results = %+v, want one named flag", span.Results)
	}
	found := false
	for _, c := range span.Results[0].SourceColumns {
		if c == (ColumnRef{Column: "a"}) {
			found = true
		}
	}
	if !found {
		t.Errorf("SourceColumns = %+v, want to contain {Column:a}", span.Results[0].SourceColumns)
	}
	sigs := toSigs(span.AccessTables)
	if !containsSig(sigs, tableSig{Table: "t"}) || !containsSig(sigs, tableSig{Table: "u"}) {
		t.Errorf("AccessTables = %+v, want [t, u]", sigs)
	}
}

func TestGetQuerySpan_QualifiedSourceColumn(t *testing.T) {
	// A dotted column reference t.a is a Dereference over a ColumnRef in the
	// trino AST; the extractor must flatten it back into a qualified ColumnRef.
	span, err := GetQuerySpan("SELECT t.a FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results len = %d, want 1", len(span.Results))
	}
	r := span.Results[0]
	if r.Name != "a" {
		t.Errorf("Results[0].Name = %q, want a", r.Name)
	}
	if len(r.SourceColumns) != 1 {
		t.Fatalf("SourceColumns = %+v, want 1", r.SourceColumns)
	}
	if got := r.SourceColumns[0]; got != (ColumnRef{Table: "t", Column: "a"}) {
		t.Errorf("SourceColumns[0] = %+v, want {Table:t Column:a}", got)
	}
}

func TestGetQuerySpan_CTEShadowsTable(t *testing.T) {
	// A CTE named real_t shadows any physical table of the same name; the outer
	// reference resolves to the CTE, so AccessTables holds only the CTE body's
	// tables.
	span, err := GetQuerySpan(`
		WITH real_t AS (SELECT id FROM underlying)
		SELECT id FROM real_t
	`)
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	sigs := toSigs(span.AccessTables)
	if len(sigs) != 1 || sigs[0].Table != "underlying" {
		t.Errorf("AccessTables = %+v, want [underlying]", sigs)
	}
}

func TestGetQuerySpan_JSONFunctionSourceColumns(t *testing.T) {
	// A SQL/JSON function input is a real column-reference position; the
	// extractor must surface it as a source column.
	span, err := GetQuerySpan("SELECT JSON_VALUE(t.data, 'lax $.x') AS v FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results len = %d, want 1 (%+v)", len(span.Results), span.Results)
	}
	r := span.Results[0]
	if r.Name != "v" {
		t.Errorf("Results[0].Name = %q, want v", r.Name)
	}
	want := ColumnRef{Table: "t", Column: "data"}
	found := false
	for _, got := range r.SourceColumns {
		if got == want {
			found = true
		}
	}
	if !found {
		t.Errorf("SourceColumns = %+v, want to contain %+v", r.SourceColumns, want)
	}
}

func TestGetQuerySpan_SubqueryInheritsOuterCTEScope(t *testing.T) {
	// A subquery embedded in WHERE references an outer-scope CTE (c); because the
	// re-parsed subquery is walked with the current CTE scope as parent, c is
	// recognized as a CTE and excluded — only the CTE body's base table and the
	// outer FROM table appear.
	span, err := GetQuerySpan("WITH c AS (SELECT y FROM base) SELECT a FROM t WHERE x IN (SELECT y FROM c)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	sigs := toSigs(span.AccessTables)
	if len(sigs) != 2 {
		t.Fatalf("AccessTables = %+v, want [base, t]", sigs)
	}
	if !containsSig(sigs, tableSig{Table: "base"}) || !containsSig(sigs, tableSig{Table: "t"}) {
		t.Errorf("AccessTables = %+v, want [base, t] (CTE c excluded)", sigs)
	}
}

func TestGetQuerySpan_DeeplyNestedFromSubquery(t *testing.T) {
	span, err := GetQuerySpan("SELECT a FROM (SELECT b FROM (SELECT c FROM deep) AS i1) AS i2")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	sigs := toSigs(span.AccessTables)
	if len(sigs) != 1 || sigs[0].Table != "deep" {
		t.Errorf("AccessTables = %+v, want [deep]", sigs)
	}
}

func TestGetQuerySpan_AliasedCTEReferenceExcluded(t *testing.T) {
	// A CTE referenced WITH an alias (`FROM c AS x`) is still a CTE and must be
	// excluded from AccessTables — only the CTE body's base table appears.
	span, err := GetQuerySpan("WITH c AS (SELECT a FROM real_t) SELECT x.a FROM c AS x")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	sigs := toSigs(span.AccessTables)
	if len(sigs) != 1 || sigs[0].Table != "real_t" {
		t.Errorf("AccessTables = %+v, want only [real_t] (CTE c excluded)", sigs)
	}
}

func TestGetQuerySpan_TableQueryResultStar(t *testing.T) {
	// `TABLE orders` is shorthand for `SELECT * FROM orders`; its result is a
	// single star column.
	span, err := GetQuerySpan("TABLE orders")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 || span.Results[0].Name != "*" {
		t.Errorf("Results = %+v, want [*]", span.Results)
	}
}

func TestGetQuerySpan_JoinUsingPredicateColumns(t *testing.T) {
	// JOIN ... USING (col) columns are join keys and must appear in
	// PredicateColumns for row-level masking parity.
	span, err := GetQuerySpan("SELECT * FROM orders JOIN lineitem USING (orderkey)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	found := false
	for _, c := range span.PredicateColumns {
		if c == (ColumnRef{Column: "orderkey"}) {
			found = true
		}
	}
	if !found {
		t.Errorf("PredicateColumns = %+v, want to contain {Column:orderkey}", span.PredicateColumns)
	}
}

func TestGetQuerySpan_LambdaParamNotALineageColumn(t *testing.T) {
	// A lambda-bound parameter (x in `x -> x + 1`) is NOT a table column and
	// must not appear as a source column; an outer column referenced inside the
	// lambda body still should.
	span, err := GetQuerySpan("SELECT transform(xs, x -> x + outer_col) AS ys FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results = %+v, want 1", span.Results)
	}
	for _, c := range span.Results[0].SourceColumns {
		if c == (ColumnRef{Column: "x"}) {
			t.Errorf("SourceColumns = %+v, must NOT contain lambda param {Column:x}", span.Results[0].SourceColumns)
		}
	}
	// The free column `xs` (and ideally `outer_col`) should still be present.
	foundXS := false
	for _, c := range span.Results[0].SourceColumns {
		if c == (ColumnRef{Column: "xs"}) {
			foundXS = true
		}
	}
	if !foundXS {
		t.Errorf("SourceColumns = %+v, want to contain free column {Column:xs}", span.Results[0].SourceColumns)
	}
}

func TestGetQuerySpan_LambdaParamDereferenceNotALineageColumn(t *testing.T) {
	// A field access rooted at a lambda-bound parameter (`x.price` in
	// `x -> x.price + …`) is access on the bound variable, not a table column,
	// and must not appear as a source column; the free column `outer_col` still
	// should.
	span, err := GetQuerySpan("SELECT transform(items, x -> x.price + outer_col) AS prices FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results = %+v, want 1", span.Results)
	}
	for _, c := range span.Results[0].SourceColumns {
		if c.Table == "x" {
			t.Errorf("SourceColumns = %+v, must NOT contain a ref rooted at lambda param x", span.Results[0].SourceColumns)
		}
	}
	foundFree := false
	foundItems := false
	for _, c := range span.Results[0].SourceColumns {
		if c == (ColumnRef{Column: "outer_col"}) {
			foundFree = true
		}
		if c == (ColumnRef{Column: "items"}) {
			foundItems = true
		}
	}
	if !foundFree {
		t.Errorf("SourceColumns = %+v, want free column {Column:outer_col}", span.Results[0].SourceColumns)
	}
	if !foundItems {
		t.Errorf("SourceColumns = %+v, want collection arg {Column:items}", span.Results[0].SourceColumns)
	}
}

func TestGetQuerySpan_GroupingArgsAreColumns(t *testing.T) {
	// GROUPING(a) references column a; it must be surfaced as a source column.
	span, err := GetQuerySpan("SELECT GROUPING(a) AS g FROM t GROUP BY ROLLUP(a)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results = %+v, want 1", span.Results)
	}
	found := false
	for _, c := range span.Results[0].SourceColumns {
		if c == (ColumnRef{Column: "a"}) {
			found = true
		}
	}
	if !found {
		t.Errorf("SourceColumns = %+v, want to contain {Column:a}", span.Results[0].SourceColumns)
	}
}

func TestGetQuerySpan_TableFunctionArgColumns(t *testing.T) {
	// A table-function table argument's PARTITION BY and ORDER BY columns are
	// real column references and must both be captured.
	span, err := GetQuerySpan("SELECT * FROM TABLE(my_function(input => TABLE(orders) PARTITION BY orderstatus ORDER BY orderdate))")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if !containsSig(toSigs(span.AccessTables), tableSig{Table: "orders"}) {
		t.Errorf("AccessTables = %+v, want to contain orders", toSigs(span.AccessTables))
	}
	for _, want := range []ColumnRef{{Column: "orderstatus"}, {Column: "orderdate"}} {
		found := false
		for _, c := range span.PredicateColumns {
			if c == want {
				found = true
			}
		}
		if !found {
			t.Errorf("PredicateColumns = %+v, want to contain %+v", span.PredicateColumns, want)
		}
	}
}

func TestGetQuerySpan_NamedWindowColumnsCaptured(t *testing.T) {
	// The PARTITION BY / ORDER BY columns of a named WINDOW definition are real
	// column references and must be captured (as predicate columns) so masking
	// sees them — even though attributing them to the specific OVER-w result
	// column is out of best-effort scope.
	span, err := GetQuerySpan("SELECT row_number() OVER w AS rn FROM t WINDOW w AS (PARTITION BY y ORDER BY z)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	for _, want := range []ColumnRef{{Column: "y"}, {Column: "z"}} {
		found := false
		for _, c := range span.PredicateColumns {
			if c == want {
				found = true
			}
		}
		if !found {
			t.Errorf("PredicateColumns = %+v, want to contain %+v", span.PredicateColumns, want)
		}
	}
}

func TestClassify_ExplainAnalyzeMutating(t *testing.T) {
	// EXPLAIN ANALYZE executes its inner statement, so EXPLAIN ANALYZE of a
	// mutating statement is NOT read-only — it must classify by the inner
	// statement's nature (oracle-confirmed: Trino 481 runs the INSERT).
	tests := []struct {
		sql  string
		want QueryType
	}{
		{"EXPLAIN ANALYZE SELECT a FROM t", Explain},
		{"EXPLAIN ANALYZE VERBOSE SELECT a FROM t", Explain},
		{"EXPLAIN SELECT a FROM t", Explain},
		{"EXPLAIN ANALYZE INSERT INTO t VALUES (1)", DML},
		{"EXPLAIN ANALYZE DELETE FROM t WHERE a = 1", DML},
		{"EXPLAIN ANALYZE UPDATE t SET a = 1", DML},
		// Plain EXPLAIN (no ANALYZE) never executes, so it stays read-only even
		// for a DML inner statement.
		{"EXPLAIN INSERT INTO t VALUES (1)", Explain},
		// EXPLAIN ANALYZE of an info-schema read stays read-only/info-schema.
		{"EXPLAIN ANALYZE SELECT * FROM information_schema.tables", SelectInfoSchema},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			if got := Classify(tt.sql); got != tt.want {
				t.Errorf("Classify(%q) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestGetQuerySpan_PredicateColumns(t *testing.T) {
	// WHERE / JOIN-ON columns are collected into PredicateColumns (used for
	// row-level masking parity with the legacy extractor).
	span, err := GetQuerySpan("SELECT a FROM t1 JOIN t2 ON t1.id = t2.id WHERE t1.region = 'x'")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	want := []ColumnRef{
		{Table: "t1", Column: "id"},
		{Table: "t2", Column: "id"},
		{Table: "t1", Column: "region"},
	}
	for _, w := range want {
		found := false
		for _, got := range span.PredicateColumns {
			if got == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PredicateColumns missing %+v (got %+v)", w, span.PredicateColumns)
		}
	}
}

func TestGetQuerySpan_NoPanicOnVariedInput(t *testing.T) {
	// Robustness sweep (correctness-protocol.md item 3): GetQuerySpan must
	// terminate without panicking on a broad range of expression-heavy,
	// nested, exotic, and malformed inputs.
	corpus := []string{
		"",
		";",
		"SELECT",
		"SELECT 1",
		"SELECT CASE WHEN a > 1 THEN b ELSE c END FROM t",
		"SELECT CAST(a AS varchar), TRY_CAST(b AS bigint) FROM t",
		"SELECT EXTRACT(YEAR FROM ts), SUBSTRING(s FROM 1 FOR 3) FROM t",
		"SELECT ARRAY[a, b, c], ROW(x, y) FROM t",
		"SELECT f(x) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) FROM t",
		"SELECT a[1], m['k'] FROM t",
		"SELECT JSON_QUERY(payload, 'lax $.items') FROM t",
		"SELECT JSON_OBJECT(KEY k VALUE v) FROM t",
		"SELECT * FROM t1 NATURAL JOIN t2",
		"SELECT * FROM UNNEST(ARRAY[1,2]) WITH ORDINALITY AS u(v, ord)",
		"SELECT * FROM t1, LATERAL (SELECT * FROM t2 WHERE t2.id = t1.id) l",
		"SELECT * FROM TABLE(my_func(input => TABLE(orders) PARTITION BY region))",
		"WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM r WHERE n < 5) SELECT * FROM r",
		"SELECT a FROM t WHERE x = ANY (SELECT y FROM u)",
		"SELECT a FROM t GROUP BY ROLLUP (a, b), CUBE (c) HAVING count(*) > 1",
		"SELECT a FROM t WHERE a BETWEEN 1 AND 10 AND b LIKE '%x%' ESCAPE '\\'",
		"VALUES (1, 'a'), (2, 'b')",
		"TABLE orders",
		"((SELECT 1))",
		"SELECT a FROM (SELECT b FROM (SELECT c FROM deep) AS inner1) AS inner2",
		"garbage tokens that do not parse @@@ %%%",
		"INSERT INTO t SELECT * FROM s",
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a",
	}
	for _, sql := range corpus {
		sql := sql
		t.Run(truncate(sql), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("GetQuerySpan(%q) panicked: %v", sql, r)
				}
			}()
			if _, err := GetQuerySpan(sql); err != nil {
				t.Fatalf("GetQuerySpan(%q) returned error: %v", sql, err)
			}
		})
	}
}

func truncate(s string) string {
	if len(s) > 48 {
		return s[:48]
	}
	if s == "" {
		return "(empty)"
	}
	return s
}

func TestGetQuerySpan_ToleratesParseError(t *testing.T) {
	// A statement whose body fails to parse must not panic; GetQuerySpan returns
	// whatever was classified/extracted.
	span, err := GetQuerySpan("SELECT a FROM")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if span.Type != Select {
		t.Errorf("Type = %v, want Select", span.Type)
	}
}
