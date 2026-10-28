package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func runGUI() {
	a := app.New()
	w := a.NewWindow("8-Puzzle Solver")
	w.Resize(fyne.NewSize(520, 700))

	currentN := 3
	var entries [][]*widget.Entry
	var boardLabels []*widget.Label
	var boardView *fyne.Container
	var solution []Board

	var (
		gridHolder    *fyne.Container
		gridContainer *fyne.Container
		errorLabel    *widget.Label
		resultLabel   *widget.Label
		algoSelect    *widget.Select
		sizeSelect    *widget.Select
		solveBtn      *widget.Button
	)

	buildGrid := func(n int) (*fyne.Container, [][]*widget.Entry) {
		grid := container.NewGridWithColumns(n)
		es := make([][]*widget.Entry, n)
		for r := range n {
			es[r] = make([]*widget.Entry, n)
			for c := range n {
				e := widget.NewEntry()
				e.SetPlaceHolder("0")
				e.TextStyle = fyne.TextStyle{Monospace: true}
				es[r][c] = e
				grid.Add(e)
			}
		}
		return grid, es
	}

	buildBoardView := func(n int) *fyne.Container {
		boardLabels = make([]*widget.Label, n*n)
		grid := container.NewGridWithColumns(n)
		for i := range boardLabels {
			boardLabels[i] = widget.NewLabel("")
			boardLabels[i].Alignment = fyne.TextAlignCenter
			boardLabels[i].TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
			grid.Add(boardLabels[i])
		}
		return grid
	}

	renderBoard := func(b Board) {
		fyne.Do(func() {
			n := b.Size()
			for r := range n {
				for c := range n {
					v := b.TileAt(r, c)
					if v == 0 {
						boardLabels[r*n+c].SetText("")
					} else {
						boardLabels[r*n+c].SetText(strconv.Itoa(v))
					}
				}
			}
		})
	}

	fillGrid := func(b Board) {
		n := b.Size()
		for r := range n {
			for c := range n {
				v := b.TileAt(r, c)
				if v == 0 {
					entries[r][c].SetText("")
				} else {
					entries[r][c].SetText(strconv.Itoa(v))
				}
			}
		}
	}

	readGrid := func() (Board, error) {
		tiles := make([][]int, currentN)
		for r := range currentN {
			tiles[r] = make([]int, currentN)
			for c := range currentN {
				text := entries[r][c].Text
				if text == "" {
					tiles[r][c] = 0
					continue
				}
				v, err := strconv.Atoi(text)
				if err != nil {
					return Board{}, fmt.Errorf("cell (%d,%d): not a number", r, c)
				}
				tiles[r][c] = v
			}
		}
		if err := ValidateTiles(tiles); err != nil {
			return Board{}, err
		}
		return NewBoard(tiles), nil
	}

	gridContainer, entries = buildGrid(currentN)
	gridHolder = container.NewStack(gridContainer)

	boardView = buildBoardView(currentN)
	renderBoard(GoalBoard(currentN))

	errorLabel = widget.NewLabel("")
	errorLabel.Wrapping = fyne.TextWrapWord
	errorLabel.TextStyle = fyne.TextStyle{Bold: true}

	resultLabel = widget.NewLabel("")
	resultLabel.Wrapping = fyne.TextWrapWord

	algoSelect = widget.NewSelect(
		[]string{"A* Manhattan", "A* Hamming", "BFS", "UCS (tile cost)", "DFS (not optimal)"}, nil)
	algoSelect.SetSelected("A* Manhattan")

	algoMap := map[string]string{
		"A* Manhattan":      "astar-manhattan",
		"A* Hamming":        "astar-hamming",
		"BFS":               "bfs",
		"UCS (tile cost)":   "ucs",
		"DFS (not optimal)": "dfs",
	}

	solveBtn = widget.NewButton("Solve", func() {
		solveBtn.Disable()
		errorLabel.SetText("")
		resultLabel.SetText("")
		solution = nil

		board, err := readGrid()
		if err != nil {
			errorLabel.SetText("Invalid board: " + err.Error())
			solveBtn.Enable()
			return
		}
		algo := algoMap[algoSelect.Selected]
		go func() {
			s, err := NewSolver(board, algo)
			fyne.Do(func() {
				solveBtn.Enable()
				if err != nil {
					errorLabel.SetText("Solver error: " + err.Error())
					return
				}
				if !s.Solvable() {
					resultLabel.SetText("Unsolvable puzzle")
					return
				}
				solution = s.Solution()
				renderBoard(solution[len(solution)-1])
				switch algo {
				case "ucs":
					resultLabel.SetText(fmt.Sprintf("Number of moves = %d, Total cost = %d", s.Moves(), s.Cost()))
				case "dfs":
					resultLabel.SetText(fmt.Sprintf("Number of moves = %d", s.Moves()))
				default:
					resultLabel.SetText(fmt.Sprintf("Minimum number of moves = %d", s.Moves()))
				}
			})
		}()
	})

	loadBtn := widget.NewButton("Load File", func() {
		dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			defer rc.Close()
			board, err := readBoard(rc)
			if err != nil {
				errorLabel.SetText("Load error: " + err.Error())
				return
			}
			fillGrid(board)
			errorLabel.SetText("")
			resultLabel.SetText("")
			solution = nil
			renderBoard(board)
		}, w).Show()
	})

	randomBtn := widget.NewButton("Random", func() {
		b := GoalBoard(currentN)
		for i := 0; i < 80; i++ {
			neighbors := b.Neighbors()
			if len(neighbors) == 0 {
				break
			}
			b = neighbors[rand.IntN(len(neighbors))]
		}
		fillGrid(b)
		errorLabel.SetText("")
		resultLabel.SetText("")
		solution = nil
		renderBoard(b)
	})

	sizeSelect = widget.NewSelect([]string{"3", "4", "5"}, func(s string) {
		n, _ := strconv.Atoi(s)
		currentN = n
		gridContainer, entries = buildGrid(n)
		gridHolder.Objects = []fyne.CanvasObject{gridContainer}
		gridHolder.Refresh()
		boardView = buildBoardView(n)
		renderBoard(GoalBoard(n))
		resultLabel.SetText("")
		solution = nil
	})
	sizeSelect.SetSelected("3")

	inputSection := container.NewVBox(
		widget.NewLabelWithStyle("Input Board", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(widget.NewLabel("Size:"), sizeSelect),
		gridHolder,
		container.NewHBox(loadBtn, randomBtn),
	)
	algoSection := container.NewVBox(
		widget.NewLabelWithStyle("Algorithm", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		algoSelect,
		solveBtn,
	)
	resultSection := container.NewVBox(
		widget.NewLabelWithStyle("Solution", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		boardView,
		resultLabel,
	)
	content := container.NewVScroll(container.NewVBox(
		inputSection,
		widget.NewSeparator(),
		algoSection,
		widget.NewSeparator(),
		resultSection,
		widget.NewSeparator(),
		errorLabel,
	))

	w.SetContent(content)
	w.ShowAndRun()
}
