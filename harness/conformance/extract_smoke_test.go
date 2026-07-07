package main

import (
	"os"
	"testing"
)

// TestSmokeRealCorpus guards extraction health against the fetched corpus.
// Short-gated: omni CI runs -short and must never require the corpus.
func TestSmokeRealCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("corpus smoke test skipped in short mode")
	}
	if _, err := os.Stat("corpus/tidb"); err != nil {
		t.Skip("corpus not fetched — run ./fetch_corpus.sh")
	}
	entries, err := extractTiDBCorpus("corpus")
	if err != nil {
		t.Fatal(err)
	}
	var accepts, rejects, skips int
	for _, e := range entries {
		switch {
		case e.SkipReason != "":
			skips++
		case e.Expected == VerdictAccept:
			accepts++
		default:
			rejects++
		}
	}
	t.Logf("extracted %d entries: %d accept / %d reject / %d skip", len(entries), accepts, rejects, skips)
	if len(entries) < 3000 {
		t.Fatalf("suspiciously few entries: %d", len(entries))
	}
	if skips > len(entries)/10 {
		t.Fatalf("skip rate over 10%%: %d/%d", skips, len(entries))
	}
}
