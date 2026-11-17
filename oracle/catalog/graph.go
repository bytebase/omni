package catalog

import (
	"fmt"
	"slices"
)

// Graph is the small dependency graph definition and migration ordering use.
type Graph struct {
	nodes map[string]bool
	edges [][2]string
}

func newGraph() *Graph {
	return &Graph{nodes: make(map[string]bool)}
}

func (g *Graph) AddNode(id string) {
	g.nodes[id] = true
}

func (g *Graph) AddEdge(start, end string) {
	g.nodes[start] = true
	g.nodes[end] = true
	g.edges = append(g.edges, [2]string{start, end})
}

func (g *Graph) TopologicalSort() ([]string, error) {
	var result []string
	inDegree := make(map[string]int)
	outEdges := make(map[string][][2]string)
	for _, edge := range g.edges {
		inDegree[edge[1]]++
		outEdges[edge[0]] = append(outEdges[edge[0]], edge)
	}
	var queue []string
	for id := range g.nodes {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}
	slices.Sort(queue)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, edge := range outEdges[node] {
			inDegree[edge[1]]--
			if inDegree[edge[1]] == 0 {
				queue = append(queue, edge[1])
			}
		}
		result = append(result, node)
	}
	if len(result) != len(g.nodes) {
		return nil, fmt.Errorf("graph has cycle")
	}
	return result, nil
}
