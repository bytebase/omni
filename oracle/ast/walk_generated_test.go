package ast

import (
	goast "go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"
)

func TestWalkGeneratedCoversNodeLikeFields(t *testing.T) {
	expected := expectedWalkFields(t)
	generated := generatedWalkFields(t)

	var missing []string
	for structName, fields := range expected {
		actualFields, ok := generated[structName]
		if !ok {
			missing = append(missing, structName+": missing case")
			continue
		}
		for _, field := range fields {
			if !actualFields[field] {
				missing = append(missing, structName+"."+field)
			}
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("walk_generated.go is missing walkable fields:\n%s", strings.Join(missing, "\n"))
	}
}

func expectedWalkFields(t *testing.T) map[string][]string {
	t.Helper()

	fset := token.NewFileSet()
	var files []*goast.File
	for _, src := range []string{"parsenodes.go", "node.go"} {
		f, err := parser.ParseFile(fset, src, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", src, err)
		}
		files = append(files, f)
	}

	structNames := map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*goast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts := spec.(*goast.TypeSpec)
				if _, ok := ts.Type.(*goast.StructType); ok {
					structNames[ts.Name.Name] = true
				}
			}
		}
	}

	out := map[string][]string{}
	for _, f := range files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*goast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts := spec.(*goast.TypeSpec)
				st, ok := ts.Type.(*goast.StructType)
				if !ok {
					continue
				}
				for _, fl := range st.Fields.List {
					if len(fl.Names) == 0 {
						continue
					}
					if !isWalkableFieldType(typeStringForTest(fl.Type), structNames) {
						continue
					}
					for _, name := range fl.Names {
						out[ts.Name.Name] = append(out[ts.Name.Name], name.Name)
					}
				}
			}
		}
	}
	return out
}

func generatedWalkFields(t *testing.T) map[string]map[string]bool {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "walk_generated.go", nil, 0)
	if err != nil {
		t.Fatalf("parse walk_generated.go: %v", err)
	}

	out := map[string]map[string]bool{}
	for _, decl := range f.Decls {
		fn, ok := decl.(*goast.FuncDecl)
		if !ok || fn.Name.Name != "walkChildren" {
			continue
		}
		goast.Inspect(fn.Body, func(n goast.Node) bool {
			sw, ok := n.(*goast.TypeSwitchStmt)
			if !ok {
				return true
			}
			for _, stmt := range sw.Body.List {
				clause := stmt.(*goast.CaseClause)
				structName := walkCaseStructName(clause)
				if structName == "" {
					continue
				}
				fields := map[string]bool{}
				for _, bodyStmt := range clause.Body {
					goast.Inspect(bodyStmt, func(n goast.Node) bool {
						sel, ok := n.(*goast.SelectorExpr)
						if !ok {
							return true
						}
						if ident, ok := sel.X.(*goast.Ident); ok && ident.Name == "n" {
							fields[sel.Sel.Name] = true
						}
						return true
					})
				}
				out[structName] = fields
			}
			return false
		})
	}
	return out
}

func walkCaseStructName(clause *goast.CaseClause) string {
	if len(clause.List) != 1 {
		return ""
	}
	star, ok := clause.List[0].(*goast.StarExpr)
	if !ok {
		return ""
	}
	ident, ok := star.X.(*goast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func typeStringForTest(expr goast.Expr) string {
	switch t := expr.(type) {
	case *goast.Ident:
		return t.Name
	case *goast.StarExpr:
		return "*" + typeStringForTest(t.X)
	case *goast.SelectorExpr:
		return typeStringForTest(t.X) + "." + t.Sel.Name
	case *goast.ArrayType:
		return "[]" + typeStringForTest(t.Elt)
	default:
		return ""
	}
}

func isWalkableFieldType(typStr string, structNames map[string]bool) bool {
	switch typStr {
	case "Node", "ExprNode", "TableExpr", "StmtNode", "*List":
		return true
	}

	excludedStructs := map[string]bool{
		"Loc":     true,
		"List":    true,
		"String":  true,
		"Integer": true,
		"Float":   true,
		"Boolean": true,
	}

	if strings.HasPrefix(typStr, "*") {
		name := typStr[1:]
		return structNames[name] && !excludedStructs[name]
	}
	if strings.HasPrefix(typStr, "[]") {
		elemType := typStr[2:]
		switch elemType {
		case "Node", "ExprNode", "TableExpr", "StmtNode":
			return true
		}
		if strings.HasPrefix(elemType, "*") {
			name := elemType[1:]
			return structNames[name] && !excludedStructs[name]
		}
	}
	return false
}
