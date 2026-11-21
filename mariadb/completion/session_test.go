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
