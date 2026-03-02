package main

import "slices"

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Node[T comparable] struct {
	id   int
	data T
}

func NewNode[T comparable](id int, data T) *Node[T] {
	return &Node[T]{
		id:   id,
		data: data,
	}
}

type Graph[K comparable, V Ordered] struct {
	nodes map[int]*Node[K]
	edges map[*Node[K]]*map[*Node[K]]V
}

func NewGraph[K comparable, V Ordered]() *Graph[K, V] {
	return &Graph[K, V]{
		nodes: map[int]*Node[K]{},
		edges: map[*Node[K]]*map[*Node[K]]V{},
	}
}

func (g *Graph[K, V]) AddNode(id int, data K) bool {
	if g == nil {
		return false
	}
	// verify there's no other node with the same id
	if _, ok := g.nodes[id]; ok {
		return false
	}
	g.nodes[id] = NewNode(id, data)
	return true
}

func (g *Graph[K, V]) AddEdge(sourceId int, destId int, weight V) bool {
	if g == nil {
		return false
	}
	// verify/get both source and dest exist
	var sourceNode, destNode *Node[K]
	var ok bool
	if sourceNode, ok = g.GetNode(sourceId); !ok {
		return false
	}
	if destNode, ok = g.GetNode(destId); !ok {
		return false
	}
	if g.edges[sourceNode] == nil {
		g.edges[sourceNode] = &map[*Node[K]]V{}
	}
	(*g.edges[sourceNode])[destNode] = weight
	return true
}

func (g *Graph[K, V]) GetNode(id int) (*Node[K], bool) {
	if g == nil {
		return nil, false
	}
	var node *Node[K]
	var ok bool
	if node, ok = g.nodes[id]; !ok {
		return nil, false
	}
	return node, true
}

func (g *Graph[K, V]) RemoveNode(id int) bool {
	if g == nil {
		return false
	}
	// Get node
	var node *Node[K]
	var ok bool
	if node, ok = g.nodes[id]; !ok {
		return false
	}

	// Delete outgoing edges
	delete(g.edges, node)

	// Delete incoming edges
	for _, edges := range g.edges {
		delete(*edges, node)
	}
	// Delete node from graph
	delete(g.nodes, id)
	return true
}

func (g *Graph[K, V]) GetEdgeWeight(sourceId int, destId int) (V, bool) {
	var zero V

	if g == nil {
		return zero, false
	}

	sourceNode, ok := g.GetNode(sourceId)
	if !ok {
		return zero, false
	}
	destNode, ok := g.GetNode(destId)
	if !ok {
		return zero, false
	}

	edges, ok := g.edges[sourceNode]
	if !ok {
		return zero, false
	}

	weight, ok := (*edges)[destNode]
	if !ok {
		return zero, false
	}
	return weight, true
}

func (g *Graph[K, V]) RemoveEdge(sourceId int, destId int) bool {
	if g == nil {
		return false
	}

	sourceNode, ok := g.GetNode(sourceId)
	if !ok {
		return false
	}
	destNode, ok := g.GetNode(destId)
	if !ok {
		return false
	}

	edges, ok := g.edges[sourceNode]
	if !ok {
		return false
	}

	if _, ok := (*edges)[destNode]; !ok {
		return false
	}
	delete(*edges, destNode)
	return true
}

