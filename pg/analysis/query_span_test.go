package analysis

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/catalog"
	"github.com/stretchr/testify/require"
)

// protoDatabase answers Database lookups from a metadata proto, matching
// names exactly as Bytebase does for PostgreSQL.
type protoDatabase struct {
	metadata *metadatapb.DatabaseSchemaMetadata
}

func (d *protoDatabase) Metadata() *metadatapb.DatabaseSchemaMetadata { return d.metadata }

func (d *protoDatabase) schema(name string) *metadatapb.SchemaMetadata {
	for _, schema := range d.metadata.GetSchemas() {
		if schema.GetName() == name {
			return schema
		}
	}
	return nil
}

func (d *protoDatabase) GetTable(schema, table string) *metadatapb.TableMetadata {
	for _, t := range d.schema(schema).GetTables() {
		if t.GetName() == table || hasPartition(t.GetPartitions(), table) {
			return t
		}
	}
	return nil
}

func hasPartition(partitions []*metadatapb.TablePartitionMetadata, name string) bool {
	for _, partition := range partitions {
		if partition.GetName() == name || hasPartition(partition.GetSubpartitions(), name) {
			return true
		}
	}
	return false
}

func (d *protoDatabase) SearchRelation(searchPath []string, name string) (string, string) {
	for _, schemaName := range searchPath {
		schema := d.schema(schemaName)
		if schema == nil {
			continue
		}
		if d.GetTable(schemaName, name) != nil {
			return schemaName, name
		}
		for _, view := range schema.GetViews() {
			if view.GetName() == name {
				return schemaName, name
			}
		}
		for _, view := range schema.GetMaterializedViews() {
			if view.GetName() == name {
				return schemaName, name
			}
		}
		for _, table := range schema.GetExternalTables() {
			if table.GetName() == name {
				return schemaName, name
			}
		}
		for _, sequence := range schema.GetSequences() {
			if sequence.GetName() == name {
				return schemaName, name
			}
		}
	}
	return "", ""
}

func (d *protoDatabase) SearchFunctions(searchPath []string, name string) []*metadatapb.FunctionMetadata {
	var functions []*metadatapb.FunctionMetadata
	for _, schemaName := range searchPath {
		for _, function := range d.schema(schemaName).GetFunctions() {
			if strings.EqualFold(function.GetName(), name) {
				functions = append(functions, function)
			}
		}
	}
	return functions
}

// testCatalog serves the given databases; any other name fails the lookup.
func testCatalog(databases ...*metadatapb.DatabaseSchemaMetadata) Catalog {
	return Catalog{GetDatabase: func(name string) (Database, error) {
		for _, database := range databases {
			if database.GetName() == name {
				return &protoDatabase{metadata: database}, nil
			}
		}
		return nil, errors.New("database not found")
	}}
}

func spanFor(t *testing.T, statement string, metadata *metadatapb.DatabaseSchemaMetadata) *ResolvedSpan {
	t.Helper()
	span, err := ResolveQuerySpan(context.Background(), statement, ResolveOptions{
		DefaultDatabase: metadata.GetName(),
		Catalog:         testCatalog(metadata),
	})
	require.NoError(t, err)
	return span
}

func column(name, typ string) *metadatapb.ColumnMetadata {
	return &metadatapb.ColumnMetadata{Name: name, Type: typ}
}

func healthySchema() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("id", "int4"), column("email", "text"), column("ssn", "text")}},
				{Name: "o", Columns: []*metadatapb.ColumnMetadata{column("id", "int4"), column("amt", "numeric")}},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT id, email FROM public.t",
				Columns:    []*metadatapb.ColumnMetadata{column("id", "int4"), column("email", "text")},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv",
				Definition: "SELECT id, ssn FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{
					{Schema: "public", Table: "t", Column: "id"},
					{Schema: "public", Table: "t", Column: "ssn"},
				},
			}},
		}},
	}
}

// A table with no columns models a snapshot from a sync before #20581.
func degradedSchema() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "t"}},
		}},
	}
}

