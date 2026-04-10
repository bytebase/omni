package parsertest

import (
	"testing"

	"github.com/bytebase/omni/elasticsearch"
)

func TestClassifyRequest(t *testing.T) {
	cases := []struct {
		name   string
		method string
		url    string
		want   elasticsearch.QueryType
	}{
		// HEAD → always Select
		{
			name:   "HEAD root",
			method: "HEAD",
			url:    "/",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "HEAD index",
			method: "HEAD",
			url:    "/my-index",
			want:   elasticsearch.QueryTypeSelect,
		},

		// GET non-info-schema → Select
		{
			name:   "GET index search",
			method: "GET",
			url:    "/my-index/_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "GET document",
			method: "GET",
			url:    "/my-index/_doc/1",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "GET root",
			method: "GET",
			url:    "/",
			want:   elasticsearch.QueryTypeSelect,
		},

		// GET with info schema patterns → SelectInfoSchema
		{
			name:   "GET _cat/indices",
			method: "GET",
			url:    "/_cat/indices",
			want:   elasticsearch.QueryTypeSelectInfoSchema,
		},
		{
			name:   "GET _cluster/health",
			method: "GET",
			url:    "/_cluster/health",
			want:   elasticsearch.QueryTypeSelectInfoSchema,
		},
		{
			name:   "GET _nodes/stats",
			method: "GET",
			url:    "/_nodes/stats",
			want:   elasticsearch.QueryTypeSelectInfoSchema,
		},

		// DELETE with _doc/ → DML
		{
			name:   "DELETE document",
			method: "DELETE",
			url:    "/my-index/_doc/1",
			want:   elasticsearch.QueryTypeDML,
		},

		// DELETE without _doc/ → DDL
		{
			name:   "DELETE index",
			method: "DELETE",
			url:    "/my-index",
			want:   elasticsearch.QueryTypeDDL,
		},
		{
			name:   "DELETE alias",
			method: "DELETE",
			url:    "/_aliases",
			want:   elasticsearch.QueryTypeDDL,
		},

		// PUT with document write patterns → DML
		{
			name:   "PUT _doc with id",
			method: "PUT",
			url:    "/my-index/_doc/1",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "PUT _create with id",
			method: "PUT",
			url:    "/my-index/_create/1",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "PUT _bulk",
			method: "PUT",
			url:    "/_bulk",
			want:   elasticsearch.QueryTypeDML,
		},

		// PUT without document patterns → DDL
		{
			name:   "PUT index",
			method: "PUT",
			url:    "/my-index",
			want:   elasticsearch.QueryTypeDDL,
		},
		{
			name:   "PUT mapping",
			method: "PUT",
			url:    "/my-index/_mapping",
			want:   elasticsearch.QueryTypeDDL,
		},
		{
			name:   "PUT settings",
			method: "PUT",
			url:    "/my-index/_settings",
			want:   elasticsearch.QueryTypeDDL,
		},

		// PATCH → always DML
		{
			name:   "PATCH document",
			method: "PATCH",
			url:    "/my-index/_doc/1",
			want:   elasticsearch.QueryTypeDML,
		},

		// POST _explain → Explain
		{
			name:   "POST _explain suffix",
			method: "POST",
			url:    "/my-index/_explain",
			want:   elasticsearch.QueryTypeExplain,
		},
		{
			name:   "POST _explain with id",
			method: "POST",
			url:    "/my-index/_explain/1",
			want:   elasticsearch.QueryTypeExplain,
		},

		// POST info schema → SelectInfoSchema
		{
			name:   "POST _cat",
			method: "POST",
			url:    "/_cat/indices",
			want:   elasticsearch.QueryTypeSelectInfoSchema,
		},
		{
			name:   "POST _cluster",
			method: "POST",
			url:    "/_cluster/settings",
			want:   elasticsearch.QueryTypeSelectInfoSchema,
		},

		// POST read-only endpoints → Select
		{
			name:   "POST _search",
			method: "POST",
			url:    "/my-index/_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _msearch",
			method: "POST",
			url:    "/_msearch",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _count",
			method: "POST",
			url:    "/my-index/_count",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _field_caps",
			method: "POST",
			url:    "/_field_caps",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _validate/query",
			method: "POST",
			url:    "/my-index/_validate/query",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _terms_enum",
			method: "POST",
			url:    "/my-index/_terms_enum",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _termvectors",
			method: "POST",
			url:    "/my-index/_termvectors",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _mtermvectors",
			method: "POST",
			url:    "/_mtermvectors",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _search_template",
			method: "POST",
			url:    "/my-index/_search_template",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _msearch_template",
			method: "POST",
			url:    "/_msearch_template",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _render/template",
			method: "POST",
			url:    "/_render/template",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _rank_eval",
			method: "POST",
			url:    "/_rank_eval",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _knn_search",
			method: "POST",
			url:    "/my-index/_knn_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _search_mvt",
			method: "POST",
			url:    "/my-index/_search_mvt/field/0/0/0",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _sql",
			method: "POST",
			url:    "/_sql",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _esql/query",
			method: "POST",
			url:    "/_esql/query",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _eql/search",
			method: "POST",
			url:    "/_eql/search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _fleet/_search",
			method: "POST",
			url:    "/_fleet/_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _fleet/_msearch",
			method: "POST",
			url:    "/_fleet/_msearch",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST _graph/explore",
			method: "POST",
			url:    "/my-index/_graph/explore",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST rollup_search",
			method: "POST",
			url:    "/my-index/rollup_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "POST async_search",
			method: "POST",
			url:    "/my-index/async_search",
			want:   elasticsearch.QueryTypeSelect,
		},

		// POST DML endpoints → DML
		{
			name:   "POST _doc",
			method: "POST",
			url:    "/my-index/_doc",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _create",
			method: "POST",
			url:    "/my-index/_create/1",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _update",
			method: "POST",
			url:    "/my-index/_update/1",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _bulk",
			method: "POST",
			url:    "/_bulk",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _delete_by_query",
			method: "POST",
			url:    "/my-index/_delete_by_query",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _update_by_query",
			method: "POST",
			url:    "/my-index/_update_by_query",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "POST _reindex",
			method: "POST",
			url:    "/_reindex",
			want:   elasticsearch.QueryTypeDML,
		},

		// POST default → DDL
		{
			name:   "POST index create",
			method: "POST",
			url:    "/my-new-index",
			want:   elasticsearch.QueryTypeDDL,
		},
		{
			name:   "POST _aliases",
			method: "POST",
			url:    "/_aliases",
			want:   elasticsearch.QueryTypeDDL,
		},
		{
			name:   "POST _shrink",
			method: "POST",
			url:    "/my-index/_shrink/my-shrunk-index",
			want:   elasticsearch.QueryTypeDDL,
		},

		// Unknown method
		{
			name:   "CONNECT unknown",
			method: "CONNECT",
			url:    "/",
			want:   elasticsearch.QueryTypeUnknown,
		},
		{
			name:   "OPTIONS unknown",
			method: "OPTIONS",
			url:    "/",
			want:   elasticsearch.QueryTypeUnknown,
		},

		// Case-insensitive method matching
		{
			name:   "lowercase get",
			method: "get",
			url:    "/my-index/_search",
			want:   elasticsearch.QueryTypeSelect,
		},
		{
			name:   "mixed case Post",
			method: "Post",
			url:    "/_bulk",
			want:   elasticsearch.QueryTypeDML,
		},
		{
			name:   "lowercase delete",
			method: "delete",
			url:    "/my-index/_doc/1",
			want:   elasticsearch.QueryTypeDML,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := elasticsearch.ClassifyRequest(tc.method, tc.url)
			if got != tc.want {
				t.Errorf("ClassifyRequest(%q, %q) = %v, want %v", tc.method, tc.url, got, tc.want)
			}
		})
	}
}
