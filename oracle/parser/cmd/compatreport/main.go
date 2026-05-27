package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	oracleparser "github.com/bytebase/omni/oracle/parser"
)

const (
	reportStartMarker = "<!-- oracle-compat-report:start -->"
	reportEndMarker   = "<!-- oracle-compat-report:end -->"
)

type compatRow struct {
	Fields map[string]string
}

type manifestSummary struct {
	Total            int
	Phases           map[string]int
	Families         map[string]int
	Expectations     map[string]int
	ReferenceClasses map[string]int
}

type reportData struct {
	Manifest                    manifestSummary
	Local                       map[string]int
	Reference                   map[string]int
	ReferenceByFamily           map[string]int
	PrivilegedReference         map[string]int
	PrivilegedReferenceByFamily map[string]int
}

func main() {
	manifestPath := flag.String("manifest", "oracle/parser/testdata/coverage/compat_oracle.tsv", "compat_oracle.tsv path")
	referenceLogPath := flag.String("reference-log", "", "optional Oracle reference go test output path; use - for stdin")
	runID := flag.String("run-id", "LOCALRUN", "placeholder value for {run} in local parser report")
	includeLocal := flag.Bool("local", true, "include local parser report")
	outPath := flag.String("out", "", "optional Markdown output path; use - for stdout")
	updateDocPath := flag.String("update-doc", "", "optional Markdown document path with oracle compatibility report markers to update")
	flag.Parse()

	rows, err := readCompatTSV(*manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "compatreport: %v\n", err)
		os.Exit(1)
	}

	report := reportData{Manifest: summarizeRows(rows)}
	if *includeLocal {
		report.Local = summarizeLocal(rows, *runID)
	}
	if *referenceLogPath != "" {
		logText, err := readAllPath(*referenceLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "compatreport: %v\n", err)
			os.Exit(1)
		}
		if outcomes, ok := parseOutcomeLine(logText, "Oracle compatibility reference report:"); ok {
			report.Reference = outcomes
		}
		if outcomes, ok := parseOutcomeLine(logText, "Oracle compatibility reference report by family:"); ok {
			report.ReferenceByFamily = outcomes
		}
		if outcomes, ok := parseOutcomeLine(logText, "Oracle privileged compatibility reference report:"); ok {
			report.PrivilegedReference = outcomes
		}
		if outcomes, ok := parseOutcomeLine(logText, "Oracle privileged compatibility reference report by family:"); ok {
			report.PrivilegedReferenceByFamily = outcomes
		}
	}

	markdown := renderMarkdown(report)
	if *updateDocPath != "" {
		doc, err := readAllPath(*updateDocPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "compatreport: %v\n", err)
			os.Exit(1)
		}
		updated, err := replaceMarkedSection(doc, markdown)
		if err != nil {
			fmt.Fprintf(os.Stderr, "compatreport: update %s: %v\n", *updateDocPath, err)
			os.Exit(1)
		}
		if err := writeTextPath(*updateDocPath, updated); err != nil {
			fmt.Fprintf(os.Stderr, "compatreport: %v\n", err)
			os.Exit(1)
		}
	}
	if *outPath != "" {
		if err := writeTextPath(*outPath, markdown); err != nil {
			fmt.Fprintf(os.Stderr, "compatreport: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if *updateDocPath == "" {
		fmt.Print(markdown)
	}
}

func readCompatTSV(path string) ([]compatRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var header []string
	var rows []compatRow
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
				return nil, fmt.Errorf("%s: missing header", path)
			}
			continue
		}
		if len(parts) != len(header) {
			return nil, fmt.Errorf("%s: row has %d fields, want %d: %q", path, len(parts), len(header), line)
		}
		fields := make(map[string]string, len(header))
		for i, name := range header {
			fields[name] = parts[i]
		}
		rows = append(rows, compatRow{Fields: fields})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	if header == nil {
		return nil, fmt.Errorf("%s: empty TSV", path)
	}
	return rows, nil
}

func summarizeRows(rows []compatRow) manifestSummary {
	summary := manifestSummary{
		Total:            len(rows),
		Phases:           make(map[string]int),
		Families:         make(map[string]int),
		Expectations:     make(map[string]int),
		ReferenceClasses: make(map[string]int),
	}
	for _, row := range rows {
		summary.Phases[row.Fields["phase"]]++
		summary.Families[row.Fields["family"]]++
		summary.Expectations[row.Fields["expect"]]++
		summary.ReferenceClasses[row.Fields["reference_class"]]++
	}
	return summary
}

