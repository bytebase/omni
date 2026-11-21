// Command specgrammar generates elasticsearch/testdata/ElasticsearchConsole.g4
// from Elastic's official elasticsearch-specification schema.json (Apache-2.0).
//
// Run from the repo root:
//
//	go run ./elasticsearch/parser/cmd/specgrammar -cache /tmp/es-spec
//
// Or from a local schema.json checkout:
//
//	go run ./elasticsearch/parser/cmd/specgrammar -schema /path/to/schema.json
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	defaultCommit = "fd5d063c4824de5d33618b7fca4df0e6c1b8489f"
	defaultTag    = "8.19"
	repoSlug      = "elastic/elasticsearch-specification"
	schemaRelPath = "output/schema/schema.json"
)

type specSchema struct {
	Endpoints []specEndpoint `json:"endpoints"`
	Types     []specType     `json:"types"`
}

type specEndpoint struct {
	Name    string      `json:"name"`
	Request *specTypeID `json:"request"`
	URLs    []specURL   `json:"urls"`
}

type specURL struct {
	Methods []string `json:"methods"`
	Path    string   `json:"path"`
}

type specTypeID struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type specType struct {
	Name  specTypeID  `json:"name"`
	Query []specField `json:"query"`
	Path  []specField `json:"path"`
}

type specField struct {
	Name string        `json:"name"`
	Type specFieldType `json:"type"`
}

type specFieldType struct {
	Kind string      `json:"kind"`
	Type *specTypeID `json:"type"`
}

func main() {
	schemaPath := flag.String("schema", "", "path to local schema.json (downloads from GitHub if empty)")
	commit := flag.String("commit", defaultCommit, "elastic/elasticsearch-specification commit SHA")
	tag := flag.String("tag", defaultTag, "version/branch label recorded in pins.tsv")
	cacheDir := flag.String("cache", "", "optional directory to cache downloaded schema.json")
	outPath := flag.String("out", filepath.Join("elasticsearch", "testdata", "ElasticsearchConsole.g4"), "output .g4 path")
	pinsPath := flag.String("pins", filepath.Join("elasticsearch", "testdata", "pins.tsv"), "pins.tsv path (used with -update-pin)")
	updatePin := flag.Bool("update-pin", false, "update the default row in pins.tsv")
	flag.Parse()

	raw, err := loadSchema(*schemaPath, *commit, *cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specgrammar: %v\n", err)
		os.Exit(1)
	}
	sum := sha256.Sum256(raw)
	digest := "sha256:" + hex.EncodeToString(sum[:])

	grammar, stats, err := convert(raw, *commit, *tag, digest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specgrammar: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "specgrammar: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outPath, []byte(grammar), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "specgrammar: %v\n", err)
		os.Exit(1)
	}
	if *updatePin {
		if err := writePin(*pinsPath, *commit, digest, *tag); err != nil {
			fmt.Fprintf(os.Stderr, "specgrammar: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("wrote %s (%s): %d endpoints, %d distinct paths, %d non-GET official pairs, %d query forms, %d relative forms, %d special forms\n",
		*outPath, digest, stats.endpoints, stats.distinctPaths, stats.nonGetPairs, stats.queryForms, stats.relativeForms, stats.specialForms)
}

func loadSchema(localPath, commit, cacheDir string) ([]byte, error) {
	if localPath != "" {
		return os.ReadFile(localPath)
	}
	var cached string
	if cacheDir != "" {
		cached = filepath.Join(cacheDir, commit+"-schema.json")
		if b, err := os.ReadFile(cached); err == nil {
			return b, nil
		}
	}
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repoSlug, commit, schemaRelPath)
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: HTTP %s", url, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}
	if cached != "" {
		if err := os.MkdirAll(cacheDir, 0o755); err == nil {
			_ = os.WriteFile(cached, b, 0o644)
		}
	}
	return b, nil
}

type convertStats struct {
	endpoints     int
	distinctPaths int
	nonGetPairs   int
	queryForms    int
	relativeForms int
	specialForms  int
}

var pathParamRE = regexp.MustCompile(`\{([^}]+)\}`)

// samplePath replaces {param} templates with deterministic, representative
// Elasticsearch path segment values.
func samplePath(raw string) (string, error) {
	sampled := pathParamRE.ReplaceAllStringFunc(raw, func(m string) string {
		name := m[1 : len(m)-1]
		switch name {
		case "index", "indices", "target", "source", "dest", "data_stream":
			return "my_index"
		case "version", "page", "size", "x", "y", "zoom", "timestamp":
			return "1"
		default:
			return "a"
		}
	})
	if strings.ContainsAny(sampled, " \t\r\n'\"") || strings.Contains(sampled, "/*") || strings.Contains(sampled, "--") {
		return "", fmt.Errorf("sampled path %q (from %q) contains characters incompatible with the grammar recognizer", sampled, raw)
	}
	return sampled, nil
}

func sampleQueryValue(f specField) string {
	typeName := ""
	if f.Type.Type != nil {
		typeName = strings.ToLower(f.Type.Type.Name)
	}
	switch {
	case typeName == "boolean":
		return "true"
	case typeName == "integer" || typeName == "long" || typeName == "short" || typeName == "byte" || typeName == "int":
		return "10"
	case typeName == "float" || typeName == "double":
		return "1.5"
	case typeName == "duration" || strings.HasSuffix(f.Name, "timeout") || f.Name == "keep_alive":
		return "30s"
	case f.Type.Kind == "array_of" || f.Name == "h" || f.Name == "s" || f.Name == "fields" || f.Name == "names":
		return "a,b"
	case f.Name == "format":
		return "json"
	case f.Name == "expand_wildcards":
		return "open"
	case f.Name == "refresh":
		return "wait_for"
	default:
		return "a"
	}
}

func nsRuleName(ns string) string {
	clean := strings.TrimLeft(ns, "_")
	if clean == "" {
		clean = "root"
	}
	parts := strings.Split(clean, "_")
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(strings.ToLower(p))
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
		}
	}
	return b.String() + "Path"
}

