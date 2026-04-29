package parser

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

func TestOracleParserContract(t *testing.T) {
	files := parserSourceFiles(t)
	fset := token.NewFileSet()

	var bare []string
	var withError int
	var silentDiscards []string

	for _, path := range files {
		file, err := goparser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !isParserParseMethod(fn) {
				continue
			}
			name := fset.Position(fn.Pos()).String() + " " + fn.Name.Name
			if funcReturnsError(fn) {
				withError++
			} else {
				bare = append(bare, name)
			}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			as, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			hasParseCall := false
			for _, rhs := range as.Rhs {
				call, ok := rhs.(*ast.CallExpr)
				if ok && isParserParseCall(call) {
					hasParseCall = true
					break
				}
			}
			if !hasParseCall {
				return true
			}
			for _, lhs := range as.Lhs {
				if isBlankIdent(lhs) {
					silentDiscards = append(silentDiscards, fset.Position(as.Pos()).String())
					break
				}
			}
			return true
		})
	}

	if len(bare) > 0 {
		t.Fatalf("oracle parser has %d parse methods without error return (%d with error). First offenders:\n%s",
			len(bare), withError, strings.Join(firstN(bare, 25), "\n"))
	}
	if len(silentDiscards) > 0 {
		t.Fatalf("oracle parser silently discards parse errors at:\n%s", strings.Join(silentDiscards, "\n"))
	}
}

func parserSourceFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read oracle/parser: %v", err)
	}
	var files []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, name)
	}
	return files
}

func isParserParseMethod(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || len(fn.Recv.List) != 1 || !strings.HasPrefix(fn.Name.Name, "parse") {
		return false
	}
	switch recv := fn.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		ident, ok := recv.X.(*ast.Ident)
		return ok && ident.Name == "Parser"
	case *ast.Ident:
		return recv.Name == "Parser"
	default:
		return false
	}
}

func funcReturnsError(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, result := range fn.Type.Results.List {
		if ident, ok := result.Type.(*ast.Ident); ok && ident.Name == "error" {
			return true
		}
	}
	return false
}

func isParserParseCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && strings.HasPrefix(sel.Sel.Name, "parse")
}

func isBlankIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "_"
}

func firstN(items []string, n int) []string {
	if len(items) <= n {
		return items
	}
	result := append([]string{}, items[:n]...)
	result = append(result, "...")
	return result
}
