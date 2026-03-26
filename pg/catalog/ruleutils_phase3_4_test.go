package catalog

import (
	"strings"
	"testing"
)

// =============================================================================
// Phase 3: JSON Expression Types (PG16+)
// =============================================================================

func setupJsonTest(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE jtest (id int PRIMARY KEY, data jsonb, name text, val int);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return c
}

// --- 3.1 JSON Constructor and Predicate ---

func TestRuleutils3_JsonIsPredicate(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest WHERE data IS JSON;`)
	if !strings.Contains(def, "IS JSON") {
		t.Errorf("expected IS JSON in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonIsPredicateObject(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest WHERE data IS JSON OBJECT;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "IS JSON") {
		t.Errorf("expected IS JSON in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonIsPredicateArray(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest WHERE data IS JSON ARRAY;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "IS JSON") {
		t.Errorf("expected IS JSON in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonIsPredicateScalar(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest WHERE data IS JSON SCALAR;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "IS JSON") {
		t.Errorf("expected IS JSON in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonValueFunc(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT JSON_VALUE(data, '$.key') FROM jtest;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "JSON_VALUE") {
		t.Errorf("expected JSON_VALUE in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonQueryFunc(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT JSON_QUERY(data, '$.arr') FROM jtest;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "JSON_QUERY") {
		t.Errorf("expected JSON_QUERY in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonExistsFunc(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest WHERE JSON_EXISTS(data, '$.key');`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "JSON_EXISTS") {
		t.Errorf("expected JSON_EXISTS in view def, got: %s", def)
	}
}

func TestRuleutils3_JsonObjectConstructor(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT JSON_OBJECT('key': name) FROM jtest;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "JSON_OBJECT") {
		t.Errorf("expected JSON_OBJECT in view def, got: %s", def)
	}
}

// --- 3.2 JSON constructor/expr direct deparse tests (unit level) ---

func TestRuleutils3_JsonConstructorExprDeparseObject(t *testing.T) {
	c := New()
	expr := &JsonConstructorExprQ{
		ConstructorType: JsonConstructorObject,
		Args: []AnalyzedExpr{
			&ConstExpr{TypeOID: TEXTOID, Value: "key"},
			&ConstExpr{TypeOID: INT4OID, Value: "42"},
		},
		ResultType: JSONOID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "JSON_OBJECT") {
		t.Errorf("expected JSON_OBJECT, got: %s", result)
	}
}

func TestRuleutils3_JsonConstructorExprDeparseArray(t *testing.T) {
	c := New()
	expr := &JsonConstructorExprQ{
		ConstructorType: JsonConstructorArray,
		Args: []AnalyzedExpr{
			&ConstExpr{TypeOID: INT4OID, Value: "1"},
			&ConstExpr{TypeOID: INT4OID, Value: "2"},
			&ConstExpr{TypeOID: INT4OID, Value: "3"},
		},
		ResultType: JSONOID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "JSON_ARRAY") {
		t.Errorf("expected JSON_ARRAY, got: %s", result)
	}
}

func TestRuleutils3_JsonExprDeparseValue(t *testing.T) {
	c := New()
	expr := &JsonExprQ{
		Op:         JsonExprValue,
		Expr:       &VarExpr{RangeIdx: 0, AttNum: 1, TypeOID: JSONBOID},
		Path:       "$.key",
		ResultType: TEXTOID,
	}
	rte := &RangeTableEntry{
		Kind:     RTERelation,
		RelName:  "t",
		ERef:     "t",
		ColNames: []string{"data"},
	}
	result := c.DeparseExpr(expr, []*RangeTableEntry{rte}, true)
	if !strings.Contains(result, "JSON_VALUE") {
		t.Errorf("expected JSON_VALUE, got: %s", result)
	}
	if !strings.Contains(result, "$.key") {
		t.Errorf("expected path $.key, got: %s", result)
	}
}

func TestRuleutils3_JsonExprDeparseQuery(t *testing.T) {
	c := New()
	expr := &JsonExprQ{
		Op:         JsonExprQuery,
		Expr:       &ConstExpr{TypeOID: JSONBOID, Value: `{"arr":[1,2]}`},
		Path:       "$.arr",
		ResultType: JSONBOID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "JSON_QUERY") {
		t.Errorf("expected JSON_QUERY, got: %s", result)
	}
}

