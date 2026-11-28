package analysis

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestWalkThroughStatementKindString(t *testing.T) {
	tests := []struct {
		kind WalkThroughStatementKind
		want string
	}{
		{WalkThroughStatementOther, "OTHER"},
		{WalkThroughStatementCreateDatabase, "CREATE_DATABASE"},
		{WalkThroughStatementAlterDatabase, "ALTER_DATABASE"},
		{WalkThroughStatementDropDatabase, "DROP_DATABASE"},
		{WalkThroughStatementCreateTable, "CREATE_TABLE"},
		{WalkThroughStatementDropTable, "DROP_TABLE"},
		{WalkThroughStatementTableRefCheck, "TABLE_REF_CHECK"},
		{WalkThroughStatementRenameTable, "RENAME_TABLE"},
		{WalkThroughStatementKind(99), "OTHER"},
	}

	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("WalkThroughStatementKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

// fakeMySQLCatalog accepts every statement unless told to fail.
type fakeMySQLCatalog struct {
	loadErr, runErr error
	scripts         *[]string
}

func (c fakeMySQLCatalog) Load(context.Context, *metadata.DatabaseSchemaMetadata) error {
	return c.loadErr
}

func (c fakeMySQLCatalog) Run(script string) ([]CatalogResult, error) {
	*c.scripts = append(*c.scripts, script)
	if c.runErr != nil && len(*c.scripts) > 1 {
		return nil, c.runErr
	}
	return nil, nil
}

func (fakeMySQLCatalog) DatabaseMetadata(database string) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: database}
}

// fakeMySQLStatements classifies each statement from its text: "ctas t" is a
// CREATE TABLE IF NOT EXISTS t AS SELECT, "create t" a CREATE TABLE, "rename
// a b" a rename of a into database b.
func fakeMySQLStatements(texts ...string) []ParsedStatement {
	statements := make([]ParsedStatement, len(texts))
	for i, text := range texts {
		statements[i] = ParsedStatement{Statement: Statement{Text: text}, Node: text}
	}
	return statements
}

func fakeMySQLClassify(node any) (WalkThroughStatement, bool) {
	fields := strings.Fields(node.(string))
	switch fields[0] {
	case "ctas":
		return WalkThroughStatement{Kind: WalkThroughStatementCreateTable, IsCTAS: true, IfNotExists: true, TableRefs: []TableRef{{Name: fields[1]}}}, true
	case "create":
		return WalkThroughStatement{Kind: WalkThroughStatementCreateTable, TableRefs: []TableRef{{Name: fields[1]}}}, true
	case "rename":
		return WalkThroughStatement{Kind: WalkThroughStatementRenameTable, RenamePairs: []RenamePair{{Old: TableRef{Name: fields[1]}, New: TableRef{Schema: fields[2], Name: fields[1]}}}}, true
	default:
		return WalkThroughStatement{}, false
	}
}

func TestMySQLWalkThroughWithAFakeCatalog(t *testing.T) {
	database := &Database{Name: "test", Metadata: &metadata.DatabaseSchemaMetadata{Name: "test"}}
	run := func(w MySQLWalkThrough, c fakeMySQLCatalog, statements []ParsedStatement) WalkThroughResult {
		w.Classify = fakeMySQLClassify
		w.ErrorNumber = func(error) (int, string, bool) { return 0, "", false }
		w.NewCatalog = func() MySQLCatalog { return c }
		return w.WalkThrough(context.Background(), WalkThroughRequest{Database: database, Script: "script"}, statements)
	}
	ctas := `CREATE TABLE AS statement is used in "ctas t"`

	var scripts []string
	got := run(MySQLWalkThrough{TrackBatchTables: true}, fakeMySQLCatalog{scripts: &scripts}, fakeMySQLStatements("create t", "ctas t", "rename t b`x", "rename u b`x"))
	if got.Finding != nil || got.Database == nil || got.Database.Name != "test" || !got.Database.ObjectCaseSensitive {
		t.Fatalf("tracked batch = %+v", got)
	}
	if want := []string{"SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `b``x`;", "script"}; !reflect.DeepEqual(scripts, want) {
		t.Fatalf("scripts = %q, want %q", scripts, want)
	}

	scripts = nil
	got = run(MySQLWalkThrough{}, fakeMySQLCatalog{scripts: &scripts}, fakeMySQLStatements("create t", "ctas t"))
	if want := (&WalkThroughFinding{Code: WalkThroughCreateTableAs, Title: ctas, Content: ctas, Statement: 1}); !reflect.DeepEqual(got.Finding, want) {
		t.Fatalf("untracked batch = %+v, want %+v", got.Finding, want)
	}

	// A failure of the whole script sorts after a precheck.
	scripts = nil
	got = run(MySQLWalkThrough{}, fakeMySQLCatalog{scripts: &scripts, runErr: errors.New("boom")}, fakeMySQLStatements("ctas t"))
	if got.Finding == nil || got.Finding.Code != WalkThroughCreateTableAs {
		t.Fatalf("precheck and script failure = %+v", got.Finding)
	}
	scripts = nil
	got = run(MySQLWalkThrough{}, fakeMySQLCatalog{scripts: &scripts, runErr: errors.New("boom")}, nil)
	if want := (&WalkThroughFinding{Code: WalkThroughDDLSimulationFailed, Error: true, Title: "boom", Content: "boom", Statement: -1}); !reflect.DeepEqual(got.Finding, want) {
		t.Fatalf("script failure = %+v, want %+v", got.Finding, want)
	}

	scripts = nil
	got = run(MySQLWalkThrough{}, fakeMySQLCatalog{scripts: &scripts, loadErr: errors.New("bad schema")}, nil)
	if want := (&WalkThroughFinding{Code: WalkThroughDDLSimulationFailed, Error: true, Title: "Failed to load schema", Content: "bad schema", Statement: -1}); !reflect.DeepEqual(got.Finding, want) {
		t.Fatalf("load failure = %+v, want %+v", got.Finding, want)
	}
}
