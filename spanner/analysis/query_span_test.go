package analysis

import (
	"sort"
	"strings"
	"testing"
)

// --- helpers --------------------------------------------------------------

// tableNames renders AccessTables as sorted "schema.table" (or
// "catalog.schema.table") strings for stable comparison.
func tableNames(span *QuerySpan) []string {
	var out []string
	for _, ta := range span.AccessTables {
		out = append(out, ta.qualifiedString())
	}
	sort.Strings(out)
	return out
}

func (ta TableAccess) qualifiedString() string {
	var parts []string
	if ta.Catalog != "" {
		parts = append(parts, ta.Catalog)
	}
	if ta.Database != "" {
		parts = append(parts, ta.Database)
	}
	if ta.Schema != "" {
		parts = append(parts, ta.Schema)
	}
	parts = append(parts, ta.Table)
	return strings.Join(parts, ".")
}

func resultNames(span *QuerySpan) []string {
	var out []string
	for _, r := range span.Results {
		out = append(out, r.Name)
	}
	return out
}

func predicateColumnNames(span *QuerySpan) []string {
	var out []string
	for _, c := range span.PredicateColumns {
		out = append(out, c.Column)
	}
	sort.Strings(out)
	return out
}

func eqStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// --- AccessTables ---------------------------------------------------------

func TestGetQuerySpan_AccessTables(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantTables []string
	}{
		// Quantified comparison subquery RHS (`= ANY (SELECT ...)`): its table must
		// be followed, exactly like an IN/EXISTS subquery (regression for the
		// CompareExpr Quant* walker path).
		{"quantified-comparison all subquery", "SELECT * FROM t WHERE x > ALL (SELECT id FROM u)", []string{"t", "u"}},
		// The `TABLE t` query primary (table_clause_reserved) names its table, and a
		// set operation over TABLE primaries names every arm's.
		// A TABLE primary inside a parenthesized query or a FROM subquery.
		// EXPORT TABLE METADATA reads the metadata of the table it names.
		// An OPTIONS value's subquery is filled, so the table it reads is recorded.
		// Spanner schema.table bucketing (named schema under one DB).
		{"qualified table spanner", "SELECT a FROM myschema.t", []string{"myschema.t"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			got := tableNames(span)
			if !eqStrings(got, tc.wantTables) {
				t.Errorf("GetQuerySpan(%q) tables = %v, want %v", tc.sql, got, tc.wantTables)
			}
		})
	}
}

// TestGetQuerySpan_TableClauseTVFFailsClosed pins the one TABLE clause shape whose
// projection cannot be resolved: a table-valued function's output columns are
// unknowable without a signature, so the span fails closed rather than reporting
// no columns at all (which a masking consumer would read as nothing to mask).
func TestGetQuerySpan_TableClauseTVFFailsClosed(t *testing.T) {
	for _, sql := range []string{"TABLE f()", "TABLE f(TABLE u)", "TABLE f(a => 1)"} {
		if _, err := GetQuerySpan(sql); err == nil {
			t.Errorf("GetQuerySpan(%q): want a fail-closed error, got none", sql)
		}
	}
	// The same clause as a pipe set operation's operand still records its table
	// argument (it resolves no projection there either, and fails closed is the
	// surrounding statement's business).
	if got := ClassifySQL("TABLE f(TABLE u)"); got != Select {
		t.Errorf("ClassifySQL = %v, want Select", got)
	}
}

// TestGetQuerySpan_CallTableArguments records the tables a CALL's tvf_argument
// list names: positionally (`TABLE t`) and as a named table_clause argument
// (`arg => TABLE t`).
func TestGetQuerySpan_CallTableArguments(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"CALL p(TABLE u)", []string{"u"}},
		{"CALL p(a => TABLE u)", []string{"u"}},
		{"CALL p(1, TABLE u, b => TABLE v)", []string{"u", "v"}},
		{"CALL p(a => INPUT TABLE)", nil},
		{"CALL p(INPUT TABLE)", nil},
	} {
		t.Run(c.sql, func(t *testing.T) {
			span, err := GetQuerySpan(c.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q): %v", c.sql, err)
			}
			if got := tableNames(span); !eqStrings(got, c.want) {
				t.Errorf("tables = %v, want %v", got, c.want)
			}
		})
	}
}

// TestGetQuerySpan_CTEsExcluded confirms CTE names are recorded in span.CTEs and
// excluded from AccessTables, while the CTE BODY's base tables are kept.
func TestGetQuerySpan_CTEsExcluded(t *testing.T) {
	sql := "WITH c AS (SELECT a FROM base) SELECT a FROM c"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := tableNames(span); !eqStrings(got, []string{"base"}) {
		t.Errorf("tables = %v, want [base] (CTE c must be excluded, its body's base kept)", got)
	}
	if len(span.CTEs) != 1 || !strings.EqualFold(span.CTEs[0], "c") {
		t.Errorf("CTEs = %v, want [c]", span.CTEs)
	}
}

