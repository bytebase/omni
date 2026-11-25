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
