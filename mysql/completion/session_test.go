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
					Name: "foo",
					Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{
						{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "lower_col", Type: "int"}}},
						{Name: "lower_only", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "int"}}},
					}}},
				}, true, true
			case "Foo":
				return db, &metadatapb.DatabaseSchemaMetadata{
					Name: "Foo",
					Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{
						{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "upper_col", Type: "int"}}},
						{Name: "upper_only", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "int"}}},
					}}},
				}, true, true
			default:
				return "", nil, false, false
			}
		},
	}

	gotUpper := CompleteSession(context.Background(), caseSensCfg, "SELECT Foo.t. FROM Foo.t", 1, len("SELECT Foo.t."))
	require.True(t, containsCandidate(gotUpper, "upper_col", CandidateColumn))
	require.False(t, containsCandidate(gotUpper, "lower_col", CandidateColumn))

	// Colliding default DB (foo) must NOT merge its distinct tables (lower_only) into Foo.
	gotUpperTbl := CompleteSession(context.Background(), caseSensCfg, "SELECT * FROM Foo.", 1, len("SELECT * FROM Foo."))
	require.True(t, containsCandidate(gotUpperTbl, "upper_only", CandidateTable), "got %v", gotUpperTbl)
	require.False(t, containsCandidate(gotUpperTbl, "lower_only", CandidateTable), "foo.lower_only must not merge into Foo.; got %v", gotUpperTbl)

	// Unqualified default-DB completion keeps foo and does not merge Foo.upper_only.
	gotDefaultTbl := CompleteSession(context.Background(), caseSensCfg, "SELECT * FROM ", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(gotDefaultTbl, "lower_only", CandidateTable), "got %v", gotDefaultTbl)
	require.False(t, containsCandidate(gotDefaultTbl, "upper_only", CandidateTable), "got %v", gotDefaultTbl)

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

