package parser

import (
	"testing"

	ast "github.com/bytebase/omni/oracle/ast"
)

func TestP2AdministerKeyManagementOptions(t *testing.T) {
	result := ParseAndCheck(t, "ADMINISTER KEY MANAGEMENT SET KEY USING TAG 'quarterly_key' IDENTIFIED BY password1 WITH BACKUP USING 'Q1 key rotation' CONTAINER = ALL")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
	if !ok {
		t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
	}
	if stmt.Options == nil {
		t.Fatal("expected options, got nil")
	}

	want := []struct {
		key   string
		value string
	}{
		{"SCOPE", "KEY MANAGEMENT"},
		{"COMMAND", "SET KEY"},
		{"USING TAG", "quarterly_key"},
		{"IDENTIFIED BY", "PASSWORD1"},
		{"WITH BACKUP", ""},
		{"USING", "Q1 key rotation"},
		{"CONTAINER", "ALL"},
	}
	if len(stmt.Options.Items) != len(want) {
		t.Fatalf("options len=%d, want %d: %#v", len(stmt.Options.Items), len(want), stmt.Options.Items)
	}
	for i, expected := range want {
		opt, ok := stmt.Options.Items[i].(*ast.DDLOption)
		if !ok {
			t.Fatalf("option %d: expected *DDLOption, got %T", i, stmt.Options.Items[i])
		}
		if opt.Key != expected.key || opt.Value != expected.value {
			t.Fatalf("option %d: got %q=%q, want %q=%q", i, opt.Key, opt.Value, expected.key, expected.value)
		}
		if opt.Loc.IsUnknown() {
			t.Fatalf("option %d %s has unknown Loc", i, opt.Key)
		}
	}
}

func TestP2AdministerKeyManagementRequiresKeyManagement(t *testing.T) {
	assertParseErrorContains(t,
		"ADMINISTER KEY SET KEY IDENTIFIED BY password1",
		`syntax error at or near "SET"`,
	)
}

func TestP2AdministerKeyManagementLoc(t *testing.T) {
	CheckLocations(t, "ADMINISTER KEY MANAGEMENT SET KEY IDENTIFIED BY password1 WITH BACKUP")
}
