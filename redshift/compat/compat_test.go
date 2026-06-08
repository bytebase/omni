package compat

import (
	"os"
	"strings"
	"testing"
)

func TestAWSCommandManifestIsCompleteAndClassified(t *testing.T) {
	manifest, err := LoadAWSCommandManifest("aws_commands_manifest.json")
	if err != nil {
		t.Fatalf("LoadAWSCommandManifest returned error: %v", err)
	}

	if manifest.Source.URL == "" {
		t.Fatal("manifest source URL is empty")
	}
	if len(manifest.Commands) < 100 {
		t.Fatalf("expected at least 100 AWS Redshift commands, got %d", len(manifest.Commands))
	}

	seen := make(map[string]bool)
	for _, command := range manifest.Commands {
		if command.Name == "" {
			t.Fatal("command with empty name")
		}
		if seen[command.Name] {
			t.Fatalf("duplicate command %q", command.Name)
		}
		seen[command.Name] = true
		if command.Status == "" {
			t.Fatalf("%s: empty status", command.Name)
		}
		if !command.Status.Valid() {
			t.Fatalf("%s: invalid status %q", command.Name, command.Status)
		}
		if command.SampleSQL == "" && command.Status != StatusNotRelevant {
			t.Fatalf("%s: missing sample SQL", command.Name)
		}
	}

	for _, required := range []string{
		"CREATE TABLE",
		"COPY",
		"UNLOAD",
		"MERGE",
		"SHOW EXTERNAL TABLE",
		"VACUUM",
		"USE",
	} {
		if !seen[required] {
			t.Fatalf("required AWS command %q is missing from manifest", required)
		}
	}
}

func TestAWSCommandMatrixSamplesMatchStatus(t *testing.T) {
	manifest, err := LoadAWSCommandManifest("aws_commands_manifest.json")
	if err != nil {
		t.Fatalf("LoadAWSCommandManifest returned error: %v", err)
	}

	result := EvaluateAWSCommandManifest(manifest)
	if len(result.Commands) != len(manifest.Commands) {
		t.Fatalf("got %d command results, want %d", len(result.Commands), len(manifest.Commands))
	}

	for _, command := range result.Commands {
		switch command.Status {
		case StatusSupportedParse, StatusSupportedRuntime:
			if !command.ParseOK {
				t.Errorf("%s: status %s but sample did not parse: %v", command.Name, command.Status, command.ParseError)
			}
		case StatusExplicitUnsupported:
			if command.ParseOK {
				t.Errorf("%s: status explicit_unsupported but sample parsed", command.Name)
			}
			if command.ParseError == "" {
				t.Errorf("%s: unsupported sample has empty parse error", command.Name)
			}
		case StatusNotRelevant:
			// Not relevant commands do not need parse proof.
		default:
			t.Errorf("%s: unhandled status %q", command.Name, command.Status)
		}
	}
}

func TestLegacyCorpusStats(t *testing.T) {
	stats, err := EvaluateLegacyCorpus("../parser/testdata/legacy", "../parser/legacy_examples_test.go")
	if err != nil {
		t.Fatalf("EvaluateLegacyCorpus returned error: %v", err)
	}

	if stats.TotalFiles != 115 {
		t.Fatalf("TotalFiles = %d, want 115", stats.TotalFiles)
	}
	if stats.ExpectedFailureFiles != 6 {
		t.Fatalf("ExpectedFailureFiles = %d, want 6", stats.ExpectedFailureFiles)
	}
	if stats.NewFailureFiles != 0 {
		t.Fatalf("NewFailureFiles = %d, want 0: %v", stats.NewFailureFiles, stats.NewFailures)
	}
	if stats.PromotedFiles != 0 {
		t.Fatalf("PromotedFiles = %d, want 0: %v", stats.PromotedFiles, stats.Promoted)
	}
}

