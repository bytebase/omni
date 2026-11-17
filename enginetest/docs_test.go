package enginetest

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// TestDocumentedRecordingCommands holds every `go test ... -run` command a
// document publishes to the tests its package holds. A recording is executed by
// a reader who has the endpoint and no other way to tell it did nothing: a
// pattern that matches no test exits 0, writes no tape, and reads as a pass in
// a file whose whole purpose is to be run later.
func TestDocumentedRecordingCommands(t *testing.T) {
	problems, err := checkDocumentedCommands(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, problem := range problems {
		t.Error(problem)
	}
}

// checkDocumentedCommands returns one problem for each documented `go test
// -run` command whose pattern names no test of the package it names.
func checkDocumentedCommands(root string) ([]string, error) {
	funcs := map[string]map[string]bool{}
	var problems []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel != "." && (strings.HasPrefix(d.Name(), ".") || datedRecord(rel)) {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(rel, ".md") || batchArtifact(d.Name()) {
			return nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, cmd := range runCommands(rel, src) {
			problems = append(problems, checkCommand(root, cmd, funcs)...)
		}
		return nil
	})
	return problems, err
}

// datedRecord reports whether dir holds a record of how an engine was built,
// which docs/AGENTS.md keeps as written: a command there says what ran then,
// not what a reader executes today.
func datedRecord(rel string) bool {
	for _, dir := range []string{"docs/plans", "docs/specs", "docs/scenarios"} {
		if rel == dir || strings.HasPrefix(rel, dir+"/") {
			return true
		}
	}
	return false
}

// batchArtifact reports whether name is a pipeline's leftover in an engine
// directory, which AGENTS.md says is not guidance for changing the engine.
func batchArtifact(name string) bool {
	return name == "SKILL.md" || strings.HasPrefix(name, "SCENARIOS-")
}

// checkCommand reports why cmd names no test, and nothing when its pattern
// matches one or the command names no package this tree can resolve.
func checkCommand(root string, cmd runCommand, funcs map[string]map[string]bool) []string {
	top := topPattern(cmd.Run)
	if templateName(top) {
		return nil
	}
	re, err := regexp.Compile(top)
	if err != nil {
		return []string{fmt.Sprintf("%s: -run %q is not a regular expression: %v", cmd, cmd.Run, err)}
	}
	var dirs, names []string
	for _, arg := range cmd.Packages {
		for _, dir := range packageDirs(root, arg) {
			tests, err := testFuncs(filepath.Join(root, dir), funcs)
			if err != nil {
				return []string{fmt.Sprintf("%s: read the tests of %s: %v", cmd, packageName(dir), err)}
			}
			dirs = append(dirs, packageName(dir))
			for name := range tests {
				names = append(names, name)
			}
		}
	}
	if len(dirs) == 0 {
		return nil
	}
	for _, name := range names {
		if re.MatchString(name) {
			return nil
		}
	}
	sort.Strings(names)
	problem := fmt.Sprintf("%s names no test in %s: a -run pattern that matches nothing exits 0, records nothing, and reads as a pass",
		cmd, strings.Join(dirs, " "))
	if close := suggestions(top, names); len(close) > 0 {
		problem += "; closest: " + strings.Join(close, ", ")
	}
	return []string{problem}
}

// A runCommand is one `go test` command a document publishes.
type runCommand struct {
	Doc      string   // the document's path from the repo root
	Line     int      // the document's line, 1-based
	Packages []string // package arguments, as written
	Run      string   // the -run pattern
}

func (c runCommand) String() string {
	packages := strings.Join(c.Packages, " ")
	if packages == "" {
		packages = "."
	}
	return fmt.Sprintf("%s:%d: go test %s -run %q", c.Doc, c.Line, packages, c.Run)
}

// runCommands extracts the commands one document publishes. A command is one
// line: the repo writes each in a fenced block or a code span.
func runCommands(doc string, src []byte) []runCommand {
	var cmds []runCommand
	for i, line := range strings.Split(string(src), "\n") {
		for _, segment := range goTestSegments(line) {
			cmd, ok := parseRunCommand(segment)
			if !ok {
				continue
			}
			cmd.Doc, cmd.Line = doc, i+1
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

// goTestSegments returns the stretch of line that starts at each `go test`, so
// a line holding two commands is two.
func goTestSegments(line string) []string {
	var segments []string
	for from := 0; from < len(line); {
		start := indexGoTest(line[from:])
		if start < 0 {
			break
		}
		start += from
		end := len(line)
		if next := indexGoTest(line[start+len("go test"):]); next >= 0 {
			end = start + len("go test") + next
		}
		segments = append(segments, line[start:end])
		from = end
	}
	return segments
}

// indexGoTest returns where `go test` starts in s, or -1. A word on either side
// keeps prose from reading as a command: a sentence that says "go tests" holds
// the bytes but names none.
func indexGoTest(s string) int {
	for from := 0; from < len(s); {
		i := strings.Index(s[from:], "go test")
		if i < 0 {
			return -1
		}
		i += from
		before, after := byte(' '), byte(' ')
		if i > 0 {
			before = s[i-1]
		}
		if i+len("go test") < len(s) {
			after = s[i+len("go test")]
		}
		if !isNameByte(before) && !isNameByte(after) {
			return i
		}
		from = i + 1
	}
	return -1
}

func isNameByte(b byte) bool {
	return b == '_' || b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

// parseRunCommand reads the package arguments and the -run pattern out of one
// command, and reports whether it names a pattern at all.
func parseRunCommand(cmd string) (runCommand, bool) {
	var c runCommand
	fields := commandFields(strings.TrimSuffix(strings.TrimSpace(cmd), "`"))
	for i := 2; i < len(fields); i++ { // fields[0] is go, fields[1] is test
		switch f := fields[i]; {
		case f == "-run" && i+1 < len(fields):
			c.Run = fields[i+1]
			i++
		case strings.HasPrefix(f, "-run="):
			c.Run = strings.TrimPrefix(f, "-run=")
		case strings.HasPrefix(f, "-"):
		case strings.HasPrefix(f, "."):
			c.Packages = append(c.Packages, f)
		}
	}
	return c, c.Run != ""
}

// commandFields splits a command into fields, dropping the quotes around a
// quoted one so that a pattern's parentheses and pipes survive.
func commandFields(s string) []string {
	var fields []string
	var field strings.Builder
	inField, inSingle, inDouble := false, false, false
	for _, r := range s {
		switch {
		case inSingle:
			if r == '\'' {
				inSingle = false
			} else {
				field.WriteRune(r)
			}
		case inDouble:
			if r == '"' {
				inDouble = false
			} else {
				field.WriteRune(r)
			}
		case r == '\'':
			inSingle, inField = true, true
		case r == '"':
			inDouble, inField = true, true
		case unicode.IsSpace(r):
			if inField {
				fields = append(fields, field.String())
				field.Reset()
				inField = false
			}
		default:
			field.WriteRune(r)
			inField = true
		}
	}
	if inField {
		fields = append(fields, field.String())
	}
	return fields
}

// topPattern is the part of a pattern that names the top-level test: a
// pattern's later elements name subtests, which a static read cannot hold.
func topPattern(run string) string {
	top, _, _ := strings.Cut(run, "/")
	return top
}

// templateName reports whether a pattern is the placeholder these documents use
// for "a test's name", as in `-run '^TestName$'`: a guide that records a test of
// the reader's choosing cannot name it.
func templateName(pattern string) bool {
	pattern = strings.TrimSuffix(strings.TrimPrefix(pattern, "^"), "$")
	for _, alternative := range strings.Split(pattern, "|") {
		if alternative != "TestName" {
			return false
		}
	}
	return true
}

// packageDirs resolves a package argument to the directories of the root module
// it names, relative to root: `./x/` is x, `./x/...` is x and every package
// under it. An argument that names no directory here, such as the
// `./<engine>/driver/` a guide writes, has none.
func packageDirs(root, arg string) []string {
	dir := strings.TrimSuffix(strings.TrimPrefix(arg, "./"), "/")
	subtree := dir == "..." || strings.HasSuffix(dir, "/...")
	if dir == "..." {
		dir = "."
	}
	if subtree {
		dir = strings.TrimSuffix(dir, "/...")
	}
	if dir == "" || dir == "." {
		dir = "."
	}
	if _, err := os.Stat(filepath.Join(root, dir)); err != nil {
		return nil
	}
	// A subtree names the package itself and every package under it, and the
	// package itself need not hold a test file: the command then matches
	// nothing, which is the failure to report.
	dirs := []string{dir}
	if !subtree {
		return dirs
	}
	filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}
		if p == filepath.Join(root, dir) || !holdsTestFile(p) {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err == nil {
			dirs = append(dirs, filepath.ToSlash(rel))
		}
		return nil
	})
	return dirs
}

// packageName renders a directory as a package argument, as a command writes it.
func packageName(dir string) string {
	if dir == "." {
		return "./"
	}
	return "./" + dir + "/"
}

func holdsTestFile(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), "_test.go") {
			return true
		}
	}
	return false
}

