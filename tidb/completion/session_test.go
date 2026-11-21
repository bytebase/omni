package completion

import (
	"context"
	"testing"

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