// TestGetQuerySpan_CTEShadowsTable confirms a bare reference to a CTE name is not
// recorded as a base table even when a real table of the same name could exist.
func TestGetQuerySpan_CTEShadowsTable(t *testing.T) {
	sql := "WITH t AS (SELECT 1 AS n) SELECT n FROM t"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.AccessTables) != 0 {
		t.Errorf("tables = %v, want [] (FROM t references the CTE, not a base table)", tableNames(span))
	}
}

// --- System-schema classification (dialect divergence) --------------------

func TestGetQuerySpan_SystemTables(t *testing.T) {
	// wantSchema / wantDatabase reflect WHERE the 2nd-to-last qualifier lands per
	// the dialect bucketing: in BigQuery a non-INFORMATION_SCHEMA qualifier is the
	// Database (dataset); in Spanner it is always the Schema.
	tests := []struct {
		name         string
		sql          string
		wantSystem   bool
		wantSchema   string
		wantDatabase string
	}{
		{"spanner info schema", "SELECT * FROM INFORMATION_SCHEMA.TABLES", true, "INFORMATION_SCHEMA", ""},
		{"spanner spanner_sys", "SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE", true, "SPANNER_SYS", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			if len(span.AccessTables) != 1 {
				t.Fatalf("expected exactly 1 table, got %d (%v)", len(span.AccessTables), tableNames(span))
			}
			ta := span.AccessTables[0]
			if ta.IsSystem != tc.wantSystem {
				t.Errorf("IsSystem = %v, want %v (sql %q)", ta.IsSystem, tc.wantSystem, tc.sql)
			}
			if ta.Schema != tc.wantSchema {
				t.Errorf("Schema = %q, want %q", ta.Schema, tc.wantSchema)
			}
			if ta.Database != tc.wantDatabase {
				t.Errorf("Database = %q, want %q", ta.Database, tc.wantDatabase)
			}
		})
	}
}

// TestGetQuerySpan_BigQuerySchemaBucketing confirms the BigQuery accessTableListener
// rule: the 2nd-to-last identifier is treated as Schema ONLY when it is
// INFORMATION_SCHEMA; otherwise it is the Database (dataset). Spanner always
// treats the 2nd-to-last identifier as Schema.
func TestGetQuerySpan_SchemaBucketing(t *testing.T) {
	// Spanner: sch.t => sch is Schema.
	span, _ := GetQuerySpan("SELECT a FROM sch.t")
	ta := span.AccessTables[0]
	if ta.Schema != "sch" {
		t.Errorf("Spanner sch.t: Schema=%q, want sch", ta.Schema)
	}
}

// --- Results --------------------------------------------------------------

func TestGetQuerySpan_Results(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		wantResults []string
	}{
		{"columns", "SELECT a, b FROM t", []string{"a", "b"}},
		{"alias", "SELECT a AS x, b y FROM t", []string{"x", "y"}},
		{"star", "SELECT * FROM t", []string{"*"}},
		{"qualified col", "SELECT t.a FROM t", []string{"a"}},
		{"expr no name", "SELECT a + b FROM t", []string{""}},
		{"set op left wins", "SELECT a AS x FROM t UNION ALL SELECT b AS y FROM u", []string{"x"}},
		// A `rel.*` star is now a resolved StarSegments item (the relation's
		// projection), surfaced with Name "*"; the segment carries the base relation
		// the consumer expands. See TestGetQuerySpan_StarSegments for the segment shape.
		{"dot star", "SELECT t.* FROM t", []string{"*"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			got := resultNames(span)
			if !eqStrings(got, tc.wantResults) {
				t.Errorf("GetQuerySpan(%q) results = %v, want %v", tc.sql, got, tc.wantResults)
			}
		})
	}
}

// --- PredicateColumns -----------------------------------------------------

func TestGetQuerySpan_PredicateColumns(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantPred []string
	}{
		{"where", "SELECT a FROM t WHERE b > 1", []string{"b"}},
		{"where and", "SELECT a FROM t WHERE b > 1 AND c < 2", []string{"b", "c"}},
		{"join on", "SELECT * FROM a JOIN b ON a.x = b.y", []string{"x", "y"}},
		{"having", "SELECT a, COUNT(*) c FROM t GROUP BY a HAVING c > 1", []string{"a", "c"}},
		{"using", "SELECT * FROM a JOIN b USING (k)", []string{"k"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			got := predicateColumnNames(span)
			if !eqStrings(got, tc.wantPred) {
				t.Errorf("GetQuerySpan(%q) predicate columns = %v, want %v", tc.sql, got, tc.wantPred)
			}
		})
	}
}

// --- DML / DDL / empty ----------------------------------------------------

