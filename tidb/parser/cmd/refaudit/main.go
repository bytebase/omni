// Command refaudit is the TiDB syntax-mapping ledger: it checks that every
// "// Ref: ... parser.y:N" citation in omni's TiDB productions names a real rule
// in the vendored v8.5.5 parser.y, and reports which grammar rules no production
// cites. It is the lint-and-audit half of the rejection suite (design doc
// section 3.6): a ledger that keeps coverage visible, not a proof of rejection.
//
// It is a plain root-module command that reads text and
// writes text; it imports nothing from omni or pingcap. Run it from the repo
// root:
//
//	go run ./tidb/parser/cmd/refaudit            # rewrite the coverage report
//	go run ./tidb/parser/cmd/refaudit -check     # validate citations, write nothing
//
// A broken citation (a line number past the grammar section, or landing between
// rules) exits non-zero in either mode. The companion test in this package gates
// the same check and keeps the report current with -update.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// defaultGrammar is the vendored pkg/parser/parser.y that the compat rows are
// generated from; refaudit resolves citations against the same copy.
// defaultTree is the TiDB source tree whose productions carry the citations.
// defaultReport is the
// checked-in coverage report. All three are repo-root-relative, so a normal run
// is `go run ./tidb/parser/cmd/refaudit` from the root.
const (
	defaultGrammar = "tidb/testdata/parser.y"
	defaultTree    = "tidb"
	defaultReport  = "tidb/parser/cmd/refaudit/testdata/parser_y_coverage.md"
)

func main() {
	grammarPath := flag.String("grammar", defaultGrammar, "the vendored parser.y to resolve citations against")
	tree := flag.String("tree", defaultTree, "the Go source tree to scan for Ref citations")
	report := flag.String("report", defaultReport, "the coverage report to write")
	check := flag.Bool("check", false, "validate citations and print the summary, but write no report")
	flag.Parse()

	a, err := audit(*grammarPath, *tree)
	if err != nil {
		fail(err)
	}

	fmt.Printf("refaudit: %d rules, %d citations, %d cited, %d uncited, %d broken\n",
		len(a.Rules), a.NumCits, len(a.Rules)-len(a.uncitedRules()), len(a.uncitedRules()), len(a.Broken))

	if !*check {
		if err := os.WriteFile(*report, []byte(renderReport(a)), 0o644); err != nil {
			fail(fmt.Errorf("write %s: %w", *report, err))
		}
		fmt.Printf("wrote %s\n", *report)
	}

	if len(a.Broken) > 0 {
		for _, b := range a.Broken {
			fmt.Fprintf(os.Stderr, "broken citation %s:%d: %s\n", b.File, b.Line, b.Reason)
		}
		fail(fmt.Errorf("%d citation(s) do not resolve to a v8.5.5 rule", len(a.Broken)))
	}
}

// audit reads the grammar and the source tree and resolves every citation.
func audit(grammarPath, tree string) (Audit, error) {
	src, err := os.ReadFile(grammarPath)
	if err != nil {
		return Audit{}, fmt.Errorf("read %s: %w", grammarPath, err)
	}
	rules, _, _, err := parseGrammarRules(string(src))
	if err != nil {
		return Audit{}, fmt.Errorf("parse %s: %w", filepath.Base(grammarPath), err)
	}
	cits, named, err := scanTree(tree)
	if err != nil {
		return Audit{}, fmt.Errorf("scan %s: %w", tree, err)
	}
	return runAudit(rules, cits, named), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "refaudit:", err)
	os.Exit(1)
}
