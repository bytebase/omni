package analysis

import (
	"encoding/json"
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
	extractSource(parsed, result)
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
func extractSource(parsed map[string]any, result *RequestAnalysis) {
	src, ok := parsed["_source"]
	if !ok {
		return
	}

	switch v := src.(type) {
	case bool:
		if !v {
			result.SourceDisabled = true
		}
	case []any:
		result.SourceFields = toStringSlice(v)
	case map[string]any:
		if includes, ok := v["includes"]; ok {
			if arr, ok := includes.([]any); ok {
				result.SourceFields = toStringSlice(arr)
			}
		}
	default:
	}
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