func TestCompleteSession_StatementIsolationUTF16AndFunctions(t *testing.T) {
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
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "my_udf", Definition: "CREATE FUNCTION my_udf() RETURNS INT RETURN 1"},
					{Name: "my-func", Definition: "CREATE FUNCTION `my-func`() RETURNS INT RETURN 1"},
					{Name: "f()", Definition: "CREATE FUNCTION `f()`() RETURNS INT RETURN 1"},
				},
			}},
		}, true),
	}

	gotIso := CompleteSession(context.Background(), cfg, "SELECT * FROM t1; SELECT  FROM t2", 1, len("SELECT * FROM t1; SELECT "))
	require.True(t, containsCandidate(gotIso, "b1", CandidateColumn))
	require.False(t, containsCandidate(gotIso, "a1", CandidateColumn))

	gotFunc := CompleteSession(context.Background(), cfg, "SELECT  FROM t1", 1, len("SELECT "))
	require.True(t, containsCandidate(gotFunc, "COUNT()", CandidateFunction))
	require.True(t, containsCandidate(gotFunc, "LEFT()", CandidateFunction), "reserved built-in must stay unquoted; got %v", gotFunc)
	require.False(t, containsCandidate(gotFunc, "`LEFT`()", CandidateFunction))
	require.True(t, containsCandidate(gotFunc, "my_udf()", CandidateFunction))
	require.True(t, containsCandidate(gotFunc, "`my-func`()", CandidateFunction), "got %v", gotFunc)
	require.True(t, containsCandidate(gotFunc, "`f()`()", CandidateFunction), "got %v", gotFunc)
	require.False(t, containsCandidate(gotFunc, "my-func()", CandidateFunction))
	require.False(t, containsCandidate(gotFunc, "f()", CandidateFunction))
	require.False(t, containsCandidate(gotFunc, "COUNT", CandidateFunction))

	gotFuncParen := CompleteSession(context.Background(), cfg, "SELECT CO(a1) FROM t1", 1, len("SELECT CO"))
	require.True(t, containsCandidate(gotFuncParen, "COUNT", CandidateFunction), "got %v", gotFuncParen)
	require.False(t, containsCandidate(gotFuncParen, "COUNT()", CandidateFunction), "must not insert duplicate () before existing (; got %v", gotFuncParen)

	// DDL and SHOW routine-name contexts must NOT append () to function candidates,
	// while still backtick-quoting non-identifier stored function names.
	for _, routineStmt := range []string{
		"DROP FUNCTION ",
		"DROP FUNCTION IF EXISTS ",
		"ALTER FUNCTION ",
		"SHOW CREATE FUNCTION ",
	} {
		gotRoutine := CompleteSession(context.Background(), cfg, routineStmt, 1, len(routineStmt))
		require.True(t, containsCandidate(gotRoutine, "my_udf", CandidateFunction), "%q should offer bare my_udf; got %v", routineStmt, gotRoutine)
		require.True(t, containsCandidate(gotRoutine, "`my-func`", CandidateFunction), "%q should quote `my-func`; got %v", routineStmt, gotRoutine)
		require.True(t, containsCandidate(gotRoutine, "`f()`", CandidateFunction), "%q should quote `f()`; got %v", routineStmt, gotRoutine)
		require.False(t, containsCandidate(gotRoutine, "f()", CandidateFunction), "%q must preserve parens as quoted identifier; got %v", routineStmt, gotRoutine)
		require.False(t, containsCandidate(gotRoutine, "my_udf()", CandidateFunction), "%q must not append (); got %v", routineStmt, gotRoutine)
		require.False(t, containsCandidate(gotRoutine, "`my-func`()", CandidateFunction), "%q must not append (); got %v", routineStmt, gotRoutine)
	}

	gotQualifiedFunc := postProcessSessionCandidates("DROP FUNCTION otherdb.", len("DROP FUNCTION otherdb."), false, []Candidate{
		{Type: CandidateFunction, Text: "LEFT"},
		{Type: CandidateFunction, Text: "my-func"},
	}, func(string, int) bool { return false })
	require.True(t, containsCandidate(gotQualifiedFunc, "`my-func`", CandidateFunction), "got %v", gotQualifiedFunc)
	require.True(t, containsCandidate(gotQualifiedFunc, "`LEFT`", CandidateFunction), "qualified stored function matching builtin keyword should quote; got %v", gotQualifiedFunc)

	gotQualifiedExprFunc := postProcessSessionCandidates("SELECT otherdb.", len("SELECT otherdb."), false, []Candidate{
		{Type: CandidateFunction, Text: "LEFT"},
		{Type: CandidateFunction, Text: "my-func"},
	}, func(string, int) bool { return true })
	require.True(t, containsCandidate(gotQualifiedExprFunc, "`LEFT`()", CandidateFunction), "got %v", gotQualifiedExprFunc)
	require.True(t, containsCandidate(gotQualifiedExprFunc, "`my-func`()", CandidateFunction), "got %v", gotQualifiedExprFunc)

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
	require.Equal(t, "a``b", QuoteIdentifierIfNeeded("a`b", true))
	require.Equal(t, "`a``b`", QuoteIdentifierIfNeeded("a`b", false))
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