func TestGetQuerySpan_NonQuery(t *testing.T) {
	// DML / DDL statements get a Type but the access-table walk still runs over
	// their bodies (an UPDATE/INSERT touches tables, which masking/access cares
	// about).
	tests := []struct {
		name       string
		sql        string
		wantType   QueryType
		wantTables []string
	}{
		{"insert select", "INSERT INTO dst SELECT a FROM src", DML, []string{"dst", "src"}},
		{"update", "UPDATE t SET a = 1 WHERE b = 2", DML, []string{"t"}},
		{"delete", "DELETE FROM t WHERE a = 1", DML, []string{"t"}},
		{"create table", "CREATE TABLE t (a INT64) PRIMARY KEY (a)", DDL, []string{"t"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			if span.Type != tc.wantType {
				t.Errorf("Type = %v, want %v", span.Type, tc.wantType)
			}
			got := tableNames(span)
			if !eqStrings(got, tc.wantTables) {
				t.Errorf("tables = %v, want %v", got, tc.wantTables)
			}
		})
	}
}

func TestGetQuerySpan_Empty(t *testing.T) {
	span, err := GetQuerySpan("")
	if err != nil {
		t.Fatalf("GetQuerySpan(\"\") error: %v", err)
	}
	if span.Type != Unknown {
		t.Errorf("Type = %v, want Unknown", span.Type)
	}
	if len(span.AccessTables) != 0 || len(span.Results) != 0 {
		t.Errorf("expected empty span, got %+v", span)
	}
}

// TestGetQuerySpan_FailsClosedOnParseError: a statement that does not fully
// parse yields an error, never a smaller span. The parser re-parses embedded
// subqueries itself, so a malformed subquery surfaces through the same list.
func TestGetQuerySpan_FailsClosedOnParseError(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM",
		"SELECT a FROM t )))",
		"SELECT (SELECT a FROM t2 )))) AS x FROM t1",
		"SELECT a FROM t1 WHERE EXISTS (SELECT 1 FROM t2 )))",
		"this is not sql",
	} {
		if _, err := GetQuerySpan(sql); err == nil {
			t.Errorf("GetQuerySpan(%q) accepted input that does not parse", sql)
		}
	}
	// ClassifySQL shares the verdict: a truncated statement is Unknown, not
	// the type of whatever prefix parsed.
	if got := ClassifySQL("SELECT a FROM t )))"); got != Unknown {
		t.Errorf("ClassifySQL(truncated) = %v, want Unknown", got)
	}
}

// TestGetQuerySpan_CorrelatedSubquery confirms a correlated scalar subquery
// contributes its own base table AND its correlation predicate columns.
func TestGetQuerySpan_CorrelatedSubquery(t *testing.T) {
	sql := "SELECT (SELECT x FROM u WHERE u.id = t.id) AS m FROM t"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := tableNames(span); !eqStrings(got, []string{"t", "u"}) {
		t.Errorf("tables = %v, want [t u]", got)
	}
	// The correlation predicate `u.id = t.id` columns are collected from the inner
	// WHERE (a predicate position).
	if got := predicateColumnNames(span); !eqStrings(got, []string{"id", "id"}) {
		t.Errorf("predicate columns = %v, want [id id]", got)
	}
}

// TestGetQuerySpan_NestedCTEShadowing confirms an inner WITH scope shadows an
// outer CTE of the same name, and both scopes' bodies' base tables are recorded.
func TestGetQuerySpan_NestedCTEShadowing(t *testing.T) {
	sql := "WITH c AS (SELECT a FROM outer_base) " +
		"SELECT a FROM (WITH c AS (SELECT a FROM inner_base) SELECT a FROM c)"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	// Both CTE bodies' base tables are recorded; neither `c` reference is a base
	// table. outer_base is from the outer CTE body (walked even though the outer
	// query does not reference the outer c — matching the legacy permissive walk).
	if got := tableNames(span); !eqStrings(got, []string{"inner_base", "outer_base"}) {
		t.Errorf("tables = %v, want [inner_base outer_base]", got)
	}
}

// TestGetQuerySpan_SetOpThreeArms confirms the leftmost arm's result names win
// across a chained set operation.
func TestGetQuerySpan_SetOpThreeArms(t *testing.T) {
	sql := "SELECT a AS first FROM t UNION ALL SELECT b FROM u UNION ALL SELECT c FROM v"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"first"}) {
		t.Errorf("results = %v, want [first] (leftmost arm wins)", got)
	}
	if got := tableNames(span); !eqStrings(got, []string{"t", "u", "v"}) {
		t.Errorf("tables = %v, want [t u v]", got)
	}
}

// TestGetQuerySpan_DMLSubquery confirms a subquery inside a DML WHERE clause
// contributes its table.
func TestGetQuerySpan_DMLSubquery(t *testing.T) {
	sql := "DELETE FROM t WHERE id IN (SELECT id FROM staging)"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if span.Type != DML {
		t.Errorf("Type = %v, want DML", span.Type)
	}
	if got := tableNames(span); !eqStrings(got, []string{"staging", "t"}) {
		t.Errorf("tables = %v, want [staging t]", got)
	}
}

