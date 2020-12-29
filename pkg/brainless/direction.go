package brainless

type Direction int

const (
	East Direction = 1 << iota
	North
	West
	South
)

func (d Direction) Rotate(times int) Direction {
	return Direction((d<<times + d<<times>>4) & 15)
}

func (d Direction) Flip() Direction {
	return d.Rotate(2)
}
