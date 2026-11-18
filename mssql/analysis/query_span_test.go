package analysis

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

type testDatabase struct {
	schemas map[string]*metadata.SchemaMetadata
	order   []string
}

func (d *testDatabase) ListSchemaNames() []string { return d.order }

func (d *testDatabase) GetSchema(name string) Schema {
	for _, schemaName := range d.order {
		if strings.EqualFold(schemaName, name) {
			return testSchema{d.schemas[schemaName]}
		}
	}
	return nil
}

type testSchema struct{ schema *metadata.SchemaMetadata }

func (s testSchema) ListTableNames() []string {
	var names []string
	for _, table := range s.schema.Tables {
		names = append(names, table.Name)
	}
	return names
}

func (s testSchema) GetTable(name string) *metadata.TableMetadata {
	for _, table := range s.schema.Tables {
		if strings.EqualFold(table.Name, name) {
			return table
		}
	}
	return nil
}

func (s testSchema) ListViewNames() []string {
	var names []string
	for _, view := range s.schema.Views {
		names = append(names, view.Name)
	}
	return names
}

func (s testSchema) GetView(name string) *metadata.ViewMetadata {
	for _, view := range s.schema.Views {
		if strings.EqualFold(view.Name, name) {
			return view
		}
	}
	return nil
}

func columns(names ...string) []*metadata.ColumnMetadata {
	var result []*metadata.ColumnMetadata
	for _, name := range names {
		result = append(result, &metadata.ColumnMetadata{Name: name})
	}
	return result
}

func testCatalog() Catalog {
	databases := []*metadata.DatabaseSchemaMetadata{
		{
			Name: "db",
			Schemas: []*metadata.SchemaMetadata{
				{
					Name: "dbo",
					Tables: []*metadata.TableMetadata{
						{Name: "t", Columns: columns("a", "b", "c")},
						{Name: "t1", Columns: columns("a", "b")},
						{Name: "ident_t", Columns: []*metadata.ColumnMetadata{{Name: "id", IsIdentity: true}, {Name: "payload"}}},
						{Name: "bare_edge"},
					},
					Views: []*metadata.ViewMetadata{
						{Name: "v", Definition: "CREATE VIEW v (x, y) AS SELECT a, b + c FROM t"},
						{Name: "v_cycle", Definition: "CREATE VIEW v_cycle AS SELECT a FROM v_cycle"},
						{Name: "v_filtered", Definition: "CREATE VIEW v_filtered AS SELECT a FROM t WHERE b > 0"},
					},
				},
			},
		},
		{
			Name:    "master",
			Schemas: []*metadata.SchemaMetadata{{Name: "dbo", Tables: []*metadata.TableMetadata{{Name: "spt_values", Columns: columns("name")}}}},
		},
	}
	return Catalog{
		ListDatabaseNames: func() ([]string, error) {
			var names []string
			for _, database := range databases {
				names = append(names, database.Name)
			}
			return names, nil
		},
		GetDatabase: func(name string) (Database, error) {
			for _, database := range databases {
				if database.Name != name {
					continue
				}
				d := &testDatabase{schemas: map[string]*metadata.SchemaMetadata{}}
				for _, schema := range database.Schemas {
					d.schemas[schema.Name] = schema
					d.order = append(d.order, schema.Name)
				}
				return d, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
	}
}

func resolve(t *testing.T, statement string) *ResolvedSpan {
	t.Helper()
	span, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
	if err != nil {
		t.Fatalf("ResolveQuerySpan(%q): %v", statement, err)
	}
	return span
}

func column(table, name string) ColumnResource {
	return ColumnResource{Database: "db", Schema: "dbo", Table: table, Column: name}
}

func set(resources ...ColumnResource) ColumnSet {
	result := ColumnSet{}
	for _, resource := range resources {
		result[resource] = true
	}
	return result
}

func equalSet(a, b ColumnSet) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func TestResolveQuerySpanLineage(t *testing.T) {
	tests := []struct {
		statement string
		want      []QuerySpanResult
	}{
		{
			statement: "SELECT a, b + c AS s FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: set(column("t", "a")), IsPlainField: true},
				{Name: "s", SourceColumns: set(column("t", "b"), column("t", "c"))},
			},
		},
		{
			// A view's column list renames its columns. A column read through a
			// view is plain even where the view computes it.
			statement: "SELECT x, y FROM v",
			want: []QuerySpanResult{
				{Name: "x", SourceColumns: set(column("t", "a")), IsPlainField: true},
				{Name: "y", SourceColumns: set(column("t", "b"), column("t", "c")), IsPlainField: true},
			},
		},
		{
			statement: "WITH r AS (SELECT a FROM t UNION ALL SELECT r.a FROM r JOIN t1 ON r.a = t1.b) SELECT a FROM r",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: set(column("t", "a")), IsPlainField: true},
			},
		},
		{
			// The recursive arm reads the CTE by its column list's names.
			statement: "WITH r(n) AS (SELECT a FROM t UNION ALL SELECT n + 1 FROM r WHERE n < 10) SELECT n FROM r",
			want: []QuerySpanResult{
				{Name: "n", SourceColumns: set(column("t", "a")), IsPlainField: true},
			},
		},
		{
			statement: "SELECT PARSE(b AS int) AS p, TRY_PARSE(c AS int USING 'en-US') AS q, JSON_OBJECT('k': a) AS j FROM t",
			want: []QuerySpanResult{
				{Name: "p", SourceColumns: set(column("t", "b"))},
				{Name: "q", SourceColumns: set(column("t", "c"))},
				{Name: "j", SourceColumns: set(column("t", "a"))},
			},
		},
		{
			statement: "SELECT COALESCE((SELECT a FROM t1), 0) AS c, IIF(1 = 1, (SELECT b FROM t1), 0) AS i FROM t",
			want: []QuerySpanResult{
				{Name: "c", SourceColumns: set(column("t1", "a"))},
				{Name: "i", SourceColumns: set(column("t1", "b"))},
			},
		},
		{
			statement: "SELECT $IDENTITY FROM ident_t",
			want: []QuerySpanResult{
				{Name: "id", SourceColumns: set(column("ident_t", "id")), IsPlainField: true},
			},
		},
		{
			statement: "SELECT value FROM STRING_SPLIT('a,b', ',')",
			want: []QuerySpanResult{
				{Name: "value", SourceColumns: set(), IsPlainField: true},
			},
		},
	}
	for _, tc := range tests {
		span := resolve(t, tc.statement)
		if span.Type != Select || len(span.Results) != len(tc.want) {
			t.Fatalf("%q: got %+v", tc.statement, span)
		}
		for i, want := range tc.want {
			got := span.Results[i]
			if got.Name != want.Name || !equalSet(got.SourceColumns, want.SourceColumns) || got.IsPlainField != want.IsPlainField {
				t.Errorf("%q result %d: got %+v, want %+v", tc.statement, i, got, want)
			}
		}
	}
}

