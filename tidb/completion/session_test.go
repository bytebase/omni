package completion

import (
	"context"
	"fmt"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/catalog"
	"github.com/stretchr/testify/require"
)

func sessionMetaLoader(meta *metadatapb.DatabaseSchemaMetadata, caseSensitive bool) func(context.Context, string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
	return func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
		return db, meta, caseSensitive, true
	}
}

func TestCompleteSession_ReservedWordIdentifiersSurface(t *testing.T) {
	cfg := SessionConfig{
		DefaultDatabase: "testdb",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "testdb",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{{
					Name: "order",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "select"},
						{Name: "key"},
					},
				}},
			}},
		}, true),
	}

	got := CompleteSession(context.Background(), cfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(got, "`order`", CandidateTable), "got %v", got)

	got2 := CompleteSession(context.Background(), cfg, "SELECT  FROM `order`", 1, len("SELECT "))
	require.True(t, containsCandidate(got2, "`select`", CandidateColumn), "got %v", got2)
	require.True(t, containsCandidate(got2, "`key`", CandidateColumn), "got %v", got2)
}

func TestCompleteSession_ReadOnlyFiltersWriteKeywordsAndTiDBBatch(t *testing.T) {
	loader := sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
		Name:    "testdb",
		Schemas: []*metadatapb.SchemaMetadata{{}},
	}, true)
	const stmt = "SELECT 1; "
	all := CompleteSession(context.Background(), SessionConfig{DefaultDatabase: "testdb", LoadDatabase: loader}, stmt, 1, len(stmt))
	ro := CompleteSession(context.Background(), SessionConfig{ReadOnly: true, DefaultDatabase: "testdb", LoadDatabase: loader}, stmt, 1, len(stmt))

	for _, kw := range []string{"SELECT", "WITH", "SHOW"} {
		require.True(t, containsCandidate(all, kw, CandidateKeyword), "%s in ALL", kw)
		require.True(t, containsCandidate(ro, kw, CandidateKeyword), "%s in ReadOnly", kw)
	}
	for _, kw := range []string{"INSERT", "BATCH"} {
		require.True(t, containsCandidate(all, kw, CandidateKeyword), "%s in ALL", kw)
		require.False(t, containsCandidate(ro, kw, CandidateKeyword), "%s must be dropped in ReadOnly", kw)
	}

	roCfg := SessionConfig{ReadOnly: true, DefaultDatabase: "testdb", LoadDatabase: loader}
	gotShow := CompleteSession(context.Background(), roCfg, "SHOW ", 1, len("SHOW "))
	require.True(t, containsCandidate(gotShow, "CREATE", CandidateKeyword))

	gotFor := CompleteSession(context.Background(), roCfg, "SELECT 1 FOR ", 1, len("SELECT 1 FOR "))
	require.True(t, containsCandidate(gotFor, "UPDATE", CandidateKeyword))

	gotCre := CompleteSession(context.Background(), roCfg, "CRE", 1, len("CRE"))
	require.False(t, containsCandidate(gotCre, "CREATE", CandidateKeyword))

	gotIns := CompleteSession(context.Background(), roCfg, "SELECT 1; INS", 1, len("SELECT 1; INS"))
	require.False(t, containsCandidate(gotIns, "INSERT", CandidateKeyword))

	gotShowCre := CompleteSession(context.Background(), roCfg, "SHOW CRE", 1, len("SHOW CRE"))
	require.True(t, containsCandidate(gotShowCre, "CREATE", CandidateKeyword))
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	cases := []struct {
		name            string
		caretInBacktick bool
		want            string
	}{
		{"users", false, "users"},
		{"t1", false, "t1"},
		{"comment", false, "comment"},
		{"order", false, "`order`"},
		{"select", false, "`select`"},
		{"Order", false, "`Order`"},
		{"foo-bar", false, "`foo-bar`"},
		{"1col", false, "`1col`"},
		{"a`b", false, "`a``b`"},
		{"a`b", true, "a``b"},
		{"order", true, "order"},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, QuoteIdentifierIfNeeded(tc.name, tc.caretInBacktick),
			"QuoteIdentifierIfNeeded(%q, %v)", tc.name, tc.caretInBacktick)
	}
}

