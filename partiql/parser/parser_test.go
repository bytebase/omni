package parser

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// update regenerates golden files. Run with -update to refresh.
var update = flag.Bool("update", false, "update golden files")

// TestParser_Machinery verifies the low-level token buffer helpers
// without any expression parsing. Each case constructs a Parser
// directly and drives the helpers through their state transitions.
func TestParser_Machinery(t *testing.T) {
	t.Run("new_parser_primes_cur", func(t *testing.T) {
		p := NewParser("foo")
		if p.cur.Type != tokIDENT {
			t.Errorf("cur.Type = %d, want tokIDENT", p.cur.Type)
		}
		if p.cur.Str != "foo" {
			t.Errorf("cur.Str = %q, want %q", p.cur.Str, "foo")
		}
	})

	t.Run("advance_moves_cur_forward", func(t *testing.T) {
		p := NewParser("foo bar")
		if p.peek().Str != "foo" {
			t.Errorf("first peek() = %q, want foo", p.peek().Str)
		}
		p.advance()
		if p.peek().Str != "bar" {
			t.Errorf("second peek() = %q, want bar", p.peek().Str)
		}
		if p.prev.Str != "foo" {
			t.Errorf("prev.Str = %q, want foo", p.prev.Str)
		}
	})

	t.Run("peek_next_lookahead", func(t *testing.T) {
		p := NewParser("foo bar baz")
		if p.peek().Str != "foo" {
			t.Errorf("peek() = %q, want foo", p.peek().Str)
		}
		if p.peekNext().Str != "bar" {
			t.Errorf("peekNext() = %q, want bar", p.peekNext().Str)
		}
		if p.peek().Str != "foo" {
			t.Errorf("peek() after peekNext = %q, want foo", p.peek().Str)
		}
		if p.peekNext().Str != "bar" {
			t.Errorf("second peekNext() = %q, want bar", p.peekNext().Str)
		}
		p.advance()
		if p.peek().Str != "bar" {
			t.Errorf("peek() after advance = %q, want bar", p.peek().Str)
		}
		if p.peekNext().Str != "baz" {
			t.Errorf("peekNext() after advance = %q, want baz", p.peekNext().Str)
		}
	})

	t.Run("match_consumes_on_hit", func(t *testing.T) {
		p := NewParser("AND OR")
		if !p.match(tokAND) {
			t.Errorf("match(tokAND) returned false")
		}
		if p.cur.Type != tokOR {
			t.Errorf("cur.Type after match = %d, want tokOR", p.cur.Type)
		}
		if p.match(tokAND) {
			t.Errorf("match(tokAND) returned true for non-matching token")
		}
		if p.match(tokOR, tokAND) {
			if p.cur.Type != tokEOF {
				t.Errorf("cur.Type after match = %d, want tokEOF", p.cur.Type)
			}
		} else {
			t.Errorf("match(tokOR, tokAND) returned false when tokOR was cur")
		}
	})

	t.Run("expect_returns_error_on_miss", func(t *testing.T) {
		p := NewParser("foo")
		_, err := p.expect(tokCOMMA)
		if err == nil {
			t.Fatal("expect(tokCOMMA) returned nil error on non-matching token")
		}
		perr, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("error type = %T, want *ParseError", err)
		}
		if !strings.Contains(perr.Message, "expected") {
			t.Errorf("error message = %q, want to contain 'expected'", perr.Message)
		}
		if perr.Loc.Start != 0 {
			t.Errorf("error Loc.Start = %d, want 0", perr.Loc.Start)
		}
	})

	t.Run("expect_consumes_on_hit", func(t *testing.T) {
		p := NewParser("foo, bar")
		tok, err := p.expect(tokIDENT)
		if err != nil {
			t.Fatalf("expect(tokIDENT) error: %v", err)
		}
		if tok.Str != "foo" {
			t.Errorf("expect returned %q, want foo", tok.Str)
		}
		if p.cur.Type != tokCOMMA {
			t.Errorf("cur.Type after expect = %d, want tokCOMMA", p.cur.Type)
		}
	})

	t.Run("lexer_error_propagation", func(t *testing.T) {
		p := NewParser("'unclosed")
		if p.cur.Type != tokEOF {
			t.Errorf("cur.Type = %d, want tokEOF for unterminated string", p.cur.Type)
		}
		err := p.checkLexerErr()
		if err == nil {
			t.Fatal("checkLexerErr returned nil, want lexer error")
		}
		perr, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("error type = %T, want *ParseError", err)
		}
		if !strings.Contains(perr.Message, "unterminated string literal") {
			t.Errorf("error message = %q, want to contain 'unterminated string literal'", perr.Message)
		}
	})

	t.Run("parse_symbol_primitive_bare", func(t *testing.T) {
		p := NewParser("foo")
		name, caseSensitive, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			t.Fatalf("parseSymbolPrimitive error: %v", err)
		}
		if name != "foo" {
			t.Errorf("name = %q, want foo", name)
		}
		if caseSensitive {
			t.Error("caseSensitive = true, want false for bare ident")
		}
		if loc.Start != 0 || loc.End != 3 {
			t.Errorf("loc = %+v, want {0, 3}", loc)
		}
	})

	t.Run("parse_symbol_primitive_quoted", func(t *testing.T) {
		p := NewParser(`"Foo"`)
		name, caseSensitive, _, err := p.parseSymbolPrimitive()
		if err != nil {
			t.Fatalf("parseSymbolPrimitive error: %v", err)
		}
		if name != "Foo" {
			t.Errorf("name = %q, want Foo", name)
		}
		if !caseSensitive {
			t.Error("caseSensitive = false, want true for quoted ident")
		}
	})

	t.Run("parse_var_ref_bare", func(t *testing.T) {
		p := NewParser("foo")
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v, ok := expr.(*ast.VarRef)
		if !ok {
			t.Fatalf("parseVarRef returned %T, want *ast.VarRef", expr)
		}
		if v.Name != "foo" || v.AtPrefixed || v.CaseSensitive {
			t.Errorf("VarRef = %+v, want {Name:foo AtPrefixed:false CaseSensitive:false}", v)
		}
	})

	t.Run("parse_var_ref_at_prefixed", func(t *testing.T) {
		p := NewParser("@x")
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v, ok := expr.(*ast.VarRef)
		if !ok {
			t.Fatalf("parseVarRef returned %T, want *ast.VarRef", expr)
		}
		if v.Name != "x" || !v.AtPrefixed {
			t.Errorf("VarRef = %+v, want {Name:x AtPrefixed:true}", v)
		}
	})

	t.Run("parse_var_ref_quoted", func(t *testing.T) {
		p := NewParser(`"Foo"`)
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v, ok := expr.(*ast.VarRef)
		if !ok {
			t.Fatalf("parseVarRef returned %T, want *ast.VarRef", expr)
		}
		if v.Name != "Foo" || !v.CaseSensitive {
			t.Errorf("VarRef = %+v, want {Name:Foo CaseSensitive:true}", v)
		}
	})
}

