package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

type onlineMigrationStatement struct {
	node ast.Node
	text string
}

type onlineMigrationFinding struct {
	Code      v1review.Code
	Table     string
	Value     string
	Statement int
	Text      string
}

// onlineMigrationSnapshot is a snapshot whose first schema holds the tables.
func onlineMigrationSnapshot(tables ...*metadata.TableMetadata) v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Name:    "test",
		Schemas: []*metadata.SchemaMetadata{{Tables: tables}},
	}}
}

func tableWithRows(name string, rows int64) *metadata.TableMetadata {
	return &metadata.TableMetadata{Name: name, RowCount: rows}
}

// judgeOnlineMigration splits sql the way Bytebase does, parses each statement
// on its own, and runs OnlineMigration over them.
func judgeOnlineMigration(t *testing.T, sql string, schema v1review.Schema, currentDatabase string, minRows int64) []onlineMigrationFinding {
	t.Helper()
	var stmts []onlineMigrationStatement
	for _, segment := range parser.Split(sql) {
		if segment.Empty() {
			continue
		}
		list, err := parser.Parse(segment.Text)
		if err != nil {
			t.Fatalf("Parse %q: %v", segment.Text, err)
		}
		for _, node := range list.Items {
			stmts = append(stmts, onlineMigrationStatement{node: node, text: segment.Text})
		}
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
		nodes = append(nodes, stmt.node)
	}
	var got []onlineMigrationFinding
	for _, f := range OnlineMigration(nodes, schema, currentDatabase, minRows) {
		got = append(got, onlineMigrationFinding{
			Code:      f.Code,
			Table:     f.Table,
			Value:     f.Value,
			Statement: f.Statement,
			Text:      stmts[f.Statement].text[f.Start:f.End],
		})
	}
	return got
}

func TestOnlineMigration(t *testing.T) {
	needed := func(table, rows, text string, stmt int) onlineMigrationFinding {
		return onlineMigrationFinding{Code: v1review.OnlineMigrationNeeded, Table: table, Value: rows, Statement: stmt, Text: text}
	}
	big := []*metadata.TableMetadata{tableWithRows("t", 1000), tableWithRows("small", 3)}
	for _, tt := range []struct {
		name            string
		sql             string
		schema          v1review.Schema
		currentDatabase string
		minRows         int64
		want            []onlineMigrationFinding
	}{
		{name: "under the limit", sql: "ALTER TABLE t ADD b INT;", schema: onlineMigrationSnapshot(big...), currentDatabase: "test", minRows: 1001},
		{
			name:            "at the limit",
			sql:             "ALTER TABLE t ADD b INT;",
			schema:          onlineMigrationSnapshot(big...),
			currentDatabase: "test",
			minRows:         1000,
			want:            []onlineMigrationFinding{needed("test.t", "1000", "ALTER TABLE t ADD b INT", 0)},
		},
		{
			name:            "over the limit",
			sql:             "SELECT 1;\nALTER TABLE t ADD b INT;",
			schema:          onlineMigrationSnapshot(big...),
			currentDatabase: "test",
			minRows:         100,
			want:            []onlineMigrationFinding{needed("test.t", "1000", "ALTER TABLE t ADD b INT", 1)},
		},
		{name: "no command", sql: "ALTER TABLE t;", schema: onlineMigrationSnapshot(big...), currentDatabase: "test", minRows: 100},
		{name: "unknown table", sql: "ALTER TABLE nope ADD b INT;", schema: onlineMigrationSnapshot(big...), currentDatabase: "test", minRows: 100},
		{name: "other statements", sql: "CREATE TABLE u (a INT);\nRENAME TABLE t TO u;", schema: onlineMigrationSnapshot(big...), currentDatabase: "test", minRows: 100},
		{
			name:            "a qualified table of this database",
			sql:             "ALTER TABLE test.t ADD b INT;",
			schema:          onlineMigrationSnapshot(big...),
			currentDatabase: "other",
			minRows:         100,
			want:            []onlineMigrationFinding{needed("test.t", "1000", "ALTER TABLE test.t ADD b INT", 0)},
		},
		{
			name:            "a qualified name differing in case",
			sql:             "ALTER TABLE TEST.t ADD b INT;",
			schema:          onlineMigrationSnapshot(big...),
			currentDatabase: "active",
			minRows:         100,
			want:            []onlineMigrationFinding{needed("TEST.t", "1000", "ALTER TABLE TEST.t ADD b INT", 0)},
		},
		{
			name:            "another database's table is not judged",
			sql:             "ALTER TABLE other.t ADD b INT;",
			schema:          onlineMigrationSnapshot(tableWithRows("small", 3), big[0]),
			currentDatabase: "test",
			minRows:         100,
		},
		{
			name:    "an unnamed snapshot judges no qualified table",
			sql:     "ALTER TABLE test.t ADD b INT;",
			schema:  v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Tables: big}}}},
			minRows: 100,
		},
		{
			name:    "without a current database",
			sql:     "ALTER TABLE t ADD b INT;",
			schema:  onlineMigrationSnapshot(big...),
			minRows: 100,
			want:    []onlineMigrationFinding{needed(".t", "1000", "ALTER TABLE t ADD b INT", 0)},
		},
		{
			// V1 read the snapshot's first schema only.
			name: "the first schema only",
			sql:  "ALTER TABLE t ADD b INT;",
			schema: v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
				{Tables: []*metadata.TableMetadata{tableWithRows("u", 1000)}},
				{Tables: []*metadata.TableMetadata{tableWithRows("t", 1000)}},
			}}},
			currentDatabase: "test",
			minRows:         100,
		},
		{
			// V1 stopped at the first table of that name.
			name:            "the first table of that name wins",
			sql:             "ALTER TABLE t ADD b INT;",
			schema:          onlineMigrationSnapshot(tableWithRows("t", 1), tableWithRows("t", 1000)),
			currentDatabase: "test",
			minRows:         100,
		},
		{name: "no schema", sql: "ALTER TABLE t ADD b INT;", schema: v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: "test"}}, currentDatabase: "test", minRows: 100},
		{name: "no snapshot", sql: "ALTER TABLE t ADD b INT;", currentDatabase: "test", minRows: 100},
		{name: "no statements", schema: onlineMigrationSnapshot(big...), currentDatabase: "test", minRows: 100},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeOnlineMigration(t, tt.sql, tt.schema, tt.currentDatabase, tt.minRows)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
