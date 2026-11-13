package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestIndexNotRedundant(t *testing.T) {
	snapshot := &metadata.DatabaseSchemaMetadata{
		Name: "master",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{Name: "t", Indexes: []*metadata.IndexMetadata{
						{Name: "idx_ab", Expressions: []string{"a", "b"}},
						{Name: "idx_empty"},
						{Name: "idx_a", Expressions: []string{"a"}},
					}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{{Name: "first", Expressions: []string{"a"}}}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{{Name: "second", Expressions: []string{"b"}}}},
					{Name: "parted", Partitions: []*metadata.TablePartitionMetadata{{Name: "p1"}}, Indexes: []*metadata.IndexMetadata{{Name: "idx_p", Expressions: []string{"a"}}}},
					{Name: "", Indexes: []*metadata.IndexMetadata{{Name: "idx_nameless", Expressions: []string{"a"}}}},
				},
			},
			{Name: "Sales", Tables: []*metadata.TableMetadata{{Name: "Orders", Indexes: []*metadata.IndexMetadata{{Name: "idx_c", Expressions: []string{"c"}}}}}},
			{Name: "", Tables: []*metadata.TableMetadata{{Name: "t", Indexes: []*metadata.IndexMetadata{{Name: "idx_x", Expressions: []string{"x"}}}}}},
		},
	}
	type finding struct {
		Name, Other, Schema, Table, Text string
		Statement                        int
	}
	redundant := func(name, other, schema, table, text string) finding {
		return finding{name, other, schema, table, text, 0}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		noMeta bool
		want   []finding
	}{
		{
			name: "first column matches an existing index's first expression",
			sql:  "CREATE INDEX idx_new ON t (a, c)",
			want: []finding{redundant("idx_new", "idx_ab", "dbo", "t", "CREATE INDEX idx_new ON t (a, c)")},
		},
		{name: "a different first column", sql: "CREATE INDEX idx_new ON t (b, a)"},
		{
			name: "schema and current database qualifiers",
			sql:  "CREATE UNIQUE NONCLUSTERED INDEX idx_new ON MASTER.Sales.Orders (c DESC) INCLUDE (d) WHERE c > 0",
			want: []finding{redundant("idx_new", "idx_c", "Sales", "Orders", "CREATE UNIQUE NONCLUSTERED INDEX idx_new ON MASTER.Sales.Orders (c DESC) INCLUDE (d) WHERE c > 0")},
		},
		{name: "another database is not judged", sql: "CREATE INDEX idx_new ON other.dbo.t (a)"},
		{name: "schema names compare exactly", sql: "CREATE INDEX idx_new ON sales.Orders (c)"},
		{name: "table names compare exactly", sql: "CREATE INDEX idx_new ON T (a)"},
		{name: "column names compare exactly", sql: "CREATE INDEX idx_new ON t (A)"},
		{
			name: "the last table of a name is used",
			sql:  "CREATE INDEX idx_new ON dup (b)",
			want: []finding{redundant("idx_new", "second", "dbo", "dup", "CREATE INDEX idx_new ON dup (b)")},
		},
		{name: "the earlier table of a name is not", sql: "CREATE INDEX idx_new ON dup (a)"},
		{name: "a partition does not find its table", sql: "CREATE INDEX idx_new ON p1 (a)"},
		{name: "an empty-named schema is skipped", sql: "CREATE INDEX idx_new ON t (x)"},
		{name: "an unknown table", sql: "CREATE INDEX idx_new ON u (a)"},
		{name: "nil metadata finds nothing", sql: "CREATE INDEX idx_new ON t (a)", noMeta: true},
		{
			name: "one finding per statement, in statement order",
			sql:  "CREATE INDEX i1 ON t (a);\nSELECT 1;\nCREATE INDEX i2 ON dbo.dup (b)",
			want: []finding{
				redundant("i1", "idx_ab", "dbo", "t", "CREATE INDEX i1 ON t (a)"),
				{"i2", "second", "dbo", "dup", "CREATE INDEX i2 ON dbo.dup (b)", 2},
			},
		},
		{name: "a routine body is not entered", sql: "CREATE PROCEDURE p AS BEGIN CREATE INDEX idx_new ON t (a) END"},
		{name: "a block is not entered", sql: "BEGIN CREATE INDEX idx_new ON t (a) END"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			meta := snapshot
			if tt.noMeta {
				meta = nil
			}
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []finding
			for _, f := range IndexNotRedundant(list.Items, v1review.Schema{Metadata: meta}, "master") {
				if f.Code != v1review.IndexRedundant {
					t.Errorf("code %q, want %q", f.Code, v1review.IndexRedundant)
				}
				got = append(got, finding{f.Name, f.Other, f.Schema, f.Table, tt.sql[f.Start:f.End], f.Statement})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIndexNotRedundantReadsTheWholeSnapshot(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a nil table in a named schema did not panic")
		}
	}()
	IndexNotRedundant(nil, v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{{Name: "dbo", Tables: []*metadata.TableMetadata{nil}}},
	}}, "master")
}
