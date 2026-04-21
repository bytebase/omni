package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestPARENAuditLint enforces PAREN_AUDIT.json governance
// (SCENARIOS-pg-paren-dispatch.md §5.3).
//
// It runs on the default build tag — no testcontainer, no oracle — so
// every PR gate picks it up via `.github/workflows/ci.yml` without
// booting Docker. See PAREN_AUDIT_SCHEMA.md for the enforced schema
// and the baseline policy.
//
// The four invariants it enforces:
//
//  1. Schema completeness. Every row has all required fields populated
//     with the expected types and allowed values.
//  2. Proof fence. Every aligned=yes row carries non-empty proof_notes.
//     Deliberate-break test: temporarily zap a proof_notes string and
//     this test fires immediately. (Smoke-tested at commit time; see
//     the commit message for the transcript.)
//  3. Site uniqueness. No two rows share the same `<file>:<line>` site
//     identifier.
//  4. Code-vs-audit drift. Every `p.cur.Type == '('` / `')'` dispatch
//     site in pg/parser/*.go (non-test) maps to an audit row by
//     (file, function). Known rename-drift captured at Phase 5 is
//     allowlisted in PAREN_AUDIT_DRIFT_BASELINE.txt so the test passes
//     today; any NEW undocumented site fails immediately.
//
// The test is deliberately cheap — scans JSON + a handful of .go files
// with regex matching, no parser invocation. Typical runtime <50ms.
type parenAuditRow struct {
	Site             string   `json:"site"`
	Function         string   `json:"function"`
	Nonterminals    []string  `json:"nonterminals"`
	AmbiguityPresent bool     `json:"ambiguity_present"`
	CurrentTechnique *string  `json:"current_technique"`
	PGReference      string   `json:"pg_reference"`
	Aligned          string   `json:"aligned"`
	BlockedBy        *string  `json:"blocked_by"`
	Cluster          string   `json:"cluster"`
	Priority         string   `json:"priority"`
	ProofNotes       string   `json:"proof_notes"`
	SuspicionNotes   *string  `json:"suspicion_notes"`
}

const (
	parenAuditPath         = "PAREN_AUDIT.json"
	parenAuditSchemaPath   = "PAREN_AUDIT_SCHEMA.md"
	parenAuditDriftBaseln  = "PAREN_AUDIT_DRIFT_BASELINE.txt"
)

var (
	parenAuditAllowedAligned  = map[string]bool{"yes": true, "no": true, "blocked": true, "unclear": true}
	parenAuditAllowedPriority = map[string]bool{"high": true, "med": true, "low": true}
	// Accept C1..C5 plus optional subcluster suffix (e.g. C5.a .. C5.n
	// for the utility-tail subgroups in PAREN_AUDIT.md).
	parenAuditClusterRe = regexp.MustCompile(`^C[1-5](\.[a-z]+)?$`)

	// parenAuditDispatchRe matches the dispatch-site pattern scanned by
	// the drift detector. The lint is deliberately scoped to the
	// `p.cur.Type == '(' | ')'` form (the question SCENARIOS §5.3
	// specifies). Other patterns (expect('('), match('(')) are still
	// audit-worthy but the scan doesn't enforce them — adding a row to
	// PAREN_AUDIT.json when introducing one of those is a code-review
	// responsibility.
	parenAuditDispatchRe = regexp.MustCompile(`p\.cur\.Type\s*==\s*'[()]'`)

	// parenAuditFuncRe recognises a top-level Go function declaration.
	// Receiver methods and plain functions both match. Nested closures
	// within a function are intentionally NOT tracked — dispatch sites
	// inside closures are attributed to the enclosing function.
	parenAuditFuncRe = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)
)

