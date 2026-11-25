package catalog

import (
	"errors"
	"slices"
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

// Sort orders the graph with its soft edges, or without them when they make
// a cycle.
func (g *definitionGraph) Sort() ([]string, error) {
	if len(g.soft) > 0 {
		full := &definitionGraph{nodes: g.nodes, edges: make(map[string][]string)}
		for start, ends := range g.edges {
			full.edges[start] = append([]string(nil), ends...)
		}
		for _, edge := range g.soft {
			full.edges[edge[0]] = append(full.edges[edge[0]], edge[1])
		}
		if ordered, err := full.TopologicalSort(); err == nil {
			return ordered, nil
		}
	}
	return g.TopologicalSort()
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
