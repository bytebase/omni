package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/cosmosdb/analysis"
	"github.com/bytebase/omni/cosmosdb/ast"
)

func parseSelect(t *testing.T, sql string) *ast.SelectStmt {
	t.Helper()
	stmts, err := cosmosdb.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, len(stmts))
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): AST is %T, want *ast.SelectStmt", sql, stmts[0].AST)
	}
	return sel
}

func pathStr(fp analysis.FieldPath) string {
	return fp.String()
}

func TestAnalyzeSelectStar(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM c")
	qa := analysis.Analyze(sel)
	if !qa.SelectStar {
		t.Error("SelectStar: got false, want true")
	}
	if len(qa.Projections) != 0 {
		t.Errorf("Projections: got %d, want 0", len(qa.Projections))
	}
}

func TestAnalyzeSimpleProjection(t *testing.T) {
	sel := parseSelect(t, "SELECT c.name FROM c")
	qa := analysis.Analyze(sel)
	if qa.SelectStar {
		t.Error("SelectStar: got true, want false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "name" {
		t.Errorf("Name: got %q, want %q", p.Name, "name")
	}
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.name")
	}
}

func TestAnalyzeAliasedProjection(t *testing.T) {
	sel := parseSelect(t, "SELECT c.name AS n FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if qa.Projections[0].Name != "n" {
		t.Errorf("Name: got %q, want %q", qa.Projections[0].Name, "n")
	}
}

func TestAnalyzeAliasResolution(t *testing.T) {
	sel := parseSelect(t, "SELECT t.name FROM container t")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "container.name" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "container.name")
	}
}

func TestAnalyzeDotChain(t *testing.T) {
	sel := parseSelect(t, "SELECT c.address.city FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "city" {
		t.Errorf("Name: got %q, want %q", p.Name, "city")
	}
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.address.city" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.address.city")
	}
}

func TestAnalyzeBracketAccess(t *testing.T) {
	sel := parseSelect(t, "SELECT c.addresses[1].country FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.addresses[1].country" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.addresses[1].country")
	}
}

func TestAnalyzeBracketStringAccess(t *testing.T) {
	sel := parseSelect(t, `SELECT c["field"] FROM c`)
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.field" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.field")
	}
}

func TestAnalyzeExpressionMultiplePaths(t *testing.T) {
	sel := parseSelect(t, "SELECT c.name ?? c.nickname AS displayName FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "displayName" {
		t.Errorf("Name: got %q, want %q", p.Name, "displayName")
	}
	if len(p.SourcePaths) != 2 {
		t.Fatalf("SourcePaths: got %d, want 2", len(p.SourcePaths))
	}
}

func TestAnalyzeFunctionCall(t *testing.T) {
	sel := parseSelect(t, "SELECT CONCAT(c.first, c.last) AS fullName FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "fullName" {
		t.Errorf("Name: got %q, want %q", p.Name, "fullName")
	}
	if len(p.SourcePaths) != 2 {
		t.Fatalf("SourcePaths: got %d, want 2", len(p.SourcePaths))
	}
}

func TestAnalyzeUDFCall(t *testing.T) {
	sel := parseSelect(t, "SELECT udf.myFunc(c.name) AS result FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "result" {
		t.Errorf("Name: got %q, want %q", p.Name, "result")
	}
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.name")
	}
}

