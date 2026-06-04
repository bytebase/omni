package parser

import "testing"

// This file is the parser-select node's correctness slice for the set-operation
// layer (setops.go): UNION / INTERSECT / EXCEPT precedence & associativity (S1)
// and the four query primaries. The oracle differential is authoritative;
// structural tests pin the INTERSECT-binds-tighter and left-associative shapes.
// Helpers live in oracle_foundation_test.go and select_test.go.

// setopAcceptCorpus is the set-operation surface Trino 481 accepts.
var setopAcceptCorpus = []string{
	"SELECT 1 UNION SELECT 2",
	"SELECT 1 UNION ALL SELECT 2",
	"SELECT 1 UNION DISTINCT SELECT 2",
	"SELECT 1 INTERSECT SELECT 2",
	"SELECT 1 INTERSECT ALL SELECT 2",
	"SELECT 1 INTERSECT DISTINCT SELECT 2",
	"SELECT 1 EXCEPT SELECT 2",
	"SELECT 1 EXCEPT ALL SELECT 2",
	"SELECT 1 EXCEPT DISTINCT SELECT 2",
	// chains & precedence
	"SELECT 1 UNION SELECT 2 UNION SELECT 3",
	"SELECT 1 UNION SELECT 2 UNION ALL SELECT 3",
	"SELECT 1 INTERSECT SELECT 2 UNION SELECT 3",
	"SELECT 1 UNION SELECT 2 INTERSECT SELECT 3",
	"SELECT 1 EXCEPT SELECT 2 EXCEPT SELECT 3",
	// legacy examples
	"SELECT 123 UNION DISTINCT SELECT 123 UNION ALL SELECT 123",
	"SELECT 123 INTERSECT DISTINCT SELECT 123 INTERSECT ALL SELECT 123",
	// primaries as operands
	"TABLE foo UNION TABLE bar",
	"VALUES 1, 2 UNION VALUES 3",
	"VALUES (1) UNION ALL SELECT 2",
	// parenthesized operands & whole-query ORDER BY (S2)
	"(SELECT 1) UNION (SELECT 2)",
	"SELECT 1 UNION (SELECT 2 INTERSECT SELECT 3)",
	"(SELECT 1 UNION SELECT 2) INTERSECT SELECT 3",
	"SELECT 1 UNION SELECT 2 ORDER BY 1",
	"SELECT 1 UNION SELECT 2 ORDER BY 1 LIMIT 5",
	"(SELECT 1 ORDER BY 1) UNION (SELECT 2 LIMIT 3)",
}

// setopRejectCorpus is malformed set-operation input Trino 481 rejects.
var setopRejectCorpus = []string{
	"SELECT 1 UNION",                // no right operand
	"SELECT 1 INTERSECT",            // no right operand
	"SELECT 1 EXCEPT",               // no right operand
	"UNION SELECT 1",                // no left operand
	"SELECT 1 UNION ALL",            // quantifier but no operand
	"SELECT 1 UNION UNION SELECT 2", // doubled operator
	// a parenthesized query PRIMARY is `( queryNoWith )` — WITH is NOT allowed
	// there (a WITH must be a relation subquery `( query )`). Oracle-confirmed.
	"(WITH c AS (SELECT 1) SELECT * FROM c)",
	"(WITH c AS (SELECT 1) SELECT * FROM c) UNION SELECT 2",
}

func TestSetop_AcceptCorpusParses(t *testing.T) {
	for _, sql := range setopAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) != 0 {
				t.Errorf("Parse(%q) should accept, got: %v", sql, errs)
			}
		})
	}
}

func TestSetop_RejectCorpusRejected(t *testing.T) {
	for _, sql := range setopRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) == 0 {
				t.Errorf("Parse(%q) should reject, but accepted", sql)
			}
		})
	}
}

