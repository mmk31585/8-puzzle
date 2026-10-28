package main

import "container/heap"

type Node struct {
	board  Board
	moves  int
	cost   int
	prio   int
	parent *Node
}

type Frontier interface {
	Push(*Node)
	Pop() *Node
	Len() int
}

type pqFrontier []*Node

func (pq pqFrontier) Len() int { return len(pq) }

func (pq pqFrontier) Less(i, j int) bool {
	if pq[i].prio != pq[j].prio {
		return pq[i].prio < pq[j].prio
	}
	return pq[i].moves < pq[j].moves
}

func (pq pqFrontier) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pqFrontier) Push(x any) {
	*pq = append(*pq, x.(*Node))
}

func (pq *pqFrontier) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type pqFrontierWrapper struct {
	pq pqFrontier
}

func (fw *pqFrontierWrapper) Push(n *Node) {
	heap.Push(&fw.pq, n)
}

func (fw *pqFrontierWrapper) Pop() *Node {
	return heap.Pop(&fw.pq).(*Node)
}

func (fw *pqFrontierWrapper) Len() int {
	return fw.pq.Len()
}

type fifoFrontier []*Node

func (f *fifoFrontier) Push(n *Node) {
	*f = append(*f, n)
}

func (f *fifoFrontier) Pop() *Node {
	n := (*f)[0]
	*f = (*f)[1:]
	return n
}

func (f *fifoFrontier) Len() int {
	return len(*f)
}

type lifoFrontier []*Node

func (f *lifoFrontier) Push(n *Node) {
	*f = append(*f, n)
}

func (f *lifoFrontier) Pop() *Node {
	n := (*f)[len(*f)-1]
	*f = (*f)[:len(*f)-1]
	return n
}

func (f *lifoFrontier) Len() int {
	return len(*f)
}

func NewFrontier(algo string) Frontier {
	switch algo {
	case "bfs":
		return &fifoFrontier{}
	case "dfs":
		return &lifoFrontier{}
	case "ucs", "astar-manhattan", "astar-hamming":
		fw := &pqFrontierWrapper{pq: make(pqFrontier, 0)}
		heap.Init(&fw.pq)
		return fw
	default:
		return &pqFrontierWrapper{pq: make(pqFrontier, 0)}
	}
}