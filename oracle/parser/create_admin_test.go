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

// TestParseCreateMLEModuleSourceVerbatim checks that the module_text of
// CREATE MLE MODULE ... AS is kept verbatim, without lexing JavaScript as SQL.
func TestParseCreateMLEModuleSourceVerbatim(t *testing.T) {
	source := "export function f(a) {\n" +
		"  // don't stop (here; a stray quote and paren\n" +
		"  let s = \"it's /\";\n" +
		"  return a--;\n" +
		"}"
	tests := []string{
		"CREATE MLE MODULE m LANGUAGE JAVASCRIPT AS\n" + source,
		"CREATE OR REPLACE MLE MODULE s.m LANGUAGE JAVASCRIPT VERSION '1.0' AS " + source + "\n",
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
	ParseShouldFail(t, "CREATE MLE MODULE m LANGUAGE JAVASCRIPT AS")
}

func TestParseCreateJavaSourceRequiresSource(t *testing.T) {
	ParseShouldFail(t, "CREATE JAVA SOURCE NAMED q AS")
	ParseShouldFail(t, "CREATE JAVA SOURCE NAMED q AS \n  ")
}

// TestParseCreateJavaHead pins the CREATE JAVA head to Oracle 23ai. The image
// has no Java VM, so a statement whose head Oracle accepts fails later with
// ORA-29538 (or ORA-65021 for SHARING outside an application container); a
// rejected head reports a syntax error first, quoted keywords included.
func TestParseCreateJavaHead(t *testing.T) {
	accept := []string{
		"CREATE JAVA SOURCE NAMED j AS class X {}",
		"CREATE JAVA CLASS USING BFILE (d, 'x.class')",
		"CREATE JAVA RESOURCE NAMED j USING BLOB (SELECT NULL FROM dual)",
		"CREATE JAVA SOURCE NAMED j USING CLOB (SELECT 'x' FROM dual)",
		"CREATE JAVA SOURCE NAMED j USING CLOB SELECT 'x' FROM dual",
		"CREATE JAVA CLASS SCHEMA p1 USING BFILE (d, 'x.class')",
		`CREATE JAVA CLASS SCHEMA "P1" USING BFILE (d, 'x.class')`,
		"CREATE JAVA SOURCE NAMED j SHARING = METADATA AS class X {}",
		"CREATE JAVA SOURCE NAMED j AUTHID CURRENT_USER AS class X {}",
		"CREATE JAVA SOURCE NAMED j RESOLVER ((* public)(* -)) AS class X {}",
		"CREATE JAVA RESOURCE NAMED j USING 'key'",
		"CREATE OR REPLACE AND COMPILE NOFORCE JAVA SOURCE NAMED j AS class X {}",
		"CREATE JAVA SOURCE NAMED j USING BFILE (d, 'x.java')",
		"CREATE JAVA CLASS USING CLOB (SELECT 'x' FROM dual)",
		// IF NOT EXISTS follows the kind.
		"CREATE JAVA SOURCE IF NOT EXISTS NAMED j AS class X {}",
		"CREATE JAVA SOURCE IF NOT EXISTS NAMED p1.j AS class X {}",
		"CREATE AND COMPILE JAVA SOURCE IF NOT EXISTS NAMED j AS class X {}",
		"CREATE NOFORCE JAVA SOURCE IF NOT EXISTS NAMED j AS class X {}",
		"CREATE JAVA RESOURCE IF NOT EXISTS NAMED j USING 'key'",
		"CREATE JAVA CLASS IF NOT EXISTS USING BFILE (d, 'x.class')",
		"CREATE JAVA CLASS IF NOT EXISTS SCHEMA p1 USING BFILE (d, 'x.class')",
		// The subquery may omit SELECT, bare or in parentheses.
		"CREATE JAVA SOURCE NAMED j USING CLOB 'x' FROM dual",
		"CREATE JAVA SOURCE NAMED j USING CLOB ('x' FROM dual)",
		"CREATE JAVA RESOURCE NAMED j USING BLOB NULL FROM dual WHERE 1 = 1",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
	reject := []string{
		`CREATE JAVA "SOURCE" NAMED j AS class X {}`,
		`CREATE JAVA SOURCE "NAMED" j AS class X {}`,
		"CREATE JAVA SOURCE NAMED j garbage AS class X {}",
		"CREATE JAVA SOURCE NAMED j",
		"CREATE JAVA SOURCE j AS class X {}",
		"CREATE JAVA CLASS NAMED j USING BFILE (d, 'x.class')",
		"CREATE JAVA NAMED j AS class X {}",
		"CREATE JAVA SOURCE NAMED j SHARING METADATA AS class X {}",
		"CREATE JAVA SOURCE NAMED j AUTHID nobody AS class X {}",
		"CREATE JAVA CLASS USING BFILE (d)",
		"CREATE JAVA CLASS USING BFILE (d, 'x', 'y')",
		"CREATE JAVA CLASS USING BFILE (SELECT BFILENAME('D', 'x') FROM dual)",
		"CREATE JAVA SOURCE NAMED j RESOLVER ((* public)) AUTHID CURRENT_USER AS class X {}",
		"CREATE JAVA SOURCE NAMED j AUTHID DEFINER SHARING = NONE AS class X {}",
		"CREATE AND RESOLVE TABLE t (a NUMBER)",
		"CREATE AND JAVA SOURCE NAMED j AS class X {}",
		// IF NOT EXISTS elsewhere: ORA-00905 after JAVA, ORA-00901 before
		// it, ORA-00922 after SCHEMA, after NAMED, or repeated.
		"CREATE JAVA IF NOT EXISTS SOURCE NAMED j AS class X {}",
		"CREATE IF NOT EXISTS JAVA SOURCE NAMED j AS class X {}",
		"CREATE IF NOT EXISTS AND COMPILE JAVA SOURCE NAMED j AS class X {}",
		"CREATE JAVA CLASS SCHEMA p1 IF NOT EXISTS USING BFILE (d, 'x.class')",
		"CREATE JAVA SOURCE NAMED IF NOT EXISTS j AS class X {}",
		"CREATE JAVA SOURCE IF NOT EXISTS IF NOT EXISTS NAMED j AS class X {}",
		// ORA-11543: any other spelling; ORA-11541: beside OR REPLACE.
		"CREATE JAVA SOURCE IF EXISTS NAMED j AS class X {}",
		"CREATE JAVA SOURCE IF NOT NAMED j AS class X {}",
		`CREATE JAVA SOURCE IF NOT "EXISTS" NAMED j AS class X {}`,
		"CREATE OR REPLACE JAVA SOURCE IF NOT EXISTS NAMED j AS class X {}",
		"CREATE OR REPLACE AND COMPILE JAVA SOURCE IF NOT EXISTS NAMED j AS class X {}",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
	result := ParseAndCheck(t, "CREATE JAVA CLASS IF NOT EXISTS USING BFILE (d, 'x.class')")
	if stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AdminDDLStmt); !stmt.IfNotExists {
		t.Fatalf("IfNotExists = false, want true")
	}
}
