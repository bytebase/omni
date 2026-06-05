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
		dialect    Dialect
		wantTables []string
	}{
		{"single table", "SELECT a FROM t", DialectBigQuery, []string{"t"}},
		{"qualified table bq", "SELECT a FROM ds.t", DialectBigQuery, []string{"ds.t"}},
		{"three-part bq", "SELECT a FROM proj.ds.t", DialectBigQuery, []string{"proj.ds.t"}},
		{"join", "SELECT * FROM a JOIN b ON a.x = b.x", DialectBigQuery, []string{"a", "b"}},
		{"comma join", "SELECT * FROM a, b", DialectBigQuery, []string{"a", "b"}},
		{"chained join", "SELECT * FROM a JOIN b ON a.x=b.x JOIN c ON b.y=c.y", DialectBigQuery, []string{"a", "b", "c"}},
		{"aliased table excluded-from-name", "SELECT * FROM t AS x", DialectBigQuery, []string{"t"}},
		{"subquery in from", "SELECT * FROM (SELECT a FROM inner_t) AS s", DialectBigQuery, []string{"inner_t"}},
		{"scalar subquery", "SELECT (SELECT MAX(x) FROM u) AS m FROM t", DialectBigQuery, []string{"t", "u"}},
		{"exists subquery", "SELECT * FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id = t.id)", DialectBigQuery, []string{"t", "u"}},
		{"in subquery", "SELECT * FROM t WHERE id IN (SELECT id FROM u)", DialectBigQuery, []string{"t", "u"}},
		{"set op", "SELECT a FROM t UNION ALL SELECT a FROM u", DialectBigQuery, []string{"t", "u"}},
		{"dedup same table", "SELECT * FROM t WHERE a IN (SELECT a FROM t)", DialectBigQuery, []string{"t"}},
		// Spanner schema.table bucketing (named schema under one DB).
		{"qualified table spanner", "SELECT a FROM myschema.t", DialectSpanner, []string{"myschema.t"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql, tc.dialect)
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

// TestGetQuerySpan_CTEsExcluded confirms CTE names are recorded in span.CTEs and
// excluded from AccessTables, while the CTE BODY's base tables are kept.
func TestGetQuerySpan_CTEsExcluded(t *testing.T) {
	sql := "WITH c AS (SELECT a FROM base) SELECT a FROM c"
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
		dialect      Dialect
		wantSystem   bool
		wantSchema   string
		wantDatabase string
	}{
		{"bq info schema", "SELECT * FROM INFORMATION_SCHEMA.TABLES", DialectBigQuery, true, "INFORMATION_SCHEMA", ""},
		{"bq dataset.info_schema.view", "SELECT * FROM ds.INFORMATION_SCHEMA.COLUMNS", DialectBigQuery, true, "INFORMATION_SCHEMA", "ds"},
		{"bq spanner_sys not system", "SELECT * FROM SPANNER_SYS.X", DialectBigQuery, false, "", "SPANNER_SYS"},
		{"spanner info schema", "SELECT * FROM INFORMATION_SCHEMA.TABLES", DialectSpanner, true, "INFORMATION_SCHEMA", ""},
		{"spanner spanner_sys", "SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE", DialectSpanner, true, "SPANNER_SYS", ""},
		{"user table bq", "SELECT * FROM ds.users", DialectBigQuery, false, "", "ds"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql, tc.dialect)
			if err != nil {
				t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
			}
			if len(span.AccessTables) != 1 {
				t.Fatalf("expected exactly 1 table, got %d (%v)", len(span.AccessTables), tableNames(span))
			}
			ta := span.AccessTables[0]
			if ta.IsSystem != tc.wantSystem {
				t.Errorf("IsSystem = %v, want %v (sql %q dialect %v)", ta.IsSystem, tc.wantSystem, tc.sql, tc.dialect)
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
	// BigQuery: ds.t => ds is Database, Schema empty.
	span, _ := GetQuerySpan("SELECT a FROM ds.t", DialectBigQuery)
	ta := span.AccessTables[0]
	if ta.Database != "ds" || ta.Schema != "" {
		t.Errorf("BigQuery ds.t: Database=%q Schema=%q, want Database=ds Schema=\"\"", ta.Database, ta.Schema)
	}

	// Spanner: sch.t => sch is Schema.
	span, _ = GetQuerySpan("SELECT a FROM sch.t", DialectSpanner)
	ta = span.AccessTables[0]
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
		{"dot star", "SELECT t.* FROM t", []string{"t.*"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql, DialectBigQuery)
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
		{"qualify", "SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY b) = 1", []string{"b"}},
		{"using", "SELECT * FROM a JOIN b USING (k)", []string{"k"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql, DialectBigQuery)
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
		{"ctas", "CREATE TABLE t AS SELECT a FROM src", DDL, []string{"src", "t"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := GetQuerySpan(tc.sql, DialectBigQuery)
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
	span, err := GetQuerySpan("", DialectBigQuery)
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

// TestGetQuerySpan_PartialParse confirms tolerance of a partial/invalid parse:
// a malformed statement yields whatever was discoverable without panicking.
func TestGetQuerySpan_PartialParse(t *testing.T) {
	// Missing FROM target — invalid, but must not panic.
	span, err := GetQuerySpan("SELECT * FROM", DialectBigQuery)
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	_ = span // no panic / nil-deref is the assertion
}

// TestGetQuerySpan_CorrelatedSubquery confirms a correlated scalar subquery
// contributes its own base table AND its correlation predicate columns.
func TestGetQuerySpan_CorrelatedSubquery(t *testing.T) {
	sql := "SELECT (SELECT x FROM u WHERE u.id = t.id) AS m FROM t"
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
	span, err := GetQuerySpan("SELECT a + b AS s, t.c FROM t", DialectBigQuery)
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
	span, err := GetQuerySpan(sql, DialectBigQuery)
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