// testFuncs returns the tests the package in dir holds, which are the names
// -run matches: TestXxx functions taking one *testing.T, in the files go test
// compiles for this platform.
func testFuncs(dir string, cache map[string]map[string]bool) (map[string]bool, error) {
	if tests, ok := cache[dir]; ok {
		return tests, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	tests := map[string]bool{}
	fset := token.NewFileSet()
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		if match, err := build.Default.MatchFile(dir, name); err == nil && !match {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && isTestFunc(fn) {
				tests[fn.Name.Name] = true
			}
		}
	}
	cache[dir] = tests
	return tests, nil
}

// isTestFunc reports whether fn is a test go test runs: named TestXxx, where
// the character after Test is not a lower-case letter, with one *testing.T
// parameter.
func isTestFunc(fn *ast.FuncDecl) bool {
	if fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	if r, _ := utf8.DecodeRuneInString(fn.Name.Name[len("Test"):]); unicode.IsLower(r) {
		return false
	}
	params := fn.Type.Params.List
	if len(params) != 1 {
		return false
	}
	star, ok := params[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "testing" && sel.Sel.Name == "T"
}

// suggestions returns the test names most like a pattern's, so that a failure
// also names the command that would work: a renamed test leaves its old name in
// the documents that record it.
func suggestions(pattern string, names []string) []string {
	core := strings.TrimSuffix(strings.TrimPrefix(pattern, "^"), "$")
	if i := strings.IndexAny(core, `\|()[]*+?{}`); i >= 0 {
		core = core[:i]
	}
	if len(core) <= len("Test") {
		return nil
	}
	type candidate struct {
		name string
		dist int
	}
	var candidates []candidate
	for _, name := range names {
		dist := editDistance(core, name)
		if dist <= 4 || strings.Contains(name, core) || len(name) > len("Test")+2 && strings.Contains(core, name) {
			candidates = append(candidates, candidate{name, dist})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].dist != candidates[j].dist {
			return candidates[i].dist < candidates[j].dist
		}
		return candidates[i].name < candidates[j].name
	})
	var out []string
	for _, c := range candidates {
		if len(out) == 3 {
			break
		}
		out = append(out, c.name)
	}
	return out
}

