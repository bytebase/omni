package parser

import (
	"testing"

	"github.com/bytebase/omni/pg/plpgsql/ast"
)

// --------------------------------------------------------------------------
// Section 1.5: DECLARE Section — Variable Declarations
// --------------------------------------------------------------------------

func TestDeclScalar(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantName string
		wantType string
	}{
		{"scalar integer", "DECLARE x integer; BEGIN END", "x", "integer"},
		{"record type", "DECLARE r record; BEGIN END", "r", "record"},
		{"qualified type", "DECLARE x public.my_type; BEGIN END", "x", "public.my_type"},
		{"multiword type", "DECLARE x character varying; BEGIN END", "x", "character varying"},
		{"parameterized type", "DECLARE x character varying(255); BEGIN END", "x", "character varying(255)"},
		{"timestamp type", "DECLARE x timestamp with time zone; BEGIN END", "x", "timestamp with time zone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			d, ok := result.Declarations[0].(*ast.PLDeclare)
			if !ok {
				t.Fatalf("expected PLDeclare, got %T", result.Declarations[0])
			}
			if d.Name != tt.wantName {
				t.Errorf("name: want %q, got %q", tt.wantName, d.Name)
			}
			if d.TypeName != tt.wantType {
				t.Errorf("type: want %q, got %q", tt.wantType, d.TypeName)
			}
		})
	}
}

func TestDeclDefault(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantDefault string
	}{
		{"default :=", "DECLARE x integer := 0; BEGIN END", "0"},
		{"default =", "DECLARE x integer = 0; BEGIN END", "0"},
		{"default DEFAULT", "DECLARE x integer DEFAULT 0; BEGIN END", "0"},
		{"expression default", "DECLARE x int := 1 + 2 * 3; BEGIN END", "1 + 2 * 3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			d := result.Declarations[0].(*ast.PLDeclare)
			if d.Default != tt.wantDefault {
				t.Errorf("default: want %q, got %q", tt.wantDefault, d.Default)
			}
		})
	}
}

func TestDeclConstant(t *testing.T) {
	result := parseOK(t, "DECLARE x CONSTANT integer := 42; BEGIN END")
	if len(result.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
	}
	d := result.Declarations[0].(*ast.PLDeclare)
	if !d.Constant {
		t.Error("expected Constant=true")
	}
	if d.Default != "42" {
		t.Errorf("default: want %q, got %q", "42", d.Default)
	}
}

func TestDeclNotNull(t *testing.T) {
	result := parseOK(t, "DECLARE x integer NOT NULL := 1; BEGIN END")
	if len(result.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
	}
	d := result.Declarations[0].(*ast.PLDeclare)
	if !d.NotNull {
		t.Error("expected NotNull=true")
	}
}

func TestDeclConstantNotNull(t *testing.T) {
	result := parseOK(t, "DECLARE x CONSTANT integer NOT NULL := 1; BEGIN END")
	if len(result.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
	}
	d := result.Declarations[0].(*ast.PLDeclare)
	if !d.Constant {
		t.Error("expected Constant=true")
	}
	if !d.NotNull {
		t.Error("expected NotNull=true")
	}
}

func TestDeclCollate(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		wantCollation string
	}{
		{"collate quoted", `DECLARE x text COLLATE "en_US"; BEGIN END`, "en_US"},
		{"collate with not null", `DECLARE x text COLLATE "C" NOT NULL := 'a'; BEGIN END`, "C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			d := result.Declarations[0].(*ast.PLDeclare)
			if d.Collation != tt.wantCollation {
				t.Errorf("collation: want %q, got %q", tt.wantCollation, d.Collation)
			}
		})
	}
}

func TestDeclTypeRef(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantType string
	}{
		{"%TYPE", "DECLARE x my_table.my_column%TYPE; BEGIN END", "my_table.my_column%TYPE"},
		{"%ROWTYPE", "DECLARE r my_table%ROWTYPE; BEGIN END", "my_table%ROWTYPE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			d := result.Declarations[0].(*ast.PLDeclare)
			if d.TypeName != tt.wantType {
				t.Errorf("type: want %q, got %q", tt.wantType, d.TypeName)
			}
		})
	}
}