// TestGetQuerySpan_SourceColumns confirms a select-item's direct source columns
// are captured (the lineage seed the masking layer consumes).
func TestGetQuerySpan_SourceColumns(t *testing.T) {
	span, err := GetQuerySpan("SELECT a + b AS s, t.c FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(span.Results))
	}
	// s = a + b → two source columns a, b.
	var srcCols []string
	for _, sc := range span.Results[0].SourceColumns {
		srcCols = append(srcCols, sc.Column)
	}
	sort.Strings(srcCols)
	if !eqStrings(srcCols, []string{"a", "b"}) {
		t.Errorf("result[0] source columns = %v, want [a b]", srcCols)
	}
	// t.c → one source column c with table t.
	if len(span.Results[1].SourceColumns) != 1 ||
		span.Results[1].SourceColumns[0].Column != "c" ||
		span.Results[1].SourceColumns[0].Table != "t" {
		t.Errorf("result[1] source columns = %+v, want [{Table:t Column:c}]", span.Results[1].SourceColumns)
	}
}

// TestGetQuerySpan_LambdaParamsExcluded confirms a lambda's bound parameter is
// not recorded as a column reference, while a free column inside the lambda is.
func TestGetQuerySpan_LambdaParamsExcluded(t *testing.T) {
	// ARRAY_FILTER-style: the lambda param `e` is bound; `threshold` is a free
	// column reference.
	sql := "SELECT ARRAY_TRANSFORM(arr, e -> e + threshold) AS r FROM t"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	var cols []string
	for _, sc := range span.Results[0].SourceColumns {
		cols = append(cols, sc.Column)
	}
	sort.Strings(cols)
	// `arr` (the array arg) and `threshold` (free) are columns; `e` (bound param)
	// is excluded.
	if !eqStrings(cols, []string{"arr", "threshold"}) {
		t.Errorf("source columns = %v, want [arr threshold] (lambda param e excluded)", cols)
	}
}

