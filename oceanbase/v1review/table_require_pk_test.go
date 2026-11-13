package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

type pkFinding struct {
	Name      string
	Statement int
	Text      string
}

func requirePK(t *testing.T, original *v1review.Schema, sqls ...string) []pkFinding {
	t.Helper()
	var stmts []ast.Node
	for _, sql := range sqls {
		list, err := parser.Parse(sql)
		if err != nil {
			t.Fatalf("Parse %q: %v", sql, err)
		}
		stmts = append(stmts, list.Items...)
	}
	var got []pkFinding
	for _, f := range TableRequirePK(stmts, original) {
		if f.Code != v1review.TablePrimaryKeyMissing {
			t.Fatalf("code %v", f.Code)
		}
		text := ""
		if f.Statement >= 0 {
			text = sqls[f.Statement][f.Start:f.End]
		}
		got = append(got, pkFinding{f.Name, f.Statement, text})
	}
	return got
}

func TestTableRequirePK(t *testing.T) {
	pk := func(name string, cols ...string) *metadata.IndexMetadata {
		return &metadata.IndexMetadata{Name: name, Primary: true, Expressions: cols}
	}
	original := &v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "book", Indexes: []*metadata.IndexMetadata{pk("PRIMARY", "id", "name")}},
			{Name: "lower", Indexes: []*metadata.IndexMetadata{pk("primary", "a")}},
			{Name: "bare", Indexes: []*metadata.IndexMetadata{pk("PRIMARY")}},
			{Name: "named", Indexes: []*metadata.IndexMetadata{pk("pk_x", "a")}},
		},
	}}}}
	for _, tt := range []struct {
		name string
		sqls []string
		want []pkFinding
	}{
		{name: "column key", sqls: []string{"CREATE TABLE t (id INT PRIMARY KEY)"}},
		{name: "table key", sqls: []string{"CREATE TABLE t (id INT, PRIMARY KEY (id))"}},
		{name: "no key", sqls: []string{"CREATE TABLE t (id INT)"}, want: []pkFinding{{"t", 0, "CREATE TABLE t (id INT)"}}},
		{name: "like and as-select have no key", sqls: []string{"CREATE TABLE t2 LIKE book", "CREATE TABLE t3 AS SELECT 1"},
			want: []pkFinding{{"t2", 0, "CREATE TABLE t2 LIKE book"}, {"t3", 1, "CREATE TABLE t3 AS SELECT 1"}}},
		{name: "dropped table forgotten", sqls: []string{"CREATE TABLE t (id INT)", "DROP TABLE t"}},
		{name: "key added later", sqls: []string{"CREATE TABLE t (id INT)", "ALTER TABLE t ADD PRIMARY KEY (id)"}},
		{name: "key column added later", sqls: []string{"CREATE TABLE t (id INT)", "ALTER TABLE t ADD COLUMN (x INT, pid INT PRIMARY KEY)"}},
		{name: "drop primary key anchors the alter", sqls: []string{"ALTER TABLE book DROP PRIMARY KEY"}, want: []pkFinding{{"book", 0, "ALTER TABLE book DROP PRIMARY KEY"}}},
		{name: "drop index primary in any case", sqls: []string{"ALTER TABLE book DROP INDEX `primary`"}, want: []pkFinding{{"book", 0, "ALTER TABLE book DROP INDEX `primary`"}}},
		{name: "dropping one key column leaves the other", sqls: []string{"ALTER TABLE book DROP COLUMN name"}},
		{name: "dropping every key column", sqls: []string{"ALTER TABLE book DROP COLUMN id", "ALTER TABLE book DROP COLUMN name"},
			want: []pkFinding{{"book", 1, "ALTER TABLE book DROP COLUMN name"}}},
		{name: "a renamed key column stays in the key", sqls: []string{"ALTER TABLE book CHANGE COLUMN id uid INT", "ALTER TABLE book DROP COLUMN name"}},
		{name: "table names as written, original looked up case-insensitively", sqls: []string{"ALTER TABLE BOOK DROP COLUMN id", "ALTER TABLE book DROP COLUMN name"}},
		{name: "index name compared case-insensitively", sqls: []string{"ALTER TABLE lower DROP COLUMN a"}, want: []pkFinding{{"lower", 0, "ALTER TABLE lower DROP COLUMN a"}}},
		{name: "primary key not named PRIMARY not found", sqls: []string{"ALTER TABLE named DROP COLUMN a"}},
		{name: "empty original key has no anchor", sqls: []string{"ALTER TABLE bare DROP COLUMN a"}, want: []pkFinding{{"bare", -1, ""}}},
		{name: "findings in table name order", sqls: []string{"CREATE TABLE b (x INT)", "CREATE TABLE a (x INT)"},
			want: []pkFinding{{"a", 1, "CREATE TABLE a (x INT)"}, {"b", 0, "CREATE TABLE b (x INT)"}}},
		{name: "database not part of the name", sqls: []string{"CREATE TABLE d1.t (x INT)", "CREATE TABLE d2.t (y INT PRIMARY KEY)"}},
		{name: "routine body not judged", sqls: []string{"CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT); END"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, requirePK(t, original, tt.sqls...)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableRequirePKNilOriginalPanicsOnlyWhenLookingUp(t *testing.T) {
	if got := requirePK(t, nil, "CREATE TABLE t (id INT)", "ALTER TABLE t DROP COLUMN id"); len(got) != 1 {
		t.Errorf("got %v", got)
	}
	defer func() {
		if recover() == nil {
			t.Error("want a panic")
		}
	}()
	requirePK(t, nil, "ALTER TABLE book DROP COLUMN id")
}
