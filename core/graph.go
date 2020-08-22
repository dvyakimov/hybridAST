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
