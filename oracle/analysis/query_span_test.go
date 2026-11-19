package analysis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	oracleast "github.com/bytebase/omni/oracle/ast"
)

type fakeDatabase struct {
	tables            map[string]*metadata.TableMetadata
	views             map[string]*metadata.ViewMetadata
	materializedViews map[string]*metadata.MaterializedViewMetadata
	externalTables    map[string]*metadata.ExternalTableMetadata
}

func (*fakeDatabase) HasSchema() bool { return true }

func (d *fakeDatabase) GetTable(name string) *metadata.TableMetadata { return d.tables[name] }

func (d *fakeDatabase) GetView(name string) *metadata.ViewMetadata { return d.views[name] }

func (d *fakeDatabase) GetMaterializedView(name string) *metadata.MaterializedViewMetadata {
	return d.materializedViews[name]
}

func (d *fakeDatabase) GetExternalTable(name string) *metadata.ExternalTableMetadata {
	return d.externalTables[name]
}

func columns(names ...string) []*metadata.ColumnMetadata {
	var result []*metadata.ColumnMetadata
	for _, name := range names {
		result = append(result, &metadata.ColumnMetadata{Name: name})
	}
	return result
}

// testCatalog holds schema DB on instance A, and schema R on instance B,
// which link LINK reaches from A.
func testCatalog() Catalog {
	instances := map[string]map[string]*fakeDatabase{
		"A": {
			"DB": {
				tables: map[string]*metadata.TableMetadata{
					"T": {Name: "T", Columns: []*metadata.ColumnMetadata{{Name: "A"}, {Name: "H", IsInvisible: true}, {Name: "B"}}},
				},
				views: map[string]*metadata.ViewMetadata{
					"V":       {Name: "V", Definition: "SELECT A, A + B AS S FROM T"},
					"AV_COLS": {Name: "AV_COLS", Columns: columns("A", "S")},
					"V_EMPTY": {Name: "V_EMPTY"},
					"V_MIXED": {Name: "V_MIXED", Definition: "SELECT * FROM T, SYS.DBA_USERS"},
				},
				materializedViews: map[string]*metadata.MaterializedViewMetadata{
					"MV": {Name: "MV", Definition: "SELECT B FROM T"},
				},
				externalTables: map[string]*metadata.ExternalTableMetadata{
					"EXT": {Name: "EXT", Columns: columns("E")},
				},
			},
			"SYS": {tables: map[string]*metadata.TableMetadata{"DBA_USERS": {Name: "DBA_USERS", Columns: columns("USERNAME")}}},
		},
		"B": {
			"R": {
				tables: map[string]*metadata.TableMetadata{"RT": {Name: "RT", Columns: columns("X")}},
				views:  map[string]*metadata.ViewMetadata{"RV": {Name: "RV", Definition: "SELECT X FROM RT"}},
			},
		},
	}
	return Catalog{
		GetDatabase: func(instance, database string) (Database, error) {
			if db, ok := instances[instance][database]; ok {
				return db, nil
			}
			return nil, errors.New("not synced")
		},
		GetLinkedDatabase: func(_, link, schema string) (string, string, Database, error) {
			if link != "LINK" {
				return "", "", nil, nil
			}
			if schema == "" {
				schema = "R"
			}
			if db, ok := instances["B"][schema]; ok {
				return "B", schema, db, nil
			}
			return "", "", nil, nil
		},
	}
}

func resolve(t *testing.T, statement string) *ResolvedSpan {
	t.Helper()
	span, err := ResolveQuerySpan(statement, ResolveOptions{Instance: "A", DefaultDatabase: "DB", Catalog: testCatalog()})
	require.NoError(t, err)
	return span
}

func column(database, table, name string) ColumnSet {
	return ColumnSet{{Database: database, Table: table, Column: name}: true}
}

func TestResolveQuerySpanLineage(t *testing.T) {
	span := resolve(t, "SELECT * FROM T")
	require.Equal(t, Select, span.Type)
	require.Equal(t, ColumnSet{{Database: "DB", Table: "T"}: true}, span.SourceColumns)
	require.Equal(t, []QuerySpanResult{
		{Name: "A", SourceColumns: column("DB", "T", "A"), IsPlainField: true},
		{Name: "B", SourceColumns: column("DB", "T", "B"), IsPlainField: true},
	}, span.Results)

	span = resolve(t, "SELECT H FROM T")
	require.Equal(t, []QuerySpanResult{{Name: "H", SourceColumns: column("DB", "T", "H")}}, span.Results)

	span = resolve(t, "SELECT S FROM V")
	require.Equal(t, "S", span.Results[0].Name)
	require.Equal(t, ColumnSet{
		{Database: "DB", Table: "T", Column: "A"}: true,
		{Database: "DB", Table: "T", Column: "B"}: true,
	}, span.Results[0].SourceColumns)

	span = resolve(t, "SELECT * FROM MV, EXT")
	require.Equal(t, []string{"B", "E"}, []string{span.Results[0].Name, span.Results[1].Name})
	require.Equal(t, column("DB", "EXT", "E"), span.Results[1].SourceColumns)
}