func degradedViewSchema() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "t"}},
			Views:  []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT id, email FROM public.t"}},
		}},
	}
}

func partialSchema() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("email", "text"), column("ssn", "text")}}},
		}},
	}
}

func columnNames(columns ColumnSet) []string {
	names := []string{}
	for column := range columns {
		names = append(names, column.String())
	}
	slices.Sort(names)
	return names
}

func TestResolveQuerySpanLineage(t *testing.T) {
	span := spanFor(t, "SELECT t.id, upper(t.email) AS e, o.amt FROM t JOIN o ON o.id = t.id WHERE t.ssn = '1'", healthySchema())
	require.Equal(t, Select, span.Type)
	require.False(t, span.ExplainAnalyze)
	require.Equal(t, []string{"db.public.o", "db.public.t"}, columnNames(span.SourceColumns))
	require.Equal(t, []string{"db.public.o.id", "db.public.t.id", "db.public.t.ssn"}, columnNames(span.PredicateColumns))
	require.Nil(t, span.UnresolvedColumnsError)
	require.Len(t, span.Results, 3)
	want := []struct {
		name    string
		sources []string
		plain   bool
	}{
		{"id", []string{"db.public.t.id"}, false},
		{"e", []string{"db.public.t.email"}, false},
		{"amt", []string{"db.public.o.amt"}, false},
	}
	for i, w := range want {
		require.Equal(t, w.name, span.Results[i].Name)
		require.Equal(t, w.sources, columnNames(span.Results[i].SourceColumns))
		require.Equal(t, w.plain, span.Results[i].IsPlainField, w.name)
	}
}

// Only the columns a star expands are plain fields.
func TestResolveQuerySpanPlainFields(t *testing.T) {
	span := spanFor(t, "SELECT *, id AS x FROM t", healthySchema())
	require.Len(t, span.Results, 4)
	for i, plain := range []bool{true, true, true, false} {
		require.Equal(t, plain, span.Results[i].IsPlainField, span.Results[i].Name)
	}
}

// A view's columns are its own sources; lineage stops at the view.
func TestResolveQuerySpanThroughViews(t *testing.T) {
	span := spanFor(t, "SELECT email FROM v", healthySchema())
	require.Equal(t, []string{"db.public.v"}, columnNames(span.SourceColumns))
	require.Len(t, span.Results, 1)
	require.Equal(t, []string{"db.public.v.email"}, columnNames(span.Results[0].SourceColumns))

	span = spanFor(t, "SELECT ssn FROM mv", healthySchema())
	require.Len(t, span.Results, 1)
	require.Equal(t, []string{"db.public.mv.ssn"}, columnNames(span.Results[0].SourceColumns))
}

func TestResolveQuerySpanStatements(t *testing.T) {
	testCases := []struct {
		statement string
		want      QueryType
		analyze   bool
		sources   []string
	}{
		{"INSERT INTO t (id) SELECT id FROM o", DML, false, []string{"db.public.o", "db.public.t"}},
		{"CREATE TABLE n (id int)", DDL, false, []string{}},
		{"EXPLAIN SELECT * FROM t", Explain, false, []string{"db.public.t"}},
		{"EXPLAIN ANALYZE DELETE FROM t", DML, true, []string{"db.public.t"}},
		{"SELECT * FROM pg_catalog.pg_class", SelectInfoSchema, false, []string{"db.pg_catalog.pg_class"}},
		{"SELECT * FROM pg_database", SelectInfoSchema, false, []string{"db.pg_database"}},
		{"SHOW search_path", SelectInfoSchema, false, []string{}},
		{"SET search_path TO public", Select, false, []string{}},
		{"SELECT * FROM missing", Select, false, []string{}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, healthySchema())
			require.Equal(t, tc.want, span.Type)
			require.Equal(t, tc.analyze, span.ExplainAnalyze)
			require.Equal(t, tc.sources, columnNames(span.SourceColumns))
		})
	}
}