func TestCompleteSession_QualifiedTableAcrossDatabases(t *testing.T) {
	appMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "appdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}}}},
		}},
	}
	otherMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "otherdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "c2"}}}},
		}},
	}
	cfg := SessionConfig{
		DefaultDatabase: "appdb",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			if db == "otherdb" {
				return db, otherMeta, true, true
			}
			return db, appMeta, true, true
		},
	}

	got := CompleteSession(context.Background(), cfg, "SELECT * FROM otherdb.", 1, len("SELECT * FROM otherdb."))
	require.True(t, containsCandidate(got, "t2", CandidateTable), "got %v", got)
	require.False(t, containsCandidate(got, "t1", CandidateTable), "got %v", got)

	gotFuncParen := CompleteSession(context.Background(), cfg, "SELECT CO(c1) FROM t1", 1, len("SELECT CO"))
	require.True(t, containsCandidate(gotFuncParen, "COUNT", CandidateFunction), "got %v", gotFuncParen)
	require.False(t, containsCandidate(gotFuncParen, "COUNT()", CandidateFunction), "got %v", gotFuncParen)

	gotFuncQuote := postProcessSessionCandidates("SELECT ", len("SELECT "), false, []Candidate{
		{Type: CandidateFunction, Text: "LEFT"},
		{Type: CandidateFunction, Text: "my-func"},
		{Type: CandidateFunction, Text: "f()"},
	}, func(string, int) bool { return true })
	require.True(t, containsCandidate(gotFuncQuote, "LEFT()", CandidateFunction), "got %v", gotFuncQuote)
	require.False(t, containsCandidate(gotFuncQuote, "`LEFT`()", CandidateFunction), "got %v", gotFuncQuote)
	require.True(t, containsCandidate(gotFuncQuote, "`my-func`()", CandidateFunction), "got %v", gotFuncQuote)
	require.True(t, containsCandidate(gotFuncQuote, "`f()`()", CandidateFunction), "got %v", gotFuncQuote)

	gotRoutineQuote := postProcessSessionCandidates("DROP FUNCTION ", len("DROP FUNCTION "), false, []Candidate{
		{Type: CandidateFunction, Text: "my-func"},
		{Type: CandidateFunction, Text: "f()"},
	}, func(string, int) bool { return false })
	require.True(t, containsCandidate(gotRoutineQuote, "`my-func`", CandidateFunction), "got %v", gotRoutineQuote)
	require.True(t, containsCandidate(gotRoutineQuote, "`f()`", CandidateFunction), "got %v", gotRoutineQuote)
	require.False(t, containsCandidate(gotRoutineQuote, "`my-func`()", CandidateFunction), "got %v", gotRoutineQuote)
	require.False(t, containsCandidate(gotRoutineQuote, "f()", CandidateFunction), "got %v", gotRoutineQuote)
}

func TestCaretInsideBacktickIdentifier(t *testing.T) {
	cases := []struct {
		statement string
		pos       int
		want      bool
	}{
		{"SELECT * FROM ", 14, false},
		{"SELECT * FROM `o", 16, true},
		{"SELECT `a`, ", 12, false},
		{"`", 1, true},
		{"ab", 99, false},
		{"SELECT '`' FROM ", 16, false},
		{"SELECT /* `x` */ FROM ", 22, false},
		{"SELECT 1 -- `\nFROM ", 19, false},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, CaretInsideBacktickIdentifier(tc.statement, tc.pos),
			"CaretInsideBacktickIdentifier(%q, %d)", tc.statement, tc.pos)
	}
}

func TestStatementReferencesDatabase(t *testing.T) {
	cases := []struct {
		statement string
		db        string
		want      bool
	}{
		{"SELECT * FROM otherdb.t", "otherdb", true},
		{"SELECT * FROM `otherdb`.t", "otherdb", true},
		{"SELECT OTHERDB.t FROM x", "otherdb", true},
		{"SELECT 'otherdb.' FROM t", "otherdb", false},
		{"SELECT /* otherdb. */ * FROM t", "otherdb", false},
		{"SELECT 1 -- otherdb.\n FROM t", "otherdb", false},
		{"SELECT nototherdb.t FROM x", "otherdb", false},
		{"SELECT a FROM t", "otherdb", false},
		{"SELECT * FROM otherdb . t2", "otherdb", true},
		{"SELECT * FROM `otherdb` . t2", "otherdb", true},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, StatementReferencesDatabase(tc.statement, tc.db),
			"StatementReferencesDatabase(%q, %q)", tc.statement, tc.db)
	}
}

