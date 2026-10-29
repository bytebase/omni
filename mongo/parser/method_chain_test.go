package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/mongo/ast"
)

// TestMethodChainsStayLinear: a method chain inside an argument,
// `x.f().f()...`, becomes one HelperCall named by the whole dotted chain.
// The name was rebuilt at every step, copying the name so far, so a chain of
// n calls took time quadratic in n: 40,000 took 1.5 seconds. Each chain
// below takes milliseconds now; each parse runs against a deadline.
func TestMethodChainsStayLinear(t *testing.T) {
	const steps = 200000
	chain := strings.Repeat(".f()", steps-1) + ".g(1, 2)"
	for _, tc := range []struct {
		name  string
		input string
		root  string
		value func(stmt *ast.CollectionStatement) ast.Node
	}{
		{"in a document", "db.c.find({a: x" + chain + "})", "x", func(s *ast.CollectionStatement) ast.Node {
			return s.Args[0].(*ast.Document).Pairs[0].Value
		}},
		{"on a helper", "db.c.find({a: Long(1)" + chain + "})", "Long", func(s *ast.CollectionStatement) ast.Node {
			return s.Args[0].(*ast.Document).Pairs[0].Value
		}},
		{"as an argument", "db.c.find(x" + chain + ")", "x", func(s *ast.CollectionStatement) ast.Node {
			return s.Args[0]
		}},
	} {
		type result struct {
			nodes []ast.Node
			err   error
		}
		done := make(chan result, 1)
		go func() {
			nodes, err := Parse(tc.input)
			done <- result{nodes, err}
		}()
		select {
		case r := <-done:
			if r.err != nil {
				t.Fatalf("%s: %v", tc.name, r.err)
			}
			call, ok := tc.value(r.nodes[0].(*ast.CollectionStatement)).(*ast.HelperCall)
			if !ok {
				t.Fatalf("%s: value is not a HelperCall", tc.name)
			}
			wantName := tc.root + strings.Repeat(".f", steps-1) + ".g"
			if call.Name != wantName || len(call.Args) != 2 {
				t.Errorf("%s: name of %d bytes with %d args, want %d bytes with 2", tc.name, len(call.Name), len(call.Args), len(wantName))
			}
			if got := tc.input[call.Loc.Start:call.Loc.End]; !strings.HasPrefix(got, tc.root) || !strings.HasSuffix(got, ".g(1, 2)") {
				t.Errorf("%s: call spans %q...%q", tc.name, got[:10], got[len(got)-10:])
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse through %d chained calls took over 2s", tc.name, steps)
		}
	}
}
