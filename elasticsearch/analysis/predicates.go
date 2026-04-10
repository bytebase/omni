package analysis

import "strings"

// directFieldClauses are clause types where the keys inside the clause object
// are field names. For example, {"match": {"email": "alice"}} -> field "email".
var directFieldClauses = map[string]bool{
	"match": true, "match_phrase": true, "match_phrase_prefix": true, "match_bool_prefix": true,
	"term": true, "terms": true, "terms_set": true,
	"range":    true,
	"wildcard": true, "prefix": true, "fuzzy": true, "regexp": true,
	"geo_distance": true, "geo_bounding_box": true, "geo_shape": true, "geo_polygon": true,
	"span_term": true,
}

// compoundClauses maps clause types to the keys that contain sub-queries.
var compoundClauses = map[string][]string{
	"bool":           {"must", "must_not", "should", "filter"},
	"nested":         {"query"},
	"has_child":      {"query"},
	"has_parent":     {"query"},
	"dis_max":        {"queries"},
	"constant_score": {"filter"},
	"boosting":       {"positive", "negative"},
	"function_score": {"query"},
}

// extractPredicateFields extracts field names from ES query clauses.
// It looks for the top-level "query" key, then recursively walks all recognized
// clause types to collect the dot-notation field paths used in predicates.
func extractPredicateFields(parsed map[string]any) []string {
	queryVal, ok := parsed["query"]
	if !ok {
		return nil
	}
	queryMap, ok := queryVal.(map[string]any)
	if !ok {
		return nil
	}
	var fields []string
	extractFieldsFromQuery(queryMap, &fields)
	return fields
}

// extractFieldsFromQuery recursively walks ES query clauses to collect field names.
func extractFieldsFromQuery(queryMap map[string]any, fields *[]string) {
	for clauseType, clauseVal := range queryMap {
		clauseObj, ok := clauseVal.(map[string]any)
		if !ok {
			// Value might be an array (e.g. bare array at top level).
			extractFieldsFromQueryArray(clauseVal, fields)
			continue
		}
		switch {
		case directFieldClauses[clauseType]:
			extractDirectFieldNames(clauseObj, fields)
		case clauseType == "exists":
			if fieldVal, ok := clauseObj["field"].(string); ok {
				*fields = append(*fields, fieldVal)
			}
		case compoundClauses[clauseType] != nil:
			extractCompoundClauseFields(clauseObj, compoundClauses[clauseType], fields)
		default:
			// Unknown clause type — recurse in case it wraps sub-queries.
			extractFieldsFromQuery(clauseObj, fields)
		}
	}
}

// extractFieldsFromQueryArray handles the case where a clause value is an array
// of query maps (e.g. the value of "must" in a bool query).
func extractFieldsFromQueryArray(val any, fields *[]string) {
	arr, ok := val.([]any)
	if !ok {
		return
	}
	for _, item := range arr {
		if m, ok := item.(map[string]any); ok {
			extractFieldsFromQuery(m, fields)
		}
	}
}

// extractDirectFieldNames extracts field names from a direct field clause object.
// Keys starting with "_" and the "boost" parameter are excluded.
func extractDirectFieldNames(clauseObj map[string]any, fields *[]string) {
	for fieldName := range clauseObj {
		if !strings.HasPrefix(fieldName, "_") && fieldName != "boost" {
			*fields = append(*fields, fieldName)
		}
	}
}

// extractCompoundClauseFields recurses into sub-query keys of a compound clause.
func extractCompoundClauseFields(clauseObj map[string]any, subQueryKeys []string, fields *[]string) {
	for _, sqKey := range subQueryKeys {
		sqVal, ok := clauseObj[sqKey]
		if !ok {
			continue
		}
		switch sq := sqVal.(type) {
		case map[string]any:
			extractFieldsFromQuery(sq, fields)
		case []any:
			extractFieldsFromQueryArray(sq, fields)
		default:
		}
	}
}
