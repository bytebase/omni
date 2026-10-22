package catalog

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// everyAnalyzedExpr holds one zero value of each AnalyzedExpr node.
var everyAnalyzedExpr = []AnalyzedExpr{
	&AggExpr{}, &AlternativeSubPlanExpr{}, &ArrayCoerceExprQ{}, &ArrayExprQ{},
	&BoolExprQ{}, &BooleanTestExpr{}, &CaseExprQ{}, &CaseTestExprQ{},
	&CoalesceExprQ{}, &CoerceToDomainExpr{}, &CoerceToDomainValueExpr{},
	&CoerceViaIOExpr{}, &CollateExprQ{}, &ConstExpr{}, &ConvertRowtypeExprQ{},
	&DistinctExprQ{}, &FieldSelectExprQ{}, &FieldStoreExprQ{}, &FuncCallExpr{},
	&GroupingFuncExpr{}, &InferenceElemExpr{}, &JsonConstructorExprQ{},
	&JsonExprQ{}, &JsonIsPredicateExpr{}, &JsonValueExprQ{},
	&MergeSupportFuncExpr{}, &MinMaxExprQ{}, &NamedArgExprQ{}, &NextValueExprQ{},
	&NullIfExprQ{}, &NullTestExpr{}, &OpExpr{}, &ParamExpr{},
	&PartitionBoundSpecExpr{}, &RelabelExpr{}, &RowCompareExprQ{}, &RowExprQ{},
	&SQLValueFuncExpr{}, &ScalarArrayOpExpr{}, &SetOpColumnExpr{},
	&SetToDefaultExpr{}, &SubLinkExpr{}, &SubPlanExpr{}, &SubscriptingRefExpr{},
	&VarExpr{}, &WindowFuncExpr{}, &XmlExprQ{},
}

// TestEveryAnalyzedExprListed keeps everyAnalyzedExpr in step with the
// package: every type with an exprType method must be in it.
func TestEveryAnalyzedExprListed(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	var declared []string
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				fd, ok := decl.(*ast.FuncDecl)
				if !ok || fd.Recv == nil || fd.Name.Name != "exprType" || len(fd.Recv.List) != 1 {
					continue
				}
				if star, ok := fd.Recv.List[0].Type.(*ast.StarExpr); ok {
					if id, ok := star.X.(*ast.Ident); ok {
						declared = append(declared, id.Name)
					}
				}
			}
		}
	}
	listed := make(map[string]bool)
	for _, e := range everyAnalyzedExpr {
		listed[reflect.TypeOf(e).Elem().Name()] = true
	}
	sort.Strings(declared)
	for _, name := range declared {
		if !listed[name] {
			t.Errorf("%s implements AnalyzedExpr but is not in everyAnalyzedExpr; add it there and to AnalyzedExprChildren", name)
		}
	}
	if len(declared) != len(listed) {
		t.Errorf("everyAnalyzedExpr lists %d nodes, the package declares %d", len(listed), len(declared))
	}
}

// TestAnalyzedExprChildrenCoversEveryField fills every expression field of
// every node, including those inside nested helper structs such as
// CaseWhenQ, and checks AnalyzedExprChildren returns each one in field order.
func TestAnalyzedExprChildrenCoversEveryField(t *testing.T) {
	exprIface := reflect.TypeOf((*AnalyzedExpr)(nil)).Elem()
	queryPtr := reflect.TypeOf((*Query)(nil))
	for _, proto := range everyAnalyzedExpr {
		typ := reflect.TypeOf(proto).Elem()
		t.Run(typ.Name(), func(t *testing.T) {
			var want []AnalyzedExpr
			next := func() AnalyzedExpr {
				e := &ConstExpr{Value: fmt.Sprintf("child%d", len(want))}
				want = append(want, e)
				return e
			}
			var fill func(v reflect.Value)
			fill = func(v reflect.Value) {
				for i := 0; i < v.NumField(); i++ {
					f := v.Field(i)
					ft := f.Type()
					switch {
					case ft == exprIface:
						f.Set(reflect.ValueOf(next()))
					case ft.Kind() == reflect.Slice && ft.Elem() == exprIface:
						s := reflect.MakeSlice(ft, 2, 2)
						s.Index(0).Set(reflect.ValueOf(next()))
						s.Index(1).Set(reflect.ValueOf(next()))
						f.Set(s)
					case ft == queryPtr:
						// Subqueries are not entered.
					case ft.Kind() == reflect.Ptr && ft.Elem().Kind() == reflect.Struct:
						p := reflect.New(ft.Elem())
						fill(p.Elem())
						f.Set(p)
					case ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr && ft.Elem().Elem().Kind() == reflect.Struct:
						s := reflect.MakeSlice(ft, 2, 2)
						for j := 0; j < 2; j++ {
							p := reflect.New(ft.Elem().Elem())
							fill(p.Elem())
							s.Index(j).Set(p)
						}
						f.Set(s)
					}
				}
			}
			node := reflect.New(typ)
			fill(node.Elem())
			got := AnalyzedExprChildren(node.Interface().(AnalyzedExpr))
			if len(got) != len(want) {
				t.Fatalf("got %d children, want %d", len(got), len(want))
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("child %d is %v, want %v", i, got[i], want[i])
				}
			}
		})
	}
}

func TestWalkAnalyzedExprPrunes(t *testing.T) {
	inner := &VarExpr{RangeIdx: 0, AttNum: 1}
	expr := &OpExpr{
		Left:  &FuncCallExpr{Args: []AnalyzedExpr{inner}},
		Right: &ConstExpr{Value: "1"},
	}
	var seen []string
	WalkAnalyzedExpr(expr, func(e AnalyzedExpr) bool {
		seen = append(seen, reflect.TypeOf(e).Elem().Name())
		_, isFunc := e.(*FuncCallExpr)
		return !isFunc
	})
	if got, want := strings.Join(seen, ","), "OpExpr,FuncCallExpr,ConstExpr"; got != want {
		t.Errorf("visited %s, want %s", got, want)
	}
}
