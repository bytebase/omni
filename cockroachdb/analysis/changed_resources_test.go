package analysis

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractChangedResources(t *testing.T) {
	const statement = `CREATE TABLE t1 (c1 INT);
DROP TABLE t1;
ALTER TABLE t1 ADD COLUMN c2 INT;
ALTER TABLE t1 RENAME TO t2;
INSERT INTO t1 (c1) VALUES (1), (5);
UPDATE t1
SET c1 = 5
WHERE c1 = 1;
DELETE FROM t1; -- trailing comment
`
	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if diff := cmp.Diff([]string{"db.public.t1", "db.public.t2"}, resourceNames(got)); diff != "" {
		t.Errorf("resources mismatch (-want +got):\n%s", diff)
	}
	if !resourcesAffect("db.public.t1", got) {
		t.Error("db.public.t1 should affect data")
	}
	if resourcesAffect("db.public.t2", got) {
		t.Error("db.public.t2 should not affect data")
	}
	if got.DMLCount != 2 {
		t.Errorf("DMLCount = %d, want 2", got.DMLCount)
	}
	if diff := cmp.Diff([]string{"UPDATE t1\nSET c1 = 5\nWHERE c1 = 1", "DELETE FROM t1"}, dmlTexts(got)); diff != "" {
		t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
	}
	if got.InsertCount != 2 {
		t.Errorf("InsertCount = %d, want 2", got.InsertCount)
	}
}

func TestExtractChangedResourcesDMLCounts(t *testing.T) {
	for _, tc := range []struct {
		name            string
		statement       string
		wantDMLCount    int
		wantInsertCount int
		wantTables      []string
	}{
		{
			name:            "insert values",
			statement:       `INSERT INTO t VALUES (1), (2), (3)`,
			wantInsertCount: 3,
			wantTables:      []string{"t"},
		},
		{
			name:            "parenthesized insert values",
			statement:       `INSERT INTO t (VALUES (1), (2))`,
			wantInsertCount: 2,
			wantTables:      []string{"t"},
		},
		{
			name:            "upsert values",
			statement:       `UPSERT INTO t VALUES (1), (2)`,
			wantInsertCount: 2,
			wantTables:      []string{"t"},
		},
		{
			name:            "insert values on conflict",
			statement:       `INSERT INTO t VALUES (1), (2) ON CONFLICT (id) DO UPDATE SET c = excluded.c`,
			wantInsertCount: 2,
			wantTables:      []string{"t"},
		},
		{
			name:            "insert default values",
			statement:       `INSERT INTO t DEFAULT VALUES`,
			wantInsertCount: 1,
			wantTables:      []string{"t"},
		},
		{
			name:         "insert select",
			statement:    `INSERT INTO t SELECT * FROM t2 ON CONFLICT DO NOTHING`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "upsert select",
			statement:    `UPSERT INTO t SELECT * FROM t2`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "insert select with read-only CTE",
			statement:    `WITH x AS (SELECT id FROM t2) INSERT INTO t SELECT id FROM x`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "update",
			statement:    `UPDATE t SET c = 1 FROM t2 WHERE t.id = t2.id`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "delete",
			statement:    `DELETE FROM t USING t2 WHERE t.id = t2.id ORDER BY t.id LIMIT 10`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "data-modifying CTE feeding insert",
			statement:    `WITH d AS (DELETE FROM t WHERE id > 0 RETURNING *) INSERT INTO t2 SELECT * FROM d`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2"},
		},
		{
			name:         "data-modifying CTE feeding insert values",
			statement:    `WITH i AS (INSERT INTO t VALUES (1) RETURNING id) INSERT INTO t2 VALUES (2), (3)`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2"},
		},
		{
			name:         "data-modifying CTE in select",
			statement:    `WITH d AS (DELETE FROM t RETURNING 1) SELECT count(*) FROM d`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "data-modifying CTEs feeding delete",
			statement:    `WITH u AS (UPDATE t SET c = 1 RETURNING id), i AS (UPSERT INTO t3 SELECT * FROM t2 RETURNING id) DELETE FROM t2 WHERE id IN (SELECT id FROM u)`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2", "t3"},
		},
		{
			name:         "statement source in select",
			statement:    `SELECT * FROM [DELETE FROM t WHERE c < 10 RETURNING id]`,
			wantDMLCount: 1,
			wantTables:   []string{"t"},
		},
		{
			name:         "statement source in scalar subquery of insert values",
			statement:    `INSERT INTO t2 VALUES ((SELECT max(id) FROM [DELETE FROM t RETURNING id]))`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2"},
		},
		{
			name:         "statement source in update subquery",
			statement:    `UPDATE t SET c = 1 WHERE id IN (SELECT id FROM [DELETE FROM t2 WHERE c < 3 RETURNING id])`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2"},
		},
		{
			name:         "statement source joined in delete using",
			statement:    `DELETE FROM t USING t3 JOIN [UPDATE t2 SET c = 1 RETURNING id] AS u ON t3.id = u.id WHERE t.id = t3.id`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t2"},
		},
		{
			name:         "statement source in create table as",
			statement:    `CREATE TABLE t_moved AS SELECT * FROM [DELETE FROM t RETURNING *]`,
			wantDMLCount: 1,
			wantTables:   []string{"t", "t_moved"},
		},
		{
			name:      "create table as",
			statement: `CREATE TABLE t_copy AS SELECT * FROM t`,
			wantTables: []string{
				"t_copy",
			},
		},
		{
			name:      "read-only CTE in select",
			statement: `WITH x AS (SELECT id FROM t) SELECT * FROM x WHERE id IN (SELECT id FROM [SHOW TABLES])`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractChangedResources(mustParse(t, tc.statement), "db", []string{"public"}, nil)
			var want []string
			for _, table := range tc.wantTables {
				want = append(want, "db.public."+table)
			}
			slices.Sort(want)
			if diff := cmp.Diff(want, resourceNames(got)); diff != "" {
				t.Errorf("resources mismatch (-want +got):\n%s", diff)
			}
			var wantDMLStatements []string
			if tc.wantDMLCount > 0 {
				wantDMLStatements = []string{tc.statement}
			}
			if diff := cmp.Diff(wantDMLStatements, dmlTexts(got)); diff != "" {
				t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
			}
			if got.DMLCount != tc.wantDMLCount {
				t.Errorf("DMLCount = %d, want %d", got.DMLCount, tc.wantDMLCount)
			}
			if got.InsertCount != tc.wantInsertCount {
				t.Errorf("InsertCount = %d, want %d", got.InsertCount, tc.wantInsertCount)
			}
		})
	}
}

