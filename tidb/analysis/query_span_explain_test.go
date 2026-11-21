package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/analysis"
)

func buildTiDBCatalog(dbs []*metadata.DatabaseSchemaMetadata) analysis.Catalog {
	byName := make(map[string]*metadata.DatabaseSchemaMetadata, len(dbs))
	names := make([]string, 0, len(dbs))
	for _, db := range dbs {
		byName[db.GetName()] = db
		names = append(names, db.GetName())
	}
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) { return names, nil },
		GetDatabase: func(name string) (analysis.Database, error) {
			if db, ok := byName[name]; ok {
				return &tidbMetadataDB{proto: db}, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
	}
}

type tidbMetadataDB struct {
	proto *metadata.DatabaseSchemaMetadata
}

func (d *tidbMetadataDB) Name() string { return d.proto.GetName() }

func (d *tidbMetadataDB) HasSchema() bool { return d.schema() != nil }

func (d *tidbMetadataDB) ListTableNames() []string {
	var names []string
	for _, t := range d.schema().GetTables() {
		names = append(names, t.GetName())
	}
	return names
}

func (d *tidbMetadataDB) GetTable(name string) *metadata.TableMetadata {
	for _, t := range d.schema().GetTables() {
		if strings.EqualFold(t.GetName(), name) {
			return t
		}
	}
	return nil
}

func (d *tidbMetadataDB) ListViewNames() []string {
	var names []string
	for _, v := range d.schema().GetViews() {
		names = append(names, v.GetName())
	}
	return names
}

func (d *tidbMetadataDB) GetView(name string) *metadata.ViewMetadata {
	for _, v := range d.schema().GetViews() {
		if strings.EqualFold(v.GetName(), name) {
			return v
		}
	}
	return nil
}

func (d *tidbMetadataDB) schema() *metadata.SchemaMetadata {
	for _, s := range d.proto.GetSchemas() {
		if s.GetName() == "" {
			return s
		}
	}
	return nil
}

func TestExplainByDigestOrConnectionRejected(t *testing.T) {
	span, err := analysis.ResolveQuerySpan("EXPLAIN FOR CONNECTION 5", analysis.ResolveOptions{DefaultDatabase: "db"})
	require.ErrorContains(t, err, "EXPLAIN by digest or connection is not supported")
	require.Nil(t, span)

	span, err = analysis.ResolveQuerySpan("EXPLAIN ANALYZE 'sqldigest'", analysis.ResolveOptions{DefaultDatabase: "db"})
	require.Error(t, err)
	require.Nil(t, span)
}