func TestCompleteSession_All12ReproScenarios(t *testing.T) {
	remoteMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "remote",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{
				Name:    "orders",
				Columns: []*metadatapb.ColumnMetadata{{Name: "order_id"}, {Name: "amount"}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "order_view",
				Definition: "SELECT order_id FROM orders",
			}},
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "sync_orders",
				Definition: "CREATE PROCEDURE sync_orders() BEGIN SELECT 1; END",
			}},
		}},
	}
	appMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "appdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{
				{Name: "orders", Columns: []*metadatapb.ColumnMetadata{{Name: "local_col"}}},
				{Name: "café_tbl", Columns: []*metadatapb.ColumnMetadata{{Name: "id"}}},
				{Name: "sales$archive", Columns: []*metadatapb.ColumnMetadata{{Name: "id"}}},
			},
		}},
	}
	upperFoo := &metadatapb.DatabaseSchemaMetadata{
		Name: "Foo",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{
				Name:    "up",
				Columns: []*metadatapb.ColumnMetadata{{Name: "up_col"}},
			}},
		}},
	}
	lowerFoo := &metadatapb.DatabaseSchemaMetadata{
		Name: "foo",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{
				Name:    "lo",
				Columns: []*metadatapb.ColumnMetadata{{Name: "lo_col"}},
			}},
		}},
	}

	cfg := SessionConfig{
		DefaultDatabase: "appdb",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"appdb", "remote"}, nil
		},
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			switch db {
			case "appdb":
				return db, appMeta, true, true
			case "remote":
				return db, remoteMeta, true, true
			default:
				return "", nil, true, false
			}
		},
	}

	callCat := catalog.New()
	_, _ = callCat.LoadMetadataInto(context.Background(), remoteMeta, "remote")
	if db := callCat.GetDatabase("remote"); db != nil {
		db.Procedures["sync_orders"] = &catalog.Routine{Name: "sync_orders"}
	}
	gotCall := Complete("CALL remote.", len("CALL remote."), callCat)
	require.True(t, containsCandidate(gotCall, "sync_orders", CandidateProcedure), "got %v", gotCall)
	require.False(t, containsCandidate(gotCall, "orders", CandidateTable), "got %v", gotCall)
	require.False(t, containsCandidate(gotCall, "order_view", CandidateView), "got %v", gotCall)

	gotDropView := CompleteSession(context.Background(), cfg, "DROP VIEW remote.", 1, len("DROP VIEW remote."))
	require.True(t, containsCandidate(gotDropView, "order_view", CandidateView), "got %v", gotDropView)
	require.False(t, containsCandidate(gotDropView, "orders", CandidateTable), "got %v", gotDropView)

	gotCTE := CompleteSession(context.Background(), cfg, "WITH c AS (SELECT * FROM remote.orders) SELECT c. FROM c", 1, len("WITH c AS (SELECT * FROM remote.orders) SELECT c."))
	require.True(t, containsCandidate(gotCTE, "order_id", CandidateColumn), "got %v", gotCTE)
	require.True(t, containsCandidate(gotCTE, "amount", CandidateColumn), "got %v", gotCTE)
	require.False(t, containsCandidate(gotCTE, "local_col", CandidateColumn), "got %v", gotCTE)

	gotBacktick := CompleteSession(context.Background(), cfg, "SELECT * FROM `or", 1, len("SELECT * FROM `or"))
	require.True(t, containsCandidate(gotBacktick, "orders", CandidateTable), "got %v", gotBacktick)

	gotUnicode := CompleteSession(context.Background(), cfg, "SELECT * FROM café", 1, len("SELECT * FROM café"))
	require.True(t, containsCandidate(gotUnicode, "café_tbl", CandidateTable), "got %v", gotUnicode)

	gotDollar := CompleteSession(context.Background(), cfg, "SELECT * FROM sales$a", 1, len("SELECT * FROM sales$a"))
	require.True(t, containsCandidate(gotDollar, "sales$archive", CandidateTable), "got %v", gotDollar)

	csCfg := SessionConfig{
		DefaultDatabase: "foo",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			switch db {
			case "Foo":
				return db, upperFoo, true, true
			case "foo":
				return db, lowerFoo, true, true
			default:
				return "", nil, true, false
			}
		},
	}
	gotBothCS := CompleteSession(context.Background(), csCfg, "SELECT  FROM Foo.up JOIN foo.lo", 1, len("SELECT "))
	require.True(t, containsCandidate(gotBothCS, "up_col", CandidateColumn), "got %v", gotBothCS)
	require.True(t, containsCandidate(gotBothCS, "lo_col", CandidateColumn), "got %v", gotBothCS)

	gotWrongCase := CompleteSession(context.Background(), csCfg, "SELECT * FROM FOO.", 1, len("SELECT * FROM FOO."))
	require.False(t, containsCandidate(gotWrongCase, "lo", CandidateTable), "got %v", gotWrongCase)

	gotCommentDot := CompleteSession(context.Background(), cfg, "SELECT * FROM remote /* note */ .", 1, len("SELECT * FROM remote /* note */ ."))
	require.True(t, containsCandidate(gotCommentDot, "orders", CandidateTable), "got %v", gotCommentDot)
	require.False(t, containsCandidate(gotCommentDot, "café_tbl", CandidateTable), "got %v", gotCommentDot)

	gotUsingComma := CompleteSession(context.Background(), cfg, "SELECT * FROM t JOIN u USING (id), remote.", 1, len("SELECT * FROM t JOIN u USING (id), remote."))
	require.True(t, containsCandidate(gotUsingComma, "orders", CandidateTable), "got %v", gotUsingComma)

	gotOnComma := CompleteSession(context.Background(), cfg, "SELECT * FROM t JOIN u ON t.id=u.id, remote.", 1, len("SELECT * FROM t JOIN u ON t.id=u.id, remote."))
	require.True(t, containsCandidate(gotOnComma, "orders", CandidateTable), "got %v", gotOnComma)

	noDefCfg := cfg
	noDefCfg.DefaultDatabase = ""
	gotNoDef := CompleteSession(context.Background(), noDefCfg, "SELECT * FROM  JOIN remote.orders", 1, len("SELECT * FROM "))
	require.False(t, containsCandidate(gotNoDef, "orders", CandidateTable), "got %v", gotNoDef)

	gotMissing := CompleteSession(context.Background(), cfg, "SELECT * FROM  JOIN missing.t", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(gotMissing, "remote", CandidateDatabase), "got %v", gotMissing)
	require.False(t, containsCandidate(gotMissing, "missing", CandidateDatabase), "got %v", gotMissing)

	gotDropLeft := postProcessSessionCandidates("DROP FUNCTION LEF", len("DROP FUNCTION LEF"), false, []Candidate{
		{Type: CandidateFunction, Text: "LEFT"},
	}, func(string, int) bool { return false })
	require.True(t, containsCandidate(gotDropLeft, "`LEFT`", CandidateFunction), "got %v", gotDropLeft)
	require.False(t, containsCandidate(gotDropLeft, "LEFT", CandidateFunction), "got %v", gotDropLeft)

	roCfg := cfg
	roCfg.ReadOnly = true
	gotNestedRO := CompleteSession(context.Background(), roCfg, "WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) ", 1, len("WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) "))
	require.True(t, containsCandidate(gotNestedRO, "SELECT", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "UPDATE", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "DELETE", CandidateKeyword), "got %v", gotNestedRO)

	gotCommentParen := postProcessSessionCandidates("SELECT my /* note */ ()", len("SELECT my"), false, []Candidate{
		{Type: CandidateFunction, Text: "my-func"},
	}, func(string, int) bool { return true })
	require.True(t, containsCandidate(gotCommentParen, "`my-func`", CandidateFunction), "got %v", gotCommentParen)
	require.False(t, containsCandidate(gotCommentParen, "`my-func`()", CandidateFunction), "got %v", gotCommentParen)

	require.False(t, StatementReferencesDatabase("SELECT `old.db.foo` FROM t", "db"))
	require.True(t, StatementReferencesDatabase("SELECT * FROM `a``b`.t", "a`b"))
}