func TestResolveQuerySpanSearchPath(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("a", "int4")}}}},
			{Name: "app", Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("b", "int4")}}}},
		},
	}
	span, err := ResolveQuerySpan(context.Background(), "SELECT * FROM t", ResolveOptions{
		DefaultDatabase: "db",
		SearchPath:      SearchPathForSelectedSchema("app", []string{"public"}),
		Catalog:         testCatalog(metadata),
	})
	require.NoError(t, err)
	require.Equal(t, []string{"db.app.t"}, columnNames(span.SourceColumns))
	require.Len(t, span.Results, 1)
	require.Equal(t, "b", span.Results[0].Name)

	require.Equal(t, []string{"public"}, SearchPathForSelectedSchema("", []string{"public"}))
	require.Equal(t, []string{"public"}, SearchPathForSelectedSchema("public", []string{"x"}))
	require.Equal(t, []string{"app", "public"}, SearchPathForSelectedSchema("app", []string{"x"}))
}

func TestResolveQuerySpanErrors(t *testing.T) {
	_, err := ResolveQuerySpan(context.Background(), "SELECT FROM WHERE", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(healthySchema())})
	var syntaxErr *SyntaxError
	require.ErrorAs(t, err, &syntaxErr)

	_, err = ResolveQuerySpan(context.Background(), "SELECT 1; SELECT 2", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(healthySchema())})
	require.EqualError(t, err, "expected exactly 1 statement, got 2")

	_, err = ResolveQuerySpan(context.Background(), "SELECT * FROM t, pg_catalog.pg_class", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(healthySchema())})
	require.ErrorIs(t, err, ErrMixUserSystemTables)

	_, err = ResolveQuerySpan(context.Background(), "SELECT 1", ResolveOptions{DefaultDatabase: "other", Catalog: testCatalog(healthySchema())})
	require.EqualError(t, err, "failed to initialize catalog: failed to get database metadata for catalog init: failed to get database metadata for database: other: database not found")
}

func TestResolveQuerySpanUnsyncedDatabase(t *testing.T) {
	span, err := ResolveQuerySpan(context.Background(), "SELECT 1 AS x", ResolveOptions{
		DefaultDatabase: "db",
		Catalog:         Catalog{GetDatabase: func(string) (Database, error) { return nil, nil }},
	})
	require.NoError(t, err)
	require.Len(t, span.Results, 1)
	require.Equal(t, "x", span.Results[0].Name)
}

func TestResolveQuerySpanUnresolvedColumns(t *testing.T) {
	span := spanFor(t, "SELECT * FROM public.t", degradedSchema())
	require.NotNil(t, span.UnresolvedColumnsError)
	require.Equal(t, []ColumnResource{{Database: "db", Schema: "public", Table: "t"}}, span.UnresolvedColumnsError.Relations)
	require.EqualError(t, span.UnresolvedColumnsError, "the synced schema describes no columns for db.public.t")
}

func TestResolveQuerySpanThroughAPartition(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:       "orders",
				Columns:    []*metadatapb.ColumnMetadata{column("id", "int4"), column("ssn", "text")},
				Partitions: []*metadatapb.TablePartitionMetadata{{Name: "orders_2024"}},
			}},
		}},
	}
	span := spanFor(t, "SELECT ssn FROM orders_2024", metadata)
	require.Equal(t, []string{"db.public.orders_2024"}, columnNames(span.SourceColumns))
	require.Len(t, span.Results, 1)
}