func TestDatabaseQualifiersInStatement(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		defaultDB     string
		caseSensitive bool
		want          []string
	}{
		{name: "table alias column qualifier ignored", statement: "SELECT a. FROM users a", defaultDB: "appdb"},
		{name: "join alias qualifiers in ON ignored", statement: "SELECT a2. FROM t1 JOIN t2 a2 ON t1.c1 = a2.c1", defaultDB: "appdb"},
		{name: "update set and where table qualifiers ignored", statement: "UPDATE users u SET u.name = 'x' WHERE u.id = 1", defaultDB: "appdb"},
		{name: "extract from inside scalar function ignored", statement: "SELECT EXTRACT(YEAR FROM u.created_at) FROM users u", defaultDB: "appdb"},
		{name: "cross-db table and three-part column extracts db only", statement: "SELECT otherdb.t2. FROM otherdb.t2", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "cross-db in FROM with alias qualifier in SELECT extracts db only", statement: "SELECT u. FROM otherdb.users u", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "comma-joined cross-db table", statement: "SELECT a.id FROM users a, otherdb.orders o WHERE a.id = o.user_id", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "derived table subquery cross-db reference", statement: "SELECT sub.c1 FROM (SELECT a.c1 FROM otherdb.t2 a) sub WHERE sub.c1 = 1", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "default database excluded case-insensitively", statement: "SELECT * FROM APPDB.users u WHERE u.id = 1", defaultDB: "appdb"},
		{name: "preserves backtick-quoted database casing", statement: "SELECT * FROM `OtherDB`.`users` u WHERE u.id = 1", defaultDB: "appdb", want: []string{"OtherDB"}},
		{name: "table qualifier in insert column list ignored while target db extracted", statement: "INSERT INTO otherdb.t2 (t2.c1) VALUES (1)", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "insert with priority/ignore modifiers keeps target state before INTO", statement: "INSERT LOW_PRIORITY IGNORE INTO otherdb.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "replace with delayed modifier keeps target state before INTO", statement: "REPLACE DELAYED INTO otherdb.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "drop view preserves object list across comma", statement: "DROP VIEW appdb.v1, otherdb.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "drop temporary table preserves object list across comma", statement: "DROP TEMPORARY TABLE IF EXISTS appdb.t1, otherdb.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "analyze table maintenance statement preserves list across comma and modifier", statement: "ANALYZE NO_WRITE_TO_BINLOG TABLE appdb.t1, otherdb.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "check optimize repair checksum flush preserve list across comma", statement: "CHECK TABLE appdb.t1, db2.; OPTIMIZE LOCAL TABLE appdb.t1, db3.; REPAIR TABLE appdb.t1, db4.; CHECKSUM TABLE appdb.t1, db5.; FLUSH TABLES appdb.t1, db6.", defaultDB: "appdb", want: []string{"db2", "db3", "db4", "db5", "db6"}},
		{name: "on duplicate key update and for update do not treat qualifiers as databases", statement: "INSERT INTO appdb.t1 (id) VALUES (1) ON DUPLICATE KEY UPDATE id = a.id; SELECT * FROM t1 u FOR UPDATE OF u.", defaultDB: "appdb"},
		{name: "grant and revoke extract target db on ON and ignore role after FROM", statement: "GRANT SELECT, INSERT ON otherdb.t1 TO u; REVOKE UPDATE ON third_db.t2 FROM r.user", defaultDB: "appdb", want: []string{"otherdb", "third_db"}},
		{name: "show full columns and delete with modifiers extract db", statement: "SHOW FULL COLUMNS FROM otherdb.; DELETE LOW_PRIORITY IGNORE FROM third_db.t1 WHERE id = u.id", defaultDB: "appdb", want: []string{"otherdb", "third_db"}},
		{name: "ddl create trigger on cross-db table", statement: "CREATE TRIGGER tr BEFORE INSERT ON otherdb.t2 FOR EACH ROW SET NEW.c1 = 1", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "join with incomplete qualifier in FROM and ON", statement: "SELECT a. FROM users a JOIN otherdb.orders o ON a. = o.user_id WHERE o.", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "desc in order by ignored", statement: "SELECT * FROM users u ORDER BY u.id DESC", defaultDB: "appdb"},
		{name: "truncate without TABLE keyword extracts db while TRUNCATE function call does not", statement: "TRUNCATE otherdb.; SELECT TRUNCATE(a.amount, 2) FROM users a", defaultDB: "appdb", want: []string{"otherdb"}},
		{name: "drop sequence comma list and tidb import into extract dbs", statement: "DROP SEQUENCE IF EXISTS db1.s1, otherdb.; IMPORT INTO third_db.t1 FROM 's3://bucket/path'", defaultDB: "appdb", want: []string{"db1", "otherdb", "third_db"}},
		{name: "cache index and load index into cache extract table lists", statement: "CACHE INDEX db1.t1, db2. IN hot_cache; LOAD INDEX INTO CACHE db3.t1 INDEX (idx1), db4.", defaultDB: "appdb", want: []string{"db1", "db2", "db3", "db4"}},
		{name: "function return expression does not treat qualifier as database", statement: "CREATE FUNCTION fn() RETURNS INT RETURN a.x", defaultDB: "appdb"},
		{name: "dollar and unicode identifiers preserved as single database qualifiers", statement: "SELECT sales$archive.t1., café$2. FROM café$2.users u JOIN users a$1 ON a$1.id = u.id", defaultDB: "appdb", want: []string{"sales$archive", "café$2"}},
		{name: "parenthesized table references in UPDATE and DELETE USING extract dbs", statement: "UPDATE LOW_PRIORITY (otherdb.t1 JOIN db2.t2 ON t1.id = t2.id) SET t1.c = a.c; DELETE FROM db3.t1 USING (db4.t1 JOIN db5.t2 ON t1.id = t2.id) WHERE a.id = 1", defaultDB: "appdb", want: []string{"otherdb", "db2", "db3", "db4", "db5"}},
		{name: "case-sensitive mode preserves case-distinct database qualifiers", statement: "SELECT Foo.t. FROM Foo.t", defaultDB: "foo", caseSensitive: true, want: []string{"Foo"}},
		{name: "delete multi-table target before FROM extracts db while ignoring alias wildcard", statement: "DELETE otherdb., db2.t1 FROM otherdb.t1 JOIN db2.t1; DELETE a.* FROM users a", defaultDB: "appdb", want: []string{"otherdb", "db2"}},
		{name: "grant wildcard on database still extracts db", statement: "GRANT SELECT ON otherdb.* TO 'u'@'%'", defaultDB: "appdb", want: []string{"otherdb"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := databaseQualifiersInStatement(tc.statement, tc.defaultDB, tc.caseSensitive)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestPopulateCatalog_SkipsEnumerationAndCollidingMerges(t *testing.T) {
	listCalls := 0
	var loaded []string
	var registered []string

	cfg := SessionConfig{
		DefaultDatabase: "appdb",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			listCalls++
			return []string{"appdb", "otherdb", "unuseddb"}, nil
		},
		LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
			loaded = append(loaded, db)
			return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
		},
	}

	populateCatalog(
		context.Background(),
		cfg,
		"SELECT a. FROM users a",
		len("SELECT a."),
		func(string, int) bool { return false },
		func(name string) { registered = append(registered, name) },
		func(string, *metadatapb.DatabaseSchemaMetadata) {},
		func(string) {},
	)
	require.Zero(t, listCalls)
	require.Empty(t, registered)
	require.Equal(t, []string{"appdb"}, loaded)

	loaded = nil
	populateCatalog(
		context.Background(),
		cfg,
		"SELECT otherdb.t2. FROM otherdb.t2",
		len("SELECT otherdb.t2."),
		func(string, int) bool { return false },
		func(name string) { registered = append(registered, name) },
		func(string, *metadatapb.DatabaseSchemaMetadata) {},
		func(string) {},
	)
	require.Zero(t, listCalls)
	require.Empty(t, registered)
	require.Equal(t, []string{"appdb", "otherdb"}, loaded)

	loaded = nil
	populateCatalog(
		context.Background(),
		cfg,
		"SELECT * FROM ",
		len("SELECT * FROM "),
		func(string, int) bool { return true },
		func(name string) { registered = append(registered, name) },
		func(string, *metadatapb.DatabaseSchemaMetadata) {},
		func(string) {},
	)
	require.Equal(t, 1, listCalls)
	require.Equal(t, []string{"otherdb", "unuseddb"}, registered)
	require.Equal(t, []string{"appdb"}, loaded)

	// Case-sensitive mode loads only the prioritized database per lowercased catalog slot.
	caseSensListCalls := 0
	var caseSensLoaded []string
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "foo",
			ListDatabaseNames: func(context.Context) ([]string, error) {
				caseSensListCalls++
				return []string{"foo", "Foo"}, nil
			},
			LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				if db != "foo" && db != "Foo" {
					return "", nil, false, false
				}
				return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
			},
		},
		"SELECT Foo.t. FROM Foo.t",
		len("SELECT Foo.t."),
		func(string, int) bool { return false },
		func(string) {},
		func(name string, _ *metadatapb.DatabaseSchemaMetadata) { caseSensLoaded = append(caseSensLoaded, name) },
		func(string) {},
	)
	require.Zero(t, caseSensListCalls)
	require.Equal(t, []string{"Foo"}, caseSensLoaded)

	// Case-folded qualifier resolves canonically via lazy enumeration on miss.
	caseFoldListCalls := 0
	var caseFoldLoaded []string
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "appdb",
			ListDatabaseNames: func(context.Context) ([]string, error) {
				caseFoldListCalls++
				return []string{"appdb", "db_b"}, nil
			},
			LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				if db != "appdb" && db != "db_b" {
					return "", nil, false, false
				}
				return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, false, true
			},
		},
		"SELECT DB_B.orders. FROM DB_B.orders",
		len("SELECT DB_B.orders."),
		func(string, int) bool { return false },
		func(string) {},
		func(name string, _ *metadatapb.DatabaseSchemaMetadata) { caseFoldLoaded = append(caseFoldLoaded, name) },
		func(string) {},
	)
	require.Equal(t, 1, caseFoldListCalls)
	require.Equal(t, []string{"appdb", "db_b"}, caseFoldLoaded)

	var multiCaseLoaded []string
	stmtMultiCase := "SELECT * FROM Foo.t JOIN foo.t ON Foo.t.id = foo.t.id WHERE foo.t."
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "appdb",
			LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
			},
		},
		stmtMultiCase,
		len(stmtMultiCase),
		func(string, int) bool { return false },
		func(string) {},
		func(name string, _ *metadatapb.DatabaseSchemaMetadata) {
			multiCaseLoaded = append(multiCaseLoaded, name)
		},
		func(string) {},
	)
	require.Equal(t, []string{"appdb", "foo"}, multiCaseLoaded)

	missingDefaultListCalls := 0
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "missing_db",
			ListDatabaseNames: func(context.Context) ([]string, error) {
				missingDefaultListCalls++
				return []string{"missing_db", "otherdb"}, nil
			},
			LoadDatabase: func(_ context.Context, _ string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				return "", nil, false, false
			},
		},
		"SELECT a. FROM users a",
		len("SELECT a."),
		func(string, int) bool { return false },
		func(string) {},
		func(string, *metadatapb.DatabaseSchemaMetadata) {},
		func(string) {},
	)
	require.Zero(t, missingDefaultListCalls)

	var aliasCaseLoaded []string
	stmtAliasCase := "SELECT x. FROM foo.t y JOIN Foo.t x"
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "appdb",
			LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
			},
		},
		stmtAliasCase,
		len("SELECT x."),
		func(string, int) bool { return false },
		func(string) {},
		func(name string, _ *metadatapb.DatabaseSchemaMetadata) {
			aliasCaseLoaded = append(aliasCaseLoaded, name)
		},
		func(string) {},
	)
	require.Equal(t, []string{"appdb", "Foo"}, aliasCaseLoaded)

	var unqualifiedLoaded []string
	populateCatalog(
		context.Background(),
		SessionConfig{
			DefaultDatabase: "foo",
			LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
				return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
			},
		},
		"SELECT  FROM local JOIN Foo.remote",
		len("SELECT "),
		func(string, int) bool { return false },
		func(string) {},
		func(name string, _ *metadatapb.DatabaseSchemaMetadata) {
			unqualifiedLoaded = append(unqualifiedLoaded, name)
		},
		func(string) {},
	)
	require.Equal(t, []string{"foo"}, unqualifiedLoaded)

	for _, tc := range []struct {
		name string
		stmt string
		pos  int
	}{
		{name: "derived table alias", stmt: "SELECT x. FROM (SELECT * FROM Foo.t) x", pos: len("SELECT x.")},
		{name: "cte alias", stmt: "WITH c AS (SELECT * FROM Foo.t) SELECT x. FROM c x", pos: len("WITH c AS (SELECT * FROM Foo.t) SELECT x.")},
		{name: "cte alias case-insensitive and nested CTE isolation", stmt: "WITH C AS (SELECT * FROM Foo.t) SELECT x., (WITH c AS (SELECT * FROM foo.t) SELECT 1 FROM c) FROM c x", pos: len("WITH C AS (SELECT * FROM Foo.t) SELECT x.")},
		{name: "multi-db CTE projects non-default case-colliding db", stmt: "WITH c AS (SELECT Foo.t.remote FROM foo.local JOIN Foo.t) SELECT x. FROM c x", pos: len("WITH c AS (SELECT Foo.t.remote FROM foo.local JOIN Foo.t) SELECT x.")},
		{name: "multi-db derived table projects non-default case-colliding db", stmt: "SELECT x. FROM (SELECT Foo.t.remote FROM foo.local JOIN Foo.t) x", pos: len("SELECT x.")},
		{name: "multi-db derived table projects aliased star", stmt: "SELECT x. FROM (SELECT y.* FROM foo.t z JOIN Foo.t y) x", pos: len("SELECT x.")},
		{name: "multi-db CTE projects aliased column inside expression and ignores WHERE", stmt: "WITH c AS (SELECT COALESCE(y.remote, 0) FROM foo.local z JOIN Foo.t y WHERE foo.local.id = 1) SELECT x. FROM c x", pos: len("WITH c AS (SELECT COALESCE(y.remote, 0) FROM foo.local z JOIN Foo.t y WHERE foo.local.id = 1) SELECT x.")},
		{name: "multi-db double-parenthesized CTE projects aliased star", stmt: "WITH c AS ((SELECT y.* FROM foo.t z JOIN Foo.t y)) SELECT x. FROM c x", pos: len("WITH c AS ((SELECT y.* FROM foo.t z JOIN Foo.t y)) SELECT x.")},
		{name: "unqualified caret with sole colliding explicit source", stmt: "SELECT  FROM Foo.t", pos: len("SELECT ")},
		{name: "unqualified caret inside subquery ignores outer unqualified table", stmt: "SELECT (SELECT  FROM Foo.t) FROM local", pos: len("SELECT (SELECT ")},
		{name: "unqualified caret with sole CTE backed by colliding source", stmt: "WITH c AS (SELECT * FROM Foo.t) SELECT  FROM c", pos: len("WITH c AS (SELECT * FROM Foo.t) SELECT ")},
		{name: "unqualified caret in INSERT column list of colliding source", stmt: "INSERT INTO Foo.t (", pos: len("INSERT INTO Foo.t (")},
		{name: "constant CTE without backing DB before colliding explicit source", stmt: "WITH c AS (SELECT 1 AS n) SELECT  FROM c JOIN Foo.t", pos: len("WITH c AS (SELECT 1 AS n) SELECT ")},
		{name: "recursive constant CTE without backing DB before colliding explicit source", stmt: "WITH RECURSIVE c AS (SELECT 1 AS n UNION ALL SELECT n+1 FROM c WHERE n < 2) SELECT  FROM c JOIN Foo.t", pos: len("WITH RECURSIVE c AS (SELECT 1 AS n UNION ALL SELECT n+1 FROM c WHERE n < 2) SELECT ")},
		{name: "lateral derived table on colliding explicit source", stmt: "SELECT  FROM LATERAL (SELECT * FROM Foo.t) x", pos: len("SELECT ")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var loadedOrder []string
			populateCatalog(
				context.Background(),
				SessionConfig{
					DefaultDatabase: "foo",
					LoadDatabase: func(_ context.Context, db string) (string, *metadatapb.DatabaseSchemaMetadata, bool, bool) {
						return db, &metadatapb.DatabaseSchemaMetadata{Name: db}, true, true
					},
				},
				tc.stmt,
				tc.pos,
				func(string, int) bool { return false },
				func(string) {},
				func(name string, _ *metadatapb.DatabaseSchemaMetadata) { loadedOrder = append(loadedOrder, name) },
				func(string) {},
			)
			require.Equal(t, []string{"Foo"}, loadedOrder)
		})
	}
}