func TestCompleteSession_CodexRound4Fixes(t *testing.T) {
	dbs := map[string]*metadatapb.DatabaseSchemaMetadata{
		"appdb": {
			Name: "appdb",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "t1",
						Columns: []*metadatapb.ColumnMetadata{{Name: "outer_col"}, {Name: "shared_id"}},
					},
					{
						Name:    "t2",
						Columns: []*metadatapb.ColumnMetadata{{Name: "inner_col"}, {Name: "shared_id"}},
					},
				},
			}},
		},
		"remote": {
			Name: "remote",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "orders",
						Columns: []*metadatapb.ColumnMetadata{{Name: "order_id"}, {Name: "amount"}},
					},
				},
			}},
		},
	}
	var loaded []string
	cfg := SessionConfig{
		DefaultDatabase: "appdb",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			loaded = append(loaded, db)
			meta, ok := dbs[strings.ToLower(db)]
			if !ok {
				return db, nil, false, false
			}
			return meta.Name, meta, false, true
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"appdb", "remote"}, nil
		},
	}

	// 1. Correlated outer table alias resolves inside subquery SELECT & EXISTS,
	//    while matching inner aliases shadow outer ones.
	gotCorr := CompleteSession(context.Background(), cfg, "SELECT (SELECT o. FROM t2 i) FROM t1 o", 1, len("SELECT (SELECT o."))
	require.True(t, containsCandidate(gotCorr, "outer_col", CandidateColumn), "got %v", gotCorr)
	require.False(t, containsCandidate(gotCorr, "inner_col", CandidateColumn), "got %v", gotCorr)

	gotExists := CompleteSession(context.Background(), cfg, "SELECT * FROM t1 o WHERE EXISTS (SELECT 1 FROM t2 i WHERE o.)", 1, len("SELECT * FROM t1 o WHERE EXISTS (SELECT 1 FROM t2 i WHERE o."))
	require.True(t, containsCandidate(gotExists, "outer_col", CandidateColumn), "got %v", gotExists)
	require.False(t, containsCandidate(gotExists, "inner_col", CandidateColumn), "got %v", gotExists)

	gotShadow := CompleteSession(context.Background(), cfg, "SELECT (SELECT o. FROM t2 o) FROM t1 o", 1, len("SELECT (SELECT o."))
	require.True(t, containsCandidate(gotShadow, "inner_col", CandidateColumn), "got %v", gotShadow)
	require.False(t, containsCandidate(gotShadow, "outer_col", CandidateColumn), "got %v", gotShadow)

	// 2. Executable comments retain qualifier discovery and target lookup.
	loaded = nil
	gotExecWhere := CompleteSession(context.Background(), cfg, "SELECT * FROM /*! remote.orders */ WHERE ", 1, len("SELECT * FROM /*! remote.orders */ WHERE "))
	require.Contains(t, loaded, "remote")
	require.True(t, containsCandidate(gotExecWhere, "order_id", CandidateColumn), "got %v", gotExecWhere)
	require.True(t, containsCandidate(gotExecWhere, "amount", CandidateColumn), "got %v", gotExecWhere)

	gotExecVerWhere := CompleteSession(context.Background(), cfg, "SELECT * FROM /*!50000 remote.orders */ WHERE ", 1, len("SELECT * FROM /*!50000 remote.orders */ WHERE "))
	require.True(t, containsCandidate(gotExecVerWhere, "order_id", CandidateColumn), "got %v", gotExecVerWhere)

	gotExecDot := CompleteSession(context.Background(), cfg, "SELECT * FROM /*! remote. */", 1, len("SELECT * FROM /*! remote."))
	require.True(t, containsCandidate(gotExecDot, "orders", CandidateTable), "got %v", gotExecDot)

	// 3. Repeated duplicate wildcard expansion in chained CTEs stays bounded.
	var sb strings.Builder
	sb.WriteString("WITH c0 AS (SELECT 1 AS col)")
	for i := 1; i <= 25; i++ {
		fmt.Fprintf(&sb, ", c%d AS (SELECT c%d.*, c%d.* FROM c%d)", i, i-1, i-1, i-1)
	}
	sb.WriteString(" SELECT c25. FROM c25")
	sql := sb.String()
	gotCTEChain := CompleteSession(context.Background(), cfg, sql, 1, len(sql)-len(" FROM c25"))
	require.Len(t, gotCTEChain, 1)
	require.True(t, containsCandidate(gotCTEChain, "col", CandidateColumn), "got %v", gotCTEChain)

	// 4. Wrapper CTEs on set-operation derived tables project columns accurately.
	unionSQL := "SELECT d. FROM (WITH x AS (SELECT 1 AS n) SELECT * FROM x UNION ALL SELECT * FROM x) d"
	gotUnionDerived := CompleteSession(context.Background(), cfg, unionSQL, 1, len("SELECT d."))
	require.True(t, containsCandidate(gotUnionDerived, "n", CandidateColumn), "got %v", gotUnionDerived)

	// 5. Recursive CTE columns are visible inside the recursive arm (both explicit and implicit).
	recExplicit := "WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT c. FROM c) SELECT * FROM c"
	gotRecExplicit := CompleteSession(context.Background(), cfg, recExplicit, 1, len("WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT c."))
	require.True(t, containsCandidate(gotRecExplicit, "n", CandidateColumn), "got %v", gotRecExplicit)

	recImplicit := "WITH RECURSIVE c AS (SELECT 1 AS n UNION ALL SELECT c. FROM c) SELECT * FROM c"
	gotRecImplicit := CompleteSession(context.Background(), cfg, recImplicit, 1, len("WITH RECURSIVE c AS (SELECT 1 AS n UNION ALL SELECT c."))
	require.True(t, containsCandidate(gotRecImplicit, "n", CandidateColumn), "got %v", gotRecImplicit)
}

