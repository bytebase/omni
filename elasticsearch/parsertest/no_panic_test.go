package parsertest

import (
	"fmt"
	"testing"

	es "github.com/bytebase/omni/elasticsearch"
)

// noPanicScripts are multi-request scripts whose prefixes and one-byte
// corruptions TestEntryPointsNeverPanic feeds to every entry point.
var noPanicScripts = []string{
	"GET /_search\n{\"a\": x}\nGET /x\n",
	"GET /idx/_search // trailing\n{\n  // comment\n  \"query\": {\"match\": {\"f\": \"v\"}} /* block */\n}\n\n" +
		"POST /idx/_doc\n{\"a\": \"\"\"lit\neral\"\"\"}\n# hash\nPUT /x\n{}\n{}\nDELETE /y",
	"POST _bulk\n{\"index\":{\"_index\":\"a\"}}\n{\"f\":1}\n{\"index\":{}}\n{\"f\":\"é😀\\u00e9\"}\n" +
		"GET /a/_search?q=\"x y\"\n{\"size\": -1.5e+3, \"b\": [true, false, null]}\n",
	"get /a\r\n{\"a\":1}\r\nhead /b\r\npatch /c\n{\"x\": [1, {\"y\": \"}\"}]}\n",
	"GET /a\n{} junk\nGET /b\n3333\nPUT /c\n{\"a\": \"\\\"\"}\n",
}

// TestEntryPointsNeverPanic: no input panics an entry point. A request whose
// body failed to parse kept EndOffset -1, and SplitMultiSQL sliced the text
// with it: "GET /_search\n{\"a\": x}\nGET /x\n", a body truncated at the end of
// the input, and "GET \x00" all panicked with "slice bounds out of range
// [:-1]". Every prefix of each script and every one-byte substitution of it
// runs through each entry point, and the statements SplitMultiSQL returns must
// cover disjoint ranges, in order, whose text they hold.
func TestEntryPointsNeverPanic(t *testing.T) {
	substitutes := []byte("{}[]\":,\n #/x\x00\xff")
	seen := map[string]bool{}
	var inputs []string
	add := func(s string) {
		if !seen[s] {
			seen[s] = true
			inputs = append(inputs, s)
		}
	}
	for _, script := range noPanicScripts {
		for i := 0; i <= len(script); i++ {
			add(script[:i])
		}
		for i := 0; i < len(script); i++ {
			for _, c := range substitutes {
				b := []byte(script)
				b[i] = c
				add(string(b))
			}
		}
	}
	failures := 0
	for _, in := range inputs {
		if err := checkEntryPoints(in); err != nil {
			t.Errorf("%q: %v", in, err)
			if failures++; failures == 10 {
				t.Fatal("too many failures")
			}
		}
	}
}

func checkEntryPoints(in string) (err error) {
	entry := ""
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", entry, r)
		}
	}()
	entry = "ParseElasticsearchREST"
	_, _ = es.ParseElasticsearchREST(in)
	entry = "Diagnose"
	_, _ = es.Diagnose(in)
	entry = "GetStatementRanges"
	_, _ = es.GetStatementRanges(in)
	entry = "GetQuerySpan"
	_, _ = es.GetQuerySpan(in)
	entry = "SplitMultiSQL"
	stmts, splitErr := es.SplitMultiSQL(in)
	if splitErr != nil {
		return nil
	}
	if len(stmts) == 1 && stmts[0].Range == (es.Range{}) && stmts[0].Text == in {
		// No request parsed: the whole text is one statement.
		return nil
	}
	prevEnd := 0
	for i, s := range stmts {
		if s.Range.Start < prevEnd || s.Range.End < s.Range.Start || s.Range.End > len(in) {
			return fmt.Errorf("statement %d has range %+v after offset %d in %d bytes", i, s.Range, prevEnd, len(in))
		}
		if s.Text != in[s.Range.Start:s.Range.End] {
			return fmt.Errorf("statement %d has text %q, want %q", i, s.Text, in[s.Range.Start:s.Range.End])
		}
		prevEnd = s.Range.End
	}
	return nil
}
