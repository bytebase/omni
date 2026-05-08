package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var allowedP2Statuses = map[string]struct{}{
	"p2_done":               {},
	"p2_partial":            {},
	"p2_stub":               {},
	"p2_deferred":           {},
	"p2_not_parser_visible": {},
	"p2_unknown":            {},
}

var p2PlaceholderValues = map[string]struct{}{
	"":        {},
	"todo":    {},
	"unknown": {},
	"n/a":     {},
	"none":    {},
}

func TestOracleP2BNFManifestCompleteness(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "p2_coverage.tsv"))
	historicalRows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "bnf_coverage.tsv"))
	historical := make(map[string]coverageRow, len(historicalRows))
	for _, row := range historicalRows {
		historical[row.Key] = row
	}

	manifest := make(map[string]coverageRow, len(rows))
	statusCounts := make(map[string]int)
	for _, row := range rows {
		if _, exists := manifest[row.Key]; exists {
			t.Fatalf("%s: duplicate P2 BNF manifest row", row.Key)
		}
		manifest[row.Key] = row

		if _, err := os.Stat(filepath.Join("bnf", row.Key)); err != nil {
			t.Fatalf("%s: P2 BNF manifest row has no matching BNF file: %v", row.Key, err)
		}

		status := row.Fields["p2_status"]
		if _, ok := allowedP2Statuses[status]; !ok {
			t.Fatalf("%s: unknown P2 status %q", row.Key, status)
		}
		if status == "p2_unknown" {
			t.Fatalf("%s: P2 status must be classified, got p2_unknown", row.Key)
		}
		if status == "catalog" {
			t.Fatalf("%s: historical catalog status must not be reused as P2 status", row.Key)
		}
		if historicalRow, ok := historical[row.Key]; ok && historicalRow.Fields["status"] == "catalog" && status == historicalRow.Fields["status"] {
			t.Fatalf("%s: historical catalog row was carried over unchanged", row.Key)
		}
		statusCounts[status]++

		requireP2ProofFields(t, row)
		switch status {
		case "p2_done":
			verifyP2DoneEvidence(t, row)
		case "p2_partial", "p2_stub", "p2_deferred":
			if isP2Placeholder(row.Fields["owner_phase"]) {
				t.Fatalf("%s: %s row must have owner_phase", row.Key, status)
			}
		case "p2_not_parser_visible":
			if !strings.Contains(row.Fields["current_gap"], "semantic") &&
				!strings.Contains(row.Fields["current_gap"], "catalog") &&
				!strings.Contains(row.Fields["current_gap"], "parser_visible_complete") {
				t.Fatalf("%s: p2_not_parser_visible row must name parser-visible AST evidence", row.Key)
			}
		}
	}

	files, err := filepath.Glob(filepath.Join("bnf", "*.bnf"))
	if err != nil {
		t.Fatalf("glob BNF files: %v", err)
	}
	sort.Strings(files)
	if len(rows) != len(files) {
		t.Fatalf("P2 BNF manifest rows=%d, want one row per BNF file=%d", len(rows), len(files))
	}
	for _, path := range files {
		name := filepath.Base(path)
		if _, ok := manifest[name]; !ok {
			t.Fatalf("BNF file %s is missing from P2 coverage manifest", name)
		}
	}

	t.Logf("Oracle P2 BNF coverage: rows=%d bnf_files=%d status=%v", len(rows), len(files), statusCounts)
}

func TestOracleP2SkipInventoryDrift(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "p2_skip_inventory.tsv"))
	p2Rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "p2_coverage.tsv"))

	p2ByBNF := make(map[string]coverageRow, len(p2Rows))
	for _, row := range p2Rows {
		p2ByBNF[row.Key] = row
	}

	inventory := make(map[string]coverageRow, len(rows))
	for _, row := range rows {
		if _, exists := inventory[row.Key]; exists {
			t.Fatalf("%s: duplicate P2 skip inventory row", row.Key)
		}
		inventory[row.Key] = row
		bnf := row.Fields["bnf_file"]
		if _, ok := p2ByBNF[bnf]; !ok {
			t.Fatalf("%s: skip inventory references BNF file absent from p2_coverage.tsv: %s", row.Key, bnf)
		}
		if p2ByBNF[bnf].Fields["p2_status"] == "p2_done" {
			t.Fatalf("%s: p2_done BNF %s still owns skip/stub inventory", row.Key, bnf)
		}
	}

	sites := scanP2SkipSites(t)
	for site := range sites {
		if _, ok := inventory[site]; !ok {
			t.Fatalf("%s: skip/stub site is missing from p2_skip_inventory.tsv", site)
		}
	}
	for site := range inventory {
		if _, ok := sites[site]; !ok {
			t.Fatalf("%s: stale skip/stub inventory row no longer matches scanner output", site)
		}
	}

	t.Logf("Oracle P2 skip inventory: rows=%d scanned_sites=%d", len(rows), len(sites))
}

func TestOracleP2CompletionHasNoParserDebt(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "p2_coverage.tsv"))
	var debt []string
	for _, row := range rows {
		switch row.Fields["p2_status"] {
		case "p2_done", "p2_not_parser_visible":
		default:
			debt = append(debt, row.Key+"="+row.Fields["p2_status"])
		}
	}
	if len(debt) > 0 {
		t.Fatalf("P2 completion requires zero parser debt rows, found %d: %s", len(debt), strings.Join(debt, ", "))
	}

	inventoryRows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "p2_skip_inventory.tsv"))
	if len(inventoryRows) > 0 {
		var sites []string
		for _, row := range inventoryRows {
			sites = append(sites, row.Key)
		}
		t.Fatalf("P2 completion requires zero skip/stub inventory rows, found %d: %s", len(inventoryRows), strings.Join(sites, ", "))
	}
}