func TestRuleutils3_JsonExprDeparseExists(t *testing.T) {
	c := New()
	expr := &JsonExprQ{
		Op:         JsonExprExists,
		Expr:       &ConstExpr{TypeOID: JSONBOID, Value: `{"key":1}`},
		Path:       "$.key",
		ResultType: BOOLOID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "JSON_EXISTS") {
		t.Errorf("expected JSON_EXISTS, got: %s", result)
	}
}

func TestRuleutils3_JsonIsPredicateDeparseIsJson(t *testing.T) {
	c := New()
	expr := &JsonIsPredicateExpr{
		Expr:     &ConstExpr{TypeOID: TEXTOID, Value: "{}"},
		ItemType: 0, // ANY
		IsNot:    false,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "IS JSON") {
		t.Errorf("expected IS JSON, got: %s", result)
	}
}

func TestRuleutils3_JsonIsPredicateDeparseNotJsonObject(t *testing.T) {
	c := New()
	expr := &JsonIsPredicateExpr{
		Expr:     &ConstExpr{TypeOID: TEXTOID, Value: "[]"},
		ItemType: 1, // OBJECT
		IsNot:    true,
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "IS NOT JSON") {
		t.Errorf("expected IS NOT JSON, got: %s", result)
	}
	if !strings.Contains(result, "OBJECT") {
		t.Errorf("expected OBJECT, got: %s", result)
	}
}

func TestRuleutils3_JsonValueExprDeparse(t *testing.T) {
	c := New()
	expr := &JsonValueExprQ{
		Expr:       &ConstExpr{TypeOID: TEXTOID, Value: `{"a":1}`},
		ResultType: TEXTOID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if result == "" {
		t.Error("expected non-empty deparse for JsonValueExprQ")
	}
}

// =============================================================================
// Phase 4: Query Feature Completeness
// =============================================================================

// --- 4.1 FOR UPDATE/SHARE ---

func TestRuleutils4_ForUpdate(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR UPDATE;`)
	if !strings.Contains(strings.ToUpper(def), "FOR UPDATE") {
		t.Errorf("expected FOR UPDATE in view def, got: %s", def)
	}
}

func TestRuleutils4_ForShare(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR SHARE;`)
	if !strings.Contains(strings.ToUpper(def), "FOR SHARE") {
		t.Errorf("expected FOR SHARE in view def, got: %s", def)
	}
}

func TestRuleutils4_ForUpdateNowait(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR UPDATE NOWAIT;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "FOR UPDATE") {
		t.Errorf("expected FOR UPDATE in view def, got: %s", def)
	}
	if !strings.Contains(upper, "NOWAIT") {
		t.Errorf("expected NOWAIT in view def, got: %s", def)
	}
}

func TestRuleutils4_ForUpdateSkipLocked(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR UPDATE SKIP LOCKED;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "FOR UPDATE") {
		t.Errorf("expected FOR UPDATE in view def, got: %s", def)
	}
	if !strings.Contains(upper, "SKIP LOCKED") {
		t.Errorf("expected SKIP LOCKED in view def, got: %s", def)
	}
}

func TestRuleutils4_ForUpdateOfTable(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR UPDATE OF jtest;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "FOR UPDATE") {
		t.Errorf("expected FOR UPDATE in view def, got: %s", def)
	}
	if !strings.Contains(upper, "OF") {
		t.Errorf("expected OF in view def, got: %s", def)
	}
}

func TestRuleutils4_ForKeyShare(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT id FROM jtest FOR KEY SHARE;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "FOR KEY SHARE") {
		t.Errorf("expected FOR KEY SHARE in view def, got: %s", def)
	}
}

// --- 4.1 DISTINCT ON ---

