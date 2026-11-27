package analysis

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1 FROM DUAL;\n;\nSELECT 2 FROM DUAL;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1 FROM DUAL" || !segments[1].Empty || segments[2].Text != "\nSELECT 2 FROM DUAL" {
		t.Fatalf("Split = %+v", segments)
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELECT 2 FROM DUAL" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T, want the statement inside the RawStmt", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(""); err != nil || parsed != nil {
		t.Fatalf("Parse of an empty script = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorIsPositionedInTheScript(t *testing.T) {
	for _, tc := range []struct {
		script string
		want   analysis.Position
	}{
		// A column on the segment's first line moves by the segment's column.
		{"SELECT 1 FROM DUAL; SELECT * FROM", analysis.Position{Line: 1, Column: 34}},
		{"SELECT 1 FROM DUAL;\n  SELECT * FROM", analysis.Position{Line: 2, Column: 16}},
	} {
		_, err := Engine.Parse(tc.script)
		var syntaxErr *analysis.SyntaxError
		if !errors.As(err, &syntaxErr) {
			t.Fatalf("Parse(%q) err = %v", tc.script, err)
		}
		if syntaxErr.Position != tc.want || syntaxErr.RawMessage != "syntax error at end of input" {
			t.Fatalf("Parse(%q) err = %+v", tc.script, syntaxErr)
		}
	}
	_, err := Engine.Parse("SELECT 1 FROM DUAL;\n  SELECT * FROM")
	if err.Error() != "Syntax error at line 2:16: syntax error at end of input" {
		t.Fatalf("message = %q", err.Error())
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1 FROM DUAL;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1 FROM DUAL;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 2, Column: 14}) || diagnostics[0].Message != "syntax error at end of input" {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateQuery(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1 FROM DUAL"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	_, _, err := Engine.ValidateQuery("SELECT 1 FROM DUAL;\nSELECT * FROM")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) || syntaxErr.Message != "Syntax error at line 2:14: syntax error at end of input" {
		t.Fatalf("ValidateQuery err = %#v", err)
	}
}

func TestEngineResultLimit(t *testing.T) {
	for _, tc := range []struct {
		statement, version, want string
	}{
		{"SELECT * FROM t;", "19.3.0.0.0", "SELECT * FROM t FETCH NEXT 10 ROWS ONLY"},
		{"SELECT * FROM t;", "11.2.0.4.0", "SELECT * FROM (SELECT * FROM t) WHERE ROWNUM <= 10"},
		{"SELECT * FROM t", "", "SELECT * FROM (SELECT * FROM t) WHERE ROWNUM <= 10"},
		// A FETCH NEXT rewrite that fails falls back to ROWNUM.
		{"SELECT * FROM t FETCH FIRST 50 PERCENT ROWS ONLY", "19.3.0.0.0", "SELECT * FROM (SELECT * FROM t FETCH FIRST 50 PERCENT ROWS ONLY) WHERE ROWNUM <= 10"},
	} {
		if got := Engine.ResultLimit(tc.statement, 10, tc.version); got != tc.want {
			t.Fatalf("ResultLimit(%q, %q) = %q", tc.statement, tc.version, got)
		}
	}
}

func oracleChangeStatements(t *testing.T, script string) []analysis.ParsedStatement {
	t.Helper()
	statements, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	return statements
}

func oracleChangeDatabase() *analysis.Database {
	table := &metadata.TableMetadata{Name: "T", Indexes: []*metadata.IndexMetadata{{Name: "T_IDX"}}}
	return &analysis.Database{Name: "ORCL", ObjectCaseSensitive: true, Metadata: &metadata.DatabaseSchemaMetadata{
		Name: "ORCL", Schemas: []*metadata.SchemaMetadata{{Name: "S1", Tables: []*metadata.TableMetadata{table}}},
	}}
}

func TestEngineExtractChangesTakeTheDefaultDatabaseAsTheSchema(t *testing.T) {
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{DefaultDatabase: "S1", DefaultSchema: "S2"}, oracleChangeStatements(t, "UPDATE t SET a = 1;"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(changes.Resources, []analysis.ChangedResource{{Database: "S1", Table: "T"}}) || changes.DMLCount != 1 {
		t.Fatalf("changes = %+v", changes)
	}
	if !reflect.DeepEqual(changes.DMLStatements, []analysis.DMLStatement{{Text: "UPDATE t SET a = 1;"}}) {
		t.Fatalf("DML statements = %+v", changes.DMLStatements)
	}
}

func TestEngineExtractChangesFindTheTableOfADroppedIndex(t *testing.T) {
	request := analysis.ChangeRequest{DefaultDatabase: "S1", Database: oracleChangeDatabase()}
	changes, err := Engine.ExtractChanges(request, oracleChangeStatements(t, `DROP INDEX t_idx; DROP INDEX s1.t_idx; DROP INDEX "t_idx"; DROP INDEX s2.t_idx;`))
	if err != nil {
		t.Fatal(err)
	}
	want := []analysis.ChangedResource{{Database: "S1", Table: "T"}, {Database: "S1", Table: "T"}}
	if !reflect.DeepEqual(changes.Resources, want) {
		t.Fatalf("changes = %+v", changes.Resources)
	}
	request.Database = nil
	if changes, err := Engine.ExtractChanges(request, oracleChangeStatements(t, "DROP INDEX t_idx;")); err != nil || len(changes.Resources) != 0 {
		t.Fatalf("without metadata: changes = %+v, err = %v", changes, err)
	}
}

func TestEngineExtractChangesSkipANilOrForeignNode(t *testing.T) {
	changes, err := Engine.ExtractChanges(analysis.ChangeRequest{}, []analysis.ParsedStatement{{Node: nil}, {Node: "not a node"}})
	if err != nil || len(changes.Resources) != 0 || changes.DMLCount != 0 {
		t.Fatalf("changes = %+v, err = %v", changes, err)
	}
}

func oracleBackupCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "NUMBER"}, {Name: "a", Type: "NUMBER"}},
			Indexes: []*metadata.IndexMetadata{{Name: "t_pk", Primary: true, Unique: true, Expressions: []string{"id"}}},
		}},
	}}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: db, ObjectCaseSensitive: true}, nil
	}}
}