func TestResolveQuerySpanPredicates(t *testing.T) {
	span := resolve(t, "SELECT a FROM t WHERE b IN (SELECT a FROM t1 WHERE t1.b > 0) ORDER BY c OFFSET 0 ROWS FETCH NEXT (SELECT COUNT(*) FROM t1) ROWS ONLY")
	want := set(column("t", "b"), column("t1", "a"), column("t1", "b"))
	if !equalSet(span.PredicateColumns, want) {
		t.Errorf("got %v, want %v", span.PredicateColumns, want)
	}
}

func TestResolveQuerySpanPredicatesThroughExpressionsAndViews(t *testing.T) {
	tests := []struct {
		statement string
		want      ColumnSet
	}{
		{"SELECT a FROM t WHERE PARSE(b AS int) > 0", set(column("t", "b"))},
		{"SELECT a FROM t WHERE JSON_OBJECT('k': c) IS NOT NULL", set(column("t", "c"))},
		// A graph pattern's node names are read as the nodes' columns, as
		// LAST_NODE (a) = LAST_NODE (b) was before MATCH had a shape of its own.
		{"SELECT a FROM t WHERE MATCH(LAST_NODE(a) = LAST_NODE(b))", set(column("t", "a"), column("t", "b"))},
		{"SELECT a FROM v_filtered", set(column("t", "b"))},
	}
	for _, tc := range tests {
		if got := resolve(t, tc.statement).PredicateColumns; !equalSet(got, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.statement, got, tc.want)
		}
	}
}

func TestResolveQuerySpanCaseSensitiveAliases(t *testing.T) {
	for _, statement := range []string{
		"SELECT D.X FROM (SELECT a AS X FROM t) AS D",
		"WITH C(X) AS (SELECT a FROM t) SELECT C.X FROM C",
	} {
		span, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
		if err != nil {
			t.Fatalf("%q: %v", statement, err)
		}
		if len(span.Results) != 1 || span.Results[0].Name != "X" || !equalSet(span.Results[0].SourceColumns, set(column("t", "a"))) {
			t.Errorf("%q: got %+v", statement, span.Results)
		}
	}
}

