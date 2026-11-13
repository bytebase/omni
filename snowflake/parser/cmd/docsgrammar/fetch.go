package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// docsSite is where the command reference lives. The index lists every
// command page by a link relative to the site root.
const (
	docsSite  = "https://docs.snowflake.com/en"
	indexPath = "/sql-reference/sql-all"
)

var commandLink = regexp.MustCompile(`href="(/(?:en/)?sql-reference/sql/[a-z0-9_-]+)"`)

// commandPages returns the command pages the index links to, in the
// index's order, each once.
func commandPages(index string) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range commandLink.FindAllStringSubmatch(index, -1) {
		p := strings.TrimPrefix(m[1], "/en")
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

// fetcher downloads pages, and keeps what it read in cache when cache is
// set, so a rerun of the converter reads the same pages without the network.
type fetcher struct {
	cache  string
	client *http.Client
	pause  time.Duration
}

func (f *fetcher) get(path string) (string, error) {
	var cached string
	if f.cache != "" {
		cached = filepath.Join(f.cache, strings.ReplaceAll(strings.Trim(path, "/"), "/", "__")+".html")
		if b, err := os.ReadFile(cached); err == nil {
			return string(b), nil
		}
	}
	url := docsSite + path
	var body string
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 || f.pause > 0 {
			time.Sleep(f.pause * time.Duration(attempt+1))
		}
		body, err = f.download(url)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	if cached != "" {
		if err := os.MkdirAll(f.cache, 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(cached, []byte(body), 0o644); err != nil {
			return "", err
		}
	}
	return body, nil
}

func (f *fetcher) download(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "omni-docsgrammar")
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

var (
	preBlock   = regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`)
	anyTag     = regexp.MustCompile(`(?s)<[^>]+>`)
	nextH2     = regexp.MustCompile(`<h2[ >]`)
	syntaxSecs = regexp.MustCompile(`data-section="syntax[a-z0-9-]*"`)
)

// syntaxBlocks returns the text of every code block in a page's Syntax
// sections: Syntax, or Syntax for ... where a page has several, each from
// its heading to the next second-level heading, so the subsections some
// pages split their syntax into are included. A page with none, such as an
// overview, returns none.
func syntaxBlocks(doc string) []string {
	var out []string
	end := 0
	for _, loc := range syntaxSecs.FindAllStringIndex(doc, -1) {
		if loc[0] < end {
			continue // a subsection of the section just read
		}
		sec := doc[loc[0]:]
		// Skip the section's own <h2>, then stop at the next one.
		if h := strings.Index(sec, "</h2>"); h >= 0 {
			if j := nextH2.FindStringIndex(sec[h:]); j != nil {
				sec = sec[:h+j[0]]
			}
		}
		end = loc[0] + len(sec)
		for _, m := range preBlock.FindAllStringSubmatch(sec, -1) {
			text := html.UnescapeString(anyTag.ReplaceAllString(m[1], ""))
			text = strings.ReplaceAll(text, "\u00a0", " ")
			out = append(out, strings.TrimRight(text, " \n\t"))
		}
	}
	return out
}