func summarizeLocal(rows []compatRow, runID string) map[string]int {
	outcomes := make(map[string]int)
	for _, row := range rows {
		expect := row.Fields["expect"]
		if expect == "unsafe" || expect == "version_dependent" {
			outcomes[expect]++
			continue
		}

		_, err := oracleparser.Parse(strings.ReplaceAll(row.Fields["sql"], "{run}", runID))
		omniAccepts := err == nil
		switch {
		case expect == "accept" && omniAccepts:
			outcomes["local_accept_match"]++
		case expect == "accept" && !omniAccepts:
			outcomes["local_too_strict"]++
		case expect == "reject" && !omniAccepts:
			outcomes["local_reject_match"]++
		case expect == "reject" && omniAccepts:
			outcomes["local_too_lenient"]++
		}
	}
	return outcomes
}

func parseOutcomeLine(logText string, prefix string) (map[string]int, bool) {
	for _, line := range strings.Split(logText, "\n") {
		idx := strings.Index(line, prefix)
		if idx < 0 {
			continue
		}
		counts := make(map[string]int)
		payload := strings.TrimSpace(line[idx+len(prefix):])
		for _, part := range strings.Split(payload, ",") {
			key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok || key == "" {
				continue
			}
			n, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			counts[key] = n
		}
		return counts, true
	}
	return nil, false
}

func renderMarkdown(report reportData) string {
	var b strings.Builder
	b.WriteString("# Oracle Compatibility Summary\n\n")

	b.WriteString("## Manifest Shape\n\n")
	b.WriteString("| Metric | Count |\n|---|---:|\n")
	writeMetric(&b, "total rows", report.Manifest.Total)
	writeCountMapWithSuffix(&b, "phase ", " rows", report.Manifest.Phases, orderedPhaseKeys(report.Manifest.Phases))
	writeCountMapWithSuffix(&b, "expected ", " rows", report.Manifest.Expectations, orderedExpectationKeys(report.Manifest.Expectations))
	writeCountMapWithSuffix(&b, "", " reference rows", report.Manifest.ReferenceClasses, orderedReferenceClassKeys(report.Manifest.ReferenceClasses))

	b.WriteString("\n## Family Distribution\n\n")
	b.WriteString("| Family | Rows |\n|---|---:|\n")
	writeCountMap(&b, report.Manifest.Families, orderedFamilyKeys(report.Manifest.Families))

	if len(report.Local) > 0 {
		b.WriteString("\n## Local Parser Report\n\n")
		b.WriteString("| Outcome | Rows |\n|---|---:|\n")
		writeCountMap(&b, report.Local, orderedOutcomeKeys(report.Local))
	}

	if len(report.Reference) > 0 {
		b.WriteString("\n## Oracle Reference Report\n\n")
		b.WriteString("| Outcome | Rows |\n|---|---:|\n")
		writeCountMap(&b, report.Reference, orderedOutcomeKeys(report.Reference))
	}

	if len(report.ReferenceByFamily) > 0 {
		b.WriteString("\n## Oracle Reference Report By Family\n\n")
		b.WriteString("| Family | Outcome | Rows |\n|---|---|---:|\n")
		writeByFamilyOutcomeMap(&b, report.ReferenceByFamily)
	}

	if len(report.PrivilegedReference) > 0 {
		b.WriteString("\n## Oracle Privileged Reference Report\n\n")
		b.WriteString("| Outcome | Rows |\n|---|---:|\n")
		writeCountMap(&b, report.PrivilegedReference, orderedOutcomeKeys(report.PrivilegedReference))
	}

	if len(report.PrivilegedReferenceByFamily) > 0 {
		b.WriteString("\n## Oracle Privileged Reference Report By Family\n\n")
		b.WriteString("| Family | Outcome | Rows |\n|---|---|---:|\n")
		writeByFamilyOutcomeMap(&b, report.PrivilegedReferenceByFamily)
	}

	return b.String()
}

func writeMetric(b *strings.Builder, label string, count int) {
	fmt.Fprintf(b, "| %s | %d |\n", label, count)
}

func writeCountMap(b *strings.Builder, counts map[string]int, keys []string) {
	for _, key := range keys {
		writeMetric(b, displayLabel(key), counts[key])
	}
}