func TestResolveQuerySpanAccessTables(t *testing.T) {
	table := func(database, name string) ColumnResource {
		return ColumnResource{Database: database, Schema: "dbo", Table: name}
	}
	tests := []struct {
		statement string
		want      ColumnSet
	}{
		{"SELECT COALESCE((SELECT a FROM t1), 0)", set(table("db", "t1"))},
		{"SELECT NULLIF(a, (SELECT b FROM t1)) FROM t", set(table("db", "t"), table("db", "t1"))},
		{"SELECT TOP ((SELECT COUNT(*) FROM t1)) a FROM t", set(table("db", "t"), table("db", "t1"))},
		// A CTE reference names no physical table, in every arm of a set
		// operation, while a schema-qualified name still does.
		{"WITH x AS (SELECT a FROM t) SELECT a FROM x UNION SELECT a FROM x UNION SELECT a FROM x", set(table("db", "t"))},
		{"WITH x AS (SELECT a FROM t) SELECT x.a FROM x JOIN dbo.x AS y ON x.a = y.a", set(table("db", "t"), table("db", "x"))},
		{"WITH x AS (SELECT a FROM t1) DELETE FROM t WHERE a IN (SELECT a FROM x)", set(table("db", "t"), table("db", "t1"))},
		{"UPDATE t SET a = (SELECT a FROM t1)", set(table("db", "t"), table("db", "t1"))},
		{"INSERT INTO t (a) VALUES ((SELECT a FROM t1))", set(table("db", "t"), table("db", "t1"))},
	}
	for _, tc := range tests {
		span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
		if err != nil {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		if !equalSet(span.SourceColumns, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.statement, span.SourceColumns, tc.want)
		}
	}

	span := resolve(t, "WITH x AS (SELECT name FROM master.dbo.spt_values) SELECT * FROM x")
	if span.Type != SelectInfoSchema || !equalSet(span.SourceColumns, set(table("master", "spt_values"))) {
		t.Errorf("system CTE: got %+v", span)
	}
	for _, statement := range []string{
		"SELECT COALESCE((SELECT name FROM master.dbo.spt_values), a) FROM t",
		"SELECT IIF(a > 0, (SELECT name FROM master.dbo.spt_values), '') FROM t",
		"SELECT TOP ((SELECT COUNT(*) FROM master.dbo.spt_values)) a FROM t",
	} {
		if _, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()}); err != ErrMixUserSystemTables {
			t.Errorf("%q: got %v, want ErrMixUserSystemTables", statement, err)
		}
	}
}

func TestResolveQuerySpanNotFound(t *testing.T) {
	span := resolve(t, "SELECT a FROM missing")
	if span.NotFoundError == nil || *span.NotFoundError.Table != "missing" || len(span.Results) != 0 {
		t.Fatalf("got %+v", span)
	}
	if got := span.NotFoundError.Error(); got != "resource not found: database: db, schema: dbo, table: missing" {
		t.Errorf("got %q", got)
	}
	for _, statement := range []string{
		"SELECT missing FROM t",
		"SELECT a FROM t WHERE t.missing > 0",
		"SELECT a FROM t WHERE b IN (SELECT x FROM t1)",
	} {
		_, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
		var notFound *NotFoundError
		if err == nil || errors.As(err, &notFound) {
			t.Errorf("%q: got %v, want a hard error", statement, err)
		}
	}
}