// TestResolveQuerySpanDerivedAnalyticView holds both forms of a derived
// analytic view to the persistent analytic view each names: an inline one in
// a FROM clause, and a WITH element, whose span is the base view's under the
// element's name.
func TestResolveQuerySpanDerivedAnalyticView(t *testing.T) {
	span := resolve(t, "SELECT * FROM ANALYTIC VIEW(USING V)")
	require.Equal(t, ColumnSet{{Database: "DB", Table: "V"}: true}, span.SourceColumns)
	require.Equal(t, []string{"A", "S"}, []string{span.Results[0].Name, span.Results[1].Name})
	require.Equal(t, column("DB", "T", "A"), span.Results[0].SourceColumns)
	require.Equal(t, ColumnSet{
		{Database: "DB", Table: "T", Column: "A"}: true,
		{Database: "DB", Table: "T", Column: "B"}: true,
	}, span.Results[1].SourceColumns)

	span = resolve(t, "SELECT * FROM ANALYTIC VIEW(USING V) AV WHERE AV.A = 1")
	require.Equal(t, ColumnSet{{Database: "DB", Table: "V"}: true}, span.SourceColumns)
	require.Equal(t, []string{"A", "S"}, []string{span.Results[0].Name, span.Results[1].Name})

	span = resolve(t, "WITH AV ANALYTIC VIEW AS(USING V) SELECT A FROM AV")
	require.Equal(t, ColumnSet{{Database: "DB", Table: "V"}: true}, span.SourceColumns)
	require.Equal(t, []QuerySpanResult{{Name: "A", SourceColumns: column("DB", "T", "A")}}, span.Results)

	// ADD MEASURES adds a column to the derived view, beside the base's, and
	// keeps the base view's qualifier when no alias renames it.
	span = resolve(t, "SELECT V.S, V.LAG_SALES FROM ANALYTIC VIEW(USING V ADD MEASURES(LAG_SALES AS(1)))")
	require.Equal(t, []string{"S", "LAG_SALES"}, []string{span.Results[0].Name, span.Results[1].Name})
	require.NotEmpty(t, span.Results[0].SourceColumns)
	require.Empty(t, span.Results[1].SourceColumns)

	// A parsed measure carries the lineage of what it reads, the measures
	// before it included; one the parser read whole carries none.
	span = resolve(t, "SELECT M FROM ANALYTIC VIEW(USING V ADD MEASURES(M AS(A)))")
	require.Equal(t, []QuerySpanResult{{Name: "M", SourceColumns: column("DB", "T", "A")}}, span.Results)

	span = resolve(t, "SELECT N FROM ANALYTIC VIEW(USING V ADD MEASURES(M AS(A), N AS(M + 1)))")
	require.Equal(t, []QuerySpanResult{{Name: "N", SourceColumns: column("DB", "T", "A")}}, span.Results)

	span = resolve(t, "SELECT M FROM ANALYTIC VIEW(USING V ADD MEASURES(M AS(LAG(S) OVER (HIERARCHY H OFFSET 1))))")
	require.Equal(t, []QuerySpanResult{{Name: "M", SourceColumns: ColumnSet{}}}, span.Results)

	// USING names a persistent analytic view in the schema, never a CTE of the
	// same name.
	span = resolve(t, "WITH V AS (SELECT 1 AS X FROM T) SELECT A FROM ANALYTIC VIEW(USING V)")
	require.Equal(t, []QuerySpanResult{{Name: "A", SourceColumns: column("DB", "T", "A")}}, span.Results)

	span = resolve(t, "WITH NOPE AS (SELECT 1 AS X FROM T) SELECT X FROM NOPE WHERE EXISTS (SELECT 1 FROM ANALYTIC VIEW(USING NOPE))")
	require.NotNil(t, span.NotFoundError)
	require.Equal(t, "NOPE", *span.NotFoundError.Table)

	// A base view with no stored definition yields an empty span rather than
	// cloning a nil source; one carrying synced column metadata resolves them.
	span = resolve(t, "WITH AV ANALYTIC VIEW AS(USING V_EMPTY) SELECT * FROM AV")
	require.Equal(t, ColumnSet{{Database: "DB", Table: "V_EMPTY"}: true}, span.SourceColumns)
	require.Empty(t, span.Results)

	span = resolve(t, "SELECT S FROM ANALYTIC VIEW(USING AV_COLS)")
	require.Equal(t, []QuerySpanResult{{Name: "S", SourceColumns: column("DB", "AV_COLS", "S")}}, span.Results)

	// USING DUAL names a persistent analytic view rather than the pseudo-table.
	span = resolve(t, "SELECT * FROM ANALYTIC VIEW(USING DUAL)")
	require.NotNil(t, span.NotFoundError)
	require.Equal(t, "DUAL", *span.NotFoundError.Table)

	span = resolve(t, "DELETE T FROM ANALYTIC VIEW(USING DUAL)")
	require.Equal(t, ColumnSet{{Database: "DB", Table: "T"}: true, {Database: "DB", Table: "DUAL"}: true}, span.SourceColumns)
}

