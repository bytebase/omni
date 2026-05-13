package ast

import (
	goast "go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestLocMethods(t *testing.T) {
	if !(Loc{Start: 1, End: 3}).IsValid() {
		t.Fatal("expected loc with non-negative start/end to be valid")
	}
	if (Loc{Start: 1, End: -1}).IsValid() {
		t.Fatal("expected loc with negative end to be invalid")
	}
	if got := (Loc{Start: 4, End: 8}).Merge(Loc{Start: 1, End: 5}); got != (Loc{Start: 1, End: 8}) {
		t.Fatalf("Merge = %+v, want {Start:1 End:8}", got)
	}
}

func TestNodeLocExpressions(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want Loc
	}{
		{"BinaryExpr", &BinaryExpr{Loc: Loc{Start: 1, End: 2}}, Loc{Start: 1, End: 2}},
		{"FuncCallExpr", &FuncCallExpr{Loc: Loc{Start: 2, End: 5}}, Loc{Start: 2, End: 5}},
		{"DatePart", &DatePart{Loc: Loc{Start: 2, End: 6}}, Loc{Start: 2, End: 6}},
		{"NextValueForExpr", &NextValueForExpr{Loc: Loc{Start: 2, End: 7}}, Loc{Start: 2, End: 7}},
		{"ParseExpr", &ParseExpr{Loc: Loc{Start: 2, End: 8}}, Loc{Start: 2, End: 8}},
		{"JsonKeyValueExpr", &JsonKeyValueExpr{Loc: Loc{Start: 2, End: 9}}, Loc{Start: 2, End: 9}},
		{"ColumnRef", &ColumnRef{Loc: Loc{Start: 3, End: 4}}, Loc{Start: 3, End: 4}},
		{"LikeExpr", &LikeExpr{Loc: Loc{Start: 4, End: 9}}, Loc{Start: 4, End: 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NodeLoc(tt.node); got != tt.want {
				t.Fatalf("NodeLoc = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestNodeLocCoversLocBearingNodes(t *testing.T) {
	locBearing := locBearingStructs(t)
	if len(locBearing) == 0 {
		t.Fatal("expected at least one Loc-bearing AST node")
	}
	if genericNodeLocFallbackWorks() {
		return
	}

	covered := nodeLocTypeSwitchCases(t)

	var missing []string
	for _, name := range locBearing {
		if !covered[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("NodeLoc missing Loc-bearing node types: %s", strings.Join(missing, ", "))
	}
}

type nodeLocFallbackStub struct {
	Loc Loc
}

func (*nodeLocFallbackStub) nodeTag() {}

func genericNodeLocFallbackWorks() bool {
	loc := Loc{Start: 11, End: 17}
	return NodeLoc(&nodeLocFallbackStub{Loc: loc}) == loc
}

func locBearingStructs(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*goast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*goast.TypeSpec)
				if !ok || ts.Name.Name == "Loc" {
					continue
				}
				st, ok := ts.Type.(*goast.StructType)
				if !ok || !structHasLocField(st) {
					continue
				}
				names = append(names, ts.Name.Name)
			}
		}
	}
	sort.Strings(names)
	return names
}

func structHasLocField(st *goast.StructType) bool {
	for _, field := range st.Fields.List {
		ident, ok := field.Type.(*goast.Ident)
		if !ok || ident.Name != "Loc" {
			continue
		}
		for _, name := range field.Names {
			if name.Name == "Loc" {
				return true
			}
		}
	}
	return false
}

func nodeLocTypeSwitchCases(t *testing.T) map[string]bool {
	t.Helper()

	f, err := parser.ParseFile(token.NewFileSet(), "loc.go", nil, 0)
	if err != nil {
		t.Fatalf("parse loc.go: %v", err)
	}
	covered := make(map[string]bool)
	goast.Inspect(f, func(n goast.Node) bool {
		cc, ok := n.(*goast.CaseClause)
		if !ok {
			return true
		}
		for _, expr := range cc.List {
			star, ok := expr.(*goast.StarExpr)
			if !ok {
				continue
			}
			ident, ok := star.X.(*goast.Ident)
			if ok {
				covered[ident.Name] = true
			}
		}
		return true
	})
	return covered
}