func requireP2ProofFields(t *testing.T, row coverageRow) {
	t.Helper()

	required := []string{
		"p2_surface",
		"ast_target",
		"parser_entrypoint",
		"current_gap",
		"next_action",
		"positive_test",
		"negative_test",
		"loc_test",
	}
	for _, field := range required {
		if isP2Placeholder(row.Fields[field]) {
			t.Fatalf("%s: required P2 field %s has placeholder value %q", row.Key, field, row.Fields[field])
		}
	}
}

func verifyP2DoneEvidence(t *testing.T, row coverageRow) {
	t.Helper()

	for _, field := range []string{"positive_test", "negative_test", "loc_test"} {
		verifyP2TestEvidence(t, row, field, row.Fields[field])
	}
	verifyP2SymbolEvidence(t, row, "parser_entrypoint", "oracle/parser")
	verifyP2SymbolEvidence(t, row, "ast_target", "oracle/ast")
}

func verifyP2TestEvidence(t *testing.T, row coverageRow, field string, evidence string) {
	t.Helper()

	path, name, ok := strings.Cut(evidence, "::")
	if !ok || path == "" || name == "" {
		t.Fatalf("%s: %s evidence must be path::test_name, got %q", row.Key, field, evidence)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if strings.HasPrefix(path, "oracle/parser/") {
			path = strings.TrimPrefix(path, "oracle/parser/")
			data, err = os.ReadFile(path)
		}
		if err != nil {
			t.Fatalf("%s: %s evidence file %s cannot be read: %v", row.Key, field, path, err)
		}
	}
	if !strings.Contains(string(data), name) {
		t.Fatalf("%s: %s evidence %s does not contain test name %s", row.Key, field, path, name)
	}
}

func verifyP2SymbolEvidence(t *testing.T, row coverageRow, field string, root string) {
	t.Helper()

	symbol := row.Fields[field]
	if strings.Contains(symbol, ".") {
		parts := strings.Split(symbol, ".")
		symbol = parts[len(parts)-1]
	}
	if symbol == "" {
		t.Fatalf("%s: %s evidence is empty", row.Key, field)
	}
	if _, err := os.Stat(root); err != nil && strings.HasPrefix(root, "oracle/parser") {
		root = strings.TrimPrefix(root, "oracle/parser")
		if root == "" {
			root = "."
		}
	}
	if _, err := os.Stat(root); err != nil && strings.HasPrefix(root, "oracle/ast") {
		root = filepath.Join("..", "ast")
	}
	var found bool
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), symbol) {
			found = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%s: scan %s for %s evidence: %v", row.Key, root, field, err)
	}
	if !found {
		t.Fatalf("%s: %s symbol %q was not found under %s", row.Key, field, row.Fields[field], root)
	}
}

func isP2Placeholder(value string) bool {
	_, ok := p2PlaceholderValues[strings.ToLower(strings.TrimSpace(value))]
	return ok
}

func scanP2SkipSites(t *testing.T) map[string]struct{} {
	t.Helper()

	files, err := filepath.Glob(filepath.Join("*.go"))
	if err != nil {
		t.Fatalf("glob parser Go files: %v", err)
	}
	sites := make(map[string]struct{})
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		scanP2SkipSitesInFile(sites, filepath.ToSlash(filepath.Join("oracle", "parser", file)), string(data))
	}
	return sites
}

func scanP2SkipSitesInFile(sites map[string]struct{}, file string, src string) {
	functionName := "file_scope"
	functionHasSkipToSemicolon := false
	functionHasNilNilReturn := false
	functionSites := make(map[string]struct{})

	flushFunction := func() {
		if functionHasSkipToSemicolon && functionHasNilNilReturn {
			functionSites[siteID(file, functionName, "nil_nil_after_skipToSemicolon")] = struct{}{}
		}
		for site := range functionSites {
			sites[site] = struct{}{}
		}
		functionSites = make(map[string]struct{})
		functionHasSkipToSemicolon = false
		functionHasNilNilReturn = false
	}

	for _, line := range strings.Split(src, "\n") {
		if name, ok := p2FunctionName(line); ok {
			flushFunction()
			functionName = name
		}
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		for _, pattern := range []string{
			"skipToSemicolon",
			"skipAlterTableClauseDetails",
			"skipParenthesizedBlock",
			"skipToNextClause",
			"parseAdminDDLStmt",
			"parseAlterGeneric",
		} {
			if strings.Contains(line, pattern) {
				functionSites[siteID(file, functionName, pattern)] = struct{}{}
			}
		}
		if strings.Contains(line, "skipToSemicolon") {
			functionHasSkipToSemicolon = true
		}
		if strings.Contains(line, "return nil, nil") {
			functionHasNilNilReturn = true
		}
		if strings.HasPrefix(trimmed, "//") {
			switch {
			case strings.Contains(lower, "placeholder"):
				functionSites[siteID(file, functionName, "comment_placeholder")] = struct{}{}
			case strings.Contains(lower, "skip remaining"):
				functionSites[siteID(file, functionName, "comment_skip_remaining")] = struct{}{}
			case strings.Contains(lower, "skip details"):
				functionSites[siteID(file, functionName, "comment_skip_details")] = struct{}{}
			case strings.Contains(lower, "skip unrecognized"):
				functionSites[siteID(file, functionName, "comment_skip_unrecognized")] = struct{}{}
			}
		}
	}
	flushFunction()
}

func p2FunctionName(line string) (string, bool) {
	re := regexp.MustCompile(`^func\s+(?:\([^)]*\)\s+)?([A-Za-z0-9_]+)\s*\(`)
	matches := re.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func siteID(file string, function string, pattern string) string {
	return file + ":" + function + ":" + pattern
}
