// Package analysis classifies Elasticsearch REST API requests for masking
// and extracts field references from request bodies.
package analysis

import "strings"

// MaskableAPI classifies an Elasticsearch endpoint for data-masking purposes.
type MaskableAPI int

const (
	// APIUnsupported means the request cannot return user data (e.g. PUT, DELETE).
	APIUnsupported MaskableAPI = iota
	// APIMaskSearch means the request is a search that may return masked fields.
	APIMaskSearch
	// APIMaskGetDoc means the request fetches a single document via _doc.
	APIMaskGetDoc
	// APIMaskGetSource means the request fetches document source via _source.
	APIMaskGetSource
	// APIMaskMGet means the request uses the multi-get API (_mget).
	APIMaskMGet
	// APIMaskExplain means the request uses the explain API (_explain).
	APIMaskExplain
	// APIBlocked means the endpoint cannot be safely masked (scroll, SQL, etc.).
	APIBlocked
)

// BlockedFeature identifies a request body feature that prevents safe masking.
type BlockedFeature int

const (
	BlockedFeatureAggs            BlockedFeature = iota // "aggs" / "aggregations"
	BlockedFeatureSuggest                               // "suggest"
	BlockedFeatureScriptFields                          // "script_fields"
	BlockedFeatureRuntimeMappings                       // "runtime_mappings"
	BlockedFeatureStoredFields                          // "stored_fields"
	BlockedFeatureDocvalueFields                        // "docvalue_fields"
)

// BlockedFeatureNames maps BlockedFeature values to human-readable names.
var BlockedFeatureNames = map[BlockedFeature]string{
	BlockedFeatureAggs:            "aggregations",
	BlockedFeatureSuggest:         "suggest",
	BlockedFeatureScriptFields:    "script_fields",
	BlockedFeatureRuntimeMappings: "runtime_mappings",
	BlockedFeatureStoredFields:    "stored_fields",
	BlockedFeatureDocvalueFields:  "docvalue_fields",
}

// RequestAnalysis is the result of analysing a single Elasticsearch request.
type RequestAnalysis struct {
	// API is the classified endpoint type.
	API MaskableAPI
	// Index is the target index name extracted from the URL (may be empty).
	Index string
	// BlockedFeatures lists body features that prevent safe masking.
	BlockedFeatures []BlockedFeature
	// SourceDisabled is true when "_source": false was set in the body.
	SourceDisabled bool
	// SourceFields lists fields from the _source include list.
	SourceFields []string
	// RequestedFields lists fields from the top-level "fields" array.
	RequestedFields []string
	// HighlightFields lists field names from highlight.fields.
	HighlightFields []string
	// SortFields lists non-metadata field names from the sort array.
	SortFields []string
	// HasInnerHits is true when any nested inner_hits key was found.
	HasInnerHits bool
	// PredicateFields lists field names referenced in the query predicate.
	PredicateFields []string
}

// blockedEndpoints lists URL patterns that must be rejected when masking is
// active.  These overlap with read patterns but cannot be safely masked.
// Includes both Elasticsearch and OpenSearch equivalents.
var blockedEndpoints = []string{
	// Elasticsearch endpoints.
	"_async_search",
	"_search/scroll",
	"_search_template",
	"_msearch/template",
	"_sql",
	"_eql/search",
	"_esql/query",
	"_terms_enum",
	"_termvectors",
	"_mtermvectors",
	"_knn_search",
	// OpenSearch equivalents (plugin-based paths).
	"_plugins/_asynchronous_search",
	"_plugins/_sql",
	"_plugins/_ppl",
}

// ClassifyMaskableAPI determines the MaskableAPI type and target index for an
// Elasticsearch request.
func ClassifyMaskableAPI(method, url string) (MaskableAPI, string) {
	method = strings.ToUpper(method)

	// Strip query parameters.
	path := url
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}

	// Strip leading slash.
	path = strings.TrimPrefix(path, "/")

	// Only GET and POST can return user data.
	switch method {
	case "GET":
		if blocked, index := isBlockedEndpoint(path); blocked {
			return APIBlocked, index
		}
		return classifyGetForMasking(path)
	case "POST":
		if blocked, index := isBlockedEndpoint(path); blocked {
			return APIBlocked, index
		}
		return classifyPostForMasking(path)
	default:
		// HEAD, PUT, DELETE, PATCH do not return user data.
		return APIUnsupported, ""
	}
}

// isBlockedEndpoint checks whether the path matches a blocked endpoint pattern.
// Returns true and the extracted index if blocked.
func isBlockedEndpoint(path string) (bool, string) {
	for _, pattern := range blockedEndpoints {
		if strings.Contains(path, pattern) {
			return true, extractIndex(path)
		}
	}
	return false, ""
}

// classifyGetForMasking classifies a GET request path for masking.
func classifyGetForMasking(path string) (MaskableAPI, string) {
	for part := range strings.SplitSeq(path, "/") {
		switch part {
		case "_search", "_msearch":
			return APIMaskSearch, extractIndex(path)
		case "_doc":
			return APIMaskGetDoc, extractIndex(path)
		case "_source":
			return APIMaskGetSource, extractIndex(path)
		case "_mget":
			return APIMaskMGet, extractIndex(path)
		}
	}
	return APIUnsupported, ""
}

// classifyPostForMasking classifies a POST request path for masking.
func classifyPostForMasking(path string) (MaskableAPI, string) {
	for part := range strings.SplitSeq(path, "/") {
		switch part {
		case "_search", "_msearch":
			return APIMaskSearch, extractIndex(path)
		case "_mget":
			return APIMaskMGet, extractIndex(path)
		case "_explain":
			return APIMaskExplain, extractIndex(path)
		}
	}
	return APIUnsupported, ""
}

// extractIndex returns the first path segment if it does not start with "_".
// This extracts the index name from ES URL paths like "<index>/_search".
func extractIndex(path string) string {
	if path == "" {
		return ""
	}
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 0 {
		return ""
	}
	first := parts[0]
	if strings.HasPrefix(first, "_") {
		return ""
	}
	return first
}

// AnalyzeRequest combines URL classification with body analysis for an ES
// request.  Body analysis is only performed for search and explain APIs.
func AnalyzeRequest(method, url, body string) *RequestAnalysis {
	api, index := ClassifyMaskableAPI(method, url)
	result := &RequestAnalysis{
		API:   api,
		Index: index,
	}

	// Only analyze the body for APIs that accept a query body.
	switch api {
	case APIMaskSearch, APIMaskExplain:
		bodyResult := analyzeRequestBody(body)
		result.BlockedFeatures = bodyResult.BlockedFeatures
		result.SourceFields = bodyResult.SourceFields
		result.SourceDisabled = bodyResult.SourceDisabled
		result.RequestedFields = bodyResult.RequestedFields
		result.HighlightFields = bodyResult.HighlightFields
		result.SortFields = bodyResult.SortFields
		result.HasInnerHits = bodyResult.HasInnerHits
		result.PredicateFields = bodyResult.PredicateFields
	default:
		// No body analysis for other API types.
	}

	return result
}
