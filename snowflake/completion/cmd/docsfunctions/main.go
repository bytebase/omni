// Command docsfunctions records the built-in functions Snowflake's function
// index (https://docs.snowflake.com/en/sql-reference/functions-all) lists in
// snowflake/completion/functions.tsv, which completion offers by name. Run it
// from the repo root to bump the list:
//
//	go run ./snowflake/completion/cmd/docsfunctions
//
// With -page it reads a saved copy of the index instead of the network. The
// file's pin header holds the index URL, the SHA-256 of the rows it wrote,
// and the date it read them; the date moves only when the rows do.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const indexURL = "https://docs.snowflake.com/en/sql-reference/functions-all"

func main() {
	page := flag.String("page", "", "a saved copy of the function index; the network when empty")
	out := flag.String("out", "snowflake/completion/functions.tsv", "the file to write")
	date := flag.String("date", "", "the date to pin when the rows change; today when empty")
	flag.Parse()

	body, err := read(*page)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := functionRows(body)
	if err != nil {
		log.Fatal(err)
	}
	text, err := render(*out, rows, *date)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeFile(*out, text); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d functions\n", len(rows))
}

func writeFile(path, text string) error {
	return os.WriteFile(path, []byte(text), 0o644)
}

func read(page string) (string, error) {
	if page != "" {
		b, err := os.ReadFile(page)
		return string(b), err
	}
	req, err := http.NewRequest(http.MethodGet, indexURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "omni-docsfunctions")
	resp, err := (&http.Client{Timeout: time.Minute}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: %s", indexURL, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

var (
	indexRow = regexp.MustCompile(`(?s)<tr><td><a href="[^"]*">(.*?)</a></td><td>(.*?)</td><td>`)
	tag      = regexp.MustCompile(`<[^>]*>`)
	space    = regexp.MustCompile(`\s+`)
	name     = regexp.MustCompile(`^[A-Z_][A-Z0-9_$]*$`)
	variant  = regexp.MustCompile(`^(.*?) \(([^()]*)\)$`)
)

// functionRows returns name and summary rows for the functions the index
// lists that are called by a bare name. An index entry can name several
// functions (LENGTH, LEN) and a variant in parentheses (AI_COMPLETE (Single
// string)); an entry whose parentheses name a namespace (SNOWFLAKE.CORTEX, a
// system data metric function) is called qualified and is left out.
func functionRows(body string) ([][2]string, error) {
	matches := indexRow.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no rows in the function index; has its layout changed?")
	}
	seen := map[string]bool{}
	var rows [][2]string
	for _, m := range matches {
		entry := clean(m[1])
		summary := clean(m[2])
		if strings.HasSuffix(entry, "— Obsoleted") {
			continue
		}
		entry = strings.TrimSuffix(entry, " — Deprecated")
		for _, part := range strings.FieldsFunc(entry, func(r rune) bool { return r == ',' || r == '/' }) {
			part = strings.TrimSpace(part)
			if v := variant.FindStringSubmatch(part); v != nil {
				if namespaced(v[2]) {
					continue
				}
				part = v[1]
			}
			if !name.MatchString(part) || seen[part] {
				continue
			}
			seen[part] = true
			rows = append(rows, [2]string{part, summary})
		}
	}
	return rows, nil
}

func namespaced(paren string) bool {
	return strings.Contains(paren, "SNOWFLAKE") || strings.Contains(paren, "namespace") ||
		strings.Contains(paren, "system data metric function") || strings.Contains(paren, "SYS_CONTEXT")
}

func clean(s string) string {
	s = html.UnescapeString(tag.ReplaceAllString(s, ""))
	return strings.TrimSpace(space.ReplaceAllString(s, " "))
}

// render returns the file for rows, keeping the pinned date of out when the
// rows are the same.
func render(out string, rows [][2]string, date string) (string, error) {
	var b strings.Builder
	for _, r := range rows {
		b.WriteString(r[0] + "\t" + r[1] + "\n")
	}
	sum := sha256.Sum256([]byte(b.String()))
	digest := "sha256:" + hex.EncodeToString(sum[:])
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
		if old, err := os.ReadFile(out); err == nil {
			if pin, ok := readPin(string(old)); ok && pin.digest == digest {
				date = pin.date
			}
		}
	}
	return header(digest, date) + b.String(), nil
}

func header(digest, date string) string {
	return "# Snowflake's built-in functions called by a bare name, with the summary the\n" +
		"# function index gives, written by snowflake/completion/cmd/docsfunctions.\n" +
		"# Do not edit it: run the command again to bump the pin.\n" +
		"# pin: " + indexURL + " " + digest + " " + date + "\n" +
		"# columns: name, summary\n"
}

type pin struct{ digest, date string }

func readPin(text string) (pin, bool) {
	for _, line := range strings.Split(text, "\n") {
		if f := strings.Fields(strings.TrimPrefix(line, "# pin: ")); strings.HasPrefix(line, "# pin: ") && len(f) == 3 {
			return pin{digest: f[1], date: f[2]}, true
		}
	}
	return pin{}, false
}
