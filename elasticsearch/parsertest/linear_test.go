package parsertest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	es "github.com/bytebase/omni/elasticsearch"
)

// entryPoint is one public function of the package that parses a script.
type entryPoint struct {
	name string
	call func(string) (int, error)
}

var (
	parseREST          = entryPoint{"ParseElasticsearchREST", func(s string) (int, error) { r, err := es.ParseElasticsearchREST(s); return countRequests(r), err }}
	splitMultiSQL      = entryPoint{"SplitMultiSQL", func(s string) (int, error) { r, err := es.SplitMultiSQL(s); return len(r), err }}
	diagnose           = entryPoint{"Diagnose", func(s string) (int, error) { r, err := es.Diagnose(s); return len(r), err }}
	getStatementRanges = entryPoint{"GetStatementRanges", func(s string) (int, error) { r, err := es.GetStatementRanges(s); return len(r), err }}
	getQuerySpan       = entryPoint{"GetQuerySpan", func(s string) (int, error) { _, err := es.GetQuerySpan(s); return 0, err }}
	allEntryPoints     = []entryPoint{parseREST, splitMultiSQL, diagnose, getStatementRanges, getQuerySpan}
)

func countRequests(r *es.ParseResult) int {
	if r == nil {
		return 0
	}
	return len(r.Requests)
}

// runWithin calls each entry point on text in a goroutine and fails the test
// when they take longer than limit together, so that a regression fails the
// test instead of hanging it. It returns what each entry point counted: the
// requests, statements, diagnostics or ranges it returned.
func runWithin(t *testing.T, name, text string, limit time.Duration, entries ...entryPoint) map[string]int {
	t.Helper()
	type result struct {
		counts map[string]int
		err    error
	}
	done := make(chan result, 1)
	go func() {
		counts := map[string]int{}
		for _, e := range entries {
			n, err := e.call(text)
			if err != nil {
				done <- result{err: fmt.Errorf("%s: %w", e.name, err)}
				return
			}
			counts[e.name] = n
		}
		done <- result{counts: counts}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("%s: %v", name, r.err)
		}
		return r.counts
	case <-time.After(limit):
		t.Fatalf("%s: the entry points took over %v on %d bytes", name, limit, len(text))
		return nil
	}
}

// TestRequestBodiesStayLinear guards the strings built while a request is
// read: the parser built each string, number and URL, and
// SplitDataIntoJSONObjects each document of the body, by appending one
// character at a time to a string, which copied it whole every time. A
// 100,000-character string took the parser 1.7s and a 117KB body
// SplitDataIntoJSONObjects 2.2s, and every entry point does both. Each request
// below, 450KB to 1MB, takes milliseconds when linear and minutes otherwise.
func TestRequestBodiesStayLinear(t *testing.T) {
	keys := func(n int) string {
		var b strings.Builder
		b.WriteString("{")
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteString(",\n")
			}
			fmt.Fprintf(&b, "\"k%d\": %d", i, i)
		}
		b.WriteString("}")
		return b.String()
	}
	for _, tc := range []struct {
		name string
		text string
	}{
		{"long string", "GET /_search\n{\"a\": \"" + strings.Repeat("x", 1000000) + "\"}\n"},
		{"long literal string", "GET /_search\n{\"a\": \"\"\"" + strings.Repeat("x", 1000000) + "\"\"\"}\n"},
		{"long number", "GET /_search\n{\"a\": 0." + strings.Repeat("1", 1000000) + "}\n"},
		{"long URL", "GET /" + strings.Repeat("x", 1000000) + "\n"},
		{"many keys", "GET /_search\n" + keys(40000) + "\n"},
		{"long array", "GET /_search\n{\"a\": [" + strings.Repeat("1, ", 150000) + "1]}\n"},
		{"many comments", "GET /_search\n{" + strings.Repeat("/* c */ // d\n", 80000) + "\"a\": 1}\n"},
	} {
		counts := runWithin(t, tc.name, tc.text, 10*time.Second, allEntryPoints...)
		if counts["SplitMultiSQL"] != 1 || counts["Diagnose"] != 0 {
			t.Errorf("%s: %v, want one statement and no diagnostics", tc.name, counts)
		}
	}
}

// TestCommentedBodiesStayLinear guards IndentData, which rewrites a
// document that holds comments as indented JSON. hjson.Unmarshal, which it
// called, converts the HJSON tree to JSON by marshaling each object once for
// every object it is nested in, so a 1MB string under 1,000 objects took
// 18s, and 2,000 objects 37s. ParseElasticsearchREST and GetQuerySpan
// rewrite this body in well under a second now.
func TestCommentedBodiesStayLinear(t *testing.T) {
	const depth = 1000
	text := "GET /_search\n" + strings.Repeat("{\"a\": // c\n", depth) +
		"\"" + strings.Repeat("x", 1000000) + "\"" + strings.Repeat("}", depth) + "\n"
	counts := runWithin(t, "body with comments", text, 10*time.Second, parseREST, getQuerySpan)
	if counts["ParseElasticsearchREST"] != 1 {
		t.Errorf("%v, want one request", counts)
	}
}

// TestDeepCommentedBodiesStayLinear: the indented JSON IndentData writes
// grows with the square of the depth, its indentation alone 200MB for a
// document nested 10,000 deep, and SplitMultiSQL, Diagnose and
// GetStatementRanges, which do not read it, no longer have it written. They
// took 1.5s each at 5,000 levels and four times that at 10,000.
func TestDeepCommentedBodiesStayLinear(t *testing.T) {
	const depth = 10000 // the parser's limit
	text := "GET /_search\n" + strings.Repeat("{\"a\": // c\n", depth) + "1" + strings.Repeat("}", depth) + "\n"
	counts := runWithin(t, "deep body with comments", text, 5*time.Second, splitMultiSQL, diagnose, getStatementRanges)
	if counts["SplitMultiSQL"] != 1 || counts["Diagnose"] != 0 {
		t.Errorf("%v, want one statement and no diagnostics", counts)
	}
}
