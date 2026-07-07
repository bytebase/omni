package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

// RunMeta is the run-level metadata block (design doc §3/§4). One per run file.
type RunMeta struct {
	Engine            string `json:"engine"`
	EngineVersion     string `json:"engine_version"`
	OmniSHA           string `json:"omni_sha"`
	CorpusTag         string `json:"corpus_tag"`
	ContainerDigest   string `json:"container_digest,omitempty"`
	ClassifierVersion string `json:"classifier_version"`
	DuplicatesDropped int    `json:"duplicates_dropped"`
	// DuplicateLabelConflicts counts dropped duplicates whose upstream label
	// conflicted with the kept row's; the kept row flips to INDETERMINATE.
	DuplicateLabelConflicts int `json:"duplicate_label_conflicts"`
}

const classifierVersion = "1" // bump when classify()/clusterKey() semantics change

// writeJSONL writes one meta line followed by all rows sorted by stmt_hash,
// so run files diff cleanly across harvests. Sorts rows in place. Buffered;
// flush and close errors are propagated (a run file must never be silently
// truncated).
func writeJSONL(path string, meta RunMeta, rows []Row) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	if err := enc.Encode(map[string]any{"meta": meta}); err != nil {
		return err
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].StmtHash < rows[j].StmtHash })
	for i := range rows {
		if err := enc.Encode(&rows[i]); err != nil {
			return err
		}
	}
	return w.Flush()
}

type cluster struct {
	Family      string
	Key         string
	Count       int
	Exemplar    string // shortest SQL in the cluster
	ExemplarSrc string
}

// clusterRows groups upstream-lane rows of one class by family+divergence key,
// matching the headline counts filter. The exemplar choice is total — shortest
// SQL, then lexicographic, with Count==1 marking the first row (an empty-SQL
// sentinel would be order-dependent) — so rendering never depends on row order.
func clusterRows(rows []Row, class Class) []cluster {
	m := map[string]*cluster{}
	for _, r := range rows {
		if r.Lane != "upstream" || r.Class != class {
			continue
		}
		k := r.Family + "\x00" + r.DivergenceKey
		c, ok := m[k]
		if !ok {
			c = &cluster{Family: r.Family, Key: r.DivergenceKey}
			m[k] = c
		}
		c.Count++
		if c.Count == 1 || len(r.SQL) < len(c.Exemplar) ||
			(len(r.SQL) == len(c.Exemplar) && r.SQL < c.Exemplar) {
			c.Exemplar = r.SQL
			c.ExemplarSrc = fmt.Sprintf("%s:%d", r.SourcePath, r.Line)
		}
	}
	out := make([]cluster, 0, len(m))
	for _, c := range m {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].Family != out[j].Family {
			return out[i].Family < out[j].Family
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// renderScoreboard produces the committed markdown summary. Upstream lane only
// in the headline (design §2). Deterministic: stable ordering everywhere.
func renderScoreboard(meta RunMeta, rows []Row) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Conformance scoreboard: %s\n\n", meta.Engine)
	fmt.Fprintf(&b, "| meta | value |\n|---|---|\n")
	fmt.Fprintf(&b, "| engine_version | %s |\n", meta.EngineVersion)
	fmt.Fprintf(&b, "| omni_sha | %s |\n", meta.OmniSHA)
	fmt.Fprintf(&b, "| corpus_tag | %s |\n", meta.CorpusTag)
	fmt.Fprintf(&b, "| container_digest | %s |\n", orDash(meta.ContainerDigest))
	fmt.Fprintf(&b, "| classifier_version | %s |\n\n", meta.ClassifierVersion)

	counts := map[Class]int{}
	var total int
	for _, r := range rows {
		if r.Lane != "upstream" {
			continue
		}
		counts[r.Class]++
		total++
	}
	gapClusters := clusterRows(rows, ClassGap)
	overClusters := clusterRows(rows, ClassOver)

	fmt.Fprintf(&b, "## Counts (upstream lane)\n\n| class | statements |\n|---|---|\n")
	for _, c := range []Class{ClassAgreeAccept, ClassAgreeReject, ClassGap, ClassOver, ClassIndeterminate, ClassSkip} {
		fmt.Fprintf(&b, "| %s | %d |\n", c, counts[c])
	}
	fmt.Fprintf(&b, "| duplicates_dropped | %d |\n", meta.DuplicatesDropped)
	fmt.Fprintf(&b, "| duplicate_label_conflicts | %d |\n", meta.DuplicateLabelConflicts)
	fmt.Fprintf(&b, "| total | %d |\n\n", total)
	fmt.Fprintf(&b, "GAP clusters: %d\n\nOVER clusters: %d\n\nClusters are the work unit; statement counts are coverage context.\n\n", len(gapClusters), len(overClusters))

	writeClusterSection(&b, "GAP clusters (engine accepts, omni rejects) — the burn-down list", gapClusters)
	writeClusterSection(&b, "OVER clusters (engine rejects, omni accepts) — triage: structural vs leniency", overClusters)
	return b.String()
}

func writeClusterSection(b *strings.Builder, title string, cs []cluster) {
	fmt.Fprintf(b, "## %s\n\n", title)
	if len(cs) == 0 {
		fmt.Fprintf(b, "none\n\n")
		return
	}
	fmt.Fprintf(b, "| n | family | count | divergence | exemplar | source |\n|---|---|---|---|---|---|\n")
	for i, c := range cs {
		fmt.Fprintf(b, "| %d | %s | %d | %s | `%s` | %s |\n",
			i+1, c.Family, c.Count, mdEscape(truncate(c.Key, 80)), mdEscape(truncate(c.Exemplar, 100)), c.ExemplarSrc)
	}
	fmt.Fprintf(b, "\n")
}

// truncate cuts at a rune boundary (backing off from a mid-rune position) so
// display strings stay valid UTF-8; the JSONL keeps the raw text.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n] + "..."
}

// mdEscape keeps table cells intact: pipes escaped, CR/LF flattened (both are
// CommonMark line endings), and backticks mapped to ' because exemplars render
// inside backtick code spans. Lossy but display-only; the JSONL keeps raw SQL.
func mdEscape(s string) string {
	return strings.NewReplacer("|", "\\|", "\r", " ", "\n", " ", "`", "'").Replace(s)
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
