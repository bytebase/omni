package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/dynamodb/ast"
)

// TestParser_CreateIndexBracketKeyLiteral verifies the literal keys a
// bracket path step takes in a CREATE INDEX path list.
//
// Grammar (PartiQLParser.g4:141):
//
//	pathSimpleSteps : BRACKET_LEFT key=literal BRACKET_RIGHT  # PathSimpleLiteral
//
// The grammar takes any literal, and the reference parser's visitor takes a
// number or a string and refuses the rest (PartiQLParserDefault.kt:859-876,
// "Expected an integer or string literal"): see
// TestParser_CreateIndexBracketKeyReject.
func TestParser_CreateIndexBracketKeyLiteral(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "string_key",
			input: "CREATE INDEX ON t (a['k'])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:StringLit{Val:\"k\"}}]}]}",
		},
		{
			name:  "int_key",
			input: "CREATE INDEX ON t (a[42])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NumberLit{Val:42}}]}]}",
		},
		{
			name:  "decimal_key",
			input: "CREATE INDEX ON t (a[3.14])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NumberLit{Val:3.14}}]}]}",
		},
		{
			name:  "symbol_key",
			input: "CREATE INDEX ON t (a[b])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:VarRef{Name:b}}]}]}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tc.want {
				t.Errorf("AST mismatch\ngot:  %s\nwant: %s", got, tc.want)
			}
		})
	}
}

// TestParser_CreateIndexBracketKeyReject verifies that the bracket-key path
// step still rejects everything the grammar forbids: a non-literal/
// non-symbol bracket key is not a valid `pathSimpleSteps`. Oracle (the
// generated ANTLR parser) rejects every input below.
//
// Negative coverage is required: routing the literal set through
// parseLiteral must NOT loosen the gate into accepting general expressions
// (e.g. `a[1+1]`, `a[(null)]`), an empty bracket, or a literal that is no
// number or string.
func TestParser_CreateIndexBracketKeyReject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "empty_bracket",
			input:     "CREATE INDEX ON t (a[])",
			wantErrIn: "expected",
		},
		{
			name:      "missing_rbracket",
			input:     "CREATE INDEX ON t (a['k')",
			wantErrIn: "BRACKET_RIGHT",
		},
		{
			name:      "expr_not_literal",
			input:     "CREATE INDEX ON t (a[1+1])",
			wantErrIn: "BRACKET_RIGHT",
		},
		{
			name:      "paren_literal",
			input:     "CREATE INDEX ON t (a[(null)])",
			wantErrIn: "expected",
		},
		// Literals the grammar takes as a key and the reference parser's
		// visitor refuses (PartiQLParserDefault.kt:859-876).
		{name: "null_key", input: "CREATE INDEX ON t (a[null])", wantErrIn: "integer or string literal"},
		{name: "true_key", input: "CREATE INDEX ON t (a[true])", wantErrIn: "integer or string literal"},
		{name: "false_key", input: "CREATE INDEX ON t (a[false])", wantErrIn: "integer or string literal"},
		{name: "missing_key", input: "CREATE INDEX ON t (a[missing])", wantErrIn: "integer or string literal"},
		{name: "ion_key", input: "CREATE INDEX ON t (a[`42`])", wantErrIn: "integer or string literal"},
		{name: "dot_then_null_key", input: "CREATE INDEX ON t (a.b[null])", wantErrIn: "integer or string literal"},
		{name: "date_key", input: "CREATE INDEX ON t (a[DATE '2020-01-01'])", wantErrIn: "integer or string literal"},
		{name: "time_key", input: "CREATE INDEX ON t (a[TIME '12:00:00'])", wantErrIn: "integer or string literal"},
		{name: "timestamp_key", input: "CREATE INDEX ON t (a[TIMESTAMP '2020-01-01 00:00:00'])", wantErrIn: "integer or string literal"},
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

// TestParser_CreateTableDefinition covers CREATE TABLE's qualified name and
// column declarations, and DROP TABLE's qualified name
// (PartiQLParser.g4:75, 87-112), which the reference parser's visitor takes
// (visitCreateTable, visitColumnDeclaration, visitColumnConstraint).
func TestParser_CreateTableDefinition(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			"CREATE TABLE a.b",
			"CreateTableStmt{Name:VarRef{Name:b} Qualifier:[VarRef{Name:a}]}",
		},
		{
			`CREATE TABLE "s".a.b`,
			"CreateTableStmt{Name:VarRef{Name:b} Qualifier:[VarRef{Name:s CaseSensitive:true} VarRef{Name:a}]}",
		},
		{
			"CREATE TABLE t(a INT)",
			"CreateTableStmt{Name:VarRef{Name:t} Columns:[ColumnDef{Name:VarRef{Name:a} Type:TypeRef{Name:INT}}]}",
		},
		{
			"CREATE TABLE t(a NULL NULL, b DECIMAL(5, 2) NOT NULL, c zzz CONSTRAINT nn NOT NULL NULL)",
			"CreateTableStmt{Name:VarRef{Name:t} Columns:[" +
				"ColumnDef{Name:VarRef{Name:a} Type:TypeRef{Name:NULL} Constraints:[ColumnConstraint{NotNull:false}]} " +
				"ColumnDef{Name:VarRef{Name:b} Type:TypeRef{Name:DECIMAL Args:[5,2]} Constraints:[ColumnConstraint{NotNull:true}]} " +
				"ColumnDef{Name:VarRef{Name:c} Type:TypeRef{Name:zzz} Constraints:[ColumnConstraint{Name:VarRef{Name:nn} NotNull:true} ColumnConstraint{NotNull:false}]}]}",
		},
		{
			"DROP TABLE a.b",
			"DropTableStmt{Name:VarRef{Name:b} Qualifier:[VarRef{Name:a}]}",
		},
	}
	for _, tc := range cases {
		stmt, err := NewParser(tc.input).ParseStatement()
		if err != nil {
			t.Errorf("ParseStatement(%q): %v", tc.input, err)
			continue
		}
		if got := ast.NodeToString(stmt); got != tc.want {
			t.Errorf("ParseStatement(%q)\ngot:  %s\nwant: %s", tc.input, got, tc.want)
		}
		if loc := stmt.GetLoc(); loc != (ast.Loc{Start: 0, End: len(tc.input)}) {
			t.Errorf("ParseStatement(%q).Loc = %+v, want 0..%d", tc.input, loc, len(tc.input))
		}
	}
}