func TestPostProcessSessionCandidates_PreservesDistinctDefinitionsAndParens(t *testing.T) {
	got := postProcessSessionCandidates(
		"SELECT '😀', ",
		len("SELECT '😀', "),
		false,
		[]Candidate{
			{Type: CandidateColumn, Text: "id", Definition: "appdb.t2.id | int"},
			{Type: CandidateColumn, Text: "id", Definition: "appdb.t1.id | int"},
			{Type: CandidateFunction, Text: "COUNT"},
			{Type: CandidateFunction, Text: "f()"},
			{Type: CandidateFunction, Text: "SUM"},
		},
		func(string, int) bool { return true },
	)

	require.Equal(t, []Candidate{
		{Type: CandidateColumn, Text: "id", Definition: "appdb.t1.id | int"},
		{Type: CandidateColumn, Text: "id", Definition: "appdb.t2.id | int"},
		{Type: CandidateFunction, Text: "COUNT()"},
		{Type: CandidateFunction, Text: "SUM()"},
		{Type: CandidateFunction, Text: "`f()`()"},
	}, got)

	for _, tc := range []struct {
		name string
		stmt string
		pos  int
	}{
		{name: "immediate paren", stmt: "SELECT CO(id) FROM t", pos: len("SELECT CO")},
		{name: "mid-ident before paren", stmt: "SELECT COUNT(id) FROM t", pos: len("SELECT CO")},
		{name: "space before paren", stmt: "SELECT CO (id) FROM t", pos: len("SELECT CO")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotParen := postProcessSessionCandidates(
				tc.stmt,
				tc.pos,
				false,
				[]Candidate{{Type: CandidateFunction, Text: "COUNT"}},
				func(string, int) bool { return true },
			)
			require.Equal(t, []Candidate{{Type: CandidateFunction, Text: "COUNT"}}, gotParen)
		})
	}
}

