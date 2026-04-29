package parser

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestOracleParserProgress(t *testing.T) {
	files := parserSourceFiles(t)
	fset := token.NewFileSet()

	var parseMethods int
	var parseMethodsWithError int
	var silentDiscards int

	for _, path := range files {
		file, err := goparser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !isParserParseMethod(fn) {
				continue
			}
			parseMethods++
			if funcReturnsError(fn) {
				parseMethodsWithError++
			}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			as, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			hasParseCall := false
			for _, rhs := range as.Rhs {
				call, ok := rhs.(*ast.CallExpr)
				if ok && isParserParseCall(call) {
					hasParseCall = true
					break
				}
			}
			if !hasParseCall {
				return true
			}
			for _, lhs := range as.Lhs {
				if isBlankIdent(lhs) {
					silentDiscards++
					break
				}
			}
			return true
		})
	}

	const softFailScenarios = 62
	strictScenarios := oracleCoverageManifestRows(t, "strictness_v2.tsv")
	keywordGoldenScenarios := oracleCoverageManifestRows(t, "oracle_keywords.tsv")
	bnfRows, bnfUnknown := oracleBNFCoverageProgress(t)
	locRows, locUnknown := oracleStatusCoverageProgress(t, "loc_node_coverage.tsv")
	locFixtureRows := oracleLocFixtureCoverageProgress(t)
	referenceRows := oracleCoverageManifestRows(t, "reference_oracle.tsv")
	targets := oracleScenarioTargetProgress(t)

	t.Logf("Oracle parser progress: parse_methods=%d with_error=%d bare=%d silent_discards=%d",
		parseMethods, parseMethodsWithError, parseMethods-parseMethodsWithError, silentDiscards)
	t.Logf("Oracle parser progress: soft_fail_scenarios=%d strict_scenarios=%d keyword_golden_scenarios=%d",
		softFailScenarios, strictScenarios, keywordGoldenScenarios)
	t.Logf("Oracle parser progress: bnf_rows_classified=%d bnf_rows_unknown=%d", bnfRows, bnfUnknown)
	t.Logf("Oracle parser progress: loc_node_rows=%d loc_node_fixture_rows=%d loc_node_rows_unknown=%d", locRows, locFixtureRows, locUnknown)
	t.Logf("Oracle parser progress: reference_oracle_rows=%d reference_oracle_mismatches=gated", referenceRows)
	t.Logf("Oracle parser progress: Loc corpus hard gate is enforced by TestVerifyCorpus")

	if parseMethods == 0 {
		t.Fatal("no parser methods found")
	}
	if parseMethodsWithError != parseMethods {
		t.Fatalf("parse methods with error = %d/%d, want all", parseMethodsWithError, parseMethods)
	}
	if silentDiscards != 0 {
		t.Fatalf("silent parse error discards = %d, want 0", silentDiscards)
	}
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
}

func oracleScenarioTargetProgress(t *testing.T) map[string]int {
	t.Helper()
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "scenario_targets.tsv"))
	targets := make(map[string]int, len(rows))
	for _, row := range rows {
		value, ok := parseNonNegativeInt(row.Fields["minimum"])
		if !ok {
			t.Fatalf("%s: invalid minimum %q", row.Fields["metric"], row.Fields["minimum"])
		}
		targets[row.Fields["metric"]] = value
	}
	return targets
}

func oracleBNFCoverageProgress(t *testing.T) (rows int, unknown int) {
	t.Helper()
	return oracleStatusCoverageProgress(t, "bnf_coverage.tsv")
}

func oracleStatusCoverageProgress(t *testing.T, file string) (rows int, unknown int) {
	t.Helper()
	coverageRows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", file))
	for _, row := range coverageRows {
		rows++
		if coverageStatus(row) == "unknown" {
			unknown++
		}
	}
	return rows, unknown
}

func oracleCoverageManifestRows(t *testing.T, file string) int {
	t.Helper()
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", file))
	return len(rows)
}

func oracleLocFixtureCoverageProgress(t *testing.T) int {
	t.Helper()
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "loc_node_coverage.tsv"))
	var covered int
	for _, row := range rows {
		if coverageStatus(row) == "covered" && row.Fields["fixture"] != "" {
			covered++
		}
	}
	return covered
}
