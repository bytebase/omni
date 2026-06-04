package parser

import "testing"

// This file is the parser-select node's correctness slice for the FROM-clause
// relation grammar (relation.go): joins, sampledRelation/TABLESAMPLE,
// aliasedRelation, relationPrimary (table / subquery / UNNEST / LATERAL /
// TABLE(func) / parenthesized / queryPeriod), and tableFunctionCall. The oracle
// differential is authoritative; structural tests pin join associativity and
// relation shapes. Helpers (connectOracle / oracleAccepts / truncateName) live
// in oracle_foundation_test.go; parseOneQuery / querySpec live in select_test.go.

// relationAcceptCorpus is the FROM-clause relation surface Trino 481 accepts.
// Each is a complete statement (the differential needs a whole statement);
// table/column references resolve to semantic errors, not syntax rejections.
var relationAcceptCorpus = []string{
	// --- comma (implicit cross) joins ---
	"SELECT * FROM a, b",
	"SELECT * FROM a, b, c",

	// --- qualified joins + criteria ---
	"SELECT * FROM a JOIN b ON a.x = b.x",
	"SELECT * FROM a INNER JOIN b ON a.x = b.x",
	"SELECT * FROM a LEFT JOIN b ON a.x = b.x",
	"SELECT * FROM a LEFT OUTER JOIN b ON a.x = b.x",
	"SELECT * FROM a RIGHT JOIN b ON true",
	"SELECT * FROM a RIGHT OUTER JOIN b ON true",
	"SELECT * FROM a FULL JOIN b ON true",
	"SELECT * FROM a FULL OUTER JOIN b ON true",
	"SELECT * FROM a JOIN b USING (x)",
	"SELECT * FROM a JOIN b USING (x, y)",

	// --- cross / natural joins ---
	"SELECT * FROM a CROSS JOIN b",
	"SELECT * FROM a NATURAL JOIN b",
	"SELECT * FROM a NATURAL LEFT JOIN b",
	"SELECT * FROM a NATURAL FULL OUTER JOIN b",

	// --- join chains (left-associative) + the legacy join_precedence examples ---
	"SELECT * FROM a JOIN b ON a.x=b.x JOIN c ON b.y=c.y",
	"SELECT * FROM a CROSS JOIN b LEFT JOIN c ON true",
	"SELECT * FROM a CROSS JOIN b NATURAL JOIN c CROSS JOIN d NATURAL JOIN e",
	"SELECT * FROM t1 a JOIN t2 b ON a.id=b.id JOIN t3 c ON b.id=c.id",

	// --- parenthesized relation (explicit grouping) ---
	"SELECT * FROM (a JOIN b ON a.x=b.x)",
	"SELECT * FROM (a CROSS JOIN b) JOIN c ON true",
	// nested parens: `((join))` is a parenthesized relation, NOT a subquery — the
	// deciding token lies arbitrarily deep, so it is resolved by speculation.
	"SELECT * FROM ((a JOIN b ON true) JOIN c ON true)",
	"SELECT * FROM (((a)))",
	"SELECT * FROM ((SELECT 1))", // nested parens that ARE a subquery

	// --- aliased relations ---
	"SELECT * FROM t a",
	"SELECT * FROM t AS a",
	"SELECT * FROM t AS a (x, y)",
	"SELECT * FROM t a (x, y)",
	"SELECT * FROM nation AS n CROSS JOIN region AS r",

	// --- subquery relations ---
	"SELECT * FROM (SELECT 1)",
	"SELECT * FROM (SELECT 1) t",
	"SELECT * FROM (SELECT 1) AS t (x)",
	"SELECT * FROM (SELECT a FROM u WHERE a > 1) t",
	"SELECT * FROM (VALUES (1, '1'), (2, '2'))",
	"SELECT * FROM (TABLE foo)",
	"SELECT * FROM (SELECT 1 UNION SELECT 2) t",
	// a relation subquery is `( query )` — WITH is allowed inside it (unlike the
	// parenthesized query PRIMARY `( queryNoWith )`, which rejects WITH; see
	// setopRejectCorpus). Oracle-confirmed.
	"SELECT * FROM (WITH c AS (SELECT 1) SELECT * FROM c)",

	// --- UNNEST ---
	"SELECT * FROM UNNEST(ARRAY[1,2])",
	"SELECT * FROM UNNEST(ARRAY[1,2]) WITH ORDINALITY",
	"SELECT * FROM UNNEST(ARRAY[1,2]) AS t(number)",
	"SELECT * FROM UNNEST(a, b) WITH ORDINALITY AS t(x, y, ord)",
	"SELECT * FROM t CROSS JOIN UNNEST(a)",
	"SELECT * FROM t CROSS JOIN UNNEST(a, b) WITH ORDINALITY",
	"SELECT * FROM t FULL JOIN UNNEST(a) AS tmp (c) ON true",

	// --- LATERAL ---
	"SELECT * FROM LATERAL (SELECT 1)",
	"SELECT * FROM t, LATERAL (VALUES 1) a(x)",
	"SELECT * FROM t CROSS JOIN LATERAL (VALUES 1)",
	"SELECT * FROM t FULL JOIN LATERAL (VALUES 1) ON true",
	// LATERAL is non-reserved: not followed by '(' it is an ordinary table name.
	"SELECT * FROM LATERAL",
	"SELECT * FROM LATERAL x",

	// --- TABLESAMPLE ---
	"SELECT * FROM t TABLESAMPLE BERNOULLI (10)",
	"SELECT * FROM t TABLESAMPLE SYSTEM (50)",
	"SELECT * FROM t a TABLESAMPLE BERNOULLI (10)",

	// --- queryPeriod (time travel) ---
	"SELECT * FROM t FOR TIMESTAMP AS OF TIMESTAMP '2020-01-01 00:00:00'",
	"SELECT * FROM t FOR VERSION AS OF 3",
	"SELECT * FROM cat.sch.t FOR VERSION AS OF 'snap-1'",

	// --- table function invocation ---
	"SELECT * FROM TABLE(my_function(1, 100))",
	"SELECT * FROM TABLE(my_function(row_count => 100, column_count => 1))",
	"SELECT * FROM TABLE(schema_name.my_function(1, 100))",
	"SELECT * FROM TABLE(catalog_name.schema_name.my_function(1, 100))",
	"SELECT * FROM TABLE(sequence(start => 1000000, stop => -2000000, step => -3))",
	"SELECT * FROM TABLE(exclude_columns(input => TABLE(orders), columns => DESCRIPTOR(clerk, comment)))",
	"SELECT * FROM TABLE(my_function(input => TABLE(orders) PARTITION BY orderstatus ORDER BY orderdate))",
	"SELECT * FROM TABLE(my_function(input => TABLE(SELECT * FROM orders) PARTITION BY (a, b) PRUNE WHEN EMPTY))",
	"SELECT * FROM TABLE(f(TABLE(orders) KEEP WHEN EMPTY))",
	"SELECT * FROM TABLE(f(d => DESCRIPTOR(a integer, b varchar)))",
	"SELECT * FROM TABLE(f(d => CAST(NULL AS DESCRIPTOR)))",
	"SELECT * FROM TABLE(f(TABLE(a) PARTITION BY x, TABLE(b) PARTITION BY y COPARTITION (a, b)))",
	// COPARTITION and DESCRIPTOR are NON-RESERVED, so they are also ordinary
	// argument names / column references (not only the clause / descriptor forms).
	// Oracle-confirmed; these were cross-review (Codex) catches.
	"SELECT * FROM TABLE(f(copartition => 1))",
	"SELECT * FROM TABLE(f(descriptor))",
	"SELECT * FROM TABLE(f(descriptor => 1))",
}

