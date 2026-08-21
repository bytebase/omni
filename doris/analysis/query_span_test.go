package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/doris/parser"
)

// tableSig is a compact form of TableAccess used for order-insensitive
// comparison in tests.
type tableSig struct {
	Database string
	Table    string
	Alias    string
}

func toSigs(tables []TableAccess) []tableSig {
	out := make([]tableSig, len(tables))
	for i, t := range tables {
		out[i] = tableSig{Database: t.Database, Table: t.Table, Alias: t.Alias}
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
			wantType:    QueryTypeSelect,
			wantTables:  []tableSig{{Table: "t"}},
			wantResults: []string{"a"},
		},
		{
			name:        "qualified table",
			sql:         "SELECT * FROM db.t",
			wantType:    QueryTypeSelect,
			wantTables:  []tableSig{{Database: "db", Table: "t"}},
			wantResults: []string{"*"},
		},
		{
			name:        "three-part table name",
			sql:         "SELECT * FROM catalog1.db1.t1",
			wantType:    QueryTypeSelect,
			wantTables:  []tableSig{{Database: "db1", Table: "t1"}},
			wantResults: []string{"*"},
		},
		{
			name:        "alias",
			sql:         "SELECT a FROM t AS x",
			wantType:    QueryTypeSelect,
			wantTables:  []tableSig{{Table: "t", Alias: "x"}},
			wantResults: []string{"a"},
		},
		{
			name:        "select item alias",
			sql:         "SELECT a AS col_a, b FROM t",
			wantType:    QueryTypeSelect,
			wantTables:  []tableSig{{Table: "t"}},
			wantResults: []string{"col_a", "b"},
		},
		{
			name:     "join of two tables",
			sql:      "SELECT a, b FROM t1 JOIN t2 ON t1.id = t2.id",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"a", "b"},
		},
		{
			name:     "inner and left join",
			sql:      "SELECT a FROM t1 INNER JOIN t2 ON t1.id = t2.id LEFT JOIN db.t3 ON t1.x = t3.x",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
				{Database: "db", Table: "t3"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "cte exclusion",
			sql:      "WITH c AS (SELECT * FROM real_t) SELECT * FROM c",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "real_t"},
			},
			wantResults: []string{"*"},
			wantCTEs:    []string{"c"},
		},
		{
			name: "multiple ctes chained",
			sql: `WITH a AS (SELECT id FROM t1),
			      b AS (SELECT id FROM a JOIN t2 ON a.id = t2.id)
			      SELECT id FROM b`,
			wantType: QueryTypeSelect,
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
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "other"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "scalar subquery in select",
			sql:      "SELECT a, (SELECT MAX(x) FROM agg) FROM t",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "agg"},
			},
		},
		{
			name:     "exists subquery",
			sql:      "SELECT a FROM t WHERE EXISTS (SELECT 1 FROM other WHERE other.id = t.id)",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t"},
				{Table: "other"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "from subquery",
			sql:      "SELECT x.a FROM (SELECT a FROM inner_t) AS x",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "inner_t"},
			},
		},
		{
			name:     "union two tables",
			sql:      "SELECT a FROM t1 UNION SELECT a FROM t2",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
			},
			wantResults: []string{"a"},
		},
		{
			name:     "union all three tables",
			sql:      "SELECT a FROM t1 UNION ALL SELECT a FROM t2 UNION ALL SELECT a FROM db.t3",
			wantType: QueryTypeSelect,
			wantTables: []tableSig{
				{Table: "t1"},
				{Table: "t2"},
				{Database: "db", Table: "t3"},
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

func TestGetQuerySpan_NonSelect(t *testing.T) {
	tests := []struct {
		sql      string
		wantType QueryType
	}{
		{"CREATE TABLE t (id INT)", QueryTypeDDL},
		{"SHOW TABLES", QueryTypeSelectInfoSchema},
		{"INSERT INTO t VALUES (1)", QueryTypeDML},
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
			// Non-SELECT statements should not populate AccessTables via our
			// SELECT walker. (They may still land here in future when DML
			// parsing is wired through — for now we only assert Type.)
		})
	}
}

