package completion

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	metadatapb "github.com/bytebase/omni/metadata"
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

	gotRoutine := CompleteSession(context.Background(), SessionConfig{
		DefaultDatabase: "testdb",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "testdb",
			Schemas: []*metadatapb.SchemaMetadata{{Functions: []*metadatapb.FunctionMetadata{
				{Name: "my_udf", Definition: "CREATE FUNCTION my_udf() RETURNS INT RETURN 1"},
				{Name: "my-func", Definition: "CREATE FUNCTION `my-func`() RETURNS INT RETURN 1"},
				{Name: "f()", Definition: "CREATE FUNCTION `f()`() RETURNS INT RETURN 1"},
			}}},
		}, true),
	}, "DROP FUNCTION ", 1, len("DROP FUNCTION "))
	require.True(t, containsCandidate(gotRoutine, "my_udf", CandidateFunction), "got %v", gotRoutine)
	require.True(t, containsCandidate(gotRoutine, "`my-func`", CandidateFunction), "got %v", gotRoutine)
	require.True(t, containsCandidate(gotRoutine, "`f()`", CandidateFunction), "got %v", gotRoutine)
	require.False(t, containsCandidate(gotRoutine, "my_udf()", CandidateFunction), "got %v", gotRoutine)
	require.False(t, containsCandidate(gotRoutine, "`my-func`()", CandidateFunction), "got %v", gotRoutine)
	require.False(t, containsCandidate(gotRoutine, "f()", CandidateFunction), "got %v", gotRoutine)

	gotFuncExpr := CompleteSession(context.Background(), SessionConfig{
		DefaultDatabase: "testdb",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "testdb",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}}}},
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "my-func", Definition: "CREATE FUNCTION `my-func`() RETURNS INT RETURN 1"},
					{Name: "f()", Definition: "CREATE FUNCTION `f()`() RETURNS INT RETURN 1"},
				},
			}},
		}, true),
	}, "SELECT  FROM t1", 1, len("SELECT "))
	require.True(t, containsCandidate(gotFuncExpr, "`my-func`()", CandidateFunction), "got %v", gotFuncExpr)
	require.True(t, containsCandidate(gotFuncExpr, "`f()`()", CandidateFunction), "got %v", gotFuncExpr)
	require.True(t, containsCandidate(gotFuncExpr, "LEFT()", CandidateFunction), "got %v", gotFuncExpr)
	require.False(t, containsCandidate(gotFuncExpr, "`LEFT`()", CandidateFunction), "got %v", gotFuncExpr)
	require.False(t, containsCandidate(gotFuncExpr, "f()", CandidateFunction), "got %v", gotFuncExpr)
}

func TestCompleteSession_ReadOnlyFiltersWriteKeywords(t *testing.T) {
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
	for _, kw := range []string{"INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP"} {
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
}

func TestCompleteSession_CrossDBAndHelpers(t *testing.T) {
	appMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "appdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}}}},
		}},
	}
	otherMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "otherdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "c2"}, {Name: "c3"}}}},
		}},
	}
	listCalls := 0
	cfg := SessionConfig{
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

	got := CompleteSession(context.Background(), cfg, "SELECT otherdb.t2. FROM otherdb.t2", 1, len("SELECT otherdb.t2."))
	require.True(t, containsCandidate(got, "c2", CandidateColumn))
	require.True(t, containsCandidate(got, "c3", CandidateColumn))
	require.Zero(t, listCalls)

	gotOtherTbl := CompleteSession(context.Background(), cfg, "SELECT * FROM otherdb.", 1, len("SELECT * FROM otherdb."))
	require.True(t, containsCandidate(gotOtherTbl, "t2", CandidateTable), "got %v", gotOtherTbl)
	require.False(t, containsCandidate(gotOtherTbl, "t1", CandidateTable), "got %v", gotOtherTbl)

	got2 := CompleteSession(context.Background(), cfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(got2, "otherdb", CandidateDatabase))
	require.Equal(t, 1, listCalls)

	require.Equal(t, "`order`", QuoteIdentifierIfNeeded("order", false))
	require.Equal(t, "order", QuoteIdentifierIfNeeded("order", true))
	require.True(t, CaretInsideBacktickIdentifier("SELECT * FROM `o", 16))
	require.True(t, StatementReferencesDatabase("SELECT * FROM `otherdb`.t", "otherdb"))
}

func TestCompleteSession_BrokenTrailingStatementScalesLinearly(t *testing.T) {
	var sheet strings.Builder
	sheet.WriteString("SELECT a2. FROM t1 JOIN t2 a2 ON t1.c1 = a2.c1;\n")
	sheet.WriteString("SELEC broken FROM oops;\n")
	for i := range 800 {
		fmt.Fprintf(&sheet, "SELECT col_a, col_b, col_c FROM table_%04d WHERE col_a = %d AND col_b LIKE 'pattern%%' ORDER BY col_c LIMIT 100;\n", i, i)
	}

	cfg := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{
					{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}}},
					{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "c1"}, {Name: "c2"}}},
				},
			}},
		}, true),
	}

	started := time.Now()
	got := CompleteSession(context.Background(), cfg, sheet.String(), 1, 10)
	elapsed := time.Since(started)

	require.True(t, containsCandidate(got, "c1", CandidateColumn))
	require.True(t, containsCandidate(got, "c2", CandidateColumn))
	require.Less(t, elapsed, 2*time.Second)
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
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "LEFT", Definition: "CREATE FUNCTION `LEFT`() RETURNS INT RETURN 1"},
				{Name: "my-func", Definition: "CREATE FUNCTION `my-func`() RETURNS INT RETURN 1"},
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

	gotCall := CompleteSession(context.Background(), cfg, "CALL remote.", 1, len("CALL remote."))
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

	gotDropLeft := CompleteSession(context.Background(), cfg, "DROP FUNCTION LEF", 1, len("DROP FUNCTION LEF"))
	require.True(t, containsCandidate(gotDropLeft, "`LEFT`", CandidateFunction), "got %v", gotDropLeft)
	require.False(t, containsCandidate(gotDropLeft, "LEFT", CandidateFunction), "got %v", gotDropLeft)

	roCfg := cfg
	roCfg.ReadOnly = true
	gotNestedRO := CompleteSession(context.Background(), roCfg, "WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) ", 1, len("WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) "))
	require.True(t, containsCandidate(gotNestedRO, "SELECT", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "UPDATE", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "DELETE", CandidateKeyword), "got %v", gotNestedRO)

	gotCommentParen := CompleteSession(context.Background(), cfg, "SELECT my /* note */ ()", 1, len("SELECT my"))
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
