package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest/cputime"
)

// TestNestingLimit: every value a request body nests recursively stops at
// maxNestingDepth with a syntax error instead of overflowing the Go stack,
// which killed the process: a body nested about 3 million levels deep was
// enough. The same values nested 1,000 deep still parse. Each deep parse runs
// against a deadline, so a parse that retried at the limit fails the test
// instead of hanging it.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	// A request follows the deep one: the parser records a body's error and
	// resumes at the next request, and without one it fails the whole text
	// with a bare "Syntax error", as Kibana's console parser does.
	request := func(body string) string { return "GET /_search\n" + body + "\nGET /_search\n" }
	for _, tc := range []struct {
		name string
		text func(n int) string
	}{
		{"objects", func(n int) string { return request(nest(n, `{"a":`, "1", "}")) }},
		{"arrays", func(n int) string { return request(`{"a":` + nest(n, "[", "1", "]") + "}") }},
		{"objects in arrays", func(n int) string { return request(nest(n, `{"a":[`, "1", "]}")) }},
		{"second document", func(n int) string { return request("{}\n" + nest(n, `{"a":`, "1", "}")) }},
		{"unclosed objects", func(n int) string { return request(strings.Repeat(`{"a":`, n)) }},
		{"unclosed arrays", func(n int) string { return request(`{"a":` + strings.Repeat("[", n)) }},
	} {
		if !strings.HasPrefix(tc.name, "unclosed") {
			res, err := Parse(tc.text(1000))
			if err != nil || len(res.Errors) != 0 {
				t.Errorf("%s: 1,000 levels: err = %v, syntax errors = %v", tc.name, err, res)
			}
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			text := tc.text(n)
			var res *ParseResult
			var err error
			cputime.Within(t, 10*time.Second, fmt.Sprintf("%s: %d levels", tc.name, n), func() {
				res, err = Parse(text)
			})
			if err != nil || len(res.Errors) == 0 || res.Errors[0].Message != errNestedTooDeeply.Error() {
				t.Errorf("%s: %d levels: err = %v, result = %+v, want the nesting error", tc.name, n, err, res)
			}
		}
	}
}