// TestClassifyQueryTypeExplain pins how EXPLAIN is classified, which decides the
// permission an explain request needs and whether the driver will run it. The
// security-critical case is EXPLAIN ANALYZE of a write: PostgreSQL executes the
// wrapped statement, so it must classify as the write it performs, not as a
// plain EXPLAIN.
func TestClassifyQueryTypeExplain(t *testing.T) {
	testCases := []struct {
		statement string
		wantType  QueryType
		wantAlyze bool
	}{
		// Plain EXPLAIN only plans; it never executes, whatever it wraps.
		{"EXPLAIN SELECT * FROM t", Explain, false},
		{"EXPLAIN DELETE FROM t", Explain, false},
		{"EXPLAIN INSERT INTO t VALUES (1)", Explain, false},

		// EXPLAIN ANALYZE executes the wrapped statement, so it inherits its type.
		{"EXPLAIN ANALYZE SELECT * FROM t", Select, true},
		{"EXPLAIN ANALYZE DELETE FROM t", DML, true},
		{"EXPLAIN ANALYZE INSERT INTO t VALUES (1)", DML, true},
		{"EXPLAIN ANALYZE UPDATE t SET a = 1", DML, true},

		// Bare statements keep classifying as before.
		{"SELECT * FROM t", Select, false},
		{"DELETE FROM t", DML, false},
	}

	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			stmts, err := omnipg.Parse(tc.statement)
			require.NoError(t, err)
			require.Len(t, stmts, 1)

			gotType, gotAnalyze := classifyQueryType(stmts[0].AST, false /* allSystems */)
			require.Equal(t, tc.wantType, gotType)
			require.Equal(t, tc.wantAlyze, gotAnalyze)
		})
	}
}

func TestPolarDBSystemSchema(t *testing.T) {
	require.True(t, IsSystemSchema("pg_bitmapindex"))
	require.Contains(t, strings.Split(SystemSchemaWhereClause, ","), "'pg_bitmapindex'")
}

func TestBuiltinFunctions(t *testing.T) {
	functions := BuiltinFunctions()
	require.True(t, slices.IsSorted(functions))
	require.Contains(t, functions, "jsonb_build_object")
	functions[0] = "changed"
	require.NotEqual(t, "changed", BuiltinFunctions()[0], "callers get a copy")
}

func TestAppendQueryDoesNotMutateSharedStack(t *testing.T) {
	outer := &catalog.Query{}
	parent := &catalog.Query{}
	current := &catalog.Query{}
	nested := &catalog.Query{}
	queryStack := []*catalog.Query{outer, parent, current}

	nextStack := appendQuery(queryStack[:2], nested)

	if queryStack[2] != current {
		t.Fatalf("appendQuery mutated caller stack: got %p, want %p", queryStack[2], current)
	}
	if len(nextStack) != 3 || nextStack[2] != nested {
		t.Fatalf("appendQuery returned unexpected stack: %+v", nextStack)
	}
}

func TestWalkExprStopsOnSelfReferentialFunction(t *testing.T) {
	mustNotOverflow(t, func() {
		varExpr := &catalog.VarExpr{RangeIdx: 0, AttNum: 1}
		q := &catalog.Query{
			TargetList: []*catalog.TargetEntry{{
				Expr: varExpr,
			}},
			RangeTable: []*catalog.RangeTableEntry{{
				Kind: catalog.RTEFunction,
				FuncExprs: []catalog.AnalyzedExpr{&catalog.FuncCallExpr{
					Args: []catalog.AnalyzedExpr{&catalog.OpExpr{
						Left: varExpr,
						Right: &catalog.ConstExpr{
							Value: "items",
						},
					}},
				}},
			}},
		}
		extractor := newQuerySpanExtractor(testDB, []string{testSchema}, Catalog{})
		extractor.cat = catalog.New()
		extractor.walkExpr(q, q.TargetList[0].Expr, make(ColumnSet))
	})
}

func TestIsUltimatelyPlainColumnStopsOnCyclicSubquery(t *testing.T) {
	mustNotOverflow(t, func() {
		varExpr := &catalog.VarExpr{RangeIdx: 0, AttNum: 1}
		q := &catalog.Query{
			TargetList: []*catalog.TargetEntry{{
				Expr: varExpr,
			}},
			RangeTable: []*catalog.RangeTableEntry{{
				Kind: catalog.RTESubquery,
			}},
		}
		q.RangeTable[0].Subquery = q
		_ = isUltimatelyPlainColumn(q, varExpr)
	})
}