// TestParseType exhaustively tests parseType across the 30+ type
// forms in PartiQLParser.g4's `type` rule. Uses direct table-driven
// assertions because TypeRef is a leaf node (no child recursion) and
// inline expected values are clearer than filesystem goldens for
// exhaustive enumeration.
func TestParseType(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantName   string
		wantArgs   []int
		wantWithTZ bool
	}{
		// Atomic types.
		{"null", "NULL", "NULL", nil, false},
		{"bool", "BOOL", "BOOL", nil, false},
		{"boolean", "BOOLEAN", "BOOLEAN", nil, false},
		{"smallint", "SMALLINT", "SMALLINT", nil, false},
		{"int2", "INT2", "INT2", nil, false},
		{"integer2", "INTEGER2", "INTEGER2", nil, false},
		{"int4", "INT4", "INT4", nil, false},
		{"integer4", "INTEGER4", "INTEGER4", nil, false},
		{"int8", "INT8", "INT8", nil, false},
		{"integer8", "INTEGER8", "INTEGER8", nil, false},
		{"int", "INT", "INT", nil, false},
		{"integer", "INTEGER", "INTEGER", nil, false},
		{"bigint", "BIGINT", "BIGINT", nil, false},
		{"real", "REAL", "REAL", nil, false},
		{"timestamp", "TIMESTAMP", "TIMESTAMP", nil, false},
		{"missing", "MISSING", "MISSING", nil, false},
		{"string", "STRING", "STRING", nil, false},
		{"symbol", "SYMBOL", "SYMBOL", nil, false},
		{"blob", "BLOB", "BLOB", nil, false},
		{"clob", "CLOB", "CLOB", nil, false},
		{"date", "DATE", "DATE", nil, false},
		{"struct", "STRUCT", "STRUCT", nil, false},
		{"tuple", "TUPLE", "TUPLE", nil, false},
		{"list", "LIST", "LIST", nil, false},
		{"sexp", "SEXP", "SEXP", nil, false},
		{"bag", "BAG", "BAG", nil, false},
		{"any", "ANY", "ANY", nil, false},

		// Two-token DOUBLE PRECISION.
		{"double_precision", "DOUBLE PRECISION", "DOUBLE PRECISION", nil, false},

		// Parameterized single-arg types.
		{"char", "CHAR", "CHAR", nil, false},
		{"char_n", "CHAR(10)", "CHAR", []int{10}, false},
		{"character_n", "CHARACTER(20)", "CHARACTER", []int{20}, false},
		{"varchar", "VARCHAR", "VARCHAR", nil, false},
		{"varchar_n", "VARCHAR(255)", "VARCHAR", []int{255}, false},
		{"float", "FLOAT", "FLOAT", nil, false},
		{"float_p", "FLOAT(53)", "FLOAT", []int{53}, false},

		// CHARACTER VARYING two-token form.
		{"character_varying", "CHARACTER VARYING", "CHARACTER VARYING", nil, false},
		{"character_varying_n", "CHARACTER VARYING(80)", "CHARACTER VARYING", []int{80}, false},

		// Parameterized two-arg types.
		{"decimal", "DECIMAL", "DECIMAL", nil, false},
		{"decimal_p", "DECIMAL(10)", "DECIMAL", []int{10}, false},
		{"decimal_p_s", "DECIMAL(10,2)", "DECIMAL", []int{10, 2}, false},
		{"dec_p_s", "DEC(5,0)", "DEC", []int{5, 0}, false},
		{"numeric_p_s", "NUMERIC(18,4)", "NUMERIC", []int{18, 4}, false},

		// TIME with precision and WITH TIME ZONE.
		{"time", "TIME", "TIME", nil, false},
		{"time_p", "TIME(6)", "TIME", []int{6}, false},
		{"time_wtz", "TIME WITH TIME ZONE", "TIME", nil, true},
		{"time_p_wtz", "TIME(3) WITH TIME ZONE", "TIME", []int{3}, true},

		// Custom types (symbolPrimitive fallback).
		{"custom_ident", "MyType", "MyType", nil, false},
		{"custom_quoted", `"MyType"`, "MyType", nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			typeRef, err := p.parseType()
			if err != nil {
				t.Fatalf("parseType error: %v", err)
			}
			if typeRef.Name != tc.wantName {
				t.Errorf("Name = %q, want %q", typeRef.Name, tc.wantName)
			}
			if !intSliceEq(typeRef.Args, tc.wantArgs) {
				t.Errorf("Args = %v, want %v", typeRef.Args, tc.wantArgs)
			}
			if typeRef.WithTimeZone != tc.wantWithTZ {
				t.Errorf("WithTimeZone = %v, want %v", typeRef.WithTimeZone, tc.wantWithTZ)
			}
		})
	}
}

func intSliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestParser_Goldens iterates every .partiql file under
// testdata/parser-foundation/ and compares the parser's pretty-printed
// output (via ast.NodeToString) against the matching .golden file.
//
// Run with `go test -update -run TestParser_Goldens ./partiql/parser/...`
// to regenerate goldens after intentional AST shape changes.
func TestParser_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-foundation/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/parser-foundation/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got+"\n" != string(want) {
				t.Errorf("AST mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}

// TestParser_Stubs_Task9 locks in the deferred-feature error messages
// for SELECT and UNION — the two stubs added in Task 9. Later tasks
// will add more stub error cases and Task 12 consolidates them all
// into TestParser_Errors.
func TestParser_Stubs_Task9(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "select_stub",
			input:     "SELECT * FROM t",
			wantErrIn: "SELECT is deferred to parser-select (DAG node 5)",
		},
		{
			name:      "union_stub",
			input:     "a UNION b",
			wantErrIn: "UNION is deferred to parser-select (DAG node 5)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_Stubs_Task10 locks in the deferred-feature error messages
// for function calls and graph MATCH — the two stubs upgraded in
// Task 10.
func TestParser_Stubs_Task10(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "funcall_stub",
			input:     "foo(x)",
			wantErrIn: `function call "foo" is deferred to parser-builtins (DAG node 15)`,
		},
		{
			name:      "graph_match_stub",
			input:     "(a MATCH (b))",
			wantErrIn: "graph MATCH expression is deferred to parser-graph (DAG node 16)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_AWSCorpus loads every .partiql file from
// testdata/aws-corpus/, filters out the 2 known-bad syntax-skeleton
// files, and asserts each one either (a) fully parses or (b) hits a
// deferred-feature stub error. Any other error (or a panic) indicates
// a parser bug.
//
// At foundation milestone (DAG node 4), most corpus files start with
// SELECT and hit the parser-select stub — that's expected. The
// summary log reports how many fully parsed vs stubbed.
func TestParser_AWSCorpus(t *testing.T) {
	skip := map[string]bool{
		"select-001.partiql": true, // syntax skeleton (bracket placeholders)
		"insert-002.partiql": true, // syntax skeleton (backtick placeholders)
	}
	files, err := filepath.Glob("testdata/aws-corpus/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no AWS corpus files found under testdata/aws-corpus/")
	}
	var fullyParsed, stubbed, skipped int
	for _, f := range files {
		name := filepath.Base(f)
		if skip[name] {
			skipped++
			continue
		}
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(data))
			_, err = p.ParseExpr()
			if err == nil {
				fullyParsed++
				return
			}
			if !strings.Contains(err.Error(), "deferred to") {
				t.Errorf("unexpected parse error (not a deferred-feature stub): %v", err)
				return
			}
			stubbed++
		})
	}
	t.Logf("AWS corpus: %d fully parsed, %d stubbed, %d skipped",
		fullyParsed, stubbed, skipped)
}
