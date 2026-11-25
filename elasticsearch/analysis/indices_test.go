package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/elasticsearch/analysis"
)

func TestRequestIndices(t *testing.T) {
	for _, tc := range []struct {
		method, url, body string
		want              []string
	}{
		{"GET", "/users/_search", "", []string{"users"}},
		{"GET", "/_search", "", nil},
		{"GET", "/users,orders/_search", "", []string{"users", "orders"}},
		{"GET", "/users%2Corders/_search", "", []string{"users", "orders"}},
		{"GET", "/user*,-users_archive/_search", "", []string{"user*", "-users_archive"}},
		{"GET", "/users/_doc/1", "", []string{"users"}},
		{"POST", "/_mget", `{"docs":[{"_index":"users","_id":"1"},{"_index":"orders","_id":"2"}]}`, []string{"users", "orders"}},
		{"POST", "/_msearch", "{\"index\":\"users\"}\n{}\n{\"index\":[\"a\",\"b\"]}\n{}\n", []string{"users", "a", "b"}},
		{"POST", "/x/_msearch", "{}\n{}\n{\"index\":\"y,z\"}\n{}\n", []string{"x", "y", "z"}},
		{"PUT", "/users/_doc/1", `{}`, nil},
	} {
		got := analysis.AnalyzeRequest(tc.method, tc.url, tc.body).Indices
		if !slices.Equal(got, tc.want) {
			t.Errorf("%s %s: Indices = %q, want %q", tc.method, tc.url, got, tc.want)
		}
	}
}

func TestExpandIndices(t *testing.T) {
	known := []string{".internal", "orders", "users", "users_archive"}
	for _, tc := range []struct {
		names, want []string
	}{
		{nil, []string{"orders", "users", "users_archive"}},
		{[]string{"_all"}, []string{"orders", "users", "users_archive"}},
		{[]string{"*"}, []string{"orders", "users", "users_archive"}},
		{[]string{".*"}, []string{".internal"}},
		{[]string{"user*"}, []string{"users", "users_archive"}},
		{[]string{"user*", "-users_archive"}, []string{"users"}},
		{[]string{"*", "-user*"}, []string{"orders"}},
		{[]string{"people", "orders"}, []string{"orders", "people"}},
		{[]string{"nothing*"}, nil},
	} {
		if got := analysis.ExpandIndices(tc.names, known); !slices.Equal(got, tc.want) {
			t.Errorf("ExpandIndices(%q) = %q, want %q", tc.names, got, tc.want)
		}
	}
}

func TestSourceFiltering(t *testing.T) {
	for _, tc := range []struct {
		url, body          string
		includes, excludes []string
		disabled           bool
	}{
		{"/u/_search", `{"_source":"name"}`, []string{"name"}, nil, false},
		{"/u/_search", `{"_source":{"includes":"a","excludes":["b"]}}`, []string{"a"}, []string{"b"}, false},
		{"/u/_search", `{"_source":{"include":["a"],"exclude":"b"}}`, []string{"a"}, []string{"b"}, false},
		{"/u/_search?_source=false", `{}`, nil, nil, true},
		{"/u/_search?_source=a,b", `{}`, []string{"a", "b"}, nil, false},
		{"/u/_search?_source_excludes=b", `{"_source":["a","b"]}`, nil, []string{"b"}, false},
		{"/u/_search?_source=true", `{"_source":false}`, nil, nil, false},
		{"/u/_doc/1?_source_includes=a", "", []string{"a"}, nil, false},
		{"/u/_source/1?_source_excludes=a", "", nil, []string{"a"}, false},
	} {
		method := "POST"
		if tc.body == "" {
			method = "GET"
		}
		a := analysis.AnalyzeRequest(method, tc.url, tc.body)
		if !slices.Equal(a.SourceFields, tc.includes) || !slices.Equal(a.SourceExcludes, tc.excludes) || a.SourceDisabled != tc.disabled {
			t.Errorf("%s %s: includes %q excludes %q disabled %v, want %q %q %v", tc.url, tc.body, a.SourceFields, a.SourceExcludes, a.SourceDisabled, tc.includes, tc.excludes, tc.disabled)
		}
	}
}