func TestResolveQuerySpanDatabaseLink(t *testing.T) {
	span := resolve(t, "SELECT * FROM RV@LINK")
	require.Equal(t, ColumnSet{{Instance: "B", Server: "LINK", Database: "R", Table: "RV"}: true}, span.SourceColumns)
	require.Equal(t, []QuerySpanResult{{Name: "X", SourceColumns: column("R", "RT", "X")}}, span.Results)

	span = resolve(t, "SELECT * FROM T@NOWHERE")
	require.Equal(t, ColumnSet{{Server: "NOWHERE", Table: "T"}: true}, span.SourceColumns)
	require.NotNil(t, span.NotFoundError)
	require.Equal(t, "NOWHERE", *span.NotFoundError.DatabaseLink)
}

func TestResolveQuerySpanNotFound(t *testing.T) {
	span := resolve(t, "SELECT * FROM NOPE")
	require.Empty(t, span.Results)
	require.Equal(t, ColumnSet{{Database: "DB"}: true}, span.SourceColumns)
	require.Equal(t, "resource not found: database: DB, table: NOPE", span.NotFoundError.Error())
}

func TestResolveQuerySpanStatementKinds(t *testing.T) {
	span, err := ResolveQuerySpan("", ResolveOptions{Instance: "A", DefaultDatabase: "DB", Catalog: testCatalog()})
	require.NoError(t, err)
	require.Nil(t, span)

	span = resolve(t, "SELECT USERNAME FROM SYS.DBA_USERS")
	require.Equal(t, SelectInfoSchema, span.Type)

	span = resolve(t, "DELETE FROM T WHERE A = 1")
	require.Equal(t, DML, span.Type)
	require.Equal(t, ColumnSet{{Database: "DB", Table: "T"}: true}, span.SourceColumns)
	require.Empty(t, span.Results)

	_, err = ResolveQuerySpan("SELECT * FROM T, SYS.DBA_USERS", ResolveOptions{Instance: "A", DefaultDatabase: "DB", Catalog: testCatalog()})
	require.Same(t, ErrMixUserSystemTables, err)

	_, err = ResolveQuerySpan("SELECT * FROM V_MIXED", ResolveOptions{Instance: "A", DefaultDatabase: "DB", Catalog: testCatalog()})
	require.ErrorIs(t, err, ErrMixUserSystemTables)
	require.NotSame(t, ErrMixUserSystemTables, err)
}