func TestSetop_OracleDifferential(t *testing.T) {
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
	for _, sql := range setopAcceptCorpus {
		t.Run("accept/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
	for _, sql := range setopRejectCorpus {
		t.Run("reject/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
}

// ---------------------------------------------------------------------------
// structural tests (precedence & associativity, S1)
// ---------------------------------------------------------------------------

func TestSetop_StructureUnionLeftAssoc(t *testing.T) {
	// `a UNION b UNION c` groups left: top.Left is itself a SetOperation.
	q := parseOneQuery(t, "SELECT 1 UNION SELECT 2 UNION SELECT 3")
	top, ok := q.Body.(*SetOperation)
	if !ok {
		t.Fatalf("body is %T, want *SetOperation", q.Body)
	}
	if top.Op != "UNION" {
		t.Errorf("top op=%q, want UNION", top.Op)
	}
	if _, ok := top.Left.(*SetOperation); !ok {
		t.Errorf("UNION not left-associative: Left is %T, want *SetOperation", top.Left)
	}
	if _, ok := top.Right.(*QuerySpec); !ok {
		t.Errorf("top.Right is %T, want *QuerySpec", top.Right)
	}
}

func TestSetop_StructureIntersectBindsTighter(t *testing.T) {
	// `a INTERSECT b UNION c` == `(a INTERSECT b) UNION c`: the top op is UNION,
	// its Left is the INTERSECT.
	q := parseOneQuery(t, "SELECT 1 INTERSECT SELECT 2 UNION SELECT 3")
	top := q.Body.(*SetOperation)
	if top.Op != "UNION" {
		t.Fatalf("top op=%q, want UNION", top.Op)
	}
	left, ok := top.Left.(*SetOperation)
	if !ok || left.Op != "INTERSECT" {
		t.Errorf("top.Left = %T (op?), want INTERSECT SetOperation", top.Left)
	}

	// `a UNION b INTERSECT c` == `a UNION (b INTERSECT c)`: top is UNION, Right is
	// the INTERSECT.
	q = parseOneQuery(t, "SELECT 1 UNION SELECT 2 INTERSECT SELECT 3")
	top = q.Body.(*SetOperation)
	if top.Op != "UNION" {
		t.Fatalf("top op=%q, want UNION", top.Op)
	}
	right, ok := top.Right.(*SetOperation)
	if !ok || right.Op != "INTERSECT" {
		t.Errorf("top.Right = %T, want INTERSECT SetOperation", top.Right)
	}
}

func TestSetop_StructureQuantifier(t *testing.T) {
	q := parseOneQuery(t, "SELECT 1 UNION ALL SELECT 2")
	if op := q.Body.(*SetOperation); op.Quantifier != "ALL" {
		t.Errorf("Quantifier=%q, want ALL", op.Quantifier)
	}
	q = parseOneQuery(t, "SELECT 1 INTERSECT DISTINCT SELECT 2")
	if op := q.Body.(*SetOperation); op.Quantifier != "DISTINCT" {
		t.Errorf("Quantifier=%q, want DISTINCT", op.Quantifier)
	}
	q = parseOneQuery(t, "SELECT 1 UNION SELECT 2")
	if op := q.Body.(*SetOperation); op.Quantifier != "" {
		t.Errorf("Quantifier=%q, want empty", op.Quantifier)
	}
}

func TestSetop_StructureOrderByAttachesToWhole(t *testing.T) {
	// S2: `SELECT 1 UNION SELECT 2 ORDER BY 1` — the ORDER BY is on the queryNoWith
	// (the whole union), NOT on the right SELECT.
	q := parseOneQuery(t, "SELECT 1 UNION SELECT 2 ORDER BY 1")
	if _, ok := q.Body.(*SetOperation); !ok {
		t.Fatalf("body is %T, want *SetOperation", q.Body)
	}
	if len(q.OrderBy) != 1 {
		t.Errorf("query ORDER BY items=%d, want 1 (attached to the whole union)", len(q.OrderBy))
	}
}

func TestSetop_StructurePrimaries(t *testing.T) {
	q := parseOneQuery(t, "TABLE foo UNION TABLE bar")
	top := q.Body.(*SetOperation)
	if _, ok := top.Left.(*TableQuery); !ok {
		t.Errorf("Left is %T, want *TableQuery", top.Left)
	}
	if _, ok := top.Right.(*TableQuery); !ok {
		t.Errorf("Right is %T, want *TableQuery", top.Right)
	}

	q = parseOneQuery(t, "(SELECT 1) UNION (SELECT 2)")
	top = q.Body.(*SetOperation)
	if _, ok := top.Left.(*ParenQuery); !ok {
		t.Errorf("Left is %T, want *ParenQuery", top.Left)
	}
}
