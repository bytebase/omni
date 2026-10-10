package parser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the known-gap manifests of the Oracle dictionary tests from the current results")

// oracleDictionaryVersion is the release of the pinned Oracle image whose
// dictionary the known-gap manifests were recorded on.
const oracleDictionaryVersion = "23.26.3.0.0"

var plsqlDictionaryManifest = filepath.Join("testdata", "coverage", "plsql_dictionary_known_failures.tsv")

// dictionaryUnit is one PL/SQL unit read from DBA_SOURCE.
type dictionaryUnit struct {
	key  string // OWNER.NAME:TYPE
	kind string // DBA_SOURCE type
	text string // CREATE OR REPLACE followed by the stored source
}

// TestOraclePLSQLDictionary parses the PL/SQL that ships with Oracle: every
// unwrapped, VALID procedure, function, package, package body, type body, and
// trigger of an Oracle-maintained schema in the pinned image, about 875 units
// of code the engine compiled. Each unit goes through Split and ParseRange,
// as Bytebase runs a script, and must come out as one segment that parses.
//
// The units that fail today are recorded with a category in
// plsql_dictionary_known_failures.tsv. A unit that starts failing is a
// regression; a recorded unit that now parses must leave the manifest, so the
// list only shrinks. Regenerate it with -update and review the diff. The
// source text itself is Oracle's and is read at test time, never committed.
func TestOraclePLSQLDictionary(t *testing.T) {
	units := loadOracleDictionaryUnits(t)
	if len(units) < 800 {
		t.Fatalf("read %d dictionary units, want at least 800: the DBA_SOURCE query no longer matches the image", len(units))
	}

	known := map[string]coverageRow{}
	if !*update {
		rows, _ := readCoverageTSV(t, plsqlDictionaryManifest)
		for _, row := range rows {
			known[row.Key] = row
		}
	}

	type failure struct {
		unit     dictionaryUnit
		category string
		where    string
		message  string
	}
	var failures []failure
	byCategory := map[string]int{}
	var newFailures, fixed int
	for _, u := range units {
		category, where, message, ok := parseDictionaryUnit(u)
		if ok {
			if _, wasKnown := known[u.key]; wasKnown {
				fixed++
				t.Errorf("FIXED %s now parses: remove it from %s", u.key, plsqlDictionaryManifest)
			}
			continue
		}
		failures = append(failures, failure{unit: u, category: category, where: where, message: message})
		byCategory[category]++
		if _, wasKnown := known[u.key]; !wasKnown && !*update {
			newFailures++
			t.Errorf("NEW FAILURE %s at %s: %s (%s)", u.key, where, message, category)
		}
	}

	t.Logf("Oracle %s dictionary: %d units, %d fail (%.1f%%), %d new, %d fixed",
		oracleDictionaryVersion, len(units), len(failures),
		100*float64(len(failures))/float64(len(units)), newFailures, fixed)
	var categories []string
	for c := range byCategory {
		categories = append(categories, c)
	}
	sort.Slice(categories, func(i, j int) bool {
		if byCategory[categories[i]] != byCategory[categories[j]] {
			return byCategory[categories[i]] > byCategory[categories[j]]
		}
		return categories[i] < categories[j]
	})
	for _, c := range categories {
		t.Logf("  %4d  %s", byCategory[c], c)
	}

	if *update {
		var b strings.Builder
		b.WriteString("unit\tcategory\tposition\terror\n")
		for _, f := range failures {
			fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", f.unit.key, f.category, f.where, f.message)
		}
		if err := os.WriteFile(plsqlDictionaryManifest, []byte(b.String()), 0o644); err != nil {
			t.Fatalf("write %s: %v", plsqlDictionaryManifest, err)
		}
		t.Logf("wrote %s (%d units)", plsqlDictionaryManifest, len(failures))
	}
}

// loadOracleDictionaryUnits reads the units from the container started by
// startOracleDB, after checking that it runs the release the manifest was
// recorded on. It does not honor ORACLE_PARSER_REF_* DSNs: the manifest
// describes the pinned image's dictionary, not another server's.
func loadOracleDictionaryUnits(t *testing.T) []dictionaryUnit {
	t.Helper()
	o := startOracleDB(t)

	var version string
	if err := o.adminDB.QueryRowContext(o.ctx,
		`SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'`).Scan(&version); err != nil {
		t.Fatalf("read Oracle version: %v", err)
	}
	if version != oracleDictionaryVersion {
		t.Fatalf("Oracle %s is running; the dictionary manifests were recorded on %s. Rerun with -update and review the drift.",
			version, oracleDictionaryVersion)
	}

	// The dba_objects and dba_users filters are semi-joins: a join repeats
	// rows for objects listed in more than one container.
	rows, err := o.adminDB.QueryContext(o.ctx, `
SELECT s.owner, s.name, s.type, s.line, s.text
  FROM dba_source s
 WHERE s.type IN ('PROCEDURE', 'FUNCTION', 'PACKAGE', 'PACKAGE BODY', 'TYPE BODY', 'TRIGGER')
   AND s.owner IN (SELECT username FROM dba_users WHERE oracle_maintained = 'Y')
   AND (s.owner, s.name, s.type) IN
       (SELECT owner, object_name, object_type FROM dba_objects WHERE status = 'VALID')
   AND (s.owner, s.name, s.type) NOT IN
       (SELECT owner, name, type FROM dba_source WHERE line = 1 AND LOWER(text) LIKE '%wrapped%')
 ORDER BY s.owner, s.type, s.name, s.line`)
	if err != nil {
		t.Fatalf("query DBA_SOURCE: %v", err)
	}
	defer rows.Close()

	var units []dictionaryUnit
	var cur *dictionaryUnit
	var body strings.Builder
	lastLine := 0
	flush := func() {
		if cur != nil {
			cur.text = "CREATE OR REPLACE " + body.String()
			units = append(units, *cur)
		}
		body.Reset()
	}
	for rows.Next() {
		var owner, name, kind, text string
		var line int
		if err := rows.Scan(&owner, &name, &kind, &line, &text); err != nil {
			t.Fatalf("scan DBA_SOURCE: %v", err)
		}
		key := owner + "." + name + ":" + kind
		if cur == nil || cur.key != key {
			flush()
			cur = &dictionaryUnit{key: key, kind: kind}
			lastLine = 0
		}
		if line == lastLine {
			continue
		}
		lastLine = line
		body.WriteString(text)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate DBA_SOURCE: %v", err)
	}
	flush()
	return units
}