func TestAnalyzeWherePredicate(t *testing.T) {
	sel := parseSelect(t, `SELECT * FROM c WHERE c.country = "US"`)
	qa := analysis.Analyze(sel)
	if !qa.SelectStar {
		t.Error("SelectStar: got false, want true")
	}
	if len(qa.Predicates) != 1 {
		t.Fatalf("Predicates: got %d, want 1", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.country" {
		t.Errorf("Predicates[0]: got %q, want %q", got, "c.country")
	}
}

func TestAnalyzeWhereMultiplePredicates(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM c WHERE c.age > 18 AND c.active = true")
	qa := analysis.Analyze(sel)
	if len(qa.Predicates) != 2 {
		t.Fatalf("Predicates: got %d, want 2", len(qa.Predicates))
	}
	paths := map[string]bool{}
	for _, p := range qa.Predicates {
		paths[pathStr(p)] = true
	}
	if !paths["c.age"] {
		t.Error("missing predicate c.age")
	}
	if !paths["c.active"] {
		t.Error("missing predicate c.active")
	}
}

func TestAnalyzeWhereBetween(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM c WHERE c.population BETWEEN 100000 AND 5000000")
	qa := analysis.Analyze(sel)
	if len(qa.Predicates) != 1 {
		t.Fatalf("Predicates: got %d, want 1", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.population" {
		t.Errorf("Predicates[0]: got %q, want %q", got, "c.population")
	}
}

func TestAnalyzeWhereIn(t *testing.T) {
	sel := parseSelect(t, `SELECT * FROM c WHERE c.country IN ("US", "UK")`)
	qa := analysis.Analyze(sel)
	if len(qa.Predicates) != 1 {
		t.Fatalf("Predicates: got %d, want 1", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.country" {
		t.Errorf("Predicates[0]: got %q, want %q", got, "c.country")
	}
}

func TestAnalyzeJoin(t *testing.T) {
	sel := parseSelect(t, "SELECT p.name, t.tag FROM products p JOIN t IN p.tags")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 2 {
		t.Fatalf("Projections: got %d, want 2", len(qa.Projections))
	}
	// p.name -> products.name
	if len(qa.Projections[0].SourcePaths) != 1 {
		t.Fatalf("Projections[0].SourcePaths: got %d, want 1", len(qa.Projections[0].SourcePaths))
	}
	if got := pathStr(qa.Projections[0].SourcePaths[0]); got != "products.name" {
		t.Errorf("Projections[0] path: got %q, want %q", got, "products.name")
	}
	// t.tag -> products.tags.tag
	if len(qa.Projections[1].SourcePaths) != 1 {
		t.Fatalf("Projections[1].SourcePaths: got %d, want 1", len(qa.Projections[1].SourcePaths))
	}
	if got := pathStr(qa.Projections[1].SourcePaths[0]); got != "products.tags.tag" {
		t.Errorf("Projections[1] path: got %q, want %q", got, "products.tags.tag")
	}
}

func TestAnalyzeJoinWherePredicate(t *testing.T) {
	sel := parseSelect(t, "SELECT p.name FROM products p JOIN t IN p.tags WHERE t.active = true")
	qa := analysis.Analyze(sel)
	if len(qa.Predicates) != 1 {
		t.Fatalf("Predicates: got %d, want 1", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "products.tags.active" {
		t.Errorf("Predicates[0]: got %q, want %q", got, "products.tags.active")
	}
}

func TestAnalyzeNestedJoin(t *testing.T) {
	sel := parseSelect(t, "SELECT s.size FROM products p JOIN t IN p.tags JOIN s IN t.sizes")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "products.tags.sizes.size" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "products.tags.sizes.size")
	}
}

func TestAnalyzeSelectValue(t *testing.T) {
	sel := parseSelect(t, "SELECT VALUE c.name FROM c")
	qa := analysis.Analyze(sel)
	if qa.SelectStar {
		t.Error("SelectStar: got true, want false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("SourcePaths: got %d, want 1", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("SourcePaths[0]: got %q, want %q", got, "c.name")
	}
}

func TestAnalyzeNoFrom(t *testing.T) {
	sel := parseSelect(t, "SELECT 1")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if len(qa.Projections[0].SourcePaths) != 0 {
		t.Errorf("SourcePaths: got %d, want 0", len(qa.Projections[0].SourcePaths))
	}
}

func TestAnalyzeValueCount(t *testing.T) {
	sel := parseSelect(t, "SELECT VALUE COUNT(1) FROM c")
	qa := analysis.Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if len(qa.Projections[0].SourcePaths) != 0 {
		t.Errorf("SourcePaths: got %d, want 0", len(qa.Projections[0].SourcePaths))
	}
}
