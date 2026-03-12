package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// =============================================================================
// Phase 4: Analyze + Deparse Extension Tests
// =============================================================================

// setupPhase4 creates a catalog with common tables for testing.
func setupPhase4(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE t1 (id int, name text, val int);
		CREATE TABLE t2 (id int, label text, amount numeric);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return c
}

// viewDef creates a view and returns its definition.
func viewDef(t *testing.T, c *Catalog, sql string) string {
	t.Helper()
	stmts := parseStmts(t, sql)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	// Extract view name from SQL.
	list, err := pgparser.Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, item := range list.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			item = raw.Stmt
		}
		if vs, ok := item.(*nodes.ViewStmt); ok {
			def, err := c.GetViewDefinition("", vs.View.Relname)
			if err != nil {
				t.Fatalf("get view def: %v", err)
			}
			return def
		}
	}
	t.Fatal("no view statement found")
	return ""
}

// --- 4f. Additional expression types ---

func TestPhase4_ArrayExpr(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_array AS SELECT ARRAY[1, 2, 3] AS arr")
	if !strings.Contains(def, "ARRAY[") {
		t.Errorf("expected ARRAY constructor in view def, got: %s", def)
	}
}

func TestPhase4_RowExpr(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_row AS SELECT ROW(1, 'hello') AS r")
	if !strings.Contains(def, "ROW(") {
		t.Errorf("expected ROW constructor in view def, got: %s", def)
	}
}

func TestPhase4_CollateExpr(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, `CREATE VIEW v_collate AS SELECT name COLLATE "C" AS cname FROM t1`)
	if !strings.Contains(def, "COLLATE") {
		t.Errorf("expected COLLATE in view def, got: %s", def)
	}
}

func TestPhase4_SQLValueFunc(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_svf AS SELECT CURRENT_DATE AS d, CURRENT_USER AS u")
	if !strings.Contains(def, "CURRENT_DATE") {
		t.Errorf("expected CURRENT_DATE in view def, got: %s", def)
	}
	if !strings.Contains(def, "CURRENT_USER") {
		t.Errorf("expected CURRENT_USER in view def, got: %s", def)
	}
}

// --- 4d. DISTINCT ON ---

func TestPhase4_DistinctOn(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_distinct_on AS SELECT DISTINCT ON (id) id, name FROM t1 ORDER BY id, name")
	if !strings.Contains(def, "DISTINCT ON (") {
		t.Errorf("expected DISTINCT ON in view def, got: %s", def)
	}
}

func TestPhase4_DistinctPlain(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_distinct AS SELECT DISTINCT name FROM t1")
	// Should have DISTINCT but NOT "DISTINCT ON"
	if !strings.Contains(def, "DISTINCT") {
		t.Errorf("expected DISTINCT in view def, got: %s", def)
	}
	if strings.Contains(def, "DISTINCT ON") {
		t.Errorf("should be plain DISTINCT, not DISTINCT ON, got: %s", def)
	}
}

// --- 4e. LATERAL ---

func TestPhase4_LateralSubquery(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_lateral AS SELECT t1.id, sub.total FROM t1, LATERAL (SELECT sum(val) AS total FROM t1 t1b WHERE t1b.id = t1.id) sub")
	if !strings.Contains(def, "LATERAL") {
		t.Errorf("expected LATERAL in view def, got: %s", def)
	}
}

// --- 4a. CTE Support ---

func TestPhase4_SimpleCTE(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, `
		CREATE VIEW v_cte AS
		WITH cte AS (SELECT id, name FROM t1 WHERE val > 0)
		SELECT id, name FROM cte
	`)
	if !strings.Contains(def, "WITH") {
		t.Errorf("expected WITH in view def, got: %s", def)
	}
	if !strings.Contains(def, "cte") {
		t.Errorf("expected CTE name in view def, got: %s", def)
	}
}

func TestPhase4_CTEWithAliases(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, `
		CREATE VIEW v_cte_alias AS
		WITH cte(a, b) AS (SELECT id, name FROM t1)
		SELECT a, b FROM cte
	`)
	if !strings.Contains(def, "WITH") {
		t.Errorf("expected WITH in view def, got: %s", def)
	}
}

func TestPhase4_MultipleCTEs(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, `
		CREATE VIEW v_multi_cte AS
		WITH
			cte1 AS (SELECT id, name FROM t1),
			cte2 AS (SELECT id, label FROM t2)
		SELECT cte1.id, cte1.name, cte2.label
		FROM cte1 JOIN cte2 ON cte1.id = cte2.id
	`)
	if !strings.Contains(def, "WITH") {
		t.Errorf("expected WITH in view def, got: %s", def)
	}
	if !strings.Contains(def, "cte1") || !strings.Contains(def, "cte2") {
		t.Errorf("expected both CTE names in view def, got: %s", def)
	}
}

