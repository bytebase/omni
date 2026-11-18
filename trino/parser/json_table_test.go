package parser

import "testing"

// This file pins the JSON_TABLE relation primary (json_table.go): the
// jsonTableColumn alternatives, the specific and default plans, the path
// invocation's optional AS name, and the trailing ON ERROR behavior. The
// grammar-generated compat.tsv rows are the byte-level gate; these tests keep
// the accepted forms parsing and the AST shape readable.

var jsonTableAcceptCorpus = []string{
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a FOR ORDINALITY))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP PATH '$.a'))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP NULL ON EMPTY ERROR ON ERROR))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP FORMAT JSON))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP FORMAT JSON PATH '$.a' WITH ARRAY WRAPPER KEEP QUOTES ON SCALAR STRING EMPTY ARRAY ON EMPTY ERROR ON ERROR))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(NESTED '$.a' COLUMNS(a TIMESTAMP)))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(NESTED PATH '$.a' AS n COLUMNS(a TIMESTAMP, NESTED '$.b' COLUMNS(b TIMESTAMP))))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' AS p COLUMNS(a TIMESTAMP))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' PASSING 1 AS x COLUMNS(a TIMESTAMP))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' AS p PASSING 1 AS x COLUMNS(a TIMESTAMP))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p OUTER q))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p INNER (q CROSS r)))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p UNION q UNION r))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p CROSS q CROSS r))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT(OUTER))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT(INNER, UNION))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT(UNION, OUTER))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) ERROR ON ERROR)",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) EMPTY ON ERROR) jt",
}

var jsonTableRejectCorpus = []string{
	"SELECT * FROM JSON_TABLE('{}', 'lax $')",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS())",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a FOR))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(NESTED '$.a'))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN())",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p OUTER))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN(p UNION))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT())",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT(OUTER, OUTER))",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) ON ERROR)",
	"SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) ERROR ON)",
}

func TestJSONTable_AcceptCorpusParses(t *testing.T) {
	for _, sql := range jsonTableAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := parseForTest(sql); len(errs) != 0 {
				t.Errorf("parseForTest(%q) should accept, got: %v", sql, errs)
			}
		})
	}
}

func TestJSONTable_RejectCorpusRejected(t *testing.T) {
	for _, sql := range jsonTableRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := parseForTest(sql); len(errs) == 0 {
				t.Errorf("parseForTest(%q) should reject, but accepted", sql)
			}
		})
	}
}

func TestJSONTable_Structure(t *testing.T) {
	sql := "SELECT * FROM JSON_TABLE('{}', 'lax $' AS root COLUMNS(" +
		"ord FOR ORDINALITY, " +
		"v VARCHAR PATH '$.v' NULL ON EMPTY DEFAULT 'd' ON ERROR, " +
		"q TIMESTAMP FORMAT JSON PATH '$.q' WITH ARRAY WRAPPER KEEP QUOTES ON SCALAR STRING ERROR ON EMPTY EMPTY ARRAY ON ERROR, " +
		"NESTED PATH '$.n' AS n COLUMNS(nv TIMESTAMP)" +
		") PLAN(root OUTER (a UNION b)) EMPTY ON ERROR)"
	rel := unwrapAlias(fromRel(t, sql))
	jt, ok := rel.(*JSONTableRelation)
	if !ok {
		t.Fatalf("relation is %T, want *JSONTableRelation", rel)
	}
	if jt.Path.PathName == nil || jt.Path.PathName.Value != "root" {
		t.Errorf("Path.PathName=%v, want root", jt.Path.PathName)
	}
	if jt.OnError != "EMPTY" {
		t.Errorf("OnError=%q, want EMPTY", jt.OnError)
	}
	if len(jt.Columns) != 4 {
		t.Fatalf("want 4 columns, got %d", len(jt.Columns))
	}
	if jt.Columns[0].Kind != JSONTableOrdinalityColumn || jt.Columns[0].Name.Value != "ord" {
		t.Errorf("column 0=%+v, want ordinality ord", jt.Columns[0])
	}
	v := jt.Columns[1]
	if v.Kind != JSONTableValueColumn || v.Path == nil || *v.Path != "$.v" || v.OnEmpty == nil || v.OnEmpty.Kind != "NULL" || v.OnError == nil || v.OnError.Kind != "DEFAULT" {
		t.Errorf("column 1=%+v, want value column with PATH/NULL ON EMPTY/DEFAULT ON ERROR", v)
	}
	q := jt.Columns[2]
	if q.Kind != JSONTableQueryColumn || q.Format == nil || q.Wrapper != "WITH" || !q.WrapperArray || q.Quotes != "KEEP" || !q.QuotesOnScalar || q.QueryOnEmpty == nil || q.QueryOnEmpty.Kind != "ERROR" || q.QueryOnError == nil || q.QueryOnError.Kind != "EMPTY ARRAY" {
		t.Errorf("column 2=%+v, want query column with every clause", q)
	}
	if q.Loc.End <= q.QueryOnError.Loc.End || v.Loc.End <= v.OnError.Loc.End {
		t.Errorf("column locations must include their behavior's ON EMPTY/ON ERROR tails: query=%+v value=%+v", q.Loc, v.Loc)
	}
	n := jt.Columns[3]
	if n.Kind != JSONTableNestedColumn || n.Path == nil || *n.Path != "$.n" || n.Alias == nil || n.Alias.Value != "n" || len(n.Columns) != 1 {
		t.Errorf("column 3=%+v, want nested columns with alias n", n)
	}
	plan := jt.Plan
	if plan == nil || plan.Kind != JSONTableJoinPlan || plan.PathName.Value != "root" || plan.JoinType != "OUTER" || plan.Right == nil || plan.Right.Plan == nil || plan.Right.Plan.Kind != JSONTableUnionPlan {
		t.Errorf("Plan=%+v, want root OUTER (a UNION b)", plan)
	}
}

func TestJSONTable_DefaultPlanStructure(t *testing.T) {
	rel := unwrapAlias(fromRel(t, "SELECT * FROM JSON_TABLE('{}', 'lax $' COLUMNS(a TIMESTAMP) PLAN DEFAULT(UNION, OUTER))"))
	jt := rel.(*JSONTableRelation)
	if jt.Plan != nil || jt.DefaultPlan == nil || jt.DefaultPlan.First != "UNION" || jt.DefaultPlan.Second != "OUTER" {
		t.Errorf("Plan=%v DefaultPlan=%+v, want default UNION, OUTER", jt.Plan, jt.DefaultPlan)
	}
}