// editDistance is the number of insertions, deletions and substitutions that
// turn a into b.
func editDistance(a, b string) int {
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = min(previous[j]+1, current[j-1]+1, previous[j-1]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

// repoRoot is the directory holding go.mod, which a command's package
// arguments are relative to.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		if filepath.Dir(dir) == dir {
			t.Fatalf("no go.mod above %s", dir)
		}
		dir = filepath.Dir(dir)
	}
}

func TestRunCommands(t *testing.T) {
	cases := []struct {
		name string
		doc  string
		want []runCommand
	}{
		{
			name: "a fenced recording command",
			doc:  "Record it:\n\n```bash\nOMNI_RECORD=1 go test ./mysql/parser/ -run '^TestCompatibility$'\n```\n",
			want: []runCommand{{Line: 4, Packages: []string{"./mysql/parser/"}, Run: "^TestCompatibility$"}},
		},
		{
			name: "a quoted alternation in a code span",
			doc:  "Record one: `go test ./x/ -run \"TestXxx|TestYyy\" -count=1 -v` first.\n",
			want: []runCommand{{Line: 1, Packages: []string{"./x/"}, Run: "TestXxx|TestYyy"}},
		},
		{
			name: "a pattern before the package",
			doc:  "go test -run '^TestMigration_AddColumn$' ./oceanbase/catalog/\n",
			want: []runCommand{{Line: 1, Packages: []string{"./oceanbase/catalog/"}, Run: "^TestMigration_AddColumn$"}},
		},
		{
			name: "the -run= form",
			doc:  "go test -run=^TestSkip$ ./x/...\n",
			want: []runCommand{{Line: 1, Packages: []string{"./x/..."}, Run: "^TestSkip$"}},
		},
		{
			name: "two commands on one line",
			doc:  "go test ./a/ -run '^TestA$' && go test ./b/ -run '^TestB$'\n",
			want: []runCommand{
				{Line: 1, Packages: []string{"./a/"}, Run: "^TestA$"},
				{Line: 1, Packages: []string{"./b/"}, Run: "^TestB$"},
			},
		},
		{
			name: "a command with no pattern",
			doc:  "go test ./x/...\n",
		},
		{
			name: "a mention of the command",
			doc:  "Two go tests commands are byte-identical in the tapes.\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := runCommands("doc.md", []byte(c.doc))
			if len(got) != len(c.want) {
				t.Fatalf("got %d commands, want %d: %v", len(got), len(c.want), got)
			}
			for i := range got {
				if got[i].Line != c.want[i].Line || got[i].Run != c.want[i].Run || strings.Join(got[i].Packages, " ") != strings.Join(c.want[i].Packages, " ") {
					t.Errorf("command %d = %+v, want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
}

// TestCheckDocumentedCommands holds the checker to the documents it must fail
// and the ones it must leave alone.
func TestCheckDocumentedCommands(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "x/x_test.go", `package x

import "testing"

func TestExists(t *testing.T) {}

func TestMain(m *testing.M) {}

func Testlower(t *testing.T) {}
`)
	cases := []struct {
		doc  string
		body string
		want string // a substring of the one problem; "" for none
	}{
		{doc: "x/existing-test.md", body: "go test ./x/ -run '^TestExists$'\n"},
		{doc: "x/missing-test.md", body: "go test ./x/ -run '^TestGone$'\n", want: "names no test in ./x/: a -run pattern that matches nothing exits 0"},
		{doc: "x/close-name.md", body: "go test ./x/ -run '^TestExist$'\n", want: "closest: TestExists"},
		{doc: "x/subtree.md", body: "go test ./x/... -run '^TestExists$'\n"},
		{doc: "x/subtree-missing.md", body: "go test ./x/... -run '^TestGone$'\n", want: "names no test in ./x/"},
		{doc: "root-subtree.md", body: "go test ./... -run '^TestExists$'\n"},
		{doc: "root-subtree-missing.md", body: "go test ./... -run '^TestGone$'\n", want: "names no test in ./ "},
		{doc: "x/placeholder.md", body: "go test ./x/ -run '^TestName$'\n"},
		{doc: "x/unresolved-package.md", body: "go test ./<engine>/driver/ -run '^TestGone$'\n"},
		{doc: "x/testmain.md", body: "go test ./x/ -run '^TestMain$'\n", want: "names no test"},
		{doc: "x/lowercase.md", body: "go test ./x/ -run '^Testlower$'\n", want: "names no test"},
		{doc: "docs/plans/dated-record.md", body: "go test ./x/ -run '^TestGone$'\n"},
		{doc: "x/SKILL.md", body: "go test ./x/ -run '^TestGone$'\n"},
	}
	for _, c := range cases {
		writeFile(t, root, c.doc, c.body)
	}
	problems, err := checkDocumentedCommands(root)
	if err != nil {
		t.Fatal(err)
	}
	byDoc := map[string][]string{}
	for _, problem := range problems {
		doc, _, _ := strings.Cut(problem, ":")
		byDoc[doc] = append(byDoc[doc], problem)
	}
	for _, c := range cases {
		t.Run(c.doc, func(t *testing.T) {
			got := byDoc[c.doc]
			if c.want == "" {
				if len(got) != 0 {
					t.Fatalf("got %d problems, want none: %v", len(got), got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("got %d problems, want 1: %v", len(got), got)
			}
			if !strings.Contains(got[0], c.want) {
				t.Errorf("problem = %q, want it to hold %q", got[0], c.want)
			}
		})
	}
}

func writeFile(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
