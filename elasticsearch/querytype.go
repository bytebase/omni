package elasticsearch

import "strings"

// QueryType classifies an Elasticsearch REST request for cost and routing decisions.
type QueryType int

const (
	QueryTypeUnknown         QueryType = iota
	QueryTypeSelect                    // Read-only data query
	QueryTypeSelectInfoSchema          // Read-only cluster/node metadata
	QueryTypeDML                       // Document write (index, update, delete)
	QueryTypeDDL                       // Index/alias/mapping management
	QueryTypeExplain                   // Query explanation
)

// readOnlyPostEndpoints lists URL patterns that are read-only despite using POST.
// These endpoints accept a request body for query parameters but do not modify data.
var readOnlyPostEndpoints = []string{
	"_search",
	"_msearch",
	"_count",
	"_field_caps",
	"_validate/query",
	"_search_shards",
	"_terms_enum",
	"_termvectors",
	"_mtermvectors",
	"_search_template",
	"_msearch_template",
	"_render/template",
	"_rank_eval",
	"_knn_search",
	"_search_mvt",
	"_sql",
	"_esql/query",
	"_eql/search",
	"_fleet/_search",
	"_fleet/_msearch",
	"_graph/explore",
	"rollup_search",
	"async_search",
}

// infoSchemaEndpoints lists URL patterns for cluster/node metadata queries.
var infoSchemaEndpoints = []string{
	"_cat/",
	"_cluster/",
	"_nodes/",
}

// dmlPostEndpoints lists URL patterns for document write operations via POST.
var dmlPostEndpoints = []string{
	"_doc",
	"_create/",
	"_update/",
	"_bulk",
	"_delete_by_query",
	"_update_by_query",
	"_reindex",
}

// ClassifyRequest determines the QueryType for an Elasticsearch REST API request.
func ClassifyRequest(method, url string) QueryType {
	method = strings.ToUpper(method)
	urlLower := strings.ToLower(url)

	switch method {
	case "HEAD":
		return QueryTypeSelect

	case "GET":
		if isInfoSchemaURL(urlLower) {
			return QueryTypeSelectInfoSchema
		}
		return QueryTypeSelect

	case "DELETE":
		if isDocumentURL(urlLower) {
			return QueryTypeDML
		}
		return QueryTypeDDL

	case "PUT":
		if isDocumentWriteURL(urlLower) {
			return QueryTypeDML
		}
		return QueryTypeDDL

	case "PATCH":
		return QueryTypeDML

	case "POST":
		return classifyPostRequest(urlLower)

	default:
		return QueryTypeUnknown
	}
}

func classifyPostRequest(url string) QueryType {
	// Explain
	if strings.Contains(url, "_explain/") || strings.HasSuffix(url, "_explain") {
		return QueryTypeExplain
	}

	// Info schema (metadata reads)
	if isInfoSchemaURL(url) {
		return QueryTypeSelectInfoSchema
	}

	// Read-only POST endpoints
	if isReadOnlyPostURL(url) {
		return QueryTypeSelect
	}

	// Document writes (DML)
	if isDMLPostURL(url) {
		return QueryTypeDML
	}

	// Default: DDL for other POST operations (index management, admin, etc.)
	return QueryTypeDDL
}

func isInfoSchemaURL(url string) bool {
	for _, pattern := range infoSchemaEndpoints {
		if strings.Contains(url, pattern) {
			return true
		}
	}
	return false
}

func isReadOnlyPostURL(url string) bool {
	for _, pattern := range readOnlyPostEndpoints {
		if strings.Contains(url, pattern) {
			return true
		}
	}
	return false
}

func isDMLPostURL(url string) bool {
	for _, pattern := range dmlPostEndpoints {
		if strings.Contains(url, pattern) {
			return true
		}
	}
	return false
}

func isDocumentURL(url string) bool {
	// DELETE /{index}/_doc/{id} is DML (document deletion)
	return strings.Contains(url, "_doc/")
}

func isDocumentWriteURL(url string) bool {
	// PUT /{index}/_doc/{id}, PUT /{index}/_create/{id}, PUT /_bulk are DML
	return strings.Contains(url, "_doc/") ||
		strings.Contains(url, "_doc") ||
		strings.Contains(url, "_create/") ||
		strings.Contains(url, "_bulk")
}