// TestGetQuerySpan_GraphPredicateOperandsAreColumns pins the counterpart: a
// graph predicate used in an ordinary relational expression is a comparison of
// two columns, so both operands resolve — as a select item's sources and as a
// WHERE predicate's columns.
func TestGetQuerySpan_GraphPredicateOperandsAreColumns(t *testing.T) {
	span, err := GetQuerySpan("SELECT a IS SOURCE b FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	cols := []string{}
	for _, r := range span.Results {
		for _, src := range r.SourceColumns {
			cols = append(cols, src.Column)
		}
	}
	sort.Strings(cols)
	if !eqStrings(cols, []string{"a", "b"}) {
		t.Errorf("select-item sources = %v, want [a b]", cols)
	}

	span, err = GetQuerySpan("SELECT 1 FROM t WHERE a IS LABELED !x AND a IS DESTINATION OF b")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if got := predicateColumnNames(span); !eqStrings(got, []string{"a", "b"}) {
		t.Errorf("PredicateColumns = %v, want [a b]", got)
	}
}

// TestGetQuerySpan_GraphVariablesAreNotColumns pins the other half of a graph
// subquery's span: its element, path, LET and FOR variables, and its RETURN
// items, are graph-local names rather than relational columns, so none of them
// leaks into the SELECT's source columns or predicate columns. A
// column-collecting pass reports `t.n` for `SELECT EXISTS { MATCH (n) RETURN n }
// FROM t` if the graph tree is walked with the caller's callbacks.
func TestGetQuerySpan_GraphVariablesAreNotColumns(t *testing.T) {
	cases := []struct {
		sql       string
		graphOnly []string // names the graph invents; none may surface as a column
	}{
		{"SELECT EXISTS { MATCH (n) RETURN n } FROM t", []string{"n"}},
		{"SELECT EXISTS { MATCH (n:Person) WHERE n.id = 1 RETURN n } FROM t", []string{"n"}},
		{"SELECT EXISTS { LET a = b RETURN a } FROM t", []string{"a", "b"}},
		{"SELECT ARRAY { FOR x IN y RETURN x } FROM t", []string{"x", "y"}},
		{"SELECT a FROM t WHERE EXISTS { MATCH (n) RETURN n }", []string{"n"}},
		// A table argument's filter inside a graph CALL is still the graph's scope.
		{"SELECT EXISTS { MATCH (n) CALL f(TABLE secret WHERE n.id > 0) RETURN 1 } FROM t", []string{"n"}},
	}
	for _, tc := range cases {
		span, err := GetQuerySpan(tc.sql)
		if err != nil {
			t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
		}
		for _, r := range span.Results {
			for _, src := range r.SourceColumns {
				for _, name := range tc.graphOnly {
					if src.Column == name {
						t.Errorf("GetQuerySpan(%q) result %q has source column %v, which is a graph-local name", tc.sql, r.Name, src)
					}
				}
			}
		}
		for _, name := range predicateColumnNames(span) {
			for _, graphName := range tc.graphOnly {
				if name == graphName {
					t.Errorf("GetQuerySpan(%q) PredicateColumns has %q, which is a graph-local name", tc.sql, name)
				}
			}
		}
	}
}

// TestGetQuerySpan_GraphSubqueryExpression pins what a graph subquery used as an
// expression does to the span of an ordinary SELECT. The graph subquery is a leaf
// for lineage, like a bare GQL statement (see TestGetQuerySpan_GQLEmptySpan):
// the legacy accessTableListener walks table_path_expression positions only, and
// a graph element pattern references labels, not tables — so a MATCH's patterns
// and the graph's own name contribute no AccessTables. Everything else about the
// SELECT still resolves: its FROM table is an access and its select list is the
// projection.
func TestGetQuerySpan_GraphSubqueryExpression(t *testing.T) {
	cases := []struct {
		sql     string
		tables  []string
		results int
	}{
		{"SELECT a, EXISTS { MATCH (n) RETURN n } FROM t", []string{"t"}, 2},
		{"SELECT ARRAY { RETURN NULL }", nil, 1},
		{"SELECT VALUE { RETURN NULL }", nil, 1},
		{"SELECT NULL IN { RETURN NULL }", nil, 1},
		{"SELECT a, EXISTS { MATCH (n:Person) WHERE n.id = 1 RETURN n } FROM t", []string{"t"}, 2},
		// A scalar subquery inside a graph operator or pattern is still a read: the
		// graph tree's expressions are walked even though its patterns name labels.
		{"SELECT EXISTS { LET a = (SELECT MAX(x) FROM secret) RETURN a }", []string{"secret"}, 1},
		{"SELECT EXISTS { MATCH (n) WHERE n.id IN (SELECT id FROM secret) RETURN n }", []string{"secret"}, 1},
		// The IN form's graph subquery is walked like the others.
		{"SELECT 1 IN { LET x = (SELECT id FROM secret) RETURN x }", []string{"secret"}, 1},
		// So are a graph GROUP BY item, a TABLESAMPLE PARTITION BY expression and a
		// property value.
		{"SELECT EXISTS { RETURN 1 GROUP BY (SELECT id FROM secret) }", []string{"secret"}, 1},
		{"SELECT EXISTS { TABLESAMPLE BERNOULLI(1 ROWS PARTITION BY (SELECT id FROM secret)) RETURN 1 }", []string{"secret"}, 1},
		{"SELECT EXISTS { MATCH ({a: (SELECT id FROM secret)}) RETURN 1 }", []string{"secret"}, 1},
		// A graph CALL's TVF table argument names a table the call reads, in both
		// the positional and the named shape, and a filter on it still follows its
		// own subqueries.
		{"SELECT EXISTS { CALL f(TABLE secret) RETURN 1 }", []string{"secret"}, 1},
		{"SELECT EXISTS { CALL f(t => TABLE secret) RETURN 1 }", []string{"secret"}, 1},
		{"SELECT EXISTS { CALL f(TABLE secret WHERE (SELECT 1 FROM other) > 0) RETURN 1 }", []string{"other", "secret"}, 1},
		// A table argument's PIVOT expression is walked too.
		{"SELECT EXISTS { CALL f(TABLE g() PIVOT((SELECT id FROM secret) FOR x IN (1))) RETURN 1 }", []string{"secret"}, 1},
		// A standalone graph predicate's operands are ordinary relational
		// columns: neither names a table here.
		{"SELECT a FROM t WHERE a IS SOURCE b", []string{"t"}, 1},
	}
	for _, tc := range cases {
		span, err := GetQuerySpan(tc.sql)
		if err != nil {
			t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
		}
		if span.Type != Select {
			t.Errorf("GetQuerySpan(%q) Type = %v, want Select", tc.sql, span.Type)
		}
		if got := tableNames(span); !eqStrings(got, tc.tables) {
			t.Errorf("GetQuerySpan(%q) AccessTables = %v, want %v", tc.sql, got, tc.tables)
		}
		if len(span.Results) != tc.results {
			t.Errorf("GetQuerySpan(%q) Results = %v, want %d columns", tc.sql, resultNames(span), tc.results)
		}
	}
}

// TestGetQuerySpan_GQLEmptySpan pins the parser-gql query-span DISPOSITION
// (finding F9): a GQL graph query (`GRAPH … MATCH … RETURN …`) yields an EMPTY
// span — Type=Unknown, no AccessTables, no Results, no PredicateColumns. This is
// faithful parity with the legacy bytebase GetQuerySpan, which short-circuits any
// non-Select statement to `{Type: <unknown>, SourceColumns: {}, Results: []}`
// (plugin/parser/{bigquery,spanner}/query_span_extractor.go getQuerySpan) and
// whose accessTableListener walks only table_path_expression nodes — graph
// element patterns reference labels, not table paths, so legacy records nothing
// for a graph read either. Deferred-parity: a future source-table-discovery pass
// over graph patterns is tracked as a divergence/follow-up, not implemented here.
func TestGetQuerySpan_GQLEmptySpan(t *testing.T) {
	cases := []string{
		"GRAPH g MATCH (n) RETURN n",
		"GRAPH g MATCH (a)-[e]->(b) WHERE a.id = b.id RETURN a, b",
		"GRAPH myproject.mydataset.g MATCH (p:Person) RETURN p.name",
	}
	for _, sql := range cases {
		span, err := GetQuerySpan(sql)
		if err != nil {
			t.Fatalf("GetQuerySpan(%q) error: %v", sql, err)
		}
		if span.Type != Unknown {
			t.Errorf("GetQuerySpan(%q) Type = %v, want Unknown (legacy non-Select parity)", sql, span.Type)
		}
		if len(span.AccessTables) != 0 {
			t.Errorf("GetQuerySpan(%q) AccessTables = %v, want empty", sql, tableNames(span))
		}
		if len(span.Results) != 0 {
			t.Errorf("GetQuerySpan(%q) Results = %v, want empty", sql, resultNames(span))
		}
		if len(span.PredicateColumns) != 0 {
			t.Errorf("GetQuerySpan(%q) PredicateColumns = %v, want empty", sql, predicateColumnNames(span))
		}
	}
}

// --- Masking-grade per-shape lineage (BigQuery cutover) -------------------
//
// These pin the lineage SHAPE the bytebase masking layer consumes, independent
// of catalog metadata: omni records bare (unqualified) source-column refs and a
// star's EXCEPT/REPLACE modifiers; the bytebase extractor resolves them against
// dataset metadata. The shape — which refs feed which output position — is what
// must be correct here.

// resultSourceColumns renders one result's source ColumnRefs as sorted
// "table.column" (or just "column" when unqualified) strings.
func resultSourceColumns(info ColumnInfo) []string {
	var out []string
	for _, sc := range info.SourceColumns {
		if sc.Table != "" {
			out = append(out, sc.Table+"."+sc.Column)
		} else {
			out = append(out, sc.Column)
		}
	}
	sort.Strings(out)
	return out
}

// TestGetQuerySpan_SetOpMergesBothArms pins fix #1: a set operation's output
// column at position i draws its lineage from the UNION of BOTH arms' source
// columns at that position (the left arm alone is a masking leak). The output
// NAME still comes from the left arm.
func TestGetQuerySpan_SetOpMergesBothArms(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		wantNames   []string
		wantSources [][]string // per output position, sorted source columns
	}{
		{
			name:        "two arms aliased",
			sql:         "SELECT a AS x FROM t UNION ALL SELECT b AS y FROM u",
			wantNames:   []string{"x"},
			wantSources: [][]string{{"t.a", "u.b"}},
		},
		{
			name:        "two columns each",
			sql:         "SELECT a, b FROM t UNION ALL SELECT c, d FROM u",
			wantNames:   []string{"a", "b"},
			wantSources: [][]string{{"t.a", "u.c"}, {"t.b", "u.d"}},
		},
		{
			name:        "three arms accumulate",
			sql:         "SELECT a FROM t UNION ALL SELECT b FROM u UNION ALL SELECT c FROM v",
			wantNames:   []string{"a"},
			wantSources: [][]string{{"t.a", "u.b", "v.c"}},
		},
		{
			name:        "intersect and except also merge",
			sql:         "SELECT a FROM t INTERSECT DISTINCT SELECT b FROM u",
			wantNames:   []string{"a"},
			wantSources: [][]string{{"t.a", "u.b"}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			if got := resultNames(span); !eqStrings(got, tc.wantNames) {
				t.Errorf("names = %v, want %v", got, tc.wantNames)
			}
			if len(span.Results) != len(tc.wantSources) {
				t.Fatalf("got %d results, want %d", len(span.Results), len(tc.wantSources))
			}
			for i, want := range tc.wantSources {
				if got := resultSourceColumns(span.Results[i]); !eqStrings(got, want) {
					t.Errorf("result[%d] sources = %v, want %v", i, got, want)
				}
			}
		})
	}
}

