// Command mongoshgrammar converts the official MongoDB Shell (mongodb-js/mongosh,
// Apache-2.0) TypeScript shell-api and shell-bson definitions into the ANTLR 4
// call-shape grammar mongo/testdata/MongoShell.g4 that mongo/parser's
// compat.tsv is generated from and judged by.
//
// From the repo root, to regenerate MongoShell.g4 (and update the mongosh row
// in mongo/testdata/pins.tsv) from a local mongosh checkout or tarball extract:
//
//	go run ./mongo/parser/cmd/mongoshgrammar -src /path/to/mongosh -update-pin
//
// Or fetch the pinned commit directly into a local cache directory:
//
//	go run ./mongo/parser/cmd/mongoshgrammar -cache /tmp/mongosh-v2.14.3 -update-pin
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultTag    = "v2.14.3"
	defaultCommit = "82c228b597c3221de59bdbfaff4c346622780cdd"
	defaultOut    = "mongo/testdata/MongoShell.g4"
	defaultPins   = "mongo/testdata/pins.tsv"
)

func main() {
	srcDir := flag.String("src", "", "directory containing a mongodb-js/mongosh source tree")
	cacheDir := flag.String("cache", "", "cache directory to download pinned mongodb-js/mongosh files when -src is empty")
	tag := flag.String("tag", defaultTag, "mongodb-js/mongosh release tag")
	commit := flag.String("commit", defaultCommit, "mongodb-js/mongosh commit SHA for -tag")
	outPath := flag.String("out", defaultOut, "output path for MongoShell.g4")
	pinsPath := flag.String("pins", defaultPins, "path to mongo/testdata/pins.tsv")
	updatePin := flag.Bool("update-pin", false, "update the mongosh role row in -pins")
	flag.Parse()

	sources, err := loadSources(*srcDir, *cacheDir, *commit)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
		os.Exit(1)
	}
	digest, err := filesDigest(sources)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
		os.Exit(1)
	}
	ex, err := extract(sources)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
		os.Exit(1)
	}
	grammarText, err := renderGrammar(*tag, *commit, digest, ex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outPath, []byte(grammarText), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
		os.Exit(1)
	}
	if *updatePin {
		b, err := os.ReadFile(*pinsPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
			os.Exit(1)
		}
		updated := updatePinRow(string(b), *commit, digest, *tag)
		if err := os.WriteFile(*pinsPath, []byte(updated), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "mongoshgrammar:", err)
			os.Exit(1)
		}
	}
	fmt.Printf("wrote %s (mongosh %s @ %s, %s)\n", *outPath, *tag, *commit, digest)
}

func loadSources(srcDir, cacheDir, commit string) (map[string][]byte, error) {
	if srcDir == "" && cacheDir == "" {
		return nil, fmt.Errorf("give -src <mongosh-dir> or -cache <dir>")
	}
	client := &http.Client{Timeout: 30 * time.Second}
	sources := make(map[string][]byte, len(vendorFiles))
	for _, rel := range vendorFiles {
		if srcDir != "" {
			b, err := os.ReadFile(filepath.Join(srcDir, filepath.FromSlash(rel)))
			if err != nil {
				return nil, err
			}
			sources[rel] = b
			continue
		}
		local := filepath.Join(cacheDir, commit, filepath.FromSlash(rel))
		if b, err := os.ReadFile(local); err == nil {
			sources[rel] = b
			continue
		}
		url := fmt.Sprintf("https://raw.githubusercontent.com/mongodb-js/mongosh/%s/%s", commit, rel)
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetch %s: %w", url, err)
		}
		b, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch %s: HTTP %d", url, resp.StatusCode)
		}
		if readErr != nil {
			return nil, fmt.Errorf("read %s: %w", url, readErr)
		}
		if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(local, b, 0o644); err != nil {
			return nil, err
		}
		sources[rel] = b
	}
	return sources, nil
}
