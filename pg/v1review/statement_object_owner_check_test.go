package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

func TestStatementObjectOwnerCheck(t *testing.T) {
	schema := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Name:  "shop",
		Owner: "alice",
		Schemas: []*metadata.SchemaMetadata{
			{Name: "public", Owner: "pg_database_owner", Tables: []*metadata.TableMetadata{{Name: "orders", Owner: "bob"}}},
			{Name: "sales", Owner: "bob"},
		},
	}}
	owner := func(kind, name, owner string, statement int) v1review.Finding {
		return v1review.Finding{Code: v1review.OwnerMismatch, Name: kind, Table: name, Detail: owner, Statement: statement}
	}
	for _, tt := range []struct {
		name, sql, role string
		want            []v1review.Finding
	}{
		{name: "public owned by the database owner", sql: "CREATE TABLE t (a int)", role: "alice"},
		{name: "public resolved to its database owner", sql: "CREATE TABLE t (a int)", role: "bob", want: []v1review.Finding{owner("Schema", "public", "alice", 0)}},
		{name: "table owner", sql: "ALTER TABLE orders ADD COLUMN a int", role: "alice", want: []v1review.Finding{owner("Table", "orders", "bob", 0)}},
		{name: "schema owner", sql: "CREATE INDEX i ON sales.t (a)", role: "alice", want: []v1review.Finding{owner("Schema", "sales", "bob", 0)}},
		{name: "database owner", sql: "SELECT 1; CREATE SCHEMA s", role: "bob", want: []v1review.Finding{owner("Database", "shop", "alice", 1)}},
		{name: "unknown schema not judged", sql: "CREATE TABLE nosuch.t (a int)", role: "carol"},
		{name: "drop lists each name", sql: "DROP TABLE sales.a, b", role: "carol", want: []v1review.Finding{owner("Schema", "sales", "bob", 0), owner("Schema", "public", "alice", 0)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := StatementObjectOwnerCheck(parseStatements(t, tt.sql), schema, tt.role)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