func writeCountMapWithSuffix(b *strings.Builder, prefix string, suffix string, counts map[string]int, keys []string) {
	for _, key := range keys {
		writeMetric(b, prefix+displayLabel(key)+suffix, counts[key])
	}
}

func writeByFamilyOutcomeMap(b *strings.Builder, counts map[string]int) {
	for _, key := range orderedByFamilyOutcomeKeys(counts) {
		family, outcome, ok := strings.Cut(key, "/")
		if !ok {
			family = key
		}
		fmt.Fprintf(b, "| %s | %s | %d |\n", displayLabel(family), displayLabel(outcome), counts[key])
	}
}

func orderedPhaseKeys(counts map[string]int) []string {
	return orderedKeys(counts, []string{"A", "B", "C", "D", "E", "F"})
}

func orderedExpectationKeys(counts map[string]int) []string {
	return orderedKeys(counts, []string{"accept", "reject", "unsafe", "version_dependent"})
}

func orderedReferenceClassKeys(counts map[string]int) []string {
	return orderedKeys(counts, []string{"syntax", "semantic_catalog", "privilege_catalog", "fixture_sensitive", "version_dependent"})
}

func orderedFamilyKeys(counts map[string]int) []string {
	return orderedKeys(counts, []string{"select", "dml", "merge", "ddl", "plsql", "keyword", "admin", "utility"})
}

func orderedOutcomeKeys(counts map[string]int) []string {
	return orderedKeys(counts, []string{
		"local_accept_match",
		"local_reject_match",
		"match_accept",
		"match_reject",
		"local_too_strict",
		"local_too_lenient",
		"omni_too_strict",
		"omni_too_lenient",
		"oracle_semantic_catalog",
		"oracle_privilege_catalog",
		"unsafe",
		"version_dependent",
	})
}

func orderedByFamilyOutcomeKeys(counts map[string]int) []string {
	familyOrder := orderedFamilyKeys(mapFamilies(counts))
	outcomeOrder := []string{
		"match_accept",
		"match_reject",
		"omni_too_strict",
		"omni_too_lenient",
		"oracle_semantic_catalog",
		"oracle_privilege_catalog",
		"unsafe",
		"version_dependent",
	}
	seen := make(map[string]struct{}, len(counts))
	var keys []string
	for _, family := range familyOrder {
		for _, outcome := range outcomeOrder {
			key := family + "/" + outcome
			if _, ok := counts[key]; ok {
				keys = append(keys, key)
				seen[key] = struct{}{}
			}
		}
	}
	var rest []string
	for key := range counts {
		if _, ok := seen[key]; !ok {
			rest = append(rest, key)
		}
	}
	sort.Strings(rest)
	return append(keys, rest...)
}

func mapFamilies(counts map[string]int) map[string]int {
	families := make(map[string]int)
	for key, count := range counts {
		family, _, ok := strings.Cut(key, "/")
		if !ok {
			family = key
		}
		families[family] += count
	}
	return families
}

func orderedKeys(counts map[string]int, preferred []string) []string {
	seen := make(map[string]struct{}, len(counts))
	var keys []string
	for _, key := range preferred {
		if _, ok := counts[key]; ok {
			keys = append(keys, key)
			seen[key] = struct{}{}
		}
	}
	var rest []string
	for key := range counts {
		if _, ok := seen[key]; !ok {
			rest = append(rest, key)
		}
	}
	sort.Strings(rest)
	return append(keys, rest...)
}

func displayLabel(key string) string {
	return strings.ReplaceAll(key, "_", " ")
}

func replaceMarkedSection(doc string, report string) (string, error) {
	startIdx := strings.Index(doc, reportStartMarker)
	if startIdx < 0 {
		return "", errors.New("missing start marker")
	}
	endIdx := strings.Index(doc, reportEndMarker)
	if endIdx < 0 {
		return "", errors.New("missing end marker")
	}
	if endIdx < startIdx {
		return "", errors.New("end marker appears before start marker")
	}

	report = strings.TrimRight(report, "\n") + "\n"
	prefix := doc[:startIdx+len(reportStartMarker)]
	suffix := doc[endIdx:]
	return prefix + "\n" + report + suffix, nil
}

func readAllPath(path string) (string, error) {
	var r io.Reader
	if path == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("open %s: %w", path, err)
		}
		defer f.Close()
		r = f
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeTextPath(path string, text string) error {
	if path == "-" {
		_, err := fmt.Print(text)
		return err
	}
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
