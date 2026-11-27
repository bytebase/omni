package analysis

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1;\nGO\nSELECT 2;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || !segments[1].Empty || segments[2].Text != "\nSELECT 2;" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELECT 2;" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(""); err != nil || parsed != nil {
		t.Fatalf("Parse of an empty script = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorKeepsTheSegmentsColumn(t *testing.T) {
	_, err := Engine.Parse("SELECT 1;\nSELECT 2; SELECT * FROM;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 15}) || syntaxErr.Message != `Syntax error at line 2: syntax error at or near ";"` || syntaxErr.RawMessage != `syntax error at or near ";"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1; SELECT * FROM;")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 1, Column: 24}) || diagnostics[0].Message != `syntax error at or near ";"` {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateQueryAndLimit(t *testing.T) {
	if valid, allQuery, err := Engine.ValidateQuery("SET NOCOUNT ON; SELECT 1;"); err != nil || !valid || allQuery {
		t.Fatalf("ValidateQuery = %v %v %v", valid, allQuery, err)
	}
	var parseErr *parser.ParseError
	if _, _, err := Engine.ValidateQuery("SELECT FROM ;"); !errors.As(err, &parseErr) {
		t.Fatalf("ValidateQuery err = %#v, want the parser's own error", err)
	}
	if got := Engine.ResultLimit("SELECT * FROM dbo.t;", 10, ""); got != "SELECT TOP 10 * FROM dbo.t;" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT FROM ;", 10, ""); got != "WITH result AS (SELECT FROM) SELECT TOP 10 * FROM result;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
}

func mssqlSpanCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{
		{Name: "dbo", Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}}},
		{Name: "s1", Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}}},
	}}
	return analysis.Catalog{
		ListDatabaseNames: func(context.Context) ([]string, error) { return []string{"db"}, nil },
		GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
			if name != "db" {
				return nil, nil
			}
			return &analysis.Database{Name: name, Metadata: db}, nil
		},
	}
}

func mssqlQuerySpans(t *testing.T, schema string, texts ...string) []*analysis.ResolvedSpan {
	t.Helper()
	statements := make([]analysis.QuerySpanStatement, 0, len(texts))
	for _, text := range texts {
		statements = append(statements, analysis.QuerySpanStatement{Text: text})
	}
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", DefaultSchema: schema, Catalog: mssqlSpanCatalog()}, statements)
	if err != nil {
		t.Fatal(err)
	}
	return spans
}

func TestEngineQuerySpansResolveInTheDefaultSchema(t *testing.T) {
	for schema, want := range map[string]string{"": "dbo", "s1": "s1"} {
		spans := mssqlQuerySpans(t, schema, "SELECT id FROM t")
		if len(spans[0].Results) != 1 || !spans[0].Results[0].SourceColumns[analysis.ColumnResource{Database: "db", Schema: want, Table: "t", Column: "id"}] {
			t.Fatalf("schema %q: spans = %+v", schema, spans[0])
		}
	}
}

func TestEngineQuerySpansCarryTheBatchStateAcrossStatements(t *testing.T) {
	spans := mssqlQuerySpans(t, "", "CREATE TABLE #t (c1 INT);", "SELECT c1 FROM #t;", "SET SHOWPLAN_XML ON;", "SELECT id FROM t;", "SET SHOWPLAN_XML OFF;", "SET STATISTICS XML ON;", "SELECT id FROM t;")
	if len(spans[1].Results) != 1 || spans[1].Results[0].Name != "c1" {
		t.Fatalf("temp table read = %+v", spans[1])
	}
	if spans[3].Type != Explain || spans[6].Type != Select || !spans[6].ExplainAnalyze {
		t.Fatalf("plan output: %+v, %+v", spans[3], spans[6])
	}
	if spans := mssqlQuerySpans(t, "", "SELECT id FROM t;"); spans[0].Type != Select || spans[0].ExplainAnalyze {
		t.Fatalf("a new script starts with no plan option on: %+v", spans[0])
	}
}