func TestDeclRejectsMalformedTypeName(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"bad percent type ref", "DECLARE x my_table%BAD; BEGIN END", "expected TYPE or ROWTYPE"},
		{"unclosed typmod", "DECLARE x numeric(10; BEGIN END", "expected closing delimiter"},
		{"unclosed array", "DECLARE x integer[; BEGIN END", "expected closing delimiter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseErr(t, tt.body, tt.want)
		})
	}
}

func TestDeclMultiple(t *testing.T) {
	result := parseOK(t, "DECLARE x int; y text; z boolean; BEGIN END")
	if len(result.Declarations) != 3 {
		t.Fatalf("expected 3 declarations, got %d", len(result.Declarations))
	}
	names := []string{"x", "y", "z"}
	for i, n := range names {
		d := result.Declarations[i].(*ast.PLDeclare)
		if d.Name != n {
			t.Errorf("decl[%d]: want name %q, got %q", i, n, d.Name)
		}
	}
}

func TestDeclCursor(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantScroll ast.ScrollOption
		wantArgs   int
		wantQuery  string
	}{
		{"basic cursor", "DECLARE c CURSOR FOR SELECT 1; BEGIN END", ast.ScrollNone, 0, "SELECT 1"},
		{"scroll cursor", "DECLARE c SCROLL CURSOR FOR SELECT 1; BEGIN END", ast.ScrollYes, 0, "SELECT 1"},
		{"no scroll cursor", "DECLARE c NO SCROLL CURSOR FOR SELECT 1; BEGIN END", ast.ScrollNo, 0, "SELECT 1"},
		{"cursor with params", "DECLARE c CURSOR (p1 int, p2 text) FOR SELECT p1, p2; BEGIN END", ast.ScrollNone, 2, "SELECT p1, p2"},
		{"cursor with IS", "DECLARE c CURSOR IS SELECT 1; BEGIN END", ast.ScrollNone, 0, "SELECT 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			c, ok := result.Declarations[0].(*ast.PLCursorDecl)
			if !ok {
				t.Fatalf("expected PLCursorDecl, got %T", result.Declarations[0])
			}
			if c.Scroll != tt.wantScroll {
				t.Errorf("scroll: want %d, got %d", tt.wantScroll, c.Scroll)
			}
			if len(c.Args) != tt.wantArgs {
				t.Errorf("args: want %d, got %d", tt.wantArgs, len(c.Args))
			}
			if c.Query != tt.wantQuery {
				t.Errorf("query: want %q, got %q", tt.wantQuery, c.Query)
			}
		})
	}
}

func TestDeclAlias(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantRef string
	}{
		{"alias for $1", "DECLARE a ALIAS FOR $1; BEGIN END", "$1"},
		{"alias for named", "DECLARE a ALIAS FOR param_name; BEGIN END", "param_name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			a, ok := result.Declarations[0].(*ast.PLAliasDecl)
			if !ok {
				t.Fatalf("expected PLAliasDecl, got %T", result.Declarations[0])
			}
			if a.RefName != tt.wantRef {
				t.Errorf("ref: want %q, got %q", tt.wantRef, a.RefName)
			}
		})
	}
}

func TestDeclArray(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantType string
	}{
		{"array type", "DECLARE a integer[]; BEGIN END", "integer[]"},
		{"array with dimension", "DECLARE a integer[10]; BEGIN END", "integer[10]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Declarations) != 1 {
				t.Fatalf("expected 1 declaration, got %d", len(result.Declarations))
			}
			d := result.Declarations[0].(*ast.PLDeclare)
			if d.TypeName != tt.wantType {
				t.Errorf("type: want %q, got %q", tt.wantType, d.TypeName)
			}
		})
	}
}

func TestDeclError(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantContains string
	}{
		{"NOT NULL without default", "DECLARE x integer NOT NULL; BEGIN END", "NOT NULL"},
		{"CONSTANT without default", "DECLARE x CONSTANT integer; BEGIN END", "CONSTANT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseErr(t, tt.body, tt.wantContains)
		})
	}
}
