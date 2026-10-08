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

// TestParseCreateJavaSourceVerbatim checks that the source_char of CREATE
// JAVA ... AS is kept verbatim, without lexing Java as SQL.
func TestParseCreateJavaSourceVerbatim(t *testing.T) {
	source := "public class Q {\n" +
		"  // don't stop (here; a stray quote and paren\n" +
		"  static char c = '\\'';\n" +
		"  static int g(int i) { i--; return i / 2; }\n" +
		"}"
	tests := []string{
		"CREATE JAVA SOURCE NAMED \"S\".\"Q\" AS\n" + source,
		"CREATE OR REPLACE AND COMPILE JAVA SOURCE NAMED q AS " + source + "\n",
		"CREATE OR REPLACE AND RESOLVE NOFORCE JAVA SOURCE NAMED q AUTHID CURRENT_USER AS\n\n" + source,
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AdminDDLStmt)
			var got string
			for _, item := range stmt.Options.Items {
				if opt := item.(*ast.DDLOption); opt.Key == "AS" {
					got = opt.Value
				}
			}
			if got != source {
				t.Fatalf("AS source = %q, want %q", got, source)
			}
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

func TestParseCreateJavaSourceRequiresSource(t *testing.T) {
	ParseShouldFail(t, "CREATE JAVA SOURCE NAMED q AS")
	ParseShouldFail(t, "CREATE JAVA SOURCE NAMED q AS \n  ")
}
