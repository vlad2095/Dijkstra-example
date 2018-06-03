package main

import (
	"bytes"
	"fmt"
)

// A Vertex represents a node in a directed multi graph.
type Vertex struct {
	id  string
	In  map[string]*Arc // map of ingoing arcs
	Out map[string]*Arc // map of outgoing arcs
}

func NewVertex(id string) *Vertex {
	v := &Vertex{id: id,
		In:  make(map[string]*Arc),
		Out: make(map[string]*Arc),
	}
	return v
}

func (v *Vertex) AddIngoingArc(arc *Arc) error {
	if _, ok := v.In[arc.from]; !ok {
		v.In[arc.from] = arc
	} else {
		return fmt.Errorf("arc <%s --> %s>  already exist in vertex <%s> ", arc.from, arc.to, v.id)
	}
	return nil
}

func (v *Vertex) AddOutgoingArc(arc *Arc) error {
	if _, ok := v.Out[arc.to]; !ok {
		v.Out[arc.to] = arc
	} else {
		return fmt.Errorf("arc <%s --> %s>  already exist in vertex <%s> ", arc.from, arc.to, v.id)
	}
	return nil
}

func (v *Vertex) GetIngoingArc(from string) (*Arc, bool) {
	arc, ok := v.In[from]
	return arc, ok
}

func (v *Vertex) GetOutgoingArc(to string) (*Arc, bool) {
	arc, ok := v.Out[to]
	return arc, ok
}

// An Arc represents a parallel edges in the directed multi graph.
type Arc struct {
	from  string
	to    string
	edges map[string]*Edge
}

func NewArc(from, to string) *Arc {
	a := &Arc{from: from,
		to:    to,
		edges: make(map[string]*Edge)}
	return a
}

func (a *Arc) AddEdge(edge *Edge) error {
	if _, ok := a.edges[edge.id]; !ok {
		a.edges[edge.id] = edge
	} else {
		return fmt.Errorf("edge <%s> already exist in arc <%s --> %s> ", edge.id, a.from, a.to)
	}
	return nil
}

func (a *Arc) GetEdges() map[string]*Edge {
	return a.edges
}

// An Edge represents a single connection between two vertices
// an Arc contains from one or more edges
type Edge struct {
	id   string
	data interface{}
}

func NewEdge(id string, data interface{}) *Edge {
	return &Edge{id: id, data: data}
}

// Graph describes the methods of graph operations.
// It assumes that the identifier of a Vertex is unique.
type Graph struct {
	Vertices map[string]*Vertex //Map of vertices by their id
}

func NewGraph() *Graph {
	g := &Graph{Vertices: make(map[string]*Vertex)}
	return g
}

// AddVertex adds new vertex, return an error if vertex with same id already exist
func (g *Graph) AddVertex(v *Vertex) error {
	if _, ok := g.GetVertex(v.id); !ok {
		g.Vertices[v.id] = v
	} else {
		return fmt.Errorf("vertex <%s> is already exist", v.id)
	}
	return nil
}

// GetVertex return a vertex and true if its present in graph
// nil and false otherwise
func (g *Graph) GetVertex(id string) (*Vertex, bool) {
	v, ok := g.Vertices[id]
	return v, ok
}

// returns all vertices
func (g *Graph) GetVertices() map[string]*Vertex {
	return g.Vertices
}

// converts all graph data into string
func (g *Graph) String() string {
	buf := new(bytes.Buffer)
	// iterate over all vertices of graph
	for _, node := range g.Vertices {
		fmt.Fprintf(buf, "%s\n", node.id)
		// iterate over all arcs out from vertex
		for _, arc := range node.Out {
			fmt.Fprintf(buf, "%s --- > %s\n", arc.from, arc.to)
			// iterate over all edges of arc
			for _, edge := range arc.edges {
				fmt.Fprintf(buf, "%s ---- %v\n", edge.id, edge.data)
			}
		}
	}
	return buf.String()
}
