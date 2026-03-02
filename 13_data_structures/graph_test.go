package main

import (
	"testing"
)

func TestAddAndGetNode(t *testing.T) {
	g := NewGraph[string, int]()

	if !g.AddNode(1, "A") {
		t.Fatal("expected AddNode to succeed")
	}

	if g.AddNode(1, "Duplicate") {
		t.Fatal("expected AddNode to fail for duplicate id")
	}

	node, ok := g.GetNode(1)
	if !ok {
		t.Fatal("expected GetNode to find node")
	}

	if node.data != "A" {
		t.Fatalf("expected data A, got %v", node.data)
	}
}

func TestRemoveNode(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddEdge(1, 2, 10)

	if !g.RemoveNode(1) {
		t.Fatal("expected RemoveNode to succeed")
	}

	if _, ok := g.GetNode(1); ok {
		t.Fatal("expected node 1 to be removed")
	}

	if g.HasEdge(1, 2) {
		t.Fatal("expected edge to be removed when node deleted")
	}
}

func TestAddAndGetEdge(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	if !g.AddEdge(1, 2, 5) {
		t.Fatal("expected AddEdge to succeed")
	}

	weight, ok := g.GetEdgeWeight(1, 2)
	if !ok {
		t.Fatal("expected GetEdgeWeight to succeed")
	}

	if weight != 5 {
		t.Fatalf("expected weight 5, got %v", weight)
	}
}

func TestRemoveEdge(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddEdge(1, 2, 5)

	if !g.RemoveEdge(1, 2) {
		t.Fatal("expected RemoveEdge to succeed")
	}

	if g.HasEdge(1, 2) {
		t.Fatal("expected edge to be removed")
	}
}

func TestGetNodes(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	nodes := g.GetNodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestBFS(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(1, 3, 1)

	visited := []int{}

	g.BFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 3 {
		t.Fatalf("expected 3 visited nodes, got %d", len(visited))
	}
}

func TestBFS_NilGraph(t *testing.T) {
	var g *Graph[string, int]

	called := false
	g.BFS(1, func(n *Node[string]) bool {
		called = true
		return true
	})

	if called {
		t.Fatal("expected BFS on nil graph not to call visit")
	}
}

func TestBFS_StartNodeDoesNotExist(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")

	called := false
	g.BFS(999, func(n *Node[string]) bool {
		called = true
		return true
	})

	if called {
		t.Fatal("expected BFS not to run when start node does not exist")
	}
}

func TestBFS_SingleNode(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")

	var visited []int
	g.BFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 1 || visited[0] != 1 {
		t.Fatal("expected only the single node to be visited")
	}
}

func TestBFS_DisconnectedGraph(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C") // disconnected

	g.AddEdge(1, 2, 1)

	var visited []int
	g.BFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 2 {
		t.Fatalf("expected 2 reachable nodes, got %d", len(visited))
	}
}

func TestBFS_Cycle(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)
	g.AddEdge(3, 1, 1) // cycle

	count := 0
	g.BFS(1, func(n *Node[string]) bool {
		count++
		return true
	})

	if count != 3 {
		t.Fatalf("expected each node visited once in cycle, got %d", count)
	}
}

func TestBFS_EarlyStop(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)

	count := 0
	g.BFS(1, func(n *Node[string]) bool {
		count++
		return false // stop immediately
	})

	if count != 1 {
		t.Fatalf("expected traversal to stop after first visit, got %d visits", count)
	}
}

func TestBFS_NoOutgoingEdges(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	// No edges

	var visited []int
	g.BFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 1 {
		t.Fatalf("expected only start node visited, got %d", len(visited))
	}
}

func TestDFS(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)

	visited := []int{}

	g.DFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 3 {
		t.Fatalf("expected 3 visited nodes, got %d", len(visited))
	}
}

func TestDFS_NilGraph(t *testing.T) {
	var g *Graph[string, int]

	called := false
	g.DFS(1, func(n *Node[string]) bool {
		called = true
		return true
	})

	if called {
		t.Fatal("expected DFS on nil graph not to call visit")
	}
}

func TestDFS_StartNodeDoesNotExist(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")

	called := false
	g.DFS(999, func(n *Node[string]) bool {
		called = true
		return true
	})

	if called {
		t.Fatal("expected DFS not to run when start node does not exist")
	}
}

func TestDFS_SingleNode(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")

	var visited []int
	g.DFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 1 || visited[0] != 1 {
		t.Fatal("expected only the single node to be visited")
	}
}

func TestDFS_DisconnectedGraph(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C") // disconnected

	g.AddEdge(1, 2, 1)

	var visited []int
	g.DFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 2 {
		t.Fatalf("expected 2 reachable nodes, got %d", len(visited))
	}
}

func TestDFS_Cycle(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)
	g.AddEdge(3, 1, 1) // cycle

	count := 0
	g.DFS(1, func(n *Node[string]) bool {
		count++
		return true
	})

	if count != 3 {
		t.Fatalf("expected each node visited once in cycle, got %d", count)
	}
}

func TestDFS_EarlyStop(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)

	count := 0
	g.DFS(1, func(n *Node[string]) bool {
		count++
		return false // stop immediately
	})

	if count != 1 {
		t.Fatalf("expected traversal to stop after first visit, got %d visits", count)
	}
}

func TestDFS_NoOutgoingEdges(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	var visited []int
	g.DFS(1, func(n *Node[string]) bool {
		visited = append(visited, n.id)
		return true
	})

	if len(visited) != 1 {
		t.Fatalf("expected only start node visited, got %d", len(visited))
	}
}

func TestShortestPath(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")
	g.AddNode(3, "C")

	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 2)
	g.AddEdge(1, 3, 10)

	path := g.ShortestPath(1, 3)
	if path == nil {
		t.Fatal("expected path, got nil")
	}

	if len(path) != 3 {
		t.Fatalf("expected path length 3, got %d", len(path))
	}

	if path[0].id != 1 || path[1].id != 2 || path[2].id != 3 {
		t.Fatal("unexpected shortest path")
	}
}

func TestShortestPathNoPath(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	path := g.ShortestPath(1, 2)
	if path != nil {
		t.Fatal("expected nil path when no connection exists")
	}
}

func TestShortestPathNegativeCycle(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")
	g.AddNode(2, "B")

	g.AddEdge(1, 2, -1)
	g.AddEdge(2, 1, -1)

	path := g.ShortestPath(1, 2)
	if path != nil {
		t.Fatal("expected nil due to negative cycle")
	}
}

func TestShortestPathSameSourceDest(t *testing.T) {
	g := NewGraph[string, int]()
	g.AddNode(1, "A")

	path := g.ShortestPath(1, 1)
	if path == nil || len(path) != 1 || path[0].id != 1 {
		t.Fatal("expected single-node path")
	}
}

func TestNilGraphSafety(t *testing.T) {
	var g *Graph[string, int]

	if g.AddNode(1, "A") {
		t.Fatal("expected AddNode on nil graph to fail")
	}

	if g.AddEdge(1, 2, 1) {
		t.Fatal("expected AddEdge on nil graph to fail")
	}

	if g.RemoveNode(1) {
		t.Fatal("expected RemoveNode on nil graph to fail")
	}

	if g.RemoveEdge(1, 2) {
		t.Fatal("expected RemoveEdge on nil graph to fail")
	}

	if g.GetNodes() != nil {
		t.Fatal("expected GetNodes on nil graph to return nil")
	}

	if g.ShortestPath(1, 2) != nil {
		t.Fatal("expected ShortestPath on nil graph to return nil")
	}
}
