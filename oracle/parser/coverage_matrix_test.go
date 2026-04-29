package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type coverageRow struct {
	Key    string
	Fields map[string]string
}

var allowedCoverageStatuses = map[string]struct{}{
	"covered":     {},
	"partial":     {},
	"missing":     {},
	"unsupported": {},
	"catalog":     {},
	"deferred":    {},
	"unknown":     {},
}

func readCoverageTSV(t *testing.T, path string) ([]coverageRow, []string) {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	var header []string
	var rows []coverageRow
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if header == nil {
			header = parts
			if len(header) == 0 || header[0] == "" {
				t.Fatalf("%s: missing header", path)
			}
			continue
		}
		if len(parts) != len(header) {
			t.Fatalf("%s: row has %d fields, want %d: %q", path, len(parts), len(header), line)
		}
		fields := make(map[string]string, len(header))
		for i, name := range header {
			fields[name] = parts[i]
		}
		rows = append(rows, coverageRow{Key: parts[0], Fields: fields})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan %s: %v", path, err)
	}
	if header == nil {
		t.Fatalf("%s: empty TSV", path)
	}
	return rows, header
}

func coverageStatus(row coverageRow) string {
	return row.Fields["status"]
}

func TestOracleBNFCoverageManifestCompleteness(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "bnf_coverage.tsv"))

	manifest := make(map[string]coverageRow, len(rows))
	statusCounts := make(map[string]int)
	familyCounts := make(map[string]int)
	for _, row := range rows {
		status := coverageStatus(row)
		if _, ok := allowedCoverageStatuses[status]; !ok {
			t.Fatalf("%s: unknown status %q", row.Key, status)
		}
		if status == "unknown" {
			t.Fatalf("%s: coverage status must be classified, got unknown", row.Key)
		}
		if _, exists := manifest[row.Key]; exists {
			t.Fatalf("%s: duplicate BNF manifest row", row.Key)
		}
		manifest[row.Key] = row
		statusCounts[status]++
		familyCounts[row.Fields["family"]]++
	}

	files, err := filepath.Glob(filepath.Join("bnf", "*.bnf"))
	if err != nil {
		t.Fatalf("glob BNF files: %v", err)
	}
	sort.Strings(files)
	for _, path := range files {
		name := filepath.Base(path)
		if _, ok := manifest[name]; !ok {
			t.Fatalf("BNF file %s is missing from coverage manifest", name)
		}
	}
	for name := range manifest {
		if _, err := os.Stat(filepath.Join("bnf", name)); err != nil {
			t.Fatalf("BNF manifest row %s has no matching file: %v", name, err)
		}
	}

	t.Logf("Oracle BNF coverage: rows=%d bnf_files=%d status=%v family=%v", len(rows), len(files), statusCounts, familyCounts)
}

func TestOracleHighValueBNFGapsClosed(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "bnf_coverage.tsv"))
	highValue := map[string]struct{}{
		"select":   {},
		"insert":   {},
		"update":   {},
		"delete":   {},
		"merge":    {},
		"create":   {},
		"alter":    {},
		"drop":     {},
		"comment":  {},
		"grant":    {},
		"revoke":   {},
		"truncate": {},
		"set":      {},
	}
	closed := 0
	for _, row := range rows {
		if _, ok := highValue[row.Fields["family"]]; !ok {
			continue
		}
		switch coverageStatus(row) {
		case "unknown", "missing":
			t.Fatalf("%s: high-value BNF family %s has status %s", row.Key, row.Fields["family"], coverageStatus(row))
		}
		closed++
	}
	if closed == 0 {
		t.Fatal("no high-value BNF rows checked")
	}
	t.Logf("Oracle high-value BNF rows closed=%d", closed)
}

func TestOracleBNFImplementationDebtRequiresApproval(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "bnf_coverage.tsv"))
	statusCounts := make(map[string]int)
	for _, row := range rows {
		status := coverageStatus(row)
		statusCounts[status]++
		switch status {
		case "covered":
			continue
		case "partial", "deferred", "catalog", "unsupported":
		default:
			t.Fatalf("%s: status %q is not allowed in strict BNF debt gate", row.Key, status)
		}
		if row.Fields["debt_class"] == "" {
			t.Fatalf("%s: non-covered BNF row must have debt_class", row.Key)
		}
		if row.Fields["approval"] == "" {
			t.Fatalf("%s: non-covered BNF row must have approval", row.Key)
		}
		if row.Fields["next_action"] == "" {
			t.Fatalf("%s: non-covered BNF row must have next_action", row.Key)
		}
	}
	t.Logf("Oracle BNF implementation debt approval status=%v", statusCounts)
}

