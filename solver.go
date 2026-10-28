package main

import (
	"errors"
	"slices"
)

var ErrUnknownAlgo = errors.New("unknown algorithm")

type Solver struct {
	goal     *Node
	solvable bool
	expanded int
}

func NewSolver(initial Board, algo string) (*Solver, error) {
	var h func(Board) int

	switch algo {
	case "bfs", "dfs":

	case "ucs":

	case "astar-manhattan":
		h = func(b Board) int { return b.Manhattan() }
	case "astar-hamming":
		h = func(b Board) int { return b.Hamming() }
	default:
		return nil, ErrUnknownAlgo
	}

	s := &Solver{solvable: false, expanded: 0}

	if !initial.IsSolvable() {
		return s, nil
	}

	s.solvable = true

	frontier := NewFrontier(algo)

	var best map[string]int
	var visited map[string]bool

	startPrio := 0
	if h != nil {
		startPrio = h(initial)
	}
	startNode := &Node{board: initial, moves: 0, cost: 0, prio: startPrio, parent: nil}
	frontier.Push(startNode)

	if algo == "bfs" || algo == "dfs" {
		visited = make(map[string]bool)
		visited[initial.Key()] = true
	} else {
		best = make(map[string]int)
		best[initial.Key()] = 0
	}

	for frontier.Len() > 0 {
		n := frontier.Pop()
		s.expanded++

		if n.board.IsGoal() {
			s.goal = n
			break
		}

		for _, nb := range n.board.Neighbors() {

			if n.parent != nil && nb.Equals(n.parent.board) {
				continue
			}

			newMoves := n.moves + 1

			if algo == "ucs" {

				moved := nb.TileAt(n.board.blankR, n.board.blankC)
				newCost := n.cost + moved
				key := nb.Key()
				if prevCost, ok := best[key]; ok && prevCost <= newCost {
					continue
				}
				best[key] = newCost
				prio := newCost
				frontier.Push(&Node{board: nb, moves: newMoves, cost: newCost, prio: prio, parent: n})
			} else if algo == "bfs" || algo == "dfs" {
				key := nb.Key()
				if visited[key] {
					continue
				}
				visited[key] = true

				prio := 0
				frontier.Push(&Node{board: nb, moves: newMoves, cost: newMoves, prio: prio, parent: n})
			} else {

				newCost := n.cost + 1
				key := nb.Key()
				if prevCost, ok := best[key]; ok && prevCost <= newCost {
					continue
				}
				best[key] = newCost
				prio := newMoves + h(nb)
				frontier.Push(&Node{board: nb, moves: newMoves, cost: newCost, prio: prio, parent: n})
			}
		}
	}

	return s, nil
}

func (s *Solver) Solvable() bool {
	return s.solvable
}

func (s *Solver) Moves() int {
	if !s.solvable || s.goal == nil {
		return -1
	}
	return s.goal.moves
}

func (s *Solver) Cost() int {
	if !s.solvable || s.goal == nil {
		return -1
	}
	return s.goal.cost
}

func (s *Solver) Expanded() int {
	return s.expanded
}

func (s *Solver) Solution() []Board {
	if !s.solvable || s.goal == nil {
		return nil
	}
	path := make([]Board, 0)
	for n := s.goal; n != nil; n = n.parent {
		path = append(path, n.board)
	}
	slices.Reverse(path)
	return path
}