func TestResolveQuerySpanErrors(t *testing.T) {
	_, err := ResolveQuerySpan("SELECT * FROM t, master.dbo.spt_values", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
	if err != ErrMixUserSystemTables {
		t.Errorf("mixed tables: got %v", err)
	}
	_, err = ResolveQuerySpan("SELECT * FROM v_cycle", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
	if err == nil || !strings.Contains(err.Error(), "cyclic view reference") {
		t.Errorf("cycle: got %v", err)
	}
	_, err = ResolveQuerySpan("SELECT * FROM t PIVOT (SUM(a) FOR b IN ([1])) AS p", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
	var typeNotSupported *TypeNotSupportedError
	if !errors.As(err, &typeNotSupported) || typeNotSupported.Type != "*ast.PivotExpr" {
		t.Errorf("pivot: got %v", err)
	}
	_, err = newQuerySpanExtractor("db", "", testCatalog(), nil, false).resolveExpressionNode(&ast.WindowDef{})
	if !errors.As(err, &typeNotSupported) || typeNotSupported.Type != "*ast.WindowDef" {
		t.Errorf("unmodeled expression: got %v", err)
	}
	_, err = ResolveQuerySpan("SELECT $edge_id FROM bare_edge", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
	if err == nil || !strings.Contains(err.Error(), "cannot derive lineage") {
		t.Errorf("graph pseudo-column: got %v", err)
	}
	catalog := testCatalog()
	catalog.ListDatabaseNames = func() ([]string, error) { return nil, errors.New("boom") }
	_, err = ResolveQuerySpan("SELECT a FROM t", ResolveOptions{Catalog: catalog})
	if err == nil || err.Error() != "failed to list databases: boom" {
		t.Errorf("list: got %v", err)
	}
	_, err = ResolveQuerySpan("SELECT FROM", ResolveOptions{Catalog: catalog})
	if err == nil {
		t.Error("syntax error: got nil")
	}
}

func TestResolveQuerySpanNonQueries(t *testing.T) {
	tests := []struct {
		statement string
		want      QueryType
	}{
		{"INSERT INTO t SELECT a, b, c FROM t1", DML},
		{"EXEC sp_who", DML},
		{"CREATE TABLE t2 (a INT)", DDL},
		{"SELECT a INTO t2 FROM t", DDL},
		{"SET NOCOUNT ON", Select},
		{"SELECT * FROM master.dbo.spt_values", SelectInfoSchema},
		{"BEGIN TRAN", QueryTypeUnknown},
	}
	for _, tc := range tests {
		if got := resolve(t, tc.statement).Type; got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.statement, got, tc.want)
		}
	}
}

func TestResolveQuerySpanTempTables(t *testing.T) {
	tempTables := map[string]*TempTable{}
	opts := ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog(), TempTables: tempTables}
	for _, statement := range []string{
		"DECLARE @v TABLE (id INT, name NVARCHAR(10))",
		"CREATE TABLE #c (x INT)",
		"SELECT a, b + c INTO #s FROM t",
	} {
		if _, err := ResolveQuerySpan(statement, opts); err != nil {
			t.Fatalf("%q: %v", statement, err)
		}
	}
	for name, want := range map[string][]string{"@v": {"id", "name"}, "#c": {"x"}, "#s": {"a", "b+c"}} {
		if got := tempTables[name]; got == nil || strings.Join(got.Columns, ",") != strings.Join(want, ",") {
			t.Errorf("%s: got %+v, want columns %v", name, got, want)
		}
	}
	span, err := ResolveQuerySpan("SELECT id FROM @V", opts)
	if err != nil || len(span.Results) != 1 || !span.Results[0].SourceColumns[ColumnResource{Table: "@v", Column: "id"}] {
		t.Errorf("table variable: got %+v, %v", span, err)
	}
}

func TestResolveQuerySpanPlanOutput(t *testing.T) {
	output := map[string]bool{}
	opts := ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(), PlanOutput: output}
	steps := []struct {
		statement   string
		wantType    QueryType
		wantAnalyze bool
	}{
		{"SELECT a FROM t", Select, false},
		{"SET SHOWPLAN_XML ON", Explain, false},
		{"SELECT a FROM t", Explain, false},
		{"INSERT INTO t1 VALUES (1, 2)", DML, false},
		{"SET SHOWPLAN_XML OFF", Explain, false},
		{"SET STATISTICS XML ON", Select, true},
		{"INSERT INTO t1 VALUES (1, 2)", DML, true},
	}
	for _, step := range steps {
		span, err := ResolveQuerySpan(step.statement, opts)
		if err != nil {
			t.Fatalf("%q: %v", step.statement, err)
		}
		if span.Type != step.wantType || span.ExplainAnalyze != step.wantAnalyze {
			t.Errorf("%q: got type %v analyze %v, want %v %v", step.statement, span.Type, span.ExplainAnalyze, step.wantType, step.wantAnalyze)
		}
	}
	if output["SHOWPLAN_XML"] || !output["STATISTICS XML"] {
		t.Errorf("plan output: got %v", output)
	}
}

func TestClassifyQueryTypeSelectInto(t *testing.T) {
	tests := []struct {
		statement string
		want      QueryType
	}{
		// SELECT ... INTO creates a table, so it is a write, including when
		// INTO sits on the first arm of a set operation.
		{"SELECT a INTO dbo.new_t FROM dbo.src", DDL},
		{"SELECT a INTO new_t FROM x UNION SELECT b FROM y", DDL},
		// INTO a session-scoped temp table stays a read.
		{"SELECT a INTO #t FROM dbo.src", Select},
		{"SELECT a INTO ##g FROM dbo.src", Select},
		{"SELECT a FROM dbo.src", Select},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		if got := classifyQueryType(stmts[0].AST, false); got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.statement, got, tc.want)
		}
	}
}

func TestSystemObjects(t *testing.T) {
	if !isSystemResource(ColumnResource{Database: "MASTER"}, true) || isSystemResource(ColumnResource{Database: "MASTER"}, false) {
		t.Error("system database case rules")
	}
	if !isSystemResource(ColumnResource{Schema: "sys"}, false) || isSystemResource(ColumnResource{Schema: "SYS"}, false) {
		t.Error("system schema case rules")
	}
}
