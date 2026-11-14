// Command docsgrammar converts the Syntax sections of Snowflake's SQL
// command reference to snowflake/testdata/SnowflakeDocs.g4, the grammar
// snowflake/parser's compat.tsv is generated from and judged by, and pins
// what it read in snowflake/testdata/pins.tsv: the index page, a SHA-256 of
// every Syntax block it read, and the date it read them. Run it from the
// repo root to bump the pin, or after a change to docs-fixes.tsv:
//
//	go run ./snowflake/parser/cmd/docsgrammar -cache /tmp/snowflake-docs
//
// With -cache it reads pages from the directory, and keeps the ones it
// fetches there. The date in the pin moves only when the blocks do.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	cache := flag.String("cache", "", "read pages from, and keep fetched pages in, this directory")
	testdata := flag.String("testdata", filepath.Join("snowflake", "testdata"), "the directory of SnowflakeDocs.g4, docs-fixes.tsv and pins.tsv")
	date := flag.String("date", "", "the date the pages were read, when it is not today (UTC)")
	flag.Parse()

	fixes, err := readFixes(filepath.Join(*testdata, "docs-fixes.tsv"))
	if err != nil {
		fail(err)
	}
	f := &fetcher{cache: *cache, client: &http.Client{Timeout: 60 * time.Second}, pause: 200 * time.Millisecond}
	index, err := f.get(indexPath)
	if err != nil {
		fail(err)
	}
	digest := sha256.New()
	var pages []pageGrammar
	for _, p := range commandPages(index) {
		doc, err := f.get(p)
		if err != nil {
			fail(err)
		}
		slug := path.Base(p)
		blocks := syntaxBlocks(doc)
		fmt.Fprintf(digest, "%s\n%s\n", slug, strings.Join(blocks, "\n\x00\n"))
		pages = append(pages, convertPage(slug, blocks, fixes))
	}
	for _, fx := range fixes {
		if !fx.used {
			fail(fmt.Errorf("docs-fixes.tsv: %s: %q is not on the page any more", fx.page, fx.find))
		}
	}

	pin, err := pinFor(filepath.Join(*testdata, "pins.tsv"), "sha256:"+hex.EncodeToString(digest.Sum(nil)), *date)
	if err != nil {
		fail(err)
	}
	g4 := writeG4(strings.ReplaceAll(pin, "\t", " "), pages)
	if err := os.WriteFile(filepath.Join(*testdata, "SnowflakeDocs.g4"), []byte(g4), 0o644); err != nil {
		fail(err)
	}
	report(pages)
}

// readFixes reads docs-fixes.tsv: page, text, replacement, reason. \n and
// \t in the text and replacement stand for a newline and a tab.
func readFixes(file string) ([]*fix, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	unescape := strings.NewReplacer(`\n`, "\n", `\t`, "\t")
	var out []*fix
	for i, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) != 4 {
			return nil, fmt.Errorf("%s:%d: want page, text, replacement and reason", file, i+1)
		}
		out = append(out, &fix{page: cols[0], find: unescape.Replace(cols[1]), replace: unescape.Replace(cols[2]), reason: cols[3]})
	}
	return out, nil
}

// pinFor returns the default row of pins.tsv for blocks with digest, and
// writes it: the pinned row as it is when the digest is the same, or a
// row dated today (or date) when the docs changed.
func pinFor(file, digest, date string) (string, error) {
	header := "# The docs this engine's committed answers come from. default judges compat.tsv:\n" +
		"# the Syntax blocks of the pages its image (the command index) lists, their\n" +
		"# SHA-256 as snowflake/parser/cmd/docsgrammar reads them, and the date they\n" +
		"# were read, which is the pin's source tag.\n" +
		"# columns: role, image, digest, source tag\n"
	if b, err := os.ReadFile(file); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			cols := strings.Split(line, "\t")
			if len(cols) == 4 && cols[0] == "default" && cols[2] == digest && date == "" {
				return line, nil
			}
		}
	}
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	row := strings.Join([]string{"default", docsSite + indexPath, digest, date}, "\t")
	return row, os.WriteFile(file, []byte(header+row+"\n"), 0o644)
}

// report prints what was converted, and why the rest was not.
func report(pages []pageGrammar) {
	var withSyntax, withRule, forms, frags, failed int
	reasons := map[string]int{}
	for _, pg := range pages {
		if len(pg.forms)+len(pg.failed) > 0 {
			withSyntax++
		}
		if len(pg.forms) > 0 {
			withRule++
		}
		forms += len(pg.forms)
		frags += len(pg.frags)
		failed += len(pg.failed)
		for _, f := range pg.failed {
			r := f.reason
			if i := strings.Index(r, " fragment "); i >= 0 {
				r = r[:i] + " fragment"
			}
			reasons[r]++
		}
	}
	fmt.Printf("%d pages, %d with syntax, %d with a rule; %d alternatives, %d fragment rules; %d forms or definitions not converted\n",
		len(pages), withSyntax, withRule, forms, frags, failed)
	var rs []string
	for r := range reasons {
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool { return reasons[rs[i]] > reasons[rs[j]] })
	for _, r := range rs {
		fmt.Printf("  %4d %s\n", reasons[r], r)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "docsgrammar:", err)
	os.Exit(1)
}
