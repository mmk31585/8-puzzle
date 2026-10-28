# 8-Puzzle Solver

A Go implementation of the n×n sliding puzzle solver with multiple search algorithms.

## Quick Start

```bash
go run .                          # Read board from stdin
go run . examples/solvable.txt    # Read from file, default algorithm (A* Manhattan)
go run . examples/solvable.txt astar-manhattan  # Specify algorithm
go run . --gui                    # Launch GUI (requires display)
```

## Input Format

Board format (file or stdin):
```
n
row1_col1 row1_col2 ... row1_coln
row2_col1 row2_col2 ... row2_coln
...
rown_col1 rown_col2 ... rown_coln
```

- `n` = board size (2..10)
- Values: 0 = blank, 1..n²-1 = tiles
- Example (3×3):
```
3
0 1 3
4 2 5
7 8 6
```

## Algorithms

| Algorithm      | Description                              | Optimality |
|----------------|------------------------------------------|------------|
| `bfs`          | Breadth-first search                     | Optimal (moves) |
| `dfs`          | Depth-first search                       | NOT optimal |
| `ucs`          | Uniform-cost search (tile cost = value)  | Optimal (cost) |
| `astar-manhattan` | A* with Manhattan distance           | Optimal (moves) |
| `astar-hamming`   | A* with Hamming distance              | Optimal (moves) |

## Output Format

### Solvable
```
Minimum number of moves = 4          # bfs, astar-*
Number of moves = 4                  # dfs
Number of moves = 4                  # ucs
Total cost = 14                      # ucs only

3
0 1 3
4 2 5
7 8 6

3
1 0 3
4 2 5
7 8 6
... (all boards from start to goal, separated by blank line)
```

### Unsolvable
```
Unsolvable puzzle
```

### Stderr (all cases)
```
Algorithm: astar-manhattan, expanded nodes: 5
```

## Examples

```bash
# Solvable puzzle (4 moves optimal)
go run . examples/solvable.txt bfs
go run . examples/solvable.txt dfs
go run . examples/solvable.txt ucs
go run . examples/solvable.txt astar-manhattan
go run . examples/solvable.txt astar-hamming

# Unsolvable puzzle
go run . examples/unsolvable.txt astar-manhattan

# Stdin input
printf "3\n0 1 3\n4 2 5\n7 8 6\n" | go run .
```

## Running Tests

```bash
# Build and vet
go build ./... && go vet ./...

# All algorithms on solvable example
for algo in bfs dfs ucs astar-manhattan astar-hamming; do
  echo "=== $algo ==="
  go run . examples/solvable.txt $algo
done

# Unsolvable
go run . examples/unsolvable.txt astar-manhattan
```

## Project Structure

```
puzzle/
├── board.go      # Board type, validation, neighbors, solvability
├── frontier.go   # Frontier implementations (FIFO, LIFO, priority queue)
├── solver.go     # Search algorithms
├── main.go       # CLI entry point
├── gui.go        # Fyne GUI
├── go.mod
└── examples/
    ├── solvable.txt
    └── unsolvable.txt
```

## Notes

- The solver never mutates a Board after creation
- Solvability is checked before any search (instant for unsolvable)
- DFS is NOT optimal — path may be very long
- UCS minimizes total tile cost, not move count
- GUI requires a display (X11/Wayland/macOS/Windows); headless CI will panic