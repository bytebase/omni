package analysis

import (
	"encoding/json"
	"net/url"
	"path"
	"slices"
	"strings"
)

// requestIndices returns every index, alias, data stream or pattern a request
// names, as written and URL-decoded: the path's comma-separated target first,
// then those a multi-get's docs or a multi-search's header lines name.
func requestIndices(api MaskableAPI, urlPath, body string) []string {
	var names []string
	if target := extractIndex(urlPath); target != "" {
		// Elasticsearch decodes the path segment before it splits the list,
		// so users%2Corders names two indices.
		if decoded, err := url.PathUnescape(target); err == nil {
			target = decoded
		}
		for name := range strings.SplitSeq(target, ",") {
			if name != "" {
				names = append(names, name)
			}
		}
	}
	switch api {
	case APIMaskMGet:
		var parsed struct {
			Docs []struct {
				Index string `json:"_index"`
			} `json:"docs"`
		}
		if json.Unmarshal([]byte(body), &parsed) == nil {
			for _, doc := range parsed.Docs {
				if doc.Index != "" {
					names = append(names, doc.Index)
				}
			}
		}
	case APIMaskSearch:
		if !strings.Contains(urlPath, "_msearch") {
			break
		}
		// A multi-search body alternates a header line and a search body;
		// a header's "index" is a name, a comma-separated list or an array.
		lines := nonEmptyLines(body)
		for i := 0; i < len(lines); i += 2 {
			var header map[string]any
			if json.Unmarshal([]byte(lines[i]), &header) != nil {
				continue
			}
			switch v := header["index"].(type) {
			case string:
				for name := range strings.SplitSeq(v, ",") {
					if name != "" {
						names = append(names, name)
					}
				}
			case []any:
				names = append(names, toStringSlice(v)...)
			default:
			}
		}
	default:
	}
	return names
}

func nonEmptyLines(body string) []string {
	var lines []string
	for line := range strings.SplitSeq(body, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// ExpandIndices resolves the names a request reads (RequestAnalysis.Indices)
// against the indices a sync stores, as Elasticsearch expands a target with
// its default expand_wildcards=open: no name, _all or * is every known index;
// a name with a wildcard is the known indices it matches, those starting with
// a dot only when the pattern does; a name starting with - removes what it
// matches from the names before it; any other name is kept as written, since
// it may be an index, an alias or a data stream the sync does not store.
func ExpandIndices(names, known []string) []string {
	if len(names) == 0 {
		names = []string{"_all"}
	}
	var out []string
	for i, name := range names {
		switch {
		case name == "_all":
			out = append(out, visible(known)...)
		case i > 0 && strings.HasPrefix(name, "-"):
			pattern := strings.TrimPrefix(name, "-")
			out = slices.DeleteFunc(out, func(n string) bool { return wildcardMatch(pattern, n) })
		case strings.Contains(name, "*"):
			for _, n := range known {
				if strings.HasPrefix(n, ".") && !strings.HasPrefix(name, ".") {
					continue
				}
				if wildcardMatch(name, n) {
					out = append(out, n)
				}
			}
		default:
			out = append(out, name)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func visible(known []string) []string {
	var out []string
	for _, n := range known {
		if !strings.HasPrefix(n, ".") {
			out = append(out, n)
		}
	}
	return out
}

// wildcardMatch reports whether name matches pattern, where * is any run of
// characters. Elasticsearch's index patterns take * alone; ? is a literal
// that no index name may hold.
func wildcardMatch(pattern, name string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == name
	}
	// path.Match treats [ and \ specially; index names hold neither.
	ok, err := path.Match(strings.ReplaceAll(pattern, "?", `\?`), name)
	return err == nil && ok
}