// --- 4b. Window Functions ---

func TestPhase4_WindowFuncRowNumber(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_rownum AS SELECT id, name, row_number() OVER (ORDER BY id) AS rn FROM t1")
	if !strings.Contains(def, "OVER") {
		t.Errorf("expected OVER in view def, got: %s", def)
	}
	if !strings.Contains(def, "row_number") {
		t.Errorf("expected row_number in view def, got: %s", def)
	}
}

func TestPhase4_WindowFuncPartitionBy(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_wf_part AS SELECT id, name, count(*) OVER (PARTITION BY name ORDER BY id) AS cnt FROM t1")
	if !strings.Contains(def, "PARTITION BY") {
		t.Errorf("expected PARTITION BY in view def, got: %s", def)
	}
}

// --- Analyze correctness ---

func TestPhase4_CTEColumnTypes(t *testing.T) {
	c := setupPhase4(t)
	stmts := parseStmts(t, `
		CREATE VIEW v_cte_types AS
		WITH cte AS (SELECT id, name FROM t1)
		SELECT id, name FROM cte
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_cte_types")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 2 {
		t.Fatalf("expected at least 2 columns, got %d", len(rel.Columns))
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("expected column 0 type INT4OID, got %d", rel.Columns[0].TypeOID)
	}
	if rel.Columns[1].TypeOID != TEXTOID {
		t.Errorf("expected column 1 type TEXTOID, got %d", rel.Columns[1].TypeOID)
	}
}

func TestPhase4_ArrayExprType(t *testing.T) {
	c := setupPhase4(t)
	stmts := parseStmts(t, "CREATE VIEW v_arr_type AS SELECT ARRAY[1, 2, 3] AS arr")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_arr_type")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 1 {
		t.Fatal("expected at least 1 column")
	}
	// Array of INT4 should be INT4ARRAYOID (1007).
	arrType := c.findArrayType(INT4OID)
	if arrType == 0 {
		t.Skip("no array type for INT4OID in catalog")
	}
	if rel.Columns[0].TypeOID != arrType {
		t.Errorf("expected column type %d (int4[]), got %d", arrType, rel.Columns[0].TypeOID)
	}
}

func TestPhase4_DistinctOnAnalysis(t *testing.T) {
	c := setupPhase4(t)
	stmts := parseStmts(t, "CREATE VIEW v_don_test AS SELECT DISTINCT ON (id) id, name FROM t1 ORDER BY id, name")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_don_test")
	if rel == nil {
		t.Fatal("view not found")
	}
	if rel.AnalyzedQuery == nil {
		t.Fatal("no analyzed query")
	}
	if !rel.AnalyzedQuery.Distinct {
		t.Error("expected Distinct=true")
	}
	if len(rel.AnalyzedQuery.DistinctOn) == 0 {
		t.Error("expected DistinctOn to be populated")
	}
}

func TestPhase4_LateralAnalysis(t *testing.T) {
	c := setupPhase4(t)
	sql := "CREATE VIEW v_lat_test AS SELECT t1.id, sub.total FROM t1, LATERAL (SELECT sum(val) AS total FROM t1 t1b WHERE t1b.id = t1.id) sub"
	stmts := parseStmts(t, sql)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_lat_test")
	if rel == nil {
		t.Fatal("view not found")
	}
	if rel.AnalyzedQuery == nil {
		t.Fatal("no analyzed query")
	}

	// Check that one RTE is marked LATERAL.
	foundLateral := false
	for _, rte := range rel.AnalyzedQuery.RangeTable {
		if rte.Lateral {
			foundLateral = true
			break
		}
	}
	if !foundLateral {
		t.Error("expected one RTE with Lateral=true")
	}
}

func TestPhase4_WindowClauseAnalysis(t *testing.T) {
	c := setupPhase4(t)
	stmts := parseStmts(t, "CREATE VIEW v_wc_test AS SELECT id, name, row_number() OVER (ORDER BY id) AS rn FROM t1")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_wc_test")
	if rel == nil {
		t.Fatal("view not found")
	}
	if rel.AnalyzedQuery == nil {
		t.Fatal("no analyzed query")
	}
	if len(rel.AnalyzedQuery.WindowClause) == 0 {
		t.Error("expected WindowClause to be populated")
	}
}