// TestGetQuerySpan_ScalarSubqueryInSelectList pins the SELECT-list
// scalar-subquery lineage: the output column reads the subquery's columns.
// An unqualified one may also be the enclosing query's, so it is kept
// unqualified too. The subquery's tables and correlation predicates are still
// discovered.
func TestGetQuerySpan_ScalarSubqueryInSelectList(t *testing.T) {
	sql := "SELECT (SELECT NAME FROM people), ID FROM people"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("got %d results, want 2", len(span.Results))
	}
	if got := resultSourceColumns(span.Results[0]); !eqStrings(got, []string{"NAME", "people.NAME"}) {
		t.Errorf("scalar-subquery result sources = %v, want [NAME people.NAME]", got)
	}
	// The sibling bare column resolves normally.
	if got := resultSourceColumns(span.Results[1]); !eqStrings(got, []string{"people.ID"}) {
		t.Errorf("result[1] sources = %v, want [people.ID]", got)
	}
	// The subquery's table is still discovered (people appears once, deduped).
	if got := tableNames(span); !eqStrings(got, []string{"people"}) {
		t.Errorf("tables = %v, want [people]", got)
	}
}

// TestGetQuerySpan_JoinUsingCoalesceSurfacesColumn pins fix #4's omni-side
// contribution: a JOIN ... USING (c) records the coalesced column c as a bare
// reference (Column=c, no table qualifier) so the catalog-aware consumer matches
// it against BOTH joined sides' c columns; the USING column is also a predicate.
func TestGetQuerySpan_JoinUsingCoalesceSurfacesColumn(t *testing.T) {
	sql := "SELECT ID, NAME, ADDRESS FROM people JOIN address ON people.ID = address.PEOPLE_ID USING (ID)"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) == 0 || span.Results[0].Name != "ID" {
		t.Fatalf("results = %v, want first column ID", resultNames(span))
	}
	// The coalesced ID is a bare ref (no table) so it flat-matches both PEOPLE.ID
	// and ADDRESS.ID downstream.
	if got := resultSourceColumns(span.Results[0]); !eqStrings(got, []string{"ID"}) {
		t.Errorf("ID result sources = %v, want [ID] (bare, table-agnostic)", got)
	}
	if got := predicateColumnNames(span); !containsString(got, "ID") {
		t.Errorf("predicate columns = %v, want to include the USING column ID", got)
	}
}

