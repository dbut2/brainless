package sudoku

import (
	"github.com/puzzle-solvers/brainless/pkg/brainless"
	"net/http"
)

func Solve(w http.ResponseWriter, r *http.Request) {
	game := brainless.NewGame(w, r, brainless.Config{
		PuzzleURL:   "https://www.puzzle-sudoku.com",
		FunctionURL: "http://localhost:8080/sudoku",
	})
	game.ConnectBrain(&SudokuBrain{})
	game.HandleRequest()
}

type SudokuBrain struct {
	Task [][]int
	Game *Game
}

func (s *SudokuBrain) Setup() {
	g := &Game{}

	if len(s.Task) != 9 || len(s.Task[0]) != 9 {
		return
	}

	g.Cells = make([][]*Cell, 9)
	for i := 0; i < 9; i++ {
		g.Cells[i] = make([]*Cell, 9)
		for j := 0; j < 9; j++ {
			g.Cells[i][j] = &Cell{
				Game:   g,
				Row:    i,
				Col:    j,
				Box:    3*(i/3) + j/3,
				Number: s.Task[i][j],
			}
		}
	}

	s.Game = g
}

// check what numbers can exist in cell if 1 set

// check where number can sit in row if 1 set
// check where number can sit in col if 1 set
// check where number can sit in box if 1 set

func (s *SudokuBrain) Step() {
	for _, row := range s.Game.Cells {
		for _, cell := range row {
			if cell.Number != 0 {
				continue
			}

			potential := []int{}
			for i := 1; i <= 9; i++ {
				if cell.CanBe(i) {
					potential = append(potential, i)
				}
			}
			if len(potential) == 1 {
				cell.Number = potential[0]
				continue
			}
		}
	}
}

func (s *SudokuBrain) CheckDone() bool {
	for _, row := range s.Game.Cells {
		for _, cell := range row {
			if cell.Number == 0 {
				return false
			}
		}
	}
	return true
}

func (s *SudokuBrain) GetTask(task [][]int) {
	s.Task = task
}

type Response struct {
	CellStatus [][]CellResponse `json:"cellStatus"`
}

type CellResponse struct {
	Immutable     bool  `json:"immutable"`
	Number        int   `json:"number"`
	Pencil        bool  `json:"pencil"`
	PencilNumbers []int `json:"pencilNumbers"`
}

func (s *SudokuBrain) ToResponse() interface{} {
	cellStatus := make([][]CellResponse, 9)
	for i, row := range s.Game.Cells {
		cellStatus[i] = make([]CellResponse, 9)
		for j, cell := range row {
			cellStatus[i][j] = CellResponse{
				Immutable:     false,
				Number:        cell.Number,
				Pencil:        false,
				PencilNumbers: []int{},
			}
		}
	}
	return Response{
		CellStatus: cellStatus,
	}
}

type Game struct {
	Cells [][]*Cell
}

func (g *Game) Row(i int) []*Cell {
	return g.Cells[i]
}

func (g *Game) Col(j int) []*Cell {
	cells := []*Cell{}
	for _, row := range g.Cells {
		cells = append(cells, row[j])
	}
	return cells
}

func (g *Game) Box(b int) []*Cell {
	cells := []*Cell{}
	for i := b - (b % 3); i < b-(b%3)+3; i++ {
		for j := (b % 3) * 3; j < (b%3)*3+3; j++ {
			cells = append(cells, g.Cells[i][j])
		}
	}
	return cells
}

func (g *Game) NumberInRow(number int, row int) bool {
	for _, cell := range g.Row(row) {
		if cell.Number == number {
			return true
		}
	}
	return false
}

func (g *Game) NumberInCol(number int, col int) bool {
	for _, cell := range g.Col(col) {
		if cell.Number == number {
			return true
		}
	}
	return false
}

func (g *Game) NumberInBox(number int, box int) bool {
	for _, cell := range g.Box(box) {
		if cell.Number == number {
			return true
		}
	}
	return false
}

type Cell struct {
	Game          *Game
	Row, Col, Box int
	Number        int
}

func (c *Cell) CanBe(number int) bool {
	if c.Game.NumberInRow(number, c.Row) {
		return false
	}
	if c.Game.NumberInCol(number, c.Col) {
		return false
	}
	if c.Game.NumberInBox(number, c.Box) {
		return false
	}
	return true
}
