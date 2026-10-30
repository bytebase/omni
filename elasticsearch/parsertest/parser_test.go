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
			input:       "GET _search\n{\n  \"query\": {\n    \"match_all\": {}\n  }\n}",
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

// TestParseNULIsAnOrdinaryCharacter: the parser took rune 0 for the end of
// the text, so a NUL byte ended the parse there and every request after it
// was dropped without an error. The upstream ends at the empty string that
// text.charAt returns past the end, and reads a NUL byte as any other
// character: whitespace between tokens, part of a URL, string or comment.
// Each want is what the upstream parser (Kibana's console parser.js) returns
// for the same input.
func TestParseNULIsAnOrdinaryCharacter(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		want        []parser.ParsedRequest
	}{
		{
			description: "in a URL",
			input:       "GET /a\x00b\nGET /c",
			want:        []parser.ParsedRequest{{StartOffset: 0, EndOffset: 8}, {StartOffset: 9, EndOffset: 15}},
		},
		{
			description: "at the end of a URL",
			input:       "GET /a\x00",
			want:        []parser.ParsedRequest{{StartOffset: 0, EndOffset: 7}},
		},
		{
			description: "before the first request",
			input:       "\x00GET /x",
			want:        []parser.ParsedRequest{{StartOffset: 1, EndOffset: 7}},
		},
		{
			description: "in a body string",
			input:       "GET /x\n{\"a\": \"b\x00c\"}\nGET /y",
			want:        []parser.ParsedRequest{{StartOffset: 0, EndOffset: 19}, {StartOffset: 20, EndOffset: 26}},
		},
		{
			description: "between body tokens",
			input:       "GET /x\n{\"a\":\x001}\nGET /y",
			want:        []parser.ParsedRequest{{StartOffset: 0, EndOffset: 15}, {StartOffset: 16, EndOffset: 22}},
		},
		{
			description: "in a line comment",
			input:       "# c\x00omment\nGET /z",
			want:        []parser.ParsedRequest{{StartOffset: 11, EndOffset: 17}},
		},
		{
			description: "in a block comment",
			input:       "GET /x\n/* a\x00b */\nGET /y",
			want:        []parser.ParsedRequest{{StartOffset: 0, EndOffset: 6}, {StartOffset: 17, EndOffset: 23}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := parser.Parse(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, result.Requests)
			require.Empty(t, result.Errors)

			stmts, err := es.SplitMultiSQL(tc.input)
			require.NoError(t, err)
			require.Len(t, stmts, len(tc.want))
			for i, r := range tc.want {
				require.Equal(t, tc.input[r.StartOffset:r.EndOffset], stmts[i].Text)
			}
		})
	}
}

// TestParseRecordsTextLeftAfterTheLastError: when a request fails and no
// request follows to resume at, text is left over. Parse returned "Syntax
// error" as its error then, and lost the requests and errors it had read, so
// SplitMultiSQL failed and Diagnose reported nothing. The upstream records the
// leftover as one more error, "Syntax error", and returns the rest. Each want
// holds the request offsets and error offsets the upstream parser (Kibana's
// console parser.js) returns for the same input.
func TestParseRecordsTextLeftAfterTheLastError(t *testing.T) {
	testCases := []struct {
		description  string
		input        string
		want         []parser.ParsedRequest
		errorOffsets []int
	}{
		{
			description:  "body of the only request fails",
			input:        "GET /_search\n{\"a\": x}\n",
			want:         []parser.ParsedRequest{{StartOffset: 0, EndOffset: -1}},
			errorOffsets: []int{20, 20},
		},
		{
			description:  "body of the last request fails",
			input:        "GET /a\n{\"a\": 1}\nGET /b\n{\"c\": x}",
			want:         []parser.ParsedRequest{{StartOffset: 0, EndOffset: 15}, {StartOffset: 16, EndOffset: -1}},
			errorOffsets: []int{30, 30},
		},
		{
			description:  "a line that is not a request after the last request",
			input:        "GET /x\n{}\n}\n",
			want:         []parser.ParsedRequest{{StartOffset: 0, EndOffset: 9}, {StartOffset: 10, EndOffset: -1}},
			errorOffsets: []int{11, 11},
		},
		{
			description:  "no request at all",
			input:        "db.users.find({})",
			want:         []parser.ParsedRequest{{StartOffset: 0, EndOffset: -1}},
			errorOffsets: []int{2, 2},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := parser.Parse(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, result.Requests)
			var offsets []int
			for _, e := range result.Errors {
				offsets = append(offsets, e.ByteOffset)
			}
			require.Equal(t, tc.errorOffsets, offsets)
			require.Equal(t, "Syntax error", result.Errors[len(result.Errors)-1].Message)
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

// TestIndentData pins what IndentData returns, which main's
// hjson.Unmarshal-based implementation produced for each input: comments
// gone, keys sorted, the last duplicate kept, numbers as float64, and the
// input itself back where the body is not a JSON object once HJSON is read.
func TestIndentData(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		want        string
	}{
		{
			description: "comments go and keys sort",
			input:       "{\"query\": {\"match\": {\"f\": \"v\"}}, // c\n \"size\": [true, null, \"s\", -1.5e+3]}",
			want:        "{\n  \"query\": {\n    \"match\": {\n      \"f\": \"v\"\n    }\n  },\n  \"size\": [\n    true,\n    null,\n    \"s\",\n    -1500\n  ]\n}",
		},
		{
			description: "block comment",
			input:       "{\"b\": 1, \"a\": 2} /* c */",
			want:        "{\n  \"a\": 2,\n  \"b\": 1\n}",
		},
		{
			description: "the last duplicate key wins",
			input:       "{\"a\": 1, \"a\": 2} // c",
			want:        "{\n  \"a\": 2\n}",
		},
		{
			description: "keys equal once made valid UTF-8 are one key",
			input:       "{\"\xff\": 1, \"\xfe\": 2} // c",
			want:        "{\n  \"�\": 2\n}",
		},
		{
			description: "a number HJSON accepts and JSON does not",
			input:       "{\"a\": 1.} // c",
			want:        "{\"a\": 1.} // c",
		},
		{
			description: "a number out of float64 range",
			input:       "{\"a\": 1e999} // c",
			want:        "{\"a\": 1e999} // c",
		},
		{
			description: "a number past float64 precision",
			input:       "{\"a\": 123456789012345678901234567890} // c",
			want:        "{\n  \"a\": 1.2345678901234568e+29\n}",
		},
		{
			description: "an HJSON quoteless string runs to the end of the line",
			input:       "{\"a\": quoteless // c\n}",
			want:        "{\n  \"a\": \"quoteless // c\"\n}",
		},
		{
			description: "HTML characters are escaped",
			input:       "{\"a\": \"<&>\"} // c",
			want:        "{\n  \"a\": \"\\u003c\\u0026\\u003e\"\n}",
		},
		{
			description: "empty object and array",
			input:       "{\"a\": {}, \"b\": []} // c",
			want:        "{\n  \"a\": {},\n  \"b\": []\n}",
		},
		{
			description: "a literal string HJSON cannot read",
			input:       "{\"a\": \"\"\"lit\"\"\"} // c",
			want:        "{\"a\": \"\"\"lit\"\"\"} // c",
		},
		{
			description: "null",
			input:       "null // c",
			want:        "null",
		},
		{
			description: "an array is not an object",
			input:       "[1] // c",
			want:        "[1] // c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.want, parser.IndentData(tc.input), "description: %s", tc.description)
		})
	}
}
