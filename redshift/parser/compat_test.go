package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/RedshiftDocs.g4, the Syntax sections of Amazon
// Redshift's SQL command reference converted by cmd/docsgrammar, and its
// near-misses, with omni's verdict and the docs'. Redshift has no server CI
// can start, so the docs judge: a statement the grammar derives is
// accepted, and a near-miss it does not is undocumented, which passes
// whatever omni says. The judge is a recognizer, so every run asks it
// again and needs no network.
//
// The docs are the weakest judge (§2.3 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md). A
// recording that asks a dedicated Redshift cluster, REDSHIFT_COMPAT_DSN, is
// the upgrade: it would add a pin row for the cluster and give this suite
// compat.SourceEngine and a Judge that reads compat.RawParseVerdict.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join("..", "testdata", "RedshiftDocs.g4"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceDocs,
	})
}

// omniVerdict accepts a row that parses as exactly one statement.
func omniVerdict(sql string) string {
	stmts, err := Parse(sql)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return compat.Reject
	}
	return compat.Accept
}