type pathEntry struct {
	sampled   string
	endpoints []string
}

func convert(raw []byte, commit, tag, digest string) (string, convertStats, error) {
	var schema specSchema
	if err := json.Unmarshal(raw, &schema); err != nil {
		return "", convertStats{}, fmt.Errorf("parse schema.json: %w", err)
	}
	if len(schema.Endpoints) == 0 {
		return "", convertStats{}, fmt.Errorf("schema.json has no endpoints")
	}

	typesByKey := make(map[specTypeID]specType, len(schema.Types))
	for _, t := range schema.Types {
		typesByKey[t.Name] = t
	}

	// Sort endpoints deterministically by name.
	endpoints := slices.Clone(schema.Endpoints)
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Name < endpoints[j].Name
	})

	// Track distinct sampled paths globally so each path belongs to exactly one
	// namespace rule (the first namespace alphabetically that declares it).
	globalPathOwner := map[string]string{}
	nsPaths := map[string]map[string][]string{}
	byMethod := map[string]map[string][]string{
		"POST":   {},
		"PUT":    {},
		"DELETE": {},
		"HEAD":   {},
	}

	// Collect one representative single-param and multi-param query path per namespace.
	type queryCandidate struct {
		ns      string
		path    string
		epName  string
		fields  []specField
	}
	nsQuery := map[string][]queryCandidate{}

	for _, ep := range endpoints {
		ns := "root"
		if dot := strings.IndexByte(ep.Name, '.'); dot >= 0 {
			ns = ep.Name[:dot]
		}
		var qFields []specField
		if ep.Request != nil {
			if t, ok := typesByKey[*ep.Request]; ok && len(t.Query) > 0 {
				qFields = slices.Clone(t.Query)
				sort.Slice(qFields, func(i, j int) bool { return qFields[i].Name < qFields[j].Name })
			}
		}

		urls := slices.Clone(ep.URLs)
		sort.Slice(urls, func(i, j int) bool { return urls[i].Path < urls[j].Path })

		for _, u := range urls {
			sp, err := samplePath(u.Path)
			if err != nil {
				return "", convertStats{}, err
			}
			owner, owned := globalPathOwner[sp]
			if !owned {
				globalPathOwner[sp] = ns
				owner = ns
			}
			if nsPaths[owner] == nil {
				nsPaths[owner] = map[string][]string{}
			}
			if !slices.Contains(nsPaths[owner][sp], ep.Name) {
				nsPaths[owner][sp] = append(nsPaths[owner][sp], ep.Name)
			}

			methods := slices.Clone(u.Methods)
			sort.Strings(methods)
			for _, m := range methods {
				m = strings.ToUpper(m)
				switch m {
				case "GET":
					// Covered by httpMethod specPath (GET is the first httpMethod).
				case "POST", "PUT", "DELETE", "HEAD":
					if !slices.Contains(byMethod[m][sp], ep.Name) {
						byMethod[m][sp] = append(byMethod[m][sp], ep.Name)
					}
				default:
					return "", convertStats{}, fmt.Errorf("endpoint %s: unexpected HTTP method %q", ep.Name, m)
				}
			}

			if len(qFields) > 0 {
				nsQuery[ns] = append(nsQuery[ns], queryCandidate{
					ns:     ns,
					path:   sp,
					epName: ep.Name,
					fields: qFields,
				})
			}
		}
	}

	// Build sorted namespace list, placing "root" first so '/' is the shortest
	// path derivation of specPath (producing clean rows like "GET /", "POST /").
	var namespaces []string
	for ns := range nsPaths {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)
	if idx := slices.Index(namespaces, "root"); idx > 0 {
		namespaces = append([]string{"root"}, append(namespaces[:idx], namespaces[idx+1:]...)...)
	}

	// Build deterministic query-parameter paths (1 single-param + up to 1 multi-param per namespace).
	querySet := map[string]string{}
	for _, ns := range namespaces {
		cands := nsQuery[ns]
		if len(cands) == 0 {
			continue
		}
		sort.Slice(cands, func(i, j int) bool {
			if cands[i].epName != cands[j].epName {
				return cands[i].epName < cands[j].epName
			}
			return cands[i].path < cands[j].path
		})
		// Pick first candidate for single-param form.
		c0 := cands[0]
		q1 := fmt.Sprintf("%s?%s=%s", c0.path, c0.fields[0].Name, sampleQueryValue(c0.fields[0]))
		querySet[q1] = c0.epName
		// Pick first candidate with >=2 query params for multi-param form.
		for _, c := range cands {
			if len(c.fields) >= 2 {
				q2 := fmt.Sprintf("%s?%s=%s&%s=%s",
					c.path,
					c.fields[0].Name, sampleQueryValue(c.fields[0]),
					c.fields[1].Name, sampleQueryValue(c.fields[1]))
				querySet[q2] = c.epName
				break
			}
		}
	}
	// Add canonical Kibana Dev Console query forms: bare flag (?v), comma-separated
	// list (?h=index,docs.count), and quoted query value (?q="a b").
	querySet["/_cat/indices?v"] = "cat.indices (bare boolean flag)"
	querySet["/_cat/indices?h=index,docs.count&s=index"] = "cat.indices (comma list)"
	querySet["/_search?q=\"a b\""] = "search (quoted query value)"

	var queryPaths []string
	for qp := range querySet {
		queryPaths = append(queryPaths, qp)
	}
	sort.Strings(queryPaths)

	// Dev Console relative paths (omitting leading '/'), derived from core
	// root/cat/cluster/indices/nodes endpoints in schema.json.
	relativeCandidates := []string{
		"/_aliases",
		"/_bulk",
		"/_cat/health",
		"/_cat/indices",
		"/_cat/nodes",
		"/_cluster/health",
		"/_cluster/settings",
		"/_cluster/state",
		"/_count",
		"/_mget",
		"/_msearch",
		"/_nodes/stats",
		"/_refresh",
		"/_search",
		"/_sql",
		"/my_index",
		"/my_index/_count",
		"/my_index/_doc",
		"/my_index/_doc/a",
		"/my_index/_mapping",
		"/my_index/_refresh",
		"/my_index/_search",
		"/my_index/_settings",
		"/my_index/_source/a",
		"/my_index/_update/a",
	}
	var relativePaths []string
	for _, p := range relativeCandidates {
		if _, ok := globalPathOwner[p]; ok {
			relativePaths = append(relativePaths, strings.TrimPrefix(p, "/"))
		}
	}
	sort.Strings(relativePaths)

	// Elasticsearch REST path forms exercising URL characters valid in index
	// names, multi-targets, cross-cluster search, and document IDs.
	specialPaths := []string{
		"/%3Clogs-%7Bnow%2Fd%7D%3E/_search",
		"/index1,index2/_search",
		"/logs-2026.*/_search",
		"/logs.2026.11/_search",
		"/my_index/_doc/a+b",
		"/my_index/_doc/a~1",
		"/my_index/_doc/user%40example.com",
		"/remote:my_index/_search",
	}

	var officialMethodRules []string
	for _, m := range []string{"POST", "PUT", "DELETE", "HEAD"} {
		if len(byMethod[m]) > 0 {
			officialMethodRules = append(officialMethodRules, strings.ToLower(m)+"Endpoint")
		}
	}

	// Render ANTLR 4 grammar.
	var b strings.Builder
	fmt.Fprintf(&b, `/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Elasticsearch Kibana Dev Console REST request-line grammar, generated by
// elasticsearch/parser/cmd/specgrammar from Elastic's official
// %s (%s @ %s,
// %s, %s).
//
// Each official REST endpoint path is emitted as a single literal token so
// enginetest/grammar renders single-line Console requests without inserting
// whitespace inside URL paths. httpMethod + consolePath models the Kibana Dev
// Console request-line grammar (any of the six Console HTTP methods with any
// official, query-parameterized, relative, or special-character path), while
// postEndpoint, putEndpoint, deleteEndpoint, and headEndpoint exercise every
// non-GET official (method, path) pair from the specification.
//
// Do not edit by hand: rerun elasticsearch/parser/cmd/specgrammar.

grammar ElasticsearchConsole;

request
    : httpMethod consolePath
`, repoSlug, tag, commit, schemaRelPath, digest)
	if len(officialMethodRules) > 0 {
		b.WriteString("    | officialEndpoint\n")
	}
	b.WriteString(`    ;

httpMethod
    : GET
    | POST
    | PUT
    | DELETE
    | HEAD
    | PATCH
    ;

consolePath
    : specPath
`)
	if len(queryPaths) > 0 {
		b.WriteString("    | queryPath\n")
	}
	if len(relativePaths) > 0 {
		b.WriteString("    | relativePath\n")
	}
	if len(specialPaths) > 0 {
		b.WriteString("    | specialPath\n")
	}
	b.WriteString("    ;\n")

	if len(officialMethodRules) > 0 {
		b.WriteString("\nofficialEndpoint\n")
		for i, r := range officialMethodRules {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s%s\n", sep, r)
		}
		b.WriteString("    ;\n")
	}

	b.WriteString("\nspecPath\n")

	for i, ns := range namespaces {
		sep := "    : "
		if i > 0 {
			sep = "    | "
		}
		fmt.Fprintf(&b, "%s%s\n", sep, nsRuleName(ns))
	}
	b.WriteString("    ;\n")

	totalDistinctPaths := 0
	for _, ns := range namespaces {
		entries := make([]pathEntry, 0, len(nsPaths[ns]))
		for p, eps := range nsPaths[ns] {
			sort.Strings(eps)
			entries = append(entries, pathEntry{sampled: p, endpoints: eps})
		}
		sort.Slice(entries, func(i, j int) bool {
			// Keep "/" first in rootPath so the shortest path derivation is "/".
			if entries[i].sampled == "/" {
				return true
			}
			if entries[j].sampled == "/" {
				return false
			}
			return entries[i].sampled < entries[j].sampled
		})
		totalDistinctPaths += len(entries)
		fmt.Fprintf(&b, "\n// Namespace: %s (%d paths)\n%s\n", ns, len(entries), nsRuleName(ns))
		for i, e := range entries {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s'%s' // %s\n", sep, escapeANTLR(e.sampled), strings.Join(e.endpoints, ", "))
		}
		b.WriteString("    ;\n")
	}

	nonGetPairs := 0
	for _, m := range []string{"POST", "PUT", "DELETE", "HEAD"} {
		if len(byMethod[m]) == 0 {
			continue
		}
		ruleName := strings.ToLower(m) + "Endpoint"
		entries := make([]pathEntry, 0, len(byMethod[m]))
		for p, eps := range byMethod[m] {
			sort.Strings(eps)
			entries = append(entries, pathEntry{sampled: p, endpoints: eps})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].sampled < entries[j].sampled })
		nonGetPairs += len(entries)
		fmt.Fprintf(&b, "\n// Official %s endpoints (%d paths)\n%s\n", m, len(entries), ruleName)
		for i, e := range entries {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s%s '%s' // %s\n", sep, m, escapeANTLR(e.sampled), strings.Join(e.endpoints, ", "))
		}
		b.WriteString("    ;\n")
	}

	if len(queryPaths) > 0 {
		fmt.Fprintf(&b, "\n// Official query-parameter call shapes (%d paths)\nqueryPath\n", len(queryPaths))
		for i, qp := range queryPaths {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s'%s' // %s\n", sep, escapeANTLR(qp), querySet[qp])
		}
		b.WriteString("    ;\n")
	}

	if len(relativePaths) > 0 {
		fmt.Fprintf(&b, "\n// Kibana Dev Console relative paths without leading slash (%d paths)\nrelativePath\n", len(relativePaths))
		for i, rp := range relativePaths {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s'%s'\n", sep, escapeANTLR(rp))
		}
		b.WriteString("    ;\n")
	}

	if len(specialPaths) > 0 {
		fmt.Fprintf(&b, "\n// Elasticsearch REST paths with wildcards, multi-targets, cross-cluster colons, and encoding (%d paths)\nspecialPath\n", len(specialPaths))
		for i, sp := range specialPaths {
			sep := "    : "
			if i > 0 {
				sep = "    | "
			}
			fmt.Fprintf(&b, "%s'%s'\n", sep, escapeANTLR(sp))
		}
		b.WriteString("    ;\n")
	}

	b.WriteString(`
// Kibana Dev Console HTTP methods (case-insensitive in parser.go).
GET    : [Gg] [Ee] [Tt] ;
POST   : [Pp] [Oo] [Ss] [Tt] ;
PUT    : [Pp] [Uu] [Tt] ;
DELETE : [Dd] [Ee] [Ll] [Ee] [Tt] [Ee] ;
HEAD   : [Hh] [Ee] [Aa] [Dd] ;
PATCH  : [Pp] [Aa] [Tt] [Cc] [Hh] ;
`)

	return b.String(), convertStats{
		endpoints:     len(endpoints),
		distinctPaths: totalDistinctPaths,
		nonGetPairs:   nonGetPairs,
		queryForms:    len(queryPaths),
		relativeForms: len(relativePaths),
		specialForms:  len(specialPaths),
	}, nil
}

func escapeANTLR(s string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s)
}

func writePin(path, commit, digest, tag string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	newRow := fmt.Sprintf("default\tgithub.com/%s@%s\t%s\t%s", repoSlug, commit, digest, tag)
	lines := strings.Split(string(b), "\n")
	var out []string
	replaced := false
	for _, line := range lines {
		if strings.HasPrefix(line, "default\t") {
			out = append(out, newRow)
			replaced = true
			continue
		}
		out = append(out, line)
	}
	if !replaced {
		// Insert default before the first non-comment row (e.g. driver).
		var inserted []string
		for _, line := range out {
			if !replaced && line != "" && !strings.HasPrefix(line, "#") {
				inserted = append(inserted, newRow)
				replaced = true
			}
			inserted = append(inserted, line)
		}
		if !replaced {
			inserted = append(inserted, newRow)
		}
		out = inserted
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}
