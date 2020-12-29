package pipes

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/puzzle-solvers/brainless/pkg/brainless"
)

func Solve(w http.ResponseWriter, r *http.Request) {
	game := brainless.NewGame(w, r, brainless.Config{
		PuzzleURL:   "https://www.puzzle-pipes.com",
		FunctionURL: "http://localhost:8080/solve",
		//FunctionURL: "https://australia-southeast1-puzzle-solvers.cloudfunctions.net/pipes",
		//FunctionURL: "https://australia-southeast1-dylanbutler.cloudfunctions.net/solve-pipes",
	})
	brain := &PipeBrain{}
	game.ConnectBrain(brain)
	game.HandleRequest()
}

type Option func()

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

	g.Locked = make([][]bool, rows)
	for i := range g.Locked {
		g.Locked[i] = make([]bool, cols)
		for j := range g.Locked[i] {
			g.Locked[i][j] = false
		}
	}

	g.Rotations = make([][]int, rows)
	for i := range g.Rotations {
		g.Rotations[i] = make([]int, cols)
		for j := range g.Rotations[i] {
			g.Rotations[i][j] = 0
		}
	}

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

			if cell.IsLocked() {
				continue
			}

			shape := Shape(cell.Val())

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

func (p *PipeBrain) ToResponse() interface{} {
	return Response{
		CellStatus: p.Game.CellStatus(),
		Pinned:     p.Game.Pinned(),
	}
}

type Response struct {
	CellStatus [][]int `json:"cellStatus"`
	Pinned     [][]int `json:"pinned"`
}

type Game struct {
	Rows, Cols int
	Task       [][]int
	Locked     [][]bool
	Rotations  [][]int
	Cells      [][]*Cell
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
	return s.rotate(times)
}

func (s Shape) rotate(times int) Shape {
	return Shape((s<<times + s<<times>>4) & 15)
}

func (s Shape) Flip(along brainless.Direction) brainless.Orientator {
	return s.flip(along)
}

func (s Shape) flip(along brainless.Direction) Shape {
	axis := int(along + along.Flip())
	return Shape(axis&int(s) + ^axis&int(s.rotate(2)))
}

func (s Shape) Orientate(o brainless.Orientation) brainless.Orientator {
	return s.orientate(o)
}

func (s Shape) orientate(o brainless.Orientation) brainless.Orientator {
	s = s.rotate(o.GetRotations())
	if o.ShouldFlip() {
		s = s.flip(o.GetFlipDirection())
	}
	return s
}

// Rotates anti-clockwise
func (g *Game) Rotate(i, j, times int) {
	if !g.InBounds(i, j) {
		return
	}
	if !g.Locked[i][j] {
		g.Rotations[i][j] = (g.Rotations[i][j] + times) & 3
	}
}

func (g *Game) Faces(i, j int, direction brainless.Direction) bool {
	if !g.InBounds(i, j) {
		return false
	}
	return ((g.Task[i][j] << g.Rotations[i][j]) & (int(direction) + int(direction)<<4)) > 0
}

func (g *Game) Lock(i, j int) {
	if !g.InBounds(i, j) {
		return
	}
	g.Locked[i][j] = true
}

func (g *Game) IsLocked(i, j int) bool {
	if !g.InBounds(i, j) {
		return true
	}
	return g.Locked[i][j]
}

func (g *Game) Val(i, j int) int {
	if !g.InBounds(i, j) {
		return -1
	}
	return g.Task[i][j]
}

func (g *Game) InBounds(i, j int) bool {
	return i >= 0 && j >= 0 && i < g.Rows && j < g.Cols
}

func (g *Game) CellStatus() [][]int {
	return g.Rotations
}

func (g *Game) Pinned() [][]int {
	p := make([][]int, g.Rows)
	for i := range p {
		p[i] = make([]int, g.Cols)
		for j := range p[i] {
			p[i][j] = 0
			if g.IsLocked(i, j) {
				p[i][j] = 1
			}
		}
	}
	return p
}

func (g *Game) IsComplete() bool {
	for _, row := range g.Locked {
		for _, cell := range row {
			if !cell {
				return false
			}
		}
	}
	return true
}

type Cell struct {
	Game     *Game
	Row, Col int
}

func (g *Game) NewCell(i, j int) *Cell {
	return &Cell{
		Game: g,
		Row:  i,
		Col:  j,
	}
}

func (g *Game) CellAt(i, j int) *Cell {
	if !g.InBounds(i, j) {
		return g.NewCell(i, j)
	}
	return g.Cells[i][j]
}

func (c *Cell) Rotate(times int) {
	c.Game.Rotate(c.Row, c.Col, times)
}

func (c *Cell) Faces(direction brainless.Direction) bool {
	return c.Game.Faces(c.Row, c.Col, direction)
}

func (c *Cell) Lock() {
	c.Game.Lock(c.Row, c.Col)
}

func (c *Cell) IsLocked() bool {
	return c.Game.IsLocked(c.Row, c.Col)
}

func (c *Cell) MustFace(direction brainless.Direction) bool {
	return c.IsLocked() && c.Faces(direction)
}

func (c *Cell) CantFace(direction brainless.Direction) bool {
	return c.IsLocked() && !c.Faces(direction)
}

func (c *Cell) Val() int {
	return c.Game.Val(c.Row, c.Col)
}

func (c *Cell) CanLock(direction brainless.Direction) bool {
	return !(c.Neighbour(direction).IsLocked() && c.Neighbour(direction).Faces(direction.Rotate(2)) != c.Faces(direction))
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

func dump(w ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if w == nil {
		fmt.Println()
	} else {
		fmt.Printf("Dumped at %s:%d %v\n", file, line, w)
	}
}

func bin(i int) string {
	return fmt.Sprintf("%04s", strconv.FormatInt(int64(i), 2))
}
