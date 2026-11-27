package analysis

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/parser"
)

func TestEngineSplit(t *testing.T) {
	script := "SELECT 1;\n-- lead\nSELECT 2;\n;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || segments[1].Text != "\n-- lead\nSELECT 2;" || !segments[2].Empty {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
	}
	if segments := Engine.Split(""); segments != nil {
		t.Fatalf("Split(\"\") = %#v", segments)
	}
}

func TestEngineDiagnoseKeepsTheSegmentsColumn(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1; SELEC 2;")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 1, Column: 2}) || diagnostics[0].Severity != analysis.SeverityError || diagnostics[0].Message != `syntax error at or near "SELEC"` {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
	diagnostics = Engine.Diagnose("SELECT 1;\n  SELEC 2;")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 2, Column: 3}) {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateQuery(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1; SELECT 2"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	if readOnly, _, err := Engine.ValidateQuery("DELETE FROM t"); err != nil || readOnly {
		t.Fatalf("ValidateQuery(DELETE) = %v %v", readOnly, err)
	}
	_, _, err := Engine.ValidateQuery("SELECT 1; SELEC 2")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 1, Column: 2}) || syntaxErr.Message != "Syntax error at line 1:2: syntax error at or near \"SELEC\"\nrelated text:  SELEC 2" || syntaxErr.RawMessage != `syntax error at or near "SELEC"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineExplainLimitAndRanges(t *testing.T) {
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN FORMAT=tidb_json SELECT 1",
		analysis.ExplainFormatXML:     "EXPLAIN SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("EXPLAIN ANALYZE SELECT 1", analysis.ExplainFormatDefault); !errors.Is(err, analysis.ErrExplainAnalyzeExecutes) {
		t.Fatalf("Explain(EXPLAIN ANALYZE) err = %v", err)
	}
	if _, err := Engine.Explain("ANALYZE TABLE t", analysis.ExplainFormatDefault); !errors.Is(err, analysis.ErrNotExplainable) {
		t.Fatalf("Explain(ANALYZE TABLE) err = %v", err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, err := Engine.Explain("SELEC 1", analysis.ExplainFormatDefault); !errors.As(err, &syntaxErr) {
		t.Fatalf("Explain(SELEC 1) err = %v", err)
	}

	if got := Engine.ResultLimit("SELECT 1", 10, ""); got != "SELECT 1 LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT 1; SELECT 2;", 10, ""); got != "WITH result AS (SELECT 1; SELECT 2) SELECT * FROM result LIMIT 10;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}

	ranges := Engine.StatementRanges("SELECT 1;\n  SELECT '🎉';")
	if len(ranges) != 2 || ranges[1].Start != (analysis.LSPPosition{Line: 1, Character: 2}) || ranges[1].End != (analysis.LSPPosition{Line: 1, Character: 14}) {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
	if ranges := Engine.StatementRanges("   -- comment\n"); ranges != nil {
		t.Fatalf("StatementRanges of a comment = %#v", ranges)
	}
}

func tidbSpanCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}},
	}}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: "resolved", Metadata: db, ObjectCaseSensitive: true}, nil
	}}
}

func TestEngineQuerySpansResolveInTheDefaultDatabase(t *testing.T) {
	for _, ignoreCase := range []bool{false, true} {
		request := analysis.QuerySpanRequest{DefaultDatabase: "db", DefaultSchema: "s1", IgnoreCaseSensitive: ignoreCase, Catalog: tidbSpanCatalog()}
		spans, err := Engine.QuerySpans(context.Background(), request, []analysis.QuerySpanStatement{{Text: "SELECT id FROM t"}, {Text: "SELECT id FROM T"}})
		if err != nil {
			t.Fatal(err)
		}
		for _, span := range spans {
			if span.NotFoundError != nil || !span.Results[0].SourceColumns[analysis.ColumnResource{Database: "db", Table: "t", Column: "id"}] {
				t.Fatalf("ignoreCase %v: span = %+v", ignoreCase, span)
			}
		}
	}
}

func TestEngineQuerySpansWrapTheParseError(t *testing.T) {
	text := "SELECT * FORM t"
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", Catalog: tidbSpanCatalog()}, []analysis.QuerySpanStatement{{Text: text, Start: analysis.Position{Line: 3, Column: 5}}})
	var parseErr *parser.ParseError
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &parseErr) || errors.As(err, &syntaxErr) || !strings.HasPrefix(err.Error(), "failed to get query span from statement: "+text+": ") {
		t.Fatalf("err = %v", err)
	}
}

func parseTiDBForChanges(t *testing.T, texts ...string) []analysis.ParsedStatement {
	t.Helper()
	var statements []analysis.ParsedStatement
	for _, text := range texts {
		list, err := parser.Parse(text)
		if err != nil {
			t.Fatal(err)
		}
		for _, node := range list.Items {
			statements = append(statements, analysis.ParsedStatement{Statement: analysis.Statement{Text: text}, Node: node})
		}
	}
	return statements
}

func TestEngineExtractChangesInTheDefaultDatabase(t *testing.T) {
	statements := append(parseTiDBForChanges(t, "INSERT INTO t VALUES (1), (2);", "\nUPDATE other.t1 SET a = 1;"), analysis.ParsedStatement{})
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{DefaultDatabase: "db", DefaultSchema: "s1"}, statements)
	if err != nil {
		t.Fatal(err)
	}
	want := &analysis.Changes{
		Resources:     []analysis.ChangedResource{{Database: "db", Table: "t"}, {Database: "other", Table: "t1"}},
		DMLStatements: []analysis.DMLStatement{{Text: "UPDATE other.t1 SET a = 1"}},
		DMLCount:      1,
		InsertCount:   2,
	}
	if !reflect.DeepEqual(changes, want) {
		t.Fatalf("changes = %+v", changes)
	}
}

func TestEngineExtractChangesRejectAForeignNode(t *testing.T) {
	_, err := Engine.ExtractChanges(analysis.ChangeRequest{DefaultDatabase: "db"}, []analysis.ParsedStatement{{Node: "not a node"}})
	if err == nil || err.Error() != "expected TiDB AST" {
		t.Fatalf("err = %v", err)
	}
}

func TestEngineExtractChangesKeepATableWithoutAName(t *testing.T) {
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{DefaultDatabase: "db"}, parseTiDBForChanges(t, "TRUNCATE ``"))
	if err != nil || len(changes.Resources) != 1 || changes.Resources[0] != (analysis.ChangedResource{Database: "db", AffectsData: true}) {
		t.Fatalf("changes = %+v, err = %v", changes, err)
	}
}

func tidbBackupCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "a", Type: "int"}, {Name: "g", Type: "int", Generation: &metadata.GenerationMetadata{Expression: "a + 1"}}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}},
		}},
	}}}
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

func TestEngineBackupsRoundTripThroughRestoreSQL(t *testing.T) {
	script := "UPDATE t SET a = 1 WHERE id = 1;\nUPDATE t SET a = 2 WHERE id = 2;"
	backups, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db", BackupDatabase: "bbdataarchive", TablePrefix: "p", Catalog: tidbBackupCatalog()}, script)
	if err != nil {
		t.Fatal(err)
	}
	// Each statement has its own backup, which ends where the next one starts.
	if len(backups) != 2 || backups[0].Start != (analysis.Position{Line: 1, Column: 1}) || backups[0].End != (analysis.Position{Line: 1, Column: 33}) || backups[1].Start != backups[0].End {
		t.Fatalf("backups = %+v", backups)
	}
	sql, err := Engine.RestoreSQL(context.Background(), analysis.RestoreRequest{
		Script: script, SourceDatabase: "db", SourceTable: "t", BackupDatabase: "bbdataarchive", BackupTable: backups[0].BackupTable,
		Start: backups[0].Start, End: backups[0].End, Catalog: tidbBackupCatalog(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sql, "/*\nOriginal SQL:\nUPDATE t SET a = 1 WHERE id = 1;\n*/\nINSERT INTO `db`.`t` (`id`, `a`) SELECT") {
		t.Fatalf("sql = %q", sql)
	}
}

func TestEngineRestoreSQLReadsTheSourceTablesColumnsFirst(t *testing.T) {
	request := analysis.RestoreRequest{Script: "DELETE FROM t;", SourceDatabase: "missing", SourceTable: "t", Start: analysis.Position{Line: 1, Column: 1}, End: analysis.Position{Line: 1, Column: 15}, Catalog: tidbBackupCatalog()}
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || err.Error() != "failed to get target table columns: database metadata is nil for missing" {
		t.Fatalf("err = %v", err)
	}
	if _, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "missing", Catalog: tidbBackupCatalog()}, request.Script); err == nil ||
		err.Error() != `failed to classify columns: failed to get database metadata for Database "missing"` {
		t.Fatalf("Backups err = %v", err)
	}
}