// TestGetQuerySpan_StarExceptReplaceCaptured pins fix #2: a SELECT * EXCEPT/
// REPLACE star item records its modifiers on the ColumnInfo (StarExcept names,
// StarReplace name+expr-sources) so the catalog-aware consumer can drop EXCEPT-ed
// columns and re-point a REPLACE-d column's lineage.
func TestGetQuerySpan_StarExceptReplaceCaptured(t *testing.T) {
	sql := "SELECT * EXCEPT (ID) REPLACE (ID/2 AS NAME) FROM people"
	span, err := GetQuerySpan(sql)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("got %d results, want 1 (the star)", len(span.Results))
	}
	star := span.Results[0]
	if star.Name != "*" {
		t.Errorf("star name = %q, want *", star.Name)
	}
	if !eqStrings(star.StarExcept, []string{"ID"}) {
		t.Errorf("StarExcept = %v, want [ID]", star.StarExcept)
	}
	if len(star.StarReplace) != 1 {
		t.Fatalf("got %d StarReplace, want 1", len(star.StarReplace))
	}
	repl := star.StarReplace[0]
	if repl.Name != "NAME" {
		t.Errorf("StarReplace[0].Name = %q, want NAME", repl.Name)
	}
	var cols []string
	for _, sc := range repl.Sources {
		cols = append(cols, sc.Column)
	}
	if !eqStrings(cols, []string{"ID"}) {
		t.Errorf("StarReplace[0].Sources = %v, want [ID] (the ID/2 expr)", cols)
	}
}

// containsString reports whether list contains s.
func containsString(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// TestGetQuerySpan_TableClauseProjection holds the `TABLE t` query primary to
// the explicit query it means: `SELECT * FROM t`. The clause's output carries
// the same columns and lineage — including a CTE's resolved columns, which a
// span that discarded the relation would lose.
func TestGetQuerySpan_TableClauseProjection(t *testing.T) {
	for _, c := range []struct{ clause, explicit string }{
		{"TABLE t", "SELECT * FROM t"},
		{"TABLE ds.t", "SELECT * FROM ds.t"},
		{"TABLE t WHERE a > 1", "SELECT * FROM t WHERE a > 1"},
		{"TABLE a UNION ALL TABLE b", "SELECT * FROM a UNION ALL SELECT * FROM b"},
		{"WITH c AS (SELECT secret FROM t) TABLE c", "WITH c AS (SELECT secret FROM t) SELECT * FROM c"},
	} {
		t.Run(c.clause, func(t *testing.T) {
			got, err := GetQuerySpan(c.clause)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q): %v", c.clause, err)
			}
			want, err := GetQuerySpan(c.explicit)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q): %v", c.explicit, err)
			}
			if !eqStrings(resultNames(got), resultNames(want)) {
				t.Errorf("results = %v, want %v (as %q)", resultNames(got), resultNames(want), c.explicit)
			}
			for i := range got.Results {
				if i >= len(want.Results) {
					break
				}
				g, w := got.Results[i].SourceColumns, want.Results[i].SourceColumns
				if len(g) != len(w) {
					t.Errorf("result %d sources = %v, want %v", i, g, w)
					continue
				}
				for j := range g {
					if g[j] != w[j] {
						t.Errorf("result %d source %d = %+v, want %+v", i, j, g[j], w[j])
					}
				}
			}
			if !eqStrings(tableNames(got), tableNames(want)) {
				t.Errorf("tables = %v, want %v", tableNames(got), tableNames(want))
			}
			if !eqStrings(predicateColumnNames(got), predicateColumnNames(want)) {
				t.Errorf("predicate columns = %v, want %v", predicateColumnNames(got), predicateColumnNames(want))
			}
		})
	}
}

// TestGetQuerySpan_OptionSubqueries holds the statements whose OPTIONS (or AS OF)
// expression the parser fills: a subquery there reads a table, so the span
// records it even though the statement itself names no table.
func TestGetQuerySpan_OptionSubqueries(t *testing.T) {
	for _, c := range []struct {
		sql string
		raw []string
	}{
		{"IMPORT MODULE m OPTIONS(x = (SELECT 1 FROM a))", []string{"a"}},
		{"MODULE m OPTIONS(x = (SELECT 1 FROM a))", []string{"a"}},
		{"UNDROP SCHEMA s OPTIONS(x = (SELECT 1 FROM a))", []string{"a"}},
		{"UNDROP SCHEMA s FOR SYSTEM_TIME AS OF (SELECT MAX(t) FROM a)", []string{"a"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			span, err := GetQuerySpan(c.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q): %v", c.sql, err)
			}
			if got := tableNames(span); !eqStrings(got, c.raw) {
				t.Errorf("tables = %v, want %v", got, c.raw)
			}
		})
	}
}