func TestOracleScenarioTargets(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "scenario_targets.tsv"))
	targets := make(map[string]int, len(rows))
	for _, row := range rows {
		name := row.Fields["metric"]
		value, ok := parseNonNegativeInt(row.Fields["minimum"])
		if !ok {
			t.Fatalf("%s: invalid minimum %q", name, row.Fields["minimum"])
		}
		targets[name] = value
	}

	required := []string{
		"soft_fail_scenarios",
		"strict_scenarios",
		"keyword_golden_scenarios",
		"bnf_rows_classified",
		"bnf_rows_unknown",
		"loc_node_fixture_rows",
		"loc_node_rows_unknown",
		"reference_oracle_rows",
	}
	for _, metric := range required {
		if _, ok := targets[metric]; !ok {
			t.Fatalf("scenario target %q is missing", metric)
		}
	}

	t.Logf("Oracle scenario targets: %v", targets)
}

func TestOracleReferenceManifest(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "reference_oracle.tsv"))
	if len(rows) == 0 {
		t.Fatal("reference oracle manifest has no rows")
	}
	familyCounts := make(map[string]int)
	for _, row := range rows {
		switch row.Fields["expect"] {
		case "accept", "reject", "catalog", "unsafe":
		default:
			t.Fatalf("%s: unknown reference expectation %q", row.Key, row.Fields["expect"])
		}
		if row.Fields["sql"] == "" {
			t.Fatalf("%s: reference row has empty SQL", row.Key)
		}
		familyCounts[row.Fields["family"]]++
	}
	t.Logf("Oracle reference manifest rows=%d family=%v", len(rows), familyCounts)
}

func TestOracleCoverage(t *testing.T) {
	const softFailScenarios = 62

	targets := oracleScenarioTargetProgress(t)
	strictScenarios := oracleCoverageManifestRows(t, "strictness_v2.tsv")
	keywordGoldenScenarios := oracleCoverageManifestRows(t, "oracle_keywords.tsv")
	bnfRows, bnfUnknown := oracleBNFCoverageProgress(t)
	locRows, locUnknown := oracleStatusCoverageProgress(t, "loc_node_coverage.tsv")
	locFixtureRows := oracleLocFixtureCoverageProgress(t)
	referenceRows := oracleCoverageManifestRows(t, "reference_oracle.tsv")

	if softFailScenarios < targets["soft_fail_scenarios"] {
		t.Fatalf("soft-fail scenarios = %d, want at least %d", softFailScenarios, targets["soft_fail_scenarios"])
	}
	if strictScenarios < targets["strict_scenarios"] {
		t.Fatalf("strict scenarios = %d, want at least %d", strictScenarios, targets["strict_scenarios"])
	}
	if keywordGoldenScenarios < targets["keyword_golden_scenarios"] {
		t.Fatalf("keyword golden scenarios = %d, want at least %d", keywordGoldenScenarios, targets["keyword_golden_scenarios"])
	}
	if bnfRows < targets["bnf_rows_classified"] {
		t.Fatalf("BNF rows classified = %d, want at least %d", bnfRows, targets["bnf_rows_classified"])
	}
	if bnfUnknown != targets["bnf_rows_unknown"] {
		t.Fatalf("BNF rows unknown = %d, want %d", bnfUnknown, targets["bnf_rows_unknown"])
	}
	if locUnknown != targets["loc_node_rows_unknown"] {
		t.Fatalf("Loc node rows unknown = %d, want %d", locUnknown, targets["loc_node_rows_unknown"])
	}
	if locFixtureRows < targets["loc_node_fixture_rows"] {
		t.Fatalf("Loc node fixture rows = %d, want at least %d", locFixtureRows, targets["loc_node_fixture_rows"])
	}
	if referenceRows < targets["reference_oracle_rows"] {
		t.Fatalf("reference oracle rows = %d, want at least %d", referenceRows, targets["reference_oracle_rows"])
	}

	t.Logf("Oracle parser coverage gates: soft_fail=%d strict=%d keyword=%d bnf_rows=%d bnf_unknown=%d loc_rows=%d loc_fixture_rows=%d loc_unknown=%d reference_rows=%d",
		softFailScenarios, strictScenarios, keywordGoldenScenarios, bnfRows, bnfUnknown, locRows, locFixtureRows, locUnknown, referenceRows)
}

func parseNonNegativeInt(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
	}
	return n, true
}
