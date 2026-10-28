package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func readBoard(r io.Reader) (Board, error) {
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return Board{}, errors.New("invalid board size")
	}
	if n < 2 || n > 10 {
		return Board{}, errors.New("size must be 2..10")
	}
	tiles := make([][]int, n)
	for i := range n {
		tiles[i] = make([]int, n)
		for j := range n {
			if _, err := fmt.Fscan(r, &tiles[i][j]); err != nil {
				return Board{}, fmt.Errorf("expected %d numbers", n*n)
			}
		}
	}
	if err := ValidateTiles(tiles); err != nil {
		return Board{}, err
	}
	return NewBoard(tiles), nil
}

func run() error {
	args := os.Args[1:]

	if len(args) == 1 && args[0] == "--gui" {
		runGUI()
		return nil
	}

	var initial Board
	var algo string

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Enter n, then n rows (0 = blank):")
		board, err := readBoard(os.Stdin)
		if err != nil {
			return err
		}
		initial = board
		algo = "astar-manhattan"
	} else if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		board, err := readBoard(f)
		if err != nil {
			return err
		}
		initial = board
		algo = "astar-manhattan"
	} else if len(args) == 2 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		board, err := readBoard(f)
		if err != nil {
			return err
		}
		initial = board
		algo = args[1]
	} else {
		return errors.New("usage: go run . [file] [algo]")
	}

	s, err := NewSolver(initial, algo)
	if err != nil {
		return err
	}

	if !s.Solvable() {
		fmt.Println("Unsolvable puzzle")
		fmt.Fprintf(os.Stderr, "Algorithm: %s, expanded nodes: %d\n", algo, s.Expanded())
		return nil
	}

	moves := s.Moves()
	cost := s.Cost()

	switch algo {
	case "bfs", "astar-manhattan", "astar-hamming":
		fmt.Printf("Minimum number of moves = %d\n", moves)
	case "ucs":
		fmt.Printf("Number of moves = %d\n", moves)
		fmt.Printf("Total cost = %d\n", cost)
	case "dfs":
		fmt.Printf("Number of moves = %d\n", moves)
	}

	sol := s.Solution()
	for i, b := range sol {
		fmt.Println(b.String())
		if i < len(sol)-1 {
			fmt.Println()
		}
	}

	fmt.Fprintf(os.Stderr, "Algorithm: %s, expanded nodes: %d\n", algo, s.Expanded())
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
