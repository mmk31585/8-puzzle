package main

import (
	"errors"
	"fmt"
	"strings"
)

type Board struct {
	n              int
	tiles          [][]int
	manhattan      int 
	blankR, blankC int
}

var (
	ErrInvalidSize     = errors.New("board size must be >= 2")
	ErrNotSquare       = errors.New("board must be square")
	ErrInvalidValue    = errors.New("tile value out of range")
	ErrDuplicateValue  = errors.New("duplicate tile value")
)

func ValidateTiles(tiles [][]int) error {
	n := len(tiles)
	if n < 2 {
		return ErrInvalidSize
	}
	for _, row := range tiles {
		if len(row) != n {
			return ErrNotSquare
		}
	}
	seen := make(map[int]bool, n*n)
	for _, row := range tiles {
		for _, v := range row {
			if v < 0 || v >= n*n {
				return ErrInvalidValue
			}
			if seen[v] {
				return ErrDuplicateValue
			}
			seen[v] = true
		}
	}
	return nil
}

func NewBoard(tiles [][]int) Board {
	n := len(tiles)
	copyTiles := make([][]int, n)
	for i := range tiles {
		copyTiles[i] = make([]int, n)
		copy(copyTiles[i], tiles[i])
	}
	b := Board{n: n, tiles: copyTiles}
	sum := 0
	for r := range n {
		for c := range n {
			v := copyTiles[r][c]
			if v == 0 {
				b.blankR, b.blankC = r, c
				continue
			}
			goalR := (v - 1) / n
			goalC := (v - 1) % n
			sum += abs(r-goalR) + abs(c-goalC)
		}
	}
	b.manhattan = sum
	return b
}

func GoalBoard(n int) Board {
	tiles := make([][]int, n)
	val := 1
	for r := range n {
		tiles[r] = make([]int, n)
		for c := range n {
			if r == n-1 && c == n-1 {
				tiles[r][c] = 0
			} else {
				tiles[r][c] = val
				val++
			}
		}
	}
	return NewBoard(tiles)
}

func (b Board) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", b.n))
	for r := range b.n {
		sb.WriteString("\n")
		for c := range b.n {
			if c > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%d", b.tiles[r][c]))
		}
	}
	return sb.String()
}

func (b Board) Size() int {
	return b.n
}

func (b Board) TileAt(row, col int) int {
	if row < 0 || row >= b.n || col < 0 || col >= b.n {
		panic("tileAt: index out of range")
	}
	return b.tiles[row][col]
}

func (b Board) Hamming() int {
	count := 0
	for r := range b.n {
		for c := range b.n {
			v := b.tiles[r][c]
			if v == 0 {
				continue
			}
			goalR := (v - 1) / b.n
			goalC := (v - 1) % b.n
			if r != goalR || c != goalC {
				count++
			}
		}
	}
	return count
}

func (b Board) Manhattan() int {
	return b.manhattan
}

func (b Board) IsGoal() bool {
	return b.manhattan == 0
}

func (b Board) Equals(o Board) bool {
	if b.n != o.n {
		return false
	}
	for r := range b.n {
		for c := range b.n {
			if b.tiles[r][c] != o.tiles[r][c] {
				return false
			}
		}
	}
	return true
}

func (b Board) Key() string {
	var sb strings.Builder
	for r := range b.n {
		for c := range b.n {
			if r > 0 || c > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("%d", b.tiles[r][c]))
		}
	}
	return sb.String()
}

func (b Board) Neighbors() []Board {
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	neighbors := make([]Board, 0, 4)
	for _, d := range dirs {
		nr := b.blankR + d[0]
		nc := b.blankC + d[1]
		if nr < 0 || nr >= b.n || nc < 0 || nc >= b.n {
			continue
		}
		copyTiles := make([][]int, b.n)
		for i := range b.tiles {
			copyTiles[i] = make([]int, b.n)
			copy(copyTiles[i], b.tiles[i])
		}
		copyTiles[b.blankR][b.blankC], copyTiles[nr][nc] = copyTiles[nr][nc], copyTiles[b.blankR][b.blankC]
		neighbors = append(neighbors, NewBoard(copyTiles))
	}
	return neighbors
}

func (b Board) IsSolvable() bool {
	flat := make([]int, 0, b.n*b.n-1)
	for r := range b.n {
		for c := range b.n {
			v := b.tiles[r][c]
			if v != 0 {
				flat = append(flat, v)
			}
		}
	}
	inv := 0
	for i := range flat {
		for j := i + 1; j < len(flat); j++ {
			if flat[i] > flat[j] {
				inv++
			}
		}
	}
	if b.n%2 == 1 {
		return inv%2 == 0
	}
	return (inv+b.blankR)%2 == 1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}