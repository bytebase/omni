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

// TestParser_StmtGoldens iterates every .partiql file under
// testdata/parser-ddl/, testdata/parser-select/, and testdata/parser-dml/
// and compares ParseStatement output (via ast.NodeToString) against the
// matching .golden file.
//
// Run with `go test -update -run TestParser_StmtGoldens ./partiql/parser/...`
// to regenerate goldens after intentional AST shape changes.
func TestParser_StmtGoldens(t *testing.T) {
	dirs := []string{
		"testdata/parser-ddl",
		"testdata/parser-select",
		"testdata/parser-dml",
	}
	var allFiles []string
	for _, dir := range dirs {
		files, err := filepath.Glob(dir + "/*.partiql")
		if err != nil {
			t.Fatal(err)
		}
		allFiles = append(allFiles, files...)
	}
	if len(allFiles) == 0 {
		t.Fatal("no golden inputs found")
	}
	for _, inPath := range allFiles {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(stmt)
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

// TestParser_DDLErrors verifies that malformed DDL statements produce
// the expected parse errors.
func TestParser_DDLErrors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "create_alone",
			input:     "CREATE",
			wantErrIn: "expected TABLE or INDEX after CREATE",
		},
		{
			name:      "create_unknown_keyword",
			input:     "CREATE SLAB t",
			wantErrIn: "expected TABLE or INDEX after CREATE",
		},
		{
			name:      "drop_alone",
			input:     "DROP",
			wantErrIn: "expected TABLE or INDEX after DROP",
		},
		{
			name:      "drop_unknown_keyword",
			input:     "DROP SLAB t",
			wantErrIn: "expected TABLE or INDEX after DROP",
		},
		{
			name:      "drop_index_missing_on",
			input:     "DROP INDEX idx t",
			wantErrIn: "expected ON",
		},
		{
			name:      "create_index_missing_paren",
			input:     "CREATE INDEX ON t name",
			wantErrIn: "expected PAREN_LEFT",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_DMLErrors verifies that malformed DML statements produce
// the expected parse errors.
func TestParser_DMLErrors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "insert_missing_into",
			input:     "INSERT t VALUE 1",
			wantErrIn: "expected INTO",
		},
		{
			name:      "insert_missing_value_or_expr",
			input:     "INSERT INTO",
			wantErrIn: "expected identifier",
		},
		{
			name:      "delete_missing_from",
			input:     "DELETE WHERE x = 1",
			wantErrIn: "expected FROM",
		},
		{
			name:      "update_missing_set",
			input:     "UPDATE t WHERE x = 1",
			wantErrIn: "expected SET",
		},
		{
			name:      "replace_missing_into",
			input:     "REPLACE t {'id': 1}",
			wantErrIn: "expected INTO",
		},
		{
			name:      "upsert_missing_into",
			input:     "UPSERT t {'id': 1}",
			wantErrIn: "expected INTO",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_Errors is the consolidated error-case test. It covers:
//
//  1. Deferred-feature stubs — one case per stub owner node, locking
//     in the exact error message for the grep contract (future DAG
//     node implementers grep for "deferred to parser-<name>" to find
//     their work items).
//  2. Real syntax errors — malformed inputs the parser must reject.
//
// This test replaces the per-task TestParser_Stubs_Task9 and
// TestParser_Stubs_Task10 stubs from earlier in the plan.
func TestParser_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// --- Deferred-feature stubs (one per owner node) ---
		{
			name:      "insert_stub",
			input:     "INSERT INTO t VALUE 1",
			wantErrIn: "INSERT is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "update_stub",
			input:     "UPDATE t SET x = 1",
			wantErrIn: "UPDATE is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "delete_stub",
			input:     "DELETE FROM t",
			wantErrIn: "DELETE is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "values_stub",
			input:     "VALUES (1, 2)",
			wantErrIn: "VALUES is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "valuelist_stub",
			input:     "(1, 2, 3)",
			wantErrIn: "valueList is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "lag_stub",
			input:     "LAG(x)",
			wantErrIn: "LAG() window is deferred to parser-window (DAG node 13)",
		},
		{
			name:      "count_stub",
			input:     "COUNT(x)",
			wantErrIn: "COUNT() aggregate is deferred to parser-aggregates (DAG node 14)",
		},
		{
			name:      "cast_stub",
			input:     "CAST(x AS INT)",
			wantErrIn: "CAST is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "substring_stub",
			input:     "SUBSTRING(s, 1, 2)",
			wantErrIn: "SUBSTRING is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "coalesce_stub",
			input:     "COALESCE(a, b)",
			wantErrIn: "COALESCE is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "list_constructor_stub",
			input:     "LIST(1, 2, 3)",
			wantErrIn: "LIST() constructor is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "graph_match_stub",
			input:     "(a MATCH (b))",
			wantErrIn: "graph MATCH expression is deferred to parser-graph (DAG node 16)",
		},
		{
			name:      "date_literal_stub",
			input:     "DATE '2026-01-01'",
			wantErrIn: "DATE literal is deferred to parser-datetime-literals (DAG node 18)",
		},
		{
			name:      "time_literal_stub",
			input:     "TIME '12:00:00'",
			wantErrIn: "TIME literal is deferred to parser-datetime-literals (DAG node 18)",
		},

		// --- Real syntax errors ---
		{
			name:      "case_no_when",
			input:     "CASE 1 + 1 END",
			wantErrIn: "CASE expression requires at least one WHEN clause",
		},
		{
			name:      "case_unclosed",
			input:     "CASE WHEN x > 0 THEN 1",
			wantErrIn: "expected END",
		},
		{
			name:      "case_missing_then",
			input:     "CASE WHEN x > 0",
			wantErrIn: "expected THEN",
		},
		{
			name:      "unclosed_paren",
			input:     "(1 + 2",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "unclosed_array",
			input:     "[1, 2",
			wantErrIn: "expected BRACKET_RIGHT",
		},
		{
			name:      "unclosed_bag",
			input:     "<<1, 2",
			wantErrIn: "expected ANGLE_DOUBLE_RIGHT",
		},
		{
			name:      "unclosed_tuple",
			input:     "{'a': 1",
			wantErrIn: "expected BRACE_RIGHT",
		},
		{
			name:      "tuple_missing_colon",
			input:     "{'a' 1}",
			wantErrIn: "expected COLON",
		},
		{
			name:      "between_missing_and",
			input:     "a BETWEEN 1 10",
			wantErrIn: "expected AND",
		},
		{
			name:      "is_invalid_type",
			input:     "a IS INT",
			wantErrIn: "IS predicate requires NULL, MISSING, TRUE, or FALSE",
		},
		{
			name:      "bare_comma",
			input:     ",",
			wantErrIn: "unexpected token",
		},
		{
			name:      "funcall_at_prefix",
			input:     "@foo(x)",
			wantErrIn: "@-prefix is not allowed before a function call",
		},
		{
			name:      "funcall_trailing_comma",
			input:     "foo(x,)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "funcall_unclosed",
			input:     "foo(x",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "funcall_leading_comma",
			input:     "foo(,)",
			wantErrIn: "unexpected token",
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

// TestParser_SelectErrors verifies that malformed SELECT statements
// produce the expected parse errors.
func TestParser_SelectErrors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "select_missing_from",
			input:     "SELECT *",
			wantErrIn: "expected FROM",
		},
		{
			name:      "join_missing_on",
			input:     "SELECT * FROM t1 INNER JOIN t2",
			wantErrIn: "expected ON after JOIN",
		},
		{
			name:      "order_by_missing_by",
			input:     "SELECT * FROM t ORDER a",
			wantErrIn: "expected BY",
		},
		{
			name:      "nulls_missing_direction",
			input:     "SELECT * FROM t ORDER BY a NULLS MAYBE",
			wantErrIn: "expected FIRST or LAST after NULLS",
		},
		// LET clause rejects (letBinding requires `expr AS symbolPrimitive`;
		// AS is mandatory, the alias must be an identifier).
		{
			name:      "let_missing_binding",
			input:     "SELECT x FROM t LET",
			wantErrIn: "unexpected token",
		},
		{
			name:      "let_missing_as",
			input:     "SELECT x FROM t LET a x",
			wantErrIn: "expected AS",
		},
		{
			name:      "let_missing_alias",
			input:     "SELECT x FROM t LET a AS",
			wantErrIn: "expected identifier",
		},
		{
			name:      "let_alias_not_identifier",
			input:     "SELECT x FROM t LET a AS 1",
			wantErrIn: "expected identifier",
		},
		{
			name:      "let_trailing_comma",
			input:     "SELECT x FROM t LET a AS x,",
			wantErrIn: "unexpected token",
		},
		// PIVOT projection rejects (PIVOT expr AT expr; AT is mandatory, and
		// PIVOT is a SELECT-replacement so a FROM clause is still required).
		{
			name:      "pivot_missing_at",
			input:     "PIVOT v FROM t",
			wantErrIn: "expected AT",
		},
		{
			name:      "pivot_missing_value",
			input:     "PIVOT AT k FROM t",
			wantErrIn: "unexpected token",
		},
		{
			name:      "pivot_missing_at_expr",
			input:     "PIVOT v AT FROM t",
			wantErrIn: "unexpected token",
		},
		{
			name:      "pivot_missing_from",
			input:     "PIVOT v AT k",
			wantErrIn: "expected FROM",
		},
		// UNPIVOT rejects (UNPIVOT expr asIdent? atIdent? byIdent?; the source
		// expression is mandatory and a bare alias without AS is not allowed).
		{
			name:      "unpivot_missing_expr",
			input:     "SELECT * FROM UNPIVOT",
			wantErrIn: "unexpected token",
		},
		{
			name:      "unpivot_bare_alias",
			input:     "SELECT * FROM UNPIVOT t v",
			wantErrIn: "unexpected token",
		},
		{
			name:      "unpivot_missing_alias_after_as",
			input:     "SELECT * FROM UNPIVOT t AS",
			wantErrIn: "expected identifier",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_FuncCall verifies that generic IDENT(args) function calls
// produce *ast.FuncCall nodes per DAG node 15a (parser-builtins-generic-call).
// Compared against ast.NodeToString output for compact shape assertion.
func TestParser_FuncCall(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single_arg",
			input: "foo(x)",
			want:  "FuncCall{Name:foo Args:[VarRef{Name:x}]}",
		},
		{
			name:  "zero_args",
			input: "foo()",
			want:  "FuncCall{Name:foo Args:[]}",
		},
		{
			name:  "multi_arg_literals",
			input: "foo(1, 2, 3)",
			want:  "FuncCall{Name:foo Args:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}]}",
		},
		{
			name:  "nested_calls",
			input: "f(g(x))",
			want:  "FuncCall{Name:f Args:[FuncCall{Name:g Args:[VarRef{Name:x}]}]}",
		},
		{
			name:  "call_with_path_step",
			input: "f(x).bar",
			want:  "PathExpr{Root:FuncCall{Name:f Args:[VarRef{Name:x}]} Steps:[DotStep{Field:bar}]}",
		},
		{
			name:  "quoted_arg_name",
			input: `attribute_exists("a")`,
			want:  "FuncCall{Name:attribute_exists Args:[VarRef{Name:a CaseSensitive:true}]}",
		},
		{
			name:  "dynamodb_begins_with",
			input: `begins_with(addr, '7834')`,
			want:  `FuncCall{Name:begins_with Args:[VarRef{Name:addr} StringLit{Val:"7834"}]}`,
		},
		// DAG node 15b.1: reserved-name function calls (SIZE and EXISTS).
		{
			name:  "reserved_size",
			input: "SIZE(Awards)",
			want:  "FuncCall{Name:SIZE Args:[VarRef{Name:Awards}]}",
		},
		{
			name:  "reserved_size_path_arg",
			input: "SIZE(Music.Albums)",
			want:  "FuncCall{Name:SIZE Args:[PathExpr{Root:VarRef{Name:Music} Steps:[DotStep{Field:Albums}]}]}",
		},
		{
			name:  "reserved_size_lowercase",
			input: "size(Awards)",
			want:  "FuncCall{Name:SIZE Args:[VarRef{Name:Awards}]}",
		},
		{
			name:  "reserved_exists_subquery",
			input: "EXISTS(SELECT * FROM Music)",
			want:  "FuncCall{Name:EXISTS Args:[SubLink{Stmt:SelectStmt{Star:true From:VarRef{Name:Music}}}]}",
		},
		{
			name:  "reserved_exists_via_path",
			input: "EXISTS(SELECT Awards FROM Music WHERE Artist = 'X')",
			want:  `FuncCall{Name:EXISTS Args:[SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:Awards}}] From:VarRef{Name:Music} Where:BinaryExpr{Op:= Left:VarRef{Name:Artist} Right:StringLit{Val:"X"}}}}]}`,
		},
		{
			name:  "reserved_upper",
			input: "UPPER(Artist)",
			want:  "FuncCall{Name:UPPER Args:[VarRef{Name:Artist}]}",
		},
		{
			name:  "reserved_lower_lowercase_keyword",
			input: "lower(Artist)",
			want:  "FuncCall{Name:LOWER Args:[VarRef{Name:Artist}]}",
		},
		{
			name:  "reserved_char_length",
			input: "CHAR_LENGTH(SongTitle)",
			want:  "FuncCall{Name:CHAR_LENGTH Args:[VarRef{Name:SongTitle}]}",
		},
		{
			name:  "reserved_character_length",
			input: "CHARACTER_LENGTH(SongTitle)",
			want:  "FuncCall{Name:CHARACTER_LENGTH Args:[VarRef{Name:SongTitle}]}",
		},
		{
			name:  "reserved_octet_length",
			input: "OCTET_LENGTH(SongTitle)",
			want:  "FuncCall{Name:OCTET_LENGTH Args:[VarRef{Name:SongTitle}]}",
		},
		{
			name:  "reserved_bit_length",
			input: "BIT_LENGTH(SongTitle)",
			want:  "FuncCall{Name:BIT_LENGTH Args:[VarRef{Name:SongTitle}]}",
		},
		// DAG node 15b.3: CASE expressions (searched and simple forms).
		{
			name:  "case_searched_one_when",
			input: "CASE WHEN x > 0 THEN 1 END",
			want:  "CaseExpr{Whens:[CaseWhen{When:BinaryExpr{Op:> Left:VarRef{Name:x} Right:NumberLit{Val:0}} Then:NumberLit{Val:1}}]}",
		},
		{
			name:  "case_searched_with_else",
			input: "CASE WHEN x > 0 THEN 1 ELSE 0 END",
			want:  "CaseExpr{Whens:[CaseWhen{When:BinaryExpr{Op:> Left:VarRef{Name:x} Right:NumberLit{Val:0}} Then:NumberLit{Val:1}}] Else:NumberLit{Val:0}}",
		},
		{
			name:  "case_searched_multi_when",
			input: "CASE WHEN x > 0 THEN 1 WHEN x < 0 THEN -1 ELSE 0 END",
			want:  "CaseExpr{Whens:[CaseWhen{When:BinaryExpr{Op:> Left:VarRef{Name:x} Right:NumberLit{Val:0}} Then:NumberLit{Val:1}} CaseWhen{When:BinaryExpr{Op:< Left:VarRef{Name:x} Right:NumberLit{Val:0}} Then:UnaryExpr{Op:- Operand:NumberLit{Val:1}}}] Else:NumberLit{Val:0}}",
		},
		{
			name:  "case_simple",
			input: "CASE Artist WHEN 'Acme' THEN 1 ELSE 0 END",
			want:  `CaseExpr{Operand:VarRef{Name:Artist} Whens:[CaseWhen{When:StringLit{Val:"Acme"} Then:NumberLit{Val:1}}] Else:NumberLit{Val:0}}`,
		},
		{
			name:  "case_simple_multi_when",
			input: "CASE Artist WHEN 'X' THEN 1 WHEN 'Y' THEN 2 ELSE 0 END",
			want:  `CaseExpr{Operand:VarRef{Name:Artist} Whens:[CaseWhen{When:StringLit{Val:"X"} Then:NumberLit{Val:1}} CaseWhen{When:StringLit{Val:"Y"} Then:NumberLit{Val:2}}] Else:NumberLit{Val:0}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("AST mismatch\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}
