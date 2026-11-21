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
}

func TestCompleteSession_UnparseableColumnTypeFallsBackToStandIn(t *testing.T) {
	cfg := SessionConfig{
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

	got := CompleteSession(context.Background(), cfg, "SELECT  FROM t", 1, len("SELECT "))
	require.True(t, containsCandidate(got, "good", CandidateColumn), "got %v", got)
	require.True(t, containsCandidate(got, "weird", CandidateColumn), "got %v", got)
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
	for _, tc := range []struct{ stmt, keep string }{
		{"SHOW ", "CREATE"},
		{"SELECT 1 FOR ", "UPDATE"},
	} {
		got := CompleteSession(context.Background(), roCfg, tc.stmt, 1, len(tc.stmt))
		require.True(t, containsCandidate(got, tc.keep, CandidateKeyword), "%q after %q", tc.keep, tc.stmt)
	}

	gotCre := CompleteSession(context.Background(), roCfg, "CRE", 1, len("CRE"))
	require.False(t, containsCandidate(gotCre, "CREATE", CandidateKeyword))

	gotIns := CompleteSession(context.Background(), roCfg, "SELECT 1; INS", 1, len("SELECT 1; INS"))
	require.False(t, containsCandidate(gotIns, "INSERT", CandidateKeyword))

	gotShowCre := CompleteSession(context.Background(), roCfg, "SHOW CRE", 1, len("SHOW CRE"))
	require.True(t, containsCandidate(gotShowCre, "CREATE", CandidateKeyword))

	gotEmpty := CompleteSession(context.Background(), roCfg, "", 1, 0)
	require.False(t, containsCandidate(gotEmpty, "SHUTDOWN", CandidateKeyword))

	for _, postCTE := range []string{"WITH c AS (SELECT 1) ", "WITH c AS (SELECT 1) UP", "WITH c AS (SELECT 1) DEL"} {
		gotPostCTE := CompleteSession(context.Background(), roCfg, postCTE, 1, len(postCTE))
		require.False(t, containsCandidate(gotPostCTE, "UPDATE", CandidateKeyword), "postCTE %q", postCTE)
		require.False(t, containsCandidate(gotPostCTE, "DELETE", CandidateKeyword), "postCTE %q", postCTE)
	}
}

func TestCompleteSession_CrossDBAndLazyEnumeration(t *testing.T) {
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
	loads := map[string]int{}
	listCalls := 0
	cfg := SessionConfig{
		DefaultDatabase: "appdb",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			loads[db]++
			switch db {
			case "appdb":
				return db, appMeta, true, true
			case "otherdb":
				return db, otherMeta, true, true
			default:
				return "", nil, false, false
			}
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			listCalls++
			return []string{"appdb", "otherdb", "unuseddb"}, nil
		},
	}

	got := CompleteSession(context.Background(), cfg, "SELECT otherdb.t2. FROM otherdb.t2", 1, len("SELECT otherdb.t2."))
	require.True(t, containsCandidate(got, "c2", CandidateColumn))
	require.True(t, containsCandidate(got, "c3", CandidateColumn))
	require.Zero(t, loads["unuseddb"])
	require.Zero(t, listCalls)

	// Qualified object completion resolves against the referenced DB (otherdb), not DefaultDatabase (appdb).
	otherFullMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "otherdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "other_tbl", Columns: []*metadatapb.ColumnMetadata{{Name: "c2"}}}},
			Views:  []*metadatapb.ViewMetadata{{Name: "other_view", Definition: "SELECT c2 FROM other_tbl"}},
		}},
	}
	otherOnlyCfg := SessionConfig{
		DefaultDatabase: "appdb",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			if db == "otherdb" {
				return db, otherFullMeta, true, true
			}
			return db, appMeta, true, true
		},
	}
	gotOtherTbl := CompleteSession(context.Background(), otherOnlyCfg, "SELECT * FROM otherdb.", 1, len("SELECT * FROM otherdb."))
	require.True(t, containsCandidate(gotOtherTbl, "other_tbl", CandidateTable), "got %v", gotOtherTbl)
	require.True(t, containsCandidate(gotOtherTbl, "other_view", CandidateView), "got %v", gotOtherTbl)
	require.False(t, containsCandidate(gotOtherTbl, "t1", CandidateTable), "must not leak appdb.t1 into otherdb.; got %v", gotOtherTbl)

	gotOtherView := CompleteSession(context.Background(), otherOnlyCfg, "DROP VIEW otherdb.", 1, len("DROP VIEW otherdb."))
	require.True(t, containsCandidate(gotOtherView, "other_view", CandidateView), "got %v", gotOtherView)

	got2 := CompleteSession(context.Background(), cfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(got2, "otherdb", CandidateDatabase))
	require.True(t, containsCandidate(got2, "unuseddb", CandidateDatabase))
	require.Equal(t, 1, listCalls)
}