func TestGetQuerySpan_Empty(t *testing.T) {
	span, err := GetQuerySpan("")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if span.Type != QueryTypeUnknown {
		t.Errorf("Type = %v, want QueryTypeUnknown", span.Type)
	}
	if len(span.AccessTables) != 0 {
		t.Errorf("AccessTables = %v, want empty", span.AccessTables)
	}
	if len(span.Results) != 0 {
		t.Errorf("Results = %v, want empty", span.Results)
	}
}

func TestGetQuerySpan_DedupTables(t *testing.T) {
	// Same unaliased table referenced twice — should only appear once.
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

func TestGetQuerySpan_SelectItemSourceColumns(t *testing.T) {
	// Parser quirk: a SELECT item that starts with a qualified identifier
	// (t.a) takes a fast path and only produces a ColumnRef, so we use bare
	// identifiers here to exercise the full expression parser.
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
	want := []ColumnRef{
		{Column: "a"},
		{Column: "b"},
	}
	for i, w := range want {
		if r.SourceColumns[i] != w {
			t.Errorf("SourceColumns[%d] = %+v, want %+v", i, r.SourceColumns[i], w)
		}
	}
}

func TestGetQuerySpan_CTEShadowsTable(t *testing.T) {
	// Even if a physical table exists elsewhere named "real_t", a CTE with
	// the same name inside the WITH clause means the outer reference resolves
	// to the CTE — AccessTables should only contain the CTE's body tables.
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

func TestGetQuerySpan_ExtractSourceColumns(t *testing.T) {
	// EXTRACT(unit FROM col) must expose col as a source column so masking and
	// lineage see through the expression.
	span, err := GetQuerySpan("SELECT EXTRACT(YEAR FROM created_at) AS y FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results len = %d, want 1 (%+v)", len(span.Results), span.Results)
	}
	r := span.Results[0]
	if r.Name != "y" {
		t.Errorf("Results[0].Name = %q, want y", r.Name)
	}
	want := []ColumnRef{{Column: "created_at"}}
	if len(r.SourceColumns) != len(want) {
		t.Fatalf("SourceColumns = %+v, want %+v", r.SourceColumns, want)
	}
	for i, w := range want {
		if r.SourceColumns[i] != w {
			t.Errorf("SourceColumns[%d] = %+v, want %+v", i, r.SourceColumns[i], w)
		}
	}
}

func TestGetQuerySpan_LambdaSourceColumns(t *testing.T) {
	// Columns referenced inside a lambda body must reach lineage, but the
	// lambda's own parameters must not: they are local bindings, not columns,
	// and reporting one would feed a nonexistent column to masking.
	tests := []struct {
		name    string
		sql     string
		want    []string
		notWant []string
	}{
		{
			name:    "single parameter",
			sql:     "SELECT array_map(x -> x + bonus, scores) AS adj FROM t",
			want:    []string{"bonus", "scores"},
			notWant: []string{"x"},
		},
		{
			name:    "multiple parameters",
			sql:     "SELECT array_map((x, y) -> x + y + bonus, a, b) AS adj FROM t",
			want:    []string{"bonus", "a", "b"},
			notWant: []string{"x", "y"},
		},
		{
			name: "qualified name is a column even when it shadows a parameter",
			sql:  "SELECT array_map(x -> t.x + bonus, scores) AS adj FROM t",
			want: []string{"x", "bonus", "scores"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span, err := GetQuerySpan(tt.sql)
			if err != nil {
				t.Fatalf("GetQuerySpan returned error: %v", err)
			}
			if len(span.Results) != 1 {
				t.Fatalf("Results len = %d, want 1 (%+v)", len(span.Results), span.Results)
			}
			got := map[string]bool{}
			for _, sc := range span.Results[0].SourceColumns {
				got[sc.Column] = true
			}
			for _, want := range tt.want {
				if !got[want] {
					t.Errorf("SourceColumns %+v missing %q", span.Results[0].SourceColumns, want)
				}
			}
			for _, notWant := range tt.notWant {
				if got[notWant] {
					t.Errorf("SourceColumns %+v contains lambda parameter %q", span.Results[0].SourceColumns, notWant)
				}
			}
		})
	}
}

func TestGetQuerySpan_SubscriptKeepsTableAccess(t *testing.T) {
	// Before BYT-10084 the parser silently dropped everything after `[`, so
	// this query lost its FROM clause and the span reported no table access —
	// a fail-open blind spot for masking.
	span, err := GetQuerySpan("SELECT secret_col[1] FROM sensitive_table")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.AccessTables) != 1 || span.AccessTables[0].Table != "sensitive_table" {
		t.Fatalf("AccessTables = %+v, want [sensitive_table]", span.AccessTables)
	}
	if len(span.Results) != 1 {
		t.Fatalf("Results = %+v, want one column", span.Results)
	}
	got := map[string]bool{}
	for _, sc := range span.Results[0].SourceColumns {
		got[sc.Column] = true
	}
	if !got["secret_col"] {
		t.Errorf("SourceColumns %+v missing secret_col", span.Results[0].SourceColumns)
	}
}

