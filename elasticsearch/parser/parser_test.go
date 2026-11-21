package parser

import (
	"testing"
)

func TestRequestLineValidation(t *testing.T) {
	valid := []struct {
		name  string
		input string
		wantM string
		wantU string
		docs  int
	}{
		{name: "space separated", input: "GET /_search", wantM: "GET", wantU: "/_search"},
		{name: "tab separated", input: "GET\t/_search", wantM: "GET", wantU: "/_search"},
		{name: "HEAD root", input: "HEAD /", wantM: "HEAD", wantU: "/"},
		{name: "HEAD index", input: "head /my_index", wantM: "HEAD", wantU: "/my_index"},
		{name: "PATCH doc", input: "PATCH /my_index/_doc/1", wantM: "PATCH", wantU: "/my_index/_doc/1"},
		{name: "trailing spaces", input: "GET /_search   ", wantM: "GET", wantU: "/_search"},
		{name: "inline slash comment", input: "GET /_search // comment", wantM: "GET", wantU: "/_search"},
		{name: "inline tab slash comment", input: "GET\t/_search\t// comment", wantM: "GET", wantU: "/_search"},
		{name: "inline hash comment", input: "GET /_search # comment", wantM: "GET", wantU: "/_search"},
		{name: "inline block comment", input: "GET /_search /* comment */", wantM: "GET", wantU: "/_search"},
		{name: "multiple inline comments", input: "GET /_search /* c1 */ /* c2 */ // end", wantM: "GET", wantU: "/_search"},
		{name: "quoted space in query", input: `GET /_search?q="a b"`, wantM: "GET", wantU: `/_search?q="a b"`},
		{name: "quoted space with trailing comment", input: `GET /_search?q="a b" // comment`, wantM: "GET", wantU: `/_search?q="a b"`},
		{name: "CRLF single body", input: "POST /my_index/_search\r\n{\"query\": {}}", wantM: "POST", wantU: "/my_index/_search", docs: 1},
		{name: "CRLF multi-doc bulk", input: "POST /_bulk\r\n{\"index\":{\"_index\":\"a\"}}\r\n{\"f\":1}\r\n", wantM: "POST", wantU: "/_bulk", docs: 2},
	}

	for _, tc := range valid {
		t.Run("valid/"+tc.name, func(t *testing.T) {
			res, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned err: %v", tc.input, err)
			}
			if len(res.Errors) != 0 {
				t.Fatalf("Parse(%q) returned errors: %+v", tc.input, res.Errors)
			}
			if len(res.Requests) != 1 {
				t.Fatalf("Parse(%q) returned %d requests, want 1", tc.input, len(res.Requests))
			}
			lines := NewLines(tc.input)
			adj := lines.AdjustedParsedRequest(res.Requests[0], nil)
			req := lines.EditorRequest(adj)
			if req == nil {
				t.Fatalf("EditorRequest(%q) = nil", tc.input)
			}
			if req.Method != tc.wantM || req.URL != tc.wantU {
				t.Fatalf("EditorRequest(%q) = (%q, %q), want (%q, %q)", tc.input, req.Method, req.URL, tc.wantM, tc.wantU)
			}
			if len(req.Data) != tc.docs {
				t.Fatalf("EditorRequest(%q) data len = %d, want %d", tc.input, len(req.Data), tc.docs)
			}
		})
	}

	invalid := []struct {
		name  string
		input string
	}{
		{name: "missing space before slash", input: "GET/_search"},
		{name: "method prefix word", input: "GETfoo"},
		{name: "HEAD prefix word", input: "HEADING"},
		{name: "missing url", input: "GET"},
		{name: "missing url with spaces", input: "GET   "},
		{name: "missing url with newline", input: "GET\n"},
		{name: "missing method", input: "/_search"},
		{name: "unknown method", input: "OPTIONS /"},
		{name: "stray trailing token", input: "GET /_search zzz"},
		{name: "same-line JSON body", input: `GET /_search {"query":{}}`},
		{name: "duplicate method", input: "GET GET /_search"},
		{name: "comment-only slash", input: "GET //comment"},
		{name: "comment-only hash", input: "GET #comment"},
		{name: "comment-only block", input: "GET /* comment */"},
		{name: "unterminated block comment", input: "GET /_search /* unclosed"},
		{name: "unterminated quote in query", input: `GET /_search?q="unterminated`},
		{name: "bare CR without LF", input: "GET /_search\r{\"query\":{}}"},
	}

	for _, tc := range invalid {
		t.Run("invalid/"+tc.name, func(t *testing.T) {
			res, err := Parse(tc.input)
			if err == nil && res != nil && len(res.Errors) == 0 && len(res.Requests) == 1 {
				t.Fatalf("Parse(%q) unexpectedly accepted", tc.input)
			}
		})
	}
}

func TestAdjustedParsedRequestStopsAtHead(t *testing.T) {
	// When a request has EndOffset == -1 and no nextRequest is passed,
	// AdjustedParsedRequest scans ahead for a line starting with any Console
	// HTTP method, including HEAD.
	text := "GET /_search\n{\"a\": x}\nHEAD /my_index"
	lines := NewLines(text)
	adj := lines.AdjustedParsedRequest(ParsedRequest{StartOffset: 0, EndOffset: -1}, nil)
	if adj.StartLineNumber != 0 || adj.EndLineNumber != 1 {
		t.Fatalf("AdjustedParsedRequest = %+v, want {StartLineNumber:0 EndLineNumber:1}", adj)
	}
}
