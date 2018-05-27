package main

import (
	"container/heap"
	"fmt"
	"math"
)

func main() {
	// Parse info from xml to graph
	g := NewGraphFromXML("data.xml")
	// find and output all paths fom all nodes
	PrintAllRoutes(g)
}

func NewGraphFromXML(filename string) Graph {
	trains := GetTrainsFromXML(filename)

	g := NewGraph()

	for _, train := range trains {
		id1 := train.DepartureStationId
		node1 := g.GetNode(StringID(id1))
		if node1 == nil {
			//	Add new node
			node1 = NewNode(id1)
			g.AddNode(node1)
		}

		id2 := train.ArrivalStationId
		node2 := g.GetNode(StringID(id2))
		if node2 == nil {
			//	Add new node
			node2 = NewNode(id2)
		}
		edge := NewEdge(node1, node2, train)
		existingEdge := g.GetEdge(node1.ID(), node2.ID())

		if existingEdge == nil || existingEdge.Weight() > edge.Weight() {
			g.ReplaceEdge(node1.ID(), node2.ID(), edge)
		}
	}
	return g
}

func GetTrainsFromXML(filename string) []Train {
	var trains []Train
	b, _ := parseXmlFile(filename)
	t, _ := unmarshalTickets(b)
	for _, train := range t {
		train.FindDuration()
		trains = append(trains, train)
	}
	return trains
}

// Dijkstra returns the shortest path using Dijkstra
// algorithm with a min-priority queue. This algorithm
// does not work with negative weight edges.
// (https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)
//
//	 0. Dijkstra(G, source, target)
//	 1.
//	 2. 	let Q be a priority queue
//	 3. 	distance[source] = 0
//	 4.
//	 5. 	for each vertex v in G:
//	 6.
//	 7. 		if v ≠ source:
//	 8. 			distance[v] = ∞
//	 9. 			prev[v] = undefined
//	10.
//	11. 		Q.add_with_priority(v, distance[v])
//	12.
//	13. 	while Q is not empty:
//	14.
//	15. 		u = Q.extract_min()
//	16. 		if u == target:
//	17. 			break
//	18.
//	19. 		for each child vertex v of u:
//	20.
//	21. 			alt = distance[u] + weight(u, v)
//	22. 			if distance[v] > alt:
//	23. 				distance[v] = alt
//	24. 				prev[v] = u
//	25. 				Q.decrease_priority(v, alt)
//	26.
//	27. 		reheapify(Q)
//	28.
//	29.
//	30. 	path = []
//	31. 	u = target
//	32. 	while prev[u] is defined:
//	33. 		path.push_front(u)
//	34. 		u = prev[u]
//	35.
//	36. 	return path, prev


func Dijkstra(g Graph, source, target ID) ([]ID, error) {
	// let Q be a priority queue
	minHeap := new(nodeDistanceHeap)

	// distance[source] = 0
	distance := make(map[ID]float64)
	distance[source] = 0.0

	// for each vertex v in G:
	for id := range g.GetNodes() {
		// if v ≠ source:
		if id != source {
			// distance[v] = ∞
			distance[id] = math.MaxFloat64
		}

		// Q.add_with_priority(v, distance[v])
		nds := nodeDistance{}
		nds.id = id
		nds.distance = distance[id]

		heap.Push(minHeap, nds)
	}

	heap.Init(minHeap)
	prev := make(map[ID]ID)

	// while Q is not empty:
	for minHeap.Len() != 0 {
		// u = Q.extract_min()
		u := heap.Pop(minHeap).(nodeDistance)
		// if u == target:
		if u.id == target {
			break
		}

		// for each child vertex v of u:
		chNodes, err := g.GetTargets(u.id)
		if err != nil {
			return nil, err
		}
		for v := range chNodes {

			// alt = distance[u] + weight(u, v)
			weight := g.GetEdge(u.id, v).Weight()
			alt := distance[u.id] + weight

			// if distance[v] > alt:
			if distance[v] > alt {

				// distance[v] = alt
				distance[v] = alt

				// prev[v] = u
				prev[v] = u.id

				// Q.decrease_priority(v, alt)
				minHeap.updateDistance(v, alt)
			}
		}
		heap.Init(minHeap)
	}

	// path = []
	var path []ID

	// u = target
	u := target

	// while prev[u] is defined:
	for {
		if _, ok := prev[u]; !ok {
			break
		}
		// path.push_front(u)
		temp := make([]ID, len(path)+1)
		temp[0] = u
		copy(temp[1:], path)
		path = temp

		// u = prev[u]
		u = prev[u]

	}
	// add the source
	temp := make([]ID, len(path)+1)
	temp[0] = source
	copy(temp[1:], path)
	path = temp
	return path, nil
}

func PrintPath(graph Graph, path []ID) {
	for i := 1; i < len(path); i++ {
		fmt.Println(graph.GetEdge(path[i-1], path[i]).Train())
	}
	fmt.Println()
}

// Prints shortest paths from every station to every station
func PrintAllRoutes(graph Graph) {
	nodes := graph.GetNodes()
	for id1 := range nodes {
		fmt.Printf("<--- %s --- > \n", id1)
		for id2 := range nodes {
			if id1 != id2 {
				fmt.Printf("%s --- > %s \n",id1,id2)
				path, _ := Dijkstra(graph, id1, id2)
				PrintPath(graph, path)
			}
		}
	}
}