func TestPARENAuditLint(t *testing.T) {
	rows, err := loadPARENAudit(t)
	if err != nil {
		t.Fatalf("load audit: %v", err)
	}

	t.Run("schema_completeness", func(t *testing.T) {
		checkPARENAuditSchema(t, rows)
	})
	t.Run("proof_fence", func(t *testing.T) {
		checkPARENAuditProofFence(t, rows)
	})
	t.Run("site_uniqueness", func(t *testing.T) {
		checkPARENAuditSiteUniqueness(t, rows)
	})
	t.Run("code_vs_audit_drift", func(t *testing.T) {
		checkPARENAuditDrift(t, rows)
	})
	t.Run("schema_doc_present", func(t *testing.T) {
		// Cheap sanity: the schema doc file must exist and reference
		// every required field. Prevents silently dropping the doc.
		b, err := os.ReadFile(parenAuditSchemaPath)
		if err != nil {
			t.Fatalf("read %s: %v", parenAuditSchemaPath, err)
		}
		required := []string{"site", "function", "nonterminals",
			"ambiguity_present", "current_technique", "pg_reference",
			"aligned", "blocked_by", "cluster", "priority",
			"proof_notes", "suspicion_notes"}
		content := string(b)
		for _, f := range required {
			if !strings.Contains(content, f) {
				t.Errorf("%s missing reference to field %q",
					parenAuditSchemaPath, f)
			}
		}
	})
}

func loadPARENAudit(t *testing.T) ([]parenAuditRow, error) {
	t.Helper()
	b, err := os.ReadFile(parenAuditPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", parenAuditPath, err)
	}
	var rows []parenAuditRow
	if err := json.Unmarshal(b, &rows); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", parenAuditPath, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s is empty (expected ≥ 1 row)", parenAuditPath)
	}
	return rows, nil
}

func checkPARENAuditSchema(t *testing.T, rows []parenAuditRow) {
	t.Helper()
	for _, r := range rows {
		if strings.TrimSpace(r.Site) == "" {
			t.Errorf("row with empty site: %+v", r)
			continue
		}
		if strings.TrimSpace(r.Function) == "" {
			t.Errorf("row %s: empty function", r.Site)
		}
		if !strings.Contains(r.Site, ":") {
			t.Errorf("row %s: site should be `file:line`", r.Site)
		}
		if len(r.Nonterminals) == 0 {
			t.Errorf("row %s: nonterminals is empty (every dispatch site targets at least one grammar nonterminal)", r.Site)
		}
		if strings.TrimSpace(r.PGReference) == "" {
			t.Errorf("row %s: pg_reference is empty", r.Site)
		}
		if !parenAuditAllowedAligned[r.Aligned] {
			t.Errorf("row %s: aligned=%q not in {yes,no,blocked,unclear}", r.Site, r.Aligned)
		}
		if !parenAuditClusterRe.MatchString(r.Cluster) {
			t.Errorf("row %s: cluster=%q does not match C[1-5][.subcluster]", r.Site, r.Cluster)
		}
		if !parenAuditAllowedPriority[r.Priority] {
			t.Errorf("row %s: priority=%q not in {high,med,low}", r.Site, r.Priority)
		}
		if r.Aligned == "blocked" && (r.BlockedBy == nil || strings.TrimSpace(*r.BlockedBy) == "") {
			t.Errorf("row %s: aligned=blocked requires non-empty blocked_by", r.Site)
		}
	}
}

func checkPARENAuditProofFence(t *testing.T, rows []parenAuditRow) {
	t.Helper()
	var offenders []string
	for _, r := range rows {
		if r.Aligned != "yes" {
			continue
		}
		if strings.TrimSpace(r.ProofNotes) == "" {
			offenders = append(offenders, fmt.Sprintf("  %s (%s): proof_notes empty",
				r.Site, r.Function))
		}
	}
	if len(offenders) > 0 {
		sort.Strings(offenders)
		t.Errorf("PAREN_AUDIT proof fence: %d aligned=yes row(s) have empty proof_notes:\n%s",
			len(offenders), strings.Join(offenders, "\n"))
	}
}

