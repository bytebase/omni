package parsertest

import (
	"testing"

	"github.com/stretchr/testify/require"

	es "github.com/bytebase/omni/elasticsearch"
	"github.com/bytebase/omni/elasticsearch/parser"
)

func TestParseElasticsearchREST(t *testing.T) {
	type testCase struct {
		Description string          `yaml:"description,omitempty"`
		Statement   string          `yaml:"statement,omitempty"`
		Result      *es.ParseResult `yaml:"result,omitempty"`
	}

	var (
		filename = "parse-elasticsearch-rest.yaml"
		record   = false
	)

	cases := loadYAML[testCase](t, filename)

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			got, err := es.ParseElasticsearchREST(tc.Statement)
			require.NoErrorf(t, err, "description: %s", tc.Description)
			if record {
				cases[i].Result = got
			} else {
				require.Equalf(t, tc.Result, got, "description: %s", tc.Description)
			}
		})
	}

	if record {
		writeYAML(t, filename, cases)
	}
}

func TestParse(t *testing.T) {
	// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/packages/shared/kbn-monaco/src/languages/console/parser.test.ts.
	testCases := []struct {
		description string
		input       string
		got         []parser.ParsedRequest
	}{
		{
			description: "returns parsedRequests if the input is correct",
			input:       "GET _search",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   11,
				},
			},
		},
		{
			description: "parses several requests",
			input:       "GET _search\nPOST _test_index",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   11,
				},
				{
					StartOffset: 12,
					EndOffset:   28,
				},
			},
		},
		{
			description: "parses a request with a request body",
			input: "GET _search\n{\n  \"query\": {\n    \"match_all\": {}\n  }\n}",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   52,
				},
			},
		},
		{
			description: "allows upper case methods",
			input:       "GET _search\nPOST _search\nPATCH _search\nPUT _search\nHEAD _search",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   11,
				},
				{
					StartOffset: 12,
					EndOffset:   24,
				},
				{
					StartOffset: 25,
					EndOffset:   38,
				},
				{
					StartOffset: 39,
					EndOffset:   50,
				},
				{
					StartOffset: 51,
					EndOffset:   63,
				},
			},
		},
		{
			description: "allows lower case methods",
			input:       "get _search\npost _search\npatch _search\nput _search\nhead _search",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   11,
				},
				{
					StartOffset: 12,
					EndOffset:   24,
				},
				{
					StartOffset: 25,
					EndOffset:   38,
				},
				{
					StartOffset: 39,
					EndOffset:   50,
				},
				{
					StartOffset: 51,
					EndOffset:   63,
				},
			},
		},
		{
			description: "allows mixed case methods",
			input:       "GeT _search\npOSt _search\nPaTch _search\nPut _search\nheAD _search",
			got: []parser.ParsedRequest{
				{
					StartOffset: 0,
					EndOffset:   11,
				},
				{
					StartOffset: 12,
					EndOffset:   24,
				},
				{
					StartOffset: 25,
					EndOffset:   38,
				},
				{
					StartOffset: 39,
					EndOffset:   50,
				},
				{
					StartOffset: 51,
					EndOffset:   63,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := parser.Parse(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.got, result.Requests, "description: %s", tc.description)
		})
	}
}

func TestGetEditorRequest(t *testing.T) {
	testCases := []struct {
		description string
		content     string
		adjusted    parser.AdjustedParsedRequest
		want        *parser.EditorRequest
	}{
		{
			description: "cleans up any text following the url",
			content:     "GET _search // inline comment",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   0,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
			},
		},
		{
			description: "doesn't incorrectly removes parts of url params that include whitespaces",
			content:     `GET _search?query="test test"`,
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   0,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    `_search?query="test test"`,
			},
		},
		{
			description: "correctly includes the request body",
			content:     "GET _search\n{\n  \"query\": {}\n}",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   3,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
				Data: []string{
					"{\n  \"query\": {}\n}",
				},
			},
		},
		{
			description: "correctly handles nested braces",
			content:     "GET _search\n{\n  \"query\": \"{a} {b}\"\n}\n{\n  \"query\": {}\n}",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   6,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
				Data: []string{
					"{\n  \"query\": \"{a} {b}\"\n}",
					"{\n  \"query\": {}\n}",
				},
			},
		},
		{
			description: "works for several request bodies",
			content:     "GET _search\n{\n  \"query\": {}\n}\n{\n  \"query\": {}\n}",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   6,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
				Data: []string{
					"{\n  \"query\": {}\n}",
					"{\n  \"query\": {}\n}",
				},
			},
		},
		{
			description: "splits several json objects",
			content:     "GET _search\n{\"query\":\"test\"}\n{\n  \"query\": \"test\"\n}\n{\"query\":\"test\"}",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   5,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
				Data: []string{
					`{"query":"test"}`,
					"{\n  \"query\": \"test\"\n}",
					`{"query":"test"}`,
				},
			},
		},
		{
			description: "works for invalid json objects",
			content:     "GET _search\n{\"query\":\"test\"}\n{\n  \"query\":\n{",
			adjusted: parser.AdjustedParsedRequest{
				StartLineNumber: 0,
				EndLineNumber:   4,
			},
			want: &parser.EditorRequest{
				Method: "GET",
				URL:    "_search",
				Data: []string{
					`{"query":"test"}`,
					"{\n  \"query\":\n{",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := parser.GetEditorRequest(tc.content, tc.adjusted)
			require.Equal(t, tc.want, got, "description: %s", tc.description)
		})
	}
}

func TestContainsComments(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		want        bool
	}{
		{
			description: "should return false for JSON with // and /* inside strings",
			input: `{
      "docs": [
        {
          "_source": {
            "trace": {
              "name": "GET /actuator/health/**"
            },
            "transaction": {
              "outcome": "success"
            }
          }
        },
        {
          "_source": {
            "vulnerability": {
              "reference": [
                "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-15778"
              ]
            }
          }
        }
      ]
    }`,
			want: false,
		},
		{
			description: "should return true for text with actual line comment",
			input: `{
      // This is a comment
      "query": { "match_all": {} }
    }`,
			want: true,
		},
		{
			description: "should return true for text with actual block comment",
			input: `{
      /* Bulk insert */
      "index": { "_index": "test" },
      "field1": "value1"
    }`,
			want: true,
		},
		{
			description: "should return false for text without any comments",
			input: `{
      "field": "value"
    }`,
			want: false,
		},
		{
			description: "should return false for empty string",
			input:       ``,
			want:        false,
		},
		{
			description: "should correctly handle escaped quotes within strings",
			input: `{
      "field": "value with \"escaped quotes\""
    }`,
			want: false,
		},
		{
			description: "should return true if comment is outside of strings",
			input: `{
      "field": "value" // comment here
    }`,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := parser.ContainsComments(tc.input)
			require.Equal(t, tc.want, got, "description: %s", tc.description)
		})
	}
}
