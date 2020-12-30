package pipes

import (
	"net/http"

	"github.com/puzzle-solvers/brainless/pkg/brainless"
)

func Solve(w http.ResponseWriter, r *http.Request) {
	game := brainless.NewGame(w, r, brainless.Config{
		PuzzleURL:   "https://www.puzzle-pipes.com",
		FunctionURL: "http://localhost:8080/pipes",
		//FunctionURL: "https://australia-southeast1-puzzle-solvers.cloudfunctions.net/pipes",
		//FunctionURL: "https://australia-southeast1-dylanbutler.cloudfunctions.net/solve-pipes",
	})
	brain := &PipeBrain{}
	game.ConnectBrain(brain)
	game.HandleRequest()
}

type PipeBrain struct {
	Game *Game
	task [][]int
}

func (p *PipeBrain) Setup() {
	g := &Game{}

	if len(p.task) == 0 {
		return
	}

	if len(p.task[0]) == 0 {
		return
	}

	rows, cols := len(p.task), len(p.task[0])

	g.Rows, g.Cols = rows, cols

	g.Task = p.task

	g.Cells = make([][]*Cell, rows)
	for i := range g.Cells {
		g.Cells[i] = make([]*Cell, cols)
		for j := range g.Cells[i] {
			g.Cells[i][j] = g.NewCell(i, j)
		}
	}

	p.Game = g
}

func (p *PipeBrain) Step() {
	for _, row := range p.Game.Cells {
		for _, cell := range row {

			if cell.Locked {
				continue
			}

			shape := Shape(cell.Value)

			brainless.OrientatorRule(shape, P, func(c brainless.Compass, r int) {
				if cell.Neighbour(c.East).MustFace(c.West) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.North).CantFace(c.South) && cell.Neighbour(c.West).CantFace(c.East) && cell.Neighbour(c.South).CantFace(c.North) {
					cell.Rotate(r)
					cell.Lock()
				}
			}, ShapeChecker)

			brainless.OrientatorRule(shape, L, func(c brainless.Compass, r int) {
				if cell.Neighbour(c.East).MustFace(c.West) && cell.Neighbour(c.North).MustFace(c.South) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.North).MustFace(c.South) && cell.Neighbour(c.West).CantFace(c.East) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.West).CantFace(c.East) && cell.Neighbour(c.South).CantFace(c.North) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.South).CantFace(c.North) && cell.Neighbour(c.East).MustFace(c.West) {
					cell.Rotate(r)
					cell.Lock()
				}
			}, ShapeChecker)

			brainless.OrientatorRule(shape, I, func(c brainless.Compass, r int) {
				if cell.Neighbour(c.East).MustFace(c.West) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.North).CantFace(c.South) {
					cell.Rotate(r)
					cell.Lock()
				}
			}, ShapeChecker)

			brainless.OrientatorRule(shape, T, func(c brainless.Compass, r int) {
				if cell.Neighbour(c.East).MustFace(c.West) && cell.Neighbour(c.North).MustFace(c.South) && cell.Neighbour(c.West).MustFace(c.East) {
					cell.Rotate(r)
					cell.Lock()
				}

				if cell.Neighbour(c.South).CantFace(c.North) {
					cell.Rotate(r)
					cell.Lock()
				}
			}, ShapeChecker)
		}
	}
}

func (p *PipeBrain) CheckDone() bool {
	return p.Game.IsComplete()
}

func (p *PipeBrain) GetTask(task [][]int) {
	p.task = task
}

type Response struct {
	CellStatus [][]int `json:"cellStatus"`
	Pinned     [][]int `json:"pinned"`
}

func (p *PipeBrain) ToResponse() interface{} {
	return Response{
		CellStatus: p.Game.CellStatus(),
		Pinned:     p.Game.Pinned(),
	}
}

type Shape int

const (
	P Shape = 1
	L Shape = 3
	I Shape = 5
	T Shape = 7
)

func ShapeChecker(a, b brainless.Orientator) bool {
	return a.(Shape) == b.(Shape)
}

func (s Shape) Rotate(times int) brainless.Orientator {
	return Shape((s<<times + s<<times>>4) & 15)
}

func (s Shape) Flip(along brainless.Direction) brainless.Orientator {
	axis := int(along + along.Flip())
	return Shape(axis&int(s) + ^axis&int(s.Rotate(2).(Shape)))
}

func (s Shape) Orientate(o brainless.Orientation) brainless.Orientator {
	s = s.Rotate(o.GetRotations()).(Shape)
	if o.ShouldFlip() {
		s = s.Flip(o.GetFlipDirection()).(Shape)
	}
	return s
}

type Game struct {
	Rows, Cols int
	Task       [][]int
	Cells      [][]*Cell
}

func (g *Game) InBounds(i, j int) bool {
	return i >= 0 && j >= 0 && i < g.Rows && j < g.Cols
}

func (g *Game) CellStatus() [][]int {
	c := make([][]int, g.Rows)
	for i := range c {
		c[i] = make([]int, g.Cols)
		for j := range c[i] {
			c[i][j] = g.Cells[i][j].Rotations
		}
	}
	return c
}

func (g *Game) Pinned() [][]int {
	p := make([][]int, g.Rows)
	for i := range p {
		p[i] = make([]int, g.Cols)
		for j := range p[i] {
			p[i][j] = 0
			if g.Cells[i][j].Locked {
				p[i][j] = 1
			}
		}
	}
	return p
}

func (g *Game) IsComplete() bool {
	for _, row := range g.Cells {
		for _, cell := range row {
			if !cell.Locked {
				return false
			}
		}
	}
	return true
}

type Cell struct {
	Game      *Game
	Row, Col  int
	Value     int
	Rotations int
	Locked    bool
}

func (g *Game) NewCell(i, j int) *Cell {
	c := &Cell{
		Game:      g,
		Row:       i,
		Col:       j,
		Rotations: 0,
	}
	if !g.InBounds(i, j) {
		c.Value = 0
		c.Locked = true
		return c
	}
	c.Value = g.Task[i][j]
	c.Locked = false
	return c
}

func (g *Game) CellAt(i, j int) *Cell {
	if !g.InBounds(i, j) {
		return g.NewCell(i, j)
	}
	return g.Cells[i][j]
}

func (c *Cell) Rotate(times int) {
	if !c.Locked {
		c.Rotations = (c.Rotations + times) & 3
	}
}

func (c *Cell) Lock() {
	c.Locked = true
}

func (c *Cell) Faces(direction brainless.Direction) bool {
	return ((c.Value << c.Rotations) & (int(direction) + int(direction)<<4)) > 0
}

func (c *Cell) MustFace(direction brainless.Direction) bool {
	return c.Locked && c.Faces(direction)
}

func (c *Cell) CantFace(direction brainless.Direction) bool {
	return c.Locked && !c.Faces(direction)
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
	return c.Game.CellAt(i, j)
}