func TestEngineQuerySpansReturnTheParsersError(t *testing.T) {
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", Catalog: mssqlSpanCatalog()}, []analysis.QuerySpanStatement{{Text: "SELECT * FROM", Start: analysis.Position{Line: 3, Column: 1}}})
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) || err != error(parseErr) {
		t.Fatalf("err = %T %v", err, err)
	}
}

func mssqlChangeStatements(t *testing.T, script string) []analysis.ParsedStatement {
	t.Helper()
	statements, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	return statements
}

func TestEngineExtractChangesResolveInTheDefaultDatabaseAndSchema(t *testing.T) {
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{DefaultDatabase: "db", DefaultSchema: "s1"}, mssqlChangeStatements(t, "UPDATE t SET a = 1;\nGO\nINSERT INTO other.dbo.t1 VALUES (1), (2);"))
	if err != nil {
		t.Fatal(err)
	}
	want := []analysis.ChangedResource{{Database: "db", Schema: "s1", Table: "t"}, {Database: "other", Schema: "dbo", Table: "t1"}}
	if !reflect.DeepEqual(changes.Resources, want) || changes.DMLCount != 1 || changes.InsertCount != 2 {
		t.Fatalf("changes = %+v", changes)
	}
	if !reflect.DeepEqual(changes.DMLStatements, []analysis.DMLStatement{{Text: "UPDATE t SET a = 1;"}}) {
		t.Fatalf("DML statements = %+v", changes.DMLStatements)
	}
}

func TestEngineExtractChangesSkipANilNodeAndRejectAForeignOne(t *testing.T) {
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{}, []analysis.ParsedStatement{{Node: nil}})
	if err != nil || len(changes.Resources) != 0 {
		t.Fatalf("nil node: changes = %+v, err = %v", changes, err)
	}
	if _, err := Engine.ExtractChanges(analysis.ChangeRequest{}, []analysis.ParsedStatement{{Node: "not a node"}}); err == nil || err.Error() != "expected AST for MSSQL" {
		t.Fatalf("foreign node: err = %v", err)
	}
}

func mssqlBackupCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "a", Type: "int"}},
			Indexes: []*metadata.IndexMetadata{{Name: "PK_T", Primary: true, Unique: true, Expressions: []string{"id"}}},
		}},
	}}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: db}, nil
	}}
}

func mssqlRestoreRequest(script string, backup analysis.Backup) analysis.RestoreRequest {
	return analysis.RestoreRequest{
		Script: script, SourceDatabase: "db", SourceTable: backup.SourceTable, BackupDatabase: "bbdataarchive", BackupTable: backup.BackupTable,
		Start: backup.Start, End: backup.End, Catalog: mssqlBackupCatalog(),
	}
}

// The restore quotes every UPDATE and DELETE in the range, on any table.
func TestEngineBackupsRoundTripThroughRestoreSQL(t *testing.T) {
	script := "UPDATE t SET a = 1 WHERE id = 1\nGO\nDELETE FROM t1;\nUPDATE t SET a = 2 WHERE id = 2;"
	backups, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db", BackupDatabase: "bbdataarchive", TablePrefix: "p"}, script)
	if err != nil {
		t.Fatal(err)
	}
	if len(backups) != 2 || backups[0].SourceTable != "t" || backups[0].SourceSchema != "dbo" ||
		backups[0].Start != (analysis.Position{Line: 1, Column: 0}) || backups[0].End != (analysis.Position{Line: 4, Column: 31}) ||
		backups[1].SourceTable != "t1" || backups[1].Start != (analysis.Position{Line: 3, Column: 0}) {
		t.Fatalf("backups = %+v", backups)
	}
	sql, err := Engine.RestoreSQL(context.Background(), mssqlRestoreRequest(script, backups[0]))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sql, "/*\nOriginal SQL:\nUPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t1;\nUPDATE t SET a = 2 WHERE id = 2;\n*/\n") ||
		!strings.Contains(sql, "[db].[dbo].[t]") || !strings.Contains(sql, "[bbdataarchive].[dbo].["+backups[0].BackupTable+"]") {
		t.Fatalf("sql = %s", sql)
	}
}

