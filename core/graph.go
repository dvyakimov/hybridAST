package core

import (
	"fmt"
)

// Node a single node that composes the tree
type Node struct {
	Value  string
	IsFunc bool
	ID     int
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}

// ItemGraph the Items graph
type ItemGraph struct {
	nodes []*Node
	edges map[Node][]*Node
}

// AddNode adds a node to the graph
func (g *ItemGraph) AddNode(n *Node) {
	g.nodes = append(g.nodes, n)
}
func LastNode(g *ItemGraph) int {
	return len(g.nodes) - 1
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	//g.edges[*n2] = append(g.edges[*n2], n1)
}

// for debugging
func (g *ItemGraph) String() {
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
}

/* ------ Search--------- */

func FindInGraphByIndex(g *ItemGraph, index int, search string) int {
	near := g.edges[*g.nodes[index]]
	for j := 0; j < len(near); j++ {
		if near[j].String() == search {
			return near[j].ID
		}
	}
	return -1
}

func (g *ItemGraph) GetValueString(id int) string {
	if g.nodes[id] != nil {
		return g.nodes[id].String()
	} else {
		return ""
	}
}

func (g *ItemGraph) GetTheLastNodeValueString(id int) string {
	near := g.edges[*g.nodes[id]]
	if near[0] != nil {
		return near[0].String()
	} else {
		return ""
	}
}

// for debug
func (g *ItemGraph) FindInGraph(search string) int {
	for i := 0; i < len(g.nodes); i++ {
		if g.nodes[i].String() == search {
			near := g.edges[*g.nodes[i]]
			for j := 0; j < len(near); j++ {
				fmt.Println(near[j].String())
			}
			return g.nodes[i].ID
		}
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			if near[j].String() == search {
				return g.nodes[i].ID
			}
		}
	}
	return -1
}

/*-------BFS----------*/

// NodeQueue the queue of Nodes
type NodeQueue struct {
	items []Node
}

// New creates a new NodeQueue
func (s *NodeQueue) New() *NodeQueue {
	s.items = []Node{}
	return s
}

// Enqueue adds an Node to the end of the queue
func (s *NodeQueue) Enqueue(t Node) {
	s.items = append(s.items, t)
}

// Dequeue removes an Node from the start of the queue
func (s *NodeQueue) Dequeue() *Node {
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	return &item
}

// Front returns the item next in the queue, without removing it
func (s *NodeQueue) Front() *Node {
	item := s.items[0]
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *NodeQueue) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of Nodes in the queue
func (s *NodeQueue) Size() int {
	return len(s.items)
}

// Traverse implements the BFS traversing algorithm
func (g *ItemGraph) Traverse(f func(*Node)) {
	q := NodeQueue{}
	q.New()
	n := g.nodes[0]
	q.Enqueue(*n)
	visited := make(map[*Node]bool)
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		near := g.edges[*node]

		for i := 0; i < len(near); i++ {
			j := near[i]
			if !visited[j] {
				q.Enqueue(*j)
				visited[j] = true
			}
		}
		if f != nil {
			f(node)
		}
	}
}
