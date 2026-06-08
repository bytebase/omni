package catalog

import (
	"strings"
	"testing"
)

// =============================================================================
// Ruleutils Phase 2: GroupingFunc, GROUPING SETS, XmlExpr, TABLESAMPLE
// =============================================================================

// setupRuleutils2 creates a catalog with tables for grouping/xml tests.
func setupRuleutils2(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE sales (dept text, product text, amount int);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return c
}

// --- 2.1 Grouping and Window Enhancements ---

func TestRuleutils2_GroupByRollup(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT dept, sum(amount) FROM sales GROUP BY ROLLUP(dept);`)
	if !strings.Contains(def, "ROLLUP(") {
		t.Errorf("expected ROLLUP in view def, got: %s", def)
	}
}

func TestRuleutils2_GroupByCube(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT dept, product, sum(amount) FROM sales GROUP BY CUBE(dept, product);`)
	if !strings.Contains(def, "CUBE(") {
		t.Errorf("expected CUBE in view def, got: %s", def)
	}
}

func TestRuleutils2_GroupByGroupingSets(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT dept, product, sum(amount) FROM sales GROUP BY GROUPING SETS((dept), (product), ());`)
	if !strings.Contains(def, "GROUPING SETS(") {
		t.Errorf("expected GROUPING SETS in view def, got: %s", def)
	}
}

func TestRuleutils2_GroupingFunc(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT dept, GROUPING(dept), sum(amount) FROM sales GROUP BY ROLLUP(dept);`)
	if !strings.Contains(def, "GROUPING(") {
		t.Errorf("expected GROUPING() in view def, got: %s", def)
	}
}

func TestRuleutils2_GroupingFuncMultipleArgs(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT dept, product, GROUPING(dept, product), sum(amount) FROM sales GROUP BY CUBE(dept, product);`)
	if !strings.Contains(def, "GROUPING(") {
		t.Errorf("expected GROUPING() in view def, got: %s", def)
	}
}

// --- 2.2 XML Expressions ---

func TestRuleutils2_XmlElement(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlelement(name "row", dept) FROM sales;`)
	if !strings.Contains(strings.ToUpper(def), "XMLELEMENT") {
		t.Errorf("expected XMLELEMENT in view def, got: %s", def)
	}
}

func TestRuleutils2_XmlForest(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlforest(dept AS "department", product) FROM sales;`)
	if !strings.Contains(strings.ToUpper(def), "XMLFOREST") {
		t.Errorf("expected XMLFOREST in view def, got: %s", def)
	}
}

func TestRuleutils2_XmlConcat(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlconcat(xmlelement(name "a", dept), xmlelement(name "b", product)) FROM sales;`)
	if !strings.Contains(strings.ToUpper(def), "XMLCONCAT") {
		t.Errorf("expected XMLCONCAT in view def, got: %s", def)
	}
}

func TestRuleutils2_XmlParse(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlparse(document '<doc/>');`)
	if !strings.Contains(strings.ToUpper(def), "XMLPARSE") {
		t.Errorf("expected XMLPARSE in view def, got: %s", def)
	}
}

func TestRuleutils2_XmlPI(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlpi(name "php");`)
	if !strings.Contains(strings.ToUpper(def), "XMLPI") {
		t.Errorf("expected XMLPI in view def, got: %s", def)
	}
}

func TestRuleutils2_XmlSerialize(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT xmlserialize(content xmlelement(name "r", dept) AS text) FROM sales;`)
	if !strings.Contains(strings.ToUpper(def), "XMLSERIALIZE") {
		t.Errorf("expected XMLSERIALIZE in view def, got: %s", def)
	}
}

// --- 2.3 TABLESAMPLE ---

func TestRuleutils2_TablesampleBernoulli(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT * FROM sales TABLESAMPLE BERNOULLI(10);`)
	if !strings.Contains(strings.ToUpper(def), "TABLESAMPLE") {
		t.Errorf("expected TABLESAMPLE in view def, got: %s", def)
	}
	if !strings.Contains(strings.ToLower(def), "bernoulli") {
		t.Errorf("expected bernoulli in view def, got: %s", def)
	}
}

func TestRuleutils2_TablesampleSystem(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT * FROM sales TABLESAMPLE SYSTEM(50);`)
	if !strings.Contains(strings.ToUpper(def), "TABLESAMPLE") {
		t.Errorf("expected TABLESAMPLE in view def, got: %s", def)
	}
	if !strings.Contains(strings.ToLower(def), "system") {
		t.Errorf("expected system in view def, got: %s", def)
	}
}

func TestRuleutils2_TablesampleRepeatable(t *testing.T) {
	c := setupRuleutils2(t)
	def := viewDef(t, c, `CREATE VIEW v AS SELECT * FROM sales TABLESAMPLE BERNOULLI(10) REPEATABLE (42);`)
	if !strings.Contains(strings.ToUpper(def), "REPEATABLE") {
		t.Errorf("expected REPEATABLE in view def, got: %s", def)
	}
}