// parseDictionaryUnit splits and parses one unit. On failure it returns the
// category of the construct at the error, the error's line:column within the
// unit, and the message.
func parseDictionaryUnit(u dictionaryUnit) (category, where, message string, ok bool) {
	var segments []Segment
	for _, seg := range Split(u.text) {
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}
	if len(segments) != 1 {
		return classifyDictionaryFailure(u, ""), "1:1", fmt.Sprintf("split into %d segments", len(segments)), false
	}
	_, err := ParseRange(u.text, segments[0].ByteStart, segments[0].ByteEnd)
	if err == nil {
		return "", "", "", true
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		return "unclassified", "1:1", err.Error(), false
	}
	pos := pe.Position
	if pos < 0 || pos > len(u.text) {
		pos = 0
	}
	lineStart := strings.LastIndexByte(u.text[:pos], '\n') + 1
	lineEnd := strings.IndexByte(u.text[pos:], '\n')
	if lineEnd < 0 {
		lineEnd = len(u.text)
	} else {
		lineEnd += pos
	}
	line := strings.Count(u.text[:pos], "\n") + 1
	where = fmt.Sprintf("%d:%d", line, len([]rune(u.text[lineStart:pos]))+1)
	return classifyDictionaryFailure(u, u.text[lineStart:lineEnd]), where, pe.Message, false
}

// dictionaryCategories name the construct behind a failure from the source
// line the parser stopped on. First match wins.
var dictionaryCategories = []struct {
	name    string
	kinds   string // DBA_SOURCE types the rule applies to; empty for all
	pattern *regexp.Regexp
}{
	{"conditional compilation ($IF/$END)", "", regexp.MustCompile(`(?i)^\s*\$|\$(if|then|elsif|else|end|error)\b|\$\$`)},
	{"ACCESSIBLE BY clause", "", regexp.MustCompile(`(?i)accessible\s+by`)},
	{"CHARACTER SET ANY_CS / %CHARSET", "", regexp.MustCompile(`(?i)character\s+set`)},
	{"PIPELINED / AGGREGATE USING implementation type", "", regexp.MustCompile(`(?i)(pipelined|aggregate|parallel_enable)\s+using|^\s*using\b`)},
	{"PARALLEL_ENABLE (PARTITION ...)", "", regexp.MustCompile(`(?i)parallel_enable\s*\(`)},
	{"package initialization section", "PACKAGE BODY", regexp.MustCompile(`(?i)^\s*begin\b`)},
	{"REF type", "", regexp.MustCompile(`(?i)\bref\s`)},
	{"Oracle-internal syntax (PRAGMA INTERFACE, ... parameters)", "", regexp.MustCompile(`(?i)pragma\s+interface|\.\.\.`)},
	{"JSON function clause", "", regexp.MustCompile(`(?i)\breturning\b|json_|^\s*key\b`)},
	{"trigger CALL with arguments", "TRIGGER", regexp.MustCompile(`(?i)^\s*call\b|\(\s*:(new|old)\.`)},
	{"qualified expression (=>)", "", regexp.MustCompile(`=>`)},
	{"type length expression", "", regexp.MustCompile(`(?i)\(\s*[a-z_0-9]*\s*[*+-]`)},
	{"SQL%BULK_EXCEPTIONS", "", regexp.MustCompile(`(?i)sql%bulk`)},
	{"type body method modifiers", "TYPE BODY", regexp.MustCompile(`(?i)\b(final|instantiable|overriding)\b`)},
}

func classifyDictionaryFailure(u dictionaryUnit, line string) string {
	if line == "" && strings.Contains(strings.ToUpper(u.text), "$END") {
		return "conditional compilation ($IF/$END)"
	}
	for _, c := range dictionaryCategories {
		if c.kinds != "" && c.kinds != u.kind {
			continue
		}
		if c.pattern.MatchString(line) {
			return c.name
		}
	}
	return "unclassified"
}