func TestCompleteSession_CaseCollidingAndFoldedQualifiers(t *testing.T) {
	caseSensCfg := SessionConfig{
		DefaultDatabase: "foo",
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			switch db {
			case "foo":
				return db, &metadatapb.DatabaseSchemaMetadata{
					Name:    "foo",
					Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "lower_col", Type: "int"}}}}}},
				}, true, true
			case "Foo":
				return db, &metadatapb.DatabaseSchemaMetadata{
					Name:    "Foo",
					Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "upper_col", Type: "int"}}}}}},
				}, true, true
			default:
				return "", nil, false, false
			}
		},
	}

	gotUpper := CompleteSession(context.Background(), caseSensCfg, "SELECT Foo.t. FROM Foo.t", 1, len("SELECT Foo.t."))
	require.True(t, containsCandidate(gotUpper, "upper_col", CandidateColumn))

	stmtMulti := "SELECT * FROM Foo.t JOIN foo.t ON Foo.t.lower_col = foo.t.lower_col WHERE foo.t."
	gotLower := CompleteSession(context.Background(), caseSensCfg, stmtMulti, 1, len(stmtMulti))
	require.True(t, containsCandidate(gotLower, "lower_col", CandidateColumn))
	require.False(t, containsCandidate(gotLower, "upper_col", CandidateColumn))

	gotDerived := CompleteSession(context.Background(), caseSensCfg, "SELECT x. FROM (SELECT * FROM Foo.t) x", 1, len("SELECT x."))
	require.True(t, containsCandidate(gotDerived, "upper_col", CandidateColumn))
	require.False(t, containsCandidate(gotDerived, "lower_col", CandidateColumn))

	gotFolded := CompleteSession(context.Background(), SessionConfig{
		DefaultDatabase: "",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"appdb", "db_b"}, nil
		},
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			if db != "db_b" {
				return "", nil, false, false
			}
			return db, &metadatapb.DatabaseSchemaMetadata{
				Name:    "db_b",
				Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "orders", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "int"}}}}}},
			}, false, true
		},
	}, "SELECT DB_B.orders. FROM DB_B.orders", 1, len("SELECT DB_B.orders."))
	require.True(t, containsCandidate(gotFolded, "id", CandidateColumn))
}

func TestCompleteSession_StatementIsolationUTF16AndHelpers(t *testing.T) {
	cfg := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: sessionMetaLoader(&metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{
					{Name: "t1", Columns: []*metadatapb.ColumnMetadata{{Name: "a1"}}},
					{Name: "t2", Columns: []*metadatapb.ColumnMetadata{{Name: "b1"}}},
					{Name: "order", Columns: []*metadatapb.ColumnMetadata{{Name: "id"}}},
				},
				Views: []*metadatapb.ViewMetadata{{Name: "v1", Definition: "SELECT a1 FROM t1"}},
			}},
		}, true),
	}

	gotIso := CompleteSession(context.Background(), cfg, "SELECT * FROM t1; SELECT  FROM t2", 1, len("SELECT * FROM t1; SELECT "))
	require.True(t, containsCandidate(gotIso, "b1", CandidateColumn))
	require.False(t, containsCandidate(gotIso, "a1", CandidateColumn))

	gotFunc := CompleteSession(context.Background(), cfg, "SELECT  FROM t1", 1, len("SELECT "))
	require.True(t, containsCandidate(gotFunc, "COUNT()", CandidateFunction))
	require.False(t, containsCandidate(gotFunc, "COUNT", CandidateFunction))

	gotFuncParen := CompleteSession(context.Background(), cfg, "SELECT CO(a1) FROM t1", 1, len("SELECT CO"))
	require.True(t, containsCandidate(gotFuncParen, "COUNT", CandidateFunction), "got %v", gotFuncParen)
	require.False(t, containsCandidate(gotFuncParen, "COUNT()", CandidateFunction), "must not insert duplicate () before existing (; got %v", gotFuncParen)

	gotView := CompleteSession(context.Background(), cfg, "SELECT  FROM v1", 1, len("SELECT "))
	require.True(t, containsCandidate(gotView, "a1", CandidateColumn))

	gotQuote := CompleteSession(context.Background(), cfg, "SELECT '`' FROM ", 1, len("SELECT '`' FROM "))
	require.True(t, containsCandidate(gotQuote, "`order`", CandidateTable))

	roCfg := cfg
	roCfg.ReadOnly = true
	gotGap := CompleteSession(context.Background(), roCfg, "SELECT * FROM t1; ; SELEC broken", 1, len("SELECT * FROM t1; "))
	require.True(t, containsCandidate(gotGap, "SELECT", CandidateKeyword))
	require.False(t, containsCandidate(gotGap, "INSERT", CandidateKeyword))

	stmtEmoji := "SELECT '😀',  FROM t1"
	gotEmoji := CompleteSession(context.Background(), cfg, stmtEmoji, 1, len("SELECT '")+2+len("', "))
	require.True(t, containsCandidate(gotEmoji, "a1", CandidateColumn))

	require.Equal(t, "`order`", QuoteIdentifierIfNeeded("order", false))
	require.Equal(t, "order", QuoteIdentifierIfNeeded("order", true))
	require.Equal(t, "users", QuoteIdentifierIfNeeded("users", false))
	require.True(t, CaretInsideBacktickIdentifier("SELECT * FROM `o", 16))
	require.False(t, CaretInsideBacktickIdentifier("SELECT '`' FROM ", 16))
	require.True(t, StatementReferencesDatabase("SELECT * FROM `otherdb` . t2", "otherdb"))
	require.False(t, StatementReferencesDatabase("SELECT /* otherdb. */ * FROM t", "otherdb"))
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
