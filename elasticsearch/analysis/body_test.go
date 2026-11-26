package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/elasticsearch/analysis"
)

// TestMultiSearchBody holds a multi-search to the analysis of each of its
// searches: a blocked feature in any search blocks the request, and the one
// _source filter keeps every field some search keeps.
func TestMultiSearchBody(t *testing.T) {
	for _, tc := range []struct {
		body              string
		blocked           []analysis.BlockedFeature
		disabled          bool
		includes, exclude []string
		predicates        []string
	}{
		{body: "{}\n{\"_source\":[\"name\"]}\n", includes: []string{"name"}},
		{body: "{}\n{\"query\":{\"term\":{\"ssn\":\"1\"}}}\n{}\n{\"aggs\":{\"a\":{\"terms\":{\"field\":\"ssn\"}}}}\n", blocked: []analysis.BlockedFeature{analysis.BlockedFeatureAggs}, predicates: []string{"ssn"}},
		{body: "{}\n{\"_source\":[\"name\"]}\n{}\n{\"_source\":[\"email\",\"name\"]}\n", includes: []string{"name", "email"}},
		{body: "{}\n{\"_source\":[\"name\"]}\n{}\n{}\n"},
		{body: "{}\n{\"_source\":false}\n{}\n{\"_source\":[\"name\"]}\n", includes: []string{"name"}},
		{body: "{}\n{\"_source\":false}\n", disabled: true},
		{body: "{}\n{\"_source\":{\"excludes\":[\"ssn\",\"age\"]}}\n{}\n{\"_source\":{\"excludes\":[\"ssn\"]}}\n", exclude: []string{"ssn"}},
		{body: "{}\n{\"_source\":{\"excludes\":[\"ssn\"]}}\n{}\n{}\n"},
		{body: "{}\n{\"_source\":[]}\n{}\n{\"_source\":[\"name\"]}\n"},
		{body: "{}\n{\"_source\":{\"includes\":[]}}\n{}\n{\"_source\":[\"name\"]}\n"},
	} {
		a := analysis.AnalyzeRequest("POST", "/users/_msearch", tc.body)
		if !slices.Equal(a.BlockedFeatures, tc.blocked) || a.SourceDisabled != tc.disabled ||
			!slices.Equal(a.SourceFields, tc.includes) || !slices.Equal(a.SourceExcludes, tc.exclude) ||
			!slices.Equal(a.PredicateFields, tc.predicates) {
			t.Errorf("%q: blocked %v, disabled %v, includes %q, excludes %q, predicates %q; want %v, %v, %q, %q, %q",
				tc.body, a.BlockedFeatures, a.SourceDisabled, a.SourceFields, a.SourceExcludes, a.PredicateFields,
				tc.blocked, tc.disabled, tc.includes, tc.exclude, tc.predicates)
		}
	}
}

// TestSearchOfAnIndexNamedMultiSearch holds a search of an index whose name
// holds _msearch to the analysis of one JSON body.
func TestSearchOfAnIndexNamedMultiSearch(t *testing.T) {
	a := analysis.AnalyzeRequest("POST", "/logs_msearch/_search", `{"aggs":{"a":{"terms":{"field":"ssn"}}}}`)
	if !slices.Equal(a.BlockedFeatures, []analysis.BlockedFeature{analysis.BlockedFeatureAggs}) {
		t.Errorf("blocked %v, want aggregations", a.BlockedFeatures)
	}
	if want := [][]string{{"logs_msearch"}}; !slices.EqualFunc(a.Targets, want, slices.Equal) {
		t.Errorf("Targets = %q, want %q", a.Targets, want)
	}
}