func TestDollarAndUnicodeIdentifiersConsistentAcrossHelpers(t *testing.T) {
	require.Equal(t, len("sales$archive"), extractIdentPrefixLen("USE sales$archive", len("USE sales$archive")))
	require.Equal(t, len("sales$"), extractIdentPrefixLen("SELECT * FROM sales$", len("SELECT * FROM sales$")))
	require.Equal(t, len("café$2"), extractIdentPrefixLen("USE café$2", len("USE café$2")))
	require.Equal(t, 1, extractIdentPrefixLen("sales$archive.t", len("sales$archive.t")))

	require.True(t, StatementReferencesDatabase("SELECT * FROM sales$archive.t1", "sales$archive"))
	require.False(t, StatementReferencesDatabase("SELECT * FROM sales$archive.t1", "archive"))
	require.False(t, StatementReferencesDatabase("SELECT * FROM café.t1", "fé"))

	require.False(t, identifierNeedsQuoting("sales$archive", nil))
	require.True(t, identifierNeedsQuoting("1$bad", nil))

	require.Equal(t, "a``b", quoteIdentifierIfNeeded("a`b", true, nil))
	require.Equal(t, "`a``b`", quoteIdentifierIfNeeded("a`b", false, nil))

	ident, _, ok := identBeforeDot("SELECT * FROM otherdb.", len("SELECT * FROM otherdb."))
	require.True(t, ok)
	require.Equal(t, "otherdb", ident)
	ident, _, ok = identBeforeDot("SELECT * FROM `a``b` . ", len("SELECT * FROM `a``b` . "))
	require.True(t, ok)
	require.Equal(t, "a`b", ident)
	_, _, ok = identBeforeDot("SELECT * FROM ", len("SELECT * FROM "))
	require.False(t, ok)

	require.True(t, caretAtStatementStart("WITH c AS (SELECT 1) ", len("WITH c AS (SELECT 1) ")))
	require.True(t, caretAtStatementStart("WITH RECURSIVE c1(a) AS (SELECT 1), c2 AS (SELECT 2) UP", len("WITH RECURSIVE c1(a) AS (SELECT 1), c2 AS (SELECT 2) UP")))
	require.False(t, caretAtStatementStart("WITH c AS (SELECT 1), ", len("WITH c AS (SELECT 1), ")))
	require.False(t, caretAtStatementStart("WITH c AS (SELECT 1) SELECT * FROM c WHERE a = 1 FOR ", len("WITH c AS (SELECT 1) SELECT * FROM c WHERE a = 1 FOR ")))
	require.False(t, caretAtStatementStart("WITH c AS (SELECT 1) SELECT (1) ", len("WITH c AS (SELECT 1) SELECT (1) ")))
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

	// 1. Qualified routine/view names resolve as procedures/views, not tables.
	gotCall := CompleteSession(context.Background(), cfg, "CALL remote.", 1, len("CALL remote."))
	require.True(t, containsCandidate(gotCall, "sync_orders", CandidateProcedure), "got %v", gotCall)
	require.False(t, containsCandidate(gotCall, "orders", CandidateTable), "got %v", gotCall)
	require.False(t, containsCandidate(gotCall, "order_view", CandidateView), "got %v", gotCall)

	gotDropView := CompleteSession(context.Background(), cfg, "DROP VIEW remote.", 1, len("DROP VIEW remote."))
	require.True(t, containsCandidate(gotDropView, "order_view", CandidateView), "got %v", gotDropView)
	require.False(t, containsCandidate(gotDropView, "orders", CandidateTable), "got %v", gotDropView)

	// 2. CTE columns resolve from remote tables without falling back to defaultDB.
	gotCTE := CompleteSession(context.Background(), cfg, "WITH c AS (SELECT * FROM remote.orders) SELECT c. FROM c", 1, len("WITH c AS (SELECT * FROM remote.orders) SELECT c."))
	require.True(t, containsCandidate(gotCTE, "order_id", CandidateColumn), "got %v", gotCTE)
	require.True(t, containsCandidate(gotCTE, "amount", CandidateColumn), "got %v", gotCTE)
	require.False(t, containsCandidate(gotCTE, "local_col", CandidateColumn), "got %v", gotCTE)

	// 3. Quoted, Unicode, and $ identifiers complete mid-token.
	gotBacktick := CompleteSession(context.Background(), cfg, "SELECT * FROM `or", 1, len("SELECT * FROM `or"))
	require.True(t, containsCandidate(gotBacktick, "orders", CandidateTable), "got %v", gotBacktick)

	gotUnicode := CompleteSession(context.Background(), cfg, "SELECT * FROM café", 1, len("SELECT * FROM café"))
	require.True(t, containsCandidate(gotUnicode, "café_tbl", CandidateTable), "got %v", gotUnicode)

	gotDollar := CompleteSession(context.Background(), cfg, "SELECT * FROM sales$a", 1, len("SELECT * FROM sales$a"))
	require.True(t, containsCandidate(gotDollar, "sales$archive", CandidateTable), "got %v", gotDollar)

	// 4. Case-sensitive databases never collapse and missing case variant stays empty.
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

	// 5. Comment between qualifier and dot still resolves target database.
	gotCommentDot := CompleteSession(context.Background(), cfg, "SELECT * FROM remote /* note */ .", 1, len("SELECT * FROM remote /* note */ ."))
	require.True(t, containsCandidate(gotCommentDot, "orders", CandidateTable), "got %v", gotCommentDot)
	require.False(t, containsCandidate(gotCommentDot, "café_tbl", CandidateTable), "got %v", gotCommentDot)

	// 6. Comma after JOIN USING(...) or ON ... still loads subsequent database.
	gotUsingComma := CompleteSession(context.Background(), cfg, "SELECT * FROM t JOIN u USING (id), remote.", 1, len("SELECT * FROM t JOIN u USING (id), remote."))
	require.True(t, containsCandidate(gotUsingComma, "orders", CandidateTable), "got %v", gotUsingComma)

	gotOnComma := CompleteSession(context.Background(), cfg, "SELECT * FROM t JOIN u ON t.id=u.id, remote.", 1, len("SELECT * FROM t JOIN u ON t.id=u.id, remote."))
	require.True(t, containsCandidate(gotOnComma, "orders", CandidateTable), "got %v", gotOnComma)

	// 7. Empty DefaultDatabase must not adopt remote DB as current default.
	noDefCfg := cfg
	noDefCfg.DefaultDatabase = ""
	gotNoDef := CompleteSession(context.Background(), noDefCfg, "SELECT * FROM  JOIN remote.orders", 1, len("SELECT * FROM "))
	require.False(t, containsCandidate(gotNoDef, "orders", CandidateTable), "got %v", gotNoDef)

	// 8. Failed LoadDatabase lookup must not invent a database suggestion.
	gotMissing := CompleteSession(context.Background(), cfg, "SELECT * FROM  JOIN missing.t", 1, len("SELECT * FROM "))
	require.True(t, containsCandidate(gotMissing, "remote", CandidateDatabase), "got %v", gotMissing)
	require.False(t, containsCandidate(gotMissing, "missing", CandidateDatabase), "got %v", gotMissing)

	// 9. Stored function named like reserved built-in keyword quotes in DDL.
	gotDropLeft := CompleteSession(context.Background(), cfg, "DROP FUNCTION LEF", 1, len("DROP FUNCTION LEF"))
	require.True(t, containsCandidate(gotDropLeft, "`LEFT`", CandidateFunction), "got %v", gotDropLeft)
	require.False(t, containsCandidate(gotDropLeft, "LEFT", CandidateFunction), "got %v", gotDropLeft)

	// 10. Nested CTE still filters write keywords in read-only mode.
	roCfg := cfg
	roCfg.ReadOnly = true
	gotNestedRO := CompleteSession(context.Background(), roCfg, "WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) ", 1, len("WITH c AS (WITH d AS (SELECT 1) SELECT * FROM d) "))
	require.True(t, containsCandidate(gotNestedRO, "SELECT", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "UPDATE", CandidateKeyword), "got %v", gotNestedRO)
	require.False(t, containsCandidate(gotNestedRO, "DELETE", CandidateKeyword), "got %v", gotNestedRO)

	// 11. Comment before existing '(' suppresses extra '()'.
	gotCommentParen := CompleteSession(context.Background(), cfg, "SELECT my /* note */ ()", 1, len("SELECT my"))
	require.True(t, containsCandidate(gotCommentParen, "`my-func`", CandidateFunction), "got %v", gotCommentParen)
	require.False(t, containsCandidate(gotCommentParen, "`my-func`()", CandidateFunction), "got %v", gotCommentParen)

	// 12. StatementReferencesDatabase respects quoted identifier token boundaries.
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
}