func TestExtractChangedResourcesFollowsSearchPath(t *testing.T) {
	const statement = `DELETE FROM t1;
SET search_path = app, public;
DELETE FROM t2;
SET search_path TO '$user', "Reports";
DELETE FROM t3;
RESET search_path;
DELETE FROM t4;`
	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if diff := cmp.Diff([]string{"db.Reports.t3", "db.app.t2", "db.public.t1", "db.public.t4"}, resourceNames(got)); diff != "" {
		t.Errorf("resources mismatch (-want +got):\n%s", diff)
	}
	// A statement after the search path changes carries it for its EXPLAIN.
	if diff := cmp.Diff([]string{
		"DELETE FROM t1",
		"SET LOCAL search_path TO \"app\", \"public\";\nDELETE FROM t2",
		"SET LOCAL search_path TO \"Reports\";\nDELETE FROM t3",
		"DELETE FROM t4",
	}, dmlTexts(got)); diff != "" {
		t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
	}

	t.Run("SET LOCAL lasts until the transaction ends, SET outlasts ROLLBACK, and DISCARD ALL resets", func(t *testing.T) {
		const statement = `BEGIN;
SET LOCAL search_path = app;
DELETE FROM t1;
COMMIT;
DELETE FROM t2;
BEGIN;
SET search_path = app;
SET LOCAL search_path = other;
ROLLBACK;
DELETE FROM t3;
DISCARD ALL;
DELETE FROM t4;`
		got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
		if diff := cmp.Diff([]string{"db.app.t1", "db.app.t3", "db.public.t2", "db.public.t4"}, resourceNames(got)); diff != "" {
			t.Errorf("resources mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("the synced search path applies without a current schema", func(t *testing.T) {
		metadata := newTestMetadata(map[string][]testTable{
			"app":    {{name: "t"}, {name: "mv"}},
			"public": {{name: "t2"}},
		})
		const statement = `DELETE FROM t; DELETE FROM t2; DROP MATERIALIZED VIEW mv; DROP VIEW v;`
		got := ExtractChangedResources(mustParse(t, statement), "db", []string{"app", "public"}, metadata)
		if diff := cmp.Diff([]string{"db.app.mv", "db.app.t", "db.public.t2"}, resourceNames(got)); diff != "" {
			t.Errorf("resources mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("existing tables resolve to the first schema that has them", func(t *testing.T) {
		metadata := newTestMetadata(map[string][]testTable{
			"public": {{name: "t", rows: 1000}},
			"empty":  nil,
		})
		const statement = `SET search_path = empty, public;
DELETE FROM t;
CREATE TABLE t2 (id INT PRIMARY KEY);`
		got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, metadata)
		if diff := cmp.Diff([]string{"db.empty.t2", "db.public.t"}, resourceNames(got)); diff != "" {
			t.Errorf("resources mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestExtractChangedResourcesNames(t *testing.T) {
	for _, tc := range []struct {
		name          string
		statement     string
		currentSchema string
		want          string
	}{
		{name: "unqualified defaults to public", statement: `UPDATE t SET c = 1`, want: "db.public.t"},
		{name: "unqualified uses the current schema", statement: `UPDATE t SET c = 1`, currentSchema: "app", want: "db.app.t"},
		{name: "schema qualified with alias", statement: `INSERT INTO myschema.t AS x VALUES (1)`, currentSchema: "app", want: "db.myschema.t"},
		{name: "database qualified", statement: `DELETE FROM otherdb.myschema.t`, want: "otherdb.myschema.t"},
		{name: "quoted identifiers keep their case", statement: `UPDATE "MySchema"."T" SET c = 1`, want: "db.MySchema.T"},
		{name: "unquoted identifiers are lowercased", statement: `TRUNCATE MySchema.T`, want: "db.myschema.t"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			searchPath := []string{"public"}
			if tc.currentSchema != "" {
				searchPath = []string{tc.currentSchema}
			}
			got := ExtractChangedResources(mustParse(t, tc.statement), "db", searchPath, nil)
			if diff := cmp.Diff([]string{tc.want}, resourceNames(got)); diff != "" {
				t.Errorf("resources mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtractChangedResourcesDDL(t *testing.T) {
	metadata := newTestMetadata(map[string][]testTable{
		"public": {
			{name: "t", rows: 10},
			{name: "t2", rows: 1000, indexes: []string{"t2_c_idx"}},
		},
		"app": {
			{name: "t", rows: 100_000, indexes: []string{"t_c_idx"}},
		},
	})
	for _, tc := range []struct {
		statement        string
		wantTables       []string
		wantAffectedRows int64
	}{
		{statement: `CREATE TABLE t3 (id INT PRIMARY KEY)`, wantTables: []string{"db.public.t3"}},
		{statement: `CREATE INDEX idx ON t2 (c)`, wantTables: []string{"db.public.t2"}},
		{statement: `CREATE MATERIALIZED VIEW mv AS SELECT * FROM t`, wantTables: []string{"db.public.mv"}},
		{statement: `CREATE VIEW v AS SELECT * FROM t`},
		{statement: `IMPORT INTO t2 (c) CSV DATA ('nodelocal://1/t2.csv')`, wantTables: []string{"db.public.t2"}},
		{statement: `ALTER TABLE t2 CONFIGURE ZONE USING num_replicas = 5`, wantTables: []string{"db.public.t2"}, wantAffectedRows: 1000},
		{statement: `ALTER INDEX t2_c_idx CONFIGURE ZONE USING num_replicas = 5`, wantTables: []string{"db.public.t2"}},
		{statement: `ALTER RANGE default CONFIGURE ZONE USING num_replicas = 5`},
		{statement: `DROP INDEX t2@t2_c_idx, t2_c_idx, app.t_c_idx, missing_idx`, wantTables: []string{"db.app.t", "db.public.t2"}},
		{statement: `DROP INDEX db.public.t2_c_idx`, wantTables: []string{"db.public.t2"}},
		{statement: `ALTER INDEX app.t@t_c_idx PARTITION BY NOTHING`, wantTables: []string{"db.app.t"}},
		{statement: `ALTER INDEX t2_c_idx NOT VISIBLE`, wantTables: []string{"db.public.t2"}},
		{statement: `ALTER INDEX t2_c_idx RENAME TO t2_idx`, wantTables: []string{"db.public.t2"}},
		{statement: `ALTER TABLE t ADD COLUMN c2 INT, DROP COLUMN c`, wantTables: []string{"db.public.t"}, wantAffectedRows: 10},
		{statement: `ALTER TABLE t RENAME COLUMN c TO c2`, wantTables: []string{"db.public.t"}, wantAffectedRows: 10},
		{statement: `ALTER TABLE t SET LOCALITY REGIONAL BY ROW`, wantTables: []string{"db.public.t"}, wantAffectedRows: 10},
		{statement: `ALTER TABLE t SET SCHEMA app`, wantTables: []string{"db.public.t"}, wantAffectedRows: 10},
		{statement: `ALTER TABLE t OWNER TO u`, wantTables: []string{"db.public.t"}, wantAffectedRows: 10},
		{statement: `ALTER VIEW t2 SET SCHEMA app`},
		{statement: `ALTER MATERIALIZED VIEW t2 SET SCHEMA app`, wantTables: []string{"db.public.t2"}, wantAffectedRows: 1000},
		{statement: `ALTER MATERIALIZED VIEW t2 OWNER TO u`, wantTables: []string{"db.public.t2"}, wantAffectedRows: 1000},
		{statement: `ALTER TABLE t RENAME TO t_new`, wantTables: []string{"db.public.t", "db.public.t_new"}, wantAffectedRows: 10},
		{statement: `ALTER TABLE app.t RENAME TO t_new`, wantTables: []string{"db.app.t", "db.app.t_new"}, wantAffectedRows: 100_000},
		{statement: `ALTER MATERIALIZED VIEW t2 RENAME TO mv`, wantTables: []string{"db.public.mv", "db.public.t2"}, wantAffectedRows: 1000},
		{statement: `ALTER VIEW t2 RENAME TO v`},
		{statement: `ALTER SEQUENCE t2 RENAME TO s`},
		{statement: `DROP TABLE t, app.t`, wantTables: []string{"db.app.t", "db.public.t"}, wantAffectedRows: 100_010},
		{statement: `TRUNCATE t2`, wantTables: []string{"db.public.t2"}, wantAffectedRows: 1000},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got := ExtractChangedResources(mustParse(t, tc.statement), "db", []string{"public"}, metadata)
			if diff := cmp.Diff(tc.wantTables, resourceNames(got)); diff != "" {
				t.Errorf("resources mismatch (-want +got):\n%s", diff)
			}
			if rows := affectedRows(metadata, got); rows != tc.wantAffectedRows {
				t.Errorf("affected rows = %d, want %d", rows, tc.wantAffectedRows)
			}
			if got.DMLCount != 0 {
				t.Errorf("DMLCount = %d, want 0", got.DMLCount)
			}
		})
	}

	t.Run("an index of another database is recorded by its database", func(t *testing.T) {
		const statement = `DROP INDEX otherdb.public.t2_c_idx`
		got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, metadata)
		if diff := cmp.Diff([]string(nil), resourceNames(got)); diff != "" {
			t.Errorf("resources mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]string{"otherdb"}, got.Databases); diff != "" {
			t.Errorf("databases mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestExtractChangedResourcesExplainAnalyze(t *testing.T) {
	const statement = `EXPLAIN ANALYZE (VERBOSE) DELETE FROM t WHERE id = 1;
EXPLAIN ANALYZE UPSERT INTO t2 SELECT * FROM t;
EXPLAIN UPDATE t3 SET c = 1;`
	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if diff := cmp.Diff([]string{"db.public.t", "db.public.t2"}, resourceNames(got)); diff != "" {
		t.Errorf("resources mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"DELETE FROM t WHERE id = 1", "UPSERT INTO t2 SELECT * FROM t"}, dmlTexts(got)); diff != "" {
		t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
	}
}

func TestExtractChangedResourcesListsEveryDMLStatement(t *testing.T) {
	var statements []string
	for i := range 12 {
		statements = append(statements, fmt.Sprintf("DELETE FROM t WHERE id = %d", i))
	}
	statement := strings.Join(statements, ";\n")

	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if got.DMLCount != 12 {
		t.Errorf("DMLCount = %d, want 12", got.DMLCount)
	}
	if diff := cmp.Diff(statements, dmlTexts(got)); diff != "" {
		t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
	}
}

// resourceNames returns the changed tables as sorted database.schema.table names.
func resourceNames(changes Changes) []string {
	if len(changes.Resources) == 0 {
		return nil
	}
	names := make([]string, 0, len(changes.Resources))
	for _, resource := range changes.Resources {
		names = append(names, resource.Database+"."+resource.Schema+"."+resource.Table)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func resourcesAffect(name string, changes Changes) bool {
	for _, resource := range changes.Resources {
		if resource.Database+"."+resource.Schema+"."+resource.Table == name && resource.AffectsData {
			return true
		}
	}
	return false
}

// dmlTexts returns the DML statements with the search path in force, as the
// EXPLAIN connection receives them.
func dmlTexts(changes Changes) []string {
	if len(changes.DMLStatements) == 0 {
		return nil
	}
	texts := make([]string, 0, len(changes.DMLStatements))
	for _, statement := range changes.DMLStatements {
		text := statement.Text
		if len(statement.SearchPath) > 0 {
			quoted := make([]string, len(statement.SearchPath))
			for i, schema := range statement.SearchPath {
				quoted[i] = `"` + strings.ReplaceAll(schema, `"`, `""`) + `"`
			}
			text = "SET LOCAL search_path TO " + strings.Join(quoted, ", ") + ";\n" + text
		}
		texts = append(texts, text)
	}
	return texts
}

type testTable struct {
	name    string
	rows    int64
	indexes []string
}

type testMetadata struct {
	schemas map[string]map[string]testTable
	indexes map[string]map[string]testTable
}

func newTestMetadata(schemas map[string][]testTable) *testMetadata {
	metadata := &testMetadata{
		schemas: map[string]map[string]testTable{},
		indexes: map[string]map[string]testTable{},
	}
	for schemaName, tables := range schemas {
		metadata.schemas[schemaName] = map[string]testTable{}
		metadata.indexes[schemaName] = map[string]testTable{}
		for _, table := range tables {
			metadata.schemas[schemaName][table.name] = table
			for _, index := range table.indexes {
				metadata.indexes[schemaName][index] = table
			}
		}
	}
	return metadata
}

func (m *testMetadata) RelationSchema(searchPath []string, name string) string {
	for _, schema := range searchPath {
		if _, ok := m.schemas[schema][name]; ok {
			return schema
		}
	}
	return ""
}

func (m *testMetadata) IndexTable(searchPath []string, name string) (string, string, bool) {
	for _, schema := range searchPath {
		if table, ok := m.indexes[schema][name]; ok {
			return schema, table.name, true
		}
	}
	return "", "", false
}

func affectedRows(metadata *testMetadata, changes Changes) int64 {
	var rows int64
	seen := map[string]bool{}
	for _, resource := range changes.Resources {
		if !resource.AffectsData {
			continue
		}
		key := resource.Database + "." + resource.Schema + "." + resource.Table
		if seen[key] {
			continue
		}
		seen[key] = true
		rows += metadata.schemas[resource.Schema][resource.Table].rows
	}
	return rows
}

func TestExtractChangedResourcesSavepoint(t *testing.T) {
	const statement = `BEGIN;
SAVEPOINT s;
SET LOCAL search_path = app;
ROLLBACK TO SAVEPOINT s;
DELETE FROM t;
COMMIT;`
	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if diff := cmp.Diff([]string{"db.public.t"}, resourceNames(got)); diff != "" {
		t.Errorf("resources mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"DELETE FROM t"}, dmlTexts(got)); diff != "" {
		t.Errorf("DMLStatements mismatch (-want +got):\n%s", diff)
	}
}

func TestExtractChangedResourcesInsertWithLimit(t *testing.T) {
	for _, statement := range []string{
		`INSERT INTO t VALUES (1), (2) LIMIT 1`,
		`INSERT INTO t VALUES (1), (2) OFFSET 1`,
		`INSERT INTO t (VALUES (1), (2)) LIMIT 0`,
	} {
		t.Run(statement, func(t *testing.T) {
			got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
			if got.InsertCount != 0 {
				t.Errorf("InsertCount = %d, want 0", got.InsertCount)
			}
			if got.DMLCount != 1 {
				t.Errorf("DMLCount = %d, want 1", got.DMLCount)
			}
		})
	}
}

func TestExtractChangedResourcesReleaseSavepoint(t *testing.T) {
	const statement = `BEGIN;
SAVEPOINT s;
SET LOCAL search_path = app;
RELEASE SAVEPOINT s;
DELETE FROM t;
COMMIT;`
	got := ExtractChangedResources(mustParse(t, statement), "db", []string{"public"}, nil)
	if diff := cmp.Diff([]string{"db.public.t"}, resourceNames(got)); diff != "" {
		t.Errorf("resources mismatch (-want +got):\n%s", diff)
	}
}
