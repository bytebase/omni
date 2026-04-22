package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestDispatchStrictnessGate is a GATING lint (Phase G1-GATE of the
// pg-paren-dispatch starmap) that scans pg/parser/*.go for
// `default: return nil, nil` sites and fails the build if any lacks
// an intent annotation. This prevents regression of the KB-2 class
// (silent-accept via grammar-gap default) and its Class-A family.
//
// Expected annotations on the comment line(s) immediately preceding
// the `default:` line:
//
//	// optional-probe: <reason>          — legit nil=absence pattern
//	// exhaustive: gram.y:N-M            — all productions enumerated
//	// known-gap: PAREN-KB-N              — known bug, tracked for fix
//
// See pg/parser/PARSER_DISPATCH_AUDIT.md for the classification
// reference. Annotations are informational — they don't change
// runtime behavior but force intent declaration at new sites.
func TestDispatchStrictnessGate(t *testing.T) {
	sites := findDefaultReturnNilNilSites(t)
	unannotated := 0
	annotated := 0
	for _, s := range sites {
		if s.Annotation == "" {
			unannotated++
			t.Errorf("UNANNOTATED %s:%d in %s — annotate with // optional-probe, // exhaustive, or // known-gap",
				s.File, s.Line, s.Function)
		} else {
			annotated++
		}
	}
	t.Logf("dispatch-strictness gate: %d sites total, %d annotated, %d unannotated",
		len(sites), annotated, unannotated)
	t.Logf("see pg/parser/PARSER_DISPATCH_AUDIT.md for classification")
}

type dispatchStrictnessSite struct {
	File       string
	Line       int
	Function   string
	Annotation string // "optional-probe" / "exhaustive" / "known-gap" / "" if none
}

func findDefaultReturnNilNilSites(t *testing.T) []dispatchStrictnessSite {
	t.Helper()
	var sites []dispatchStrictnessSite
	root := "."
	defaultRE := regexp.MustCompile(`^\s*default:\s*$`)
	returnNilRE := regexp.MustCompile(`^\s*return nil, nil\b`)
	funcRE := regexp.MustCompile(`^func\s+\([^)]+\)\s+(\w+)\s*\(`)
	annotRE := regexp.MustCompile(`//\s*(optional-probe|exhaustive|known-gap)`)

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(root, e.Name())
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		var lines []string
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		f.Close()
		currentFunc := ""
		for i, line := range lines {
			if m := funcRE.FindStringSubmatch(line); m != nil {
				currentFunc = m[1]
			}
			if defaultRE.MatchString(line) && i+1 < len(lines) && returnNilRE.MatchString(lines[i+1]) {
				site := dispatchStrictnessSite{
					File:     e.Name(),
					Line:     i + 1, // 1-based
					Function: currentFunc,
				}
				// Look backward up to 3 lines for an annotation comment.
				for j := i - 1; j >= 0 && j >= i-3; j-- {
					if m := annotRE.FindStringSubmatch(lines[j]); m != nil {
						site.Annotation = m[1]
						break
					}
					// Stop scanning back if we hit a non-comment, non-blank line.
					s := strings.TrimSpace(lines[j])
					if s != "" && !strings.HasPrefix(s, "//") {
						break
					}
				}
				sites = append(sites, site)
			}
		}
	}
	if len(sites) == 0 {
		t.Logf("warning: lint scanner found 0 `default: return nil, nil` sites — likely regex mismatch")
	}
	return sites
}

// TestDispatchStrictnessCountInAuditRange is a sanity check — expect
// ~38 sites based on the current audit. If the count drops dramatically,
// that's a signal the grep-pattern drifted or the sites were mass-fixed;
// investigate before the lint is flipped to gating (G1-GATE).
func TestDispatchStrictnessCountInAuditRange(t *testing.T) {
	got := len(findDefaultReturnNilNilSites(t))
	const minExpected = 10
	const maxExpected = 50
	if got < minExpected || got > maxExpected {
		t.Fatalf("dispatch site count out of expected range [%d,%d]: got %d. "+
			"Update PARSER_DISPATCH_AUDIT.md and adjust bounds in this test.",
			minExpected, maxExpected, got)
	}
	t.Logf("dispatch site count: %d (within [%d,%d] expected range)", got, minExpected, maxExpected)
}