func TestEngineRestoreSQLCoversEveryUpdateOfTheTable(t *testing.T) {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "a", Type: "int"}, {Name: "b", Type: "int"}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PK_T", Primary: true, Unique: true, Expressions: []string{"id"}},
				{Name: "UQ_T_B", Unique: true, Expressions: []string{"b"}},
			},
		}},
	}}}
	request := mssqlRestoreRequest("", analysis.Backup{SourceTable: "t", BackupTable: "b"})
	request.Catalog = analysis.Catalog{GetDatabase: func(context.Context, string) (*analysis.Database, error) {
		return &analysis.Database{Name: "db", Metadata: db}, nil
	}}
	// The second UPDATE of t changes the primary key, so UQ_T_B keys the
	// MERGE, and t1's column is no column of t.
	// other.dbo.t is another database's t. A is a, as the case-insensitive
	// collation reads it.
	request.Script = "UPDATE t SET a = 1 WHERE id = 1;\nUPDATE t1 SET z = 1;\nUPDATE other.dbo.t SET z = 1;\nUPDATE x SET id = 5 FROM t AS x WHERE x.id = 2;\nUPDATE t SET A = 2;"
	request.Start, request.End = analysis.Position{Line: 1, Column: 0}, analysis.Position{Line: 5, Column: 19}
	sql, err := Engine.RestoreSQL(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	want := "MERGE INTO [db].[dbo].[t] AS t\nUSING [bbdataarchive].[dbo].[b] AS b\n  ON t.[b] = b.[b]\n" +
		"WHEN MATCHED THEN\n  UPDATE SET t.[a] = b.[a], t.[id] = b.[id]\nWHEN NOT MATCHED THEN\n INSERT ([id], [a], [b]) VALUES (b.[id], b.[a], b.[b]);"
	if !strings.Contains(sql, "UPDATE t SET A = 2;\n*/\n") || !strings.HasSuffix(sql, "\n*/\n"+want) {
		t.Fatalf("sql = %s", sql)
	}
}

func TestEngineRestoreSQLErrors(t *testing.T) {
	script := "UPDATE t SET a = 1;"
	request := mssqlRestoreRequest(script, analysis.Backup{SourceTable: "t", BackupTable: "p_t", Start: analysis.Position{Line: 2}, End: analysis.Position{Line: 3}})
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || err.Error() != "no original SQL" {
		t.Fatalf("err = %v", err)
	}
	request.Start, request.End, request.SourceDatabase = analysis.Position{Line: 1}, analysis.Position{Line: 2}, "other"
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || err.Error() != "database metadata is nil" {
		t.Fatalf("err = %v", err)
	}
	request.Script = "UPDATE t SET a = ("
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || !strings.HasPrefix(err.Error(), "failed to extract single SQL: failed to parse statement: ") {
		t.Fatalf("err = %v", err)
	}
	if _, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db"}, "UPDATE t SET a = ("); err == nil ||
		!strings.HasPrefix(err.Error(), "failed to prepare transformation: failed to parse statement: ") {
		t.Fatalf("err = %v", err)
	}
}

func TestEngineBackupsCompareTheUseDatabaseByTheRequestsCaseSensitivity(t *testing.T) {
	script := "USE DB;\nUPDATE t SET a = 1;"
	if backups, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db", BackupDatabase: "bbdataarchive"}, script); err != nil || len(backups) != 1 {
		t.Fatalf("case-insensitive: backups = %+v, err = %v", backups, err)
	}
	if _, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db", BackupDatabase: "bbdataarchive", CaseSensitive: true}, script); err == nil {
		t.Fatal("a case-sensitive USE of another database must fail")
	}
}