func TestLegacyStatementParity(t *testing.T) {
	stats, err := EvaluateLegacyStatementParity("../parser/testdata/legacy")
	if err != nil {
		t.Fatalf("EvaluateLegacyStatementParity returned error: %v", err)
	}

	if stats.TotalStatements < 100 {
		t.Fatalf("TotalStatements = %d, want at least 100", stats.TotalStatements)
	}
	if stats.OldAcceptsOmniRejects == 0 {
		t.Fatal("expected current Redshift fork to have at least one old-accepts/omni-rejects gap")
	}
	if stats.BothAccept == 0 {
		t.Fatal("expected at least one statement accepted by both legacy and omni parsers")
	}
	if len(stats.OldAcceptsOmniRejectExamples) == 0 {
		t.Fatal("expected old-accepts/omni-rejects examples for gap reporting")
	}
}

func TestRedshiftReferenceHarnessRequiresDSN(t *testing.T) {
	result := EvaluateReferenceRedshift(EvaluateReferenceOptions{})
	if result.Enabled {
		t.Fatal("reference Redshift harness should be disabled without a DSN")
	}
	if result.SkipReason == "" {
		t.Fatal("disabled reference Redshift harness should report a skip reason")
	}
}

func TestRuntimeSemanticStats(t *testing.T) {
	stats := EvaluateRuntimeSemantics()
	if stats.TotalChecks < 8 {
		t.Fatalf("TotalChecks = %d, want at least 8", stats.TotalChecks)
	}
	if stats.FailedChecks != 0 {
		t.Fatalf("FailedChecks = %d, want 0: %v", stats.FailedChecks, stats.Failures)
	}
	if stats.PassedChecks != stats.TotalChecks {
		t.Fatalf("PassedChecks = %d, want %d", stats.PassedChecks, stats.TotalChecks)
	}
}

func TestCompatibilityReportMarkdown(t *testing.T) {
	manifest, err := LoadAWSCommandManifest("aws_commands_manifest.json")
	if err != nil {
		t.Fatalf("LoadAWSCommandManifest returned error: %v", err)
	}
	commandResult := EvaluateAWSCommandManifest(manifest)

	legacyStats, err := EvaluateLegacyCorpus("../parser/testdata/legacy", "../parser/legacy_examples_test.go")
	if err != nil {
		t.Fatalf("EvaluateLegacyCorpus returned error: %v", err)
	}
	parityStats, err := EvaluateLegacyStatementParity("../parser/testdata/legacy")
	if err != nil {
		t.Fatalf("EvaluateLegacyStatementParity returned error: %v", err)
	}
	referenceResult := EvaluateReferenceRedshift(EvaluateReferenceOptions{})
	runtimeStats := EvaluateRuntimeSemantics()

	report := RenderFullMarkdownReport(commandResult, legacyStats, parityStats, referenceResult, runtimeStats)
	for _, required := range []string{
		"# Redshift Compatibility Report",
		"AWS Command Coverage",
		"Legacy Corpus",
		"Legacy Statement Parity",
		"Runtime Semantics",
		"Reference Redshift",
		"supported_parse",
		"explicit_unsupported",
	} {
		if !strings.Contains(report, required) {
			t.Fatalf("report missing %q:\n%s", required, report)
		}
	}
}

func TestCompatibilityReportFileIsCurrent(t *testing.T) {
	manifest, err := LoadAWSCommandManifest("aws_commands_manifest.json")
	if err != nil {
		t.Fatalf("LoadAWSCommandManifest returned error: %v", err)
	}
	commandResult := EvaluateAWSCommandManifest(manifest)

	legacyStats, err := EvaluateLegacyCorpus("../parser/testdata/legacy", "../parser/legacy_examples_test.go")
	if err != nil {
		t.Fatalf("EvaluateLegacyCorpus returned error: %v", err)
	}
	parityStats, err := EvaluateLegacyStatementParity("../parser/testdata/legacy")
	if err != nil {
		t.Fatalf("EvaluateLegacyStatementParity returned error: %v", err)
	}
	referenceResult := EvaluateReferenceRedshift(EvaluateReferenceOptions{})
	runtimeStats := EvaluateRuntimeSemantics()

	want := RenderFullMarkdownReport(commandResult, legacyStats, parityStats, referenceResult, runtimeStats)
	got, err := os.ReadFile("report.md")
	if err != nil {
		t.Fatalf("read report.md: %v", err)
	}
	if string(got) != want {
		t.Fatalf("report.md is stale; regenerate with go run ./redshift/compat/cmd/redshift-compat-report")
	}
}