func (g *Graph[K, V]) GetNodes() []*Node[K] {
	if g == nil {
		return nil
	}
	nodes := make([]*Node[K], 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *Graph[K, V]) GetEdges() map[*Node[K]]*map[*Node[K]]V {
	if g == nil {
		return nil
	}
	return g.edges
}

func (g *Graph[K, V]) HasEdge(sourceId int, destId int) bool {
	if g == nil {
		return false
	}

	sourceNode, ok := g.GetNode(sourceId)
	if !ok {
		return false
	}
	destNode, ok := g.GetNode(destId)
	if !ok {
		return false
	}

	edges, ok := g.edges[sourceNode]
	if !ok {
		return false
	}

	if _, ok := (*edges)[destNode]; !ok {
		return false
	}

	return true

}

// DFS performs a depth-first traversal starting from startId.
//
// The visit function is called when a node is dequeued for processing.
// If visit returns false, the traversal stops immediately.
// If visit returns true, traversal continues.
//
// Behavior:
//
//   - If the graph is nil, BFS returns immediately.
//   - If startId does not exist in the graph, BFS returns immediately.
//   - Each reachable node is visited at most once.
func (g *Graph[K, V]) DFS(startId int, visit func(*Node[K]) bool) {
	if g == nil {
		return
	}
	startNode, ok := g.GetNode(startId)
	if !ok {
		return
	}
	toVisit := []*Node[K]{startNode}
	visited := make(map[int]struct{})
	var curNode *Node[K]
	for len(toVisit) != 0 {
		curNode, toVisit = toVisit[len(toVisit)-1], toVisit[:len(toVisit)-1]
		if _, ok := visited[curNode.id]; ok {
			continue
		}
		// Visit curNode
		visited[curNode.id] = struct{}{}
		if shouldContinue := visit(curNode); !shouldContinue {
			return
		}
		// Explore neighbors
		outgoingEdges, ok := g.edges[curNode]
		if !ok {
			continue
		}
		for neighbor := range *outgoingEdges {
			if _, ok := visited[neighbor.id]; ok {
				continue
			}
			toVisit = append(toVisit, neighbor)
		}
	}
}

// BFS performs a breadth-first traversal starting from startId.
//
// The visit function is called when a node is dequeued for processing.
// If visit returns false, the traversal stops immediately.
// If visit returns true, traversal continues.
//
// Behavior:
//
//   - If the graph is nil, BFS returns immediately.
//   - If startId does not exist in the graph, BFS returns immediately.
//   - Each reachable node is visited at most once.
func (g *Graph[K, V]) BFS(startId int, visit func(*Node[K]) bool) {
	if g == nil {
		return
	}
	startNode, ok := g.GetNode(startId)
	if !ok {
		return
	}
	toVisit := []*Node[K]{startNode}
	visited := map[int]struct{}{}
	var curNode *Node[K]
	for len(toVisit) != 0 {
		curNode, toVisit = toVisit[0], toVisit[1:]
		if _, ok := visited[curNode.id]; ok {
			continue
		}
		// Visit curNode
		visited[curNode.id] = struct{}{}
		if shouldContinue := visit(curNode); !shouldContinue {
			return
		}
		// Explore neighbors
		outgoingEdges, ok := g.edges[curNode]
		if !ok {
			continue
		}
		for neighbor := range *outgoingEdges {
			if _, ok := visited[neighbor.id]; ok {
				continue
			}
			toVisit = append(toVisit, neighbor)
		}
	}

}

// ShortestPath computes the shortest path from sourceId to destId
// using the SPFA (Shortest Path Faster Algorithm).
//
// It returns the sequence of nodes representing the path from source
// to destination (inclusive). If no path exists, it returns nil.
//
// Behavior:
//
//   - If the graph is nil, it returns nil.
//   - If sourceId does not exist, it returns nil.
//   - If destId does not exist, it returns nil.
//   - If a negative-weight cycle is detected the function returns nil.
func (g *Graph[K, V]) ShortestPath(sourceId int, destId int) []*Node[K] {
	if g == nil {
		return nil
	}

	sourceNode, ok := g.GetNode(sourceId)
	if !ok {
		return nil
	}

	if sourceId == destId {
		return []*Node[K]{sourceNode}
	}

	_, ok = g.GetNode(destId)
	if !ok {
		return nil
	}

	// Simple queue
	toVisit := []*Node[K]{
		sourceNode,
	}
	dists := map[int]V{
		sourceId: V(0),
	}
	inQueue := map[int]struct{}{
		sourceId: struct{}{},
	}
	preds := map[int]int{
		sourceId: -1,
	}

	relaxCount := map[int]int{}

	var curNode *Node[K]
	for len(toVisit) != 0 {
		curNode, toVisit = toVisit[0], toVisit[1:]
		delete(inQueue, curNode.id)

		outgoingEdges, ok := g.edges[curNode]
		if !ok {
			continue
		}
		for neighbor, weight := range *outgoingEdges {
			alt := dists[curNode.id] + weight
			_, ok := dists[neighbor.id]
			if !ok || alt < dists[neighbor.id] {
				dists[neighbor.id] = alt
				preds[neighbor.id] = curNode.id
				// Increment relaxCount
				if _, ok := relaxCount[neighbor.id]; !ok {
					relaxCount[neighbor.id] = 0
				}
				relaxCount[neighbor.id]++

				// Check for negative cycle
				if relaxCount[neighbor.id] >= len(g.nodes) {
					return nil
				}

				// Only add to queue if not queued yet
				if _, ok := inQueue[neighbor.id]; ok {
					continue
				}
				inQueue[neighbor.id] = struct{}{}
				toVisit = append(toVisit, neighbor)
			}
		}
	}
	// Build path
	if _, ok := dists[destId]; !ok {
		return nil
	}
	path := []*Node[K]{}
	for cur := destId; cur != -1; cur = preds[cur] {
		curNode, ok := g.GetNode(cur)
		if !ok {
			panic("Node not found")
		}
		path = append(path, curNode)
	}
	slices.Reverse(path)
	return path
}
