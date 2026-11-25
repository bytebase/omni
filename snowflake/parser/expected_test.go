package parser

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestExpectedKeywords(t *testing.T) {
	tests := []struct {
		prefix string
		want   []string
		absent []string
	}{
		{prefix: "", want: []string{"SELECT", "WITH", "CREATE", "EXECUTE", "LS", "UNDROP", "WHILE"}},
		{prefix: "SELECT ", want: []string{"DISTINCT", "TOP", "CASE", "NOT", "EXISTS"}, absent: []string{"FROM", "WHERE"}},
		{prefix: "SELECT a ", want: []string{"FROM", "AS", "QUALIFY", "UNION", "MINUS"}},
		{prefix: "SELECT a FROM t ", want: []string{"WHERE", "QUALIFY", "JOIN", "ASOF", "AT", "BEFORE", "CHANGES", "MATCH_RECOGNIZE", "PIVOT", "SAMPLE", "LIMIT"}},
		{prefix: "SELECT a FROM t, ", want: []string{"LATERAL", "TABLE"}, absent: []string{"WHERE"}},
		{prefix: "SELECT a FROM t ORDER ", want: []string{"BY"}, absent: []string{"WHERE"}},
		{prefix: "SELECT a FROM t JOIN u ", want: []string{"ON", "USING"}},
		{prefix: "SELECT a FROM t QUALIFY ROW_NUMBER() OVER (", want: []string{"PARTITION", "ORDER"}},
		{prefix: "CREATE ", want: []string{"TABLE", "OR", "TRANSIENT", "WAREHOUSE", "COMPUTE", "STAGE"}},
		{prefix: "CREATE OR ", want: []string{"REPLACE", "ALTER"}},
		{prefix: "ALTER ", want: []string{"TABLE", "WAREHOUSE", "COMPUTE"}, absent: []string{"TEMPORARY"}},
		{prefix: "DROP ", want: []string{"TABLE", "COMPUTE"}},
		{prefix: "INSERT ", want: []string{"INTO", "OVERWRITE", "ALL", "FIRST"}},
		{prefix: "SELECT SUM(a) OVER (ORDER BY a ROWS BETWEEN ", want: []string{"UNBOUNDED", "CURRENT"}},
		// A prefix that fails before its end says nothing about what follows.
		{prefix: "SELECT FROM FROM ", want: nil},
	}
	for _, tt := range tests {
		got := ExpectedKeywords(tt.prefix)
		for _, w := range tt.want {
			if !slices.Contains(got, w) {
				t.Errorf("ExpectedKeywords(%q) lacks %s; got %v", tt.prefix, w, got)
			}
		}
		for _, w := range tt.absent {
			if slices.Contains(got, w) {
				t.Errorf("ExpectedKeywords(%q) has %s; got %v", tt.prefix, w, got)
			}
		}
		if tt.want == nil && tt.absent == nil && got != nil {
			t.Errorf("ExpectedKeywords(%q) = %v, want nil", tt.prefix, got)
		}
	}
}

// TestTokenTestsAreRecorded keeps collection mode complete: parser code tests
// the current or next token through curIs, nextIs and the word helpers of
// expected.go, which record the test, never by comparing a token type.
func TestTokenTestsAreRecorded(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") || name == "expected.go" {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		fset := token.NewFileSet()
		file, err := goparser.ParseFile(fset, name, src, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			var bad ast.Expr
			switch n := n.(type) {
			case *ast.BinaryExpr:
				if (n.Op == token.EQL || n.Op == token.NEQ) && (parserToken(n.X) != "" || parserToken(n.Y) != "") {
					bad = n
				}
			case *ast.SwitchStmt:
				if n.Tag != nil && parserToken(n.Tag) != "" {
					bad = n.Tag
				}
			case *ast.CallExpr:
				if id, ok := n.Fun.(*ast.Ident); ok && (id.Name == "isWordToken" || id.Name == "tokIsWord") && len(n.Args) > 0 && parserTokenValue(n.Args[0]) {
					bad = n
				}
			}
			if bad != nil {
				t.Errorf("%s: test the token with curIs, nextIs or a word helper of expected.go", fset.Position(bad.Pos()))
			}
			return true
		})
	}
}

// parserToken reports whether e is p.cur.Type or p.peekNext().Type.
func parserToken(e ast.Expr) string {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Type" || !parserTokenValue(sel.X) {
		return ""
	}
	return sel.Sel.Name
}

// parserTokenValue reports whether e is p.cur or p.peekNext().
func parserTokenValue(e ast.Expr) bool {
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 0 {
		e = call.Fun
		sel, ok := e.(*ast.SelectorExpr)
		return ok && sel.Sel.Name == "peekNext" && isIdent(sel.X, "p")
	}
	sel, ok := e.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "cur" && isIdent(sel.X, "p")
}

func isIdent(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == name
}
