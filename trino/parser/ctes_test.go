package parser

import "testing"

// This file is the parser-select node's correctness slice for the CTE grammar
// (ctes.go): WITH / WITH RECURSIVE / namedQuery, including the WITH SESSION
// divergence (D2). The oracle differential is authoritative; structural tests
// pin the CTE list / column-alias / recursive shapes. Helpers live in
// oracle_foundation_test.go and select_test.go.

// cteAcceptCorpus is the WITH / CTE surface Trino 481 accepts.
var cteAcceptCorpus = []string{
	"WITH t AS (SELECT 1) SELECT * FROM t",
	"WITH t AS (SELECT 1) TABLE t",
	"WITH t (a, b) AS (SELECT 1, 2) SELECT * FROM t",
	"WITH a AS (SELECT 1), b AS (SELECT 2) SELECT * FROM a, b",
	"WITH a AS (SELECT 1), b AS (SELECT 2), c AS (SELECT 3) SELECT * FROM a, b, c",
	"WITH RECURSIVE t(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM t WHERE n<5) SELECT * FROM t",
	"WITH RECURSIVE a AS (SELECT * FROM x) TABLE y",
	"WITH a (t, u) AS (SELECT * FROM x), b AS (SELECT * FROM y) TABLE z",
	// a CTE whose body is itself a set operation / has ORDER BY / its own WITH
	"WITH t AS (SELECT 1 UNION SELECT 2) SELECT * FROM t",
	"WITH t AS (SELECT 1 ORDER BY 1 LIMIT 1) SELECT * FROM t",
	"WITH outer_cte AS (WITH inner_cte AS (SELECT 1) SELECT * FROM inner_cte) SELECT * FROM outer_cte",
	// WITH feeding a set operation / VALUES body
	"WITH t AS (SELECT 1) SELECT * FROM t UNION SELECT 2",
	"WITH t AS (VALUES 1, 2) SELECT * FROM t",
}

// cteRejectCorpus is malformed WITH input Trino 481 rejects with a SYNTAX_ERROR.
// Includes the WITH SESSION divergence (D2): omni rejects it (legacy grammar
// scope) and Trino 481 ALSO rejects `WITH SESSION` as a *CTE* — the 481 SESSION
// prefix is `WITH SESSION prop = v <query>` which is exercised in the divergence
// test, not here.
var cteRejectCorpus = []string{
	"WITH SELECT 1",                           // WITH with no CTE before the query
	"WITH t SELECT 1",                         // CTE name with no AS (…)
	"WITH t AS SELECT 1",                      // CTE body without parentheses
	"WITH t AS () SELECT 1",                   // empty CTE body
	"WITH t AS (SELECT 1)",                    // CTE with no trailing query
	"WITH t AS (SELECT 1),",                   // trailing comma, no second CTE
	"WITH RECURSIVE SELECT 1",                 // RECURSIVE with no CTE
	"WITH t () AS (SELECT 1) SELECT * FROM t", // empty column-alias list
}

func TestCTE_AcceptCorpusParses(t *testing.T) {
	for _, sql := range cteAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) != 0 {
				t.Errorf("Parse(%q) should accept, got: %v", sql, errs)
			}
		})
	}
}

func TestCTE_RejectCorpusRejected(t *testing.T) {
	for _, sql := range cteRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) == 0 {
				t.Errorf("Parse(%q) should reject, but accepted", sql)
			}
		})
	}
}

func TestCTE_OracleDifferential(t *testing.T) {
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
	for _, sql := range cteAcceptCorpus {
		t.Run("accept/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
	for _, sql := range cteRejectCorpus {
		t.Run("reject/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
}

// ---------------------------------------------------------------------------
// structural tests
// ---------------------------------------------------------------------------

func TestCTE_StructureWithList(t *testing.T) {
	q := parseOneQuery(t, "WITH a AS (SELECT 1), b (x, y) AS (SELECT 1, 2) SELECT * FROM a, b")
	if q.With == nil {
		t.Fatal("With is nil, want a CTE list")
	}
	if q.With.Recursive {
		t.Error("Recursive=true, want false")
	}
	if len(q.With.CTEs) != 2 {
		t.Fatalf("CTEs=%d, want 2", len(q.With.CTEs))
	}
	if q.With.CTEs[0].Name.Value != "a" || len(q.With.CTEs[0].ColumnAliases) != 0 {
		t.Errorf("CTE 0 = %+v, want name a no aliases", q.With.CTEs[0])
	}
	if q.With.CTEs[1].Name.Value != "b" || len(q.With.CTEs[1].ColumnAliases) != 2 {
		t.Errorf("CTE 1 = %+v, want name b with 2 aliases", q.With.CTEs[1])
	}
	if q.With.CTEs[0].Query == nil {
		t.Error("CTE 0 body is nil")
	}
}

func TestCTE_StructureRecursive(t *testing.T) {
	q := parseOneQuery(t, "WITH RECURSIVE t(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM t WHERE n<5) SELECT * FROM t")
	if q.With == nil || !q.With.Recursive {
		t.Fatalf("With=%+v, want Recursive", q.With)
	}
	if len(q.With.CTEs) != 1 {
		t.Fatalf("CTEs=%d, want 1", len(q.With.CTEs))
	}
	// The recursive CTE body is a set operation.
	if _, ok := q.With.CTEs[0].Query.Body.(*SetOperation); !ok {
		t.Errorf("CTE body is %T, want *SetOperation", q.With.CTEs[0].Query.Body)
	}
}

func TestCTE_StructureNestedWith(t *testing.T) {
	// A CTE whose body itself has a WITH.
	q := parseOneQuery(t, "WITH a AS (WITH b AS (SELECT 1) SELECT * FROM b) SELECT * FROM a")
	if q.With == nil || len(q.With.CTEs) != 1 {
		t.Fatalf("With=%+v, want one outer CTE", q.With)
	}
	if q.With.CTEs[0].Query.With == nil {
		t.Error("inner CTE query has no WITH, want a nested WITH")
	}
}