// relationRejectCorpus is malformed FROM-clause input Trino 481 rejects with a
// SYNTAX_ERROR (required negative coverage).
var relationRejectCorpus = []string{
	"SELECT * FROM a JOIN b",                // qualified join needs ON/USING
	"SELECT * FROM a INNER JOIN b",          // ditto
	"SELECT * FROM a JOIN b ON",             // ON with no predicate
	"SELECT * FROM a JOIN b USING",          // USING with no column list
	"SELECT * FROM a JOIN b USING ()",       // empty USING list
	"SELECT * FROM a CROSS b",               // CROSS without JOIN
	"SELECT * FROM ,",                       // empty relation list
	"SELECT * FROM a,",                      // trailing comma
	"SELECT * FROM UNNEST",                  // UNNEST (reserved) without parens
	"SELECT * FROM UNNEST()",                // UNNEST with no expression
	"SELECT * FROM LATERAL (a)",             // LATERAL ( … ) must wrap a query, not a name
	"SELECT * FROM t TABLESAMPLE (10)",      // TABLESAMPLE without method
	"SELECT * FROM t TABLESAMPLE BERNOULLI", // TABLESAMPLE without percentage
	"SELECT * FROM t FOR AS OF 3",           // queryPeriod without range type
	"SELECT * FROM t FOR VERSION 3",         // queryPeriod without AS OF
	"SELECT * FROM t AS",                    // AS with no alias
	// a COPARTITION group requires at least two tables (the first comma is
	// mandatory); a single-table group is a syntax error. Oracle-confirmed.
	"SELECT * FROM TABLE(f(TABLE(a) PARTITION BY x COPARTITION (a)))",
}

