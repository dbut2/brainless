package binairo

import (
	"net/http"

	"github.com/puzzle-solvers/brainless/pkg/brainless"
)

func Solve(w http.ResponseWriter, r *http.Request) {
	game := brainless.NewGame(w, r, brainless.Config{
		PuzzleURL:   "https://www.puzzle-binairo.com",
		FunctionURL: "http://localhost:8080/binairo",
	})
	brain := &BinairoBrain{}
	game.ConnectBrain(brain)
	game.HandleRequest()
}

type BinairoBrain struct {
	Task [][]int
	Game *Game
}

func (b *BinairoBrain) Setup() {
	g := &Game{}

	if len(b.Task) == 0 || len(b.Task[0]) == 0 {
		return
	}

	g.Task = b.Task

	rows, cols := len(b.Task), len(b.Task[0])

	g.Rows, g.Cols = rows, cols

	g.Cells = make([][]*Cell, rows)
	for i := range g.Cells {
		g.Cells[i] = make([]*Cell, cols)
		for j := range g.Cells {
			g.Cells[i][j] = g.NewCell(i, j)
		}
	}

	b.Game = g
}

func (b *BinairoBrain) Step() {
	for _, row := range b.Game.Cells {
		for _, cell := range row {
			if cell.Colour != Blank {
				continue
			}

			tile := NewTile()

			brainless.FlipperRule(tile, func(c brainless.Compass) {
				if cell.Neighbour(c.South).Colour == tile.Black && cell.Neighbour(c.South).Neighbour(c.South).Colour == tile.Black {
					cell.Colour = tile.White
				}

				if cell.Neighbour(c.South).Colour == tile.Black && cell.Neighbour(c.North).Colour == tile.Black {
					cell.Colour = tile.White
				}
			})
		}
	}

	for i := 0; i < b.Game.Rows; i++ {
		tile := NewTile()

		brainless.FlipperRule(tile, func(c brainless.Compass) {
			if b.Game.ColourInRow(tile.Black, i) == b.Game.Cols/2 {
				for _, cell := range b.Game.GetRow(i) {
					if cell.Colour == Blank {
						cell.Colour = tile.White
					}
				}
			}
		})
	}

	for j := 0; j < b.Game.Cols; j++ {
		tile := NewTile()

		brainless.FlipperRule(tile, func(c brainless.Compass) {
			if b.Game.ColourInCol(tile.Black, j) == b.Game.Rows/2 {
				for _, cell := range b.Game.GetCol(j) {
					if cell.Colour == Blank {
						cell.Colour = tile.White
					}
				}
			}
		})
	}

}

func (b *BinairoBrain) CheckDone() bool {
	for _, row := range b.Game.Cells {
		for _, cell := range row {
			if cell.Colour == Blank {
				return false
			}
		}
	}
	return true
}

func (b *BinairoBrain) GetTask(task [][]int) {
	b.Task = task
}

type Response struct {
	CellStatus [][]int `json:"cellStatus"`
}

func (b *BinairoBrain) ToResponse() interface{} {
	r := Response{}
	r.CellStatus = make([][]int, b.Game.Rows)
	for i, row := range b.Game.Cells {
		r.CellStatus[i] = make([]int, b.Game.Cols)
		for j, cell := range row {
			r.CellStatus[i][j] = cell.GetValue()
		}
	}
	return r
}

type Game struct {
	Rows, Cols int
	Task       [][]int
	Cells      [][]*Cell
}

func (g *Game) InBounds(i, j int) bool {
	return i >= 0 && j >= 0 && i < g.Rows && j < g.Cols
}

func (g *Game) ColourInRow(c Colour, row int) int {
	count := 0
	for _, cell := range g.GetRow(row) {
		if cell.Colour == c {
			count++
		}
	}
	return count
}

func (g *Game) ColourInCol(c Colour, col int) int {
	count := 0
	for _, cell := range g.GetCol(col) {
		if cell.Colour == c {
			count++
		}
	}
	return count
}

func (g *Game) GetRow(i int) []*Cell {
	return g.Cells[i]
}

func (g *Game) GetCol(j int) []*Cell {
	cells := []*Cell{}
	for _, row := range g.Cells {
		cells = append(cells, row[j])
	}
	return cells
}

type Colour int

const (
	Blank Colour = -1
	White Colour = 0
	Black Colour = 1
)

type Cell struct {
	Game     *Game
	Row, Col int
	Colour   Colour
}

func (g *Game) NewCell(i, j int) *Cell {
	c := &Cell{
		Game:   g,
		Row:    i,
		Col:    j,
		Colour: Blank,
	}
	if !g.InBounds(i, j) {
		return c
	}
	c.Colour = Colour(g.Task[i][j])
	return c
}

func (g *Game) GetCell(i, j int) *Cell {
	if !g.InBounds(i, j) {
		return g.NewCell(i, j)
	}
	return g.Cells[i][j]
}

func (c *Cell) GetValue() int {
	switch c.Colour {
	case Blank:
		return 0
	case Black:
		return 1
	case White:
		return 2
	}
	return -1
}

func (c *Cell) Neighbour(direction brainless.Direction) *Cell {
	i, j := -1, -1
	switch direction {
	case brainless.East:
		i, j = c.Row, c.Col+1
	case brainless.North:
		i, j = c.Row-1, c.Col
	case brainless.West:
		i, j = c.Row, c.Col-1
	case brainless.South:
		i, j = c.Row+1, c.Col
	}
	return c.Game.GetCell(i, j)
}

type Tile struct {
	Black Colour
	White Colour
}

func NewTile() *Tile {
	return &Tile{
		Black: Black,
		White: White,
	}
}

func (t *Tile) Flip() brainless.Flipper {
	t.Black = t.Black ^ 1
	t.White = t.White ^ 1
	return t
}