func TestCompleteSession_UnparseableColumnInstanceDBsQualifiersAndCrossDB(t *testing.T) {
	// Unparseable column type falls back to stand-in so all columns still surface.
	standInCfg := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{{
					Name: "t",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "good", Type: "int"},
						{Name: "weird", Type: ")"},
					},
				}},
			}},
		}, true),
	}
	gotStandIn := CompleteSession(context.Background(), standInCfg, "SELECT  FROM t", 1, len("SELECT "))
	require.True(t, containsCandidate(gotStandIn, "good", CandidateColumn), "got %v", gotStandIn)
	require.True(t, containsCandidate(gotStandIn, "weird", CandidateColumn), "got %v", gotStandIn)

	// Instance-level completion without DefaultDatabase surfaces database names.
	instCfg := SessionConfig{
		DefaultDatabase: "",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			return db, &metadatapb.DatabaseSchemaMetadata{Name: db, Schemas: []*metadatapb.SchemaMetadata{{}}}, true, true
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"appdb", "otherdb"}, nil
		},
	}
	gotInst := CompleteSession(context.Background(), instCfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(gotInst, "appdb", CandidateDatabase), "got %v", gotInst)
	require.True(t, containsCandidate(gotInst, "otherdb", CandidateDatabase), "got %v", gotInst)

	// Statement isolation, views, core FROM/JOIN/column contexts, qualifier scoping,
	// function parens, and backtick inside string literal.
	cfg := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{
					{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "a1"}, {Name: "a2"}}},
					{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "b1"}, {Name: "b2"}}},
					{Name: "order", Columns: []*metadatapb.ColumnMetadata{{Name: "id"}}},
				},
				Views: []*metadatapb.ViewMetadata{{Name: "v1", Definition: "SELECT a1 FROM t1"}},
			}},
		}, true),
	}

	gotIso := CompleteSession(context.Background(), cfg, "SELECT * FROM t1; SELECT  FROM t2", 1, len("SELECT * FROM t1; SELECT "))
	require.True(t, containsCandidate(gotIso, "b1", CandidateColumn), "got %v", gotIso)
	require.True(t, containsCandidate(gotIso, "b2", CandidateColumn), "got %v", gotIso)
	require.False(t, containsCandidate(gotIso, "a1", CandidateColumn), "got %v", gotIso)
	require.False(t, containsCandidate(gotIso, "a2", CandidateColumn), "got %v", gotIso)

	// Table/view contexts: FROM, JOIN, comma-join.
	for _, stmt := range []string{"SELECT * FROM ", "SELECT * FROM t1 JOIN ", "SELECT * FROM t1, "} {
		got := CompleteSession(context.Background(), cfg, stmt, 1, len(stmt))
		require.True(t, containsCandidate(got, "t1", CandidateTable), "%q: got %v", stmt, got)
		require.True(t, containsCandidate(got, "t2", CandidateTable), "%q: got %v", stmt, got)
		require.True(t, containsCandidate(got, "v1", CandidateView), "%q: got %v", stmt, got)
	}

	// Column contexts: projection, table-qualified, WHERE, JOIN ON.
	for _, tc := range []struct {
		stmt string
		pos  int
		want []string
	}{
		{"SELECT  FROM t1", len("SELECT "), []string{"a1", "a2"}},
		{"SELECT t1. FROM t1", len("SELECT t1."), []string{"a1", "a2"}},
		{"SELECT * FROM t1 WHERE ", len("SELECT * FROM t1 WHERE "), []string{"a1", "a2"}},
		{"SELECT  FROM t1 JOIN t2 ON t1.a1 = t2.b1", len("SELECT "), []string{"a1", "a2", "b1", "b2"}},
	} {
		got := CompleteSession(context.Background(), cfg, tc.stmt, 1, tc.pos)
		for _, col := range tc.want {
			require.True(t, containsCandidate(got, col, CandidateColumn), "%q missing %s: got %v", tc.stmt, col, got)
		}
	}

	// Qualified column scoping in JOIN.
	for _, tc := range []struct{ stmt, at, want, notWant string }{
		{"SELECT a. FROM t1 AS a JOIN t2 AS b", "SELECT a.", "a1", "b1"},
		{"SELECT b. FROM t1 AS a JOIN t2 AS b", "SELECT b.", "b1", "a1"},
		{"SELECT t1. FROM t1 JOIN t2", "SELECT t1.", "a1", "b1"},
	} {
		got := CompleteSession(context.Background(), cfg, tc.stmt, 1, len(tc.at))
		require.True(t, containsCandidate(got, tc.want, CandidateColumn), "%q: got %v", tc.stmt, got)
		require.False(t, containsCandidate(got, tc.notWant, CandidateColumn), "%q: got %v", tc.stmt, got)
	}

	// Unknown qualifier offers no columns.
	for _, tc := range []struct{ stmt, at string }{
		{"SELECT x. FROM t1 AS a JOIN t2 AS b", "SELECT x."},
		{"SELECT t3. FROM t1 JOIN t2", "SELECT t3."},
	} {
		got := CompleteSession(context.Background(), cfg, tc.stmt, 1, len(tc.at))
		require.False(t, containsCandidate(got, "a1", CandidateColumn), "%q: got %v", tc.stmt, got)
		require.False(t, containsCandidate(got, "b1", CandidateColumn), "%q: got %v", tc.stmt, got)
	}

	// Function parens and backtick inside string literal.
	gotFunc := CompleteSession(context.Background(), cfg, "SELECT  FROM t1", 1, len("SELECT "))
	require.True(t, containsCandidate(gotFunc, "COUNT()", CandidateFunction), "got %v", gotFunc)
	require.False(t, containsCandidate(gotFunc, "COUNT", CandidateFunction), "got %v", gotFunc)

	gotBacktickStr := CompleteSession(context.Background(), cfg, "SELECT '`' FROM ", 1, len("SELECT '`' FROM "))
	require.True(t, containsCandidate(gotBacktickStr, "`order`", CandidateTable), "got %v", gotBacktickStr)

	// Cross-DB qualified column resolution without enumerating databases.
	appMeta := &metadatapb.DatabaseSchemaMetadata{
		Name:    "appdb",
		Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}}}}}},
	}
	otherMeta := &metadatapb.DatabaseSchemaMetadata{
		Name:    "otherdb",
		Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "c2"}, {Name: "c3"}}}}}},
	}
	listCalls := 0
	crossCfg := SessionConfig{
		DefaultDatabase: "appdb",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			if db == "otherdb" {
				return db, otherMeta, true, true
			}
			return db, appMeta, true, true
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			listCalls++
			return []string{"appdb", "otherdb"}, nil
		},
	}
	gotCross := CompleteSession(context.Background(), crossCfg, "SELECT otherdb.t2. FROM otherdb.t2", 1, len("SELECT otherdb.t2."))
	require.True(t, containsCandidate(gotCross, "c2", CandidateColumn), "got %v", gotCross)
	require.True(t, containsCandidate(gotCross, "c3", CandidateColumn), "got %v", gotCross)
	require.Zero(t, listCalls)

	gotDBList := CompleteSession(context.Background(), crossCfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(gotDBList, "otherdb", CandidateDatabase), "got %v", gotDBList)
	require.Equal(t, 1, listCalls)
}
