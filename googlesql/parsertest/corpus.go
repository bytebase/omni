// Package parsertest provides corpus-extraction helpers shared by the
// googlesql parser's corpus-closure gates (legacy_corpus_test.go and
// official_corpus_test.go).
//
// It deliberately depends on NOTHING in googlesql/parser — it only knows how to
// (a) lift fenced ```sql blocks out of the truth1 markdown docs and (b) walk a
// corpus tree for .sql / .md files. The corpus tests live in package parser and
// drive parser.Parse themselves; keeping the extraction here avoids an import
// cycle and keeps the two corpus gates byte-identical in how they read inputs.
package parsertest

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SQLBlock is one fenced ```sql code block lifted from a truth1 markdown file,
// tagged with the file it came from, its 0-based block index within that file
// (the stable key the skip-list uses), and the nearest preceding `## ID: title`
// section header (for human-readable diagnostics — e.g. "DDL-013: CREATE
// EXTERNAL SCHEMA").
//
// Text is the raw block body verbatim (newlines preserved, fences stripped). It
// may hold one statement, several `;`-separated statements, leading/trailing
// `--`/`#` comments, or — in the lexical/datatype reference pages —
// illustrative fragments that are not complete statements at all. The caller
// (official_corpus_test.go) feeds Text to parser.Parse, which splits and parses
// every contained statement; non-statement fragments are handled by the
// test's explicit skip-list.
type SQLBlock struct {
	File         string // path relative to the corpus root, slash-separated
	Index        int    // 0-based ```sql block index within File
	SectionID    string // e.g. "DDL-013" (empty if the block has no `## ID:` header)
	SectionTitle string // e.g. "CREATE EXTERNAL SCHEMA"
	Text         string // raw block body (fences stripped)
}

// Key is the stable skip-list key for a block: "<file>#<index>", e.g.
// "bigquery/ddl.md#13". Indices are assigned in document order and are stable as
// long as no ```sql block is inserted/removed earlier in the same file — the
// corpus gate's block-count assertion guards against that drift.
func (b SQLBlock) Key() string {
	return b.File + "#" + itoa(b.Index)
}

// Label is a human-readable identifier for test output: the Key plus the section
// header when present, e.g. "bigquery/ddl.md#13 (DDL-013: CREATE EXTERNAL
// SCHEMA)".
func (b SQLBlock) Label() string {
	if b.SectionID == "" {
		return b.Key()
	}
	return b.Key() + " (" + b.SectionID + ": " + b.SectionTitle + ")"
}

// ExtractSQLBlocks lifts every fenced ```sql block from a markdown document, in
// document order, tagging each with the nearest preceding `## ID: title`
// section header. Only blocks opened by a line whose trimmed text is exactly
// "```sql" are returned — bare ``` blocks (grammar/EBNF snippets in the truth1
// docs) are ignored. relFile is recorded as each block's File.
//
// The scanner is line-oriented and fence-matched on the trimmed line so leading
// indentation inside a block never terminates it; only a line that trims to
// "```" closes the current block.
func ExtractSQLBlocks(relFile, markdown string) []SQLBlock {
	var (
		blocks   []SQLBlock
		curID    string
		curTitle string
		idx      int
		inBlock  bool
		body     []string
	)
	sc := bufio.NewScanner(strings.NewReader(markdown))
	sc.Buffer(make([]byte, 0, 1<<20), 4<<20) // truth1 blocks are small, but be safe
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.TrimSpace(line)

		if inBlock {
			if trimmed == "```" {
				blocks = append(blocks, SQLBlock{
					File:         relFile,
					Index:        idx,
					SectionID:    curID,
					SectionTitle: curTitle,
					Text:         strings.Join(body, "\n"),
				})
				idx++
				inBlock = false
				body = nil
				continue
			}
			body = append(body, line)
			continue
		}

		switch {
		case trimmed == "```sql":
			inBlock = true
			body = nil
		case strings.HasPrefix(trimmed, "## "):
			curID, curTitle = parseSectionHeader(strings.TrimPrefix(trimmed, "## "))
		}
	}
	return blocks
}

// parseSectionHeader splits a truth1 section header body of the form
// "DDL-013: CREATE EXTERNAL SCHEMA" into its ID and title. The left side counts
// as a form ID only when it matches the truth1 convention `<UPPER>-<digits>`
// (DDL-013 / QUERY-001 / LEX-001 / DCL-002 — every form ID in the corpus).
// A prose heading ("## Note: ...", "## Totals") yields an empty ID and the whole
// text as the title, so it never masquerades as a form anchor.
func parseSectionHeader(h string) (id, title string) {
	if i := strings.Index(h, ":"); i >= 0 {
		left := strings.TrimSpace(h[:i])
		if isFormID(left) {
			return left, strings.TrimSpace(h[i+1:])
		}
	}
	return "", strings.TrimSpace(h)
}

// isFormID reports whether s matches the truth1 form-ID shape: one or more
// ASCII uppercase letters, a single '-', then one or more ASCII digits
// (e.g. "DDL-013"). No regexp — this is on the hot path of corpus extraction.
func isFormID(s string) bool {
	dash := strings.IndexByte(s, '-')
	if dash <= 0 || dash == len(s)-1 {
		return false
	}
	for _, c := range s[:dash] {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	for _, c := range s[dash+1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// WalkMarkdownFiles returns every .md file under root (recursively), excluding
// INDEX.md, as paths relative to root (slash-separated), sorted for
// deterministic iteration.
func WalkMarkdownFiles(root string) ([]string, error) {
	return walk(root, func(name string) bool {
		return strings.HasSuffix(name, ".md") && !strings.EqualFold(filepath.Base(name), "INDEX.md")
	})
}

// WalkSQLFiles returns every .sql file under root (recursively) as paths
// relative to root (slash-separated), sorted for deterministic iteration.
func WalkSQLFiles(root string) ([]string, error) {
	return walk(root, func(name string) bool {
		return strings.HasSuffix(name, ".sql")
	})
}

func walk(root string, keep func(string) bool) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !keep(path) {
			return nil
		}
		rel, rerr := filepath.Rel(root, path)
		if rerr != nil {
			return rerr
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// itoa is a tiny non-allocating-ish int formatter (avoids importing strconv just
// for two call sites, keeping the helper dependency-light).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