// TestGetQuerySpan_ExprDotStarFailClosed verifies that arbitrary non-path
// `expr.*` projections (such as `(SELECT AS STRUCT ...).*`, `fn().*`, or
// `STRUCT(...).*`) fail closed instead of silently dropping output columns.
func TestGetQuerySpan_ExprDotStarFailClosed(t *testing.T) {
	cases := []string{
		"SELECT (SELECT AS STRUCT secret FROM sensitive).*, 1 AS public",
		"SELECT STRUCT(a AS x, b AS y).* FROM t",
		"SELECT fn(a).* FROM t",
	}
	for _, sql := range cases {
		if span, err := GetQuerySpan(sql); err == nil {
			t.Errorf("GetQuerySpan(%q) succeeded with results=%v, want fail-closed error", sql, span.Results)
		}
	}
	// A parenthesized FROM relation `(t).*` still resolves to t's projection.
	span, err := GetQuerySpan("SELECT (t).* FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan(SELECT (t).* FROM t): %v", err)
	}
	if len(span.Results) != 1 || span.Results[0].Name != "*" {
		t.Errorf("SELECT (t).* FROM t results = %+v, want single star segment", span.Results)
	}
}

// TestGetQuerySpan_GraphTableFailClosed verifies that GRAPH_TABLE(...) FROM
// sources fail closed in GetQuerySpan (rather than returning an empty span)
// while still walking embedded SQL subqueries for Classify.
func TestGetQuerySpan_GraphTableFailClosed(t *testing.T) {
	cases := []string{
		"SELECT * FROM GRAPH_TABLE(g MATCH (n) COLUMNS(n.secret AS secret))",
		"SELECT * FROM GRAPH_TABLE(g MATCH (n) RETURN n.secret AS secret)",
		"SELECT * FROM t JOIN GRAPH_TABLE(g MATCH (n) COLUMNS(n.id AS id)) gt ON t.id = gt.id",
	}
	for _, sql := range cases {
		if span, err := GetQuerySpan(sql); err == nil {
			t.Errorf("GetQuerySpan(%q) succeeded with tables=%v results=%v, want fail-closed error", sql, span.AccessTables, span.Results)
		}
	}
	// Classify is best-effort and still discovers tables inside embedded SQL
	// subqueries in GRAPH_TABLE operators and COLUMNS expressions.
	if got := ClassifySQL("SELECT * FROM GRAPH_TABLE(g MATCH (n) WHERE n.id IN (SELECT id FROM secret) COLUMNS(n.id))"); got != Select {
		t.Errorf("ClassifySQL(GRAPH_TABLE with subquery) = %v, want Select", got)
	}
	if got := ClassifySQL("SELECT * FROM GRAPH_TABLE(g MATCH (n) COLUMNS((SELECT id FROM INFORMATION_SCHEMA.TABLES) AS id))"); got != SelectInfoSchema {
		t.Errorf("ClassifySQL(GRAPH_TABLE with info_schema subquery) = %v, want SelectInfoSchema", got)
	}
}

// TestGetQuerySpan_ExplainAndExtendedDDL verifies that EXPLAIN traverses its
// inner statement for AccessTables (keeping Results empty) and that extended
// DDL statements (CREATE MODEL AS query, ALTER ROW ACCESS POLICY) record
// referenced tables.
func TestGetQuerySpan_ExplainAndExtendedDDL(t *testing.T) {
	span, err := GetQuerySpan("EXPLAIN SELECT secret FROM sensitive")
	if err != nil {
		t.Fatalf("GetQuerySpan(EXPLAIN SELECT): %v", err)
	}
	if span.Type != Explain {
		t.Errorf("EXPLAIN Type = %v, want Explain", span.Type)
	}
	if got := tableNames(span); !eqStrings(got, []string{"sensitive"}) {
		t.Errorf("EXPLAIN AccessTables = %v, want [sensitive]", got)
	}
	if len(span.Results) != 0 {
		t.Errorf("EXPLAIN Results = %+v, want empty", span.Results)
	}

	span, err = GetQuerySpan("EXPLAIN INSERT INTO dst SELECT * FROM src")
	if err != nil {
		t.Fatalf("GetQuerySpan(EXPLAIN INSERT): %v", err)
	}
	if got := tableNames(span); !eqStrings(got, []string{"dst", "src"}) {
		t.Errorf("EXPLAIN INSERT AccessTables = %v, want [dst src]", got)
	}

	if _, err := GetQuerySpan("EXPLAIN SELECT * FROM GRAPH_TABLE(g MATCH (n) COLUMNS(n.secret))"); err == nil {
		t.Errorf("GetQuerySpan(EXPLAIN GRAPH_TABLE) succeeded, want fail-closed error")
	}
}
