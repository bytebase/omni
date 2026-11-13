package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereDisallowFunctionsAndCalculations(t *testing.T) {
	table := func(name string, indexed ...string) *metadata.TableMetadata {
		return &metadata.TableMetadata{
			Name:    name,
			Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "total"}, {Name: "note"}},
			Indexes: []*metadata.IndexMetadata{{Name: "ix", Expressions: indexed}},
		}
	}
	schema := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{
		{Name: "zeta", Tables: []*metadata.TableMetadata{table("orders", "note")}},
		{Name: "dbo", Tables: []*metadata.TableMetadata{table("orders", "id")}},
		{Name: "alpha", Tables: []*metadata.TableMetadata{table("orders", "total"), table("logs", "id")}},
	}}}
	fn := func(col, name, text string) statementFinding {
		return statementFinding{Code: v1review.FunctionOnIndexedColumn, Name: col, Value: name, Text: text}
	}
	calc := func(col, text string) statementFinding {
		return statementFinding{Code: v1review.CalculationOnIndexedColumn, Name: col, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		want   []statementFinding
	}{
		{name: "unqualified resolves in dbo", sql: "SELECT * FROM orders WHERE abs(id) = 1 AND total + 1 = 2", schema: schema, want: []statementFinding{fn("id", "ABS", "abs(id)")}},
		{name: "named schema", sql: "SELECT * FROM alpha.orders WHERE total * 2 = 1", schema: schema, want: []statementFinding{calc("total", "total * 2")}},
		{name: "fallback in sorted schema order", sql: "SELECT * FROM logs WHERE -id = 1", schema: schema, want: []statementFinding{calc("id", "-id")}},
		{name: "qualified function name printed", sql: "SELECT * FROM orders WHERE dbo.MyFn(id) = 1", schema: schema, want: []statementFinding{fn("id", "DBO.MYFN", "dbo.MyFn(id)")}},
		{name: "other database not resolved", sql: "SELECT * FROM other.dbo.orders WHERE abs(id) = 1", schema: schema},
		{name: "no snapshot", sql: "SELECT * FROM orders WHERE abs(id) = 1"},
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
