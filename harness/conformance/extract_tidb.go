package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CorpusEntry is a raw extracted statement before omni evaluation.
type CorpusEntry struct {
	SQL        string
	Expected   Verdict
	SourcePath string
	Line       int
	TestName   string
	SkipReason string
}

// extractTiDBCorpus walks the sparse-cloned pkg/parser dir (top level only —
// sub-packages don't use the {src, ok, restore} table shape) and extracts
// every testCase composite literal.
func extractTiDBCorpus(corpusDir string) ([]CorpusEntry, error) {
	dir := filepath.Join(corpusDir, "tidb", "pkg", "parser")
	files, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no test files under %s — run ./fetch_corpus.sh first", dir)
	}
	var all []CorpusEntry
	for _, f := range files {
		entries, err := extractTiDBFile(f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		all = append(all, entries...)
	}
	return all, nil
}

func extractTiDBFile(path string) ([]CorpusEntry, error) {
	fset := token.NewFileSet()
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return nil, err
	}
	var out []CorpusEntry
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		ast.Inspect(fn, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			switch {
			case isTestCaseSlice(cl):
				for _, elt := range cl.Elts {
					var entry CorpusEntry
					if e, ok := elt.(*ast.CompositeLit); ok {
						entry = parseTestCaseLit(e, fset)
					} else {
						// Zero-loss invariant: every slice element yields an entry
						// with provenance, even ones we can't read statically.
						entry = CorpusEntry{
							Line:       fset.Position(elt.Pos()).Line,
							SkipReason: "non_composite_element",
						}
					}
					entry.SourcePath = relCorpusPath(path)
					entry.TestName = fn.Name.Name
					out = append(out, entry)
				}
				return false // don't descend into the slice again
			case isTestCaseLit(cl):
				// Bare testCase{...} literals — append-form rows like
				// append(testcases, testCase{...}). Elements inside a matched
				// slice never reach this arm: the slice arm returns false.
				entry := parseTestCaseLit(cl, fset)
				entry.SourcePath = relCorpusPath(path)
				entry.TestName = fn.Name.Name
				out = append(out, entry)
				return false
			}
			return true
		})
	}
	return out, nil
}

// isTestCaseSlice matches []testCase{...} composite literals. Purely
// syntactic: it matches the type identifier by name with no type resolution —
// correct while the corpus package declares a single testCase type; re-verify
// that assumption before reusing this for another engine's corpus.
func isTestCaseSlice(cl *ast.CompositeLit) bool {
	arr, ok := cl.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	id, ok := arr.Elt.(*ast.Ident)
	return ok && id.Name == "testCase"
}

// isTestCaseLit matches bare testCase{...} composite literals (append-form
// rows). Same name-only caveat as isTestCaseSlice. Untyped elements inside a
// []testCase{...} slice have a nil Type and can never match.
func isTestCaseLit(cl *ast.CompositeLit) bool {
	id, ok := cl.Type.(*ast.Ident)
	return ok && id.Name == "testCase"
}

// parseTestCaseLit reads one {src, ok, restore} literal — positional or keyed.
func parseTestCaseLit(e *ast.CompositeLit, fset *token.FileSet) CorpusEntry {
	entry := CorpusEntry{Line: fset.Position(e.Pos()).Line}
	var srcExpr, okExpr ast.Expr
	if len(e.Elts) > 0 {
		if _, keyed := e.Elts[0].(*ast.KeyValueExpr); keyed {
			for _, el := range e.Elts {
				kv, ok := el.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok {
					continue
				}
				switch key.Name {
				case "src":
					srcExpr = kv.Value
				case "ok":
					okExpr = kv.Value
				}
			}
		} else if len(e.Elts) >= 2 {
			srcExpr, okExpr = e.Elts[0], e.Elts[1]
		}
	}
	sql, ok1 := stringValue(srcExpr)
	okVal, ok2 := boolValue(okExpr)
	if !ok1 || !ok2 {
		entry.SkipReason = "non_literal"
		return entry
	}
	entry.SQL = sql
	if okVal {
		entry.Expected = VerdictAccept
	} else {
		entry.Expected = VerdictReject
	}
	return entry
}

// stringValue resolves string literals and simple "a" + "b" concatenations.
func stringValue(e ast.Expr) (string, bool) {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(v.Value)
		return s, err == nil
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}
		l, ok1 := stringValue(v.X)
		r, ok2 := stringValue(v.Y)
		return l + r, ok1 && ok2
	}
	return "", false
}

func boolValue(e ast.Expr) (bool, bool) {
	id, ok := e.(*ast.Ident)
	if !ok {
		return false, false
	}
	switch id.Name {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

// relCorpusPath trims the absolute prefix so provenance is machine-independent.
func relCorpusPath(path string) string {
	if i := strings.Index(path, "corpus/"); i >= 0 {
		return path[i:]
	}
	return path
}
