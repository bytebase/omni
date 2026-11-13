package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereDisallowFunctionsAndCalculations(t *testing.T) {
	db := &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{
				Name:       "orders",
				Columns:    []*metadata.ColumnMetadata{{Name: "id"}, {Name: "total"}, {Name: "note"}},
				Indexes:    []*metadata.IndexMetadata{{Name: "PRIMARY", Expressions: []string{"id"}}, {Name: "idx_total", Expressions: []string{"TOTAL"}}},
				Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}},
			},
			{Name: "items", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "qty"}}},
		},
	}}}
	schema := v1review.Schema{Metadata: db}
	fn := func(col, name, text string) stmtFinding {
		return stmtFinding{Code: v1review.FunctionOnIndexedColumn, Name: col, Detail: name, Text: text}
	}
	calc := func(col, text string) stmtFinding {
		return stmtFinding{Code: v1review.CalculationOnIndexedColumn, Name: col, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		want   []stmtFinding
	}{
		{name: "function upper-cased", sql: "SELECT * FROM orders WHERE abs(id) = 1;", schema: schema, want: []stmtFinding{fn("id", "ABS", "abs(id)")}},
		{name: "index expression lower-cased", sql: "SELECT * FROM orders WHERE total * 2 > 1;", schema: schema, want: []stmtFinding{calc("total", "total * 2")}},
		{name: "right operand", sql: "SELECT * FROM orders WHERE 1 + id = 2;", schema: schema, want: []stmtFinding{calc("id", "1 + id")}},
		{name: "unary minus", sql: "SELECT * FROM orders WHERE -id = 1;", schema: schema, want: []stmtFinding{calc("id", "-id")}},
		{name: "unindexed column", sql: "SELECT * FROM orders WHERE upper(note) = 'a';", schema: schema},
		{name: "no snapshot", sql: "SELECT * FROM orders WHERE abs(id) = 1;"},
		{name: "partition finds its table", sql: "SELECT * FROM p0 WHERE abs(id) = 1;", schema: schema, want: []stmtFinding{fn("id", "ABS", "abs(id)")}},
		{name: "other database not resolved", sql: "SELECT * FROM other.orders WHERE abs(id) = 1;", schema: schema},
		{name: "own database resolved", sql: "SELECT * FROM SHOP.orders WHERE abs(id) = 1;", schema: schema, want: []stmtFinding{fn("id", "ABS", "abs(id)")}},
		{
			name:   "correlated subquery falls back to the outer table",
			sql:    "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.qty = abs(o.id));",
			schema: schema,
			want:   []stmtFinding{fn("id", "ABS", "abs(o.id)")},
		},
		{
			name:   "unknown local source blocks the outer fallback",
			sql:    "SELECT * FROM orders WHERE EXISTS (SELECT 1 FROM missing WHERE abs(id) = 1);",
			schema: schema,
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN SELECT * FROM orders WHERE abs(id) = 1; END;", schema: schema},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runStatement(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return StatementWhereDisallowFunctionsAndCalculations(stmts, tt.schema)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
