package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/bytebase/omni/mongo/parser"
)

// mustParse parses input expecting exactly one AST node.
func mustParse(t *testing.T, input string) ast.Node {
	t.Helper()
	nodes, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	return nodes[0]
}

// mustParseN parses input expecting exactly n AST nodes.
func mustParseN(t *testing.T, input string, n int) []ast.Node {
	t.Helper()
	nodes, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(nodes) != n {
		t.Fatalf("expected %d nodes, got %d", n, len(nodes))
	}
	return nodes
}

// mustFail parses input and expects a parse error.
func mustFail(t *testing.T, input string) {
	t.Helper()
	_, err := parser.Parse(input)
	if err == nil {
		t.Fatalf("expected parse error for input: %s", input)
	}
}

// assertLoc checks that a node's location matches the expected start and end offsets.
func assertLoc(t *testing.T, node ast.Node, start, end int) {
	t.Helper()
	loc := node.GetLoc()
	if loc.Start != start || loc.End != end {
		t.Errorf("expected loc [%d, %d), got [%d, %d)", start, end, loc.Start, loc.End)
	}
}
