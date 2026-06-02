package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func rawStmt(t *testing.T, sql string) ast.StmtNode {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	return raw.Stmt
}

func TestBytebaseCreateTriggerReferencingOldNew(t *testing.T) {
	stmt, ok := rawStmt(t, `CREATE OR REPLACE TRIGGER audit_emp
AFTER UPDATE ON employees
REFERENCING OLD AS old_row NEW AS new_row
FOR EACH ROW
BEGIN
  NULL;
END;`).(*ast.CreateTriggerStmt)
	if !ok {
		t.Fatalf("expected CreateTriggerStmt")
	}

	if stmt.Referencing == nil {
		t.Fatal("expected REFERENCING clause")
	}
	if stmt.Referencing.OldAlias != "OLD_ROW" {
		t.Fatalf("old alias = %q, want OLD_ROW", stmt.Referencing.OldAlias)
	}
	if stmt.Referencing.NewAlias != "NEW_ROW" {
		t.Fatalf("new alias = %q, want NEW_ROW", stmt.Referencing.NewAlias)
	}
	if !stmt.ForEachRow {
		t.Fatal("expected FOR EACH ROW")
	}
}

func TestBytebasePivotTypedOutputShape(t *testing.T) {
	stmt, ok := rawStmt(t, `SELECT * FROM sales PIVOT (SUM(amount) AS total, COUNT(*) AS cnt FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))`).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}

	if stmt.Pivot == nil {
		t.Fatal("expected PIVOT")
	}
	if stmt.Pivot.Aggregates == nil || stmt.Pivot.Aggregates.Len() != 2 {
		t.Fatalf("expected 2 typed aggregates, got %v", stmt.Pivot.Aggregates)
	}
	agg := stmt.Pivot.Aggregates.Items[0].(*ast.PivotAggregate)
	if agg.Alias != "TOTAL" {
		t.Fatalf("aggregate alias = %q, want TOTAL", agg.Alias)
	}
	if stmt.Pivot.InItems == nil || stmt.Pivot.InItems.Len() != 2 {
		t.Fatalf("expected 2 pivot in-items, got %v", stmt.Pivot.InItems)
	}
	in := stmt.Pivot.InItems.Items[0].(*ast.PivotInItem)
	if in.Alias != "Q1" || in.Values == nil || in.Values.Len() != 1 {
		t.Fatalf("unexpected pivot in-item: %+v", in)
	}
	if stmt.Pivot.Source == nil {
		t.Fatal("expected PIVOT source table expression")
	}
}

func TestBytebaseUnpivotTypedMappings(t *testing.T) {
	stmt, ok := rawStmt(t, `SELECT * FROM quarterly_sales UNPIVOT INCLUDE NULLS ((sales, cost) FOR quarter IN ((q1_sales, q1_cost) AS 'Q1', (q2_sales, q2_cost) AS 'Q2'))`).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}

	if stmt.Unpivot == nil {
		t.Fatal("expected UNPIVOT")
	}
	if stmt.Unpivot.ValueColumns == nil || stmt.Unpivot.ValueColumns.Len() != 2 {
		t.Fatalf("expected 2 value columns, got %v", stmt.Unpivot.ValueColumns)
	}
	if stmt.Unpivot.InputMappings == nil || stmt.Unpivot.InputMappings.Len() != 2 {
		t.Fatalf("expected 2 input mappings, got %v", stmt.Unpivot.InputMappings)
	}
	mapping := stmt.Unpivot.InputMappings.Items[0].(*ast.UnpivotInItem)
	if mapping.InputColumns == nil || mapping.InputColumns.Len() != 2 {
		t.Fatalf("expected 2 input columns, got %v", mapping.InputColumns)
	}
	if mapping.Alias != "Q1" {
		t.Fatalf("mapping alias = %q, want Q1", mapping.Alias)
	}
	if stmt.Unpivot.Source == nil {
		t.Fatal("expected UNPIVOT source table expression")
	}
}

func TestBytebaseModelRuleCellReference(t *testing.T) {
	stmt, ok := rawStmt(t, `SELECT * FROM sales MODEL DIMENSION BY (product AS p) MEASURES (amount AS amt) RULES (amt[p] = amt[p])`).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}

	if stmt.ModelClause == nil || stmt.ModelClause.MainModel == nil {
		t.Fatal("expected MODEL main model")
	}
	measures := stmt.ModelClause.MainModel.ColumnClauses.Measures
	if measures == nil || measures.Len() != 1 {
		t.Fatalf("expected one measure, got %v", measures)
	}
	rules := stmt.ModelClause.MainModel.RulesClause.Rules
	if rules == nil || rules.Len() != 1 {
		t.Fatalf("expected one rule, got %v", rules)
	}
	rule := rules.Items[0].(*ast.ModelRule)
	if rule.Cell == nil {
		t.Fatal("expected typed model cell reference")
	}
	if rule.Cell.Measure != "AMT" || rule.Cell.Dimensions == nil || rule.Cell.Dimensions.Len() != 1 {
		t.Fatalf("unexpected model cell: %+v", rule.Cell)
	}
}

func TestBytebaseTableCollectionTypedCallAndColumnAliases(t *testing.T) {
	stmt, ok := rawStmt(t, `SELECT * FROM TABLE(pkg.get_employees(10)) e (employee_id, employee_name)`).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}
	ref := stmt.FromClause.Items[0].(*ast.TableCollectionExpr)

	if ref.FunctionCall == nil {
		t.Fatal("expected TABLE() function call")
	}
	if ref.Alias == nil || ref.Alias.Name != "E" {
		t.Fatalf("alias = %+v, want E", ref.Alias)
	}
	if ref.ColumnAliases == nil || ref.ColumnAliases.Len() != 2 {
		t.Fatalf("expected 2 column aliases, got %v", ref.ColumnAliases)
	}
}

func TestBytebaseMatchRecognizeSourceAndMeasures(t *testing.T) {
	stmt, ok := rawStmt(t, `SELECT * FROM trades MATCH_RECOGNIZE (
  PARTITION BY account_id
  ORDER BY trade_time
  MEASURES FIRST(price) AS first_price, LAST(price) AS last_price
  ONE ROW PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.price > A.price
) mr`).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}
	ref := stmt.FromClause.Items[0].(*ast.MatchRecognizeClause)

	if ref.Source == nil {
		t.Fatal("expected MATCH_RECOGNIZE source table expression")
	}
	if ref.Measures == nil || ref.Measures.Len() != 2 {
		t.Fatalf("expected 2 measures, got %v", ref.Measures)
	}
	if ref.Alias == nil || ref.Alias.Name != "MR" {
		t.Fatalf("alias = %+v, want MR", ref.Alias)
	}
}

func TestBytebaseCreateIndexTypedAdvancedMetadata(t *testing.T) {
	stmt, ok := rawStmt(t, `CREATE INDEX ix_sales ON sales (sale_date) LOCAL STORE IN (ts1, ts2) STORAGE (INITIAL 64K) ANNOTATIONS (classification 'pii')`).(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected CreateIndexStmt")
	}

	if stmt.LocalPartition == "" {
		t.Fatal("expected typed LOCAL partition payload")
	}
	if stmt.StorageSpec == "" {
		t.Fatal("expected typed STORAGE payload")
	}
	if stmt.AnnotationsSpec == "" {
		t.Fatal("expected typed ANNOTATIONS payload")
	}
}
