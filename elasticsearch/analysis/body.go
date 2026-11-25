package analysis

import (
	"bytes"
	"encoding/json"
	"net/url"
	"slices"
	"strings"
)

// analyzeRequestBody parses the JSON body and extracts fields and blocked
// features.
func analyzeRequestBody(body string) *RequestAnalysis {
	result := &RequestAnalysis{}

	parsed := make(map[string]any)
	if json.Unmarshal([]byte(body), &parsed) != nil {
		return result
	}

	result.BlockedFeatures = detectBlockedFeatures(parsed)
	extractSource(body, parsed, result)
	result.RequestedFields = extractStringArray(parsed, "fields")
	result.HighlightFields = extractHighlightFields(parsed)
	result.SortFields = extractSortFields(parsed)
	result.HasInnerHits = containsKey(parsed, "inner_hits")
	result.PredicateFields = extractPredicateFields(parsed)

	return result
}

// detectBlockedFeatures checks for top-level keys that indicate blocked
// features.
func detectBlockedFeatures(parsed map[string]any) []BlockedFeature {
	var blocked []BlockedFeature

	if _, ok := parsed["aggs"]; ok {
		blocked = append(blocked, BlockedFeatureAggs)
	} else if _, ok := parsed["aggregations"]; ok {
		blocked = append(blocked, BlockedFeatureAggs)
	}
	if _, ok := parsed["suggest"]; ok {
		blocked = append(blocked, BlockedFeatureSuggest)
	}
	if _, ok := parsed["script_fields"]; ok {
		blocked = append(blocked, BlockedFeatureScriptFields)
	}
	if _, ok := parsed["runtime_mappings"]; ok {
		blocked = append(blocked, BlockedFeatureRuntimeMappings)
	}
	if _, ok := parsed["stored_fields"]; ok {
		blocked = append(blocked, BlockedFeatureStoredFields)
	}
	if _, ok := parsed["docvalue_fields"]; ok {
		blocked = append(blocked, BlockedFeatureDocvalueFields)
	}

	return blocked
}

// extractSource handles the three forms of _source in the request body.
func extractSource(body string, parsed map[string]any, result *RequestAnalysis) {
	src, ok := parsed["_source"]
	if !ok {
		return
	}

	switch v := src.(type) {
	case bool:
		if !v {
			result.SourceDisabled = true
		}
	case string:
		result.SourceFields = []string{v}
	case []any:
		result.SourceFields = toStringSlice(v)
	case map[string]any:
		result.SourceFields, result.SourceExcludes = sourceObject(body)
	default:
	}
}

// sourceObject reads a _source object's include and exclude lists from the
// body as written. Each spelling (includes or include, excludes or exclude)
// sets its list, so the last one written wins, as Elasticsearch's
// FetchSourceContext parses it; a generic decode loses the keys' order.
func sourceObject(body string) (includes, excludes []string) {
	var top struct {
		Source json.RawMessage `json:"_source"`
	}
	if json.Unmarshal([]byte(body), &top) != nil {
		return nil, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(top.Source))
	if tok, err := decoder.Token(); err != nil || tok != json.Delim('{') {
		return nil, nil
	}
	for decoder.More() {
		tok, err := decoder.Token()
		if err != nil {
			return includes, excludes
		}
		key, _ := tok.(string)
		var value any
		if err := decoder.Decode(&value); err != nil {
			return includes, excludes
		}
		switch key {
		case "includes", "include":
			includes = stringOrList(value)
		case "excludes", "exclude":
			excludes = stringOrList(value)
		}
	}
	return includes, excludes
}

// stringOrList reads a _source filter value: one field or a list of them.
func stringOrList(v any) []string {
	switch v := v.(type) {
	case string:
		return []string{v}
	case []any:
		return toStringSlice(v)
	default:
		return nil
	}
}

// sourceParams reads the URL's _source filtering, which Elasticsearch takes on
// every API that returns a document: _source (true, false or a
// comma-separated include list), _source_includes and _source_excludes. When
// the URL holds any of them, they replace the body's _source.
func sourceParams(query string, result *RequestAnalysis) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return
	}
	list := func(key string) []string {
		var out []string
		for _, v := range values[key] {
			for field := range strings.SplitSeq(v, ",") {
				if field != "" {
					out = append(out, field)
				}
			}
		}
		return out
	}
	_, source := values["_source"]
	_, includes := values["_source_includes"]
	_, excludes := values["_source_excludes"]
	if !source && !includes && !excludes {
		return
	}
	result.SourceDisabled, result.SourceFields, result.SourceExcludes = false, nil, nil
	if v := values["_source"]; len(v) > 0 {
		switch v[len(v)-1] {
		case "false":
			result.SourceDisabled = true
		case "true", "":
		default:
			result.SourceFields = list("_source")
		}
	}
	if includes {
		result.SourceFields = list("_source_includes")
	}
	result.SourceExcludes = list("_source_excludes")
}

// extractStringArray extracts a string array value from a top-level key.
func extractStringArray(parsed map[string]any, key string) []string {
	val, ok := parsed[key]
	if !ok {
		return nil
	}
	arr, ok := val.([]any)
	if !ok {
		return nil
	}
	return toStringSlice(arr)
}

// extractHighlightFields extracts field names from highlight.fields.
func extractHighlightFields(parsed map[string]any) []string {
	highlight, ok := parsed["highlight"]
	if !ok {
		return nil
	}
	hlMap, ok := highlight.(map[string]any)
	if !ok {
		return nil
	}
	fields, ok := hlMap["fields"]
	if !ok {
		return nil
	}
	fieldsMap, ok := fields.(map[string]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(fieldsMap))
	for k := range fieldsMap {
		result = append(result, k)
	}
	slices.Sort(result)
	return result
}

// extractSortFields extracts field names from the sort array.
func extractSortFields(parsed map[string]any) []string {
	sortVal, ok := parsed["sort"]
	if !ok {
		return nil
	}
	arr, ok := sortVal.([]any)
	if !ok {
		return nil
	}
	var fields []string
	for _, item := range arr {
		switch v := item.(type) {
		case string:
			if strings.HasPrefix(v, "_") {
				continue
			}
			fields = append(fields, v)
		case map[string]any:
			for k := range v {
				if !strings.HasPrefix(k, "_") {
					fields = append(fields, k)
				}
			}
		default:
		}
	}
	return fields
}

// containsKey recursively searches for a key in a nested map structure.
func containsKey(data map[string]any, key string) bool {
	for k, v := range data {
		if k == key {
			return true
		}
		if containsKeyInValue(v, key) {
			return true
		}
	}
	return false
}

// containsKeyInValue checks if any nested map within v contains the given key.
func containsKeyInValue(v any, key string) bool {
	switch child := v.(type) {
	case map[string]any:
		return containsKey(child, key)
	case []any:
		for _, item := range child {
			if m, ok := item.(map[string]any); ok && containsKey(m, key) {
				return true
			}
		}
	default:
	}
	return false
}

// toStringSlice converts a []any to a []string, skipping non-string elements.
func toStringSlice(arr []any) []string {
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