func TestRuleutils4_DistinctOn(t *testing.T) {
	c := setupJsonTest(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT DISTINCT ON (name) name, val FROM jtest ORDER BY name, val;`)
	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "DISTINCT ON") {
		t.Errorf("expected DISTINCT ON in view def, got: %s", def)
	}
}

// --- 4.2 Plan-Internal Stubs ---

func TestRuleutils4_SubPlanStub(t *testing.T) {
	c := New()
	expr := &SubPlanExpr{}
	result := c.DeparseExpr(expr, nil, true)
	if result != "/* SubPlan */" {
		t.Errorf("expected /* SubPlan */, got: %s", result)
	}
}

func TestRuleutils4_AlternativeSubPlanStub(t *testing.T) {
	c := New()
	expr := &AlternativeSubPlanExpr{}
	result := c.DeparseExpr(expr, nil, true)
	if result != "/* AlternativeSubPlan */" {
		t.Errorf("expected /* AlternativeSubPlan */, got: %s", result)
	}
}

func TestRuleutils4_MergeSupportFuncStub(t *testing.T) {
	c := New()
	expr := &MergeSupportFuncExpr{ResultType: INT4OID}
	result := c.DeparseExpr(expr, nil, true)
	if result != "merge_action()" {
		t.Errorf("expected merge_action(), got: %s", result)
	}
}

func TestRuleutils4_InferenceElemStub(t *testing.T) {
	c := New()
	expr := &InferenceElemExpr{
		Expr:       &ConstExpr{TypeOID: INT4OID, Value: "42"},
		ResultType: INT4OID,
	}
	result := c.DeparseExpr(expr, nil, true)
	if result != "42" {
		t.Errorf("expected 42, got: %s", result)
	}
}

func TestRuleutils4_PartitionBoundSpecList(t *testing.T) {
	c := New()
	expr := &PartitionBoundSpecExpr{
		Strategy: 'l',
		ListValues: []AnalyzedExpr{
			&ConstExpr{TypeOID: INT4OID, Value: "1"},
			&ConstExpr{TypeOID: INT4OID, Value: "2"},
		},
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "FOR VALUES IN") {
		t.Errorf("expected FOR VALUES IN, got: %s", result)
	}
}

func TestRuleutils4_PartitionBoundSpecRange(t *testing.T) {
	c := New()
	expr := &PartitionBoundSpecExpr{
		Strategy: 'r',
		LowerBound: []AnalyzedExpr{
			&ConstExpr{TypeOID: INT4OID, Value: "1"},
		},
		UpperBound: []AnalyzedExpr{
			&ConstExpr{TypeOID: INT4OID, Value: "100"},
		},
	}
	result := c.DeparseExpr(expr, nil, true)
	if !strings.Contains(result, "FOR VALUES FROM") {
		t.Errorf("expected FOR VALUES FROM, got: %s", result)
	}
	if !strings.Contains(result, "TO") {
		t.Errorf("expected TO, got: %s", result)
	}
}

// --- 4.3 Integration Tests ---

func TestRuleutils4_IntegrationComplexView(t *testing.T) {
	// Complex view with 5+ expression types: CASE, COALESCE, aggregate, operator, constant, boolean.
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE orders (
			id int PRIMARY KEY,
			status text,
			amount numeric,
			customer text,
			created_at timestamptz DEFAULT CURRENT_TIMESTAMP
		);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	def := viewDef(t, c, `
		CREATE VIEW order_summary AS
		SELECT
			id,
			COALESCE(status, 'unknown') AS effective_status,
			CASE WHEN amount > 100 THEN 'high' WHEN amount > 50 THEN 'medium' ELSE 'low' END AS tier,
			amount * 1.1 AS amount_with_tax,
			status IS NOT NULL AS has_status,
			NULLIF(customer, '') AS customer_or_null,
			GREATEST(amount, 0) AS nonneg_amount
		FROM orders;
	`)

	// Check multiple expression types are present
	checks := map[string]bool{
		"COALESCE":        false,
		"CASE":            false,
		"IS NOT NULL":     false,
		"NULLIF":          false,
		"GREATEST":        false,
	}
	upper := strings.ToUpper(def)
	for kw := range checks {
		if strings.Contains(upper, kw) {
			checks[kw] = true
		}
	}
	for kw, found := range checks {
		if !found {
			t.Errorf("expected %s in complex view def, got: %s", kw, def)
		}
	}
}

func TestRuleutils4_IntegrationCTEWindowCoalesce(t *testing.T) {
	// View with CTE + window function + COALESCE
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE events (id int PRIMARY KEY, category text, value int);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	def := viewDef(t, c, `
		CREATE VIEW event_analysis AS
		WITH ranked AS (
			SELECT id, category, value,
				row_number() OVER (PARTITION BY category ORDER BY value DESC) AS rn
			FROM events
		)
		SELECT id, category, COALESCE(category, 'none') AS safe_category, rn
		FROM ranked
		WHERE rn <= 10;
	`)

	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "WITH") {
		t.Errorf("expected WITH in view def, got: %s", def)
	}
	if !strings.Contains(upper, "ROW_NUMBER") {
		t.Errorf("expected ROW_NUMBER in view def, got: %s", def)
	}
	if !strings.Contains(upper, "COALESCE") {
		t.Errorf("expected COALESCE in view def, got: %s", def)
	}
	if !strings.Contains(upper, "PARTITION BY") {
		t.Errorf("expected PARTITION BY in view def, got: %s", def)
	}
}

func TestRuleutils4_IntegrationGroupByRollupGroupingCase(t *testing.T) {
	// View with GROUP BY ROLLUP + GROUPING() + CASE WHEN
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE sales_data (region text, product text, revenue int);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	def := viewDef(t, c, `
		CREATE VIEW sales_report AS
		SELECT
			CASE WHEN GROUPING(region) = 1 THEN 'All Regions' ELSE region END AS region_label,
			CASE WHEN GROUPING(product) = 1 THEN 'All Products' ELSE product END AS product_label,
			sum(revenue) AS total_revenue
		FROM sales_data
		GROUP BY ROLLUP(region, product);
	`)

	upper := strings.ToUpper(def)
	if !strings.Contains(upper, "ROLLUP") {
		t.Errorf("expected ROLLUP in view def, got: %s", def)
	}
	if !strings.Contains(upper, "GROUPING") {
		t.Errorf("expected GROUPING in view def, got: %s", def)
	}
	if !strings.Contains(upper, "CASE") {
		t.Errorf("expected CASE in view def, got: %s", def)
	}
}

func TestRuleutils4_IntegrationArraySubscript(t *testing.T) {
	// View with array subscript
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE arrtable (id int PRIMARY KEY, tags text[]);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	def := viewDef(t, c, `
		CREATE VIEW v AS SELECT id, tags[1] AS first_tag FROM arrtable;
	`)

	if !strings.Contains(def, "[1]") && !strings.Contains(def, "[") {
		t.Errorf("expected array subscript in view def, got: %s", def)
	}
}

func TestRuleutils4_RegressionExistingViews(t *testing.T) {
	// Ensure basic view types still work (regression).
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE r (a int, b text, c boolean);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	tests := []struct {
		name string
		sql  string
		want string
	}{
		{"simple_select", `CREATE VIEW v AS SELECT a, b FROM r;`, "SELECT"},
		{"where_clause", `CREATE VIEW v AS SELECT a FROM r WHERE c;`, "WHERE"},
		{"order_by", `CREATE VIEW v AS SELECT a FROM r ORDER BY a;`, "ORDER BY"},
		{"limit", `CREATE VIEW v AS SELECT a FROM r LIMIT 10;`, "LIMIT"},
		{"distinct", `CREATE VIEW v AS SELECT DISTINCT a FROM r;`, "DISTINCT"},
		{"group_by", `CREATE VIEW v AS SELECT a, count(*) FROM r GROUP BY a;`, "GROUP BY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c2 := New()
			s2 := parseStmts(t, `CREATE TABLE r (a int, b text, c boolean);`)
			for _, s := range s2 {
				if err := c2.ProcessUtility(s); err != nil {
					t.Fatalf("setup: %v", err)
				}
			}
			def := viewDef(t, c2, tt.sql)
			if !strings.Contains(strings.ToUpper(def), tt.want) {
				t.Errorf("expected %s in view def, got: %s", tt.want, def)
			}
		})
	}
}

// --- LockingClause direct deparse test ---

func TestRuleutils4_LockingClauseDeparse(t *testing.T) {
	c := New()
	ctx := &deparseCtx{
		catalog:      c,
		buf:          &strings.Builder{},
		prettyIndent: false,
		prettyParen:  true,
	}

	tests := []struct {
		name   string
		clause *LockingClauseQ
		want   string
	}{
		{
			"for_update",
			&LockingClauseQ{Strength: LockForUpdate},
			"FOR UPDATE",
		},
		{
			"for_share_nowait",
			&LockingClauseQ{Strength: LockForShare, WaitPolicy: LockWaitError},
			"FOR SHARE NOWAIT",
		},
		{
			"for_update_skip_locked",
			&LockingClauseQ{Strength: LockForUpdate, WaitPolicy: LockWaitSkip},
			"FOR UPDATE SKIP LOCKED",
		},
		{
			"for_update_of_table",
			&LockingClauseQ{Strength: LockForUpdate, Tables: []string{"t1"}},
			"FOR UPDATE OF t1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.buf.Reset()
			ctx.getLockingClause(tt.clause)
			got := ctx.buf.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected %q in %q", tt.want, got)
			}
		})
	}
}
