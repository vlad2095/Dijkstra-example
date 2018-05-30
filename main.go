package main

import (
	"container/heap"
	"fmt"
	"math"
	"time"
)

func main() {
	// Parse info from xml to graph
	g := NewGraphFromXML("data.xml")
	// find and output all paths fom all nodes
	PrintAllRoutes(g)
}

func NewGraphFromXML(filename string) *Graph {
	trains := GetTrainsFromXML(filename)

	g := NewGraph()

	for _, train := range trains {
		AddTrainToGraph(g, train)
	}
	return g
}

func AddTrainToGraph(g *Graph, train Train) {
	fromVertexId := train.DepartureStationId
	toVertexId := train.ArrivalStationId

	// Get vertex, create new one if not exist
	if _, ok := g.GetVertex(fromVertexId); !ok {
		v := NewVertex(fromVertexId)
		g.AddVertex(v)

	}
	fromVertex, _ := g.GetVertex(fromVertexId)

	// Get vertex, create new one if not exist
	if _, ok := g.GetVertex(toVertexId); !ok {
		v := NewVertex(toVertexId)
		g.AddVertex(v)
	}
	toVertex, _ := g.GetVertex(toVertexId)

	// create arc between vertexes if not exist
	if _, ok := fromVertex.GetOutgoingArc(toVertexId); !ok {
		arc := NewArc(fromVertexId, toVertexId)
		fromVertex.AddOutgoingArc(arc)
	}
	if _, ok := toVertex.GetIngoingArc(fromVertexId); !ok {
		arc := NewArc(fromVertexId, toVertexId)
		toVertex.AddIngoingArc(arc)
	}

	// add new edge to vertices arc
	// it's safe because we checked before and sure that arc is exist
	edge := NewEdge(train.ID, train)
	fromVertex.Out[toVertexId].AddEdge(edge)
	toVertex.In[fromVertexId].AddEdge(edge)
}

func GetTrainsFromXML(filename string) []Train {
	var trains []Train
	b, _ := parseXmlFile(filename)
	t, _ := unmarshalTickets(b)
	for _, train := range t {
		// convert time from string to time.Time
		train.ConvertTime()
		// duration = arrival time - departure time
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

func Dijkstra(g *Graph, source, target string) ([]string, error) {
	// let Q be a priority queue
	minHeap := new(nodeDistanceHeap)

	// distance[source] = 0
	distance := make(map[string]float64)
	distance[source] = 0.0

	// for each vertex v in G:
	for id := range g.Vertices {
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
	prev := make(map[string]string)
	//previousEdges := make(map[string]*Edge)

	// while Q is not empty:
	for minHeap.Len() != 0 {
		// u = Q.extract_min()
		u := heap.Pop(minHeap).(nodeDistance)
		// if u == target:
		if u.id == target {
			break
		}

		// for each child vertex v of u:
		uV, _ := g.GetVertex(u.id)
		chNodes := uV.Out
		for v, arc := range chNodes {
			weight := math.MaxFloat64
			// find smallest weight edge in arc
			for _, edge := range arc.edges {
				if ew := CalculateWeight(edge, nil); ew < weight {
					weight = ew
				}
			}
			// alt = distance[u] + weight(u, v)
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
	var path []string

	// u = target
	u := target

	// while prev[u] is defined:
	for {
		if _, ok := prev[u]; !ok {
			break
		}
		// path.push_front(u)
		temp := make([]string, len(path)+1)
		temp[0] = u
		copy(temp[1:], path)
		path = temp

		// u = prev[u]
		u = prev[u]

	}
	// add the source
	temp := make([]string, len(path)+1)
	temp[0] = source
	copy(temp[1:], path)
	path = temp
	return path, nil
}

// recieves one or two edges
// and returns relative weight
func CalculateWeight(edge1, edge2 *Edge) float64 {
	if edge2 == nil {
		train := edge1.data.(Train)
		//fmt.Println(train)
		return train.Price * float64(train.Duration)
	}

	train1 := edge1.data.(Train)
	train2 := edge2.data.(Train)

	// If second train outgoing before first comes (5 minute is min amount to be in time)
	// Then add 24 hour to second one departure time
	if !(train2.DepartureTime.After(train1.ArrivalTime.Add(5 * time.Minute))) {
		train2.DepartureTime = train2.DepartureTime.Add(time.Hour * 24)
		train2.ArrivalTime = train2.ArrivalTime.Add(time.Hour * 24)
	}
	totalCost := train1.Price + train2.Price
	// total duration contains duration of train routes and time beetween them
	totalDuration := train1.Duration + train2.Duration + train2.DepartureTime.Sub(train1.ArrivalTime)
	//	weight is a direct proportion of price and time
	weight := totalCost * float64(totalDuration)
	return weight
}

// BuildBetterRoute finds better possible route from given path
// return trains list (slice), total cost and total duration of route
func BuildBetterRoute(graph *Graph, path []string) (trains []*Train, totalCost float64, totalDuration time.Duration) {
	pathLen := len(path)
	if pathLen < 2 {
		// there is no path
		return
	}
	source := path[0]
	next := path[1]

	v, _ := graph.GetVertex(source)
	arc := v.Out[next]

	id := ""
	weight := math.MaxFloat64
	for _, edge := range arc.edges {
		if w := CalculateWeight(edge, nil); w < weight {
			id = edge.id
			weight = w
		}
	}
	train := arc.edges[id].data.(Train)
	trains = append(trains, &train)

	totalCost += train.Price
	totalDuration += train.Duration

	if pathLen == 2 {
		return
	}

	// set up previous edge and time  for calculations of weight in loop
	previousEdge := arc.edges[id]
	previousArrivalTime := train.ArrivalTime

	for i := 1; i < (pathLen - 1); i++ {
		source := path[i]
		next := path[i+1]
		v, _ := graph.GetVertex(source)
		arc := v.Out[next]

		id := ""
		weight := math.MaxFloat64
		for _, edge := range arc.edges {
			if w := CalculateWeight(previousEdge, edge); w < weight {
				id = edge.id
				weight = w
			}
		}
		// update train time
		train := arc.edges[id].data.(Train)
		if !(train.DepartureTime.After(previousArrivalTime.Add(5 * time.Minute))) {
			train.DepartureTime = train.DepartureTime.Add(time.Hour * 24)
			train.ArrivalTime = train.ArrivalTime.Add(time.Hour * 24)
		}


		totalCost += train.Price
		totalDuration += (train.Duration + train.DepartureTime.Sub(previousArrivalTime))
		trains = append(trains, &train)

		//	 Update previous edge and time
		previousEdge = arc.edges[id]
		previousArrivalTime = train.ArrivalTime

	}
	return
}

//Prints shortest paths from every station to every station
func PrintAllRoutes(graph *Graph) {
	nodes := graph.Vertices
	for id1 := range nodes {
		// station id
		fmt.Printf("<--- %s --- > \n", id1)
		for id2 := range nodes {
			if id1 != id2 {
				// route direction
				fmt.Printf("%s --- > %s \n", id1, id2)
				// find path
				path, err := Dijkstra(graph, id1, id2)
				if err != nil {
					fmt.Printf("Error: %s", err.Error())
					continue
				}
				// prints all trains of path
				route, cost, duration := BuildBetterRoute(graph, path)
				for _, train := range route {
					fmt.Printf("%s\n", train.String())
				}
				fmt.Printf("Total Cost: %.2f$, Total Duration:  %s\n---\n", cost, duration)
			}
		}
	}
}
