package brainless

type Orientator interface {
	Rotate(int) Orientator
	Flip(Direction) Orientator
	Orientate(Orientation) Orientator
}

type Compass struct {
	East, North, West, South Direction
}

func NewCompass() Compass {
	return Compass{
		East:  East,
		North: North,
		West:  West,
		South: South,
	}
}

func (c Compass) Rotate(times int) Orientator {
	c.East = c.East.Rotate(times)
	c.North = c.North.Rotate(times)
	c.West = c.West.Rotate(times)
	c.South = c.South.Rotate(times)
	return c
}

func (c Compass) Flip(along Direction) Orientator {
	switch along {
	case East, West:
		c.North = c.North.Flip()
		c.South = c.South.Flip()
	case North, South:
		c.East = c.East.Flip()
		c.West = c.West.Flip()
	}
	return c
}

func (c Compass) Orientate(orientation Orientation) Orientator {
	c = c.Rotate(orientation.GetRotations()).(Compass)
	if orientation.ShouldFlip() {
		c = c.Flip(orientation.GetFlipDirection()).(Compass)
	}
	return c
}

type Orientation int

func (o Orientation) GetRotations() int {
	return int(o >> 1)
}

func (o Orientation) ShouldFlip() bool {
	return int(o)&1 == 1
}

func (o Orientation) GetFlipDirection() Direction {
	return 1 << (o >> 1)
}
