package analysis

import (
	"encoding/json"
	"net/url"
	"slices"
	"strings"
)

// requestTargets returns the index expressions a request reads, each the
// list of names one search or document resolves, as written: the path's
// comma-separated target, URL-decoded; for a multi-get, each doc's _index
// (the path's target for a doc that names none); for a multi-search, each
// header's index (the path's target for a header that names none). An empty
// expression reads every index.
func requestTargets(api MaskableAPI, urlPath, body string) [][]string {
	var pathTarget []string
	if target := extractIndex(urlPath); target != "" {
		// Elasticsearch decodes the path segment before it splits the list,
		// so users%2Corders names two indices.
		if decoded, err := url.PathUnescape(target); err == nil {
			target = decoded
		}
		pathTarget = splitNames(target)
	}
	switch {
	case api == APIMaskMGet:
		var parsed struct {
			Docs []struct {
				Index string `json:"_index"`
			} `json:"docs"`
		}
		if json.Unmarshal([]byte(body), &parsed) != nil || len(parsed.Docs) == 0 {
			return [][]string{pathTarget}
		}
		var targets [][]string
		for _, doc := range parsed.Docs {
			if doc.Index != "" {
				targets = append(targets, []string{doc.Index})
			} else {
				targets = append(targets, pathTarget)
			}
		}
		return targets
	case api == APIMaskSearch && strings.Contains(urlPath, "_msearch"):
		// A multi-search body alternates a header line and a search body;
		// a header's "index" is a name, a comma-separated list or an array.
		lines := nonEmptyLines(body)
		var targets [][]string
		for i := 0; i < len(lines); i += 2 {
			var header map[string]any
			if json.Unmarshal([]byte(lines[i]), &header) != nil {
				header = nil
			}
			switch v := header["index"].(type) {
			case string:
				targets = append(targets, splitNames(v))
			case []any:
				targets = append(targets, toStringSlice(v))
			default:
				targets = append(targets, pathTarget)
			}
		}
		if len(targets) == 0 {
			return [][]string{pathTarget}
		}
		return targets
	default:
		return [][]string{pathTarget}
	}
}

func splitNames(list string) []string {
	var names []string
	for name := range strings.SplitSeq(list, ",") {
		if name != "" {
			names = append(names, name)
		}
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

// ExpandIndices resolves the index expressions a request reads
// (RequestAnalysis.Targets) against the indices a sync stores, which are the
// ones a search's default expand_wildcards=open reaches, and returns their
// union. Within one expression, as Elasticsearch expands it: no name, _all or
// * is every known index; a name with a wildcard is the known indices it
// matches; a name starting with - after the first removes what it matches
// from the names before it; any other name is kept as written, since it may
// be an index, an alias or a data stream the sync does not store.
func ExpandIndices(targets [][]string, known []string) []string {
	if len(targets) == 0 {
		targets = [][]string{nil}
	}
	var out []string
	for _, names := range targets {
		out = append(out, expandTarget(names, known)...)
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func expandTarget(names, known []string) []string {
	if len(names) == 0 {
		names = []string{"_all"}
	}
	var out []string
	for i, name := range names {
		switch {
		case name == "_all":
			out = append(out, known...)
		case i > 0 && strings.HasPrefix(name, "-"):
			pattern := strings.TrimPrefix(name, "-")
			out = slices.DeleteFunc(out, func(n string) bool { return wildcardMatch(pattern, n) })
		case strings.Contains(name, "*"):
			for _, n := range known {
				if wildcardMatch(name, n) {
					out = append(out, n)
				}
			}
		default:
			out = append(out, name)
		}
	}
	return out
}

// wildcardMatch reports whether name matches pattern, where * is any run of
// characters and every other character stands for itself, as in
// Elasticsearch's Regex.simpleMatch.
func wildcardMatch(pattern, name string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == name
	}
	if !strings.HasPrefix(name, parts[0]) {
		return false
	}
	name = name[len(parts[0]):]
	last := parts[len(parts)-1]
	for _, part := range parts[1 : len(parts)-1] {
		i := strings.Index(name, part)
		if i < 0 {
			return false
		}
		name = name[i+len(part):]
	}
	return strings.HasSuffix(name, last) && len(name) >= len(last)
}
