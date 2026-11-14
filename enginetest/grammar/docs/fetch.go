package docs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// fetcher downloads pages, and keeps what it read in cache when cache is
// set, so a rerun of the converter reads the same pages without the network.
type fetcher struct {
	site   string // what a page's path is relative to
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
	url := f.site + path
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
