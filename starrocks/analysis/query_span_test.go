package analysis

import (
	"testing"
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
			name:       "qualified table",
			sql:        "SELECT * FROM db.t",
			wantType:   QueryTypeSelect,
			wantTables: []tableSig{{Database: "db", Table: "t"}},
			wantResults: []string{"*"},
		},
		{
			name:       "three-part table name",
			sql:        "SELECT * FROM catalog1.db1.t1",
			wantType:   QueryTypeSelect,
			wantTables: []tableSig{{Database: "db1", Table: "t1"}},
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
