package parsertest

import (
	"testing"

	"github.com/bytebase/omni/elasticsearch"
	"github.com/stretchr/testify/require"
)

func TestGetQuerySpan(t *testing.T) {
	type testCase struct {
		Description string                `yaml:"description"`
		Statement   string                `yaml:"statement"`
		QueryType   elasticsearch.QueryType `yaml:"queryType"`
	}

	cases := loadYAML[testCase](t, "query_span.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			span, err := elasticsearch.GetQuerySpan(tc.Statement)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Equalf(t, tc.QueryType, span.Type,
				"description: %s, statement: %s", tc.Description, tc.Statement)
		})
	}
}

func TestGetQuerySpanPredicatePaths(t *testing.T) {
	tests := []struct {
		description    string
		statement      string
		predicatePaths map[string]bool
	}{
		{
			description: "search with match predicate",
			statement:   `GET /users/_search` + "\n" + `{"query":{"match":{"email":"alice"}}}`,
			predicatePaths: map[string]bool{
				"email": true,
			},
		},
		{
			description: "search with nested bool query",
			statement: `GET /users/_search` + "\n" +
				`{"query":{"bool":{"must":[{"term":{"status":"active"}},{"range":{"age":{"gte":18}}}]}}}`,
			predicatePaths: map[string]bool{
				"status": true,
				"age":    true,
			},
		},
		{
			description:    "get doc returns no predicates",
			statement:      "GET /users/_doc/123",
			predicatePaths: map[string]bool{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			span, err := elasticsearch.GetQuerySpan(tc.statement)
			require.NoError(t, err, tc.description)
			require.Equal(t, len(tc.predicatePaths), len(span.PredicatePaths), tc.description)
			for path := range tc.predicatePaths {
				_, ok := span.PredicatePaths[path]
				require.Truef(t, ok, "%s: missing path %q", tc.description, path)
			}
		})
	}
}

func TestGetQuerySpanDotPathAST(t *testing.T) {
	tests := []struct {
		dotPath  string
		segments []string
	}{
		{"email", []string{"email"}},
		{"user.name", []string{"user", "name"}},
		{"contact.address.city", []string{"contact", "address", "city"}},
	}

	for _, tc := range tests {
		t.Run(tc.dotPath, func(t *testing.T) {
			ast := elasticsearch.DotPathToPathAST(tc.dotPath)
			require.NotNil(t, ast)
			require.NotNil(t, ast.Root)

			// Walk the chain and collect segment names.
			var got []string
			cur := ast.Root
			for cur != nil {
				got = append(got, cur.Name)
				cur = cur.Next
			}
			require.Equal(t, tc.segments, got)
		})
	}
}