func oracleBackups(t *testing.T, version any, script string) []analysis.Backup {
	t.Helper()
	backups, err := Engine.Backups(context.Background(), analysis.BackupRequest{SourceDatabase: "db", BackupDatabase: "BBDATAARCHIVE", TablePrefix: "p", Version: version}, script)
	if err != nil {
		t.Fatal(err)
	}
	return backups
}

func oracleRestoreRequest(script string, backup analysis.Backup) analysis.RestoreRequest {
	return analysis.RestoreRequest{
		Script: script, SourceDatabase: "db", SourceTable: backup.SourceTable, BackupDatabase: "BBDATAARCHIVE", BackupTable: backup.BackupTable,
		Start: backup.Start, End: backup.End, Catalog: oracleBackupCatalog(),
	}
}

// The restore quotes only the range's statements on the source table.
func TestEngineBackupsRoundTripThroughRestoreSQL(t *testing.T) {
	script := "UPDATE \"t\" SET \"a\" = 1 WHERE \"id\" = 1;\nDELETE FROM \"t1\";\nUPDATE \"t\" SET \"a\" = 2;"
	backups := oracleBackups(t, nil, script)
	if len(backups) != 2 || backups[0].SourceTable != "t" || backups[0].Start != (analysis.Position{Line: 1, Column: 0}) || backups[0].End != (analysis.Position{Line: 3, Column: 21}) {
		t.Fatalf("backups = %+v", backups)
	}
	sql, err := Engine.RestoreSQL(context.Background(), oracleRestoreRequest(script, backups[0]))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sql, "/*\nOriginal SQL:\nUPDATE \"t\" SET \"a\" = 1 WHERE \"id\" = 1\nUPDATE \"t\" SET \"a\" = 2\n*/\n") ||
		!strings.Contains(sql, `MERGE INTO "db"."t"`) || !strings.Contains(sql, `"BBDATAARCHIVE"."`+backups[0].BackupTable+`"`) {
		t.Fatalf("sql = %s", sql)
	}
}

func TestEngineBackupsEndAtTheLastCharacterOrAfterAMultiByteOne(t *testing.T) {
	for script, want := range map[string]analysis.Position{
		`DELETE FROM "t" WHERE 2 = a;`: {Line: 1, Column: 26},
		`DELETE FROM "t" WHERE 2 = 名;`: {Line: 1, Column: 27},
	} {
		if backups := oracleBackups(t, nil, script); len(backups) != 1 || backups[0].End != want {
			t.Fatalf("%s: backups = %+v", script, backups)
		}
	}
}

func TestEngineBackupsCapTheTableNameByTheVersion(t *testing.T) {
	script := `UPDATE "` + strings.Repeat("a", 40) + `" SET "a" = 1;`
	for _, version := range []any{nil, &Version{First: 12, Second: 1}, Version{First: 19}} {
		if backups := oracleBackups(t, version, script); len(backups) != 1 || len(backups[0].BackupTable) > 30 {
			t.Fatalf("version %v: backups = %+v", version, backups)
		}
	}
	if backups := oracleBackups(t, &Version{First: 12, Second: 2}, script); len(backups) != 1 || len(backups[0].BackupTable) <= 30 {
		t.Fatalf("12.2: backups = %+v", backups)
	}
}

func TestEngineRestoreSQLErrors(t *testing.T) {
	script := `UPDATE "t" SET "a" = 1;`
	request := oracleRestoreRequest(script, analysis.Backup{SourceTable: "t1", BackupTable: "p_t", Start: analysis.Position{Line: 1}, End: analysis.Position{Line: 2}})
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || err.Error() != "no DML statement found in extracted SQL" {
		t.Fatalf("err = %v", err)
	}
	request.SourceTable, request.SourceDatabase = "t", "other"
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || err.Error() != `database metadata for database "other" not found (database not synced)` {
		t.Fatalf("err = %v", err)
	}
	request.Script = `UPDATE "t" SET "a" = (`
	if _, err := Engine.RestoreSQL(context.Background(), request); err == nil || !strings.HasPrefix(err.Error(), "failed to extract single SQL: failed to prepare transformation: ") {
		t.Fatalf("err = %v", err)
	}
}