func TestCollectAccessTables(t *testing.T) {
	tests := []struct {
		statement string
		expected  []accessResource
		linked    []accessResource
	}{
		{
			statement: "SELECT * FROM t1 WHERE c1 = 1",
			expected: []accessResource{
				{
					Database: "DB",
					Table:    "T1",
				},
			},
		},
		{
			statement: "SELECT * FROM schema1.t1 JOIN schema2.t2 ON t1.c1 = t2.c1;",
			expected: []accessResource{
				{
					Database: "SCHEMA1",
					Table:    "T1",
				},
				{
					Database: "SCHEMA2",
					Table:    "T2",
				},
			},
		},
		{
			statement: "SELECT a > (select max(a) from t1) FROM t2;",
			expected: []accessResource{
				{
					Database: "DB",
					Table:    "T1",
				},
				{
					Database: "DB",
					Table:    "T2",
				},
			},
		},
		{
			statement: "SELECT * FROM schema1.t1@remote;",
			linked: []accessResource{
				{
					Database:     "SCHEMA1",
					Table:        "T1",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT * FROM t1@remote;",
			linked: []accessResource{
				{
					Table:        "T1",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT * FROM t1, schema1.t2@remote WHERE EXISTS (SELECT 1 FROM t3@remote);",
			expected: []accessResource{
				{
					Database: "DB",
					Table:    "T1",
				},
			},
			linked: []accessResource{
				{
					Database:     "SCHEMA1",
					Table:        "T2",
					LinkedServer: "REMOTE",
				},
				{
					Table:        "T3",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT * FROM sys.dba_users@remote;",
			linked: []accessResource{
				{
					Database:     "SYS",
					Table:        "DBA_USERS",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT sysdate FROM dual;",
		},
		// Through a link an unqualified DUAL resolves in the link user's schema first.
		{
			statement: "SELECT sysdate FROM dual@remote;",
			linked: []accessResource{
				{
					Table:        "DUAL",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT * FROM t1, sys.dual;",
			expected: []accessResource{
				{
					Database: "DB",
					Table:    "T1",
				},
			},
		},
		{
			statement: "SELECT sysdate FROM sys.dual@remote;",
		},
		// Oracle accepts the public synonym's owner only quoted (ORA-00903 unquoted).
		{
			statement: `SELECT sysdate FROM "PUBLIC".dual;`,
		},
		// A schema-qualified DUAL is a table in that schema, local or remote.
		{
			statement: "SELECT * FROM s.dual@remote;",
			linked: []accessResource{
				{
					Database:     "S",
					Table:        "DUAL",
					LinkedServer: "REMOTE",
				},
			},
		},
		{
			statement: "SELECT * FROM app.dual;",
			expected: []accessResource{
				{
					Database: "APP",
					Table:    "DUAL",
				},
			},
		},
		// Oracle rejects a link after a partition clause (ORA-03048), so it names none.
		{
			statement: "SELECT * FROM t1 PARTITION (p1)@remote;",
			expected: []accessResource{
				{
					Database: "DB",
					Table:    "T1",
				},
			},
		},
		// `t@remote@q` names the link REMOTE@Q.
		{
			statement: "SELECT * FROM s.t@remote@q;",
			linked: []accessResource{
				{
					Database:     "S",
					Table:        "T",
					LinkedServer: "REMOTE@Q",
				},
			},
		},
		{
			statement: "SELECT * FROM t1@remote@q, t1@remote;",
			linked: []accessResource{
				{
					Table:        "T1",
					LinkedServer: "REMOTE@Q",
				},
				{
					Table:        "T1",
					LinkedServer: "REMOTE",
				},
			},
		},
	}

	for _, test := range tests {
		results, err := Parse(test.statement)
		require.NoError(t, err)
		require.NotEmpty(t, results)
		require.NotEmpty(t, results.Items)
		raw, ok := results.Items[0].(*oracleast.RawStmt)
		require.True(t, ok)
		local, linked := collectAccessTables("DB", raw.Stmt)
		require.Equal(t, test.expected, local, test.statement)
		require.Equal(t, test.linked, linked, test.statement)
	}
}

func TestQueryType(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      QueryType
	}{
		{
			name:      "create table",
			statement: "CREATE TABLE T(A NUMBER)",
			want:      DDL,
		},
		{
			name:      "create database admin ddl",
			statement: "CREATE DATABASE",
			want:      DDL,
		},
		{
			name:      "alter database admin ddl",
			statement: "ALTER DATABASE OPEN",
			want:      DDL,
		},
		{
			name:      "drop database admin ddl",
			statement: "DROP DATABASE",
			want:      DDL,
		},
		{
			name:      "lock table is dml",
			statement: "LOCK TABLE T IN EXCLUSIVE MODE",
			want:      DML,
		},
		{
			name:      "transaction remains unknown",
			statement: "COMMIT",
			want:      QueryTypeUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list, err := Parse(test.statement)
			require.NoError(t, err)
			require.Len(t, list.Items, 1)
			raw, ok := list.Items[0].(*oracleast.RawStmt)
			require.True(t, ok)
			require.NotNil(t, raw.Stmt)

			require.Equal(t, test.want, queryType(raw.Stmt, false))
		})
	}
}
