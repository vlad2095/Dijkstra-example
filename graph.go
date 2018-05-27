package main

import (
	"fmt"
	"time"
	"bytes"
)

// ID is unique identifier.
type ID interface {
	// String returns the string ID.
	String() string
}

type StringID string

func (s StringID) String() string {
	return string(s)
}

// Node is vertex. The ID must be unique within the graph.
type Node interface {
	// ID returns the ID.
	ID() ID
	String() string
}

type node struct {
	id string
}

func NewNode(id string) Node {
	return &node{
		id: id,
	}
}

func (n *node) ID() ID         { return StringID(n.id) }
func (n *node) String() string { return n.id }


type Train struct {
	ID                  string  `xml:"TrainId,attr"`
	DepartureStationId  string  `xml:"DepartureStationId,attr"`
	ArrivalStationId    string  `xml:"ArrivalStationId,attr"`
	DepartureTimeString string  `xml:"DepartureTimeString,attr"`
	ArrivalTimeString   string  `xml:"ArrivalTimeString,attr"`
	Price               float64 `xml:"Price,attr"`
	Duration            time.Duration
}

func (t Train) String() string {
	f := `№%s Dep.Station: %s Dep.Time: %s Arr.Station: %s Arr.Time: %s Price: %.2f$" Duration: %s`
	return fmt.Sprintf(f, t.ID,
		t.DepartureStationId, t.DepartureTimeString,
		t.ArrivalStationId, t.ArrivalTimeString,
		t.Price, t.Duration.String())
}

const timeFormat = "15:04:05"

func (t *Train) FindDuration() {
	depTime, _ := time.Parse(timeFormat, t.DepartureTimeString)
	arrTime, _ := time.Parse(timeFormat, t.ArrivalTimeString)
	if !arrTime.After(depTime) {
		arrTime = arrTime.Add(time.Hour * 24)
	}
	t.Duration = arrTime.Sub(depTime)
}

func (t *Train) GetWeight() float64 {
	return t.Price * float64(t.Duration)
}

// Edge connects between two Nodes.
type Edge interface {
	Source() Node
	Target() Node
	Train() Train
	Weight() float64
}

// edge is an Edge from Source to Target.
type edge struct {
	src   Node
	tgt   Node
	wgt   float64
	train Train
}

func (e *edge) Source() Node    { return e.src }
func (e *edge) Target() Node    { return e.tgt }
func (e *edge) Weight() float64 { return e.wgt }
func (e *edge) Train() Train    { return e.train }

func NewEdge(src, tgt Node, train Train) Edge {
	return &edge{
		src:   src,
		tgt:   tgt,
		wgt:   train.GetWeight(),
		train: train,
	}
}

// Graph describes the methods of graph operations.
// It assumes that the identifier of a Node is unique.
type Graph interface {
	ExistNode(id ID) bool

	GetNode(id ID) Node

	// GetNodes returns a map from node ID to empty struct value.
	GetNodes() map[ID]Node

	// AddNode adds a node to a graph, and returns false
	AddNode(nd Node) bool

	// GetEdge returns the edge from id1 to id2 (nil if not exist)
	GetEdge(id1, id2 ID) Edge

	// ReplaceEdge replaces an edge from id1 to id2.
	ReplaceEdge(id1, id2 ID, edge Edge) error

	// GetWeight returns the weight from id1 to id2.
	//GetWeight(id1, id2 ID) (float64, error)

	// GetSources returns the map of parent Nodes.
	// (Nodes that come towards the argument vertex.)
	GetSources(id ID) (map[ID]Node, error)

	// GetTargets returns the map of child Nodes.
	// (Nodes that go out of the argument vertex.)
	GetTargets(id ID) (map[ID]Node, error)

	// String describes the Graph.
	String() string
}
// graph is an internal default graph type that
// implements all methods in Graph interface.
type graph struct {
	// idToNodes stores all nodes.
	idToNodes map[ID]Node
	// nodeToSources maps a Node identifer to sources(parents) with edge weights.
	nodeToSources map[ID]map[ID]Edge
	// nodeToTargets maps a Node identifer to targets(children) with edge weights.
	nodeToTargets map[ID]map[ID]Edge
}

func newGraph() *graph {
	return &graph{
		idToNodes:     make(map[ID]Node),
		nodeToSources: make(map[ID]map[ID]Edge),
		nodeToTargets: make(map[ID]map[ID]Edge),
	}
}

// NewGraph returns a new graph.
func NewGraph() Graph {
	return newGraph()
}

func (g *graph) ExistNode(id ID) bool {
	_, ok := g.idToNodes[id]
	return ok
}

func (g *graph) GetNode(id ID) Node    { return g.idToNodes[id] }
func (g *graph) GetNodes() map[ID]Node { return g.idToNodes }

func (g *graph) AddNode(nd Node) bool {
	if g.ExistNode(nd.ID()) {
		return false
	}
	id := nd.ID()
	g.idToNodes[id] = nd
	return true
}

func (g *graph) GetEdge(id1, id2 ID) (Edge) {
	if _, ok := g.nodeToTargets[id1]; ok {
		if v, ok := g.nodeToTargets[id1][id2]; ok {
			return v
		}
	}
	return nil
}

func (g *graph) ReplaceEdge(id1, id2 ID, edge Edge) error {
	if !g.ExistNode(id1) {
		return fmt.Errorf("%s does not exist in the graph", id1)
	}
	if !g.ExistNode(id2) {
		return fmt.Errorf("%s does not exist in the graph", id2)
	}

	if _, ok := g.nodeToTargets[id1]; ok {
		g.nodeToTargets[id1][id2] = edge
	} else {
		edges := make(map[ID]Edge)
		edges[id2] = edge
		g.nodeToTargets[id1] = edges
	}

	if _, ok := g.nodeToSources[id2]; ok {
		g.nodeToSources[id2][id1] = edge
	} else {
		edges := make(map[ID]Edge)
		edges[id1] = edge
		g.nodeToSources[id2] = edges
	}
	return nil
}

func (g *graph) GetSources(id ID) (map[ID]Node, error) {
	if !g.ExistNode(id) {
		return nil, fmt.Errorf("%s does not exist in the graph.", id)
	}

	rs := make(map[ID]Node)
	if _, ok := g.nodeToSources[id]; ok {
		for n := range g.nodeToSources[id] {
			rs[n] = g.idToNodes[n]
		}
	}
	return rs, nil
}

func (g *graph) GetTargets(id ID) (map[ID]Node, error) {
	if !g.ExistNode(id) {
		return nil, fmt.Errorf("%s does not exist in the graph", id)
	}

	rs := make(map[ID]Node)
	if _, ok := g.nodeToTargets[id]; ok {
		for n := range g.nodeToTargets[id] {
			rs[n] = g.idToNodes[n]
		}
	}
	return rs, nil
}

func (g *graph) String() string {
	buf := new(bytes.Buffer)
	for id1, nd1 := range g.idToNodes {
		nmap, _ := g.GetTargets(id1)
		for id2, nd2 := range nmap {
			edge := g.GetEdge(id1, id2)
			fmt.Fprintf(buf, "%s -→ %s : %s \n", nd1, nd2, edge.Train().String())
		}
	}
	return buf.String()
}