func TestRelation_AcceptCorpusParses(t *testing.T) {
	for _, sql := range relationAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) != 0 {
				t.Errorf("Parse(%q) should accept, got: %v", sql, errs)
			}
		})
	}
}

func TestRelation_RejectCorpusRejected(t *testing.T) {
	for _, sql := range relationRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) == 0 {
				t.Errorf("Parse(%q) should reject, but accepted", sql)
			}
		})
	}
}

func TestRelation_OracleDifferential(t *testing.T) {
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
	for _, sql := range relationAcceptCorpus {
		t.Run("accept/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
	for _, sql := range relationRejectCorpus {
		t.Run("reject/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
}

// ---------------------------------------------------------------------------
// structural tests
// ---------------------------------------------------------------------------

// fromRel extracts the single FROM relation of a plain query specification.
func fromRel(t *testing.T, sql string) Relation {
	t.Helper()
	spec := querySpec(t, parseOneQuery(t, sql))
	if len(spec.From) != 1 {
		t.Fatalf("Parse(%q): want 1 FROM relation, got %d", sql, len(spec.From))
	}
	return spec.From[0]
}

// unwrapAlias returns the inner relation primary of an AliasedRelation (the
// shape parseAliasedRelation always produces), or the relation itself otherwise.
func unwrapAlias(r Relation) Relation {
	if ar, ok := r.(*AliasedRelation); ok {
		return ar.Inner
	}
	return r
}

func TestRelation_StructureJoinLeftAssoc(t *testing.T) {
	// `a JOIN b ... JOIN c ...` groups left: the top Join's Left is itself a Join.
	rel := fromRel(t, "SELECT * FROM a JOIN b ON a.x=b.x JOIN c ON b.y=c.y")
	top, ok := rel.(*Join)
	if !ok {
		t.Fatalf("top relation is %T, want *Join", rel)
	}
	if _, ok := top.Left.(*Join); !ok {
		t.Errorf("join not left-associative: Left is %T, want *Join", top.Left)
	}
}

func TestRelation_StructureJoinTypes(t *testing.T) {
	cases := []struct {
		sql   string
		jt    JoinType
		outer bool
	}{
		{"SELECT * FROM a JOIN b ON true", JoinInner, false},
		{"SELECT * FROM a INNER JOIN b ON true", JoinInner, false},
		{"SELECT * FROM a LEFT JOIN b ON true", JoinLeft, false},
		{"SELECT * FROM a LEFT OUTER JOIN b ON true", JoinLeft, true},
		{"SELECT * FROM a RIGHT JOIN b ON true", JoinRight, false},
		{"SELECT * FROM a FULL OUTER JOIN b ON true", JoinFull, true},
		{"SELECT * FROM a CROSS JOIN b", JoinCross, false},
	}
	for _, c := range cases {
		t.Run(truncateName(c.sql), func(t *testing.T) {
			j, ok := fromRel(t, c.sql).(*Join)
			if !ok {
				t.Fatalf("relation is %T, want *Join", fromRel(t, c.sql))
			}
			if j.Type != c.jt {
				t.Errorf("Type=%v, want %v", j.Type, c.jt)
			}
			if j.Outer != c.outer {
				t.Errorf("Outer=%v, want %v", j.Outer, c.outer)
			}
		})
	}
}

func TestRelation_StructureJoinCriteria(t *testing.T) {
	j := fromRel(t, "SELECT * FROM a JOIN b ON a.x = b.x").(*Join)
	if j.On == nil || j.Using != nil {
		t.Errorf("ON join: On=%v Using=%v, want On set", j.On, j.Using)
	}
	j = fromRel(t, "SELECT * FROM a JOIN b USING (x, y)").(*Join)
	if j.On != nil || len(j.Using) != 2 {
		t.Errorf("USING join: On=%v Using=%v, want 2-col Using", j.On, j.Using)
	}
	j = fromRel(t, "SELECT * FROM a NATURAL LEFT JOIN b").(*Join)
	if !j.Natural || j.Type != JoinLeft {
		t.Errorf("natural left join: Natural=%v Type=%v", j.Natural, j.Type)
	}
}

func TestRelation_StructurePrimaries(t *testing.T) {
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM t")).(*TableRelation); !ok {
		t.Errorf("`FROM t`: want *TableRelation")
	}
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM (SELECT 1) x")).(*SubqueryRelation); !ok {
		t.Errorf("`FROM (SELECT 1) x`: want *SubqueryRelation")
	}
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM UNNEST(a) WITH ORDINALITY")).(*UnnestRelation); !ok {
		t.Errorf("`FROM UNNEST(a)`: want *UnnestRelation")
	}
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM LATERAL (SELECT 1)")).(*LateralRelation); !ok {
		t.Errorf("`FROM LATERAL (...)`: want *LateralRelation")
	}
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM TABLE(f(1))")).(*TableFunctionRelation); !ok {
		t.Errorf("`FROM TABLE(f(1))`: want *TableFunctionRelation")
	}
	if _, ok := unwrapAlias(fromRel(t, "SELECT * FROM (a JOIN b ON true)")).(*ParenRelation); !ok {
		t.Errorf("`FROM (a JOIN b)`: want *ParenRelation")
	}
}

