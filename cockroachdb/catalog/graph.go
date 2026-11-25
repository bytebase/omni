package catalog

import (
	"errors"
	"slices"
	"strings"
)

// definitionGraph orders object IDs so that each object follows the objects it
// depends on. Of the objects ready at each step the least ID goes first.
type definitionGraph struct {
	nodes map[string]bool
	edges map[string][]string
	// soft holds edges that may be false, such as a call that more than one
	// overload could answer; Sort drops them when they make a cycle.
	soft [][2]string
}

func newDefinitionGraph() *definitionGraph {
	return &definitionGraph{nodes: make(map[string]bool), edges: make(map[string][]string)}
}

func (g *definitionGraph) AddNode(id string) {
	g.nodes[id] = true
}

// AddEdge records that end depends on start, adding either as a node.
func (g *definitionGraph) AddEdge(start, end string) {
	if start == end {
		return
	}
	g.nodes[start] = true
	g.nodes[end] = true
	g.edges[start] = append(g.edges[start], end)
}

// AddSoftEdge records that end may depend on start.
func (g *definitionGraph) AddSoftEdge(start, end string) {
	if start == end {
		return
	}
	g.nodes[start] = true
	g.nodes[end] = true
	g.soft = append(g.soft, [2]string{start, end})
}

// Sort orders the graph with as many of its soft edges as fit: each, in a
// stable order, is kept unless it would close a cycle, so a false edge drops
// without taking a real one with it.
func (g *definitionGraph) Sort() ([]string, error) {
	if len(g.soft) == 0 {
		return g.TopologicalSort()
	}
	full := &definitionGraph{nodes: g.nodes, edges: make(map[string][]string)}
	for start, ends := range g.edges {
		full.edges[start] = append([]string(nil), ends...)
	}
	soft := slices.Clone(g.soft)
	slices.SortFunc(soft, func(a, b [2]string) int {
		if c := strings.Compare(a[0], b[0]); c != 0 {
			return c
		}
		return strings.Compare(a[1], b[1])
	})
	for _, edge := range soft {
		if !full.reaches(edge[1], edge[0]) {
			full.edges[edge[0]] = append(full.edges[edge[0]], edge[1])
		}
	}
	return full.TopologicalSort()
}

// reaches reports whether a path of edges leads from start to end.
func (g *definitionGraph) reaches(start, end string) bool {
	seen := map[string]bool{start: true}
	stack := []string{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == end {
			return true
		}
		for _, next := range g.edges[n] {
			if !seen[next] {
				seen[next] = true
				stack = append(stack, next)
			}
		}
	}
	return false
}

// TopologicalSort returns the IDs in dependency order, or an error when the
// graph has a cycle.
func (g *definitionGraph) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int)
	for _, ends := range g.edges {
		for _, end := range ends {
			inDegree[end]++
		}
	}
	var ready []string
	for id := range g.nodes {
		if inDegree[id] == 0 {
			ready = append(ready, id)
		}
	}
	slices.Sort(ready)
	result := make([]string, 0, len(g.nodes))
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		result = append(result, id)
		for _, end := range g.edges[id] {
			inDegree[end]--
			if inDegree[end] == 0 {
				i, _ := slices.BinarySearch(ready, end)
				ready = slices.Insert(ready, i, end)
			}
		}
	}
	if len(result) != len(g.nodes) {
		return nil, errors.New("graph has cycle")
	}
	return result, nil
}
