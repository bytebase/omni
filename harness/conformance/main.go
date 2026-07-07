package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		engine     = flag.String("engine", "tidb", "engine to sweep (tidb)")
		corpus     = flag.String("corpus", "corpus", "corpus root (from fetch_corpus.sh)")
		outDir     = flag.String("out", "out", "JSONL output dir (gitignored)")
		boardDir   = flag.String("scoreboards", "scoreboards", "committed scoreboard dir")
		omniSHA    = flag.String("omni-sha", "unknown", "omni commit under test (git rev-parse HEAD)")
		adjudicate = flag.Bool("adjudicate", false, "probe divergences against a live container (Task 7)")
	)
	flag.Parse()
	if *engine != "tidb" {
		log.Fatalf("engine %q not implemented in slice 1", *engine)
	}
	if *adjudicate {
		log.Fatal("adjudicate mode lands in Task 7")
	}

	entries, err := extractTiDBCorpus(*corpus)
	if err != nil {
		log.Fatal(err)
	}
	meta := RunMeta{Engine: "tidb", EngineVersion: "v8.5.5", OmniSHA: *omniSHA, CorpusTag: "v8.5.5", ClassifierVersion: classifierVersion}
	rows, stats := buildRows(entries)
	meta.DuplicatesDropped = stats.duplicatesDropped
	meta.DuplicateLabelConflicts = stats.duplicateLabelConflicts

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := writeJSONL(filepath.Join(*outDir, "tidb.jsonl"), meta, rows); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*boardDir, 0o755); err != nil {
		log.Fatal(err)
	}
	board := renderScoreboard(meta, rows)
	if err := os.WriteFile(filepath.Join(*boardDir, "tidb.md"), []byte(board), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Print(board)
}

// buildStats summarizes what buildRows dropped or flagged, for run meta.
type buildStats struct {
	duplicatesDropped       int
	duplicateLabelConflicts int
}

// buildRows converts extracted corpus entries into evaluated, classified rows.
//
// Skip entries (SkipReason != "") pass through as Class=SKIP without omni
// evaluation and without classify() — classify would overwrite Class — and are
// never deduped. Non-skip entries are deduped by stmt_hash: the first
// occurrence (extraction order is deterministic) keeps its provenance;
// subsequent duplicates are dropped and counted. A dropped duplicate whose
// upstream label conflicts with the kept row's flips the kept row to
// INDETERMINATE: the label is context-dependent, so it is not ground truth.
func buildRows(entries []CorpusEntry) ([]Row, buildStats) {
	kept := map[string]int{} // stmt hash -> index of kept row
	rows := make([]Row, 0, len(entries))
	var stats buildStats
	for _, e := range entries {
		r := Row{
			Engine: "tidb", Lane: "upstream",
			SourcePath: e.SourcePath, Line: e.Line, TestName: e.TestName,
			SQL: e.SQL, StmtHash: stmtHash(e.SQL),
			Expected: e.Expected, Family: classifyFamily(e.SQL), SkipReason: e.SkipReason,
		}
		if e.SkipReason != "" {
			r.Class = ClassSkip
			rows = append(rows, r)
			continue
		}
		if i, dup := kept[r.StmtHash]; dup {
			stats.duplicatesDropped++
			if rows[i].Expected != e.Expected {
				stats.duplicateLabelConflicts++
				rows[i].Class = ClassIndeterminate
				rows[i].ClassifierReason = "duplicate_label_conflict"
				rows[i].DivergenceKey = "" // INDETERMINATE rows are not clustered
			}
			continue
		}
		kept[r.StmtHash] = len(rows)
		r.OmniVerdict, r.OmniError = omniTiDBVerdict(e.SQL)
		classify(&r)
		rows = append(rows, r)
	}
	return rows, stats
}