func TestRelation_StructureAliasAndSample(t *testing.T) {
	ar, ok := fromRel(t, "SELECT * FROM t AS a (x, y) TABLESAMPLE BERNOULLI (10)").(*AliasedRelation)
	if !ok {
		t.Fatalf("want *AliasedRelation")
	}
	if ar.Alias == nil || ar.Alias.Value != "a" {
		t.Errorf("Alias=%v, want a", ar.Alias)
	}
	if len(ar.ColumnAliases) != 2 {
		t.Errorf("ColumnAliases=%d, want 2", len(ar.ColumnAliases))
	}
	if ar.Sample == nil || ar.Sample.Method != "BERNOULLI" {
		t.Errorf("Sample=%v, want BERNOULLI", ar.Sample)
	}
}

func TestRelation_StructureUnnest(t *testing.T) {
	ur := unwrapAlias(fromRel(t, "SELECT * FROM UNNEST(a, b) WITH ORDINALITY AS t(x,y,o)")).(*UnnestRelation)
	if len(ur.Exprs) != 2 {
		t.Errorf("UNNEST exprs=%d, want 2", len(ur.Exprs))
	}
	if !ur.WithOrdinality {
		t.Error("WithOrdinality=false, want true")
	}
}

func TestRelation_StructureQueryPeriod(t *testing.T) {
	tr := unwrapAlias(fromRel(t, "SELECT * FROM t FOR VERSION AS OF 3")).(*TableRelation)
	if tr.Period == nil || tr.Period.RangeType != "VERSION" {
		t.Errorf("Period=%v, want VERSION", tr.Period)
	}
	tr = unwrapAlias(fromRel(t, "SELECT * FROM t FOR TIMESTAMP AS OF TIMESTAMP '2020-01-01 00:00:00'")).(*TableRelation)
	if tr.Period == nil || tr.Period.RangeType != "TIMESTAMP" {
		t.Errorf("Period=%v, want TIMESTAMP", tr.Period)
	}
}

func TestRelation_StructureTableFunction(t *testing.T) {
	tf := unwrapAlias(fromRel(t,
		"SELECT * FROM TABLE(exclude_columns(input => TABLE(orders), columns => DESCRIPTOR(clerk, comment)))")).(*TableFunctionRelation)
	if tf.Call.Name.Normalize() != "exclude_columns" {
		t.Errorf("func name=%q, want exclude_columns", tf.Call.Name.Normalize())
	}
	if len(tf.Call.Args) != 2 {
		t.Fatalf("args=%d, want 2", len(tf.Call.Args))
	}
	if tf.Call.Args[0].Name == nil || tf.Call.Args[0].Name.Value != "input" || tf.Call.Args[0].Kind != TFArgTable {
		t.Errorf("arg0 = %+v, want named 'input' table arg", tf.Call.Args[0])
	}
	if tf.Call.Args[1].Kind != TFArgDescriptor {
		t.Errorf("arg1 kind=%v, want TFArgDescriptor", tf.Call.Args[1].Kind)
	}
	if tf.Call.Args[1].Descriptor == nil || len(tf.Call.Args[1].Descriptor.Fields) != 2 {
		t.Errorf("descriptor fields = %+v, want 2", tf.Call.Args[1].Descriptor)
	}
}

func TestRelation_StructureCopartition(t *testing.T) {
	tf := unwrapAlias(fromRel(t,
		"SELECT * FROM TABLE(f(TABLE(a) PARTITION BY x, TABLE(b) PARTITION BY y COPARTITION (a, b)))")).(*TableFunctionRelation)
	if len(tf.Call.Copartition) != 1 {
		t.Fatalf("copartition groups=%d, want 1", len(tf.Call.Copartition))
	}
	if len(tf.Call.Copartition[0]) != 2 {
		t.Errorf("copartition group tables=%d, want 2", len(tf.Call.Copartition[0]))
	}
}
