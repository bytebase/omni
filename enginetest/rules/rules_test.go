package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

const compatText = `# pin: default img:1 sha256:abc v1
# columns: omni, expected, source, rule, statement
accept	accept	engine	r/A	A
accept	accept	engine	r/B	B
accept	accept	engine	r/A!drop-1	C
accept	accept	engine	r(X)?/X	D
accept	accept	engine	suite:mine	E
accept	accept	engine	issue:BYT-1	F
`

const parserText = `package parser

type Parser struct{}

func parseA() {}
func (p *Parser) parseB() {}
func (p *Parser) parseR() {}
`

func suite(t *testing.T, rulesText string) Suite {
	t.Helper()
	dir := t.TempDir()
	write(t, filepath.Join(dir, "compat", "testdata", "compat.tsv"), compatText)
	write(t, filepath.Join(dir, "engine", "parser", "parser.go"), parserText)
	write(t, filepath.Join(dir, "engine", "testdata", "rules.tsv"), rulesText)
	return Suite{
		Compat: filepath.Join(dir, "compat", "testdata", "compat.tsv"),
		Rules:  filepath.Join(dir, "engine", "testdata", "rules.tsv"),
		Root:   dir,
	}
}

const goodRules = "# columns: rule, file, function\n" +
	"r\tengine/parser/parser.go\tparseR\n" +
	"r(X)?\tengine/parser/parser.go\tparseR\n"

func TestProblemsPasses(t *testing.T) {
	problems, err := Problems(suite(t, goodRules))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) > 0 {
		t.Errorf("want a pass, got %q", problems)
	}
}

func TestProblemsFailsARuleWithoutARow(t *testing.T) {
	rules := strings.ReplaceAll(goodRules, "r(X)?\tengine/parser/parser.go\tparseR\n", "")
	problems, err := Problems(suite(t, rules))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "rule r(X)? has no rules.tsv row") {
		t.Errorf("want the missing rule, got %q", problems)
	}
}

func TestProblemsFailsAFileThatDoesNotExist(t *testing.T) {
	rules := strings.Replace(goodRules, "engine/parser/parser.go\tparseR", "engine/parser/missing.go\tparseR", 1)
	problems, err := Problems(suite(t, rules))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "missing.go") {
		t.Errorf("want the missing file, got %q", problems)
	}
}

func TestProblemsFailsAFunctionTheFileDoesNotDeclare(t *testing.T) {
	rules := strings.Replace(goodRules, "parseR", "parseMissing", 1)
	problems, err := Problems(suite(t, rules))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "does not declare func parseMissing") {
		t.Errorf("want the missing function, got %q", problems)
	}
}

func TestProblemsFailsARowForNoRule(t *testing.T) {
	rules := goodRules + "r(OLD)?\tengine/parser/parser.go\tparseR\n"
	problems, err := Problems(suite(t, rules))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "rule r(OLD)? is named by no compat.tsv row") {
		t.Errorf("want the orphan row, got %q", problems)
	}
}

func TestReadFileRequiresTheHeader(t *testing.T) {
	_, err := Parse("rules.tsv", []byte("r\tengine/parser/parser.go\tparseR\n"))
	if err == nil || !strings.Contains(err.Error(), "header") {
		t.Errorf("want a missing-header error, got %v", err)
	}
}
