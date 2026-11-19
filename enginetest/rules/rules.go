// Package rules holds a compat.tsv's grammar rules to the parser code that
// reads them. An engine's testdata/rules.tsv names, for each rule its rows
// exercise, the parser file and the function a reader should start at.
package rules

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/grammar"
)

// columnsHeader is the header a rules.tsv starts with.
const columnsHeader = "# columns: rule, file, function"

// A Suite is one engine's rules.tsv, the compat.tsv it describes and the
// repository root its file paths are relative to.
type Suite struct {
	Compat string // <engine>/parser/testdata/compat.tsv
	Rules  string // <engine>/testdata/rules.tsv
	Root   string // the repository root, so a row's file is repo-relative
}

// An Entry is one rules.tsv row: the grammar rule, the parser file that
// reads it, and the function in that file a reader should start at.
type Entry struct {
	Rule     string
	File     string
	Function string
	line     int
}

// A File is an engine's rules.tsv.
type File struct {
	Entries []Entry
}

// ReadFile reads and validates a rules.tsv.
func ReadFile(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, b)
}

// Parse reads and validates the text of a rules.tsv; name is where it came
// from, for errors.
func Parse(name string, b []byte) (*File, error) {
	f := &File{}
	seen := map[string]int{}
	header := false
	for i, raw := range strings.Split(string(b), "\n") {
		line := i + 1
		if raw == "" {
			continue
		}
		if strings.HasPrefix(raw, "#") {
			if len(f.Entries) > 0 {
				return nil, fmt.Errorf("%s:%d: a comment after the first row; comments belong in the header", name, line)
			}
			if raw == columnsHeader {
				header = true
			}
			continue
		}
		if !header {
			return nil, fmt.Errorf("%s:%d: the header %q is missing", name, line, columnsHeader)
		}
		cols := strings.Split(raw, "\t")
		if len(cols) != 3 {
			return nil, fmt.Errorf("%s:%d: want 3 tab-separated columns (rule, file, function), got %d", name, line, len(cols))
		}
		e := Entry{Rule: cols[0], File: cols[1], Function: cols[2], line: line}
		if err := e.validate(); err != nil {
			return nil, fmt.Errorf("%s:%d: %v", name, line, err)
		}
		if prev, dup := seen[e.Rule]; dup {
			return nil, fmt.Errorf("%s:%d: repeats line %d (%s)", name, line, prev, e.Rule)
		}
		seen[e.Rule] = line
		f.Entries = append(f.Entries, e)
	}
	if len(f.Entries) > 0 && !header {
		return nil, fmt.Errorf("%s: the header %q is missing", name, columnsHeader)
	}
	return f, nil
}

func (e Entry) validate() error {
	if e.Rule == "" {
		return fmt.Errorf("empty rule")
	}
	if strings.ContainsAny(e.Rule, " \t") {
		return fmt.Errorf("rule %q holds a space", e.Rule)
	}
	if e.File == "" {
		return fmt.Errorf("%s: empty file", e.Rule)
	}
	clean := filepath.Clean(filepath.FromSlash(e.File))
	if filepath.IsAbs(e.File) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%s: file %q is not repo-relative", e.Rule, e.File)
	}
	if e.Function == "" {
		return fmt.Errorf("%s: empty function", e.Rule)
	}
	if strings.ContainsAny(e.Function, " \t()") {
		return fmt.Errorf("%s: function %q is not a name", e.Rule, e.Function)
	}
	return nil
}

// Check fails t when s's rules.tsv and compat.tsv disagree, or when a row
// names a file or function that does not read its rule.
func Check(t testing.TB, s Suite) {
	t.Helper()
	problems, err := Problems(s)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, p := range problems {
		t.Error(p)
	}
}

// Problems returns what is wrong with s: a rule a compat row names with no
// rules.tsv row, a rules.tsv file that does not exist or a function it does
// not declare, and a rules.tsv row for a rule no compat row names.
func Problems(s Suite) ([]string, error) {
	cf, err := compat.ReadFile(s.Compat)
	if err != nil {
		return nil, err
	}
	rf, err := ReadFile(s.Rules)
	if err != nil {
		return nil, err
	}
	want := map[string]bool{}
	for _, r := range cf.Rows {
		if rule := grammar.RuleName(r.Rule); rule != "" {
			want[rule] = true
		}
	}
	have := map[string]Entry{}
	for _, e := range rf.Entries {
		have[e.Rule] = e
	}
	var problems []string
	for _, rule := range sortedKeys(want) {
		if _, ok := have[rule]; !ok {
			problems = append(problems, fmt.Sprintf("%s: rule %s has no rules.tsv row", s.Compat, rule))
		}
	}
	// Many rows share one parser file; read its declarations once.
	declared := map[string]map[string]bool{}
	cached := map[string]error{}
	for _, e := range rf.Entries {
		if !want[e.Rule] {
			problems = append(problems, fmt.Sprintf("%s:%d: rule %s is named by no compat.tsv row", s.Rules, e.line, e.Rule))
			continue
		}
		path := filepath.Join(s.Root, filepath.FromSlash(e.File))
		funcs, ok := declared[path]
		if !ok {
			funcs, cached[path] = declaredFuncs(path)
			declared[path] = funcs
		}
		if err := cached[path]; err != nil {
			problems = append(problems, fmt.Sprintf("%s:%d: rule %s: %v", s.Rules, e.line, e.Rule, err))
			continue
		}
		if !funcs[e.Function] {
			problems = append(problems, fmt.Sprintf("%s:%d: rule %s: %s does not declare func %s or func (p *Parser) %s",
				s.Rules, e.line, e.Rule, path, e.Function, e.Function))
		}
	}
	return problems, nil
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// declaredFuncs returns the functions path declares, either as functions or
// as methods on Parser.
func declaredFuncs(path string) (map[string]bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	funcs := map[string]bool{}
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fd.Recv == nil {
			funcs[fd.Name.Name] = true
			continue
		}
		if len(fd.Recv.List) == 1 && receiverName(fd.Recv.List[0].Type) == "Parser" {
			funcs[fd.Name.Name] = true
		}
	}
	return funcs, nil
}

func receiverName(e ast.Expr) string {
	if star, ok := e.(*ast.StarExpr); ok {
		e = star.X
	}
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}