func TestExplainAnalyzeAccessTables(t *testing.T) {
	var dbs []*metadata.DatabaseSchemaMetadata
	for _, database := range []string{"db", "other"} {
		dbs = append(dbs, &metadata.DatabaseSchemaMetadata{
			Name: database,
			Schemas: []*metadata.SchemaMetadata{{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
					{Name: "u", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
				},
			}},
		})
	}
	dbs = append(dbs, &metadata.DatabaseSchemaMetadata{
		Name:    "information_schema",
		Schemas: []*metadata.SchemaMetadata{{Name: "", Tables: []*metadata.TableMetadata{{Name: "tables"}}}},
	})
	catalog := buildTiDBCatalog(dbs)

	for _, tc := range []struct {
		name      string
		statement string
		queryType analysis.QueryType
		sources   analysis.ColumnSet
		wantMixed bool
	}{
		{
			name:      "count with predicate",
			statement: "EXPLAIN ANALYZE SELECT count(*) FROM t WHERE id = 839195265",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true},
		},
		{
			name:      "join and cross database subquery",
			statement: "EXPLAIN ANALYZE SELECT a.id FROM t a JOIN u b ON a.id = b.id WHERE a.id IN (SELECT id FROM other.t)",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "db", Table: "u"}: true, {Database: "other", Table: "t"}: true},
		},
		{
			name:      "derived join retains every table",
			statement: "EXPLAIN ANALYZE SELECT * FROM (SELECT t.id FROM t JOIN u ON t.id = u.id) x",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "db", Table: "u"}: true},
		},
		{
			name:      "nested derived join across databases",
			statement: "EXPLAIN ANALYZE SELECT * FROM (SELECT * FROM (SELECT t.id FROM t JOIN other.t u ON t.id = u.id) x) y",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "other", Table: "t"}: true},
		},
		{
			name:      "derived predicate subquery",
			statement: "EXPLAIN ANALYZE SELECT * FROM (SELECT id FROM t WHERE id IN (SELECT id FROM other.u)) x",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "other", Table: "u"}: true},
		},
		{
			name:      "subqueries inside compound expressions",
			statement: "EXPLAIN ANALYZE SELECT COALESCE((SELECT max(id) FROM other.t), 0) FROM t WHERE id > 0 AND id IN (SELECT id FROM u)",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "db", Table: "u"}: true, {Database: "other", Table: "t"}: true},
		},
		{
			name:      "qualified table",
			statement: "EXPLAIN ANALYZE SELECT count(*) FROM other.t",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "other", Table: "t"}: true},
		},
		{
			name:      "uncached joined tables",
			statement: "EXPLAIN ANALYZE SELECT * FROM t JOIN missing ON TRUE JOIN other.missing2 ON TRUE",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "db", Table: "missing"}: true, {Database: "other", Table: "missing2"}: true},
		},
		{
			name:      "uncached nested predicate",
			statement: "EXPLAIN ANALYZE SELECT * FROM t WHERE id IN (SELECT id FROM unsynced.missing)",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "unsynced", Table: "missing"}: true},
		},
		{
			name:      "CTE references are not physical targets",
			statement: "EXPLAIN ANALYZE WITH c AS (SELECT * FROM missing) SELECT * FROM c JOIN t ON TRUE",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "missing"}: true, {Database: "db", Table: "t"}: true},
		},
		{
			name:      "qualified table sharing CTE name",
			statement: "EXPLAIN ANALYZE WITH c AS (SELECT * FROM t) SELECT * FROM c JOIN other.c ON TRUE",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "other", Table: "c"}: true},
		},
		{
			name:      "non recursive CTE reads same named physical table",
			statement: "EXPLAIN ANALYZE WITH t AS (SELECT * FROM t) SELECT * FROM t",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true},
		},
		{
			name:      "recursive CTE reference",
			statement: "EXPLAIN ANALYZE WITH RECURSIVE c AS (SELECT id FROM missing UNION ALL SELECT id + 1 FROM c WHERE id < 5) SELECT * FROM c",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "missing"}: true},
		},
		{
			name:      "nested CTE does not hide outer physical table",
			statement: "EXPLAIN ANALYZE SELECT * FROM (WITH missing AS (SELECT * FROM t) SELECT * FROM missing) x JOIN missing ON TRUE",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "db", Table: "missing"}: true},
		},
		{
			name:      "plain explain permits mixed tables",
			statement: "EXPLAIN SELECT * FROM t JOIN information_schema.tables ON TRUE",
			queryType: analysis.Explain,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true},
		},
		{
			name:      "analyze still rejects mixed tables",
			statement: "EXPLAIN ANALYZE SELECT * FROM t JOIN information_schema.tables ON TRUE",
			wantMixed: true,
		},
		{
			name:      "locking aliases are not physical tables",
			statement: "EXPLAIN ANALYZE SELECT * FROM t AS a JOIN other.t AS b ON a.id = b.id FOR UPDATE OF a, b",
			queryType: analysis.Select,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true, {Database: "other", Table: "t"}: true},
		},
		{
			name:      "plain explain",
			statement: "EXPLAIN SELECT * FROM t",
			queryType: analysis.Explain,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true},
		},
		{
			name:      "analyzed write remains dml",
			statement: "EXPLAIN ANALYZE DELETE FROM t WHERE id = 1",
			queryType: analysis.DML,
			sources:   analysis.ColumnSet{{Database: "db", Table: "t"}: true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			span, err := analysis.ResolveQuerySpan(tc.statement, analysis.ResolveOptions{
				DefaultDatabase: "db",
				Catalog:         catalog,
			})
			if tc.wantMixed {
				require.ErrorIs(t, err, analysis.ErrMixUserSystemTables)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.queryType, span.Type)
			require.Equal(t, strings.HasPrefix(tc.statement, "EXPLAIN ANALYZE"), span.ExplainAnalyze)
			require.Equal(t, tc.sources, span.SourceColumns)
			require.Empty(t, span.Results)
		})
	}
}

func TestRecursiveCTEWithInnerWith(t *testing.T) {
	catalog := buildTiDBCatalog([]*metadata.DatabaseSchemaMetadata{{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "n"}}}},
		}},
	}})
	statement := "WITH RECURSIVE r(n) AS (WITH seed AS (SELECT n FROM t) SELECT n FROM seed UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r"
	span, err := analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{DefaultDatabase: "db", Catalog: catalog})
	require.NoError(t, err)
	require.Nil(t, span.NotFoundError)
	require.Len(t, span.Results, 1)
	require.Equal(t, analysis.ColumnSet{{Database: "db", Table: "t", Column: "n"}: true}, span.Results[0].SourceColumns)
}

func TestRecursiveCTEShadowedInitialReference(t *testing.T) {
	catalog := buildTiDBCatalog([]*metadata.DatabaseSchemaMetadata{{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id"}}}},
		}},
	}})
	statement := "WITH RECURSIVE c(id) AS (SELECT COALESCE((SELECT id FROM (WITH c AS (SELECT id FROM t) SELECT id FROM c) x), 0) UNION ALL SELECT id + 1 FROM c WHERE id < 10) SELECT id FROM c"
	span, err := analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{DefaultDatabase: "db", Catalog: catalog})
	require.NoError(t, err)
	require.Nil(t, span.NotFoundError)
	require.Len(t, span.Results, 1)
	require.Equal(t, analysis.ColumnSet{{Database: "db", Table: "t"}: true}, span.SourceColumns)
}
