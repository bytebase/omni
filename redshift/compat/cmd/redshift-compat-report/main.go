package main

import (
	"fmt"
	"os"

	"github.com/bytebase/omni/redshift/compat"
)

func main() {
	manifestPath := "redshift/compat/aws_commands_manifest.json"
	legacyDir := "redshift/parser/testdata/legacy"
	legacyExpectedPath := "redshift/parser/legacy_examples_test.go"

	manifest, err := compat.LoadAWSCommandManifest(manifestPath)
	if err != nil {
		exitf("load AWS command manifest: %v", err)
	}
	commandResult := compat.EvaluateAWSCommandManifest(manifest)
	legacyStats, err := compat.EvaluateLegacyCorpus(legacyDir, legacyExpectedPath)
	if err != nil {
		exitf("evaluate legacy corpus: %v", err)
	}
	parityStats, err := compat.EvaluateLegacyStatementParity(legacyDir)
	if err != nil {
		exitf("evaluate legacy statement parity: %v", err)
	}
	referenceResult := compat.EvaluateReferenceRedshift(compat.EvaluateReferenceOptions{
		DSN: os.Getenv("REDSHIFT_COMPAT_DSN"),
	})
	runtimeStats := compat.EvaluateRuntimeSemantics()

	fmt.Print(compat.RenderFullMarkdownReport(commandResult, legacyStats, parityStats, referenceResult, runtimeStats))
}

func exitf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
