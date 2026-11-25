package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/elasticsearch/analysis"
)

func TestRequestTargets(t *testing.T) {
	for _, tc := range []struct {
		method, url, body string
		want              [][]string
	}{
		{"GET", "/users/_search", "", [][]string{{"users"}}},
		{"GET", "/_search", "", [][]string{nil}},
		{"GET", "/users,orders/_search", "", [][]string{{"users", "orders"}}},
		{"GET", "/users%2Corders/_search", "", [][]string{{"users", "orders"}}},
		{"GET", "/user*,-users_archive/_search", "", [][]string{{"user*", "-users_archive"}}},
		{"GET", "/users/_doc/1", "", [][]string{{"users"}}},
		{"POST", "/_mget", `{"docs":[{"_index":"users","_id":"1"},{"_index":"orders","_id":"2"}]}`, [][]string{{"users"}, {"orders"}}},
		{"POST", "/u/_mget", `{"docs":[{"_id":"1"},{"_index":"orders","_id":"2"}]}`, [][]string{{"u"}, {"orders"}}},
		{"POST", "/_msearch", "{\"index\":\"users\"}\n{}\n{\"index\":[\"a\",\"b\"]}\n{}\n{}\n{}\n", [][]string{{"users"}, {"a", "b"}, nil}},
		{"POST", "/x/_msearch", "{}\n{}\n{\"index\":\"y,z\"}\n{}\n", [][]string{{"x"}, {"y", "z"}}},
		{"PUT", "/users/_doc/1", `{}`, nil},
	} {
		got := analysis.AnalyzeRequest(tc.method, tc.url, tc.body).Targets
		if !slices.EqualFunc(got, tc.want, slices.Equal) {
			t.Errorf("%s %s: Targets = %q, want %q", tc.method, tc.url, got, tc.want)
		}
	}
}

func TestExpandIndices(t *testing.T) {
	known := []string{".customers", "logs[eu]-1", "logsp-1", "orders", "users", "users_archive"}
	for _, tc := range []struct {
		targets [][]string
		want    []string
	}{
		{nil, known},
		{[][]string{{"_all"}}, known},
		{[][]string{{"*"}}, known},
		{[][]string{{".*"}}, []string{".customers"}},
		{[][]string{{"user*"}}, []string{"users", "users_archive"}},
		{[][]string{{"user*", "-users_archive"}}, []string{"users"}},
		{[][]string{{"*", "-user*", "-logs*", "-.*"}}, []string{"orders"}},
		{[][]string{{"people", "orders"}}, []string{"orders", "people"}},
		{[][]string{{"nothing*"}}, nil},
		{[][]string{{"logs[eu]*"}}, []string{"logs[eu]-1"}},
		{[][]string{{"u*s"}}, []string{"users"}},
		// Each expression expands on its own: one search's exclusion
		// leaves another's names alone.
		{[][]string{{"users"}, {"*", "-users"}}, known},
		{[][]string{{"orders"}, nil}, known},
	} {
		if got := analysis.ExpandIndices(tc.targets, known); !slices.Equal(got, tc.want) {
			t.Errorf("ExpandIndices(%q) = %q, want %q", tc.targets, got, tc.want)
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
		{"/u/_search", `{"_source":{"excludes":["ssn"],"exclude":["email"]}}`, nil, []string{"email"}, false},
		{"/u/_search", `{"_source":{"include":"a","includes":["b","c"]}}`, []string{"b", "c"}, nil, false},
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