func TestGetQuerySpan_EmptyStringAlias(t *testing.T) {
	// SELECT c AS '' names the result column with the explicit empty alias;
	// it must not fall back to the expression name.
	span, err := GetQuerySpan("SELECT c AS '' FROM t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 1 || span.Results[0].Name != "" {
		t.Fatalf("Results = %+v, want one column named \"\"", span.Results)
	}
}

func TestGetQuerySpan_QualifiedSubscriptKeepsTableAccess(t *testing.T) {
	// The qualified twin of the subscript fail-open: the select-item fast
	// path used to bypass the expression parser for t.col and drop the rest.
	span, err := GetQuerySpan("SELECT t.secret_col[1] FROM sensitive_table t")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.AccessTables) != 1 || span.AccessTables[0].Table != "sensitive_table" {
		t.Fatalf("AccessTables = %+v, want [sensitive_table]", span.AccessTables)
	}
}

func TestGetQuerySpan_TableFunctionIsNotTableAccess(t *testing.T) {
	// A table-valued function is not a physical table: BACKENDS() must not
	// surface as an accessed table for authorization or lineage.
	span, err := GetQuerySpan("SELECT * FROM BACKENDS()")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.AccessTables) != 0 {
		t.Fatalf("AccessTables = %+v, want none for a table function", span.AccessTables)
	}
}

func TestGetQuerySpan_FailsClosedOnParseError(t *testing.T) {
	// Engine-invalid SQL the parser used to swallow must now yield an error,
	// never a silently smaller span (BYT-10085).
	if _, err := GetQuerySpan("SELECT j->'$.a' FROM t"); err == nil {
		t.Fatal("GetQuerySpan accepted a statement with trailing junk")
	}

	// A subquery that does not fully parse fails the whole span: its table
	// reads would otherwise vanish from AccessTables.
	if _, err := GetQuerySpan("SELECT (SELECT j->'$.a' FROM t2) x FROM t1"); err == nil {
		t.Fatal("GetQuerySpan accepted an unparseable subquery")
	}

	// Empty input keeps the zero-span contract.
	span, err := GetQuerySpan("")
	if err != nil || span == nil || span.Type != QueryTypeUnknown {
		t.Fatalf("empty input: span=%+v err=%v, want zero span and nil error", span, err)
	}
}

func TestGetQuerySpan_TableNamesResemblingKeywords(t *testing.T) {
	// The FROM-subquery detection is an explicit AST discriminator, not a
	// prefix heuristic: tables named selected/within (and a quoted `select`)
	// are ordinary tables, not query text to re-parse.
	for _, tc := range []struct{ sql, table string }{
		{"SELECT * FROM selected", "selected"},
		{"SELECT * FROM within", "within"},
		{"SELECT * FROM `select`", "select"},
	} {
		span, err := GetQuerySpan(tc.sql)
		if err != nil {
			t.Fatalf("GetQuerySpan(%q) error: %v", tc.sql, err)
		}
		if len(span.AccessTables) != 1 || span.AccessTables[0].Table != tc.table {
			t.Errorf("GetQuerySpan(%q) AccessTables = %+v, want [%s]", tc.sql, span.AccessTables, tc.table)
		}
	}
}

