package catalog

import (
	"errors"
	"slices"
)

// graph is the topological sort Bytebase's schema package uses, ported so that
// a dump's object order does not change with the move into omni.
type graph struct {
	nodes map[string]bool
	edges []edge
}

type edge struct {
	start string
	end   string
}

func newGraph() *graph {
	return &graph{nodes: make(map[string]bool)}
}

func (g *graph) addNode(id string) {
	g.nodes[id] = true
}

func (g *graph) addEdge(start, end string) {
	g.addNode(start)
	g.addNode(end)
	g.edges = append(g.edges, edge{start: start, end: end})
}

func (g *graph) topologicalSort() ([]string, error) {
	inDegree := make(map[string]int)
	outEdges := make(map[string][]edge)
	for _, e := range g.edges {
		inDegree[e.end]++
		outEdges[e.start] = append(outEdges[e.start], e)
	}
	var queue []string
	for id := range g.nodes {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}
	slices.Sort(queue)
	var result []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, e := range outEdges[node] {
			inDegree[e.end]--
			if inDegree[e.end] == 0 {
				queue = append(queue, e.end)
			}
		}
		result = append(result, node)
	}
	if len(result) != len(g.nodes) {
		return nil, errors.New("graph has cycle")
	}
	return result, nil
}

func objectID(schema, name string) string {
	if schema == "" {
		return name
	}
	return schema + "." + name
}