func TestCollectPredicateColumnsStopsOnCyclicQuery(t *testing.T) {
	mustNotOverflow(t, func() {
		q := &catalog.Query{}
		q.CTEList = []*catalog.CommonTableExprQ{{Query: q}}
		extractor := newQuerySpanExtractor(testDB, []string{testSchema}, Catalog{})
		extractor.collectPredicateColumns(q)
	})
}

func TestExtractLineageStopsOnCyclicSetOp(t *testing.T) {
	mustNotOverflow(t, func() {
		q := &catalog.Query{SetOp: catalog.SetOpUnion}
		q.LArg = q
		q.RArg = &catalog.Query{TargetList: []*catalog.TargetEntry{{
			Expr: &catalog.ConstExpr{
				Value: "1",
			},
		}}}
		selStmt := &ast.SelectStmt{}
		selStmt.Larg = selStmt
		selStmt.Rarg = &ast.SelectStmt{}
		extractor := newQuerySpanExtractor(testDB, []string{testSchema}, Catalog{})
		extractor.extractLineage(q, selStmt)
	})
}

func TestResolveThroughCTEStopsOnCyclicCTE(t *testing.T) {
	mustNotOverflow(t, func() {
		cteSel := &ast.SelectStmt{
			TargetList: &ast.List{Items: []ast.Node{&ast.ResTarget{
				Name: "x",
				Val: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{
					&ast.String{Str: "c"},
					&ast.String{Str: "x"},
				}}},
			}}},
			FromClause: &ast.List{Items: []ast.Node{&ast.RangeVar{Relname: "c"}}},
		}
		extractor := newQuerySpanExtractor(testDB, []string{testSchema}, Catalog{})
		extractor.cat = catalog.New()
		analyzer := &plpgsqlAnalyzer{
			extractor: extractor,
			scope:     newVariableScope(nil),
			cteMap:    map[string]*ast.SelectStmt{"c": cteSel},
		}
		analyzer.resolveThroughCTE(cteSel, "x", nil, make(ColumnSet))
	})
}

// ---------- helpers ----------

// testDB / testSchema are the fixed database and schema names these tests use;
// as constants they keep the call sites terse and satisfy the unparam linter.
const (
	testDB     = "db"
	testSchema = "public"
)

// noOverflowEnv marks the subprocess mustNotOverflow spawns.
const noOverflowEnv = "OMNI_TEST_NO_OVERFLOW"

// mustNotOverflow runs fn in a subprocess with a 1 MB stack, so recursion that
// does not terminate fails this test instead of taking the whole run down. The
// subprocess re-runs this same test, which then takes the fn branch.
func mustNotOverflow(t *testing.T, fn func()) {
	t.Helper()
	if os.Getenv(noOverflowEnv) == "1" {
		debug.SetMaxStack(1 << 20)
		fn()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$")
	cmd.Env = append(os.Environ(), noOverflowEnv+"=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("recursion did not terminate: %v\n%s", err, output)
	}
}

// TestRangeVarFallbackSkipsCompositeTypes proves the query-span fallback does
// not expand a standalone composite type as a FROM source — PostgreSQL
// rejects composite types as table sources.
func TestRangeVarFallbackSkipsCompositeTypes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{
						Name: "addr",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "street", Type: "text"},
						},
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "t",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
						},
					},
				},
			},
		},
	}

	cat := catalog.New()
	cat.SetSearchPath([]string{"public"})
	if _, err := cat.LoadMetadata(context.Background(), meta, catalog.LoadMetadataOptions{}); err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}

	e := &querySpanExtractor{
		cat:             cat,
		searchPath:      []string{"public"},
		defaultDatabase: "db",
	}
	if results := e.extractColumnsFromRangeVar(&ast.RangeVar{Schemaname: "public", Relname: "addr"}); results != nil {
		t.Errorf("composite type must not expand as a FROM source, got %d columns", len(results))
	}
	if results := e.extractColumnsFromRangeVar(&ast.RangeVar{Schemaname: "public", Relname: "t"}); len(results) != 1 {
		t.Errorf("real table must still expand, got %d columns", len(results))
	}
}