func TestGetQuerySpan_DMLSubqueriesFailClosed(t *testing.T) {
	// The fail-closed contract holds for statements that produce no lineage
	// of their own: a malformed subquery inside DML surfaces as an error.
	if _, err := GetQuerySpan("DELETE FROM t WHERE id IN (SELECT 1 */ 2 FROM secret)"); err == nil {
		t.Fatal("malformed DML subquery accepted")
	}
	// A well-formed one contributes its table reads.
	span, err := GetQuerySpan("DELETE FROM t WHERE id IN (SELECT id FROM other)")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	found := false
	for _, a := range span.AccessTables {
		if a.Table == "other" {
			found = true
		}
	}
	if !found {
		t.Errorf("AccessTables = %+v, want to include other", span.AccessTables)
	}
}

func TestGetQuerySpan_SubqueryErrorLocationsAreOuter(t *testing.T) {
	// Errors from re-parsing a subquery's raw text must be shifted into the
	// outer statement's coordinates so diagnostics highlight the right spot.
	sql := "SELECT (SELECT 1 */ 2 FROM t2) x FROM t1"
	_, err := GetQuerySpan(sql)
	if err == nil {
		t.Fatal("malformed subquery accepted")
	}
	pe, ok := err.(*parser.ParseError)
	if !ok {
		t.Fatalf("err = %T, want *parser.ParseError", err)
	}
	want := strings.Index(sql, "*/")
	if pe.Loc.Start != want {
		t.Errorf("error Loc.Start = %d, want %d (the */ in the outer text)", pe.Loc.Start, want)
	}
}

func TestGetQuerySpan_TableFunctionArgSubqueriesFailClosed(t *testing.T) {
	// Subqueries embedded in table-function arguments are validated too: the
	// function call is walked, so a malformed one fails the span.
	if _, err := GetQuerySpan("SELECT * FROM numbers(EXISTS (SELECT 1 */ 2 FROM secret)) x"); err == nil {
		t.Fatal("malformed subquery in table-function argument accepted")
	}
}

func TestGetQuerySpan_NestedSubqueryErrorLocationsAccumulate(t *testing.T) {
	// Each nesting level's TextStart is relative to its own extracted text;
	// the walker must accumulate ancestor bases so a deeply nested error
	// still points into the original statement.
	sql := "SELECT (SELECT (SELECT 1 */ 2 FROM t3) FROM t2) FROM t1"
	_, err := GetQuerySpan(sql)
	if err == nil {
		t.Fatal("malformed nested subquery accepted")
	}
	pe, ok := err.(*parser.ParseError)
	if !ok {
		t.Fatalf("err = %T, want *parser.ParseError", err)
	}
	if want := strings.Index(sql, "*/"); pe.Loc.Start != want {
		t.Errorf("error Loc.Start = %d, want %d", pe.Loc.Start, want)
	}
}

func TestGetQuerySpan_NonQuerySubqueryFailsClosed(t *testing.T) {
	// A placeholder body that parses cleanly as a non-query statement is not
	// a valid subquery; accepting it would hide its reads from the span.
	if _, err := GetQuerySpan("SELECT EXISTS (DELETE FROM secret) FROM public"); err == nil {
		t.Fatal("EXISTS with a DML body accepted")
	}
}

func TestGetQuerySpan_EmptySubqueryPlaceholdersFailClosed(t *testing.T) {
	// Empty or comment-only placeholder bodies parse to zero statements and
	// used to bypass the query-node validation entirely.
	for _, sql := range []string{
		"SELECT EXISTS () FROM public",
		"SELECT EXISTS (/*comment*/) FROM public",
	} {
		if _, err := GetQuerySpan(sql); err == nil {
			t.Errorf("GetQuerySpan(%q) accepted an empty subquery body", sql)
		}
	}
}

func TestGetQuerySpan_InsertSelectRecordsReads(t *testing.T) {
	// A parsed query child inside DML carries real table reads: the SELECT
	// side of INSERT ... SELECT must land in AccessTables.
	span, err := GetQuerySpan("INSERT INTO dest SELECT * FROM secret")
	if err != nil {
		t.Fatalf("GetQuerySpan error: %v", err)
	}
	found := false
	for _, a := range span.AccessTables {
		if a.Table == "secret" {
			found = true
		}
	}
	if !found {
		t.Errorf("AccessTables = %+v, want to include secret", span.AccessTables)
	}
}