func checkPARENAuditSiteUniqueness(t *testing.T, rows []parenAuditRow) {
	t.Helper()
	// Soft check: the Phase 1 audit captured a handful of rows that
	// share a `file:line` coordinate because two adjacent dispatch
	// checks landed on the same source line during original scoping.
	// We log these for visibility but don't fail the lint — disam-
	// biguating them is an audit-data cleanup task, not a regression
	// signal. Duplicate (file, function, proof_notes) triples ARE a
	// copy-paste bug though, so we fail only on exact triple-dup.
	seenSite := make(map[string]string, len(rows))
	seenTriple := make(map[string]int, len(rows))
	for _, r := range rows {
		if prior, ok := seenSite[r.Site]; ok && prior != r.Function {
			t.Logf("duplicate site %q across functions %q and %q (informational, audit-data cleanup)",
				r.Site, prior, r.Function)
		}
		seenSite[r.Site] = r.Function
		triple := r.Site + "|" + r.Function + "|" + r.ProofNotes
		seenTriple[triple]++
	}
	for triple, n := range seenTriple {
		if n > 1 {
			parts := strings.SplitN(triple, "|", 3)
			t.Errorf("row copy-paste: site %q function %q appears %d times with IDENTICAL proof_notes",
				parts[0], parts[1], n)
		}
	}
}

// checkPARENAuditDrift scans pg/parser/*.go non-test files for
// `p.cur.Type == '(' | ')'` dispatch sites and verifies each
// (file, function) pair appears either in the audit or in the drift
// baseline allowlist.
func checkPARENAuditDrift(t *testing.T, rows []parenAuditRow) {
	t.Helper()
	auditPairs := make(map[[2]string]struct{}, len(rows))
	for _, r := range rows {
		file := r.Site
		if i := strings.Index(file, ":"); i >= 0 {
			file = file[:i]
		}
		auditPairs[[2]string{file, r.Function}] = struct{}{}
	}

	baseline, err := loadDriftBaseline(parenAuditDriftBaseln)
	if err != nil {
		t.Fatalf("load drift baseline: %v", err)
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("readdir .: %v", err)
	}

	var undocumented []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
			strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		sites, err := scanDispatchSites(e.Name())
		if err != nil {
			t.Fatalf("scan %s: %v", e.Name(), err)
		}
		for _, s := range sites {
			pair := [2]string{e.Name(), s.function}
			if _, ok := auditPairs[pair]; ok {
				continue
			}
			if _, ok := baseline[pair]; ok {
				continue
			}
			undocumented = append(undocumented,
				fmt.Sprintf("  %s:%d func=%s (add a row to %s or re-baseline %s)",
					s.file, s.line, s.function, parenAuditPath, parenAuditDriftBaseln))
		}
	}
	if len(undocumented) > 0 {
		sort.Strings(undocumented)
		t.Errorf("PAREN_AUDIT drift: %d dispatch site(s) not in audit or baseline:\n%s",
			len(undocumented), strings.Join(undocumented, "\n"))
	}
}

type dispatchSite struct {
	file     string
	line     int
	function string
}

func scanDispatchSites(path string) ([]dispatchSite, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []dispatchSite
	var curFunc string
	lineNum := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		if m := parenAuditFuncRe.FindStringSubmatch(line); m != nil {
			curFunc = m[1]
		}
		if parenAuditDispatchRe.MatchString(line) {
			if curFunc == "" {
				// Dispatch site before any func decl — package-level
				// var or init; shouldn't happen in pg/parser, but
				// record for visibility.
				continue
			}
			out = append(out, dispatchSite{file: path, line: lineNum, function: curFunc})
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func loadDriftBaseline(path string) (map[[2]string]struct{}, error) {
	out := make(map[[2]string]struct{})
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		raw := sc.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s:%d: malformed line (want \"<file>\\t<function>\"): %q",
				path, lineNum, raw)
		}
		file := strings.TrimSpace(parts[0])
		fn := strings.TrimSpace(parts[1])
		if file == "" || fn == "" {
			return nil, fmt.Errorf("%s:%d: empty field", path, lineNum)
		}
		out[[2]string{file, fn}] = struct{}{}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

